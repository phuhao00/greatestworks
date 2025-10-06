package minigame

import (
	"fmt"
	"time"
)

// MinigameError 小游戏错误接口
type MinigameError interface {
	error
	GetCode() string
	GetMessage() string
	GetSeverity() ErrorSeverity
	GetTimestamp() time.Time
	GetContext() map[string]interface{}
	SetContext(key string, value interface{})
	GetCause() error
	IsRetryable() bool
	GetRetryAfter() *time.Duration
}

// ErrorSeverity 错误严重程度
type ErrorSeverity int32

const (
	ErrorSeverityLow      ErrorSeverity = iota + 1 // 低严重程度
	ErrorSeverityMedium                            // 中等严重程度
	ErrorSeverityHigh                              // 高严重程度
	ErrorSeverityCritical                          // 严重程度
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

// BaseMinigameError 基础小游戏错误
type BaseMinigameError struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Severity   ErrorSeverity          `json:"severity"`
	Timestamp  time.Time              `json:"timestamp"`
	Context    map[string]interface{} `json:"context"`
	Cause      error                  `json:"cause,omitempty"`
	Retryable  bool                   `json:"retryable"`
	RetryAfter *time.Duration         `json:"retry_after,omitempty"`
}

// Error 实现error接口
func (e *BaseMinigameError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// GetCode 获取错误代码
func (e *BaseMinigameError) GetCode() string {
	return e.Code
}

// GetMessage 获取错误消息
func (e *BaseMinigameError) GetMessage() string {
	return e.Message
}

// GetSeverity 获取错误严重程度
func (e *BaseMinigameError) GetSeverity() ErrorSeverity {
	return e.Severity
}

// GetTimestamp 获取错误时间戳
func (e *BaseMinigameError) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetContext 获取错误上下文
func (e *BaseMinigameError) GetContext() map[string]interface{} {
	return e.Context
}

// SetContext 设置错误上下文
func (e *BaseMinigameError) SetContext(key string, value interface{}) {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
}

// GetCause 获取错误原因
func (e *BaseMinigameError) GetCause() error {
	return e.Cause
}

// IsRetryable 检查是否可重试
func (e *BaseMinigameError) IsRetryable() bool {
	return e.Retryable
}

// GetRetryAfter 获取重试延迟时间
func (e *BaseMinigameError) GetRetryAfter() *time.Duration {
	return e.RetryAfter
}

// 具体错误类型

// GameNotFoundError 游戏未找到错误
type GameNotFoundError struct {
	*BaseMinigameError
	GameID string `json:"game_id"`
}

// PlayerNotFoundError 玩家未找到错误
type PlayerNotFoundError struct {
	*BaseMinigameError
	PlayerID uint64 `json:"player_id"`
}

// SessionNotFoundError 会话未找到错误
type SessionNotFoundError struct {
	*BaseMinigameError
	SessionID string `json:"session_id"`
}

// ScoreNotFoundError 分数未找到错误
type ScoreNotFoundError struct {
	*BaseMinigameError
	ScoreID string `json:"score_id"`
}

// RewardNotFoundError 奖励未找到错误
type RewardNotFoundError struct {
	*BaseMinigameError
	RewardID string `json:"reward_id"`
}

// AchievementNotFoundError 成就未找到错误
type AchievementNotFoundError struct {
	*BaseMinigameError
	AchievementID string `json:"achievement_id"`
}

// InvalidGameTypeError 无效游戏类型错误
type InvalidGameTypeError struct {
	*BaseMinigameError
	GameType GameType `json:"game_type"`
}

// InvalidGameStatusError 无效游戏状态错误
type InvalidGameStatusError struct {
	*BaseMinigameError
	CurrentStatus GameStatus `json:"current_status"`
	TargetStatus  GameStatus `json:"target_status"`
}

// InvalidPlayerStatusError 无效玩家状态错误
type InvalidPlayerStatusError struct {
	*BaseMinigameError
	PlayerID      uint64       `json:"player_id"`
	CurrentStatus PlayerStatus `json:"current_status"`
	TargetStatus  PlayerStatus `json:"target_status"`
}

// GameNotJoinableError 游戏不可加入错误
type GameNotJoinableError struct {
	*BaseMinigameError
	GameID string     `json:"game_id"`
	Status GameStatus `json:"status"`
	Reason string     `json:"reason"`
}

// GameFullError 游戏已满错误
type GameFullError struct {
	*BaseMinigameError
	GameID         string `json:"game_id"`
	CurrentPlayers int32  `json:"current_players"`
	MaxPlayers     int32  `json:"max_players"`
}

// PlayerAlreadyInGameError 玩家已在游戏中错误
type PlayerAlreadyInGameError struct {
	*BaseMinigameError
	GameID   string `json:"game_id"`
	PlayerID uint64 `json:"player_id"`
}

// PlayerNotInGameError 玩家不在游戏中错误
type PlayerNotInGameError struct {
	*BaseMinigameError
	GameID   string `json:"game_id"`
	PlayerID uint64 `json:"player_id"`
}

// GameNotRunningError 游戏未运行错误
type GameNotRunningError struct {
	*BaseMinigameError
	GameID string     `json:"game_id"`
	Status GameStatus `json:"status"`
}

// InvalidScoreError 无效分数错误
type InvalidScoreError struct {
	*BaseMinigameError
	Score     int64     `json:"score"`
	ScoreType ScoreType `json:"score_type"`
	Reason    string    `json:"reason"`
}

// InvalidRewardError 无效奖励错误
type InvalidRewardError struct {
	*BaseMinigameError
	RewardType RewardType `json:"reward_type"`
	ItemID     string     `json:"item_id"`
	Quantity   int64      `json:"quantity"`
	Reason     string     `json:"reason"`
}

// RewardAlreadyClaimedError 奖励已领取错误
type RewardAlreadyClaimedError struct {
	*BaseMinigameError
	RewardID  string    `json:"reward_id"`
	ClaimedAt time.Time `json:"claimed_at"`
}

// RewardExpiredError 奖励已过期错误
type RewardExpiredError struct {
	*BaseMinigameError
	RewardID  string    `json:"reward_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// InvalidConfigError 无效配置错误
type InvalidConfigError struct {
	*BaseMinigameError
	ConfigField string      `json:"config_field"`
	ConfigValue interface{} `json:"config_value"`
	Reason      string      `json:"reason"`
}

// InvalidSessionError 无效会话错误
type InvalidSessionError struct {
	*BaseMinigameError
	SessionID string `json:"session_id"`
	Reason    string `json:"reason"`
}

// PermissionDeniedError 权限拒绝错误
type PermissionDeniedError struct {
	*BaseMinigameError
	UserID     uint64 `json:"user_id"`
	Operation  string `json:"operation"`
	ResourceID string `json:"resource_id"`
}

// RateLimitExceededError 速率限制超出错误
type RateLimitExceededError struct {
	*BaseMinigameError
	UserID     uint64        `json:"user_id"`
	Operation  string        `json:"operation"`
	Limit      int64         `json:"limit"`
	WindowSize time.Duration `json:"window_size"`
	ResetTime  time.Time     `json:"reset_time"`
}

// ConcurrencyLimitError 并发限制错误
type ConcurrencyLimitError struct {
	*BaseMinigameError
	CurrentCount int64  `json:"current_count"`
	MaxCount     int64  `json:"max_count"`
	ResourceType string `json:"resource_type"`
}

// ValidationError 验证错误
type ValidationError struct {
	*BaseMinigameError
	Field  string      `json:"field"`
	Value  interface{} `json:"value"`
	Rule   string      `json:"rule"`
	Reason string      `json:"reason"`
}

// RepositoryError 仓储错误
type RepositoryError struct {
	*BaseMinigameError
	Operation string `json:"operation"`
	Entity    string `json:"entity"`
	EntityID  string `json:"entity_id"`
}

// NetworkError 网络错误
type NetworkError struct {
	*BaseMinigameError
	Endpoint   string        `json:"endpoint"`
	Method     string        `json:"method"`
	StatusCode int           `json:"status_code,omitempty"`
	Timeout    time.Duration `json:"timeout,omitempty"`
}

// TimeoutError 超时错误
type TimeoutError struct {
	*BaseMinigameError
	Operation string        `json:"operation"`
	Timeout   time.Duration `json:"timeout"`
	Elapsed   time.Duration `json:"elapsed"`
}

// ResourceExhaustedError 资源耗尽错误
type ResourceExhaustedError struct {
	*BaseMinigameError
	ResourceType string `json:"resource_type"`
	Limit        int64  `json:"limit"`
	Used         int64  `json:"used"`
}

// InternalError 内部错误
type InternalError struct {
	*BaseMinigameError
	Component string `json:"component"`
	Function  string `json:"function"`
}

// 错误代码常量

const (
	// 通用错误代码
	ErrorCodeUnknown           = "MINIGAME_UNKNOWN"
	ErrorCodeInternalError     = "MINIGAME_INTERNAL_ERROR"
	ErrorCodeValidationError   = "MINIGAME_VALIDATION_ERROR"
	ErrorCodePermissionDenied  = "MINIGAME_PERMISSION_DENIED"
	ErrorCodeRateLimitExceeded = "MINIGAME_RATE_LIMIT_EXCEEDED"
	ErrorCodeConcurrencyLimit  = "MINIGAME_CONCURRENCY_LIMIT"
	ErrorCodeResourceExhausted = "MINIGAME_RESOURCE_EXHAUSTED"
	ErrorCodeTimeout           = "MINIGAME_TIMEOUT"
	ErrorCodeNetworkError      = "MINIGAME_NETWORK_ERROR"

	// 游戏相关错误代码
	ErrorCodeGameNotFound      = "MINIGAME_GAME_NOT_FOUND"
	ErrorCodeInvalidGameType   = "MINIGAME_INVALID_GAME_TYPE"
	ErrorCodeInvalidGameStatus = "MINIGAME_INVALID_GAME_STATUS"
	ErrorCodeGameNotJoinable   = "MINIGAME_GAME_NOT_JOINABLE"
	ErrorCodeGameFull          = "MINIGAME_GAME_FULL"
	ErrorCodeGameNotRunning    = "MINIGAME_GAME_NOT_RUNNING"
	ErrorCodeInvalidOperation  = "MINIGAME_INVALID_OPERATION"
	ErrorCodeInvalidConfig     = "MINIGAME_INVALID_CONFIG"

	// 玩家相关错误代码
	ErrorCodePlayerNotFound      = "MINIGAME_PLAYER_NOT_FOUND"
	ErrorCodeInvalidPlayer       = "MINIGAME_INVALID_PLAYER"
	ErrorCodeInvalidPlayerStatus = "MINIGAME_INVALID_PLAYER_STATUS"
	ErrorCodePlayerAlreadyInGame = "MINIGAME_PLAYER_ALREADY_IN_GAME"
	ErrorCodePlayerNotInGame     = "MINIGAME_PLAYER_NOT_IN_GAME"

	// 会话相关错误代码
	ErrorCodeSessionNotFound = "MINIGAME_SESSION_NOT_FOUND"
	ErrorCodeInvalidSession  = "MINIGAME_INVALID_SESSION"

	// 分数相关错误代码
	ErrorCodeScoreNotFound = "MINIGAME_SCORE_NOT_FOUND"
	ErrorCodeInvalidScore  = "MINIGAME_INVALID_SCORE"

	// 奖励相关错误代码
	ErrorCodeRewardNotFound       = "MINIGAME_REWARD_NOT_FOUND"
	ErrorCodeInvalidReward        = "MINIGAME_INVALID_REWARD"
	ErrorCodeRewardAlreadyClaimed = "MINIGAME_REWARD_ALREADY_CLAIMED"
	ErrorCodeRewardExpired        = "MINIGAME_REWARD_EXPIRED"

	// 成就相关错误代码
	ErrorCodeAchievementNotFound = "MINIGAME_ACHIEVEMENT_NOT_FOUND"
	ErrorCodeInvalidAchievement  = "MINIGAME_INVALID_ACHIEVEMENT"

	// 仓储相关错误代码
	ErrorCodeRepositoryError = "MINIGAME_REPOSITORY_ERROR"
	ErrorCodeDatabaseError   = "MINIGAME_DATABASE_ERROR"
	ErrorCodeCacheError      = "MINIGAME_CACHE_ERROR"
)

// 错误工厂函数

// NewMinigameError 创建基础小游戏错误
func NewMinigameError(code, message string, cause error) *BaseMinigameError {
	return &BaseMinigameError{
		Code:      code,
		Message:   message,
		Severity:  ErrorSeverityMedium,
		Timestamp: time.Now(),
		Context:   make(map[string]interface{}),
		Cause:     cause,
		Retryable: false,
	}
}

// NewMinigameInvalidStateError 创建无效状态错误
func NewMinigameInvalidStateError(gameID string, currentStatus, expectedStatus GameStatus, operation string) *BaseMinigameError {
	return &BaseMinigameError{
		Code:      ErrorCodeInvalidState,
		Message:   fmt.Sprintf("Invalid state for operation %s: current=%s, expected=%s", operation, currentStatus, expectedStatus),
		Severity:  ErrorSeverityMedium,
		Timestamp: time.Now(),
		Context:   map[string]interface{}{"game_id": gameID, "current_status": currentStatus, "expected_status": expectedStatus, "operation": operation},
		Retryable: false,
	}
}

// NewMinigameInsufficientPlayersError 创建玩家不足错误
func NewMinigameInsufficientPlayersError(gameID string, currentPlayers, minPlayers int32) *BaseMinigameError {
	return &BaseMinigameError{
		Code:      ErrorCodeInsufficientPlayers,
		Message:   fmt.Sprintf("Insufficient players: current=%d, minimum=%d", currentPlayers, minPlayers),
		Severity:  ErrorSeverityMedium,
		Timestamp: time.Now(),
		Context:   map[string]interface{}{"game_id": gameID, "current_players": currentPlayers, "min_players": minPlayers},
		Retryable: false,
	}
}

// NewPlayerAlreadyInGameError 创建玩家已在游戏中错误
func NewPlayerAlreadyInGameError(playerID uint64, gameID string) *BaseMinigameError {
	return &BaseMinigameError{
		Code:      ErrorCodePlayerAlreadyInGame,
		Message:   fmt.Sprintf("Player %d is already in game %s", playerID, gameID),
		Severity:  ErrorSeverityMedium,
		Timestamp: time.Now(),
		Context:   map[string]interface{}{"player_id": playerID, "game_id": gameID},
		Retryable: false,
	}
}

// NewGameNotFoundError 创建游戏未找到错误
func NewGameNotFoundError(gameID string) *GameNotFoundError {
	return &GameNotFoundError{
		BaseMinigameError: &BaseMinigameError{
			Code:      ErrorCodeGameNotFound,
			Message:   fmt.Sprintf("Game not found: %s", gameID),
			Severity:  ErrorSeverityMedium,
			Timestamp: time.Now(),
			Context:   make(map[string]interface{}),
			Retryable: false,
		},
		GameID: gameID,
	}
}

// NewPlayerNotFoundError 创建玩家未找到错误
func NewPlayerNotFoundError(playerID uint64) *PlayerNotFoundError {
	return &PlayerNotFoundError{
		BaseMinigameError: &BaseMinigameError{
			Code:      ErrorCodePlayerNotFound,
			Message:   fmt.Sprintf("Player not found: %d", playerID),
			Severity:  ErrorSeverityMedium,
			Timestamp: time.Now(),
			Context:   make(map[string]interface{}),
			Retryable: false,
		},
		PlayerID: playerID,
	}
}

// NewSessionNotFoundError 创建会话未找到错误
func NewSessionNotFoundError(sessionID string) *SessionNotFoundError {
	return &SessionNotFoundError{
		BaseMinigameError: &BaseMinigameError{
			Code:      ErrorCodeSessionNotFound,
			Message:   fmt.Sprintf("Session not found: %s", sessionID),
			Severity:  ErrorSeverityMedium,
			Timestamp: time.Now(),
			Context:   make(map[string]interface{}),
			Retryable: false,
		},
		SessionID: sessionID,
	}
}

// NewInvalidGameTypeError 创建无效游戏类型错误
func NewInvalidGameTypeError(gameType GameType) *InvalidGameTypeError {
	return &InvalidGameTypeError{
		BaseMinigameError: &BaseMinigameError{
			Code:      ErrorCodeInvalidGameType,
			Message:   fmt.Sprintf("Invalid game type: %v", gameType),
			Severity:  ErrorSeverityMedium,
			Timestamp: time.Now(),
			Context:   make(map[string]interface{}),
			Retryable: false,
		},
		GameType: gameType,
	}
}

// NewGameNotJoinableError 创建游戏不可加入错误
func NewGameNotJoinableError(gameID string, status GameStatus, reason string) *GameNotJoinableError {
	return &GameNotJoinableError{
		BaseMinigameError: &BaseMinigameError{
			Code:      ErrorCodeGameNotJoinable,
			Message:   fmt.Sprintf("Game %s is not joinable: %s", gameID, reason),
			Severity:  ErrorSeverityMedium,
			Timestamp: time.Now(),
			Context:   make(map[string]interface{}),
			Retryable: false,
		},
		GameID: gameID,
		Status: status,
		Reason: reason,
	}
}

// NewGameFullError 创建游戏已满错误
func NewGameFullError(gameID string, currentPlayers, maxPlayers int32) *GameFullError {
	return &GameFullError{
		BaseMinigameError: &BaseMinigameError{
			Code:      ErrorCodeGameFull,
			Message:   fmt.Sprintf("Game %s is full: %d/%d players", gameID, currentPlayers, maxPlayers),
			Severity:  ErrorSeverityMedium,
			Timestamp: time.Now(),
			Context:   make(map[string]interface{}),
			Retryable: true,
		},
		GameID:         gameID,
		CurrentPlayers: currentPlayers,
		MaxPlayers:     maxPlayers,
	}
}

// NewPlayerAlreadyInGameError 创建玩家已在游戏中错误 (duplicate removed - see line 428)

// NewPlayerNotInGameError 创建玩家不在游戏中错误
func NewPlayerNotInGameError(gameID string, playerID uint64) *PlayerNotInGameError {
	return &PlayerNotInGameError{
		BaseMinigameError: &BaseMinigameError{
			Code:      ErrorCodePlayerNotInGame,
			Message:   fmt.Sprintf("Player %d is not in game %s", playerID, gameID),
			Severity:  ErrorSeverityMedium,
			Timestamp: time.Now(),
			Context:   make(map[string]interface{}),
			Retryable: false,
		},
		GameID:   gameID,
		PlayerID: playerID,
	}
}

// NewInvalidScoreError 创建无效分数错误
func NewInvalidScoreError(score int64, scoreType ScoreType, reason string) *InvalidScoreError {
	return &InvalidScoreError{
		BaseMinigameError: &BaseMinigameError{
			Code:      ErrorCodeInvalidScore,
			Message:   fmt.Sprintf("Invalid score %d for type %v: %s", score, scoreType, reason),
			Severity:  ErrorSeverityMedium,
			Timestamp: time.Now(),
			Context:   make(map[string]interface{}),
			Retryable: false,
		},
		Score:     score,
		ScoreType: scoreType,
		Reason:    reason,
	}
}

// NewRewardAlreadyClaimedError 创建奖励已领取错误
func NewRewardAlreadyClaimedError(rewardID string, claimedAt time.Time) *RewardAlreadyClaimedError {
	return &RewardAlreadyClaimedError{
		BaseMinigameError: &BaseMinigameError{
			Code:      ErrorCodeRewardAlreadyClaimed,
			Message:   fmt.Sprintf("Reward %s has already been claimed at %v", rewardID, claimedAt),
			Severity:  ErrorSeverityMedium,
			Timestamp: time.Now(),
			Context:   make(map[string]interface{}),
			Retryable: false,
		},
		RewardID:  rewardID,
		ClaimedAt: claimedAt,
	}
}

// NewRewardExpiredError 创建奖励已过期错误
func NewRewardExpiredError(rewardID string, expiresAt time.Time) *RewardExpiredError {
	return &RewardExpiredError{
		BaseMinigameError: &BaseMinigameError{
			Code:      ErrorCodeRewardExpired,
			Message:   fmt.Sprintf("Reward %s has expired at %v", rewardID, expiresAt),
			Severity:  ErrorSeverityMedium,
			Timestamp: time.Now(),
			Context:   make(map[string]interface{}),
			Retryable: false,
		},
		RewardID:  rewardID,
		ExpiresAt: expiresAt,
	}
}

// NewPermissionDeniedError 创建权限拒绝错误
func NewPermissionDeniedError(userID uint64, operation, resourceID string) *PermissionDeniedError {
	return &PermissionDeniedError{
		BaseMinigameError: &BaseMinigameError{
			Code:      ErrorCodePermissionDenied,
			Message:   fmt.Sprintf("Permission denied for user %d to perform %s on %s", userID, operation, resourceID),
			Severity:  ErrorSeverityHigh,
			Timestamp: time.Now(),
			Context:   make(map[string]interface{}),
			Retryable: false,
		},
		UserID:     userID,
		Operation:  operation,
		ResourceID: resourceID,
	}
}

// NewRateLimitExceededError 创建速率限制超出错误
func NewRateLimitExceededError(userID uint64, operation string, limit int64, windowSize time.Duration, resetTime time.Time) *RateLimitExceededError {
	retryAfter := time.Until(resetTime)
	return &RateLimitExceededError{
		BaseMinigameError: &BaseMinigameError{
			Code:       ErrorCodeRateLimitExceeded,
			Message:    fmt.Sprintf("Rate limit exceeded for user %d on %s: %d requests per %v", userID, operation, limit, windowSize),
			Severity:   ErrorSeverityMedium,
			Timestamp:  time.Now(),
			Context:    make(map[string]interface{}),
			Retryable:  true,
			RetryAfter: &retryAfter,
		},
		UserID:     userID,
		Operation:  operation,
		Limit:      limit,
		WindowSize: windowSize,
		ResetTime:  resetTime,
	}
}

// NewTimeoutError 创建超时错误
func NewTimeoutError(operation string, timeout, elapsed time.Duration) *TimeoutError {
	return &TimeoutError{
		BaseMinigameError: &BaseMinigameError{
			Code:      ErrorCodeTimeout,
			Message:   fmt.Sprintf("Operation %s timed out after %v (elapsed: %v)", operation, timeout, elapsed),
			Severity:  ErrorSeverityHigh,
			Timestamp: time.Now(),
			Context:   make(map[string]interface{}),
			Retryable: true,
		},
		Operation: operation,
		Timeout:   timeout,
		Elapsed:   elapsed,
	}
}

// NewRepositoryError 创建仓储错误
func NewRepositoryError(operation, entity, entityID string, cause error) *RepositoryError {
	return &RepositoryError{
		BaseMinigameError: &BaseMinigameError{
			Code:      ErrorCodeRepositoryError,
			Message:   fmt.Sprintf("Repository error during %s on %s %s", operation, entity, entityID),
			Severity:  ErrorSeverityHigh,
			Timestamp: time.Now(),
			Context:   make(map[string]interface{}),
			Cause:     cause,
			Retryable: true,
		},
		Operation: operation,
		Entity:    entity,
		EntityID:  entityID,
	}
}

// NewInternalError 创建内部错误
func NewInternalError(component, function string, cause error) *InternalError {
	return &InternalError{
		BaseMinigameError: &BaseMinigameError{
			Code:      ErrorCodeInternalError,
			Message:   fmt.Sprintf("Internal error in %s.%s", component, function),
			Severity:  ErrorSeverityCritical,
			Timestamp: time.Now(),
			Context:   make(map[string]interface{}),
			Cause:     cause,
			Retryable: false,
		},
		Component: component,
		Function:  function,
	}
}

// 错误分类函数

// IsNotFoundError 检查是否为未找到错误
func IsNotFoundError(err error) bool {
	switch err.(type) {
	case *GameNotFoundError, *PlayerNotFoundError, *SessionNotFoundError,
		*ScoreNotFoundError, *RewardNotFoundError, *AchievementNotFoundError:
		return true
	default:
		return false
	}
}

// IsValidationError 检查是否为验证错误
func IsValidationError(err error) bool {
	switch err.(type) {
	case *InvalidGameTypeError, *InvalidGameStatusError, *InvalidPlayerStatusError,
		*InvalidScoreError, *InvalidRewardError, *InvalidConfigError, *InvalidSessionError, *ValidationError:
		return true
	default:
		return false
	}
}

// IsPermissionError 检查是否为权限错误
func IsPermissionError(err error) bool {
	switch err.(type) {
	case *PermissionDeniedError:
		return true
	default:
		return false
	}
}

// IsRateLimitError 检查是否为速率限制错误
func IsRateLimitError(err error) bool {
	switch err.(type) {
	case *RateLimitExceededError:
		return true
	default:
		return false
	}
}

// IsRetryableError 检查是否为可重试错误
func IsRetryableError(err error) bool {
	if minigameErr, ok := err.(MinigameError); ok {
		return minigameErr.IsRetryable()
	}
	return false
}

// IsTemporaryError 检查是否为临时错误
func IsTemporaryError(err error) bool {
	switch err.(type) {
	case *NetworkError, *TimeoutError, *RateLimitExceededError, *ConcurrencyLimitError, *ResourceExhaustedError:
		return true
	default:
		return false
	}
}

// IsCriticalError 检查是否为严重错误
func IsCriticalError(err error) bool {
	if minigameErr, ok := err.(MinigameError); ok {
		return minigameErr.GetSeverity() == ErrorSeverityCritical
	}
	return false
}

// 错误恢复策略

// ErrorRecoveryStrategy 错误恢复策略
type ErrorRecoveryStrategy int32

const (
	RecoveryStrategyNone                ErrorRecoveryStrategy = iota + 1 // 无恢复策略
	RecoveryStrategyRetry                                                // 重试
	RecoveryStrategyFallback                                             // 降级
	RecoveryStrategyCircuitBreaker                                       // 熔断
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
	switch err.(type) {
	case *NetworkError, *TimeoutError:
		return RecoveryStrategyRetry
	case *RateLimitExceededError:
		return RecoveryStrategyRetry
	case *ConcurrencyLimitError, *ResourceExhaustedError:
		return RecoveryStrategyCircuitBreaker
	case *GameFullError:
		return RecoveryStrategyFallback
	case *RepositoryError:
		return RecoveryStrategyRetry
	case *InternalError:
		return RecoveryStrategyGracefulDegradation
	default:
		return RecoveryStrategyNone
	}
}

// GetRetryDelay 获取重试延迟时间
func GetRetryDelay(err error, attempt int) time.Duration {
	baseDelay := time.Second
	maxDelay := 30 * time.Second

	// 指数退避
	delay := time.Duration(attempt) * baseDelay
	if delay > maxDelay {
		delay = maxDelay
	}

	// 特殊错误类型的延迟调整
	switch e := err.(type) {
	case *RateLimitExceededError:
		if e.GetRetryAfter() != nil {
			return *e.GetRetryAfter()
		}
	case *TimeoutError:
		// 超时错误使用更长的延迟
		delay = delay * 2
	case *NetworkError:
		// 网络错误使用较短的延迟
		delay = delay / 2
	}

	return delay
}

// GetMaxRetryAttempts 获取最大重试次数
func GetMaxRetryAttempts(err error) int {
	switch err.(type) {
	case *NetworkError:
		return 3
	case *TimeoutError:
		return 2
	case *RateLimitExceededError:
		return 5
	case *RepositoryError:
		return 3
	case *ResourceExhaustedError:
		return 1
	default:
		return 0
	}
}

// 错误聚合和统计

// ErrorStatistics 错误统计
type ErrorStatistics struct {
	TotalErrors      int64                   `json:"total_errors"`
	ErrorsByCode     map[string]int64        `json:"errors_by_code"`
	ErrorsBySeverity map[ErrorSeverity]int64 `json:"errors_by_severity"`
	ErrorsByType     map[string]int64        `json:"errors_by_type"`
	RetryableErrors  int64                   `json:"retryable_errors"`
	CriticalErrors   int64                   `json:"critical_errors"`
	LastError        *time.Time              `json:"last_error,omitempty"`
	ErrorRate        float64                 `json:"error_rate"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
}

// NewErrorStatistics 创建错误统计
func NewErrorStatistics() *ErrorStatistics {
	now := time.Now()
	return &ErrorStatistics{
		TotalErrors:      0,
		ErrorsByCode:     make(map[string]int64),
		ErrorsBySeverity: make(map[ErrorSeverity]int64),
		ErrorsByType:     make(map[string]int64),
		RetryableErrors:  0,
		CriticalErrors:   0,
		ErrorRate:        0.0,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// RecordError 记录错误
func (es *ErrorStatistics) RecordError(err error) {
	es.TotalErrors++
	now := time.Now()
	es.LastError = &now
	es.UpdatedAt = now

	if minigameErr, ok := err.(MinigameError); ok {
		// 按错误代码统计
		code := minigameErr.GetCode()
		es.ErrorsByCode[code]++

		// 按严重程度统计
		severity := minigameErr.GetSeverity()
		es.ErrorsBySeverity[severity]++

		// 统计可重试错误
		if minigameErr.IsRetryable() {
			es.RetryableErrors++
		}

		// 统计严重错误
		if severity == ErrorSeverityCritical {
			es.CriticalErrors++
		}
	}

	// 按错误类型统计
	errorType := fmt.Sprintf("%T", err)
	es.ErrorsByType[errorType]++
}

// CalculateErrorRate 计算错误率
func (es *ErrorStatistics) CalculateErrorRate(totalRequests int64) {
	if totalRequests > 0 {
		es.ErrorRate = float64(es.TotalErrors) / float64(totalRequests) * 100
	} else {
		es.ErrorRate = 0.0
	}
	es.UpdatedAt = time.Now()
}

// Reset 重置统计
func (es *ErrorStatistics) Reset() {
	es.TotalErrors = 0
	es.ErrorsByCode = make(map[string]int64)
	es.ErrorsBySeverity = make(map[ErrorSeverity]int64)
	es.ErrorsByType = make(map[string]int64)
	es.RetryableErrors = 0
	es.CriticalErrors = 0
	es.LastError = nil
	es.ErrorRate = 0.0
	es.UpdatedAt = time.Now()
}

// 辅助函数

// WrapError 包装错误
func WrapError(err error, code, message string) MinigameError {
	if err == nil {
		return nil
	}

	if minigameErr, ok := err.(MinigameError); ok {
		return minigameErr
	}

	return NewMinigameError(code, message, err)
}

// UnwrapError 解包错误
func UnwrapError(err error) error {
	if minigameErr, ok := err.(MinigameError); ok {
		return minigameErr.GetCause()
	}
	return err
}

// FormatError 格式化错误信息
func FormatError(err error) string {
	if minigameErr, ok := err.(MinigameError); ok {
		return fmt.Sprintf("[%s] %s (severity: %s, timestamp: %v)",
			minigameErr.GetCode(),
			minigameErr.GetMessage(),
			minigameErr.GetSeverity(),
			minigameErr.GetTimestamp().Format(time.RFC3339))
	}
	return err.Error()
}

// LogError 记录错误日志
func LogError(err error, context map[string]interface{}) {
	if minigameErr, ok := err.(MinigameError); ok {
		// 合并上下文
		for k, v := range context {
			minigameErr.SetContext(k, v)
		}

		// 根据严重程度记录日志
		switch minigameErr.GetSeverity() {
		case ErrorSeverityCritical:
			// 记录严重错误日志
			fmt.Printf("CRITICAL ERROR: %s\n", FormatError(err))
		case ErrorSeverityHigh:
			// 记录高级错误日志
			fmt.Printf("HIGH ERROR: %s\n", FormatError(err))
		case ErrorSeverityMedium:
			// 记录中级错误日志
			fmt.Printf("MEDIUM ERROR: %s\n", FormatError(err))
		case ErrorSeverityLow:
			// 记录低级错误日志
			fmt.Printf("LOW ERROR: %s\n", FormatError(err))
		}
	} else {
		// 记录普通错误日志
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}

// NewMinigameValidationError 创建小游戏验证错误
func NewMinigameValidationError(message string) error {
	return &MinigameError{
		Code:      "VALIDATION_ERROR",
		Message:   message,
		Severity:  ErrorSeverityMedium,
		Retryable: false,
		Details:   make(map[string]interface{}),
	}
}
