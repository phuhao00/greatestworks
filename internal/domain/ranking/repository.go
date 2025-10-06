package ranking

import (
	"time"
)

// RankingRepository 排行榜仓储接口
type RankingRepository interface {
	// 基础CRUD操作
	Save(ranking *RankingAggregate) error
	FindByID(id string) (*RankingAggregate, error)
	FindByRankID(rankID uint32) (*RankingAggregate, error)
	Update(ranking *RankingAggregate) error
	Delete(id string) error

	// 查询操作
	FindByType(rankType RankType) ([]*RankingAggregate, error)
	FindByCategory(category RankCategory) ([]*RankingAggregate, error)
	FindByStatus(status RankStatus) ([]*RankingAggregate, error)
	FindByPeriod(period RankPeriod) ([]*RankingAggregate, error)
	FindActive() ([]*RankingAggregate, error)
	FindExpired() ([]*RankingAggregate, error)

	// 分页查询
	FindWithPagination(query *RankingQuery) (*RankingPageResult, error)
	FindByPlayerID(playerID uint64) ([]*RankingAggregate, error)

	// 统计操作
	Count() (int64, error)
	CountByType(rankType RankType) (int64, error)
	CountByCategory(category RankCategory) (int64, error)
	CountByStatus(status RankStatus) (int64, error)

	// 批量操作
	SaveBatch(rankings []*RankingAggregate) error
	UpdateBatch(rankings []*RankingAggregate) error
	DeleteBatch(ids []string) error

	// 高级查询
	FindByTimeRange(startTime, endTime time.Time) ([]*RankingAggregate, error)
	FindByScoreRange(minScore, maxScore int64) ([]*RankingAggregate, error)
	FindTopRankings(limit int) ([]*RankingAggregate, error)
	FindRecentlyUpdated(duration time.Duration) ([]*RankingAggregate, error)

	// 搜索操作
	Search(keyword string, filters map[string]interface{}) ([]*RankingAggregate, error)
	FindSimilar(ranking *RankingAggregate, limit int) ([]*RankingAggregate, error)
}

// RankEntryRepository 排行榜条目仓储接口
type RankEntryRepository interface {
	// 基础CRUD操作
	Save(entry *RankEntry) error
	FindByID(id string) (*RankEntry, error)
	FindByPlayerAndRank(playerID uint64, rankID uint32) (*RankEntry, error)
	Update(entry *RankEntry) error
	Delete(id string) error

	// 查询操作
	FindByRankID(rankID uint32) ([]*RankEntry, error)
	FindByPlayerID(playerID uint64) ([]*RankEntry, error)
	FindByRankRange(rankID uint32, startRank, endRank int64) ([]*RankEntry, error)
	FindByScoreRange(rankID uint32, minScore, maxScore int64) ([]*RankEntry, error)

	// 分页查询
	FindWithPagination(query *RankEntryQuery) (*RankEntryPageResult, error)

	// 排序查询
	FindTopEntries(rankID uint32, limit int) ([]*RankEntry, error)
	FindBottomEntries(rankID uint32, limit int) ([]*RankEntry, error)
	FindAroundPlayer(rankID uint32, playerID uint64, range_ int) ([]*RankEntry, error)

	// 统计操作
	Count() (int64, error)
	CountByRankID(rankID uint32) (int64, error)
	CountByPlayerID(playerID uint64) (int64, error)
	CountActive(rankID uint32) (int64, error)

	// 批量操作
	SaveBatch(entries []*RankEntry) error
	UpdateBatch(entries []*RankEntry) error
	DeleteBatch(ids []string) error
	DeleteByRankID(rankID uint32) error
	DeleteByPlayerID(playerID uint64) error

	// 高级查询
	FindByLastActive(rankID uint32, since time.Time) ([]*RankEntry, error)
	FindByConsecutiveDays(rankID uint32, minDays int32) ([]*RankEntry, error)
	FindByTotalUpdates(rankID uint32, minUpdates int64) ([]*RankEntry, error)
	FindRecentlyImproved(rankID uint32, duration time.Duration) ([]*RankEntry, error)

	// 历史数据
	FindScoreHistory(entryID string, limit int) ([]*ScoreHistoryEntry, error)
	FindRewardHistory(entryID string, limit int) ([]*RankRewardEarned, error)

	// 聚合查询
	GetAverageScore(rankID uint32) (float64, error)
	GetTopScore(rankID uint32) (int64, error)
	GetScoreDistribution(rankID uint32, buckets int) (map[string]int64, error)
}

// BlacklistRepository 黑名单仓储接口
type BlacklistRepository interface {
	// 基础CRUD操作
	Save(blacklist *Blacklist) error
	FindByID(id string) (*Blacklist, error)
	FindByRankID(rankID uint32) (*Blacklist, error)
	Update(blacklist *Blacklist) error
	Delete(id string) error

	// 玩家操作
	AddPlayer(rankID uint32, entry *BlacklistEntry) error
	RemovePlayer(rankID uint32, playerID uint64) error
	IsPlayerBlacklisted(rankID uint32, playerID uint64) (bool, error)
	GetBlacklistEntry(rankID uint32, playerID uint64) (*BlacklistEntry, error)

	// 查询操作
	FindBlacklistedPlayers(rankID uint32) ([]*BlacklistEntry, error)
	FindByReason(rankID uint32, reason string) ([]*BlacklistEntry, error)
	FindExpiredEntries(rankID uint32) ([]*BlacklistEntry, error)
	FindPermanentEntries(rankID uint32) ([]*BlacklistEntry, error)

	// 分页查询
	FindWithPagination(query *BlacklistQuery) (*BlacklistPageResult, error)

	// 统计操作
	Count() (int64, error)
	CountByRankID(rankID uint32) (int64, error)
	CountByReason(rankID uint32, reason string) (int64, error)
	CountExpired(rankID uint32) (int64, error)

	// 批量操作
	AddPlayersBatch(rankID uint32, entries []*BlacklistEntry) error
	RemovePlayersBatch(rankID uint32, playerIDs []uint64) error
	CleanupExpired(rankID uint32) (int64, error)

	// 全局操作
	FindPlayerInAllRankings(playerID uint64) ([]*BlacklistEntry, error)
	IsPlayerGloballyBlacklisted(playerID uint64) (bool, error)
}

// RankingStatisticsRepository 排行榜统计仓储接口
type RankingStatisticsRepository interface {
	// 基础操作
	Save(stats *RankingStatistics) error
	FindByRankID(rankID uint32) (*RankingStatistics, error)
	Update(stats *RankingStatistics) error
	Delete(rankID uint32) error

	// 历史统计
	SaveHistorySnapshot(rankID uint32, stats *RankingStatistics) error
	GetHistoryStatistics(rankID uint32, period RankPeriod, points int) ([]*RankingStatistics, error)
	GetStatisticsTrend(rankID uint32, startTime, endTime time.Time) ([]*RankingStatistics, error)

	// 聚合统计
	GetGlobalStatistics() (*GlobalRankingStatistics, error)
	GetCategoryStatistics(category RankCategory) (*CategoryRankingStatistics, error)
	GetTypeStatistics(rankType RankType) (*TypeRankingStatistics, error)

	// 排行榜统计
	GetTopRankingsByPlayers(limit int) ([]*RankingStatistics, error)
	GetTopRankingsByActivity(limit int) ([]*RankingStatistics, error)
	GetMostCompetitiveRankings(limit int) ([]*RankingStatistics, error)

	// 时间范围统计
	GetStatisticsByTimeRange(startTime, endTime time.Time) ([]*RankingStatistics, error)
	GetDailyStatistics(rankID uint32, days int) ([]*DailyRankingStats, error)
	GetWeeklyStatistics(rankID uint32, weeks int) ([]*WeeklyRankingStats, error)
	GetMonthlyStatistics(rankID uint32, months int) ([]*MonthlyRankingStats, error)

	// 比较统计
	CompareRankings(rankID1, rankID2 uint32) (*RankingComparison, error)
	GetRankingPerformance(rankID uint32, period RankPeriod) (*RankingPerformance, error)

	// 清理操作
	CleanupOldStatistics(before time.Time) (int64, error)
	ArchiveStatistics(before time.Time) error
}

// RankingCacheRepository 排行榜缓存仓储接口
type RankingCacheRepository interface {
	// 排行榜缓存
	SetRanking(rankID uint32, ranking *RankingAggregate, ttl time.Duration) error
	GetRanking(rankID uint32) (*RankingAggregate, error)
	DeleteRanking(rankID uint32) error

	// 排行榜数据缓存
	SetRankingData(rankID uint32, start, end int64, entries []*RankEntry, ttl time.Duration) error
	GetRankingData(rankID uint32, start, end int64) ([]*RankEntry, error)
	DeleteRankingData(rankID uint32) error

	// 玩家排名缓存
	SetPlayerRank(rankID uint32, playerID uint64, entry *RankEntry, rank int64, ttl time.Duration) error
	GetPlayerRank(rankID uint32, playerID uint64) (*RankEntry, int64, error)
	DeletePlayerRank(rankID uint32, playerID uint64) error

	// 前N名缓存
	SetTopPlayers(rankID uint32, count int, entries []*RankEntry, ttl time.Duration) error
	GetTopPlayers(rankID uint32, count int) ([]*RankEntry, error)
	DeleteTopPlayers(rankID uint32) error

	// 统计缓存
	SetStatistics(rankID uint32, stats *RankingStatistics, ttl time.Duration) error
	GetStatistics(rankID uint32) (*RankingStatistics, error)
	DeleteStatistics(rankID uint32) error

	// 黑名单缓存
	SetBlacklist(rankID uint32, blacklist *Blacklist, ttl time.Duration) error
	GetBlacklist(rankID uint32) (*Blacklist, error)
	DeleteBlacklist(rankID uint32) error

	// 玩家黑名单状态缓存
	SetPlayerBlacklistStatus(rankID uint32, playerID uint64, isBlacklisted bool, ttl time.Duration) error
	GetPlayerBlacklistStatus(rankID uint32, playerID uint64) (bool, error)
	DeletePlayerBlacklistStatus(rankID uint32, playerID uint64) error

	// 批量操作
	SetBatch(items map[string]interface{}, ttl time.Duration) error
	GetBatch(keys []string) (map[string]interface{}, error)
	DeleteBatch(keys []string) error

	// 缓存管理
	Clear() error
	ClearByPattern(pattern string) error
	Exists(key string) (bool, error)
	SetTTL(key string, ttl time.Duration) error
	GetTTL(key string) (time.Duration, error)
	GetCacheInfo() (*CacheInfo, error)

	// 预热和刷新
	WarmupRanking(rankID uint32) error
	RefreshRanking(rankID uint32) error
	ScheduleRefresh(rankID uint32, interval time.Duration) error
	CancelScheduledRefresh(rankID uint32) error
}

// RankingEventRepository 排行榜事件仓储接口
type RankingEventRepository interface {
	// 事件存储
	Save(event RankingEvent) error
	SaveBatch(events []RankingEvent) error

	// 事件查询
	FindByID(eventID string) (RankingEvent, error)
	FindByAggregateID(aggregateID string) ([]RankingEvent, error)
	FindByEventType(eventType string) ([]RankingEvent, error)
	FindByPlayerID(playerID uint64) ([]RankingEvent, error)
	FindByRankID(rankID uint32) ([]RankingEvent, error)

	// 时间范围查询
	FindByTimeRange(startTime, endTime time.Time) ([]RankingEvent, error)
	FindRecent(limit int) ([]RankingEvent, error)

	// 分页查询
	FindWithPagination(query *RankingEventQuery) (*RankingEventPageResult, error)

	// 事件流
	GetEventStream(aggregateID string, fromVersion int) ([]RankingEvent, error)
	GetEventStreamByType(eventType string, limit int) ([]RankingEvent, error)

	// 快照管理
	SaveSnapshot(aggregateID string, version int, data interface{}) error
	GetSnapshot(aggregateID string) (interface{}, int, error)
	DeleteSnapshot(aggregateID string) error

	// 事件清理
	DeleteByAggregateID(aggregateID string) error
	DeleteByTimeRange(startTime, endTime time.Time) error
	ArchiveEvents(before time.Time) error
	CleanupEvents(before time.Time) error

	// 统计
	CountEvents() (int64, error)
	CountByEventType(eventType string) (int64, error)
	CountByAggregateID(aggregateID string) (int64, error)
}

// RankingSearchRepository 排行榜搜索仓储接口
type RankingSearchRepository interface {
	// 全文搜索
	SearchRankings(query string, filters map[string]interface{}) ([]*RankingAggregate, error)
	SearchEntries(query string, filters map[string]interface{}) ([]*RankEntry, error)
	SearchPlayers(query string, filters map[string]interface{}) ([]*RankEntry, error)

	// 智能推荐
	RecommendRankings(playerID uint64, limit int) ([]*RankingAggregate, error)
	RecommendCompetitors(rankID uint32, playerID uint64, limit int) ([]*RankEntry, error)
	RecommendSimilarPlayers(rankID uint32, playerID uint64, limit int) ([]*RankEntry, error)

	// 相似度搜索
	FindSimilarRankings(rankID uint32, limit int) ([]*RankingAggregate, error)
	FindSimilarPlayers(playerID uint64, limit int) ([]*RankEntry, error)

	// 趋势分析
	AnalyzeTrends(rankID uint32, period RankPeriod) (*RankingTrend, error)
	PredictRankings(rankID uint32, days int) (*RankingPrediction, error)

	// 索引管理
	RebuildIndex() error
	UpdateIndex(entityType string, entityID string, data interface{}) error
	DeleteFromIndex(entityType string, entityID string) error
	OptimizeIndex() error
}

// 查询条件结构体

// RankEntryQuery 排行榜条目查询条件
type RankEntryQuery struct {
	RankID           *uint32    `json:"rank_id,omitempty"`
	PlayerID         *uint64    `json:"player_id,omitempty"`
	PlayerName       string     `json:"player_name,omitempty"`
	MinScore         *int64     `json:"min_score,omitempty"`
	MaxScore         *int64     `json:"max_score,omitempty"`
	MinRank          *int64     `json:"min_rank,omitempty"`
	MaxRank          *int64     `json:"max_rank,omitempty"`
	IsActive         *bool      `json:"is_active,omitempty"`
	MinLevel         *uint32    `json:"min_level,omitempty"`
	MaxLevel         *uint32    `json:"max_level,omitempty"`
	LastActiveAfter  *time.Time `json:"last_active_after,omitempty"`
	LastActiveBefore *time.Time `json:"last_active_before,omitempty"`
	CreatedAfter     *time.Time `json:"created_after,omitempty"`
	CreatedBefore    *time.Time `json:"created_before,omitempty"`
	UpdatedAfter     *time.Time `json:"updated_after,omitempty"`
	UpdatedBefore    *time.Time `json:"updated_before,omitempty"`
	Tags             []string   `json:"tags,omitempty"`
	OrderBy          string     `json:"order_by,omitempty"`
	OrderDesc        bool       `json:"order_desc,omitempty"`
	Offset           int        `json:"offset,omitempty"`
	Limit            int        `json:"limit,omitempty"`
}

// BlacklistQuery 黑名单查询条件
type BlacklistQuery struct {
	RankID        *uint32    `json:"rank_id,omitempty"`
	PlayerID      *uint64    `json:"player_id,omitempty"`
	Reason        string     `json:"reason,omitempty"`
	IsPermanent   *bool      `json:"is_permanent,omitempty"`
	IsExpired     *bool      `json:"is_expired,omitempty"`
	AddedBy       string     `json:"added_by,omitempty"`
	AddedAfter    *time.Time `json:"added_after,omitempty"`
	AddedBefore   *time.Time `json:"added_before,omitempty"`
	ExpiresAfter  *time.Time `json:"expires_after,omitempty"`
	ExpiresBefore *time.Time `json:"expires_before,omitempty"`
	OrderBy       string     `json:"order_by,omitempty"`
	OrderDesc     bool       `json:"order_desc,omitempty"`
	Offset        int        `json:"offset,omitempty"`
	Limit         int        `json:"limit,omitempty"`
}

// RankingEventQuery 排行榜事件查询条件
type RankingEventQuery struct {
	EventID         string     `json:"event_id,omitempty"`
	EventType       string     `json:"event_type,omitempty"`
	AggregateID     string     `json:"aggregate_id,omitempty"`
	PlayerID        *uint64    `json:"player_id,omitempty"`
	RankID          *uint32    `json:"rank_id,omitempty"`
	MinVersion      *int       `json:"min_version,omitempty"`
	MaxVersion      *int       `json:"max_version,omitempty"`
	TimestampAfter  *time.Time `json:"timestamp_after,omitempty"`
	TimestampBefore *time.Time `json:"timestamp_before,omitempty"`
	OrderBy         string     `json:"order_by,omitempty"`
	OrderDesc       bool       `json:"order_desc,omitempty"`
	Offset          int        `json:"offset,omitempty"`
	Limit           int        `json:"limit,omitempty"`
}

// 分页结果结构体

// RankingPageResult 排行榜分页结果
type RankingPageResult struct {
	Items   []*RankingAggregate `json:"items"`
	Total   int64               `json:"total"`
	Offset  int                 `json:"offset"`
	Limit   int                 `json:"limit"`
	HasMore bool                `json:"has_more"`
}

// RankEntryPageResult 排行榜条目分页结果
type RankEntryPageResult struct {
	Items   []*RankEntry `json:"items"`
	Total   int64        `json:"total"`
	Offset  int          `json:"offset"`
	Limit   int          `json:"limit"`
	HasMore bool         `json:"has_more"`
}

// BlacklistPageResult 黑名单分页结果
type BlacklistPageResult struct {
	Items   []*BlacklistEntry `json:"items"`
	Total   int64             `json:"total"`
	Offset  int               `json:"offset"`
	Limit   int               `json:"limit"`
	HasMore bool              `json:"has_more"`
}

// RankingEventPageResult 排行榜事件分页结果
type RankingEventPageResult struct {
	Items   []RankingEvent `json:"items"`
	Total   int64          `json:"total"`
	Offset  int            `json:"offset"`
	Limit   int            `json:"limit"`
	HasMore bool           `json:"has_more"`
}

// 统计数据结构体

// GlobalRankingStatistics 全局排行榜统计
type GlobalRankingStatistics struct {
	TotalRankings            int64                  `json:"total_rankings"`
	ActiveRankings           int64                  `json:"active_rankings"`
	TotalPlayers             int64                  `json:"total_players"`
	TotalEntries             int64                  `json:"total_entries"`
	AveragePlayersPerRanking float64                `json:"average_players_per_ranking"`
	CategoryDistribution     map[RankCategory]int64 `json:"category_distribution"`
	TypeDistribution         map[RankType]int64     `json:"type_distribution"`
	StatusDistribution       map[RankStatus]int64   `json:"status_distribution"`
	MostPopularCategory      RankCategory           `json:"most_popular_category"`
	MostPopularType          RankType               `json:"most_popular_type"`
	TotalBlacklisted         int64                  `json:"total_blacklisted"`
	LastUpdated              time.Time              `json:"last_updated"`
}

// CategoryRankingStatistics 分类排行榜统计
type CategoryRankingStatistics struct {
	Category                 RankCategory `json:"category"`
	TotalRankings            int64        `json:"total_rankings"`
	ActiveRankings           int64        `json:"active_rankings"`
	TotalPlayers             int64        `json:"total_players"`
	AveragePlayersPerRanking float64      `json:"average_players_per_ranking"`
	MostActiveRanking        uint32       `json:"most_active_ranking"`
	HighestScore             int64        `json:"highest_score"`
	AverageScore             float64      `json:"average_score"`
	LastUpdated              time.Time    `json:"last_updated"`
}

// TypeRankingStatistics 类型排行榜统计
type TypeRankingStatistics struct {
	RankType                 RankType  `json:"rank_type"`
	TotalRankings            int64     `json:"total_rankings"`
	ActiveRankings           int64     `json:"active_rankings"`
	TotalPlayers             int64     `json:"total_players"`
	AveragePlayersPerRanking float64   `json:"average_players_per_ranking"`
	MostCompetitiveRanking   uint32    `json:"most_competitive_ranking"`
	HighestScore             int64     `json:"highest_score"`
	AverageScore             float64   `json:"average_score"`
	LastUpdated              time.Time `json:"last_updated"`
}

// DailyRankingStats 日排行榜统计
type DailyRankingStats struct {
	Date             time.Time `json:"date"`
	RankID           uint32    `json:"rank_id"`
	TotalPlayers     int64     `json:"total_players"`
	ActiveEntries    int64     `json:"active_entries"`
	NewPlayers       int64     `json:"new_players"`
	TopScore         int64     `json:"top_score"`
	AverageScore     float64   `json:"average_score"`
	ScoreUpdates     int64     `json:"score_updates"`
	RankChanges      int64     `json:"rank_changes"`
	BlacklistChanges int64     `json:"blacklist_changes"`
}

// WeeklyRankingStats 周排行榜统计
type WeeklyRankingStats struct {
	WeekStart       time.Time `json:"week_start"`
	WeekEnd         time.Time `json:"week_end"`
	RankID          uint32    `json:"rank_id"`
	TotalPlayers    int64     `json:"total_players"`
	ActiveEntries   int64     `json:"active_entries"`
	NewPlayers      int64     `json:"new_players"`
	TopScore        int64     `json:"top_score"`
	AverageScore    float64   `json:"average_score"`
	ScoreUpdates    int64     `json:"score_updates"`
	RankChanges     int64     `json:"rank_changes"`
	Competitiveness float64   `json:"competitiveness"`
	GrowthRate      float64   `json:"growth_rate"`
}

// MonthlyRankingStats 月排行榜统计
type MonthlyRankingStats struct {
	Month           time.Time `json:"month"`
	RankID          uint32    `json:"rank_id"`
	TotalPlayers    int64     `json:"total_players"`
	ActiveEntries   int64     `json:"active_entries"`
	NewPlayers      int64     `json:"new_players"`
	TopScore        int64     `json:"top_score"`
	AverageScore    float64   `json:"average_score"`
	ScoreUpdates    int64     `json:"score_updates"`
	RankChanges     int64     `json:"rank_changes"`
	Competitiveness float64   `json:"competitiveness"`
	GrowthRate      float64   `json:"growth_rate"`
	RetentionRate   float64   `json:"retention_rate"`
	ChurnRate       float64   `json:"churn_rate"`
}

// RankingComparison 排行榜比较
type RankingComparison struct {
	RankID1             uint32    `json:"rank_id_1"`
	RankID2             uint32    `json:"rank_id_2"`
	PlayerDiff          int64     `json:"player_diff"`
	ScoreDiff           float64   `json:"score_diff"`
	ActivityDiff        float64   `json:"activity_diff"`
	CompetitivenessDiff float64   `json:"competitiveness_diff"`
	SimilarityScore     float64   `json:"similarity_score"`
	ComparedAt          time.Time `json:"compared_at"`
}

// RankingPerformance 排行榜性能
type RankingPerformance struct {
	RankID          uint32     `json:"rank_id"`
	Period          RankPeriod `json:"period"`
	PlayerGrowth    float64    `json:"player_growth"`
	ScoreGrowth     float64    `json:"score_growth"`
	ActivityLevel   float64    `json:"activity_level"`
	Competitiveness float64    `json:"competitiveness"`
	RetentionRate   float64    `json:"retention_rate"`
	EngagementScore float64    `json:"engagement_score"`
	HealthScore     float64    `json:"health_score"`
	TrendDirection  string     `json:"trend_direction"`
	CalculatedAt    time.Time  `json:"calculated_at"`
}

// RankingPrediction 排行榜预测
type RankingPrediction struct {
	RankID             uint32    `json:"rank_id"`
	PredictionDays     int       `json:"prediction_days"`
	PredictedPlayers   int64     `json:"predicted_players"`
	PredictedTopScore  int64     `json:"predicted_top_score"`
	PredictedAvgScore  float64   `json:"predicted_avg_score"`
	ConfidenceLevel    float64   `json:"confidence_level"`
	PredictionAccuracy float64   `json:"prediction_accuracy"`
	RiskFactors        []string  `json:"risk_factors"`
	Recommendations    []string  `json:"recommendations"`
	CreatedAt          time.Time `json:"created_at"`
	ValidUntil         time.Time `json:"valid_until"`
}

// CacheInfo 缓存信息
type CacheInfo struct {
	TotalKeys     int64         `json:"total_keys"`
	UsedMemory    int64         `json:"used_memory"`
	HitRate       float64       `json:"hit_rate"`
	MissRate      float64       `json:"miss_rate"`
	EvictionCount int64         `json:"eviction_count"`
	ExpiredCount  int64         `json:"expired_count"`
	AverageTTL    time.Duration `json:"average_ttl"`
	OldestKey     string        `json:"oldest_key"`
	NewestKey     string        `json:"newest_key"`
	LastCleanup   time.Time     `json:"last_cleanup"`
	CacheHealth   string        `json:"cache_health"`
}

// 仓储工厂接口

// RankingRepositoryFactory 排行榜仓储工厂接口
type RankingRepositoryFactory interface {
	// 创建仓储实例
	CreateRankingRepository() RankingRepository
	CreateRankEntryRepository() RankEntryRepository
	CreateBlacklistRepository() BlacklistRepository
	CreateStatisticsRepository() RankingStatisticsRepository
	CreateCacheRepository() RankingCacheRepository
	CreateEventRepository() RankingEventRepository
	CreateSearchRepository() RankingSearchRepository

	// 健康检查
	HealthCheck() error

	// 关闭连接
	Close() error
}

// 事务接口

// RankingTransactionRepository 排行榜事务仓储接口
type RankingTransactionRepository interface {
	// 事务管理
	BeginTransaction() (RankingTransaction, error)
	CommitTransaction(tx RankingTransaction) error
	RollbackTransaction(tx RankingTransaction) error

	// 在事务中执行操作
	ExecuteInTransaction(fn func(tx RankingTransaction) error) error
}

// RankingTransaction 排行榜事务接口
type RankingTransaction interface {
	// 排行榜操作
	SaveRanking(ranking *RankingAggregate) error
	UpdateRanking(ranking *RankingAggregate) error
	DeleteRanking(id string) error

	// 条目操作
	SaveEntry(entry *RankEntry) error
	UpdateEntry(entry *RankEntry) error
	DeleteEntry(id string) error

	// 黑名单操作
	SaveBlacklist(blacklist *Blacklist) error
	UpdateBlacklist(blacklist *Blacklist) error
	DeleteBlacklist(id string) error

	// 统计操作
	SaveStatistics(stats *RankingStatistics) error
	UpdateStatistics(stats *RankingStatistics) error

	// 事件操作
	SaveEvent(event RankingEvent) error
	SaveEvents(events []RankingEvent) error

	// 事务状态
	IsActive() bool
	GetID() string
}
