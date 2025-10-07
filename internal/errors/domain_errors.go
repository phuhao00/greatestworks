// Package errors 提供统一的错误处理机制
package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误代码类型
type ErrorCode string

const (
	// 通用错误
	ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
	ErrCodeInvalidInput ErrorCode = "INVALID_INPUT"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeConflict     ErrorCode = "CONFLICT"
	ErrCodeTimeout      ErrorCode = "TIMEOUT"

	// 领域特定错误
	ErrCodePlayerNotFound      ErrorCode = "PLAYER_NOT_FOUND"
	ErrCodePlayerOffline       ErrorCode = "PLAYER_OFFLINE"
	ErrCodePlayerAlreadyExists ErrorCode = "PLAYER_ALREADY_EXISTS"
	ErrCodeInvalidPlayerName   ErrorCode = "INVALID_PLAYER_NAME"
	ErrCodeInvalidPosition     ErrorCode = "INVALID_POSITION"
	ErrCodeVersionMismatch     ErrorCode = "VERSION_MISMATCH"

	ErrCodeBattleNotFound           ErrorCode = "BATTLE_NOT_FOUND"
	ErrCodeBattleAlreadyStarted     ErrorCode = "BATTLE_ALREADY_STARTED"
	ErrCodeBattleNotInProgress      ErrorCode = "BATTLE_NOT_IN_PROGRESS"
	ErrCodePlayerNotInBattle        ErrorCode = "PLAYER_NOT_IN_BATTLE"
	ErrCodePlayerAlreadyInBattle    ErrorCode = "PLAYER_ALREADY_IN_BATTLE"
	ErrCodeInsufficientParticipants ErrorCode = "INSUFFICIENT_PARTICIPANTS"
	ErrCodeInsufficientMana         ErrorCode = "INSUFFICIENT_MANA"
	ErrCodeInvalidTarget            ErrorCode = "INVALID_TARGET"
	ErrCodeBattleFinished           ErrorCode = "BATTLE_FINISHED"
)

// DomainError 领域错误结构
type DomainError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	HTTPStatus int       `json:"-"`
	Cause      error     `json:"-"`
}

// Error 实现 error 接口
func (e *DomainError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 支持错误链
func (e *DomainError) Unwrap() error {
	return e.Cause
}

// NewDomainError 创建领域错误
func NewDomainError(code ErrorCode, message string) *DomainError {
	return &DomainError{
		Code:       code,
		Message:    message,
		HTTPStatus: getHTTPStatus(code),
	}
}

// NewDomainErrorWithDetails 创建带详情的领域错误
func NewDomainErrorWithDetails(code ErrorCode, message, details string) *DomainError {
	return &DomainError{
		Code:       code,
		Message:    message,
		Details:    details,
		HTTPStatus: getHTTPStatus(code),
	}
}

// NewDomainErrorWithCause 创建带原因的领域错误
func NewDomainErrorWithCause(code ErrorCode, message string, cause error) *DomainError {
	return &DomainError{
		Code:       code,
		Message:    message,
		HTTPStatus: getHTTPStatus(code),
		Cause:      cause,
	}
}

// getHTTPStatus 根据错误代码获取HTTP状态码
func getHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrCodeNotFound, ErrCodePlayerNotFound, ErrCodeBattleNotFound:
		return http.StatusNotFound
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeConflict, ErrCodePlayerAlreadyExists, ErrCodePlayerAlreadyInBattle:
		return http.StatusConflict
	case ErrCodeInvalidInput, ErrCodeInvalidPlayerName, ErrCodeInvalidPosition, ErrCodeInvalidTarget:
		return http.StatusBadRequest
	case ErrCodeTimeout:
		return http.StatusRequestTimeout
	default:
		return http.StatusInternalServerError
	}
}

// 预定义的领域错误
var (
	// 玩家相关错误
	ErrPlayerNotFound      = NewDomainError(ErrCodePlayerNotFound, "玩家不存在")
	ErrPlayerOffline       = NewDomainError(ErrCodePlayerOffline, "玩家已离线")
	ErrPlayerAlreadyExists = NewDomainError(ErrCodePlayerAlreadyExists, "玩家已存在")
	ErrInvalidPlayerName   = NewDomainError(ErrCodeInvalidPlayerName, "无效的玩家名称")
	ErrInvalidPosition     = NewDomainError(ErrCodeInvalidPosition, "无效的位置")
	ErrVersionMismatch     = NewDomainError(ErrCodeVersionMismatch, "版本不匹配")

	// 战斗相关错误
	ErrBattleNotFound           = NewDomainError(ErrCodeBattleNotFound, "战斗不存在")
	ErrBattleAlreadyStarted     = NewDomainError(ErrCodeBattleAlreadyStarted, "战斗已开始")
	ErrBattleNotInProgress      = NewDomainError(ErrCodeBattleNotInProgress, "战斗未进行中")
	ErrPlayerNotInBattle        = NewDomainError(ErrCodePlayerNotInBattle, "玩家不在战斗中")
	ErrPlayerAlreadyInBattle    = NewDomainError(ErrCodePlayerAlreadyInBattle, "玩家已在战斗中")
	ErrInsufficientParticipants = NewDomainError(ErrCodeInsufficientParticipants, "参与者不足")
	ErrInsufficientMana         = NewDomainError(ErrCodeInsufficientMana, "魔法值不足")
	ErrInvalidTarget            = NewDomainError(ErrCodeInvalidTarget, "无效的目标")
	ErrBattleFinished           = NewDomainError(ErrCodeBattleFinished, "战斗已结束")

	// 通用错误
	ErrInternal     = NewDomainError(ErrCodeInternal, "内部错误")
	ErrInvalidInput = NewDomainError(ErrCodeInvalidInput, "无效输入")
	ErrNotFound     = NewDomainError(ErrCodeNotFound, "资源不存在")
	ErrUnauthorized = NewDomainError(ErrCodeUnauthorized, "未授权")
	ErrForbidden    = NewDomainError(ErrCodeForbidden, "禁止访问")
	ErrConflict     = NewDomainError(ErrCodeConflict, "冲突")
	ErrTimeout      = NewDomainError(ErrCodeTimeout, "超时")
)

// IsDomainError 检查是否为领域错误
func IsDomainError(err error) bool {
	_, ok := err.(*DomainError)
	return ok
}

// GetDomainError 获取领域错误
func GetDomainError(err error) (*DomainError, bool) {
	domainErr, ok := err.(*DomainError)
	return domainErr, ok
}

// WrapError 包装错误为领域错误
func WrapError(err error, code ErrorCode, message string) *DomainError {
	return NewDomainErrorWithCause(code, message, err)
}
