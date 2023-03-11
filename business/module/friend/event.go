package friend

import (
	"greatestworks/aop/event"
	manager2 "greatestworks/business/module/hub"
)

type EventCategory int

const (
	CountOfFriend EventCategory = iota + 1
)

var (
	category2Event map[EventCategory]event.IEvent
)

func init() {
	category2Event[CountOfFriend] = &AddOrDelFriendEvent{}

	//
	manager2.MManager.AddModuleName2ModuleGetEventFunc(GetModule(), GetEvent)
}

func GetEvent(category int) event.IEvent {
	return category2Event[EventCategory(category)]
}

type AddOrDelFriendEvent struct {
	CurFriendCount int
	event.Base
}

func (e *AddOrDelFriendEvent) GetDesc() string {
	return ""
}

func (s *System) PublishAddOrDelFriend() {
	e := &AddOrDelFriendEvent{CurFriendCount: len(s.friends)}
	s.DataAsPublisher.Publish(e)
}
