package template

import (
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
	internal.MManager.RegisterModule(ModuleName, Mod)
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
