package chat

import (
	"github.com/phuhao00/greatestworks-proto/module"
	"greatestworks/internal"
	"sync"
)

const (
	ModuleName = "chat"
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
}

func GetMod() *Module {
	Mod = &Module{internal.NewBaseModule()}

	return Mod
}

func (m *Module) GetName() string {
	return module.Module_Chat.String()
}
