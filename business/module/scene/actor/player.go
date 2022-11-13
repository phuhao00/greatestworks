package actor

import "google.golang.org/protobuf/proto"

type Player struct {
	*Base
	real  PlayerReal
	Title string
}

func NewPlayer() *Player {
	return &Player{
		Base: &Base{
			Hp:     0,
			Damage: 0,
		},
	}
}

func (p *Player) OnDamage(delta int64) {
	p.Hp -= delta
}

func (p *Player) OnAttack() {
	//TODO implement me
	panic("implement me")
}

func (p *Player) OnMove() {
	//TODO implement me
	panic("implement me")
}

func (p *Player) SendMsg(message proto.Message) {
	p.real.SendMsg(message)
}
