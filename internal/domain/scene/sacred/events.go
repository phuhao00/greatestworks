package sacred

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
	GetVersion() int
	GetPayload() map[string]interface{}
}

// BaseDomainEvent 基础领域事件
type BaseDomainEvent struct {
	EventID     string
	EventType   string
	AggregateID string
	OccurredAt  time.Time
	Version     int
	Payload     map[string]interface{}
}

// GetEventID 获取事件ID
func (e *BaseDomainEvent) GetEventID() string {
	return e.EventID
}

// GetEventType 获取事件类型
func (e *BaseDomainEvent) GetEventType() string {
	return e.EventType
}

// GetAggregateID 获取聚合ID
func (e *BaseDomainEvent) GetAggregateID() string {
	return e.AggregateID
}

// GetOccurredAt 获取发生时间
func (e *BaseDomainEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetVersion 获取版本
func (e *BaseDomainEvent) GetVersion() int {
	return e.Version
}

// GetPayload 获取载荷
func (e *BaseDomainEvent) GetPayload() map[string]interface{} {
	return e.Payload
}

// SacredNameChangedEvent 圣地名称变更事件
type SacredNameChangedEvent struct {
	*BaseDomainEvent
	SacredID string
	OldName  string
	NewName  string
}

// NewSacredNameChangedEvent 创建圣地名称变更事件
func NewSacredNameChangedEvent(sacredID, oldName, newName string) *SacredNameChangedEvent {
	now := time.Now()
	event := &SacredNameChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("sacred_name_changed_%d", now.UnixNano()),
			EventType:   "sacred.name_changed",
			AggregateID: sacredID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SacredID: sacredID,
		OldName:  oldName,
		NewName:  newName,
	}

	// 设置载荷
	event.Payload["sacred_id"] = sacredID
	event.Payload["old_name"] = oldName
	event.Payload["new_name"] = newName

	return event
}

// SacredLevelUpEvent 圣地升级事件
type SacredLevelUpEvent struct {
	*BaseDomainEvent
	SacredID   string
	OldLevel   int
	NewLevel   int
	Experience int
	Rewards    map[string]interface{}
}

// NewSacredLevelUpEvent 创建圣地升级事件
func NewSacredLevelUpEvent(sacredID string, oldLevel, newLevel, experience int) *SacredLevelUpEvent {
	now := time.Now()
	event := &SacredLevelUpEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("sacred_level_up_%d", now.UnixNano()),
			EventType:   "sacred.level_up",
			AggregateID: sacredID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SacredID:   sacredID,
		OldLevel:   oldLevel,
		NewLevel:   newLevel,
		Experience: experience,
		Rewards:    make(map[string]interface{}),
	}

	// 设置载荷
	event.Payload["sacred_id"] = sacredID
	event.Payload["old_level"] = oldLevel
	event.Payload["new_level"] = newLevel
	event.Payload["experience"] = experience
	event.Payload["level_gain"] = newLevel - oldLevel

	return event
}

// ChallengeAddedEvent 挑战添加事件
type ChallengeAddedEvent struct {
	*BaseDomainEvent
	SacredID      string
	ChallengeID   string
	ChallengeType ChallengeType
	Difficulty    ChallengeDifficulty
}

// NewChallengeAddedEvent 创建挑战添加事件
func NewChallengeAddedEvent(sacredID, challengeID string, challengeType ChallengeType, difficulty ChallengeDifficulty) *ChallengeAddedEvent {
	now := time.Now()
	event := &ChallengeAddedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("challenge_added_%d", now.UnixNano()),
			EventType:   "sacred.challenge_added",
			AggregateID: sacredID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SacredID:      sacredID,
		ChallengeID:   challengeID,
		ChallengeType: challengeType,
		Difficulty:    difficulty,
	}

	// 设置载荷
	event.Payload["sacred_id"] = sacredID
	event.Payload["challenge_id"] = challengeID
	event.Payload["challenge_type"] = challengeType.String()
	event.Payload["difficulty"] = difficulty.String()

	return event
}

// ChallengeRemovedEvent 挑战移除事件
type ChallengeRemovedEvent struct {
	*BaseDomainEvent
	SacredID      string
	ChallengeID   string
	ChallengeType ChallengeType
	Reason        string
}

// NewChallengeRemovedEvent 创建挑战移除事件
func NewChallengeRemovedEvent(sacredID, challengeID string, challengeType ChallengeType) *ChallengeRemovedEvent {
	now := time.Now()
	event := &ChallengeRemovedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("challenge_removed_%d", now.UnixNano()),
			EventType:   "sacred.challenge_removed",
			AggregateID: sacredID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SacredID:      sacredID,
		ChallengeID:   challengeID,
		ChallengeType: challengeType,
		Reason:        "manual_removal",
	}

	// 设置载荷
	event.Payload["sacred_id"] = sacredID
	event.Payload["challenge_id"] = challengeID
	event.Payload["challenge_type"] = challengeType.String()
	event.Payload["reason"] = event.Reason

	return event
}

// ChallengeStartedEvent 挑战开始事件
type ChallengeStartedEvent struct {
	*BaseDomainEvent
	SacredID      string
	ChallengeID   string
	PlayerID      string
	ChallengeType ChallengeType
	StartTime     time.Time
}

// NewChallengeStartedEvent 创建挑战开始事件
func NewChallengeStartedEvent(sacredID, challengeID, playerID string, challengeType ChallengeType) *ChallengeStartedEvent {
	now := time.Now()
	event := &ChallengeStartedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("challenge_started_%d", now.UnixNano()),
			EventType:   "sacred.challenge_started",
			AggregateID: sacredID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SacredID:      sacredID,
		ChallengeID:   challengeID,
		PlayerID:      playerID,
		ChallengeType: challengeType,
		StartTime:     now,
	}

	// 设置载荷
	event.Payload["sacred_id"] = sacredID
	event.Payload["challenge_id"] = challengeID
	event.Payload["player_id"] = playerID
	event.Payload["challenge_type"] = challengeType.String()
	event.Payload["start_time"] = now

	return event
}

// ChallengeCompletedEvent 挑战完成事件
type ChallengeCompletedEvent struct {
	*BaseDomainEvent
	SacredID       string
	ChallengeID    string
	PlayerID       string
	Success        bool
	Score          int
	Reward         *ChallengeReward
	CompletionTime time.Time
	Duration       time.Duration
}

// NewChallengeCompletedEvent 创建挑战完成事件
func NewChallengeCompletedEvent(sacredID, challengeID, playerID string, success bool, score int, reward *ChallengeReward) *ChallengeCompletedEvent {
	now := time.Now()
	event := &ChallengeCompletedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("challenge_completed_%d", now.UnixNano()),
			EventType:   "sacred.challenge_completed",
			AggregateID: sacredID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SacredID:       sacredID,
		ChallengeID:    challengeID,
		PlayerID:       playerID,
		Success:        success,
		Score:          score,
		Reward:         reward,
		CompletionTime: now,
	}

	// 设置载荷
	event.Payload["sacred_id"] = sacredID
	event.Payload["challenge_id"] = challengeID
	event.Payload["player_id"] = playerID
	event.Payload["success"] = success
	event.Payload["score"] = score
	event.Payload["completion_time"] = now
	if reward != nil {
		event.Payload["reward_gold"] = reward.Gold
		event.Payload["reward_experience"] = reward.Experience
		event.Payload["reward_items"] = len(reward.Items)
	}

	return event
}

// BlessingAddedEvent 祝福添加事件
type BlessingAddedEvent struct {
	*BaseDomainEvent
	SacredID     string
	BlessingID   string
	BlessingType BlessingType
	Duration     time.Duration
}

// NewBlessingAddedEvent 创建祝福添加事件
func NewBlessingAddedEvent(sacredID, blessingID string, blessingType BlessingType, duration time.Duration) *BlessingAddedEvent {
	now := time.Now()
	event := &BlessingAddedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("blessing_added_%d", now.UnixNano()),
			EventType:   "sacred.blessing_added",
			AggregateID: sacredID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SacredID:     sacredID,
		BlessingID:   blessingID,
		BlessingType: blessingType,
		Duration:     duration,
	}

	// 设置载荷
	event.Payload["sacred_id"] = sacredID
	event.Payload["blessing_id"] = blessingID
	event.Payload["blessing_type"] = blessingType.String()
	event.Payload["duration"] = duration.String()

	return event
}

// BlessingRemovedEvent 祝福移除事件
type BlessingRemovedEvent struct {
	*BaseDomainEvent
	SacredID     string
	BlessingID   string
	BlessingType BlessingType
	Reason       string
}

// NewBlessingRemovedEvent 创建祝福移除事件
func NewBlessingRemovedEvent(sacredID, blessingID string, blessingType BlessingType) *BlessingRemovedEvent {
	now := time.Now()
	event := &BlessingRemovedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("blessing_removed_%d", now.UnixNano()),
			EventType:   "sacred.blessing_removed",
			AggregateID: sacredID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SacredID:     sacredID,
		BlessingID:   blessingID,
		BlessingType: blessingType,
		Reason:       "manual_removal",
	}

	// 设置载荷
	event.Payload["sacred_id"] = sacredID
	event.Payload["blessing_id"] = blessingID
	event.Payload["blessing_type"] = blessingType.String()
	event.Payload["reason"] = event.Reason

	return event
}

// BlessingActivatedEvent 祝福激活事件
type BlessingActivatedEvent struct {
	*BaseDomainEvent
	SacredID     string
	BlessingID   string
	PlayerID     string
	BlessingType BlessingType
	Effect       *BlessingEffect
	ActivatedAt  time.Time
	ExpiresAt    time.Time
}

// NewBlessingActivatedEvent 创建祝福激活事件
func NewBlessingActivatedEvent(sacredID, blessingID, playerID string, blessingType BlessingType, effect *BlessingEffect) *BlessingActivatedEvent {
	now := time.Now()
	event := &BlessingActivatedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("blessing_activated_%d", now.UnixNano()),
			EventType:   "sacred.blessing_activated",
			AggregateID: sacredID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SacredID:     sacredID,
		BlessingID:   blessingID,
		PlayerID:     playerID,
		BlessingType: blessingType,
		Effect:       effect,
		ActivatedAt:  now,
		ExpiresAt:    effect.ExpiresAt,
	}

	// 设置载荷
	event.Payload["sacred_id"] = sacredID
	event.Payload["blessing_id"] = blessingID
	event.Payload["player_id"] = playerID
	event.Payload["blessing_type"] = blessingType.String()
	event.Payload["activated_at"] = now
	event.Payload["expires_at"] = effect.ExpiresAt
	event.Payload["duration"] = effect.ExpiresAt.Sub(now).String()

	return event
}

// SacredStatusChangedEvent 圣地状态变更事件
type SacredStatusChangedEvent struct {
	*BaseDomainEvent
	SacredID  string
	OldStatus SacredStatus
	NewStatus SacredStatus
	Reason    string
}

// NewSacredStatusChangedEvent 创建圣地状态变更事件
func NewSacredStatusChangedEvent(sacredID string, oldStatus, newStatus SacredStatus) *SacredStatusChangedEvent {
	now := time.Now()
	event := &SacredStatusChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("sacred_status_changed_%d", now.UnixNano()),
			EventType:   "sacred.status_changed",
			AggregateID: sacredID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SacredID:  sacredID,
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Reason:    "status_update",
	}

	// 设置载荷
	event.Payload["sacred_id"] = sacredID
	event.Payload["old_status"] = oldStatus.String()
	event.Payload["new_status"] = newStatus.String()
	event.Payload["reason"] = event.Reason

	return event
}

// RelicObtainedEvent 圣物获得事件
type RelicObtainedEvent struct {
	*BaseDomainEvent
	SacredID   string
	PlayerID   string
	RelicID    string
	RelicType  RelicType
	Rarity     RelicRarity
	Source     string
	ObtainedAt time.Time
}

// NewRelicObtainedEvent 创建圣物获得事件
func NewRelicObtainedEvent(sacredID, playerID, relicID string, relicType RelicType, rarity RelicRarity, source string) *RelicObtainedEvent {
	now := time.Now()
	event := &RelicObtainedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("relic_obtained_%d", now.UnixNano()),
			EventType:   "sacred.relic_obtained",
			AggregateID: sacredID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SacredID:   sacredID,
		PlayerID:   playerID,
		RelicID:    relicID,
		RelicType:  relicType,
		Rarity:     rarity,
		Source:     source,
		ObtainedAt: now,
	}

	// 设置载荷
	event.Payload["sacred_id"] = sacredID
	event.Payload["player_id"] = playerID
	event.Payload["relic_id"] = relicID
	event.Payload["relic_type"] = relicType.String()
	event.Payload["rarity"] = rarity.String()
	event.Payload["source"] = source
	event.Payload["obtained_at"] = now

	return event
}

// RelicUpgradedEvent 圣物升级事件
type RelicUpgradedEvent struct {
	*BaseDomainEvent
	SacredID   string
	PlayerID   string
	RelicID    string
	OldLevel   int
	NewLevel   int
	OldPower   float64
	NewPower   float64
	UpgradedAt time.Time
}

// NewRelicUpgradedEvent 创建圣物升级事件
func NewRelicUpgradedEvent(sacredID, playerID, relicID string, oldLevel, newLevel int, oldPower, newPower float64) *RelicUpgradedEvent {
	now := time.Now()
	event := &RelicUpgradedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("relic_upgraded_%d", now.UnixNano()),
			EventType:   "sacred.relic_upgraded",
			AggregateID: sacredID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SacredID:   sacredID,
		PlayerID:   playerID,
		RelicID:    relicID,
		OldLevel:   oldLevel,
		NewLevel:   newLevel,
		OldPower:   oldPower,
		NewPower:   newPower,
		UpgradedAt: now,
	}

	// 设置载荷
	event.Payload["sacred_id"] = sacredID
	event.Payload["player_id"] = playerID
	event.Payload["relic_id"] = relicID
	event.Payload["old_level"] = oldLevel
	event.Payload["new_level"] = newLevel
	event.Payload["level_gain"] = newLevel - oldLevel
	event.Payload["old_power"] = oldPower
	event.Payload["new_power"] = newPower
	event.Payload["power_gain"] = newPower - oldPower
	event.Payload["upgraded_at"] = now

	return event
}

// PlayerEnteredSacredEvent 玩家进入圣地事件
type PlayerEnteredSacredEvent struct {
	*BaseDomainEvent
	SacredID  string
	PlayerID  string
	EnteredAt time.Time
	Source    string
}

// NewPlayerEnteredSacredEvent 创建玩家进入圣地事件
func NewPlayerEnteredSacredEvent(sacredID, playerID, source string) *PlayerEnteredSacredEvent {
	now := time.Now()
	event := &PlayerEnteredSacredEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("player_entered_sacred_%d", now.UnixNano()),
			EventType:   "sacred.player_entered",
			AggregateID: sacredID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SacredID:  sacredID,
		PlayerID:  playerID,
		EnteredAt: now,
		Source:    source,
	}

	// 设置载荷
	event.Payload["sacred_id"] = sacredID
	event.Payload["player_id"] = playerID
	event.Payload["entered_at"] = now
	event.Payload["source"] = source

	return event
}

// PlayerLeftSacredEvent 玩家离开圣地事件
type PlayerLeftSacredEvent struct {
	*BaseDomainEvent
	SacredID   string
	PlayerID   string
	LeftAt     time.Time
	Duration   time.Duration
	Activities []string
}

// NewPlayerLeftSacredEvent 创建玩家离开圣地事件
func NewPlayerLeftSacredEvent(sacredID, playerID string, enteredAt time.Time, activities []string) *PlayerLeftSacredEvent {
	now := time.Now()
	duration := now.Sub(enteredAt)

	event := &PlayerLeftSacredEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("player_left_sacred_%d", now.UnixNano()),
			EventType:   "sacred.player_left",
			AggregateID: sacredID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SacredID:   sacredID,
		PlayerID:   playerID,
		LeftAt:     now,
		Duration:   duration,
		Activities: activities,
	}

	// 设置载荷
	event.Payload["sacred_id"] = sacredID
	event.Payload["player_id"] = playerID
	event.Payload["left_at"] = now
	event.Payload["duration"] = duration.String()
	event.Payload["activities"] = activities
	event.Payload["activity_count"] = len(activities)

	return event
}

// SacredMaintenanceEvent 圣地维护事件
type SacredMaintenanceEvent struct {
	*BaseDomainEvent
	SacredID         string
	MaintenanceType  string
	StartTime        time.Time
	EstimatedEnd     time.Time
	Reason           string
	AffectedFeatures []string
}

// NewSacredMaintenanceEvent 创建圣地维护事件
func NewSacredMaintenanceEvent(sacredID, maintenanceType, reason string, duration time.Duration, affectedFeatures []string) *SacredMaintenanceEvent {
	now := time.Now()
	event := &SacredMaintenanceEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("sacred_maintenance_%d", now.UnixNano()),
			EventType:   "sacred.maintenance",
			AggregateID: sacredID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SacredID:         sacredID,
		MaintenanceType:  maintenanceType,
		StartTime:        now,
		EstimatedEnd:     now.Add(duration),
		Reason:           reason,
		AffectedFeatures: affectedFeatures,
	}

	// 设置载荷
	event.Payload["sacred_id"] = sacredID
	event.Payload["maintenance_type"] = maintenanceType
	event.Payload["start_time"] = now
	event.Payload["estimated_end"] = event.EstimatedEnd
	event.Payload["duration"] = duration.String()
	event.Payload["reason"] = reason
	event.Payload["affected_features"] = affectedFeatures
	event.Payload["affected_count"] = len(affectedFeatures)

	return event
}

// SacredAchievementUnlockedEvent 圣地成就解锁事件
type SacredAchievementUnlockedEvent struct {
	*BaseDomainEvent
	SacredID        string
	PlayerID        string
	AchievementID   string
	AchievementName string
	Description     string
	Rewards         map[string]interface{}
	UnlockedAt      time.Time
}

// NewSacredAchievementUnlockedEvent 创建圣地成就解锁事件
func NewSacredAchievementUnlockedEvent(sacredID, playerID, achievementID, achievementName, description string, rewards map[string]interface{}) *SacredAchievementUnlockedEvent {
	now := time.Now()
	event := &SacredAchievementUnlockedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("sacred_achievement_unlocked_%d", now.UnixNano()),
			EventType:   "sacred.achievement_unlocked",
			AggregateID: sacredID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SacredID:        sacredID,
		PlayerID:        playerID,
		AchievementID:   achievementID,
		AchievementName: achievementName,
		Description:     description,
		Rewards:         rewards,
		UnlockedAt:      now,
	}

	// 设置载荷
	event.Payload["sacred_id"] = sacredID
	event.Payload["player_id"] = playerID
	event.Payload["achievement_id"] = achievementID
	event.Payload["achievement_name"] = achievementName
	event.Payload["description"] = description
	event.Payload["rewards"] = rewards
	event.Payload["unlocked_at"] = now

	return event
}

// SacredSeasonChangedEvent 圣地季节变化事件
type SacredSeasonChangedEvent struct {
	*BaseDomainEvent
	SacredID      string
	OldSeason     string
	NewSeason     string
	SeasonEffects map[string]interface{}
	ChangedAt     time.Time
}

// NewSacredSeasonChangedEvent 创建圣地季节变化事件
func NewSacredSeasonChangedEvent(sacredID, oldSeason, newSeason string, seasonEffects map[string]interface{}) *SacredSeasonChangedEvent {
	now := time.Now()
	event := &SacredSeasonChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("sacred_season_changed_%d", now.UnixNano()),
			EventType:   "sacred.season_changed",
			AggregateID: sacredID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SacredID:      sacredID,
		OldSeason:     oldSeason,
		NewSeason:     newSeason,
		SeasonEffects: seasonEffects,
		ChangedAt:     now,
	}

	// 设置载荷
	event.Payload["sacred_id"] = sacredID
	event.Payload["old_season"] = oldSeason
	event.Payload["new_season"] = newSeason
	event.Payload["season_effects"] = seasonEffects
	event.Payload["changed_at"] = now

	return event
}

// SacredRankingUpdatedEvent 圣地排名更新事件
type SacredRankingUpdatedEvent struct {
	*BaseDomainEvent
	SacredID    string
	RankingType string
	OldRank     int
	NewRank     int
	Score       float64
	UpdatedAt   time.Time
}

// NewSacredRankingUpdatedEvent 创建圣地排名更新事件
func NewSacredRankingUpdatedEvent(sacredID, rankingType string, oldRank, newRank int, score float64) *SacredRankingUpdatedEvent {
	now := time.Now()
	event := &SacredRankingUpdatedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("sacred_ranking_updated_%d", now.UnixNano()),
			EventType:   "sacred.ranking_updated",
			AggregateID: sacredID,
			OccurredAt:  now,
			Version:     1,
			Payload:     make(map[string]interface{}),
		},
		SacredID:    sacredID,
		RankingType: rankingType,
		OldRank:     oldRank,
		NewRank:     newRank,
		Score:       score,
		UpdatedAt:   now,
	}

	// 设置载荷
	event.Payload["sacred_id"] = sacredID
	event.Payload["ranking_type"] = rankingType
	event.Payload["old_rank"] = oldRank
	event.Payload["new_rank"] = newRank
	event.Payload["rank_change"] = newRank - oldRank
	event.Payload["score"] = score
	event.Payload["updated_at"] = now

	return event
}

// 事件处理器接口

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(event DomainEvent) error
	CanHandle(eventType string) bool
	GetHandlerName() string
}

// EventBus 事件总线接口
type EventBus interface {
	Publish(event DomainEvent) error
	PublishBatch(events []DomainEvent) error
	Subscribe(eventType string, handler EventHandler) error
	Unsubscribe(eventType string, handler EventHandler) error
	GetSubscribers(eventType string) []EventHandler
}

// EventStore 事件存储接口
type EventStore interface {
	Save(event DomainEvent) error
	SaveBatch(events []DomainEvent) error
	Load(aggregateID string) ([]DomainEvent, error)
	LoadFromVersion(aggregateID string, version int) ([]DomainEvent, error)
	LoadByEventType(eventType string, limit int) ([]DomainEvent, error)
	LoadByTimeRange(startTime, endTime time.Time) ([]DomainEvent, error)
	LoadByAggregateType(aggregateType string, limit int) ([]DomainEvent, error)
	Delete(eventID string) error
	DeleteByAggregate(aggregateID string) error
	Count() (int64, error)
	CountByEventType(eventType string) (int64, error)
	CountByAggregate(aggregateID string) (int64, error)
}

// EventProjector 事件投影器接口
type EventProjector interface {
	Project(event DomainEvent) error
	ProjectBatch(events []DomainEvent) error
	Rebuild(aggregateID string) error
	GetProjectionName() string
	GetLastProcessedVersion(aggregateID string) (int, error)
	SetLastProcessedVersion(aggregateID string, version int) error
}

// EventSnapshot 事件快照接口
type EventSnapshot interface {
	SaveSnapshot(aggregateID string, version int, data interface{}) error
	LoadSnapshot(aggregateID string) (interface{}, int, error)
	DeleteSnapshot(aggregateID string) error
	GetSnapshotFrequency() int
	ShouldCreateSnapshot(aggregateID string, currentVersion int) bool
}

// 事件中间件

// EventMiddleware 事件中间件接口
type EventMiddleware interface {
	Before(event DomainEvent) error
	After(event DomainEvent) error
	OnError(event DomainEvent, err error) error
}

// EventValidator 事件验证器
type EventValidator struct {
	rules map[string][]ValidationRule
}

// ValidationRule 验证规则
type ValidationRule interface {
	Validate(event DomainEvent) error
	GetRuleName() string
}

// NewEventValidator 创建事件验证器
func NewEventValidator() *EventValidator {
	return &EventValidator{
		rules: make(map[string][]ValidationRule),
	}
}

// AddRule 添加验证规则
func (ev *EventValidator) AddRule(eventType string, rule ValidationRule) {
	if ev.rules[eventType] == nil {
		ev.rules[eventType] = make([]ValidationRule, 0)
	}
	ev.rules[eventType] = append(ev.rules[eventType], rule)
}

// Validate 验证事件
func (ev *EventValidator) Validate(event DomainEvent) error {
	rules, exists := ev.rules[event.GetEventType()]
	if !exists {
		return nil // 没有规则则通过
	}

	for _, rule := range rules {
		if err := rule.Validate(event); err != nil {
			return fmt.Errorf("validation failed for rule %s: %w", rule.GetRuleName(), err)
		}
	}

	return nil
}

// EventMetrics 事件指标
type EventMetrics struct {
	EventType      string
	Count          int64
	LastOccurred   time.Time
	AverageSize    float64
	ProcessingTime time.Duration
}

// EventMonitor 事件监控器
type EventMonitor interface {
	RecordEvent(event DomainEvent, processingTime time.Duration)
	GetMetrics(eventType string) (*EventMetrics, error)
	GetAllMetrics() (map[string]*EventMetrics, error)
	Reset() error
}
