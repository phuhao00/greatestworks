package player

import "errors"

// Player命令相关错误
var (
	ErrInvalidPlayerName       = errors.New("invalid player name")
	ErrInvalidPlayerNameLength = errors.New("player name must be between 2 and 20 characters")
	ErrPlayerNotFound          = errors.New("player not found")
	ErrPlayerAlreadyExists     = errors.New("player already exists")
	ErrInvalidPlayerID         = errors.New("invalid player id")
	ErrInvalidPosition         = errors.New("invalid position")
	ErrInvalidExperience       = errors.New("invalid experience amount")
	ErrInvalidHealAmount       = errors.New("invalid heal amount")
	ErrPlayerOffline           = errors.New("player is offline")
	ErrInsufficientPermission  = errors.New("insufficient permission")
	ErrInvalidUsername         = errors.New("invalid username")
	ErrInvalidPassword         = errors.New("invalid password")
	ErrInvalidEmail            = errors.New("invalid email")
	ErrInvalidRequest          = errors.New("invalid request")
)
