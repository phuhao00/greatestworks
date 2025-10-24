package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoaderMergesFilesAndAppliesDefaults(t *testing.T) {
	dir := t.TempDir()

	base := `app:
  name: GreatestWorks Test
  version: 0.1.0
logging:
  level: warn
`

	service := `server:
  http:
    port: 9090
  rpc:
    port: 19080
`

	if err := os.WriteFile(filepath.Join(dir, "config.base.yaml"), []byte(base), 0o644); err != nil {
		t.Fatalf("write base config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "game-service.yaml"), []byte(service), 0o644); err != nil {
		t.Fatalf("write service config: %v", err)
	}

	loader := NewLoader(
		WithBaseDir(dir),
		WithService("game-service"),
		WithEnvironment("development"),
	)

	cfg, files, err := loader.Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 config files to be used, got %d", len(files))
	}

	if cfg.App.Name != "GreatestWorks Test" {
		t.Fatalf("unexpected app name: %s", cfg.App.Name)
	}

	if cfg.Server.HTTP.Port != 9090 {
		t.Fatalf("expected HTTP port override to be 9090, got %d", cfg.Server.HTTP.Port)
	}

	if cfg.Server.RPC.Port != 19080 {
		t.Fatalf("expected RPC port override to be 19080, got %d", cfg.Server.RPC.Port)
	}

	if cfg.Logging.Format != "json" {
		t.Fatalf("expected logging format default json, got %s", cfg.Logging.Format)
	}

	if cfg.Security.JWT.Secret == "" {
		t.Fatalf("expected JWT secret to have default value")
	}
}

func TestLoaderRespectsEnvironmentOverrides(t *testing.T) {
	dir := t.TempDir()

	fileContent := `server:
  http:
    port: 9090
`
	if err := os.WriteFile(filepath.Join(dir, "game-service.yaml"), []byte(fileContent), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Setenv("SERVER_HTTP_PORT", "8088")
	t.Setenv("MONGODB_URI", "mongodb://example:27017")

	loader := NewLoader(
		WithBaseDir(dir),
		WithService("game-service"),
	)

	cfg, _, err := loader.Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.Server.HTTP.Port != 8088 {
		t.Fatalf("expected env override for HTTP port, got %d", cfg.Server.HTTP.Port)
	}

	if cfg.Database.MongoDB.URI != "mongodb://example:27017" {
		t.Fatalf("expected env override for mongo uri, got %s", cfg.Database.MongoDB.URI)
	}
}
