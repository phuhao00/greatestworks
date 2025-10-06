package minigame

import "time"

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
