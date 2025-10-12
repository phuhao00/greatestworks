package sacred

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// SacredService 圣地领域服务
type SacredService struct {
	challengeTemplates map[ChallengeType]*ChallengeTemplate
	blessingTemplates  map[BlessingType]*BlessingTemplate
	relicTemplates     map[RelicType][]*RelicTemplate
	difficultyCurves   map[ChallengeDifficulty]*DifficultyCurve
	rewardCalculator   *RewardCalculator
	balanceRules       *BalanceRules
}

// NewSacredService 创建圣地服务
func NewSacredService() *SacredService {
	service := &SacredService{
		challengeTemplates: make(map[ChallengeType]*ChallengeTemplate),
		blessingTemplates:  make(map[BlessingType]*BlessingTemplate),
		relicTemplates:     make(map[RelicType][]*RelicTemplate),
		difficultyCurves:   make(map[ChallengeDifficulty]*DifficultyCurve),
		rewardCalculator:   NewRewardCalculator(),
		balanceRules:       NewBalanceRules(),
	}

	// 初始化默认模板和规则
	service.initializeDefaultTemplates()
	service.initializeDifficultyCurves()

	return service
}

// CreateSacredPlace 创建圣地
func (s *SacredService) CreateSacredPlace(id, name, description, owner string) (*SacredPlaceAggregate, error) {
	if id == "" || name == "" || owner == "" {
		return nil, fmt.Errorf("invalid parameters for sacred place creation")
	}

	sacredPlace := NewSacredPlaceAggregate(id, name, description, owner)

	// 添加默认挑战
	defaultChallenges := s.generateDefaultChallenges(sacredPlace.GetLevel().Level)
	for _, challenge := range defaultChallenges {
		sacredPlace.AddChallenge(challenge)
	}

	// 添加默认祝福
	defaultBlessings := s.generateDefaultBlessings(sacredPlace.GetLevel().Level)
	for _, blessing := range defaultBlessings {
		sacredPlace.AddBlessing(blessing)
	}

	return sacredPlace, nil
}

// GenerateChallenge 生成挑战
func (s *SacredService) GenerateChallenge(challengeType ChallengeType, difficulty ChallengeDifficulty, sacredLevel int) (*Challenge, error) {
	template, exists := s.challengeTemplates[challengeType]
	if !exists {
		return nil, fmt.Errorf("challenge template not found for type: %s", challengeType.String())
	}

	// 生成唯一ID
	id := fmt.Sprintf("challenge_%s_%s_%d", challengeType.String(), difficulty.String(), time.Now().UnixNano())

	// 根据模板创建挑战
	challenge := NewChallenge(
		id,
		template.GenerateName(difficulty),
		template.GenerateDescription(difficulty),
		challengeType,
		difficulty,
		difficulty.GetRequiredLevel(),
	)

	// 设置持续时间和冷却时间
	challenge.SetDuration(template.GetDuration(difficulty))
	challenge.SetCooldown(template.GetCooldown(difficulty))

	// 添加条件
	conditions := template.GenerateConditions(difficulty, sacredLevel)
	for _, condition := range conditions {
		challenge.AddCondition(condition)
	}

	return challenge, nil
}

// GenerateBlessing 生成祝福
func (s *SacredService) GenerateBlessing(blessingType BlessingType, sacredLevel int) (*Blessing, error) {
	template, exists := s.blessingTemplates[blessingType]
	if !exists {
		return nil, fmt.Errorf("blessing template not found for type: %s", blessingType.String())
	}

	// 生成唯一ID
	id := fmt.Sprintf("blessing_%s_%d", blessingType.String(), time.Now().UnixNano())

	// 根据模板创建祝福
	blessing := NewBlessing(
		id,
		template.GenerateName(sacredLevel),
		template.GenerateDescription(sacredLevel),
		blessingType,
		template.GetDuration(sacredLevel),
	)

	// 设置冷却时间和最大使用次数
	blessing.SetCooldown(template.GetCooldown(sacredLevel))
	blessing.SetMaxUsage(template.GetMaxUsage(sacredLevel))

	// 添加效果
	effects := template.GenerateEffects(sacredLevel)
	for _, effect := range effects {
		blessing.AddEffect(effect)
	}

	return blessing, nil
}

// GenerateRelic 生成圣物
func (s *SacredService) GenerateRelic(relicType RelicType, rarity RelicRarity) (*SacredRelic, error) {
	templates, exists := s.relicTemplates[relicType]
	if !exists || len(templates) == 0 {
		return nil, fmt.Errorf("relic templates not found for type: %s", relicType.String())
	}

	// 随机选择模板
	template := templates[rand.Intn(len(templates))]

	// 生成唯一ID
	id := fmt.Sprintf("relic_%s_%s_%d", relicType.String(), rarity.String(), time.Now().UnixNano())

	// 根据模板创建圣物
	relic := NewSacredRelic(
		id,
		template.GenerateName(rarity),
		template.GenerateDescription(rarity),
		relicType,
		rarity,
	)

	// 添加属性
	attributes := template.GenerateAttributes(rarity)
	for name, value := range attributes {
		relic.AddAttribute(name, value)
	}

	// 添加效果
	effects := template.GenerateEffects(rarity)
	for _, effect := range effects {
		relic.AddEffect(effect)
	}

	// 添加需求
	requirements := template.GenerateRequirements(rarity)
	for name, value := range requirements {
		relic.AddRequirement(name, value)
	}

	return relic, nil
}

// CalculateChallengeReward 计算挑战奖励
func (s *SacredService) CalculateChallengeReward(challengeType ChallengeType, difficulty ChallengeDifficulty, success bool, score int, playerLevel int) *ChallengeReward {
	return s.rewardCalculator.CalculateChallengeReward(challengeType, difficulty, success, score, playerLevel)
}

// CalculateBlessingEffect 计算祝福效果
func (s *SacredService) CalculateBlessingEffect(blessingType BlessingType, sacredLevel int, playerLevel int) *BlessingEffect {
	effect := NewBlessingEffect(
		fmt.Sprintf("blessing_effect_%d", time.Now().UnixNano()),
		"", // playerID will be set when activated
		blessingType,
		time.Hour, // default duration
	)

	// 根据类型计算效果
	switch blessingType {
	case BlessingTypeAttribute:
		effect.AddAttribute("strength", float64(sacredLevel*5+playerLevel))
		effect.AddAttribute("agility", float64(sacredLevel*3+playerLevel/2))
	case BlessingTypeSkill:
		effect.AddModifier("skill_damage", 1.0+float64(sacredLevel)*0.1)
		effect.AddModifier("skill_cooldown", 0.9-float64(sacredLevel)*0.05)
	case BlessingTypeExperience:
		effect.AddModifier("exp_multiplier", 1.0+float64(sacredLevel)*0.2)
	case BlessingTypeWealth:
		effect.AddModifier("gold_multiplier", 1.0+float64(sacredLevel)*0.15)
	case BlessingTypeProtection:
		effect.AddAttribute("defense", float64(sacredLevel*10+playerLevel*2))
		effect.AddModifier("damage_reduction", float64(sacredLevel)*0.05)
	case BlessingTypeHealing:
		effect.AddAttribute("health_regen", float64(sacredLevel*2+playerLevel/5))
		effect.AddModifier("healing_received", 1.0+float64(sacredLevel)*0.1)
	case BlessingTypeSpeed:
		effect.AddAttribute("movement_speed", float64(sacredLevel*5))
		effect.AddModifier("action_speed", 1.0+float64(sacredLevel)*0.08)
	case BlessingTypeLuck:
		effect.AddAttribute("luck", float64(sacredLevel*3+playerLevel/3))
		effect.AddModifier("critical_chance", float64(sacredLevel)*0.02)
	}

	return effect
}

// ValidateChallenge 验证挑战
func (s *SacredService) ValidateChallenge(challenge *Challenge, playerData map[string]interface{}) error {
	if challenge == nil {
		return fmt.Errorf("challenge is nil")
	}

	// 检查挑战状态
	if !challenge.CanStart() {
		return fmt.Errorf("challenge cannot be started")
	}

	// 检查玩家等级
	playerLevel, ok := playerData["level"].(int)
	if !ok || playerLevel < challenge.GetRequiredLevel() {
		return fmt.Errorf("insufficient player level")
	}

	// 检查挑战条件
	if !challenge.CheckConditions(playerData) {
		return fmt.Errorf("challenge conditions not met")
	}

	return nil
}

// ValidateBlessing 验证祝福
func (s *SacredService) ValidateBlessing(blessing *Blessing, playerData map[string]interface{}) error {
	if blessing == nil {
		return fmt.Errorf("blessing is nil")
	}

	// 检查祝福状态
	if !blessing.IsAvailable() {
		return fmt.Errorf("blessing is not available")
	}

	// 检查平衡规则
	if !s.balanceRules.CanActivateBlessing(blessing.GetType(), playerData) {
		return fmt.Errorf("blessing activation violates balance rules")
	}

	return nil
}

// CalculateOptimalDifficulty 计算最佳难度
func (s *SacredService) CalculateOptimalDifficulty(playerLevel int, playerSkill float64) ChallengeDifficulty {
	// 基于玩家等级和技能计算推荐难度
	baseScore := float64(playerLevel) + playerSkill*10

	if baseScore < 20 {
		return ChallengeDifficultyEasy
	} else if baseScore < 50 {
		return ChallengeDifficultyNormal
	} else if baseScore < 100 {
		return ChallengeDifficultyHard
	} else if baseScore < 200 {
		return ChallengeDifficultyExpert
	} else {
		return ChallengeDifficultyLegendary
	}
}

// GetRecommendedChallenges 获取推荐挑战
func (s *SacredService) GetRecommendedChallenges(playerData map[string]interface{}, sacredLevel int) []*Challenge {
	playerLevel, _ := playerData["level"].(int)
	playerSkill, _ := playerData["skill"].(float64)

	optimalDifficulty := s.CalculateOptimalDifficulty(playerLevel, playerSkill)

	var recommendations []*Challenge

	// 为每种挑战类型生成推荐
	for challengeType := ChallengeTypeCombat; challengeType <= ChallengeTypeSpecial; challengeType++ {
		if challenge, err := s.GenerateChallenge(challengeType, optimalDifficulty, sacredLevel); err == nil {
			recommendations = append(recommendations, challenge)
		}
	}

	return recommendations
}

// GetAvailableBlessings 获取可用祝福
func (s *SacredService) GetAvailableBlessings(playerData map[string]interface{}, sacredLevel int) []*Blessing {
	var available []*Blessing

	// 为每种祝福类型生成可用祝福
	for blessingType := BlessingTypeAttribute; blessingType <= BlessingTypeLuck; blessingType++ {
		if blessing, err := s.GenerateBlessing(blessingType, sacredLevel); err == nil {
			if s.ValidateBlessing(blessing, playerData) == nil {
				available = append(available, blessing)
			}
		}
	}

	return available
}

// 私有方法

// generateDefaultChallenges 生成默认挑战
func (s *SacredService) generateDefaultChallenges(sacredLevel int) []*Challenge {
	var challenges []*Challenge

	// 根据圣地等级生成适当的挑战
	difficulties := []ChallengeDifficulty{ChallengeDifficultyEasy, ChallengeDifficultyNormal}
	if sacredLevel >= 10 {
		difficulties = append(difficulties, ChallengeDifficultyHard)
	}
	if sacredLevel >= 25 {
		difficulties = append(difficulties, ChallengeDifficultyExpert)
	}
	if sacredLevel >= 50 {
		difficulties = append(difficulties, ChallengeDifficultyLegendary)
	}

	// 为每种难度生成战斗挑战
	for _, difficulty := range difficulties {
		if challenge, err := s.GenerateChallenge(ChallengeTypeCombat, difficulty, sacredLevel); err == nil {
			challenges = append(challenges, challenge)
		}
	}

	return challenges
}

// generateDefaultBlessings 生成默认祝福
func (s *SacredService) generateDefaultBlessings(sacredLevel int) []*Blessing {
	var blessings []*Blessing

	// 生成基础祝福
	basicTypes := []BlessingType{BlessingTypeAttribute, BlessingTypeExperience, BlessingTypeWealth}
	for _, blessingType := range basicTypes {
		if blessing, err := s.GenerateBlessing(blessingType, sacredLevel); err == nil {
			blessings = append(blessings, blessing)
		}
	}

	// 根据等级解锁高级祝福
	if sacredLevel >= 10 {
		advancedTypes := []BlessingType{BlessingTypeProtection, BlessingTypeHealing}
		for _, blessingType := range advancedTypes {
			if blessing, err := s.GenerateBlessing(blessingType, sacredLevel); err == nil {
				blessings = append(blessings, blessing)
			}
		}
	}

	if sacredLevel >= 25 {
		specialTypes := []BlessingType{BlessingTypeSpeed, BlessingTypeLuck}
		for _, blessingType := range specialTypes {
			if blessing, err := s.GenerateBlessing(blessingType, sacredLevel); err == nil {
				blessings = append(blessings, blessing)
			}
		}
	}

	return blessings
}

// initializeDefaultTemplates 初始化默认模板
func (s *SacredService) initializeDefaultTemplates() {
	// 初始化挑战模板
	s.challengeTemplates[ChallengeTypeCombat] = NewChallengeTemplate(
		"战斗挑战",
		"测试战斗技巧的挑战",
		time.Minute*30,
		time.Hour*6,
	)

	s.challengeTemplates[ChallengeTypePuzzle] = NewChallengeTemplate(
		"解谜挑战",
		"需要智慧解决的谜题",
		time.Minute*15,
		time.Hour*4,
	)

	s.challengeTemplates[ChallengeTypeEndurance] = NewChallengeTemplate(
		"耐力挑战",
		"考验持久力的挑战",
		time.Hour,
		time.Hour*12,
	)

	// 初始化祝福模板
	s.blessingTemplates[BlessingTypeAttribute] = NewBlessingTemplate(
		"属性祝福",
		"提升基础属性",
		time.Hour*2,
		time.Hour*24,
		3,
	)

	s.blessingTemplates[BlessingTypeExperience] = NewBlessingTemplate(
		"经验祝福",
		"增加经验获取",
		time.Hour,
		time.Hour*12,
		5,
	)

	// 初始化圣物模板
	s.initializeRelicTemplates()
}

// initializeRelicTemplates 初始化圣物模板
func (s *SacredService) initializeRelicTemplates() {
	// 武器模板
	weaponTemplates := []*RelicTemplate{
		NewRelicTemplate("圣剑", "神圣的武器", []string{"attack", "critical"}, []string{"增加攻击力", "提高暴击率"}),
		NewRelicTemplate("法杖", "魔法武器", []string{"magic_power", "mana"}, []string{"增加魔法攻击", "提高法力值"}),
	}
	s.relicTemplates[RelicTypeWeapon] = weaponTemplates

	// 护甲模板
	armorTemplates := []*RelicTemplate{
		NewRelicTemplate("圣甲", "神圣的护甲", []string{"defense", "health"}, []string{"增加防御力", "提高生命值"}),
		NewRelicTemplate("法袍", "魔法护甲", []string{"magic_defense", "mana_regen"}, []string{"增加魔法防御", "提高法力回复"}),
	}
	s.relicTemplates[RelicTypeArmor] = armorTemplates

	// 饰品模板
	accessoryTemplates := []*RelicTemplate{
		NewRelicTemplate("圣环", "神圣的戒指", []string{"luck", "experience"}, []string{"增加幸运值", "提高经验获取"}),
		NewRelicTemplate("护符", "保护饰品", []string{"resistance", "health_regen"}, []string{"增加抗性", "提高生命回复"}),
	}
	s.relicTemplates[RelicTypeAccessory] = accessoryTemplates
}

// initializeDifficultyCurves 初始化难度曲线
func (s *SacredService) initializeDifficultyCurves() {
	s.difficultyCurves[ChallengeDifficultyEasy] = &DifficultyCurve{
		HealthMultiplier: 0.5,
		DamageMultiplier: 0.7,
		SpeedMultiplier:  0.8,
		RewardMultiplier: 0.5,
		ExpMultiplier:    0.3,
	}

	s.difficultyCurves[ChallengeDifficultyNormal] = &DifficultyCurve{
		HealthMultiplier: 1.0,
		DamageMultiplier: 1.0,
		SpeedMultiplier:  1.0,
		RewardMultiplier: 1.0,
		ExpMultiplier:    1.0,
	}

	s.difficultyCurves[ChallengeDifficultyHard] = &DifficultyCurve{
		HealthMultiplier: 1.5,
		DamageMultiplier: 1.3,
		SpeedMultiplier:  1.2,
		RewardMultiplier: 1.5,
		ExpMultiplier:    1.8,
	}

	s.difficultyCurves[ChallengeDifficultyExpert] = &DifficultyCurve{
		HealthMultiplier: 2.0,
		DamageMultiplier: 1.8,
		SpeedMultiplier:  1.5,
		RewardMultiplier: 2.5,
		ExpMultiplier:    3.0,
	}

	s.difficultyCurves[ChallengeDifficultyLegendary] = &DifficultyCurve{
		HealthMultiplier: 3.0,
		DamageMultiplier: 2.5,
		SpeedMultiplier:  2.0,
		RewardMultiplier: 5.0,
		ExpMultiplier:    8.0,
	}
}

// 辅助结构体

// ChallengeTemplate 挑战模板
type ChallengeTemplate struct {
	Name         string
	Description  string
	BaseDuration time.Duration
	BaseCooldown time.Duration
}

// NewChallengeTemplate 创建挑战模板
func NewChallengeTemplate(name, description string, baseDuration, baseCooldown time.Duration) *ChallengeTemplate {
	return &ChallengeTemplate{
		Name:         name,
		Description:  description,
		BaseDuration: baseDuration,
		BaseCooldown: baseCooldown,
	}
}

// GenerateName 生成名称
func (ct *ChallengeTemplate) GenerateName(difficulty ChallengeDifficulty) string {
	return fmt.Sprintf("%s (%s)", ct.Name, difficulty.String())
}

// GenerateDescription 生成描述
func (ct *ChallengeTemplate) GenerateDescription(difficulty ChallengeDifficulty) string {
	return fmt.Sprintf("%s - 难度: %s", ct.Description, difficulty.String())
}

// GetDuration 获取持续时间
func (ct *ChallengeTemplate) GetDuration(difficulty ChallengeDifficulty) time.Duration {
	multiplier := difficulty.GetMultiplier()
	return time.Duration(float64(ct.BaseDuration) * multiplier)
}

// GetCooldown 获取冷却时间
func (ct *ChallengeTemplate) GetCooldown(difficulty ChallengeDifficulty) time.Duration {
	multiplier := difficulty.GetMultiplier()
	return time.Duration(float64(ct.BaseCooldown) * multiplier)
}

// GenerateConditions 生成条件
func (ct *ChallengeTemplate) GenerateConditions(difficulty ChallengeDifficulty, sacredLevel int) []*ChallengeCondition {
	conditions := []*ChallengeCondition{
		NewChallengeCondition("level", "level", "gte", difficulty.GetRequiredLevel(), "等级不足"),
	}

	// 根据难度添加额外条件
	if difficulty >= ChallengeDifficultyHard {
		conditions = append(conditions, NewChallengeCondition("equipment", "power", "gte", 1000, "装备威力不足"))
	}

	return conditions
}

// BlessingTemplate 祝福模板
type BlessingTemplate struct {
	Name         string
	Description  string
	BaseDuration time.Duration
	BaseCooldown time.Duration
	BaseMaxUsage int
}

// NewBlessingTemplate 创建祝福模板
func NewBlessingTemplate(name, description string, baseDuration, baseCooldown time.Duration, baseMaxUsage int) *BlessingTemplate {
	return &BlessingTemplate{
		Name:         name,
		Description:  description,
		BaseDuration: baseDuration,
		BaseCooldown: baseCooldown,
		BaseMaxUsage: baseMaxUsage,
	}
}

// GenerateName 生成名称
func (bt *BlessingTemplate) GenerateName(sacredLevel int) string {
	return fmt.Sprintf("%s (Lv.%d)", bt.Name, sacredLevel)
}

// GenerateDescription 生成描述
func (bt *BlessingTemplate) GenerateDescription(sacredLevel int) string {
	return fmt.Sprintf("%s - 圣地等级: %d", bt.Description, sacredLevel)
}

// GetDuration 获取持续时间
func (bt *BlessingTemplate) GetDuration(sacredLevel int) time.Duration {
	multiplier := 1.0 + float64(sacredLevel)*0.1
	return time.Duration(float64(bt.BaseDuration) * multiplier)
}

// GetCooldown 获取冷却时间
func (bt *BlessingTemplate) GetCooldown(sacredLevel int) time.Duration {
	multiplier := math.Max(0.5, 1.0-float64(sacredLevel)*0.02)
	return time.Duration(float64(bt.BaseCooldown) * multiplier)
}

// GetMaxUsage 获取最大使用次数
func (bt *BlessingTemplate) GetMaxUsage(sacredLevel int) int {
	return bt.BaseMaxUsage + sacredLevel/10
}

// GenerateEffects 生成效果
func (bt *BlessingTemplate) GenerateEffects(sacredLevel int) []*BlessingEffect {
	// 这里可以根据圣地等级生成不同的效果
	return []*BlessingEffect{}
}

// RelicTemplate 圣物模板
type RelicTemplate struct {
	Name        string
	Description string
	Attributes  []string
	Effects     []string
}

// NewRelicTemplate 创建圣物模板
func NewRelicTemplate(name, description string, attributes, effects []string) *RelicTemplate {
	return &RelicTemplate{
		Name:        name,
		Description: description,
		Attributes:  attributes,
		Effects:     effects,
	}
}

// GenerateName 生成名称
func (rt *RelicTemplate) GenerateName(rarity RelicRarity) string {
	return fmt.Sprintf("%s (%s)", rt.Name, rarity.String())
}

// GenerateDescription 生成描述
func (rt *RelicTemplate) GenerateDescription(rarity RelicRarity) string {
	return fmt.Sprintf("%s - 稀有度: %s", rt.Description, rarity.String())
}

// GenerateAttributes 生成属性
func (rt *RelicTemplate) GenerateAttributes(rarity RelicRarity) map[string]float64 {
	attributes := make(map[string]float64)
	basePower := rarity.GetBasePower()

	for _, attr := range rt.Attributes {
		attributes[attr] = basePower * (0.5 + rand.Float64())
	}

	return attributes
}

// GenerateEffects 生成效果
func (rt *RelicTemplate) GenerateEffects(rarity RelicRarity) []string {
	// 根据稀有度返回部分效果
	maxEffects := int(rarity) // 稀有度越高，效果越多
	if maxEffects > len(rt.Effects) {
		maxEffects = len(rt.Effects)
	}

	return rt.Effects[:maxEffects]
}

// GenerateRequirements 生成需求
func (rt *RelicTemplate) GenerateRequirements(rarity RelicRarity) map[string]interface{} {
	requirements := make(map[string]interface{})
	requirements["level"] = int(rarity) * 10 // 稀有度越高，等级要求越高
	return requirements
}

// DifficultyCurve 难度曲线
type DifficultyCurve struct {
	HealthMultiplier float64
	DamageMultiplier float64
	SpeedMultiplier  float64
	RewardMultiplier float64
	ExpMultiplier    float64
}

// RewardCalculator 奖励计算器
type RewardCalculator struct{}

// NewRewardCalculator 创建奖励计算器
func NewRewardCalculator() *RewardCalculator {
	return &RewardCalculator{}
}

// CalculateChallengeReward 计算挑战奖励
func (rc *RewardCalculator) CalculateChallengeReward(challengeType ChallengeType, difficulty ChallengeDifficulty, success bool, score int, playerLevel int) *ChallengeReward {
	if !success {
		return &ChallengeReward{
			Gold:       0,
			Experience: 0,
			Items:      make(map[string]int),
			Special:    make(map[string]interface{}),
		}
	}

	baseGold := 100
	baseExp := 50

	// 难度倍数
	difficultyMultiplier := difficulty.GetMultiplier()

	// 分数倍数
	scoreMultiplier := math.Max(0.1, math.Min(2.0, float64(score)/100.0))

	// 玩家等级影响
	levelMultiplier := 1.0 + float64(playerLevel)*0.05

	// 计算最终奖励
	finalGold := int(float64(baseGold) * difficultyMultiplier * scoreMultiplier * levelMultiplier)
	finalExp := int(float64(baseExp) * difficultyMultiplier * scoreMultiplier * levelMultiplier)

	reward := &ChallengeReward{
		Gold:       finalGold,
		Experience: finalExp,
		Items:      make(map[string]int),
		Special:    make(map[string]interface{}),
	}

	// 根据挑战类型添加特殊奖励
	switch challengeType {
	case ChallengeTypeCombat:
		reward.AddItem("combat_token", 1)
	case ChallengeTypePuzzle:
		reward.AddItem("wisdom_crystal", 1)
	case ChallengeTypeEndurance:
		reward.AddItem("endurance_potion", 1)
	}

	return reward
}

// BalanceRules 平衡规则
type BalanceRules struct {
	maxActiveBlessings int
	blessingCooldowns  map[BlessingType]time.Duration
}

// NewBalanceRules 创建平衡规则
func NewBalanceRules() *BalanceRules {
	return &BalanceRules{
		maxActiveBlessings: 3,
		blessingCooldowns:  make(map[BlessingType]time.Duration),
	}
}

// CanActivateBlessing 检查是否可以激活祝福
func (br *BalanceRules) CanActivateBlessing(blessingType BlessingType, playerData map[string]interface{}) bool {
	// 检查激活的祝福数量
	activeBlessings, _ := playerData["active_blessings"].(int)
	if activeBlessings >= br.maxActiveBlessings {
		return false
	}

	// 检查特定类型的冷却时间
	if cooldown, exists := br.blessingCooldowns[blessingType]; exists {
		lastUsed, _ := playerData[fmt.Sprintf("last_used_%s", blessingType.String())].(time.Time)
		if !lastUsed.IsZero() && time.Since(lastUsed) < cooldown {
			return false
		}
	}

	return true
}
