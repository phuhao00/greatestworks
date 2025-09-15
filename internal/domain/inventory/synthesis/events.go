package synthesis

import (
	"time"
	"github.com/google/uuid"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	GetEventID() string
	GetEventType() string
	GetAggregateID() string
	GetOccurredAt() time.Time
	GetEventData() interface{}
}

// BaseDomainEvent 基础领域事件
type BaseDomainEvent struct {
	EventID     string    `json:"event_id"`
	EventType   string    `json:"event_type"`
	AggregateID string    `json:"aggregate_id"`
	OccurredAt  time.Time `json:"occurred_at"`
}

// GetEventID 获取事件ID
func (e *BaseDomainEvent) GetEventID() string {
	return e.EventID
}

// GetEventType 获取事件类型
func (e *BaseDomainEvent) GetEventType() string {
	return e.EventType
}

// GetAggregateID 获取聚合根ID
func (e *BaseDomainEvent) GetAggregateID() string {
	return e.AggregateID
}

// GetOccurredAt 获取发生时间
func (e *BaseDomainEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// RecipeLearnedEvent 配方学习事件
type RecipeLearnedEvent struct {
	BaseDomainEvent
	PlayerID string  `json:"player_id"`
	Recipe   *Recipe `json:"recipe"`
	Source   string  `json:"source"` // 学习来源：quest, shop, drop, etc.
}

// NewRecipeLearnedEvent 创建配方学习事件
func NewRecipeLearnedEvent(playerID string, recipe *Recipe, source string) *RecipeLearnedEvent {
	return &RecipeLearnedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "RecipeLearned",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID: playerID,
		Recipe:   recipe,
		Source:   source,
	}
}

// GetEventData 获取事件数据
func (e *RecipeLearnedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id": e.PlayerID,
		"recipe":    e.Recipe,
		"source":    e.Source,
	}
}

// MaterialObtainedEvent 材料获得事件
type MaterialObtainedEvent struct {
	BaseDomainEvent
	PlayerID string    `json:"player_id"`
	Material *Material `json:"material"`
	Quantity int       `json:"quantity"`
	Source   string    `json:"source"`
}

// NewMaterialObtainedEvent 创建材料获得事件
func NewMaterialObtainedEvent(playerID string, material *Material, quantity int, source string) *MaterialObtainedEvent {
	return &MaterialObtainedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "MaterialObtained",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID: playerID,
		Material: material,
		Quantity: quantity,
		Source:   source,
	}
}

// GetEventData 获取事件数据
func (e *MaterialObtainedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id": e.PlayerID,
		"material":  e.Material,
		"quantity":  e.Quantity,
		"source":    e.Source,
	}
}

// SynthesisStartedEvent 合成开始事件
type SynthesisStartedEvent struct {
	BaseDomainEvent
	PlayerID   string    `json:"player_id"`
	RecipeID   string    `json:"recipe_id"`
	Quantity   int       `json:"quantity"`
	StartTime  time.Time `json:"start_time"`
	FinishTime time.Time `json:"finish_time"`
}

// NewSynthesisStartedEvent 创建合成开始事件
func NewSynthesisStartedEvent(playerID, recipeID string, quantity int, craftTime time.Duration) *SynthesisStartedEvent {
	startTime := time.Now()
	return &SynthesisStartedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "SynthesisStarted",
			AggregateID: playerID,
			OccurredAt:  startTime,
		},
		PlayerID:   playerID,
		RecipeID:   recipeID,
		Quantity:   quantity,
		StartTime:  startTime,
		FinishTime: startTime.Add(craftTime),
	}
}

// GetEventData 获取事件数据
func (e *SynthesisStartedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id":   e.PlayerID,
		"recipe_id":   e.RecipeID,
		"quantity":    e.Quantity,
		"start_time":  e.StartTime,
		"finish_time": e.FinishTime,
	}
}

// SynthesisCompletedEvent 合成完成事件
type SynthesisCompletedEvent struct {
	BaseDomainEvent
	PlayerID string           `json:"player_id"`
	RecipeID string           `json:"recipe_id"`
	Quantity int              `json:"quantity"`
	Result   *SynthesisResult `json:"result"`
	Record   *SynthesisRecord `json:"record"`
}

// NewSynthesisCompletedEvent 创建合成完成事件
func NewSynthesisCompletedEvent(playerID, recipeID string, quantity int, result *SynthesisResult, record *SynthesisRecord) *SynthesisCompletedEvent {
	return &SynthesisCompletedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "SynthesisCompleted",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID: playerID,
		RecipeID: recipeID,
		Quantity: quantity,
		Result:   result,
		Record:   record,
	}
}

// GetEventData 获取事件数据
func (e *SynthesisCompletedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id": e.PlayerID,
		"recipe_id": e.RecipeID,
		"quantity":  e.Quantity,
		"result":    e.Result,
		"record":    e.Record,
	}
}

// MaterialConsumedEvent 材料消耗事件
type MaterialConsumedEvent struct {
	BaseDomainEvent
	PlayerID   string `json:"player_id"`
	MaterialID string `json:"material_id"`
	Quantity   int    `json:"quantity"`
	Reason     string `json:"reason"` // 消耗原因：synthesis, upgrade, etc.
}

// NewMaterialConsumedEvent 创建材料消耗事件
func NewMaterialConsumedEvent(playerID, materialID string, quantity int, reason string) *MaterialConsumedEvent {
	return &MaterialConsumedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "MaterialConsumed",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:   playerID,
		MaterialID: materialID,
		Quantity:   quantity,
		Reason:     reason,
	}
}

// GetEventData 获取事件数据
func (e *MaterialConsumedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id":   e.PlayerID,
		"material_id": e.MaterialID,
		"quantity":    e.Quantity,
		"reason":      e.Reason,
	}
}

// SynthesisBonusAppliedEvent 合成加成应用事件
type SynthesisBonusAppliedEvent struct {
	BaseDomainEvent
	PlayerID string            `json:"player_id"`
	RecipeID string            `json:"recipe_id"`
	Bonuses  []*SynthesisBonus `json:"bonuses"`
}

// NewSynthesisBonusAppliedEvent 创建合成加成应用事件
func NewSynthesisBonusAppliedEvent(playerID, recipeID string, bonuses []*SynthesisBonus) *SynthesisBonusAppliedEvent {
	return &SynthesisBonusAppliedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "SynthesisBonusApplied",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID: playerID,
		RecipeID: recipeID,
		Bonuses:  bonuses,
	}
}

// GetEventData 获取事件数据
func (e *SynthesisBonusAppliedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id": e.PlayerID,
		"recipe_id": e.RecipeID,
		"bonuses":   e.Bonuses,
	}
}

// RareItemSynthesizedEvent 稀有物品合成事件
type RareItemSynthesizedEvent struct {
	BaseDomainEvent
	PlayerID string  `json:"player_id"`
	RecipeID string  `json:"recipe_id"`
	ItemID   string  `json:"item_id"`
	Quality  Quality `json:"quality"`
	Quantity int     `json:"quantity"`
}

// NewRareItemSynthesizedEvent 创建稀有物品合成事件
func NewRareItemSynthesizedEvent(playerID, recipeID, itemID string, quality Quality, quantity int) *RareItemSynthesizedEvent {
	return &RareItemSynthesizedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "RareItemSynthesized",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID: playerID,
		RecipeID: recipeID,
		ItemID:   itemID,
		Quality:  quality,
		Quantity: quantity,
	}
}

// GetEventData 获取事件数据
func (e *RareItemSynthesizedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id": e.PlayerID,
		"recipe_id": e.RecipeID,
		"item_id":   e.ItemID,
		"quality":   e.Quality,
		"quantity":  e.Quantity,
	}
}

// SynthesisFailedEvent 合成失败事件
type SynthesisFailedEvent struct {
	BaseDomainEvent
	PlayerID string `json:"player_id"`
	RecipeID string `json:"recipe_id"`
	Quantity int    `json:"quantity"`
	Reason   string `json:"reason"`
}

// NewSynthesisFailedEvent 创建合成失败事件
func NewSynthesisFailedEvent(playerID, recipeID string, quantity int, reason string) *SynthesisFailedEvent {
	return &SynthesisFailedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "SynthesisFailed",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID: playerID,
		RecipeID: recipeID,
		Quantity: quantity,
		Reason:   reason,
	}
}

// GetEventData 获取事件数据
func (e *SynthesisFailedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id": e.PlayerID,
		"recipe_id": e.RecipeID,
		"quantity":  e.Quantity,
		"reason":    e.Reason,
	}
}