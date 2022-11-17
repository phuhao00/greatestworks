package shop

type Manager struct {
	Shops map[uint32]Shop
	Owner
}

func (s *Manager) SetOwner(owner Owner) {
	s.Owner = owner
}
