package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"greatestworks/internal/domain/player/hangup"
	"greatestworks/internal/infrastructure/cache"
	"greatestworks/internal/infrastructure/logger"
)

// MongoHangupRepository MongoDB挂机仓储实现
type MongoHangupRepository struct {
	db         *mongo.Database
	cache      cache.Cache
	logger     logger.Logger
	collection *mongo.Collection
}

// HangupDocument MongoDB挂机文档结构
type HangupDocument struct {
	ID         bson.ObjectID `bson:"_id,omitempty"`
	HangupID   string        `bson:"hangup_id"`
	PlayerID   string        `bson:"player_id"`
	LocationID string        `bson:"location_id"`
	StartTime  time.Time     `bson:"start_time"`
	EndTime    *time.Time    `bson:"end_time,omitempty"`
	Duration   int64         `bson:"duration"` // 秒数
	Efficiency float64       `bson:"efficiency"`
	BaseRate   float64       `bson:"base_rate"`
	Status     string        `bson:"status"`
	Rewards    []RewardItem  `bson:"rewards"`
	CreatedAt  time.Time     `bson:"created_at"`
	UpdatedAt  time.Time     `bson:"updated_at"`
}

// RewardItem 奖励项目
type RewardItem struct {
	Type     string `bson:"type"`
	ItemID   string `bson:"item_id"`
	Quantity int64  `bson:"quantity"`
	Quality  string `bson:"quality"`
}

// Save 保存挂机记录
func (r *MongoHangupRepository) Save(hangupAggregate *hangup.HangupAggregate) error {
	ctx := context.Background()
	doc := r.aggregateToDocument(hangupAggregate)
	doc.UpdatedAt = time.Now()

	if doc.ID.IsZero() {
		doc.CreatedAt = time.Now()
		result, err := r.collection.InsertOne(ctx, doc)
		if err != nil {
			r.logger.Error("Failed to insert hangup", "error", err, "player_id", hangupAggregate.GetPlayerID())
			return fmt.Errorf("failed to insert hangup: %w", err)
		}

		if oid, ok := result.InsertedID.(bson.ObjectID); ok {
			doc.ID = oid
		}
	} else {
		filter := bson.M{"player_id": hangupAggregate.GetPlayerID()}
		update := bson.M{"$set": doc}

		_, err := r.collection.UpdateOne(ctx, filter, update)
		if err != nil {
			r.logger.Error("Failed to update hangup", "error", err, "player_id", hangupAggregate.GetPlayerID())
			return fmt.Errorf("failed to update hangup: %w", err)
		}
	}

	// 更新缓存
	cacheKey := fmt.Sprintf("player:%s", hangupAggregate.GetPlayerID())
	if err := r.cache.Set(ctx, cacheKey, hangupAggregate, time.Hour); err != nil {
		r.logger.Warn("Failed to cache hangup", "error", err, "player_id", hangupAggregate.GetPlayerID())
	}

	return nil
}

// FindByID 根据ID查找挂机记录
func (r *MongoHangupRepository) FindByID(hangupID string) (*hangup.HangupAggregate, error) {
	ctx := context.Background()

	// 先从缓存获取
	cacheKey := fmt.Sprintf("hangup:%s", hangupID)
	var cachedHangup *hangup.HangupAggregate
	if err := r.cache.Get(ctx, cacheKey, &cachedHangup); err == nil && cachedHangup != nil {
		return cachedHangup, nil
	}

	// 从数据库获取
	filter := bson.M{"hangup_id": hangupID}
	var doc HangupDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("hangup not found")
		}
		r.logger.Error("Failed to find hangup", "error", err, "hangup_id", hangupID)
		return nil, fmt.Errorf("failed to find hangup: %w", err)
	}

	hangupAggregate := r.documentToAggregate(&doc)

	// 更新缓存
	if err := r.cache.Set(ctx, cacheKey, hangupAggregate, time.Hour); err != nil {
		r.logger.Warn("Failed to cache hangup", "error", err, "hangup_id", hangupID)
	}

	return hangupAggregate, nil
}

// FindActiveByPlayer 查找玩家的活跃挂机记录
func (r *MongoHangupRepository) FindActiveByPlayer(playerID string) (*hangup.HangupAggregate, error) {
	ctx := context.Background()

	// 先从缓存获取
	cacheKey := fmt.Sprintf("hangup:active:%s", playerID)
	var cachedHangup *hangup.HangupAggregate
	if err := r.cache.Get(ctx, cacheKey, &cachedHangup); err == nil && cachedHangup != nil {
		return cachedHangup, nil
	}

	// 从数据库获取
	filter := bson.M{
		"player_id": playerID,
		"status":    "active",
		"end_time":  bson.M{"$exists": false},
	}

	var doc HangupDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 没有活跃的挂机记录
		}
		r.logger.Error("Failed to find active hangup", "error", err, "player_id", playerID)
		return nil, fmt.Errorf("failed to find active hangup: %w", err)
	}

	hangupAggregate := r.documentToAggregate(&doc)

	// 更新缓存
	if err := r.cache.Set(ctx, cacheKey, hangupAggregate, time.Minute*30); err != nil {
		r.logger.Warn("Failed to cache active hangup", "error", err, "player_id", playerID)
	}

	return hangupAggregate, nil
}

// FindHistoryByPlayer 查找玩家的挂机历�?
func (r *MongoHangupRepository) FindHistoryByPlayer(playerID string, limit int) ([]*hangup.HangupAggregate, error) {
	ctx := context.Background()

	filter := bson.M{
		"player_id": playerID,
		"status":    "completed",
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "end_time", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to find hangup history", "error", err, "player_id", playerID)
		return nil, fmt.Errorf("failed to find hangup history: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []HangupDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode hangup history", "error", err, "player_id", playerID)
		return nil, fmt.Errorf("failed to decode hangup history: %w", err)
	}

	hangups := make([]*hangup.HangupAggregate, len(docs))
	for i, doc := range docs {
		hangups[i] = r.documentToAggregate(&doc)
	}

	return hangups, nil
}

// Update 更新挂机记录
func (r *MongoHangupRepository) Update(hangupAggregate *hangup.HangupAggregate) error {
	return r.Save(hangupAggregate)
}

// Delete 删除挂机记录
func (r *MongoHangupRepository) Delete(hangupID string) error {
	ctx := context.Background()

	filter := bson.M{"hangup_id": hangupID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete hangup", "error", err, "hangup_id", hangupID)
		return fmt.Errorf("failed to delete hangup: %w", err)
	}

	if result.DeletedCount == 0 {
		return hangup.ErrHangupNotFound
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("hangup:%s", hangupID)
	if err := r.cache.Delete(ctx, cacheKey); err != nil {
		r.logger.Warn("Failed to delete hangup cache", "error", err, "hangup_id", hangupID)
	}

	return nil
}

// FindByLocation 根据地点查找挂机记录
func (r *MongoHangupRepository) FindByLocation(locationID string, limit int) ([]*hangup.HangupAggregate, error) {
	ctx := context.Background()

	filter := bson.M{"location_id": locationID}
	opts := options.Find().
		SetSort(bson.D{{Key: "start_time", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to find hangups by location", "error", err, "location_id", locationID)
		return nil, fmt.Errorf("failed to find hangups by location: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []HangupDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode hangups by location", "error", err, "location_id", locationID)
		return nil, fmt.Errorf("failed to decode hangups by location: %w", err)
	}

	hangups := make([]*hangup.HangupAggregate, len(docs))
	for i, doc := range docs {
		hangups[i] = r.documentToAggregate(&doc)
	}

	return hangups, nil
}

// FindByTimeRange 根据时间范围查找挂机记录
func (r *MongoHangupRepository) FindByTimeRange(startTime, endTime time.Time) ([]*hangup.HangupAggregate, error) {
	ctx := context.Background()

	filter := bson.M{
		"start_time": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find hangups by time range", "error", err)
		return nil, fmt.Errorf("failed to find hangups by time range: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []HangupDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode hangups by time range", "error", err)
		return nil, fmt.Errorf("failed to decode hangups by time range: %w", err)
	}

	hangups := make([]*hangup.HangupAggregate, len(docs))
	for i, doc := range docs {
		hangups[i] = r.documentToAggregate(&doc)
	}

	return hangups, nil
}

// Count 计数查询
func (r *MongoHangupRepository) Count() (int64, error) {
	ctx := context.Background()

	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		r.logger.Error("Failed to count hangups", "error", err)
		return 0, fmt.Errorf("failed to count hangups: %w", err)
	}

	return count, nil
}

// CountByPlayer 根据玩家计数
func (r *MongoHangupRepository) CountByPlayer(playerID string) (int64, error) {
	ctx := context.Background()

	filter := bson.M{"player_id": playerID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to count hangups by player", "error", err, "player_id", playerID)
		return 0, fmt.Errorf("failed to count hangups by player: %w", err)
	}

	return count, nil
}

// CountByLocation 根据地点计数
func (r *MongoHangupRepository) CountByLocation(locationID string) (int64, error) {
	ctx := context.Background()

	filter := bson.M{"location_id": locationID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to count hangups by location", "error", err, "location_id", locationID)
		return 0, fmt.Errorf("failed to count hangups by location: %w", err)
	}

	return count, nil
}

// 私有方法

// aggregateToDocument 聚合根转文档
func (r *MongoHangupRepository) aggregateToDocument(hangupAggregate *hangup.HangupAggregate) *HangupDocument {
	rewards := make([]RewardItem, 0)
	for _, reward := range hangupAggregate.GetRewards() {
		rewards = append(rewards, RewardItem{
			Type:     reward.Type,
			ItemID:   reward.ItemID,
			Quantity: reward.Quantity,
			Quality:  reward.Quality,
		})
	}

	doc := &HangupDocument{
		HangupID:   hangupAggregate.GetID(),
		PlayerID:   hangupAggregate.GetPlayerID(),
		LocationID: hangupAggregate.GetLocationID(),
		StartTime:  hangupAggregate.GetStartTime(),
		Duration:   int64(hangupAggregate.GetDuration().Seconds()),
		Efficiency: hangupAggregate.GetEfficiency(),
		BaseRate:   hangupAggregate.GetBaseRate(),
		Status:     string(hangupAggregate.GetStatus()),
		Rewards:    rewards,
		CreatedAt:  hangupAggregate.GetCreatedAt(),
		UpdatedAt:  hangupAggregate.GetUpdatedAt(),
	}

	if !hangupAggregate.GetEndTime().IsZero() {
		endTime := hangupAggregate.GetEndTime()
		doc.EndTime = &endTime
	}

	return doc
}

// stringToHangupStatus 字符串转挂机状态
func (r *MongoHangupRepository) stringToHangupStatus(status string) hangup.HangupStatus {
	switch status {
	case "offline":
		return hangup.HangupStatusOffline
	case "online":
		return hangup.HangupStatusOnline
	case "paused":
		return hangup.HangupStatusPaused
	default:
		return hangup.HangupStatusOffline
	}
}

// documentToAggregate 文档转聚合根
func (r *MongoHangupRepository) documentToAggregate(doc *HangupDocument) *hangup.HangupAggregate {
	rewards := make([]hangup.RewardItem, len(doc.Rewards))
	for i, reward := range doc.Rewards {
		rewards[i] = hangup.RewardItem{
			Type:     reward.Type,
			ItemID:   reward.ItemID,
			Quantity: reward.Quantity,
			Quality:  reward.Quality,
		}
	}

	endTime := time.Time{}
	if doc.EndTime != nil {
		endTime = *doc.EndTime
	}

	// 这里需要根据实际的HangupAggregate构造函数来实现
	return hangup.ReconstructHangupAggregate(
		doc.HangupID,
		doc.PlayerID,
		doc.LocationID,
		doc.StartTime,
		endTime,
		time.Duration(doc.Duration)*time.Second,
		doc.Efficiency,
		doc.BaseRate,
		r.stringToHangupStatus(doc.Status),
		rewards,
		doc.CreatedAt,
		doc.UpdatedAt,
	)
}

// CreateIndexes 创建索引
func (r *MongoHangupRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "hangup_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "player_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "location_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "start_time", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "end_time", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "player_id", Value: 1}, {Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "location_id", Value: 1}, {Key: "start_time", Value: -1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		r.logger.Error("Failed to create hangup indexes", "error", err)
		return fmt.Errorf("failed to create hangup indexes: %w", err)
	}

	r.logger.Info("Hangup indexes created successfully")
	return nil
}
