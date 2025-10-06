package weather

import (
	// "fmt"
	"math"
	"math/rand"
	"time"
)

// WeatherService 天气领域服务
type WeatherService struct {
	weatherTemplates   map[WeatherType]*WeatherTemplate
	seasonalPatterns   map[string]*SeasonalPattern
	weatherEvents      map[WeatherEventType]*WeatherEventTemplate
	effectCalculators  map[string]EffectCalculator
	forecastGenerators map[string]ForecastGenerator
	weatherRules       []*WeatherRule
	climateZones       map[string]*ClimateZone
	randomSeed         int64
	createdAt          time.Time
	updatedAt          time.Time
}

// NewWeatherService 创建天气服务
func NewWeatherService() *WeatherService {
	now := time.Now()
	service := &WeatherService{
		weatherTemplates:   make(map[WeatherType]*WeatherTemplate),
		seasonalPatterns:   make(map[string]*SeasonalPattern),
		weatherEvents:      make(map[WeatherEventType]*WeatherEventTemplate),
		effectCalculators:  make(map[string]EffectCalculator),
		forecastGenerators: make(map[string]ForecastGenerator),
		weatherRules:       make([]*WeatherRule, 0),
		climateZones:       make(map[string]*ClimateZone),
		randomSeed:         now.UnixNano(),
		createdAt:          now,
		updatedAt:          now,
	}

	// 初始化默认模板和规则
	service.initializeDefaultTemplates()
	service.initializeDefaultRules()
	service.initializeDefaultClimateZones()
	service.initializeEffectCalculators()
	service.initializeForecastGenerators()

	return service
}

// RegisterWeatherTemplate 注册天气模板
func (ws *WeatherService) RegisterWeatherTemplate(weatherType WeatherType, template *WeatherTemplate) {
	ws.weatherTemplates[weatherType] = template
	ws.updatedAt = time.Now()
}

// GetWeatherTemplate 获取天气模板
func (ws *WeatherService) GetWeatherTemplate(weatherType WeatherType) *WeatherTemplate {
	return ws.weatherTemplates[weatherType]
}

// RegisterSeasonalPattern 注册季节模式
func (ws *WeatherService) RegisterSeasonalPattern(zoneID string, pattern *SeasonalPattern) {
	ws.seasonalPatterns[zoneID] = pattern
	ws.updatedAt = time.Now()
}

// GetSeasonalPattern 获取季节模式
func (ws *WeatherService) GetSeasonalPattern(zoneID string) *SeasonalPattern {
	return ws.seasonalPatterns[zoneID]
}

// RegisterWeatherEvent 注册天气事件
func (ws *WeatherService) RegisterWeatherEvent(eventType WeatherEventType, template *WeatherEventTemplate) {
	ws.weatherEvents[eventType] = template
	ws.updatedAt = time.Now()
}

// GetWeatherEventTemplate 获取天气事件模板
func (ws *WeatherService) GetWeatherEventTemplate(eventType WeatherEventType) *WeatherEventTemplate {
	return ws.weatherEvents[eventType]
}

// AddWeatherRule 添加天气规则
func (ws *WeatherService) AddWeatherRule(rule *WeatherRule) {
	ws.weatherRules = append(ws.weatherRules, rule)
	ws.updatedAt = time.Now()
}

// GetWeatherRules 获取天气规则
func (ws *WeatherService) GetWeatherRules() []*WeatherRule {
	return ws.weatherRules
}

// RegisterClimateZone 注册气候区域
func (ws *WeatherService) RegisterClimateZone(zoneID string, zone *ClimateZone) {
	ws.climateZones[zoneID] = zone
	ws.updatedAt = time.Now()
}

// GetClimateZone 获取气候区域
func (ws *WeatherService) GetClimateZone(zoneID string) *ClimateZone {
	return ws.climateZones[zoneID]
}

// CalculateNextWeather 计算下一个天气
func (ws *WeatherService) CalculateNextWeather(currentWeather *WeatherState, zoneID string) (*WeatherState, error) {
	if currentWeather == nil {
		return nil, ErrInvalidWeatherState
	}

	// 获取气候区域和季节模式
	climateZone := ws.GetClimateZone(zoneID)
	if climateZone == nil {
		climateZone = ws.getDefaultClimateZone()
	}

	seasonalPattern := ws.GetSeasonalPattern(zoneID)
	if seasonalPattern == nil {
		seasonalPattern = NewSeasonalPattern()
	}

	// 获取当前季节
	currentSeason := seasonalPattern.GetCurrentSeason(time.Now())

	// 获取天气转换概率
	transitionProbs := ws.calculateWeatherTransitionProbabilities(currentWeather.WeatherType, currentSeason, climateZone)

	// 应用天气规则
	adjustedProbs := ws.applyWeatherRules(transitionProbs, currentWeather, climateZone)

	// 选择下一个天气
	nextWeatherType := ws.selectWeatherByProbability(adjustedProbs)
	nextIntensity := ws.calculateWeatherIntensity(nextWeatherType, currentSeason, climateZone)

	// 创建新的天气状态
	nextWeather := NewWeatherState(nextWeatherType, nextIntensity)

	// 应用气候区域的影响
	ws.applyClimateZoneEffects(nextWeather, climateZone)

	return nextWeather, nil
}

// GenerateWeatherForecast 生成天气预报
func (ws *WeatherService) GenerateWeatherForecast(currentWeather *WeatherState, zoneID string, hours int) ([]*WeatherForecast, error) {
	if currentWeather == nil {
		return nil, ErrInvalidWeatherState
	}

	if hours <= 0 || hours > 168 {
		return nil, ErrInvalidForecastPeriod
	}

	forecasts := make([]*WeatherForecast, 0, hours)
	currentTime := time.Now()
	predictedWeather := currentWeather

	for i := 1; i <= hours; i++ {
		forecastTime := currentTime.Add(time.Duration(i) * time.Hour)

		// 预测下一个小时的天气
		nextWeather, err := ws.CalculateNextWeather(predictedWeather, zoneID)
		if err != nil {
			return nil, err
		}

		// 计算预报置信度
		confidence := ws.calculateForecastConfidence(i)

		// 创建预报
		forecast := NewWeatherForecast(forecastTime, nextWeather.WeatherType, nextWeather.Intensity)
		forecast.Temperature = nextWeather.Temperature
		forecast.Humidity = nextWeather.Humidity
		forecast.WindSpeed = nextWeather.WindSpeed
		forecast.Visibility = nextWeather.Visibility
		forecast.Confidence = confidence

		forecasts = append(forecasts, forecast)
		predictedWeather = nextWeather
	}

	return forecasts, nil
}

// CalculateWeatherEffects 计算天气效果
func (ws *WeatherService) CalculateWeatherEffects(weather *WeatherState, targetType string) ([]*WeatherEffect, error) {
	if weather == nil {
		return nil, ErrInvalidWeatherState
	}

	effects := make([]*WeatherEffect, 0)

	// 获取天气模板
	template := ws.GetWeatherTemplate(weather.WeatherType)
	if template == nil {
		return effects, nil
	}

	// 根据目标类型计算效果
	for effectType, baseMultiplier := range template.BaseEffects {
		if targetType != "" && effectType != targetType {
			continue
		}

		// 应用强度调整
		adjustedMultiplier := baseMultiplier * weather.Intensity.GetMultiplier()

		// 创建效果
		effect := NewWeatherEffect(effectType, adjustedMultiplier, weather.Duration)
		effects = append(effects, effect)
	}

	return effects, nil
}

// CheckWeatherEventTrigger 检查天气事件触发
func (ws *WeatherService) CheckWeatherEventTrigger(weather *WeatherState, zoneID string) (*WeatherEvent, error) {
	if weather == nil {
		return nil, ErrInvalidWeatherState
	}

	// 检查每种天气事件的触发条件
	for eventType, template := range ws.weatherEvents {
		if ws.shouldTriggerWeatherEvent(weather, eventType, template, zoneID) {
			// 创建天气事件
			event := ws.createWeatherEvent(eventType, template, weather)
			return event, nil
		}
	}

	return nil, nil
}

// CalculateWeatherInfluence 计算天气对属性的影响
func (ws *WeatherService) CalculateWeatherInfluence(weather *WeatherState, attributeType string) (float64, error) {
	if weather == nil {
		return 1.0, ErrInvalidWeatherState
	}

	// 使用效果计算器
	if calculator, exists := ws.effectCalculators[attributeType]; exists {
		return calculator.Calculate(weather), nil
	}

	// 默认计算逻辑
	return ws.calculateDefaultInfluence(weather, attributeType), nil
}

// ValidateWeatherTransition 验证天气转换
func (ws *WeatherService) ValidateWeatherTransition(from, to WeatherType, intensity WeatherIntensity) error {
	if !from.IsValid() || !to.IsValid() {
		return ErrInvalidWeatherType
	}

	if !intensity.IsValid() {
		return ErrInvalidWeatherIntensity
	}

	// 检查转换规则
	for _, rule := range ws.weatherRules {
		if rule.FromWeather == from && rule.ToWeather == to {
			if rule.MinIntensity != 0 && intensity < rule.MinIntensity {
				return ErrInvalidWeatherTransition
			}
			if rule.MaxIntensity != 0 && intensity > rule.MaxIntensity {
				return ErrInvalidWeatherTransition
			}
			return nil
		}
	}

	return nil
}

// GetOptimalWeatherForActivity 获取活动的最佳天气
func (ws *WeatherService) GetOptimalWeatherForActivity(activityType string) (*WeatherCondition, error) {
	// 根据活动类型返回最佳天气条件
	switch activityType {
	case "farming":
		return NewWeatherCondition(WeatherTypeRainy, WeatherIntensityLight, 2*time.Hour), nil
	case "combat":
		return NewWeatherCondition(WeatherTypeSunny, WeatherIntensityNormal, 4*time.Hour), nil
	case "exploration":
		return NewWeatherCondition(WeatherTypeCloudy, WeatherIntensityNormal, 3*time.Hour), nil
	case "crafting":
		return NewWeatherCondition(WeatherTypeSunny, WeatherIntensityLight, 6*time.Hour), nil
	default:
		return NewWeatherCondition(WeatherTypeSunny, WeatherIntensityNormal, 2*time.Hour), nil
	}
}

// 私有方法

// initializeDefaultTemplates 初始化默认模板
func (ws *WeatherService) initializeDefaultTemplates() {
	// 晴天模板
	sunnyTemplate := &WeatherTemplate{
		WeatherType: WeatherTypeSunny,
		BaseEffects: map[string]float64{
			"visibility":     1.2,
			"movement_speed": 1.1,
			"energy_regen":   1.15,
			"mood":           1.1,
		},
		DurationRange: DurationRange{Min: 2 * time.Hour, Max: 8 * time.Hour},
		TransitionRules: map[WeatherType]float64{
			WeatherTypeSunny:  0.6,
			WeatherTypeCloudy: 0.3,
			WeatherTypeRainy:  0.1,
		},
	}
	ws.RegisterWeatherTemplate(WeatherTypeSunny, sunnyTemplate)

	// 雨天模板
	rainyTemplate := &WeatherTemplate{
		WeatherType: WeatherTypeRainy,
		BaseEffects: map[string]float64{
			"visibility":   0.8,
			"fire_damage":  0.7,
			"water_damage": 1.3,
			"plant_growth": 1.5,
		},
		DurationRange: DurationRange{Min: 1 * time.Hour, Max: 4 * time.Hour},
		TransitionRules: map[WeatherType]float64{
			WeatherTypeRainy:  0.5,
			WeatherTypeCloudy: 0.3,
			WeatherTypeStormy: 0.2,
		},
	}
	ws.RegisterWeatherTemplate(WeatherTypeRainy, rainyTemplate)

	// 暴风雨模板
	stormyTemplate := &WeatherTemplate{
		WeatherType: WeatherTypeStormy,
		BaseEffects: map[string]float64{
			"visibility":       0.5,
			"lightning_damage": 1.8,
			"movement_speed":   0.7,
			"accuracy":         0.8,
		},
		DurationRange: DurationRange{Min: 30 * time.Minute, Max: 2 * time.Hour},
		TransitionRules: map[WeatherType]float64{
			WeatherTypeStormy: 0.3,
			WeatherTypeRainy:  0.5,
			WeatherTypeCloudy: 0.2,
		},
	}
	ws.RegisterWeatherTemplate(WeatherTypeStormy, stormyTemplate)
}

// initializeDefaultRules 初始化默认规则
func (ws *WeatherService) initializeDefaultRules() {
	// 添加一些基本的天气转换规则
	ws.AddWeatherRule(&WeatherRule{
		FromWeather:  WeatherTypeSunny,
		ToWeather:    WeatherTypeStormy,
		Probability:  0.05, // 晴天直接转暴风雨概率很低
		MinIntensity: WeatherIntensityNormal,
		SeasonFactor: map[Season]float64{SeasonSummer: 0.1, SeasonWinter: 0.01},
	})

	ws.AddWeatherRule(&WeatherRule{
		FromWeather:  WeatherTypeRainy,
		ToWeather:    WeatherTypeSnowy,
		Probability:  0.3,
		MinIntensity: WeatherIntensityLight,
		SeasonFactor: map[Season]float64{SeasonWinter: 0.8, SeasonSummer: 0.0},
	})
}

// initializeDefaultClimateZones 初始化默认气候区域
func (ws *WeatherService) initializeDefaultClimateZones() {
	// 温带气候
	temperateZone := &ClimateZone{
		ZoneID:           "temperate",
		Name:             "温带气候",
		Description:      "四季分明的温带气候区域",
		BaseTemperature:  15.0,
		TemperatureRange: TemperatureRange{Min: -10, Max: 35, Average: 15},
		HumidityRange:    HumidityRange{Min: 30, Max: 80, Average: 55},
		WeatherModifiers: map[WeatherType]float64{
			WeatherTypeSunny: 1.0,
			WeatherTypeRainy: 1.0,
			WeatherTypeSnowy: 0.8,
		},
	}
	ws.RegisterClimateZone("temperate", temperateZone)

	// 热带气候
	tropicalZone := &ClimateZone{
		ZoneID:           "tropical",
		Name:             "热带气候",
		Description:      "高温多雨的热带气候区域",
		BaseTemperature:  28.0,
		TemperatureRange: TemperatureRange{Min: 20, Max: 40, Average: 28},
		HumidityRange:    HumidityRange{Min: 60, Max: 95, Average: 75},
		WeatherModifiers: map[WeatherType]float64{
			WeatherTypeSunny:  1.2,
			WeatherTypeRainy:  1.5,
			WeatherTypeStormy: 1.3,
			WeatherTypeSnowy:  0.0, // 热带不下雪
		},
	}
	ws.RegisterClimateZone("tropical", tropicalZone)
}

// initializeEffectCalculators 初始化效果计算器
func (ws *WeatherService) initializeEffectCalculators() {
	// 能见度计算器
	ws.effectCalculators["visibility"] = EffectCalculatorFunc(func(weather *WeatherState) float64 {
		base := weather.WeatherType.GetBaseVisibility() / 10.0 // 标准化到0-2范围
		intensityFactor := weather.Intensity.GetMultiplier()
		return math.Max(0.1, base*intensityFactor)
	})

	// 移动速度计算器
	ws.effectCalculators["movement_speed"] = EffectCalculatorFunc(func(weather *WeatherState) float64 {
		switch weather.WeatherType {
		case WeatherTypeSunny:
			return 1.0 + 0.1*weather.Intensity.GetMultiplier()
		case WeatherTypeSnowy, WeatherTypeStormy:
			return 1.0 - 0.2*weather.Intensity.GetMultiplier()
		default:
			return 1.0
		}
	})
}

// initializeForecastGenerators 初始化预报生成器
func (ws *WeatherService) initializeForecastGenerators() {
	// 短期预报生成器（1-6小时）
	ws.forecastGenerators["short_term"] = ForecastGeneratorFunc(func(current *WeatherState, hours int) []*WeatherForecast {
		// 实现短期预报逻辑
		return make([]*WeatherForecast, 0)
	})

	// 长期预报生成器（1-7天）
	ws.forecastGenerators["long_term"] = ForecastGeneratorFunc(func(current *WeatherState, hours int) []*WeatherForecast {
		// 实现长期预报逻辑
		return make([]*WeatherForecast, 0)
	})
}

// calculateWeatherTransitionProbabilities 计算天气转换概率
func (ws *WeatherService) calculateWeatherTransitionProbabilities(currentType WeatherType, season Season, zone *ClimateZone) map[WeatherType]float64 {
	probs := make(map[WeatherType]float64)

	// 获取基础转换概率
	template := ws.GetWeatherTemplate(currentType)
	if template != nil {
		for weatherType, prob := range template.TransitionRules {
			probs[weatherType] = prob
		}
	}

	// 应用气候区域修正
	if zone != nil {
		for weatherType, modifier := range zone.WeatherModifiers {
			if prob, exists := probs[weatherType]; exists {
				probs[weatherType] = prob * modifier
			}
		}
	}

	// 标准化概率
	total := 0.0
	for _, prob := range probs {
		total += prob
	}

	if total > 0 {
		for weatherType := range probs {
			probs[weatherType] /= total
		}
	}

	return probs
}

// applyWeatherRules 应用天气规则
func (ws *WeatherService) applyWeatherRules(probs map[WeatherType]float64, current *WeatherState, zone *ClimateZone) map[WeatherType]float64 {
	adjusted := make(map[WeatherType]float64)
	for k, v := range probs {
		adjusted[k] = v
	}

	// 应用每个规则
	for _, rule := range ws.weatherRules {
		if rule.FromWeather == current.WeatherType {
			if prob, exists := adjusted[rule.ToWeather]; exists {
				// 应用规则修正
				adjusted[rule.ToWeather] = prob * rule.Probability
			}
		}
	}

	return adjusted
}

// selectWeatherByProbability 根据概率选择天气
func (ws *WeatherService) selectWeatherByProbability(probs map[WeatherType]float64) WeatherType {
	rand.Seed(time.Now().UnixNano() + ws.randomSeed)
	randomValue := rand.Float64()

	cumulativeProb := 0.0
	for weatherType, prob := range probs {
		cumulativeProb += prob
		if randomValue <= cumulativeProb {
			return weatherType
		}
	}

	// 默认返回晴天
	return WeatherTypeSunny
}

// calculateWeatherIntensity 计算天气强度
func (ws *WeatherService) calculateWeatherIntensity(weatherType WeatherType, season Season, zone *ClimateZone) WeatherIntensity {
	rand.Seed(time.Now().UnixNano() + ws.randomSeed)

	// 基础强度概率
	intensityProbs := map[WeatherIntensity]float64{
		WeatherIntensityLight:  0.3,
		WeatherIntensityNormal: 0.5,
		WeatherIntensityHeavy:  0.2,
	}

	// 根据天气类型调整
	switch weatherType {
	case WeatherTypeStormy:
		intensityProbs[WeatherIntensityHeavy] = 0.4
		intensityProbs[WeatherIntensityExtreme] = 0.1
	case WeatherTypeSunny:
		intensityProbs[WeatherIntensityLight] = 0.4
		intensityProbs[WeatherIntensityNormal] = 0.6
		intensityProbs[WeatherIntensityHeavy] = 0.0
	}

	// 选择强度
	randomValue := rand.Float64()
	cumulativeProb := 0.0

	for intensity, prob := range intensityProbs {
		cumulativeProb += prob
		if randomValue <= cumulativeProb {
			return intensity
		}
	}

	return WeatherIntensityNormal
}

// applyClimateZoneEffects 应用气候区域效果
func (ws *WeatherService) applyClimateZoneEffects(weather *WeatherState, zone *ClimateZone) {
	if zone == nil {
		return
	}

	// 调整温度
	temperatureOffset := zone.BaseTemperature - weather.WeatherType.GetBaseTemperature()
	weather.UpdateTemperature(weather.Temperature + temperatureOffset)

	// 调整湿度
	if weather.Humidity < zone.HumidityRange.Min {
		weather.UpdateHumidity(zone.HumidityRange.Min)
	} else if weather.Humidity > zone.HumidityRange.Max {
		weather.UpdateHumidity(zone.HumidityRange.Max)
	}
}

// calculateForecastConfidence 计算预报置信度
func (ws *WeatherService) calculateForecastConfidence(hoursAhead int) float64 {
	baseConfidence := 0.95
	decayRate := 0.02

	confidence := baseConfidence - float64(hoursAhead)*decayRate
	if confidence < 0.3 {
		confidence = 0.3
	}

	return confidence
}

// shouldTriggerWeatherEvent 检查是否应该触发天气事件
func (ws *WeatherService) shouldTriggerWeatherEvent(weather *WeatherState, eventType WeatherEventType, template *WeatherEventTemplate, zoneID string) bool {
	// 检查天气类型匹配
	if !template.CanTriggerWith(weather.WeatherType, weather.Intensity) {
		return false
	}

	// 检查概率
	rand.Seed(time.Now().UnixNano() + ws.randomSeed)
	return rand.Float64() < template.TriggerProbability
}

// createWeatherEvent 创建天气事件
func (ws *WeatherService) createWeatherEvent(eventType WeatherEventType, template *WeatherEventTemplate, weather *WeatherState) *WeatherEvent {
	severity := ws.calculateEventSeverity(weather.Intensity)
	duration := template.BaseDuration

	event := NewWeatherEvent(eventType, severity, template.Title, template.Description, duration)

	// 添加效果
	for effectType, multiplier := range template.Effects {
		effect := NewWeatherEffect(effectType, multiplier*weather.Intensity.GetMultiplier(), duration)
		event.AddEffect(effect)
	}

	return event
}

// calculateEventSeverity 计算事件严重程度
func (ws *WeatherService) calculateEventSeverity(intensity WeatherIntensity) WeatherEventSeverity {
	switch intensity {
	case WeatherIntensityLight:
		return WeatherEventSeverityMinor
	case WeatherIntensityNormal:
		return WeatherEventSeverityModerate
	case WeatherIntensityHeavy:
		return WeatherEventSeverityMajor
	case WeatherIntensityExtreme:
		return WeatherEventSeverityCritical
	default:
		return WeatherEventSeverityMinor
	}
}

// calculateDefaultInfluence 计算默认影响
func (ws *WeatherService) calculateDefaultInfluence(weather *WeatherState, attributeType string) float64 {
	// 默认的影响计算逻辑
	switch attributeType {
	case "visibility":
		return weather.Visibility / 10.0 // 标准化
	case "movement_speed":
		if weather.WeatherType == WeatherTypeSnowy || weather.WeatherType == WeatherTypeStormy {
			return 0.8
		}
		return 1.0
	default:
		return 1.0
	}
}

// getDefaultClimateZone 获取默认气候区域
func (ws *WeatherService) getDefaultClimateZone() *ClimateZone {
	return ws.GetClimateZone("temperate")
}

// 辅助类型和接口

// WeatherTemplate 天气模板
type WeatherTemplate struct {
	WeatherType     WeatherType
	BaseEffects     map[string]float64
	DurationRange   DurationRange
	TransitionRules map[WeatherType]float64
}

// DurationRange 持续时间范围
type DurationRange struct {
	Min time.Duration
	Max time.Duration
}

// WeatherRule 天气规则
type WeatherRule struct {
	FromWeather  WeatherType
	ToWeather    WeatherType
	Probability  float64
	MinIntensity WeatherIntensity
	MaxIntensity WeatherIntensity
	SeasonFactor map[Season]float64
	Conditions   []string
}

// ClimateZone 气候区域
type ClimateZone struct {
	ZoneID           string
	Name             string
	Description      string
	BaseTemperature  float64
	TemperatureRange TemperatureRange
	HumidityRange    HumidityRange
	WeatherModifiers map[WeatherType]float64
}

// HumidityRange 湿度范围
type HumidityRange struct {
	Min     float64
	Max     float64
	Average float64
}

// WeatherEventTemplate 天气事件模板
type WeatherEventTemplate struct {
	EventType          WeatherEventType
	Title              string
	Description        string
	TriggerProbability float64
	BaseDuration       time.Duration
	Effects            map[string]float64
	TriggerConditions  []WeatherCondition
}

// CanTriggerWith 检查是否可以触发
func (wet *WeatherEventTemplate) CanTriggerWith(weatherType WeatherType, intensity WeatherIntensity) bool {
	for _, condition := range wet.TriggerConditions {
		if condition.Matches(weatherType, intensity) {
			return true
		}
	}
	return len(wet.TriggerConditions) == 0 // 如果没有条件，默认可以触发
}

// EffectCalculator 效果计算器接口
type EffectCalculator interface {
	Calculate(weather *WeatherState) float64
}

// EffectCalculatorFunc 效果计算器函数类型
type EffectCalculatorFunc func(weather *WeatherState) float64

// Calculate 实现EffectCalculator接口
func (f EffectCalculatorFunc) Calculate(weather *WeatherState) float64 {
	return f(weather)
}

// ForecastGenerator 预报生成器接口
type ForecastGenerator interface {
	Generate(current *WeatherState, hours int) []*WeatherForecast
}

// ForecastGeneratorFunc 预报生成器函数类型
type ForecastGeneratorFunc func(current *WeatherState, hours int) []*WeatherForecast

// Generate 实现ForecastGenerator接口
func (f ForecastGeneratorFunc) Generate(current *WeatherState, hours int) []*WeatherForecast {
	return f(current, hours)
}
