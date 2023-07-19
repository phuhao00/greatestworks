package chat

type SystemMsgHandler struct {
	Consumer *nsq.Consumer
	Handler
	SystemTransfer
}
