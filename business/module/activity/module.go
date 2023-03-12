package activity

import (
	"greatestworks/business/module"
	"sync"
)

func init() {
	module.MManager.RegisterModule("", GetMe())
}

var (
	onceInit sync.Once
	Mod      *Module
)

func GetMe() *Module {
	onceInit.Do(func() {
		Mod = &Module{}
	})
	return Mod
}

type Module struct {
	*module.MetricsBase
	*module.DBActionBase
}

func (a *Module) OnStart() {
	//TODO implement me
	panic("implement me")
}

func (a *Module) AfterStart() {
	//TODO implement me
	panic("implement me")
}

func (a *Module) OnStop() {
	//TODO implement me
	panic("implement me")
}

func (a *Module) AfterStop() {
	//TODO implement me
	panic("implement me")
}
