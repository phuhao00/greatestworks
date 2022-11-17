package pet

type Model struct {
	ConfigId uint32
	Id       uint64
	Name     string
	Category int
}

func (m *Model) ToInstance() *Pet {
	return &Pet{
		Category: m.Category,
		ConfigId: m.ConfigId,
		Id:       m.Id,
		Name:     m.Name,
	}
}
