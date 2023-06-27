package skill

import (
	"greatestworks/internal"
	"sync"
)

const (
	ModuleName = "skill"
)

var (
	Mod         *Module
	OnceInitMod sync.Once
)

func init() {
	internal.ModuleManager.RegisterModule(ModuleName, GetMod())
}

type Module struct {
	*internal.BaseModule
}

func GetMod() *Module {
	OnceInitMod.Do(func() {
		Mod = &Module{
			BaseModule: internal.NewBaseModule(),
		}
	})
	return Mod
}

func (m *Module) SetName(name string) {
	m.BaseModule.SetName(name)
}
