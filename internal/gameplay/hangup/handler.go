package hangup

import "greatestworks/aop/module_router"

func RegisterHandler() {
	module_router.RegisterModuleMessageHandler(0, 0, nil)
}
