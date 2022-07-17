package player

import (
	"greatestworks/network"
	"greatestworks/network/protocol/gen/messageId"
)

type Player struct {
	UId            uint64
	FriendList     []uint64 //朋友
	HandlerParamCh chan *network.Message
	handlers       map[messageId.MessageId]Handler
	session        *network.Session
}

func NewPlayer() *Player {
	p := &Player{
		UId:        0,
		FriendList: make([]uint64, 100),
		handlers:   make(map[messageId.MessageId]Handler),
	}
	p.HandlerRegister()
	return p
}

func (p *Player) Run() {
	for {
		select {
		case handlerParam := <-p.HandlerParamCh:
			if fn, ok := p.handlers[messageId.MessageId(handlerParam.ID)]; ok {
				fn(handlerParam)
			}
		}
	}
}
