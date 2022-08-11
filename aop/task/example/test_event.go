package example

import "greatestworks/aop/task"

type TEvent struct {
	Data        int
	Subscribers []task.Target
}

func (e *TEvent) Notify() {
	for _, subscriber := range e.Subscribers {
		subscriber.OnNotify(e)
	}
}

func (e *TEvent) Attach(target task.Target) {
	e.Subscribers = append(e.Subscribers, target)
}

func (e *TEvent) Detach(id uint32) {
	for i, subscriber := range e.Subscribers {
		if subscriber.GetTargetId() == id {
			e.Subscribers = append(e.Subscribers[:i], e.Subscribers[i+1:]...)
		}
	}
}
