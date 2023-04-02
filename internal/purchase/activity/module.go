package activity

import (
	"greatestworks/internal"
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
	internal.MManager.RegisterModule(ModuleName, GetMod())
}

type Module struct {
	*internal.BaseModule
	*internal.MetricsBase
	*internal.DBActionBase
}

func GetMod() *Module {
	Mod = &Module{BaseModule: internal.NewBaseModule()}

	return Mod
}
