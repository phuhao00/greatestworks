package condition

import (
	"greatestworks/aop/condition/event"
)

type Condition interface {
	CheckArrived() bool
	OnNotify(event.Event)
	GetId() uint32
	SetCB(func())
}

type Base struct {
	Cb func()
}

func NewTargetBase() *Base {
	return &Base{}
}

func (t *Base) CheckArrived() bool {
	return false
}

func (t *Base) OnNotify(event event.Event) {

}

func (t *Base) GetId() uint32 {
	return 0
}

func (t *Base) SetCB(f func()) {
	t.Cb = f
}
