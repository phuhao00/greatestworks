package persistence

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"greatestworks/internal/domain/battle"
	"greatestworks/internal/domain/player"
)

// MongoBattleRepo MongoDB战斗仓储实现
type MongoBattleRepo struct {
	collection *mongo.Collection
}

// NewMongoBattleRepository 创建MongoDB战斗仓储
func NewMongoBattleRepository(db *mongo.Database) battle.Repository {
	return &MongoBattleRepo{
		collection: db.Collection("battles"),
	}
}

// BattleDoc MongoDB战斗文档
type BattleDoc struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty"`
	BattleID     string                 `bson:"battle_id"`
	BattleType   int                    `bson:"battle_type"`
	Status       int                    `bson:"status"`
	Participants []BattleParticipantDoc `bson:"participants"`
	Rounds       []BattleRoundDoc       `bson:"rounds"`
	Winner       *string                `bson:"winner,omitempty"`
	StartTime    time.Time              `bson:"start_time"`
	EndTime      *time.Time             `bson:"end_time,omitempty"`
	CreatedAt    time.Time              `bson:"created_at"`
	UpdatedAt    time.Time              `bson:"updated_at"`
	Version      int64                  `bson:"version"`
}

// BattleParticipantDoc 战斗参与者文档
type BattleParticipantDoc struct {
	PlayerID    string    `bson:"player_id"`
	Team        int       `bson:"team"`
	CurrentHP   int       `bson:"current_hp"`
	CurrentMP   int       `bson:"current_mp"`
	IsAlive     bool      `bson:"is_alive"`
	DamageDealt int       `bson:"damage_dealt"`
	DamageTaken int       `bson:"damage_taken"`
	JoinedAt    time.Time `bson:"joined_at"`
}

// BattleRoundDoc 战斗回合文档
type BattleRoundDoc struct {
	RoundNumber int               `bson:"round_number"`
	Actions     []BattleActionDoc `bson:"actions"`
	StartTime   time.Time         `bson:"start_time"`
	EndTime     *time.Time        `bson:"end_time,omitempty"`
}

// BattleActionDoc 战斗行动文档
type BattleActionDoc struct {
	ActionID   string    `bson:"action_id"`
	ActorID    string    `bson:"actor_id"`
	TargetID   *string   `bson:"target_id,omitempty"`
	ActionType int       `bson:"action_type"`
	SkillID    *string   `bson:"skill_id,omitempty"`
	Damage     int       `bson:"damage"`
	Healing    int       `bson:"healing"`
	Critical   bool      `bson:"critical"`
	Timestamp  time.Time `bson:"timestamp"`
}

// Save 保存战斗
func (r *MongoBattleRepo) Save(ctx context.Context, b *battle.Battle) error {
	doc := r.battleToDoc(b)

	filter := bson.M{"battle_id": b.ID().String()}
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save battle: %w", err)
	}

	return nil
}

// FindByID 根据ID查找战斗
func (r *MongoBattleRepo) FindByID(ctx context.Context, id battle.BattleID) (*battle.Battle, error) {
	filter := bson.M{"battle_id": id.String()}

	var doc BattleDoc
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, battle.ErrBattleNotFound
		}
		return nil, fmt.Errorf("failed to find battle: %w", err)
	}

	return r.docToBattle(&doc), nil
}

// Update 更新战斗
func (r *MongoBattleRepo) Update(ctx context.Context, b *battle.Battle) error {
	return r.Save(ctx, b)
}

// Delete 删除战斗
func (r *MongoBattleRepo) Delete(ctx context.Context, id battle.BattleID) error {
	filter := bson.M{"battle_id": id.String()}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete battle: %w", err)
	}

	if result.DeletedCount == 0 {
		return battle.ErrBattleNotFound
	}

	return nil
}

// FindByPlayerID 根据玩家ID查找战斗
func (r *MongoBattleRepo) FindByPlayerID(ctx context.Context, playerID player.PlayerID, limit int) ([]*battle.Battle, error) {
	filter := bson.M{
		"participants.player_id": playerID.String(),
	}
	opts := options.Find().SetLimit(int64(limit)).SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find battles by player: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []BattleDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode battles: %w", err)
	}

	battles := make([]*battle.Battle, len(docs))
	for i, doc := range docs {
		battles[i] = r.docToBattle(&doc)
	}

	return battles, nil
}

// FindActiveBattles 查找进行中的战斗
func (r *MongoBattleRepo) FindActiveBattles(ctx context.Context, limit int) ([]*battle.Battle, error) {
	filter := bson.M{
		"status": int(battle.BattleStatusInProgress),
	}
	opts := options.Find().SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find active battles: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []BattleDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode active battles: %w", err)
	}

	battles := make([]*battle.Battle, len(docs))
	for i, doc := range docs {
		battles[i] = r.docToBattle(&doc)
	}

	return battles, nil
}

// FindByStatus 根据状态查找战斗
func (r *MongoBattleRepo) FindByStatus(ctx context.Context, status battle.BattleStatus, limit int) ([]*battle.Battle, error) {
	filter := bson.M{
		"status": int(status),
	}
	opts := options.Find().SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find battles by status: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []BattleDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode battles by status: %w", err)
	}

	battles := make([]*battle.Battle, len(docs))
	for i, doc := range docs {
		battles[i] = r.docToBattle(&doc)
	}

	return battles, nil
}

// FindByType 根据类型查找战斗
func (r *MongoBattleRepo) FindByType(ctx context.Context, battleType battle.BattleType, limit int) ([]*battle.Battle, error) {
	filter := bson.M{
		"battle_type": int(battleType),
	}
	opts := options.Find().SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find battles by type: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []BattleDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode battles by type: %w", err)
	}

	battles := make([]*battle.Battle, len(docs))
	for i, doc := range docs {
		battles[i] = r.docToBattle(&doc)
	}

	return battles, nil
}

// CountByPlayerID 统计玩家参与的战斗数量
func (r *MongoBattleRepo) CountByPlayerID(ctx context.Context, playerID player.PlayerID) (int64, error) {
	filter := bson.M{
		"participants.player_id": playerID.String(),
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count battles by player: %w", err)
	}

	return count, nil
}

// battleToDoc 将战斗聚合根转换为文档
func (r *MongoBattleRepo) battleToDoc(b *battle.Battle) *BattleDoc {
	// 转换参与者
	participants := make([]BattleParticipantDoc, len(b.Participants()))
	for i, p := range b.Participants() {
		participants[i] = BattleParticipantDoc{
			PlayerID:    p.PlayerID.String(),
			Team:        p.Team,
			CurrentHP:   p.CurrentHP,
			CurrentMP:   p.CurrentMP,
			IsAlive:     p.IsAlive,
			DamageDealt: p.DamageDealt,
			DamageTaken: p.DamageTaken,
			JoinedAt:    p.JoinedAt,
		}
	}

	// 转换回合（这里需要Battle提供获取回合的方法）
	rounds := make([]BattleRoundDoc, 0)
	// TODO: 实现回合转换

	doc := &BattleDoc{
		BattleID:     b.ID().String(),
		BattleType:   int(b.GetBattleType()),
		Status:       int(b.Status()),
		Participants: participants,
		Rounds:       rounds,
		StartTime:    b.StartTime(),
		CreatedAt:    b.CreatedAt(),
		UpdatedAt:    time.Now(),
		Version:      b.Version(),
	}

	if b.Winner() != nil {
		winnerStr := b.Winner().String()
		doc.Winner = &winnerStr
	}

	if b.EndTime() != nil {
		doc.EndTime = b.EndTime()
	}

	return doc
}

// docToBattle 将文档转换为战斗聚合根
func (r *MongoBattleRepo) docToBattle(doc *BattleDoc) *battle.Battle {
	// 这里需要Battle提供重建方法
	// 暂时返回新创建的战斗
	b := battle.NewBattle(battle.BattleType(doc.BattleType))

	// TODO: 设置其他属性

	return b
}

// CreateIndexes 创建索引
func (r *MongoBattleRepo) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "battle_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "battle_type", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "participants.player_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "start_time", Value: -1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create battle indexes: %w", err)
	}

	return nil
}
