package honour

import (
	"github.com/phuhao00/greatestworks-proto/module"
	"greatestworks/aop/module_router"
	"greatestworks/internal"
	"greatestworks/internal/note/event"
)

type Module struct {
	*internal.BaseModule
}

func (m *Module) onEvent(event event.IEvent) {
	//todo
}

func (m *Module) RegisterHandler() {
	module_router.RegisterModuleMessageHandler(module.Module_Honour, 0, nil)
}

func init() {
	internal.ModuleManager.RegisterModule(module.Module_Honour.String(), GetMod())
}

func GetMod() *Module {
	return nil
}

func (m *Module) GetName() string {
	return module.Module_Honour.String()
}
