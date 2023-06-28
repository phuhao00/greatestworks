package scene

import (
	"google.golang.org/protobuf/proto"
)

type IScene interface {
	OnCreate()
	Run()
	OnDestroy()
	loop()
	monitor()
}

type Notify interface {
	NotifyAll(message proto.Message)
	NotifyNearby(actor IActor, message proto.Message)
	NotifyPlayer(playerId uint64, message proto.Message)
}

type Action interface {
	OnNextWave()
	OnMonsterDie()
	OnWaveEnd()
}

type FightScene interface {
	IScene
	Action
}
