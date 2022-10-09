package example

import (
	"greatestworks/business/module/condition"
)

type TEvent struct {
	Data        int
	Subscribers []condition.Condition
}

func (e *TEvent) Notify() {
	for _, subscriber := range e.Subscribers {
		subscriber.OnNotify(e)
	}
}

func (e *TEvent) Attach(target condition.Condition) {
	e.Subscribers = append(e.Subscribers, target)
}

func (e *TEvent) Detach(id uint32) {
	for i, subscriber := range e.Subscribers {
		if subscriber.GetId() == id {
			e.Subscribers = append(e.Subscribers[:i], e.Subscribers[i+1:]...)
		}
	}
}
