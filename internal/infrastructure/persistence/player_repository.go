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
	"greatestworks/internal/infrastructure/logging"
)

// MongoPlayerRepository MongoDB玩家仓储实现
type MongoPlayerRepository struct {
	collection *mongo.Collection
	cache      cache.Cache
	logger     logging.Logger
}

// NewMongoPlayerRepository 创建MongoDB玩家仓储
func NewMongoPlayerRepository(db *mongo.Database, cache cache.Cache, logger logging.Logger) *MongoPlayerRepository {
	return &MongoPlayerRepository{
		collection: db.Collection("players"),
		cache:      cache,
		logger:     logger,
	}
}

// PlayerDocument 玩家文档结构
type PlayerDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Level     int                `bson:"level" json:"level"`
	Exp       int64              `bson:"exp" json:"exp"`
	Status    int                `bson:"status" json:"status"`
	Position  PlayerPosition     `bson:"position" json:"position"`
	LastMapID int32              `bson:"last_map_id" json:"last_map_id"`
	Stats     PlayerStats        `bson:"stats" json:"stats"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Version   int64              `bson:"version" json:"version"`
}

// PlayerPosition 玩家位置值对象
type PlayerPosition struct {
	X float64 `bson:"x" json:"x"`
	Y float64 `bson:"y" json:"y"`
	Z float64 `bson:"z" json:"z"`
}

// PlayerStats 玩家属性值对象
type PlayerStats struct {
	HP      int `bson:"hp" json:"hp"`
	MaxHP   int `bson:"max_hp" json:"max_hp"`
	MP      int `bson:"mp" json:"mp"`
	MaxMP   int `bson:"max_mp" json:"max_mp"`
	Attack  int `bson:"attack" json:"attack"`
	Defense int `bson:"defense" json:"defense"`
	Speed   int `bson:"speed" json:"speed"`
}

// Save 保存玩家
func (r *MongoPlayerRepository) Save(ctx context.Context, p *player.Player) error {
	doc := r.toDocument(p)

	// 设置时间戳
	now := time.Now()
	if doc.ID.IsZero() {
		doc.CreatedAt = now
	}
	doc.UpdatedAt = now

	// 使用Upsert操作
	filter := bson.M{"name": p.Name()}
	opts := options.Replace().SetUpsert(true)

	_, err := r.collection.ReplaceOne(ctx, filter, doc, opts)
	if err != nil {
		r.logger.Error("保存玩家失败", err, logging.Fields{
			"name": p.Name(),
		})
		return fmt.Errorf("保存玩家失败: %w", err)
	}

	// 更新缓存
	cacheKey := fmt.Sprintf("player:id:%s", p.ID().String())
	if err := r.cache.Set(ctx, cacheKey, p, time.Hour); err != nil {
		r.logger.Warn("更新玩家缓存失败", map[string]interface{}{
			"name":  p.Name(),
			"error": err.Error(),
		})
	}

	r.logger.Info("玩家保存成功", map[string]interface{}{
		"name":  p.Name(),
		"level": p.Level(),
	})

	return nil
}

// FindByID 根据ID查找玩家
func (r *MongoPlayerRepository) FindByID(ctx context.Context, id string) (*player.Player, error) {
	// 先从缓存获取
	cacheKey := fmt.Sprintf("player:id:%s", id)
	var cachedPlayer *player.Player
	if err := r.cache.Get(ctx, cacheKey, &cachedPlayer); err == nil && cachedPlayer != nil {
		return cachedPlayer, nil
	}

	// 从数据库获取
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的ID格式: %w", err)
	}

	filter := bson.M{"_id": objectID}
	var doc PlayerDocument
	err = r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, player.ErrPlayerNotFound
		}
		r.logger.Error("查找玩家失败", err, logging.Fields{
			"id": id,
		})
		return nil, fmt.Errorf("查找玩家失败: %w", err)
	}

	// 转换为领域对象
	playerAggregate, err := r.toAggregate(&doc)
	if err != nil {
		return nil, fmt.Errorf("转换玩家对象失败: %w", err)
	}

	// 更新缓存
	if err := r.cache.Set(ctx, cacheKey, playerAggregate, time.Hour); err != nil {
		r.logger.Warn("更新玩家缓存失败", map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		})
	}

	return playerAggregate, nil
}

// FindByName 根据名称查找玩家
func (r *MongoPlayerRepository) FindByName(ctx context.Context, name string) (*player.Player, error) {
	// 先从缓存获取
	cacheKey := fmt.Sprintf("player:name:%s", name)
	var cachedPlayer *player.Player
	if err := r.cache.Get(ctx, cacheKey, &cachedPlayer); err == nil && cachedPlayer != nil {
		return cachedPlayer, nil
	}

	// 从数据库获取
	filter := bson.M{"name": name}
	var doc PlayerDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, player.ErrPlayerNotFound
		}
		r.logger.Error("根据名称查找玩家失败", err, logging.Fields{
			"name": name,
		})
		return nil, fmt.Errorf("根据名称查找玩家失败: %w", err)
	}

	// 转换为领域对象
	playerAggregate, err := r.toAggregate(&doc)
	if err != nil {
		return nil, fmt.Errorf("转换玩家对象失败: %w", err)
	}

	// 更新缓存
	if err := r.cache.Set(ctx, cacheKey, playerAggregate, time.Hour); err != nil {
		r.logger.Warn("更新玩家缓存失败", map[string]interface{}{
			"name":  name,
			"error": err.Error(),
		})
	}

	return playerAggregate, nil
}

// Delete 删除玩家
func (r *MongoPlayerRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	filter := bson.M{"_id": objectID}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("删除玩家失败", err, logging.Fields{
			"id": id,
		})
		return fmt.Errorf("删除玩家失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return player.ErrPlayerNotFound
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("player:id:%s", id)
	if err := r.cache.Delete(ctx, cacheKey); err != nil {
		r.logger.Warn("清除玩家缓存失败", map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		})
	}

	r.logger.Info("玩家删除成功", map[string]interface{}{
		"id": id,
	})

	return nil
}

// List 获取玩家列表
func (r *MongoPlayerRepository) List(ctx context.Context, limit, offset int) ([]*player.Player, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		r.logger.Error("获取玩家列表失败", err, logging.Fields{})
		return nil, fmt.Errorf("获取玩家列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []PlayerDocument
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("解析玩家列表失败: %w", err)
	}

	players := make([]*player.Player, 0, len(docs))
	for _, doc := range docs {
		playerAggregate, err := r.toAggregate(&doc)
		if err != nil {
			r.logger.Error("转换玩家对象失败", err, logging.Fields{
				"id": doc.ID.Hex(),
			})
			continue
		}
		players = append(players, playerAggregate)
	}

	r.logger.Info("获取玩家列表成功", map[string]interface{}{
		"count": len(players),
	})

	return players, nil
}

// Count 获取玩家总数
func (r *MongoPlayerRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		r.logger.Error("获取玩家总数失败", err, logging.Fields{})
		return 0, fmt.Errorf("获取玩家总数失败: %w", err)
	}

	return count, nil
}

// toDocument 转换为文档
func (r *MongoPlayerRepository) toDocument(p *player.Player) *PlayerDocument {
	position := p.GetPosition()
	stats := p.Stats()

	doc := &PlayerDocument{
		Name:      p.Name(),
		Level:     p.Level(),
		Exp:       p.Exp(),
		Status:    int(p.Status()),
		Position:  PlayerPosition{X: position.X, Y: position.Y, Z: position.Z},
		LastMapID: p.LastMapID(),
		Stats:     PlayerStats{HP: stats.HP, MaxHP: stats.MaxHP, MP: stats.MP, MaxMP: stats.MaxMP, Attack: stats.Attack, Defense: stats.Defense, Speed: stats.Speed},
		CreatedAt: p.CreatedAt(),
		UpdatedAt: p.UpdatedAt(),
		Version:   p.Version(),
	}

	// 如果有ID，转换为ObjectID
	if p.ID().String() != "" {
		if objectID, err := primitive.ObjectIDFromHex(p.ID().String()); err == nil {
			doc.ID = objectID
		}
	}

	return doc
}

// toAggregate 转换为聚合根
func (r *MongoPlayerRepository) toAggregate(doc *PlayerDocument) (*player.Player, error) {
	// 使用ReconstructPlayer方法从持久化数据重建玩家聚合根
	playerID := player.PlayerIDFromString(doc.ID.Hex())
	status := player.PlayerStatus(doc.Status)
	position := player.Position{X: doc.Position.X, Y: doc.Position.Y, Z: doc.Position.Z}
	stats := player.PlayerStats{HP: doc.Stats.HP, MaxHP: doc.Stats.MaxHP, MP: doc.Stats.MP, MaxMP: doc.Stats.MaxMP, Attack: doc.Stats.Attack, Defense: doc.Stats.Defense, Speed: doc.Stats.Speed}

	p := player.ReconstructPlayer(
		playerID,
		doc.Name,
		doc.Level,
		doc.Exp,
		status,
		position,
		doc.LastMapID,
		stats,
		doc.CreatedAt,
		doc.UpdatedAt,
		doc.Version,
	)

	return p, nil
}
