package player

import (
	"context"
	"greatestworks/internal/application/interfaces"
	"time"
)

// CreatePlayerCommand 创建玩家命令
type CreatePlayerCommand struct {
	Name   string `json:"name" validate:"required,min=2,max=20"`
	Avatar string `json:"avatar,omitempty"`
	Gender int    `json:"gender,omitempty" validate:"min=0,max=2"`
}

// CreatePlayerResult 创建玩家结果
type CreatePlayerResult struct {
	PlayerID  string    `json:"player_id"`
	Name      string    `json:"name"`
	Level     int       `json:"level"`
	CreatedAt time.Time `json:"created_at"`
}

// CreatePlayerWithAccountCommand 创建带账户的玩家命令
type CreatePlayerWithAccountCommand struct {
	Username     string `json:"username" validate:"required,min=3,max=50"`
	PasswordHash string `json:"password_hash" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	PlayerName   string `json:"player_name" validate:"required,min=2,max=50"`
	Avatar       string `json:"avatar,omitempty"`
	Gender       int    `json:"gender,omitempty" validate:"min=0,max=2"`
}

// CommandType 返回命令类型
func (cmd *CreatePlayerWithAccountCommand) CommandType() string {
	return "CreatePlayerWithAccount"
}

// Validate 验证命令
func (cmd *CreatePlayerWithAccountCommand) Validate() error {
	if cmd.Username == "" {
		return ErrInvalidUsername
	}
	if cmd.PasswordHash == "" {
		return ErrInvalidPassword
	}
	if cmd.Email == "" {
		return ErrInvalidEmail
	}
	if cmd.PlayerName == "" {
		return ErrInvalidPlayerName
	}
	return nil
}

// CreatePlayerWithAccountResult 创建带账户的玩家结果
type CreatePlayerWithAccountResult struct {
	PlayerID  string    `json:"player_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// CreatePlayerHandler 创建玩家命令处理器
type CreatePlayerHandler struct {
	playerService PlayerService
}

// 确保实现了接口
var _ interfaces.CommandHandler[*CreatePlayerCommand, *CreatePlayerResult] = (*CreatePlayerHandler)(nil)

// PlayerService 玩家服务接口
type PlayerService interface {
	CreatePlayer(ctx context.Context, name string) (*CreatePlayerResult, error)
	MovePlayer(ctx context.Context, playerID string, position Position) error
	LevelUpPlayer(ctx context.Context, cmd *LevelUpPlayerCommand) (*LevelUpPlayerResult, error)
	DeletePlayer(ctx context.Context, cmd *DeletePlayerCommand) (*DeletePlayerResult, error)
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
