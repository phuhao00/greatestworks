package minigame

import (
	"fmt"
	"time"
)

// RewardType 奖励类型
type RewardType string

const (
	RewardTypeCoin     RewardType = "coin"
	RewardTypeExp      RewardType = "exp"
	RewardTypeItem     RewardType = "item"
	RewardTypeCurrency RewardType = "currency"
)

// GameSession 游戏会话实体
type GameSession struct {
	ID           string                 `json:"id" bson:"_id"`
	GameID       string                 `json:"game_id" bson:"game_id"`
	PlayerID     uint64                 `json:"player_id" bson:"player_id"`
	SessionToken string                 `json:"session_token" bson:"session_token"`
	Status       PlayerStatus           `json:"status" bson:"status"`
	JoinedAt     time.Time              `json:"joined_at" bson:"joined_at"`
	LeftAt       *time.Time             `json:"left_at,omitempty" bson:"left_at,omitempty"`
	LeaveReason  *PlayerLeaveReason     `json:"leave_reason,omitempty" bson:"leave_reason,omitempty"`
	Score        int64                  `json:"score" bson:"score"`
	HighScore    int64                  `json:"high_score" bson:"high_score"`
	PlayTime     time.Duration          `json:"play_time" bson:"play_time"`
	Moves        int32                  `json:"moves" bson:"moves"`
	Level        int32                  `json:"level" bson:"level"`
	Progress     float64                `json:"progress" bson:"progress"`
	Achievements []string               `json:"achievements" bson:"achievements"`
	Statistics   map[string]interface{} `json:"statistics" bson:"statistics"`
	GameData     map[string]interface{} `json:"game_data" bson:"game_data"`
	LastActivity time.Time              `json:"last_activity" bson:"last_activity"`
	CreatedAt    time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" bson:"updated_at"`
}

// GameReward 游戏奖励
type GameReward struct {
	RewardID   string     `json:"reward_id" bson:"reward_id"`
	PlayerID   uint64     `json:"player_id" bson:"player_id"`
	RewardType RewardType `json:"reward_type" bson:"reward_type"`
	Amount     int64      `json:"amount" bson:"amount"`
	ItemID     string     `json:"item_id,omitempty" bson:"item_id,omitempty"`
	ItemCount  int        `json:"item_count,omitempty" bson:"item_count,omitempty"`
	GameID     string     `json:"game_id" bson:"game_id"`
	Timestamp  time.Time  `json:"timestamp" bson:"timestamp"`
}

// NewGameSession 创建新的游戏会话
func NewGameSession(gameID string, playerID uint64, sessionToken string) *GameSession {
	now := time.Now()
	return &GameSession{
		ID:           fmt.Sprintf("%s_%d_%d", gameID, playerID, now.Unix()),
		GameID:       gameID,
		PlayerID:     playerID,
		SessionToken: sessionToken,
		Status:       PlayerStatusWaiting,
		JoinedAt:     now,
		Score:        0,
		HighScore:    0,
		PlayTime:     0,
		Moves:        0,
		Level:        1,
		Progress:     0.0,
		Achievements: make([]string, 0),
		Statistics:   make(map[string]interface{}),
		GameData:     make(map[string]interface{}),
		LastActivity: now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// UpdateStatus 更新玩家状态
func (gs *GameSession) UpdateStatus(status PlayerStatus) error {
	if !status.IsValid() {
		return fmt.Errorf("invalid player status: %v", status)
	}

	gs.Status = status
	gs.LastActivity = time.Now()
	gs.UpdatedAt = time.Now()
	return nil
}

// UpdateScore 更新分数
func (gs *GameSession) UpdateScore(score int64) {
	gs.Score = score
	if score > gs.HighScore {
		gs.HighScore = score
	}
	gs.LastActivity = time.Now()
	gs.UpdatedAt = time.Now()
}

// AddScore 增加分数
func (gs *GameSession) AddScore(points int64) {
	gs.Score += points
	if gs.Score > gs.HighScore {
		gs.HighScore = gs.Score
	}
	gs.LastActivity = time.Now()
	gs.UpdatedAt = time.Now()
}

// UpdateProgress 更新进度
func (gs *GameSession) UpdateProgress(progress float64) error {
	if progress < 0 || progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100")
	}

	gs.Progress = progress
	gs.LastActivity = time.Now()
	gs.UpdatedAt = time.Now()
	return nil
}

// UpdateLevel 更新等级
func (gs *GameSession) UpdateLevel(level int32) error {
	if level < 1 {
		return fmt.Errorf("level must be positive")
	}

	gs.Level = level
	gs.LastActivity = time.Now()
	gs.UpdatedAt = time.Now()
	return nil
}

// AddMove 增加移动次数
func (gs *GameSession) AddMove() {
	gs.Moves++
	gs.LastActivity = time.Now()
	gs.UpdatedAt = time.Now()
}

// AddPlayTime 增加游戏时间
func (gs *GameSession) AddPlayTime(duration time.Duration) {
	gs.PlayTime += duration
	gs.LastActivity = time.Now()
	gs.UpdatedAt = time.Now()
}

// AddAchievement 添加成就
func (gs *GameSession) AddAchievement(achievement string) {
	// 检查是否已存在
	for _, existing := range gs.Achievements {
		if existing == achievement {
			return // 已存在，不重复添加
		}
	}

	gs.Achievements = append(gs.Achievements, achievement)
	gs.LastActivity = time.Now()
	gs.UpdatedAt = time.Now()
}

// SetGameData 设置游戏数据
func (gs *GameSession) SetGameData(key string, value interface{}) {
	if gs.GameData == nil {
		gs.GameData = make(map[string]interface{})
	}
	gs.GameData[key] = value
	gs.LastActivity = time.Now()
	gs.UpdatedAt = time.Now()
}

// GetGameData 获取游戏数据
func (gs *GameSession) GetGameData(key string) (interface{}, bool) {
	if gs.GameData == nil {
		return nil, false
	}
	value, exists := gs.GameData[key]
	return value, exists
}

// SetStatistic 设置统计数据
func (gs *GameSession) SetStatistic(key string, value interface{}) {
	if gs.Statistics == nil {
		gs.Statistics = make(map[string]interface{})
	}
	gs.Statistics[key] = value
	gs.LastActivity = time.Now()
	gs.UpdatedAt = time.Now()
}

// GetStatistic 获取统计数据
func (gs *GameSession) GetStatistic(key string) (interface{}, bool) {
	if gs.Statistics == nil {
		return nil, false
	}
	value, exists := gs.Statistics[key]
	return value, exists
}

// Leave 离开游戏
func (gs *GameSession) Leave(reason PlayerLeaveReason) error {
	if !reason.IsValid() {
		return fmt.Errorf("invalid leave reason: %v", reason)
	}

	now := time.Now()
	gs.LeftAt = &now
	gs.LeaveReason = &reason
	gs.Status = PlayerStatusLeft
	gs.LastActivity = now
	gs.UpdatedAt = now
	return nil
}

// IsActive 检查会话是否活跃
func (gs *GameSession) IsActive() bool {
	return gs.Status == PlayerStatusWaiting || gs.Status == PlayerStatusReady || gs.Status == PlayerStatusPlaying
}

// IsFinished 检查会话是否已结束
func (gs *GameSession) IsFinished() bool {
	return gs.Status == PlayerStatusFinished || gs.Status == PlayerStatusLeft || gs.Status == PlayerStatusKicked
}

// GetDuration 获取会话持续时间
func (gs *GameSession) GetDuration() time.Duration {
	if gs.LeftAt != nil {
		return gs.LeftAt.Sub(gs.JoinedAt)
	}
	return time.Since(gs.JoinedAt)
}

// Clone 克隆游戏会话
func (gs *GameSession) Clone() *GameSession {
	clone := &GameSession{
		ID:           gs.ID,
		GameID:       gs.GameID,
		PlayerID:     gs.PlayerID,
		SessionToken: gs.SessionToken,
		Status:       gs.Status,
		JoinedAt:     gs.JoinedAt,
		Score:        gs.Score,
		HighScore:    gs.HighScore,
		PlayTime:     gs.PlayTime,
		Moves:        gs.Moves,
		Level:        gs.Level,
		Progress:     gs.Progress,
		Achievements: make([]string, len(gs.Achievements)),
		Statistics:   make(map[string]interface{}),
		GameData:     make(map[string]interface{}),
		LastActivity: gs.LastActivity,
		CreatedAt:    gs.CreatedAt,
		UpdatedAt:    gs.UpdatedAt,
	}

	// 深拷贝切片
	copy(clone.Achievements, gs.Achievements)

	// 深拷贝map
	for k, v := range gs.Statistics {
		clone.Statistics[k] = v
	}
	for k, v := range gs.GameData {
		clone.GameData[k] = v
	}

	// 深拷贝指针
	if gs.LeftAt != nil {
		leftAt := *gs.LeftAt
		clone.LeftAt = &leftAt
	}
	if gs.LeaveReason != nil {
		leaveReason := *gs.LeaveReason
		clone.LeaveReason = &leaveReason
	}

	return clone
}

// GameScore 游戏分数实体
type GameScore struct {
	ID         string                 `json:"id" bson:"_id"`
	GameID     string                 `json:"game_id" bson:"game_id"`
	PlayerID   uint64                 `json:"player_id" bson:"player_id"`
	SessionID  string                 `json:"session_id" bson:"session_id"`
	ScoreType  ScoreType              `json:"score_type" bson:"score_type"`
	Value      int64                  `json:"value" bson:"value"`
	MaxValue   int64                  `json:"max_value" bson:"max_value"`
	Multiplier float64                `json:"multiplier" bson:"multiplier"`
	Bonus      int64                  `json:"bonus" bson:"bonus"`
	Penalty    int64                  `json:"penalty" bson:"penalty"`
	FinalScore int64                  `json:"final_score" bson:"final_score"`
	Rank       int32                  `json:"rank" bson:"rank"`
	Percentile float64                `json:"percentile" bson:"percentile"`
	Breakdown  map[string]int64       `json:"breakdown" bson:"breakdown"`
	Metadata   map[string]interface{} `json:"metadata" bson:"metadata"`
	AchievedAt time.Time              `json:"achieved_at" bson:"achieved_at"`
	CreatedAt  time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at" bson:"updated_at"`
}

// NewGameScore 创建新的游戏分数
func NewGameScore(gameID string, playerID uint64, sessionID string, scoreType ScoreType) *GameScore {
	now := time.Now()
	return &GameScore{
		ID:         fmt.Sprintf("%s_%d_%s_%d", gameID, playerID, scoreType.String(), now.Unix()),
		GameID:     gameID,
		PlayerID:   playerID,
		SessionID:  sessionID,
		ScoreType:  scoreType,
		Value:      0,
		MaxValue:   0,
		Multiplier: 1.0,
		Bonus:      0,
		Penalty:    0,
		FinalScore: 0,
		Rank:       0,
		Percentile: 0.0,
		Breakdown:  make(map[string]int64),
		Metadata:   make(map[string]interface{}),
		AchievedAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// UpdateValue 更新分数值
func (gs *GameScore) UpdateValue(value int64) {
	gs.Value = value
	if value > gs.MaxValue {
		gs.MaxValue = value
	}
	gs.calculateFinalScore()
	gs.UpdatedAt = time.Now()
}

// AddValue 增加分数值
func (gs *GameScore) AddValue(points int64) {
	gs.Value += points
	if gs.Value > gs.MaxValue {
		gs.MaxValue = gs.Value
	}
	gs.calculateFinalScore()
	gs.UpdatedAt = time.Now()
}

// SetMultiplier 设置倍数
func (gs *GameScore) SetMultiplier(multiplier float64) error {
	if multiplier < 0 {
		return fmt.Errorf("multiplier cannot be negative")
	}

	gs.Multiplier = multiplier
	gs.calculateFinalScore()
	gs.UpdatedAt = time.Now()
	return nil
}

// AddBonus 添加奖励分数
func (gs *GameScore) AddBonus(bonus int64) {
	gs.Bonus += bonus
	gs.calculateFinalScore()
	gs.UpdatedAt = time.Now()
}

// AddPenalty 添加惩罚分数
func (gs *GameScore) AddPenalty(penalty int64) {
	gs.Penalty += penalty
	gs.calculateFinalScore()
	gs.UpdatedAt = time.Now()
}

// SetRank 设置排名
func (gs *GameScore) SetRank(rank int32, percentile float64) error {
	if rank < 0 {
		return fmt.Errorf("rank cannot be negative")
	}
	if percentile < 0 || percentile > 100 {
		return fmt.Errorf("percentile must be between 0 and 100")
	}

	gs.Rank = rank
	gs.Percentile = percentile
	gs.UpdatedAt = time.Now()
	return nil
}

// SetBreakdown 设置分数明细
func (gs *GameScore) SetBreakdown(breakdown map[string]int64) {
	gs.Breakdown = make(map[string]int64)
	for k, v := range breakdown {
		gs.Breakdown[k] = v
	}
	gs.UpdatedAt = time.Now()
}

// AddBreakdownItem 添加分数明细项
func (gs *GameScore) AddBreakdownItem(key string, value int64) {
	if gs.Breakdown == nil {
		gs.Breakdown = make(map[string]int64)
	}
	gs.Breakdown[key] = value
	gs.UpdatedAt = time.Now()
}

// SetMetadata 设置元数据
func (gs *GameScore) SetMetadata(key string, value interface{}) {
	if gs.Metadata == nil {
		gs.Metadata = make(map[string]interface{})
	}
	gs.Metadata[key] = value
	gs.UpdatedAt = time.Now()
}

// GetMetadata 获取元数据
func (gs *GameScore) GetMetadata(key string) (interface{}, bool) {
	if gs.Metadata == nil {
		return nil, false
	}
	value, exists := gs.Metadata[key]
	return value, exists
}

// calculateFinalScore 计算最终分数
func (gs *GameScore) calculateFinalScore() {
	baseScore := float64(gs.Value) * gs.Multiplier
	gs.FinalScore = int64(baseScore) + gs.Bonus - gs.Penalty
	if gs.FinalScore < 0 {
		gs.FinalScore = 0
	}
}

// Clone 克隆游戏分数
func (gs *GameScore) Clone() *GameScore {
	clone := &GameScore{
		ID:         gs.ID,
		GameID:     gs.GameID,
		PlayerID:   gs.PlayerID,
		SessionID:  gs.SessionID,
		ScoreType:  gs.ScoreType,
		Value:      gs.Value,
		MaxValue:   gs.MaxValue,
		Multiplier: gs.Multiplier,
		Bonus:      gs.Bonus,
		Penalty:    gs.Penalty,
		FinalScore: gs.FinalScore,
		Rank:       gs.Rank,
		Percentile: gs.Percentile,
		Breakdown:  make(map[string]int64),
		Metadata:   make(map[string]interface{}),
		AchievedAt: gs.AchievedAt,
		CreatedAt:  gs.CreatedAt,
		UpdatedAt:  gs.UpdatedAt,
	}

	// 深拷贝map
	for k, v := range gs.Breakdown {
		clone.Breakdown[k] = v
	}
	for k, v := range gs.Metadata {
		clone.Metadata[k] = v
	}

	return clone
}

// 注意：GameReward已经在文件前面定义，这里删除重复定义

// NewGameReward 创建新的游戏奖励
func NewGameReward(gameID string, playerID uint64, sessionID string, rewardType RewardType, itemID string, quantity int64) *GameReward {
	now := time.Now()
	return &GameReward{
		RewardID:   fmt.Sprintf("%s_%d_%s_%s_%d", gameID, playerID, rewardType.String(), itemID, now.Unix()),
		GameID:     gameID,
		PlayerID:   playerID,
		RewardType: rewardType,
		Amount:     quantity,
		ItemID:     itemID,
		ItemCount:  int(quantity),
		Timestamp:  now,
	}
}

// SetRarity 设置稀有度
func (gr *GameReward) SetRarity(rarity string) {
	// TODO: 实现稀有度设置
	// gr.Rarity = rarity
	// gr.UpdatedAt = time.Now()
}

// SetSource 设置来源
func (gr *GameReward) SetSource(source string) {
	// TODO: 实现来源设置
	// gr.Source = source
	// gr.UpdatedAt = time.Now()
}

// SetReason 设置原因
func (gr *GameReward) SetReason(reason string) {
	// TODO: 实现原因设置
	// gr.Reason = reason
	// gr.UpdatedAt = time.Now()
}

// SetExpiration 设置过期时间
func (gr *GameReward) SetExpiration(expiresAt time.Time) {
	// TODO: 实现过期时间设置
	// gr.ExpiresAt = &expiresAt
	// gr.UpdatedAt = time.Now()
}

// Claim 领取奖励
func (gr *GameReward) Claim() error {
	// TODO: 实现奖励领取
	// if gr.Claimed {
	// 	return fmt.Errorf("reward already claimed")
	// }

	// if gr.IsExpired() {
	// 	return fmt.Errorf("reward has expired")
	// }

	// now := time.Now()
	// gr.Claimed = true
	// gr.ClaimedAt = &now
	// gr.UpdatedAt = now
	return nil
}

// IsExpired 检查是否已过期
func (gr *GameReward) IsExpired() bool {
	// TODO: 实现过期检查
	// if gr.ExpiresAt == nil {
	// 	return false
	// }
	// return time.Now().After(*gr.ExpiresAt)
	return false
}

// IsClaimable 检查是否可领取
func (gr *GameReward) IsClaimable() bool {
	// TODO: 实现可领取检查
	// return !gr.Claimed && !gr.IsExpired()
	return true
}

// SetMetadata 设置元数据
func (gr *GameReward) SetMetadata(key string, value interface{}) {
	// TODO: 实现元数据设置
	// if gr.Metadata == nil {
	// 	gr.Metadata = make(map[string]interface{})
	// }
	// gr.Metadata[key] = value
	// gr.UpdatedAt = time.Now()
}

// GetMetadata 获取元数据
func (gr *GameReward) GetMetadata(key string) (interface{}, bool) {
	// TODO: 实现元数据获取
	// if gr.Metadata == nil {
	// 	return nil, false
	// }
	// value, exists := gr.Metadata[key]
	// return value, exists
	return nil, false
}

// Clone 克隆游戏奖励
func (gr *GameReward) Clone() *GameReward {
	clone := &GameReward{
		RewardID:   gr.RewardID,
		GameID:     gr.GameID,
		PlayerID:   gr.PlayerID,
		RewardType: gr.RewardType,
		Amount:     gr.Amount,
		ItemID:     gr.ItemID,
		ItemCount:  gr.ItemCount,
		Timestamp:  gr.Timestamp,
	}

	// TODO: 实现深拷贝
	// 深拷贝map
	// for k, v := range gr.Metadata {
	// 	clone.Metadata[k] = v
	// }

	// 深拷贝指针
	// if gr.ClaimedAt != nil {
	// 	claimedAt := *gr.ClaimedAt
	// 	clone.ClaimedAt = &claimedAt
	// }
	// if gr.ExpiresAt != nil {
	// 	expiresAt := *gr.ExpiresAt
	// 	clone.ExpiresAt = &expiresAt
	// }

	return clone
}

// GameAchievement 游戏成就实体
type GameAchievement struct {
	ID            string                 `json:"id" bson:"_id"`
	GameID        string                 `json:"game_id" bson:"game_id"`
	PlayerID      uint64                 `json:"player_id" bson:"player_id"`
	SessionID     string                 `json:"session_id" bson:"session_id"`
	AchievementID string                 `json:"achievement_id" bson:"achievement_id"`
	Name          string                 `json:"name" bson:"name"`
	Description   string                 `json:"description" bson:"description"`
	Category      string                 `json:"category" bson:"category"`
	Rarity        string                 `json:"rarity" bson:"rarity"`
	Points        int64                  `json:"points" bson:"points"`
	Progress      float64                `json:"progress" bson:"progress"`
	MaxProgress   float64                `json:"max_progress" bson:"max_progress"`
	Completed     bool                   `json:"completed" bson:"completed"`
	CompletedAt   *time.Time             `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
	Unlocked      bool                   `json:"unlocked" bson:"unlocked"`
	UnlockedAt    *time.Time             `json:"unlocked_at,omitempty" bson:"unlocked_at,omitempty"`
	Conditions    map[string]interface{} `json:"conditions" bson:"conditions"`
	Rewards       []string               `json:"rewards" bson:"rewards"`
	Metadata      map[string]interface{} `json:"metadata" bson:"metadata"`
	CreatedAt     time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at" bson:"updated_at"`
}

// NewGameAchievement 创建新的游戏成就
func NewGameAchievement(gameID string, playerID uint64, sessionID string, achievementID string, name string) *GameAchievement {
	now := time.Now()
	return &GameAchievement{
		ID:            fmt.Sprintf("%s_%d_%s_%d", gameID, playerID, achievementID, now.Unix()),
		GameID:        gameID,
		PlayerID:      playerID,
		SessionID:     sessionID,
		AchievementID: achievementID,
		Name:          name,
		Description:   "",
		Category:      "general",
		Rarity:        "common",
		Points:        0,
		Progress:      0.0,
		MaxProgress:   100.0,
		Completed:     false,
		Unlocked:      false,
		Conditions:    make(map[string]interface{}),
		Rewards:       make([]string, 0),
		Metadata:      make(map[string]interface{}),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// UpdateProgress 更新进度
func (ga *GameAchievement) UpdateProgress(progress float64) error {
	if progress < 0 {
		return fmt.Errorf("progress cannot be negative")
	}
	if progress > ga.MaxProgress {
		progress = ga.MaxProgress
	}

	ga.Progress = progress

	// 检查是否完成
	if progress >= ga.MaxProgress && !ga.Completed {
		ga.Complete()
	}

	ga.UpdatedAt = time.Now()
	return nil
}

// AddProgress 增加进度
func (ga *GameAchievement) AddProgress(delta float64) error {
	return ga.UpdateProgress(ga.Progress + delta)
}

// Complete 完成成就
func (ga *GameAchievement) Complete() {
	if ga.Completed {
		return
	}

	now := time.Now()
	ga.Completed = true
	ga.CompletedAt = &now
	ga.Progress = ga.MaxProgress
	ga.UpdatedAt = now
}

// Unlock 解锁成就
func (ga *GameAchievement) Unlock() {
	if ga.Unlocked {
		return
	}

	now := time.Now()
	ga.Unlocked = true
	ga.UnlockedAt = &now
	ga.UpdatedAt = now
}

// SetDescription 设置描述
func (ga *GameAchievement) SetDescription(description string) {
	ga.Description = description
	ga.UpdatedAt = time.Now()
}

// SetCategory 设置分类
func (ga *GameAchievement) SetCategory(category string) {
	ga.Category = category
	ga.UpdatedAt = time.Now()
}

// SetRarity 设置稀有度
func (ga *GameAchievement) SetRarity(rarity string) {
	ga.Rarity = rarity
	ga.UpdatedAt = time.Now()
}

// SetPoints 设置积分
func (ga *GameAchievement) SetPoints(points int64) {
	ga.Points = points
	ga.UpdatedAt = time.Now()
}

// SetMaxProgress 设置最大进度
func (ga *GameAchievement) SetMaxProgress(maxProgress float64) error {
	if maxProgress <= 0 {
		return fmt.Errorf("max progress must be positive")
	}

	ga.MaxProgress = maxProgress
	ga.UpdatedAt = time.Now()
	return nil
}

// AddReward 添加奖励
func (ga *GameAchievement) AddReward(reward string) {
	// 检查是否已存在
	for _, existing := range ga.Rewards {
		if existing == reward {
			return // 已存在，不重复添加
		}
	}

	ga.Rewards = append(ga.Rewards, reward)
	ga.UpdatedAt = time.Now()
}

// SetCondition 设置条件
func (ga *GameAchievement) SetCondition(key string, value interface{}) {
	if ga.Conditions == nil {
		ga.Conditions = make(map[string]interface{})
	}
	ga.Conditions[key] = value
	ga.UpdatedAt = time.Now()
}

// GetCondition 获取条件
func (ga *GameAchievement) GetCondition(key string) (interface{}, bool) {
	if ga.Conditions == nil {
		return nil, false
	}
	value, exists := ga.Conditions[key]
	return value, exists
}

// SetMetadata 设置元数据
func (ga *GameAchievement) SetMetadata(key string, value interface{}) {
	if ga.Metadata == nil {
		ga.Metadata = make(map[string]interface{})
	}
	ga.Metadata[key] = value
	ga.UpdatedAt = time.Now()
}

// GetMetadata 获取元数据
func (ga *GameAchievement) GetMetadata(key string) (interface{}, bool) {
	if ga.Metadata == nil {
		return nil, false
	}
	value, exists := ga.Metadata[key]
	return value, exists
}

// GetProgressPercentage 获取进度百分比
func (ga *GameAchievement) GetProgressPercentage() float64 {
	if ga.MaxProgress == 0 {
		return 0
	}
	return (ga.Progress / ga.MaxProgress) * 100
}

// IsCompleted 检查是否已完成
func (ga *GameAchievement) IsCompleted() bool {
	return ga.Completed
}

// IsUnlocked 检查是否已解锁
func (ga *GameAchievement) IsUnlocked() bool {
	return ga.Unlocked
}

// Clone 克隆游戏成就
func (ga *GameAchievement) Clone() *GameAchievement {
	clone := &GameAchievement{
		ID:            ga.ID,
		GameID:        ga.GameID,
		PlayerID:      ga.PlayerID,
		SessionID:     ga.SessionID,
		AchievementID: ga.AchievementID,
		Name:          ga.Name,
		Description:   ga.Description,
		Category:      ga.Category,
		Rarity:        ga.Rarity,
		Points:        ga.Points,
		Progress:      ga.Progress,
		MaxProgress:   ga.MaxProgress,
		Completed:     ga.Completed,
		Unlocked:      ga.Unlocked,
		Conditions:    make(map[string]interface{}),
		Rewards:       make([]string, len(ga.Rewards)),
		Metadata:      make(map[string]interface{}),
		CreatedAt:     ga.CreatedAt,
		UpdatedAt:     ga.UpdatedAt,
	}

	// 深拷贝切片
	copy(clone.Rewards, ga.Rewards)

	// 深拷贝map
	for k, v := range ga.Conditions {
		clone.Conditions[k] = v
	}
	for k, v := range ga.Metadata {
		clone.Metadata[k] = v
	}

	// 深拷贝指针
	if ga.CompletedAt != nil {
		completedAt := *ga.CompletedAt
		clone.CompletedAt = &completedAt
	}
	if ga.UnlockedAt != nil {
		unlockedAt := *ga.UnlockedAt
		clone.UnlockedAt = &unlockedAt
	}

	return clone
}

// 常量定义

const (
	// 会话相关常量
	MaxSessionDuration    = 24 * time.Hour   // 最大会话持续时间
	SessionTimeoutWarning = 5 * time.Minute  // 会话超时警告时间
	MaxInactiveDuration   = 30 * time.Minute // 最大非活跃时间

	// 分数相关常量
	MaxScoreValue     = int64(999999999) // 最大分数值
	MinScoreValue     = int64(0)         // 最小分数值
	DefaultMultiplier = 1.0              // 默认倍数
	MaxMultiplier     = 10.0             // 最大倍数

	// 奖励相关常量
	MaxRewardQuantity = int64(999999)      // 最大奖励数量
	DefaultRewardTTL  = 7 * 24 * time.Hour // 默认奖励过期时间

	// 成就相关常量
	MaxAchievementProgress = 100.0        // 最大成就进度
	MaxAchievementPoints   = int64(10000) // 最大成就积分
)

// 验证函数

// ValidateGameSession 验证游戏会话
func ValidateGameSession(session *GameSession) error {
	if session == nil {
		return fmt.Errorf("session cannot be nil")
	}

	if session.GameID == "" {
		return fmt.Errorf("game_id cannot be empty")
	}

	if session.PlayerID == 0 {
		return fmt.Errorf("player_id cannot be zero")
	}

	if session.SessionToken == "" {
		return fmt.Errorf("session_token cannot be empty")
	}

	if !session.Status.IsValid() {
		return fmt.Errorf("invalid player status: %v", session.Status)
	}

	if session.Progress < 0 || session.Progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100")
	}

	if session.Level < 1 {
		return fmt.Errorf("level must be positive")
	}

	return nil
}

// ValidateGameScore 验证游戏分数
func ValidateGameScore(score *GameScore) error {
	if score == nil {
		return fmt.Errorf("score cannot be nil")
	}

	if score.GameID == "" {
		return fmt.Errorf("game_id cannot be empty")
	}

	if score.PlayerID == 0 {
		return fmt.Errorf("player_id cannot be zero")
	}

	if !score.ScoreType.IsValid() {
		return fmt.Errorf("invalid score type: %v", score.ScoreType)
	}

	if score.Value < MinScoreValue || score.Value > MaxScoreValue {
		return fmt.Errorf("score value must be between %d and %d", MinScoreValue, MaxScoreValue)
	}

	if score.Multiplier < 0 || score.Multiplier > MaxMultiplier {
		return fmt.Errorf("multiplier must be between 0 and %f", MaxMultiplier)
	}

	if score.Percentile < 0 || score.Percentile > 100 {
		return fmt.Errorf("percentile must be between 0 and 100")
	}

	return nil
}

// ValidateGameReward 验证游戏奖励
func ValidateGameReward(reward *GameReward) error {
	if reward == nil {
		return fmt.Errorf("reward cannot be nil")
	}

	if reward.GameID == "" {
		return fmt.Errorf("game_id cannot be empty")
	}

	if reward.PlayerID == 0 {
		return fmt.Errorf("player_id cannot be zero")
	}

	if !reward.RewardType.IsValid() {
		return fmt.Errorf("invalid reward type: %v", reward.RewardType)
	}

	if reward.ItemID == "" {
		return fmt.Errorf("item_id cannot be empty")
	}

	if reward.Amount <= 0 || reward.Amount > MaxRewardQuantity {
		return fmt.Errorf("amount must be between 1 and %d", MaxRewardQuantity)
	}

	// TODO: 实现过期时间检查
	// if reward.ExpiresAt != nil && reward.ExpiresAt.Before(time.Now()) {
	// 	return fmt.Errorf("expiration time cannot be in the past")
	// }

	return nil
}

// ValidateGameAchievement 验证游戏成就
func ValidateGameAchievement(achievement *GameAchievement) error {
	if achievement == nil {
		return fmt.Errorf("achievement cannot be nil")
	}

	if achievement.GameID == "" {
		return fmt.Errorf("game_id cannot be empty")
	}

	if achievement.PlayerID == 0 {
		return fmt.Errorf("player_id cannot be zero")
	}

	if achievement.AchievementID == "" {
		return fmt.Errorf("achievement_id cannot be empty")
	}

	if achievement.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if achievement.Progress < 0 || achievement.Progress > achievement.MaxProgress {
		return fmt.Errorf("progress must be between 0 and %f", achievement.MaxProgress)
	}

	if achievement.MaxProgress <= 0 || achievement.MaxProgress > MaxAchievementProgress {
		return fmt.Errorf("max_progress must be between 0 and %f", MaxAchievementProgress)
	}

	if achievement.Points < 0 || achievement.Points > MaxAchievementPoints {
		return fmt.Errorf("points must be between 0 and %d", MaxAchievementPoints)
	}

	return nil
}
