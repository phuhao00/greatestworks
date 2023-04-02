package recharge

import (
	"greatestworks/internal"
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
	internal.MManager.RegisterModule(ModuleName, GetMod())
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
