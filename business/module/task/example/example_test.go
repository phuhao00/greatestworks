package example

import (
	"fmt"
	"greatestworks/business/module/task"
	"testing"
)

func TestEvent(t *testing.T) {
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

func TestTask(t *testing.T) {
	NewTTask(nil)

}
