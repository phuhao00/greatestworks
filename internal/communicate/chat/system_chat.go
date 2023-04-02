package chat

import "github.com/nsqio/go-nsq"

type SystemMsgHandler struct {
	Consumer *nsq.Consumer
	Handler
	SystemTransfer
}
