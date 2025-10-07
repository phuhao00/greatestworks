// Package container provides service providers for dependency injection
package container

import (
	"fmt"
	"greatestworks/internal/infrastructure/config"
	"greatestworks/internal/infrastructure/logging"
	"greatestworks/internal/infrastructure/monitoring"
	"greatestworks/internal/infrastructure/persistence"
	"greatestworks/internal/infrastructure/protocol"
	"greatestworks/internal/infrastructure/weave"
)

// ConfigProvider provides configuration services
type ConfigProvider struct {
	configPath string
}

// NewConfigProvider creates a new config provider
func NewConfigProvider(configPath string) *ConfigProvider {
	return &ConfigProvider{
		configPath: configPath,
	}
}

// RegisterServices registers configuration services
func (cp *ConfigProvider) RegisterServices(container *Container) error {
	// Register config loader
	container.RegisterSingleton("config.loader", func(c *Container) (interface{}, error) {
		return config.NewConfigLoader(""), nil
	})

	// Register environment manager
	container.RegisterSingleton("config.env_manager", func(c *Container) (interface{}, error) {
		return config.GetEnvManager(), nil
	})

	// Register main configuration
	container.RegisterSingleton("config", func(c *Container) (interface{}, error) {
		loader, err := c.Resolve("config.loader")
		if err != nil {
			return nil, err
		}

		configLoader := loader.(*config.ConfigLoader)
		if cp.configPath != "" {
			return configLoader.LoadFromFile(cp.configPath)
		}

		envManager, err := c.Resolve("config.env_manager")
		if err != nil {
			return nil, err
		}

		envMgr := envManager.(*config.EnvManager)
		return configLoader.LoadFromFile(config.GetConfigPath())
	}, "config.loader", "config.env_manager")

	// Register hot reload manager
	// TODO: 实现热重载功能
	container.RegisterSingleton("config.hot_reload", func(c *Container) (interface{}, error) {
		return nil, fmt.Errorf("hot reload not implemented")
	}, "config")

	return nil
}

// LoggingProvider provides logging services
type LoggingProvider struct{}

// NewLoggingProvider creates a new logging provider
func NewLoggingProvider() *LoggingProvider {
	return &LoggingProvider{}
}

// RegisterServices registers logging services
func (lp *LoggingProvider) RegisterServices(container *Container) error {
	// Register logger
	container.RegisterSingleton("logger", func(c *Container) (interface{}, error) {
		cfg, err := c.Resolve("config")
		if err != nil {
			return nil, err
		}

		config := cfg.(*config.Config)
		logConfig := config.Logging

		// TODO: 实现日志配置
		return nil, fmt.Errorf("logging not implemented")
	}, "config")

	// Register HTTP middleware
	container.RegisterTransient("logging.http_middleware", func(c *Container) (interface{}, error) {
		logger, err := c.Resolve("logger")
		if err != nil {
			return nil, err
		}

		cfg, err := c.Resolve("config")
		if err != nil {
			return nil, err
		}

		config := cfg.(*config.Config)
		return logging.NewHTTPMiddleware(logger.(logging.Logger), &config.Logging.Middleware), nil
	}, "logger", "config")

	// Register Game middleware
	container.RegisterTransient("logging.game_middleware", func(c *Container) (interface{}, error) {
		logger, err := c.Resolve("logger")
		if err != nil {
			return nil, err
		}

		cfg, err := c.Resolve("config")
		if err != nil {
			return nil, err
		}

		config := cfg.(*config.Config)
		return logging.NewGameMiddleware(logger.(logging.Logger), &config.Logging.Middleware), nil
	}, "logger", "config")

	return nil
}

// MonitoringProvider provides monitoring services
type MonitoringProvider struct{}

// NewMonitoringProvider creates a new monitoring provider
func NewMonitoringProvider() *MonitoringProvider {
	return &MonitoringProvider{}
}

// RegisterServices registers monitoring services
func (mp *MonitoringProvider) RegisterServices(container *Container) error {
	// Register metrics registry
	container.RegisterSingleton("metrics.registry", func(c *Container) (interface{}, error) {
		cfg, err := c.Resolve("config")
		if err != nil {
			return nil, err
		}

		config := cfg.(*config.Config)
		return monitoring.NewPrometheusRegistry(&config.Monitoring)
	}, "config")

	// Register metrics factory
	container.RegisterSingleton("metrics.factory", func(c *Container) (interface{}, error) {
		registry, err := c.Resolve("metrics.registry")
		if err != nil {
			return nil, err
		}

		return monitoring.NewPrometheusFactory(registry.(monitoring.MetricsRegistry)), nil
	}, "metrics.registry")

	// Register collectors
	container.RegisterSingleton("metrics.system_collector", func(c *Container) (interface{}, error) {
		factory, err := c.Resolve("metrics.factory")
		if err != nil {
			return nil, err
		}

		return monitoring.NewSystemCollector(factory.(monitoring.MetricsFactory)), nil
	}, "metrics.factory")

	container.RegisterSingleton("metrics.game_collector", func(c *Container) (interface{}, error) {
		factory, err := c.Resolve("metrics.factory")
		if err != nil {
			return nil, err
		}

		return monitoring.NewGameCollector(factory.(monitoring.MetricsFactory)), nil
	}, "metrics.factory")

	// Register metrics server
	container.RegisterSingleton("metrics.server", func(c *Container) (interface{}, error) {
		registry, err := c.Resolve("metrics.registry")
		if err != nil {
			return nil, err
		}

		cfg, err := c.Resolve("config")
		if err != nil {
			return nil, err
		}

		config := cfg.(*config.Config)
		return monitoring.NewPrometheusServer(registry.(monitoring.MetricsRegistry), &config.Monitoring), nil
	}, "metrics.registry", "config")

	return nil
}

// PersistenceProvider provides persistence services
type PersistenceProvider struct{}

// NewPersistenceProvider creates a new persistence provider
func NewPersistenceProvider() *PersistenceProvider {
	return &PersistenceProvider{}
}

// RegisterServices registers persistence services
func (pp *PersistenceProvider) RegisterServices(container *Container) error {
	// Register MongoDB
	container.RegisterSingleton("database.mongodb", func(c *Container) (interface{}, error) {
		cfg, err := c.Resolve("config")
		if err != nil {
			return nil, err
		}

		config := cfg.(*config.Config)
		return persistence.NewMongoDB(&config.Database.MongoDB)
	}, "config")

	// Register Redis
	container.RegisterSingleton("database.redis", func(c *Container) (interface{}, error) {
		cfg, err := c.Resolve("config")
		if err != nil {
			return nil, err
		}

		config := cfg.(*config.Config)
		return persistence.NewRedisClient(&config.Cache.Redis)
	}, "config")

	// Register repositories
	container.RegisterScoped("repository.player", func(c *Container) (interface{}, error) {
		mongoDB, err := c.Resolve("database.mongodb")
		if err != nil {
			return nil, err
		}

		return persistence.NewPlayerRepository(mongoDB.(*persistence.MongoDB)), nil
	}, "database.mongodb")

	return nil
}

// ProtocolProvider provides protocol services
type ProtocolProvider struct{}

// NewProtocolProvider creates a new protocol provider
func NewProtocolProvider() *ProtocolProvider {
	return &ProtocolProvider{}
}

// RegisterServices registers protocol services
func (pp *ProtocolProvider) RegisterServices(container *Container) error {
	// Register binary codec
	container.RegisterSingleton("protocol.binary_codec", func(c *Container) (interface{}, error) {
		cfg, err := c.Resolve("config")
		if err != nil {
			return nil, err
		}

		config := cfg.(*config.Config)
		return protocol.NewBinaryCodec(&config.Protocol), nil
	}, "config")

	// Register JSON codec
	container.RegisterSingleton("protocol.json_codec", func(c *Container) (interface{}, error) {
		cfg, err := c.Resolve("config")
		if err != nil {
			return nil, err
		}

		config := cfg.(*config.Config)
		return protocol.NewJSONCodec(&config.Protocol), nil
	}, "config")

	// Register protocol manager
	container.RegisterSingleton("protocol.manager", func(c *Container) (interface{}, error) {
		binaryCodec, err := c.Resolve("protocol.binary_codec")
		if err != nil {
			return nil, err
		}

		jsonCodec, err := c.Resolve("protocol.json_codec")
		if err != nil {
			return nil, err
		}

		return protocol.NewProtocolManager(
			binaryCodec.(protocol.Codec),
			jsonCodec.(protocol.Codec),
		), nil
	}, "protocol.binary_codec", "protocol.json_codec")

	return nil
}

// WeaveProvider provides Service Weaver integration services
type WeaveProvider struct{}

// NewWeaveProvider creates a new weave provider
func NewWeaveProvider() *WeaveProvider {
	return &WeaveProvider{}
}

// RegisterServices registers weave services
func (wp *WeaveProvider) RegisterServices(container *Container) error {
	// Register weavelet manager
	container.RegisterSingleton("weave.manager", func(c *Container) (interface{}, error) {
		cfg, err := c.Resolve("config")
		if err != nil {
			return nil, err
		}

		config := cfg.(*config.Config)
		if config.Weave != nil {
			return weave.NewWeaveletManager(config.Weave), nil
		}

		// Use default config if not specified
		defaultConfig := weave.DefaultWeaveletConfig()
		return weave.NewWeaveletManager(&defaultConfig), nil
	}, "config")

	return nil
}

// AllProvidersProvider combines all service providers
type AllProvidersProvider struct {
	configPath string
}

// NewAllProvidersProvider creates a provider that registers all services
func NewAllProvidersProvider(configPath string) *AllProvidersProvider {
	return &AllProvidersProvider{
		configPath: configPath,
	}
}

// RegisterServices registers all services from all providers
func (app *AllProvidersProvider) RegisterServices(container *Container) error {
	providers := []ServiceProvider{
		NewConfigProvider(app.configPath),
		NewLoggingProvider(),
		NewMonitoringProvider(),
		NewPersistenceProvider(),
		NewProtocolProvider(),
		NewWeaveProvider(),
	}

	for _, provider := range providers {
		if err := container.RegisterProvider(provider); err != nil {
			return err
		}
	}

	return nil
}
