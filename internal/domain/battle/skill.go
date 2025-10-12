package battle

import (
	"errors"
	"time"
)

// SkillID 技能ID值对象
type SkillID struct {
	value string
}

// NewSkillID 创建技能ID
func NewSkillID(id string) SkillID {
	return SkillID{value: id}
}

// String 返回字符串表示
func (id SkillID) String() string {
	return id.value
}

// SkillType 技能类型枚举
type SkillType int

const (
	SkillTypeAttack  SkillType = iota // 攻击技能
	SkillTypeDefense                  // 防御技能
	SkillTypeHeal                     // 治疗技能
	SkillTypeBuff                     // 增益技能
	SkillTypeDebuff                   // 减益技能
	SkillTypeUtility                  // 辅助技能
)

// SkillTarget 技能目标类型
type SkillTarget int

const (
	SkillTargetSelf     SkillTarget = iota // 自己
	SkillTargetEnemy                       // 敌人
	SkillTargetAlly                        // 队友
	SkillTargetAll                         // 所有人
	SkillTargetEnemyAll                    // 所有敌人
	SkillTargetAllyAll                     // 所有队友
)

// Skill 技能实体
type Skill struct {
	id          SkillID
	name        string
	description string
	skillType   SkillType
	targetType  SkillTarget
	manaCost    int
	cooldown    time.Duration
	damage      int
	healing     int
	effects     []*SkillEffect
	range_      float64
	level       int
	maxLevel    int
}

// NewSkill 创建新技能
func NewSkill(id, name, description string, skillType SkillType, targetType SkillTarget) *Skill {
	return &Skill{
		id:          NewSkillID(id),
		name:        name,
		description: description,
		skillType:   skillType,
		targetType:  targetType,
		manaCost:    10,
		cooldown:    time.Second * 3,
		damage:      0,
		healing:     0,
		effects:     make([]*SkillEffect, 0),
		range_:      5.0,
		level:       1,
		maxLevel:    10,
	}
}

// ID 获取技能ID
func (s *Skill) ID() SkillID {
	return s.id
}

// Name 获取技能名称
func (s *Skill) Name() string {
	return s.name
}

// Description 获取技能描述
func (s *Skill) Description() string {
	return s.description
}

// SkillType 获取技能类型
func (s *Skill) GetSkillType() SkillType {
	return s.skillType
}

// TargetType 获取目标类型
func (s *Skill) TargetType() SkillTarget {
	return s.targetType
}

// ManaCost 获取魔法消耗
func (s *Skill) ManaCost() int {
	return s.manaCost
}

// Cooldown 获取冷却时间
func (s *Skill) Cooldown() time.Duration {
	return s.cooldown
}

// Damage 获取伤害值
func (s *Skill) Damage() int {
	return s.damage
}

// Healing 获取治疗值
func (s *Skill) Healing() int {
	return s.healing
}

// Effects 获取技能效果
func (s *Skill) Effects() []*SkillEffect {
	return s.effects
}

// Range 获取技能范围
func (s *Skill) Range() float64 {
	return s.range_
}

// Level 获取技能等级
func (s *Skill) Level() int {
	return s.level
}

// MaxLevel 获取最大等级
func (s *Skill) MaxLevel() int {
	return s.maxLevel
}

// CanUpgrade 是否可以升级
func (s *Skill) CanUpgrade() bool {
	return s.level < s.maxLevel
}

// Upgrade 升级技能
func (s *Skill) Upgrade() error {
	if !s.CanUpgrade() {
		return errors.New("skill already at max level")
	}

	s.level++
	// 升级时提升技能属性
	s.damage = int(float64(s.damage) * 1.1)
	s.healing = int(float64(s.healing) * 1.1)

	return nil
}

// SkillEffect 技能效果
type SkillEffect struct {
	effectType EffectType
	value      int
	duration   time.Duration
	stackable  bool
}

// EffectType 效果类型枚举
type EffectType int

const (
	EffectTypePoison           EffectType = iota // 中毒
	EffectTypeBurn                               // 燃烧
	EffectTypeFreeze                             // 冰冻
	EffectTypeStun                               // 眩晕
	EffectTypeAttackBoost                        // 攻击力提升
	EffectTypeDefenseBoost                       // 防御力提升
	EffectTypeSpeedBoost                         // 速度提升
	EffectTypeAttackReduction                    // 攻击力降低
	EffectTypeDefenseReduction                   // 防御力降低
	EffectTypeSpeedReduction                     // 速度降低
)

// NewSkillEffect 创建技能效果
func NewSkillEffect(effectType EffectType, value int, duration time.Duration, stackable bool) *SkillEffect {
	return &SkillEffect{
		effectType: effectType,
		value:      value,
		duration:   duration,
		stackable:  stackable,
	}
}

// EffectType 获取效果类型
func (e *SkillEffect) GetEffectType() EffectType {
	return e.effectType
}

// Value 获取效果值
func (e *SkillEffect) Value() int {
	return e.value
}

// Duration 获取持续时间
func (e *SkillEffect) Duration() time.Duration {
	return e.duration
}

// IsStackable 是否可叠加
func (e *SkillEffect) IsStackable() bool {
	return e.stackable
}

// PlayerSkill 玩家技能
type PlayerSkill struct {
	skill        *Skill
	lastUsedTime time.Time
	uses         int
}

// NewPlayerSkill 创建玩家技能
func NewPlayerSkill(skill *Skill) *PlayerSkill {
	return &PlayerSkill{
		skill:        skill,
		lastUsedTime: time.Time{},
		uses:         0,
	}
}

// Skill 获取技能
func (ps *PlayerSkill) GetSkill() *Skill {
	return ps.skill
}

// CanUse 是否可以使用
func (ps *PlayerSkill) CanUse(currentMana int) bool {
	// 检查魔法值是否足够
	if currentMana < ps.skill.ManaCost() {
		return false
	}

	// 检查冷却时间
	if time.Since(ps.lastUsedTime) < ps.skill.Cooldown() {
		return false
	}

	return true
}

// Use 使用技能
func (ps *PlayerSkill) Use() error {
	if time.Since(ps.lastUsedTime) < ps.skill.Cooldown() {
		return errors.New("skill is on cooldown")
	}

	ps.lastUsedTime = time.Now()
	ps.uses++

	return nil
}

// GetCooldownRemaining 获取剩余冷却时间
func (ps *PlayerSkill) GetCooldownRemaining() time.Duration {
	elapsed := time.Since(ps.lastUsedTime)
	if elapsed >= ps.skill.Cooldown() {
		return 0
	}
	return ps.skill.Cooldown() - elapsed
}

// Uses 获取使用次数
func (ps *PlayerSkill) Uses() int {
	return ps.uses
}

// SkillRegistry 技能注册表
type SkillRegistry struct {
	skills map[string]*Skill
}

// NewSkillRegistry 创建技能注册表
func NewSkillRegistry() *SkillRegistry {
	return &SkillRegistry{
		skills: make(map[string]*Skill),
	}
}

// RegisterSkill 注册技能
func (sr *SkillRegistry) RegisterSkill(skill *Skill) {
	sr.skills[skill.ID().String()] = skill
}

// GetSkill 获取技能
func (sr *SkillRegistry) GetSkill(skillID string) (*Skill, error) {
	skill, exists := sr.skills[skillID]
	if !exists {
		return nil, errors.New("skill not found")
	}
	return skill, nil
}

// GetAllSkills 获取所有技能
func (sr *SkillRegistry) GetAllSkills() []*Skill {
	skills := make([]*Skill, 0, len(sr.skills))
	for _, skill := range sr.skills {
		skills = append(skills, skill)
	}
	return skills
}

// InitializeDefaultSkills 初始化默认技能
func (sr *SkillRegistry) InitializeDefaultSkills() {
	// 基础攻击技能
	basicAttack := NewSkill("basic_attack", "基础攻击", "普通的物理攻击", SkillTypeAttack, SkillTargetEnemy)
	basicAttack.damage = 20
	basicAttack.manaCost = 0
	basicAttack.cooldown = time.Second * 1
	sr.RegisterSkill(basicAttack)

	// 火球术
	fireball := NewSkill("fireball", "火球术", "发射一个火球攻击敌人", SkillTypeAttack, SkillTargetEnemy)
	fireball.damage = 35
	fireball.manaCost = 15
	fireball.cooldown = time.Second * 3
	fireball.effects = append(fireball.effects, NewSkillEffect(EffectTypeBurn, 5, time.Second*3, false))
	sr.RegisterSkill(fireball)

	// 治疗术
	heal := NewSkill("heal", "治疗术", "恢复目标的生命值", SkillTypeHeal, SkillTargetAlly)
	heal.healing = 30
	heal.manaCost = 20
	heal.cooldown = time.Second * 2
	sr.RegisterSkill(heal)

	// 防御姿态
	defense := NewSkill("defense_stance", "防御姿态", "提高防御力", SkillTypeDefense, SkillTargetSelf)
	defense.manaCost = 10
	defense.cooldown = time.Second * 5
	defense.effects = append(defense.effects, NewSkillEffect(EffectTypeDefenseBoost, 10, time.Second*10, false))
	sr.RegisterSkill(defense)
}
