package event

type Publisher interface {
	AddSubscriber(e IEvent, subscriber Subscriber)
	Publish(e IEvent)
}

type Subscriber interface {
	OnEvent(e IEvent)
}
