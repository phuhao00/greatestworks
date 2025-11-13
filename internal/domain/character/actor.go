package character

import (
	"context"
	"fmt"
	"sync"
)

// Actor 角色实体 - 具有战斗属性的实体（玩家、怪物等）
// 继承自Entity，添加战斗相关的属性和行为
type Actor struct {
	*Entity // 组合Entity

	mu sync.RWMutex

	// 基础信息
	name  string
	level int32

	// 战斗属性
	hp    float32 // 当前生命值
	mp    float32 // 当前魔法值
	speed float32 // 移动速度

	// 状态
	flagState FlagState // 状态标志位

	// 伤害来源信息
	damageSourceInfo *DamageInfo

	// 子系统（聚合其他值对象或服务）
	attributeManager *AttributeManager // 属性管理器
	skillManager     *SkillManager     // 技能管理器
	buffManager      *BuffManager      // Buff管理器
	spell            *Spell            // 施法器

	// 领域事件发布器（可选注入）
	publisher EventPublisher
}

// DamageInfo 伤害信息
type DamageInfo struct {
	TargetID     EntityID     // 目标ID
	AttackerInfo AttackerInfo // 攻击者信息
	Amount       int32        // 伤害数值
	DamageType   DamageType   // 伤害类型
	IsCrit       bool         // 是否暴击
	IsMiss       bool         // 是否未命中
}

// AttackerInfo 攻击者信息
type AttackerInfo struct {
	AttackerID   EntityID     // 攻击者ID
	AttackerType AttackerType // 攻击者类型
	SkillID      int32        // 技能ID
	BuffID       int32        // BuffID
}

// DamageType 伤害类型
type DamageType int32

const (
	DamageTypeUnknown  DamageType = 0 // 未知
	DamageTypePhysical DamageType = 1 // 物理伤害
	DamageTypeMagical  DamageType = 2 // 魔法伤害
	DamageTypeReal     DamageType = 3 // 真实伤害
	DamageTypeHeal     DamageType = 4 // 治疗
)

// AttackerType 攻击者类型
type AttackerType int32

const (
	AttackerTypeSkill       AttackerType = 0 // 技能攻击
	AttackerTypeBuff        AttackerType = 1 // Buff伤害
	AttackerTypeNormal      AttackerType = 2 // 普通攻击
	AttackerTypeEnvironment AttackerType = 3 // 环境伤害
)

// NewActor 创建新Actor（工厂方法）
func NewActor(
	entityID EntityID,
	entityType EntityType,
	unitID int32,
	position Vector3,
	direction Vector3,
	name string,
	level int32,
) *Actor {
	entity := NewEntity(entityID, entityType, unitID, position, direction)

	actor := &Actor{
		Entity:    entity,
		name:      name,
		level:     level,
		flagState: FlagStateZero,
	}

	// 初始化子系统
	actor.attributeManager = NewAttributeManager(actor)
	actor.skillManager = NewSkillManager(actor)
	actor.buffManager = NewBuffManager(actor)
	actor.spell = NewSpell(actor)

	return actor
}

// ========== 基础信息 ==========

// Name 获取名称
func (a *Actor) Name() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.name
}

// Level 获取等级
func (a *Actor) Level() int32 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.level
}

// SetLevel 设置等级
func (a *Actor) SetLevel(level int32) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.level = level
}

// ========== 战斗属性 ==========

// HP 获取当前生命值
func (a *Actor) HP() float32 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.hp
}

// MP 获取当前魔法值
func (a *Actor) MP() float32 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.mp
}

// Speed 获取移动速度
func (a *Actor) Speed() float32 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.speed
}

// ChangeHP 改变生命值
func (a *Actor) ChangeHP(amount float32) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.hp += amount
	if a.hp <= 0 {
		a.hp = 0
	}
	maxHP := a.attributeManager.Final().MaxHP
	if a.hp > maxHP {
		a.hp = maxHP
	}

	// TODO: 同步属性变化到客户端
	// a.syncAttributeEntry(AttributeTypeHP, int32(a.hp))
}

// ChangeMP 改变魔法值
func (a *Actor) ChangeMP(amount float32) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.mp += amount
	if a.mp <= 0 {
		a.mp = 0
	}
	maxMP := a.attributeManager.Final().MaxMP
	if a.mp > maxMP {
		a.mp = maxMP
	}

	// TODO: 同步属性变化到客户端
	// a.syncAttributeEntry(AttributeTypeMP, int32(a.mp))
}

// IsDeath 是否死亡
func (a *Actor) IsDeath() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.hp <= 0
}

// Revive 复活（由子类重写）
func (a *Actor) Revive(ctx context.Context) error {
	// 基类默认实现：恢复满血满蓝
	maxHP := a.attributeManager.Final().MaxHP
	maxMP := a.attributeManager.Final().MaxMP
	a.ChangeHP(maxHP)
	a.ChangeMP(maxMP)
	return nil
}

// ========== 状态标志 ==========

// FlagState 获取状态标志
func (a *Actor) GetFlagState() FlagState {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.flagState
}

// AddFlagState 添加状态标志
func (a *Actor) AddFlagState(state FlagState) {
	a.mu.Lock()
	defer a.mu.Unlock()

	newState := a.flagState.AddFlag(state)
	if newState == a.flagState {
		return
	}
	a.flagState = newState

	// TODO: 同步状态变化到客户端
	// a.syncAttributeEntry(AttributeTypeFlagState, int32(a.flagState))
}

// RemoveFlagState 移除状态标志
func (a *Actor) RemoveFlagState(state FlagState) {
	a.mu.Lock()
	defer a.mu.Unlock()

	newState := a.flagState.RemoveFlag(state)
	if newState == a.flagState {
		return
	}
	a.flagState = newState

	// TODO: 同步状态变化到客户端
	// a.syncAttributeEntry(AttributeTypeFlagState, int32(a.flagState))
}

// ZeroFlagState 清空状态标志
func (a *Actor) ZeroFlagState() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.flagState == FlagStateZero {
		return
	}
	a.flagState = FlagStateZero

	// TODO: 同步状态变化到客户端
	// a.syncAttributeEntry(AttributeTypeFlagState, int32(a.flagState))
}

// SetFlagStateExact 设置为精确的状态标志（用于根据 Buff 汇总结果覆盖）
func (a *Actor) SetFlagStateExact(state FlagState) {
	a.mu.Lock()
	if a.flagState == state {
		a.mu.Unlock()
		return
	}
	a.flagState = state
	a.mu.Unlock()
	// TODO: 同步状态变化到客户端
}

// ========== 伤害处理 ==========

// OnHurt 受到伤害
func (a *Actor) OnHurt(ctx context.Context, info *DamageInfo) error {
	a.mu.Lock()
	a.damageSourceInfo = info
	a.mu.Unlock()

	// TODO: 广播受伤消息到地图内的玩家
	// mapRef := a.GetMap()
	// if gameMap, ok := mapRef.(*Map); ok {
	//     gameMap.BroadcastEntityHurt(info)
	// }

	// 扣除生命值
	a.ChangeHP(-float32(info.Amount))

	// 发布造成伤害事件（由被伤害方触发发布，聚合ID为攻击者）
	if a.publisher != nil {
		evt := NewDamageDealtEvent(info.AttackerInfo.AttackerID, a.ID(), info.Amount, info.DamageType, info.IsCrit)
		a.publisher.Publish(evt)
	}

	// TODO: 记录日志
	// logger.Debug("%s受到%d的%s攻击, 扣除%d血量, 剩余血量:%f",
	//     a.String(), info.AttackerInfo.AttackerID, info.AttackerInfo.AttackerType, info.Amount, a.hp)

	return nil
}

// DamageSource 获取伤害来源信息
func (a *Actor) DamageSource() *DamageInfo {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.damageSourceInfo
}

// ========== 子系统访问 ==========

// AttributeManager 获取属性管理器
func (a *Actor) GetAttributeManager() *AttributeManager {
	return a.attributeManager
}

// SkillManager 获取技能管理器
func (a *Actor) GetSkillManager() *SkillManager {
	return a.skillManager
}

// BuffManager 获取Buff管理器
func (a *Actor) GetBuffManager() *BuffManager {
	return a.buffManager
}

// Spell 获取施法器
func (a *Actor) GetSpell() *Spell {
	return a.spell
}

// ========== 事件发布 ==========

// EventPublisher 领域事件发布器接口
type EventPublisher interface {
	Publish(event DomainEvent)
}

// SetEventPublisher 注入事件发布器
func (a *Actor) SetEventPublisher(p EventPublisher) { a.publisher = p }

// GetEventPublisher 获取事件发布器
func (a *Actor) GetEventPublisher() EventPublisher { return a.publisher }

// ========== 生命周期 ==========

// Start 初始化Actor
func (a *Actor) Start(ctx context.Context) error {
	// 调用Entity的Start
	if err := a.Entity.Start(ctx); err != nil {
		return err
	}

	// 初始化子系统
	if err := a.attributeManager.Start(ctx); err != nil {
		return fmt.Errorf("attributeManager start failed: %w", err)
	}

	if err := a.skillManager.Start(ctx); err != nil {
		return fmt.Errorf("skillManager start failed: %w", err)
	}

	if err := a.buffManager.Start(ctx); err != nil {
		return fmt.Errorf("buffManager start failed: %w", err)
	}

	// 根据最终属性初始化当前生命、魔法与移动速度
	fin := a.attributeManager.Final()
	a.mu.Lock()
	a.hp = fin.MaxHP
	a.mp = fin.MaxMP
	a.speed = fin.Speed
	a.mu.Unlock()

	return nil
}

// Update 每帧更新
func (a *Actor) Update(ctx context.Context, deltaTime float32) error {
	// 调用Entity的Update
	if err := a.Entity.Update(ctx, deltaTime); err != nil {
		return err
	}

	// 更新子系统
	if err := a.skillManager.Update(ctx, deltaTime); err != nil {
		return err
	}

	if err := a.buffManager.Update(ctx, deltaTime); err != nil {
		return err
	}

	// 按回复速度恢复生命与魔法，并刷新移动速度
	fin := a.attributeManager.Final()
	if fin.HPRegen != 0 || fin.MPRegen != 0 {
		a.ChangeHP(fin.HPRegen * deltaTime)
		a.ChangeMP(fin.MPRegen * deltaTime)
	}
	a.mu.Lock()
	a.speed = fin.Speed
	a.mu.Unlock()

	return nil
}

// String 字符串表示
func (a *Actor) String() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return fmt.Sprintf("%s:\"%s(%d)\"", a.Type().String(), a.name, a.ID())
}
