// Package services 应用服务层 - 场景服务
package services

import (
	"context"
	"fmt"
	"time"

	"greatestworks/internal/domain/scene"
	"greatestworks/internal/infrastructure/logging"
)

// SceneService 场景应用服务
// 编排场景聚合根，协调仓储、事件总线等基础设施
type SceneService struct {
	sceneRepo scene.Repository
	eventBus  EventPublisher
	logger    logging.Logger
}

// EventPublisher 事件发布器接口
type EventPublisher interface {
	Publish(ctx context.Context, event interface{}) error
}

// NewSceneService 创建场景应用服务
func NewSceneService(
	sceneRepo scene.Repository,
	eventBus EventPublisher,
	logger logging.Logger,
) *SceneService {
	return &SceneService{
		sceneRepo: sceneRepo,
		eventBus:  eventBus,
		logger:    logger,
	}
}

// CreateSceneCommand 创建场景命令
type CreateSceneCommand struct {
	ID         string          `json:"id" validate:"required"`
	Name       string          `json:"name" validate:"required,min=2,max=50"`
	SceneType  scene.SceneType `json:"scene_type" validate:"required"`
	Width      float64         `json:"width" validate:"required,gt=0"`
	Height     float64         `json:"height" validate:"required,gt=0"`
	MaxPlayers int             `json:"max_players" validate:"required,gt=0"`
}

// CreateSceneResult 创建场景结果
type CreateSceneResult struct {
	SceneID   string            `json:"scene_id"`
	Name      string            `json:"name"`
	SceneType scene.SceneType   `json:"scene_type"`
	Status    scene.SceneStatus `json:"status"`
}

// CreateScene 创建场景
func (s *SceneService) CreateScene(ctx context.Context, cmd *CreateSceneCommand) (*CreateSceneResult, error) {
	// 检查场景是否已存在
	exists, err := s.sceneRepo.Exists(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("检查场景是否存在失败: %w", err)
	}
	if exists {
		return nil, scene.ErrSceneAlreadyExists
	}

	// 创建场景聚合根
	newScene := scene.NewScene(
		cmd.ID,
		cmd.Name,
		cmd.SceneType,
		cmd.Width,
		cmd.Height,
		cmd.MaxPlayers,
	)

	// 保存场景
	if err := s.sceneRepo.Save(ctx, newScene); err != nil {
		return nil, fmt.Errorf("保存场景失败: %w", err)
	}

	s.logger.Info("创建场景成功", logging.Fields{
		"scene_id":   newScene.ID(),
		"scene_name": newScene.Name(),
		"scene_type": cmd.SceneType,
	})

	return &CreateSceneResult{
		SceneID:   newScene.ID(),
		Name:      newScene.Name(),
		SceneType: cmd.SceneType,
		Status:    newScene.Status(),
	}, nil
}

// PlayerEnterSceneCommand 玩家进入场景命令
type PlayerEnterSceneCommand struct {
	SceneID    string          `json:"scene_id" validate:"required"`
	PlayerID   string          `json:"player_id" validate:"required"`
	PlayerName string          `json:"player_name" validate:"required"`
	Level      int             `json:"level" validate:"required,gt=0"`
	Position   *scene.Position `json:"position"`
}

// PlayerEnterScene 玩家进入场景
func (s *SceneService) PlayerEnterScene(ctx context.Context, cmd *PlayerEnterSceneCommand) error {
	// 获取场景
	sceneObj, err := s.sceneRepo.FindByID(ctx, cmd.SceneID)
	if err != nil {
		return fmt.Errorf("获取场景失败: %w", err)
	}
	if sceneObj == nil {
		return scene.ErrSceneNotFound
	}

	// 创建玩家实体（这里简化，实际应从玩家服务获取）
	player := &scene.Player{
		// 填充玩家数据
	}

	// 添加玩家到场景
	if err := sceneObj.AddPlayer(player); err != nil {
		return fmt.Errorf("玩家进入场景失败: %w", err)
	}

	// 保存场景状态
	if err := s.sceneRepo.Save(ctx, sceneObj); err != nil {
		return fmt.Errorf("保存场景状态失败: %w", err)
	}

	// 发布领域事件
	events := sceneObj.GetEvents()
	for _, event := range events {
		if err := s.eventBus.Publish(ctx, event); err != nil {
			s.logger.Error("发布事件失败", err, logging.Fields{
				"event_type": event.EventType(),
				"scene_id":   event.SceneID(),
			})
		}
	}
	sceneObj.ClearEvents()

	s.logger.Info("玩家进入场景", logging.Fields{
		"scene_id":    cmd.SceneID,
		"player_id":   cmd.PlayerID,
		"player_name": cmd.PlayerName,
	})

	return nil
}

// PlayerLeaveSceneCommand 玩家离开场景命令
type PlayerLeaveSceneCommand struct {
	SceneID  string `json:"scene_id" validate:"required"`
	PlayerID string `json:"player_id" validate:"required"`
}

// PlayerLeaveScene 玩家离开场景
func (s *SceneService) PlayerLeaveScene(ctx context.Context, cmd *PlayerLeaveSceneCommand) error {
	// 获取场景
	sceneObj, err := s.sceneRepo.FindByID(ctx, cmd.SceneID)
	if err != nil {
		return fmt.Errorf("获取场景失败: %w", err)
	}
	if sceneObj == nil {
		return scene.ErrSceneNotFound
	}

	// 移除玩家
	if err := sceneObj.RemovePlayer(cmd.PlayerID); err != nil {
		return fmt.Errorf("玩家离开场景失败: %w", err)
	}

	// 保存场景状态
	if err := s.sceneRepo.Save(ctx, sceneObj); err != nil {
		return fmt.Errorf("保存场景状态失败: %w", err)
	}

	// 发布领域事件
	events := sceneObj.GetEvents()
	for _, event := range events {
		if err := s.eventBus.Publish(ctx, event); err != nil {
			s.logger.Error("发布事件失败", err, logging.Fields{
				"event_type": event.EventType(),
				"scene_id":   event.SceneID(),
			})
		}
	}
	sceneObj.ClearEvents()

	s.logger.Info("玩家离开场景", logging.Fields{
		"scene_id":  cmd.SceneID,
		"player_id": cmd.PlayerID,
	})

	return nil
}

// UpdateSceneCommand 更新场景命令
type UpdateSceneCommand struct {
	SceneID   string `json:"scene_id" validate:"required"`
	DeltaTime int64  `json:"delta_time" validate:"required,gt=0"` // 毫秒
}

// UpdateScene 更新场景（tick）
func (s *SceneService) UpdateScene(ctx context.Context, cmd *UpdateSceneCommand) error {
	// 获取场景
	sceneObj, err := s.sceneRepo.FindByID(ctx, cmd.SceneID)
	if err != nil {
		return fmt.Errorf("获取场景失败: %w", err)
	}
	if sceneObj == nil {
		return scene.ErrSceneNotFound
	}

	// 更新场景逻辑（AI、刷怪、AOI等）
	deltaTime := time.Duration(cmd.DeltaTime) * time.Millisecond
	sceneObj.Update(deltaTime)

	// 保存场景状态
	if err := s.sceneRepo.Save(ctx, sceneObj); err != nil {
		return fmt.Errorf("保存场景状态失败: %w", err)
	}

	// 发布领域事件
	events := sceneObj.GetEvents()
	for _, event := range events {
		if err := s.eventBus.Publish(ctx, event); err != nil {
			s.logger.Error("发布事件失败", err, logging.Fields{
				"event_type": event.EventType(),
				"scene_id":   event.SceneID(),
			})
		}
	}
	sceneObj.ClearEvents()

	return nil
}

// GetSceneInfoQuery 获取场景信息查询
type GetSceneInfoQuery struct {
	SceneID string `json:"scene_id" validate:"required"`
}

// SceneInfoDTO 场景信息DTO
type SceneInfoDTO struct {
	SceneID        string            `json:"scene_id"`
	Name           string            `json:"name"`
	SceneType      scene.SceneType   `json:"scene_type"`
	Status         scene.SceneStatus `json:"status"`
	CurrentPlayers int               `json:"current_players"`
	MaxPlayers     int               `json:"max_players"`
}

// GetSceneInfo 获取场景信息
func (s *SceneService) GetSceneInfo(ctx context.Context, query *GetSceneInfoQuery) (*SceneInfoDTO, error) {
	// 获取场景
	sceneObj, err := s.sceneRepo.FindByID(ctx, query.SceneID)
	if err != nil {
		return nil, fmt.Errorf("获取场景失败: %w", err)
	}
	if sceneObj == nil {
		return nil, scene.ErrSceneNotFound
	}

	return &SceneInfoDTO{
		SceneID:        sceneObj.ID(),
		Name:           sceneObj.Name(),
		SceneType:      sceneObj.Type(),
		Status:         sceneObj.Status(),
		CurrentPlayers: sceneObj.PlayerCount(),
		MaxPlayers:     sceneObj.GetMaxPlayers(),
	}, nil
}

// ListAvailableScenesQuery 列出可用场景查询
type ListAvailableScenesQuery struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// ListAvailableScenes 列出可用场景
func (s *SceneService) ListAvailableScenes(ctx context.Context, query *ListAvailableScenesQuery) ([]*SceneInfoDTO, error) {
	// 获取可用场景列表
	scenes, err := s.sceneRepo.FindAvailableScenes(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取可用场景列表失败: %w", err)
	}

	// 转换为DTO
	result := make([]*SceneInfoDTO, 0, len(scenes))
	for _, sceneObj := range scenes {
		result = append(result, &SceneInfoDTO{
			SceneID:        sceneObj.ID(),
			Name:           sceneObj.Name(),
			SceneType:      sceneObj.Type(),
			Status:         sceneObj.Status(),
			CurrentPlayers: sceneObj.PlayerCount(),
			MaxPlayers:     sceneObj.GetMaxPlayers(),
		})
	}

	return result, nil
}
