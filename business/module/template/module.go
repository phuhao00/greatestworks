package template

import (
	"greatestworks/business/module"
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
	module.MManager.RegisterModule(ModuleName, Mod)
}

type Module struct {
	*module.BaseModule
}

func GetMod() *Module {
	onceInit.Do(func() {
		Mod = &Module{module.NewBaseModule()}
	})
	return Mod
}

func (m *Module) SetName(name string) {
	m.BaseModule.SetName(name)
}
