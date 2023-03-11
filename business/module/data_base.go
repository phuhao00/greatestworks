package module

import "greatestworks/aop/event"

type DataAsPublisher struct {
	event.BasePublisher
}

type DataAsSubscriber struct {
	event.BaseSubscriber
}
