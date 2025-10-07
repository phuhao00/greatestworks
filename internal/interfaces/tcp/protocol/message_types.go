package protocol

import (
	"fmt"
	"greatestworks/internal/proto/messages"
	"greatestworks/internal/proto/protocol"
)

// import (
// 	"time"
// )

// 消息类型常量定义 - 使用proto生成的常量
const (
	// 系统消息 (0x0000 - 0x00FF)
	MsgHeartbeat  uint32 = uint32(messages.SystemMessageID_MSG_HEARTBEAT)
	MsgHandshake  uint32 = uint32(messages.SystemMessageID_MSG_HANDSHAKE)
	MsgAuth       uint32 = uint32(messages.SystemMessageID_MSG_AUTH)
	MsgDisconnect uint32 = uint32(messages.SystemMessageID_MSG_DISCONNECT)
	MsgError      uint32 = uint32(messages.SystemMessageID_MSG_ERROR)
	MsgPing       uint32 = uint32(messages.SystemMessageID_MSG_PING)
	MsgPong       uint32 = uint32(messages.SystemMessageID_MSG_PONG)

	// 玩家相关消息 (0x0100 - 0x01FF) - 定义在game_protocol.go中
	// 战斗相关消息 (0x0200 - 0x02FF) - 定义在game_protocol.go中

	// 宠物相关消息 (0x0300 - 0x03FF) - 使用proto生成的常量
	MsgPetSummon    uint32 = uint32(messages.PetMessageID_MSG_PET_SUMMON)    // 召唤宠物
	MsgPetDismiss   uint32 = uint32(messages.PetMessageID_MSG_PET_DISMISS)   // 收回宠物
	MsgPetInfo      uint32 = uint32(messages.PetMessageID_MSG_PET_INFO)      // 宠物信息
	MsgPetMove      uint32 = uint32(messages.PetMessageID_MSG_PET_MOVE)      // 宠物移动
	MsgPetAction    uint32 = uint32(messages.PetMessageID_MSG_PET_ACTION)    // 宠物行动
	MsgPetLevelUp   uint32 = uint32(messages.PetMessageID_MSG_PET_LEVEL_UP)  // 宠物升级
	MsgPetEvolution uint32 = uint32(messages.PetMessageID_MSG_PET_EVOLUTION) // 宠物进化
	MsgPetTrain     uint32 = uint32(messages.PetMessageID_MSG_PET_TRAIN)     // 宠物训练
	MsgPetFeed      uint32 = uint32(messages.PetMessageID_MSG_PET_FEED)      // 宠物喂养
	MsgPetStatus    uint32 = uint32(messages.PetMessageID_MSG_PET_STATUS)    // 宠物状态

	// 建筑相关消息 (0x0400 - 0x04FF) - 使用proto生成的常量
	MsgBuildingCreate  uint32 = uint32(messages.BuildingMessageID_MSG_BUILDING_CREATE)  // 创建建筑
	MsgBuildingUpgrade uint32 = uint32(messages.BuildingMessageID_MSG_BUILDING_UPGRADE) // 升级建筑
	MsgBuildingDestroy uint32 = uint32(messages.BuildingMessageID_MSG_BUILDING_DESTROY) // 摧毁建筑
	MsgBuildingInfo    uint32 = uint32(messages.BuildingMessageID_MSG_BUILDING_INFO)    // 建筑信息
	MsgBuildingProduce uint32 = uint32(messages.BuildingMessageID_MSG_BUILDING_PRODUCE) // 建筑生产
	MsgBuildingCollect uint32 = uint32(messages.BuildingMessageID_MSG_BUILDING_COLLECT) // 收集资源
	MsgBuildingRepair  uint32 = uint32(messages.BuildingMessageID_MSG_BUILDING_REPAIR)  // 修复建筑
	MsgBuildingStatus  uint32 = uint32(messages.BuildingMessageID_MSG_BUILDING_STATUS)  // 建筑状态

	// 社交相关消息 (0x0500 - 0x05FF) - 使用proto生成的常量
	MsgChatMessage   uint32 = uint32(messages.SocialMessageID_MSG_CHAT_MESSAGE)   // 聊天消息
	MsgFriendRequest uint32 = uint32(messages.SocialMessageID_MSG_FRIEND_REQUEST) // 好友请求
	MsgFriendAccept  uint32 = uint32(messages.SocialMessageID_MSG_FRIEND_ACCEPT)  // 接受好友
	MsgFriendReject  uint32 = uint32(messages.SocialMessageID_MSG_FRIEND_REJECT)  // 拒绝好友
	MsgFriendRemove  uint32 = uint32(messages.SocialMessageID_MSG_FRIEND_REMOVE)  // 删除好友
	MsgFriendList    uint32 = uint32(messages.SocialMessageID_MSG_FRIEND_LIST)    // 好友列表
	MsgGuildCreate   uint32 = uint32(messages.SocialMessageID_MSG_GUILD_CREATE)   // 创建公会
	MsgGuildJoin     uint32 = uint32(messages.SocialMessageID_MSG_GUILD_JOIN)     // 加入公会
	MsgGuildLeave    uint32 = uint32(messages.SocialMessageID_MSG_GUILD_LEAVE)    // 离开公会
	MsgGuildInfo     uint32 = uint32(messages.SocialMessageID_MSG_GUILD_INFO)     // 公会信息
	MsgTeamCreate    uint32 = uint32(messages.SocialMessageID_MSG_TEAM_CREATE)    // 创建队伍
	MsgTeamJoin      uint32 = uint32(messages.SocialMessageID_MSG_TEAM_JOIN)      // 加入队伍
	MsgTeamLeave     uint32 = uint32(messages.SocialMessageID_MSG_TEAM_LEAVE)     // 离开队伍
	MsgTeamInfo      uint32 = uint32(messages.SocialMessageID_MSG_TEAM_INFO)      // 队伍信息

	// 物品相关消息 (0x0600 - 0x06FF) - 使用proto生成的常量
	MsgItemUse       uint32 = uint32(messages.ItemMessageID_MSG_ITEM_USE)       // 使用物品
	MsgItemEquip     uint32 = uint32(messages.ItemMessageID_MSG_ITEM_EQUIP)     // 装备物品
	MsgItemUnequip   uint32 = uint32(messages.ItemMessageID_MSG_ITEM_UNEQUIP)   // 卸下装备
	MsgItemDrop      uint32 = uint32(messages.ItemMessageID_MSG_ITEM_DROP)      // 丢弃物品
	MsgItemPickup    uint32 = uint32(messages.ItemMessageID_MSG_ITEM_PICKUP)    // 拾取物品
	MsgItemTrade     uint32 = uint32(messages.ItemMessageID_MSG_ITEM_TRADE)     // 交易物品
	MsgInventoryInfo uint32 = uint32(messages.ItemMessageID_MSG_INVENTORY_INFO) // 背包信息
	MsgItemInfo      uint32 = uint32(messages.ItemMessageID_MSG_ITEM_INFO)      // 物品信息
	MsgItemCraft     uint32 = uint32(messages.ItemMessageID_MSG_ITEM_CRAFT)     // 制作物品
	MsgItemEnhance   uint32 = uint32(messages.ItemMessageID_MSG_ITEM_ENHANCE)   // 强化物品

	// 任务相关消息 (0x0700 - 0x07FF) - 使用proto生成的常量
	MsgQuestAccept   uint32 = uint32(messages.QuestMessageID_MSG_QUEST_ACCEPT)   // 接受任务
	MsgQuestComplete uint32 = uint32(messages.QuestMessageID_MSG_QUEST_COMPLETE) // 完成任务
	MsgQuestCancel   uint32 = uint32(messages.QuestMessageID_MSG_QUEST_CANCEL)   // 取消任务
	MsgQuestProgress uint32 = uint32(messages.QuestMessageID_MSG_QUEST_PROGRESS) // 任务进度
	MsgQuestList     uint32 = uint32(messages.QuestMessageID_MSG_QUEST_LIST)     // 任务列表
	MsgQuestInfo     uint32 = uint32(messages.QuestMessageID_MSG_QUEST_INFO)     // 任务信息
	MsgQuestReward   uint32 = uint32(messages.QuestMessageID_MSG_QUEST_REWARD)   // 任务奖励

	// 查询相关消息 (0x0800 - 0x08FF) - 定义在game_protocol.go中
)

// 消息魔数
const MessageMagic uint32 = 0x47574B53 // "GWKS" - GreatestWorks

// 消息头大小
const MessageHeaderSize = 32 // 消息头固定大小

// ParseMessageHeader 解析消息头
func ParseMessageHeader(data []byte) (*MessageHeader, error) {
	if len(data) < MessageHeaderSize {
		return nil, fmt.Errorf("invalid message header size: %d", len(data))
	}

	header := &MessageHeader{}

	// 解析魔数 (4字节)
	header.Magic = uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])
	if header.Magic != MessageMagic {
		return nil, fmt.Errorf("invalid message magic: 0x%08X", header.Magic)
	}

	// 解析消息ID (4字节)
	header.MessageID = uint32(data[4])<<24 | uint32(data[5])<<16 | uint32(data[6])<<8 | uint32(data[7])

	// 解析消息类型 (4字节)
	header.MessageType = uint32(data[8])<<24 | uint32(data[9])<<16 | uint32(data[10])<<8 | uint32(data[11])

	// 解析标志位 (2字节)
	header.Flags = uint16(data[12])<<8 | uint16(data[13])

	// 解析玩家ID (8字节)
	header.PlayerID = uint64(data[14])<<56 | uint64(data[15])<<48 | uint64(data[16])<<40 | uint64(data[17])<<32 |
		uint64(data[18])<<24 | uint64(data[19])<<16 | uint64(data[20])<<8 | uint64(data[21])

	// 解析时间戳 (8字节)
	header.Timestamp = int64(data[22])<<56 | int64(data[23])<<48 | int64(data[24])<<40 | int64(data[25])<<32 |
		int64(data[26])<<24 | int64(data[27])<<16 | int64(data[28])<<8 | int64(data[29])

	// 解析序列号 (4字节)
	header.Sequence = uint32(data[30])<<8 | uint32(data[31])

	// 解析消息体长度 (4字节) - 这个在消息头中不包含，需要从其他地方获取
	header.Length = 0

	return header, nil
}

// MessageHeader TCP消息头
type MessageHeader struct {
	Magic       uint32 `json:"magic"`        // 魔数标识
	MessageID   uint32 `json:"message_id"`   // 消息ID（用于请求响应匹配）
	MessageType uint32 `json:"message_type"` // 消息类型
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

// 消息标志位 - 使用proto生成的常量
const (
	FlagRequest    uint16 = uint16(protocol.MessageFlag_MESSAGE_FLAG_REQUEST)    // 请求消息
	FlagResponse   uint16 = uint16(protocol.MessageFlag_MESSAGE_FLAG_RESPONSE)   // 响应消息
	FlagError      uint16 = uint16(protocol.MessageFlag_MESSAGE_FLAG_ERROR)      // 错误消息
	FlagAsync      uint16 = uint16(protocol.MessageFlag_MESSAGE_FLAG_ASYNC)      // 异步消息
	FlagBroadcast  uint16 = uint16(protocol.MessageFlag_MESSAGE_FLAG_BROADCAST)  // 广播消息
	FlagEncrypted  uint16 = uint16(protocol.MessageFlag_MESSAGE_FLAG_ENCRYPTED)  // 加密消息
	FlagCompressed uint16 = uint16(protocol.MessageFlag_MESSAGE_FLAG_COMPRESSED) // 压缩消息
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
