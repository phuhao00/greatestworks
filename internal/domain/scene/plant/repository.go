package plant

import (
	"context"
	"time"
)

// FarmRepository 农场仓储接口
type FarmRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, farm *FarmAggregate) error
	FindByID(ctx context.Context, farmID string) (*FarmAggregate, error)
	FindByOwner(ctx context.Context, owner string) ([]*FarmAggregate, error)
	FindBySceneID(ctx context.Context, sceneID string) ([]*FarmAggregate, error)
	Update(ctx context.Context, farm *FarmAggregate) error
	Delete(ctx context.Context, farmID string) error

	// 查询操作
	FindBySize(ctx context.Context, size FarmSize, limit int) ([]*FarmAggregate, error)
	FindByStatus(ctx context.Context, status FarmStatus, limit int) ([]*FarmAggregate, error)
	FindActiveByOwner(ctx context.Context, owner string) ([]*FarmAggregate, error)
	FindByClimateZone(ctx context.Context, climateZone string, limit int) ([]*FarmAggregate, error)

	// 统计操作
	GetFarmStatistics(ctx context.Context, farmID string) (*FarmStatistics, error)
	GetOwnerStatistics(ctx context.Context, owner string) (*OwnerStatistics, error)
	GetFarmCount(ctx context.Context) (int64, error)
	GetFarmCountByOwner(ctx context.Context, owner string) (int64, error)

	// 批量操作
	SaveBatch(ctx context.Context, farms []*FarmAggregate) error
	UpdateBatch(ctx context.Context, farms []*FarmAggregate) error
	DeleteBatch(ctx context.Context, farmIDs []string) error

	// 排行榜操作
	GetTopFarmsByValue(ctx context.Context, limit int) ([]*FarmRanking, error)
	GetTopFarmsByProductivity(ctx context.Context, limit int) ([]*FarmRanking, error)
	GetTopFarmsByYield(ctx context.Context, period time.Duration, limit int) ([]*FarmRanking, error)
}

// CropRepository 作物仓储接口
type CropRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, crop *Crop) error
	FindByID(ctx context.Context, cropID string) (*Crop, error)
	Update(ctx context.Context, crop *Crop) error
	Delete(ctx context.Context, cropID string) error

	// 查询操作
	FindByFarmID(ctx context.Context, farmID string) ([]*Crop, error)
	FindBySeedType(ctx context.Context, seedType SeedType, limit int) ([]*Crop, error)
	FindByGrowthStage(ctx context.Context, stage GrowthStage, limit int) ([]*Crop, error)
	FindHarvestable(ctx context.Context, farmID string) ([]*Crop, error)
	FindNeedsCare(ctx context.Context, farmID string) ([]*Crop, error)
	FindByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*Crop, error)

	// 状态查询
	FindByHealthRange(ctx context.Context, minHealth, maxHealth float64, limit int) ([]*Crop, error)
	FindByProgressRange(ctx context.Context, minProgress, maxProgress float64, limit int) ([]*Crop, error)
	FindExpiredCrops(ctx context.Context, beforeTime time.Time) ([]*Crop, error)

	// 统计操作
	GetCropStatistics(ctx context.Context, farmID string) (*CropStatistics, error)
	GetCropCountByType(ctx context.Context, seedType SeedType) (int64, error)
	GetCropCountByStage(ctx context.Context, stage GrowthStage) (int64, error)
	GetAverageGrowthProgress(ctx context.Context, seedType SeedType) (float64, error)

	// 批量操作
	SaveBatch(ctx context.Context, crops []*Crop) error
	UpdateBatch(ctx context.Context, crops []*Crop) error
	DeleteBatch(ctx context.Context, cropIDs []string) error

	// 清理操作
	CleanupExpiredCrops(ctx context.Context, beforeTime time.Time) (int64, error)
}

// PlotRepository 地块仓储接口
type PlotRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, plot *Plot) error
	FindByID(ctx context.Context, plotID string) (*Plot, error)
	Update(ctx context.Context, plot *Plot) error
	Delete(ctx context.Context, plotID string) error

	// 查询操作
	FindByFarmID(ctx context.Context, farmID string) ([]*Plot, error)
	FindAvailable(ctx context.Context, farmID string) ([]*Plot, error)
	FindOccupied(ctx context.Context, farmID string) ([]*Plot, error)
	FindBySize(ctx context.Context, size PlotSize, limit int) ([]*Plot, error)
	FindBySoilType(ctx context.Context, soilType SoilType, limit int) ([]*Plot, error)

	// 统计操作
	GetPlotStatistics(ctx context.Context, farmID string) (*PlotStatistics, error)
	GetAvailablePlotCount(ctx context.Context, farmID string) (int64, error)
	GetOccupiedPlotCount(ctx context.Context, farmID string) (int64, error)

	// 批量操作
	SaveBatch(ctx context.Context, plots []*Plot) error
	UpdateBatch(ctx context.Context, plots []*Plot) error
	DeleteBatch(ctx context.Context, plotIDs []string) error
}

// FarmToolRepository 农具仓储接口
type FarmToolRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, tool *FarmTool) error
	FindByID(ctx context.Context, toolID string) (*FarmTool, error)
	Update(ctx context.Context, tool *FarmTool) error
	Delete(ctx context.Context, toolID string) error

	// 查询操作
	FindByFarmID(ctx context.Context, farmID string) ([]*FarmTool, error)
	FindByType(ctx context.Context, toolType ToolType, limit int) ([]*FarmTool, error)
	FindByLevel(ctx context.Context, minLevel, maxLevel int, limit int) ([]*FarmTool, error)
	FindUsable(ctx context.Context, farmID string) ([]*FarmTool, error)
	FindNeedsMaintenance(ctx context.Context, farmID string, durabilityThreshold float64) ([]*FarmTool, error)

	// 统计操作
	GetToolStatistics(ctx context.Context, farmID string) (*ToolStatistics, error)
	GetToolCountByType(ctx context.Context, toolType ToolType) (int64, error)
	GetAverageToolLevel(ctx context.Context, toolType ToolType) (float64, error)

	// 批量操作
	SaveBatch(ctx context.Context, tools []*FarmTool) error
	UpdateBatch(ctx context.Context, tools []*FarmTool) error
	DeleteBatch(ctx context.Context, toolIDs []string) error
}

// SoilRepository 土壤仓储接口
type SoilRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, soil *Soil) error
	FindByID(ctx context.Context, soilID string) (*Soil, error)
	Update(ctx context.Context, soil *Soil) error
	Delete(ctx context.Context, soilID string) error

	// 查询操作
	FindByFarmID(ctx context.Context, farmID string) (*Soil, error)
	FindByType(ctx context.Context, soilType SoilType, limit int) ([]*Soil, error)
	FindByFertilityRange(ctx context.Context, minFertility, maxFertility float64, limit int) ([]*Soil, error)
	FindByPHRange(ctx context.Context, minPH, maxPH float64, limit int) ([]*Soil, error)
	FindHighQuality(ctx context.Context, qualityThreshold float64, limit int) ([]*Soil, error)

	// 统计操作
	GetSoilStatistics(ctx context.Context, farmID string) (*SoilStatistics, error)
	GetAverageFertility(ctx context.Context, soilType SoilType) (float64, error)
	GetAveragePH(ctx context.Context, soilType SoilType) (float64, error)

	// 历史记录
	SaveSoilHistory(ctx context.Context, farmID string, soil *Soil) error
	GetSoilHistory(ctx context.Context, farmID string, limit int) ([]*SoilHistoryRecord, error)

	// 批量操作
	SaveBatch(ctx context.Context, soils []*Soil) error
	UpdateBatch(ctx context.Context, soils []*Soil) error
}

// HarvestRepository 收获仓储接口
type HarvestRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, harvest *HarvestResult) error
	FindByID(ctx context.Context, harvestID string) (*HarvestResult, error)
	Update(ctx context.Context, harvest *HarvestResult) error
	Delete(ctx context.Context, harvestID string) error

	// 查询操作
	FindByFarmID(ctx context.Context, farmID string, limit int) ([]*HarvestResult, error)
	FindByCropID(ctx context.Context, cropID string) (*HarvestResult, error)
	FindBySeedType(ctx context.Context, seedType SeedType, limit int) ([]*HarvestResult, error)
	FindByQuality(ctx context.Context, quality CropQuality, limit int) ([]*HarvestResult, error)
	FindByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*HarvestResult, error)

	// 统计操作
	GetHarvestStatistics(ctx context.Context, farmID string, period time.Duration) (*HarvestStatistics, error)
	GetTotalYield(ctx context.Context, farmID string, seedType SeedType, period time.Duration) (int, error)
	GetAverageQuality(ctx context.Context, farmID string, seedType SeedType, period time.Duration) (float64, error)
	GetHarvestTrend(ctx context.Context, farmID string, period time.Duration) (*HarvestTrend, error)

	// 排行榜操作
	GetTopHarvestsByYield(ctx context.Context, period time.Duration, limit int) ([]*HarvestRanking, error)
	GetTopHarvestsByQuality(ctx context.Context, period time.Duration, limit int) ([]*HarvestRanking, error)

	// 批量操作
	SaveBatch(ctx context.Context, harvests []*HarvestResult) error
	DeleteBatch(ctx context.Context, harvestIDs []string) error

	// 清理操作
	CleanupOldHarvests(ctx context.Context, beforeTime time.Time) (int64, error)
}

// PlantEventRepository 种植事件仓储接口
type PlantEventRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, event *PlantEvent) error
	FindByID(ctx context.Context, eventID string) (*PlantEvent, error)
	Update(ctx context.Context, event *PlantEvent) error
	Delete(ctx context.Context, eventID string) error

	// 查询操作
	FindByFarmID(ctx context.Context, farmID string, limit int) ([]*PlantEvent, error)
	FindByCropID(ctx context.Context, cropID string, limit int) ([]*PlantEvent, error)
	FindByEventType(ctx context.Context, eventType string, limit int) ([]*PlantEvent, error)
	FindByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*PlantEvent, error)
	FindActiveEvents(ctx context.Context, farmID string) ([]*PlantEvent, error)

	// 统计操作
	GetEventStatistics(ctx context.Context, farmID string, period time.Duration) (*EventStatistics, error)
	GetEventCountByType(ctx context.Context, eventType string, period time.Duration) (int64, error)

	// 批量操作
	SaveBatch(ctx context.Context, events []*PlantEvent) error
	DeleteBatch(ctx context.Context, eventIDs []string) error

	// 清理操作
	CleanupExpiredEvents(ctx context.Context, beforeTime time.Time) (int64, error)
}

// OwnerStatistics 所有者统计信息
type OwnerStatistics struct {
	Owner               string
	TotalFarms          int
	TotalPlots          int
	TotalCrops          int
	TotalHarvests       int
	TotalYield          int
	TotalExperience     int
	TotalValue          float64
	AverageProductivity float64
	BestFarmID          string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// CropStatistics 作物统计信息
type CropStatistics struct {
	FarmID                string
	TotalCrops            int
	CropsByType           map[SeedType]int
	CropsByStage          map[GrowthStage]int
	AverageGrowthProgress float64
	AverageHealthScore    float64
	HarvestableCrops      int
	CropsNeedingCare      int
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// PlotStatistics 地块统计信息
type PlotStatistics struct {
	FarmID          string
	TotalPlots      int
	AvailablePlots  int
	OccupiedPlots   int
	PlotsBySize     map[PlotSize]int
	PlotsBySoil     map[SoilType]int
	UtilizationRate float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ToolStatistics 工具统计信息
type ToolStatistics struct {
	FarmID             string
	TotalTools         int
	ToolsByType        map[ToolType]int
	ToolsByLevel       map[int]int
	UsableTools        int
	ToolsNeedingRepair int
	AverageEfficiency  float64
	AverageDurability  float64
	TotalValue         float64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// SoilStatistics 土壤统计信息
type SoilStatistics struct {
	FarmID            string
	SoilType          SoilType
	Fertility         float64
	PH                float64
	Moisture          float64
	Organic           float64
	Nitrogen          float64
	Phosphorus        float64
	Potassium         float64
	QualityScore      float64
	ProductivityScore float64
	LastTested        time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// HarvestStatistics 收获统计信息
type HarvestStatistics struct {
	FarmID            string
	Period            time.Duration
	StartTime         time.Time
	EndTime           time.Time
	TotalHarvests     int
	TotalYield        int
	HarvestsByType    map[SeedType]int
	YieldByType       map[SeedType]int
	HarvestsByQuality map[CropQuality]int
	AverageYield      float64
	AverageQuality    float64
	BestHarvest       *HarvestResult
	TotalExperience   int
	TotalValue        float64
	CreatedAt         time.Time
}

// EventStatistics 事件统计信息
type EventStatistics struct {
	FarmID         string
	Period         time.Duration
	StartTime      time.Time
	EndTime        time.Time
	TotalEvents    int
	EventsByType   map[string]int
	ActiveEvents   int
	ResolvedEvents int
	CriticalEvents int
	CreatedAt      time.Time
}

// 排行榜结构体

// FarmRanking 农场排行
type FarmRanking struct {
	Rank        int
	FarmID      string
	Owner       string
	FarmName    string
	Score       float64
	Metric      string
	Value       interface{}
	LastUpdated time.Time
}

// HarvestRanking 收获排行
type HarvestRanking struct {
	Rank        int
	FarmID      string
	Owner       string
	HarvestID   string
	SeedType    SeedType
	Yield       int
	Quality     CropQuality
	Score       float64
	HarvestTime time.Time
}

// 趋势分析结构体

// HarvestTrend 收获趋势
type HarvestTrend struct {
	FarmID       string
	Period       time.Duration
	StartTime    time.Time
	EndTime      time.Time
	TrendType    TrendType
	YieldTrend   YieldTrend
	QualityTrend QualityTrend
	DataPoints   []*TrendDataPoint
	Prediction   *TrendPrediction
	Confidence   float64
	CreatedAt    time.Time
}

// TrendType 趋势类型
type TrendType int

const (
	TrendTypeIncreasing TrendType = iota + 1
	TrendTypeDecreasing
	TrendTypeStable
	TrendTypeVolatile
	TrendTypeCyclical
)

// String 返回趋势类型字符串
func (tt TrendType) String() string {
	switch tt {
	case TrendTypeIncreasing:
		return "increasing"
	case TrendTypeDecreasing:
		return "decreasing"
	case TrendTypeStable:
		return "stable"
	case TrendTypeVolatile:
		return "volatile"
	case TrendTypeCyclical:
		return "cyclical"
	default:
		return "unknown"
	}
}

// YieldTrend 产量趋势
type YieldTrend struct {
	Direction    TrendDirection
	ChangeRate   float64
	AverageYield float64
	MinYield     int
	MaxYield     int
	Variance     float64
}

// QualityTrend 品质趋势
type QualityTrend struct {
	Direction      TrendDirection
	ChangeRate     float64
	AverageQuality float64
	BestQuality    CropQuality
	WorstQuality   CropQuality
	Variance       float64
}

// TrendDirection 趋势方向
type TrendDirection int

const (
	TrendDirectionUp TrendDirection = iota + 1
	TrendDirectionDown
	TrendDirectionStable
	TrendDirectionVolatile
)

// String 返回趋势方向字符串
func (td TrendDirection) String() string {
	switch td {
	case TrendDirectionUp:
		return "up"
	case TrendDirectionDown:
		return "down"
	case TrendDirectionStable:
		return "stable"
	case TrendDirectionVolatile:
		return "volatile"
	default:
		return "unknown"
	}
}

// TrendDataPoint 趋势数据点
type TrendDataPoint struct {
	Time     time.Time
	Yield    int
	Quality  float64
	Value    float64
	Metadata map[string]interface{}
}

// TrendPrediction 趋势预测
type TrendPrediction struct {
	PredictedYield   int
	PredictedQuality float64
	PredictedValue   float64
	Timeframe        time.Duration
	Confidence       float64
	Factors          []string
}

// 查询条件结构体

// FarmQuery 农场查询条件
type FarmQuery struct {
	Owners        []string
	SceneIDs      []string
	Sizes         []FarmSize
	Statuses      []FarmStatus
	ClimateZones  []string
	MinValue      *float64
	MaxValue      *float64
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	Limit         int
	Offset        int
	OrderBy       string
	OrderDesc     bool
}

// CropQuery 作物查询条件
type CropQuery struct {
	FarmIDs       []string
	SeedTypes     []SeedType
	GrowthStages  []GrowthStage
	MinHealth     *float64
	MaxHealth     *float64
	MinProgress   *float64
	MaxProgress   *float64
	IsHarvestable *bool
	NeedsCare     *bool
	PlantedAfter  *time.Time
	PlantedBefore *time.Time
	Limit         int
	Offset        int
	OrderBy       string
	OrderDesc     bool
}

// HarvestQuery 收获查询条件
type HarvestQuery struct {
	FarmIDs         []string
	SeedTypes       []SeedType
	Qualities       []CropQuality
	MinYield        *int
	MaxYield        *int
	HarvestedAfter  *time.Time
	HarvestedBefore *time.Time
	Limit           int
	Offset          int
	OrderBy         string
	OrderDesc       bool
}

// 历史记录结构体

// SoilHistoryRecord 土壤历史记录
type SoilHistoryRecord struct {
	ID         string
	FarmID     string
	Soil       *Soil
	ChangeType string
	Changes    map[string]interface{}
	Reason     string
	RecordedAt time.Time
}

// PlantEvent 种植事件
type PlantEvent struct {
	ID          string
	FarmID      string
	CropID      string
	EventType   string
	Title       string
	Description string
	Severity    string
	Status      string
	Data        map[string]interface{}
	OccurredAt  time.Time
	ResolvedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// 缓存接口

// PlantCacheRepository 种植缓存仓储接口
type PlantCacheRepository interface {
	// 缓存操作
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)

	// 批量操作
	SetBatch(ctx context.Context, items map[string]interface{}, expiration time.Duration) error
	GetBatch(ctx context.Context, keys []string) (map[string]interface{}, error)
	DeleteBatch(ctx context.Context, keys []string) error

	// 模式操作
	DeleteByPattern(ctx context.Context, pattern string) error
	GetKeysByPattern(ctx context.Context, pattern string) ([]string, error)

	// 缓存管理
	Flush(ctx context.Context) error
	GetStats(ctx context.Context) (*CacheStats, error)
}

// CacheStats 缓存统计
type CacheStats struct {
	Hits        int64
	Misses      int64
	Keys        int64
	MemoryUsage int64
	HitRate     float64
	CreatedAt   time.Time
}

// 事务接口

// PlantTransaction 种植事务接口
type PlantTransaction interface {
	// 事务控制
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	// 获取仓储
	FarmRepository() FarmRepository
	CropRepository() CropRepository
	PlotRepository() PlotRepository
	FarmToolRepository() FarmToolRepository
	SoilRepository() SoilRepository
	HarvestRepository() HarvestRepository
	PlantEventRepository() PlantEventRepository
}

// 仓储工厂接口

// PlantRepositoryFactory 种植仓储工厂接口
type PlantRepositoryFactory interface {
	// 创建仓储
	CreateFarmRepository() FarmRepository
	CreateCropRepository() CropRepository
	CreatePlotRepository() PlotRepository
	CreateFarmToolRepository() FarmToolRepository
	CreateSoilRepository() SoilRepository
	CreateHarvestRepository() HarvestRepository
	CreatePlantEventRepository() PlantEventRepository
	CreatePlantCacheRepository() PlantCacheRepository

	// 创建事务
	CreateTransaction() PlantTransaction

	// 健康检查
	HealthCheck(ctx context.Context) error

	// 关闭连接
	Close() error
}
