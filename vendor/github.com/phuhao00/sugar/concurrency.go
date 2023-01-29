package sugar

import "sync"

type synchronize struct {
	locker sync.Locker
}

func (s *synchronize) Do(cb func() error) {
	s.locker.Lock()
	Try(cb)
	s.locker.Unlock()
}

func Synchronize(opt ...sync.Locker) synchronize {
	if len(opt) > 1 {
		panic("unexpected arguments")
	} else if len(opt) == 0 {
		opt = append(opt, &sync.Mutex{})
	}

	return synchronize{locker: opt[0]}
}

func Async[A any](f func() A) chan A {
	ch := make(chan A)

	go func() {
		ch <- f()
	}()

	return ch
}
