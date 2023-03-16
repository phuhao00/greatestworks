package world

import (
	"fmt"
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"github.com/phuhao00/greatestworks-proto/gen/player"
	"github.com/phuhao00/network"
	"google.golang.org/protobuf/proto"
	logicPlayer "greatestworks/internal/module/player"
)

func (w *World) CreatePlayer(message *network.Packet) {
	msg := &player.CSCreateUser{}
	err := proto.Unmarshal(message.Msg.Data, msg)
	if err != nil {
		return
	}
	fmt.Println("[World.CreatePlayer]", msg)
	w.SendMsg(uint64(messageId.MessageId_SCCreatePlayer), &player.SCCreateUser{}, message.Conn)

}

func (w *World) UserLogin(message *network.Packet) {
	msg := &player.CSLogin{}
	err := proto.Unmarshal(message.Msg.Data, msg)
	if err != nil {
		return
	}
	newPlayer := logicPlayer.NewPlayer()
	newPlayer.UId = 111
	newPlayer.Session = message.Conn
	w.pm.Add(newPlayer)

}

func (w *World) SendMsg(id uint64, message proto.Message, session *network.TcpConnX) {
	session.AsyncSend(uint16(id), message)
}
