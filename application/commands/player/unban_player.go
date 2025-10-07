package player

import (
	"context"
	"time"
)

// UnbanPlayerCommand GM解封玩家命令
type UnbanPlayerCommand struct {
	PlayerID       string `json:"player_id" validate:"required"`
	UnbannedBy     string `json:"unbanned_by" validate:"required"`
	UnbannedByName string `json:"unbanned_by_name" validate:"required"`
	Reason         string `json:"reason" validate:"required"`
}

// UnbanPlayerResult GM解封玩家结果
type UnbanPlayerResult struct {
	PlayerID   string    `json:"player_id"`
	Success    bool      `json:"success"`
	UnbannedAt time.Time `json:"unbanned_at"`
}

// UnbanPlayerHandler GM解封玩家命令处理器
type UnbanPlayerHandler struct {
	playerService PlayerService
}

// NewUnbanPlayerHandler 创建命令处理器
func NewUnbanPlayerHandler(playerService PlayerService) *UnbanPlayerHandler {
	return &UnbanPlayerHandler{
		playerService: playerService,
	}
}

// Handle 处理GM解封玩家命令
func (h *UnbanPlayerHandler) Handle(ctx context.Context, cmd *UnbanPlayerCommand) (*UnbanPlayerResult, error) {
	// 验证命令
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	// TODO: 实现GM解封玩家逻辑
	// err := h.playerService.UnbanPlayer(ctx, cmd)
	// if err != nil {
	// 	return nil, err
	// }

	return &UnbanPlayerResult{
		PlayerID:   cmd.PlayerID,
		Success:    true,
		UnbannedAt: time.Now(),
	}, nil
}

// CommandType 返回命令类型
func (cmd *UnbanPlayerCommand) CommandType() string {
	return "UnbanPlayer"
}

// Validate 验证命令
func (cmd *UnbanPlayerCommand) Validate() error {
	if cmd.PlayerID == "" {
		return ErrInvalidPlayerID
	}
	if cmd.UnbannedBy == "" || cmd.UnbannedByName == "" {
		return ErrInvalidRequest
	}
	if cmd.Reason == "" {
		return ErrInvalidRequest
	}
	return nil
}
