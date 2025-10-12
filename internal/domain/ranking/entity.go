package ranking

import (
	"fmt"
	"sync"
	"time"
)

// RankEntry 排行榜条目实体
type RankEntry struct {
	// 基础信息
	ID       string `json:"id" bson:"_id"`
	PlayerID uint64 `json:"player_id" bson:"player_id"`
	RankID   uint32 `json:"rank_id" bson:"rank_id"`

	// 分数信息
	Score     int64 `json:"score" bson:"score"`           // 真实分数
	TimeScore int64 `json:"time_score" bson:"time_score"` // 带时间因子的分数
	Rank      int64 `json:"rank" bson:"rank"`             // 当前排名

	// 玩家信息
	PlayerName   string `json:"player_name" bson:"player_name"`
	PlayerLevel  uint32 `json:"player_level" bson:"player_level"`
	PlayerAvatar string `json:"player_avatar" bson:"player_avatar"`
	PlayerTitle  string `json:"player_title" bson:"player_title"`

	// 状态信息
	IsActive     bool                 `json:"is_active" bson:"is_active"`
	LastActive   time.Time            `json:"last_active" bson:"last_active"`
	ScoreHistory []*ScoreHistoryEntry `json:"score_history" bson:"score_history"`

	// 排名变化
	PreviousRank *int64 `json:"previous_rank,omitempty" bson:"previous_rank,omitempty"`
	RankChange   int64  `json:"rank_change" bson:"rank_change"`
	BestRank     int64  `json:"best_rank" bson:"best_rank"`
	WorstRank    int64  `json:"worst_rank" bson:"worst_rank"`

	// 统计信息
	TotalUpdates    int64     `json:"total_updates" bson:"total_updates"`
	ConsecutiveDays int32     `json:"consecutive_days" bson:"consecutive_days"`
	FirstEntryTime  time.Time `json:"first_entry_time" bson:"first_entry_time"`
	LastUpdateTime  time.Time `json:"last_update_time" bson:"last_update_time"`

	// 奖励信息
	RewardsEarned    []*RankRewardEarned `json:"rewards_earned" bson:"rewards_earned"`
	LastRewardTime   *time.Time          `json:"last_reward_time,omitempty" bson:"last_reward_time,omitempty"`
	TotalRewardValue int64               `json:"total_reward_value" bson:"total_reward_value"`

	// 元数据
	Metadata   map[string]interface{} `json:"metadata" bson:"metadata"`
	Tags       []string               `json:"tags" bson:"tags"`
	CustomData map[string]interface{} `json:"custom_data" bson:"custom_data"`

	// 时间戳
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	// 内部状态
	mutex sync.RWMutex `json:"-" bson:"-"`
}

// GetID 获取条目ID
func (re *RankEntry) GetID() string {
	re.mutex.RLock()
	defer re.mutex.RUnlock()
	return re.ID
}

// GetRankingID 获取排行榜ID
func (re *RankEntry) GetRankingID() uint32 {
	re.mutex.RLock()
	defer re.mutex.RUnlock()
	return re.RankID
}

// GetPlayerID 获取玩家ID
func (re *RankEntry) GetPlayerID() uint64 {
	re.mutex.RLock()
	defer re.mutex.RUnlock()
	return re.PlayerID
}

// GetRank 获取排名
func (re *RankEntry) GetRank() int64 {
	re.mutex.RLock()
	defer re.mutex.RUnlock()
	return re.Rank
}

// GetScore 获取分数
func (re *RankEntry) GetScore() int64 {
	re.mutex.RLock()
	defer re.mutex.RUnlock()
	return re.Score
}

// GetPrevRank 获取前一排名
func (re *RankEntry) GetPrevRank() *int64 {
	re.mutex.RLock()
	defer re.mutex.RUnlock()
	return re.PreviousRank
}

// GetPrevScore 获取前一分数 (暂时返回0，可根据需要实现)
func (re *RankEntry) GetPrevScore() int64 {
	re.mutex.RLock()
	defer re.mutex.RUnlock()
	// 从历史记录中获取前一分数
	if len(re.ScoreHistory) > 0 {
		return re.ScoreHistory[len(re.ScoreHistory)-1].Score
	}
	return 0
}

// GetMetadata 获取元数据
func (re *RankEntry) GetMetadata() map[string]interface{} {
	re.mutex.RLock()
	defer re.mutex.RUnlock()
	return re.Metadata
}

// GetCreatedAt 获取创建时间
func (re *RankEntry) GetCreatedAt() time.Time {
	re.mutex.RLock()
	defer re.mutex.RUnlock()
	return re.CreatedAt
}

// GetUpdatedAt 获取更新时间
func (re *RankEntry) GetUpdatedAt() time.Time {
	re.mutex.RLock()
	defer re.mutex.RUnlock()
	return re.UpdatedAt
}

// SetRank 设置排名
func (re *RankEntry) SetRank(rank int64) {
	re.mutex.Lock()
	defer re.mutex.Unlock()
	re.Rank = rank
	re.UpdatedAt = time.Now()
}

// SetPrevious 设置前一排名和分数
func (re *RankEntry) SetPrevious(prevRank *int64, prevScore int64) {
	re.mutex.Lock()
	defer re.mutex.Unlock()
	re.PreviousRank = prevRank
	// 可以添加前一分数的存储逻辑
	re.UpdatedAt = time.Now()
}

// SetMetadata 设置元数据
func (re *RankEntry) SetMetadata(metadata map[string]interface{}) {
	re.mutex.Lock()
	defer re.mutex.Unlock()
	re.Metadata = metadata
	re.UpdatedAt = time.Now()
}

// NewRankEntry 创建新的排行榜条目
func NewRankEntry(playerID uint64, timeScore, realScore int64, metadata map[string]interface{}) *RankEntry {
	now := time.Now()
	return &RankEntry{
		ID:               generateRankEntryID(playerID, now),
		PlayerID:         playerID,
		Score:            realScore,
		TimeScore:        timeScore,
		Rank:             0, // 将在排序后设置
		IsActive:         true,
		LastActive:       now,
		ScoreHistory:     make([]*ScoreHistoryEntry, 0),
		RankChange:       0,
		BestRank:         0,
		WorstRank:        0,
		TotalUpdates:     1,
		ConsecutiveDays:  1,
		FirstEntryTime:   now,
		LastUpdateTime:   now,
		RewardsEarned:    make([]*RankRewardEarned, 0),
		TotalRewardValue: 0,
		Metadata:         metadata,
		Tags:             make([]string, 0),
		CustomData:       make(map[string]interface{}),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// NewRankEntryFromRepository 创建新的排行榜条目（用于repository）
func NewRankEntryFromRepository(entryID string, rankingID uint32, playerID uint64, score int64) *RankEntry {
	now := time.Now()
	return &RankEntry{
		ID:               entryID,
		PlayerID:         playerID,
		RankID:           rankingID,
		Score:            score,
		TimeScore:        score,
		Rank:             0,
		IsActive:         true,
		LastActive:       now,
		ScoreHistory:     make([]*ScoreHistoryEntry, 0),
		RankChange:       0,
		BestRank:         0,
		WorstRank:        0,
		TotalUpdates:     1,
		ConsecutiveDays:  1,
		FirstEntryTime:   now,
		LastUpdateTime:   now,
		RewardsEarned:    make([]*RankRewardEarned, 0),
		TotalRewardValue: 0,
		Metadata:         make(map[string]interface{}),
		Tags:             make([]string, 0),
		CustomData:       make(map[string]interface{}),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// UpdateScore 更新分数
func (re *RankEntry) UpdateScore(timeScore, realScore int64, metadata map[string]interface{}) {
	re.mutex.Lock()
	defer re.mutex.Unlock()

	// 记录历史分数
	re.addScoreHistory(re.Score, re.TimeScore)

	// 更新分数
	re.Score = realScore
	re.TimeScore = timeScore
	re.TotalUpdates++
	re.LastUpdateTime = time.Now()
	re.UpdatedAt = time.Now()
	re.LastActive = time.Now()

	// 更新元数据
	if metadata != nil {
		for k, v := range metadata {
			re.Metadata[k] = v
		}
	}

	// 更新连续天数
	re.updateConsecutiveDays()
}

// UpdateRank 更新排名
func (re *RankEntry) UpdateRank(newRank int64) {
	re.mutex.Lock()
	defer re.mutex.Unlock()

	// 记录排名变化
	if re.Rank != 0 {
		re.PreviousRank = &re.Rank
		re.RankChange = re.Rank - newRank
	}

	// 更新排名
	re.Rank = newRank

	// 更新最佳和最差排名
	if re.BestRank == 0 || newRank < re.BestRank {
		re.BestRank = newRank
	}
	if re.WorstRank == 0 || newRank > re.WorstRank {
		re.WorstRank = newRank
	}

	re.UpdatedAt = time.Now()
}

// UpdatePlayerInfo 更新玩家信息
func (re *RankEntry) UpdatePlayerInfo(name string, level uint32, avatar, title string) {
	re.mutex.Lock()
	defer re.mutex.Unlock()

	re.PlayerName = name
	re.PlayerLevel = level
	re.PlayerAvatar = avatar
	re.PlayerTitle = title
	re.UpdatedAt = time.Now()
}

// AddReward 添加奖励
func (re *RankEntry) AddReward(reward *RankRewardEarned) {
	re.mutex.Lock()
	defer re.mutex.Unlock()

	re.RewardsEarned = append(re.RewardsEarned, reward)
	re.TotalRewardValue += reward.Value
	now := time.Now()
	re.LastRewardTime = &now
	re.UpdatedAt = now
}

// SetActive 设置活跃状态
func (re *RankEntry) SetActive(active bool) {
	re.mutex.Lock()
	defer re.mutex.Unlock()

	re.IsActive = active
	if active {
		re.LastActive = time.Now()
	}
	re.UpdatedAt = time.Now()
}

// AddTag 添加标签
func (re *RankEntry) AddTag(tag string) {
	re.mutex.Lock()
	defer re.mutex.Unlock()

	// 检查是否已存在
	for _, existingTag := range re.Tags {
		if existingTag == tag {
			return
		}
	}

	re.Tags = append(re.Tags, tag)
	re.UpdatedAt = time.Now()
}

// RemoveTag 移除标签
func (re *RankEntry) RemoveTag(tag string) {
	re.mutex.Lock()
	defer re.mutex.Unlock()

	for i, existingTag := range re.Tags {
		if existingTag == tag {
			re.Tags = append(re.Tags[:i], re.Tags[i+1:]...)
			re.UpdatedAt = time.Now()
			return
		}
	}
}

// SetCustomData 设置自定义数据
func (re *RankEntry) SetCustomData(key string, value interface{}) {
	re.mutex.Lock()
	defer re.mutex.Unlock()

	if re.CustomData == nil {
		re.CustomData = make(map[string]interface{})
	}
	re.CustomData[key] = value
	re.UpdatedAt = time.Now()
}

// GetCustomData 获取自定义数据
func (re *RankEntry) GetCustomData(key string) (interface{}, bool) {
	re.mutex.RLock()
	defer re.mutex.RUnlock()

	value, exists := re.CustomData[key]
	return value, exists
}

// GetScoreHistory 获取分数历史
func (re *RankEntry) GetScoreHistory(limit int) []*ScoreHistoryEntry {
	re.mutex.RLock()
	defer re.mutex.RUnlock()

	if limit <= 0 || limit > len(re.ScoreHistory) {
		limit = len(re.ScoreHistory)
	}

	// 返回最近的记录
	start := len(re.ScoreHistory) - limit
	history := make([]*ScoreHistoryEntry, limit)
	copy(history, re.ScoreHistory[start:])
	return history
}

// GetRecentRewards 获取最近的奖励
func (re *RankEntry) GetRecentRewards(limit int) []*RankRewardEarned {
	re.mutex.RLock()
	defer re.mutex.RUnlock()

	if limit <= 0 || limit > len(re.RewardsEarned) {
		limit = len(re.RewardsEarned)
	}

	// 返回最近的奖励
	start := len(re.RewardsEarned) - limit
	rewards := make([]*RankRewardEarned, limit)
	copy(rewards, re.RewardsEarned[start:])
	return rewards
}

// IsRankImproved 检查排名是否提升
func (re *RankEntry) IsRankImproved() bool {
	re.mutex.RLock()
	defer re.mutex.RUnlock()

	return re.RankChange > 0
}

// GetRankChangeDirection 获取排名变化方向
func (re *RankEntry) GetRankChangeDirection() string {
	re.mutex.RLock()
	defer re.mutex.RUnlock()

	if re.RankChange > 0 {
		return "up"
	} else if re.RankChange < 0 {
		return "down"
	}
	return "unchanged"
}

// 私有方法

// addScoreHistory 添加分数历史
func (re *RankEntry) addScoreHistory(score, timeScore int64) {
	history := &ScoreHistoryEntry{
		Score:     score,
		TimeScore: timeScore,
		Timestamp: time.Now(),
	}

	re.ScoreHistory = append(re.ScoreHistory, history)

	// 限制历史记录数量
	if len(re.ScoreHistory) > MaxScoreHistorySize {
		re.ScoreHistory = re.ScoreHistory[1:]
	}
}

// updateConsecutiveDays 更新连续天数
func (re *RankEntry) updateConsecutiveDays() {
	now := time.Now()
	lastUpdate := re.LastUpdateTime

	// 检查是否是连续的天
	if now.Sub(lastUpdate) <= 24*time.Hour {
		// 同一天或连续天
		if now.Day() != lastUpdate.Day() {
			re.ConsecutiveDays++
		}
	} else {
		// 中断了连续性
		re.ConsecutiveDays = 1
	}
}

// ScoreHistoryEntry 分数历史条目
type ScoreHistoryEntry struct {
	Score     int64                  `json:"score" bson:"score"`
	TimeScore int64                  `json:"time_score" bson:"time_score"`
	Timestamp time.Time              `json:"timestamp" bson:"timestamp"`
	Reason    string                 `json:"reason,omitempty" bson:"reason,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

// RankRewardEarned 已获得的排行榜奖励
type RankRewardEarned struct {
	RewardID    string                 `json:"reward_id" bson:"reward_id"`
	RewardType  string                 `json:"reward_type" bson:"reward_type"`
	Quantity    int64                  `json:"quantity" bson:"quantity"`
	Value       int64                  `json:"value" bson:"value"`
	Rank        int64                  `json:"rank" bson:"rank"`
	RankTier    string                 `json:"rank_tier" bson:"rank_tier"`
	EarnedAt    time.Time              `json:"earned_at" bson:"earned_at"`
	ClaimedAt   *time.Time             `json:"claimed_at,omitempty" bson:"claimed_at,omitempty"`
	IsClaimed   bool                   `json:"is_claimed" bson:"is_claimed"`
	Description string                 `json:"description" bson:"description"`
	Metadata    map[string]interface{} `json:"metadata" bson:"metadata"`
}

// NewRankRewardEarned 创建新的已获得奖励
func NewRankRewardEarned(rewardID, rewardType string, quantity, value, rank int64, tier, description string) *RankRewardEarned {
	return &RankRewardEarned{
		RewardID:    rewardID,
		RewardType:  rewardType,
		Quantity:    quantity,
		Value:       value,
		Rank:        rank,
		RankTier:    tier,
		EarnedAt:    time.Now(),
		IsClaimed:   false,
		Description: description,
		Metadata:    make(map[string]interface{}),
	}
}

// Claim 领取奖励
func (rre *RankRewardEarned) Claim() {
	if !rre.IsClaimed {
		now := time.Now()
		rre.ClaimedAt = &now
		rre.IsClaimed = true
	}
}

// Blacklist 黑名单实体
type Blacklist struct {
	// 基础信息
	ID     string `json:"id" bson:"_id"`
	RankID uint32 `json:"rank_id" bson:"rank_id"`

	// 黑名单玩家
	Players map[uint64]*BlacklistEntry `json:"players" bson:"players"`

	// 统计信息
	TotalBlacklisted int64     `json:"total_blacklisted" bson:"total_blacklisted"`
	LastUpdated      time.Time `json:"last_updated" bson:"last_updated"`

	// 时间戳
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	// 内部状态
	mutex sync.RWMutex `json:"-" bson:"-"`
}

// NewBlacklist 创建新的黑名单
func NewBlacklist(rankID uint32) *Blacklist {
	now := time.Now()
	return &Blacklist{
		ID:               generateBlacklistID(rankID),
		RankID:           rankID,
		Players:          make(map[uint64]*BlacklistEntry),
		TotalBlacklisted: 0,
		LastUpdated:      now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// AddPlayer 添加玩家到黑名单
func (bl *Blacklist) AddPlayer(playerID uint64, reason string) error {
	bl.mutex.Lock()
	defer bl.mutex.Unlock()

	// 检查是否已在黑名单中
	if _, exists := bl.Players[playerID]; exists {
		return NewPlayerAlreadyBlacklistedError(playerID, bl.RankID)
	}

	// 添加到黑名单
	entry := NewBlacklistEntry(playerID, reason)
	bl.Players[playerID] = entry
	bl.TotalBlacklisted++
	bl.LastUpdated = time.Now()
	bl.UpdatedAt = time.Now()

	return nil
}

// RemovePlayer 从黑名单移除玩家
func (bl *Blacklist) RemovePlayer(playerID uint64) error {
	bl.mutex.Lock()
	defer bl.mutex.Unlock()

	// 检查是否在黑名单中
	if _, exists := bl.Players[playerID]; !exists {
		return NewPlayerNotBlacklistedError(playerID, bl.RankID)
	}

	// 从黑名单移除
	delete(bl.Players, playerID)
	bl.TotalBlacklisted--
	bl.LastUpdated = time.Now()
	bl.UpdatedAt = time.Now()

	return nil
}

// IsBlacklisted 检查玩家是否在黑名单中
func (bl *Blacklist) IsBlacklisted(playerID uint64) bool {
	bl.mutex.RLock()
	defer bl.mutex.RUnlock()

	_, exists := bl.Players[playerID]
	return exists
}

// GetBlacklistEntry 获取黑名单条目
func (bl *Blacklist) GetBlacklistEntry(playerID uint64) (*BlacklistEntry, bool) {
	bl.mutex.RLock()
	defer bl.mutex.RUnlock()

	entry, exists := bl.Players[playerID]
	return entry, exists
}

// GetAllPlayers 获取所有黑名单玩家
func (bl *Blacklist) GetAllPlayers() []*BlacklistEntry {
	bl.mutex.RLock()
	defer bl.mutex.RUnlock()

	entries := make([]*BlacklistEntry, 0, len(bl.Players))
	for _, entry := range bl.Players {
		entries = append(entries, entry)
	}
	return entries
}

// GetPlayersByReason 根据原因获取黑名单玩家
func (bl *Blacklist) GetPlayersByReason(reason string) []*BlacklistEntry {
	bl.mutex.RLock()
	defer bl.mutex.RUnlock()

	entries := make([]*BlacklistEntry, 0)
	for _, entry := range bl.Players {
		if entry.Reason == reason {
			entries = append(entries, entry)
		}
	}
	return entries
}

// GetPlayerIDs 获取所有黑名单玩家ID
func (bl *Blacklist) GetPlayerIDs() []uint64 {
	bl.mutex.RLock()
	defer bl.mutex.RUnlock()

	ids := make([]uint64, 0, len(bl.Players))
	for playerID := range bl.Players {
		ids = append(ids, playerID)
	}
	return ids
}

// Clear 清空黑名单
func (bl *Blacklist) Clear() {
	bl.mutex.Lock()
	defer bl.mutex.Unlock()

	bl.Players = make(map[uint64]*BlacklistEntry)
	bl.TotalBlacklisted = 0
	bl.LastUpdated = time.Now()
	bl.UpdatedAt = time.Now()
}

// Clone 克隆黑名单
func (bl *Blacklist) Clone() *Blacklist {
	bl.mutex.RLock()
	defer bl.mutex.RUnlock()

	clone := &Blacklist{
		ID:               bl.ID,
		RankID:           bl.RankID,
		Players:          make(map[uint64]*BlacklistEntry),
		TotalBlacklisted: bl.TotalBlacklisted,
		LastUpdated:      bl.LastUpdated,
		CreatedAt:        bl.CreatedAt,
		UpdatedAt:        bl.UpdatedAt,
	}

	// 深拷贝玩家条目
	for playerID, entry := range bl.Players {
		entryCopy := *entry
		clone.Players[playerID] = &entryCopy
	}

	return clone
}

// BlacklistEntry 黑名单条目
type BlacklistEntry struct {
	PlayerID    uint64                 `json:"player_id" bson:"player_id"`
	Reason      string                 `json:"reason" bson:"reason"`
	AddedAt     time.Time              `json:"added_at" bson:"added_at"`
	AddedBy     string                 `json:"added_by" bson:"added_by"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
	IsPermanent bool                   `json:"is_permanent" bson:"is_permanent"`
	Metadata    map[string]interface{} `json:"metadata" bson:"metadata"`
}

// NewBlacklistEntry 创建新的黑名单条目
func NewBlacklistEntry(playerID uint64, reason string) *BlacklistEntry {
	return &BlacklistEntry{
		PlayerID:    playerID,
		Reason:      reason,
		AddedAt:     time.Now(),
		IsPermanent: true,
		Metadata:    make(map[string]interface{}),
	}
}

// NewTemporaryBlacklistEntry 创建临时黑名单条目
func NewTemporaryBlacklistEntry(playerID uint64, reason string, duration time.Duration) *BlacklistEntry {
	expiresAt := time.Now().Add(duration)
	return &BlacklistEntry{
		PlayerID:    playerID,
		Reason:      reason,
		AddedAt:     time.Now(),
		ExpiresAt:   &expiresAt,
		IsPermanent: false,
		Metadata:    make(map[string]interface{}),
	}
}

// IsExpired 检查是否已过期
func (ble *BlacklistEntry) IsExpired() bool {
	if ble.IsPermanent || ble.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*ble.ExpiresAt)
}

// GetRemainingTime 获取剩余时间
func (ble *BlacklistEntry) GetRemainingTime() time.Duration {
	if ble.IsPermanent || ble.ExpiresAt == nil {
		return 0
	}

	remaining := ble.ExpiresAt.Sub(time.Now())
	if remaining < 0 {
		return 0
	}
	return remaining
}

// SetExpiration 设置过期时间
func (ble *BlacklistEntry) SetExpiration(expiresAt time.Time) {
	ble.ExpiresAt = &expiresAt
	ble.IsPermanent = false
}

// SetPermanent 设置为永久
func (ble *BlacklistEntry) SetPermanent() {
	ble.ExpiresAt = nil
	ble.IsPermanent = true
}

// SetMetadata 设置元数据
func (ble *BlacklistEntry) SetMetadata(key string, value interface{}) {
	if ble.Metadata == nil {
		ble.Metadata = make(map[string]interface{})
	}
	ble.Metadata[key] = value
}

// GetMetadata 获取元数据
func (ble *BlacklistEntry) GetMetadata(key string) (interface{}, bool) {
	value, exists := ble.Metadata[key]
	return value, exists
}

// 常量定义

const (
	// 历史记录限制
	MaxScoreHistorySize  = 100
	MaxRewardHistorySize = 50

	// 黑名单限制
	MaxBlacklistSize = 10000

	// 标签限制
	MaxTagsPerEntry = 10
	MaxTagLength    = 50

	// 元数据限制
	MaxMetadataSize = 1024 * 1024 // 1MB
)

// 辅助函数

// generateRankEntryID 生成排行榜条目ID
func generateRankEntryID(playerID uint64, timestamp time.Time) string {
	return fmt.Sprintf("rank_entry_%d_%d", playerID, timestamp.Unix())
}

// generateBlacklistID 生成黑名单ID
func generateBlacklistID(rankID uint32) string {
	return fmt.Sprintf("blacklist_%d", rankID)
}

// 验证方法

// ValidateRankEntry 验证排行榜条目
func ValidateRankEntry(entry *RankEntry) error {
	if entry == nil {
		return fmt.Errorf("rank entry cannot be nil")
	}

	if entry.PlayerID == 0 {
		return fmt.Errorf("player ID cannot be zero")
	}

	if entry.RankID == 0 {
		return fmt.Errorf("rank ID cannot be zero")
	}

	if entry.Score < 0 {
		return fmt.Errorf("score cannot be negative")
	}

	if entry.Rank < 0 {
		return fmt.Errorf("rank cannot be negative")
	}

	if len(entry.Tags) > MaxTagsPerEntry {
		return fmt.Errorf("too many tags: max %d, got %d", MaxTagsPerEntry, len(entry.Tags))
	}

	for _, tag := range entry.Tags {
		if len(tag) > MaxTagLength {
			return fmt.Errorf("tag too long: max %d, got %d", MaxTagLength, len(tag))
		}
	}

	return nil
}

// ValidateBlacklistEntry 验证黑名单条目
func ValidateBlacklistEntry(entry *BlacklistEntry) error {
	if entry == nil {
		return fmt.Errorf("blacklist entry cannot be nil")
	}

	if entry.PlayerID == 0 {
		return fmt.Errorf("player ID cannot be zero")
	}

	if entry.Reason == "" {
		return fmt.Errorf("reason cannot be empty")
	}

	if len(entry.Reason) > 500 {
		return fmt.Errorf("reason too long: max 500, got %d", len(entry.Reason))
	}

	if !entry.IsPermanent && entry.ExpiresAt == nil {
		return fmt.Errorf("temporary blacklist entry must have expiration time")
	}

	if entry.ExpiresAt != nil && entry.ExpiresAt.Before(entry.AddedAt) {
		return fmt.Errorf("expiration time cannot be before added time")
	}

	return nil
}
