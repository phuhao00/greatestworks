package chat

type CrossSrvChatHandler struct {
	Consumer *nsq.Consumer
	Handler
	ServerTransfer
}
