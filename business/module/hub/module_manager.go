package hub

import "greatestworks/aop/event"

var (
	MManager ModuleManager
)

type GetEventFunc func(Category int) event.IEvent

type ModuleManager struct {
	ModuleName2ModuleGetEventFunc map[string]GetEventFunc
}

func (m *ModuleManager) AddModuleName2ModuleGetEventFunc(moduleName string, fn GetEventFunc) {
	m.ModuleName2ModuleGetEventFunc[moduleName] = fn
}

func (m *ModuleManager) GetEvent(moduleName string, Category int) event.IEvent {
	if fn, exist := m.ModuleName2ModuleGetEventFunc[moduleName]; exist {
		return fn(Category)
	}
	return nil
}
