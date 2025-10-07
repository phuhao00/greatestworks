package building

import (
	"context"
	"fmt"
	"time"
)

// BuildingRepository 建筑仓储接口
type BuildingRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, building *BuildingAggregate) error
	FindByID(ctx context.Context, id string) (*BuildingAggregate, error)
	FindByIDs(ctx context.Context, ids []string) ([]*BuildingAggregate, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)

	// 查询操作
	FindByOwner(ctx context.Context, ownerID uint64) ([]*BuildingAggregate, error)
	FindByType(ctx context.Context, buildingType BuildingType) ([]*BuildingAggregate, error)
	FindByCategory(ctx context.Context, category BuildingCategory) ([]*BuildingAggregate, error)
	FindByStatus(ctx context.Context, status BuildingStatus) ([]*BuildingAggregate, error)
	FindByPosition(ctx context.Context, position *Position) (*BuildingAggregate, error)
	FindByPlayerAndPosition(ctx context.Context, playerID uint64, position *Position) (*BuildingAggregate, error)
	FindByArea(ctx context.Context, area *Area) ([]*BuildingAggregate, error)
	FindByQuery(ctx context.Context, query *BuildingQuery) ([]*BuildingAggregate, int64, error)

	// 统计操作
	Count(ctx context.Context) (int64, error)
	CountByOwner(ctx context.Context, ownerID uint64) (int64, error)
	CountByType(ctx context.Context, buildingType BuildingType) (int64, error)
	CountByCategory(ctx context.Context, category BuildingCategory) (int64, error)
	CountByStatus(ctx context.Context, status BuildingStatus) (int64, error)
	GetStatistics(ctx context.Context, ownerID uint64) (*BuildingStatistics, error)

	// 批量操作
	SaveAll(ctx context.Context, buildings []*BuildingAggregate) error
	DeleteAll(ctx context.Context, ids []string) error
	UpdateStatus(ctx context.Context, ids []string, status BuildingStatus) error
	UpdateHealth(ctx context.Context, id string, health float64) error
	UpdateLevel(ctx context.Context, id string, level int32) error
}

// ConstructionRepository 建造仓储接口
type ConstructionRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, construction *ConstructionInfo) error
	FindByID(ctx context.Context, id string) (*ConstructionInfo, error)
	FindByIDs(ctx context.Context, ids []string) ([]*ConstructionInfo, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)

	// 查询操作
	FindByBuildingID(ctx context.Context, buildingID string) ([]*ConstructionInfo, error)
	FindByStatus(ctx context.Context, status ConstructionStatus) ([]*ConstructionInfo, error)
	FindByWorker(ctx context.Context, workerID uint64) ([]*ConstructionInfo, error)
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*ConstructionInfo, error)
	FindByQuery(ctx context.Context, query *ConstructionQuery) ([]*ConstructionInfo, int64, error)

	// 统计操作
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status ConstructionStatus) (int64, error)
	CountByBuildingID(ctx context.Context, buildingID string) (int64, error)
	GetStatistics(ctx context.Context, buildingID string) (*ConstructionStatistics, error)

	// 批量操作
	SaveAll(ctx context.Context, constructions []*ConstructionInfo) error
	DeleteAll(ctx context.Context, ids []string) error
	UpdateStatus(ctx context.Context, ids []string, status ConstructionStatus) error
	UpdateProgress(ctx context.Context, id string, progress float64) error
}

// UpgradeRepository 升级仓储接口
type UpgradeRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, upgrade *UpgradeInfo) error
	FindByID(ctx context.Context, id string) (*UpgradeInfo, error)
	FindByIDs(ctx context.Context, ids []string) ([]*UpgradeInfo, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)

	// 查询操作
	FindByBuildingID(ctx context.Context, buildingID string) ([]*UpgradeInfo, error)
	FindByStatus(ctx context.Context, status UpgradeStatus) ([]*UpgradeInfo, error)
	FindByLevel(ctx context.Context, fromLevel, toLevel int32) ([]*UpgradeInfo, error)
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*UpgradeInfo, error)
	FindByQuery(ctx context.Context, query *UpgradeQuery) ([]*UpgradeInfo, int64, error)

	// 统计操作
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status UpgradeStatus) (int64, error)
	CountByBuildingID(ctx context.Context, buildingID string) (int64, error)
	GetStatistics(ctx context.Context, buildingID string) (*UpgradeStatistics, error)

	// 批量操作
	SaveAll(ctx context.Context, upgrades []*UpgradeInfo) error
	DeleteAll(ctx context.Context, ids []string) error
	UpdateStatus(ctx context.Context, ids []string, status UpgradeStatus) error
	UpdateProgress(ctx context.Context, id string, progress float64) error
}

// BlueprintRepository 蓝图仓储接口
type BlueprintRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, blueprint *Blueprint) error
	FindByID(ctx context.Context, id string) (*Blueprint, error)
	FindByIDs(ctx context.Context, ids []string) ([]*Blueprint, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)

	// 查询操作
	FindByName(ctx context.Context, name string) ([]*Blueprint, error)
	FindByAuthor(ctx context.Context, author string) ([]*Blueprint, error)
	FindByCategory(ctx context.Context, category BuildingCategory) ([]*Blueprint, error)
	FindByDifficulty(ctx context.Context, minDifficulty, maxDifficulty int32) ([]*Blueprint, error)
	FindByTags(ctx context.Context, tags []string) ([]*Blueprint, error)
	FindByQuery(ctx context.Context, query *BlueprintQuery) ([]*Blueprint, int64, error)

	// 统计操作
	Count(ctx context.Context) (int64, error)
	CountByCategory(ctx context.Context, category BuildingCategory) (int64, error)
	CountByAuthor(ctx context.Context, author string) (int64, error)
	GetStatistics(ctx context.Context) (*BlueprintStatistics, error)

	// 批量操作
	SaveAll(ctx context.Context, blueprints []*Blueprint) error
	DeleteAll(ctx context.Context, ids []string) error
}

// 查询条件结构体

// BuildingQuery 建筑查询条件
type BuildingQuery struct {
	// 基础字段
	OwnerID  *uint64           `json:"owner_id,omitempty"`
	Name     *string           `json:"name,omitempty"`
	Type     *BuildingType     `json:"type,omitempty"`
	Category *BuildingCategory `json:"category,omitempty"`
	Status   *BuildingStatus   `json:"status,omitempty"`

	// 等级范围
	MinLevel *int32 `json:"min_level,omitempty"`
	MaxLevel *int32 `json:"max_level,omitempty"`

	// 健康度范围
	MinHealth *float64 `json:"min_health,omitempty"`
	MaxHealth *float64 `json:"max_health,omitempty"`

	// 位置查询
	Position *Position `json:"position,omitempty"`
	Area     *Area     `json:"area,omitempty"`

	// 时间范围
	CreatedAfter  *time.Time `json:"created_after,omitempty"`
	CreatedBefore *time.Time `json:"created_before,omitempty"`
	UpdatedAfter  *time.Time `json:"updated_after,omitempty"`
	UpdatedBefore *time.Time `json:"updated_before,omitempty"`

	// 标签
	Tags []string `json:"tags,omitempty"`

	// 排序
	SortBy    string `json:"sort_by,omitempty"`    // name, type, level, health, created_at, updated_at
	SortOrder string `json:"sort_order,omitempty"` // asc, desc

	// 分页
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}

// ConstructionQuery 建造查询条件
type ConstructionQuery struct {
	// 基础字段
	BuildingID *string             `json:"building_id,omitempty"`
	Status     *ConstructionStatus `json:"status,omitempty"`
	WorkerID   *uint64             `json:"worker_id,omitempty"`

	// 进度范围
	MinProgress *float64 `json:"min_progress,omitempty"`
	MaxProgress *float64 `json:"max_progress,omitempty"`

	// 时间范围
	StartedAfter    *time.Time `json:"started_after,omitempty"`
	StartedBefore   *time.Time `json:"started_before,omitempty"`
	CompletedAfter  *time.Time `json:"completed_after,omitempty"`
	CompletedBefore *time.Time `json:"completed_before,omitempty"`
	CreatedAfter    *time.Time `json:"created_after,omitempty"`
	CreatedBefore   *time.Time `json:"created_before,omitempty"`

	// 排序
	SortBy    string `json:"sort_by,omitempty"`    // progress, started_at, duration, created_at
	SortOrder string `json:"sort_order,omitempty"` // asc, desc

	// 分页
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}

// UpgradeQuery 升级查询条件
type UpgradeQuery struct {
	// 基础字段
	BuildingID *string        `json:"building_id,omitempty"`
	Status     *UpgradeStatus `json:"status,omitempty"`

	// 等级范围
	FromLevel *int32 `json:"from_level,omitempty"`
	ToLevel   *int32 `json:"to_level,omitempty"`

	// 进度范围
	MinProgress *float64 `json:"min_progress,omitempty"`
	MaxProgress *float64 `json:"max_progress,omitempty"`

	// 时间范围
	StartedAfter    *time.Time `json:"started_after,omitempty"`
	StartedBefore   *time.Time `json:"started_before,omitempty"`
	CompletedAfter  *time.Time `json:"completed_after,omitempty"`
	CompletedBefore *time.Time `json:"completed_before,omitempty"`
	CreatedAfter    *time.Time `json:"created_after,omitempty"`
	CreatedBefore   *time.Time `json:"created_before,omitempty"`

	// 排序
	SortBy    string `json:"sort_by,omitempty"`    // from_level, to_level, progress, started_at, created_at
	SortOrder string `json:"sort_order,omitempty"` // asc, desc

	// 分页
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}

// BlueprintQuery 蓝图查询条件
type BlueprintQuery struct {
	// 基础字段
	Name     *string           `json:"name,omitempty"`
	Author   *string           `json:"author,omitempty"`
	Category *BuildingCategory `json:"category,omitempty"`
	Version  *string           `json:"version,omitempty"`

	// 难度范围
	MinDifficulty *int32 `json:"min_difficulty,omitempty"`
	MaxDifficulty *int32 `json:"max_difficulty,omitempty"`

	// 尺寸范围
	MinWidth  *int32 `json:"min_width,omitempty"`
	MaxWidth  *int32 `json:"max_width,omitempty"`
	MinHeight *int32 `json:"min_height,omitempty"`
	MaxHeight *int32 `json:"max_height,omitempty"`
	MinDepth  *int32 `json:"min_depth,omitempty"`
	MaxDepth  *int32 `json:"max_depth,omitempty"`

	// 时间范围
	MinDuration   *time.Duration `json:"min_duration,omitempty"`
	MaxDuration   *time.Duration `json:"max_duration,omitempty"`
	CreatedAfter  *time.Time     `json:"created_after,omitempty"`
	CreatedBefore *time.Time     `json:"created_before,omitempty"`
	UpdatedAfter  *time.Time     `json:"updated_after,omitempty"`
	UpdatedBefore *time.Time     `json:"updated_before,omitempty"`

	// 标签
	Tags []string `json:"tags,omitempty"`

	// 搜索关键词
	Keyword *string `json:"keyword,omitempty"`

	// 排序
	SortBy    string `json:"sort_by,omitempty"`    // name, author, difficulty, duration, created_at
	SortOrder string `json:"sort_order,omitempty"` // asc, desc

	// 分页
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}

// 分页结果结构体

// PaginationResult 分页结果
type PaginationResult struct {
	Total       int64 `json:"total"`
	Page        int   `json:"page"`
	PageSize    int   `json:"page_size"`
	TotalPages  int64 `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}

// 统计数据结构体

// BuildingStatistics 建筑统计数据
type BuildingStatistics struct {
	OwnerID uint64 `json:"owner_id"`

	// 总体统计
	TotalBuildings    int64 `json:"total_buildings"`
	ActiveBuildings   int64 `json:"active_buildings"`
	InactiveBuildings int64 `json:"inactive_buildings"`
	UnderConstruction int64 `json:"under_construction"`
	UnderUpgrade      int64 `json:"under_upgrade"`
	Damaged           int64 `json:"damaged"`
	Destroyed         int64 `json:"destroyed"`

	// 按类型统计
	ByType map[BuildingType]int64 `json:"by_type"`

	// 按分类统计
	ByCategory map[BuildingCategory]int64 `json:"by_category"`

	// 按等级统计
	ByLevel map[int32]int64 `json:"by_level"`

	// 健康度统计
	AverageHealth float64 `json:"average_health"`
	MinHealth     float64 `json:"min_health"`
	MaxHealth     float64 `json:"max_health"`

	// 等级统计
	AverageLevel float64 `json:"average_level"`
	MinLevel     int32   `json:"min_level"`
	MaxLevel     int32   `json:"max_level"`

	// 时间统计
	OldestBuilding time.Time `json:"oldest_building"`
	NewestBuilding time.Time `json:"newest_building"`

	// 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// ConstructionStatistics 建造统计数据
type ConstructionStatistics struct {
	BuildingID string `json:"building_id"`

	// 总体统计
	TotalConstructions      int64 `json:"total_constructions"`
	CompletedConstructions  int64 `json:"completed_constructions"`
	InProgressConstructions int64 `json:"in_progress_constructions"`
	CancelledConstructions  int64 `json:"cancelled_constructions"`
	FailedConstructions     int64 `json:"failed_constructions"`

	// 进度统计
	AverageProgress float64 `json:"average_progress"`
	MinProgress     float64 `json:"min_progress"`
	MaxProgress     float64 `json:"max_progress"`

	// 时间统计
	AverageDuration time.Duration `json:"average_duration"`
	MinDuration     time.Duration `json:"min_duration"`
	MaxDuration     time.Duration `json:"max_duration"`

	// 效率统计
	AverageEfficiency float64 `json:"average_efficiency"`
	MinEfficiency     float64 `json:"min_efficiency"`
	MaxEfficiency     float64 `json:"max_efficiency"`

	// 成本统计
	TotalCost   int64   `json:"total_cost"`
	AverageCost float64 `json:"average_cost"`
	MinCost     int64   `json:"min_cost"`
	MaxCost     int64   `json:"max_cost"`

	// 工人统计
	TotalWorkers   int64   `json:"total_workers"`
	AverageWorkers float64 `json:"average_workers"`
	MinWorkers     int64   `json:"min_workers"`
	MaxWorkers     int64   `json:"max_workers"`

	// 材料统计
	TotalMaterials   int64   `json:"total_materials"`
	AverageMaterials float64 `json:"average_materials"`
	MinMaterials     int64   `json:"min_materials"`
	MaxMaterials     int64   `json:"max_materials"`

	// 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// UpgradeStatistics 升级统计数据
type UpgradeStatistics struct {
	BuildingID string `json:"building_id"`

	// 总体统计
	TotalUpgrades      int64 `json:"total_upgrades"`
	CompletedUpgrades  int64 `json:"completed_upgrades"`
	InProgressUpgrades int64 `json:"in_progress_upgrades"`
	CancelledUpgrades  int64 `json:"cancelled_upgrades"`
	FailedUpgrades     int64 `json:"failed_upgrades"`

	// 等级统计
	AverageFromLevel float64 `json:"average_from_level"`
	AverageToLevel   float64 `json:"average_to_level"`
	MaxLevelReached  int32   `json:"max_level_reached"`

	// 进度统计
	AverageProgress float64 `json:"average_progress"`
	MinProgress     float64 `json:"min_progress"`
	MaxProgress     float64 `json:"max_progress"`

	// 时间统计
	AverageDuration time.Duration `json:"average_duration"`
	MinDuration     time.Duration `json:"min_duration"`
	MaxDuration     time.Duration `json:"max_duration"`

	// 成本统计
	TotalCost   int64   `json:"total_cost"`
	AverageCost float64 `json:"average_cost"`
	MinCost     int64   `json:"min_cost"`
	MaxCost     int64   `json:"max_cost"`

	// 收益统计
	TotalBenefits   int64   `json:"total_benefits"`
	AverageBenefits float64 `json:"average_benefits"`

	// 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// BlueprintStatistics 蓝图统计数据
type BlueprintStatistics struct {
	// 总体统计
	TotalBlueprints int64 `json:"total_blueprints"`

	// 按分类统计
	ByCategory map[BuildingCategory]int64 `json:"by_category"`

	// 按作者统计
	ByAuthor map[string]int64 `json:"by_author"`

	// 难度统计
	AverageDifficulty float64 `json:"average_difficulty"`
	MinDifficulty     int32   `json:"min_difficulty"`
	MaxDifficulty     int32   `json:"max_difficulty"`

	// 尺寸统计
	AverageWidth  float64 `json:"average_width"`
	AverageHeight float64 `json:"average_height"`
	AverageDepth  float64 `json:"average_depth"`
	MinWidth      int32   `json:"min_width"`
	MaxWidth      int32   `json:"max_width"`
	MinHeight     int32   `json:"min_height"`
	MaxHeight     int32   `json:"max_height"`
	MinDepth      int32   `json:"min_depth"`
	MaxDepth      int32   `json:"max_depth"`

	// 时间统计
	AverageDuration time.Duration `json:"average_duration"`
	MinDuration     time.Duration `json:"min_duration"`
	MaxDuration     time.Duration `json:"max_duration"`

	// 成本统计
	AverageCost float64 `json:"average_cost"`
	MinCost     int64   `json:"min_cost"`
	MaxCost     int64   `json:"max_cost"`

	// 材料统计
	AverageMaterials float64 `json:"average_materials"`
	MinMaterials     int64   `json:"min_materials"`
	MaxMaterials     int64   `json:"max_materials"`

	// 标签统计
	PopularTags []TagStatistic `json:"popular_tags"`

	// 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// TagStatistic 标签统计
type TagStatistic struct {
	Tag   string `json:"tag"`
	Count int64  `json:"count"`
}

// Area 区域
type Area struct {
	MinX int32 `json:"min_x"`
	MaxX int32 `json:"max_x"`
	MinY int32 `json:"min_y"`
	MaxY int32 `json:"max_y"`
	MinZ int32 `json:"min_z"`
	MaxZ int32 `json:"max_z"`
}

// NewArea 创建新区域
func NewArea(minX, maxX, minY, maxY, minZ, maxZ int32) *Area {
	return &Area{
		MinX: minX,
		MaxX: maxX,
		MinY: minY,
		MaxY: maxY,
		MinZ: minZ,
		MaxZ: maxZ,
	}
}

// IsValid 检查区域是否有效
func (a *Area) IsValid() bool {
	return a.MinX <= a.MaxX && a.MinY <= a.MaxY && a.MinZ <= a.MaxZ
}

// Contains 检查是否包含位置
func (a *Area) Contains(pos *Position) bool {
	if pos == nil {
		return false
	}
	return pos.X >= a.MinX && pos.X <= a.MaxX &&
		pos.Y >= a.MinY && pos.Y <= a.MaxY &&
		pos.Z >= a.MinZ && pos.Z <= a.MaxZ
}

// Overlaps 检查是否与另一个区域重叠
func (a *Area) Overlaps(other *Area) bool {
	if other == nil {
		return false
	}
	return !(a.MaxX < other.MinX || a.MinX > other.MaxX ||
		a.MaxY < other.MinY || a.MinY > other.MaxY ||
		a.MaxZ < other.MinZ || a.MinZ > other.MaxZ)
}

// GetVolume 获取区域体积
func (a *Area) GetVolume() int64 {
	width := int64(a.MaxX - a.MinX + 1)
	height := int64(a.MaxY - a.MinY + 1)
	depth := int64(a.MaxZ - a.MinZ + 1)
	return width * height * depth
}

// GetCenter 获取区域中心点
func (a *Area) GetCenter() *Position {
	return &Position{
		X: (a.MinX + a.MaxX) / 2,
		Y: (a.MinY + a.MaxY) / 2,
		Z: (a.MinZ + a.MaxZ) / 2,
	}
}

// 查询辅助函数

// NewBuildingQuery 创建新建筑查询
func NewBuildingQuery() *BuildingQuery {
	return &BuildingQuery{
		Page:      1,
		PageSize:  20,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}

// WithOwner 设置所有者
func (q *BuildingQuery) WithOwner(ownerID uint64) *BuildingQuery {
	q.OwnerID = &ownerID
	return q
}

// WithType 设置建筑类型
func (q *BuildingQuery) WithType(buildingType BuildingType) *BuildingQuery {
	q.Type = &buildingType
	return q
}

// WithCategory 设置建筑分类
func (q *BuildingQuery) WithCategory(category BuildingCategory) *BuildingQuery {
	q.Category = &category
	return q
}

// WithStatus 设置建筑状态
func (q *BuildingQuery) WithStatus(status BuildingStatus) *BuildingQuery {
	q.Status = &status
	return q
}

// WithLevelRange 设置等级范围
func (q *BuildingQuery) WithLevelRange(minLevel, maxLevel int32) *BuildingQuery {
	q.MinLevel = &minLevel
	q.MaxLevel = &maxLevel
	return q
}

// WithHealthRange 设置健康度范围
func (q *BuildingQuery) WithHealthRange(minHealth, maxHealth float64) *BuildingQuery {
	q.MinHealth = &minHealth
	q.MaxHealth = &maxHealth
	return q
}

// WithPosition 设置位置
func (q *BuildingQuery) WithPosition(position *Position) *BuildingQuery {
	q.Position = position
	return q
}

// WithArea 设置区域
func (q *BuildingQuery) WithArea(area *Area) *BuildingQuery {
	q.Area = area
	return q
}

// WithTags 设置标签
func (q *BuildingQuery) WithTags(tags []string) *BuildingQuery {
	q.Tags = tags
	return q
}

// WithSort 设置排序
func (q *BuildingQuery) WithSort(sortBy, sortOrder string) *BuildingQuery {
	q.SortBy = sortBy
	q.SortOrder = sortOrder
	return q
}

// WithPagination 设置分页
func (q *BuildingQuery) WithPagination(page, pageSize int) *BuildingQuery {
	q.Page = page
	q.PageSize = pageSize
	return q
}

// NewConstructionQuery 创建新建造查询
func NewConstructionQuery() *ConstructionQuery {
	return &ConstructionQuery{
		Page:      1,
		PageSize:  20,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}

// WithBuildingID 设置建筑ID
func (q *ConstructionQuery) WithBuildingID(buildingID string) *ConstructionQuery {
	q.BuildingID = &buildingID
	return q
}

// WithStatus 设置状态
func (q *ConstructionQuery) WithStatus(status ConstructionStatus) *ConstructionQuery {
	q.Status = &status
	return q
}

// WithWorker 设置工人
func (q *ConstructionQuery) WithWorker(workerID uint64) *ConstructionQuery {
	q.WorkerID = &workerID
	return q
}

// WithProgressRange 设置进度范围
func (q *ConstructionQuery) WithProgressRange(minProgress, maxProgress float64) *ConstructionQuery {
	q.MinProgress = &minProgress
	q.MaxProgress = &maxProgress
	return q
}

// WithSort 设置排序
func (q *ConstructionQuery) WithSort(sortBy, sortOrder string) *ConstructionQuery {
	q.SortBy = sortBy
	q.SortOrder = sortOrder
	return q
}

// WithPagination 设置分页
func (q *ConstructionQuery) WithPagination(page, pageSize int) *ConstructionQuery {
	q.Page = page
	q.PageSize = pageSize
	return q
}

// NewUpgradeQuery 创建新升级查询
func NewUpgradeQuery() *UpgradeQuery {
	return &UpgradeQuery{
		Page:      1,
		PageSize:  20,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}

// WithBuildingID 设置建筑ID
func (q *UpgradeQuery) WithBuildingID(buildingID string) *UpgradeQuery {
	q.BuildingID = &buildingID
	return q
}

// WithStatus 设置状态
func (q *UpgradeQuery) WithStatus(status UpgradeStatus) *UpgradeQuery {
	q.Status = &status
	return q
}

// WithLevelRange 设置等级范围
func (q *UpgradeQuery) WithLevelRange(fromLevel, toLevel int32) *UpgradeQuery {
	q.FromLevel = &fromLevel
	q.ToLevel = &toLevel
	return q
}

// WithProgressRange 设置进度范围
func (q *UpgradeQuery) WithProgressRange(minProgress, maxProgress float64) *UpgradeQuery {
	q.MinProgress = &minProgress
	q.MaxProgress = &maxProgress
	return q
}

// WithSort 设置排序
func (q *UpgradeQuery) WithSort(sortBy, sortOrder string) *UpgradeQuery {
	q.SortBy = sortBy
	q.SortOrder = sortOrder
	return q
}

// WithPagination 设置分页
func (q *UpgradeQuery) WithPagination(page, pageSize int) *UpgradeQuery {
	q.Page = page
	q.PageSize = pageSize
	return q
}

// NewBlueprintQuery 创建新蓝图查询
func NewBlueprintQuery() *BlueprintQuery {
	return &BlueprintQuery{
		Page:      1,
		PageSize:  20,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}

// WithName 设置名称
func (q *BlueprintQuery) WithName(name string) *BlueprintQuery {
	q.Name = &name
	return q
}

// WithAuthor 设置作者
func (q *BlueprintQuery) WithAuthor(author string) *BlueprintQuery {
	q.Author = &author
	return q
}

// WithCategory 设置分类
func (q *BlueprintQuery) WithCategory(category BuildingCategory) *BlueprintQuery {
	q.Category = &category
	return q
}

// WithDifficultyRange 设置难度范围
func (q *BlueprintQuery) WithDifficultyRange(minDifficulty, maxDifficulty int32) *BlueprintQuery {
	q.MinDifficulty = &minDifficulty
	q.MaxDifficulty = &maxDifficulty
	return q
}

// WithTags 设置标签
func (q *BlueprintQuery) WithTags(tags []string) *BlueprintQuery {
	q.Tags = tags
	return q
}

// WithKeyword 设置关键词
func (q *BlueprintQuery) WithKeyword(keyword string) *BlueprintQuery {
	q.Keyword = &keyword
	return q
}

// WithSort 设置排序
func (q *BlueprintQuery) WithSort(sortBy, sortOrder string) *BlueprintQuery {
	q.SortBy = sortBy
	q.SortOrder = sortOrder
	return q
}

// WithPagination 设置分页
func (q *BlueprintQuery) WithPagination(page, pageSize int) *BlueprintQuery {
	q.Page = page
	q.PageSize = pageSize
	return q
}

// 常量定义

const (
	// 默认分页大小
	DefaultPageSize = 20
	MaxPageSize     = 100

	// 默认排序
	DefaultSortBy    = "created_at"
	DefaultSortOrder = "desc"

	// 查询限制
	MaxTagsCount  = 10
	MaxKeywordLen = 100
	MaxNameLen    = 100
	MaxAuthorLen  = 50
	MaxVersionLen = 20
)

// 验证函数

// ValidateQuery 验证查询参数
func ValidateQuery(page, pageSize int, sortBy, sortOrder string, validSortFields []string) error {
	if page < 1 {
		return fmt.Errorf("page must be at least 1")
	}
	if pageSize < 1 || pageSize > MaxPageSize {
		return fmt.Errorf("page size must be between 1 and %d", MaxPageSize)
	}

	if sortBy != "" {
		valid := false
		for _, field := range validSortFields {
			if sortBy == field {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid sort field: %s", sortBy)
		}
	}

	if sortOrder != "" && sortOrder != "asc" && sortOrder != "desc" {
		return fmt.Errorf("sort order must be 'asc' or 'desc'")
	}

	return nil
}

// ValidateBuildingQuery 验证建筑查询
func ValidateBuildingQuery(query *BuildingQuery) error {
	if query == nil {
		return fmt.Errorf("query cannot be nil")
	}

	validSortFields := []string{"name", "type", "category", "level", "health", "created_at", "updated_at"}
	return ValidateQuery(query.Page, query.PageSize, query.SortBy, query.SortOrder, validSortFields)
}

// ValidateConstructionQuery 验证建造查询
func ValidateConstructionQuery(query *ConstructionQuery) error {
	if query == nil {
		return fmt.Errorf("query cannot be nil")
	}

	validSortFields := []string{"progress", "started_at", "duration", "created_at"}
	return ValidateQuery(query.Page, query.PageSize, query.SortBy, query.SortOrder, validSortFields)
}

// ValidateUpgradeQuery 验证升级查询
func ValidateUpgradeQuery(query *UpgradeQuery) error {
	if query == nil {
		return fmt.Errorf("query cannot be nil")
	}

	validSortFields := []string{"from_level", "to_level", "progress", "started_at", "created_at"}
	return ValidateQuery(query.Page, query.PageSize, query.SortBy, query.SortOrder, validSortFields)
}

// ValidateBlueprintQuery 验证蓝图查询
func ValidateBlueprintQuery(query *BlueprintQuery) error {
	if query == nil {
		return fmt.Errorf("query cannot be nil")
	}

	validSortFields := []string{"name", "author", "difficulty", "duration", "created_at"}
	return ValidateQuery(query.Page, query.PageSize, query.SortBy, query.SortOrder, validSortFields)
}
