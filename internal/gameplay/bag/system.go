package bag

type System struct {
	Normal Bag
}

func NewSystem() *System {
	return &System{
		Normal: nil,
	}
}
