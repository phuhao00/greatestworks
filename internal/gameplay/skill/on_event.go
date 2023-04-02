package skill

import (
	"greatestworks/internal"
	"greatestworks/internal/note/event"
)

type EventHandle func(iEvent event.IEvent)

type EventWrap struct {
	IPlayer
	event.IEvent
}

func (m *Module) OnEvent(c internal.Character, event event.IEvent) {
	//TODO implement me
	panic("implement me")
}

func (m *Module) SetEventCategoryActive(eventCategory int) {
	//TODO implement me
	panic("implement me")
}
