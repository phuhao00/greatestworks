package scene

import (
	"greatestworks/business/module"
	"sync"
)

const (
	ModuleName = "scene"
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
		Mod = &Module{module.NewBaseModule()}
	})

	return Mod
}
