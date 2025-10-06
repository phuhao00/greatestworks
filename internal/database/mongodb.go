// Package database 数据库连接池和操作封装
// Author: MMO Server Team
// Created: 2024

package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDB MongoDB数据库管理器
type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
	config   *MongoConfig
}

// MongoConfig MongoDB配置
type MongoConfig struct {
	URI            string `json:"uri"`
	Database       string `json:"database"`
	MaxPoolSize    uint64 `json:"max_pool_size"`
	MinPoolSize    uint64 `json:"min_pool_size"`
	MaxIdleTime    int    `json:"max_idle_time"`
	ConnectTimeout int    `json:"connect_timeout"`
	SocketTimeout  int    `json:"socket_timeout"`
}

// NewMongoDB 创建MongoDB管理器
func NewMongoDB(config *MongoConfig) *MongoDB {
	return &MongoDB{
		config: config,
	}
}

// Connect 连接到MongoDB
func (m *MongoDB) Connect(ctx context.Context) error {
	// 设置连接选项
	clientOptions := options.Client().ApplyURI(m.config.URI)

	if m.config.MaxPoolSize > 0 {
		clientOptions.SetMaxPoolSize(m.config.MaxPoolSize)
	}
	if m.config.MinPoolSize > 0 {
		clientOptions.SetMinPoolSize(m.config.MinPoolSize)
	}
	if m.config.MaxIdleTime > 0 {
		clientOptions.SetMaxConnIdleTime(time.Duration(m.config.MaxIdleTime) * time.Second)
	}
	if m.config.ConnectTimeout > 0 {
		clientOptions.SetConnectTimeout(time.Duration(m.config.ConnectTimeout) * time.Second)
	}
	if m.config.SocketTimeout > 0 {
		clientOptions.SetSocketTimeout(time.Duration(m.config.SocketTimeout) * time.Second)
	}

	// 创建客户端
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// 测试连接
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	m.client = client
	m.database = client.Database(m.config.Database)

	return nil
}

// Disconnect 断开连接
func (m *MongoDB) Disconnect(ctx context.Context) error {
	if m.client != nil {
		if err := m.client.Disconnect(ctx); err != nil {
			return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
		}
	}
	return nil
}

// GetCollection 获取集合
func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.database.Collection(name)
}

// GetDatabase 获取数据库
func (m *MongoDB) GetDatabase() *mongo.Database {
	return m.database
}

// GetClient 获取客户端
func (m *MongoDB) GetClient() *mongo.Client {
	return m.client
}
