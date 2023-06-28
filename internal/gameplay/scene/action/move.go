package action

import (
	b3 "github.com/magicsea/behavior3go"
	b3Cfg "github.com/magicsea/behavior3go/config"
	b3Core "github.com/magicsea/behavior3go/core"
	"greatestworks/internal/gameplay/scene"
)

type Move struct {
	b3Core.Action
	owner scene.IActor
}

func (m *Move) Initialize(setting *b3Cfg.BTNodeCfg) {
	m.Action.Initialize(setting)
}

func (m *Move) OnOpen(tick *b3Core.Tick) {
	owner := tick.Blackboard.Get("owner", tick.GetTree().GetID(), "")
	m.owner = owner.(scene.IActor)
}

func (m *Move) OnTick(tick *b3Core.Tick) b3.Status {
	return m.owner.MoveToTarget()
}
