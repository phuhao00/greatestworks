package sacred

import (
	"fmt"
	"time"
)

// SacredPlaceAggregate 圣地聚合根
type SacredPlaceAggregate struct {
	id          string
	name        string
	description string
	level       *SacredLevel
	challenges  map[string]*Challenge
	blessings   map[string]*Blessing
	status      SacredStatus
	owner       string
	createdAt   time.Time
	updatedAt   time.Time
	version     int
	events      []DomainEvent
}

// NewSacredPlaceAggregate 创建圣地聚合根
func NewSacredPlaceAggregate(id, name, description, owner string) *SacredPlaceAggregate {
	now := time.Now()
	return &SacredPlaceAggregate{
		id:          id,
		name:        name,
		description: description,
		level:       NewSacredLevel(1, 0),
		challenges:  make(map[string]*Challenge),
		blessings:   make(map[string]*Blessing),
		status:      SacredStatusActive,
		owner:       owner,
		createdAt:   now,
		updatedAt:   now,
		version:     1,
		events:      make([]DomainEvent, 0),
	}
}

// GetID 获取ID
func (s *SacredPlaceAggregate) GetID() string {
	return s.id
}

// GetName 获取名称
func (s *SacredPlaceAggregate) GetName() string {
	return s.name
}

// GetDescription 获取描述
func (s *SacredPlaceAggregate) GetDescription() string {
	return s.description
}

// GetLevel 获取等级
func (s *SacredPlaceAggregate) GetLevel() *SacredLevel {
	return s.level
}

// GetStatus 获取状态
func (s *SacredPlaceAggregate) GetStatus() SacredStatus {
	return s.status
}

// GetOwner 获取拥有者
func (s *SacredPlaceAggregate) GetOwner() string {
	return s.owner
}

// GetVersion 获取版本
func (s *SacredPlaceAggregate) GetVersion() int {
	return s.version
}

// GetEvents 获取领域事件
func (s *SacredPlaceAggregate) GetEvents() []DomainEvent {
	return s.events
}

// ClearEvents 清除领域事件
func (s *SacredPlaceAggregate) ClearEvents() {
	s.events = make([]DomainEvent, 0)
}

// SetName 设置名称
func (s *SacredPlaceAggregate) SetName(name string) error {
	if name == "" {
		return ErrInvalidSacredName
	}
	
	oldName := s.name
	s.name = name
	s.updatedAt = time.Now()
	s.version++
	
	// 发布名称变更事件
	event := NewSacredNameChangedEvent(s.id, oldName, name)
	s.addEvent(event)
	
	return nil
}

// SetDescription 设置描述
func (s *SacredPlaceAggregate) SetDescription(description string) {
	s.description = description
	s.updatedAt = time.Now()
	s.version++
}

// UpgradeLevel 升级等级
func (s *SacredPlaceAggregate) UpgradeLevel(experience int) error {
	if experience <= 0 {
		return ErrInvalidExperience
	}
	
	oldLevel := s.level.Level
	newLevel, err := s.level.AddExperience(experience)
	if err != nil {
		return err
	}
	
	s.updatedAt = time.Now()
	s.version++
	
	// 如果等级提升，发布升级事件
	if newLevel > oldLevel {
		event := NewSacredLevelUpEvent(s.id, oldLevel, newLevel, s.level.Experience)
		s.addEvent(event)
	}
	
	return nil
}

// AddChallenge 添加挑战
func (s *SacredPlaceAggregate) AddChallenge(challenge *Challenge) error {
	if challenge == nil {
		return ErrInvalidChallenge
	}
	
	if _, exists := s.challenges[challenge.GetID()]; exists {
		return ErrChallengeAlreadyExists
	}
	
	// 检查等级要求
	if challenge.GetRequiredLevel() > s.level.Level {
		return ErrInsufficientSacredLevel
	}
	
	s.challenges[challenge.GetID()] = challenge
	s.updatedAt = time.Now()
	s.version++
	
	// 发布挑战添加事件
	event := NewChallengeAddedEvent(s.id, challenge.GetID(), challenge.GetType(), challenge.GetDifficulty())
	s.addEvent(event)
	
	return nil
}

// RemoveChallenge 移除挑战
func (s *SacredPlaceAggregate) RemoveChallenge(challengeID string) error {
	challenge, exists := s.challenges[challengeID]
	if !exists {
		return ErrChallengeNotFound
	}
	
	// 检查挑战是否正在进行
	if challenge.GetStatus() == ChallengeStatusInProgress {
		return ErrChallengeInProgress
	}
	
	delete(s.challenges, challengeID)
	s.updatedAt = time.Now()
	s.version++
	
	// 发布挑战移除事件
	event := NewChallengeRemovedEvent(s.id, challengeID, challenge.GetType())
	s.addEvent(event)
	
	return nil
}

// StartChallenge 开始挑战
func (s *SacredPlaceAggregate) StartChallenge(challengeID, playerID string) (*ChallengeResult, error) {
	challenge, exists := s.challenges[challengeID]
	if !exists {
		return nil, ErrChallengeNotFound
	}
	
	// 检查挑战状态
	if challenge.GetStatus() != ChallengeStatusAvailable {
		return nil, ErrChallengeNotAvailable
	}
	
	// 检查冷却时间
	if !challenge.CanStart() {
		return nil, ErrChallengeOnCooldown
	}
	
	// 开始挑战
	result, err := challenge.Start(playerID)
	if err != nil {
		return nil, err
	}
	
	s.updatedAt = time.Now()
	s.version++
	
	// 发布挑战开始事件
	event := NewChallengeStartedEvent(s.id, challengeID, playerID, challenge.GetType())
	s.addEvent(event)
	
	return result, nil
}

// CompleteChallenge 完成挑战
func (s *SacredPlaceAggregate) CompleteChallenge(challengeID, playerID string, success bool, score int) (*ChallengeReward, error) {
	challenge, exists := s.challenges[challengeID]
	if !exists {
		return nil, ErrChallengeNotFound
	}
	
	// 完成挑战
	reward, err := challenge.Complete(playerID, success, score)
	if err != nil {
		return nil, err
	}
	
	s.updatedAt = time.Now()
	s.version++
	
	// 如果成功，增加经验
	if success {
		s.UpgradeLevel(reward.Experience)
	}
	
	// 发布挑战完成事件
	event := NewChallengeCompletedEvent(s.id, challengeID, playerID, success, score, reward)
	s.addEvent(event)
	
	return reward, nil
}

// AddBlessing 添加祝福
func (s *SacredPlaceAggregate) AddBlessing(blessing *Blessing) error {
	if blessing == nil {
		return ErrInvalidBlessing
	}
	
	if _, exists := s.blessings[blessing.GetID()]; exists {
		return ErrBlessingAlreadyExists
	}
	
	s.blessings[blessing.GetID()] = blessing
	s.updatedAt = time.Now()
	s.version++
	
	// 发布祝福添加事件
	event := NewBlessingAddedEvent(s.id, blessing.GetID(), blessing.GetType(), blessing.GetDuration())
	s.addEvent(event)
	
	return nil
}

// RemoveBlessing 移除祝福
func (s *SacredPlaceAggregate) RemoveBlessing(blessingID string) error {
	blessing, exists := s.blessings[blessingID]
	if !exists {
		return ErrBlessingNotFound
	}
	
	delete(s.blessings, blessingID)
	s.updatedAt = time.Now()
	s.version++
	
	// 发布祝福移除事件
	event := NewBlessingRemovedEvent(s.id, blessingID, blessing.GetType())
	s.addEvent(event)
	
	return nil
}

// ActivateBlessing 激活祝福
func (s *SacredPlaceAggregate) ActivateBlessing(blessingID, playerID string) (*BlessingEffect, error) {
	blessing, exists := s.blessings[blessingID]
	if !exists {
		return nil, ErrBlessingNotFound
	}
	
	// 检查祝福状态
	if !blessing.IsAvailable() {
		return nil, ErrBlessingNotAvailable
	}
	
	// 激活祝福
	effect, err := blessing.Activate(playerID)
	if err != nil {
		return nil, err
	}
	
	s.updatedAt = time.Now()
	s.version++
	
	// 发布祝福激活事件
	event := NewBlessingActivatedEvent(s.id, blessingID, playerID, blessing.GetType(), effect)
	s.addEvent(event)
	
	return effect, nil
}

// GetChallenge 获取挑战
func (s *SacredPlaceAggregate) GetChallenge(challengeID string) (*Challenge, error) {
	challenge, exists := s.challenges[challengeID]
	if !exists {
		return nil, ErrChallengeNotFound
	}
	return challenge, nil
}

// GetAllChallenges 获取所有挑战
func (s *SacredPlaceAggregate) GetAllChallenges() map[string]*Challenge {
	return s.challenges
}

// GetAvailableChallenges 获取可用挑战
func (s *SacredPlaceAggregate) GetAvailableChallenges() []*Challenge {
	var available []*Challenge
	for _, challenge := range s.challenges {
		if challenge.GetStatus() == ChallengeStatusAvailable && challenge.CanStart() {
			available = append(available, challenge)
		}
	}
	return available
}

// GetBlessing 获取祝福
func (s *SacredPlaceAggregate) GetBlessing(blessingID string) (*Blessing, error) {
	blessing, exists := s.blessings[blessingID]
	if !exists {
		return nil, ErrBlessingNotFound
	}
	return blessing, nil
}

// GetAllBlessings 获取所有祝福
func (s *SacredPlaceAggregate) GetAllBlessings() map[string]*Blessing {
	return s.blessings
}

// GetActiveBlessings 获取激活的祝福
func (s *SacredPlaceAggregate) GetActiveBlessings() []*Blessing {
	var active []*Blessing
	for _, blessing := range s.blessings {
		if blessing.IsActive() {
			active = append(active, blessing)
		}
	}
	return active
}

// SetStatus 设置状态
func (s *SacredPlaceAggregate) SetStatus(status SacredStatus) error {
	if !status.IsValid() {
		return ErrInvalidSacredStatus
	}
	
	oldStatus := s.status
	s.status = status
	s.updatedAt = time.Now()
	s.version++
	
	// 发布状态变更事件
	event := NewSacredStatusChangedEvent(s.id, oldStatus, status)
	s.addEvent(event)
	
	return nil
}

// Activate 激活圣地
func (s *SacredPlaceAggregate) Activate() error {
	return s.SetStatus(SacredStatusActive)
}

// Deactivate 停用圣地
func (s *SacredPlaceAggregate) Deactivate() error {
	return s.SetStatus(SacredStatusInactive)
}

// Lock 锁定圣地
func (s *SacredPlaceAggregate) Lock() error {
	return s.SetStatus(SacredStatusLocked)
}

// IsActive 检查是否激活
func (s *SacredPlaceAggregate) IsActive() bool {
	return s.status == SacredStatusActive
}

// CanAccess 检查是否可访问
func (s *SacredPlaceAggregate) CanAccess(playerID string) bool {
	if !s.IsActive() {
		return false
	}
	
	// 检查是否为拥有者或有权限
	return s.owner == playerID || s.hasAccessPermission(playerID)
}

// hasAccessPermission 检查访问权限
func (s *SacredPlaceAggregate) hasAccessPermission(playerID string) bool {
	// 这里可以实现更复杂的权限逻辑
	// 例如：公会成员、好友、VIP等
	return true // 暂时允许所有人访问
}

// GetStatistics 获取统计信息
func (s *SacredPlaceAggregate) GetStatistics() *SacredStatistics {
	totalChallenges := len(s.challenges)
	completedChallenges := 0
	activeBlessings := 0
	
	for _, challenge := range s.challenges {
		if challenge.GetStatus() == ChallengeStatusCompleted {
			completedChallenges++
		}
	}
	
	for _, blessing := range s.blessings {
		if blessing.IsActive() {
			activeBlessings++
		}
	}
	
	return &SacredStatistics{
		SacredID:            s.id,
		Level:               s.level.Level,
		Experience:          s.level.Experience,
		TotalChallenges:     totalChallenges,
		CompletedChallenges: completedChallenges,
		ActiveBlessings:     activeBlessings,
		CreatedAt:           s.createdAt,
		LastActiveAt:        s.updatedAt,
	}
}

// UpdateActivity 更新活动时间
func (s *SacredPlaceAggregate) UpdateActivity() {
	s.updatedAt = time.Now()
	s.version++
}

// addEvent 添加领域事件
func (s *SacredPlaceAggregate) addEvent(event DomainEvent) {
	s.events = append(s.events, event)
}

// ToMap 转换为映射
func (s *SacredPlaceAggregate) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":          s.id,
		"name":        s.name,
		"description": s.description,
		"level":       s.level.ToMap(),
		"status":      s.status.String(),
		"owner":       s.owner,
		"created_at":  s.createdAt,
		"updated_at":  s.updatedAt,
		"version":     s.version,
		"challenges":  len(s.challenges),
		"blessings":   len(s.blessings),
	}
}

// SacredStatus 圣地状态
type SacredStatus int

const (
	SacredStatusActive   SacredStatus = iota + 1 // 激活
	SacredStatusInactive                         // 未激活
	SacredStatusLocked                           // 锁定
	SacredStatusMaintenance                      // 维护中
)

// String 返回状态字符串
func (s SacredStatus) String() string {
	switch s {
	case SacredStatusActive:
		return "active"
	case SacredStatusInactive:
		return "inactive"
	case SacredStatusLocked:
		return "locked"
	case SacredStatusMaintenance:
		return "maintenance"
	default:
		return "unknown"
	}
}

// IsValid 检查状态是否有效
func (s SacredStatus) IsValid() bool {
	return s >= SacredStatusActive && s <= SacredStatusMaintenance
}

// SacredStatistics 圣地统计信息
type SacredStatistics struct {
	SacredID            string
	Level               int
	Experience          int
	TotalChallenges     int
	CompletedChallenges int
	ActiveBlessings     int
	CreatedAt           time.Time
	LastActiveAt        time.Time
}

// 相关错误定义
var (
	ErrInvalidSacredName      = fmt.Errorf("invalid sacred name")
	ErrInvalidExperience      = fmt.Errorf("invalid experience")
	ErrInvalidChallenge       = fmt.Errorf("invalid challenge")
	ErrChallengeAlreadyExists = fmt.Errorf("challenge already exists")
	ErrInsufficientSacredLevel = fmt.Errorf("insufficient sacred level")
	ErrChallengeNotFound      = fmt.Errorf("challenge not found")
	ErrChallengeInProgress    = fmt.Errorf("challenge in progress")
	ErrChallengeNotAvailable  = fmt.Errorf("challenge not available")
	ErrChallengeOnCooldown    = fmt.Errorf("challenge on cooldown")
	ErrInvalidBlessing        = fmt.Errorf("invalid blessing")
	ErrBlessingAlreadyExists  = fmt.Errorf("blessing already exists")
	ErrBlessingNotFound       = fmt.Errorf("blessing not found")
	ErrBlessingNotAvailable   = fmt.Errorf("blessing not available")
	ErrInvalidSacredStatus    = fmt.Errorf("invalid sacred status")
)