package activity

import (
	"greatestworks/aop/event"
	"greatestworks/business/module"
	"sync"
)

func init() {
	module.MManager.RegisterModule("", GetMe())
}

var (
	onceInit sync.Once
	Mod      *Module
)

func GetMe() *Module {
	onceInit.Do(func() {
		Mod = &Module{}
	})
	return Mod
}

type Module struct {
	*module.MetricsBase
	*module.DBActionBase
}

func (m *Module) OnEvent(c module.Character, event event.IEvent) {
	//TODO implement me
	panic("implement me")
}

func (m *Module) SetEventCategoryActive(eventCategory int) {
	//TODO implement me
	panic("implement me")
}

func (a *Module) OnStart() {
	//TODO implement me
	panic("implement me")
}

func (a *Module) AfterStart() {
	//TODO implement me
	panic("implement me")
}

func (a *Module) OnStop() {
	//TODO implement me
	panic("implement me")
}

func (a *Module) AfterStop() {
	//TODO implement me
	panic("implement me")
}
