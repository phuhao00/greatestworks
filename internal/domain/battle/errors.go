package battle

import "errors"

// 战斗领域错误定义
var (
	ErrBattleNotFound           = errors.New("battle not found")
	ErrBattleAlreadyStarted     = errors.New("battle already started")
	ErrBattleNotInProgress      = errors.New("battle not in progress")
	ErrPlayerNotInBattle        = errors.New("player not in battle")
	ErrPlayerAlreadyInBattle    = errors.New("player already in battle")
	ErrInsufficientParticipants = errors.New("insufficient participants")
	ErrPlayerDead               = errors.New("player is dead")
	ErrInvalidAction            = errors.New("invalid action")
	ErrActionOnCooldown         = errors.New("action on cooldown")
	ErrInsufficientMana         = errors.New("insufficient mana")
	ErrInvalidTarget            = errors.New("invalid target")
	ErrBattleFinished           = errors.New("battle is finished")
	ErrBattleAlreadyFinished    = errors.New("battle already finished")
	ErrBattleNotFinished        = errors.New("battle not finished")
)