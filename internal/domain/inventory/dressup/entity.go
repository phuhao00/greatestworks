package dressup

import (
	"time"

	"github.com/google/uuid"
)

// Outfit 服装实体
type Outfit struct {
	id             string
	name           string
	description    string
	outfitType     OutfitType
	rarity         Rarity
	quality        OutfitQuality
	source         OutfitSource
	attributes     map[string]int
	slots          []OutfitSlot
	isLocked       bool
	isEquipped     bool
	level          int
	exp            int
	maxExp         int
	tags           []string
	setID          string // 所属套装ID
	appearance     *AppearanceConfig
	dyeColors      map[string]*DyeColor // 部位 -> 染色
	enhanceLevel   int
	enhanceBonuses map[string]int
	obtainedAt     time.Time
	lastUsedAt     *time.Time
	useCount       int
	metadata       map[string]interface{}
}

// NewOutfit 创建服装实体
func NewOutfit(name string, outfitType OutfitType, rarity Rarity) *Outfit {
	return &Outfit{
		id:             uuid.New().String(),
		name:           name,
		description:    "",
		outfitType:     outfitType,
		rarity:         rarity,
		quality:        QualityNormal,
		source:         SourceShop,
		attributes:     make(map[string]int),
		slots:          []OutfitSlot{},
		isLocked:       false,
		isEquipped:     false,
		level:          1,
		exp:            0,
		maxExp:         100,
		tags:           make([]string, 0),
		setID:          "",
		appearance:     NewAppearanceConfig(),
		dyeColors:      make(map[string]*DyeColor),
		enhanceLevel:   0,
		enhanceBonuses: make(map[string]int),
		obtainedAt:     time.Now(),
		lastUsedAt:     nil,
		useCount:       0,
		metadata:       make(map[string]interface{}),
	}
}

// GetID 获取服装ID
func (o *Outfit) GetID() string {
	return o.id
}

// GetName 获取服装名称
func (o *Outfit) GetName() string {
	return o.name
}

// GetType 获取服装类型
func (o *Outfit) GetType() OutfitType {
	return o.outfitType
}

// GetRarity 获取稀有度
func (o *Outfit) GetRarity() Rarity {
	return o.rarity
}

// AddAttribute 添加属性
func (o *Outfit) AddAttribute(attr string, value int) {
	o.attributes[attr] = value
}

// GetAttributes 获取所有属性
func (o *Outfit) GetAttributes() map[string]int {
	return o.attributes
}

// AddSlot 添加可装备槽位
func (o *Outfit) AddSlot(slot OutfitSlot) {
	o.slots = append(o.slots, slot)
}

// CanEquipToSlot 检查是否可以装备到指定槽位
func (o *Outfit) CanEquipToSlot(slot OutfitSlot) bool {
	for _, s := range o.slots {
		if s == slot {
			return true
		}
	}
	return false
}

// Lock 锁定服装
func (o *Outfit) Lock() {
	o.isLocked = true
}

// Unlock 解锁服装
func (o *Outfit) Unlock() {
	o.isLocked = false
}

// IsLocked 检查是否锁定
func (o *Outfit) IsLocked() bool {
	return o.isLocked
}

// GetObtainedAt 获取获得时间
func (o *Outfit) GetObtainedAt() time.Time {
	return o.obtainedAt
}

// GetDescription 获取描述
func (o *Outfit) GetDescription() string {
	return o.description
}

// SetDescription 设置描述
func (o *Outfit) SetDescription(description string) {
	o.description = description
}

// GetQuality 获取品质
func (o *Outfit) GetQuality() OutfitQuality {
	return o.quality
}

// SetQuality 设置品质
func (o *Outfit) SetQuality(quality OutfitQuality) {
	o.quality = quality
}

// GetSource 获取来源
func (o *Outfit) GetSource() OutfitSource {
	return o.source
}

// SetSource 设置来源
func (o *Outfit) SetSource(source OutfitSource) {
	o.source = source
}

// IsEquipped 是否已装备
func (o *Outfit) IsEquipped() bool {
	return o.isEquipped
}

// SetEquipped 设置装备状态
func (o *Outfit) SetEquipped(equipped bool) {
	o.isEquipped = equipped
	if equipped {
		now := time.Now()
		o.lastUsedAt = &now
		o.useCount++
	}
}

// GetLevel 获取等级
func (o *Outfit) GetLevel() int {
	return o.level
}

// GetExp 获取经验值
func (o *Outfit) GetExp() int {
	return o.exp
}

// GetMaxExp 获取最大经验值
func (o *Outfit) GetMaxExp() int {
	return o.maxExp
}

// AddExp 增加经验值
func (o *Outfit) AddExp(exp int) bool {
	o.exp += exp
	leveledUp := false

	// 检查是否升级
	for o.exp >= o.maxExp && o.level < 100 {
		o.exp -= o.maxExp
		o.level++
		o.maxExp = o.calculateMaxExp(o.level)
		leveledUp = true
	}

	return leveledUp
}

// calculateMaxExp 计算最大经验值
func (o *Outfit) calculateMaxExp(level int) int {
	return 100 + (level-1)*50 // 基础100，每级增加50
}

// GetTags 获取标签
func (o *Outfit) GetTags() []string {
	return o.tags
}

// AddTag 添加标签
func (o *Outfit) AddTag(tag string) {
	// 检查是否已存在
	for _, existingTag := range o.tags {
		if existingTag == tag {
			return
		}
	}
	o.tags = append(o.tags, tag)
}

// RemoveTag 移除标签
func (o *Outfit) RemoveTag(tag string) {
	for i, existingTag := range o.tags {
		if existingTag == tag {
			o.tags = append(o.tags[:i], o.tags[i+1:]...)
			return
		}
	}
}

// GetSetID 获取套装ID
func (o *Outfit) GetSetID() string {
	return o.setID
}

// SetSetID 设置套装ID
func (o *Outfit) SetSetID(setID string) {
	o.setID = setID
}

// GetAppearance 获取外观配置
func (o *Outfit) GetAppearance() *AppearanceConfig {
	return o.appearance
}

// SetAppearance 设置外观配置
func (o *Outfit) SetAppearance(appearance *AppearanceConfig) {
	o.appearance = appearance
}

// GetDyeColors 获取染色配置
func (o *Outfit) GetDyeColors() map[string]*DyeColor {
	return o.dyeColors
}

// SetDyeColor 设置部位染色
func (o *Outfit) SetDyeColor(part string, color *DyeColor) {
	o.dyeColors[part] = color
}

// RemoveDyeColor 移除部位染色
func (o *Outfit) RemoveDyeColor(part string) {
	delete(o.dyeColors, part)
}

// GetEnhanceLevel 获取强化等级
func (o *Outfit) GetEnhanceLevel() int {
	return o.enhanceLevel
}

// Enhance 强化服装
func (o *Outfit) Enhance() bool {
	if o.enhanceLevel >= 20 { // 最大强化等级
		return false
	}

	o.enhanceLevel++

	// 计算强化加成
	for attr, baseValue := range o.attributes {
		enhanceBonus := int(float64(baseValue) * 0.1 * float64(o.enhanceLevel))
		o.enhanceBonuses[attr] = enhanceBonus
	}

	return true
}

// GetEnhanceBonuses 获取强化加成
func (o *Outfit) GetEnhanceBonuses() map[string]int {
	return o.enhanceBonuses
}

// GetTotalAttributes 获取总属性（基础+强化）
func (o *Outfit) GetTotalAttributes() map[string]int {
	total := make(map[string]int)

	// 基础属性
	for attr, value := range o.attributes {
		total[attr] = value
	}

	// 强化加成
	for attr, bonus := range o.enhanceBonuses {
		total[attr] += bonus
	}

	// 品质加成
	qualityMultiplier := o.quality.GetQualityMultiplier()
	for attr, value := range total {
		total[attr] = int(float64(value) * qualityMultiplier)
	}

	// 稀有度加成
	rarityMultiplier := o.rarity.GetRarityMultiplier()
	for attr, value := range total {
		total[attr] = int(float64(value) * rarityMultiplier)
	}

	return total
}

// GetLastUsedAt 获取最后使用时间
func (o *Outfit) GetLastUsedAt() *time.Time {
	return o.lastUsedAt
}

// GetUseCount 获取使用次数
func (o *Outfit) GetUseCount() int {
	return o.useCount
}

// GetMetadata 获取元数据
func (o *Outfit) GetMetadata() map[string]interface{} {
	return o.metadata
}

// SetMetadata 设置元数据
func (o *Outfit) SetMetadata(key string, value interface{}) {
	o.metadata[key] = value
}

// GetMetadataValue 获取元数据值
func (o *Outfit) GetMetadataValue(key string) (interface{}, bool) {
	value, exists := o.metadata[key]
	return value, exists
}

// CanUpgrade 是否可以升级
func (o *Outfit) CanUpgrade() bool {
	return o.level < 100 && o.exp >= o.maxExp
}

// CanEnhance 是否可以强化
func (o *Outfit) CanEnhance() bool {
	return o.enhanceLevel < 20
}

// GetPower 获取战力值
func (o *Outfit) GetPower() int {
	power := 0
	totalAttrs := o.GetTotalAttributes()

	for _, value := range totalAttrs {
		power += value
	}

	// 等级加成
	power += o.level * 10

	// 强化加成
	power += o.enhanceLevel * 20

	return power
}

// Clone 克隆服装
func (o *Outfit) Clone() *Outfit {
	cloned := &Outfit{
		id:             uuid.New().String(), // 新ID
		name:           o.name,
		description:    o.description,
		outfitType:     o.outfitType,
		rarity:         o.rarity,
		quality:        o.quality,
		source:         o.source,
		attributes:     make(map[string]int),
		slots:          make([]OutfitSlot, len(o.slots)),
		isLocked:       o.isLocked,
		isEquipped:     false, // 克隆的不装备
		level:          o.level,
		exp:            o.exp,
		maxExp:         o.maxExp,
		tags:           make([]string, len(o.tags)),
		setID:          o.setID,
		appearance:     NewAppearanceConfig(),
		dyeColors:      make(map[string]*DyeColor),
		enhanceLevel:   o.enhanceLevel,
		enhanceBonuses: make(map[string]int),
		obtainedAt:     time.Now(),
		lastUsedAt:     nil,
		useCount:       0,
		metadata:       make(map[string]interface{}),
	}

	// 复制属性
	for attr, value := range o.attributes {
		cloned.attributes[attr] = value
	}

	// 复制槽位
	copy(cloned.slots, o.slots)

	// 复制标签
	copy(cloned.tags, o.tags)

	// 复制外观配置
	if o.appearance != nil {
		cloned.appearance = NewAppearanceConfig()
		for part, color := range o.appearance.GetColorScheme() {
			cloned.appearance.SetColor(part, color)
		}
		for _, effect := range o.appearance.GetEffects() {
			cloned.appearance.AddEffect(effect)
		}
		for _, animation := range o.appearance.GetAnimations() {
			cloned.appearance.AddAnimation(animation)
		}
		for part, texture := range o.appearance.GetTextures() {
			cloned.appearance.SetTexture(part, texture)
		}
		cloned.appearance.SetScale(o.appearance.GetScale())
		cloned.appearance.SetTransparency(o.appearance.GetTransparency())
		cloned.appearance.SetGlowIntensity(o.appearance.GetGlowIntensity())
	}

	// 复制染色
	for part, color := range o.dyeColors {
		cloned.dyeColors[part] = &DyeColor{
			colorID:    color.colorID,
			colorName:  color.colorName,
			hexValue:   color.hexValue,
			rarity:     color.rarity,
			isUnlocked: color.isUnlocked,
		}
	}

	// 复制强化加成
	for attr, bonus := range o.enhanceBonuses {
		cloned.enhanceBonuses[attr] = bonus
	}

	// 复制元数据
	for key, value := range o.metadata {
		cloned.metadata[key] = value
	}

	return cloned
}

// OutfitSet 套装实体
type OutfitSet struct {
	equippedOutfits map[OutfitSlot]*Outfit
	setBonuses      map[string]int
	fashionSets     map[string]*FashionSetBonus // 套装ID -> 套装加成信息
	styleBonus      *DressupStyle
	totalPower      int
	lastUpdated     time.Time
}

// FashionSetBonus 时装套装加成信息
type FashionSetBonus struct {
	setID         string
	setName       string
	equippedCount int
	totalCount    int
	activeBonus   map[string]int
}

// NewOutfitSet 创建套装
func NewOutfitSet() *OutfitSet {
	return &OutfitSet{
		equippedOutfits: make(map[OutfitSlot]*Outfit),
		setBonuses:      make(map[string]int),
		fashionSets:     make(map[string]*FashionSetBonus),
		styleBonus:      nil,
		totalPower:      0,
		lastUpdated:     time.Now(),
	}
}

// EquipToSlot 装备到槽位
func (os *OutfitSet) EquipToSlot(slot OutfitSlot, outfit *Outfit) {
	os.equippedOutfits[slot] = outfit
	os.calculateSetBonuses()
}

// UnequipFromSlot 从槽位卸下
func (os *OutfitSet) UnequipFromSlot(slot OutfitSlot) {
	delete(os.equippedOutfits, slot)
	os.calculateSetBonuses()
}

// GetEquippedOutfit 获取指定槽位的装备
func (os *OutfitSet) GetEquippedOutfit(slot OutfitSlot) *Outfit {
	return os.equippedOutfits[slot]
}

// GetAllEquipped 获取所有装备
func (os *OutfitSet) GetAllEquipped() map[OutfitSlot]*Outfit {
	return os.equippedOutfits
}

// GetSetBonuses 获取套装加成
func (os *OutfitSet) GetSetBonuses() map[string]int {
	return os.setBonuses
}

// calculateSetBonuses 计算套装加成
func (os *OutfitSet) calculateSetBonuses() {
	// 清空现有加成
	os.setBonuses = make(map[string]int)
	os.fashionSets = make(map[string]*FashionSetBonus)
	os.totalPower = 0

	// 统计各类型服装数量
	typeCount := make(map[OutfitType]int)
	setCount := make(map[string]int)    // 套装ID -> 装备数量
	setNames := make(map[string]string) // 套装ID -> 套装名称

	for _, outfit := range os.equippedOutfits {
		if outfit != nil {
			typeCount[outfit.GetType()]++

			// 统计套装
			if outfit.GetSetID() != "" {
				setCount[outfit.GetSetID()]++
				// 这里应该从套装配置中获取名称，暂时使用ID
				setNames[outfit.GetSetID()] = outfit.GetSetID()
			}

			// 累计战力
			os.totalPower += outfit.GetPower()
		}
	}

	// 根据类型数量计算基础加成
	for _, count := range typeCount {
		if count >= 2 {
			os.setBonuses["attack"] += count * 10
		}
		if count >= 4 {
			os.setBonuses["defense"] += count * 15
		}
		if count >= 6 {
			os.setBonuses["hp"] += count * 20
		}
	}

	// 计算时装套装加成
	for setID, equippedCount := range setCount {
		setBonus := &FashionSetBonus{
			setID:         setID,
			setName:       setNames[setID],
			equippedCount: equippedCount,
			totalCount:    6, // 假设每套装有6件，实际应该从配置获取
			activeBonus:   make(map[string]int),
		}

		// 根据装备数量计算套装加成
		if equippedCount >= 2 {
			setBonus.activeBonus["attack"] = equippedCount * 15
			os.setBonuses["attack"] += setBonus.activeBonus["attack"]
		}
		if equippedCount >= 4 {
			setBonus.activeBonus["defense"] = equippedCount * 20
			os.setBonuses["defense"] += setBonus.activeBonus["defense"]
		}
		if equippedCount >= 6 {
			setBonus.activeBonus["hp"] = equippedCount * 30
			setBonus.activeBonus["crit_rate"] = 10 // 暴击率+10%
			os.setBonuses["hp"] += setBonus.activeBonus["hp"]
			os.setBonuses["crit_rate"] += setBonus.activeBonus["crit_rate"]
		}

		os.fashionSets[setID] = setBonus
	}

	// 应用风格加成
	if os.styleBonus != nil {
		for attr, bonus := range os.styleBonus.GetBonuses() {
			os.setBonuses[attr] += bonus
		}
	}

	os.lastUpdated = time.Now()
}

// GetFashionSets 获取时装套装信息
func (os *OutfitSet) GetFashionSets() map[string]*FashionSetBonus {
	return os.fashionSets
}

// GetFashionSetBonus 获取指定套装的加成信息
func (os *OutfitSet) GetFashionSetBonus(setID string) *FashionSetBonus {
	return os.fashionSets[setID]
}

// SetStyleBonus 设置风格加成
func (os *OutfitSet) SetStyleBonus(style *DressupStyle) {
	os.styleBonus = style
	os.calculateSetBonuses() // 重新计算加成
}

// GetStyleBonus 获取风格加成
func (os *OutfitSet) GetStyleBonus() *DressupStyle {
	return os.styleBonus
}

// GetTotalPower 获取总战力
func (os *OutfitSet) GetTotalPower() int {
	return os.totalPower
}

// GetLastUpdated 获取最后更新时间
func (os *OutfitSet) GetLastUpdated() time.Time {
	return os.lastUpdated
}

// GetEquippedCount 获取已装备数量
func (os *OutfitSet) GetEquippedCount() int {
	count := 0
	for _, outfit := range os.equippedOutfits {
		if outfit != nil {
			count++
		}
	}
	return count
}

// GetEmptySlots 获取空槽位
func (os *OutfitSet) GetEmptySlots() []OutfitSlot {
	allSlots := []OutfitSlot{
		SlotWeapon, SlotArmor, SlotHelmet, SlotShoes,
		SlotRing, SlotNecklace, SlotFashionWeapon,
		SlotFashionArmor, SlotFashionHelmet, SlotPet, SlotMount,
	}

	emptySlots := make([]OutfitSlot, 0)
	for _, slot := range allSlots {
		if os.equippedOutfits[slot] == nil {
			emptySlots = append(emptySlots, slot)
		}
	}

	return emptySlots
}

// GetOutfitsByType 根据类型获取装备
func (os *OutfitSet) GetOutfitsByType(outfitType OutfitType) []*Outfit {
	outfits := make([]*Outfit, 0)
	for _, outfit := range os.equippedOutfits {
		if outfit != nil && outfit.GetType() == outfitType {
			outfits = append(outfits, outfit)
		}
	}
	return outfits
}

// GetOutfitsByRarity 根据稀有度获取装备
func (os *OutfitSet) GetOutfitsByRarity(rarity Rarity) []*Outfit {
	outfits := make([]*Outfit, 0)
	for _, outfit := range os.equippedOutfits {
		if outfit != nil && outfit.GetRarity() == rarity {
			outfits = append(outfits, outfit)
		}
	}
	return outfits
}

// GetTotalAttributes 获取套装总属性
func (os *OutfitSet) GetTotalAttributes() map[string]int {
	total := make(map[string]int)

	// 累计装备属性
	for _, outfit := range os.equippedOutfits {
		if outfit != nil {
			for attr, value := range outfit.GetTotalAttributes() {
				total[attr] += value
			}
		}
	}

	// 加上套装加成
	for attr, bonus := range os.setBonuses {
		total[attr] += bonus
	}

	return total
}

// CanEquipOutfit 检查是否可以装备服装
func (os *OutfitSet) CanEquipOutfit(outfit *Outfit, slot OutfitSlot) bool {
	if outfit == nil {
		return false
	}

	// 检查服装是否支持该槽位
	if !outfit.CanEquipToSlot(slot) {
		return false
	}

	// 检查是否已装备
	if outfit.IsEquipped() {
		return false
	}

	// 检查是否锁定
	if outfit.IsLocked() {
		return false
	}

	return true
}

// GetSetCompletionRate 获取套装完成度
func (os *OutfitSet) GetSetCompletionRate(setID string) float64 {
	setBonus := os.fashionSets[setID]
	if setBonus == nil {
		return 0.0
	}

	return float64(setBonus.equippedCount) / float64(setBonus.totalCount) * 100.0
}

// GetHighestQualityOutfit 获取品质最高的装备
func (os *OutfitSet) GetHighestQualityOutfit() *Outfit {
	var highest *Outfit
	highestQuality := QualityNormal

	for _, outfit := range os.equippedOutfits {
		if outfit != nil && outfit.GetQuality() > highestQuality {
			highest = outfit
			highestQuality = outfit.GetQuality()
		}
	}

	return highest
}

// GetAverageLevel 获取平均等级
func (os *OutfitSet) GetAverageLevel() float64 {
	totalLevel := 0
	count := 0

	for _, outfit := range os.equippedOutfits {
		if outfit != nil {
			totalLevel += outfit.GetLevel()
			count++
		}
	}

	if count == 0 {
		return 0.0
	}

	return float64(totalLevel) / float64(count)
}

// GetAverageEnhanceLevel 获取平均强化等级
func (os *OutfitSet) GetAverageEnhanceLevel() float64 {
	totalEnhance := 0
	count := 0

	for _, outfit := range os.equippedOutfits {
		if outfit != nil {
			totalEnhance += outfit.GetEnhanceLevel()
			count++
		}
	}

	if count == 0 {
		return 0.0
	}

	return float64(totalEnhance) / float64(count)
}

// Clone 克隆套装
func (os *OutfitSet) Clone() *OutfitSet {
	cloned := NewOutfitSet()

	// 复制装备（注意：这里复制引用，不是深拷贝装备本身）
	for slot, outfit := range os.equippedOutfits {
		cloned.equippedOutfits[slot] = outfit
	}

	// 复制套装加成
	for attr, bonus := range os.setBonuses {
		cloned.setBonuses[attr] = bonus
	}

	// 复制时装套装信息
	for setID, setBonus := range os.fashionSets {
		cloned.fashionSets[setID] = &FashionSetBonus{
			setID:         setBonus.setID,
			setName:       setBonus.setName,
			equippedCount: setBonus.equippedCount,
			totalCount:    setBonus.totalCount,
			activeBonus:   make(map[string]int),
		}

		// 复制激活加成
		for attr, bonus := range setBonus.activeBonus {
			cloned.fashionSets[setID].activeBonus[attr] = bonus
		}
	}

	// 复制风格加成
	cloned.styleBonus = os.styleBonus
	cloned.totalPower = os.totalPower
	cloned.lastUpdated = os.lastUpdated

	return cloned
}
