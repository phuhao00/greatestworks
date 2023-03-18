package recharge

import (
	"greatestworks/internal/module"
	"sync"
)

const (
	ModuleName = "recharge"
)

var (
	Mod         *Module
	onceInitMod sync.Once
)

func init() {
	module.MManager.RegisterModule(ModuleName, GetMod())
}

type Module struct {
	*module.BaseModule
}

func GetMod() *Module {
	Mod = &Module{module.NewBaseModule()}

	return Mod
}

// CreateOrder 创建支付订单
func (m *Module) CreateOrder() {

}

// OnSdkOrderRsp sdk 订单返回逻辑
func (m *Module) OnSdkOrderRsp() {

}
