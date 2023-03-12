package activity

import (
	"greatestworks/business/module"
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
	module.MManager.RegisterModule(ModuleName, GetMod())
}

type Module struct {
	*module.BaseModule
	*module.MetricsBase
	*module.DBActionBase
}

func GetMod() *Module {
	Mod = &Module{BaseModule: module.NewBaseModule()}

	return Mod
}
