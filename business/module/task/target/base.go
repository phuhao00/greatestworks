package target

import "greatestworks/business/module/task"

type Base struct {
	ConfigId uint32
}

func (b *Base) GetCategory() task.TargetCategory {
	//TODO implement me
	panic("implement me")
}

func (b *Base) CheckDone() bool {
	//TODO implement me
	panic("implement me")
}

func (b *Base) OnNotify(param interface{}) {
	//TODO implement me
	panic("implement me")
}

func (b *Base) GetTargetConfigId() uint32 {
	return b.ConfigId
}
