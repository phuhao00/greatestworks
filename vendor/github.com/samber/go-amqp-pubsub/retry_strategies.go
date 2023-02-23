package pubsub

import (
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RetryStrategy interface {
	NextBackOff(*amqp.Delivery, int) (time.Duration, bool)
}

type ConstantRetryStrategy struct {
	maxRetry int
	interval time.Duration
}

func NewConstantRetryStrategy(maxRetry int, interval time.Duration) RetryStrategy {
	return &ConstantRetryStrategy{
		maxRetry: maxRetry,
		interval: interval,
	}
}

func (rs *ConstantRetryStrategy) NextBackOff(msg *amqp.Delivery, attempts int) (time.Duration, bool) {
	if attempts >= rs.maxRetry {
		return 0, false
	}

	return rs.interval, true
}

type ExponentialRetryStrategy struct {
	maxRetry           int
	initialInterval    time.Duration
	intervalMultiplier float64
}

func NewExponentialRetryStrategy(maxRetry int, initialInterval time.Duration, intervalMultiplier float64) RetryStrategy {
	return &ExponentialRetryStrategy{
		maxRetry:           maxRetry,
		initialInterval:    initialInterval,
		intervalMultiplier: intervalMultiplier,
	}
}

func (rs *ExponentialRetryStrategy) NextBackOff(msg *amqp.Delivery, attempts int) (time.Duration, bool) {
	if attempts >= rs.maxRetry {
		return 0, false
	}

	ns := float64(rs.initialInterval.Nanoseconds()) * math.Pow(rs.intervalMultiplier, float64(attempts))
	return time.Duration(ns), true
}
