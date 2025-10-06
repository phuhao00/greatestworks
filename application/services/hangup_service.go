package services

import (
	"context"
	"fmt"
	"time"

	"greatestworks/internal/domain/player/hangup"
)

// HangupService 挂机应用服务
type HangupService struct {
	hangupRepo hangup.HangupRepository
	// TODO: 实现这些仓储接口
	// locationRepo   hangup.LocationRepository
	// rewardRepo     hangup.RewardRepository
	// statisticsRepo hangup.StatisticsRepository
	// cacheRepo      hangup.CacheRepository
	hangupService *hangup.HangupService
}

// NewHangupService 创建挂机应用服务
func NewHangupService(
	hangupRepo hangup.HangupRepository,
	// TODO: 实现这些仓储接口
	// locationRepo hangup.LocationRepository,
	// rewardRepo hangup.RewardRepository,
	// statisticsRepo hangup.StatisticsRepository,
	// cacheRepo hangup.CacheRepository,
	hangupService *hangup.HangupService,
) *HangupService {
	return &HangupService{
		hangupRepo: hangupRepo,
		// TODO: 实现这些仓储接口
		// locationRepo:   locationRepo,
		// rewardRepo:     rewardRepo,
		// statisticsRepo: statisticsRepo,
		// cacheRepo:      cacheRepo,
		hangupService: hangupService,
	}
}

// StartHangup 开始挂机
func (s *HangupService) StartHangup(ctx context.Context, playerID string, locationID string) error {
	// 检查玩家是否已在挂机
	existingHangup, err := s.hangupRepo.FindActiveByPlayer(playerID)
	if err != nil && !hangup.IsNotFoundError(err) {
		return fmt.Errorf("failed to check existing hangup: %w", err)
	}

	if existingHangup != nil {
		return hangup.ErrAlreadyHanging
	}

	// 获取挂机地点信息
	location, err := s.locationRepo.FindByID(locationID)
	if err != nil {
		return fmt.Errorf("failed to get hangup location: %w", err)
	}

	// 创建挂机记录
	hangupRecord, err := s.hangupService.StartHangup(playerID, location)
	if err != nil {
		return fmt.Errorf("failed to start hangup: %w", err)
	}

	// 保存挂机记录
	if err := s.hangupRepo.Save(hangupRecord); err != nil {
		return fmt.Errorf("failed to save hangup record: %w", err)
	}

	// 更新缓存
	if err := s.cacheRepo.SetActiveHangup(playerID, hangupRecord, time.Hour); err != nil {
		// 缓存失败不影响主流程，只记录日志
		// TODO: 添加日志记录
	}

	return nil
}

// StopHangup 停止挂机
func (s *HangupService) StopHangup(ctx context.Context, playerID string) (*hangup.OfflineReward, error) {
	// 获取当前挂机记录
	hangupRecord, err := s.hangupRepo.FindActiveByPlayer(playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hangup record: %w", err)
	}

	if hangupRecord == nil {
		return nil, hangup.ErrNotHanging
	}

	// 计算离线奖励
	reward, err := s.hangupService.CalculateOfflineReward(hangupRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate offline reward: %w", err)
	}

	// 停止挂机
	if err := hangupRecord.Stop(); err != nil {
		return nil, fmt.Errorf("failed to stop hangup: %w", err)
	}

	// 更新挂机记录
	if err := s.hangupRepo.Update(hangupRecord); err != nil {
		return nil, fmt.Errorf("failed to update hangup record: %w", err)
	}

	// 保存奖励记录
	if err := s.rewardRepo.Save(reward); err != nil {
		return nil, fmt.Errorf("failed to save reward record: %w", err)
	}

	// 更新统计数据
	if err := s.updateStatistics(ctx, playerID, hangupRecord, reward); err != nil {
		// 统计更新失败不影响主流程
		// TODO: 添加日志记录
	}

	// 清除缓存
	if err := s.cacheRepo.DeleteActiveHangup(playerID); err != nil {
		// 缓存清除失败不影响主流程
		// TODO: 添加日志记录
	}

	return reward, nil
}

// GetHangupStatus 获取挂机状态
func (s *HangupService) GetHangupStatus(ctx context.Context, playerID string) (*HangupStatusDTO, error) {
	// 先从缓存获取
	cachedHangup, err := s.cacheRepo.GetActiveHangup(playerID)
	if err == nil && cachedHangup != nil {
		return s.buildHangupStatusDTO(cachedHangup), nil
	}

	// 从数据库获取
	hangupRecord, err := s.hangupRepo.FindActiveByPlayer(playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hangup record: %w", err)
	}

	if hangupRecord == nil {
		return &HangupStatusDTO{
			PlayerID:  playerID,
			IsHanging: false,
		}, nil
	}

	// 更新缓存
	if err := s.cacheRepo.SetActiveHangup(playerID, hangupRecord, time.Hour); err != nil {
		// 缓存更新失败不影响主流程
		// TODO: 添加日志记录
	}

	return s.buildHangupStatusDTO(hangupRecord), nil
}

// GetAvailableLocations 获取可用挂机地点
func (s *HangupService) GetAvailableLocations(ctx context.Context, playerID string) ([]*HangupLocationDTO, error) {
	// 先从缓存获取
	cachedLocations, err := s.cacheRepo.GetAvailableLocations(playerID)
	if err == nil && len(cachedLocations) > 0 {
		return s.buildLocationDTOs(cachedLocations), nil
	}

	// 从数据库获取
	locations, err := s.locationRepo.FindAvailableForPlayer(playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get available locations: %w", err)
	}

	// 更新缓存
	if err := s.cacheRepo.SetAvailableLocations(playerID, locations, time.Hour*2); err != nil {
		// 缓存更新失败不影响主流程
		// TODO: 添加日志记录
	}

	return s.buildLocationDTOs(locations), nil
}

// GetHangupHistory 获取挂机历史
func (s *HangupService) GetHangupHistory(ctx context.Context, playerID string, limit int) ([]*HangupHistoryDTO, error) {
	history, err := s.hangupRepo.FindHistoryByPlayer(playerID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get hangup history: %w", err)
	}

	return s.buildHistoryDTOs(history), nil
}

// GetHangupStatistics 获取挂机统计
func (s *HangupService) GetHangupStatistics(ctx context.Context, playerID string) (*HangupStatisticsDTO, error) {
	stats, err := s.statisticsRepo.FindByPlayer(playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hangup statistics: %w", err)
	}

	return s.buildStatisticsDTO(stats), nil
}

// ClaimOfflineReward 领取离线奖励
func (s *HangupService) ClaimOfflineReward(ctx context.Context, playerID string, rewardID string) error {
	// 获取奖励记录
	reward, err := s.rewardRepo.FindByID(rewardID)
	if err != nil {
		return fmt.Errorf("failed to get reward record: %w", err)
	}

	if reward.PlayerID != playerID {
		return hangup.ErrUnauthorized
	}

	if reward.IsClaimed() {
		return hangup.ErrRewardAlreadyClaimed
	}

	// 领取奖励
	if err := reward.Claim(); err != nil {
		return fmt.Errorf("failed to claim reward: %w", err)
	}

	// 更新奖励记录
	if err := s.rewardRepo.Update(reward); err != nil {
		return fmt.Errorf("failed to update reward record: %w", err)
	}

	return nil
}

// 私有方法

// updateStatistics 更新统计数据
func (s *HangupService) updateStatistics(ctx context.Context, playerID string, hangupRecord *hangup.HangupAggregate, reward *hangup.OfflineReward) error {
	stats, err := s.statisticsRepo.FindByPlayer(playerID)
	if err != nil && !hangup.IsNotFoundError(err) {
		return err
	}

	if stats == nil {
		stats = hangup.NewHangupStatistics(playerID)
	}

	// 更新统计数据
	stats.AddHangupSession(hangupRecord.GetDuration(), reward.GetTotalValue())
	stats.AddLocationTime(hangupRecord.GetLocationID(), hangupRecord.GetDuration())

	// 保存统计数据
	return s.statisticsRepo.Save(stats)
}

// buildHangupStatusDTO 构建挂机状态DTO
func (s *HangupService) buildHangupStatusDTO(hangupRecord *hangup.HangupAggregate) *HangupStatusDTO {
	return &HangupStatusDTO{
		PlayerID:        hangupRecord.GetPlayerID(),
		IsHanging:       hangupRecord.IsActive(),
		LocationID:      hangupRecord.GetLocationID(),
		LocationName:    hangupRecord.GetLocationName(),
		StartTime:       hangupRecord.GetStartTime(),
		Duration:        hangupRecord.GetDuration(),
		Efficiency:      hangupRecord.GetEfficiency(),
		EstimatedReward: s.calculateEstimatedReward(hangupRecord),
	}
}

// buildLocationDTOs 构建地点DTO列表
func (s *HangupService) buildLocationDTOs(locations []*hangup.HangupLocation) []*HangupLocationDTO {
	dtos := make([]*HangupLocationDTO, len(locations))
	for i, location := range locations {
		dtos[i] = &HangupLocationDTO{
			ID:            location.GetID(),
			Name:          location.GetName(),
			Description:   location.GetDescription(),
			BaseRate:      location.GetBaseRate(),
			RequiredLevel: location.GetRequiredLevel(),
			IsUnlocked:    location.IsUnlocked(),
			RewardTypes:   location.GetRewardTypes(),
		}
	}
	return dtos
}

// buildHistoryDTOs 构建历史DTO列表
func (s *HangupService) buildHistoryDTOs(history []*hangup.HangupAggregate) []*HangupHistoryDTO {
	dtos := make([]*HangupHistoryDTO, len(history))
	for i, record := range history {
		dtos[i] = &HangupHistoryDTO{
			ID:           record.GetID(),
			LocationID:   record.GetLocationID(),
			LocationName: record.GetLocationName(),
			StartTime:    record.GetStartTime(),
			EndTime:      record.GetEndTime(),
			Duration:     record.GetDuration(),
			Efficiency:   record.GetEfficiency(),
			TotalReward:  record.GetTotalReward(),
		}
	}
	return dtos
}

// buildStatisticsDTO 构建统计DTO
func (s *HangupService) buildStatisticsDTO(stats *hangup.HangupStatistics) *HangupStatisticsDTO {
	return &HangupStatisticsDTO{
		PlayerID:         stats.GetPlayerID(),
		TotalSessions:    stats.GetTotalSessions(),
		TotalDuration:    stats.GetTotalDuration(),
		TotalReward:      stats.GetTotalReward(),
		AverageDuration:  stats.GetAverageDuration(),
		AverageReward:    stats.GetAverageReward(),
		FavoriteLocation: stats.GetFavoriteLocation(),
		LocationStats:    stats.GetLocationStats(),
		LastHangupTime:   stats.GetLastHangupTime(),
	}
}

// calculateEstimatedReward 计算预估奖励
func (s *HangupService) calculateEstimatedReward(hangupRecord *hangup.HangupAggregate) int64 {
	// 基于当前挂机时长和效率计算预估奖励
	duration := hangupRecord.GetDuration()
	efficiency := hangupRecord.GetEfficiency()
	baseRate := hangupRecord.GetBaseRate()

	return int64(duration.Hours() * efficiency * baseRate)
}

// DTO 定义

// HangupStatusDTO 挂机状态DTO
type HangupStatusDTO struct {
	PlayerID        string        `json:"player_id"`
	IsHanging       bool          `json:"is_hanging"`
	LocationID      string        `json:"location_id,omitempty"`
	LocationName    string        `json:"location_name,omitempty"`
	StartTime       time.Time     `json:"start_time,omitempty"`
	Duration        time.Duration `json:"duration,omitempty"`
	Efficiency      float64       `json:"efficiency,omitempty"`
	EstimatedReward int64         `json:"estimated_reward,omitempty"`
}

// HangupLocationDTO 挂机地点DTO
type HangupLocationDTO struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	BaseRate      float64  `json:"base_rate"`
	RequiredLevel int      `json:"required_level"`
	IsUnlocked    bool     `json:"is_unlocked"`
	RewardTypes   []string `json:"reward_types"`
}

// HangupHistoryDTO 挂机历史DTO
type HangupHistoryDTO struct {
	ID           string        `json:"id"`
	LocationID   string        `json:"location_id"`
	LocationName string        `json:"location_name"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
	Duration     time.Duration `json:"duration"`
	Efficiency   float64       `json:"efficiency"`
	TotalReward  int64         `json:"total_reward"`
}

// HangupStatisticsDTO 挂机统计DTO
type HangupStatisticsDTO struct {
	PlayerID         string                    `json:"player_id"`
	TotalSessions    int64                     `json:"total_sessions"`
	TotalDuration    time.Duration             `json:"total_duration"`
	TotalReward      int64                     `json:"total_reward"`
	AverageDuration  time.Duration             `json:"average_duration"`
	AverageReward    float64                   `json:"average_reward"`
	FavoriteLocation string                    `json:"favorite_location"`
	LocationStats    map[string]*LocationStats `json:"location_stats"`
	LastHangupTime   time.Time                 `json:"last_hangup_time"`
}

// LocationStats 地点统计
type LocationStats struct {
	LocationID    string        `json:"location_id"`
	LocationName  string        `json:"location_name"`
	TotalTime     time.Duration `json:"total_time"`
	TotalSessions int64         `json:"total_sessions"`
	TotalReward   int64         `json:"total_reward"`
	AverageReward float64       `json:"average_reward"`
}
