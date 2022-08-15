package task

type Target interface {
	CheckDone() bool
	OnNotify(Event)
	GetTargetId() uint32
	SetTaskCB(func())
}

type TargetBase struct {
	TaskCB func()
}

func NewTargetBase() *TargetBase {
	return &TargetBase{}
}

func (t *TargetBase) CheckDone() bool {
	return false
}

func (t *TargetBase) OnNotify(event Event) {

}

func (t TargetBase) GetTargetId() uint32 {
	return 0
}

func (t *TargetBase) SetTaskCB(f func()) {
	t.TaskCB = f
}
