package sacred

import (
	"fmt"
	"time"
)

// Challenge 挑战实体
type Challenge struct {
	id            string
	name          string
	description   string
	type_         ChallengeType
	difficulty    ChallengeDifficulty
	requiredLevel int
	status        ChallengeStatus
	duration      time.Duration
	cooldown      time.Duration
	lastStartTime time.Time
	lastEndTime   time.Time
	participants  map[string]*ChallengeParticipant
	rewards       *ChallengeReward
	conditions    []*ChallengeCondition
	createdAt     time.Time
	updatedAt     time.Time
}

// NewChallenge 创建挑战
func NewChallenge(id, name, description string, challengeType ChallengeType, difficulty ChallengeDifficulty, requiredLevel int) *Challenge {
	now := time.Now()
	return &Challenge{
		id:            id,
		name:          name,
		description:   description,
		type_:         challengeType,
		difficulty:    difficulty,
		requiredLevel: requiredLevel,
		status:        ChallengeStatusAvailable,
		duration:      time.Hour, // 默认1小时
		cooldown:      time.Hour * 24, // 默认24小时冷却
		participants:  make(map[string]*ChallengeParticipant),
		rewards:       NewChallengeReward(challengeType, difficulty),
		conditions:    make([]*ChallengeCondition, 0),
		createdAt:     now,
		updatedAt:     now,
	}
}

// GetID 获取ID
func (c *Challenge) GetID() string {
	return c.id
}

// GetName 获取名称
func (c *Challenge) GetName() string {
	return c.name
}

// GetDescription 获取描述
func (c *Challenge) GetDescription() string {
	return c.description
}

// GetType 获取类型
func (c *Challenge) GetType() ChallengeType {
	return c.type_
}

// GetDifficulty 获取难度
func (c *Challenge) GetDifficulty() ChallengeDifficulty {
	return c.difficulty
}

// GetRequiredLevel 获取所需等级
func (c *Challenge) GetRequiredLevel() int {
	return c.requiredLevel
}

// GetStatus 获取状态
func (c *Challenge) GetStatus() ChallengeStatus {
	return c.status
}

// GetDuration 获取持续时间
func (c *Challenge) GetDuration() time.Duration {
	return c.duration
}

// GetCooldown 获取冷却时间
func (c *Challenge) GetCooldown() time.Duration {
	return c.cooldown
}

// GetRewards 获取奖励
func (c *Challenge) GetRewards() *ChallengeReward {
	return c.rewards
}

// CanStart 检查是否可以开始
func (c *Challenge) CanStart() bool {
	if c.status != ChallengeStatusAvailable {
		return false
	}
	
	// 检查冷却时间
	if !c.lastEndTime.IsZero() && time.Since(c.lastEndTime) < c.cooldown {
		return false
	}
	
	return true
}

// Start 开始挑战
func (c *Challenge) Start(playerID string) (*ChallengeResult, error) {
	if !c.CanStart() {
		return nil, fmt.Errorf("challenge cannot be started")
	}
	
	// 创建参与者
	participant := NewChallengeParticipant(playerID, c.id)
	c.participants[playerID] = participant
	
	// 更新状态
	c.status = ChallengeStatusInProgress
	c.lastStartTime = time.Now()
	c.updatedAt = time.Now()
	
	// 创建挑战结果
	result := &ChallengeResult{
		ChallengeID: c.id,
		PlayerID:    playerID,
		StartTime:   c.lastStartTime,
		Status:      "started",
	}
	
	return result, nil
}

// Complete 完成挑战
func (c *Challenge) Complete(playerID string, success bool, score int) (*ChallengeReward, error) {
	participant, exists := c.participants[playerID]
	if !exists {
		return nil, fmt.Errorf("participant not found")
	}
	
	// 完成参与者记录
	participant.Complete(success, score)
	
	// 更新挑战状态
	c.status = ChallengeStatusCompleted
	c.lastEndTime = time.Now()
	c.updatedAt = time.Now()
	
	// 计算奖励
	reward := c.calculateReward(success, score)
	
	return reward, nil
}

// calculateReward 计算奖励
func (c *Challenge) calculateReward(success bool, score int) *ChallengeReward {
	if !success {
		return &ChallengeReward{
			Gold:       0,
			Experience: 0,
			Items:      make(map[string]int),
		}
	}
	
	// 基础奖励
	baseReward := c.rewards
	
	// 根据分数调整奖励
	multiplier := float64(score) / 100.0
	if multiplier > 2.0 {
		multiplier = 2.0
	}
	if multiplier < 0.1 {
		multiplier = 0.1
	}
	
	return &ChallengeReward{
		Gold:       int(float64(baseReward.Gold) * multiplier),
		Experience: int(float64(baseReward.Experience) * multiplier),
		Items:      baseReward.Items,
	}
}

// AddCondition 添加条件
func (c *Challenge) AddCondition(condition *ChallengeCondition) {
	c.conditions = append(c.conditions, condition)
	c.updatedAt = time.Now()
}

// CheckConditions 检查条件
func (c *Challenge) CheckConditions(playerData map[string]interface{}) bool {
	for _, condition := range c.conditions {
		if !condition.Check(playerData) {
			return false
		}
	}
	return true
}

// SetDuration 设置持续时间
func (c *Challenge) SetDuration(duration time.Duration) {
	c.duration = duration
	c.updatedAt = time.Now()
}

// SetCooldown 设置冷却时间
func (c *Challenge) SetCooldown(cooldown time.Duration) {
	c.cooldown = cooldown
	c.updatedAt = time.Now()
}

// GetRemainingCooldown 获取剩余冷却时间
func (c *Challenge) GetRemainingCooldown() time.Duration {
	if c.lastEndTime.IsZero() {
		return 0
	}
	
	elapsed := time.Since(c.lastEndTime)
	if elapsed >= c.cooldown {
		return 0
	}
	
	return c.cooldown - elapsed
}

// ToMap 转换为映射
func (c *Challenge) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":             c.id,
		"name":           c.name,
		"description":    c.description,
		"type":           c.type_.String(),
		"difficulty":     c.difficulty.String(),
		"required_level": c.requiredLevel,
		"status":         c.status.String(),
		"duration":       c.duration.String(),
		"cooldown":       c.cooldown.String(),
		"participants":   len(c.participants),
		"created_at":     c.createdAt,
		"updated_at":     c.updatedAt,
	}
}

// ChallengeParticipant 挑战参与者
type ChallengeParticipant struct {
	PlayerID    string
	ChallengeID string
	StartTime   time.Time
	EndTime     time.Time
	Success     bool
	Score       int
	Attempts    int
	Status      string
}

// NewChallengeParticipant 创建挑战参与者
func NewChallengeParticipant(playerID, challengeID string) *ChallengeParticipant {
	return &ChallengeParticipant{
		PlayerID:    playerID,
		ChallengeID: challengeID,
		StartTime:   time.Now(),
		Attempts:    1,
		Status:      "in_progress",
	}
}

// Complete 完成挑战
func (cp *ChallengeParticipant) Complete(success bool, score int) {
	cp.EndTime = time.Now()
	cp.Success = success
	cp.Score = score
	cp.Status = "completed"
}

// GetDuration 获取用时
func (cp *ChallengeParticipant) GetDuration() time.Duration {
	if cp.EndTime.IsZero() {
		return time.Since(cp.StartTime)
	}
	return cp.EndTime.Sub(cp.StartTime)
}

// ChallengeCondition 挑战条件
type ChallengeCondition struct {
	Type      string
	Field     string
	Operator  string
	Value     interface{}
	Message   string
}

// NewChallengeCondition 创建挑战条件
func NewChallengeCondition(conditionType, field, operator string, value interface{}, message string) *ChallengeCondition {
	return &ChallengeCondition{
		Type:     conditionType,
		Field:    field,
		Operator: operator,
		Value:    value,
		Message:  message,
	}
}

// Check 检查条件
func (cc *ChallengeCondition) Check(data map[string]interface{}) bool {
	value, exists := data[cc.Field]
	if !exists {
		return false
	}
	
	switch cc.Operator {
	case "eq":
		return value == cc.Value
	case "ne":
		return value != cc.Value
	case "gt":
		if v1, ok := value.(int); ok {
			if v2, ok := cc.Value.(int); ok {
				return v1 > v2
			}
		}
	case "gte":
		if v1, ok := value.(int); ok {
			if v2, ok := cc.Value.(int); ok {
				return v1 >= v2
			}
		}
	case "lt":
		if v1, ok := value.(int); ok {
			if v2, ok := cc.Value.(int); ok {
				return v1 < v2
			}
		}
	case "lte":
		if v1, ok := value.(int); ok {
			if v2, ok := cc.Value.(int); ok {
				return v1 <= v2
			}
		}
	}
	
	return false
}

// ChallengeResult 挑战结果
type ChallengeResult struct {
	ChallengeID string
	PlayerID    string
	StartTime   time.Time
	EndTime     time.Time
	Success     bool
	Score       int
	Status      string
	Message     string
}

// ChallengeReward 挑战奖励
type ChallengeReward struct {
	Gold       int
	Experience int
	Items      map[string]int
	Special    map[string]interface{}
}

// NewChallengeReward 创建挑战奖励
func NewChallengeReward(challengeType ChallengeType, difficulty ChallengeDifficulty) *ChallengeReward {
	baseGold := 100
	baseExp := 50
	
	// 根据难度调整基础奖励
	multiplier := difficulty.GetMultiplier()
	
	return &ChallengeReward{
		Gold:       int(float64(baseGold) * multiplier),
		Experience: int(float64(baseExp) * multiplier),
		Items:      make(map[string]int),
		Special:    make(map[string]interface{}),
	}
}

// AddItem 添加物品奖励
func (cr *ChallengeReward) AddItem(itemID string, quantity int) {
	cr.Items[itemID] = quantity
}

// AddSpecial 添加特殊奖励
func (cr *ChallengeReward) AddSpecial(key string, value interface{}) {
	cr.Special[key] = value
}

// Blessing 祝福实体
type Blessing struct {
	id          string
	name        string
	description string
	type_       BlessingType
	effects     []*BlessingEffect
	duration    time.Duration
	cooldown    time.Duration
	status      BlessingStatus
	activatedAt time.Time
	expiresAt   time.Time
	lastUsedAt  time.Time
	usageCount  int
	maxUsage    int
	createdAt   time.Time
	updatedAt   time.Time
}

// NewBlessing 创建祝福
func NewBlessing(id, name, description string, blessingType BlessingType, duration time.Duration) *Blessing {
	now := time.Now()
	return &Blessing{
		id:          id,
		name:        name,
		description: description,
		type_:       blessingType,
		effects:     make([]*BlessingEffect, 0),
		duration:    duration,
		cooldown:    time.Hour * 24, // 默认24小时冷却
		status:      BlessingStatusAvailable,
		maxUsage:    1, // 默认只能使用一次
		createdAt:   now,
		updatedAt:   now,
	}
}

// GetID 获取ID
func (b *Blessing) GetID() string {
	return b.id
}

// GetName 获取名称
func (b *Blessing) GetName() string {
	return b.name
}

// GetDescription 获取描述
func (b *Blessing) GetDescription() string {
	return b.description
}

// GetType 获取类型
func (b *Blessing) GetType() BlessingType {
	return b.type_
}

// GetDuration 获取持续时间
func (b *Blessing) GetDuration() time.Duration {
	return b.duration
}

// GetStatus 获取状态
func (b *Blessing) GetStatus() BlessingStatus {
	return b.status
}

// IsAvailable 检查是否可用
func (b *Blessing) IsAvailable() bool {
	if b.status != BlessingStatusAvailable {
		return false
	}
	
	// 检查使用次数
	if b.usageCount >= b.maxUsage {
		return false
	}
	
	// 检查冷却时间
	if !b.lastUsedAt.IsZero() && time.Since(b.lastUsedAt) < b.cooldown {
		return false
	}
	
	return true
}

// IsActive 检查是否激活
func (b *Blessing) IsActive() bool {
	return b.status == BlessingStatusActive && time.Now().Before(b.expiresAt)
}

// Activate 激活祝福
func (b *Blessing) Activate(playerID string) (*BlessingEffect, error) {
	if !b.IsAvailable() {
		return nil, fmt.Errorf("blessing is not available")
	}
	
	// 更新状态
	b.status = BlessingStatusActive
	b.activatedAt = time.Now()
	b.expiresAt = b.activatedAt.Add(b.duration)
	b.lastUsedAt = b.activatedAt
	b.usageCount++
	b.updatedAt = time.Now()
	
	// 创建祝福效果
	effect := &BlessingEffect{
		BlessingID:  b.id,
		PlayerID:    playerID,
		Type:        b.type_,
		ActivatedAt: b.activatedAt,
		ExpiresAt:   b.expiresAt,
		Effects:     b.effects,
	}
	
	return effect, nil
}

// Deactivate 停用祝福
func (b *Blessing) Deactivate() {
	b.status = BlessingStatusInactive
	b.updatedAt = time.Now()
}

// AddEffect 添加效果
func (b *Blessing) AddEffect(effect *BlessingEffect) {
	b.effects = append(b.effects, effect)
	b.updatedAt = time.Now()
}

// GetRemainingDuration 获取剩余持续时间
func (b *Blessing) GetRemainingDuration() time.Duration {
	if b.status != BlessingStatusActive {
		return 0
	}
	
	if time.Now().After(b.expiresAt) {
		return 0
	}
	
	return b.expiresAt.Sub(time.Now())
}

// GetRemainingCooldown 获取剩余冷却时间
func (b *Blessing) GetRemainingCooldown() time.Duration {
	if b.lastUsedAt.IsZero() {
		return 0
	}
	
	elapsed := time.Since(b.lastUsedAt)
	if elapsed >= b.cooldown {
		return 0
	}
	
	return b.cooldown - elapsed
}

// SetMaxUsage 设置最大使用次数
func (b *Blessing) SetMaxUsage(maxUsage int) {
	b.maxUsage = maxUsage
	b.updatedAt = time.Now()
}

// SetCooldown 设置冷却时间
func (b *Blessing) SetCooldown(cooldown time.Duration) {
	b.cooldown = cooldown
	b.updatedAt = time.Now()
}

// ToMap 转换为映射
func (b *Blessing) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":           b.id,
		"name":         b.name,
		"description":  b.description,
		"type":         b.type_.String(),
		"duration":     b.duration.String(),
		"cooldown":     b.cooldown.String(),
		"status":       b.status.String(),
		"usage_count":  b.usageCount,
		"max_usage":    b.maxUsage,
		"activated_at": b.activatedAt,
		"expires_at":   b.expiresAt,
		"created_at":   b.createdAt,
		"updated_at":   b.updatedAt,
	}
}

// BlessingEffect 祝福效果
type BlessingEffect struct {
	BlessingID  string
	PlayerID    string
	Type        BlessingType
	ActivatedAt time.Time
	ExpiresAt   time.Time
	Effects     []*BlessingEffect
	Attributes  map[string]float64
	Modifiers   map[string]float64
}

// NewBlessingEffect 创建祝福效果
func NewBlessingEffect(blessingID, playerID string, blessingType BlessingType, duration time.Duration) *BlessingEffect {
	now := time.Now()
	return &BlessingEffect{
		BlessingID:  blessingID,
		PlayerID:    playerID,
		Type:        blessingType,
		ActivatedAt: now,
		ExpiresAt:   now.Add(duration),
		Attributes:  make(map[string]float64),
		Modifiers:   make(map[string]float64),
	}
}

// IsActive 检查是否激活
func (be *BlessingEffect) IsActive() bool {
	return time.Now().Before(be.ExpiresAt)
}

// GetRemainingDuration 获取剩余时间
func (be *BlessingEffect) GetRemainingDuration() time.Duration {
	if !be.IsActive() {
		return 0
	}
	return be.ExpiresAt.Sub(time.Now())
}

// AddAttribute 添加属性加成
func (be *BlessingEffect) AddAttribute(name string, value float64) {
	be.Attributes[name] = value
}

// AddModifier 添加修饰符
func (be *BlessingEffect) AddModifier(name string, value float64) {
	be.Modifiers[name] = value
}

// GetAttribute 获取属性值
func (be *BlessingEffect) GetAttribute(name string) float64 {
	return be.Attributes[name]
}

// GetModifier 获取修饰符值
func (be *BlessingEffect) GetModifier(name string) float64 {
	return be.Modifiers[name]
}