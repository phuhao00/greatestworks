package hangup

import (
	"context"
	"time"
)

// HangupRepository 挂机仓储接口
type HangupRepository interface {
	// Save 保存挂机聚合根
	Save(ctx context.Context, hangup *HangupAggregate) error
	
	// FindByPlayerID 根据玩家ID查找挂机信息
	FindByPlayerID(ctx context.Context, playerID string) (*HangupAggregate, error)
	
	// Delete 删除挂机信息
	Delete(ctx context.Context, playerID string) error
	
	// FindActiveHangups 查找活跃的挂机玩家
	FindActiveHangups(ctx context.Context, limit, offset int) ([]*HangupAggregate, error)
	
	// FindByLocation 根据地点查找挂机玩家
	FindByLocation(ctx context.Context, locationID string, limit, offset int) ([]*HangupAggregate, error)
	
	// UpdateLastOnlineTime 更新最后在线时间
	UpdateLastOnlineTime(ctx context.Context, playerID string, timestamp time.Time) error
	
	// UpdateLastOfflineTime 更新最后离线时间
	UpdateLastOfflineTime(ctx context.Context, playerID string, timestamp time.Time) error
	
	// GetHangupStatistics 获取挂机统计信息
	GetHangupStatistics(ctx context.Context, playerID string) (*HangupStatistics, error)
	
	// UpdateHangupStatistics 更新挂机统计信息
	UpdateHangupStatistics(ctx context.Context, stats *HangupStatistics) error
	
	// FindTopHangupPlayers 查找挂机排行榜
	FindTopHangupPlayers(ctx context.Context, rankType string, limit int) ([]*HangupRank, error)
	
	// GetPlayerHangupRank 获取玩家挂机排名
	GetPlayerHangupRank(ctx context.Context, playerID string, rankType string) (int, error)
}

// HangupLocationRepository 挂机地点仓储接口
type HangupLocationRepository interface {
	// SaveLocation 保存挂机地点
	SaveLocation(ctx context.Context, location *HangupLocation) error
	
	// FindLocationByID 根据ID查找地点
	FindLocationByID(ctx context.Context, locationID string) (*HangupLocation, error)
	
	// FindAllLocations 查找所有地点
	FindAllLocations(ctx context.Context) ([]*HangupLocation, error)
	
	// FindLocationsByType 根据类型查找地点
	FindLocationsByType(ctx context.Context, locationType LocationType) ([]*HangupLocation, error)
	
	// FindUnlockedLocations 查找已解锁的地点
	FindUnlockedLocations(ctx context.Context, playerID string) ([]*HangupLocation, error)
	
	// FindLocationsByLevel 根据等级要求查找地点
	FindLocationsByLevel(ctx context.Context, minLevel, maxLevel int) ([]*HangupLocation, error)
	
	// UpdateLocationStatus 更新地点状态
	UpdateLocationStatus(ctx context.Context, locationID string, isActive bool) error
	
	// DeleteLocation 删除地点
	DeleteLocation(ctx context.Context, locationID string) error
	
	// GetLocationStatistics 获取地点统计信息
	GetLocationStatistics(ctx context.Context, locationID string) (*LocationStatistics, error)
}

// HangupSessionRepository 挂机会话仓储接口
type HangupSessionRepository interface {
	// SaveSession 保存挂机会话
	SaveSession(ctx context.Context, session *HangupSession) error
	
	// FindSessionByID 根据ID查找会话
	FindSessionByID(ctx context.Context, sessionID string) (*HangupSession, error)
	
	// FindSessionsByPlayer 根据玩家查找会话
	FindSessionsByPlayer(ctx context.Context, playerID string, limit, offset int) ([]*HangupSession, error)
	
	// FindActiveSessions 查找活跃会话
	FindActiveSessions(ctx context.Context, limit, offset int) ([]*HangupSession, error)
	
	// FindSessionsByLocation 根据地点查找会话
	FindSessionsByLocation(ctx context.Context, locationID string, limit, offset int) ([]*HangupSession, error)
	
	// FindSessionsByTimeRange 根据时间范围查找会话
	FindSessionsByTimeRange(ctx context.Context, startTime, endTime time.Time, limit, offset int) ([]*HangupSession, error)
	
	// EndSession 结束会话
	EndSession(ctx context.Context, sessionID string, reward *BaseReward) error
	
	// DeleteSession 删除会话
	DeleteSession(ctx context.Context, sessionID string) error
	
	// GetSessionStatistics 获取会话统计信息
	GetSessionStatistics(ctx context.Context, playerID string, timeRange string) (*SessionStatistics, error)
}

// HangupRewardRepository 挂机奖励仓储接口
type HangupRewardRepository interface {
	// SaveOfflineReward 保存离线奖励
	SaveOfflineReward(ctx context.Context, playerID string, reward *OfflineReward) error
	
	// FindOfflineReward 查找离线奖励
	FindOfflineReward(ctx context.Context, playerID string) (*OfflineReward, error)
	
	// ClaimOfflineReward 领取离线奖励
	ClaimOfflineReward(ctx context.Context, playerID string) (*OfflineReward, error)
	
	// DeleteOfflineReward 删除离线奖励
	DeleteOfflineReward(ctx context.Context, playerID string) error
	
	// FindUnclaimedRewards 查找未领取的奖励
	FindUnclaimedRewards(ctx context.Context, limit, offset int) ([]*OfflineReward, error)
	
	// GetRewardHistory 获取奖励历史
	GetRewardHistory(ctx context.Context, playerID string, limit, offset int) ([]*OfflineReward, error)
	
	// GetTotalRewards 获取总奖励统计
	GetTotalRewards(ctx context.Context, playerID string) (*RewardSummary, error)
}

// HangupConfigRepository 挂机配置仓储接口
type HangupConfigRepository interface {
	// SaveConfig 保存配置
	SaveConfig(ctx context.Context, config *HangupConfig) error
	
	// GetConfig 获取配置
	GetConfig(ctx context.Context) (*HangupConfig, error)
	
	// UpdateConfig 更新配置
	UpdateConfig(ctx context.Context, config *HangupConfig) error
	
	// GetConfigHistory 获取配置历史
	GetConfigHistory(ctx context.Context, limit, offset int) ([]*HangupConfig, error)
}

// 统计信息结构体

// LocationStatistics 地点统计信息
type LocationStatistics struct {
	LocationID       string        `json:"location_id"`
	LocationName     string        `json:"location_name"`
	TotalPlayers     int64         `json:"total_players"`
	ActiveUsers      int64         `json:"active_users"`
	AverageSession   time.Duration `json:"average_session"`
	TotalExperience  int64         `json:"total_experience"`
	TotalGold        int64         `json:"total_gold"`
	PopularityRank   int           `json:"popularity_rank"`
	LastUpdated      time.Time     `json:"last_updated"`
}

// SessionStatistics 会话统计信息
type SessionStatistics struct {
	PlayerID         string        `json:"player_id"`
	TotalSessions    int64         `json:"total_sessions"`
	TotalDuration    time.Duration `json:"total_duration"`
	AverageSession   time.Duration `json:"average_session"`
	LongestSession   time.Duration `json:"longest_session"`
	ShortestSession  time.Duration `json:"shortest_session"`
	OnlineSessions   int64         `json:"online_sessions"`
	OfflineSessions  int64         `json:"offline_sessions"`
	FavoriteLocation string        `json:"favorite_location"`
	LastSession      time.Time     `json:"last_session"`
}

// RewardSummary 奖励汇总
type RewardSummary struct {
	PlayerID         string    `json:"player_id"`
	TotalExperience  int64     `json:"total_experience"`
	TotalGold        int64     `json:"total_gold"`
	TotalItems       int       `json:"total_items"`
	TotalRewards     int64     `json:"total_rewards"`
	ClaimedRewards   int64     `json:"claimed_rewards"`
	UnclaimedRewards int64     `json:"unclaimed_rewards"`
	LastReward       time.Time `json:"last_reward"`
	LastClaimed      time.Time `json:"last_claimed"`
}

// HangupQuery 挂机查询条件
type HangupQuery struct {
	PlayerIDs    []string      `json:"player_ids,omitempty"`
	LocationIDs  []string      `json:"location_ids,omitempty"`
	Status       []string      `json:"status,omitempty"`
	MinLevel     int           `json:"min_level,omitempty"`
	MaxLevel     int           `json:"max_level,omitempty"`
	StartTime    time.Time     `json:"start_time,omitempty"`
	EndTime      time.Time     `json:"end_time,omitempty"`
	MinDuration  time.Duration `json:"min_duration,omitempty"`
	MaxDuration  time.Duration `json:"max_duration,omitempty"`
	IsOnline     *bool         `json:"is_online,omitempty"`
	HasReward    *bool         `json:"has_reward,omitempty"`
	Limit        int           `json:"limit,omitempty"`
	Offset       int           `json:"offset,omitempty"`
	SortBy       string        `json:"sort_by,omitempty"`
	SortOrder    string        `json:"sort_order,omitempty"`
}

// HangupQueryRepository 挂机查询仓储接口
type HangupQueryRepository interface {
	// FindByQuery 根据查询条件查找挂机信息
	FindByQuery(ctx context.Context, query *HangupQuery) ([]*HangupAggregate, error)
	
	// CountByQuery 根据查询条件统计数量
	CountByQuery(ctx context.Context, query *HangupQuery) (int64, error)
	
	// FindPlayersInLocation 查找在指定地点挂机的玩家
	FindPlayersInLocation(ctx context.Context, locationID string, isOnline bool) ([]*HangupAggregate, error)
	
	// FindLongTermHangups 查找长期挂机的玩家
	FindLongTermHangups(ctx context.Context, minDuration time.Duration) ([]*HangupAggregate, error)
	
	// FindInactiveHangups 查找不活跃的挂机
	FindInactiveHangups(ctx context.Context, inactiveDuration time.Duration) ([]*HangupAggregate, error)
	
	// GetHangupTrends 获取挂机趋势数据
	GetHangupTrends(ctx context.Context, timeRange string, groupBy string) ([]TrendData, error)
	
	// GetLocationPopularity 获取地点热度数据
	GetLocationPopularity(ctx context.Context, timeRange string) ([]LocationPopularity, error)
	
	// GetPlayerHangupSummary 获取玩家挂机汇总
	GetPlayerHangupSummary(ctx context.Context, playerID string, timeRange string) (*PlayerHangupSummary, error)
}

// 趋势和分析数据结构体

// TrendData 趋势数据
type TrendData struct {
	Timestamp    time.Time `json:"timestamp"`
	ActiveUsers  int64     `json:"active_users"`
	TotalSessions int64    `json:"total_sessions"`
	AverageDuration time.Duration `json:"average_duration"`
	TotalRewards int64     `json:"total_rewards"`
}

// LocationPopularity 地点热度
type LocationPopularity struct {
	LocationID   string  `json:"location_id"`
	LocationName string  `json:"location_name"`
	UserCount    int64   `json:"user_count"`
	SessionCount int64   `json:"session_count"`
	PopularityScore float64 `json:"popularity_score"`
	Rank         int     `json:"rank"`
}

// PlayerHangupSummary 玩家挂机汇总
type PlayerHangupSummary struct {
	PlayerID           string                    `json:"player_id"`
	TotalHangupTime    time.Duration             `json:"total_hangup_time"`
	TotalSessions      int64                     `json:"total_sessions"`
	TotalExperience    int64                     `json:"total_experience"`
	TotalGold          int64                     `json:"total_gold"`
	TotalItems         int                       `json:"total_items"`
	FavoriteLocations  []LocationUsage           `json:"favorite_locations"`
	HangupEfficiency   float64                   `json:"hangup_efficiency"`
	Rank               int                       `json:"rank"`
	Achievements       []string                  `json:"achievements"`
	LastActive         time.Time                 `json:"last_active"`
	WeeklyStats        map[string]interface{}    `json:"weekly_stats"`
	MonthlyStats       map[string]interface{}    `json:"monthly_stats"`
}

// LocationUsage 地点使用情况
type LocationUsage struct {
	LocationID   string        `json:"location_id"`
	LocationName string        `json:"location_name"`
	UsageCount   int64         `json:"usage_count"`
	TotalTime    time.Duration `json:"total_time"`
	Percentage   float64       `json:"percentage"`
}

// HangupAnalyticsRepository 挂机分析仓储接口
type HangupAnalyticsRepository interface {
	// GetDailyStats 获取每日统计
	GetDailyStats(ctx context.Context, date time.Time) (*DailyHangupStats, error)
	
	// GetWeeklyStats 获取周统计
	GetWeeklyStats(ctx context.Context, year, week int) (*WeeklyHangupStats, error)
	
	// GetMonthlyStats 获取月统计
	GetMonthlyStats(ctx context.Context, year, month int) (*MonthlyHangupStats, error)
	
	// GetPlayerAnalytics 获取玩家分析数据
	GetPlayerAnalytics(ctx context.Context, playerID string, timeRange string) (*PlayerAnalytics, error)
	
	// GetLocationAnalytics 获取地点分析数据
	GetLocationAnalytics(ctx context.Context, locationID string, timeRange string) (*LocationAnalytics, error)
	
	// GetSystemAnalytics 获取系统分析数据
	GetSystemAnalytics(ctx context.Context, timeRange string) (*SystemAnalytics, error)
	
	// GenerateReport 生成报告
	GenerateReport(ctx context.Context, reportType string, params map[string]interface{}) ([]byte, error)
}

// 分析数据结构体

// DailyHangupStats 每日挂机统计
type DailyHangupStats struct {
	Date            time.Time `json:"date"`
	ActiveUsers     int64     `json:"active_users"`
	NewUsers        int64     `json:"new_users"`
	TotalSessions   int64     `json:"total_sessions"`
	AverageDuration time.Duration `json:"average_duration"`
	TotalRewards    int64     `json:"total_rewards"`
	TopLocations    []string  `json:"top_locations"`
}

// WeeklyHangupStats 周挂机统计
type WeeklyHangupStats struct {
	Year            int       `json:"year"`
	Week            int       `json:"week"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	ActiveUsers     int64     `json:"active_users"`
	NewUsers        int64     `json:"new_users"`
	TotalSessions   int64     `json:"total_sessions"`
	AverageDuration time.Duration `json:"average_duration"`
	TotalRewards    int64     `json:"total_rewards"`
	GrowthRate      float64   `json:"growth_rate"`
}

// MonthlyHangupStats 月挂机统计
type MonthlyHangupStats struct {
	Year            int       `json:"year"`
	Month           int       `json:"month"`
	ActiveUsers     int64     `json:"active_users"`
	NewUsers        int64     `json:"new_users"`
	TotalSessions   int64     `json:"total_sessions"`
	AverageDuration time.Duration `json:"average_duration"`
	TotalRewards    int64     `json:"total_rewards"`
	RetentionRate   float64   `json:"retention_rate"`
	ChurnRate       float64   `json:"churn_rate"`
}

// PlayerAnalytics 玩家分析数据
type PlayerAnalytics struct {
	PlayerID        string                 `json:"player_id"`
	HangupPattern   map[string]interface{} `json:"hangup_pattern"`
	EfficiencyTrend []float64              `json:"efficiency_trend"`
	LocationPreference []LocationUsage    `json:"location_preference"`
	RewardTrend     []int64                `json:"reward_trend"`
	BehaviorScore   float64                `json:"behavior_score"`
	Predictions     map[string]interface{} `json:"predictions"`
}

// LocationAnalytics 地点分析数据
type LocationAnalytics struct {
	LocationID      string                 `json:"location_id"`
	UsagePattern    map[string]interface{} `json:"usage_pattern"`
	UserDistribution map[string]int64      `json:"user_distribution"`
	EfficiencyStats map[string]float64     `json:"efficiency_stats"`
	PopularityTrend []float64              `json:"popularity_trend"`
	OptimizationSuggestions []string      `json:"optimization_suggestions"`
}

// SystemAnalytics 系统分析数据
type SystemAnalytics struct {
	OverallHealth   float64                `json:"overall_health"`
	PerformanceMetrics map[string]float64 `json:"performance_metrics"`
	UsageDistribution map[string]int64    `json:"usage_distribution"`
	TrendAnalysis   map[string]interface{} `json:"trend_analysis"`
	Anomalies       []string               `json:"anomalies"`
	Recommendations []string               `json:"recommendations"`
}