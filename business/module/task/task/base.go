package task

import "greatestworks/business/module/task"

type Base struct {
	Id       uint64
	ConfigId uint32
}

func (b *Base) SetStatus(status task.Status) {
	//TODO implement me
	panic("implement me")
}

func (b *Base) TargetDoneCallBack() {
	//TODO implement me
	panic("implement me")
}
