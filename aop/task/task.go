package task

type Task struct {
	Conf   *Config
	Next   *Task
	Status Status
}

func NewTask(config *Config) *Task {
	t := &Task{
		Conf: config,
	}
	return t

}

func (t *Task) Accept(config *Config) {
	t.Status = ACCEPT
}

func (t *Task) Finish() {
	t.Status = FINISH

}
