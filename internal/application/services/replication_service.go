// Package services 应用服务层
// ReplicationService 编排副本实例的创建、管理和销毁
package services

import (
	"context"
	"fmt"
	"time"

	"greatestworks/internal/domain/replication"
	"greatestworks/internal/infrastructure/logging"

	"github.com/google/uuid"
)

// ReplicationService 副本应用服务
type ReplicationService struct {
	replicationRepo replication.Repository
	eventBus        EventPublisher
	logger          logging.Logger
}

// NewReplicationService 创建副本应用服务
func NewReplicationService(
	replicationRepo replication.Repository,
	eventBus EventPublisher,
	logger logging.Logger,
) *ReplicationService {
	return &ReplicationService{
		replicationRepo: replicationRepo,
		eventBus:        eventBus,
		logger:          logger,
	}
}

// CreateInstanceCommand 创建实例命令
type CreateInstanceCommand struct {
	TemplateID    string
	InstanceType  int
	OwnerPlayerID string
	OwnerName     string
	OwnerLevel    int
	MaxPlayers    int
	Difficulty    int
	Lifetime      time.Duration
}

// JoinInstanceCommand 加入实例命令
type JoinInstanceCommand struct {
	InstanceID string
	PlayerID   string
	PlayerName string
	Level      int
	Role       string
}

// LeaveInstanceCommand 离开实例命令
type LeaveInstanceCommand struct {
	InstanceID string
	PlayerID   string
}

// UpdateInstanceProgressCommand 更新实例进度命令
type UpdateInstanceProgressCommand struct {
	InstanceID    string
	Progress      int
	CompletedTask string
}

// InstanceInfoDTO 实例信息DTO
type InstanceInfoDTO struct {
	InstanceID   string
	TemplateID   string
	InstanceType int
	Status       int
	PlayerCount  int
	MaxPlayers   int
	Progress     int
	SceneID      string
	CreatedAt    time.Time
	OwnerID      string
	Difficulty   int
}

// CreateInstance 创建副本实例
func (s *ReplicationService) CreateInstance(ctx context.Context, cmd *CreateInstanceCommand) (*InstanceInfoDTO, error) {
	// 基础校验
	if cmd == nil {
		return nil, fmt.Errorf("invalid command")
	}
	if cmd.TemplateID == "" || cmd.OwnerPlayerID == "" {
		return nil, fmt.Errorf("template_id and owner_player_id are required")
	}
	if cmd.MaxPlayers <= 0 {
		return nil, fmt.Errorf("max_players must be > 0")
	}
	if cmd.Lifetime <= 0 {
		return nil, fmt.Errorf("lifetime must be > 0")
	}
	s.logger.Info("创建副本实例", logging.Fields{
		"template_id": cmd.TemplateID,
		"owner_id":    cmd.OwnerPlayerID,
	})

	// 生成实例ID
	instanceID := uuid.New().String()

	// 创建领域对象
	instance := replication.NewInstance(
		instanceID,
		cmd.TemplateID,
		replication.InstanceType(cmd.InstanceType),
		cmd.OwnerPlayerID,
		cmd.MaxPlayers,
		cmd.Lifetime,
	)

	// 添加创建者作为第一个玩家
	if err := instance.AddPlayer(cmd.OwnerPlayerID, cmd.OwnerName, cmd.OwnerLevel, ""); err != nil {
		return nil, fmt.Errorf("添加创建者失败: %w", err)
	}

	// 保存到仓储
	if err := s.replicationRepo.Save(ctx, instance); err != nil {
		return nil, fmt.Errorf("保存实例失败: %w", err)
	}

	// 发布领域事件
	s.publishEvents(ctx, instance)

	s.logger.Info("副本实例创建成功", logging.Fields{
		"instance_id": instanceID,
		"template_id": cmd.TemplateID,
	})

	return &InstanceInfoDTO{
		InstanceID:   instance.ID(),
		TemplateID:   instance.TemplateID(),
		InstanceType: int(instance.Type()),
		Status:       int(instance.Status()),
		PlayerCount:  instance.PlayerCount(),
		MaxPlayers:   instance.MaxPlayers(),
		Progress:     instance.Progress(),
		SceneID:      instance.SceneID(),
		CreatedAt:    instance.CreatedAt(),
		Difficulty:   instance.Difficulty(),
	}, nil
}

// JoinInstance 加入副本实例
func (s *ReplicationService) JoinInstance(ctx context.Context, cmd *JoinInstanceCommand) error {
	if cmd == nil || cmd.InstanceID == "" || cmd.PlayerID == "" {
		return fmt.Errorf("instance_id and player_id are required")
	}
	s.logger.Info("玩家加入实例", logging.Fields{
		"instance_id": cmd.InstanceID,
		"player_id":   cmd.PlayerID,
	})

	// 查找实例
	instance, err := s.replicationRepo.FindByID(ctx, cmd.InstanceID)
	if err != nil {
		return fmt.Errorf("查找实例失败: %w", err)
	}
	if instance == nil {
		return fmt.Errorf("实例不存在: %s", cmd.InstanceID)
	}

	// 添加玩家
	if err := instance.AddPlayer(cmd.PlayerID, cmd.PlayerName, cmd.Level, cmd.Role); err != nil {
		return fmt.Errorf("添加玩家失败: %w", err)
	}

	// 保存
	if err := s.replicationRepo.Save(ctx, instance); err != nil {
		return fmt.Errorf("保存实例失败: %w", err)
	}

	// 发布事件
	s.publishEvents(ctx, instance)

	s.logger.Info("玩家加入实例成功", logging.Fields{
		"instance_id": cmd.InstanceID,
		"player_id":   cmd.PlayerID,
	})

	return nil
}

// LeaveInstance 离开副本实例
func (s *ReplicationService) LeaveInstance(ctx context.Context, cmd *LeaveInstanceCommand) error {
	if cmd == nil || cmd.InstanceID == "" || cmd.PlayerID == "" {
		return fmt.Errorf("instance_id and player_id are required")
	}
	s.logger.Info("玩家离开实例", logging.Fields{
		"instance_id": cmd.InstanceID,
		"player_id":   cmd.PlayerID,
	})

	// 查找实例
	instance, err := s.replicationRepo.FindByID(ctx, cmd.InstanceID)
	if err != nil {
		return fmt.Errorf("查找实例失败: %w", err)
	}
	if instance == nil {
		return fmt.Errorf("实例不存在: %s", cmd.InstanceID)
	}

	// 移除玩家
	if err := instance.RemovePlayer(cmd.PlayerID); err != nil {
		return fmt.Errorf("移除玩家失败: %w", err)
	}

	// 保存
	if err := s.replicationRepo.Save(ctx, instance); err != nil {
		return fmt.Errorf("保存实例失败: %w", err)
	}

	// 发布事件
	s.publishEvents(ctx, instance)

	s.logger.Info("玩家离开实例成功", logging.Fields{
		"instance_id": cmd.InstanceID,
		"player_id":   cmd.PlayerID,
	})

	return nil
}

// UpdateInstanceProgress 更新实例进度
func (s *ReplicationService) UpdateInstanceProgress(ctx context.Context, cmd *UpdateInstanceProgressCommand) error {
	if cmd == nil || cmd.InstanceID == "" {
		return fmt.Errorf("instance_id is required")
	}
	if cmd.Progress < 0 || cmd.Progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100")
	}
	s.logger.Info("更新实例进度", logging.Fields{
		"instance_id": cmd.InstanceID,
		"progress":    cmd.Progress,
	})

	// 查找实例
	instance, err := s.replicationRepo.FindByID(ctx, cmd.InstanceID)
	if err != nil {
		return fmt.Errorf("查找实例失败: %w", err)
	}
	if instance == nil {
		return fmt.Errorf("实例不存在: %s", cmd.InstanceID)
	}

	// 更新进度
	instance.UpdateProgress(cmd.Progress, cmd.CompletedTask)

	// 保存
	if err := s.replicationRepo.Save(ctx, instance); err != nil {
		return fmt.Errorf("保存实例失败: %w", err)
	}

	// 发布事件
	s.publishEvents(ctx, instance)

	return nil
}

// GetInstanceInfo 获取实例信息
func (s *ReplicationService) GetInstanceInfo(ctx context.Context, instanceID string) (*InstanceInfoDTO, error) {
	instance, err := s.replicationRepo.FindByID(ctx, instanceID)
	if err != nil {
		return nil, fmt.Errorf("查找实例失败: %w", err)
	}
	if instance == nil {
		return nil, fmt.Errorf("实例不存在: %s", instanceID)
	}

	return &InstanceInfoDTO{
		InstanceID:   instance.ID(),
		TemplateID:   instance.TemplateID(),
		InstanceType: int(instance.Type()),
		Status:       int(instance.Status()),
		PlayerCount:  instance.PlayerCount(),
		MaxPlayers:   instance.MaxPlayers(),
		Progress:     instance.Progress(),
		SceneID:      instance.SceneID(),
		CreatedAt:    instance.CreatedAt(),
		Difficulty:   instance.Difficulty(),
	}, nil
}

// ListActiveInstances 列出所有活跃实例
func (s *ReplicationService) ListActiveInstances(ctx context.Context) ([]*InstanceInfoDTO, error) {
	instances, err := s.replicationRepo.FindActiveInstances(ctx)
	if err != nil {
		return nil, fmt.Errorf("查找活跃实例失败: %w", err)
	}

	result := make([]*InstanceInfoDTO, 0, len(instances))
	for _, instance := range instances {
		result = append(result, &InstanceInfoDTO{
			InstanceID:   instance.ID(),
			TemplateID:   instance.TemplateID(),
			InstanceType: int(instance.Type()),
			Status:       int(instance.Status()),
			PlayerCount:  instance.PlayerCount(),
			MaxPlayers:   instance.MaxPlayers(),
			Progress:     instance.Progress(),
			SceneID:      instance.SceneID(),
			CreatedAt:    instance.CreatedAt(),
			Difficulty:   instance.Difficulty(),
		})
	}

	return result, nil
}

// CleanupExpiredInstances 清理过期实例
func (s *ReplicationService) CleanupExpiredInstances(ctx context.Context) (int, error) {
	s.logger.Info("清理过期实例")

	instances, err := s.replicationRepo.FindExpiredInstances(ctx)
	if err != nil {
		return 0, fmt.Errorf("查找过期实例失败: %w", err)
	}

	count := 0
	for _, instance := range instances {
		instance.Close()
		if err := s.replicationRepo.Save(ctx, instance); err != nil {
			s.logger.Error("保存关闭实例失败", err, logging.Fields{
				"instance_id": instance.ID(),
			})
			continue
		}

		// 发布事件
		s.publishEvents(ctx, instance)
		count++
	}

	s.logger.Info("过期实例清理完成", logging.Fields{
		"count": count,
	})

	return count, nil
}

// publishEvents 发布领域事件
func (s *ReplicationService) publishEvents(ctx context.Context, instance *replication.Instance) {
	events := instance.GetEvents()
	for _, event := range events {
		if err := s.eventBus.Publish(ctx, event); err != nil {
			s.logger.Error("发布事件失败", err, logging.Fields{
				"event_type": fmt.Sprintf("%T", event),
			})
		}
	}
}
