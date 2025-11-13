package scene

import (
	// "errors"
	"fmt"
	"math"
	"sync"
	"time"
)

// Scene 场景聚合根
type Scene struct {
	id             string
	name           string
	sceneType      SceneType
	status         SceneStatus
	width          float64
	height         float64
	maxPlayers     int
	currentPlayers int
	entities       map[string]Entity
	players        map[string]*Player
	npcs           map[string]*NPC
	monsters       map[string]*Monster
	items          map[string]*Item
	portals        map[string]*Portal
	aoi            *AOIManager
	spawnPoints    []*SpawnPoint
	lastUpdate     time.Time
	events         []DomainEvent
	mu             sync.RWMutex
}

// NewScene 创建新场景
func NewScene(id, name string, sceneType SceneType, width, height float64, maxPlayers int) *Scene {
	return &Scene{
		id:             id,
		name:           name,
		sceneType:      sceneType,
		status:         SceneStatusActive,
		width:          width,
		height:         height,
		maxPlayers:     maxPlayers,
		currentPlayers: 0,
		entities:       make(map[string]Entity),
		players:        make(map[string]*Player),
		npcs:           make(map[string]*NPC),
		monsters:       make(map[string]*Monster),
		items:          make(map[string]*Item),
		portals:        make(map[string]*Portal),
		aoi:            NewAOIManager(width, height, 100.0), // 默认AOI半径100
		spawnPoints:    make([]*SpawnPoint, 0),
		lastUpdate:     time.Now(),
		events:         make([]DomainEvent, 0),
	}
}

// SceneType 场景类型
type SceneType int

const (
	SceneTypeCity SceneType = iota + 1
	SceneTypeDungeon
	SceneTypeBattlefield
	SceneTypeWilderness
	SceneTypeInstance
	SceneTypeGuild
	SceneTypePvP
	SceneTypeRaid
)

// SceneStatus 场景状态
type SceneStatus int

const (
	SceneStatusActive SceneStatus = iota + 1
	SceneStatusMaintenance
	SceneStatusClosed
	SceneStatusFull
)

// Entity 实体接口
type Entity interface {
	GetID() string
	GetPosition() *Position
	SetPosition(*Position)
	GetEntityType() EntityType
	Update(deltaTime time.Duration)
	IsActive() bool
}

// EntityType 实体类型
type EntityType int

const (
	EntityTypePlayer EntityType = iota + 1
	EntityTypeNPC
	EntityTypeMonster
	EntityTypeItem
	EntityTypePortal
	EntityTypeBuilding
	EntityTypeProjectile
)

// Position 位置
type Position struct {
	X         float64   `json:"x"`
	Y         float64   `json:"y"`
	Z         float64   `json:"z"`
	Direction float64   `json:"direction"` // 朝向角度
	UpdatedAt time.Time `json:"updated_at"`
}

// NewPosition 创建新位置
func NewPosition(x, y, z, direction float64) *Position {
	return &Position{
		X:         x,
		Y:         y,
		Z:         z,
		Direction: direction,
		UpdatedAt: time.Now(),
	}
}

// Distance 计算距离
func (p *Position) Distance(other *Position) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	dz := p.Z - other.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// Player 玩家实体
type Player struct {
	id         string
	name       string
	level      int
	position   *Position
	health     int64
	maxHealth  int64
	mana       int64
	maxMana    int64
	status     PlayerStatus
	lastAction time.Time
	active     bool
}

// PlayerStatus 玩家状态
type PlayerStatus int

const (
	PlayerStatusNormal PlayerStatus = iota + 1
	PlayerStatusCombat
	PlayerStatusDead
	PlayerStatusAFK
	PlayerStatusTrading
	PlayerStatusCasting
)

// NPC 非玩家角色实体
type NPC struct {
	id        string
	name      string
	npcType   NPCType
	position  *Position
	health    int64
	maxHealth int64
	status    NPCStatus
	ai        *AIBehavior
	active    bool
}

// NPCType NPC类型
type NPCType int

const (
	NPCTypeVendor NPCType = iota + 1
	NPCTypeGuard
	NPCTypeQuest
	NPCTypeTrainer
	NPCTypeBanker
	NPCTypeTransporter
)

// NPCStatus NPC状态
type NPCStatus int

const (
	NPCStatusIdle NPCStatus = iota + 1
	NPCStatusPatrolling
	NPCStatusCombat
	NPCStatusInteracting
	NPCStatusDead
)

// Monster 怪物实体
type Monster struct {
	id          string
	name        string
	monsterType MonsterType
	position    *Position
	health      int64
	maxHealth   int64
	level       int
	status      MonsterStatus
	ai          *AIBehavior
	spawnPoint  *SpawnPoint
	lastAttack  time.Time
	active      bool
}

// MonsterType 怪物类型
type MonsterType int

const (
	MonsterTypeNormal MonsterType = iota + 1
	MonsterTypeElite
	MonsterTypeBoss
	MonsterTypeWorldBoss
	MonsterTypeMinion
)

// MonsterStatus 怪物状态
type MonsterStatus int

const (
	MonsterStatusIdle MonsterStatus = iota + 1
	MonsterStatusPatrolling
	MonsterStatusCombat
	MonsterStatusChasing
	MonsterStatusReturning
	MonsterStatusDead
	MonsterStatusRespawning
)

// Item 物品实体
type Item struct {
	id         string
	itemID     string // 物品模板ID
	position   *Position
	quantity   int64
	owner      string // 拾取者限制
	expireTime *time.Time
	active     bool
}

// Portal 传送门实体
type Portal struct {
	id             string
	name           string
	position       *Position
	targetSceneID  string
	targetPosition *Position
	requiredLevel  int
	requiredItems  []string
	cost           int64
	active         bool
}

// SpawnPoint 刷新点
type SpawnPoint struct {
	id           string
	position     *Position
	spawnType    SpawnType
	targetID     string // 刷新的实体ID
	interval     time.Duration
	lastSpawn    time.Time
	maxCount     int
	currentCount int
	active       bool
}

// SpawnType 刷新类型
type SpawnType int

const (
	SpawnTypeMonster SpawnType = iota + 1
	SpawnTypeNPC
	SpawnTypeItem
	SpawnTypePlayer
)

// AIBehavior AI行为
type AIBehavior struct {
	behaviorType  BehaviorType
	patrolPath    []*Position
	currentTarget string
	aggroRange    float64
	chaseRange    float64
	returnRange   float64
	attackRange   float64
	lastUpdate    time.Time
}

// BehaviorType 行为类型
type BehaviorType int

const (
	BehaviorTypeIdle BehaviorType = iota + 1
	BehaviorTypePatrol
	BehaviorTypeGuard
	BehaviorTypeAggressive
	BehaviorTypeDefensive
	BehaviorTypeFlee
)

// AOIManager AOI管理器
type AOIManager struct {
	width    float64
	height   float64
	gridSize float64
	grids    map[string]*AOIGrid
	entities map[string]*AOIEntity
	mu       sync.RWMutex
}

// NewAOIManager 创建AOI管理器
func NewAOIManager(width, height, gridSize float64) *AOIManager {
	return &AOIManager{
		width:    width,
		height:   height,
		gridSize: gridSize,
		grids:    make(map[string]*AOIGrid),
		entities: make(map[string]*AOIEntity),
	}
}

// AOIGrid AOI网格
type AOIGrid struct {
	x        int
	y        int
	entities map[string]*AOIEntity
}

// AOIEntity AOI实体
type AOIEntity struct {
	id       string
	entity   Entity
	gridX    int
	gridY    int
	watchers map[string]bool // 观察者列表
}

// DomainEvent 领域事件接口
type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
	SceneID() string
}

// PlayerEnteredEvent 玩家进入场景事件
type PlayerEnteredEvent struct {
	sceneID    string
	playerID   string
	playerName string
	position   *Position
	occurredAt time.Time
}

func (e PlayerEnteredEvent) EventType() string     { return "player.entered" }
func (e PlayerEnteredEvent) OccurredAt() time.Time { return e.occurredAt }
func (e PlayerEnteredEvent) SceneID() string       { return e.sceneID }

// PlayerLeftEvent 玩家离开场景事件
type PlayerLeftEvent struct {
	sceneID    string
	playerID   string
	playerName string
	occurredAt time.Time
}

func (e PlayerLeftEvent) EventType() string     { return "player.left" }
func (e PlayerLeftEvent) OccurredAt() time.Time { return e.occurredAt }
func (e PlayerLeftEvent) SceneID() string       { return e.sceneID }

// EntityMovedEvent 实体移动事件
type EntityMovedEvent struct {
	sceneID     string
	entityID    string
	entityType  EntityType
	oldPosition *Position
	newPosition *Position
	occurredAt  time.Time
}

func (e EntityMovedEvent) EventType() string     { return "entity.moved" }
func (e EntityMovedEvent) OccurredAt() time.Time { return e.occurredAt }
func (e EntityMovedEvent) SceneID() string       { return e.sceneID }

// MonsterSpawnedEvent 怪物刷新事件
type MonsterSpawnedEvent struct {
	sceneID     string
	monsterID   string
	monsterType MonsterType
	position    *Position
	occurredAt  time.Time
}

func (e MonsterSpawnedEvent) EventType() string     { return "monster.spawned" }
func (e MonsterSpawnedEvent) OccurredAt() time.Time { return e.occurredAt }
func (e MonsterSpawnedEvent) SceneID() string       { return e.sceneID }

// ItemDroppedEvent 物品掉落事件
type ItemDroppedEvent struct {
	sceneID    string
	itemID     string
	itemType   string
	position   *Position
	quantity   int64
	occurredAt time.Time
}

func (e ItemDroppedEvent) EventType() string     { return "item.dropped" }
func (e ItemDroppedEvent) OccurredAt() time.Time { return e.occurredAt }
func (e ItemDroppedEvent) SceneID() string       { return e.sceneID }

// Scene 业务方法实现

// ID 获取场景ID
func (s *Scene) ID() string {
	return s.id
}

// Name 获取场景名称
func (s *Scene) Name() string {
	return s.name
}

// Status 获取场景状态
func (s *Scene) Status() SceneStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}

// Type 获取场景类型
func (s *Scene) Type() SceneType {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sceneType
}

// Width 获取场景宽度
func (s *Scene) GetWidth() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.width
}

// Height 获取场景高度
func (s *Scene) GetHeight() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.height
}

// MaxPlayers 获取最大玩家数
func (s *Scene) GetMaxPlayers() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.maxPlayers
}

// PlayerCount 获取当前玩家数量
func (s *Scene) PlayerCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentPlayers
}

// AddPlayer 添加玩家到场景
func (s *Scene) AddPlayer(player *Player) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status != SceneStatusActive {
		return ErrSceneNotActive
	}

	if s.currentPlayers >= s.maxPlayers {
		return ErrSceneFull
	}

	if _, exists := s.players[player.id]; exists {
		return ErrPlayerAlreadyInScene
	}

	// 添加玩家
	s.players[player.id] = player
	s.entities[player.id] = player
	s.currentPlayers++

	// 添加到AOI
	s.aoi.AddEntity(player.id, player)

	s.lastUpdate = time.Now()

	// 发布事件
	s.addEvent(PlayerEnteredEvent{
		sceneID:    s.id,
		playerID:   player.id,
		playerName: player.name,
		position:   player.position,
		occurredAt: time.Now(),
	})

	return nil
}

// RemovePlayer 从场景移除玩家
func (s *Scene) RemovePlayer(playerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	player, exists := s.players[playerID]
	if !exists {
		return ErrPlayerNotInScene
	}

	// 从AOI移除
	s.aoi.RemoveEntity(playerID)

	// 移除玩家
	delete(s.players, playerID)
	delete(s.entities, playerID)
	s.currentPlayers--

	s.lastUpdate = time.Now()

	// 发布事件
	s.addEvent(PlayerLeftEvent{
		sceneID:    s.id,
		playerID:   playerID,
		playerName: player.name,
		occurredAt: time.Now(),
	})

	return nil
}

// MoveEntity 移动实体
func (s *Scene) MoveEntity(entityID string, newPosition *Position) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entity, exists := s.entities[entityID]
	if !exists {
		return ErrEntityNotFound
	}

	// 检查位置是否有效
	if !s.isValidPosition(newPosition) {
		return ErrInvalidPosition
	}

	oldPosition := entity.GetPosition()
	entity.SetPosition(newPosition)

	// 更新AOI
	s.aoi.UpdateEntity(entityID, newPosition)

	s.lastUpdate = time.Now()

	// 发布事件
	s.addEvent(EntityMovedEvent{
		sceneID:     s.id,
		entityID:    entityID,
		entityType:  entity.GetEntityType(),
		oldPosition: oldPosition,
		newPosition: newPosition,
		occurredAt:  time.Now(),
	})

	return nil
}

// SpawnMonster 刷新怪物
func (s *Scene) SpawnMonster(monster *Monster, spawnPoint *SpawnPoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.monsters[monster.id]; exists {
		return ErrMonsterAlreadyExists
	}

	// 设置刷新点
	monster.spawnPoint = spawnPoint
	monster.position = &Position{
		X:         spawnPoint.position.X,
		Y:         spawnPoint.position.Y,
		Z:         spawnPoint.position.Z,
		Direction: spawnPoint.position.Direction,
		UpdatedAt: time.Now(),
	}

	// 添加怪物
	s.monsters[monster.id] = monster
	s.entities[monster.id] = monster

	// 添加到AOI
	s.aoi.AddEntity(monster.id, monster)

	// 更新刷新点
	spawnPoint.currentCount++
	spawnPoint.lastSpawn = time.Now()

	s.lastUpdate = time.Now()

	// 发布事件
	s.addEvent(MonsterSpawnedEvent{
		sceneID:     s.id,
		monsterID:   monster.id,
		monsterType: monster.monsterType,
		position:    monster.position,
		occurredAt:  time.Now(),
	})

	return nil
}

// DropItem 掉落物品
func (s *Scene) DropItem(item *Item, position *Position) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[item.id]; exists {
		return ErrItemAlreadyExists
	}

	item.position = position
	s.items[item.id] = item
	s.entities[item.id] = item

	// 添加到AOI
	s.aoi.AddEntity(item.id, item)

	s.lastUpdate = time.Now()

	// 发布事件
	s.addEvent(ItemDroppedEvent{
		sceneID:    s.id,
		itemID:     item.id,
		itemType:   item.itemID,
		position:   position,
		quantity:   item.quantity,
		occurredAt: time.Now(),
	})

	return nil
}

// GetNearbyEntities 获取附近实体
func (s *Scene) GetNearbyEntities(entityID string, radius float64) ([]Entity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entity, exists := s.entities[entityID]
	if !exists {
		return nil, ErrEntityNotFound
	}

	return s.aoi.GetNearbyEntities(entity.GetPosition(), radius), nil
}

// Update 更新场景
func (s *Scene) Update(deltaTime time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	// 更新所有实体
	for _, entity := range s.entities {
		if entity.IsActive() {
			entity.Update(deltaTime)
		}
	}

	// 处理刷新点
	s.processSpawnPoints(now)

	// 清理过期物品
	s.cleanupExpiredItems(now)

	s.lastUpdate = now
}

// processSpawnPoints 处理刷新点
func (s *Scene) processSpawnPoints(now time.Time) {
	for _, spawnPoint := range s.spawnPoints {
		if !spawnPoint.active {
			continue
		}

		if spawnPoint.currentCount >= spawnPoint.maxCount {
			continue
		}

		if now.Sub(spawnPoint.lastSpawn) < spawnPoint.interval {
			continue
		}

		// 这里应该根据刷新点类型创建对应实体
		// 简化实现，实际应该从配置或工厂创建
		switch spawnPoint.spawnType {
		case SpawnTypeMonster:
			// 创建怪物逻辑
		case SpawnTypeNPC:
			// 创建NPC逻辑
		case SpawnTypeItem:
			// 创建物品逻辑
		}
	}
}

// cleanupExpiredItems 清理过期物品
func (s *Scene) cleanupExpiredItems(now time.Time) {
	for itemID, item := range s.items {
		if item.expireTime != nil && now.After(*item.expireTime) {
			s.aoi.RemoveEntity(itemID)
			delete(s.items, itemID)
			delete(s.entities, itemID)
		}
	}
}

// isValidPosition 检查位置是否有效
func (s *Scene) isValidPosition(pos *Position) bool {
	return pos.X >= 0 && pos.X <= s.width && pos.Y >= 0 && pos.Y <= s.height
}

// addEvent 添加领域事件
func (s *Scene) addEvent(event DomainEvent) {
	s.events = append(s.events, event)
}

// GetEvents 获取领域事件
func (s *Scene) GetEvents() []DomainEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.events
}

// ClearEvents 清除领域事件
func (s *Scene) ClearEvents() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = make([]DomainEvent, 0)
}

// AOIManager 方法实现

// AddEntity 添加实体到AOI
func (aoi *AOIManager) AddEntity(entityID string, entity Entity) {
	aoi.mu.Lock()
	defer aoi.mu.Unlock()

	pos := entity.GetPosition()
	gridX, gridY := aoi.getGridCoords(pos.X, pos.Y)

	aoiEntity := &AOIEntity{
		id:       entityID,
		entity:   entity,
		gridX:    gridX,
		gridY:    gridY,
		watchers: make(map[string]bool),
	}

	aoi.entities[entityID] = aoiEntity

	// 添加到网格
	gridKey := aoi.getGridKey(gridX, gridY)
	if _, exists := aoi.grids[gridKey]; !exists {
		aoi.grids[gridKey] = &AOIGrid{
			x:        gridX,
			y:        gridY,
			entities: make(map[string]*AOIEntity),
		}
	}
	aoi.grids[gridKey].entities[entityID] = aoiEntity
}

// RemoveEntity 从AOI移除实体
func (aoi *AOIManager) RemoveEntity(entityID string) {
	aoi.mu.Lock()
	defer aoi.mu.Unlock()

	aoiEntity, exists := aoi.entities[entityID]
	if !exists {
		return
	}

	// 从网格移除
	gridKey := aoi.getGridKey(aoiEntity.gridX, aoiEntity.gridY)
	if grid, exists := aoi.grids[gridKey]; exists {
		delete(grid.entities, entityID)
		if len(grid.entities) == 0 {
			delete(aoi.grids, gridKey)
		}
	}

	delete(aoi.entities, entityID)
}

// UpdateEntity 更新实体位置
func (aoi *AOIManager) UpdateEntity(entityID string, newPosition *Position) {
	aoi.mu.Lock()
	defer aoi.mu.Unlock()

	aoiEntity, exists := aoi.entities[entityID]
	if !exists {
		return
	}

	newGridX, newGridY := aoi.getGridCoords(newPosition.X, newPosition.Y)

	// 如果网格没有变化，直接返回
	if newGridX == aoiEntity.gridX && newGridY == aoiEntity.gridY {
		return
	}

	// 从旧网格移除
	oldGridKey := aoi.getGridKey(aoiEntity.gridX, aoiEntity.gridY)
	if grid, exists := aoi.grids[oldGridKey]; exists {
		delete(grid.entities, entityID)
		if len(grid.entities) == 0 {
			delete(aoi.grids, oldGridKey)
		}
	}

	// 添加到新网格
	newGridKey := aoi.getGridKey(newGridX, newGridY)
	if _, exists := aoi.grids[newGridKey]; !exists {
		aoi.grids[newGridKey] = &AOIGrid{
			x:        newGridX,
			y:        newGridY,
			entities: make(map[string]*AOIEntity),
		}
	}
	aoi.grids[newGridKey].entities[entityID] = aoiEntity

	// 更新实体网格坐标
	aoiEntity.gridX = newGridX
	aoiEntity.gridY = newGridY
}

// GetNearbyEntities 获取附近实体
func (aoi *AOIManager) GetNearbyEntities(position *Position, radius float64) []Entity {
	aoi.mu.RLock()
	defer aoi.mu.RUnlock()

	var entities []Entity
	gridRadius := int(math.Ceil(radius / aoi.gridSize))
	centerGridX, centerGridY := aoi.getGridCoords(position.X, position.Y)

	// 遍历周围网格
	for x := centerGridX - gridRadius; x <= centerGridX+gridRadius; x++ {
		for y := centerGridY - gridRadius; y <= centerGridY+gridRadius; y++ {
			gridKey := aoi.getGridKey(x, y)
			if grid, exists := aoi.grids[gridKey]; exists {
				for _, aoiEntity := range grid.entities {
					entityPos := aoiEntity.entity.GetPosition()
					if position.Distance(entityPos) <= radius {
						entities = append(entities, aoiEntity.entity)
					}
				}
			}
		}
	}

	return entities
}

// getGridCoords 获取网格坐标
func (aoi *AOIManager) getGridCoords(x, y float64) (int, int) {
	gridX := int(x / aoi.gridSize)
	gridY := int(y / aoi.gridSize)
	return gridX, gridY
}

// getGridKey 获取网格键
func (aoi *AOIManager) getGridKey(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

// 实体接口实现

// Player 实现 Entity 接口
func (p *Player) GetID() string {
	return p.id
}

func (p *Player) GetPosition() *Position {
	return p.position
}

func (p *Player) SetPosition(pos *Position) {
	p.position = pos
}

func (p *Player) GetEntityType() EntityType {
	return EntityTypePlayer
}

func (p *Player) Update(deltaTime time.Duration) {
	// 玩家更新逻辑
	p.lastAction = time.Now()
}

func (p *Player) IsActive() bool {
	return p.active && p.status != PlayerStatusDead
}

// Monster 实现 Entity 接口
func (m *Monster) GetID() string {
	return m.id
}

func (m *Monster) GetPosition() *Position {
	return m.position
}

func (m *Monster) SetPosition(pos *Position) {
	m.position = pos
}

func (m *Monster) GetEntityType() EntityType {
	return EntityTypeMonster
}

func (m *Monster) Update(deltaTime time.Duration) {
	// 怪物AI更新逻辑
	if m.ai != nil {
		m.updateAI(deltaTime)
	}
}

func (m *Monster) IsActive() bool {
	return m.active && m.status != MonsterStatusDead
}

func (m *Monster) updateAI(deltaTime time.Duration) {
	// AI行为更新逻辑
	switch m.ai.behaviorType {
	case BehaviorTypePatrol:
		// 巡逻逻辑
	case BehaviorTypeAggressive:
		// 攻击逻辑
	case BehaviorTypeGuard:
		// 守卫逻辑
	}
}

// NPC 实现 Entity 接口
func (n *NPC) GetID() string {
	return n.id
}

func (n *NPC) GetPosition() *Position {
	return n.position
}

func (n *NPC) SetPosition(pos *Position) {
	n.position = pos
}

func (n *NPC) GetEntityType() EntityType {
	return EntityTypeNPC
}

func (n *NPC) Update(deltaTime time.Duration) {
	// NPC更新逻辑
	if n.ai != nil {
		n.updateAI(deltaTime)
	}
}

func (n *NPC) IsActive() bool {
	return n.active && n.status != NPCStatusDead
}

func (n *NPC) updateAI(deltaTime time.Duration) {
	// NPC AI逻辑
}

// Item 实现 Entity 接口
func (i *Item) GetID() string {
	return i.id
}

func (i *Item) GetPosition() *Position {
	return i.position
}

func (i *Item) SetPosition(pos *Position) {
	i.position = pos
}

func (i *Item) GetEntityType() EntityType {
	return EntityTypeItem
}

func (i *Item) Update(deltaTime time.Duration) {
	// 物品更新逻辑（检查过期等）
}

func (i *Item) IsActive() bool {
	return i.active
}

// Portal 实现 Entity 接口
func (p *Portal) GetID() string {
	return p.id
}

func (p *Portal) GetPosition() *Position {
	return p.position
}

func (p *Portal) SetPosition(pos *Position) {
	p.position = pos
}

func (p *Portal) GetEntityType() EntityType {
	return EntityTypePortal
}

func (p *Portal) Update(deltaTime time.Duration) {
	// 传送门更新逻辑
}

func (p *Portal) IsActive() bool {
	return p.active
}
