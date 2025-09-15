package weather

import (
	"fmt"
	"time"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	GetEventID() string
	GetEventType() string
	GetAggregateID() string
	GetOccurredAt() time.Time
	GetVersion() int
	GetPayload() map[string]interface{}
}

// BaseDomainEvent 基础领域事件
type BaseDomainEvent struct {
	EventID     string
	EventType   string
	AggregateID string
	OccurredAt  time.Time
	Version     int
	Payload     map[string]interface{}
}

// GetEventID 获取事件ID
func (e *BaseDomainEvent) GetEventID() string {
	return e.EventID
}

// GetEventType 获取事件类型
func (e *BaseDomainEvent) GetEventType() string {
	return e.EventType
}

// GetAggregateID 获取聚合ID
func (e *BaseDomainEvent) GetAggregateID() string {
	return e.AggregateID
}

// GetOccurredAt 获取发生时间
func (e *BaseDomainEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetVersion 获取版本
func (e *BaseDomainEvent) GetVersion() int {
	return e.Version
}

// GetPayload 获取载荷
func (e *BaseDomainEvent) GetPayload() map[string]interface{} {
	return e.Payload
}

// WeatherChangedEvent 天气变化事件
type WeatherChangedEvent struct {
	*BaseDomainEvent
	SceneID         string
	PreviousWeather *WeatherState
	CurrentWeather  *WeatherState
	ChangeReason    string
	Effects         []*WeatherEffect
}

// NewWeatherChangedEvent 创建天气变化事件
func NewWeatherChangedEvent(sceneID string, previous, current *WeatherState, reason string, effects []*WeatherEffect) *WeatherChangedEvent {
	now := time.Now()
	event := &WeatherChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("weather_changed_%d", now.UnixNano()),
			EventType:   "weather.changed",
			AggregateID: sceneID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SceneID:         sceneID,
		PreviousWeather: previous,
		CurrentWeather:  current,
		ChangeReason:    reason,
		Effects:         effects,
	}
	
	// 设置载荷
	event.Payload["scene_id"] = sceneID
	event.Payload["change_reason"] = reason
	if previous != nil {
		event.Payload["previous_weather_type"] = previous.WeatherType.String()
		event.Payload["previous_intensity"] = previous.Intensity.String()
	}
	if current != nil {
		event.Payload["current_weather_type"] = current.WeatherType.String()
		event.Payload["current_intensity"] = current.Intensity.String()
		event.Payload["current_temperature"] = current.Temperature
		event.Payload["current_humidity"] = current.Humidity
	}
	event.Payload["effects_count"] = len(effects)
	
	return event
}

// WeatherIntensityChangedEvent 天气强度变化事件
type WeatherIntensityChangedEvent struct {
	*BaseDomainEvent
	SceneID           string
	WeatherType       WeatherType
	PreviousIntensity WeatherIntensity
	CurrentIntensity  WeatherIntensity
	ChangeReason      string
}

// NewWeatherIntensityChangedEvent 创建天气强度变化事件
func NewWeatherIntensityChangedEvent(sceneID string, weatherType WeatherType, previous, current WeatherIntensity, reason string) *WeatherIntensityChangedEvent {
	now := time.Now()
	event := &WeatherIntensityChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("weather_intensity_changed_%d", now.UnixNano()),
			EventType:   "weather.intensity_changed",
			AggregateID: sceneID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SceneID:           sceneID,
		WeatherType:       weatherType,
		PreviousIntensity: previous,
		CurrentIntensity:  current,
		ChangeReason:      reason,
	}
	
	// 设置载荷
	event.Payload["scene_id"] = sceneID
	event.Payload["weather_type"] = weatherType.String()
	event.Payload["previous_intensity"] = previous.String()
	event.Payload["current_intensity"] = current.String()
	event.Payload["change_reason"] = reason
	
	return event
}

// WeatherEffectActivatedEvent 天气效果激活事件
type WeatherEffectActivatedEvent struct {
	*BaseDomainEvent
	SceneID    string
	Effect     *WeatherEffect
	WeatherType WeatherType
	Intensity  WeatherIntensity
	Duration   time.Duration
}

// NewWeatherEffectActivatedEvent 创建天气效果激活事件
func NewWeatherEffectActivatedEvent(sceneID string, effect *WeatherEffect, weatherType WeatherType, intensity WeatherIntensity) *WeatherEffectActivatedEvent {
	now := time.Now()
	event := &WeatherEffectActivatedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("weather_effect_activated_%d", now.UnixNano()),
			EventType:   "weather.effect_activated",
			AggregateID: sceneID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SceneID:     sceneID,
		Effect:      effect,
		WeatherType: weatherType,
		Intensity:   intensity,
		Duration:    effect.Duration,
	}
	
	// 设置载荷
	event.Payload["scene_id"] = sceneID
	event.Payload["effect_id"] = effect.ID
	event.Payload["effect_type"] = effect.EffectType
	event.Payload["multiplier"] = effect.Multiplier
	event.Payload["duration"] = effect.Duration
	event.Payload["weather_type"] = weatherType.String()
	event.Payload["intensity"] = intensity.String()
	
	return event
}

// WeatherEffectDeactivatedEvent 天气效果停用事件
type WeatherEffectDeactivatedEvent struct {
	*BaseDomainEvent
	SceneID    string
	Effect     *WeatherEffect
	Reason     string
	Duration   time.Duration
}

// NewWeatherEffectDeactivatedEvent 创建天气效果停用事件
func NewWeatherEffectDeactivatedEvent(sceneID string, effect *WeatherEffect, reason string) *WeatherEffectDeactivatedEvent {
	now := time.Now()
	event := &WeatherEffectDeactivatedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("weather_effect_deactivated_%d", now.UnixNano()),
			EventType:   "weather.effect_deactivated",
			AggregateID: sceneID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SceneID:  sceneID,
		Effect:   effect,
		Reason:   reason,
		Duration: now.Sub(effect.StartTime),
	}
	
	// 设置载荷
	event.Payload["scene_id"] = sceneID
	event.Payload["effect_id"] = effect.ID
	event.Payload["effect_type"] = effect.EffectType
	event.Payload["reason"] = reason
	event.Payload["total_duration"] = now.Sub(effect.StartTime)
	
	return event
}

// WeatherEventTriggeredEvent 天气事件触发事件
type WeatherEventTriggeredEvent struct {
	*BaseDomainEvent
	SceneID      string
	WeatherEvent *WeatherEvent
	TriggerWeather *WeatherState
	Severity     WeatherEventSeverity
}

// NewWeatherEventTriggeredEvent 创建天气事件触发事件
func NewWeatherEventTriggeredEvent(sceneID string, weatherEvent *WeatherEvent, triggerWeather *WeatherState) *WeatherEventTriggeredEvent {
	now := time.Now()
	event := &WeatherEventTriggeredEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("weather_event_triggered_%d", now.UnixNano()),
			EventType:   "weather.event_triggered",
			AggregateID: sceneID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SceneID:        sceneID,
		WeatherEvent:   weatherEvent,
		TriggerWeather: triggerWeather,
		Severity:       weatherEvent.Severity,
	}
	
	// 设置载荷
	event.Payload["scene_id"] = sceneID
	event.Payload["event_id"] = weatherEvent.ID
	event.Payload["event_type"] = weatherEvent.EventType.String()
	event.Payload["severity"] = weatherEvent.Severity.String()
	event.Payload["title"] = weatherEvent.Title
	event.Payload["duration"] = weatherEvent.Duration
	if triggerWeather != nil {
		event.Payload["trigger_weather_type"] = triggerWeather.WeatherType.String()
		event.Payload["trigger_intensity"] = triggerWeather.Intensity.String()
	}
	
	return event
}

// WeatherEventEndedEvent 天气事件结束事件
type WeatherEventEndedEvent struct {
	*BaseDomainEvent
	SceneID      string
	WeatherEvent *WeatherEvent
	EndReason    string
	TotalDuration time.Duration
}

// NewWeatherEventEndedEvent 创建天气事件结束事件
func NewWeatherEventEndedEvent(sceneID string, weatherEvent *WeatherEvent, endReason string) *WeatherEventEndedEvent {
	now := time.Now()
	event := &WeatherEventEndedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("weather_event_ended_%d", now.UnixNano()),
			EventType:   "weather.event_ended",
			AggregateID: sceneID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SceneID:       sceneID,
		WeatherEvent:  weatherEvent,
		EndReason:     endReason,
		TotalDuration: now.Sub(weatherEvent.StartTime),
	}
	
	// 设置载荷
	event.Payload["scene_id"] = sceneID
	event.Payload["event_id"] = weatherEvent.ID
	event.Payload["event_type"] = weatherEvent.EventType.String()
	event.Payload["end_reason"] = endReason
	event.Payload["total_duration"] = now.Sub(weatherEvent.StartTime)
	event.Payload["planned_duration"] = weatherEvent.Duration
	
	return event
}

// WeatherForecastGeneratedEvent 天气预报生成事件
type WeatherForecastGeneratedEvent struct {
	*BaseDomainEvent
	SceneID   string
	Forecasts []*WeatherForecast
	Hours     int
	Confidence float64
}

// NewWeatherForecastGeneratedEvent 创建天气预报生成事件
func NewWeatherForecastGeneratedEvent(sceneID string, forecasts []*WeatherForecast, hours int) *WeatherForecastGeneratedEvent {
	now := time.Now()
	
	// 计算平均置信度
	totalConfidence := 0.0
	for _, forecast := range forecasts {
		totalConfidence += forecast.Confidence
	}
	averageConfidence := totalConfidence / float64(len(forecasts))
	
	event := &WeatherForecastGeneratedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("weather_forecast_generated_%d", now.UnixNano()),
			EventType:   "weather.forecast_generated",
			AggregateID: sceneID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SceneID:    sceneID,
		Forecasts:  forecasts,
		Hours:      hours,
		Confidence: averageConfidence,
	}
	
	// 设置载荷
	event.Payload["scene_id"] = sceneID
	event.Payload["forecast_count"] = len(forecasts)
	event.Payload["hours"] = hours
	event.Payload["average_confidence"] = averageConfidence
	if len(forecasts) > 0 {
		event.Payload["first_forecast_time"] = forecasts[0].Time
		event.Payload["last_forecast_time"] = forecasts[len(forecasts)-1].Time
	}
	
	return event
}

// SeasonChangedEvent 季节变化事件
type SeasonChangedEvent struct {
	*BaseDomainEvent
	ZoneID         string
	PreviousSeason Season
	CurrentSeason  Season
	ChangeTime     time.Time
	Pattern        *SeasonalPattern
}

// NewSeasonChangedEvent 创建季节变化事件
func NewSeasonChangedEvent(zoneID string, previous, current Season, pattern *SeasonalPattern) *SeasonChangedEvent {
	now := time.Now()
	event := &SeasonChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("season_changed_%d", now.UnixNano()),
			EventType:   "weather.season_changed",
			AggregateID: zoneID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		ZoneID:         zoneID,
		PreviousSeason: previous,
		CurrentSeason:  current,
		ChangeTime:     now,
		Pattern:        pattern,
	}
	
	// 设置载荷
	event.Payload["zone_id"] = zoneID
	event.Payload["previous_season"] = previous.String()
	event.Payload["current_season"] = current.String()
	event.Payload["change_time"] = now
	
	return event
}

// WeatherSystemInitializedEvent 天气系统初始化事件
type WeatherSystemInitializedEvent struct {
	*BaseDomainEvent
	SceneID        string
	InitialWeather *WeatherState
	ClimateZone    string
	SeasonalPattern *SeasonalPattern
}

// NewWeatherSystemInitializedEvent 创建天气系统初始化事件
func NewWeatherSystemInitializedEvent(sceneID string, initialWeather *WeatherState, climateZone string, pattern *SeasonalPattern) *WeatherSystemInitializedEvent {
	now := time.Now()
	event := &WeatherSystemInitializedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("weather_system_initialized_%d", now.UnixNano()),
			EventType:   "weather.system_initialized",
			AggregateID: sceneID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SceneID:         sceneID,
		InitialWeather:  initialWeather,
		ClimateZone:     climateZone,
		SeasonalPattern: pattern,
	}
	
	// 设置载荷
	event.Payload["scene_id"] = sceneID
	event.Payload["climate_zone"] = climateZone
	if initialWeather != nil {
		event.Payload["initial_weather_type"] = initialWeather.WeatherType.String()
		event.Payload["initial_intensity"] = initialWeather.Intensity.String()
		event.Payload["initial_temperature"] = initialWeather.Temperature
	}
	if pattern != nil {
		event.Payload["current_season"] = pattern.CurrentSeason.String()
	}
	
	return event
}

// WeatherUpdateFailedEvent 天气更新失败事件
type WeatherUpdateFailedEvent struct {
	*BaseDomainEvent
	SceneID      string
	Error        error
	ErrorMessage string
	RetryCount   int
	LastAttempt  time.Time
}

// NewWeatherUpdateFailedEvent 创建天气更新失败事件
func NewWeatherUpdateFailedEvent(sceneID string, err error, retryCount int) *WeatherUpdateFailedEvent {
	now := time.Now()
	event := &WeatherUpdateFailedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("weather_update_failed_%d", now.UnixNano()),
			EventType:   "weather.update_failed",
			AggregateID: sceneID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SceneID:      sceneID,
		Error:        err,
		ErrorMessage: err.Error(),
		RetryCount:   retryCount,
		LastAttempt:  now,
	}
	
	// 设置载荷
	event.Payload["scene_id"] = sceneID
	event.Payload["error_message"] = err.Error()
	event.Payload["retry_count"] = retryCount
	event.Payload["last_attempt"] = now
	
	return event
}

// WeatherDataCorruptedEvent 天气数据损坏事件
type WeatherDataCorruptedEvent struct {
	*BaseDomainEvent
	SceneID        string
	CorruptedData  string
	CorruptionType string
	RecoveryAction string
}

// NewWeatherDataCorruptedEvent 创建天气数据损坏事件
func NewWeatherDataCorruptedEvent(sceneID, corruptedData, corruptionType, recoveryAction string) *WeatherDataCorruptedEvent {
	now := time.Now()
	event := &WeatherDataCorruptedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("weather_data_corrupted_%d", now.UnixNano()),
			EventType:   "weather.data_corrupted",
			AggregateID: sceneID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SceneID:        sceneID,
		CorruptedData:  corruptedData,
		CorruptionType: corruptionType,
		RecoveryAction: recoveryAction,
	}
	
	// 设置载荷
	event.Payload["scene_id"] = sceneID
	event.Payload["corrupted_data"] = corruptedData
	event.Payload["corruption_type"] = corruptionType
	event.Payload["recovery_action"] = recoveryAction
	
	return event
}

// WeatherAnomalyDetectedEvent 天气异常检测事件
type WeatherAnomalyDetectedEvent struct {
	*BaseDomainEvent
	SceneID       string
	AnomalyType   string
	AnomalyData   map[string]interface{}
	SeverityLevel int
	DetectionTime time.Time
}

// NewWeatherAnomalyDetectedEvent 创建天气异常检测事件
func NewWeatherAnomalyDetectedEvent(sceneID, anomalyType string, anomalyData map[string]interface{}, severityLevel int) *WeatherAnomalyDetectedEvent {
	now := time.Now()
	event := &WeatherAnomalyDetectedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("weather_anomaly_detected_%d", now.UnixNano()),
			EventType:   "weather.anomaly_detected",
			AggregateID: sceneID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SceneID:       sceneID,
		AnomalyType:   anomalyType,
		AnomalyData:   anomalyData,
		SeverityLevel: severityLevel,
		DetectionTime: now,
	}
	
	// 设置载荷
	event.Payload["scene_id"] = sceneID
	event.Payload["anomaly_type"] = anomalyType
	event.Payload["severity_level"] = severityLevel
	event.Payload["detection_time"] = now
	for k, v := range anomalyData {
		event.Payload[k] = v
	}
	
	return event
}

// WeatherConfigurationChangedEvent 天气配置变化事件
type WeatherConfigurationChangedEvent struct {
	*BaseDomainEvent
	SceneID           string
	ConfigurationType string
	OldConfiguration  map[string]interface{}
	NewConfiguration  map[string]interface{}
	ChangedBy         string
}

// NewWeatherConfigurationChangedEvent 创建天气配置变化事件
func NewWeatherConfigurationChangedEvent(sceneID, configurationType string, oldConfig, newConfig map[string]interface{}, changedBy string) *WeatherConfigurationChangedEvent {
	now := time.Now()
	event := &WeatherConfigurationChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("weather_configuration_changed_%d", now.UnixNano()),
			EventType:   "weather.configuration_changed",
			AggregateID: sceneID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SceneID:           sceneID,
		ConfigurationType: configurationType,
		OldConfiguration:  oldConfig,
		NewConfiguration:  newConfig,
		ChangedBy:         changedBy,
	}
	
	// 设置载荷
	event.Payload["scene_id"] = sceneID
	event.Payload["configuration_type"] = configurationType
	event.Payload["changed_by"] = changedBy
	event.Payload["old_configuration"] = oldConfig
	event.Payload["new_configuration"] = newConfig
	
	return event
}

// 事件处理器接口

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(event DomainEvent) error
	CanHandle(eventType string) bool
}

// EventBus 事件总线接口
type EventBus interface {
	Publish(event DomainEvent) error
	Subscribe(eventType string, handler EventHandler) error
	Unsubscribe(eventType string, handler EventHandler) error
}

// 事件存储接口

// EventStore 事件存储接口
type EventStore interface {
	Save(event DomainEvent) error
	Load(aggregateID string) ([]DomainEvent, error)
	LoadFromVersion(aggregateID string, version int) ([]DomainEvent, error)
	LoadByEventType(eventType string, limit int) ([]DomainEvent, error)
	LoadByTimeRange(startTime, endTime time.Time) ([]DomainEvent, error)
}