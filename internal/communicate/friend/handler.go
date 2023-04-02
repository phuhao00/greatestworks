package friend

import (
	"errors"
	"github.com/phuhao00/greatestworks-proto/messageId"
	"github.com/phuhao00/greatestworks-proto/player"
	"github.com/phuhao00/network"
	"github.com/phuhao00/sugar"
	"google.golang.org/protobuf/proto"
	"sync"
)

type Handler struct {
	Id messageId.MessageId
	Fn func(s *System, packet *network.Message)
}

var (
	handlers     []*Handler
	onceInit     sync.Once
	MinMessageId messageId.MessageId
	MaxMessageId messageId.MessageId //handle 的消息范围
)

func IsBelongToHere(id messageId.MessageId) bool {
	return id > MinMessageId && id < MaxMessageId
}

func GetHandler(id messageId.MessageId) (*Handler, error) {

	if id > MinMessageId && id < MaxMessageId {
		return nil, errors.New("not in")
	}
	for _, handler := range handlers {
		if handler.Id == id {
			return handler, nil
		}
	}
	return nil, errors.New("not exist")
}

func init() {
	onceInit.Do(func() {
		HandlerFriendRegister()
	})
}

func HandlerFriendRegister() {
	handlers[0] = &Handler{
		messageId.MessageId_CSAddFriend,
		AddFriend,
	}
	handlers[1] = &Handler{
		messageId.MessageId_CSDelFriend,
		DelFriend,
	}
}

func GetFriendList(s *System, packet *network.Message) {

}

func GetFriendInfo(s *System, packet *network.Message) {

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
	s.IPlayer.SendMsg(messageId.MessageId_SCAddFriend, &player.SCSendChatMsg{})

}

func DelFriend(s *System, packet *network.Message) {
	req := &player.CSDelFriend{}
	err := proto.Unmarshal(packet.Data, req)
	if err != nil {
		return
	}
	s.FriendList = sugar.DelOneInSlice(req.UId, s.FriendList)

	s.IPlayer.SendMsg(messageId.MessageId_SCDelFriend, &player.SCDelFriend{})
}

func GiveFriendItem(s *System, packet *network.Message) {

}

func AddApply(s *System, packet *network.Message) {

}

func ManagerApply(s *System, packet *network.Message) {

}
