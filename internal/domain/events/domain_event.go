// Package events 领域事件定义
package events

import (
	"encoding/json"
	"time"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	GetEventType() string
	GetAggregateID() string
	GetVersion() int64
	GetOccurredAt() time.Time
	GetData() interface{}
	ToJSON() ([]byte, error)
}

// BaseDomainEvent 基础领域事件
type BaseDomainEvent struct {
	EventType   string      `json:"event_type"`
	AggregateID string      `json:"aggregate_id"`
	Version     int64       `json:"version"`
	OccurredAt  time.Time   `json:"occurred_at"`
	Data        interface{} `json:"data"`
}

// NewBaseDomainEvent 创建基础领域事件
func NewBaseDomainEvent(eventType, aggregateID string, version int64, data interface{}) *BaseDomainEvent {
	return &BaseDomainEvent{
		EventType:   eventType,
		AggregateID: aggregateID,
		Version:     version,
		OccurredAt:  time.Now(),
		Data:        data,
	}
}

// GetEventType 获取事件类型
func (e *BaseDomainEvent) GetEventType() string {
	return e.EventType
}

// GetAggregateID 获取聚合根ID
func (e *BaseDomainEvent) GetAggregateID() string {
	return e.AggregateID
}

// GetVersion 获取版本号
func (e *BaseDomainEvent) GetVersion() int64 {
	return e.Version
}

// GetOccurredAt 获取发生时间
func (e *BaseDomainEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetData 获取事件数据
func (e *BaseDomainEvent) GetData() interface{} {
	return e.Data
}

// ToJSON 转换为JSON
func (e *BaseDomainEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// 玩家相关事件
const (
	EventTypePlayerCreated   = "player.created"
	EventTypePlayerLoggedIn  = "player.logged_in"
	EventTypePlayerLoggedOut = "player.logged_out"
	EventTypePlayerMoved     = "player.moved"
	EventTypePlayerLeveledUp = "player.leveled_up"
	EventTypePlayerDied      = "player.died"
	EventTypePlayerHealed    = "player.healed"
	EventTypePlayerGainedExp = "player.gained_exp"
)

// 战斗相关事件
const (
	EventTypeBattleCreated        = "battle.created"
	EventTypeBattleStarted        = "battle.started"
	EventTypeBattleFinished       = "battle.finished"
	EventTypeBattleCancelled      = "battle.cancelled"
	EventTypePlayerJoinedBattle   = "battle.player_joined"
	EventTypePlayerLeftBattle     = "battle.player_left"
	EventTypeBattleActionExecuted = "battle.action_executed"
)

// PlayerCreatedEvent 玩家创建事件
type PlayerCreatedEvent struct {
	*BaseDomainEvent
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	Level      int    `json:"level"`
}

// NewPlayerCreatedEvent 创建玩家创建事件
func NewPlayerCreatedEvent(playerID, playerName string, level int) *PlayerCreatedEvent {
	baseEvent := NewBaseDomainEvent(EventTypePlayerCreated, playerID, 1, nil)
	return &PlayerCreatedEvent{
		BaseDomainEvent: baseEvent,
		PlayerID:        playerID,
		PlayerName:      playerName,
		Level:           level,
	}
}

// PlayerLoggedInEvent 玩家登录事件
type PlayerLoggedInEvent struct {
	*BaseDomainEvent
	PlayerID string `json:"player_id"`
}

// NewPlayerLoggedInEvent 创建玩家登录事件
func NewPlayerLoggedInEvent(playerID string) *PlayerLoggedInEvent {
	baseEvent := NewBaseDomainEvent(EventTypePlayerLoggedIn, playerID, 1, nil)
	return &PlayerLoggedInEvent{
		BaseDomainEvent: baseEvent,
		PlayerID:        playerID,
	}
}

// PlayerLoggedOutEvent 玩家登出事件
type PlayerLoggedOutEvent struct {
	*BaseDomainEvent
	PlayerID string `json:"player_id"`
}

// NewPlayerLoggedOutEvent 创建玩家登出事件
func NewPlayerLoggedOutEvent(playerID string) *PlayerLoggedOutEvent {
	baseEvent := NewBaseDomainEvent(EventTypePlayerLoggedOut, playerID, 1, nil)
	return &PlayerLoggedOutEvent{
		BaseDomainEvent: baseEvent,
		PlayerID:        playerID,
	}
}

// PlayerMovedEvent 玩家移动事件
type PlayerMovedEvent struct {
	*BaseDomainEvent
	PlayerID string  `json:"player_id"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        float64 `json:"z"`
}

// NewPlayerMovedEvent 创建玩家移动事件
func NewPlayerMovedEvent(playerID string, x, y, z float64) *PlayerMovedEvent {
	baseEvent := NewBaseDomainEvent(EventTypePlayerMoved, playerID, 1, nil)
	return &PlayerMovedEvent{
		BaseDomainEvent: baseEvent,
		PlayerID:        playerID,
		X:               x,
		Y:               y,
		Z:               z,
	}
}

// PlayerLeveledUpEvent 玩家升级事件
type PlayerLeveledUpEvent struct {
	*BaseDomainEvent
	PlayerID string `json:"player_id"`
	OldLevel int    `json:"old_level"`
	NewLevel int    `json:"new_level"`
}

// NewPlayerLeveledUpEvent 创建玩家升级事件
func NewPlayerLeveledUpEvent(playerID string, oldLevel, newLevel int) *PlayerLeveledUpEvent {
	baseEvent := NewBaseDomainEvent(EventTypePlayerLeveledUp, playerID, 1, nil)
	return &PlayerLeveledUpEvent{
		BaseDomainEvent: baseEvent,
		PlayerID:        playerID,
		OldLevel:        oldLevel,
		NewLevel:        newLevel,
	}
}

// BattleCreatedEvent 战斗创建事件
type BattleCreatedEvent struct {
	*BaseDomainEvent
	BattleID   string `json:"battle_id"`
	BattleType string `json:"battle_type"`
	CreatorID  string `json:"creator_id"`
}

// NewBattleCreatedEvent 创建战斗创建事件
func NewBattleCreatedEvent(battleID, battleType, creatorID string) *BattleCreatedEvent {
	baseEvent := NewBaseDomainEvent(EventTypeBattleCreated, battleID, 1, nil)
	return &BattleCreatedEvent{
		BaseDomainEvent: baseEvent,
		BattleID:        battleID,
		BattleType:      battleType,
		CreatorID:       creatorID,
	}
}

// BattleStartedEvent 战斗开始事件
type BattleStartedEvent struct {
	*BaseDomainEvent
	BattleID string `json:"battle_id"`
}

// NewBattleStartedEvent 创建战斗开始事件
func NewBattleStartedEvent(battleID string) *BattleStartedEvent {
	baseEvent := NewBaseDomainEvent(EventTypeBattleStarted, battleID, 1, nil)
	return &BattleStartedEvent{
		BaseDomainEvent: baseEvent,
		BattleID:        battleID,
	}
}

// BattleFinishedEvent 战斗结束事件
type BattleFinishedEvent struct {
	*BaseDomainEvent
	BattleID string  `json:"battle_id"`
	WinnerID *string `json:"winner_id,omitempty"`
	Duration int64   `json:"duration"` // 战斗持续时间（秒）
}

// NewBattleFinishedEvent 创建战斗结束事件
func NewBattleFinishedEvent(battleID string, winnerID *string, duration int64) *BattleFinishedEvent {
	baseEvent := NewBaseDomainEvent(EventTypeBattleFinished, battleID, 1, nil)
	return &BattleFinishedEvent{
		BaseDomainEvent: baseEvent,
		BattleID:        battleID,
		WinnerID:        winnerID,
		Duration:        duration,
	}
}
