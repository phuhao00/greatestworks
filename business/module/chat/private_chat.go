package chat

import "github.com/nsqio/go-nsq"

type PrivateChat struct {
	Consumer nsq.Consumer
	Handler
	PrivateTransfer
}

func (p *PrivateChat) SetHandler(handler Handler) {
	p.Handler = handler
}
