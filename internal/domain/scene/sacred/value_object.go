package sacred

import (
	"fmt"
	"time"
)

// SacredLevel 圣地等级值对象
type SacredLevel struct {
	Level      int
	Experience int
	MaxExp     int
}

// NewSacredLevel 创建圣地等级
func NewSacredLevel(level, experience int) *SacredLevel {
	return &SacredLevel{
		Level:      level,
		Experience: experience,
		MaxExp:     calculateMaxExp(level),
	}
}

// AddExperience 添加经验
func (sl *SacredLevel) AddExperience(exp int) (int, error) {
	if exp <= 0 {
		return sl.Level, fmt.Errorf("experience must be positive")
	}
	
	sl.Experience += exp
	oldLevel := sl.Level
	
	// 检查是否可以升级
	for sl.Experience >= sl.MaxExp {
		sl.Experience -= sl.MaxExp
		sl.Level++
		sl.MaxExp = calculateMaxExp(sl.Level)
	}
	
	return sl.Level, nil
}

// GetProgress 获取升级进度
func (sl *SacredLevel) GetProgress() float64 {
	if sl.MaxExp == 0 {
		return 0
	}
	return float64(sl.Experience) / float64(sl.MaxExp)
}

// GetRemainingExp 获取升级所需经验
func (sl *SacredLevel) GetRemainingExp() int {
	return sl.MaxExp - sl.Experience
}

// CanUpgrade 检查是否可以升级
func (sl *SacredLevel) CanUpgrade() bool {
	return sl.Experience >= sl.MaxExp
}

// ToMap 转换为映射
func (sl *SacredLevel) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"level":      sl.Level,
		"experience": sl.Experience,
		"max_exp":    sl.MaxExp,
		"progress":   sl.GetProgress(),
		"remaining":  sl.GetRemainingExp(),
	}
}

// calculateMaxExp 计算等级所需最大经验
func calculateMaxExp(level int) int {
	// 经验公式：level * 100 + (level-1) * 50
	return level*100 + (level-1)*50
}

// ChallengeType 挑战类型
type ChallengeType int

const (
	ChallengeTypeCombat     ChallengeType = iota + 1 // 战斗挑战
	ChallengeTypePuzzle                              // 解谜挑战
	ChallengeTypeEndurance                           // 耐力挑战
	ChallengeTypeSpeed                               // 速度挑战
	ChallengeTypeStrategy                            // 策略挑战
	ChallengeTypeCooperation                         // 合作挑战
	ChallengeTypeSpecial                             // 特殊挑战
)

// String 返回类型字符串
func (ct ChallengeType) String() string {
	switch ct {
	case ChallengeTypeCombat:
		return "combat"
	case ChallengeTypePuzzle:
		return "puzzle"
	case ChallengeTypeEndurance:
		return "endurance"
	case ChallengeTypeSpeed:
		return "speed"
	case ChallengeTypeStrategy:
		return "strategy"
	case ChallengeTypeCooperation:
		return "cooperation"
	case ChallengeTypeSpecial:
		return "special"
	default:
		return "unknown"
	}
}

// IsValid 检查类型是否有效
func (ct ChallengeType) IsValid() bool {
	return ct >= ChallengeTypeCombat && ct <= ChallengeTypeSpecial
}

// GetDescription 获取类型描述
func (ct ChallengeType) GetDescription() string {
	switch ct {
	case ChallengeTypeCombat:
		return "测试战斗技巧和策略的挑战"
	case ChallengeTypePuzzle:
		return "需要智慧和逻辑思维的解谜挑战"
	case ChallengeTypeEndurance:
		return "考验持久力和毅力的挑战"
	case ChallengeTypeSpeed:
		return "需要快速反应和敏捷的挑战"
	case ChallengeTypeStrategy:
		return "需要深度思考和规划的策略挑战"
	case ChallengeTypeCooperation:
		return "需要团队合作完成的挑战"
	case ChallengeTypeSpecial:
		return "独特的特殊挑战"
	default:
		return "未知类型的挑战"
	}
}

// ChallengeDifficulty 挑战难度
type ChallengeDifficulty int

const (
	ChallengeDifficultyEasy      ChallengeDifficulty = iota + 1 // 简单
	ChallengeDifficultyNormal                                    // 普通
	ChallengeDifficultyHard                                      // 困难
	ChallengeDifficultyExpert                                    // 专家
	ChallengeDifficultyLegendary                                 // 传奇
)

// String 返回难度字符串
func (cd ChallengeDifficulty) String() string {
	switch cd {
	case ChallengeDifficultyEasy:
		return "easy"
	case ChallengeDifficultyNormal:
		return "normal"
	case ChallengeDifficultyHard:
		return "hard"
	case ChallengeDifficultyExpert:
		return "expert"
	case ChallengeDifficultyLegendary:
		return "legendary"
	default:
		return "unknown"
	}
}

// IsValid 检查难度是否有效
func (cd ChallengeDifficulty) IsValid() bool {
	return cd >= ChallengeDifficultyEasy && cd <= ChallengeDifficultyLegendary
}

// GetMultiplier 获取难度倍数
func (cd ChallengeDifficulty) GetMultiplier() float64 {
	switch cd {
	case ChallengeDifficultyEasy:
		return 0.5
	case ChallengeDifficultyNormal:
		return 1.0
	case ChallengeDifficultyHard:
		return 1.5
	case ChallengeDifficultyExpert:
		return 2.0
	case ChallengeDifficultyLegendary:
		return 3.0
	default:
		return 1.0
	}
}

// GetRequiredLevel 获取所需等级
func (cd ChallengeDifficulty) GetRequiredLevel() int {
	switch cd {
	case ChallengeDifficultyEasy:
		return 1
	case ChallengeDifficultyNormal:
		return 5
	case ChallengeDifficultyHard:
		return 10
	case ChallengeDifficultyExpert:
		return 20
	case ChallengeDifficultyLegendary:
		return 50
	default:
		return 1
	}
}

// GetColor 获取难度颜色
func (cd ChallengeDifficulty) GetColor() string {
	switch cd {
	case ChallengeDifficultyEasy:
		return "green"
	case ChallengeDifficultyNormal:
		return "blue"
	case ChallengeDifficultyHard:
		return "yellow"
	case ChallengeDifficultyExpert:
		return "red"
	case ChallengeDifficultyLegendary:
		return "purple"
	default:
		return "gray"
	}
}

// ChallengeStatus 挑战状态
type ChallengeStatus int

const (
	ChallengeStatusAvailable  ChallengeStatus = iota + 1 // 可用
	ChallengeStatusInProgress                             // 进行中
	ChallengeStatusCompleted                              // 已完成
	ChallengeStatusFailed                                 // 失败
	ChallengeStatusLocked                                 // 锁定
	ChallengeStatusExpired                                // 过期
)

// String 返回状态字符串
func (cs ChallengeStatus) String() string {
	switch cs {
	case ChallengeStatusAvailable:
		return "available"
	case ChallengeStatusInProgress:
		return "in_progress"
	case ChallengeStatusCompleted:
		return "completed"
	case ChallengeStatusFailed:
		return "failed"
	case ChallengeStatusLocked:
		return "locked"
	case ChallengeStatusExpired:
		return "expired"
	default:
		return "unknown"
	}
}

// IsValid 检查状态是否有效
func (cs ChallengeStatus) IsValid() bool {
	return cs >= ChallengeStatusAvailable && cs <= ChallengeStatusExpired
}

// CanStart 检查是否可以开始
func (cs ChallengeStatus) CanStart() bool {
	return cs == ChallengeStatusAvailable
}

// IsFinished 检查是否已结束
func (cs ChallengeStatus) IsFinished() bool {
	return cs == ChallengeStatusCompleted || cs == ChallengeStatusFailed || cs == ChallengeStatusExpired
}

// BlessingType 祝福类型
type BlessingType int

const (
	BlessingTypeAttribute BlessingType = iota + 1 // 属性祝福
	BlessingTypeSkill                             // 技能祝福
	BlessingTypeExperience                        // 经验祝福
	BlessingTypeWealth                            // 财富祝福
	BlessingTypeProtection                        // 保护祝福
	BlessingTypeHealing                           // 治疗祝福
	BlessingTypeSpeed                             // 速度祝福
	BlessingTypeLuck                              // 幸运祝福
)

// String 返回类型字符串
func (bt BlessingType) String() string {
	switch bt {
	case BlessingTypeAttribute:
		return "attribute"
	case BlessingTypeSkill:
		return "skill"
	case BlessingTypeExperience:
		return "experience"
	case BlessingTypeWealth:
		return "wealth"
	case BlessingTypeProtection:
		return "protection"
	case BlessingTypeHealing:
		return "healing"
	case BlessingTypeSpeed:
		return "speed"
	case BlessingTypeLuck:
		return "luck"
	default:
		return "unknown"
	}
}

// IsValid 检查类型是否有效
func (bt BlessingType) IsValid() bool {
	return bt >= BlessingTypeAttribute && bt <= BlessingTypeLuck
}

// GetDescription 获取类型描述
func (bt BlessingType) GetDescription() string {
	switch bt {
	case BlessingTypeAttribute:
		return "提升角色基础属性的祝福"
	case BlessingTypeSkill:
		return "增强技能效果的祝福"
	case BlessingTypeExperience:
		return "增加经验获取的祝福"
	case BlessingTypeWealth:
		return "增加财富收入的祝福"
	case BlessingTypeProtection:
		return "提供保护效果的祝福"
	case BlessingTypeHealing:
		return "提供治疗效果的祝福"
	case BlessingTypeSpeed:
		return "提升移动和行动速度的祝福"
	case BlessingTypeLuck:
		return "增加幸运值的祝福"
	default:
		return "未知类型的祝福"
	}
}

// GetIcon 获取图标
func (bt BlessingType) GetIcon() string {
	switch bt {
	case BlessingTypeAttribute:
		return "💪"
	case BlessingTypeSkill:
		return "⚡"
	case BlessingTypeExperience:
		return "📚"
	case BlessingTypeWealth:
		return "💰"
	case BlessingTypeProtection:
		return "🛡️"
	case BlessingTypeHealing:
		return "❤️"
	case BlessingTypeSpeed:
		return "💨"
	case BlessingTypeLuck:
		return "🍀"
	default:
		return "❓"
	}
}

// BlessingStatus 祝福状态
type BlessingStatus int

const (
	BlessingStatusAvailable BlessingStatus = iota + 1 // 可用
	BlessingStatusActive                               // 激活
	BlessingStatusInactive                             // 未激活
	BlessingStatusExpired                              // 过期
	BlessingStatusLocked                               // 锁定
)

// String 返回状态字符串
func (bs BlessingStatus) String() string {
	switch bs {
	case BlessingStatusAvailable:
		return "available"
	case BlessingStatusActive:
		return "active"
	case BlessingStatusInactive:
		return "inactive"
	case BlessingStatusExpired:
		return "expired"
	case BlessingStatusLocked:
		return "locked"
	default:
		return "unknown"
	}
}

// IsValid 检查状态是否有效
func (bs BlessingStatus) IsValid() bool {
	return bs >= BlessingStatusAvailable && bs <= BlessingStatusLocked
}

// CanActivate 检查是否可以激活
func (bs BlessingStatus) CanActivate() bool {
	return bs == BlessingStatusAvailable
}

// IsActive 检查是否激活
func (bs BlessingStatus) IsActive() bool {
	return bs == BlessingStatusActive
}

// SacredRelic 圣物值对象
type SacredRelic struct {
	ID          string
	Name        string
	Description string
	Type        RelicType
	Rarity      RelicRarity
	Level       int
	Attributes  map[string]float64
	Effects     []string
	Requirements map[string]interface{}
	ObtainedAt  time.Time
}

// NewSacredRelic 创建圣物
func NewSacredRelic(id, name, description string, relicType RelicType, rarity RelicRarity) *SacredRelic {
	return &SacredRelic{
		ID:          id,
		Name:        name,
		Description: description,
		Type:        relicType,
		Rarity:      rarity,
		Level:       1,
		Attributes:  make(map[string]float64),
		Effects:     make([]string, 0),
		Requirements: make(map[string]interface{}),
		ObtainedAt:  time.Now(),
	}
}

// GetPower 获取圣物威力
func (sr *SacredRelic) GetPower() float64 {
	basePower := sr.Rarity.GetBasePower()
	levelMultiplier := float64(sr.Level)
	return basePower * levelMultiplier
}

// CanUpgrade 检查是否可以升级
func (sr *SacredRelic) CanUpgrade() bool {
	return sr.Level < sr.Rarity.GetMaxLevel()
}

// Upgrade 升级圣物
func (sr *SacredRelic) Upgrade() error {
	if !sr.CanUpgrade() {
		return fmt.Errorf("relic cannot be upgraded further")
	}
	
	sr.Level++
	// 升级时增强属性
	for attr, value := range sr.Attributes {
		sr.Attributes[attr] = value * 1.1 // 每级增加10%
	}
	
	return nil
}

// AddAttribute 添加属性
func (sr *SacredRelic) AddAttribute(name string, value float64) {
	sr.Attributes[name] = value
}

// AddEffect 添加效果
func (sr *SacredRelic) AddEffect(effect string) {
	sr.Effects = append(sr.Effects, effect)
}

// AddRequirement 添加需求
func (sr *SacredRelic) AddRequirement(name string, value interface{}) {
	sr.Requirements[name] = value
}

// CheckRequirements 检查需求
func (sr *SacredRelic) CheckRequirements(playerData map[string]interface{}) bool {
	for req, reqValue := range sr.Requirements {
		playerValue, exists := playerData[req]
		if !exists {
			return false
		}
		
		// 简单的数值比较
		if reqInt, ok := reqValue.(int); ok {
			if playerInt, ok := playerValue.(int); ok {
				if playerInt < reqInt {
					return false
				}
			}
		}
	}
	return true
}

// ToMap 转换为映射
func (sr *SacredRelic) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":           sr.ID,
		"name":         sr.Name,
		"description":  sr.Description,
		"type":         sr.Type.String(),
		"rarity":       sr.Rarity.String(),
		"level":        sr.Level,
		"power":        sr.GetPower(),
		"attributes":   sr.Attributes,
		"effects":      sr.Effects,
		"requirements": sr.Requirements,
		"obtained_at":  sr.ObtainedAt,
	}
}

// RelicType 圣物类型
type RelicType int

const (
	RelicTypeWeapon     RelicType = iota + 1 // 武器
	RelicTypeArmor                           // 护甲
	RelicTypeAccessory                       // 饰品
	RelicTypeConsumable                      // 消耗品
	RelicTypeSpecial                         // 特殊
)

// String 返回类型字符串
func (rt RelicType) String() string {
	switch rt {
	case RelicTypeWeapon:
		return "weapon"
	case RelicTypeArmor:
		return "armor"
	case RelicTypeAccessory:
		return "accessory"
	case RelicTypeConsumable:
		return "consumable"
	case RelicTypeSpecial:
		return "special"
	default:
		return "unknown"
	}
}

// RelicRarity 圣物稀有度
type RelicRarity int

const (
	RelicRarityCommon    RelicRarity = iota + 1 // 普通
	RelicRarityUncommon                          // 不常见
	RelicRarityRare                              // 稀有
	RelicRarityEpic                              // 史诗
	RelicRarityLegendary                         // 传奇
	RelicRarityMythic                            // 神话
)

// String 返回稀有度字符串
func (rr RelicRarity) String() string {
	switch rr {
	case RelicRarityCommon:
		return "common"
	case RelicRarityUncommon:
		return "uncommon"
	case RelicRarityRare:
		return "rare"
	case RelicRarityEpic:
		return "epic"
	case RelicRarityLegendary:
		return "legendary"
	case RelicRarityMythic:
		return "mythic"
	default:
		return "unknown"
	}
}

// GetBasePower 获取基础威力
func (rr RelicRarity) GetBasePower() float64 {
	switch rr {
	case RelicRarityCommon:
		return 10.0
	case RelicRarityUncommon:
		return 25.0
	case RelicRarityRare:
		return 50.0
	case RelicRarityEpic:
		return 100.0
	case RelicRarityLegendary:
		return 200.0
	case RelicRarityMythic:
		return 500.0
	default:
		return 1.0
	}
}

// GetMaxLevel 获取最大等级
func (rr RelicRarity) GetMaxLevel() int {
	switch rr {
	case RelicRarityCommon:
		return 10
	case RelicRarityUncommon:
		return 20
	case RelicRarityRare:
		return 30
	case RelicRarityEpic:
		return 50
	case RelicRarityLegendary:
		return 80
	case RelicRarityMythic:
		return 100
	default:
		return 1
	}
}

// GetColor 获取颜色
func (rr RelicRarity) GetColor() string {
	switch rr {
	case RelicRarityCommon:
		return "gray"
	case RelicRarityUncommon:
		return "green"
	case RelicRarityRare:
		return "blue"
	case RelicRarityEpic:
		return "purple"
	case RelicRarityLegendary:
		return "orange"
	case RelicRarityMythic:
		return "red"
	default:
		return "white"
	}
}

// SacredPortal 圣地传送门值对象
type SacredPortal struct {
	ID            string
	Name          string
	Destination   string
	RequiredLevel int
	Cost          int
	Cooldown      time.Duration
	LastUsed      time.Time
	Active        bool
}

// NewSacredPortal 创建传送门
func NewSacredPortal(id, name, destination string, requiredLevel, cost int, cooldown time.Duration) *SacredPortal {
	return &SacredPortal{
		ID:            id,
		Name:          name,
		Destination:   destination,
		RequiredLevel: requiredLevel,
		Cost:          cost,
		Cooldown:      cooldown,
		Active:        true,
	}
}

// CanUse 检查是否可以使用
func (sp *SacredPortal) CanUse(playerLevel int, playerGold int) bool {
	if !sp.Active {
		return false
	}
	
	if playerLevel < sp.RequiredLevel {
		return false
	}
	
	if playerGold < sp.Cost {
		return false
	}
	
	// 检查冷却时间
	if !sp.LastUsed.IsZero() && time.Since(sp.LastUsed) < sp.Cooldown {
		return false
	}
	
	return true
}

// Use 使用传送门
func (sp *SacredPortal) Use() error {
	if !sp.Active {
		return fmt.Errorf("portal is not active")
	}
	
	sp.LastUsed = time.Now()
	return nil
}

// GetRemainingCooldown 获取剩余冷却时间
func (sp *SacredPortal) GetRemainingCooldown() time.Duration {
	if sp.LastUsed.IsZero() {
		return 0
	}
	
	elapsed := time.Since(sp.LastUsed)
	if elapsed >= sp.Cooldown {
		return 0
	}
	
	return sp.Cooldown - elapsed
}

// Activate 激活传送门
func (sp *SacredPortal) Activate() {
	sp.Active = true
}

// Deactivate 停用传送门
func (sp *SacredPortal) Deactivate() {
	sp.Active = false
}

// ToMap 转换为映射
func (sp *SacredPortal) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":             sp.ID,
		"name":           sp.Name,
		"destination":    sp.Destination,
		"required_level": sp.RequiredLevel,
		"cost":           sp.Cost,
		"cooldown":       sp.Cooldown.String(),
		"last_used":      sp.LastUsed,
		"active":         sp.Active,
		"remaining_cooldown": sp.GetRemainingCooldown().String(),
	}
}

// SacredAura 圣地光环值对象
type SacredAura struct {
	Type       AuraType
	Intensity  float64
	Radius     float64
	Effects    map[string]float64
	Duration   time.Duration
	ActivatedAt time.Time
}

// NewSacredAura 创建圣地光环
func NewSacredAura(auraType AuraType, intensity, radius float64, duration time.Duration) *SacredAura {
	return &SacredAura{
		Type:      auraType,
		Intensity: intensity,
		Radius:    radius,
		Effects:   make(map[string]float64),
		Duration:  duration,
		ActivatedAt: time.Now(),
	}
}

// IsActive 检查是否激活
func (sa *SacredAura) IsActive() bool {
	return time.Since(sa.ActivatedAt) < sa.Duration
}

// GetRemainingDuration 获取剩余时间
func (sa *SacredAura) GetRemainingDuration() time.Duration {
	if !sa.IsActive() {
		return 0
	}
	return sa.Duration - time.Since(sa.ActivatedAt)
}

// AddEffect 添加效果
func (sa *SacredAura) AddEffect(name string, value float64) {
	sa.Effects[name] = value
}

// GetEffect 获取效果值
func (sa *SacredAura) GetEffect(name string) float64 {
	return sa.Effects[name] * sa.Intensity
}

// AuraType 光环类型
type AuraType int

const (
	AuraTypeHealing     AuraType = iota + 1 // 治疗光环
	AuraTypeProtection                      // 保护光环
	AuraTypeStrength                        // 力量光环
	AuraTypeWisdom                          // 智慧光环
	AuraTypeSpeed                           // 速度光环
	AuraTypeLuck                            // 幸运光环
)

// String 返回光环类型字符串
func (at AuraType) String() string {
	switch at {
	case AuraTypeHealing:
		return "healing"
	case AuraTypeProtection:
		return "protection"
	case AuraTypeStrength:
		return "strength"
	case AuraTypeWisdom:
		return "wisdom"
	case AuraTypeSpeed:
		return "speed"
	case AuraTypeLuck:
		return "luck"
	default:
		return "unknown"
	}
}