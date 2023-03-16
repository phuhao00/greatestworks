package vip

import (
	module2 "greatestworks/internal/module"
	"sync"
)

const (
	ModuleName = "vip"
)

var (
	Mod         *Module
	onceInitMod sync.Once
)

func init() {
	module2.MManager.RegisterModule(ModuleName, GetMod())
}

type Module struct {
	*module2.BaseModule
}

func GetMod() *Module {
	onceInitMod.Do(func() {
		Mod = &Module{
			BaseModule: module2.NewBaseModule(),
		}
	})
	return Mod
}

func (m *Module) SetName(name string) {
	m.BaseModule.SetName(name)
}
