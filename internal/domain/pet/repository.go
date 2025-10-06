package pet

import (
	"time"
)

// PetRepository 宠物仓储接口
type PetRepository interface {
	// 基础CRUD操作
	Save(pet *PetAggregate) error
	FindByID(id string) (*PetAggregate, error)
	FindByPlayer(playerID string) ([]*PetAggregate, error)
	FindByPlayerAndCategory(playerID string, category PetCategory) ([]*PetAggregate, error)
	Update(pet *PetAggregate) error
	Delete(id string) error

	// 分页查询
	FindWithPagination(query *PetQuery) (*PetPageResult, error)

	// 统计操作
	Count() (int64, error)
	CountByPlayer(playerID string) (int64, error)
	CountByCategory(category PetCategory) (int64, error)

	// 状态查询
	FindByState(state PetState) ([]*PetAggregate, error)
	FindActiveByPlayer(playerID string) ([]*PetAggregate, error)
	FindDeadPets() ([]*PetAggregate, error)

	// 等级和星级查询
	FindByLevelRange(minLevel, maxLevel uint32) ([]*PetAggregate, error)
	FindByStarRange(minStar, maxStar uint32) ([]*PetAggregate, error)

	// 批量操作
	SaveBatch(pets []*PetAggregate) error
	DeleteBatch(ids []string) error

	// 高级查询
	FindTopPetsByPower(limit int) ([]*PetAggregate, error)
	FindRecentlyCreated(duration time.Duration) ([]*PetAggregate, error)
}

// PetFragmentRepository 宠物碎片仓储接口
type PetFragmentRepository interface {
	// 基础CRUD操作
	Save(fragment *PetFragment) error
	FindByID(id string) (*PetFragment, error)
	FindByPlayer(playerID string) ([]*PetFragment, error)
	FindByPlayerAndFragmentID(playerID string, fragmentID uint32) (*PetFragment, error)
	Update(fragment *PetFragment) error
	Delete(id string) error

	// 分页查询
	FindWithPagination(query *FragmentQuery) (*FragmentPageResult, error)

	// 统计操作
	Count() (int64, error)
	CountByPlayer(playerID string) (int64, error)
	GetTotalQuantityByPlayer(playerID string, fragmentID uint32) (uint64, error)

	// 碎片相关查询
	FindByRelatedPet(relatedPetID uint32) ([]*PetFragment, error)
	FindSufficientFragments(playerID string, fragmentID uint32, requiredQuantity uint64) ([]*PetFragment, error)

	// 批量操作
	SaveBatch(fragments []*PetFragment) error
	UpdateBatch(fragments []*PetFragment) error
}

// PetSkinRepository 宠物皮肤仓储接口
type PetSkinRepository interface {
	// 基础CRUD操作
	Save(skin *PetSkin) error
	FindByID(id string) (*PetSkin, error)
	FindBySkinID(skinID string) (*PetSkin, error)
	Update(skin *PetSkin) error
	Delete(id string) error

	// 分页查询
	FindWithPagination(query *SkinQuery) (*SkinPageResult, error)

	// 统计操作
	Count() (int64, error)
	CountByRarity(rarity PetRarity) (int64, error)

	// 皮肤相关查询
	FindByRarity(rarity PetRarity) ([]*PetSkin, error)
	FindUnlockedSkins() ([]*PetSkin, error)
	FindEquippedSkins() ([]*PetSkin, error)

	// 批量操作
	SaveBatch(skins []*PetSkin) error
	UpdateBatch(skins []*PetSkin) error
}

// PetSkillRepository 宠物技能仓储接口
type PetSkillRepository interface {
	// 基础CRUD操作
	Save(skill *PetSkill) error
	FindByID(id string) (*PetSkill, error)
	FindBySkillID(skillID string) (*PetSkill, error)
	Update(skill *PetSkill) error
	Delete(id string) error

	// 分页查询
	FindWithPagination(query *SkillQuery) (*SkillPageResult, error)

	// 统计操作
	Count() (int64, error)
	CountByType(skillType SkillType) (int64, error)

	// 技能相关查询
	FindByType(skillType SkillType) ([]*PetSkill, error)
	FindByLevelRange(minLevel, maxLevel uint32) ([]*PetSkill, error)
	FindReadySkills() ([]*PetSkill, error)

	// 批量操作
	SaveBatch(skills []*PetSkill) error
	UpdateBatch(skills []*PetSkill) error
}

// PetBondsRepository 宠物羁绊仓储接口
type PetBondsRepository interface {
	// 基础CRUD操作
	Save(bonds *PetBonds) error
	FindByID(id string) (*PetBonds, error)
	Update(bonds *PetBonds) error
	Delete(id string) error

	// 羁绊相关查询
	FindActiveBonds() ([]*PetBonds, error)
	FindByBondID(bondID string) ([]*PetBonds, error)

	// 统计操作
	Count() (int64, error)
	CountActiveBonds() (int64, error)
}

// PetPictorialRepository 宠物图鉴仓储接口
type PetPictorialRepository interface {
	// 基础CRUD操作
	Save(pictorial *PetPictorial) error
	FindByID(id string) (*PetPictorial, error)
	FindByPlayer(playerID string) ([]*PetPictorial, error)
	FindByPlayerAndPetConfig(playerID string, petConfigID uint32) (*PetPictorial, error)
	Update(pictorial *PetPictorial) error
	Delete(id string) error

	// 分页查询
	FindWithPagination(query *PictorialQuery) (*PictorialPageResult, error)

	// 统计操作
	Count() (int64, error)
	CountByPlayer(playerID string) (int64, error)
	CountUnlockedByPlayer(playerID string) (int64, error)

	// 图鉴相关查询
	FindUnlockedByPlayer(playerID string) ([]*PetPictorial, error)
	FindRecentlyUnlocked(duration time.Duration) ([]*PetPictorial, error)

	// 批量操作
	SaveBatch(pictorials []*PetPictorial) error
	UpdateBatch(pictorials []*PetPictorial) error
}

// PetStatisticsRepository 宠物统计仓储接口
type PetStatisticsRepository interface {
	// 保存统计数据
	SaveStatistics(stats *PetStatistics) error
	UpdateStatistics(stats *PetStatistics) error

	// 查询统计数据
	FindStatistics(playerID string) (*PetStatistics, error)
	FindStatisticsByCategory(category PetCategory) ([]*PetStatistics, error)

	// 全局统计
	GetGlobalStatistics() (*GlobalPetStatistics, error)
	GetCategoryStatistics(category PetCategory) (*CategoryPetStatistics, error)

	// 趋势分析
	GetLevelTrend(playerID string, days int) ([]*LevelTrendData, error)
	GetPowerTrend(playerID string, days int) ([]*PowerTrendData, error)

	// 排行榜数据
	GetTopPlayersByPetCount(limit int) ([]*PlayerPetRanking, error)
	GetTopPlayersByTotalPower(limit int) ([]*PlayerPetRanking, error)
}

// 查询条件结构体

// PetQuery 宠物查询条件
type PetQuery struct {
	PlayerID      string
	Name          string
	Category      *PetCategory
	State         *PetState
	MinLevel      *uint32
	MaxLevel      *uint32
	MinStar       *uint32
	MaxStar       *uint32
	MinPower      *int64
	MaxPower      *int64
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	UpdatedAfter  *time.Time
	UpdatedBefore *time.Time
	OrderBy       string
	OrderDesc     bool
	Offset        int
	Limit         int
}

// FragmentQuery 碎片查询条件
type FragmentQuery struct {
	PlayerID      string
	FragmentID    *uint32
	RelatedPetID  *uint32
	MinQuantity   *uint64
	MaxQuantity   *uint64
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	OrderBy       string
	OrderDesc     bool
	Offset        int
	Limit         int
}

// SkinQuery 皮肤查询条件
type SkinQuery struct {
	SkinID        string
	Name          string
	Rarity        *PetRarity
	Unlocked      *bool
	Equipped      *bool
	MinPowerBonus *int64
	MaxPowerBonus *int64
	OrderBy       string
	OrderDesc     bool
	Offset        int
	Limit         int
}

// SkillQuery 技能查询条件
type SkillQuery struct {
	SkillID   string
	Name      string
	SkillType *SkillType
	MinLevel  *uint32
	MaxLevel  *uint32
	MinDamage *int64
	MaxDamage *int64
	Ready     *bool
	OrderBy   string
	OrderDesc bool
	Offset    int
	Limit     int
}

// PictorialQuery 图鉴查询条件
type PictorialQuery struct {
	PlayerID      string
	PetConfigID   *uint32
	Unlocked      *bool
	MinLevel      *uint32
	MaxLevel      *uint32
	MinStar       *uint32
	MaxStar       *uint32
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	OrderBy       string
	OrderDesc     bool
	Offset        int
	Limit         int
}

// 分页结果结构体

// PetPageResult 宠物分页结果
type PetPageResult struct {
	Items   []*PetAggregate
	Total   int64
	Offset  int
	Limit   int
	HasMore bool
}

// FragmentPageResult 碎片分页结果
type FragmentPageResult struct {
	Items   []*PetFragment
	Total   int64
	Offset  int
	Limit   int
	HasMore bool
}

// SkinPageResult 皮肤分页结果
type SkinPageResult struct {
	Items   []*PetSkin
	Total   int64
	Offset  int
	Limit   int
	HasMore bool
}

// SkillPageResult 技能分页结果
type SkillPageResult struct {
	Items   []*PetSkill
	Total   int64
	Offset  int
	Limit   int
	HasMore bool
}

// PictorialPageResult 图鉴分页结果
type PictorialPageResult struct {
	Items   []*PetPictorial
	Total   int64
	Offset  int
	Limit   int
	HasMore bool
}

// 统计数据结构体

// PetStatistics 宠物统计
type PetStatistics struct {
	PlayerID          string
	TotalPets         int64
	AlivePets         int64
	DeadPets          int64
	MaxLevel          uint32
	MaxStar           uint32
	TotalPower        int64
	AveragePower      float64
	CategoryStats     map[PetCategory]*CategoryStats
	FavoritePet       string
	MostUsedSkill     string
	TotalTrainingTime time.Duration
	TotalFeedingCount int64
	LastActivityTime  time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// CategoryStats 类别统计
type CategoryStats struct {
	Category     PetCategory
	Count        int64
	TotalPower   int64
	AveragePower float64
	MaxLevel     uint32
	MaxStar      uint32
}

// GlobalPetStatistics 全局宠物统计
type GlobalPetStatistics struct {
	TotalPets            int64
	TotalPlayers         int64
	AveragePetsPerPlayer float64
	CategoryDistribution map[PetCategory]int64
	RarityDistribution   map[PetRarity]int64
	MostPopularCategory  PetCategory
	MostPopularRarity    PetRarity
	TotalPower           int64
	AveragePower         float64
	TopPet               string
	UpdatedAt            time.Time
}

// CategoryPetStatistics 类别宠物统计
type CategoryPetStatistics struct {
	Category         PetCategory
	TotalCount       int64
	ActiveCount      int64
	AverageLevel     float64
	AverageStar      float64
	TotalPower       int64
	AveragePower     float64
	TopPet           string
	MostActivePlayer string
	UpdatedAt        time.Time
}

// LevelTrendData 等级趋势数据
type LevelTrendData struct {
	Date         time.Time
	AverageLevel float64
	MaxLevel     uint32
	LevelUps     int64
}

// PowerTrendData 战力趋势数据
type PowerTrendData struct {
	Date         time.Time
	TotalPower   int64
	AveragePower float64
	MaxPower     int64
}

// PlayerPetRanking 玩家宠物排行
type PlayerPetRanking struct {
	PlayerID     string
	PlayerName   string
	PetCount     int64
	TotalPower   int64
	AveragePower float64
	TopPetName   string
	Rank         int
}

// 缓存接口

// PetCacheRepository 宠物缓存仓储接口
type PetCacheRepository interface {
	// 宠物缓存
	SetPet(id string, pet *PetAggregate, ttl time.Duration) error
	GetPet(id string) (*PetAggregate, error)
	DeletePet(id string) error

	// 碎片缓存
	SetFragment(id string, fragment *PetFragment, ttl time.Duration) error
	GetFragment(id string) (*PetFragment, error)
	DeleteFragment(id string) error

	// 皮肤缓存
	SetSkin(id string, skin *PetSkin, ttl time.Duration) error
	GetSkin(id string) (*PetSkin, error)
	DeleteSkin(id string) error

	// 技能缓存
	SetSkill(id string, skill *PetSkill, ttl time.Duration) error
	GetSkill(id string) (*PetSkill, error)
	DeleteSkill(id string) error

	// 图鉴缓存
	SetPictorial(id string, pictorial *PetPictorial, ttl time.Duration) error
	GetPictorial(id string) (*PetPictorial, error)
	DeletePictorial(id string) error

	// 统计缓存
	SetStatistics(key string, stats interface{}, ttl time.Duration) error
	GetStatistics(key string, result interface{}) error
	DeleteStatistics(key string) error

	// 玩家宠物列表缓存
	SetPlayerPets(playerID string, pets []*PetAggregate, ttl time.Duration) error
	GetPlayerPets(playerID string) ([]*PetAggregate, error)
	DeletePlayerPets(playerID string) error

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

// PetTransactionRepository 宠物事务仓储接口
type PetTransactionRepository interface {
	// 事务管理
	BeginTransaction() (PetTransaction, error)
	CommitTransaction(tx PetTransaction) error
	RollbackTransaction(tx PetTransaction) error

	// 在事务中执行操作
	ExecuteInTransaction(fn func(tx PetTransaction) error) error
}

// PetTransaction 宠物事务接口
type PetTransaction interface {
	// 宠物操作
	SavePet(pet *PetAggregate) error
	UpdatePet(pet *PetAggregate) error
	DeletePet(id string) error

	// 碎片操作
	SaveFragment(fragment *PetFragment) error
	UpdateFragment(fragment *PetFragment) error
	DeleteFragment(id string) error

	// 皮肤操作
	SaveSkin(skin *PetSkin) error
	UpdateSkin(skin *PetSkin) error
	DeleteSkin(id string) error

	// 技能操作
	SaveSkill(skill *PetSkill) error
	UpdateSkill(skill *PetSkill) error
	DeleteSkill(id string) error

	// 羁绊操作
	SaveBonds(bonds *PetBonds) error
	UpdateBonds(bonds *PetBonds) error
	DeleteBonds(id string) error

	// 图鉴操作
	SavePictorial(pictorial *PetPictorial) error
	UpdatePictorial(pictorial *PetPictorial) error
	DeletePictorial(id string) error

	// 统计操作
	UpdateStatistics(stats *PetStatistics) error

	// 事务状态
	IsActive() bool
	GetID() string
}

// 仓储工厂接口

// PetRepositoryFactory 宠物仓储工厂接口
type PetRepositoryFactory interface {
	// 创建仓储实例
	CreatePetRepository() PetRepository
	CreateFragmentRepository() PetFragmentRepository
	CreateSkinRepository() PetSkinRepository
	CreateSkillRepository() PetSkillRepository
	CreateBondsRepository() PetBondsRepository
	CreatePictorialRepository() PetPictorialRepository
	CreateStatisticsRepository() PetStatisticsRepository
	CreateCacheRepository() PetCacheRepository
	CreateTransactionRepository() PetTransactionRepository

	// 健康检查
	HealthCheck() error

	// 关闭连接
	Close() error
}

// 搜索接口

// PetSearchRepository 宠物搜索仓储接口
type PetSearchRepository interface {
	// 全文搜索
	SearchPets(query string, filters map[string]interface{}) ([]*PetAggregate, error)
	SearchFragments(query string, filters map[string]interface{}) ([]*PetFragment, error)
	SearchSkins(query string, filters map[string]interface{}) ([]*PetSkin, error)
	SearchSkills(query string, filters map[string]interface{}) ([]*PetSkill, error)

	// 智能推荐
	RecommendPets(playerID string, category PetCategory, limit int) ([]*PetAggregate, error)
	RecommendSkills(petID string, limit int) ([]*PetSkill, error)
	RecommendSkins(petID string, limit int) ([]*PetSkin, error)

	// 相似度搜索
	FindSimilarPets(petID string, limit int) ([]*PetAggregate, error)

	// 索引管理
	RebuildIndex() error
	UpdateIndex(entity *PetAggregate) error
}
