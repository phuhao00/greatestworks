package bag

import "greatestworks/business/module/bag/item"

type Bag interface {
	AddItem(item item.Item)
	DelItem(item item.Item)
}
