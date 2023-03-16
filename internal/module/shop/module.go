package shop

import (
	module2 "greatestworks/internal/module"
	"sync"
)

const (
	ModuleName = "shop"
)

var (
	Mod         *Module
	onceInitMod sync.Once
)

func init() {
	module2.MManager.RegisterModule(ModuleName, Mod)
}

type Module struct {
	*module2.BaseModule
}

func GetMod() *Module {
	Mod = &Module{module2.NewBaseModule()}

	return Mod
}
