package chat

import (
	"fmt"
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"github.com/phuhao00/greatestworks-proto/gen/player"
	"github.com/phuhao00/network"
	"google.golang.org/protobuf/proto"
	"sync"
)

var (
	handlers map[messageId.MessageId]func(p *PrivateChat, packet *network.Message)
	onceInit sync.Once
)

func init() {
	onceInit.Do(func() {
		HandlerChatRegister()
	})
}

func GetHandler(id messageId.MessageId) func(p *PrivateChat, packet *network.Message) {
	return handlers[id]
}

func HandlerChatRegister() {
	handlers[messageId.MessageId_CSSendChatMsg] = ResolvePrivateChatMsg
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
