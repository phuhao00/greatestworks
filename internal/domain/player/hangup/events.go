package hangup

import (
	"fmt"
	"time"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	GetEventID() string
	GetEventType() string
	GetAggregateID() string
	GetOccurredAt() time.Time
	GetEventData() interface{}
}

// BaseDomainEvent 基础领域事件
type BaseDomainEvent struct {
	EventID     string      `json:"event_id"`
	EventType   string      `json:"event_type"`
	AggregateID string      `json:"aggregate_id"`
	OccurredAt  time.Time   `json:"occurred_at"`
	EventData   interface{} `json:"event_data"`
}

// GetEventID 获取事件ID
func (e *BaseDomainEvent) GetEventID() string {
	return e.EventID
}

// GetEventType 获取事件类型
func (e *BaseDomainEvent) GetEventType() string {
	return e.EventType
}

// GetAggregateID 获取聚合根ID
func (e *BaseDomainEvent) GetAggregateID() string {
	return e.AggregateID
}

// GetOccurredAt 获取发生时间
func (e *BaseDomainEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetEventData 获取事件数据
func (e *BaseDomainEvent) GetEventData() interface{} {
	return e.EventData
}

// HangupStartedEvent 开始挂机事件
type HangupStartedEvent struct {
	*BaseDomainEvent
	PlayerID   string `json:"player_id"`
	LocationID string `json:"location_id"`
	LocationName string `json:"location_name"`
	StartTime  time.Time `json:"start_time"`
	IsOnline   bool   `json:"is_online"`
}

// NewHangupStartedEvent 创建开始挂机事件
func NewHangupStartedEvent(playerID, locationID, locationName string, isOnline bool) *HangupStartedEvent {
	return &HangupStartedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "HangupStarted",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:     playerID,
		LocationID:   locationID,
		LocationName: locationName,
		StartTime:    time.Now(),
		IsOnline:     isOnline,
	}
}

// HangupStoppedEvent 停止挂机事件
type HangupStoppedEvent struct {
	*BaseDomainEvent
	PlayerID   string        `json:"player_id"`
	LocationID string        `json:"location_id"`
	LocationName string      `json:"location_name"`
	StartTime  time.Time     `json:"start_time"`
	EndTime    time.Time     `json:"end_time"`
	Duration   time.Duration `json:"duration"`
	IsOnline   bool          `json:"is_online"`
}

// NewHangupStoppedEvent 创建停止挂机事件
func NewHangupStoppedEvent(playerID, locationID, locationName string, startTime time.Time, isOnline bool) *HangupStoppedEvent {
	endTime := time.Now()
	return &HangupStoppedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "HangupStopped",
			AggregateID: playerID,
			OccurredAt:  endTime,
		},
		PlayerID:     playerID,
		LocationID:   locationID,
		LocationName: locationName,
		StartTime:    startTime,
		EndTime:      endTime,
		Duration:     endTime.Sub(startTime),
		IsOnline:     isOnline,
	}
}

// HangupLocationChangedEvent 挂机地点变更事件
type HangupLocationChangedEvent struct {
	*BaseDomainEvent
	PlayerID        string `json:"player_id"`
	PreviousLocationID string `json:"previous_location_id,omitempty"`
	PreviousLocationName string `json:"previous_location_name,omitempty"`
	NewLocationID   string `json:"new_location_id"`
	NewLocationName string `json:"new_location_name"`
	Reason          string `json:"reason,omitempty"`
}

// NewHangupLocationChangedEvent 创建挂机地点变更事件
func NewHangupLocationChangedEvent(playerID, prevLocationID, prevLocationName, newLocationID, newLocationName, reason string) *HangupLocationChangedEvent {
	return &HangupLocationChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "HangupLocationChanged",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:             playerID,
		PreviousLocationID:   prevLocationID,
		PreviousLocationName: prevLocationName,
		NewLocationID:        newLocationID,
		NewLocationName:      newLocationName,
		Reason:               reason,
	}
}

// OfflineRewardCalculatedEvent 离线奖励计算事件
type OfflineRewardCalculatedEvent struct {
	*BaseDomainEvent
	PlayerID        string        `json:"player_id"`
	LocationID      string        `json:"location_id"`
	LocationName    string        `json:"location_name"`
	OfflineDuration time.Duration `json:"offline_duration"`
	Experience      int64         `json:"experience"`
	Gold            int64         `json:"gold"`
	Items           []RewardItem  `json:"items"`
	EfficiencyBonus float64       `json:"efficiency_bonus"`
	CalculatedAt    time.Time     `json:"calculated_at"`
}

// NewOfflineRewardCalculatedEvent 创建离线奖励计算事件
func NewOfflineRewardCalculatedEvent(playerID, locationID, locationName string, reward *OfflineReward, efficiencyBonus float64) *OfflineRewardCalculatedEvent {
	return &OfflineRewardCalculatedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "OfflineRewardCalculated",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:        playerID,
		LocationID:      locationID,
		LocationName:    locationName,
		OfflineDuration: reward.OfflineDuration,
		Experience:      reward.Experience,
		Gold:            reward.Gold,
		Items:           reward.Items,
		EfficiencyBonus: efficiencyBonus,
		CalculatedAt:    reward.CalculatedAt,
	}
}

// OfflineRewardClaimedEvent 离线奖励领取事件
type OfflineRewardClaimedEvent struct {
	*BaseDomainEvent
	PlayerID        string        `json:"player_id"`
	LocationID      string        `json:"location_id"`
	LocationName    string        `json:"location_name"`
	OfflineDuration time.Duration `json:"offline_duration"`
	Experience      int64         `json:"experience"`
	Gold            int64         `json:"gold"`
	Items           []RewardItem  `json:"items"`
	ClaimedAt       time.Time     `json:"claimed_at"`
}

// NewOfflineRewardClaimedEvent 创建离线奖励领取事件
func NewOfflineRewardClaimedEvent(playerID, locationID, locationName string, reward *OfflineReward) *OfflineRewardClaimedEvent {
	return &OfflineRewardClaimedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "OfflineRewardClaimed",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:        playerID,
		LocationID:      locationID,
		LocationName:    locationName,
		OfflineDuration: reward.OfflineDuration,
		Experience:      reward.Experience,
		Gold:            reward.Gold,
		Items:           reward.Items,
		ClaimedAt:       reward.ClaimedAt,
	}
}

// HangupEfficiencyUpdatedEvent 挂机效率更新事件
type HangupEfficiencyUpdatedEvent struct {
	*BaseDomainEvent
	PlayerID        string  `json:"player_id"`
	PreviousBonus   float64 `json:"previous_bonus"`
	NewBonus        float64 `json:"new_bonus"`
	VipBonus        float64 `json:"vip_bonus"`
	EquipmentBonus  float64 `json:"equipment_bonus"`
	SkillBonus      float64 `json:"skill_bonus"`
	GuildBonus      float64 `json:"guild_bonus"`
	EventBonus      float64 `json:"event_bonus"`
	UpdateReason    string  `json:"update_reason"`
}

// NewHangupEfficiencyUpdatedEvent 创建挂机效率更新事件
func NewHangupEfficiencyUpdatedEvent(playerID string, previousBonus, newBonus float64, efficiency *EfficiencyBonus, reason string) *HangupEfficiencyUpdatedEvent {
	return &HangupEfficiencyUpdatedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "HangupEfficiencyUpdated",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:       playerID,
		PreviousBonus:  previousBonus,
		NewBonus:       newBonus,
		VipBonus:       efficiency.GetVipBonus(),
		EquipmentBonus: efficiency.GetEquipmentBonus(),
		SkillBonus:     efficiency.GetSkillBonus(),
		GuildBonus:     efficiency.GetGuildBonus(),
		EventBonus:     efficiency.GetEventBonus(),
		UpdateReason:   reason,
	}
}

// HangupLocationUnlockedEvent 挂机地点解锁事件
type HangupLocationUnlockedEvent struct {
	*BaseDomainEvent
	PlayerID     string       `json:"player_id"`
	LocationID   string       `json:"location_id"`
	LocationName string       `json:"location_name"`
	LocationType LocationType `json:"location_type"`
	RequiredLevel int         `json:"required_level"`
	PlayerLevel  int          `json:"player_level"`
	UnlockMethod string       `json:"unlock_method"`
}

// NewHangupLocationUnlockedEvent 创建挂机地点解锁事件
func NewHangupLocationUnlockedEvent(playerID string, location *HangupLocation, playerLevel int, unlockMethod string) *HangupLocationUnlockedEvent {
	return &HangupLocationUnlockedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "HangupLocationUnlocked",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:      playerID,
		LocationID:    location.GetID(),
		LocationName:  location.GetName(),
		LocationType:  location.GetLocationType(),
		RequiredLevel: location.GetRequiredLevel(),
		PlayerLevel:   playerLevel,
		UnlockMethod:  unlockMethod,
	}
}

// HangupMilestoneReachedEvent 挂机里程碑达成事件
type HangupMilestoneReachedEvent struct {
	*BaseDomainEvent
	PlayerID        string        `json:"player_id"`
	MilestoneType   string        `json:"milestone_type"`
	MilestoneName   string        `json:"milestone_name"`
	CurrentValue    int64         `json:"current_value"`
	TargetValue     int64         `json:"target_value"`
	Reward          *BaseReward   `json:"reward,omitempty"`
	TotalHangupTime time.Duration `json:"total_hangup_time"`
}

// NewHangupMilestoneReachedEvent 创建挂机里程碑达成事件
func NewHangupMilestoneReachedEvent(playerID, milestoneType, milestoneName string, currentValue, targetValue int64, reward *BaseReward, totalTime time.Duration) *HangupMilestoneReachedEvent {
	return &HangupMilestoneReachedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "HangupMilestoneReached",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:        playerID,
		MilestoneType:   milestoneType,
		MilestoneName:   milestoneName,
		CurrentValue:    currentValue,
		TargetValue:     targetValue,
		Reward:          reward,
		TotalHangupTime: totalTime,
	}
}

// HangupRankingChangedEvent 挂机排名变化事件
type HangupRankingChangedEvent struct {
	*BaseDomainEvent
	PlayerID     string        `json:"player_id"`
	PlayerName   string        `json:"player_name"`
	RankType     string        `json:"rank_type"`
	PreviousRank int           `json:"previous_rank"`
	NewRank      int           `json:"new_rank"`
	TotalTime    time.Duration `json:"total_time"`
	TotalRewards int64         `json:"total_rewards"`
	Efficiency   float64       `json:"efficiency"`
}

// NewHangupRankingChangedEvent 创建挂机排名变化事件
func NewHangupRankingChangedEvent(playerID, playerName, rankType string, prevRank, newRank int, totalTime time.Duration, totalRewards int64, efficiency float64) *HangupRankingChangedEvent {
	return &HangupRankingChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "HangupRankingChanged",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:     playerID,
		PlayerName:   playerName,
		RankType:     rankType,
		PreviousRank: prevRank,
		NewRank:      newRank,
		TotalTime:    totalTime,
		TotalRewards: totalRewards,
		Efficiency:   efficiency,
	}
}

// HangupSystemInitializedEvent 挂机系统初始化事件
type HangupSystemInitializedEvent struct {
	*BaseDomainEvent
	PlayerID          string `json:"player_id"`
	AvailableLocations int   `json:"available_locations"`
	UnlockedLocations int    `json:"unlocked_locations"`
	InitialEfficiency float64 `json:"initial_efficiency"`
	PlayerLevel       int     `json:"player_level"`
}

// NewHangupSystemInitializedEvent 创建挂机系统初始化事件
func NewHangupSystemInitializedEvent(playerID string, availableLocations, unlockedLocations int, initialEfficiency float64, playerLevel int) *HangupSystemInitializedEvent {
	return &HangupSystemInitializedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "HangupSystemInitialized",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:           playerID,
		AvailableLocations: availableLocations,
		UnlockedLocations:  unlockedLocations,
		InitialEfficiency:  initialEfficiency,
		PlayerLevel:        playerLevel,
	}
}

// HangupAnomalyDetectedEvent 挂机异常检测事件
type HangupAnomalyDetectedEvent struct {
	*BaseDomainEvent
	PlayerID     string                 `json:"player_id"`
	AnomalyType  string                 `json:"anomaly_type"`
	Description  string                 `json:"description"`
	Severity     string                 `json:"severity"`
	Metrics      map[string]interface{} `json:"metrics"`
	Suggestions  []string               `json:"suggestions"`
	AutoResolved bool                   `json:"auto_resolved"`
}

// NewHangupAnomalyDetectedEvent 创建挂机异常检测事件
func NewHangupAnomalyDetectedEvent(playerID, anomalyType, description, severity string, metrics map[string]interface{}, suggestions []string, autoResolved bool) *HangupAnomalyDetectedEvent {
	return &HangupAnomalyDetectedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "HangupAnomalyDetected",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:     playerID,
		AnomalyType:  anomalyType,
		Description:  description,
		Severity:     severity,
		Metrics:      metrics,
		Suggestions:  suggestions,
		AutoResolved: autoResolved,
	}
}

// HangupConfigChangedEvent 挂机配置变更事件
type HangupConfigChangedEvent struct {
	*BaseDomainEvent
	ConfigType      string                 `json:"config_type"`
	PreviousConfig  map[string]interface{} `json:"previous_config"`
	NewConfig       map[string]interface{} `json:"new_config"`
	ChangedBy       string                 `json:"changed_by"`
	ChangeReason    string                 `json:"change_reason"`
	AffectedPlayers []string               `json:"affected_players,omitempty"`
}

// NewHangupConfigChangedEvent 创建挂机配置变更事件
func NewHangupConfigChangedEvent(configType string, prevConfig, newConfig map[string]interface{}, changedBy, reason string, affectedPlayers []string) *HangupConfigChangedEvent {
	return &HangupConfigChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "HangupConfigChanged",
			AggregateID: "system",
			OccurredAt:  time.Now(),
		},
		ConfigType:      configType,
		PreviousConfig:  prevConfig,
		NewConfig:       newConfig,
		ChangedBy:       changedBy,
		ChangeReason:    reason,
		AffectedPlayers: affectedPlayers,
	}
}

// HangupSessionCompletedEvent 挂机会话完成事件
type HangupSessionCompletedEvent struct {
	*BaseDomainEvent
	SessionID    string        `json:"session_id"`
	PlayerID     string        `json:"player_id"`
	LocationID   string        `json:"location_id"`
	LocationName string        `json:"location_name"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
	Duration     time.Duration `json:"duration"`
	IsOnline     bool          `json:"is_online"`
	Reward       *BaseReward   `json:"reward"`
	Efficiency   float64       `json:"efficiency"`
	Quality      string        `json:"quality"` // "excellent", "good", "normal", "poor"
}

// NewHangupSessionCompletedEvent 创建挂机会话完成事件
func NewHangupSessionCompletedEvent(session *HangupSession, locationName string, efficiency float64, quality string) *HangupSessionCompletedEvent {
	return &HangupSessionCompletedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "HangupSessionCompleted",
			AggregateID: session.GetPlayerID(),
			OccurredAt:  time.Now(),
		},
		SessionID:    session.GetSessionID(),
		PlayerID:     session.GetPlayerID(),
		LocationID:   session.GetLocationID(),
		LocationName: locationName,
		StartTime:    session.GetStartTime(),
		EndTime:      session.GetEndTime(),
		Duration:     session.GetDuration(),
		IsOnline:     session.IsOnlineSession(),
		Reward:       session.GetReward(),
		Efficiency:   efficiency,
		Quality:      quality,
	}
}

// generateEventID 生成事件ID
func generateEventID() string {
	// 这里可以使用UUID或其他唯一ID生成方式
	// 为了简化，使用时间戳
	return fmt.Sprintf("hangup_%d", time.Now().UnixNano())
}