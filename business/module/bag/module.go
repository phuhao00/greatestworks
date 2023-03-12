package bag

import (
	"greatestworks/aop/event"
	"greatestworks/business/module"
	"sync"
)

func init() {
	module.MManager.RegisterModule("", GetMe())
}

type Module struct {
}

var (
	onceInitMod sync.Once
	Mod         *Module
)

func GetMe() *Module {
	onceInitMod.Do(func() {
		Mod = &Module{}
	})
	return Mod
}

func (m Module) OnEvent(c module.Character, event event.IEvent) {
	//TODO implement me
	panic("implement me")
}

func (m Module) SetEventCategoryActive(eventCategory int) {
	//TODO implement me
	panic("implement me")
}
