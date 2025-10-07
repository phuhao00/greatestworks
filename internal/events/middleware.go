package events

import (
	"context"
	"fmt"
	"log"
	"time"
)

// LoggingMiddleware 日志中间件
type LoggingMiddleware struct {
	logger *log.Logger
}

// NewLoggingMiddleware 创建日志中间件
func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: log.New(log.Writer(), "[EventLogging] ", log.LstdFlags),
	}
}

// Process 处理事件
func (lm *LoggingMiddleware) Process(ctx context.Context, event Event, next func(context.Context, Event) error) error {
	start := time.Now()
	lm.logger.Printf("Processing event: %s (type: %s, player: %s)",
		event.GetID(), event.GetType(), event.GetPlayerID())

	err := next(ctx, event)

	duration := time.Since(start)
	if err != nil {
		lm.logger.Printf("Event processing failed: %s, duration: %v, error: %v",
			event.GetID(), duration, err)
	} else {
		lm.logger.Printf("Event processing completed: %s, duration: %v",
			event.GetID(), duration)
	}

	return err
}

// ValidationMiddleware 验证中间件
type ValidationMiddleware struct {
	logger *log.Logger
}

// NewValidationMiddleware 创建验证中间件
func NewValidationMiddleware() *ValidationMiddleware {
	return &ValidationMiddleware{
		logger: log.New(log.Writer(), "[EventValidation] ", log.LstdFlags),
	}
}

// Process 处理事件
func (vm *ValidationMiddleware) Process(ctx context.Context, event Event, next func(context.Context, Event) error) error {
	// 验证事件基本信息
	if err := vm.validateEvent(event); err != nil {
		vm.logger.Printf("Event validation failed: %s, error: %v", event.GetID(), err)
		return fmt.Errorf("event validation failed: %w", err)
	}

	return next(ctx, event)
}

// validateEvent 验证事件
func (vm *ValidationMiddleware) validateEvent(event Event) error {
	if event.GetID() == "" {
		return fmt.Errorf("event ID is required")
	}

	if event.GetType() == "" {
		return fmt.Errorf("event type is required")
	}

	if event.GetTimestamp().IsZero() {
		return fmt.Errorf("event timestamp is required")
	}

	// 检查时间戳是否合理（不能是未来时间，不能太久以前）
	now := time.Now()
	if event.GetTimestamp().After(now.Add(5 * time.Minute)) {
		return fmt.Errorf("event timestamp is in the future")
	}

	if event.GetTimestamp().Before(now.Add(-24 * time.Hour)) {
		return fmt.Errorf("event timestamp is too old")
	}

	return nil
}

// RateLimitMiddleware 限流中间件
type RateLimitMiddleware struct {
	limiter map[string]*TokenBucket
	logger  *log.Logger
}

// TokenBucket 令牌桶
type TokenBucket struct {
	capacity   int
	tokens     int
	refillRate int // 每秒补充的令牌数
	lastRefill time.Time
}

// NewRateLimitMiddleware 创建限流中间件
func NewRateLimitMiddleware() *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiter: make(map[string]*TokenBucket),
		logger:  log.New(log.Writer(), "[EventRateLimit] ", log.LstdFlags),
	}
}

// Process 处理事件
func (rlm *RateLimitMiddleware) Process(ctx context.Context, event Event, next func(context.Context, Event) error) error {
	// 基于玩家ID进行限流
	playerID := event.GetPlayerID()
	if playerID == "" {
		// 系统事件不限流
		return next(ctx, event)
	}

	if !rlm.allowRequest(playerID) {
		rlm.logger.Printf("Rate limit exceeded for player: %s, event: %s", playerID, event.GetID())
		return fmt.Errorf("rate limit exceeded for player: %s", playerID)
	}

	return next(ctx, event)
}

// allowRequest 检查是否允许请求
func (rlm *RateLimitMiddleware) allowRequest(playerID string) bool {
	bucket, exists := rlm.limiter[playerID]
	if !exists {
		// 为新玩家创建令牌桶
		bucket = &TokenBucket{
			capacity:   10, // 容量10个令牌
			tokens:     10, // 初始10个令牌
			refillRate: 5,  // 每秒补充5个令牌
			lastRefill: time.Now(),
		}
		rlm.limiter[playerID] = bucket
	}

	// 补充令牌
	now := time.Now()
	elapsed := now.Sub(bucket.lastRefill).Seconds()
	if elapsed > 0 {
		newTokens := int(elapsed * float64(bucket.refillRate))
		bucket.tokens += newTokens
		if bucket.tokens > bucket.capacity {
			bucket.tokens = bucket.capacity
		}
		bucket.lastRefill = now
	}

	// 检查是否有可用令牌
	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

// AuthenticationMiddleware 认证中间件
type AuthenticationMiddleware struct {
	logger *log.Logger
}

// NewAuthenticationMiddleware 创建认证中间件
func NewAuthenticationMiddleware() *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		logger: log.New(log.Writer(), "[EventAuth] ", log.LstdFlags),
	}
}

// Process 处理事件
func (am *AuthenticationMiddleware) Process(ctx context.Context, event Event, next func(context.Context, Event) error) error {
	// 检查事件是否需要认证
	if am.requiresAuthentication(EventType(event.GetType())) {
		if err := am.authenticateEvent(ctx, event); err != nil {
			am.logger.Printf("Event authentication failed: %s, error: %v", event.GetID(), err)
			return fmt.Errorf("event authentication failed: %w", err)
		}
	}

	return next(ctx, event)
}

// requiresAuthentication 检查事件类型是否需要认证
func (am *AuthenticationMiddleware) requiresAuthentication(eventType EventType) bool {
	// 系统事件不需要认证
	if eventType == EventTypeSystemStart || eventType == EventTypeSystemStop || eventType == EventTypeSystemHealth {
		return false
	}

	// 其他事件需要认证
	return true
}

// authenticateEvent 认证事件
func (am *AuthenticationMiddleware) authenticateEvent(ctx context.Context, event Event) error {
	playerID := event.GetPlayerID()
	if playerID == "" {
		return fmt.Errorf("player ID is required for authenticated events")
	}

	// 这里可以添加更复杂的认证逻辑
	// 例如：验证JWT token、检查玩家状态等

	return nil
}

// MetricsMiddleware 指标中间件
type MetricsMiddleware struct {
	metrics *EventMetrics
	logger  *log.Logger
}

// NewMetricsMiddleware 创建指标中间件
func NewMetricsMiddleware(metrics *EventMetrics) *MetricsMiddleware {
	return &MetricsMiddleware{
		metrics: metrics,
		logger:  log.New(log.Writer(), "[EventMetrics] ", log.LstdFlags),
	}
}

// Process 处理事件
func (mm *MetricsMiddleware) Process(ctx context.Context, event Event, next func(context.Context, Event) error) error {
	start := time.Now()

	err := next(ctx, event)

	duration := time.Since(start)
	mm.metrics.RecordProcessingTime(EventType(event.GetType()), duration)

	if err != nil {
		mm.metrics.IncrementErrorCount(EventType(event.GetType()))
	} else {
		mm.metrics.IncrementSuccessCount(EventType(event.GetType()))
	}

	return err
}

// CircuitBreakerMiddleware 熔断器中间件
type CircuitBreakerMiddleware struct {
	breakers map[EventType]*CircuitBreaker
	logger   *log.Logger
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	failureCount     int
	failureThreshold int
	timeout          time.Duration
	lastFailureTime  time.Time
	state            CircuitBreakerState
}

// CircuitBreakerState 熔断器状态
type CircuitBreakerState int

const (
	Closed CircuitBreakerState = iota
	Open
	HalfOpen
)

// NewCircuitBreakerMiddleware 创建熔断器中间件
func NewCircuitBreakerMiddleware() *CircuitBreakerMiddleware {
	return &CircuitBreakerMiddleware{
		breakers: make(map[EventType]*CircuitBreaker),
		logger:   log.New(log.Writer(), "[EventCircuitBreaker] ", log.LstdFlags),
	}
}

// Process 处理事件
func (cbm *CircuitBreakerMiddleware) Process(ctx context.Context, event Event, next func(context.Context, Event) error) error {
	breaker := cbm.getBreaker(EventType(event.GetType()))

	if !breaker.allowRequest() {
		cbm.logger.Printf("Circuit breaker is open for event type: %s", event.GetType())
		return fmt.Errorf("circuit breaker is open for event type: %s", event.GetType())
	}

	err := next(ctx, event)

	if err != nil {
		breaker.recordFailure()
	} else {
		breaker.recordSuccess()
	}

	return err
}

// getBreaker 获取熔断器
func (cbm *CircuitBreakerMiddleware) getBreaker(eventType EventType) *CircuitBreaker {
	breaker, exists := cbm.breakers[eventType]
	if !exists {
		breaker = &CircuitBreaker{
			failureThreshold: 5,                // 失败阈值
			timeout:          30 * time.Second, // 超时时间
			state:            Closed,
		}
		cbm.breakers[eventType] = breaker
	}
	return breaker
}

// allowRequest 检查是否允许请求
func (cb *CircuitBreaker) allowRequest() bool {
	now := time.Now()

	switch cb.state {
	case Closed:
		return true
	case Open:
		if now.Sub(cb.lastFailureTime) > cb.timeout {
			cb.state = HalfOpen
			return true
		}
		return false
	case HalfOpen:
		return true
	default:
		return false
	}
}

// recordFailure 记录失败
func (cb *CircuitBreaker) recordFailure() {
	cb.failureCount++
	cb.lastFailureTime = time.Now()

	if cb.state == HalfOpen {
		cb.state = Open
	} else if cb.failureCount >= cb.failureThreshold {
		cb.state = Open
	}
}

// recordSuccess 记录成功
func (cb *CircuitBreaker) recordSuccess() {
	if cb.state == HalfOpen {
		cb.state = Closed
		cb.failureCount = 0
	}
}
