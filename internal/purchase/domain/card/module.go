package card

import (
	"greatestworks/aop/module_router"
	"greatestworks/internal"
)

type Module struct {
	*internal.BaseModule
}

func init() {
	internal.ModuleManager.RegisterModule("", GetMod())
}

func GetMod() *Module {
	return nil
}

func (m *Module) RegisterHandler() {
	module_router.RegisterModuleMessageHandler(0, 0, nil)
}
