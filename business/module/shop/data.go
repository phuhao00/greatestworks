package shop

type Data struct {
	Shops map[uint32]Shop
}

func NewData() *Data {
	return &Data{
		Shops: nil,
	}
}
