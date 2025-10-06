package minigame

import (
	"fmt"
	"time"
)

// 小游戏类型相关值对象

// GameType 游戏类型
type GameType int32

const (
	GameTypeSaveDog  GameType = iota + 1 // 拯救小狗
	GameTypePuzzle                       // 拼图游戏
	GameTypeRacing                       // 赛车游戏
	GameTypeMemory                       // 记忆游戏
	GameTypeMatch                        // 消除游戏
	GameTypeJump                         // 跳跃游戏
	GameTypeShoot                        // 射击游戏
	GameTypeStrategy                     // 策略游戏
	GameTypeCard                         // 卡牌游戏
	GameTypeCustom                       // 自定义游戏
)

// String 返回游戏类型的字符串表示
func (gt GameType) String() string {
	switch gt {
	case GameTypeSaveDog:
		return "save_dog"
	case GameTypePuzzle:
		return "puzzle"
	case GameTypeRacing:
		return "racing"
	case GameTypeMemory:
		return "memory"
	case GameTypeMatch:
		return "match"
	case GameTypeJump:
		return "jump"
	case GameTypeShoot:
		return "shoot"
	case GameTypeStrategy:
		return "strategy"
	case GameTypeCard:
		return "card"
	case GameTypeCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// IsValid 检查游戏类型是否有效
func (gt GameType) IsValid() bool {
	return gt >= GameTypeSaveDog && gt <= GameTypeCustom
}

// 注意：GameCategory已经在types.go中定义，这里删除重复定义

// String 返回游戏分类的字符串表示
func (gc GameCategory) String() string {
	switch gc {
	case GameCategoryNormal:
		return "normal"
	case GameCategoryCompetitive:
		return "competitive"
	case GameCategoryCasual:
		return "casual"
	case GameCategoryRanked:
		return "ranked"
	default:
		return "unknown"
	}
}

// IsValid 检查游戏分类是否有效
func (gc GameCategory) IsValid() bool {
	switch gc {
	case GameCategoryNormal, GameCategoryCompetitive, GameCategoryCasual, GameCategoryRanked:
		return true
	default:
		return false
	}
}

// GameStatus 游戏状态
type GameStatus int32

const (
	GameStatusWaiting   GameStatus = iota + 1 // 等待开始
	GameStatusRunning                         // 进行中
	GameStatusPaused                          // 暂停
	GameStatusFinished                        // 已结束
	GameStatusCancelled                       // 已取消
	GameStatusError                           // 错误状态
)

// String 返回游戏状态的字符串表示
func (gs GameStatus) String() string {
	switch gs {
	case GameStatusWaiting:
		return "waiting"
	case GameStatusRunning:
		return "running"
	case GameStatusPaused:
		return "paused"
	case GameStatusFinished:
		return "finished"
	case GameStatusCancelled:
		return "cancelled"
	case GameStatusError:
		return "error"
	default:
		return "unknown"
	}
}

// IsValid 检查游戏状态是否有效
func (gs GameStatus) IsValid() bool {
	return gs >= GameStatusWaiting && gs <= GameStatusError
}

// CanTransitionTo 检查是否可以转换到目标状态
func (gs GameStatus) CanTransitionTo(target GameStatus) bool {
	switch gs {
	case GameStatusWaiting:
		return target == GameStatusRunning || target == GameStatusCancelled
	case GameStatusRunning:
		return target == GameStatusPaused || target == GameStatusFinished || target == GameStatusCancelled || target == GameStatusError
	case GameStatusPaused:
		return target == GameStatusRunning || target == GameStatusFinished || target == GameStatusCancelled
	case GameStatusFinished, GameStatusCancelled, GameStatusError:
		return false // 终态，不能转换
	default:
		return false
	}
}

// GamePhase 游戏阶段
type GamePhase int32

const (
	GamePhaseWaiting  GamePhase = iota + 1 // 等待阶段
	GamePhaseStarting                      // 开始阶段
	GamePhaseRunning                       // 运行阶段
	GamePhasePaused                        // 暂停阶段
	GamePhaseEnding                        // 结束阶段
	GamePhaseFinished                      // 完成阶段
)

// String 返回游戏阶段的字符串表示
func (gp GamePhase) String() string {
	switch gp {
	case GamePhaseWaiting:
		return "waiting"
	case GamePhaseStarting:
		return "starting"
	case GamePhaseRunning:
		return "running"
	case GamePhasePaused:
		return "paused"
	case GamePhaseEnding:
		return "ending"
	case GamePhaseFinished:
		return "finished"
	default:
		return "unknown"
	}
}

// IsValid 检查游戏阶段是否有效
func (gp GamePhase) IsValid() bool {
	return gp >= GamePhaseWaiting && gp <= GamePhaseFinished
}

// GameEndReason 游戏结束原因
type GameEndReason int32

const (
	GameEndReasonCompleted           GameEndReason = iota + 1 // 正常完成
	GameEndReasonTimeout                                      // 超时
	GameEndReasonCancelled                                    // 取消
	GameEndReasonInsufficientPlayers                          // 玩家不足
	GameEndReasonError                                        // 错误
	GameEndReasonForceQuit                                    // 强制退出
	GameEndReasonNetworkError                                 // 网络错误
	GameEndReasonSystemMaintenance                            // 系统维护
)

// String 返回游戏结束原因的字符串表示
func (ger GameEndReason) String() string {
	switch ger {
	case GameEndReasonCompleted:
		return "completed"
	case GameEndReasonTimeout:
		return "timeout"
	case GameEndReasonCancelled:
		return "cancelled"
	case GameEndReasonInsufficientPlayers:
		return "insufficient_players"
	case GameEndReasonError:
		return "error"
	case GameEndReasonForceQuit:
		return "force_quit"
	case GameEndReasonNetworkError:
		return "network_error"
	case GameEndReasonSystemMaintenance:
		return "system_maintenance"
	default:
		return "unknown"
	}
}

// IsValid 检查游戏结束原因是否有效
func (ger GameEndReason) IsValid() bool {
	return ger >= GameEndReasonCompleted && ger <= GameEndReasonSystemMaintenance
}

// PlayerStatus 玩家状态
type PlayerStatus int32

const (
	PlayerStatusWaiting      PlayerStatus = iota + 1 // 等待中
	PlayerStatusReady                                // 准备就绪
	PlayerStatusPlaying                              // 游戏中
	PlayerStatusPaused                               // 暂停
	PlayerStatusFinished                             // 已完成
	PlayerStatusLeft                                 // 已离开
	PlayerStatusDisconnected                         // 断线
	PlayerStatusKicked                               // 被踢出
)

// String 返回玩家状态的字符串表示
func (ps PlayerStatus) String() string {
	switch ps {
	case PlayerStatusWaiting:
		return "waiting"
	case PlayerStatusReady:
		return "ready"
	case PlayerStatusPlaying:
		return "playing"
	case PlayerStatusPaused:
		return "paused"
	case PlayerStatusFinished:
		return "finished"
	case PlayerStatusLeft:
		return "left"
	case PlayerStatusDisconnected:
		return "disconnected"
	case PlayerStatusKicked:
		return "kicked"
	default:
		return "unknown"
	}
}

// IsValid 检查玩家状态是否有效
func (ps PlayerStatus) IsValid() bool {
	return ps >= PlayerStatusWaiting && ps <= PlayerStatusKicked
}

// PlayerLeaveReason 玩家离开原因
type PlayerLeaveReason int32

const (
	PlayerLeaveReasonVoluntary    PlayerLeaveReason = iota + 1 // 主动离开
	PlayerLeaveReasonKicked                                    // 被踢出
	PlayerLeaveReasonDisconnected                              // 断线
	PlayerLeaveReasonTimeout                                   // 超时
	PlayerLeaveReasonError                                     // 错误
	PlayerLeaveReasonGameEnded                                 // 游戏结束
)

// String 返回玩家离开原因的字符串表示
func (plr PlayerLeaveReason) String() string {
	switch plr {
	case PlayerLeaveReasonVoluntary:
		return "voluntary"
	case PlayerLeaveReasonKicked:
		return "kicked"
	case PlayerLeaveReasonDisconnected:
		return "disconnected"
	case PlayerLeaveReasonTimeout:
		return "timeout"
	case PlayerLeaveReasonError:
		return "error"
	case PlayerLeaveReasonGameEnded:
		return "game_ended"
	default:
		return "unknown"
	}
}

// IsValid 检查玩家离开原因是否有效
func (plr PlayerLeaveReason) IsValid() bool {
	return plr >= PlayerLeaveReasonVoluntary && plr <= PlayerLeaveReasonGameEnded
}

// ScoreType 分数类型
type ScoreType int32

const (
	ScoreTypePoints   ScoreType = iota + 1 // 积分
	ScoreTypeTime                          // 时间
	ScoreTypeDistance                      // 距离
	ScoreTypeAccuracy                      // 准确度
	ScoreTypeCombo                         // 连击
	ScoreTypeLevel                         // 等级
	ScoreTypeProgress                      // 进度
	ScoreTypeCustom                        // 自定义
)

// String 返回分数类型的字符串表示
func (st ScoreType) String() string {
	switch st {
	case ScoreTypePoints:
		return "points"
	case ScoreTypeTime:
		return "time"
	case ScoreTypeDistance:
		return "distance"
	case ScoreTypeAccuracy:
		return "accuracy"
	case ScoreTypeCombo:
		return "combo"
	case ScoreTypeLevel:
		return "level"
	case ScoreTypeProgress:
		return "progress"
	case ScoreTypeCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// IsValid 检查分数类型是否有效
func (st ScoreType) IsValid() bool {
	return st >= ScoreTypePoints && st <= ScoreTypeCustom
}

// 注意：RewardType已经在entity.go中定义，这里删除重复定义

// String 返回奖励类型的字符串表示
func (rt RewardType) String() string {
	switch rt {
	case RewardTypeCoin:
		return "coin"
	case RewardTypeExp:
		return "exp"
	case RewardTypeItem:
		return "item"
	case RewardTypeCurrency:
		return "currency"
	default:
		return "unknown"
	}
}

// IsValid 检查奖励类型是否有效
func (rt RewardType) IsValid() bool {
	switch rt {
	case RewardTypeCoin, RewardTypeExp, RewardTypeItem, RewardTypeCurrency:
		return true
	default:
		return false
	}
}

// 游戏配置相关值对象

// GameConfig 游戏配置
type GameConfig struct {
	MaxPlayers      int32          `json:"max_players" bson:"max_players"`
	MinPlayers      int32          `json:"min_players" bson:"min_players"`
	MaxDuration     time.Duration  `json:"max_duration" bson:"max_duration"`
	MinDuration     time.Duration  `json:"min_duration" bson:"min_duration"`
	AutoStart       bool           `json:"auto_start" bson:"auto_start"`
	AutoEnd         bool           `json:"auto_end" bson:"auto_end"`
	AllowSpectators bool           `json:"allow_spectators" bson:"allow_spectators"`
	AllowReconnect  bool           `json:"allow_reconnect" bson:"allow_reconnect"`
	Difficulty      GameDifficulty `json:"difficulty" bson:"difficulty"`
	CreatedAt       time.Time      `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" bson:"updated_at"`
}

// NewGameConfig 创建新的游戏配置
func NewGameConfig() *GameConfig {
	now := time.Now()
	return &GameConfig{
		MaxPlayers:      10,
		MinPlayers:      1,
		MaxDuration:     30 * time.Minute,
		MinDuration:     1 * time.Minute,
		AutoStart:       false,
		AutoEnd:         true,
		AllowSpectators: true,
		AllowReconnect:  true,
		Difficulty:      GameDifficultyNormal,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// Clone 克隆游戏配置
func (gc *GameConfig) Clone() *GameConfig {
	return &GameConfig{
		MaxPlayers:      gc.MaxPlayers,
		MinPlayers:      gc.MinPlayers,
		MaxDuration:     gc.MaxDuration,
		MinDuration:     gc.MinDuration,
		AutoStart:       gc.AutoStart,
		AutoEnd:         gc.AutoEnd,
		AllowSpectators: gc.AllowSpectators,
		AllowReconnect:  gc.AllowReconnect,
		Difficulty:      gc.Difficulty,
		CreatedAt:       gc.CreatedAt,
		UpdatedAt:       gc.UpdatedAt,
	}
}

// GameDifficulty 游戏难度
type GameDifficulty int32

const (
	GameDifficultyEasy   GameDifficulty = iota + 1 // 简单
	GameDifficultyNormal                           // 普通
	GameDifficultyHard                             // 困难
	GameDifficultyExpert                           // 专家
	GameDifficultyCustom                           // 自定义
)

// String 返回游戏难度的字符串表示
func (gd GameDifficulty) String() string {
	switch gd {
	case GameDifficultyEasy:
		return "easy"
	case GameDifficultyNormal:
		return "normal"
	case GameDifficultyHard:
		return "hard"
	case GameDifficultyExpert:
		return "expert"
	case GameDifficultyCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// IsValid 检查游戏难度是否有效
func (gd GameDifficulty) IsValid() bool {
	return gd >= GameDifficultyEasy && gd <= GameDifficultyCustom
}

// GetScoreMultiplier 获取分数倍数
func (gd GameDifficulty) GetScoreMultiplier() float64 {
	switch gd {
	case GameDifficultyEasy:
		return 0.8
	case GameDifficultyNormal:
		return 1.0
	case GameDifficultyHard:
		return 1.5
	case GameDifficultyExpert:
		return 2.0
	case GameDifficultyCustom:
		return 1.0
	default:
		return 1.0
	}
}

// GameRules 游戏规则
type GameRules struct {
	WinConditions  []WinCondition         `json:"win_conditions" bson:"win_conditions"`
	LoseConditions []LoseCondition        `json:"lose_conditions" bson:"lose_conditions"`
	ScoringRules   []ScoringRule          `json:"scoring_rules" bson:"scoring_rules"`
	TimeLimit      *time.Duration         `json:"time_limit,omitempty" bson:"time_limit,omitempty"`
	MoveLimit      *int32                 `json:"move_limit,omitempty" bson:"move_limit,omitempty"`
	SpecialRules   map[string]interface{} `json:"special_rules" bson:"special_rules"`
	Penalties      []Penalty              `json:"penalties" bson:"penalties"`
	Bonuses        []Bonus                `json:"bonuses" bson:"bonuses"`
	CreatedAt      time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at" bson:"updated_at"`
}

// NewGameRules 创建新的游戏规则
func NewGameRules() *GameRules {
	now := time.Now()
	return &GameRules{
		WinConditions:  make([]WinCondition, 0),
		LoseConditions: make([]LoseCondition, 0),
		ScoringRules:   make([]ScoringRule, 0),
		SpecialRules:   make(map[string]interface{}),
		Penalties:      make([]Penalty, 0),
		Bonuses:        make([]Bonus, 0),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// Clone 克隆游戏规则
func (gr *GameRules) Clone() *GameRules {
	clone := &GameRules{
		WinConditions:  make([]WinCondition, len(gr.WinConditions)),
		LoseConditions: make([]LoseCondition, len(gr.LoseConditions)),
		ScoringRules:   make([]ScoringRule, len(gr.ScoringRules)),
		SpecialRules:   make(map[string]interface{}),
		Penalties:      make([]Penalty, len(gr.Penalties)),
		Bonuses:        make([]Bonus, len(gr.Bonuses)),
		CreatedAt:      gr.CreatedAt,
		UpdatedAt:      gr.UpdatedAt,
	}

	// 深拷贝切片
	copy(clone.WinConditions, gr.WinConditions)
	copy(clone.LoseConditions, gr.LoseConditions)
	copy(clone.ScoringRules, gr.ScoringRules)
	copy(clone.Penalties, gr.Penalties)
	copy(clone.Bonuses, gr.Bonuses)

	// 深拷贝map
	for k, v := range gr.SpecialRules {
		clone.SpecialRules[k] = v
	}

	// 深拷贝指针
	if gr.TimeLimit != nil {
		timeLimit := *gr.TimeLimit
		clone.TimeLimit = &timeLimit
	}
	if gr.MoveLimit != nil {
		moveLimit := *gr.MoveLimit
		clone.MoveLimit = &moveLimit
	}

	return clone
}

// WinCondition 胜利条件
type WinCondition struct {
	Type        string      `json:"type" bson:"type"`
	Description string      `json:"description" bson:"description"`
	Target      interface{} `json:"target" bson:"target"`
	Operator    string      `json:"operator" bson:"operator"`
	Priority    int32       `json:"priority" bson:"priority"`
}

// LoseCondition 失败条件
type LoseCondition struct {
	Type        string      `json:"type" bson:"type"`
	Description string      `json:"description" bson:"description"`
	Target      interface{} `json:"target" bson:"target"`
	Operator    string      `json:"operator" bson:"operator"`
	Priority    int32       `json:"priority" bson:"priority"`
}

// ScoringRule 计分规则
type ScoringRule struct {
	Action      string  `json:"action" bson:"action"`
	Points      int64   `json:"points" bson:"points"`
	Multiplier  float64 `json:"multiplier" bson:"multiplier"`
	Description string  `json:"description" bson:"description"`
}

// Penalty 惩罚
type Penalty struct {
	Trigger     string `json:"trigger" bson:"trigger"`
	Penalty     int64  `json:"penalty" bson:"penalty"`
	Description string `json:"description" bson:"description"`
}

// Bonus 奖励
type Bonus struct {
	Trigger     string  `json:"trigger" bson:"trigger"`
	Bonus       int64   `json:"bonus" bson:"bonus"`
	Multiplier  float64 `json:"multiplier" bson:"multiplier"`
	Description string  `json:"description" bson:"description"`
}

// GameSettings 游戏设置
type GameSettings struct {
	SoundEnabled   bool                   `json:"sound_enabled" bson:"sound_enabled"`
	MusicEnabled   bool                   `json:"music_enabled" bson:"music_enabled"`
	EffectsEnabled bool                   `json:"effects_enabled" bson:"effects_enabled"`
	Language       string                 `json:"language" bson:"language"`
	Theme          string                 `json:"theme" bson:"theme"`
	Quality        GameQuality            `json:"quality" bson:"quality"`
	Controls       map[string]interface{} `json:"controls" bson:"controls"`
	CustomSettings map[string]interface{} `json:"custom_settings" bson:"custom_settings"`
	CreatedAt      time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at" bson:"updated_at"`
}

// NewGameSettings 创建新的游戏设置
func NewGameSettings() *GameSettings {
	now := time.Now()
	return &GameSettings{
		SoundEnabled:   true,
		MusicEnabled:   true,
		EffectsEnabled: true,
		Language:       "zh-CN",
		Theme:          "default",
		Quality:        GameQualityMedium,
		Controls:       make(map[string]interface{}),
		CustomSettings: make(map[string]interface{}),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// Clone 克隆游戏设置
func (gs *GameSettings) Clone() *GameSettings {
	clone := &GameSettings{
		SoundEnabled:   gs.SoundEnabled,
		MusicEnabled:   gs.MusicEnabled,
		EffectsEnabled: gs.EffectsEnabled,
		Language:       gs.Language,
		Theme:          gs.Theme,
		Quality:        gs.Quality,
		Controls:       make(map[string]interface{}),
		CustomSettings: make(map[string]interface{}),
		CreatedAt:      gs.CreatedAt,
		UpdatedAt:      gs.UpdatedAt,
	}

	// 深拷贝map
	for k, v := range gs.Controls {
		clone.Controls[k] = v
	}
	for k, v := range gs.CustomSettings {
		clone.CustomSettings[k] = v
	}

	return clone
}

// GameQuality 游戏质量
type GameQuality int32

const (
	GameQualityLow    GameQuality = iota + 1 // 低质量
	GameQualityMedium                        // 中等质量
	GameQualityHigh                          // 高质量
	GameQualityUltra                         // 超高质量
)

// String 返回游戏质量的字符串表示
func (gq GameQuality) String() string {
	switch gq {
	case GameQualityLow:
		return "low"
	case GameQualityMedium:
		return "medium"
	case GameQualityHigh:
		return "high"
	case GameQualityUltra:
		return "ultra"
	default:
		return "unknown"
	}
}

// IsValid 检查游戏质量是否有效
func (gq GameQuality) IsValid() bool {
	return gq >= GameQualityLow && gq <= GameQualityUltra
}

// 游戏查询相关值对象

// GameQuery 游戏查询条件
type GameQuery struct {
	GameID        *string        `json:"game_id,omitempty"`
	GameType      *GameType      `json:"game_type,omitempty"`
	Category      *GameCategory  `json:"category,omitempty"`
	Status        *GameStatus    `json:"status,omitempty"`
	CreatorID     *uint64        `json:"creator_id,omitempty"`
	PlayerID      *uint64        `json:"player_id,omitempty"`
	MinPlayers    *int32         `json:"min_players,omitempty"`
	MaxPlayers    *int32         `json:"max_players,omitempty"`
	MinDuration   *time.Duration `json:"min_duration,omitempty"`
	MaxDuration   *time.Duration `json:"max_duration,omitempty"`
	StartedAfter  *time.Time     `json:"started_after,omitempty"`
	StartedBefore *time.Time     `json:"started_before,omitempty"`
	EndedAfter    *time.Time     `json:"ended_after,omitempty"`
	EndedBefore   *time.Time     `json:"ended_before,omitempty"`
	CreatedAfter  *time.Time     `json:"created_after,omitempty"`
	CreatedBefore *time.Time     `json:"created_before,omitempty"`
	Keywords      []string       `json:"keywords,omitempty"`
	Tags          []string       `json:"tags,omitempty"`
	OrderBy       string         `json:"order_by,omitempty"`
	OrderDesc     bool           `json:"order_desc,omitempty"`
	Offset        int            `json:"offset,omitempty"`
	Limit         int            `json:"limit,omitempty"`
}

// GameFilter 游戏过滤器
type GameFilter struct {
	IncludeFinished  bool                   `json:"include_finished"`
	IncludeCancelled bool                   `json:"include_cancelled"`
	OnlyActive       bool                   `json:"only_active"`
	OnlyJoinable     bool                   `json:"only_joinable"`
	ExcludeGameIDs   []string               `json:"exclude_game_ids,omitempty"`
	ExcludePlayerIDs []uint64               `json:"exclude_player_ids,omitempty"`
	MinScore         *int64                 `json:"min_score,omitempty"`
	MaxScore         *int64                 `json:"max_score,omitempty"`
	Difficulties     []GameDifficulty       `json:"difficulties,omitempty"`
	CustomFilters    map[string]interface{} `json:"custom_filters,omitempty"`
}

// NewGameFilter 创建新的游戏过滤器
func NewGameFilter() *GameFilter {
	return &GameFilter{
		IncludeFinished:  false,
		IncludeCancelled: false,
		OnlyActive:       true,
		OnlyJoinable:     false,
		ExcludeGameIDs:   make([]string, 0),
		ExcludePlayerIDs: make([]uint64, 0),
		Difficulties:     make([]GameDifficulty, 0),
		CustomFilters:    make(map[string]interface{}),
	}
}

// 游戏操作相关值对象

// GameOperation 游戏操作类型
type GameOperation int32

const (
	GameOperationStart  GameOperation = iota + 1 // 开始游戏
	GameOperationPause                           // 暂停游戏
	GameOperationResume                          // 恢复游戏
	GameOperationEnd                             // 结束游戏
	GameOperationCancel                          // 取消游戏
	GameOperationReset                           // 重置游戏
	GameOperationKick                            // 踢出玩家
	GameOperationJoin                            // 加入游戏
	GameOperationLeave                           // 离开游戏
)

// String 返回游戏操作的字符串表示
func (go_ GameOperation) String() string {
	switch go_ {
	case GameOperationStart:
		return "start"
	case GameOperationPause:
		return "pause"
	case GameOperationResume:
		return "resume"
	case GameOperationEnd:
		return "end"
	case GameOperationCancel:
		return "cancel"
	case GameOperationReset:
		return "reset"
	case GameOperationKick:
		return "kick"
	case GameOperationJoin:
		return "join"
	case GameOperationLeave:
		return "leave"
	default:
		return "unknown"
	}
}

// IsValid 检查游戏操作是否有效
func (go_ GameOperation) IsValid() bool {
	return go_ >= GameOperationStart && go_ <= GameOperationLeave
}

// RequiresPermission 检查操作是否需要权限
func (go_ GameOperation) RequiresPermission() bool {
	switch go_ {
	case GameOperationStart, GameOperationPause, GameOperationResume,
		GameOperationEnd, GameOperationCancel, GameOperationReset, GameOperationKick:
		return true
	default:
		return false
	}
}

// GameOperationResult 游戏操作结果
type GameOperationResult struct {
	Success       bool                   `json:"success"`
	Operation     GameOperation          `json:"operation"`
	GameID        string                 `json:"game_id"`
	PlayerID      *uint64                `json:"player_id,omitempty"`
	OldStatus     *GameStatus            `json:"old_status,omitempty"`
	NewStatus     *GameStatus            `json:"new_status,omitempty"`
	AffectedCount int64                  `json:"affected_count"`
	Message       string                 `json:"message"`
	Error         string                 `json:"error,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	Duration      time.Duration          `json:"duration"`
}

// NewGameOperationResult 创建游戏操作结果
func NewGameOperationResult(operation GameOperation, gameID string, success bool) *GameOperationResult {
	return &GameOperationResult{
		Success:   success,
		Operation: operation,
		GameID:    gameID,
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
	}
}

// SetPlayerInfo 设置玩家信息
func (gor *GameOperationResult) SetPlayerInfo(playerID uint64) {
	gor.PlayerID = &playerID
}

// SetStatusChange 设置状态变更
func (gor *GameOperationResult) SetStatusChange(oldStatus, newStatus GameStatus) {
	gor.OldStatus = &oldStatus
	gor.NewStatus = &newStatus
}

// SetError 设置错误信息
func (gor *GameOperationResult) SetError(err error) {
	gor.Success = false
	gor.Error = err.Error()
}

// SetMessage 设置消息
func (gor *GameOperationResult) SetMessage(message string) {
	gor.Message = message
}

// SetDuration 设置持续时间
func (gor *GameOperationResult) SetDuration(start time.Time) {
	gor.Duration = time.Since(start)
}

// AddMetadata 添加元数据
func (gor *GameOperationResult) AddMetadata(key string, value interface{}) {
	if gor.Metadata == nil {
		gor.Metadata = make(map[string]interface{})
	}
	gor.Metadata[key] = value
}

// 验证函数

// ValidateGameQuery 验证游戏查询
func ValidateGameQuery(query *GameQuery) error {
	if query == nil {
		return fmt.Errorf("query cannot be nil")
	}

	if query.Limit <= 0 {
		return fmt.Errorf("limit must be positive")
	}

	if query.Limit > 1000 {
		return fmt.Errorf("limit cannot exceed 1000")
	}

	if query.Offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}

	if query.MinPlayers != nil && query.MaxPlayers != nil && *query.MinPlayers > *query.MaxPlayers {
		return fmt.Errorf("min_players cannot be greater than max_players")
	}

	if query.MinDuration != nil && query.MaxDuration != nil && *query.MinDuration > *query.MaxDuration {
		return fmt.Errorf("min_duration cannot be greater than max_duration")
	}

	if query.StartedAfter != nil && query.StartedBefore != nil && query.StartedAfter.After(*query.StartedBefore) {
		return fmt.Errorf("started_after cannot be after started_before")
	}

	if query.EndedAfter != nil && query.EndedBefore != nil && query.EndedAfter.After(*query.EndedBefore) {
		return fmt.Errorf("ended_after cannot be after ended_before")
	}

	if query.CreatedAfter != nil && query.CreatedBefore != nil && query.CreatedAfter.After(*query.CreatedBefore) {
		return fmt.Errorf("created_after cannot be after created_before")
	}

	return nil
}

// 辅助函数

// GetGameTypeByString 根据字符串获取游戏类型
func GetGameTypeByString(s string) (GameType, error) {
	switch s {
	case "save_dog":
		return GameTypeSaveDog, nil
	case "puzzle":
		return GameTypePuzzle, nil
	case "racing":
		return GameTypeRacing, nil
	case "memory":
		return GameTypeMemory, nil
	case "match":
		return GameTypeMatch, nil
	case "jump":
		return GameTypeJump, nil
	case "shoot":
		return GameTypeShoot, nil
	case "strategy":
		return GameTypeStrategy, nil
	case "card":
		return GameTypeCard, nil
	case "custom":
		return GameTypeCustom, nil
	default:
		return 0, fmt.Errorf("unknown game type: %s", s)
	}
}

// GetGameCategoryByString 根据字符串获取游戏分类
func GetGameCategoryByString(s string) (GameCategory, error) {
	switch s {
	case "normal":
		return GameCategoryNormal, nil
	case "competitive":
		return GameCategoryCompetitive, nil
	case "casual":
		return GameCategoryCasual, nil
	case "ranked":
		return GameCategoryRanked, nil
	default:
		return "", fmt.Errorf("unknown game category: %s", s)
	}
}

// GetGameStatusByString 根据字符串获取游戏状态
func GetGameStatusByString(s string) (GameStatus, error) {
	switch s {
	case "waiting":
		return GameStatusWaiting, nil
	case "running":
		return GameStatusRunning, nil
	case "paused":
		return GameStatusPaused, nil
	case "finished":
		return GameStatusFinished, nil
	case "cancelled":
		return GameStatusCancelled, nil
	case "error":
		return GameStatusError, nil
	default:
		return 0, fmt.Errorf("unknown game status: %s", s)
	}
}

// GetGameDifficultyByString 根据字符串获取游戏难度
func GetGameDifficultyByString(s string) (GameDifficulty, error) {
	switch s {
	case "easy":
		return GameDifficultyEasy, nil
	case "normal":
		return GameDifficultyNormal, nil
	case "hard":
		return GameDifficultyHard, nil
	case "expert":
		return GameDifficultyExpert, nil
	case "custom":
		return GameDifficultyCustom, nil
	default:
		return 0, fmt.Errorf("unknown game difficulty: %s", s)
	}
}
