package pubsub

import (
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/samber/lo"
)

type RetryConsistency int

const (
	ConsistentRetry           RetryConsistency = 0 // slow
	EventuallyConsistentRetry RetryConsistency = 1 // fast, at *least* once
)

func getAttemptsFromHeaders(msg *amqp.Delivery) int {
	result := 0

	ok := lo.Try0(func() {
		header, ok := msg.Headers["x-retry-attempts"]
		if !ok {
			return
		}

		attempts, err := strconv.ParseInt(header.(string), 10, 32)
		if err != nil {
			return
		}

		result = int(attempts)
	})
	if !ok {
		logger("could not parse x-retry-attempts header")
		return 0
	}

	return result
}

type retryAcknowledger struct {
	retryProducer    *Producer
	retryQueue       string
	retryer          RetryStrategy
	retryConsistency RetryConsistency
	msg              amqp.Delivery
	parent           amqp.Acknowledger
}

func newRetryAcknowledger(retryProducer *Producer, retryQueue string, retryer RetryStrategy, retryConsistency RetryConsistency, msg amqp.Delivery) amqp.Acknowledger {
	return &retryAcknowledger{
		retryProducer:    retryProducer,
		retryQueue:       retryQueue,
		retryer:          retryer,
		retryConsistency: retryConsistency,
		msg:              msg,
		parent:           msg.Acknowledger,
	}
}

func (a *retryAcknowledger) Ack(tag uint64, multiple bool) error {
	return a.parent.Ack(tag, multiple)
}

func (a *retryAcknowledger) Nack(tag uint64, multiple bool, requeue bool) error {
	if multiple {
		panic("multiple nack is not available with retry strategy")
	}

	if requeue {
		panic("requeue is not available with retry strategy")
	}

	attempts := getAttemptsFromHeaders(&a.msg)

	ttl, ok := a.retryer.NextBackOff(&a.msg, attempts)
	if ok {
		return a.retry(tag, attempts, ttl)
	}

	return a.parent.Nack(tag, false, requeue)
}

func (a *retryAcknowledger) Reject(tag uint64, requeue bool) error {
	if requeue {
		panic("requeue is not available with retry strategy")
	}

	attempts := getAttemptsFromHeaders(&a.msg)

	ttl, ok := a.retryer.NextBackOff(&a.msg, attempts)
	if ok {
		return a.retry(tag, attempts, ttl)
	}

	return a.parent.Reject(tag, requeue)
}

func (a *retryAcknowledger) retry(tag uint64, attempts int, ttl time.Duration) error {
	headers := a.msg.Headers
	if headers == nil {
		headers = amqp.Table{}
	}

	headers["x-retry-attempts"] = strconv.FormatInt(int64(attempts+1), 10)

	if _, ok := headers["x-first-retry-exchange"]; !ok {
		headers["x-first-retry-exchange"] = a.msg.Exchange
	}

	if _, ok := headers["x-first-retry-routing-key"]; !ok {
		headers["x-first-retry-routing-key"] = a.msg.RoutingKey
	}

	msg := amqp.Publishing{
		Headers:         headers,
		ContentType:     a.msg.ContentType,
		ContentEncoding: a.msg.ContentEncoding,
		DeliveryMode:    a.msg.DeliveryMode,
		Priority:        a.msg.Priority,
		CorrelationId:   a.msg.CorrelationId,
		ReplyTo:         a.msg.ReplyTo,
		Expiration:      strconv.FormatInt(ttl.Milliseconds(), 10),
		MessageId:       a.msg.MessageId,
		Timestamp:       a.msg.Timestamp,
		Type:            a.msg.Type,
		UserId:          a.msg.UserId,
		AppId:           a.msg.AppId,
		Body:            a.msg.Body,
	}

	switch a.retryConsistency {
	case ConsistentRetry:
		err := a.retryProducer.channel.Tx()
		if err != nil {
			return err
		}

		err = a.retryProducer.Publish(a.retryQueue, true, false, msg)
		if err != nil {
			a.retryProducer.channel.TxRollback()
			return err
		}

		err = a.parent.Ack(tag, false)
		if err != nil {
			a.retryProducer.channel.TxRollback()
			return err
		}

		return a.retryProducer.channel.TxCommit()

	case EventuallyConsistentRetry:
		err := a.retryProducer.Publish(a.retryQueue, true, false, msg)
		if err != nil {
			return err
		}

		return a.parent.Ack(tag, false)

	default:
		panic("unsupported retry consistency")
	}
}
