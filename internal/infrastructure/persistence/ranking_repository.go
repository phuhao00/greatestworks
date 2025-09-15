package persistence

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/greatestworks/internal/domain/ranking"
)

// RankingRepository MongoDB排行榜仓储实现
type RankingRepository struct {
	db          *mongo.Database
	rankingColl *mongo.Collection
	entryColl   *mongo.Collection
}

// NewRankingRepository 创建排行榜仓储
func NewRankingRepository(db *mongo.Database) *RankingRepository {
	return &RankingRepository{
		db:          db,
		rankingColl: db.Collection("rankings"),
		entryColl:   db.Collection("rank_entries"),
	}
}

// RankingDocument 排行榜文档结构
type RankingDocument struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	RankingID   string             `bson:"ranking_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	RankType    string             `bson:"rank_type"`
	PeriodType  string             `bson:"period_type"`
	MaxEntries  int32              `bson:"max_entries"`
	IsActive    bool               `bson:"is_active"`
	Blacklist   []uint64           `bson:"blacklist"`
	Settings    map[string]interface{} `bson:"settings"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
	ResetAt     *time.Time         `bson:"reset_at,omitempty"`
	Version     int64              `bson:"version"`
}

// RankEntryDocument 排名条目文档结构
type RankEntryDocument struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty"`
	EntryID   string                 `bson:"entry_id"`
	RankingID string                 `bson:"ranking_id"`
	PlayerID  uint64                 `bson:"player_id"`
	Rank      int32                  `bson:"rank"`
	Score     int64                  `bson:"score"`
	PrevRank  int32                  `bson:"prev_rank"`
	PrevScore int64                  `bson:"prev_score"`
	Metadata  map[string]interface{} `bson:"metadata"`
	CreatedAt time.Time              `bson:"created_at"`
	UpdatedAt time.Time              `bson:"updated_at"`
}

// Save 保存排行榜聚合根
func (r *RankingRepository) Save(ctx context.Context, rankingAggregate *ranking.RankingAggregate) error {
	doc := r.toRankingDocument(rankingAggregate)
	
	filter := bson.M{"ranking_id": doc.RankingID}
	update := bson.M{
		"$set": doc,
		"$inc": bson.M{"version": 1},
	}
	opts := options.Update().SetUpsert(true)
	
	_, err := r.rankingColl.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save ranking: %w", err)
	}
	
	return nil
}

// FindByID 根据ID查找排行榜
func (r *RankingRepository) FindByID(ctx context.Context, rankingID string) (*ranking.RankingAggregate, error) {
	filter := bson.M{"ranking_id": rankingID}
	
	var doc RankingDocument
	err := r.rankingColl.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find ranking: %w", err)
	}
	
	return r.fromRankingDocument(&doc), nil
}

// FindByType 根据类型查找排行榜
func (r *RankingRepository) FindByType(ctx context.Context, rankType ranking.RankType) ([]*ranking.RankingAggregate, error) {
	filter := bson.M{
		"rank_type": rankType.String(),
		"is_active": true,
	}
	
	cursor, err := r.rankingColl.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find rankings by type: %w", err)
	}
	defer cursor.Close(ctx)
	
	var rankings []*ranking.RankingAggregate
	for cursor.Next(ctx) {
		var doc RankingDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode ranking document: %w", err)
		}
		rankings = append(rankings, r.fromRankingDocument(&doc))
	}
	
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	
	return rankings, nil
}

// FindActive 查找激活的排行榜
func (r *RankingRepository) FindActive(ctx context.Context) ([]*ranking.RankingAggregate, error) {
	filter := bson.M{"is_active": true}
	
	cursor, err := r.rankingColl.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find active rankings: %w", err)
	}
	defer cursor.Close(ctx)
	
	var rankings []*ranking.RankingAggregate
	for cursor.Next(ctx) {
		var doc RankingDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode ranking document: %w", err)
		}
		rankings = append(rankings, r.fromRankingDocument(&doc))
	}
	
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	
	return rankings, nil
}

// Delete 删除排行榜
func (r *RankingRepository) Delete(ctx context.Context, rankingID string) error {
	filter := bson.M{"ranking_id": rankingID}
	update := bson.M{
		"$set": bson.M{
			"is_active": false,
			"updated_at": time.Now(),
		},
		"$inc": bson.M{"version": 1},
	}
	
	_, err := r.rankingColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete ranking: %w", err)
	}
	
	return nil
}

// RankEntryRepository 排名条目仓储实现
type RankEntryRepository struct {
	db        *mongo.Database
	entryColl *mongo.Collection
}

// NewRankEntryRepository 创建排名条目仓储
func NewRankEntryRepository(db *mongo.Database) *RankEntryRepository {
	return &RankEntryRepository{
		db:        db,
		entryColl: db.Collection("rank_entries"),
	}
}

// Save 保存排名条目
func (r *RankEntryRepository) Save(ctx context.Context, entry *ranking.RankEntry) error {
	doc := r.toRankEntryDocument(entry)
	
	filter := bson.M{"entry_id": doc.EntryID}
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)
	
	_, err := r.entryColl.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save rank entry: %w", err)
	}
	
	return nil
}

// FindByID 根据ID查找排名条目
func (r *RankEntryRepository) FindByID(ctx context.Context, entryID string) (*ranking.RankEntry, error) {
	filter := bson.M{"entry_id": entryID}
	
	var doc RankEntryDocument
	err := r.entryColl.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find rank entry: %w", err)
	}
	
	return r.fromRankEntryDocument(&doc), nil
}

// FindByRankingAndPlayer 根据排行榜和玩家查找条目
func (r *RankEntryRepository) FindByRankingAndPlayer(ctx context.Context, rankingID string, playerID uint64) (*ranking.RankEntry, error) {
	filter := bson.M{
		"ranking_id": rankingID,
		"player_id":  playerID,
	}
	
	var doc RankEntryDocument
	err := r.entryColl.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find rank entry: %w", err)
	}
	
	return r.fromRankEntryDocument(&doc), nil
}

// FindByRanking 根据排行榜查找条目
func (r *RankEntryRepository) FindByRanking(ctx context.Context, rankingID string, limit int) ([]*ranking.RankEntry, error) {
	filter := bson.M{"ranking_id": rankingID}
	opts := options.Find().
		SetSort(bson.D{{Key: "rank", Value: 1}}).
		SetLimit(int64(limit))
	
	cursor, err := r.entryColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find rank entries: %w", err)
	}
	defer cursor.Close(ctx)
	
	var entries []*ranking.RankEntry
	for cursor.Next(ctx) {
		var doc RankEntryDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode rank entry document: %w", err)
		}
		entries = append(entries, r.fromRankEntryDocument(&doc))
	}
	
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	
	return entries, nil
}

// FindByQuery 根据查询条件查找条目
func (r *RankEntryRepository) FindByQuery(ctx context.Context, query *ranking.RankEntryQuery) ([]*ranking.RankEntry, int64, error) {
	filter := r.buildRankEntryFilter(query)
	
	// 计算总数
	total, err := r.entryColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count rank entries: %w", err)
	}
	
	// 构建查询选项
	opts := options.Find()
	if query.GetSort() != "" {
		sortOrder := 1
		if query.GetSortOrder() == "desc" {
			sortOrder = -1
		}
		opts.SetSort(bson.D{{Key: query.GetSort(), Value: sortOrder}})
	}
	if query.GetLimit() > 0 {
		opts.SetLimit(int64(query.GetLimit()))
	}
	if query.GetOffset() > 0 {
		opts.SetSkip(int64(query.GetOffset()))
	}
	
	cursor, err := r.entryColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find rank entries: %w", err)
	}
	defer cursor.Close(ctx)
	
	var entries []*ranking.RankEntry
	for cursor.Next(ctx) {
		var doc RankEntryDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, fmt.Errorf("failed to decode rank entry document: %w", err)
		}
		entries = append(entries, r.fromRankEntryDocument(&doc))
	}
	
	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor error: %w", err)
	}
	
	return entries, total, nil
}

// DeleteByRanking 删除排行榜的所有条目
func (r *RankEntryRepository) DeleteByRanking(ctx context.Context, rankingID string) (int64, error) {
	filter := bson.M{"ranking_id": rankingID}
	
	result, err := r.entryColl.DeleteMany(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to delete rank entries: %w", err)
	}
	
	return result.DeletedCount, nil
}

// UpdateRanks 批量更新排名
func (r *RankEntryRepository) UpdateRanks(ctx context.Context, rankingID string) error {
	// 使用聚合管道重新计算排名
	pipeline := []bson.M{
		{"$match": bson.M{"ranking_id": rankingID}},
		{"$sort": bson.M{"score": -1, "updated_at": 1}},
		{"$group": bson.M{
			"_id": "$ranking_id",
			"entries": bson.M{"$push": "$$ROOT"},
		}},
		{"$unwind": bson.M{
			"path":              "$entries",
			"includeArrayIndex": "rank",
		}},
		{"$addFields": bson.M{
			"entries.prev_rank": "$entries.rank",
			"entries.rank":      bson.M{"$add": []interface{}{"$rank", 1}},
			"entries.updated_at": time.Now(),
		}},
		{"$replaceRoot": bson.M{"newRoot": "$entries"}},
		{"$merge": bson.M{
			"into": "rank_entries",
			"on":   "_id",
			"whenMatched": "replace",
		}},
	}
	
	_, err := r.entryColl.Aggregate(ctx, pipeline)
	if err != nil {
		return fmt.Errorf("failed to update ranks: %w", err)
	}
	
	return nil
}

// 私有方法

// toRankingDocument 转换为排行榜文档
func (r *RankingRepository) toRankingDocument(rankingAggregate *ranking.RankingAggregate) *RankingDocument {
	doc := &RankingDocument{
		RankingID:   rankingAggregate.GetID(),
		Name:        rankingAggregate.GetName(),
		Description: rankingAggregate.GetDescription(),
		RankType:    rankingAggregate.GetRankType().String(),
		PeriodType:  rankingAggregate.GetPeriodType().String(),
		MaxEntries:  rankingAggregate.GetMaxEntries(),
		IsActive:    rankingAggregate.IsActive(),
		Blacklist:   rankingAggregate.GetBlacklist(),
		Settings:    rankingAggregate.GetSettings(),
		CreatedAt:   rankingAggregate.GetCreatedAt(),
		UpdatedAt:   rankingAggregate.GetUpdatedAt(),
		Version:     rankingAggregate.GetVersion(),
	}
	
	if !rankingAggregate.GetResetAt().IsZero() {
		resetAt := rankingAggregate.GetResetAt()
		doc.ResetAt = &resetAt
	}
	
	return doc
}

// fromRankingDocument 从排行榜文档转换
func (r *RankingRepository) fromRankingDocument(doc *RankingDocument) *ranking.RankingAggregate {
	// 解析枚举值
	rankType := ranking.ParseRankType(doc.RankType)
	periodType := ranking.ParsePeriodType(doc.PeriodType)
	
	// 重建聚合根
	rankingAggregate := ranking.NewRankingAggregate(doc.Name, rankType, periodType)
	rankingAggregate.SetID(doc.RankingID)
	rankingAggregate.SetDescription(doc.Description)
	rankingAggregate.SetMaxEntries(doc.MaxEntries)
	rankingAggregate.SetBlacklist(doc.Blacklist)
	rankingAggregate.SetSettings(doc.Settings)
	rankingAggregate.SetVersion(doc.Version)
	
	if doc.IsActive {
		rankingAggregate.Activate()
	} else {
		rankingAggregate.Deactivate()
	}
	
	if doc.ResetAt != nil {
		rankingAggregate.SetResetAt(*doc.ResetAt)
	}
	
	return rankingAggregate
}

// toRankEntryDocument 转换为排名条目文档
func (r *RankEntryRepository) toRankEntryDocument(entry *ranking.RankEntry) *RankEntryDocument {
	return &RankEntryDocument{
		EntryID:   entry.GetID(),
		RankingID: entry.GetRankingID(),
		PlayerID:  entry.GetPlayerID(),
		Rank:      entry.GetRank(),
		Score:     entry.GetScore(),
		PrevRank:  entry.GetPrevRank(),
		PrevScore: entry.GetPrevScore(),
		Metadata:  entry.GetMetadata(),
		CreatedAt: entry.GetCreatedAt(),
		UpdatedAt: entry.GetUpdatedAt(),
	}
}

// fromRankEntryDocument 从排名条目文档转换
func (r *RankEntryRepository) fromRankEntryDocument(doc *RankEntryDocument) *ranking.RankEntry {
	entry := ranking.NewRankEntry(
		doc.EntryID,
		doc.RankingID,
		doc.PlayerID,
		doc.Score,
	)
	
	entry.SetRank(doc.Rank)
	entry.SetPrevious(doc.PrevRank, doc.PrevScore)
	entry.SetMetadata(doc.Metadata)
	
	return entry
}

// buildRankEntryFilter 构建排名条目查询过滤器
func (r *RankEntryRepository) buildRankEntryFilter(query *ranking.RankEntryQuery) bson.M {
	filter := bson.M{}
	
	if query.GetRankingID() != "" {
		filter["ranking_id"] = query.GetRankingID()
	}
	
	if query.GetPlayerID() > 0 {
		filter["player_id"] = query.GetPlayerID()
	}
	
	if query.GetMinRank() > 0 {
		filter["rank"] = bson.M{"$gte": query.GetMinRank()}
	}
	
	if query.GetMaxRank() > 0 {
		if rankFilter, exists := filter["rank"]; exists {
			rankFilter.(bson.M)["$lte"] = query.GetMaxRank()
		} else {
			filter["rank"] = bson.M{"$lte": query.GetMaxRank()}
		}
	}
	
	if query.GetMinScore() > 0 {
		filter["score"] = bson.M{"$gte": query.GetMinScore()}
	}
	
	if query.GetMaxScore() > 0 {
		if scoreFilter, exists := filter["score"]; exists {
			scoreFilter.(bson.M)["$lte"] = query.GetMaxScore()
		} else {
			filter["score"] = bson.M{"$lte": query.GetMaxScore()}
		}
	}
	
	return filter
}

// CreateIndexes 创建索引
func (r *RankingRepository) CreateIndexes(ctx context.Context) error {
	// 排行榜索引
	rankingIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "ranking_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "rank_type", Value: 1}, {Key: "is_active", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "period_type", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_active", Value: 1}},
		},
	}
	
	if _, err := r.rankingColl.Indexes().CreateMany(ctx, rankingIndexes); err != nil {
		return fmt.Errorf("failed to create ranking indexes: %w", err)
	}
	
	return nil
}

// CreateIndexes 创建排名条目索引
func (r *RankEntryRepository) CreateIndexes(ctx context.Context) error {
	// 排名条目索引
	entryIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "entry_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "ranking_id", Value: 1}, {Key: "player_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "ranking_id", Value: 1}, {Key: "rank", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "ranking_id", Value: 1}, {Key: "score", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "player_id", Value: 1}},
		},
	}
	
	if _, err := r.entryColl.Indexes().CreateMany(ctx, entryIndexes); err != nil {
		return fmt.Errorf("failed to create rank entry indexes: %w", err)
	}
	
	return nil
}