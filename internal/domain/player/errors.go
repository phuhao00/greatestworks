package player

import "errors"

// 玩家领域错误定义
var (
	ErrPlayerNotFound    = errors.New("player not found")
	ErrPlayerOffline     = errors.New("player is offline")
	ErrInvalidPlayerName = errors.New("invalid player name")
	ErrPlayerAlreadyExists = errors.New("player already exists")
	ErrInvalidLevel      = errors.New("invalid level")
	ErrInsufficientExp   = errors.New("insufficient experience")
	ErrPlayerDead        = errors.New("player is dead")
	ErrInvalidPosition   = errors.New("invalid position")
	ErrVersionMismatch   = errors.New("version mismatch")
)