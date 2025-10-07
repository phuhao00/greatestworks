package player

import (
	"context"
	"time"
	// TODO: Add necessary imports when implementing Delete functionality
	// "greatestworks/internal/domain/player"
)

// DeletePlayerCommand 删除玩家命令
type DeletePlayerCommand struct {
	PlayerID string `json:"player_id" validate:"required"`
	Reason   string `json:"reason,omitempty"`
}

// DeletePlayerResult 删除玩家结果
type DeletePlayerResult struct {
	PlayerID  string    `json:"player_id"`
	Deleted   bool      `json:"deleted"`
	DeletedAt time.Time `json:"deleted_at"`
}

// DeletePlayerHandler 删除玩家命令处理器
type DeletePlayerHandler struct {
	playerService PlayerService
}

// NewDeletePlayerHandler 创建删除玩家命令处理器
func NewDeletePlayerHandler(playerService PlayerService) *DeletePlayerHandler {
	return &DeletePlayerHandler{
		playerService: playerService,
	}
}

// Handle 处理删除玩家命令
func (h *DeletePlayerHandler) Handle(ctx context.Context, cmd *DeletePlayerCommand) (*DeletePlayerResult, error) {
	// 验证命令
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	// 调用服务层执行删除
	result, err := h.playerService.DeletePlayer(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CommandType 返回命令类型
func (cmd *DeletePlayerCommand) CommandType() string {
	return "DeletePlayer"
}

// Validate 验证命令
func (cmd *DeletePlayerCommand) Validate() error {
	if cmd.PlayerID == "" {
		return ErrInvalidPlayerID
	}
	return nil
}
