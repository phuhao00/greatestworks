package character

import (
	"context"
	"fmt"
)

// Monster 怪物实体
type Monster struct {
	*Actor // 继承Actor

	// 初始位置（用于AI返回出生点）
	initPosition Vector3

	// 刷新点定义
	spawnDefine *SpawnDefine

	// AI系统
	ai AI
}

// SpawnDefine 刷新点定义（从配置加载）
type SpawnDefine struct {
	SpawnID     int32   // 刷新点ID
	WalkRange   float32 // 巡逻范围
	ChaseRange  float32 // 追击范围
	AttackRange float32 // 攻击范围
}

// AI 怪物AI接口
type AI interface {
	Start(ctx context.Context) error
	Update(ctx context.Context, deltaTime float32) error
	OnDeath(ctx context.Context) error
}

// NewMonster 创建怪物（工厂方法）
func NewMonster(
	entityID EntityID,
	unitID int32,
	position Vector3,
	direction Vector3,
	name string,
	level int32,
	spawnDefine *SpawnDefine,
) *Monster {
	actor := NewActor(
		entityID,
		EntityTypeMonster,
		unitID,
		position,
		direction,
		name,
		level,
	)

	monster := &Monster{
		Actor:        actor,
		initPosition: position,
		spawnDefine:  spawnDefine,
	}

	// TODO: 根据怪物类型创建对应的AI
	// monster.ai = NewMonsterAI(monster)

	return monster
}

// InitPosition 获取初始位置
func (m *Monster) InitPosition() Vector3 {
	return m.initPosition
}

// SpawnDefine 获取刷新点定义
func (m *Monster) GetSpawnDefine() *SpawnDefine {
	return m.spawnDefine
}

// GetAI 获取AI
func (m *Monster) GetAI() AI {
	return m.ai
}

// Start 初始化怪物
func (m *Monster) Start(ctx context.Context) error {
	// 调用Actor的Start
	if err := m.Actor.Start(ctx); err != nil {
		return err
	}

	// 启动AI
	if m.ai != nil {
		if err := m.ai.Start(ctx); err != nil {
			return fmt.Errorf("ai start failed: %w", err)
		}
	}

	return nil
}

// Update 更新怪物
func (m *Monster) Update(ctx context.Context, deltaTime float32) error {
	// 调用Actor的Update
	if err := m.Actor.Update(ctx, deltaTime); err != nil {
		return err
	}

	// 更新AI
	if m.ai != nil {
		if err := m.ai.Update(ctx, deltaTime); err != nil {
			return err
		}
	}

	return nil
}

// Revive 怪物复活
func (m *Monster) Revive(ctx context.Context) error {
	// 调用Actor的复活逻辑
	if err := m.Actor.Revive(ctx); err != nil {
		return err
	}

	// 重置到出生点
	m.SetPosition(m.initPosition)

	// 清空伤害来源
	m.damageSourceInfo = nil

	return nil
}

// String 字符串表示
func (m *Monster) String() string {
	return fmt.Sprintf("Monster:\"%s(%d)\"", m.Name(), m.ID())
}

// ========== NPC 实体 ==========

// NPC NPC实体
type NPC struct {
	*Entity // NPC不是Actor，因为不参与战斗

	// NPC信息
	npcID int32  // NPC定义ID
	name  string // NPC名称

	// 功能定义
	functions []NPCFunction // NPC功能列表（对话、商店、任务等）
}

// NPCFunction NPC功能类型
type NPCFunction int32

const (
	NPCFunctionDialogue NPCFunction = 0 // 对话
	NPCFunctionShop     NPCFunction = 1 // 商店
	NPCFunctionQuest    NPCFunction = 2 // 任务
	NPCFunctionTeleport NPCFunction = 3 // 传送
	NPCFunctionCraft    NPCFunction = 4 // 制作
)

// NewNPC 创建NPC（工厂方法）
func NewNPC(
	entityID EntityID,
	npcID int32,
	unitID int32,
	position Vector3,
	direction Vector3,
	name string,
	functions []NPCFunction,
) *NPC {
	entity := NewEntity(entityID, EntityTypeNPC, unitID, position, direction)

	return &NPC{
		Entity:    entity,
		npcID:     npcID,
		name:      name,
		functions: functions,
	}
}

// NPCID 获取NPC定义ID
func (n *NPC) NPCID() int32 {
	return n.npcID
}

// Name 获取NPC名称
func (n *NPC) Name() string {
	return n.name
}

// HasFunction 检查是否具有某个功能
func (n *NPC) HasFunction(function NPCFunction) bool {
	for _, f := range n.functions {
		if f == function {
			return true
		}
	}
	return false
}

// Functions 获取所有功能
func (n *NPC) Functions() []NPCFunction {
	return n.functions
}

// String 字符串表示
func (n *NPC) String() string {
	return fmt.Sprintf("NPC:\"%s(%d)\"[NPCID:%d]", n.name, n.ID(), n.npcID)
}

// ========== Missile 投射物实体 ==========

// Missile 投射物（技能子弹等）
type Missile struct {
	*Entity // 投射物是简单实体

	// 投射物信息
	casterID EntityID // 施法者ID
	targetID EntityID // 目标ID（单体目标）
	target   Vector3  // 目标位置（范围技能）

	// 运动参数
	speed    float32 // 飞行速度
	lifetime float32 // 生命周期
	elapsed  float32 // 已存在时间

	// 技能信息
	skillID int32 // 关联的技能ID
}

// NewMissile 创建投射物
func NewMissile(
	entityID EntityID,
	unitID int32,
	position Vector3,
	direction Vector3,
	casterID EntityID,
	skillID int32,
	speed float32,
	lifetime float32,
) *Missile {
	entity := NewEntity(entityID, EntityTypeMissile, unitID, position, direction)

	return &Missile{
		Entity:   entity,
		casterID: casterID,
		skillID:  skillID,
		speed:    speed,
		lifetime: lifetime,
		elapsed:  0,
	}
}

// Update 更新投射物
func (m *Missile) Update(ctx context.Context, deltaTime float32) error {
	// 调用Entity的Update
	if err := m.Entity.Update(ctx, deltaTime); err != nil {
		return err
	}

	// 更新飞行时间
	m.elapsed += deltaTime

	// 检查是否超时
	if m.elapsed >= m.lifetime {
		m.Invalidate()
		return nil
	}

	// TODO: 更新位置
	// 沿方向移动
	// newPos := m.Position().Add(m.Direction().Mul(m.speed * deltaTime))
	// m.SetPosition(newPos)

	// TODO: 检查碰撞
	// if hit detected {
	//     m.Invalidate()
	//     trigger skill effect
	// }

	return nil
}

// String 字符串表示
func (m *Missile) String() string {
	return fmt.Sprintf("Missile(%d)[Skill:%d,Caster:%d]",
		m.ID(), m.skillID, m.casterID)
}
