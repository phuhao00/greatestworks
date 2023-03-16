package vip

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

}

func (m *Module) SetEventCategoryActive(eventCategory int) {

}

func (m *Module) AddSubscriber(e event.IEvent, subscriber event.Subscriber) {

}

func (m *Module) Publish(iEvent event.IEvent) {

}
