package player

import (
	"greatestworks/internal/event"
	module2 "greatestworks/internal/module"
)

func (p *Player) OnEvent(e event.IEvent) {
	module2.MManager.GetModule(e.GetToModuleName()).OnEvent(p, e)
}

func (m *Module) OnEvent(c module2.Character, event event.IEvent) {
	//TODO implement me
	panic("implement me")
}

func (m *Module) SetEventCategoryActive(eventCategory int) {
	//TODO implement me
	panic("implement me")
}
