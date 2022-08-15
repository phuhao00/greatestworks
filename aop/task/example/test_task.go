package example

import "greatestworks/aop/task"

type TTask struct {
	Conf    *task.Config
	Next    *TTask
	Status  task.Status
	Targets []task.Target
}

func NewTTask(config *task.Config) *TTask {
	t := &TTask{
		Conf: config,
	}
	return t

}

func (t *TTask) Accept(config *task.Config) {
	t.Status = task.ACCEPT
}

func (t *TTask) Finish() {
	t.Status = task.FINISH

}

func (t *TTask) TargetDoneCallBack() {
	count := 0
	for _, target := range t.Targets {
		if target.CheckDone() {
			count++
		}
	}
	if count == len(t.Targets) {
		t.Finish()
	}
}
