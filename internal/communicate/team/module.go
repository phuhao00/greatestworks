package template

import (
	"greatestworks/aop/module_router"
	"greatestworks/internal"
	"sync"
)

const (
	ModuleName = "team"
)

var (
	Mod         *Module
	OnceInitMod sync.Once
)

func init() {
	internal.ModuleManager.RegisterModule(ModuleName, Mod)
}

func GetMod() *Module {
	OnceInitMod.Do(func() {
		Mod = &Module{internal.NewBaseModule()}
	})
	return Mod
}

type Module struct {
	*internal.BaseModule
}

func (m *Module) SetName(name string) {
	m.BaseModule.SetName(name)
}

func (m *Module) RegisterHandler() {
	module_router.RegisterModuleMessageHandler(0, 0, nil)
}
