package player

import (
	"context"
	"time"

	"greatestworks/internal/domain/player"
)

// LevelUpPlayerCommand 玩家升级命令
type LevelUpPlayerCommand struct {
	PlayerID string `json:"player_id" validate:"required"`
	ExpGain  int64  `json:"exp_gain,omitempty"`
}

// LevelUpPlayerResult 玩家升级结果
type LevelUpPlayerResult struct {
	PlayerID  string              `json:"player_id"`
	OldLevel  int                 `json:"old_level"`
	NewLevel  int                 `json:"new_level"`
	OldExp    int64               `json:"old_exp"`
	NewExp    int64               `json:"new_exp"`
	LeveledUp bool                `json:"leveled_up"`
	Status    player.PlayerStatus `json:"status"`
	UpdatedAt time.Time           `json:"updated_at"`
}

// LevelUpPlayerHandler 玩家升级命令处理器
type LevelUpPlayerHandler struct {
	playerService PlayerService
}

// NewLevelUpPlayerHandler 创建玩家升级命令处理器
func NewLevelUpPlayerHandler(playerService PlayerService) *LevelUpPlayerHandler {
	return &LevelUpPlayerHandler{
		playerService: playerService,
	}
}

// Handle 处理玩家升级命令
func (h *LevelUpPlayerHandler) Handle(ctx context.Context, cmd *LevelUpPlayerCommand) (*LevelUpPlayerResult, error) {
	// 验证命令
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	// 调用服务层执行升级
	result, err := h.playerService.LevelUpPlayer(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CommandType 返回命令类型
func (cmd *LevelUpPlayerCommand) CommandType() string {
	return "LevelUpPlayer"
}

// Validate 验证命令
func (cmd *LevelUpPlayerCommand) Validate() error {
	if cmd.PlayerID == "" {
		return ErrInvalidPlayerID
	}
	return nil
}
