package synthetise

import (
	module2 "greatestworks/internal/module"
	"sync"
)

const (
	ModuleName = "synthetise"
)

var (
	Mod         *Module
	OnceInitMod sync.Once
)

func init() {
	module2.MManager.RegisterModule(ModuleName, GetMod())
}

type Module struct {
	*module2.BaseModule
}

func GetMod() *Module {
	OnceInitMod.Do(func() {
		Mod = &Module{
			BaseModule: module2.NewBaseModule(),
		}
	})
	return Mod
}

func (m *Module) SetName(name string) {
	m.BaseModule.SetName(name)
}
