package friendevent

import (
	"greatestworks/internal/note/event"
)

type AddOrDelFriendEvent struct {
	CurFriendCount int
	event.Base
}

func (e *AddOrDelFriendEvent) GetDesc() string {
	return ""
}
