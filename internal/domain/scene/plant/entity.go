package plant

import (
	"fmt"
	"time"
)

// Crop 作物实体
type Crop struct {
	ID                  string
	PlayerID            string
	SeedType            SeedType
	Quantity            int
	GrowthStage         GrowthStage
	HealthPoints        float64
	MaxHealthPoints     float64
	GrowthProgress      float64
	WaterLevel          float64
	NutrientLevel       float64
	PlantedTime         time.Time
	LastWateredTime     time.Time
	LastFertilizedTime  time.Time
	ExpectedHarvestTime time.Time
	SoilCondition       *Soil
	ClimateZone         string
	CareHistory         []*CareRecord
	Problems            []string
	Bonuses             []*GrowthBonus
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// NewCrop 创建作物
func NewCrop(id string, playerID string, seedType SeedType, quantity int, soil *Soil, climateZone string) *Crop {
	now := time.Now()
	growthDuration := seedType.GetGrowthDuration()

	return &Crop{
		ID:                  id,
		PlayerID:            playerID,
		SeedType:            seedType,
		Quantity:            quantity,
		GrowthStage:         GrowthStageSeed,
		HealthPoints:        100.0,
		MaxHealthPoints:     100.0,
		GrowthProgress:      0.0,
		WaterLevel:          50.0,
		NutrientLevel:       50.0,
		PlantedTime:         now,
		LastWateredTime:     now,
		LastFertilizedTime:  now,
		ExpectedHarvestTime: now.Add(growthDuration),
		SoilCondition:       soil,
		ClimateZone:         climateZone,
		CareHistory:         make([]*CareRecord, 0),
		Problems:            make([]string, 0),
		Bonuses:             make([]*GrowthBonus, 0),
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}

// GetID 获取ID
func (c *Crop) GetID() string {
	return c.ID
}

// GetPlayerID 获取玩家ID
func (c *Crop) GetPlayerID() string {
	return c.PlayerID
}

// GetSeedType 获取种子类型
func (c *Crop) GetSeedType() SeedType {
	return c.SeedType
}

// GetQuantity 获取数量
func (c *Crop) GetQuantity() int {
	return c.Quantity
}

// GetGrowthStage 获取生长阶段
func (c *Crop) GetGrowthStage() GrowthStage {
	return c.GrowthStage
}

// GetGrowthProgress 获取生长进度
func (c *Crop) GetGrowthProgress() float64 {
	return c.GrowthProgress
}

// GetHealthPoints 获取健康值
func (c *Crop) GetHealthPoints() float64 {
	return c.HealthPoints
}

// GetHealthScore 获取健康分数
func (c *Crop) GetHealthScore() float64 {
	return (c.HealthPoints / c.MaxHealthPoints) * 100
}

// GetWaterLevel 获取水分等级
func (c *Crop) GetWaterLevel() float64 {
	return c.WaterLevel
}

// GetNutrientLevel 获取营养等级
func (c *Crop) GetNutrientLevel() float64 {
	return c.NutrientLevel
}

// GetPlotID 获取地块ID (暂时返回空字符串，需要在Crop结构体中添加PlotID字段)
func (c *Crop) GetPlotID() string {
	return "" // TODO: Add PlotID field to Crop struct
}

// GetSeedID 获取种子ID (暂时返回空字符串，需要在Crop结构体中添加SeedID字段)
func (c *Crop) GetSeedID() string {
	return "" // TODO: Add SeedID field to Crop struct
}

// GetCropType 获取作物类型 (返回种子类型的字符串表示)
func (c *Crop) GetCropType() string {
	return c.SeedType.String()
}

// GetCurrentStage 获取当前阶段 (返回生长阶段的字符串表示)
func (c *Crop) GetCurrentStage() string {
	return c.GrowthStage.String()
}

// IsHarvestable 检查是否可收获
func (c *Crop) IsHarvestable() bool {
	return c.GrowthStage == GrowthStageMature && c.GrowthProgress >= 100.0
}

// NeedsCare 检查是否需要照料
func (c *Crop) NeedsCare() bool {
	return c.WaterLevel < 30.0 || c.NutrientLevel < 30.0 || c.HealthPoints < 70.0 || len(c.Problems) > 0
}

// Water 浇水
func (c *Crop) Water(amount float64) {
	c.WaterLevel += amount
	if c.WaterLevel > 100.0 {
		c.WaterLevel = 100.0
	}

	c.LastWateredTime = time.Now()
	c.addCareRecord("watering", amount)
	c.UpdatedAt = time.Now()
}

// Fertilize 施肥
func (c *Crop) Fertilize(fertilizer *Fertilizer) {
	c.NutrientLevel += fertilizer.GetNutrientValue()
	if c.NutrientLevel > 100.0 {
		c.NutrientLevel = 100.0
	}

	// 应用肥料的额外效果
	if bonus := fertilizer.GetGrowthBonus(); bonus != nil {
		c.AddBonus(bonus)
	}

	c.LastFertilizedTime = time.Now()
	c.addCareRecord("fertilizing", fertilizer.GetNutrientValue())
	c.UpdatedAt = time.Now()
}

// Update 更新作物状态
func (c *Crop) Update(currentTime time.Time) {
	// 计算时间差
	timeDiff := currentTime.Sub(c.UpdatedAt)
	hours := timeDiff.Hours()

	if hours <= 0 {
		return
	}

	// 更新生长进度
	c.updateGrowthProgress(hours)

	// 更新生长阶段
	c.updateGrowthStage()

	// 消耗水分和营养
	c.consumeResources(hours)

	// 更新健康状态
	c.updateHealth()

	// 处理奖励效果
	c.processGrowthBonuses(hours)

	// 检查问题
	c.checkForProblems()

	c.UpdatedAt = currentTime
}

// AddBonus 添加生长奖励
func (c *Crop) AddBonus(bonus *GrowthBonus) {
	c.Bonuses = append(c.Bonuses, bonus)
	c.UpdatedAt = time.Now()
}

// RemoveBonus 移除生长奖励
func (c *Crop) RemoveBonus(bonusID string) {
	for i, bonus := range c.Bonuses {
		if bonus.ID == bonusID {
			c.Bonuses = append(c.Bonuses[:i], c.Bonuses[i+1:]...)
			break
		}
	}
	c.UpdatedAt = time.Now()
}

// AddProblem 添加问题
func (c *Crop) AddProblem(problem string) {
	for _, p := range c.Problems {
		if p == problem {
			return // 问题已存在
		}
	}
	c.Problems = append(c.Problems, problem)
	c.UpdatedAt = time.Now()
}

// RemoveProblem 移除问题
func (c *Crop) RemoveProblem(problem string) {
	for i, p := range c.Problems {
		if p == problem {
			c.Problems = append(c.Problems[:i], c.Problems[i+1:]...)
			break
		}
	}
	c.UpdatedAt = time.Now()
}

// GetProblems 获取问题列表
func (c *Crop) GetProblems() []string {
	return c.Problems
}

// ApplyGrowthBoost 应用生长加速
func (c *Crop) ApplyGrowthBoost(multiplier float64) {
	bonus := &GrowthBonus{
		ID:         fmt.Sprintf("boost_%d", time.Now().UnixNano()),
		Type:       "growth_speed",
		Multiplier: multiplier,
		Duration:   24 * time.Hour, // 持续24小时
		StartTime:  time.Now(),
	}
	c.AddBonus(bonus)
}

// GetBaseYield 获取基础产量
func (c *Crop) GetBaseYield() int {
	return c.SeedType.GetBaseYield() * c.Quantity
}

// GetBaseExperience 获取基础经验
func (c *Crop) GetBaseExperience() int {
	return c.SeedType.GetBaseExperience() * c.Quantity
}

// GetValue 获取作物价值
func (c *Crop) GetValue() float64 {
	baseValue := c.SeedType.GetBaseValue() * float64(c.Quantity)

	// 生长进度影响
	progressMultiplier := c.GrowthProgress / 100.0

	// 健康状态影响
	healthMultiplier := c.GetHealthScore() / 100.0

	return baseValue * progressMultiplier * healthMultiplier
}

// GetCareQualityMultiplier 获取照料质量倍率
func (c *Crop) GetCareQualityMultiplier() float64 {
	if len(c.CareHistory) == 0 {
		return 1.0
	}

	// 基于照料历史计算质量倍率
	totalCare := 0.0
	for _, record := range c.CareHistory {
		totalCare += record.Quality
	}

	averageQuality := totalCare / float64(len(c.CareHistory))
	return 0.8 + (averageQuality/100.0)*0.4 // 0.8-1.2倍率
}

// GetCareQualityScore 获取照料质量分数
func (c *Crop) GetCareQualityScore() float64 {
	return c.GetCareQualityMultiplier() * 100.0
}

// 私有方法

// updateGrowthProgress 更新生长进度
func (c *Crop) updateGrowthProgress(hours float64) {
	baseGrowthRate := c.SeedType.GetGrowthRate()

	// 应用环境因素
	environmentMultiplier := c.calculateEnvironmentMultiplier()

	// 应用奖励效果
	bonusMultiplier := c.calculateBonusMultiplier()

	// 计算实际生长速度
	actualGrowthRate := baseGrowthRate * environmentMultiplier * bonusMultiplier

	// 更新进度
	c.GrowthProgress += actualGrowthRate * hours
	if c.GrowthProgress > 100.0 {
		c.GrowthProgress = 100.0
	}
}

// updateGrowthStage 更新生长阶段
func (c *Crop) updateGrowthStage() {
	if c.GrowthProgress >= 100.0 {
		c.GrowthStage = GrowthStageMature
	} else if c.GrowthProgress >= 75.0 {
		c.GrowthStage = GrowthStageFlowering
	} else if c.GrowthProgress >= 50.0 {
		c.GrowthStage = GrowthStageGrowing
	} else if c.GrowthProgress >= 25.0 {
		c.GrowthStage = GrowthStageSeedling
	} else {
		c.GrowthStage = GrowthStageSeed
	}
}

// consumeResources 消耗资源
func (c *Crop) consumeResources(hours float64) {
	// 水分消耗
	waterConsumption := c.SeedType.GetWaterConsumption() * hours
	c.WaterLevel -= waterConsumption
	if c.WaterLevel < 0 {
		c.WaterLevel = 0
	}

	// 营养消耗
	nutrientConsumption := c.SeedType.GetNutrientConsumption() * hours
	c.NutrientLevel -= nutrientConsumption
	if c.NutrientLevel < 0 {
		c.NutrientLevel = 0
	}
}

// updateHealth 更新健康状态
func (c *Crop) updateHealth() {
	// 基于水分和营养水平调整健康值
	if c.WaterLevel < 20.0 || c.NutrientLevel < 20.0 {
		c.HealthPoints -= 5.0 // 缺水或缺营养会降低健康值
	} else if c.WaterLevel > 80.0 && c.NutrientLevel > 80.0 {
		c.HealthPoints += 2.0 // 充足的水分和营养会恢复健康值
	}

	// 限制健康值范围
	if c.HealthPoints < 0 {
		c.HealthPoints = 0
	} else if c.HealthPoints > c.MaxHealthPoints {
		c.HealthPoints = c.MaxHealthPoints
	}
}

// processGrowthBonuses 处理生长奖励
func (c *Crop) processGrowthBonuses(hours float64) {
	// 移除过期的奖励
	now := time.Now()
	for i := len(c.Bonuses) - 1; i >= 0; i-- {
		bonus := c.Bonuses[i]
		if now.After(bonus.StartTime.Add(bonus.Duration)) {
			c.Bonuses = append(c.Bonuses[:i], c.Bonuses[i+1:]...)
		}
	}
}

// checkForProblems 检查问题
func (c *Crop) checkForProblems() {
	// 清除旧问题
	c.Problems = c.Problems[:0]

	// 检查缺水
	if c.WaterLevel < 20.0 {
		c.AddProblem("drought_stress")
	}

	// 检查缺营养
	if c.NutrientLevel < 20.0 {
		c.AddProblem("nutrient_deficiency")
	}

	// 检查健康状态
	if c.HealthPoints < 30.0 {
		c.AddProblem("poor_health")
	}

	// 检查过度浇水
	if c.WaterLevel > 95.0 {
		c.AddProblem("overwatering")
	}
}

// calculateEnvironmentMultiplier 计算环境倍率
func (c *Crop) calculateEnvironmentMultiplier() float64 {
	multiplier := 1.0

	// 土壤影响
	if c.SoilCondition != nil {
		multiplier *= c.SoilCondition.GetGrowthMultiplier(c.SeedType)
	}

	// 水分影响
	if c.WaterLevel < 30.0 {
		multiplier *= 0.7 // 缺水减慢生长
	} else if c.WaterLevel > 80.0 {
		multiplier *= 1.1 // 充足水分加速生长
	}

	// 营养影响
	if c.NutrientLevel < 30.0 {
		multiplier *= 0.8 // 缺营养减慢生长
	} else if c.NutrientLevel > 80.0 {
		multiplier *= 1.2 // 充足营养加速生长
	}

	return multiplier
}

// calculateBonusMultiplier 计算奖励倍率
func (c *Crop) calculateBonusMultiplier() float64 {
	multiplier := 1.0

	for _, bonus := range c.Bonuses {
		if bonus.IsActive() {
			multiplier *= bonus.Multiplier
		}
	}

	return multiplier
}

// addCareRecord 添加照料记录
func (c *Crop) addCareRecord(careType string, value float64) {
	record := &CareRecord{
		Type:      careType,
		Value:     value,
		Quality:   c.calculateCareQuality(careType, value),
		Timestamp: time.Now(),
	}

	c.CareHistory = append(c.CareHistory, record)

	// 限制历史记录数量
	if len(c.CareHistory) > 50 {
		c.CareHistory = c.CareHistory[1:]
	}
}

// calculateCareQuality 计算照料质量
func (c *Crop) calculateCareQuality(careType string, value float64) float64 {
	// 基于照料类型和数值计算质量分数
	switch careType {
	case "watering":
		if value >= 20.0 && value <= 30.0 {
			return 100.0 // 完美浇水量
		} else if value >= 10.0 && value <= 40.0 {
			return 80.0 // 良好浇水量
		} else {
			return 60.0 // 一般浇水量
		}
	case "fertilizing":
		if value >= 15.0 && value <= 25.0 {
			return 100.0 // 完美施肥量
		} else if value >= 10.0 && value <= 30.0 {
			return 80.0 // 良好施肥量
		} else {
			return 60.0 // 一般施肥量
		}
	default:
		return 70.0 // 默认质量
	}
}

// Plot 地块实体
type Plot struct {
	ID          string
	Name        string
	Size        PlotSize
	SoilType    SoilType
	Fertility   float64
	Moisture    float64
	Crop        *Crop
	IsAvailable bool
	LastUsed    time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewPlot 创建地块
func NewPlot(id, name string, size PlotSize, soilType SoilType) *Plot {
	now := time.Now()
	return &Plot{
		ID:          id,
		Name:        name,
		Size:        size,
		SoilType:    soilType,
		Fertility:   50.0, // 默认肥力
		Moisture:    30.0, // 默认湿度
		Crop:        nil,
		IsAvailable: true,
		LastUsed:    now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// GetID 获取ID
func (p *Plot) GetID() string {
	return p.ID
}

// GetName 获取名称
func (p *Plot) GetName() string {
	return p.Name
}

// GetSize 获取大小
func (p *Plot) GetSize() PlotSize {
	return p.Size
}

// GetSoilType 获取土壤类型
func (p *Plot) GetSoilType() SoilType {
	return p.SoilType
}

// GetFertility 获取肥力
func (p *Plot) GetFertility() float64 {
	return p.Fertility
}

// GetMoisture 获取湿度
func (p *Plot) GetMoisture() float64 {
	return p.Moisture
}

// GetCrop 获取作物
func (p *Plot) GetCrop() *Crop {
	return p.Crop
}

// GetCropID 获取作物ID
func (p *Plot) GetCropID() string {
	if p.Crop != nil {
		return p.Crop.GetID()
	}
	return ""
}

// HasCrop 检查是否有作物
func (p *Plot) HasCrop() bool {
	return p.Crop != nil
}

// PlantCrop 种植作物
func (p *Plot) PlantCrop(crop *Crop) error {
	if !p.IsAvailable {
		return fmt.Errorf("plot is not available")
	}

	p.Crop = crop
	p.IsAvailable = false
	p.LastUsed = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

// ClearCrop 清除作物
func (p *Plot) ClearCrop() {
	p.Crop = nil
	p.IsAvailable = true
	p.UpdatedAt = time.Now()
}

// FarmTool 农具实体
type FarmTool struct {
	ID            string
	Name          string
	Type          ToolType
	Level         int
	Durability    float64
	MaxDurability float64
	Efficiency    float64
	IsActive      bool
	LastUsed      time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewFarmTool 创建农具
func NewFarmTool(id, name string, toolType ToolType, level int) *FarmTool {
	now := time.Now()
	maxDurability := float64(100 + level*20) // 等级越高耐久越高

	return &FarmTool{
		ID:            id,
		Name:          name,
		Type:          toolType,
		Level:         level,
		Durability:    maxDurability,
		MaxDurability: maxDurability,
		Efficiency:    1.0 + float64(level)*0.1, // 等级越高效率越高
		IsActive:      true,
		LastUsed:      now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// GetID 获取ID
func (ft *FarmTool) GetID() string {
	return ft.ID
}

// GetName 获取名称
func (ft *FarmTool) GetName() string {
	return ft.Name
}

// GetType 获取类型
func (ft *FarmTool) GetType() ToolType {
	return ft.Type
}

// GetLevel 获取等级
func (ft *FarmTool) GetLevel() int {
	return ft.Level
}

// GetDurability 获取耐久度
func (ft *FarmTool) GetDurability() float64 {
	return ft.Durability
}

// GetEfficiency 获取效率
func (ft *FarmTool) GetEfficiency() float64 {
	return ft.Efficiency
}

// IsUsable 检查是否可用
func (ft *FarmTool) IsUsable() bool {
	return ft.IsActive && ft.Durability >= 10.0 // 至少需要10点耐久
}

// Use 使用农具
func (ft *FarmTool) Use() *ToolEffect {
	if !ft.IsUsable() {
		return nil
	}

	// 消耗耐久度
	ft.Durability -= 5.0
	if ft.Durability < 0 {
		ft.Durability = 0
	}

	ft.LastUsed = time.Now()
	ft.UpdatedAt = time.Now()

	// 返回工具效果
	return ft.Type.GetEffect(ft.Level, ft.Efficiency)
}

// Repair 修理农具
func (ft *FarmTool) Repair(amount float64) {
	ft.Durability += amount
	if ft.Durability > ft.MaxDurability {
		ft.Durability = ft.MaxDurability
	}
	ft.UpdatedAt = time.Now()
}

// Upgrade 升级农具
func (ft *FarmTool) Upgrade() {
	ft.Level++
	ft.MaxDurability += 20.0
	ft.Durability = ft.MaxDurability // 升级后恢复满耐久
	ft.Efficiency += 0.1
	ft.UpdatedAt = time.Now()
}

// GetValue 获取价值
func (ft *FarmTool) GetValue() float64 {
	baseValue := ft.Type.GetBaseValue()
	levelMultiplier := 1.0 + float64(ft.Level)*0.2
	durabilityMultiplier := ft.Durability / ft.MaxDurability

	return baseValue * levelMultiplier * durabilityMultiplier
}

// GetProductivityBonus 获取生产力奖励
func (ft *FarmTool) GetProductivityBonus() float64 {
	if !ft.IsActive {
		return 1.0
	}

	return ft.Efficiency
}

// CareRecord 照料记录
type CareRecord struct {
	Type      string
	Value     float64
	Quality   float64
	Timestamp time.Time
}

// GrowthBonus 生长奖励
type GrowthBonus struct {
	ID         string
	Type       string
	Multiplier float64
	Duration   time.Duration
	StartTime  time.Time
}

// IsActive 检查奖励是否激活
func (gb *GrowthBonus) IsActive() bool {
	return time.Now().Before(gb.StartTime.Add(gb.Duration))
}

// ToolEffect 工具效果
type ToolEffect struct {
	Type  string
	Value float64
}

// GetType 获取类型
func (te *ToolEffect) GetType() string {
	return te.Type
}

// GetValue 获取数值
func (te *ToolEffect) GetValue() float64 {
	return te.Value
}
