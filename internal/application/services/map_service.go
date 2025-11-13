package services

import (
	"context"
	"errors"
	"sync"
	"time"

	"greatestworks/internal/domain/character"
	"greatestworks/internal/domain/mapmanager"
	"greatestworks/internal/infrastructure/datamanager"
)

// MapService 地图服务
type MapService struct {
	mu   sync.RWMutex
	maps map[int32]*mapmanager.Map
	// 应用层广播适配器
	broadcaster mapmanager.BroadcastFn
	// 异步刷怪/掉落等任务
	spawnMgr *SpawnManager
}

// NewMapService 创建地图服务
func NewMapService() *MapService {
	return &MapService{
		maps: make(map[int32]*mapmanager.Map),
	}
}

// SetBroadcaster 设置地图广播函数（通常由接口层注入，用于向会话发送消息）
func (s *MapService) SetBroadcaster(fn mapmanager.BroadcastFn) {
	s.mu.Lock()
	s.broadcaster = fn
	// 将已加载地图也设置广播器
	for _, m := range s.maps {
		m.SetBroadcaster(fn)
	}
	s.mu.Unlock()
}

// SetSpawnManager 注入SpawnManager以在地图生命周期中投递异步任务
func (s *MapService) SetSpawnManager(sm *SpawnManager) {
	s.mu.Lock()
	s.spawnMgr = sm
	s.mu.Unlock()
}

// LoadMap 加载地图
func (s *MapService) LoadMap(ctx context.Context, mapID int32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查是否已加载
	if _, exists := s.maps[mapID]; exists {
		return nil
	}

	// 获取地图配置
	mapDefine := datamanager.GetInstance().GetMap(mapID)
	if mapDefine == nil {
		return errors.New("map not found")
	}

	// 创建地图实例
	gameMap := mapmanager.NewMap(mapID, mapDefine.Name, mapDefine.Width, mapDefine.Height)
	if s.broadcaster != nil {
		gameMap.SetBroadcaster(s.broadcaster)
	}
	s.maps[mapID] = gameMap

	return nil
}

// GetMap 获取地图
func (s *MapService) GetMap(mapID int32) (*mapmanager.Map, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	gameMap, exists := s.maps[mapID]
	if !exists {
		return nil, errors.New("map not loaded")
	}

	return gameMap, nil
}

// EnterMap 进入地图
func (s *MapService) EnterMap(ctx context.Context, entity *character.Entity, mapID int32, x, y, z float32) error {
	gameMap, err := s.GetMap(mapID)
	if err != nil {
		// 尝试加载地图
		if err := s.LoadMap(ctx, mapID); err != nil {
			return err
		}
		gameMap, err = s.GetMap(mapID)
		if err != nil {
			return err
		}
	}

	// 进入地图
	// 设置初始位置
	entity.SetPosition(character.NewVector3(x, y, z))
	_ = gameMap.Enter(ctx, entity)

	// 可选：在进入地图时投递一次异步任务（示例）
	if s.spawnMgr != nil {
		s.spawnMgr.Enqueue(func(ctx context.Context) {
			// 占位：后续可在此触发进入地图后的刷怪或欢迎事件
		})
	}
	return nil
}

// LeaveMap 离开地图
func (s *MapService) LeaveMap(ctx context.Context, entity *character.Entity, mapID int32) error {
	gameMap, err := s.GetMap(mapID)
	if err != nil {
		return err
	}

	_ = gameMap.Leave(ctx, entity.ID())
	return nil
}

// LeaveMapByID 使用地图ID和实体ID离开地图（便于接口/清理逻辑调用）
func (s *MapService) LeaveMapByID(ctx context.Context, mapID int32, entityID int32) error {
	gameMap, err := s.GetMap(mapID)
	if err != nil {
		return err
	}
	return gameMap.Leave(ctx, character.EntityID(entityID))
}

// UpdatePosition 更新位置
func (s *MapService) UpdatePosition(ctx context.Context, entity *character.Entity, mapID int32, x, y, z float32) error {
	gameMap, err := s.GetMap(mapID)
	if err != nil {
		return err
	}

	return gameMap.UpdatePosition(entity.ID(), character.NewVector3(x, y, z))
}

// UpdatePositionByID 使用地图ID和实体ID更新位置（便于接口层调用）
func (s *MapService) UpdatePositionByID(ctx context.Context, mapID int32, entityID int32, x, y, z float32) error {
	gameMap, err := s.GetMap(mapID)
	if err != nil {
		return err
	}
	return gameMap.UpdatePosition(character.EntityID(entityID), character.NewVector3(x, y, z))
}

// GetEntitiesInRange 获取范围内的实体
func (s *MapService) GetEntitiesInRange(ctx context.Context, mapID int32, x, y, z, range_ float32) ([]*character.Entity, error) {
	gameMap, err := s.GetMap(mapID)
	if err != nil {
		return nil, err
	}

	return gameMap.GetEntitiesInRange(x, z, range_), nil
}

// TransferMap 传送到其他地图
func (s *MapService) TransferMap(ctx context.Context, entity *character.Entity, fromMapID, toMapID int32, x, y, z float32) error {
	// 离开当前地图
	if err := s.LeaveMap(ctx, entity, fromMapID); err != nil {
		return err
	}

	// 进入目标地图
	if err := s.EnterMap(ctx, entity, toMapID, x, y, z); err != nil {
		// 回滚：重新进入原地图
		_ = s.EnterMap(ctx, entity, fromMapID, 0, 0, 0)
		return err
	}

	return nil
}

// BroadcastToMap 向地图广播消息
func (s *MapService) BroadcastToMap(ctx context.Context, mapID int32, message interface{}) error {
	gameMap, err := s.GetMap(mapID)
	if err != nil {
		return err
	}

	// 获取地图内所有实体并广播
	entities := gameMap.GetAllEntities()
	recipients := make([]character.EntityID, 0, len(entities))
	for _, e := range entities {
		recipients = append(recipients, e.ID())
	}
	gameMap.BroadcastTo(recipients, "map_broadcast", message)

	return nil
}

// BroadcastToRange 向范围内广播消息
func (s *MapService) BroadcastToRange(ctx context.Context, mapID int32, x, y, z, range_ float32, message interface{}) error {
	gameMap, err := s.GetMap(mapID)
	if err != nil {
		return err
	}
	// 使用AOI广播到范围内
	gameMap.BroadcastInRange(x, z, range_, "range_broadcast", message)

	return nil
}

// Tick 地图更新（供 UpdateManager 调用）
func (s *MapService) Tick(ctx context.Context, delta time.Duration) {
	// 目前无需要逐帧更新的逻辑；保留入口以便未来接入地图定时器/效果
	// 示例：遍历地图执行周期性任务
	s.mu.RLock()
	defer s.mu.RUnlock()
	_ = delta
	_ = ctx
	// for _, m := range s.maps { /* m.DoPeriodicStuff() */ }
}
