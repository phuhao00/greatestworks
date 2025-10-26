package persistence

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepository 用户仓储
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, user *DbUser) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID 根据ID查找用户
func (r *UserRepository) FindByID(ctx context.Context, userID int64) (*DbUser, error) {
	var user DbUser
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*DbUser, error) {
	var user DbUser
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *UserRepository) Update(ctx context.Context, user *DbUser) error {
	user.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"user_id": user.UserID},
		bson.M{"$set": user},
	)
	return err
}

// UpdateLastLogin 更新最后登录时间
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID int64) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"user_id": userID},
		bson.M{"$set": bson.M{"last_login_at": time.Now()}},
	)
	return err
}

// CharacterRepository 角色仓储
type CharacterRepository struct {
	collection *mongo.Collection
}

// NewCharacterRepository 创建角色仓储
func NewCharacterRepository(db *mongo.Database) *CharacterRepository {
	return &CharacterRepository{
		collection: db.Collection("characters"),
	}
}

// Create 创建角色
func (r *CharacterRepository) Create(ctx context.Context, character *DbCharacter) error {
	character.CreatedAt = time.Now()
	character.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, character)
	if err != nil {
		return err
	}
	character.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID 根据ID查找角色
func (r *CharacterRepository) FindByID(ctx context.Context, characterID int64) (*DbCharacter, error) {
	var character DbCharacter
	err := r.collection.FindOne(ctx, bson.M{
		"character_id": characterID,
		"deleted_at":   bson.M{"$exists": false},
	}).Decode(&character)
	if err != nil {
		return nil, err
	}
	return &character, nil
}

// FindByUserID 根据用户ID查找所有角色
func (r *CharacterRepository) FindByUserID(ctx context.Context, userID int64) ([]*DbCharacter, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"user_id":    userID,
		"deleted_at": bson.M{"$exists": false},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var characters []*DbCharacter
	if err := cursor.All(ctx, &characters); err != nil {
		return nil, err
	}
	return characters, nil
}

// Update 更新角色
func (r *CharacterRepository) Update(ctx context.Context, character *DbCharacter) error {
	character.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"character_id": character.CharacterID},
		bson.M{"$set": character},
	)
	return err
}

// Delete 软删除角色
func (r *CharacterRepository) Delete(ctx context.Context, characterID int64) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"character_id": characterID},
		bson.M{"$set": bson.M{
			"deleted_at": time.Now(),
		}},
	)
	return err
}

// UpdatePosition 更新角色位置
func (r *CharacterRepository) UpdatePosition(ctx context.Context, characterID int64, mapID int32, x, y, z, dir float32) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"character_id": characterID},
		bson.M{"$set": bson.M{
			"map_id":     mapID,
			"position_x": x,
			"position_y": y,
			"position_z": z,
			"direction":  dir,
			"updated_at": time.Now(),
		}},
	)
	return err
}

// ItemRepository 物品仓储
type ItemRepository struct {
	collection *mongo.Collection
}

// NewItemRepository 创建物品仓储
func NewItemRepository(db *mongo.Database) *ItemRepository {
	return &ItemRepository{
		collection: db.Collection("items"),
	}
}

// Create 创建物品
func (r *ItemRepository) Create(ctx context.Context, item *DbItem) error {
	item.CreatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return err
	}
	item.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByCharacterID 查找角色的所有物品
func (r *ItemRepository) FindByCharacterID(ctx context.Context, characterID int64) ([]*DbItem, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"character_id": characterID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*DbItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// FindByUID 根据唯一ID查找物品
func (r *ItemRepository) FindByUID(ctx context.Context, itemUID int64) (*DbItem, error) {
	var item DbItem
	err := r.collection.FindOne(ctx, bson.M{"item_uid": itemUID}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Update 更新物品
func (r *ItemRepository) Update(ctx context.Context, item *DbItem) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"item_uid": item.ItemUID},
		bson.M{"$set": item},
	)
	return err
}

// Delete 删除物品
func (r *ItemRepository) Delete(ctx context.Context, itemUID int64) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"item_uid": itemUID})
	return err
}

// QuestRepository 任务仓储
type QuestRepository struct {
	collection *mongo.Collection
}

// NewQuestRepository 创建任务仓储
func NewQuestRepository(db *mongo.Database) *QuestRepository {
	return &QuestRepository{
		collection: db.Collection("quests"),
	}
}

// Create 创建任务进度
func (r *QuestRepository) Create(ctx context.Context, quest *DbQuest) error {
	quest.AcceptedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, quest)
	if err != nil {
		return err
	}
	quest.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByCharacterID 查找角色的所有任务
func (r *QuestRepository) FindByCharacterID(ctx context.Context, characterID int64) ([]*DbQuest, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"character_id": characterID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var quests []*DbQuest
	if err := cursor.All(ctx, &quests); err != nil {
		return nil, err
	}
	return quests, nil
}

// Update 更新任务进度
func (r *QuestRepository) Update(ctx context.Context, quest *DbQuest) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"character_id": quest.CharacterID,
			"quest_id":     quest.QuestID,
		},
		bson.M{"$set": quest},
	)
	return err
}

// MailRepository 邮件仓储
type MailRepository struct {
	collection *mongo.Collection
}

// NewMailRepository 创建邮件仓储
func NewMailRepository(db *mongo.Database) *MailRepository {
	return &MailRepository{
		collection: db.Collection("mails"),
	}
}

// Create 创建邮件
func (r *MailRepository) Create(ctx context.Context, mail *DbMail) error {
	mail.CreatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, mail)
	if err != nil {
		return err
	}
	mail.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByReceiverID 查找收件人的邮件
func (r *MailRepository) FindByReceiverID(ctx context.Context, receiverID int64, limit int) ([]*DbMail, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit))
	cursor, err := r.collection.Find(ctx, bson.M{"receiver_id": receiverID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mails []*DbMail
	if err := cursor.All(ctx, &mails); err != nil {
		return nil, err
	}
	return mails, nil
}

// MarkAsRead 标记为已读
func (r *MailRepository) MarkAsRead(ctx context.Context, mailID int64) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"mail_id": mailID},
		bson.M{"$set": bson.M{"is_read": true}},
	)
	return err
}

// Delete 删除邮件
func (r *MailRepository) Delete(ctx context.Context, mailID int64) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"mail_id": mailID})
	return err
}

// DeleteExpired 删除过期邮件
func (r *MailRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{
		"expire_at": bson.M{"$lt": time.Now()},
	})
	return err
}
