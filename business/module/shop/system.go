package shop

type System struct {
	Shops map[uint32]Shop
	Owner
}

func (s *System) SetOwner(owner Owner) {
	s.Owner = owner
}
