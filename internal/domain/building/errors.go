package building

import (
	"fmt"
	"time"
)

// BuildingError 建筑错误接口
type BuildingError interface {
	error
	GetCode() string
	GetMessage() string
	GetSeverity() ErrorSeverity
	GetTimestamp() time.Time
	GetContext() map[string]interface{}
	SetContext(key string, value interface{})
	IsRetryable() bool
	GetRetryAfter() time.Duration
}

// ErrorSeverity 错误严重程度
type ErrorSeverity int32

const (
	ErrorSeverityLow      ErrorSeverity = iota + 1 // 低严重程度
	ErrorSeverityMedium                            // 中等严重程度
	ErrorSeverityHigh                              // 高严重程度
	ErrorSeverityCritical                          // 关键严重程度
)

// String 返回错误严重程度的字符串表示
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

// IsValid 检查错误严重程度是否有效
func (es ErrorSeverity) IsValid() bool {
	return es >= ErrorSeverityLow && es <= ErrorSeverityCritical
}

// BaseBuildingError 建筑错误基础结构
type BaseBuildingError struct {
	code       string
	message    string
	severity   ErrorSeverity
	timestamp  time.Time
	context    map[string]interface{}
	retryable  bool
	retryAfter time.Duration
}

// NewBuildingError 创建新建筑错误
func NewBuildingError(code, message string, severity ErrorSeverity) *BaseBuildingError {
	return &BaseBuildingError{
		code:      code,
		message:   message,
		severity:  severity,
		timestamp: time.Now(),
		context:   make(map[string]interface{}),
		retryable: false,
	}
}

// Error 实现error接口
func (e *BaseBuildingError) Error() string {
	return fmt.Sprintf("[%s] %s (severity: %s)", e.code, e.message, e.severity.String())
}

// GetCode 获取错误代码
func (e *BaseBuildingError) GetCode() string {
	return e.code
}

// GetMessage 获取错误消息
func (e *BaseBuildingError) GetMessage() string {
	return e.message
}

// GetSeverity 获取错误严重程度
func (e *BaseBuildingError) GetSeverity() ErrorSeverity {
	return e.severity
}

// GetTimestamp 获取错误时间戳
func (e *BaseBuildingError) GetTimestamp() time.Time {
	return e.timestamp
}

// GetContext 获取错误上下文
func (e *BaseBuildingError) GetContext() map[string]interface{} {
	return e.context
}

// SetContext 设置错误上下文
func (e *BaseBuildingError) SetContext(key string, value interface{}) {
	if e.context == nil {
		e.context = make(map[string]interface{})
	}
	e.context[key] = value
}

// IsRetryable 检查错误是否可重试
func (e *BaseBuildingError) IsRetryable() bool {
	return e.retryable
}

// GetRetryAfter 获取重试间隔
func (e *BaseBuildingError) GetRetryAfter() time.Duration {
	return e.retryAfter
}

// SetRetryable 设置错误可重试
func (e *BaseBuildingError) SetRetryable(retryable bool, retryAfter time.Duration) {
	e.retryable = retryable
	e.retryAfter = retryAfter
}

// 具体错误类型

// BuildingNotFoundError 建筑未找到错误
type BuildingNotFoundError struct {
	*BaseBuildingError
	BuildingID string `json:"building_id"`
}

// NewBuildingNotFoundError 创建建筑未找到错误
func NewBuildingNotFoundError(buildingID string) *BuildingNotFoundError {
	err := &BuildingNotFoundError{
		BaseBuildingError: NewBuildingError(
			ErrCodeBuildingNotFound,
			fmt.Sprintf("building with ID %s not found", buildingID),
			ErrorSeverityHigh,
		),
		BuildingID: buildingID,
	}
	err.SetContext("building_id", buildingID)
	return err
}

// ConstructionNotFoundError 建造未找到错误
type ConstructionNotFoundError struct {
	*BaseBuildingError
	ConstructionID string `json:"construction_id"`
}

// NewConstructionNotFoundError 创建建造未找到错误
func NewConstructionNotFoundError(constructionID string) *ConstructionNotFoundError {
	err := &ConstructionNotFoundError{
		BaseBuildingError: NewBuildingError(
			ErrCodeConstructionNotFound,
			fmt.Sprintf("construction with ID %s not found", constructionID),
			ErrorSeverityHigh,
		),
		ConstructionID: constructionID,
	}
	err.SetContext("construction_id", constructionID)
	return err
}

// UpgradeNotFoundError 升级未找到错误
type UpgradeNotFoundError struct {
	*BaseBuildingError
	UpgradeID string `json:"upgrade_id"`
}

// NewUpgradeNotFoundError 创建升级未找到错误
func NewUpgradeNotFoundError(upgradeID string) *UpgradeNotFoundError {
	err := &UpgradeNotFoundError{
		BaseBuildingError: NewBuildingError(
			ErrCodeUpgradeNotFound,
			fmt.Sprintf("upgrade with ID %s not found", upgradeID),
			ErrorSeverityHigh,
		),
		UpgradeID: upgradeID,
	}
	err.SetContext("upgrade_id", upgradeID)
	return err
}

// BlueprintNotFoundError 蓝图未找到错误
type BlueprintNotFoundError struct {
	*BaseBuildingError
	BlueprintID string `json:"blueprint_id"`
}

// NewBlueprintNotFoundError 创建蓝图未找到错误
func NewBlueprintNotFoundError(blueprintID string) *BlueprintNotFoundError {
	err := &BlueprintNotFoundError{
		BaseBuildingError: NewBuildingError(
			ErrCodeBlueprintNotFound,
			fmt.Sprintf("blueprint with ID %s not found", blueprintID),
			ErrorSeverityHigh,
		),
		BlueprintID: blueprintID,
	}
	err.SetContext("blueprint_id", blueprintID)
	return err
}

// InvalidBuildingStateError 无效建筑状态错误
type InvalidBuildingStateError struct {
	*BaseBuildingError
	BuildingID    string         `json:"building_id"`
	CurrentState  BuildingStatus `json:"current_state"`
	ExpectedState BuildingStatus `json:"expected_state"`
	Operation     string         `json:"operation"`
}

// NewInvalidBuildingStateError 创建无效建筑状态错误
func NewInvalidBuildingStateError(buildingID string, currentState, expectedState BuildingStatus, operation string) *InvalidBuildingStateError {
	err := &InvalidBuildingStateError{
		BaseBuildingError: NewBuildingError(
			ErrCodeInvalidBuildingState,
			fmt.Sprintf("building %s is in state %s, expected %s for operation %s", buildingID, currentState.String(), expectedState.String(), operation),
			ErrorSeverityHigh,
		),
		BuildingID:    buildingID,
		CurrentState:  currentState,
		ExpectedState: expectedState,
		Operation:     operation,
	}
	err.SetContext("building_id", buildingID)
	err.SetContext("current_state", currentState.String())
	err.SetContext("expected_state", expectedState.String())
	err.SetContext("operation", operation)
	return err
}

// InsufficientResourcesError 资源不足错误
type InsufficientResourcesError struct {
	*BaseBuildingError
	ResourceType string `json:"resource_type"`
	Required     int64  `json:"required"`
	Available    int64  `json:"available"`
	OwnerID      uint64 `json:"owner_id"`
}

// NewInsufficientResourcesError 创建资源不足错误
func NewInsufficientResourcesError(resourceType string, required, available int64, ownerID uint64) *InsufficientResourcesError {
	err := &InsufficientResourcesError{
		BaseBuildingError: NewBuildingError(
			ErrCodeInsufficientResources,
			fmt.Sprintf("insufficient %s: required %d, available %d", resourceType, required, available),
			ErrorSeverityHigh,
		),
		ResourceType: resourceType,
		Required:     required,
		Available:    available,
		OwnerID:      ownerID,
	}
	err.SetContext("resource_type", resourceType)
	err.SetContext("required", required)
	err.SetContext("available", available)
	err.SetContext("owner_id", ownerID)
	return err
}

// PositionOccupiedError 位置被占用错误
type PositionOccupiedError struct {
	*BaseBuildingError
	Position          *Position `json:"position"`
	OccupyingBuilding string    `json:"occupying_building"`
}

// NewPositionOccupiedError 创建位置被占用错误
func NewPositionOccupiedError(position *Position, occupyingBuilding string) *PositionOccupiedError {
	err := &PositionOccupiedError{
		BaseBuildingError: NewBuildingError(
			ErrCodePositionOccupied,
			fmt.Sprintf("position (%d,%d,%d) is occupied by building %s", position.X, position.Y, position.Z, occupyingBuilding),
			ErrorSeverityHigh,
		),
		Position:          position,
		OccupyingBuilding: occupyingBuilding,
	}
	err.SetContext("position", fmt.Sprintf("%d,%d,%d", position.X, position.Y, position.Z))
	err.SetContext("occupying_building", occupyingBuilding)
	return err
}

// ConstructionFailedError 建造失败错误
type ConstructionFailedError struct {
	*BaseBuildingError
	BuildingID     string `json:"building_id"`
	ConstructionID string `json:"construction_id"`
	Reason         string `json:"reason"`
}

// NewConstructionFailedError 创建建造失败错误
func NewConstructionFailedError(buildingID, constructionID, reason string) *ConstructionFailedError {
	err := &ConstructionFailedError{
		BaseBuildingError: NewBuildingError(
			ErrCodeConstructionFailed,
			fmt.Sprintf("construction %s for building %s failed: %s", constructionID, buildingID, reason),
			ErrorSeverityHigh,
		),
		BuildingID:     buildingID,
		ConstructionID: constructionID,
		Reason:         reason,
	}
	err.SetContext("building_id", buildingID)
	err.SetContext("construction_id", constructionID)
	err.SetContext("reason", reason)
	err.SetRetryable(true, 5*time.Minute)
	return err
}

// UpgradeFailedError 升级失败错误
type UpgradeFailedError struct {
	*BaseBuildingError
	BuildingID string `json:"building_id"`
	UpgradeID  string `json:"upgrade_id"`
	Reason     string `json:"reason"`
}

// NewUpgradeFailedError 创建升级失败错误
func NewUpgradeFailedError(buildingID, upgradeID, reason string) *UpgradeFailedError {
	err := &UpgradeFailedError{
		BaseBuildingError: NewBuildingError(
			ErrCodeUpgradeFailed,
			fmt.Sprintf("upgrade %s for building %s failed: %s", upgradeID, buildingID, reason),
			ErrorSeverityHigh,
		),
		BuildingID: buildingID,
		UpgradeID:  upgradeID,
		Reason:     reason,
	}
	err.SetContext("building_id", buildingID)
	err.SetContext("upgrade_id", upgradeID)
	err.SetContext("reason", reason)
	err.SetRetryable(true, 5*time.Minute)
	return err
}

// RepairFailedError 修复失败错误
type RepairFailedError struct {
	*BaseBuildingError
	BuildingID string `json:"building_id"`
	Reason     string `json:"reason"`
}

// NewRepairFailedError 创建修复失败错误
func NewRepairFailedError(buildingID, reason string) *RepairFailedError {
	err := &RepairFailedError{
		BaseBuildingError: NewBuildingError(
			ErrCodeRepairFailed,
			fmt.Sprintf("repair for building %s failed: %s", buildingID, reason),
			ErrorSeverityMedium,
		),
		BuildingID: buildingID,
		Reason:     reason,
	}
	err.SetContext("building_id", buildingID)
	err.SetContext("reason", reason)
	err.SetRetryable(true, 1*time.Minute)
	return err
}

// DestroyFailedError 摧毁失败错误
type DestroyFailedError struct {
	*BaseBuildingError
	BuildingID string `json:"building_id"`
	Reason     string `json:"reason"`
}

// NewDestroyFailedError 创建摧毁失败错误
func NewDestroyFailedError(buildingID, reason string) *DestroyFailedError {
	err := &DestroyFailedError{
		BaseBuildingError: NewBuildingError(
			ErrCodeDestroyFailed,
			fmt.Sprintf("destroy for building %s failed: %s", buildingID, reason),
			ErrorSeverityMedium,
		),
		BuildingID: buildingID,
		Reason:     reason,
	}
	err.SetContext("building_id", buildingID)
	err.SetContext("reason", reason)
	return err
}

// WorkerNotAvailableError 工人不可用错误
type WorkerNotAvailableError struct {
	*BaseBuildingError
	WorkerID uint64 `json:"worker_id"`
	Reason   string `json:"reason"`
}

// NewWorkerNotAvailableError 创建工人不可用错误
func NewWorkerNotAvailableError(workerID uint64, reason string) *WorkerNotAvailableError {
	err := &WorkerNotAvailableError{
		BaseBuildingError: NewBuildingError(
			ErrCodeWorkerNotAvailable,
			fmt.Sprintf("worker %d is not available: %s", workerID, reason),
			ErrorSeverityMedium,
		),
		WorkerID: workerID,
		Reason:   reason,
	}
	err.SetContext("worker_id", workerID)
	err.SetContext("reason", reason)
	err.SetRetryable(true, 10*time.Minute)
	return err
}

// InvalidInputError 无效输入错误
type InvalidInputError struct {
	*BaseBuildingError
	Field string `json:"field"`
	Value string `json:"value"`
}

// NewInvalidInputError 创建无效输入错误
func NewInvalidInputError(field, value, reason string) *InvalidInputError {
	err := &InvalidInputError{
		BaseBuildingError: NewBuildingError(
			ErrCodeInvalidInput,
			fmt.Sprintf("invalid input for field %s with value %s: %s", field, value, reason),
			ErrorSeverityHigh,
		),
		Field: field,
		Value: value,
	}
	err.SetContext("field", field)
	err.SetContext("value", value)
	return err
}

// RepositoryError 仓储错误
type RepositoryError struct {
	*BaseBuildingError
	Operation string `json:"operation"`
	Entity    string `json:"entity"`
}

// NewRepositoryError 创建仓储错误
func NewRepositoryError(operation, entity, reason string) *RepositoryError {
	err := &RepositoryError{
		BaseBuildingError: NewBuildingError(
			ErrCodeRepositoryError,
			fmt.Sprintf("repository error during %s operation on %s: %s", operation, entity, reason),
			ErrorSeverityMedium,
		),
		Operation: operation,
		Entity:    entity,
	}
	err.SetContext("operation", operation)
	err.SetContext("entity", entity)
	err.SetRetryable(true, 30*time.Second)
	return err
}

// ConcurrencyError 并发错误
type ConcurrencyError struct {
	*BaseBuildingError
	ResourceID string `json:"resource_id"`
	Operation  string `json:"operation"`
}

// NewConcurrencyError 创建并发错误
func NewConcurrencyError(resourceID, operation string) *ConcurrencyError {
	err := &ConcurrencyError{
		BaseBuildingError: NewBuildingError(
			ErrCodeConcurrencyError,
			fmt.Sprintf("concurrency conflict for resource %s during operation %s", resourceID, operation),
			ErrorSeverityMedium,
		),
		ResourceID: resourceID,
		Operation:  operation,
	}
	err.SetContext("resource_id", resourceID)
	err.SetContext("operation", operation)
	err.SetRetryable(true, 1*time.Second)
	return err
}

// 错误代码常量

const (
	// 通用错误
	ErrCodeInvalidInput     = "BUILDING_INVALID_INPUT"
	ErrCodeRepositoryError  = "BUILDING_REPOSITORY_ERROR"
	ErrCodeConcurrencyError = "BUILDING_CONCURRENCY_ERROR"
	ErrCodeInvalidOperation = "BUILDING_INVALID_OPERATION"

	// 建筑相关错误
	ErrCodeBuildingNotFound      = "BUILDING_NOT_FOUND"
	ErrCodeInvalidBuildingState  = "BUILDING_INVALID_STATE"
	ErrCodePositionOccupied      = "BUILDING_POSITION_OCCUPIED"
	ErrCodeInsufficientResources = "BUILDING_INSUFFICIENT_RESOURCES"

	// 建造相关错误
	ErrCodeConstructionNotFound = "CONSTRUCTION_NOT_FOUND"
	ErrCodeConstructionFailed   = "CONSTRUCTION_FAILED"

	// 升级相关错误
	ErrCodeUpgradeNotFound = "UPGRADE_NOT_FOUND"
	ErrCodeUpgradeFailed   = "UPGRADE_FAILED"
	ErrCodeInvalidUpgrade  = "UPGRADE_INVALID"

	// 维护相关错误
	ErrCodeRepairFailed  = "REPAIR_FAILED"
	ErrCodeDestroyFailed = "DESTROY_FAILED"

	// 工人相关错误
	ErrCodeWorkerNotFound     = "WORKER_NOT_FOUND"
	ErrCodeWorkerNotAvailable = "WORKER_NOT_AVAILABLE"

	// 蓝图相关错误
	ErrCodeBlueprintNotFound = "BLUEPRINT_NOT_FOUND"
	ErrCodeBlueprintInvalid  = "BLUEPRINT_INVALID"

	// 系统错误
	ErrCodeSystemError  = "BUILDING_SYSTEM_ERROR"
	ErrCodeConfigError  = "BUILDING_CONFIG_ERROR"
	ErrCodeNetworkError = "BUILDING_NETWORK_ERROR"
	ErrCodeTimeoutError = "BUILDING_TIMEOUT_ERROR"
)

// 工具函数

// IsBuildingError 检查是否为建筑错误
func IsBuildingError(err error) bool {
	_, ok := err.(BuildingError)
	return ok
}

// GetBuildingError 获取建筑错误
func GetBuildingError(err error) (BuildingError, bool) {
	buildingErr, ok := err.(BuildingError)
	return buildingErr, ok
}

// IsRetryableError 检查错误是否可重试
func IsRetryableError(err error) bool {
	if buildingErr, ok := GetBuildingError(err); ok {
		return buildingErr.IsRetryable()
	}
	return false
}

// GetErrorSeverity 获取错误严重程度
func GetErrorSeverity(err error) ErrorSeverity {
	if buildingErr, ok := GetBuildingError(err); ok {
		return buildingErr.GetSeverity()
	}
	return ErrorSeverityLow
}

// GetErrorCode 获取错误代码
func GetErrorCode(err error) string {
	if buildingErr, ok := GetBuildingError(err); ok {
		return buildingErr.GetCode()
	}
	return "UNKNOWN_ERROR"
}

// 错误分类函数

// IsNotFoundError 检查是否为未找到错误
func IsNotFoundError(err error) bool {
	code := GetErrorCode(err)
	return code == ErrCodeBuildingNotFound ||
		code == ErrCodeConstructionNotFound ||
		code == ErrCodeUpgradeNotFound ||
		code == ErrCodeBlueprintNotFound ||
		code == ErrCodeWorkerNotFound
}

// IsValidationError 检查是否为验证错误
func IsValidationError(err error) bool {
	code := GetErrorCode(err)
	return code == ErrCodeInvalidInput ||
		code == ErrCodeInvalidBuildingState ||
		code == ErrCodeInvalidUpgrade ||
		code == ErrCodeBlueprintInvalid
}

// IsResourceError 检查是否为资源错误
func IsResourceError(err error) bool {
	code := GetErrorCode(err)
	return code == ErrCodeInsufficientResources ||
		code == ErrCodePositionOccupied ||
		code == ErrCodeWorkerNotAvailable
}

// IsSystemError 检查是否为系统错误
func IsSystemError(err error) bool {
	code := GetErrorCode(err)
	return code == ErrCodeSystemError ||
		code == ErrCodeConfigError ||
		code == ErrCodeNetworkError ||
		code == ErrCodeTimeoutError ||
		code == ErrCodeRepositoryError ||
		code == ErrCodeConcurrencyError
}

// IsOperationError 检查是否为操作错误
func IsOperationError(err error) bool {
	code := GetErrorCode(err)
	return code == ErrCodeConstructionFailed ||
		code == ErrCodeUpgradeFailed ||
		code == ErrCodeRepairFailed ||
		code == ErrCodeDestroyFailed
}

// 错误恢复策略

// ErrorRecoveryStrategy 错误恢复策略
type ErrorRecoveryStrategy int32

const (
	RecoveryStrategyNone                ErrorRecoveryStrategy = iota + 1 // 无恢复策略
	RecoveryStrategyRetry                                                // 重试
	RecoveryStrategyFallback                                             // 回退
	RecoveryStrategyCircuitBreaker                                       // 熔断器
	RecoveryStrategyGracefulDegradation                                  // 优雅降级
)

// String 返回恢复策略的字符串表示
func (ers ErrorRecoveryStrategy) String() string {
	switch ers {
	case RecoveryStrategyNone:
		return "none"
	case RecoveryStrategyRetry:
		return "retry"
	case RecoveryStrategyFallback:
		return "fallback"
	case RecoveryStrategyCircuitBreaker:
		return "circuit_breaker"
	case RecoveryStrategyGracefulDegradation:
		return "graceful_degradation"
	default:
		return "unknown"
	}
}

// GetRecoveryStrategy 获取错误的恢复策略
func GetRecoveryStrategy(err error) ErrorRecoveryStrategy {
	code := GetErrorCode(err)
	severity := GetErrorSeverity(err)

	// 根据错误代码和严重程度确定恢复策略
	switch {
	case IsSystemError(err) && severity == ErrorSeverityCritical:
		return RecoveryStrategyCircuitBreaker
	case IsSystemError(err) && severity == ErrorSeverityHigh:
		return RecoveryStrategyGracefulDegradation
	case IsResourceError(err):
		return RecoveryStrategyRetry
	case IsOperationError(err):
		return RecoveryStrategyFallback
	case code == ErrCodeConcurrencyError:
		return RecoveryStrategyRetry
	case code == ErrCodeRepositoryError:
		return RecoveryStrategyRetry
	default:
		return RecoveryStrategyNone
	}
}

// 错误统计

// ErrorStatistics 错误统计
type ErrorStatistics struct {
	TotalErrors      int64                   `json:"total_errors"`
	ErrorsByCode     map[string]int64        `json:"errors_by_code"`
	ErrorsBySeverity map[ErrorSeverity]int64 `json:"errors_by_severity"`
	ErrorsByHour     map[string]int64        `json:"errors_by_hour"`
	ErrorsByDay      map[string]int64        `json:"errors_by_day"`
	RetryableErrors  int64                   `json:"retryable_errors"`
	LastErrorTime    time.Time               `json:"last_error_time"`
	UpdatedAt        time.Time               `json:"updated_at"`
}

// NewErrorStatistics 创建新错误统计
func NewErrorStatistics() *ErrorStatistics {
	return &ErrorStatistics{
		TotalErrors:      0,
		ErrorsByCode:     make(map[string]int64),
		ErrorsBySeverity: make(map[ErrorSeverity]int64),
		ErrorsByHour:     make(map[string]int64),
		ErrorsByDay:      make(map[string]int64),
		RetryableErrors:  0,
		UpdatedAt:        time.Now(),
	}
}

// AddError 添加错误到统计
func (es *ErrorStatistics) AddError(err error) {
	es.TotalErrors++

	if buildingErr, ok := GetBuildingError(err); ok {
		code := buildingErr.GetCode()
		severity := buildingErr.GetSeverity()
		timestamp := buildingErr.GetTimestamp()

		es.ErrorsByCode[code]++
		es.ErrorsBySeverity[severity]++

		hourKey := timestamp.Format("2006-01-02-15")
		dayKey := timestamp.Format("2006-01-02")
		es.ErrorsByHour[hourKey]++
		es.ErrorsByDay[dayKey]++

		if buildingErr.IsRetryable() {
			es.RetryableErrors++
		}

		if timestamp.After(es.LastErrorTime) {
			es.LastErrorTime = timestamp
		}
	} else {
		// 非建筑错误
		es.ErrorsByCode["UNKNOWN_ERROR"]++
		es.ErrorsBySeverity[ErrorSeverityLow]++

		now := time.Now()
		hourKey := now.Format("2006-01-02-15")
		dayKey := now.Format("2006-01-02")
		es.ErrorsByHour[hourKey]++
		es.ErrorsByDay[dayKey]++

		if now.After(es.LastErrorTime) {
			es.LastErrorTime = now
		}
	}

	es.UpdatedAt = time.Now()
}

// GetMostFrequentError 获取最频繁的错误
func (es *ErrorStatistics) GetMostFrequentError() string {
	maxCount := int64(0)
	mostFrequent := ""

	for code, count := range es.ErrorsByCode {
		if count > maxCount {
			maxCount = count
			mostFrequent = code
		}
	}

	return mostFrequent
}

// GetErrorRate 获取错误率（每小时）
func (es *ErrorStatistics) GetErrorRate() float64 {
	if es.TotalErrors == 0 {
		return 0.0
	}

	// 计算最近24小时的错误数
	recentErrors := int64(0)
	now := time.Now()
	for i := 0; i < 24; i++ {
		hourKey := now.Add(-time.Duration(i) * time.Hour).Format("2006-01-02-15")
		recentErrors += es.ErrorsByHour[hourKey]
	}

	return float64(recentErrors) / 24.0
}

// GetRetryableErrorRate 获取可重试错误率
func (es *ErrorStatistics) GetRetryableErrorRate() float64 {
	if es.TotalErrors == 0 {
		return 0.0
	}
	return float64(es.RetryableErrors) / float64(es.TotalErrors) * 100.0
}

// 辅助函数

// WrapError 包装错误
func WrapError(err error, code, message string, severity ErrorSeverity) BuildingError {
	buildingErr := NewBuildingError(code, message, severity)
	buildingErr.SetContext("wrapped_error", err.Error())
	return buildingErr
}

// ChainErrors 链接错误
func ChainErrors(errors []error) BuildingError {
	if len(errors) == 0 {
		return nil
	}

	if len(errors) == 1 {
		if buildingErr, ok := GetBuildingError(errors[0]); ok {
			return buildingErr
		}
		return WrapError(errors[0], "CHAINED_ERROR", errors[0].Error(), ErrorSeverityMedium)
	}

	// 创建链式错误
	messages := make([]string, len(errors))
	for i, err := range errors {
		messages[i] = err.Error()
	}

	chainedErr := NewBuildingError(
		"CHAINED_ERROR",
		fmt.Sprintf("multiple errors occurred: %v", messages),
		ErrorSeverityHigh,
	)

	for i, err := range errors {
		chainedErr.SetContext(fmt.Sprintf("error_%d", i), err.Error())
	}

	return chainedErr
}

// ValidateErrorCode 验证错误代码
func ValidateErrorCode(code string) bool {
	validCodes := []string{
		ErrCodeInvalidInput,
		ErrCodeRepositoryError,
		ErrCodeConcurrencyError,
		ErrCodeInvalidOperation,
		ErrCodeBuildingNotFound,
		ErrCodeInvalidBuildingState,
		ErrCodePositionOccupied,
		ErrCodeInsufficientResources,
		ErrCodeConstructionNotFound,
		ErrCodeConstructionFailed,
		ErrCodeUpgradeNotFound,
		ErrCodeUpgradeFailed,
		ErrCodeInvalidUpgrade,
		ErrCodeRepairFailed,
		ErrCodeDestroyFailed,
		ErrCodeWorkerNotFound,
		ErrCodeWorkerNotAvailable,
		ErrCodeBlueprintNotFound,
		ErrCodeBlueprintInvalid,
		ErrCodeSystemError,
		ErrCodeConfigError,
		ErrCodeNetworkError,
		ErrCodeTimeoutError,
	}

	for _, validCode := range validCodes {
		if code == validCode {
			return true
		}
	}
	return false
}

// FormatError 格式化错误信息
func FormatError(err error) string {
	if buildingErr, ok := GetBuildingError(err); ok {
		return fmt.Sprintf("[%s] %s (severity: %s, timestamp: %s)",
			buildingErr.GetCode(),
			buildingErr.GetMessage(),
			buildingErr.GetSeverity().String(),
			buildingErr.GetTimestamp().Format(time.RFC3339),
		)
	}
	return err.Error()
}

// GetErrorDetails 获取错误详情
func GetErrorDetails(err error) map[string]interface{} {
	details := make(map[string]interface{})

	if buildingErr, ok := GetBuildingError(err); ok {
		details["code"] = buildingErr.GetCode()
		details["message"] = buildingErr.GetMessage()
		details["severity"] = buildingErr.GetSeverity().String()
		details["timestamp"] = buildingErr.GetTimestamp()
		details["retryable"] = buildingErr.IsRetryable()
		details["retry_after"] = buildingErr.GetRetryAfter()
		details["context"] = buildingErr.GetContext()
		details["recovery_strategy"] = GetRecoveryStrategy(err).String()
	} else {
		details["code"] = "UNKNOWN_ERROR"
		details["message"] = err.Error()
		details["severity"] = ErrorSeverityLow.String()
		details["timestamp"] = time.Now()
		details["retryable"] = false
		details["retry_after"] = time.Duration(0)
		details["context"] = make(map[string]interface{})
		details["recovery_strategy"] = RecoveryStrategyNone.String()
	}

	return details
}
