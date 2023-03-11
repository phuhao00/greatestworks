package event

import (
	"greatestworks/aop/condition"
)

type Base struct {
	Subscribers []condition.Condition
}

func (b *Base) Notify() {
	for _, subscriber := range b.Subscribers {
		subscriber.OnNotify(b)
	}
}

func (b *Base) Attach(c condition.Condition) {
	b.Subscribers = append(b.Subscribers, c)
}

func (b *Base) Detach(id uint32) {
	for i, subscriber := range b.Subscribers {
		if subscriber.GetId() == id {
			b.Subscribers = append(b.Subscribers[:i], b.Subscribers[i+1:]...)
		}
	}
}
