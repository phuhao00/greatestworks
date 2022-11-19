package example

import (
	"greatestworks/business/module/condition"
	task2 "greatestworks/business/module/task"
)

type TTask struct {
	Conf    *task2.Config
	Next    *TTask
	Status  task2.Status
	Targets []condition.Condition
}

func NewTTask(config *task2.Config) *TTask {
	t := &TTask{
		Conf: config,
	}
	return t

}

func (t *TTask) SetStatus(status task2.Status) {
	t.Status = status
}

func (t *TTask) TargetDoneCallBack() {
	count := 0
	for _, target := range t.Targets {
		if target.CheckArrived() {
			count++
		}
	}
	if count == len(t.Targets) {
		t.SetStatus(task2.FINISH)
	}
}

func (t *TTask) GetTargets() []task2.Target {
	return nil
}
