package world

import (
	"greatestworks/manager"
	"greatestworks/network"
)

type MgrMgr struct {
	Pm     *manager.PlayerMgr
	Server *network.Server
}

func NewMgrMgr() *MgrMgr {
	m := &MgrMgr{Pm: &manager.PlayerMgr{}}
	return m
}

var MM *MgrMgr

func (mm *MgrMgr) name() {

}
