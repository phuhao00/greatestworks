package sacred

import (
	"math"
	"time"
)

// SacredPlaceRepository 圣地仓储接口
type SacredPlaceRepository interface {
	// 基础CRUD操作
	Save(sacredPlace *SacredPlaceAggregate) error
	FindByID(id string) (*SacredPlaceAggregate, error)
	FindByOwner(ownerID string) ([]*SacredPlaceAggregate, error)
	Update(sacredPlace *SacredPlaceAggregate) error
	Delete(id string) error

	// 查询操作
	FindByStatus(status SacredStatus) ([]*SacredPlaceAggregate, error)
	FindByLevel(minLevel, maxLevel int) ([]*SacredPlaceAggregate, error)
	FindByName(name string) ([]*SacredPlaceAggregate, error)
	FindActive() ([]*SacredPlaceAggregate, error)

	// 分页查询
	FindWithPagination(query *SacredPlaceQuery) (*SacredPlacePageResult, error)

	// 统计操作
	Count() (int64, error)
	CountByStatus(status SacredStatus) (int64, error)
	CountByOwner(ownerID string) (int64, error)

	// 批量操作
	SaveBatch(sacredPlaces []*SacredPlaceAggregate) error
	DeleteBatch(ids []string) error

	// 高级查询
	FindNearby(location *Location, radius float64) ([]*SacredPlaceAggregate, error)
	FindTopByLevel(limit int) ([]*SacredPlaceAggregate, error)
	FindRecentlyActive(since time.Time) ([]*SacredPlaceAggregate, error)
}

// ChallengeRepository 挑战仓储接口
type ChallengeRepository interface {
	// 基础CRUD操作
	Save(challenge *Challenge) error
	FindByID(id string) (*Challenge, error)
	FindBySacredPlace(sacredPlaceID string) ([]*Challenge, error)
	Update(challenge *Challenge) error
	Delete(id string) error

	// 查询操作
	FindByType(challengeType ChallengeType) ([]*Challenge, error)
	FindByDifficulty(difficulty ChallengeDifficulty) ([]*Challenge, error)
	FindByStatus(status ChallengeStatus) ([]*Challenge, error)
	FindAvailable() ([]*Challenge, error)
	FindCompleted(playerID string) ([]*Challenge, error)

	// 分页查询
	FindWithPagination(query *ChallengeQuery) (*ChallengePageResult, error)

	// 统计操作
	Count() (int64, error)
	CountByType(challengeType ChallengeType) (int64, error)
	CountByDifficulty(difficulty ChallengeDifficulty) (int64, error)
	CountCompleted(playerID string) (int64, error)

	// 参与者相关
	FindParticipants(challengeID string) ([]*ChallengeParticipant, error)
	SaveParticipant(participant *ChallengeParticipant) error
	UpdateParticipant(participant *ChallengeParticipant) error

	// 排行榜
	GetLeaderboard(challengeID string, limit int) ([]*ChallengeParticipant, error)
	GetPlayerRanking(challengeID, playerID string) (int, error)
}

// BlessingRepository 祝福仓储接口
type BlessingRepository interface {
	// 基础CRUD操作
	Save(blessing *Blessing) error
	FindByID(id string) (*Blessing, error)
	FindBySacredPlace(sacredPlaceID string) ([]*Blessing, error)
	Update(blessing *Blessing) error
	Delete(id string) error

	// 查询操作
	FindByType(blessingType BlessingType) ([]*Blessing, error)
	FindByStatus(status BlessingStatus) ([]*Blessing, error)
	FindAvailable() ([]*Blessing, error)
	FindActive(playerID string) ([]*Blessing, error)

	// 分页查询
	FindWithPagination(query *BlessingQuery) (*BlessingPageResult, error)

	// 统计操作
	Count() (int64, error)
	CountByType(blessingType BlessingType) (int64, error)
	CountActive(playerID string) (int64, error)

	// 效果相关
	SaveBlessingEffect(effect *BlessingEffect) error
	FindEffectsByPlayer(playerID string) ([]*BlessingEffect, error)
	DeleteExpiredEffects() error
}

// RelicRepository 圣物仓储接口
type RelicRepository interface {
	// 基础CRUD操作
	Save(relic *SacredRelic) error
	FindByID(id string) (*SacredRelic, error)
	FindByOwner(ownerID string) ([]*SacredRelic, error)
	Update(relic *SacredRelic) error
	Delete(id string) error

	// 查询操作
	FindByType(relicType RelicType) ([]*SacredRelic, error)
	FindByRarity(rarity RelicRarity) ([]*SacredRelic, error)
	FindByLevel(minLevel, maxLevel int) ([]*SacredRelic, error)

	// 分页查询
	FindWithPagination(query *RelicQuery) (*RelicPageResult, error)

	// 统计操作
	Count() (int64, error)
	CountByType(relicType RelicType) (int64, error)
	CountByRarity(rarity RelicRarity) (int64, error)
	CountByOwner(ownerID string) (int64, error)

	// 高级查询
	FindTopByPower(limit int) ([]*SacredRelic, error)
	FindByAttributes(attributes map[string]float64) ([]*SacredRelic, error)
	FindUpgradeable(ownerID string) ([]*SacredRelic, error)
}

// SacredStatisticsRepository 圣地统计仓储接口
type SacredStatisticsRepository interface {
	// 保存统计数据
	SaveStatistics(stats *SacredStatistics) error
	UpdateStatistics(stats *SacredStatistics) error

	// 查询统计数据
	FindStatistics(sacredID string) (*SacredStatistics, error)
	FindStatisticsByOwner(ownerID string) ([]*SacredStatistics, error)

	// 排行榜统计
	GetLevelRanking(limit int) ([]*SacredStatistics, error)
	GetExperienceRanking(limit int) ([]*SacredStatistics, error)
	GetChallengeRanking(limit int) ([]*SacredStatistics, error)

	// 趋势分析
	GetLevelTrend(sacredID string, days int) ([]*LevelTrendData, error)
	GetActivityTrend(sacredID string, days int) ([]*ActivityTrendData, error)

	// 聚合统计
	GetGlobalStatistics() (*GlobalSacredStatistics, error)
	GetOwnerStatistics(ownerID string) (*OwnerSacredStatistics, error)
}

// 查询条件结构体

// SacredPlaceQuery 圣地查询条件
type SacredPlaceQuery struct {
	OwnerID       string
	Name          string
	Status        *SacredStatus
	MinLevel      *int
	MaxLevel      *int
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	UpdatedAfter  *time.Time
	UpdatedBefore *time.Time
	Location      *Location
	Radius        *float64
	OrderBy       string
	OrderDesc     bool
	Offset        int
	Limit         int
}

// ChallengeQuery 挑战查询条件
type ChallengeQuery struct {
	SacredPlaceID string
	Type          *ChallengeType
	Difficulty    *ChallengeDifficulty
	Status        *ChallengeStatus
	MinLevel      *int
	MaxLevel      *int
	PlayerID      string
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	OrderBy       string
	OrderDesc     bool
	Offset        int
	Limit         int
}

// BlessingQuery 祝福查询条件
type BlessingQuery struct {
	SacredPlaceID string
	Type          *BlessingType
	Status        *BlessingStatus
	PlayerID      string
	ActiveOnly    bool
	AvailableOnly bool
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	OrderBy       string
	OrderDesc     bool
	Offset        int
	Limit         int
}

// RelicQuery 圣物查询条件
type RelicQuery struct {
	OwnerID        string
	Type           *RelicType
	Rarity         *RelicRarity
	MinLevel       *int
	MaxLevel       *int
	MinPower       *float64
	MaxPower       *float64
	Attributes     map[string]float64
	ObtainedAfter  *time.Time
	ObtainedBefore *time.Time
	OrderBy        string
	OrderDesc      bool
	Offset         int
	Limit          int
}

// 分页结果结构体

// SacredPlacePageResult 圣地分页结果
type SacredPlacePageResult struct {
	Items   []*SacredPlaceAggregate
	Total   int64
	Offset  int
	Limit   int
	HasMore bool
}

// ChallengePageResult 挑战分页结果
type ChallengePageResult struct {
	Items   []*Challenge
	Total   int64
	Offset  int
	Limit   int
	HasMore bool
}

// BlessingPageResult 祝福分页结果
type BlessingPageResult struct {
	Items   []*Blessing
	Total   int64
	Offset  int
	Limit   int
	HasMore bool
}

// RelicPageResult 圣物分页结果
type RelicPageResult struct {
	Items   []*SacredRelic
	Total   int64
	Offset  int
	Limit   int
	HasMore bool
}

// 统计数据结构体

// LevelTrendData 等级趋势数据
type LevelTrendData struct {
	Date       time.Time
	Level      int
	Experience int
	Growth     int
}

// ActivityTrendData 活动趋势数据
type ActivityTrendData struct {
	Date                time.Time
	ChallengesStarted   int
	ChallengesCompleted int
	BlessingsActivated  int
	PlayersActive       int
}

// GlobalSacredStatistics 全局圣地统计
type GlobalSacredStatistics struct {
	TotalSacredPlaces    int64
	ActiveSacredPlaces   int64
	TotalChallenges      int64
	CompletedChallenges  int64
	TotalBlessings       int64
	ActiveBlessings      int64
	TotalRelics          int64
	AverageLevel         float64
	TopLevel             int
	MostActiveSacred     string
	MostPopularChallenge ChallengeType
	MostUsedBlessing     BlessingType
	UpdatedAt            time.Time
}

// OwnerSacredStatistics 拥有者圣地统计
type OwnerSacredStatistics struct {
	OwnerID             string
	TotalSacredPlaces   int
	ActiveSacredPlaces  int
	TotalLevel          int
	AverageLevel        float64
	HighestLevel        int
	TotalChallenges     int
	CompletedChallenges int
	SuccessRate         float64
	TotalBlessings      int
	ActiveBlessings     int
	TotalRelics         int
	TotalExperience     int64
	Ranking             int
	LastActiveAt        time.Time
	CreatedAt           time.Time
}

// Location 位置信息
type Location struct {
	Latitude  float64
	Longitude float64
	Region    string
	Zone      string
}

// NewLocation 创建位置
func NewLocation(latitude, longitude float64, region, zone string) *Location {
	return &Location{
		Latitude:  latitude,
		Longitude: longitude,
		Region:    region,
		Zone:      zone,
	}
}

// DistanceTo 计算到另一个位置的距离
func (l *Location) DistanceTo(other *Location) float64 {
	// 简化的距离计算（实际应用中可能需要更精确的地理计算）
	dx := l.Latitude - other.Latitude
	dy := l.Longitude - other.Longitude
	return math.Sqrt(dx*dx + dy*dy)
}

// IsWithinRadius 检查是否在指定半径内
func (l *Location) IsWithinRadius(center *Location, radius float64) bool {
	return l.DistanceTo(center) <= radius
}

// 缓存接口

// SacredCacheRepository 圣地缓存仓储接口
type SacredCacheRepository interface {
	// 圣地缓存
	SetSacredPlace(id string, sacredPlace *SacredPlaceAggregate, ttl time.Duration) error
	GetSacredPlace(id string) (*SacredPlaceAggregate, error)
	DeleteSacredPlace(id string) error

	// 挑战缓存
	SetChallenge(id string, challenge *Challenge, ttl time.Duration) error
	GetChallenge(id string) (*Challenge, error)
	DeleteChallenge(id string) error

	// 祝福缓存
	SetBlessing(id string, blessing *Blessing, ttl time.Duration) error
	GetBlessing(id string) (*Blessing, error)
	DeleteBlessing(id string) error

	// 圣物缓存
	SetRelic(id string, relic *SacredRelic, ttl time.Duration) error
	GetRelic(id string) (*SacredRelic, error)
	DeleteRelic(id string) error

	// 统计缓存
	SetStatistics(key string, stats interface{}, ttl time.Duration) error
	GetStatistics(key string, result interface{}) error
	DeleteStatistics(key string) error

	// 排行榜缓存
	SetRanking(key string, ranking interface{}, ttl time.Duration) error
	GetRanking(key string, result interface{}) error
	DeleteRanking(key string) error

	// 批量操作
	SetBatch(items map[string]interface{}, ttl time.Duration) error
	GetBatch(keys []string) (map[string]interface{}, error)
	DeleteBatch(keys []string) error

	// 缓存管理
	Clear() error
	Exists(key string) (bool, error)
	SetTTL(key string, ttl time.Duration) error
	GetTTL(key string) (time.Duration, error)
}

// 事务接口

// SacredTransactionRepository 圣地事务仓储接口
type SacredTransactionRepository interface {
	// 事务管理
	BeginTransaction() (SacredTransaction, error)
	CommitTransaction(tx SacredTransaction) error
	RollbackTransaction(tx SacredTransaction) error

	// 在事务中执行操作
	ExecuteInTransaction(fn func(tx SacredTransaction) error) error
}

// SacredTransaction 圣地事务接口
type SacredTransaction interface {
	// 圣地操作
	SaveSacredPlace(sacredPlace *SacredPlaceAggregate) error
	UpdateSacredPlace(sacredPlace *SacredPlaceAggregate) error
	DeleteSacredPlace(id string) error

	// 挑战操作
	SaveChallenge(challenge *Challenge) error
	UpdateChallenge(challenge *Challenge) error
	DeleteChallenge(id string) error

	// 祝福操作
	SaveBlessing(blessing *Blessing) error
	UpdateBlessing(blessing *Blessing) error
	DeleteBlessing(id string) error

	// 圣物操作
	SaveRelic(relic *SacredRelic) error
	UpdateRelic(relic *SacredRelic) error
	DeleteRelic(id string) error

	// 统计操作
	UpdateStatistics(stats *SacredStatistics) error

	// 事务状态
	IsActive() bool
	GetID() string
}

// 仓储工厂接口

// SacredRepositoryFactory 圣地仓储工厂接口
type SacredRepositoryFactory interface {
	// 创建仓储实例
	CreateSacredPlaceRepository() SacredPlaceRepository
	CreateChallengeRepository() ChallengeRepository
	CreateBlessingRepository() BlessingRepository
	CreateRelicRepository() RelicRepository
	CreateStatisticsRepository() SacredStatisticsRepository
	CreateCacheRepository() SacredCacheRepository
	CreateTransactionRepository() SacredTransactionRepository

	// 健康检查
	HealthCheck() error

	// 关闭连接
	Close() error
}
