package sugar

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitSignal(fn func(sig os.Signal) bool) {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGPIPE)

	for sig := range ch {
		if !fn(sig) {
			close(ch)
			break
		}
	}
}
