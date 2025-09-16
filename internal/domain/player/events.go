package player

import (
	"time"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	EventID() string
	EventType() string
	AggregateID() string
	OccurredAt() time.Time
	Version() int
}

// BaseEvent 基础事件
type BaseEvent struct {
	eventID     string
	eventType   string
	aggregateID string
	occurredAt  time.Time
	version     int
}

// EventID 获取事件ID
func (e BaseEvent) EventID() string {
	return e.eventID
}

// EventType 获取事件类型
func (e BaseEvent) EventType() string {
	return e.eventType
}

// AggregateID 获取聚合根ID
func (e BaseEvent) AggregateID() string {
	return e.aggregateID
}

// OccurredAt 获取发生时间
func (e BaseEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// Version 获取版本
func (e BaseEvent) Version() int {
	return e.version
}

// PlayerCreatedEvent 玩家创建事件
type PlayerCreatedEvent struct {
	BaseEvent
	PlayerID PlayerID `json:"player_id"`
	Name     string   `json:"name"`
}

// NewPlayerCreatedEvent 创建玩家创建事件
func NewPlayerCreatedEvent(playerID PlayerID, name string) *PlayerCreatedEvent {
	return &PlayerCreatedEvent{
		BaseEvent: BaseEvent{
			eventID:     generateEventID(),
			eventType:   "PlayerCreated",
			aggregateID: playerID.String(),
			occurredAt:  time.Now(),
			version:     1,
		},
		PlayerID: playerID,
		Name:     name,
	}
}

// PlayerLevelUpEvent 玩家升级事件
type PlayerLevelUpEvent struct {
	BaseEvent
	PlayerID PlayerID `json:"player_id"`
	OldLevel int      `json:"old_level"`
	NewLevel int      `json:"new_level"`
}

// NewPlayerLevelUpEvent 创建玩家升级事件
func NewPlayerLevelUpEvent(playerID PlayerID, oldLevel, newLevel int) *PlayerLevelUpEvent {
	return &PlayerLevelUpEvent{
		BaseEvent: BaseEvent{
			eventID:     generateEventID(),
			eventType:   "PlayerLevelUp",
			aggregateID: playerID.String(),
			occurredAt:  time.Now(),
			version:     1,
		},
		PlayerID: playerID,
		OldLevel: oldLevel,
		NewLevel: newLevel,
	}
}

// PlayerOnlineEvent 玩家上线事件
type PlayerOnlineEvent struct {
	BaseEvent
	PlayerID PlayerID `json:"player_id"`
	Position Position `json:"position"`
}

// NewPlayerOnlineEvent 创建玩家上线事件
func NewPlayerOnlineEvent(playerID PlayerID, position Position) *PlayerOnlineEvent {
	return &PlayerOnlineEvent{
		BaseEvent: BaseEvent{
			eventID:     generateEventID(),
			eventType:   "PlayerOnline",
			aggregateID: playerID.String(),
			occurredAt:  time.Now(),
			version:     1,
		},
		PlayerID: playerID,
		Position: position,
	}
}

// PlayerOfflineEvent 玩家下线事件
type PlayerOfflineEvent struct {
	BaseEvent
	PlayerID PlayerID `json:"player_id"`
	Position Position `json:"position"`
}

// NewPlayerOfflineEvent 创建玩家下线事件
func NewPlayerOfflineEvent(playerID PlayerID, position Position) *PlayerOfflineEvent {
	return &PlayerOfflineEvent{
		BaseEvent: BaseEvent{
			eventID:     generateEventID(),
			eventType:   "PlayerOffline",
			aggregateID: playerID.String(),
			occurredAt:  time.Now(),
			version:     1,
		},
		PlayerID: playerID,
		Position: position,
	}
}

// PlayerMovedEvent 玩家移动事件
type PlayerMovedEvent struct {
	BaseEvent
	PlayerID    PlayerID `json:"player_id"`
	OldPosition Position `json:"old_position"`
	NewPosition Position `json:"new_position"`
}

// NewPlayerMovedEvent 创建玩家移动事件
func NewPlayerMovedEvent(playerID PlayerID, oldPos, newPos Position) *PlayerMovedEvent {
	return &PlayerMovedEvent{
		BaseEvent: BaseEvent{
			eventID:     generateEventID(),
			eventType:   "PlayerMoved",
			aggregateID: playerID.String(),
			occurredAt:  time.Now(),
			version:     1,
		},
		PlayerID:    playerID,
		OldPosition: oldPos,
		NewPosition: newPos,
	}
}

// PlayerDiedEvent 玩家死亡事件
type PlayerDiedEvent struct {
	BaseEvent
	PlayerID PlayerID `json:"player_id"`
	Position Position `json:"position"`
	KillerID *PlayerID `json:"killer_id,omitempty"`
}

// NewPlayerDiedEvent 创建玩家死亡事件
func NewPlayerDiedEvent(playerID PlayerID, position Position, killerID *PlayerID) *PlayerDiedEvent {
	return &PlayerDiedEvent{
		BaseEvent: BaseEvent{
			eventID:     generateEventID(),
			eventType:   "PlayerDied",
			aggregateID: playerID.String(),
			occurredAt:  time.Now(),
			version:     1,
		},
		PlayerID: playerID,
		Position: position,
		KillerID: killerID,
	}
}

// generateEventID 生成事件ID
func generateEventID() string {
	return NewPlayerID().String()
}