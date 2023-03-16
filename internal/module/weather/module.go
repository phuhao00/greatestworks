package weather

import (
	"greatestworks/internal/module"
)

var (
	Mod *Module
)

func init() {
	module.MManager.RegisterModule("", Mod)
}

type Module struct {
}

