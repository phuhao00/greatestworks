package action

import (
	b3 "github.com/magicsea/behavior3go"
	b3Cfg "github.com/magicsea/behavior3go/config"
	b3Core "github.com/magicsea/behavior3go/core"
	"greatestworks/internal/gameplay/scene"
)

type RandomStatus struct {
	owner scene.IActor
	b3Core.Action
}

func (setStatus *RandomStatus) Initialize(setting *b3Cfg.BTNodeCfg) {
	setStatus.Action.Initialize(setting)
}

func (setStatus *RandomStatus) OnOpen(tick *b3Core.Tick) {
	owner := tick.Blackboard.Get("owner", tick.GetTree().GetID(), "")
	setStatus.owner = owner.(scene.IActor)
}

func (setStatus *RandomStatus) OnTick(tick *b3Core.Tick) b3.Status {
	return setStatus.owner.RandomStatus()
}
