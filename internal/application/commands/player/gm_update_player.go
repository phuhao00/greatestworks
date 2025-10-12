package player

import (
	"context"
	"time"
)

// GMUpdatePlayerCommand GM更新玩家命令
type GMUpdatePlayerCommand struct {
	PlayerID string                 `json:"player_id" validate:"required"`
	GMUserID string                 `json:"gm_user_id" validate:"required"`
	GMUser   string                 `json:"gm_user" validate:"required"`
	Reason   string                 `json:"reason" validate:"required"`
	Updates  map[string]interface{} `json:"updates"`
}

// GMUpdatePlayerResult GM更新玩家结果
type GMUpdatePlayerResult struct {
	PlayerID  string                 `json:"player_id"`
	Success   bool                   `json:"success"`
	Updates   map[string]interface{} `json:"updates"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// GMUpdatePlayerHandler GM更新玩家命令处理器
type GMUpdatePlayerHandler struct {
	playerService PlayerService
}

// NewGMUpdatePlayerHandler 创建命令处理器
func NewGMUpdatePlayerHandler(playerService PlayerService) *GMUpdatePlayerHandler {
	return &GMUpdatePlayerHandler{
		playerService: playerService,
	}
}

// Handle 处理GM更新玩家命令
func (h *GMUpdatePlayerHandler) Handle(ctx context.Context, cmd *GMUpdatePlayerCommand) (*GMUpdatePlayerResult, error) {
	// 验证命令
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	// TODO: 实现GM更新玩家逻辑
	// err := h.playerService.GMUpdatePlayer(ctx, cmd)
	// if err != nil {
	// 	return nil, err
	// }

	return &GMUpdatePlayerResult{
		PlayerID:  cmd.PlayerID,
		Success:   true,
		Updates:   cmd.Updates,
		UpdatedAt: time.Now(),
	}, nil
}

// CommandType 返回命令类型
func (cmd *GMUpdatePlayerCommand) CommandType() string {
	return "GMUpdatePlayer"
}

// Validate 验证命令
func (cmd *GMUpdatePlayerCommand) Validate() error {
	if cmd.PlayerID == "" {
		return ErrInvalidPlayerID
	}
	if cmd.GMUserID == "" || cmd.GMUser == "" {
		return ErrInvalidRequest
	}
	if cmd.Reason == "" {
		return ErrInvalidRequest
	}
	return nil
}
