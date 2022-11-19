package event

import (
	"greatestworks/business/module/task"
)

var (
	base *Base
)

func NewBase() *Base {
	return &Base{
		targets: make([]uint32, 0, 100),
	}
}

type Base struct {
	targets []uint32
}

func (b *Base) GetEventCategory() task.EventCategory {
	return task.BaseEvent
}

func (b *Base) Attach(targetId uint32) error {
	b.targets = append(b.targets, targetId)
	return nil
}

func (b *Base) Detach(id uint64) error {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Notify(param interface{}, player task.Player) {
	data := player.GetTaskData()
	taskCache := data.GetTaskCache(b.GetEventCategory())
	for _, targetId := range b.targets {
		for _, taskId := range taskCache {
			targets := data.GetTask(taskId).GetTargets()
			for _, target := range targets {
				if target.GetTargetConfigId() == targetId {
					target.OnNotify(param)
				}
			}
		}
	}
}
