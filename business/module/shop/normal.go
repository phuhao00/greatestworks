package shop

type Normal struct {
	Id       uint32
	Category Category
	Pool     []Item
}

func (n *Normal) Refresh() {

}
