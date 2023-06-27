package hangup

import (
	"github.com/phuhao00/greatestworks-proto/module"
	"greatestworks/aop/module_router"
	"greatestworks/internal"
)

type Module struct {
	*internal.BaseModule
}

func (m *Module) RegisterHandler() {
	module_router.RegisterModuleMessageHandler(module.Module_HangUp, 0, nil)
}

func init() {
	internal.ModuleManager.RegisterModule(module.Module_HangUp.String(), GetMod())
}

func GetMod() *Module {
	return nil
}

func (m *Module) GetName() string {
	return module.Module_HangUp.String()
}
