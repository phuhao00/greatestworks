package player

import (
	"greatestworks/internal"
	"sync"
)

const (
	ModuleName = "player"
)

var (
	Mod         *Module
	onceInitMod sync.Once
)

func init() {
	internal.MManager.RegisterModule(ModuleName, GetMod())
}

type Module struct {
	*internal.BaseModule
	*internal.MetricsBase
	players map[uint64]*Player
	addPCh  chan *Player
}

func GetMod() *Module {
	onceInitMod.Do(func() {
		Mod = &Module{BaseModule: internal.NewBaseModule()}
	})

	return Mod
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
