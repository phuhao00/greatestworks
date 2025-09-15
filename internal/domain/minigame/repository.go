package minigame

import (
	"context"
	"time"
)

// MinigameRepository 小游戏仓储接口
type MinigameRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, minigame *MinigameAggregate) error
	FindByID(ctx context.Context, id string) (*MinigameAggregate, error)
	FindByIDs(ctx context.Context, ids []string) ([]*MinigameAggregate, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
	
	// 查询操作
	FindByQuery(ctx context.Context, query *GameQuery) ([]*MinigameAggregate, error)
	FindByCreator(ctx context.Context, creatorID uint64, limit, offset int32) ([]*MinigameAggregate, error)
	FindByType(ctx context.Context, gameType GameType, limit, offset int32) ([]*MinigameAggregate, error)
	FindByStatus(ctx context.Context, status GameStatus, limit, offset int32) ([]*MinigameAggregate, error)
	FindByPlayer(ctx context.Context, playerID uint64, limit, offset int32) ([]*MinigameAggregate, error)
	FindActive(ctx context.Context, limit, offset int32) ([]*MinigameAggregate, error)
	FindJoinable(ctx context.Context, limit, offset int32) ([]*MinigameAggregate, error)
	
	// 分页查询
	FindWithPagination(ctx context.Context, filter *GameFilter, limit, offset int32) (*MinigamePaginationResult, error)
	
	// 统计操作
	Count(ctx context.Context, filter *GameFilter) (int64, error)
	CountByCreator(ctx context.Context, creatorID uint64) (int64, error)
	CountByType(ctx context.Context, gameType GameType) (int64, error)
	CountByStatus(ctx context.Context, status GameStatus) (int64, error)
	GetStatistics(ctx context.Context, filter *GameFilter) (*MinigameStatistics, error)
	
	// 批量操作
	SaveBatch(ctx context.Context, minigames []*MinigameAggregate) error
	DeleteBatch(ctx context.Context, ids []string) error
	UpdateStatusBatch(ctx context.Context, ids []string, status GameStatus) error
	
	// 清理操作
	CleanupExpired(ctx context.Context, expiredBefore time.Time) (int64, error)
	CleanupFinished(ctx context.Context, finishedBefore time.Time) (int64, error)
}

// GameSessionRepository 游戏会话仓储接口
type GameSessionRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, session *GameSession) error
	FindByID(ctx context.Context, id string) (*GameSession, error)
	FindByIDs(ctx context.Context, ids []string) ([]*GameSession, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
	
	// 查询操作
	FindByQuery(ctx context.Context, query *GameSessionQuery) ([]*GameSession, error)
	FindByGame(ctx context.Context, gameID string) ([]*GameSession, error)
	FindByPlayer(ctx context.Context, playerID uint64) ([]*GameSession, error)
	FindByGameAndPlayer(ctx context.Context, gameID string, playerID uint64) (*GameSession, error)
	FindByStatus(ctx context.Context, status PlayerStatus, limit, offset int32) ([]*GameSession, error)
	FindActive(ctx context.Context, limit, offset int32) ([]*GameSession, error)
	FindByToken(ctx context.Context, sessionToken string) (*GameSession, error)
	
	// 分页查询
	FindWithPagination(ctx context.Context, query *GameSessionQuery, limit, offset int32) (*SessionPaginationResult, error)
	
	// 统计操作
	Count(ctx context.Context, query *GameSessionQuery) (int64, error)
	CountByGame(ctx context.Context, gameID string) (int64, error)
	CountByPlayer(ctx context.Context, playerID uint64) (int64, error)
	CountByStatus(ctx context.Context, status PlayerStatus) (int64, error)
	GetPlayerStatistics(ctx context.Context, playerID uint64, gameType *GameType) (*PlayerStatistics, error)
	GetGameStatistics(ctx context.Context, gameID string) (*GameSessionStatistics, error)
	
	// 批量操作
	SaveBatch(ctx context.Context, sessions []*GameSession) error
	DeleteBatch(ctx context.Context, ids []string) error
	UpdateStatusBatch(ctx context.Context, ids []string, status PlayerStatus) error
	
	// 清理操作
	CleanupExpired(ctx context.Context, expiredBefore time.Time) (int64, error)
	CleanupInactive(ctx context.Context, inactiveBefore time.Time) (int64, error)
}

// GameScoreRepository 游戏分数仓储接口
type GameScoreRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, score *GameScore) error
	FindByID(ctx context.Context, id string) (*GameScore, error)
	FindByIDs(ctx context.Context, ids []string) ([]*GameScore, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
	
	// 查询操作
	FindByGame(ctx context.Context, gameID string) ([]*GameScore, error)
	FindByPlayer(ctx context.Context, playerID uint64) ([]*GameScore, error)
	FindByGameAndPlayer(ctx context.Context, gameID string, playerID uint64) ([]*GameScore, error)
	FindByGamePlayerAndType(ctx context.Context, gameID string, playerID uint64, scoreType ScoreType) (*GameScore, error)
	FindByType(ctx context.Context, scoreType ScoreType, limit, offset int32) ([]*GameScore, error)
	FindBySession(ctx context.Context, sessionID string) ([]*GameScore, error)
	
	// 排行榜操作
	FindTopScores(ctx context.Context, scoreType ScoreType, limit int32) ([]*GameScore, error)
	FindTopScoresByGame(ctx context.Context, gameID string, scoreType ScoreType, limit int32) ([]*GameScore, error)
	FindTopScoresByPlayer(ctx context.Context, playerID uint64, scoreType ScoreType, limit int32) ([]*GameScore, error)
	GetPlayerRank(ctx context.Context, gameID string, playerID uint64, scoreType ScoreType) (int32, error)
	GetScorePercentile(ctx context.Context, gameID string, score int64, scoreType ScoreType) (float64, error)
	
	// 分页查询
	FindWithPagination(ctx context.Context, query *ScoreQuery, limit, offset int32) (*ScorePaginationResult, error)
	
	// 统计操作
	Count(ctx context.Context, query *ScoreQuery) (int64, error)
	CountByGame(ctx context.Context, gameID string) (int64, error)
	CountByPlayer(ctx context.Context, playerID uint64) (int64, error)
	GetScoreStatistics(ctx context.Context, gameID string, scoreType ScoreType) (*ScoreStatistics, error)
	GetPlayerScoreStatistics(ctx context.Context, playerID uint64, scoreType ScoreType) (*PlayerScoreStatistics, error)
	
	// 批量操作
	SaveBatch(ctx context.Context, scores []*GameScore) error
	DeleteBatch(ctx context.Context, ids []string) error
	UpdateRanksBatch(ctx context.Context, gameID string, scoreType ScoreType) error
	
	// 清理操作
	CleanupOldScores(ctx context.Context, olderThan time.Time) (int64, error)
}

// GameRewardRepository 游戏奖励仓储接口
type GameRewardRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, reward *GameReward) error
	FindByID(ctx context.Context, id string) (*GameReward, error)
	FindByIDs(ctx context.Context, ids []string) ([]*GameReward, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
	
	// 查询操作
	FindByGame(ctx context.Context, gameID string) ([]*GameReward, error)
	FindByPlayer(ctx context.Context, playerID uint64) ([]*GameReward, error)
	FindByGameAndPlayer(ctx context.Context, gameID string, playerID uint64) ([]*GameReward, error)
	FindByType(ctx context.Context, rewardType RewardType, limit, offset int32) ([]*GameReward, error)
	FindBySession(ctx context.Context, sessionID string) ([]*GameReward, error)
	FindClaimable(ctx context.Context, playerID uint64) ([]*GameReward, error)
	FindClaimed(ctx context.Context, playerID uint64, limit, offset int32) ([]*GameReward, error)
	FindExpired(ctx context.Context, expiredBefore time.Time) ([]*GameReward, error)
	
	// 分页查询
	FindWithPagination(ctx context.Context, query *RewardQuery, limit, offset int32) (*RewardPaginationResult, error)
	
	// 统计操作
	Count(ctx context.Context, query *RewardQuery) (int64, error)
	CountByGame(ctx context.Context, gameID string) (int64, error)
	CountByPlayer(ctx context.Context, playerID uint64) (int64, error)
	CountClaimable(ctx context.Context, playerID uint64) (int64, error)
	CountClaimed(ctx context.Context, playerID uint64) (int64, error)
	GetRewardStatistics(ctx context.Context, gameID string) (*RewardStatistics, error)
	GetPlayerRewardStatistics(ctx context.Context, playerID uint64) (*PlayerRewardStatistics, error)
	
	// 批量操作
	SaveBatch(ctx context.Context, rewards []*GameReward) error
	DeleteBatch(ctx context.Context, ids []string) error
	ClaimBatch(ctx context.Context, ids []string, playerID uint64) error
	
	// 清理操作
	CleanupExpired(ctx context.Context, expiredBefore time.Time) (int64, error)
	CleanupClaimed(ctx context.Context, claimedBefore time.Time) (int64, error)
}

// GameAchievementRepository 游戏成就仓储接口
type GameAchievementRepository interface {
	// 基础CRUD操作
	Save(ctx context.Context, achievement *GameAchievement) error
	FindByID(ctx context.Context, id string) (*GameAchievement, error)
	FindByIDs(ctx context.Context, ids []string) ([]*GameAchievement, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
	
	// 查询操作
	FindByGame(ctx context.Context, gameID string) ([]*GameAchievement, error)
	FindByPlayer(ctx context.Context, playerID uint64) ([]*GameAchievement, error)
	FindByGameAndPlayer(ctx context.Context, gameID string, playerID uint64) ([]*GameAchievement, error)
	FindByAchievementID(ctx context.Context, achievementID string) ([]*GameAchievement, error)
	FindByCategory(ctx context.Context, category string, limit, offset int32) ([]*GameAchievement, error)
	FindByRarity(ctx context.Context, rarity string, limit, offset int32) ([]*GameAchievement, error)
	FindCompleted(ctx context.Context, playerID uint64) ([]*GameAchievement, error)
	FindInProgress(ctx context.Context, playerID uint64) ([]*GameAchievement, error)
	FindUnlocked(ctx context.Context, playerID uint64) ([]*GameAchievement, error)
	
	// 分页查询
	FindWithPagination(ctx context.Context, query *AchievementQuery, limit, offset int32) (*AchievementPaginationResult, error)
	
	// 统计操作
	Count(ctx context.Context, query *AchievementQuery) (int64, error)
	CountByGame(ctx context.Context, gameID string) (int64, error)
	CountByPlayer(ctx context.Context, playerID uint64) (int64, error)
	CountCompleted(ctx context.Context, playerID uint64) (int64, error)
	CountUnlocked(ctx context.Context, playerID uint64) (int64, error)
	GetAchievementStatistics(ctx context.Context, gameID string) (*AchievementStatistics, error)
	GetPlayerAchievementStatistics(ctx context.Context, playerID uint64) (*PlayerAchievementStatistics, error)
	
	// 批量操作
	SaveBatch(ctx context.Context, achievements []*GameAchievement) error
	DeleteBatch(ctx context.Context, ids []string) error
	CompleteBatch(ctx context.Context, ids []string) error
	UnlockBatch(ctx context.Context, ids []string) error
	
	// 清理操作
	CleanupOldAchievements(ctx context.Context, olderThan time.Time) (int64, error)
}

// 查询条件结构体

// ScoreQuery 分数查询条件
type ScoreQuery struct {
	GameID       *string     `json:"game_id,omitempty"`
	PlayerID     *uint64     `json:"player_id,omitempty"`
	SessionID    *string     `json:"session_id,omitempty"`
	ScoreType    *ScoreType  `json:"score_type,omitempty"`
	MinValue     *int64      `json:"min_value,omitempty"`
	MaxValue     *int64      `json:"max_value,omitempty"`
	MinRank      *int32      `json:"min_rank,omitempty"`
	MaxRank      *int32      `json:"max_rank,omitempty"`
	AchievedAfter *time.Time `json:"achieved_after,omitempty"`
	AchievedBefore *time.Time `json:"achieved_before,omitempty"`
	OrderBy      string      `json:"order_by,omitempty"`
	OrderDesc    bool        `json:"order_desc,omitempty"`
}

// RewardQuery 奖励查询条件
type RewardQuery struct {
	GameID       *string      `json:"game_id,omitempty"`
	PlayerID     *uint64      `json:"player_id,omitempty"`
	SessionID    *string      `json:"session_id,omitempty"`
	RewardType   *RewardType  `json:"reward_type,omitempty"`
	ItemID       *string      `json:"item_id,omitempty"`
	Rarity       *string      `json:"rarity,omitempty"`
	Source       *string      `json:"source,omitempty"`
	Claimed      *bool        `json:"claimed,omitempty"`
	Expired      *bool        `json:"expired,omitempty"`
	CreatedAfter *time.Time   `json:"created_after,omitempty"`
	CreatedBefore *time.Time  `json:"created_before,omitempty"`
	ExpiresAfter *time.Time   `json:"expires_after,omitempty"`
	ExpiresBefore *time.Time  `json:"expires_before,omitempty"`
	OrderBy      string       `json:"order_by,omitempty"`
	OrderDesc    bool         `json:"order_desc,omitempty"`
}

// AchievementQuery 成就查询条件
type AchievementQuery struct {
	GameID        *string    `json:"game_id,omitempty"`
	PlayerID      *uint64    `json:"player_id,omitempty"`
	SessionID     *string    `json:"session_id,omitempty"`
	AchievementID *string    `json:"achievement_id,omitempty"`
	Category      *string    `json:"category,omitempty"`
	Rarity        *string    `json:"rarity,omitempty"`
	Completed     *bool      `json:"completed,omitempty"`
	Unlocked      *bool      `json:"unlocked,omitempty"`
	MinProgress   *float64   `json:"min_progress,omitempty"`
	MaxProgress   *float64   `json:"max_progress,omitempty"`
	MinPoints     *int64     `json:"min_points,omitempty"`
	MaxPoints     *int64     `json:"max_points,omitempty"`
	CreatedAfter  *time.Time `json:"created_after,omitempty"`
	CreatedBefore *time.Time `json:"created_before,omitempty"`
	CompletedAfter *time.Time `json:"completed_after,omitempty"`
	CompletedBefore *time.Time `json:"completed_before,omitempty"`
	OrderBy       string     `json:"order_by,omitempty"`
	OrderDesc     bool       `json:"order_desc,omitempty"`
}

// 分页结果结构体

// MinigamePaginationResult 小游戏分页结果
type MinigamePaginationResult struct {
	Items      []*MinigameAggregate `json:"items"`
	Total      int64                `json:"total"`
	Limit      int32                `json:"limit"`
	Offset     int32                `json:"offset"`
	HasMore    bool                 `json:"has_more"`
	TotalPages int32                `json:"total_pages"`
}

// SessionPaginationResult 会话分页结果
type SessionPaginationResult struct {
	Items      []*GameSession `json:"items"`
	Total      int64          `json:"total"`
	Limit      int32          `json:"limit"`
	Offset     int32          `json:"offset"`
	HasMore    bool           `json:"has_more"`
	TotalPages int32          `json:"total_pages"`
}

// ScorePaginationResult 分数分页结果
type ScorePaginationResult struct {
	Items      []*GameScore `json:"items"`
	Total      int64        `json:"total"`
	Limit      int32        `json:"limit"`
	Offset     int32        `json:"offset"`
	HasMore    bool         `json:"has_more"`
	TotalPages int32        `json:"total_pages"`
}

// RewardPaginationResult 奖励分页结果
type RewardPaginationResult struct {
	Items      []*GameReward `json:"items"`
	Total      int64         `json:"total"`
	Limit      int32         `json:"limit"`
	Offset     int32         `json:"offset"`
	HasMore    bool          `json:"has_more"`
	TotalPages int32         `json:"total_pages"`
}

// AchievementPaginationResult 成就分页结果
type AchievementPaginationResult struct {
	Items      []*GameAchievement `json:"items"`
	Total      int64              `json:"total"`
	Limit      int32              `json:"limit"`
	Offset     int32              `json:"offset"`
	HasMore    bool               `json:"has_more"`
	TotalPages int32              `json:"total_pages"`
}

// 统计数据结构体

// MinigameStatistics 小游戏统计
type MinigameStatistics struct {
	TotalGames       int64                    `json:"total_games"`
	ActiveGames      int64                    `json:"active_games"`
	FinishedGames    int64                    `json:"finished_games"`
	CancelledGames   int64                    `json:"cancelled_games"`
	TotalPlayers     int64                    `json:"total_players"`
	ActivePlayers    int64                    `json:"active_players"`
	AveragePlayTime  time.Duration            `json:"average_play_time"`
	AverageScore     float64                  `json:"average_score"`
	GamesByType      map[GameType]int64       `json:"games_by_type"`
	GamesByStatus    map[GameStatus]int64     `json:"games_by_status"`
	GamesByDifficulty map[GameDifficulty]int64 `json:"games_by_difficulty"`
	CreatedAt        time.Time                `json:"created_at"`
	UpdatedAt        time.Time                `json:"updated_at"`
}

// GameSessionStatistics 游戏会话统计
type GameSessionStatistics struct {
	GameID           string        `json:"game_id"`
	TotalSessions    int64         `json:"total_sessions"`
	ActiveSessions   int64         `json:"active_sessions"`
	FinishedSessions int64         `json:"finished_sessions"`
	AveragePlayTime  time.Duration `json:"average_play_time"`
	AverageScore     float64       `json:"average_score"`
	AverageMoves     float64       `json:"average_moves"`
	AverageLevel     float64       `json:"average_level"`
	HighestScore     int64         `json:"highest_score"`
	HighestLevel     int32         `json:"highest_level"`
	TotalPlayTime    time.Duration `json:"total_play_time"`
	TotalMoves       int64         `json:"total_moves"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

// ScoreStatistics 分数统计
type ScoreStatistics struct {
	GameID       string    `json:"game_id"`
	ScoreType    ScoreType `json:"score_type"`
	TotalScores  int64     `json:"total_scores"`
	AverageScore float64   `json:"average_score"`
	HighestScore int64     `json:"highest_score"`
	LowestScore  int64     `json:"lowest_score"`
	MedianScore  int64     `json:"median_score"`
	TotalValue   int64     `json:"total_value"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PlayerScoreStatistics 玩家分数统计
type PlayerScoreStatistics struct {
	PlayerID     uint64    `json:"player_id"`
	ScoreType    ScoreType `json:"score_type"`
	TotalScores  int64     `json:"total_scores"`
	AverageScore float64   `json:"average_score"`
	HighestScore int64     `json:"highest_score"`
	LowestScore  int64     `json:"lowest_score"`
	BestRank     int32     `json:"best_rank"`
	CurrentRank  int32     `json:"current_rank"`
	TotalValue   int64     `json:"total_value"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RewardStatistics 奖励统计
type RewardStatistics struct {
	GameID         string                 `json:"game_id"`
	TotalRewards   int64                  `json:"total_rewards"`
	ClaimedRewards int64                  `json:"claimed_rewards"`
	ExpiredRewards int64                  `json:"expired_rewards"`
	PendingRewards int64                  `json:"pending_rewards"`
	RewardsByType  map[RewardType]int64   `json:"rewards_by_type"`
	RewardsByRarity map[string]int64      `json:"rewards_by_rarity"`
	TotalValue     map[RewardType]int64   `json:"total_value"`
	ClaimRate      float64                `json:"claim_rate"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// PlayerRewardStatistics 玩家奖励统计
type PlayerRewardStatistics struct {
	PlayerID       uint64                 `json:"player_id"`
	TotalRewards   int64                  `json:"total_rewards"`
	ClaimedRewards int64                  `json:"claimed_rewards"`
	ExpiredRewards int64                  `json:"expired_rewards"`
	PendingRewards int64                  `json:"pending_rewards"`
	RewardsByType  map[RewardType]int64   `json:"rewards_by_type"`
	RewardsByRarity map[string]int64      `json:"rewards_by_rarity"`
	TotalValue     map[RewardType]int64   `json:"total_value"`
	ClaimRate      float64                `json:"claim_rate"`
	LastClaimedAt  *time.Time             `json:"last_claimed_at,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// AchievementStatistics 成就统计
type AchievementStatistics struct {
	GameID              string            `json:"game_id"`
	TotalAchievements   int64             `json:"total_achievements"`
	CompletedAchievements int64           `json:"completed_achievements"`
	UnlockedAchievements int64            `json:"unlocked_achievements"`
	InProgressAchievements int64          `json:"in_progress_achievements"`
	AchievementsByCategory map[string]int64 `json:"achievements_by_category"`
	AchievementsByRarity map[string]int64 `json:"achievements_by_rarity"`
	AverageProgress     float64           `json:"average_progress"`
	CompletionRate      float64           `json:"completion_rate"`
	TotalPoints         int64             `json:"total_points"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
}

// PlayerAchievementStatistics 玩家成就统计
type PlayerAchievementStatistics struct {
	PlayerID            uint64            `json:"player_id"`
	TotalAchievements   int64             `json:"total_achievements"`
	CompletedAchievements int64           `json:"completed_achievements"`
	UnlockedAchievements int64            `json:"unlocked_achievements"`
	InProgressAchievements int64          `json:"in_progress_achievements"`
	AchievementsByCategory map[string]int64 `json:"achievements_by_category"`
	AchievementsByRarity map[string]int64 `json:"achievements_by_rarity"`
	AverageProgress     float64           `json:"average_progress"`
	CompletionRate      float64           `json:"completion_rate"`
	TotalPoints         int64             `json:"total_points"`
	LastCompletedAt     *time.Time        `json:"last_completed_at,omitempty"`
	LastUnlockedAt      *time.Time        `json:"last_unlocked_at,omitempty"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
}

// 辅助函数

// NewMinigamePaginationResult 创建小游戏分页结果
func NewMinigamePaginationResult(items []*MinigameAggregate, total int64, limit, offset int32) *MinigamePaginationResult {
	totalPages := int32((total + int64(limit) - 1) / int64(limit))
	hasMore := offset+limit < int32(total)
	
	return &MinigamePaginationResult{
		Items:      items,
		Total:      total,
		Limit:      limit,
		Offset:     offset,
		HasMore:    hasMore,
		TotalPages: totalPages,
	}
}

// NewSessionPaginationResult 创建会话分页结果
func NewSessionPaginationResult(items []*GameSession, total int64, limit, offset int32) *SessionPaginationResult {
	totalPages := int32((total + int64(limit) - 1) / int64(limit))
	hasMore := offset+limit < int32(total)
	
	return &SessionPaginationResult{
		Items:      items,
		Total:      total,
		Limit:      limit,
		Offset:     offset,
		HasMore:    hasMore,
		TotalPages: totalPages,
	}
}

// NewScorePaginationResult 创建分数分页结果
func NewScorePaginationResult(items []*GameScore, total int64, limit, offset int32) *ScorePaginationResult {
	totalPages := int32((total + int64(limit) - 1) / int64(limit))
	hasMore := offset+limit < int32(total)
	
	return &ScorePaginationResult{
		Items:      items,
		Total:      total,
		Limit:      limit,
		Offset:     offset,
		HasMore:    hasMore,
		TotalPages: totalPages,
	}
}

// NewRewardPaginationResult 创建奖励分页结果
func NewRewardPaginationResult(items []*GameReward, total int64, limit, offset int32) *RewardPaginationResult {
	totalPages := int32((total + int64(limit) - 1) / int64(limit))
	hasMore := offset+limit < int32(total)
	
	return &RewardPaginationResult{
		Items:      items,
		Total:      total,
		Limit:      limit,
		Offset:     offset,
		HasMore:    hasMore,
		TotalPages: totalPages,
	}
}

// NewAchievementPaginationResult 创建成就分页结果
func NewAchievementPaginationResult(items []*GameAchievement, total int64, limit, offset int32) *AchievementPaginationResult {
	totalPages := int32((total + int64(limit) - 1) / int64(limit))
	hasMore := offset+limit < int32(total)
	
	return &AchievementPaginationResult{
		Items:      items,
		Total:      total,
		Limit:      limit,
		Offset:     offset,
		HasMore:    hasMore,
		TotalPages: totalPages,
	}
}

// NewMinigameStatistics 创建小游戏统计
func NewMinigameStatistics() *MinigameStatistics {
	now := time.Now()
	return &MinigameStatistics{
		TotalGames:        0,
		ActiveGames:       0,
		FinishedGames:     0,
		CancelledGames:    0,
		TotalPlayers:      0,
		ActivePlayers:     0,
		AveragePlayTime:   0,
		AverageScore:      0.0,
		GamesByType:       make(map[GameType]int64),
		GamesByStatus:     make(map[GameStatus]int64),
		GamesByDifficulty: make(map[GameDifficulty]int64),
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// NewGameSessionStatistics 创建游戏会话统计
func NewGameSessionStatistics(gameID string) *GameSessionStatistics {
	now := time.Now()
	return &GameSessionStatistics{
		GameID:           gameID,
		TotalSessions:    0,
		ActiveSessions:   0,
		FinishedSessions: 0,
		AveragePlayTime:  0,
		AverageScore:     0.0,
		AverageMoves:     0.0,
		AverageLevel:     0.0,
		HighestScore:     0,
		HighestLevel:     0,
		TotalPlayTime:    0,
		TotalMoves:       0,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// NewScoreStatistics 创建分数统计
func NewScoreStatistics(gameID string, scoreType ScoreType) *ScoreStatistics {
	now := time.Now()
	return &ScoreStatistics{
		GameID:       gameID,
		ScoreType:    scoreType,
		TotalScores:  0,
		AverageScore: 0.0,
		HighestScore: 0,
		LowestScore:  0,
		MedianScore:  0,
		TotalValue:   0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// NewPlayerScoreStatistics 创建玩家分数统计
func NewPlayerScoreStatistics(playerID uint64, scoreType ScoreType) *PlayerScoreStatistics {
	now := time.Now()
	return &PlayerScoreStatistics{
		PlayerID:     playerID,
		ScoreType:    scoreType,
		TotalScores:  0,
		AverageScore: 0.0,
		HighestScore: 0,
		LowestScore:  0,
		BestRank:     0,
		CurrentRank:  0,
		TotalValue:   0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// NewRewardStatistics 创建奖励统计
func NewRewardStatistics(gameID string) *RewardStatistics {
	now := time.Now()
	return &RewardStatistics{
		GameID:          gameID,
		TotalRewards:    0,
		ClaimedRewards:  0,
		ExpiredRewards:  0,
		PendingRewards:  0,
		RewardsByType:   make(map[RewardType]int64),
		RewardsByRarity: make(map[string]int64),
		TotalValue:      make(map[RewardType]int64),
		ClaimRate:       0.0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// NewPlayerRewardStatistics 创建玩家奖励统计
func NewPlayerRewardStatistics(playerID uint64) *PlayerRewardStatistics {
	now := time.Now()
	return &PlayerRewardStatistics{
		PlayerID:        playerID,
		TotalRewards:    0,
		ClaimedRewards:  0,
		ExpiredRewards:  0,
		PendingRewards:  0,
		RewardsByType:   make(map[RewardType]int64),
		RewardsByRarity: make(map[string]int64),
		TotalValue:      make(map[RewardType]int64),
		ClaimRate:       0.0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// NewAchievementStatistics 创建成就统计
func NewAchievementStatistics(gameID string) *AchievementStatistics {
	now := time.Now()
	return &AchievementStatistics{
		GameID:                 gameID,
		TotalAchievements:      0,
		CompletedAchievements:  0,
		UnlockedAchievements:   0,
		InProgressAchievements: 0,
		AchievementsByCategory: make(map[string]int64),
		AchievementsByRarity:   make(map[string]int64),
		AverageProgress:        0.0,
		CompletionRate:         0.0,
		TotalPoints:            0,
		CreatedAt:              now,
		UpdatedAt:              now,
	}
}

// NewPlayerAchievementStatistics 创建玩家成就统计
func NewPlayerAchievementStatistics(playerID uint64) *PlayerAchievementStatistics {
	now := time.Now()
	return &PlayerAchievementStatistics{
		PlayerID:               playerID,
		TotalAchievements:      0,
		CompletedAchievements:  0,
		UnlockedAchievements:   0,
		InProgressAchievements: 0,
		AchievementsByCategory: make(map[string]int64),
		AchievementsByRarity:   make(map[string]int64),
		AverageProgress:        0.0,
		CompletionRate:         0.0,
		TotalPoints:            0,
		CreatedAt:              now,
		UpdatedAt:              now,
	}
}