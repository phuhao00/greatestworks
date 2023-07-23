package weather

import (
	"greatestworks/aop/module_router"
	"greatestworks/internal"
)

var (
	Mod *Module
)

func init() {
	internal.ModuleManager.RegisterModule("", Mod)
}

type Module struct {
	*internal.BaseModule
}

func (m *Module) RegisterHandler() {
	module_router.RegisterModuleMessageHandler(0, 0, nil)
}
