package npc

import (
	"fmt"
	"strings"
	"time"
)

// 基础错误变量
var (
	// NPC相关错误
	ErrNPCNotFound         = fmt.Errorf("NPC not found")
	ErrNPCAlreadyExists    = fmt.Errorf("NPC already exists")
	ErrNPCInvalidID        = fmt.Errorf("invalid NPC ID")
	ErrNPCInvalidName      = fmt.Errorf("invalid NPC name")
	ErrNPCInvalidType      = fmt.Errorf("invalid NPC type")
	ErrNPCInvalidStatus    = fmt.Errorf("invalid NPC status")
	ErrNPCInvalidLocation  = fmt.Errorf("invalid NPC location")
	ErrNPCNotActive        = fmt.Errorf("NPC is not active")
	ErrNPCBusy             = fmt.Errorf("NPC is busy")
	ErrNPCUnavailable      = fmt.Errorf("NPC is unavailable")
	ErrNPCMaxInteractions  = fmt.Errorf("NPC has reached maximum interactions")
	
	// 对话相关错误
	ErrDialogueNotFound       = fmt.Errorf("dialogue not found")
	ErrDialogueAlreadyExists  = fmt.Errorf("dialogue already exists")
	ErrDialogueInvalidID      = fmt.Errorf("invalid dialogue ID")
	ErrDialogueInvalidType    = fmt.Errorf("invalid dialogue type")
	ErrDialogueNotAvailable   = fmt.Errorf("dialogue not available")
	ErrDialogueConditionFailed = fmt.Errorf("dialogue condition not met")
	ErrDialogueSessionExpired = fmt.Errorf("dialogue session expired")
	ErrDialogueSessionNotFound = fmt.Errorf("dialogue session not found")
	ErrDialogueInProgress     = fmt.Errorf("dialogue already in progress")
	ErrDialogueNodeNotFound   = fmt.Errorf("dialogue node not found")
	ErrDialogueInvalidChoice  = fmt.Errorf("invalid dialogue choice")
	
	// 任务相关错误
	ErrQuestNotFound          = fmt.Errorf("quest not found")
	ErrQuestAlreadyExists     = fmt.Errorf("quest already exists")
	ErrQuestInvalidID         = fmt.Errorf("invalid quest ID")
	ErrQuestInvalidType       = fmt.Errorf("invalid quest type")
	ErrQuestNotAvailable      = fmt.Errorf("quest not available")
	ErrQuestAlreadyAccepted   = fmt.Errorf("quest already accepted")
	ErrQuestAlreadyCompleted  = fmt.Errorf("quest already completed")
	ErrQuestNotAccepted       = fmt.Errorf("quest not accepted")
	ErrQuestNotCompleted      = fmt.Errorf("quest not completed")
	ErrQuestConditionFailed   = fmt.Errorf("quest condition not met")
	ErrQuestObjectiveNotFound = fmt.Errorf("quest objective not found")
	ErrQuestRewardClaimed     = fmt.Errorf("quest reward already claimed")
	ErrQuestExpired           = fmt.Errorf("quest has expired")
	ErrQuestCooldown          = fmt.Errorf("quest is on cooldown")
	ErrQuestMaxAttempts       = fmt.Errorf("quest maximum attempts reached")
	
	// 商店相关错误
	ErrShopNotFound        = fmt.Errorf("shop not found")
	ErrShopAlreadyExists   = fmt.Errorf("shop already exists")
	ErrShopInvalidID       = fmt.Errorf("invalid shop ID")
	ErrShopNotOpen         = fmt.Errorf("shop is not open")
	ErrShopItemNotFound    = fmt.Errorf("shop item not found")
	ErrShopItemOutOfStock  = fmt.Errorf("shop item out of stock")
	ErrShopInsufficientFunds = fmt.Errorf("insufficient funds")
	ErrShopInvalidQuantity = fmt.Errorf("invalid quantity")
	ErrShopInvalidPrice    = fmt.Errorf("invalid price")
	ErrShopTransactionFailed = fmt.Errorf("shop transaction failed")
	ErrShopInventoryFull   = fmt.Errorf("shop inventory is full")
	
	// 关系相关错误
	ErrRelationshipNotFound    = fmt.Errorf("relationship not found")
	ErrRelationshipAlreadyExists = fmt.Errorf("relationship already exists")
	ErrRelationshipInvalidValue = fmt.Errorf("invalid relationship value")
	ErrRelationshipInvalidLevel = fmt.Errorf("invalid relationship level")
	ErrRelationshipMaxValue    = fmt.Errorf("relationship value at maximum")
	ErrRelationshipMinValue    = fmt.Errorf("relationship value at minimum")
	ErrRelationshipLocked      = fmt.Errorf("relationship is locked")
	
	// 行为相关错误
	ErrBehaviorNotFound     = fmt.Errorf("behavior not found")
	ErrBehaviorInvalidType  = fmt.Errorf("invalid behavior type")
	ErrBehaviorInvalidState = fmt.Errorf("invalid behavior state")
	ErrBehaviorCooldown     = fmt.Errorf("behavior is on cooldown")
	ErrBehaviorConditionFailed = fmt.Errorf("behavior condition not met")
	
	// 位置相关错误
	ErrLocationInvalid     = fmt.Errorf("invalid location")
	ErrLocationOutOfBounds = fmt.Errorf("location out of bounds")
	ErrLocationNotAccessible = fmt.Errorf("location not accessible")
	ErrLocationOccupied    = fmt.Errorf("location is occupied")
	
	// 权限相关错误
	ErrPermissionDenied    = fmt.Errorf("permission denied")
	ErrUnauthorized        = fmt.Errorf("unauthorized access")
	ErrInsufficientLevel   = fmt.Errorf("insufficient level")
	ErrInsufficientReputation = fmt.Errorf("insufficient reputation")
	
	// 系统相关错误
	ErrSystemBusy          = fmt.Errorf("system is busy")
	ErrSystemMaintenance   = fmt.Errorf("system under maintenance")
	ErrRateLimitExceeded   = fmt.Errorf("rate limit exceeded")
	ErrTimeout             = fmt.Errorf("operation timeout")
	ErrConcurrencyConflict = fmt.Errorf("concurrency conflict")
)

// NPCError NPC错误类型
type NPCError struct {
	Code      string
	Message   string
	Details   map[string]interface{}
	Cause     error
	Timestamp time.Time
	Context   map[string]string
}

// Error 实现error接口
func (e *NPCError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap 解包错误
func (e *NPCError) Unwrap() error {
	return e.Cause
}

// WithDetail 添加详细信息
func (e *NPCError) WithDetail(key string, value interface{}) *NPCError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithContext 添加上下文信息
func (e *NPCError) WithContext(key, value string) *NPCError {
	if e.Context == nil {
		e.Context = make(map[string]string)
	}
	e.Context[key] = value
	return e
}

// WithCause 添加原因错误
func (e *NPCError) WithCause(cause error) *NPCError {
	e.Cause = cause
	return e
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

// ValidationErrors 多个验证错误
type ValidationErrors struct {
	Errors []ValidationError
}

// Error 实现error接口
func (e *ValidationErrors) Error() string {
	if len(e.Errors) == 0 {
		return "no validation errors"
	}
	
	var messages []string
	for _, err := range e.Errors {
		messages = append(messages, err.Error())
	}
	return fmt.Sprintf("validation errors: %s", strings.Join(messages, "; "))
}

// Add 添加验证错误
func (e *ValidationErrors) Add(field, rule, message string, value interface{}) {
	e.Errors = append(e.Errors, ValidationError{
		Field:   field,
		Value:   value,
		Rule:    rule,
		Message: message,
	})
}

// HasErrors 是否有错误
func (e *ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
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
	return fmt.Sprintf("business rule violation '%s': %s (%s)", e.Rule, e.Description, e.Violation)
}

// WithContext 添加上下文
func (e *BusinessRuleError) WithContext(key string, value interface{}) *BusinessRuleError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// ConcurrencyError 并发错误
type ConcurrencyError struct {
	Resource    string
	Operation   string
	ConflictID  string
	ExpectedVersion int
	ActualVersion   int
}

// Error 实现error接口
func (e *ConcurrencyError) Error() string {
	return fmt.Sprintf("concurrency conflict on %s during %s: expected version %d, got %d (conflict ID: %s)", 
		e.Resource, e.Operation, e.ExpectedVersion, e.ActualVersion, e.ConflictID)
}

// TimeoutError 超时错误
type TimeoutError struct {
	Operation string
	Timeout   time.Duration
	Elapsed   time.Duration
}

// Error 实现error接口
func (e *TimeoutError) Error() string {
	return fmt.Sprintf("operation '%s' timed out after %v (timeout: %v)", e.Operation, e.Elapsed, e.Timeout)
}

// RateLimitError 限流错误
type RateLimitError struct {
	Resource    string
	Limit       int
	Window      time.Duration
	RetryAfter  time.Duration
	CurrentRate int
}

// Error 实现error接口
func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limit exceeded for %s: %d/%d requests in %v, retry after %v", 
		e.Resource, e.CurrentRate, e.Limit, e.Window, e.RetryAfter)
}

// 错误工厂函数

// NewNPCError 创建NPC错误
func NewNPCError(code, message string) *NPCError {
	return &NPCError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// NewValidationError 创建验证错误
func NewValidationError(field, rule, message string, value interface{}) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Rule:    rule,
		Message: message,
	}
}

// NewValidationErrors 创建验证错误集合
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make([]ValidationError, 0),
	}
}

// NewBusinessRuleError 创建业务规则错误
func NewBusinessRuleError(rule, description, violation string) *BusinessRuleError {
	return &BusinessRuleError{
		Rule:        rule,
		Description: description,
		Violation:   violation,
	}
}

// NewConcurrencyError 创建并发错误
func NewConcurrencyError(resource, operation, conflictID string, expectedVersion, actualVersion int) *ConcurrencyError {
	return &ConcurrencyError{
		Resource:        resource,
		Operation:       operation,
		ConflictID:      conflictID,
		ExpectedVersion: expectedVersion,
		ActualVersion:   actualVersion,
	}
}

// NewTimeoutError 创建超时错误
func NewTimeoutError(operation string, timeout, elapsed time.Duration) *TimeoutError {
	return &TimeoutError{
		Operation: operation,
		Timeout:   timeout,
		Elapsed:   elapsed,
	}
}

// NewRateLimitError 创建限流错误
func NewRateLimitError(resource string, limit int, window, retryAfter time.Duration, currentRate int) *RateLimitError {
	return &RateLimitError{
		Resource:    resource,
		Limit:       limit,
		Window:      window,
		RetryAfter:  retryAfter,
		CurrentRate: currentRate,
	}
}

// 错误检查函数

// IsNPCError 检查是否为NPC错误
func IsNPCError(err error) bool {
	_, ok := err.(*NPCError)
	return ok
}

// IsValidationError 检查是否为验证错误
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	if ok {
		return true
	}
	_, ok = err.(*ValidationErrors)
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

// IsTimeoutError 检查是否为超时错误
func IsTimeoutError(err error) bool {
	_, ok := err.(*TimeoutError)
	return ok
}

// IsRateLimitError 检查是否为限流错误
func IsRateLimitError(err error) bool {
	_, ok := err.(*RateLimitError)
	return ok
}

// 错误分类函数

// ErrorCategory 错误类别
type ErrorCategory string

const (
	ErrorCategoryValidation   ErrorCategory = "validation"
	ErrorCategoryBusinessRule ErrorCategory = "business_rule"
	ErrorCategoryNotFound     ErrorCategory = "not_found"
	ErrorCategoryConflict     ErrorCategory = "conflict"
	ErrorCategoryPermission   ErrorCategory = "permission"
	ErrorCategorySystem       ErrorCategory = "system"
	ErrorCategoryTimeout      ErrorCategory = "timeout"
	ErrorCategoryRateLimit    ErrorCategory = "rate_limit"
	ErrorCategoryUnknown      ErrorCategory = "unknown"
)

// CategorizeError 错误分类
func CategorizeError(err error) ErrorCategory {
	if err == nil {
		return ErrorCategoryUnknown
	}
	
	switch {
	case IsValidationError(err):
		return ErrorCategoryValidation
	case IsBusinessRuleError(err):
		return ErrorCategoryBusinessRule
	case IsConcurrencyError(err):
		return ErrorCategoryConflict
	case IsTimeoutError(err):
		return ErrorCategoryTimeout
	case IsRateLimitError(err):
		return ErrorCategoryRateLimit
	default:
		// 根据错误消息进行分类
		msg := err.Error()
		switch {
		case strings.Contains(msg, "not found"):
			return ErrorCategoryNotFound
		case strings.Contains(msg, "permission"), strings.Contains(msg, "unauthorized"):
			return ErrorCategoryPermission
		case strings.Contains(msg, "system"), strings.Contains(msg, "maintenance"):
			return ErrorCategorySystem
		default:
			return ErrorCategoryUnknown
		}
	}
}

// 错误恢复策略

// RecoveryStrategy 恢复策略
type RecoveryStrategy string

const (
	RecoveryStrategyRetry     RecoveryStrategy = "retry"
	RecoveryStrategyFallback  RecoveryStrategy = "fallback"
	RecoveryStrategyCircuit   RecoveryStrategy = "circuit_breaker"
	RecoveryStrategyIgnore    RecoveryStrategy = "ignore"
	RecoveryStrategyEscalate  RecoveryStrategy = "escalate"
)

// GetRecoveryStrategy 获取恢复策略
func GetRecoveryStrategy(err error) RecoveryStrategy {
	category := CategorizeError(err)
	
	switch category {
	case ErrorCategoryTimeout, ErrorCategoryRateLimit:
		return RecoveryStrategyRetry
	case ErrorCategorySystem:
		return RecoveryStrategyCircuit
	case ErrorCategoryValidation, ErrorCategoryBusinessRule:
		return RecoveryStrategyEscalate
	case ErrorCategoryNotFound:
		return RecoveryStrategyFallback
	case ErrorCategoryConflict:
		return RecoveryStrategyRetry
	default:
		return RecoveryStrategyEscalate
	}
}

// 错误统计

// ErrorStats 错误统计
type ErrorStats struct {
	TotalErrors    int64
	ErrorsByType   map[string]int64
	ErrorsByCode   map[string]int64
	LastError      error
	LastErrorTime  time.Time
	ErrorRate      float64
	RecoveryRate   float64
}

// NewErrorStats 创建错误统计
func NewErrorStats() *ErrorStats {
	return &ErrorStats{
		ErrorsByType: make(map[string]int64),
		ErrorsByCode: make(map[string]int64),
	}
}

// RecordError 记录错误
func (s *ErrorStats) RecordError(err error) {
	s.TotalErrors++
	s.LastError = err
	s.LastErrorTime = time.Now()
	
	// 按类型统计
	category := string(CategorizeError(err))
	s.ErrorsByType[category]++
	
	// 按错误码统计
	if npcErr, ok := err.(*NPCError); ok {
		s.ErrorsByCode[npcErr.Code]++
	} else {
		s.ErrorsByCode["unknown"]++
	}
}

// GetMostCommonError 获取最常见错误
func (s *ErrorStats) GetMostCommonError() (string, int64) {
	var maxType string
	var maxCount int64
	
	for errorType, count := range s.ErrorsByType {
		if count > maxCount {
			maxType = errorType
			maxCount = count
		}
	}
	
	return maxType, maxCount
}

// 错误处理器

// ErrorHandler 错误处理器接口
type ErrorHandler interface {
	Handle(err error) error
	CanHandle(err error) bool
	GetHandlerName() string
}

// DefaultErrorHandler 默认错误处理器
type DefaultErrorHandler struct {
	name     string
	handlers map[ErrorCategory]func(error) error
}

// NewDefaultErrorHandler 创建默认错误处理器
func NewDefaultErrorHandler(name string) *DefaultErrorHandler {
	return &DefaultErrorHandler{
		name:     name,
		handlers: make(map[ErrorCategory]func(error) error),
	}
}

// RegisterHandler 注册处理器
func (h *DefaultErrorHandler) RegisterHandler(category ErrorCategory, handler func(error) error) {
	h.handlers[category] = handler
}

// Handle 处理错误
func (h *DefaultErrorHandler) Handle(err error) error {
	category := CategorizeError(err)
	if handler, exists := h.handlers[category]; exists {
		return handler(err)
	}
	return err
}

// CanHandle 是否可以处理
func (h *DefaultErrorHandler) CanHandle(err error) bool {
	category := CategorizeError(err)
	_, exists := h.handlers[category]
	return exists
}

// GetHandlerName 获取处理器名称
func (h *DefaultErrorHandler) GetHandlerName() string {
	return h.name
}

// 错误上下文

// ErrorContext 错误上下文
type ErrorContext struct {
	OperationID   string
	UserID        string
	NPCID         string
	RequestID     string
	SessionID     string
	Timestamp     time.Time
	StackTrace    string
	AdditionalInfo map[string]interface{}
}

// NewErrorContext 创建错误上下文
func NewErrorContext() *ErrorContext {
	return &ErrorContext{
		Timestamp:      time.Now(),
		AdditionalInfo: make(map[string]interface{}),
	}
}

// WithOperation 设置操作ID
func (c *ErrorContext) WithOperation(operationID string) *ErrorContext {
	c.OperationID = operationID
	return c
}

// WithUser 设置用户ID
func (c *ErrorContext) WithUser(userID string) *ErrorContext {
	c.UserID = userID
	return c
}

// WithNPC 设置NPC ID
func (c *ErrorContext) WithNPC(npcID string) *ErrorContext {
	c.NPCID = npcID
	return c
}

// WithRequest 设置请求ID
func (c *ErrorContext) WithRequest(requestID string) *ErrorContext {
	c.RequestID = requestID
	return c
}

// WithSession 设置会话ID
func (c *ErrorContext) WithSession(sessionID string) *ErrorContext {
	c.SessionID = sessionID
	return c
}

// WithInfo 添加额外信息
func (c *ErrorContext) WithInfo(key string, value interface{}) *ErrorContext {
	c.AdditionalInfo[key] = value
	return c
}

// 辅助函数

// WrapError 包装错误
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// WrapErrorWithContext 带上下文包装错误
func WrapErrorWithContext(err error, message string, context *ErrorContext) error {
	if err == nil {
		return nil
	}
	
	npcErr := NewNPCError("WRAPPED_ERROR", message).WithCause(err)
	if context != nil {
		npcErr.WithContext("operation_id", context.OperationID)
		npcErr.WithContext("user_id", context.UserID)
		npcErr.WithContext("npc_id", context.NPCID)
		npcErr.WithContext("request_id", context.RequestID)
		npcErr.WithContext("session_id", context.SessionID)
	}
	
	return npcErr
}

// IsRetryableError 是否可重试错误
func IsRetryableError(err error) bool {
	strategy := GetRecoveryStrategy(err)
	return strategy == RecoveryStrategyRetry
}

// IsFatalError 是否致命错误
func IsFatalError(err error) bool {
	strategy := GetRecoveryStrategy(err)
	return strategy == RecoveryStrategyEscalate
}

// GetErrorSeverity 获取错误严重程度
func GetErrorSeverity(err error) string {
	category := CategorizeError(err)
	
	switch category {
	case ErrorCategoryValidation:
		return "low"
	case ErrorCategoryNotFound, ErrorCategoryRateLimit:
		return "medium"
	case ErrorCategoryBusinessRule, ErrorCategoryConflict:
		return "high"
	case ErrorCategoryPermission, ErrorCategorySystem, ErrorCategoryTimeout:
		return "critical"
	default:
		return "unknown"
	}
}