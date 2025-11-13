//nolint:staticcheck // Quest rich model contains fields reserved for future use; suppress U1000 unused field warnings in early stage
package quest

import (
	// "errors"
	"time"
)

// QuestManager 任务管理器聚合根
type QuestManager struct {
	playerID        string
	activeQuests    map[string]*Quest
	completedQuests map[string]*Quest
	dailyQuests     map[string]*Quest
	weeklyQuests    map[string]*Quest
	achievements    map[string]*Achievement
	lastUpdate      time.Time
	events          []DomainEvent
}

// NewQuestManager 创建新任务管理器
func NewQuestManager(playerID string) *QuestManager {
	return &QuestManager{
		playerID:        playerID,
		activeQuests:    make(map[string]*Quest),
		completedQuests: make(map[string]*Quest),
		dailyQuests:     make(map[string]*Quest),
		weeklyQuests:    make(map[string]*Quest),
		achievements:    make(map[string]*Achievement),
		lastUpdate:      time.Now(),
		events:          make([]DomainEvent, 0),
	}
}

// Quest 任务实体
type Quest struct {
	ID                string
	Name              string
	Description       string
	QuestType         QuestType
	Category          QuestCategory
	Status            QuestStatus
	Priority          QuestPriority
	Objectives        []*QuestObjective
	Rewards           []*QuestReward
	Prerequisites     []string // 前置任务ID
	StartTime         *time.Time
	ExpireTime        *time.Time
	CompletedTime     *time.Time
	TimeLimit         *time.Duration
	RepeatType        RepeatType
	RepeatCount       int
	MaxRepeats        int
	Level             int
	MinLevel          int
	MaxLevel          int
	ClassRestrictions []string
	RaceRestrictions  []string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// NewQuest 创建新任务
func NewQuest(id, name string, questType QuestType) *Quest {
	return &Quest{
		ID:          id,
		Name:        name,
		QuestType:   questType,
		Status:      QuestStatusAvailable,
		Priority:    QuestPriorityNormal,
		Objectives:  make([]*QuestObjective, 0),
		Rewards:     make([]*QuestReward, 0),
		RepeatType:  RepeatTypeNone,
		RepeatCount: 0,
		MaxRepeats:  1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// QuestType 任务类型
type QuestType int

const (
	QuestTypeMain QuestType = iota + 1
	QuestTypeSide
	QuestTypeDaily
	QuestTypeWeekly
	QuestTypeEvent
	QuestTypeGuild
	QuestTypePvP
	QuestTypeDungeon
	QuestTypeRaid
)

// QuestCategory 任务分类
type QuestCategory int

const (
	QuestCategoryKill QuestCategory = iota + 1
	QuestCategoryCollect
	QuestCategoryDeliver
	QuestCategoryEscort
	QuestCategoryExplore
	QuestCategoryCraft
	QuestCategoryTalk
	QuestCategoryUse
	QuestCategoryReach
)

// QuestStatus 任务状态
type QuestStatus int

const (
	QuestStatusAvailable QuestStatus = iota + 1
	QuestStatusAccepted
	QuestStatusInProgress
	QuestStatusCompleted
	QuestStatusFailed
	QuestStatusExpired
	QuestStatusAbandoned
)

// QuestPriority 任务优先级
type QuestPriority int

const (
	QuestPriorityLow QuestPriority = iota + 1
	QuestPriorityNormal
	QuestPriorityHigh
	QuestPriorityUrgent
)

// RepeatType 重复类型
type RepeatType int

const (
	RepeatTypeNone RepeatType = iota
	RepeatTypeDaily
	RepeatTypeWeekly
	RepeatTypeMonthly
	RepeatTypeUnlimited
)

// QuestObjective 任务目标
type QuestObjective struct {
	ID            string
	Description   string
	ObjectiveType ObjectiveType
	Target        string // 目标ID或名称
	Current       int64
	Required      int64
	Completed     bool
	Optional      bool
	Order         int
	Metadata      map[string]interface{}
}

// ObjectiveType 目标类型
type ObjectiveType int

const (
	ObjectiveTypeKill ObjectiveType = iota + 1
	ObjectiveTypeCollect
	ObjectiveTypeDeliver
	ObjectiveTypeReach
	ObjectiveTypeUse
	ObjectiveTypeTalk
	ObjectiveTypeCraft
	ObjectiveTypeLevel
	ObjectiveTypeEquip
	ObjectiveTypeSpend
	ObjectiveTypeEarn
)

// QuestReward 任务奖励
type QuestReward struct {
	RewardType RewardType
	RewardID   string
	Quantity   int64
	Optional   bool
	Condition  *RewardCondition
}

// RewardType 奖励类型
type RewardType int

const (
	RewardTypeExperience RewardType = iota + 1
	RewardTypeGold
	RewardTypeItem
	RewardTypeSkillPoint
	RewardTypeReputation
	RewardTypeTitle
	RewardTypeAchievement
	RewardTypeBuff
)

// RewardCondition 奖励条件
type RewardCondition struct {
	ConditionType ConditionType
	Value         interface{}
}

// ConditionType 条件类型
type ConditionType int

const (
	ConditionTypeLevel ConditionType = iota + 1
	ConditionTypeClass
	ConditionTypeRace
	ConditionTypeGuild
	ConditionTypeTime
	ConditionTypeRandom
)

// Achievement 成就实体
type Achievement struct {
	ID            string
	Name          string
	Description   string
	Category      AchievementCategory
	Points        int64
	Requirements  []*AchievementRequirement
	Rewards       []*QuestReward
	Unlocked      bool
	Progress      int64
	TotalProgress int64
	UnlockedAt    *time.Time
	Hidden        bool
	Rare          bool
}

// AchievementCategory 成就分类
type AchievementCategory int

const (
	AchievementCategoryGeneral AchievementCategory = iota + 1
	AchievementCategoryCombat
	AchievementCategoryExploration
	AchievementCategoryCrafting
	AchievementCategorySocial
	AchievementCategoryPvP
	AchievementCategoryPvE
	AchievementCategoryCollection
)

// AchievementRequirement 成就要求
type AchievementRequirement struct {
	RequirementType RequirementType
	Target          string
	Value           int64
	Current         int64
	Completed       bool
}

// RequirementType 要求类型
type RequirementType int

const (
	RequirementTypeKill RequirementType = iota + 1
	RequirementTypeCollect
	RequirementTypeComplete
	RequirementTypeReach
	RequirementTypeSpend
	RequirementTypeEarn
	RequirementTypeUse
	RequirementTypeCraft
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
	PlayerID() string
}

// QuestAcceptedEvent 任务接受事件
type QuestAcceptedEvent struct {
	Player         string
	QuestID        string
	QuestName      string
	OccurredAtTime time.Time
}

func (e QuestAcceptedEvent) EventType() string     { return "quest.accepted" }
func (e QuestAcceptedEvent) OccurredAt() time.Time { return e.OccurredAtTime }
func (e QuestAcceptedEvent) PlayerID() string      { return e.Player }

// QuestCompletedEvent 任务完成事件
type QuestCompletedEvent struct {
	Player         string
	QuestID        string
	QuestName      string
	Rewards        []*QuestReward
	OccurredAtTime time.Time
}

func (e QuestCompletedEvent) EventType() string     { return "quest.completed" }
func (e QuestCompletedEvent) OccurredAt() time.Time { return e.OccurredAtTime }
func (e QuestCompletedEvent) PlayerID() string      { return e.Player }

// QuestFailedEvent 任务失败事件
type QuestFailedEvent struct {
	Player         string
	QuestID        string
	QuestName      string
	Reason         string
	OccurredAtTime time.Time
}

func (e QuestFailedEvent) EventType() string     { return "quest.failed" }
func (e QuestFailedEvent) OccurredAt() time.Time { return e.OccurredAtTime }
func (e QuestFailedEvent) PlayerID() string      { return e.Player }

// ObjectiveCompletedEvent 目标完成事件
type ObjectiveCompletedEvent struct {
	Player         string
	QuestID        string
	ObjectiveID    string
	ObjectiveName  string
	OccurredAtTime time.Time
}

func (e ObjectiveCompletedEvent) EventType() string     { return "objective.completed" }
func (e ObjectiveCompletedEvent) OccurredAt() time.Time { return e.OccurredAtTime }
func (e ObjectiveCompletedEvent) PlayerID() string      { return e.Player }

// AchievementUnlockedEvent 成就解锁事件
type AchievementUnlockedEvent struct {
	Player          string
	AchievementID   string
	AchievementName string
	Points          int64
	OccurredAtTime  time.Time
}

func (e AchievementUnlockedEvent) EventType() string     { return "achievement.unlocked" }
func (e AchievementUnlockedEvent) OccurredAt() time.Time { return e.OccurredAtTime }
func (e AchievementUnlockedEvent) PlayerID() string      { return e.Player }

// QuestManager 业务方法

// PlayerID 获取玩家ID
func (qm *QuestManager) PlayerID() string {
	return qm.playerID
}

// ActiveQuests 获取活跃任务
func (qm *QuestManager) ActiveQuests() map[string]*Quest {
	return qm.activeQuests
}

// CompletedQuests 获取已完成任务
func (qm *QuestManager) CompletedQuests() map[string]*Quest {
	return qm.completedQuests
}

// Achievements 获取成就
func (qm *QuestManager) Achievements() map[string]*Achievement {
	return qm.achievements
}

// AcceptQuest 接受任务
func (qm *QuestManager) AcceptQuest(quest *Quest) error {
	// 检查任务是否可接受
	if quest.Status != QuestStatusAvailable {
		return ErrQuestNotAvailable
	}

	// 检查是否已接受
	if _, exists := qm.activeQuests[quest.ID]; exists {
		return ErrQuestAlreadyAccepted
	}

	// 检查前置条件
	if !qm.checkPrerequisites(quest.Prerequisites) {
		return ErrPrerequisitesNotMet
	}

	// 检查等级限制
	if quest.MinLevel > 0 || quest.MaxLevel > 0 {
		// 这里需要获取玩家等级，暂时跳过
	}

	// 接受任务
	quest.Status = QuestStatusAccepted
	now := time.Now()
	quest.StartTime = &now
	quest.UpdatedAt = time.Now()

	// 设置过期时间
	if quest.TimeLimit != nil {
		expireTime := now.Add(*quest.TimeLimit)
		quest.ExpireTime = &expireTime
	}

	qm.activeQuests[quest.ID] = quest
	qm.lastUpdate = time.Now()

	// 发布事件
	qm.addEvent(QuestAcceptedEvent{
		Player:         qm.playerID,
		QuestID:        quest.ID,
		QuestName:      quest.Name,
		OccurredAtTime: time.Now(),
	})

	return nil
}

// UpdateObjectiveProgress 更新目标进度
func (qm *QuestManager) UpdateObjectiveProgress(questID string, objectiveID string, progress int64) error {
	quest, exists := qm.activeQuests[questID]
	if !exists {
		return ErrQuestNotFound
	}

	// 查找目标
	var objective *QuestObjective
	for _, obj := range quest.Objectives {
		if obj.ID == objectiveID {
			objective = obj
			break
		}
	}

	if objective == nil {
		return ErrObjectiveNotFound
	}

	if objective.Completed {
		return ErrObjectiveAlreadyCompleted
	}

	// 更新进度
	objective.Current += progress
	if objective.Current >= objective.Required {
		objective.Current = objective.Required
		objective.Completed = true

		// 发布目标完成事件
		qm.addEvent(ObjectiveCompletedEvent{
			Player:         qm.playerID,
			QuestID:        questID,
			ObjectiveID:    objectiveID,
			ObjectiveName:  objective.Description,
			OccurredAtTime: time.Now(),
		})
	}

	quest.UpdatedAt = time.Now()
	qm.lastUpdate = time.Now()

	// 检查任务是否完成
	if qm.checkQuestCompletion(quest) {
		return qm.CompleteQuest(questID)
	}

	return nil
}

// CompleteQuest 完成任务
func (qm *QuestManager) CompleteQuest(questID string) error {
	quest, exists := qm.activeQuests[questID]
	if !exists {
		return ErrQuestNotFound
	}

	if quest.Status != QuestStatusAccepted && quest.Status != QuestStatusInProgress {
		return ErrQuestNotActive
	}

	// 检查所有必需目标是否完成
	if !qm.checkQuestCompletion(quest) {
		return ErrQuestNotCompleted
	}

	// 完成任务
	quest.Status = QuestStatusCompleted
	now := time.Now()
	quest.CompletedTime = &now
	quest.UpdatedAt = time.Now()

	// 移动到已完成任务
	delete(qm.activeQuests, questID)
	qm.completedQuests[questID] = quest

	// 处理重复任务
	if quest.RepeatType != RepeatTypeNone {
		quest.RepeatCount++
		if quest.MaxRepeats == 0 || quest.RepeatCount < quest.MaxRepeats {
			// 重置任务状态以便重复
			qm.resetQuestForRepeat(quest)
		}
	}

	qm.lastUpdate = time.Now()

	// 发布事件
	qm.addEvent(QuestCompletedEvent{
		Player:         qm.playerID,
		QuestID:        questID,
		QuestName:      quest.Name,
		Rewards:        quest.Rewards,
		OccurredAtTime: time.Now(),
	})

	return nil
}

// AbandonQuest 放弃任务
func (qm *QuestManager) AbandonQuest(questID string) error {
	quest, exists := qm.activeQuests[questID]
	if !exists {
		return ErrQuestNotFound
	}

	// 检查是否可以放弃
	if quest.QuestType == QuestTypeMain {
		return ErrCannotAbandonMainQuest
	}

	quest.Status = QuestStatusAbandoned
	quest.UpdatedAt = time.Now()
	delete(qm.activeQuests, questID)
	qm.lastUpdate = time.Now()

	return nil
}

// UnlockAchievement 解锁成就
func (qm *QuestManager) UnlockAchievement(achievementID string, achievement *Achievement) error {
	if _, exists := qm.achievements[achievementID]; exists {
		return ErrAchievementAlreadyUnlocked
	}

	// 检查成就要求
	if !qm.checkAchievementRequirements(achievement) {
		return ErrAchievementRequirementsNotMet
	}

	achievement.Unlocked = true
	now := time.Now()
	achievement.UnlockedAt = &now
	qm.achievements[achievementID] = achievement
	qm.lastUpdate = time.Now()

	// 发布事件
	qm.addEvent(AchievementUnlockedEvent{
		Player:          qm.playerID,
		AchievementID:   achievementID,
		AchievementName: achievement.Name,
		Points:          achievement.Points,
		OccurredAtTime:  time.Now(),
	})

	return nil
}

// checkPrerequisites 检查前置条件
func (qm *QuestManager) checkPrerequisites(prerequisites []string) bool {
	for _, prereq := range prerequisites {
		if _, exists := qm.completedQuests[prereq]; !exists {
			return false
		}
	}
	return true
}

// checkQuestCompletion 检查任务完成条件
func (qm *QuestManager) checkQuestCompletion(quest *Quest) bool {
	for _, objective := range quest.Objectives {
		if !objective.Optional && !objective.Completed {
			return false
		}
	}
	return true
}

// checkAchievementRequirements 检查成就要求
func (qm *QuestManager) checkAchievementRequirements(achievement *Achievement) bool {
	for _, req := range achievement.Requirements {
		if !req.Completed {
			return false
		}
	}
	return true
}

// resetQuestForRepeat 重置任务以便重复
func (qm *QuestManager) resetQuestForRepeat(quest *Quest) {
	quest.Status = QuestStatusAvailable
	quest.StartTime = nil
	quest.CompletedTime = nil
	quest.ExpireTime = nil

	// 重置所有目标
	for _, objective := range quest.Objectives {
		objective.Current = 0
		objective.Completed = false
	}
}

// addEvent 添加领域事件
func (qm *QuestManager) addEvent(event DomainEvent) {
	qm.events = append(qm.events, event)
}

// GetEvents 获取领域事件
func (qm *QuestManager) GetEvents() []DomainEvent {
	return qm.events
}

// ClearEvents 清除领域事件
func (qm *QuestManager) ClearEvents() {
	qm.events = make([]DomainEvent, 0)
}
