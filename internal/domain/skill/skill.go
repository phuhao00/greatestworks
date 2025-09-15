package skill

import (
	"errors"
	"time"
)

// SkillTree 技能树聚合根
type SkillTree struct {
	playerID     string
	skills       map[string]*Skill
	skillPoints  int64
	totalPoints  int64
	lastUpdate   time.Time
	events       []DomainEvent
}

// NewSkillTree 创建新技能树
func NewSkillTree(playerID string) *SkillTree {
	return &SkillTree{
		playerID:    playerID,
		skills:      make(map[string]*Skill),
		skillPoints: 0,
		totalPoints: 0,
		lastUpdate:  time.Now(),
		events:      make([]DomainEvent, 0),
	}
}

// Skill 技能实体
type Skill struct {
	id            string
	name          string
	description   string
	skillType     SkillType
	level         int
	maxLevel      int
	prerequisites []string // 前置技能ID
	effects       []*SkillEffect
	cooldown      time.Duration
	lastUsed      *time.Time
	manaCost      int64
	castTime      time.Duration
	range_        float64
	damageType    DamageType
	baseDamage    int64
	scaling       map[AttributeType]float64
	createdAt     time.Time
	updatedAt     time.Time
}

// NewSkill 创建新技能
func NewSkill(id, name string, skillType SkillType) *Skill {
	return &Skill{
		id:          id,
		name:        name,
		skillType:   skillType,
		level:       0,
		maxLevel:    10,
		effects:     make([]*SkillEffect, 0),
		scaling:     make(map[AttributeType]float64),
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}
}

// SkillType 技能类型
type SkillType int

const (
	SkillTypeActive SkillType = iota + 1
	SkillTypePassive
	SkillTypeToggle
	SkillTypeChanneled
	SkillTypeInstant
)

// DamageType 伤害类型
type DamageType int

const (
	DamageTypePhysical DamageType = iota + 1
	DamageTypeMagical
	DamageTypeTrue
	DamageTypeHealing
)

// AttributeType 属性类型
type AttributeType int

const (
	AttributeTypeStrength AttributeType = iota + 1
	AttributeTypeIntelligence
	AttributeTypeAgility
	AttributeTypeVitality
	AttributeTypeSpirit
)

// SkillEffect 技能效果
type SkillEffect struct {
	effectType EffectType
	value      float64
	duration   time.Duration
	target     TargetType
	condition  *EffectCondition
}

// EffectType 效果类型
type EffectType int

const (
	EffectTypeDamage EffectType = iota + 1
	EffectTypeHeal
	EffectTypeBuff
	EffectTypeDebuff
	EffectTypeStun
	EffectTypeSilence
	EffectTypeRoot
	EffectTypeSlow
	EffectTypeHaste
	EffectTypeShield
	EffectTypeReflect
)

// TargetType 目标类型
type TargetType int

const (
	TargetTypeSelf TargetType = iota + 1
	TargetTypeEnemy
	TargetTypeAlly
	TargetTypeAll
	TargetTypeArea
)

// EffectCondition 效果条件
type EffectCondition struct {
	conditionType ConditionType
	value         interface{}
}

// ConditionType 条件类型
type ConditionType int

const (
	ConditionTypeHealthBelow ConditionType = iota + 1
	ConditionTypeManaBelow
	ConditionTypeEnemyCount
	ConditionTypeBuffActive
	ConditionTypeDebuffActive
)

// SkillCombo 技能连击
type SkillCombo struct {
	id         string
	name       string
	skills     []string // 技能ID序列
	timeWindow time.Duration
	bonusEffect *SkillEffect
}

// DomainEvent 领域事件接口
type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
	PlayerID() string
}

// SkillLearnedEvent 技能学习事件
type SkillLearnedEvent struct {
	playerID   string
	skillID    string
	occurredAt time.Time
}

func (e SkillLearnedEvent) EventType() string   { return "skill.learned" }
func (e SkillLearnedEvent) OccurredAt() time.Time { return e.occurredAt }
func (e SkillLearnedEvent) PlayerID() string    { return e.playerID }

// SkillUpgradedEvent 技能升级事件
type SkillUpgradedEvent struct {
	playerID   string
	skillID    string
	oldLevel   int
	newLevel   int
	occurredAt time.Time
}

func (e SkillUpgradedEvent) EventType() string   { return "skill.upgraded" }
func (e SkillUpgradedEvent) OccurredAt() time.Time { return e.occurredAt }
func (e SkillUpgradedEvent) PlayerID() string    { return e.playerID }

// SkillUsedEvent 技能使用事件
type SkillUsedEvent struct {
	playerID   string
	skillID    string
	targetID   string
	damage     int64
	occurredAt time.Time
}

func (e SkillUsedEvent) EventType() string   { return "skill.used" }
func (e SkillUsedEvent) OccurredAt() time.Time { return e.occurredAt }
func (e SkillUsedEvent) PlayerID() string    { return e.playerID }

// SkillPointsGainedEvent 技能点获得事件
type SkillPointsGainedEvent struct {
	playerID   string
	points     int64
	reason     string
	occurredAt time.Time
}

func (e SkillPointsGainedEvent) EventType() string   { return "skill.points.gained" }
func (e SkillPointsGainedEvent) OccurredAt() time.Time { return e.occurredAt }
func (e SkillPointsGainedEvent) PlayerID() string    { return e.playerID }

// SkillTree 业务方法

// PlayerID 获取玩家ID
func (st *SkillTree) PlayerID() string {
	return st.playerID
}

// SkillPoints 获取技能点
func (st *SkillTree) SkillPoints() int64 {
	return st.skillPoints
}

// TotalPoints 获取总技能点
func (st *SkillTree) TotalPoints() int64 {
	return st.totalPoints
}

// Skills 获取所有技能
func (st *SkillTree) Skills() map[string]*Skill {
	return st.skills
}

// GetSkill 获取指定技能
func (st *SkillTree) GetSkill(skillID string) (*Skill, bool) {
	skill, exists := st.skills[skillID]
	return skill, exists
}

// LearnSkill 学习技能
func (st *SkillTree) LearnSkill(skillID string, skillData *Skill) error {
	if st.skillPoints <= 0 {
		return ErrInsufficientSkillPoints
	}

	// 检查是否已学习
	if _, exists := st.skills[skillID]; exists {
		return ErrSkillAlreadyLearned
	}

	// 检查前置技能
	if !st.checkPrerequisites(skillData.prerequisites) {
		return ErrPrerequisitesNotMet
	}

	// 学习技能
	skillData.level = 1
	skillData.updatedAt = time.Now()
	st.skills[skillID] = skillData
	st.skillPoints--
	st.lastUpdate = time.Now()

	// 发布事件
	st.addEvent(SkillLearnedEvent{
		playerID:   st.playerID,
		skillID:    skillID,
		occurredAt: time.Now(),
	})

	return nil
}

// UpgradeSkill 升级技能
func (st *SkillTree) UpgradeSkill(skillID string) error {
	skill, exists := st.skills[skillID]
	if !exists {
		return ErrSkillNotLearned
	}

	if skill.level >= skill.maxLevel {
		return ErrSkillMaxLevel
	}

	requiredPoints := st.calculateUpgradeCost(skill.level)
	if st.skillPoints < requiredPoints {
		return ErrInsufficientSkillPoints
	}

	oldLevel := skill.level
	skill.level++
	skill.updatedAt = time.Now()
	st.skillPoints -= requiredPoints
	st.lastUpdate = time.Now()

	// 发布事件
	st.addEvent(SkillUpgradedEvent{
		playerID:   st.playerID,
		skillID:    skillID,
		oldLevel:   oldLevel,
		newLevel:   skill.level,
		occurredAt: time.Now(),
	})

	return nil
}

// UseSkill 使用技能
func (st *SkillTree) UseSkill(skillID string, targetID string) (*SkillResult, error) {
	skill, exists := st.skills[skillID]
	if !exists {
		return nil, ErrSkillNotLearned
	}

	if skill.skillType == SkillTypePassive {
		return nil, ErrPassiveSkillNotUsable
	}

	// 检查冷却时间
	if skill.lastUsed != nil && time.Since(*skill.lastUsed) < skill.cooldown {
		return nil, ErrSkillOnCooldown
	}

	// 计算伤害和效果
	result := st.calculateSkillResult(skill, targetID)

	// 更新使用时间
	now := time.Now()
	skill.lastUsed = &now
	st.lastUpdate = time.Now()

	// 发布事件
	st.addEvent(SkillUsedEvent{
		playerID:   st.playerID,
		skillID:    skillID,
		targetID:   targetID,
		damage:     result.Damage,
		occurredAt: time.Now(),
	})

	return result, nil
}

// AddSkillPoints 添加技能点
func (st *SkillTree) AddSkillPoints(points int64, reason string) error {
	if points <= 0 {
		return ErrInvalidSkillPoints
	}

	st.skillPoints += points
	st.totalPoints += points
	st.lastUpdate = time.Now()

	// 发布事件
	st.addEvent(SkillPointsGainedEvent{
		playerID:   st.playerID,
		points:     points,
		reason:     reason,
		occurredAt: time.Now(),
	})

	return nil
}

// ResetSkills 重置技能
func (st *SkillTree) ResetSkills() error {
	// 计算返还的技能点
	refundPoints := int64(0)
	for _, skill := range st.skills {
		for level := 1; level <= skill.level; level++ {
			refundPoints += st.calculateUpgradeCost(level - 1)
		}
	}

	// 重置所有技能
	st.skills = make(map[string]*Skill)
	st.skillPoints += refundPoints
	st.lastUpdate = time.Now()

	return nil
}

// checkPrerequisites 检查前置技能
func (st *SkillTree) checkPrerequisites(prerequisites []string) bool {
	for _, prereq := range prerequisites {
		if _, exists := st.skills[prereq]; !exists {
			return false
		}
	}
	return true
}

// calculateUpgradeCost 计算升级消耗
func (st *SkillTree) calculateUpgradeCost(currentLevel int) int64 {
	// 基础消耗 + 等级加成
	return int64(currentLevel + 1)
}

// calculateSkillResult 计算技能结果
func (st *SkillTree) calculateSkillResult(skill *Skill, targetID string) *SkillResult {
	// 基础伤害计算
	damage := skill.baseDamage * int64(skill.level)
	
	// 这里可以添加更复杂的伤害计算逻辑
	// 包括属性加成、暴击、抗性等
	
	return &SkillResult{
		SkillID:  skill.id,
		TargetID: targetID,
		Damage:   damage,
		Effects:  skill.effects,
		Success:  true,
	}
}

// addEvent 添加领域事件
func (st *SkillTree) addEvent(event DomainEvent) {
	st.events = append(st.events, event)
}

// GetEvents 获取领域事件
func (st *SkillTree) GetEvents() []DomainEvent {
	return st.events
}

// ClearEvents 清除领域事件
func (st *SkillTree) ClearEvents() {
	st.events = make([]DomainEvent, 0)
}

// SkillResult 技能使用结果
type SkillResult struct {
	SkillID  string
	TargetID string
	Damage   int64
	Effects  []*SkillEffect
	Success  bool
	Message  string
}