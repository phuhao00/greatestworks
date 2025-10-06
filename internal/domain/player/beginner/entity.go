package beginner

import (
	// "fmt" // 未使用
	// "strings" // 未使用
	"time"

	"github.com/google/uuid"
)

// GuideStep 引导步骤实体
type GuideStep struct {
	id          string
	stepID      int
	guideID     string
	title       string
	description string
	stepType    StepType
	conditions  []*StepCondition
	reward      *BeginnerReward
	isCompleted bool
	completedAt *time.Time
	createdAt   time.Time
}

// NewGuideStep 创建引导步骤
func NewGuideStep(stepID int, guideID, title, description string, stepType StepType) *GuideStep {
	return &GuideStep{
		id:          uuid.New().String(),
		stepID:      stepID,
		guideID:     guideID,
		title:       title,
		description: description,
		stepType:    stepType,
		conditions:  make([]*StepCondition, 0),
		isCompleted: false,
		createdAt:   time.Now(),
	}
}

// GetID 获取步骤ID
func (gs *GuideStep) GetID() string {
	return gs.id
}

// GetStepID 获取步骤序号
func (gs *GuideStep) GetStepID() int {
	return gs.stepID
}

// GetGuideID 获取引导ID
func (gs *GuideStep) GetGuideID() string {
	return gs.guideID
}

// GetTitle 获取标题
func (gs *GuideStep) GetTitle() string {
	return gs.title
}

// GetDescription 获取描述
func (gs *GuideStep) GetDescription() string {
	return gs.description
}

// GetStepType 获取步骤类型
func (gs *GuideStep) GetStepType() StepType {
	return gs.stepType
}

// AddCondition 添加完成条件
func (gs *GuideStep) AddCondition(condition *StepCondition) {
	gs.conditions = append(gs.conditions, condition)
}

// GetConditions 获取完成条件
func (gs *GuideStep) GetConditions() []*StepCondition {
	return gs.conditions
}

// SetReward 设置奖励
func (gs *GuideStep) SetReward(reward *BeginnerReward) {
	gs.reward = reward
}

// GetReward 获取奖励
func (gs *GuideStep) GetReward() *BeginnerReward {
	return gs.reward
}

// HasReward 是否有奖励
func (gs *GuideStep) HasReward() bool {
	return gs.reward != nil
}

// Complete 完成步骤
func (gs *GuideStep) Complete() {
	gs.isCompleted = true
	now := time.Now()
	gs.completedAt = &now
}

// IsCompleted 是否已完成
func (gs *GuideStep) IsCompleted() bool {
	return gs.isCompleted
}

// GetCompletedAt 获取完成时间
func (gs *GuideStep) GetCompletedAt() *time.Time {
	return gs.completedAt
}

// GetCreatedAt 获取创建时间
func (gs *GuideStep) GetCreatedAt() time.Time {
	return gs.createdAt
}

// CheckConditions 检查完成条件
func (gs *GuideStep) CheckConditions(playerData map[string]interface{}) bool {
	for _, condition := range gs.conditions {
		if !condition.IsMet(playerData) {
			return false
		}
	}
	return true
}

// Tutorial 教程实体
type Tutorial struct {
	id          string
	name        string
	category    TutorialCategory
	content     string
	mediaURL    string
	duration    time.Duration
	reward      *BeginnerReward
	isCompleted bool
	completedAt *time.Time
	createdAt   time.Time
}

// NewTutorial 创建教程
func NewTutorial(name string, category TutorialCategory, content string) *Tutorial {
	return &Tutorial{
		id:          uuid.New().String(),
		name:        name,
		category:    category,
		content:     content,
		duration:    time.Minute * 5, // 默认5分钟
		isCompleted: false,
		createdAt:   time.Now(),
	}
}

// GetID 获取教程ID
func (t *Tutorial) GetID() string {
	return t.id
}

// GetName 获取教程名称
func (t *Tutorial) GetName() string {
	return t.name
}

// GetCategory 获取教程分类
func (t *Tutorial) GetCategory() TutorialCategory {
	return t.category
}

// GetContent 获取教程内容
func (t *Tutorial) GetContent() string {
	return t.content
}

// SetContent 设置教程内容
func (t *Tutorial) SetContent(content string) {
	t.content = content
}

// GetMediaURL 获取媒体URL
func (t *Tutorial) GetMediaURL() string {
	return t.mediaURL
}

// SetMediaURL 设置媒体URL
func (t *Tutorial) SetMediaURL(url string) {
	t.mediaURL = url
}

// GetDuration 获取时长
func (t *Tutorial) GetDuration() time.Duration {
	return t.duration
}

// SetDuration 设置时长
func (t *Tutorial) SetDuration(duration time.Duration) {
	t.duration = duration
}

// SetReward 设置奖励
func (t *Tutorial) SetReward(reward *BeginnerReward) {
	t.reward = reward
}

// GetReward 获取奖励
func (t *Tutorial) GetReward() *BeginnerReward {
	return t.reward
}

// HasReward 是否有奖励
func (t *Tutorial) HasReward() bool {
	return t.reward != nil
}

// Complete 完成教程
func (t *Tutorial) Complete() {
	t.isCompleted = true
	now := time.Now()
	t.completedAt = &now
}

// IsCompleted 是否已完成
func (t *Tutorial) IsCompleted() bool {
	return t.isCompleted
}

// GetCompletedAt 获取完成时间
func (t *Tutorial) GetCompletedAt() *time.Time {
	return t.completedAt
}

// GetCreatedAt 获取创建时间
func (t *Tutorial) GetCreatedAt() time.Time {
	return t.createdAt
}

// BeginnerReward 新手奖励实体
type BeginnerReward struct {
	id          string
	rewardType  RewardType
	rewardData  map[string]interface{}
	description string
	isClaimed   bool
	claimedAt   *time.Time
	createdAt   time.Time
}

// NewBeginnerReward 创建新手奖励
func NewBeginnerReward(id string, rewardType RewardType, rewardData map[string]interface{}, description string) *BeginnerReward {
	return &BeginnerReward{
		id:          id,
		rewardType:  rewardType,
		rewardData:  rewardData,
		description: description,
		isClaimed:   false,
		createdAt:   time.Now(),
	}
}

// GetID 获取奖励ID
func (br *BeginnerReward) GetID() string {
	return br.id
}

// GetRewardType 获取奖励类型
func (br *BeginnerReward) GetRewardType() RewardType {
	return br.rewardType
}

// GetRewardData 获取奖励数据
func (br *BeginnerReward) GetRewardData() map[string]interface{} {
	return br.rewardData
}

// GetDescription 获取描述
func (br *BeginnerReward) GetDescription() string {
	return br.description
}

// Claim 领取奖励
func (br *BeginnerReward) Claim() {
	br.isClaimed = true
	now := time.Now()
	br.claimedAt = &now
}

// IsClaimed 是否已领取
func (br *BeginnerReward) IsClaimed() bool {
	return br.isClaimed
}

// GetClaimedAt 获取领取时间
func (br *BeginnerReward) GetClaimedAt() *time.Time {
	return br.claimedAt
}

// GetCreatedAt 获取创建时间
func (br *BeginnerReward) GetCreatedAt() time.Time {
	return br.createdAt
}

// GetGoldAmount 获取金币数量
func (br *BeginnerReward) GetGoldAmount() int {
	if gold, exists := br.rewardData["gold"]; exists {
		if amount, ok := gold.(int); ok {
			return amount
		}
	}
	return 0
}

// GetExperienceAmount 获取经验数量
func (br *BeginnerReward) GetExperienceAmount() int {
	if exp, exists := br.rewardData["experience"]; exists {
		if amount, ok := exp.(int); ok {
			return amount
		}
	}
	return 0
}

// GetItems 获取物品列表
func (br *BeginnerReward) GetItems() []string {
	if items, exists := br.rewardData["items"]; exists {
		if itemList, ok := items.([]string); ok {
			return itemList
		}
	}
	return []string{}
}

// GuideProgress 引导进度实体
type GuideProgress struct {
	GuideID        string  `json:"guide_id"`
	TotalSteps     int     `json:"total_steps"`
	CompletedSteps int     `json:"completed_steps"`
	IsCompleted    bool    `json:"is_completed"`
	Progress       float64 `json:"progress"`
}

// GetProgressPercentage 获取进度百分比
func (gp *GuideProgress) GetProgressPercentage() int {
	return int(gp.Progress * 100)
}

// GetRemainingSteps 获取剩余步骤数
func (gp *GuideProgress) GetRemainingSteps() int {
	return gp.TotalSteps - gp.CompletedSteps
}
