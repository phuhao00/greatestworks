package player

import (
	"context"
	"time"
)

// BanPlayerCommand GM封禁玩家命令
type BanPlayerCommand struct {
	PlayerID     string    `json:"player_id" validate:"required"`
	BannedBy     string    `json:"banned_by" validate:"required"`
	BannedByName string    `json:"banned_by_name" validate:"required"`
	Reason       string    `json:"reason" validate:"required"`
	BanType      string    `json:"ban_type" validate:"required"` // permanent, temporary
	BanUntil     time.Time `json:"ban_until,omitempty"`
}

// BanPlayerResult GM封禁玩家结果
type BanPlayerResult struct {
	PlayerID string    `json:"player_id"`
	Success  bool      `json:"success"`
	BannedAt time.Time `json:"banned_at"`
	BanUntil time.Time `json:"ban_until,omitempty"`
}

// BanPlayerHandler GM封禁玩家命令处理器
type BanPlayerHandler struct {
	playerService PlayerService
}

// NewBanPlayerHandler 创建命令处理器
func NewBanPlayerHandler(playerService PlayerService) *BanPlayerHandler {
	return &BanPlayerHandler{
		playerService: playerService,
	}
}

// Handle 处理GM封禁玩家命令
func (h *BanPlayerHandler) Handle(ctx context.Context, cmd *BanPlayerCommand) (*BanPlayerResult, error) {
	// 验证命令
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	// TODO: 实现GM封禁玩家逻辑
	// err := h.playerService.BanPlayer(ctx, cmd)
	// if err != nil {
	// 	return nil, err
	// }

	return &BanPlayerResult{
		PlayerID: cmd.PlayerID,
		Success:  true,
		BannedAt: time.Now(),
		BanUntil: cmd.BanUntil,
	}, nil
}

// CommandType 返回命令类型
func (cmd *BanPlayerCommand) CommandType() string {
	return "BanPlayer"
}

// Validate 验证命令
func (cmd *BanPlayerCommand) Validate() error {
	if cmd.PlayerID == "" {
		return ErrInvalidPlayerID
	}
	if cmd.BannedBy == "" || cmd.BannedByName == "" {
		return ErrInvalidRequest
	}
	if cmd.Reason == "" {
		return ErrInvalidRequest
	}
	if cmd.BanType != "permanent" && cmd.BanType != "temporary" {
		return ErrInvalidRequest
	}
	if cmd.BanType == "temporary" && cmd.BanUntil.IsZero() {
		return ErrInvalidRequest
	}
	return nil
}
