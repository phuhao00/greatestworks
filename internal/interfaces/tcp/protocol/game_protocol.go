package protocol

import (
	"fmt"
	"greatestworks/internal/proto/messages"
	"time"
)

// 协议消息类型定义 - 使用proto生成的常量
const (
	// 玩家相关协议 (0x1000 - 0x1FFF) - 使用proto生成的常量
	MsgPlayerLogin      uint32 = uint32(messages.PlayerMessageID_MSG_PLAYER_LOGIN)
	MsgPlayerLogout     uint32 = uint32(messages.PlayerMessageID_MSG_PLAYER_LOGOUT)
	MsgPlayerMove       uint32 = uint32(messages.PlayerMessageID_MSG_PLAYER_MOVE)
	MsgPlayerInfo       uint32 = uint32(messages.PlayerMessageID_MSG_PLAYER_INFO)
	MsgPlayerCreate     uint32 = uint32(messages.PlayerMessageID_MSG_PLAYER_CREATE)
	MsgPlayerUpdate     uint32 = uint32(messages.PlayerMessageID_MSG_PLAYER_UPDATE)
	MsgPlayerDelete     uint32 = uint32(messages.PlayerMessageID_MSG_PLAYER_DELETE)
	MsgPlayerLevelUp    uint32 = uint32(messages.PlayerMessageID_MSG_PLAYER_LEVEL)
	MsgPlayerExpGain    uint32 = uint32(messages.PlayerMessageID_MSG_PLAYER_EXP_GAIN)
	MsgPlayerStatusSync uint32 = uint32(messages.PlayerMessageID_MSG_PLAYER_SYNC)
	MsgPlayerStatus     uint32 = uint32(messages.PlayerMessageID_MSG_PLAYER_STATUS)
	MsgPlayerStats      uint32 = uint32(messages.PlayerMessageID_MSG_PLAYER_STATS)

	// 战斗相关协议 (0x2000 - 0x2FFF) - 使用proto生成的常量
	MsgCreateBattle uint32 = uint32(messages.BattleMessageID_MSG_CREATE_BATTLE)
	MsgJoinBattle   uint32 = uint32(messages.BattleMessageID_MSG_JOIN_BATTLE)
	MsgStartBattle  uint32 = uint32(messages.BattleMessageID_MSG_START_BATTLE)
	MsgBattleAction uint32 = uint32(messages.BattleMessageID_MSG_BATTLE_ACTION)
	MsgLeaveBattle  uint32 = uint32(messages.BattleMessageID_MSG_LEAVE_BATTLE)
	MsgBattleResult uint32 = uint32(messages.BattleMessageID_MSG_BATTLE_RESULT)
	MsgBattleStatus uint32 = uint32(messages.BattleMessageID_MSG_BATTLE_STATUS)
	MsgBattleRound  uint32 = uint32(messages.BattleMessageID_MSG_BATTLE_ROUND)
	MsgBattleSkill  uint32 = uint32(messages.BattleMessageID_MSG_SKILL_CAST)
	MsgBattleDamage uint32 = uint32(messages.BattleMessageID_MSG_DAMAGE_DEALT)

	// 查询相关协议 (0x3000 - 0x3FFF) - 使用proto生成的常量
	MsgGetPlayerInfo    uint32 = uint32(messages.QueryMessageID_MSG_GET_PLAYER_INFO)
	MsgGetOnlinePlayers uint32 = uint32(messages.QueryMessageID_MSG_GET_ONLINE_PLAYERS)
	MsgGetBattleInfo    uint32 = uint32(messages.QueryMessageID_MSG_GET_BATTLE_INFO)
	MsgGetPlayerStats   uint32 = uint32(messages.QueryMessageID_MSG_GET_PLAYER_INFO) // 使用玩家信息代替
	MsgGetBattleList    uint32 = uint32(messages.QueryMessageID_MSG_GET_BATTLE_INFO) // 使用战斗信息代替
	MsgGetRankings      uint32 = uint32(messages.QueryMessageID_MSG_GET_RANKING_LIST)
	MsgGetServerInfo    uint32 = uint32(messages.QueryMessageID_MSG_GET_SERVER_INFO)

	// 系统相关协议 (0x9000 - 0x9FFF)
	// 使用message_types.go中定义的消息类型
)

// 基础协议结构

// BaseRequest 基础请求结构
type BaseRequest struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

// BaseResponse 基础响应结构
type BaseResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	BaseResponse
	ErrorCode int    `json:"error_code,omitempty"`
	ErrorType string `json:"error_type,omitempty"`
}

// 玩家协议结构

// PlayerLoginRequest 玩家登录请求
type PlayerLoginRequest struct {
	BaseRequest
	PlayerID string `json:"player_id"`
	Token    string `json:"token"`
	Version  string `json:"version,omitempty"`
}

// PlayerLoginResponse 玩家登录响应
type PlayerLoginResponse struct {
	BaseResponse
	Player      *PlayerInfo `json:"player,omitempty"`
	SessionID   string      `json:"session_id,omitempty"`
	ServerTime  int64       `json:"server_time,omitempty"`
	Permissions []string    `json:"permissions,omitempty"`
}

// PlayerCreateRequest 创建玩家请求
type PlayerCreateRequest struct {
	BaseRequest
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
	Gender int    `json:"gender,omitempty"`
}

// PlayerCreateResponse 创建玩家响应
type PlayerCreateResponse struct {
	BaseResponse
	PlayerID  string    `json:"player_id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Level     int       `json:"level,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// PlayerMoveRequest 玩家移动请求
type PlayerMoveRequest struct {
	BaseRequest
	Position Position `json:"position"`
	Speed    float64  `json:"speed,omitempty"`
}

// PlayerMoveResponse 玩家移动响应
type PlayerMoveResponse struct {
	BaseResponse
	OldPosition Position `json:"old_position,omitempty"`
	NewPosition Position `json:"new_position,omitempty"`
	MoveTime    int64    `json:"move_time,omitempty"`
}

// PlayerInfoRequest 获取玩家信息请求
type PlayerInfoRequest struct {
	BaseRequest
	TargetPlayerID string `json:"target_player_id,omitempty"`
}

// PlayerInfoResponse 获取玩家信息响应
type PlayerInfoResponse struct {
	BaseResponse
	Player *PlayerInfo `json:"player,omitempty"`
}

// 战斗协议结构

// CreateBattleRequest 创建战斗请求
type CreateBattleRequest struct {
	BaseRequest
	BattleType int               `json:"battle_type"`
	Settings   *BattleSettings   `json:"settings,omitempty"`
	Players    []string          `json:"players,omitempty"`
	Options    map[string]string `json:"options,omitempty"`
}

// CreateBattleResponse 创建战斗响应
type CreateBattleResponse struct {
	BaseResponse
	BattleID   string          `json:"battle_id,omitempty"`
	BattleType int             `json:"battle_type,omitempty"`
	Status     string          `json:"status,omitempty"`
	Settings   *BattleSettings `json:"settings,omitempty"`
	CreatedAt  time.Time       `json:"created_at,omitempty"`
}

// JoinBattleRequest 加入战斗请求
type JoinBattleRequest struct {
	BaseRequest
	BattleID string `json:"battle_id"`
	Team     int    `json:"team,omitempty"`
	Position int    `json:"position,omitempty"`
}

// JoinBattleResponse 加入战斗响应
type JoinBattleResponse struct {
	BaseResponse
	BattleID     string        `json:"battle_id,omitempty"`
	PlayerTeam   int           `json:"player_team,omitempty"`
	PlayerPos    int           `json:"player_position,omitempty"`
	BattleInfo   *BattleInfo   `json:"battle_info,omitempty"`
	OtherPlayers []*PlayerInfo `json:"other_players,omitempty"`
}

// BattleActionRequest 战斗行动请求
type BattleActionRequest struct {
	BaseRequest
	BattleID   string                 `json:"battle_id"`
	ActionType string                 `json:"action_type"`
	TargetID   string                 `json:"target_id,omitempty"`
	SkillID    string                 `json:"skill_id,omitempty"`
	Params     map[string]interface{} `json:"params,omitempty"`
}

// BattleActionResponse 战斗行动响应
type BattleActionResponse struct {
	BaseResponse
	BattleID     string        `json:"battle_id,omitempty"`
	ActionResult *ActionResult `json:"action_result,omitempty"`
	BattleState  *BattleState  `json:"battle_state,omitempty"`
	NextTurn     string        `json:"next_turn,omitempty"`
}

// 查询协议结构

// GetOnlinePlayersRequest 获取在线玩家请求
type GetOnlinePlayersRequest struct {
	BaseRequest
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}

// GetOnlinePlayersResponse 获取在线玩家响应
type GetOnlinePlayersResponse struct {
	BaseResponse
	Players    []*PlayerInfo `json:"players,omitempty"`
	Total      int           `json:"total,omitempty"`
	Page       int           `json:"page,omitempty"`
	PageSize   int           `json:"page_size,omitempty"`
	ServerTime int64         `json:"server_time,omitempty"`
}

// GetBattleInfoRequest 获取战斗信息请求
type GetBattleInfoRequest struct {
	BaseRequest
	BattleID string `json:"battle_id"`
}

// GetBattleInfoResponse 获取战斗信息响应
type GetBattleInfoResponse struct {
	BaseResponse
	BattleInfo *BattleInfo `json:"battle_info,omitempty"`
}

// 数据结构定义

// PlayerInfo 玩家信息
type PlayerInfo struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Level      int       `json:"level"`
	Exp        int64     `json:"exp"`
	Status     string    `json:"status"`
	Position   Position  `json:"position"`
	Stats      Stats     `json:"stats"`
	Avatar     string    `json:"avatar,omitempty"`
	Gender     int       `json:"gender,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	LastActive time.Time `json:"last_active,omitempty"`
}

// Position 位置信息
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Stats 属性信息
type Stats struct {
	HP      int `json:"hp"`
	MaxHP   int `json:"max_hp"`
	MP      int `json:"mp"`
	MaxMP   int `json:"max_mp"`
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	Speed   int `json:"speed"`
	Crit    int `json:"crit,omitempty"`
	Hit     int `json:"hit,omitempty"`
	Dodge   int `json:"dodge,omitempty"`
}

// BattleInfo 战斗信息
type BattleInfo struct {
	ID        string          `json:"id"`
	Type      int             `json:"type"`
	Status    string          `json:"status"`
	Players   []*PlayerInfo   `json:"players"`
	Settings  *BattleSettings `json:"settings"`
	State     *BattleState    `json:"state,omitempty"`
	StartTime *time.Time      `json:"start_time,omitempty"`
	EndTime   *time.Time      `json:"end_time,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// BattleSettings 战斗设置
type BattleSettings struct {
	MaxPlayers    int           `json:"max_players"`
	TurnTimeout   time.Duration `json:"turn_timeout"`
	MaxRounds     int           `json:"max_rounds"`
	AutoStart     bool          `json:"auto_start"`
	AllowSpectate bool          `json:"allow_spectate"`
	TeamMode      bool          `json:"team_mode"`
}

// BattleState 战斗状态
type BattleState struct {
	CurrentRound int                     `json:"current_round"`
	CurrentTurn  string                  `json:"current_turn"`
	TurnOrder    []string                `json:"turn_order"`
	PlayerStates map[string]*PlayerState `json:"player_states"`
	RoundHistory []*RoundResult          `json:"round_history,omitempty"`
}

// PlayerState 玩家战斗状态
type PlayerState struct {
	PlayerID   string         `json:"player_id"`
	HP         int            `json:"hp"`
	MP         int            `json:"mp"`
	Status     []string       `json:"status,omitempty"`
	Buffs      []*Buff        `json:"buffs,omitempty"`
	Debuffs    []*Debuff      `json:"debuffs,omitempty"`
	SkillCDs   map[string]int `json:"skill_cds,omitempty"`
	Position   int            `json:"position"`
	Team       int            `json:"team"`
	IsAlive    bool           `json:"is_alive"`
	LastAction *ActionResult  `json:"last_action,omitempty"`
}

// ActionResult 行动结果
type ActionResult struct {
	ActionID   string                 `json:"action_id"`
	PlayerID   string                 `json:"player_id"`
	ActionType string                 `json:"action_type"`
	TargetID   string                 `json:"target_id,omitempty"`
	SkillID    string                 `json:"skill_id,omitempty"`
	Damage     int                    `json:"damage,omitempty"`
	Healing    int                    `json:"healing,omitempty"`
	Effects    []*Effect              `json:"effects,omitempty"`
	Success    bool                   `json:"success"`
	Message    string                 `json:"message,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	Params     map[string]interface{} `json:"params,omitempty"`
}

// RoundResult 回合结果
type RoundResult struct {
	RoundNumber int             `json:"round_number"`
	Actions     []*ActionResult `json:"actions"`
	StartTime   time.Time       `json:"start_time"`
	EndTime     time.Time       `json:"end_time"`
	Winner      string          `json:"winner,omitempty"`
}

// Buff 增益效果
type Buff struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Duration    int       `json:"duration"`
	Effect      *Effect   `json:"effect"`
	StartTime   time.Time `json:"start_time"`
}

// Debuff 减益效果
type Debuff struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Duration    int       `json:"duration"`
	Effect      *Effect   `json:"effect"`
	StartTime   time.Time `json:"start_time"`
}

// Effect 效果
type Effect struct {
	Type     string                 `json:"type"`
	Value    int                    `json:"value,omitempty"`
	Target   string                 `json:"target,omitempty"`
	Duration int                    `json:"duration,omitempty"`
	Params   map[string]interface{} `json:"params,omitempty"`
	Trigger  string                 `json:"trigger,omitempty"`
}

// 系统协议结构

// HeartbeatRequest 心跳请求
type HeartbeatRequest struct {
	BaseRequest
	ClientTime int64 `json:"client_time"`
}

// HeartbeatResponse 心跳响应
type HeartbeatResponse struct {
	BaseResponse
	ServerTime int64 `json:"server_time"`
	Latency    int64 `json:"latency,omitempty"`
}

// PingRequest Ping请求
type PingRequest struct {
	BaseRequest
	Sequence   int   `json:"sequence"`
	ClientTime int64 `json:"client_time"`
}

// PingResponse Ping响应
type PingResponse struct {
	BaseResponse
	Sequence   int   `json:"sequence"`
	ServerTime int64 `json:"server_time"`
	RoundTrip  int64 `json:"round_trip,omitempty"`
}

// 协议工具函数

// NewBaseRequest 创建基础请求
func NewBaseRequest() BaseRequest {
	return BaseRequest{
		Timestamp: time.Now().Unix(),
	}
}

// NewBaseResponse 创建基础响应
func NewBaseResponse(success bool, message string) BaseResponse {
	return BaseResponse{
		Success:   success,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(message string, errorCode int, errorType string) ErrorResponse {
	return ErrorResponse{
		BaseResponse: NewBaseResponse(false, message),
		ErrorCode:    errorCode,
		ErrorType:    errorType,
	}
}

// IsValidMessageType 检查消息类型是否有效
func IsValidMessageType(msgType uint32) bool {
	switch {
	case msgType >= 0x1000 && msgType <= 0x1FFF: // 玩家协议
		return true
	case msgType >= 0x2000 && msgType <= 0x2FFF: // 战斗协议
		return true
	case msgType >= 0x3000 && msgType <= 0x3FFF: // 查询协议
		return true
	case msgType >= 0x9000 && msgType <= 0x9FFF: // 系统协议
		return true
	default:
		return false
	}
}

// GetMessageTypeName 获取消息类型名称
func GetMessageTypeName(msgType uint32) string {
	msgNames := map[uint32]string{
		MsgPlayerLogin:       "PlayerLogin",
		MsgPlayerLogout:      "PlayerLogout",
		MsgPlayerMove:        "PlayerMove",
		MsgPlayerInfo:        "PlayerInfo",
		MsgPlayerCreate:      "PlayerCreate",
		MsgCreateBattle:      "CreateBattle",
		MsgJoinBattle:        "JoinBattle",
		MsgStartBattle:       "StartBattle",
		MsgBattleAction:      "BattleAction",
		MsgLeaveBattle:       "LeaveBattle",
		MsgGetPlayerInfo:     "GetPlayerInfo",
		MsgGetOnlinePlayers:  "GetOnlinePlayers",
		MsgGetBattleInfo:     "GetBattleInfo",
		uint32(MsgHeartbeat): "Heartbeat",
		uint32(MsgPing):      "Ping",
		uint32(MsgError):     "Error",
	}

	if name, exists := msgNames[msgType]; exists {
		return name
	}
	return fmt.Sprintf("Unknown(0x%04X)", msgType)
}
