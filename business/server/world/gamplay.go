package world

import (
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"google.golang.org/protobuf/proto"
	"greatestworks/business/module/activity"
	"greatestworks/business/module/bag"
	"greatestworks/business/module/chat"
	"greatestworks/business/module/email"
	"greatestworks/business/module/friend"
	"greatestworks/business/module/minigame"
	"greatestworks/business/module/player"
	"greatestworks/business/module/rank"
	"greatestworks/business/module/recharge"
	"greatestworks/business/module/task"
)

type GamePlay struct {
	activity activity.Abstract
	bag      bag.Abstract
	chat     chat.Abstract
	rank     rank.Abstract
	email    email.Abstract
	friend   friend.Abstract
	minigame minigame.Abstract
	recharge recharge.Abstract
	task     task.Abstract
}
type Option func(play *GamePlay) *GamePlay

func WithActivity(activity activity.Abstract) Option {
	return func(play *GamePlay) *GamePlay {
		play.activity = activity
		return play
	}
}

func NewGamePlay(option ...Option) *GamePlay {
	g := &GamePlay{}
	for _, op := range option {
		op(g)
	}
	return g
}

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
