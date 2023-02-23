package pubsub

import "log"

var logger func(format string, values ...any) = DefaultLogger

func SetLogger(cb func(format string, values ...any)) {
	logger = cb
}

func DefaultLogger(format string, values ...any) {
	log.Printf(format, values...)
}
