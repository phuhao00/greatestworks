package vip

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
	internal.ModuleManager.RegisterModule(module.Module_Vip.String(), GetMod())
}

type Module struct {
	*internal.BaseModule
}

func GetMod() *Module {
	onceInitMod.Do(func() {
		Mod = &Module{
			BaseModule: internal.NewBaseModule(),
		}
	})
	return Mod
}

func (m *Module) SetName(name string) {
	m.BaseModule.SetName(name)
}

func (m *Module) OnDailyReset() {
	//
}

func (m *Module) OnRecharge() {
	//add exp
}

func (m *Module) GetName() string {
	return module.Module_Vip.String()
}

func (m *Module) RegisterHandler() {
	module_router.RegisterModuleMessageHandler(module.Module_Vip, 0, nil)
}
