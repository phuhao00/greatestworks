package scene

import (
	"greatestworks/internal"
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
	internal.MManager.RegisterModule(ModuleName, GetMod())
}

type Module struct {
	*internal.BaseModule
}

func GetMod() *Module {
	onceInitMod.Do(func() {
		Mod = &Module{internal.NewBaseModule()}
	})

	return Mod
}
