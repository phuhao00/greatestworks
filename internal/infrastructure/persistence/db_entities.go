package persistence

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DbUser 用户数据库实体
type DbUser struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       int64              `bson:"user_id"`
	Username     string             `bson:"username"`
	PasswordHash string             `bson:"password_hash"`
	Status       int32              `bson:"status"` // 0:正常 1:封禁
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	LastLoginAt  time.Time          `bson:"last_login_at"`
}

// DbCharacter 角色数据库实体
type DbCharacter struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	CharacterID int64              `bson:"character_id"`
	UserID      int64              `bson:"user_id"`
	Name        string             `bson:"name"`
	Race        int32              `bson:"race"`
	Class       int32              `bson:"class"`
	Level       int32              `bson:"level"`
	Exp         int64              `bson:"exp"`
	Gold        int64              `bson:"gold"`

	// 位置信息
	MapID     int32   `bson:"map_id"`
	PositionX float32 `bson:"position_x"`
	PositionY float32 `bson:"position_y"`
	PositionZ float32 `bson:"position_z"`
	Direction float32 `bson:"direction"`

	// 属性
	HP    int32 `bson:"hp"`
	MP    int32 `bson:"mp"`
	MaxHP int32 `bson:"max_hp"`
	MaxMP int32 `bson:"max_mp"`

	// 基础属性
	STR int32 `bson:"str"` // 力量
	INT int32 `bson:"int"` // 智力
	AGI int32 `bson:"agi"` // 敏捷
	VIT int32 `bson:"vit"` // 体力
	SPR int32 `bson:"spr"` // 精神

	// 战斗属性
	AD  int32 `bson:"ad"`  // 物理攻击
	AP  int32 `bson:"ap"`  // 魔法攻击
	DEF int32 `bson:"def"` // 物理防御
	RES int32 `bson:"res"` // 魔法抗性
	SPD int32 `bson:"spd"` // 速度

	CRI     int32 `bson:"cri"`      // 暴击率
	CRID    int32 `bson:"crid"`     // 暴击伤害
	HitRate int32 `bson:"hit_rate"` // 命中率
	Dodge   int32 `bson:"dodge"`    // 闪避率

	// 时间戳
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	DeletedAt time.Time `bson:"deleted_at,omitempty"`
}

// DbItem 物品数据库实体（背包/装备/仓库）
type DbItem struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ItemUID     int64              `bson:"item_uid"`     // 物品唯一ID
	CharacterID int64              `bson:"character_id"` // 所属角色
	ItemID      int32              `bson:"item_id"`      // 物品配置ID
	Count       int32              `bson:"count"`        // 数量
	Slot        int32              `bson:"slot"`         // 槽位
	Location    int32              `bson:"location"`     // 位置（背包/装备/仓库）
	Bound       bool               `bson:"bound"`        // 是否绑定
	Expire      int64              `bson:"expire"`       // 过期时间戳
	CreatedAt   time.Time          `bson:"created_at"`
}

// DbQuest 任务进度数据库实体
type DbQuest struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	CharacterID int64              `bson:"character_id"`
	QuestID     int32              `bson:"quest_id"`
	Status      int32              `bson:"status"` // 0:进行中 1:已完成 2:已领取
	Objectives  []DbObjective      `bson:"objectives"`
	AcceptedAt  time.Time          `bson:"accepted_at"`
	CompletedAt time.Time          `bson:"completed_at,omitempty"`
}

// DbObjective 任务目标
type DbObjective struct {
	Type     int32 `bson:"type"`      // 目标类型
	TargetID int32 `bson:"target_id"` // 目标ID
	Required int32 `bson:"required"`  // 需要数量
	Current  int32 `bson:"current"`   // 当前数量
}

// DbMail 邮件数据库实体
type DbMail struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	MailID      int64              `bson:"mail_id"`
	ReceiverID  int64              `bson:"receiver_id"`
	SenderName  string             `bson:"sender_name"`
	Title       string             `bson:"title"`
	Content     string             `bson:"content"`
	IsRead      bool               `bson:"is_read"`
	HasItems    bool               `bson:"has_items"`
	Attachments []DbAttachment     `bson:"attachments"`
	ExpireAt    time.Time          `bson:"expire_at"`
	CreatedAt   time.Time          `bson:"created_at"`
}

// DbAttachment 邮件附件
type DbAttachment struct {
	ItemID int32 `bson:"item_id"`
	Count  int32 `bson:"count"`
}

// DbGuild 公会数据库实体
type DbGuild struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	GuildID     int64              `bson:"guild_id"`
	Name        string             `bson:"name"`
	LeaderID    int64              `bson:"leader_id"`
	Level       int32              `bson:"level"`
	Exp         int64              `bson:"exp"`
	Notice      string             `bson:"notice"`
	MemberCount int32              `bson:"member_count"`
	MaxMembers  int32              `bson:"max_members"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

// DbGuildMember 公会成员数据库实体
type DbGuildMember struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	GuildID      int64              `bson:"guild_id"`
	CharacterID  int64              `bson:"character_id"`
	Rank         int32              `bson:"rank"` // 0:会长 1:副会长 2:精英 3:成员
	Contribution int64              `bson:"contribution"`
	JoinedAt     time.Time          `bson:"joined_at"`
}
