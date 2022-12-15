//go:build wireinject
// +build wireinject

package player

import (
	"greatestworks/business/module/bag"
	"greatestworks/business/module/chat"
	"greatestworks/business/module/friend"
	"greatestworks/business/module/pet"
	"greatestworks/business/module/shop"
	"greatestworks/business/module/task"
	"greatestworks/business/module/vip"
)
import "github.com/google/wire"

var MegaSet = wire.NewSet(friend.NewSystem, chat.NewPrivateChat, task.NewTaskData, pet.NewSystem, shop.NewData, bag.NewSystem, vip.NewVip,
	wire.Struct(new(GamePlay), "friendSystem", "privateChat", "taskData", "petSystem", "shopData", "bagSystem", "vip"))

func NewGamePlay() *GamePlay {
	wire.Build(MegaSet)
	return &GamePlay{}
}
