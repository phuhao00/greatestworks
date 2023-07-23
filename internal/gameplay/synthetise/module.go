package synthetise

import (
	"github.com/phuhao00/greatestworks-proto/module"
	"greatestworks/aop/module_router"
	"greatestworks/internal"
	"sync"
)

var (
	Mod         *Module
	OnceInitMod sync.Once
)

func init() {
	internal.ModuleManager.RegisterModule(module.Module_Synthetise.String(), GetMod())
}

type Module struct {
	*internal.BaseModule
}

func GetMod() *Module {
	OnceInitMod.Do(func() {
		Mod = &Module{
			BaseModule: internal.NewBaseModule(),
		}
	})
	return Mod
}

func (m *Module) SetName(name string) {
	m.BaseModule.SetName(name)
}

func (m *Module) RegisterHandler() {
	module_router.RegisterModuleMessageHandler(module.Module_Synthetise, 0, nil)
}
