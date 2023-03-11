package event

import "sync"

type Normal struct {
	subscribers sync.Map
}

func (n *Normal) AddSubscriber(e IEvent, subscriber Subscriber) {
	n.subscribers.LoadOrStore(e, subscriber)
}

func (n *Normal) Publish(e IEvent) {
	n.subscribers.Range(func(key, value any) bool {
		s := value.(Subscriber)
		s.OnEvent(nil)
		return true
	})
}
