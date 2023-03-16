package family

import (
	module2 "greatestworks/internal/module"
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
	module2.MManager.RegisterModule(ModuleName, GetMod())
}

type Module struct {
	*module2.BaseModule
	families map[uint64]*Family
	IWorld
	ChIn  chan ManagerHandlerParam
	ChOut chan interface{}
}

func GetMod() *Module {
	Mod = &Module{BaseModule: module2.NewBaseModule()}
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
