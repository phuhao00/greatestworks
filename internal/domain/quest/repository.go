package quest

import (
	"context"
	"time"
)

// Repository 任务仓储接口
type Repository interface {
	// 基础CRUD操作
	Save(ctx context.Context, questManager *QuestManager) error
	FindByPlayerID(ctx context.Context, playerID string) (*QuestManager, error)
	Delete(ctx context.Context, playerID string) error
	Exists(ctx context.Context, playerID string) (bool, error)
	
	// 批量操作
	SaveBatch(ctx context.Context, questManagers []*QuestManager) error
	FindByPlayerIDs(ctx context.Context, playerIDs []string) ([]*QuestManager, error)
	
	// 任务查询
	FindQuestsByType(ctx context.Context, playerID string, questType QuestType) ([]*Quest, error)
	FindActiveQuests(ctx context.Context, playerID string) ([]*Quest, error)
	FindCompletedQuests(ctx context.Context, playerID string) ([]*Quest, error)
	FindAvailableQuests(ctx context.Context, playerID string) ([]*Quest, error)
	FindExpiredQuests(ctx context.Context, playerID string) ([]*Quest, error)
	
	// 成就查询
	FindAchievements(ctx context.Context, playerID string) ([]*Achievement, error)
	FindUnlockedAchievements(ctx context.Context, playerID string) ([]*Achievement, error)
	FindAchievementsByCategory(ctx context.Context, playerID string, category AchievementCategory) ([]*Achievement, error)
	
	// 统计查询
	GetQuestStats(ctx context.Context, playerID string) (*QuestStats, error)
	GetAchievementStats(ctx context.Context, playerID string) (*AchievementStats, error)
	GetQuestHistory(ctx context.Context, playerID string, limit int) ([]*QuestHistoryRecord, error)
	
	// 配置管理
	GetQuestConfig(ctx context.Context, questID string) (*QuestConfig, error)
	GetAllQuestConfigs(ctx context.Context) ([]*QuestConfig, error)
	SaveQuestConfig(ctx context.Context, config *QuestConfig) error
	GetAchievementConfig(ctx context.Context, achievementID string) (*AchievementConfig, error)
	GetAllAchievementConfigs(ctx context.Context) ([]*AchievementConfig, error)
	SaveAchievementConfig(ctx context.Context, config *AchievementConfig) error
}

// QuestStats 任务统计信息
type QuestStats struct {
	PlayerID         string                    `json:"player_id"`
	TotalQuests      int                       `json:"total_quests"`
	ActiveQuests     int                       `json:"active_quests"`
	CompletedQuests  int                       `json:"completed_quests"`
	FailedQuests     int                       `json:"failed_quests"`
	AbandonedQuests  int                       `json:"abandoned_quests"`
	QuestsByType     map[QuestType]int         `json:"quests_by_type"`
	QuestsByCategory map[QuestCategory]int     `json:"quests_by_category"`
	CompletionRate   float64                   `json:"completion_rate"`
	AverageTime      time.Duration             `json:"average_completion_time"`
	LastUpdate       time.Time                 `json:"last_update"`
}

// AchievementStats 成就统计信息
type AchievementStats struct {
	PlayerID            string                           `json:"player_id"`
	TotalAchievements   int                              `json:"total_achievements"`
	UnlockedAchievements int                             `json:"unlocked_achievements"`
	TotalPoints         int64                            `json:"total_points"`
	AchievementsByCategory map[AchievementCategory]int   `json:"achievements_by_category"`
	CompletionRate      float64                          `json:"completion_rate"`
	RareAchievements    int                              `json:"rare_achievements"`
	HiddenAchievements  int                              `json:"hidden_achievements"`
	LastUnlocked        *time.Time                       `json:"last_unlocked"`
	LastUpdate          time.Time                        `json:"last_update"`
}

// QuestHistoryRecord 任务历史记录
type QuestHistoryRecord struct {
	ID           string      `json:"id"`
	PlayerID     string      `json:"player_id"`
	QuestID      string      `json:"quest_id"`
	QuestName    string      `json:"quest_name"`
	QuestType    QuestType   `json:"quest_type"`
	Action       string      `json:"action"` // accepted, completed, failed, abandoned
	StartTime    *time.Time  `json:"start_time"`
	EndTime      *time.Time  `json:"end_time"`
	Duration     *time.Duration `json:"duration"`
	Rewards      []*QuestReward `json:"rewards"`
	OccurredAt   time.Time   `json:"occurred_at"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// QuestConfig 任务配置
type QuestConfig struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	QuestType         QuestType              `json:"quest_type"`
	Category          QuestCategory          `json:"category"`
	Priority          QuestPriority          `json:"priority"`
	Objectives        []*ObjectiveConfig     `json:"objectives"`
	Rewards           []*QuestReward         `json:"rewards"`
	Prerequisites     []string               `json:"prerequisites"`
	TimeLimit         *time.Duration         `json:"time_limit"`
	RepeatType        RepeatType             `json:"repeat_type"`
	MaxRepeats        int                    `json:"max_repeats"`
	Level             int                    `json:"level"`
	MinLevel          int                    `json:"min_level"`
	MaxLevel          int                    `json:"max_level"`
	ClassRestrictions []string               `json:"class_restrictions"`
	RaceRestrictions  []string               `json:"race_restrictions"`
	Enabled           bool                   `json:"enabled"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

// ObjectiveConfig 目标配置
type ObjectiveConfig struct {
	ID            string                 `json:"id"`
	Description   string                 `json:"description"`
	ObjectiveType ObjectiveType          `json:"objective_type"`
	Target        string                 `json:"target"`
	Required      int64                  `json:"required"`
	Optional      bool                   `json:"optional"`
	Order         int                    `json:"order"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// AchievementConfig 成就配置
type AchievementConfig struct {
	ID           string                     `json:"id"`
	Name         string                     `json:"name"`
	Description  string                     `json:"description"`
	Category     AchievementCategory        `json:"category"`
	Points       int64                      `json:"points"`
	Requirements []*RequirementConfig       `json:"requirements"`
	Rewards      []*QuestReward             `json:"rewards"`
	Hidden       bool                       `json:"hidden"`
	Rare         bool                       `json:"rare"`
	Enabled      bool                       `json:"enabled"`
	CreatedAt    time.Time                  `json:"created_at"`
	UpdatedAt    time.Time                  `json:"updated_at"`
}

// RequirementConfig 要求配置
type RequirementConfig struct {
	RequirementType RequirementType `json:"requirement_type"`
	Target          string          `json:"target"`
	Value           int64           `json:"value"`
}

// QuestQueryFilter 任务查询过滤器
type QuestQueryFilter struct {
	PlayerID          string          `json:"player_id"`
	QuestTypes        []QuestType     `json:"quest_types,omitempty"`
	Categories        []QuestCategory `json:"categories,omitempty"`
	Statuses          []QuestStatus   `json:"statuses,omitempty"`
	Priorities        []QuestPriority `json:"priorities,omitempty"`
	MinLevel          *int            `json:"min_level,omitempty"`
	MaxLevel          *int            `json:"max_level,omitempty"`
	AvailableOnly     bool            `json:"available_only"`
	ActiveOnly        bool            `json:"active_only"`
	CompletedOnly     bool            `json:"completed_only"`
	IncludeExpired    bool            `json:"include_expired"`
	SortBy            string          `json:"sort_by"` // priority, level, created_at
	SortOrder         string          `json:"sort_order"` // asc, desc
	Limit             int             `json:"limit"`
	Offset            int             `json:"offset"`
}

// AchievementQueryFilter 成就查询过滤器
type AchievementQueryFilter struct {
	PlayerID       string                    `json:"player_id"`
	Categories     []AchievementCategory     `json:"categories,omitempty"`
	UnlockedOnly   bool                      `json:"unlocked_only"`
	LockedOnly     bool                      `json:"locked_only"`
	HiddenOnly     bool                      `json:"hidden_only"`
	RareOnly       bool                      `json:"rare_only"`
	MinPoints      *int64                    `json:"min_points,omitempty"`
	MaxPoints      *int64                    `json:"max_points,omitempty"`
	SortBy         string                    `json:"sort_by"` // points, unlocked_at, category
	SortOrder      string                    `json:"sort_order"` // asc, desc
	Limit          int                       `json:"limit"`
	Offset         int                       `json:"offset"`
}