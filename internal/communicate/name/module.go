package name

import "greatestworks/aop/module_router"

type Name struct {
}

func (n *Name) RandomName() {

}

func (m *Name) RegisterHandler() {
	module_router.RegisterModuleMessageHandler(0, 0, nil)
}
