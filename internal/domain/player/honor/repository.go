package honor

import "context"

// HonorRepository 荣誉仓储接口
type HonorRepository interface {
	// Save 保存荣誉聚合根
	Save(ctx context.Context, honor *HonorAggregate) error

	// FindByPlayerID 根据玩家ID查找荣誉信息
	FindByPlayerID(ctx context.Context, playerID string) (*HonorAggregate, error)

	// Delete 删除荣誉信息
	Delete(ctx context.Context, playerID string) error

	// FindByHonorLevel 根据荣誉等级查找玩家
	FindByHonorLevel(ctx context.Context, level int, limit, offset int) ([]*HonorAggregate, error)

	// FindTopHonorPlayers 查找荣誉排行榜
	FindTopHonorPlayers(ctx context.Context, limit int) ([]*HonorAggregate, error)

	// CountByHonorLevel 统计指定荣誉等级的玩家数量
	CountByHonorLevel(ctx context.Context, level int) (int64, error)

	// FindByTitleEquipped 查找装备指定称号的玩家
	FindByTitleEquipped(ctx context.Context, titleID string, limit, offset int) ([]*HonorAggregate, error)

	// FindByAchievementUnlocked 查找解锁指定成就的玩家
	FindByAchievementUnlocked(ctx context.Context, achievementID string, limit, offset int) ([]*HonorAggregate, error)

	// UpdateStatistics 批量更新玩家统计数据
	UpdateStatistics(ctx context.Context, playerID string, statistics map[StatisticType]int) error

	// FindByReputation 根据声望查找玩家
	FindByReputation(ctx context.Context, faction string, minReputation int, limit, offset int) ([]*HonorAggregate, error)

	// GetPlayerRank 获取玩家荣誉排名
	GetPlayerRank(ctx context.Context, playerID string) (int, error)

	// FindRecentlyUnlockedTitles 查找最近解锁的称号
	FindRecentlyUnlockedTitles(ctx context.Context, playerID string, limit int) ([]*Title, error)

	// FindRecentlyUnlockedAchievements 查找最近解锁的成就
	FindRecentlyUnlockedAchievements(ctx context.Context, playerID string, limit int) ([]*Achievement, error)

	// GetHonorStatistics 获取荣誉系统统计信息
	GetHonorStatistics(ctx context.Context) (*HonorStatistics, error)
}

// TitleRepository 称号仓储接口
type TitleRepository interface {
	// SaveTemplate 保存称号模板
	SaveTemplate(ctx context.Context, template *TitleTemplate) error

	// FindTemplateByID 根据ID查找称号模板
	FindTemplateByID(ctx context.Context, id string) (*TitleTemplate, error)

	// FindAllTemplates 查找所有称号模板
	FindAllTemplates(ctx context.Context) ([]*TitleTemplate, error)

	// FindTemplatesByCategory 根据分类查找称号模板
	FindTemplatesByCategory(ctx context.Context, category TitleCategory) ([]*TitleTemplate, error)

	// FindTemplatesByRarity 根据稀有度查找称号模板
	FindTemplatesByRarity(ctx context.Context, rarity TitleRarity) ([]*TitleTemplate, error)

	// DeleteTemplate 删除称号模板
	DeleteTemplate(ctx context.Context, id string) error

	// GetTitleStatistics 获取称号统计信息
	GetTitleStatistics(ctx context.Context, titleID string) (*TitleStatistics, error)
}

// AchievementRepository 成就仓储接口
type AchievementRepository interface {
	// SaveTemplate 保存成就模板
	SaveTemplate(ctx context.Context, template *AchievementTemplate) error

	// FindTemplateByID 根据ID查找成就模板
	FindTemplateByID(ctx context.Context, id string) (*AchievementTemplate, error)

	// FindAllTemplates 查找所有成就模板
	FindAllTemplates(ctx context.Context) ([]*AchievementTemplate, error)

	// FindTemplatesByCategory 根据分类查找成就模板
	FindTemplatesByCategory(ctx context.Context, category AchievementCategory) ([]*AchievementTemplate, error)

	// FindTemplatesByType 根据类型查找成就模板
	FindTemplatesByType(ctx context.Context, achievementType AchievementType) ([]*AchievementTemplate, error)

	// DeleteTemplate 删除成就模板
	DeleteTemplate(ctx context.Context, id string) error

	// GetAchievementStatistics 获取成就统计信息
	GetAchievementStatistics(ctx context.Context, achievementID string) (*AchievementStatistics, error)
}

// HonorStatistics 荣誉系统统计信息
type HonorStatistics struct {
	TotalPlayers             int64              `json:"total_players"`
	AverageHonorPoints       float64            `json:"average_honor_points"`
	HonorLevelDistribution   map[int]int64      `json:"honor_level_distribution"`
	MostPopularTitles        []TitleStats       `json:"most_popular_titles"`
	MostUnlockedAchievements []AchievementStats `json:"most_unlocked_achievements"`
	TopReputationFactions    []ReputationStats  `json:"top_reputation_factions"`
}

// TitleStatistics 称号统计信息
type TitleStatistics struct {
	TitleID        string  `json:"title_id"`
	TitleName      string  `json:"title_name"`
	UnlockCount    int64   `json:"unlock_count"`
	EquipCount     int64   `json:"equip_count"`
	UnlockRate     float64 `json:"unlock_rate"`
	PopularityRank int     `json:"popularity_rank"`
}

// AchievementStatistics 成就统计信息
type AchievementStatistics struct {
	AchievementID     string  `json:"achievement_id"`
	AchievementName   string  `json:"achievement_name"`
	UnlockCount       int64   `json:"unlock_count"`
	UnlockRate        float64 `json:"unlock_rate"`
	AverageUnlockTime float64 `json:"average_unlock_time"` // 平均解锁时间（小时）
	DifficultyRank    int     `json:"difficulty_rank"`
}

// TitleStats 称号统计
type TitleStats struct {
	TitleID     string `json:"title_id"`
	TitleName   string `json:"title_name"`
	UnlockCount int64  `json:"unlock_count"`
	EquipCount  int64  `json:"equip_count"`
}

// AchievementStats 成就统计
type AchievementStats struct {
	AchievementID   string `json:"achievement_id"`
	AchievementName string `json:"achievement_name"`
	UnlockCount     int64  `json:"unlock_count"`
}

// ReputationStats 声望统计
type ReputationStats struct {
	Faction           string  `json:"faction"`
	AverageReputation float64 `json:"average_reputation"`
	MaxReputation     int     `json:"max_reputation"`
	PlayerCount       int64   `json:"player_count"`
}

// HonorQuery 荣誉查询条件
type HonorQuery struct {
	PlayerIDs      []string              `json:"player_ids,omitempty"`
	MinHonorLevel  int                   `json:"min_honor_level,omitempty"`
	MaxHonorLevel  int                   `json:"max_honor_level,omitempty"`
	MinHonorPoints int                   `json:"min_honor_points,omitempty"`
	MaxHonorPoints int                   `json:"max_honor_points,omitempty"`
	TitleEquipped  string                `json:"title_equipped,omitempty"`
	Achievements   []string              `json:"achievements,omitempty"`
	Reputation     map[string]int        `json:"reputation,omitempty"`
	Statistics     map[StatisticType]int `json:"statistics,omitempty"`
	Limit          int                   `json:"limit,omitempty"`
	Offset         int                   `json:"offset,omitempty"`
	SortBy         string                `json:"sort_by,omitempty"`
	SortOrder      string                `json:"sort_order,omitempty"`
}

// HonorQueryRepository 荣誉查询仓储接口
type HonorQueryRepository interface {
	// FindByQuery 根据查询条件查找荣誉信息
	FindByQuery(ctx context.Context, query *HonorQuery) ([]*HonorAggregate, error)

	// CountByQuery 根据查询条件统计数量
	CountByQuery(ctx context.Context, query *HonorQuery) (int64, error)

	// FindPlayersWithTitle 查找拥有指定称号的玩家
	FindPlayersWithTitle(ctx context.Context, titleID string, unlocked bool, equipped bool) ([]*HonorAggregate, error)

	// FindPlayersWithAchievement 查找拥有指定成就的玩家
	FindPlayersWithAchievement(ctx context.Context, achievementID string, unlocked bool) ([]*HonorAggregate, error)

	// FindTopPlayersByStatistic 根据统计数据查找排行榜
	FindTopPlayersByStatistic(ctx context.Context, statType StatisticType, limit int) ([]*HonorAggregate, error)

	// FindPlayersByReputationRange 根据声望范围查找玩家
	FindPlayersByReputationRange(ctx context.Context, faction string, minRep, maxRep int) ([]*HonorAggregate, error)
}
