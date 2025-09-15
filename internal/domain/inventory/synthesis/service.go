package synthesis

import (
	"math/rand"
	"time"
)

// SynthesisService 合成领域服务
type SynthesisService struct {
	recipeFactory   *RecipeFactory
	materialFactory *MaterialFactory
	bonusCalculator *BonusCalculator
}

// NewSynthesisService 创建合成服务
func NewSynthesisService() *SynthesisService {
	return &SynthesisService{
		recipeFactory:   NewRecipeFactory(),
		materialFactory: NewMaterialFactory(),
		bonusCalculator: NewBonusCalculator(),
	}
}

// ValidateRecipe 验证配方
func (ss *SynthesisService) ValidateRecipe(recipe *Recipe) error {
	if recipe == nil {
		return ErrInvalidRecipe
	}
	
	if len(recipe.GetRequirements()) == 0 {
		return ErrInvalidRecipe
	}
	
	if len(recipe.GetOutputs()) == 0 {
		return ErrInvalidRecipe
	}
	
	if recipe.GetSuccessRate() < 0 || recipe.GetSuccessRate() > 1 {
		return ErrInvalidRecipe
	}
	
	return nil
}

// CalculateEnhancedSuccessRate 计算增强成功率
func (ss *SynthesisService) CalculateEnhancedSuccessRate(baseRate float64, bonuses []*SynthesisBonus, playerLevel int) float64 {
	enhancedRate := baseRate
	
	// 应用加成
	for _, bonus := range bonuses {
		if bonus.GetBonusType() == BonusTypeSuccessRate {
			enhancedRate += bonus.GetBonusValue()
		}
	}
	
	// 玩家等级加成
	levelBonus := float64(playerLevel) * 0.001 // 每级0.1%加成
	enhancedRate += levelBonus
	
	// 确保在合理范围内
	if enhancedRate > 1.0 {
		enhancedRate = 1.0
	} else if enhancedRate < 0.0 {
		enhancedRate = 0.0
	}
	
	return enhancedRate
}

// CalculateCraftTime 计算制作时间
func (ss *SynthesisService) CalculateCraftTime(baseCraftTime time.Duration, bonuses []*SynthesisBonus) time.Duration {
	speedMultiplier := 1.0
	
	// 应用速度加成
	for _, bonus := range bonuses {
		if bonus.GetBonusType() == BonusTypeCraftSpeed {
			speedMultiplier += bonus.GetBonusValue()
		}
	}
	
	// 计算最终时间
	finalTime := time.Duration(float64(baseCraftTime) / speedMultiplier)
	
	// 最小时间限制
	minTime := time.Second * 1
	if finalTime < minTime {
		finalTime = minTime
	}
	
	return finalTime
}

// CalculateMaterialConsumption 计算材料消耗
func (ss *SynthesisService) CalculateMaterialConsumption(requirements []*MaterialRequirement, bonuses []*SynthesisBonus) []*MaterialRequirement {
	materialSaveRate := 0.0
	
	// 计算材料节省率
	for _, bonus := range bonuses {
		if bonus.GetBonusType() == BonusTypeMaterialSave {
			materialSaveRate += bonus.GetBonusValue()
		}
	}
	
	// 应用材料节省
	adjustedRequirements := make([]*MaterialRequirement, len(requirements))
	for i, req := range requirements {
		adjustedQuantity := int(float64(req.GetQuantity()) * (1.0 - materialSaveRate))
		if adjustedQuantity < 1 {
			adjustedQuantity = 1 // 至少需要1个
		}
		adjustedRequirements[i] = NewMaterialRequirement(req.GetMaterialID(), adjustedQuantity)
	}
	
	return adjustedRequirements
}

// GenerateRandomMaterial 生成随机材料
func (ss *SynthesisService) GenerateRandomMaterial(materialType MaterialType, playerLevel int) *Material {
	return ss.materialFactory.CreateRandomMaterial(materialType, playerLevel)
}

// CreateRecipeFromTemplate 从模板创建配方
func (ss *SynthesisService) CreateRecipeFromTemplate(templateID string, playerLevel int) *Recipe {
	return ss.recipeFactory.CreateFromTemplate(templateID, playerLevel)
}

// RecipeFactory 配方工厂
type RecipeFactory struct {
	random *rand.Rand
}

// NewRecipeFactory 创建配方工厂
func NewRecipeFactory() *RecipeFactory {
	return &RecipeFactory{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CreateFromTemplate 从模板创建配方
func (rf *RecipeFactory) CreateFromTemplate(templateID string, playerLevel int) *Recipe {
	// 这里应该从配置或数据库加载模板
	// 为了演示，创建一个示例配方
	recipe := NewRecipe("示例配方", RecipeCategoryWeapon, 0.8)
	recipe.AddRequirement("iron_ore", 5)
	recipe.AddRequirement("coal", 2)
	recipe.AddOutput("iron_sword", 1, 1.0)
	recipe.SetRequireLevel(playerLevel)
	
	return recipe
}

// CreateRandomRecipe 创建随机配方
func (rf *RecipeFactory) CreateRandomRecipe(category RecipeCategory, playerLevel int) *Recipe {
	successRate := 0.5 + rf.random.Float64()*0.4 // 50%-90%成功率
	recipe := NewRecipe(rf.generateRecipeName(category), category, successRate)
	
	// 添加随机材料需求
	materialCount := rf.random.Intn(3) + 2 // 2-4种材料
	for i := 0; i < materialCount; i++ {
		materialID := rf.generateMaterialID(category)
		quantity := rf.random.Intn(5) + 1
		recipe.AddRequirement(materialID, quantity)
	}
	
	// 添加产出
	outputID := rf.generateOutputID(category)
	recipe.AddOutput(outputID, 1, 1.0)
	
	recipe.SetRequireLevel(playerLevel)
	return recipe
}

// generateRecipeName 生成配方名称
func (rf *RecipeFactory) generateRecipeName(category RecipeCategory) string {
	names := map[RecipeCategory][]string{
		RecipeCategoryWeapon:     {"铁剑制作", "钢刀锻造", "魔法杖合成"},
		RecipeCategoryArmor:      {"皮甲制作", "铁甲锻造", "法袍缝制"},
		RecipeCategoryAccessory:  {"戒指打造", "项链制作", "护符合成"},
		RecipeCategoryConsumable: {"生命药水", "魔法药水", "解毒剂"},
	}
	
	nameList := names[category]
	if len(nameList) == 0 {
		return "未知配方"
	}
	
	return nameList[rf.random.Intn(len(nameList))]
}

// generateMaterialID 生成材料ID
func (rf *RecipeFactory) generateMaterialID(category RecipeCategory) string {
	materials := map[RecipeCategory][]string{
		RecipeCategoryWeapon:     {"iron_ore", "coal", "leather"},
		RecipeCategoryArmor:      {"leather", "cloth", "metal_plate"},
		RecipeCategoryAccessory:  {"gem", "gold", "silver"},
		RecipeCategoryConsumable: {"herb", "water", "magic_essence"},
	}
	
	materialList := materials[category]
	if len(materialList) == 0 {
		return "unknown_material"
	}
	
	return materialList[rf.random.Intn(len(materialList))]
}

// generateOutputID 生成产出ID
func (rf *RecipeFactory) generateOutputID(category RecipeCategory) string {
	outputs := map[RecipeCategory][]string{
		RecipeCategoryWeapon:     {"iron_sword", "steel_blade", "magic_wand"},
		RecipeCategoryArmor:      {"leather_armor", "iron_armor", "magic_robe"},
		RecipeCategoryAccessory:  {"power_ring", "magic_necklace", "protection_amulet"},
		RecipeCategoryConsumable: {"health_potion", "mana_potion", "antidote"},
	}
	
	outputList := outputs[category]
	if len(outputList) == 0 {
		return "unknown_item"
	}
	
	return outputList[rf.random.Intn(len(outputList))]
}

// MaterialFactory 材料工厂
type MaterialFactory struct {
	random *rand.Rand
}

// NewMaterialFactory 创建材料工厂
func NewMaterialFactory() *MaterialFactory {
	return &MaterialFactory{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CreateRandomMaterial 创建随机材料
func (mf *MaterialFactory) CreateRandomMaterial(materialType MaterialType, playerLevel int) *Material {
	quality := mf.determineQuality(playerLevel)
	quantity := mf.random.Intn(10) + 1
	
	materialID := mf.generateMaterialID(materialType)
	materialName := mf.generateMaterialName(materialType, quality)
	
	return NewMaterial(materialID, materialName, materialType, quality, quantity)
}

// determineQuality 确定品质
func (mf *MaterialFactory) determineQuality(playerLevel int) Quality {
	roll := mf.random.Float64()
	levelBonus := float64(playerLevel) / 100.0
	
	switch {
	case roll < 0.5-levelBonus*0.1:
		return QualityCommon
	case roll < 0.75-levelBonus*0.05:
		return QualityUncommon
	case roll < 0.9:
		return QualityRare
	case roll < 0.97:
		return QualityEpic
	case roll < 0.995:
		return QualityLegendary
	default:
		return QualityMythic
	}
}

// generateMaterialID 生成材料ID
func (mf *MaterialFactory) generateMaterialID(materialType MaterialType) string {
	ids := map[MaterialType][]string{
		MaterialTypeOre:     {"iron_ore", "copper_ore", "gold_ore"},
		MaterialTypeHerb:    {"healing_herb", "mana_herb", "poison_herb"},
		MaterialTypeLeather: {"wolf_leather", "bear_leather", "dragon_leather"},
		MaterialTypeCloth:   {"cotton_cloth", "silk_cloth", "magic_cloth"},
		MaterialTypeWood:    {"oak_wood", "pine_wood", "magic_wood"},
	}
	
	idList := ids[materialType]
	if len(idList) == 0 {
		return "unknown_material"
	}
	
	return idList[mf.random.Intn(len(idList))]
}

// generateMaterialName 生成材料名称
func (mf *MaterialFactory) generateMaterialName(materialType MaterialType, quality Quality) string {
	qualityPrefix := map[Quality]string{
		QualityCommon:    "普通的",
		QualityUncommon:  "优质的",
		QualityRare:      "稀有的",
		QualityEpic:      "史诗的",
		QualityLegendary: "传说的",
		QualityMythic:    "神话的",
	}
	
	typeNames := map[MaterialType]string{
		MaterialTypeOre:     "矿石",
		MaterialTypeHerb:    "草药",
		MaterialTypeLeather: "皮革",
		MaterialTypeCloth:   "布料",
		MaterialTypeWood:    "木材",
	}
	
	prefix := qualityPrefix[quality]
	typeName := typeNames[materialType]
	
	return prefix + typeName
}

// BonusCalculator 加成计算器
type BonusCalculator struct{}

// NewBonusCalculator 创建加成计算器
func NewBonusCalculator() *BonusCalculator {
	return &BonusCalculator{}
}

// CalculateTotalBonus 计算总加成
func (bc *BonusCalculator) CalculateTotalBonus(bonuses []*SynthesisBonus, bonusType BonusType) float64 {
	total := 0.0
	for _, bonus := range bonuses {
		if bonus.GetBonusType() == bonusType {
			total += bonus.GetBonusValue()
		}
	}
	return total
}

// ApplyQualityBonus 应用品质加成
func (bc *BonusCalculator) ApplyQualityBonus(baseValue int, quality Quality) int {
	multiplier := quality.GetQualityMultiplier()
	return int(float64(baseValue) * multiplier)
}