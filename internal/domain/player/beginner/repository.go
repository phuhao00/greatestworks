package beginner

import (
	"context"
	"time"
)

// BeginnerRepository 新手仓储接口
type BeginnerRepository interface {
	// SaveBeginnerAggregate 保存新手聚合根
	SaveBeginnerAggregate(ctx context.Context, aggregate *BeginnerAggregate) error

	// GetBeginnerAggregate 获取新手聚合根
	GetBeginnerAggregate(ctx context.Context, playerID string) (*BeginnerAggregate, error)

	// DeleteBeginnerAggregate 删除新手聚合根
	DeleteBeginnerAggregate(ctx context.Context, playerID string) error

	// SaveGuideStep 保存引导步骤
	SaveGuideStep(ctx context.Context, playerID string, step *GuideStep) error

	// GetGuideStep 获取引导步骤
	GetGuideStep(ctx context.Context, playerID, guideID string, stepID int) (*GuideStep, error)

	// GetGuideSteps 获取引导的所有步骤
	GetGuideSteps(ctx context.Context, playerID, guideID string) ([]*GuideStep, error)

	// UpdateGuideStep 更新引导步骤
	UpdateGuideStep(ctx context.Context, playerID string, step *GuideStep) error

	// SaveTutorial 保存教程
	SaveTutorial(ctx context.Context, playerID string, tutorial *Tutorial) error

	// GetTutorial 获取教程
	GetTutorial(ctx context.Context, playerID, tutorialID string) (*Tutorial, error)

	// GetPlayerTutorials 获取玩家所有教程
	GetPlayerTutorials(ctx context.Context, playerID string) ([]*Tutorial, error)

	// UpdateTutorial 更新教程
	UpdateTutorial(ctx context.Context, playerID string, tutorial *Tutorial) error

	// SaveBeginnerReward 保存新手奖励
	SaveBeginnerReward(ctx context.Context, playerID string, reward *BeginnerReward) error

	// GetBeginnerReward 获取新手奖励
	GetBeginnerReward(ctx context.Context, playerID, rewardID string) (*BeginnerReward, error)

	// GetPlayerRewards 获取玩家所有奖励
	GetPlayerRewards(ctx context.Context, playerID string) ([]*BeginnerReward, error)

	// GetUnclaimedRewards 获取未领取的奖励
	GetUnclaimedRewards(ctx context.Context, playerID string) ([]*BeginnerReward, error)

	// UpdateRewardStatus 更新奖励状态
	UpdateRewardStatus(ctx context.Context, playerID, rewardID string, claimed bool) error

	// GetGuideProgress 获取引导进度
	GetGuideProgress(ctx context.Context, playerID, guideID string) (*GuideProgress, error)

	// GetCompletedGuides 获取已完成的引导
	GetCompletedGuides(ctx context.Context, playerID string) ([]string, error)

	// IsGuideCompleted 检查引导是否完成
	IsGuideCompleted(ctx context.Context, playerID, guideID string) (bool, error)

	// IsTutorialCompleted 检查教程是否完成
	IsTutorialCompleted(ctx context.Context, playerID, tutorialID string) (bool, error)

	// GetTutorialsByCategory 根据分类获取教程
	GetTutorialsByCategory(ctx context.Context, playerID string, category TutorialCategory) ([]*Tutorial, error)

	// GetCompletedTutorials 获取已完成的教程
	GetCompletedTutorials(ctx context.Context, playerID string) ([]*Tutorial, error)
}

// GuideTemplateRepository 引导模板仓储接口
type GuideTemplateRepository interface {
	// GetGuideTemplate 获取引导模板
	GetGuideTemplate(ctx context.Context, templateID string) (*GuideTemplate, error)

	// GetGuideTemplatesByCategory 根据分类获取引导模板
	GetGuideTemplatesByCategory(ctx context.Context, category string) ([]*GuideTemplate, error)

	// SaveGuideTemplate 保存引导模板
	SaveGuideTemplate(ctx context.Context, template *GuideTemplate) error

	// DeleteGuideTemplate 删除引导模板
	DeleteGuideTemplate(ctx context.Context, templateID string) error

	// GetAllGuideTemplates 获取所有引导模板
	GetAllGuideTemplates(ctx context.Context) ([]*GuideTemplate, error)
}

// TutorialTemplateRepository 教程模板仓储接口
type TutorialTemplateRepository interface {
	// GetTutorialTemplate 获取教程模板
	GetTutorialTemplate(ctx context.Context, templateID string) (*TutorialTemplate, error)

	// GetTutorialTemplatesByCategory 根据分类获取教程模板
	GetTutorialTemplatesByCategory(ctx context.Context, category TutorialCategory) ([]*TutorialTemplate, error)

	// SaveTutorialTemplate 保存教程模板
	SaveTutorialTemplate(ctx context.Context, template *TutorialTemplate) error

	// DeleteTutorialTemplate 删除教程模板
	DeleteTutorialTemplate(ctx context.Context, templateID string) error

	// GetAllTutorialTemplates 获取所有教程模板
	GetAllTutorialTemplates(ctx context.Context) ([]*TutorialTemplate, error)
}

// GuideTemplate 引导模板
type GuideTemplate struct {
	ID            string               `json:"id"`
	Name          string               `json:"name"`
	Category      string               `json:"category"`
	Description   string               `json:"description"`
	Steps         []*GuideStepTemplate `json:"steps"`
	Prerequisites []string             `json:"prerequisites"`
	RequireLevel  int                  `json:"require_level"`
	Reward        *RewardTemplate      `json:"reward"`
	IsActive      bool                 `json:"is_active"`
}

// GuideStepTemplate 引导步骤模板
type GuideStepTemplate struct {
	StepID      int                      `json:"step_id"`
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	StepType    StepType                 `json:"step_type"`
	Conditions  []*StepConditionTemplate `json:"conditions"`
	Reward      *RewardTemplate          `json:"reward"`
}

// StepConditionTemplate 步骤条件模板
type StepConditionTemplate struct {
	ConditionType ConditionType      `json:"condition_type"`
	Target        string             `json:"target"`
	Value         interface{}        `json:"value"`
	Operator      ComparisonOperator `json:"operator"`
	Description   string             `json:"description"`
}

// TutorialTemplate 教程模板
type TutorialTemplate struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Category TutorialCategory `json:"category"`
	Content  string           `json:"content"`
	MediaURL string           `json:"media_url"`
	Duration int64            `json:"duration"` // 毫秒
	Reward   *RewardTemplate  `json:"reward"`
	IsActive bool             `json:"is_active"`
}

// RewardTemplate 奖励模板
type RewardTemplate struct {
	RewardType  RewardType             `json:"reward_type"`
	RewardData  map[string]interface{} `json:"reward_data"`
	Description string                 `json:"description"`
}

// CreateGuideFromTemplate 从模板创建引导步骤
func (gt *GuideTemplate) CreateGuideFromTemplate() []*GuideStep {
	steps := make([]*GuideStep, len(gt.Steps))

	for i, stepTemplate := range gt.Steps {
		step := NewGuideStep(
			stepTemplate.StepID,
			gt.ID,
			stepTemplate.Title,
			stepTemplate.Description,
			stepTemplate.StepType,
		)

		// 添加条件
		for _, condTemplate := range stepTemplate.Conditions {
			condition := NewStepCondition(
				condTemplate.ConditionType,
				condTemplate.Target,
				condTemplate.Value,
				condTemplate.Operator,
				condTemplate.Description,
			)
			step.AddCondition(condition)
		}

		// 添加奖励
		if stepTemplate.Reward != nil {
			reward := NewBeginnerReward(
				gt.ID+"_step_"+string(rune(stepTemplate.StepID))+"_reward",
				stepTemplate.Reward.RewardType,
				stepTemplate.Reward.RewardData,
				stepTemplate.Reward.Description,
			)
			step.SetReward(reward)
		}

		steps[i] = step
	}

	return steps
}

// CreateTutorialFromTemplate 从模板创建教程
func (tt *TutorialTemplate) CreateTutorialFromTemplate() *Tutorial {
	tutorial := NewTutorial(tt.Name, tt.Category, tt.Content)
	tutorial.SetMediaURL(tt.MediaURL)
	tutorial.SetDuration(time.Duration(tt.Duration) * time.Millisecond)

	// 添加奖励
	if tt.Reward != nil {
		reward := NewBeginnerReward(
			tt.ID+"_reward",
			tt.Reward.RewardType,
			tt.Reward.RewardData,
			tt.Reward.Description,
		)
		tutorial.SetReward(reward)
	}

	return tutorial
}
