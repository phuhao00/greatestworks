package event

import "greatestworks/business/module/condition"

type Event interface {
	Notify()
	Attach(condition condition.Condition)
	Detach(id uint32)
}
