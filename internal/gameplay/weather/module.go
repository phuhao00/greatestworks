package weather

import (
	"greatestworks/internal"
)

var (
	Mod *Module
)

func init() {
	internal.ModuleManager.RegisterModule("", Mod)
}

type Module struct {
}
