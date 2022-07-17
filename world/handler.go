package world

import (
	"fmt"
	"greatestworks/network"
	"greatestworks/network/protocol/gen/player"
	logicPlayer "greatestworks/player"
	"time"

	"google.golang.org/protobuf/proto"
)

func (mm *MgrMgr) CreatePlayer(message *network.SessionPacket) {
	msg := &player.CSCreateUser{}
	err := proto.Unmarshal(message.Msg.Data, msg)
	if err != nil {
		return
	}
	fmt.Println("[MgrMgr.CreatePlayer]", msg)
	mm.SendMsg(message.Msg.ID, &player.SCCreateUser{}, message.Sess)

}

func (mm *MgrMgr) UserLogin(message *network.SessionPacket) {
	msg := &player.CSLogin{}
	err := proto.Unmarshal(message.Msg.Data, msg)
	if err != nil {
		return
	}
	newPlayer := logicPlayer.NewPlayer()
	newPlayer.UId = uint64(time.Now().Unix())
	newPlayer.HandlerParamCh = message.Sess.WriteCh
	message.Sess.IsPlayerOnline = true
	mm.Pm.Add(newPlayer)
	newPlayer.Run()
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
