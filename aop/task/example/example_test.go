package example

import (
	"fmt"
	"greatestworks/aop/task"
	"testing"
)

func TestName(t *testing.T) {
	te := TEvent{
		Subscribers: make([]task.Target, 0),
	}
	tg := &TTarget{
		Id:   111,
		Data: 1,
	}
	te.Attach(tg)
	te.Data = 1
	te.Notify()
	fmt.Println("CheckDone:", tg.CheckDone())
}
