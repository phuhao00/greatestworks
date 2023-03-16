package bag

import (
	"greatestworks/internal/module/bag/item"
)

type Bag interface {
	AddItem(item item.Item)
	DelItem(item item.Item)
}
