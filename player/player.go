package player

import "greatestworks/define"

type Player struct {
	UId            uint64
	FriendList     []uint64 //朋友
	HandlerParamCh chan define.HandlerParam
	handlers       map[string]Handler
}

func NewPlayer() *Player {
	p := &Player{
		UId:        0,
		FriendList: make([]uint64, 100),
		handlers:   make(map[string]Handler),
	}
	p.HandlerRegister()
	return p
}

func (p *Player) Run() {
	for {
		select {
		case handlerParam := <-p.HandlerParamCh:
			if fn, ok := p.handlers[handlerParam.HandlerKey]; ok {
				fn(handlerParam.Data)
			}
		}
	}
}
