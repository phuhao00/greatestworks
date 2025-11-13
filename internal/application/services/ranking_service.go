package services

import (
	"context"
	"fmt"
	"time"

	"greatestworks/internal/domain/ranking"
)

// RankingApplicationService 排行榜应用服务
type RankingApplicationService struct {
	rankingRepo    ranking.RankingRepository
	rankEntryRepo  ranking.RankEntryRepository
	rankingService *ranking.RankingService
	eventBus       ranking.RankingEventBus
}

// NewRankingApplicationService 创建排行榜应用服务
func NewRankingApplicationService(
	rankingRepo ranking.RankingRepository,
	rankEntryRepo ranking.RankEntryRepository,
	rankingService *ranking.RankingService,
	eventBus ranking.RankingEventBus,
) *RankingApplicationService {
	return &RankingApplicationService{
		rankingRepo:    rankingRepo,
		rankEntryRepo:  rankEntryRepo,
		rankingService: rankingService,
		eventBus:       eventBus,
	}
}

// CreateRankingRequest 创建排行榜请求
type CreateRankingRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	RankType    string `json:"rank_type"`
	PeriodType  string `json:"period_type"`
	MaxEntries  int32  `json:"max_entries"`
	IsActive    bool   `json:"is_active"`
}

// CreateRankingResponse 创建排行榜响应
type CreateRankingResponse struct {
	RankingID   string    `json:"ranking_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	RankType    string    `json:"rank_type"`
	PeriodType  string    `json:"period_type"`
	MaxEntries  int32     `json:"max_entries"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateRanking 创建排行榜
func (s *RankingApplicationService) CreateRanking(ctx context.Context, req *CreateRankingRequest) (*CreateRankingResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	if err := s.validateCreateRankingRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// 转换排行榜类型
	// rankType, err := s.parseRankType(req.RankType)
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid rank type: %w", err)
	// }

	// 转换周期类型
	// periodType, err := s.parsePeriodType(req.PeriodType)
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid period type: %w", err)
	// }

	// 创建排行榜聚合根
	// TODO: 修复NewRankingAggregate方法调用
	// rankingAggregate := ranking.NewRankingAggregate(req.Name, rankType, periodType)
	// rankingAggregate := &ranking.RankingAggregate{}
	// TODO: 修复SetDescription方法调用
	// rankingAggregate.SetDescription(req.Description)
	// TODO: 修复SetMaxEntries方法调用
	// rankingAggregate.SetMaxEntries(req.MaxEntries)
	// TODO: 修复Activate方法调用
	// if req.IsActive {
	// 	rankingAggregate.Activate()
	// }

	// 保存排行榜
	// TODO: 修复Save方法调用
	// if err := s.rankingRepo.Save(ctx, rankingAggregate); err != nil {
	// 	return nil, fmt.Errorf("failed to save ranking: %w", err)
	// }

	// 发布事件
	// TODO: 修复NewRankingCreatedEvent方法调用
	// event := ranking.NewRankingCreatedEvent(rankingAggregate.GetID(), req.Name, rankType, periodType)
	// event := &ranking.RankingCreatedEvent{}
	// TODO: 修复Publish方法调用
	// if err := s.eventBus.Publish(ctx, event); err != nil {
	// 	fmt.Printf("failed to publish ranking created event: %v\n", err)
	// }

	return &CreateRankingResponse{
		RankingID:   "",         // TODO: rankingAggregate.GetID(),
		Name:        "",         // TODO: rankingAggregate.GetName(),
		Description: "",         // TODO: rankingAggregate.GetDescription(),
		RankType:    "",         // TODO: rankingAggregate.GetRankType().String(),
		PeriodType:  "",         // TODO: rankingAggregate.GetPeriodType().String(),
		MaxEntries:  0,          // TODO: rankingAggregate.GetMaxEntries(),
		IsActive:    false,      // TODO: rankingAggregate.IsActive(),
		CreatedAt:   time.Now(), // TODO: rankingAggregate.GetCreatedAt(),
	}, nil
}

// UpdatePlayerScoreRequest 更新玩家分数请求
type UpdatePlayerScoreRequest struct {
	RankingID string                 `json:"ranking_id"`
	PlayerID  uint64                 `json:"player_id"`
	Score     int64                  `json:"score"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// UpdatePlayerScoreResponse 更新玩家分数响应
type UpdatePlayerScoreResponse struct {
	RankingID    string `json:"ranking_id"`
	PlayerID     uint64 `json:"player_id"`
	OldScore     int64  `json:"old_score"`
	NewScore     int64  `json:"new_score"`
	OldRank      int32  `json:"old_rank"`
	NewRank      int32  `json:"new_rank"`
	RankChanged  bool   `json:"rank_changed"`
	ScoreChanged bool   `json:"score_changed"`
}

// UpdatePlayerScore 更新玩家分数
func (s *RankingApplicationService) UpdatePlayerScore(ctx context.Context, req *UpdatePlayerScoreRequest) (*UpdatePlayerScoreResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	if err := s.validateUpdatePlayerScoreRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// 获取排行榜
	// rankingAggregate, err := s.rankingRepo.FindByID(ctx, req.RankingID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to find ranking: %w", err)
	// }
	// if rankingAggregate == nil {
	// 	return nil, fmt.Errorf("ranking not found")
	// }
	// rankingAggregate := &ranking.RankingAggregate{}

	// 检查排行榜是否激活
	// TODO: 修复IsActive方法调用
	// if !rankingAggregate.IsActive() {
	// 	return nil, fmt.Errorf("ranking is not active")
	// }

	// 获取当前玩家排名条目
	// TODO: 修复FindByRankingAndPlayer方法调用
	// currentEntry, err := s.rankEntryRepo.FindByRankingAndPlayer(ctx, req.RankingID, req.PlayerID)
	// TODO: 修复err变量
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to find current entry: %w", err)
	// }

	oldScore := int64(0)
	oldRank := int64(0)

	// 更新分数
	// entry, err := s.rankingService.UpdatePlayerScore(ctx, req.RankingID, req.PlayerID, req.Score)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to update player score: %w", err)
	// }
	entry := &ranking.RankEntry{}

	// 设置元数据
	// if req.Metadata != nil {
	// 	for key, value := range req.Metadata {
	// 		entry.SetMetadata(key, value)
	// 	}
	// 	if err := s.rankEntryRepo.Save(ctx, entry); err != nil {
	// 		return nil, fmt.Errorf("failed to save entry metadata: %w", err)
	// 	}
	// }

	newScore := entry.GetScore()
	newRank := entry.GetRank()
	rankChanged := oldRank != newRank
	scoreChanged := oldScore != newScore

	// 发布事件
	// if scoreChanged {
	// 	event := ranking.NewPlayerScoreUpdatedEvent(req.RankingID, req.PlayerID, oldScore, newScore)
	// 	if err := s.eventBus.Publish(ctx, event); err != nil {
	// 		fmt.Printf("failed to publish score updated event: %v\n", err)
	// 	}
	// }

	// if rankChanged {
	// 	event := ranking.NewPlayerRankChangedEvent(req.RankingID, req.PlayerID, oldRank, newRank)
	// 	if err := s.eventBus.Publish(ctx, event); err != nil {
	// 		fmt.Printf("failed to publish rank changed event: %v\n", err)
	// 	}
	// }

	return &UpdatePlayerScoreResponse{
		RankingID:    req.RankingID,
		PlayerID:     req.PlayerID,
		OldScore:     oldScore,
		NewScore:     newScore,
		OldRank:      int32(oldRank),
		NewRank:      int32(newRank),
		RankChanged:  rankChanged,
		ScoreChanged: scoreChanged,
	}, nil
}

// GetRankingRequest 获取排行榜请求
type GetRankingRequest struct {
	RankingID string `json:"ranking_id"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
}

// RankEntryResponse 排名条目响应
type RankEntryResponse struct {
	PlayerID  uint64                 `json:"player_id"`
	Rank      int32                  `json:"rank"`
	Score     int64                  `json:"score"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// GetRankingResponse 获取排行榜响应
type GetRankingResponse struct {
	RankingID   string               `json:"ranking_id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	RankType    string               `json:"rank_type"`
	PeriodType  string               `json:"period_type"`
	Entries     []*RankEntryResponse `json:"entries"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	PageSize    int                  `json:"page_size"`
	TotalPages  int64                `json:"total_pages"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

// GetRanking 获取排行榜
func (s *RankingApplicationService) GetRanking(ctx context.Context, req *GetRankingRequest) (*GetRankingResponse, error) {
	if req == nil || req.RankingID == "" {
		return nil, fmt.Errorf("ranking ID is required")
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 50
	}

	// 获取排行榜
	// rankingAggregate, err := s.rankingRepo.FindByID(ctx, req.RankingID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to find ranking: %w", err)
	// }
	// if rankingAggregate == nil {
	// 	return nil, fmt.Errorf("ranking not found")
	// }

	// 构建查询
	// query := ranking.NewRankEntryQuery().
	// 	WithRanking(req.RankingID).
	// 	WithSort("rank", "asc").
	// 	WithPagination(req.Page, req.PageSize)

	// 查询排名条目
	// entries, total, err := s.rankEntryRepo.FindByQuery(ctx, query)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to find rank entries: %w", err)
	// }
	entries := []*ranking.RankEntry{}
	total := 0

	// 转换响应
	entryResponses := make([]*RankEntryResponse, len(entries))
	for i, entry := range entries {
		entryResponses[i] = &RankEntryResponse{
			PlayerID:  entry.GetPlayerID(),
			Rank:      int32(entry.GetRank()),
			Score:     entry.GetScore(),
			Metadata:  entry.GetMetadata(),
			UpdatedAt: entry.GetUpdatedAt(),
		}
	}

	totalPages := (int64(total) + int64(req.PageSize) - 1) / int64(req.PageSize)

	return &GetRankingResponse{
		RankingID:   "", // TODO: rankingAggregate.GetID(),
		Name:        "", // TODO: rankingAggregate.GetName(),
		Description: "", // TODO: rankingAggregate.GetDescription(),
		RankType:    "", // TODO: rankingAggregate.GetRankType().String(),
		PeriodType:  "", // TODO: rankingAggregate.GetPeriodType().String(),
		Entries:     entryResponses,
		Total:       int64(total),
		Page:        req.Page,
		PageSize:    req.PageSize,
		TotalPages:  totalPages,
		UpdatedAt:   time.Now(), // TODO: rankingAggregate.GetUpdatedAt(),
	}, nil
}

// GetPlayerRankRequest 获取玩家排名请求
type GetPlayerRankRequest struct {
	RankingID string `json:"ranking_id"`
	PlayerID  uint64 `json:"player_id"`
}

// GetPlayerRankResponse 获取玩家排名响应
type GetPlayerRankResponse struct {
	RankingID string                 `json:"ranking_id"`
	PlayerID  uint64                 `json:"player_id"`
	Rank      int32                  `json:"rank"`
	Score     int64                  `json:"score"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	UpdatedAt time.Time              `json:"updated_at"`
	Found     bool                   `json:"found"`
}

// GetPlayerRank 获取玩家排名
func (s *RankingApplicationService) GetPlayerRank(ctx context.Context, req *GetPlayerRankRequest) (*GetPlayerRankResponse, error) {
	if req == nil || req.RankingID == "" || req.PlayerID == 0 {
		return nil, fmt.Errorf("ranking ID and player ID are required")
	}

	// 获取玩家排名条目
	// entry, err := s.rankEntryRepo.FindByRankingAndPlayer(ctx, req.RankingID, req.PlayerID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to find player rank: %w", err)
	// }
	// 占位实现：未接入存储，默认未找到
	return &GetPlayerRankResponse{
		RankingID: req.RankingID,
		PlayerID:  req.PlayerID,
		Found:     false,
	}, nil

	// 以下为真实实现路径，待存储接入后启用
	// return &GetPlayerRankResponse{
	//     RankingID: req.RankingID,
	//     PlayerID:  req.PlayerID,
	//     Rank:      int32(entry.GetRank()),
	//     Score:     entry.GetScore(),
	//     Metadata:  entry.GetMetadata(),
	//     UpdatedAt: entry.GetUpdatedAt(),
	//     Found:     true,
	// }, nil
}

// GetTopPlayersRequest 获取排行榜前N名请求
type GetTopPlayersRequest struct {
	RankingID string `json:"ranking_id"`
	Limit     int    `json:"limit"`
}

// GetTopPlayersResponse 获取排行榜前N名响应
type GetTopPlayersResponse struct {
	RankingID string               `json:"ranking_id"`
	Entries   []*RankEntryResponse `json:"entries"`
	UpdatedAt time.Time            `json:"updated_at"`
}

// GetTopPlayers 获取排行榜前N名
func (s *RankingApplicationService) GetTopPlayers(ctx context.Context, req *GetTopPlayersRequest) (*GetTopPlayersResponse, error) {
	if req == nil || req.RankingID == "" {
		return nil, fmt.Errorf("ranking ID is required")
	}

	// 设置默认值
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100 // 限制最大数量
	}

	// 获取排行榜
	// rankingAggregate, err := s.rankingRepo.FindByID(ctx, req.RankingID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to find ranking: %w", err)
	// }
	// if rankingAggregate == nil {
	// 	return nil, fmt.Errorf("ranking not found")
	// }

	// 获取前N名
	// entries, err := s.rankingService.GetTopPlayers(ctx, req.RankingID, req.Limit)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get top players: %w", err)
	// }
	entries := []*ranking.RankEntry{}

	// 转换响应
	entryResponses := make([]*RankEntryResponse, len(entries))
	for i, entry := range entries {
		entryResponses[i] = &RankEntryResponse{
			PlayerID:  entry.GetPlayerID(),
			Rank:      int32(entry.GetRank()),
			Score:     entry.GetScore(),
			Metadata:  entry.GetMetadata(),
			UpdatedAt: entry.GetUpdatedAt(),
		}
	}

	return &GetTopPlayersResponse{
		RankingID: req.RankingID,
		Entries:   entryResponses,
		UpdatedAt: time.Now(), // TODO: rankingAggregate.GetUpdatedAt(),
	}, nil
}

// ResetRankingRequest 重置排行榜请求
type ResetRankingRequest struct {
	RankingID string `json:"ranking_id"`
	Reason    string `json:"reason"`
}

// ResetRankingResponse 重置排行榜响应
type ResetRankingResponse struct {
	RankingID      string    `json:"ranking_id"`
	EntriesCleared int64     `json:"entries_cleared"`
	ResetAt        time.Time `json:"reset_at"`
}

// ResetRanking 重置排行榜
func (s *RankingApplicationService) ResetRanking(ctx context.Context, req *ResetRankingRequest) (*ResetRankingResponse, error) {
	if req == nil || req.RankingID == "" {
		return nil, fmt.Errorf("ranking ID is required")
	}

	// 获取排行榜
	// rankingAggregate, err := s.rankingRepo.FindByID(ctx, req.RankingID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to find ranking: %w", err)
	// }
	// if rankingAggregate == nil {
	// 	return nil, fmt.Errorf("ranking not found")
	// }

	// 重置排行榜
	// entriesCleared, err := s.rankingService.ResetRanking(ctx, req.RankingID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to reset ranking: %w", err)
	// }
	// entriesCleared := &ranking.RankingOperationResult{}

	// 更新排行榜状态
	// rankingAggregate.Reset()
	// if err := s.rankingRepo.Save(ctx, rankingAggregate); err != nil {
	// 	return nil, fmt.Errorf("failed to save ranking: %w", err)
	// }

	// 发布事件
	// event := ranking.NewRankingResetEvent(req.RankingID, req.Reason, entriesCleared)
	// if err := s.eventBus.Publish(ctx, event); err != nil {
	// 	fmt.Printf("failed to publish ranking reset event: %v\n", err)
	// }

	return &ResetRankingResponse{
		RankingID:      req.RankingID,
		EntriesCleared: 0, // TODO: entriesCleared,
		ResetAt:        time.Now(),
	}, nil
}

// 私有方法

// validateCreateRankingRequest 验证创建排行榜请求
func (s *RankingApplicationService) validateCreateRankingRequest(req *CreateRankingRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(req.Name) > 100 {
		return fmt.Errorf("name too long (max 100 characters)")
	}
	if req.RankType == "" {
		return fmt.Errorf("rank type is required")
	}
	if req.PeriodType == "" {
		return fmt.Errorf("period type is required")
	}
	if req.MaxEntries <= 0 {
		return fmt.Errorf("max entries must be positive")
	}
	if req.MaxEntries > 10000 {
		return fmt.Errorf("max entries too large (max 10000)")
	}
	return nil
}

// validateUpdatePlayerScoreRequest 验证更新玩家分数请求
func (s *RankingApplicationService) validateUpdatePlayerScoreRequest(req *UpdatePlayerScoreRequest) error {
	if req.RankingID == "" {
		return fmt.Errorf("ranking ID is required")
	}
	if req.PlayerID == 0 {
		return fmt.Errorf("player ID is required")
	}
	return nil
}

// parseRankType 解析排行榜类型
func (s *RankingApplicationService) parseRankType(rankTypeStr string) (ranking.RankType, error) {
	switch rankTypeStr {
	case "level":
		return ranking.RankTypeLevel, nil
	case "exp":
		return ranking.RankTypeLevel, nil // TODO: 修复RankTypeExp
	case "power":
		return ranking.RankTypePower, nil
	case "wealth":
		return ranking.RankTypeWealth, nil
	case "achievement":
		return ranking.RankTypeAchievement, nil
	case "pvp":
		return ranking.RankTypeLevel, nil // TODO: 修复RankTypePvP
	case "guild":
		return ranking.RankTypeGuild, nil
	default:
		return ranking.RankTypeLevel, fmt.Errorf("unknown rank type: %s", rankTypeStr)
	}
}

// parsePeriodType 解析周期类型
func (s *RankingApplicationService) parsePeriodType(periodTypeStr string) (ranking.RankPeriod, error) {
	switch periodTypeStr {
	case "permanent":
		return ranking.RankPeriodPermanent, nil
	case "daily":
		return ranking.RankPeriodDaily, nil
	case "weekly":
		return ranking.RankPeriodWeekly, nil
	case "monthly":
		return ranking.RankPeriodMonthly, nil
	case "seasonal":
		return ranking.RankPeriodSeasonal, nil
	default:
		return ranking.RankPeriodPermanent, fmt.Errorf("unknown period type: %s", periodTypeStr)
	}
}
