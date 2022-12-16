package minigame

import "sync"

type MiniGame struct {
	configs sync.Map
	games   sync.Map
}

type GameKey struct {
	OpenTime       int64  `json:"open_time"`
	CreatePlayerId uint64 `json:"create_player_id"`
	Category       uint16 `json:"category"`
}

// GetGame 获取游戏
func (g *MiniGame) GetGame(key GameKey) Abstract {
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

// AddGame 添加游戏
func (g *MiniGame) AddGame(key GameKey, game Abstract) {
	g.games.Store(key, game)
}

// DeleteGame 删除游戏
func (g *MiniGame) DeleteGame(key GameKey) {
	g.games.Delete(key)
}

// Load 加载
func (g *MiniGame) Load() {

}

// Save 保存
func (g *MiniGame) Save() {

}
