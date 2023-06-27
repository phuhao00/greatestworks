package template

import (
	"greatestworks/internal"
	"sync"
)

const (
	ModuleName = "template"
)

var (
	Mod         *Module
	OnceInitMod sync.Once
)

func init() {
	internal.ModuleManager.RegisterModule(ModuleName, Mod)
}

type Module struct {
	*internal.BaseModule
}

func GetMod() *Module {
	onceInit.Do(func() {
		Mod = &Module{internal.NewBaseModule()}
	})
	return Mod
}

func (m *Module) SetName(name string) {
	m.BaseModule.SetName(name)
}
