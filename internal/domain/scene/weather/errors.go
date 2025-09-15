package weather

import (
	"errors"
	"fmt"
	"time"
)

// 基础错误定义
var (
	// 天气类型相关错误
	ErrInvalidWeatherType      = errors.New("invalid weather type")
	ErrUnsupportedWeatherType  = errors.New("unsupported weather type")
	ErrWeatherTypeNotFound     = errors.New("weather type not found")
	ErrWeatherTypeConflict     = errors.New("weather type conflict")
	
	// 天气强度相关错误
	ErrInvalidWeatherIntensity = errors.New("invalid weather intensity")
	ErrIntensityOutOfRange     = errors.New("weather intensity out of range")
	ErrIntensityNotSupported   = errors.New("weather intensity not supported")
	
	// 天气状态相关错误
	ErrInvalidWeatherState     = errors.New("invalid weather state")
	ErrWeatherStateExpired     = errors.New("weather state expired")
	ErrWeatherStateNotActive   = errors.New("weather state not active")
	ErrWeatherStateConflict    = errors.New("weather state conflict")
	ErrWeatherStateNotFound    = errors.New("weather state not found")
	
	// 天气转换相关错误
	ErrInvalidWeatherTransition = errors.New("invalid weather transition")
	ErrWeatherTransitionBlocked = errors.New("weather transition blocked")
	ErrTransitionTooFrequent    = errors.New("weather transition too frequent")
	ErrTransitionNotAllowed     = errors.New("weather transition not allowed")
	
	// 天气效果相关错误
	ErrInvalidWeatherEffect    = errors.New("invalid weather effect")
	ErrWeatherEffectExpired    = errors.New("weather effect expired")
	ErrWeatherEffectNotActive  = errors.New("weather effect not active")
	ErrWeatherEffectConflict   = errors.New("weather effect conflict")
	ErrWeatherEffectNotFound   = errors.New("weather effect not found")
	ErrEffectMultiplierInvalid = errors.New("effect multiplier invalid")
	
	// 天气事件相关错误
	ErrInvalidWeatherEvent     = errors.New("invalid weather event")
	ErrWeatherEventExpired     = errors.New("weather event expired")
	ErrWeatherEventNotActive   = errors.New("weather event not active")
	ErrWeatherEventConflict    = errors.New("weather event conflict")
	ErrWeatherEventNotFound    = errors.New("weather event not found")
	ErrEventTriggerFailed      = errors.New("weather event trigger failed")
	ErrEventSeverityInvalid    = errors.New("weather event severity invalid")
	
	// 天气预报相关错误
	ErrInvalidForecastPeriod   = errors.New("invalid forecast period")
	ErrForecastNotAvailable    = errors.New("weather forecast not available")
	ErrForecastExpired         = errors.New("weather forecast expired")
	ErrForecastGenerationFailed = errors.New("weather forecast generation failed")
	ErrForecastAccuracyLow     = errors.New("weather forecast accuracy too low")
	
	// 季节模式相关错误
	ErrInvalidSeason           = errors.New("invalid season")
	ErrSeasonalPatternNotFound = errors.New("seasonal pattern not found")
	ErrSeasonalPatternInvalid  = errors.New("seasonal pattern invalid")
	ErrSeasonTransitionFailed  = errors.New("season transition failed")
	
	// 气候区域相关错误
	ErrInvalidClimateZone      = errors.New("invalid climate zone")
	ErrClimateZoneNotFound     = errors.New("climate zone not found")
	ErrClimateZoneConflict     = errors.New("climate zone conflict")
	
	// 时间相关错误
	ErrInvalidTimeRange        = errors.New("invalid time range")
	ErrInvalidDuration         = errors.New("invalid duration")
	ErrInvalidChangeInterval   = errors.New("invalid change interval")
	ErrTimeoutExceeded         = errors.New("timeout exceeded")
	
	// 配置相关错误
	ErrInvalidConfiguration    = errors.New("invalid weather configuration")
	ErrConfigurationNotFound   = errors.New("weather configuration not found")
	ErrConfigurationConflict   = errors.New("weather configuration conflict")
	
	// 数据相关错误
	ErrDataCorrupted           = errors.New("weather data corrupted")
	ErrDataNotFound            = errors.New("weather data not found")
	ErrDataValidationFailed    = errors.New("weather data validation failed")
	ErrDataSerializationFailed = errors.New("weather data serialization failed")
	
	// 系统相关错误
	ErrSystemNotInitialized    = errors.New("weather system not initialized")
	ErrSystemOverloaded        = errors.New("weather system overloaded")
	ErrSystemMaintenance       = errors.New("weather system under maintenance")
	ErrResourceExhausted       = errors.New("weather system resources exhausted")
	
	// 并发相关错误
	ErrConcurrentModification  = errors.New("concurrent weather modification")
	ErrLockAcquisitionFailed   = errors.New("weather lock acquisition failed")
	ErrVersionMismatch         = errors.New("weather version mismatch")
	
	// 权限相关错误
	ErrPermissionDenied        = errors.New("weather operation permission denied")
	ErrUnauthorizedAccess      = errors.New("unauthorized weather access")
	ErrInsufficientPrivileges  = errors.New("insufficient weather privileges")
)

// WeatherError 天气错误结构体
type WeatherError struct {
	Code      string
	Message   string
	Cause     error
	Context   map[string]interface{}
	Timestamp int64
}

// Error 实现error接口
func (e *WeatherError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 返回原始错误
func (e *WeatherError) Unwrap() error {
	return e.Cause
}

// WithContext 添加上下文信息
func (e *WeatherError) WithContext(key string, value interface{}) *WeatherError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// GetContext 获取上下文信息
func (e *WeatherError) GetContext(key string) interface{} {
	if e.Context == nil {
		return nil
	}
	return e.Context[key]
}

// NewWeatherError 创建天气错误
func NewWeatherError(code, message string, cause error) *WeatherError {
	return &WeatherError{
		Code:      code,
		Message:   message,
		Cause:     cause,
		Context:   make(map[string]interface{}),
		Timestamp: getCurrentTimestamp(),
	}
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Value   interface{}
	Rule    string
	Message string
}

// Error 实现error接口
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s (value: %v, rule: %s)", e.Field, e.Message, e.Value, e.Rule)
}

// NewValidationError 创建验证错误
func NewValidationError(field string, value interface{}, rule, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Rule:    rule,
		Message: message,
	}
}

// BusinessRuleError 业务规则错误
type BusinessRuleError struct {
	Rule        string
	Description string
	Violation   string
	Context     map[string]interface{}
}

// Error 实现error接口
func (e *BusinessRuleError) Error() string {
	return fmt.Sprintf("business rule violation: %s - %s (%s)", e.Rule, e.Description, e.Violation)
}

// NewBusinessRuleError 创建业务规则错误
func NewBusinessRuleError(rule, description, violation string) *BusinessRuleError {
	return &BusinessRuleError{
		Rule:        rule,
		Description: description,
		Violation:   violation,
		Context:     make(map[string]interface{}),
	}
}

// ConcurrencyError 并发错误
type ConcurrencyError struct {
	Operation   string
	Resource    string
	ConflictID  string
	RetryCount  int
	MaxRetries  int
	LastAttempt int64
}

// Error 实现error接口
func (e *ConcurrencyError) Error() string {
	return fmt.Sprintf("concurrency error in operation '%s' on resource '%s': conflict with %s (retry %d/%d)", 
		e.Operation, e.Resource, e.ConflictID, e.RetryCount, e.MaxRetries)
}

// CanRetry 检查是否可以重试
func (e *ConcurrencyError) CanRetry() bool {
	return e.RetryCount < e.MaxRetries
}

// NewConcurrencyError 创建并发错误
func NewConcurrencyError(operation, resource, conflictID string, retryCount, maxRetries int) *ConcurrencyError {
	return &ConcurrencyError{
		Operation:   operation,
		Resource:    resource,
		ConflictID:  conflictID,
		RetryCount:  retryCount,
		MaxRetries:  maxRetries,
		LastAttempt: getCurrentTimestamp(),
	}
}

// ConfigurationError 配置错误
type ConfigurationError struct {
	ConfigType  string
	ConfigKey   string
	ConfigValue interface{}
	Expected    string
	Actual      string
	Suggestion  string
}

// Error 实现error接口
func (e *ConfigurationError) Error() string {
	return fmt.Sprintf("configuration error in %s.%s: expected %s, got %s (value: %v). Suggestion: %s", 
		e.ConfigType, e.ConfigKey, e.Expected, e.Actual, e.ConfigValue, e.Suggestion)
}

// NewConfigurationError 创建配置错误
func NewConfigurationError(configType, configKey string, configValue interface{}, expected, actual, suggestion string) *ConfigurationError {
	return &ConfigurationError{
		ConfigType:  configType,
		ConfigKey:   configKey,
		ConfigValue: configValue,
		Expected:    expected,
		Actual:      actual,
		Suggestion:  suggestion,
	}
}

// SystemError 系统错误
type SystemError struct {
	Component   string
	Operation   string
	ErrorCode   string
	Description string
	Cause       error
	Severity    ErrorSeverity
	Recoverable bool
	Timestamp   int64
}

// Error 实现error接口
func (e *SystemError) Error() string {
	msg := fmt.Sprintf("system error in %s.%s [%s]: %s (severity: %s, recoverable: %t)", 
		e.Component, e.Operation, e.ErrorCode, e.Description, e.Severity, e.Recoverable)
	if e.Cause != nil {
		msg += fmt.Sprintf(", cause: %v", e.Cause)
	}
	return msg
}

// Unwrap 返回原始错误
func (e *SystemError) Unwrap() error {
	return e.Cause
}

// NewSystemError 创建系统错误
func NewSystemError(component, operation, errorCode, description string, cause error, severity ErrorSeverity, recoverable bool) *SystemError {
	return &SystemError{
		Component:   component,
		Operation:   operation,
		ErrorCode:   errorCode,
		Description: description,
		Cause:       cause,
		Severity:    severity,
		Recoverable: recoverable,
		Timestamp:   getCurrentTimestamp(),
	}
}

// ErrorSeverity 错误严重程度
type ErrorSeverity int

const (
	ErrorSeverityLow ErrorSeverity = iota + 1
	ErrorSeverityMedium
	ErrorSeverityHigh
	ErrorSeverityCritical
)

// String 返回严重程度字符串
func (es ErrorSeverity) String() string {
	switch es {
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

// ErrorCollection 错误集合
type ErrorCollection struct {
	Errors    []error
	Context   string
	Timestamp int64
}

// Error 实现error接口
func (ec *ErrorCollection) Error() string {
	if len(ec.Errors) == 0 {
		return "no errors"
	}
	
	if len(ec.Errors) == 1 {
		return fmt.Sprintf("error in %s: %v", ec.Context, ec.Errors[0])
	}
	
	msg := fmt.Sprintf("%d errors in %s:", len(ec.Errors), ec.Context)
	for i, err := range ec.Errors {
		msg += fmt.Sprintf("\n  %d. %v", i+1, err)
	}
	return msg
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

// Clear 清空错误
func (ec *ErrorCollection) Clear() {
	ec.Errors = ec.Errors[:0]
}

// NewErrorCollection 创建错误集合
func NewErrorCollection(context string) *ErrorCollection {
	return &ErrorCollection{
		Errors:    make([]error, 0),
		Context:   context,
		Timestamp: getCurrentTimestamp(),
	}
}

// 错误工厂函数

// NewInvalidWeatherTypeError 创建无效天气类型错误
func NewInvalidWeatherTypeError(weatherType interface{}) *WeatherError {
	return NewWeatherError("INVALID_WEATHER_TYPE", "Invalid weather type provided", ErrInvalidWeatherType).
		WithContext("weather_type", weatherType)
}

// NewInvalidWeatherIntensityError 创建无效天气强度错误
func NewInvalidWeatherIntensityError(intensity interface{}) *WeatherError {
	return NewWeatherError("INVALID_WEATHER_INTENSITY", "Invalid weather intensity provided", ErrInvalidWeatherIntensity).
		WithContext("intensity", intensity)
}

// NewWeatherStateExpiredError 创建天气状态过期错误
func NewWeatherStateExpiredError(stateID string, expiredAt int64) *WeatherError {
	return NewWeatherError("WEATHER_STATE_EXPIRED", "Weather state has expired", ErrWeatherStateExpired).
		WithContext("state_id", stateID).
		WithContext("expired_at", expiredAt)
}

// NewWeatherTransitionBlockedError 创建天气转换被阻止错误
func NewWeatherTransitionBlockedError(from, to WeatherType, reason string) *WeatherError {
	return NewWeatherError("WEATHER_TRANSITION_BLOCKED", "Weather transition is blocked", ErrWeatherTransitionBlocked).
		WithContext("from_weather", from.String()).
		WithContext("to_weather", to.String()).
		WithContext("reason", reason)
}

// NewForecastGenerationFailedError 创建预报生成失败错误
func NewForecastGenerationFailedError(sceneID string, hours int, cause error) *WeatherError {
	return NewWeatherError("FORECAST_GENERATION_FAILED", "Failed to generate weather forecast", cause).
		WithContext("scene_id", sceneID).
		WithContext("hours", hours)
}

// NewClimateZoneNotFoundError 创建气候区域未找到错误
func NewClimateZoneNotFoundError(zoneID string) *WeatherError {
	return NewWeatherError("CLIMATE_ZONE_NOT_FOUND", "Climate zone not found", ErrClimateZoneNotFound).
		WithContext("zone_id", zoneID)
}

// NewWeatherSystemOverloadedError 创建天气系统过载错误
func NewWeatherSystemOverloadedError(currentLoad, maxLoad int) *WeatherError {
	return NewWeatherError("WEATHER_SYSTEM_OVERLOADED", "Weather system is overloaded", ErrSystemOverloaded).
		WithContext("current_load", currentLoad).
		WithContext("max_load", maxLoad)
}

// NewConcurrentWeatherModificationError 创建并发天气修改错误
func NewConcurrentWeatherModificationError(sceneID, operation string, conflictVersion int) *WeatherError {
	return NewWeatherError("CONCURRENT_WEATHER_MODIFICATION", "Concurrent weather modification detected", ErrConcurrentModification).
		WithContext("scene_id", sceneID).
		WithContext("operation", operation).
		WithContext("conflict_version", conflictVersion)
}

// 错误检查函数

// IsWeatherError 检查是否为天气错误
func IsWeatherError(err error) bool {
	_, ok := err.(*WeatherError)
	return ok
}

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

// IsSystemError 检查是否为系统错误
func IsSystemError(err error) bool {
	_, ok := err.(*SystemError)
	return ok
}

// IsRecoverableError 检查错误是否可恢复
func IsRecoverableError(err error) bool {
	if sysErr, ok := err.(*SystemError); ok {
		return sysErr.Recoverable
	}
	if concErr, ok := err.(*ConcurrencyError); ok {
		return concErr.CanRetry()
	}
	return false
}

// GetErrorSeverity 获取错误严重程度
func GetErrorSeverity(err error) ErrorSeverity {
	if sysErr, ok := err.(*SystemError); ok {
		return sysErr.Severity
	}
	return ErrorSeverityMedium // 默认中等严重程度
}

// 辅助函数

// getCurrentTimestamp 获取当前时间戳
func getCurrentTimestamp() int64 {
	return getCurrentTime().Unix()
}

// getCurrentTime 获取当前时间（可用于测试时的时间注入）
var getCurrentTime = func() int64 {
	return time.Now().Unix()
}