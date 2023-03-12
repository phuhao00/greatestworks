package email

import (
	"greatestworks/business/module"
)

var (
	Mod *Module
)

func init() {
	module.MManager.RegisterModule("", Mod)
}

type Module struct {
}

func (m *Module) Loop() {
	//TODO implement me
	panic("implement me")
}

func (m *Module) Monitor() {
	//TODO implement me
	panic("implement me")
}
