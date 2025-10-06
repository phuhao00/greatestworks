package persistence

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"greatestworks/internal/domain/npc"
	"greatestworks/internal/infrastructure/cache"
	"greatestworks/internal/infrastructure/logger"
)

// MongoNPCRepository MongoDB NPC仓储实现
type MongoNPCRepository struct {
	db         *mongo.Database
	cache      cache.Cache
	logger     logger.Logger
	collection *mongo.Collection
}

// NewMongoNPCRepository 创建MongoDB NPC仓储
func NewMongoNPCRepository(db *mongo.Database, cache cache.Cache, logger logger.Logger) npc.NPCRepository {
	return &MongoNPCRepository{
		db:         db,
		cache:      cache,
		logger:     logger,
		collection: db.Collection("npcs"),
	}
}

// NPCDocument MongoDB NPC文档结构
type NPCDocument struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	NPCID       string             `bson:"npc_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Type        string             `bson:"type"`
	Status      string             `bson:"status"`
	Location    LocationDoc        `bson:"location"`
	Attributes  NPCAttributesDoc   `bson:"attributes"`
	Behavior    BehaviorDoc        `bson:"behavior"`
	Dialogues   []string           `bson:"dialogues"`
	Quests      []string           `bson:"quests"`
	ShopID      string             `bson:"shop_id,omitempty"`
	Statistics  StatisticsDoc      `bson:"statistics"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

// LocationDoc 位置文档
type LocationDoc struct {
	X      float64 `bson:"x"`
	Y      float64 `bson:"y"`
	Z      float64 `bson:"z"`
	Region string  `bson:"region"`
	Zone   string  `bson:"zone"`
}

// NPCAttributesDoc NPC属性文档
type NPCAttributesDoc struct {
	Level        int     `bson:"level"`
	Health       int64   `bson:"health"`
	MaxHealth    int64   `bson:"max_health"`
	Attack       int64   `bson:"attack"`
	Defense      int64   `bson:"defense"`
	Speed        float64 `bson:"speed"`
	Intelligence int64   `bson:"intelligence"`
	Charisma     int64   `bson:"charisma"`
}

// BehaviorDoc 行为文档
type BehaviorDoc struct {
	Type        string            `bson:"type"`
	State       string            `bson:"state"`
	LastAction  time.Time         `bson:"last_action"`
	CooldownEnd time.Time         `bson:"cooldown_end"`
	PatrolRoute []LocationDoc     `bson:"patrol_route"`
	Schedule    map[string]string `bson:"schedule"`
}

// StatisticsDoc 统计文档
type StatisticsDoc struct {
	TotalInteractions int64     `bson:"total_interactions"`
	TotalDialogues    int64     `bson:"total_dialogues"`
	TotalQuests       int64     `bson:"total_quests"`
	TotalTrades       int64     `bson:"total_trades"`
	LastInteraction   time.Time `bson:"last_interaction"`
	PopularityScore   float64   `bson:"popularity_score"`
}

// Save 保存NPC
func (r *MongoNPCRepository) Save(npcAggregate *npc.NPCAggregate) error {
	ctx := context.Background()
	doc := r.aggregateToDocument(npcAggregate)
	doc.UpdatedAt = time.Now()

	if doc.ID.IsZero() {
		doc.CreatedAt = time.Now()
		result, err := r.collection.InsertOne(ctx, doc)
		if err != nil {
			r.logger.Error("Failed to insert NPC", "error", err, "npc_id", npcAggregate.GetID())
			return fmt.Errorf("failed to insert NPC: %w", err)
		}

		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
			doc.ID = oid
		}
	} else {
		filter := bson.M{"npc_id": npcAggregate.GetID()}
		update := bson.M{"$set": doc}

		_, err := r.collection.UpdateOne(ctx, filter, update)
		if err != nil {
			r.logger.Error("Failed to update NPC", "error", err, "npc_id", npcAggregate.GetID())
			return fmt.Errorf("failed to update NPC: %w", err)
		}
	}

	// 更新缓存
	cacheKey := fmt.Sprintf("npc:%s", npcAggregate.GetID())
	if err := r.cache.Set(ctx, cacheKey, npcAggregate, time.Hour); err != nil {
		r.logger.Warn("Failed to cache NPC", "error", err, "npc_id", npcAggregate.GetID())
	}

	return nil
}

// FindByID 根据ID查找NPC
func (r *MongoNPCRepository) FindByID(npcID string) (*npc.NPCAggregate, error) {
	ctx := context.Background()

	// 先从缓存获取
	cacheKey := fmt.Sprintf("npc:%s", npcID)
	var cachedNPC *npc.NPCAggregate
	if err := r.cache.Get(ctx, cacheKey, &cachedNPC); err == nil && cachedNPC != nil {
		return cachedNPC, nil
	}

	// 从数据库获取
	filter := bson.M{"npc_id": npcID}
	var doc NPCDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, npc.ErrNPCNotFound
		}
		r.logger.Error("Failed to find NPC", "error", err, "npc_id", npcID)
		return nil, fmt.Errorf("failed to find NPC: %w", err)
	}

	npcAggregate := r.documentToAggregate(&doc)

	// 更新缓存
	if err := r.cache.Set(ctx, cacheKey, npcAggregate, time.Hour); err != nil {
		r.logger.Warn("Failed to cache NPC", "error", err, "npc_id", npcID)
	}

	return npcAggregate, nil
}

// FindByType 根据类型查找NPC
func (r *MongoNPCRepository) FindByType(npcType npc.NPCType) ([]*npc.NPCAggregate, error) {
	ctx := context.Background()

	filter := bson.M{"type": string(npcType)}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find NPCs by type", "error", err, "type", npcType)
		return nil, fmt.Errorf("failed to find NPCs by type: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []NPCDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode NPCs by type", "error", err, "type", npcType)
		return nil, fmt.Errorf("failed to decode NPCs by type: %w", err)
	}

	npcs := make([]*npc.NPCAggregate, len(docs))
	for i, doc := range docs {
		npcs[i] = r.documentToAggregate(&doc)
	}

	return npcs, nil
}

// FindByStatus 根据状态查找NPC
func (r *MongoNPCRepository) FindByStatus(status npc.NPCStatus) ([]*npc.NPCAggregate, error) {
	ctx := context.Background()

	filter := bson.M{"status": string(status)}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find NPCs by status", "error", err, "status", status)
		return nil, fmt.Errorf("failed to find NPCs by status: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []NPCDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode NPCs by status", "error", err, "status", status)
		return nil, fmt.Errorf("failed to decode NPCs by status: %w", err)
	}

	npcs := make([]*npc.NPCAggregate, len(docs))
	for i, doc := range docs {
		npcs[i] = r.documentToAggregate(&doc)
	}

	return npcs, nil
}

// FindByLocation 根据位置查找NPC
func (r *MongoNPCRepository) FindByLocation(location *npc.Location, radius float64) ([]*npc.NPCAggregate, error) {
	ctx := context.Background()

	// 使用地理位置查询
	filter := bson.M{
		"location.x": bson.M{
			"$gte": location.GetX() - radius,
			"$lte": location.GetX() + radius,
		},
		"location.y": bson.M{
			"$gte": location.GetY() - radius,
			"$lte": location.GetY() + radius,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find NPCs by location", "error", err)
		return nil, fmt.Errorf("failed to find NPCs by location: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []NPCDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode NPCs by location", "error", err)
		return nil, fmt.Errorf("failed to decode NPCs by location: %w", err)
	}

	npcs := make([]*npc.NPCAggregate, len(docs))
	for i, doc := range docs {
		npcs[i] = r.documentToAggregate(&doc)
	}

	return npcs, nil
}

// FindByRegion 根据区域查找NPC
func (r *MongoNPCRepository) FindByRegion(region string) ([]*npc.NPCAggregate, error) {
	ctx := context.Background()

	filter := bson.M{"location.region": region}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find NPCs by region", "error", err, "region", region)
		return nil, fmt.Errorf("failed to find NPCs by region: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []NPCDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode NPCs by region", "error", err, "region", region)
		return nil, fmt.Errorf("failed to decode NPCs by region: %w", err)
	}

	npcs := make([]*npc.NPCAggregate, len(docs))
	for i, doc := range docs {
		npcs[i] = r.documentToAggregate(&doc)
	}

	return npcs, nil
}

// Update 更新NPC
func (r *MongoNPCRepository) Update(npcAggregate *npc.NPCAggregate) error {
	return r.Save(npcAggregate)
}

// Delete 删除NPC
func (r *MongoNPCRepository) Delete(npcID string) error {
	ctx := context.Background()

	filter := bson.M{"npc_id": npcID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete NPC", "error", err, "npc_id", npcID)
		return fmt.Errorf("failed to delete NPC: %w", err)
	}

	if result.DeletedCount == 0 {
		return npc.ErrNPCNotFound
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("npc:%s", npcID)
	if err := r.cache.Delete(ctx, cacheKey); err != nil {
		r.logger.Warn("Failed to delete NPC cache", "error", err, "npc_id", npcID)
	}

	return nil
}

// FindActiveNPCs 查找活跃的NPC
func (r *MongoNPCRepository) FindActiveNPCs() ([]*npc.NPCAggregate, error) {
	ctx := context.Background()

	filter := bson.M{"status": string(npc.NPCStatusActive)}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find active NPCs", "error", err)
		return nil, fmt.Errorf("failed to find active NPCs: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []NPCDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode active NPCs", "error", err)
		return nil, fmt.Errorf("failed to decode active NPCs: %w", err)
	}

	npcs := make([]*npc.NPCAggregate, len(docs))
	for i, doc := range docs {
		npcs[i] = r.documentToAggregate(&doc)
	}

	return npcs, nil
}

// FindNPCsWithShops 查找有商店的NPC
func (r *MongoNPCRepository) FindNPCsWithShops() ([]*npc.NPCAggregate, error) {
	ctx := context.Background()

	filter := bson.M{"shop_id": bson.M{"$exists": true, "$ne": ""}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find NPCs with shops", "error", err)
		return nil, fmt.Errorf("failed to find NPCs with shops: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []NPCDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode NPCs with shops", "error", err)
		return nil, fmt.Errorf("failed to decode NPCs with shops: %w", err)
	}

	npcs := make([]*npc.NPCAggregate, len(docs))
	for i, doc := range docs {
		npcs[i] = r.documentToAggregate(&doc)
	}

	return npcs, nil
}

// FindNPCsWithQuests 查找有任务的NPC
func (r *MongoNPCRepository) FindNPCsWithQuests() ([]*npc.NPCAggregate, error) {
	ctx := context.Background()

	filter := bson.M{"quests": bson.M{"$exists": true, "$not": bson.M{"$size": 0}}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find NPCs with quests", "error", err)
		return nil, fmt.Errorf("failed to find NPCs with quests: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []NPCDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode NPCs with quests", "error", err)
		return nil, fmt.Errorf("failed to decode NPCs with quests: %w", err)
	}

	npcs := make([]*npc.NPCAggregate, len(docs))
	for i, doc := range docs {
		npcs[i] = r.documentToAggregate(&doc)
	}

	return npcs, nil
}

// Count 计数查询
func (r *MongoNPCRepository) Count() (int64, error) {
	ctx := context.Background()

	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		r.logger.Error("Failed to count NPCs", "error", err)
		return 0, fmt.Errorf("failed to count NPCs: %w", err)
	}

	return count, nil
}

// CountByType 根据类型计数
func (r *MongoNPCRepository) CountByType(npcType npc.NPCType) (int64, error) {
	ctx := context.Background()

	filter := bson.M{"type": string(npcType)}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to count NPCs by type", "error", err, "type", npcType)
		return 0, fmt.Errorf("failed to count NPCs by type: %w", err)
	}

	return count, nil
}

// 私有方法

// aggregateToDocument 聚合根转文档
func (r *MongoNPCRepository) aggregateToDocument(npcAggregate *npc.NPCAggregate) *NPCDocument {
	location := npcAggregate.GetLocation()
	attributes := npcAggregate.GetAttributes()
	behavior := npcAggregate.GetBehavior()
	statistics := npcAggregate.GetStatistics()

	// 转换巡逻路线
	patrolRoute := make([]LocationDoc, 0)
	for _, loc := range behavior.GetPatrolRoute() {
		patrolRoute = append(patrolRoute, LocationDoc{
			X:      loc.GetX(),
			Y:      loc.GetY(),
			Z:      loc.GetZ(),
			Region: loc.GetRegion(),
			Zone:   loc.GetZone(),
		})
	}

	return &NPCDocument{
		NPCID:       npcAggregate.GetID(),
		Name:        npcAggregate.GetName(),
		Description: npcAggregate.GetDescription(),
		Type:        string(npcAggregate.GetType()),
		Status:      string(npcAggregate.GetStatus()),
		Location: LocationDoc{
			X:      location.GetX(),
			Y:      location.GetY(),
			Z:      location.GetZ(),
			Region: location.GetRegion(),
			Zone:   location.GetZone(),
		},
		Attributes: NPCAttributesDoc{
			Level:        attributes.GetLevel(),
			Health:       attributes.GetHealth(),
			MaxHealth:    int64(attributes.GetMaxHealth()),
			Attack:       attributes.GetAttack(),
			Defense:      attributes.GetDefense(),
			Speed:        attributes.GetSpeed(),
			Intelligence: int64(attributes.GetIntelligence()),
			Charisma:     attributes.GetCharisma(),
		},
		Behavior: BehaviorDoc{
			Type:        string(behavior.GetType()),
			State:       string(behavior.GetState()),
			LastAction:  behavior.GetLastAction(),
			CooldownEnd: behavior.GetCooldownEnd(),
			PatrolRoute: patrolRoute,
			Schedule:    behavior.GetSchedule(),
		},
		Dialogues: npcAggregate.GetDialogueIDs(),
		Quests:    npcAggregate.GetQuestIDs(),
		ShopID:    npcAggregate.GetShopID(),
		Statistics: StatisticsDoc{
			TotalInteractions: statistics.GetTotalInteractions(),
			TotalDialogues:    statistics.GetTotalDialogues(),
			TotalQuests:       statistics.GetTotalQuests(),
			TotalTrades:       statistics.GetTotalTrades(),
			LastInteraction:   statistics.GetLastInteraction(),
			PopularityScore:   statistics.GetPopularityScore(),
		},
		CreatedAt: npcAggregate.GetCreatedAt(),
		UpdatedAt: npcAggregate.GetUpdatedAt(),
	}
}

// documentToAggregate 文档转聚合根
func (r *MongoNPCRepository) documentToAggregate(doc *NPCDocument) *npc.NPCAggregate {
	// 转换巡逻路线
	patrolRoute := make([]*npc.Location, len(doc.Behavior.PatrolRoute))
	for i, loc := range doc.Behavior.PatrolRoute {
		patrolRoute[i] = npc.NewLocation(loc.X, loc.Y, loc.Z, loc.Region, loc.Zone)
	}

	// 这里需要根据实际的NPCAggregate构造函数来实现
	return npc.ReconstructNPCAggregate(
		doc.NPCID,
		doc.Name,
		doc.Description,
		npc.NPCType(doc.Type),
		npc.NPCStatus(doc.Status),
		npc.NewLocation(doc.Location.X, doc.Location.Y, doc.Location.Z, doc.Location.Region, doc.Location.Zone),
		npc.NewNPCAttributes(
			doc.Attributes.Level,
			doc.Attributes.Health,
			doc.Attributes.MaxHealth,
			doc.Attributes.Attack,
			doc.Attributes.Defense,
			doc.Attributes.Speed,
			doc.Attributes.Intelligence,
			doc.Attributes.Charisma,
		),
		npc.NewNPCBehavior(
			npc.BehaviorType(doc.Behavior.Type),
			npc.BehaviorState(doc.Behavior.State),
			doc.Behavior.LastAction,
			doc.Behavior.CooldownEnd,
			patrolRoute,
			doc.Behavior.Schedule,
		),
		doc.Dialogues,
		doc.Quests,
		doc.ShopID,
		npc.NewNPCStatistics(
			doc.Statistics.TotalInteractions,
			doc.Statistics.TotalDialogues,
			doc.Statistics.TotalQuests,
			doc.Statistics.TotalTrades,
			doc.Statistics.LastInteraction,
			doc.Statistics.PopularityScore,
		),
		doc.CreatedAt,
		doc.UpdatedAt,
	)
}

// CreateIndexes 创建索引
func (r *MongoNPCRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "npc_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "type", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "location.region", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "location.zone", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "shop_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "quests", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "dialogues", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "location.x", Value: 1}, {Key: "location.y", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "type", Value: 1}, {Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "statistics.popularity_score", Value: -1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		r.logger.Error("Failed to create NPC indexes", "error", err)
		return fmt.Errorf("failed to create NPC indexes: %w", err)
	}

	r.logger.Info("NPC indexes created successfully")
	return nil
}
