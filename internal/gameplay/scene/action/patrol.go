package action

import (
	b3 "github.com/magicsea/behavior3go"
	b3Cfg "github.com/magicsea/behavior3go/config"
	b3Core "github.com/magicsea/behavior3go/core"
	"greatestworks/internal/gameplay/scene"
)

type Patrol struct {
	owner scene.IActor
	b3Core.Action
}

func (pa *Patrol) Initialize(setting *b3Cfg.BTNodeCfg) {
	pa.Action.Initialize(setting)
}

func (pa *Patrol) OnOpen(tick *b3Core.Tick) {
	owner := tick.Blackboard.Get("owner", tick.GetTree().GetID(), "")
	pa.owner = owner.(scene.IActor)
}

func (pa *Patrol) OnTick(tick *b3Core.Tick) b3.Status {
	if pa.owner != nil {
		return pa.owner.Patrol()
	}
	return b3.SUCCESS
}
