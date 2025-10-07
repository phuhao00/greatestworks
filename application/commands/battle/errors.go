package battle

import (
	"errors"
	protoerrors "greatestworks/internal/proto/errors"
)

// Battle命令相关错误码 - 使用proto生成的常量
const (
	ErrCodeInvalidCreatorID     = int32(protoerrors.BattleErrorCode_ERR_INVALID_CREATOR_ID)
	ErrCodeInvalidBattleType    = int32(protoerrors.BattleErrorCode_ERR_INVALID_BATTLE_TYPE)
	ErrCodeInvalidBattleID      = int32(protoerrors.BattleErrorCode_ERR_INVALID_BATTLE_ID)
	ErrCodeInvalidPlayerID      = int32(protoerrors.BattleErrorCode_ERR_INVALID_PLAYER_ID)
	ErrCodeInvalidTargetID      = int32(protoerrors.BattleErrorCode_ERR_INVALID_TARGET_ID)
	ErrCodeInvalidSkillID       = int32(protoerrors.BattleErrorCode_ERR_INVALID_SKILL_ID)
	ErrCodeInvalidTeam          = int32(protoerrors.BattleErrorCode_ERR_INVALID_TEAM)
	ErrCodeBattleNotFound       = int32(protoerrors.CommonErrorCode_ERR_BATTLE_NOT_FOUND)
	ErrCodeBattleAlreadyStarted = int32(protoerrors.BattleErrorCode_ERR_BATTLE_ALREADY_STARTED)
	ErrCodeBattleNotStarted     = int32(protoerrors.BattleErrorCode_ERR_BATTLE_NOT_STARTED)
	ErrCodePlayerNotInBattle    = int32(protoerrors.BattleErrorCode_ERR_PLAYER_NOT_IN_BATTLE)
	ErrCodeInsufficientMana     = int32(protoerrors.BattleErrorCode_ERR_INSUFFICIENT_MANA)
	ErrCodeSkillOnCooldown      = int32(protoerrors.BattleErrorCode_ERR_SKILL_ON_COOLDOWN)
	ErrCodeInvalidAction        = int32(protoerrors.BattleErrorCode_ERR_INVALID_ACTION)
)

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
	ErrInvalidAction        = errors.New("invalid action")
)
