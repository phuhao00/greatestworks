package player

import (
	"greatestworks/internal"
	"greatestworks/internal/note/event"
)

func (p *Player) OnEvent(e event.IEvent) {
	internal.ModuleManager.GetModule(e.GetToModuleName()).OnEvent(p, e)
}

func (m *Module) OnEvent(c internal.Character, event event.IEvent) {
	//TODO implement me
	panic("implement me")
}

func (m *Module) SetEventCategoryActive(eventCategory int) {
	//TODO implement me
	panic("implement me")
}
