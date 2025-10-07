package player

import (
	"context"
)

// GetPlayerDetailQuery GM获取玩家详情查询
type GetPlayerDetailQuery struct {
	PlayerID string `json:"player_id" validate:"required"`
}

// GetPlayerDetailResult GM获取玩家详情结果
type GetPlayerDetailResult struct {
	Player *PlayerDTO `json:"player"`
	Found  bool       `json:"found"`
}

// GetPlayerDetailHandler GM获取玩家详情处理器
type GetPlayerDetailHandler struct {
	playerQueryService PlayerQueryService
}

// NewGetPlayerDetailHandler 创建GM获取玩家详情处理器
func NewGetPlayerDetailHandler(playerQueryService PlayerQueryService) *GetPlayerDetailHandler {
	return &GetPlayerDetailHandler{
		playerQueryService: playerQueryService,
	}
}

// Handle 处理GM获取玩家详情请求
func (h *GetPlayerDetailHandler) Handle(ctx context.Context, query *GetPlayerDetailQuery) (*GetPlayerDetailResult, error) {
	// 验证查询参数
	if err := query.Validate(); err != nil {
		return nil, err
	}

	// TODO: 调用服务层获取玩家详情
	// player, err := h.playerQueryService.GetPlayerDetail(ctx, query.PlayerID)
	// if err != nil {
	// 	return nil, err
	// }

	return &GetPlayerDetailResult{
		Player: nil,
		Found:  false,
	}, nil
}

// Validate 验证查询
func (q *GetPlayerDetailQuery) Validate() error {
	if q.PlayerID == "" {
		return ErrInvalidPlayerID
	}
	return nil
}
