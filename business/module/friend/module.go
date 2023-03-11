package friend

import (
	"greatestworks/aop/event"
	"greatestworks/business/module"
)

var (
	Mod *Module
)

func init() {
	module.MManager.RegisterModule("", Mod)
}

func GetName() string {
	return ""
}

type Module struct {
}

func (s *Module) OnEvent(player module.Character, event event.IEvent) {

}

func (s *Module) SetEventCategoryActive(eventCategory int) {

}
