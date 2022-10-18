package instance

import (
	"greatestworks/business/module/bag/item"
	"greatestworks/business/module/bag/item/template"
	"sync"
)

type Normal struct {
	Data sync.Map
}

func (n *Normal) AddItem(item item.Item) {
	value, ok := n.Data.Load(item.GetId())
	if ok {
		value.(*template.ItemBase).Add(item.GetNum())
	}
}

func (n *Normal) DelItem(item item.Item) {
	value, ok := n.Data.Load(item.GetId())
	if ok {
		value.(*template.ItemBase).Delete(item.GetNum())
	}
}
