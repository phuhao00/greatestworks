package instance

import "greatestworks/business/module/condition"

func NewItemNum(str string) condition.Condition {
	return &ItemNum{}
}

func init() {
	condition.GetMe().Reg(111, NewItemNum)
}

type ItemNum struct {
	condition.Base
	ItemId uint32
	Num    uint64
}

func (i *ItemNum) CheckArrived() bool {

	return false
}
