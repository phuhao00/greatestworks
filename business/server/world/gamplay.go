package world

import (
	"greatestworks/business/module/activity"
	"greatestworks/business/module/bag"
	"greatestworks/business/module/chat"
	"greatestworks/business/module/email"
	"greatestworks/business/module/friend"
	"greatestworks/business/module/minigame"
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
