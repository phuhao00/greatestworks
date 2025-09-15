package honor

import (
	"errors"
	"fmt"
)

// 荣誉系统相关错误定义

// 称号相关错误
var (
	// ErrInvalidTitle 无效的称号
	ErrInvalidTitle = errors.New("invalid title")
	
	// ErrTitleNotFound 称号未找到
	ErrTitleNotFound = errors.New("title not found")
	
	// ErrTitleAlreadyOwned 称号已拥有
	ErrTitleAlreadyOwned = errors.New("title already owned")
	
	// ErrTitleAlreadyUnlocked 称号已解锁
	ErrTitleAlreadyUnlocked = errors.New("title already unlocked")
	
	// ErrTitleNotUnlocked 称号未解锁
	ErrTitleNotUnlocked = errors.New("title not unlocked")
	
	// ErrTitleConditionNotMet 称号解锁条件未满足
	ErrTitleConditionNotMet = errors.New("title unlock condition not met")
	
	// ErrNoTitleEquipped 没有装备称号
	ErrNoTitleEquipped = errors.New("no title equipped")
	
	// ErrTitleAlreadyEquipped 称号已装备
	ErrTitleAlreadyEquipped = errors.New("title already equipped")
	
	// ErrCannotEquipTitle 无法装备称号
	ErrCannotEquipTitle = errors.New("cannot equip title")
	
	// ErrTitleTemplateNotFound 称号模板未找到
	ErrTitleTemplateNotFound = errors.New("title template not found")
	
	// ErrTitleTemplateAlreadyExists 称号模板已存在
	ErrTitleTemplateAlreadyExists = errors.New("title template already exists")
)

// 成就相关错误
var (
	// ErrInvalidAchievement 无效的成就
	ErrInvalidAchievement = errors.New("invalid achievement")
	
	// ErrAchievementNotFound 成就未找到
	ErrAchievementNotFound = errors.New("achievement not found")
	
	// ErrAchievementAlreadyOwned 成就已拥有
	ErrAchievementAlreadyOwned = errors.New("achievement already owned")
	
	// ErrAchievementAlreadyUnlocked 成就已解锁
	ErrAchievementAlreadyUnlocked = errors.New("achievement already unlocked")
	
	// ErrAchievementConditionNotMet 成就解锁条件未满足
	ErrAchievementConditionNotMet = errors.New("achievement unlock condition not met")
	
	// ErrAchievementTemplateNotFound 成就模板未找到
	ErrAchievementTemplateNotFound = errors.New("achievement template not found")
	
	// ErrAchievementTemplateAlreadyExists 成就模板已存在
	ErrAchievementTemplateAlreadyExists = errors.New("achievement template already exists")
)

// 荣誉系统相关错误
var (
	// ErrInvalidPlayerID 无效的玩家ID
	ErrInvalidPlayerID = errors.New("invalid player id")
	
	// ErrHonorNotFound 荣誉信息未找到
	ErrHonorNotFound = errors.New("honor not found")
	
	// ErrHonorAlreadyExists 荣誉信息已存在
	ErrHonorAlreadyExists = errors.New("honor already exists")
	
	// ErrInvalidHonorPoints 无效的荣誉点数
	ErrInvalidHonorPoints = errors.New("invalid honor points")
	
	// ErrInvalidHonorLevel 无效的荣誉等级
	ErrInvalidHonorLevel = errors.New("invalid honor level")
	
	// ErrHonorLevelNotReached 荣誉等级未达到
	ErrHonorLevelNotReached = errors.New("honor level not reached")
	
	// ErrInsufficientHonorPoints 荣誉点数不足
	ErrInsufficientHonorPoints = errors.New("insufficient honor points")
)

// 声望相关错误
var (
	// ErrInvalidFaction 无效的阵营
	ErrInvalidFaction = errors.New("invalid faction")
	
	// ErrFactionNotFound 阵营未找到
	ErrFactionNotFound = errors.New("faction not found")
	
	// ErrInvalidReputation 无效的声望值
	ErrInvalidReputation = errors.New("invalid reputation")
	
	// ErrReputationTooLow 声望过低
	ErrReputationTooLow = errors.New("reputation too low")
	
	// ErrReputationTooHigh 声望过高
	ErrReputationTooHigh = errors.New("reputation too high")
	
	// ErrCannotChangeReputation 无法改变声望
	ErrCannotChangeReputation = errors.New("cannot change reputation")
)

// 统计数据相关错误
var (
	// ErrInvalidStatisticType 无效的统计数据类型
	ErrInvalidStatisticType = errors.New("invalid statistic type")
	
	// ErrInvalidStatisticValue 无效的统计数据值
	ErrInvalidStatisticValue = errors.New("invalid statistic value")
	
	// ErrStatisticNotFound 统计数据未找到
	ErrStatisticNotFound = errors.New("statistic not found")
	
	// ErrCannotUpdateStatistic 无法更新统计数据
	ErrCannotUpdateStatistic = errors.New("cannot update statistic")
)

// 解锁条件相关错误
var (
	// ErrInvalidUnlockCondition 无效的解锁条件
	ErrInvalidUnlockCondition = errors.New("invalid unlock condition")
	
	// ErrUnlockConditionNotMet 解锁条件未满足
	ErrUnlockConditionNotMet = errors.New("unlock condition not met")
	
	// ErrInvalidConditionType 无效的条件类型
	ErrInvalidConditionType = errors.New("invalid condition type")
	
	// ErrInvalidConditionValue 无效的条件值
	ErrInvalidConditionValue = errors.New("invalid condition value")
	
	// ErrCircularDependency 循环依赖
	ErrCircularDependency = errors.New("circular dependency in unlock conditions")
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
)

// 配置相关错误
var (
	// ErrInvalidConfiguration 无效配置
	ErrInvalidConfiguration = errors.New("invalid configuration")
	
	// ErrConfigurationNotFound 配置未找到
	ErrConfigurationNotFound = errors.New("configuration not found")
	
	// ErrConfigurationLoadFailure 配置加载失败
	ErrConfigurationLoadFailure = errors.New("configuration load failure")
)

// 事件相关错误
var (
	// ErrInvalidEvent 无效事件
	ErrInvalidEvent = errors.New("invalid event")
	
	// ErrEventPublishFailure 事件发布失败
	ErrEventPublishFailure = errors.New("event publish failure")
	
	// ErrEventHandlerNotFound 事件处理器未找到
	ErrEventHandlerNotFound = errors.New("event handler not found")
	
	// ErrEventProcessingFailure 事件处理失败
	ErrEventProcessingFailure = errors.New("event processing failure")
)

// HonorError 荣誉系统错误类型
type HonorError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Cause   error  `json:"-"`
}

// Error 实现error接口
func (e *HonorError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 返回原始错误
func (e *HonorError) Unwrap() error {
	return e.Cause
}

// NewHonorError 创建荣誉系统错误
func NewHonorError(code, message string) *HonorError {
	return &HonorError{
		Code:    code,
		Message: message,
	}
}

// NewHonorErrorWithDetails 创建带详情的荣誉系统错误
func NewHonorErrorWithDetails(code, message, details string) *HonorError {
	return &HonorError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewHonorErrorWithCause 创建带原因的荣誉系统错误
func NewHonorErrorWithCause(code, message string, cause error) *HonorError {
	return &HonorError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// 预定义的错误代码常量
const (
	// 称号相关错误代码
	ErrCodeTitleNotFound         = "TITLE_NOT_FOUND"
	ErrCodeTitleAlreadyUnlocked  = "TITLE_ALREADY_UNLOCKED"
	ErrCodeTitleConditionNotMet  = "TITLE_CONDITION_NOT_MET"
	ErrCodeTitleCannotEquip      = "TITLE_CANNOT_EQUIP"
	
	// 成就相关错误代码
	ErrCodeAchievementNotFound        = "ACHIEVEMENT_NOT_FOUND"
	ErrCodeAchievementAlreadyUnlocked = "ACHIEVEMENT_ALREADY_UNLOCKED"
	ErrCodeAchievementConditionNotMet = "ACHIEVEMENT_CONDITION_NOT_MET"
	
	// 荣誉系统错误代码
	ErrCodeHonorNotFound           = "HONOR_NOT_FOUND"
	ErrCodeInsufficientHonorPoints = "INSUFFICIENT_HONOR_POINTS"
	ErrCodeInvalidHonorLevel       = "INVALID_HONOR_LEVEL"
	
	// 声望相关错误代码
	ErrCodeFactionNotFound      = "FACTION_NOT_FOUND"
	ErrCodeInvalidReputation    = "INVALID_REPUTATION"
	ErrCodeReputationTooLow     = "REPUTATION_TOO_LOW"
	
	// 数据相关错误代码
	ErrCodeDataNotFound       = "DATA_NOT_FOUND"
	ErrCodeDataCorrupted      = "DATA_CORRUPTED"
	ErrCodeVersionConflict    = "VERSION_CONFLICT"
	ErrCodeSaveFailure        = "SAVE_FAILURE"
	
	// 业务逻辑错误代码
	ErrCodeOperationNotAllowed = "OPERATION_NOT_ALLOWED"
	ErrCodePermissionDenied    = "PERMISSION_DENIED"
	ErrCodeRateLimitExceeded   = "RATE_LIMIT_EXCEEDED"
	ErrCodeMaintenanceMode     = "MAINTENANCE_MODE"
)

// IsHonorError 检查是否为荣誉系统错误
func IsHonorError(err error) bool {
	_, ok := err.(*HonorError)
	return ok
}

// GetHonorErrorCode 获取荣誉系统错误代码
func GetHonorErrorCode(err error) string {
	if honorErr, ok := err.(*HonorError); ok {
		return honorErr.Code
	}
	return ""
}

// WrapError 包装错误为荣誉系统错误
func WrapError(err error, code, message string) *HonorError {
	return &HonorError{
		Code:    code,
		Message: message,
		Cause:   err,
	}
}