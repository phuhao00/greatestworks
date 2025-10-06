package hangup

import (
	"fmt"
	"math"
	"time"
)

// HangupService 挂机领域服务
type HangupService struct {
	config            *HangupConfig
	locationTemplates map[string]*LocationTemplate
	efficiencyRules   []EfficiencyRule
	rewardCalculators map[LocationType]RewardCalculator
}

// NewHangupService 创建挂机服务
func NewHangupService(config *HangupConfig) *HangupService {
	return &HangupService{
		config:            config,
		locationTemplates: make(map[string]*LocationTemplate),
		efficiencyRules:   make([]EfficiencyRule, 0),
		rewardCalculators: make(map[LocationType]RewardCalculator),
	}
}

// LocationTemplate 地点模板
type LocationTemplate struct {
	ID               string
	Name             string
	Description      string
	LocationType     LocationType
	RequiredLevel    int
	RequiredQuests   []string
	BaseExpRate      float64
	BaseGoldRate     float64
	SpecialItems     []ItemDrop
	MaxOfflineHours  int
	UnlockConditions []UnlockCondition
}

// CreateLocation 根据模板创建地点
func (lt *LocationTemplate) CreateLocation() *HangupLocation {
	location := NewHangupLocation(lt.ID, lt.Name, lt.Description, lt.LocationType)
	location.SetRequiredLevel(lt.RequiredLevel)
	location.SetBaseExpRate(lt.BaseExpRate)
	location.SetBaseGoldRate(lt.BaseGoldRate)
	location.SetMaxOfflineHours(lt.MaxOfflineHours)

	// 添加所需任务
	for _, questID := range lt.RequiredQuests {
		location.AddRequiredQuest(questID)
	}

	// 添加特殊物品
	for _, item := range lt.SpecialItems {
		location.AddSpecialItem(item)
	}

	return location
}

// UnlockCondition 解锁条件
type UnlockCondition struct {
	ConditionType string
	RequiredValue interface{}
	Description   string
}

// EfficiencyRule 效率规则
type EfficiencyRule struct {
	Name        string
	Condition   func(*HangupAggregate) bool
	BonusType   string // "vip", "equipment", "skill", "guild", "event"
	BonusValue  float64
	Description string
}

// RewardCalculator 奖励计算器接口
type RewardCalculator interface {
	CalculateReward(location *HangupLocation, duration time.Duration, bonus *EfficiencyBonus) *BaseReward
}

// DefaultRewardCalculator 默认奖励计算器
type DefaultRewardCalculator struct{}

// CalculateReward 计算奖励
func (drc *DefaultRewardCalculator) CalculateReward(location *HangupLocation, duration time.Duration, bonus *EfficiencyBonus) *BaseReward {
	hours := duration.Hours()

	// 基础奖励计算
	baseExp := int64(hours * 100 * location.GetBaseExpRate())
	baseGold := int64(hours * 50 * location.GetBaseGoldRate())

	// 应用地点类型倍率
	locationMultiplier := location.GetLocationType().GetExpMultiplier()
	baseExp = int64(float64(baseExp) * locationMultiplier)
	baseGold = int64(float64(baseGold) * location.GetLocationType().GetGoldMultiplier())

	baseReward := NewBaseReward(baseExp, baseGold)

	// 计算物品掉落
	for _, itemDrop := range location.GetSpecialItems() {
		if itemDrop.ShouldDrop(hours) {
			baseReward.AddItem(NewRewardItem(itemDrop.ItemID, itemDrop.CalculateQuantity(hours)))
		}
	}

	return baseReward
}

// RegisterLocationTemplate 注册地点模板
func (hs *HangupService) RegisterLocationTemplate(template *LocationTemplate) {
	hs.locationTemplates[template.ID] = template
}

// GetLocationTemplate 获取地点模板
func (hs *HangupService) GetLocationTemplate(id string) *LocationTemplate {
	return hs.locationTemplates[id]
}

// GetAllLocationTemplates 获取所有地点模板
func (hs *HangupService) GetAllLocationTemplates() map[string]*LocationTemplate {
	return hs.locationTemplates
}

// AddEfficiencyRule 添加效率规则
func (hs *HangupService) AddEfficiencyRule(rule EfficiencyRule) {
	hs.efficiencyRules = append(hs.efficiencyRules, rule)
}

// RegisterRewardCalculator 注册奖励计算器
func (hs *HangupService) RegisterRewardCalculator(locationType LocationType, calculator RewardCalculator) {
	hs.rewardCalculators[locationType] = calculator
}

// CreatePlayerHangup 为玩家创建挂机系统
func (hs *HangupService) CreatePlayerHangup(playerID string) *HangupAggregate {
	return NewHangupAggregate(playerID)
}

// CalculateOfflineReward 计算离线奖励
func (hs *HangupService) CalculateOfflineReward(hangup *HangupAggregate, offlineDuration time.Duration) (*OfflineReward, error) {
	if hangup.GetCurrentLocation() == nil {
		return nil, ErrNoHangupLocationSet
	}

	// 限制最大离线时间
	maxOfflineTime := time.Duration(hs.config.GetMaxOfflineHours()) * time.Hour
	if offlineDuration > maxOfflineTime {
		offlineDuration = maxOfflineTime
	}

	// 应用离线衰减
	offlineMultiplier := hs.config.GetOfflineDecayRate()
	effectiveDuration := time.Duration(float64(offlineDuration) * offlineMultiplier)

	// 获取奖励计算器
	location := hangup.GetCurrentLocation()
	calculator := hs.getRewardCalculator(location.GetLocationType())

	// 计算基础奖励
	baseReward := calculator.CalculateReward(location, effectiveDuration, hangup.GetEfficiencyBonus())

	// 应用效率加成
	finalReward := hangup.GetEfficiencyBonus().ApplyBonus(baseReward)

	// 创建离线奖励
	offlineReward := &OfflineReward{
		Experience:      finalReward.Experience,
		Gold:            finalReward.Gold,
		Items:           finalReward.Items,
		OfflineDuration: offlineDuration,
		LocationID:      location.GetID(),
		CalculatedAt:    time.Now(),
		IsClaimed:       false,
	}

	return offlineReward, nil
}

// UpdateEfficiencyBonus 更新效率加成
func (hs *HangupService) UpdateEfficiencyBonus(hangup *HangupAggregate, playerLevel int, vipLevel int, equipmentBonus float64) {
	bonus := NewEfficiencyBonus()

	// 应用所有效率规则
	for _, rule := range hs.efficiencyRules {
		if rule.Condition(hangup) {
			switch rule.BonusType {
			case "vip":
				bonus.SetVipBonus(bonus.GetVipBonus() + rule.BonusValue)
			case "equipment":
				bonus.SetEquipmentBonus(bonus.GetEquipmentBonus() + rule.BonusValue)
			case "skill":
				bonus.SetSkillBonus(bonus.GetSkillBonus() + rule.BonusValue)
			case "guild":
				bonus.SetGuildBonus(bonus.GetGuildBonus() + rule.BonusValue)
			case "event":
				bonus.SetEventBonus(bonus.GetEventBonus() + rule.BonusValue)
			default:
				bonus.SetSpecialBonus(rule.BonusType, rule.BonusValue)
			}
		}
	}

	// 设置VIP加成
	vipBonus := hs.calculateVipBonus(vipLevel)
	bonus.SetVipBonus(vipBonus)

	// 设置装备加成
	bonus.SetEquipmentBonus(equipmentBonus)

	hangup.UpdateEfficiencyBonus(bonus)
}

// ValidateLocationUnlock 验证地点解锁
func (hs *HangupService) ValidateLocationUnlock(playerID string, locationID string, playerLevel int, completedQuests []string) error {
	template := hs.GetLocationTemplate(locationID)
	if template == nil {
		return fmt.Errorf("location template not found: %s", locationID)
	}

	// 检查等级要求
	if playerLevel < template.RequiredLevel {
		return fmt.Errorf("player level %d is below required level %d", playerLevel, template.RequiredLevel)
	}

	// 检查任务要求
	completedQuestMap := make(map[string]bool)
	for _, questID := range completedQuests {
		completedQuestMap[questID] = true
	}

	for _, requiredQuest := range template.RequiredQuests {
		if !completedQuestMap[requiredQuest] {
			return fmt.Errorf("required quest not completed: %s", requiredQuest)
		}
	}

	return nil
}

// CalculateOptimalLocation 计算最优挂机地点
func (hs *HangupService) CalculateOptimalLocation(hangup *HangupAggregate, availableLocations []*HangupLocation, targetType string) *HangupLocation {
	if len(availableLocations) == 0 {
		return nil
	}

	var bestLocation *HangupLocation
	var bestScore float64

	for _, location := range availableLocations {
		if !location.IsUnlocked() || !location.IsActive() {
			continue
		}

		score := hs.calculateLocationScore(location, targetType)
		if bestLocation == nil || score > bestScore {
			bestLocation = location
			bestScore = score
		}
	}

	return bestLocation
}

// CalculateHangupEfficiency 计算挂机效率
func (hs *HangupService) CalculateHangupEfficiency(hangup *HangupAggregate) float64 {
	if hangup.GetCurrentLocation() == nil {
		return 0
	}

	location := hangup.GetCurrentLocation()
	bonus := hangup.GetEfficiencyBonus()

	// 基础效率
	baseEfficiency := location.GetBaseExpRate() + location.GetBaseGoldRate()

	// 应用加成
	totalBonus := bonus.GetTotalBonus()

	// 地点类型加成
	locationBonus := location.GetLocationType().GetExpMultiplier() + location.GetLocationType().GetGoldMultiplier()

	return baseEfficiency * totalBonus * locationBonus
}

// GetHangupRecommendations 获取挂机建议
func (hs *HangupService) GetHangupRecommendations(hangup *HangupAggregate, playerLevel int) []HangupRecommendation {
	recommendations := make([]HangupRecommendation, 0)

	// 检查当前地点是否最优
	currentLocation := hangup.GetCurrentLocation()
	if currentLocation != nil {
		efficiency := hs.CalculateHangupEfficiency(hangup)
		if efficiency < 2.0 { // 效率较低
			recommendations = append(recommendations, HangupRecommendation{
				Type:        "location_change",
				Title:       "建议更换挂机地点",
				Description: "当前地点效率较低，建议选择更高效的地点",
				Priority:    "medium",
			})
		}
	}

	// 检查效率加成
	bonus := hangup.GetEfficiencyBonus()
	if bonus.GetTotalBonus() < 1.5 {
		recommendations = append(recommendations, HangupRecommendation{
			Type:        "efficiency_boost",
			Title:       "提升挂机效率",
			Description: "通过提升VIP等级、装备或技能来增加挂机效率",
			Priority:    "high",
		})
	}

	// 检查每日挂机时间
	dailyTime := hangup.GetDailyHangupTime()
	maxDailyTime := time.Duration(hs.config.GetMaxDailyHangupHours()) * time.Hour
	if dailyTime < time.Duration(float64(maxDailyTime)*0.8) { //todo need check it correct
		recommendations = append(recommendations, HangupRecommendation{
			Type:        "time_optimization",
			Title:       "增加挂机时间",
			Description: "今日挂机时间较少，建议增加挂机时间以获得更多收益",
			Priority:    "low",
		})
	}

	return recommendations
}

// 私有方法

// getRewardCalculator 获取奖励计算器
func (hs *HangupService) getRewardCalculator(locationType LocationType) RewardCalculator {
	if calculator, exists := hs.rewardCalculators[locationType]; exists {
		return calculator
	}
	return &DefaultRewardCalculator{}
}

// calculateVipBonus 计算VIP加成
func (hs *HangupService) calculateVipBonus(vipLevel int) float64 {
	if vipLevel <= 0 {
		return 0
	}

	// VIP等级越高，加成越大，但有递减效应
	bonus := float64(vipLevel) * 0.1 // 每级10%加成
	if vipLevel > 10 {
		bonus = 1.0 + math.Log(float64(vipLevel-10))*0.1 // 超过10级后递减
	}

	return bonus
}

// calculateLocationScore 计算地点评分
func (hs *HangupService) calculateLocationScore(location *HangupLocation, targetType string) float64 {
	baseScore := location.GetBaseExpRate() + location.GetBaseGoldRate()

	// 根据目标类型调整评分
	switch targetType {
	case "experience":
		baseScore = location.GetBaseExpRate()*2 + location.GetBaseGoldRate()*0.5
	case "gold":
		baseScore = location.GetBaseGoldRate()*2 + location.GetBaseExpRate()*0.5
	case "items":
		baseScore += float64(len(location.GetSpecialItems())) * 0.5
	default:
		// 平衡型
		baseScore = location.GetBaseExpRate() + location.GetBaseGoldRate()
	}

	// 地点类型加成
	locationMultiplier := location.GetLocationType().GetExpMultiplier() + location.GetLocationType().GetGoldMultiplier()
	baseScore *= locationMultiplier

	return baseScore
}

// InitializeDefaultLocations 初始化默认地点
func (hs *HangupService) InitializeDefaultLocations() {
	// 新手森林
	beginnerForest := &LocationTemplate{
		ID:              "beginner_forest",
		Name:            "新手森林",
		Description:     "适合新手的安全森林区域",
		LocationType:    LocationTypeForest,
		RequiredLevel:   1,
		RequiredQuests:  []string{},
		BaseExpRate:     1.0,
		BaseGoldRate:    1.0,
		MaxOfflineHours: 12,
		SpecialItems: []ItemDrop{
			NewItemDrop("wood", 0.3, 1, 3),
			NewItemDrop("herb", 0.2, 1, 2),
		},
	}
	hs.RegisterLocationTemplate(beginnerForest)

	// 魔法洞穴
	magicCave := &LocationTemplate{
		ID:              "magic_cave",
		Name:            "魔法洞穴",
		Description:     "充满魔法能量的神秘洞穴",
		LocationType:    LocationTypeCave,
		RequiredLevel:   10,
		RequiredQuests:  []string{"explore_forest"},
		BaseExpRate:     1.4,
		BaseGoldRate:    1.2,
		MaxOfflineHours: 18,
		SpecialItems: []ItemDrop{
			NewItemDrop("magic_crystal", 0.1, 1, 1),
			NewItemDrop("rare_ore", 0.15, 1, 2),
		},
	}
	hs.RegisterLocationTemplate(magicCave)

	// 古代遗迹
	ancientRuins := &LocationTemplate{
		ID:              "ancient_ruins",
		Name:            "古代遗迹",
		Description:     "蕴含古老力量的神秘遗迹",
		LocationType:    LocationTypeSpecial,
		RequiredLevel:   25,
		RequiredQuests:  []string{"cave_exploration", "ancient_key"},
		BaseExpRate:     2.0,
		BaseGoldRate:    1.8,
		MaxOfflineHours: 24,
		SpecialItems: []ItemDrop{
			NewItemDrop("ancient_artifact", 0.05, 1, 1),
			NewItemDrop("legendary_gem", 0.02, 1, 1),
		},
	}
	hs.RegisterLocationTemplate(ancientRuins)
}

// InitializeDefaultRules 初始化默认规则
func (hs *HangupService) InitializeDefaultRules() {
	// VIP加成规则
	hs.AddEfficiencyRule(EfficiencyRule{
		Name: "VIP Bonus",
		Condition: func(hangup *HangupAggregate) bool {
			return true // 所有玩家都适用
		},
		BonusType:   "vip",
		BonusValue:  0.0, // 将在UpdateEfficiencyBonus中计算
		Description: "VIP等级加成",
	})

	// 长时间挂机加成
	hs.AddEfficiencyRule(EfficiencyRule{
		Name: "Long Session Bonus",
		Condition: func(hangup *HangupAggregate) bool {
			return hangup.GetDailyHangupTime() > 6*time.Hour
		},
		BonusType:   "special",
		BonusValue:  0.1, // 10%加成
		Description: "长时间挂机加成",
	})

	// 连续挂机加成
	hs.AddEfficiencyRule(EfficiencyRule{
		Name: "Consecutive Days Bonus",
		Condition: func(hangup *HangupAggregate) bool {
			// 这里需要外部提供连续挂机天数信息
			return true // 简化实现
		},
		BonusType:   "special",
		BonusValue:  0.05, // 5%加成
		Description: "连续挂机加成",
	})
}

// HangupRecommendation 挂机建议
type HangupRecommendation struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"` // "high", "medium", "low"
	ActionURL   string `json:"action_url,omitempty"`
}
