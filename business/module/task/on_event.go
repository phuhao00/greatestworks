package task

import (
	"fmt"
	"greatestworks/aop/event"
	"greatestworks/business/module/friend"
)

type EventHandle func(iEvent event.IEvent)

func (d *Data) OnEvent(event event.IEvent) {
	d.eventHandles[event](event)
}

func (d *Data) HandleAddOrDelFriendEvent(iEvent event.IEvent) {
	e := iEvent.(*friend.AddOrDelFriendEvent)
	fmt.Println(e)
}
