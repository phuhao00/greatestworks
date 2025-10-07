package player

import (
	"context"
)

// SearchPlayersQuery GM搜索玩家查询
type SearchPlayersQuery struct {
	Keyword   string `json:"keyword,omitempty"`
	PlayerID  string `json:"player_id,omitempty"`
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	Status    string `json:"status,omitempty"`
	MinLevel  int    `json:"min_level,omitempty"`
	MaxLevel  int    `json:"max_level,omitempty"`
	Page      int    `json:"page" validate:"min=1"`
	PageSize  int    `json:"page_size" validate:"min=1,max=100"`
	SortBy    string `json:"sort_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}

// SearchPlayersResult GM搜索玩家结果
type SearchPlayersResult struct {
	Players []*PlayerDTO `json:"players"`
	Total   int64        `json:"total"`
	Page    int          `json:"page"`
	Size    int          `json:"size"`
}

// SearchPlayersHandler GM搜索玩家处理器
type SearchPlayersHandler struct {
	playerQueryService PlayerQueryService
}

// NewSearchPlayersHandler 创建GM搜索玩家处理器
func NewSearchPlayersHandler(playerQueryService PlayerQueryService) *SearchPlayersHandler {
	return &SearchPlayersHandler{
		playerQueryService: playerQueryService,
	}
}

// Handle 处理GM搜索玩家请求
func (h *SearchPlayersHandler) Handle(ctx context.Context, query *SearchPlayersQuery) (*SearchPlayersResult, error) {
	// 验证查询参数
	if err := query.Validate(); err != nil {
		return nil, err
	}

	// TODO: 调用服务层搜索玩家
	// players, total, err := h.playerQueryService.SearchPlayers(ctx, query)
	// if err != nil {
	// 	return nil, err
	// }

	return &SearchPlayersResult{
		Players: make([]*PlayerDTO, 0),
		Total:   0,
		Page:    query.Page,
		Size:    query.PageSize,
	}, nil
}

// Validate 验证查询
func (q *SearchPlayersQuery) Validate() error {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
	return nil
}
