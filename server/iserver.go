package server

type IService interface {
	Start()
	Reload()
	Init(config interface{}, processId int)
	Stop()
}
