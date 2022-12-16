package minigame

import "sync"

type MiniGame struct {
	games sync.Map
}

type GameKey struct {
	OpenTime       int64  `json:"open_time"`
	CreatePlayerId uint64 `json:"create_player_id"`
	Category       uint16 `json:"category"`
}

func (g *MiniGame) Load(key GameKey) Abstract {
	value, ok := g.games.Load(key)
	if ok {
		switch key.Category {
		case 1:
			saveDog := value.(*SaveDog)
			return saveDog
		}
	}
	return nil
}
