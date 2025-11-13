package messaging

import (
	"greatestworks/internal/infrastructure/logging"
)

// EventLoggerAdapter adapts infrastructure logging.Logger to events.Logger
type EventLoggerAdapter struct{ logger logging.Logger }

func NewEventLoggerAdapter(logger logging.Logger) *EventLoggerAdapter {
	return &EventLoggerAdapter{logger: logger}
}

func (a *EventLoggerAdapter) Info(msg string, args ...interface{}) {
	a.logger.Info(msg, logging.Fields{"data": args})
}
func (a *EventLoggerAdapter) Debug(msg string, args ...interface{}) {
	a.logger.Debug(msg, logging.Fields{"data": args})
}
func (a *EventLoggerAdapter) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		if err, ok := args[0].(error); ok {
			a.logger.Error(msg, err, logging.Fields{"data": args[1:]})
			return
		}
	}
	a.logger.Error(msg, nil, logging.Fields{"data": args})
}
