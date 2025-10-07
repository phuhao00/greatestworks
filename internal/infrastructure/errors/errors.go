package errors

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"greatestworks/internal/infrastructure/logging"

	"github.com/gin-gonic/gin"
)

// ErrorCode 错误码类型
type ErrorCode int

// 预定义错误码
const (
	// 通用错误码(1000-1999)
	ErrUnknown ErrorCode = 1000 + iota
	ErrInternal
	ErrInvalidInput
	ErrNotFound
	ErrUnauthorized
	ErrForbidden
	ErrConflict
	ErrTimeout
	ErrRateLimit
	ErrServiceUnavailable

	// 认证相关错误码(2000-2999)
	ErrAuthTokenMissing
	ErrAuthTokenInvalid
	ErrAuthTokenExpired
	ErrAuthUserNotFound
	ErrAuthPasswordIncorrect
	ErrAuthUserAlreadyExists
	ErrAuthUserDisabled

	// 玩家相关错误码(3000-3999)
	ErrPlayerNotFound
	ErrPlayerOffline
	ErrPlayerAlreadyExists
	ErrPlayerInvalidName
	ErrPlayerInvalidLevel
	ErrPlayerInsufficientExp
	ErrPlayerDead
	ErrPlayerInvalidPosition
	ErrPlayerVersionMismatch

	// 战斗相关错误码(4000-4999)
	ErrBattleNotFound
	ErrBattleAlreadyStarted
	ErrBattleNotInProgress
	ErrPlayerNotInBattle
	ErrPlayerAlreadyInBattle
	ErrInsufficientParticipants
	ErrPlayerDeadInBattle
	ErrInvalidAction
	ErrActionOnCooldown
	ErrInsufficientMana
	ErrInvalidTarget
	ErrBattleFinished
	ErrBattleAlreadyFinished
	ErrBattleNotFinished

	// 数据库相关错误码(5000-5999)
	ErrDatabaseConnection
	ErrDatabaseQuery
	ErrDatabaseTransaction
	ErrDatabaseConstraint
	ErrDatabaseTimeout

	// 网络相关错误码(6000-6999)
	ErrNetworkConnection
	ErrNetworkTimeout
	ErrNetworkUnreachable
	ErrNetworkInvalidResponse
)

// Error 自定义错误结构
type Error struct {
	Code      ErrorCode `json:"code"`
	Message   string    `json:"message"`
	Details   string    `json:"details,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	File      string    `json:"file,omitempty"`
	Line      int       `json:"line,omitempty"`
	Stack     string    `json:"stack,omitempty"`
}

// Error 实现error接口
func (e *Error) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// NewError 创建新错误
func NewError(code ErrorCode, message string) *Error {
	_, file, line, _ := runtime.Caller(1)
	return &Error{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		File:      file,
		Line:      line,
	}
}

// NewErrorWithDetails 创建带详情的错误
func NewErrorWithDetails(code ErrorCode, message, details string) *Error {
	_, file, line, _ := runtime.Caller(1)
	return &Error{
		Code:      code,
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
		File:      file,
		Line:      line,
	}
}

// NewErrorWithStack 创建带堆栈的错误
func NewErrorWithStack(code ErrorCode, message string) *Error {
	_, file, line, _ := runtime.Caller(1)
	stack := make([]byte, 1024)
	length := runtime.Stack(stack, false)

	return &Error{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		File:      file,
		Line:      line,
		Stack:     string(stack[:length]),
	}
}

// GetHTTPStatus 获取HTTP状态码
func (e *Error) GetHTTPStatus() int {
	switch e.Code {
	case ErrNotFound, ErrPlayerNotFound, ErrBattleNotFound:
		return http.StatusNotFound
	case ErrUnauthorized, ErrAuthTokenMissing, ErrAuthTokenInvalid, ErrAuthTokenExpired:
		return http.StatusUnauthorized
	case ErrForbidden, ErrAuthUserDisabled:
		return http.StatusForbidden
	case ErrConflict, ErrPlayerAlreadyExists, ErrPlayerAlreadyInBattle:
		return http.StatusConflict
	case ErrInvalidInput, ErrPlayerInvalidName, ErrPlayerInvalidLevel, ErrPlayerInvalidPosition:
		return http.StatusBadRequest
	case ErrTimeout, ErrDatabaseTimeout, ErrNetworkTimeout:
		return http.StatusRequestTimeout
	case ErrRateLimit:
		return http.StatusTooManyRequests
	case ErrServiceUnavailable, ErrDatabaseConnection, ErrNetworkConnection:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// 预定义错误
var (
	// 通用错误
	ErrUnknownError            = NewError(ErrUnknown, "未知错误")
	ErrInternalError           = NewError(ErrInternal, "内部错误")
	ErrInvalidInputError       = NewError(ErrInvalidInput, "无效输入")
	ErrNotFoundError           = NewError(ErrNotFound, "资源不存在")
	ErrUnauthorizedError       = NewError(ErrUnauthorized, "未授权")
	ErrForbiddenError          = NewError(ErrForbidden, "禁止访问")
	ErrConflictError           = NewError(ErrConflict, "冲突")
	ErrTimeoutError            = NewError(ErrTimeout, "超时")
	ErrRateLimitError          = NewError(ErrRateLimit, "请求过于频繁")
	ErrServiceUnavailableError = NewError(ErrServiceUnavailable, "服务不可用")

	// 认证相关错误
	ErrAuthTokenMissingError      = NewError(ErrAuthTokenMissing, "缺少认证令牌")
	ErrAuthTokenInvalidError      = NewError(ErrAuthTokenInvalid, "无效的认证令牌")
	ErrAuthTokenExpiredError      = NewError(ErrAuthTokenExpired, "认证令牌已过期")
	ErrAuthUserNotFoundError      = NewError(ErrAuthUserNotFound, "用户不存在")
	ErrAuthPasswordIncorrectError = NewError(ErrAuthPasswordIncorrect, "密码错误")
	ErrAuthUserAlreadyExistsError = NewError(ErrAuthUserAlreadyExists, "用户已存在")
	ErrAuthUserDisabledError      = NewError(ErrAuthUserDisabled, "用户已被禁用")

	// 玩家相关错误
	ErrPlayerNotFoundError        = NewError(ErrPlayerNotFound, "玩家不存在")
	ErrPlayerOfflineError         = NewError(ErrPlayerOffline, "玩家已离线")
	ErrPlayerAlreadyExistsError   = NewError(ErrPlayerAlreadyExists, "玩家已存在")
	ErrPlayerInvalidNameError     = NewError(ErrPlayerInvalidName, "无效的玩家名称")
	ErrPlayerInvalidLevelError    = NewError(ErrPlayerInvalidLevel, "无效的等级")
	ErrPlayerInsufficientExpError = NewError(ErrPlayerInsufficientExp, "经验值不足")
	ErrPlayerDeadError            = NewError(ErrPlayerDead, "玩家已死亡")
	ErrPlayerInvalidPositionError = NewError(ErrPlayerInvalidPosition, "无效的位置")
	ErrPlayerVersionMismatchError = NewError(ErrPlayerVersionMismatch, "版本不匹配")

	// 战斗相关错误
	ErrBattleNotFoundError           = NewError(ErrBattleNotFound, "战斗不存在")
	ErrBattleAlreadyStartedError     = NewError(ErrBattleAlreadyStarted, "战斗已开始")
	ErrBattleNotInProgressError      = NewError(ErrBattleNotInProgress, "战斗未进行中")
	ErrPlayerNotInBattleError        = NewError(ErrPlayerNotInBattle, "玩家不在战斗中")
	ErrPlayerAlreadyInBattleError    = NewError(ErrPlayerAlreadyInBattle, "玩家已在战斗中")
	ErrInsufficientParticipantsError = NewError(ErrInsufficientParticipants, "参与者不足")
	ErrInvalidActionError            = NewError(ErrInvalidAction, "无效的行动")
	ErrActionOnCooldownError         = NewError(ErrActionOnCooldown, "行动冷却中")
	ErrInsufficientManaError         = NewError(ErrInsufficientMana, "魔法值不足")
	ErrInvalidTargetError            = NewError(ErrInvalidTarget, "无效的目标")
	ErrBattleFinishedError           = NewError(ErrBattleFinished, "战斗已结束")
	ErrBattleAlreadyFinishedError    = NewError(ErrBattleAlreadyFinished, "战斗已结束")
	ErrBattleNotFinishedError        = NewError(ErrBattleNotFinished, "战斗未结束")

	// 数据库相关错误
	ErrDatabaseConnectionError  = NewError(ErrDatabaseConnection, "数据库连接失败")
	ErrDatabaseQueryError       = NewError(ErrDatabaseQuery, "数据库查询失败")
	ErrDatabaseTransactionError = NewError(ErrDatabaseTransaction, "数据库事务失败")
	ErrDatabaseConstraintError  = NewError(ErrDatabaseConstraint, "数据库约束违反")
	ErrDatabaseTimeoutError     = NewError(ErrDatabaseTimeout, "数据库操作超时")

	// 网络相关错误
	ErrNetworkConnectionError      = NewError(ErrNetworkConnection, "网络连接失败")
	ErrNetworkTimeoutError         = NewError(ErrNetworkTimeout, "网络超时")
	ErrNetworkUnreachableError     = NewError(ErrNetworkUnreachable, "网络不可达")
	ErrNetworkInvalidResponseError = NewError(ErrNetworkInvalidResponse, "无效的网络响应")
)

// ErrorHandler 错误处理器
type ErrorHandler struct {
	logger logging.Logger
}

// NewErrorHandler 创建错误处理器
func NewErrorHandler(logger logging.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// HandleError 处理错误
func (h *ErrorHandler) HandleError(c *gin.Context, err error) {
	// 记录错误日志
	h.logger.Error("Request error", err, logging.Fields{
		"path":   c.Request.URL.Path,
		"method": c.Request.Method,
	})

	// 检查是否为自定义错误
	if customErr, ok := err.(*Error); ok {
		c.JSON(customErr.GetHTTPStatus(), gin.H{
			"error": gin.H{
				"code":      customErr.Code,
				"message":   customErr.Message,
				"details":   customErr.Details,
				"timestamp": customErr.Timestamp,
			},
		})
		return
	}

	// 处理其他错误
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"code":      ErrInternal,
			"message":   "内部错误",
			"timestamp": time.Now(),
		},
	})
}

// HandlePanic 处理panic
func (h *ErrorHandler) HandlePanic(c *gin.Context, recovered interface{}) {
	// 记录panic日志
	h.logger.Error("Panic recovered", fmt.Errorf("panic: %v", recovered), logging.Fields{
		"path":   c.Request.URL.Path,
		"method": c.Request.Method,
	})

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"code":      ErrInternal,
			"message":   "内部错误",
			"timestamp": time.Now(),
		},
	})
}

// IsError 检查是否为特定错误
func IsError(err error, code ErrorCode) bool {
	if customErr, ok := err.(*Error); ok {
		return customErr.Code == code
	}
	return false
}

// GetErrorCode 获取错误码
func GetErrorCode(err error) ErrorCode {
	if customErr, ok := err.(*Error); ok {
		return customErr.Code
	}
	return ErrUnknown
}
