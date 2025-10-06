package synthesis

import (
	"time"

	"github.com/google/uuid"
	// "math/rand"
	// "github.com/google/uuid"
)

// Recipe 配方实体
type Recipe struct {
	id           string
	name         string
	category     RecipeCategory
	requirements []*MaterialRequirement
	outputs      []*ItemOutput
	failOutputs  []*ItemOutput
	successRate  float64
	craftTime    time.Duration
	requireLevel int
	description  string
	unlockedAt   time.Time
}

// NewRecipe 创建配方
func NewRecipe(name string, category RecipeCategory, successRate float64) *Recipe {
	return &Recipe{
		id:           uuid.New().String(),
		name:         name,
		category:     category,
		requirements: make([]*MaterialRequirement, 0),
		outputs:      make([]*ItemOutput, 0),
		failOutputs:  make([]*ItemOutput, 0),
		successRate:  successRate,
		craftTime:    time.Minute * 5, // 默认5分钟
		requireLevel: 1,
		unlockedAt:   time.Now(),
	}
}

// GetID 获取配方ID
func (r *Recipe) GetID() string {
	return r.id
}

// GetName 获取配方名称
func (r *Recipe) GetName() string {
	return r.name
}

// GetCategory 获取配方分类
func (r *Recipe) GetCategory() RecipeCategory {
	return r.category
}

// AddRequirement 添加材料需求
func (r *Recipe) AddRequirement(materialID string, quantity int) {
	r.requirements = append(r.requirements, &MaterialRequirement{
		MaterialID: materialID,
		Quantity:   quantity,
	})
}

// GetRequirements 获取材料需求
func (r *Recipe) GetRequirements() []*MaterialRequirement {
	return r.requirements
}

// AddOutput 添加产出物品
func (r *Recipe) AddOutput(itemID string, quantity int, probability float64) {
	r.outputs = append(r.outputs, &ItemOutput{
		ItemID:      itemID,
		Quantity:    quantity,
		Probability: probability,
	})
}

// GetOutputs 获取产出物品
func (r *Recipe) GetOutputs() []*ItemOutput {
	return r.outputs
}

// AddFailOutput 添加失败产出
func (r *Recipe) AddFailOutput(itemID string, quantity int, probability float64) {
	r.failOutputs = append(r.failOutputs, &ItemOutput{
		ItemID:      itemID,
		Quantity:    quantity,
		Probability: probability,
	})
}

// GetFailOutputs 获取失败产出
func (r *Recipe) GetFailOutputs() []*ItemOutput {
	return r.failOutputs
}

// GetSuccessRate 获取成功率
func (r *Recipe) GetSuccessRate() float64 {
	return r.successRate
}

// SetSuccessRate 设置成功率
func (r *Recipe) SetSuccessRate(rate float64) {
	if rate < 0 {
		rate = 0
	} else if rate > 1 {
		rate = 1
	}
	r.successRate = rate
}

// GetCraftTime 获取制作时间
func (r *Recipe) GetCraftTime() time.Duration {
	return r.craftTime
}

// SetCraftTime 设置制作时间
func (r *Recipe) SetCraftTime(duration time.Duration) {
	r.craftTime = duration
}

// GetRequireLevel 获取需求等级
func (r *Recipe) GetRequireLevel() int {
	return r.requireLevel
}

// SetRequireLevel 设置需求等级
func (r *Recipe) SetRequireLevel(level int) {
	r.requireLevel = level
}

// GetDescription 获取描述
func (r *Recipe) GetDescription() string {
	return r.description
}

// SetDescription 设置描述
func (r *Recipe) SetDescription(desc string) {
	r.description = desc
}

// GetUnlockedAt 获取解锁时间
func (r *Recipe) GetUnlockedAt() time.Time {
	return r.unlockedAt
}

// Material 材料实体
type Material struct {
	id           string
	name         string
	materialType MaterialType
	quality      Quality
	quantity     int
	maxStack     int
	description  string
	obtainedAt   time.Time
}

// NewMaterial 创建材料
func NewMaterial(id, name string, materialType MaterialType, quality Quality, quantity int) *Material {
	return &Material{
		id:           id,
		name:         name,
		materialType: materialType,
		quality:      quality,
		quantity:     quantity,
		maxStack:     999, // 默认最大堆叠999
		obtainedAt:   time.Now(),
	}
}

// GetID 获取材料ID
func (m *Material) GetID() string {
	return m.id
}

// GetName 获取材料名称
func (m *Material) GetName() string {
	return m.name
}

// GetType 获取材料类型
func (m *Material) GetType() MaterialType {
	return m.materialType
}

// GetQuality 获取品质
func (m *Material) GetQuality() Quality {
	return m.quality
}

// GetQuantity 获取数量
func (m *Material) GetQuantity() int {
	return m.quantity
}

// AddQuantity 增加数量
func (m *Material) AddQuantity(amount int) {
	m.quantity += amount
	if m.quantity > m.maxStack {
		m.quantity = m.maxStack
	}
}

// ConsumeQuantity 消耗数量
func (m *Material) ConsumeQuantity(amount int) error {
	if m.quantity < amount {
		return ErrInsufficientMaterial
	}
	m.quantity -= amount
	return nil
}

// GetMaxStack 获取最大堆叠
func (m *Material) GetMaxStack() int {
	return m.maxStack
}

// SetMaxStack 设置最大堆叠
func (m *Material) SetMaxStack(maxStack int) {
	m.maxStack = maxStack
}

// GetDescription 获取描述
func (m *Material) GetDescription() string {
	return m.description
}

// SetDescription 设置描述
func (m *Material) SetDescription(desc string) {
	m.description = desc
}

// GetObtainedAt 获取获得时间
func (m *Material) GetObtainedAt() time.Time {
	return m.obtainedAt
}

// SynthesisRecord 合成记录实体
type SynthesisRecord struct {
	id        string
	playerID  string
	recipeID  string
	quantity  int
	result    *SynthesisResult
	createdAt time.Time
}

// NewSynthesisRecord 创建合成记录
func NewSynthesisRecord(playerID, recipeID string, quantity int, result *SynthesisResult) *SynthesisRecord {
	return &SynthesisRecord{
		id:        uuid.New().String(),
		playerID:  playerID,
		recipeID:  recipeID,
		quantity:  quantity,
		result:    result,
		createdAt: time.Now(),
	}
}

// GetID 获取记录ID
func (sr *SynthesisRecord) GetID() string {
	return sr.id
}

// GetPlayerID 获取玩家ID
func (sr *SynthesisRecord) GetPlayerID() string {
	return sr.playerID
}

// GetRecipeID 获取配方ID
func (sr *SynthesisRecord) GetRecipeID() string {
	return sr.recipeID
}

// GetQuantity 获取合成数量
func (sr *SynthesisRecord) GetQuantity() int {
	return sr.quantity
}

// GetResult 获取合成结果
func (sr *SynthesisRecord) GetResult() *SynthesisResult {
	return sr.result
}

// GetCreatedAt 获取创建时间
func (sr *SynthesisRecord) GetCreatedAt() time.Time {
	return sr.createdAt
}

// SynthesisResult 合成结果
type SynthesisResult struct {
	successItems map[string]int // 成功获得的物品
	failItems    map[string]int // 失败获得的物品
	successCount int            // 成功次数
	failCount    int            // 失败次数
}

// NewSynthesisResult 创建合成结果
func NewSynthesisResult() *SynthesisResult {
	return &SynthesisResult{
		successItems: make(map[string]int),
		failItems:    make(map[string]int),
		successCount: 0,
		failCount:    0,
	}
}

// AddSuccessItem 添加成功物品
func (sr *SynthesisResult) AddSuccessItem(itemID string, quantity int) {
	sr.successItems[itemID] += quantity
	sr.successCount++
}

// AddFailItem 添加失败物品
func (sr *SynthesisResult) AddFailItem(itemID string, quantity int) {
	sr.failItems[itemID] += quantity
	sr.failCount++
}

// GetSuccessItems 获取成功物品
func (sr *SynthesisResult) GetSuccessItems() map[string]int {
	return sr.successItems
}

// GetFailItems 获取失败物品
func (sr *SynthesisResult) GetFailItems() map[string]int {
	return sr.failItems
}

// GetSuccessCount 获取成功次数
func (sr *SynthesisResult) GetSuccessCount() int {
	return sr.successCount
}

// GetFailCount 获取失败次数
func (sr *SynthesisResult) GetFailCount() int {
	return sr.failCount
}

// GetTotalCount 获取总次数
func (sr *SynthesisResult) GetTotalCount() int {
	return sr.successCount + sr.failCount
}

// GetSuccessRate 获取成功率
func (sr *SynthesisResult) GetSuccessRate() float64 {
	total := sr.GetTotalCount()
	if total == 0 {
		return 0
	}
	return float64(sr.successCount) / float64(total)
}
