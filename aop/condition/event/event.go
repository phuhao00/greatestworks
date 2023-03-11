package event

import (
	"greatestworks/aop/condition"
)

type Event interface {
	Notify()
	Attach(condition condition.Condition)
	Detach(id uint32)
}
