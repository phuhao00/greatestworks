package template

import (
	"greatestworks/aop/event"
	"greatestworks/internal/module"
)

type EventHandle func(iEvent event.IEvent)

type EventWrap struct {
	Player
	event.IEvent
}

func (m *Module) OnEvent(c module.Character, event event.IEvent) {
	//TODO implement me
	panic("implement me")
}

func (m *Module) SetEventCategoryActive(eventCategory int) {
	//TODO implement me
	panic("implement me")
}
