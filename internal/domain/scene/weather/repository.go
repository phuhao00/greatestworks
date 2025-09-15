package weather

import (
	"context"
	"time"
)

// WeatherRepository 天气仓储接口
type WeatherRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, weather *WeatherAggregate) error
	FindByID(ctx context.Context, id string) (*WeatherAggregate, error)
	FindBySceneID(ctx context.Context, sceneID string) (*WeatherAggregate, error)
	Update(ctx context.Context, weather *WeatherAggregate) error
	Delete(ctx context.Context, id string) error
	
	// 查询操作
	FindBySceneIDs(ctx context.Context, sceneIDs []string) ([]*WeatherAggregate, error)
	FindActiveWeather(ctx context.Context, sceneID string) (*WeatherAggregate, error)
	FindWeatherByTimeRange(ctx context.Context, sceneID string, startTime, endTime time.Time) ([]*WeatherAggregate, error)
	FindWeatherByType(ctx context.Context, weatherType WeatherType, limit int) ([]*WeatherAggregate, error)
	
	// 统计操作
	GetWeatherStatistics(ctx context.Context, sceneID string, period time.Duration) (*WeatherStatistics, error)
	GetWeatherCount(ctx context.Context, sceneID string) (int64, error)
	GetWeatherCountByType(ctx context.Context, sceneID string, weatherType WeatherType) (int64, error)
	
	// 批量操作
	SaveBatch(ctx context.Context, weathers []*WeatherAggregate) error
	DeleteBatch(ctx context.Context, ids []string) error
	UpdateBatch(ctx context.Context, weathers []*WeatherAggregate) error
	
	// 清理操作
	CleanupExpiredWeather(ctx context.Context, beforeTime time.Time) (int64, error)
	CleanupOldHistory(ctx context.Context, sceneID string, keepDays int) (int64, error)
}

// WeatherStateRepository 天气状态仓储接口
type WeatherStateRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, state *WeatherState) error
	FindByID(ctx context.Context, id string) (*WeatherState, error)
	Update(ctx context.Context, state *WeatherState) error
	Delete(ctx context.Context, id string) error
	
	// 查询操作
	FindBySceneID(ctx context.Context, sceneID string, limit int) ([]*WeatherState, error)
	FindCurrentState(ctx context.Context, sceneID string) (*WeatherState, error)
	FindActiveStates(ctx context.Context, sceneIDs []string) ([]*WeatherState, error)
	FindByTimeRange(ctx context.Context, sceneID string, startTime, endTime time.Time) ([]*WeatherState, error)
	FindByWeatherType(ctx context.Context, weatherType WeatherType, limit int) ([]*WeatherState, error)
	FindByIntensity(ctx context.Context, intensity WeatherIntensity, limit int) ([]*WeatherState, error)
	
	// 历史记录操作
	SaveHistory(ctx context.Context, sceneID string, states []*WeatherState) error
	GetHistory(ctx context.Context, sceneID string, limit int) ([]*WeatherState, error)
	GetHistoryByTimeRange(ctx context.Context, sceneID string, startTime, endTime time.Time) ([]*WeatherState, error)
	
	// 统计操作
	GetStateStatistics(ctx context.Context, sceneID string, period time.Duration) (*StateStatistics, error)
	GetAverageTemperature(ctx context.Context, sceneID string, period time.Duration) (float64, error)
	GetAverageHumidity(ctx context.Context, sceneID string, period time.Duration) (float64, error)
	
	// 清理操作
	CleanupExpiredStates(ctx context.Context, beforeTime time.Time) (int64, error)
}

// WeatherEffectRepository 天气效果仓储接口
type WeatherEffectRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, effect *WeatherEffect) error
	FindByID(ctx context.Context, id string) (*WeatherEffect, error)
	Update(ctx context.Context, effect *WeatherEffect) error
	Delete(ctx context.Context, id string) error
	
	// 查询操作
	FindBySceneID(ctx context.Context, sceneID string) ([]*WeatherEffect, error)
	FindActiveEffects(ctx context.Context, sceneID string) ([]*WeatherEffect, error)
	FindByEffectType(ctx context.Context, effectType string, limit int) ([]*WeatherEffect, error)
	FindByTimeRange(ctx context.Context, sceneID string, startTime, endTime time.Time) ([]*WeatherEffect, error)
	
	// 批量操作
	SaveBatch(ctx context.Context, effects []*WeatherEffect) error
	DeleteBatch(ctx context.Context, ids []string) error
	UpdateBatch(ctx context.Context, effects []*WeatherEffect) error
	
	// 清理操作
	CleanupExpiredEffects(ctx context.Context, beforeTime time.Time) (int64, error)
	DeactivateEffects(ctx context.Context, sceneID string, effectTypes []string) error
}

// WeatherEventRepository 天气事件仓储接口
type WeatherEventRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, event *WeatherEvent) error
	FindByID(ctx context.Context, id string) (*WeatherEvent, error)
	Update(ctx context.Context, event *WeatherEvent) error
	Delete(ctx context.Context, id string) error
	
	// 查询操作
	FindBySceneID(ctx context.Context, sceneID string, limit int) ([]*WeatherEvent, error)
	FindActiveEvents(ctx context.Context, sceneID string) ([]*WeatherEvent, error)
	FindByEventType(ctx context.Context, eventType WeatherEventType, limit int) ([]*WeatherEvent, error)
	FindBySeverity(ctx context.Context, severity WeatherEventSeverity, limit int) ([]*WeatherEvent, error)
	FindByTimeRange(ctx context.Context, sceneID string, startTime, endTime time.Time) ([]*WeatherEvent, error)
	
	// 统计操作
	GetEventStatistics(ctx context.Context, sceneID string, period time.Duration) (*EventStatistics, error)
	GetEventCount(ctx context.Context, sceneID string) (int64, error)
	GetEventCountByType(ctx context.Context, eventType WeatherEventType) (int64, error)
	
	// 批量操作
	SaveBatch(ctx context.Context, events []*WeatherEvent) error
	DeleteBatch(ctx context.Context, ids []string) error
	
	// 清理操作
	CleanupExpiredEvents(ctx context.Context, beforeTime time.Time) (int64, error)
}

// WeatherForecastRepository 天气预报仓储接口
type WeatherForecastRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, forecast *WeatherForecast) error
	FindByID(ctx context.Context, id string) (*WeatherForecast, error)
	Update(ctx context.Context, forecast *WeatherForecast) error
	Delete(ctx context.Context, id string) error
	
	// 查询操作
	FindBySceneID(ctx context.Context, sceneID string, limit int) ([]*WeatherForecast, error)
	FindByTimeRange(ctx context.Context, sceneID string, startTime, endTime time.Time) ([]*WeatherForecast, error)
	FindLatestForecast(ctx context.Context, sceneID string, hours int) ([]*WeatherForecast, error)
	FindByWeatherType(ctx context.Context, weatherType WeatherType, limit int) ([]*WeatherForecast, error)
	
	// 批量操作
	SaveBatch(ctx context.Context, forecasts []*WeatherForecast) error
	DeleteBatch(ctx context.Context, ids []string) error
	UpdateBatch(ctx context.Context, forecasts []*WeatherForecast) error
	
	// 预报管理
	ReplaceForecast(ctx context.Context, sceneID string, forecasts []*WeatherForecast) error
	GetForecastAccuracy(ctx context.Context, sceneID string, period time.Duration) (float64, error)
	
	// 清理操作
	CleanupOldForecasts(ctx context.Context, beforeTime time.Time) (int64, error)
}

// SeasonalPatternRepository 季节模式仓储接口
type SeasonalPatternRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, pattern *SeasonalPattern) error
	FindByZoneID(ctx context.Context, zoneID string) (*SeasonalPattern, error)
	Update(ctx context.Context, pattern *SeasonalPattern) error
	Delete(ctx context.Context, zoneID string) error
	
	// 查询操作
	FindAll(ctx context.Context) ([]*SeasonalPattern, error)
	FindBySeason(ctx context.Context, season Season) ([]*SeasonalPattern, error)
	
	// 批量操作
	SaveBatch(ctx context.Context, patterns []*SeasonalPattern) error
	UpdateBatch(ctx context.Context, patterns []*SeasonalPattern) error
}

// 统计信息结构体

// WeatherStatistics 天气统计信息
type WeatherStatistics struct {
	SceneID              string
	Period               time.Duration
	StartTime            time.Time
	EndTime              time.Time
	TotalWeatherChanges  int64
	WeatherTypeCount     map[WeatherType]int64
	IntensityCount       map[WeatherIntensity]int64
	AverageTemperature   float64
	AverageHumidity      float64
	AverageWindSpeed     float64
	AverageVisibility    float64
	MostCommonWeather    WeatherType
	MostCommonIntensity  WeatherIntensity
	LongestWeatherPeriod time.Duration
	ShortestWeatherPeriod time.Duration
	CreatedAt            time.Time
}

// StateStatistics 状态统计信息
type StateStatistics struct {
	SceneID            string
	Period             time.Duration
	StartTime          time.Time
	EndTime            time.Time
	TotalStates        int64
	ActiveStates       int64
	ExpiredStates      int64
	AverageTemperature float64
	MinTemperature     float64
	MaxTemperature     float64
	AverageHumidity    float64
	MinHumidity        float64
	MaxHumidity        float64
	AverageWindSpeed   float64
	MaxWindSpeed       float64
	AverageVisibility  float64
	MinVisibility      float64
	CreatedAt          time.Time
}

// EventStatistics 事件统计信息
type EventStatistics struct {
	SceneID              string
	Period               time.Duration
	StartTime            time.Time
	EndTime              time.Time
	TotalEvents          int64
	ActiveEvents         int64
	ExpiredEvents        int64
	EventTypeCount       map[WeatherEventType]int64
	SeverityCount        map[WeatherEventSeverity]int64
	MostCommonEventType  WeatherEventType
	MostCommonSeverity   WeatherEventSeverity
	AverageEventDuration time.Duration
	LongestEventDuration time.Duration
	CreatedAt            time.Time
}

// 查询条件结构体

// WeatherQuery 天气查询条件
type WeatherQuery struct {
	SceneIDs     []string
	WeatherTypes []WeatherType
	Intensities  []WeatherIntensity
	StartTime    *time.Time
	EndTime      *time.Time
	IsActive     *bool
	Limit        int
	Offset       int
	OrderBy      string
	OrderDesc    bool
}

// EffectQuery 效果查询条件
type EffectQuery struct {
	SceneIDs    []string
	EffectTypes []string
	IsActive    *bool
	StartTime   *time.Time
	EndTime     *time.Time
	MinMultiplier *float64
	MaxMultiplier *float64
	Limit       int
	Offset      int
	OrderBy     string
	OrderDesc   bool
}

// EventQuery 事件查询条件
type EventQuery struct {
	SceneIDs   []string
	EventTypes []WeatherEventType
	Severities []WeatherEventSeverity
	StartTime  *time.Time
	EndTime    *time.Time
	IsActive   *bool
	Limit      int
	Offset     int
	OrderBy    string
	OrderDesc  bool
}

// ForecastQuery 预报查询条件
type ForecastQuery struct {
	SceneIDs     []string
	WeatherTypes []WeatherType
	Intensities  []WeatherIntensity
	StartTime    *time.Time
	EndTime      *time.Time
	MinConfidence *float64
	Limit        int
	Offset       int
	OrderBy      string
	OrderDesc    bool
}

// 趋势分析结构体

// WeatherTrend 天气趋势
type WeatherTrend struct {
	SceneID           string
	Period            time.Duration
	StartTime         time.Time
	EndTime           time.Time
	TrendType         TrendType
	WeatherTypeChanges map[WeatherType]float64 // 变化率
	TemperatureTrend  TemperatureTrend
	HumidityTrend     HumidityTrend
	VisibilityTrend   VisibilityTrend
	PredictedChanges  []PredictedChange
	Confidence        float64
	CreatedAt         time.Time
}

// TrendType 趋势类型
type TrendType int

const (
	TrendTypeStable TrendType = iota + 1
	TrendTypeIncreasing
	TrendTypeDecreasing
	TrendTypeVolatile
	TrendTypeCyclical
)

// String 返回趋势类型字符串
func (tt TrendType) String() string {
	switch tt {
	case TrendTypeStable:
		return "stable"
	case TrendTypeIncreasing:
		return "increasing"
	case TrendTypeDecreasing:
		return "decreasing"
	case TrendTypeVolatile:
		return "volatile"
	case TrendTypeCyclical:
		return "cyclical"
	default:
		return "unknown"
	}
}

// TemperatureTrend 温度趋势
type TemperatureTrend struct {
	Direction    TrendDirection
	ChangeRate   float64 // 每小时变化率
	AverageValue float64
	MinValue     float64
	MaxValue     float64
	Variance     float64
}

// HumidityTrend 湿度趋势
type HumidityTrend struct {
	Direction    TrendDirection
	ChangeRate   float64
	AverageValue float64
	MinValue     float64
	MaxValue     float64
	Variance     float64
}

// VisibilityTrend 能见度趋势
type VisibilityTrend struct {
	Direction    TrendDirection
	ChangeRate   float64
	AverageValue float64
	MinValue     float64
	MaxValue     float64
	Variance     float64
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

// PredictedChange 预测变化
type PredictedChange struct {
	Time        time.Time
	WeatherType WeatherType
	Intensity   WeatherIntensity
	Probability float64
	Confidence  float64
	Reason      string
}

// 缓存接口

// WeatherCacheRepository 天气缓存仓储接口
type WeatherCacheRepository interface {
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

// WeatherTransaction 天气事务接口
type WeatherTransaction interface {
	// 事务控制
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	
	// 获取仓储
	WeatherRepository() WeatherRepository
	WeatherStateRepository() WeatherStateRepository
	WeatherEffectRepository() WeatherEffectRepository
	WeatherEventRepository() WeatherEventRepository
	WeatherForecastRepository() WeatherForecastRepository
	SeasonalPatternRepository() SeasonalPatternRepository
}

// 仓储工厂接口

// WeatherRepositoryFactory 天气仓储工厂接口
type WeatherRepositoryFactory interface {
	// 创建仓储
	CreateWeatherRepository() WeatherRepository
	CreateWeatherStateRepository() WeatherStateRepository
	CreateWeatherEffectRepository() WeatherEffectRepository
	CreateWeatherEventRepository() WeatherEventRepository
	CreateWeatherForecastRepository() WeatherForecastRepository
	CreateSeasonalPatternRepository() SeasonalPatternRepository
	CreateWeatherCacheRepository() WeatherCacheRepository
	
	// 创建事务
	CreateTransaction() WeatherTransaction
	
	// 健康检查
	HealthCheck(ctx context.Context) error
	
	// 关闭连接
	Close() error
}