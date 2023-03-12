package npc

import (
	"greatestworks/aop/event"
	"greatestworks/business/module"
)

func init() {
	module.MManager.RegisterModule("", Mod)
}

var (
	Mod *Module
)

type Module struct {
}

func (m Module) OnEvent(c module.Character, event event.IEvent) {
	//TODO implement me
	panic("implement me")
}

func (m Module) SetEventCategoryActive(eventCategory int) {
	//TODO implement me
	panic("implement me")
}
