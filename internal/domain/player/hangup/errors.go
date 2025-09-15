package hangup

import (
	"errors"
	"fmt"
)

// 挂机系统相关错误定义

// 挂机地点相关错误
var (
	// ErrInvalidHangupLocation 无效的挂机地点
	ErrInvalidHangupLocation = errors.New("invalid hangup location")
	
	// ErrHangupLocationNotFound 挂机地点未找到
	ErrHangupLocationNotFound = errors.New("hangup location not found")
	
	// ErrHangupLocationNotUnlocked 挂机地点未解锁
	ErrHangupLocationNotUnlocked = errors.New("hangup location not unlocked")
	
	// ErrHangupLocationRequirementNotMet 挂机地点要求未满足
	ErrHangupLocationRequirementNotMet = errors.New("hangup location requirement not met")
	
	// ErrHangupLocationInactive 挂机地点未激活
	ErrHangupLocationInactive = errors.New("hangup location inactive")
	
	// ErrHangupLocationFull 挂机地点已满
	ErrHangupLocationFull = errors.New("hangup location full")
	
	// ErrNoHangupLocationSet 未设置挂机地点
	ErrNoHangupLocationSet = errors.New("no hangup location set")
	
	// ErrCannotChangeLocationWhileHanging 挂机中无法更换地点
	ErrCannotChangeLocationWhileHanging = errors.New("cannot change location while hanging up")
)

// 挂机状态相关错误
var (
	// ErrAlreadyHangingUp 已经在挂机
	ErrAlreadyHangingUp = errors.New("already hanging up")
	
	// ErrNotHangingUp 未在挂机
	ErrNotHangingUp = errors.New("not hanging up")
	
	// ErrHangupPaused 挂机已暂停
	ErrHangupPaused = errors.New("hangup paused")
	
	// ErrCannotStartHangup 无法开始挂机
	ErrCannotStartHangup = errors.New("cannot start hangup")
	
	// ErrCannotStopHangup 无法停止挂机
	ErrCannotStopHangup = errors.New("cannot stop hangup")
	
	// ErrHangupCooldown 挂机冷却中
	ErrHangupCooldown = errors.New("hangup cooldown")
	
	// ErrHangupLimitExceeded 挂机限制超出
	ErrHangupLimitExceeded = errors.New("hangup limit exceeded")
)

// 离线奖励相关错误
var (
	// ErrNoOfflineRewardAvailable 没有可用的离线奖励
	ErrNoOfflineRewardAvailable = errors.New("no offline reward available")
	
	// ErrOfflineRewardAlreadyClaimed 离线奖励已领取
	ErrOfflineRewardAlreadyClaimed = errors.New("offline reward already claimed")
	
	// ErrOfflineRewardExpired 离线奖励已过期
	ErrOfflineRewardExpired = errors.New("offline reward expired")
	
	// ErrInvalidOfflineReward 无效的离线奖励
	ErrInvalidOfflineReward = errors.New("invalid offline reward")
	
	// ErrOfflineRewardCalculationFailed 离线奖励计算失败
	ErrOfflineRewardCalculationFailed = errors.New("offline reward calculation failed")
	
	// ErrOfflineTimeTooShort 离线时间太短
	ErrOfflineTimeTooShort = errors.New("offline time too short")
	
	// ErrOfflineTimeTooLong 离线时间太长
	ErrOfflineTimeTooLong = errors.New("offline time too long")
)

// 效率加成相关错误
var (
	// ErrInvalidEfficiencyBonus 无效的效率加成
	ErrInvalidEfficiencyBonus = errors.New("invalid efficiency bonus")
	
	// ErrEfficiencyBonusNotFound 效率加成未找到
	ErrEfficiencyBonusNotFound = errors.New("efficiency bonus not found")
	
	// ErrCannotUpdateEfficiencyBonus 无法更新效率加成
	ErrCannotUpdateEfficiencyBonus = errors.New("cannot update efficiency bonus")
	
	// ErrEfficiencyBonusExpired 效率加成已过期
	ErrEfficiencyBonusExpired = errors.New("efficiency bonus expired")
	
	// ErrMaxEfficiencyBonusReached 已达到最大效率加成
	ErrMaxEfficiencyBonusReached = errors.New("max efficiency bonus reached")
)

// 玩家相关错误
var (
	// ErrInvalidPlayerID 无效的玩家ID
	ErrInvalidPlayerID = errors.New("invalid player id")
	
	// ErrPlayerNotFound 玩家未找到
	ErrPlayerNotFound = errors.New("player not found")
	
	// ErrPlayerLevelTooLow 玩家等级太低
	ErrPlayerLevelTooLow = errors.New("player level too low")
	
	// ErrPlayerNotOnline 玩家不在线
	ErrPlayerNotOnline = errors.New("player not online")
	
	// ErrPlayerBanned 玩家被封禁
	ErrPlayerBanned = errors.New("player banned")
	
	// ErrPlayerInCombat 玩家在战斗中
	ErrPlayerInCombat = errors.New("player in combat")
)

// 配置相关错误
var (
	// ErrInvalidHangupConfig 无效的挂机配置
	ErrInvalidHangupConfig = errors.New("invalid hangup config")
	
	// ErrHangupConfigNotFound 挂机配置未找到
	ErrHangupConfigNotFound = errors.New("hangup config not found")
	
	// ErrCannotUpdateHangupConfig 无法更新挂机配置
	ErrCannotUpdateHangupConfig = errors.New("cannot update hangup config")
	
	// ErrHangupConfigValidationFailed 挂机配置验证失败
	ErrHangupConfigValidationFailed = errors.New("hangup config validation failed")
)

// 会话相关错误
var (
	// ErrInvalidHangupSession 无效的挂机会话
	ErrInvalidHangupSession = errors.New("invalid hangup session")
	
	// ErrHangupSessionNotFound 挂机会话未找到
	ErrHangupSessionNotFound = errors.New("hangup session not found")
	
	// ErrHangupSessionAlreadyEnded 挂机会话已结束
	ErrHangupSessionAlreadyEnded = errors.New("hangup session already ended")
	
	// ErrCannotEndHangupSession 无法结束挂机会话
	ErrCannotEndHangupSession = errors.New("cannot end hangup session")
	
	// ErrHangupSessionExpired 挂机会话已过期
	ErrHangupSessionExpired = errors.New("hangup session expired")
)

// 统计相关错误
var (
	// ErrInvalidHangupStatistics 无效的挂机统计
	ErrInvalidHangupStatistics = errors.New("invalid hangup statistics")
	
	// ErrHangupStatisticsNotFound 挂机统计未找到
	ErrHangupStatisticsNotFound = errors.New("hangup statistics not found")
	
	// ErrCannotUpdateHangupStatistics 无法更新挂机统计
	ErrCannotUpdateHangupStatistics = errors.New("cannot update hangup statistics")
	
	// ErrStatisticsCalculationFailed 统计计算失败
	ErrStatisticsCalculationFailed = errors.New("statistics calculation failed")
)

// 数据持久化相关错误
var (
	// ErrDatabaseConnection 数据库连接错误
	ErrDatabaseConnection = errors.New("database connection error")
	
	// ErrDataNotFound 数据未找到
	ErrDataNotFound = errors.New("data not found")
	
	// ErrDataCorrupted 数据损坏
	ErrDataCorrupted = errors.New("data corrupted")
	
	// ErrSaveFailure 保存失败
	ErrSaveFailure = errors.New("save failure")
	
	// ErrLoadFailure 加载失败
	ErrLoadFailure = errors.New("load failure")
	
	// ErrDeleteFailure 删除失败
	ErrDeleteFailure = errors.New("delete failure")
	
	// ErrVersionConflict 版本冲突
	ErrVersionConflict = errors.New("version conflict")
	
	// ErrConcurrentModification 并发修改冲突
	ErrConcurrentModification = errors.New("concurrent modification")
	
	// ErrTransactionFailed 事务失败
	ErrTransactionFailed = errors.New("transaction failed")
)

// 业务逻辑相关错误
var (
	// ErrOperationNotAllowed 操作不被允许
	ErrOperationNotAllowed = errors.New("operation not allowed")
	
	// ErrInvalidOperation 无效操作
	ErrInvalidOperation = errors.New("invalid operation")
	
	// ErrPermissionDenied 权限被拒绝
	ErrPermissionDenied = errors.New("permission denied")
	
	// ErrResourceLocked 资源被锁定
	ErrResourceLocked = errors.New("resource locked")
	
	// ErrRateLimitExceeded 速率限制超出
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	
	// ErrMaintenanceMode 维护模式
	ErrMaintenanceMode = errors.New("system in maintenance mode")
	
	// ErrServiceUnavailable 服务不可用
	ErrServiceUnavailable = errors.New("service unavailable")
)

// 验证相关错误
var (
	// ErrInvalidInput 无效输入
	ErrInvalidInput = errors.New("invalid input")
	
	// ErrValidationFailed 验证失败
	ErrValidationFailed = errors.New("validation failed")
	
	// ErrInvalidTimeRange 无效时间范围
	ErrInvalidTimeRange = errors.New("invalid time range")
	
	// ErrInvalidDuration 无效持续时间
	ErrInvalidDuration = errors.New("invalid duration")
	
	// ErrInvalidRewardAmount 无效奖励数量
	ErrInvalidRewardAmount = errors.New("invalid reward amount")
)

// HangupError 挂机系统错误类型
type HangupError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Cause   error  `json:"-"`
}

// Error 实现error接口
func (e *HangupError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 返回原始错误
func (e *HangupError) Unwrap() error {
	return e.Cause
}

// NewHangupError 创建挂机系统错误
func NewHangupError(code, message string) *HangupError {
	return &HangupError{
		Code:    code,
		Message: message,
	}
}

// NewHangupErrorWithDetails 创建带详情的挂机系统错误
func NewHangupErrorWithDetails(code, message, details string) *HangupError {
	return &HangupError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewHangupErrorWithCause 创建带原因的挂机系统错误
func NewHangupErrorWithCause(code, message string, cause error) *HangupError {
	return &HangupError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// 预定义的错误代码常量
const (
	// 地点相关错误代码
	ErrCodeLocationNotFound         = "LOCATION_NOT_FOUND"
	ErrCodeLocationNotUnlocked      = "LOCATION_NOT_UNLOCKED"
	ErrCodeLocationRequirementNotMet = "LOCATION_REQUIREMENT_NOT_MET"
	ErrCodeLocationInactive         = "LOCATION_INACTIVE"
	ErrCodeLocationFull             = "LOCATION_FULL"
	ErrCodeNoLocationSet            = "NO_LOCATION_SET"
	
	// 挂机状态错误代码
	ErrCodeAlreadyHangingUp      = "ALREADY_HANGING_UP"
	ErrCodeNotHangingUp          = "NOT_HANGING_UP"
	ErrCodeHangupPaused          = "HANGUP_PAUSED"
	ErrCodeCannotStartHangup     = "CANNOT_START_HANGUP"
	ErrCodeCannotStopHangup      = "CANNOT_STOP_HANGUP"
	ErrCodeHangupCooldown        = "HANGUP_COOLDOWN"
	ErrCodeHangupLimitExceeded   = "HANGUP_LIMIT_EXCEEDED"
	
	// 奖励相关错误代码
	ErrCodeNoOfflineReward           = "NO_OFFLINE_REWARD"
	ErrCodeOfflineRewardClaimed      = "OFFLINE_REWARD_CLAIMED"
	ErrCodeOfflineRewardExpired      = "OFFLINE_REWARD_EXPIRED"
	ErrCodeRewardCalculationFailed   = "REWARD_CALCULATION_FAILED"
	ErrCodeOfflineTimeTooShort       = "OFFLINE_TIME_TOO_SHORT"
	ErrCodeOfflineTimeTooLong        = "OFFLINE_TIME_TOO_LONG"
	
	// 效率加成错误代码
	ErrCodeInvalidEfficiencyBonus    = "INVALID_EFFICIENCY_BONUS"
	ErrCodeEfficiencyBonusNotFound   = "EFFICIENCY_BONUS_NOT_FOUND"
	ErrCodeEfficiencyBonusExpired    = "EFFICIENCY_BONUS_EXPIRED"
	ErrCodeMaxEfficiencyBonusReached = "MAX_EFFICIENCY_BONUS_REACHED"
	
	// 玩家相关错误代码
	ErrCodePlayerNotFound    = "PLAYER_NOT_FOUND"
	ErrCodePlayerLevelTooLow = "PLAYER_LEVEL_TOO_LOW"
	ErrCodePlayerNotOnline   = "PLAYER_NOT_ONLINE"
	ErrCodePlayerBanned      = "PLAYER_BANNED"
	ErrCodePlayerInCombat    = "PLAYER_IN_COMBAT"
	
	// 配置相关错误代码
	ErrCodeInvalidConfig           = "INVALID_CONFIG"
	ErrCodeConfigNotFound          = "CONFIG_NOT_FOUND"
	ErrCodeConfigValidationFailed  = "CONFIG_VALIDATION_FAILED"
	
	// 会话相关错误代码
	ErrCodeSessionNotFound      = "SESSION_NOT_FOUND"
	ErrCodeSessionAlreadyEnded  = "SESSION_ALREADY_ENDED"
	ErrCodeCannotEndSession     = "CANNOT_END_SESSION"
	ErrCodeSessionExpired       = "SESSION_EXPIRED"
	
	// 统计相关错误代码
	ErrCodeStatisticsNotFound        = "STATISTICS_NOT_FOUND"
	ErrCodeStatisticsCalculationFailed = "STATISTICS_CALCULATION_FAILED"
	
	// 数据相关错误代码
	ErrCodeDataNotFound          = "DATA_NOT_FOUND"
	ErrCodeDataCorrupted         = "DATA_CORRUPTED"
	ErrCodeVersionConflict       = "VERSION_CONFLICT"
	ErrCodeSaveFailure           = "SAVE_FAILURE"
	ErrCodeLoadFailure           = "LOAD_FAILURE"
	ErrCodeDeleteFailure         = "DELETE_FAILURE"
	ErrCodeConcurrentModification = "CONCURRENT_MODIFICATION"
	ErrCodeTransactionFailed     = "TRANSACTION_FAILED"
	
	// 业务逻辑错误代码
	ErrCodeOperationNotAllowed = "OPERATION_NOT_ALLOWED"
	ErrCodePermissionDenied    = "PERMISSION_DENIED"
	ErrCodeRateLimitExceeded   = "RATE_LIMIT_EXCEEDED"
	ErrCodeMaintenanceMode     = "MAINTENANCE_MODE"
	ErrCodeServiceUnavailable  = "SERVICE_UNAVAILABLE"
	
	// 验证相关错误代码
	ErrCodeInvalidInput       = "INVALID_INPUT"
	ErrCodeValidationFailed   = "VALIDATION_FAILED"
	ErrCodeInvalidTimeRange   = "INVALID_TIME_RANGE"
	ErrCodeInvalidDuration    = "INVALID_DURATION"
	ErrCodeInvalidRewardAmount = "INVALID_REWARD_AMOUNT"
)

// IsHangupError 检查是否为挂机系统错误
func IsHangupError(err error) bool {
	_, ok := err.(*HangupError)
	return ok
}

// GetHangupErrorCode 获取挂机系统错误代码
func GetHangupErrorCode(err error) string {
	if hangupErr, ok := err.(*HangupError); ok {
		return hangupErr.Code
	}
	return ""
}

// WrapError 包装错误为挂机系统错误
func WrapError(err error, code, message string) *HangupError {
	return &HangupError{
		Code:    code,
		Message: message,
		Cause:   err,
	}
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Value   interface{} `json:"value"`
	Message string `json:"message"`
}

// Error 实现error接口
func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", ve.Field, ve.Message)
}

// NewValidationError 创建验证错误
func NewValidationError(field string, value interface{}, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

// BusinessRuleError 业务规则错误
type BusinessRuleError struct {
	Rule        string `json:"rule"`
	Description string `json:"description"`
	Suggestion  string `json:"suggestion,omitempty"`
}

// Error 实现error接口
func (bre *BusinessRuleError) Error() string {
	if bre.Suggestion != "" {
		return fmt.Sprintf("business rule violation '%s': %s. Suggestion: %s", bre.Rule, bre.Description, bre.Suggestion)
	}
	return fmt.Sprintf("business rule violation '%s': %s", bre.Rule, bre.Description)
}

// NewBusinessRuleError 创建业务规则错误
func NewBusinessRuleError(rule, description, suggestion string) *BusinessRuleError {
	return &BusinessRuleError{
		Rule:        rule,
		Description: description,
		Suggestion:  suggestion,
	}
}

// ConcurrencyError 并发错误
type ConcurrencyError struct {
	Resource    string `json:"resource"`
	Operation   string `json:"operation"`
	ConflictType string `json:"conflict_type"`
	Message     string `json:"message"`
}

// Error 实现error接口
func (ce *ConcurrencyError) Error() string {
	return fmt.Sprintf("concurrency error on %s during %s (%s): %s", ce.Resource, ce.Operation, ce.ConflictType, ce.Message)
}

// NewConcurrencyError 创建并发错误
func NewConcurrencyError(resource, operation, conflictType, message string) *ConcurrencyError {
	return &ConcurrencyError{
		Resource:     resource,
		Operation:    operation,
		ConflictType: conflictType,
		Message:      message,
	}
}

// ErrorCollection 错误集合
type ErrorCollection struct {
	Errors []error `json:"errors"`
}

// Error 实现error接口
func (ec *ErrorCollection) Error() string {
	if len(ec.Errors) == 0 {
		return "no errors"
	}
	if len(ec.Errors) == 1 {
		return ec.Errors[0].Error()
	}
	return fmt.Sprintf("multiple errors occurred: %d errors", len(ec.Errors))
}

// Add 添加错误
func (ec *ErrorCollection) Add(err error) {
	if err != nil {
		ec.Errors = append(ec.Errors, err)
	}
}

// HasErrors 是否有错误
func (ec *ErrorCollection) HasErrors() bool {
	return len(ec.Errors) > 0
}

// Count 错误数量
func (ec *ErrorCollection) Count() int {
	return len(ec.Errors)
}

// First 获取第一个错误
func (ec *ErrorCollection) First() error {
	if len(ec.Errors) > 0 {
		return ec.Errors[0]
	}
	return nil
}

// NewErrorCollection 创建错误集合
func NewErrorCollection() *ErrorCollection {
	return &ErrorCollection{
		Errors: make([]error, 0),
	}
}