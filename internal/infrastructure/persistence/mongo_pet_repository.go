package persistence

import (
	"context"
	"fmt"
	"time"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	"greatestworks/internal/domain/pet"
	"greatestworks/internal/infrastructure/logger"
)

// MongoPetRepository MongoDB宠物仓储实现
type MongoPetRepository struct {
	collection *mongo.Collection
	logger     logger.Logger
}

// NewMongoPetRepository 创建MongoDB宠物仓储
func NewMongoPetRepository(db *mongo.Database, logger logger.Logger) *MongoPetRepository {
	collection := db.Collection("pets")
	return &MongoPetRepository{
		collection: collection,
		logger:     logger,
	}
}

// PetDocument 宠物文档结构
type PetDocument struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	PetID      string             `bson:"pet_id"`
	PlayerID   string             `bson:"player_id"`
	ConfigID   uint32             `bson:"config_id"`
	Name       string             `bson:"name"`
	Category   int                `bson:"category"`
	Star       uint32             `bson:"star"`
	Level      uint32             `bson:"level"`
	Experience uint64             `bson:"experience"`
	State      int                `bson:"state"`
	Attributes AttributesDoc     `bson:"attributes"`
	Skills     []SkillDoc         `bson:"skills"`
	Bonds      BondsDoc           `bson:"bonds"`
	Skins      []SkinDoc          `bson:"skins"`
	ReviveTime time.Time          `bson:"revive_time"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
	Version    int                `bson:"version"`
}

// AttributesDoc 属性文档结构
type AttributesDoc struct {
	Health   int64 `bson:"health"`
	Attack   int64 `bson:"attack"`
	Defense  int64 `bson:"defense"`
	Speed    int64 `bson:"speed"`
	Critical int64 `bson:"critical"`
	Hit      int64 `bson:"hit"`
	Dodge    int64 `bson:"dodge"`
}

// SkillDoc 技能文档结构
type SkillDoc struct {
	SkillID     string    `bson:"skill_id"`
	Name        string    `bson:"name"`
	Level       uint32    `bson:"level"`
	Type        int       `bson:"type"`
	Cooldown    int64     `bson:"cooldown"`
	LastUsed    time.Time `bson:"last_used"`
	Damage      int64     `bson:"damage"`
	Description string    `bson:"description"`
}

// BondsDoc 羁绊文档结构
type BondsDoc struct {
	ActiveBonds []ActiveBondDoc `bson:"active_bonds"`
	BondPoints  int64           `bson:"bond_points"`
}

// ActiveBondDoc 激活羁绊文档结构
type ActiveBondDoc struct {
	BondID      string    `bson:"bond_id"`
	Name        string    `bson:"name"`
	Level       uint32    `bson:"level"`
	Effect      string    `bson:"effect"`
	ActivatedAt time.Time `bson:"activated_at"`
}

// SkinDoc 皮肤文档结构
type SkinDoc struct {
	SkinID         string             `bson:"skin_id"`
	Name           string             `bson:"name"`
	Rarity         int                `bson:"rarity"`
	Equipped       bool               `bson:"equipped"`
	PowerBonus     int64              `bson:"power_bonus"`
	AttributeBonus map[string]float64 `bson:"attribute_bonus"`
	Unlocked       bool               `bson:"unlocked"`
	UnlockTime     time.Time          `bson:"unlock_time"`
}

// Save 保存宠物
func (r *MongoPetRepository) Save(petAggregate *pet.PetAggregate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	doc := r.aggregateToDocument(petAggregate)
	
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		r.logger.Error("Failed to save pet", "error", err, "pet_id", petAggregate.GetID())
		return fmt.Errorf("failed to save pet: %w", err)
	}
	
	r.logger.Debug("Pet saved successfully", "pet_id", petAggregate.GetID())
	return nil
}

// FindByID 根据ID查找宠物
func (r *MongoPetRepository) FindByID(id string) (*pet.PetAggregate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	filter := bson.M{"pet_id": id}
	var doc PetDocument
	
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, pet.ErrPetNotFound
		}
		r.logger.Error("Failed to find pet by ID", "error", err, "pet_id", id)
		return nil, fmt.Errorf("failed to find pet: %w", err)
	}
	
	return r.documentToAggregate(&doc)
}

// FindByPlayer 根据玩家ID查找宠物
func (r *MongoPetRepository) FindByPlayer(playerID string) ([]*pet.PetAggregate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	filter := bson.M{"player_id": playerID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find pets by player", "error", err, "player_id", playerID)
		return nil, fmt.Errorf("failed to find pets: %w", err)
	}
	defer cursor.Close(ctx)
	
	var pets []*pet.PetAggregate
	for cursor.Next(ctx) {
		var doc PetDocument
		if err := cursor.Decode(&doc); err != nil {
			r.logger.Error("Failed to decode pet document", "error", err)
			continue
		}
		
		petAggregate, err := r.documentToAggregate(&doc)
		if err != nil {
			r.logger.Error("Failed to convert document to aggregate", "error", err)
			continue
		}
		
		pets = append(pets, petAggregate)
	}
	
	return pets, nil
}

// FindByPlayerAndCategory 根据玩家ID和类别查找宠物
func (r *MongoPetRepository) FindByPlayerAndCategory(playerID string, category pet.PetCategory) ([]*pet.PetAggregate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	filter := bson.M{
		"player_id": playerID,
		"category":  int(category),
	}
	
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find pets by player and category", "error", err, "player_id", playerID, "category", category)
		return nil, fmt.Errorf("failed to find pets: %w", err)
	}
	defer cursor.Close(ctx)
	
	var pets []*pet.PetAggregate
	for cursor.Next(ctx) {
		var doc PetDocument
		if err := cursor.Decode(&doc); err != nil {
			r.logger.Error("Failed to decode pet document", "error", err)
			continue
		}
		
		petAggregate, err := r.documentToAggregate(&doc)
		if err != nil {
			r.logger.Error("Failed to convert document to aggregate", "error", err)
			continue
		}
		
		pets = append(pets, petAggregate)
	}
	
	return pets, nil
}

// Update 更新宠物
func (r *MongoPetRepository) Update(petAggregate *pet.PetAggregate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	filter := bson.M{"pet_id": petAggregate.GetID()}
	doc := r.aggregateToDocument(petAggregate)
	
	update := bson.M{"$set": doc}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("Failed to update pet", "error", err, "pet_id", petAggregate.GetID())
		return fmt.Errorf("failed to update pet: %w", err)
	}
	
	if result.MatchedCount == 0 {
		return pet.ErrPetNotFound
	}
	
	r.logger.Debug("Pet updated successfully", "pet_id", petAggregate.GetID())
	return nil
}

// Delete 删除宠物
func (r *MongoPetRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	filter := bson.M{"pet_id": id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete pet", "error", err, "pet_id", id)
		return fmt.Errorf("failed to delete pet: %w", err)
	}
	
	if result.DeletedCount == 0 {
		return pet.ErrPetNotFound
	}
	
	r.logger.Debug("Pet deleted successfully", "pet_id", id)
	return nil
}

// FindByState 根据状态查找宠物
func (r *MongoPetRepository) FindByState(state pet.PetState) ([]*pet.PetAggregate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	filter := bson.M{"state": int(state)}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find pets by state", "error", err, "state", state)
		return nil, fmt.Errorf("failed to find pets: %w", err)
	}
	defer cursor.Close(ctx)
	
	var pets []*pet.PetAggregate
	for cursor.Next(ctx) {
		var doc PetDocument
		if err := cursor.Decode(&doc); err != nil {
			r.logger.Error("Failed to decode pet document", "error", err)
			continue
		}
		
		petAggregate, err := r.documentToAggregate(&doc)
		if err != nil {
			r.logger.Error("Failed to convert document to aggregate", "error", err)
			continue
		}
		
		pets = append(pets, petAggregate)
	}
	
	return pets, nil
}

// FindActiveByPlayer 查找玩家的活跃宠物
func (r *MongoPetRepository) FindActiveByPlayer(playerID string) ([]*pet.PetAggregate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	filter := bson.M{
		"player_id": playerID,
		"state": bson.M{"$ne": int(pet.PetStateDead)},
	}
	
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find active pets by player", "error", err, "player_id", playerID)
		return nil, fmt.Errorf("failed to find active pets: %w", err)
	}
	defer cursor.Close(ctx)
	
	var pets []*pet.PetAggregate
	for cursor.Next(ctx) {
		var doc PetDocument
		if err := cursor.Decode(&doc); err != nil {
			r.logger.Error("Failed to decode pet document", "error", err)
			continue
		}
		
		petAggregate, err := r.documentToAggregate(&doc)
		if err != nil {
			r.logger.Error("Failed to convert document to aggregate", "error", err)
			continue
		}
		
		pets = append(pets, petAggregate)
	}
	
	return pets, nil
}

// Count 统计宠物总数
func (r *MongoPetRepository) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		r.logger.Error("Failed to count pets", "error", err)
		return 0, fmt.Errorf("failed to count pets: %w", err)
	}
	
	return count, nil
}

// CountByPlayer 统计玩家宠物数量
func (r *MongoPetRepository) CountByPlayer(playerID string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	filter := bson.M{"player_id": playerID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to count pets by player", "error", err, "player_id", playerID)
		return 0, fmt.Errorf("failed to count pets: %w", err)
	}
	
	return count, nil
}

// CreateIndexes 创建索引
func (r *MongoPetRepository) CreateIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{"pet_id", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{"player_id", 1}},
		},
		{
			Keys: bson.D{{"player_id", 1}, {"category", 1}},
		},
		{
			Keys: bson.D{{"state", 1}},
		},
		{
			Keys: bson.D{{"level", 1}},
		},
		{
			Keys: bson.D{{"star", 1}},
		},
		{
			Keys: bson.D{{"created_at", -1}},
		},
	}
	
	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		r.logger.Error("Failed to create pet indexes", "error", err)
		return fmt.Errorf("failed to create indexes: %w", err)
	}
	
	r.logger.Info("Pet indexes created successfully")
	return nil
}

// 转换方法

// aggregateToDocument 将聚合根转换为文档
func (r *MongoPetRepository) aggregateToDocument(petAggregate *pet.PetAggregate) *PetDocument {
	// 转换属性
	attributes := AttributesDoc{
		Health:   petAggregate.GetAttributes().GetHealth(),
		Attack:   petAggregate.GetAttributes().GetAttack(),
		Defense:  petAggregate.GetAttributes().GetDefense(),
		Speed:    petAggregate.GetAttributes().GetSpeed(),
		Critical: petAggregate.GetAttributes().GetCritical(),
		Hit:      petAggregate.GetAttributes().GetHit(),
		Dodge:    petAggregate.GetAttributes().GetDodge(),
	}
	
	// 转换技能
	var skills []SkillDoc
	for _, skill := range petAggregate.GetSkills() {
		skills = append(skills, SkillDoc{
			SkillID:     skill.GetSkillID(),
			Name:        skill.GetName(),
			Level:       skill.GetLevel(),
			Type:        int(skill.GetType()),
			Cooldown:    skill.GetCooldown().Milliseconds(),
			LastUsed:    skill.GetLastUsed(),
			Damage:      skill.GetDamage(),
			Description: skill.GetDescription(),
		})
	}
	
	// 转换羁绊
	var activeBonds []ActiveBondDoc
	for _, bond := range petAggregate.GetBonds().GetActiveBonds() {
		activeBonds = append(activeBonds, ActiveBondDoc{
			BondID:      bond.GetBondID(),
			Name:        bond.GetName(),
			Level:       bond.GetLevel(),
			Effect:      bond.GetEffect(),
			ActivatedAt: bond.GetActivatedAt(),
		})
	}
	
	bonds := BondsDoc{
		ActiveBonds: activeBonds,
		BondPoints:  petAggregate.GetBonds().GetBondPoints(),
	}
	
	// 转换皮肤
	var skins []SkinDoc
	for _, skin := range petAggregate.GetSkins() {
		skins = append(skins, SkinDoc{
			SkinID:         skin.GetSkinID(),
			Name:           skin.GetName(),
			Rarity:         int(skin.GetRarity()),
			Equipped:       skin.IsEquipped(),
			PowerBonus:     skin.GetPowerBonus(),
			AttributeBonus: skin.GetAttributeBonus(),
			Unlocked:       skin.IsUnlocked(),
		})
	}
	
	return &PetDocument{
		PetID:      petAggregate.GetID(),
		PlayerID:   petAggregate.GetPlayerID(),
		ConfigID:   petAggregate.GetConfigID(),
		Name:       petAggregate.GetName(),
		Category:   int(petAggregate.GetCategory()),
		Star:       petAggregate.GetStar(),
		Level:      petAggregate.GetLevel(),
		Experience: petAggregate.GetExperience(),
		State:      int(petAggregate.GetState()),
		Attributes: attributes,
		Skills:     skills,
		Bonds:      bonds,
		Skins:      skins,
		ReviveTime: petAggregate.GetReviveTime(),
		CreatedAt:  petAggregate.GetCreatedAt(),
		UpdatedAt:  petAggregate.GetUpdatedAt(),
		Version:    petAggregate.GetVersion(),
	}
}

// documentToAggregate 将文档转换为聚合根
func (r *MongoPetRepository) documentToAggregate(doc *PetDocument) (*pet.PetAggregate, error) {
	// 重建属性
	attributes := pet.NewPetAttributes()
	attributes.AddHealth(doc.Attributes.Health - attributes.GetHealth())
	attributes.AddAttack(doc.Attributes.Attack - attributes.GetAttack())
	attributes.AddDefense(doc.Attributes.Defense - attributes.GetDefense())
	attributes.AddSpeed(doc.Attributes.Speed - attributes.GetSpeed())
	attributes.AddCritical(doc.Attributes.Critical - attributes.GetCritical())
	attributes.AddHit(doc.Attributes.Hit - attributes.GetHit())
	attributes.AddDodge(doc.Attributes.Dodge - attributes.GetDodge())
	
	// 重建技能
	var skills []*pet.PetSkill
	for _, skillDoc := range doc.Skills {
		skill := pet.NewPetSkill(
			skillDoc.SkillID,
			skillDoc.Name,
			pet.SkillType(skillDoc.Type),
			time.Duration(skillDoc.Cooldown)*time.Millisecond,
			skillDoc.Damage,
			skillDoc.Description,
		)
		skills = append(skills, skill)
	}
	
	// 重建羁绊
	bonds := pet.NewPetBonds()
	for _, bondDoc := range doc.Bonds.ActiveBonds {
		bond := pet.NewActiveBond(
			bondDoc.BondID,
			bondDoc.Name,
			bondDoc.Level,
			bondDoc.Effect,
		)
		bonds.AddActiveBond(bond)
	}
	bonds.AddBondPoints(doc.Bonds.BondPoints)
	
	// 重建皮肤
	var skins []*pet.PetSkin
	for _, skinDoc := range doc.Skins {
		skin := pet.NewPetSkin(
			skinDoc.SkinID,
			skinDoc.Name,
			pet.PetRarity(skinDoc.Rarity),
			skinDoc.PowerBonus,
		)
		if skinDoc.Unlocked {
			skin.Unlock()
		}
		if skinDoc.Equipped {
			skin.Equip()
		}
		skins = append(skins, skin)
	}
	
	// 使用重建函数创建聚合根
	return pet.ReconstructPetAggregate(
		doc.PetID,
		doc.PlayerID,
		doc.ConfigID,
		doc.Name,
		pet.PetCategory(doc.Category),
		doc.Star,
		doc.Level,
		doc.Experience,
		pet.PetState(doc.State),
		attributes,
		skills,
		bonds,
		skins,
		doc.ReviveTime,
		doc.CreatedAt,
		doc.UpdatedAt,
		doc.Version,
	), nil
}

// 实现其他接口方法的占位符
func (r *MongoPetRepository) FindWithPagination(query *pet.PetQuery) (*pet.PetPageResult, error) {
	// TODO: 实现分页查询
	return nil, fmt.Errorf("not implemented")
}

func (r *MongoPetRepository) CountByCategory(category pet.PetCategory) (int64, error) {
	// TODO: 实现按类别统计
	return 0, fmt.Errorf("not implemented")
}

func (r *MongoPetRepository) FindDeadPets() ([]*pet.PetAggregate, error) {
	return r.FindByState(pet.PetStateDead)
}

func (r *MongoPetRepository) FindByLevelRange(minLevel, maxLevel uint32) ([]*pet.PetAggregate, error) {
	// TODO: 实现等级范围查询
	return nil, fmt.Errorf("not implemented")
}

func (r *MongoPetRepository) FindByStarRange(minStar, maxStar uint32) ([]*pet.PetAggregate, error) {
	// TODO: 实现星级范围查询
	return nil, fmt.Errorf("not implemented")
}

func (r *MongoPetRepository) SaveBatch(pets []*pet.PetAggregate) error {
	// TODO: 实现批量保存
	return fmt.Errorf("not implemented")
}

func (r *MongoPetRepository) DeleteBatch(ids []string) error {
	// TODO: 实现批量删除
	return fmt.Errorf("not implemented")
}

func (r *MongoPetRepository) FindTopPetsByPower(limit int) ([]*pet.PetAggregate, error) {
	// TODO: 实现按战力排序查询
	return nil, fmt.Errorf("not implemented")
}

func (r *MongoPetRepository) FindRecentlyCreated(duration time.Duration) ([]*pet.PetAggregate, error) {
	// TODO: 实现最近创建查询
	return nil, fmt.Errorf("not implemented")
}