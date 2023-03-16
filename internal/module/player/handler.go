package player

import (
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"github.com/phuhao00/network"
	"google.golang.org/protobuf/proto"
	"greatestworks/business/module/task"
	"greatestworks/internal/module/bag"
	"greatestworks/internal/module/chat"
	"greatestworks/internal/module/friend"
	task2 "greatestworks/internal/module/task"
)

func (p *Player) SendMsg(ID messageId.MessageId, message proto.Message) {
	id := uint64(ID)
	p.Session.AsyncSend(uint16(id), message)
}

func (p *Player) Handler(id messageId.MessageId, msg *network.Message) {
	if handler, _ := friend.GetHandler(id); handler != nil {
		handler.Fn(p.friendSystem, msg)
	}
	if handler, _ := chat.GetHandler(id); handler != nil {
		handler.Fn(p.privateChat, msg)
	}

	if handler, _ := bag.GetHandler(id); handler != nil {
		handler.Fn(p, msg)
	}

	if task2.IsBelongToHere(id) {
		task.GetManager().ChIn <- &task2.PlayerActionParam{
			MessageId: id,
			Player:    p,
			Packet:    msg,
		}
	}
}
