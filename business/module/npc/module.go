package npc

import (
	"greatestworks/business/module"
)

func init() {
	module.MManager.RegisterModule("", Mod)
}

var (
	Mod *Module
)

type Module struct {
}
