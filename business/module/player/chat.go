package player

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"github.com/phuhao00/greatestworks-proto/gen/player"
	"github.com/phuhao00/network"
	"google.golang.org/protobuf/proto"
)

func HandlerChatRegister() {
	handlers[messageId.MessageId_CSSendChatMsg] = ResolveChatMsg
}

func ResolveChatMsg(p *Player, packet *network.Message) {
	req := &player.CSSendChatMsg{}
	err := proto.Unmarshal(packet.Data, req)
	if err != nil {
		return
	}
	fmt.Println(req.Msg.Content)
	p.SendMsg(messageId.MessageId_SCSendChatMsg, &player.SCSendChatMsg{})
}

func (p *Player) InitNsqHandler(channel string) {
	//TODO implement me
	panic("implement me")
}

func (p *Player) HandleMessage(message nsq.Message) error {
	//TODO implement me
	panic("implement me")
}

func (p *Player) PublishChatMsg(chatMsg interface{}) error {
	//TODO implement me
	panic("implement me")
}
