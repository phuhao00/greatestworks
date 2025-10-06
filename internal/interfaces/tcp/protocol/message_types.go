package protocol

// import (
// 	"time"
// )

// 消息类型常量定义
const (
	// 系统消息 (0x0000 - 0x00FF)
	MsgHeartbeat  uint16 = 0x0001 // 心跳消息
	MsgHandshake  uint16 = 0x0002 // 握手消息
	MsgAuth       uint16 = 0x0003 // 认证消息
	MsgDisconnect uint16 = 0x0004 // 断开连接
	MsgError      uint16 = 0x0005 // 错误消息
	MsgPing       uint16 = 0x0006 // Ping消息
	MsgPong       uint16 = 0x0007 // Pong消息

	// 玩家相关消息 (0x0100 - 0x01FF) - 定义在game_protocol.go中
	// MsgPlayerLogin   uint16 = 0x0101 // 玩家登录
	// MsgPlayerLogout  uint16 = 0x0102 // 玩家登出
	// MsgPlayerInfo    uint16 = 0x0103 // 玩家信息
	// MsgPlayerMove    uint16 = 0x0104 // 玩家移动
	// MsgPlayerCreate  uint16 = 0x0105 // 创建玩家
	// MsgPlayerUpdate  uint16 = 0x0106 // 更新玩家
	// MsgPlayerDelete  uint16 = 0x0107 // 删除玩家
	MsgPlayerStatus uint16 = 0x0108 // 玩家状态
	MsgPlayerStats  uint16 = 0x0109 // 玩家属性
	MsgPlayerLevel  uint16 = 0x010A // 玩家升级

	// 战斗相关消息 (0x0200 - 0x02FF) - 定义在game_protocol.go中
	// MsgCreateBattle  uint16 = 0x0201 // 创建战斗
	// MsgJoinBattle    uint16 = 0x0202 // 加入战斗
	// MsgLeaveBattle   uint16 = 0x0203 // 离开战斗
	// MsgStartBattle   uint16 = 0x0204 // 开始战斗 - 定义在game_protocol.go中
	MsgEndBattle uint16 = 0x0205 // 结束战斗
	// MsgBattleAction  uint16 = 0x0206 // 战斗行动 - 定义在game_protocol.go中
	// MsgBattleResult  uint16 = 0x0207 // 战斗结果 - 定义在game_protocol.go中
	// MsgBattleStatus  uint16 = 0x0208 // 战斗状态 - 定义在game_protocol.go中
	MsgSkillCast   uint16 = 0x0209 // 技能释放
	MsgDamageDealt uint16 = 0x020A // 伤害计算

	// 宠物相关消息 (0x0300 - 0x03FF)
	MsgPetSummon    uint16 = 0x0301 // 召唤宠物
	MsgPetDismiss   uint16 = 0x0302 // 收回宠物
	MsgPetInfo      uint16 = 0x0303 // 宠物信息
	MsgPetMove      uint16 = 0x0304 // 宠物移动
	MsgPetAction    uint16 = 0x0305 // 宠物行动
	MsgPetLevelUp   uint16 = 0x0306 // 宠物升级
	MsgPetEvolution uint16 = 0x0307 // 宠物进化
	MsgPetTrain     uint16 = 0x0308 // 宠物训练
	MsgPetFeed      uint16 = 0x0309 // 宠物喂养
	MsgPetStatus    uint16 = 0x030A // 宠物状态

	// 建筑相关消息 (0x0400 - 0x04FF)
	MsgBuildingCreate  uint16 = 0x0401 // 创建建筑
	MsgBuildingUpgrade uint16 = 0x0402 // 升级建筑
	MsgBuildingDestroy uint16 = 0x0403 // 摧毁建筑
	MsgBuildingInfo    uint16 = 0x0404 // 建筑信息
	MsgBuildingProduce uint16 = 0x0405 // 建筑生产
	MsgBuildingCollect uint16 = 0x0406 // 收集资源
	MsgBuildingRepair  uint16 = 0x0407 // 修复建筑
	MsgBuildingStatus  uint16 = 0x0408 // 建筑状态

	// 社交相关消息 (0x0500 - 0x05FF)
	MsgChatMessage   uint16 = 0x0501 // 聊天消息
	MsgFriendRequest uint16 = 0x0502 // 好友请求
	MsgFriendAccept  uint16 = 0x0503 // 接受好友
	MsgFriendReject  uint16 = 0x0504 // 拒绝好友
	MsgFriendRemove  uint16 = 0x0505 // 删除好友
	MsgFriendList    uint16 = 0x0506 // 好友列表
	MsgGuildCreate   uint16 = 0x0507 // 创建公会
	MsgGuildJoin     uint16 = 0x0508 // 加入公会
	MsgGuildLeave    uint16 = 0x0509 // 离开公会
	MsgGuildInfo     uint16 = 0x050A // 公会信息
	MsgTeamCreate    uint16 = 0x050B // 创建队伍
	MsgTeamJoin      uint16 = 0x050C // 加入队伍
	MsgTeamLeave     uint16 = 0x050D // 离开队伍
	MsgTeamInfo      uint16 = 0x050E // 队伍信息

	// 物品相关消息 (0x0600 - 0x06FF)
	MsgItemUse       uint16 = 0x0601 // 使用物品
	MsgItemEquip     uint16 = 0x0602 // 装备物品
	MsgItemUnequip   uint16 = 0x0603 // 卸下装备
	MsgItemDrop      uint16 = 0x0604 // 丢弃物品
	MsgItemPickup    uint16 = 0x0605 // 拾取物品
	MsgItemTrade     uint16 = 0x0606 // 交易物品
	MsgInventoryInfo uint16 = 0x0607 // 背包信息
	MsgItemInfo      uint16 = 0x0608 // 物品信息
	MsgItemCraft     uint16 = 0x0609 // 制作物品
	MsgItemEnhance   uint16 = 0x060A // 强化物品

	// 任务相关消息 (0x0700 - 0x07FF)
	MsgQuestAccept   uint16 = 0x0701 // 接受任务
	MsgQuestComplete uint16 = 0x0702 // 完成任务
	MsgQuestCancel   uint16 = 0x0703 // 取消任务
	MsgQuestProgress uint16 = 0x0704 // 任务进度
	MsgQuestList     uint16 = 0x0705 // 任务列表
	MsgQuestInfo     uint16 = 0x0706 // 任务信息
	MsgQuestReward   uint16 = 0x0707 // 任务奖励

	// 查询相关消息 (0x0800 - 0x08FF)
	// MsgGetPlayerInfo    uint16 = 0x0801 // 获取玩家信息 - 定义在game_protocol.go中
	// MsgGetOnlinePlayers uint16 = 0x0802 // 获取在线玩家 - 定义在game_protocol.go中
	// MsgGetBattleInfo    uint16 = 0x0803 // 获取战斗信息 - 定义在game_protocol.go中
	MsgGetRankingList uint16 = 0x0804 // 获取排行榜
	MsgGetServerInfo  uint16 = 0x0805 // 获取服务器信息
)

// 消息魔数
const MessageMagic uint32 = 0x47574B53 // "GWKS" - GreatestWorks

// MessageHeader TCP消息头
type MessageHeader struct {
	Magic       uint32 `json:"magic"`        // 魔数标识
	MessageID   uint32 `json:"message_id"`   // 消息ID（用于请求响应匹配）
	MessageType uint16 `json:"message_type"` // 消息类型
	Flags       uint16 `json:"flags"`        // 标志位
	PlayerID    uint64 `json:"player_id"`    // 玩家ID
	Timestamp   int64  `json:"timestamp"`    // 时间戳
	Sequence    uint32 `json:"sequence"`     // 序列号
	Length      uint32 `json:"length"`       // 消息体长度
}

// Message TCP消息
type Message struct {
	Header  MessageHeader `json:"header"`
	Payload interface{}   `json:"payload"`
}

// 消息标志位
const (
	FlagRequest    uint16 = 0x0001 // 请求消息
	FlagResponse   uint16 = 0x0002 // 响应消息
	FlagError      uint16 = 0x0004 // 错误消息
	FlagAsync      uint16 = 0x0008 // 异步消息
	FlagBroadcast  uint16 = 0x0010 // 广播消息
	FlagEncrypted  uint16 = 0x0020 // 加密消息
	FlagCompressed uint16 = 0x0040 // 压缩消息
)

// BaseResponse 基础响应 - 定义在game_protocol.go中
// type BaseResponse struct {
// 	Success   bool   `json:"success"`
// 	Message   string `json:"message,omitempty"`
// 	ErrorCode int    `json:"error_code,omitempty"`
// }

// NewBaseResponse 创建基础响应 - 定义在game_protocol.go中
// func NewBaseResponse(success bool, message string) BaseResponse {
// 	return BaseResponse{
// 		Success: success,
// 		Message: message,
// 	}
// }

// Position 位置信息 - 定义在game_protocol.go中
// type Position struct {
// 	X float64 `json:"x"`
// 	Y float64 `json:"y"`
// 	Z float64 `json:"z"`
// }

// Stats 属性信息 - 定义在game_protocol.go中
// type Stats struct {
// 	HP      int `json:"hp"`
// 	MaxHP   int `json:"max_hp"`
// 	MP      int `json:"mp"`
// 	MaxMP   int `json:"max_mp"`
// 	Attack  int `json:"attack"`
// 	Defense int `json:"defense"`
// 	Speed   int `json:"speed"`
// }

// PlayerInfo 玩家信息 - 定义在game_protocol.go中
// type PlayerInfo struct {
// 	ID        string    `json:"id"`
// 	Name      string    `json:"name"`
// 	Level     int       `json:"level"`
// 	Exp       int64     `json:"exp"`
// 	Status    string    `json:"status"`
// 	Position  Position  `json:"position"`
// 	Stats     Stats     `json:"stats"`
// 	Avatar    string    `json:"avatar,omitempty"`
// 	Gender    int       `json:"gender,omitempty"`
// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at"`
// }

// 系统消息

// HeartbeatRequest 心跳请求 - 定义在game_protocol.go中
// type HeartbeatRequest struct {
// 	Timestamp int64 `json:"timestamp"`
// }

// HeartbeatResponse 心跳响应 - 定义在game_protocol.go中
// type HeartbeatResponse struct {
// 	BaseResponse
// 	ServerTime int64 `json:"server_time"`
// }

// AuthRequest 认证请求
type AuthRequest struct {
	Token      string     `json:"token"`
	PlayerID   string     `json:"player_id"`
	ClientInfo ClientInfo `json:"client_info"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	BaseResponse
	SessionID  string      `json:"session_id,omitempty"`
	PlayerInfo *PlayerInfo `json:"player_info,omitempty"`
	ServerTime int64       `json:"server_time"`
}

// ClientInfo 客户端信息
type ClientInfo struct {
	Version   string `json:"version"`
	Platform  string `json:"platform"`
	DeviceID  string `json:"device_id"`
	IPAddress string `json:"ip_address,omitempty"`
}

// 玩家相关消息

// PlayerLoginRequest 玩家登录请求 - 定义在game_protocol.go中
// type PlayerLoginRequest struct {
// 	PlayerID string `json:"player_id"`
// 	Token    string `json:"token"`
// }

// PlayerLoginResponse 玩家登录响应 - 定义在game_protocol.go中
// type PlayerLoginResponse struct {
// 	BaseResponse
// 	Player     *PlayerInfo `json:"player,omitempty"`
// 	SessionID  string      `json:"session_id,omitempty"`
// 	ServerTime int64       `json:"server_time"`
// }

// PlayerMoveRequest 玩家移动请求 - 定义在game_protocol.go中
// type PlayerMoveRequest struct {
// 	Position Position `json:"position"`
// 	Speed    float64  `json:"speed,omitempty"`
// }

// PlayerMoveResponse 玩家移动响应 - 定义在game_protocol.go中
// type PlayerMoveResponse struct {
// 	BaseResponse
// 	OldPosition Position `json:"old_position"`
// 	NewPosition Position `json:"new_position"`
// 	MoveTime    int64    `json:"move_time"`
// }

// PlayerCreateRequest 创建玩家请求 - 定义在game_protocol.go中
// type PlayerCreateRequest struct {
// 	Name   string `json:"name"`
// 	Avatar string `json:"avatar,omitempty"`
// 	Gender int    `json:"gender,omitempty"`
// }

// PlayerCreateResponse 创建玩家响应 - 定义在game_protocol.go中
// type PlayerCreateResponse struct {
// 	BaseResponse
// 	PlayerID  string    `json:"player_id,omitempty"`
// 	Name      string    `json:"name,omitempty"`
// 	Level     int       `json:"level,omitempty"`
// 	CreatedAt time.Time `json:"created_at,omitempty"`
// }

// 战斗相关消息

// CreateBattleRequest 创建战斗请求 - 定义在game_protocol.go中
// type CreateBattleRequest struct {
// 	BattleType string   `json:"battle_type"`
// 	MaxPlayers int      `json:"max_players"`
// 	Settings   BattleSettings `json:"settings,omitempty"`
// }

// CreateBattleResponse 创建战斗响应 - 定义在game_protocol.go中
// type CreateBattleResponse struct {
// 	BaseResponse
// 	BattleID   string `json:"battle_id,omitempty"`
// 	BattleInfo *BattleInfo `json:"battle_info,omitempty"`
// }

// BattleSettings 战斗设置 - 定义在game_protocol.go中
// type BattleSettings struct {
// 	TimeLimit    int  `json:"time_limit,omitempty"`
// 	AllowPets    bool `json:"allow_pets"`
// 	AllowItems   bool `json:"allow_items"`
// 	FriendlyFire bool `json:"friendly_fire"`
// }

// BattleInfo 战斗信息 - 定义在game_protocol.go中
// type BattleInfo struct {
// 	ID         string        `json:"id"`
// 	Type       string        `json:"type"`
// 	Status     string        `json:"status"`
// 	Players    []PlayerInfo  `json:"players"`
// 	Settings   BattleSettings `json:"settings"`
// 	StartTime  *time.Time    `json:"start_time,omitempty"`
// 	EndTime    *time.Time    `json:"end_time,omitempty"`
// 	CreatedAt  time.Time     `json:"created_at"`
// }

// JoinBattleRequest 加入战斗请求 - 定义在game_protocol.go中
// type JoinBattleRequest struct {
// 	BattleID string `json:"battle_id"`
// 	Team     int    `json:"team,omitempty"`
// }

// JoinBattleResponse 加入战斗响应 - 定义在game_protocol.go中
// type JoinBattleResponse struct {
// 	BaseResponse
// 	BattleInfo *BattleInfo `json:"battle_info,omitempty"`
// 	PlayerTeam int         `json:"player_team,omitempty"`
// }

// BattleActionRequest 战斗行动请求 - 定义在game_protocol.go中
// type BattleActionRequest struct {
// 	BattleID   string      `json:"battle_id"`
// 	ActionType string      `json:"action_type"`
// 	TargetID   string      `json:"target_id,omitempty"`
// 	SkillID    string      `json:"skill_id,omitempty"`
// 	ItemID     string      `json:"item_id,omitempty"`
// 	Position   *Position   `json:"position,omitempty"`
// 	Params     interface{} `json:"params,omitempty"`
// }

// BattleActionResponse 战斗行动响应 - 定义在game_protocol.go中
// type BattleActionResponse struct {
// 	BaseResponse
// 	ActionResult *ActionResult `json:"action_result,omitempty"`
// 	BattleState  *BattleState  `json:"battle_state,omitempty"`
// }

// ActionResult 行动结果 - 定义在game_protocol.go中
// type ActionResult struct {
// 	ActionID   string      `json:"action_id"`
// 	PlayerID   string      `json:"player_id"`
// 	ActionType string      `json:"action_type"`
// 	TargetID   string      `json:"target_id,omitempty"`
// 	Damage     int         `json:"damage,omitempty"`
// 	Healing    int         `json:"healing,omitempty"`
// 	Effects    []Effect    `json:"effects,omitempty"`
// 	Success    bool        `json:"success"`
// 	Message    string      `json:"message,omitempty"`
// 	Timestamp  time.Time   `json:"timestamp"`
// }

// Effect 效果 - 定义在game_protocol.go中
// type Effect struct {
// 	Type     string      `json:"type"`
// 	Value    int         `json:"value"`
// 	Duration int         `json:"duration,omitempty"`
// 	Params   interface{} `json:"params,omitempty"`
// }

// BattleState 战斗状态 - 定义在game_protocol.go中
// type BattleState struct {
// 	BattleID    string              `json:"battle_id"`
// 	Status      string              `json:"status"`
// 	CurrentTurn string              `json:"current_turn,omitempty"`
// 	TurnNumber  int                 `json:"turn_number"`
// 	Players     map[string]*PlayerState `json:"players"`
// 	TimeLeft    int                 `json:"time_left,omitempty"`
// 	UpdatedAt   time.Time           `json:"updated_at"`
// }

// PlayerState 玩家状态 - 定义在game_protocol.go中
// type PlayerState struct {
// 	PlayerID  string    `json:"player_id"`
// 	HP        int       `json:"hp"`
// 	MaxHP     int       `json:"max_hp"`
// 	MP        int       `json:"mp"`
// 	MaxMP     int       `json:"max_mp"`
// 	Position  Position  `json:"position"`
// 	Status    string    `json:"status"`
// 	Effects   []Effect  `json:"effects,omitempty"`
// 	IsAlive   bool      `json:"is_alive"`
// 	Team      int       `json:"team"`
// 	UpdatedAt time.Time `json:"updated_at"`
// }

// 以下所有结构体都定义在game_protocol.go中，这里注释掉避免重复定义
/*

// 宠物相关消息

// PetSummonRequest 召唤宠物请求
type PetSummonRequest struct {
	PetID string `json:"pet_id"`
	Slot  int    `json:"slot,omitempty"`
}

// PetSummonResponse 召唤宠物响应
type PetSummonResponse struct {
	BaseResponse
	PetInfo *PetInfo `json:"pet_info,omitempty"`
}

// PetInfo 宠物信息
type PetInfo struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Level     int       `json:"level"`
	Exp       int64     `json:"exp"`
	Stats     Stats     `json:"stats"`
	Skills    []string  `json:"skills,omitempty"`
	Status    string    `json:"status"`
	Position  Position  `json:"position"`
	OwnerID   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 聊天相关消息

// ChatMessageRequest 聊天消息请求
type ChatMessageRequest struct {
	Channel   string `json:"channel"`
	Content   string `json:"content"`
	TargetID  string `json:"target_id,omitempty"`
	MessageType string `json:"message_type,omitempty"`
}

// ChatMessageResponse 聊天消息响应
type ChatMessageResponse struct {
	BaseResponse
	MessageID string    `json:"message_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	ID        string    `json:"id"`
	Channel   string    `json:"channel"`
	SenderID  string    `json:"sender_id"`
	SenderName string   `json:"sender_name"`
	Content   string    `json:"content"`
	TargetID  string    `json:"target_id,omitempty"`
	MessageType string  `json:"message_type"`
	Timestamp time.Time `json:"timestamp"`
}

// 错误消息

// ErrorMessage 错误消息
type ErrorMessage struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// 错误码定义
const (
	ErrSuccess           = 0     // 成功
	ErrUnknown          = 1000  // 未知错误
	ErrInvalidMessage   = 1001  // 无效消息
	ErrInvalidPlayer    = 1002  // 无效玩家
	ErrPlayerNotFound   = 1003  // 玩家未找到
	ErrPlayerOffline    = 1004  // 玩家离线
	ErrAuthFailed       = 1005  // 认证失败
	ErrPermissionDenied = 1006  // 权限不足
	ErrRateLimited      = 1007  // 请求过于频繁
	ErrServerBusy       = 1008  // 服务器繁忙
	ErrMaintenance      = 1009  // 服务器维护

	// 战斗相关错误
	ErrBattleNotFound   = 2001  // 战斗未找到
	ErrBattleFull       = 2002  // 战斗已满
	ErrBattleStarted    = 2003  // 战斗已开始
	ErrBattleEnded      = 2004  // 战斗已结束
	ErrInvalidAction    = 2005  // 无效行动
	ErrNotYourTurn      = 2006  // 不是你的回合
	ErrSkillCooldown    = 2007  // 技能冷却中
	ErrInsufficientMP   = 2008  // MP不足

	// 宠物相关错误
	ErrPetNotFound      = 3001  // 宠物未找到
	ErrPetAlreadyActive = 3002  // 宠物已激活
	ErrPetNotActive     = 3003  // 宠物未激活
	ErrPetLevelTooLow   = 3004  // 宠物等级过低
	ErrPetEvolutionFail = 3005  // 宠物进化失败

	// 物品相关错误
	ErrItemNotFound     = 4001  // 物品未找到
	ErrItemNotUsable    = 4002  // 物品不可使用
	ErrInventoryFull    = 4003  // 背包已满
	ErrInsufficientItem = 4004  // 物品数量不足
	ErrItemEquipFailed  = 4005  // 装备失败
)
*/
