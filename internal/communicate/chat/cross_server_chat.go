package chat

import "github.com/nsqio/go-nsq"

type CrossSrvChatHandler struct {
	Consumer *nsq.Consumer
	Handler
	ServerTransfer
}
