package protocol

import (
	"errors"
	protoerrors "greatestworks/internal/proto/errors"
)

// Error codes - 使用proto生成的常量
const (
	ErrCodeInvalidMessage = int32(protoerrors.CommonErrorCode_ERR_INVALID_MESSAGE)
	ErrCodeAuthFailed     = int32(protoerrors.CommonErrorCode_ERR_AUTH_FAILED)
	ErrCodePlayerNotFound = int32(protoerrors.CommonErrorCode_ERR_PLAYER_NOT_FOUND)
	ErrCodeBattleNotFound = int32(protoerrors.CommonErrorCode_ERR_BATTLE_NOT_FOUND)
	ErrCodeUnknownMessage = int32(protoerrors.CommonErrorCode_ERR_UNKNOWN_MESSAGE)
	ErrCodeServerBusy     = int32(protoerrors.CommonErrorCode_ERR_SERVER_BUSY)
	ErrCodeInvalidPlayer  = int32(protoerrors.CommonErrorCode_ERR_INVALID_PLAYER)
	ErrCodeUnknown        = int32(protoerrors.CommonErrorCode_ERR_UNKNOWN)
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
