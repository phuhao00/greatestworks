package task

import (
	"greatestworks/internal/note/event"
)

type BaseTask struct {
	Id       uint64
	ConfigId uint32
	Targets  []ITarget
}

func (b *BaseTask) SetStatus(status Status) {
}

func (b *BaseTask) OnEvent(event event.IEvent) {

}

func (b *BaseTask) GetTaskData() ITaskData {
	return nil
}
