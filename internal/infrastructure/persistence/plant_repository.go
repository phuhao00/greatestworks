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

// PlantRepository 植物仓储
type PlantRepository struct {
	collection *mongo.Collection
	logger     logging.Logger
}

// NewPlantRepository 创建植物仓储
func NewPlantRepository(db *mongo.Database, logger logging.Logger) *PlantRepository {
	return &PlantRepository{
		collection: db.Collection("plants"),
		logger:     logger,
	}
}

// PlantRecord 植物记录
type PlantRecord struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PlayerID   string             `bson:"player_id" json:"player_id"`
	PlantType  string             `bson:"plant_type" json:"plant_type"`
	Position   Position           `bson:"position" json:"position"`
	Level      int                `bson:"level" json:"level"`
	Growth     int                `bson:"growth" json:"growth"`
	MaxGrowth  int                `bson:"max_growth" json:"max_growth"`
	WaterLevel int                `bson:"water_level" json:"water_level"`
	Fertilizer int                `bson:"fertilizer" json:"fertilizer"`
	Status     string             `bson:"status" json:"status"`
	PlantedAt  time.Time          `bson:"planted_at" json:"planted_at"`
	HarvestAt  *time.Time         `bson:"harvest_at,omitempty" json:"harvest_at,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

// CreatePlant 创建植物
func (r *PlantRepository) CreatePlant(ctx context.Context, plant *PlantRecord) error {
	plant.CreatedAt = time.Now()
	plant.UpdatedAt = time.Now()
	plant.PlantedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, plant)
	if err != nil {
		r.logger.Error("创建植物失败", err, logging.Fields{
			"player_id":  plant.PlayerID,
			"plant_type": plant.PlantType,
		})
		return fmt.Errorf("创建植物失败: %w", err)
	}

	r.logger.Info("植物创建成功", map[string]interface{}{
		"player_id":  plant.PlayerID,
		"plant_type": plant.PlantType,
		"level":      plant.Level,
	})

	return nil
}

// GetPlant 获取植物
func (r *PlantRepository) GetPlant(ctx context.Context, id string) (*PlantRecord, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的ID格式: %w", err)
	}

	var plant PlantRecord
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&plant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("植物不存在")
		}
		r.logger.Error("获取植物失败", err, logging.Fields{
			"id": id,
		})
		return nil, fmt.Errorf("获取植物失败: %w", err)
	}

	return &plant, nil
}

// UpdatePlant 更新植物
func (r *PlantRepository) UpdatePlant(ctx context.Context, id string, updates map[string]interface{}) error {
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
		r.logger.Error("更新植物失败", err, logging.Fields{
			"id": id,
		})
		return fmt.Errorf("更新植物失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("植物不存在")
	}

	r.logger.Info("植物更新成功", map[string]interface{}{
		"id": id,
	})

	return nil
}

// DeletePlant 删除植物
func (r *PlantRepository) DeletePlant(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		r.logger.Error("删除植物失败", err, logging.Fields{
			"id": id,
		})
		return fmt.Errorf("删除植物失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("植物不存在")
	}

	r.logger.Info("植物删除成功", map[string]interface{}{
		"id": id,
	})

	return nil
}

// GetPlayerPlants 获取玩家的植物列表
func (r *PlantRepository) GetPlayerPlants(ctx context.Context, playerID string, limit, offset int) ([]*PlantRecord, error) {
	filter := bson.M{"player_id": playerID}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("获取玩家植物失败", err, logging.Fields{
			"player_id": playerID,
		})
		return nil, fmt.Errorf("获取玩家植物失败: %w", err)
	}
	defer cursor.Close(ctx)

	var plants []*PlantRecord
	if err = cursor.All(ctx, &plants); err != nil {
		return nil, fmt.Errorf("解析植物列表失败: %w", err)
	}

	r.logger.Info("获取玩家植物成功", map[string]interface{}{
		"player_id": playerID,
		"count":     len(plants),
	})

	return plants, nil
}

// GetPlantsByType 根据类型获取植物列表
func (r *PlantRepository) GetPlantsByType(ctx context.Context, plantType string, limit, offset int) ([]*PlantRecord, error) {
	filter := bson.M{"plant_type": plantType}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("根据类型获取植物失败", err, logging.Fields{
			"plant_type": plantType,
		})
		return nil, fmt.Errorf("根据类型获取植物失败: %w", err)
	}
	defer cursor.Close(ctx)

	var plants []*PlantRecord
	if err = cursor.All(ctx, &plants); err != nil {
		return nil, fmt.Errorf("解析植物列表失败: %w", err)
	}

	r.logger.Info("根据类型获取植物成功", map[string]interface{}{
		"plant_type": plantType,
		"count":      len(plants),
	})

	return plants, nil
}

// WaterPlant 浇水
func (r *PlantRepository) WaterPlant(ctx context.Context, id string, waterAmount int) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	updates := bson.M{
		"$inc": bson.M{
			"water_level": waterAmount,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		updates,
	)
	if err != nil {
		r.logger.Error("植物浇水失败", err, logging.Fields{
			"id":           id,
			"water_amount": waterAmount,
		})
		return fmt.Errorf("植物浇水失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("植物不存在")
	}

	r.logger.Info("植物浇水成功", map[string]interface{}{
		"id":           id,
		"water_amount": waterAmount,
	})

	return nil
}

// FertilizePlant 施肥
func (r *PlantRepository) FertilizePlant(ctx context.Context, id string, fertilizerAmount int) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	updates := bson.M{
		"$inc": bson.M{
			"fertilizer": fertilizerAmount,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		updates,
	)
	if err != nil {
		r.logger.Error("植物施肥失败", err, logging.Fields{
			"id":                id,
			"fertilizer_amount": fertilizerAmount,
		})
		return fmt.Errorf("植物施肥失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("植物不存在")
	}

	r.logger.Info("植物施肥成功", map[string]interface{}{
		"id":                id,
		"fertilizer_amount": fertilizerAmount,
	})

	return nil
}

// HarvestPlant 收获植物
func (r *PlantRepository) HarvestPlant(ctx context.Context, id string) (*PlantRecord, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的ID格式: %w", err)
	}

	harvestTime := time.Now()
	updates := bson.M{
		"harvest_at": harvestTime,
		"status":     "harvested",
		"updated_at": harvestTime,
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)
	if err != nil {
		r.logger.Error("收获植物失败", err, logging.Fields{
			"id": id,
		})
		return nil, fmt.Errorf("收获植物失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("植物不存在")
	}

	// 获取更新后的植物记录
	plant, err := r.GetPlant(ctx, id)
	if err != nil {
		return nil, err
	}

	r.logger.Info("植物收获成功", map[string]interface{}{
		"id":         id,
		"harvest_at": harvestTime,
	})

	return plant, nil
}
