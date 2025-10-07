package persistence

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"greatestworks/internal/domain/scene/weather"
	"greatestworks/internal/infrastructure/cache"
	"greatestworks/internal/infrastructure/logger"
)

// MongoWeatherRepository MongoDB天气仓储实现
type MongoWeatherRepository struct {
	db         *mongo.Database
	cache      cache.Cache
	logger     logger.Logger
	collection *mongo.Collection
}

// NewMongoWeatherRepository 创建MongoDB天气仓储
func NewMongoWeatherRepository(db *mongo.Database, cache cache.Cache, logger logger.Logger) weather.WeatherRepository {
	return &MongoWeatherRepository{
		db:         db,
		cache:      cache,
		logger:     logger,
		collection: db.Collection("weather"),
	}
}

// WeatherDocument MongoDB天气文档结构
type WeatherDocument struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	WeatherID   string             `bson:"weather_id"`
	RegionID    string             `bson:"region_id"`
	WeatherType string             `bson:"weather_type"`
	Intensity   float64            `bson:"intensity"`
	Temperature float64            `bson:"temperature"`
	Humidity    float64            `bson:"humidity"`
	WindSpeed   float64            `bson:"wind_speed"`
	Visibility  float64            `bson:"visibility"`
	StartTime   time.Time          `bson:"start_time"`
	EndTime     time.Time          `bson:"end_time"`
	Duration    int64              `bson:"duration"` // 秒数
	IsSpecial   bool               `bson:"is_special"`
	Description string             `bson:"description"`
	Effects     []WeatherEffect    `bson:"effects"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

// WeatherEffect 天气影响
type WeatherEffect struct {
	EffectType string  `bson:"effect_type"`
	TargetType string  `bson:"target_type"`
	Modifier   float64 `bson:"modifier"`
	Duration   int64   `bson:"duration"` // 秒数
}

// Save 保存天气记录
func (r *MongoWeatherRepository) Save(weatherAggregate *weather.WeatherAggregate) error {
	ctx := context.Background()
	doc := r.aggregateToDocument(weatherAggregate)
	doc.UpdatedAt = time.Now()

	if doc.ID.IsZero() {
		doc.CreatedAt = time.Now()
		result, err := r.collection.InsertOne(ctx, doc)
		if err != nil {
			r.logger.Error("Failed to insert weather", "error", err, "weather_id", weatherAggregate.GetID())
			return fmt.Errorf("failed to insert weather: %w", err)
		}

		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
			doc.ID = oid
		}
	} else {
		filter := bson.M{"weather_id": weatherAggregate.GetID()}
		update := bson.M{"$set": doc}

		_, err := r.collection.UpdateOne(ctx, filter, update)
		if err != nil {
			r.logger.Error("Failed to update weather", "error", err, "weather_id", weatherAggregate.GetID())
			return fmt.Errorf("failed to update weather: %w", err)
		}
	}

	// 更新缓存
	cacheKey := fmt.Sprintf("weather:%s", weatherAggregate.GetID())
	if err := r.cache.Set(ctx, cacheKey, weatherAggregate, time.Hour); err != nil {
		r.logger.Warn("Failed to cache weather", "error", err, "weather_id", weatherAggregate.GetID())
	}

	return nil
}

// FindByID 根据ID查找天气记录
func (r *MongoWeatherRepository) FindByID(weatherID string) (*weather.WeatherAggregate, error) {
	ctx := context.Background()

	// 先从缓存获取
	cacheKey := fmt.Sprintf("weather:%s", weatherID)
	var cachedWeather *weather.WeatherAggregate
	if err := r.cache.Get(ctx, cacheKey, &cachedWeather); err == nil && cachedWeather != nil {
		return cachedWeather, nil
	}

	// 从数据库获取
	filter := bson.M{"weather_id": weatherID}
	var doc WeatherDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, weather.ErrWeatherNotFound
		}
		r.logger.Error("Failed to find weather", "error", err, "weather_id", weatherID)
		return nil, fmt.Errorf("failed to find weather: %w", err)
	}

	weatherAggregate := r.documentToAggregate(&doc)

	// 更新缓存
	if err := r.cache.Set(ctx, cacheKey, weatherAggregate, time.Hour); err != nil {
		r.logger.Warn("Failed to cache weather", "error", err, "weather_id", weatherID)
	}

	return weatherAggregate, nil
}

// FindCurrentByRegion 查找区域当前天气
func (r *MongoWeatherRepository) FindCurrentByRegion(regionID string) (*weather.WeatherAggregate, error) {
	ctx := context.Background()

	// 先从缓存获取
	cacheKey := fmt.Sprintf("weather:current:%s", regionID)
	var cachedWeather *weather.WeatherAggregate
	if err := r.cache.Get(ctx, cacheKey, &cachedWeather); err == nil && cachedWeather != nil {
		return cachedWeather, nil
	}

	// 从数据库获取当前时间的天气
	now := time.Now()
	filter := bson.M{
		"region_id":  regionID,
		"start_time": bson.M{"$lte": now},
		"end_time":   bson.M{"$gte": now},
	}

	var doc WeatherDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, weather.ErrWeatherNotFound
		}
		r.logger.Error("Failed to find current weather", "error", err, "region_id", regionID)
		return nil, fmt.Errorf("failed to find current weather: %w", err)
	}

	weatherAggregate := r.documentToAggregate(&doc)

	// 更新缓存（较短时间，因为天气会变化）
	if err := r.cache.Set(ctx, cacheKey, weatherAggregate, time.Minute*10); err != nil {
		r.logger.Warn("Failed to cache current weather", "error", err, "region_id", regionID)
	}

	return weatherAggregate, nil
}

// FindByRegionAndTimeRange 根据区域和时间范围查找天气
func (r *MongoWeatherRepository) FindByRegionAndTimeRange(regionID string, startTime, endTime time.Time) ([]*weather.WeatherAggregate, error) {
	ctx := context.Background()

	filter := bson.M{
		"region_id": regionID,
		"$or": []bson.M{
			{
				"start_time": bson.M{
					"$gte": startTime,
					"$lte": endTime,
				},
			},
			{
				"end_time": bson.M{
					"$gte": startTime,
					"$lte": endTime,
				},
			},
			{
				"start_time": bson.M{"$lte": startTime},
				"end_time":   bson.M{"$gte": endTime},
			},
		},
	}

	opts := options.Find().SetSort(bson.D{{Key: "start_time", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to find weather by time range", "error", err, "region_id", regionID)
		return nil, fmt.Errorf("failed to find weather by time range: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []WeatherDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode weather by time range", "error", err, "region_id", regionID)
		return nil, fmt.Errorf("failed to decode weather by time range: %w", err)
	}

	weathers := make([]*weather.WeatherAggregate, len(docs))
	for i, doc := range docs {
		weathers[i] = r.documentToAggregate(&doc)
	}

	return weathers, nil
}

// FindAllCurrent 查找所有区域的当前天气
func (r *MongoWeatherRepository) FindAllCurrent() ([]*weather.WeatherAggregate, error) {
	ctx := context.Background()

	// 先从缓存获取
	cacheKey := "weather:all:current"
	var cachedWeathers []*weather.WeatherAggregate
	if err := r.cache.Get(ctx, cacheKey, &cachedWeathers); err == nil && len(cachedWeathers) > 0 {
		return cachedWeathers, nil
	}

	// 从数据库获取
	now := time.Now()
	filter := bson.M{
		"start_time": bson.M{"$lte": now},
		"end_time":   bson.M{"$gte": now},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find all current weather", "error", err)
		return nil, fmt.Errorf("failed to find all current weather: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []WeatherDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode all current weather", "error", err)
		return nil, fmt.Errorf("failed to decode all current weather: %w", err)
	}

	weathers := make([]*weather.WeatherAggregate, len(docs))
	for i, doc := range docs {
		weathers[i] = r.documentToAggregate(&doc)
	}

	// 更新缓存
	if err := r.cache.Set(ctx, cacheKey, weathers, time.Minute*5); err != nil {
		r.logger.Warn("Failed to cache all current weather", "error", err)
	}

	return weathers, nil
}

// Update 更新天气记录
func (r *MongoWeatherRepository) Update(weatherAggregate *weather.WeatherAggregate) error {
	return r.Save(weatherAggregate)
}

// Delete 删除天气记录
func (r *MongoWeatherRepository) Delete(ctx context.Context, weatherID string) error {

	filter := bson.M{"weather_id": weatherID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete weather", "error", err, "weather_id", weatherID)
		return fmt.Errorf("failed to delete weather: %w", err)
	}

	if result.DeletedCount == 0 {
		return weather.ErrWeatherNotFound
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("weather:%s", weatherID)
	if err := r.cache.Delete(ctx, cacheKey); err != nil {
		r.logger.Warn("Failed to delete weather cache", "error", err, "weather_id", weatherID)
	}

	return nil
}

// DeleteBatch 批量删除天气记录
func (r *MongoWeatherRepository) DeleteBatch(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	filter := bson.M{"weather_id": bson.M{"$in": ids}}

	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete weather batch", "error", err, "ids", ids)
		return fmt.Errorf("failed to delete weather batch: %w", err)
	}

	// 批量清除缓存
	cacheKeys := make([]string, len(ids))
	for i, id := range ids {
		cacheKeys[i] = fmt.Sprintf("weather:%s", id)
	}
	if err := r.cache.DeleteBatch(ctx, cacheKeys); err != nil {
		r.logger.Warn("Failed to delete weather cache batch", "error", err, "ids", ids)
	}

	r.logger.Info("Batch deleted weather records", "deleted_count", result.DeletedCount, "requested_count", len(ids))
	return nil
}

// FindByWeatherType 根据天气类型查找
func (r *MongoWeatherRepository) FindByWeatherType(weatherType weather.WeatherType, limit int) ([]*weather.WeatherAggregate, error) {
	ctx := context.Background()

	filter := bson.M{"weather_type": string(weatherType)}
	opts := options.Find().
		SetSort(bson.D{{Key: "start_time", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to find weather by type", "error", err, "weather_type", weatherType)
		return nil, fmt.Errorf("failed to find weather by type: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []WeatherDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode weather by type", "error", err, "weather_type", weatherType)
		return nil, fmt.Errorf("failed to decode weather by type: %w", err)
	}

	weathers := make([]*weather.WeatherAggregate, len(docs))
	for i, doc := range docs {
		weathers[i] = r.documentToAggregate(&doc)
	}

	return weathers, nil
}

// FindActiveWeather 查找活跃天气
func (r *MongoWeatherRepository) FindActiveWeather(ctx context.Context, sceneID string) (*weather.WeatherAggregate, error) {
	// 先从缓存获取
	cacheKey := fmt.Sprintf("weather:active:%s", sceneID)
	var cachedWeather *weather.WeatherAggregate
	if err := r.cache.Get(ctx, cacheKey, &cachedWeather); err == nil && cachedWeather != nil {
		return cachedWeather, nil
	}

	// 从数据库获取当前时间的活跃天气
	now := time.Now()
	filter := bson.M{
		"scene_id":   sceneID,
		"start_time": bson.M{"$lte": now},
		"end_time":   bson.M{"$gte": now},
		"is_active":  true,
	}

	var doc WeatherDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 没有找到活跃天气
		}
		r.logger.Error("Failed to find active weather", "error", err, "scene_id", sceneID)
		return nil, fmt.Errorf("failed to find active weather: %w", err)
	}

	weatherAggregate := r.documentToAggregate(&doc)

	// 更新缓存
	if err := r.cache.Set(ctx, cacheKey, weatherAggregate, 5*time.Minute); err != nil {
		r.logger.Warn("Failed to cache active weather", "error", err, "scene_id", sceneID)
	}

	return weatherAggregate, nil
}

// FindSpecialWeather 查找特殊天气
func (r *MongoWeatherRepository) FindSpecialWeather(limit int) ([]*weather.WeatherAggregate, error) {
	ctx := context.Background()

	filter := bson.M{"is_special": true}
	opts := options.Find().
		SetSort(bson.D{{Key: "start_time", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to find special weather", "error", err)
		return nil, fmt.Errorf("failed to find special weather: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []WeatherDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode special weather", "error", err)
		return nil, fmt.Errorf("failed to decode special weather: %w", err)
	}

	weathers := make([]*weather.WeatherAggregate, len(docs))
	for i, doc := range docs {
		weathers[i] = r.documentToAggregate(&doc)
	}

	return weathers, nil
}

// Count 计数查询
func (r *MongoWeatherRepository) Count() (int64, error) {
	ctx := context.Background()

	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		r.logger.Error("Failed to count weather", "error", err)
		return 0, fmt.Errorf("failed to count weather: %w", err)
	}

	return count, nil
}

// CountByRegion 根据区域计数
func (r *MongoWeatherRepository) CountByRegion(regionID string) (int64, error) {
	ctx := context.Background()

	filter := bson.M{"region_id": regionID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to count weather by region", "error", err, "region_id", regionID)
		return 0, fmt.Errorf("failed to count weather by region: %w", err)
	}

	return count, nil
}

// CountByType 根据类型计数
func (r *MongoWeatherRepository) CountByType(weatherType weather.WeatherType) (int64, error) {
	ctx := context.Background()

	filter := bson.M{"weather_type": string(weatherType)}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to count weather by type", "error", err, "weather_type", weatherType)
		return 0, fmt.Errorf("failed to count weather by type: %w", err)
	}

	return count, nil
}

// 私有方法

// aggregateToDocument 聚合根转文档
func (r *MongoWeatherRepository) aggregateToDocument(weatherAggregate *weather.WeatherAggregate) *WeatherDocument {
	effects := make([]WeatherEffect, 0)
	for _, effect := range weatherAggregate.GetEffects() {
		effects = append(effects, WeatherEffect{
			EffectType: string(effect.GetEffectType()),
			TargetType: string(effect.GetTargetType()),
			Modifier:   effect.GetModifier(),
			Duration:   int64(effect.GetDuration().Seconds()),
		})
	}

	return &WeatherDocument{
			WeatherID:   weatherAggregate.GetID(),
			RegionID:    weatherAggregate.GetRegionID(),
			WeatherType: string(weatherAggregate.GetWeatherType()),
			Intensity:   weatherAggregate.GetIntensity().GetMultiplier(),
			Temperature: weatherAggregate.GetTemperature(),
		Humidity:    weatherAggregate.GetHumidity(),
		WindSpeed:   weatherAggregate.GetWindSpeed(),
		Visibility:  weatherAggregate.GetVisibility(),
		StartTime:   weatherAggregate.GetStartTime(),
		EndTime:     weatherAggregate.GetEndTime(),
		Duration:    int64(weatherAggregate.GetDuration().Seconds()),
		IsSpecial:   weatherAggregate.IsSpecialWeather(),
		Description: weatherAggregate.GetDescription(),
		Effects:     effects,
		CreatedAt:   weatherAggregate.GetCreatedAt(),
		UpdatedAt:   weatherAggregate.GetUpdatedAt(),
	}
}

// documentToAggregate 文档转聚合根
func (r *MongoWeatherRepository) documentToAggregate(doc *WeatherDocument) *weather.WeatherAggregate {
	effects := make([]*weather.WeatherEffect, len(doc.Effects))
	for i, effect := range doc.Effects {
		effects[i] = weather.NewWeatherEffect(
			effect.EffectType,
			effect.TargetType,
			effect.Modifier,
			time.Duration(effect.Duration)*time.Second,
		)
	}

	// 这里需要根据实际的WeatherAggregate构造函数来实现
	return weather.ReconstructWeatherAggregate(
		doc.WeatherID,
		doc.RegionID,
		weather.ParseWeatherType(doc.WeatherType),
		doc.Intensity,
		doc.Temperature,
		doc.Humidity,
		doc.WindSpeed,
		doc.Visibility,
		doc.StartTime,
		doc.EndTime,
		time.Duration(doc.Duration)*time.Second,
		doc.IsSpecial,
		doc.Description,
		effects,
		doc.CreatedAt,
		doc.UpdatedAt,
	)
}

// CleanupExpiredWeather 清理过期天气数据
func (r *MongoWeatherRepository) CleanupExpiredWeather(ctx context.Context, beforeTime time.Time) (int64, error) {
	filter := bson.M{
		"updated_at": bson.M{"$lt": beforeTime},
		"is_active":  false,
	}

	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to cleanup expired weather", "error", err, "before_time", beforeTime)
		return 0, fmt.Errorf("failed to cleanup expired weather: %w", err)
	}

	r.logger.Info("Cleaned up expired weather", "deleted_count", result.DeletedCount, "before_time", beforeTime)
	return result.DeletedCount, nil
}

// CleanupOldHistory 清理旧的天气历史记录
func (r *MongoWeatherRepository) CleanupOldHistory(ctx context.Context, sceneID string, keepDays int) (int64, error) {
	cutoffTime := time.Now().AddDate(0, 0, -keepDays)
	filter := bson.M{
		"scene_id":   sceneID,
		"created_at": bson.M{"$lt": cutoffTime},
	}

	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to cleanup old weather history", "error", err, "scene_id", sceneID, "keep_days", keepDays)
		return 0, fmt.Errorf("failed to cleanup old weather history: %w", err)
	}

	r.logger.Info("Cleaned up old weather history", "deleted_count", result.DeletedCount, "scene_id", sceneID, "keep_days", keepDays)
	return result.DeletedCount, nil
}

// CreateIndexes 创建索引
func (r *MongoWeatherRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{"scene_id", 1}},
		},
		{
			Keys: bson.D{{"region", 1}},
		},
		{
			Keys: bson.D{{"weather_type", 1}},
		},
		{
			Keys: bson.D{{"start_time", 1}},
		},
		{
			Keys: bson.D{{"end_time", 1}},
		},
		{
			Keys: bson.D{{"is_active", 1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
