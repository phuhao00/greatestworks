package weather

import (
	"fmt"
	"time"
)

// WeatherState 天气状态实体
type WeatherState struct {
	ID          string
	WeatherType WeatherType
	Intensity   WeatherIntensity
	StartTime   time.Time
	EndTime     time.Time
	Duration    time.Duration
	Temperature float64 // 温度（摄氏度）
	Humidity    float64 // 湿度（百分比）
	WindSpeed   float64 // 风速（km/h）
	Visibility  float64 // 能见度（km）
	Pressure    float64 // 气压（hPa）
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewWeatherState 创建天气状态
func NewWeatherState(weatherType WeatherType, intensity WeatherIntensity) *WeatherState {
	now := time.Now()
	return &WeatherState{
		ID:          generateWeatherID(""),
		WeatherType: weatherType,
		Intensity:   intensity,
		StartTime:   now,
		Temperature: weatherType.GetBaseTemperature(),
		Humidity:    weatherType.GetBaseHumidity(),
		WindSpeed:   weatherType.GetBaseWindSpeed(),
		Visibility:  weatherType.GetBaseVisibility(),
		Pressure:    1013.25, // 标准大气压
		Description: fmt.Sprintf("%s %s", intensity.GetDescription(), weatherType.GetDescription()),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// GetID 获取ID
func (ws *WeatherState) GetID() string {
	return ws.ID
}

// GetWeatherType 获取天气类型
func (ws *WeatherState) GetWeatherType() WeatherType {
	return ws.WeatherType
}

// GetIntensity 获取强度
func (ws *WeatherState) GetIntensity() WeatherIntensity {
	return ws.Intensity
}

// GetStartTime 获取开始时间
func (ws *WeatherState) GetStartTime() time.Time {
	return ws.StartTime
}

// GetEndTime 获取结束时间
func (ws *WeatherState) GetEndTime() time.Time {
	return ws.EndTime
}

// GetDuration 获取持续时间
func (ws *WeatherState) GetDuration() time.Duration {
	return ws.Duration
}

// IsActive 检查天气是否激活
func (ws *WeatherState) IsActive() bool {
	now := time.Now()
	return now.After(ws.StartTime) && now.Before(ws.EndTime)
}

// IsExpired 检查天气是否过期
func (ws *WeatherState) IsExpired() bool {
	return time.Now().After(ws.EndTime)
}

// GetRemainingTime 获取剩余时间
func (ws *WeatherState) GetRemainingTime() time.Duration {
	now := time.Now()
	if now.After(ws.EndTime) {
		return 0
	}
	return ws.EndTime.Sub(now)
}

// GetElapsedTime 获取已过时间
func (ws *WeatherState) GetElapsedTime() time.Duration {
	now := time.Now()
	if now.Before(ws.StartTime) {
		return 0
	}
	if now.After(ws.EndTime) {
		return ws.Duration
	}
	return now.Sub(ws.StartTime)
}

// GetProgress 获取进度（0-1）
func (ws *WeatherState) GetProgress() float64 {
	if ws.Duration == 0 {
		return 1.0
	}

	elapsed := ws.GetElapsedTime()
	progress := float64(elapsed) / float64(ws.Duration)

	if progress > 1.0 {
		return 1.0
	}
	if progress < 0.0 {
		return 0.0
	}

	return progress
}

// UpdateTemperature 更新温度
func (ws *WeatherState) UpdateTemperature(temperature float64) {
	ws.Temperature = temperature
	ws.UpdatedAt = time.Now()
}

// UpdateHumidity 更新湿度
func (ws *WeatherState) UpdateHumidity(humidity float64) {
	ws.Humidity = humidity
	ws.UpdatedAt = time.Now()
}

// UpdateWindSpeed 更新风速
func (ws *WeatherState) UpdateWindSpeed(windSpeed float64) {
	ws.WindSpeed = windSpeed
	ws.UpdatedAt = time.Now()
}

// UpdateVisibility 更新能见度
func (ws *WeatherState) UpdateVisibility(visibility float64) {
	ws.Visibility = visibility
	ws.UpdatedAt = time.Now()
}

// UpdatePressure 更新气压
func (ws *WeatherState) UpdatePressure(pressure float64) {
	ws.Pressure = pressure
	ws.UpdatedAt = time.Now()
}

// GetTemperature 获取温度
func (ws *WeatherState) GetTemperature() float64 {
	return ws.Temperature
}

// GetHumidity 获取湿度
func (ws *WeatherState) GetHumidity() float64 {
	return ws.Humidity
}

// GetWindSpeed 获取风速
func (ws *WeatherState) GetWindSpeed() float64 {
	return ws.WindSpeed
}

// GetVisibility 获取能见度
func (ws *WeatherState) GetVisibility() float64 {
	return ws.Visibility
}

// GetPressure 获取气压
func (ws *WeatherState) GetPressure() float64 {
	return ws.Pressure
}

// GetDescription 获取描述
func (ws *WeatherState) GetDescription() string {
	return ws.Description
}

// UpdateDescription 更新描述
func (ws *WeatherState) UpdateDescription(description string) {
	ws.Description = description
	ws.UpdatedAt = time.Now()
}

// GetCreatedAt 获取创建时间
func (ws *WeatherState) GetCreatedAt() time.Time {
	return ws.CreatedAt
}

// GetUpdatedAt 获取更新时间
func (ws *WeatherState) GetUpdatedAt() time.Time {
	return ws.UpdatedAt
}

// ToMap 转换为映射
func (ws *WeatherState) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":           ws.ID,
		"weather_type": ws.WeatherType.String(),
		"intensity":    ws.Intensity.String(),
		"start_time":   ws.StartTime,
		"end_time":     ws.EndTime,
		"duration":     ws.Duration,
		"temperature":  ws.Temperature,
		"humidity":     ws.Humidity,
		"wind_speed":   ws.WindSpeed,
		"visibility":   ws.Visibility,
		"pressure":     ws.Pressure,
		"description":  ws.Description,
		"is_active":    ws.IsActive(),
		"progress":     ws.GetProgress(),
		"created_at":   ws.CreatedAt,
		"updated_at":   ws.UpdatedAt,
	}
}

// WeatherEffect 天气效果实体
type WeatherEffect struct {
	ID         string
	EffectType string
	TargetType string
	Multiplier float64
	Duration   time.Duration
	StartTime  time.Time
	EndTime    time.Time
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewWeatherEffect 创建天气效果
func NewWeatherEffect(effectType, targetType string, multiplier float64, duration time.Duration) *WeatherEffect {
	now := time.Now()
	return &WeatherEffect{
		ID:         generateEffectID(),
		EffectType: effectType,
		TargetType: targetType,
		Multiplier: multiplier,
		Duration:   duration,
		StartTime:  now,
		EndTime:    now.Add(duration),
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// GetID 获取ID
func (we *WeatherEffect) GetID() string {
	return we.ID
}

// GetEffectType 获取效果类型
func (we *WeatherEffect) GetEffectType() string {
	return we.EffectType
}

// GetTargetType 获取目标类型
func (we *WeatherEffect) GetTargetType() string {
	return we.TargetType
}

// GetMultiplier 获取倍率
func (we *WeatherEffect) GetMultiplier() float64 {
	return we.Multiplier
}

// GetModifier 获取修饰符（别名方法，与GetMultiplier相同）
func (we *WeatherEffect) GetModifier() float64 {
	return we.Multiplier
}

// GetDuration 获取持续时间
func (we *WeatherEffect) GetDuration() time.Duration {
	return we.Duration
}

// GetStartTime 获取开始时间
func (we *WeatherEffect) GetStartTime() time.Time {
	return we.StartTime
}

// GetEndTime 获取结束时间
func (we *WeatherEffect) GetEndTime() time.Time {
	return we.EndTime
}

// IsEffectActive 检查效果是否激活
func (we *WeatherEffect) IsEffectActive() bool {
	now := time.Now()
	return we.IsActive && now.After(we.StartTime) && now.Before(we.EndTime)
}

// IsEffectExpired 检查效果是否过期
func (we *WeatherEffect) IsEffectExpired() bool {
	return time.Now().After(we.EndTime)
}

// GetRemainingTime 获取剩余时间
func (we *WeatherEffect) GetRemainingTime() time.Duration {
	now := time.Now()
	if now.After(we.EndTime) {
		return 0
	}
	return we.EndTime.Sub(now)
}

// Activate 激活效果
func (we *WeatherEffect) Activate() {
	we.IsActive = true
	we.UpdatedAt = time.Now()
}

// Deactivate 停用效果
func (we *WeatherEffect) Deactivate() {
	we.IsActive = false
	we.UpdatedAt = time.Now()
}

// ExtendDuration 延长持续时间
func (we *WeatherEffect) ExtendDuration(additionalDuration time.Duration) {
	we.Duration += additionalDuration
	we.EndTime = we.EndTime.Add(additionalDuration)
	we.UpdatedAt = time.Now()
}

// UpdateMultiplier 更新倍率
func (we *WeatherEffect) UpdateMultiplier(multiplier float64) {
	we.Multiplier = multiplier
	we.UpdatedAt = time.Now()
}

// GetCreatedAt 获取创建时间
func (we *WeatherEffect) GetCreatedAt() time.Time {
	return we.CreatedAt
}

// GetUpdatedAt 获取更新时间
func (we *WeatherEffect) GetUpdatedAt() time.Time {
	return we.UpdatedAt
}

// ToMap 转换为映射
func (we *WeatherEffect) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":               we.ID,
		"effect_type":      we.EffectType,
		"target_type":      we.TargetType,
		"multiplier":       we.Multiplier,
		"duration":         we.Duration,
		"start_time":       we.StartTime,
		"end_time":         we.EndTime,
		"is_active":        we.IsActive,
		"is_effect_active": we.IsEffectActive(),
		"remaining_time":   we.GetRemainingTime(),
		"created_at":       we.CreatedAt,
		"updated_at":       we.UpdatedAt,
	}
}

// WeatherEvent 天气事件实体
type WeatherEvent struct {
	ID          string
	EventType   WeatherEventType
	Severity    WeatherEventSeverity
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	Duration    time.Duration
	Effects     []*WeatherEffect
	Triggers    []string // 触发条件
	Rewards     []string // 奖励
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewWeatherEvent 创建天气事件
func NewWeatherEvent(eventType WeatherEventType, severity WeatherEventSeverity, title, description string, duration time.Duration) *WeatherEvent {
	now := time.Now()
	return &WeatherEvent{
		ID:          generateEventID(),
		EventType:   eventType,
		Severity:    severity,
		Title:       title,
		Description: description,
		StartTime:   now,
		EndTime:     now.Add(duration),
		Duration:    duration,
		Effects:     make([]*WeatherEffect, 0),
		Triggers:    make([]string, 0),
		Rewards:     make([]string, 0),
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// GetID 获取ID
func (we *WeatherEvent) GetID() string {
	return we.ID
}

// GetEventType 获取事件类型
func (we *WeatherEvent) GetEventType() WeatherEventType {
	return we.EventType
}

// GetSeverity 获取严重程度
func (we *WeatherEvent) GetSeverity() WeatherEventSeverity {
	return we.Severity
}

// GetTitle 获取标题
func (we *WeatherEvent) GetTitle() string {
	return we.Title
}

// GetDescription 获取描述
func (we *WeatherEvent) GetDescription() string {
	return we.Description
}

// GetStartTime 获取开始时间
func (we *WeatherEvent) GetStartTime() time.Time {
	return we.StartTime
}

// GetEndTime 获取结束时间
func (we *WeatherEvent) GetEndTime() time.Time {
	return we.EndTime
}

// GetDuration 获取持续时间
func (we *WeatherEvent) GetDuration() time.Duration {
	return we.Duration
}

// GetEffects 获取效果列表
func (we *WeatherEvent) GetEffects() []*WeatherEffect {
	return we.Effects
}

// AddEffect 添加效果
func (we *WeatherEvent) AddEffect(effect *WeatherEffect) {
	we.Effects = append(we.Effects, effect)
	we.UpdatedAt = time.Now()
}

// RemoveEffect 移除效果
func (we *WeatherEvent) RemoveEffect(effectID string) {
	for i, effect := range we.Effects {
		if effect.GetID() == effectID {
			we.Effects = append(we.Effects[:i], we.Effects[i+1:]...)
			we.UpdatedAt = time.Now()
			break
		}
	}
}

// GetTriggers 获取触发条件
func (we *WeatherEvent) GetTriggers() []string {
	return we.Triggers
}

// AddTrigger 添加触发条件
func (we *WeatherEvent) AddTrigger(trigger string) {
	we.Triggers = append(we.Triggers, trigger)
	we.UpdatedAt = time.Now()
}

// GetRewards 获取奖励
func (we *WeatherEvent) GetRewards() []string {
	return we.Rewards
}

// AddReward 添加奖励
func (we *WeatherEvent) AddReward(reward string) {
	we.Rewards = append(we.Rewards, reward)
	we.UpdatedAt = time.Now()
}

// IsEventActive 检查事件是否激活
func (we *WeatherEvent) IsEventActive() bool {
	now := time.Now()
	return we.IsActive && now.After(we.StartTime) && now.Before(we.EndTime)
}

// IsEventExpired 检查事件是否过期
func (we *WeatherEvent) IsEventExpired() bool {
	return time.Now().After(we.EndTime)
}

// Activate 激活事件
func (we *WeatherEvent) Activate() {
	we.IsActive = true
	we.UpdatedAt = time.Now()
}

// Deactivate 停用事件
func (we *WeatherEvent) Deactivate() {
	we.IsActive = false
	we.UpdatedAt = time.Now()
}

// GetCreatedAt 获取创建时间
func (we *WeatherEvent) GetCreatedAt() time.Time {
	return we.CreatedAt
}

// GetUpdatedAt 获取更新时间
func (we *WeatherEvent) GetUpdatedAt() time.Time {
	return we.UpdatedAt
}

// ToMap 转换为映射
func (we *WeatherEvent) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":              we.ID,
		"event_type":      we.EventType.String(),
		"severity":        we.Severity.String(),
		"title":           we.Title,
		"description":     we.Description,
		"start_time":      we.StartTime,
		"end_time":        we.EndTime,
		"duration":        we.Duration,
		"effects":         len(we.Effects),
		"triggers":        we.Triggers,
		"rewards":         we.Rewards,
		"is_active":       we.IsActive,
		"is_event_active": we.IsEventActive(),
		"created_at":      we.CreatedAt,
		"updated_at":      we.UpdatedAt,
	}
}

// 辅助函数

// generateWeatherID 生成天气ID - 已在aggregate.go中定义

// generateEffectID 生成效果ID
func generateEffectID() string {
	return fmt.Sprintf("effect_%d", time.Now().UnixNano())
}

// generateEventID 生成事件ID
func generateEventID() string {
	return fmt.Sprintf("event_%d", time.Now().UnixNano())
}
