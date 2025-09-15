package plant

import (
	"fmt"
	"strings"
	"time"
)

// 基础错误变量
var (
	// 种子相关错误
	ErrInvalidSeedType     = fmt.Errorf("invalid seed type")
	ErrSeedNotFound        = fmt.Errorf("seed not found")
	ErrInsufficientSeeds   = fmt.Errorf("insufficient seeds")
	ErrSeedExpired         = fmt.Errorf("seed has expired")
	ErrSeedAlreadyPlanted  = fmt.Errorf("seed already planted")
	ErrSeedNotCompatible   = fmt.Errorf("seed not compatible with soil")

	// 作物相关错误
	ErrCropNotFound        = fmt.Errorf("crop not found")
	ErrCropNotMature       = fmt.Errorf("crop is not mature")
	ErrCropAlreadyHarvested = fmt.Errorf("crop already harvested")
	ErrCropDead            = fmt.Errorf("crop is dead")
	ErrCropDiseased        = fmt.Errorf("crop is diseased")
	ErrCropNotHealthy      = fmt.Errorf("crop is not healthy")
	ErrCropOverwatered     = fmt.Errorf("crop is overwatered")
	ErrCropUnderwatered    = fmt.Errorf("crop is underwatered")

	// 地块相关错误
	ErrPlotNotFound        = fmt.Errorf("plot not found")
	ErrPlotOccupied        = fmt.Errorf("plot is occupied")
	ErrPlotNotEmpty        = fmt.Errorf("plot is not empty")
	ErrPlotNotReady        = fmt.Errorf("plot is not ready for planting")
	ErrPlotTooSmall        = fmt.Errorf("plot is too small")
	ErrPlotDamaged         = fmt.Errorf("plot is damaged")

	// 土壤相关错误
	ErrInvalidSoilType     = fmt.Errorf("invalid soil type")
	ErrSoilTooAcidic       = fmt.Errorf("soil is too acidic")
	ErrSoilTooAlkaline     = fmt.Errorf("soil is too alkaline")
	ErrSoilPoorFertility   = fmt.Errorf("soil has poor fertility")
	ErrSoilContaminated    = fmt.Errorf("soil is contaminated")
	ErrSoilTooWet          = fmt.Errorf("soil is too wet")
	ErrSoilTooDry          = fmt.Errorf("soil is too dry")

	// 肥料相关错误
	ErrInvalidFertilizer   = fmt.Errorf("invalid fertilizer")
	ErrInsufficientFertilizer = fmt.Errorf("insufficient fertilizer")
	ErrFertilizerExpired   = fmt.Errorf("fertilizer has expired")
	ErrOverFertilization   = fmt.Errorf("over fertilization")
	ErrFertilizerNotCompatible = fmt.Errorf("fertilizer not compatible")

	// 工具相关错误
	ErrToolNotFound        = fmt.Errorf("tool not found")
	ErrToolBroken          = fmt.Errorf("tool is broken")
	ErrToolNotSuitable     = fmt.Errorf("tool is not suitable for this operation")
	ErrToolInUse           = fmt.Errorf("tool is in use")
	ErrInsufficientDurability = fmt.Errorf("insufficient tool durability")

	// 农场相关错误
	ErrFarmNotFound        = fmt.Errorf("farm not found")
	ErrFarmFull            = fmt.Errorf("farm is full")
	ErrFarmLocked          = fmt.Errorf("farm is locked")
	ErrInsufficientSpace   = fmt.Errorf("insufficient space")
	ErrFarmNotOwned        = fmt.Errorf("farm is not owned by player")

	// 季节相关错误
	ErrInvalidSeason       = fmt.Errorf("invalid season")
	ErrSeasonNotSuitable   = fmt.Errorf("season is not suitable for this crop")
	ErrSeasonTransition    = fmt.Errorf("season transition in progress")

	// 天气相关错误
	ErrBadWeather          = fmt.Errorf("bad weather conditions")
	ErrWeatherNotSuitable  = fmt.Errorf("weather not suitable for operation")
	ErrExtremeWeather      = fmt.Errorf("extreme weather conditions")

	// 时间相关错误
	ErrInvalidTime         = fmt.Errorf("invalid time")
	ErrTooEarly            = fmt.Errorf("too early for this operation")
	ErrTooLate             = fmt.Errorf("too late for this operation")
	ErrTimeExpired         = fmt.Errorf("time has expired")

	// 资源相关错误
	ErrInsufficientResources = fmt.Errorf("insufficient resources")
	ErrInsufficientWater   = fmt.Errorf("insufficient water")
	ErrInsufficientGold    = fmt.Errorf("insufficient gold")
	ErrResourceNotFound    = fmt.Errorf("resource not found")

	// 权限相关错误
	ErrPermissionDenied    = fmt.Errorf("permission denied")
	ErrUnauthorized        = fmt.Errorf("unauthorized operation")
	ErrAccessRestricted    = fmt.Errorf("access restricted")

	// 配置相关错误
	ErrInvalidConfiguration = fmt.Errorf("invalid configuration")
	ErrConfigurationNotFound = fmt.Errorf("configuration not found")
	ErrConfigurationCorrupted = fmt.Errorf("configuration corrupted")

	// 数据相关错误
	ErrDataCorrupted       = fmt.Errorf("data corrupted")
	ErrDataNotFound        = fmt.Errorf("data not found")
	ErrDataInconsistent    = fmt.Errorf("data inconsistent")

	// 系统相关错误
	ErrSystemError         = fmt.Errorf("system error")
	ErrServiceUnavailable  = fmt.Errorf("service unavailable")
	ErrTimeout             = fmt.Errorf("operation timeout")

	// 并发相关错误
	ErrConcurrentModification = fmt.Errorf("concurrent modification")
	ErrResourceLocked      = fmt.Errorf("resource is locked")
	ErrDeadlock            = fmt.Errorf("deadlock detected")
)

// PlantError 种植系统错误
type PlantError struct {
	Code      string
	Message   string
	Details   map[string]interface{}
	Cause     error
	Timestamp time.Time
	Context   map[string]string
}

// Error 实现error接口
func (e *PlantError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap 返回原始错误
func (e *PlantError) Unwrap() error {
	return e.Cause
}

// WithDetail 添加详细信息
func (e *PlantError) WithDetail(key string, value interface{}) *PlantError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithContext 添加上下文信息
func (e *PlantError) WithContext(key, value string) *PlantError {
	if e.Context == nil {
		e.Context = make(map[string]string)
	}
	e.Context[key] = value
	return e
}

// ValidationError 验证错误
type ValidationError struct {
	*PlantError
	Field      string
	Value      interface{}
	Constraint string
	Rule       string
}

// NewValidationError 创建验证错误
func NewValidationError(field, constraint, rule string, value interface{}) *ValidationError {
	return &ValidationError{
		PlantError: &PlantError{
			Code:      "VALIDATION_ERROR",
			Message:   fmt.Sprintf("validation failed for field '%s': %s", field, constraint),
			Timestamp: time.Now(),
			Details:   make(map[string]interface{}),
			Context:   make(map[string]string),
		},
		Field:      field,
		Value:      value,
		Constraint: constraint,
		Rule:       rule,
	}
}

// BusinessRuleError 业务规则错误
type BusinessRuleError struct {
	*PlantError
	Rule        string
	Violation   string
	Expected    interface{}
	Actual      interface{}
	Suggestion  string
}

// NewBusinessRuleError 创建业务规则错误
func NewBusinessRuleError(rule, violation string, expected, actual interface{}) *BusinessRuleError {
	return &BusinessRuleError{
		PlantError: &PlantError{
			Code:      "BUSINESS_RULE_ERROR",
			Message:   fmt.Sprintf("business rule violation: %s - %s", rule, violation),
			Timestamp: time.Now(),
			Details:   make(map[string]interface{}),
			Context:   make(map[string]string),
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
	*PlantError
	Resource    string
	Operation   string
	ConflictID  string
	RetryAfter  time.Duration
}

// NewConcurrencyError 创建并发错误
func NewConcurrencyError(resource, operation, conflictID string) *ConcurrencyError {
	return &ConcurrencyError{
		PlantError: &PlantError{
			Code:      "CONCURRENCY_ERROR",
			Message:   fmt.Sprintf("concurrent access conflict on %s during %s", resource, operation),
			Timestamp: time.Now(),
			Details:   make(map[string]interface{}),
			Context:   make(map[string]string),
		},
		Resource:   resource,
		Operation:  operation,
		ConflictID: conflictID,
	}
}

// WithRetryAfter 设置重试时间
func (e *ConcurrencyError) WithRetryAfter(duration time.Duration) *ConcurrencyError {
	e.RetryAfter = duration
	return e
}

// ConfigurationError 配置错误
type ConfigurationError struct {
	*PlantError
	ConfigKey   string
	ConfigValue interface{}
	ExpectedType string
	ValidValues []interface{}
}

// NewConfigurationError 创建配置错误
func NewConfigurationError(configKey string, configValue interface{}, expectedType string) *ConfigurationError {
	return &ConfigurationError{
		PlantError: &PlantError{
			Code:      "CONFIGURATION_ERROR",
			Message:   fmt.Sprintf("invalid configuration for key '%s': expected %s", configKey, expectedType),
			Timestamp: time.Now(),
			Details:   make(map[string]interface{}),
			Context:   make(map[string]string),
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
	*PlantError
	Component   string
	Operation   string
	ErrorCode   int
	Recoverable bool
	RetryCount  int
	MaxRetries  int
}

// NewSystemError 创建系统错误
func NewSystemError(component, operation string, errorCode int, cause error) *SystemError {
	return &SystemError{
		PlantError: &PlantError{
			Code:      "SYSTEM_ERROR",
			Message:   fmt.Sprintf("system error in %s during %s (code: %d)", component, operation, errorCode),
			Cause:     cause,
			Timestamp: time.Now(),
			Details:   make(map[string]interface{}),
			Context:   make(map[string]string),
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

// ErrorCollection 错误集合
type ErrorCollection struct {
	Errors    []error
	Context   string
	Timestamp time.Time
}

// NewErrorCollection 创建错误集合
func NewErrorCollection(context string) *ErrorCollection {
	return &ErrorCollection{
		Errors:    make([]error, 0),
		Context:   context,
		Timestamp: time.Now(),
	}
}

// Add 添加错误
func (ec *ErrorCollection) Add(err error) {
	if err != nil {
		ec.Errors = append(ec.Errors, err)
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

// 错误工厂函数

// NewPlantError 创建种植错误
func NewPlantError(code, message string) *PlantError {
	return &PlantError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Details:   make(map[string]interface{}),
		Context:   make(map[string]string),
	}
}

// NewPlantErrorWithCause 创建带原因的种植错误
func NewPlantErrorWithCause(code, message string, cause error) *PlantError {
	return &PlantError{
		Code:      code,
		Message:   message,
		Cause:     cause,
		Timestamp: time.Now(),
		Details:   make(map[string]interface{}),
		Context:   make(map[string]string),
	}
}

// WrapError 包装错误
func WrapError(err error, code, message string) *PlantError {
	return &PlantError{
		Code:      code,
		Message:   message,
		Cause:     err,
		Timestamp: time.Now(),
		Details:   make(map[string]interface{}),
		Context:   make(map[string]string),
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

// IsPlantError 检查是否为种植错误
func IsPlantError(err error) bool {
	_, ok := err.(*PlantError)
	return ok
}

// 错误分类函数

// IsRetryableError 检查错误是否可重试
func IsRetryableError(err error) bool {
	if sysErr, ok := err.(*SystemError); ok {
		return sysErr.CanRetry()
	}
	if concErr, ok := err.(*ConcurrencyError); ok {
		return concErr.RetryAfter > 0
	}
	return false
}

// IsTemporaryError 检查错误是否为临时错误
func IsTemporaryError(err error) bool {
	switch err {
	case ErrServiceUnavailable, ErrTimeout, ErrResourceLocked:
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

// 辅助函数

// GetErrorCode 获取错误代码
func GetErrorCode(err error) string {
	if plantErr, ok := err.(*PlantError); ok {
		return plantErr.Code
	}
	return "UNKNOWN_ERROR"
}

// GetErrorDetails 获取错误详情
func GetErrorDetails(err error) map[string]interface{} {
	if plantErr, ok := err.(*PlantError); ok {
		return plantErr.Details
	}
	return nil
}

// GetErrorContext 获取错误上下文
func GetErrorContext(err error) map[string]string {
	if plantErr, ok := err.(*PlantError); ok {
		return plantErr.Context
	}
	return nil
}

// FormatError 格式化错误信息
func FormatError(err error) string {
	if err == nil {
		return "no error"
	}
	
	if plantErr, ok := err.(*PlantError); ok {
		var parts []string
		parts = append(parts, fmt.Sprintf("Code: %s", plantErr.Code))
		parts = append(parts, fmt.Sprintf("Message: %s", plantErr.Message))
		parts = append(parts, fmt.Sprintf("Time: %s", plantErr.Timestamp.Format(time.RFC3339)))
		
		if len(plantErr.Details) > 0 {
			parts = append(parts, fmt.Sprintf("Details: %+v", plantErr.Details))
		}
		
		if len(plantErr.Context) > 0 {
			parts = append(parts, fmt.Sprintf("Context: %+v", plantErr.Context))
		}
		
		if plantErr.Cause != nil {
			parts = append(parts, fmt.Sprintf("Cause: %v", plantErr.Cause))
		}
		
		return strings.Join(parts, ", ")
	}
	
	return err.Error()
}