package template

import (
	module2 "greatestworks/internal/module"
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
	module2.MManager.RegisterModule(ModuleName, Mod)
}

type Module struct {
	*module2.BaseModule
}

func GetMod() *Module {
	onceInit.Do(func() {
		Mod = &Module{module2.NewBaseModule()}
	})
	return Mod
}

func (m *Module) SetName(name string) {
	m.BaseModule.SetName(name)
}
