package server

import (
	"fmt"
	"github.com/phuhao00/greatestworks-proto/messageId"
	"github.com/phuhao00/greatestworks-proto/player"
	"github.com/phuhao00/network"
	"google.golang.org/protobuf/proto"
	logicPlayer "greatestworks/internal/communicate/player"
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
	w.playerManager.Add(newPlayer)

}

func (w *World) SendMsg(id uint64, message proto.Message, session *network.TcpSession) {
	session.AsyncSend(id, message)
}
