package internal

import (
	"greatestworks/internal/note/event"
)

type DataAsPublisher struct {
	event.BasePublisher
}

type DataAsSubscriber struct {
	event.BaseSubscriber
}
