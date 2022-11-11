package actor

type Actor interface {
	OnDamage(delta int64)
	Attack()
	OnMove()
}

type Base struct {
	Hp     int64
	Damage int64
	Pos    []float64
}

func (b Base) OnDamage(delta int64) {
	//TODO implement me
	panic("implement me")
}

func (b Base) Attack() {
	//TODO implement me
	panic("implement me")
}

func (b Base) OnMove() {
	//TODO implement me
	panic("implement me")
}
