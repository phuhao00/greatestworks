package persistence

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/greatestworks/internal/domain/minigame"
)

// MinigameRepository MongoDB小游戏仓储实现
type MinigameRepository struct {
	db           *mongo.Database
	minigameColl *mongo.Collection
	sessionColl  *mongo.Collection
}

// NewMinigameRepository 创建小游戏仓储
func NewMinigameRepository(db *mongo.Database) *MinigameRepository {
	return &MinigameRepository{
		db:           db,
		minigameColl: db.Collection("minigames"),
		sessionColl:  db.Collection("game_sessions"),
	}
}

// MinigameDocument 小游戏文档结构
type MinigameDocument struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	MinigameID  string             `bson:"minigame_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	GameType    string             `bson:"game_type"`
	Difficulty  string             `bson:"difficulty"`
	MaxPlayers  int32              `bson:"max_players"`
	TimeLimit   int32              `bson:"time_limit"`
	IsActive    bool               `bson:"is_active"`
	Rules       map[string]interface{} `bson:"rules"`
	Rewards     []RewardDocument   `bson:"rewards"`
	Settings    map[string]interface{} `bson:"settings"`
	Statistics  StatisticsDocument `bson:"statistics"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
	Version     int64              `bson:"version"`
}

// RewardDocument 奖励文档结构
type RewardDocument struct {
	Type     string `bson:"type"`
	ItemID   string `bson:"item_id,omitempty"`
	Quantity int32  `bson:"quantity"`
	Reason   string `bson:"reason"`
}

// StatisticsDocument 统计文档结构
type StatisticsDocument struct {
	TotalPlays    int64   `bson:"total_plays"`
	TotalPlayers  int64   `bson:"total_players"`
	AverageScore  float64 `bson:"average_score"`
	HighestScore  int64   `bson:"highest_score"`
	AverageTime   float64 `bson:"average_time"`
	CompletionRate float64 `bson:"completion_rate"`
}

// GameSessionDocument 游戏会话文档结构
type GameSessionDocument struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty"`
	SessionID   string                 `bson:"session_id"`
	MinigameID  string                 `bson:"minigame_id"`
	PlayerID    uint64                 `bson:"player_id"`
	Status      string                 `bson:"status"`
	Score       int64                  `bson:"score"`
	TimeLimit   int32                  `bson:"time_limit"`
	TimeElapsed int32                  `bson:"time_elapsed"`
	Settings    map[string]interface{} `bson:"settings"`
	GameData    map[string]interface{} `bson:"game_data"`
	Rewards     []RewardDocument       `bson:"rewards"`
	StartedAt   time.Time              `bson:"started_at"`
	ExpiresAt   time.Time              `bson:"expires_at"`
	CompletedAt time.Time              `bson:"completed_at,omitempty"`
	CreatedAt   time.Time              `bson:"created_at"`
	UpdatedAt   time.Time              `bson:"updated_at"`
}

// Save 保存小游戏聚合根
func (r *MinigameRepository) Save(ctx context.Context, minigameAggregate *minigame.MinigameAggregate) error {
	doc := r.toMinigameDocument(minigameAggregate)
	
	filter := bson.M{"minigame_id": doc.MinigameID}
	update := bson.M{
		"$set": doc,
		"$inc": bson.M{"version": 1},
	}
	opts := options.Update().SetUpsert(true)
	
	_, err := r.minigameColl.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save minigame: %w", err)
	}
	
	return nil
}

// FindByID 根据ID查找小游戏
func (r *MinigameRepository) FindByID(ctx context.Context, minigameID string) (*minigame.MinigameAggregate, error) {
	filter := bson.M{"minigame_id": minigameID}
	
	var doc MinigameDocument
	err := r.minigameColl.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find minigame: %w", err)
	}
	
	return r.fromMinigameDocument(&doc), nil
}

// FindByType 根据游戏类型查找小游戏
func (r *MinigameRepository) FindByType(ctx context.Context, gameType minigame.GameType) ([]*minigame.MinigameAggregate, error) {
	filter := bson.M{
		"game_type": gameType.String(),
		"is_active": true,
	}
	
	cursor, err := r.minigameColl.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find minigames by type: %w", err)
	}
	defer cursor.Close(ctx)
	
	var minigames []*minigame.MinigameAggregate
	for cursor.Next(ctx) {
		var doc MinigameDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode minigame document: %w", err)
		}
		minigames = append(minigames, r.fromMinigameDocument(&doc))
	}
	
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	
	return minigames, nil
}

// FindActive 查找激活的小游戏
func (r *MinigameRepository) FindActive(ctx context.Context) ([]*minigame.MinigameAggregate, error) {
	filter := bson.M{"is_active": true}
	
	cursor, err := r.minigameColl.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find active minigames: %w", err)
	}
	defer cursor.Close(ctx)
	
	var minigames []*minigame.MinigameAggregate
	for cursor.Next(ctx) {
		var doc MinigameDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode minigame document: %w", err)
		}
		minigames = append(minigames, r.fromMinigameDocument(&doc))
	}
	
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	
	return minigames, nil
}

// Delete 删除小游戏
func (r *MinigameRepository) Delete(ctx context.Context, minigameID string) error {
	filter := bson.M{"minigame_id": minigameID}
	update := bson.M{
		"$set": bson.M{
			"is_active": false,
			"updated_at": time.Now(),
		},
		"$inc": bson.M{"version": 1},
	}
	
	_, err := r.minigameColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete minigame: %w", err)
	}
	
	return nil
}

// GameSessionRepository 游戏会话仓储实现
type GameSessionRepository struct {
	db          *mongo.Database
	sessionColl *mongo.Collection
}

// NewGameSessionRepository 创建游戏会话仓储
func NewGameSessionRepository(db *mongo.Database) *GameSessionRepository {
	return &GameSessionRepository{
		db:          db,
		sessionColl: db.Collection("game_sessions"),
	}
}

// Save 保存游戏会话
func (r *GameSessionRepository) Save(ctx context.Context, session *minigame.GameSession) error {
	doc := r.toGameSessionDocument(session)
	
	filter := bson.M{"session_id": doc.SessionID}
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)
	
	_, err := r.sessionColl.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save game session: %w", err)
	}
	
	return nil
}

// FindByID 根据ID查找游戏会话
func (r *GameSessionRepository) FindByID(ctx context.Context, sessionID string) (*minigame.GameSession, error) {
	filter := bson.M{"session_id": sessionID}
	
	var doc GameSessionDocument
	err := r.sessionColl.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find game session: %w", err)
	}
	
	return r.fromGameSessionDocument(&doc), nil
}

// FindActiveByPlayer 根据玩家查找激活的游戏会话
func (r *GameSessionRepository) FindActiveByPlayer(ctx context.Context, playerID uint64) (*minigame.GameSession, error) {
	filter := bson.M{
		"player_id": playerID,
		"status":    minigame.SessionStatusActive.String(),
	}
	
	var doc GameSessionDocument
	err := r.sessionColl.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find active game session: %w", err)
	}
	
	return r.fromGameSessionDocument(&doc), nil
}

// FindByPlayer 根据玩家查找游戏会话
func (r *GameSessionRepository) FindByPlayer(ctx context.Context, playerID uint64, limit int) ([]*minigame.GameSession, error) {
	filter := bson.M{"player_id": playerID}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit))
	
	cursor, err := r.sessionColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find game sessions by player: %w", err)
	}
	defer cursor.Close(ctx)
	
	var sessions []*minigame.GameSession
	for cursor.Next(ctx) {
		var doc GameSessionDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode game session document: %w", err)
		}
		sessions = append(sessions, r.fromGameSessionDocument(&doc))
	}
	
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	
	return sessions, nil
}

// FindByMinigame 根据小游戏查找会话
func (r *GameSessionRepository) FindByMinigame(ctx context.Context, minigameID string, limit int) ([]*minigame.GameSession, error) {
	filter := bson.M{"minigame_id": minigameID}
	opts := options.Find().
		SetSort(bson.D{{Key: "score", Value: -1}}).
		SetLimit(int64(limit))
	
	cursor, err := r.sessionColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find game sessions by minigame: %w", err)
	}
	defer cursor.Close(ctx)
	
	var sessions []*minigame.GameSession
	for cursor.Next(ctx) {
		var doc GameSessionDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode game session document: %w", err)
		}
		sessions = append(sessions, r.fromGameSessionDocument(&doc))
	}
	
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	
	return sessions, nil
}

// FindByQuery 根据查询条件查找会话
func (r *GameSessionRepository) FindByQuery(ctx context.Context, query *minigame.GameSessionQuery) ([]*minigame.GameSession, int64, error) {
	filter := r.buildGameSessionFilter(query)
	
	// 计算总数
	total, err := r.sessionColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count game sessions: %w", err)
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
	
	cursor, err := r.sessionColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find game sessions: %w", err)
	}
	defer cursor.Close(ctx)
	
	var sessions []*minigame.GameSession
	for cursor.Next(ctx) {
		var doc GameSessionDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, fmt.Errorf("failed to decode game session document: %w", err)
		}
		sessions = append(sessions, r.fromGameSessionDocument(&doc))
	}
	
	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor error: %w", err)
	}
	
	return sessions, total, nil
}

// CleanupExpiredSessions 清理过期会话
func (r *GameSessionRepository) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	filter := bson.M{
		"status": minigame.SessionStatusActive.String(),
		"expires_at": bson.M{"$lt": time.Now()},
	}
	update := bson.M{
		"$set": bson.M{
			"status": minigame.SessionStatusExpired.String(),
			"updated_at": time.Now(),
		},
	}
	
	result, err := r.sessionColl.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}
	
	return result.ModifiedCount, nil
}

// 私有方法

// toMinigameDocument 转换为小游戏文档
func (r *MinigameRepository) toMinigameDocument(minigameAggregate *minigame.MinigameAggregate) *MinigameDocument {
	// 转换奖励
	rewards := make([]RewardDocument, 0)
	for _, reward := range minigameAggregate.GetRewards() {
		rewards = append(rewards, RewardDocument{
			Type:     reward.GetType().String(),
			ItemID:   reward.GetItemID(),
			Quantity: reward.GetQuantity(),
			Reason:   reward.GetReason(),
		})
	}
	
	// 转换统计信息
	stats := minigameAggregate.GetStatistics()
	statistics := StatisticsDocument{
		TotalPlays:     stats.GetTotalPlays(),
		TotalPlayers:   stats.GetTotalPlayers(),
		AverageScore:   stats.GetAverageScore(),
		HighestScore:   stats.GetHighestScore(),
		AverageTime:    stats.GetAverageTime(),
		CompletionRate: stats.GetCompletionRate(),
	}
	
	return &MinigameDocument{
		MinigameID:  minigameAggregate.GetID(),
		Name:        minigameAggregate.GetName(),
		Description: minigameAggregate.GetDescription(),
		GameType:    minigameAggregate.GetGameType().String(),
		Difficulty:  minigameAggregate.GetDifficulty().String(),
		MaxPlayers:  minigameAggregate.GetMaxPlayers(),
		TimeLimit:   minigameAggregate.GetTimeLimit(),
		IsActive:    minigameAggregate.IsActive(),
		Rules:       minigameAggregate.GetRules(),
		Rewards:     rewards,
		Settings:    minigameAggregate.GetSettings(),
		Statistics:  statistics,
		CreatedAt:   minigameAggregate.GetCreatedAt(),
		UpdatedAt:   minigameAggregate.GetUpdatedAt(),
		Version:     minigameAggregate.GetVersion(),
	}
}

// fromMinigameDocument 从小游戏文档转换
func (r *MinigameRepository) fromMinigameDocument(doc *MinigameDocument) *minigame.MinigameAggregate {
	// 解析枚举值
	gameType := minigame.ParseGameType(doc.GameType)
	difficulty := minigame.ParseDifficulty(doc.Difficulty)
	
	// 重建聚合根
	minigameAggregate := minigame.NewMinigameAggregate(doc.Name, gameType, difficulty)
	minigameAggregate.SetID(doc.MinigameID)
	minigameAggregate.SetDescription(doc.Description)
	minigameAggregate.SetMaxPlayers(doc.MaxPlayers)
	minigameAggregate.SetTimeLimit(doc.TimeLimit)
	minigameAggregate.SetRules(doc.Rules)
	minigameAggregate.SetSettings(doc.Settings)
	minigameAggregate.SetVersion(doc.Version)
	
	// 转换奖励
	for _, rewardDoc := range doc.Rewards {
		rewardType := minigame.ParseRewardType(rewardDoc.Type)
		reward := minigame.NewGameReward(rewardType, rewardDoc.ItemID, rewardDoc.Quantity, rewardDoc.Reason)
		minigameAggregate.AddReward(reward)
	}
	
	// 设置统计信息
	stats := minigame.NewGameStatistics()
	stats.SetTotalPlays(doc.Statistics.TotalPlays)
	stats.SetTotalPlayers(doc.Statistics.TotalPlayers)
	stats.SetAverageScore(doc.Statistics.AverageScore)
	stats.SetHighestScore(doc.Statistics.HighestScore)
	stats.SetAverageTime(doc.Statistics.AverageTime)
	stats.SetCompletionRate(doc.Statistics.CompletionRate)
	minigameAggregate.SetStatistics(stats)
	
	if doc.IsActive {
		minigameAggregate.Activate()
	} else {
		minigameAggregate.Deactivate()
	}
	
	return minigameAggregate
}

// toGameSessionDocument 转换为游戏会话文档
func (r *GameSessionRepository) toGameSessionDocument(session *minigame.GameSession) *GameSessionDocument {
	// 转换奖励
	rewards := make([]RewardDocument, 0)
	for _, reward := range session.GetRewards() {
		rewards = append(rewards, RewardDocument{
			Type:     reward.GetType().String(),
			ItemID:   reward.GetItemID(),
			Quantity: reward.GetQuantity(),
			Reason:   reward.GetReason(),
		})
	}
	
	doc := &GameSessionDocument{
		SessionID:   session.GetID(),
		MinigameID:  session.GetMinigameID(),
		PlayerID:    session.GetPlayerID(),
		Status:      session.GetStatus().String(),
		Score:       session.GetScore(),
		TimeLimit:   session.GetTimeLimit(),
		TimeElapsed: session.GetTimeElapsed(),
		Settings:    session.GetSettings(),
		GameData:    session.GetGameData(),
		Rewards:     rewards,
		StartedAt:   session.GetStartedAt(),
		ExpiresAt:   session.GetExpiresAt(),
		CreatedAt:   session.GetCreatedAt(),
		UpdatedAt:   session.GetUpdatedAt(),
	}
	
	if !session.GetCompletedAt().IsZero() {
		doc.CompletedAt = session.GetCompletedAt()
	}
	
	return doc
}

// fromGameSessionDocument 从游戏会话文档转换
func (r *GameSessionRepository) fromGameSessionDocument(doc *GameSessionDocument) *minigame.GameSession {
	// 解析状态
	status := minigame.ParseSessionStatus(doc.Status)
	
	// 重建会话
	session := minigame.NewGameSession(
		doc.SessionID,
		doc.MinigameID,
		doc.PlayerID,
		doc.TimeLimit,
	)
	
	session.SetStatus(status)
	session.SetScore(doc.Score)
	session.SetTimeElapsed(doc.TimeElapsed)
	session.SetSettings(doc.Settings)
	session.SetGameData(doc.GameData)
	session.SetTimestamps(doc.StartedAt, doc.ExpiresAt, doc.CompletedAt)
	
	// 转换奖励
	for _, rewardDoc := range doc.Rewards {
		rewardType := minigame.ParseRewardType(rewardDoc.Type)
		reward := minigame.NewGameReward(rewardType, rewardDoc.ItemID, rewardDoc.Quantity, rewardDoc.Reason)
		session.AddReward(reward)
	}
	
	return session
}

// buildGameSessionFilter 构建游戏会话查询过滤器
func (r *GameSessionRepository) buildGameSessionFilter(query *minigame.GameSessionQuery) bson.M {
	filter := bson.M{}
	
	if query.GetPlayerID() > 0 {
		filter["player_id"] = query.GetPlayerID()
	}
	
	if query.GetMinigameID() != "" {
		filter["minigame_id"] = query.GetMinigameID()
	}
	
	if query.GetStatus() != nil {
		filter["status"] = query.GetStatus().String()
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
	
	if !query.GetStartTime().IsZero() {
		filter["started_at"] = bson.M{"$gte": query.GetStartTime()}
	}
	
	if !query.GetEndTime().IsZero() {
		if timeFilter, exists := filter["started_at"]; exists {
			timeFilter.(bson.M)["$lte"] = query.GetEndTime()
		} else {
			filter["started_at"] = bson.M{"$lte": query.GetEndTime()}
		}
	}
	
	return filter
}

// CreateIndexes 创建索引
func (r *MinigameRepository) CreateIndexes(ctx context.Context) error {
	// 小游戏索引
	minigameIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "minigame_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "game_type", Value: 1}, {Key: "is_active", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "difficulty", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_active", Value: 1}},
		},
	}
	
	if _, err := r.minigameColl.Indexes().CreateMany(ctx, minigameIndexes); err != nil {
		return fmt.Errorf("failed to create minigame indexes: %w", err)
	}
	
	return nil
}

// CreateIndexes 创建游戏会话索引
func (r *GameSessionRepository) CreateIndexes(ctx context.Context) error {
	// 游戏会话索引
	sessionIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "session_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "player_id", Value: 1}, {Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "minigame_id", Value: 1}, {Key: "score", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}, {Key: "expires_at", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "started_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "completed_at", Value: -1}},
		},
	}
	
	if _, err := r.sessionColl.Indexes().CreateMany(ctx, sessionIndexes); err != nil {
		return fmt.Errorf("failed to create game session indexes: %w", err)
	}
	
	return nil
}