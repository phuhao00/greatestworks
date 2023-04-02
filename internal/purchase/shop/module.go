package shop

import (
	"greatestworks/internal"
	"sync"
)

const (
	ModuleName = "shop"
)

var (
	Mod         *Module
	onceInitMod sync.Once
)

func init() {
	internal.MManager.RegisterModule(ModuleName, Mod)
}

type Module struct {
	*internal.BaseModule
}

func GetMod() *Module {
	Mod = &Module{internal.NewBaseModule()}

	return Mod
}

// ResourceBuy eg:金币，钻石
func (m *Module) ResourceBuy() {

}

// MoneyBuy 直购
func (m *Module) MoneyBuy() {

}

// TokenMoneyBuy 代金券
func (m *Module) TokenMoneyBuy() {

}
