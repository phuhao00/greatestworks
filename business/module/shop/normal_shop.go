package shop

type Normal struct {
	Id          uint32
	Category    Category
	RefreshTime int64
	Items       map[uint32]Item
}

func (n *Normal) Refresh() {

}
