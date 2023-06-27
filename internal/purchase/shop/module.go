package shop

import (
	"github.com/phuhao00/greatestworks-proto/module"
	"greatestworks/internal"
	"sync"
)

var (
	Mod         *Module
	onceInitMod sync.Once
)

func init() {
	internal.ModuleManager.RegisterModule(module.Module_Shop.String(), Mod)
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

func (m *Module) GetName() string {
	return module.Module_Shop.String()
}
