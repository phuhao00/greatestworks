package item

type Item interface {
	Add(delta int64)
	Delete(delta int64)
	GetNum() int64
	GetId() uint32
}

type Trans interface {
	ToPB()
}

type ResultEffect interface {
	AddAttr() //添加属性
	Extra()   //额外的
}

// Limit 限制
type Limit interface {
	DayGetCheck() bool
	WeekGetCheck() bool
	UseCdCheck() bool
}
