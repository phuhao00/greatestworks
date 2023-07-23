package pet

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

func init() {
	internal.ModuleManager.RegisterModule(module.Module_Pet.String(), GetMod())
}

type Module struct {
	*internal.BaseModule
}

func GetMod() *Module {
	Mod = &Module{internal.NewBaseModule()}

	return Mod
}

func (m *Module) RegisterHandler() {
	module_router.RegisterModuleMessageHandler(module.Module_Pet, 0, nil)
}

func (m *Module) GetName() string {
	return module.Module_Pet.String()
}
