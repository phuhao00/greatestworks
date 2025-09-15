package beginner

import (
	"time"
	"github.com/google/uuid"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	GetEventID() string
	GetEventType() string
	GetAggregateID() string
	GetOccurredAt() time.Time
	GetEventData() interface{}
}

// BaseDomainEvent 基础领域事件
type BaseDomainEvent struct {
	EventID     string    `json:"event_id"`
	EventType   string    `json:"event_type"`
	AggregateID string    `json:"aggregate_id"`
	OccurredAt  time.Time `json:"occurred_at"`
}

// GetEventID 获取事件ID
func (e *BaseDomainEvent) GetEventID() string {
	return e.EventID
}

// GetEventType 获取事件类型
func (e *BaseDomainEvent) GetEventType() string {
	return e.EventType
}

// GetAggregateID 获取聚合根ID
func (e *BaseDomainEvent) GetAggregateID() string {
	return e.AggregateID
}

// GetOccurredAt 获取发生时间
func (e *BaseDomainEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// BeginnerGuideStartedEvent 新手引导开始事件
type BeginnerGuideStartedEvent struct {
	BaseDomainEvent
	PlayerID string `json:"player_id"`
	GuideID  string `json:"guide_id"`
}

// NewBeginnerGuideStartedEvent 创建新手引导开始事件
func NewBeginnerGuideStartedEvent(playerID, guideID string) *BeginnerGuideStartedEvent {
	return &BeginnerGuideStartedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "BeginnerGuideStarted",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID: playerID,
		GuideID:  guideID,
	}
}

// GetEventData 获取事件数据
func (e *BeginnerGuideStartedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id": e.PlayerID,
		"guide_id":  e.GuideID,
	}
}

// GuideStepCompletedEvent 引导步骤完成事件
type GuideStepCompletedEvent struct {
	BaseDomainEvent
	PlayerID string     `json:"player_id"`
	GuideID  string     `json:"guide_id"`
	StepID   int        `json:"step_id"`
	Step     *GuideStep `json:"step"`
	Reward   *BeginnerReward `json:"reward"`
}

// NewGuideStepCompletedEvent 创建引导步骤完成事件
func NewGuideStepCompletedEvent(playerID, guideID string, stepID int, step *GuideStep, reward *BeginnerReward) *GuideStepCompletedEvent {
	return &GuideStepCompletedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "GuideStepCompleted",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID: playerID,
		GuideID:  guideID,
		StepID:   stepID,
		Step:     step,
		Reward:   reward,
	}
}

// GetEventData 获取事件数据
func (e *GuideStepCompletedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id": e.PlayerID,
		"guide_id":  e.GuideID,
		"step_id":   e.StepID,
		"step":      e.Step,
		"reward":    e.Reward,
	}
}

// BeginnerGuideCompletedEvent 新手引导完成事件
type BeginnerGuideCompletedEvent struct {
	BaseDomainEvent
	PlayerID string          `json:"player_id"`
	GuideID  string          `json:"guide_id"`
	Progress *GuideProgress  `json:"progress"`
	Reward   *BeginnerReward `json:"reward"`
}

// NewBeginnerGuideCompletedEvent 创建新手引导完成事件
func NewBeginnerGuideCompletedEvent(playerID, guideID string, progress *GuideProgress, reward *BeginnerReward) *BeginnerGuideCompletedEvent {
	return &BeginnerGuideCompletedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "BeginnerGuideCompleted",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID: playerID,
		GuideID:  guideID,
		Progress: progress,
		Reward:   reward,
	}
}

// GetEventData 获取事件数据
func (e *BeginnerGuideCompletedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id": e.PlayerID,
		"guide_id":  e.GuideID,
		"progress":  e.Progress,
		"reward":    e.Reward,
	}
}

// TutorialStartedEvent 教程开始事件
type TutorialStartedEvent struct {
	BaseDomainEvent
	PlayerID   string   `json:"player_id"`
	TutorialID string   `json:"tutorial_id"`
	Tutorial   *Tutorial `json:"tutorial"`
}

// NewTutorialStartedEvent 创建教程开始事件
func NewTutorialStartedEvent(playerID, tutorialID string, tutorial *Tutorial) *TutorialStartedEvent {
	return &TutorialStartedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "TutorialStarted",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:   playerID,
		TutorialID: tutorialID,
		Tutorial:   tutorial,
	}
}

// GetEventData 获取事件数据
func (e *TutorialStartedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id":   e.PlayerID,
		"tutorial_id": e.TutorialID,
		"tutorial":    e.Tutorial,
	}
}

// TutorialCompletedEvent 教程完成事件
type TutorialCompletedEvent struct {
	BaseDomainEvent
	PlayerID   string          `json:"player_id"`
	TutorialID string          `json:"tutorial_id"`
	Tutorial   *Tutorial       `json:"tutorial"`
	Reward     *BeginnerReward `json:"reward"`
}

// NewTutorialCompletedEvent 创建教程完成事件
func NewTutorialCompletedEvent(playerID, tutorialID string, tutorial *Tutorial, reward *BeginnerReward) *TutorialCompletedEvent {
	return &TutorialCompletedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "TutorialCompleted",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:   playerID,
		TutorialID: tutorialID,
		Tutorial:   tutorial,
		Reward:     reward,
	}
}

// GetEventData 获取事件数据
func (e *TutorialCompletedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id":   e.PlayerID,
		"tutorial_id": e.TutorialID,
		"tutorial":    e.Tutorial,
		"reward":      e.Reward,
	}
}

// BeginnerRewardClaimedEvent 新手奖励领取事件
type BeginnerRewardClaimedEvent struct {
	BaseDomainEvent
	PlayerID string          `json:"player_id"`
	RewardID string          `json:"reward_id"`
	Reward   *BeginnerReward `json:"reward"`
	Source   string          `json:"source"` // guide, tutorial, final
}

// NewBeginnerRewardClaimedEvent 创建新手奖励领取事件
func NewBeginnerRewardClaimedEvent(playerID, rewardID string, reward *BeginnerReward, source string) *BeginnerRewardClaimedEvent {
	return &BeginnerRewardClaimedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "BeginnerRewardClaimed",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID: playerID,
		RewardID: rewardID,
		Reward:   reward,
		Source:   source,
	}
}

// GetEventData 获取事件数据
func (e *BeginnerRewardClaimedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id": e.PlayerID,
		"reward_id": e.RewardID,
		"reward":    e.Reward,
		"source":    e.Source,
	}
}

// AllBeginnerGuidesCompletedEvent 所有新手引导完成事件
type AllBeginnerGuidesCompletedEvent struct {
	BaseDomainEvent
	PlayerID        string          `json:"player_id"`
	CompletedGuides []string        `json:"completed_guides"`
	FinalReward     *BeginnerReward `json:"final_reward"`
	CompletionTime  time.Duration   `json:"completion_time"`
}

// NewAllBeginnerGuidesCompletedEvent 创建所有新手引导完成事件
func NewAllBeginnerGuidesCompletedEvent(playerID string, completedGuides []string, finalReward *BeginnerReward, completionTime time.Duration) *AllBeginnerGuidesCompletedEvent {
	return &AllBeginnerGuidesCompletedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "AllBeginnerGuidesCompleted",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:        playerID,
		CompletedGuides: completedGuides,
		FinalReward:     finalReward,
		CompletionTime:  completionTime,
	}
}

// GetEventData 获取事件数据
func (e *AllBeginnerGuidesCompletedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id":        e.PlayerID,
		"completed_guides": e.CompletedGuides,
		"final_reward":     e.FinalReward,
		"completion_time":  e.CompletionTime,
	}
}

// BeginnerProgressUpdatedEvent 新手进度更新事件
type BeginnerProgressUpdatedEvent struct {
	BaseDomainEvent
	PlayerID      string         `json:"player_id"`
	GuideID       string         `json:"guide_id"`
	CurrentStep   int            `json:"current_step"`
	Progress      *GuideProgress `json:"progress"`
	NextStepHint  string         `json:"next_step_hint"`
}

// NewBeginnerProgressUpdatedEvent 创建新手进度更新事件
func NewBeginnerProgressUpdatedEvent(playerID, guideID string, currentStep int, progress *GuideProgress, nextStepHint string) *BeginnerProgressUpdatedEvent {
	return &BeginnerProgressUpdatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "BeginnerProgressUpdated",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:     playerID,
		GuideID:      guideID,
		CurrentStep:  currentStep,
		Progress:     progress,
		NextStepHint: nextStepHint,
	}
}

// GetEventData 获取事件数据
func (e *BeginnerProgressUpdatedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id":      e.PlayerID,
		"guide_id":       e.GuideID,
		"current_step":   e.CurrentStep,
		"progress":       e.Progress,
		"next_step_hint": e.NextStepHint,
	}
}

// BeginnerSkippedEvent 新手引导跳过事件
type BeginnerSkippedEvent struct {
	BaseDomainEvent
	PlayerID string `json:"player_id"`
	GuideID  string `json:"guide_id"`
	StepID   int    `json:"step_id"`
	Reason   string `json:"reason"`
}

// NewBeginnerSkippedEvent 创建新手引导跳过事件
func NewBeginnerSkippedEvent(playerID, guideID string, stepID int, reason string) *BeginnerSkippedEvent {
	return &BeginnerSkippedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "BeginnerSkipped",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID: playerID,
		GuideID:  guideID,
		StepID:   stepID,
		Reason:   reason,
	}
}

// GetEventData 获取事件数据
func (e *BeginnerSkippedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id": e.PlayerID,
		"guide_id":  e.GuideID,
		"step_id":   e.StepID,
		"reason":    e.Reason,
	}
}

// BeginnerHelpRequestedEvent 新手帮助请求事件
type BeginnerHelpRequestedEvent struct {
	BaseDomainEvent
	PlayerID    string `json:"player_id"`
	GuideID     string `json:"guide_id"`
	StepID      int    `json:"step_id"`
	HelpType    string `json:"help_type"` // hint, tutorial, support
	Description string `json:"description"`
}

// NewBeginnerHelpRequestedEvent 创建新手帮助请求事件
func NewBeginnerHelpRequestedEvent(playerID, guideID string, stepID int, helpType, description string) *BeginnerHelpRequestedEvent {
	return &BeginnerHelpRequestedEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:     uuid.New().String(),
			EventType:   "BeginnerHelpRequested",
			AggregateID: playerID,
			OccurredAt:  time.Now(),
		},
		PlayerID:    playerID,
		GuideID:     guideID,
		StepID:      stepID,
		HelpType:    helpType,
		Description: description,
	}
}

// GetEventData 获取事件数据
func (e *BeginnerHelpRequestedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"player_id":   e.PlayerID,
		"guide_id":    e.GuideID,
		"step_id":     e.StepID,
		"help_type":   e.HelpType,
		"description": e.Description,
	}
}