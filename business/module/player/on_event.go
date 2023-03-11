package player

import (
	"greatestworks/aop/event"
	"greatestworks/business/module"
)

func (p *Player) OnEvent(e event.IEvent) {
	module.MManager.GetModule(e.GetToModuleName()).OnEvent(p, e)
}
