package example

import (
	"greatestworks/business/module/condition"
	"greatestworks/business/module/condition/event"
)

type TTarget struct {
	Id   uint32
	Data int
	Done bool
	*condition.Base
}

func NewTTarget() *TTarget {
	tt := &TTarget{
		Id:   0,
		Data: 0,
		Done: false,
		Base: condition.NewTargetBase(),
	}
	return tt
}

func (T TTarget) CheckArrived() bool {
	return T.Done
}

func (T *TTarget) OnNotify(event event.Event) {
	e := event.(*TEvent)
	if e.Data == T.Data {
		T.Done = true
	}
	if T.Done {
		T.Cb()
	}
}

func (T TTarget) GetId() uint32 {
	return T.Id
}
