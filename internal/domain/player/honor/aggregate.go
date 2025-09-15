package honor

import (
	"time"
	"github.com/google/uuid"
)

// HonorAggregate 荣誉聚合根
type HonorAggregate struct {
	playerID       string
	titles         map[string]*Title
	achievements   map[string]*Achievement
	currentTitle   string
	honorPoints    int
	honorLevel     int
	reputation     map[string]int // 声望系统
	statistics     *PlayerStatistics
	updatedAt      time.Time
	version        int
}

// NewHonorAggregate 创建荣誉聚合根
func NewHonorAggregate(playerID string) *HonorAggregate {
	return &HonorAggregate{
		playerID:     playerID,
		titles:       make(map[string]*Title),
		achievements: make(map[string]*Achievement),
		currentTitle: "",
		honorPoints:  0,
		honorLevel:   1,
		reputation:   make(map[string]int),
		statistics:   NewPlayerStatistics(),
		updatedAt:    time.Now(),
		version:      1,
	}
}

// GetPlayerID 获取玩家ID
func (h *HonorAggregate) GetPlayerID() string {
	return h.playerID
}

// AddTitle 添加称号
func (h *HonorAggregate) AddTitle(title *Title) error {
	if title == nil {
		return ErrInvalidTitle
	}
	
	// 检查是否已拥有该称号
	if _, exists := h.titles[title.GetID()]; exists {
		return ErrTitleAlreadyOwned
	}
	
	h.titles[title.GetID()] = title
	h.updateVersion()
	return nil
}

// UnlockTitle 解锁称号
func (h *HonorAggregate) UnlockTitle(titleID string) error {
	title, exists := h.titles[titleID]
	if !exists {
		return ErrTitleNotFound
	}
	
	if title.IsUnlocked() {
		return ErrTitleAlreadyUnlocked
	}
	
	// 检查解锁条件
	if !h.checkTitleUnlockConditions(title) {
		return ErrTitleConditionNotMet
	}
	
	title.Unlock()
	h.updateVersion()
	return nil
}

// EquipTitle 装备称号
func (h *HonorAggregate) EquipTitle(titleID string) error {
	title, exists := h.titles[titleID]
	if !exists {
		return ErrTitleNotFound
	}
	
	if !title.IsUnlocked() {
		return ErrTitleNotUnlocked
	}
	
	// 卸下当前称号
	if h.currentTitle != "" {
		if currentTitle, exists := h.titles[h.currentTitle]; exists {
			currentTitle.Unequip()
		}
	}
	
	// 装备新称号
	title.Equip()
	h.currentTitle = titleID
	h.updateVersion()
	return nil
}

// UnequipTitle 卸下称号
func (h *HonorAggregate) UnequipTitle() error {
	if h.currentTitle == "" {
		return ErrNoTitleEquipped
	}
	
	if title, exists := h.titles[h.currentTitle]; exists {
		title.Unequip()
	}
	
	h.currentTitle = ""
	h.updateVersion()
	return nil
}

// GetCurrentTitle 获取当前称号
func (h *HonorAggregate) GetCurrentTitle() *Title {
	if h.currentTitle == "" {
		return nil
	}
	return h.titles[h.currentTitle]
}

// GetAllTitles 获取所有称号
func (h *HonorAggregate) GetAllTitles() map[string]*Title {
	return h.titles
}

// GetUnlockedTitles 获取已解锁的称号
func (h *HonorAggregate) GetUnlockedTitles() []*Title {
	var unlocked []*Title
	for _, title := range h.titles {
		if title.IsUnlocked() {
			unlocked = append(unlocked, title)
		}
	}
	return unlocked
}

// AddAchievement 添加成就
func (h *HonorAggregate) AddAchievement(achievement *Achievement) error {
	if achievement == nil {
		return ErrInvalidAchievement
	}
	
	h.achievements[achievement.GetID()] = achievement
	h.updateVersion()
	return nil
}

// UnlockAchievement 解锁成就
func (h *HonorAggregate) UnlockAchievement(achievementID string) error {
	achievement, exists := h.achievements[achievementID]
	if !exists {
		return ErrAchievementNotFound
	}
	
	if achievement.IsUnlocked() {
		return ErrAchievementAlreadyUnlocked
	}
	
	// 检查解锁条件
	if !h.checkAchievementUnlockConditions(achievement) {
		return ErrAchievementConditionNotMet
	}
	
	achievement.Unlock()
	
	// 给予荣誉点数奖励
	h.AddHonorPoints(achievement.GetHonorReward())
	
	h.updateVersion()
	return nil
}

// GetAchievement 获取成就
func (h *HonorAggregate) GetAchievement(achievementID string) *Achievement {
	return h.achievements[achievementID]
}

// GetAllAchievements 获取所有成就
func (h *HonorAggregate) GetAllAchievements() map[string]*Achievement {
	return h.achievements
}

// GetUnlockedAchievements 获取已解锁的成就
func (h *HonorAggregate) GetUnlockedAchievements() []*Achievement {
	var unlocked []*Achievement
	for _, achievement := range h.achievements {
		if achievement.IsUnlocked() {
			unlocked = append(unlocked, achievement)
		}
	}
	return unlocked
}

// AddHonorPoints 增加荣誉点数
func (h *HonorAggregate) AddHonorPoints(points int) {
	h.honorPoints += points
	
	// 检查是否升级
	h.checkHonorLevelUp()
	
	h.updateVersion()
}

// GetHonorPoints 获取荣誉点数
func (h *HonorAggregate) GetHonorPoints() int {
	return h.honorPoints
}

// GetHonorLevel 获取荣誉等级
func (h *HonorAggregate) GetHonorLevel() int {
	return h.honorLevel
}

// AddReputation 增加声望
func (h *HonorAggregate) AddReputation(faction string, points int) {
	h.reputation[faction] += points
	h.updateVersion()
}

// GetReputation 获取声望
func (h *HonorAggregate) GetReputation(faction string) int {
	return h.reputation[faction]
}

// GetAllReputation 获取所有声望
func (h *HonorAggregate) GetAllReputation() map[string]int {
	return h.reputation
}

// UpdateStatistics 更新统计数据
func (h *HonorAggregate) UpdateStatistics(statType StatisticType, value int) {
	h.statistics.UpdateStatistic(statType, value)
	
	// 检查是否触发成就或称号解锁
	h.checkUnlockConditions()
	
	h.updateVersion()
}

// GetStatistics 获取统计数据
func (h *HonorAggregate) GetStatistics() *PlayerStatistics {
	return h.statistics
}

// GetTitlesByCategory 根据分类获取称号
func (h *HonorAggregate) GetTitlesByCategory(category TitleCategory) []*Title {
	var titles []*Title
	for _, title := range h.titles {
		if title.GetCategory() == category {
			titles = append(titles, title)
		}
	}
	return titles
}

// GetAchievementsByCategory 根据分类获取成就
func (h *HonorAggregate) GetAchievementsByCategory(category AchievementCategory) []*Achievement {
	var achievements []*Achievement
	for _, achievement := range h.achievements {
		if achievement.GetCategory() == category {
			achievements = append(achievements, achievement)
		}
	}
	return achievements
}

// GetVersion 获取版本
func (h *HonorAggregate) GetVersion() int {
	return h.version
}

// GetUpdatedAt 获取更新时间
func (h *HonorAggregate) GetUpdatedAt() time.Time {
	return h.updatedAt
}

// 私有方法

// checkTitleUnlockConditions 检查称号解锁条件
func (h *HonorAggregate) checkTitleUnlockConditions(title *Title) bool {
	for _, condition := range title.GetUnlockConditions() {
		if !h.checkCondition(condition) {
			return false
		}
	}
	return true
}

// checkAchievementUnlockConditions 检查成就解锁条件
func (h *HonorAggregate) checkAchievementUnlockConditions(achievement *Achievement) bool {
	for _, condition := range achievement.GetUnlockConditions() {
		if !h.checkCondition(condition) {
			return false
		}
	}
	return true
}

// checkCondition 检查单个条件
func (h *HonorAggregate) checkCondition(condition *UnlockCondition) bool {
	switch condition.GetConditionType() {
	case ConditionTypeLevel:
		// 这里需要从外部获取玩家等级，暂时返回true
		return true
	case ConditionTypeStatistic:
		statValue := h.statistics.GetStatistic(condition.GetStatisticType())
		return statValue >= condition.GetRequiredValue()
	case ConditionTypeReputation:
		repValue := h.GetReputation(condition.GetFaction())
		return repValue >= condition.GetRequiredValue()
	case ConditionTypeAchievement:
		achievement := h.GetAchievement(condition.GetRequiredAchievement())
		return achievement != nil && achievement.IsUnlocked()
	case ConditionTypeTitle:
		title := h.titles[condition.GetRequiredTitle()]
		return title != nil && title.IsUnlocked()
	default:
		return false
	}
}

// checkHonorLevelUp 检查荣誉等级提升
func (h *HonorAggregate) checkHonorLevelUp() {
	requiredPoints := h.getRequiredPointsForLevel(h.honorLevel + 1)
	if h.honorPoints >= requiredPoints {
		h.honorLevel++
		// 可以在这里触发等级提升事件
	}
}

// getRequiredPointsForLevel 获取等级所需点数
func (h *HonorAggregate) getRequiredPointsForLevel(level int) int {
	// 简单的等级计算公式
	return level * level * 100
}

// checkUnlockConditions 检查解锁条件
func (h *HonorAggregate) checkUnlockConditions() {
	// 检查所有未解锁的称号
	for _, title := range h.titles {
		if !title.IsUnlocked() && h.checkTitleUnlockConditions(title) {
			title.Unlock()
		}
	}
	
	// 检查所有未解锁的成就
	for _, achievement := range h.achievements {
		if !achievement.IsUnlocked() && h.checkAchievementUnlockConditions(achievement) {
			achievement.Unlock()
			h.AddHonorPoints(achievement.GetHonorReward())
		}
	}
}

// updateVersion 更新版本
func (h *HonorAggregate) updateVersion() {
	h.version++
	h.updatedAt = time.Now()
}