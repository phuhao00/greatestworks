package honor

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

// TitleUnlockedEvent 称号解锁事件
type TitleUnlockedEvent struct {
	*BaseDomainEvent
	PlayerID    string `json:"player_id"`
	TitleID     string `json:"title_id"`
	TitleName   string `json:"title_name"`
	TitleRarity TitleRarity `json:"title_rarity"`
}

// NewTitleUnlockedEvent 创建称号解锁事件
func NewTitleUnlockedEvent(playerID, titleID, titleName string, rarity TitleRarity) *TitleUnlockedEvent {
	return &TitleUnlockedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "TitleUnlocked",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:    playerID,
		TitleID:     titleID,
		TitleName:   titleName,
		TitleRarity: rarity,
	}
}

// TitleEquippedEvent 称号装备事件
type TitleEquippedEvent struct {
	*BaseDomainEvent
	PlayerID      string `json:"player_id"`
	TitleID       string `json:"title_id"`
	TitleName     string `json:"title_name"`
	PreviousTitle string `json:"previous_title,omitempty"`
}

// NewTitleEquippedEvent 创建称号装备事件
func NewTitleEquippedEvent(playerID, titleID, titleName, previousTitle string) *TitleEquippedEvent {
	return &TitleEquippedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "TitleEquipped",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:      playerID,
		TitleID:       titleID,
		TitleName:     titleName,
		PreviousTitle: previousTitle,
	}
}

// TitleUnequippedEvent 称号卸下事件
type TitleUnequippedEvent struct {
	*BaseDomainEvent
	PlayerID  string `json:"player_id"`
	TitleID   string `json:"title_id"`
	TitleName string `json:"title_name"`
}

// NewTitleUnequippedEvent 创建称号卸下事件
func NewTitleUnequippedEvent(playerID, titleID, titleName string) *TitleUnequippedEvent {
	return &TitleUnequippedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "TitleUnequipped",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:  playerID,
		TitleID:   titleID,
		TitleName: titleName,
	}
}

// AchievementUnlockedEvent 成就解锁事件
type AchievementUnlockedEvent struct {
	*BaseDomainEvent
	PlayerID        string              `json:"player_id"`
	AchievementID   string              `json:"achievement_id"`
	AchievementName string              `json:"achievement_name"`
	Category        AchievementCategory `json:"category"`
	HonorReward     int                 `json:"honor_reward"`
	ItemRewards     []string            `json:"item_rewards"`
}

// NewAchievementUnlockedEvent 创建成就解锁事件
func NewAchievementUnlockedEvent(playerID, achievementID, achievementName string, category AchievementCategory, honorReward int, itemRewards []string) *AchievementUnlockedEvent {
	return &AchievementUnlockedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "AchievementUnlocked",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:        playerID,
		AchievementID:   achievementID,
		AchievementName: achievementName,
		Category:        category,
		HonorReward:     honorReward,
		ItemRewards:     itemRewards,
	}
}

// HonorLevelUpEvent 荣誉等级提升事件
type HonorLevelUpEvent struct {
	*BaseDomainEvent
	PlayerID     string `json:"player_id"`
	PreviousLevel int   `json:"previous_level"`
	NewLevel     int    `json:"new_level"`
	HonorPoints  int    `json:"honor_points"`
	LevelTitle   string `json:"level_title"`
}

// NewHonorLevelUpEvent 创建荣誉等级提升事件
func NewHonorLevelUpEvent(playerID string, previousLevel, newLevel, honorPoints int, levelTitle string) *HonorLevelUpEvent {
	return &HonorLevelUpEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "HonorLevelUp",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:     playerID,
		PreviousLevel: previousLevel,
		NewLevel:     newLevel,
		HonorPoints:  honorPoints,
		LevelTitle:   levelTitle,
	}
}

// HonorPointsEarnedEvent 荣誉点数获得事件
type HonorPointsEarnedEvent struct {
	*BaseDomainEvent
	PlayerID      string `json:"player_id"`
	PointsEarned  int    `json:"points_earned"`
	TotalPoints   int    `json:"total_points"`
	Source        string `json:"source"` // 来源：achievement, quest, event等
	SourceID      string `json:"source_id,omitempty"`
}

// NewHonorPointsEarnedEvent 创建荣誉点数获得事件
func NewHonorPointsEarnedEvent(playerID string, pointsEarned, totalPoints int, source, sourceID string) *HonorPointsEarnedEvent {
	return &HonorPointsEarnedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "HonorPointsEarned",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:     playerID,
		PointsEarned: pointsEarned,
		TotalPoints:  totalPoints,
		Source:       source,
		SourceID:     sourceID,
	}
}

// ReputationChangedEvent 声望变化事件
type ReputationChangedEvent struct {
	*BaseDomainEvent
	PlayerID         string `json:"player_id"`
	Faction          string `json:"faction"`
	PreviousReputation int  `json:"previous_reputation"`
	NewReputation    int    `json:"new_reputation"`
	Change           int    `json:"change"`
	Reason           string `json:"reason"`
}

// NewReputationChangedEvent 创建声望变化事件
func NewReputationChangedEvent(playerID, faction string, previousRep, newRep int, reason string) *ReputationChangedEvent {
	return &ReputationChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "ReputationChanged",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:           playerID,
		Faction:            faction,
		PreviousReputation: previousRep,
		NewReputation:      newRep,
		Change:             newRep - previousRep,
		Reason:             reason,
	}
}

// StatisticUpdatedEvent 统计数据更新事件
type StatisticUpdatedEvent struct {
	*BaseDomainEvent
	PlayerID      string        `json:"player_id"`
	StatisticType StatisticType `json:"statistic_type"`
	PreviousValue int           `json:"previous_value"`
	NewValue      int           `json:"new_value"`
	Change        int           `json:"change"`
}

// NewStatisticUpdatedEvent 创建统计数据更新事件
func NewStatisticUpdatedEvent(playerID string, statType StatisticType, previousValue, newValue int) *StatisticUpdatedEvent {
	return &StatisticUpdatedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "StatisticUpdated",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:      playerID,
		StatisticType: statType,
		PreviousValue: previousValue,
		NewValue:      newValue,
		Change:        newValue - previousValue,
	}
}

// HonorSystemInitializedEvent 荣誉系统初始化事件
type HonorSystemInitializedEvent struct {
	*BaseDomainEvent
	PlayerID         string `json:"player_id"`
	InitialLevel     int    `json:"initial_level"`
	InitialPoints    int    `json:"initial_points"`
	TitlesCount      int    `json:"titles_count"`
	AchievementsCount int   `json:"achievements_count"`
}

// NewHonorSystemInitializedEvent 创建荣誉系统初始化事件
func NewHonorSystemInitializedEvent(playerID string, initialLevel, initialPoints, titlesCount, achievementsCount int) *HonorSystemInitializedEvent {
	return &HonorSystemInitializedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "HonorSystemInitialized",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:          playerID,
		InitialLevel:      initialLevel,
		InitialPoints:     initialPoints,
		TitlesCount:       titlesCount,
		AchievementsCount: achievementsCount,
	}
}

// RareAchievementUnlockedEvent 稀有成就解锁事件（特殊事件）
type RareAchievementUnlockedEvent struct {
	*BaseDomainEvent
	PlayerID        string              `json:"player_id"`
	PlayerName      string              `json:"player_name"`
	AchievementID   string              `json:"achievement_id"`
	AchievementName string              `json:"achievement_name"`
	Category        AchievementCategory `json:"category"`
	Rarity          string              `json:"rarity"`
	UnlockRate      float64             `json:"unlock_rate"` // 解锁率，用于判断稀有度
}

// NewRareAchievementUnlockedEvent 创建稀有成就解锁事件
func NewRareAchievementUnlockedEvent(playerID, playerName, achievementID, achievementName string, category AchievementCategory, rarity string, unlockRate float64) *RareAchievementUnlockedEvent {
	return &RareAchievementUnlockedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "RareAchievementUnlocked",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:        playerID,
		PlayerName:      playerName,
		AchievementID:   achievementID,
		AchievementName: achievementName,
		Category:        category,
		Rarity:          rarity,
		UnlockRate:      unlockRate,
	}
}

// HonorRankingChangedEvent 荣誉排名变化事件
type HonorRankingChangedEvent struct {
	*BaseDomainEvent
	PlayerID     string `json:"player_id"`
	PlayerName   string `json:"player_name"`
	PreviousRank int    `json:"previous_rank"`
	NewRank      int    `json:"new_rank"`
	HonorPoints  int    `json:"honor_points"`
	HonorLevel   int    `json:"honor_level"`
}

// NewHonorRankingChangedEvent 创建荣誉排名变化事件
func NewHonorRankingChangedEvent(playerID, playerName string, previousRank, newRank, honorPoints, honorLevel int) *HonorRankingChangedEvent {
	return &HonorRankingChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     generateEventID(),
			EventType:   "HonorRankingChanged",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:     playerID,
		PlayerName:   playerName,
		PreviousRank: previousRank,
		NewRank:      newRank,
		HonorPoints:  honorPoints,
		HonorLevel:   honorLevel,
	}
}

// generateEventID 生成事件ID
func generateEventID() string {
	// 这里可以使用UUID或其他唯一ID生成方式
	// 为了简化，使用时间戳
	return fmt.Sprintf("honor_%d", time.Now().UnixNano())
}