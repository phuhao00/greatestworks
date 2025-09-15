package honor

import (
	"time"
	"github.com/google/uuid"
)

// Title 称号实体
type Title struct {
	id               string
	name             string
	description      string
	category         TitleCategory
	rarity           TitleRarity
	unlockConditions []*UnlockCondition
	attributeBonus   map[string]int
	isUnlocked       bool
	isEquipped       bool
	unlockedAt       *time.Time
	equippedAt       *time.Time
	createdAt        time.Time
}

// NewTitle 创建新称号
func NewTitle(id, name, description string, category TitleCategory, rarity TitleRarity) *Title {
	return &Title{
		id:               id,
		name:             name,
		description:      description,
		category:         category,
		rarity:           rarity,
		unlockConditions: make([]*UnlockCondition, 0),
		attributeBonus:   make(map[string]int),
		isUnlocked:       false,
		isEquipped:       false,
		createdAt:        time.Now(),
	}
}

// GetID 获取称号ID
func (t *Title) GetID() string {
	return t.id
}

// GetName 获取称号名称
func (t *Title) GetName() string {
	return t.name
}

// GetDescription 获取称号描述
func (t *Title) GetDescription() string {
	return t.description
}

// GetCategory 获取称号分类
func (t *Title) GetCategory() TitleCategory {
	return t.category
}

// GetRarity 获取称号稀有度
func (t *Title) GetRarity() TitleRarity {
	return t.rarity
}

// GetUnlockConditions 获取解锁条件
func (t *Title) GetUnlockConditions() []*UnlockCondition {
	return t.unlockConditions
}

// AddUnlockCondition 添加解锁条件
func (t *Title) AddUnlockCondition(condition *UnlockCondition) {
	t.unlockConditions = append(t.unlockConditions, condition)
}

// GetAttributeBonus 获取属性加成
func (t *Title) GetAttributeBonus() map[string]int {
	return t.attributeBonus
}

// SetAttributeBonus 设置属性加成
func (t *Title) SetAttributeBonus(attribute string, bonus int) {
	t.attributeBonus[attribute] = bonus
}

// IsUnlocked 是否已解锁
func (t *Title) IsUnlocked() bool {
	return t.isUnlocked
}

// IsEquipped 是否已装备
func (t *Title) IsEquipped() bool {
	return t.isEquipped
}

// Unlock 解锁称号
func (t *Title) Unlock() {
	if !t.isUnlocked {
		t.isUnlocked = true
		now := time.Now()
		t.unlockedAt = &now
	}
}

// Equip 装备称号
func (t *Title) Equip() {
	if t.isUnlocked && !t.isEquipped {
		t.isEquipped = true
		now := time.Now()
		t.equippedAt = &now
	}
}

// Unequip 卸下称号
func (t *Title) Unequip() {
	if t.isEquipped {
		t.isEquipped = false
		t.equippedAt = nil
	}
}

// GetUnlockedAt 获取解锁时间
func (t *Title) GetUnlockedAt() *time.Time {
	return t.unlockedAt
}

// GetEquippedAt 获取装备时间
func (t *Title) GetEquippedAt() *time.Time {
	return t.equippedAt
}

// Achievement 成就实体
type Achievement struct {
	id               string
	name             string
	description      string
	category         AchievementCategory
	type_            AchievementType
	unlockConditions []*UnlockCondition
	honorReward      int
	itemRewards      []string
	isUnlocked       bool
	unlockedAt       *time.Time
	createdAt        time.Time
}

// NewAchievement 创建新成就
func NewAchievement(id, name, description string, category AchievementCategory, achievementType AchievementType) *Achievement {
	return &Achievement{
		id:               id,
		name:             name,
		description:      description,
		category:         category,
		type_:            achievementType,
		unlockConditions: make([]*UnlockCondition, 0),
		itemRewards:      make([]string, 0),
		isUnlocked:       false,
		createdAt:        time.Now(),
	}
}

// GetID 获取成就ID
func (a *Achievement) GetID() string {
	return a.id
}

// GetName 获取成就名称
func (a *Achievement) GetName() string {
	return a.name
}

// GetDescription 获取成就描述
func (a *Achievement) GetDescription() string {
	return a.description
}

// GetCategory 获取成就分类
func (a *Achievement) GetCategory() AchievementCategory {
	return a.category
}

// GetType 获取成就类型
func (a *Achievement) GetType() AchievementType {
	return a.type_
}

// GetUnlockConditions 获取解锁条件
func (a *Achievement) GetUnlockConditions() []*UnlockCondition {
	return a.unlockConditions
}

// AddUnlockCondition 添加解锁条件
func (a *Achievement) AddUnlockCondition(condition *UnlockCondition) {
	a.unlockConditions = append(a.unlockConditions, condition)
}

// GetHonorReward 获取荣誉点数奖励
func (a *Achievement) GetHonorReward() int {
	return a.honorReward
}

// SetHonorReward 设置荣誉点数奖励
func (a *Achievement) SetHonorReward(reward int) {
	a.honorReward = reward
}

// GetItemRewards 获取物品奖励
func (a *Achievement) GetItemRewards() []string {
	return a.itemRewards
}

// AddItemReward 添加物品奖励
func (a *Achievement) AddItemReward(itemID string) {
	a.itemRewards = append(a.itemRewards, itemID)
}

// IsUnlocked 是否已解锁
func (a *Achievement) IsUnlocked() bool {
	return a.isUnlocked
}

// Unlock 解锁成就
func (a *Achievement) Unlock() {
	if !a.isUnlocked {
		a.isUnlocked = true
		now := time.Now()
		a.unlockedAt = &now
	}
}

// GetUnlockedAt 获取解锁时间
func (a *Achievement) GetUnlockedAt() *time.Time {
	return a.unlockedAt
}

// PlayerStatistics 玩家统计数据实体
type PlayerStatistics struct {
	statistics map[StatisticType]int
	updatedAt  time.Time
}

// NewPlayerStatistics 创建新的玩家统计数据
func NewPlayerStatistics() *PlayerStatistics {
	return &PlayerStatistics{
		statistics: make(map[StatisticType]int),
		updatedAt:  time.Now(),
	}
}

// UpdateStatistic 更新统计数据
func (ps *PlayerStatistics) UpdateStatistic(statType StatisticType, value int) {
	ps.statistics[statType] = value
	ps.updatedAt = time.Now()
}

// IncrementStatistic 增加统计数据
func (ps *PlayerStatistics) IncrementStatistic(statType StatisticType, increment int) {
	ps.statistics[statType] += increment
	ps.updatedAt = time.Now()
}

// GetStatistic 获取统计数据
func (ps *PlayerStatistics) GetStatistic(statType StatisticType) int {
	return ps.statistics[statType]
}

// GetAllStatistics 获取所有统计数据
func (ps *PlayerStatistics) GetAllStatistics() map[StatisticType]int {
	return ps.statistics
}

// GetUpdatedAt 获取更新时间
func (ps *PlayerStatistics) GetUpdatedAt() time.Time {
	return ps.updatedAt
}