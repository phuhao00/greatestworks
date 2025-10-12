package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Loader merges configuration files for a service.
type Loader struct {
	baseDir       string
	env           string
	service       string
	explicitFiles []string
}

// Option customises a Loader instance.
type Option func(*Loader)

// NewLoader constructs a Loader with optional overrides.
func NewLoader(opts ...Option) *Loader {
	loader := &Loader{
		baseDir: "configs",
		env:     normalizeEnv(os.Getenv("APP_ENV")),
	}

	if loader.env == "" {
		loader.env = "development"
	}

	for _, opt := range opts {
		opt(loader)
	}

	if loader.service == "" {
		loader.service = strings.TrimSpace(os.Getenv("SERVICE_NAME"))
	}

	if path := strings.TrimSpace(os.Getenv("CONFIG_PATH")); path != "" {
		loader.explicitFiles = []string{path}
	} else if path := strings.TrimSpace(os.Getenv("CONFIG_FILE")); path != "" {
		loader.explicitFiles = []string{path}
	}

	if loader.baseDir == "" {
		loader.baseDir = "."
	}

	return loader
}

// WithBaseDir sets the base search directory for configuration files.
func WithBaseDir(dir string) Option {
	return func(l *Loader) {
		if dir == "" {
			return
		}
		l.baseDir = filepath.Clean(dir)
	}
}

// WithEnvironment sets the active environment (e.g. development, production).
func WithEnvironment(env string) Option {
	return func(l *Loader) {
		l.env = normalizeEnv(env)
	}
}

// WithService sets the logical service name whose configuration should be loaded.
func WithService(service string) Option {
	return func(l *Loader) {
		l.service = strings.TrimSpace(service)
	}
}

// WithExplicitFiles supplies explicit configuration files and bypasses discovery.
func WithExplicitFiles(files ...string) Option {
	return func(l *Loader) {
		cleaned := make([]string, 0, len(files))
		for _, file := range files {
			file = strings.TrimSpace(file)
			if file != "" {
				cleaned = append(cleaned, file)
			}
		}
		if len(cleaned) > 0 {
			l.explicitFiles = cleaned
		}
	}
}

// Load gathers and merges configuration files into a Config instance.
func (l *Loader) Load() (*Config, []string, error) {
	candidates := l.resolveCandidates()
	if len(candidates) == 0 {
		return nil, nil, fmt.Errorf("config: no configuration candidates resolved")
	}

	var (
		cfg  Config
		used []string
		errs []error
	)

	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			errs = append(errs, fmt.Errorf("config: read %s: %w", path, err))
			continue
		}

		if err := yaml.Unmarshal(data, &cfg); err != nil {
			errs = append(errs, fmt.Errorf("config: parse %s: %w", path, err))
			continue
		}

		used = append(used, path)
	}

	if len(used) == 0 {
		if len(errs) > 0 {
			return nil, nil, errors.Join(errs...)
		}
		return nil, nil, fmt.Errorf(
			"config: no configuration files found for service %q (env=%s) under %s",
			l.service,
			l.env,
			l.baseDir,
		)
	}

	cfg.ApplyDefaults()
	l.applyEnvOverrides(&cfg)

	if l.service != "" && cfg.Service.Name == "" {
		cfg.Service.Name = l.service
	}
	if cfg.App.Environment == "" {
		cfg.App.Environment = l.env
	}

	if err := cfg.Validate(); err != nil {
		return nil, nil, err
	}

	return &cfg, used, nil
}

// LoadInto hydrates the supplied target structure with the merged configuration.
func (l *Loader) LoadInto(target any) ([]string, error) {
	cfg, files, err := l.Load()
	if err != nil {
		return nil, err
	}

	// Marshal then unmarshal to allow the caller to provide a tailored struct shape.
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("config: marshal combined config: %w", err)
	}

	if err := yaml.Unmarshal(data, target); err != nil {
		return nil, fmt.Errorf("config: hydrate target: %w", err)
	}

	return files, nil
}

// Environment exposes the current loader environment value.
func (l *Loader) Environment() string {
	return l.env
}

// Service exposes the service name for the loader.
func (l *Loader) Service() string {
	return l.service
}

// BaseDir exposes the root search directory.
func (l *Loader) BaseDir() string {
	return l.baseDir
}

func (l *Loader) resolveCandidates() []string {
	if len(l.explicitFiles) > 0 {
		return normalizePaths(l.explicitFiles)
	}

	names := []string{
		"config.base.yaml",
		"config.yaml",
	}

	if l.env != "" {
		names = append(names, fmt.Sprintf("config.%s.yaml", l.env))
	}

	if l.service != "" {
		names = append(names, fmt.Sprintf("%s.yaml", l.service))
		if l.env != "" {
			names = append(names, fmt.Sprintf("%s.%s.yaml", l.service, l.env))
		}
	}

	paths := make([]string, 0, len(names))
	for _, name := range names {
		if name == "" {
			continue
		}
		paths = append(paths, filepath.Join(l.baseDir, name))
	}

	return normalizePaths(paths)
}

func (l *Loader) applyEnvOverrides(cfg *Config) {
	overrideString := func(target *string, keys ...string) {
		for _, key := range keys {
			if value := strings.TrimSpace(os.Getenv(key)); value != "" {
				*target = value
				return
			}
		}
	}

	overrideInt := func(target *int, keys ...string) {
		for _, key := range keys {
			if value := strings.TrimSpace(os.Getenv(key)); value != "" {
				if v, err := strconv.Atoi(value); err == nil {
					*target = v
					return
				}
			}
		}
	}

	overrideString(&cfg.Server.HTTP.Host, "SERVER_HTTP_HOST", "SERVER_HOST")
	overrideInt(&cfg.Server.HTTP.Port, "SERVER_HTTP_PORT", "SERVER_PORT")

	overrideString(&cfg.Database.MongoDB.URI, "MONGODB_URI")
	overrideString(&cfg.Database.MongoDB.Database, "MONGODB_DATABASE")

	overrideString(&cfg.Database.Redis.Addr, "REDIS_ADDR")
	overrideString(&cfg.Database.Redis.Password, "REDIS_PASSWORD")

	overrideString(&cfg.Security.JWT.Secret, "JWT_SECRET")

	overrideString(&cfg.Logging.Level, "LOG_LEVEL")

	overrideString(&cfg.Messaging.NATS.URL, "NATS_URL")
	overrideString(&cfg.Messaging.NATS.ClusterID, "NATS_CLUSTER_ID")
	overrideString(&cfg.Messaging.NATS.ClientID, "NATS_CLIENT_ID")

	overrideString(&cfg.Service.NodeID, "SERVICE_NODE_ID", "POD_NAME")
}

func normalizePaths(paths []string) []string {
	seen := make(map[string]struct{}, len(paths))
	normalized := make([]string, 0, len(paths))

	for _, path := range paths {
		if path == "" {
			continue
		}
		cleaned := filepath.Clean(path)
		if !filepath.IsAbs(cleaned) {
			abs, err := filepath.Abs(cleaned)
			if err == nil {
				cleaned = abs
			}
		}
		if _, exists := seen[cleaned]; exists {
			continue
		}
		seen[cleaned] = struct{}{}
		normalized = append(normalized, cleaned)
	}

	sort.Strings(normalized)
	return normalized
}

func normalizeEnv(env string) string {
	return strings.ToLower(strings.TrimSpace(env))
}
