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

// WeatherRepository 天气仓储
type WeatherRepository struct {
	collection *mongo.Collection
	logger     logging.Logger
}

// NewWeatherRepository 创建天气仓储
func NewWeatherRepository(db *mongo.Database, logger logging.Logger) *WeatherRepository {
	return &WeatherRepository{
		collection: db.Collection("weather"),
		logger:     logger,
	}
}

// WeatherRecord 天气记录
type WeatherRecord struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Region        string             `bson:"region" json:"region"`
	WeatherType   string             `bson:"weather_type" json:"weather_type"`
	Temperature   int                `bson:"temperature" json:"temperature"`
	Humidity      int                `bson:"humidity" json:"humidity"`
	WindSpeed     int                `bson:"wind_speed" json:"wind_speed"`
	WindDirection string             `bson:"wind_direction" json:"wind_direction"`
	Pressure      int                `bson:"pressure" json:"pressure"`
	Visibility    int                `bson:"visibility" json:"visibility"`
	Description   string             `bson:"description" json:"description"`
	StartTime     time.Time          `bson:"start_time" json:"start_time"`
	EndTime       time.Time          `bson:"end_time" json:"end_time"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

// CreateWeather 创建天气记录
func (r *WeatherRepository) CreateWeather(ctx context.Context, weather *WeatherRecord) error {
	weather.CreatedAt = time.Now()
	weather.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, weather)
	if err != nil {
		r.logger.Error("创建天气记录失败", err, logging.Fields{
			"region":       weather.Region,
			"weather_type": weather.WeatherType,
		})
		return fmt.Errorf("创建天气记录失败: %w", err)
	}

	r.logger.Info("天气记录创建成功", map[string]interface{}{
		"region":       weather.Region,
		"weather_type": weather.WeatherType,
		"temperature":  weather.Temperature,
	})

	return nil
}

// GetWeather 获取天气记录
func (r *WeatherRepository) GetWeather(ctx context.Context, id string) (*WeatherRecord, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的ID格式: %w", err)
	}

	var weather WeatherRecord
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&weather)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("天气记录不存在")
		}
		r.logger.Error("获取天气记录失败", err, logging.Fields{
			"id": id,
		})
		return nil, fmt.Errorf("获取天气记录失败: %w", err)
	}

	return &weather, nil
}

// UpdateWeather 更新天气记录
func (r *WeatherRepository) UpdateWeather(ctx context.Context, id string, updates map[string]interface{}) error {
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
		r.logger.Error("更新天气记录失败", err, logging.Fields{
			"id": id,
		})
		return fmt.Errorf("更新天气记录失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("天气记录不存在")
	}

	r.logger.Info("天气记录更新成功", map[string]interface{}{
		"id": id,
	})

	return nil
}

// DeleteWeather 删除天气记录
func (r *WeatherRepository) DeleteWeather(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		r.logger.Error("删除天气记录失败", err, logging.Fields{
			"id": id,
		})
		return fmt.Errorf("删除天气记录失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("天气记录不存在")
	}

	r.logger.Info("天气记录删除成功", map[string]interface{}{
		"id": id,
	})

	return nil
}

// GetCurrentWeather 获取当前天气
func (r *WeatherRepository) GetCurrentWeather(ctx context.Context, region string) (*WeatherRecord, error) {
	now := time.Now()
	filter := bson.M{
		"region":     region,
		"start_time": bson.M{"$lte": now},
		"end_time":   bson.M{"$gte": now},
	}

	opts := options.FindOne().SetSort(bson.D{{Key: "start_time", Value: -1}})

	var weather WeatherRecord
	err := r.collection.FindOne(ctx, filter, opts).Decode(&weather)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("当前没有天气记录")
		}
		r.logger.Error("获取当前天气失败", err, logging.Fields{
			"region": region,
		})
		return nil, fmt.Errorf("获取当前天气失败: %w", err)
	}

	return &weather, nil
}

// GetWeatherHistory 获取天气历史
func (r *WeatherRepository) GetWeatherHistory(ctx context.Context, region string, startTime, endTime time.Time, limit, offset int) ([]*WeatherRecord, error) {
	filter := bson.M{
		"region":     region,
		"start_time": bson.M{"$gte": startTime, "$lte": endTime},
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "start_time", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("获取天气历史失败", err, logging.Fields{
			"region":     region,
			"start_time": startTime,
			"end_time":   endTime,
		})
		return nil, fmt.Errorf("获取天气历史失败: %w", err)
	}
	defer cursor.Close(ctx)

	var weathers []*WeatherRecord
	if err = cursor.All(ctx, &weathers); err != nil {
		return nil, fmt.Errorf("解析天气历史失败: %w", err)
	}

	r.logger.Info("获取天气历史成功", map[string]interface{}{
		"region": region,
		"count":  len(weathers),
	})

	return weathers, nil
}

// GetWeatherByType 根据天气类型获取记录
func (r *WeatherRepository) GetWeatherByType(ctx context.Context, weatherType string, limit, offset int) ([]*WeatherRecord, error) {
	filter := bson.M{"weather_type": weatherType}
	opts := options.Find().
		SetSort(bson.D{{Key: "start_time", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("根据天气类型获取记录失败", err, logging.Fields{
			"weather_type": weatherType,
		})
		return nil, fmt.Errorf("根据天气类型获取记录失败: %w", err)
	}
	defer cursor.Close(ctx)

	var weathers []*WeatherRecord
	if err = cursor.All(ctx, &weathers); err != nil {
		return nil, fmt.Errorf("解析天气记录失败: %w", err)
	}

	r.logger.Info("根据天气类型获取记录成功", map[string]interface{}{
		"weather_type": weatherType,
		"count":        len(weathers),
	})

	return weathers, nil
}

// GetWeatherStats 获取天气统计
func (r *WeatherRepository) GetWeatherStats(ctx context.Context, region string, startTime, endTime time.Time) (map[string]interface{}, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"region":     region,
				"start_time": bson.M{"$gte": startTime, "$lte": endTime},
			},
		},
		{
			"$group": bson.M{
				"_id":             "$weather_type",
				"count":           bson.M{"$sum": 1},
				"avg_temperature": bson.M{"$avg": "$temperature"},
				"avg_humidity":    bson.M{"$avg": "$humidity"},
				"avg_wind_speed":  bson.M{"$avg": "$wind_speed"},
			},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Error("获取天气统计失败", err, logging.Fields{
			"region":     region,
			"start_time": startTime,
			"end_time":   endTime,
		})
		return nil, fmt.Errorf("获取天气统计失败: %w", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("解析天气统计失败: %w", err)
	}

	stats := make(map[string]interface{})
	for _, result := range results {
		weatherType := result["_id"].(string)
		stats[weatherType] = map[string]interface{}{
			"count":           result["count"],
			"avg_temperature": result["avg_temperature"],
			"avg_humidity":    result["avg_humidity"],
			"avg_wind_speed":  result["avg_wind_speed"],
		}
	}

	r.logger.Info("获取天气统计成功", map[string]interface{}{
		"region": region,
		"stats":  stats,
	})

	return stats, nil
}

// GetWeatherForecast 获取天气预报
func (r *WeatherRepository) GetWeatherForecast(ctx context.Context, region string, days int) ([]*WeatherRecord, error) {
	startTime := time.Now()
	endTime := startTime.AddDate(0, 0, days)

	filter := bson.M{
		"region":     region,
		"start_time": bson.M{"$gte": startTime, "$lte": endTime},
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "start_time", Value: 1}}).
		SetLimit(int64(days * 24)) // 假设每小时一条记录
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("获取天气预报失败", err, logging.Fields{
			"region": region,
			"days":   days,
		})
		return nil, fmt.Errorf("获取天气预报失败: %w", err)
	}
	defer cursor.Close(ctx)

	var weathers []*WeatherRecord
	if err = cursor.All(ctx, &weathers); err != nil {
		return nil, fmt.Errorf("解析天气预报失败: %w", err)
	}

	r.logger.Info("获取天气预报成功", map[string]interface{}{
		"region": region,
		"days":   days,
		"count":  len(weathers),
	})

	return weathers, nil
}
