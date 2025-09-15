// Package game 游戏核心逻辑
// Author: MMO Server Team
// Created: 2024

package game

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Player 玩家数据结构
type Player struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"user_id"`
	Username    string             `bson:"username" json:"username"`
	Level       int32              `bson:"level" json:"level"`
	Experience  int64              `bson:"experience" json:"experience"`
	Gold        int64              `bson:"gold" json:"gold"`
	Diamond     int64              `bson:"diamond" json:"diamond"`
	Position    Position           `bson:"position" json:"position"`
	Attributes  PlayerAttributes   `bson:"attributes" json:"attributes"`
	Inventory   []Item             `bson:"inventory" json:"inventory"`
	Equipment   Equipment          `bson:"equipment" json:"equipment"`
	Skills      []Skill            `bson:"skills" json:"skills"`
	Quests      []Quest            `bson:"quests" json:"quests"`
	GuildID     string             `bson:"guild_id,omitempty" json:"guild_id,omitempty"`
	LastLogin   time.Time          `bson:"last_login" json:"last_login"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	Online      bool               `bson:"online" json:"online"`
	ServerID    string             `bson:"server_id" json:"server_id"`
}

// Position 位置信息
type Position struct {
	X     float64 `bson:"x" json:"x"`
	Y     float64 `bson:"y" json:"y"`
	Z     float64 `bson:"z" json:"z"`
	MapID string  `bson:"map_id" json:"map_id"`
}

// PlayerAttributes 玩家属性
type PlayerAttributes struct {
	HP           int32 `bson:"hp" json:"hp"`
	MaxHP        int32 `bson:"max_hp" json:"max_hp"`
	MP           int32 `bson:"mp" json:"mp"`
	MaxMP        int32 `bson:"max_mp" json:"max_mp"`
	Attack       int32 `bson:"attack" json:"attack"`
	Defense      int32 `bson:"defense" json:"defense"`
	Speed        int32 `bson:"speed" json:"speed"`
	CriticalRate int32 `bson:"critical_rate" json:"critical_rate"`
	DodgeRate    int32 `bson:"dodge_rate" json:"dodge_rate"`
}

// Item 物品
type Item struct {
	ID       string `bson:"id" json:"id"`
	ItemID   string `bson:"item_id" json:"item_id"`
	Quantity int32  `bson:"quantity" json:"quantity"`
	Slot     int32  `bson:"slot" json:"slot"`
}

// Equipment 装备
type Equipment struct {
	Weapon     *Item `bson:"weapon,omitempty" json:"weapon,omitempty"`
	Armor      *Item `bson:"armor,omitempty" json:"armor,omitempty"`
	Helmet     *Item `bson:"helmet,omitempty" json:"helmet,omitempty"`
	Boots      *Item `bson:"boots,omitempty" json:"boots,omitempty"`
	Gloves     *Item `bson:"gloves,omitempty" json:"gloves,omitempty"`
	Accessory1 *Item `bson:"accessory1,omitempty" json:"accessory1,omitempty"`
	Accessory2 *Item `bson:"accessory2,omitempty" json:"accessory2,omitempty"`
}

// Skill 技能
type Skill struct {
	SkillID string `bson:"skill_id" json:"skill_id"`
	Level   int32  `bson:"level" json:"level"`
	Exp     int64  `bson:"exp" json:"exp"`
}

// Quest 任务
type Quest struct {
	QuestID   string            `bson:"quest_id" json:"quest_id"`
	Status    string            `bson:"status" json:"status"` // pending, active, completed, failed
	Progress  map[string]int32  `bson:"progress" json:"progress"`
	StartTime time.Time         `bson:"start_time" json:"start_time"`
	EndTime   *time.Time        `bson:"end_time,omitempty" json:"end_time,omitempty"`
}

// PlayerManager 玩家管理器
type PlayerManager struct {
	collection *mongo.Collection
	mu         sync.RWMutex
	onlinePlayers map[string]*Player
}

// NewPlayerManager 创建玩家管理器
func NewPlayerManager(collection *mongo.Collection) *PlayerManager {
	return &PlayerManager{
		collection:    collection,
		onlinePlayers: make(map[string]*Player),
	}
}

// CreatePlayer 创建新玩家
func (pm *PlayerManager) CreatePlayer(ctx context.Context, userID, username string) (*Player, error) {
	player := &Player{
		UserID:   userID,
		Username: username,
		Level:    1,
		Gold:     1000,
		Diamond:  0,
		Position: Position{
			X:     0,
			Y:     0,
			Z:     0,
			MapID: "starter_map",
		},
		Attributes: PlayerAttributes{
			HP:           100,
			MaxHP:        100,
			MP:           50,
			MaxMP:        50,
			Attack:       10,
			Defense:      5,
			Speed:        10,
			CriticalRate: 5,
			DodgeRate:    5,
		},
		Inventory: make([]Item, 0),
		Skills:    make([]Skill, 0),
		Quests:    make([]Quest, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Online:    false,
	}

	result, err := pm.collection.InsertOne(ctx, player)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	player.ID = result.InsertedID.(primitive.ObjectID)
	return player, nil
}

// GetPlayer 获取玩家信息
func (pm *PlayerManager) GetPlayer(ctx context.Context, userID string) (*Player, error) {
	var player Player
	err := pm.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&player)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("player not found")
		}
		return nil, fmt.Errorf("failed to get player: %w", err)
	}
	return &player, nil
}

// UpdatePlayer 更新玩家信息
func (pm *PlayerManager) UpdatePlayer(ctx context.Context, player *Player) error {
	player.UpdatedAt = time.Now()
	_, err := pm.collection.UpdateOne(
		ctx,
		bson.M{"user_id": player.UserID},
		bson.M{"$set": player},
	)
	if err != nil {
		return fmt.Errorf("failed to update player: %w", err)
	}
	return nil
}

// SetPlayerOnline 设置玩家在线状态
func (pm *PlayerManager) SetPlayerOnline(ctx context.Context, userID string, online bool, serverID string) error {
	update := bson.M{
		"online":     online,
		"updated_at": time.Now(),
	}
	if online {
		update["last_login"] = time.Now()
		update["server_id"] = serverID
	} else {
		update["server_id"] = ""
	}

	_, err := pm.collection.UpdateOne(
		ctx,
		bson.M{"user_id": userID},
		bson.M{"$set": update},
	)
	if err != nil {
		return fmt.Errorf("failed to set player online status: %w", err)
	}

	// 更新内存中的在线玩家列表
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if online {
		player, err := pm.GetPlayer(ctx, userID)
		if err == nil {
			pm.onlinePlayers[userID] = player
		}
	} else {
		delete(pm.onlinePlayers, userID)
	}

	return nil
}

// GetOnlinePlayers 获取在线玩家列表
func (pm *PlayerManager) GetOnlinePlayers() map[string]*Player {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	result := make(map[string]*Player)
	for k, v := range pm.onlinePlayers {
		result[k] = v
	}
	return result
}