package mapmanager

import (
	"context"
	"fmt"
	"sync"

	character "greatestworks/internal/domain/character"
)

// Map 地图聚合根
type Map struct {
	mu sync.RWMutex

	id       int32                                    // 地图ID
	name     string                                   // 地图名称
	width    int32                                    // 地图宽度
	height   int32                                    // 地图高度
	entities map[character.EntityID]*character.Entity // 地图内的所有实体

	// AOI系统（简化实现）
	aoiGrid *AOIGrid

	// 视野与广播
	viewRadius  float32
	visibleSets map[character.EntityID]map[character.EntityID]struct{}
	broadcaster BroadcastFn
}

// NewMap 创建地图
func NewMap(id int32, name string, width, height int32) *Map {
	return &Map{
		id:          id,
		name:        name,
		width:       width,
		height:      height,
		entities:    make(map[character.EntityID]*character.Entity),
		aoiGrid:     NewAOIGrid(width, height, 100), // 100单位网格大小
		viewRadius:  200,
		visibleSets: make(map[character.EntityID]map[character.EntityID]struct{}),
	}
}

// BroadcastFn 应用层注入的广播函数：向指定接收者发送某个主题的消息
type BroadcastFn func(recipients []character.EntityID, topic string, payload interface{})

// SetBroadcaster 设置广播函数
func (m *Map) SetBroadcaster(fn BroadcastFn) {
	m.mu.Lock()
	m.broadcaster = fn
	m.mu.Unlock()
}

// Enter 实体进入地图
func (m *Map) Enter(ctx context.Context, entity *character.Entity) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if entity == nil {
		return fmt.Errorf("entity is nil")
	}

	entityID := entity.ID()
	if _, exists := m.entities[entityID]; exists {
		return fmt.Errorf("entity already in map: %d", entityID)
	}

	m.entities[entityID] = entity
	entity.SetMap(m)

	// 添加到AOI网格
	pos := entity.Position2D()
	m.aoiGrid.Add(entityID, pos.X, pos.Y)

	// 初始化可见集并广播出现
	m.refreshVisibilityFor(entityID)
	return nil
}

// Leave 实体离开地图
func (m *Map) Leave(ctx context.Context, entityID character.EntityID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entity, exists := m.entities[entityID]
	if !exists {
		return fmt.Errorf("entity not in map: %d", entityID)
	}

	// 从AOI网格移除
	pos := entity.Position2D()
	m.aoiGrid.Remove(entityID, pos.X, pos.Y)

	// 广播消失并清理可见集
	m.broadcastDisappear(entityID)

	delete(m.entities, entityID)
	entity.SetMap(nil)

	// TODO: 广播实体离开事件给周围玩家
	return nil
}

// UpdatePosition 更新实体位置
func (m *Map) UpdatePosition(entityID character.EntityID, newPos character.Vector3) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entity, exists := m.entities[entityID]
	if !exists {
		return fmt.Errorf("entity not in map: %d", entityID)
	}

	oldPos := entity.Position2D()
	newPos2D := newPos.ToVector2()

	// 更新AOI
	m.aoiGrid.Move(entityID, oldPos.X, oldPos.Y, newPos2D.X, newPos2D.Y)

	// 更新实体位置
	entity.SetPosition(newPos)

	// 刷新视野与广播移动
	m.refreshVisibilityFor(entityID)
	m.broadcastMove(entityID, newPos)
	return nil
}

// GetEntitiesInRange 获取范围内的实体
func (m *Map) GetEntitiesInRange(centerX, centerY, radius float32) []*character.Entity {
	m.mu.RLock()
	defer m.mu.RUnlock()

	nearbyIDs := m.aoiGrid.GetNearby(centerX, centerY, radius)
	entities := make([]*character.Entity, 0, len(nearbyIDs))

	for _, id := range nearbyIDs {
		if entity, exists := m.entities[id]; exists {
			entities = append(entities, entity)
		}
	}

	return entities
}

// GetEntity 获取实体
func (m *Map) GetEntity(entityID character.EntityID) *character.Entity {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.entities[entityID]
}

// GetAllEntities 获取所有实体
func (m *Map) GetAllEntities() []*character.Entity {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entities := make([]*character.Entity, 0, len(m.entities))
	for _, entity := range m.entities {
		entities = append(entities, entity)
	}
	return entities
}

// ID 获取地图ID
func (m *Map) ID() int32 {
	return m.id
}

// Name 获取地图名称
func (m *Map) Name() string {
	return m.name
}

// ===== 视野与广播辅助 =====

// refreshVisibilityFor 计算并更新某实体的可见集，并广播出现/消失
func (m *Map) refreshVisibilityFor(entityID character.EntityID) {
	if m.broadcaster == nil {
		// 未注入广播器，仍然维护可见集以备将来使用
	}
	entity, ok := m.entities[entityID]
	if !ok {
		return
	}
	pos := entity.Position2D()
	nearby := m.aoiGrid.GetNearby(pos.X, pos.Y, m.viewRadius)

	// 构建新集合（不包含自身）
	newSet := make(map[character.EntityID]struct{})
	filtered := make([]character.EntityID, 0, len(nearby))
	for _, id := range nearby {
		if id == entityID {
			continue
		}
		newSet[id] = struct{}{}
		filtered = append(filtered, id)
	}

	oldSet := m.visibleSets[entityID]
	if oldSet == nil {
		oldSet = make(map[character.EntityID]struct{})
	}

	// 计算出现与消失
	appear := make([]character.EntityID, 0)
	disappear := make([]character.EntityID, 0)
	for id := range newSet {
		if _, ok := oldSet[id]; !ok {
			appear = append(appear, id)
		}
	}
	for id := range oldSet {
		if _, ok := newSet[id]; !ok {
			disappear = append(disappear, id)
		}
	}

	// 保存新集
	m.visibleSets[entityID] = newSet

	// 广播给自身：别人出现/消失
	if m.broadcaster != nil {
		if len(appear) > 0 {
			payload := m.buildAppearPayload(appear)
			m.broadcaster([]character.EntityID{entityID}, "entity_appear", payload)
		}
		if len(disappear) > 0 {
			payload := m.buildDisappearPayload(disappear)
			m.broadcaster([]character.EntityID{entityID}, "entity_disappear", payload)
		}
	}

	// 对于 appear 集合中的每个实体，也需要让对方看到“我出现”
	if m.broadcaster != nil && len(appear) > 0 {
		me := m.buildAppearPayload([]character.EntityID{entityID})
		m.broadcaster(appear, "entity_appear", me)
	}
}

// broadcastDisappear 在实体离开地图时通知其可见范围内的对象
func (m *Map) broadcastDisappear(entityID character.EntityID) {
	// 通知所有将其视为可见的实体（这里简化：对所有实体检查其可见集）
	recipients := make([]character.EntityID, 0)
	for viewer, set := range m.visibleSets {
		if viewer == entityID {
			continue
		}
		if _, ok := set[entityID]; ok {
			recipients = append(recipients, viewer)
			delete(set, entityID)
		}
	}
	delete(m.visibleSets, entityID)
	if m.broadcaster != nil && len(recipients) > 0 {
		payload := m.buildDisappearPayload([]character.EntityID{entityID})
		m.broadcaster(recipients, "entity_disappear", payload)
	}
}

// broadcastMove 将移动事件广播给当前可见该实体的对象
func (m *Map) broadcastMove(entityID character.EntityID, pos character.Vector3) {
	viewers := make([]character.EntityID, 0)
	for viewer, set := range m.visibleSets {
		if viewer == entityID {
			continue
		}
		if _, ok := set[entityID]; ok {
			viewers = append(viewers, viewer)
		}
	}
	if m.broadcaster != nil && len(viewers) > 0 {
		payload := EntityMove{ID: entityID, Position: pos}
		m.broadcaster(viewers, "entity_move", payload)
	}
}

// buildAppearPayload 构造出现列表payload
func (m *Map) buildAppearPayload(ids []character.EntityID) []EntityAppear {
	res := make([]EntityAppear, 0, len(ids))
	for _, id := range ids {
		if e, ok := m.entities[id]; ok {
			t := e.GetTransform()
			res = append(res, EntityAppear{ID: id, Position: t.Position, Direction: t.Direction})
		}
	}
	return res
}

// buildDisappearPayload 构造消失列表payload
func (m *Map) buildDisappearPayload(ids []character.EntityID) []EntityDisappear {
	res := make([]EntityDisappear, 0, len(ids))
	for _, id := range ids {
		res = append(res, EntityDisappear{ID: id})
	}
	return res
}

// BroadcastInRange 使用AOI在范围内广播
func (m *Map) BroadcastInRange(x, y, radius float32, topic string, payload interface{}) {
	if m.broadcaster == nil {
		return
	}
	ids := m.aoiGrid.GetNearby(x, y, radius)
	if len(ids) == 0 {
		return
	}
	m.broadcaster(ids, topic, payload)
}

// BroadcastTo 向给定实体列表广播
func (m *Map) BroadcastTo(recipients []character.EntityID, topic string, payload interface{}) {
	if m.broadcaster == nil || len(recipients) == 0 {
		return
	}
	m.broadcaster(recipients, topic, payload)
}

// ===== 广播数据结构 =====
type EntityAppear struct {
	ID        character.EntityID
	Position  character.Vector3
	Direction character.Vector3
}

type EntityDisappear struct {
	ID character.EntityID
}

type EntityMove struct {
	ID       character.EntityID
	Position character.Vector3
}

// AOIGrid AOI网格系统
type AOIGrid struct {
	mu         sync.RWMutex
	width      int32
	height     int32
	gridSize   int32
	gridsX     int32
	gridsY     int32
	grids      map[int32]map[character.EntityID]bool // 网格索引 -> 实体ID集合
	entityGrid map[character.EntityID]int32          // 实体ID -> 网格索引
}

// NewAOIGrid 创建AOI网格
func NewAOIGrid(width, height, gridSize int32) *AOIGrid {
	gridsX := (width + gridSize - 1) / gridSize
	gridsY := (height + gridSize - 1) / gridSize

	return &AOIGrid{
		width:      width,
		height:     height,
		gridSize:   gridSize,
		gridsX:     gridsX,
		gridsY:     gridsY,
		grids:      make(map[int32]map[character.EntityID]bool),
		entityGrid: make(map[character.EntityID]int32),
	}
}

// Add 添加实体到网格
func (aoi *AOIGrid) Add(entityID character.EntityID, x, y float32) {
	aoi.mu.Lock()
	defer aoi.mu.Unlock()

	gridIndex := aoi.getGridIndex(x, y)
	if _, exists := aoi.grids[gridIndex]; !exists {
		aoi.grids[gridIndex] = make(map[character.EntityID]bool)
	}

	aoi.grids[gridIndex][entityID] = true
	aoi.entityGrid[entityID] = gridIndex
}

// Remove 从网格移除实体
func (aoi *AOIGrid) Remove(entityID character.EntityID, x, y float32) {
	aoi.mu.Lock()
	defer aoi.mu.Unlock()

	if gridIndex, exists := aoi.entityGrid[entityID]; exists {
		if grid, ok := aoi.grids[gridIndex]; ok {
			delete(grid, entityID)
		}
		delete(aoi.entityGrid, entityID)
	}
}

// Move 移动实体
func (aoi *AOIGrid) Move(entityID character.EntityID, oldX, oldY, newX, newY float32) {
	oldGridIndex := aoi.getGridIndex(oldX, oldY)
	newGridIndex := aoi.getGridIndex(newX, newY)

	if oldGridIndex == newGridIndex {
		return // 未跨网格
	}

	aoi.mu.Lock()
	defer aoi.mu.Unlock()

	// 从旧网格移除
	if grid, exists := aoi.grids[oldGridIndex]; exists {
		delete(grid, entityID)
	}

	// 添加到新网格
	if _, exists := aoi.grids[newGridIndex]; !exists {
		aoi.grids[newGridIndex] = make(map[character.EntityID]bool)
	}
	aoi.grids[newGridIndex][entityID] = true
	aoi.entityGrid[entityID] = newGridIndex
}

// GetNearby 获取附近的实体
func (aoi *AOIGrid) GetNearby(x, y, radius float32) []character.EntityID {
	aoi.mu.RLock()
	defer aoi.mu.RUnlock()

	// 计算需要检查的网格范围
	gridRadius := int32(radius/float32(aoi.gridSize)) + 1
	centerGridX, centerGridY := aoi.posToGrid(x, y)

	nearbyEntities := make([]character.EntityID, 0)
	visited := make(map[character.EntityID]bool)

	for dx := -gridRadius; dx <= gridRadius; dx++ {
		for dy := -gridRadius; dy <= gridRadius; dy++ {
			gridX := centerGridX + dx
			gridY := centerGridY + dy

			if gridX < 0 || gridX >= aoi.gridsX || gridY < 0 || gridY >= aoi.gridsY {
				continue
			}

			gridIndex := gridY*aoi.gridsX + gridX
			if grid, exists := aoi.grids[gridIndex]; exists {
				for entityID := range grid {
					if !visited[entityID] {
						visited[entityID] = true
						nearbyEntities = append(nearbyEntities, entityID)
					}
				}
			}
		}
	}

	return nearbyEntities
}

// getGridIndex 获取网格索引
func (aoi *AOIGrid) getGridIndex(x, y float32) int32 {
	gridX, gridY := aoi.posToGrid(x, y)
	return gridY*aoi.gridsX + gridX
}

// posToGrid 位置转网格坐标
func (aoi *AOIGrid) posToGrid(x, y float32) (int32, int32) {
	gridX := int32(x) / aoi.gridSize
	gridY := int32(y) / aoi.gridSize

	if gridX < 0 {
		gridX = 0
	}
	if gridX >= aoi.gridsX {
		gridX = aoi.gridsX - 1
	}
	if gridY < 0 {
		gridY = 0
	}
	if gridY >= aoi.gridsY {
		gridY = aoi.gridsY - 1
	}

	return gridX, gridY
}
