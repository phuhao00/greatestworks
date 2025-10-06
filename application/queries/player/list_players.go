package player

import (
	"context"
)

// ListPlayersQuery 列表查询玩家请求
type ListPlayersQuery struct {
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
	Name     string `json:"name,omitempty"`
	Status   string `json:"status,omitempty"`
	Level    int    `json:"level,omitempty"`
}

// ListPlayersResult 列表查询玩家结果
type ListPlayersResult struct {
	Players []*PlayerDTO `json:"players"`
	Total   int64        `json:"total"`
	Page    int          `json:"page"`
	Size    int          `json:"size"`
}

// ListPlayersHandler 列表查询玩家处理器
type ListPlayersHandler struct {
	playerQueryService PlayerQueryService
}

// NewListPlayersHandler 创建列表查询玩家处理器
func NewListPlayersHandler(playerQueryService PlayerQueryService) *ListPlayersHandler {
	return &ListPlayersHandler{
		playerQueryService: playerQueryService,
	}
}

// Handle 处理列表查询玩家请求
func (h *ListPlayersHandler) Handle(ctx context.Context, query *ListPlayersQuery) (*ListPlayersResult, error) {
	// 验证查询参数
	if err := query.Validate(); err != nil {
		return nil, err
	}

	// 调用服务层获取玩家列表
	players, total, err := h.playerQueryService.ListPlayers(ctx, query)
	if err != nil {
		return nil, err
	}

	return &ListPlayersResult{
		Players: players,
		Total:   total,
		Page:    query.Page,
		Size:    len(players),
	}, nil
}

// Validate 验证查询参数
func (q *ListPlayersQuery) Validate() error {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
	return nil
}
