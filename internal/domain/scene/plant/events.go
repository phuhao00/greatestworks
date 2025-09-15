package plant

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

// CropPlantedEvent 作物种植事件
type CropPlantedEvent struct {
	*BaseDomainEvent
	FarmID   string
	PlotID   string
	CropID   string
	SeedType SeedType
	Quantity int
	Soil     *Soil
	Season   Season
}

// NewCropPlantedEvent 创建作物种植事件
func NewCropPlantedEvent(farmID, plotID, cropID string, seedType SeedType, quantity int, soil *Soil, season Season) *CropPlantedEvent {
	now := time.Now()
	event := &CropPlantedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("crop_planted_%d", now.UnixNano()),
			EventType:   "plant.crop_planted",
			AggregateID: farmID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		FarmID:   farmID,
		PlotID:   plotID,
		CropID:   cropID,
		SeedType: seedType,
		Quantity: quantity,
		Soil:     soil,
		Season:   season,
	}
	
	// 设置载荷
	event.Payload["farm_id"] = farmID
	event.Payload["plot_id"] = plotID
	event.Payload["crop_id"] = cropID
	event.Payload["seed_type"] = seedType.String()
	event.Payload["quantity"] = quantity
	event.Payload["season"] = season.String()
	if soil != nil {
		event.Payload["soil_type"] = soil.Type.String()
		event.Payload["soil_fertility"] = soil.Fertility
	}
	
	return event
}

// CropHarvestedEvent 作物收获事件
type CropHarvestedEvent struct {
	*BaseDomainEvent
	FarmID        string
	CropID        string
	SeedType      SeedType
	Yield         int
	Quality       CropQuality
	Experience    int
	HarvestResult *HarvestResult
}

// NewCropHarvestedEvent 创建作物收获事件
func NewCropHarvestedEvent(farmID, cropID string, harvestResult *HarvestResult) *CropHarvestedEvent {
	now := time.Now()
	event := &CropHarvestedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("crop_harvested_%d", now.UnixNano()),
			EventType:   "plant.crop_harvested",
			AggregateID: farmID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		FarmID:        farmID,
		CropID:        cropID,
		SeedType:      harvestResult.SeedType,
		Yield:         harvestResult.Yield,
		Quality:       harvestResult.Quality,
		Experience:    harvestResult.Experience,
		HarvestResult: harvestResult,
	}
	
	// 设置载荷
	event.Payload["farm_id"] = farmID
	event.Payload["crop_id"] = cropID
	event.Payload["seed_type"] = harvestResult.SeedType.String()
	event.Payload["yield"] = harvestResult.Yield
	event.Payload["quality"] = harvestResult.Quality.String()
	event.Payload["experience"] = harvestResult.Experience
	event.Payload["harvest_time"] = harvestResult.HarvestTime
	
	return event
}

// CropWateredEvent 作物浇水事件
type CropWateredEvent struct {
	*BaseDomainEvent
	FarmID      string
	PlotIDs     []string
	WaterAmount float64
	TotalCrops  int
}

// NewCropWateredEvent 创建作物浇水事件
func NewCropWateredEvent(farmID string, plotIDs []string, waterAmount float64, totalCrops int) *CropWateredEvent {
	now := time.Now()
	event := &CropWateredEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("crop_watered_%d", now.UnixNano()),
			EventType:   "plant.crop_watered",
			AggregateID: farmID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		FarmID:      farmID,
		PlotIDs:     plotIDs,
		WaterAmount: waterAmount,
		TotalCrops:  totalCrops,
	}
	
	// 设置载荷
	event.Payload["farm_id"] = farmID
	event.Payload["plot_ids"] = plotIDs
	event.Payload["water_amount"] = waterAmount
	event.Payload["total_crops"] = totalCrops
	event.Payload["plots_count"] = len(plotIDs)
	
	return event
}

// SoilFertilizedEvent 土壤施肥事件
type SoilFertilizedEvent struct {
	*BaseDomainEvent
	FarmID          string
	Fertilizer      *Fertilizer
	PreviousSoil    *Soil
	UpdatedSoil     *Soil
	FertilityChange float64
}

// NewSoilFertilizedEvent 创建土壤施肥事件
func NewSoilFertilizedEvent(farmID string, fertilizer *Fertilizer, previousSoil, updatedSoil *Soil) *SoilFertilizedEvent {
	now := time.Now()
	fertilityChange := updatedSoil.Fertility - previousSoil.Fertility
	
	event := &SoilFertilizedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("soil_fertilized_%d", now.UnixNano()),
			EventType:   "plant.soil_fertilized",
			AggregateID: farmID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		FarmID:          farmID,
		Fertilizer:      fertilizer,
		PreviousSoil:    previousSoil,
		UpdatedSoil:     updatedSoil,
		FertilityChange: fertilityChange,
	}
	
	// 设置载荷
	event.Payload["farm_id"] = farmID
	event.Payload["fertilizer_type"] = fertilizer.Type.String()
	event.Payload["fertilizer_amount"] = fertilizer.Amount
	event.Payload["previous_fertility"] = previousSoil.Fertility
	event.Payload["updated_fertility"] = updatedSoil.Fertility
	event.Payload["fertility_change"] = fertilityChange
	
	return event
}

// CropGrowthStageChangedEvent 作物生长阶段变化事件
type CropGrowthStageChangedEvent struct {
	*BaseDomainEvent
	FarmID        string
	CropID        string
	SeedType      SeedType
	PreviousStage GrowthStage
	CurrentStage  GrowthStage
	GrowthProgress float64
}

// NewCropGrowthStageChangedEvent 创建作物生长阶段变化事件
func NewCropGrowthStageChangedEvent(farmID, cropID string, seedType SeedType, previousStage, currentStage GrowthStage, progress float64) *CropGrowthStageChangedEvent {
	now := time.Now()
	event := &CropGrowthStageChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("crop_growth_stage_changed_%d", now.UnixNano()),
			EventType:   "plant.crop_growth_stage_changed",
			AggregateID: farmID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		FarmID:         farmID,
		CropID:         cropID,
		SeedType:       seedType,
		PreviousStage:  previousStage,
		CurrentStage:   currentStage,
		GrowthProgress: progress,
	}
	
	// 设置载荷
	event.Payload["farm_id"] = farmID
	event.Payload["crop_id"] = cropID
	event.Payload["seed_type"] = seedType.String()
	event.Payload["previous_stage"] = previousStage.String()
	event.Payload["current_stage"] = currentStage.String()
	event.Payload["growth_progress"] = progress
	
	return event
}

// ToolUsedEvent 工具使用事件
type ToolUsedEvent struct {
	*BaseDomainEvent
	FarmID         string
	ToolID         string
	ToolType       ToolType
	Operation      string
	TargetID       string
	Efficiency     float64
	DurabilityLoss float64
}

// NewToolUsedEvent 创建工具使用事件
func NewToolUsedEvent(farmID, toolID string, toolType ToolType, operation, targetID string, efficiency, durabilityLoss float64) *ToolUsedEvent {
	now := time.Now()
	event := &ToolUsedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("tool_used_%d", now.UnixNano()),
			EventType:   "plant.tool_used",
			AggregateID: farmID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		FarmID:         farmID,
		ToolID:         toolID,
		ToolType:       toolType,
		Operation:      operation,
		TargetID:       targetID,
		Efficiency:     efficiency,
		DurabilityLoss: durabilityLoss,
	}
	
	// 设置载荷
	event.Payload["farm_id"] = farmID
	event.Payload["tool_id"] = toolID
	event.Payload["tool_type"] = toolType.String()
	event.Payload["operation"] = operation
	event.Payload["target_id"] = targetID
	event.Payload["efficiency"] = efficiency
	event.Payload["durability_loss"] = durabilityLoss
	
	return event
}

// FarmExpandedEvent 农场扩展事件
type FarmExpandedEvent struct {
	*BaseDomainEvent
	FarmID      string
	PreviousSize FarmSize
	NewSize     FarmSize
	ExpansionCost *ExpansionCost
	NewPlots    []*Plot
}

// NewFarmExpandedEvent 创建农场扩展事件
func NewFarmExpandedEvent(farmID string, previousSize, newSize FarmSize, cost *ExpansionCost, newPlots []*Plot) *FarmExpandedEvent {
	now := time.Now()
	event := &FarmExpandedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("farm_expanded_%d", now.UnixNano()),
			EventType:   "plant.farm_expanded",
			AggregateID: farmID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		FarmID:        farmID,
		PreviousSize:  previousSize,
		NewSize:       newSize,
		ExpansionCost: cost,
		NewPlots:      newPlots,
	}
	
	// 设置载荷
	event.Payload["farm_id"] = farmID
	event.Payload["previous_size"] = previousSize.String()
	event.Payload["new_size"] = newSize.String()
	event.Payload["new_plots_count"] = len(newPlots)
	if cost != nil {
		event.Payload["expansion_cost_gold"] = cost.Gold
		event.Payload["expansion_cost_materials"] = cost.Materials
		event.Payload["expansion_time"] = cost.Time
	}
	
	return event
}

// PestDiseaseDetectedEvent 病虫害检测事件
type PestDiseaseDetectedEvent struct {
	*BaseDomainEvent
	FarmID      string
	CropID      string
	SeedType    SeedType
	PestDisease *PestDiseaseEvent
	Severity    string
	AffectedArea float64
}

// NewPestDiseaseDetectedEvent 创建病虫害检测事件
func NewPestDiseaseDetectedEvent(farmID, cropID string, seedType SeedType, pestDisease *PestDiseaseEvent, affectedArea float64) *PestDiseaseDetectedEvent {
	now := time.Now()
	event := &PestDiseaseDetectedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("pest_disease_detected_%d", now.UnixNano()),
			EventType:   "plant.pest_disease_detected",
			AggregateID: farmID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		FarmID:       farmID,
		CropID:       cropID,
		SeedType:     seedType,
		PestDisease:  pestDisease,
		Severity:     pestDisease.Severity,
		AffectedArea: affectedArea,
	}
	
	// 设置载荷
	event.Payload["farm_id"] = farmID
	event.Payload["crop_id"] = cropID
	event.Payload["seed_type"] = seedType.String()
	event.Payload["pest_disease_name"] = pestDisease.Name
	event.Payload["severity"] = pestDisease.Severity
	event.Payload["affected_area"] = affectedArea
	event.Payload["occurred_at"] = pestDisease.OccurredAt
	
	return event
}

// SeasonChangedEvent 季节变化事件
type SeasonChangedEvent struct {
	*BaseDomainEvent
	FarmID         string
	PreviousSeason Season
	CurrentSeason  Season
	SeasonModifier *SeasonModifier
	AffectedCrops  []string
}

// NewSeasonChangedEvent 创建季节变化事件
func NewSeasonChangedEvent(farmID string, previousSeason, currentSeason Season, modifier *SeasonModifier, affectedCrops []string) *SeasonChangedEvent {
	now := time.Now()
	event := &SeasonChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("season_changed_%d", now.UnixNano()),
			EventType:   "plant.season_changed",
			AggregateID: farmID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		FarmID:         farmID,
		PreviousSeason: previousSeason,
		CurrentSeason:  currentSeason,
		SeasonModifier: modifier,
		AffectedCrops:  affectedCrops,
	}
	
	// 设置载荷
	event.Payload["farm_id"] = farmID
	event.Payload["previous_season"] = previousSeason.String()
	event.Payload["current_season"] = currentSeason.String()
	event.Payload["affected_crops_count"] = len(affectedCrops)
	if modifier != nil {
		event.Payload["growth_multiplier"] = modifier.GrowthMultiplier
		event.Payload["yield_multiplier"] = modifier.YieldMultiplier
		event.Payload["quality_multiplier"] = modifier.QualityMultiplier
	}
	
	return event
}

// CropHealthChangedEvent 作物健康变化事件
type CropHealthChangedEvent struct {
	*BaseDomainEvent
	FarmID         string
	CropID         string
	SeedType       SeedType
	PreviousHealth float64
	CurrentHealth  float64
	HealthChange   float64
	Cause          string
	Problems       []string
}

// NewCropHealthChangedEvent 创建作物健康变化事件
func NewCropHealthChangedEvent(farmID, cropID string, seedType SeedType, previousHealth, currentHealth float64, cause string, problems []string) *CropHealthChangedEvent {
	now := time.Now()
	healthChange := currentHealth - previousHealth
	
	event := &CropHealthChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("crop_health_changed_%d", now.UnixNano()),
			EventType:   "plant.crop_health_changed",
			AggregateID: farmID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		FarmID:         farmID,
		CropID:         cropID,
		SeedType:       seedType,
		PreviousHealth: previousHealth,
		CurrentHealth:  currentHealth,
		HealthChange:   healthChange,
		Cause:          cause,
		Problems:       problems,
	}
	
	// 设置载荷
	event.Payload["farm_id"] = farmID
	event.Payload["crop_id"] = cropID
	event.Payload["seed_type"] = seedType.String()
	event.Payload["previous_health"] = previousHealth
	event.Payload["current_health"] = currentHealth
	event.Payload["health_change"] = healthChange
	event.Payload["cause"] = cause
	event.Payload["problems"] = problems
	event.Payload["problems_count"] = len(problems)
	
	return event
}

// AutomationTriggeredEvent 自动化触发事件
type AutomationTriggeredEvent struct {
	*BaseDomainEvent
	FarmID         string
	AutomationType string
	TriggerReason  string
	TargetCrops    []string
	ActionTaken    string
	ResourcesUsed  map[string]float64
}

// NewAutomationTriggeredEvent 创建自动化触发事件
func NewAutomationTriggeredEvent(farmID, automationType, triggerReason string, targetCrops []string, actionTaken string, resourcesUsed map[string]float64) *AutomationTriggeredEvent {
	now := time.Now()
	event := &AutomationTriggeredEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("automation_triggered_%d", now.UnixNano()),
			EventType:   "plant.automation_triggered",
			AggregateID: farmID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		FarmID:         farmID,
		AutomationType: automationType,
		TriggerReason:  triggerReason,
		TargetCrops:    targetCrops,
		ActionTaken:    actionTaken,
		ResourcesUsed:  resourcesUsed,
	}
	
	// 设置载荷
	event.Payload["farm_id"] = farmID
	event.Payload["automation_type"] = automationType
	event.Payload["trigger_reason"] = triggerReason
	event.Payload["target_crops"] = targetCrops
	event.Payload["target_crops_count"] = len(targetCrops)
	event.Payload["action_taken"] = actionTaken
	event.Payload["resources_used"] = resourcesUsed
	
	return event
}

// FarmValueChangedEvent 农场价值变化事件
type FarmValueChangedEvent struct {
	*BaseDomainEvent
	FarmID        string
	PreviousValue float64
	CurrentValue  float64
	ValueChange   float64
	ChangeReason  string
	Contributors  map[string]float64
}

// NewFarmValueChangedEvent 创建农场价值变化事件
func NewFarmValueChangedEvent(farmID string, previousValue, currentValue float64, changeReason string, contributors map[string]float64) *FarmValueChangedEvent {
	now := time.Now()
	valueChange := currentValue - previousValue
	
	event := &FarmValueChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("farm_value_changed_%d", now.UnixNano()),
			EventType:   "plant.farm_value_changed",
			AggregateID: farmID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		FarmID:        farmID,
		PreviousValue: previousValue,
		CurrentValue:  currentValue,
		ValueChange:   valueChange,
		ChangeReason:  changeReason,
		Contributors:  contributors,
	}
	
	// 设置载荷
	event.Payload["farm_id"] = farmID
	event.Payload["previous_value"] = previousValue
	event.Payload["current_value"] = currentValue
	event.Payload["value_change"] = valueChange
	event.Payload["change_reason"] = changeReason
	event.Payload["contributors"] = contributors
	
	return event
}

// PlotStatusChangedEvent 地块状态变化事件
type PlotStatusChangedEvent struct {
	*BaseDomainEvent
	FarmID           string
	PlotID           string
	PreviousStatus   string
	CurrentStatus    string
	CropID           string
	StatusChangeTime time.Time
}

// NewPlotStatusChangedEvent 创建地块状态变化事件
func NewPlotStatusChangedEvent(farmID, plotID, previousStatus, currentStatus, cropID string) *PlotStatusChangedEvent {
	now := time.Now()
	event := &PlotStatusChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("plot_status_changed_%d", now.UnixNano()),
			EventType:   "plant.plot_status_changed",
			AggregateID: farmID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		FarmID:           farmID,
		PlotID:           plotID,
		PreviousStatus:   previousStatus,
		CurrentStatus:    currentStatus,
		CropID:           cropID,
		StatusChangeTime: now,
	}
	
	// 设置载荷
	event.Payload["farm_id"] = farmID
	event.Payload["plot_id"] = plotID
	event.Payload["previous_status"] = previousStatus
	event.Payload["current_status"] = currentStatus
	event.Payload["crop_id"] = cropID
	event.Payload["status_change_time"] = now
	
	return event
}

// ResourcesConsumedEvent 资源消耗事件
type ResourcesConsumedEvent struct {
	*BaseDomainEvent
	FarmID           string
	ResourceType     string
	AmountConsumed   float64
	RemainingAmount  float64
	ConsumptionReason string
	RelatedCropID    string
}

// NewResourcesConsumedEvent 创建资源消耗事件
func NewResourcesConsumedEvent(farmID, resourceType string, amountConsumed, remainingAmount float64, reason, relatedCropID string) *ResourcesConsumedEvent {
	now := time.Now()
	event := &ResourcesConsumedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("resources_consumed_%d", now.UnixNano()),
			EventType:   "plant.resources_consumed",
			AggregateID: farmID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		FarmID:            farmID,
		ResourceType:      resourceType,
		AmountConsumed:    amountConsumed,
		RemainingAmount:   remainingAmount,
		ConsumptionReason: reason,
		RelatedCropID:     relatedCropID,
	}
	
	// 设置载荷
	event.Payload["farm_id"] = farmID
	event.Payload["resource_type"] = resourceType
	event.Payload["amount_consumed"] = amountConsumed
	event.Payload["remaining_amount"] = remainingAmount
	event.Payload["consumption_reason"] = reason
	event.Payload["related_crop_id"] = relatedCropID
	
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