package minigame

import (
	"context"
	"fmt"
	"time"
)

// MinigameService 小游戏领域服务
type MinigameService struct {
	minigameRepo    MinigameRepository
	sessionRepo     GameSessionRepository
	scoreRepo       GameScoreRepository
	rewardRepo      GameRewardRepository
	achievementRepo GameAchievementRepository
	eventBus        MinigameEventBus
}

// NewMinigameService 创建小游戏服务
func NewMinigameService(
	minigameRepo MinigameRepository,
	sessionRepo GameSessionRepository,
	scoreRepo GameScoreRepository,
	rewardRepo GameRewardRepository,
	achievementRepo GameAchievementRepository,
	eventBus MinigameEventBus,
) *MinigameService {
	return &MinigameService{
		minigameRepo:    minigameRepo,
		sessionRepo:     sessionRepo,
		scoreRepo:       scoreRepo,
		rewardRepo:      rewardRepo,
		achievementRepo: achievementRepo,
		eventBus:        eventBus,
	}
}

// CreateMinigame 创建小游戏
func (s *MinigameService) CreateMinigame(ctx context.Context, gameType GameType, creatorID uint64, config *GameConfig) (*MinigameAggregate, error) {
	// 验证游戏类型
	if !gameType.IsValid() {
		return nil, NewMinigameError(ErrorCodeInvalidGameType, "Invalid game type", fmt.Errorf("invalid game type: %v", gameType))
	}

	// 验证创建者
	if creatorID == 0 {
		return nil, NewMinigameError(ErrorCodeInvalidPlayer, "Invalid creator", fmt.Errorf("creator_id cannot be zero"))
	}

	// 验证配置
	if config == nil {
		config = NewGameConfig()
	}

	if err := s.validateGameConfig(config); err != nil {
		return nil, NewMinigameError(ErrorCodeInvalidConfig, "Invalid game config", err)
	}

	// 创建小游戏聚合
	minigame := NewMinigameAggregate(gameType, creatorID)
	minigame.SetConfig(config.Clone())

	// 保存到仓储
	if err := s.minigameRepo.Save(ctx, minigame); err != nil {
		return nil, NewMinigameError(ErrorCodeRepositoryError, "Failed to save minigame", err)
	}

	// 发布事件
	event := NewMinigameCreatedEvent(minigame.ID, minigame.GameType, creatorID)
	s.eventBus.Publish(ctx, event)

	return minigame, nil
}

// StartGame 开始游戏
func (s *MinigameService) StartGame(ctx context.Context, gameID string, operatorID uint64) error {
	// 获取游戏
	minigame, err := s.minigameRepo.FindByID(ctx, gameID)
	if err != nil {
		return NewMinigameError(ErrorCodeGameNotFound, "Game not found", err)
	}

	// 检查权限
	if !s.canOperateGame(minigame, operatorID) {
		return NewMinigameError(ErrorCodePermissionDenied, "Permission denied", fmt.Errorf("user %d cannot operate game %s", operatorID, gameID))
	}

	// 开始游戏
	if err := minigame.StartGame(); err != nil {
		return NewMinigameError(ErrorCodeInvalidOperation, "Cannot start game", err)
	}

	// 保存更新
	if err := s.minigameRepo.Save(ctx, minigame); err != nil {
		return NewMinigameError(ErrorCodeRepositoryError, "Failed to save minigame", err)
	}

	// 发布事件
	event := NewGameStartedEvent(gameID, operatorID)
	s.eventBus.Publish(ctx, event)

	return nil
}

// EndGame 结束游戏
func (s *MinigameService) EndGame(ctx context.Context, gameID string, operatorID uint64, reason GameEndReason) error {
	// 获取游戏
	minigame, err := s.minigameRepo.FindByID(ctx, gameID)
	if err != nil {
		return NewMinigameError(ErrorCodeGameNotFound, "Game not found", err)
	}

	// 检查权限
	if !s.canOperateGame(minigame, operatorID) {
		return NewMinigameError(ErrorCodePermissionDenied, "Permission denied", fmt.Errorf("user %d cannot operate game %s", operatorID, gameID))
	}

	// 结束游戏
	if err := minigame.EndGame(reason); err != nil {
		return NewMinigameError(ErrorCodeInvalidOperation, "Cannot end game", err)
	}

	// 处理游戏结束后的逻辑
	if err := s.handleGameEnd(ctx, minigame); err != nil {
		return err
	}

	// 保存更新
	if err := s.minigameRepo.Save(ctx, minigame); err != nil {
		return NewMinigameError(ErrorCodeRepositoryError, "Failed to save minigame", err)
	}

	// 发布事件
	event := NewGameEndedEvent(gameID, reason, operatorID)
	s.eventBus.Publish(ctx, event)

	return nil
}

// JoinGame 加入游戏
func (s *MinigameService) JoinGame(ctx context.Context, gameID string, playerID uint64, sessionToken string) (*GameSession, error) {
	// 获取游戏
	minigame, err := s.minigameRepo.FindByID(ctx, gameID)
	if err != nil {
		return nil, NewMinigameError(ErrorCodeGameNotFound, "Game not found", err)
	}

	// 检查是否可以加入
	if !minigame.CanJoin() {
		return nil, NewMinigameError(ErrorCodeGameNotJoinable, "Game is not joinable", fmt.Errorf("game %s is not joinable", gameID))
	}

	// 检查玩家是否已在游戏中
	if minigame.HasPlayer(playerID) {
		return nil, NewMinigameError(ErrorCodePlayerAlreadyInGame, "Player already in game", fmt.Errorf("player %d already in game %s", playerID, gameID))
	}

	// 检查游戏是否已满
	if minigame.IsFull() {
		return nil, NewMinigameError(ErrorCodeGameFull, "Game is full", fmt.Errorf("game %s is full", gameID))
	}

	// 创建游戏会话
	session := NewGameSession(gameID, playerID, sessionToken)

	// 验证会话
	if err := ValidateGameSession(session); err != nil {
		return nil, NewMinigameError(ErrorCodeInvalidSession, "Invalid session", err)
	}

	// 加入游戏
	if err := minigame.AddPlayer(playerID, session); err != nil {
		return nil, NewMinigameError(ErrorCodeInvalidOperation, "Cannot add player", err)
	}

	// 保存会话
	if err := s.sessionRepo.Save(ctx, session); err != nil {
		return nil, NewMinigameError(ErrorCodeRepositoryError, "Failed to save session", err)
	}

	// 保存游戏更新
	if err := s.minigameRepo.Save(ctx, minigame); err != nil {
		return nil, NewMinigameError(ErrorCodeRepositoryError, "Failed to save minigame", err)
	}

	// 发布事件
	event := NewPlayerJoinedEvent(gameID, playerID, session.ID)
	s.eventBus.Publish(ctx, event)

	return session, nil
}

// LeaveGame 离开游戏
func (s *MinigameService) LeaveGame(ctx context.Context, gameID string, playerID uint64, reason PlayerLeaveReason) error {
	// 获取游戏
	minigame, err := s.minigameRepo.FindByID(ctx, gameID)
	if err != nil {
		return NewMinigameError(ErrorCodeGameNotFound, "Game not found", err)
	}

	// 检查玩家是否在游戏中
	if !minigame.HasPlayer(playerID) {
		return NewMinigameError(ErrorCodePlayerNotInGame, "Player not in game", fmt.Errorf("player %d not in game %s", playerID, gameID))
	}

	// 获取会话
	session, err := s.sessionRepo.FindByGameAndPlayer(ctx, gameID, playerID)
	if err != nil {
		return NewMinigameError(ErrorCodeSessionNotFound, "Session not found", err)
	}

	// 离开游戏
	if err := session.Leave(reason); err != nil {
		return NewMinigameError(ErrorCodeInvalidOperation, "Cannot leave game", err)
	}

	// 从游戏中移除玩家
	if err := minigame.RemovePlayer(playerID, reason); err != nil {
		return NewMinigameError(ErrorCodeInvalidOperation, "Cannot remove player", err)
	}

	// 保存会话更新
	if err := s.sessionRepo.Save(ctx, session); err != nil {
		return NewMinigameError(ErrorCodeRepositoryError, "Failed to save session", err)
	}

	// 保存游戏更新
	if err := s.minigameRepo.Save(ctx, minigame); err != nil {
		return NewMinigameError(ErrorCodeRepositoryError, "Failed to save minigame", err)
	}

	// 发布事件
	event := NewPlayerLeftEvent(gameID, playerID, reason)
	s.eventBus.Publish(ctx, event)

	return nil
}

// UpdatePlayerScore 更新玩家分数
func (s *MinigameService) UpdatePlayerScore(ctx context.Context, gameID string, playerID uint64, scoreType ScoreType, value int64) (*GameScore, error) {
	// 获取游戏
	minigame, err := s.minigameRepo.FindByID(ctx, gameID)
	if err != nil {
		return nil, NewMinigameError(ErrorCodeGameNotFound, "Game not found", err)
	}

	// 检查游戏状态
	if !minigame.IsRunning() {
		return nil, NewMinigameError(ErrorCodeGameNotRunning, "Game is not running", fmt.Errorf("game %s is not running", gameID))
	}

	// 检查玩家是否在游戏中
	if !minigame.HasPlayer(playerID) {
		return nil, NewMinigameError(ErrorCodePlayerNotInGame, "Player not in game", fmt.Errorf("player %d not in game %s", playerID, gameID))
	}

	// 获取会话
	session, err := s.sessionRepo.FindByGameAndPlayer(ctx, gameID, playerID)
	if err != nil {
		return nil, NewMinigameError(ErrorCodeSessionNotFound, "Session not found", err)
	}

	// 更新会话分数
	session.UpdateScore(value)

	// 创建或更新分数记录
	score, err := s.scoreRepo.FindByGamePlayerAndType(ctx, gameID, playerID, scoreType)
	if err != nil {
		// 创建新分数记录
		score = NewGameScore(gameID, playerID, session.ID, scoreType)
	}

	score.UpdateValue(value)

	// 应用难度倍数
	if minigame.Config != nil {
		multiplier := minigame.Config.Difficulty.GetScoreMultiplier()
		score.SetMultiplier(multiplier)
	}

	// 验证分数
	if err := ValidateGameScore(score); err != nil {
		return nil, NewMinigameError(ErrorCodeInvalidScore, "Invalid score", err)
	}

	// 保存分数
	if err := s.scoreRepo.Save(ctx, score); err != nil {
		return nil, NewMinigameError(ErrorCodeRepositoryError, "Failed to save score", err)
	}

	// 保存会话更新
	if err := s.sessionRepo.Save(ctx, session); err != nil {
		return nil, NewMinigameError(ErrorCodeRepositoryError, "Failed to save session", err)
	}

	// 更新游戏统计
	minigame.UpdatePlayerScore(playerID, value)

	// 保存游戏更新
	if err := s.minigameRepo.Save(ctx, minigame); err != nil {
		return nil, NewMinigameError(ErrorCodeRepositoryError, "Failed to save minigame", err)
	}

	// 发布事件
	event := NewScoreUpdatedEvent(gameID, playerID, scoreType, value, score.FinalScore)
	s.eventBus.Publish(ctx, event)

	// 检查成就
	go s.checkAchievements(context.Background(), gameID, playerID, session)

	return score, nil
}

// GrantReward 授予奖励
func (s *MinigameService) GrantReward(ctx context.Context, gameID string, playerID uint64, rewardType RewardType, itemID string, quantity int64, source string) (*GameReward, error) {
	// 获取游戏
	minigame, err := s.minigameRepo.FindByID(ctx, gameID)
	if err != nil {
		return nil, NewMinigameError(ErrorCodeGameNotFound, "Game not found", err)
	}

	// 检查玩家是否参与过游戏
	if !minigame.HasPlayer(playerID) {
		return nil, NewMinigameError(ErrorCodePlayerNotInGame, "Player not in game", fmt.Errorf("player %d not in game %s", playerID, gameID))
	}

	// 获取会话
	session, err := s.sessionRepo.FindByGameAndPlayer(ctx, gameID, playerID)
	if err != nil {
		return nil, NewMinigameError(ErrorCodeSessionNotFound, "Session not found", err)
	}

	// 创建奖励
	reward := NewGameReward(gameID, playerID, session.ID, rewardType, itemID, quantity)
	reward.SetSource(source)
	reward.SetReason(fmt.Sprintf("Game reward from %s", source))

	// 设置过期时间
	expiresAt := time.Now().Add(DefaultRewardTTL)
	reward.SetExpiration(expiresAt)

	// 验证奖励
	if err := ValidateGameReward(reward); err != nil {
		return nil, NewMinigameError(ErrorCodeInvalidReward, "Invalid reward", err)
	}

	// 保存奖励
	if err := s.rewardRepo.Save(ctx, reward); err != nil {
		return nil, NewMinigameError(ErrorCodeRepositoryError, "Failed to save reward", err)
	}

	// 发布事件
	event := NewRewardGrantedEvent(gameID, playerID, rewardType, itemID, quantity)
	s.eventBus.Publish(ctx, event)

	return reward, nil
}

// ClaimReward 领取奖励
func (s *MinigameService) ClaimReward(ctx context.Context, rewardID string, playerID uint64) error {
	// 获取奖励
	reward, err := s.rewardRepo.FindByID(ctx, rewardID)
	if err != nil {
		return NewMinigameError(ErrorCodeRewardNotFound, "Reward not found", err)
	}

	// 检查玩家权限
	if reward.PlayerID != playerID {
		return NewMinigameError(ErrorCodePermissionDenied, "Permission denied", fmt.Errorf("player %d cannot claim reward %s", playerID, rewardID))
	}

	// 领取奖励
	if err := reward.Claim(); err != nil {
		return NewMinigameError(ErrorCodeInvalidOperation, "Cannot claim reward", err)
	}

	// 保存更新
	if err := s.rewardRepo.Save(ctx, reward); err != nil {
		return NewMinigameError(ErrorCodeRepositoryError, "Failed to save reward", err)
	}

	// 发布事件
	event := NewRewardClaimedEvent(reward.GameID, playerID, reward.RewardType, reward.ItemID, reward.Quantity)
	s.eventBus.Publish(ctx, event)

	return nil
}

// GetGameLeaderboard 获取游戏排行榜
func (s *MinigameService) GetGameLeaderboard(ctx context.Context, gameID string, scoreType ScoreType, limit int32) ([]*GameScore, error) {
	// 获取游戏
	minigame, err := s.minigameRepo.FindByID(ctx, gameID)
	if err != nil {
		return nil, NewMinigameError(ErrorCodeGameNotFound, "Game not found", err)
	}

	// 获取排行榜
	scores, err := s.scoreRepo.FindTopScoresByGame(ctx, gameID, scoreType, limit)
	if err != nil {
		return nil, NewMinigameError(ErrorCodeRepositoryError, "Failed to get leaderboard", err)
	}

	// 更新排名
	for i, score := range scores {
		rank := int32(i + 1)
		percentile := float64(rank) / float64(len(scores)) * 100
		score.SetRank(rank, percentile)
	}

	return scores, nil
}

// GetPlayerGameHistory 获取玩家游戏历史
func (s *MinigameService) GetPlayerGameHistory(ctx context.Context, playerID uint64, gameType *GameType, limit int32, offset int32) ([]*GameSession, error) {
	// 构建查询条件
	query := &GameSessionQuery{
		PlayerID:  &playerID,
		GameType:  gameType,
		Limit:     &limit,
		Offset:    &offset,
		OrderBy:   "created_at",
		OrderDesc: true,
	}

	// 获取会话历史
	sessions, err := s.sessionRepo.FindByQuery(ctx, query)
	if err != nil {
		return nil, NewMinigameError(ErrorCodeRepositoryError, "Failed to get game history", err)
	}

	return sessions, nil
}

// GetPlayerStatistics 获取玩家统计
func (s *MinigameService) GetPlayerStatistics(ctx context.Context, playerID uint64, gameType *GameType) (*PlayerStatistics, error) {
	// 获取玩家会话统计
	stats, err := s.sessionRepo.GetPlayerStatistics(ctx, playerID, gameType)
	if err != nil {
		return nil, NewMinigameError(ErrorCodeRepositoryError, "Failed to get player statistics", err)
	}

	return stats, nil
}

// 私有方法

// validateGameConfig 验证游戏配置
func (s *MinigameService) validateGameConfig(config *GameConfig) error {
	if config.MaxPlayers <= 0 {
		return fmt.Errorf("max_players must be positive")
	}

	if config.MinPlayers <= 0 {
		return fmt.Errorf("min_players must be positive")
	}

	if config.MinPlayers > config.MaxPlayers {
		return fmt.Errorf("min_players cannot be greater than max_players")
	}

	if config.MaxDuration <= 0 {
		return fmt.Errorf("max_duration must be positive")
	}

	if config.MinDuration <= 0 {
		return fmt.Errorf("min_duration must be positive")
	}

	if config.MinDuration > config.MaxDuration {
		return fmt.Errorf("min_duration cannot be greater than max_duration")
	}

	if !config.Difficulty.IsValid() {
		return fmt.Errorf("invalid difficulty: %v", config.Difficulty)
	}

	return nil
}

// canOperateGame 检查是否可以操作游戏
func (s *MinigameService) canOperateGame(minigame *MinigameAggregate, operatorID uint64) bool {
	// 创建者可以操作
	if minigame.CreatorID == operatorID {
		return true
	}

	// TODO: 添加其他权限检查逻辑，如管理员权限等

	return false
}

// handleGameEnd 处理游戏结束
func (s *MinigameService) handleGameEnd(ctx context.Context, minigame *MinigameAggregate) error {
	// 结算所有玩家会话
	for playerID := range minigame.Players {
		session, err := s.sessionRepo.FindByGameAndPlayer(ctx, minigame.ID, playerID)
		if err != nil {
			continue // 忽略错误，继续处理其他玩家
		}

		// 更新会话状态
		if session.IsActive() {
			session.UpdateStatus(PlayerStatusFinished)
			s.sessionRepo.Save(ctx, session)
		}

		// 发放完成奖励
		go s.grantCompletionRewards(context.Background(), minigame.ID, playerID, session)
	}

	return nil
}

// grantCompletionRewards 发放完成奖励
func (s *MinigameService) grantCompletionRewards(ctx context.Context, gameID string, playerID uint64, session *GameSession) {
	// 基础完成奖励
	baseReward := int64(100)

	// 根据分数计算奖励倍数
	multiplier := 1.0
	if session.Score > 1000 {
		multiplier = 1.5
	}
	if session.Score > 5000 {
		multiplier = 2.0
	}

	finalReward := int64(float64(baseReward) * multiplier)

	// 发放金币奖励
	s.GrantReward(ctx, gameID, playerID, RewardTypeCoins, "coins", finalReward, "game_completion")

	// 发放经验奖励
	expReward := finalReward / 2
	s.GrantReward(ctx, gameID, playerID, RewardTypeExperience, "exp", expReward, "game_completion")
}

// checkAchievements 检查成就
func (s *MinigameService) checkAchievements(ctx context.Context, gameID string, playerID uint64, session *GameSession) {
	// 获取玩家成就
	achievements, err := s.achievementRepo.FindByGameAndPlayer(ctx, gameID, playerID)
	if err != nil {
		return
	}

	// 检查各种成就条件
	for _, achievement := range achievements {
		if achievement.IsCompleted() {
			continue
		}

		// 根据成就类型检查条件
		switch achievement.Category {
		case "score":
			s.checkScoreAchievement(ctx, achievement, session)
		case "time":
			s.checkTimeAchievement(ctx, achievement, session)
		case "moves":
			s.checkMovesAchievement(ctx, achievement, session)
		case "level":
			s.checkLevelAchievement(ctx, achievement, session)
		}

		// 保存成就更新
		if achievement.IsCompleted() {
			s.achievementRepo.Save(ctx, achievement)

			// 发布成就完成事件
			event := NewAchievementCompletedEvent(gameID, playerID, achievement.AchievementID, achievement.Points)
			s.eventBus.Publish(ctx, event)

			// 发放成就奖励
			go s.grantAchievementRewards(context.Background(), gameID, playerID, achievement)
		}
	}
}

// checkScoreAchievement 检查分数成就
func (s *MinigameService) checkScoreAchievement(ctx context.Context, achievement *GameAchievement, session *GameSession) {
	targetScore, ok := achievement.GetCondition("target_score")
	if !ok {
		return
	}

	target, ok := targetScore.(float64)
	if !ok {
		return
	}

	if float64(session.Score) >= target {
		achievement.Complete()
	} else {
		progress := float64(session.Score) / target * 100
		achievement.UpdateProgress(progress)
	}
}

// checkTimeAchievement 检查时间成就
func (s *MinigameService) checkTimeAchievement(ctx context.Context, achievement *GameAchievement, session *GameSession) {
	targetTime, ok := achievement.GetCondition("target_time")
	if !ok {
		return
	}

	target, ok := targetTime.(float64)
	if !ok {
		return
	}

	playTimeSeconds := session.PlayTime.Seconds()
	if playTimeSeconds >= target {
		achievement.Complete()
	} else {
		progress := playTimeSeconds / target * 100
		achievement.UpdateProgress(progress)
	}
}

// checkMovesAchievement 检查移动成就
func (s *MinigameService) checkMovesAchievement(ctx context.Context, achievement *GameAchievement, session *GameSession) {
	targetMoves, ok := achievement.GetCondition("target_moves")
	if !ok {
		return
	}

	target, ok := targetMoves.(float64)
	if !ok {
		return
	}

	if float64(session.Moves) >= target {
		achievement.Complete()
	} else {
		progress := float64(session.Moves) / target * 100
		achievement.UpdateProgress(progress)
	}
}

// checkLevelAchievement 检查等级成就
func (s *MinigameService) checkLevelAchievement(ctx context.Context, achievement *GameAchievement, session *GameSession) {
	targetLevel, ok := achievement.GetCondition("target_level")
	if !ok {
		return
	}

	target, ok := targetLevel.(float64)
	if !ok {
		return
	}

	if float64(session.Level) >= target {
		achievement.Complete()
	} else {
		progress := float64(session.Level) / target * 100
		achievement.UpdateProgress(progress)
	}
}

// grantAchievementRewards 发放成就奖励
func (s *MinigameService) grantAchievementRewards(ctx context.Context, gameID string, playerID uint64, achievement *GameAchievement) {
	// 发放成就积分奖励
	if achievement.Points > 0 {
		s.GrantReward(ctx, gameID, playerID, RewardTypeCoins, "coins", achievement.Points, "achievement")
	}

	// 发放其他奖励
	for _, rewardID := range achievement.Rewards {
		// 根据奖励ID发放对应奖励
		// 这里可以根据具体的奖励系统实现
		s.GrantReward(ctx, gameID, playerID, RewardTypeItems, rewardID, 1, "achievement")
	}
}

// 常量定义

const (
	// 服务相关常量
	// DefaultGameDuration = 30 * time.Minute // Moved to aggregate.go
	MaxConcurrentGames = 1000 // 最大并发游戏数
	MaxPlayersPerGame  = 100  // 每个游戏最大玩家数

	// 奖励相关常量
	BaseCompletionReward = int64(100) // 基础完成奖励
	MaxRewardMultiplier  = 5.0        // 最大奖励倍数

	// 成就相关常量
	MaxAchievementsPerGame = 50 // 每个游戏最大成就数
)

// 辅助结构体

// GameSessionQuery 游戏会话查询条件
type GameSessionQuery struct {
	GameID        *string       `json:"game_id,omitempty"`
	PlayerID      *uint64       `json:"player_id,omitempty"`
	GameType      *GameType     `json:"game_type,omitempty"`
	Status        *PlayerStatus `json:"status,omitempty"`
	StartedAfter  *time.Time    `json:"started_after,omitempty"`
	StartedBefore *time.Time    `json:"started_before,omitempty"`
	EndedAfter    *time.Time    `json:"ended_after,omitempty"`
	EndedBefore   *time.Time    `json:"ended_before,omitempty"`
	MinScore      *int64        `json:"min_score,omitempty"`
	MaxScore      *int64        `json:"max_score,omitempty"`
	MinLevel      *int32        `json:"min_level,omitempty"`
	MaxLevel      *int32        `json:"max_level,omitempty"`
	OrderBy       string        `json:"order_by,omitempty"`
	OrderDesc     bool          `json:"order_desc,omitempty"`
	Limit         *int32        `json:"limit,omitempty"`
	Offset        *int32        `json:"offset,omitempty"`
}

// PlayerStatistics 玩家统计
type PlayerStatistics struct {
	PlayerID          uint64        `json:"player_id"`
	TotalGames        int64         `json:"total_games"`
	CompletedGames    int64         `json:"completed_games"`
	WonGames          int64         `json:"won_games"`
	TotalScore        int64         `json:"total_score"`
	HighestScore      int64         `json:"highest_score"`
	AverageScore      float64       `json:"average_score"`
	TotalPlayTime     time.Duration `json:"total_play_time"`
	AveragePlayTime   time.Duration `json:"average_play_time"`
	TotalMoves        int64         `json:"total_moves"`
	AverageMoves      float64       `json:"average_moves"`
	HighestLevel      int32         `json:"highest_level"`
	TotalAchievements int64         `json:"total_achievements"`
	TotalRewards      int64         `json:"total_rewards"`
	WinRate           float64       `json:"win_rate"`
	CompletionRate    float64       `json:"completion_rate"`
	LastPlayedAt      *time.Time    `json:"last_played_at,omitempty"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

// NewPlayerStatistics 创建玩家统计
func NewPlayerStatistics(playerID uint64) *PlayerStatistics {
	now := time.Now()
	return &PlayerStatistics{
		PlayerID:          playerID,
		TotalGames:        0,
		CompletedGames:    0,
		WonGames:          0,
		TotalScore:        0,
		HighestScore:      0,
		AverageScore:      0.0,
		TotalPlayTime:     0,
		AveragePlayTime:   0,
		TotalMoves:        0,
		AverageMoves:      0.0,
		HighestLevel:      0,
		TotalAchievements: 0,
		TotalRewards:      0,
		WinRate:           0.0,
		CompletionRate:    0.0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// UpdateStatistics 更新统计数据
func (ps *PlayerStatistics) UpdateStatistics(session *GameSession) {
	ps.TotalGames++

	if session.IsFinished() {
		ps.CompletedGames++
	}

	ps.TotalScore += session.Score
	if session.Score > ps.HighestScore {
		ps.HighestScore = session.Score
	}

	ps.TotalPlayTime += session.PlayTime
	ps.TotalMoves += int64(session.Moves)

	if session.Level > ps.HighestLevel {
		ps.HighestLevel = session.Level
	}

	// 重新计算平均值
	if ps.TotalGames > 0 {
		ps.AverageScore = float64(ps.TotalScore) / float64(ps.TotalGames)
		ps.AveragePlayTime = ps.TotalPlayTime / time.Duration(ps.TotalGames)
		ps.AverageMoves = float64(ps.TotalMoves) / float64(ps.TotalGames)
		ps.CompletionRate = float64(ps.CompletedGames) / float64(ps.TotalGames) * 100
	}

	ps.UpdatedAt = time.Now()
}
