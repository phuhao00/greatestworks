package family

import (
	"greatestworks/business/module"
	"sync"
)

const (
	ModuleName = "family"
)

var (
	Mod         *Module
	onceInitMod sync.Once
)

func init() {
	module.MManager.RegisterModule(ModuleName, GetMod())
}

type Module struct {
	*module.BaseModule
	families map[uint64]*Family
	IWorld
	ChIn  chan ManagerHandlerParam
	ChOut chan interface{}
}

func GetMod() *Module {
	Mod = &Module{BaseModule: module.NewBaseModule()}
	return Mod
}

func (m *Module) Loop() {
	for {
		select {
		case msg := <-m.ChOut:
			m.ForwardMsg(msg)
		}
	}
}

func (m *Module) Monitor() {
	for {
		select {
		case param := <-m.ChIn:
			m.Handler(param)
		}
	}
}
