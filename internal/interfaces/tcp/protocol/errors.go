package protocol

import "errors"

// Error codes
const (
	ErrCodeInvalidMessage = 1001
	ErrCodeAuthFailed     = 1002
	ErrCodePlayerNotFound = 1003
	ErrCodeBattleNotFound = 1004
	ErrCodeUnknownMessage = 1005
	ErrCodeServerBusy     = 1006
	ErrCodeInvalidPlayer  = 1007
	ErrCodeUnknown        = 1999
)

// Error definitions for protocol
var (
	ErrInvalidMessage = errors.New("invalid message")
	ErrAuthFailed     = errors.New("authentication failed")
	ErrPlayerNotFound = errors.New("player not found")
	ErrBattleNotFound = errors.New("battle not found")
	ErrUnknownMessage = errors.New("unknown message type")
	ErrServerBusy     = errors.New("server busy")
	ErrInvalidPlayer  = errors.New("invalid player")
)
