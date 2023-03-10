package module

import "greatestworks/aop/event"

type DataBase struct {
}

func (m *DataBase) RegisterListener(e event.Enum, cb event.OnEvent) {
	//TODO implement me
	panic("implement me")
}

func (m *DataBase) Dispatch(e event.Enum, params ...interface{}) {
	//TODO implement me
	panic("implement me")
}
