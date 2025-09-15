package services

import (
	"context"
	"fmt"
	"time"
	
	"greatestworks/internal/domain/scene/weather"
)

// WeatherService 天气应用服务
type WeatherService struct {
	weatherRepo    weather.WeatherRepository
	forecastRepo   weather.ForecastRepository
	effectRepo     weather.EffectRepository
	statisticsRepo weather.StatisticsRepository
	cacheRepo      weather.CacheRepository
	weatherService *weather.WeatherService
}

// NewWeatherService 创建天气应用服务
func NewWeatherService(
	weatherRepo weather.WeatherRepository,
	forecastRepo weather.ForecastRepository,
	effectRepo weather.EffectRepository,
	statisticsRepo weather.StatisticsRepository,
	cacheRepo weather.CacheRepository,
	weatherService *weather.WeatherService,
) *WeatherService {
	return &WeatherService{
		weatherRepo:    weatherRepo,
		forecastRepo:   forecastRepo,
		effectRepo:     effectRepo,
		statisticsRepo: statisticsRepo,
		cacheRepo:      cacheRepo,
		weatherService: weatherService,
	}
}

// GetCurrentWeather 获取当前天气
func (s *WeatherService) GetCurrentWeather(ctx context.Context, regionID string) (*WeatherDTO, error) {
	// 先从缓存获取
	cachedWeather, err := s.cacheRepo.GetCurrentWeather(regionID)
	if err == nil && cachedWeather != nil {
		return s.buildWeatherDTO(cachedWeather), nil
	}
	
	// 从数据库获取
	currentWeather, err := s.weatherRepo.FindCurrentByRegion(regionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current weather: %w", err)
	}
	
	// 更新缓存
	if err := s.cacheRepo.SetCurrentWeather(regionID, currentWeather, time.Minute*10); err != nil {
		// 缓存更新失败不影响主流程
		// TODO: 添加日志记录
	}
	
	return s.buildWeatherDTO(currentWeather), nil
}

// GetWeatherForecast 获取天气预报
func (s *WeatherService) GetWeatherForecast(ctx context.Context, regionID string, days int) ([]*WeatherForecastDTO, error) {
	// 先从缓存获取
	cachedForecast, err := s.cacheRepo.GetWeatherForecast(regionID, days)
	if err == nil && len(cachedForecast) > 0 {
		return s.buildForecastDTOs(cachedForecast), nil
	}
	
	// 生成天气预报
	forecasts, err := s.weatherService.GenerateForecast(regionID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to generate weather forecast: %w", err)
	}
	
	// 保存预报数据
	for _, forecast := range forecasts {
		if err := s.forecastRepo.Save(forecast); err != nil {
			// 保存失败不影响返回结果
			// TODO: 添加日志记录
		}
	}
	
	// 更新缓存
	if err := s.cacheRepo.SetWeatherForecast(regionID, forecasts, time.Hour); err != nil {
		// 缓存更新失败不影响主流程
		// TODO: 添加日志记录
	}
	
	return s.buildForecastDTOs(forecasts), nil
}

// GetWeatherEffects 获取天气影响
func (s *WeatherService) GetWeatherEffects(ctx context.Context, regionID string, targetType weather.EffectTargetType) ([]*WeatherEffectDTO, error) {
	// 获取当前天气
	currentWeather, err := s.weatherRepo.FindCurrentByRegion(regionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current weather: %w", err)
	}
	
	// 获取天气影响
	effects, err := s.weatherService.CalculateEffects(currentWeather, targetType)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate weather effects: %w", err)
	}
	
	return s.buildEffectDTOs(effects), nil
}

// UpdateWeather 更新天气（系统调用）
func (s *WeatherService) UpdateWeather(ctx context.Context, regionID string) error {
	// 获取当前天气
	currentWeather, err := s.weatherRepo.FindCurrentByRegion(regionID)
	if err != nil && !weather.IsNotFoundError(err) {
		return fmt.Errorf("failed to get current weather: %w", err)
	}
	
	// 生成新天气
	newWeather, err := s.weatherService.GenerateNextWeather(regionID, currentWeather)
	if err != nil {
		return fmt.Errorf("failed to generate new weather: %w", err)
	}
	
	// 保存新天气
	if err := s.weatherRepo.Save(newWeather); err != nil {
		return fmt.Errorf("failed to save new weather: %w", err)
	}
	
	// 更新统计数据
	if err := s.updateStatistics(ctx, regionID, newWeather); err != nil {
		// 统计更新失败不影响主流程
		// TODO: 添加日志记录
	}
	
	// 清除相关缓存
	if err := s.cacheRepo.DeleteCurrentWeather(regionID); err != nil {
		// 缓存清除失败不影响主流程
		// TODO: 添加日志记录
	}
	
	return nil
}

// GetWeatherHistory 获取天气历史
func (s *WeatherService) GetWeatherHistory(ctx context.Context, regionID string, startTime, endTime time.Time) ([]*WeatherHistoryDTO, error) {
	weatherHistory, err := s.weatherRepo.FindByRegionAndTimeRange(regionID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather history: %w", err)
	}
	
	return s.buildHistoryDTOs(weatherHistory), nil
}

// GetWeatherStatistics 获取天气统计
func (s *WeatherService) GetWeatherStatistics(ctx context.Context, regionID string) (*WeatherStatisticsDTO, error) {
	stats, err := s.statisticsRepo.FindByRegion(regionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather statistics: %w", err)
	}
	
	return s.buildStatisticsDTO(stats), nil
}

// GetGlobalWeatherInfo 获取全局天气信息
func (s *WeatherService) GetGlobalWeatherInfo(ctx context.Context) (*GlobalWeatherDTO, error) {
	// 先从缓存获取
	cachedGlobal, err := s.cacheRepo.GetGlobalWeatherInfo()
	if err == nil && cachedGlobal != nil {
		return cachedGlobal, nil
	}
	
	// 获取所有区域的当前天气
	allWeather, err := s.weatherRepo.FindAllCurrent()
	if err != nil {
		return nil, fmt.Errorf("failed to get all current weather: %w", err)
	}
	
	// 构建全局天气信息
	globalInfo := s.buildGlobalWeatherDTO(allWeather)
	
	// 更新缓存
	if err := s.cacheRepo.SetGlobalWeatherInfo(globalInfo, time.Minute*5); err != nil {
		// 缓存更新失败不影响主流程
		// TODO: 添加日志记录
	}
	
	return globalInfo, nil
}

// TriggerSpecialWeather 触发特殊天气事件
func (s *WeatherService) TriggerSpecialWeather(ctx context.Context, regionID string, weatherType weather.WeatherType, duration time.Duration) error {
	// 创建特殊天气事件
	specialWeather, err := s.weatherService.CreateSpecialWeather(regionID, weatherType, duration)
	if err != nil {
		return fmt.Errorf("failed to create special weather: %w", err)
	}
	
	// 保存特殊天气
	if err := s.weatherRepo.Save(specialWeather); err != nil {
		return fmt.Errorf("failed to save special weather: %w", err)
	}
	
	// 清除相关缓存
	if err := s.cacheRepo.DeleteCurrentWeather(regionID); err != nil {
		// 缓存清除失败不影响主流程
		// TODO: 添加日志记录
	}
	
	return nil
}

// GetSeasonInfo 获取季节信息
func (s *WeatherService) GetSeasonInfo(ctx context.Context, regionID string) (*SeasonInfoDTO, error) {
	// 先从缓存获取
	cachedSeason, err := s.cacheRepo.GetSeasonInfo(regionID)
	if err == nil && cachedSeason != nil {
		return cachedSeason, nil
	}
	
	// 计算当前季节信息
	seasonInfo, err := s.weatherService.GetCurrentSeason(regionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get season info: %w", err)
	}
	
	seasonDTO := s.buildSeasonInfoDTO(seasonInfo)
	
	// 更新缓存
	if err := s.cacheRepo.SetSeasonInfo(regionID, seasonDTO, time.Hour*6); err != nil {
		// 缓存更新失败不影响主流程
		// TODO: 添加日志记录
	}
	
	return seasonDTO, nil
}

// 私有方法

// updateStatistics 更新统计数据
func (s *WeatherService) updateStatistics(ctx context.Context, regionID string, newWeather *weather.WeatherAggregate) error {
	stats, err := s.statisticsRepo.FindByRegion(regionID)
	if err != nil && !weather.IsNotFoundError(err) {
		return err
	}
	
	if stats == nil {
		stats = weather.NewWeatherStatistics(regionID)
	}
	
	// 更新统计数据
	stats.AddWeatherRecord(newWeather.GetWeatherType(), newWeather.GetIntensity())
	stats.UpdateLastWeatherTime(newWeather.GetStartTime())
	
	// 保存统计数据
	return s.statisticsRepo.Save(stats)
}

// buildWeatherDTO 构建天气DTO
func (s *WeatherService) buildWeatherDTO(weatherAggregate *weather.WeatherAggregate) *WeatherDTO {
	return &WeatherDTO{
		RegionID:    weatherAggregate.GetRegionID(),
		WeatherType: string(weatherAggregate.GetWeatherType()),
		Intensity:   weatherAggregate.GetIntensity(),
		Temperature: weatherAggregate.GetTemperature(),
		Humidity:    weatherAggregate.GetHumidity(),
		WindSpeed:   weatherAggregate.GetWindSpeed(),
		Visibility:  weatherAggregate.GetVisibility(),
		StartTime:   weatherAggregate.GetStartTime(),
		EndTime:     weatherAggregate.GetEndTime(),
		Duration:    weatherAggregate.GetDuration(),
		IsSpecial:   weatherAggregate.IsSpecialWeather(),
		Description: weatherAggregate.GetDescription(),
	}
}

// buildForecastDTOs 构建预报DTO列表
func (s *WeatherService) buildForecastDTOs(forecasts []*weather.WeatherForecast) []*WeatherForecastDTO {
	dtos := make([]*WeatherForecastDTO, len(forecasts))
	for i, forecast := range forecasts {
		dtos[i] = &WeatherForecastDTO{
			RegionID:     forecast.GetRegionID(),
			ForecastDate: forecast.GetForecastDate(),
			WeatherType:  string(forecast.GetWeatherType()),
			Intensity:    forecast.GetIntensity(),
			Temperature:  forecast.GetTemperature(),
			Humidity:     forecast.GetHumidity(),
			WindSpeed:    forecast.GetWindSpeed(),
			Probability:  forecast.GetProbability(),
			Description:  forecast.GetDescription(),
		}
	}
	return dtos
}

// buildEffectDTOs 构建影响DTO列表
func (s *WeatherService) buildEffectDTOs(effects []*weather.WeatherEffect) []*WeatherEffectDTO {
	dtos := make([]*WeatherEffectDTO, len(effects))
	for i, effect := range effects {
		dtos[i] = &WeatherEffectDTO{
			EffectType:   string(effect.GetEffectType()),
			TargetType:   string(effect.GetTargetType()),
			Modifier:     effect.GetModifier(),
			Duration:     effect.GetDuration(),
			IsPositive:   effect.IsPositive(),
			Description:  effect.GetDescription(),
		}
	}
	return dtos
}

// buildHistoryDTOs 构建历史DTO列表
func (s *WeatherService) buildHistoryDTOs(history []*weather.WeatherAggregate) []*WeatherHistoryDTO {
	dtos := make([]*WeatherHistoryDTO, len(history))
	for i, record := range history {
		dtos[i] = &WeatherHistoryDTO{
			RegionID:    record.GetRegionID(),
			WeatherType: string(record.GetWeatherType()),
			Intensity:   record.GetIntensity(),
			Temperature: record.GetTemperature(),
			StartTime:   record.GetStartTime(),
			EndTime:     record.GetEndTime(),
			Duration:    record.GetDuration(),
			IsSpecial:   record.IsSpecialWeather(),
		}
	}
	return dtos
}

// buildStatisticsDTO 构建统计DTO
func (s *WeatherService) buildStatisticsDTO(stats *weather.WeatherStatistics) *WeatherStatisticsDTO {
	return &WeatherStatisticsDTO{
		RegionID:           stats.GetRegionID(),
		TotalRecords:       stats.GetTotalRecords(),
		WeatherTypeStats:   stats.GetWeatherTypeStats(),
		AverageTemperature: stats.GetAverageTemperature(),
		AverageHumidity:    stats.GetAverageHumidity(),
		AverageWindSpeed:   stats.GetAverageWindSpeed(),
		MostCommonWeather:  string(stats.GetMostCommonWeather()),
		SpecialWeatherCount: stats.GetSpecialWeatherCount(),
		LastWeatherTime:    stats.GetLastWeatherTime(),
	}
}

// buildGlobalWeatherDTO 构建全局天气DTO
func (s *WeatherService) buildGlobalWeatherDTO(allWeather []*weather.WeatherAggregate) *GlobalWeatherDTO {
	regionWeather := make(map[string]*WeatherDTO)
	weatherTypeCount := make(map[string]int)
	totalRegions := len(allWeather)
	specialWeatherCount := 0
	
	for _, w := range allWeather {
		regionWeather[w.GetRegionID()] = s.buildWeatherDTO(w)
		weatherType := string(w.GetWeatherType())
		weatherTypeCount[weatherType]++
		if w.IsSpecialWeather() {
			specialWeatherCount++
		}
	}
	
	return &GlobalWeatherDTO{
		RegionWeather:       regionWeather,
		TotalRegions:        totalRegions,
		WeatherTypeCount:    weatherTypeCount,
		SpecialWeatherCount: specialWeatherCount,
		LastUpdateTime:      time.Now(),
	}
}

// buildSeasonInfoDTO 构建季节信息DTO
func (s *WeatherService) buildSeasonInfoDTO(seasonInfo *weather.SeasonInfo) *SeasonInfoDTO {
	return &SeasonInfoDTO{
		CurrentSeason:    string(seasonInfo.GetCurrentSeason()),
		SeasonProgress:   seasonInfo.GetSeasonProgress(),
		DaysRemaining:    seasonInfo.GetDaysRemaining(),
		NextSeason:       string(seasonInfo.GetNextSeason()),
		SeasonEffects:    seasonInfo.GetSeasonEffects(),
		TemperatureRange: seasonInfo.GetTemperatureRange(),
		WeatherTendency:  seasonInfo.GetWeatherTendency(),
	}
}

// DTO 定义

// WeatherDTO 天气DTO
type WeatherDTO struct {
	RegionID    string        `json:"region_id"`
	WeatherType string        `json:"weather_type"`
	Intensity   float64       `json:"intensity"`
	Temperature float64       `json:"temperature"`
	Humidity    float64       `json:"humidity"`
	WindSpeed   float64       `json:"wind_speed"`
	Visibility  float64       `json:"visibility"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	Duration    time.Duration `json:"duration"`
	IsSpecial   bool          `json:"is_special"`
	Description string        `json:"description"`
}

// WeatherForecastDTO 天气预报DTO
type WeatherForecastDTO struct {
	RegionID     string    `json:"region_id"`
	ForecastDate time.Time `json:"forecast_date"`
	WeatherType  string    `json:"weather_type"`
	Intensity    float64   `json:"intensity"`
	Temperature  float64   `json:"temperature"`
	Humidity     float64   `json:"humidity"`
	WindSpeed    float64   `json:"wind_speed"`
	Probability  float64   `json:"probability"`
	Description  string    `json:"description"`
}

// WeatherEffectDTO 天气影响DTO
type WeatherEffectDTO struct {
	EffectType  string        `json:"effect_type"`
	TargetType  string        `json:"target_type"`
	Modifier    float64       `json:"modifier"`
	Duration    time.Duration `json:"duration"`
	IsPositive  bool          `json:"is_positive"`
	Description string        `json:"description"`
}

// WeatherHistoryDTO 天气历史DTO
type WeatherHistoryDTO struct {
	RegionID    string        `json:"region_id"`
	WeatherType string        `json:"weather_type"`
	Intensity   float64       `json:"intensity"`
	Temperature float64       `json:"temperature"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	Duration    time.Duration `json:"duration"`
	IsSpecial   bool          `json:"is_special"`
}

// WeatherStatisticsDTO 天气统计DTO
type WeatherStatisticsDTO struct {
	RegionID            string             `json:"region_id"`
	TotalRecords        int64              `json:"total_records"`
	WeatherTypeStats    map[string]int64   `json:"weather_type_stats"`
	AverageTemperature  float64            `json:"average_temperature"`
	AverageHumidity     float64            `json:"average_humidity"`
	AverageWindSpeed    float64            `json:"average_wind_speed"`
	MostCommonWeather   string             `json:"most_common_weather"`
	SpecialWeatherCount int64              `json:"special_weather_count"`
	LastWeatherTime     time.Time          `json:"last_weather_time"`
}

// GlobalWeatherDTO 全局天气DTO
type GlobalWeatherDTO struct {
	RegionWeather       map[string]*WeatherDTO `json:"region_weather"`
	TotalRegions        int                    `json:"total_regions"`
	WeatherTypeCount    map[string]int         `json:"weather_type_count"`
	SpecialWeatherCount int                    `json:"special_weather_count"`
	LastUpdateTime      time.Time              `json:"last_update_time"`
}

// SeasonInfoDTO 季节信息DTO
type SeasonInfoDTO struct {
	CurrentSeason    string                 `json:"current_season"`
	SeasonProgress   float64                `json:"season_progress"`
	DaysRemaining    int                    `json:"days_remaining"`
	NextSeason       string                 `json:"next_season"`
	SeasonEffects    map[string]float64     `json:"season_effects"`
	TemperatureRange map[string]float64     `json:"temperature_range"`
	WeatherTendency  map[string]float64     `json:"weather_tendency"`
}