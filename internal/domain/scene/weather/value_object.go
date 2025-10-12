package weather

import (
	"fmt"
	"time"
)

// WeatherType 天气类型
type WeatherType int

const (
	WeatherTypeSunny WeatherType = iota + 1
	WeatherTypeCloudy
	WeatherTypeRainy
	WeatherTypeSnowy
	WeatherTypeWindy
	WeatherTypeStormy
	WeatherTypeFoggy
	WeatherTypeHazy
	WeatherTypeHail
	WeatherTypeBlizzard
)

// String 返回天气类型字符串
func (wt WeatherType) String() string {
	switch wt {
	case WeatherTypeSunny:
		return "sunny"
	case WeatherTypeCloudy:
		return "cloudy"
	case WeatherTypeRainy:
		return "rainy"
	case WeatherTypeSnowy:
		return "snowy"
	case WeatherTypeWindy:
		return "windy"
	case WeatherTypeStormy:
		return "stormy"
	case WeatherTypeFoggy:
		return "foggy"
	case WeatherTypeHazy:
		return "hazy"
	case WeatherTypeHail:
		return "hail"
	case WeatherTypeBlizzard:
		return "blizzard"
	default:
		return "unknown"
	}
}

// GetDescription 获取天气描述
func (wt WeatherType) GetDescription() string {
	switch wt {
	case WeatherTypeSunny:
		return "晴朗"
	case WeatherTypeCloudy:
		return "多云"
	case WeatherTypeRainy:
		return "下雨"
	case WeatherTypeSnowy:
		return "下雪"
	case WeatherTypeWindy:
		return "大风"
	case WeatherTypeStormy:
		return "暴风雨"
	case WeatherTypeFoggy:
		return "雾天"
	case WeatherTypeHazy:
		return "霾天"
	case WeatherTypeHail:
		return "冰雹"
	case WeatherTypeBlizzard:
		return "暴雪"
	default:
		return "未知天气"
	}
}

// IsValid 检查天气类型是否有效
func (wt WeatherType) IsValid() bool {
	return wt >= WeatherTypeSunny && wt <= WeatherTypeBlizzard
}

// ParseWeatherType 从字符串解析天气类型
func ParseWeatherType(s string) WeatherType {
	switch s {
	case "sunny":
		return WeatherTypeSunny
	case "cloudy":
		return WeatherTypeCloudy
	case "rainy":
		return WeatherTypeRainy
	case "snowy":
		return WeatherTypeSnowy
	case "windy":
		return WeatherTypeWindy
	case "stormy":
		return WeatherTypeStormy
	case "foggy":
		return WeatherTypeFoggy
	case "hazy":
		return WeatherTypeHazy
	case "hail":
		return WeatherTypeHail
	case "blizzard":
		return WeatherTypeBlizzard
	default:
		return WeatherTypeSunny // 默认返回晴天
	}
}

// GetBaseTemperature 获取基础温度
func (wt WeatherType) GetBaseTemperature() float64 {
	switch wt {
	case WeatherTypeSunny:
		return 25.0
	case WeatherTypeCloudy:
		return 20.0
	case WeatherTypeRainy:
		return 15.0
	case WeatherTypeSnowy:
		return -5.0
	case WeatherTypeWindy:
		return 18.0
	case WeatherTypeStormy:
		return 12.0
	case WeatherTypeFoggy:
		return 10.0
	case WeatherTypeHazy:
		return 22.0
	case WeatherTypeHail:
		return 8.0
	case WeatherTypeBlizzard:
		return -15.0
	default:
		return 20.0
	}
}

// GetBaseHumidity 获取基础湿度
func (wt WeatherType) GetBaseHumidity() float64 {
	switch wt {
	case WeatherTypeSunny:
		return 40.0
	case WeatherTypeCloudy:
		return 60.0
	case WeatherTypeRainy:
		return 85.0
	case WeatherTypeSnowy:
		return 70.0
	case WeatherTypeWindy:
		return 50.0
	case WeatherTypeStormy:
		return 90.0
	case WeatherTypeFoggy:
		return 95.0
	case WeatherTypeHazy:
		return 65.0
	case WeatherTypeHail:
		return 80.0
	case WeatherTypeBlizzard:
		return 85.0
	default:
		return 50.0
	}
}

// GetBaseWindSpeed 获取基础风速
func (wt WeatherType) GetBaseWindSpeed() float64 {
	switch wt {
	case WeatherTypeSunny:
		return 5.0
	case WeatherTypeCloudy:
		return 10.0
	case WeatherTypeRainy:
		return 15.0
	case WeatherTypeSnowy:
		return 20.0
	case WeatherTypeWindy:
		return 35.0
	case WeatherTypeStormy:
		return 50.0
	case WeatherTypeFoggy:
		return 3.0
	case WeatherTypeHazy:
		return 8.0
	case WeatherTypeHail:
		return 25.0
	case WeatherTypeBlizzard:
		return 60.0
	default:
		return 10.0
	}
}

// GetBaseVisibility 获取基础能见度
func (wt WeatherType) GetBaseVisibility() float64 {
	switch wt {
	case WeatherTypeSunny:
		return 20.0
	case WeatherTypeCloudy:
		return 15.0
	case WeatherTypeRainy:
		return 8.0
	case WeatherTypeSnowy:
		return 5.0
	case WeatherTypeWindy:
		return 12.0
	case WeatherTypeStormy:
		return 3.0
	case WeatherTypeFoggy:
		return 1.0
	case WeatherTypeHazy:
		return 6.0
	case WeatherTypeHail:
		return 4.0
	case WeatherTypeBlizzard:
		return 2.0
	default:
		return 10.0
	}
}

// WeatherIntensity 天气强度
type WeatherIntensity int

const (
	WeatherIntensityLight WeatherIntensity = iota + 1
	WeatherIntensityNormal
	WeatherIntensityHeavy
	WeatherIntensityExtreme
)

// String 返回强度字符串
func (wi WeatherIntensity) String() string {
	switch wi {
	case WeatherIntensityLight:
		return "light"
	case WeatherIntensityNormal:
		return "normal"
	case WeatherIntensityHeavy:
		return "heavy"
	case WeatherIntensityExtreme:
		return "extreme"
	default:
		return "unknown"
	}
}

// GetDescription 获取强度描述
func (wi WeatherIntensity) GetDescription() string {
	switch wi {
	case WeatherIntensityLight:
		return "轻微"
	case WeatherIntensityNormal:
		return "正常"
	case WeatherIntensityHeavy:
		return "强烈"
	case WeatherIntensityExtreme:
		return "极端"
	default:
		return "未知强度"
	}
}

// IsValid 检查强度是否有效
func (wi WeatherIntensity) IsValid() bool {
	return wi >= WeatherIntensityLight && wi <= WeatherIntensityExtreme
}

// GetMultiplier 获取强度倍率
func (wi WeatherIntensity) GetMultiplier() float64 {
	switch wi {
	case WeatherIntensityLight:
		return 0.5
	case WeatherIntensityNormal:
		return 1.0
	case WeatherIntensityHeavy:
		return 1.5
	case WeatherIntensityExtreme:
		return 2.0
	default:
		return 1.0
	}
}

// GetDurationFactor 获取持续时间因子
func (wi WeatherIntensity) GetDurationFactor() float64 {
	switch wi {
	case WeatherIntensityLight:
		return 1.5 // 轻微天气持续更久
	case WeatherIntensityNormal:
		return 1.0
	case WeatherIntensityHeavy:
		return 0.7 // 强烈天气持续较短
	case WeatherIntensityExtreme:
		return 0.5 // 极端天气持续很短
	default:
		return 1.0
	}
}

// Season 季节
type Season int

const (
	SeasonSpring Season = iota + 1
	SeasonSummer
	SeasonAutumn
	SeasonWinter
)

// String 返回季节字符串
func (s Season) String() string {
	switch s {
	case SeasonSpring:
		return "spring"
	case SeasonSummer:
		return "summer"
	case SeasonAutumn:
		return "autumn"
	case SeasonWinter:
		return "winter"
	default:
		return "unknown"
	}
}

// GetDescription 获取季节描述
func (s Season) GetDescription() string {
	switch s {
	case SeasonSpring:
		return "春季"
	case SeasonSummer:
		return "夏季"
	case SeasonAutumn:
		return "秋季"
	case SeasonWinter:
		return "冬季"
	default:
		return "未知季节"
	}
}

// IsValid 检查季节是否有效
func (s Season) IsValid() bool {
	return s >= SeasonSpring && s <= SeasonWinter
}

// SeasonalPattern 季节模式
type SeasonalPattern struct {
	CurrentSeason        Season
	SeasonStartTime      time.Time
	SeasonDuration       time.Duration
	WeatherProbabilities map[Season]map[WeatherType]float64
	TemperatureRanges    map[Season]TemperatureRange
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// NewSeasonalPattern 创建季节模式
func NewSeasonalPattern() *SeasonalPattern {
	now := time.Now()
	pattern := &SeasonalPattern{
		CurrentSeason:        getCurrentSeason(now),
		SeasonStartTime:      getSeasonStartTime(now),
		SeasonDuration:       90 * 24 * time.Hour, // 90天
		WeatherProbabilities: make(map[Season]map[WeatherType]float64),
		TemperatureRanges:    make(map[Season]TemperatureRange),
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	// 初始化默认概率和温度范围
	pattern.initializeDefaultProbabilities()
	pattern.initializeTemperatureRanges()

	return pattern
}

// GetCurrentSeason 获取当前季节
func (sp *SeasonalPattern) GetCurrentSeason(currentTime time.Time) Season {
	return getCurrentSeason(currentTime)
}

// GetWeatherProbabilities 获取天气概率
func (sp *SeasonalPattern) GetWeatherProbabilities(season Season) map[WeatherType]float64 {
	return sp.WeatherProbabilities[season]
}

// GetTemperatureRange 获取温度范围
func (sp *SeasonalPattern) GetTemperatureRange(season Season) TemperatureRange {
	return sp.TemperatureRanges[season]
}

// UpdateWeatherProbability 更新天气概率
func (sp *SeasonalPattern) UpdateWeatherProbability(season Season, weatherType WeatherType, probability float64) {
	if sp.WeatherProbabilities[season] == nil {
		sp.WeatherProbabilities[season] = make(map[WeatherType]float64)
	}
	sp.WeatherProbabilities[season][weatherType] = probability
	sp.UpdatedAt = time.Now()
}

// UpdateTemperatureRange 更新温度范围
func (sp *SeasonalPattern) UpdateTemperatureRange(season Season, tempRange TemperatureRange) {
	sp.TemperatureRanges[season] = tempRange
	sp.UpdatedAt = time.Now()
}

// initializeDefaultProbabilities 初始化默认概率
func (sp *SeasonalPattern) initializeDefaultProbabilities() {
	// 春季
	sp.WeatherProbabilities[SeasonSpring] = map[WeatherType]float64{
		WeatherTypeSunny:  0.4,
		WeatherTypeCloudy: 0.3,
		WeatherTypeRainy:  0.2,
		WeatherTypeWindy:  0.1,
	}

	// 夏季
	sp.WeatherProbabilities[SeasonSummer] = map[WeatherType]float64{
		WeatherTypeSunny:  0.6,
		WeatherTypeCloudy: 0.2,
		WeatherTypeRainy:  0.1,
		WeatherTypeStormy: 0.1,
	}

	// 秋季
	sp.WeatherProbabilities[SeasonAutumn] = map[WeatherType]float64{
		WeatherTypeSunny:  0.3,
		WeatherTypeCloudy: 0.4,
		WeatherTypeRainy:  0.2,
		WeatherTypeWindy:  0.1,
	}

	// 冬季
	sp.WeatherProbabilities[SeasonWinter] = map[WeatherType]float64{
		WeatherTypeCloudy: 0.4,
		WeatherTypeSnowy:  0.3,
		WeatherTypeSunny:  0.2,
		WeatherTypeFoggy:  0.1,
	}
}

// initializeTemperatureRanges 初始化温度范围
func (sp *SeasonalPattern) initializeTemperatureRanges() {
	sp.TemperatureRanges[SeasonSpring] = TemperatureRange{Min: 10, Max: 25, Average: 17.5}
	sp.TemperatureRanges[SeasonSummer] = TemperatureRange{Min: 20, Max: 35, Average: 27.5}
	sp.TemperatureRanges[SeasonAutumn] = TemperatureRange{Min: 5, Max: 20, Average: 12.5}
	sp.TemperatureRanges[SeasonWinter] = TemperatureRange{Min: -10, Max: 10, Average: 0}
}

// TemperatureRange 温度范围
type TemperatureRange struct {
	Min     float64
	Max     float64
	Average float64
}

// IsInRange 检查温度是否在范围内
func (tr TemperatureRange) IsInRange(temperature float64) bool {
	return temperature >= tr.Min && temperature <= tr.Max
}

// GetRandomTemperature 获取随机温度
func (tr TemperatureRange) GetRandomTemperature() float64 {
	// 简化的随机温度生成
	return tr.Min + (tr.Max-tr.Min)*0.5 // 返回中间值，实际可以使用随机数
}

// WeatherForecast 天气预报
type WeatherForecast struct {
	Time        time.Time
	WeatherType WeatherType
	Intensity   WeatherIntensity
	Temperature float64
	Humidity    float64
	WindSpeed   float64
	Visibility  float64
	Confidence  float64 // 预报置信度 0-1
	Description string
	CreatedAt   time.Time
}

// NewWeatherForecast 创建天气预报
func NewWeatherForecast(time time.Time, weatherType WeatherType, intensity WeatherIntensity) *WeatherForecast {
	return &WeatherForecast{
		Time:        time,
		WeatherType: weatherType,
		Intensity:   intensity,
		Temperature: weatherType.GetBaseTemperature(),
		Humidity:    weatherType.GetBaseHumidity(),
		WindSpeed:   weatherType.GetBaseWindSpeed(),
		Visibility:  weatherType.GetBaseVisibility(),
		Confidence:  0.8, // 默认置信度
		Description: fmt.Sprintf("%s %s", intensity.GetDescription(), weatherType.GetDescription()),
		CreatedAt:   time,
	}
}

// GetTime 获取时间
func (wf *WeatherForecast) GetTime() time.Time {
	return wf.Time
}

// GetWeatherType 获取天气类型
func (wf *WeatherForecast) GetWeatherType() WeatherType {
	return wf.WeatherType
}

// GetIntensity 获取强度
func (wf *WeatherForecast) GetIntensity() WeatherIntensity {
	return wf.Intensity
}

// GetConfidence 获取置信度
func (wf *WeatherForecast) GetConfidence() float64 {
	return wf.Confidence
}

// GetDescription 获取描述
func (wf *WeatherForecast) GetDescription() string {
	return wf.Description
}

// IsHighConfidence 是否高置信度
func (wf *WeatherForecast) IsHighConfidence() bool {
	return wf.Confidence >= 0.8
}

// ToMap 转换为映射
func (wf *WeatherForecast) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"time":         wf.Time,
		"weather_type": wf.WeatherType.String(),
		"intensity":    wf.Intensity.String(),
		"temperature":  wf.Temperature,
		"humidity":     wf.Humidity,
		"wind_speed":   wf.WindSpeed,
		"visibility":   wf.Visibility,
		"confidence":   wf.Confidence,
		"description":  wf.Description,
		"created_at":   wf.CreatedAt,
	}
}

// WeatherEventType 天气事件类型
type WeatherEventType int

const (
	WeatherEventTypeStorm WeatherEventType = iota + 1
	WeatherEventTypeBlizzard
	WeatherEventTypeHeatWave
	WeatherEventTypeColdWave
	WeatherEventTypeDrought
	WeatherEventTypeFlood
	WeatherEventTypeHurricane
	WeatherEventTypeTornado
)

// String 返回事件类型字符串
func (wet WeatherEventType) String() string {
	switch wet {
	case WeatherEventTypeStorm:
		return "storm"
	case WeatherEventTypeBlizzard:
		return "blizzard"
	case WeatherEventTypeHeatWave:
		return "heat_wave"
	case WeatherEventTypeColdWave:
		return "cold_wave"
	case WeatherEventTypeDrought:
		return "drought"
	case WeatherEventTypeFlood:
		return "flood"
	case WeatherEventTypeHurricane:
		return "hurricane"
	case WeatherEventTypeTornado:
		return "tornado"
	default:
		return "unknown"
	}
}

// GetDescription 获取事件描述
func (wet WeatherEventType) GetDescription() string {
	switch wet {
	case WeatherEventTypeStorm:
		return "暴风雨"
	case WeatherEventTypeBlizzard:
		return "暴雪"
	case WeatherEventTypeHeatWave:
		return "热浪"
	case WeatherEventTypeColdWave:
		return "寒潮"
	case WeatherEventTypeDrought:
		return "干旱"
	case WeatherEventTypeFlood:
		return "洪水"
	case WeatherEventTypeHurricane:
		return "飓风"
	case WeatherEventTypeTornado:
		return "龙卷风"
	default:
		return "未知事件"
	}
}

// WeatherEventSeverity 天气事件严重程度
type WeatherEventSeverity int

const (
	WeatherEventSeverityMinor WeatherEventSeverity = iota + 1
	WeatherEventSeverityModerate
	WeatherEventSeverityMajor
	WeatherEventSeverityCritical
	WeatherEventSeverityCatastrophic
)

// String 返回严重程度字符串
func (wes WeatherEventSeverity) String() string {
	switch wes {
	case WeatherEventSeverityMinor:
		return "minor"
	case WeatherEventSeverityModerate:
		return "moderate"
	case WeatherEventSeverityMajor:
		return "major"
	case WeatherEventSeverityCritical:
		return "critical"
	case WeatherEventSeverityCatastrophic:
		return "catastrophic"
	default:
		return "unknown"
	}
}

// GetDescription 获取严重程度描述
func (wes WeatherEventSeverity) GetDescription() string {
	switch wes {
	case WeatherEventSeverityMinor:
		return "轻微"
	case WeatherEventSeverityModerate:
		return "中等"
	case WeatherEventSeverityMajor:
		return "严重"
	case WeatherEventSeverityCritical:
		return "危急"
	case WeatherEventSeverityCatastrophic:
		return "灾难性"
	default:
		return "未知程度"
	}
}

// GetMultiplier 获取严重程度倍率
func (wes WeatherEventSeverity) GetMultiplier() float64 {
	switch wes {
	case WeatherEventSeverityMinor:
		return 1.2
	case WeatherEventSeverityModerate:
		return 1.5
	case WeatherEventSeverityMajor:
		return 2.0
	case WeatherEventSeverityCritical:
		return 3.0
	case WeatherEventSeverityCatastrophic:
		return 5.0
	default:
		return 1.0
	}
}

// WeatherCondition 天气条件
type WeatherCondition struct {
	WeatherType WeatherType
	Intensity   WeatherIntensity
	Duration    time.Duration
	Effects     []string
}

// NewWeatherCondition 创建天气条件
func NewWeatherCondition(weatherType WeatherType, intensity WeatherIntensity, duration time.Duration) *WeatherCondition {
	return &WeatherCondition{
		WeatherType: weatherType,
		Intensity:   intensity,
		Duration:    duration,
		Effects:     make([]string, 0),
	}
}

// AddEffect 添加效果
func (wc *WeatherCondition) AddEffect(effect string) {
	wc.Effects = append(wc.Effects, effect)
}

// GetEffects 获取效果列表
func (wc *WeatherCondition) GetEffects() []string {
	return wc.Effects
}

// Matches 检查是否匹配条件
func (wc *WeatherCondition) Matches(weatherType WeatherType, intensity WeatherIntensity) bool {
	return wc.WeatherType == weatherType && wc.Intensity == intensity
}

// 辅助函数

// getCurrentSeason 获取当前季节
func getCurrentSeason(currentTime time.Time) Season {
	month := currentTime.Month()
	switch {
	case month >= 3 && month <= 5:
		return SeasonSpring
	case month >= 6 && month <= 8:
		return SeasonSummer
	case month >= 9 && month <= 11:
		return SeasonAutumn
	default:
		return SeasonWinter
	}
}

// getSeasonStartTime 获取季节开始时间
func getSeasonStartTime(currentTime time.Time) time.Time {
	year := currentTime.Year()
	month := currentTime.Month()

	switch {
	case month >= 3 && month <= 5:
		return time.Date(year, 3, 1, 0, 0, 0, 0, currentTime.Location())
	case month >= 6 && month <= 8:
		return time.Date(year, 6, 1, 0, 0, 0, 0, currentTime.Location())
	case month >= 9 && month <= 11:
		return time.Date(year, 9, 1, 0, 0, 0, 0, currentTime.Location())
	default:
		if month == 12 {
			return time.Date(year, 12, 1, 0, 0, 0, 0, currentTime.Location())
		} else {
			return time.Date(year-1, 12, 1, 0, 0, 0, 0, currentTime.Location())
		}
	}
}
