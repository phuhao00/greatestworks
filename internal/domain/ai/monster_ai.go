package ai

import (
	"context"
	"greatestworks/internal/domain/character"
	"math"
)

// AIState AI状态
type AIState int32

const (
	AIStateIdle   AIState = 0 // 空闲
	AIStateWalk   AIState = 1 // 巡逻
	AIStateChase  AIState = 2 // 追击
	AIStateCast   AIState = 3 // 施法
	AIStateGoback AIState = 4 // 返回
	AIStateHurt   AIState = 5 // 受伤
	AIStateDeath  AIState = 6 // 死亡
)

// MonsterAI 怪物AI
type MonsterAI struct {
	owner *character.Monster

	state          AIState
	stateTime      float32 // 当前状态持续时间
	target         *character.Actor
	patrolRadius   float32 // 巡逻半径
	chaseRadius    float32 // 追击半径
	attackRadius   float32 // 攻击半径
	initPosition   character.Vector3
	currentSkillID int32
}

// NewMonsterAI 创建怪物AI
func NewMonsterAI(owner *character.Monster, patrolRadius, chaseRadius, attackRadius float32) *MonsterAI {
	return &MonsterAI{
		owner:        owner,
		state:        AIStateIdle,
		patrolRadius: patrolRadius,
		chaseRadius:  chaseRadius,
		attackRadius: attackRadius,
	}
}

// Start 初始化AI
func (ai *MonsterAI) Start(ctx context.Context, initPos character.Vector3) error {
	ai.initPosition = initPos
	ai.state = AIStateIdle
	ai.stateTime = 0
	return nil
}

// Update AI更新
func (ai *MonsterAI) Update(ctx context.Context, deltaTime float32) error {
	if ai.owner == nil || ai.owner.IsDeath() {
		ai.state = AIStateDeath
		return nil
	}

	ai.stateTime += deltaTime

	switch ai.state {
	case AIStateIdle:
		ai.updateIdle(ctx, deltaTime)
	case AIStateWalk:
		ai.updateWalk(ctx, deltaTime)
	case AIStateChase:
		ai.updateChase(ctx, deltaTime)
	case AIStateCast:
		ai.updateCast(ctx, deltaTime)
	case AIStateGoback:
		ai.updateGoback(ctx, deltaTime)
	case AIStateHurt:
		ai.updateHurt(ctx, deltaTime)
	}

	return nil
}

// updateIdle 更新空闲状态
func (ai *MonsterAI) updateIdle(ctx context.Context, deltaTime float32) {
	// 检测周围目标
	if ai.detectTarget() {
		ai.changeState(AIStateChase)
		return
	}

	// 空闲一段时间后开始巡逻
	if ai.stateTime > 3.0 {
		ai.changeState(AIStateWalk)
	}
}

// updateWalk 更新巡逻状态
func (ai *MonsterAI) updateWalk(ctx context.Context, deltaTime float32) {
	// 检测目标
	if ai.detectTarget() {
		ai.changeState(AIStateChase)
		return
	}

	// 检查是否离出生点过远
	dist := ai.owner.Position().Distance(ai.initPosition)
	if dist > ai.patrolRadius {
		ai.changeState(AIStateGoback)
		return
	}

	// 巡逻一段时间后停止
	if ai.stateTime > 5.0 {
		ai.changeState(AIStateIdle)
	}

	// TODO: 实际移动逻辑
}

// updateChase 更新追击状态
func (ai *MonsterAI) updateChase(ctx context.Context, deltaTime float32) {
	if ai.target == nil || ai.target.IsDeath() {
		ai.target = nil
		ai.changeState(AIStateIdle)
		return
	}

	// 检查目标是否超出追击范围
	dist := ai.owner.DistanceTo(ai.target.Entity)
	if dist > ai.chaseRadius {
		ai.target = nil
		ai.changeState(AIStateGoback)
		return
	}

	// 检查是否在攻击范围内
	if dist <= ai.attackRadius {
		ai.changeState(AIStateCast)
		return
	}

	// TODO: 追击移动逻辑
}

// updateCast 更新施法状态
func (ai *MonsterAI) updateCast(ctx context.Context, deltaTime float32) {
	if ai.target == nil || ai.target.IsDeath() {
		ai.target = nil
		ai.changeState(AIStateIdle)
		return
	}

	// 检查目标距离
	dist := ai.owner.DistanceTo(ai.target.Entity)
	if dist > ai.attackRadius {
		ai.changeState(AIStateChase)
		return
	}

	// 尝试释放技能
	if ai.stateTime > 0.5 { // 施法间隔
		ai.trySkillCast()
		ai.changeState(AIStateChase)
	}
}

// updateGoback 更新返回状态
func (ai *MonsterAI) updateGoback(ctx context.Context, deltaTime float32) {
	dist := ai.owner.Position().Distance(ai.initPosition)
	if dist < 1.0 {
		// 到达出生点，回满血
		ai.owner.Revive(ctx)
		ai.changeState(AIStateIdle)
		return
	}

	// TODO: 返回移动逻辑
}

// updateHurt 更新受伤状态
func (ai *MonsterAI) updateHurt(ctx context.Context, deltaTime float32) {
	if ai.stateTime > 0.5 {
		ai.changeState(AIStateChase)
	}
}

// OnHurt 受到伤害回调
func (ai *MonsterAI) OnHurt(attacker *character.Actor) {
	if attacker != nil && !attacker.IsDeath() {
		ai.target = attacker
		ai.changeState(AIStateHurt)
	}
}

// detectTarget 检测目标
func (ai *MonsterAI) detectTarget() bool {
	// TODO: 从地图获取附近的玩家
	// 这里需要地图系统支持
	return false
}

// trySkillCast 尝试释放技能
func (ai *MonsterAI) trySkillCast() {
	if ai.target == nil {
		return
	}

	// 获取第一个可用技能
	sm := ai.owner.GetSkillManager()
	// TODO: 从技能管理器选择合适的技能
	_ = sm
}

// changeState 切换状态
func (ai *MonsterAI) changeState(newState AIState) {
	ai.state = newState
	ai.stateTime = 0
}

// GetState 获取当前状态
func (ai *MonsterAI) GetState() AIState {
	return ai.state
}

// MissileAI 投射物AI
type MissileAI struct {
	owner      *character.Missile
	target     *character.Actor
	speed      float32
	maxRange   float32
	travelDist float32
}

// NewMissileAI 创建投射物AI
func NewMissileAI(owner *character.Missile, target *character.Actor, speed, maxRange float32) *MissileAI {
	return &MissileAI{
		owner:    owner,
		target:   target,
		speed:    speed,
		maxRange: maxRange,
	}
}

// Update 更新投射物
func (mai *MissileAI) Update(ctx context.Context, deltaTime float32) error {
	if mai.owner == nil {
		return nil
	}

	// 检查目标
	if mai.target == nil || mai.target.IsDeath() {
		// 目标消失，销毁投射物
		mai.owner.Destroy(ctx)
		return nil
	}

	// 计算移动
	currentPos := mai.owner.Position()
	targetPos := mai.target.Position()

	dx := targetPos.X - currentPos.X
	dy := targetPos.Y - currentPos.Y
	dz := targetPos.Z - currentPos.Z

	dist := float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))

	if dist < 0.5 {
		// 命中目标
		mai.onHit()
		mai.owner.Destroy(ctx)
		return nil
	}

	// 移动
	moveSpeed := mai.speed * deltaTime
	if moveSpeed > dist {
		moveSpeed = dist
	}

	ratio := moveSpeed / dist
	newPos := character.NewVector3(
		currentPos.X+dx*ratio,
		currentPos.Y+dy*ratio,
		currentPos.Z+dz*ratio,
	)

	mai.owner.SetPosition(newPos)
	mai.travelDist += moveSpeed

	// 检查是否超出最大射程
	if mai.travelDist >= mai.maxRange {
		mai.owner.Destroy(ctx)
	}

	return nil
}

// onHit 命中处理
func (mai *MissileAI) onHit() {
	// TODO: 应用伤害或效果到目标
}
