package persistence

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"greatestworks/internal/infrastructure/logging"
)

// NPCRepository NPC仓储
type NPCRepository struct {
	collection *mongo.Collection
	logger     logging.Logger
}

// NewNPCRepository 创建NPC仓储
func NewNPCRepository(db *mongo.Database, logger logging.Logger) *NPCRepository {
	return &NPCRepository{
		collection: db.Collection("npcs"),
		logger:     logger,
	}
}

// NPCRecord NPC记录
type NPCRecord struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Type      string             `bson:"type" json:"type"`
	Level     int                `bson:"level" json:"level"`
	Health    int64              `bson:"health" json:"health"`
	MaxHealth int64              `bson:"max_health" json:"max_health"`
	Attack    int64              `bson:"attack" json:"attack"`
	Defense   int64              `bson:"defense" json:"defense"`
	Position  Position           `bson:"position" json:"position"`
	Status    string             `bson:"status" json:"status"`
	LastSeen  time.Time          `bson:"last_seen" json:"last_seen"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// Position 位置信息
type Position struct {
	X float64 `bson:"x" json:"x"`
	Y float64 `bson:"y" json:"y"`
	Z float64 `bson:"z" json:"z"`
}

// CreateNPC 创建NPC
func (r *NPCRepository) CreateNPC(ctx context.Context, npc *NPCRecord) error {
	npc.CreatedAt = time.Now()
	npc.UpdatedAt = time.Now()
	npc.LastSeen = time.Now()

	_, err := r.collection.InsertOne(ctx, npc)
	if err != nil {
		r.logger.Error("创建NPC失败", err, logging.Fields{
			"name": npc.Name,
			"type": npc.Type,
		})
		return fmt.Errorf("创建NPC失败: %w", err)
	}

	r.logger.Info("NPC创建成功", map[string]interface{}{
		"name":  npc.Name,
		"type":  npc.Type,
		"level": npc.Level,
	})

	return nil
}

// GetNPC 获取NPC
func (r *NPCRepository) GetNPC(ctx context.Context, id string) (*NPCRecord, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的ID格式: %w", err)
	}

	var npc NPCRecord
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&npc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("NPC不存在")
		}
		r.logger.Error("获取NPC失败", err, logging.Fields{
			"id": id,
		})
		return nil, fmt.Errorf("获取NPC失败: %w", err)
	}

	return &npc, nil
}

// UpdateNPC 更新NPC
func (r *NPCRepository) UpdateNPC(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	updates["updated_at"] = time.Now()

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)
	if err != nil {
		r.logger.Error("更新NPC失败", err, logging.Fields{
			"id": id,
		})
		return fmt.Errorf("更新NPC失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("NPC不存在")
	}

	r.logger.Info("NPC更新成功", map[string]interface{}{
		"id": id,
	})

	return nil
}

// DeleteNPC 删除NPC
func (r *NPCRepository) DeleteNPC(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		r.logger.Error("删除NPC失败", err, logging.Fields{
			"id": id,
		})
		return fmt.Errorf("删除NPC失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("NPC不存在")
	}

	r.logger.Info("NPC删除成功", map[string]interface{}{
		"id": id,
	})

	return nil
}

// GetNPCsByType 根据类型获取NPC列表
func (r *NPCRepository) GetNPCsByType(ctx context.Context, npcType string, limit, offset int) ([]*NPCRecord, error) {
	filter := bson.M{"type": npcType}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("根据类型获取NPC失败", err, logging.Fields{
			"type": npcType,
		})
		return nil, fmt.Errorf("根据类型获取NPC失败: %w", err)
	}
	defer cursor.Close(ctx)

	var npcs []*NPCRecord
	if err = cursor.All(ctx, &npcs); err != nil {
		return nil, fmt.Errorf("解析NPC列表失败: %w", err)
	}

	r.logger.Info("根据类型获取NPC成功", map[string]interface{}{
		"type":  npcType,
		"count": len(npcs),
	})

	return npcs, nil
}

// GetNPCsByPosition 根据位置获取NPC列表
func (r *NPCRepository) GetNPCsByPosition(ctx context.Context, x, y, z float64, radius float64) ([]*NPCRecord, error) {
	filter := bson.M{
		"position.x": bson.M{
			"$gte": x - radius,
			"$lte": x + radius,
		},
		"position.y": bson.M{
			"$gte": y - radius,
			"$lte": y + radius,
		},
		"position.z": bson.M{
			"$gte": z - radius,
			"$lte": z + radius,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("根据位置获取NPC失败", err, logging.Fields{
			"x":      x,
			"y":      y,
			"z":      z,
			"radius": radius,
		})
		return nil, fmt.Errorf("根据位置获取NPC失败: %w", err)
	}
	defer cursor.Close(ctx)

	var npcs []*NPCRecord
	if err = cursor.All(ctx, &npcs); err != nil {
		return nil, fmt.Errorf("解析NPC列表失败: %w", err)
	}

	r.logger.Info("根据位置获取NPC成功", map[string]interface{}{
		"x":      x,
		"y":      y,
		"z":      z,
		"radius": radius,
		"count":  len(npcs),
	})

	return npcs, nil
}

// UpdateNPCPosition 更新NPC位置
func (r *NPCRepository) UpdateNPCPosition(ctx context.Context, id string, position Position) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	updates := bson.M{
		"position":   position,
		"last_seen":  time.Now(),
		"updated_at": time.Now(),
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)
	if err != nil {
		r.logger.Error("更新NPC位置失败", err, logging.Fields{
			"id":       id,
			"position": position,
		})
		return fmt.Errorf("更新NPC位置失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("NPC不存在")
	}

	r.logger.Info("NPC位置更新成功", map[string]interface{}{
		"id":       id,
		"position": position,
	})

	return nil
}
