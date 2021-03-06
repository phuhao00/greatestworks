package world

import (
	"greatestworks/manager"
	"greatestworks/network"
	"greatestworks/network/protocol/gen/messageId"
)

type MgrMgr struct {
	Pm              *manager.PlayerMgr
	Server          *network.Server
	Handlers        map[messageId.MessageId]func(message *network.SessionPacket)
	chSessionPacket chan *network.SessionPacket
}

func NewMgrMgr() *MgrMgr {
	m := &MgrMgr{Pm: manager.NewPlayerMgr()}
	m.Server = network.NewServer(":8023")
	m.Server.OnSessionPacket = m.OnSessionPacket
	m.Handlers = make(map[messageId.MessageId]func(message *network.SessionPacket))

	return m
}

var MM *MgrMgr

func (mm *MgrMgr) Run() {
	mm.HandlerRegister()
	go mm.Server.Run()
	go mm.Pm.Run()
}

func (mm *MgrMgr) OnSessionPacket(packet *network.SessionPacket) {
	if handler, ok := mm.Handlers[messageId.MessageId(packet.Msg.ID)]; ok {
		handler(packet)
		return
	}
	if p := mm.Pm.GetPlayer(packet.Sess.UId); p != nil {
		p.HandlerParamCh <- packet.Msg
	}
}
