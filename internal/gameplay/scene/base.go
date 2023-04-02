package scene

import (
	"google.golang.org/protobuf/proto"
	actor2 "greatestworks/internal/gameplay/scene/actor"
	"sync"
	"sync/atomic"
)

type Base struct {
	Id         uint64
	ConfId     uint32
	Category   uint32
	CreateTime uint64
	FinishTime uint64
	Players    sync.Map
	Status     atomic.Value
}

func NewBase() *Base {
	return &Base{
		Id:         0,
		ConfId:     0,
		Category:   0,
		CreateTime: 0,
		FinishTime: 0,
		Players:    sync.Map{},
		Status:     atomic.Value{},
	}
}

func (b *Base) NotifyAll(message proto.Message) {
	b.Players.Range(func(key, value any) bool {
		player := value.(*actor2.Player)
		player.SendMsg(message)
		return true
	})
}

func (b *Base) NotifyNearby(actor actor2.Actor, message proto.Message) {
	//TODO implement me
	panic("implement me")
}

func (b *Base) NotifyPlayer(playerId uint64, message proto.Message) {
	v, ok := b.Players.Load(playerId)
	if ok {
		player := v.(actor2.Player)
		player.SendMsg(message)
	}
}

func (b *Base) OnCreate() {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Run() {
	//TODO implement me
	panic("implement me")
}

func (b *Base) OnDestroy() {
	//TODO implement me
	panic("implement me")
}

func (b *Base) loop() {

}

func (b *Base) monitor() {

}
