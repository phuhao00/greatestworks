package world

import (
	"fmt"
	logicPlayer "greatestworks/business/player"
	"greatestworks/network/protocol/gen/messageId"
	"greatestworks/network/protocol/gen/player"

	"github.com/phuhao00/network"
	"google.golang.org/protobuf/proto"
)

func (mm *MgrMgr) CreatePlayer(message *network.SessionPacket) {
	msg := &player.CSCreateUser{}
	err := proto.Unmarshal(message.Msg.Data, msg)
	if err != nil {
		return
	}
	fmt.Println("[MgrMgr.CreatePlayer]", msg)
	mm.SendMsg(uint64(messageId.MessageId_SCCreatePlayer), &player.SCCreateUser{}, message.Sess)

}

func (mm *MgrMgr) UserLogin(message *network.SessionPacket) {
	msg := &player.CSLogin{}
	err := proto.Unmarshal(message.Msg.Data, msg)
	if err != nil {
		return
	}
	newPlayer := logicPlayer.NewPlayer()
	newPlayer.UId = 111
	//newPlayer.UId = uint64(time.Now().Unix())
	newPlayer.HandlerParamCh = message.Sess.WriteCh
	message.Sess.IsPlayerOnline = true
	message.Sess.UId = newPlayer.UId
	newPlayer.Session = message.Sess
	mm.Pm.Add(newPlayer)

}

func (mm *MgrMgr) SendMsg(id uint64, message proto.Message, session *network.Session) {
	bytes, err := proto.Marshal(message)
	if err != nil {
		return
	}
	rsp := &network.Message{
		ID:   id,
		Data: bytes,
	}
	session.SendMsg(rsp)
}
