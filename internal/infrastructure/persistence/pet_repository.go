package persistence

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"greatestworks/internal/domain/pet"
)

// PetRepository MongoDB宠物仓储实现
type PetRepository struct {
	db       *mongo.Database
	petColl  *mongo.Collection
	fragColl *mongo.Collection
	skinColl *mongo.Collection
	bondColl *mongo.Collection
	pictColl *mongo.Collection
}

// NewPetRepository 创建宠物仓储
func NewPetRepository(db *mongo.Database) *PetRepository {
	return &PetRepository{
		db:       db,
		petColl:  db.Collection("pets"),
		fragColl: db.Collection("pet_fragments"),
		skinColl: db.Collection("pet_skins"),
		bondColl: db.Collection("pet_bonds"),
		pictColl: db.Collection("pet_pictorials"),
	}
}

// PetDocument 宠物文档结构
type PetDocument struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	PetID        string             `bson:"pet_id"`
	PlayerID     uint64             `bson:"player_id"`
	SpeciesID    string             `bson:"species_id"`
	Name         string             `bson:"name"`
	Level        int32              `bson:"level"`
	Exp          int64              `bson:"exp"`
	MaxExp       int64              `bson:"max_exp"`
	Rarity       string             `bson:"rarity"`
	Quality      string             `bson:"quality"`
	Attributes   map[string]int64   `bson:"attributes"`
	Skills       []string           `bson:"skills"`
	EquippedSkin string             `bson:"equipped_skin"`
	Mood         string             `bson:"mood"`
	Hunger       int32              `bson:"hunger"`
	Energy       int32              `bson:"energy"`
	Health       int32              `bson:"health"`
	Happiness    int32              `bson:"happiness"`
	IsActive     bool               `bson:"is_active"`
	LastFedAt    time.Time          `bson:"last_fed_at"`
	LastPlayedAt time.Time          `bson:"last_played_at"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	Version      int64              `bson:"version"`
}

// PetFragmentDocument 宠物碎片文档结构
type PetFragmentDocument struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	FragmentID string             `bson:"fragment_id"`
	PlayerID   uint64             `bson:"player_id"`
	SpeciesID  string             `bson:"species_id"`
	Quantity   int32              `bson:"quantity"`
	Required   int32              `bson:"required"`
	Source     string             `bson:"source"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

// PetSkinDocument 宠物皮肤文档结构
type PetSkinDocument struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	SkinID     string             `bson:"skin_id"`
	PlayerID   uint64             `bson:"player_id"`
	SpeciesID  string             `bson:"species_id"`
	Name       string             `bson:"name"`
	Rarity     string             `bson:"rarity"`
	Effects    map[string]float64 `bson:"effects"`
	IsUnlocked bool               `bson:"is_unlocked"`
	UnlockedAt time.Time          `bson:"unlocked_at"`
	CreatedAt  time.Time          `bson:"created_at"`
}

// Save 保存宠物聚合根
func (r *PetRepository) Save(ctx context.Context, petAggregate *pet.PetAggregate) error {
	doc := r.toPetDocument(petAggregate)

	filter := bson.M{"pet_id": doc.PetID}
	update := bson.M{
		"$set": doc,
		"$inc": bson.M{"version": 1},
	}
	opts := options.Update().SetUpsert(true)

	_, err := r.petColl.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save pet: %w", err)
	}

	return nil
}

// FindByID 根据ID查找宠物
func (r *PetRepository) FindByID(ctx context.Context, petID string) (*pet.PetAggregate, error) {
	filter := bson.M{"pet_id": petID}

	var doc PetDocument
	err := r.petColl.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find pet: %w", err)
	}

	return r.fromPetDocument(&doc), nil
}

// FindByPlayer 根据玩家ID查找宠物列表
func (r *PetRepository) FindByPlayer(ctx context.Context, playerID uint64) ([]*pet.PetAggregate, error) {
	filter := bson.M{"player_id": playerID, "is_active": true}

	cursor, err := r.petColl.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find pets by player: %w", err)
	}
	defer cursor.Close(ctx)

	var pets []*pet.PetAggregate
	for cursor.Next(ctx) {
		var doc PetDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode pet document: %w", err)
		}
		pets = append(pets, r.fromPetDocument(&doc))
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return pets, nil
}

// FindByQuery 根据查询条件查找宠物
func (r *PetRepository) FindByQuery(ctx context.Context, query *pet.PetQuery) ([]*pet.PetAggregate, int64, error) {
	filter := r.buildPetFilter(query)

	// 计算总数
	total, err := r.petColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pets: %w", err)
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

	cursor, err := r.petColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find pets: %w", err)
	}
	defer cursor.Close(ctx)

	var pets []*pet.PetAggregate
	for cursor.Next(ctx) {
		var doc PetDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, fmt.Errorf("failed to decode pet document: %w", err)
		}
		pets = append(pets, r.fromPetDocument(&doc))
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor error: %w", err)
	}

	return pets, total, nil
}

// Delete 删除宠物
func (r *PetRepository) Delete(ctx context.Context, petID string) error {
	filter := bson.M{"pet_id": petID}
	update := bson.M{
		"$set": bson.M{
			"is_active":  false,
			"updated_at": time.Now(),
		},
		"$inc": bson.M{"version": 1},
	}

	_, err := r.petColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete pet: %w", err)
	}

	return nil
}

// SaveFragment 保存宠物碎片
func (r *PetRepository) SaveFragment(ctx context.Context, fragment *pet.PetFragment) error {
	doc := r.toFragmentDocument(fragment)

	filter := bson.M{"fragment_id": doc.FragmentID}
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)

	_, err := r.fragColl.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save pet fragment: %w", err)
	}

	return nil
}

// FindFragmentsByPlayer 根据玩家ID查找宠物碎片
func (r *PetRepository) FindFragmentsByPlayer(ctx context.Context, playerID uint64) ([]*pet.PetFragment, error) {
	filter := bson.M{"player_id": playerID}

	cursor, err := r.fragColl.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find pet fragments: %w", err)
	}
	defer cursor.Close(ctx)

	var fragments []*pet.PetFragment
	for cursor.Next(ctx) {
		var doc PetFragmentDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode fragment document: %w", err)
		}
		fragments = append(fragments, r.fromFragmentDocument(&doc))
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return fragments, nil
}

// SaveSkin 保存宠物皮肤
func (r *PetRepository) SaveSkin(ctx context.Context, skin *pet.PetSkin) error {
	doc := r.toSkinDocument(skin)

	filter := bson.M{"skin_id": doc.SkinID}
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)

	_, err := r.skinColl.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save pet skin: %w", err)
	}

	return nil
}

// FindSkinsByPlayer 根据玩家ID查找宠物皮肤
func (r *PetRepository) FindSkinsByPlayer(ctx context.Context, playerID uint64) ([]*pet.PetSkin, error) {
	filter := bson.M{"player_id": playerID}

	cursor, err := r.skinColl.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find pet skins: %w", err)
	}
	defer cursor.Close(ctx)

	var skins []*pet.PetSkin
	for cursor.Next(ctx) {
		var doc PetSkinDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode skin document: %w", err)
		}
		skins = append(skins, r.fromSkinDocument(&doc))
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return skins, nil
}

// 私有方法

// toPetDocument 转换为宠物文档
func (r *PetRepository) toPetDocument(petAggregate *pet.PetAggregate) *PetDocument {
	return &PetDocument{
		PetID:        petAggregate.GetID(),
		PlayerID:     petAggregate.GetPlayerID(),
		SpeciesID:    petAggregate.GetSpeciesID(),
		Name:         petAggregate.GetName(),
		Level:        petAggregate.GetLevel(),
		Exp:          petAggregate.GetExp(),
		MaxExp:       petAggregate.GetMaxExp(),
		Rarity:       petAggregate.GetRarity().String(),
		Quality:      petAggregate.GetQuality().String(),
		Attributes:   petAggregate.GetAttributes().ToMap(),
		Skills:       petAggregate.GetSkills(),
		EquippedSkin: petAggregate.GetEquippedSkin(),
		Mood:         petAggregate.GetMood().String(),
		Hunger:       petAggregate.GetHunger(),
		Energy:       petAggregate.GetEnergy(),
		Health:       petAggregate.GetHealth(),
		Happiness:    petAggregate.GetHappiness(),
		IsActive:     petAggregate.IsActive(),
		LastFedAt:    petAggregate.GetLastFedAt(),
		LastPlayedAt: petAggregate.GetLastPlayedAt(),
		CreatedAt:    petAggregate.GetCreatedAt(),
		UpdatedAt:    petAggregate.GetUpdatedAt(),
		Version:      petAggregate.GetVersion(),
	}
}

// fromPetDocument 从宠物文档转换
func (r *PetRepository) fromPetDocument(doc *PetDocument) *pet.PetAggregate {
	// 解析枚举值
	rarity := pet.ParseRarity(doc.Rarity)
	quality := pet.ParseQuality(doc.Quality)
	mood := pet.ParseMood(doc.Mood)

	// 创建属性对象
	attributes := pet.NewPetAttributes()
	for key, value := range doc.Attributes {
		attributes.Set(key, value)
	}

	// 重建聚合根
	petAggregate := pet.NewPetAggregate(
		doc.PetID,
		doc.PlayerID,
		doc.SpeciesID,
		doc.Name,
		rarity,
	)

	// 设置其他属性
	petAggregate.SetLevel(doc.Level)
	petAggregate.SetExp(doc.Exp, doc.MaxExp)
	petAggregate.SetQuality(quality)
	petAggregate.SetAttributes(attributes)
	petAggregate.SetSkills(doc.Skills)
	petAggregate.SetEquippedSkin(doc.EquippedSkin)
	petAggregate.SetMood(mood)
	petAggregate.SetStats(doc.Hunger, doc.Energy, doc.Health, doc.Happiness)
	petAggregate.SetTimestamps(doc.LastFedAt, doc.LastPlayedAt)
	petAggregate.SetVersion(doc.Version)

	if !doc.IsActive {
		petAggregate.Deactivate()
	}

	return petAggregate
}

// toFragmentDocument 转换为碎片文档
func (r *PetRepository) toFragmentDocument(fragment *pet.PetFragment) *PetFragmentDocument {
	return &PetFragmentDocument{
		FragmentID: fragment.GetID(),
		PlayerID:   fragment.GetPlayerID(),
		SpeciesID:  fragment.GetSpeciesID(),
		Quantity:   fragment.GetQuantity(),
		Required:   fragment.GetRequired(),
		Source:     fragment.GetSource(),
		CreatedAt:  fragment.GetCreatedAt(),
		UpdatedAt:  fragment.GetUpdatedAt(),
	}
}

// fromFragmentDocument 从碎片文档转换
func (r *PetRepository) fromFragmentDocument(doc *PetFragmentDocument) *pet.PetFragment {
	return pet.NewPetFragment(
		doc.FragmentID,
		doc.PlayerID,
		doc.SpeciesID,
		doc.Quantity,
		doc.Required,
		doc.Source,
	)
}

// toSkinDocument 转换为皮肤文档
func (r *PetRepository) toSkinDocument(skin *pet.PetSkin) *PetSkinDocument {
	return &PetSkinDocument{
		SkinID:     skin.GetID(),
		PlayerID:   skin.GetPlayerID(),
		SpeciesID:  skin.GetSpeciesID(),
		Name:       skin.GetName(),
		Rarity:     skin.GetRarity().String(),
		Effects:    skin.GetEffects(),
		IsUnlocked: skin.IsUnlocked(),
		UnlockedAt: skin.GetUnlockedAt(),
		CreatedAt:  skin.GetCreatedAt(),
	}
}

// fromSkinDocument 从皮肤文档转换
func (r *PetRepository) fromSkinDocument(doc *PetSkinDocument) *pet.PetSkin {
	rarity := pet.ParseRarity(doc.Rarity)

	skin := pet.NewPetSkin(
		doc.SkinID,
		doc.PlayerID,
		doc.SpeciesID,
		doc.Name,
		rarity,
		doc.Effects,
	)

	if doc.IsUnlocked {
		skin.Unlock()
	}

	return skin
}

// buildPetFilter 构建宠物查询过滤器
func (r *PetRepository) buildPetFilter(query *pet.PetQuery) bson.M {
	filter := bson.M{}

	if query.GetPlayerID() > 0 {
		filter["player_id"] = query.GetPlayerID()
	}

	if query.GetSpeciesID() != "" {
		filter["species_id"] = query.GetSpeciesID()
	}

	if query.GetRarity() != nil {
		filter["rarity"] = query.GetRarity().String()
	}

	if query.GetQuality() != nil {
		filter["quality"] = query.GetQuality().String()
	}

	if query.GetMinLevel() > 0 {
		filter["level"] = bson.M{"$gte": query.GetMinLevel()}
	}

	if query.GetMaxLevel() > 0 {
		if levelFilter, exists := filter["level"]; exists {
			levelFilter.(bson.M)["$lte"] = query.GetMaxLevel()
		} else {
			filter["level"] = bson.M{"$lte": query.GetMaxLevel()}
		}
	}

	if query.IsActiveOnly() {
		filter["is_active"] = true
	}

	return filter
}

// CreateIndexes 创建索引
func (r *PetRepository) CreateIndexes(ctx context.Context) error {
	// 宠物索引
	petIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "pet_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "player_id", Value: 1}, {Key: "is_active", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "species_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "level", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "rarity", Value: 1}},
		},
	}

	if _, err := r.petColl.Indexes().CreateMany(ctx, petIndexes); err != nil {
		return fmt.Errorf("failed to create pet indexes: %w", err)
	}

	// 碎片索引
	fragmentIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "fragment_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "player_id", Value: 1}, {Key: "species_id", Value: 1}},
		},
	}

	if _, err := r.fragColl.Indexes().CreateMany(ctx, fragmentIndexes); err != nil {
		return fmt.Errorf("failed to create fragment indexes: %w", err)
	}

	// 皮肤索引
	skinIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "skin_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "player_id", Value: 1}, {Key: "species_id", Value: 1}},
		},
	}

	if _, err := r.skinColl.Indexes().CreateMany(ctx, skinIndexes); err != nil {
		return fmt.Errorf("failed to create skin indexes: %w", err)
	}

	return nil
}
