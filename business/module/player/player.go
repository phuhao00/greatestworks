package player

import (
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"greatestworks/business/module/chat"
	"greatestworks/business/module/friend"

	"github.com/phuhao00/network"
)

type Player struct {
	UId            uint64
	FriendList     []uint64 //朋友
	HandlerParamCh chan *network.Message
	handlers       map[messageId.MessageId]Handler
	Session        *network.TcpConnX
	FriendSystem   friend.System
	PrivateChat    chat.PrivateChat
}

func NewPlayer() *Player {
	p := &Player{
		UId:        0,
		FriendList: make([]uint64, 100),
		handlers:   make(map[messageId.MessageId]Handler),
	}
	p.HandlerRegister()
	p.FriendSystem.SetOwner(p)
	p.PrivateChat.SetHandler(p)
	return p
}

func (p *Player) Start() {
	for {
		select {
		case handlerParam := <-p.HandlerParamCh:
			if fn, ok := p.handlers[messageId.MessageId(handlerParam.ID)]; ok {
				fn(handlerParam)
			}
		}
	}
}

func (p *Player) Stop() {

}

func (p *Player) OnLogin() {
	//从db加载数据初始化
	//同步数据给客户端

}

func (p *Player) OnLogout() {
	//存db
}
