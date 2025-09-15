// Package persistence 数据持久化基础设施
package persistence

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoConfig MongoDB配置
type MongoConfig struct {
	URI            string        `json:"uri"`
	Database       string        `json:"database"`
	ConnectTimeout time.Duration `json:"connect_timeout"`
	MaxPoolSize    uint64        `json:"max_pool_size"`
	MinPoolSize    uint64        `json:"min_pool_size"`
}

// DefaultMongoConfig 默认MongoDB配置
func DefaultMongoConfig() *MongoConfig {
	return &MongoConfig{
		URI:            "mongodb://localhost:27017",
		Database:       "greatestworks",
		ConnectTimeout: 10 * time.Second,
		MaxPoolSize:    100,
		MinPoolSize:    5,
	}
}

// MongoDB MongoDB客户端包装器
type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
	config   *MongoConfig
}

// NewMongoDB 创建新的MongoDB客户端
func NewMongoDB(config *MongoConfig) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectTimeout)
	defer cancel()
	
	// 设置客户端选项
	clientOptions := options.Client().
		ApplyURI(config.URI).
		SetMaxPoolSize(config.MaxPoolSize).
		SetMinPoolSize(config.MinPoolSize).
		SetMaxConnIdleTime(30 * time.Minute)
	
	// 连接到MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("连接MongoDB失败: %w", err)
	}
	
	// 测试连接
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("MongoDB连接测试失败: %w", err)
	}
	
	database := client.Database(config.Database)
	
	return &MongoDB{
		client:   client,
		database: database,
		config:   config,
	}, nil
}

// GetDatabase 获取数据库实例
func (m *MongoDB) GetDatabase() *mongo.Database {
	return m.database
}

// GetCollection 获取集合
func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.database.Collection(name)
}

// Close 关闭连接
func (m *MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

// Ping 测试连接
func (m *MongoDB) Ping(ctx context.Context) error {
	return m.client.Ping(ctx, readpref.Primary())
}

// StartSession 开始会话
func (m *MongoDB) StartSession() (mongo.Session, error) {
	return m.client.StartSession()
}

// WithTransaction 执行事务
func (m *MongoDB) WithTransaction(ctx context.Context, fn func(mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	session, err := m.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)
	
	return session.WithTransaction(ctx, fn)
}

// CreateIndexes 创建索引
func (m *MongoDB) CreateIndexes(ctx context.Context, collectionName string, indexes []mongo.IndexModel) error {
	collection := m.GetCollection(collectionName)
	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// DropCollection 删除集合
func (m *MongoDB) DropCollection(ctx context.Context, collectionName string) error {
	return m.GetCollection(collectionName).Drop(ctx)
}

// ListCollections 列出所有集合
func (m *MongoDB) ListCollections(ctx context.Context) ([]string, error) {
	cursor, err := m.database.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return cursor, nil
}