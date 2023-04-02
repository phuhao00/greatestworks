package task

import (
	"greatestworks/internal/note/event"
)

type ITask interface {
	SetStatus(Status)
	OnEvent(event event.IEvent)
}
