package character

import (
	"context"
	"sync"
)

// AttributeManager 属性管理器 - 管理Actor的属性计算
type AttributeManager struct {
	owner *Actor
	mu    sync.RWMutex

	base  *Attributes // 基础属性
	final *Attributes // 最终属性（经过装备、Buff等加成）
}

// Attributes 属性集合
type Attributes struct {
	MaxHP   float32 // 最大生命值
	MaxMP   float32 // 最大魔法值
	HPRegen float32 // 生命回复
	MPRegen float32 // 魔法回复

	AD   float32 // 物理攻击力
	AP   float32 // 法术攻击力
	Def  float32 // 物理防御
	MDef float32 // 法术防御

	Cri       float32 // 暴击率
	Crd       float32 // 暴击伤害
	HitRate   float32 // 命中率
	DodgeRate float32 // 闪避率

	Speed       float32 // 移动速度
	AttackSpeed float32 // 攻击速度
}

// AttributeModifier 属性修饰器（用于 Buff/装备 等对属性的加成）
// 约定：Add 为加法；Mul 为乘法（叠加时相加后一次性乘以 1+总和）
type AttributeModifier struct {
	MaxHPAdd, MaxHPMul     float32
	MaxMPAdd, MaxMPMul     float32
	HPRegenAdd, MPRegenAdd float32

	ADAdd, ADMul     float32
	APAdd, APMul     float32
	DefAdd, DefMul   float32
	MDefAdd, MDefMul float32

	CriAdd, CrdAdd                 float32
	HitRateAdd                     float32
	DodgeRateAdd                   float32
	SpeedAdd, SpeedMul             float32
	AttackSpeedAdd, AttackSpeedMul float32
}

// NewAttributeManager 创建属性管理器
func NewAttributeManager(owner *Actor) *AttributeManager {
	return &AttributeManager{
		owner: owner,
		base:  &Attributes{},
		final: &Attributes{},
	}
}

// Start 初始化
func (am *AttributeManager) Start(ctx context.Context) error {
	// TODO: 从配置中加载基础属性（占位：按等级给出默认值，待 DataManager 接入后替换）
	am.mu.Lock()
	if am.base.MaxHP == 0 {
		lvl := float32(am.owner.Level())
		am.base.MaxHP = 100 + 10*lvl
		am.base.MaxMP = 50 + 5*lvl
		am.base.HPRegen = 1
		am.base.MPRegen = 0.5
		am.base.AD = 10 + 2*lvl
		am.base.AP = 5 + 1*lvl
		am.base.Def = 2 + 1*lvl
		am.base.MDef = 1 + 0.5*lvl
		am.base.Cri = 0.05
		am.base.Crd = 1.5
		am.base.HitRate = 0.9
		am.base.DodgeRate = 0.05
		am.base.Speed = 5
		am.base.AttackSpeed = 1
	}
	am.mu.Unlock()

	am.Recalculate()
	return nil
}

// Base 获取基础属性
func (am *AttributeManager) Base() *Attributes {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return am.base
}

// Final 获取最终属性
func (am *AttributeManager) Final() *Attributes {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return am.final
}

// Recalculate 重新计算最终属性
func (am *AttributeManager) Recalculate() {
	am.mu.Lock()
	defer am.mu.Unlock()

	// 基础拷贝
	*am.final = *am.base

	// TODO: 装备加成（占位）
	// 例如：am.final.AD += equip.ADAdd; am.final.AD *= (1 + equip.ADMul)

	// Buff 加成（叠加所有 Buff 的属性修饰器）
	if am.owner != nil && am.owner.buffManager != nil {
		mods := am.owner.buffManager.CollectModifiers()
		var m AttributeModifier // 汇总
		for _, mod := range mods {
			m.MaxHPAdd += mod.MaxHPAdd
			m.MaxHPMul += mod.MaxHPMul
			m.MaxMPAdd += mod.MaxMPAdd
			m.MaxMPMul += mod.MaxMPMul
			m.HPRegenAdd += mod.HPRegenAdd
			m.MPRegenAdd += mod.MPRegenAdd

			m.ADAdd += mod.ADAdd
			m.ADMul += mod.ADMul
			m.APAdd += mod.APAdd
			m.APMul += mod.APMul
			m.DefAdd += mod.DefAdd
			m.DefMul += mod.DefMul
			m.MDefAdd += mod.MDefAdd
			m.MDefMul += mod.MDefMul

			m.CriAdd += mod.CriAdd
			m.CrdAdd += mod.CrdAdd
			m.HitRateAdd += mod.HitRateAdd
			m.DodgeRateAdd += mod.DodgeRateAdd
			m.SpeedAdd += mod.SpeedAdd
			m.SpeedMul += mod.SpeedMul
			m.AttackSpeedAdd += mod.AttackSpeedAdd
			m.AttackSpeedMul += mod.AttackSpeedMul
		}

		// 应用到最终属性
		am.final.MaxHP = (am.final.MaxHP + m.MaxHPAdd) * (1 + m.MaxHPMul)
		am.final.MaxMP = (am.final.MaxMP + m.MaxMPAdd) * (1 + m.MaxMPMul)
		am.final.HPRegen = am.final.HPRegen + m.HPRegenAdd
		am.final.MPRegen = am.final.MPRegen + m.MPRegenAdd

		am.final.AD = (am.final.AD + m.ADAdd) * (1 + m.ADMul)
		am.final.AP = (am.final.AP + m.APAdd) * (1 + m.APMul)
		am.final.Def = (am.final.Def + m.DefAdd) * (1 + m.DefMul)
		am.final.MDef = (am.final.MDef + m.MDefAdd) * (1 + m.MDefMul)

		am.final.Cri += m.CriAdd
		am.final.Crd += m.CrdAdd
		am.final.HitRate += m.HitRateAdd
		am.final.DodgeRate += m.DodgeRateAdd
		am.final.Speed = (am.final.Speed + m.SpeedAdd) * (1 + m.SpeedMul)
		am.final.AttackSpeed = (am.final.AttackSpeed + m.AttackSpeedAdd) * (1 + m.AttackSpeedMul)
	}

	// 下游：可在此应用被动/天赋等
}

// SetBase 设置基础属性（整体替换）
func (am *AttributeManager) SetBase(attrs Attributes) {
	am.mu.Lock()
	am.base = &attrs
	am.mu.Unlock()
	am.Recalculate()
}

// ModifyBase 对基础属性进行增量修改（加法）
func (am *AttributeManager) ModifyBase(mod func(a *Attributes)) {
	am.mu.Lock()
	mod(am.base)
	am.mu.Unlock()
	am.Recalculate()
}

// ========== SkillManager 技能管理器 ==========

// SkillManager 技能管理器
type SkillManager struct {
	owner *Actor
	mu    sync.RWMutex

	skills map[int32]*Skill // 技能ID -> 技能实例
}

// NewSkillManager 创建技能管理器
func NewSkillManager(owner *Actor) *SkillManager {
	return &SkillManager{
		owner:  owner,
		skills: make(map[int32]*Skill),
	}
}

// Start 初始化
func (sm *SkillManager) Start(ctx context.Context) error {
	// TODO: 从配置中加载技能
	return nil
}

// Update 每帧更新
func (sm *SkillManager) Update(ctx context.Context, deltaTime float32) error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// 更新所有技能
	for _, skill := range sm.skills {
		if err := skill.Update(ctx, deltaTime); err != nil {
			return err
		}
	}
	return nil
}

// GetSkill 获取技能
func (sm *SkillManager) GetSkill(skillID int32) *Skill {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.skills[skillID]
}

// AddSkill 添加技能
func (sm *SkillManager) AddSkill(skill *Skill) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.skills[skill.ID()] = skill
}

// ========== Skill 技能 ==========

// Skill 技能实例
type Skill struct {
	id    int32
	owner *Actor
	mu    sync.RWMutex

	// 技能状态
	state         SkillState
	cooldownTimer float32 // 冷却计时器
	castTimer     float32 // 施法计时器

	// 配置（占位，后续由 DataManager 驱动）
	castTime     float32 // 吟唱时间
	activeTime   float32 // 生效窗口时间（如持续伤害/命中帧窗口）
	cooldownTime float32 // 冷却时长

	// 伤害配置（简化版）
	baseDamage float32
	scaleAD    float32
	scaleAP    float32
	dmgType    DamageType
}

// SkillState 技能状态
type SkillState int32

const (
	SkillStateIdle     SkillState = 0 // 空闲
	SkillStateReady    SkillState = 1 // 就绪
	SkillStateIntonate SkillState = 2 // 吟唱中
	SkillStateActive   SkillState = 3 // 激活中
	SkillStateCooling  SkillState = 4 // 冷却中
)

// NewSkill 创建技能
func NewSkill(id int32, owner *Actor) *Skill {
	return &Skill{
		id:    id,
		owner: owner,
		state: SkillStateReady,
		// 默认占位：瞬发、短冷却
		castTime:     0,
		activeTime:   0.1,
		cooldownTime: 1.0,
	}
}

// ID 获取技能ID
func (s *Skill) ID() int32 {
	return s.id
}

// SetDamage 配置技能伤害参数
func (s *Skill) SetDamage(base, scaleAD, scaleAP float32, dmgType DamageType) {
	s.mu.Lock()
	s.baseDamage = base
	s.scaleAD = scaleAD
	s.scaleAP = scaleAP
	s.dmgType = dmgType
	s.mu.Unlock()
}

// Update 更新技能
func (s *Skill) Update(ctx context.Context, deltaTime float32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch s.state {
	case SkillStateIntonate:
		s.castTimer -= deltaTime
		if s.castTimer <= 0 {
			// 进入激活
			s.state = SkillStateActive
			s.castTimer = s.activeTime
			// 命中/效果应用（命中目标、生成投射物等）
			if s.owner != nil && s.owner.spell != nil {
				s.owner.spell.ApplySkillEffect(s)
			}
		}
	case SkillStateActive:
		s.castTimer -= deltaTime
		if s.castTimer <= 0 {
			// 进入冷却
			s.state = SkillStateCooling
			s.cooldownTimer = s.cooldownTime
		}
	case SkillStateCooling:
		s.cooldownTimer -= deltaTime
		if s.cooldownTimer <= 0 {
			s.state = SkillStateReady
			s.cooldownTimer = 0
		}
	}
	return nil
}

// State 获取当前技能状态
func (s *Skill) State() SkillState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

// SetTimings 设置技能关键时序（吟唱/激活/冷却）
func (s *Skill) SetTimings(cast, active, cooldown float32) {
	s.mu.Lock()
	s.castTime = cast
	s.activeTime = active
	s.cooldownTime = cooldown
	s.mu.Unlock()
}

// StartCast 尝试开始施法（由施法器或应用层触发）
func (s *Skill) StartCast() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.owner == nil || s.owner.IsDeath() {
		return false
	}
	// 眩晕/沉默等状态约束
	flags := s.owner.GetFlagState()
	if flags.HasFlag(FlagStateStun) || flags.HasFlag(FlagStateSilence) {
		return false
	}
	// 冷却中不可释放
	if s.state == SkillStateCooling || s.state == SkillStateActive || s.state == SkillStateIntonate {
		return false
	}

	if s.castTime > 0 {
		s.state = SkillStateIntonate
		s.castTimer = s.castTime
	} else {
		s.state = SkillStateActive
		s.castTimer = s.activeTime
	}
	return true
}

// ========== BuffManager Buff管理器 ==========

// BuffManager Buff管理器
type BuffManager struct {
	owner *Actor
	mu    sync.RWMutex

	buffs []*Buff // Buff列表
}

// NewBuffManager 创建Buff管理器
func NewBuffManager(owner *Actor) *BuffManager {
	return &BuffManager{
		owner: owner,
		buffs: make([]*Buff, 0),
	}
}

// Start 初始化
func (bm *BuffManager) Start(ctx context.Context) error {
	return nil
}

// Update 每帧更新
func (bm *BuffManager) Update(ctx context.Context, deltaTime float32) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	// 更新所有Buff
	toRemove := make([]int, 0)
	for i, buff := range bm.buffs {
		if err := buff.Update(ctx, deltaTime); err != nil {
			return err
		}
		if buff.IsExpired() {
			toRemove = append(toRemove, i)
		}
	}

	// 移除过期的Buff
	for i := len(toRemove) - 1; i >= 0; i-- {
		idx := toRemove[i]
		bm.buffs = append(bm.buffs[:idx], bm.buffs[idx+1:]...)
	}

	if len(toRemove) > 0 {
		// Buff 变化后触发属性重算
		if bm.owner != nil && bm.owner.attributeManager != nil {
			bm.owner.attributeManager.Recalculate()
		}
		// 刷新基于 Buff 的状态标志
		bm.refreshActorFlags()
	}

	return nil
}

// AddBuff 添加Buff
func (bm *BuffManager) AddBuff(buff *Buff) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.buffs = append(bm.buffs, buff)
	// Buff 变化后触发属性重算
	if bm.owner != nil && bm.owner.attributeManager != nil {
		bm.owner.attributeManager.Recalculate()
	}
	// 刷新状态标志
	bm.refreshActorFlags()
}

// RemoveBuff 移除Buff
func (bm *BuffManager) RemoveBuff(buff *Buff) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	for i, b := range bm.buffs {
		if b == buff {
			bm.buffs = append(bm.buffs[:i], bm.buffs[i+1:]...)
			break
		}
	}
	// Buff 变化后触发属性重算
	if bm.owner != nil && bm.owner.attributeManager != nil {
		bm.owner.attributeManager.Recalculate()
	}
	// 刷新状态标志
	bm.refreshActorFlags()
}

// GetBuffByID 按ID获取Buff
func (bm *BuffManager) GetBuffByID(id int32) *Buff {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	for _, b := range bm.buffs {
		if b != nil && b.id == id {
			return b
		}
	}
	return nil
}

// RemoveBuffByID 按ID移除Buff
func (bm *BuffManager) RemoveBuffByID(id int32) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	for i := 0; i < len(bm.buffs); i++ {
		if bm.buffs[i] != nil && bm.buffs[i].id == id {
			bm.buffs = append(bm.buffs[:i], bm.buffs[i+1:]...)
			i--
		}
	}
	if bm.owner != nil && bm.owner.attributeManager != nil {
		bm.owner.attributeManager.Recalculate()
	}
	bm.refreshActorFlags()
}

// CollectModifiers 汇总当前 Buff 的属性修饰器快照
func (bm *BuffManager) CollectModifiers() []AttributeModifier {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	mods := make([]AttributeModifier, 0, len(bm.buffs))
	for _, b := range bm.buffs {
		mods = append(mods, b.Modifier())
	}
	return mods
}

// 汇总 Buff 的状态标志位（位或）
func (bm *BuffManager) collectFlags() FlagState {
	var flags FlagState = FlagStateZero
	for _, b := range bm.buffs {
		flags = flags.AddFlag(b.FlagAdd())
	}
	return flags
}

// 刷新 Actor 的状态标志，基于当前 Buff 汇总
func (bm *BuffManager) refreshActorFlags() {
	if bm.owner == nil {
		return
	}
	flags := bm.collectFlags()
	bm.owner.SetFlagStateExact(flags)
}

// ========== Buff ==========

// Buff Buff实例
type Buff struct {
	id       int32
	owner    *Actor
	caster   *Actor
	duration float32
	elapsed  float32

	modifier AttributeModifier

	// 状态效果：为简化，使用位或累加的 FlagState
	addFlags FlagState
}

// NewBuff 创建Buff
func NewBuff(id int32, owner, caster *Actor, duration float32) *Buff {
	return &Buff{
		id:       id,
		owner:    owner,
		caster:   caster,
		duration: duration,
		elapsed:  0,
	}
}

// Update 更新Buff
func (b *Buff) Update(ctx context.Context, deltaTime float32) error {
	b.elapsed += deltaTime
	return nil
}

// IsExpired 是否过期
func (b *Buff) IsExpired() bool {
	return b.elapsed >= b.duration
}

// SetModifier 设置属性修饰器
func (b *Buff) SetModifier(mod AttributeModifier) { b.modifier = mod }

// Modifier 获取属性修饰器
func (b *Buff) Modifier() AttributeModifier { return b.modifier }

// SetFlagAdd 设置该 Buff 施加的状态标志
func (b *Buff) SetFlagAdd(flags FlagState) { b.addFlags = flags }

// FlagAdd 获取该 Buff 施加的状态标志
func (b *Buff) FlagAdd() FlagState { return b.addFlags }

// ========== Spell 施法器 ==========

// Spell 施法器 - 管理当前正在施放的技能
type Spell struct {
	owner *Actor
	mu    sync.RWMutex

	currentSkill *Skill // 当前正在施放的技能
	target       *Actor // 当前施法目标（简化：单体）
}

// NewSpell 创建施法器
func NewSpell(owner *Actor) *Spell {
	return &Spell{
		owner: owner,
	}
}

// CurrentSkill 获取当前技能
func (s *Spell) CurrentSkill() *Skill {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentSkill
}

// SetCurrentSkill 设置当前技能
func (s *Spell) SetCurrentSkill(skill *Skill) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentSkill = skill
}

// Cast 根据技能ID对目标施放技能
func (s *Spell) Cast(skillID int32, target *Actor) bool {
	if s.owner == nil {
		return false
	}
	sm := s.owner.GetSkillManager()
	if sm == nil {
		return false
	}
	sk := sm.GetSkill(skillID)
	if sk == nil {
		return false
	}
	if !sk.StartCast() {
		return false
	}
	s.mu.Lock()
	s.currentSkill = sk
	s.target = target
	s.mu.Unlock()
	return true
}

// ApplySkillEffect 在技能进入 Active 时调用，应用技能效果到目标
func (s *Spell) ApplySkillEffect(skill *Skill) {
	s.mu.RLock()
	target := s.target
	s.mu.RUnlock()
	if target == nil || target.IsDeath() || s.owner == nil {
		return
	}
	// 计算伤害并应用
	dmg := computeDamage(s.owner, target, skill)
	if dmg <= 0 {
		return
	}
	info := &DamageInfo{
		TargetID:     target.ID(),
		AttackerInfo: AttackerInfo{AttackerID: s.owner.ID(), AttackerType: AttackerTypeSkill, SkillID: skill.ID()},
		Amount:       int32(dmg),
		DamageType:   skill.dmgType,
		IsCrit:       false, IsMiss: false,
	}
	_ = target.OnHurt(context.Background(), info)
}

// computeDamage 伤害计算（简化且确定性）
func computeDamage(attacker *Actor, defender *Actor, skill *Skill) float32 {
	if attacker == nil || defender == nil || skill == nil {
		return 0
	}
	af := attacker.GetAttributeManager().Final()
	df := defender.GetAttributeManager().Final()
	// 攻击力构成
	atk := skill.baseDamage + af.AD*skill.scaleAD + af.AP*skill.scaleAP
	if atk <= 0 {
		return 0
	}
	// 防御减伤： dmg * (1 - def/(def+100))
	var def float32
	if skill.dmgType == DamageTypeMagical {
		def = df.MDef
	} else {
		def = df.Def
	}
	reduction := float32(1.0)
	if def > 0 {
		reduction = 1 - def/(def+100)
	}
	dmg := atk * reduction
	if dmg < 0 {
		dmg = 0
	}
	return dmg
}

// SetTarget 设置当前施法目标
func (s *Spell) SetTarget(target *Actor) {
	s.mu.Lock()
	s.target = target
	s.mu.Unlock()
}

// Target 获取当前施法目标
func (s *Spell) Target() *Actor {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.target
}
