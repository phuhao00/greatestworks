package weather

import (
	"time"
	"math/rand"
	"errors"
)

// WeatherAggregate 天气聚合根
type WeatherAggregate struct {
	sceneID         string
	currentWeather  *WeatherState
	weatherHistory  []*WeatherState
	weatherForecast []*WeatherForecast
	weatherEffects  map[string]*WeatherEffect
	seasonalPattern *SeasonalPattern
	lastUpdateTime  time.Time
	nextChangeTime  time.Time
	changeInterval  time.Duration
	randomSeed      int64
	updatedAt       time.Time
	version         int
}

// NewWeatherAggregate 创建天气聚合根
func NewWeatherAggregate(sceneID string) *WeatherAggregate {
	now := time.Now()
	return &WeatherAggregate{
		sceneID:         sceneID,
		currentWeather:  NewWeatherState(WeatherTypeSunny, WeatherIntensityNormal),
		weatherHistory:  make([]*WeatherState, 0),
		weatherForecast: make([]*WeatherForecast, 0),
		weatherEffects:  make(map[string]*WeatherEffect),
		seasonalPattern: NewSeasonalPattern(),
		lastUpdateTime:  now,
		nextChangeTime:  now.Add(30 * time.Minute), // 默认30分钟变化一次
		changeInterval:  30 * time.Minute,
		randomSeed:      now.UnixNano(),
		updatedAt:       now,
		version:         1,
	}
}

// GetSceneID 获取场景ID
func (w *WeatherAggregate) GetSceneID() string {
	return w.sceneID
}

// GetCurrentWeather 获取当前天气
func (w *WeatherAggregate) GetCurrentWeather() *WeatherState {
	return w.currentWeather
}

// ChangeWeather 改变天气
func (w *WeatherAggregate) ChangeWeather(weatherType WeatherType, intensity WeatherIntensity) error {
	if !weatherType.IsValid() {
		return ErrInvalidWeatherType
	}
	
	if !intensity.IsValid() {
		return ErrInvalidWeatherIntensity
	}
	
	// 保存当前天气到历史记录
	if w.currentWeather != nil {
		w.addToHistory(w.currentWeather)
	}
	
	// 创建新的天气状态
	newWeather := NewWeatherState(weatherType, intensity)
	newWeather.StartTime = time.Now()
	
	// 计算持续时间
	duration := w.calculateWeatherDuration(weatherType, intensity)
	newWeather.Duration = duration
	newWeather.EndTime = newWeather.StartTime.Add(duration)
	
	w.currentWeather = newWeather
	w.lastUpdateTime = time.Now()
	w.nextChangeTime = newWeather.EndTime
	
	// 更新天气效果
	w.updateWeatherEffects()
	
	w.updateVersion()
	return nil
}

// UpdateWeather 更新天气（自动变化）
func (w *WeatherAggregate) UpdateWeather() error {
	now := time.Now()
	
	// 检查是否需要变化天气
	if now.Before(w.nextChangeTime) {
		return nil // 还未到变化时间
	}
	
	// 根据季节模式和随机因素决定下一个天气
	nextWeather := w.calculateNextWeather()
	
	return w.ChangeWeather(nextWeather.WeatherType, nextWeather.Intensity)
}

// GetWeatherEffects 获取天气效果
func (w *WeatherAggregate) GetWeatherEffects() map[string]*WeatherEffect {
	return w.weatherEffects
}

// GetWeatherEffect 获取指定的天气效果
func (w *WeatherAggregate) GetWeatherEffect(effectType string) *WeatherEffect {
	return w.weatherEffects[effectType]
}

// AddWeatherEffect 添加天气效果
func (w *WeatherAggregate) AddWeatherEffect(effect *WeatherEffect) {
	w.weatherEffects[effect.GetEffectType()] = effect
	w.updateVersion()
}

// RemoveWeatherEffect 移除天气效果
func (w *WeatherAggregate) RemoveWeatherEffect(effectType string) {
	delete(w.weatherEffects, effectType)
	w.updateVersion()
}

// GetWeatherHistory 获取天气历史
func (w *WeatherAggregate) GetWeatherHistory() []*WeatherState {
	return w.weatherHistory
}

// GetWeatherForecast 获取天气预报
func (w *WeatherAggregate) GetWeatherForecast() []*WeatherForecast {
	return w.weatherForecast
}

// GenerateForecast 生成天气预报
func (w *WeatherAggregate) GenerateForecast(hours int) error {
	if hours <= 0 || hours > 168 { // 最多预报一周
		return ErrInvalidForecastPeriod
	}
	
	w.weatherForecast = make([]*WeatherForecast, 0)
	
	currentTime := time.Now()
	currentWeather := w.currentWeather
	
	for i := 1; i <= hours; i++ {
		forecastTime := currentTime.Add(time.Duration(i) * time.Hour)
		
		// 基于当前天气和季节模式预测
		predictedWeather := w.predictWeatherAt(forecastTime)
		
		forecast := &WeatherForecast{
			Time:        forecastTime,
			WeatherType: predictedWeather.WeatherType,
			Intensity:   predictedWeather.Intensity,
			Confidence:  w.calculateForecastConfidence(i),
			Description: w.generateWeatherDescription(predictedWeather.WeatherType, predictedWeather.Intensity),
		}
		
		w.weatherForecast = append(w.weatherForecast, forecast)
	}
	
	w.updateVersion()
	return nil
}

// GetSeasonalPattern 获取季节模式
func (w *WeatherAggregate) GetSeasonalPattern() *SeasonalPattern {
	return w.seasonalPattern
}

// UpdateSeasonalPattern 更新季节模式
func (w *WeatherAggregate) UpdateSeasonalPattern(pattern *SeasonalPattern) {
	w.seasonalPattern = pattern
	w.updateVersion()
}

// GetLastUpdateTime 获取最后更新时间
func (w *WeatherAggregate) GetLastUpdateTime() time.Time {
	return w.lastUpdateTime
}

// GetNextChangeTime 获取下次变化时间
func (w *WeatherAggregate) GetNextChangeTime() time.Time {
	return w.nextChangeTime
}

// SetChangeInterval 设置变化间隔
func (w *WeatherAggregate) SetChangeInterval(interval time.Duration) error {
	if interval < time.Minute || interval > 24*time.Hour {
		return ErrInvalidChangeInterval
	}
	
	w.changeInterval = interval
	w.nextChangeTime = w.lastUpdateTime.Add(interval)
	w.updateVersion()
	return nil
}

// GetChangeInterval 获取变化间隔
func (w *WeatherAggregate) GetChangeInterval() time.Duration {
	return w.changeInterval
}

// IsWeatherActive 检查天气是否激活
func (w *WeatherAggregate) IsWeatherActive() bool {
	return w.currentWeather != nil && !w.currentWeather.IsExpired()
}

// GetRemainingDuration 获取当前天气剩余时间
func (w *WeatherAggregate) GetRemainingDuration() time.Duration {
	if w.currentWeather == nil {
		return 0
	}
	
	now := time.Now()
	if now.After(w.currentWeather.EndTime) {
		return 0
	}
	
	return w.currentWeather.EndTime.Sub(now)
}

// CalculateWeatherInfluence 计算天气对指定属性的影响
func (w *WeatherAggregate) CalculateWeatherInfluence(attributeType string) float64 {
	if w.currentWeather == nil {
		return 1.0 // 无影响
	}
	
	influence := 1.0
	
	// 基于天气类型的影响
	switch w.currentWeather.WeatherType {
	case WeatherTypeSunny:
		if attributeType == "visibility" {
			influence = 1.2
		} else if attributeType == "movement_speed" {
			influence = 1.1
		}
	case WeatherTypeRainy:
		if attributeType == "visibility" {
			influence = 0.8
		} else if attributeType == "fire_damage" {
			influence = 0.7
		}
	case WeatherTypeSnowy:
		if attributeType == "movement_speed" {
			influence = 0.8
		} else if attributeType == "ice_damage" {
			influence = 1.3
		}
	case WeatherTypeStormy:
		if attributeType == "lightning_damage" {
			influence = 1.5
		} else if attributeType == "accuracy" {
			influence = 0.9
		}
	case WeatherTypeFoggy:
		if attributeType == "visibility" {
			influence = 0.5
		} else if attributeType == "detection_range" {
			influence = 0.6
		}
	}
	
	// 基于强度调整影响
	intensityMultiplier := w.currentWeather.Intensity.GetMultiplier()
	if influence != 1.0 {
		// 只对有影响的属性应用强度倍率
		if influence > 1.0 {
			influence = 1.0 + (influence-1.0)*intensityMultiplier
		} else {
			influence = 1.0 - (1.0-influence)*intensityMultiplier
		}
	}
	
	return influence
}

// GetVersion 获取版本
func (w *WeatherAggregate) GetVersion() int {
	return w.version
}

// GetUpdatedAt 获取更新时间
func (w *WeatherAggregate) GetUpdatedAt() time.Time {
	return w.updatedAt
}

// 私有方法

// addToHistory 添加到历史记录
func (w *WeatherAggregate) addToHistory(weather *WeatherState) {
	w.weatherHistory = append(w.weatherHistory, weather)
	
	// 限制历史记录数量（保留最近100条）
	if len(w.weatherHistory) > 100 {
		w.weatherHistory = w.weatherHistory[1:]
	}
}

// calculateWeatherDuration 计算天气持续时间
func (w *WeatherAggregate) calculateWeatherDuration(weatherType WeatherType, intensity WeatherIntensity) time.Duration {
	baseDuration := w.changeInterval
	
	// 根据天气类型调整持续时间
	switch weatherType {
	case WeatherTypeSunny:
		baseDuration = baseDuration * 2 // 晴天持续更久
	case WeatherTypeStormy:
		baseDuration = baseDuration / 2 // 暴风雨持续较短
	case WeatherTypeFoggy:
		baseDuration = baseDuration / 3 // 雾天持续很短
	}
	
	// 根据强度调整
	intensityFactor := intensity.GetDurationFactor()
	baseDuration = time.Duration(float64(baseDuration) * intensityFactor)
	
	// 添加随机因素（±20%）
	randomFactor := 0.8 + rand.Float64()*0.4
	baseDuration = time.Duration(float64(baseDuration) * randomFactor)
	
	return baseDuration
}

// calculateNextWeather 计算下一个天气
func (w *WeatherAggregate) calculateNextWeather() *WeatherState {
	now := time.Now()
	currentSeason := w.seasonalPattern.GetCurrentSeason(now)
	
	// 获取季节天气概率
	weatherProbabilities := w.seasonalPattern.GetWeatherProbabilities(currentSeason)
	
	// 考虑当前天气的影响（天气转换规律）
	currentType := w.currentWeather.WeatherType
	transitionProbabilities := w.getWeatherTransitionProbabilities(currentType)
	
	// 合并概率
	combinedProbabilities := w.combineWeatherProbabilities(weatherProbabilities, transitionProbabilities)
	
	// 随机选择下一个天气
	nextWeatherType := w.selectWeatherByProbability(combinedProbabilities)
	nextIntensity := w.selectRandomIntensity(nextWeatherType)
	
	return NewWeatherState(nextWeatherType, nextIntensity)
}

// updateWeatherEffects 更新天气效果
func (w *WeatherAggregate) updateWeatherEffects() {
	// 清除旧的效果
	w.weatherEffects = make(map[string]*WeatherEffect)
	
	if w.currentWeather == nil {
		return
	}
	
	// 根据当前天气添加效果
	effects := w.generateWeatherEffects(w.currentWeather.WeatherType, w.currentWeather.Intensity)
	for _, effect := range effects {
		w.weatherEffects[effect.GetEffectType()] = effect
	}
}

// generateWeatherEffects 生成天气效果
func (w *WeatherAggregate) generateWeatherEffects(weatherType WeatherType, intensity WeatherIntensity) []*WeatherEffect {
	effects := make([]*WeatherEffect, 0)
	
	switch weatherType {
	case WeatherTypeSunny:
		effects = append(effects, NewWeatherEffect("visibility_boost", 1.2*intensity.GetMultiplier(), w.currentWeather.Duration))
		effects = append(effects, NewWeatherEffect("movement_speed_boost", 1.1*intensity.GetMultiplier(), w.currentWeather.Duration))
	
	case WeatherTypeRainy:
		effects = append(effects, NewWeatherEffect("visibility_reduction", 0.8/intensity.GetMultiplier(), w.currentWeather.Duration))
		effects = append(effects, NewWeatherEffect("fire_damage_reduction", 0.7/intensity.GetMultiplier(), w.currentWeather.Duration))
		effects = append(effects, NewWeatherEffect("water_damage_boost", 1.2*intensity.GetMultiplier(), w.currentWeather.Duration))
	
	case WeatherTypeSnowy:
		effects = append(effects, NewWeatherEffect("movement_speed_reduction", 0.8/intensity.GetMultiplier(), w.currentWeather.Duration))
		effects = append(effects, NewWeatherEffect("ice_damage_boost", 1.3*intensity.GetMultiplier(), w.currentWeather.Duration))
		effects = append(effects, NewWeatherEffect("cold_resistance_reduction", 0.9/intensity.GetMultiplier(), w.currentWeather.Duration))
	
	case WeatherTypeStormy:
		effects = append(effects, NewWeatherEffect("lightning_damage_boost", 1.5*intensity.GetMultiplier(), w.currentWeather.Duration))
		effects = append(effects, NewWeatherEffect("accuracy_reduction", 0.9/intensity.GetMultiplier(), w.currentWeather.Duration))
		effects = append(effects, NewWeatherEffect("wind_resistance_reduction", 0.8/intensity.GetMultiplier(), w.currentWeather.Duration))
	
	case WeatherTypeFoggy:
		effects = append(effects, NewWeatherEffect("visibility_severe_reduction", 0.5/intensity.GetMultiplier(), w.currentWeather.Duration))
		effects = append(effects, NewWeatherEffect("detection_range_reduction", 0.6/intensity.GetMultiplier(), w.currentWeather.Duration))
		effects = append(effects, NewWeatherEffect("stealth_boost", 1.3*intensity.GetMultiplier(), w.currentWeather.Duration))
	}
	
	return effects
}

// predictWeatherAt 预测指定时间的天气
func (w *WeatherAggregate) predictWeatherAt(targetTime time.Time) *WeatherState {
	// 简化的预测算法，实际可以更复杂
	hoursDiff := int(targetTime.Sub(time.Now()).Hours())
	
	// 基于小时数和季节模式预测
	currentSeason := w.seasonalPattern.GetCurrentSeason(targetTime)
	weatherProbabilities := w.seasonalPattern.GetWeatherProbabilities(currentSeason)
	
	// 添加时间因素的随机性
	rand.Seed(w.randomSeed + int64(hoursDiff))
	weatherType := w.selectWeatherByProbability(weatherProbabilities)
	intensity := w.selectRandomIntensity(weatherType)
	
	return NewWeatherState(weatherType, intensity)
}

// calculateForecastConfidence 计算预报置信度
func (w *WeatherAggregate) calculateForecastConfidence(hoursAhead int) float64 {
	// 预报时间越远，置信度越低
	baseConfidence := 0.95
	decayRate := 0.02
	
	confidence := baseConfidence - float64(hoursAhead)*decayRate
	if confidence < 0.3 {
		confidence = 0.3 // 最低置信度
	}
	
	return confidence
}

// generateWeatherDescription 生成天气描述
func (w *WeatherAggregate) generateWeatherDescription(weatherType WeatherType, intensity WeatherIntensity) string {
	baseDescription := weatherType.GetDescription()
	intensityDescription := intensity.GetDescription()
	
	return intensityDescription + baseDescription
}

// getWeatherTransitionProbabilities 获取天气转换概率
func (w *WeatherAggregate) getWeatherTransitionProbabilities(currentType WeatherType) map[WeatherType]float64 {
	// 定义天气转换规律
	transitions := make(map[WeatherType]float64)
	
	switch currentType {
	case WeatherTypeSunny:
		transitions[WeatherTypeSunny] = 0.6
		transitions[WeatherTypeCloudy] = 0.25
		transitions[WeatherTypeRainy] = 0.1
		transitions[WeatherTypeWindy] = 0.05
	
	case WeatherTypeCloudy:
		transitions[WeatherTypeCloudy] = 0.4
		transitions[WeatherTypeSunny] = 0.3
		transitions[WeatherTypeRainy] = 0.2
		transitions[WeatherTypeStormy] = 0.1
	
	case WeatherTypeRainy:
		transitions[WeatherTypeRainy] = 0.5
		transitions[WeatherTypeCloudy] = 0.3
		transitions[WeatherTypeStormy] = 0.15
		transitions[WeatherTypeSunny] = 0.05
	
	case WeatherTypeStormy:
		transitions[WeatherTypeStormy] = 0.3
		transitions[WeatherTypeRainy] = 0.4
		transitions[WeatherTypeCloudy] = 0.2
		transitions[WeatherTypeWindy] = 0.1
	
	default:
		// 默认均匀分布
		transitions[WeatherTypeSunny] = 0.3
		transitions[WeatherTypeCloudy] = 0.25
		transitions[WeatherTypeRainy] = 0.2
		transitions[WeatherTypeWindy] = 0.15
		transitions[WeatherTypeStormy] = 0.1
	}
	
	return transitions
}

// combineWeatherProbabilities 合并天气概率
func (w *WeatherAggregate) combineWeatherProbabilities(seasonal, transition map[WeatherType]float64) map[WeatherType]float64 {
	combined := make(map[WeatherType]float64)
	
	// 季节概率权重0.7，转换概率权重0.3
	seasonalWeight := 0.7
	transitionWeight := 0.3
	
	allWeatherTypes := []WeatherType{WeatherTypeSunny, WeatherTypeCloudy, WeatherTypeRainy, WeatherTypeWindy, WeatherTypeStormy, WeatherTypeSnowy, WeatherTypeFoggy}
	
	for _, weatherType := range allWeatherTypes {
		seasonalProb := seasonal[weatherType]
		transitionProb := transition[weatherType]
		
		combined[weatherType] = seasonalProb*seasonalWeight + transitionProb*transitionWeight
	}
	
	return combined
}

// selectWeatherByProbability 根据概率选择天气
func (w *WeatherAggregate) selectWeatherByProbability(probabilities map[WeatherType]float64) WeatherType {
	rand.Seed(time.Now().UnixNano() + w.randomSeed)
	randomValue := rand.Float64()
	
	cumulativeProbability := 0.0
	for weatherType, probability := range probabilities {
		cumulativeProbability += probability
		if randomValue <= cumulativeProbability {
			return weatherType
		}
	}
	
	// 默认返回晴天
	return WeatherTypeSunny
}

// selectRandomIntensity 随机选择强度
func (w *WeatherAggregate) selectRandomIntensity(weatherType WeatherType) WeatherIntensity {
	rand.Seed(time.Now().UnixNano() + w.randomSeed)
	
	// 根据天气类型调整强度概率
	switch weatherType {
	case WeatherTypeSunny, WeatherTypeCloudy:
		// 温和天气更可能是正常强度
		if rand.Float64() < 0.7 {
			return WeatherIntensityNormal
		} else if rand.Float64() < 0.9 {
			return WeatherIntensityLight
		} else {
			return WeatherIntensityHeavy
		}
	
	case WeatherTypeStormy:
		// 暴风雨更可能是强烈的
		if rand.Float64() < 0.5 {
			return WeatherIntensityHeavy
		} else if rand.Float64() < 0.8 {
			return WeatherIntensityNormal
		} else {
			return WeatherIntensityExtreme
		}
	
	default:
		// 其他天气均匀分布
		intensities := []WeatherIntensity{WeatherIntensityLight, WeatherIntensityNormal, WeatherIntensityHeavy}
		return intensities[rand.Intn(len(intensities))]
	}
}

// updateVersion 更新版本
func (w *WeatherAggregate) updateVersion() {
	w.version++
	w.updatedAt = time.Now()
}

// 天气相关错误
var (
	ErrInvalidWeatherType      = errors.New("invalid weather type")
	ErrInvalidWeatherIntensity = errors.New("invalid weather intensity")
	ErrInvalidForecastPeriod   = errors.New("invalid forecast period")
	ErrInvalidChangeInterval   = errors.New("invalid change interval")
)