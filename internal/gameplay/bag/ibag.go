package bag

import (
	"greatestworks/internal/gameplay/bag/item"
)

type Bag interface {
	AddItem(item item.Item)
	DelItem(item item.Item)
}
