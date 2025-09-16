package battle

import "errors"

// Battle命令相关错误
var (
	ErrInvalidCreatorID     = errors.New("invalid creator id")
	ErrInvalidBattleType    = errors.New("invalid battle type")
	ErrInvalidBattleID      = errors.New("invalid battle id")
	ErrInvalidPlayerID      = errors.New("invalid player id")
	ErrInvalidTargetID      = errors.New("invalid target id")
	ErrInvalidSkillID       = errors.New("invalid skill id")
	ErrInvalidTeam          = errors.New("invalid team")
	ErrBattleNotFound       = errors.New("battle not found")
	ErrBattleAlreadyStarted = errors.New("battle already started")
	ErrBattleNotStarted     = errors.New("battle not started")
	ErrPlayerNotInBattle    = errors.New("player not in battle")
	ErrInsufficientMana     = errors.New("insufficient mana")
	ErrSkillOnCooldown      = errors.New("skill on cooldown")
	ErrInvalidAction        = errors.New("invali