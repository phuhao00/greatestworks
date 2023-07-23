package player

import (
	"greatestworks/internal/communicate/chat"
	email2 "greatestworks/internal/communicate/email"
	"greatestworks/internal/communicate/friend"
	bag2 "greatestworks/internal/gameplay/bag"
	building2 "greatestworks/internal/gameplay/building"
	pet2 "greatestworks/internal/gameplay/pet"
	plant2 "greatestworks/internal/gameplay/plant"
	task2 "greatestworks/internal/gameplay/task"
	shop2 "greatestworks/internal/purchase/shop"
	vip2 "greatestworks/internal/purchase/vip"
)

var (
	_ pet2.Player       = (*Player)(nil)
	_ shop2.IPlayer     = (*Player)(nil)
	_ task2.Player      = (*Player)(nil)
	_ bag2.IPlayer      = (*Player)(nil)
	_ plant2.Player     = (*Player)(nil)
	_ building2.IPlayer = (*Player)(nil)
	_ email2.IPlayer    = (*Player)(nil)
	_ vip2.Player       = (*Player)(nil)
)

type GamePlay struct {
	friendSystem   *friend.System
	privateChat    *chat.PrivateChat
	taskData       *task2.Data
	petSystem      *pet2.System
	shopData       *shop2.Data
	bagSystem      *bag2.System
	vip            *vip2.Vip
	buildingSystem *building2.System
	plantSystem    *plant2.System
	emailData      *email2.Data
}

func InitGamePlay() GamePlay {
	return GamePlay{
		friendSystem:   nil,
		privateChat:    nil,
		taskData:       nil,
		petSystem:      nil,
		shopData:       nil,
		bagSystem:      nil,
		vip:            nil,
		buildingSystem: nil,
	}
}

func (p *GamePlay) GetTaskData() *task2.Data {
	return p.taskData
}

func (p *GamePlay) GetPetSystem() *pet2.System {
	return p.petSystem
}

func (p *GamePlay) GetShopData() *shop2.Data {
	return p.shopData
}

func (p *GamePlay) GetBagSystem() *bag2.System {
	return p.bagSystem
}

func (p *GamePlay) GetVip() *vip2.Vip {
	return p.vip
}

func (p *GamePlay) GetBuildingSystem() *building2.System {
	return p.buildingSystem
}

func (p *GamePlay) GetPlantSystem() *plant2.System {
	return p.plantSystem
}

func (p *GamePlay) GetEmailData() *email2.Data {
	return p.emailData
}
