package persistence

import (
	"context"
	"fmt"
	"time"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	"greatestworks/internal/domain/player"
	"greatestworks/internal/infrastructure/cache"
	"greatestworks/aop/logger"
)

// MongoPlayerRepository MongoDB玩家仓储实现
type MongoPlayerRepository struct {
	db         *mongo.Database
	cache      cache.Cache
	logger     logger.Logger
	collection *mongo.Collection
}

// NewMongoPlayerRepository 创建MongoDB玩家仓储
func NewMongoPlayerRepository(db *mongo.Database, cache cache.Cache, logger logger.Logger) player.Repository {
	return &MongoPlayerRepository{
		db:         db,
		cache:      cache,
		logger:     logger,
		collection: db.Collection("players"),
	}
}

// PlayerDocument MongoDB玩家文档结构
type PlayerDocument struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	PlayerID    string             `bson:"player_id"`
	Username    string             `bson:"username"`
	Nickname    string             `bson:"nickname"`
	Level       int                `bson:"level"`
	Experience  int64              `bson:"experience"`
	Gold        int64              `bson:"gold"`
	Diamond     int64              `bson:"diamond"`
	VIPLevel    int                `bson:"vip_level"`
	Status      string             `bson:"status"`
	LastLoginAt time.Time          `bson:"last_login_at"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
	Attributes  map[string]int64   `bson:"attributes"`
	Settings    map[string]string  `bson:"settings"`
}

// Save 保存玩家
func (r *MongoPlayerRepository) Save(ctx context.Context, playerAggregate *player.PlayerAggregate) error {
	doc := r.aggregateToDocument(playerAggregate)
	doc.UpdatedAt = time.Now()
	
	if doc.ID.IsZero() {
		doc.CreatedAt = time.Now()
		result, err := r.collection.InsertOne(ctx, doc)
		if err != nil {
			r.logger.Error("Failed to insert player", "error", err, "player_id", playerAggregate.GetID())
			return fmt.Errorf("failed to insert player: %w", err)
		}
		
		// 更新聚合根ID
		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
			doc.ID = oid
		}
	} else {
		filter := bson.M{"player_id": playerAggregate.GetID()}
		update := bson.M{"$set": doc}
		
		_, err := r.collection.UpdateOne(ctx, filter, update)
		if err != nil {
			r.logger.Error("Failed to update player", "error", err, "player_id", playerAggregate.GetID())
			return fmt.Errorf("failed to update player: %w", err)
		}
	}
	
	// 更新缓存
	cacheKey := fmt.Sprintf("player:%s", playerAggregate.GetID())
	if err := r.cache.Set(ctx, cacheKey, playerAggregate, time.Hour); err != nil {
		r.logger.Warn("Failed to cache player", "error", err, "player_id", playerAggregate.GetID())
	}
	
	return nil
}

// FindByID 根据ID查找玩家
func (r *MongoPlayerRepository) FindByID(ctx context.Context, playerID string) (*player.PlayerAggregate, error) {
	// 先从缓存获取
	cacheKey := fmt.Sprintf("player:%s", playerID)
	var cachedPlayer *player.PlayerAggregate
	if err := r.cache.Get(ctx, cacheKey, &cachedPlayer); err == nil && cachedPlayer != nil {
		return cachedPlayer, nil
	}
	
	// 从数据库获取
	filter := bson.M{"player_id": playerID}
	var doc PlayerDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, player.ErrPlayerNotFound
		}
		r.logger.Error("Failed to find player", "error", err, "player_id", playerID)
		return nil, fmt.Errorf("failed to find player: %w", err)
	}
	
	playerAggregate := r.documentToAggregate(&doc)
	
	// 更新缓存
	if err := r.cache.Set(ctx, cacheKey, playerAggregate, time.Hour); err != nil {
		r.logger.Warn("Failed to cache player", "error", err, "player_id", playerID)
	}
	
	return playerAggregate, nil
}

// FindByUsername 根据用户名查找玩家
func (r *MongoPlayerRepository) FindByUsername(ctx context.Context, username string) (*player.PlayerAggregate, error) {
	// 先从缓存获取
	cacheKey := fmt.Sprintf("player:username:%s", username)
	var cachedPlayer *player.PlayerAggregate
	if err := r.cache.Get(ctx, cacheKey, &cachedPlayer); err == nil && cachedPlayer != nil {
		return cachedPlayer, nil
	}
	
	// 从数据库获取
	filter := bson.M{"username": username}
	var doc PlayerDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, player.ErrPlayerNotFound
		}
		r.logger.Error("Failed to find player by username", "error", err, "username", username)
		return nil, fmt.Errorf("failed to find player by username: %w", err)
	}
	
	playerAggregate := r.documentToAggregate(&doc)
	
	// 更新缓存
	if err := r.cache.Set(ctx, cacheKey, playerAggregate, time.Hour); err != nil {
		r.logger.Warn("Failed to cache player by username", "error", err, "username", username)
	}
	
	return playerAggregate, nil
}

// Update 更新玩家
func (r *MongoPlayerRepository) Update(ctx context.Context, playerAggregate *player.PlayerAggregate) error {
	return r.Save(ctx, playerAggregate)
}

// Delete 删除玩家
func (r *MongoPlayerRepository) Delete(ctx context.Context, playerID string) error {
	filter := bson.M{"player_id": playerID}
	
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete player", "error", err, "player_id", playerID)
		return fmt.Errorf("failed to delete player: %w", err)
	}
	
	if result.DeletedCount == 0 {
		return player.ErrPlayerNotFound
	}
	
	// 清除缓存
	cacheKey := fmt.Sprintf("player:%s", playerID)
	if err := r.cache.Delete(ctx, cacheKey); err != nil {
		r.logger.Warn("Failed to delete player cache", "error", err, "player_id", playerID)
	}
	
	return nil
}

// List 列表查询玩家
func (r *MongoPlayerRepository) List(ctx context.Context, query *player.PlayerQuery) ([]*player.PlayerAggregate, error) {
	filter := r.buildFilter(query)
	opts := r.buildOptions(query)
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to list players", "error", err)
		return nil, fmt.Errorf("failed to list players: %w", err)
	}
	defer cursor.Close(ctx)
	
	var docs []PlayerDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode players", "error", err)
		return nil, fmt.Errorf("failed to decode players: %w", err)
	}
	
	players := make([]*player.PlayerAggregate, len(docs))
	for i, doc := range docs {
		players[i] = r.documentToAggregate(&doc)
	}
	
	return players, nil
}

// Count 计数查询
func (r *MongoPlayerRepository) Count(ctx context.Context, query *player.PlayerQuery) (int64, error) {
	filter := r.buildFilter(query)
	
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to count players", "error", err)
		return 0, fmt.Errorf("failed to count players: %w", err)
	}
	
	return count, nil
}

// FindByLevel 根据等级范围查找玩家
func (r *MongoPlayerRepository) FindByLevel(ctx context.Context, minLevel, maxLevel int) ([]*player.PlayerAggregate, error) {
	filter := bson.M{
		"level": bson.M{
			"$gte": minLevel,
			"$lte": maxLevel,
		},
	}
	
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find players by level", "error", err)
		return nil, fmt.Errorf("failed to find players by level: %w", err)
	}
	defer cursor.Close(ctx)
	
	var docs []PlayerDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode players by level", "error", err)
		return nil, fmt.Errorf("failed to decode players by level: %w", err)
	}
	
	players := make([]*player.PlayerAggregate, len(docs))
	for i, doc := range docs {
		players[i] = r.documentToAggregate(&doc)
	}
	
	return players, nil
}

// FindOnlinePlayers 查找在线玩家
func (r *MongoPlayerRepository) FindOnlinePlayers(ctx context.Context) ([]*player.PlayerAggregate, error) {
	filter := bson.M{
		"status": "online",
		"last_login_at": bson.M{
			"$gte": time.Now().Add(-time.Hour), // 1小时内登录的视为在线
		},
	}
	
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find online players", "error", err)
		return nil, fmt.Errorf("failed to find online players: %w", err)
	}
	defer cursor.Close(ctx)
	
	var docs []PlayerDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode online players", "error", err)
		return nil, fmt.Errorf("failed to decode online players: %w", err)
	}
	
	players := make([]*player.PlayerAggregate, len(docs))
	for i, doc := range docs {
		players[i] = r.documentToAggregate(&doc)
	}
	
	return players, nil
}

// UpdateLastLogin 更新最后登录时间
func (r *MongoPlayerRepository) UpdateLastLogin(ctx context.Context, playerID string) error {
	filter := bson.M{"player_id": playerID}
	update := bson.M{
		"$set": bson.M{
			"last_login_at": time.Now(),
			"updated_at":    time.Now(),
		},
	}
	
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("Failed to update last login", "error", err, "player_id", playerID)
		return fmt.Errorf("failed to update last login: %w", err)
	}
	
	// 清除缓存
	cacheKey := fmt.Sprintf("player:%s", playerID)
	if err := r.cache.Delete(ctx, cacheKey); err != nil {
		r.logger.Warn("Failed to delete player cache after login update", "error", err, "player_id", playerID)
	}
	
	return nil
}

// 私有方法

// aggregateToDocument 聚合根转文档
func (r *MongoPlayerRepository) aggregateToDocument(playerAggregate *player.PlayerAggregate) *PlayerDocument {
	return &PlayerDocument{
		PlayerID:    playerAggregate.GetID(),
		Username:    playerAggregate.GetUsername(),
		Nickname:    playerAggregate.GetNickname(),
		Level:       playerAggregate.GetLevel(),
		Experience:  playerAggregate.GetExperience(),
		Gold:        playerAggregate.GetGold(),
		Diamond:     playerAggregate.GetDiamond(),
		VIPLevel:    playerAggregate.GetVIPLevel(),
		Status:      string(playerAggregate.GetStatus()),
		LastLoginAt: playerAggregate.GetLastLoginAt(),
		CreatedAt:   playerAggregate.GetCreatedAt(),
		UpdatedAt:   playerAggregate.GetUpdatedAt(),
		Attributes:  playerAggregate.GetAttributes(),
		Settings:    playerAggregate.GetSettings(),
	}
}

// documentToAggregate 文档转聚合根
func (r *MongoPlayerRepository) documentToAggregate(doc *PlayerDocument) *player.PlayerAggregate {
	// 这里需要根据实际的PlayerAggregate构造函数来实现
	// 假设有一个重建聚合根的方法
	return player.ReconstructPlayerAggregate(
		doc.PlayerID,
		doc.Username,
		doc.Nickname,
		doc.Level,
		doc.Experience,
		doc.Gold,
		doc.Diamond,
		doc.VIPLevel,
		player.PlayerStatus(doc.Status),
		doc.LastLoginAt,
		doc.CreatedAt,
		doc.UpdatedAt,
		doc.Attributes,
		doc.Settings,
	)
}

// buildFilter 构建查询过滤器
func (r *MongoPlayerRepository) buildFilter(query *player.PlayerQuery) bson.M {
	filter := bson.M{}
	
	if query.Username != "" {
		filter["username"] = bson.M{"$regex": query.Username, "$options": "i"}
	}
	
	if query.Nickname != "" {
		filter["nickname"] = bson.M{"$regex": query.Nickname, "$options": "i"}
	}
	
	if query.MinLevel > 0 {
		if filter["level"] == nil {
			filter["level"] = bson.M{}
		}
		filter["level"].(bson.M)["$gte"] = query.MinLevel
	}
	
	if query.MaxLevel > 0 {
		if filter["level"] == nil {
			filter["level"] = bson.M{}
		}
		filter["level"].(bson.M)["$lte"] = query.MaxLevel
	}
	
	if query.Status != "" {
		filter["status"] = query.Status
	}
	
	if query.VIPLevel > 0 {
		filter["vip_level"] = bson.M{"$gte": query.VIPLevel}
	}
	
	if !query.CreatedAfter.IsZero() {
		if filter["created_at"] == nil {
			filter["created_at"] = bson.M{}
		}
		filter["created_at"].(bson.M)["$gte"] = query.CreatedAfter
	}
	
	if !query.CreatedBefore.IsZero() {
		if filter["created_at"] == nil {
			filter["created_at"] = bson.M{}
		}
		filter["created_at"].(bson.M)["$lte"] = query.CreatedBefore
	}
	
	return filter
}

// buildOptions 构建查询选项
func (r *MongoPlayerRepository) buildOptions(query *player.PlayerQuery) *options.FindOptions {
	opts := options.Find()
	
	if query.Limit > 0 {
		opts.SetLimit(int64(query.Limit))
	}
	
	if query.Offset > 0 {
		opts.SetSkip(int64(query.Offset))
	}
	
	if query.OrderBy != "" {
		sortOrder := 1
		if query.OrderDesc {
			sortOrder = -1
		}
		opts.SetSort(bson.D{{Key: query.OrderBy, Value: sortOrder}})
	}
	
	return opts
}

// CreateIndexes 创建索引
func (r *MongoPlayerRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "player_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "level", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "last_login_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "vip_level", Value: -1}},
		},
	}
	
	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		r.logger.Error("Failed to create player indexes", "error", err)
		return fmt.Errorf("failed to create player indexes: %w", err)
	}
	
	r.logger.Info("Player indexes created successfully")
	return nil
}