package event

import "sync"

type BaseSubscriber struct {
	event2Handle sync.Map
}

func (b *BaseSubscriber) OnEvent(e IEvent) {
	//TODO implement me
	panic("implement me")
}
