package scene

import (
	"google.golang.org/protobuf/proto"
	"greatestworks/internal/module/scene/actor"
)

type Abstract interface {
	OnCreate()
	Run()
	OnDestroy()
	loop()
	monitor()
}

type Notify interface {
	NotifyAll(message proto.Message)
	NotifyNearby(actor actor.Actor, message proto.Message)
	NotifyPlayer(playerId uint64, message proto.Message)
}

type Action interface {
	OnNextWave()
	OnMonsterDie()
	OnWaveEnd()
}

type FightScene interface {
	Abstract
	Action
}
