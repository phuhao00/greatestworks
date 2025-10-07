// Package container 提供简化的依赖注入容器
package container

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

// SimpleContainer 简化的依赖注入容器
type SimpleContainer struct {
	mu        sync.RWMutex
	services  map[string]interface{}
	factories map[string]func() (interface{}, error)
}

// NewSimpleContainer 创建新的简化容器
func NewSimpleContainer() *SimpleContainer {
	return &SimpleContainer{
		services:  make(map[string]interface{}),
		factories: make(map[string]func() (interface{}, error)),
	}
}

// RegisterSingleton 注册单例服务
func (c *SimpleContainer) RegisterSingleton(name string, factory func() (interface{}, error)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.factories[name] = factory
}

// RegisterInstance 注册服务实例
func (c *SimpleContainer) RegisterInstance(name string, instance interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = instance
}

// RegisterTransient 注册瞬态服务
func (c *SimpleContainer) RegisterTransient(name string, factory func() (interface{}, error)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 对于瞬态服务，每次都调用工厂函数
	c.factories[name] = factory
}

// Resolve 解析服务
func (c *SimpleContainer) Resolve(name string) (interface{}, error) {
	c.mu.RLock()
	instance, exists := c.services[name]
	c.mu.RUnlock()

	if exists {
		return instance, nil
	}

	c.mu.RLock()
	factory, exists := c.factories[name]
	c.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("service '%s' not registered", name)
	}

	// 创建实例
	instance, err := factory()
	if err != nil {
		return nil, fmt.Errorf("failed to create service '%s': %w", name, err)
	}

	// 检查是否为单例
	c.mu.Lock()
	if _, isSingleton := c.services[name]; !isSingleton {
		// 如果是单例，缓存实例
		c.services[name] = instance
	}
	c.mu.Unlock()

	return instance, nil
}

// ResolveTyped 解析类型化服务
func ResolveTyped[T any](c *SimpleContainer, name string) (T, error) {
	instance, err := c.Resolve(name)
	if err != nil {
		var zero T
		return zero, err
	}

	typed, ok := instance.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("service '%s' is not of type %T", name, zero)
	}

	return typed, nil
}

// AutoRegister 自动注册服务（通过反射）
func (c *SimpleContainer) AutoRegister(service interface{}) error {
	serviceType := reflect.TypeOf(service)
	if serviceType.Kind() == reflect.Ptr {
		serviceType = serviceType.Elem()
	}

	name := serviceType.Name()
	factory := func() (interface{}, error) {
		return service, nil
	}

	c.RegisterSingleton(name, factory)
	return nil
}

// RegisterWithDependencies 注册带依赖的服务
func (c *SimpleContainer) RegisterWithDependencies(name string, factory func(container *SimpleContainer) (interface{}, error)) {
	c.RegisterSingleton(name, func() (interface{}, error) {
		return factory(c)
	})
}

// IsRegistered 检查服务是否已注册
func (c *SimpleContainer) IsRegistered(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.services[name]
	if !exists {
		_, exists = c.factories[name]
	}
	return exists
}

// GetServiceNames 获取所有已注册的服务名称
func (c *SimpleContainer) GetServiceNames() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	names := make([]string, 0, len(c.services)+len(c.factories))
	for name := range c.services {
		names = append(names, name)
	}
	for name := range c.factories {
		if _, exists := c.services[name]; !exists {
			names = append(names, name)
		}
	}
	return names
}

// Clear 清空容器
func (c *SimpleContainer) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services = make(map[string]interface{})
	c.factories = make(map[string]func() (interface{}, error))
}

// StartServices 启动所有实现Lifecycle接口的服务
func (c *SimpleContainer) StartServices(ctx context.Context) error {
	c.mu.RLock()
	services := make([]interface{}, 0, len(c.services))
	for _, service := range c.services {
		services = append(services, service)
	}
	c.mu.RUnlock()

	for _, service := range services {
		if lifecycle, ok := service.(Lifecycle); ok {
			if err := lifecycle.Start(ctx); err != nil {
				return fmt.Errorf("failed to start service: %w", err)
			}
		}
	}

	return nil
}

// StopServices 停止所有实现Lifecycle接口的服务
func (c *SimpleContainer) StopServices(ctx context.Context) error {
	c.mu.RLock()
	services := make([]interface{}, 0, len(c.services))
	for _, service := range c.services {
		services = append(services, service)
	}
	c.mu.RUnlock()

	// 逆序停止服务
	for i := len(services) - 1; i >= 0; i-- {
		if lifecycle, ok := services[i].(Lifecycle); ok {
			if err := lifecycle.Stop(ctx); err != nil {
				// 记录错误但继续停止其他服务
				continue
			}
		}
	}

	return nil
}
