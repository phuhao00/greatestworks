package player

import (
	"context"
	"time"
)

// CreatePlayerCommand 创建玩家命令
type CreatePlayerCommand struct {
	Name string `json:"name" validate:"required,min=2,max=20"`
}

// CreatePlayerResult 创建玩家结果
type CreatePlayerResult struct {
	PlayerID  string    `json:"player_id"`
	Name      string    `json:"name"`
	Level     int       `json:"level"`
	CreatedAt time.Time `json:"created_at"`
}

// CreatePlayerHandler 创建玩家命令处理器
type CreatePlayerHandler struct {
	playerService PlayerService
}

// PlayerService 玩家服务接口
type PlayerService interface {
	CreatePlayer(ctx context.Context, name string) (*CreatePlayerResult, error)
}

// NewCreatePlayerHandler 创建命令处理器
func NewCreatePlayerHandler(playerService PlayerService) *CreatePlayerHandler {
	return &CreatePlayerHandler{
		playerService: playerService,
	}
}

// Handle 处理创建玩家命令
func (h *CreatePlayerHandler) Handle(ctx context.Context, cmd *CreatePlayerCommand) (*CreatePlayerResult, error) {
	return h.playerService.CreatePlayer(ctx, cmd.Name)
}

// CommandType 返回命令类型
func (cmd *CreatePlayerCommand) CommandType() string {
	return "CreatePlayer"
}

// Validate 验证命令
func (cmd *CreatePlayerCommand) Validate() error {
	if cmd.Name == "" {
		return ErrInvalidPlayerName
	}
	if len(cmd.Name) < 2 || len(cmd.Name) > 20 {
		return ErrInvalidPlayerNameLength
	}
	return nil
}