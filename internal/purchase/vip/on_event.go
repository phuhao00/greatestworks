package vip

import (
	"greatestworks/internal"
	"greatestworks/internal/note/event"
)

type EventHandle func(iEvent event.IEvent)

type EventWrap struct {
	Player
	event.IEvent
}

func (m *Module) OnEvent(c internal.Character, event event.IEvent) {

}

func (m *Module) SetEventCategoryActive(eventCategory int) {

}

func (m *Module) AddSubscriber(e event.IEvent, subscriber event.Subscriber) {

}

func (m *Module) Publish(iEvent event.IEvent) {

}
