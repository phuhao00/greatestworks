package task

import (
	"greatestworks/aop/event"
	"greatestworks/business/module/friend"
)

type EventHandle func(iEvent event.IEvent)

type EventWrap struct {
	Player
	event.IEvent
}

func (m *Module) OnEvent(event *EventWrap) {
	m.eventHandles[event](event)
}

func (m *Module) HandleAddOrDelFriendEvent(eventWrap *EventWrap) {
	e := eventWrap.IEvent.(*friend.AddOrDelFriendEvent)
	player := eventWrap.Player
	taskData := player.GetTaskData()
	taskGroup := taskData.GetTaskGroup(e.GetModuleName(), e.GetCategory())
	for _, task := range taskGroup {
		task.OnEvent(e)
	}
}
