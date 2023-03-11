package task

import (
	"greatestworks/business/module/hub"
	"sync"
)

type Data struct {
	Tasks        sync.Map
	Achievements sync.Map
	hub.DataAsSubscriber
}

type GroupKey struct {
	ModuleName string
	Category   int
}

func NewTaskData() *Data {
	return &Data{
		Tasks:        sync.Map{},
		Achievements: sync.Map{},
	}
}

func (d *Data) ToDB() {

}

func (d *Data) LoadFromDB() {

}

func (d *Data) GetTaskGroup(moduleName string, category int) []ITask {
	value, ok := d.Tasks.Load(GroupKey{ModuleName: moduleName, Category: category})
	if ok {
		return value.([]ITask)
	}
	return nil
}

func (d *Data) SyncAllTasks(player Player) {
	//todo
}
