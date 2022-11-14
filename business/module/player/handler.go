package player

import (
	"fmt"
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"github.com/phuhao00/greatestworks-proto/gen/player"
	"github.com/phuhao00/network"

	"google.golang.org/protobuf/proto"
)

type Handler func(packet *network.Message)

func (p *Player) ResolveChatMsg(packet *network.Message) {

	req := &player.CSSendChatMsg{}
	err := proto.Unmarshal(packet.Data, req)
	if err != nil {
		return
	}
	fmt.Println(req.Msg.Content)
	p.SendMsg(messageId.MessageId_SCSendChatMsg, &player.SCSendChatMsg{})
}

func (p *Player) SendMsg(ID messageId.MessageId, message proto.Message) {
	id := uint64(ID)
	p.Session.AsyncSend(uint16(id), message)
}
