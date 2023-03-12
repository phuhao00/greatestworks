package player

import (
	"greatestworks/aop/event"
	"greatestworks/business/module"
)

var (
	Mod *Module
)

func init() {
	module.MManager.RegisterModule("", Mod)
}

type Module struct {
	*module.MetricsBase
	players map[uint64]*Player
	addPCh  chan *Player
}

func (m *Module) OnEvent(c module.Character, event event.IEvent) {
	//TODO implement me
	panic("implement me")
}

func (m *Module) SetEventCategoryActive(eventCategory int) {
	//TODO implement me
	panic("implement me")
}

func (pm *Module) OnStart() {
	//TODO implement me
	panic("implement me")
}

func (pm *Module) AfterStart() {
	//TODO implement me
	panic("implement me")
}

func (pm *Module) OnStop() {
	//TODO implement me
	panic("implement me")
}

func (pm *Module) AfterStop() {
	//TODO implement me
	panic("implement me")
}

func NewPlayerMgr() *Module {
	return &Module{
		players: make(map[uint64]*Player),
		addPCh:  make(chan *Player, 1),
	}
}

// Add ...
func (pm *Module) Add(p *Player) {
	if pm.players[p.UId] != nil {
		return
	}
	pm.players[p.UId] = p
	go p.Start()
}

// Del ...
func (pm *Module) Del(p Player) {
	delete(pm.players, p.UId)
}

func (pm *Module) Run() {
	for {
		select {
		case p := <-pm.addPCh:
			pm.Add(p)
		}
	}
}

func (pm *Module) GetPlayer(uId uint64) *Player {
	p, ok := pm.players[uId]
	if ok {
		return p
	}
	return nil
}
