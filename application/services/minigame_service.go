package services

import (
	"context"
	"fmt"
	"time"

	"github.com/greatestworks/internal/domain/minigame"
)

// MinigameApplicationService 小游戏应用服务
type MinigameApplicationService struct {
	minigameRepo    minigame.MinigameRepository
	sessionRepo     minigame.GameSessionRepository
	minigameService *minigame.MinigameService
	eventBus        minigame.MinigameEventBus
}

// NewMinigameApplicationService 创建小游戏应用服务
func NewMinigameApplicationService(
	minigameRepo minigame.MinigameRepository,
	sessionRepo minigame.GameSessionRepository,
	minigameService *minigame.MinigameService,
	eventBus minigame.MinigameEventBus,
) *MinigameApplicationService {
	return &MinigameApplicationService{
		minigameRepo:    minigameRepo,
		sessionRepo:     sessionRepo,
		minigameService: minigameService,
		eventBus:        eventBus,
	}
}

// CreateMinigameRequest 创建小游戏请求
type CreateMinigameRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	GameType    string `json:"game_type"`
	Difficulty  string `json:"difficulty"`
	MaxPlayers  int32  `json:"max_players"`
	TimeLimit   int32  `json:"time_limit"`
	IsActive    bool   `json:"is_active"`
}

// CreateMinigameResponse 创建小游戏响应
type CreateMinigameResponse struct {
	MinigameID  string    `json:"minigame_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	GameType    string    `json:"game_type"`
	Difficulty  string    `json:"difficulty"`
	MaxPlayers  int32     `json:"max_players"`
	TimeLimit   int32     `json:"time_limit"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateMinigame 创建小游戏
func (s *MinigameApplicationService) CreateMinigame(ctx context.Context, req *CreateMinigameRequest) (*CreateMinigameResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	
	if err := s.validateCreateMinigameRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	
	// 转换游戏类型
	gameType, err := s.parseGameType(req.GameType)
	if err != nil {
		return nil, fmt.Errorf("invalid game type: %w", err)
	}
	
	// 转换难度
	difficulty, err := s.parseDifficulty(req.Difficulty)
	if err != nil {
		return nil, fmt.Errorf("invalid difficulty: %w", err)
	}
	
	// 创建小游戏聚合根
	minigameAggregate := minigame.NewMinigameAggregate(req.Name, gameType, difficulty)
	minigameAggregate.SetDescription(req.Description)
	minigameAggregate.SetMaxPlayers(req.MaxPlayers)
	minigameAggregate.SetTimeLimit(req.TimeLimit)
	if req.IsActive {
		minigameAggregate.Activate()
	}
	
	// 保存小游戏
	if err := s.minigameRepo.Save(ctx, minigameAggregate); err != nil {
		return nil, fmt.Errorf("failed to save minigame: %w", err)
	}
	
	// 发布事件
	event := minigame.NewMinigameCreatedEvent(minigameAggregate.GetID(), req.Name, gameType, difficulty)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish minigame created event: %v\n", err)
	}
	
	return &CreateMinigameResponse{
		MinigameID:  minigameAggregate.GetID(),
		Name:        minigameAggregate.GetName(),
		Description: minigameAggregate.GetDescription(),
		GameType:    minigameAggregate.GetGameType().String(),
		Difficulty:  minigameAggregate.GetDifficulty().String(),
		MaxPlayers:  minigameAggregate.GetMaxPlayers(),
		TimeLimit:   minigameAggregate.GetTimeLimit(),
		IsActive:    minigameAggregate.IsActive(),
		CreatedAt:   minigameAggregate.GetCreatedAt(),
	}, nil
}

// StartGameSessionRequest 开始游戏会话请求
type StartGameSessionRequest struct {
	MinigameID string `json:"minigame_id"`
	PlayerID   uint64 `json:"player_id"`
	Settings   map[string]interface{} `json:"settings,omitempty"`
}

// StartGameSessionResponse 开始游戏会话响应
type StartGameSessionResponse struct {
	SessionID   string    `json:"session_id"`
	MinigameID  string    `json:"minigame_id"`
	PlayerID    uint64    `json:"player_id"`
	Status      string    `json:"status"`
	TimeLimit   int32     `json:"time_limit"`
	StartedAt   time.Time `json:"started_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// StartGameSession 开始游戏会话
func (s *MinigameApplicationService) StartGameSession(ctx context.Context, req *StartGameSessionRequest) (*StartGameSessionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	
	if err := s.validateStartGameSessionRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	
	// 获取小游戏
	minigameAggregate, err := s.minigameRepo.FindByID(ctx, req.MinigameID)
	if err != nil {
		return nil, fmt.Errorf("failed to find minigame: %w", err)
	}
	if minigameAggregate == nil {
		return nil, fmt.Errorf("minigame not found")
	}
	
	// 检查游戏是否激活
	if !minigameAggregate.IsActive() {
		return nil, fmt.Errorf("minigame is not active")
	}
	
	// 检查玩家是否有进行中的会话
	activeSession, err := s.sessionRepo.FindActiveByPlayer(ctx, req.PlayerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check active session: %w", err)
	}
	if activeSession != nil {
		return nil, fmt.Errorf("player already has an active game session")
	}
	
	// 创建游戏会话
	session, err := s.minigameService.StartGameSession(ctx, req.MinigameID, req.PlayerID)
	if err != nil {
		return nil, fmt.Errorf("failed to start game session: %w", err)
	}
	
	// 设置游戏设置
	if req.Settings != nil {
		for key, value := range req.Settings {
			session.SetSetting(key, value)
		}
		if err := s.sessionRepo.Save(ctx, session); err != nil {
			return nil, fmt.Errorf("failed to save session settings: %w", err)
		}
	}
	
	// 发布事件
	event := minigame.NewGameSessionStartedEvent(session.GetID(), req.MinigameID, req.PlayerID)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish game session started event: %v\n", err)
	}
	
	return &StartGameSessionResponse{
		SessionID:  session.GetID(),
		MinigameID: session.GetMinigameID(),
		PlayerID:   session.GetPlayerID(),
		Status:     session.GetStatus().String(),
		TimeLimit:  session.GetTimeLimit(),
		StartedAt:  session.GetStartedAt(),
		ExpiresAt:  session.GetExpiresAt(),
	}, nil
}

// SubmitGameScoreRequest 提交游戏分数请求
type SubmitGameScoreRequest struct {
	SessionID string `json:"session_id"`
	Score     int64  `json:"score"`
	GameData  map[string]interface{} `json:"game_data,omitempty"`
}

// SubmitGameScoreResponse 提交游戏分数响应
type SubmitGameScoreResponse struct {
	SessionID    string                 `json:"session_id"`
	Score        int64                  `json:"score"`
	BestScore    int64                  `json:"best_score"`
	IsNewRecord  bool                   `json:"is_new_record"`
	Rewards      []*GameRewardResponse  `json:"rewards,omitempty"`
	Achievements []string               `json:"achievements,omitempty"`
	CompletedAt  time.Time              `json:"completed_at"`
}

// GameRewardResponse 游戏奖励响应
type GameRewardResponse struct {
	Type     string `json:"type"`
	ItemID   string `json:"item_id,omitempty"`
	Quantity int32  `json:"quantity"`
	Reason   string `json:"reason"`
}

// SubmitGameScore 提交游戏分数
func (s *MinigameApplicationService) SubmitGameScore(ctx context.Context, req *SubmitGameScoreRequest) (*SubmitGameScoreResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	
	if err := s.validateSubmitGameScoreRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	
	// 获取游戏会话
	session, err := s.sessionRepo.FindByID(ctx, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, fmt.Errorf("session not found")
	}
	
	// 检查会话状态
	if !session.IsActive() {
		return nil, fmt.Errorf("session is not active")
	}
	
	// 检查会话是否过期
	if session.IsExpired() {
		return nil, fmt.Errorf("session has expired")
	}
	
	// 提交分数并完成会话
	result, err := s.minigameService.SubmitGameScore(ctx, req.SessionID, req.Score)
	if err != nil {
		return nil, fmt.Errorf("failed to submit game score: %w", err)
	}
	
	// 设置游戏数据
	if req.GameData != nil {
		for key, value := range req.GameData {
			session.SetGameData(key, value)
		}
		if err := s.sessionRepo.Save(ctx, session); err != nil {
			return nil, fmt.Errorf("failed to save session game data: %w", err)
		}
	}
	
	// 转换奖励响应
	rewardResponses := make([]*GameRewardResponse, len(result.Rewards))
	for i, reward := range result.Rewards {
		rewardResponses[i] = &GameRewardResponse{
			Type:     reward.GetType().String(),
			ItemID:   reward.GetItemID(),
			Quantity: reward.GetQuantity(),
			Reason:   reward.GetReason(),
		}
	}
	
	// 发布事件
	event := minigame.NewGameScoreSubmittedEvent(
		req.SessionID,
		session.GetMinigameID(),
		session.GetPlayerID(),
		req.Score,
		result.IsNewRecord,
	)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish game score submitted event: %v\n", err)
	}
	
	return &SubmitGameScoreResponse{
		SessionID:    req.SessionID,
		Score:        req.Score,
		BestScore:    result.BestScore,
		IsNewRecord:  result.IsNewRecord,
		Rewards:      rewardResponses,
		Achievements: result.Achievements,
		CompletedAt:  session.GetCompletedAt(),
	}, nil
}

// GetGameSessionRequest 获取游戏会话请求
type GetGameSessionRequest struct {
	SessionID string `json:"session_id"`
}

// GetGameSessionResponse 获取游戏会话响应
type GetGameSessionResponse struct {
	SessionID   string                 `json:"session_id"`
	MinigameID  string                 `json:"minigame_id"`
	PlayerID    uint64                 `json:"player_id"`
	Status      string                 `json:"status"`
	Score       int64                  `json:"score"`
	TimeLimit   int32                  `json:"time_limit"`
	TimeElapsed int32                  `json:"time_elapsed"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	GameData    map[string]interface{} `json:"game_data,omitempty"`
	StartedAt   time.Time              `json:"started_at"`
	ExpiresAt   time.Time              `json:"expires_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
}

// GetGameSession 获取游戏会话
func (s *MinigameApplicationService) GetGameSession(ctx context.Context, req *GetGameSessionRequest) (*GetGameSessionResponse, error) {
	if req == nil || req.SessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}
	
	// 获取游戏会话
	session, err := s.sessionRepo.FindByID(ctx, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, fmt.Errorf("session not found")
	}
	
	// 计算已用时间
	timeElapsed := int32(0)
	if session.IsActive() {
		timeElapsed = int32(time.Since(session.GetStartedAt()).Seconds())
	} else if !session.GetCompletedAt().IsZero() {
		timeElapsed = int32(session.GetCompletedAt().Sub(session.GetStartedAt()).Seconds())
	}
	
	response := &GetGameSessionResponse{
		SessionID:   session.GetID(),
		MinigameID:  session.GetMinigameID(),
		PlayerID:    session.GetPlayerID(),
		Status:      session.GetStatus().String(),
		Score:       session.GetScore(),
		TimeLimit:   session.GetTimeLimit(),
		TimeElapsed: timeElapsed,
		Settings:    session.GetSettings(),
		GameData:    session.GetGameData(),
		StartedAt:   session.GetStartedAt(),
		ExpiresAt:   session.GetExpiresAt(),
	}
	
	if !session.GetCompletedAt().IsZero() {
		completedAt := session.GetCompletedAt()
		response.CompletedAt = &completedAt
	}
	
	return response, nil
}

// GetPlayerScoresRequest 获取玩家分数请求
type GetPlayerScoresRequest struct {
	PlayerID   uint64 `json:"player_id"`
	MinigameID string `json:"minigame_id,omitempty"`
	Limit      int    `json:"limit"`
}

// PlayerScoreResponse 玩家分数响应
type PlayerScoreResponse struct {
	SessionID   string    `json:"session_id"`
	MinigameID  string    `json:"minigame_id"`
	Score       int64     `json:"score"`
	Rank        int32     `json:"rank"`
	CompletedAt time.Time `json:"completed_at"`
}

// GetPlayerScoresResponse 获取玩家分数响应
type GetPlayerScoresResponse struct {
	PlayerID uint64                 `json:"player_id"`
	Scores   []*PlayerScoreResponse `json:"scores"`
	Total    int64                  `json:"total"`
}

// GetPlayerScores 获取玩家分数
func (s *MinigameApplicationService) GetPlayerScores(ctx context.Context, req *GetPlayerScoresRequest) (*GetPlayerScoresResponse, error) {
	if req == nil || req.PlayerID == 0 {
		return nil, fmt.Errorf("player ID is required")
	}
	
	// 设置默认值
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	
	// 构建查询
	query := minigame.NewGameSessionQuery().
		WithPlayer(req.PlayerID).
		WithStatus(minigame.SessionStatusCompleted).
		WithSort("completed_at", "desc").
		WithLimit(req.Limit)
	
	if req.MinigameID != "" {
		query = query.WithMinigame(req.MinigameID)
	}
	
	// 查询会话
	sessions, total, err := s.sessionRepo.FindByQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find sessions: %w", err)
	}
	
	// 转换响应
	scoreResponses := make([]*PlayerScoreResponse, len(sessions))
	for i, session := range sessions {
		// 获取排名（这里简化处理，实际可能需要更复杂的排名计算）
		rank, _ := s.minigameService.GetPlayerRank(ctx, session.GetMinigameID(), req.PlayerID)
		
		scoreResponses[i] = &PlayerScoreResponse{
			SessionID:   session.GetID(),
			MinigameID:  session.GetMinigameID(),
			Score:       session.GetScore(),
			Rank:        rank,
			CompletedAt: session.GetCompletedAt(),
		}
	}
	
	return &GetPlayerScoresResponse{
		PlayerID: req.PlayerID,
		Scores:   scoreResponses,
		Total:    total,
	}, nil
}

// GetMinigameLeaderboardRequest 获取小游戏排行榜请求
type GetMinigameLeaderboardRequest struct {
	MinigameID string `json:"minigame_id"`
	Period     string `json:"period"`
	Limit      int    `json:"limit"`
}

// LeaderboardEntryResponse 排行榜条目响应
type LeaderboardEntryResponse struct {
	PlayerID    uint64    `json:"player_id"`
	Rank        int32     `json:"rank"`
	Score       int64     `json:"score"`
	SessionID   string    `json:"session_id"`
	CompletedAt time.Time `json:"completed_at"`
}

// GetMinigameLeaderboardResponse 获取小游戏排行榜响应
type GetMinigameLeaderboardResponse struct {
	MinigameID string                      `json:"minigame_id"`
	Period     string                      `json:"period"`
	Entries    []*LeaderboardEntryResponse `json:"entries"`
	UpdatedAt  time.Time                   `json:"updated_at"`
}

// GetMinigameLeaderboard 获取小游戏排行榜
func (s *MinigameApplicationService) GetMinigameLeaderboard(ctx context.Context, req *GetMinigameLeaderboardRequest) (*GetMinigameLeaderboardResponse, error) {
	if req == nil || req.MinigameID == "" {
		return nil, fmt.Errorf("minigame ID is required")
	}
	
	// 设置默认值
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.Period == "" {
		req.Period = "all_time"
	}
	
	// 获取排行榜
	leaderboard, err := s.minigameService.GetLeaderboard(ctx, req.MinigameID, req.Period, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	
	// 转换响应
	entryResponses := make([]*LeaderboardEntryResponse, len(leaderboard.Entries))
	for i, entry := range leaderboard.Entries {
		entryResponses[i] = &LeaderboardEntryResponse{
			PlayerID:    entry.PlayerID,
			Rank:        entry.Rank,
			Score:       entry.Score,
			SessionID:   entry.SessionID,
			CompletedAt: entry.CompletedAt,
		}
	}
	
	return &GetMinigameLeaderboardResponse{
		MinigameID: req.MinigameID,
		Period:     req.Period,
		Entries:    entryResponses,
		UpdatedAt:  leaderboard.UpdatedAt,
	}, nil
}

// 私有方法

// validateCreateMinigameRequest 验证创建小游戏请求
func (s *MinigameApplicationService) validateCreateMinigameRequest(req *CreateMinigameRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(req.Name) > 100 {
		return fmt.Errorf("name too long (max 100 characters)")
	}
	if req.GameType == "" {
		return fmt.Errorf("game type is required")
	}
	if req.Difficulty == "" {
		return fmt.Errorf("difficulty is required")
	}
	if req.MaxPlayers <= 0 {
		return fmt.Errorf("max players must be positive")
	}
	if req.TimeLimit <= 0 {
		return fmt.Errorf("time limit must be positive")
	}
	return nil
}

// validateStartGameSessionRequest 验证开始游戏会话请求
func (s *MinigameApplicationService) validateStartGameSessionRequest(req *StartGameSessionRequest) error {
	if req.MinigameID == "" {
		return fmt.Errorf("minigame ID is required")
	}
	if req.PlayerID == 0 {
		return fmt.Errorf("player ID is required")
	}
	return nil
}

// validateSubmitGameScoreRequest 验证提交游戏分数请求
func (s *MinigameApplicationService) validateSubmitGameScoreRequest(req *SubmitGameScoreRequest) error {
	if req.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if req.Score < 0 {
		return fmt.Errorf("score cannot be negative")
	}
	return nil
}

// parseGameType 解析游戏类型
func (s *MinigameApplicationService) parseGameType(gameTypeStr string) (minigame.GameType, error) {
	switch gameTypeStr {
	case "puzzle":
		return minigame.GameTypePuzzle, nil
	case "action":
		return minigame.GameTypeAction, nil
	case "strategy":
		return minigame.GameTypeStrategy, nil
	case "arcade":
		return minigame.GameTypeArcade, nil
	case "card":
		return minigame.GameTypeCard, nil
	case "quiz":
		return minigame.GameTypeQuiz, nil
	default:
		return minigame.GameTypePuzzle, fmt.Errorf("unknown game type: %s", gameTypeStr)
	}
}

// parseDifficulty 解析难度
func (s *MinigameApplicationService) parseDifficulty(difficultyStr string) (minigame.Difficulty, error) {
	switch difficultyStr {
	case "easy":
		return minigame.DifficultyEasy, nil
	case "normal":
		return minigame.DifficultyNormal, nil
	case "hard":
		return minigame.DifficultyHard, nil
	case "expert":
		return minigame.DifficultyExpert, nil
	default:
		return minigame.DifficultyNormal, fmt.Errorf("unknown difficulty: %s", difficultyStr)
	}
}