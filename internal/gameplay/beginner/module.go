package beginner

import (
	"github.com/phuhao00/greatestworks-proto/module"
	"greatestworks/aop/module_router"
	"greatestworks/internal"
)

type Module struct {
	internal.BaseModule
}

func init() {
	internal.ModuleManager.RegisterModule(module.Module_Beginner.String(), GetMod())
}

func GetMod() *Module {
	return nil
}

func (m *Module) RegisterHandler() {
	module_router.RegisterModuleMessageHandler(module.Module_Beginner, 0, nil)
}

func (m *Module) GetName() string {
	return module.Module_Beginner.String()
}
