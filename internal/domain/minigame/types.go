package minigame

import (
	"fmt"
	"time"
)

// GamePlayer 游戏玩家
type GamePlayer struct {
	PlayerID string    `json:"player_id"`
	Username string    `json:"username"`
	Score    int64     `json:"score"`
	JoinTime time.Time `json:"join_time"`
	IsActive bool      `json:"is_active"`
	IsWinner bool      `json:"is_winner"`
	Rank     int       `json:"rank"`
}

// GameData 游戏数据
type GameData struct {
	GameID      string                 `json:"game_id"`
	GameType    GameType               `json:"game_type"`
	Config      map[string]interface{} `json:"config"`
	StartTime   time.Time              `json:"start_time"`
	EndTime     time.Time              `json:"end_time"`
	Duration    time.Duration          `json:"duration"`
	MaxPlayers  int                    `json:"max_players"`
	MinPlayers  int                    `json:"min_players"`
	IsActive    bool                   `json:"is_active"`
	IsCompleted bool                   `json:"is_completed"`
}

// GameResult 游戏结果
type GameResult struct {
	GameID      string       `json:"game_id"`
	WinnerID    string       `json:"winner_id"`
	WinnerName  string       `json:"winner_name"`
	FinalScore  int64        `json:"final_score"`
	Rankings    []GamePlayer `json:"rankings"`
	Rewards     []Reward     `json:"rewards"`
	CompletedAt time.Time    `json:"completed_at"`
}

// RewardPool 奖励池
type RewardPool struct {
	PoolID      string    `json:"pool_id"`
	GameID      string    `json:"game_id"`
	TotalReward int64     `json:"total_reward"`
	Rewards     []Reward  `json:"rewards"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Reward 奖励
type Reward struct {
	RewardID   string `json:"reward_id"`
	PlayerID   string `json:"player_id"`
	RewardType string `json:"reward_type"`
	Amount     int64  `json:"amount"`
	ItemID     string `json:"item_id,omitempty"`
	ItemCount  int    `json:"item_count,omitempty"`
}

// GameStatistics 游戏统计
type GameStatistics struct {
	GameID          string        `json:"game_id"`
	TotalPlayers    int           `json:"total_players"`
	ActivePlayers   int           `json:"active_players"`
	CompletedGames  int           `json:"completed_games"`
	AverageScore    float64       `json:"average_score"`
	HighestScore    int64         `json:"highest_score"`
	AverageDuration time.Duration `json:"average_duration"`
	LastPlayedAt    time.Time     `json:"last_played_at"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

// NewGameStatistics 创建游戏统计信息
func NewGameStatistics() *GameStatistics {
	now := time.Now()
	return &GameStatistics{
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewGameData 创建游戏数据
func NewGameData() *GameData {
	now := time.Now()
	return &GameData{
		StartTime: now,
		Config:    make(map[string]interface{}),
	}
}

// generateScoreID 生成分数ID
func generateScoreID() string {
	return fmt.Sprintf("score_%d", time.Now().UnixNano())
}

// SetData 设置游戏数据
func (gd *GameData) SetData(key string, value interface{}) {
	if gd.Config == nil {
		gd.Config = make(map[string]interface{})
	}
	gd.Config[key] = value
}

// GetData 获取游戏数据
func (gd *GameData) GetData(key string) (interface{}, bool) {
	if gd.Config == nil {
		return nil, false
	}
	value, exists := gd.Config[key]
	return value, exists
}

// Clone 克隆游戏数据
func (gd *GameData) Clone() *GameData {
	clone := &GameData{
		GameID:     gd.GameID,
		GameType:   gd.GameType,
		StartTime:  gd.StartTime,
		EndTime:    gd.EndTime,
		Duration:   gd.Duration,
		MaxPlayers: gd.MaxPlayers,
		MinPlayers: gd.MinPlayers,
		Config:     make(map[string]interface{}),
	}

	// 复制配置
	for k, v := range gd.Config {
		clone.Config[k] = v
	}

	return clone
}

// Clone 克隆游戏统计信息
func (gs *GameStatistics) Clone() *GameStatistics {
	return &GameStatistics{
		GameID:          gs.GameID,
		TotalPlayers:    gs.TotalPlayers,
		ActivePlayers:   gs.ActivePlayers,
		CompletedGames:  gs.CompletedGames,
		AverageScore:    gs.AverageScore,
		HighestScore:    gs.HighestScore,
		AverageDuration: gs.AverageDuration,
		LastPlayedAt:    gs.LastPlayedAt,
		CreatedAt:       gs.CreatedAt,
		UpdatedAt:       gs.UpdatedAt,
	}
}
