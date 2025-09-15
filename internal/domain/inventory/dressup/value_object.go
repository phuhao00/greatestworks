package dressup

// OutfitType 服装类型
type OutfitType int

const (
	OutfitTypeWeapon OutfitType = iota + 1
	OutfitTypeArmor
	OutfitTypeHelmet
	OutfitTypeShoes
	OutfitTypeAccessory
	OutfitTypeFashion
	OutfitTypePet
	OutfitTypeMount
)

// String 返回服装类型字符串
func (ot OutfitType) String() string {
	switch ot {
	case OutfitTypeWeapon:
		return "weapon"
	case OutfitTypeArmor:
		return "armor"
	case OutfitTypeHelmet:
		return "helmet"
	case OutfitTypeShoes:
		return "shoes"
	case OutfitTypeAccessory:
		return "accessory"
	case OutfitTypeFashion:
		return "fashion"
	case OutfitTypePet:
		return "pet"
	case OutfitTypeMount:
		return "mount"
	default:
		return "unknown"
	}
}

// OutfitSlot 装备槽位
type OutfitSlot int

const (
	SlotWeapon OutfitSlot = iota + 1
	SlotArmor
	SlotHelmet
	SlotShoes
	SlotRing
	SlotNecklace
	SlotFashionWeapon
	SlotFashionArmor
	SlotFashionHelmet
	SlotPet
	SlotMount
)

// String 返回槽位字符串
func (os OutfitSlot) String() string {
	switch os {
	case SlotWeapon:
		return "weapon"
	case SlotArmor:
		return "armor"
	case SlotHelmet:
		return "helmet"
	case SlotShoes:
		return "shoes"
	case SlotRing:
		return "ring"
	case SlotNecklace:
		return "necklace"
	case SlotFashionWeapon:
		return "fashion_weapon"
	case SlotFashionArmor:
		return "fashion_armor"
	case SlotFashionHelmet:
		return "fashion_helmet"
	case SlotPet:
		return "pet"
	case SlotMount:
		return "mount"
	default:
		return "unknown"
	}
}

// Rarity 稀有度
type Rarity int

const (
	RarityCommon Rarity = iota + 1
	RarityUncommon
	RarityRare
	RarityEpic
	RarityLegendary
	RarityMythic
)

// String 返回稀有度字符串
func (r Rarity) String() string {
	switch r {
	case RarityCommon:
		return "common"
	case RarityUncommon:
		return "uncommon"
	case RarityRare:
		return "rare"
	case RarityEpic:
		return "epic"
	case RarityLegendary:
		return "legendary"
	case RarityMythic:
		return "mythic"
	default:
		return "unknown"
	}
}

// GetRarityMultiplier 获取稀有度属性倍数
func (r Rarity) GetRarityMultiplier() float64 {
	switch r {
	case RarityCommon:
		return 1.0
	case RarityUncommon:
		return 1.2
	case RarityRare:
		return 1.5
	case RarityEpic:
		return 2.0
	case RarityLegendary:
		return 3.0
	case RarityMythic:
		return 5.0
	default:
		return 1.0
	}
}

// DressupStyle 换装风格
type DressupStyle struct {
	styleID   string
	styleName string
	theme     string
	bonuses   map[string]int
}

// NewDressupStyle 创建换装风格
func NewDressupStyle(styleID, styleName, theme string) *DressupStyle {
	return &DressupStyle{
		styleID:   styleID,
		styleName: styleName,
		theme:     theme,
		bonuses:   make(map[string]int),
	}
}

// GetStyleID 获取风格ID
func (ds *DressupStyle) GetStyleID() string {
	return ds.styleID
}

// GetStyleName 获取风格名称
func (ds *DressupStyle) GetStyleName() string {
	return ds.styleName
}

// GetTheme 获取主题
func (ds *DressupStyle) GetTheme() string {
	return ds.theme
}

// AddBonus 添加风格加成
func (ds *DressupStyle) AddBonus(attr string, value int) {
	ds.bonuses[attr] = value
}

// GetBonuses 获取所有加成
func (ds *DressupStyle) GetBonuses() map[string]int {
	return ds.bonuses
}

// OutfitQuality 服装品质
type OutfitQuality int

const (
	QualityNormal OutfitQuality = iota + 1 // 普通
	QualityGood                            // 良好
	QualityExcellent                       // 优秀
	QualityPerfect                         // 完美
	QualityMasterwork                      // 大师级
)

// String 返回品质字符串
func (oq OutfitQuality) String() string {
	switch oq {
	case QualityNormal:
		return "normal"
	case QualityGood:
		return "good"
	case QualityExcellent:
		return "excellent"
	case QualityPerfect:
		return "perfect"
	case QualityMasterwork:
		return "masterwork"
	default:
		return "unknown"
	}
}

// GetQualityMultiplier 获取品质属性倍数
func (oq OutfitQuality) GetQualityMultiplier() float64 {
	switch oq {
	case QualityNormal:
		return 1.0
	case QualityGood:
		return 1.1
	case QualityExcellent:
		return 1.25
	case QualityPerfect:
		return 1.5
	case QualityMasterwork:
		return 2.0
	default:
		return 1.0
	}
}

// OutfitSource 服装来源
type OutfitSource int

const (
	SourceShop     OutfitSource = iota + 1 // 商店购买
	SourceCraft                            // 制作
	SourceDrop                             // 掉落
	SourceEvent                            // 活动
	SourceGift                             // 礼品
	SourceAchievement                      // 成就
	SourceVIP                              // VIP
	SourceLimitedTime                      // 限时
)

// String 返回来源字符串
func (os OutfitSource) String() string {
	switch os {
	case SourceShop:
		return "shop"
	case SourceCraft:
		return "craft"
	case SourceDrop:
		return "drop"
	case SourceEvent:
		return "event"
	case SourceGift:
		return "gift"
	case SourceAchievement:
		return "achievement"
	case SourceVIP:
		return "vip"
	case SourceLimitedTime:
		return "limited_time"
	default:
		return "unknown"
	}
}

// FashionSet 时装套装
type FashionSet struct {
	setID       string
	setName     string
	description string
	pieces      []string // 套装部件ID列表
	setBonuses  map[int]map[string]int // 件数 -> 属性加成
	theme       string
	season      string
	isLimited   bool
}

// NewFashionSet 创建时装套装
func NewFashionSet(setID, setName, description string) *FashionSet {
	return &FashionSet{
		setID:       setID,
		setName:     setName,
		description: description,
		pieces:      make([]string, 0),
		setBonuses:  make(map[int]map[string]int),
		isLimited:   false,
	}
}

// GetSetID 获取套装ID
func (fs *FashionSet) GetSetID() string {
	return fs.setID
}

// GetSetName 获取套装名称
func (fs *FashionSet) GetSetName() string {
	return fs.setName
}

// GetDescription 获取描述
func (fs *FashionSet) GetDescription() string {
	return fs.description
}

// AddPiece 添加套装部件
func (fs *FashionSet) AddPiece(pieceID string) {
	fs.pieces = append(fs.pieces, pieceID)
}

// GetPieces 获取套装部件
func (fs *FashionSet) GetPieces() []string {
	return fs.pieces
}

// AddSetBonus 添加套装加成
func (fs *FashionSet) AddSetBonus(pieceCount int, attribute string, value int) {
	if fs.setBonuses[pieceCount] == nil {
		fs.setBonuses[pieceCount] = make(map[string]int)
	}
	fs.setBonuses[pieceCount][attribute] = value
}

// GetSetBonuses 获取套装加成
func (fs *FashionSet) GetSetBonuses() map[int]map[string]int {
	return fs.setBonuses
}

// GetBonusForPieceCount 获取指定件数的加成
func (fs *FashionSet) GetBonusForPieceCount(pieceCount int) map[string]int {
	return fs.setBonuses[pieceCount]
}

// SetTheme 设置主题
func (fs *FashionSet) SetTheme(theme string) {
	fs.theme = theme
}

// GetTheme 获取主题
func (fs *FashionSet) GetTheme() string {
	return fs.theme
}

// SetSeason 设置季节
func (fs *FashionSet) SetSeason(season string) {
	fs.season = season
}

// GetSeason 获取季节
func (fs *FashionSet) GetSeason() string {
	return fs.season
}

// SetLimited 设置限定状态
func (fs *FashionSet) SetLimited(limited bool) {
	fs.isLimited = limited
}

// IsLimited 是否限定
func (fs *FashionSet) IsLimited() bool {
	return fs.isLimited
}

// AppearanceConfig 外观配置
type AppearanceConfig struct {
	colorScheme   map[string]string // 颜色方案
	effects       []string          // 特效列表
	animations    []string          // 动画列表
	textures      map[string]string // 材质贴图
	scale         float64           // 缩放比例
	transparency  float64           // 透明度
	glowIntensity float64           // 发光强度
}

// NewAppearanceConfig 创建外观配置
func NewAppearanceConfig() *AppearanceConfig {
	return &AppearanceConfig{
		colorScheme:   make(map[string]string),
		effects:       make([]string, 0),
		animations:    make([]string, 0),
		textures:      make(map[string]string),
		scale:         1.0,
		transparency:  1.0,
		glowIntensity: 0.0,
	}
}

// SetColor 设置颜色
func (ac *AppearanceConfig) SetColor(part, color string) {
	ac.colorScheme[part] = color
}

// GetColor 获取颜色
func (ac *AppearanceConfig) GetColor(part string) string {
	return ac.colorScheme[part]
}

// GetColorScheme 获取颜色方案
func (ac *AppearanceConfig) GetColorScheme() map[string]string {
	return ac.colorScheme
}

// AddEffect 添加特效
func (ac *AppearanceConfig) AddEffect(effect string) {
	ac.effects = append(ac.effects, effect)
}

// GetEffects 获取特效列表
func (ac *AppearanceConfig) GetEffects() []string {
	return ac.effects
}

// AddAnimation 添加动画
func (ac *AppearanceConfig) AddAnimation(animation string) {
	ac.animations = append(ac.animations, animation)
}

// GetAnimations 获取动画列表
func (ac *AppearanceConfig) GetAnimations() []string {
	return ac.animations
}

// SetTexture 设置材质
func (ac *AppearanceConfig) SetTexture(part, texture string) {
	ac.textures[part] = texture
}

// GetTexture 获取材质
func (ac *AppearanceConfig) GetTexture(part string) string {
	return ac.textures[part]
}

// GetTextures 获取所有材质
func (ac *AppearanceConfig) GetTextures() map[string]string {
	return ac.textures
}

// SetScale 设置缩放
func (ac *AppearanceConfig) SetScale(scale float64) {
	ac.scale = scale
}

// GetScale 获取缩放
func (ac *AppearanceConfig) GetScale() float64 {
	return ac.scale
}

// SetTransparency 设置透明度
func (ac *AppearanceConfig) SetTransparency(transparency float64) {
	ac.transparency = transparency
}

// GetTransparency 获取透明度
func (ac *AppearanceConfig) GetTransparency() float64 {
	return ac.transparency
}

// SetGlowIntensity 设置发光强度
func (ac *AppearanceConfig) SetGlowIntensity(intensity float64) {
	ac.glowIntensity = intensity
}

// GetGlowIntensity 获取发光强度
func (ac *AppearanceConfig) GetGlowIntensity() float64 {
	return ac.glowIntensity
}

// AttributeBonus 属性加成
type AttributeBonus struct {
	attribute string  // 属性名称
	baseValue int     // 基础值
	bonus     float64 // 加成倍数
	bonusType string  // 加成类型：percentage, fixed
}

// NewAttributeBonus 创建属性加成
func NewAttributeBonus(attribute string, baseValue int, bonus float64, bonusType string) *AttributeBonus {
	return &AttributeBonus{
		attribute: attribute,
		baseValue: baseValue,
		bonus:     bonus,
		bonusType: bonusType,
	}
}

// GetAttribute 获取属性名称
func (ab *AttributeBonus) GetAttribute() string {
	return ab.attribute
}

// GetBaseValue 获取基础值
func (ab *AttributeBonus) GetBaseValue() int {
	return ab.baseValue
}

// GetBonus 获取加成倍数
func (ab *AttributeBonus) GetBonus() float64 {
	return ab.bonus
}

// GetBonusType 获取加成类型
func (ab *AttributeBonus) GetBonusType() string {
	return ab.bonusType
}

// CalculateFinalValue 计算最终值
func (ab *AttributeBonus) CalculateFinalValue() int {
	switch ab.bonusType {
	case "percentage":
		return int(float64(ab.baseValue) * (1.0 + ab.bonus))
	case "fixed":
		return ab.baseValue + int(ab.bonus)
	default:
		return ab.baseValue
	}
}

// DyeColor 染色颜色
type DyeColor struct {
	colorID   string
	colorName string
	hexValue  string
	rarity    Rarity
	isUnlocked bool
}

// NewDyeColor 创建染色颜色
func NewDyeColor(colorID, colorName, hexValue string, rarity Rarity) *DyeColor {
	return &DyeColor{
		colorID:   colorID,
		colorName: colorName,
		hexValue:  hexValue,
		rarity:    rarity,
		isUnlocked: false,
	}
}

// GetColorID 获取颜色ID
func (dc *DyeColor) GetColorID() string {
	return dc.colorID
}

// GetColorName 获取颜色名称
func (dc *DyeColor) GetColorName() string {
	return dc.colorName
}

// GetHexValue 获取十六进制值
func (dc *DyeColor) GetHexValue() string {
	return dc.hexValue
}

// GetRarity 获取稀有度
func (dc *DyeColor) GetRarity() Rarity {
	return dc.rarity
}

// Unlock 解锁颜色
func (dc *DyeColor) Unlock() {
	dc.isUnlocked = true
}

// IsUnlocked 是否已解锁
func (dc *DyeColor) IsUnlocked() bool {
	return dc.isUnlocked
}

// OutfitFilter 服装筛选器
type OutfitFilter struct {
	outfitType   *OutfitType
	rarity       *Rarity
	quality      *OutfitQuality
	source       *OutfitSource
	slot         *OutfitSlot
	isLocked     *bool
	hasSetBonus  *bool
	minLevel     *int
	maxLevel     *int
	searchText   string
	tags         []string
}

// NewOutfitFilter 创建服装筛选器
func NewOutfitFilter() *OutfitFilter {
	return &OutfitFilter{
		tags: make([]string, 0),
	}
}

// SetOutfitType 设置服装类型筛选
func (of *OutfitFilter) SetOutfitType(outfitType OutfitType) {
	of.outfitType = &outfitType
}

// SetRarity 设置稀有度筛选
func (of *OutfitFilter) SetRarity(rarity Rarity) {
	of.rarity = &rarity
}

// SetQuality 设置品质筛选
func (of *OutfitFilter) SetQuality(quality OutfitQuality) {
	of.quality = &quality
}

// SetSource 设置来源筛选
func (of *OutfitFilter) SetSource(source OutfitSource) {
	of.source = &source
}

// SetSlot 设置槽位筛选
func (of *OutfitFilter) SetSlot(slot OutfitSlot) {
	of.slot = &slot
}

// SetLocked 设置锁定状态筛选
func (of *OutfitFilter) SetLocked(locked bool) {
	of.isLocked = &locked
}

// SetHasSetBonus 设置套装加成筛选
func (of *OutfitFilter) SetHasSetBonus(hasSetBonus bool) {
	of.hasSetBonus = &hasSetBonus
}

// SetLevelRange 设置等级范围筛选
func (of *OutfitFilter) SetLevelRange(minLevel, maxLevel int) {
	of.minLevel = &minLevel
	of.maxLevel = &maxLevel
}

// SetSearchText 设置搜索文本
func (of *OutfitFilter) SetSearchText(text string) {
	of.searchText = text
}

// AddTag 添加标签筛选
func (of *OutfitFilter) AddTag(tag string) {
	of.tags = append(of.tags, tag)
}

// GetOutfitType 获取服装类型筛选
func (of *OutfitFilter) GetOutfitType() *OutfitType {
	return of.outfitType
}

// GetRarity 获取稀有度筛选
func (of *OutfitFilter) GetRarity() *Rarity {
	return of.rarity
}

// GetQuality 获取品质筛选
func (of *OutfitFilter) GetQuality() *OutfitQuality {
	return of.quality
}

// GetSource 获取来源筛选
func (of *OutfitFilter) GetSource() *OutfitSource {
	return of.source
}

// GetSlot 获取槽位筛选
func (of *OutfitFilter) GetSlot() *OutfitSlot {
	return of.slot
}

// GetLocked 获取锁定状态筛选
func (of *OutfitFilter) GetLocked() *bool {
	return of.isLocked
}

// GetHasSetBonus 获取套装加成筛选
func (of *OutfitFilter) GetHasSetBonus() *bool {
	return of.hasSetBonus
}

// GetMinLevel 获取最小等级
func (of *OutfitFilter) GetMinLevel() *int {
	return of.minLevel
}

// GetMaxLevel 获取最大等级
func (of *OutfitFilter) GetMaxLevel() *int {
	return of.maxLevel
}

// GetSearchText 获取搜索文本
func (of *OutfitFilter) GetSearchText() string {
	return of.searchText
}

// GetTags 获取标签列表
func (of *OutfitFilter) GetTags() []string {
	return of.tags
}