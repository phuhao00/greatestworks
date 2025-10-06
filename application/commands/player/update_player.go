package player

import (
	"context"
	"time"
)

// UpdatePlayerCommand 更新玩家命令
type UpdatePlayerCommand struct {
	PlayerID string `json:"player_id" validate:"required"`
	Name     string `json:"name,omitempty" validate:"omitempty,min=2,max=20"`
	Avatar   string `json:"avatar,omitempty"`
	Gender   int    `json:"gender,omitempty" validate:"min=0,max=2"`
}

// UpdatePlayerResult 更新玩家结果
type UpdatePlayerResult struct {
	PlayerID  string    `json:"player_id"`
	Name      string    `json:"name"`
	Level     int       `json:"level"`
	Exp       int64     `json:"exp"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdatePlayerHandler 更新玩家命令处理器
type UpdatePlayerHandler struct {
	playerService PlayerService
}

// NewUpdatePlayerHandler 创建命令处理器
func NewUpdatePlayerHandler(playerService PlayerService) *UpdatePlayerHandler {
	return &UpdatePlayerHandler{
		playerService: playerService,
	}
}

// Handle 处理更新玩家命令
func (h *UpdatePlayerHandler) Handle(ctx context.Context, cmd *UpdatePlayerCommand) (*UpdatePlayerResult, error) {
	// TODO: 实现更新玩家逻辑
	return nil, nil
}

// CommandType 返回命令类型
func (cmd *UpdatePlayerCommand) CommandType() string {
	return "UpdatePlayer"
}

// Validate 验证命令
func (cmd *UpdatePlayerCommand) Validate() error {
	if cmd.PlayerID == "" {
		return ErrInvalidPlayerID
	}
	if cmd.Name != "" && (len(cmd.Name) < 2 || len(cmd.Name) > 20) {
		return ErrInvalidPlayerNameLength
	}
	return nil
}
