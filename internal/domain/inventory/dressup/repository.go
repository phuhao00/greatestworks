package dressup

import "context"

// DressupRepository 换装仓储接口
type DressupRepository interface {
	// SaveDressupAggregate 保存换装聚合根
	SaveDressupAggregate(ctx context.Context, aggregate *DressupAggregate) error
	
	// GetDressupAggregate 获取换装聚合根
	GetDressupAggregate(ctx context.Context, playerID string) (*DressupAggregate, error)
	
	// DeleteDressupAggregate 删除换装聚合根
	DeleteDressupAggregate(ctx context.Context, playerID string) error
	
	// SaveOutfit 保存服装
	SaveOutfit(ctx context.Context, playerID string, outfit *Outfit) error
	
	// GetOutfit 获取服装
	GetOutfit(ctx context.Context, playerID, outfitID string) (*Outfit, error)
	
	// GetPlayerOutfits 获取玩家所有服装
	GetPlayerOutfits(ctx context.Context, playerID string) ([]*Outfit, error)
	
	// DeleteOutfit 删除服装
	DeleteOutfit(ctx context.Context, playerID, outfitID string) error
	
	// SaveOutfitSet 保存套装配置
	SaveOutfitSet(ctx context.Context, playerID string, outfitSet *OutfitSet) error
	
	// GetOutfitSet 获取套装配置
	GetOutfitSet(ctx context.Context, playerID string) (*OutfitSet, error)
	
	// GetOutfitsByType 根据类型获取服装
	GetOutfitsByType(ctx context.Context, playerID string, outfitType OutfitType) ([]*Outfit, error)
	
	// GetOutfitsByRarity 根据稀有度获取服装
	GetOutfitsByRarity(ctx context.Context, playerID string, rarity Rarity) ([]*Outfit, error)
	
	// UpdateOutfitLockStatus 更新服装锁定状态
	UpdateOutfitLockStatus(ctx context.Context, playerID, outfitID string, locked bool) error
	
	// GetOutfitCount 获取服装数量
	GetOutfitCount(ctx context.Context, playerID string) (int, error)
	
	// GetOutfitsBySlot 根据槽位获取可装备的服装
	GetOutfitsBySlot(ctx context.Context, playerID string, slot OutfitSlot) ([]*Outfit, error)
}

// OutfitTemplateRepository 服装模板仓储接口
type OutfitTemplateRepository interface {
	// GetOutfitTemplate 获取服装模板
	GetOutfitTemplate(ctx context.Context, templateID string) (*OutfitTemplate, error)
	
	// GetOutfitTemplatesByType 根据类型获取服装模板
	GetOutfitTemplatesByType(ctx context.Context, outfitType OutfitType) ([]*OutfitTemplate, error)
	
	// GetOutfitTemplatesByRarity 根据稀有度获取服装模板
	GetOutfitTemplatesByRarity(ctx context.Context, rarity Rarity) ([]*OutfitTemplate, error)
	
	// SaveOutfitTemplate 保存服装模板
	SaveOutfitTemplate(ctx context.Context, template *OutfitTemplate) error
	
	// DeleteOutfitTemplate 删除服装模板
	DeleteOutfitTemplate(ctx context.Context, templateID string) error
}

// OutfitTemplate 服装模板
type OutfitTemplate struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        OutfitType        `json:"type"`
	Rarity      Rarity            `json:"rarity"`
	BaseAttrs   map[string]int    `json:"base_attrs"`
	Slots       []OutfitSlot      `json:"slots"`
	RequireLevel int              `json:"require_level"`
	Description string            `json:"description"`
	IconURL     string            `json:"icon_url"`
	ModelURL    string            `json:"model_url"`
}

// CreateOutfitFromTemplate 从模板创建服装
func (ot *OutfitTemplate) CreateOutfitFromTemplate() *Outfit {
	outfit := NewOutfit(ot.Name, ot.Type, ot.Rarity)
	
	// 复制基础属性
	for attr, value := range ot.BaseAttrs {
		outfit.AddAttribute(attr, value)
	}
	
	// 添加槽位
	for _, slot := range ot.Slots {
		outfit.AddSlot(slot)
	}
	
	return outfit
}