package friend

import (
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"github.com/phuhao00/greatestworks-proto/gen/player"
	"github.com/phuhao00/network"
	"github.com/phuhao00/sugar"
	"google.golang.org/protobuf/proto"
	"sync"
)

type Handler func(s *System, packet *network.Message)

var (
	handlers map[messageId.MessageId]Handler
	onceInit sync.Once
)

func GetHandler(id messageId.MessageId) Handler {
	return handlers[id]
}

func init() {
	onceInit.Do(func() {
		HandlerFriendRegister()
	})
}

func HandlerFriendRegister() {
	handlers[messageId.MessageId_CSAddFriend] = AddFriend
	handlers[messageId.MessageId_CSDelFriend] = DelFriend
}

func GetFriendList(p Owner, s *System, packet *network.Message) {

}

func GetFriendInfo(p Owner, s *System, packet *network.Message) {

}

func AddFriend(s *System, packet *network.Message) {
	req := &player.CSAddFriend{}

	err := proto.Unmarshal(packet.Data, req)
	if err != nil {
		return
	}

	if !sugar.CheckInSlice(req.UId, s.FriendList) {
		s.FriendList = append(s.FriendList, req.UId)
	}
	s.Owner.SendMsg(messageId.MessageId_SCAddFriend, &player.SCSendChatMsg{})

}

func DelFriend(s *System, packet *network.Message) {
	req := &player.CSDelFriend{}
	err := proto.Unmarshal(packet.Data, req)
	if err != nil {
		return
	}
	s.FriendList = sugar.DelOneInSlice(req.UId, s.FriendList)

	s.Owner.SendMsg(messageId.MessageId_SCDelFriend, &player.SCDelFriend{})
}

func GiveFriendItem(p Owner, s *System, packet *network.Message) {

}

func AddApply(p Owner, s *System, packet *network.Message) {

}

func ManagerApply(p Owner, s *System, packet *network.Message) {

}
