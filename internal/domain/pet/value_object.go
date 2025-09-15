package pet

import (
	"fmt"
	"math/rand"
	"time"
)

// PetCategory 宠物类别
type PetCategory int

const (
	PetCategoryFire   PetCategory = 1 // 火系
	PetCategoryWater  PetCategory = 2 // 水系
	PetCategoryEarth  PetCategory = 3 // 土系
	PetCategoryAir    PetCategory = 4 // 风系
	PetCategoryLight  PetCategory = 5 // 光系
	PetCategoryDark   PetCategory = 6 // 暗系
	PetCategoryNormal PetCategory = 7 // 普通系
)

// String 返回宠物类别字符串
func (pc PetCategory) String() string {
	switch pc {
	case PetCategoryFire:
		return "Fire"
	case PetCategoryWater:
		return "Water"
	case PetCategoryEarth:
		return "Earth"
	case PetCategoryAir:
		return "Air"
	case PetCategoryLight:
		return "Light"
	case PetCategoryDark:
		return "Dark"
	case PetCategoryNormal:
		return "Normal"
	default:
		return "Unknown"
	}
}

// IsValid 检查宠物类别是否有效
func (pc PetCategory) IsValid() bool {
	return pc >= PetCategoryFire && pc <= PetCategoryNormal
}

// PetState 宠物状态
type PetState int

const (
	PetStateIdle     PetState = 1 // 空闲
	PetStateBattle   PetState = 2 // 战斗中
	PetStateTraining PetState = 3 // 训练中
	PetStateDead     PetState = 4 // 死亡
)

// String 返回宠物状态字符串
func (ps PetState) String() string {
	switch ps {
	case PetStateIdle:
		return "Idle"
	case PetStateBattle:
		return "Battle"
	case PetStateTraining:
		return "Training"
	case PetStateDead:
		return "Dead"
	default:
		return "Unknown"
	}
}

// IsValid 检查宠物状态是否有效
func (ps PetState) IsValid() bool {
	return ps >= PetStateIdle && ps <= PetStateDead
}

// PetAttributes 宠物属性
type PetAttributes struct {
	health   int64
	attack   int64
	defense  int64
	speed    int64
	critical int64
	hit      int64
	dodge    int64
}

// NewPetAttributes 创建新的宠物属性
func NewPetAttributes() *PetAttributes {
	return &PetAttributes{
		health:   100,
		attack:   50,
		defense:  30,
		speed:    40,
		critical: 10,
		hit:      80,
		dodge:    20,
	}
}

// GetHealth 获取生命值
func (pa *PetAttributes) GetHealth() int64 {
	return pa.health
}

// GetAttack 获取攻击力
func (pa *PetAttributes) GetAttack() int64 {
	return pa.attack
}

// GetDefense 获取防御力
func (pa *PetAttributes) GetDefense() int64 {
	return pa.defense
}

// GetSpeed 获取速度
func (pa *PetAttributes) GetSpeed() int64 {
	return pa.speed
}

// GetCritical 获取暴击
func (pa *PetAttributes) GetCritical() int64 {
	return pa.critical
}

// GetHit 获取命中
func (pa *PetAttributes) GetHit() int64 {
	return pa.hit
}

// GetDodge 获取闪避
func (pa *PetAttributes) GetDodge() int64 {
	return pa.dodge
}

// AddHealth 增加生命值
func (pa *PetAttributes) AddHealth(value int64) {
	pa.health += value
	if pa.health < 0 {
		pa.health = 0
	}
}

// AddAttack 增加攻击力
func (pa *PetAttributes) AddAttack(value int64) {
	pa.attack += value
	if pa.attack < 0 {
		pa.attack = 0
	}
}

// AddDefense 增加防御力
func (pa *PetAttributes) AddDefense(value int64) {
	pa.defense += value
	if pa.defense < 0 {
		pa.defense = 0
	}
}

// AddSpeed 增加速度
func (pa *PetAttributes) AddSpeed(value int64) {
	pa.speed += value
	if pa.speed < 0 {
		pa.speed = 0
	}
}

// AddCritical 增加暴击
func (pa *PetAttributes) AddCritical(value int64) {
	pa.critical += value
	if pa.critical < 0 {
		pa.critical = 0
	}
}

// AddHit 增加命中
func (pa *PetAttributes) AddHit(value int64) {
	pa.hit += value
	if pa.hit < 0 {
		pa.hit = 0
	}
}

// AddDodge 增加闪避
func (pa *PetAttributes) AddDodge(value int64) {
	pa.dodge += value
	if pa.dodge < 0 {
		pa.dodge = 0
	}
}

// UpgradeOnLevelUp 升级时属性提升
func (pa *PetAttributes) UpgradeOnLevelUp(level uint32) {
	// 每级提升基础属性
	levelBonus := int64(level)
	pa.health += levelBonus * 10
	pa.attack += levelBonus * 5
	pa.defense += levelBonus * 3
	pa.speed += levelBonus * 2
	pa.critical += levelBonus * 1
	pa.hit += levelBonus * 1
	pa.dodge += levelBonus * 1
}

// UpgradeOnStarUp 升星时属性提升
func (pa *PetAttributes) UpgradeOnStarUp(star uint32) {
	// 每星提升大量属性
	starBonus := int64(star)
	pa.health += starBonus * 50
	pa.attack += starBonus * 25
	pa.defense += starBonus * 15
	pa.speed += starBonus * 10
	pa.critical += starBonus * 5
	pa.hit += starBonus * 5
	pa.dodge += starBonus * 5
}

// AddRandomAttribute 随机增加属性
func (pa *PetAttributes) AddRandomAttribute(value int64) {
	attributeType := rand.Intn(7)
	switch attributeType {
	case 0:
		pa.AddHealth(value)
	case 1:
		pa.AddAttack(value)
	case 2:
		pa.AddDefense(value)
	case 3:
		pa.AddSpeed(value)
	case 4:
		pa.AddCritical(value)
	case 5:
		pa.AddHit(value)
	case 6:
		pa.AddDodge(value)
	}
}

// CalculatePower 计算战力
func (pa *PetAttributes) CalculatePower() int64 {
	return pa.health*2 + pa.attack*3 + pa.defense*2 + pa.speed*1 + pa.critical*2 + pa.hit*1 + pa.dodge*1
}

// Clone 克隆属性
func (pa *PetAttributes) Clone() *PetAttributes {
	return &PetAttributes{
		health:   pa.health,
		attack:   pa.attack,
		defense:  pa.defense,
		speed:    pa.speed,
		critical: pa.critical,
		hit:      pa.hit,
		dodge:    pa.dodge,
	}
}

// FoodType 食物类型
type FoodType int

const (
	FoodTypeExperience FoodType = 1 // 经验食物
	FoodTypeHealth     FoodType = 2 // 生命食物
	FoodTypeAttack     FoodType = 3 // 攻击食物
	FoodTypeDefense    FoodType = 4 // 防御食物
)

// String 返回食物类型字符串
func (ft FoodType) String() string {
	switch ft {
	case FoodTypeExperience:
		return "Experience"
	case FoodTypeHealth:
		return "Health"
	case FoodTypeAttack:
		return "Attack"
	case FoodTypeDefense:
		return "Defense"
	default:
		return "Unknown"
	}
}

// IsValid 检查食物类型是否有效
func (ft FoodType) IsValid() bool {
	return ft >= FoodTypeExperience && ft <= FoodTypeDefense
}

// TrainingType 训练类型
type TrainingType int

const (
	TrainingTypeExperience TrainingType = 1 // 经验训练
	TrainingTypeAttribute  TrainingType = 2 // 属性训练
	TrainingTypeSkill      TrainingType = 3 // 技能训练
)

// String 返回训练类型字符串
func (tt TrainingType) String() string {
	switch tt {
	case TrainingTypeExperience:
		return "Experience"
	case TrainingTypeAttribute:
		return "Attribute"
	case TrainingTypeSkill:
		return "Skill"
	default:
		return "Unknown"
	}
}

// IsValid 检查训练类型是否有效
func (tt TrainingType) IsValid() bool {
	return tt >= TrainingTypeExperience && tt <= TrainingTypeSkill
}

// GetTrainingDuration 获取训练持续时间
func (tt TrainingType) GetTrainingDuration() time.Duration {
	switch tt {
	case TrainingTypeExperience:
		return 30 * time.Minute
	case TrainingTypeAttribute:
		return 60 * time.Minute
	case TrainingTypeSkill:
		return 120 * time.Minute
	default:
		return 30 * time.Minute
	}
}

// GetTrainingCost 获取训练消耗
func (tt TrainingType) GetTrainingCost() int64 {
	switch tt {
	case TrainingTypeExperience:
		return 100
	case TrainingTypeAttribute:
		return 200
	case TrainingTypeSkill:
		return 500
	default:
		return 100
	}
}

// PetRarity 宠物稀有度
type PetRarity int

const (
	PetRarityCommon    PetRarity = 1 // 普通
	PetRarityUncommon  PetRarity = 2 // 不常见
	PetRarityRare      PetRarity = 3 // 稀有
	PetRarityEpic      PetRarity = 4 // 史诗
	PetRarityLegendary PetRarity = 5 // 传说
	PetRarityMythic    PetRarity = 6 // 神话
)

// String 返回稀有度字符串
func (pr PetRarity) String() string {
	switch pr {
	case PetRarityCommon:
		return "Common"
	case PetRarityUncommon:
		return "Uncommon"
	case PetRarityRare:
		return "Rare"
	case PetRarityEpic:
		return "Epic"
	case PetRarityLegendary:
		return "Legendary"
	case PetRarityMythic:
		return "Mythic"
	default:
		return "Unknown"
	}
}

// IsValid 检查稀有度是否有效
func (pr PetRarity) IsValid() bool {
	return pr >= PetRarityCommon && pr <= PetRarityMythic
}

// GetAttributeMultiplier 获取属性倍数
func (pr PetRarity) GetAttributeMultiplier() float64 {
	switch pr {
	case PetRarityCommon:
		return 1.0
	case PetRarityUncommon:
		return 1.2
	case PetRarityRare:
		return 1.5
	case PetRarityEpic:
		return 2.0
	case PetRarityLegendary:
		return 3.0
	case PetRarityMythic:
		return 5.0
	default:
		return 1.0
	}
}

// PetGender 宠物性别
type PetGender int

const (
	PetGenderMale   PetGender = 1 // 雄性
	PetGenderFemale PetGender = 2 // 雌性
	PetGenderNone   PetGender = 3 // 无性别
)

// String 返回性别字符串
func (pg PetGender) String() string {
	switch pg {
	case PetGenderMale:
		return "Male"
	case PetGenderFemale:
		return "Female"
	case PetGenderNone:
		return "None"
	default:
		return "Unknown"
	}
}

// IsValid 检查性别是否有效
func (pg PetGender) IsValid() bool {
	return pg >= PetGenderMale && pg <= PetGenderNone
}

// PetSize 宠物体型
type PetSize int

const (
	PetSizeTiny   PetSize = 1 // 微型
	PetSizeSmall  PetSize = 2 // 小型
	PetSizeMedium PetSize = 3 // 中型
	PetSizeLarge  PetSize = 4 // 大型
	PetSizeHuge   PetSize = 5 // 巨型
)

// String 返回体型字符串
func (ps PetSize) String() string {
	switch ps {
	case PetSizeTiny:
		return "Tiny"
	case PetSizeSmall:
		return "Small"
	case PetSizeMedium:
		return "Medium"
	case PetSizeLarge:
		return "Large"
	case PetSizeHuge:
		return "Huge"
	default:
		return "Unknown"
	}
}

// IsValid 检查体型是否有效
func (ps PetSize) IsValid() bool {
	return ps >= PetSizeTiny && ps <= PetSizeHuge
}

// GetSpeedModifier 获取速度修正
func (ps PetSize) GetSpeedModifier() float64 {
	switch ps {
	case PetSizeTiny:
		return 1.3
	case PetSizeSmall:
		return 1.1
	case PetSizeMedium:
		return 1.0
	case PetSizeLarge:
		return 0.9
	case PetSizeHuge:
		return 0.7
	default:
		return 1.0
	}
}

// GetHealthModifier 获取生命值修正
func (ps PetSize) GetHealthModifier() float64 {
	switch ps {
	case PetSizeTiny:
		return 0.7
	case PetSizeSmall:
		return 0.9
	case PetSizeMedium:
		return 1.0
	case PetSizeLarge:
		return 1.2
	case PetSizeHuge:
		return 1.5
	default:
		return 1.0
	}
}

// PetPersonality 宠物性格
type PetPersonality int

const (
	PetPersonalityBrave     PetPersonality = 1  // 勇敢
	PetPersonalityTimid     PetPersonality = 2  // 胆小
	PetPersonalityAggressive PetPersonality = 3  // 好斗
	PetPersonalityGentle    PetPersonality = 4  // 温和
	PetPersonalityPlayful   PetPersonality = 5  // 顽皮
	PetPersonalityLazy      PetPersonality = 6  // 懒惰
	PetPersonalityLoyal     PetPersonality = 7  // 忠诚
	PetPersonalityStubborn  PetPersonality = 8  // 固执
	PetPersonalityCurious   PetPersonality = 9  // 好奇
	PetPersonalityCalm      PetPersonality = 10 // 冷静
)

// String 返回性格字符串
func (pp PetPersonality) String() string {
	switch pp {
	case PetPersonalityBrave:
		return "Brave"
	case PetPersonalityTimid:
		return "Timid"
	case PetPersonalityAggressive:
		return "Aggressive"
	case PetPersonalityGentle:
		return "Gentle"
	case PetPersonalityPlayful:
		return "Playful"
	case PetPersonalityLazy:
		return "Lazy"
	case PetPersonalityLoyal:
		return "Loyal"
	case PetPersonalityStubborn:
		return "Stubborn"
	case PetPersonalityCurious:
		return "Curious"
	case PetPersonalityCalm:
		return "Calm"
	default:
		return "Unknown"
	}
}

// IsValid 检查性格是否有效
func (pp PetPersonality) IsValid() bool {
	return pp >= PetPersonalityBrave && pp <= PetPersonalityCalm
}

// GetAttributeBonus 获取性格属性加成
func (pp PetPersonality) GetAttributeBonus() map[string]float64 {
	switch pp {
	case PetPersonalityBrave:
		return map[string]float64{"attack": 1.1, "defense": 0.9}
	case PetPersonalityTimid:
		return map[string]float64{"dodge": 1.2, "attack": 0.8}
	case PetPersonalityAggressive:
		return map[string]float64{"attack": 1.2, "critical": 1.1, "defense": 0.8}
	case PetPersonalityGentle:
		return map[string]float64{"health": 1.1, "attack": 0.9}
	case PetPersonalityPlayful:
		return map[string]float64{"speed": 1.2, "dodge": 1.1}
	case PetPersonalityLazy:
		return map[string]float64{"health": 1.2, "speed": 0.7}
	case PetPersonalityLoyal:
		return map[string]float64{"defense": 1.2, "health": 1.1}
	case PetPersonalityStubborn:
		return map[string]float64{"defense": 1.3, "speed": 0.8}
	case PetPersonalityCurious:
		return map[string]float64{"hit": 1.2, "critical": 1.1}
	case PetPersonalityCalm:
		return map[string]float64{"hit": 1.1, "dodge": 1.1}
	default:
		return map[string]float64{}
	}
}

// PetMood 宠物心情
type PetMood int

const (
	PetMoodHappy   PetMood = 1 // 开心
	PetMoodNormal  PetMood = 2 // 普通
	PetMoodSad     PetMood = 3 // 伤心
	PetMoodAngry   PetMood = 4 // 愤怒
	PetMoodExcited PetMood = 5 // 兴奋
	PetMoodTired   PetMood = 6 // 疲惫
)

// String 返回心情字符串
func (pm PetMood) String() string {
	switch pm {
	case PetMoodHappy:
		return "Happy"
	case PetMoodNormal:
		return "Normal"
	case PetMoodSad:
		return "Sad"
	case PetMoodAngry:
		return "Angry"
	case PetMoodExcited:
		return "Excited"
	case PetMoodTired:
		return "Tired"
	default:
		return "Unknown"
	}
}

// IsValid 检查心情是否有效
func (pm PetMood) IsValid() bool {
	return pm >= PetMoodHappy && pm <= PetMoodTired
}

// GetEfficiencyModifier 获取效率修正
func (pm PetMood) GetEfficiencyModifier() float64 {
	switch pm {
	case PetMoodHappy:
		return 1.2
	case PetMoodNormal:
		return 1.0
	case PetMoodSad:
		return 0.8
	case PetMoodAngry:
		return 0.9
	case PetMoodExcited:
		return 1.3
	case PetMoodTired:
		return 0.7
	default:
		return 1.0
	}
}

// GetExperienceModifier 获取经验修正
func (pm PetMood) GetExperienceModifier() float64 {
	switch pm {
	case PetMoodHappy:
		return 1.1
	case PetMoodNormal:
		return 1.0
	case PetMoodSad:
		return 0.9
	case PetMoodAngry:
		return 0.8
	case PetMoodExcited:
		return 1.2
	case PetMoodTired:
		return 0.8
	default:
		return 1.0
	}
}