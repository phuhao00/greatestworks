package chat

import "github.com/nsqio/go-nsq"

type PrivateChat struct {
	Consumer nsq.Consumer
	Handler
	PrivateTransfer
}
