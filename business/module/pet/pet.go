package pet

type Pet struct {
	Category int
	ConfigId uint32
	Id       uint64
	Name     string
}

func (p *Pet) ToModel() *Model {
	return &Model{
		Category: p.Category,
		Id:       p.Id,
		ConfigId: p.ConfigId,
		Name:     p.Name,
	}
}
