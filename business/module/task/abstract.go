package task

type Abstract interface {
}

type Task interface {
	SetStatus(Status)
	TargetDoneCallBack()
}

type Target interface {
	GetCategory() TargetCategory
	CheckDone() bool
	OnEvent(event Event)
}

type Event interface {
	Attach(target Target) error
	Detach(id uint64) error
	Notify()
}
