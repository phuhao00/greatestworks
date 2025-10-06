package honor

import (
	"fmt"
	// "time" // 未使用
)

// HonorService 荣誉领域服务
type HonorService struct {
	titleTemplates       map[string]*TitleTemplate
	achievementTemplates map[string]*AchievementTemplate
	honorLevels          map[int]*HonorLevel
}

// NewHonorService 创建荣誉服务
func NewHonorService() *HonorService {
	return &HonorService{
		titleTemplates:       make(map[string]*TitleTemplate),
		achievementTemplates: make(map[string]*AchievementTemplate),
		honorLevels:          make(map[int]*HonorLevel),
	}
}

// TitleTemplate 称号模板
type TitleTemplate struct {
	id               string
	name             string
	description      string
	category         TitleCategory
	rarity           TitleRarity
	unlockConditions []*UnlockCondition
	attributeBonus   map[string]int
}

// NewTitleTemplate 创建称号模板
func NewTitleTemplate(id, name, description string, category TitleCategory, rarity TitleRarity) *TitleTemplate {
	return &TitleTemplate{
		id:               id,
		name:             name,
		description:      description,
		category:         category,
		rarity:           rarity,
		unlockConditions: make([]*UnlockCondition, 0),
		attributeBonus:   make(map[string]int),
	}
}

// CreateTitle 根据模板创建称号
func (tt *TitleTemplate) CreateTitle() *Title {
	title := NewTitle(tt.id, tt.name, tt.description, tt.category, tt.rarity)

	// 复制解锁条件
	for _, condition := range tt.unlockConditions {
		title.AddUnlockCondition(condition)
	}

	// 复制属性加成
	for attr, bonus := range tt.attributeBonus {
		title.SetAttributeBonus(attr, bonus)
	}

	return title
}

// AchievementTemplate 成就模板
type AchievementTemplate struct {
	id               string
	name             string
	description      string
	category         AchievementCategory
	type_            AchievementType
	unlockConditions []*UnlockCondition
	honorReward      int
	itemRewards      []string
}

// NewAchievementTemplate 创建成就模板
func NewAchievementTemplate(id, name, description string, category AchievementCategory, achievementType AchievementType) *AchievementTemplate {
	return &AchievementTemplate{
		id:               id,
		name:             name,
		description:      description,
		category:         category,
		type_:            achievementType,
		unlockConditions: make([]*UnlockCondition, 0),
		itemRewards:      make([]string, 0),
	}
}

// CreateAchievement 根据模板创建成就
func (at *AchievementTemplate) CreateAchievement() *Achievement {
	achievement := NewAchievement(at.id, at.name, at.description, at.category, at.type_)

	// 复制解锁条件
	for _, condition := range at.unlockConditions {
		achievement.AddUnlockCondition(condition)
	}

	// 设置奖励
	achievement.SetHonorReward(at.honorReward)
	for _, itemID := range at.itemRewards {
		achievement.AddItemReward(itemID)
	}

	return achievement
}

// RegisterTitleTemplate 注册称号模板
func (hs *HonorService) RegisterTitleTemplate(template *TitleTemplate) {
	hs.titleTemplates[template.id] = template
}

// RegisterAchievementTemplate 注册成就模板
func (hs *HonorService) RegisterAchievementTemplate(template *AchievementTemplate) {
	hs.achievementTemplates[template.id] = template
}

// RegisterHonorLevel 注册荣誉等级
func (hs *HonorService) RegisterHonorLevel(level *HonorLevel) {
	hs.honorLevels[level.GetLevel()] = level
}

// CreatePlayerHonor 为玩家创建荣誉系统
func (hs *HonorService) CreatePlayerHonor(playerID string) *HonorAggregate {
	honor := NewHonorAggregate(playerID)

	// 添加所有称号模板
	for _, template := range hs.titleTemplates {
		title := template.CreateTitle()
		honor.AddTitle(title)
	}

	// 添加所有成就模板
	for _, template := range hs.achievementTemplates {
		achievement := template.CreateAchievement()
		honor.AddAchievement(achievement)
	}

	return honor
}

// CalculateHonorLevel 计算荣誉等级
func (hs *HonorService) CalculateHonorLevel(honorPoints int) int {
	level := 1
	for lvl, honorLevel := range hs.honorLevels {
		if honorPoints >= honorLevel.GetRequiredXP() && lvl > level {
			level = lvl
		}
	}
	return level
}

// GetHonorLevel 获取荣誉等级信息
func (hs *HonorService) GetHonorLevel(level int) *HonorLevel {
	return hs.honorLevels[level]
}

// GetNextHonorLevel 获取下一个荣誉等级
func (hs *HonorService) GetNextHonorLevel(currentLevel int) *HonorLevel {
	return hs.honorLevels[currentLevel+1]
}

// ValidateTitleUnlock 验证称号解锁
func (hs *HonorService) ValidateTitleUnlock(honor *HonorAggregate, titleID string) error {
	title := honor.titles[titleID]
	if title == nil {
		return ErrTitleNotFound
	}

	if title.IsUnlocked() {
		return ErrTitleAlreadyUnlocked
	}

	// 检查所有解锁条件
	for _, condition := range title.GetUnlockConditions() {
		if !hs.checkUnlockCondition(honor, condition) {
			return fmt.Errorf("条件未满足: %s", condition.GetDescription())
		}
	}

	return nil
}

// ValidateAchievementUnlock 验证成就解锁
func (hs *HonorService) ValidateAchievementUnlock(honor *HonorAggregate, achievementID string) error {
	achievement := honor.achievements[achievementID]
	if achievement == nil {
		return ErrAchievementNotFound
	}

	if achievement.IsUnlocked() {
		return ErrAchievementAlreadyUnlocked
	}

	// 检查所有解锁条件
	for _, condition := range achievement.GetUnlockConditions() {
		if !hs.checkUnlockCondition(honor, condition) {
			return fmt.Errorf("条件未满足: %s", condition.GetDescription())
		}
	}

	return nil
}

// checkUnlockCondition 检查解锁条件
func (hs *HonorService) checkUnlockCondition(honor *HonorAggregate, condition *UnlockCondition) bool {
	switch condition.GetConditionType() {
	case ConditionTypeLevel:
		// 这里需要从外部获取玩家等级，暂时返回true
		return true
	case ConditionTypeStatistic:
		statValue := honor.statistics.GetStatistic(condition.GetStatisticType())
		return statValue >= condition.GetRequiredValue()
	case ConditionTypeReputation:
		repValue := honor.GetReputation(condition.GetFaction())
		return repValue >= condition.GetRequiredValue()
	case ConditionTypeAchievement:
		achievement := honor.GetAchievement(condition.GetRequiredAchievement())
		return achievement != nil && achievement.IsUnlocked()
	case ConditionTypeTitle:
		title := honor.titles[condition.GetRequiredTitle()]
		return title != nil && title.IsUnlocked()
	default:
		return false
	}
}

// CalculateTitleAttributeBonus 计算称号属性加成
func (hs *HonorService) CalculateTitleAttributeBonus(honor *HonorAggregate) map[string]int {
	bonuses := make(map[string]int)

	// 只计算当前装备的称号
	currentTitle := honor.GetCurrentTitle()
	if currentTitle != nil {
		for attr, bonus := range currentTitle.GetAttributeBonus() {
			bonuses[attr] += bonus
		}
	}

	return bonuses
}

// GetTitlesByRarity 根据稀有度获取称号
func (hs *HonorService) GetTitlesByRarity(honor *HonorAggregate, rarity TitleRarity) []*Title {
	var titles []*Title
	for _, title := range honor.GetAllTitles() {
		if title.GetRarity() == rarity {
			titles = append(titles, title)
		}
	}
	return titles
}

// GetAchievementsByType 根据类型获取成就
func (hs *HonorService) GetAchievementsByType(honor *HonorAggregate, achievementType AchievementType) []*Achievement {
	var achievements []*Achievement
	for _, achievement := range honor.GetAllAchievements() {
		if achievement.GetType() == achievementType {
			achievements = append(achievements, achievement)
		}
	}
	return achievements
}

// CalculateHonorRank 计算荣誉排名（需要外部数据支持）
func (hs *HonorService) CalculateHonorRank(honor *HonorAggregate, allPlayers []*HonorAggregate) int {
	rank := 1
	currentPoints := honor.GetHonorPoints()

	for _, otherHonor := range allPlayers {
		if otherHonor.GetPlayerID() != honor.GetPlayerID() && otherHonor.GetHonorPoints() > currentPoints {
			rank++
		}
	}

	return rank
}

// GetRecommendedTitles 获取推荐称号（接近解锁的称号）
func (hs *HonorService) GetRecommendedTitles(honor *HonorAggregate) []*Title {
	var recommended []*Title

	for _, title := range honor.GetAllTitles() {
		if !title.IsUnlocked() {
			// 检查是否接近解锁
			meetsConditions := 0
			totalConditions := len(title.GetUnlockConditions())

			for _, condition := range title.GetUnlockConditions() {
				if hs.checkUnlockCondition(honor, condition) {
					meetsConditions++
				}
			}

			// 如果满足80%以上的条件，则推荐
			if totalConditions > 0 && float64(meetsConditions)/float64(totalConditions) >= 0.8 {
				recommended = append(recommended, title)
			}
		}
	}

	return recommended
}

// InitializeDefaultTemplates 初始化默认模板
func (hs *HonorService) InitializeDefaultTemplates() {
	// 初始化默认称号模板
	hs.initializeDefaultTitles()

	// 初始化默认成就模板
	hs.initializeDefaultAchievements()

	// 初始化荣誉等级
	hs.initializeHonorLevels()
}

// initializeDefaultTitles 初始化默认称号
func (hs *HonorService) initializeDefaultTitles() {
	// 新手称号
	newbieTitle := NewTitleTemplate("newbie", "新手冒险者", "刚踏上冒险之路的勇士", TitleCategorySpecial, TitleRarityCommon)
	newbieTitle.unlockConditions = append(newbieTitle.unlockConditions, NewLevelCondition(1))
	hs.RegisterTitleTemplate(newbieTitle)

	// 战斗称号
	warriorTitle := NewTitleTemplate("warrior", "勇敢战士", "在战斗中展现勇气的战士", TitleCategoryCombat, TitleRarityUncommon)
	warriorTitle.unlockConditions = append(warriorTitle.unlockConditions, NewStatisticCondition(StatisticTypeKillCount, 100))
	warriorTitle.attributeBonus["attack"] = 10
	hs.RegisterTitleTemplate(warriorTitle)

	// 探索称号
	explorerTitle := NewTitleTemplate("explorer", "大陆探索者", "足迹遍布大陆的探索者", TitleCategoryExploration, TitleRarityRare)
	explorerTitle.unlockConditions = append(explorerTitle.unlockConditions, NewStatisticCondition(StatisticTypeDistanceTraveled, 10000))
	explorerTitle.attributeBonus["speed"] = 15
	hs.RegisterTitleTemplate(explorerTitle)
}

// initializeDefaultAchievements 初始化默认成就
func (hs *HonorService) initializeDefaultAchievements() {
	// 首次击杀成就
	firstKill := NewAchievementTemplate("first_kill", "初次胜利", "获得第一次击杀", AchievementCategoryCombat, AchievementTypeNormal)
	firstKill.unlockConditions = append(firstKill.unlockConditions, NewStatisticCondition(StatisticTypeKillCount, 1))
	firstKill.honorReward = 10
	hs.RegisterAchievementTemplate(firstKill)

	// 连续登录成就
	loginStreak := NewAchievementTemplate("login_streak_7", "坚持不懈", "连续登录7天", AchievementCategoryProgression, AchievementTypeNormal)
	loginStreak.unlockConditions = append(loginStreak.unlockConditions, NewStatisticCondition(StatisticTypeLoginDays, 7))
	loginStreak.honorReward = 50
	hs.RegisterAchievementTemplate(loginStreak)

	// 收集成就
	collector := NewAchievementTemplate("collector", "收集家", "收集100个不同的物品", AchievementCategoryCollection, AchievementTypeNormal)
	collector.unlockConditions = append(collector.unlockConditions, NewStatisticCondition(StatisticTypeItemsCollected, 100))
	collector.honorReward = 100
	hs.RegisterAchievementTemplate(collector)
}

// initializeHonorLevels 初始化荣誉等级
func (hs *HonorService) initializeHonorLevels() {
	// 创建荣誉等级
	levels := []struct {
		level       int
		requiredXP  int
		title       string
		description string
	}{
		{1, 0, "平民", "普通的冒险者"},
		{2, 100, "见习者", "初入江湖的新人"},
		{3, 300, "冒险者", "有一定经验的冒险者"},
		{4, 600, "勇士", "勇敢的战士"},
		{5, 1000, "英雄", "受人尊敬的英雄"},
		{6, 1500, "传奇", "传说中的人物"},
		{7, 2100, "神话", "神话般的存在"},
		{8, 2800, "不朽", "不朽的传说"},
		{9, 3600, "至尊", "至高无上的存在"},
		{10, 4500, "神明", "如神明般的力量"},
	}

	for _, lvl := range levels {
		honorLevel := NewHonorLevel(lvl.level, lvl.requiredXP, lvl.title, lvl.description)
		hs.RegisterHonorLevel(honorLevel)
	}
}
