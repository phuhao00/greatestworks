package plant

import (
	"errors"
	"fmt"
	"math"
	"time"
)

// PlantService 种植领域服务
type PlantService struct {
	seedTemplates       map[SeedType]*SeedTemplate
	soilTemplates       map[SoilType]*SoilTemplate
	fertilizerTemplates map[FertilizerType]*FertilizerTemplate
	toolTemplates       map[ToolType]*ToolTemplate
	growthRules         []*GrowthRule
	seasonalEffects     map[Season]map[SeedType]float64
	climateZones        map[string]*ClimateZone
	pestDiseaseRules    []*PestDiseaseRule
	qualityRules        []*QualityRule
	createdAt           time.Time
	updatedAt           time.Time
}

// NewPlantService 创建种植服务
func NewPlantService() *PlantService {
	now := time.Now()
	service := &PlantService{
		seedTemplates:       make(map[SeedType]*SeedTemplate),
		soilTemplates:       make(map[SoilType]*SoilTemplate),
		fertilizerTemplates: make(map[FertilizerType]*FertilizerTemplate),
		toolTemplates:       make(map[ToolType]*ToolTemplate),
		growthRules:         make([]*GrowthRule, 0),
		seasonalEffects:     make(map[Season]map[SeedType]float64),
		climateZones:        make(map[string]*ClimateZone),
		pestDiseaseRules:    make([]*PestDiseaseRule, 0),
		qualityRules:        make([]*QualityRule, 0),
		createdAt:           now,
		updatedAt:           now,
	}

	// 初始化默认模板和规则
	service.initializeDefaultTemplates()
	service.initializeDefaultRules()
	service.initializeSeasonalEffects()
	service.initializeClimateZones()

	return service
}

// RegisterSeedTemplate 注册种子模板
func (ps *PlantService) RegisterSeedTemplate(seedType SeedType, template *SeedTemplate) {
	ps.seedTemplates[seedType] = template
	ps.updatedAt = time.Now()
}

// GetSeedTemplate 获取种子模板
func (ps *PlantService) GetSeedTemplate(seedType SeedType) *SeedTemplate {
	return ps.seedTemplates[seedType]
}

// RegisterSoilTemplate 注册土壤模板
func (ps *PlantService) RegisterSoilTemplate(soilType SoilType, template *SoilTemplate) {
	ps.soilTemplates[soilType] = template
	ps.updatedAt = time.Now()
}

// GetSoilTemplate 获取土壤模板
func (ps *PlantService) GetSoilTemplate(soilType SoilType) *SoilTemplate {
	return ps.soilTemplates[soilType]
}

// RegisterFertilizerTemplate 注册肥料模板
func (ps *PlantService) RegisterFertilizerTemplate(fertilizerType FertilizerType, template *FertilizerTemplate) {
	ps.fertilizerTemplates[fertilizerType] = template
	ps.updatedAt = time.Now()
}

// GetFertilizerTemplate 获取肥料模板
func (ps *PlantService) GetFertilizerTemplate(fertilizerType FertilizerType) *FertilizerTemplate {
	return ps.fertilizerTemplates[fertilizerType]
}

// RegisterToolTemplate 注册工具模板
func (ps *PlantService) RegisterToolTemplate(toolType ToolType, template *ToolTemplate) {
	ps.toolTemplates[toolType] = template
	ps.updatedAt = time.Now()
}

// GetToolTemplate 获取工具模板
func (ps *PlantService) GetToolTemplate(toolType ToolType) *ToolTemplate {
	return ps.toolTemplates[toolType]
}

// AddGrowthRule 添加生长规则
func (ps *PlantService) AddGrowthRule(rule *GrowthRule) {
	ps.growthRules = append(ps.growthRules, rule)
	ps.updatedAt = time.Now()
}

// GetGrowthRules 获取生长规则
func (ps *PlantService) GetGrowthRules() []*GrowthRule {
	return ps.growthRules
}

// RegisterClimateZone 注册气候区域
func (ps *PlantService) RegisterClimateZone(zoneID string, zone *ClimateZone) {
	ps.climateZones[zoneID] = zone
	ps.updatedAt = time.Now()
}

// GetClimateZone 获取气候区域
func (ps *PlantService) GetClimateZone(zoneID string) *ClimateZone {
	return ps.climateZones[zoneID]
}

// CalculateOptimalPlantingTime 计算最佳种植时间
func (ps *PlantService) CalculateOptimalPlantingTime(seedType SeedType, climateZone string) (time.Time, error) {
	if !seedType.IsValid() {
		return time.Time{}, ErrInvalidSeedType
	}

	zone := ps.GetClimateZone(climateZone)
	if zone == nil {
		zone = ps.getDefaultClimateZone()
	}

	now := time.Now()
	currentSeason := getCurrentSeason(now)

	// 获取种子的最佳种植季节
	optimalSeasons := ps.getOptimalSeasonsForSeed(seedType)

	// 如果当前季节是最佳季节，返回当前时间
	for _, season := range optimalSeasons {
		if season == currentSeason {
			return now, nil
		}
	}

	// 计算下一个最佳季节的开始时间
	nextOptimalTime := ps.calculateNextSeasonTime(now, optimalSeasons[0])
	return nextOptimalTime, nil
}

// CalculateGrowthProgress 计算生长进度
func (ps *PlantService) CalculateGrowthProgress(crop *Crop, deltaTime time.Duration) (float64, error) {
	if crop == nil {
		return 0, ErrInvalidCrop
	}

	// 基础生长速度
	baseGrowthRate := crop.SeedType.GetGrowthRate()

	// 应用生长规则
	actualGrowthRate := ps.applyGrowthRules(crop, baseGrowthRate)

	// 计算进度增量
	progressDelta := actualGrowthRate * deltaTime.Hours()

	return progressDelta, nil
}

// CalculateYield 计算产量
func (ps *PlantService) CalculateYield(crop *Crop, soil *Soil, season Season) (int, error) {
	if crop == nil {
		return 0, ErrInvalidCrop
	}

	if soil == nil {
		return 0, ErrInvalidSoil
	}

	// 基础产量
	baseYield := crop.GetBaseYield()

	// 土壤影响
	soilMultiplier := soil.GetYieldMultiplier(crop.SeedType)

	// 季节影响
	seasonMultiplier := ps.getSeasonalEffect(season, crop.SeedType)

	// 健康状态影响
	healthMultiplier := crop.GetHealthScore() / 100.0

	// 照料质量影响
	careMultiplier := crop.GetCareQualityMultiplier()

	// 计算最终产量
	finalYield := float64(baseYield) * soilMultiplier * seasonMultiplier * healthMultiplier * careMultiplier

	return int(math.Round(finalYield)), nil
}

// CalculateQuality 计算品质
func (ps *PlantService) CalculateQuality(crop *Crop, soil *Soil, season Season) (CropQuality, error) {
	if crop == nil {
		return CropQualityCommon, ErrInvalidCrop
	}

	if soil == nil {
		return CropQualityCommon, ErrInvalidSoil
	}

	// 计算质量分数
	qualityScore := 0.0

	// 土壤质量贡献（30%）
	qualityScore += soil.GetQualityScore() * 0.3

	// 作物健康贡献（25%）
	qualityScore += crop.GetHealthScore() * 0.25

	// 照料质量贡献（25%）
	qualityScore += crop.GetCareQualityScore() * 0.25

	// 季节影响贡献（20%）
	seasonBonus := season.GetQualityMultiplier() * 20.0
	qualityScore += seasonBonus

	// 应用质量规则
	qualityScore = ps.applyQualityRules(crop, soil, qualityScore)

	// 转换为品质等级
	return ps.scoreToQuality(qualityScore), nil
}

// ValidatePlantingConditions 验证种植条件
func (ps *PlantService) ValidatePlantingConditions(seedType SeedType, soil *Soil, season Season, climateZone string) error {
	if !seedType.IsValid() {
		return ErrInvalidSeedType
	}

	if soil == nil {
		return ErrInvalidSoil
	}

	// 检查土壤适宜性
	if !soil.IsSuitableFor(seedType) {
		return errors.New("soil is not suitable")
	}

	// 检查季节适宜性
	if !ps.isSeasonSuitableForSeed(seedType, season) {
		return ErrSeasonNotSuitable
	}

	// 检查气候区域适宜性
	zone := ps.GetClimateZone(climateZone)
	if zone != nil && !zone.IsSuitableFor(seedType) {
		return ErrClimateNotSuitable
	}

	return nil
}

// CalculateFertilizerEffect 计算肥料效果
func (ps *PlantService) CalculateFertilizerEffect(fertilizer *Fertilizer, soil *Soil, crop *Crop) (*FertilizerEffect, error) {
	if fertilizer == nil {
		return nil, ErrInvalidFertilizer
	}

	if soil == nil {
		return nil, ErrInvalidSoil
	}

	template := ps.GetFertilizerTemplate(fertilizer.GetType())
	if template == nil {
		return nil, ErrFertilizerTemplateNotFound
	}

	// 计算基础效果
	baseEffect := template.GetBaseEffect()

	// 土壤类型影响
	soilMultiplier := template.GetSoilMultiplier(soil.GetType())

	// 作物类型影响
	cropMultiplier := 1.0
	if crop != nil {
		cropMultiplier = template.GetCropMultiplier(crop.GetSeedType())
	}

	// 计算最终效果
	finalEffect := baseEffect * soilMultiplier * cropMultiplier * fertilizer.GetAmount()

	return &FertilizerEffect{
		FertilityBoost:   finalEffect * 0.4,
		NutrientBoost:    finalEffect * 0.6,
		GrowthSpeedBoost: finalEffect * 0.2,
		Duration:         template.GetEffectDuration(),
	}, nil
}

// CalculateToolEfficiency 计算工具效率
func (ps *PlantService) CalculateToolEfficiency(tool *FarmTool, operation string, target interface{}) (float64, error) {
	if tool == nil {
		return 1.0, errors.New("Farm tool not found")
	}

	if !tool.IsUsable() {
		return 1.0, errors.New("Farm tool is not usable")
	}

	template := ps.GetToolTemplate(tool.GetType())
	if template == nil {
		return 1.0, ErrToolTemplateNotFound
	}

	// 基础效率
	baseEfficiency := tool.GetEfficiency()

	// 操作类型影响
	operationMultiplier := template.GetOperationMultiplier(operation)

	// 工具等级影响
	levelMultiplier := 1.0 + float64(tool.GetLevel())*0.1

	// 耐久度影响
	durabilityMultiplier := tool.GetDurability() / tool.MaxDurability
	if durabilityMultiplier < 0.5 {
		durabilityMultiplier = 0.5 // 最低50%效率
	}

	// 计算最终效率
	finalEfficiency := baseEfficiency * operationMultiplier * levelMultiplier * durabilityMultiplier

	return finalEfficiency, nil
}

// DetectPestsAndDiseases 检测病虫害
func (ps *PlantService) DetectPestsAndDiseases(crop *Crop, soil *Soil, season Season, climateZone string) ([]*PestDiseaseEvent, error) {
	if crop == nil {
		return nil, ErrInvalidCrop
	}

	events := make([]*PestDiseaseEvent, 0)

	// 应用病虫害规则
	for _, rule := range ps.pestDiseaseRules {
		if rule.ShouldTrigger(crop, soil, season, climateZone) {
			event := rule.CreateEvent(crop)
			events = append(events, event)
		}
	}

	return events, nil
}

// CalculateWaterRequirement 计算水分需求
func (ps *PlantService) CalculateWaterRequirement(crop *Crop, soil *Soil, season Season, climateZone string) (float64, error) {
	if crop == nil {
		return 0, ErrInvalidCrop
	}

	// 基础水分消耗
	baseConsumption := crop.SeedType.GetWaterConsumption()

	// 生长阶段影响
	stageMultiplier := ps.getGrowthStageWaterMultiplier(crop.GetGrowthStage())

	// 季节影响
	seasonMultiplier := season.GetWaterConsumptionMultiplier()

	// 土壤影响
	soilMultiplier := 1.0
	if soil != nil {
		// 排水好的土壤需要更多水分
		drainageRate := soil.GetType().GetDrainageRate()
		soilMultiplier = 0.8 + drainageRate*0.4 // 0.8-1.2倍率
	}

	// 气候区域影响
	climateMultiplier := 1.0
	zone := ps.GetClimateZone(climateZone)
	if zone != nil {
		climateMultiplier = zone.GetWaterRequirementMultiplier()
	}

	// 计算最终需求
	finalRequirement := baseConsumption * stageMultiplier * seasonMultiplier * soilMultiplier * climateMultiplier

	return finalRequirement, nil
}

// CalculateNutrientRequirement 计算营养需求
func (ps *PlantService) CalculateNutrientRequirement(crop *Crop, soil *Soil, season Season) (map[string]float64, error) {
	if crop == nil {
		return nil, ErrInvalidCrop
	}

	// 基础营养消耗
	baseConsumption := crop.SeedType.GetNutrientConsumption()

	// 生长阶段影响
	stageMultiplier := ps.getGrowthStageNutrientMultiplier(crop.GetGrowthStage())

	// 季节影响
	seasonMultiplier := season.GetNutrientConsumptionMultiplier()

	// 土壤营养保持率影响
	soilMultiplier := 1.0
	if soil != nil {
		retentionRate := soil.GetType().GetNutrientRetention()
		soilMultiplier = 2.0 - retentionRate // 保持率低需要更多营养
	}

	// 计算各种营养需求
	requirements := map[string]float64{
		"nitrogen":   baseConsumption * 0.4 * stageMultiplier * seasonMultiplier * soilMultiplier,
		"phosphorus": baseConsumption * 0.3 * stageMultiplier * seasonMultiplier * soilMultiplier,
		"potassium":  baseConsumption * 0.3 * stageMultiplier * seasonMultiplier * soilMultiplier,
	}

	return requirements, nil
}

// GetOptimalHarvestTime 获取最佳收获时间
func (ps *PlantService) GetOptimalHarvestTime(crop *Crop) (time.Time, error) {
	if crop == nil {
		return time.Time{}, ErrInvalidCrop
	}

	// 基础收获时间
	baseHarvestTime := crop.ExpectedHarvestTime

	// 考虑生长进度调整
	if crop.GetGrowthProgress() < 100.0 {
		// 根据当前进度和生长速度估算剩余时间
		remainingProgress := 100.0 - crop.GetGrowthProgress()
		growthRate := crop.SeedType.GetGrowthRate()

		// 应用当前环境因素
		environmentMultiplier := ps.calculateCurrentEnvironmentMultiplier(crop)
		actualGrowthRate := growthRate * environmentMultiplier

		remainingHours := remainingProgress / actualGrowthRate
		adjustedHarvestTime := time.Now().Add(time.Duration(remainingHours) * time.Hour)

		return adjustedHarvestTime, nil
	}

	return baseHarvestTime, nil
}

// 私有方法

// initializeDefaultTemplates 初始化默认模板
func (ps *PlantService) initializeDefaultTemplates() {
	// 初始化种子模板
	ps.initializeSeedTemplates()

	// 初始化土壤模板
	ps.initializeSoilTemplates()

	// 初始化肥料模板
	ps.initializeFertilizerTemplates()

	// 初始化工具模板
	ps.initializeToolTemplates()
}

// initializeSeedTemplates 初始化种子模板
func (ps *PlantService) initializeSeedTemplates() {
	// 小麦模板
	wheatTemplate := &SeedTemplate{
		SeedType:         SeedTypeWheat,
		OptimalSeasons:   []Season{SeasonSpring, SeasonAutumn},
		OptimalSoilTypes: []SoilType{SoilTypeLoam, SoilTypeSilt},
		MinTemperature:   5.0,
		MaxTemperature:   30.0,
		OptimalPH:        6.5,
		WaterTolerance:   0.8,
		NutrientNeeds:    map[string]float64{"nitrogen": 0.6, "phosphorus": 0.3, "potassium": 0.4},
		GrowthModifiers:  map[string]float64{"temperature": 1.2, "moisture": 1.1},
	}
	ps.RegisterSeedTemplate(SeedTypeWheat, wheatTemplate)

	// 玉米模板
	cornTemplate := &SeedTemplate{
		SeedType:         SeedTypeCorn,
		OptimalSeasons:   []Season{SeasonSummer},
		OptimalSoilTypes: []SoilType{SoilTypeLoam, SoilTypeSandy},
		MinTemperature:   15.0,
		MaxTemperature:   35.0,
		OptimalPH:        6.8,
		WaterTolerance:   0.9,
		NutrientNeeds:    map[string]float64{"nitrogen": 0.8, "phosphorus": 0.4, "potassium": 0.6},
		GrowthModifiers:  map[string]float64{"temperature": 1.3, "sunlight": 1.2},
	}
	ps.RegisterSeedTemplate(SeedTypeCorn, cornTemplate)

	// 番茄模板
	tomatoTemplate := &SeedTemplate{
		SeedType:         SeedTypeTomato,
		OptimalSeasons:   []Season{SeasonSpring, SeasonSummer},
		OptimalSoilTypes: []SoilType{SoilTypeLoam},
		MinTemperature:   18.0,
		MaxTemperature:   28.0,
		OptimalPH:        6.2,
		WaterTolerance:   0.7,
		NutrientNeeds:    map[string]float64{"nitrogen": 0.5, "phosphorus": 0.6, "potassium": 0.8},
		GrowthModifiers:  map[string]float64{"temperature": 1.1, "moisture": 1.3},
	}
	ps.RegisterSeedTemplate(SeedTypeTomato, tomatoTemplate)
}

// initializeSoilTemplates 初始化土壤模板
func (ps *PlantService) initializeSoilTemplates() {
	// 壤土模板
	loamTemplate := &SoilTemplate{
		SoilType:         SoilTypeLoam,
		BaseProductivity: 1.2,
		WaterRetention:   0.7,
		NutrientCapacity: 0.8,
		DrainageRate:     0.6,
		OptimalPHRange:   PHRange{Min: 6.0, Max: 7.5},
		SuitableCrops:    []SeedType{SeedTypeWheat, SeedTypeCorn, SeedTypeTomato, SeedTypeCabbage},
	}
	ps.RegisterSoilTemplate(SoilTypeLoam, loamTemplate)

	// 沙土模板
	sandyTemplate := &SoilTemplate{
		SoilType:         SoilTypeSandy,
		BaseProductivity: 0.8,
		WaterRetention:   0.3,
		NutrientCapacity: 0.4,
		DrainageRate:     0.9,
		OptimalPHRange:   PHRange{Min: 5.5, Max: 7.0},
		SuitableCrops:    []SeedType{SeedTypePotato, SeedTypeCarrot},
	}
	ps.RegisterSoilTemplate(SoilTypeSandy, sandyTemplate)

	// 粘土模板
	clayTemplate := &SoilTemplate{
		SoilType:         SoilTypeClay,
		BaseProductivity: 0.9,
		WaterRetention:   0.9,
		NutrientCapacity: 0.9,
		DrainageRate:     0.2,
		OptimalPHRange:   PHRange{Min: 6.5, Max: 8.0},
		SuitableCrops:    []SeedType{SeedTypeRice},
	}
	ps.RegisterSoilTemplate(SoilTypeClay, clayTemplate)
}

// initializeFertilizerTemplates 初始化肥料模板
func (ps *PlantService) initializeFertilizerTemplates() {
	// 有机肥模板
	organicTemplate := &FertilizerTemplate{
		FertilizerType:  FertilizerTypeOrganic,
		BaseEffect:      1.0,
		EffectDuration:  72 * time.Hour,
		SoilMultipliers: map[SoilType]float64{SoilTypeLoam: 1.2, SoilTypeSandy: 1.1, SoilTypeClay: 1.0},
		CropMultipliers: map[SeedType]float64{SeedTypeWheat: 1.1, SeedTypeCorn: 1.0, SeedTypeTomato: 1.2},
		NutrientProfile: map[string]float64{"nitrogen": 0.3, "phosphorus": 0.2, "potassium": 0.25, "organic": 0.8},
	}
	ps.RegisterFertilizerTemplate(FertilizerTypeOrganic, organicTemplate)

	// 化学肥料模板
	chemicalTemplate := &FertilizerTemplate{
		FertilizerType:  FertilizerTypeChemical,
		BaseEffect:      1.5,
		EffectDuration:  48 * time.Hour,
		SoilMultipliers: map[SoilType]float64{SoilTypeLoam: 1.0, SoilTypeSandy: 1.3, SoilTypeClay: 0.9},
		CropMultipliers: map[SeedType]float64{SeedTypeWheat: 1.3, SeedTypeCorn: 1.4, SeedTypeTomato: 1.1},
		NutrientProfile: map[string]float64{"nitrogen": 0.6, "phosphorus": 0.4, "potassium": 0.5, "organic": 0.1},
	}
	ps.RegisterFertilizerTemplate(FertilizerTypeChemical, chemicalTemplate)
}

// initializeToolTemplates 初始化工具模板
func (ps *PlantService) initializeToolTemplates() {
	// 锄头模板
	hoeTemplate := &ToolTemplate{
		ToolType:               ToolTypeHoe,
		BaseEfficiency:         1.0,
		OperationMultipliers:   map[string]float64{"soil_preparation": 1.5, "weeding": 1.3, "planting": 1.1},
		DurabilityConsumption:  5.0,
		MaintenanceRequirement: 0.1,
	}
	ps.RegisterToolTemplate(ToolTypeHoe, hoeTemplate)

	// 洒水壶模板
	wateringCanTemplate := &ToolTemplate{
		ToolType:               ToolTypeWateringCan,
		BaseEfficiency:         1.0,
		OperationMultipliers:   map[string]float64{"watering": 1.8, "fertilizing": 1.2},
		DurabilityConsumption:  3.0,
		MaintenanceRequirement: 0.05,
	}
	ps.RegisterToolTemplate(ToolTypeWateringCan, wateringCanTemplate)

	// 收割机模板
	harvesterTemplate := &ToolTemplate{
		ToolType:               ToolTypeHarvester,
		BaseEfficiency:         2.0,
		OperationMultipliers:   map[string]float64{"harvesting": 2.5, "processing": 1.5},
		DurabilityConsumption:  8.0,
		MaintenanceRequirement: 0.2,
	}
	ps.RegisterToolTemplate(ToolTypeHarvester, harvesterTemplate)
}

// initializeDefaultRules 初始化默认规则
func (ps *PlantService) initializeDefaultRules() {
	// 添加基础生长规则
	ps.AddGrowthRule(&GrowthRule{
		Name:        "optimal_temperature",
		Description: "最适温度加速生长",
		Condition:   "temperature_in_range",
		Multiplier:  1.2,
		Priority:    1,
	})

	ps.AddGrowthRule(&GrowthRule{
		Name:        "adequate_water",
		Description: "充足水分促进生长",
		Condition:   "water_level_high",
		Multiplier:  1.15,
		Priority:    2,
	})

	ps.AddGrowthRule(&GrowthRule{
		Name:        "nutrient_deficiency",
		Description: "营养不足减缓生长",
		Condition:   "nutrient_level_low",
		Multiplier:  0.7,
		Priority:    3,
	})

	// 添加病虫害规则
	ps.pestDiseaseRules = append(ps.pestDiseaseRules, &PestDiseaseRule{
		Name:        "aphid_infestation",
		Description: "蚜虫侵害",
		TriggerConditions: map[string]interface{}{
			"temperature_min": 20.0,
			"humidity_min":    70.0,
			"season":          []Season{SeasonSpring, SeasonSummer},
		},
		Probability:   0.15,
		Severity:      "medium",
		AffectedCrops: []SeedType{SeedTypeTomato, SeedTypeCabbage},
	})

	// 添加质量规则
	ps.qualityRules = append(ps.qualityRules, &QualityRule{
		Name:        "perfect_conditions",
		Description: "完美条件提升品质",
		Conditions: map[string]interface{}{
			"soil_quality_min": 80.0,
			"health_min":       90.0,
			"care_quality_min": 85.0,
		},
		QualityBonus: 15.0,
	})
}

// initializeSeasonalEffects 初始化季节效果
func (ps *PlantService) initializeSeasonalEffects() {
	// 春季效果
	ps.seasonalEffects[SeasonSpring] = map[SeedType]float64{
		SeedTypeWheat:      1.2,
		SeedTypeTomato:     1.3,
		SeedTypeCarrot:     1.1,
		SeedTypeCabbage:    1.2,
		SeedTypeStrawberry: 1.4,
	}

	// 夏季效果
	ps.seasonalEffects[SeasonSummer] = map[SeedType]float64{
		SeedTypeCorn:       1.4,
		SeedTypeTomato:     1.2,
		SeedTypeStrawberry: 1.1,
		SeedTypeApple:      1.1,
		SeedTypeOrange:     1.2,
	}

	// 秋季效果
	ps.seasonalEffects[SeasonAutumn] = map[SeedType]float64{
		SeedTypeWheat:  1.1,
		SeedTypeRice:   1.2,
		SeedTypePotato: 1.3,
		SeedTypeApple:  1.4,
		SeedTypeOrange: 1.3,
	}

	// 冬季效果
	ps.seasonalEffects[SeasonWinter] = map[SeedType]float64{
		SeedTypeCabbage: 0.9,
		SeedTypeCarrot:  0.8,
	}
}

// initializeClimateZones 初始化气候区域
func (ps *PlantService) initializeClimateZones() {
	// 温带气候
	temperateZone := &ClimateZone{
		ZoneID:                     "temperate",
		Name:                       "温带气候",
		Description:                "四季分明，适合多种作物",
		TemperatureRange:           TemperatureRange{Min: -5, Max: 35, Average: 15},
		HumidityRange:              HumidityRange{Min: 40, Max: 80, Average: 60},
		SuitableCrops:              []SeedType{SeedTypeWheat, SeedTypeCorn, SeedTypeTomato, SeedTypePotato},
		SeasonalModifiers:          map[Season]float64{SeasonSpring: 1.2, SeasonSummer: 1.1, SeasonAutumn: 1.0, SeasonWinter: 0.7},
		WaterRequirementMultiplier: 1.0,
	}
	ps.RegisterClimateZone("temperate", temperateZone)

	// 热带气候
	tropicalZone := &ClimateZone{
		ZoneID:                     "tropical",
		Name:                       "热带气候",
		Description:                "高温多湿，适合热带作物",
		TemperatureRange:           TemperatureRange{Min: 20, Max: 40, Average: 28},
		HumidityRange:              HumidityRange{Min: 70, Max: 95, Average: 85},
		SuitableCrops:              []SeedType{SeedTypeRice, SeedTypeOrange, SeedTypeStrawberry},
		SeasonalModifiers:          map[Season]float64{SeasonSpring: 1.1, SeasonSummer: 1.3, SeasonAutumn: 1.1, SeasonWinter: 1.0},
		WaterRequirementMultiplier: 1.3,
	}
	ps.RegisterClimateZone("tropical", tropicalZone)
}

// 辅助方法

// getOptimalSeasonsForSeed 获取种子的最佳季节
func (ps *PlantService) getOptimalSeasonsForSeed(seedType SeedType) []Season {
	template := ps.GetSeedTemplate(seedType)
	if template != nil {
		return template.OptimalSeasons
	}

	// 默认春季和夏季
	return []Season{SeasonSpring, SeasonSummer}
}

// calculateNextSeasonTime 计算下一个季节时间
func (ps *PlantService) calculateNextSeasonTime(currentTime time.Time, targetSeason Season) time.Time {
	year := currentTime.Year()

	switch targetSeason {
	case SeasonSpring:
		return time.Date(year, 3, 1, 0, 0, 0, 0, currentTime.Location())
	case SeasonSummer:
		return time.Date(year, 6, 1, 0, 0, 0, 0, currentTime.Location())
	case SeasonAutumn:
		return time.Date(year, 9, 1, 0, 0, 0, 0, currentTime.Location())
	case SeasonWinter:
		return time.Date(year, 12, 1, 0, 0, 0, 0, currentTime.Location())
	default:
		return currentTime
	}
}

// applyGrowthRules 应用生长规则
func (ps *PlantService) applyGrowthRules(crop *Crop, baseGrowthRate float64) float64 {
	actualRate := baseGrowthRate

	for _, rule := range ps.growthRules {
		if rule.AppliesTo(crop) {
			actualRate *= rule.Multiplier
		}
	}

	return actualRate
}

// getSeasonalEffect 获取季节效果
func (ps *PlantService) getSeasonalEffect(season Season, seedType SeedType) float64 {
	if effects, exists := ps.seasonalEffects[season]; exists {
		if effect, exists := effects[seedType]; exists {
			return effect
		}
	}
	return 1.0 // 默认无效果
}

// isSeasonSuitableForSeed 检查季节是否适合种子
func (ps *PlantService) isSeasonSuitableForSeed(seedType SeedType, season Season) bool {
	optimalSeasons := ps.getOptimalSeasonsForSeed(seedType)
	for _, optimalSeason := range optimalSeasons {
		if optimalSeason == season {
			return true
		}
	}
	return false
}

// applyQualityRules 应用质量规则
func (ps *PlantService) applyQualityRules(crop *Crop, soil *Soil, baseScore float64) float64 {
	adjustedScore := baseScore

	for _, rule := range ps.qualityRules {
		if rule.AppliesTo(crop, soil) {
			adjustedScore += rule.QualityBonus
		}
	}

	return adjustedScore
}

// scoreToQuality 分数转换为品质
func (ps *PlantService) scoreToQuality(score float64) CropQuality {
	if score >= 95 {
		return CropQualityLegendary
	} else if score >= 85 {
		return CropQualityEpic
	} else if score >= 75 {
		return CropQualityRare
	} else if score >= 65 {
		return CropQualityUncommon
	} else {
		return CropQualityCommon
	}
}

// getGrowthStageWaterMultiplier 获取生长阶段水分倍率
func (ps *PlantService) getGrowthStageWaterMultiplier(stage GrowthStage) float64 {
	switch stage {
	case GrowthStageSeed:
		return 0.8
	case GrowthStageSeedling:
		return 1.2
	case GrowthStageGrowing:
		return 1.5
	case GrowthStageFlowering:
		return 1.3
	case GrowthStageMature:
		return 1.0
	default:
		return 1.0
	}
}

// getGrowthStageNutrientMultiplier 获取生长阶段营养倍率
func (ps *PlantService) getGrowthStageNutrientMultiplier(stage GrowthStage) float64 {
	switch stage {
	case GrowthStageSeed:
		return 0.5
	case GrowthStageSeedling:
		return 1.3
	case GrowthStageGrowing:
		return 1.8
	case GrowthStageFlowering:
		return 1.4
	case GrowthStageMature:
		return 0.8
	default:
		return 1.0
	}
}

// calculateCurrentEnvironmentMultiplier 计算当前环境倍率
func (ps *PlantService) calculateCurrentEnvironmentMultiplier(crop *Crop) float64 {
	multiplier := 1.0

	// 健康状态影响
	healthMultiplier := crop.GetHealthScore() / 100.0
	multiplier *= healthMultiplier

	// 水分和营养影响
	if crop.GetWaterLevel() > 70.0 && crop.GetNutrientLevel() > 70.0 {
		multiplier *= 1.2
	} else if crop.GetWaterLevel() < 30.0 || crop.GetNutrientLevel() < 30.0 {
		multiplier *= 0.8
	}

	return multiplier
}

// getDefaultClimateZone 获取默认气候区域
func (ps *PlantService) getDefaultClimateZone() *ClimateZone {
	return ps.GetClimateZone("temperate")
}

// 模板和规则结构体定义

// SeedTemplate 种子模板
type SeedTemplate struct {
	SeedType         SeedType
	OptimalSeasons   []Season
	OptimalSoilTypes []SoilType
	MinTemperature   float64
	MaxTemperature   float64
	OptimalPH        float64
	WaterTolerance   float64
	NutrientNeeds    map[string]float64
	GrowthModifiers  map[string]float64
}

// SoilTemplate 土壤模板
type SoilTemplate struct {
	SoilType         SoilType
	BaseProductivity float64
	WaterRetention   float64
	NutrientCapacity float64
	DrainageRate     float64
	OptimalPHRange   PHRange
	SuitableCrops    []SeedType
}

// FertilizerTemplate 肥料模板
type FertilizerTemplate struct {
	FertilizerType  FertilizerType
	BaseEffect      float64
	EffectDuration  time.Duration
	SoilMultipliers map[SoilType]float64
	CropMultipliers map[SeedType]float64
	NutrientProfile map[string]float64
}

// GetBaseEffect 获取基础效果
func (ft *FertilizerTemplate) GetBaseEffect() float64 {
	return ft.BaseEffect
}

// GetSoilMultiplier 获取土壤倍率
func (ft *FertilizerTemplate) GetSoilMultiplier(soilType SoilType) float64 {
	if multiplier, exists := ft.SoilMultipliers[soilType]; exists {
		return multiplier
	}
	return 1.0
}

// GetCropMultiplier 获取作物倍率
func (ft *FertilizerTemplate) GetCropMultiplier(seedType SeedType) float64 {
	if multiplier, exists := ft.CropMultipliers[seedType]; exists {
		return multiplier
	}
	return 1.0
}

// GetEffectDuration 获取效果持续时间
func (ft *FertilizerTemplate) GetEffectDuration() time.Duration {
	return ft.EffectDuration
}

// ToolTemplate 工具模板
type ToolTemplate struct {
	ToolType               ToolType
	BaseEfficiency         float64
	OperationMultipliers   map[string]float64
	DurabilityConsumption  float64
	MaintenanceRequirement float64
}

// GetOperationMultiplier 获取操作倍率
func (tt *ToolTemplate) GetOperationMultiplier(operation string) float64 {
	if multiplier, exists := tt.OperationMultipliers[operation]; exists {
		return multiplier
	}
	return 1.0
}

// GrowthRule 生长规则
type GrowthRule struct {
	Name        string
	Description string
	Condition   string
	Multiplier  float64
	Priority    int
}

// AppliesTo 检查规则是否适用
func (gr *GrowthRule) AppliesTo(crop *Crop) bool {
	switch gr.Condition {
	case "temperature_in_range":
		return true // 简化实现
	case "water_level_high":
		return crop.GetWaterLevel() > 70.0
	case "nutrient_level_low":
		return crop.GetNutrientLevel() < 30.0
	default:
		return false
	}
}

// PestDiseaseRule 病虫害规则
type PestDiseaseRule struct {
	Name              string
	Description       string
	TriggerConditions map[string]interface{}
	Probability       float64
	Severity          string
	AffectedCrops     []SeedType
}

// ShouldTrigger 检查是否应该触发
func (pdr *PestDiseaseRule) ShouldTrigger(crop *Crop, soil *Soil, season Season, climateZone string) bool {
	// 简化实现
	for _, affectedCrop := range pdr.AffectedCrops {
		if affectedCrop == crop.GetSeedType() {
			return true
		}
	}
	return false
}

// CreateEvent 创建事件
func (pdr *PestDiseaseRule) CreateEvent(crop *Crop) *PestDiseaseEvent {
	return &PestDiseaseEvent{
		Name:        pdr.Name,
		Description: pdr.Description,
		Severity:    pdr.Severity,
		CropID:      crop.GetID(),
		OccurredAt:  time.Now(),
	}
}

// QualityRule 质量规则
type QualityRule struct {
	Name         string
	Description  string
	Conditions   map[string]interface{}
	QualityBonus float64
}

// AppliesTo 检查规则是否适用
func (qr *QualityRule) AppliesTo(crop *Crop, soil *Soil) bool {
	// 简化实现
	if soilQualityMin, exists := qr.Conditions["soil_quality_min"]; exists {
		if soil.GetQualityScore() < soilQualityMin.(float64) {
			return false
		}
	}

	if healthMin, exists := qr.Conditions["health_min"]; exists {
		if crop.GetHealthScore() < healthMin.(float64) {
			return false
		}
	}

	return true
}

// ClimateZone 气候区域
type ClimateZone struct {
	ZoneID                     string
	Name                       string
	Description                string
	TemperatureRange           TemperatureRange
	HumidityRange              HumidityRange
	SuitableCrops              []SeedType
	SeasonalModifiers          map[Season]float64
	WaterRequirementMultiplier float64
}

// IsSuitableFor 检查是否适合作物
func (cz *ClimateZone) IsSuitableFor(seedType SeedType) bool {
	for _, suitableCrop := range cz.SuitableCrops {
		if suitableCrop == seedType {
			return true
		}
	}
	return false
}

// GetWaterRequirementMultiplier 获取水分需求倍率
func (cz *ClimateZone) GetWaterRequirementMultiplier() float64 {
	return cz.WaterRequirementMultiplier
}

// TemperatureRange 温度范围
type TemperatureRange struct {
	Min     float64
	Max     float64
	Average float64
}

// HumidityRange 湿度范围
type HumidityRange struct {
	Min     float64
	Max     float64
	Average float64
}

// PHRange pH范围
type PHRange struct {
	Min float64
	Max float64
}

// FertilizerEffect 肥料效果
type FertilizerEffect struct {
	FertilityBoost   float64
	NutrientBoost    float64
	GrowthSpeedBoost float64
	Duration         time.Duration
}

// PestDiseaseEvent 病虫害事件
type PestDiseaseEvent struct {
	Name        string
	Description string
	Severity    string
	CropID      string
	OccurredAt  time.Time
}

// 错误定义
var (
	ErrInvalidCrop                = fmt.Errorf("invalid crop")
	ErrInvalidSoil                = fmt.Errorf("invalid soil")
	ErrClimateNotSuitable         = fmt.Errorf("climate not suitable")
	ErrFertilizerTemplateNotFound = fmt.Errorf("fertilizer template not found")
	ErrToolTemplateNotFound       = fmt.Errorf("tool template not found")
)
