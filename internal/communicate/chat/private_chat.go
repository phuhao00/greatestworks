package chat

import (
	"github.com/nsqio/go-nsq"
	"github.com/phuhao00/greatestworks-proto/messageId"
	"google.golang.org/protobuf/proto"
)

type PrivateChat struct {
	Consumer nsq.Consumer
	Owner
}

func NewPrivateChat() *PrivateChat {
	return &PrivateChat{
		Consumer: nsq.Consumer{},
		Owner:    nil,
	}
}

func (p *PrivateChat) ForwardPlayer(message proto.Message) {
	//TODO implement me
	panic("implement me")
}

func (p *PrivateChat) InitNsqHandler(channel string) {
	//TODO implement me
	panic("implement me")
}

func (p *PrivateChat) HandleMessage(message nsq.Message) error {
	//TODO implement me
	panic("implement me")
}

func (p *PrivateChat) PublishChatMsg(chatMsg interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (p *PrivateChat) Stop() {
	//TODO implement me
	panic("implement me")
}

func (p *PrivateChat) SendMsg(ID messageId.MessageId, message proto.Message) {
	//TODO implement me
	panic("implement me")
}

func (p *PrivateChat) SetHandler(handler Handler) {

}
