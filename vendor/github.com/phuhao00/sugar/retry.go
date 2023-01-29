package sugar

import (
	"sync"
	"time"
)

type Debounce struct {
	after     time.Duration
	mu        *sync.Mutex
	timer     *time.Timer
	done      bool
	callbacks []func()
}

func NewDebounce(duration time.Duration, fns ...func()) (func(), func()) {
	d := &Debounce{
		after:     duration,
		mu:        new(sync.Mutex),
		timer:     nil,
		done:      false,
		callbacks: fns,
	}
	return func() {
		d.reset()
	}, d.cancel
}

func (d *Debounce) reset() *Debounce {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.done {
		return d
	}
	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.after, func() {
		for _, cb := range d.callbacks {
			cb()
		}
	})
	return d
}

func (d *Debounce) cancel() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
	d.done = true
}

// AttemptWithDelay invokes a function N times until it returns valid output,
// with a pause between each call. Returning either the caught error or nil.
// When first argument is less than `1`, the function runs until a successful
// response is returned.

func AttemptWithDelay(maxIteration int, delay time.Duration, f func(int, time.Duration) error) (int, time.Duration, error) {
	var err error
	start := time.Now()
	for i := 0; maxIteration <= 0 || i < maxIteration; i++ {
		err = f(i, time.Since(start))
		if err == nil {
			return i + 1, time.Since(start), nil
		}
		if maxIteration <= 0 || i+1 < maxIteration {
			time.Sleep(delay)
		}
	}
	return maxIteration, time.Since(start), err
}
