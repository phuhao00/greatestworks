package handlers

import (
	"context"
	"reflect"

	"greatestworks/internal/domain/replication"
	"greatestworks/internal/events"
	"greatestworks/internal/infrastructure/logging"
)

// RegisterReplicationSubscribers subscribes to replication-related events and logs them.
func RegisterReplicationSubscribers(bus *events.EventBus, logger logging.Logger) {
	types := []string{
		reflect.TypeOf(&replication.PlayerJoinedEvent{}).String(),
		reflect.TypeOf(&replication.PlayerLeftEvent{}).String(),
		reflect.TypeOf(&replication.InstanceStartedEvent{}).String(),
		reflect.TypeOf(&replication.InstanceFullEvent{}).String(),
		reflect.TypeOf(&replication.InstanceProgressUpdatedEvent{}).String(),
		reflect.TypeOf(&replication.InstanceCompletedEvent{}).String(),
		reflect.TypeOf(&replication.InstanceClosingEvent{}).String(),
		reflect.TypeOf(&replication.InstanceClosedEvent{}).String(),
	}

	for _, t := range types {
		bus.Subscribe(t, func(ctx context.Context, e events.Event) error {
			logger.Info("replication event", logging.Fields{"type": t, "data": e.GetData()})
			return nil
		})
	}
}
