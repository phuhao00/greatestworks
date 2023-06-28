package actor

type Base struct {
	Hp     int64     `json:"hp"`
	Damage int64     `json:"damage"`
	Pos    []float64 `json:"pos"`
}

func (b Base) OnDamage(delta int64) {
}

func (b Base) Attack() {
}

func (b Base) OnMove() {
}
