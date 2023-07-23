package action

import (
	b3 "github.com/magicsea/behavior3go"
	b3Cfg "github.com/magicsea/behavior3go/config"
	b3Core "github.com/magicsea/behavior3go/core"
	"greatestworks/internal/gameplay/scene"
)

type DisappearAction struct {
	b3Core.Action
	owner scene.IActor
}

func (action *DisappearAction) Initialize(setting *b3Cfg.BTNodeCfg) {
	action.Action.Initialize(setting)
}

func (action *DisappearAction) OnOpen(tick *b3Core.Tick) {
	owner := tick.Blackboard.Get("owner", tick.GetTree().GetID(), "")
	action.owner = owner.(scene.IActor)
}

func (action *DisappearAction) OnTick(tick *b3Core.Tick) b3.Status {
	return action.owner.DisappearTrigger()
}
