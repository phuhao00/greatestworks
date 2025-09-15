package sacred

import (
	"fmt"
	"strings"
	"time"
)

// 基础错误变量
var (
	// 圣地相关错误
	ErrSacredNotFound        = fmt.Errorf("sacred place not found")
	ErrSacredAlreadyExists   = fmt.Errorf("sacred place already exists")
	ErrSacredNotActive       = fmt.Errorf("sacred place is not active")
	ErrSacredLocked          = fmt.Errorf("sacred place is locked")
	ErrSacredMaintenance     = fmt.Errorf("sacred place is under maintenance")
	ErrSacredAccessDenied    = fmt.Errorf("access to sacred place denied")
	ErrSacredCapacityFull    = fmt.Errorf("sacred place capacity is full")
	ErrInvalidSacredName     = fmt.Errorf("invalid sacred place name")
	ErrInvalidSacredStatus   = fmt.Errorf("invalid sacred place status")
	ErrSacredOwnerMismatch   = fmt.Errorf("sacred place owner mismatch")

	// 挑战相关错误
	ErrChallengeNotFound      = fmt.Errorf("challenge not found")
	ErrChallengeAlreadyExists = fmt.Errorf("challenge already exists")
	ErrChallengeNotAvailable  = fmt.Errorf("challenge is not available")
	ErrChallengeInProgress    = fmt.Errorf("challenge is in progress")
	ErrChallengeCompleted     = fmt.Errorf("challenge already completed")
	ErrChallengeFailed        = fmt.Errorf("challenge failed")
	ErrChallengeExpired       = fmt.Errorf("challenge has expired")
	ErrChallengeOnCooldown    = fmt.Errorf("challenge is on cooldown")
	ErrInvalidChallenge       = fmt.Errorf("invalid challenge")
	ErrInvalidChallengeType   = fmt.Errorf("invalid challenge type")
	ErrInvalidDifficulty      = fmt.Errorf("invalid challenge difficulty")
	ErrInsufficientLevel      = fmt.Errorf("insufficient level for challenge")
	ErrChallengeConditionsNotMet = fmt.Errorf("challenge conditions not met")

	// 祝福相关错误
	ErrBlessingNotFound       = fmt.Errorf("blessing not found")
	ErrBlessingAlreadyExists  = fmt.Errorf("blessing already exists")
	ErrBlessingNotAvailable   = fmt.Errorf("blessing is not available")
	ErrBlessingExpired        = fmt.Errorf("blessing has expired")
	ErrBlessingOnCooldown     = fmt.Errorf("blessing is on cooldown")
	ErrBlessingLimitReached   = fmt.Errorf("blessing usage limit reached")
	ErrInvalidBlessing        = fmt.Errorf("invalid blessing")
	ErrInvalidBlessingType    = fmt.Errorf("invalid blessing type")
	ErrBlessingConflict       = fmt.Errorf("blessing conflicts with existing effects")
	ErrMaxActiveBlessings     = fmt.Errorf("maximum active blessings reached")

	// 圣物相关错误
	ErrRelicNotFound          = fmt.Errorf("relic not found")
	ErrRelicAlreadyExists     = fmt.Errorf("relic already exists")
	ErrRelicNotOwned          = fmt.Errorf("relic is not owned by player")
	ErrRelicCannotUpgrade     = fmt.Errorf("relic cannot be upgraded")
	ErrRelicMaxLevel          = fmt.Errorf("relic is at maximum level")
	ErrRelicRequirementsNotMet = fmt.Errorf("relic requirements not met")
	ErrInvalidRelic           = fmt.Errorf("invalid relic")
	ErrInvalidRelicType       = fmt.Errorf("invalid relic type")
	ErrInvalidRelicRarity     = fmt.Errorf("invalid relic rarity")
	ErrRelicInventoryFull     = fmt.Errorf("relic inventory is full")

	// 等级和经验相关错误
	ErrInvalidLevel           = fmt.Errorf("invalid level")
	ErrInvalidExperience      = fmt.Errorf("invalid experience")
	ErrMaxLevelReached        = fmt.Errorf("maximum level reached")
	ErrInsufficientExperience = fmt.Errorf("insufficient experience")
	ErrLevelDowngrade         = fmt.Errorf("level downgrade not allowed")

	// 权限相关错误
	ErrPermissionDenied       = fmt.Errorf("permission denied")
	ErrUnauthorized           = fmt.Errorf("unauthorized operation")
	ErrAccessRestricted       = fmt.Errorf("access restricted")
	ErrInsufficientPrivileges = fmt.Errorf("insufficient privileges")
	ErrOwnershipRequired      = fmt.Errorf("ownership required")

	// 资源相关错误
	ErrInsufficientResources  = fmt.Errorf("insufficient resources")
	ErrInsufficientGold       = fmt.Errorf("insufficient gold")
	ErrInsufficientMana       = fmt.Errorf("insufficient mana")
	ErrInsufficientEnergy     = fmt.Errorf("insufficient energy")
	ErrResourceNotFound       = fmt.Errorf("resource not found")
	ErrResourceLocked         = fmt.Errorf("resource is locked")

	// 时间相关错误
	ErrInvalidTime            = fmt.Errorf("invalid time")
	ErrTimeExpired            = fmt.Errorf("time has expired")
	ErrTooEarly               = fmt.Errorf("too early for this operation")
	ErrTooLate                = fmt.Errorf("too late for this operation")
	ErrCooldownActive         = fmt.Errorf("cooldown is active")
	ErrDurationTooShort       = fmt.Errorf("duration is too short")
	ErrDurationTooLong        = fmt.Errorf("duration is too long")

	// 配置相关错误
	ErrInvalidConfiguration   = fmt.Errorf("invalid configuration")
	ErrConfigurationNotFound  = fmt.Errorf("configuration not found")
	ErrConfigurationCorrupted = fmt.Errorf("configuration corrupted")
	ErrMissingConfiguration   = fmt.Errorf("missing configuration")

	// 数据相关错误
	ErrDataCorrupted          = fmt.Errorf("data corrupted")
	ErrDataNotFound           = fmt.Errorf("data not found")
	ErrDataInconsistent       = fmt.Errorf("data inconsistent")
	ErrInvalidData            = fmt.Errorf("invalid data")
	ErrDataConflict           = fmt.Errorf("data conflict")

	// 系统相关错误
	ErrSystemError            = fmt.Errorf("system error")
	ErrServiceUnavailable     = fmt.Errorf("service unavailable")
	ErrTimeout                = fmt.Errorf("operation timeout")
	ErrInternalError          = fmt.Errorf("internal error")
	ErrExternalServiceError   = fmt.Errorf("external service error")

	// 并发相关错误
	ErrConcurrentModification = fmt.Errorf("concurrent modification")
	ErrResourceBusy           = fmt.Errorf("resource is busy")
	ErrDeadlock               = fmt.Errorf("deadlock detected")
	ErrRaceCondition          = fmt.Errorf("race condition detected")

	// 验证相关错误
	ErrValidationFailed       = fmt.Errorf("validation failed")
	ErrInvalidInput           = fmt.Errorf("invalid input")
	ErrMissingRequiredField   = fmt.Errorf("missing required field")
	ErrFieldTooLong           = fmt.Errorf("field is too long")
	ErrFieldTooShort          = fmt.Errorf("field is too short")
	ErrInvalidFormat          = fmt.Errorf("invalid format")

	// 业务规则相关错误
	ErrBusinessRuleViolation  = fmt.Errorf("business rule violation")
	ErrOperationNotAllowed    = fmt.Errorf("operation not allowed")
	ErrStateTransitionInvalid = fmt.Errorf("invalid state transition")
	ErrPreconditionFailed     = fmt.Errorf("precondition failed")
	ErrPostconditionFailed    = fmt.Errorf("postcondition failed")
)

// SacredError 圣地系统错误
type SacredError struct {
	Code      string
	Message   string
	Details   map[string]interface{}
	Cause     error
	Timestamp time.Time
	Context   map[string]string
	Severity  ErrorSeverity
	Category  ErrorCategory
}

// Error 实现error接口
func (e *SacredError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap 返回原始错误
func (e *SacredError) Unwrap() error {
	return e.Cause
}

// WithDetail 添加详细信息
func (e *SacredError) WithDetail(key string, value interface{}) *SacredError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithContext 添加上下文信息
func (e *SacredError) WithContext(key, value string) *SacredError {
	if e.Context == nil {
		e.Context = make(map[string]string)
	}
	e.Context[key] = value
	return e
}

// WithSeverity 设置严重程度
func (e *SacredError) WithSeverity(severity ErrorSeverity) *SacredError {
	e.Severity = severity
	return e
}

// WithCategory 设置错误类别
func (e *SacredError) WithCategory(category ErrorCategory) *SacredError {
	e.Category = category
	return e
}

// IsRetryable 检查是否可重试
func (e *SacredError) IsRetryable() bool {
	return e.Category == ErrorCategoryTemporary || e.Category == ErrorCategoryNetwork
}

// IsCritical 检查是否为关键错误
func (e *SacredError) IsCritical() bool {
	return e.Severity == ErrorSeverityCritical || e.Severity == ErrorSeverityFatal
}

// ErrorSeverity 错误严重程度
type ErrorSeverity int

const (
	ErrorSeverityInfo     ErrorSeverity = iota + 1 // 信息
	ErrorSeverityWarning                           // 警告
	ErrorSeverityError                             // 错误
	ErrorSeverityCritical                          // 关键
	ErrorSeverityFatal                             // 致命
)

// String 返回严重程度字符串
func (es ErrorSeverity) String() string {
	switch es {
	case ErrorSeverityInfo:
		return "info"
	case ErrorSeverityWarning:
		return "warning"
	case ErrorSeverityError:
		return "error"
	case ErrorSeverityCritical:
		return "critical"
	case ErrorSeverityFatal:
		return "fatal"
	default:
		return "unknown"
	}
}

// ErrorCategory 错误类别
type ErrorCategory int

const (
	ErrorCategoryValidation ErrorCategory = iota + 1 // 验证错误
	ErrorCategoryBusiness                             // 业务错误
	ErrorCategorySystem                               // 系统错误
	ErrorCategoryNetwork                              // 网络错误
	ErrorCategoryDatabase                             // 数据库错误
	ErrorCategoryPermission                           // 权限错误
	ErrorCategoryResource                             // 资源错误
	ErrorCategoryTemporary                            // 临时错误
	ErrorCategoryConfiguration                        // 配置错误
	ErrorCategoryConcurrency                          // 并发错误
)

// String 返回类别字符串
func (ec ErrorCategory) String() string {
	switch ec {
	case ErrorCategoryValidation:
		return "validation"
	case ErrorCategoryBusiness:
		return "business"
	case ErrorCategorySystem:
		return "system"
	case ErrorCategoryNetwork:
		return "network"
	case ErrorCategoryDatabase:
		return "database"
	case ErrorCategoryPermission:
		return "permission"
	case ErrorCategoryResource:
		return "resource"
	case ErrorCategoryTemporary:
		return "temporary"
	case ErrorCategoryConfiguration:
		return "configuration"
	case ErrorCategoryConcurrency:
		return "concurrency"
	default:
		return "unknown"
	}
}

// ValidationError 验证错误
type ValidationError struct {
	*SacredError
	Field      string
	Value      interface{}
	Constraint string
	Rule       string
}

// NewValidationError 创建验证错误
func NewValidationError(field, constraint, rule string, value interface{}) *ValidationError {
	return &ValidationError{
		SacredError: &SacredError{
			Code:      "VALIDATION_ERROR",
			Message:   fmt.Sprintf("validation failed for field '%s': %s", field, constraint),
			Timestamp: time.Now(),
			Details:   make(map[string]interface{}),
			Context:   make(map[string]string),
			Severity:  ErrorSeverityError,
			Category:  ErrorCategoryValidation,
		},
		Field:      field,
		Value:      value,
		Constraint: constraint,
		Rule:       rule,
	}
}

// BusinessRuleError 业务规则错误
type BusinessRuleError struct {
	*SacredError
	Rule        string
	Violation   string
	Expected    interface{}
	Actual      interface{}
	Suggestion  string
}

// NewBusinessRuleError 创建业务规则错误
func NewBusinessRuleError(rule, violation string, expected, actual interface{}) *BusinessRuleError {
	return &BusinessRuleError{
		SacredError: &SacredError{
			Code:      "BUSINESS_RULE_ERROR",
			Message:   fmt.Sprintf("business rule violation: %s - %s", rule, violation),
			Timestamp: time.Now(),
			Details:   make(map[string]interface{}),
			Context:   make(map[string]string),
			Severity:  ErrorSeverityError,
			Category:  ErrorCategoryBusiness,
		},
		Rule:      rule,
		Violation: violation,
		Expected:  expected,
		Actual:    actual,
	}
}

// WithSuggestion 添加建议
func (e *BusinessRuleError) WithSuggestion(suggestion string) *BusinessRuleError {
	e.Suggestion = suggestion
	return e
}

// ConcurrencyError 并发错误
type ConcurrencyError struct {
	*SacredError
	Resource    string
	Operation   string
	ConflictID  string
	RetryAfter  time.Duration
	MaxRetries  int
	CurrentTry  int
}

// NewConcurrencyError 创建并发错误
func NewConcurrencyError(resource, operation, conflictID string) *ConcurrencyError {
	return &ConcurrencyError{
		SacredError: &SacredError{
			Code:      "CONCURRENCY_ERROR",
			Message:   fmt.Sprintf("concurrent access conflict on %s during %s", resource, operation),
			Timestamp: time.Now(),
			Details:   make(map[string]interface{}),
			Context:   make(map[string]string),
			Severity:  ErrorSeverityWarning,
			Category:  ErrorCategoryConcurrency,
		},
		Resource:   resource,
		Operation:  operation,
		ConflictID: conflictID,
		MaxRetries: 3,
		CurrentTry: 1,
	}
}

// WithRetryAfter 设置重试时间
func (e *ConcurrencyError) WithRetryAfter(duration time.Duration) *ConcurrencyError {
	e.RetryAfter = duration
	return e
}

// CanRetry 检查是否可以重试
func (e *ConcurrencyError) CanRetry() bool {
	return e.CurrentTry < e.MaxRetries
}

// IncrementTry 增加尝试次数
func (e *ConcurrencyError) IncrementTry() {
	e.CurrentTry++
}

// ConfigurationError 配置错误
type ConfigurationError struct {
	*SacredError
	ConfigKey    string
	ConfigValue  interface{}
	ExpectedType string
	ValidValues  []interface{}
}

// NewConfigurationError 创建配置错误
func NewConfigurationError(configKey string, configValue interface{}, expectedType string) *ConfigurationError {
	return &ConfigurationError{
		SacredError: &SacredError{
			Code:      "CONFIGURATION_ERROR",
			Message:   fmt.Sprintf("invalid configuration for key '%s': expected %s", configKey, expectedType),
			Timestamp: time.Now(),
			Details:   make(map[string]interface{}),
			Context:   make(map[string]string),
			Severity:  ErrorSeverityError,
			Category:  ErrorCategoryConfiguration,
		},
		ConfigKey:    configKey,
		ConfigValue:  configValue,
		ExpectedType: expectedType,
	}
}

// WithValidValues 设置有效值
func (e *ConfigurationError) WithValidValues(values ...interface{}) *ConfigurationError {
	e.ValidValues = values
	return e
}

// SystemError 系统错误
type SystemError struct {
	*SacredError
	Component   string
	Operation   string
	ErrorCode   int
	Recoverable bool
	RetryCount  int
	MaxRetries  int
	StackTrace  string
}

// NewSystemError 创建系统错误
func NewSystemError(component, operation string, errorCode int, cause error) *SystemError {
	return &SystemError{
		SacredError: &SacredError{
			Code:      "SYSTEM_ERROR",
			Message:   fmt.Sprintf("system error in %s during %s (code: %d)", component, operation, errorCode),
			Cause:     cause,
			Timestamp: time.Now(),
			Details:   make(map[string]interface{}),
			Context:   make(map[string]string),
			Severity:  ErrorSeverityCritical,
			Category:  ErrorCategorySystem,
		},
		Component:  component,
		Operation:  operation,
		ErrorCode:  errorCode,
		MaxRetries: 3,
	}
}

// SetRecoverable 设置是否可恢复
func (e *SystemError) SetRecoverable(recoverable bool) *SystemError {
	e.Recoverable = recoverable
	return e
}

// IncrementRetry 增加重试次数
func (e *SystemError) IncrementRetry() *SystemError {
	e.RetryCount++
	return e
}

// CanRetry 检查是否可以重试
func (e *SystemError) CanRetry() bool {
	return e.Recoverable && e.RetryCount < e.MaxRetries
}

// WithStackTrace 添加堆栈跟踪
func (e *SystemError) WithStackTrace(stackTrace string) *SystemError {
	e.StackTrace = stackTrace
	return e
}

// PermissionError 权限错误
type PermissionError struct {
	*SacredError
	UserID         string
	Resource       string
	RequiredRole   string
	CurrentRole    string
	RequiredPermissions []string
	CurrentPermissions  []string
}

// NewPermissionError 创建权限错误
func NewPermissionError(userID, resource, requiredRole, currentRole string) *PermissionError {
	return &PermissionError{
		SacredError: &SacredError{
			Code:      "PERMISSION_ERROR",
			Message:   fmt.Sprintf("permission denied for user %s on resource %s", userID, resource),
			Timestamp: time.Now(),
			Details:   make(map[string]interface{}),
			Context:   make(map[string]string),
			Severity:  ErrorSeverityError,
			Category:  ErrorCategoryPermission,
		},
		UserID:       userID,
		Resource:     resource,
		RequiredRole: requiredRole,
		CurrentRole:  currentRole,
	}
}

// WithPermissions 设置权限信息
func (e *PermissionError) WithPermissions(required, current []string) *PermissionError {
	e.RequiredPermissions = required
	e.CurrentPermissions = current
	return e
}

// ResourceError 资源错误
type ResourceError struct {
	*SacredError
	ResourceType string
	ResourceID   string
	Required     interface{}
	Available    interface{}
	Unit         string
}

// NewResourceError 创建资源错误
func NewResourceError(resourceType, resourceID string, required, available interface{}, unit string) *ResourceError {
	return &ResourceError{
		SacredError: &SacredError{
			Code:      "RESOURCE_ERROR",
			Message:   fmt.Sprintf("insufficient %s: required %v, available %v %s", resourceType, required, available, unit),
			Timestamp: time.Now(),
			Details:   make(map[string]interface{}),
			Context:   make(map[string]string),
			Severity:  ErrorSeverityError,
			Category:  ErrorCategoryResource,
		},
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Required:     required,
		Available:    available,
		Unit:         unit,
	}
}

// ErrorCollection 错误集合
type ErrorCollection struct {
	Errors    []error
	Context   string
	Timestamp time.Time
	Severity  ErrorSeverity
}

// NewErrorCollection 创建错误集合
func NewErrorCollection(context string) *ErrorCollection {
	return &ErrorCollection{
		Errors:    make([]error, 0),
		Context:   context,
		Timestamp: time.Now(),
		Severity:  ErrorSeverityError,
	}
}

// Add 添加错误
func (ec *ErrorCollection) Add(err error) {
	if err != nil {
		ec.Errors = append(ec.Errors, err)
		
		// 更新严重程度
		if sacredErr, ok := err.(*SacredError); ok {
			if sacredErr.Severity > ec.Severity {
				ec.Severity = sacredErr.Severity
			}
		}
	}
}

// HasErrors 检查是否有错误
func (ec *ErrorCollection) HasErrors() bool {
	return len(ec.Errors) > 0
}

// Count 获取错误数量
func (ec *ErrorCollection) Count() int {
	return len(ec.Errors)
}

// Error 实现error接口
func (ec *ErrorCollection) Error() string {
	if len(ec.Errors) == 0 {
		return "no errors"
	}
	
	var messages []string
	for i, err := range ec.Errors {
		messages = append(messages, fmt.Sprintf("%d: %v", i+1, err))
	}
	
	return fmt.Sprintf("multiple errors in %s: [%s]", ec.Context, strings.Join(messages, "; "))
}

// First 获取第一个错误
func (ec *ErrorCollection) First() error {
	if len(ec.Errors) > 0 {
		return ec.Errors[0]
	}
	return nil
}

// Last 获取最后一个错误
func (ec *ErrorCollection) Last() error {
	if len(ec.Errors) > 0 {
		return ec.Errors[len(ec.Errors)-1]
	}
	return nil
}

// FilterBySeverity 按严重程度过滤
func (ec *ErrorCollection) FilterBySeverity(severity ErrorSeverity) []error {
	var filtered []error
	for _, err := range ec.Errors {
		if sacredErr, ok := err.(*SacredError); ok {
			if sacredErr.Severity == severity {
				filtered = append(filtered, err)
			}
		}
	}
	return filtered
}

// FilterByCategory 按类别过滤
func (ec *ErrorCollection) FilterByCategory(category ErrorCategory) []error {
	var filtered []error
	for _, err := range ec.Errors {
		if sacredErr, ok := err.(*SacredError); ok {
			if sacredErr.Category == category {
				filtered = append(filtered, err)
			}
		}
	}
	return filtered
}

// 错误工厂函数

// NewSacredError 创建圣地错误
func NewSacredError(code, message string) *SacredError {
	return &SacredError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Details:   make(map[string]interface{}),
		Context:   make(map[string]string),
		Severity:  ErrorSeverityError,
		Category:  ErrorCategoryBusiness,
	}
}

// NewSacredErrorWithCause 创建带原因的圣地错误
func NewSacredErrorWithCause(code, message string, cause error) *SacredError {
	return &SacredError{
		Code:      code,
		Message:   message,
		Cause:     cause,
		Timestamp: time.Now(),
		Details:   make(map[string]interface{}),
		Context:   make(map[string]string),
		Severity:  ErrorSeverityError,
		Category:  ErrorCategoryBusiness,
	}
}

// WrapError 包装错误
func WrapError(err error, code, message string) *SacredError {
	return &SacredError{
		Code:      code,
		Message:   message,
		Cause:     err,
		Timestamp: time.Now(),
		Details:   make(map[string]interface{}),
		Context:   make(map[string]string),
		Severity:  ErrorSeverityError,
		Category:  ErrorCategorySystem,
	}
}

// 错误检查函数

// IsValidationError 检查是否为验证错误
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// IsBusinessRuleError 检查是否为业务规则错误
func IsBusinessRuleError(err error) bool {
	_, ok := err.(*BusinessRuleError)
	return ok
}

// IsConcurrencyError 检查是否为并发错误
func IsConcurrencyError(err error) bool {
	_, ok := err.(*ConcurrencyError)
	return ok
}

// IsConfigurationError 检查是否为配置错误
func IsConfigurationError(err error) bool {
	_, ok := err.(*ConfigurationError)
	return ok
}

// IsSystemError 检查是否为系统错误
func IsSystemError(err error) bool {
	_, ok := err.(*SystemError)
	return ok
}

// IsPermissionError 检查是否为权限错误
func IsPermissionError(err error) bool {
	_, ok := err.(*PermissionError)
	return ok
}

// IsResourceError 检查是否为资源错误
func IsResourceError(err error) bool {
	_, ok := err.(*ResourceError)
	return ok
}

// IsSacredError 检查是否为圣地错误
func IsSacredError(err error) bool {
	_, ok := err.(*SacredError)
	return ok
}

// 错误分类函数

// IsRetryableError 检查错误是否可重试
func IsRetryableError(err error) bool {
	if sacredErr, ok := err.(*SacredError); ok {
		return sacredErr.IsRetryable()
	}
	if sysErr, ok := err.(*SystemError); ok {
		return sysErr.CanRetry()
	}
	if concErr, ok := err.(*ConcurrencyError); ok {
		return concErr.CanRetry()
	}
	return false
}

// IsTemporaryError 检查错误是否为临时错误
func IsTemporaryError(err error) bool {
	switch err {
	case ErrServiceUnavailable, ErrTimeout, ErrResourceBusy:
		return true
	default:
		return IsRetryableError(err)
	}
}

// IsPermanentError 检查错误是否为永久错误
func IsPermanentError(err error) bool {
	switch err {
	case ErrPermissionDenied, ErrUnauthorized, ErrDataCorrupted:
		return true
	default:
		return IsValidationError(err) || IsConfigurationError(err)
	}
}

// IsCriticalError 检查错误是否为关键错误
func IsCriticalError(err error) bool {
	if sacredErr, ok := err.(*SacredError); ok {
		return sacredErr.IsCritical()
	}
	return false
}

// 辅助函数

// GetErrorCode 获取错误代码
func GetErrorCode(err error) string {
	if sacredErr, ok := err.(*SacredError); ok {
		return sacredErr.Code
	}
	return "UNKNOWN_ERROR"
}

// GetErrorSeverity 获取错误严重程度
func GetErrorSeverity(err error) ErrorSeverity {
	if sacredErr, ok := err.(*SacredError); ok {
		return sacredErr.Severity
	}
	return ErrorSeverityError
}

// GetErrorCategory 获取错误类别
func GetErrorCategory(err error) ErrorCategory {
	if sacredErr, ok := err.(*SacredError); ok {
		return sacredErr.Category
	}
	return ErrorCategorySystem
}

// GetErrorDetails 获取错误详情
func GetErrorDetails(err error) map[string]interface{} {
	if sacredErr, ok := err.(*SacredError); ok {
		return sacredErr.Details
	}
	return nil
}

// GetErrorContext 获取错误上下文
func GetErrorContext(err error) map[string]string {
	if sacredErr, ok := err.(*SacredError); ok {
		return sacredErr.Context
	}
	return nil
}

// FormatError 格式化错误信息
func FormatError(err error) string {
	if err == nil {
		return "no error"
	}
	
	if sacredErr, ok := err.(*SacredError); ok {
		var parts []string
		parts = append(parts, fmt.Sprintf("Code: %s", sacredErr.Code))
		parts = append(parts, fmt.Sprintf("Message: %s", sacredErr.Message))
		parts = append(parts, fmt.Sprintf("Severity: %s", sacredErr.Severity.String()))
		parts = append(parts, fmt.Sprintf("Category: %s", sacredErr.Category.String()))
		parts = append(parts, fmt.Sprintf("Time: %s", sacredErr.Timestamp.Format(time.RFC3339)))
		
		if len(sacredErr.Details) > 0 {
			parts = append(parts, fmt.Sprintf("Details: %+v", sacredErr.Details))
		}
		
		if len(sacredErr.Context) > 0 {
			parts = append(parts, fmt.Sprintf("Context: %+v", sacredErr.Context))
		}
		
		if sacredErr.Cause != nil {
			parts = append(parts, fmt.Sprintf("Cause: %v", sacredErr.Cause))
		}
		
		return strings.Join(parts, ", ")
	}
	
	return err.Error()
}

// CreateErrorResponse 创建错误响应
func CreateErrorResponse(err error) map[string]interface{} {
	response := map[string]interface{}{
		"error":     true,
		"message":   err.Error(),
		"timestamp": time.Now(),
	}
	
	if sacredErr, ok := err.(*SacredError); ok {
		response["code"] = sacredErr.Code
		response["severity"] = sacredErr.Severity.String()
		response["category"] = sacredErr.Category.String()
		response["retryable"] = sacredErr.IsRetryable()
		
		if len(sacredErr.Details) > 0 {
			response["details"] = sacredErr.Details
		}
		
		if len(sacredErr.Context) > 0 {
			response["context"] = sacredErr.Context
		}
	}
	
	return response
}