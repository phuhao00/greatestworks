package activity

import (
	"greatestworks/business/module/base"
	"sync"
)

type Manager struct {
	*base.MetricsBase
	*base.DBActionBase
}

var (
	instance *Manager
	onceInit sync.Once
)

func GetMe() *Manager {
	onceInit.Do(func() {
		instance = &Manager{}
	})
	return instance
}

func (a *Manager) OnStart() {
	//TODO implement me
	panic("implement me")
}

func (a *Manager) AfterStart() {
	//TODO implement me
	panic("implement me")
}

func (a *Manager) OnStop() {
	//TODO implement me
	panic("implement me")
}

func (a *Manager) AfterStop() {
	//TODO implement me
	panic("implement me")
}
