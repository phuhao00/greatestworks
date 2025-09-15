package beginner

import (
	"time"
)

// BeginnerService 新手领域服务
type BeginnerService struct {
	guideFactory    *GuideFactory
	tutorialFactory *TutorialFactory
	rewardCalculator *RewardCalculator
}

// NewBeginnerService 创建新手服务
func NewBeginnerService() *BeginnerService {
	return &BeginnerService{
		guideFactory:    NewGuideFactory(),
		tutorialFactory: NewTutorialFactory(),
		rewardCalculator: NewRewardCalculator(),
	}
}

// ValidateStepCompletion 验证步骤完成
func (bs *BeginnerService) ValidateStepCompletion(step *GuideStep, playerData map[string]interface{}) error {
	if step == nil {
		return ErrInvalidStep
	}
	
	if step.IsCompleted() {
		return ErrStepAlreadyCompleted
	}
	
	// 检查完成条件
	if !step.CheckConditions(playerData) {
		return ErrConditionNotMet
	}
	
	return nil
}

// CalculateGuideReward 计算引导奖励
func (bs *BeginnerService) CalculateGuideReward(guideID string, playerLevel int) *BeginnerReward {
	return bs.rewardCalculator.CalculateGuideReward(guideID, playerLevel)
}

// CalculateTutorialReward 计算教程奖励
func (bs *BeginnerService) CalculateTutorialReward(category TutorialCategory, playerLevel int) *BeginnerReward {
	return bs.rewardCalculator.CalculateTutorialReward(category, playerLevel)
}

// CreateGuideFromTemplate 从模板创建引导
func (bs *BeginnerService) CreateGuideFromTemplate(templateID string, playerLevel int) ([]*GuideStep, error) {
	return bs.guideFactory.CreateFromTemplate(templateID, playerLevel)
}

// CreateTutorialFromTemplate 从模板创建教程
func (bs *BeginnerService) CreateTutorialFromTemplate(templateID string) (*Tutorial, error) {
	return bs.tutorialFactory.CreateFromTemplate(templateID)
}

// CheckPrerequisites 检查前置条件
func (bs *BeginnerService) CheckPrerequisites(guideID string, completedGuides []string) bool {
	prerequisites := bs.getGuidePrerequisites(guideID)
	
	for _, prerequisite := range prerequisites {
		found := false
		for _, completed := range completedGuides {
			if completed == prerequisite {
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

// GetRecommendedNextGuide 获取推荐的下一个引导
func (bs *BeginnerService) GetRecommendedNextGuide(completedGuides []string, playerLevel int) string {
	allGuides := []string{"main_guide", "combat_guide", "inventory_guide", "social_guide", "economy_guide"}
	
	for _, guide := range allGuides {
		// 检查是否已完成
		alreadyCompleted := false
		for _, completed := range completedGuides {
			if completed == guide {
				alreadyCompleted = true
				break
			}
		}
		
		if !alreadyCompleted && bs.CheckPrerequisites(guide, completedGuides) {
			return guide
		}
	}
	
	return ""
}

// getGuidePrerequisites 获取引导前置条件
func (bs *BeginnerService) getGuidePrerequisites(guideID string) []string {
	prerequisites := map[string][]string{
		"main_guide":      {},
		"combat_guide":    {"main_guide"},
		"inventory_guide": {"main_guide"},
		"social_guide":    {"main_guide", "combat_guide"},
		"economy_guide":   {"main_guide", "inventory_guide"},
	}
	
	return prerequisites[guideID]
}

// GuideFactory 引导工厂
type GuideFactory struct{}

// NewGuideFactory 创建引导工厂
func NewGuideFactory() *GuideFactory {
	return &GuideFactory{}
}

// CreateFromTemplate 从模板创建引导
func (gf *GuideFactory) CreateFromTemplate(templateID string, playerLevel int) ([]*GuideStep, error) {
	switch templateID {
	case "main_guide":
		return gf.createMainGuide(), nil
	case "combat_guide":
		return gf.createCombatGuide(), nil
	case "inventory_guide":
		return gf.createInventoryGuide(), nil
	case "social_guide":
		return gf.createSocialGuide(), nil
	case "economy_guide":
		return gf.createEconomyGuide(), nil
	default:
		return nil, ErrInvalidGuide
	}
}

// createMainGuide 创建主引导
func (gf *GuideFactory) createMainGuide() []*GuideStep {
	steps := []*GuideStep{
		NewGuideStep(1, "main_guide", "欢迎来到游戏", "欢迎来到这个奇幻世界！让我们开始你的冒险之旅。", StepTypeDialog),
		NewGuideStep(2, "main_guide", "移动角色", "使用WASD键或方向键移动你的角色。", StepTypeNavigation),
		NewGuideStep(3, "main_guide", "查看角色信息", "按C键打开角色面板，查看你的属性。", StepTypeAction),
		NewGuideStep(4, "main_guide", "完成第一个任务", "前往村长处接受你的第一个任务。", StepTypeQuest),
	}
	
	// 添加条件和奖励
	steps[1].AddCondition(NewStepCondition(ConditionTypeLocation, "starting_area", "moved", OperatorEqual, "移动到指定位置"))
	steps[2].AddCondition(NewStepCondition(ConditionTypeInteract, "character_panel", "opened", OperatorEqual, "打开角色面板"))
	steps[3].AddCondition(NewStepCondition(ConditionTypeQuest, "first_quest", "accepted", OperatorEqual, "接受第一个任务"))
	
	// 设置奖励
	steps[0].SetReward(NewBeginnerReward("welcome_reward", RewardTypeGold, map[string]interface{}{"gold": 100}, "欢迎奖励"))
	steps[3].SetReward(NewBeginnerReward("quest_reward", RewardTypeMultiple, map[string]interface{}{
		"gold":       200,
		"experience": 100,
		"items":      []string{"beginner_sword"},
	}, "任务完成奖励"))
	
	return steps
}

// createCombatGuide 创建战斗引导
func (gf *GuideFactory) createCombatGuide() []*GuideStep {
	steps := []*GuideStep{
		NewGuideStep(1, "combat_guide", "学习攻击", "点击鼠标左键或按空格键攻击敌人。", StepTypeCombat),
		NewGuideStep(2, "combat_guide", "使用技能", "按1-4数字键使用技能攻击敌人。", StepTypeCombat),
		NewGuideStep(3, "combat_guide", "击败怪物", "击败3只史莱姆来完成战斗训练。", StepTypeCombat),
	}
	
	// 添加条件
	steps[0].AddCondition(NewStepCondition(ConditionTypeInteract, "attack", "used", OperatorEqual, "使用攻击"))
	steps[1].AddCondition(NewStepCondition(ConditionTypeSkill, "basic_skill", "used", OperatorEqual, "使用技能"))
	steps[2].AddCondition(NewStepCondition(ConditionTypeKill, "slime", 3, OperatorGreaterEqual, "击败3只史莱姆"))
	
	// 设置奖励
	steps[2].SetReward(NewBeginnerReward("combat_reward", RewardTypeMultiple, map[string]interface{}{
		"experience": 200,
		"items":      []string{"health_potion", "mana_potion"},
	}, "战斗训练奖励"))
	
	return steps
}

// createInventoryGuide 创建背包引导
func (gf *GuideFactory) createInventoryGuide() []*GuideStep {
	steps := []*GuideStep{
		NewGuideStep(1, "inventory_guide", "打开背包", "按I键打开背包界面。", StepTypeInventory),
		NewGuideStep(2, "inventory_guide", "装备物品", "将获得的武器装备到武器槽。", StepTypeInventory),
		NewGuideStep(3, "inventory_guide", "使用消耗品", "右键点击药水来恢复生命值。", StepTypeInventory),
	}
	
	// 添加条件
	steps[0].AddCondition(NewStepCondition(ConditionTypeInteract, "inventory", "opened", OperatorEqual, "打开背包"))
	steps[1].AddCondition(NewStepCondition(ConditionTypeEquip, "weapon", "beginner_sword", OperatorEqual, "装备新手剑"))
	steps[2].AddCondition(NewStepCondition(ConditionTypeInteract, "use_item", "health_potion", OperatorEqual, "使用生命药水"))
	
	return steps
}

// createSocialGuide 创建社交引导
func (gf *GuideFactory) createSocialGuide() []*GuideStep {
	steps := []*GuideStep{
		NewGuideStep(1, "social_guide", "添加好友", "学习如何添加其他玩家为好友。", StepTypeSocial),
		NewGuideStep(2, "social_guide", "发送消息", "向好友发送一条消息。", StepTypeSocial),
		NewGuideStep(3, "social_guide", "加入公会", "申请加入一个公会。", StepTypeSocial),
	}
	
	return steps
}

// createEconomyGuide 创建经济引导
func (gf *GuideFactory) createEconomyGuide() []*GuideStep {
	steps := []*GuideStep{
		NewGuideStep(1, "economy_guide", "访问商店", "前往商店购买物品。", StepTypeShop),
		NewGuideStep(2, "economy_guide", "出售物品", "将不需要的物品出售给商人。", StepTypeShop),
		NewGuideStep(3, "economy_guide", "交易系统", "学习如何与其他玩家交易。", StepTypeShop),
	}
	
	return steps
}

// TutorialFactory 教程工厂
type TutorialFactory struct{}

// NewTutorialFactory 创建教程工厂
func NewTutorialFactory() *TutorialFactory {
	return &TutorialFactory{}
}

// CreateFromTemplate 从模板创建教程
func (tf *TutorialFactory) CreateFromTemplate(templateID string) (*Tutorial, error) {
	switch templateID {
	case "basic_controls":
		return tf.createBasicControlsTutorial(), nil
	case "combat_basics":
		return tf.createCombatBasicsTutorial(), nil
	case "inventory_management":
		return tf.createInventoryManagementTutorial(), nil
	default:
		return nil, ErrInvalidTutorial
	}
}

// createBasicControlsTutorial 创建基础操作教程
func (tf *TutorialFactory) createBasicControlsTutorial() *Tutorial {
	tutorial := NewTutorial("基础操作", TutorialCategoryBasic, "学习游戏的基础操作方法")
	tutorial.SetDuration(time.Minute * 3)
	tutorial.SetMediaURL("https://example.com/tutorials/basic_controls.mp4")
	tutorial.SetReward(NewBeginnerReward("basic_tutorial_reward", RewardTypeExperience, map[string]interface{}{"experience": 50}, "基础教程奖励"))
	return tutorial
}

// createCombatBasicsTutorial 创建战斗基础教程
func (tf *TutorialFactory) createCombatBasicsTutorial() *Tutorial {
	tutorial := NewTutorial("战斗基础", TutorialCategoryCombat, "学习战斗系统的基本操作")
	tutorial.SetDuration(time.Minute * 5)
	tutorial.SetMediaURL("https://example.com/tutorials/combat_basics.mp4")
	tutorial.SetReward(NewBeginnerReward("combat_tutorial_reward", RewardTypeSkillPoint, map[string]interface{}{"skill_points": 1}, "战斗教程奖励"))
	return tutorial
}

// createInventoryManagementTutorial 创建背包管理教程
func (tf *TutorialFactory) createInventoryManagementTutorial() *Tutorial {
	tutorial := NewTutorial("背包管理", TutorialCategoryInventory, "学习如何有效管理背包空间")
	tutorial.SetDuration(time.Minute * 4)
	tutorial.SetMediaURL("https://example.com/tutorials/inventory_management.mp4")
	tutorial.SetReward(NewBeginnerReward("inventory_tutorial_reward", RewardTypeItem, map[string]interface{}{"items": []string{"bag_expansion"}}, "背包教程奖励"))
	return tutorial
}

// RewardCalculator 奖励计算器
type RewardCalculator struct{}

// NewRewardCalculator 创建奖励计算器
func NewRewardCalculator() *RewardCalculator {
	return &RewardCalculator{}
}

// CalculateGuideReward 计算引导奖励
func (rc *RewardCalculator) CalculateGuideReward(guideID string, playerLevel int) *BeginnerReward {
	baseReward := map[string]interface{}{
		"gold":       100,
		"experience": 50,
	}
	
	// 根据引导类型调整奖励
	switch guideID {
	case "main_guide":
		baseReward["gold"] = 200
		baseReward["experience"] = 100
		baseReward["items"] = []string{"beginner_weapon"}
	case "combat_guide":
		baseReward["experience"] = 150
		baseReward["items"] = []string{"health_potion", "mana_potion"}
	case "inventory_guide":
		baseReward["gold"] = 150
		baseReward["items"] = []string{"bag_expansion"}
	}
	
	// 根据玩家等级调整奖励
	levelMultiplier := 1.0 + float64(playerLevel)*0.1
	if gold, ok := baseReward["gold"].(int); ok {
		baseReward["gold"] = int(float64(gold) * levelMultiplier)
	}
	if exp, ok := baseReward["experience"].(int); ok {
		baseReward["experience"] = int(float64(exp) * levelMultiplier)
	}
	
	return NewBeginnerReward(guideID+"_reward", RewardTypeMultiple, baseReward, "引导完成奖励")
}

// CalculateTutorialReward 计算教程奖励
func (rc *RewardCalculator) CalculateTutorialReward(category TutorialCategory, playerLevel int) *BeginnerReward {
	baseReward := map[string]interface{}{
		"experience": 25,
	}
	
	// 根据教程分类调整奖励
	switch category {
	case TutorialCategoryBasic:
		baseReward["experience"] = 50
	case TutorialCategoryCombat:
		baseReward["experience"] = 75
		baseReward["skill_points"] = 1
	case TutorialCategoryInventory:
		baseReward["gold"] = 100
	case TutorialCategorySkills:
		baseReward["skill_points"] = 2
	}
	
	return NewBeginnerReward("tutorial_"+category.String()+"_reward", RewardTypeMultiple, baseReward, "教程完成奖励")
}