package name

import (
	"github.com/phuhao00/greatestworks-proto/module"
	"greatestworks/aop/module_router"
	"greatestworks/internal"
)

type Name struct {
	*internal.BaseModule
}

func init() {
	internal.ModuleManager.RegisterModule(module.Module_Name.String(), GetMod())
}

func GetMod() internal.IModule {
	return nil
}

func (n *Name) RandomName() {

}

func (m *Name) RegisterHandler() {
	module_router.RegisterModuleMessageHandler(module.Module_Name, 0, nil)
}

func (pm *Name) GetName() string {
	return module.Module_Name.String()
}
