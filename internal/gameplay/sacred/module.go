package sacred

import (
	"github.com/phuhao00/greatestworks-proto/module"
	"greatestworks/aop/module_router"
	"greatestworks/internal"
	"sync"
)

var (
	Mod         *Module
	onceInitMod sync.Once
)

type Module struct {
	*internal.BaseModule
}

func init() {
	internal.ModuleManager.RegisterModule(module.Module_Sacred.String(), GetMod())
}

func GetMod() *Module {
	Mod = &Module{internal.NewBaseModule()}

	return Mod
}

func (m *Module) RegisterHandler() {
	module_router.RegisterModuleMessageHandler(module.Module_Sacred, 0, nil)
}
