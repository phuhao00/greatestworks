package synthesis

// RecipeCategory 配方分类
type RecipeCategory int

const (
	RecipeCategoryWeapon RecipeCategory = iota + 1
	RecipeCategoryArmor
	RecipeCategoryAccessory
	RecipeCategoryConsumable
	RecipeCategoryMaterial
	RecipeCategorySpecial
	RecipeCategoryEnchant
	RecipeCategoryGem
)

// String 返回配方分类字符串
func (rc RecipeCategory) String() string {
	switch rc {
	case RecipeCategoryWeapon:
		return "weapon"
	case RecipeCategoryArmor:
		return "armor"
	case RecipeCategoryAccessory:
		return "accessory"
	case RecipeCategoryConsumable:
		return "consumable"
	case RecipeCategoryMaterial:
		return "material"
	case RecipeCategorySpecial:
		return "special"
	case RecipeCategoryEnchant:
		return "enchant"
	case RecipeCategoryGem:
		return "gem"
	default:
		return "unknown"
	}
}

// MaterialType 材料类型
type MaterialType int

const (
	MaterialTypeOre MaterialType = iota + 1
	MaterialTypeHerb
	MaterialTypeLeather
	MaterialTypeCloth
	MaterialTypeWood
	MaterialTypeStone
	MaterialTypeMetal
	MaterialTypeGem
	MaterialTypeMagic
	MaterialTypeRare
)

// String 返回材料类型字符串
func (mt MaterialType) String() string {
	switch mt {
	case MaterialTypeOre:
		return "ore"
	case MaterialTypeHerb:
		return "herb"
	case MaterialTypeLeather:
		return "leather"
	case MaterialTypeCloth:
		return "cloth"
	case MaterialTypeWood:
		return "wood"
	case MaterialTypeStone:
		return "stone"
	case MaterialTypeMetal:
		return "metal"
	case MaterialTypeGem:
		return "gem"
	case MaterialTypeMagic:
		return "magic"
	case MaterialTypeRare:
		return "rare"
	default:
		return "unknown"
	}
}

// Quality 品质
type Quality int

const (
	QualityCommon Quality = iota + 1
	QualityUncommon
	QualityRare
	QualityEpic
	QualityLegendary
	QualityMythic
)

// String 返回品质字符串
func (q Quality) String() string {
	switch q {
	case QualityCommon:
		return "common"
	case QualityUncommon:
		return "uncommon"
	case QualityRare:
		return "rare"
	case QualityEpic:
		return "epic"
	case QualityLegendary:
		return "legendary"
	case QualityMythic:
		return "mythic"
	default:
		return "unknown"
	}
}

// GetQualityMultiplier 获取品质倍数
func (q Quality) GetQualityMultiplier() float64 {
	switch q {
	case QualityCommon:
		return 1.0
	case QualityUncommon:
		return 1.2
	case QualityRare:
		return 1.5
	case QualityEpic:
		return 2.0
	case QualityLegendary:
		return 3.0
	case QualityMythic:
		return 5.0
	default:
		return 1.0
	}
}

// MaterialRequirement 材料需求值对象
type MaterialRequirement struct {
	MaterialID string `json:"material_id"`
	Quantity   int    `json:"quantity"`
}

// NewMaterialRequirement 创建材料需求
func NewMaterialRequirement(materialID string, quantity int) *MaterialRequirement {
	return &MaterialRequirement{
		MaterialID: materialID,
		Quantity:   quantity,
	}
}

// GetMaterialID 获取材料ID
func (mr *MaterialRequirement) GetMaterialID() string {
	return mr.MaterialID
}

// GetQuantity 获取数量
func (mr *MaterialRequirement) GetQuantity() int {
	return mr.Quantity
}

// ItemOutput 物品产出值对象
type ItemOutput struct {
	ItemID      string  `json:"item_id"`
	Quantity    int     `json:"quantity"`
	Probability float64 `json:"probability"` // 产出概率 0-1
}

// NewItemOutput 创建物品产出
func NewItemOutput(itemID string, quantity int, probability float64) *ItemOutput {
	return &ItemOutput{
		ItemID:      itemID,
		Quantity:    quantity,
		Probability: probability,
	}
}

// GetItemID 获取物品ID
func (io *ItemOutput) GetItemID() string {
	return io.ItemID
}

// GetQuantity 获取数量
func (io *ItemOutput) GetQuantity() int {
	return io.Quantity
}

// GetProbability 获取概率
func (io *ItemOutput) GetProbability() float64 {
	return io.Probability
}

// SynthesisBonus 合成加成值对象
type SynthesisBonus struct {
	bonusType   BonusType
	bonusValue  float64
	duration    int // 持续时间（秒），0表示永久
	description string
}

// BonusType 加成类型
type BonusType int

const (
	BonusTypeSuccessRate BonusType = iota + 1
	BonusTypeCraftSpeed
	BonusTypeMaterialSave
	BonusTypeQualityUp
	BonusTypeExtraOutput
)

// String 返回加成类型字符串
func (bt BonusType) String() string {
	switch bt {
	case BonusTypeSuccessRate:
		return "success_rate"
	case BonusTypeCraftSpeed:
		return "craft_speed"
	case BonusTypeMaterialSave:
		return "material_save"
	case BonusTypeQualityUp:
		return "quality_up"
	case BonusTypeExtraOutput:
		return "extra_output"
	default:
		return "unknown"
	}
}

// NewSynthesisBonus 创建合成加成
func NewSynthesisBonus(bonusType BonusType, bonusValue float64, duration int, description string) *SynthesisBonus {
	return &SynthesisBonus{
		bonusType:   bonusType,
		bonusValue:  bonusValue,
		duration:    duration,
		description: description,
	}
}

// GetBonusType 获取加成类型
func (sb *SynthesisBonus) GetBonusType() BonusType {
	return sb.bonusType
}

// GetBonusValue 获取加成值
func (sb *SynthesisBonus) GetBonusValue() float64 {
	return sb.bonusValue
}

// GetDuration 获取持续时间
func (sb *SynthesisBonus) GetDuration() int {
	return sb.duration
}

// GetDescription 获取描述
func (sb *SynthesisBonus) GetDescription() string {
	return sb.description
}

// IsPermanent 是否永久
func (sb *SynthesisBonus) IsPermanent() bool {
	return sb.duration == 0
}

// CraftingCondition 制作条件值对象
type CraftingCondition struct {
	conditionType ConditionType
	value         interface{}
	description   string
}

// ConditionType 条件类型
type ConditionType int

const (
	ConditionTypeLevel ConditionType = iota + 1
	ConditionTypeSkill
	ConditionTypeItem
	ConditionTypeQuest
	ConditionTypeAchievement
	ConditionTypeTime
	ConditionTypeLocation
)

// String 返回条件类型字符串
func (ct ConditionType) String() string {
	switch ct {
	case ConditionTypeLevel:
		return "level"
	case ConditionTypeSkill:
		return "skill"
	case ConditionTypeItem:
		return "item"
	case ConditionTypeQuest:
		return "quest"
	case ConditionTypeAchievement:
		return "achievement"
	case ConditionTypeTime:
		return "time"
	case ConditionTypeLocation:
		return "location"
	default:
		return "unknown"
	}
}

// NewCraftingCondition 创建制作条件
func NewCraftingCondition(conditionType ConditionType, value interface{}, description string) *CraftingCondition {
	return &CraftingCondition{
		conditionType: conditionType,
		value:         value,
		description:   description,
	}
}

// GetConditionType 获取条件类型
func (cc *CraftingCondition) GetConditionType() ConditionType {
	return cc.conditionType
}

// GetValue 获取条件值
func (cc *CraftingCondition) GetValue() interface{} {
	return cc.value
}

// GetDescription 获取描述
func (cc *CraftingCondition) GetDescription() string {
	return cc.description
}
