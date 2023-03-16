package vip

import (
	event2 "greatestworks/internal/event"
	"greatestworks/internal/module"
)

type EventHandle func(iEvent event2.IEvent)

type EventWrap struct {
	Player
	event2.IEvent
}

func (m *Module) OnEvent(c module.Character, event event2.IEvent) {

}

func (m *Module) SetEventCategoryActive(eventCategory int) {

}

func (m *Module) AddSubscriber(e event2.IEvent, subscriber event2.Subscriber) {

}

func (m *Module) Publish(iEvent event2.IEvent) {

}
