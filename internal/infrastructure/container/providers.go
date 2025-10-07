// Package container provides service providers for dependency injection
package container

import (
	"fmt"
	"greatestworks/internal/infrastructure/config"
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
		// 使用简单的配置加载器
		loader := &config.ConfigLoader{}
		return loader, nil
	})

	// Register main configuration
	container.RegisterSingleton("config", func(c *Container) (interface{}, error) {
		loader, err := c.Resolve("config.loader")
		if err != nil {
			return nil, err
		}

		configLoader := loader.(*config.ConfigLoader)
		return configLoader.Load()
	}, "config.loader")

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
		// TODO: 实现日志配置
		return nil, fmt.Errorf("logging not implemented")
	})

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
	// TODO: 实现监控配置
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
	// TODO: 实现持久化配置
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
	// TODO: 实现协议配置
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
	// TODO: 实现Weave配置
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
