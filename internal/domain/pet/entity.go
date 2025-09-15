package pet

import (
	"fmt"
	"time"
)

// PetFragment 宠物碎片实体
type PetFragment struct {
	id           string
	playerID     string
	fragmentID   uint32
	relatedPetID uint32
	quantity     uint64
	createdAt    time.Time
	updatedAt    time.Time
}

// NewPetFragment 创建新的宠物碎片
func NewPetFragment(playerID string, fragmentID, relatedPetID uint32, quantity uint64) *PetFragment {
	now := time.Now()
	return &PetFragment{
		id:           fmt.Sprintf("fragment_%d", now.UnixNano()),
		playerID:     playerID,
		fragmentID:   fragmentID,
		relatedPetID: relatedPetID,
		quantity:     quantity,
		createdAt:    now,
		updatedAt:    now,
	}
}

// GetID 获取碎片ID
func (pf *PetFragment) GetID() string {
	return pf.id
}

// GetPlayerID 获取玩家ID
func (pf *PetFragment) GetPlayerID() string {
	return pf.playerID
}

// GetFragmentID 获取碎片配置ID
func (pf *PetFragment) GetFragmentID() uint32 {
	return pf.fragmentID
}

// GetRelatedPetID 获取关联宠物ID
func (pf *PetFragment) GetRelatedPetID() uint32 {
	return pf.relatedPetID
}

// GetQuantity 获取数量
func (pf *PetFragment) GetQuantity() uint64 {
	return pf.quantity
}

// AddQuantity 增加数量
func (pf *PetFragment) AddQuantity(amount uint64) {
	pf.quantity += amount
	pf.updatedAt = time.Now()
}

// ConsumeQuantity 消耗数量
func (pf *PetFragment) ConsumeQuantity(amount uint64) error {
	if pf.quantity < amount {
		return ErrInsufficientFragments
	}
	
	pf.quantity -= amount
	pf.updatedAt = time.Now()
	return nil
}

// CanSummon 是否可以召唤宠物
func (pf *PetFragment) CanSummon(requiredQuantity uint64) bool {
	return pf.quantity >= requiredQuantity
}

// GetCreatedAt 获取创建时间
func (pf *PetFragment) GetCreatedAt() time.Time {
	return pf.createdAt
}

// GetUpdatedAt 获取更新时间
func (pf *PetFragment) GetUpdatedAt() time.Time {
	return pf.updatedAt
}

// PetSkin 宠物皮肤实体
type PetSkin struct {
	id        string
	skinID    string
	name      string
	rarity    PetRarity
	equipped  bool
	powerBonus int64
	attributeBonus map[string]float64
	unlocked  bool
	unlockTime time.Time
	createdAt time.Time
	updatedAt time.Time
}

// NewPetSkin 创建新的宠物皮肤
func NewPetSkin(skinID, name string, rarity PetRarity, powerBonus int64) *PetSkin {
	now := time.Now()
	return &PetSkin{
		id:             fmt.Sprintf("skin_%d", now.UnixNano()),
		skinID:         skinID,
		name:           name,
		rarity:         rarity,
		equipped:       false,
		powerBonus:     powerBonus,
		attributeBonus: make(map[string]float64),
		unlocked:       false,
		createdAt:      now,
		updatedAt:      now,
	}
}

// GetID 获取皮肤实体ID
func (ps *PetSkin) GetID() string {
	return ps.id
}

// GetSkinID 获取皮肤配置ID
func (ps *PetSkin) GetSkinID() string {
	return ps.skinID
}

// GetName 获取皮肤名称
func (ps *PetSkin) GetName() string {
	return ps.name
}

// GetRarity 获取稀有度
func (ps *PetSkin) GetRarity() PetRarity {
	return ps.rarity
}

// IsEquipped 是否已装备
func (ps *PetSkin) IsEquipped() bool {
	return ps.equipped
}

// GetPowerBonus 获取战力加成
func (ps *PetSkin) GetPowerBonus() int64 {
	return ps.powerBonus
}

// GetAttributeBonus 获取属性加成
func (ps *PetSkin) GetAttributeBonus() map[string]float64 {
	return ps.attributeBonus
}

// IsUnlocked 是否已解锁
func (ps *PetSkin) IsUnlocked() bool {
	return ps.unlocked
}

// Unlock 解锁皮肤
func (ps *PetSkin) Unlock() error {
	if ps.unlocked {
		return ErrSkinAlreadyUnlocked
	}
	
	ps.unlocked = true
	ps.unlockTime = time.Now()
	ps.updatedAt = time.Now()
	return nil
}

// Equip 装备皮肤
func (ps *PetSkin) Equip() error {
	if !ps.unlocked {
		return ErrSkinNotUnlocked
	}
	
	if ps.equipped {
		return ErrSkinAlreadyEquipped
	}
	
	ps.equipped = true
	ps.updatedAt = time.Now()
	return nil
}

// Unequip 卸下皮肤
func (ps *PetSkin) Unequip() {
	ps.equipped = false
	ps.updatedAt = time.Now()
}

// SetAttributeBonus 设置属性加成
func (ps *PetSkin) SetAttributeBonus(attribute string, bonus float64) {
	ps.attributeBonus[attribute] = bonus
	ps.updatedAt = time.Now()
}

// GetUnlockTime 获取解锁时间
func (ps *PetSkin) GetUnlockTime() time.Time {
	return ps.unlockTime
}

// GetCreatedAt 获取创建时间
func (ps *PetSkin) GetCreatedAt() time.Time {
	return ps.createdAt
}

// GetUpdatedAt 获取更新时间
func (ps *PetSkin) GetUpdatedAt() time.Time {
	return ps.updatedAt
}

// PetSkill 宠物技能实体
type PetSkill struct {
	id          string
	skillID     string
	name        string
	level       uint32
	experience  uint64
	cooldown    time.Duration
	lastUsed    time.Time
	skillType   SkillType
	damage      int64
	effects     []SkillEffect
	createdAt   time.Time
	updatedAt   time.Time
}

// SkillType 技能类型
type SkillType int

const (
	SkillTypeAttack  SkillType = 1 // 攻击技能
	SkillTypeDefense SkillType = 2 // 防御技能
	SkillTypeHeal    SkillType = 3 // 治疗技能
	SkillTypeBuff    SkillType = 4 // 增益技能
	SkillTypeDebuff  SkillType = 5 // 减益技能
)

// SkillEffect 技能效果
type SkillEffect struct {
	EffectType string
	Value      float64
	Duration   time.Duration
}

// NewPetSkill 创建新的宠物技能
func NewPetSkill(skillID, name string, skillType SkillType, damage int64, cooldown time.Duration) *PetSkill {
	now := time.Now()
	return &PetSkill{
		id:        fmt.Sprintf("skill_%d", now.UnixNano()),
		skillID:   skillID,
		name:      name,
		level:     1,
		experience: 0,
		cooldown:  cooldown,
		skillType: skillType,
		damage:    damage,
		effects:   make([]SkillEffect, 0),
		createdAt: now,
		updatedAt: now,
	}
}

// GetID 获取技能实体ID
func (ps *PetSkill) GetID() string {
	return ps.id
}

// GetSkillID 获取技能配置ID
func (ps *PetSkill) GetSkillID() string {
	return ps.skillID
}

// GetName 获取技能名称
func (ps *PetSkill) GetName() string {
	return ps.name
}

// GetLevel 获取技能等级
func (ps *PetSkill) GetLevel() uint32 {
	return ps.level
}

// GetExperience 获取技能经验
func (ps *PetSkill) GetExperience() uint64 {
	return ps.experience
}

// GetCooldown 获取冷却时间
func (ps *PetSkill) GetCooldown() time.Duration {
	return ps.cooldown
}

// GetLastUsed 获取上次使用时间
func (ps *PetSkill) GetLastUsed() time.Time {
	return ps.lastUsed
}

// GetSkillType 获取技能类型
func (ps *PetSkill) GetSkillType() SkillType {
	return ps.skillType
}

// GetDamage 获取技能伤害
func (ps *PetSkill) GetDamage() int64 {
	return ps.damage
}

// GetEffects 获取技能效果
func (ps *PetSkill) GetEffects() []SkillEffect {
	return ps.effects
}

// AddExperience 增加技能经验
func (ps *PetSkill) AddExperience(exp uint64) {
	ps.experience += exp
	ps.updatedAt = time.Now()
	
	// 检查是否可以升级
	if ps.canLevelUp() {
		ps.levelUp()
	}
}

// canLevelUp 检查是否可以升级
func (ps *PetSkill) canLevelUp() bool {
	requiredExp := ps.calculateRequiredExperience()
	return ps.experience >= requiredExp && ps.level < MaxSkillLevel
}

// levelUp 技能升级
func (ps *PetSkill) levelUp() {
	ps.level++
	ps.damage = int64(float64(ps.damage) * 1.1) // 每级增加10%伤害
	ps.updatedAt = time.Now()
}

// calculateRequiredExperience 计算升级所需经验
func (ps *PetSkill) calculateRequiredExperience() uint64 {
	return uint64(ps.level * ps.level * 50)
}

// Upgrade 升级技能
func (ps *PetSkill) Upgrade() error {
	if ps.level >= MaxSkillLevel {
		return ErrMaxSkillLevelReached
	}
	
	if !ps.canLevelUp() {
		return ErrInsufficientSkillExperience
	}
	
	ps.levelUp()
	return nil
}

// Use 使用技能
func (ps *PetSkill) Use() error {
	if !ps.IsReady() {
		return ErrSkillOnCooldown
	}
	
	ps.lastUsed = time.Now()
	ps.updatedAt = time.Now()
	return nil
}

// IsReady 技能是否准备就绪
func (ps *PetSkill) IsReady() bool {
	return time.Since(ps.lastUsed) >= ps.cooldown
}

// GetRemainingCooldown 获取剩余冷却时间
func (ps *PetSkill) GetRemainingCooldown() time.Duration {
	elapsed := time.Since(ps.lastUsed)
	if elapsed >= ps.cooldown {
		return 0
	}
	return ps.cooldown - elapsed
}

// AddEffect 添加技能效果
func (ps *PetSkill) AddEffect(effect SkillEffect) {
	ps.effects = append(ps.effects, effect)
	ps.updatedAt = time.Now()
}

// GetCreatedAt 获取创建时间
func (ps *PetSkill) GetCreatedAt() time.Time {
	return ps.createdAt
}

// GetUpdatedAt 获取更新时间
func (ps *PetSkill) GetUpdatedAt() time.Time {
	return ps.updatedAt
}

// PetBonds 宠物羁绊实体
type PetBonds struct {
	id          string
	activeBonds []string
	bondEffects map[string]*BondEffect
	createdAt   time.Time
	updatedAt   time.Time
}

// BondEffect 羁绊效果
type BondEffect struct {
	BondID      string
	Name        string
	Description string
	PowerBonus  int64
	Attributes  map[string]float64
	Active      bool
	ActivatedAt time.Time
}

// NewPetBonds 创建新的宠物羁绊
func NewPetBonds() *PetBonds {
	now := time.Now()
	return &PetBonds{
		id:          fmt.Sprintf("bonds_%d", now.UnixNano()),
		activeBonds: make([]string, 0),
		bondEffects: make(map[string]*BondEffect),
		createdAt:   now,
		updatedAt:   now,
	}
}

// GetID 获取羁绊ID
func (pb *PetBonds) GetID() string {
	return pb.id
}

// GetActiveBonds 获取激活的羁绊
func (pb *PetBonds) GetActiveBonds() []string {
	return pb.activeBonds
}

// GetBondEffects 获取羁绊效果
func (pb *PetBonds) GetBondEffects() map[string]*BondEffect {
	return pb.bondEffects
}

// ActivateBond 激活羁绊
func (pb *PetBonds) ActivateBond(bondID string) error {
	// 检查是否已激活
	for _, activeBond := range pb.activeBonds {
		if activeBond == bondID {
			return ErrBondAlreadyActive
		}
	}
	
	// 检查是否达到最大激活数量
	if len(pb.activeBonds) >= MaxActiveBonds {
		return ErrMaxActiveBondsReached
	}
	
	pb.activeBonds = append(pb.activeBonds, bondID)
	
	// 激活羁绊效果
	if effect, exists := pb.bondEffects[bondID]; exists {
		effect.Active = true
		effect.ActivatedAt = time.Now()
	}
	
	pb.updatedAt = time.Now()
	return nil
}

// DeactivateBond 取消羁绊
func (pb *PetBonds) DeactivateBond(bondID string) error {
	// 查找并移除激活的羁绊
	for i, activeBond := range pb.activeBonds {
		if activeBond == bondID {
			pb.activeBonds = append(pb.activeBonds[:i], pb.activeBonds[i+1:]...)
			
			// 取消羁绊效果
			if effect, exists := pb.bondEffects[bondID]; exists {
				effect.Active = false
			}
			
			pb.updatedAt = time.Now()
			return nil
		}
	}
	
	return ErrBondNotActive
}

// AddBondEffect 添加羁绊效果
func (pb *PetBonds) AddBondEffect(effect *BondEffect) {
	pb.bondEffects[effect.BondID] = effect
	pb.updatedAt = time.Now()
}

// GetPowerBonus 获取羁绊战力加成
func (pb *PetBonds) GetPowerBonus() int64 {
	var totalBonus int64
	for _, bondID := range pb.activeBonds {
		if effect, exists := pb.bondEffects[bondID]; exists && effect.Active {
			totalBonus += effect.PowerBonus
		}
	}
	return totalBonus
}

// GetAttributeBonus 获取羁绊属性加成
func (pb *PetBonds) GetAttributeBonus() map[string]float64 {
	attributeBonus := make(map[string]float64)
	
	for _, bondID := range pb.activeBonds {
		if effect, exists := pb.bondEffects[bondID]; exists && effect.Active {
			for attr, bonus := range effect.Attributes {
				attributeBonus[attr] += bonus
			}
		}
	}
	
	return attributeBonus
}

// IsBondActive 检查羁绊是否激活
func (pb *PetBonds) IsBondActive(bondID string) bool {
	for _, activeBond := range pb.activeBonds {
		if activeBond == bondID {
			return true
		}
	}
	return false
}

// GetActiveCount 获取激活羁绊数量
func (pb *PetBonds) GetActiveCount() int {
	return len(pb.activeBonds)
}

// GetCreatedAt 获取创建时间
func (pb *PetBonds) GetCreatedAt() time.Time {
	return pb.createdAt
}

// GetUpdatedAt 获取更新时间
func (pb *PetBonds) GetUpdatedAt() time.Time {
	return pb.updatedAt
}

// PetPictorial 宠物图鉴实体
type PetPictorial struct {
	id           string
	playerID     string
	petConfigID  uint32
	unlocked     bool
	highestLevel uint32
	highestStar  uint32
	firstSeen    time.Time
	lastSeen     time.Time
	createdAt    time.Time
	updatedAt    time.Time
}

// NewPetPictorial 创建新的宠物图鉴
func NewPetPictorial(playerID string, petConfigID uint32) *PetPictorial {
	now := time.Now()
	return &PetPictorial{
		id:          fmt.Sprintf("pictorial_%d", now.UnixNano()),
		playerID:    playerID,
		petConfigID: petConfigID,
		unlocked:    false,
		createdAt:   now,
		updatedAt:   now,
	}
}

// GetID 获取图鉴ID
func (pp *PetPictorial) GetID() string {
	return pp.id
}

// GetPlayerID 获取玩家ID
func (pp *PetPictorial) GetPlayerID() string {
	return pp.playerID
}

// GetPetConfigID 获取宠物配置ID
func (pp *PetPictorial) GetPetConfigID() uint32 {
	return pp.petConfigID
}

// IsUnlocked 是否已解锁
func (pp *PetPictorial) IsUnlocked() bool {
	return pp.unlocked
}

// GetHighestLevel 获取最高等级
func (pp *PetPictorial) GetHighestLevel() uint32 {
	return pp.highestLevel
}

// GetHighestStar 获取最高星级
func (pp *PetPictorial) GetHighestStar() uint32 {
	return pp.highestStar
}

// Unlock 解锁图鉴
func (pp *PetPictorial) Unlock() {
	if !pp.unlocked {
		pp.unlocked = true
		pp.firstSeen = time.Now()
	}
	pp.lastSeen = time.Now()
	pp.updatedAt = time.Now()
}

// UpdateRecord 更新记录
func (pp *PetPictorial) UpdateRecord(level, star uint32) {
	if level > pp.highestLevel {
		pp.highestLevel = level
	}
	if star > pp.highestStar {
		pp.highestStar = star
	}
	pp.lastSeen = time.Now()
	pp.updatedAt = time.Now()
}

// GetFirstSeen 获取首次见到时间
func (pp *PetPictorial) GetFirstSeen() time.Time {
	return pp.firstSeen
}

// GetLastSeen 获取最后见到时间
func (pp *PetPictorial) GetLastSeen() time.Time {
	return pp.lastSeen
}

// GetCreatedAt 获取创建时间
func (pp *PetPictorial) GetCreatedAt() time.Time {
	return pp.createdAt
}

// GetUpdatedAt 获取更新时间
func (pp *PetPictorial) GetUpdatedAt() time.Time {
	return pp.updatedAt
}

// 常量定义
const (
	MaxSkillLevel    = 10
	MaxActiveBonds   = 3
)