package synthesis

import "context"

// SynthesisRepository 合成仓储接口
type SynthesisRepository interface {
	// SaveSynthesisAggregate 保存合成聚合根
	SaveSynthesisAggregate(ctx context.Context, aggregate *SynthesisAggregate) error
	
	// GetSynthesisAggregate 获取合成聚合根
	GetSynthesisAggregate(ctx context.Context, playerID string) (*SynthesisAggregate, error)
	
	// DeleteSynthesisAggregate 删除合成聚合根
	DeleteSynthesisAggregate(ctx context.Context, playerID string) error
	
	// SaveRecipe 保存配方
	SaveRecipe(ctx context.Context, playerID string, recipe *Recipe) error
	
	// GetRecipe 获取配方
	GetRecipe(ctx context.Context, playerID, recipeID string) (*Recipe, error)
	
	// GetPlayerRecipes 获取玩家所有配方
	GetPlayerRecipes(ctx context.Context, playerID string) ([]*Recipe, error)
	
	// DeleteRecipe 删除配方
	DeleteRecipe(ctx context.Context, playerID, recipeID string) error
	
	// SaveMaterial 保存材料
	SaveMaterial(ctx context.Context, playerID string, material *Material) error
	
	// GetMaterial 获取材料
	GetMaterial(ctx context.Context, playerID, materialID string) (*Material, error)
	
	// GetPlayerMaterials 获取玩家所有材料
	GetPlayerMaterials(ctx context.Context, playerID string) ([]*Material, error)
	
	// UpdateMaterialQuantity 更新材料数量
	UpdateMaterialQuantity(ctx context.Context, playerID, materialID string, quantity int) error
	
	// DeleteMaterial 删除材料
	DeleteMaterial(ctx context.Context, playerID, materialID string) error
	
	// SaveSynthesisRecord 保存合成记录
	SaveSynthesisRecord(ctx context.Context, record *SynthesisRecord) error
	
	// GetSynthesisRecords 获取合成记录
	GetSynthesisRecords(ctx context.Context, playerID string, limit, offset int) ([]*SynthesisRecord, error)
	
	// GetRecipesByCategory 根据分类获取配方
	GetRecipesByCategory(ctx context.Context, playerID string, category RecipeCategory) ([]*Recipe, error)
	
	// GetMaterialsByType 根据类型获取材料
	GetMaterialsByType(ctx context.Context, playerID string, materialType MaterialType) ([]*Material, error)
	
	// GetMaterialsByQuality 根据品质获取材料
	GetMaterialsByQuality(ctx context.Context, playerID string, quality Quality) ([]*Material, error)
	
	// GetRecipeCount 获取配方数量
	GetRecipeCount(ctx context.Context, playerID string) (int, error)
	
	// GetMaterialCount 获取材料数量
	GetMaterialCount(ctx context.Context, playerID string) (int, error)
}

// RecipeTemplateRepository 配方模板仓储接口
type RecipeTemplateRepository interface {
	// GetRecipeTemplate 获取配方模板
	GetRecipeTemplate(ctx context.Context, templateID string) (*RecipeTemplate, error)
	
	// GetRecipeTemplatesByCategory 根据分类获取配方模板
	GetRecipeTemplatesByCategory(ctx context.Context, category RecipeCategory) ([]*RecipeTemplate, error)
	
	// GetRecipeTemplatesByLevel 根据等级获取配方模板
	GetRecipeTemplatesByLevel(ctx context.Context, minLevel, maxLevel int) ([]*RecipeTemplate, error)
	
	// SaveRecipeTemplate 保存配方模板
	SaveRecipeTemplate(ctx context.Context, template *RecipeTemplate) error
	
	// DeleteRecipeTemplate 删除配方模板
	DeleteRecipeTemplate(ctx context.Context, templateID string) error
	
	// GetAllRecipeTemplates 获取所有配方模板
	GetAllRecipeTemplates(ctx context.Context) ([]*RecipeTemplate, error)
}

// MaterialTemplateRepository 材料模板仓储接口
type MaterialTemplateRepository interface {
	// GetMaterialTemplate 获取材料模板
	GetMaterialTemplate(ctx context.Context, templateID string) (*MaterialTemplate, error)
	
	// GetMaterialTemplatesByType 根据类型获取材料模板
	GetMaterialTemplatesByType(ctx context.Context, materialType MaterialType) ([]*MaterialTemplate, error)
	
	// GetMaterialTemplatesByQuality 根据品质获取材料模板
	GetMaterialTemplatesByQuality(ctx context.Context, quality Quality) ([]*MaterialTemplate, error)
	
	// SaveMaterialTemplate 保存材料模板
	SaveMaterialTemplate(ctx context.Context, template *MaterialTemplate) error
	
	// DeleteMaterialTemplate 删除材料模板
	DeleteMaterialTemplate(ctx context.Context, templateID string) error
	
	// GetAllMaterialTemplates 获取所有材料模板
	GetAllMaterialTemplates(ctx context.Context) ([]*MaterialTemplate, error)
}

// RecipeTemplate 配方模板
type RecipeTemplate struct {
	ID           string                  `json:"id"`
	Name         string                  `json:"name"`
	Category     RecipeCategory          `json:"category"`
	Requirements []*MaterialRequirement  `json:"requirements"`
	Outputs      []*ItemOutput           `json:"outputs"`
	FailOutputs  []*ItemOutput           `json:"fail_outputs"`
	SuccessRate  float64                 `json:"success_rate"`
	CraftTime    int64                   `json:"craft_time"` // 毫秒
	RequireLevel int                     `json:"require_level"`
	Conditions   []*CraftingCondition    `json:"conditions"`
	Description  string                  `json:"description"`
	IconURL      string                  `json:"icon_url"`
}

// CreateRecipeFromTemplate 从模板创建配方
func (rt *RecipeTemplate) CreateRecipeFromTemplate() *Recipe {
	recipe := NewRecipe(rt.Name, rt.Category, rt.SuccessRate)
	
	// 复制材料需求
	for _, req := range rt.Requirements {
		recipe.AddRequirement(req.MaterialID, req.Quantity)
	}
	
	// 复制产出
	for _, output := range rt.Outputs {
		recipe.AddOutput(output.ItemID, output.Quantity, output.Probability)
	}
	
	// 复制失败产出
	for _, failOutput := range rt.FailOutputs {
		recipe.AddFailOutput(failOutput.ItemID, failOutput.Quantity, failOutput.Probability)
	}
	
	recipe.SetRequireLevel(rt.RequireLevel)
	recipe.SetDescription(rt.Description)
	
	return recipe
}

// MaterialTemplate 材料模板
type MaterialTemplate struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Type         MaterialType `json:"type"`
	Quality      Quality      `json:"quality"`
	MaxStack     int          `json:"max_stack"`
	Description  string       `json:"description"`
	IconURL      string       `json:"icon_url"`
	ObtainMethods []string    `json:"obtain_methods"` // 获取方式
}

// CreateMaterialFromTemplate 从模板创建材料
func (mt *MaterialTemplate) CreateMaterialFromTemplate(quantity int) *Material {
	material := NewMaterial(mt.ID, mt.Name, mt.Type, mt.Quality, quantity)
	material.SetMaxStack(mt.MaxStack)
	material.SetDescription(mt.Description)
	return material
}