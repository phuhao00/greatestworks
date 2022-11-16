package chat

import (
	"errors"
	"fmt"
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"github.com/phuhao00/greatestworks-proto/gen/player"
	"github.com/phuhao00/network"
	"google.golang.org/protobuf/proto"
	"sync"
)

type PrivateChatHandler struct {
	Id messageId.MessageId
	Fn func(p *PrivateChat, packet *network.Message)
}

var (
	handlers                   []*PrivateChatHandler
	onceInit                   sync.Once
	MinMessageId, MaxMessageId messageId.MessageId
)

func init() {
	onceInit.Do(func() {
		HandlerChatRegister()
	})
}

func GetHandler(id messageId.MessageId) (*PrivateChatHandler, error) {
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

func HandlerChatRegister() {
	handlers[0] = &PrivateChatHandler{
		Id: messageId.MessageId_SCSendChatMsg,
		Fn: ResolvePrivateChatMsg,
	}
}

func ResolvePrivateChatMsg(p *PrivateChat, packet *network.Message) {
	req := &player.CSSendChatMsg{}
	err := proto.Unmarshal(packet.Data, req)
	if err != nil {
		return
	}
	fmt.Println(req.Msg.Content)
	p.SendMsg(messageId.MessageId_SCSendChatMsg, &player.SCSendChatMsg{})
}
