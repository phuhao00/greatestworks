package pubsub

import (
	"sync"

	"github.com/samber/lo"
)

/**
 * See https://github.com/samber/lo/pull/270
 */

type rpc[T any, R any] struct {
	C chan lo.Tuple2[T, func(R)]
}

// NewRPC synchronizes goroutines for a bidirectionnal request-response communication.
func newRPC[T any, R any](ch chan<- T) *rpc[T, R] {
	return &rpc[T, R]{
		C: make(chan lo.Tuple2[T, func(R)]),
	}
}

// Send blocks until response is triggered.
func (rpc *rpc[T, R]) Send(request T) R {
	done := make(chan R)
	defer close(done)

	once := sync.Once{}

	rpc.C <- lo.T2(request, func(response R) {
		once.Do(func() {
			done <- response
		})
	})

	return <-done
}

/**
 * See https://github.com/samber/lo/pull/268
 */

// SafeClose protects against double-close panic.
// Returns true on first close, only.
// May be equivalent to calling `sync.Once{}.Do(func() { close(ch) })`.`
func safeClose[T any](ch chan<- T) (justClosed bool) {
	defer func() {
		if recover() != nil {
			justClosed = false
		}
	}()

	close(ch) // may panic
	return true
}
