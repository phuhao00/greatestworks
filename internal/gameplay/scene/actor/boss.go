package actor

type Boss struct {
	Category int
	*Base
	real BossReal
}

func NewBoss() *Boss {
	return &Boss{
		Category: 0,
		Base: &Base{
			Hp:     0,
			Damage: 0,
		},
	}
}

func (b *Boss) OnDamage(delta int64) {
}

func (b *Boss) Attack() {
}

func (b *Boss) OnMove() {
}
