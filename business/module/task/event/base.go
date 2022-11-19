package event

import "greatestworks/business/module/task"

type Base struct {
	targets []task.Target
}

func (b *Base) Attach(target task.Target) error {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Detach(id uint64) error {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Notify() {
	//TODO implement me
	panic("implement me")
}
