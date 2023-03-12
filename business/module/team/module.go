package template

import (
	"greatestworks/business/module"
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
	module.MManager.RegisterModule(ModuleName, Mod)
}

func GetMod() *Module {
	OnceInitMod.Do(func() {
		Mod = &Module{module.NewBaseModule()}
	})
	return Mod
}

type Module struct {
	*module.BaseModule
}

func (m *Module) SetName(name string) {
	m.BaseModule.SetName(name)
}
