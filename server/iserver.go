package server

type IServer interface {
	Start()
	Loop()    //处理内部消息转发
	Monitor() //处理外部消息进入
	Stop()
}
