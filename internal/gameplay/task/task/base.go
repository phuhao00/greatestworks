package task

import (
	"greatestworks/internal/gameplay/task"
	"greatestworks/internal/note/event"
)

type Base struct {
	Id       uint64
	ConfigId uint32
	Targets  []task.Target
}

func (b *Base) SetStatus(status task.Status) {
}

func (b *Base) OnEvent(event event.IEvent) {

}

func (b *Base) GetTaskData() task.ITaskData {
	return nil
}
