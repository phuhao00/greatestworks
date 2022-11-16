package player

import (
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"github.com/phuhao00/network"
	"google.golang.org/protobuf/proto"
	"greatestworks/business/module/chat"
	"greatestworks/business/module/friend"
)

func (p *Player) SendMsg(ID messageId.MessageId, message proto.Message) {
	id := uint64(ID)
	p.Session.AsyncSend(uint16(id), message)
}

func (p *Player) Handler(id messageId.MessageId, msg *network.Message) {
	if handler, _ := friend.GetHandler(id); handler != nil {
		handler.Fn(p.FriendSystem, msg)
	}
	if handler, _ := chat.GetHandler(id); handler != nil {
		handler.Fn(p.PrivateChat, msg)
	}

}
