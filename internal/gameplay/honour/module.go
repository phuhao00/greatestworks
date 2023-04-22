package honour

import (
	"greatestworks/internal"
	"greatestworks/internal/note/event"
)

type Module struct {
	internal.BaseModule
}

func (m *Module) onEvent(event event.IEvent) {
	//todo
}
