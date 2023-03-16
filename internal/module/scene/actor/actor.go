package actor

type Actor interface {
	OnDamage(delta int64)
	Attack()
	OnMove()
}

type Base struct {
	Hp     int64     `json:"hp"`
	Damage int64     `json:"damage"`
	Pos    []float64 `json:"pos"`
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
