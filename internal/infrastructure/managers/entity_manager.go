package managers

import (
	"context"
	"fmt"
	"greatestworks/internal/domain/character"
	"sync"
	"sync/atomic"
)

// EntityManager 实体管理器
type EntityManager struct {
	mu sync.RWMutex

	entities     map[character.EntityID]*character.Entity
	nextEntityID atomic.Int64
}

var entityManagerInstance *EntityManager
var entityManagerOnce sync.Once

// GetEntityManager 获取实体管理器单例
func GetEntityManager() *EntityManager {
	entityManagerOnce.Do(func() {
		entityManagerInstance = &EntityManager{
			entities: make(map[character.EntityID]*character.Entity),
		}
		entityManagerInstance.nextEntityID.Store(1000)
	})
	return entityManagerInstance
}

// Register 注册实体
func (em *EntityManager) Register(entity *character.Entity) error {
	if entity == nil {
		return fmt.Errorf("entity is nil")
	}

	em.mu.Lock()
	defer em.mu.Unlock()

	entityID := entity.ID()
	if _, exists := em.entities[entityID]; exists {
		return fmt.Errorf("entity already registered: %d", entityID)
	}

	em.entities[entityID] = entity
	return nil
}

// Unregister 注销实体
func (em *EntityManager) Unregister(entityID character.EntityID) {
	em.mu.Lock()
	defer em.mu.Unlock()

	delete(em.entities, entityID)
}

// Get 获取实体
func (em *EntityManager) Get(entityID character.EntityID) *character.Entity {
	em.mu.RLock()
	defer em.mu.RUnlock()

	return em.entities[entityID]
}

// GetAll 获取所有实体
func (em *EntityManager) GetAll() []*character.Entity {
	em.mu.RLock()
	defer em.mu.RUnlock()

	entities := make([]*character.Entity, 0, len(em.entities))
	for _, entity := range em.entities {
		entities = append(entities, entity)
	}
	return entities
}

// AllocateEntityID 分配实体ID
func (em *EntityManager) AllocateEntityID() character.EntityID {
	return character.EntityID(em.nextEntityID.Add(1))
}

// Count 获取实体数量
func (em *EntityManager) Count() int {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return len(em.entities)
}

// SpawnManager 生成管理器
type SpawnManager struct {
	mu sync.RWMutex

	spawnPoints map[int32]*SpawnPoint // 刷新点ID -> 刷新点
}

// SpawnPoint 刷新点
type SpawnPoint struct {
	ID           int32 // 刷新点ID
	MapID        int32 // 地图ID
	UnitDefineID int32 // 单位定义ID
	Position     character.Vector3
	RespawnTime  float32 // 重生时间
	MaxCount     int32   // 最大数量

	currentCount    int32                // 当前数量
	respawnTimer    float32              // 重生计时器
	spawnedEntities []character.EntityID // 已生成的实体ID
}

var spawnManagerInstance *SpawnManager
var spawnManagerOnce sync.Once

// GetSpawnManager 获取生成管理器单例
func GetSpawnManager() *SpawnManager {
	spawnManagerOnce.Do(func() {
		spawnManagerInstance = &SpawnManager{
			spawnPoints: make(map[int32]*SpawnPoint),
		}
	})
	return spawnManagerInstance
}

// AddSpawnPoint 添加刷新点
func (sm *SpawnManager) AddSpawnPoint(sp *SpawnPoint) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sp.spawnedEntities = make([]character.EntityID, 0)
	sm.spawnPoints[sp.ID] = sp
}

// Update 更新刷新点
func (sm *SpawnManager) Update(ctx context.Context, deltaTime float32) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for _, sp := range sm.spawnPoints {
		if sp.currentCount < sp.MaxCount {
			sp.respawnTimer += deltaTime
			if sp.respawnTimer >= sp.RespawnTime {
				sm.spawnEntity(ctx, sp)
				sp.respawnTimer = 0
			}
		}
	}
}

// spawnEntity 生成实体
func (sm *SpawnManager) spawnEntity(ctx context.Context, sp *SpawnPoint) {
	// TODO: 根据UnitDefineID创建实体
	// 这里需要EntityFactory
	sp.currentCount++
}

// OnEntityDestroyed 实体销毁回调
func (sm *SpawnManager) OnEntityDestroyed(entityID character.EntityID) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// 查找并更新对应的刷新点
	for _, sp := range sm.spawnPoints {
		for i, id := range sp.spawnedEntities {
			if id == entityID {
				sp.spawnedEntities = append(sp.spawnedEntities[:i], sp.spawnedEntities[i+1:]...)
				sp.currentCount--
				return
			}
		}
	}
}
