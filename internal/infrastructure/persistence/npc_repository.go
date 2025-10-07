package persistence

import (
	"context"
	"fmt"
	"math"
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

// CountByRegion 根据区域计数
func (r *MongoNPCRepository) CountByRegion(region string) (int64, error) {
	ctx := context.Background()

	filter := bson.M{"location.region": region}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to count NPCs by region", "error", err, "region", region)
		return 0, fmt.Errorf("failed to count NPCs by region: %w", err)
	}

	return count, nil
}

// CountByStatus 根据状态统计NPC数量
func (r *MongoNPCRepository) CountByStatus(status npc.NPCStatus) (int64, error) {
	ctx := context.Background()

	filter := bson.M{"status": string(status)}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to count NPCs by status", "error", err, "status", status)
		return 0, fmt.Errorf("failed to count NPCs by status: %w", err)
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
	// Access patrol points directly from behavior
	for _, loc := range behavior.PatrolPoints {
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
			Health:       int64(attributes.GetHealth()),
			MaxHealth:    int64(attributes.GetMaxHealth()),
			Attack:       int64(attributes.GetAttack()),
			Defense:      int64(attributes.GetDefense()),
			Speed:        attributes.GetSpeed(),
			Intelligence: int64(attributes.GetIntelligence()),
			Charisma:     int64(attributes.GetIntelligence()), // Use Intelligence as fallback
		},
		Behavior: BehaviorDoc{
			Type:        string(behavior.Type),
			State:       string(behavior.State),
			LastAction:  behavior.LastMove,
			CooldownEnd: behavior.LastMove.Add(behavior.PauseTime),
			PatrolRoute: patrolRoute,
			Schedule:    make(map[string]string),
		},
		Dialogues: r.extractDialogueIDs(npcAggregate),
		Quests:    r.extractQuestIDs(npcAggregate),
		ShopID:    r.extractShopID(npcAggregate),
		Statistics: StatisticsDoc{
			TotalInteractions: 0, // NPCStatistics doesn't have this field
			TotalDialogues:    int64(statistics.TotalDialogues),
			TotalQuests:       int64(statistics.TotalQuests),
			TotalTrades:       0, // NPCStatistics doesn't have this field
			LastInteraction:   statistics.LastActiveAt,
			PopularityScore:   0.0, // NPCStatistics doesn't have this field
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

	// Parse NPC type and status
	npcType, err := r.parseNPCType(doc.Type)
	if err != nil {
		return nil
	}

	npcStatus, err := r.parseNPCStatus(doc.Status)
	if err != nil {
		return nil
	}

	// Create new NPC aggregate with basic info
	npcAggregate := npc.NewNPCAggregate(
		doc.NPCID,
		doc.Name,
		doc.Description,
		npcType,
	)

	// Set status
	err = npcAggregate.SetStatus(npcStatus)
	if err != nil {
		return nil
	}

	// Set location
	location := npc.NewLocation(doc.Location.X, doc.Location.Y, doc.Location.Z, doc.Location.Region, doc.Location.Zone)
	err = npcAggregate.MoveTo(location)
	if err != nil {
		return nil
	}

	return npcAggregate
}

// CreateIndexes 创建索引
func (r *MongoNPCRepository) CreateIndexes() error {
	ctx := context.Background()

	// 创建复合索引
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{"npc_id", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{"type", 1}, {"status", 1}},
		},
		{
			Keys: bson.D{{"location.region", 1}, {"location.zone", 1}},
		},
		{
			Keys: bson.D{{"status", 1}, {"updated_at", -1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		r.logger.Error("Failed to create indexes", "error", err)
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	r.logger.Info("NPC indexes created successfully")
	return nil
}

// Helper methods for extracting IDs

// extractDialogueIDs 提取对话ID列表
func (r *MongoNPCRepository) extractDialogueIDs(npcAggregate *npc.NPCAggregate) []string {
	dialogues := npcAggregate.GetAllDialogues()
	ids := make([]string, 0, len(dialogues))
	for id := range dialogues {
		ids = append(ids, id)
	}
	return ids
}

// extractQuestIDs 提取任务ID列表
func (r *MongoNPCRepository) extractQuestIDs(npcAggregate *npc.NPCAggregate) []string {
	quests := npcAggregate.GetAllQuests()
	ids := make([]string, 0, len(quests))
	for id := range quests {
		ids = append(ids, id)
	}
	return ids
}

// extractShopID 提取商店ID
func (r *MongoNPCRepository) extractShopID(npcAggregate *npc.NPCAggregate) string {
	shop := npcAggregate.GetShop()
	if shop == nil {
		return ""
	}
	return shop.GetID()
}

// parseNPCType 解析NPC类型
func (r *MongoNPCRepository) parseNPCType(typeStr string) (npc.NPCType, error) {
	switch typeStr {
	case "merchant":
		return npc.NPCTypeMerchant, nil
	case "guard":
		return npc.NPCTypeGuard, nil
	case "villager":
		return npc.NPCTypeVillager, nil
	case "quest_giver":
		return npc.NPCTypeQuestGiver, nil
	case "trainer":
		return npc.NPCTypeTrainer, nil
	default:
		return npc.NPCTypeVillager, nil // default type
	}
}

// parseNPCStatus 解析NPC状态
func (r *MongoNPCRepository) parseNPCStatus(statusStr string) (npc.NPCStatus, error) {
	switch statusStr {
	case "active":
		return npc.NPCStatusActive, nil
	case "inactive":
		return npc.NPCStatusInactive, nil
	case "hidden":
		return npc.NPCStatusHidden, nil
	case "busy":
		return npc.NPCStatusBusy, nil
	default:
		return npc.NPCStatusActive, nil // default status
	}
}

// DeleteBatch 批量删除NPC
func (r *MongoNPCRepository) DeleteBatch(ids []string) error {
	ctx := context.Background()

	filter := bson.M{"npc_id": bson.M{"$in": ids}}
	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete NPCs in batch", "error", err, "ids", ids)
		return fmt.Errorf("failed to delete NPCs in batch: %w", err)
	}

	r.logger.Info("NPCs deleted in batch", "count", result.DeletedCount, "ids", ids)
	return nil
}

// FindWithPagination 分页查找NPC
func (r *MongoNPCRepository) FindWithPagination(query *npc.NPCQuery) (*npc.NPCPageResult, error) {
	ctx := context.Background()

	// 构建查询条件
	filter := bson.M{}

	if query.Name != "" {
		filter["name"] = bson.M{"$regex": query.Name, "$options": "i"}
	}
	if query.Type != nil {
		filter["type"] = string(*query.Type)
	}
	if query.Status != nil {
		filter["status"] = string(*query.Status)
	}
	if query.Region != "" {
		filter["location.region"] = query.Region
	}
	if query.Zone != "" {
		filter["location.zone"] = query.Zone
	}

	// 计算总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to count NPCs", "error", err)
		return nil, fmt.Errorf("failed to count NPCs: %w", err)
	}

	// 构建查询选项
	opts := options.Find()
	if query.Limit > 0 {
		opts.SetLimit(int64(query.Limit))
	}
	if query.Offset > 0 {
		opts.SetSkip(int64(query.Offset))
	}
	if query.OrderBy != "" {
		order := 1
		if query.OrderDesc {
			order = -1
		}
		opts.SetSort(bson.D{{query.OrderBy, order}})
	}

	// 执行查询
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to find NPCs with pagination", "error", err)
		return nil, fmt.Errorf("failed to find NPCs with pagination: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []NPCDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode NPCs", "error", err)
		return nil, fmt.Errorf("failed to decode NPCs: %w", err)
	}

	// 转换为聚合根
	npcs := make([]*npc.NPCAggregate, 0, len(docs))
	for _, doc := range docs {
		npcAggregate := r.documentToAggregate(&doc)
		if npcAggregate != nil {
			npcs = append(npcs, npcAggregate)
		}
	}

	// 计算是否还有更多数据
	hasMore := int64(query.Offset+len(npcs)) < total

	return &npc.NPCPageResult{
		Items:   npcs,
		Total:   total,
		Offset:  query.Offset,
		Limit:   query.Limit,
		HasMore: hasMore,
	}, nil
}

// SaveBatch 批量保存NPC
func (r *MongoNPCRepository) SaveBatch(npcs []*npc.NPCAggregate) error {
	ctx := context.Background()

	if len(npcs) == 0 {
		return nil
	}

	// 准备批量操作
	var operations []mongo.WriteModel
	for _, npcAggregate := range npcs {
		doc := r.aggregateToDocument(npcAggregate)
		doc.UpdatedAt = time.Now()

		if doc.ID.IsZero() {
			doc.CreatedAt = time.Now()
			insertModel := mongo.NewInsertOneModel().SetDocument(doc)
			operations = append(operations, insertModel)
		} else {
			filter := bson.M{"npc_id": npcAggregate.GetID()}
			update := bson.M{"$set": doc}
			updateModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
			operations = append(operations, updateModel)
		}
	}

	// 执行批量操作
	result, err := r.collection.BulkWrite(ctx, operations)
	if err != nil {
		r.logger.Error("Failed to save NPCs in batch", "error", err)
		return fmt.Errorf("failed to save NPCs in batch: %w", err)
	}

	r.logger.Info("NPCs saved in batch", "inserted", result.InsertedCount, "modified", result.ModifiedCount)
	return nil
}

// FindNearbyNPCs 查找附近的NPC
func (r *MongoNPCRepository) FindNearbyNPCs(location *npc.Location, radius float64, npcType npc.NPCType) ([]*npc.NPCAggregate, error) {
	ctx := context.Background()

	// 构建查询条件
	filter := bson.M{
		"location.region": location.GetRegion(),
		"location.zone":   location.GetZone(),
	}

	// 如果指定了NPC类型，添加类型过滤
	if string(npcType) != "" {
		filter["type"] = string(npcType)
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find nearby NPCs", "error", err, "location", location, "radius", radius)
		return nil, fmt.Errorf("failed to find nearby NPCs: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []NPCDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode nearby NPCs", "error", err)
		return nil, fmt.Errorf("failed to decode nearby NPCs: %w", err)
	}

	// 过滤距离范围内的NPC
	npcs := make([]*npc.NPCAggregate, 0)
	for _, doc := range docs {
		// 计算距离
		dx := doc.Location.X - location.GetX()
		dy := doc.Location.Y - location.GetY()
		dz := doc.Location.Z - location.GetZ()
		distance := math.Sqrt(dx*dx + dy*dy + dz*dz)

		if distance <= radius {
			npcAggregate := r.documentToAggregate(&doc)
			if npcAggregate != nil {
				npcs = append(npcs, npcAggregate)
			}
		}
	}

	return npcs, nil
}



// FindByZone 根据区域查找NPC
func (r *MongoNPCRepository) FindByZone(zone string) ([]*npc.NPCAggregate, error) {
	ctx := context.Background()

	filter := bson.M{"location.zone": zone}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find NPCs by zone", "error", err, "zone", zone)
		return nil, fmt.Errorf("failed to find NPCs by zone: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []NPCDocument
	if err := cursor.All(ctx, &docs); err != nil {
		r.logger.Error("Failed to decode NPCs by zone", "error", err, "zone", zone)
		return nil, fmt.Errorf("failed to decode NPCs by zone: %w", err)
	}

	npcs := make([]*npc.NPCAggregate, 0, len(docs))
	for _, doc := range docs {
		npcAggregate := r.documentToAggregate(&doc)
		if npcAggregate != nil {
			npcs = append(npcs, npcAggregate)
		}
	}

	return npcs, nil
}
