package battle

import (
	"context"
	"time"
)

// BattleType 战斗类型
type BattleType int

const (
	BattleTypePvP BattleType = iota
	BattleTypePvE
	BattleTypeTeamPvP
	BattleTypeRaid
)

// CreateBattleCommand 创建战斗命令
type CreateBattleCommand struct {
	BattleType BattleType `json:"battle_type" validate:"required"`
	CreatorID  string     `json:"creator_id" validate:"required"`
}

// CreateBattleResult 创建战斗结果
type CreateBattleResult struct {
	BattleID   string     `json:"battle_id"`
	BattleType BattleType `json:"battle_type"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
}

// CreateBattleHandler 创建战斗命令处理器
type CreateBattleHandler struct {
	battleService BattleService
}

// BattleService 战斗服务接口
type BattleService interface {
	CreateBattle(ctx context.Context, battleType BattleType, creatorID string) (*CreateBattleResult, error)
}

// NewCreateBattleHandler 创建命令处理器
func NewCreateBattleHandler(battleService BattleService) *CreateBattleHandler {
	return &CreateBattleHandler{
		battleService: battleService,
	}
}

// Handle 处理创建战斗命令
func (h *CreateBattleHandler) Handle(ctx context.Context, cmd *CreateBattleCommand) (*CreateBattleResult, error) {
	return h.battleService.CreateBattle(ctx, cmd.BattleType, cmd.CreatorID)
}

// CommandType 返回命令类型
func (cmd *CreateBattleCommand) CommandType() string {
	return "CreateBattle"
}

// Validate 验证命令
func (cmd *CreateBattleCommand) Validate() error {
	if cmd.CreatorID == "" {
		return ErrInvalidCreatorID
	}

	if cmd.BattleType < BattleTypePvP || cmd.BattleType > BattleTypeRaid {
		return ErrInvalidBattleType
	}

	return nil
}
