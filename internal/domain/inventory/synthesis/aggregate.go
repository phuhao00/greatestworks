package synthesis

import (
	"time"
	"github.com/google/uuid"
)

// SynthesisAggregate 合成聚合根
type SynthesisAggregate struct {
	playerID    string
	recipes     map[string]*Recipe
	materials   map[string]*Material
	history     []*SynthesisRecord
	updatedAt   time.Time
	version     int
}

// NewSynthesisAggregate 创建合成聚合根
func NewSynthesisAggregate(playerID string) *SynthesisAggregate {
	return &SynthesisAggregate{
		playerID:  playerID,
		recipes:   make(map[string]*Recipe),
		materials: make(map[string]*Material),
		history:   make([]*SynthesisRecord, 0),
		updatedAt: time.Now(),
		version:   1,
	}
}

// GetPlayerID 获取玩家ID
func (s *SynthesisAggregate) GetPlayerID() string {
	return s.playerID
}

// AddRecipe 添加配方
func (s *SynthesisAggregate) AddRecipe(recipe *Recipe) error {
	if recipe == nil {
		return ErrInvalidRecipe
	}
	
	s.recipes[recipe.GetID()] = recipe
	s.updateVersion()
	return nil
}

// GetRecipe 获取配方
func (s *SynthesisAggregate) GetRecipe(recipeID string) *Recipe {
	return s.recipes[recipeID]
}

// GetAllRecipes 获取所有配方
func (s *SynthesisAggregate) GetAllRecipes() map[string]*Recipe {
	return s.recipes
}

// AddMaterial 添加材料
func (s *SynthesisAggregate) AddMaterial(material *Material) error {
	if material == nil {
		return ErrInvalidMaterial
	}
	
	existing, exists := s.materials[material.GetID()]
	if exists {
		existing.AddQuantity(material.GetQuantity())
	} else {
		s.materials[material.GetID()] = material
	}
	
	s.updateVersion()
	return nil
}

// ConsumeMaterial 消耗材料
func (s *SynthesisAggregate) ConsumeMaterial(materialID string, quantity int) error {
	material, exists := s.materials[materialID]
	if !exists {
		return ErrMaterialNotFound
	}
	
	if material.GetQuantity() < quantity {
		return ErrInsufficientMaterial
	}
	
	material.ConsumeQuantity(quantity)
	if material.GetQuantity() <= 0 {
		delete(s.materials, materialID)
	}
	
	s.updateVersion()
	return nil
}

// GetMaterial 获取材料
func (s *SynthesisAggregate) GetMaterial(materialID string) *Material {
	return s.materials[materialID]
}

// GetAllMaterials 获取所有材料
func (s *SynthesisAggregate) GetAllMaterials() map[string]*Material {
	return s.materials
}

// CanSynthesize 检查是否可以合成
func (s *SynthesisAggregate) CanSynthesize(recipeID string) error {
	recipe, exists := s.recipes[recipeID]
	if !exists {
		return ErrRecipeNotFound
	}
	
	// 检查材料是否足够
	for _, requirement := range recipe.GetRequirements() {
		material, exists := s.materials[requirement.MaterialID]
		if !exists || material.GetQuantity() < requirement.Quantity {
			return ErrInsufficientMaterial
		}
	}
	
	return nil
}

// Synthesize 执行合成
func (s *SynthesisAggregate) Synthesize(recipeID string, quantity int) (*SynthesisResult, error) {
	recipe, exists := s.recipes[recipeID]
	if !exists {
		return nil, ErrRecipeNotFound
	}
	
	// 检查材料是否足够
	for _, requirement := range recipe.GetRequirements() {
		requiredQuantity := requirement.Quantity * quantity
		material, exists := s.materials[requirement.MaterialID]
		if !exists || material.GetQuantity() < requiredQuantity {
			return nil, ErrInsufficientMaterial
		}
	}
	
	// 消耗材料
	for _, requirement := range recipe.GetRequirements() {
		requiredQuantity := requirement.Quantity * quantity
		err := s.ConsumeMaterial(requirement.MaterialID, requiredQuantity)
		if err != nil {
			return nil, err
		}
	}
	
	// 计算合成结果
	result := s.calculateSynthesisResult(recipe, quantity)
	
	// 记录合成历史
	record := NewSynthesisRecord(s.playerID, recipeID, quantity, result)
	s.history = append(s.history, record)
	
	s.updateVersion()
	return result, nil
}

// calculateSynthesisResult 计算合成结果
func (s *SynthesisAggregate) calculateSynthesisResult(recipe *Recipe, quantity int) *SynthesisResult {
	result := NewSynthesisResult()
	
	for i := 0; i < quantity; i++ {
		// 计算成功率
		if s.rollSuccess(recipe.GetSuccessRate()) {
			// 成功，添加产出物品
			for _, output := range recipe.GetOutputs() {
				result.AddSuccessItem(output.ItemID, output.Quantity)
			}
		} else {
			// 失败，可能有失败产出
			for _, failOutput := range recipe.GetFailOutputs() {
				result.AddFailItem(failOutput.ItemID, failOutput.Quantity)
			}
		}
	}
	
	return result
}

// rollSuccess 计算成功率
func (s *SynthesisAggregate) rollSuccess(successRate float64) bool {
	// 这里可以加入更复杂的成功率计算逻辑
	// 比如玩家技能等级、装备加成等
	return rand.Float64() < successRate
}

// GetSynthesisHistory 获取合成历史
func (s *SynthesisAggregate) GetSynthesisHistory() []*SynthesisRecord {
	return s.history
}

// GetRecentHistory 获取最近的合成历史
func (s *SynthesisAggregate) GetRecentHistory(limit int) []*SynthesisRecord {
	if len(s.history) <= limit {
		return s.history
	}
	return s.history[len(s.history)-limit:]
}

// updateVersion 更新版本
func (s *SynthesisAggregate) updateVersion() {
	s.version++
	s.updatedAt = time.Now()
}

// GetVersion 获取版本
func (s *SynthesisAggregate) GetVersion() int {
	return s.version
}

// GetUpdatedAt 获取更新时间
func (s *SynthesisAggregate) GetUpdatedAt() time.Time {
	return s.updatedAt
}

// GetMaterialQuantity 获取材料数量
func (s *SynthesisAggregate) GetMaterialQuantity(materialID string) int {
	material, exists := s.materials[materialID]
	if !exists {
		return 0
	}
	return material.GetQuantity()
}

// HasRecipe 检查是否拥有配方
func (s *SynthesisAggregate) HasRecipe(recipeID string) bool {
	_, exists := s.recipes[recipeID]
	return exists
}

// GetRecipesByCategory 根据分类获取配方
func (s *SynthesisAggregate) GetRecipesByCategory(category RecipeCategory) []*Recipe {
	var recipes []*Recipe
	for _, recipe := range s.recipes {
		if recipe.GetCategory() == category {
			recipes = append(recipes, recipe)
		}
	}
	return recipes
}