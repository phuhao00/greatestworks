package card

import (
	"errors"
	"github.com/phuhao00/greatestworks-proto/messageId"
	"github.com/phuhao00/network"
	"google.golang.org/protobuf/proto"
	"sync"
)

type Handler struct {
	Id messageId.MessageId
	Fn func(s *Data, packet *network.Message)
}

var (
	handlers     []*Handler
	onceInit     sync.Once
	MinMessageId messageId.MessageId
	MaxMessageId messageId.MessageId //handle 的消息范围
)

func IsBelongToHere(id messageId.MessageId) bool {
	return id > MinMessageId && id < MaxMessageId
}

func GetHandler(id messageId.MessageId) (*Handler, error) {

	if id > MinMessageId && id < MaxMessageId {
		return nil, errors.New("not in")
	}
	for _, handler := range handlers {
		if handler.Id == id {
			return handler, nil
		}
	}
	return nil, errors.New("not exist")
}

func init() {
	onceInit.Do(func() {
		HandlerCardRegister()
	})
}

func HandlerCardRegister() {
	handlers[0] = &Handler{
		messageId.MessageId_CSCardBuy,
		Buy,
	}
	handlers[1] = &Handler{
		messageId.MessageId_CSCardRenew,
		Renew,
	}
	handlers[1] = &Handler{
		messageId.MessageId_CSCardDailyReceive,
		Receive,
	}

}

// Buy 购买
func Buy(s *Data, packet *network.Message) {
	var req proto.Message
	err := proto.Unmarshal(packet.Data, req)
	if err != nil {
		return
	}
}

// Renew   续费
func Renew(s *Data, packet *network.Message) {
	var req proto.Message
	err := proto.Unmarshal(packet.Data, req)
	if err != nil {
		return
	}
}

// Receive 领取每日奖励
func Receive(s *Data, packet *network.Message) {
	var req proto.Message
	err := proto.Unmarshal(packet.Data, req)
	if err != nil {
		return
	}
}
