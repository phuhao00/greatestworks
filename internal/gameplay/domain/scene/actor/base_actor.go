package actor

import (
	b3 "github.com/magicsea/behavior3go"
	"greatestworks/internal/gameplay/scene"
)

type Base struct {
	Hp     int64     `json:"hp"`
	Damage int64     `json:"damage"`
	Pos    []float64 `json:"pos"`
}

func (b *Base) Patrol() b3.Status {
	panic("implement me")
}

func (b *Base) FollowTarget() scene.IActor {
	panic("implement me")
}

func (b *Base) Follow(actor scene.IActor) b3.Status {
	//TODO implement me
	panic("implement me")
}

func (b *Base) RandomStatus() b3.Status {
	//TODO implement me
	panic("implement me")
}

func (b *Base) AppearTrigger() b3.Status {
	//TODO implement me
	panic("implement me")
}

func (b *Base) DisappearTrigger() b3.Status {
	//TODO implement me
	panic("implement me")
}

func (b *Base) MoveToTarget() b3.Status {
	//TODO implement me
	panic("implement me")
}

func (b *Base) OnDamage(delta int64) {
}

func (b *Base) Attack() {
}

func (b *Base) OnMove() {
}
