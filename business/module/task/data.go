package task

import "sync"

type Data struct {
	Tasks        sync.Map
	Achievements sync.Map
	taskCache    sync.Map // map[EventCategory][]uint64
}

func NewData() *Data {
	return &Data{
		Tasks:        sync.Map{},
		Achievements: sync.Map{},
		taskCache:    sync.Map{},
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

func (d *Data) AddTaskCache(eNum EventCategory, taskId uint64) {
	value, ok := d.taskCache.Load(eNum)
	uint64s := make([]uint64, 0)
	if ok {
		uint64s = value.([]uint64)
	}
	uint64s = append(uint64s, taskId)
	d.taskCache.Store(eNum, uint64s)
}

func (d *Data) GetTaskCache(eNum EventCategory) []uint64 {
	if value, ok := d.taskCache.Load(eNum); ok {
		return value.([]uint64)
	}
	return nil
}
