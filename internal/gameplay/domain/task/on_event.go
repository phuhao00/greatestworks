package task

import (
	"greatestworks/internal/note/event"
)

type EventHandle func(iEvent event.IEvent)

type EventWrap struct {
	Player
	event.IEvent
}
