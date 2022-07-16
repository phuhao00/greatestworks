package player

import (
	"greatestworks/network"
)

type Player struct {
	UId            uint64
	FriendList     []uint64 //朋友
	HandlerParamCh chan *network.Message
	handlers       map[uint64]Handler
	session        *network.Session
}

func NewPlayer() *Player {
	p := &Player{
		UId:        0,
		FriendList: make([]uint64, 100),
		handlers:   make(map[uint64]Handler),
	}
	p.HandlerRegister()
	return p
}

func (p *Player) Run() {
	for {
		select {
		case handlerParam := <-p.HandlerParamCh:
			if fn, ok := p.handlers[handlerParam.ID]; ok {
				fn(handlerParam.Data)
			}
		}
	}
}
