package plant

import (
	"fmt"
	"time"
)

// SeedType 种子类型
type SeedType int

const (
	SeedTypeWheat SeedType = iota + 1
	SeedTypeCorn
	SeedTypeRice
	SeedTypeTomato
	SeedTypePotato
	SeedTypeCarrot
	SeedTypeCabbage
	SeedTypeStrawberry
	SeedTypeApple
	SeedTypeOrange
)

// String 返回种子类型字符串
func (st SeedType) String() string {
	switch st {
	case SeedTypeWheat:
		return "wheat"
	case SeedTypeCorn:
		return "corn"
	case SeedTypeRice:
		return "rice"
	case SeedTypeTomato:
		return "tomato"
	case SeedTypePotato:
		return "potato"
	case SeedTypeCarrot:
		return "carrot"
	case SeedTypeCabbage:
		return "cabbage"
	case SeedTypeStrawberry:
		return "strawberry"
	case SeedTypeApple:
		return "apple"
	case SeedTypeOrange:
		return "orange"
	default:
		return "unknown"
	}
}

// GetDescription 获取描述
func (st SeedType) GetDescription() string {
	switch st {
	case SeedTypeWheat:
		return "小麦"
	case SeedTypeCorn:
		return "玉米"
	case SeedTypeRice:
		return "水稻"
	case SeedTypeTomato:
		return "番茄"
	case SeedTypePotato:
		return "土豆"
	case SeedTypeCarrot:
		return "胡萝卜"
	case SeedTypeCabbage:
		return "卷心菜"
	case SeedTypeStrawberry:
		return "草莓"
	case SeedTypeApple:
		return "苹果"
	case SeedTypeOrange:
		return "橙子"
	default:
		return "未知种子"
	}
}

// IsValid 检查种子类型是否有效
func (st SeedType) IsValid() bool {
	return st >= SeedTypeWheat && st <= SeedTypeOrange
}

// GetGrowthDuration 获取生长周期
func (st SeedType) GetGrowthDuration() time.Duration {
	switch st {
	case SeedTypeWheat:
		return 120 * time.Hour // 5天
	case SeedTypeCorn:
		return 168 * time.Hour // 7天
	case SeedTypeRice:
		return 144 * time.Hour // 6天
	case SeedTypeTomato:
		return 96 * time.Hour // 4天
	case SeedTypePotato:
		return 72 * time.Hour // 3天
	case SeedTypeCarrot:
		return 48 * time.Hour // 2天
	case SeedTypeCabbage:
		return 60 * time.Hour // 2.5天
	case SeedTypeStrawberry:
		return 36 * time.Hour // 1.5天
	case SeedTypeApple:
		return 240 * time.Hour // 10天
	case SeedTypeOrange:
		return 216 * time.Hour // 9天
	default:
		return 72 * time.Hour // 默认3天
	}
}

// GetGrowthRate 获取生长速度（每小时进度百分比）
func (st SeedType) GetGrowthRate() float64 {
	duration := st.GetGrowthDuration()
	return 100.0 / duration.Hours() // 100%进度除以总小时数
}

// GetBaseYield 获取基础产量
func (st SeedType) GetBaseYield() int {
	switch st {
	case SeedTypeWheat:
		return 8
	case SeedTypeCorn:
		return 12
	case SeedTypeRice:
		return 10
	case SeedTypeTomato:
		return 6
	case SeedTypePotato:
		return 4
	case SeedTypeCarrot:
		return 3
	case SeedTypeCabbage:
		return 5
	case SeedTypeStrawberry:
		return 2
	case SeedTypeApple:
		return 15
	case SeedTypeOrange:
		return 12
	default:
		return 5
	}
}

// GetBaseValue 获取基础价值
func (st SeedType) GetBaseValue() float64 {
	switch st {
	case SeedTypeWheat:
		return 10.0
	case SeedTypeCorn:
		return 15.0
	case SeedTypeRice:
		return 12.0
	case SeedTypeTomato:
		return 8.0
	case SeedTypePotato:
		return 6.0
	case SeedTypeCarrot:
		return 5.0
	case SeedTypeCabbage:
		return 7.0
	case SeedTypeStrawberry:
		return 4.0
	case SeedTypeApple:
		return 25.0
	case SeedTypeOrange:
		return 20.0
	default:
		return 8.0
	}
}

// GetBaseExperience 获取基础经验
func (st SeedType) GetBaseExperience() int {
	switch st {
	case SeedTypeWheat:
		return 20
	case SeedTypeCorn:
		return 30
	case SeedTypeRice:
		return 25
	case SeedTypeTomato:
		return 15
	case SeedTypePotato:
		return 10
	case SeedTypeCarrot:
		return 8
	case SeedTypeCabbage:
		return 12
	case SeedTypeStrawberry:
		return 6
	case SeedTypeApple:
		return 50
	case SeedTypeOrange:
		return 40
	default:
		return 15
	}
}

// GetWaterConsumption 获取水分消耗（每小时）
func (st SeedType) GetWaterConsumption() float64 {
	switch st {
	case SeedTypeWheat:
		return 1.5
	case SeedTypeCorn:
		return 2.0
	case SeedTypeRice:
		return 3.0 // 水稻需要更多水
	case SeedTypeTomato:
		return 2.5
	case SeedTypePotato:
		return 1.8
	case SeedTypeCarrot:
		return 1.2
	case SeedTypeCabbage:
		return 1.6
	case SeedTypeStrawberry:
		return 2.2
	case SeedTypeApple:
		return 1.0
	case SeedTypeOrange:
		return 1.2
	default:
		return 1.5
	}
}

// GetNutrientConsumption 获取营养消耗（每小时）
func (st SeedType) GetNutrientConsumption() float64 {
	switch st {
	case SeedTypeWheat:
		return 1.0
	case SeedTypeCorn:
		return 1.5
	case SeedTypeRice:
		return 1.2
	case SeedTypeTomato:
		return 1.8
	case SeedTypePotato:
		return 1.3
	case SeedTypeCarrot:
		return 0.8
	case SeedTypeCabbage:
		return 1.1
	case SeedTypeStrawberry:
		return 1.6
	case SeedTypeApple:
		return 0.8
	case SeedTypeOrange:
		return 0.9
	default:
		return 1.2
	}
}

// GetPreferredSoilType 获取偏好土壤类型
func (st SeedType) GetPreferredSoilType() SoilType {
	switch st {
	case SeedTypeWheat:
		return SoilTypeLoam
	case SeedTypeCorn:
		return SoilTypeLoam
	case SeedTypeRice:
		return SoilTypeClay // 水稻喜欢粘土
	case SeedTypeTomato:
		return SoilTypeLoam
	case SeedTypePotato:
		return SoilTypeSandy
	case SeedTypeCarrot:
		return SoilTypeSandy
	case SeedTypeCabbage:
		return SoilTypeLoam
	case SeedTypeStrawberry:
		return SoilTypeLoam
	case SeedTypeApple:
		return SoilTypeLoam
	case SeedTypeOrange:
		return SoilTypeLoam
	default:
		return SoilTypeLoam
	}
}

// GetCategory 获取作物类别
func (st SeedType) GetCategory() CropCategory {
	switch st {
	case SeedTypeWheat, SeedTypeCorn, SeedTypeRice:
		return CropCategoryGrain
	case SeedTypeTomato, SeedTypeCabbage:
		return CropCategoryVegetable
	case SeedTypePotato, SeedTypeCarrot:
		return CropCategoryRoot
	case SeedTypeStrawberry:
		return CropCategoryBerry
	case SeedTypeApple, SeedTypeOrange:
		return CropCategoryFruit
	default:
		return CropCategoryVegetable
	}
}

// GrowthStage 生长阶段
type GrowthStage int

const (
	GrowthStageSeed GrowthStage = iota + 1
	GrowthStageSeedling
	GrowthStageGrowing
	GrowthStageFlowering
	GrowthStageMature
)

// String 返回生长阶段字符串
func (gs GrowthStage) String() string {
	switch gs {
	case GrowthStageSeed:
		return "seed"
	case GrowthStageSeedling:
		return "seedling"
	case GrowthStageGrowing:
		return "growing"
	case GrowthStageFlowering:
		return "flowering"
	case GrowthStageMature:
		return "mature"
	default:
		return "unknown"
	}
}

// GetDescription 获取描述
func (gs GrowthStage) GetDescription() string {
	switch gs {
	case GrowthStageSeed:
		return "种子期"
	case GrowthStageSeedling:
		return "幼苗期"
	case GrowthStageGrowing:
		return "生长期"
	case GrowthStageFlowering:
		return "开花期"
	case GrowthStageMature:
		return "成熟期"
	default:
		return "未知阶段"
	}
}

// GetProgressRange 获取进度范围
func (gs GrowthStage) GetProgressRange() (float64, float64) {
	switch gs {
	case GrowthStageSeed:
		return 0.0, 25.0
	case GrowthStageSeedling:
		return 25.0, 50.0
	case GrowthStageGrowing:
		return 50.0, 75.0
	case GrowthStageFlowering:
		return 75.0, 100.0
	case GrowthStageMature:
		return 100.0, 100.0
	default:
		return 0.0, 0.0
	}
}

// SoilType 土壤类型
type SoilType int

const (
	SoilTypeSandy SoilType = iota + 1
	SoilTypeClay
	SoilTypeLoam
	SoilTypeSilt
	SoilTypePeat
	SoilTypeChalk
)

// String 返回土壤类型字符串
func (st SoilType) String() string {
	switch st {
	case SoilTypeSandy:
		return "sandy"
	case SoilTypeClay:
		return "clay"
	case SoilTypeLoam:
		return "loam"
	case SoilTypeSilt:
		return "silt"
	case SoilTypePeat:
		return "peat"
	case SoilTypeChalk:
		return "chalk"
	default:
		return "unknown"
	}
}

// GetDescription 获取描述
func (st SoilType) GetDescription() string {
	switch st {
	case SoilTypeSandy:
		return "沙土"
	case SoilTypeClay:
		return "粘土"
	case SoilTypeLoam:
		return "壤土"
	case SoilTypeSilt:
		return "淤泥土"
	case SoilTypePeat:
		return "泥炭土"
	case SoilTypeChalk:
		return "白垩土"
	default:
		return "未知土壤"
	}
}

// GetDrainageRate 获取排水率
func (st SoilType) GetDrainageRate() float64 {
	switch st {
	case SoilTypeSandy:
		return 0.9 // 沙土排水快
	case SoilTypeClay:
		return 0.2 // 粘土排水慢
	case SoilTypeLoam:
		return 0.6 // 壤土排水适中
	case SoilTypeSilt:
		return 0.4
	case SoilTypePeat:
		return 0.3
	case SoilTypeChalk:
		return 0.8
	default:
		return 0.5
	}
}

// GetNutrientRetention 获取营养保持率
func (st SoilType) GetNutrientRetention() float64 {
	switch st {
	case SoilTypeSandy:
		return 0.3 // 沙土营养流失快
	case SoilTypeClay:
		return 0.8 // 粘土营养保持好
	case SoilTypeLoam:
		return 0.7 // 壤土营养保持较好
	case SoilTypeSilt:
		return 0.6
	case SoilTypePeat:
		return 0.9 // 泥炭土营养丰富
	case SoilTypeChalk:
		return 0.4
	default:
		return 0.5
	}
}

// GetBaseProductivity 获取基础生产力
func (st SoilType) GetBaseProductivity() float64 {
	switch st {
	case SoilTypeSandy:
		return 0.8
	case SoilTypeClay:
		return 0.9
	case SoilTypeLoam:
		return 1.2 // 壤土最适合种植
	case SoilTypeSilt:
		return 1.0
	case SoilTypePeat:
		return 1.1
	case SoilTypeChalk:
		return 0.7
	default:
		return 1.0
	}
}

// Soil 土壤值对象
type Soil struct {
	Type       SoilType
	Fertility  float64 // 肥力 0-100
	PH         float64 // 酸碱度 0-14
	Moisture   float64 // 湿度 0-100
	Organic    float64 // 有机物含量 0-100
	Nitrogen   float64 // 氮含量 0-100
	Phosphorus float64 // 磷含量 0-100
	Potassium  float64 // 钾含量 0-100
	LastTested time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewSoil 创建土壤
func NewSoil(soilType SoilType, fertility, ph, moisture float64) *Soil {
	now := time.Now()
	return &Soil{
		Type:       soilType,
		Fertility:  fertility,
		PH:         ph,
		Moisture:   moisture,
		Organic:    30.0, // 默认有机物含量
		Nitrogen:   40.0, // 默认氮含量
		Phosphorus: 35.0, // 默认磷含量
		Potassium:  45.0, // 默认钾含量
		LastTested: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// GetType 获取土壤类型
func (s *Soil) GetType() SoilType {
	return s.Type
}

// GetFertility 获取肥力
func (s *Soil) GetFertility() float64 {
	return s.Fertility
}

// GetPH 获取酸碱度
func (s *Soil) GetPH() float64 {
	return s.PH
}

// GetMoisture 获取湿度
func (s *Soil) GetMoisture() float64 {
	return s.Moisture
}

// IsSuitableFor 检查是否适合种植指定作物
func (s *Soil) IsSuitableFor(seedType SeedType) bool {
	preferredSoil := seedType.GetPreferredSoilType()

	// 完全匹配最好
	if s.Type == preferredSoil {
		return true
	}

	// 壤土适合大多数作物
	if s.Type == SoilTypeLoam {
		return true
	}

	// 检查土壤条件是否满足最低要求
	return s.Fertility >= 30.0 && s.PH >= 5.5 && s.PH <= 8.5
}

// GetProductivityMultiplier 获取生产力倍率
func (s *Soil) GetProductivityMultiplier() float64 {
	baseProductivity := s.Type.GetBaseProductivity()

	// 肥力影响
	fertilityMultiplier := 0.5 + (s.Fertility/100.0)*0.8 // 0.5-1.3倍率

	// pH影响（6.0-7.5为最佳范围）
	phMultiplier := 1.0
	if s.PH < 5.0 || s.PH > 9.0 {
		phMultiplier = 0.6 // 极端pH值
	} else if s.PH < 6.0 || s.PH > 8.0 {
		phMultiplier = 0.8 // 偏酸或偏碱
	} else {
		phMultiplier = 1.2 // 最佳pH范围
	}

	// 有机物影响
	organicMultiplier := 0.8 + (s.Organic/100.0)*0.4 // 0.8-1.2倍率

	return baseProductivity * fertilityMultiplier * phMultiplier * organicMultiplier
}

// GetGrowthMultiplier 获取生长倍率
func (s *Soil) GetGrowthMultiplier(seedType SeedType) float64 {
	baseMultiplier := s.GetProductivityMultiplier()

	// 土壤类型匹配度
	preferredSoil := seedType.GetPreferredSoilType()
	typeMultiplier := 1.0
	if s.Type == preferredSoil {
		typeMultiplier = 1.2 // 完全匹配
	} else if s.Type == SoilTypeLoam {
		typeMultiplier = 1.1 // 壤土通用性好
	} else {
		typeMultiplier = 0.9 // 不匹配
	}

	return baseMultiplier * typeMultiplier
}

// GetYieldMultiplier 获取产量倍率
func (s *Soil) GetYieldMultiplier(seedType SeedType) float64 {
	// 基础倍率
	baseMultiplier := s.GetProductivityMultiplier()

	// 营养元素影响
	nutrientMultiplier := (s.Nitrogen + s.Phosphorus + s.Potassium) / 300.0 // 平均值
	if nutrientMultiplier > 1.0 {
		nutrientMultiplier = 1.0
	}
	nutrientMultiplier = 0.7 + nutrientMultiplier*0.6 // 0.7-1.3倍率

	return baseMultiplier * nutrientMultiplier
}

// GetQualityScore 获取质量分数
func (s *Soil) GetQualityScore() float64 {
	score := 0.0

	// 肥力贡献（30%）
	score += s.Fertility * 0.3

	// pH贡献（20%）
	if s.PH >= 6.0 && s.PH <= 7.5 {
		score += 20.0 // 最佳pH
	} else if s.PH >= 5.5 && s.PH <= 8.0 {
		score += 15.0 // 良好pH
	} else {
		score += 10.0 // 一般pH
	}

	// 有机物贡献（25%）
	score += s.Organic * 0.25

	// 营养元素贡献（25%）
	averageNutrient := (s.Nitrogen + s.Phosphorus + s.Potassium) / 3.0
	score += averageNutrient * 0.25

	return score
}

// GetValue 获取土壤价值
func (s *Soil) GetValue() float64 {
	baseValue := s.Type.GetBaseProductivity() * 100.0
	qualityMultiplier := s.GetQualityScore() / 100.0

	return baseValue * qualityMultiplier
}

// ApplyFertilizer 应用肥料
func (s *Soil) ApplyFertilizer(fertilizer *Fertilizer) {
	s.Fertility += fertilizer.GetFertilityBoost()
	s.Nitrogen += fertilizer.GetNitrogenContent()
	s.Phosphorus += fertilizer.GetPhosphorusContent()
	s.Potassium += fertilizer.GetPotassiumContent()
	s.Organic += fertilizer.GetOrganicContent()

	// 限制数值范围
	s.limitValues()
	s.UpdatedAt = time.Now()
}

// AddMoisture 增加湿度
func (s *Soil) AddMoisture(amount float64) {
	s.Moisture += amount
	if s.Moisture > 100.0 {
		s.Moisture = 100.0
	}
	s.UpdatedAt = time.Now()
}

// ApplyToCrop 应用到作物
func (s *Soil) ApplyToCrop(crop *Crop) {
	// 土壤会影响作物的营养水平
	if s.GetQualityScore() > 80.0 {
		crop.NutrientLevel += 1.0 // 高质量土壤缓慢提升营养
	} else if s.GetQualityScore() < 40.0 {
		crop.NutrientLevel -= 0.5 // 低质量土壤降低营养
	}

	// 限制营养水平范围
	if crop.NutrientLevel > 100.0 {
		crop.NutrientLevel = 100.0
	} else if crop.NutrientLevel < 0.0 {
		crop.NutrientLevel = 0.0
	}
}

// ApplyImprovement 应用改良
func (s *Soil) ApplyImprovement(value float64) {
	s.Fertility += value
	s.limitValues()
	s.UpdatedAt = time.Now()
}

// limitValues 限制数值范围
func (s *Soil) limitValues() {
	if s.Fertility > 100.0 {
		s.Fertility = 100.0
	}
	if s.Nitrogen > 100.0 {
		s.Nitrogen = 100.0
	}
	if s.Phosphorus > 100.0 {
		s.Phosphorus = 100.0
	}
	if s.Potassium > 100.0 {
		s.Potassium = 100.0
	}
	if s.Organic > 100.0 {
		s.Organic = 100.0
	}
}

// Fertilizer 肥料值对象
type Fertilizer struct {
	Type              FertilizerType
	Amount            float64
	NutrientValue     float64
	FertilityBoost    float64
	NitrogenContent   float64
	PhosphorusContent float64
	PotassiumContent  float64
	OrganicContent    float64
	GrowthBonus       *GrowthBonus
}

// NewFertilizer 创建肥料
func NewFertilizer(fertilizerType FertilizerType, amount float64) *Fertilizer {
	return &Fertilizer{
		Type:              fertilizerType,
		Amount:            amount,
		NutrientValue:     fertilizerType.GetNutrientValue() * amount,
		FertilityBoost:    fertilizerType.GetFertilityBoost() * amount,
		NitrogenContent:   fertilizerType.GetNitrogenContent() * amount,
		PhosphorusContent: fertilizerType.GetPhosphorusContent() * amount,
		PotassiumContent:  fertilizerType.GetPotassiumContent() * amount,
		OrganicContent:    fertilizerType.GetOrganicContent() * amount,
		GrowthBonus:       fertilizerType.GetGrowthBonus(),
	}
}

// GetType 获取类型
func (f *Fertilizer) GetType() FertilizerType {
	return f.Type
}

// GetAmount 获取数量
func (f *Fertilizer) GetAmount() float64 {
	return f.Amount
}

// GetNutrientValue 获取营养价值
func (f *Fertilizer) GetNutrientValue() float64 {
	return f.NutrientValue
}

// GetFertilityBoost 获取肥力提升
func (f *Fertilizer) GetFertilityBoost() float64 {
	return f.FertilityBoost
}

// GetNitrogenContent 获取氮含量
func (f *Fertilizer) GetNitrogenContent() float64 {
	return f.NitrogenContent
}

// GetPhosphorusContent 获取磷含量
func (f *Fertilizer) GetPhosphorusContent() float64 {
	return f.PhosphorusContent
}

// GetPotassiumContent 获取钾含量
func (f *Fertilizer) GetPotassiumContent() float64 {
	return f.PotassiumContent
}

// GetOrganicContent 获取有机物含量
func (f *Fertilizer) GetOrganicContent() float64 {
	return f.OrganicContent
}

// GetGrowthBonus 获取生长奖励
func (f *Fertilizer) GetGrowthBonus() *GrowthBonus {
	return f.GrowthBonus
}

// FertilizerType 肥料类型
type FertilizerType int

const (
	FertilizerTypeOrganic FertilizerType = iota + 1
	FertilizerTypeChemical
	FertilizerTypeCompost
	FertilizerTypeManure
	FertilizerTypeLiquid
	FertilizerTypeGranular
)

// String 返回肥料类型字符串
func (ft FertilizerType) String() string {
	switch ft {
	case FertilizerTypeOrganic:
		return "organic"
	case FertilizerTypeChemical:
		return "chemical"
	case FertilizerTypeCompost:
		return "compost"
	case FertilizerTypeManure:
		return "manure"
	case FertilizerTypeLiquid:
		return "liquid"
	case FertilizerTypeGranular:
		return "granular"
	default:
		return "unknown"
	}
}

// GetDescription 获取描述
func (ft FertilizerType) GetDescription() string {
	switch ft {
	case FertilizerTypeOrganic:
		return "有机肥"
	case FertilizerTypeChemical:
		return "化学肥料"
	case FertilizerTypeCompost:
		return "堆肥"
	case FertilizerTypeManure:
		return "粪肥"
	case FertilizerTypeLiquid:
		return "液体肥料"
	case FertilizerTypeGranular:
		return "颗粒肥料"
	default:
		return "未知肥料"
	}
}

// GetNutrientValue 获取营养价值
func (ft FertilizerType) GetNutrientValue() float64 {
	switch ft {
	case FertilizerTypeOrganic:
		return 15.0
	case FertilizerTypeChemical:
		return 25.0
	case FertilizerTypeCompost:
		return 12.0
	case FertilizerTypeManure:
		return 18.0
	case FertilizerTypeLiquid:
		return 20.0
	case FertilizerTypeGranular:
		return 22.0
	default:
		return 15.0
	}
}

// GetFertilityBoost 获取肥力提升
func (ft FertilizerType) GetFertilityBoost() float64 {
	switch ft {
	case FertilizerTypeOrganic:
		return 10.0
	case FertilizerTypeChemical:
		return 15.0
	case FertilizerTypeCompost:
		return 8.0
	case FertilizerTypeManure:
		return 12.0
	case FertilizerTypeLiquid:
		return 13.0
	case FertilizerTypeGranular:
		return 14.0
	default:
		return 10.0
	}
}

// GetNitrogenContent 获取氮含量
func (ft FertilizerType) GetNitrogenContent() float64 {
	switch ft {
	case FertilizerTypeOrganic:
		return 8.0
	case FertilizerTypeChemical:
		return 15.0
	case FertilizerTypeCompost:
		return 6.0
	case FertilizerTypeManure:
		return 10.0
	case FertilizerTypeLiquid:
		return 12.0
	case FertilizerTypeGranular:
		return 13.0
	default:
		return 8.0
	}
}

// GetPhosphorusContent 获取磷含量
func (ft FertilizerType) GetPhosphorusContent() float64 {
	switch ft {
	case FertilizerTypeOrganic:
		return 5.0
	case FertilizerTypeChemical:
		return 10.0
	case FertilizerTypeCompost:
		return 4.0
	case FertilizerTypeManure:
		return 6.0
	case FertilizerTypeLiquid:
		return 8.0
	case FertilizerTypeGranular:
		return 9.0
	default:
		return 5.0
	}
}

// GetPotassiumContent 获取钾含量
func (ft FertilizerType) GetPotassiumContent() float64 {
	switch ft {
	case FertilizerTypeOrganic:
		return 7.0
	case FertilizerTypeChemical:
		return 12.0
	case FertilizerTypeCompost:
		return 5.0
	case FertilizerTypeManure:
		return 8.0
	case FertilizerTypeLiquid:
		return 10.0
	case FertilizerTypeGranular:
		return 11.0
	default:
		return 7.0
	}
}

// GetOrganicContent 获取有机物含量
func (ft FertilizerType) GetOrganicContent() float64 {
	switch ft {
	case FertilizerTypeOrganic:
		return 20.0
	case FertilizerTypeChemical:
		return 2.0
	case FertilizerTypeCompost:
		return 25.0
	case FertilizerTypeManure:
		return 18.0
	case FertilizerTypeLiquid:
		return 5.0
	case FertilizerTypeGranular:
		return 3.0
	default:
		return 10.0
	}
}

// GetGrowthBonus 获取生长奖励
func (ft FertilizerType) GetGrowthBonus() *GrowthBonus {
	switch ft {
	case FertilizerTypeChemical:
		return &GrowthBonus{
			ID:         fmt.Sprintf("chemical_boost_%d", time.Now().UnixNano()),
			Type:       "growth_speed",
			Multiplier: 1.3,
			Duration:   48 * time.Hour,
			StartTime:  time.Now(),
		}
	case FertilizerTypeLiquid:
		return &GrowthBonus{
			ID:         fmt.Sprintf("liquid_boost_%d", time.Now().UnixNano()),
			Type:       "nutrient_absorption",
			Multiplier: 1.2,
			Duration:   24 * time.Hour,
			StartTime:  time.Now(),
		}
	default:
		return nil
	}
}

// 其他值对象

// FarmSize 农场大小
type FarmSize int

const (
	FarmSizeSmall FarmSize = iota + 1
	FarmSizeMedium
	FarmSizeLarge
	FarmSizeHuge
)

// String 返回农场大小字符串
func (fs FarmSize) String() string {
	switch fs {
	case FarmSizeSmall:
		return "small"
	case FarmSizeMedium:
		return "medium"
	case FarmSizeLarge:
		return "large"
	case FarmSizeHuge:
		return "huge"
	default:
		return "unknown"
	}
}

// GetMaxPlots 获取最大地块数
func (fs FarmSize) GetMaxPlots() int {
	switch fs {
	case FarmSizeSmall:
		return 4
	case FarmSizeMedium:
		return 9
	case FarmSizeLarge:
		return 16
	case FarmSizeHuge:
		return 25
	default:
		return 4
	}
}

// GetBaseValue 获取基础价值
func (fs FarmSize) GetBaseValue() float64 {
	switch fs {
	case FarmSizeSmall:
		return 1000.0
	case FarmSizeMedium:
		return 2500.0
	case FarmSizeLarge:
		return 5000.0
	case FarmSizeHuge:
		return 10000.0
	default:
		return 1000.0
	}
}

// GetExpansionCost 获取扩展成本
func (fs FarmSize) GetExpansionCost(currentSize FarmSize) *ExpansionCost {
	if fs <= currentSize {
		return nil
	}

	baseCost := fs.GetBaseValue() - currentSize.GetBaseValue()
	return &ExpansionCost{
		Gold:      baseCost,
		Materials: int(baseCost / 100),
		Time:      time.Duration(int(baseCost/500)) * time.Hour,
	}
}

// PlotSize 地块大小
type PlotSize int

const (
	PlotSizeSmall PlotSize = iota + 1
	PlotSizeMedium
	PlotSizeLarge
)

// String 返回地块大小字符串
func (ps PlotSize) String() string {
	switch ps {
	case PlotSizeSmall:
		return "small"
	case PlotSizeMedium:
		return "medium"
	case PlotSizeLarge:
		return "large"
	default:
		return "unknown"
	}
}

// GetCapacity 获取容量
func (ps PlotSize) GetCapacity() int {
	switch ps {
	case PlotSizeSmall:
		return 1
	case PlotSizeMedium:
		return 4
	case PlotSizeLarge:
		return 9
	default:
		return 1
	}
}

// ToolType 工具类型
type ToolType int

const (
	ToolTypeHoe ToolType = iota + 1
	ToolTypeWateringCan
	ToolTypeFertilizerSpreader
	ToolTypeHarvester
	ToolTypePesticide
	ToolTypeTractor
)

// String 返回工具类型字符串
func (tt ToolType) String() string {
	switch tt {
	case ToolTypeHoe:
		return "hoe"
	case ToolTypeWateringCan:
		return "watering_can"
	case ToolTypeFertilizerSpreader:
		return "fertilizer_spreader"
	case ToolTypeHarvester:
		return "harvester"
	case ToolTypePesticide:
		return "pesticide"
	case ToolTypeTractor:
		return "tractor"
	default:
		return "unknown"
	}
}

// GetDescription 获取描述
func (tt ToolType) GetDescription() string {
	switch tt {
	case ToolTypeHoe:
		return "锄头"
	case ToolTypeWateringCan:
		return "洒水壶"
	case ToolTypeFertilizerSpreader:
		return "施肥器"
	case ToolTypeHarvester:
		return "收割机"
	case ToolTypePesticide:
		return "杀虫剂"
	case ToolTypeTractor:
		return "拖拉机"
	default:
		return "未知工具"
	}
}

// GetBaseValue 获取基础价值
func (tt ToolType) GetBaseValue() float64 {
	switch tt {
	case ToolTypeHoe:
		return 50.0
	case ToolTypeWateringCan:
		return 30.0
	case ToolTypeFertilizerSpreader:
		return 80.0
	case ToolTypeHarvester:
		return 200.0
	case ToolTypePesticide:
		return 40.0
	case ToolTypeTractor:
		return 500.0
	default:
		return 50.0
	}
}

// GetEffect 获取效果
func (tt ToolType) GetEffect(level int, efficiency float64) *ToolEffect {
	baseValue := float64(level) * efficiency

	switch tt {
	case ToolTypeHoe:
		return &ToolEffect{Type: "soil_improvement", Value: baseValue * 2.0}
	case ToolTypeWateringCan:
		return &ToolEffect{Type: "watering_efficiency", Value: baseValue * 1.5}
	case ToolTypeFertilizerSpreader:
		return &ToolEffect{Type: "fertilizer_efficiency", Value: baseValue * 1.8}
	case ToolTypeHarvester:
		return &ToolEffect{Type: "harvest_speed", Value: baseValue * 2.5}
	case ToolTypePesticide:
		return &ToolEffect{Type: "pest_control", Value: baseValue * 3.0}
	case ToolTypeTractor:
		return &ToolEffect{Type: "overall_efficiency", Value: baseValue * 1.2}
	default:
		return &ToolEffect{Type: "unknown", Value: baseValue}
	}
}

// CropQuality 作物品质
type CropQuality int

const (
	CropQualityCommon CropQuality = iota + 1
	CropQualityUncommon
	CropQualityRare
	CropQualityEpic
	CropQualityLegendary
)

// String 返回品质字符串
func (cq CropQuality) String() string {
	switch cq {
	case CropQualityCommon:
		return "common"
	case CropQualityUncommon:
		return "uncommon"
	case CropQualityRare:
		return "rare"
	case CropQualityEpic:
		return "epic"
	case CropQualityLegendary:
		return "legendary"
	default:
		return "unknown"
	}
}

// GetDescription 获取描述
func (cq CropQuality) GetDescription() string {
	switch cq {
	case CropQualityCommon:
		return "普通"
	case CropQualityUncommon:
		return "优良"
	case CropQualityRare:
		return "稀有"
	case CropQualityEpic:
		return "史诗"
	case CropQualityLegendary:
		return "传说"
	default:
		return "未知品质"
	}
}

// GetValueMultiplier 获取价值倍率
func (cq CropQuality) GetValueMultiplier() float64 {
	switch cq {
	case CropQualityCommon:
		return 1.0
	case CropQualityUncommon:
		return 1.5
	case CropQualityRare:
		return 2.0
	case CropQualityEpic:
		return 3.0
	case CropQualityLegendary:
		return 5.0
	default:
		return 1.0
	}
}

// GetExperienceMultiplier 获取经验倍率
func (cq CropQuality) GetExperienceMultiplier() float64 {
	switch cq {
	case CropQualityCommon:
		return 0.0
	case CropQualityUncommon:
		return 0.2
	case CropQualityRare:
		return 0.5
	case CropQualityEpic:
		return 1.0
	case CropQualityLegendary:
		return 2.0
	default:
		return 0.0
	}
}

// CropCategory 作物类别
type CropCategory int

const (
	CropCategoryGrain CropCategory = iota + 1
	CropCategoryVegetable
	CropCategoryFruit
	CropCategoryRoot
	CropCategoryBerry
	CropCategoryHerb
)

// String 返回类别字符串
func (cc CropCategory) String() string {
	switch cc {
	case CropCategoryGrain:
		return "grain"
	case CropCategoryVegetable:
		return "vegetable"
	case CropCategoryFruit:
		return "fruit"
	case CropCategoryRoot:
		return "root"
	case CropCategoryBerry:
		return "berry"
	case CropCategoryHerb:
		return "herb"
	default:
		return "unknown"
	}
}

// GetDescription 获取描述
func (cc CropCategory) GetDescription() string {
	switch cc {
	case CropCategoryGrain:
		return "谷物"
	case CropCategoryVegetable:
		return "蔬菜"
	case CropCategoryFruit:
		return "水果"
	case CropCategoryRoot:
		return "根茎类"
	case CropCategoryBerry:
		return "浆果"
	case CropCategoryHerb:
		return "草药"
	default:
		return "未知类别"
	}
}

// 辅助结构体

// ExpansionCost 扩展成本
type ExpansionCost struct {
	Gold      float64
	Materials int
	Time      time.Duration
}

// FarmResources 农场资源
type FarmResources struct {
	Gold       float64
	Seeds      map[SeedType]int
	Fertilizer map[FertilizerType]float64
	Water      float64
	Harvest    map[SeedType]map[CropQuality]int
	UpdatedAt  time.Time
}

// NewFarmResources 创建农场资源
func NewFarmResources() *FarmResources {
	return &FarmResources{
		Gold:       1000.0, // 初始金币
		Seeds:      make(map[SeedType]int),
		Fertilizer: make(map[FertilizerType]float64),
		Water:      100.0, // 初始水量
		Harvest:    make(map[SeedType]map[CropQuality]int),
		UpdatedAt:  time.Now(),
	}
}

// HasEnoughSeeds 检查种子是否足够
func (fr *FarmResources) HasEnoughSeeds(seedType SeedType, quantity int) bool {
	return fr.Seeds[seedType] >= quantity
}

// ConsumeSeeds 消耗种子
func (fr *FarmResources) ConsumeSeeds(seedType SeedType, quantity int) {
	fr.Seeds[seedType] -= quantity
	if fr.Seeds[seedType] < 0 {
		fr.Seeds[seedType] = 0
	}
	fr.UpdatedAt = time.Now()
}

// HasEnoughFertilizer 检查肥料是否足够
func (fr *FarmResources) HasEnoughFertilizer(fertilizerType FertilizerType, amount float64) bool {
	return fr.Fertilizer[fertilizerType] >= amount
}

// ConsumeFertilizer 消耗肥料
func (fr *FarmResources) ConsumeFertilizer(fertilizerType FertilizerType, amount float64) {
	fr.Fertilizer[fertilizerType] -= amount
	if fr.Fertilizer[fertilizerType] < 0 {
		fr.Fertilizer[fertilizerType] = 0
	}
	fr.UpdatedAt = time.Now()
}

// HasEnoughWater 检查水是否足够
func (fr *FarmResources) HasEnoughWater(amount float64) bool {
	return fr.Water >= amount
}

// ConsumeWater 消耗水
func (fr *FarmResources) ConsumeWater(amount float64) {
	fr.Water -= amount
	if fr.Water < 0 {
		fr.Water = 0
	}
	fr.UpdatedAt = time.Now()
}

// AddHarvest 添加收获物
func (fr *FarmResources) AddHarvest(seedType SeedType, quantity int, quality CropQuality) {
	if fr.Harvest[seedType] == nil {
		fr.Harvest[seedType] = make(map[CropQuality]int)
	}
	fr.Harvest[seedType][quality] += quantity
	fr.UpdatedAt = time.Now()
}

// CanAfford 检查是否能承担成本
func (fr *FarmResources) CanAfford(cost *ExpansionCost) bool {
	return fr.Gold >= cost.Gold
}

// GetTotalValue 获取总价值
func (fr *FarmResources) GetTotalValue() float64 {
	totalValue := fr.Gold

	// 种子价值
	for seedType, quantity := range fr.Seeds {
		totalValue += seedType.GetBaseValue() * float64(quantity) * 0.5 // 种子价值为作物价值的一半
	}

	// 收获物价值
	for seedType, qualityMap := range fr.Harvest {
		for quality, quantity := range qualityMap {
			baseValue := seedType.GetBaseValue()
			qualityMultiplier := quality.GetValueMultiplier()
			totalValue += baseValue * qualityMultiplier * float64(quantity)
		}
	}

	return totalValue
}

// FarmStatistics 农场统计
type FarmStatistics struct {
	TotalPlantings  int
	TotalHarvests   int
	TotalYield      int
	TotalExperience int
	PlantingsByType map[SeedType]int
	HarvestsByType  map[SeedType]int
	FertilizerUsage map[FertilizerType]float64
	WateringCount   int
	ToolUsage       map[ToolType]int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewFarmStatistics 创建农场统计
func NewFarmStatistics() *FarmStatistics {
	now := time.Now()
	return &FarmStatistics{
		PlantingsByType: make(map[SeedType]int),
		HarvestsByType:  make(map[SeedType]int),
		FertilizerUsage: make(map[FertilizerType]float64),
		ToolUsage:       make(map[ToolType]int),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// AddPlantingActivity 添加种植活动
func (fs *FarmStatistics) AddPlantingActivity(seedType SeedType, quantity int) {
	fs.TotalPlantings += quantity
	fs.PlantingsByType[seedType] += quantity
	fs.UpdatedAt = time.Now()
}

// AddHarvestActivity 添加收获活动
func (fs *FarmStatistics) AddHarvestActivity(seedType SeedType, yield int, quality CropQuality) {
	fs.TotalHarvests++
	fs.TotalYield += yield
	fs.HarvestsByType[seedType] += yield
	fs.UpdatedAt = time.Now()
}

// AddFertilizerUsage 添加肥料使用
func (fs *FarmStatistics) AddFertilizerUsage(fertilizerType FertilizerType, amount float64) {
	fs.FertilizerUsage[fertilizerType] += amount
	fs.UpdatedAt = time.Now()
}

// AddWateringActivity 添加浇水活动
func (fs *FarmStatistics) AddWateringActivity(plotCount int, waterAmount float64) {
	fs.WateringCount += plotCount
	fs.UpdatedAt = time.Now()
}

// AddToolUsage 添加工具使用
func (fs *FarmStatistics) AddToolUsage(toolType ToolType) {
	fs.ToolUsage[toolType]++
	fs.UpdatedAt = time.Now()
}

// SeasonModifier 季节修正
type SeasonModifier struct {
	CurrentSeason       Season
	GrowthMultiplier    float64
	YieldMultiplier     float64
	QualityMultiplier   float64
	WaterConsumption    float64
	NutrientConsumption float64
	SeasonEffects       map[SeedType]float64
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// NewSeasonModifier 创建季节修正
func NewSeasonModifier() *SeasonModifier {
	now := time.Now()
	currentSeason := getCurrentSeason(now)

	return &SeasonModifier{
		CurrentSeason:       currentSeason,
		GrowthMultiplier:    currentSeason.GetGrowthMultiplier(),
		YieldMultiplier:     currentSeason.GetYieldMultiplier(),
		QualityMultiplier:   currentSeason.GetQualityMultiplier(),
		WaterConsumption:    currentSeason.GetWaterConsumptionMultiplier(),
		NutrientConsumption: currentSeason.GetNutrientConsumptionMultiplier(),
		SeasonEffects:       make(map[SeedType]float64),
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}

// ApplyToCrop 应用到作物
func (sm *SeasonModifier) ApplyToCrop(crop *Crop) {
	// 应用季节效果到作物
	if effect, exists := sm.SeasonEffects[crop.SeedType]; exists {
		bonus := &GrowthBonus{
			ID:         fmt.Sprintf("season_%s_%d", sm.CurrentSeason.String(), time.Now().UnixNano()),
			Type:       "seasonal_effect",
			Multiplier: effect,
			Duration:   24 * time.Hour,
			StartTime:  time.Now(),
		}
		crop.AddBonus(bonus)
	}
}

// GetProductivityMultiplier 获取生产力倍率
func (sm *SeasonModifier) GetProductivityMultiplier() float64 {
	return sm.GrowthMultiplier * sm.YieldMultiplier
}

// GetYieldMultiplier 获取产量倍率
func (sm *SeasonModifier) GetYieldMultiplier(seedType SeedType) float64 {
	baseMultiplier := sm.YieldMultiplier
	if effect, exists := sm.SeasonEffects[seedType]; exists {
		return baseMultiplier * effect
	}
	return baseMultiplier
}

// Season 季节
type Season int

const (
	SeasonSpring Season = iota + 1
	SeasonSummer
	SeasonAutumn
	SeasonWinter
)

// String 返回季节字符串
func (s Season) String() string {
	switch s {
	case SeasonSpring:
		return "spring"
	case SeasonSummer:
		return "summer"
	case SeasonAutumn:
		return "autumn"
	case SeasonWinter:
		return "winter"
	default:
		return "unknown"
	}
}

// GetGrowthMultiplier 获取生长倍率
func (s Season) GetGrowthMultiplier() float64 {
	switch s {
	case SeasonSpring:
		return 1.3 // 春季生长最快
	case SeasonSummer:
		return 1.1
	case SeasonAutumn:
		return 0.9
	case SeasonWinter:
		return 0.6 // 冬季生长最慢
	default:
		return 1.0
	}
}

// GetYieldMultiplier 获取产量倍率
func (s Season) GetYieldMultiplier() float64 {
	switch s {
	case SeasonSpring:
		return 1.2
	case SeasonSummer:
		return 1.3 // 夏季产量最高
	case SeasonAutumn:
		return 1.1
	case SeasonWinter:
		return 0.7
	default:
		return 1.0
	}
}

// GetQualityMultiplier 获取品质倍率
func (s Season) GetQualityMultiplier() float64 {
	switch s {
	case SeasonSpring:
		return 1.1
	case SeasonSummer:
		return 1.0
	case SeasonAutumn:
		return 1.2 // 秋季品质最好
	case SeasonWinter:
		return 0.8
	default:
		return 1.0
	}
}

// GetWaterConsumptionMultiplier 获取水分消耗倍率
func (s Season) GetWaterConsumptionMultiplier() float64 {
	switch s {
	case SeasonSpring:
		return 1.0
	case SeasonSummer:
		return 1.5 // 夏季需要更多水分
	case SeasonAutumn:
		return 0.8
	case SeasonWinter:
		return 0.6
	default:
		return 1.0
	}
}

// GetNutrientConsumptionMultiplier 获取营养消耗倍率
func (s Season) GetNutrientConsumptionMultiplier() float64 {
	switch s {
	case SeasonSpring:
		return 1.2 // 春季生长需要更多营养
	case SeasonSummer:
		return 1.1
	case SeasonAutumn:
		return 0.9
	case SeasonWinter:
		return 0.7
	default:
		return 1.0
	}
}

// AutomationSettings 自动化设置
type AutomationSettings struct {
	AutoWatering      bool
	AutoFertilizing   bool
	AutoHarvesting    bool
	AutoPestControl   bool
	WaterThreshold    float64
	NutrientThreshold float64
	HarvestDelay      time.Duration
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// NewAutomationSettings 创建自动化设置
func NewAutomationSettings() *AutomationSettings {
	now := time.Now()
	return &AutomationSettings{
		AutoWatering:      false,
		AutoFertilizing:   false,
		AutoHarvesting:    false,
		AutoPestControl:   false,
		WaterThreshold:    30.0,
		NutrientThreshold: 30.0,
		HarvestDelay:      1 * time.Hour,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// 辅助函数

// getCurrentSeason 获取当前季节
func getCurrentSeason(currentTime time.Time) Season {
	month := currentTime.Month()
	switch {
	case month >= 3 && month <= 5:
		return SeasonSpring
	case month >= 6 && month <= 8:
		return SeasonSummer
	case month >= 9 && month <= 11:
		return SeasonAutumn
	default:
		return SeasonWinter
	}
}
