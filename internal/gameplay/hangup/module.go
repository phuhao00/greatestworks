package hangup

import "greatestworks/aop/module_router"

type Module struct {
}

func (m *Module) RegisterHandler() {
	module_router.RegisterModuleMessageHandler(0, 0, nil)
}
