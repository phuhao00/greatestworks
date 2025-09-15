package beginner

import "errors"

// 新手系统相关错误
var (
	ErrInvalidGuide             = errors.New("invalid guide")
	ErrGuideNotFound            = errors.New("guide not found")
	ErrGuideAlreadyCompleted    = errors.New("guide already completed")
	ErrInvalidStep              = errors.New("invalid step")
	ErrStepNotFound             = errors.New("step not found")
	ErrStepAlreadyCompleted     = errors.New("step already completed")
	ErrConditionNotMet          = errors.New("step condition not met")
	ErrInvalidTutorial          = errors.New("invalid tutorial")
	ErrTutorialNotFound         = errors.New("tutorial not found")
	ErrTutorialAlreadyCompleted = errors.New("tutorial already completed")
	ErrInvalidReward            = errors.New("invalid reward")
	ErrRewardNotFound           = errors.New("reward not found")
	ErrRewardAlreadyClaimed     = errors.New("reward already claimed")
	ErrBeginnerAlreadyCompleted = errors.New("beginner guide already completed")
	ErrInvalidCondition         = errors.New("invalid step condition")
	ErrInvalidRewardType        = errors.New("invalid reward type")
	ErrInsufficientLevel        = errors.New("insufficient level for guide")
	ErrPrerequisiteNotMet       = errors.New("prerequisite guide not completed")
)