package persistence

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"greatestworks/internal/domain/scene"
	"greatestworks/internal/infrastructure/cache"
	"greatestworks/internal/infrastructure/logging"
)

// MongoSceneRepository MongoDB场景仓储实现
type MongoSceneRepository struct {
	collection *mongo.Collection
	cache      cache.Cache
	logger     logging.Logger
}

// NewMongoSceneRepository 创建MongoDB场景仓储
func NewMongoSceneRepository(db *mongo.Database, cache cache.Cache, logger logging.Logger) *MongoSceneRepository {
	return &MongoSceneRepository{
		collection: db.Collection("scenes"),
		cache:      cache,
		logger:     logger,
	}
}

// SceneDocument 场景文档结构
type SceneDocument struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SceneID        string             `bson:"scene_id" json:"scene_id"`
	Name           string             `bson:"name" json:"name"`
	SceneType      int                `bson:"scene_type" json:"scene_type"`
	Status         int                `bson:"status" json:"status"`
	Width          float64            `bson:"width" json:"width"`
	Height         float64            `bson:"height" json:"height"`
	MaxPlayers     int                `bson:"max_players" json:"max_players"`
	CurrentPlayers int                `bson:"current_players" json:"current_players"`
	Players        []string           `bson:"players" json:"players"` // 玩家ID列表
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
	Version        int64              `bson:"version" json:"version"`
}

// Save 保存场景
func (r *MongoSceneRepository) Save(ctx context.Context, s *scene.Scene) error {
	doc := r.toDocument(s)

	now := time.Now()
	if doc.ID.IsZero() {
		doc.CreatedAt = now
	}
	doc.UpdatedAt = now

	filter := bson.M{"scene_id": s.ID()}
	opts := options.Replace().SetUpsert(true)

	_, err := r.collection.ReplaceOne(ctx, filter, doc, opts)
	if err != nil {
		return fmt.Errorf("保存场景到MongoDB失败: %w", err)
	}

	// 清理缓存
	if r.cache != nil {
		cacheKey := fmt.Sprintf("scene:%s", s.ID())
		_ = r.cache.Delete(ctx, cacheKey)
	}

	return nil
}

// FindByID 根据ID查找场景
func (r *MongoSceneRepository) FindByID(ctx context.Context, sceneID string) (*scene.Scene, error) {
	// 先尝试从缓存获取
	if r.cache != nil {
		cacheKey := fmt.Sprintf("scene:%s", sceneID)
		var doc SceneDocument
		if err := r.cache.Get(ctx, cacheKey, &doc); err == nil {
			return r.toDomain(&doc), nil
		}
	}

	// 从MongoDB获取
	filter := bson.M{"scene_id": sceneID}
	var doc SceneDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("从MongoDB查找场景失败: %w", err)
	}

	// 写入缓存
	if r.cache != nil {
		cacheKey := fmt.Sprintf("scene:%s", sceneID)
		_ = r.cache.Set(ctx, cacheKey, &doc, 5*time.Minute)
	}

	return r.toDomain(&doc), nil
}

// Delete 删除场景
func (r *MongoSceneRepository) Delete(ctx context.Context, sceneID string) error {
	filter := bson.M{"scene_id": sceneID}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("从MongoDB删除场景失败: %w", err)
	}

	// 清理缓存
	if r.cache != nil {
		cacheKey := fmt.Sprintf("scene:%s", sceneID)
		_ = r.cache.Delete(ctx, cacheKey)
	}

	return nil
}

// Exists 检查场景是否存在
func (r *MongoSceneRepository) Exists(ctx context.Context, sceneID string) (bool, error) {
	filter := bson.M{"scene_id": sceneID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("检查场景是否存在失败: %w", err)
	}
	return count > 0, nil
}

// SaveBatch 批量保存场景
func (r *MongoSceneRepository) SaveBatch(ctx context.Context, scenes []*scene.Scene) error {
	if len(scenes) == 0 {
		return nil
	}

	models := make([]mongo.WriteModel, 0, len(scenes))
	for _, s := range scenes {
		doc := r.toDocument(s)
		now := time.Now()
		if doc.ID.IsZero() {
			doc.CreatedAt = now
		}
		doc.UpdatedAt = now

		filter := bson.M{"scene_id": s.ID()}
		update := bson.M{"$set": doc}
		model := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		models = append(models, model)
	}

	_, err := r.collection.BulkWrite(ctx, models)
	if err != nil {
		return fmt.Errorf("批量保存场景失败: %w", err)
	}

	return nil
}

// FindByIDs 根据ID列表查找场景
func (r *MongoSceneRepository) FindByIDs(ctx context.Context, sceneIDs []string) ([]*scene.Scene, error) {
	filter := bson.M{"scene_id": bson.M{"$in": sceneIDs}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查找场景列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []SceneDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("解析场景列表失败: %w", err)
	}

	scenes := make([]*scene.Scene, 0, len(docs))
	for _, doc := range docs {
		scenes = append(scenes, r.toDomain(&doc))
	}

	return scenes, nil
}

// FindAll 查找所有场景
func (r *MongoSceneRepository) FindAll(ctx context.Context) ([]*scene.Scene, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("查找所有场景失败: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []SceneDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("解析场景列表失败: %w", err)
	}

	scenes := make([]*scene.Scene, 0, len(docs))
	for _, doc := range docs {
		scenes = append(scenes, r.toDomain(&doc))
	}

	return scenes, nil
}

// FindByType 根据类型查找场景
func (r *MongoSceneRepository) FindByType(ctx context.Context, sceneType scene.SceneType) ([]*scene.Scene, error) {
	filter := bson.M{"scene_type": int(sceneType)}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查找场景失败: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []SceneDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("解析场景列表失败: %w", err)
	}

	scenes := make([]*scene.Scene, 0, len(docs))
	for _, doc := range docs {
		scenes = append(scenes, r.toDomain(&doc))
	}

	return scenes, nil
}

// FindByStatus 根据状态查找场景
func (r *MongoSceneRepository) FindByStatus(ctx context.Context, status scene.SceneStatus) ([]*scene.Scene, error) {
	filter := bson.M{"status": int(status)}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查找场景失败: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []SceneDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("解析场景列表失败: %w", err)
	}

	scenes := make([]*scene.Scene, 0, len(docs))
	for _, doc := range docs {
		scenes = append(scenes, r.toDomain(&doc))
	}

	return scenes, nil
}

// FindAvailableScenes 查找可用场景
func (r *MongoSceneRepository) FindAvailableScenes(ctx context.Context) ([]*scene.Scene, error) {
	filter := bson.M{
		"status": int(scene.SceneStatusActive),
		"$expr": bson.M{
			"$lt": []interface{}{"$current_players", "$max_players"},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查找可用场景失败: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []SceneDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("解析场景列表失败: %w", err)
	}

	scenes := make([]*scene.Scene, 0, len(docs))
	for _, doc := range docs {
		scenes = append(scenes, r.toDomain(&doc))
	}

	return scenes, nil
}

// FindScenesWithSpace 查找有空位的场景
func (r *MongoSceneRepository) FindScenesWithSpace(ctx context.Context, minSpace int) ([]*scene.Scene, error) {
	filter := bson.M{
		"$expr": bson.M{
			"$gte": []interface{}{
				bson.M{"$subtract": []interface{}{"$max_players", "$current_players"}},
				minSpace,
			},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查找有空位的场景失败: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []SceneDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("解析场景列表失败: %w", err)
	}

	scenes := make([]*scene.Scene, 0, len(docs))
	for _, doc := range docs {
		scenes = append(scenes, r.toDomain(&doc))
	}

	return scenes, nil
}

// SaveEntity 保存实体（简化实现，实际应该独立存储）
func (r *MongoSceneRepository) SaveEntity(ctx context.Context, sceneID string, entity scene.Entity) error {
	// TODO: 实现实体持久化逻辑
	return nil
}

// RemoveEntity 移除实体
func (r *MongoSceneRepository) RemoveEntity(ctx context.Context, sceneID string, entityID string) error {
	// TODO: 实现实体移除逻辑
	return nil
}

// FindEntitiesByType 根据类型查找实体
func (r *MongoSceneRepository) FindEntitiesByType(ctx context.Context, sceneID string, entityType scene.EntityType) ([]scene.Entity, error) {
	// TODO: 实现实体查询逻辑
	return nil, nil
}

// FindEntitiesInRadius 查找半径内的实体
func (r *MongoSceneRepository) FindEntitiesInRadius(ctx context.Context, sceneID string, center *scene.Position, radius float64) ([]scene.Entity, error) {
	// TODO: 实现空间查询逻辑
	return nil, nil
}

// AddPlayerToScene 添加玩家到场景
func (r *MongoSceneRepository) AddPlayerToScene(ctx context.Context, sceneID string, playerID string) error {
	filter := bson.M{"scene_id": sceneID}
	update := bson.M{
		"$addToSet": bson.M{"players": playerID},
		"$inc":      bson.M{"current_players": 1},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// RemovePlayerFromScene 从场景移除玩家
func (r *MongoSceneRepository) RemovePlayerFromScene(ctx context.Context, sceneID string, playerID string) error {
	filter := bson.M{"scene_id": sceneID}
	update := bson.M{
		"$pull": bson.M{"players": playerID},
		"$inc":  bson.M{"current_players": -1},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// FindPlayerScene 查找玩家所在场景
func (r *MongoSceneRepository) FindPlayerScene(ctx context.Context, playerID string) (*scene.Scene, error) {
	filter := bson.M{"players": playerID}
	var doc SceneDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查找玩家场景失败: %w", err)
	}
	return r.toDomain(&doc), nil
}

// GetScenePlayerCount 获取场景玩家数
func (r *MongoSceneRepository) GetScenePlayerCount(ctx context.Context, sceneID string) (int, error) {
	filter := bson.M{"scene_id": sceneID}
	var doc SceneDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return 0, err
	}
	return doc.CurrentPlayers, nil
}

// GetScenePlayers 获取场景玩家列表
func (r *MongoSceneRepository) GetScenePlayers(ctx context.Context, sceneID string) ([]string, error) {
	filter := bson.M{"scene_id": sceneID}
	var doc SceneDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return nil, err
	}
	return doc.Players, nil
}

// GetSceneStats 获取场景统计信息
func (r *MongoSceneRepository) GetSceneStats(ctx context.Context, sceneID string) (*scene.SceneStats, error) {
	// TODO: 实现场景统计信息
	return nil, nil
}

// GetSceneHistory 获取场景历史记录
func (r *MongoSceneRepository) GetSceneHistory(ctx context.Context, sceneID string, limit int) ([]*scene.SceneHistoryRecord, error) {
	// TODO: 实现场景历史记录
	return nil, nil
}

// GetPopularScenes 获取热门场景
func (r *MongoSceneRepository) GetPopularScenes(ctx context.Context, limit int) ([]*scene.ScenePopularity, error) {
	// TODO: 实现热门场景查询
	return nil, nil
}

// GetSceneConfig 获取场景配置
func (r *MongoSceneRepository) GetSceneConfig(ctx context.Context, sceneID string) (*scene.SceneConfig, error) {
	// TODO: 实现场景配置查询
	return nil, nil
}

// SaveSceneConfig 保存场景配置
func (r *MongoSceneRepository) SaveSceneConfig(ctx context.Context, config *scene.SceneConfig) error {
	// TODO: 实现场景配置保存
	return nil
}

// GetAllSceneConfigs 获取所有场景配置
func (r *MongoSceneRepository) GetAllSceneConfigs(ctx context.Context) ([]*scene.SceneConfig, error) {
	// TODO: 实现场景配置列表查询
	return nil, nil
}

// toDocument 转换为文档
func (r *MongoSceneRepository) toDocument(s *scene.Scene) *SceneDocument {
	return &SceneDocument{
		SceneID:        s.ID(),
		Name:           s.Name(),
		SceneType:      int(s.Type()),
		Status:         int(s.Status()),
		Width:          s.GetWidth(),
		Height:         s.GetHeight(),
		MaxPlayers:     s.GetMaxPlayers(),
		CurrentPlayers: s.PlayerCount(),
		Players:        []string{}, // TODO: 从场景中提取玩家ID列表
	}
}

// toDomain 转换为领域对象
func (r *MongoSceneRepository) toDomain(doc *SceneDocument) *scene.Scene {
	return scene.NewScene(
		doc.SceneID,
		doc.Name,
		scene.SceneType(doc.SceneType),
		doc.Width,
		doc.Height,
		doc.MaxPlayers,
	)
}
