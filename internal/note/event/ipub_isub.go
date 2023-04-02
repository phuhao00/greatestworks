package event

import "sync"

type Publisher interface {
	AddSubscriber(e IEvent, subscriber Subscriber)
	Publish(e IEvent)
}

type Subscriber interface {
	OnEvent(e IEvent)
}

type BasePublisher struct {
	event2Subscribers map[IEvent][]Subscriber
}

func (b *BasePublisher) AddSubscriber(e IEvent, subscriber Subscriber) {
	if e != nil {
		b.event2Subscribers[e] = append(b.event2Subscribers[e], subscriber)
	}
}

func (b *BasePublisher) Publish(e IEvent) {
	for event, subscribers := range b.event2Subscribers {
		for _, subscriber := range subscribers {
			subscriber.OnEvent(event)
		}
	}
}

type BaseSubscriber struct {
	event2Handle sync.Map
}

func (b *BaseSubscriber) OnEvent(e IEvent) {
	//TODO implement me
	panic("implement me")
}
