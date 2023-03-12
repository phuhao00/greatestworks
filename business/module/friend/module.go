package friend

import (
	"greatestworks/business/module"
)

var (
	Mod *Module
)

func init() {
	module.MManager.RegisterModule("", Mod)
}

func GetName() string {
	return ""
}

type Module struct {
}
