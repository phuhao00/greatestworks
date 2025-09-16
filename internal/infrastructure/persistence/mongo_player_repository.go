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
)

// MongoPlayerRepo MongoDB玩家仓储实现
type MongoPlayerRepo struct {
	collection *mongo.Collection
}

// NewMongoPlayerRepository 创建MongoDB玩家仓储
func NewMongoPlayerRepository(db *mongo.Database) player.Repository {
	return &MongoPlayerRepo{
		collection: db.Collection("players"),
	}
}

// PlayerDoc MongoDB玩家文档
type PlayerDoc struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	PlayerID  string             `bson:"player_id"`
	Name      string             `bson:"name"`
	Level     int                `bson:"level"`
	Exp       int64              `bson:"exp"`
	Status    string             `bson:"status"`
	Position  PositionDoc        `bson:"position"`
	Stats     StatsDoc           `bson:"stats"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Version   int64              `bson:"version"`
}

// PositionDoc 位置文档
type PositionDoc struct {
	X float64 `bson:"x"`
	Y float64 `bson:"y"`
	Z float64 `bson:"z"`
}

// StatsDoc 属性文档
type StatsDoc struct {
	HP      int `bson:"hp"`
	MaxHP   int `bson:"max_hp"`
	MP      int `bson:"mp"`
	MaxMP   int `bson:"max_mp"`
	Attack  int `bson:"attack"`
	Defense int `bson:"defense"`
	Speed   int `bson:"speed"`
}

// Save 保存玩家
func (r *MongoPlayerRepo) Save(ctx context.Context, p *player.Player) error {
	doc := r.playerToDoc(p)
	
	filter := bson.M{"player_id": p.ID().String()}
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)
	
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save player: %w", err)
	}
	
	return nil
}

// FindByID 根据ID查找玩家
func (r *MongoPlayerRepo) FindByID(ctx context.Context, id player.PlayerID) (*player.Player, error) {
	filter := bson.M{"player_id": id.String()}
	
	var doc PlayerDoc
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, player.ErrPlayerNotFound
		}
		return nil, fmt.Errorf("failed to find player: %w", err)
	}
	
	return r.docToPlayer(&doc), nil
}

// FindByName 根据名称查找玩家
func (r *MongoPlayerRepo) FindByName(ctx context.Context, name string) (*player.Player, error) {
	filter := bson.M{"name": name}
	
	var doc PlayerDoc
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, player.ErrPlayerNotFound
		}
		return nil, fmt.Errorf("failed to find player by name: %w", err)
	}
	
	return r.docToPlayer(&doc), nil
}

// Update 更新玩家
func (r *MongoPlayerRepo) Update(ctx context.Context, p *player.Player) error {
	return r.Save(ctx, p)
}

// Delete 删除玩家
func (r *MongoPlayerRepo) Delete(ctx context.Context, id player.PlayerID) error {
	filter := bson.M{"player_id": id.String()}
	
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete player: %w", err)
	}
	
	if result.DeletedCount == 0 {
		return player.ErrPlayerNotFound
	}
	
	return nil
}

// FindOnlinePlayers 查找在线玩家
func (r *MongoPlayerRepo) FindOnlinePlayers(ctx context.Context, limit int) ([]*player.Player, error) {
	filter := bson.M{"status": "online"}
	opts := options.Find().SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find online players: %w", err)
	}
	defer cursor.Close(ctx)
	
	var docs []PlayerDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode online players: %w", err)
	}
	
	players := make([]*player.Player, len(docs))
	for i, doc := range docs {
		players[i] = r.docToPlayer(&doc)
	}
	
	return players, nil
}

// FindPlayersByLevel 根据等级范围查找玩家
func (r *MongoPlayerRepo) FindPlayersByLevel(ctx context.Context, minLevel, maxLevel int) ([]*player.Player, error) {
	filter := bson.M{
		"level": bson.M{
			"$gte": minLevel,
			"$lte": maxLevel,
		},
	}
	
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find players by level: %w", err)
	}
	defer cursor.Close(ctx)
	
	var docs []PlayerDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode players by level: %w", err)
	}
	
	players := make([]*player.Player, len(docs))
	for i, doc := range docs {
		players[i] = r.docToPlayer(&doc)
	}
	
	return players, nil
}

// ExistsByName 检查名称是否存在
func (r *MongoPlayerRepo) ExistsByName(ctx context.Context, name string) (bool, error) {
	filter := bson.M{"name": name}
	
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check player name exists: %w", err)
	}
	
	return count > 0, nil
}

// playerToDoc 将玩家聚合根转换为文档
func (r *MongoPlayerRepo) playerToDoc(p *player.Player) *PlayerDoc {
	position := p.GetPosition()
	stats := p.Stats()
	
	return &PlayerDoc{
		PlayerID: p.ID().String(),
		Name:     p.Name(),
		Level:    p.Level(),
		Exp:      p.Exp(),
		Status:   r.statusToString(p.Status()),
		Position: PositionDoc{
			X: position.X,
			Y: position.Y,
			Z: position.Z,
		},
		Stats: StatsDoc{
			HP:      stats.HP,
			MaxHP:   stats.MaxHP,
			MP:      stats.MP,
			MaxMP:   stats.MaxMP,
			Attack:  stats.Attack,
			Defense: stats.Defense,
			Speed:   stats.Speed,
		},
		CreatedAt: p.CreatedAt(),
		UpdatedAt: time.Now(),
		Version:   p.Version(),
	}
}

// docToPlayer 将文档转换为玩家聚合根
func (r *MongoPlayerRepo) docToPlayer(doc *PlayerDoc) *player.Player {
	// 重建PlayerID
	playerID := player.PlayerIDFromString(doc.PlayerID)
	
	// 重建Position
	position := player.Position{
		X: doc.Position.X,
		Y: doc.Position.Y,
		Z: doc.Position.Z,
	}
	
	// 重建Stats
	stats := player.PlayerStats{
		HP:      doc.Stats.HP,
		MaxHP:   doc.Stats.MaxHP,
		MP:      doc.Stats.MP,
		MaxMP:   doc.Stats.MaxMP,
		Attack:  doc.Stats.Attack,
		Defense: doc.Stats.Defense,
		Speed:   doc.Stats.Speed,
	}
	
	// 使用重建方法
	return player.ReconstructPlayer(
		playerID,
		doc.Name,
		doc.Level,
		doc.Exp,
		r.stringToStatus(doc.Status),
		position,
		stats,
		doc.CreatedAt,
		doc.UpdatedAt,
		doc.Version,
	)
}

// statusToString 状态转字符串
func (r *MongoPlayerRepo) statusToString(status player.PlayerStatus) string {
	switch status {
	case player.PlayerStatusOffline:
		return "offline"
	case player.PlayerStatusOnline:
		return "online"
	case player.PlayerStatusInBattle:
		return "in_battle"
	case player.PlayerStatusInScene:
		return "in_scene"
	default:
		return "offline"
	}
}

// stringToStatus 字符串转状态
func (r *MongoPlayerRepo) stringToStatus(status string) player.PlayerStatus {
	switch status {
	case "online":
		return player.PlayerStatusOnline
	case "in_battle":
		return player.PlayerStatusInBattle
	case "in_scene":
		return player.PlayerStatusInScene
	default:
		return player.PlayerStatusOffline
	}
}

// CreateIndexes 创建索引
func (r *MongoPlayerRepo) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "player_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "level", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}
	
	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create player indexes: %w", err)
	}
	
	return nil
}