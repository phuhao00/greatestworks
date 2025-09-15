package quest

import "errors"

var (
	// 任务管理器相关错误
	ErrQuestManagerNotFound    = errors.New("quest manager not found")
	ErrInvalidQuestManager     = errors.New("invalid quest manager")
	
	// 任务相关错误
	ErrQuestNotFound           = errors.New("quest not found")
	ErrQuestNotAvailable       = errors.New("quest is not available")
	ErrQuestAlreadyAccepted    = errors.New("quest already accepted")
	ErrQuestNotActive          = errors.New("quest is not active")
	ErrQuestNotCompleted       = errors.New("quest is not completed")
	ErrQuestExpired            = errors.New("quest has expired")
	ErrQuestFailed             = errors.New("quest has failed")
	ErrCannotAbandonMainQuest  = errors.New("cannot abandon main quest")
	ErrQuestAlreadyCompleted   = errors.New("quest already completed")
	ErrInvalidQuestType        = errors.New("invalid quest type")
	ErrInvalidQuestStatus      = errors.New("invalid quest status")
	
	// 任务目标相关错误
	ErrObjectiveNotFound       = errors.New("objective not found")
	ErrObjectiveAlreadyCompleted = errors.New("objective already completed")
	ErrInvalidObjectiveType    = errors.New("invalid objective type")
	ErrInvalidProgress         = errors.New("invalid progress value")
	ErrProgressExceedsRequired = errors.New("progress exceeds required amount")
	
	// 前置条件相关错误
	ErrPrerequisitesNotMet     = errors.New("prerequisites not met")
	ErrLevelRequirementNotMet  = errors.New("level requirement not met")
	ErrClassRestriction        = errors.New("class restriction for quest")
	ErrRaceRestriction         = errors.New("race restriction for quest")
	ErrGuildRequirement        = errors.New("guild requirement not met")
	
	// 奖励相关错误
	ErrInvalidReward           = errors.New("invalid quest reward")
	ErrRewardNotFound          = errors.New("quest reward not found")
	ErrRewardAlreadyClaimed    = errors.New("reward already claimed")
	ErrInsufficientSpace       = errors.New("insufficient inventory space for reward")
	
	// 成就相关错误
	ErrAchievementNotFound     = errors.New("achievement not found")
	ErrAchievementAlreadyUnlocked = errors.New("achievement already unlocked")
	ErrAchievementRequirementsNotMet = errors.New("achievement requirements not met")
	ErrInvalidAchievementType  = errors.New("invalid achievement type")
	ErrAchievementLocked       = errors.New("achievement is locked")
	
	// 日常/周常任务相关错误
	ErrDailyQuestLimitReached  = errors.New("daily quest limit reached")
	ErrWeeklyQuestLimitReached = errors.New("weekly quest limit reached")
	ErrQuestCooldownActive     = errors.New("quest cooldown is active")
	ErrRepeatLimitReached      = errors.New("quest repeat limit reached")
	
	// 任务链相关错误
	ErrQuestChainBroken        = errors.New("quest chain is broken")
	ErrInvalidQuestOrder       = errors.New("invalid quest order in chain")
	
	// 配置相关错误
	ErrInvalidQuestConfig      = errors.New("invalid quest configuration")
	ErrQuestConfigNotFound     = errors.New("quest configuration not found")
	ErrInvalidObjectiveConfig  = errors.New("invalid objective configuration")
	ErrInvalidRewardConfig     = errors.New("invalid reward configuration")
)