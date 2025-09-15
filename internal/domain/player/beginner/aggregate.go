package beginner

import (
	"time"
	"github.com/google/uuid"
)

// BeginnerAggregate 新手聚合根
type BeginnerAggregate struct {
	playerID       string
	guideSteps     map[string]*GuideStep
	tutorials      map[string]*Tutorial
	currentGuide   string
	currentStep    int
	completedGuides []string
	rewards        []*BeginnerReward
	isCompleted   bool
	startedAt      time.Time
	completedAt    *time.Time
	updatedAt      time.Time
	version        int
}

// NewBeginnerAggregate 创建新手聚合根
func NewBeginnerAggregate(playerID string) *BeginnerAggregate {
	return &BeginnerAggregate{
		playerID:        playerID,
		guideSteps:      make(map[string]*GuideStep),
		tutorials:       make(map[string]*Tutorial),
		currentGuide:    "main_guide", // 默认主引导
		currentStep:     1,
		completedGuides: make([]string, 0),
		rewards:         make([]*BeginnerReward, 0),
		isCompleted:     false,
		startedAt:       time.Now(),
		updatedAt:       time.Now(),
		version:         1,
	}
}

// GetPlayerID 获取玩家ID
func (b *BeginnerAggregate) GetPlayerID() string {
	return b.playerID
}

// StartGuide 开始引导
func (b *BeginnerAggregate) StartGuide(guideID string) error {
	if b.isCompleted {
		return ErrBeginnerAlreadyCompleted
	}
	
	// 检查是否已完成该引导
	for _, completed := range b.completedGuides {
		if completed == guideID {
			return ErrGuideAlreadyCompleted
		}
	}
	
	b.currentGuide = guideID
	b.currentStep = 1
	b.updateVersion()
	return nil
}

// CompleteStep 完成步骤
func (b *BeginnerAggregate) CompleteStep(guideID string, stepID int) error {
	if b.isCompleted {
		return ErrBeginnerAlreadyCompleted
	}
	
	if b.currentGuide != guideID {
		return ErrInvalidGuide
	}
	
	if b.currentStep != stepID {
		return ErrInvalidStep
	}
	
	stepKey := b.getStepKey(guideID, stepID)
	step, exists := b.guideSteps[stepKey]
	if !exists {
		return ErrStepNotFound
	}
	
	// 标记步骤完成
	step.Complete()
	
	// 给予奖励
	if step.HasReward() {
		reward := step.GetReward()
		b.rewards = append(b.rewards, reward)
	}
	
	// 检查是否完成整个引导
	if b.isGuideCompleted(guideID) {
		b.completeGuide(guideID)
	} else {
		b.currentStep++
	}
	
	b.updateVersion()
	return nil
}

// AddGuideStep 添加引导步骤
func (b *BeginnerAggregate) AddGuideStep(guideID string, step *GuideStep) error {
	if step == nil {
		return ErrInvalidStep
	}
	
	stepKey := b.getStepKey(guideID, step.GetStepID())
	b.guideSteps[stepKey] = step
	b.updateVersion()
	return nil
}

// GetCurrentStep 获取当前步骤
func (b *BeginnerAggregate) GetCurrentStep() *GuideStep {
	stepKey := b.getStepKey(b.currentGuide, b.currentStep)
	return b.guideSteps[stepKey]
}

// GetGuideProgress 获取引导进度
func (b *BeginnerAggregate) GetGuideProgress(guideID string) *GuideProgress {
	totalSteps := b.countGuideSteps(guideID)
	completedSteps := b.countCompletedSteps(guideID)
	
	return &GuideProgress{
		GuideID:        guideID,
		TotalSteps:     totalSteps,
		CompletedSteps: completedSteps,
		IsCompleted:    b.isGuideCompleted(guideID),
		Progress:       float64(completedSteps) / float64(totalSteps),
	}
}

// AddTutorial 添加教程
func (b *BeginnerAggregate) AddTutorial(tutorial *Tutorial) error {
	if tutorial == nil {
		return ErrInvalidTutorial
	}
	
	b.tutorials[tutorial.GetID()] = tutorial
	b.updateVersion()
	return nil
}

// CompleteTutorial 完成教程
func (b *BeginnerAggregate) CompleteTutorial(tutorialID string) error {
	tutorial, exists := b.tutorials[tutorialID]
	if !exists {
		return ErrTutorialNotFound
	}
	
	tutorial.Complete()
	
	// 给予教程奖励
	if tutorial.HasReward() {
		reward := tutorial.GetReward()
		b.rewards = append(b.rewards, reward)
	}
	
	b.updateVersion()
	return nil
}

// GetTutorial 获取教程
func (b *BeginnerAggregate) GetTutorial(tutorialID string) *Tutorial {
	return b.tutorials[tutorialID]
}

// GetAllTutorials 获取所有教程
func (b *BeginnerAggregate) GetAllTutorials() map[string]*Tutorial {
	return b.tutorials
}

// GetUnclaimedRewards 获取未领取的奖励
func (b *BeginnerAggregate) GetUnclaimedRewards() []*BeginnerReward {
	var unclaimed []*BeginnerReward
	for _, reward := range b.rewards {
		if !reward.IsClaimed() {
			unclaimed = append(unclaimed, reward)
		}
	}
	return unclaimed
}

// ClaimReward 领取奖励
func (b *BeginnerAggregate) ClaimReward(rewardID string) error {
	for _, reward := range b.rewards {
		if reward.GetID() == rewardID {
			if reward.IsClaimed() {
				return ErrRewardAlreadyClaimed
			}
			reward.Claim()
			b.updateVersion()
			return nil
		}
	}
	return ErrRewardNotFound
}

// IsCompleted 检查是否完成新手引导
func (b *BeginnerAggregate) IsCompleted() bool {
	return b.isCompleted
}

// GetCurrentGuide 获取当前引导
func (b *BeginnerAggregate) GetCurrentGuide() string {
	return b.currentGuide
}

// GetCurrentStepID 获取当前步骤ID
func (b *BeginnerAggregate) GetCurrentStepID() int {
	return b.currentStep
}

// GetCompletedGuides 获取已完成的引导
func (b *BeginnerAggregate) GetCompletedGuides() []string {
	return b.completedGuides
}

// GetStartedAt 获取开始时间
func (b *BeginnerAggregate) GetStartedAt() time.Time {
	return b.startedAt
}

// GetCompletedAt 获取完成时间
func (b *BeginnerAggregate) GetCompletedAt() *time.Time {
	return b.completedAt
}

// GetVersion 获取版本
func (b *BeginnerAggregate) GetVersion() int {
	return b.version
}

// GetUpdatedAt 获取更新时间
func (b *BeginnerAggregate) GetUpdatedAt() time.Time {
	return b.updatedAt
}

// 私有方法

// getStepKey 获取步骤键
func (b *BeginnerAggregate) getStepKey(guideID string, stepID int) string {
	return fmt.Sprintf("%s_%d", guideID, stepID)
}

// isGuideCompleted 检查引导是否完成
func (b *BeginnerAggregate) isGuideCompleted(guideID string) bool {
	totalSteps := b.countGuideSteps(guideID)
	completedSteps := b.countCompletedSteps(guideID)
	return completedSteps >= totalSteps
}

// countGuideSteps 统计引导步骤数
func (b *BeginnerAggregate) countGuideSteps(guideID string) int {
	count := 0
	for stepKey := range b.guideSteps {
		if strings.HasPrefix(stepKey, guideID+"_") {
			count++
		}
	}
	return count
}

// countCompletedSteps 统计已完成步骤数
func (b *BeginnerAggregate) countCompletedSteps(guideID string) int {
	count := 0
	for stepKey, step := range b.guideSteps {
		if strings.HasPrefix(stepKey, guideID+"_") && step.IsCompleted() {
			count++
		}
	}
	return count
}

// completeGuide 完成引导
func (b *BeginnerAggregate) completeGuide(guideID string) {
	b.completedGuides = append(b.completedGuides, guideID)
	
	// 检查是否完成所有必要的引导
	if b.isAllRequiredGuidesCompleted() {
		b.isCompleted = true
		now := time.Now()
		b.completedAt = &now
		
		// 给予完成新手引导的最终奖励
		finalReward := NewBeginnerReward("final_reward", RewardTypeMultiple, map[string]interface{}{
			"gold":       1000,
			"experience": 500,
			"items":      []string{"beginner_sword", "beginner_armor"},
		}, "完成新手引导奖励")
		b.rewards = append(b.rewards, finalReward)
	}
}

// isAllRequiredGuidesCompleted 检查是否完成所有必要引导
func (b *BeginnerAggregate) isAllRequiredGuidesCompleted() bool {
	requiredGuides := []string{"main_guide", "combat_guide", "inventory_guide"}
	
	for _, required := range requiredGuides {
		found := false
		for _, completed := range b.completedGuides {
			if completed == required {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	return true
}

// updateVersion 更新版本
func (b *BeginnerAggregate) updateVersion() {
	b.version++
	b.updatedAt = time.Now()
}