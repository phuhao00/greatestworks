// Package persistence 持久化层实现
package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"greatestworks/internal/domain/replication"
	"greatestworks/internal/infrastructure/cache"
	"greatestworks/internal/infrastructure/logging"
)

// MongoReplicationRepository MongoDB副本仓储实现
type MongoReplicationRepository struct {
	collection *mongo.Collection
	cache      cache.Cache
	logger     logging.Logger
}

// ReplicationInstanceDocument MongoDB文档结构
type ReplicationInstanceDocument struct {
	InstanceID      string                      `bson:"instance_id"`
	TemplateID      string                      `bson:"template_id"`
	InstanceType    int                         `bson:"instance_type"`
	SceneID         string                      `bson:"scene_id"`
	Players         []ReplicationPlayerDocument `bson:"players"`
	MaxPlayers      int                         `bson:"max_players"`
	MinPlayers      int                         `bson:"min_players"` // Updated to include default value
	OwnerPlayerID   string                      `bson:"owner_player_id"`
	Status          int                         `bson:"status"`
	Difficulty      int                         `bson:"difficulty"`
	CreatedAt       time.Time                   `bson:"created_at"`
	StartedAt       time.Time                   `bson:"started_at,omitempty"`
	ExpireAt        time.Time                   `bson:"expire_at"`
	ClosedAt        time.Time                   `bson:"closed_at,omitempty"`
	Lifetime        int64                       `bson:"lifetime"` // 毫秒
	Progress        int                         `bson:"progress"`
	CompletedTasks  []string                    `bson:"completed_tasks"`
	Metadata        map[string]string           `bson:"metadata"`
	ScoreMultiplier float64                     `bson:"score_multiplier"`
	UpdatedAt       time.Time                   `bson:"updated_at"`
}

// ReplicationPlayerDocument 实例中的玩家文档结构
type ReplicationPlayerDocument struct {
	PlayerID   string    `bson:"player_id"`
	PlayerName string    `bson:"player_name"`
	Level      int       `bson:"level"`
	JoinedAt   time.Time `bson:"joined_at"`
	IsReady    bool      `bson:"is_ready"`
	Role       string    `bson:"role"`
}

// NewMongoReplicationRepository 创建MongoDB副本仓储
func NewMongoReplicationRepository(
	db *mongo.Database,
	cache cache.Cache,
	logger logging.Logger,
) *MongoReplicationRepository {
	return &MongoReplicationRepository{
		collection: db.Collection("replication_instances"),
		cache:      cache,
		logger:     logger,
	}
}

// Save 保存实例
func (r *MongoReplicationRepository) Save(ctx context.Context, instance *replication.Instance) error {
	doc := r.toDocument(instance)

	filter := bson.M{"instance_id": instance.ID()}
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("保存实例失败: %w", err)
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("instance:%s", instance.ID())
	_ = r.cache.Delete(ctx, cacheKey)

	return nil
}

// FindByID 根据ID查找实例
func (r *MongoReplicationRepository) FindByID(ctx context.Context, instanceID string) (*replication.Instance, error) {
	// 先查缓存
	cacheKey := fmt.Sprintf("instance:%s", instanceID)
	var cachedData string
	if err := r.cache.Get(ctx, cacheKey, &cachedData); err == nil {
		var doc ReplicationInstanceDocument
		if err := json.Unmarshal([]byte(cachedData), &doc); err == nil {
			return r.toDomain(&doc), nil
		}
	}

	// 查数据库
	var doc ReplicationInstanceDocument
	filter := bson.M{"instance_id": instanceID}
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查找实例失败: %w", err)
	}

	// 写入缓存
	if data, err := json.Marshal(doc); err == nil {
		_ = r.cache.Set(ctx, cacheKey, string(data), 5*time.Minute)
	}

	return r.toDomain(&doc), nil
}

// FindByTemplateID 根据模板ID查找实例列表
func (r *MongoReplicationRepository) FindByTemplateID(ctx context.Context, templateID string) ([]*replication.Instance, error) {
	filter := bson.M{"template_id": templateID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查找实例失败: %w", err)
	}
	defer cursor.Close(ctx)

	var instances []*replication.Instance
	for cursor.Next(ctx) {
		var doc ReplicationInstanceDocument
		if err := cursor.Decode(&doc); err != nil {
			r.logger.Error("解码实例文档失败", err, logging.Fields{})
			continue
		}
		instances = append(instances, r.toDomain(&doc))
	}

	return instances, nil
}

// FindActiveInstances 查找所有活跃实例
func (r *MongoReplicationRepository) FindActiveInstances(ctx context.Context) ([]*replication.Instance, error) {
	filter := bson.M{
		"status": bson.M{"$in": []int{
			int(replication.InstanceStatusPending),
			int(replication.InstanceStatusCreating),
			int(replication.InstanceStatusActive),
			int(replication.InstanceStatusFull),
		}},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查找活跃实例失败: %w", err)
	}
	defer cursor.Close(ctx)

	var instances []*replication.Instance
	for cursor.Next(ctx) {
		var doc ReplicationInstanceDocument
		if err := cursor.Decode(&doc); err != nil {
			r.logger.Error("解码实例文档失败", err, logging.Fields{})
			continue
		}
		instances = append(instances, r.toDomain(&doc))
	}

	return instances, nil
}

// FindByPlayerID 根据玩家ID查找实例
func (r *MongoReplicationRepository) FindByPlayerID(ctx context.Context, playerID string) (*replication.Instance, error) {
	filter := bson.M{
		"players.player_id": playerID,
		"status": bson.M{"$in": []int{
			int(replication.InstanceStatusActive),
			int(replication.InstanceStatusFull),
		}},
	}

	var doc ReplicationInstanceDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查找玩家实例失败: %w", err)
	}

	return r.toDomain(&doc), nil
}

// Delete 删除实例
func (r *MongoReplicationRepository) Delete(ctx context.Context, instanceID string) error {
	filter := bson.M{"instance_id": instanceID}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除实例失败: %w", err)
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("instance:%s", instanceID)
	_ = r.cache.Delete(ctx, cacheKey)

	return nil
}

// UpdateStatus 更新实例状态
func (r *MongoReplicationRepository) UpdateStatus(ctx context.Context, instanceID string, status replication.InstanceStatus) error {
	filter := bson.M{"instance_id": instanceID}
	update := bson.M{
		"$set": bson.M{
			"status":     int(status),
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新实例状态失败: %w", err)
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("instance:%s", instanceID)
	_ = r.cache.Delete(ctx, cacheKey)

	return nil
}

// FindExpiredInstances 查找过期实例
func (r *MongoReplicationRepository) FindExpiredInstances(ctx context.Context) ([]*replication.Instance, error) {
	filter := bson.M{
		"expire_at": bson.M{"$lt": time.Now()},
		"status":    bson.M{"$ne": int(replication.InstanceStatusClosed)},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查找过期实例失败: %w", err)
	}
	defer cursor.Close(ctx)

	var instances []*replication.Instance
	for cursor.Next(ctx) {
		var doc ReplicationInstanceDocument
		if err := cursor.Decode(&doc); err != nil {
			r.logger.Error("解码实例文档失败", err, logging.Fields{})
			continue
		}
		instances = append(instances, r.toDomain(&doc))
	}

	return instances, nil
}

// toDocument 转换为文档
func (r *MongoReplicationRepository) toDocument(instance *replication.Instance) *ReplicationInstanceDocument {
	players := instance.GetPlayers()
	playerDocs := make([]ReplicationPlayerDocument, 0, len(players))
	for _, p := range players {
		playerDocs = append(playerDocs, ReplicationPlayerDocument{
			PlayerID:   p.PlayerID,
			PlayerName: p.PlayerName,
			Level:      p.Level,
			JoinedAt:   p.JoinedAt,
			IsReady:    p.IsReady,
			Role:       p.Role,
		})
	}

	return &ReplicationInstanceDocument{
		InstanceID:   instance.ID(),
		TemplateID:   instance.TemplateID(),
		InstanceType: int(instance.Type()),
		SceneID:      instance.SceneID(),
		Players:      playerDocs,
		MaxPlayers:   instance.MaxPlayers(),
		Status:       int(instance.Status()),
		Progress:     instance.Progress(),
		CreatedAt:    instance.CreatedAt(),
		Difficulty:   instance.Difficulty(),
		UpdatedAt:    time.Now(),
	}
}

// toDomain 转换为领域对象
func (r *MongoReplicationRepository) toDomain(doc *ReplicationInstanceDocument) *replication.Instance {
	// 通过快照重建领域对象
	players := make([]replication.PlayerInfo, 0, len(doc.Players))
	for _, p := range doc.Players {
		players = append(players, replication.PlayerInfo{
			PlayerID:   p.PlayerID,
			PlayerName: p.PlayerName,
			Level:      p.Level,
			JoinedAt:   p.JoinedAt,
			IsReady:    p.IsReady,
			Role:       p.Role,
		})
	}
	snap := replication.InstanceSnapshot{
		InstanceID:    doc.InstanceID,
		TemplateID:    doc.TemplateID,
		SceneID:       doc.SceneID,
		OwnerPlayerID: doc.OwnerPlayerID,
		InstanceType:  replication.InstanceType(doc.InstanceType),
		Status:        replication.InstanceStatus(doc.Status),
		MaxPlayers:    doc.MaxPlayers,
		MinPlayers:    doc.MinPlayers,
		Difficulty:    doc.Difficulty,
		CreatedAt:     doc.CreatedAt,
		StartedAt:     doc.StartedAt,
		ExpireAt:      doc.ExpireAt,
		ClosedAt:      doc.ClosedAt,
		Lifetime:      time.Duration(doc.Lifetime) * time.Millisecond,
		Progress:      doc.Progress,
		Completed:     append([]string(nil), doc.CompletedTasks...),
		Metadata:      doc.Metadata,
		Players:       players,
	}
	return replication.NewInstanceFromSnapshot(snap)
}
