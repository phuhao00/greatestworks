//go:build wireinject
// +build wireinject

package player

import (
	"greatestworks/internal/module/bag"
	"greatestworks/internal/module/chat"
	"greatestworks/internal/module/friend"
	"greatestworks/internal/module/pet"
	"greatestworks/internal/module/shop"
	"greatestworks/internal/module/task"
	"greatestworks/internal/module/vip"
)
import "github.com/google/wire"

var MegaSet = wire.NewSet(friend.NewSystem, chat.NewPrivateChat, task.NewTaskData, pet.NewSystem, shop.NewData, bag.NewSystem, vip.NewVip,
	wire.Struct(new(GamePlay), "friendSystem", "privateChat", "taskData", "petSystem", "shopData", "bagSystem", "vipevent"))

func NewGamePlay() *GamePlay {
	wire.Build(MegaSet)
	return &GamePlay{}
}
