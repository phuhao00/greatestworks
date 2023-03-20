package vip

import (
	"greatestworks/internal/module"
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
	module.MManager.RegisterModule(ModuleName, GetMod())
}

type Module struct {
	*module.BaseModule
}

func GetMod() *Module {
	onceInitMod.Do(func() {
		Mod = &Module{
			BaseModule: module.NewBaseModule(),
		}
	})
	return Mod
}

func (m *Module) SetName(name string) {
	m.BaseModule.SetName(name)
}

func (m *Module) OnDailyReset() {
	//
}

func (m *Module) OnRecharge() {
	//add exp
}
