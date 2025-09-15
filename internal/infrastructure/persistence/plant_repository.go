package persistence

import (
	"context"
	"fmt"
	"time"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	"greatestworks/internal/domain/scene/plant"
	"greatestworks/internal/infrastructure/cache"
	"greatestworks/aop/logger"
)

// MongoFarmRepository MongoDB农场仓储实现
type MongoFarmRepository struct {
	db         *mongo.Database
	cache      cache.Cache
	logger     logger.Logger
	collection *mongo.Collection
}

// NewMongoFarmRepository 创建MongoDB农场仓储
func NewMongoFarmRepository(db *mongo.Database, cache cache.Cache, logger logger.Logger) plant.FarmRepository {
	return &MongoFarmRepository{
		db:         db,
		cache:      cache,
		logger:     logger,
		collection: db.Collection("farms"),
	}
}

// FarmDocument MongoDB农场文档结构
type FarmDocument struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	FarmID     string             `bson:"farm_id"`
	PlayerID   string             `bson:"player_id"`
	Level      int                `bson:"level"`
	Experience int64              `bson:"experience"`
	Plots      []PlotDocument     `bson:"plots"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

// PlotDocument 地块文档结构
type PlotDocument struct {
	ID        string  `bson:"id"`
	Level     int     `bson:"level"`
	SoilType  string  `bson:"soil_type"`
	Fertility float64 `bson:"fertility"`
	Moisture  float64 `bson:"moisture"`
	CropID    string  `bson:"crop_id,omitempty"`
}

// Save 保存农场
func (r *MongoFarmRepository) Save(farm *plant.FarmAggregate) error {
	ctx := context.Background()
	doc := r.aggregateToDocument(farm)
	doc.UpdatedAt = time.Now()
	
	if doc.ID.IsZero() {
		doc.CreatedAt = time.Now()
		result, err := r.collection.InsertOne(ctx, doc)
		if err != nil {
			r.logger.Error("Failed to insert farm", "error", err, "player_id", farm.GetPlayerID())
			return fmt.Errorf("failed to insert farm: %w", err)
		}
		
		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
			doc.ID = oid
		}
	} else {
		filter := bson.M{"player_id": farm.GetPlayerID()}
		update := bson.M{"$set": doc}
		
		_, err := r.collection.UpdateOne(ctx, filter, update)
		if err != nil {
			r.logger.Error("Failed to update farm", "error", err, "player_id", farm.GetPlayerID())
			return fmt.Errorf("failed to update farm: %w", err)
		}
	}
	
	// 更新缓存
	cacheKey := fmt.Sprintf("farm:%s", farm.GetPlayerID())
	if err := r.cache.Set(ctx, cacheKey, farm, time.Hour); err != nil {
		r.logger.Warn("Failed to cache farm", "error", err, "player_id", farm.GetPlayerID())
	}
	
	return nil
}

// FindByPlayer 根据玩家查找农场
func (r *MongoFarmRepository) FindByPlayer(playerID string) (*plant.FarmAggregate, error) {
	ctx := context.Background()
	
	// 先从缓存获取
	cacheKey := fmt.Sprintf("farm:%s", playerID)
	var cachedFarm *plant.FarmAggregate
	if err := r.cache.Get(ctx, cacheKey, &cachedFarm); err == nil && cachedFarm != nil {
		return cachedFarm, nil
	}
	
	// 从数据库获取
	filter := bson.M{"player_id": playerID}
	var doc FarmDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, plant.ErrFarmNotFound
		}
		r.logger.Error("Failed to find farm", "error", err, "player_id", playerID)
		return nil, fmt.Errorf("failed to find farm: %w", err)
	}
	
	farm := r.documentToAggregate(&doc)
	
	// 更新缓存
	if err := r.cache.Set(ctx, cacheKey, farm, time.Hour); err != nil {
		r.logger.Warn("Failed to cache farm", "error", err, "player_id", playerID)
	}
	
	return farm, nil
}

// Update 更新农场
func (r *MongoFarmRepository) Update(farm *plant.FarmAggregate) error {
	return r.Save(farm)
}

// Delete 删除农场
func (r *MongoFarmRepository) Delete(playerID string) error {
	ctx := context.Background()
	
	filter := bson.M{"player_id": playerID}
	
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete farm", "error", err, "player_id", playerID)
		return fmt.Errorf("failed to delete farm: %w", err)
	}
	
	if result.DeletedCount == 0 {
		return plant.ErrFarmNotFound
	}
	
	// 清除缓存
	cacheKey := fmt.Sprintf("farm:%s", playerID)
	if err := r.cache.Delete(ctx, cacheKey); err != nil {
		r.logger.Warn("Failed to delete farm cache", "error", err, "player_id", playerID)
	}
	
	return nil
}

// 私有方法

// aggregateToDocument 聚合根转文档
func (r *MongoFarmRepository) aggregateToDocument(farm *plant.FarmAggregate) *FarmDocument {
	plots := make([]PlotDocument, 0)
	for _, plot := range farm.GetPlots() {
		plots = append(plots, PlotDocument{
			ID:        plot.GetID(),
			Level:     plot.GetLevel(),
			SoilType:  string(plot.GetSoilType()),
			Fertility: plot.GetFertility(),
			Moisture:  plot.GetMoisture(),
			CropID:    plot.GetCropID(),
		})
	}
	
	return &FarmDocument{
		FarmID:     farm.GetID(),
		PlayerID:   farm.GetPlayerID(),
		Level:      farm.GetLevel(),
		Experience: farm.GetExperience(),
		Plots:      plots,
		CreatedAt:  farm.GetCreatedAt(),
		UpdatedAt:  farm.GetUpdatedAt(),
	}
}

// documentToAggregate 文档转聚合根
func (r *MongoFarmRepository) documentToAggregate(doc *FarmDocument) *plant.FarmAggregate {
	plots := make([]*plant.Plot, len(doc.Plots))
	for i, plotDoc := range doc.Plots {
		plots[i] = plant.NewPlot(
			plotDoc.ID,
			plotDoc.Level,
			plant.SoilType(plotDoc.SoilType),
			plotDoc.Fertility,
			plotDoc.Moisture,
			plotDoc.CropID,
		)
	}
	
	// 这里需要根据实际的FarmAggregate构造函数来实现
	return plant.ReconstructFarmAggregate(
		doc.FarmID,
		doc.PlayerID,
		doc.Level,
		doc.Experience,
		plots,
		doc.CreatedAt,
		doc.UpdatedAt,
	)
}

// MongoCropRepository MongoDB作物仓储实现
type MongoCropRepository struct {
	db         *mongo.Database
	cache      cache.Cache
	logger     logger.Logger
	collection *mongo.Collection
}

// NewMongoCropRepository 创建MongoDB作物仓储
func NewMongoCropRepository(db *mongo.Database, cache cache.Cache, logger logger.Logger) plant.CropRepository {
	return &MongoCropRepository{
		db:         db,
		cache:      cache,
		logger:     logger,
		collection: db.Collection("crops"),
	}
}

// CropDocument MongoDB作物文档结构
type CropDocument struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty"`
	CropID               string             `bson:"crop_id"`
	PlayerID             string             `bson:"player_id"`
	PlotID               string             `bson:"plot_id"`
	SeedID               string             `bson:"seed_id"`
	CropType             string             `bson:"crop_type"`
	CurrentStage         string             `bson:"current_stage"`
	GrowthProgress       float64            `bson:"growth_progress"`
	Health               float64            `bson:"health"`
	Moisture             float64            `bson:"moisture"`
	Nutrition            float64            `bson:"nutrition"`
	Quality              float64            `bson:"quality"`
	PlantedAt            time.Time          `bson:"planted_at"`
	LastWatered          time.Time          `bson:"last_watered"`
	LastFertilized       time.Time          `bson:"last_fertilized"`
	EstimatedHarvestTime time.Time          `bson:"estimated_harvest_time"`
	CreatedAt            time.Time          `bson:"created_at"`
	UpdatedAt            time.Time          `bson:"updated_at"`
}

// Save 保存作物
func (r *MongoCropRepository) Save(crop *plant.CropAggregate) error {
	ctx := context.Background()
	doc := r.cropAggregateToDocument(crop)
	doc.UpdatedAt = time.Now()
	
	if doc.ID.IsZero() {
		doc.CreatedAt = time.Now()
		result, err := r.collection.InsertOne(ctx, doc)
		if err != nil {
			r.logger.Error("Failed to insert crop", "error", err, "crop_id", crop.GetID())
			return fmt.Errorf("failed to insert crop: %w", err)
		}
		
		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
			doc.ID = oid
		}
	} else {
		filter := bson.M{"crop_id": crop.GetID()}
		update := bson.M{"$set": doc}
		
		_, err := r.collection.UpdateOne(ctx, filter, update)
		if err != nil {
			r.logger.Error("Failed to update crop", "error", err, "crop_id", crop.GetID())
			return fmt.Errorf("failed to update crop: %w", err)
		}
	}
	
	// 更新缓存
	cacheKey := fmt.Sprintf("crop:%s", crop.GetID())
	if err := r.cache.Set(ctx, cacheKey, crop, time.Hour); err != nil {
		r.logger.Warn("Failed to cache crop", "error", err, "crop_id", crop.GetID())
	}
	
	return nil
}

// FindByID 根据ID查找作物
func (r *MongoCropRepository) FindByID(cropID string) (*plant.CropAggregate, error) {
	ctx := context.Background()
	
	// 先从缓存获取
	cacheKey := fmt.Sprintf("crop:%s", cropID)
	var cachedCrop *plant.CropAggregate
	if err := r.cache.Get(ctx, cacheKey, &cachedCrop); err == nil && cachedCrop != nil {
		return cachedCrop, nil
	}
	
	// 从数据库获取
	filter := bson.M{"crop_id": cropID}
	var doc CropDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, plant.ErrCropNotFound
		}
		r.logger.Error("Failed to find crop", "error", err, "crop_id", cropID)
		return nil, fmt.Errorf("failed to find crop: %w", err)
	}
	
	crop := r.cropDocumentToAggregate(&doc)
	
	// 更新缓存
	if err := r.cache.Set(ctx, cacheKey, crop, time.Hour); err != nil {
		r.logger.Warn("Failed to cache crop", "error", err, "crop_id", cropID)
	}
	
	return crop, nil
}

// FindByPlayer 根据玩家查找作物
func (r *MongoCropRepository) FindByPlayer(playerID string) ([]*plant.CropAggregate, error) {
	ctx := context.Background()
	
	filter := bson.M{"player_id": playerID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find crops by player", "error", err, "player_id", playerID)
		return nil, fmt.Errorf("failed to find crops by player: %w", err)
	}
	defer cursor.Close(ctx)
	
	var docs []CropDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode crops by player", "error", err, "player_id", playerID)
		return nil, fmt.Errorf("failed to decode crops by player: %w", err)
	}
	
	crops := make([]*plant.CropAggregate, len(docs))
	for i, doc := range docs {
		crops[i] = r.cropDocumentToAggregate(&doc)
	}
	
	return crops, nil
}

// FindGrowingCrops 查找正在生长的作物
func (r *MongoCropRepository) FindGrowingCrops() ([]*plant.CropAggregate, error) {
	ctx := context.Background()
	
	filter := bson.M{
		"current_stage": bson.M{"$ne": "harvested"},
		"growth_progress": bson.M{"$lt": 1.0},
	}
	
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find growing crops", "error", err)
		return nil, fmt.Errorf("failed to find growing crops: %w", err)
	}
	defer cursor.Close(ctx)
	
	var docs []CropDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode growing crops", "error", err)
		return nil, fmt.Errorf("failed to decode growing crops: %w", err)
	}
	
	crops := make([]*plant.CropAggregate, len(docs))
	for i, doc := range docs {
		crops[i] = r.cropDocumentToAggregate(&doc)
	}
	
	return crops, nil
}

// Update 更新作物
func (r *MongoCropRepository) Update(crop *plant.CropAggregate) error {
	return r.Save(crop)
}

// Delete 删除作物
func (r *MongoCropRepository) Delete(cropID string) error {
	ctx := context.Background()
	
	filter := bson.M{"crop_id": cropID}
	
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete crop", "error", err, "crop_id", cropID)
		return fmt.Errorf("failed to delete crop: %w", err)
	}
	
	if result.DeletedCount == 0 {
		return plant.ErrCropNotFound
	}
	
	// 清除缓存
	cacheKey := fmt.Sprintf("crop:%s", cropID)
	if err := r.cache.Delete(ctx, cacheKey); err != nil {
		r.logger.Warn("Failed to delete crop cache", "error", err, "crop_id", cropID)
	}
	
	return nil
}

// 私有方法

// cropAggregateToDocument 作物聚合根转文档
func (r *MongoCropRepository) cropAggregateToDocument(crop *plant.CropAggregate) *CropDocument {
	return &CropDocument{
		CropID:               crop.GetID(),
		PlayerID:             crop.GetPlayerID(),
		PlotID:               crop.GetPlotID(),
		SeedID:               crop.GetSeedID(),
		CropType:             string(crop.GetCropType()),
		CurrentStage:         string(crop.GetCurrentStage()),
		GrowthProgress:       crop.GetGrowthProgress(),
		Health:               crop.GetHealth(),
		Moisture:             crop.GetMoisture(),
		Nutrition:            crop.GetNutrition(),
		Quality:              crop.GetQuality(),
		PlantedAt:            crop.GetPlantedAt(),
		LastWatered:          crop.GetLastWatered(),
		LastFertilized:       crop.GetLastFertilized(),
		EstimatedHarvestTime: crop.GetEstimatedHarvestTime(),
		CreatedAt:            crop.GetCreatedAt(),
		UpdatedAt:            crop.GetUpdatedAt(),
	}
}

// cropDocumentToAggregate 作物文档转聚合根
func (r *MongoCropRepository) cropDocumentToAggregate(doc *CropDocument) *plant.CropAggregate {
	// 这里需要根据实际的CropAggregate构造函数来实现
	return plant.ReconstructCropAggregate(
		doc.CropID,
		doc.PlayerID,
		doc.PlotID,
		doc.SeedID,
		plant.CropType(doc.CropType),
		plant.GrowthStage(doc.CurrentStage),
		doc.GrowthProgress,
		doc.Health,
		doc.Moisture,
		doc.Nutrition,
		doc.Quality,
		doc.PlantedAt,
		doc.LastWatered,
		doc.LastFertilized,
		doc.EstimatedHarvestTime,
		doc.CreatedAt,
		doc.UpdatedAt,
	)
}

// CreateIndexes 创建索引
func (r *MongoFarmRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "player_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "level", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}
	
	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		r.logger.Error("Failed to create farm indexes", "error", err)
		return fmt.Errorf("failed to create farm indexes: %w", err)
	}
	
	r.logger.Info("Farm indexes created successfully")
	return nil
}

// CreateIndexes 创建作物索引
func (r *MongoCropRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "crop_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "player_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "plot_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "crop_type", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "current_stage", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "planted_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "estimated_harvest_time", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "player_id", Value: 1}, {Key: "current_stage", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "current_stage", Value: 1}, {Key: "growth_progress", Value: 1}},
		},
	}
	
	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		r.logger.Error("Failed to create crop indexes", "error", err)
		return fmt.Errorf("failed to create crop indexes: %w", err)
	}
	
	r.logger.Info("Crop indexes created successfully")
	return nil
}