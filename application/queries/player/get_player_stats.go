package player

import (
	"context"
	"time"

	"greatestworks/internal/domain/player"
)

// GetPlayerStatsQuery 获取玩家统计信息查询
type GetPlayerStatsQuery struct {
	PlayerID player.PlayerID `json:"player_id" validate:"required"`
}

// GetPlayerStatsResult 获取玩家统计信息结果
type GetPlayerStatsResult struct {
	Found         bool                `json:"found"`
	PlayerID      player.PlayerID     `json:"player_id"`
	TotalBattles  int                 `json:"total_battles"`
	Wins          int                 `json:"wins"`
	Losses        int                 `json:"losses"`
	WinRate       float64             `json:"win_rate"`
	TotalExp      int64               `json:"total_exp"`
	PlayTime      time.Duration       `json:"play_time"`
	LastLogin     *time.Time          `json:"last_login"`
	Achievements  []string            `json:"achievements"`
}

// GetPlayerStatsHandler 获取玩家统计信息处理器
type GetPlayerStatsHandler struct {
	playerStatsService PlayerStatsService
}

// PlayerStatsService 玩家统计服务接口
type PlayerStatsService interface {
	GetPlayerStats(ctx context.Context, playerID player.PlayerID) (*GetPlayerStatsResult, error)
}

// NewGetPlayerStatsHandler 创建获取玩家统计信息处理器
func NewGetPlayerStatsHandler(playerStatsService PlayerStatsService) *GetPlayerStatsHandler {
	return &GetPlayerStatsHandler{
		playerStatsService: playerStatsService,
	}
}

// Handle 处理获取玩家统计信息查询
func (h *GetPlayerStatsHandler) Handle(ctx context.Context, query *GetPlayerStatsQuery) (*GetPlayerStatsResult, error) {
	return h.playerStatsService.GetPlayerStats(ctx, query.PlayerID)
}

// QueryType 返回查询类型
func (query *GetPlayerStatsQuery) QueryType() string {
	return "GetPlayerStatsQuery"
}

// Validate 验证查询参数
func (query *GetPlayerStatsQuery) Validate() error {
	if query.PlayerID.String() == "" {
		return ErrInvalidPlayerID
	}
	return nil
}
