package example

import "greatestworks/aop/task"

type TTarget struct {
	Id   uint32
	Data int
	Done bool
}

func NewTTarget() {

}

func (T TTarget) CheckDone() bool {
	return T.Done
}

func (T *TTarget) OnNotify(event task.Event) {
	e := event.(*TEvent)
	if e.Data == T.Data {
		T.Done = true
	}
}

func (T TTarget) GetTargetId() uint32 {
	return T.Id
}

func (T TTarget) SetTaskCB(fn func()) {

}
