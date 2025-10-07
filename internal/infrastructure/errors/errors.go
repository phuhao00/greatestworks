package errors

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"

	"greatestworks/internal/infrastructure/logging"
)

// ErrorCode 错误码类�?
type ErrorCode int

// 预定义错误码
const (
	// 通用错误�?(1000-1999)
	ErrUnknown ErrorCode = 1000 + iota
	ErrInternalServer
	ErrInvalidRequest
	ErrInvalidParameter
	ErrMissingParameter
	ErrValidationFailed
	ErrTimeout
	ErrServiceUnavailable
	ErrRateLimitExceeded
	ErrMaintenanceMode

	// 认证错误�?(2000-2999)
	ErrUnauthorized ErrorCode = 2000 + iota
	ErrInvalidToken
	ErrTokenExpired
	ErrTokenRevoked
	ErrInvalidCredentials
	ErrAccountLocked
	ErrAccountDisabled
	ErrPermissionDenied
	ErrInsufficientPrivileges
	ErrSessionExpired

	// 用户错误�?(3000-3999)
	ErrUserNotFound ErrorCode = 3000 + iota
	ErrUserAlreadyExists
	ErrInvalidUserData
	ErrUsernameTaken
	ErrEmailTaken
	ErrWeakPassword
	ErrPasswordMismatch
	ErrUserBanned
	ErrUserSuspended
	ErrEmailNotVerified

	// 玩家错误�?(4000-4999)
	ErrPlayerNotFound ErrorCode = 4000 + iota
	ErrPlayerAlreadyExists
	ErrInvalidPlayerData
	ErrPlayerOffline
	ErrPlayerInBattle
	ErrPlayerBusy
	ErrInsufficientLevel
	ErrInsufficientExp
	ErrInsufficientGold
	ErrInventoryFull

	// 战斗错误�?(5000-5999)
	ErrBattleNotFound ErrorCode = 5000 + iota
	ErrBattleAlreadyExists
	ErrBattleFull
	ErrBattleStarted
	ErrBattleEnded
	ErrInvalidBattleAction
	ErrNotInBattle
	ErrBattleTimeout
	ErrInvalidTarget
	ErrActionCooldown

	// 数据库错误码 (6000-6999)
	ErrDatabaseConnection ErrorCode = 6000 + iota
	ErrDatabaseQuery
	ErrDatabaseTransaction
	ErrRecordNotFound
	ErrRecordAlreadyExists
	ErrConstraintViolation
	ErrDatabaseTimeout
	ErrDatabaseLock
	ErrMigrationFailed
	ErrBackupFailed

	// 网络错误�?(7000-7999)
	ErrConnectionFailed ErrorCode = 7000 + iota
	ErrConnectionTimeout
	ErrConnectionLost
	ErrInvalidProtocol
	ErrMessageTooLarge
	ErrInvalidMessage
	ErrNetworkUnavailable
	ErrDNSResolution
	ErrSSLHandshake
	ErrProxyError
)

// AppError 应用错误
type AppError struct {
	Code      ErrorCode              `json:"code"`
	Message   string                 `json:"message"`
	Details   string                 `json:"details,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Cause     error                  `json:"-"`
	Stack     string                 `json:"-"`
	Timestamp time.Time              `json:"timestamp"`
	RequestID string                 `json:"request_id,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 实现errors.Unwrap接口
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithCause 添加原因错误
func (e *AppError) WithCause(cause error) *AppError {
	e.Cause = cause
	return e
}

// WithDetails 添加详细信息
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithMetadata 添加元数�?
func (e *AppError) WithMetadata(key string, value interface{}) *AppError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
	return e
}

// WithRequestID 添加请求ID
func (e *AppError) WithRequestID(requestID string) *AppError {
	e.RequestID = requestID
	return e
}

// WithUserID 添加用户ID
func (e *AppError) WithUserID(userID string) *AppError {
	e.UserID = userID
	return e
}

// GetHTTPStatus 获取HTTP状态码
func (e *AppError) GetHTTPStatus() int {
	switch {
	case e.Code >= 2000 && e.Code < 3000:
		return http.StatusUnauthorized
	case e.Code == ErrPermissionDenied || e.Code == ErrInsufficientPrivileges:
		return http.StatusForbidden
	case e.Code >= 3000 && e.Code < 4000:
		if e.Code == ErrUserNotFound {
			return http.StatusNotFound
		}
		if e.Code == ErrUserAlreadyExists || e.Code == ErrUsernameTaken || e.Code == ErrEmailTaken {
			return http.StatusConflict
		}
		return http.StatusBadRequest
	case e.Code >= 4000 && e.Code < 5000:
		if e.Code == ErrPlayerNotFound {
			return http.StatusNotFound
		}
		return http.StatusBadRequest
	case e.Code >= 5000 && e.Code < 6000:
		if e.Code == ErrBattleNotFound {
			return http.StatusNotFound
		}
		return http.StatusBadRequest
	case e.Code == ErrRateLimitExceeded:
		return http.StatusTooManyRequests
	case e.Code == ErrTimeout || e.Code == ErrConnectionTimeout:
		return http.StatusRequestTimeout
	case e.Code == ErrServiceUnavailable || e.Code == ErrMaintenanceMode:
		return http.StatusServiceUnavailable
	case e.Code == ErrInvalidRequest || e.Code == ErrInvalidParameter || e.Code == ErrMissingParameter || e.Code == ErrValidationFailed:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// NewError 创建新的应用错误
func NewError(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Stack:     getStack(),
	}
}

// NewErrorWithCause 创建带原因的应用错误
func NewErrorWithCause(code ErrorCode, message string, cause error) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Cause:     cause,
		Timestamp: time.Now(),
		Stack:     getStack(),
	}
}

// WrapError 包装现有错误
func WrapError(err error, code ErrorCode, message string) *AppError {
	if err == nil {
		return nil
	}

	// 如果已经是AppError，更新信�?
	if appErr, ok := err.(*AppError); ok {
		appErr.Code = code
		appErr.Message = message
		return appErr
	}

	return &AppError{
		Code:      code,
		Message:   message,
		Cause:     err,
		Timestamp: time.Now(),
		Stack:     getStack(),
	}
}

// getStack 获取调用�?
func getStack() string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// 预定义错�?
var (
	// 通用错误
	ErrInternal         = NewError(ErrInternalServer, "Internal server error")
	ErrBadRequest       = NewError(ErrInvalidRequest, "Bad request")
	ErrNotFound         = NewError(ErrRecordNotFound, "Resource not found")
	ErrConflict         = NewError(ErrRecordAlreadyExists, "Resource already exists")
	ErrForbidden        = NewError(ErrPermissionDenied, "Access forbidden")
	ErrUnauth           = NewError(ErrUnauthorized, "Unauthorized")
	ErrRateLimit        = NewError(ErrRateLimitExceeded, "Rate limit exceeded")
	ErrMaintenance      = NewError(ErrMaintenanceMode, "Service under maintenance")

	// 认证错误
	ErrInvalidAuth      = NewError(ErrInvalidCredentials, "Invalid credentials")
	ErrExpiredToken     = NewError(ErrTokenExpired, "Token expired")
	ErrInvalidTokenErr  = NewError(ErrInvalidToken, "Invalid token")
	ErrSessionExp       = NewError(ErrSessionExpired, "Session expired")

	// 用户错误
	ErrUserNotFoundErr  = NewError(ErrUserNotFound, "User not found")
	ErrUserExists       = NewError(ErrUserAlreadyExists, "User already exists")
	ErrUsernameExists   = NewError(ErrUsernameTaken, "Username already taken")
	ErrEmailExists      = NewError(ErrEmailTaken, "Email already taken")
	ErrPasswordWeak     = NewError(ErrWeakPassword, "Password too weak")

	// 玩家错误
	ErrPlayerNotFoundErr = NewError(ErrPlayerNotFound, "Player not found")
	ErrPlayerExists      = NewError(ErrPlayerAlreadyExists, "Player already exists")
	ErrPlayerOfflineErr  = NewError(ErrPlayerOffline, "Player is offline")
	ErrPlayerInBattleErr = NewError(ErrPlayerInBattle, "Player is in battle")

	// 战斗错误
	ErrBattleNotFoundErr = NewError(ErrBattleNotFound, "Battle not found")
	ErrBattleFullErr     = NewError(ErrBattleFull, "Battle is full")
	ErrBattleStartedErr  = NewError(ErrBattleStarted, "Battle already started")
	ErrNotInBattleErr    = NewError(ErrNotInBattle, "Not in battle")
)

// ErrorHandler 错误处理�?
type ErrorHandler struct {
	logger logger.Logger
}

// NewErrorHandler 创建错误处理�?
func NewErrorHandler(logger logger.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// HandleError HTTP错误处理中间�?
func (h *ErrorHandler) HandleError() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		var err error
		switch x := recovered.(type) {
		case string:
			err = fmt.Errorf("%s", x)
		case error:
			err = x
		default:
			err = fmt.Errorf("unknown panic: %v", x)
		}

		// 记录panic错误
		h.logger.Error("Panic recovered", 
			"error", err,
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"stack", getStack())

		// 创建内部服务器错�?
		appErr := NewErrorWithCause(ErrInternalServer, "Internal server error", err)
		h.HandleAppError(c, appErr)
	})
}

// HandleAppError 处理应用错误
func (h *ErrorHandler) HandleAppError(c *gin.Context, err *AppError) {
	// 添加请求信息
	if err.RequestID == "" {
		if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
			err.RequestID = requestID
		}
	}

	if err.UserID == "" {
		if userID, exists := c.Get("user_id"); exists {
			if userIDStr, ok := userID.(string); ok {
				err.UserID = userIDStr
			}
		}
	}

	// 记录错误
	h.logError(c, err)

	// 构造响�?
	response := h.buildErrorResponse(err)

	// 返回错误响应
	c.JSON(err.GetHTTPStatus(), response)
	c.Abort()
}

// logError 记录错误
func (h *ErrorHandler) logError(c *gin.Context, err *AppError) {
	logFields := []interface{}{
		"error_code", err.Code,
		"error_message", err.Message,
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
		"request_id", err.RequestID,
		"user_id", err.UserID,
	}

	if err.Details != "" {
		logFields = append(logFields, "details", err.Details)
	}

	if err.Metadata != nil {
		logFields = append(logFields, "metadata", err.Metadata)
	}

	if err.Cause != nil {
		logFields = append(logFields, "cause", err.Cause.Error())
	}

	// 根据错误级别记录日志
	if err.Code >= 6000 || err.Code == ErrInternalServer {
		// 系统级错�?
		logFields = append(logFields, "stack", err.Stack)
		h.logger.Error("System error occurred", logFields...)
	} else if err.Code >= 2000 && err.Code < 3000 {
		// 认证错误
		h.logger.Warn("Authentication error", logFields...)
	} else {
		// 业务错误
		h.logger.Info("Business error", logFields...)
	}
}

// buildErrorResponse 构造错误响�?
func (h *ErrorHandler) buildErrorResponse(err *AppError) map[string]interface{} {
	response := map[string]interface{}{
		"success":   false,
		"error":     true,
		"code":      err.Code,
		"message":   err.Message,
		"timestamp": err.Timestamp.Unix(),
	}

	if err.RequestID != "" {
		response["request_id"] = err.RequestID
	}

	// 在开发环境下返回更多信息
	if gin.Mode() == gin.DebugMode {
		if err.Details != "" {
			response["details"] = err.Details
		}
		if err.Metadata != nil {
			response["metadata"] = err.Metadata
		}
		if err.Cause != nil {
			response["cause"] = err.Cause.Error()
		}
	}

	return response
}

// IsAppError 检查是否是应用错误
func IsAppError(err error) (*AppError, bool) {
	if err == nil {
		return nil, false
	}
	appErr, ok := err.(*AppError)
	return appErr, ok
}

// GetErrorCode 获取错误�?
func GetErrorCode(err error) ErrorCode {
	if appErr, ok := IsAppError(err); ok {
		return appErr.Code
	}
	return ErrUnknown
}

// IsErrorCode 检查是否是指定错误�?
func IsErrorCode(err error, code ErrorCode) bool {
	return GetErrorCode(err) == code
}