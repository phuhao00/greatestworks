package dressup

import (
	"github.com/google/uuid"
	"time"
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

// OutfitEquippedEvent 服装装备事件
type OutfitEquippedEvent struct {
	BaseDomainEvent
	PlayerID string     `json:"player_id"`
	OutfitID string     `json:"outfit_id"`
	Slot     OutfitSlot `json:"slot"`
	Outfit   *Outfit    `json:"outfit"`
}

// NewOutfitEquippedEvent 创建服装装备事件
func NewOutfitEquippedEvent(playerID, outfitID string, slot OutfitSlot, outfit *Outfit) *OutfitEquippedEvent {
	return &OutfitEquippedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "OutfitEquipped",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID: playerID,
		OutfitID: outfitID,
		Slot:     slot,
		Outfit:   outfit,
	}
}

// GetEventData 获取事件数据
func (e *OutfitEquippedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id": e.PlayerID,
		"outfit_id": e.OutfitID,
		"slot":      e.Slot,
		"outfit":    e.Outfit,
	}
}

// OutfitUnequippedEvent 服装卸下事件
type OutfitUnequippedEvent struct {
	BaseDomainEvent
	PlayerID string     `json:"player_id"`
	OutfitID string     `json:"outfit_id"`
	Slot     OutfitSlot `json:"slot"`
}

// NewOutfitUnequippedEvent 创建服装卸下事件
func NewOutfitUnequippedEvent(playerID, outfitID string, slot OutfitSlot) *OutfitUnequippedEvent {
	return &OutfitUnequippedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "OutfitUnequipped",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID: playerID,
		OutfitID: outfitID,
		Slot:     slot,
	}
}

// GetEventData 获取事件数据
func (e *OutfitUnequippedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id": e.PlayerID,
		"outfit_id": e.OutfitID,
		"slot":      e.Slot,
	}
}

// OutfitObtainedEvent 服装获得事件
type OutfitObtainedEvent struct {
	BaseDomainEvent
	PlayerID string  `json:"player_id"`
	Outfit   *Outfit `json:"outfit"`
	Source   string  `json:"source"` // 获得来源：shop, quest, drop, etc.
}

// NewOutfitObtainedEvent 创建服装获得事件
func NewOutfitObtainedEvent(playerID string, outfit *Outfit, source string) *OutfitObtainedEvent {
	return &OutfitObtainedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "OutfitObtained",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID: playerID,
		Outfit:   outfit,
		Source:   source,
	}
}

// GetEventData 获取事件数据
func (e *OutfitObtainedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id": e.PlayerID,
		"outfit":    e.Outfit,
		"source":    e.Source,
	}
}

// OutfitUpgradedEvent 服装升级事件
type OutfitUpgradedEvent struct {
	BaseDomainEvent
	PlayerID string         `json:"player_id"`
	OutfitID string         `json:"outfit_id"`
	OldLevel int            `json:"old_level"`
	NewLevel int            `json:"new_level"`
	OldAttrs map[string]int `json:"old_attrs"`
	NewAttrs map[string]int `json:"new_attrs"`
}

// NewOutfitUpgradedEvent 创建服装升级事件
func NewOutfitUpgradedEvent(playerID, outfitID string, oldLevel, newLevel int, oldAttrs, newAttrs map[string]int) *OutfitUpgradedEvent {
	return &OutfitUpgradedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "OutfitUpgraded",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID: playerID,
		OutfitID: outfitID,
		OldLevel: oldLevel,
		NewLevel: newLevel,
		OldAttrs: oldAttrs,
		NewAttrs: newAttrs,
	}
}

// GetEventData 获取事件数据
func (e *OutfitUpgradedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id": e.PlayerID,
		"outfit_id": e.OutfitID,
		"old_level": e.OldLevel,
		"new_level": e.NewLevel,
		"old_attrs": e.OldAttrs,
		"new_attrs": e.NewAttrs,
	}
}

// StyleAppliedEvent 风格应用事件
type StyleAppliedEvent struct {
	BaseDomainEvent
	PlayerID string        `json:"player_id"`
	StyleID  string        `json:"style_id"`
	Style    *DressupStyle `json:"style"`
}

// NewStyleAppliedEvent 创建风格应用事件
func NewStyleAppliedEvent(playerID, styleID string, style *DressupStyle) *StyleAppliedEvent {
	return &StyleAppliedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "StyleApplied",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID: playerID,
		StyleID:  styleID,
		Style:    style,
	}
}

// GetEventData 获取事件数据
func (e *StyleAppliedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id": e.PlayerID,
		"style_id":  e.StyleID,
		"style":     e.Style,
	}
}

// SetBonusActivatedEvent 套装加成激活事件
type SetBonusActivatedEvent struct {
	BaseDomainEvent
	PlayerID   string         `json:"player_id"`
	SetType    string         `json:"set_type"`
	PieceCount int            `json:"piece_count"`
	Bonuses    map[string]int `json:"bonuses"`
}

// NewSetBonusActivatedEvent 创建套装加成激活事件
func NewSetBonusActivatedEvent(playerID, setType string, pieceCount int, bonuses map[string]int) *SetBonusActivatedEvent {
	return &SetBonusActivatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "SetBonusActivated",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:   playerID,
		SetType:    setType,
		PieceCount: pieceCount,
		Bonuses:    bonuses,
	}
}

// GetEventData 获取事件数据
func (e *SetBonusActivatedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id":   e.PlayerID,
		"set_type":    e.SetType,
		"piece_count": e.PieceCount,
		"bonuses":     e.Bonuses,
	}
}
