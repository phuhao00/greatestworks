package player

// GameStatus 玩家游戏状态
type GameStatus uint16

const (
	Online        GameStatus = iota + 1 //在线状态
	OffLine                             //离线状态
	GmManager                           //GM管理
	OffLineResume                       //离线唤起，不是玩家真正的在线状态
)
