package shop

import (
	"greatestworks/internal/module"
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
	module.MManager.RegisterModule(ModuleName, Mod)
}

type Module struct {
	*module.BaseModule
}

func GetMod() *Module {
	Mod = &Module{module.NewBaseModule()}

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
