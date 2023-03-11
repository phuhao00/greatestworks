package player

import (
	"greatestworks/business/module"
)

// Manager 维护在线玩家
type Manager struct {
	*module.MetricsBase
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

func NewPlayerMgr() *Manager {
	return &Manager{
		players: make(map[uint64]*Player),
		addPCh:  make(chan *Player, 1),
	}
}

// Add ...
func (pm *Manager) Add(p *Player) {
	if pm.players[p.UId] != nil {
		return
	}
	pm.players[p.UId] = p
	go p.Start()
}

// Del ...
func (pm *Manager) Del(p Player) {
	delete(pm.players, p.UId)
}

func (pm *Manager) Run() {
	for {
		select {
		case p := <-pm.addPCh:
			pm.Add(p)
		}
	}
}

func (pm *Manager) GetPlayer(uId uint64) *Player {
	p, ok := pm.players[uId]
	if ok {
		return p
	}
	return nil
}
