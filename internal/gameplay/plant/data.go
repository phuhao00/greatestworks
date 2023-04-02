package plant

type Base struct {
	Id       uint64
	ConfId   uint32
	Status   Status
	SeedTime int64
}

func (b *Base) Update() {

}
