package shop

type Mystery struct {
	Id       uint32 //configId
	Category Category
	Pool     []Item
}

func (m *Mystery) Refresh() {

}
