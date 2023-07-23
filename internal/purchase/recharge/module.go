package recharge

import (
	"github.com/phuhao00/greatestworks-proto/module"
	"greatestworks/aop/module_router"
	"greatestworks/internal"
	"sync"
)

var (
	Mod         *Module
	onceInitMod sync.Once
)

func init() {
	internal.ModuleManager.RegisterModule(module.Module_Recharge.String(), GetMod())
}

type Module struct {
	*internal.BaseModule
}

func GetMod() *Module {
	Mod = &Module{internal.NewBaseModule()}

	return Mod
}

// CreateOrder 创建支付订单
func (m *Module) CreateOrder() {

}

// OnSdkOrderRsp sdk 订单返回逻辑
func (m *Module) OnSdkOrderRsp() {

}

func (m *Module) GetName() string {
	return ""
}

func (m *Module) RegisterHandler() {
	module_router.RegisterModuleMessageHandler(module.Module_Recharge, 0, nil)
}
