package player

import "errors"

// Player查询相关错误
var (
	ErrPlayerNotFound    = errors.New("player not found")
	ErrInvalidPlayerID   = errors.New("invalid player id")
	ErrInvalidPlayerName = errors.New("invalid player name")
	ErrInvalidLevel      = errors.New("invalid level")
	ErrInvalidLimit      = errors.New("invalid limit")
	ErrQueryFailed       = errors.New("query failed")
	ErrInvalidParameters = errors.New("invalid parameters")
	ErrInvalidUsername   = errors.New("invalid username")
)
