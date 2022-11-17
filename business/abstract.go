package business

type Processor interface {
	Start()
	Stop()
}

type IManager interface {
	Loop()
	Monitor()
}

type ISystem interface {
}

type CRUD interface {
}
