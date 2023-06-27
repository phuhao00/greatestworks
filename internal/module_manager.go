package internal

import "fmt"

var (
	MManager ManagerOfModule
)

type ManagerOfModule struct {
	moduleName2Module map[string]IModule
}

func (m *ManagerOfModule) GetModule(name string) IModule {
	return m.moduleName2Module[name]
}

func (m *ManagerOfModule) RegisterModule(moduleName string, module IModule) {
	if _, exist := m.moduleName2Module[moduleName]; exist {
		panic(fmt.Sprintf("repeat register module.proto :%v", moduleName))
	}
	m.moduleName2Module[moduleName] = module
}
