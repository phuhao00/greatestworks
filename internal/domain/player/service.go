package player

import (
	"context"
	"fmt"
	"math"
	"time"
)

// Service 玩家领域服务
type Service struct {
	repository Repository
}

// NewService 创建玩家领域服务
func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

// CreatePlayer 创建玩家
func (s *Service) CreatePlayer(ctx context.Context, name string) (*Player, error) {
	if name == "" {
		return nil, ErrInvalidPlayerName
	}
	
	// 检查名称是否已存在
	exists, err := s.repository.ExistsByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("check player name exists: %w", err)
	}
	if exists {
		return nil, ErrPlayerAlreadyExists
	}
	
	// 创建新玩家
	player := NewPlayer(name)
	
	// 保存到仓储
	if err := s.repository.Save(ctx, player); err != nil {
		return nil, fmt.Errorf("save player: %w", err)
	}
	
	return player, nil
}

// AuthenticatePlayer 玩家认证
func (s *Service) AuthenticatePlayer(ctx context.Context, playerID PlayerID) (*Player, error) {
	player, err := s.repository.FindByID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("find player: %w", err)
	}
	
	// 设置玩家上线
	player.SetOnline()
	
	// 更新玩家状态
	if err := s.repository.Update(ctx, player); err != nil {
		return nil, fmt.Errorf("update player: %w", err)
	}
	
	return player, nil
}

// LogoutPlayer 玩家登出
func (s *Service) LogoutPlayer(ctx context.Context, playerID PlayerID) error {
	player, err := s.repository.FindByID(ctx, playerID)
	if err != nil {
		return fmt.Errorf("find player: %w", err)
	}
	
	// 设置玩家下线
	player.SetOffline()
	
	// 更新玩家状态
	if err := s.repository.Update(ctx, player); err != nil {
		return fmt.Errorf("update player: %w", err)
	}
	
	return nil
}

// MovePlayer 移动玩家
func (s *Service) MovePlayer(ctx context.Context, playerID PlayerID, position Position) error {
	player, err := s.repository.FindByID(ctx, playerID)
	if err != nil {
		return fmt.Errorf("find player: %w", err)
	}
	
	// 验证位置有效性
	if err := s.validatePosition(position); err != nil {
		return err
	}
	
	// 移动玩家
	if err := player.MoveTo(position); err != nil {
		return err
	}
	
	// 更新玩家
	if err := s.repository.Update(ctx, player); err != nil {
		return fmt.Errorf("update player: %w", err)
	}
	
	return nil
}

// GainExperience 玩家获得经验
func (s *Service) GainExperience(ctx context.Context, playerID PlayerID, exp int64) error {
	player, err := s.repository.FindByID(ctx, playerID)
	if err != nil {
		return fmt.Errorf("find player: %w", err)
	}
	
	oldLevel := player.Level()
	player.GainExp(exp)
	newLevel := player.Level()
	
	// 更新玩家
	if err := s.repository.Update(ctx, player); err != nil {
		return fmt.Errorf("update player: %w", err)
	}
	
	// 如果升级了，发布升级事件
	if newLevel > oldLevel {
		// TODO: 发布玩家升级事件
	}
	
	return nil
}

// HealPlayer 治疗玩家
func (s *Service) HealPlayer(ctx context.Context, playerID PlayerID, amount int) error {
	player, err := s.repository.FindByID(ctx, playerID)
	if err != nil {
		return fmt.Errorf("find player: %w", err)
	}
	
	player.Heal(amount)
	
	// 更新玩家
	if err := s.repository.Update(ctx, player); err != nil {
		return fmt.Errorf("update player: %w", err)
	}
	
	return nil
}

// GetOnlinePlayers 获取在线玩家列表
func (s *Service) GetOnlinePlayers(ctx context.Context, limit int) ([]*Player, error) {
	players, err := s.repository.FindOnlinePlayers(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("find online players: %w", err)
	}
	
	return players, nil
}

// validatePosition 验证位置有效性
func (s *Service) validatePosition(pos Position) error {
	// 简单的位置验证逻辑
	if pos.X < -1000 || pos.X > 1000 {
		return ErrInvalidPosition
	}
	if pos.Y < -1000 || pos.Y > 1000 {
		return ErrInvalidPosition
	}
	if pos.Z < -100 || pos.Z > 100 {
		return ErrInvalidPosition
	}
	
	return nil
}

// CalculateDistance 计算两个位置之间的距离
func (s *Service) CalculateDistance(pos1, pos2 Position) float64 {
	dx := pos1.X - pos2.X
	dy := pos1.Y - pos2.Y
	dz := pos1.Z - pos2.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// IsPlayerNearby 检查玩家是否在附近
func (s *Service) IsPlayerNearby(ctx context.Context, playerID1, playerID2 PlayerID, maxDistance float64) (bool, error) {
	player1, err := s.repository.FindByID(ctx, playerID1)
	if err != nil {
		return false, fmt.Errorf("find player1: %w", err)
	}
	
	player2, err := s.repository.FindByID(ctx, playerID2)
	if err != nil {
		return false, fmt.Errorf("find player2: %w", err)
	}
	
	distance := s.CalculateDistance(player1.GetPosition(), player2.GetPosition())
	return distance <= maxDistance, nil
}