package player

import (
	"fmt"
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"github.com/phuhao00/greatestworks-proto/gen/player"
	"github.com/phuhao00/network"

	"github.com/phuhao00/sugar"

	"google.golang.org/protobuf/proto"
)

type Handler func(packet *network.Message)

func (p *Player) AddFriend(packet *network.Message) {
	req := &player.CSAddFriend{}

	err := proto.Unmarshal(packet.Data, req)
	if err != nil {
		return
	}

	if !sugar.CheckInSlice(req.UId, p.FriendList) {
		p.FriendList = append(p.FriendList, req.UId)
	}

	bytes, err := proto.Marshal(&player.SCSendChatMsg{})
	if err != nil {
		return
	}

	rsp := &network.Message{
		ID:   uint64(messageId.MessageId_SCAddFriend),
		Data: bytes,
	}

	p.Session.SendMsg(rsp)
}

func (p *Player) DelFriend(packet *network.Message) {
	req := &player.CSDelFriend{}
	err := proto.Unmarshal(packet.Data, req)
	if err != nil {
		return
	}
	p.FriendList = sugar.DelOneInSlice(req.UId, p.FriendList)

	bytes, err := proto.Marshal(&player.SCDelFriend{})
	if err != nil {
		return
	}

	rsp := &network.Message{
		ID:   uint64(messageId.MessageId_SCDelFriend),
		Data: bytes,
	}

	p.Session.SendMsg(rsp)
}

func (p *Player) ResolveChatMsg(packet *network.Message) {

	req := &player.CSSendChatMsg{}
	err := proto.Unmarshal(packet.Data, req)
	if err != nil {
		return
	}
	fmt.Println(req.Msg.Content)

	bytes, err := proto.Marshal(&player.SCSendChatMsg{})
	if err != nil {
		return
	}

	rsp := &network.Message{
		ID:   uint64(messageId.MessageId_SCSendChatMsg),
		Data: bytes,
	}

	p.Session.SendMsg(rsp)
}
