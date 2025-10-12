// Package container provides dependency injection container functionality
package container

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

// Container represents a dependency injection container
type Container struct {
	mu        sync.RWMutex
	services  map[string]*ServiceDescriptor
	instances map[string]interface{}
	scopes    map[string]*Scope
	parent    *Container
}

// ServiceDescriptor describes how to create a service
type ServiceDescriptor struct {
	Name         string
	ServiceType  reflect.Type
	Lifetime     Lifetime
	Factory      FactoryFunc
	Instance     interface{}
	Dependencies []string
}

// FactoryFunc is a function that creates a service instance
type FactoryFunc func(container *Container) (interface{}, error)

// Lifetime defines the lifetime of a service
type Lifetime int

const (
	// Transient creates a new instance every time
	Transient Lifetime = iota
	// Singleton creates a single instance for the container lifetime
	Singleton
	// Scoped creates a single instance per scope
	Scoped
)

// Scope represents a service scope
type Scope struct {
	mu        sync.RWMutex
	instances map[string]interface{}
	parent    *Container
	closed    bool
}

// ServiceProvider defines the interface for service providers
type ServiceProvider interface {
	RegisterServices(container *Container) error
}

// Lifecycle defines the interface for services with lifecycle management
type Lifecycle interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// NewContainer creates a new dependency injection container
func NewContainer() *Container {
	return &Container{
		services:  make(map[string]*ServiceDescriptor),
		instances: make(map[string]interface{}),
		scopes:    make(map[string]*Scope),
	}
}

// NewChildContainer creates a child container
func (c *Container) NewChildContainer() *Container {
	return &Container{
		services:  make(map[string]*ServiceDescriptor),
		instances: make(map[string]interface{}),
		scopes:    make(map[string]*Scope),
		parent:    c,
	}
}

// RegisterTransient registers a transient service
func (c *Container) RegisterTransient(name string, factory FactoryFunc, dependencies ...string) {
	c.register(name, &ServiceDescriptor{
		Name:         name,
		Lifetime:     Transient,
		Factory:      factory,
		Dependencies: dependencies,
	})
}

// RegisterSingleton registers a singleton service
func (c *Container) RegisterSingleton(name string, factory FactoryFunc, dependencies ...string) {
	c.register(name, &ServiceDescriptor{
		Name:         name,
		Lifetime:     Singleton,
		Factory:      factory,
		Dependencies: dependencies,
	})
}

// RegisterScoped registers a scoped service
func (c *Container) RegisterScoped(name string, factory FactoryFunc, dependencies ...string) {
	c.register(name, &ServiceDescriptor{
		Name:         name,
		Lifetime:     Scoped,
		Factory:      factory,
		Dependencies: dependencies,
	})
}

// RegisterInstance registers a service instance
func (c *Container) RegisterInstance(name string, instance interface{}) {
	c.register(name, &ServiceDescriptor{
		Name:        name,
		ServiceType: reflect.TypeOf(instance),
		Lifetime:    Singleton,
		Instance:    instance,
	})
}

// RegisterProvider registers a service provider
func (c *Container) RegisterProvider(provider ServiceProvider) error {
	return provider.RegisterServices(c)
}

// register registers a service descriptor
func (c *Container) register(name string, descriptor *ServiceDescriptor) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = descriptor
}

// Resolve resolves a service by name
func (c *Container) Resolve(name string) (interface{}, error) {
	return c.ResolveWithScope(name, nil)
}

// ResolveWithScope resolves a service with a specific scope
func (c *Container) ResolveWithScope(name string, scope *Scope) (interface{}, error) {
	c.mu.RLock()
	descriptor, exists := c.services[name]
	c.mu.RUnlock()

	if !exists {
		// Try parent container
		if c.parent != nil {
			return c.parent.ResolveWithScope(name, scope)
		}
		return nil, fmt.Errorf("service '%s' not registered", name)
	}

	return c.createInstance(descriptor, scope)
}

// createInstance creates a service instance based on its descriptor
func (c *Container) createInstance(descriptor *ServiceDescriptor, scope *Scope) (interface{}, error) {
	switch descriptor.Lifetime {
	case Singleton:
		return c.getSingletonInstance(descriptor)
	case Scoped:
		if scope == nil {
			return nil, fmt.Errorf("scoped service '%s' requires a scope", descriptor.Name)
		}
		return c.getScopedInstance(descriptor, scope)
	case Transient:
		return c.createTransientInstance(descriptor, scope)
	default:
		return nil, fmt.Errorf("unknown lifetime for service '%s'", descriptor.Name)
	}
}

// getSingletonInstance gets or creates a singleton instance
func (c *Container) getSingletonInstance(descriptor *ServiceDescriptor) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if descriptor.Instance != nil {
		return descriptor.Instance, nil
	}

	if instance, exists := c.instances[descriptor.Name]; exists {
		return instance, nil
	}

	instance, err := c.createInstanceInternal(descriptor, nil)
	if err != nil {
		return nil, err
	}

	c.instances[descriptor.Name] = instance
	return instance, nil
}

// getScopedInstance gets or creates a scoped instance
func (c *Container) getScopedInstance(descriptor *ServiceDescriptor, scope *Scope) (interface{}, error) {
	scope.mu.Lock()
	defer scope.mu.Unlock()

	if scope.closed {
		return nil, fmt.Errorf("scope is closed")
	}

	if instance, exists := scope.instances[descriptor.Name]; exists {
		return instance, nil
	}

	instance, err := c.createInstanceInternal(descriptor, scope)
	if err != nil {
		return nil, err
	}

	scope.instances[descriptor.Name] = instance
	return instance, nil
}

// createTransientInstance creates a new transient instance
func (c *Container) createTransientInstance(descriptor *ServiceDescriptor, scope *Scope) (interface{}, error) {
	return c.createInstanceInternal(descriptor, scope)
}

// createInstanceInternal creates an instance using the factory function
func (c *Container) createInstanceInternal(descriptor *ServiceDescriptor, scope *Scope) (interface{}, error) {
	if descriptor.Factory == nil {
		return nil, fmt.Errorf("no factory function for service '%s'", descriptor.Name)
	}

	// Resolve dependencies first
	for _, dep := range descriptor.Dependencies {
		_, err := c.ResolveWithScope(dep, scope)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve dependency '%s' for service '%s': %w", dep, descriptor.Name, err)
		}
	}

	return descriptor.Factory(c)
}

// CreateScope creates a new service scope
func (c *Container) CreateScope() *Scope {
	return &Scope{
		instances: make(map[string]interface{}),
		parent:    c,
		closed:    false,
	}
}

// Close closes the scope and disposes scoped services
func (s *Scope) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	// Dispose services that implement Lifecycle
	for _, instance := range s.instances {
		if lifecycle, ok := instance.(Lifecycle); ok {
			if err := lifecycle.Stop(context.Background()); err != nil {
				// Log error but continue disposing other services
				continue
			}
		}
	}

	s.instances = nil
	s.closed = true
	return nil
}

// IsRegistered checks if a service is registered
func (c *Container) IsRegistered(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.services[name]
	if !exists && c.parent != nil {
		return c.parent.IsRegistered(name)
	}
	return exists
}

// GetServiceNames returns all registered service names
func (c *Container) GetServiceNames() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	names := make([]string, 0, len(c.services))
	for name := range c.services {
		names = append(names, name)
	}
	return names
}

// StartServices starts all services that implement Lifecycle
func (c *Container) StartServices(ctx context.Context) error {
	c.mu.RLock()
	services := make([]*ServiceDescriptor, 0, len(c.services))
	for _, service := range c.services {
		services = append(services, service)
	}
	c.mu.RUnlock()

	for _, service := range services {
		instance, err := c.Resolve(service.Name)
		if err != nil {
			return fmt.Errorf("failed to resolve service '%s': %w", service.Name, err)
		}

		if lifecycle, ok := instance.(Lifecycle); ok {
			if err := lifecycle.Start(ctx); err != nil {
				return fmt.Errorf("failed to start service '%s': %w", service.Name, err)
			}
		}
	}

	return nil
}

// StopServices stops all services that implement Lifecycle
func (c *Container) StopServices(ctx context.Context) error {
	c.mu.RLock()
	instances := make([]interface{}, 0, len(c.instances))
	for _, instance := range c.instances {
		instances = append(instances, instance)
	}
	c.mu.RUnlock()

	// Stop services in reverse order
	for i := len(instances) - 1; i >= 0; i-- {
		if lifecycle, ok := instances[i].(Lifecycle); ok {
			if err := lifecycle.Stop(ctx); err != nil {
				// Log error but continue stopping other services
				continue
			}
		}
	}

	return nil
}

// Dispose disposes the container and all its services
func (c *Container) Dispose() error {
	ctx := context.Background()
	return c.StopServices(ctx)
}
