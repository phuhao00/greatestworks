package action

import (
	b3 "github.com/magicsea/behavior3go"
	b3Cfg "github.com/magicsea/behavior3go/config"
	b3Core "github.com/magicsea/behavior3go/core"
	"greatestworks/internal/gameplay/scene"
)

type Appear struct {
	b3Core.Action
	owner scene.IActor
}

func (a *Appear) Initialize(setting *b3Cfg.BTNodeCfg) {
	a.Action.Initialize(setting)
}

func (a *Appear) OnOpen(tick *b3Core.Tick) {
	owner := tick.Blackboard.Get("owner", tick.GetTree().GetID(), "")
	a.owner = owner.(scene.IActor)
}

func (a *Appear) OnTick(tick *b3Core.Tick) b3.Status {
	return a.owner.AppearTrigger()
}
