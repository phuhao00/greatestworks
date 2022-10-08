package example

import (
	task2 "greatestworks/business/module/task"
)

type TTarget struct {
	Id   uint32
	Data int
	Done bool
	*task2.TargetBase
}

func NewTTarget() *TTarget {
	tt := &TTarget{
		Id:         0,
		Data:       0,
		Done:       false,
		TargetBase: task2.NewTargetBase(),
	}
	return tt
}

func (T TTarget) CheckDone() bool {
	return T.Done
}

func (T *TTarget) OnNotify(event task2.Event) {
	e := event.(*TEvent)
	if e.Data == T.Data {
		T.Done = true
	}
	if T.Done {
		T.TaskCB()
	}
}

func (T TTarget) GetTargetId() uint32 {
	return T.Id
}
