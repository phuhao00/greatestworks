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
	"greatestworks/internal/infrastructure/logger"
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

// Save 保存农场聚合根
func (r *MongoFarmRepository) Save(ctx context.Context, farm *plant.FarmAggregate) error {
	doc := r.aggregateToDocument(farm)

	// 使用 upsert 操作
	filter := bson.M{"farm_id": farm.GetFarmID()}
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		r.logger.Error("Failed to save farm", "error", err, "farm_id", farm.GetFarmID())
		return fmt.Errorf("failed to save farm: %w", err)
	}

	// 缓存农场数据
	cacheKey := fmt.Sprintf("farm:%s", farm.GetFarmID())
	if err := r.cache.Set(ctx, cacheKey, farm, 10*time.Minute); err != nil {
		r.logger.Warn("Failed to cache farm", "error", err, "farm_id", farm.GetFarmID())
	}

	// 缓存玩家农场列表
	playerCacheKey := fmt.Sprintf("player_farms:%s", farm.GetOwner())
	if err := r.cache.Delete(ctx, playerCacheKey); err != nil {
		r.logger.Warn("Failed to invalidate player farms cache", "error", err, "player_id", farm.GetOwner())
	}

	r.logger.Info("Farm saved successfully", "farm_id", farm.GetFarmID(), "player_id", farm.GetOwner())
	return nil
}

// SaveBatch 批量保存农场聚合根
func (r *MongoFarmRepository) SaveBatch(ctx context.Context, farms []*plant.FarmAggregate) error {
	if len(farms) == 0 {
		return nil
	}

	var operations []mongo.WriteModel
	for _, farm := range farms {
		doc := r.aggregateToDocument(farm)
		filter := bson.M{"farm_id": farm.GetFarmID()}
		update := bson.M{"$set": doc}
		upsertModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		operations = append(operations, upsertModel)
	}

	_, err := r.collection.BulkWrite(ctx, operations)
	if err != nil {
		r.logger.Error("Failed to save farms in batch", "error", err, "count", len(farms))
		return fmt.Errorf("failed to save farms in batch: %w", err)
	}

	r.logger.Info("Farms saved in batch successfully", "count", len(farms))
	return nil
}

// UpdateBatch 批量更新农场聚合根
func (r *MongoFarmRepository) UpdateBatch(ctx context.Context, farms []*plant.FarmAggregate) error {
	if len(farms) == 0 {
		return nil
	}

	var operations []mongo.WriteModel
	for _, farm := range farms {
		doc := r.aggregateToDocument(farm)
		filter := bson.M{"farm_id": farm.GetFarmID()}
		update := bson.M{"$set": doc}
		updateModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
		operations = append(operations, updateModel)
	}

	_, err := r.collection.BulkWrite(ctx, operations)
	if err != nil {
		r.logger.Error("Failed to update farms in batch", "error", err, "count", len(farms))
		return fmt.Errorf("failed to update farms in batch: %w", err)
	}

	r.logger.Info("Farms updated in batch successfully", "count", len(farms))
	return nil
}

// FindByOwner 根据所有者查找农场
func (r *MongoFarmRepository) FindByOwner(ctx context.Context, owner string) ([]*plant.FarmAggregate, error) {
	filter := bson.M{"player_id": owner}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find farms by owner", "error", err, "owner", owner)
		return nil, fmt.Errorf("failed to find farms by owner: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []FarmDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode farms by owner", "error", err, "owner", owner)
		return nil, fmt.Errorf("failed to decode farms by owner: %w", err)
	}

	farms := make([]*plant.FarmAggregate, len(docs))
	for i, doc := range docs {
		farms[i] = r.documentToAggregate(&doc)
	}

	return farms, nil
}

// FindActiveByOwner 根据所有者查找活跃农场
func (r *MongoFarmRepository) FindActiveByOwner(ctx context.Context, owner string) ([]*plant.FarmAggregate, error) {
	// For now, assume all farms are active. In the future, you might add an "active" field
	return r.FindByOwner(ctx, owner)
}

// FindByPlayer 根据玩家查找农场 (保持向后兼容)
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

// FindByID 根据ID查找农场
func (r *MongoFarmRepository) FindByID(ctx context.Context, farmID string) (*plant.FarmAggregate, error) {
	filter := bson.M{"farm_id": farmID}
	var doc FarmDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to find farm by ID", "error", err, "farm_id", farmID)
		return nil, fmt.Errorf("failed to find farm by ID: %w", err)
	}
	return r.documentToAggregate(&doc), nil
}

// FindBySceneID 根据场景ID查找农场
func (r *MongoFarmRepository) FindBySceneID(ctx context.Context, sceneID string) ([]*plant.FarmAggregate, error) {
	filter := bson.M{"scene_id": sceneID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find farms by scene ID", "error", err, "scene_id", sceneID)
		return nil, fmt.Errorf("failed to find farms by scene ID: %w", err)
	}
	defer cursor.Close(ctx)

	var farms []*plant.FarmAggregate
	for cursor.Next(ctx) {
		var doc FarmDocument
		if err := cursor.Decode(&doc); err != nil {
			r.logger.Error("Failed to decode farm document", "error", err)
			continue
		}
		farms = append(farms, r.documentToAggregate(&doc))
	}

	return farms, nil
}

// FindByClimateZone 根据气候区域查找农场
func (r *MongoFarmRepository) FindByClimateZone(ctx context.Context, climateZone string, limit int) ([]*plant.FarmAggregate, error) {
	// 这里假设气候区域信息存储在场景中，我们通过场景ID来查找
	// 实际实现可能需要根据具体的数据模型调整
	filter := bson.M{"climate_zone": climateZone}
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to find farms by climate zone", "error", err, "climate_zone", climateZone)
		return nil, fmt.Errorf("failed to find farms by climate zone: %w", err)
	}
	defer cursor.Close(ctx)

	var farms []*plant.FarmAggregate
	for cursor.Next(ctx) {
		var doc FarmDocument
		if err := cursor.Decode(&doc); err != nil {
			r.logger.Error("Failed to decode farm document", "error", err)
			continue
		}
		farms = append(farms, r.documentToAggregate(&doc))
	}

	return farms, nil
}

// FindBySize 根据农场大小查找农场
func (r *MongoFarmRepository) FindBySize(ctx context.Context, size plant.FarmSize, limit int) ([]*plant.FarmAggregate, error) {
	filter := bson.M{"size": size}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetLimit(int64(limit)))
	if err != nil {
		r.logger.Error("Failed to find farms by size", "error", err, "size", size)
		return nil, fmt.Errorf("failed to find farms by size: %w", err)
	}
	defer cursor.Close(ctx)

	var farms []*plant.FarmAggregate
	for cursor.Next(ctx) {
		var doc FarmDocument
		if err := cursor.Decode(&doc); err != nil {
			r.logger.Error("Failed to decode farm document", "error", err)
			continue
		}

		farm := r.documentToAggregate(&doc)
		farms = append(farms, farm)
	}

	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error while finding farms by size", "error", err)
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	r.logger.Info("Found farms by size", "size", size, "count", len(farms))
	return farms, nil
}

// FindByStatus 根据农场状态查找农场
func (r *MongoFarmRepository) FindByStatus(ctx context.Context, status plant.FarmStatus, limit int) ([]*plant.FarmAggregate, error) {
	filter := bson.M{"status": status}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetLimit(int64(limit)))
	if err != nil {
		r.logger.Error("Failed to find farms by status", "error", err, "status", status)
		return nil, fmt.Errorf("failed to find farms by status: %w", err)
	}
	defer cursor.Close(ctx)

	var farms []*plant.FarmAggregate
	for cursor.Next(ctx) {
		var doc FarmDocument
		if err := cursor.Decode(&doc); err != nil {
			r.logger.Error("Failed to decode farm document", "error", err)
			continue
		}

		farm := r.documentToAggregate(&doc)
		farms = append(farms, farm)
	}

	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error while finding farms by status", "error", err)
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	r.logger.Info("Found farms by status", "status", status, "count", len(farms))
	return farms, nil
}

// GetFarmCount 获取农场总数
func (r *MongoFarmRepository) GetFarmCount(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		r.logger.Error("Failed to get farm count", "error", err)
		return 0, fmt.Errorf("failed to get farm count: %w", err)
	}

	r.logger.Info("Retrieved farm count", "count", count)
	return count, nil
}

// GetFarmCountByOwner 根据所有者获取农场数量
func (r *MongoFarmRepository) GetFarmCountByOwner(ctx context.Context, owner string) (int64, error) {
	filter := bson.M{"owner": owner}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to get farm count by owner", "error", err, "owner", owner)
		return 0, fmt.Errorf("failed to get farm count by owner: %w", err)
	}

	r.logger.Info("Retrieved farm count by owner", "owner", owner, "count", count)
	return count, nil
}

// GetFarmStatistics 获取农场统计信息
func (r *MongoFarmRepository) GetFarmStatistics(ctx context.Context, farmID string) (*plant.FarmStatistics, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("farm_stats:%s", farmID)
	var stats *plant.FarmStatistics
	if err := r.cache.Get(ctx, cacheKey, &stats); err == nil && stats != nil {
		return stats, nil
	}

	// 从数据库获取农场
	farm, err := r.FindByID(ctx, farmID)
	if err != nil {
		return nil, fmt.Errorf("failed to find farm: %w", err)
	}

	// 获取农场统计信息
	stats = farm.GetStatistics()
	if stats == nil {
		// 如果农场没有统计信息，创建默认的
		stats = plant.NewFarmStatistics()
		stats.UpdatedAt = time.Now()
	}

	// 缓存结果
	if err := r.cache.Set(ctx, cacheKey, stats, 5*time.Minute); err != nil {
		r.logger.Warn("Failed to cache farm statistics", "error", err, "farmID", farmID)
	}

	return stats, nil
}

// GetOwnerStatistics 获取所有者统计信息
func (r *MongoFarmRepository) GetOwnerStatistics(ctx context.Context, owner string) (*plant.OwnerStatistics, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("owner_stats:%s", owner)
	var stats *plant.OwnerStatistics
	if err := r.cache.Get(ctx, cacheKey, &stats); err == nil && stats != nil {
		return stats, nil
	}

	// 从数据库获取农场数量
	filter := bson.M{"owner": owner}
	totalFarms, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count farms: %w", err)
	}

	// 创建统计信息
	stats = &plant.OwnerStatistics{
		Owner:      owner,
		TotalFarms: int(totalFarms),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 缓存结果
	if err := r.cache.Set(ctx, cacheKey, stats, 5*time.Minute); err != nil {
		r.logger.Warn("Failed to cache owner statistics", "error", err, "owner", owner)
	}

	return stats, nil
}

// GetTopFarmsByValue 根据价值获取顶级农场
func (r *MongoFarmRepository) GetTopFarmsByValue(ctx context.Context, limit int) ([]*plant.FarmRanking, error) {
	// TODO: Implement farm value calculation logic
	// For now, return empty slice
	return []*plant.FarmRanking{}, nil
}

// GetTopFarmsByYield 获取按产量排序的顶级农场
func (r *MongoFarmRepository) GetTopFarmsByYield(ctx context.Context, period time.Duration, limit int) ([]*plant.FarmRanking, error) {
	// TODO: 实现农场产量排序逻辑
	return []*plant.FarmRanking{}, nil
}

// GetTopFarmsByProductivity 根据生产力获取顶级农场
func (r *MongoFarmRepository) GetTopFarmsByProductivity(ctx context.Context, limit int) ([]*plant.FarmRanking, error) {
	// 创建聚合管道来计算生产力并排序
	pipeline := []bson.M{
		{
			"$addFields": bson.M{
				"productivity": bson.M{
					"$multiply": []interface{}{
						"$total_yield",
						bson.M{"$ifNull": []interface{}{"$efficiency_bonus", 1.0}},
					},
				},
			},
		},
		{
			"$sort": bson.M{"productivity": -1},
		},
		{
			"$limit": limit,
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate farms by productivity: %w", err)
	}
	defer cursor.Close(ctx)

	var rankings []*plant.FarmRanking
	for cursor.Next(ctx) {
		var doc FarmDocument
		if err := cursor.Decode(&doc); err != nil {
			r.logger.Warn("Failed to decode farm document", "error", err)
			continue
		}

		farm := r.documentToAggregate(&doc)
		// Create FarmRanking from farm aggregate
		ranking := &plant.FarmRanking{
			Rank:        len(rankings) + 1,
			FarmID:      farm.GetFarmID(),
			Owner:       farm.GetOwner(),
			FarmName:    "",  // TODO: Add farm name to aggregate
			Score:       0.0, // TODO: Calculate actual productivity score
			Metric:      "productivity",
			Value:       0.0,
			LastUpdated: time.Now(),
		}
		rankings = append(rankings, ranking)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return rankings, nil
}

// Update 更新农场
func (r *MongoFarmRepository) Update(ctx context.Context, farm *plant.FarmAggregate) error {
	return r.Save(ctx, farm)
}

// Delete 删除农场
func (r *MongoFarmRepository) Delete(ctx context.Context, farmID string) error {
	filter := bson.M{"farm_id": farmID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete farm", "error", err, "farm_id", farmID)
		return fmt.Errorf("failed to delete farm: %w", err)
	}

	if result.DeletedCount == 0 {
		return plant.ErrFarmNotFound
	}

	// 清除缓存 - Note: We need playerID to clear cache, but it's not available in this method
	// This is a design issue that should be addressed
	cacheKey := fmt.Sprintf("farm:%s", farmID) // Using farmID as fallback
	if err := r.cache.Delete(ctx, cacheKey); err != nil {
		r.logger.Warn("Failed to delete farm cache", "error", err, "farm_id", farmID)
	}

	return nil
}

// 私有方法

// stringToSoilType 字符串转土壤类型
func stringToSoilType(s string) plant.SoilType {
	switch s {
	case "sandy":
		return plant.SoilTypeSandy
	case "clay":
		return plant.SoilTypeClay
	case "loam":
		return plant.SoilTypeLoam
	case "silt":
		return plant.SoilTypeSilt
	case "peat":
		return plant.SoilTypePeat
	case "chalk":
		return plant.SoilTypeChalk
	default:
		return plant.SoilTypeLoam // 默认壤土
	}
}

// aggregateToDocument 聚合根转文档
func (r *MongoFarmRepository) aggregateToDocument(farm *plant.FarmAggregate) *FarmDocument {
	plots := make([]PlotDocument, 0)
	for _, plot := range farm.GetPlots() {
		plotDoc := PlotDocument{
			ID:       plot.GetID(),
			SoilType: plot.GetSoilType().String(),
			// Use default values for fertility and moisture since Plot doesn't have these methods
			Fertility: 50.0, // Default fertility
			Moisture:  50.0, // Default moisture
		}

		// If plot has a crop, get its ID
		if plot.HasCrop() {
			crop := plot.GetCrop()
			if crop != nil {
				plotDoc.CropID = crop.GetID()
			}
		}

		plots = append(plots, plotDoc)
	}

	return &FarmDocument{
		FarmID:     farm.GetFarmID(),
		PlayerID:   farm.GetOwner(),
		Level:      1, // Default level, as FarmAggregate doesn't have level concept
		Experience: 0, // Default experience, as FarmAggregate doesn't have experience concept
		Plots:      plots,
		CreatedAt:  farm.GetCreatedAt(),
		UpdatedAt:  farm.GetUpdatedAt(),
	}
}

// documentToAggregate 文档转聚合根
func (r *MongoFarmRepository) documentToAggregate(doc *FarmDocument) *plant.FarmAggregate {
	// Create farm using NewFarmAggregate constructor
	farm := plant.NewFarmAggregate(
		doc.FarmID,
		"",                  // sceneID - not stored in document, using empty string
		doc.PlayerID,        // owner
		"",                  // name - will be set later if needed
		plant.FarmSizeSmall, // default size - should be stored in document in future
	)

	// Add plots to the farm
	for _, plotDoc := range doc.Plots {
		plot := plant.NewPlot(
			plotDoc.ID,
			"",                  // name - not stored in document
			plant.PlotSizeSmall, // default size - should be stored in document
			stringToSoilType(plotDoc.SoilType),
		)
		// Note: We can't set fertility, moisture, cropID directly as they're not exposed
		// This is a limitation of the current domain model design
		farm.AddPlot(plot)
	}

	return farm
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
func (r *MongoCropRepository) Save(ctx context.Context, crop *plant.Crop) error {
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
func (r *MongoCropRepository) FindByID(ctx context.Context, cropID string) (*plant.Crop, error) {

	// 先从缓存获取
	cacheKey := fmt.Sprintf("crop:%s", cropID)
	var cachedCrop *plant.Crop
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
func (r *MongoCropRepository) FindByPlayer(playerID string) ([]*plant.Crop, error) {
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

	crops := make([]*plant.Crop, len(docs))
	for i, doc := range docs {
		crops[i] = r.cropDocumentToAggregate(&doc)
	}

	return crops, nil
}

// FindGrowingCrops 查找正在生长的作物
func (r *MongoCropRepository) FindGrowingCrops() ([]*plant.Crop, error) {
	ctx := context.Background()

	filter := bson.M{
		"current_stage":   bson.M{"$ne": "harvested"},
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

	crops := make([]*plant.Crop, len(docs))
	for i, doc := range docs {
		crops[i] = r.cropDocumentToAggregate(&doc)
	}

	return crops, nil
}

// FindHarvestable 查找可收获的作物
func (r *MongoCropRepository) FindHarvestable(ctx context.Context, farmID string) ([]*plant.Crop, error) {
	filter := bson.M{
		"current_stage":   plant.GrowthStageMature.String(),
		"growth_progress": bson.M{"$gte": 100.0},
	}

	// 如果指定了农场ID，添加到过滤条件
	if farmID != "" {
		filter["plot_id"] = bson.M{"$regex": fmt.Sprintf("^%s", farmID)}
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find harvestable crops", "error", err, "farm_id", farmID)
		return nil, fmt.Errorf("failed to find harvestable crops: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []CropDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode harvestable crops", "error", err)
		return nil, fmt.Errorf("failed to decode harvestable crops: %w", err)
	}

	var crops []*plant.Crop
	for _, doc := range docs {
		crop := r.cropDocumentToAggregate(&doc)
		crops = append(crops, crop)
	}

	return crops, nil
}

// FindNeedsCare 查找需要照料的作物
func (r *MongoCropRepository) FindNeedsCare(ctx context.Context, farmID string) ([]*plant.Crop, error) {
	filter := bson.M{
		"farm_id": farmID,
		"$or": []bson.M{
			{"health_score": bson.M{"$lt": 50.0}},
			{"water_level": bson.M{"$lt": 30.0}},
			{"needs_fertilizer": true},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []CropDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	crops := make([]*plant.Crop, len(docs))
	for i, doc := range docs {
		crops[i] = r.cropDocumentToAggregate(&doc)
	}

	return crops, nil
}

// FindByFarmID 根据农场ID查找作物
func (r *MongoCropRepository) FindByFarmID(ctx context.Context, farmID string) ([]*plant.Crop, error) {
	// Since we don't have farm_id directly in crop document, we'll use plot_id pattern
	// This assumes plot_id contains farm information or we need to join with plots
	filter := bson.M{"plot_id": bson.M{"$regex": fmt.Sprintf("^%s", farmID)}}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find crops by farm ID", "error", err, "farm_id", farmID)
		return nil, fmt.Errorf("failed to find crops by farm ID: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []CropDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode crops by farm ID", "error", err)
		return nil, fmt.Errorf("failed to decode crops by farm ID: %w", err)
	}

	crops := make([]*plant.Crop, len(docs))
	for i, doc := range docs {
		crops[i] = r.cropDocumentToAggregate(&doc)
	}

	return crops, nil
}

// FindByGrowthStage 根据生长阶段查找作物
func (r *MongoCropRepository) FindByGrowthStage(ctx context.Context, stage plant.GrowthStage, limit int) ([]*plant.Crop, error) {
	filter := bson.M{"current_stage": stage.String()}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetLimit(int64(limit)))
	if err != nil {
		r.logger.Error("Failed to find crops by growth stage", "error", err, "stage", stage)
		return nil, fmt.Errorf("failed to find crops by growth stage: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []CropDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode crops by growth stage", "error", err)
		return nil, fmt.Errorf("failed to decode crops by growth stage: %w", err)
	}

	crops := make([]*plant.Crop, len(docs))
	for i, doc := range docs {
		crops[i] = r.cropDocumentToAggregate(&doc)
	}

	return crops, nil
}

// FindByHealthRange 根据健康值范围查找作物
func (r *MongoCropRepository) FindByHealthRange(ctx context.Context, minHealth, maxHealth float64, limit int) ([]*plant.Crop, error) {
	filter := bson.M{
		"health": bson.M{
			"$gte": minHealth,
			"$lte": maxHealth,
		},
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetLimit(int64(limit)))
	if err != nil {
		r.logger.Error("Failed to find crops by health range", "error", err, "min_health", minHealth, "max_health", maxHealth)
		return nil, fmt.Errorf("failed to find crops by health range: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []CropDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode crops by health range", "error", err)
		return nil, fmt.Errorf("failed to decode crops by health range: %w", err)
	}

	crops := make([]*plant.Crop, len(docs))
	for i, doc := range docs {
		crops[i] = r.cropDocumentToAggregate(&doc)
	}

	return crops, nil
}

// FindByProgressRange 根据生长进度范围查找作物
func (r *MongoCropRepository) FindByProgressRange(ctx context.Context, minProgress, maxProgress float64, limit int) ([]*plant.Crop, error) {
	filter := bson.M{
		"growth_progress": bson.M{
			"$gte": minProgress,
			"$lte": maxProgress,
		},
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetLimit(int64(limit)))
	if err != nil {
		r.logger.Error("Failed to find crops by progress range", "error", err, "min_progress", minProgress, "max_progress", maxProgress)
		return nil, fmt.Errorf("failed to find crops by progress range: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []CropDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode crops by progress range", "error", err)
		return nil, fmt.Errorf("failed to decode crops by progress range: %w", err)
	}

	crops := make([]*plant.Crop, len(docs))
	for i, doc := range docs {
		crops[i] = r.cropDocumentToAggregate(&doc)
	}

	return crops, nil
}

// FindBySeedType 根据种子类型查找作物
func (r *MongoCropRepository) FindBySeedType(ctx context.Context, seedType plant.SeedType, limit int) ([]*plant.Crop, error) {
	filter := bson.M{"crop_type": seedType.String()}

	findOptions := options.Find()
	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find crops by seed type: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []CropDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode crops: %w", err)
	}

	crops := make([]*plant.Crop, len(docs))
	for i, doc := range docs {
		crops[i] = r.cropDocumentToAggregate(&doc)
	}

	return crops, nil
}

// FindByTimeRange 根据时间范围查找作物
func (r *MongoCropRepository) FindByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*plant.Crop, error) {
	filter := bson.M{
		"created_at": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	findOptions := options.Find().SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find crops by time range: %w", err)
	}
	defer cursor.Close(ctx)

	var crops []*plant.Crop
	for cursor.Next(ctx) {
		var doc CropDocument
		if err := cursor.Decode(&doc); err != nil {
			r.logger.Warn("Failed to decode crop document", "error", err)
			continue
		}

		crop := r.cropDocumentToAggregate(&doc)
		crops = append(crops, crop)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return crops, nil
}

// Update 更新作物
func (r *MongoCropRepository) Update(ctx context.Context, crop *plant.Crop) error {
	return r.Save(ctx, crop)
}

// Delete 删除作物
func (r *MongoCropRepository) Delete(ctx context.Context, cropID string) error {
	filter := bson.M{"crop_id": cropID}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete crop", "error", err, "crop_id", cropID)
		return fmt.Errorf("failed to delete crop: %w", err)
	}
	return nil
}

// FindExpiredCrops 查找过期的作物
func (r *MongoCropRepository) FindExpiredCrops(ctx context.Context, expiredBefore time.Time) ([]*plant.Crop, error) {
	filter := bson.M{
		"planted_at":   bson.M{"$lt": expiredBefore},
		"growth_stage": bson.M{"$ne": plant.GrowthStageMature.String()},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var crops []*plant.Crop
	for cursor.Next(ctx) {
		var doc CropDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}

		crop := r.cropDocumentToAggregate(&doc)
		crops = append(crops, crop)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return crops, nil
}

// DeleteBatch 批量删除作物
func (r *MongoCropRepository) DeleteBatch(ctx context.Context, cropIDs []string) error {
	filter := bson.M{"crop_id": bson.M{"$in": cropIDs}}

	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete crops in batch", "error", err, "crop_ids", cropIDs)
		return fmt.Errorf("failed to delete crops in batch: %w", err)
	}

	// 清除缓存
	for _, cropID := range cropIDs {
		cacheKey := fmt.Sprintf("crop:%s", cropID)
		if err := r.cache.Delete(ctx, cacheKey); err != nil {
			r.logger.Warn("Failed to delete crop cache", "error", err, "crop_id", cropID)
		}
	}

	r.logger.Info("Crops deleted in batch", "count", result.DeletedCount, "crop_ids", cropIDs)
	return nil
}

// GetAverageGrowthProgress 获取平均生长进度
func (r *MongoCropRepository) GetAverageGrowthProgress(ctx context.Context, seedType plant.SeedType) (float64, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"crop_type": seedType.String()}},
		{"$group": bson.M{
			"_id":          nil,
			"avg_progress": bson.M{"$avg": "$growth_progress"},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Error("Failed to get average growth progress", "error", err, "seed_type", seedType.String())
		return 0, fmt.Errorf("failed to get average growth progress: %w", err)
	}
	defer cursor.Close(ctx)

	var result struct {
		AvgProgress float64 `bson:"avg_progress"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			r.logger.Error("Failed to decode average growth progress", "error", err)
			return 0, fmt.Errorf("failed to decode average growth progress: %w", err)
		}
		return result.AvgProgress, nil
	}

	// 如果没有作物，返回0
	return 0, nil
}

// GetCropCountByStage 根据生长阶段获取作物数量
func (r *MongoCropRepository) GetCropCountByStage(ctx context.Context, stage plant.GrowthStage) (int64, error) {
	filter := bson.M{"current_stage": stage.String()}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to get crop count by stage", "error", err, "stage", stage.String())
		return 0, fmt.Errorf("failed to get crop count by stage: %w", err)
	}

	r.logger.Debug("Got crop count by stage", "stage", stage.String(), "count", count)
	return count, nil
}

// GetCropCountByType 根据种子类型获取作物数量
func (r *MongoCropRepository) GetCropCountByType(ctx context.Context, seedType plant.SeedType) (int64, error) {
	filter := bson.M{"crop_type": seedType.String()}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to get crop count by type", "error", err, "seed_type", seedType.String())
		return 0, fmt.Errorf("failed to get crop count by type: %w", err)
	}

	r.logger.Debug("Got crop count by type", "seed_type", seedType.String(), "count", count)
	return count, nil
}

// GetCropStatistics 获取作物统计信息
func (r *MongoCropRepository) GetCropStatistics(ctx context.Context, farmID string) (*plant.CropStatistics, error) {
	// 基础过滤条件
	baseFilter := bson.M{}
	if farmID != "" {
		baseFilter["plot_id"] = bson.M{"$regex": fmt.Sprintf("^%s", farmID)}
	}

	// 获取总作物数量
	totalCrops, err := r.collection.CountDocuments(ctx, baseFilter)
	if err != nil {
		r.logger.Error("Failed to get total crop count", "error", err, "farm_id", farmID)
		return nil, fmt.Errorf("failed to get total crop count: %w", err)
	}

	// 按类型统计作物
	cropsByType := make(map[plant.SeedType]int)
	typePipeline := []bson.M{
		{"$match": baseFilter},
		{"$group": bson.M{
			"_id":   "$crop_type",
			"count": bson.M{"$sum": 1},
		}},
	}

	typeCursor, err := r.collection.Aggregate(ctx, typePipeline)
	if err != nil {
		r.logger.Error("Failed to aggregate crops by type", "error", err)
		return nil, fmt.Errorf("failed to aggregate crops by type: %w", err)
	}
	defer typeCursor.Close(ctx)

	for typeCursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := typeCursor.Decode(&result); err != nil {
			continue
		}
		cropsByType[stringToSeedType(result.ID)] = result.Count
	}

	// 按阶段统计作物
	cropsByStage := make(map[plant.GrowthStage]int)
	stagePipeline := []bson.M{
		{"$match": baseFilter},
		{"$group": bson.M{
			"_id":   "$current_stage",
			"count": bson.M{"$sum": 1},
		}},
	}

	stageCursor, err := r.collection.Aggregate(ctx, stagePipeline)
	if err != nil {
		r.logger.Error("Failed to aggregate crops by stage", "error", err)
		return nil, fmt.Errorf("failed to aggregate crops by stage: %w", err)
	}
	defer stageCursor.Close(ctx)

	for stageCursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := stageCursor.Decode(&result); err != nil {
			continue
		}
		cropsByStage[stringToGrowthStage(result.ID)] = result.Count
	}

	// 计算平均生长进度
	progressPipeline := []bson.M{
		{"$match": baseFilter},
		{"$group": bson.M{
			"_id":          nil,
			"avg_progress": bson.M{"$avg": "$growth_progress"},
			"avg_health":   bson.M{"$avg": "$health"},
		}},
	}

	progressCursor, err := r.collection.Aggregate(ctx, progressPipeline)
	if err != nil {
		r.logger.Error("Failed to aggregate progress stats", "error", err)
		return nil, fmt.Errorf("failed to aggregate progress stats: %w", err)
	}
	defer progressCursor.Close(ctx)

	var avgProgress, avgHealth float64
	if progressCursor.Next(ctx) {
		var result struct {
			AvgProgress float64 `bson:"avg_progress"`
			AvgHealth   float64 `bson:"avg_health"`
		}
		if err := progressCursor.Decode(&result); err == nil {
			avgProgress = result.AvgProgress
			avgHealth = result.AvgHealth
		}
	}

	// 统计可收获作物数量
	harvestableFilter := bson.M{
		"current_stage":   plant.GrowthStageMature.String(),
		"growth_progress": bson.M{"$gte": 100.0},
	}
	for k, v := range baseFilter {
		harvestableFilter[k] = v
	}

	harvestableCount, err := r.collection.CountDocuments(ctx, harvestableFilter)
	if err != nil {
		r.logger.Warn("Failed to count harvestable crops", "error", err)
		harvestableCount = 0
	}

	// 统计需要照料的作物数量（健康度低于50%或水分低于30%）
	needsCareFilter := bson.M{
		"$or": []bson.M{
			{"health": bson.M{"$lt": 50.0}},
			{"moisture": bson.M{"$lt": 30.0}},
		},
	}
	for k, v := range baseFilter {
		needsCareFilter[k] = v
	}

	needsCareCount, err := r.collection.CountDocuments(ctx, needsCareFilter)
	if err != nil {
		r.logger.Warn("Failed to count crops needing care", "error", err)
		needsCareCount = 0
	}

	now := time.Now()
	return &plant.CropStatistics{
		FarmID:                farmID,
		TotalCrops:            int(totalCrops),
		CropsByType:           cropsByType,
		CropsByStage:          cropsByStage,
		AverageGrowthProgress: avgProgress,
		AverageHealthScore:    avgHealth,
		HarvestableCrops:      int(harvestableCount),
		CropsNeedingCare:      int(needsCareCount),
		CreatedAt:             now,
		UpdatedAt:             now,
	}, nil
}

// 私有方法

// cropAggregateToDocument 将聚合转换为文档
func (r *MongoCropRepository) cropAggregateToDocument(crop *plant.Crop) *CropDocument {
	return &CropDocument{
		CropID:               crop.GetID(),
		PlayerID:             crop.GetPlayerID(),
		PlotID:               "", // TODO: Add PlotID to Crop entity
		SeedID:               "", // TODO: Add SeedID to Crop entity
		CropType:             crop.GetSeedType().String(),
		CurrentStage:         crop.GetGrowthStage().String(),
		GrowthProgress:       crop.GetGrowthProgress(),
		Health:               crop.GetHealthPoints(),
		Moisture:             crop.GetWaterLevel(),
		Nutrition:            crop.GetNutrientLevel(),
		Quality:              0.0, // TODO: Add Quality to Crop entity
		PlantedAt:            crop.PlantedTime,
		LastWatered:          crop.LastWateredTime,
		LastFertilized:       crop.LastFertilizedTime,
		EstimatedHarvestTime: crop.ExpectedHarvestTime,
		CreatedAt:            crop.CreatedAt,
		UpdatedAt:            crop.UpdatedAt,
	}
}

// stringToSeedType 字符串转种子类型
func stringToSeedType(s string) plant.SeedType {
	switch s {
	case "wheat":
		return plant.SeedTypeWheat
	case "corn":
		return plant.SeedTypeCorn
	case "rice":
		return plant.SeedTypeRice
	case "potato":
		return plant.SeedTypePotato
	case "carrot":
		return plant.SeedTypeCarrot
	case "tomato":
		return plant.SeedTypeTomato
	default:
		return plant.SeedTypeWheat // 默认小麦
	}
}

// stringToGrowthStage 字符串转生长阶段
func stringToGrowthStage(s string) plant.GrowthStage {
	switch s {
	case "seed":
		return plant.GrowthStageSeed
	case "seedling":
		return plant.GrowthStageSeedling
	case "growing":
		return plant.GrowthStageGrowing
	case "flowering":
		return plant.GrowthStageFlowering
	case "mature":
		return plant.GrowthStageMature
	default:
		return plant.GrowthStageSeed
	}
}

// cropDocumentToAggregate 文档转作物聚合根
func (r *MongoCropRepository) cropDocumentToAggregate(doc *CropDocument) *plant.Crop {
	// Create a basic crop using NewCrop constructor
	crop := plant.NewCrop(
		doc.CropID,
		doc.PlayerID,                   // Player ID from document
		stringToSeedType(doc.CropType), // Convert string to SeedType
		1,                              // Default quantity
		nil,                            // No soil for now
		"",                             // No climate zone for now
	)

	// Manually set the fields from document
	crop.GrowthStage = stringToGrowthStage(doc.CurrentStage)
	crop.GrowthProgress = doc.GrowthProgress
	crop.HealthPoints = doc.Health
	crop.WaterLevel = doc.Moisture
	crop.NutrientLevel = doc.Nutrition
	crop.PlantedTime = doc.PlantedAt
	crop.LastWateredTime = doc.LastWatered
	crop.LastFertilizedTime = doc.LastFertilized
	crop.ExpectedHarvestTime = doc.EstimatedHarvestTime
	crop.CreatedAt = doc.CreatedAt
	crop.UpdatedAt = doc.UpdatedAt

	return crop
}

// CreateIndexes 创建索引
func (r *MongoFarmRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "player_id", Value: 1}},
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

// CleanupExpiredCrops 清理过期作物
func (r *MongoCropRepository) CleanupExpiredCrops(ctx context.Context, beforeTime time.Time) (int64, error) {
	filter := bson.M{
		"updated_at":   bson.M{"$lt": beforeTime},
		"growth_stage": bson.M{"$in": []string{"dead", "withered", "expired"}},
	}

	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to cleanup expired crops", "error", err, "before_time", beforeTime)
		return 0, fmt.Errorf("failed to cleanup expired crops: %w", err)
	}

	r.logger.Info("Cleaned up expired crops", "deleted_count", result.DeletedCount, "before_time", beforeTime)
	return result.DeletedCount, nil
}

// CreateIndexes 创建作物索引
func (r *MongoCropRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "crop_id", Value: 1}},
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

// DeleteBatch 批量删除农场
func (r *MongoFarmRepository) DeleteBatch(ctx context.Context, farmIDs []string) error {
	if len(farmIDs) == 0 {
		return nil
	}

	filter := bson.M{"farm_id": bson.M{"$in": farmIDs}}
	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete farms in batch", "error", err, "farm_ids", farmIDs)
		return fmt.Errorf("failed to delete farms in batch: %w", err)
	}

	r.logger.Info("Batch deleted farms", "deleted_count", result.DeletedCount, "requested_count", len(farmIDs))

	// Clear cache for all deleted farms
	for _, farmID := range farmIDs {
		cacheKey := fmt.Sprintf("farm:%s", farmID)
		if err := r.cache.Delete(ctx, cacheKey); err != nil {
			r.logger.Warn("Failed to delete farm cache", "error", err, "farm_id", farmID)
		}
	}

	return nil
}
