package weather

import (
	"greatestworks/internal"
)

var (
	Mod *Module
)

func init() {
	internal.MManager.RegisterModule("", Mod)
}

type Module struct {
}
