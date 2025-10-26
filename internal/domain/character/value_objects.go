package character

import (
	"math"
)

// Vector3 三维向量（位置、方向）
type Vector3 struct {
	X float32
	Y float32
	Z float32
}

// NewVector3 创建新的三维向量
func NewVector3(x, y, z float32) Vector3 {
	return Vector3{X: x, Y: y, Z: z}
}

// Zero 返回零向量
func (v Vector3) Zero() Vector3 {
	return Vector3{X: 0, Y: 0, Z: 0}
}

// ToVector2 转换为二维向量（忽略Y轴）
func (v Vector3) ToVector2() Vector2 {
	return Vector2{X: v.X, Y: v.Z}
}

// Distance 计算两个三维向量之间的距离
func (v Vector3) Distance(other Vector3) float32 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	dz := v.Z - other.Z
	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

// Vector2 二维向量（平面位置）
type Vector2 struct {
	X float32
	Y float32
}

// NewVector2 创建新的二维向量
func NewVector2(x, y float32) Vector2 {
	return Vector2{X: x, Y: y}
}

// ToVector3 转换为三维向量（Y=0）
func (v Vector2) ToVector3() Vector3 {
	return Vector3{X: v.X, Y: 0, Z: v.Y}
}

// Distance 计算两个二维向量之间的距离
func (v Vector2) Distance(other Vector2) float32 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

// Normalized 返回归一化向量
func (v Vector2) Normalized() Vector2 {
	length := float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
	if length == 0 {
		return Vector2{X: 0, Y: 0}
	}
	return Vector2{X: v.X / length, Y: v.Y / length}
}

// Add 向量加法
func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{X: v.X + other.X, Y: v.Y + other.Y}
}

// Sub 向量减法
func (v Vector2) Sub(other Vector2) Vector2 {
	return Vector2{X: v.X - other.X, Y: v.Y - other.Y}
}

// Mul 向量数乘
func (v Vector2) Mul(scalar float32) Vector2 {
	return Vector2{X: v.X * scalar, Y: v.Y * scalar}
}

// EntityID 实体ID值对象
type EntityID int32

// IsValid 检查实体ID是否有效
func (id EntityID) IsValid() bool {
	return id > 0
}

// Int32 转换为int32
func (id EntityID) Int32() int32 {
	return int32(id)
}

// EntityType 实体类型
type EntityType int32

const (
	EntityTypePlayer      EntityType = 0 // 玩家
	EntityTypeMonster     EntityType = 1 // 怪物
	EntityTypeNPC         EntityType = 2 // NPC
	EntityTypeMissile     EntityType = 3 // 投射物
	EntityTypeDroppedItem EntityType = 4 // 掉落物
	EntityTypePet         EntityType = 5 // 宠物
	EntityTypeSummon      EntityType = 6 // 召唤物
)

// String 返回实体类型的字符串表示
func (t EntityType) String() string {
	switch t {
	case EntityTypePlayer:
		return "Player"
	case EntityTypeMonster:
		return "Monster"
	case EntityTypeNPC:
		return "NPC"
	case EntityTypeMissile:
		return "Missile"
	case EntityTypeDroppedItem:
		return "DroppedItem"
	case EntityTypePet:
		return "Pet"
	case EntityTypeSummon:
		return "Summon"
	default:
		return "Unknown"
	}
}

// AnimationState 动画状态
type AnimationState int32

const (
	AnimationStateIdle  AnimationState = 0 // 空闲
	AnimationStateMove  AnimationState = 1 // 移动
	AnimationStateSkill AnimationState = 2 // 释放技能
	AnimationStateHurt  AnimationState = 3 // 受伤
	AnimationStateDeath AnimationState = 4 // 死亡
	AnimationStateJump  AnimationState = 5 // 跳跃
	AnimationStateFall  AnimationState = 6 // 下落
)

// FlagState 状态标志位（可组合）
type FlagState int32

const (
	FlagStateZero       FlagState = 0  // 无状态
	FlagStateStun       FlagState = 1  // 眩晕
	FlagStateRoot       FlagState = 2  // 定身
	FlagStateSilence    FlagState = 4  // 沉默
	FlagStateInvincible FlagState = 8  // 无敌
	FlagStateInvisible  FlagState = 16 // 隐身
	FlagStateDisarm     FlagState = 32 // 缴械
	FlagStateSlow       FlagState = 64 // 减速
)

// HasFlag 检查是否包含某个状态标志
func (f FlagState) HasFlag(flag FlagState) bool {
	return (f & flag) != 0
}

// AddFlag 添加状态标志
func (f FlagState) AddFlag(flag FlagState) FlagState {
	return f | flag
}

// RemoveFlag 移除状态标志
func (f FlagState) RemoveFlag(flag FlagState) FlagState {
	return f & ^flag
}

// Transform 位置和方向
type Transform struct {
	Position  Vector3 // 位置
	Direction Vector3 // 方向
}

// NewTransform 创建新的Transform
func NewTransform(pos, dir Vector3) Transform {
	return Transform{
		Position:  pos,
		Direction: dir,
	}
}
