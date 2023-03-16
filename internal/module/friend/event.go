package friend

import (
	"greatestworks/aop/event"
	"greatestworks/aop/event/friendevent"
)

type (
	EventCategory int
	CreateEventFn func(system *System) event.IEvent
)

var (
	category2CreateEventFn map[EventCategory]CreateEventFn
)

const (
	CountOfFriend EventCategory = iota + 1
)

func init() {
	category2CreateEventFn[CountOfFriend] = CreateAddOrDelFriendEvent
}

func CreateAddOrDelFriendEvent(system *System) event.IEvent {
	//todo 使用原子模型做
	return &friendevent.AddOrDelFriendEvent{
		CurFriendCount: len(system.FriendList),
		Base:           event.Base{},
	}
}

func GetCreateEventFn(category EventCategory) CreateEventFn {
	return category2CreateEventFn[category]
}

func (s *System) Publish(e event.IEvent) {
	s.IPlayer.OnEvent(e)
}

func (s *System) PublishAddOrDelFriend() {
	if s.activeEventCategory[int(CountOfFriend)] {
		s.Publish(GetCreateEventFn(CountOfFriend)(s))
	}
}
