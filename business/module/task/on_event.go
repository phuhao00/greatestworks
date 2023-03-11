package task

import (
	"greatestworks/aop/event"
	"greatestworks/aop/event/friendevent"
	"greatestworks/business/module"
)

type EventHandle func(iEvent event.IEvent)

type EventWrap struct {
	Player
	event.IEvent
}

func (m *Module) HandleAddOrDelFriendEvent(eventWrap *EventWrap) {
	e := eventWrap.IEvent.(*friendevent.AddOrDelFriendEvent)
	player := eventWrap.Player
	taskData := player.GetTaskData()
	taskGroup := taskData.GetTaskGroup(e.GetFromModuleName(), e.GetCategory())
	for _, task := range taskGroup {
		task.OnEvent(e)
	}
}

func (m *Module) OnEvent(c module.Character, event event.IEvent) {
	//TODO implement me
	panic("implement me")
}

func (m *Module) SetEventCategoryActive(eventCategory int) {
	//TODO implement me
	panic("implement me")
}
