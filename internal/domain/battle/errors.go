package battle

import "greatestworks/internal/errors"

// 战斗领域错误定义 - 使用统一的错误处理机制
var (
	ErrBattleNotFound           = errors.ErrBattleNotFound
	ErrBattleAlreadyStarted     = errors.ErrBattleAlreadyStarted
	ErrBattleNotInProgress      = errors.ErrBattleNotInProgress
	ErrPlayerNotInBattle        = errors.ErrPlayerNotInBattle
	ErrPlayerAlreadyInBattle    = errors.ErrPlayerAlreadyInBattle
	ErrInsufficientParticipants = errors.ErrInsufficientParticipants
	ErrInsufficientMana         = errors.ErrInsufficientMana
	ErrInvalidTarget            = errors.ErrInvalidTarget
	ErrBattleFinished           = errors.ErrBattleFinished
	ErrPlayerDead               = errors.NewDomainError("PLAYER_DEAD", "玩家已死亡")
	ErrInvalidAction            = errors.NewDomainError("INVALID_ACTION", "无效的行动")
	ErrActionOnCooldown         = errors.NewDomainError("ACTION_ON_COOLDOWN", "行动冷却中")
	ErrBattleAlreadyFinished    = errors.NewDomainError("BATTLE_ALREADY_FINISHED", "战斗已结束")
	ErrBattleNotFinished        = errors.NewDomainError("BATTLE_NOT_FINISHED", "战斗未结束")
)
