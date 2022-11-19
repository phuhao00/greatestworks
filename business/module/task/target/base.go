package target

import "greatestworks/business/module/task"

type Base struct {
}

func (b *Base) GetCategory() task.TargetCategory {
	//TODO implement me
	panic("implement me")
}

func (b *Base) CheckDone() bool {
	//TODO implement me
	panic("implement me")
}

func (b *Base) OnEvent(event task.Event) {
	//TODO implement me
	panic("implement me")
}

