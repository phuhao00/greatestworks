package bag

import (
	"greatestworks/business/module"
	"sync"
)

func init() {
	module.MManager.RegisterModule("", GetMe())
}

type Module struct {
}

var (
	onceInitMod sync.Once
	Mod         *Module
)

func GetMe() *Module {
	onceInitMod.Do(func() {
		Mod = &Module{}
	})
	return Mod
}
