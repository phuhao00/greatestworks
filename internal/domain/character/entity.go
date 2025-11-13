package character

import (
	"context"
	"sync"
)

// Entity 实体基类 - 所有游戏实体的基础
// 遵循DDD原则：Entity是具有唯一标识的领域对象
type Entity struct {
	mu sync.RWMutex // 保护并发访问

	// 唯一标识
	entityID   EntityID
	entityType EntityType

	// 定义数据（配置ID）
	unitID int32

	// 空间属性
	transform  Transform
	position2D Vector2 // 2D位置（用于AOI）

	// 状态
	valid bool // 实体是否有效

	// 所属地图（聚合根引用）
	mapRef interface{} // 避免循环依赖，实际类型为 *Map

	// AOI实体引用（基础设施层）
	aoiEntity interface{} // 实际类型为 AOI系统的实体对象
}

// NewEntity 创建新实体（工厂方法）
func NewEntity(
	entityID EntityID,
	entityType EntityType,
	unitID int32,
	position Vector3,
	direction Vector3,
) *Entity {
	return &Entity{
		entityID:   entityID,
		entityType: entityType,
		unitID:     unitID,
		transform: Transform{
			Position:  position,
			Direction: direction,
		},
		position2D: position.ToVector2(),
		valid:      true,
	}
}

// ========== 身份标识 ==========

// ID 获取实体ID
func (e *Entity) ID() EntityID {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.entityID
}

// Type 获取实体类型
func (e *Entity) Type() EntityType {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.entityType
}

// UnitID 获取单位定义ID
func (e *Entity) UnitID() int32 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.unitID
}

// ========== 空间属性 ==========

// Position 获取3D位置
func (e *Entity) Position() Vector3 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.transform.Position
}

// Position2D 获取2D位置
func (e *Entity) Position2D() Vector2 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.position2D
}

// Direction 获取方向
func (e *Entity) Direction() Vector3 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.transform.Direction
}

// Transform 获取Transform
func (e *Entity) GetTransform() Transform {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.transform
}

// SetPosition 设置位置
func (e *Entity) SetPosition(pos Vector3) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.transform.Position = pos
	e.position2D = pos.ToVector2()
}

// SetDirection 设置方向
func (e *Entity) SetDirection(dir Vector3) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.transform.Direction = dir
}

// SetTransform 设置Transform
func (e *Entity) SetTransform(transform Transform) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.transform = transform
	e.position2D = transform.Position.ToVector2()
}

// ========== 状态管理 ==========

// IsValid 检查实体是否有效
func (e *Entity) IsValid() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.valid
}

// Invalidate 使实体失效
func (e *Entity) Invalidate() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.valid = false
}

// ========== 地图关联 ==========

// SetMap 设置所属地图（由基础设施层调用）
func (e *Entity) SetMap(mapRef interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.mapRef = mapRef
}

// GetMap 获取所属地图
func (e *Entity) GetMap() interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.mapRef
}

// ========== AOI关联 ==========

// SetAOIEntity 设置AOI实体（由基础设施层调用）
func (e *Entity) SetAOIEntity(aoiEntity interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.aoiEntity = aoiEntity
}

// GetAOIEntity 获取AOI实体
func (e *Entity) GetAOIEntity() interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.aoiEntity
}

// ========== 生命周期钩子 ==========

// Start 实体初始化（由子类重写）
func (e *Entity) Start(ctx context.Context) error {
	// 基类默认实现为空
	return nil
}

// Update 每帧更新（由子类重写）
func (e *Entity) Update(ctx context.Context, deltaTime float32) error {
	// 基类默认实现为空
	return nil
}

// Destroy 实体销毁（由子类重写）
func (e *Entity) Destroy(ctx context.Context) error {
	e.Invalidate()
	return nil
}

// ========== 工具方法 ==========

// DistanceTo 计算到另一个实体的距离
func (e *Entity) DistanceTo(other *Entity) float32 {
	e.mu.RLock()
	otherPos := other.Position2D()
	myPos := e.position2D
	e.mu.RUnlock()
	return myPos.Distance(otherPos)
}

// String 字符串表示
func (e *Entity) String() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.entityType.String() + ":" + string(rune(e.entityID))
}
