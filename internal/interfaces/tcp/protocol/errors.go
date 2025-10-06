package protocol

import "errors"

// Error definitions for protocol
var (
	ErrInvalidMessage = errors.New("invalid message")
	ErrAuthFailed     = errors.New("authentication failed")
	ErrPlayerNotFound = errors.New("player not found")
	ErrBattleNotFound = errors.New("battle not found")
	ErrUnknownMessage = errors.New("unknown message type")
)
