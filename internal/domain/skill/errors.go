package skill

import "errors"

var (
	// 技能树相关错误
	ErrSkillTreeNotFound = errors.New("skill tree not found")
	ErrInvalidSkillTree  = errors.New("invalid skill tree")

	// 技能相关错误
	ErrSkillNotFound         = errors.New("skill not found")
	ErrSkillNotLearned       = errors.New("skill not learned")
	ErrSkillAlreadyLearned   = errors.New("skill already learned")
	ErrSkillMaxLevel         = errors.New("skill is at maximum level")
	ErrSkillOnCooldown       = errors.New("skill is on cooldown")
	ErrPassiveSkillNotUsable = errors.New("passive skill cannot be used actively")
	ErrInvalidSkillType      = errors.New("invalid skill type")
	ErrSkillNotAvailable     = errors.New("skill is not available")

	// 技能点相关错误
	ErrInsufficientSkillPoints = errors.New("insufficient skill points")
	ErrInvalidSkillPoints      = errors.New("invalid skill points")
	ErrSkillPointsOverflow     = errors.New("skill points overflow")

	// 前置条件相关错误
	ErrPrerequisitesNotMet    = errors.New("prerequisites not met")
	ErrLevelRequirementNotMet = errors.New("level requirement not met")
	ErrClassRestriction       = errors.New("class restriction for skill")
	ErrRaceRestriction        = errors.New("race restriction for skill")

	// 技能使用相关错误
	ErrInsufficientMana    = errors.New("insufficient mana")
	ErrInsufficientStamina = errors.New("insufficient stamina")
	ErrInvalidTarget       = errors.New("invalid target")
	ErrTargetOutOfRange    = errors.New("target out of range")
	ErrTargetDead          = errors.New("target is dead")
	ErrCastInterrupted     = errors.New("cast interrupted")
	ErrSilenced            = errors.New("player is silenced")
	ErrStunned             = errors.New("player is stunned")

	// 技能效果相关错误
	ErrInvalidEffect  = errors.New("invalid skill effect")
	ErrEffectNotFound = errors.New("skill effect not found")
	ErrEffectExpired  = errors.New("skill effect expired")
	ErrEffectImmune   = errors.New("target is immune to effect")

	// 技能组合相关错误
	ErrInvalidCombo     = errors.New("invalid skill combo")
	ErrComboTimeout     = errors.New("skill combo timeout")
	ErrComboInterrupted = errors.New("skill combo interrupted")

	// 配置相关错误
	ErrInvalidSkillConfig  = errors.New("invalid skill configuration")
	ErrSkillConfigNotFound = errors.New("skill configuration not found")
)
