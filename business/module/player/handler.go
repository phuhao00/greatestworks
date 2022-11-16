package player

import (
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"github.com/phuhao00/network"
	"google.golang.org/protobuf/proto"
)

func (p *Player) SendMsg(ID messageId.MessageId, message proto.Message) {
	id := uint64(ID)
	p.Session.AsyncSend(uint16(id), message)
}

func (p *Player) ResolveFriend(packet *network.Message) {
	//friend.GetHandler()
}
