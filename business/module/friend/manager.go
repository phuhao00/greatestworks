package friend

import (
	"greatestworks/aop/event"
	"greatestworks/business/module"
)

func GetName() string {
	return ""
}

type Manager struct {
}

func (s *Manager) OnEvent(player module.Character, event event.IEvent) {

}

func (s *Manager) SetEventCategoryActive(eventCategory int) {

}
