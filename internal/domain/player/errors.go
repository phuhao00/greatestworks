package player

import "greatestworks/internal/errors"

// 玩家领域错误定义 - 使用统一的错误处理机制
var (
	ErrPlayerNotFound      = errors.ErrPlayerNotFound
	ErrPlayerOffline       = errors.ErrPlayerOffline
	ErrPlayerAlreadyExists = errors.ErrPlayerAlreadyExists
	ErrInvalidPlayerName   = errors.ErrInvalidPlayerName
	ErrInvalidPosition     = errors.ErrInvalidPosition
	ErrVersionMismatch     = errors.ErrVersionMismatch
	ErrPlayerDead          = errors.NewDomainError("PLAYER_DEAD", "玩家已死亡")
	ErrInvalidLevel        = errors.NewDomainError("INVALID_LEVEL", "无效的等级")
	ErrInsufficientExp     = errors.NewDomainError("INSUFFICIENT_EXP", "经验值不足")
)
