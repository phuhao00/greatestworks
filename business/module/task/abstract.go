package task

type Abstract interface {
}

type Task interface {
	SetStatus(Status)
	TargetDoneCallBack()
	GetTargets() []Target
}

type Target interface {
	GetCategory() TargetCategory
	CheckDone() bool
	OnNotify(param interface{})
	GetTargetConfigId() uint32
}

type Event interface {
	Attach(targetId uint32) error
	Detach(id uint64) error
	Notify(param interface{}, player Player)
}
