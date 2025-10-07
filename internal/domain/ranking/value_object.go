package ranking

import (
	"fmt"
	"time"
)

// 排行榜类型相关值对象

// RankType 排行榜类型
type RankType int32

const (
	RankTypeLevel       RankType = iota + 1 // 等级排行榜
	RankTypePower                           // 战力排行榜
	RankTypeWealth                          // 财富排行榜
	RankTypeAchievement                     // 成就排行榜
	RankTypePet                             // 宠物排行榜
	RankTypeGuild                           // 公会排行榜
	RankTypeArena                           // 竞技场排行榜
	RankTypeDungeon                         // 副本排行榜
	RankTypeActivity                        // 活动排行榜
	RankTypeCustom                          // 自定义排行榜
)

// String 返回排行榜类型的字符串表示
func (rt RankType) String() string {
	switch rt {
	case RankTypeLevel:
		return "level"
	case RankTypePower:
		return "power"
	case RankTypeWealth:
		return "wealth"
	case RankTypeAchievement:
		return "achievement"
	case RankTypePet:
		return "pet"
	case RankTypeGuild:
		return "guild"
	case RankTypeArena:
		return "arena"
	case RankTypeDungeon:
		return "dungeon"
	case RankTypeActivity:
		return "activity"
	case RankTypeCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// ParseRankType 解析排行榜类型
func ParseRankType(s string) RankType {
	switch s {
	case "level":
		return RankTypeLevel
	case "power":
		return RankTypePower
	case "wealth":
		return RankTypeWealth
	case "achievement":
		return RankTypeAchievement
	case "pet":
		return RankTypePet
	case "guild":
		return RankTypeGuild
	case "arena":
		return RankTypeArena
	case "dungeon":
		return RankTypeDungeon
	case "activity":
		return RankTypeActivity
	case "custom":
		return RankTypeCustom
	default:
		return RankTypeLevel // 默认值
	}
}

// IsValid 检查排行榜类型是否有效
func (rt RankType) IsValid() bool {
	return rt >= RankTypeLevel && rt <= RankTypeCustom
}

// RankCategory 排行榜分类
type RankCategory int32

const (
	RankCategoryPlayer RankCategory = iota + 1 // 玩家排行榜
	RankCategoryGuild                          // 公会排行榜
	RankCategoryServer                         // 服务器排行榜
	RankCategoryGlobal                         // 全球排行榜
	RankCategoryEvent                          // 活动排行榜
	RankCategorySeason                         // 赛季排行榜
)

// String 返回排行榜分类的字符串表示
func (rc RankCategory) String() string {
	switch rc {
	case RankCategoryPlayer:
		return "player"
	case RankCategoryGuild:
		return "guild"
	case RankCategoryServer:
		return "server"
	case RankCategoryGlobal:
		return "global"
	case RankCategoryEvent:
		return "event"
	case RankCategorySeason:
		return "season"
	default:
		return "unknown"
	}
}

// IsValid 检查排行榜分类是否有效
func (rc RankCategory) IsValid() bool {
	return rc >= RankCategoryPlayer && rc <= RankCategorySeason
}

// SortType 排序类型
type SortType int32

const (
	SortTypeDescending SortType = iota + 1 // 降序（从高到低）
	SortTypeAscending                      // 升序（从低到高）
)

// String 返回排序类型的字符串表示
func (st SortType) String() string {
	switch st {
	case SortTypeDescending:
		return "descending"
	case SortTypeAscending:
		return "ascending"
	default:
		return "unknown"
	}
}

// IsValid 检查排序类型是否有效
func (st SortType) IsValid() bool {
	return st == SortTypeDescending || st == SortTypeAscending
}

// RankPeriod 排行榜周期
type RankPeriod int32

const (
	RankPeriodPermanent RankPeriod = iota + 1 // 永久排行榜
	RankPeriodDaily                           // 日排行榜
	RankPeriodWeekly                          // 周排行榜
	RankPeriodMonthly                         // 月排行榜
	RankPeriodSeasonal                        // 赛季排行榜
	RankPeriodEvent                           // 活动排行榜
	RankPeriodCustom                          // 自定义周期
)

// String 返回排行榜周期的字符串表示
func (rp RankPeriod) String() string {
	switch rp {
	case RankPeriodPermanent:
		return "permanent"
	case RankPeriodDaily:
		return "daily"
	case RankPeriodWeekly:
		return "weekly"
	case RankPeriodMonthly:
		return "monthly"
	case RankPeriodSeasonal:
		return "seasonal"
	case RankPeriodEvent:
		return "event"
	case RankPeriodCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// ParsePeriodType 解析排行榜周期类型
func ParsePeriodType(s string) RankPeriod {
	switch s {
	case "permanent":
		return RankPeriodPermanent
	case "daily":
		return RankPeriodDaily
	case "weekly":
		return RankPeriodWeekly
	case "monthly":
		return RankPeriodMonthly
	case "seasonal":
		return RankPeriodSeasonal
	case "event":
		return RankPeriodEvent
	case "custom":
		return RankPeriodCustom
	default:
		return RankPeriodPermanent // 默认值
	}
}

// IsValid 检查排行榜周期是否有效
func (rp RankPeriod) IsValid() bool {
	return rp >= RankPeriodPermanent && rp <= RankPeriodCustom
}

// GetDuration 获取周期持续时间
func (rp RankPeriod) GetDuration() time.Duration {
	switch rp {
	case RankPeriodDaily:
		return 24 * time.Hour
	case RankPeriodWeekly:
		return 7 * 24 * time.Hour
	case RankPeriodMonthly:
		return 30 * 24 * time.Hour
	case RankPeriodSeasonal:
		return 90 * 24 * time.Hour
	default:
		return 0 // 永久或自定义周期
	}
}

// RankStatus 排行榜状态
type RankStatus int32

const (
	RankStatusActive      RankStatus = iota + 1 // 活跃状态
	RankStatusInactive                          // 非活跃状态
	RankStatusPaused                            // 暂停状态
	RankStatusExpired                           // 过期状态
	RankStatusMaintenance                       // 维护状态
	RankStatusArchived                          // 归档状态
)

// String 返回排行榜状态的字符串表示
func (rs RankStatus) String() string {
	switch rs {
	case RankStatusActive:
		return "active"
	case RankStatusInactive:
		return "inactive"
	case RankStatusPaused:
		return "paused"
	case RankStatusExpired:
		return "expired"
	case RankStatusMaintenance:
		return "maintenance"
	case RankStatusArchived:
		return "archived"
	default:
		return "unknown"
	}
}

// IsValid 检查排行榜状态是否有效
func (rs RankStatus) IsValid() bool {
	return rs >= RankStatusActive && rs <= RankStatusArchived
}

// CanAcceptUpdates 检查状态是否可以接受更新
func (rs RankStatus) CanAcceptUpdates() bool {
	return rs == RankStatusActive
}

// 排行榜配置相关值对象

// RankRewardConfig 排行榜奖励配置
type RankRewardConfig struct {
	Enabled        bool                   `json:"enabled" bson:"enabled"`
	RewardTiers    []*RankRewardTier      `json:"reward_tiers" bson:"reward_tiers"`
	RewardType     RankRewardType         `json:"reward_type" bson:"reward_type"`
	DistributeAt   RankRewardDistributeAt `json:"distribute_at" bson:"distribute_at"`
	AutoDistribute bool                   `json:"auto_distribute" bson:"auto_distribute"`
	CreatedAt      time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at" bson:"updated_at"`
}

// RankRewardTier 排行榜奖励档次
type RankRewardTier struct {
	MinRank  int64                  `json:"min_rank" bson:"min_rank"`
	MaxRank  int64                  `json:"max_rank" bson:"max_rank"`
	Rewards  []*RankReward          `json:"rewards" bson:"rewards"`
	Title    string                 `json:"title" bson:"title"`
	Badge    string                 `json:"badge" bson:"badge"`
	Metadata map[string]interface{} `json:"metadata" bson:"metadata"`
}

// RankReward 排行榜奖励
type RankReward struct {
	RewardID    string `json:"reward_id" bson:"reward_id"`
	RewardType  string `json:"reward_type" bson:"reward_type"`
	Quantity    int64  `json:"quantity" bson:"quantity"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

// RankRewardType 奖励类型
type RankRewardType int32

const (
	RankRewardTypeImmediate   RankRewardType = iota + 1 // 立即奖励
	RankRewardTypeDeferred                              // 延迟奖励
	RankRewardTypeMail                                  // 邮件奖励
	RankRewardTypeAchievement                           // 成就奖励
)

// String 返回奖励类型的字符串表示
func (rrt RankRewardType) String() string {
	switch rrt {
	case RankRewardTypeImmediate:
		return "immediate"
	case RankRewardTypeDeferred:
		return "deferred"
	case RankRewardTypeMail:
		return "mail"
	case RankRewardTypeAchievement:
		return "achievement"
	default:
		return "unknown"
	}
}

// RankRewardDistributeAt 奖励分发时机
type RankRewardDistributeAt int32

const (
	RankRewardDistributeAtPeriodEnd  RankRewardDistributeAt = iota + 1 // 周期结束时
	RankRewardDistributeAtRankChange                                   // 排名变化时
	RankRewardDistributeAtManual                                       // 手动分发
	RankRewardDistributeAtScheduled                                    // 定时分发
)

// String 返回奖励分发时机的字符串表示
func (rrda RankRewardDistributeAt) String() string {
	switch rrda {
	case RankRewardDistributeAtPeriodEnd:
		return "period_end"
	case RankRewardDistributeAtRankChange:
		return "rank_change"
	case RankRewardDistributeAtManual:
		return "manual"
	case RankRewardDistributeAtScheduled:
		return "scheduled"
	default:
		return "unknown"
	}
}

// RankCacheConfig 排行榜缓存配置
type RankCacheConfig struct {
	Enabled         bool          `json:"enabled" bson:"enabled"`
	CacheSize       int64         `json:"cache_size" bson:"cache_size"`
	CacheTTL        time.Duration `json:"cache_ttl" bson:"cache_ttl"`
	RefreshInterval time.Duration `json:"refresh_interval" bson:"refresh_interval"`
	PreloadTop      int64         `json:"preload_top" bson:"preload_top"`
	LazyLoad        bool          `json:"lazy_load" bson:"lazy_load"`
	Compression     bool          `json:"compression" bson:"compression"`
	CreatedAt       time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" bson:"updated_at"`
}

// 排行榜统计相关值对象

// RankingStatistics 排行榜统计信息
type RankingStatistics struct {
	RankID            uint32    `json:"rank_id" bson:"rank_id"`
	TotalPlayers      int64     `json:"total_players" bson:"total_players"`
	ActiveEntries     int64     `json:"active_entries" bson:"active_entries"`
	AverageScore      float64   `json:"average_score" bson:"average_score"`
	TopScore          int64     `json:"top_score" bson:"top_score"`
	LowestScore       int64     `json:"lowest_score" bson:"lowest_score"`
	ScoreRange        int64     `json:"score_range" bson:"score_range"`
	MedianScore       float64   `json:"median_score" bson:"median_score"`
	StandardDeviation float64   `json:"standard_deviation" bson:"standard_deviation"`
	BlacklistCount    int64     `json:"blacklist_count" bson:"blacklist_count"`
	LastUpdated       time.Time `json:"last_updated" bson:"last_updated"`
	LastScoreUpdate   time.Time `json:"last_score_update" bson:"last_score_update"`
	UpdateFrequency   float64   `json:"update_frequency" bson:"update_frequency"`
	PeakPlayers       int64     `json:"peak_players" bson:"peak_players"`
	PeakTime          time.Time `json:"peak_time" bson:"peak_time"`
}

// RankingTrend 排行榜趋势数据
type RankingTrend struct {
	RankID     uint32                  `json:"rank_id" bson:"rank_id"`
	Period     RankPeriod              `json:"period" bson:"period"`
	TrendData  []*RankingTrendPoint    `json:"trend_data" bson:"trend_data"`
	GrowthRate float64                 `json:"growth_rate" bson:"growth_rate"`
	Volatility float64                 `json:"volatility" bson:"volatility"`
	Prediction *RankingTrendPrediction `json:"prediction,omitempty" bson:"prediction,omitempty"`
	CreatedAt  time.Time               `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time               `json:"updated_at" bson:"updated_at"`
}

// RankingTrendPoint 排行榜趋势点
type RankingTrendPoint struct {
	Timestamp     time.Time `json:"timestamp" bson:"timestamp"`
	PlayerCount   int64     `json:"player_count" bson:"player_count"`
	AverageScore  float64   `json:"average_score" bson:"average_score"`
	TopScore      int64     `json:"top_score" bson:"top_score"`
	ScoreVariance float64   `json:"score_variance" bson:"score_variance"`
	NewPlayers    int64     `json:"new_players" bson:"new_players"`
	ActivePlayers int64     `json:"active_players" bson:"active_players"`
}

// RankingTrendPrediction 排行榜趋势预测
type RankingTrendPrediction struct {
	PredictedPlayerCount  int64     `json:"predicted_player_count" bson:"predicted_player_count"`
	PredictedAverageScore float64   `json:"predicted_average_score" bson:"predicted_average_score"`
	PredictedTopScore     int64     `json:"predicted_top_score" bson:"predicted_top_score"`
	ConfidenceLevel       float64   `json:"confidence_level" bson:"confidence_level"`
	PredictionTime        time.Time `json:"prediction_time" bson:"prediction_time"`
	ValidUntil            time.Time `json:"valid_until" bson:"valid_until"`
}

// 排行榜查询相关值对象

// RankingQuery 排行榜查询条件
type RankingQuery struct {
	RankID        *uint32       `json:"rank_id,omitempty"`
	RankType      *RankType     `json:"rank_type,omitempty"`
	Category      *RankCategory `json:"category,omitempty"`
	Status        *RankStatus   `json:"status,omitempty"`
	Period        *RankPeriod   `json:"period,omitempty"`
	IsActive      *bool         `json:"is_active,omitempty"`
	PlayerID      *uint64       `json:"player_id,omitempty"`
	MinScore      *int64        `json:"min_score,omitempty"`
	MaxScore      *int64        `json:"max_score,omitempty"`
	MinRank       *int64        `json:"min_rank,omitempty"`
	MaxRank       *int64        `json:"max_rank,omitempty"`
	StartTime     *time.Time    `json:"start_time,omitempty"`
	EndTime       *time.Time    `json:"end_time,omitempty"`
	CreatedAfter  *time.Time    `json:"created_after,omitempty"`
	CreatedBefore *time.Time    `json:"created_before,omitempty"`
	UpdatedAfter  *time.Time    `json:"updated_after,omitempty"`
	UpdatedBefore *time.Time    `json:"updated_before,omitempty"`
	Keywords      []string      `json:"keywords,omitempty"`
	Tags          []string      `json:"tags,omitempty"`
	OrderBy       string        `json:"order_by,omitempty"`
	OrderDesc     bool          `json:"order_desc,omitempty"`
	Offset        int           `json:"offset,omitempty"`
	Limit         int           `json:"limit,omitempty"`
}

// RankingRange 排行榜范围
type RankingRange struct {
	Start        int64 `json:"start"`
	End          int64 `json:"end"`
	IncludeStart bool  `json:"include_start"`
	IncludeEnd   bool  `json:"include_end"`
}

// NewRankingRange 创建排行榜范围
func NewRankingRange(start, end int64) *RankingRange {
	return &RankingRange{
		Start:        start,
		End:          end,
		IncludeStart: true,
		IncludeEnd:   true,
	}
}

// IsValid 检查范围是否有效
func (rr *RankingRange) IsValid() bool {
	return rr.Start >= 0 && rr.End >= rr.Start
}

// Size 获取范围大小
func (rr *RankingRange) Size() int64 {
	if !rr.IsValid() {
		return 0
	}
	return rr.End - rr.Start + 1
}

// Contains 检查是否包含指定位置
func (rr *RankingRange) Contains(position int64) bool {
	if !rr.IsValid() {
		return false
	}

	startCheck := position > rr.Start || (rr.IncludeStart && position == rr.Start)
	endCheck := position < rr.End || (rr.IncludeEnd && position == rr.End)

	return startCheck && endCheck
}

// 排行榜过滤相关值对象

// RankingFilter 排行榜过滤器
type RankingFilter struct {
	ExcludeBlacklisted bool                   `json:"exclude_blacklisted"`
	ExcludeInactive    bool                   `json:"exclude_inactive"`
	MinScore           *int64                 `json:"min_score,omitempty"`
	MaxScore           *int64                 `json:"max_score,omitempty"`
	PlayerIDs          []uint64               `json:"player_ids,omitempty"`
	ExcludePlayerIDs   []uint64               `json:"exclude_player_ids,omitempty"`
	ScoreRange         *RankingRange          `json:"score_range,omitempty"`
	TimeRange          *TimeRange             `json:"time_range,omitempty"`
	CustomFilters      map[string]interface{} `json:"custom_filters,omitempty"`
}

// TimeRange 时间范围
type TimeRange struct {
	Start *time.Time `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}

// NewTimeRange 创建时间范围
func NewTimeRange(start, end *time.Time) *TimeRange {
	return &TimeRange{
		Start: start,
		End:   end,
	}
}

// IsValid 检查时间范围是否有效
func (tr *TimeRange) IsValid() bool {
	if tr.Start == nil && tr.End == nil {
		return true // 无限制
	}
	if tr.Start != nil && tr.End != nil {
		return tr.Start.Before(*tr.End) || tr.Start.Equal(*tr.End)
	}
	return true // 单边限制
}

// Contains 检查是否包含指定时间
func (tr *TimeRange) Contains(t time.Time) bool {
	if !tr.IsValid() {
		return false
	}

	if tr.Start != nil && t.Before(*tr.Start) {
		return false
	}

	if tr.End != nil && t.After(*tr.End) {
		return false
	}

	return true
}

// 排行榜操作相关值对象

// RankingOperation 排行榜操作类型
type RankingOperation int32

const (
	RankingOperationUpdate   RankingOperation = iota + 1 // 更新分数
	RankingOperationRemove                               // 移除玩家
	RankingOperationReset                                // 重置排行榜
	RankingOperationFreeze                               // 冻结排行榜
	RankingOperationUnfreeze                             // 解冻排行榜
	RankingOperationArchive                              // 归档排行榜
	RankingOperationRestore                              // 恢复排行榜
)

// String 返回排行榜操作的字符串表示
func (ro RankingOperation) String() string {
	switch ro {
	case RankingOperationUpdate:
		return "update"
	case RankingOperationRemove:
		return "remove"
	case RankingOperationReset:
		return "reset"
	case RankingOperationFreeze:
		return "freeze"
	case RankingOperationUnfreeze:
		return "unfreeze"
	case RankingOperationArchive:
		return "archive"
	case RankingOperationRestore:
		return "restore"
	default:
		return "unknown"
	}
}

// IsValid 检查排行榜操作是否有效
func (ro RankingOperation) IsValid() bool {
	return ro >= RankingOperationUpdate && ro <= RankingOperationRestore
}

// RequiresPermission 检查操作是否需要权限
func (ro RankingOperation) RequiresPermission() bool {
	switch ro {
	case RankingOperationReset, RankingOperationFreeze, RankingOperationUnfreeze,
		RankingOperationArchive, RankingOperationRestore:
		return true
	default:
		return false
	}
}

// RankingOperationResult 排行榜操作结果
type RankingOperationResult struct {
	Success       bool                   `json:"success"`
	Operation     RankingOperation       `json:"operation"`
	RankID        uint32                 `json:"rank_id"`
	PlayerID      *uint64                `json:"player_id,omitempty"`
	OldRank       *int64                 `json:"old_rank,omitempty"`
	NewRank       *int64                 `json:"new_rank,omitempty"`
	OldScore      *int64                 `json:"old_score,omitempty"`
	NewScore      *int64                 `json:"new_score,omitempty"`
	AffectedCount int64                  `json:"affected_count"`
	Message       string                 `json:"message"`
	Error         string                 `json:"error,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	Duration      time.Duration          `json:"duration"`
}

// NewRankingOperationResult 创建排行榜操作结果
func NewRankingOperationResult(operation RankingOperation, rankID uint32, success bool) *RankingOperationResult {
	return &RankingOperationResult{
		Success:   success,
		Operation: operation,
		RankID:    rankID,
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
	}
}

// SetPlayerInfo 设置玩家信息
func (ror *RankingOperationResult) SetPlayerInfo(playerID uint64, oldRank, newRank, oldScore, newScore *int64) {
	ror.PlayerID = &playerID
	ror.OldRank = oldRank
	ror.NewRank = newRank
	ror.OldScore = oldScore
	ror.NewScore = newScore
}

// SetError 设置错误信息
func (ror *RankingOperationResult) SetError(err error) {
	ror.Success = false
	ror.Error = err.Error()
}

// SetMessage 设置消息
func (ror *RankingOperationResult) SetMessage(message string) {
	ror.Message = message
}

// SetDuration 设置持续时间
func (ror *RankingOperationResult) SetDuration(start time.Time) {
	ror.Duration = time.Since(start)
}

// AddMetadata 添加元数据
func (ror *RankingOperationResult) AddMetadata(key string, value interface{}) {
	if ror.Metadata == nil {
		ror.Metadata = make(map[string]interface{})
	}
	ror.Metadata[key] = value
}

// 验证函数

// ValidateRankingRange 验证排行榜范围
func ValidateRankingRange(start, end int64) error {
	if start < 0 {
		return fmt.Errorf("start position cannot be negative: %d", start)
	}
	if end < start {
		return fmt.Errorf("end position cannot be less than start: start=%d, end=%d", start, end)
	}
	if end-start > 1000 {
		return fmt.Errorf("range too large: max 1000, requested %d", end-start+1)
	}
	return nil
}

// ValidateRankingQuery 验证排行榜查询
func ValidateRankingQuery(query *RankingQuery) error {
	if query == nil {
		return fmt.Errorf("query cannot be nil")
	}

	if query.Limit != 0 && query.Limit <= 0 {
		return fmt.Errorf("limit must be positive")
	}

	if query.Limit != 0 && query.Limit > 1000 {
		return fmt.Errorf("limit cannot exceed 1000")
	}

	if query.Offset != 0 && query.Offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}

	if query.MinScore != nil && query.MaxScore != nil && *query.MinScore > *query.MaxScore {
		return fmt.Errorf("min_score cannot be greater than max_score")
	}

	if query.MinRank != nil && query.MaxRank != nil && *query.MinRank > *query.MaxRank {
		return fmt.Errorf("min_rank cannot be greater than max_rank")
	}

	if query.StartTime != nil && query.EndTime != nil && query.StartTime.After(*query.EndTime) {
		return fmt.Errorf("start_time cannot be after end_time")
	}

	if query.CreatedAfter != nil && query.CreatedBefore != nil && query.CreatedAfter.After(*query.CreatedBefore) {
		return fmt.Errorf("created_after cannot be after created_before")
	}

	if query.UpdatedAfter != nil && query.UpdatedBefore != nil && query.UpdatedAfter.After(*query.UpdatedBefore) {
		return fmt.Errorf("updated_after cannot be after updated_before")
	}

	return nil
}

// 辅助函数

// GetRankTypeByString 根据字符串获取排行榜类型
func GetRankTypeByString(s string) (RankType, error) {
	switch s {
	case "level":
		return RankTypeLevel, nil
	case "power":
		return RankTypePower, nil
	case "wealth":
		return RankTypeWealth, nil
	case "achievement":
		return RankTypeAchievement, nil
	case "pet":
		return RankTypePet, nil
	case "guild":
		return RankTypeGuild, nil
	case "arena":
		return RankTypeArena, nil
	case "dungeon":
		return RankTypeDungeon, nil
	case "activity":
		return RankTypeActivity, nil
	case "custom":
		return RankTypeCustom, nil
	default:
		return 0, fmt.Errorf("unknown rank type: %s", s)
	}
}

// GetRankCategoryByString 根据字符串获取排行榜分类
func GetRankCategoryByString(s string) (RankCategory, error) {
	switch s {
	case "player":
		return RankCategoryPlayer, nil
	case "guild":
		return RankCategoryGuild, nil
	case "server":
		return RankCategoryServer, nil
	case "global":
		return RankCategoryGlobal, nil
	case "event":
		return RankCategoryEvent, nil
	case "season":
		return RankCategorySeason, nil
	default:
		return 0, fmt.Errorf("unknown rank category: %s", s)
	}
}

// GetSortTypeByString 根据字符串获取排序类型
func GetSortTypeByString(s string) (SortType, error) {
	switch s {
	case "desc", "descending":
		return SortTypeDescending, nil
	case "asc", "ascending":
		return SortTypeAscending, nil
	default:
		return 0, fmt.Errorf("unknown sort type: %s", s)
	}
}

// GetRankPeriodByString 根据字符串获取排行榜周期
func GetRankPeriodByString(s string) (RankPeriod, error) {
	switch s {
	case "permanent":
		return RankPeriodPermanent, nil
	case "daily":
		return RankPeriodDaily, nil
	case "weekly":
		return RankPeriodWeekly, nil
	case "monthly":
		return RankPeriodMonthly, nil
	case "seasonal":
		return RankPeriodSeasonal, nil
	case "event":
		return RankPeriodEvent, nil
	case "custom":
		return RankPeriodCustom, nil
	default:
		return 0, fmt.Errorf("unknown rank period: %s", s)
	}
}

// GetRankStatusByString 根据字符串获取排行榜状态
func GetRankStatusByString(s string) (RankStatus, error) {
	switch s {
	case "active":
		return RankStatusActive, nil
	case "inactive":
		return RankStatusInactive, nil
	case "paused":
		return RankStatusPaused, nil
	case "expired":
		return RankStatusExpired, nil
	case "maintenance":
		return RankStatusMaintenance, nil
	case "archived":
		return RankStatusArchived, nil
	default:
		return 0, fmt.Errorf("unknown rank status: %s", s)
	}
}
