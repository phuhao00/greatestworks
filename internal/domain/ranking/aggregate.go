package ranking

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// RankingAggregate 排行榜聚合根
type RankingAggregate struct {
	// 基础信息
	ID          string      `json:"id" bson:"_id"`
	RankID      uint32      `json:"rank_id" bson:"rank_id"`
	Name        string      `json:"name" bson:"name"`
	Description string      `json:"description" bson:"description"`
	RankType    RankType    `json:"rank_type" bson:"rank_type"`
	Category    RankCategory `json:"category" bson:"category"`
	
	// 排序配置
	SortType     SortType `json:"sort_type" bson:"sort_type"`
	TimeBitLen   uint32   `json:"time_bit_len" bson:"time_bit_len"`
	TimeUnit     int64    `json:"time_unit" bson:"time_unit"`
	MaxSize      int64    `json:"max_size" bson:"max_size"`
	
	// 时间配置
	StartTime    int64     `json:"start_time" bson:"start_time"`
	EndTime      int64     `json:"end_time" bson:"end_time"`
	Period       RankPeriod `json:"period" bson:"period"`
	ResetTime    *time.Time `json:"reset_time,omitempty" bson:"reset_time,omitempty"`
	
	// 状态信息
	Status       RankStatus `json:"status" bson:"status"`
	IsActive     bool       `json:"is_active" bson:"is_active"`
	LastUpdated  time.Time  `json:"last_updated" bson:"last_updated"`
	Version      int64      `json:"version" bson:"version"`
	
	// 排行数据
	Entries      []*RankEntry `json:"entries" bson:"entries"`
	Blacklist    *Blacklist   `json:"blacklist" bson:"blacklist"`
	
	// 统计信息
	TotalPlayers    int64   `json:"total_players" bson:"total_players"`
	ActiveEntries   int64   `json:"active_entries" bson:"active_entries"`
	AverageScore    float64 `json:"average_score" bson:"average_score"`
	TopScore        int64   `json:"top_score" bson:"top_score"`
	LastScoreUpdate time.Time `json:"last_score_update" bson:"last_score_update"`
	
	// 奖励配置
	RewardConfig *RankRewardConfig `json:"reward_config,omitempty" bson:"reward_config,omitempty"`
	
	// 缓存配置
	CacheConfig  *RankCacheConfig `json:"cache_config,omitempty" bson:"cache_config,omitempty"`
	
	// 创建和更新时间
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
	
	// 内部状态
	mutex        sync.RWMutex `json:"-" bson:"-"`
	dirty        bool         `json:"-" bson:"-"`
	events       []RankingEvent `json:"-" bson:"-"`
}

// NewRankingAggregate 创建新的排行榜聚合
func NewRankingAggregate(rankID uint32, name string, rankType RankType, category RankCategory) *RankingAggregate {
	now := time.Now()
	return &RankingAggregate{
		ID:           generateRankingID(rankID),
		RankID:       rankID,
		Name:         name,
		RankType:     rankType,
		Category:     category,
		SortType:     SortTypeDescending,
		TimeBitLen:   DefaultTimeBitLen,
		TimeUnit:     DefaultTimeUnit,
		MaxSize:      DefaultMaxSize,
		StartTime:    now.Unix(),
		EndTime:      0, // 永久排行榜
		Period:       RankPeriodPermanent,
		Status:       RankStatusActive,
		IsActive:     true,
		LastUpdated:  now,
		Version:      1,
		Entries:      make([]*RankEntry, 0),
		Blacklist:    NewBlacklist(rankID),
		TotalPlayers: 0,
		ActiveEntries: 0,
		AverageScore: 0.0,
		TopScore:     0,
		LastScoreUpdate: now,
		CreatedAt:    now,
		UpdatedAt:    now,
		dirty:        true,
		events:       make([]RankingEvent, 0),
	}
}

// UpdateScore 更新玩家分数
func (r *RankingAggregate) UpdateScore(playerID uint64, score int64, metadata map[string]interface{}) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// 检查排行榜状态
	if !r.IsActive || r.Status != RankStatusActive {
		return NewRankingInactiveError(r.RankID)
	}
	
	// 检查黑名单
	if r.Blacklist.IsBlacklisted(playerID) {
		return NewPlayerBlacklistedError(playerID, r.RankID)
	}
	
	// 检查时间范围
	if !r.isInTimeRange() {
		return NewRankingTimeExpiredError(r.RankID, r.StartTime, r.EndTime)
	}
	
	// 计算带时间因子的分数
	timeScore := r.calculateTimeScore(score)
	
	// 查找现有条目
	existingEntry := r.findEntry(playerID)
	if existingEntry != nil {
		// 更新现有条目
		oldScore := existingEntry.Score
		existingEntry.UpdateScore(timeScore, score, metadata)
		
		// 发布分数更新事件
		r.addEvent(NewPlayerScoreUpdatedEvent(r.ID, playerID, oldScore, score, timeScore))
	} else {
		// 创建新条目
		newEntry := NewRankEntry(playerID, timeScore, score, metadata)
		r.Entries = append(r.Entries, newEntry)
		r.TotalPlayers++
		
		// 发布新玩家加入事件
		r.addEvent(NewPlayerJoinedRankingEvent(r.ID, playerID, score, timeScore))
	}
	
	// 重新排序
	r.sortEntries()
	
	// 限制大小
	r.limitSize()
	
	// 更新统计信息
	r.updateStatistics()
	
	// 标记为脏数据
	r.markDirty()
	
	return nil
}

// GetRanking 获取排行榜数据
func (r *RankingAggregate) GetRanking(start, end int64, excludeBlacklisted bool) ([]*RankEntry, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	if start < 0 || end < start {
		return nil, NewInvalidRangeError(start, end)
	}
	
	entries := make([]*RankEntry, 0)
	count := int64(0)
	
	for i, entry := range r.Entries {
		// 跳过黑名单玩家
		if excludeBlacklisted && r.Blacklist.IsBlacklisted(entry.PlayerID) {
			continue
		}
		
		// 检查范围
		if count >= start && count <= end {
			// 创建副本并设置排名
			entryCopy := *entry
			entryCopy.Rank = int64(i + 1)
			entries = append(entries, &entryCopy)
		}
		
		count++
		
		// 达到结束位置
		if count > end {
			break
		}
	}
	
	return entries, nil
}

// GetPlayerRank 获取玩家排名
func (r *RankingAggregate) GetPlayerRank(playerID uint64) (*RankEntry, int64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	for i, entry := range r.Entries {
		if entry.PlayerID == playerID {
			// 创建副本并设置排名
			entryCopy := *entry
			entryCopy.Rank = int64(i + 1)
			return &entryCopy, int64(i + 1), nil
		}
	}
	
	return nil, -1, NewPlayerNotInRankingError(playerID, r.RankID)
}

// AddToBlacklist 添加到黑名单
func (r *RankingAggregate) AddToBlacklist(playerID uint64, reason string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	err := r.Blacklist.AddPlayer(playerID, reason)
	if err != nil {
		return err
	}
	
	// 从排行榜中移除该玩家
	r.removePlayer(playerID)
	
	// 发布黑名单事件
	r.addEvent(NewPlayerBlacklistedEvent(r.ID, playerID, reason))
	
	// 标记为脏数据
	r.markDirty()
	
	return nil
}

// RemoveFromBlacklist 从黑名单移除
func (r *RankingAggregate) RemoveFromBlacklist(playerID uint64) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	err := r.Blacklist.RemovePlayer(playerID)
	if err != nil {
		return err
	}
	
	// 发布黑名单移除事件
	r.addEvent(NewPlayerUnblacklistedEvent(r.ID, playerID))
	
	// 标记为脏数据
	r.markDirty()
	
	return nil
}

// Reset 重置排行榜
func (r *RankingAggregate) Reset() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// 保存重置前的数据用于事件
	oldEntries := make([]*RankEntry, len(r.Entries))
	copy(oldEntries, r.Entries)
	
	// 重置数据
	r.Entries = make([]*RankEntry, 0)
	r.TotalPlayers = 0
	r.ActiveEntries = 0
	r.AverageScore = 0.0
	r.TopScore = 0
	r.LastScoreUpdate = time.Now()
	r.Version++
	
	// 更新重置时间
	now := time.Now()
	r.ResetTime = &now
	r.LastUpdated = now
	r.UpdatedAt = now
	
	// 发布重置事件
	r.addEvent(NewRankingResetEvent(r.ID, len(oldEntries)))
	
	// 标记为脏数据
	r.markDirty()
	
	return nil
}

// SetActive 设置排行榜激活状态
func (r *RankingAggregate) SetActive(active bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	oldStatus := r.IsActive
	r.IsActive = active
	
	if active {
		r.Status = RankStatusActive
	} else {
		r.Status = RankStatusInactive
	}
	
	r.LastUpdated = time.Now()
	r.UpdatedAt = time.Now()
	
	// 发布状态变更事件
	if oldStatus != active {
		r.addEvent(NewRankingStatusChangedEvent(r.ID, oldStatus, active))
	}
	
	// 标记为脏数据
	r.markDirty()
}

// SetTimeRange 设置时间范围
func (r *RankingAggregate) SetTimeRange(startTime, endTime int64) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if startTime >= endTime && endTime != 0 {
		return NewInvalidTimeRangeError(startTime, endTime)
	}
	
	r.StartTime = startTime
	r.EndTime = endTime
	r.LastUpdated = time.Now()
	r.UpdatedAt = time.Now()
	
	// 标记为脏数据
	r.markDirty()
	
	return nil
}

// SetRewardConfig 设置奖励配置
func (r *RankingAggregate) SetRewardConfig(config *RankRewardConfig) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.RewardConfig = config
	r.LastUpdated = time.Now()
	r.UpdatedAt = time.Now()
	
	// 标记为脏数据
	r.markDirty()
}

// SetCacheConfig 设置缓存配置
func (r *RankingAggregate) SetCacheConfig(config *RankCacheConfig) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.CacheConfig = config
	r.LastUpdated = time.Now()
	r.UpdatedAt = time.Now()
	
	// 标记为脏数据
	r.markDirty()
}

// GetTopPlayers 获取前N名玩家
func (r *RankingAggregate) GetTopPlayers(count int) []*RankEntry {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	if count <= 0 {
		return []*RankEntry{}
	}
	
	if count > len(r.Entries) {
		count = len(r.Entries)
	}
	
	topPlayers := make([]*RankEntry, count)
	for i := 0; i < count; i++ {
		// 创建副本并设置排名
		entryCopy := *r.Entries[i]
		entryCopy.Rank = int64(i + 1)
		topPlayers[i] = &entryCopy
	}
	
	return topPlayers
}

// GetStatistics 获取统计信息
func (r *RankingAggregate) GetStatistics() *RankingStatistics {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	return &RankingStatistics{
		RankID:          r.RankID,
		TotalPlayers:    r.TotalPlayers,
		ActiveEntries:   r.ActiveEntries,
		AverageScore:    r.AverageScore,
		TopScore:        r.TopScore,
		BlacklistCount:  int64(len(r.Blacklist.Players)),
		LastUpdated:     r.LastUpdated,
		LastScoreUpdate: r.LastScoreUpdate,
	}
}

// GetEvents 获取领域事件
func (r *RankingAggregate) GetEvents() []RankingEvent {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	events := make([]RankingEvent, len(r.events))
	copy(events, r.events)
	return events
}

// ClearEvents 清除领域事件
func (r *RankingAggregate) ClearEvents() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.events = make([]RankingEvent, 0)
}

// IsDirty 检查是否有未保存的更改
func (r *RankingAggregate) IsDirty() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	return r.dirty
}

// MarkClean 标记为已保存
func (r *RankingAggregate) MarkClean() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.dirty = false
}

// 私有方法

// calculateTimeScore 计算带时间因子的分数
func (r *RankingAggregate) calculateTimeScore(score int64) int64 {
	nowTime := time.Now().Unix()
	var timeFactor int64
	
	if r.SortType == SortTypeDescending {
		timeFactor = (nowTime - r.StartTime) / r.TimeUnit
	} else {
		timeFactor = (r.EndTime - nowTime) / r.TimeUnit
	}
	
	timeScore := (score << r.TimeBitLen) | timeFactor
	return timeScore
}

// getRealScore 获取真实分数
func (r *RankingAggregate) getRealScore(timeScore int64) int64 {
	return timeScore >> r.TimeBitLen
}

// getRealScoreTime 获取分数设置时间
func (r *RankingAggregate) getRealScoreTime(timeScore int64) int64 {
	timeFactor := timeScore & ((1 << r.TimeBitLen) - 1)
	var realTime int64
	
	if r.SortType == SortTypeDescending {
		realTime = (timeFactor * r.TimeUnit) + r.StartTime
	} else {
		realTime = r.EndTime - (timeFactor * r.TimeUnit)
	}
	
	return realTime
}

// findEntry 查找玩家条目
func (r *RankingAggregate) findEntry(playerID uint64) *RankEntry {
	for _, entry := range r.Entries {
		if entry.PlayerID == playerID {
			return entry
		}
	}
	return nil
}

// removePlayer 移除玩家
func (r *RankingAggregate) removePlayer(playerID uint64) {
	for i, entry := range r.Entries {
		if entry.PlayerID == playerID {
			r.Entries = append(r.Entries[:i], r.Entries[i+1:]...)
			r.TotalPlayers--
			break
		}
	}
}

// sortEntries 排序条目
func (r *RankingAggregate) sortEntries() {
	if r.SortType == SortTypeDescending {
		sort.Slice(r.Entries, func(i, j int) bool {
			return r.Entries[i].TimeScore > r.Entries[j].TimeScore
		})
	} else {
		sort.Slice(r.Entries, func(i, j int) bool {
			return r.Entries[i].TimeScore < r.Entries[j].TimeScore
		})
	}
}

// limitSize 限制大小
func (r *RankingAggregate) limitSize() {
	if int64(len(r.Entries)) > r.MaxSize {
		r.Entries = r.Entries[:r.MaxSize]
		r.TotalPlayers = r.MaxSize
	}
}

// updateStatistics 更新统计信息
func (r *RankingAggregate) updateStatistics() {
	r.ActiveEntries = int64(len(r.Entries))
	
	if len(r.Entries) > 0 {
		// 更新最高分
		r.TopScore = r.getRealScore(r.Entries[0].TimeScore)
		
		// 计算平均分
		totalScore := int64(0)
		for _, entry := range r.Entries {
			totalScore += r.getRealScore(entry.TimeScore)
		}
		r.AverageScore = float64(totalScore) / float64(len(r.Entries))
	} else {
		r.TopScore = 0
		r.AverageScore = 0.0
	}
	
	r.LastScoreUpdate = time.Now()
}

// isInTimeRange 检查是否在时间范围内
func (r *RankingAggregate) isInTimeRange() bool {
	if r.EndTime == 0 {
		return true // 永久排行榜
	}
	
	now := time.Now().Unix()
	return now >= r.StartTime && now <= r.EndTime
}

// addEvent 添加领域事件
func (r *RankingAggregate) addEvent(event RankingEvent) {
	r.events = append(r.events, event)
}

// markDirty 标记为脏数据
func (r *RankingAggregate) markDirty() {
	r.dirty = true
	r.LastUpdated = time.Now()
	r.UpdatedAt = time.Now()
	r.Version++
}

// 辅助函数

// generateRankingID 生成排行榜ID
func generateRankingID(rankID uint32) string {
	return fmt.Sprintf("ranking_%d", rankID)
}

// 常量定义

const (
	// 默认配置
	DefaultTimeBitLen = 24
	DefaultTimeUnit   = 60
	DefaultMaxSize    = 5000
	
	// 排行榜限制
	MaxRankingSize    = 10000
	MinRankingSize    = 10
	MaxNameLength     = 100
	MaxDescLength     = 500
)

// 验证方法

// Validate 验证排行榜聚合
func (r *RankingAggregate) Validate() error {
	if r.RankID == 0 {
		return NewRankingValidationError("rank_id", r.RankID, "rank_id cannot be zero", "required")
	}
	
	if r.Name == "" {
		return NewRankingValidationError("name", r.Name, "name cannot be empty", "required")
	}
	
	if len(r.Name) > MaxNameLength {
		return NewRankingValidationError("name", r.Name, fmt.Sprintf("name length cannot exceed %d", MaxNameLength), "max_length")
	}
	
	if len(r.Description) > MaxDescLength {
		return NewRankingValidationError("description", r.Description, fmt.Sprintf("description length cannot exceed %d", MaxDescLength), "max_length")
	}
	
	if r.MaxSize < MinRankingSize || r.MaxSize > MaxRankingSize {
		return NewRankingValidationError("max_size", r.MaxSize, fmt.Sprintf("max_size must be between %d and %d", MinRankingSize, MaxRankingSize), "range")
	}
	
	if r.TimeBitLen == 0 || r.TimeBitLen > 32 {
		return NewRankingValidationError("time_bit_len", r.TimeBitLen, "time_bit_len must be between 1 and 32", "range")
	}
	
	if r.TimeUnit <= 0 {
		return NewRankingValidationError("time_unit", r.TimeUnit, "time_unit must be positive", "positive")
	}
	
	if r.EndTime != 0 && r.StartTime >= r.EndTime {
		return NewRankingValidationError("time_range", map[string]int64{"start": r.StartTime, "end": r.EndTime}, "start_time must be less than end_time", "time_range")
	}
	
	return nil
}

// Clone 克隆排行榜聚合
func (r *RankingAggregate) Clone() *RankingAggregate {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	clone := &RankingAggregate{
		ID:              r.ID,
		RankID:          r.RankID,
		Name:            r.Name,
		Description:     r.Description,
		RankType:        r.RankType,
		Category:        r.Category,
		SortType:        r.SortType,
		TimeBitLen:      r.TimeBitLen,
		TimeUnit:        r.TimeUnit,
		MaxSize:         r.MaxSize,
		StartTime:       r.StartTime,
		EndTime:         r.EndTime,
		Period:          r.Period,
		Status:          r.Status,
		IsActive:        r.IsActive,
		LastUpdated:     r.LastUpdated,
		Version:         r.Version,
		TotalPlayers:    r.TotalPlayers,
		ActiveEntries:   r.ActiveEntries,
		AverageScore:    r.AverageScore,
		TopScore:        r.TopScore,
		LastScoreUpdate: r.LastScoreUpdate,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
		dirty:           r.dirty,
	}
	
	// 深拷贝重置时间
	if r.ResetTime != nil {
		resetTime := *r.ResetTime
		clone.ResetTime = &resetTime
	}
	
	// 深拷贝条目
	clone.Entries = make([]*RankEntry, len(r.Entries))
	for i, entry := range r.Entries {
		entryCopy := *entry
		clone.Entries[i] = &entryCopy
	}
	
	// 深拷贝黑名单
	if r.Blacklist != nil {
		clone.Blacklist = r.Blacklist.Clone()
	}
	
	// 深拷贝奖励配置
	if r.RewardConfig != nil {
		rewardConfig := *r.RewardConfig
		clone.RewardConfig = &rewardConfig
	}
	
	// 深拷贝缓存配置
	if r.CacheConfig != nil {
		cacheConfig := *r.CacheConfig
		clone.CacheConfig = &cacheConfig
	}
	
	// 深拷贝事件
	clone.events = make([]RankingEvent, len(r.events))
	copy(clone.events, r.events)
	
	return clone
}