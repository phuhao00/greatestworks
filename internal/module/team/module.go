package template

import (
	module2 "greatestworks/internal/module"
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
	module2.MManager.RegisterModule(ModuleName, Mod)
}

func GetMod() *Module {
	OnceInitMod.Do(func() {
		Mod = &Module{module2.NewBaseModule()}
	})
	return Mod
}

type Module struct {
	*module2.BaseModule
}

func (m *Module) SetName(name string) {
	m.BaseModule.SetName(name)
}
