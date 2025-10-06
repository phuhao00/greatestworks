package ranking

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// RankingService 排行榜领域服务
type RankingService struct {
	// 依赖的仓储
	rankingRepo    RankingRepository
	blacklistRepo  BlacklistRepository
	cacheRepo      RankingCacheRepository
	statisticsRepo RankingStatisticsRepository

	// 配置
	config *RankingServiceConfig

	// 内部状态
	mutex          sync.RWMutex
	activeRankings map[uint32]*RankingAggregate
	lastCleanup    time.Time
}

// RankingServiceConfig 排行榜服务配置
type RankingServiceConfig struct {
	// 缓存配置
	EnableCache          bool          `json:"enable_cache"`
	CacheTTL             time.Duration `json:"cache_ttl"`
	CacheRefreshInterval time.Duration `json:"cache_refresh_interval"`

	// 性能配置
	MaxConcurrentUpdates int           `json:"max_concurrent_updates"`
	BatchSize            int           `json:"batch_size"`
	UpdateTimeout        time.Duration `json:"update_timeout"`

	// 清理配置
	CleanupInterval      time.Duration `json:"cleanup_interval"`
	ExpiredDataRetention time.Duration `json:"expired_data_retention"`

	// 统计配置
	EnableStatistics   bool          `json:"enable_statistics"`
	StatisticsInterval time.Duration `json:"statistics_interval"`

	// 奖励配置
	EnableRewards         bool `json:"enable_rewards"`
	AutoDistributeRewards bool `json:"auto_distribute_rewards"`

	// 验证配置
	EnableValidation bool `json:"enable_validation"`
	StrictMode       bool `json:"strict_mode"`
}

// NewRankingService 创建排行榜服务
func NewRankingService(
	rankingRepo RankingRepository,
	blacklistRepo BlacklistRepository,
	cacheRepo RankingCacheRepository,
	statisticsRepo RankingStatisticsRepository,
	config *RankingServiceConfig,
) *RankingService {
	if config == nil {
		config = DefaultRankingServiceConfig()
	}

	return &RankingService{
		rankingRepo:    rankingRepo,
		blacklistRepo:  blacklistRepo,
		cacheRepo:      cacheRepo,
		statisticsRepo: statisticsRepo,
		config:         config,
		activeRankings: make(map[uint32]*RankingAggregate),
		lastCleanup:    time.Now(),
	}
}

// CreateRanking 创建排行榜
func (rs *RankingService) CreateRanking(rankID uint32, name string, rankType RankType, category RankCategory, config *RankingCreateConfig) (*RankingAggregate, error) {
	// 验证参数
	if err := rs.validateCreateRankingParams(rankID, name, rankType, category, config); err != nil {
		return nil, err
	}

	// 检查排行榜是否已存在
	existing, _ := rs.rankingRepo.FindByRankID(rankID)
	if existing != nil {
		return nil, NewRankingAlreadyExistsError(rankID)
	}

	// 创建排行榜聚合
	ranking := NewRankingAggregate(rankID, name, rankType, category)

	// 应用配置
	if config != nil {
		rs.applyCreateConfig(ranking, config)
	}

	// 验证排行榜
	if rs.config.EnableValidation {
		if err := ranking.Validate(); err != nil {
			return nil, err
		}
	}

	// 保存到仓储
	if err := rs.rankingRepo.Save(ranking); err != nil {
		return nil, NewRankingSystemError("repository", "failed to save ranking", err)
	}

	// 添加到活跃排行榜
	rs.mutex.Lock()
	rs.activeRankings[rankID] = ranking
	rs.mutex.Unlock()

	// 初始化缓存
	if rs.config.EnableCache {
		rs.initializeRankingCache(ranking)
	}

	// 初始化统计
	if rs.config.EnableStatistics {
		rs.initializeRankingStatistics(ranking)
	}

	return ranking, nil
}

// UpdatePlayerScore 更新玩家分数
func (rs *RankingService) UpdatePlayerScore(rankID uint32, playerID uint64, score int64, metadata map[string]interface{}) (*RankingOperationResult, error) {
	start := time.Now()
	result := NewRankingOperationResult(RankingOperationUpdate, rankID, false)
	defer result.SetDuration(start)

	// 获取排行榜
	ranking, err := rs.getRanking(rankID)
	if err != nil {
		result.SetError(err)
		return result, err
	}

	// 获取旧排名和分数
	oldEntry, oldRank, _ := ranking.GetPlayerRank(playerID)
	var oldScore *int64
	if oldEntry != nil {
		oldScore = &oldEntry.Score
	}

	// 更新分数
	err = ranking.UpdateScore(playerID, score, metadata)
	if err != nil {
		result.SetError(err)
		return result, err
	}

	// 获取新排名
	_, newRank, _ := ranking.GetPlayerRank(playerID)

	// 保存排行榜
	if err := rs.rankingRepo.Update(ranking); err != nil {
		err = NewRankingSystemError("repository", "failed to update ranking", err)
		result.SetError(err)
		return result, err
	}

	// 更新缓存
	if rs.config.EnableCache {
		rs.updateRankingCache(ranking)
	}

	// 更新统计
	if rs.config.EnableStatistics {
		rs.updateRankingStatistics(ranking)
	}

	// 检查奖励
	if rs.config.EnableRewards && rs.config.AutoDistributeRewards {
		rs.checkAndDistributeRewards(ranking, playerID, oldRank, newRank)
	}

	// 设置结果
	result.Success = true
	result.SetPlayerInfo(playerID, &oldRank, &newRank, oldScore, &score)
	result.AffectedCount = 1
	result.SetMessage(fmt.Sprintf("Player %d score updated from %v to %d", playerID, oldScore, score))

	return result, nil
}

// GetRanking 获取排行榜数据
func (rs *RankingService) GetRanking(rankID uint32, start, end int64, filter *RankingFilter) ([]*RankEntry, error) {
	// 验证范围
	if err := ValidateRankingRange(start, end); err != nil {
		return nil, err
	}

	// 尝试从缓存获取
	if rs.config.EnableCache {
		if entries, err := rs.getRankingFromCache(rankID, start, end, filter); err == nil {
			return entries, nil
		}
	}

	// 从仓储获取
	ranking, err := rs.getRanking(rankID)
	if err != nil {
		return nil, err
	}

	// 应用过滤器
	excludeBlacklisted := filter == nil || filter.ExcludeBlacklisted
	entries, err := ranking.GetRanking(start, end, excludeBlacklisted)
	if err != nil {
		return nil, err
	}

	// 应用额外过滤
	if filter != nil {
		entries = rs.applyRankingFilter(entries, filter)
	}

	// 更新缓存
	if rs.config.EnableCache {
		rs.cacheRankingData(rankID, start, end, entries)
	}

	return entries, nil
}

// GetPlayerRank 获取玩家排名
func (rs *RankingService) GetPlayerRank(rankID uint32, playerID uint64) (*RankEntry, int64, error) {
	// 尝试从缓存获取
	if rs.config.EnableCache {
		if entry, rank, err := rs.getPlayerRankFromCache(rankID, playerID); err == nil {
			return entry, rank, nil
		}
	}

	// 从仓储获取
	ranking, err := rs.getRanking(rankID)
	if err != nil {
		return nil, -1, err
	}

	entry, rank, err := ranking.GetPlayerRank(playerID)
	if err != nil {
		return nil, -1, err
	}

	// 更新缓存
	if rs.config.EnableCache {
		rs.cachePlayerRank(rankID, playerID, entry, rank)
	}

	return entry, rank, nil
}

// AddToBlacklist 添加到黑名单
func (rs *RankingService) AddToBlacklist(rankID uint32, playerID uint64, reason string, duration *time.Duration) error {
	// 获取排行榜
	ranking, err := rs.getRanking(rankID)
	if err != nil {
		return err
	}

	// 添加到黑名单
	err = ranking.AddToBlacklist(playerID, reason)
	if err != nil {
		return err
	}

	// 如果是临时黑名单，设置过期时间
	if duration != nil {
		blacklistEntry, exists := ranking.Blacklist.GetBlacklistEntry(playerID)
		if exists {
			blacklistEntry.SetExpiration(time.Now().Add(*duration))
		}
	}

	// 保存排行榜
	if err := rs.rankingRepo.Update(ranking); err != nil {
		return NewRankingSystemError("repository", "failed to update ranking", err)
	}

	// 清除相关缓存
	if rs.config.EnableCache {
		rs.clearPlayerCache(rankID, playerID)
		rs.clearRankingCache(rankID)
	}

	return nil
}

// RemoveFromBlacklist 从黑名单移除
func (rs *RankingService) RemoveFromBlacklist(rankID uint32, playerID uint64) error {
	// 获取排行榜
	ranking, err := rs.getRanking(rankID)
	if err != nil {
		return err
	}

	// 从黑名单移除
	err = ranking.RemoveFromBlacklist(playerID)
	if err != nil {
		return err
	}

	// 保存排行榜
	if err := rs.rankingRepo.Update(ranking); err != nil {
		return NewRankingSystemError("repository", "failed to update ranking", err)
	}

	// 清除相关缓存
	if rs.config.EnableCache {
		rs.clearPlayerCache(rankID, playerID)
		rs.clearRankingCache(rankID)
	}

	return nil
}

// ResetRanking 重置排行榜
func (rs *RankingService) ResetRanking(rankID uint32) (*RankingOperationResult, error) {
	start := time.Now()
	result := NewRankingOperationResult(RankingOperationReset, rankID, false)
	defer result.SetDuration(start)

	// 获取排行榜
	ranking, err := rs.getRanking(rankID)
	if err != nil {
		result.SetError(err)
		return result, err
	}

	// 记录重置前的玩家数量
	oldPlayerCount := ranking.TotalPlayers

	// 重置排行榜
	err = ranking.Reset()
	if err != nil {
		result.SetError(err)
		return result, err
	}

	// 保存排行榜
	if err := rs.rankingRepo.Update(ranking); err != nil {
		err = NewRankingSystemError("repository", "failed to update ranking", err)
		result.SetError(err)
		return result, err
	}

	// 清除缓存
	if rs.config.EnableCache {
		rs.clearRankingCache(rankID)
	}

	// 重置统计
	if rs.config.EnableStatistics {
		rs.resetRankingStatistics(rankID)
	}

	// 设置结果
	result.Success = true
	result.AffectedCount = oldPlayerCount
	result.SetMessage(fmt.Sprintf("Ranking %d reset, %d players removed", rankID, oldPlayerCount))

	return result, nil
}

// GetRankingStatistics 获取排行榜统计
func (rs *RankingService) GetRankingStatistics(rankID uint32) (*RankingStatistics, error) {
	// 尝试从缓存获取
	if rs.config.EnableCache {
		if stats, err := rs.getStatisticsFromCache(rankID); err == nil {
			return stats, nil
		}
	}

	// 从排行榜获取
	ranking, err := rs.getRanking(rankID)
	if err != nil {
		return nil, err
	}

	stats := ranking.GetStatistics()

	// 更新缓存
	if rs.config.EnableCache {
		rs.cacheStatistics(rankID, stats)
	}

	return stats, nil
}

// BatchUpdateScores 批量更新分数
func (rs *RankingService) BatchUpdateScores(rankID uint32, updates []*ScoreUpdate) ([]*RankingOperationResult, error) {
	if len(updates) == 0 {
		return []*RankingOperationResult{}, nil
	}

	// 获取排行榜
	ranking, err := rs.getRanking(rankID)
	if err != nil {
		return nil, err
	}

	results := make([]*RankingOperationResult, len(updates))

	// 批量处理更新
	for i, update := range updates {
		start := time.Now()
		result := NewRankingOperationResult(RankingOperationUpdate, rankID, false)

		// 获取旧排名和分数
		oldEntry, oldRank, _ := ranking.GetPlayerRank(update.PlayerID)
		var oldScore *int64
		if oldEntry != nil {
			oldScore = &oldEntry.Score
		}

		// 更新分数
		err := ranking.UpdateScore(update.PlayerID, update.Score, update.Metadata)
		if err != nil {
			result.SetError(err)
		} else {
			// 获取新排名
			_, newRank, _ := ranking.GetPlayerRank(update.PlayerID)

			result.Success = true
			result.SetPlayerInfo(update.PlayerID, &oldRank, &newRank, oldScore, &update.Score)
			result.AffectedCount = 1
		}

		result.SetDuration(start)
		results[i] = result
	}

	// 保存排行榜
	if err := rs.rankingRepo.Update(ranking); err != nil {
		return results, NewRankingSystemError("repository", "failed to update ranking", err)
	}

	// 更新缓存
	if rs.config.EnableCache {
		rs.updateRankingCache(ranking)
	}

	// 更新统计
	if rs.config.EnableStatistics {
		rs.updateRankingStatistics(ranking)
	}

	return results, nil
}

// GetTopPlayers 获取前N名玩家
func (rs *RankingService) GetTopPlayers(rankID uint32, count int) ([]*RankEntry, error) {
	if count <= 0 {
		return []*RankEntry{}, nil
	}

	// 尝试从缓存获取
	if rs.config.EnableCache {
		if entries, err := rs.getTopPlayersFromCache(rankID, count); err == nil {
			return entries, nil
		}
	}

	// 从排行榜获取
	ranking, err := rs.getRanking(rankID)
	if err != nil {
		return nil, err
	}

	entries := ranking.GetTopPlayers(count)

	// 更新缓存
	if rs.config.EnableCache {
		rs.cacheTopPlayers(rankID, count, entries)
	}

	return entries, nil
}

// CalculateRankingTrend 计算排行榜趋势
func (rs *RankingService) CalculateRankingTrend(rankID uint32, period RankPeriod, points int) (*RankingTrend, error) {
	// 获取历史统计数据
	historyStats, err := rs.statisticsRepo.GetHistoryStatistics(rankID, period, points)
	if err != nil {
		return nil, err
	}

	// 计算趋势数据
	trendData := make([]*RankingTrendPoint, len(historyStats))
	for i, stats := range historyStats {
		trendData[i] = &RankingTrendPoint{
			Timestamp:     stats.LastUpdated,
			PlayerCount:   stats.TotalPlayers,
			AverageScore:  stats.AverageScore,
			TopScore:      stats.TopScore,
			ScoreVariance: rs.calculateScoreVariance(stats),
			NewPlayers:    rs.calculateNewPlayers(stats),
			ActivePlayers: stats.ActiveEntries,
		}
	}

	// 计算增长率和波动性
	growthRate := rs.calculateGrowthRate(trendData)
	volatility := rs.calculateVolatility(trendData)

	// 生成预测
	prediction := rs.generateTrendPrediction(trendData, period)

	trend := &RankingTrend{
		RankID:     rankID,
		Period:     period,
		TrendData:  trendData,
		GrowthRate: growthRate,
		Volatility: volatility,
		Prediction: prediction,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return trend, nil
}

// CleanupExpiredData 清理过期数据
func (rs *RankingService) CleanupExpiredData() error {
	now := time.Now()

	// 检查是否需要清理
	if now.Sub(rs.lastCleanup) < rs.config.CleanupInterval {
		return nil
	}

	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	// 清理过期的黑名单条目
	for _, ranking := range rs.activeRankings {
		rs.cleanupExpiredBlacklistEntries(ranking)
	}

	// 清理过期的缓存
	if rs.config.EnableCache {
		rs.cleanupExpiredCache()
	}

	// 清理过期的统计数据
	if rs.config.EnableStatistics {
		rs.cleanupExpiredStatistics()
	}

	rs.lastCleanup = now
	return nil
}

// 私有方法

// getRanking 获取排行榜
func (rs *RankingService) getRanking(rankID uint32) (*RankingAggregate, error) {
	// 先从内存中获取
	rs.mutex.RLock()
	ranking, exists := rs.activeRankings[rankID]
	rs.mutex.RUnlock()

	if exists {
		return ranking, nil
	}

	// 从仓储加载
	ranking, err := rs.rankingRepo.FindByRankID(rankID)
	if err != nil {
		return nil, err
	}

	if ranking == nil {
		return nil, NewRankingNotFoundError(rankID)
	}

	// 添加到内存
	rs.mutex.Lock()
	rs.activeRankings[rankID] = ranking
	rs.mutex.Unlock()

	return ranking, nil
}

// validateCreateRankingParams 验证创建排行榜参数
func (rs *RankingService) validateCreateRankingParams(rankID uint32, name string, rankType RankType, category RankCategory, config *RankingCreateConfig) error {
	if rankID == 0 {
		return NewRankingValidationError("rank_id", rankID, "rank_id cannot be zero", "required")
	}

	if name == "" {
		return NewRankingValidationError("name", name, "name cannot be empty", "required")
	}

	if !rankType.IsValid() {
		return NewRankingValidationError("rank_type", rankType, "invalid rank type", "enum")
	}

	if !category.IsValid() {
		return NewRankingValidationError("category", category, "invalid category", "enum")
	}

	return nil
}

// applyCreateConfig 应用创建配置
func (rs *RankingService) applyCreateConfig(ranking *RankingAggregate, config *RankingCreateConfig) {
	if config.Description != nil {
		ranking.Description = *config.Description
	}

	if config.SortType != nil {
		ranking.SortType = *config.SortType
	}

	if config.MaxSize != nil {
		ranking.MaxSize = *config.MaxSize
	}

	if config.Period != nil {
		ranking.Period = *config.Period
	}

	if config.StartTime != nil {
		ranking.StartTime = config.StartTime.Unix()
	}

	if config.EndTime != nil {
		ranking.EndTime = config.EndTime.Unix()
	}

	if config.RewardConfig != nil {
		ranking.SetRewardConfig(config.RewardConfig)
	}

	if config.CacheConfig != nil {
		ranking.SetCacheConfig(config.CacheConfig)
	}
}

// applyRankingFilter 应用排行榜过滤器
func (rs *RankingService) applyRankingFilter(entries []*RankEntry, filter *RankingFilter) []*RankEntry {
	if filter == nil {
		return entries
	}

	filteredEntries := make([]*RankEntry, 0, len(entries))

	for _, entry := range entries {
		// 检查分数范围
		if filter.MinScore != nil && entry.Score < *filter.MinScore {
			continue
		}
		if filter.MaxScore != nil && entry.Score > *filter.MaxScore {
			continue
		}

		// 检查玩家ID过滤
		if len(filter.PlayerIDs) > 0 {
			found := false
			for _, playerID := range filter.PlayerIDs {
				if entry.PlayerID == playerID {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// 检查排除玩家ID
		if len(filter.ExcludePlayerIDs) > 0 {
			excluded := false
			for _, playerID := range filter.ExcludePlayerIDs {
				if entry.PlayerID == playerID {
					excluded = true
					break
				}
			}
			if excluded {
				continue
			}
		}

		// 检查活跃状态
		if filter.ExcludeInactive && !entry.IsActive {
			continue
		}

		// 检查时间范围
		if filter.TimeRange != nil {
			if !filter.TimeRange.Contains(entry.LastUpdateTime) {
				continue
			}
		}

		filteredEntries = append(filteredEntries, entry)
	}

	return filteredEntries
}

// 缓存相关方法

func (rs *RankingService) initializeRankingCache(ranking *RankingAggregate) {
	// 实现缓存初始化逻辑
}

func (rs *RankingService) updateRankingCache(ranking *RankingAggregate) {
	// 实现缓存更新逻辑
}

func (rs *RankingService) getRankingFromCache(rankID uint32, start, end int64, filter *RankingFilter) ([]*RankEntry, error) {
	// 实现从缓存获取排行榜数据的逻辑
	return nil, fmt.Errorf("cache miss")
}

func (rs *RankingService) cacheRankingData(rankID uint32, start, end int64, entries []*RankEntry) {
	// 实现缓存排行榜数据的逻辑
}

func (rs *RankingService) getPlayerRankFromCache(rankID uint32, playerID uint64) (*RankEntry, int64, error) {
	// 实现从缓存获取玩家排名的逻辑
	return nil, -1, fmt.Errorf("cache miss")
}

func (rs *RankingService) cachePlayerRank(rankID uint32, playerID uint64, entry *RankEntry, rank int64) {
	// 实现缓存玩家排名的逻辑
}

func (rs *RankingService) clearPlayerCache(rankID uint32, playerID uint64) {
	// 实现清除玩家缓存的逻辑
}

func (rs *RankingService) clearRankingCache(rankID uint32) {
	// 实现清除排行榜缓存的逻辑
}

func (rs *RankingService) getTopPlayersFromCache(rankID uint32, count int) ([]*RankEntry, error) {
	// 实现从缓存获取前N名玩家的逻辑
	return nil, fmt.Errorf("cache miss")
}

func (rs *RankingService) cacheTopPlayers(rankID uint32, count int, entries []*RankEntry) {
	// 实现缓存前N名玩家的逻辑
}

func (rs *RankingService) cleanupExpiredCache() {
	// 实现清理过期缓存的逻辑
}

// 统计相关方法

func (rs *RankingService) initializeRankingStatistics(ranking *RankingAggregate) {
	// 实现统计初始化逻辑
}

func (rs *RankingService) updateRankingStatistics(ranking *RankingAggregate) {
	// 实现统计更新逻辑
}

func (rs *RankingService) resetRankingStatistics(rankID uint32) {
	// 实现统计重置逻辑
}

func (rs *RankingService) getStatisticsFromCache(rankID uint32) (*RankingStatistics, error) {
	// 实现从缓存获取统计的逻辑
	return nil, fmt.Errorf("cache miss")
}

func (rs *RankingService) cacheStatistics(rankID uint32, stats *RankingStatistics) {
	// 实现缓存统计的逻辑
}

func (rs *RankingService) cleanupExpiredStatistics() {
	// 实现清理过期统计的逻辑
}

// 奖励相关方法

func (rs *RankingService) checkAndDistributeRewards(ranking *RankingAggregate, playerID uint64, oldRank, newRank int64) {
	// 实现检查和分发奖励的逻辑
}

// 清理相关方法

func (rs *RankingService) cleanupExpiredBlacklistEntries(ranking *RankingAggregate) {
	// 清理过期的黑名单条目
	expiredPlayers := make([]uint64, 0)

	for playerID, entry := range ranking.Blacklist.Players {
		if entry.IsExpired() {
			expiredPlayers = append(expiredPlayers, playerID)
		}
	}

	// 移除过期的玩家
	for _, playerID := range expiredPlayers {
		ranking.RemoveFromBlacklist(playerID)
	}

	// 如果有变更，保存排行榜
	if len(expiredPlayers) > 0 {
		rs.rankingRepo.Update(ranking)
	}
}

// 趋势计算相关方法

func (rs *RankingService) calculateScoreVariance(stats *RankingStatistics) float64 {
	// 实现分数方差计算逻辑
	return 0.0
}

func (rs *RankingService) calculateNewPlayers(stats *RankingStatistics) int64 {
	// 实现新玩家数量计算逻辑
	return 0
}

func (rs *RankingService) calculateGrowthRate(trendData []*RankingTrendPoint) float64 {
	if len(trendData) < 2 {
		return 0.0
	}

	first := trendData[0]
	last := trendData[len(trendData)-1]

	if first.PlayerCount == 0 {
		return 0.0
	}

	return float64(last.PlayerCount-first.PlayerCount) / float64(first.PlayerCount) * 100
}

func (rs *RankingService) calculateVolatility(trendData []*RankingTrendPoint) float64 {
	if len(trendData) < 2 {
		return 0.0
	}

	// 计算分数变化的标准差
	scores := make([]float64, len(trendData))
	for i, point := range trendData {
		scores[i] = point.AverageScore
	}

	// 计算平均值
	sum := 0.0
	for _, score := range scores {
		sum += score
	}
	mean := sum / float64(len(scores))

	// 计算方差
	variance := 0.0
	for _, score := range scores {
		variance += math.Pow(score-mean, 2)
	}
	variance /= float64(len(scores))

	// 返回标准差
	return math.Sqrt(variance)
}

func (rs *RankingService) generateTrendPrediction(trendData []*RankingTrendPoint, period RankPeriod) *RankingTrendPrediction {
	if len(trendData) < 3 {
		return nil
	}

	// 简单的线性预测
	last := trendData[len(trendData)-1]
	secondLast := trendData[len(trendData)-2]

	playerGrowth := last.PlayerCount - secondLast.PlayerCount
	scoreGrowth := last.AverageScore - secondLast.AverageScore
	topScoreGrowth := last.TopScore - secondLast.TopScore

	prediction := &RankingTrendPrediction{
		PredictedPlayerCount:  last.PlayerCount + playerGrowth,
		PredictedAverageScore: last.AverageScore + scoreGrowth,
		PredictedTopScore:     last.TopScore + topScoreGrowth,
		ConfidenceLevel:       0.7, // 70%置信度
		PredictionTime:        time.Now(),
		ValidUntil:            time.Now().Add(period.GetDuration()),
	}

	return prediction
}

// 辅助结构体

// RankingCreateConfig 排行榜创建配置
type RankingCreateConfig struct {
	Description  *string           `json:"description,omitempty"`
	SortType     *SortType         `json:"sort_type,omitempty"`
	MaxSize      *int64            `json:"max_size,omitempty"`
	Period       *RankPeriod       `json:"period,omitempty"`
	StartTime    *time.Time        `json:"start_time,omitempty"`
	EndTime      *time.Time        `json:"end_time,omitempty"`
	RewardConfig *RankRewardConfig `json:"reward_config,omitempty"`
	CacheConfig  *RankCacheConfig  `json:"cache_config,omitempty"`
}

// ScoreUpdate 分数更新
type ScoreUpdate struct {
	PlayerID uint64                 `json:"player_id"`
	Score    int64                  `json:"score"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// DefaultRankingServiceConfig 默认排行榜服务配置
func DefaultRankingServiceConfig() *RankingServiceConfig {
	return &RankingServiceConfig{
		EnableCache:           true,
		CacheTTL:              30 * time.Minute,
		CacheRefreshInterval:  5 * time.Minute,
		MaxConcurrentUpdates:  100,
		BatchSize:             50,
		UpdateTimeout:         30 * time.Second,
		CleanupInterval:       1 * time.Hour,
		ExpiredDataRetention:  7 * 24 * time.Hour,
		EnableStatistics:      true,
		StatisticsInterval:    10 * time.Minute,
		EnableRewards:         true,
		AutoDistributeRewards: true,
		EnableValidation:      true,
		StrictMode:            false,
	}
}
