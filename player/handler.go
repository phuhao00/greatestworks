package player

import (
	"fmt"
	"greatestworks/network"
	"greatestworks/network/protocol/gen/player"

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
}

func (p *Player) DelFriend(packet *network.Message) {
	req := &player.CSDelFriend{}
	err := proto.Unmarshal(packet.Data, req)
	if err != nil {
		return
	}
	p.FriendList = sugar.DelOneInSlice(req.UId, p.FriendList)
}

func (p *Player) ResolveChatMsg(packet *network.Message) {

	req := &player.CSSendChatMsg{}
	err := proto.Unmarshal(packet.Data, req)
	if err != nil {
		return
	}
	fmt.Println(req.Msg.Content)
	// todo 收到消息 然后转发给客户端（当你的好友给你发消息情况）
}
