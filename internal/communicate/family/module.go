package family

import (
	"github.com/phuhao00/greatestworks-proto/module"
	"greatestworks/internal"
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
	internal.MManager.RegisterModule(ModuleName, GetMod())
}

type Module struct {
	*internal.BaseModule
	families map[uint64]*Family
	IWorld
	ChIn  chan ManagerHandlerParam
	ChOut chan interface{}
}

func GetMod() *Module {
	Mod = &Module{BaseModule: internal.NewBaseModule()}
	return Mod
}

func (m *Module) Reload() {
	for {
		select {
		case msg := <-m.ChOut:
			m.ForwardMsg(msg)
		}
	}
}

func (m *Module) Init() {
	for {
		select {
		case param := <-m.ChIn:
			m.Handler(param)
		}
	}
}

func (m *Module) GetName() string {
	return module.Module_Family.String()
}
