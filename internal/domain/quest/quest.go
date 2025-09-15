package quest

import (
	"errors"
	"time"
)

// QuestManager 任务管理器聚合根
type QuestManager struct {
	playerID      string
	activeQuests  map[string]*Quest
	completedQuests map[string]*Quest
	dailyQuests   map[string]*Quest
	weeklyQuests  map[string]*Quest
	achievements  map[string]*Achievement
	lastUpdate    time.Time
	events        []DomainEvent
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
	id            string
	name          string
	description   string
	questType     QuestType
	category      QuestCategory
	status        QuestStatus
	priority      QuestPriority
	objectives    []*QuestObjective
	rewards       []*QuestReward
	prerequisites []string // 前置任务ID
	startTime     *time.Time
	expireTime    *time.Time
	completedTime *time.Time
	timeLimit     *time.Duration
	repeatType    RepeatType
	repeatCount   int
	maxRepeats    int
	level         int
	minLevel      int
	maxLevel      int
	classRestrictions []string
	raceRestrictions  []string
	createdAt     time.Time
	updatedAt     time.Time
}

// NewQuest 创建新任务
func NewQuest(id, name string, questType QuestType) *Quest {
	return &Quest{
		id:          id,
		name:        name,
		questType:   questType,
		status:      QuestStatusAvailable,
		priority:    QuestPriorityNormal,
		objectives:  make([]*QuestObjective, 0),
		rewards:     make([]*QuestReward, 0),
		repeatType:  RepeatTypeNone,
		repeatCount: 0,
		maxRepeats:  1,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
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
	id          string
	description string
	objectiveType ObjectiveType
	target      string // 目标ID或名称
	current     int64
	required    int64
	completed   bool
	optional    bool
	order       int
	metadata    map[string]interface{}
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
	rewardType RewardType
	rewardID   string
	quantity   int64
	optional   bool
	condition  *RewardCondition
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
	conditionType ConditionType
	value         interface{}
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
	id            string
	name          string
	description   string
	category      AchievementCategory
	points        int64
	requirements  []*AchievementRequirement
	rewards       []*QuestReward
	unlocked      bool
	progress      int64
	totalProgress int64
	unlockedAt    *time.Time
	hidden        bool
	rare          bool
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
	requirementType RequirementType
	target          string
	value           int64
	current         int64
	completed       bool
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
	playerID   string
	questID    string
	questName  string
	occurredAt time.Time
}

func (e QuestAcceptedEvent) EventType() string   { return "quest.accepted" }
func (e QuestAcceptedEvent) OccurredAt() time.Time { return e.occurredAt }
func (e QuestAcceptedEvent) PlayerID() string    { return e.playerID }

// QuestCompletedEvent 任务完成事件
type QuestCompletedEvent struct {
	playerID   string
	questID    string
	questName  string
	rewards    []*QuestReward
	occurredAt time.Time
}

func (e QuestCompletedEvent) EventType() string   { return "quest.completed" }
func (e QuestCompletedEvent) OccurredAt() time.Time { return e.occurredAt }
func (e QuestCompletedEvent) PlayerID() string    { return e.playerID }

// QuestFailedEvent 任务失败事件
type QuestFailedEvent struct {
	playerID   string
	questID    string
	questName  string
	reason     string
	occurredAt time.Time
}

func (e QuestFailedEvent) EventType() string   { return "quest.failed" }
func (e QuestFailedEvent) OccurredAt() time.Time { return e.occurredAt }
func (e QuestFailedEvent) PlayerID() string    { return e.playerID }

// ObjectiveCompletedEvent 目标完成事件
type ObjectiveCompletedEvent struct {
	playerID      string
	questID       string
	objectiveID   string
	objectiveName string
	occurredAt    time.Time
}

func (e ObjectiveCompletedEvent) EventType() string   { return "objective.completed" }
func (e ObjectiveCompletedEvent) OccurredAt() time.Time { return e.occurredAt }
func (e ObjectiveCompletedEvent) PlayerID() string    { return e.playerID }

// AchievementUnlockedEvent 成就解锁事件
type AchievementUnlockedEvent struct {
	playerID        string
	achievementID   string
	achievementName string
	points          int64
	occurredAt      time.Time
}

func (e AchievementUnlockedEvent) EventType() string   { return "achievement.unlocked" }
func (e AchievementUnlockedEvent) OccurredAt() time.Time { return e.occurredAt }
func (e AchievementUnlockedEvent) PlayerID() string    { return e.playerID }

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
	if quest.status != QuestStatusAvailable {
		return ErrQuestNotAvailable
	}

	// 检查是否已接受
	if _, exists := qm.activeQuests[quest.id]; exists {
		return ErrQuestAlreadyAccepted
	}

	// 检查前置条件
	if !qm.checkPrerequisites(quest.prerequisites) {
		return ErrPrerequisitesNotMet
	}

	// 检查等级限制
	if quest.minLevel > 0 || quest.maxLevel > 0 {
		// 这里需要获取玩家等级，暂时跳过
	}

	// 接受任务
	quest.status = QuestStatusAccepted
	now := time.Now()
	quest.startTime = &now
	quest.updatedAt = time.Now()

	// 设置过期时间
	if quest.timeLimit != nil {
		expireTime := now.Add(*quest.timeLimit)
		quest.expireTime = &expireTime
	}

	qm.activeQuests[quest.id] = quest
	qm.lastUpdate = time.Now()

	// 发布事件
	qm.addEvent(QuestAcceptedEvent{
		playerID:   qm.playerID,
		questID:    quest.id,
		questName:  quest.name,
		occurredAt: time.Now(),
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
	for _, obj := range quest.objectives {
		if obj.id == objectiveID {
			objective = obj
			break
		}
	}

	if objective == nil {
		return ErrObjectiveNotFound
	}

	if objective.completed {
		return ErrObjectiveAlreadyCompleted
	}

	// 更新进度
	objective.current += progress
	if objective.current >= objective.required {
		objective.current = objective.required
		objective.completed = true

		// 发布目标完成事件
		qm.addEvent(ObjectiveCompletedEvent{
			playerID:      qm.playerID,
			questID:       questID,
			objectiveID:   objectiveID,
			objectiveName: objective.description,
			occurredAt:    time.Now(),
		})
	}

	quest.updatedAt = time.Now()
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

	if quest.status != QuestStatusAccepted && quest.status != QuestStatusInProgress {
		return ErrQuestNotActive
	}

	// 检查所有必需目标是否完成
	if !qm.checkQuestCompletion(quest) {
		return ErrQuestNotCompleted
	}

	// 完成任务
	quest.status = QuestStatusCompleted
	now := time.Now()
	quest.completedTime = &now
	quest.updatedAt = time.Now()

	// 移动到已完成任务
	delete(qm.activeQuests, questID)
	qm.completedQuests[questID] = quest

	// 处理重复任务
	if quest.repeatType != RepeatTypeNone {
		quest.repeatCount++
		if quest.maxRepeats == 0 || quest.repeatCount < quest.maxRepeats {
			// 重置任务状态以便重复
			qm.resetQuestForRepeat(quest)
		}
	}

	qm.lastUpdate = time.Now()

	// 发布事件
	qm.addEvent(QuestCompletedEvent{
		playerID:   qm.playerID,
		questID:    questID,
		questName:  quest.name,
		rewards:    quest.rewards,
		occurredAt: time.Now(),
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
	if quest.questType == QuestTypeMain {
		return ErrCannotAbandonMainQuest
	}

	quest.status = QuestStatusAbandoned
	quest.updatedAt = time.Now()
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

	achievement.unlocked = true
	now := time.Now()
	achievement.unlockedAt = &now
	qm.achievements[achievementID] = achievement
	qm.lastUpdate = time.Now()

	// 发布事件
	qm.addEvent(AchievementUnlockedEvent{
		playerID:        qm.playerID,
		achievementID:   achievementID,
		achievementName: achievement.name,
		points:          achievement.points,
		occurredAt:      time.Now(),
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
	for _, objective := range quest.objectives {
		if !objective.optional && !objective.completed {
			return false
		}
	}
	return true
}

// checkAchievementRequirements 检查成就要求
func (qm *QuestManager) checkAchievementRequirements(achievement *Achievement) bool {
	for _, req := range achievement.requirements {
		if !req.completed {
			return false
		}
	}
	return true
}

// resetQuestForRepeat 重置任务以便重复
func (qm *QuestManager) resetQuestForRepeat(quest *Quest) {
	quest.status = QuestStatusAvailable
	quest.startTime = nil
	quest.completedTime = nil
	quest.expireTime = nil

	// 重置所有目标
	for _, objective := range quest.objectives {
		objective.current = 0
		objective.completed = false
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