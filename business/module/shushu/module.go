package shushu

import (
	"greatestworks/business/module"
	"sync"
)

const (
	ModuleName = "shushu"
)

var (
	Mod         *Module
	OnceInitMod sync.Once
)

func init() {
	module.MManager.RegisterModule(ModuleName, Mod)
}

type Module struct {
	*module.BaseModule
}

func GetMod() *Module {
	OnceInitMod.Do(func() {
		Mod = &Module{
			BaseModule: module.NewBaseModule(),
		}

	})
	return Mod
}

func (m *Module) SetName(name string) {
	m.BaseModule.SetName(name)
}
