// Package services 应用服务层
package services

import (
	"context"
	"fmt"
	"log"

	"greatestworks/internal/domain/player"
)

// PlayerService 玩家应用服务
type PlayerService struct {
	playerRepo player.Repository
}

// NewPlayerService 创建新的玩家应用服务
func NewPlayerService(playerRepo player.Repository) *PlayerService {
	return &PlayerService{
		playerRepo: playerRepo,
	}
}

// CreatePlayerCommand 创建玩家命令
type CreatePlayerCommand struct {
	Name string `json:"name" validate:"required,min=2,max=20"`
}

// CreatePlayerResult 创建玩家结果
type CreatePlayerResult struct {
	PlayerID string `json:"player_id"`
	Name     string `json:"name"`
	Level    int    `json:"level"`
}

// CreatePlayer 创建玩家
func (s *PlayerService) CreatePlayer(ctx context.Context, cmd *CreatePlayerCommand) (*CreatePlayerResult, error) {
	// 验证玩家名称是否已存在
	exists, err := s.playerRepo.ExistsByName(ctx, cmd.Name)
	if err != nil {
		return nil, fmt.Errorf("检查玩家名称失败: %w", err)
	}
	if exists {
		return nil, player.ErrPlayerAlreadyExists
	}

	// 创建新玩家
	newPlayer := player.NewPlayer(cmd.Name)

	// 保存玩家
	if err := s.playerRepo.Save(ctx, newPlayer); err != nil {
		return nil, fmt.Errorf("保存玩家失败: %w", err)
	}

	log.Printf("创建玩家成功: %s (ID: %s)", newPlayer.Name(), newPlayer.ID().String())

	return &CreatePlayerResult{
		PlayerID: newPlayer.ID().String(),
		Name:     newPlayer.Name(),
		Level:    newPlayer.Level(),
	}, nil
}

// MovePlayer 移动玩家
func (s *PlayerService) MovePlayer(ctx context.Context, playerID string, position player.Position) error {
	// 获取玩家
	player, err := s.playerRepo.FindByID(ctx, playerID)
	if err != nil {
		return fmt.Errorf("获取玩家失败: %w", err)
	}
	if player == nil {
		return player.ErrPlayerNotFound
	}

	// 更新玩家位置
	player.SetPosition(position)

	// 保存玩家
	if err := s.playerRepo.Save(ctx, player); err != nil {
		return fmt.Errorf("保存玩家失败: %w", err)
	}

	return nil
}

// LoginPlayerCommand 玩家登录命令
type LoginPlayerCommand struct {
	PlayerID string `json:"player_id" validate:"required"`
}

// LoginPlayerResult 玩家登录结果
type LoginPlayerResult struct {
	PlayerID string              `json:"player_id"`
	Name     string              `json:"name"`
	Level    int                 `json:"level"`
	Status   player.PlayerStatus `json:"status"`
	Position player.Position     `json:"position"`
	Stats    player.PlayerStats  `json:"stats"`
}

// LoginPlayer 玩家登录
func (s *PlayerService) LoginPlayer(ctx context.Context, cmd *LoginPlayerCommand) (*LoginPlayerResult, error) {
	// 解析玩家ID
	playerID := player.PlayerID{}
	// 注意：这里需要实现PlayerID的解析逻辑

	// 查找玩家
	p, err := s.playerRepo.FindByID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("查找玩家失败: %w", err)
	}

	// 设置玩家上线
	p.SetOnline()

	// 更新玩家状态
	if err := s.playerRepo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("更新玩家状态失败: %w", err)
	}

	log.Printf("玩家登录成功: %s (ID: %s)", p.Name(), p.ID().String())

	return &LoginPlayerResult{
		PlayerID: p.ID().String(),
		Name:     p.Name(),
		Level:    p.Level(),
		Status:   p.Status(),
		Position: p.GetPosition(),
		Stats:    p.Stats(),
	}, nil
}

// LogoutPlayerCommand 玩家登出命令
type LogoutPlayerCommand struct {
	PlayerID string `json:"player_id" validate:"required"`
}

// LogoutPlayer 玩家登出
func (s *PlayerService) LogoutPlayer(ctx context.Context, cmd *LogoutPlayerCommand) error {
	// 解析玩家ID
	playerID := player.PlayerID{}
	// 注意：这里需要实现PlayerID的解析逻辑

	// 查找玩家
	p, err := s.playerRepo.FindByID(ctx, playerID)
	if err != nil {
		return fmt.Errorf("查找玩家失败: %w", err)
	}

	// 设置玩家下线
	p.SetOffline()

	// 更新玩家状态
	if err := s.playerRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("更新玩家状态失败: %w", err)
	}

	log.Printf("玩家登出成功: %s (ID: %s)", p.Name(), p.ID().String())
	return nil
}

// MovePlayerCommand 玩家移动命令
type MovePlayerCommand struct {
	PlayerID string          `json:"player_id" validate:"required"`
	Position player.Position `json:"position" validate:"required"`
}

// MovePlayer 玩家移动
func (s *PlayerService) MovePlayer(ctx context.Context, cmd *MovePlayerCommand) error {
	// 解析玩家ID
	playerID := player.PlayerID{}
	// 注意：这里需要实现PlayerID的解析逻辑

	// 查找玩家
	p, err := s.playerRepo.FindByID(ctx, playerID)
	if err != nil {
		return fmt.Errorf("查找玩家失败: %w", err)
	}

	// 移动玩家
	if err := p.MoveTo(cmd.Position); err != nil {
		return fmt.Errorf("移动玩家失败: %w", err)
	}

	// 更新玩家位置
	if err := s.playerRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("更新玩家位置失败: %w", err)
	}

	return nil
}

// GetPlayerQuery 获取玩家查询
type GetPlayerQuery struct {
	PlayerID string `json:"player_id" validate:"required"`
}

// GetPlayerResult 获取玩家结果
type GetPlayerResult struct {
	PlayerID string              `json:"player_id"`
	Name     string              `json:"name"`
	Level    int                 `json:"level"`
	Exp      int64               `json:"exp"`
	Status   player.PlayerStatus `json:"status"`
	Position player.Position     `json:"position"`
	Stats    player.PlayerStats  `json:"stats"`
}

// GetPlayer 获取玩家信息
func (s *PlayerService) GetPlayer(ctx context.Context, query *GetPlayerQuery) (*GetPlayerResult, error) {
	// 解析玩家ID
	playerID := player.PlayerID{}
	// 注意：这里需要实现PlayerID的解析逻辑

	// 查找玩家
	p, err := s.playerRepo.FindByID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("查找玩家失败: %w", err)
	}

	return &GetPlayerResult{
		PlayerID: p.ID().String(),
		Name:     p.Name(),
		Level:    p.Level(),
		Status:   p.Status(),
		Position: p.GetPosition(),
		Stats:    p.Stats(),
	}, nil
}

// GetOnlinePlayersQuery 获取在线玩家查询
type GetOnlinePlayersQuery struct {
	Limit int `json:"limit" validate:"min=1,max=100"`
}

// GetOnlinePlayersResult 获取在线玩家结果
type GetOnlinePlayersResult struct {
	Players []*GetPlayerResult `json:"players"`
	Total   int                `json:"total"`
}

// GetOnlinePlayers 获取在线玩家列表
func (s *PlayerService) GetOnlinePlayers(ctx context.Context, query *GetOnlinePlayersQuery) (*GetOnlinePlayersResult, error) {
	if query.Limit <= 0 {
		query.Limit = 10
	}

	// 查找在线玩家
	players, err := s.playerRepo.FindOnlinePlayers(ctx, query.Limit)
	if err != nil {
		return nil, fmt.Errorf("查找在线玩家失败: %w", err)
	}

	// 转换结果
	results := make([]*GetPlayerResult, 0, len(players))
	for _, p := range players {
		results = append(results, &GetPlayerResult{
			PlayerID: p.ID().String(),
			Name:     p.Name(),
			Level:    p.Level(),
			Status:   p.Status(),
			Position: p.GetPosition(),
			Stats:    p.Stats(),
		})
	}

	return &GetOnlinePlayersResult{
		Players: results,
		Total:   len(results),
	}, nil
}

// Login 玩家登录
func (s *PlayerService) Login(ctx context.Context, playerID string) (*LoginPlayerResult, error) {
	// 解析玩家ID
	pid := player.PlayerID{}
	
	// 查找玩家
	p, err := s.playerRepo.FindByID(ctx, pid)
	if err != nil {
		return nil, fmt.Errorf("获取玩家失败: %w", err)
	}
	if p == nil {
		return nil, player.ErrPlayerNotFound
	}
	
	// 更新玩家状态为在线
	p.SetStatus(player.StatusOnline)
	
	// 保存玩家
	if err := s.playerRepo.Save(ctx, p); err != nil {
		return nil, fmt.Errorf("保存玩家失败: %w", err)
	}
	
	return &LoginPlayerResult{
		PlayerID: p.ID().String(),
		Name:     p.Name(),
		Level:    p.Level(),
		Status:   p.Status(),
		Position: p.Position(),
		Stats:    p.Stats(),
	}, nil
}

// Logout 玩家登出
func (s *PlayerService) Logout(ctx context.Context, playerID string) error {
	// 解析玩家ID
	pid := player.PlayerID{}
	
	// 查找玩家
	p, err := s.playerRepo.FindByID(ctx, pid)
	if err != nil {
		return fmt.Errorf("获取玩家失败: %w", err)
	}
	if p == nil {
		return player.ErrPlayerNotFound
	}
	
	// 更新玩家状态为离线
	p.SetStatus(player.StatusOffline)
	
	// 保存玩家
	if err := s.playerRepo.Save(ctx, p); err != nil {
		return fmt.Errorf("保存玩家失败: %w", err)
	}
	
	return nil
}

// GetPlayerInfo 获取玩家信息
func (s *PlayerService) GetPlayerInfo(ctx context.Context, playerID string) (*LoginPlayerResult, error) {
	// 解析玩家ID
	pid := player.PlayerID{}
	
	// 查找玩家
	p, err := s.playerRepo.FindByID(ctx, pid)
	if err != nil {
		return nil, fmt.Errorf("获取玩家失败: %w", err)
	}
	if p == nil {
		return nil, player.ErrPlayerNotFound
	}
	
	return &LoginPlayerResult{
		PlayerID: p.ID().String(),
		Name:     p.Name(),
		Level:    p.Level(),
		Status:   p.Status(),
		Position: p.Position(),
		Stats:    p.Stats(),
	}, nil
}

// UpdatePlayer 更新玩家信息
func (s *PlayerService) UpdatePlayer(ctx context.Context, playerID string, updates map[string]interface{}) error {
	// 解析玩家ID
	pid := player.PlayerID{}
	
	// 查找玩家
	p, err := s.playerRepo.FindByID(ctx, pid)
	if err != nil {
		return fmt.Errorf("获取玩家失败: %w", err)
	}
	if p == nil {
		return player.ErrPlayerNotFound
	}
	
	// 应用更新
	for key, value := range updates {
		switch key {
		case "name":
			if name, ok := value.(string); ok {
				p.SetName(name)
			}
		case "level":
			if level, ok := value.(int); ok {
				p.SetLevel(level)
			}
		// 可以添加更多字段的更新逻辑
		}
	}
	
	// 保存玩家
	if err := s.playerRepo.Save(ctx, p); err != nil {
		return fmt.Errorf("保存玩家失败: %w", err)
	}
	
	return nil
}
