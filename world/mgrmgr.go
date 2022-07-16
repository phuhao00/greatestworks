package world

import (
	"greatestworks/manager"
	"greatestworks/network"
)

type MgrMgr struct {
	Pm              *manager.PlayerMgr
	Server          *network.Server
	Handlers        map[uint64]func(message *network.SessionPacket)
	chSessionPacket chan *network.SessionPacket
}

func NewMgrMgr() *MgrMgr {
	m := &MgrMgr{Pm: &manager.PlayerMgr{}}
	m.Server = network.NewServer(":8023")
	m.Server.OnSessionPacket = m.OnSessionPacket
	return m
}

var MM *MgrMgr

func (mm *MgrMgr) Run() {
	go mm.Server.Run()
	go mm.Pm.Run()
}

func (mm *MgrMgr) OnSessionPacket(packet *network.SessionPacket) {
	if handler, ok := mm.Handlers[packet.Msg.ID]; ok {
		handler(packet)
	}

}
