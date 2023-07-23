package vip

type Vip struct {
	Level uint32
	Exp   uint32
}

func NewVip() *Vip {
	return &Vip{
		Level: 0,
		Exp:   0,
	}
}

func (v *Vip) Load() {
	//TODO implement me
	panic("implement me")
}

func (v *Vip) Save() {
	//TODO implement me
	panic("implement me")
}

func (v *Vip) AddExp(delta uint32) {
	v.Exp += delta
	//todo check can upgrade level
}

func (v *Vip) getReward() {

}
