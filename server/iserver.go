package server

type Processor interface {
	Start()
	Stop()
}

type IManager interface {
	Loop()    //处理内部消息转发
	Monitor() //处理外部消息进入
}

type ISystem interface {
}

type CRUD interface {
}
