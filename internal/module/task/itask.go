package task

import "greatestworks/aop/event"

type ITask interface {
	SetStatus(Status)
	OnEvent(event event.IEvent)
}
