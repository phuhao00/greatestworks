package module

import "fmt"

var (
	MManager Manager
)

type Manager struct {
	moduleName2Module map[string]IModule
}

func (m *Manager) GetModule(name string) IModule {
	return m.moduleName2Module[name]
}

func (m *Manager) RegisterModule(moduleName string, module IModule) {
	if _, exist := m.moduleName2Module[moduleName]; exist {
		panic(fmt.Sprintf("repeat register module :%v", moduleName))
	}
	m.moduleName2Module[moduleName] = module
}
