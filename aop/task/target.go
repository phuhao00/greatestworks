package task

type Target interface {
	CheckDone() bool
	OnNotify(Event)
	GetTargetId() uint32
	SetTaskCB(func())
}
