package player

import (
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"github.com/phuhao00/greatestworks-proto/gen/player"
	"github.com/phuhao00/network"
	"github.com/phuhao00/sugar"
	"google.golang.org/protobuf/proto"
)

func (p *Player) GetFriendList() {
	//TODO implement me
	panic("implement me")
}

func (p *Player) GetFriendInfo() {
	//TODO implement me
	panic("implement me")
}

func (p *Player) AddFriend(packet *network.Message) {
	req := &player.CSAddFriend{}

	err := proto.Unmarshal(packet.Data, req)
	if err != nil {
		return
	}

	if !sugar.CheckInSlice(req.UId, p.FriendList) {
		p.FriendList = append(p.FriendList, req.UId)
	}

	p.SendMsg(messageId.MessageId_SCAddFriend, &player.SCSendChatMsg{})

}

func (p *Player) DelFriend(packet *network.Message) {
	req := &player.CSDelFriend{}
	err := proto.Unmarshal(packet.Data, req)
	if err != nil {
		return
	}
	p.FriendList = sugar.DelOneInSlice(req.UId, p.FriendList)
	p.SendMsg(messageId.MessageId_SCDelFriend, &player.SCDelFriend{})
}

func (p *Player) GiveFriendItem() {
	//TODO implement me
	panic("implement me")
}

func (p *Player) FriendAddApply() {
	//TODO implement me
	panic("implement me")
}

func (p *Player) ManagerFriendApply() {
	//TODO implement me
	panic("implement me")
}
