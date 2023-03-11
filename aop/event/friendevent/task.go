package friendevent

import "greatestworks/aop/event"

type AddOrDelFriendEvent struct {
	CurFriendCount int
	event.Base
}

func (e *AddOrDelFriendEvent) GetDesc() string {
	return ""
}
