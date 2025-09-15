package ranking

import (
	"time"
)

// 排行榜相关事件定义

// RankingEvent 排行榜事件基础接口
type RankingEvent interface {
	GetEventID() string
	GetEventType() string
	GetAggregateID() string
	GetRankID() uint32
	GetTimestamp() time.Time
	GetVersion() int
	GetMetadata() map[string]interface{}
}

// BaseRankingEvent 排行榜事件基础结构
type BaseRankingEvent struct {
	EventID     string                 `json:"event_id"`
	EventType   string                 `json:"event_type"`
	AggregateID string                 `json:"aggregate_id"`
	RankID      uint32                 `json:"rank_id"`
	Timestamp   time.Time              `json:"timestamp"`
	Version     int                    `json:"version"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// GetEventID 获取事件ID
func (e *BaseRankingEvent) GetEventID() string {
	return e.EventID
}

// GetEventType 获取事件类型
func (e *BaseRankingEvent) GetEventType() string {
	return e.EventType
}

// GetAggregateID 获取聚合ID
func (e *BaseRankingEvent) GetAggregateID() string {
	return e.AggregateID
}

// GetRankID 获取排行榜ID
func (e *BaseRankingEvent) GetRankID() uint32 {
	return e.RankID
}

// GetTimestamp 获取时间戳
func (e *BaseRankingEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetVersion 获取版本
func (e *BaseRankingEvent) GetVersion() int {
	return e.Version
}

// GetMetadata 获取元数据
func (e *BaseRankingEvent) GetMetadata() map[string]interface{} {
	return e.Metadata
}

// 排行榜生命周期事件

// RankingCreatedEvent 排行榜创建事件
type RankingCreatedEvent struct {
	BaseRankingEvent
	Name        string        `json:"name"`
	RankType    RankType      `json:"rank_type"`
	Category    RankCategory  `json:"category"`
	SortType    SortType      `json:"sort_type"`
	MaxSize     int64         `json:"max_size"`
	Period      RankPeriod    `json:"period"`
	StartTime   int64         `json:"start_time"`
	EndTime     int64         `json:"end_time"`
	CreatedBy   string        `json:"created_by"`
}

// RankingDeletedEvent 排行榜删除事件
type RankingDeletedEvent struct {
	BaseRankingEvent
	Name          string `json:"name"`
	FinalPlayers  int64  `json:"final_players"`
	FinalTopScore int64  `json:"final_top_score"`
	DeletedBy     string `json:"deleted_by"`
	DeleteReason  string `json:"delete_reason"`
}

// RankingUpdatedEvent 排行榜更新事件
type RankingUpdatedEvent struct {
	BaseRankingEvent
	ChangedFields []string               `json:"changed_fields"`
	OldValues     map[string]interface{} `json:"old_values"`
	NewValues     map[string]interface{} `json:"new_values"`
	UpdatedBy     string                 `json:"updated_by"`
}

// 排行榜状态事件

// RankingStatusChangedEvent 排行榜状态改变事件
type RankingStatusChangedEvent struct {
	BaseRankingEvent
	OldStatus RankStatus `json:"old_status"`
	NewStatus RankStatus `json:"new_status"`
	OldActive bool       `json:"old_active"`
	NewActive bool       `json:"new_active"`
	Reason    string     `json:"reason"`
	ChangedBy string     `json:"changed_by"`
}

// RankingResetEvent 排行榜重置事件
type RankingResetEvent struct {
	BaseRankingEvent
	PreviousPlayerCount int   `json:"previous_player_count"`
	ResetReason         string `json:"reset_reason"`
	ResetBy             string `json:"reset_by"`
	BackupCreated       bool   `json:"backup_created"`
	BackupLocation      string `json:"backup_location,omitempty"`
}

// RankingArchivedEvent 排行榜归档事件
type RankingArchivedEvent struct {
	BaseRankingEvent
	ArchiveReason   string `json:"archive_reason"`
	ArchiveLocation string `json:"archive_location"`
	ArchivedBy      string `json:"archived_by"`
	RetentionPeriod int64  `json:"retention_period"`
}

// 玩家相关事件

// PlayerJoinedRankingEvent 玩家加入排行榜事件
type PlayerJoinedRankingEvent struct {
	BaseRankingEvent
	PlayerID     uint64 `json:"player_id"`
	PlayerName   string `json:"player_name"`
	InitialScore int64  `json:"initial_score"`
	TimeScore    int64  `json:"time_score"`
	InitialRank  int64  `json:"initial_rank"`
	JoinMethod   string `json:"join_method"` // "first_score", "migration", "restore"
}

// PlayerLeftRankingEvent 玩家离开排行榜事件
type PlayerLeftRankingEvent struct {
	BaseRankingEvent
	PlayerID    uint64 `json:"player_id"`
	PlayerName  string `json:"player_name"`
	FinalScore  int64  `json:"final_score"`
	FinalRank   int64  `json:"final_rank"`
	LeaveReason string `json:"leave_reason"` // "removed", "blacklisted", "expired", "reset"
	DaysActive  int32  `json:"days_active"`
}

// PlayerScoreUpdatedEvent 玩家分数更新事件
type PlayerScoreUpdatedEvent struct {
	BaseRankingEvent
	PlayerID      uint64                 `json:"player_id"`
	PlayerName    string                 `json:"player_name"`
	OldScore      int64                  `json:"old_score"`
	NewScore      int64                  `json:"new_score"`
	ScoreDelta    int64                  `json:"score_delta"`
	OldTimeScore  int64                  `json:"old_time_score"`
	NewTimeScore  int64                  `json:"new_time_score"`
	UpdateSource  string                 `json:"update_source"` // "game", "admin", "system", "correction"
	UpdateMetadata map[string]interface{} `json:"update_metadata"`
}

// PlayerRankChangedEvent 玩家排名改变事件
type PlayerRankChangedEvent struct {
	BaseRankingEvent
	PlayerID    uint64 `json:"player_id"`
	PlayerName  string `json:"player_name"`
	OldRank     int64  `json:"old_rank"`
	NewRank     int64  `json:"new_rank"`
	RankDelta   int64  `json:"rank_delta"`
	Score       int64  `json:"score"`
	ChangeType  string `json:"change_type"` // "improvement", "decline", "new_entry"
	TriggerEvent string `json:"trigger_event"` // "score_update", "other_player_update", "reset"
}

// 黑名单相关事件

// PlayerBlacklistedEvent 玩家被加入黑名单事件
type PlayerBlacklistedEvent struct {
	BaseRankingEvent
	PlayerID      uint64     `json:"player_id"`
	PlayerName    string     `json:"player_name"`
	Reason        string     `json:"reason"`
	BlacklistedBy string     `json:"blacklisted_by"`
	IsPermanent   bool       `json:"is_permanent"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	PreviousRank  int64      `json:"previous_rank"`
	PreviousScore int64      `json:"previous_score"`
}

// PlayerUnblacklistedEvent 玩家从黑名单移除事件
type PlayerUnblacklistedEvent struct {
	BaseRankingEvent
	PlayerID        uint64 `json:"player_id"`
	PlayerName      string `json:"player_name"`
	OriginalReason  string `json:"original_reason"`
	UnblacklistedBy string `json:"unblacklisted_by"`
	UnblacklistReason string `json:"unblacklist_reason"`
	BlacklistDuration time.Duration `json:"blacklist_duration"`
	CanRejoin       bool   `json:"can_rejoin"`
}

// BlacklistExpiredEvent 黑名单过期事件
type BlacklistExpiredEvent struct {
	BaseRankingEvent
	PlayerID         uint64        `json:"player_id"`
	PlayerName       string        `json:"player_name"`
	OriginalReason   string        `json:"original_reason"`
	BlacklistDuration time.Duration `json:"blacklist_duration"`
	AutoRemoved      bool          `json:"auto_removed"`
	CanRejoin        bool          `json:"can_rejoin"`
}

// 奖励相关事件

// RankingRewardDistributedEvent 排行榜奖励分发事件
type RankingRewardDistributedEvent struct {
	BaseRankingEvent
	RewardTier      string                 `json:"reward_tier"`
	MinRank         int64                  `json:"min_rank"`
	MaxRank         int64                  `json:"max_rank"`
	RewardCount     int64                  `json:"reward_count"`
	TotalValue      int64                  `json:"total_value"`
	DistributionMethod string              `json:"distribution_method"` // "immediate", "mail", "deferred"
	DistributedBy   string                 `json:"distributed_by"`
	RewardDetails   []RewardDistribution   `json:"reward_details"`
}

// PlayerRewardEarnedEvent 玩家获得奖励事件
type PlayerRewardEarnedEvent struct {
	BaseRankingEvent
	PlayerID      uint64   `json:"player_id"`
	PlayerName    string   `json:"player_name"`
	PlayerRank    int64    `json:"player_rank"`
	PlayerScore   int64    `json:"player_score"`
	RewardTier    string   `json:"reward_tier"`
	RewardType    string   `json:"reward_type"`
	RewardID      string   `json:"reward_id"`
	RewardQuantity int64   `json:"reward_quantity"`
	RewardValue   int64    `json:"reward_value"`
	EarnedAt      time.Time `json:"earned_at"`
	ClaimDeadline *time.Time `json:"claim_deadline,omitempty"`
}

// PlayerRewardClaimedEvent 玩家领取奖励事件
type PlayerRewardClaimedEvent struct {
	BaseRankingEvent
	PlayerID       uint64    `json:"player_id"`
	PlayerName     string    `json:"player_name"`
	RewardID       string    `json:"reward_id"`
	RewardType     string    `json:"reward_type"`
	RewardQuantity int64     `json:"reward_quantity"`
	RewardValue    int64     `json:"reward_value"`
	EarnedAt       time.Time `json:"earned_at"`
	ClaimedAt      time.Time `json:"claimed_at"`
	ClaimDelay     time.Duration `json:"claim_delay"`
}

// 统计相关事件

// RankingStatisticsUpdatedEvent 排行榜统计更新事件
type RankingStatisticsUpdatedEvent struct {
	BaseRankingEvent
	OldStatistics *RankingStatistics `json:"old_statistics"`
	NewStatistics *RankingStatistics `json:"new_statistics"`
	ChangedFields []string           `json:"changed_fields"`
	UpdateTrigger string             `json:"update_trigger"` // "score_update", "player_join", "player_leave", "scheduled"
}

// RankingMilestoneReachedEvent 排行榜里程碑达成事件
type RankingMilestoneReachedEvent struct {
	BaseRankingEvent
	MilestoneType  string    `json:"milestone_type"` // "player_count", "score_threshold", "activity_level"
	MilestoneValue int64     `json:"milestone_value"`
	CurrentValue   int64     `json:"current_value"`
	ReachedAt      time.Time `json:"reached_at"`
	IsFirstTime    bool      `json:"is_first_time"`
	Celebration    bool      `json:"celebration"`
}

// RankingRecordBrokenEvent 排行榜记录被打破事件
type RankingRecordBrokenEvent struct {
	BaseRankingEvent
	RecordType     string `json:"record_type"` // "highest_score", "most_players", "longest_streak"
	PlayerID       uint64 `json:"player_id"`
	PlayerName     string `json:"player_name"`
	OldRecord      int64  `json:"old_record"`
	NewRecord      int64  `json:"new_record"`
	RecordImprovement int64 `json:"record_improvement"`
	PreviousHolder string `json:"previous_holder"`
	RecordDuration time.Duration `json:"record_duration"`
}

// 系统事件

// RankingMaintenanceEvent 排行榜维护事件
type RankingMaintenanceEvent struct {
	BaseRankingEvent
	MaintenanceType string        `json:"maintenance_type"` // "scheduled", "emergency", "optimization"
	StartTime       time.Time     `json:"start_time"`
	EndTime         *time.Time    `json:"end_time,omitempty"`
	Duration        time.Duration `json:"duration"`
	AffectedFeatures []string     `json:"affected_features"`
	MaintenanceBy   string        `json:"maintenance_by"`
	Reason          string        `json:"reason"`
	Impact          string        `json:"impact"` // "none", "limited", "full"
}

// RankingDataMigrationEvent 排行榜数据迁移事件
type RankingDataMigrationEvent struct {
	BaseRankingEvent
	MigrationType    string `json:"migration_type"` // "version_upgrade", "schema_change", "platform_migration"
	FromVersion      string `json:"from_version"`
	ToVersion        string `json:"to_version"`
	MigratedRecords  int64  `json:"migrated_records"`
	FailedRecords    int64  `json:"failed_records"`
	MigrationStatus  string `json:"migration_status"` // "started", "in_progress", "completed", "failed"
	MigrationErrors  []string `json:"migration_errors,omitempty"`
	RollbackPlan     string `json:"rollback_plan"`
}

// RankingPerformanceEvent 排行榜性能事件
type RankingPerformanceEvent struct {
	BaseRankingEvent
	MetricType      string  `json:"metric_type"` // "response_time", "throughput", "error_rate", "memory_usage"
	MetricValue     float64 `json:"metric_value"`
	Threshold       float64 `json:"threshold"`
	Severity        string  `json:"severity"` // "info", "warning", "critical"
	Duration        time.Duration `json:"duration"`
	AffectedOperations []string `json:"affected_operations"`
	ResolutionAction string `json:"resolution_action"`
}

// 缓存相关事件

// RankingCacheUpdatedEvent 排行榜缓存更新事件
type RankingCacheUpdatedEvent struct {
	BaseRankingEvent
	CacheType     string        `json:"cache_type"` // "ranking_data", "player_rank", "statistics", "top_players"
	CacheKey      string        `json:"cache_key"`
	UpdateReason  string        `json:"update_reason"` // "data_change", "expiration", "manual_refresh"
	CacheSize     int64         `json:"cache_size"`
	TTL           time.Duration `json:"ttl"`
	HitRate       float64       `json:"hit_rate"`
}

// RankingCacheEvictedEvent 排行榜缓存驱逐事件
type RankingCacheEvictedEvent struct {
	BaseRankingEvent
	CacheType     string `json:"cache_type"`
	CacheKey      string `json:"cache_key"`
	EvictionReason string `json:"eviction_reason"` // "memory_pressure", "ttl_expired", "manual_clear", "size_limit"
	CacheAge      time.Duration `json:"cache_age"`
	AccessCount   int64  `json:"access_count"`
	LastAccessed  time.Time `json:"last_accessed"`
}

// 事件相关的辅助结构体

// RewardDistribution 奖励分发详情
type RewardDistribution struct {
	PlayerID       uint64 `json:"player_id"`
	PlayerName     string `json:"player_name"`
	PlayerRank     int64  `json:"player_rank"`
	RewardType     string `json:"reward_type"`
	RewardQuantity int64  `json:"reward_quantity"`
	RewardValue    int64  `json:"reward_value"`
	DistributionStatus string `json:"distribution_status"` // "success", "failed", "pending"
	ErrorMessage   string `json:"error_message,omitempty"`
}

// 事件常量

const (
	// 生命周期事件类型
	EventTypeRankingCreated  = "ranking.created"
	EventTypeRankingDeleted  = "ranking.deleted"
	EventTypeRankingUpdated  = "ranking.updated"
	
	// 状态事件类型
	EventTypeRankingStatusChanged = "ranking.status_changed"
	EventTypeRankingReset         = "ranking.reset"
	EventTypeRankingArchived      = "ranking.archived"
	
	// 玩家事件类型
	EventTypePlayerJoinedRanking  = "ranking.player_joined"
	EventTypePlayerLeftRanking    = "ranking.player_left"
	EventTypePlayerScoreUpdated   = "ranking.player_score_updated"
	EventTypePlayerRankChanged    = "ranking.player_rank_changed"
	
	// 黑名单事件类型
	EventTypePlayerBlacklisted   = "ranking.player_blacklisted"
	EventTypePlayerUnblacklisted = "ranking.player_unblacklisted"
	EventTypeBlacklistExpired    = "ranking.blacklist_expired"
	
	// 奖励事件类型
	EventTypeRankingRewardDistributed = "ranking.reward_distributed"
	EventTypePlayerRewardEarned       = "ranking.player_reward_earned"
	EventTypePlayerRewardClaimed      = "ranking.player_reward_claimed"
	
	// 统计事件类型
	EventTypeRankingStatisticsUpdated = "ranking.statistics_updated"
	EventTypeRankingMilestoneReached  = "ranking.milestone_reached"
	EventTypeRankingRecordBroken      = "ranking.record_broken"
	
	// 系统事件类型
	EventTypeRankingMaintenance    = "ranking.maintenance"
	EventTypeRankingDataMigration  = "ranking.data_migration"
	EventTypeRankingPerformance    = "ranking.performance"
	
	// 缓存事件类型
	EventTypeRankingCacheUpdated = "ranking.cache_updated"
	EventTypeRankingCacheEvicted = "ranking.cache_evicted"
)

// 事件工厂函数

// NewRankingCreatedEvent 创建排行榜创建事件
func NewRankingCreatedEvent(aggregateID string, rankID uint32, name string, rankType RankType, category RankCategory, createdBy string) *RankingCreatedEvent {
	return &RankingCreatedEvent{
		BaseRankingEvent: BaseRankingEvent{
			EventID:     generateEventID(),
			EventType:   EventTypeRankingCreated,
			AggregateID: aggregateID,
			RankID:      rankID,
			Timestamp:   time.Now(),
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		Name:      name,
		RankType:  rankType,
		Category:  category,
		CreatedBy: createdBy,
	}
}

// NewPlayerJoinedRankingEvent 创建玩家加入排行榜事件
func NewPlayerJoinedRankingEvent(aggregateID string, playerID uint64, score, timeScore int64) *PlayerJoinedRankingEvent {
	return &PlayerJoinedRankingEvent{
		BaseRankingEvent: BaseRankingEvent{
			EventID:     generateEventID(),
			EventType:   EventTypePlayerJoinedRanking,
			AggregateID: aggregateID,
			Timestamp:   time.Now(),
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		PlayerID:     playerID,
		InitialScore: score,
		TimeScore:    timeScore,
		JoinMethod:   "first_score",
	}
}

// NewPlayerScoreUpdatedEvent 创建玩家分数更新事件
func NewPlayerScoreUpdatedEvent(aggregateID string, playerID uint64, oldScore, newScore, timeScore int64) *PlayerScoreUpdatedEvent {
	return &PlayerScoreUpdatedEvent{
		BaseRankingEvent: BaseRankingEvent{
			EventID:     generateEventID(),
			EventType:   EventTypePlayerScoreUpdated,
			AggregateID: aggregateID,
			Timestamp:   time.Now(),
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		PlayerID:     playerID,
		OldScore:     oldScore,
		NewScore:     newScore,
		ScoreDelta:   newScore - oldScore,
		NewTimeScore: timeScore,
		UpdateSource: "game",
	}
}

// NewPlayerBlacklistedEvent 创建玩家黑名单事件
func NewPlayerBlacklistedEvent(aggregateID string, playerID uint64, reason string) *PlayerBlacklistedEvent {
	return &PlayerBlacklistedEvent{
		BaseRankingEvent: BaseRankingEvent{
			EventID:     generateEventID(),
			EventType:   EventTypePlayerBlacklisted,
			AggregateID: aggregateID,
			Timestamp:   time.Now(),
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		PlayerID:    playerID,
		Reason:      reason,
		IsPermanent: true,
	}
}

// NewPlayerUnblacklistedEvent 创建玩家解除黑名单事件
func NewPlayerUnblacklistedEvent(aggregateID string, playerID uint64) *PlayerUnblacklistedEvent {
	return &PlayerUnblacklistedEvent{
		BaseRankingEvent: BaseRankingEvent{
			EventID:     generateEventID(),
			EventType:   EventTypePlayerUnblacklisted,
			AggregateID: aggregateID,
			Timestamp:   time.Now(),
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		PlayerID:  playerID,
		CanRejoin: true,
	}
}

// NewRankingResetEvent 创建排行榜重置事件
func NewRankingResetEvent(aggregateID string, previousPlayerCount int) *RankingResetEvent {
	return &RankingResetEvent{
		BaseRankingEvent: BaseRankingEvent{
			EventID:     generateEventID(),
			EventType:   EventTypeRankingReset,
			AggregateID: aggregateID,
			Timestamp:   time.Now(),
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		PreviousPlayerCount: previousPlayerCount,
		ResetReason:         "manual_reset",
		BackupCreated:       false,
	}
}

// NewRankingStatusChangedEvent 创建排行榜状态改变事件
func NewRankingStatusChangedEvent(aggregateID string, oldActive, newActive bool) *RankingStatusChangedEvent {
	return &RankingStatusChangedEvent{
		BaseRankingEvent: BaseRankingEvent{
			EventID:     generateEventID(),
			EventType:   EventTypeRankingStatusChanged,
			AggregateID: aggregateID,
			Timestamp:   time.Now(),
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		OldActive: oldActive,
		NewActive: newActive,
		Reason:    "manual_change",
	}
}

// 事件处理器接口

// RankingEventHandler 排行榜事件处理器接口
type RankingEventHandler interface {
	Handle(event RankingEvent) error
	CanHandle(eventType string) bool
	GetHandlerName() string
}

// RankingEventBus 排行榜事件总线接口
type RankingEventBus interface {
	// 发布事件
	Publish(event RankingEvent) error
	PublishBatch(events []RankingEvent) error
	
	// 订阅事件
	Subscribe(eventType string, handler RankingEventHandler) error
	Unsubscribe(eventType string, handlerName string) error
	
	// 事件存储
	Store(event RankingEvent) error
	GetEvents(aggregateID string, fromVersion int) ([]RankingEvent, error)
	GetEventsByType(eventType string, limit int) ([]RankingEvent, error)
	GetEventsByRankID(rankID uint32, limit int) ([]RankingEvent, error)
	
	// 事件重放
	Replay(aggregateID string, fromVersion int, handler RankingEventHandler) error
	
	// 快照管理
	CreateSnapshot(aggregateID string, version int, data interface{}) error
	GetSnapshot(aggregateID string) (interface{}, int, error)
	
	// 事件清理
	CleanupEvents(beforeTime time.Time) error
	ArchiveEvents(beforeTime time.Time) error
}

// 事件聚合器接口

// RankingEventAggregator 排行榜事件聚合器接口
type RankingEventAggregator interface {
	// 聚合统计
	AggregatePlayerActivity(rankID uint32, period time.Duration) (*PlayerActivityStats, error)
	AggregateScoreChanges(rankID uint32, period time.Duration) (*ScoreChangeStats, error)
	AggregateRankingHealth(rankID uint32, period time.Duration) (*RankingHealthStats, error)
	
	// 趋势分析
	AnalyzePlayerTrends(rankID uint32, playerID uint64, period time.Duration) (*PlayerTrendAnalysis, error)
	AnalyzeRankingTrends(rankID uint32, period time.Duration) (*RankingTrendAnalysis, error)
	
	// 异常检测
	DetectAnomalies(rankID uint32, period time.Duration) ([]*RankingAnomaly, error)
	DetectSuspiciousActivity(rankID uint32, period time.Duration) ([]*SuspiciousActivity, error)
	
	// 性能分析
	AnalyzePerformance(rankID uint32, period time.Duration) (*RankingPerformanceAnalysis, error)
}

// 事件聚合统计结构体

// PlayerActivityStats 玩家活动统计
type PlayerActivityStats struct {
	RankID          uint32    `json:"rank_id"`
	Period          time.Duration `json:"period"`
	TotalPlayers    int64     `json:"total_players"`
	ActivePlayers   int64     `json:"active_players"`
	NewPlayers      int64     `json:"new_players"`
	LeavingPlayers  int64     `json:"leaving_players"`
	RetentionRate   float64   `json:"retention_rate"`
	ChurnRate       float64   `json:"churn_rate"`
	AverageSessionTime time.Duration `json:"average_session_time"`
	PeakActivityTime time.Time `json:"peak_activity_time"`
	CalculatedAt    time.Time `json:"calculated_at"`
}

// ScoreChangeStats 分数变化统计
type ScoreChangeStats struct {
	RankID              uint32    `json:"rank_id"`
	Period              time.Duration `json:"period"`
	TotalUpdates        int64     `json:"total_updates"`
	AverageScoreChange  float64   `json:"average_score_change"`
	MaxScoreIncrease    int64     `json:"max_score_increase"`
	MaxScoreDecrease    int64     `json:"max_score_decrease"`
	ScoreVolatility     float64   `json:"score_volatility"`
	TopScoreChanges     []*ScoreChange `json:"top_score_changes"`
	CalculatedAt        time.Time `json:"calculated_at"`
}

// RankingHealthStats 排行榜健康统计
type RankingHealthStats struct {
	RankID              uint32    `json:"rank_id"`
	Period              time.Duration `json:"period"`
	HealthScore         float64   `json:"health_score"`
	Competitiveness     float64   `json:"competitiveness"`
	EngagementLevel     float64   `json:"engagement_level"`
	StabilityIndex      float64   `json:"stability_index"`
	GrowthRate          float64   `json:"growth_rate"`
	ErrorRate           float64   `json:"error_rate"`
	PerformanceScore    float64   `json:"performance_score"`
	Recommendations     []string  `json:"recommendations"`
	CalculatedAt        time.Time `json:"calculated_at"`
}

// PlayerTrendAnalysis 玩家趋势分析
type PlayerTrendAnalysis struct {
	RankID           uint32    `json:"rank_id"`
	PlayerID         uint64    `json:"player_id"`
	Period           time.Duration `json:"period"`
	ScoreTrend       string    `json:"score_trend"` // "increasing", "decreasing", "stable", "volatile"
	RankTrend        string    `json:"rank_trend"`
	ActivityLevel    string    `json:"activity_level"` // "high", "medium", "low"
	Consistency      float64   `json:"consistency"`
	Improvement      float64   `json:"improvement"`
	PredictedRank    int64     `json:"predicted_rank"`
	RiskFactors      []string  `json:"risk_factors"`
	CalculatedAt     time.Time `json:"calculated_at"`
}

// RankingTrendAnalysis 排行榜趋势分析
type RankingTrendAnalysis struct {
	RankID              uint32    `json:"rank_id"`
	Period              time.Duration `json:"period"`
	OverallTrend        string    `json:"overall_trend"` // "growing", "declining", "stable"
	PlayerGrowthRate    float64   `json:"player_growth_rate"`
	ScoreInflation      float64   `json:"score_inflation"`
	CompetitionLevel    string    `json:"competition_level"` // "high", "medium", "low"
	Seasonality         []SeasonalPattern `json:"seasonality"`
	Predictions         *RankingPredictions `json:"predictions"`
	Recommendations     []string  `json:"recommendations"`
	CalculatedAt        time.Time `json:"calculated_at"`
}

// RankingAnomaly 排行榜异常
type RankingAnomaly struct {
	RankID       uint32    `json:"rank_id"`
	AnomalyType  string    `json:"anomaly_type"` // "score_spike", "mass_exodus", "unusual_pattern"
	Severity     string    `json:"severity"` // "low", "medium", "high", "critical"
	Description  string    `json:"description"`
	AffectedPlayers []uint64 `json:"affected_players"`
	DetectedAt   time.Time `json:"detected_at"`
	StartTime    time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time,omitempty"`
	Metrics      map[string]float64 `json:"metrics"`
	RecommendedActions []string `json:"recommended_actions"`
}

// SuspiciousActivity 可疑活动
type SuspiciousActivity struct {
	RankID       uint32    `json:"rank_id"`
	PlayerID     uint64    `json:"player_id"`
	ActivityType string    `json:"activity_type"` // "rapid_score_increase", "impossible_score", "bot_behavior"
	RiskLevel    string    `json:"risk_level"` // "low", "medium", "high"
	Confidence   float64   `json:"confidence"`
	Description  string    `json:"description"`
	Evidence     []string  `json:"evidence"`
	DetectedAt   time.Time `json:"detected_at"`
	RecommendedActions []string `json:"recommended_actions"`
}

// RankingPerformanceAnalysis 排行榜性能分析
type RankingPerformanceAnalysis struct {
	RankID              uint32    `json:"rank_id"`
	Period              time.Duration `json:"period"`
	AverageResponseTime time.Duration `json:"average_response_time"`
	Throughput          float64   `json:"throughput"`
	ErrorRate           float64   `json:"error_rate"`
	CacheHitRate        float64   `json:"cache_hit_rate"`
	MemoryUsage         int64     `json:"memory_usage"`
	CPUUsage            float64   `json:"cpu_usage"`
	Bottlenecks         []string  `json:"bottlenecks"`
	Optimizations       []string  `json:"optimizations"`
	CalculatedAt        time.Time `json:"calculated_at"`
}

// 辅助结构体

// ScoreChange 分数变化
type ScoreChange struct {
	PlayerID    uint64    `json:"player_id"`
	OldScore    int64     `json:"old_score"`
	NewScore    int64     `json:"new_score"`
	Change      int64     `json:"change"`
	Timestamp   time.Time `json:"timestamp"`
}

// SeasonalPattern 季节性模式
type SeasonalPattern struct {
	Period      string  `json:"period"` // "daily", "weekly", "monthly"
	Pattern     string  `json:"pattern"`
	Strength    float64 `json:"strength"`
	PeakTimes   []string `json:"peak_times"`
}

// RankingPredictions 排行榜预测
type RankingPredictions struct {
	NextWeekPlayers   int64   `json:"next_week_players"`
	NextMonthPlayers  int64   `json:"next_month_players"`
	ScoreGrowthRate   float64 `json:"score_growth_rate"`
	ChurnProbability  float64 `json:"churn_probability"`
	ConfidenceLevel   float64 `json:"confidence_level"`
}

// 辅助函数

// generateEventID 生成事件ID
func generateEventID() string {
	return fmt.Sprintf("ranking_event_%d", time.Now().UnixNano())
}

// ValidateEvent 验证事件
func ValidateEvent(event RankingEvent) error {
	if event.GetEventID() == "" {
		return fmt.Errorf("event ID cannot be empty")
	}
	if event.GetEventType() == "" {
		return fmt.Errorf("event type cannot be empty")
	}
	if event.GetAggregateID() == "" {
		return fmt.Errorf("aggregate ID cannot be empty")
	}
	if event.GetRankID() == 0 {
		return fmt.Errorf("rank ID cannot be zero")
	}
	if event.GetTimestamp().IsZero() {
		return fmt.Errorf("timestamp cannot be zero")
	}
	return nil
}

// SerializeEvent 序列化事件
func SerializeEvent(event RankingEvent) ([]byte, error) {
	// 实现事件序列化逻辑
	return nil, nil
}

// DeserializeEvent 反序列化事件
func DeserializeEvent(data []byte, eventType string) (RankingEvent, error) {
	// 实现事件反序列化逻辑
	return nil, nil
}