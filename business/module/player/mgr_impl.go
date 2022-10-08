package player

import "greatestworks/business/module/base"

//Manager 维护在线玩家
type Manager struct {
	*base.MetricsBase
	players map[uint64]*Player
	addPCh  chan *Player
}

func (pm *Manager) OnStart() {
	//TODO implement me
	panic("implement me")
}

func (pm *Manager) AfterStart() {
	//TODO implement me
	panic("implement me")
}

func (pm *Manager) OnStop() {
	//TODO implement me
	panic("implement me")
}

func (pm *Manager) AfterStop() {
	//TODO implement me
	panic("implement me")
}
