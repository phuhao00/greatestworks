package main

import (
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"google.golang.org/protobuf/proto"
	"greatestworks/internal/module/player"
)

func (w *World) BroadcastMsg(ids []uint64, msgId messageId.MessageId, msg proto.Message) {
	for _, id := range ids {
		p := w.GetPlayers(id)
		if p != nil {
			p.SendMsg(msgId, msg)
		}
	}
}

func (w *World) GetPlayers(id uint64) *player.Player {
	return w.pm.GetPlayer(id)
}
