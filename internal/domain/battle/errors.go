package battle

import "errors"

// 战斗领域错误定义
var (
	ErrBattleNotFound           = errors.New("battle not found")
	ErrBattleAlreadyStarted     = errors.New("battle already started")
	ErrBattleNotInProgress      = errors.New("battle not in progress")
	ErrPlayerNotInBattle        = errors.New("player not in battle")
	ErrPlayerAlreadyInBattle    = errors.New("player already in battle")
	ErrPlayerDead               = errors.New("player is dead")
	ErrInsufficientParticipants = errors.New("insufficient participants")
	ErrInvalidBattleType        = errors.New("invalid battle type")
	ErrInvalidActionType        = errors.New("invalid action type")
	ErrInvalidTarget            = errors.New("invalid target")
	ErrBattleFinished           = errors.New("battle is finished")
	ErrVersionMismatch          = errors.New("version mismatch")
)