package plant

import (
	"greatestworks/aop/event"
	"greatestworks/business/module"
)

type Module struct {
}

func (m Module) OnEvent(c module.Character, event event.IEvent) {
	//TODO implement me
	panic("implement me")
}

func (m Module) SetEventCategoryActive(eventCategory int) {
	//TODO implement me
	panic("implement me")
}
