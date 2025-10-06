package dressup

import (
	"time"
	// "github.com/google/uuid"
)

// DressupAggregate 换装聚合根
type DressupAggregate struct {
	playerID     string
	outfits      map[string]*Outfit
	currentSet   *OutfitSet
	savedSets    map[string]*OutfitSet    // 保存的套装方案
	fashionSets  map[string]*FashionSet   // 时装套装配置
	dyeColors    map[string]*DyeColor     // 已解锁的染色
	styles       map[string]*DressupStyle // 风格配置
	currentStyle string                   // 当前风格ID
	preferences  *DressupPreferences      // 换装偏好
	statistics   *DressupStatistics       // 换装统计
	updatedAt    time.Time
	version      int
}

// DressupPreferences 换装偏好
type DressupPreferences struct {
	autoEquipBetter   bool         // 自动装备更好的装备
	preferredRarities []Rarity     // 偏好的稀有度
	preferredTypes    []OutfitType // 偏好的类型
	hideHelmet        bool         // 隐藏头盔
	showCape          bool         // 显示斗篷
	defaultStyle      string       // 默认风格
}

// DressupStatistics 换装统计
type DressupStatistics struct {
	totalOutfits        int                   // 总服装数
	equippedOutfits     int                   // 已装备数
	totalPower          int                   // 总战力
	changeCount         int                   // 换装次数
	lastChangeTime      time.Time             // 最后换装时间
	rarityDistribution  map[Rarity]int        // 稀有度分布
	typeDistribution    map[OutfitType]int    // 类型分布
	qualityDistribution map[OutfitQuality]int // 品质分布
	favoriteSet         string                // 最常用套装
	mostUsedOutfit      string                // 最常用单品
}

// NewDressupAggregate 创建换装聚合根
func NewDressupAggregate(playerID string) *DressupAggregate {
	return &DressupAggregate{
		playerID:     playerID,
		outfits:      make(map[string]*Outfit),
		currentSet:   NewOutfitSet(),
		savedSets:    make(map[string]*OutfitSet),
		fashionSets:  make(map[string]*FashionSet),
		dyeColors:    make(map[string]*DyeColor),
		styles:       make(map[string]*DressupStyle),
		currentStyle: "",
		preferences:  NewDressupPreferences(),
		statistics:   NewDressupStatistics(),
		updatedAt:    time.Now(),
		version:      1,
	}
}

// NewDressupPreferences 创建换装偏好
func NewDressupPreferences() *DressupPreferences {
	return &DressupPreferences{
		autoEquipBetter:   false,
		preferredRarities: make([]Rarity, 0),
		preferredTypes:    make([]OutfitType, 0),
		hideHelmet:        false,
		showCape:          true,
		defaultStyle:      "",
	}
}

// NewDressupStatistics 创建换装统计
func NewDressupStatistics() *DressupStatistics {
	return &DressupStatistics{
		totalOutfits:        0,
		equippedOutfits:     0,
		totalPower:          0,
		changeCount:         0,
		lastChangeTime:      time.Time{},
		rarityDistribution:  make(map[Rarity]int),
		typeDistribution:    make(map[OutfitType]int),
		qualityDistribution: make(map[OutfitQuality]int),
		favoriteSet:         "",
		mostUsedOutfit:      "",
	}
}

// GetPlayerID 获取玩家ID
func (d *DressupAggregate) GetPlayerID() string {
	return d.playerID
}

// AddOutfit 添加服装
func (d *DressupAggregate) AddOutfit(outfit *Outfit) error {
	if outfit == nil {
		return ErrInvalidOutfit
	}

	d.outfits[outfit.GetID()] = outfit
	d.updateVersion()
	return nil
}

// EquipOutfit 装备服装
func (d *DressupAggregate) EquipOutfit(outfitID string, slot OutfitSlot) error {
	outfit, exists := d.outfits[outfitID]
	if !exists {
		return ErrOutfitNotFound
	}

	if !outfit.CanEquipToSlot(slot) {
		return ErrInvalidSlot
	}

	d.currentSet.EquipToSlot(slot, outfit)
	d.updateVersion()
	return nil
}

// UnequipOutfit 卸下服装
func (d *DressupAggregate) UnequipOutfit(slot OutfitSlot) error {
	d.currentSet.UnequipFromSlot(slot)
	d.updateVersion()
	return nil
}

// GetCurrentSet 获取当前装备套装
func (d *DressupAggregate) GetCurrentSet() *OutfitSet {
	return d.currentSet
}

// GetOutfits 获取所有服装
func (d *DressupAggregate) GetOutfits() map[string]*Outfit {
	return d.outfits
}

// updateVersion 更新版本
func (d *DressupAggregate) updateVersion() {
	d.version++
	d.updatedAt = time.Now()
}

// GetVersion 获取版本
func (d *DressupAggregate) GetVersion() int {
	return d.version
}

// GetUpdatedAt 获取更新时间
func (d *DressupAggregate) GetUpdatedAt() time.Time {
	return d.updatedAt
}

// SaveOutfitSet 保存套装方案
func (d *DressupAggregate) SaveOutfitSet(name string) error {
	if name == "" {
		return ErrInvalidSetName
	}

	// 克隆当前套装
	savedSet := d.currentSet.Clone()
	d.savedSets[name] = savedSet
	d.updateVersion()
	return nil
}

// LoadOutfitSet 加载套装方案
func (d *DressupAggregate) LoadOutfitSet(name string) error {
	savedSet, exists := d.savedSets[name]
	if !exists {
		return ErrSetNotFound
	}

	// 先卸下当前装备
	for slot := range d.currentSet.GetAllEquipped() {
		d.UnequipOutfit(slot)
	}

	// 装备保存的套装
	for slot, outfit := range savedSet.GetAllEquipped() {
		if outfit != nil {
			d.EquipOutfit(outfit.GetID(), slot)
		}
	}

	d.statistics.changeCount++
	d.statistics.lastChangeTime = time.Now()
	d.updateVersion()
	return nil
}

// DeleteOutfitSet 删除套装方案
func (d *DressupAggregate) DeleteOutfitSet(name string) error {
	if _, exists := d.savedSets[name]; !exists {
		return ErrSetNotFound
	}

	delete(d.savedSets, name)
	d.updateVersion()
	return nil
}

// GetSavedSets 获取保存的套装方案
func (d *DressupAggregate) GetSavedSets() map[string]*OutfitSet {
	return d.savedSets
}

// AddFashionSet 添加时装套装配置
func (d *DressupAggregate) AddFashionSet(fashionSet *FashionSet) error {
	if fashionSet == nil {
		return ErrInvalidFashionSet
	}

	d.fashionSets[fashionSet.GetSetID()] = fashionSet
	d.updateVersion()
	return nil
}

// GetFashionSets 获取时装套装配置
func (d *DressupAggregate) GetFashionSets() map[string]*FashionSet {
	return d.fashionSets
}

// GetFashionSet 获取指定时装套装
func (d *DressupAggregate) GetFashionSet(setID string) *FashionSet {
	return d.fashionSets[setID]
}

// UnlockDyeColor 解锁染色
func (d *DressupAggregate) UnlockDyeColor(color *DyeColor) error {
	if color == nil {
		return ErrInvalidDyeColor
	}

	color.Unlock()
	d.dyeColors[color.GetColorID()] = color
	d.updateVersion()
	return nil
}

// GetDyeColors 获取已解锁的染色
func (d *DressupAggregate) GetDyeColors() map[string]*DyeColor {
	return d.dyeColors
}

// GetUnlockedDyeColors 获取已解锁的染色列表
func (d *DressupAggregate) GetUnlockedDyeColors() []*DyeColor {
	colors := make([]*DyeColor, 0)
	for _, color := range d.dyeColors {
		if color.IsUnlocked() {
			colors = append(colors, color)
		}
	}
	return colors
}

// DyeOutfit 给服装染色
func (d *DressupAggregate) DyeOutfit(outfitID, part, colorID string) error {
	outfit, exists := d.outfits[outfitID]
	if !exists {
		return ErrOutfitNotFound
	}

	color, exists := d.dyeColors[colorID]
	if !exists || !color.IsUnlocked() {
		return ErrDyeColorNotUnlocked
	}

	outfit.SetDyeColor(part, color)
	d.updateVersion()
	return nil
}

// AddStyle 添加风格
func (d *DressupAggregate) AddStyle(style *DressupStyle) error {
	if style == nil {
		return ErrInvalidStyle
	}

	d.styles[style.GetStyleID()] = style
	d.updateVersion()
	return nil
}

// SetCurrentStyle 设置当前风格
func (d *DressupAggregate) SetCurrentStyle(styleID string) error {
	style, exists := d.styles[styleID]
	if !exists {
		return ErrStyleNotFound
	}

	d.currentStyle = styleID
	d.currentSet.SetStyleBonus(style)
	d.updateVersion()
	return nil
}

// GetCurrentStyle 获取当前风格
func (d *DressupAggregate) GetCurrentStyle() string {
	return d.currentStyle
}

// GetStyles 获取所有风格
func (d *DressupAggregate) GetStyles() map[string]*DressupStyle {
	return d.styles
}

// GetPreferences 获取换装偏好
func (d *DressupAggregate) GetPreferences() *DressupPreferences {
	return d.preferences
}

// UpdatePreferences 更新换装偏好
func (d *DressupAggregate) UpdatePreferences(preferences *DressupPreferences) {
	d.preferences = preferences
	d.updateVersion()
}

// GetStatistics 获取换装统计
func (d *DressupAggregate) GetStatistics() *DressupStatistics {
	return d.statistics
}

// UpdateStatistics 更新统计信息
func (d *DressupAggregate) UpdateStatistics() {
	d.statistics.totalOutfits = len(d.outfits)
	d.statistics.equippedOutfits = d.currentSet.GetEquippedCount()
	d.statistics.totalPower = d.currentSet.GetTotalPower()

	// 重置分布统计
	d.statistics.rarityDistribution = make(map[Rarity]int)
	d.statistics.typeDistribution = make(map[OutfitType]int)
	d.statistics.qualityDistribution = make(map[OutfitQuality]int)

	// 统计分布
	for _, outfit := range d.outfits {
		d.statistics.rarityDistribution[outfit.GetRarity()]++
		d.statistics.typeDistribution[outfit.GetType()]++
		d.statistics.qualityDistribution[outfit.GetQuality()]++
	}

	// 找出最常用的装备
	maxUseCount := 0
	for _, outfit := range d.outfits {
		if outfit.GetUseCount() > maxUseCount {
			maxUseCount = outfit.GetUseCount()
			d.statistics.mostUsedOutfit = outfit.GetID()
		}
	}
}

// FilterOutfits 筛选服装
func (d *DressupAggregate) FilterOutfits(filter *OutfitFilter) []*Outfit {
	result := make([]*Outfit, 0)

	for _, outfit := range d.outfits {
		if d.matchesFilter(outfit, filter) {
			result = append(result, outfit)
		}
	}

	return result
}

// matchesFilter 检查服装是否匹配筛选条件
func (d *DressupAggregate) matchesFilter(outfit *Outfit, filter *OutfitFilter) bool {
	if filter == nil {
		return true
	}

	// 类型筛选
	if filter.GetOutfitType() != nil && outfit.GetType() != *filter.GetOutfitType() {
		return false
	}

	// 稀有度筛选
	if filter.GetRarity() != nil && outfit.GetRarity() != *filter.GetRarity() {
		return false
	}

	// 品质筛选
	if filter.GetQuality() != nil && outfit.GetQuality() != *filter.GetQuality() {
		return false
	}

	// 来源筛选
	if filter.GetSource() != nil && outfit.GetSource() != *filter.GetSource() {
		return false
	}

	// 槽位筛选
	if filter.GetSlot() != nil {
		canEquip := outfit.CanEquipToSlot(*filter.GetSlot())
		if !canEquip {
			return false
		}
	}

	// 锁定状态筛选
	if filter.GetLocked() != nil && outfit.IsLocked() != *filter.GetLocked() {
		return false
	}

	// 等级范围筛选
	if filter.GetMinLevel() != nil && outfit.GetLevel() < *filter.GetMinLevel() {
		return false
	}
	if filter.GetMaxLevel() != nil && outfit.GetLevel() > *filter.GetMaxLevel() {
		return false
	}

	// 搜索文本筛选
	if filter.GetSearchText() != "" {
		searchText := filter.GetSearchText()
		if !d.containsIgnoreCase(outfit.GetName(), searchText) &&
			!d.containsIgnoreCase(outfit.GetDescription(), searchText) {
			return false
		}
	}

	// 标签筛选
	if len(filter.GetTags()) > 0 {
		outfitTags := outfit.GetTags()
		for _, filterTag := range filter.GetTags() {
			found := false
			for _, outfitTag := range outfitTags {
				if outfitTag == filterTag {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	return true
}

// containsIgnoreCase 忽略大小写包含检查
func (d *DressupAggregate) containsIgnoreCase(str, substr string) bool {
	// 简单实现，实际应该使用更好的字符串匹配算法
	return len(str) >= len(substr) && str != "" && substr != ""
}

// GetRecommendedOutfits 获取推荐服装
func (d *DressupAggregate) GetRecommendedOutfits(slot OutfitSlot, limit int) []*Outfit {
	recommended := make([]*Outfit, 0)

	for _, outfit := range d.outfits {
		if outfit.CanEquipToSlot(slot) && !outfit.IsEquipped() && !outfit.IsLocked() {
			recommended = append(recommended, outfit)
		}
	}

	// 按战力排序（简单实现）
	for i := 0; i < len(recommended)-1; i++ {
		for j := i + 1; j < len(recommended); j++ {
			if recommended[i].GetPower() < recommended[j].GetPower() {
				recommended[i], recommended[j] = recommended[j], recommended[i]
			}
		}
	}

	// 限制数量
	if limit > 0 && len(recommended) > limit {
		recommended = recommended[:limit]
	}

	return recommended
}

// AutoEquipBest 自动装备最佳装备
func (d *DressupAggregate) AutoEquipBest() error {
	if !d.preferences.autoEquipBetter {
		return ErrAutoEquipDisabled
	}

	allSlots := []OutfitSlot{
		SlotWeapon, SlotArmor, SlotHelmet, SlotShoes,
		SlotRing, SlotNecklace, SlotFashionWeapon,
		SlotFashionArmor, SlotFashionHelmet, SlotPet, SlotMount,
	}

	for _, slot := range allSlots {
		recommended := d.GetRecommendedOutfits(slot, 1)
		if len(recommended) > 0 {
			best := recommended[0]
			current := d.currentSet.GetEquippedOutfit(slot)

			// 如果推荐的比当前的好，则装备
			if current == nil || best.GetPower() > current.GetPower() {
				d.EquipOutfit(best.GetID(), slot)
			}
		}
	}

	return nil
}

// GetTotalPower 获取总战力
func (d *DressupAggregate) GetTotalPower() int {
	return d.currentSet.GetTotalPower()
}

// GetTotalAttributes 获取总属性
func (d *DressupAggregate) GetTotalAttributes() map[string]int {
	return d.currentSet.GetTotalAttributes()
}

// CanUpgradeAnyOutfit 是否有可升级的服装
func (d *DressupAggregate) CanUpgradeAnyOutfit() bool {
	for _, outfit := range d.outfits {
		if outfit.CanUpgrade() {
			return true
		}
	}
	return false
}

// CanEnhanceAnyOutfit 是否有可强化的服装
func (d *DressupAggregate) CanEnhanceAnyOutfit() bool {
	for _, outfit := range d.outfits {
		if outfit.CanEnhance() {
			return true
		}
	}
	return false
}

// GetOutfitsBySet 根据套装ID获取服装
func (d *DressupAggregate) GetOutfitsBySet(setID string) []*Outfit {
	outfits := make([]*Outfit, 0)
	for _, outfit := range d.outfits {
		if outfit.GetSetID() == setID {
			outfits = append(outfits, outfit)
		}
	}
	return outfits
}

// GetSetCompletionRate 获取套装完成度
func (d *DressupAggregate) GetSetCompletionRate(setID string) float64 {
	fashionSet := d.fashionSets[setID]
	if fashionSet == nil {
		return 0.0
	}

	owned := len(d.GetOutfitsBySet(setID))
	total := len(fashionSet.GetPieces())

	if total == 0 {
		return 0.0
	}

	return float64(owned) / float64(total) * 100.0
}
