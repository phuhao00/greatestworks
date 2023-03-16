package module

import (
	"greatestworks/internal/event"
)

type DataAsPublisher struct {
	event.BasePublisher
}

type DataAsSubscriber struct {
	event.BaseSubscriber
}
