package card

import (
	"errors"
	"fmt"
	"github.com/phuhao00/greatestworks-proto/messageId"
	"github.com/phuhao00/greatestworks-proto/purchase"
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
		messageId.MessageId_CSCardAction,
		Action,
	}
}

// Action  ...
func Action(s *Data, packet *network.Message) {
	var req *purchase.CSCardAction
	var rsp = &purchase.SCCardAction{
		CardId: req.CardId,
		Action: req.Action,
	}
	err := proto.Unmarshal(packet.Data, req)
	if err != nil {
		return
	}
	s.Execute(getCardCategory(req.CardId), req.Action)
	fmt.Println(rsp)
}
