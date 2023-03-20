package world

import (
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"greatestworks/aop/logger"
	"greatestworks/internal/module/chat"
	"greatestworks/internal/module/family"
	"greatestworks/internal/module/player"
	"os"
	"syscall"

	"github.com/phuhao00/network"
)

type World struct {
	Server          *network.Server
	Handlers        map[messageId.MessageId]func(message *network.Packet)
	chSessionPacket chan *network.Packet
	chatSystem      *chat.System
	familyManager   *family.Module
	pm              *player.Module
}

func NewWorld() *World {
	m := &World{pm: player.NewPlayerMgr()}
	m.Server = network.NewServer(":8023", 100, 200, logger.Logger)
	m.Server.MessageHandler = m.OnSessionPacket
	m.Handlers = make(map[messageId.MessageId]func(message *network.Packet))
	return m
}

var Oasis *World

func (w *World) OnSessionPacket(packet *network.Packet) {
	if handler, ok := w.Handlers[messageId.MessageId(packet.Msg.ID)]; ok {
		handler(packet)
		return
	}
	if p := w.pm.GetPlayer(uint64(packet.Conn.ConnID)); p != nil {
		p.HandlerParamCh <- packet.Msg
	}
}

func (w *World) OnSystemSignal(signal os.Signal) bool {
	logger.Logger.DebugF("[World] 收到信号 %v \n", signal)
	tag := true
	switch signal {
	case syscall.SIGHUP:
		//todo
	case syscall.SIGPIPE:
	default:
		logger.Logger.DebugF("[World] 收到信号准备退出...")
		tag = false

	}
	return tag
}
