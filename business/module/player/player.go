package player

import (
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"github.com/phuhao00/network"
	event2 "greatestworks/aop/event"
)

type Player struct {
	*GamePlay
	*BaseInfo
	HandlerParamCh chan *network.Message
	Session        *network.TcpConnX
}

func NewPlayer() *Player {
	p := &Player{
		GamePlay: NewGamePlay(),
	}
	return p
}

func (p *Player) Start() {
	for {
		select {
		case handlerParam := <-p.HandlerParamCh:
			p.Handler(messageId.MessageId(handlerParam.ID), handlerParam)
		}
	}
}

func (p *Player) Stop() {

}

func (p *Player) OnLogin() {
	//从db加载数据初始化
	//同步数据给客户端
	p.taskData.LoadFromDB()

}

func (p *Player) OnLogout() {
	//存db
}

func (p *Player) OnEvent(iEvent event2.IEvent) {
	//todo get module process event
	//eg:
	p.taskData.OnEvent(iEvent)
}
