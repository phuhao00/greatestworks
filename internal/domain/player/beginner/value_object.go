package beginner

// StepType 步骤类型
type StepType int

const (
	StepTypeDialog StepType = iota + 1
	StepTypeAction
	StepTypeNavigation
	StepTypeCombat
	StepTypeInventory
	StepTypeShop
	StepTypeQuest
	StepTypeSkill
	StepTypeSocial
	StepTypeCustom
)

// String 返回步骤类型字符串
func (st StepType) String() string {
	switch st {
	case StepTypeDialog:
		return "dialog"
	case StepTypeAction:
		return "action"
	case StepTypeNavigation:
		return "navigation"
	case StepTypeCombat:
		return "combat"
	case StepTypeInventory:
		return "inventory"
	case StepTypeShop:
		return "shop"
	case StepTypeQuest:
		return "quest"
	case StepTypeSkill:
		return "skill"
	case StepTypeSocial:
		return "social"
	case StepTypeCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// TutorialCategory 教程分类
type TutorialCategory int

const (
	TutorialCategoryBasic TutorialCategory = iota + 1
	TutorialCategoryCombat
	TutorialCategoryInventory
	TutorialCategorySkills
	TutorialCategorySocial
	TutorialCategoryEconomy
	TutorialCategoryAdvanced
	TutorialCategoryPvP
	TutorialCategoryGuild
	TutorialCategoryEndGame
)

// String 返回教程分类字符串
func (tc TutorialCategory) String() string {
	switch tc {
	case TutorialCategoryBasic:
		return "basic"
	case TutorialCategoryCombat:
		return "combat"
	case TutorialCategoryInventory:
		return "inventory"
	case TutorialCategorySkills:
		return "skills"
	case TutorialCategorySocial:
		return "social"
	case TutorialCategoryEconomy:
		return "economy"
	case TutorialCategoryAdvanced:
		return "advanced"
	case TutorialCategoryPvP:
		return "pvp"
	case TutorialCategoryGuild:
		return "guild"
	case TutorialCategoryEndGame:
		return "endgame"
	default:
		return "unknown"
	}
}

// RewardType 奖励类型
type RewardType int

const (
	RewardTypeGold RewardType = iota + 1
	RewardTypeExperience
	RewardTypeItem
	RewardTypeSkillPoint
	RewardTypeTitle
	RewardTypeMultiple
	RewardTypeCustom
)

// String 返回奖励类型字符串
func (rt RewardType) String() string {
	switch rt {
	case RewardTypeGold:
		return "gold"
	case RewardTypeExperience:
		return "experience"
	case RewardTypeItem:
		return "item"
	case RewardTypeSkillPoint:
		return "skill_point"
	case RewardTypeTitle:
		return "title"
	case RewardTypeMultiple:
		return "multiple"
	case RewardTypeCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// StepCondition 步骤条件值对象
type StepCondition struct {
	conditionType ConditionType
	target        string
	value         interface{}
	operator      ComparisonOperator
	description   string
}

// ConditionType 条件类型
type ConditionType int

const (
	ConditionTypeLevel ConditionType = iota + 1
	ConditionTypeLocation
	ConditionTypeItem
	ConditionTypeQuest
	ConditionTypeKill
	ConditionTypeInteract
	ConditionTypeEquip
	ConditionTypeSkill
	ConditionTypeGold
	ConditionTypeCustom
)

// String 返回条件类型字符串
func (ct ConditionType) String() string {
	switch ct {
	case ConditionTypeLevel:
		return "level"
	case ConditionTypeLocation:
		return "location"
	case ConditionTypeItem:
		return "item"
	case ConditionTypeQuest:
		return "quest"
	case ConditionTypeKill:
		return "kill"
	case ConditionTypeInteract:
		return "interact"
	case ConditionTypeEquip:
		return "equip"
	case ConditionTypeSkill:
		return "skill"
	case ConditionTypeGold:
		return "gold"
	case ConditionTypeCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// ComparisonOperator 比较操作符
type ComparisonOperator int

const (
	OperatorEqual ComparisonOperator = iota + 1
	OperatorGreaterThan
	OperatorLessThan
	OperatorGreaterEqual
	OperatorLessEqual
	OperatorNotEqual
	OperatorContains
	OperatorExists
)

// String 返回操作符字符串
func (co ComparisonOperator) String() string {
	switch co {
	case OperatorEqual:
		return "=="
	case OperatorGreaterThan:
		return ">"
	case OperatorLessThan:
		return "<"
	case OperatorGreaterEqual:
		return ">="
	case OperatorLessEqual:
		return "<="
	case OperatorNotEqual:
		return "!="
	case OperatorContains:
		return "contains"
	case OperatorExists:
		return "exists"
	default:
		return "unknown"
	}
}

// NewStepCondition 创建步骤条件
func NewStepCondition(conditionType ConditionType, target string, value interface{}, operator ComparisonOperator, description string) *StepCondition {
	return &StepCondition{
		conditionType: conditionType,
		target:        target,
		value:         value,
		operator:      operator,
		description:   description,
	}
}

// GetConditionType 获取条件类型
func (sc *StepCondition) GetConditionType() ConditionType {
	return sc.conditionType
}

// GetTarget 获取目标
func (sc *StepCondition) GetTarget() string {
	return sc.target
}

// GetValue 获取值
func (sc *StepCondition) GetValue() interface{} {
	return sc.value
}

// GetOperator 获取操作符
func (sc *StepCondition) GetOperator() ComparisonOperator {
	return sc.operator
}

// GetDescription 获取描述
func (sc *StepCondition) GetDescription() string {
	return sc.description
}

// IsMet 检查条件是否满足
func (sc *StepCondition) IsMet(playerData map[string]interface{}) bool {
	switch sc.conditionType {
	case ConditionTypeLevel:
		return sc.checkLevelCondition(playerData)
	case ConditionTypeLocation:
		return sc.checkLocationCondition(playerData)
	case ConditionTypeItem:
		return sc.checkItemCondition(playerData)
	case ConditionTypeQuest:
		return sc.checkQuestCondition(playerData)
	case ConditionTypeKill:
		return sc.checkKillCondition(playerData)
	case ConditionTypeInteract:
		return sc.checkInteractCondition(playerData)
	case ConditionTypeEquip:
		return sc.checkEquipCondition(playerData)
	case ConditionTypeSkill:
		return sc.checkSkillCondition(playerData)
	case ConditionTypeGold:
		return sc.checkGoldCondition(playerData)
	case ConditionTypeCustom:
		return sc.checkCustomCondition(playerData)
	default:
		return false
	}
}

// checkLevelCondition 检查等级条件
func (sc *StepCondition) checkLevelCondition(playerData map[string]interface{}) bool {
	if level, exists := playerData["level"]; exists {
		if playerLevel, ok := level.(int); ok {
			if targetLevel, ok := sc.value.(int); ok {
				return sc.compareValues(playerLevel, targetLevel)
			}
		}
	}
	return false
}

// checkLocationCondition 检查位置条件
func (sc *StepCondition) checkLocationCondition(playerData map[string]interface{}) bool {
	if location, exists := playerData["location"]; exists {
		if playerLocation, ok := location.(string); ok {
			if targetLocation, ok := sc.value.(string); ok {
				switch sc.operator {
				case OperatorEqual:
					return playerLocation == targetLocation
				case OperatorContains:
					return strings.Contains(playerLocation, targetLocation)
				default:
					return false
				}
			}
		}
	}
	return false
}

// checkItemCondition 检查物品条件
func (sc *StepCondition) checkItemCondition(playerData map[string]interface{}) bool {
	if inventory, exists := playerData["inventory"]; exists {
		if items, ok := inventory.(map[string]int); ok {
			if targetItem, ok := sc.value.(string); ok {
				if quantity, hasItem := items[targetItem]; hasItem {
					switch sc.operator {
					case OperatorExists:
						return quantity > 0
					case OperatorGreaterEqual:
						if requiredQuantity, ok := sc.target.(int); ok {
							return quantity >= requiredQuantity
						}
					default:
						return quantity > 0
					}
				}
			}
		}
	}
	return false
}

// checkQuestCondition 检查任务条件
func (sc *StepCondition) checkQuestCondition(playerData map[string]interface{}) bool {
	if quests, exists := playerData["quests"]; exists {
		if questMap, ok := quests.(map[string]string); ok {
			if targetQuest, ok := sc.value.(string); ok {
				if status, hasQuest := questMap[targetQuest]; hasQuest {
					if targetStatus, ok := sc.target.(string); ok {
						return status == targetStatus
					}
				}
			}
		}
	}
	return false
}

// checkKillCondition 检查击杀条件
func (sc *StepCondition) checkKillCondition(playerData map[string]interface{}) bool {
	if kills, exists := playerData["kills"]; exists {
		if killMap, ok := kills.(map[string]int); ok {
			if targetMonster, ok := sc.value.(string); ok {
				if killCount, hasKill := killMap[targetMonster]; hasKill {
					if requiredCount, ok := sc.target.(int); ok {
						return sc.compareValues(killCount, requiredCount)
					}
				}
			}
		}
	}
	return false
}

// checkInteractCondition 检查交互条件
func (sc *StepCondition) checkInteractCondition(playerData map[string]interface{}) bool {
	if interactions, exists := playerData["interactions"]; exists {
		if interactionList, ok := interactions.([]string); ok {
			if targetInteraction, ok := sc.value.(string); ok {
				for _, interaction := range interactionList {
					if interaction == targetInteraction {
						return true
					}
				}
			}
		}
	}
	return false
}

// checkEquipCondition 检查装备条件
func (sc *StepCondition) checkEquipCondition(playerData map[string]interface{}) bool {
	if equipment, exists := playerData["equipment"]; exists {
		if equipMap, ok := equipment.(map[string]string); ok {
			if targetSlot, ok := sc.target.(string); ok {
				if targetItem, ok := sc.value.(string); ok {
					if equippedItem, hasSlot := equipMap[targetSlot]; hasSlot {
						return equippedItem == targetItem
					}
				}
			}
		}
	}
	return false
}

// checkSkillCondition 检查技能条件
func (sc *StepCondition) checkSkillCondition(playerData map[string]interface{}) bool {
	if skills, exists := playerData["skills"]; exists {
		if skillMap, ok := skills.(map[string]int); ok {
			if targetSkill, ok := sc.value.(string); ok {
				if skillLevel, hasSkill := skillMap[targetSkill]; hasSkill {
					if requiredLevel, ok := sc.target.(int); ok {
						return sc.compareValues(skillLevel, requiredLevel)
					}
				}
			}
		}
	}
	return false
}

// checkGoldCondition 检查金币条件
func (sc *StepCondition) checkGoldCondition(playerData map[string]interface{}) bool {
	if gold, exists := playerData["gold"]; exists {
		if playerGold, ok := gold.(int); ok {
			if requiredGold, ok := sc.value.(int); ok {
				return sc.compareValues(playerGold, requiredGold)
			}
		}
	}
	return false
}

// checkCustomCondition 检查自定义条件
func (sc *StepCondition) checkCustomCondition(playerData map[string]interface{}) bool {
	// 自定义条件的实现可以根据具体需求来定制
	if customData, exists := playerData[sc.target]; exists {
		switch sc.operator {
		case OperatorExists:
			return customData != nil
		case OperatorEqual:
			return customData == sc.value
		default:
			return false
		}
	}
	return false
}

// compareValues 比较数值
func (sc *StepCondition) compareValues(actual, expected int) bool {
	switch sc.operator {
	case OperatorEqual:
		return actual == expected
	case OperatorGreaterThan:
		return actual > expected
	case OperatorLessThan:
		return actual < expected
	case OperatorGreaterEqual:
		return actual >= expected
	case OperatorLessEqual:
		return actual <= expected
	case OperatorNotEqual:
		return actual != expected
	default:
		return false
	}
}