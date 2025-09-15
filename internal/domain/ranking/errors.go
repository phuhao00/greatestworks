package ranking

import (
	"fmt"
)

// 排行榜领域错误定义

// RankingError 排行榜错误基础接口
type RankingError interface {
	error
	GetCode() string
	GetMessage() string
	GetDetails() map[string]interface{}
	IsRetryable() bool
	GetSeverity() ErrorSeverity
}

// ErrorSeverity 错误严重程度
type ErrorSeverity int

const (
	ErrorSeverityLow ErrorSeverity = iota
	ErrorSeverityMedium
	ErrorSeverityHigh
	ErrorSeverityCritical
)

// String 返回错误严重程度的字符串表示
func (s ErrorSeverity) String() string {
	switch s {
	case ErrorSeverityLow:
		return "low"
	case ErrorSeverityMedium:
		return "medium"
	case ErrorSeverityHigh:
		return "high"
	case ErrorSeverityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// BaseRankingError 排行榜错误基础结构
type BaseRankingError struct {
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details"`
	Retryable bool                   `json:"retryable"`
	Severity  ErrorSeverity          `json:"severity"`
}

// Error 实现error接口
func (e *BaseRankingError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// GetCode 获取错误代码
func (e *BaseRankingError) GetCode() string {
	return e.Code
}

// GetMessage 获取错误消息
func (e *BaseRankingError) GetMessage() string {
	return e.Message
}

// GetDetails 获取错误详情
func (e *BaseRankingError) GetDetails() map[string]interface{} {
	return e.Details
}

// IsRetryable 是否可重试
func (e *BaseRankingError) IsRetryable() bool {
	return e.Retryable
}

// GetSeverity 获取错误严重程度
func (e *BaseRankingError) GetSeverity() ErrorSeverity {
	return e.Severity
}

// 排行榜相关错误

// RankingNotFoundError 排行榜未找到错误
type RankingNotFoundError struct {
	*BaseRankingError
	RankID uint32 `json:"rank_id"`
}

// NewRankingNotFoundError 创建排行榜未找到错误
func NewRankingNotFoundError(rankID uint32) *RankingNotFoundError {
	return &RankingNotFoundError{
		BaseRankingError: &BaseRankingError{
			Code:      "RANKING_NOT_FOUND",
			Message:   fmt.Sprintf("Ranking with ID %d not found", rankID),
			Details:   map[string]interface{}{"rank_id": rankID},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		RankID: rankID,
	}
}

// RankingAlreadyExistsError 排行榜已存在错误
type RankingAlreadyExistsError struct {
	*BaseRankingError
	RankID uint32 `json:"rank_id"`
}

// NewRankingAlreadyExistsError 创建排行榜已存在错误
func NewRankingAlreadyExistsError(rankID uint32) *RankingAlreadyExistsError {
	return &RankingAlreadyExistsError{
		BaseRankingError: &BaseRankingError{
			Code:      "RANKING_ALREADY_EXISTS",
			Message:   fmt.Sprintf("Ranking with ID %d already exists", rankID),
			Details:   map[string]interface{}{"rank_id": rankID},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		RankID: rankID,
	}
}

// RankingInactiveError 排行榜非活跃错误
type RankingInactiveError struct {
	*BaseRankingError
	RankID uint32     `json:"rank_id"`
	Status RankStatus `json:"status"`
}

// NewRankingInactiveError 创建排行榜非活跃错误
func NewRankingInactiveError(rankID uint32) *RankingInactiveError {
	return &RankingInactiveError{
		BaseRankingError: &BaseRankingError{
			Code:      "RANKING_INACTIVE",
			Message:   fmt.Sprintf("Ranking %d is not active", rankID),
			Details:   map[string]interface{}{"rank_id": rankID},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		RankID: rankID,
	}
}

// RankingTimeExpiredError 排行榜时间过期错误
type RankingTimeExpiredError struct {
	*BaseRankingError
	RankID    uint32 `json:"rank_id"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}

// NewRankingTimeExpiredError 创建排行榜时间过期错误
func NewRankingTimeExpiredError(rankID uint32, startTime, endTime int64) *RankingTimeExpiredError {
	return &RankingTimeExpiredError{
		BaseRankingError: &BaseRankingError{
			Code:      "RANKING_TIME_EXPIRED",
			Message:   fmt.Sprintf("Ranking %d is outside valid time range", rankID),
			Details: map[string]interface{}{
				"rank_id":    rankID,
				"start_time": startTime,
				"end_time":   endTime,
			},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		RankID:    rankID,
		StartTime: startTime,
		EndTime:   endTime,
	}
}

// RankingFullError 排行榜已满错误
type RankingFullError struct {
	*BaseRankingError
	RankID      uint32 `json:"rank_id"`
	MaxSize     int64  `json:"max_size"`
	CurrentSize int64  `json:"current_size"`
}

// NewRankingFullError 创建排行榜已满错误
func NewRankingFullError(rankID uint32, maxSize, currentSize int64) *RankingFullError {
	return &RankingFullError{
		BaseRankingError: &BaseRankingError{
			Code:      "RANKING_FULL",
			Message:   fmt.Sprintf("Ranking %d is full (%d/%d)", rankID, currentSize, maxSize),
			Details: map[string]interface{}{
				"rank_id":      rankID,
				"max_size":     maxSize,
				"current_size": currentSize,
			},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		RankID:      rankID,
		MaxSize:     maxSize,
		CurrentSize: currentSize,
	}
}

// 玩家相关错误

// PlayerNotInRankingError 玩家不在排行榜错误
type PlayerNotInRankingError struct {
	*BaseRankingError
	PlayerID uint64 `json:"player_id"`
	RankID   uint32 `json:"rank_id"`
}

// NewPlayerNotInRankingError 创建玩家不在排行榜错误
func NewPlayerNotInRankingError(playerID uint64, rankID uint32) *PlayerNotInRankingError {
	return &PlayerNotInRankingError{
		BaseRankingError: &BaseRankingError{
			Code:      "PLAYER_NOT_IN_RANKING",
			Message:   fmt.Sprintf("Player %d not found in ranking %d", playerID, rankID),
			Details:   map[string]interface{}{"player_id": playerID, "rank_id": rankID},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		PlayerID: playerID,
		RankID:   rankID,
	}
}

// PlayerAlreadyInRankingError 玩家已在排行榜错误
type PlayerAlreadyInRankingError struct {
	*BaseRankingError
	PlayerID    uint64 `json:"player_id"`
	RankID      uint32 `json:"rank_id"`
	CurrentRank int64  `json:"current_rank"`
}

// NewPlayerAlreadyInRankingError 创建玩家已在排行榜错误
func NewPlayerAlreadyInRankingError(playerID uint64, rankID uint32, currentRank int64) *PlayerAlreadyInRankingError {
	return &PlayerAlreadyInRankingError{
		BaseRankingError: &BaseRankingError{
			Code:      "PLAYER_ALREADY_IN_RANKING",
			Message:   fmt.Sprintf("Player %d already exists in ranking %d at rank %d", playerID, rankID, currentRank),
			Details: map[string]interface{}{
				"player_id":    playerID,
				"rank_id":      rankID,
				"current_rank": currentRank,
			},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		PlayerID:    playerID,
		RankID:      rankID,
		CurrentRank: currentRank,
	}
}

// PlayerBlacklistedError 玩家被黑名单错误
type PlayerBlacklistedError struct {
	*BaseRankingError
	PlayerID uint64 `json:"player_id"`
	RankID   uint32 `json:"rank_id"`
	Reason   string `json:"reason"`
}

// NewPlayerBlacklistedError 创建玩家被黑名单错误
func NewPlayerBlacklistedError(playerID uint64, rankID uint32) *PlayerBlacklistedError {
	return &PlayerBlacklistedError{
		BaseRankingError: &BaseRankingError{
			Code:      "PLAYER_BLACKLISTED",
			Message:   fmt.Sprintf("Player %d is blacklisted in ranking %d", playerID, rankID),
			Details:   map[string]interface{}{"player_id": playerID, "rank_id": rankID},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		PlayerID: playerID,
		RankID:   rankID,
	}
}

// PlayerAlreadyBlacklistedError 玩家已被黑名单错误
type PlayerAlreadyBlacklistedError struct {
	*BaseRankingError
	PlayerID uint64 `json:"player_id"`
	RankID   uint32 `json:"rank_id"`
}

// NewPlayerAlreadyBlacklistedError 创建玩家已被黑名单错误
func NewPlayerAlreadyBlacklistedError(playerID uint64, rankID uint32) *PlayerAlreadyBlacklistedError {
	return &PlayerAlreadyBlacklistedError{
		BaseRankingError: &BaseRankingError{
			Code:      "PLAYER_ALREADY_BLACKLISTED",
			Message:   fmt.Sprintf("Player %d is already blacklisted in ranking %d", playerID, rankID),
			Details:   map[string]interface{}{"player_id": playerID, "rank_id": rankID},
			Retryable: false,
			Severity:  ErrorSeverityLow,
		},
		PlayerID: playerID,
		RankID:   rankID,
	}
}

// PlayerNotBlacklistedError 玩家未被黑名单错误
type PlayerNotBlacklistedError struct {
	*BaseRankingError
	PlayerID uint64 `json:"player_id"`
	RankID   uint32 `json:"rank_id"`
}

// NewPlayerNotBlacklistedError 创建玩家未被黑名单错误
func NewPlayerNotBlacklistedError(playerID uint64, rankID uint32) *PlayerNotBlacklistedError {
	return &PlayerNotBlacklistedError{
		BaseRankingError: &BaseRankingError{
			Code:      "PLAYER_NOT_BLACKLISTED",
			Message:   fmt.Sprintf("Player %d is not blacklisted in ranking %d", playerID, rankID),
			Details:   map[string]interface{}{"player_id": playerID, "rank_id": rankID},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		PlayerID: playerID,
		RankID:   rankID,
	}
}

// 范围和参数相关错误

// InvalidRangeError 无效范围错误
type InvalidRangeError struct {
	*BaseRankingError
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

// NewInvalidRangeError 创建无效范围错误
func NewInvalidRangeError(start, end int64) *InvalidRangeError {
	return &InvalidRangeError{
		BaseRankingError: &BaseRankingError{
			Code:      "INVALID_RANGE",
			Message:   fmt.Sprintf("Invalid range: start=%d, end=%d", start, end),
			Details:   map[string]interface{}{"start": start, "end": end},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		Start: start,
		End:   end,
	}
}

// InvalidTimeRangeError 无效时间范围错误
type InvalidTimeRangeError struct {
	*BaseRankingError
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

// NewInvalidTimeRangeError 创建无效时间范围错误
func NewInvalidTimeRangeError(startTime, endTime int64) *InvalidTimeRangeError {
	return &InvalidTimeRangeError{
		BaseRankingError: &BaseRankingError{
			Code:      "INVALID_TIME_RANGE",
			Message:   fmt.Sprintf("Invalid time range: start=%d, end=%d", startTime, endTime),
			Details:   map[string]interface{}{"start_time": startTime, "end_time": endTime},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		StartTime: startTime,
		EndTime:   endTime,
	}
}

// InvalidScoreError 无效分数错误
type InvalidScoreError struct {
	*BaseRankingError
	Score    int64  `json:"score"`
	MinScore *int64 `json:"min_score,omitempty"`
	MaxScore *int64 `json:"max_score,omitempty"`
}

// NewInvalidScoreError 创建无效分数错误
func NewInvalidScoreError(score int64, minScore, maxScore *int64) *InvalidScoreError {
	return &InvalidScoreError{
		BaseRankingError: &BaseRankingError{
			Code:      "INVALID_SCORE",
			Message:   fmt.Sprintf("Invalid score: %d", score),
			Details:   map[string]interface{}{"score": score, "min_score": minScore, "max_score": maxScore},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		Score:    score,
		MinScore: minScore,
		MaxScore: maxScore,
	}
}

// 权限相关错误

// InsufficientPermissionError 权限不足错误
type InsufficientPermissionError struct {
	*BaseRankingError
	UserID         string `json:"user_id"`
	RequiredPermission string `json:"required_permission"`
	Operation      string `json:"operation"`
}

// NewInsufficientPermissionError 创建权限不足错误
func NewInsufficientPermissionError(userID, operation, permission string) *InsufficientPermissionError {
	return &InsufficientPermissionError{
		BaseRankingError: &BaseRankingError{
			Code:      "INSUFFICIENT_PERMISSION",
			Message:   fmt.Sprintf("User %s lacks permission %s for operation %s", userID, permission, operation),
			Details: map[string]interface{}{
				"user_id":             userID,
				"required_permission": permission,
				"operation":           operation,
			},
			Retryable: false,
			Severity:  ErrorSeverityHigh,
		},
		UserID:             userID,
		RequiredPermission: permission,
		Operation:          operation,
	}
}

// OperationNotAllowedError 操作不允许错误
type OperationNotAllowedError struct {
	*BaseRankingError
	Operation string `json:"operation"`
	RankID    uint32 `json:"rank_id"`
	Reason    string `json:"reason"`
}

// NewOperationNotAllowedError 创建操作不允许错误
func NewOperationNotAllowedError(operation string, rankID uint32, reason string) *OperationNotAllowedError {
	return &OperationNotAllowedError{
		BaseRankingError: &BaseRankingError{
			Code:      "OPERATION_NOT_ALLOWED",
			Message:   fmt.Sprintf("Operation %s not allowed for ranking %d: %s", operation, rankID, reason),
			Details:   map[string]interface{}{"operation": operation, "rank_id": rankID, "reason": reason},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		Operation: operation,
		RankID:    rankID,
		Reason:    reason,
	}
}

// 系统错误

// RankingSystemError 排行榜系统错误
type RankingSystemError struct {
	*BaseRankingError
	SystemComponent string `json:"system_component"`
	InternalError   error  `json:"internal_error,omitempty"`
}

// NewRankingSystemError 创建排行榜系统错误
func NewRankingSystemError(component, message string, internalErr error) *RankingSystemError {
	return &RankingSystemError{
		BaseRankingError: &BaseRankingError{
			Code:      "RANKING_SYSTEM_ERROR",
			Message:   fmt.Sprintf("System error in %s: %s", component, message),
			Details:   map[string]interface{}{"system_component": component},
			Retryable: true,
			Severity:  ErrorSeverityCritical,
		},
		SystemComponent: component,
		InternalError:   internalErr,
	}
}

// RankingDatabaseError 排行榜数据库错误
type RankingDatabaseError struct {
	*BaseRankingError
	Operation     string `json:"operation"`
	Table         string `json:"table"`
	InternalError error  `json:"internal_error,omitempty"`
}

// NewRankingDatabaseError 创建排行榜数据库错误
func NewRankingDatabaseError(operation, table, message string, internalErr error) *RankingDatabaseError {
	return &RankingDatabaseError{
		BaseRankingError: &BaseRankingError{
			Code:      "RANKING_DATABASE_ERROR",
			Message:   fmt.Sprintf("Database error during %s on table %s: %s", operation, table, message),
			Details:   map[string]interface{}{"operation": operation, "table": table},
			Retryable: true,
			Severity:  ErrorSeverityHigh,
		},
		Operation:     operation,
		Table:         table,
		InternalError: internalErr,
	}
}

// RankingCacheError 排行榜缓存错误
type RankingCacheError struct {
	*BaseRankingError
	Operation     string `json:"operation"`
	Key           string `json:"key"`
	InternalError error  `json:"internal_error,omitempty"`
}

// NewRankingCacheError 创建排行榜缓存错误
func NewRankingCacheError(operation, key, message string, internalErr error) *RankingCacheError {
	return &RankingCacheError{
		BaseRankingError: &BaseRankingError{
			Code:      "RANKING_CACHE_ERROR",
			Message:   fmt.Sprintf("Cache error during %s for key %s: %s", operation, key, message),
			Details:   map[string]interface{}{"operation": operation, "key": key},
			Retryable: true,
			Severity:  ErrorSeverityMedium,
		},
		Operation:     operation,
		Key:           key,
		InternalError: internalErr,
	}
}

// 验证错误

// RankingValidationError 排行榜验证错误
type RankingValidationError struct {
	*BaseRankingError
	Field          string      `json:"field"`
	Value          interface{} `json:"value"`
	Constraint     string      `json:"constraint"`
	ValidationRule string      `json:"validation_rule"`
}

// NewRankingValidationError 创建排行榜验证错误
func NewRankingValidationError(field string, value interface{}, constraint, rule string) *RankingValidationError {
	return &RankingValidationError{
		BaseRankingError: &BaseRankingError{
			Code:      "RANKING_VALIDATION_ERROR",
			Message:   fmt.Sprintf("Validation failed for field %s: %s (rule: %s)", field, constraint, rule),
			Details: map[string]interface{}{
				"field":           field,
				"value":           value,
				"constraint":      constraint,
				"validation_rule": rule,
			},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		Field:          field,
		Value:          value,
		Constraint:     constraint,
		ValidationRule: rule,
	}
}

// 并发相关错误

// RankingConcurrencyError 排行榜并发错误
type RankingConcurrencyError struct {
	*BaseRankingError
	RankID          uint32 `json:"rank_id"`
	ExpectedVersion int64  `json:"expected_version"`
	ActualVersion   int64  `json:"actual_version"`
}

// NewRankingConcurrencyError 创建排行榜并发错误
func NewRankingConcurrencyError(rankID uint32, expectedVersion, actualVersion int64) *RankingConcurrencyError {
	return &RankingConcurrencyError{
		BaseRankingError: &BaseRankingError{
			Code:      "RANKING_CONCURRENCY_ERROR",
			Message:   fmt.Sprintf("Concurrency conflict for ranking %d: expected version %d, actual version %d", rankID, expectedVersion, actualVersion),
			Details: map[string]interface{}{
				"rank_id":          rankID,
				"expected_version": expectedVersion,
				"actual_version":   actualVersion,
			},
			Retryable: true,
			Severity:  ErrorSeverityMedium,
		},
		RankID:          rankID,
		ExpectedVersion: expectedVersion,
		ActualVersion:   actualVersion,
	}
}

// RankingLockError 排行榜锁错误
type RankingLockError struct {
	*BaseRankingError
	RankID    uint32 `json:"rank_id"`
	LockType  string `json:"lock_type"`
	LockOwner string `json:"lock_owner"`
}

// NewRankingLockError 创建排行榜锁错误
func NewRankingLockError(rankID uint32, lockType, lockOwner string) *RankingLockError {
	return &RankingLockError{
		BaseRankingError: &BaseRankingError{
			Code:      "RANKING_LOCK_ERROR",
			Message:   fmt.Sprintf("Cannot acquire %s lock for ranking %d, owned by %s", lockType, rankID, lockOwner),
			Details: map[string]interface{}{
				"rank_id":    rankID,
				"lock_type":  lockType,
				"lock_owner": lockOwner,
			},
			Retryable: true,
			Severity:  ErrorSeverityMedium,
		},
		RankID:    rankID,
		LockType:  lockType,
		LockOwner: lockOwner,
	}
}

// 配置相关错误

// RankingConfigError 排行榜配置错误
type RankingConfigError struct {
	*BaseRankingError
	ConfigKey   string      `json:"config_key"`
	ConfigValue interface{} `json:"config_value"`
	Reason      string      `json:"reason"`
}

// NewRankingConfigError 创建排行榜配置错误
func NewRankingConfigError(configKey string, configValue interface{}, reason string) *RankingConfigError {
	return &RankingConfigError{
		BaseRankingError: &BaseRankingError{
			Code:      "RANKING_CONFIG_ERROR",
			Message:   fmt.Sprintf("Configuration error for %s: %s", configKey, reason),
			Details: map[string]interface{}{
				"config_key":   configKey,
				"config_value": configValue,
				"reason":       reason,
			},
			Retryable: false,
			Severity:  ErrorSeverityHigh,
		},
		ConfigKey:   configKey,
		ConfigValue: configValue,
		Reason:      reason,
	}
}

// 限流相关错误

// RankingRateLimitError 排行榜限流错误
type RankingRateLimitError struct {
	*BaseRankingError
	PlayerID      uint64 `json:"player_id"`
	RankID        uint32 `json:"rank_id"`
	Operation     string `json:"operation"`
	CurrentRate   int64  `json:"current_rate"`
	MaxRate       int64  `json:"max_rate"`
	ResetTime     int64  `json:"reset_time"`
}

// NewRankingRateLimitError 创建排行榜限流错误
func NewRankingRateLimitError(playerID uint64, rankID uint32, operation string, currentRate, maxRate, resetTime int64) *RankingRateLimitError {
	return &RankingRateLimitError{
		BaseRankingError: &BaseRankingError{
			Code:      "RANKING_RATE_LIMIT_ERROR",
			Message:   fmt.Sprintf("Rate limit exceeded for player %d in ranking %d: %d/%d %s operations", playerID, rankID, currentRate, maxRate, operation),
			Details: map[string]interface{}{
				"player_id":    playerID,
				"rank_id":      rankID,
				"operation":    operation,
				"current_rate": currentRate,
				"max_rate":     maxRate,
				"reset_time":   resetTime,
			},
			Retryable: true,
			Severity:  ErrorSeverityMedium,
		},
		PlayerID:    playerID,
		RankID:      rankID,
		Operation:   operation,
		CurrentRate: currentRate,
		MaxRate:     maxRate,
		ResetTime:   resetTime,
	}
}

// 事件相关错误

// RankingEventError 排行榜事件错误
type RankingEventError struct {
	*BaseRankingError
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
	Reason    string `json:"reason"`
}

// NewRankingEventError 创建排行榜事件错误
func NewRankingEventError(eventID, eventType, reason string) *RankingEventError {
	return &RankingEventError{
		BaseRankingError: &BaseRankingError{
			Code:      "RANKING_EVENT_ERROR",
			Message:   fmt.Sprintf("Event error for %s (%s): %s", eventID, eventType, reason),
			Details: map[string]interface{}{
				"event_id":   eventID,
				"event_type": eventType,
				"reason":     reason,
			},
			Retryable: true,
			Severity:  ErrorSeverityMedium,
		},
		EventID:   eventID,
		EventType: eventType,
		Reason:    reason,
	}
}

// 错误代码常量

const (
	// 排行榜相关错误代码
	ErrCodeRankingNotFound      = "RANKING_NOT_FOUND"
	ErrCodeRankingAlreadyExists = "RANKING_ALREADY_EXISTS"
	ErrCodeRankingInactive      = "RANKING_INACTIVE"
	ErrCodeRankingTimeExpired   = "RANKING_TIME_EXPIRED"
	ErrCodeRankingFull          = "RANKING_FULL"
	
	// 玩家相关错误代码
	ErrCodePlayerNotInRanking      = "PLAYER_NOT_IN_RANKING"
	ErrCodePlayerAlreadyInRanking  = "PLAYER_ALREADY_IN_RANKING"
	ErrCodePlayerBlacklisted       = "PLAYER_BLACKLISTED"
	ErrCodePlayerAlreadyBlacklisted = "PLAYER_ALREADY_BLACKLISTED"
	ErrCodePlayerNotBlacklisted    = "PLAYER_NOT_BLACKLISTED"
	
	// 参数相关错误代码
	ErrCodeInvalidRange     = "INVALID_RANGE"
	ErrCodeInvalidTimeRange = "INVALID_TIME_RANGE"
	ErrCodeInvalidScore     = "INVALID_SCORE"
	
	// 权限相关错误代码
	ErrCodeInsufficientPermission = "INSUFFICIENT_PERMISSION"
	ErrCodeOperationNotAllowed    = "OPERATION_NOT_ALLOWED"
	
	// 系统错误代码
	ErrCodeRankingSystemError   = "RANKING_SYSTEM_ERROR"
	ErrCodeRankingDatabaseError = "RANKING_DATABASE_ERROR"
	ErrCodeRankingCacheError    = "RANKING_CACHE_ERROR"
	
	// 验证错误代码
	ErrCodeRankingValidationError = "RANKING_VALIDATION_ERROR"
	
	// 并发相关错误代码
	ErrCodeRankingConcurrencyError = "RANKING_CONCURRENCY_ERROR"
	ErrCodeRankingLockError        = "RANKING_LOCK_ERROR"
	
	// 配置相关错误代码
	ErrCodeRankingConfigError = "RANKING_CONFIG_ERROR"
	
	// 限流相关错误代码
	ErrCodeRankingRateLimitError = "RANKING_RATE_LIMIT_ERROR"
	
	// 事件相关错误代码
	ErrCodeRankingEventError = "RANKING_EVENT_ERROR"
)

// 错误工具函数

// IsRankingError 检查是否为排行榜错误
func IsRankingError(err error) bool {
	_, ok := err.(RankingError)
	return ok
}

// GetRankingErrorCode 获取排行榜错误代码
func GetRankingErrorCode(err error) string {
	if rankingErr, ok := err.(RankingError); ok {
		return rankingErr.GetCode()
	}
	return ""
}

// IsRetryableRankingError 检查是否为可重试的排行榜错误
func IsRetryableRankingError(err error) bool {
	if rankingErr, ok := err.(RankingError); ok {
		return rankingErr.IsRetryable()
	}
	return false
}

// GetRankingErrorSeverity 获取排行榜错误严重程度
func GetRankingErrorSeverity(err error) ErrorSeverity {
	if rankingErr, ok := err.(RankingError); ok {
		return rankingErr.GetSeverity()
	}
	return ErrorSeverityLow
}

// WrapRankingError 包装排行榜错误
func WrapRankingError(err error, code, message string) RankingError {
	return &BaseRankingError{
		Code:      code,
		Message:   fmt.Sprintf("%s: %v", message, err),
		Details:   map[string]interface{}{"wrapped_error": err.Error()},
		Retryable: IsRetryableRankingError(err),
		Severity:  GetRankingErrorSeverity(err),
	}
}

// FormatRankingError 格式化排行榜错误
func FormatRankingError(err RankingError) string {
	return fmt.Sprintf("[%s][%s] %s", err.GetSeverity(), err.GetCode(), err.GetMessage())
}

// LogRankingError 记录排行榜错误（占位符函数）
func LogRankingError(err RankingError) {
	// 实现错误日志记录逻辑
	fmt.Printf("RANKING_ERROR: %s\n", FormatRankingError(err))
}

// 错误分类函数

// IsTemporaryError 检查是否为临时错误
func IsTemporaryError(err error) bool {
	if rankingErr, ok := err.(RankingError); ok {
		code := rankingErr.GetCode()
		return code == ErrCodeRankingSystemError ||
			code == ErrCodeRankingDatabaseError ||
			code == ErrCodeRankingCacheError ||
			code == ErrCodeRankingConcurrencyError ||
			code == ErrCodeRankingLockError ||
			code == ErrCodeRankingRateLimitError
	}
	return false
}

// IsPermanentError 检查是否为永久错误
func IsPermanentError(err error) bool {
	if rankingErr, ok := err.(RankingError); ok {
		code := rankingErr.GetCode()
		return code == ErrCodeRankingNotFound ||
			code == ErrCodeRankingAlreadyExists ||
			code == ErrCodePlayerNotInRanking ||
			code == ErrCodePlayerAlreadyInRanking ||
			code == ErrCodeInvalidRange ||
			code == ErrCodeInvalidTimeRange ||
			code == ErrCodeInvalidScore ||
			code == ErrCodeRankingValidationError
	}
	return false
}

// IsUserError 检查是否为用户错误
func IsUserError(err error) bool {
	if rankingErr, ok := err.(RankingError); ok {
		code := rankingErr.GetCode()
		return code == ErrCodeInvalidRange ||
			code == ErrCodeInvalidTimeRange ||
			code == ErrCodeInvalidScore ||
			code == ErrCodeRankingValidationError ||
			code == ErrCodeInsufficientPermission ||
			code == ErrCodeOperationNotAllowed
	}
	return false
}

// IsSystemError 检查是否为系统错误
func IsSystemError(err error) bool {
	if rankingErr, ok := err.(RankingError); ok {
		code := rankingErr.GetCode()
		return code == ErrCodeRankingSystemError ||
			code == ErrCodeRankingDatabaseError ||
			code == ErrCodeRankingCacheError ||
			code == ErrCodeRankingConfigError
	}
	return false
}

// 错误恢复策略

// ErrorRecoveryStrategy 错误恢复策略
type ErrorRecoveryStrategy struct {
	MaxRetries    int           `json:"max_retries"`
	RetryInterval time.Duration `json:"retry_interval"`
	BackoffFactor float64       `json:"backoff_factor"`
	MaxInterval   time.Duration `json:"max_interval"`
	RecoveryActions []string    `json:"recovery_actions"`
}

// GetRecoveryStrategy 获取错误恢复策略
func GetRecoveryStrategy(err error) *ErrorRecoveryStrategy {
	if !IsRankingError(err) {
		return nil
	}
	
	rankingErr := err.(RankingError)
	code := rankingErr.GetCode()
	
	switch code {
	case ErrCodeRankingSystemError, ErrCodeRankingDatabaseError:
		return &ErrorRecoveryStrategy{
			MaxRetries:      5,
			RetryInterval:   time.Second,
			BackoffFactor:   2.0,
			MaxInterval:     30 * time.Second,
			RecoveryActions: []string{"retry", "fallback", "alert"},
		}
	case ErrCodeRankingCacheError:
		return &ErrorRecoveryStrategy{
			MaxRetries:      3,
			RetryInterval:   500 * time.Millisecond,
			BackoffFactor:   1.5,
			MaxInterval:     5 * time.Second,
			RecoveryActions: []string{"retry", "bypass_cache"},
		}
	case ErrCodeRankingConcurrencyError, ErrCodeRankingLockError:
		return &ErrorRecoveryStrategy{
			MaxRetries:      10,
			RetryInterval:   100 * time.Millisecond,
			BackoffFactor:   1.2,
			MaxInterval:     2 * time.Second,
			RecoveryActions: []string{"retry", "backoff"},
		}
	case ErrCodeRankingRateLimitError:
		return &ErrorRecoveryStrategy{
			MaxRetries:      1,
			RetryInterval:   time.Minute,
			BackoffFactor:   1.0,
			MaxInterval:     time.Minute,
			RecoveryActions: []string{"wait", "throttle"},
		}
	default:
		return &ErrorRecoveryStrategy{
			MaxRetries:      0,
			RetryInterval:   0,
			BackoffFactor:   1.0,
			MaxInterval:     0,
			RecoveryActions: []string{"fail"},
		}
	}
}