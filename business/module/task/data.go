package task

import "sync"

type Data struct {
	Tasks        sync.Map
	Achievements sync.Map
}

func NewData() *Data {
	return &Data{
		Tasks:        sync.Map{},
		Achievements: sync.Map{},
	}
}

func (d *Data) ToDB() {

}

func (d *Data) LoadFromDB() {

}

func (d *Data) GetTask(id uint64) Task {
	value, ok := d.Tasks.Load(id)
	if ok {
		return value.(Task)
	}
	return nil
}

func (d *Data) SyncAllTasks(player Player) {
	//todo
}
