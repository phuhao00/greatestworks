package pet

import (
	"fmt"
	"time"
)

// PetAggregate 宠物聚合根
type PetAggregate struct {
	id         string
	playerID   string
	configID   uint32
	name       string
	category   PetCategory
	star       uint32
	level      uint32
	experience uint64
	state      PetState
	attributes *PetAttributes
	skills     []*PetSkill
	bonds      *PetBonds
	skins      []*PetSkin
	reviveTime time.Time
	createdAt  time.Time
	updatedAt  time.Time
	version    int
}

// NewPetAggregate 创建新宠物聚合根
func NewPetAggregate(playerID string, configID uint32, name string, category PetCategory) *PetAggregate {
	now := time.Now()
	return &PetAggregate{
		id:         fmt.Sprintf("pet_%d", now.UnixNano()),
		playerID:   playerID,
		configID:   configID,
		name:       name,
		category:   category,
		star:       1,
		level:      1,
		experience: 0,
		state:      PetStateIdle,
		attributes: NewPetAttributes(),
		skills:     make([]*PetSkill, 0),
		bonds:      NewPetBonds(),
		skins:      make([]*PetSkin, 0),
		createdAt:  now,
		updatedAt:  now,
		version:    1,
	}
}

// GetID 获取宠物ID
func (p *PetAggregate) GetID() string {
	return p.id
}

// GetPlayerID 获取玩家ID
func (p *PetAggregate) GetPlayerID() string {
	return p.playerID
}

// GetConfigID 获取配置ID
func (p *PetAggregate) GetConfigID() uint32 {
	return p.configID
}

// GetName 获取宠物名称
func (p *PetAggregate) GetName() string {
	return p.name
}

// SetName 设置宠物名称
func (p *PetAggregate) SetName(name string) error {
	if name == "" {
		return ErrInvalidPetName
	}
	
	p.name = name
	p.updatedAt = time.Now()
	p.version++
	return nil
}

// GetCategory 获取宠物类别
func (p *PetAggregate) GetCategory() PetCategory {
	return p.category
}

// GetStar 获取宠物星级
func (p *PetAggregate) GetStar() uint32 {
	return p.star
}

// GetLevel 获取宠物等级
func (p *PetAggregate) GetLevel() uint32 {
	return p.level
}

// GetExperience 获取宠物经验
func (p *PetAggregate) GetExperience() uint64 {
	return p.experience
}

// GetState 获取宠物状态
func (p *PetAggregate) GetState() PetState {
	return p.state
}

// GetAttributes 获取宠物属性
func (p *PetAggregate) GetAttributes() *PetAttributes {
	return p.attributes
}

// GetSkills 获取宠物技能
func (p *PetAggregate) GetSkills() []*PetSkill {
	return p.skills
}

// GetBonds 获取宠物羁绊
func (p *PetAggregate) GetBonds() *PetBonds {
	return p.bonds
}

// GetSkins 获取宠物皮肤
func (p *PetAggregate) GetSkins() []*PetSkin {
	return p.skins
}

// GetReviveTime 获取复活时间
func (p *PetAggregate) GetReviveTime() time.Time {
	return p.reviveTime
}

// GetCreatedAt 获取创建时间
func (p *PetAggregate) GetCreatedAt() time.Time {
	return p.createdAt
}

// GetUpdatedAt 获取更新时间
func (p *PetAggregate) GetUpdatedAt() time.Time {
	return p.updatedAt
}

// GetVersion 获取版本号
func (p *PetAggregate) GetVersion() int {
	return p.version
}

// AddExperience 增加经验
func (p *PetAggregate) AddExperience(exp uint64) error {
	if p.state == PetStateDead {
		return ErrPetIsDead
	}
	
	p.experience += exp
	p.updatedAt = time.Now()
	p.version++
	
	// 检查是否可以升级
	if p.canLevelUp() {
		return p.levelUp()
	}
	
	return nil
}

// LevelUp 升级
func (p *PetAggregate) levelUp() error {
	if p.level >= MaxPetLevel {
		return ErrMaxLevelReached
	}
	
	p.level++
	p.attributes.UpgradeOnLevelUp(p.level)
	p.updatedAt = time.Now()
	p.version++
	
	return nil
}

// canLevelUp 检查是否可以升级
func (p *PetAggregate) canLevelUp() bool {
	requiredExp := CalculateRequiredExperience(p.level)
	return p.experience >= requiredExp && p.level < MaxPetLevel
}

// UpgradeStar 升星
func (p *PetAggregate) UpgradeStar() error {
	if p.star >= MaxPetStar {
		return ErrMaxStarReached
	}
	
	p.star++
	p.attributes.UpgradeOnStarUp(p.star)
	p.updatedAt = time.Now()
	p.version++
	
	return nil
}

// ChangeState 改变状态
func (p *PetAggregate) ChangeState(newState PetState) error {
	if !p.canChangeState(newState) {
		return ErrInvalidStateTransition
	}
	
	p.state = newState
	p.updatedAt = time.Now()
	p.version++
	
	// 如果是死亡状态，设置复活时间
	if newState == PetStateDead {
		p.reviveTime = time.Now().Add(DefaultReviveTime)
	}
	
	return nil
}

// canChangeState 检查是否可以改变状态
func (p *PetAggregate) canChangeState(newState PetState) bool {
	switch p.state {
	case PetStateIdle:
		return newState == PetStateBattle || newState == PetStateTraining || newState == PetStateDead
	case PetStateBattle:
		return newState == PetStateIdle || newState == PetStateDead
	case PetStateTraining:
		return newState == PetStateIdle
	case PetStateDead:
		return newState == PetStateIdle && time.Now().After(p.reviveTime)
	default:
		return false
	}
}

// Revive 复活宠物
func (p *PetAggregate) Revive() error {
	if p.state != PetStateDead {
		return ErrPetNotDead
	}
	
	if time.Now().Before(p.reviveTime) {
		return ErrReviveTimeNotReached
	}
	
	p.state = PetStateIdle
	p.reviveTime = time.Time{}
	p.updatedAt = time.Now()
	p.version++
	
	return nil
}

// InstantRevive 立即复活（消耗道具）
func (p *PetAggregate) InstantRevive() error {
	if p.state != PetStateDead {
		return ErrPetNotDead
	}
	
	p.state = PetStateIdle
	p.reviveTime = time.Time{}
	p.updatedAt = time.Now()
	p.version++
	
	return nil
}

// AddSkill 添加技能
func (p *PetAggregate) AddSkill(skill *PetSkill) error {
	if len(p.skills) >= MaxPetSkills {
		return ErrMaxSkillsReached
	}
	
	// 检查是否已存在相同技能
	for _, existingSkill := range p.skills {
		if existingSkill.GetSkillID() == skill.GetSkillID() {
			return ErrSkillAlreadyExists
		}
	}
	
	p.skills = append(p.skills, skill)
	p.updatedAt = time.Now()
	p.version++
	
	return nil
}

// RemoveSkill 移除技能
func (p *PetAggregate) RemoveSkill(skillID string) error {
	for i, skill := range p.skills {
		if skill.GetSkillID() == skillID {
			p.skills = append(p.skills[:i], p.skills[i+1:]...)
			p.updatedAt = time.Now()
			p.version++
			return nil
		}
	}
	
	return ErrSkillNotFound
}

// UpgradeSkill 升级技能
func (p *PetAggregate) UpgradeSkill(skillID string) error {
	for _, skill := range p.skills {
		if skill.GetSkillID() == skillID {
			if err := skill.Upgrade(); err != nil {
				return err
			}
			p.updatedAt = time.Now()
			p.version++
			return nil
		}
	}
	
	return ErrSkillNotFound
}

// AddSkin 添加皮肤
func (p *PetAggregate) AddSkin(skin *PetSkin) error {
	// 检查是否已拥有该皮肤
	for _, existingSkin := range p.skins {
		if existingSkin.GetSkinID() == skin.GetSkinID() {
			return ErrSkinAlreadyOwned
		}
	}
	
	p.skins = append(p.skins, skin)
	p.updatedAt = time.Now()
	p.version++
	
	return nil
}

// EquipSkin 装备皮肤
func (p *PetAggregate) EquipSkin(skinID string) error {
	// 先取消当前装备的皮肤
	for _, skin := range p.skins {
		if skin.IsEquipped() {
			skin.Unequip()
		}
	}
	
	// 装备新皮肤
	for _, skin := range p.skins {
		if skin.GetSkinID() == skinID {
			if err := skin.Equip(); err != nil {
				return err
			}
			p.updatedAt = time.Now()
			p.version++
			return nil
		}
	}
	
	return ErrSkinNotOwned
}

// Feed 喂食
func (p *PetAggregate) Feed(foodType FoodType, amount int) error {
	if p.state == PetStateDead {
		return ErrPetIsDead
	}
	
	if amount <= 0 {
		return ErrInvalidAmount
	}
	
	// 根据食物类型增加不同属性
	switch foodType {
	case FoodTypeExperience:
		return p.AddExperience(uint64(amount * ExperienceFoodValue))
	case FoodTypeHealth:
		p.attributes.AddHealth(int64(amount * HealthFoodValue))
	case FoodTypeAttack:
		p.attributes.AddAttack(int64(amount * AttackFoodValue))
	case FoodTypeDefense:
		p.attributes.AddDefense(int64(amount * DefenseFoodValue))
	default:
		return ErrInvalidFoodType
	}
	
	p.updatedAt = time.Now()
	p.version++
	return nil
}

// Train 训练
func (p *PetAggregate) Train(trainingType TrainingType, duration time.Duration) error {
	if p.state != PetStateIdle {
		return ErrPetNotIdle
	}
	
	p.state = PetStateTraining
	p.updatedAt = time.Now()
	p.version++
	
	// 训练完成后的效果将在训练结束时处理
	return nil
}

// FinishTraining 完成训练
func (p *PetAggregate) FinishTraining(trainingType TrainingType) error {
	if p.state != PetStateTraining {
		return ErrPetNotTraining
	}
	
	// 根据训练类型获得不同收益
	switch trainingType {
	case TrainingTypeExperience:
		p.AddExperience(TrainingExperienceGain)
	case TrainingTypeAttribute:
		p.attributes.AddRandomAttribute(TrainingAttributeGain)
	case TrainingTypeSkill:
		// 随机提升一个技能经验
		if len(p.skills) > 0 {
			randomSkill := p.skills[0] // 简化实现，实际应该随机选择
			randomSkill.AddExperience(TrainingSkillExpGain)
		}
	}
	
	p.state = PetStateIdle
	p.updatedAt = time.Now()
	p.version++
	
	return nil
}

// EnterBattle 进入战斗
func (p *PetAggregate) EnterBattle() error {
	if p.state != PetStateIdle {
		return ErrPetNotIdle
	}
	
	if p.attributes.GetHealth() <= 0 {
		return ErrPetIsDead
	}
	
	p.state = PetStateBattle
	p.updatedAt = time.Now()
	p.version++
	
	return nil
}

// ExitBattle 退出战斗
func (p *PetAggregate) ExitBattle(isDead bool) error {
	if p.state != PetStateBattle {
		return ErrPetNotInBattle
	}
	
	if isDead {
		p.state = PetStateDead
		p.reviveTime = time.Now().Add(DefaultReviveTime)
	} else {
		p.state = PetStateIdle
	}
	
	p.updatedAt = time.Now()
	p.version++
	
	return nil
}

// ActivateBond 激活羁绊
func (p *PetAggregate) ActivateBond(bondID string) error {
	return p.bonds.ActivateBond(bondID)
}

// DeactivateBond 取消羁绊
func (p *PetAggregate) DeactivateBond(bondID string) error {
	return p.bonds.DeactivateBond(bondID)
}

// GetTotalPower 获取总战力
func (p *PetAggregate) GetTotalPower() int64 {
	basePower := p.attributes.CalculatePower()
	bondBonus := p.bonds.GetPowerBonus()
	skinBonus := p.getEquippedSkinBonus()
	
	return basePower + bondBonus + skinBonus
}

// getEquippedSkinBonus 获取装备皮肤加成
func (p *PetAggregate) getEquippedSkinBonus() int64 {
	for _, skin := range p.skins {
		if skin.IsEquipped() {
			return skin.GetPowerBonus()
		}
	}
	return 0
}

// IsAlive 是否存活
func (p *PetAggregate) IsAlive() bool {
	return p.state != PetStateDead
}

// IsIdle 是否空闲
func (p *PetAggregate) IsIdle() bool {
	return p.state == PetStateIdle
}

// CanRevive 是否可以复活
func (p *PetAggregate) CanRevive() bool {
	return p.state == PetStateDead && time.Now().After(p.reviveTime)
}

// GetReviveTimeRemaining 获取剩余复活时间
func (p *PetAggregate) GetReviveTimeRemaining() time.Duration {
	if p.state != PetStateDead {
		return 0
	}
	
	remaining := p.reviveTime.Sub(time.Now())
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Validate 验证宠物数据
func (p *PetAggregate) Validate() error {
	if p.id == "" {
		return ErrInvalidPetID
	}
	
	if p.playerID == "" {
		return ErrInvalidPlayerID
	}
	
	if p.name == "" {
		return ErrInvalidPetName
	}
	
	if p.level < 1 || p.level > MaxPetLevel {
		return ErrInvalidPetLevel
	}
	
	if p.star < 1 || p.star > MaxPetStar {
		return ErrInvalidPetStar
	}
	
	if p.attributes == nil {
		return ErrInvalidPetAttributes
	}
	
	return nil
}

// 常量定义
const (
	MaxPetLevel  = 100
	MaxPetStar   = 5
	MaxPetSkills = 4
	
	DefaultReviveTime = 30 * time.Minute
	
	// 食物价值
	ExperienceFoodValue = 100
	HealthFoodValue     = 50
	AttackFoodValue     = 10
	DefenseFoodValue    = 10
	
	// 训练收益
	TrainingExperienceGain = 500
	TrainingAttributeGain  = 20
	TrainingSkillExpGain   = 100
)

// CalculateRequiredExperience 计算升级所需经验
func CalculateRequiredExperience(level uint32) uint64 {
	// 简化的经验计算公式
	return uint64(level * level * 100)
}

// ReconstructPetAggregate 从持久化数据重建宠物聚合根
func ReconstructPetAggregate(
	id string,
	playerID string,
	configID uint32,
	name string,
	category PetCategory,
	star uint32,
	level uint32,
	experience uint64,
	state PetState,
	attributes *PetAttributes,
	skills []*PetSkill,
	bonds *PetBonds,
	skins []*PetSkin,
	reviveTime time.Time,
	createdAt time.Time,
	updatedAt time.Time,
	version int,
) *PetAggregate {
	return &PetAggregate{
		id:         id,
		playerID:   playerID,
		configID:   configID,
		name:       name,
		category:   category,
		star:       star,
		level:      level,
		experience: experience,
		state:      state,
		attributes: attributes,
		skills:     skills,
		bonds:      bonds,
		skins:      skins,
		reviveTime: reviveTime,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
		version:    version,
	}
}