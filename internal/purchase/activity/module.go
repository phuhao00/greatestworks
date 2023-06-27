package activity

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
	internal.ModuleManager.RegisterModule(module.Module_Activity.String(), GetMod())
}

type Module struct {
	*internal.BaseModule
	*internal.MetricsBase
	*internal.DBActionBase
}

func GetMod() *Module {
	Mod = &Module{BaseModule: internal.NewBaseModule()}

	return Mod
}

func (m *Module) RegisterHandler() {
	module_router.RegisterModuleMessageHandler(module.Module_Activity, 0, nil)
}
