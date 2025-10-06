package plant

import (
	"errors"
	"fmt"
	"time"
)

// FarmAggregate 农场聚合根
type FarmAggregate struct {
	farmID          string
	sceneID         string
	owner           string
	name            string
	description     string
	size            FarmSize
	soil            *Soil
	crops           map[string]*Crop
	plots           map[string]*Plot
	tools           []*FarmTool
	resources       *FarmResources
	statistics      *FarmStatistics
	climateZone     string
	seasonModifier  *SeasonModifier
	automation      *AutomationSettings
	lastUpdateTime  time.Time
	lastHarvestTime time.Time
	createdAt       time.Time
	updatedAt       time.Time
	version         int
}

// NewFarmAggregate 创建农场聚合根
func NewFarmAggregate(farmID, sceneID, owner, name string, size FarmSize) *FarmAggregate {
	now := time.Now()
	return &FarmAggregate{
		farmID:          farmID,
		sceneID:         sceneID,
		owner:           owner,
		name:            name,
		size:            size,
		soil:            NewSoil(SoilTypeLoam, 50.0, 7.0, 2.0), // 默认壤土
		crops:           make(map[string]*Crop),
		plots:           make(map[string]*Plot),
		tools:           make([]*FarmTool, 0),
		resources:       NewFarmResources(),
		statistics:      NewFarmStatistics(),
		climateZone:     "temperate", // 默认温带
		seasonModifier:  NewSeasonModifier(),
		automation:      NewAutomationSettings(),
		lastUpdateTime:  now,
		lastHarvestTime: now,
		createdAt:       now,
		updatedAt:       now,
		version:         1,
	}
}

// GetFarmID 获取农场ID
func (f *FarmAggregate) GetFarmID() string {
	return f.farmID
}

// GetSceneID 获取场景ID
func (f *FarmAggregate) GetSceneID() string {
	return f.sceneID
}

// GetOwner 获取所有者
func (f *FarmAggregate) GetOwner() string {
	return f.owner
}

// GetName 获取农场名称
func (f *FarmAggregate) GetName() string {
	return f.name
}

// SetName 设置农场名称
func (f *FarmAggregate) SetName(name string) error {
	if name == "" {
		return ErrInvalidFarmName
	}

	f.name = name
	f.updateVersion()
	return nil
}

// GetDescription 获取描述
func (f *FarmAggregate) GetDescription() string {
	return f.description
}

// SetDescription 设置描述
func (f *FarmAggregate) SetDescription(description string) {
	f.description = description
	f.updateVersion()
}

// GetSize 获取农场大小
func (f *FarmAggregate) GetSize() FarmSize {
	return f.size
}

// ExpandFarm 扩展农场
func (f *FarmAggregate) ExpandFarm(newSize FarmSize) error {
	if newSize <= f.size {
		return ErrInvalidFarmSize
	}

	// 检查扩展条件
	if !f.canExpand(newSize) {
		return ErrFarmExpansionNotAllowed
	}

	f.size = newSize
	f.updateVersion()
	return nil
}

// GetSoil 获取土壤
func (f *FarmAggregate) GetSoil() *Soil {
	return f.soil
}

// UpdateSoil 更新土壤
func (f *FarmAggregate) UpdateSoil(soil *Soil) error {
	if soil == nil {
		return ErrInvalidSoil
	}

	f.soil = soil
	f.updateVersion()
	return nil
}

// ImproveSoil 改良土壤
func (f *FarmAggregate) ImproveSoil(fertilizer *Fertilizer) error {
	if fertilizer == nil {
		return ErrInvalidFertilizer
	}

	// 检查资源是否足够
	if !f.resources.HasEnoughFertilizer(fertilizer.GetType(), fertilizer.GetAmount()) {
		return ErrInsufficientResources
	}

	// 应用肥料效果
	f.soil.ApplyFertilizer(fertilizer)

	// 消耗资源
	f.resources.ConsumeFertilizer(fertilizer.GetType(), fertilizer.GetAmount())

	// 更新统计
	f.statistics.AddFertilizerUsage(fertilizer.GetType(), fertilizer.GetAmount())

	f.updateVersion()
	return nil
}

// PlantCrop 种植作物
func (f *FarmAggregate) PlantCrop(plotID string, seedType SeedType, quantity int) error {
	if plotID == "" {
		return ErrInvalidPlotID
	}

	if !seedType.IsValid() {
		return ErrInvalidSeedType
	}

	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	// 检查地块是否存在
	plot, exists := f.plots[plotID]
	if !exists {
		return ErrPlotNotFound
	}

	// 检查地块是否可用
	if !plot.IsAvailable() {
		return ErrPlotNotAvailable
	}

	// 检查种子资源
	if !f.resources.HasEnoughSeeds(seedType, quantity) {
		return ErrInsufficientSeeds
	}

	// 检查土壤适宜性
	if !f.soil.IsSuitableFor(seedType) {
		return ErrSoilNotSuitable
	}

	// 创建作物
	crop := NewCrop(generateCropID(), seedType, quantity, f.soil, f.climateZone)

	// 应用季节修正
	f.seasonModifier.ApplyToCrop(crop)

	// 种植作物
	plot.PlantCrop(crop)
	f.crops[crop.GetID()] = crop

	// 消耗种子
	f.resources.ConsumeSeeds(seedType, quantity)

	// 更新统计
	f.statistics.AddPlantingActivity(seedType, quantity)

	f.updateVersion()
	return nil
}

// WaterCrops 浇水
func (f *FarmAggregate) WaterCrops(plotIDs []string, waterAmount float64) error {
	if len(plotIDs) == 0 {
		return ErrInvalidPlotID
	}

	if waterAmount <= 0 {
		return ErrInvalidWaterAmount
	}

	// 检查水资源
	totalWaterNeeded := waterAmount * float64(len(plotIDs))
	if !f.resources.HasEnoughWater(totalWaterNeeded) {
		return ErrInsufficientWater
	}

	for _, plotID := range plotIDs {
		plot, exists := f.plots[plotID]
		if !exists {
			continue
		}

		// 浇水
		if crop := plot.GetCrop(); crop != nil {
			crop.Water(waterAmount)
		}

		// 更新土壤湿度
		f.soil.AddMoisture(waterAmount)
	}

	// 消耗水资源
	f.resources.ConsumeWater(totalWaterNeeded)

	// 更新统计
	f.statistics.AddWateringActivity(len(plotIDs), totalWaterNeeded)

	f.updateVersion()
	return nil
}

// HarvestCrop 收获作物
func (f *FarmAggregate) HarvestCrop(cropID string) (*HarvestResult, error) {
	if cropID == "" {
		return nil, ErrInvalidCropID
	}

	crop, exists := f.crops[cropID]
	if !exists {
		return nil, ErrCropNotFound
	}

	// 检查作物是否可收获
	if !crop.IsHarvestable() {
		return nil, ErrCropNotHarvestable
	}

	// 计算收获量
	yield := f.calculateYield(crop)
	quality := f.calculateQuality(crop)

	// 创建收获结果
	result := &HarvestResult{
		CropID:      cropID,
		SeedType:    crop.GetSeedType(),
		Yield:       yield,
		Quality:     quality,
		HarvestTime: time.Now(),
		Experience:  f.calculateExperience(crop, yield, quality),
	}

	// 添加收获物到资源
	f.resources.AddHarvest(crop.GetSeedType(), yield, quality)

	// 更新统计
	f.statistics.AddHarvestActivity(crop.GetSeedType(), yield, quality)

	// 移除作物
	delete(f.crops, cropID)

	// 释放地块
	for _, plot := range f.plots {
		if plot.GetCrop() != nil && plot.GetCrop().GetID() == cropID {
			plot.ClearCrop()
			break
		}
	}

	f.lastHarvestTime = time.Now()
	f.updateVersion()

	return result, nil
}

// UpdateCrops 更新所有作物
func (f *FarmAggregate) UpdateCrops() error {
	now := time.Now()

	for _, crop := range f.crops {
		// 更新作物生长
		crop.Update(now)

		// 应用环境影响
		f.applyEnvironmentalEffects(crop)

		// 检查病虫害
		f.checkPestsAndDiseases(crop)
	}

	f.lastUpdateTime = now
	f.updateVersion()
	return nil
}

// GetCrops 获取所有作物
func (f *FarmAggregate) GetCrops() map[string]*Crop {
	return f.crops
}

// GetCrop 获取指定作物
func (f *FarmAggregate) GetCrop(cropID string) *Crop {
	return f.crops[cropID]
}

// GetPlots 获取所有地块
func (f *FarmAggregate) GetPlots() map[string]*Plot {
	return f.plots
}

// AddPlot 添加地块
func (f *FarmAggregate) AddPlot(plot *Plot) error {
	if plot == nil {
		return ErrInvalidPlot
	}

	// 检查地块数量限制
	if len(f.plots) >= f.size.GetMaxPlots() {
		return ErrMaxPlotsReached
	}

	f.plots[plot.GetID()] = plot
	f.updateVersion()
	return nil
}

// RemovePlot 移除地块
func (f *FarmAggregate) RemovePlot(plotID string) error {
	plot, exists := f.plots[plotID]
	if !exists {
		return ErrPlotNotFound
	}

	// 检查地块是否有作物
	if plot.HasCrop() {
		return ErrPlotHasCrop
	}

	delete(f.plots, plotID)
	f.updateVersion()
	return nil
}

// GetTools 获取农具
func (f *FarmAggregate) GetTools() []*FarmTool {
	return f.tools
}

// AddTool 添加农具
func (f *FarmAggregate) AddTool(tool *FarmTool) error {
	if tool == nil {
		return ErrInvalidTool
	}

	f.tools = append(f.tools, tool)
	f.updateVersion()
	return nil
}

// UseTool 使用农具
func (f *FarmAggregate) UseTool(toolID string, targetID string) error {
	tool := f.findTool(toolID)
	if tool == nil {
		return ErrToolNotFound
	}

	if !tool.IsUsable() {
		return ErrToolNotUsable
	}

	// 使用农具
	effect := tool.Use()

	// 应用效果
	f.applyToolEffect(effect, targetID)

	// 更新统计
	f.statistics.AddToolUsage(tool.GetType())

	f.updateVersion()
	return nil
}

// GetResources 获取资源
func (f *FarmAggregate) GetResources() *FarmResources {
	return f.resources
}

// GetStatistics 获取统计信息
func (f *FarmAggregate) GetStatistics() *FarmStatistics {
	return f.statistics
}

// GetClimateZone 获取气候区域
func (f *FarmAggregate) GetClimateZone() string {
	return f.climateZone
}

// SetClimateZone 设置气候区域
func (f *FarmAggregate) SetClimateZone(zone string) {
	f.climateZone = zone
	f.updateVersion()
}

// GetSeasonModifier 获取季节修正
func (f *FarmAggregate) GetSeasonModifier() *SeasonModifier {
	return f.seasonModifier
}

// UpdateSeasonModifier 更新季节修正
func (f *FarmAggregate) UpdateSeasonModifier(modifier *SeasonModifier) {
	f.seasonModifier = modifier
	f.updateVersion()
}

// GetAutomationSettings 获取自动化设置
func (f *FarmAggregate) GetAutomationSettings() *AutomationSettings {
	return f.automation
}

// UpdateAutomationSettings 更新自动化设置
func (f *FarmAggregate) UpdateAutomationSettings(settings *AutomationSettings) {
	f.automation = settings
	f.updateVersion()
}

// GetLastUpdateTime 获取最后更新时间
func (f *FarmAggregate) GetLastUpdateTime() time.Time {
	return f.lastUpdateTime
}

// GetLastHarvestTime 获取最后收获时间
func (f *FarmAggregate) GetLastHarvestTime() time.Time {
	return f.lastHarvestTime
}

// GetCreatedAt 获取创建时间
func (f *FarmAggregate) GetCreatedAt() time.Time {
	return f.createdAt
}

// GetUpdatedAt 获取更新时间
func (f *FarmAggregate) GetUpdatedAt() time.Time {
	return f.updatedAt
}

// GetVersion 获取版本
func (f *FarmAggregate) GetVersion() int {
	return f.version
}

// CalculateFarmValue 计算农场价值
func (f *FarmAggregate) CalculateFarmValue() float64 {
	value := 0.0

	// 基础价值
	value += f.size.GetBaseValue()

	// 土壤价值
	value += f.soil.GetValue()

	// 作物价值
	for _, crop := range f.crops {
		value += crop.GetValue()
	}

	// 农具价值
	for _, tool := range f.tools {
		value += tool.GetValue()
	}

	// 资源价值
	value += f.resources.GetTotalValue()

	return value
}

// CalculateProductivity 计算生产力
func (f *FarmAggregate) CalculateProductivity() float64 {
	productivity := 1.0

	// 土壤影响
	productivity *= f.soil.GetProductivityMultiplier()

	// 季节影响
	productivity *= f.seasonModifier.GetProductivityMultiplier()

	// 农具影响
	for _, tool := range f.tools {
		if tool.IsActive() {
			productivity *= tool.GetProductivityBonus()
		}
	}

	return productivity
}

// GetFarmStatus 获取农场状态
func (f *FarmAggregate) GetFarmStatus() FarmStatus {
	if len(f.crops) == 0 {
		return FarmStatusIdle
	}

	// 检查是否有成熟的作物
	for _, crop := range f.crops {
		if crop.IsHarvestable() {
			return FarmStatusHarvestReady
		}
	}

	// 检查是否需要照料
	for _, crop := range f.crops {
		if crop.NeedsCare() {
			return FarmStatusNeedsCare
		}
	}

	return FarmStatusGrowing
}

// 私有方法

// canExpand 检查是否可以扩展
func (f *FarmAggregate) canExpand(newSize FarmSize) bool {
	// 检查资源是否足够
	expansionCost := newSize.GetExpansionCost(f.size)
	return f.resources.CanAfford(expansionCost)
}

// calculateYield 计算收获量
func (f *FarmAggregate) calculateYield(crop *Crop) int {
	baseYield := crop.GetBaseYield()
	multiplier := 1.0

	// 土壤影响
	multiplier *= f.soil.GetYieldMultiplier(crop.GetSeedType())

	// 季节影响
	multiplier *= f.seasonModifier.GetYieldMultiplier(crop.GetSeedType())

	// 照料质量影响
	multiplier *= crop.GetCareQualityMultiplier()

	return int(float64(baseYield) * multiplier)
}

// calculateQuality 计算品质
func (f *FarmAggregate) calculateQuality(crop *Crop) CropQuality {
	qualityScore := 0.0

	// 土壤质量影响
	qualityScore += f.soil.GetQualityScore()

	// 作物健康度影响
	qualityScore += crop.GetHealthScore()

	// 照料质量影响
	qualityScore += crop.GetCareQualityScore()

	// 转换为品质等级
	if qualityScore >= 90 {
		return CropQualityLegendary
	} else if qualityScore >= 80 {
		return CropQualityEpic
	} else if qualityScore >= 70 {
		return CropQualityRare
	} else if qualityScore >= 60 {
		return CropQualityUncommon
	} else {
		return CropQualityCommon
	}
}

// calculateExperience 计算经验值
func (f *FarmAggregate) calculateExperience(crop *Crop, yield int, quality CropQuality) int {
	baseExp := crop.GetBaseExperience()
	multiplier := 1.0

	// 产量影响
	multiplier += float64(yield) * 0.01

	// 品质影响
	multiplier += quality.GetExperienceMultiplier()

	return int(float64(baseExp) * multiplier)
}

// applyEnvironmentalEffects 应用环境影响
func (f *FarmAggregate) applyEnvironmentalEffects(crop *Crop) {
	// 土壤影响
	f.soil.ApplyToCrop(crop)

	// 季节影响
	f.seasonModifier.ApplyToCrop(crop)
}

// checkPestsAndDiseases 检查病虫害
func (f *FarmAggregate) checkPestsAndDiseases(crop *Crop) {
	// 简化的病虫害检查逻辑
	if crop.GetHealthScore() < 50 {
		// 可能有病虫害，需要处理
		crop.AddProblem("pest_infestation")
	}
}

// findTool 查找农具
func (f *FarmAggregate) findTool(toolID string) *FarmTool {
	for _, tool := range f.tools {
		if tool.GetID() == toolID {
			return tool
		}
	}
	return nil
}

// applyToolEffect 应用农具效果
func (f *FarmAggregate) applyToolEffect(effect *ToolEffect, targetID string) {
	switch effect.GetType() {
	case "soil_improvement":
		f.soil.ApplyImprovement(effect.GetValue())
	case "crop_growth_boost":
		if crop := f.crops[targetID]; crop != nil {
			crop.ApplyGrowthBoost(effect.GetValue())
		}
	case "pest_control":
		if crop := f.crops[targetID]; crop != nil {
			crop.RemoveProblem("pest_infestation")
		}
	}
}

// updateVersion 更新版本
func (f *FarmAggregate) updateVersion() {
	f.version++
	f.updatedAt = time.Now()
}

// generateCropID 生成作物ID
func generateCropID() string {
	return fmt.Sprintf("crop_%d", time.Now().UnixNano())
}

// HarvestResult 收获结果
type HarvestResult struct {
	CropID      string
	SeedType    SeedType
	Yield       int
	Quality     CropQuality
	HarvestTime time.Time
	Experience  int
}

// FarmStatus 农场状态
type FarmStatus int

const (
	FarmStatusIdle FarmStatus = iota + 1
	FarmStatusGrowing
	FarmStatusHarvestReady
	FarmStatusNeedsCare
	FarmStatusMaintenance
)

// String 返回状态字符串
func (fs FarmStatus) String() string {
	switch fs {
	case FarmStatusIdle:
		return "idle"
	case FarmStatusGrowing:
		return "growing"
	case FarmStatusHarvestReady:
		return "harvest_ready"
	case FarmStatusNeedsCare:
		return "needs_care"
	case FarmStatusMaintenance:
		return "maintenance"
	default:
		return "unknown"
	}
}

// 农场相关错误
var (
	ErrInvalidFarmName         = errors.New("invalid farm name")
	ErrInvalidFarmSize         = errors.New("invalid farm size")
	ErrFarmExpansionNotAllowed = errors.New("farm expansion not allowed")
	// 错误定义已移动到errors.go文件中
)
