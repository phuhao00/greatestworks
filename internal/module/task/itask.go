package task

import (
	"greatestworks/internal/event"
)

type ITask interface {
	SetStatus(Status)
	OnEvent(event event.IEvent)
}
