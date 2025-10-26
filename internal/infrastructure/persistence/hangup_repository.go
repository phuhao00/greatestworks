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

// HangupRepository 挂机仓储
type HangupRepository struct {
	collection *mongo.Collection
	logger     logging.Logger
}

// NewHangupRepository 创建挂机仓储
func NewHangupRepository(db *mongo.Database, logger logging.Logger) *HangupRepository {
	return &HangupRepository{
		collection: db.Collection("hangups"),
		logger:     logger,
	}
}

// HangupRecord 挂机记录
type HangupRecord struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PlayerID   string             `bson:"player_id" json:"player_id"`
	StartTime  time.Time          `bson:"start_time" json:"start_time"`
	EndTime    *time.Time         `bson:"end_time,omitempty" json:"end_time,omitempty"`
	Duration   int64              `bson:"duration" json:"duration"`
	Experience int64              `bson:"experience" json:"experience"`
	Gold       int64              `bson:"gold" json:"gold"`
	Status     string             `bson:"status" json:"status"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

// CreateHangup 创建挂机记录
func (r *HangupRepository) CreateHangup(ctx context.Context, record *HangupRecord) error {
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, record)
	if err != nil {
		r.logger.Error("创建挂机记录失败", err, logging.Fields{
			"player_id": record.PlayerID,
		})
		return fmt.Errorf("创建挂机记录失败: %w", err)
	}

	r.logger.Info("挂机记录创建成功", map[string]interface{}{
		"player_id":  record.PlayerID,
		"start_time": record.StartTime,
	})

	return nil
}

// GetHangup 获取挂机记录
func (r *HangupRepository) GetHangup(ctx context.Context, id string) (*HangupRecord, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的ID格式: %w", err)
	}

	var record HangupRecord
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("挂机记录不存在")
		}
		r.logger.Error("获取挂机记录失败", err, logging.Fields{
			"id": id,
		})
		return nil, fmt.Errorf("获取挂机记录失败: %w", err)
	}

	return &record, nil
}

// UpdateHangup 更新挂机记录
func (r *HangupRepository) UpdateHangup(ctx context.Context, id string, updates map[string]interface{}) error {
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
		r.logger.Error("更新挂机记录失败", err, logging.Fields{
			"id": id,
		})
		return fmt.Errorf("更新挂机记录失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("挂机记录不存在")
	}

	r.logger.Info("挂机记录更新成功", map[string]interface{}{
		"id": id,
	})

	return nil
}

// DeleteHangup 删除挂机记录
func (r *HangupRepository) DeleteHangup(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		r.logger.Error("删除挂机记录失败", err, logging.Fields{
			"id": id,
		})
		return fmt.Errorf("删除挂机记录失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("挂机记录不存在")
	}

	r.logger.Info("挂机记录删除成功", map[string]interface{}{
		"id": id,
	})

	return nil
}

// GetPlayerHangups 获取玩家的挂机记录
func (r *HangupRepository) GetPlayerHangups(ctx context.Context, playerID string, limit, offset int) ([]*HangupRecord, error) {
	filter := bson.M{"player_id": playerID}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("获取玩家挂机记录失败", err, logging.Fields{
			"player_id": playerID,
		})
		return nil, fmt.Errorf("获取玩家挂机记录失败: %w", err)
	}
	defer cursor.Close(ctx)

	var records []*HangupRecord
	if err = cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("解析挂机记录失败: %w", err)
	}

	r.logger.Info("获取玩家挂机记录成功", map[string]interface{}{
		"player_id": playerID,
		"count":     len(records),
	})

	return records, nil
}

// GetActiveHangups 获取活跃的挂机记录
func (r *HangupRepository) GetActiveHangups(ctx context.Context, playerID string) ([]*HangupRecord, error) {
	filter := bson.M{
		"player_id": playerID,
		"status":    "active",
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("获取活跃挂机记录失败", err, logging.Fields{
			"player_id": playerID,
		})
		return nil, fmt.Errorf("获取活跃挂机记录失败: %w", err)
	}
	defer cursor.Close(ctx)

	var records []*HangupRecord
	if err = cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("解析挂机记录失败: %w", err)
	}

	return records, nil
}

// EndHangup 结束挂机
func (r *HangupRepository) EndHangup(ctx context.Context, id string, endTime time.Time, experience, gold int64) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	updates := bson.M{
		"end_time":   endTime,
		"duration":   time.Since(endTime).Seconds(),
		"experience": experience,
		"gold":       gold,
		"status":     "completed",
		"updated_at": time.Now(),
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)
	if err != nil {
		r.logger.Error("结束挂机失败", err, logging.Fields{
			"id": id,
		})
		return fmt.Errorf("结束挂机失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("挂机记录不存在")
	}

	r.logger.Info("挂机结束成功", map[string]interface{}{
		"id":         id,
		"experience": experience,
		"gold":       gold,
	})

	return nil
}
