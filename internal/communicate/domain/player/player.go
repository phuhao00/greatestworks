package player

import (
	"github.com/phuhao00/fuse"
	"github.com/phuhao00/greatestworks-proto/messageId"
	"github.com/phuhao00/greatestworks-proto/player"
	"github.com/phuhao00/greatestworks-proto/server_common"
	"github.com/phuhao00/network"
	"google.golang.org/protobuf/proto"
	"greatestworks/aop/logger"
	"greatestworks/internal/communicate/domain/chat"
	"greatestworks/internal/communicate/domain/friend"
	"greatestworks/internal/gameplay/bag"
	"greatestworks/internal/gameplay/task"
)

type Player struct {
	*GamePlay
	*BaseInfo
	HandlerParamCh chan *network.Message
	Session        *network.TcpSession
	isOffline      bool
	PlayerID       uint64
	chanPlayerMsg  chan *player.PlayerMsgData
	chanServerMsg  chan *server_common.ServerMsgData
	LogicRouter    *fuse.LogicRouter
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

func (p *Player) GetName() string {
	return p.Name
}

func (p *Player) HandleClientMsgPacket(msgData *player.PlayerMsgData) {
	if p.isOffline {
		return
	}
	select {
	case p.chanPlayerMsg <- msgData:
	default:
		logger.Error("[消息] 转发消息到player goroutine错误 PlayerID:%v", p.PlayerID)
	}
}

func (p *Player) HandleServerMsgPacket(serverMsgData *server_common.ServerMsgData) {
	if p.isOffline {
		return
	}
	select {
	case p.chanServerMsg <- serverMsgData:
	default:
		logger.Error("[handleServerMsgPacket] PlayerID:%v", p.PlayerID)
	}
}

func (p *Player) SendMsg(ID messageId.MessageId, message proto.Message) {
	id := uint64(ID)
	p.Session.AsyncSend(id, message)
}

func (p *Player) Handler(id messageId.MessageId, msg *network.Message) {
	if handler, _ := friend.GetHandler(id); handler != nil {
		handler.Fn(p.friendSystem, msg)
	}
	if handler, _ := chat.GetHandler(id); handler != nil {
		handler.Fn(p.privateChat, msg)
	}

	if handler, _ := bag.GetHandler(id); handler != nil {
		handler.Fn(p, msg)
	}

	if task.IsBelongToHere(id) {
		task.GetMod().ChIn <- &task.PlayerActionParam{
			MessageId: id,
			Player:    p,
			Packet:    msg,
		}
	}
}
