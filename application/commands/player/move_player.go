package player

import (
	"context"
)

// Position 位置信息
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// MovePlayerCommand 移动玩家命令
type MovePlayerCommand struct {
	PlayerID string   `json:"player_id" validate:"required"`
	Position Position `json:"position" validate:"required"`
}

// MovePlayerResult 移动玩家结果
type MovePlayerResult struct {
	PlayerID    string   `json:"player_id"`
	OldPosition Position `json:"old_position"`
	NewPosition Position `json:"new_position"`
	Success     bool     `json:"success"`
}

// MovePlayerHandler 移动玩家命令处理器
type MovePlayerHandler struct {
	playerService PlayerService
}

// NewMovePlayerHandler 创建移动玩家命令处理器
func NewMovePlayerHandler(playerService PlayerService) *MovePlayerHandler {
	return &MovePlayerHandler{
		playerService: playerService,
	}
}

// Handle 处理移动玩家命令
func (h *MovePlayerHandler) Handle(ctx context.Context, cmd *MovePlayerCommand) (*MovePlayerResult, error) {
	return h.playerService.MovePlayer(ctx, cmd.PlayerID, cmd.Position)
}

// CommandType 返回命令类型
func (cmd *MovePlayerCommand) CommandType() string {
	return "MovePlayer"
}

// Validate 验证命令
func (cmd *MovePlayerCommand) Validate() error {
	if cmd.PlayerID == "" {
		return ErrInvalidPlayerID
	}
	
	// 验证位置范围
	if cmd.Position.X < -1000 || cmd.Position.X > 1000 {
		return ErrInvalidPosition
	}
	if cmd.Position.Y < -1000 || cmd.Position.Y > 1000 {
		return ErrInvalidPosition
	}
	if cmd.Position.Z < -100 || cmd.Position.Z > 100 {
		return ErrInvalidPosition
	}
	
	return nil
}