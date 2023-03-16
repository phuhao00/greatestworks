package activity

import (
	module2 "greatestworks/internal/module"
	"sync"
)

const (
	ModuleName = "activity"
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
	*module2.MetricsBase
	*module2.DBActionBase
}

func GetMod() *Module {
	Mod = &Module{BaseModule: module2.NewBaseModule()}

	return Mod
}
