package dressup

import (
	"math/rand"
	"time"
)

// DressupService 换装领域服务
type DressupService struct {
	outfitFactory *OutfitFactory
	styleManager  *StyleManager
}

// NewDressupService 创建换装服务
func NewDressupService() *DressupService {
	return &DressupService{
		outfitFactory: NewOutfitFactory(),
		styleManager:  NewStyleManager(),
	}
}

// CreateRandomOutfit 创建随机服装
func (ds *DressupService) CreateRandomOutfit(outfitType OutfitType, playerLevel int) *Outfit {
	return ds.outfitFactory.CreateRandomOutfit(outfitType, playerLevel)
}

// CalculateTotalAttributes 计算总属性
func (ds *DressupService) CalculateTotalAttributes(aggregate *DressupAggregate) map[string]int {
	totalAttrs := make(map[string]int)
	
	// 计算装备属性
	for _, outfit := range aggregate.GetCurrentSet().GetAllEquipped() {
		if outfit != nil {
			for attr, value := range outfit.GetAttributes() {
				totalAttrs[attr] += value
			}
		}
	}
	
	// 添加套装加成
	for attr, bonus := range aggregate.GetCurrentSet().GetSetBonuses() {
		totalAttrs[attr] += bonus
	}
	
	return totalAttrs
}

// ValidateOutfitCombination 验证服装搭配
func (ds *DressupService) ValidateOutfitCombination(outfits map[OutfitSlot]*Outfit) error {
	// 检查是否有冲突的装备
	for slot, outfit := range outfits {
		if outfit != nil && !outfit.CanEquipToSlot(slot) {
			return ErrInvalidSlot
		}
	}
	return nil
}

// ApplyStyle 应用换装风格
func (ds *DressupService) ApplyStyle(aggregate *DressupAggregate, styleID string) error {
	style := ds.styleManager.GetStyle(styleID)
	if style == nil {
		return ErrStyleNotFound
	}
	
	// 应用风格逻辑
	// 这里可以根据风格调整装备外观等
	return nil
}

// OutfitFactory 服装工厂
type OutfitFactory struct {
	random *rand.Rand
}

// NewOutfitFactory 创建服装工厂
func NewOutfitFactory() *OutfitFactory {
	return &OutfitFactory{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CreateRandomOutfit 创建随机服装
func (of *OutfitFactory) CreateRandomOutfit(outfitType OutfitType, playerLevel int) *Outfit {
	// 根据玩家等级确定稀有度
	rarity := of.determineRarity(playerLevel)
	
	// 创建基础服装
	outfit := NewOutfit(of.generateOutfitName(outfitType, rarity), outfitType, rarity)
	
	// 添加对应槽位
	of.addSlotsForType(outfit, outfitType)
	
	// 生成随机属性
	of.generateRandomAttributes(outfit, playerLevel)
	
	return outfit
}

// determineRarity 确定稀有度
func (of *OutfitFactory) determineRarity(playerLevel int) Rarity {
	roll := of.random.Float64()
	
	// 根据玩家等级调整稀有度概率
	levelBonus := float64(playerLevel) / 100.0
	
	switch {
	case roll < 0.5-levelBonus*0.1:
		return RarityCommon
	case roll < 0.75-levelBonus*0.05:
		return RarityUncommon
	case roll < 0.9:
		return RarityRare
	case roll < 0.97:
		return RarityEpic
	case roll < 0.995:
		return RarityLegendary
	default:
		return RarityMythic
	}
}

// generateOutfitName 生成服装名称
func (of *OutfitFactory) generateOutfitName(outfitType OutfitType, rarity Rarity) string {
	prefixes := map[Rarity][]string{
		RarityCommon:    {"普通的", "基础的", "简单的"},
		RarityUncommon:  {"精良的", "优质的", "改良的"},
		RarityRare:      {"稀有的", "精制的", "卓越的"},
		RarityEpic:      {"史诗的", "传说的", "神话的"},
		RarityLegendary: {"传奇的", "不朽的", "永恒的"},
		RarityMythic:    {"神器", "至尊", "无上"},
	}
	
	names := map[OutfitType][]string{
		OutfitTypeWeapon:    {"剑", "刀", "枪", "弓", "法杖"},
		OutfitTypeArmor:     {"护甲", "战袍", "铠甲", "法袍"},
		OutfitTypeHelmet:    {"头盔", "帽子", "头饰", "王冠"},
		OutfitTypeShoes:     {"靴子", "鞋子", "战靴", "法靴"},
		OutfitTypeAccessory: {"戒指", "项链", "手镯", "护符"},
	}
	
	prefixList := prefixes[rarity]
	nameList := names[outfitType]
	
	if len(prefixList) == 0 || len(nameList) == 0 {
		return "未知装备"
	}
	
	prefix := prefixList[of.random.Intn(len(prefixList))]
	name := nameList[of.random.Intn(len(nameList))]
	
	return prefix + name
}

// addSlotsForType 为服装类型添加槽位
func (of *OutfitFactory) addSlotsForType(outfit *Outfit, outfitType OutfitType) {
	switch outfitType {
	case OutfitTypeWeapon:
		outfit.AddSlot(SlotWeapon)
		outfit.AddSlot(SlotFashionWeapon)
	case OutfitTypeArmor:
		outfit.AddSlot(SlotArmor)
		outfit.AddSlot(SlotFashionArmor)
	case OutfitTypeHelmet:
		outfit.AddSlot(SlotHelmet)
		outfit.AddSlot(SlotFashionHelmet)
	case OutfitTypeShoes:
		outfit.AddSlot(SlotShoes)
	case OutfitTypeAccessory:
		outfit.AddSlot(SlotRing)
		outfit.AddSlot(SlotNecklace)
	case OutfitTypePet:
		outfit.AddSlot(SlotPet)
	case OutfitTypeMount:
		outfit.AddSlot(SlotMount)
	}
}

// generateRandomAttributes 生成随机属性
func (of *OutfitFactory) generateRandomAttributes(outfit *Outfit, playerLevel int) {
	baseValue := playerLevel * 2
	multiplier := outfit.GetRarity().GetRarityMultiplier()
	
	// 基础属性
	attack := int(float64(baseValue) * multiplier * (0.8 + of.random.Float64()*0.4))
	defense := int(float64(baseValue) * multiplier * (0.8 + of.random.Float64()*0.4))
	hp := int(float64(baseValue*5) * multiplier * (0.8 + of.random.Float64()*0.4))
	
	outfit.AddAttribute("attack", attack)
	outfit.AddAttribute("defense", defense)
	outfit.AddAttribute("hp", hp)
	
	// 根据稀有度添加额外属性
	if outfit.GetRarity() >= RarityRare {
		critRate := of.random.Intn(10) + 1
		outfit.AddAttribute("crit_rate", critRate)
	}
	
	if outfit.GetRarity() >= RarityEpic {
		critDamage := of.random.Intn(20) + 10
		outfit.AddAttribute("crit_damage", critDamage)
	}
}

// StyleManager 风格管理器
type StyleManager struct {
	styles map[string]*DressupStyle
}

// NewStyleManager 创建风格管理器
func NewStyleManager() *StyleManager {
	sm := &StyleManager{
		styles: make(map[string]*DressupStyle),
	}
	
	// 初始化默认风格
	sm.initDefaultStyles()
	return sm
}

// GetStyle 获取风格
func (sm *StyleManager) GetStyle(styleID string) *DressupStyle {
	return sm.styles[styleID]
}

// AddStyle 添加风格
func (sm *StyleManager) AddStyle(style *DressupStyle) {
	sm.styles[style.GetStyleID()] = style
}

// initDefaultStyles 初始化默认风格
func (sm *StyleManager) initDefaultStyles() {
	// 战士风格
	warriorStyle := NewDressupStyle("warrior", "战士风格", "战斗")
	warriorStyle.AddBonus("attack", 50)
	warriorStyle.AddBonus("defense", 30)
	sm.AddStyle(warriorStyle)
	
	// 法师风格
	mageStyle := NewDressupStyle("mage", "法师风格", "魔法")
	mageStyle.AddBonus("magic_attack", 60)
	mageStyle.AddBonus("mana", 100)
	sm.AddStyle(mageStyle)
	
	// 刺客风格
	assassinStyle := NewDressupStyle("assassin", "刺客风格", "敏捷")
	assassinStyle.AddBonus("crit_rate", 15)
	assassinStyle.AddBonus("speed", 25)
	sm.AddStyle(assassinStyle)
}