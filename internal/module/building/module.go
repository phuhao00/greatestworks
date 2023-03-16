package building

import (
	module2 "greatestworks/internal/module"
	"sync"
)

const (
	ModuleName = "building"
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
	Mod = &Module{module2.NewBaseModule()}

	return Mod
}
