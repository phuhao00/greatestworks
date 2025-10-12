package honor

import (
	"fmt"
	"time"
)

// TitleCategory 称号分类
type TitleCategory int

const (
	TitleCategoryUnknown     TitleCategory = iota
	TitleCategoryCombat                    // 战斗类
	TitleCategoryExploration               // 探索类
	TitleCategorySocial                    // 社交类
	TitleCategoryLifestyle                 // 生活类
	TitleCategorySpecial                   // 特殊类
	TitleCategoryEvent                     // 活动类
)

// String 返回称号分类的字符串表示
func (tc TitleCategory) String() string {
	switch tc {
	case TitleCategoryCombat:
		return "combat"
	case TitleCategoryExploration:
		return "exploration"
	case TitleCategorySocial:
		return "social"
	case TitleCategoryLifestyle:
		return "lifestyle"
	case TitleCategorySpecial:
		return "special"
	case TitleCategoryEvent:
		return "event"
	default:
		return "unknown"
	}
}

// TitleRarity 称号稀有度
type TitleRarity int

const (
	TitleRarityCommon TitleRarity = iota
	TitleRarityUncommon
	TitleRarityRare
	TitleRarityEpic
	TitleRarityLegendary
	TitleRarityMythic
)

// String 返回称号稀有度的字符串表示
func (tr TitleRarity) String() string {
	switch tr {
	case TitleRarityCommon:
		return "common"
	case TitleRarityUncommon:
		return "uncommon"
	case TitleRarityRare:
		return "rare"
	case TitleRarityEpic:
		return "epic"
	case TitleRarityLegendary:
		return "legendary"
	case TitleRarityMythic:
		return "mythic"
	default:
		return "common"
	}
}

// AchievementCategory 成就分类
type AchievementCategory int

const (
	AchievementCategoryUnknown     AchievementCategory = iota
	AchievementCategoryCombat                          // 战斗成就
	AchievementCategoryExploration                     // 探索成就
	AchievementCategorySocial                          // 社交成就
	AchievementCategoryCollection                      // 收集成就
	AchievementCategoryProgression                     // 进度成就
	AchievementCategorySpecial                         // 特殊成就
)

// String 返回成就分类的字符串表示
func (ac AchievementCategory) String() string {
	switch ac {
	case AchievementCategoryCombat:
		return "combat"
	case AchievementCategoryExploration:
		return "exploration"
	case AchievementCategorySocial:
		return "social"
	case AchievementCategoryCollection:
		return "collection"
	case AchievementCategoryProgression:
		return "progression"
	case AchievementCategorySpecial:
		return "special"
	default:
		return "unknown"
	}
}

// AchievementType 成就类型
type AchievementType int

const (
	AchievementTypeNormal  AchievementType = iota
	AchievementTypeHidden                  // 隐藏成就
	AchievementTypeDaily                   // 日常成就
	AchievementTypeWeekly                  // 周常成就
	AchievementTypeMonthly                 // 月常成就
	AchievementTypeEvent                   // 活动成就
)

// String 返回成就类型的字符串表示
func (at AchievementType) String() string {
	switch at {
	case AchievementTypeNormal:
		return "normal"
	case AchievementTypeHidden:
		return "hidden"
	case AchievementTypeDaily:
		return "daily"
	case AchievementTypeWeekly:
		return "weekly"
	case AchievementTypeMonthly:
		return "monthly"
	case AchievementTypeEvent:
		return "event"
	default:
		return "normal"
	}
}

// StatisticType 统计数据类型
type StatisticType int

const (
	StatisticTypeUnknown           StatisticType = iota
	StatisticTypeKillCount                       // 击杀数量
	StatisticTypeDeathCount                      // 死亡数量
	StatisticTypeDamageDealt                     // 造成伤害
	StatisticTypeDamageTaken                     // 承受伤害
	StatisticTypeHealingDone                     // 治疗量
	StatisticTypeDistanceTraveled                // 旅行距离
	StatisticTypeQuestsCompleted                 // 完成任务数
	StatisticTypeItemsCrafted                    // 制作物品数
	StatisticTypeItemsCollected                  // 收集物品数
	StatisticTypeGoldEarned                      // 获得金币
	StatisticTypeGoldSpent                       // 花费金币
	StatisticTypePlayTime                        // 游戏时间
	StatisticTypeLoginDays                       // 登录天数
	StatisticTypeFriendsCount                    // 好友数量
	StatisticTypeGuildContribution               // 公会贡献
)

// String 返回统计数据类型的字符串表示
func (st StatisticType) String() string {
	switch st {
	case StatisticTypeKillCount:
		return "kill_count"
	case StatisticTypeDeathCount:
		return "death_count"
	case StatisticTypeDamageDealt:
		return "damage_dealt"
	case StatisticTypeDamageTaken:
		return "damage_taken"
	case StatisticTypeHealingDone:
		return "healing_done"
	case StatisticTypeDistanceTraveled:
		return "distance_traveled"
	case StatisticTypeQuestsCompleted:
		return "quests_completed"
	case StatisticTypeItemsCrafted:
		return "items_crafted"
	case StatisticTypeItemsCollected:
		return "items_collected"
	case StatisticTypeGoldEarned:
		return "gold_earned"
	case StatisticTypeGoldSpent:
		return "gold_spent"
	case StatisticTypePlayTime:
		return "play_time"
	case StatisticTypeLoginDays:
		return "login_days"
	case StatisticTypeFriendsCount:
		return "friends_count"
	case StatisticTypeGuildContribution:
		return "guild_contribution"
	default:
		return "unknown"
	}
}

// ConditionType 条件类型
type ConditionType int

const (
	ConditionTypeUnknown     ConditionType = iota
	ConditionTypeLevel                     // 等级条件
	ConditionTypeStatistic                 // 统计数据条件
	ConditionTypeReputation                // 声望条件
	ConditionTypeAchievement               // 成就条件
	ConditionTypeTitle                     // 称号条件
	ConditionTypeItem                      // 物品条件
	ConditionTypeQuest                     // 任务条件
	ConditionTypeTime                      // 时间条件
)

// String 返回条件类型的字符串表示
func (ct ConditionType) String() string {
	switch ct {
	case ConditionTypeLevel:
		return "level"
	case ConditionTypeStatistic:
		return "statistic"
	case ConditionTypeReputation:
		return "reputation"
	case ConditionTypeAchievement:
		return "achievement"
	case ConditionTypeTitle:
		return "title"
	case ConditionTypeItem:
		return "item"
	case ConditionTypeQuest:
		return "quest"
	case ConditionTypeTime:
		return "time"
	default:
		return "unknown"
	}
}

// UnlockCondition 解锁条件值对象
type UnlockCondition struct {
	conditionType       ConditionType
	requiredValue       int
	statisticType       StatisticType
	faction             string
	requiredAchievement string
	requiredTitle       string
	requiredItem        string
	requiredQuest       string
	timeRequirement     time.Duration
	description         string
}

// NewUnlockCondition 创建解锁条件
func NewUnlockCondition(conditionType ConditionType, description string) *UnlockCondition {
	return &UnlockCondition{
		conditionType: conditionType,
		description:   description,
	}
}

// NewLevelCondition 创建等级条件
func NewLevelCondition(requiredLevel int) *UnlockCondition {
	return &UnlockCondition{
		conditionType: ConditionTypeLevel,
		requiredValue: requiredLevel,
		description:   fmt.Sprintf("达到等级 %d", requiredLevel),
	}
}

// NewStatisticCondition 创建统计数据条件
func NewStatisticCondition(statType StatisticType, requiredValue int) *UnlockCondition {
	return &UnlockCondition{
		conditionType: ConditionTypeStatistic,
		statisticType: statType,
		requiredValue: requiredValue,
		description:   fmt.Sprintf("%s 达到 %d", statType.String(), requiredValue),
	}
}

// NewReputationCondition 创建声望条件
func NewReputationCondition(faction string, requiredValue int) *UnlockCondition {
	return &UnlockCondition{
		conditionType: ConditionTypeReputation,
		faction:       faction,
		requiredValue: requiredValue,
		description:   fmt.Sprintf("%s 声望达到 %d", faction, requiredValue),
	}
}

// NewAchievementCondition 创建成就条件
func NewAchievementCondition(achievementID string) *UnlockCondition {
	return &UnlockCondition{
		conditionType:       ConditionTypeAchievement,
		requiredAchievement: achievementID,
		description:         fmt.Sprintf("完成成就: %s", achievementID),
	}
}

// NewTitleCondition 创建称号条件
func NewTitleCondition(titleID string) *UnlockCondition {
	return &UnlockCondition{
		conditionType: ConditionTypeTitle,
		requiredTitle: titleID,
		description:   fmt.Sprintf("获得称号: %s", titleID),
	}
}

// GetConditionType 获取条件类型
func (uc *UnlockCondition) GetConditionType() ConditionType {
	return uc.conditionType
}

// GetRequiredValue 获取所需值
func (uc *UnlockCondition) GetRequiredValue() int {
	return uc.requiredValue
}

// GetStatisticType 获取统计数据类型
func (uc *UnlockCondition) GetStatisticType() StatisticType {
	return uc.statisticType
}

// GetFaction 获取阵营
func (uc *UnlockCondition) GetFaction() string {
	return uc.faction
}

// GetRequiredAchievement 获取所需成就
func (uc *UnlockCondition) GetRequiredAchievement() string {
	return uc.requiredAchievement
}

// GetRequiredTitle 获取所需称号
func (uc *UnlockCondition) GetRequiredTitle() string {
	return uc.requiredTitle
}

// GetRequiredItem 获取所需物品
func (uc *UnlockCondition) GetRequiredItem() string {
	return uc.requiredItem
}

// GetRequiredQuest 获取所需任务
func (uc *UnlockCondition) GetRequiredQuest() string {
	return uc.requiredQuest
}

// GetTimeRequirement 获取时间要求
func (uc *UnlockCondition) GetTimeRequirement() time.Duration {
	return uc.timeRequirement
}

// GetDescription 获取描述
func (uc *UnlockCondition) GetDescription() string {
	return uc.description
}

// HonorLevel 荣誉等级值对象
type HonorLevel struct {
	level       int
	requiredXP  int
	title       string
	description string
	rewards     []string
}

// NewHonorLevel 创建荣誉等级
func NewHonorLevel(level, requiredXP int, title, description string) *HonorLevel {
	return &HonorLevel{
		level:       level,
		requiredXP:  requiredXP,
		title:       title,
		description: description,
		rewards:     make([]string, 0),
	}
}

// GetLevel 获取等级
func (hl *HonorLevel) GetLevel() int {
	return hl.level
}

// GetRequiredXP 获取所需经验
func (hl *HonorLevel) GetRequiredXP() int {
	return hl.requiredXP
}

// GetTitle 获取等级称号
func (hl *HonorLevel) GetTitle() string {
	return hl.title
}

// GetDescription 获取描述
func (hl *HonorLevel) GetDescription() string {
	return hl.description
}

// GetRewards 获取奖励
func (hl *HonorLevel) GetRewards() []string {
	return hl.rewards
}

// AddReward 添加奖励
func (hl *HonorLevel) AddReward(reward string) {
	hl.rewards = append(hl.rewards, reward)
}
