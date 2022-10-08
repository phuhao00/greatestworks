package world

import (
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"greatestworks/aop/logger"
	"greatestworks/business/module/player"
	"os"
	"syscall"

	"github.com/phuhao00/network"
)

type MgrMgr struct {
	Pm              *player.Manager
	Server          *network.Server
	Handlers        map[messageId.MessageId]func(message *network.Packet)
	chSessionPacket chan *network.Packet
}

func NewMgrMgr() *MgrMgr {
	m := &MgrMgr{Pm: player.NewPlayerMgr()}
	m.Server = network.NewServer(":8023", 100, 200, logger.Logger)
	m.Server.MessageHandler = m.OnSessionPacket
	m.Handlers = make(map[messageId.MessageId]func(message *network.Packet))

	return m
}

var MM *MgrMgr

func (mm *MgrMgr) Start() {
	mm.HandlerRegister()
	go mm.Server.Run()
	go mm.Pm.Run()
}

func (mm *MgrMgr) Stop() {

}

func (mm *MgrMgr) OnSessionPacket(packet *network.Packet) {
	if handler, ok := mm.Handlers[messageId.MessageId(packet.Msg.ID)]; ok {
		handler(packet)
		return
	}
	if p := mm.Pm.GetPlayer(uint64(packet.Conn.ConnID)); p != nil {
		p.HandlerParamCh <- packet.Msg
	}
}

func (mm *MgrMgr) OnSystemSignal(signal os.Signal) bool {
	logger.Logger.DebugF("[MgrMgr] 收到信号 %v \n", signal)
	tag := true
	switch signal {
	case syscall.SIGHUP:
		//todo
	case syscall.SIGPIPE:
	default:
		logger.Logger.DebugF("[MgrMgr] 收到信号准备退出...")
		tag = false

	}
	return tag
}
