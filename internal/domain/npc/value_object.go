package npc

import (
	"fmt"
	"math"
	"time"
)

// NPCType NPC类型
type NPCType int

const (
	NPCTypeVillager  NPCType = iota + 1 // 村民
	NPCTypeMerchant                      // 商人
	NPCTypeGuard                         // 守卫
	NPCTypeQuestGiver                    // 任务发布者
	NPCTypeTrainer                       // 训练师
	NPCTypeBlacksmith                    // 铁匠
	NPCTypeInnkeeper                     // 旅店老板
	NPCTypeLibrarian                     // 图书管理员
	NPCTypeHealer                        // 治疗师
	NPCTypeBanker                        // 银行家
	NPCTypeSpecial                       // 特殊NPC
)

// String 返回类型字符串
func (nt NPCType) String() string {
	switch nt {
	case NPCTypeVillager:
		return "villager"
	case NPCTypeMerchant:
		return "merchant"
	case NPCTypeGuard:
		return "guard"
	case NPCTypeQuestGiver:
		return "quest_giver"
	case NPCTypeTrainer:
		return "trainer"
	case NPCTypeBlacksmith:
		return "blacksmith"
	case NPCTypeInnkeeper:
		return "innkeeper"
	case NPCTypeLibrarian:
		return "librarian"
	case NPCTypeHealer:
		return "healer"
	case NPCTypeBanker:
		return "banker"
	case NPCTypeSpecial:
		return "special"
	default:
		return "unknown"
	}
}

// IsValid 检查类型是否有效
func (nt NPCType) IsValid() bool {
	return nt >= NPCTypeVillager && nt <= NPCTypeSpecial
}

// CanHaveShop 检查是否可以拥有商店
func (nt NPCType) CanHaveShop() bool {
	switch nt {
	case NPCTypeMerchant, NPCTypeBlacksmith, NPCTypeInnkeeper, NPCTypeBanker:
		return true
	default:
		return false
	}
}

// CanGiveQuests 检查是否可以发布任务
func (nt NPCType) CanGiveQuests() bool {
	switch nt {
	case NPCTypeQuestGiver, NPCTypeVillager, NPCTypeGuard, NPCTypeSpecial:
		return true
	default:
		return false
	}
}

// GetDefaultBehavior 获取默认行为
func (nt NPCType) GetDefaultBehavior() BehaviorType {
	switch nt {
	case NPCTypeGuard:
		return BehaviorTypePatrol
	case NPCTypeMerchant, NPCTypeBlacksmith, NPCTypeInnkeeper:
		return BehaviorTypeStationary
	case NPCTypeVillager:
		return BehaviorTypeWander
	default:
		return BehaviorTypeIdle
	}
}

// NPCStatus NPC状态
type NPCStatus int

const (
	NPCStatusActive   NPCStatus = iota + 1 // 激活
	NPCStatusInactive                      // 未激活
	NPCStatusHidden                        // 隐藏
	NPCStatusBusy                          // 忙碌
	NPCStatusSleeping                      // 睡眠
	NPCStatusDead                          // 死亡
)

// String 返回状态字符串
func (ns NPCStatus) String() string {
	switch ns {
	case NPCStatusActive:
		return "active"
	case NPCStatusInactive:
		return "inactive"
	case NPCStatusHidden:
		return "hidden"
	case NPCStatusBusy:
		return "busy"
	case NPCStatusSleeping:
		return "sleeping"
	case NPCStatusDead:
		return "dead"
	default:
		return "unknown"
	}
}

// IsValid 检查状态是否有效
func (ns NPCStatus) IsValid() bool {
	return ns >= NPCStatusActive && ns <= NPCStatusDead
}

// CanInteract 检查是否可以交互
func (ns NPCStatus) CanInteract() bool {
	return ns == NPCStatusActive
}

// IsVisible 检查是否可见
func (ns NPCStatus) IsVisible() bool {
	return ns != NPCStatusHidden && ns != NPCStatusDead
}

// Location 位置值对象
type Location struct {
	X      float64
	Y      float64
	Z      float64
	Region string
	Zone   string
}

// NewLocation 创建位置
func NewLocation(x, y, z float64, region, zone string) *Location {
	return &Location{
		X:      x,
		Y:      y,
		Z:      z,
		Region: region,
		Zone:   zone,
	}
}

// DistanceTo 计算到另一个位置的距离
func (l *Location) DistanceTo(other *Location) float64 {
	dx := l.X - other.X
	dy := l.Y - other.Y
	dz := l.Z - other.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// IsWithinRange 检查是否在指定范围内
func (l *Location) IsWithinRange(other *Location, range_ float64) bool {
	return l.DistanceTo(other) <= range_
}

// MoveTo 移动到指定位置
func (l *Location) MoveTo(x, y, z float64) {
	l.X = x
	l.Y = y
	l.Z = z
}

// ToMap 转换为映射
func (l *Location) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"x":      l.X,
		"y":      l.Y,
		"z":      l.Z,
		"region": l.Region,
		"zone":   l.Zone,
	}
}

// NPCAttributes NPC属性值对象
type NPCAttributes struct {
	Level       int
	Health      int
	MaxHealth   int
	Mana        int
	MaxMana     int
	Strength    int
	Agility     int
	Intelligence int
	Charisma    int
	Luck        int
	MoveSpeed   float64
	ViewRange   float64
	HearRange   float64
}

// NewNPCAttributes 创建NPC属性
func NewNPCAttributes() *NPCAttributes {
	return &NPCAttributes{
		Level:        1,
		Health:       100,
		MaxHealth:    100,
		Mana:         50,
		MaxMana:      50,
		Strength:     10,
		Agility:      10,
		Intelligence: 10,
		Charisma:     10,
		Luck:         10,
		MoveSpeed:    1.0,
		ViewRange:    10.0,
		HearRange:    5.0,
	}
}

// SetLevel 设置等级
func (na *NPCAttributes) SetLevel(level int) {
	na.Level = level
	// 根据等级调整其他属性
	na.MaxHealth = 100 + (level-1)*20
	na.Health = na.MaxHealth
	na.MaxMana = 50 + (level-1)*10
	na.Mana = na.MaxMana
}

// Heal 治疗
func (na *NPCAttributes) Heal(amount int) {
	na.Health += amount
	if na.Health > na.MaxHealth {
		na.Health = na.MaxHealth
	}
}

// TakeDamage 受到伤害
func (na *NPCAttributes) TakeDamage(damage int) {
	na.Health -= damage
	if na.Health < 0 {
		na.Health = 0
	}
}

// IsAlive 检查是否存活
func (na *NPCAttributes) IsAlive() bool {
	return na.Health > 0
}

// GetHealthPercentage 获取生命值百分比
func (na *NPCAttributes) GetHealthPercentage() float64 {
	if na.MaxHealth == 0 {
		return 0
	}
	return float64(na.Health) / float64(na.MaxHealth)
}

// GetManaPercentage 获取法力值百分比
func (na *NPCAttributes) GetManaPercentage() float64 {
	if na.MaxMana == 0 {
		return 0
	}
	return float64(na.Mana) / float64(na.MaxMana)
}

// ToMap 转换为映射
func (na *NPCAttributes) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"level":        na.Level,
		"health":       na.Health,
		"max_health":   na.MaxHealth,
		"mana":         na.Mana,
		"max_mana":     na.MaxMana,
		"strength":     na.Strength,
		"agility":      na.Agility,
		"intelligence": na.Intelligence,
		"charisma":     na.Charisma,
		"luck":         na.Luck,
		"move_speed":   na.MoveSpeed,
		"view_range":   na.ViewRange,
		"hear_range":   na.HearRange,
	}
}

// NPCBehavior NPC行为值对象
type NPCBehavior struct {
	Type         BehaviorType
	State        BehaviorState
	PatrolPoints []*Location
	CurrentPoint int
	Target       *Location
	MoveSpeed    float64
	PauseTime    time.Duration
	LastMove     time.Time
	CanMove      bool
	CanTalk      bool
	CanFight     bool
}

// NewNPCBehavior 创建NPC行为
func NewNPCBehavior() *NPCBehavior {
	return &NPCBehavior{
		Type:         BehaviorTypeIdle,
		State:        BehaviorStateIdle,
		PatrolPoints: make([]*Location, 0),
		MoveSpeed:    1.0,
		PauseTime:    time.Second * 3,
		LastMove:     time.Now(),
		CanMove:      true,
		CanTalk:      true,
		CanFight:     false,
	}
}

// SetBehaviorType 设置行为类型
func (nb *NPCBehavior) SetBehaviorType(behaviorType BehaviorType) {
	nb.Type = behaviorType
	nb.State = BehaviorStateIdle
}

// AddPatrolPoint 添加巡逻点
func (nb *NPCBehavior) AddPatrolPoint(location *Location) {
	nb.PatrolPoints = append(nb.PatrolPoints, location)
}

// GetNextPatrolPoint 获取下一个巡逻点
func (nb *NPCBehavior) GetNextPatrolPoint() *Location {
	if len(nb.PatrolPoints) == 0 {
		return nil
	}
	
	nb.CurrentPoint = (nb.CurrentPoint + 1) % len(nb.PatrolPoints)
	return nb.PatrolPoints[nb.CurrentPoint]
}

// SetTarget 设置目标
func (nb *NPCBehavior) SetTarget(target *Location) {
	nb.Target = target
	nb.State = BehaviorStateMoving
}

// ClearTarget 清除目标
func (nb *NPCBehavior) ClearTarget() {
	nb.Target = nil
	nb.State = BehaviorStateIdle
}

// CanMove 检查是否可以移动
func (nb *NPCBehavior) CanMove() bool {
	return nb.CanMove && nb.State != BehaviorStatePaused
}

// Update 更新行为
func (nb *NPCBehavior) Update(deltaTime time.Duration) {
	switch nb.Type {
	case BehaviorTypePatrol:
		nb.updatePatrol(deltaTime)
	case BehaviorTypeWander:
		nb.updateWander(deltaTime)
	case BehaviorTypeFollow:
		nb.updateFollow(deltaTime)
	default:
		nb.updateIdle(deltaTime)
	}
}

// updatePatrol 更新巡逻行为
func (nb *NPCBehavior) updatePatrol(deltaTime time.Duration) {
	if len(nb.PatrolPoints) == 0 {
		return
	}
	
	switch nb.State {
	case BehaviorStateIdle:
		if time.Since(nb.LastMove) >= nb.PauseTime {
			nb.SetTarget(nb.GetNextPatrolPoint())
		}
	case BehaviorStateMoving:
		// 移动逻辑在这里实现
		nb.LastMove = time.Now()
	}
}

// updateWander 更新漫游行为
func (nb *NPCBehavior) updateWander(deltaTime time.Duration) {
	// 简化的漫游逻辑
	if nb.State == BehaviorStateIdle && time.Since(nb.LastMove) >= nb.PauseTime {
		// 随机选择一个方向移动
		nb.State = BehaviorStateMoving
		nb.LastMove = time.Now()
	}
}

// updateFollow 更新跟随行为
func (nb *NPCBehavior) updateFollow(deltaTime time.Duration) {
	// 跟随目标的逻辑
	if nb.Target != nil {
		nb.State = BehaviorStateMoving
	}
}

// updateIdle 更新空闲行为
func (nb *NPCBehavior) updateIdle(deltaTime time.Duration) {
	// 空闲状态不需要特殊处理
}

// BehaviorType 行为类型
type BehaviorType int

const (
	BehaviorTypeIdle       BehaviorType = iota + 1 // 空闲
	BehaviorTypePatrol                             // 巡逻
	BehaviorTypeWander                             // 漫游
	BehaviorTypeFollow                             // 跟随
	BehaviorTypeStationary                         // 固定
	BehaviorTypeAggressive                         // 攻击性
	BehaviorTypeDefensive                          // 防御性
)

// String 返回行为类型字符串
func (bt BehaviorType) String() string {
	switch bt {
	case BehaviorTypeIdle:
		return "idle"
	case BehaviorTypePatrol:
		return "patrol"
	case BehaviorTypeWander:
		return "wander"
	case BehaviorTypeFollow:
		return "follow"
	case BehaviorTypeStationary:
		return "stationary"
	case BehaviorTypeAggressive:
		return "aggressive"
	case BehaviorTypeDefensive:
		return "defensive"
	default:
		return "unknown"
	}
}

// BehaviorState 行为状态
type BehaviorState int

const (
	BehaviorStateIdle    BehaviorState = iota + 1 // 空闲
	BehaviorStateMoving                            // 移动中
	BehaviorStatePaused                            // 暂停
	BehaviorStateWaiting                           // 等待
	BehaviorStateActing                            // 执行动作
)

// String 返回行为状态字符串
func (bs BehaviorState) String() string {
	switch bs {
	case BehaviorStateIdle:
		return "idle"
	case BehaviorStateMoving:
		return "moving"
	case BehaviorStatePaused:
		return "paused"
	case BehaviorStateWaiting:
		return "waiting"
	case BehaviorStateActing:
		return "acting"
	default:
		return "unknown"
	}
}

// Relationship 关系值对象
type Relationship struct {
	PlayerID  string
	NPCID     string
	Value     int
	Level     RelationshipLevel
	History   []*RelationshipEvent
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewRelationship 创建关系
func NewRelationship(playerID, npcID string) *Relationship {
	now := time.Now()
	return &Relationship{
		PlayerID:  playerID,
		NPCID:     npcID,
		Value:     0,
		Level:     RelationshipLevelNeutral,
		History:   make([]*RelationshipEvent, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// GetValue 获取关系值
func (r *Relationship) GetValue() int {
	return r.Value
}

// GetLevel 获取关系等级
func (r *Relationship) GetLevel() RelationshipLevel {
	return r.Level
}

// ChangeValue 改变关系值
func (r *Relationship) ChangeValue(change int, reason string) error {
	oldValue := r.Value
	oldLevel := r.Level
	
	r.Value += change
	
	// 限制关系值范围
	if r.Value > 1000 {
		r.Value = 1000
	} else if r.Value < -1000 {
		r.Value = -1000
	}
	
	// 更新关系等级
	r.updateLevel()
	
	// 记录历史
	event := &RelationshipEvent{
		Reason:      reason,
		Change:      change,
		OldValue:    oldValue,
		NewValue:    r.Value,
		OldLevel:    oldLevel,
		NewLevel:    r.Level,
		Timestamp:   time.Now(),
	}
	r.History = append(r.History, event)
	
	r.UpdatedAt = time.Now()
	return nil
}

// updateLevel 更新关系等级
func (r *Relationship) updateLevel() {
	switch {
	case r.Value >= 500:
		r.Level = RelationshipLevelRevered
	case r.Value >= 200:
		r.Level = RelationshipLevelFriendly
	case r.Value >= 50:
		r.Level = RelationshipLevelLiked
	case r.Value >= -50:
		r.Level = RelationshipLevelNeutral
	case r.Value >= -200:
		r.Level = RelationshipLevelDisliked
	case r.Value >= -500:
		r.Level = RelationshipLevelUnfriendly
	default:
		r.Level = RelationshipLevelHostile
	}
}

// GetRecentHistory 获取最近的历史记录
func (r *Relationship) GetRecentHistory(limit int) []*RelationshipEvent {
	if len(r.History) <= limit {
		return r.History
	}
	return r.History[len(r.History)-limit:]
}

// RelationshipLevel 关系等级
type RelationshipLevel int

const (
	RelationshipLevelHostile    RelationshipLevel = iota + 1 // 敌对
	RelationshipLevelUnfriendly                              // 不友好
	RelationshipLevelDisliked                                // 不喜欢
	RelationshipLevelNeutral                                 // 中立
	RelationshipLevelLiked                                   // 喜欢
	RelationshipLevelFriendly                                // 友好
	RelationshipLevelRevered                                 // 崇敬
)

// String 返回关系等级字符串
func (rl RelationshipLevel) String() string {
	switch rl {
	case RelationshipLevelHostile:
		return "hostile"
	case RelationshipLevelUnfriendly:
		return "unfriendly"
	case RelationshipLevelDisliked:
		return "disliked"
	case RelationshipLevelNeutral:
		return "neutral"
	case RelationshipLevelLiked:
		return "liked"
	case RelationshipLevelFriendly:
		return "friendly"
	case RelationshipLevelRevered:
		return "revered"
	default:
		return "unknown"
	}
}

// GetColor 获取关系等级颜色
func (rl RelationshipLevel) GetColor() string {
	switch rl {
	case RelationshipLevelHostile:
		return "red"
	case RelationshipLevelUnfriendly:
		return "orange"
	case RelationshipLevelDisliked:
		return "yellow"
	case RelationshipLevelNeutral:
		return "white"
	case RelationshipLevelLiked:
		return "lightgreen"
	case RelationshipLevelFriendly:
		return "green"
	case RelationshipLevelRevered:
		return "gold"
	default:
		return "gray"
	}
}

// RelationshipEvent 关系事件
type RelationshipEvent struct {
	Reason    string
	Change    int
	OldValue  int
	NewValue  int
	OldLevel  RelationshipLevel
	NewLevel  RelationshipLevel
	Timestamp time.Time
}

// NPCSchedule NPC日程值对象
type NPCSchedule struct {
	ScheduleItems []*ScheduleItem
	CurrentItem   *ScheduleItem
	TimeZone      string
}

// NewNPCSchedule 创建NPC日程
func NewNPCSchedule() *NPCSchedule {
	return &NPCSchedule{
		ScheduleItems: make([]*ScheduleItem, 0),
		TimeZone:      "UTC",
	}
}

// AddScheduleItem 添加日程项
func (ns *NPCSchedule) AddScheduleItem(item *ScheduleItem) {
	ns.ScheduleItems = append(ns.ScheduleItems, item)
}

// GetCurrentItem 获取当前日程项
func (ns *NPCSchedule) GetCurrentItem(currentTime time.Time) *ScheduleItem {
	for _, item := range ns.ScheduleItems {
		if item.IsActive(currentTime) {
			return item
		}
	}
	return nil
}

// Update 更新日程
func (ns *NPCSchedule) Update(currentTime time.Time) {
	ns.CurrentItem = ns.GetCurrentItem(currentTime)
}

// ScheduleItem 日程项
type ScheduleItem struct {
	ID          string
	Name        string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	DayOfWeek   []time.Weekday
	Location    *Location
	Behavior    BehaviorType
	Actions     []string
	Priority    int
}

// NewScheduleItem 创建日程项
func NewScheduleItem(id, name, description string, startTime, endTime time.Time) *ScheduleItem {
	return &ScheduleItem{
		ID:          id,
		Name:        name,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
		DayOfWeek:   make([]time.Weekday, 0),
		Actions:     make([]string, 0),
		Priority:    1,
	}
}

// IsActive 检查是否激活
func (si *ScheduleItem) IsActive(currentTime time.Time) bool {
	// 检查星期几
	if len(si.DayOfWeek) > 0 {
		currentWeekday := currentTime.Weekday()
		found := false
		for _, weekday := range si.DayOfWeek {
			if weekday == currentWeekday {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	// 检查时间范围
	currentHour := currentTime.Hour()
	currentMinute := currentTime.Minute()
	currentTimeOfDay := currentHour*60 + currentMinute
	
	startTimeOfDay := si.StartTime.Hour()*60 + si.StartTime.Minute()
	endTimeOfDay := si.EndTime.Hour()*60 + si.EndTime.Minute()
	
	return currentTimeOfDay >= startTimeOfDay && currentTimeOfDay <= endTimeOfDay
}

// AddDayOfWeek 添加星期几
func (si *ScheduleItem) AddDayOfWeek(weekday time.Weekday) {
	si.DayOfWeek = append(si.DayOfWeek, weekday)
}

// AddAction 添加动作
func (si *ScheduleItem) AddAction(action string) {
	si.Actions = append(si.Actions, action)
}

// ShopSchedule 商店日程值对象
type ShopSchedule struct {
	OpenTime    time.Time
	CloseTime   time.Time
	DaysOpen    []time.Weekday
	Holidays    []time.Time
	SpecialHours map[string]*SpecialHours
}

// NewShopSchedule 创建商店日程
func NewShopSchedule() *ShopSchedule {
	return &ShopSchedule{
		OpenTime:     time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),  // 9:00 AM
		CloseTime:    time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC), // 6:00 PM
		DaysOpen:     []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
		Holidays:     make([]time.Time, 0),
		SpecialHours: make(map[string]*SpecialHours),
	}
}

// IsOpen 检查是否开放
func (ss *ShopSchedule) IsOpen(currentTime time.Time) bool {
	// 检查是否为假日
	for _, holiday := range ss.Holidays {
		if currentTime.YearDay() == holiday.YearDay() && currentTime.Year() == holiday.Year() {
			return false
		}
	}
	
	// 检查特殊时间
	dateKey := currentTime.Format("2006-01-02")
	if specialHours, exists := ss.SpecialHours[dateKey]; exists {
		return specialHours.IsOpen(currentTime)
	}
	
	// 检查星期几
	currentWeekday := currentTime.Weekday()
	found := false
	for _, weekday := range ss.DaysOpen {
		if weekday == currentWeekday {
			found = true
			break
		}
	}
	if !found {
		return false
	}
	
	// 检查营业时间
	currentHour := currentTime.Hour()
	currentMinute := currentTime.Minute()
	currentTimeOfDay := currentHour*60 + currentMinute
	
	openTimeOfDay := ss.OpenTime.Hour()*60 + ss.OpenTime.Minute()
	closeTimeOfDay := ss.CloseTime.Hour()*60 + ss.CloseTime.Minute()
	
	return currentTimeOfDay >= openTimeOfDay && currentTimeOfDay <= closeTimeOfDay
}

// AddHoliday 添加假日
func (ss *ShopSchedule) AddHoliday(holiday time.Time) {
	ss.Holidays = append(ss.Holidays, holiday)
}

// SetSpecialHours 设置特殊时间
func (ss *ShopSchedule) SetSpecialHours(date string, specialHours *SpecialHours) {
	ss.SpecialHours[date] = specialHours
}

// SpecialHours 特殊时间
type SpecialHours struct {
	OpenTime  time.Time
	CloseTime time.Time
	Closed    bool
}

// NewSpecialHours 创建特殊时间
func NewSpecialHours(openTime, closeTime time.Time, closed bool) *SpecialHours {
	return &SpecialHours{
		OpenTime:  openTime,
		CloseTime: closeTime,
		Closed:    closed,
	}
}

// IsOpen 检查是否开放
func (sh *SpecialHours) IsOpen(currentTime time.Time) bool {
	if sh.Closed {
		return false
	}
	
	currentHour := currentTime.Hour()
	currentMinute := currentTime.Minute()
	currentTimeOfDay := currentHour*60 + currentMinute
	
	openTimeOfDay := sh.OpenTime.Hour()*60 + sh.OpenTime.Minute()
	closeTimeOfDay := sh.CloseTime.Hour()*60 + sh.CloseTime.Minute()
	
	return currentTimeOfDay >= openTimeOfDay && currentTimeOfDay <= closeTimeOfDay
}

// 枚举类型定义

// DialogueType 对话类型
type DialogueType int

const (
	DialogueTypeGreeting   DialogueType = iota + 1 // 问候
	DialogueTypeInformation                         // 信息
	DialogueTypeQuest                               // 任务
	DialogueTypeTrade                               // 交易
	DialogueTypeRumor                               // 传言
	DialogueTypeStory                               // 故事
	DialogueTypeSpecial                             // 特殊
)

// String 返回对话类型字符串
func (dt DialogueType) String() string {
	switch dt {
	case DialogueTypeGreeting:
		return "greeting"
	case DialogueTypeInformation:
		return "information"
	case DialogueTypeQuest:
		return "quest"
	case DialogueTypeTrade:
		return "trade"
	case DialogueTypeRumor:
		return "rumor"
	case DialogueTypeStory:
		return "story"
	case DialogueTypeSpecial:
		return "special"
	default:
		return "unknown"
	}
}

// QuestType 任务类型
type QuestType int

const (
	QuestTypeKill     QuestType = iota + 1 // 击杀
	QuestTypeCollect                       // 收集
	QuestTypeDeliver                       // 运送
	QuestTypeEscort                        // 护送
	QuestTypeExplore                       // 探索
	QuestTypeTalk                          // 对话
	QuestTypeCraft                         // 制作
	QuestTypeDaily                         // 日常
	QuestTypeWeekly                        // 周常
	QuestTypeSpecial                       // 特殊
)

// String 返回任务类型字符串
func (qt QuestType) String() string {
	switch qt {
	case QuestTypeKill:
		return "kill"
	case QuestTypeCollect:
		return "collect"
	case QuestTypeDeliver:
		return "deliver"
	case QuestTypeEscort:
		return "escort"
	case QuestTypeExplore:
		return "explore"
	case QuestTypeTalk:
		return "talk"
	case QuestTypeCraft:
		return "craft"
	case QuestTypeDaily:
		return "daily"
	case QuestTypeWeekly:
		return "weekly"
	case QuestTypeSpecial:
		return "special"
	default:
		return "unknown"
	}
}

// QuestStatus 任务状态
type QuestStatus int

const (
	QuestStatusActive    QuestStatus = iota + 1 // 激活
	QuestStatusCompleted                         // 完成
	QuestStatusFailed                            // 失败
	QuestStatusAbandoned                         // 放弃
	QuestStatusExpired                           // 过期
)

// String 返回任务状态字符串
func (qs QuestStatus) String() string {
	switch qs {
	case QuestStatusActive:
		return "active"
	case QuestStatusCompleted:
		return "completed"
	case QuestStatusFailed:
		return "failed"
	case QuestStatusAbandoned:
		return "abandoned"
	case QuestStatusExpired:
		return "expired"
	default:
		return "unknown"
	}
}

// ObjectiveType 目标类型
type ObjectiveType int

const (
	ObjectiveTypeKill     ObjectiveType = iota + 1 // 击杀
	ObjectiveTypeCollect                            // 收集
	ObjectiveTypeDeliver                            // 运送
	ObjectiveTypeReach                              // 到达
	ObjectiveTypeInteract                           // 交互
	ObjectiveTypeWait                               // 等待
	ObjectiveTypeDefend                             // 防御
	ObjectiveTypeEscape                             // 逃脱
)

// String 返回目标类型字符串
func (ot ObjectiveType) String() string {
	switch ot {
	case ObjectiveTypeKill:
		return "kill"
	case ObjectiveTypeCollect:
		return "collect"
	case ObjectiveTypeDeliver:
		return "deliver"
	case ObjectiveTypeReach:
		return "reach"
	case ObjectiveTypeInteract:
		return "interact"
	case ObjectiveTypeWait:
		return "wait"
	case ObjectiveTypeDefend:
		return "defend"
	case ObjectiveTypeEscape:
		return "escape"
	default:
		return "unknown"
	}
}

// ConditionType 条件类型
type ConditionType int

const (
	ConditionTypeLevel        ConditionType = iota + 1 // 等级
	ConditionTypeItem                                   // 物品
	ConditionTypeQuest                                  // 任务
	ConditionTypeRelationship                           // 关系
	ConditionTypeTime                                   // 时间
	ConditionTypeLocation                               // 位置
	ConditionTypeAttribute                              // 属性
	ConditionTypeCustom                                 // 自定义
)

// String 返回条件类型字符串
func (ct ConditionType) String() string {
	switch ct {
	case ConditionTypeLevel:
		return "level"
	case ConditionTypeItem:
		return "item"
	case ConditionTypeQuest:
		return "quest"
	case ConditionTypeRelationship:
		return "relationship"
	case ConditionTypeTime:
		return "time"
	case ConditionTypeLocation:
		return "location"
	case ConditionTypeAttribute:
		return "attribute"
	case ConditionTypeCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// ActionType 动作类型
type ActionType int

const (
	ActionTypeGiveItem       ActionType = iota + 1 // 给予物品
	ActionTypeTakeItem                              // 拿取物品
	ActionTypeGiveGold                              // 给予金币
	ActionTypeTakeGold                              // 拿取金币
	ActionTypeGiveExperience                        // 给予经验
	ActionTypeStartQuest                            // 开始任务
	ActionTypeCompleteQuest                         // 完成任务
	ActionTypeChangeRelationship                    // 改变关系
	ActionTypeTeleport                              // 传送
	ActionTypeCustom                                // 自定义
)

// String 返回动作类型字符串
func (at ActionType) String() string {
	switch at {
	case ActionTypeGiveItem:
		return "give_item"
	case ActionTypeTakeItem:
		return "take_item"
	case ActionTypeGiveGold:
		return "give_gold"
	case ActionTypeTakeGold:
		return "take_gold"
	case ActionTypeGiveExperience:
		return "give_experience"
	case ActionTypeStartQuest:
		return "start_quest"
	case ActionTypeCompleteQuest:
		return "complete_quest"
	case ActionTypeChangeRelationship:
		return "change_relationship"
	case ActionTypeTeleport:
		return "teleport"
	case ActionTypeCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// PrerequisiteType 前置条件类型
type PrerequisiteType int

const (
	PrerequisiteTypeLevel        PrerequisiteType = iota + 1 // 等级
	PrerequisiteTypeQuest                                     // 任务
	PrerequisiteTypeItem                                      // 物品
	PrerequisiteTypeRelationship                              // 关系
	PrerequisiteTypeAttribute                                 // 属性
	PrerequisiteTypeTime                                      // 时间
	PrerequisiteTypeCustom                                    // 自定义
)

// String 返回前置条件类型字符串
func (pt PrerequisiteType) String() string {
	switch pt {
	case PrerequisiteTypeLevel:
		return "level"
	case PrerequisiteTypeQuest:
		return "quest"
	case PrerequisiteTypeItem:
		return "item"
	case PrerequisiteTypeRelationship:
		return "relationship"
	case PrerequisiteTypeAttribute:
		return "attribute"
	case PrerequisiteTypeTime:
		return "time"
	case PrerequisiteTypeCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// DiscountType 折扣类型
type DiscountType int

const (
	DiscountTypePercentage DiscountType = iota + 1 // 百分比
	DiscountTypeFixed                               // 固定金额
)

// String 返回折扣类型字符串
func (dt DiscountType) String() string {
	switch dt {
	case DiscountTypePercentage:
		return "percentage"
	case DiscountTypeFixed:
		return "fixed"
	default:
		return "unknown"
	}
}