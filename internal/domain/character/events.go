package character

import (
	"time"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	EventName() string
	OccurredOn() time.Time
	AggregateID() interface{}
}

// BaseDomainEvent 基础领域事件
type BaseDomainEvent struct {
	eventName   string
	occurredOn  time.Time
	aggregateID interface{}
}

func NewBaseDomainEvent(eventName string, aggregateID interface{}) BaseDomainEvent {
	return BaseDomainEvent{
		eventName:   eventName,
		occurredOn:  time.Now(),
		aggregateID: aggregateID,
	}
}

func (e BaseDomainEvent) EventName() string {
	return e.eventName
}

func (e BaseDomainEvent) OccurredOn() time.Time {
	return e.occurredOn
}

func (e BaseDomainEvent) AggregateID() interface{} {
	return e.aggregateID
}

// ========== 实体相关事件 ==========

// EntityCreatedEvent 实体创建事件
type EntityCreatedEvent struct {
	BaseDomainEvent
	EntityID   EntityID
	EntityType EntityType
}

func NewEntityCreatedEvent(entityID EntityID, entityType EntityType) *EntityCreatedEvent {
	return &EntityCreatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("EntityCreated", entityID),
		EntityID:        entityID,
		EntityType:      entityType,
	}
}

// EntityDestroyedEvent 实体销毁事件
type EntityDestroyedEvent struct {
	BaseDomainEvent
	EntityID   EntityID
	EntityType EntityType
}

func NewEntityDestroyedEvent(entityID EntityID, entityType EntityType) *EntityDestroyedEvent {
	return &EntityDestroyedEvent{
		BaseDomainEvent: NewBaseDomainEvent("EntityDestroyed", entityID),
		EntityID:        entityID,
		EntityType:      entityType,
	}
}

// ========== 玩家相关事件 ==========

// PlayerCreatedEvent 玩家创建事件
type PlayerCreatedEvent struct {
	BaseDomainEvent
	CharacterID int64
	UserID      int64
	Name        string
	Level       int32
}

func NewPlayerCreatedEvent(characterID, userID int64, name string, level int32) *PlayerCreatedEvent {
	return &PlayerCreatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("PlayerCreated", characterID),
		CharacterID:     characterID,
		UserID:          userID,
		Name:            name,
		Level:           level,
	}
}

// PlayerLevelUpEvent 玩家升级事件
type PlayerLevelUpEvent struct {
	BaseDomainEvent
	CharacterID int64
	OldLevel    int32
	NewLevel    int32
}

func NewPlayerLevelUpEvent(characterID int64, oldLevel, newLevel int32) *PlayerLevelUpEvent {
	return &PlayerLevelUpEvent{
		BaseDomainEvent: NewBaseDomainEvent("PlayerLevelUp", characterID),
		CharacterID:     characterID,
		OldLevel:        oldLevel,
		NewLevel:        newLevel,
	}
}

// PlayerDeathEvent 玩家死亡事件
type PlayerDeathEvent struct {
	BaseDomainEvent
	CharacterID int64
	KillerID    EntityID
	Position    Vector3
}

func NewPlayerDeathEvent(characterID int64, killerID EntityID, position Vector3) *PlayerDeathEvent {
	return &PlayerDeathEvent{
		BaseDomainEvent: NewBaseDomainEvent("PlayerDeath", characterID),
		CharacterID:     characterID,
		KillerID:        killerID,
		Position:        position,
	}
}

// ========== 战斗相关事件 ==========

// DamageDealtEvent 造成伤害事件
type DamageDealtEvent struct {
	BaseDomainEvent
	AttackerID EntityID
	TargetID   EntityID
	Amount     int32
	DamageType DamageType
	IsCrit     bool
}

func NewDamageDealtEvent(attackerID, targetID EntityID, amount int32, damageType DamageType, isCrit bool) *DamageDealtEvent {
	return &DamageDealtEvent{
		BaseDomainEvent: NewBaseDomainEvent("DamageDealt", attackerID),
		AttackerID:      attackerID,
		TargetID:        targetID,
		Amount:          amount,
		DamageType:      damageType,
		IsCrit:          isCrit,
	}
}

// SkillCastEvent 技能释放事件
type SkillCastEvent struct {
	BaseDomainEvent
	CasterID EntityID
	SkillID  int32
	TargetID EntityID
}

func NewSkillCastEvent(casterID EntityID, skillID int32, targetID EntityID) *SkillCastEvent {
	return &SkillCastEvent{
		BaseDomainEvent: NewBaseDomainEvent("SkillCast", casterID),
		CasterID:        casterID,
		SkillID:         skillID,
		TargetID:        targetID,
	}
}

// BuffAddedEvent Buff添加事件
type BuffAddedEvent struct {
	BaseDomainEvent
	TargetID EntityID
	BuffID   int32
	CasterID EntityID
	Duration float32
}

func NewBuffAddedEvent(targetID EntityID, buffID int32, casterID EntityID, duration float32) *BuffAddedEvent {
	return &BuffAddedEvent{
		BaseDomainEvent: NewBaseDomainEvent("BuffAdded", targetID),
		TargetID:        targetID,
		BuffID:          buffID,
		CasterID:        casterID,
		Duration:        duration,
	}
}

// BuffRemovedEvent Buff移除事件
type BuffRemovedEvent struct {
	BaseDomainEvent
	TargetID EntityID
	BuffID   int32
}

func NewBuffRemovedEvent(targetID EntityID, buffID int32) *BuffRemovedEvent {
	return &BuffRemovedEvent{
		BaseDomainEvent: NewBaseDomainEvent("BuffRemoved", targetID),
		TargetID:        targetID,
		BuffID:          buffID,
	}
}

// ========== 怪物相关事件 ==========

// MonsterDeathEvent 怪物死亡事件
type MonsterDeathEvent struct {
	BaseDomainEvent
	MonsterID EntityID
	KillerID  EntityID
	Position  Vector3
	DropItems []int32 // 掉落物品ID列表
	DropExp   int32   // 掉落经验
}

func NewMonsterDeathEvent(monsterID, killerID EntityID, position Vector3, dropItems []int32, dropExp int32) *MonsterDeathEvent {
	return &MonsterDeathEvent{
		BaseDomainEvent: NewBaseDomainEvent("MonsterDeath", monsterID),
		MonsterID:       monsterID,
		KillerID:        killerID,
		Position:        position,
		DropItems:       dropItems,
		DropExp:         dropExp,
	}
}
