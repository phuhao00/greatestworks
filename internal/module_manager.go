package internal

import (
	"context"
	"fmt"
	"sync"
)

var (
	ModuleManager *ModuleManagerImpl
)

// ModuleManagerImpl 模块管理器实现
type ModuleManagerImpl struct {
	mu                sync.RWMutex
	moduleName2Module map[string]Module
	started           bool
}

// NewModuleManager 创建新的模块管理器
func NewModuleManager() *ModuleManagerImpl {
	return &ModuleManagerImpl{
		moduleName2Module: make(map[string]Module),
	}
}

// GetModule 获取模块
func (m *ModuleManagerImpl) GetModule(name string) Module {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.moduleName2Module[name]
}

// RegisterModule 注册模块
func (m *ModuleManagerImpl) RegisterModule(moduleName string, module Module) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exist := m.moduleName2Module[moduleName]; exist {
		return fmt.Errorf("重复注册模块: %v", moduleName)
	}

	m.moduleName2Module[moduleName] = module
	return nil
}

// UnregisterModule 注销模块
func (m *ModuleManagerImpl) UnregisterModule(moduleName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exist := m.moduleName2Module[moduleName]; !exist {
		return fmt.Errorf("模块不存在: %v", moduleName)
	}

	delete(m.moduleName2Module, moduleName)
	return nil
}

// GetModuleNames 获取所有模块名称
func (m *ModuleManagerImpl) GetModuleNames() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.moduleName2Module))
	for name := range m.moduleName2Module {
		names = append(names, name)
	}
	return names
}

// Start 启动所有模块
func (m *ModuleManagerImpl) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.started {
		return fmt.Errorf("模块管理器已经启动")
	}

	// 注册处理器
	for _, module := range m.moduleName2Module {
		module.RegisterHandler()
	}

	// 启动模块
	for _, module := range m.moduleName2Module {
		module.OnStart()
		// 如果模块实现了Manager接口，调用AfterStart
		if manager, ok := module.(Manager); ok {
			manager.AfterStart()
		}
	}

	m.started = true
	return nil
}

// Stop 停止所有模块
func (m *ModuleManagerImpl) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.started {
		return nil
	}

	// 停止模块（逆序）
	modules := make([]Module, 0, len(m.moduleName2Module))
	for _, module := range m.moduleName2Module {
		modules = append(modules, module)
	}

	for i := len(modules) - 1; i >= 0; i-- {
		module := modules[i]
		module.OnStop()
		// 如果模块实现了Manager接口，调用AfterStop
		if manager, ok := module.(Manager); ok {
			manager.AfterStop()
		}
	}

	m.started = false
	return nil
}

// IsStarted 检查是否已启动
func (m *ModuleManagerImpl) IsStarted() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.started
}

// 初始化全局模块管理器
func init() {
	ModuleManager = NewModuleManager()
}
