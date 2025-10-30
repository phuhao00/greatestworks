package simclient

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Duration wraps time.Duration to support YAML unmarshalling from strings like "5s".
type Duration struct {
	time.Duration
}

// NewDuration constructs a Duration from a time.Duration value.
func NewDuration(d time.Duration) Duration {
	return Duration{Duration: d}
}

// UnmarshalYAML parses a duration value expressed as a string.
func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	if value == nil {
		d.Duration = 0
		return nil
	}

	switch value.Kind {
	case yaml.ScalarNode:
		var raw string
		if err := value.Decode(&raw); err != nil {
			return fmt.Errorf("invalid duration value: %w", err)
		}
		if raw == "" {
			d.Duration = 0
			return nil
		}
		parsed, err := time.ParseDuration(raw)
		if err != nil {
			return fmt.Errorf("could not parse duration %q: %w", raw, err)
		}
		d.Duration = parsed
		return nil
	default:
		return fmt.Errorf("unsupported YAML node kind for duration: %v", value.Kind)
	}
}

// MarshalYAML renders the duration as a human-readable string.
func (d Duration) MarshalYAML() (interface{}, error) {
	return d.Duration.String(), nil
}

// AsDuration returns the embedded time.Duration, guarding zero defaults.
func (d Duration) AsDuration() time.Duration {
	return d.Duration
}

// Config contains runtime options for the simulator client.
type Config struct {
	Auth     AuthConfig     `yaml:"auth"`
	Gateway  GatewayConfig  `yaml:"gateway"`
	Scenario ScenarioConfig `yaml:"scenario"`
	Load     LoadConfig     `yaml:"load"`
	Metrics  MetricsConfig  `yaml:"metrics"`
}

// AuthConfig controls optional authentication against the auth service.
type AuthConfig struct {
	Enabled   bool     `yaml:"enabled"`
	BaseURL   string   `yaml:"base_url"`
	LoginPath string   `yaml:"login_path"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password"`
	Timeout   Duration `yaml:"timeout"`
}

// GatewayConfig holds TCP gateway connection information.
type GatewayConfig struct {
	Host           string   `yaml:"host"`
	Port           int      `yaml:"port"`
	ConnectTimeout Duration `yaml:"connect_timeout"`
	ReadTimeout    Duration `yaml:"read_timeout"`
	WriteTimeout   Duration `yaml:"write_timeout"`
}

// ScenarioConfig tunes how long and how frequently simulated actions run.
type ScenarioConfig struct {
	Name           string                 `yaml:"name"`
	Type           string                 `yaml:"type"`
	Duration       Duration               `yaml:"duration"`
	ActionInterval Duration               `yaml:"action_interval"`
	PlayerPrefix   string                 `yaml:"player_prefix"`
	StopOnError    bool                   `yaml:"stop_on_error"`
	Features       []string               `yaml:"features"`
	Actions        []ScenarioActionConfig `yaml:"actions"`
}

// ScenarioActionConfig represents a single action step in a feature scenario.
type ScenarioActionConfig struct {
	Name           string   `yaml:"name"`
	Message        string   `yaml:"message"`
	Flags          []string `yaml:"flags"`
	ExpectResponse *bool    `yaml:"expect_response"`
	Pause          Duration `yaml:"pause"`
	Repeat         int      `yaml:"repeat"`
}

// LoadConfig enables multi-user execution for pressure testing.
type LoadConfig struct {
	Enabled      bool     `yaml:"enabled"`
	VirtualUsers int      `yaml:"virtual_users"`
	Concurrency  int      `yaml:"concurrency"`
	RampUp       Duration `yaml:"ramp_up"`
	Iterations   int      `yaml:"iterations"`
	StopOnError  bool     `yaml:"stop_on_error"`
}

// MetricsConfig configures how metrics are reported.
type MetricsConfig struct {
	Output string `yaml:"output"`
}

// DefaultConfig returns a config pre-populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Auth: AuthConfig{
			Enabled:   false,
			BaseURL:   "http://localhost:8080",
			LoginPath: "/api/v1/auth/login",
			Username:  "tester",
			Password:  "tester123",
			Timeout:   NewDuration(5 * time.Second),
		},
		Gateway: GatewayConfig{
			Host:           "127.0.0.1",
			Port:           9090,
			ConnectTimeout: NewDuration(3 * time.Second),
			ReadTimeout:    NewDuration(2 * time.Second),
			WriteTimeout:   NewDuration(2 * time.Second),
		},
		Scenario: ScenarioConfig{
			Name:           "basic",
			Type:           "basic",
			Duration:       NewDuration(10 * time.Second),
			ActionInterval: NewDuration(1 * time.Second),
			PlayerPrefix:   "player",
			StopOnError:    true,
			Features:       []string{},
			Actions:        nil,
		},
		Load: LoadConfig{
			Enabled:      false,
			VirtualUsers: 50,
			Concurrency:  10,
			RampUp:       NewDuration(5 * time.Second),
			Iterations:   1,
			StopOnError:  false,
		},
		Metrics: MetricsConfig{
			Output: "console",
		},
	}
}

// LoadConfigFromFile loads a configuration file and merges it with defaults.
func LoadConfigFromFile(path string) (Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		cfg.Normalize()
		return cfg, nil
	}

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return cfg, fmt.Errorf("read config: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config: %w", err)
	}

	cfg.Normalize()
	return cfg, nil
}

// Normalize ensures required fields fall back to defaults when unset.
func (c *Config) Normalize() {
	if c.Gateway.Host == "" {
		c.Gateway.Host = "127.0.0.1"
	}
	if c.Gateway.Port == 0 {
		c.Gateway.Port = 9090
	}
	if c.Gateway.ConnectTimeout.AsDuration() == 0 {
		c.Gateway.ConnectTimeout = NewDuration(3 * time.Second)
	}
	if c.Gateway.ReadTimeout.AsDuration() == 0 {
		c.Gateway.ReadTimeout = NewDuration(2 * time.Second)
	}
	if c.Gateway.WriteTimeout.AsDuration() == 0 {
		c.Gateway.WriteTimeout = NewDuration(2 * time.Second)
	}

	if c.Scenario.Name == "" {
		c.Scenario.Name = "basic"
	}
	if c.Scenario.Type == "" {
		c.Scenario.Type = "basic"
	}
	if c.Scenario.Duration.AsDuration() == 0 {
		c.Scenario.Duration = NewDuration(10 * time.Second)
	}
	if c.Scenario.ActionInterval.AsDuration() == 0 {
		c.Scenario.ActionInterval = NewDuration(1 * time.Second)
	}
	if c.Scenario.PlayerPrefix == "" {
		c.Scenario.PlayerPrefix = "player"
	}
	if c.Scenario.Features == nil {
		c.Scenario.Features = []string{}
	}

	if c.Load.Concurrency <= 0 {
		c.Load.Concurrency = 10
	}
	if c.Load.VirtualUsers <= 0 {
		c.Load.VirtualUsers = c.Load.Concurrency
	}
	if c.Load.Iterations <= 0 {
		c.Load.Iterations = 1
	}
	if c.Load.RampUp.AsDuration() < 0 {
		c.Load.RampUp = NewDuration(0)
	}
}
