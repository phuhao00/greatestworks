// Package persistence 统一仓储基类
// Author: MMO Server Team
// Created: 2024

package persistence

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"greatestworks/internal/infrastructure/cache"
	"greatestworks/internal/infrastructure/logging"
)

// BaseRepository 基础仓储实现
type BaseRepository struct {
	db             *mongo.Database
	cache          cache.Cache
	logger         logging.Logger
	collection     *mongo.Collection
	collectionName string
}

// NewBaseRepository 创建基础仓储
func NewBaseRepository(db *mongo.Database, cache cache.Cache, logger logging.Logger, collectionName string) *BaseRepository {
	return &BaseRepository{
		db:             db,
		cache:          cache,
		logger:         logger,
		collection:     db.Collection(collectionName),
		collectionName: collectionName,
	}
}

// GetCollection 获取集合
func (r *BaseRepository) GetCollection() *mongo.Collection {
	return r.collection
}

// GetDB 获取数据库
func (r *BaseRepository) GetDB() *mongo.Database {
	return r.db
}

// GetCache 获取缓存
func (r *BaseRepository) GetCache() cache.Cache {
	return r.cache
}

// GetLogger 获取日志器
func (r *BaseRepository) GetLogger() logging.Logger {
	return r.logger
}

// Save 保存文档
func (r *BaseRepository) Save(ctx context.Context, id string, document interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid object ID: %w", err)
	}

	opts := options.Replace().SetUpsert(true)
	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objectID}, document, opts)
	if err != nil {
		r.logger.Error("Failed to save document", err, logging.Fields{
			"collection": r.collectionName,
			"id":         id,
		})
		return fmt.Errorf("failed to save document: %w", err)
	}

	// 清除缓存
	if r.cache != nil {
		cacheKey := r.buildCacheKey(id)
		r.cache.Delete(ctx, cacheKey)
	}

	r.logger.Debug("Document saved successfully", logging.Fields{
		"collection": r.collectionName,
		"id":         id,
	})

	return nil
}

// FindByID 根据ID查找文档
func (r *BaseRepository) FindByID(ctx context.Context, id string, result interface{}) error {
	// 先尝试从缓存获取
	if r.cache != nil {
		cacheKey := r.buildCacheKey(id)
		if err := r.cache.Get(ctx, cacheKey, result); err == nil {
			r.logger.Debug("Document found in cache", logging.Fields{
				"collection": r.collectionName,
				"id":         id,
			})
			return nil
		}
	}

	// 从数据库获取
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid object ID: %w", err)
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("document not found")
		}
		r.logger.Error("Failed to find document", err, logging.Fields{
			"collection": r.collectionName,
			"id":         id,
		})
		return fmt.Errorf("failed to find document: %w", err)
	}

	// 缓存结果
	if r.cache != nil {
		cacheKey := r.buildCacheKey(id)
		r.cache.Set(ctx, cacheKey, result, time.Hour)
	}

	r.logger.Debug("Document found in database", logging.Fields{
		"collection": r.collectionName,
		"id":         id,
	})

	return nil
}

// FindOne 查找单个文档
func (r *BaseRepository) FindOne(ctx context.Context, filter bson.M, result interface{}) error {
	err := r.collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("document not found")
		}
		r.logger.Error("Failed to find document", err, logging.Fields{
			"collection": r.collectionName,
			"filter":     filter,
			"error":      err,
		})
		return fmt.Errorf("failed to find document: %w", err)
	}

	return nil
}

// FindMany 查找多个文档
func (r *BaseRepository) FindMany(ctx context.Context, filter bson.M, results interface{}, opts ...*options.FindOptions) error {
	cursor, err := r.collection.Find(ctx, filter, opts...)
	if err != nil {
		r.logger.Error("Failed to find documents", err, logging.Fields{
			"collection": r.collectionName,
			"filter":     filter,
		})
		return fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, results)
	if err != nil {
		r.logger.Error("Failed to decode documents", err, logging.Fields{
			"collection": r.collectionName,
		})
		return fmt.Errorf("failed to decode documents: %w", err)
	}

	return nil
}

// Delete 删除文档
func (r *BaseRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid object ID: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		r.logger.Error("Failed to delete document", err, logging.Fields{
			"collection": r.collectionName,
			"id":         id,
			"error":      err,
		})
		return fmt.Errorf("failed to delete document: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("document not found")
	}

	// 清除缓存
	if r.cache != nil {
		cacheKey := r.buildCacheKey(id)
		r.cache.Delete(ctx, cacheKey)
	}

	r.logger.Debug("Document deleted successfully", logging.Fields{
		"collection": r.collectionName,
		"id":         id,
	})

	return nil
}

// Count 统计文档数量
func (r *BaseRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to count documents", err, logging.Fields{
			"collection": r.collectionName,
			"filter":     filter,
		})
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}

	return count, nil
}

// Exists 检查文档是否存在
func (r *BaseRepository) Exists(ctx context.Context, filter bson.M) (bool, error) {
	count, err := r.Count(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateIndex 创建索引
func (r *BaseRepository) CreateIndex(ctx context.Context, index mongo.IndexModel) error {
	_, err := r.collection.Indexes().CreateOne(ctx, index)
	if err != nil {
		r.logger.Error("Failed to create index", err, logging.Fields{
			"collection": r.collectionName,
			"error":      err,
		})
		return fmt.Errorf("failed to create index: %w", err)
	}

	r.logger.Debug("Index created successfully", logging.Fields{
		"collection": r.collectionName,
	})

	return nil
}

// CreateIndexes 创建多个索引
func (r *BaseRepository) CreateIndexes(ctx context.Context, indexes []mongo.IndexModel) error {
	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		r.logger.Error("Failed to create indexes", err, logging.Fields{
			"collection": r.collectionName,
		})
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	r.logger.Debug("Indexes created successfully", logging.Fields{
		"collection": r.collectionName,
		"count":      len(indexes),
	})

	return nil
}

// Aggregate 聚合查询
func (r *BaseRepository) Aggregate(ctx context.Context, pipeline mongo.Pipeline, results interface{}) error {
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Error("Failed to execute aggregation", err, logging.Fields{
			"collection": r.collectionName,
		})
		return fmt.Errorf("failed to execute aggregation: %w", err)
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, results)
	if err != nil {
		r.logger.Error("Failed to decode aggregation results", err, logging.Fields{
			"collection": r.collectionName,
		})
		return fmt.Errorf("failed to decode aggregation results: %w", err)
	}

	return nil
}

// UpdateOne 更新单个文档
func (r *BaseRepository) UpdateOne(ctx context.Context, filter bson.M, update bson.M) error {
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("Failed to update document", err, logging.Fields{
			"collection": r.collectionName,
			"filter":     filter,
			"error":      err,
		})
		return fmt.Errorf("failed to update document: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("document not found")
	}

	r.logger.Debug("Document updated successfully", logging.Fields{
		"collection": r.collectionName,
		"matched":    result.MatchedCount,
		"modified":   result.ModifiedCount,
	})

	return nil
}

// UpdateMany 更新多个文档
func (r *BaseRepository) UpdateMany(ctx context.Context, filter bson.M, update bson.M) error {
	result, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		r.logger.Error("Failed to update documents", err, logging.Fields{
			"collection": r.collectionName,
			"filter":     filter,
		})
		return fmt.Errorf("failed to update documents: %w", err)
	}

	r.logger.Debug("Documents updated successfully", logging.Fields{
		"collection": r.collectionName,
		"matched":    result.MatchedCount,
		"modified":   result.ModifiedCount,
	})

	return nil
}

// buildCacheKey 构建缓存键
func (r *BaseRepository) buildCacheKey(id string) string {
	return fmt.Sprintf("%s:%s", r.collectionName, id)
}

// WithTransaction 执行事务
func (r *BaseRepository) WithTransaction(ctx context.Context, fn func(mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	session, err := r.db.Client().StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	return session.WithTransaction(ctx, fn)
}

// GetStats 获取仓储统计信息
func (r *BaseRepository) GetStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 获取集合统计信息
	count, err := r.Count(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get document count: %w", err)
	}

	stats["collection"] = r.collectionName
	stats["document_count"] = count

	// 获取索引信息
	indexes, err := r.collection.Indexes().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list indexes: %w", err)
	}
	defer indexes.Close(ctx)

	var indexList []bson.M
	if err := indexes.All(ctx, &indexList); err != nil {
		return nil, fmt.Errorf("failed to decode indexes: %w", err)
	}

	stats["index_count"] = len(indexList)

	return stats, nil
}
