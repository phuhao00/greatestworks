package player

import (
	"greatestworks/aop/event"
	"greatestworks/business/module"
)

func (p *Player) OnEvent(e event.IEvent) {
	module.MManager.GetModule(e.GetToModuleName()).OnEvent(p, e)
}

func (m *Module) OnEvent(c module.Character, event event.IEvent) {
	//TODO implement me
	panic("implement me")
}

func (m *Module) SetEventCategoryActive(eventCategory int) {
	//TODO implement me
	panic("implement me")
}
