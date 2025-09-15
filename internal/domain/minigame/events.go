package minigame

import (
	"context"
	"time"
)

// MinigameEvent 小游戏事件接口
type MinigameEvent interface {
	GetEventID() string
	GetEventType() string
	GetGameID() string
	GetPlayerID() *uint64
	GetTimestamp() time.Time
	GetMetadata() map[string]interface{}
	SetMetadata(key string, value interface{})
}

// BaseMinigameEvent 基础小游戏事件
type BaseMinigameEvent struct {
	EventID   string                 `json:"event_id" bson:"event_id"`
	EventType string                 `json:"event_type" bson:"event_type"`
	GameID    string                 `json:"game_id" bson:"game_id"`
	PlayerID  *uint64                `json:"player_id,omitempty" bson:"player_id,omitempty"`
	Timestamp time.Time              `json:"timestamp" bson:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata" bson:"metadata"`
}

// GetEventID 获取事件ID
func (e *BaseMinigameEvent) GetEventID() string {
	return e.EventID
}

// GetEventType 获取事件类型
func (e *BaseMinigameEvent) GetEventType() string {
	return e.EventType
}

// GetGameID 获取游戏ID
func (e *BaseMinigameEvent) GetGameID() string {
	return e.GameID
}

// GetPlayerID 获取玩家ID
func (e *BaseMinigameEvent) GetPlayerID() *uint64 {
	return e.PlayerID
}

// GetTimestamp 获取时间戳
func (e *BaseMinigameEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetMetadata 获取元数据
func (e *BaseMinigameEvent) GetMetadata() map[string]interface{} {
	return e.Metadata
}

// SetMetadata 设置元数据
func (e *BaseMinigameEvent) SetMetadata(key string, value interface{}) {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
}

// 游戏生命周期事件

// MinigameCreatedEvent 小游戏创建事件
type MinigameCreatedEvent struct {
	BaseMinigameEvent
	GameType  GameType `json:"game_type" bson:"game_type"`
	CreatorID uint64   `json:"creator_id" bson:"creator_id"`
}

// GameStartedEvent 游戏开始事件
type GameStartedEvent struct {
	BaseMinigameEvent
	OperatorID uint64 `json:"operator_id" bson:"operator_id"`
}

// GameEndedEvent 游戏结束事件
type GameEndedEvent struct {
	BaseMinigameEvent
	EndReason  GameEndReason `json:"end_reason" bson:"end_reason"`
	OperatorID uint64        `json:"operator_id" bson:"operator_id"`
	Duration   time.Duration `json:"duration" bson:"duration"`
}

// GamePausedEvent 游戏暂停事件
type GamePausedEvent struct {
	BaseMinigameEvent
	OperatorID uint64 `json:"operator_id" bson:"operator_id"`
}

// GameResumedEvent 游戏恢复事件
type GameResumedEvent struct {
	BaseMinigameEvent
	OperatorID uint64 `json:"operator_id" bson:"operator_id"`
}

// GameCancelledEvent 游戏取消事件
type GameCancelledEvent struct {
	BaseMinigameEvent
	Reason     string `json:"reason" bson:"reason"`
	OperatorID uint64 `json:"operator_id" bson:"operator_id"`
}

// GameResetEvent 游戏重置事件
type GameResetEvent struct {
	BaseMinigameEvent
	OperatorID uint64 `json:"operator_id" bson:"operator_id"`
}

// GameDeletedEvent 游戏删除事件
type GameDeletedEvent struct {
	BaseMinigameEvent
	OperatorID uint64 `json:"operator_id" bson:"operator_id"`
}

// 游戏状态事件

// GameStatusChangedEvent 游戏状态变更事件
type GameStatusChangedEvent struct {
	BaseMinigameEvent
	OldStatus GameStatus `json:"old_status" bson:"old_status"`
	NewStatus GameStatus `json:"new_status" bson:"new_status"`
}

// GamePhaseChangedEvent 游戏阶段变更事件
type GamePhaseChangedEvent struct {
	BaseMinigameEvent
	OldPhase GamePhase `json:"old_phase" bson:"old_phase"`
	NewPhase GamePhase `json:"new_phase" bson:"new_phase"`
}

// GameConfigUpdatedEvent 游戏配置更新事件
type GameConfigUpdatedEvent struct {
	BaseMinigameEvent
	UpdatedFields []string `json:"updated_fields" bson:"updated_fields"`
	OperatorID    uint64   `json:"operator_id" bson:"operator_id"`
}

// GameRulesUpdatedEvent 游戏规则更新事件
type GameRulesUpdatedEvent struct {
	BaseMinigameEvent
	UpdatedRules []string `json:"updated_rules" bson:"updated_rules"`
	OperatorID   uint64   `json:"operator_id" bson:"operator_id"`
}

// GameSettingsUpdatedEvent 游戏设置更新事件
type GameSettingsUpdatedEvent struct {
	BaseMinigameEvent
	UpdatedSettings []string `json:"updated_settings" bson:"updated_settings"`
	OperatorID      uint64   `json:"operator_id" bson:"operator_id"`
}

// 玩家相关事件

// PlayerJoinedEvent 玩家加入事件
type PlayerJoinedEvent struct {
	BaseMinigameEvent
	SessionID string `json:"session_id" bson:"session_id"`
}

// PlayerLeftEvent 玩家离开事件
type PlayerLeftEvent struct {
	BaseMinigameEvent
	LeaveReason PlayerLeaveReason `json:"leave_reason" bson:"leave_reason"`
	SessionID   string            `json:"session_id" bson:"session_id"`
}

// PlayerKickedEvent 玩家被踢出事件
type PlayerKickedEvent struct {
	BaseMinigameEvent
	Reason     string `json:"reason" bson:"reason"`
	OperatorID uint64 `json:"operator_id" bson:"operator_id"`
	SessionID  string `json:"session_id" bson:"session_id"`
}

// PlayerStatusChangedEvent 玩家状态变更事件
type PlayerStatusChangedEvent struct {
	BaseMinigameEvent
	OldStatus PlayerStatus `json:"old_status" bson:"old_status"`
	NewStatus PlayerStatus `json:"new_status" bson:"new_status"`
	SessionID string       `json:"session_id" bson:"session_id"`
}

// PlayerReadyEvent 玩家准备事件
type PlayerReadyEvent struct {
	BaseMinigameEvent
	SessionID string `json:"session_id" bson:"session_id"`
}

// PlayerNotReadyEvent 玩家取消准备事件
type PlayerNotReadyEvent struct {
	BaseMinigameEvent
	SessionID string `json:"session_id" bson:"session_id"`
}

// PlayerDisconnectedEvent 玩家断线事件
type PlayerDisconnectedEvent struct {
	BaseMinigameEvent
	Reason    string `json:"reason" bson:"reason"`
	SessionID string `json:"session_id" bson:"session_id"`
}

// PlayerReconnectedEvent 玩家重连事件
type PlayerReconnectedEvent struct {
	BaseMinigameEvent
	SessionID string `json:"session_id" bson:"session_id"`
}

// 分数和进度事件

// ScoreUpdatedEvent 分数更新事件
type ScoreUpdatedEvent struct {
	BaseMinigameEvent
	ScoreType  ScoreType `json:"score_type" bson:"score_type"`
	OldScore   int64     `json:"old_score" bson:"old_score"`
	NewScore   int64     `json:"new_score" bson:"new_score"`
	FinalScore int64     `json:"final_score" bson:"final_score"`
	SessionID  string    `json:"session_id" bson:"session_id"`
}

// HighScoreAchievedEvent 最高分达成事件
type HighScoreAchievedEvent struct {
	BaseMinigameEvent
	ScoreType    ScoreType `json:"score_type" bson:"score_type"`
	Score        int64     `json:"score" bson:"score"`
	PreviousHigh int64     `json:"previous_high" bson:"previous_high"`
	SessionID    string    `json:"session_id" bson:"session_id"`
}

// LevelUpEvent 升级事件
type LevelUpEvent struct {
	BaseMinigameEvent
	OldLevel  int32  `json:"old_level" bson:"old_level"`
	NewLevel  int32  `json:"new_level" bson:"new_level"`
	SessionID string `json:"session_id" bson:"session_id"`
}

// ProgressUpdatedEvent 进度更新事件
type ProgressUpdatedEvent struct {
	BaseMinigameEvent
	OldProgress float64 `json:"old_progress" bson:"old_progress"`
	NewProgress float64 `json:"new_progress" bson:"new_progress"`
	SessionID   string  `json:"session_id" bson:"session_id"`
}

// MilestoneReachedEvent 里程碑达成事件
type MilestoneReachedEvent struct {
	BaseMinigameEvent
	Milestone string `json:"milestone" bson:"milestone"`
	Value     int64  `json:"value" bson:"value"`
	SessionID string `json:"session_id" bson:"session_id"`
}

// 奖励相关事件

// RewardGrantedEvent 奖励授予事件
type RewardGrantedEvent struct {
	BaseMinigameEvent
	RewardType RewardType `json:"reward_type" bson:"reward_type"`
	ItemID     string     `json:"item_id" bson:"item_id"`
	Quantity   int64      `json:"quantity" bson:"quantity"`
	Source     string     `json:"source" bson:"source"`
	RewardID   string     `json:"reward_id" bson:"reward_id"`
}

// RewardClaimedEvent 奖励领取事件
type RewardClaimedEvent struct {
	BaseMinigameEvent
	RewardType RewardType `json:"reward_type" bson:"reward_type"`
	ItemID     string     `json:"item_id" bson:"item_id"`
	Quantity   int64      `json:"quantity" bson:"quantity"`
	RewardID   string     `json:"reward_id" bson:"reward_id"`
}

// RewardExpiredEvent 奖励过期事件
type RewardExpiredEvent struct {
	BaseMinigameEvent
	RewardType RewardType `json:"reward_type" bson:"reward_type"`
	ItemID     string     `json:"item_id" bson:"item_id"`
	Quantity   int64      `json:"quantity" bson:"quantity"`
	RewardID   string     `json:"reward_id" bson:"reward_id"`
}

// BonusRewardEvent 奖励加成事件
type BonusRewardEvent struct {
	BaseMinigameEvent
	BonusType   string  `json:"bonus_type" bson:"bonus_type"`
	Multiplier  float64 `json:"multiplier" bson:"multiplier"`
	BonusAmount int64   `json:"bonus_amount" bson:"bonus_amount"`
	Reason      string  `json:"reason" bson:"reason"`
}

// 成就相关事件

// AchievementUnlockedEvent 成就解锁事件
type AchievementUnlockedEvent struct {
	BaseMinigameEvent
	AchievementID string `json:"achievement_id" bson:"achievement_id"`
	Name          string `json:"name" bson:"name"`
	Category      string `json:"category" bson:"category"`
	Rarity        string `json:"rarity" bson:"rarity"`
	Points        int64  `json:"points" bson:"points"`
}

// AchievementCompletedEvent 成就完成事件
type AchievementCompletedEvent struct {
	BaseMinigameEvent
	AchievementID string `json:"achievement_id" bson:"achievement_id"`
	Name          string `json:"name" bson:"name"`
	Category      string `json:"category" bson:"category"`
	Rarity        string `json:"rarity" bson:"rarity"`
	Points        int64  `json:"points" bson:"points"`
}

// AchievementProgressEvent 成就进度事件
type AchievementProgressEvent struct {
	BaseMinigameEvent
	AchievementID string  `json:"achievement_id" bson:"achievement_id"`
	OldProgress   float64 `json:"old_progress" bson:"old_progress"`
	NewProgress   float64 `json:"new_progress" bson:"new_progress"`
	MaxProgress   float64 `json:"max_progress" bson:"max_progress"`
}

// 游戏操作事件

// GameActionEvent 游戏动作事件
type GameActionEvent struct {
	BaseMinigameEvent
	Action     string                 `json:"action" bson:"action"`
	Parameters map[string]interface{} `json:"parameters" bson:"parameters"`
	Result     string                 `json:"result" bson:"result"`
	SessionID  string                 `json:"session_id" bson:"session_id"`
}

// GameMoveEvent 游戏移动事件
type GameMoveEvent struct {
	BaseMinigameEvent
	MoveType   string                 `json:"move_type" bson:"move_type"`
	MoveData   map[string]interface{} `json:"move_data" bson:"move_data"`
	MoveNumber int32                  `json:"move_number" bson:"move_number"`
	SessionID  string                 `json:"session_id" bson:"session_id"`
}

// GameInputEvent 游戏输入事件
type GameInputEvent struct {
	BaseMinigameEvent
	InputType string                 `json:"input_type" bson:"input_type"`
	InputData map[string]interface{} `json:"input_data" bson:"input_data"`
	SessionID string                 `json:"session_id" bson:"session_id"`
}

// GameOutputEvent 游戏输出事件
type GameOutputEvent struct {
	BaseMinigameEvent
	OutputType string                 `json:"output_type" bson:"output_type"`
	OutputData map[string]interface{} `json:"output_data" bson:"output_data"`
	SessionID  string                 `json:"session_id" bson:"session_id"`
}

// 系统事件

// GameErrorEvent 游戏错误事件
type GameErrorEvent struct {
	BaseMinigameEvent
	ErrorCode    string `json:"error_code" bson:"error_code"`
	ErrorMessage string `json:"error_message" bson:"error_message"`
	ErrorType    string `json:"error_type" bson:"error_type"`
	StackTrace   string `json:"stack_trace,omitempty" bson:"stack_trace,omitempty"`
	SessionID    string `json:"session_id,omitempty" bson:"session_id,omitempty"`
}

// GameWarningEvent 游戏警告事件
type GameWarningEvent struct {
	BaseMinigameEvent
	WarningCode    string `json:"warning_code" bson:"warning_code"`
	WarningMessage string `json:"warning_message" bson:"warning_message"`
	WarningType    string `json:"warning_type" bson:"warning_type"`
	SessionID      string `json:"session_id,omitempty" bson:"session_id,omitempty"`
}

// GameMaintenanceEvent 游戏维护事件
type GameMaintenanceEvent struct {
	BaseMinigameEvent
	MaintenanceType string    `json:"maintenance_type" bson:"maintenance_type"`
	StartTime       time.Time `json:"start_time" bson:"start_time"`
	EndTime         time.Time `json:"end_time" bson:"end_time"`
	Reason          string    `json:"reason" bson:"reason"`
}

// GameUpdateEvent 游戏更新事件
type GameUpdateEvent struct {
	BaseMinigameEvent
	UpdateType    string `json:"update_type" bson:"update_type"`
	Version       string `json:"version" bson:"version"`
	UpdateDetails string `json:"update_details" bson:"update_details"`
}

// 统计事件

// GameStatisticsEvent 游戏统计事件
type GameStatisticsEvent struct {
	BaseMinigameEvent
	StatisticsType string                 `json:"statistics_type" bson:"statistics_type"`
	Statistics     map[string]interface{} `json:"statistics" bson:"statistics"`
	Period         string                 `json:"period" bson:"period"`
}

// PlayerStatisticsEvent 玩家统计事件
type PlayerStatisticsEvent struct {
	BaseMinigameEvent
	StatisticsType string                 `json:"statistics_type" bson:"statistics_type"`
	Statistics     map[string]interface{} `json:"statistics" bson:"statistics"`
	Period         string                 `json:"period" bson:"period"`
}

// LeaderboardUpdatedEvent 排行榜更新事件
type LeaderboardUpdatedEvent struct {
	BaseMinigameEvent
	LeaderboardType string    `json:"leaderboard_type" bson:"leaderboard_type"`
	ScoreType       ScoreType `json:"score_type" bson:"score_type"`
	TopPlayers      []uint64  `json:"top_players" bson:"top_players"`
	UpdateReason    string    `json:"update_reason" bson:"update_reason"`
}

// 事件常量

const (
	// 游戏生命周期事件类型
	EventTypeMinigameCreated = "minigame.created"
	EventTypeGameStarted     = "game.started"
	EventTypeGameEnded       = "game.ended"
	EventTypeGamePaused      = "game.paused"
	EventTypeGameResumed     = "game.resumed"
	EventTypeGameCancelled   = "game.cancelled"
	EventTypeGameReset       = "game.reset"
	EventTypeGameDeleted     = "game.deleted"
	
	// 游戏状态事件类型
	EventTypeGameStatusChanged  = "game.status_changed"
	EventTypeGamePhaseChanged   = "game.phase_changed"
	EventTypeGameConfigUpdated  = "game.config_updated"
	EventTypeGameRulesUpdated   = "game.rules_updated"
	EventTypeGameSettingsUpdated = "game.settings_updated"
	
	// 玩家相关事件类型
	EventTypePlayerJoined       = "player.joined"
	EventTypePlayerLeft         = "player.left"
	EventTypePlayerKicked       = "player.kicked"
	EventTypePlayerStatusChanged = "player.status_changed"
	EventTypePlayerReady        = "player.ready"
	EventTypePlayerNotReady     = "player.not_ready"
	EventTypePlayerDisconnected = "player.disconnected"
	EventTypePlayerReconnected  = "player.reconnected"
	
	// 分数和进度事件类型
	EventTypeScoreUpdated       = "score.updated"
	EventTypeHighScoreAchieved  = "score.high_score_achieved"
	EventTypeLevelUp            = "progress.level_up"
	EventTypeProgressUpdated    = "progress.updated"
	EventTypeMilestoneReached   = "progress.milestone_reached"
	
	// 奖励相关事件类型
	EventTypeRewardGranted = "reward.granted"
	EventTypeRewardClaimed = "reward.claimed"
	EventTypeRewardExpired = "reward.expired"
	EventTypeBonusReward   = "reward.bonus"
	
	// 成就相关事件类型
	EventTypeAchievementUnlocked  = "achievement.unlocked"
	EventTypeAchievementCompleted = "achievement.completed"
	EventTypeAchievementProgress  = "achievement.progress"
	
	// 游戏操作事件类型
	EventTypeGameAction = "game.action"
	EventTypeGameMove   = "game.move"
	EventTypeGameInput  = "game.input"
	EventTypeGameOutput = "game.output"
	
	// 系统事件类型
	EventTypeGameError       = "system.error"
	EventTypeGameWarning     = "system.warning"
	EventTypeGameMaintenance = "system.maintenance"
	EventTypeGameUpdate      = "system.update"
	
	// 统计事件类型
	EventTypeGameStatistics      = "statistics.game"
	EventTypePlayerStatistics    = "statistics.player"
	EventTypeLeaderboardUpdated  = "statistics.leaderboard_updated"
)

// 事件工厂函数

// NewMinigameCreatedEvent 创建小游戏创建事件
func NewMinigameCreatedEvent(gameID string, gameType GameType, creatorID uint64) *MinigameCreatedEvent {
	return &MinigameCreatedEvent{
		BaseMinigameEvent: BaseMinigameEvent{
			EventID:   generateEventID(),
			EventType: EventTypeMinigameCreated,
			GameID:    gameID,
			Timestamp: time.Now(),
			Metadata:  make(map[string]interface{}),
		},
		GameType:  gameType,
		CreatorID: creatorID,
	}
}

// NewGameStartedEvent 创建游戏开始事件
func NewGameStartedEvent(gameID string, operatorID uint64) *GameStartedEvent {
	return &GameStartedEvent{
		BaseMinigameEvent: BaseMinigameEvent{
			EventID:   generateEventID(),
			EventType: EventTypeGameStarted,
			GameID:    gameID,
			Timestamp: time.Now(),
			Metadata:  make(map[string]interface{}),
		},
		OperatorID: operatorID,
	}
}

// NewGameEndedEvent 创建游戏结束事件
func NewGameEndedEvent(gameID string, endReason GameEndReason, operatorID uint64) *GameEndedEvent {
	return &GameEndedEvent{
		BaseMinigameEvent: BaseMinigameEvent{
			EventID:   generateEventID(),
			EventType: EventTypeGameEnded,
			GameID:    gameID,
			Timestamp: time.Now(),
			Metadata:  make(map[string]interface{}),
		},
		EndReason:  endReason,
		OperatorID: operatorID,
	}
}

// NewPlayerJoinedEvent 创建玩家加入事件
func NewPlayerJoinedEvent(gameID string, playerID uint64, sessionID string) *PlayerJoinedEvent {
	return &PlayerJoinedEvent{
		BaseMinigameEvent: BaseMinigameEvent{
			EventID:   generateEventID(),
			EventType: EventTypePlayerJoined,
			GameID:    gameID,
			PlayerID:  &playerID,
			Timestamp: time.Now(),
			Metadata:  make(map[string]interface{}),
		},
		SessionID: sessionID,
	}
}

// NewPlayerLeftEvent 创建玩家离开事件
func NewPlayerLeftEvent(gameID string, playerID uint64, leaveReason PlayerLeaveReason) *PlayerLeftEvent {
	return &PlayerLeftEvent{
		BaseMinigameEvent: BaseMinigameEvent{
			EventID:   generateEventID(),
			EventType: EventTypePlayerLeft,
			GameID:    gameID,
			PlayerID:  &playerID,
			Timestamp: time.Now(),
			Metadata:  make(map[string]interface{}),
		},
		LeaveReason: leaveReason,
	}
}

// NewScoreUpdatedEvent 创建分数更新事件
func NewScoreUpdatedEvent(gameID string, playerID uint64, scoreType ScoreType, oldScore, newScore, finalScore int64) *ScoreUpdatedEvent {
	return &ScoreUpdatedEvent{
		BaseMinigameEvent: BaseMinigameEvent{
			EventID:   generateEventID(),
			EventType: EventTypeScoreUpdated,
			GameID:    gameID,
			PlayerID:  &playerID,
			Timestamp: time.Now(),
			Metadata:  make(map[string]interface{}),
		},
		ScoreType:  scoreType,
		OldScore:   oldScore,
		NewScore:   newScore,
		FinalScore: finalScore,
	}
}

// NewRewardGrantedEvent 创建奖励授予事件
func NewRewardGrantedEvent(gameID string, playerID uint64, rewardType RewardType, itemID string, quantity int64) *RewardGrantedEvent {
	return &RewardGrantedEvent{
		BaseMinigameEvent: BaseMinigameEvent{
			EventID:   generateEventID(),
			EventType: EventTypeRewardGranted,
			GameID:    gameID,
			PlayerID:  &playerID,
			Timestamp: time.Now(),
			Metadata:  make(map[string]interface{}),
		},
		RewardType: rewardType,
		ItemID:     itemID,
		Quantity:   quantity,
	}
}

// NewRewardClaimedEvent 创建奖励领取事件
func NewRewardClaimedEvent(gameID string, playerID uint64, rewardType RewardType, itemID string, quantity int64) *RewardClaimedEvent {
	return &RewardClaimedEvent{
		BaseMinigameEvent: BaseMinigameEvent{
			EventID:   generateEventID(),
			EventType: EventTypeRewardClaimed,
			GameID:    gameID,
			PlayerID:  &playerID,
			Timestamp: time.Now(),
			Metadata:  make(map[string]interface{}),
		},
		RewardType: rewardType,
		ItemID:     itemID,
		Quantity:   quantity,
	}
}

// NewAchievementCompletedEvent 创建成就完成事件
func NewAchievementCompletedEvent(gameID string, playerID uint64, achievementID string, points int64) *AchievementCompletedEvent {
	return &AchievementCompletedEvent{
		BaseMinigameEvent: BaseMinigameEvent{
			EventID:   generateEventID(),
			EventType: EventTypeAchievementCompleted,
			GameID:    gameID,
			PlayerID:  &playerID,
			Timestamp: time.Now(),
			Metadata:  make(map[string]interface{}),
		},
		AchievementID: achievementID,
		Points:        points,
	}
}

// NewGameErrorEvent 创建游戏错误事件
func NewGameErrorEvent(gameID string, errorCode, errorMessage, errorType string) *GameErrorEvent {
	return &GameErrorEvent{
		BaseMinigameEvent: BaseMinigameEvent{
			EventID:   generateEventID(),
			EventType: EventTypeGameError,
			GameID:    gameID,
			Timestamp: time.Now(),
			Metadata:  make(map[string]interface{}),
		},
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
		ErrorType:    errorType,
	}
}

// 事件处理器接口

// MinigameEventHandler 小游戏事件处理器接口
type MinigameEventHandler interface {
	Handle(ctx context.Context, event MinigameEvent) error
	CanHandle(eventType string) bool
	GetHandlerName() string
}

// MinigameEventBus 小游戏事件总线接口
type MinigameEventBus interface {
	Publish(ctx context.Context, event MinigameEvent) error
	Subscribe(eventType string, handler MinigameEventHandler) error
	Unsubscribe(eventType string, handlerName string) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// 事件聚合器

// EventAggregator 事件聚合器
type EventAggregator struct {
	Events []MinigameEvent `json:"events"`
	Count  int64           `json:"count"`
	Period string          `json:"period"`
	From   time.Time       `json:"from"`
	To     time.Time       `json:"to"`
}

// NewEventAggregator 创建事件聚合器
func NewEventAggregator(period string, from, to time.Time) *EventAggregator {
	return &EventAggregator{
		Events: make([]MinigameEvent, 0),
		Count:  0,
		Period: period,
		From:   from,
		To:     to,
	}
}

// AddEvent 添加事件
func (ea *EventAggregator) AddEvent(event MinigameEvent) {
	ea.Events = append(ea.Events, event)
	ea.Count++
}

// GetEventsByType 根据类型获取事件
func (ea *EventAggregator) GetEventsByType(eventType string) []MinigameEvent {
	var events []MinigameEvent
	for _, event := range ea.Events {
		if event.GetEventType() == eventType {
			events = append(events, event)
		}
	}
	return events
}

// GetEventsByPlayer 根据玩家获取事件
func (ea *EventAggregator) GetEventsByPlayer(playerID uint64) []MinigameEvent {
	var events []MinigameEvent
	for _, event := range ea.Events {
		if event.GetPlayerID() != nil && *event.GetPlayerID() == playerID {
			events = append(events, event)
		}
	}
	return events
}

// GetEventsByGame 根据游戏获取事件
func (ea *EventAggregator) GetEventsByGame(gameID string) []MinigameEvent {
	var events []MinigameEvent
	for _, event := range ea.Events {
		if event.GetGameID() == gameID {
			events = append(events, event)
		}
	}
	return events
}

// GetEventStatistics 获取事件统计
func (ea *EventAggregator) GetEventStatistics() map[string]int64 {
	stats := make(map[string]int64)
	for _, event := range ea.Events {
		stats[event.GetEventType()]++
	}
	return stats
}

// Clear 清空事件
func (ea *EventAggregator) Clear() {
	ea.Events = make([]MinigameEvent, 0)
	ea.Count = 0
}

// 辅助函数

// generateEventID 生成事件ID
func generateEventID() string {
	return fmt.Sprintf("evt_%d_%d", time.Now().UnixNano(), time.Now().Unix())
}

// IsGameLifecycleEvent 检查是否为游戏生命周期事件
func IsGameLifecycleEvent(eventType string) bool {
	lifecycleEvents := []string{
		EventTypeMinigameCreated,
		EventTypeGameStarted,
		EventTypeGameEnded,
		EventTypeGamePaused,
		EventTypeGameResumed,
		EventTypeGameCancelled,
		EventTypeGameReset,
		EventTypeGameDeleted,
	}
	
	for _, event := range lifecycleEvents {
		if event == eventType {
			return true
		}
	}
	return false
}

// IsPlayerEvent 检查是否为玩家事件
func IsPlayerEvent(eventType string) bool {
	playerEvents := []string{
		EventTypePlayerJoined,
		EventTypePlayerLeft,
		EventTypePlayerKicked,
		EventTypePlayerStatusChanged,
		EventTypePlayerReady,
		EventTypePlayerNotReady,
		EventTypePlayerDisconnected,
		EventTypePlayerReconnected,
	}
	
	for _, event := range playerEvents {
		if event == eventType {
			return true
		}
	}
	return false
}

// IsScoreEvent 检查是否为分数事件
func IsScoreEvent(eventType string) bool {
	scoreEvents := []string{
		EventTypeScoreUpdated,
		EventTypeHighScoreAchieved,
		EventTypeLevelUp,
		EventTypeProgressUpdated,
		EventTypeMilestoneReached,
	}
	
	for _, event := range scoreEvents {
		if event == eventType {
			return true
		}
	}
	return false
}

// IsRewardEvent 检查是否为奖励事件
func IsRewardEvent(eventType string) bool {
	rewardEvents := []string{
		EventTypeRewardGranted,
		EventTypeRewardClaimed,
		EventTypeRewardExpired,
		EventTypeBonusReward,
	}
	
	for _, event := range rewardEvents {
		if event == eventType {
			return true
		}
	}
	return false
}

// IsAchievementEvent 检查是否为成就事件
func IsAchievementEvent(eventType string) bool {
	achievementEvents := []string{
		EventTypeAchievementUnlocked,
		EventTypeAchievementCompleted,
		EventTypeAchievementProgress,
	}
	
	for _, event := range achievementEvents {
		if event == eventType {
			return true
		}
	}
	return false
}

// IsSystemEvent 检查是否为系统事件
func IsSystemEvent(eventType string) bool {
	systemEvents := []string{
		EventTypeGameError,
		EventTypeGameWarning,
		EventTypeGameMaintenance,
		EventTypeGameUpdate,
	}
	
	for _, event := range systemEvents {
		if event == eventType {
			return true
		}
	}
	return false
}

// GetEventPriority 获取事件优先级
func GetEventPriority(eventType string) int {
	switch eventType {
	case EventTypeGameError:
		return 1 // 最高优先级
	case EventTypeGameWarning:
		return 2
	case EventTypeGameMaintenance:
		return 3
	case EventTypeMinigameCreated, EventTypeGameStarted, EventTypeGameEnded:
		return 4
	case EventTypePlayerJoined, EventTypePlayerLeft:
		return 5
	case EventTypeScoreUpdated, EventTypeHighScoreAchieved:
		return 6
	case EventTypeRewardGranted, EventTypeRewardClaimed:
		return 7
	case EventTypeAchievementCompleted, EventTypeAchievementUnlocked:
		return 8
	default:
		return 9 // 最低优先级
	}
}

// FilterEventsByTimeRange 根据时间范围过滤事件
func FilterEventsByTimeRange(events []MinigameEvent, from, to time.Time) []MinigameEvent {
	var filtered []MinigameEvent
	for _, event := range events {
		timestamp := event.GetTimestamp()
		if timestamp.After(from) && timestamp.Before(to) {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// GroupEventsByType 根据类型分组事件
func GroupEventsByType(events []MinigameEvent) map[string][]MinigameEvent {
	groups := make(map[string][]MinigameEvent)
	for _, event := range events {
		eventType := event.GetEventType()
		groups[eventType] = append(groups[eventType], event)
	}
	return groups
}

// GroupEventsByPlayer 根据玩家分组事件
func GroupEventsByPlayer(events []MinigameEvent) map[uint64][]MinigameEvent {
	groups := make(map[uint64][]MinigameEvent)
	for _, event := range events {
		if playerID := event.GetPlayerID(); playerID != nil {
			groups[*playerID] = append(groups[*playerID], event)
		}
	}
	return groups
}

// GroupEventsByGame 根据游戏分组事件
func GroupEventsByGame(events []MinigameEvent) map[string][]MinigameEvent {
	groups := make(map[string][]MinigameEvent)
	for _, event := range events {
		gameID := event.GetGameID()
		groups[gameID] = append(groups[gameID], event)
	}
	return groups
}