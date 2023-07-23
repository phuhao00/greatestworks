//go:build wireinject
// +build wireinject

package player

import (
	"greatestworks/internal/communicate/chat"
	"greatestworks/internal/communicate/friend"
	"greatestworks/internal/gameplay/bag"
	"greatestworks/internal/gameplay/pet"
	"greatestworks/internal/gameplay/task"
	"greatestworks/internal/purchase/shop"
	"greatestworks/internal/purchase/vip"
)
import "github.com/google/wire"

var MegaSet = wire.NewSet(friend.NewSystem, chat.NewPrivateChat, task.NewTaskData, pet.NewSystem, shop.NewData, bag.NewSystem, vip.NewVip,
	wire.Struct(new(GamePlay), "friendSystem", "privateChat", "taskData", "petSystem", "shopData", "bagSystem", "vipevent"))

func NewGamePlay() *GamePlay {
	wire.Build(MegaSet)
	return &GamePlay{}
}
