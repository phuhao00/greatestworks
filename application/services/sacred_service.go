package services

import (
	"context"
	"fmt"
	"time"

	"greatestworks/internal/domain/scene/sacred"
)

// SacredService 圣地应用服务
type SacredService struct {
	sacredRepo    sacred.SacredPlaceRepository
	challengeRepo sacred.ChallengeRepository
	blessingRepo  sacred.BlessingRepository
	// TODO: 实现这些仓储接口
	// artifactRepo   sacred.ArtifactRepository
	// statisticsRepo sacred.StatisticsRepository
	// cacheRepo      sacred.CacheRepository
	sacredService *sacred.SacredService
}

// NewSacredService 创建圣地应用服务
func NewSacredService(
	sacredRepo sacred.SacredPlaceRepository,
	challengeRepo sacred.ChallengeRepository,
	blessingRepo sacred.BlessingRepository,
	// TODO: 实现这些仓储接口
	// artifactRepo sacred.ArtifactRepository,
	// statisticsRepo sacred.StatisticsRepository,
	// cacheRepo sacred.CacheRepository,
	sacredService *sacred.SacredService,
) *SacredService {
	return &SacredService{
		sacredRepo:    sacredRepo,
		challengeRepo: challengeRepo,
		blessingRepo:  blessingRepo,
		// artifactRepo:   artifactRepo,
		statisticsRepo: statisticsRepo,
		cacheRepo:      cacheRepo,
		sacredService:  sacredService,
	}
}

// GetSacredPlace 获取圣地信息
func (s *SacredService) GetSacredPlace(ctx context.Context, sacredID string) (*SacredPlaceDTO, error) {
	// 先从缓存获取
	cachedSacred, err := s.cacheRepo.GetSacredPlace(sacredID)
	if err == nil && cachedSacred != nil {
		return s.buildSacredPlaceDTO(cachedSacred), nil
	}

	// 从数据库获取
	sacredPlace, err := s.sacredRepo.FindByID(sacredID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sacred place: %w", err)
	}

	// 更新缓存
	if err := s.cacheRepo.SetSacredPlace(sacredID, sacredPlace, time.Hour); err != nil {
		// 缓存更新失败不影响主流程
		// TODO: 添加日志记录
	}

	return s.buildSacredPlaceDTO(sacredPlace), nil
}

// GetAvailableSacredPlaces 获取可用圣地列表
func (s *SacredService) GetAvailableSacredPlaces(ctx context.Context, playerID string) ([]*SacredPlaceDTO, error) {
	// 先从缓存获取
	cachedPlaces, err := s.cacheRepo.GetAvailableSacredPlaces(playerID)
	if err == nil && len(cachedPlaces) > 0 {
		return s.buildSacredPlaceDTOs(cachedPlaces), nil
	}

	// 从数据库获取
	sacredPlaces, err := s.sacredRepo.FindAvailableForPlayer(playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get available sacred places: %w", err)
	}

	// 更新缓存
	if err := s.cacheRepo.SetAvailableSacredPlaces(playerID, sacredPlaces, time.Hour*2); err != nil {
		// 缓存更新失败不影响主流程
		// TODO: 添加日志记录
	}

	return s.buildSacredPlaceDTOs(sacredPlaces), nil
}

// EnterSacredPlace 进入圣地
func (s *SacredService) EnterSacredPlace(ctx context.Context, playerID string, sacredID string) error {
	// 获取圣地信息
	sacredPlace, err := s.sacredRepo.FindByID(sacredID)
	if err != nil {
		return fmt.Errorf("failed to get sacred place: %w", err)
	}

	// 检查进入条件
	if err := s.sacredService.CheckEnterConditions(playerID, sacredPlace); err != nil {
		return fmt.Errorf("failed to check enter conditions: %w", err)
	}

	// 记录进入
	if err := sacredPlace.AddVisitor(playerID); err != nil {
		return fmt.Errorf("failed to add visitor: %w", err)
	}

	// 更新圣地
	if err := s.sacredRepo.Update(sacredPlace); err != nil {
		return fmt.Errorf("failed to update sacred place: %w", err)
	}

	// 清除相关缓存
	if err := s.cacheRepo.DeleteSacredPlace(sacredID); err != nil {
		// 缓存清除失败不影响主流程
		// TODO: 添加日志记录
	}

	return nil
}

// StartChallenge 开始挑战
func (s *SacredService) StartChallenge(ctx context.Context, playerID string, challengeID string) (*ChallengeInstanceDTO, error) {
	// 获取挑战信息
	challenge, err := s.challengeRepo.FindByID(challengeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge: %w", err)
	}

	// 开始挑战
	challengeInstance, err := s.sacredService.StartChallenge(playerID, challenge)
	if err != nil {
		return nil, fmt.Errorf("failed to start challenge: %w", err)
	}

	// 保存挑战实例
	if err := s.challengeRepo.SaveInstance(challengeInstance); err != nil {
		return nil, fmt.Errorf("failed to save challenge instance: %w", err)
	}

	return s.buildChallengeInstanceDTO(challengeInstance), nil
}

// CompleteChallenge 完成挑战
func (s *SacredService) CompleteChallenge(ctx context.Context, playerID string, instanceID string, result *ChallengeResultData) (*ChallengeRewardDTO, error) {
	// 获取挑战实例
	instance, err := s.challengeRepo.FindInstanceByID(instanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge instance: %w", err)
	}

	if instance.GetPlayerID() != playerID {
		return nil, sacred.ErrUnauthorized
	}

	// 完成挑战
	reward, err := s.sacredService.CompleteChallenge(instance, result)
	if err != nil {
		return nil, fmt.Errorf("failed to complete challenge: %w", err)
	}

	// 更新挑战实例
	if err := s.challengeRepo.UpdateInstance(instance); err != nil {
		return nil, fmt.Errorf("failed to update challenge instance: %w", err)
	}

	// 更新统计数据
	if err := s.updateChallengeStatistics(ctx, playerID, instance, reward); err != nil {
		// 统计更新失败不影响主流程
		// TODO: 添加日志记录
	}

	return s.buildChallengeRewardDTO(reward), nil
}

// ActivateBlessing 激活祝福
func (s *SacredService) ActivateBlessing(ctx context.Context, playerID string, blessingID string) error {
	// 获取祝福信息
	blessing, err := s.blessingRepo.FindByID(blessingID)
	if err != nil {
		return fmt.Errorf("failed to get blessing: %w", err)
	}

	// 激活祝福
	if err := s.sacredService.ActivateBlessing(playerID, blessing); err != nil {
		return fmt.Errorf("failed to activate blessing: %w", err)
	}

	// 保存祝福状态
	if err := s.blessingRepo.SavePlayerBlessing(playerID, blessing); err != nil {
		return fmt.Errorf("failed to save player blessing: %w", err)
	}

	return nil
}

// GetPlayerBlessings 获取玩家祝福
func (s *SacredService) GetPlayerBlessings(ctx context.Context, playerID string) ([]*PlayerBlessingDTO, error) {
	blessings, err := s.blessingRepo.FindByPlayer(playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player blessings: %w", err)
	}

	return s.buildPlayerBlessingDTOs(blessings), nil
}

// GetAvailableArtifacts 获取可用圣物
func (s *SacredService) GetAvailableArtifacts(ctx context.Context, sacredID string) ([]*ArtifactDTO, error) {
	artifacts, err := s.artifactRepo.FindBySacredPlace(sacredID)
	if err != nil {
		return nil, fmt.Errorf("failed to get artifacts: %w", err)
	}

	return s.buildArtifactDTOs(artifacts), nil
}

// UseArtifact 使用圣物
func (s *SacredService) UseArtifact(ctx context.Context, playerID string, artifactID string) (*ArtifactEffectDTO, error) {
	// 获取圣物信息
	artifact, err := s.artifactRepo.FindByID(artifactID)
	if err != nil {
		return nil, fmt.Errorf("failed to get artifact: %w", err)
	}

	// 使用圣物
	effect, err := s.sacredService.UseArtifact(playerID, artifact)
	if err != nil {
		return nil, fmt.Errorf("failed to use artifact: %w", err)
	}

	// 更新圣物状态
	if err := s.artifactRepo.Update(artifact); err != nil {
		return nil, fmt.Errorf("failed to update artifact: %w", err)
	}

	return s.buildArtifactEffectDTO(effect), nil
}

// GetSacredStatistics 获取圣地统计
func (s *SacredService) GetSacredStatistics(ctx context.Context, playerID string) (*SacredStatisticsDTO, error) {
	stats, err := s.statisticsRepo.FindByPlayer(playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sacred statistics: %w", err)
	}

	return s.buildStatisticsDTO(stats), nil
}

// GetChallengeHistory 获取挑战历史
func (s *SacredService) GetChallengeHistory(ctx context.Context, playerID string, limit int) ([]*ChallengeHistoryDTO, error) {
	history, err := s.challengeRepo.FindHistoryByPlayer(playerID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge history: %w", err)
	}

	return s.buildChallengeHistoryDTOs(history), nil
}

// 私有方法

// updateChallengeStatistics 更新挑战统计
func (s *SacredService) updateChallengeStatistics(ctx context.Context, playerID string, instance *sacred.Challenge, reward *sacred.ChallengeReward) error {
	stats, err := s.statisticsRepo.FindByPlayer(playerID)
	if err != nil && !sacred.IsNotFoundError(err) {
		return err
	}

	if stats == nil {
		stats = sacred.NewSacredStatistics(playerID)
	}

	// 更新统计数据
	stats.AddChallengeResult(instance.GetChallengeID(), instance.IsCompleted(), reward.GetTotalValue())
	stats.UpdateLastChallengeTime(instance.GetCompletedAt())

	// 保存统计数据
	return s.statisticsRepo.Save(stats)
}

// buildSacredPlaceDTO 构建圣地DTO
func (s *SacredService) buildSacredPlaceDTO(sacredPlace *sacred.SacredPlaceAggregate) *SacredPlaceDTO {
	return &SacredPlaceDTO{
		ID:            sacredPlace.GetID(),
		Name:          sacredPlace.GetName(),
		Description:   sacredPlace.GetDescription(),
		Level:         sacredPlace.GetLevel(),
		Status:        string(sacredPlace.GetStatus()),
		RequiredLevel: sacredPlace.GetRequiredLevel(),
		VisitorCount:  sacredPlace.GetVisitorCount(),
		MaxVisitors:   sacredPlace.GetMaxVisitors(),
		Challenges:    sacredPlace.GetChallengeIDs(),
		Blessings:     sacredPlace.GetBlessingIDs(),
		Artifacts:     sacredPlace.GetArtifactIDs(),
		CreatedAt:     sacredPlace.GetCreatedAt(),
		UpdatedAt:     sacredPlace.GetUpdatedAt(),
	}
}

// buildSacredPlaceDTOs 构建圣地DTO列表
func (s *SacredService) buildSacredPlaceDTOs(sacredPlaces []*sacred.SacredPlaceAggregate) []*SacredPlaceDTO {
	dtos := make([]*SacredPlaceDTO, len(sacredPlaces))
	for i, place := range sacredPlaces {
		dtos[i] = s.buildSacredPlaceDTO(place)
	}
	return dtos
}

// buildChallengeInstanceDTO 构建挑战实例DTO
func (s *SacredService) buildChallengeInstanceDTO(instance *sacred.Challenge) *ChallengeInstanceDTO {
	return &ChallengeInstanceDTO{
		ID:          instance.GetID(),
		PlayerID:    instance.GetPlayerID(),
		ChallengeID: instance.GetChallengeID(),
		Status:      string(instance.GetStatus()),
		StartTime:   instance.GetStartTime(),
		EndTime:     instance.GetEndTime(),
		Duration:    instance.GetDuration(),
		Difficulty:  instance.GetDifficulty(),
		Progress:    instance.GetProgress(),
		IsCompleted: instance.IsCompleted(),
	}
}

// buildChallengeRewardDTO 构建挑战奖励DTO
func (s *SacredService) buildChallengeRewardDTO(reward *sacred.ChallengeReward) *ChallengeRewardDTO {
	return &ChallengeRewardDTO{
		Experience: reward.GetExperience(),
		Items:      reward.GetItems(),
		Blessings:  reward.GetBlessings(),
		TotalValue: reward.GetTotalValue(),
	}
}

// buildPlayerBlessingDTOs 构建玩家祝福DTO列表
func (s *SacredService) buildPlayerBlessingDTOs(blessings []*sacred.PlayerBlessing) []*PlayerBlessingDTO {
	dtos := make([]*PlayerBlessingDTO, len(blessings))
	for i, blessing := range blessings {
		dtos[i] = &PlayerBlessingDTO{
			BlessingID:    blessing.GetBlessingID(),
			Name:          blessing.GetName(),
			Description:   blessing.GetDescription(),
			EffectType:    string(blessing.GetEffectType()),
			EffectValue:   blessing.GetEffectValue(),
			Duration:      blessing.GetDuration(),
			RemainingTime: blessing.GetRemainingTime(),
			IsActive:      blessing.IsActive(),
			ActivatedAt:   blessing.GetActivatedAt(),
		}
	}
	return dtos
}

// buildArtifactDTOs 构建圣物DTO列表
func (s *SacredService) buildArtifactDTOs(artifacts []*sacred.Artifact) []*ArtifactDTO {
	dtos := make([]*ArtifactDTO, len(artifacts))
	for i, artifact := range artifacts {
		dtos[i] = &ArtifactDTO{
			ID:          artifact.GetID(),
			Name:        artifact.GetName(),
			Description: artifact.GetDescription(),
			Type:        string(artifact.GetType()),
			Rarity:      string(artifact.GetRarity()),
			Power:       artifact.GetPower(),
			Cooldown:    artifact.GetCooldown(),
			UsageCount:  artifact.GetUsageCount(),
			MaxUsage:    artifact.GetMaxUsage(),
			IsAvailable: artifact.IsAvailable(),
		}
	}
	return dtos
}

// buildArtifactEffectDTO 构建圣物效果DTO
func (s *SacredService) buildArtifactEffectDTO(effect *sacred.ArtifactEffect) *ArtifactEffectDTO {
	return &ArtifactEffectDTO{
		EffectType:  string(effect.GetEffectType()),
		Value:       effect.GetValue(),
		Duration:    effect.GetDuration(),
		Description: effect.GetDescription(),
	}
}

// buildChallengeHistoryDTOs 构建挑战历史DTO列表
func (s *SacredService) buildChallengeHistoryDTOs(history []*sacred.Challenge) []*ChallengeHistoryDTO {
	dtos := make([]*ChallengeHistoryDTO, len(history))
	for i, instance := range history {
		dtos[i] = &ChallengeHistoryDTO{
			ID:          instance.GetID(),
			ChallengeID: instance.GetChallengeID(),
			StartTime:   instance.GetStartTime(),
			EndTime:     instance.GetEndTime(),
			Duration:    instance.GetDuration(),
			Difficulty:  instance.GetDifficulty(),
			IsCompleted: instance.IsCompleted(),
			Reward:      instance.GetRewardValue(),
		}
	}
	return dtos
}

// buildStatisticsDTO 构建统计DTO
func (s *SacredService) buildStatisticsDTO(stats *sacred.SacredStatistics) *SacredStatisticsDTO {
	return &SacredStatisticsDTO{
		PlayerID:            stats.GetPlayerID(),
		TotalChallenges:     stats.GetTotalChallenges(),
		CompletedChallenges: stats.GetCompletedChallenges(),
		FailedChallenges:    stats.GetFailedChallenges(),
		CompletionRate:      stats.GetCompletionRate(),
		TotalReward:         stats.GetTotalReward(),
		AverageReward:       stats.GetAverageReward(),
		FavoriteChallenge:   stats.GetFavoriteChallenge(),
		ActiveBlessings:     stats.GetActiveBlessings(),
		LastChallengeTime:   stats.GetLastChallengeTime(),
	}
}

// DTO 定义

// SacredPlaceDTO 圣地DTO
type SacredPlaceDTO struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Level         int       `json:"level"`
	Status        string    `json:"status"`
	RequiredLevel int       `json:"required_level"`
	VisitorCount  int       `json:"visitor_count"`
	MaxVisitors   int       `json:"max_visitors"`
	Challenges    []string  `json:"challenges"`
	Blessings     []string  `json:"blessings"`
	Artifacts     []string  `json:"artifacts"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ChallengeInstanceDTO 挑战实例DTO
type ChallengeInstanceDTO struct {
	ID          string        `json:"id"`
	PlayerID    string        `json:"player_id"`
	ChallengeID string        `json:"challenge_id"`
	Status      string        `json:"status"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	Duration    time.Duration `json:"duration"`
	Difficulty  float64       `json:"difficulty"`
	Progress    float64       `json:"progress"`
	IsCompleted bool          `json:"is_completed"`
}

// ChallengeRewardDTO 挑战奖励DTO
type ChallengeRewardDTO struct {
	Experience int64          `json:"experience"`
	Items      map[string]int `json:"items"`
	Blessings  []string       `json:"blessings"`
	TotalValue int64          `json:"total_value"`
}

// PlayerBlessingDTO 玩家祝福DTO
type PlayerBlessingDTO struct {
	BlessingID    string        `json:"blessing_id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	EffectType    string        `json:"effect_type"`
	EffectValue   float64       `json:"effect_value"`
	Duration      time.Duration `json:"duration"`
	RemainingTime time.Duration `json:"remaining_time"`
	IsActive      bool          `json:"is_active"`
	ActivatedAt   time.Time     `json:"activated_at"`
}

// ArtifactDTO 圣物DTO
type ArtifactDTO struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Type        string        `json:"type"`
	Rarity      string        `json:"rarity"`
	Power       float64       `json:"power"`
	Cooldown    time.Duration `json:"cooldown"`
	UsageCount  int           `json:"usage_count"`
	MaxUsage    int           `json:"max_usage"`
	IsAvailable bool          `json:"is_available"`
}

// ArtifactEffectDTO 圣物效果DTO
type ArtifactEffectDTO struct {
	EffectType  string        `json:"effect_type"`
	Value       float64       `json:"value"`
	Duration    time.Duration `json:"duration"`
	Description string        `json:"description"`
}

// ChallengeHistoryDTO 挑战历史DTO
type ChallengeHistoryDTO struct {
	ID          string        `json:"id"`
	ChallengeID string        `json:"challenge_id"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	Duration    time.Duration `json:"duration"`
	Difficulty  float64       `json:"difficulty"`
	IsCompleted bool          `json:"is_completed"`
	Reward      int64         `json:"reward"`
}

// SacredStatisticsDTO 圣地统计DTO
type SacredStatisticsDTO struct {
	PlayerID            string    `json:"player_id"`
	TotalChallenges     int64     `json:"total_challenges"`
	CompletedChallenges int64     `json:"completed_challenges"`
	FailedChallenges    int64     `json:"failed_challenges"`
	CompletionRate      float64   `json:"completion_rate"`
	TotalReward         int64     `json:"total_reward"`
	AverageReward       float64   `json:"average_reward"`
	FavoriteChallenge   string    `json:"favorite_challenge"`
	ActiveBlessings     int       `json:"active_blessings"`
	LastChallengeTime   time.Time `json:"last_challenge_time"`
}

// ChallengeResultData 挑战结果数据
type ChallengeResultData struct {
	Score       int64              `json:"score"`
	TimeUsed    time.Duration      `json:"time_used"`
	Actions     []string           `json:"actions"`
	Performance map[string]float64 `json:"performance"`
}
