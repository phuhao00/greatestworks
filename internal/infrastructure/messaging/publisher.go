package messaging

import (
	"context"
	"reflect"
	"time"

	"greatestworks/internal/events"
)

// EventBusPublisher publishes application events to the EventBus.
// If the event doesn't implement events.Event, it is wrapped into BaseEvent.
type EventBusPublisher struct{ bus *events.EventBus }

func NewEventBusPublisher(bus *events.EventBus) *EventBusPublisher {
	return &EventBusPublisher{bus: bus}
}

func (p *EventBusPublisher) Publish(ctx context.Context, event interface{}) error {
	if e, ok := event.(events.Event); ok {
		return p.bus.Publish(ctx, e)
	}
	wrapper := &events.BaseEvent{
		Type:      reflect.TypeOf(event).String(),
		Timestamp: time.Now(),
		Data:      event,
	}
	return p.bus.Publish(ctx, wrapper)
}
