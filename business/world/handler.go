package world

import (
	"fmt"
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"github.com/phuhao00/greatestworks-proto/gen/player"
	logicPlayer "greatestworks/business/player"

	"github.com/phuhao00/network"
	"google.golang.org/protobuf/proto"
)

func (mm *MgrMgr) CreatePlayer(message *network.Packet) {
	msg := &player.CSCreateUser{}
	err := proto.Unmarshal(message.Msg.Data, msg)
	if err != nil {
		return
	}
	fmt.Println("[MgrMgr.CreatePlayer]", msg)
	mm.SendMsg(uint64(messageId.MessageId_SCCreatePlayer), &player.SCCreateUser{}, message.Conn)

}

func (mm *MgrMgr) UserLogin(message *network.Packet) {
	msg := &player.CSLogin{}
	err := proto.Unmarshal(message.Msg.Data, msg)
	if err != nil {
		return
	}
	newPlayer := logicPlayer.NewPlayer()
	newPlayer.UId = 111
	newPlayer.Session = message.Conn
	mm.Pm.Add(newPlayer)

}

func (mm *MgrMgr) SendMsg(id uint64, message proto.Message, session *network.TcpConnX) {
	session.AsyncSend(uint16(id), message)
}
