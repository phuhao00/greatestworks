package main

import (
	"fmt"
	"runtime"
	"time"
)

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor) Run(c *Obj) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}
func runJanitor(c *Obj, ci time.Duration) {
	j := &janitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	c.janitor = j
	go j.Run(c)
}

type Obj struct {
	janitor *janitor
}

func NewObj() *Obj {
	return &Obj{}
}

func (o *Obj) DeleteExpired() {
	fmt.Println("DeleteExpired")
}

func stopJanitor(c *Obj) {
	c.janitor.stop <- true
}

func newCacheWithJanitor(ci time.Duration) *Obj {
	c := NewObj()

	if ci > 0 {
		runJanitor(c, ci)
		runtime.SetFinalizer(c, stopJanitor)
	}
	return c
}

func main() {
	newCacheWithJanitor(time.Second * 3)
}
