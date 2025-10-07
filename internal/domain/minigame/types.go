package minigame

import (
	"fmt"
	"time"
)

// GameCategory 游戏分类
type GameCategory string

const (
	GameCategoryNormal      GameCategory = "normal"
	GameCategoryCompetitive GameCategory = "competitive"
	GameCategoryCasual      GameCategory = "casual"
	GameCategoryRanked      GameCategory = "ranked"
)

// GameResult 游戏结果
type GameResult struct {
	GameID      string    `json:"game_id"`
	WinnerID    uint64    `json:"winner_id"`
	WinnerName  string    `json:"winner_name"`
	FinalScore  int64     `json:"final_score"`
	CompletedAt time.Time `json:"completed_at"`
	Rank        int       `json:"rank"`
	Score       int64     `json:"score"`
	IsWinner    bool      `json:"is_winner"`
	PlayerID    uint64    `json:"player_id"`
}

// GamePlayer 游戏玩家
type GamePlayer struct {
	PlayerID uint64    `json:"player_id"`
	Username string    `json:"username"`
	Score    int64     `json:"score"`
	JoinTime time.Time `json:"join_time"`
	IsActive bool      `json:"is_active"`
	IsWinner bool      `json:"is_winner"`
	Rank     int       `json:"rank"`
}

// Clone 克隆游戏玩家
func (gp *GamePlayer) Clone() *GamePlayer {
	return &GamePlayer{
		PlayerID: gp.PlayerID,
		Username: gp.Username,
		Score:    gp.Score,
		JoinTime: gp.JoinTime,
		IsActive: gp.IsActive,
		IsWinner: gp.IsWinner,
		Rank:     gp.Rank,
	}
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
	RewardID   string    `json:"reward_id"`
	PlayerID   string    `json:"player_id"`
	RewardType string    `json:"reward_type"`
	Amount     int64     `json:"amount"`
	ItemID     string    `json:"item_id,omitempty"`
	ItemCount  int       `json:"item_count,omitempty"`
	GameID     string    `json:"game_id"`
	Timestamp  time.Time `json:"timestamp"`
}

// GameStatistics 游戏统计
type GameStatistics struct {
	GameID              string        `json:"game_id"`
	TotalPlayers        int           `json:"total_players"`
	ActivePlayers       int           `json:"active_players"`
	CompletedGames      int           `json:"completed_games"`
	AverageScore        float64       `json:"average_score"`
	HighestScore        int64         `json:"highest_score"`
	LowestScore         int64         `json:"lowest_score"`
	AverageDuration     time.Duration `json:"average_duration"`
	AverageGameDuration float64       `json:"average_game_duration"`
	TotalGames          int           `json:"total_games"`
	LastPlayedAt        time.Time     `json:"last_played_at"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
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

// GetTotalPlays 获取总游戏次数
func (gs *GameStatistics) GetTotalPlays() int {
	return gs.TotalGames
}

// GetTotalPlayers 获取总玩家数
func (gs *GameStatistics) GetTotalPlayers() int {
	return gs.TotalPlayers
}

// GetAverageScore 获取平均分数
func (gs *GameStatistics) GetAverageScore() float64 {
	return gs.AverageScore
}

// GetHighestScore 获取最高分数
func (gs *GameStatistics) GetHighestScore() int64 {
	return gs.HighestScore
}

// GetAverageTime 获取平均游戏时长
func (gs *GameStatistics) GetAverageTime() time.Duration {
	return gs.AverageDuration
}

// GetCompletionRate 获取完成率
func (gs *GameStatistics) GetCompletionRate() float64 {
	if gs.TotalGames == 0 {
		return 0.0
	}
	return float64(gs.CompletedGames) / float64(gs.TotalGames)
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

// CalculateRewards 计算奖励
func (rp *RewardPool) CalculateRewards(rank int, score int64, isWinner bool) []Reward {
	// 简单的奖励计算逻辑
	rewards := make([]Reward, 0)

	if isWinner {
		// 获胜者获得基础奖励
		rewards = append(rewards, Reward{
			RewardType: "experience",
			Amount:     int64(100 * rank),
		})
	}

	// 根据分数给予额外奖励
	if score > 1000 {
		rewards = append(rewards, Reward{
			RewardType: "coin",
			Amount:     score / 10,
		})
	}

	return rewards
}

// Clone 克隆奖励池
func (rp *RewardPool) Clone() *RewardPool {
	clone := &RewardPool{
		PoolID:      rp.PoolID,
		GameID:      rp.GameID,
		TotalReward: rp.TotalReward,
		CreatedAt:   rp.CreatedAt,
		UpdatedAt:   rp.UpdatedAt,
		Rewards:     make([]Reward, len(rp.Rewards)),
	}

	// 复制奖励列表
	for i, reward := range rp.Rewards {
		clone.Rewards[i] = reward
	}

	return clone
}

// 注意：GameReward已经在entity.go中定义，不需要重复定义
// 注意：GameStatistics和GameResult已经在文件开头定义，这里只是添加了缺失的字段
