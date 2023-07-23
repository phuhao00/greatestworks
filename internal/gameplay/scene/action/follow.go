package action

import (
	b3 "github.com/magicsea/behavior3go"
	b3Cfg "github.com/magicsea/behavior3go/config"
	b3Core "github.com/magicsea/behavior3go/core"
	"greatestworks/internal/gameplay/scene"
)

type Follow struct {
	owner scene.IActor
	b3Core.Action
}

func (fa *Follow) Initialize(setting *b3Cfg.BTNodeCfg) {
	fa.Action.Initialize(setting)
}

func (fa *Follow) OnOpen(tick *b3Core.Tick) {
	owner := tick.Blackboard.Get("owner", tick.GetTree().GetID(), "")
	fa.owner = owner.(scene.IActor)
}

func (fa *Follow) OnTick(tick *b3Core.Tick) b3.Status {
	if fa.owner != nil && fa.owner.FollowTarget != nil {
		return fa.owner.Follow(fa.owner.FollowTarget())
	}
	return b3.SUCCESS
}
