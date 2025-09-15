// Package database Redis缓存操作封装
// Author: MMO Server Team
// Created: 2024

package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis Redis缓存管理器
type Redis struct {
	client *redis.Client
	config *RedisConfig
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr         string `json:"addr"`
	Password     string `json:"password"`
	DB           int    `json:"db"`
	PoolSize     int    `json:"pool_size"`
	MinIdleConns int    `json:"min_idle_conns"`
	MaxIdleConns int    `json:"max_idle_conns"`
	ConnMaxAge   int    `json:"conn_max_age"`
	DialTimeout  int    `json:"dial_timeout"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

// NewRedis 创建Redis管理器
func NewRedis(config *RedisConfig) *Redis {
	return &Redis{
		config: config,
	}
}

// Connect 连接到Redis
func (r *Redis) Connect(ctx context.Context) error {
	// 设置连接选项
	opts := &redis.Options{
		Addr:     r.config.Addr,
		Password: r.config.Password,
		DB:       r.config.DB,
	}
	
	if r.config.PoolSize > 0 {
		opts.PoolSize = r.config.PoolSize
	}
	if r.config.MinIdleConns > 0 {
		opts.MinIdleConns = r.config.MinIdleConns
	}
	if r.config.MaxIdleConns > 0 {
		opts.MaxIdleConns = r.config.MaxIdleConns
	}
	if r.config.ConnMaxAge > 0 {
		opts.ConnMaxLifetime = time.Duration(r.config.ConnMaxAge) * time.Second
	}
	if r.config.DialTimeout > 0 {
		opts.DialTimeout = time.Duration(r.config.DialTimeout) * time.Second
	}
	if r.config.ReadTimeout > 0 {
		opts.ReadTimeout = time.Duration(r.config.ReadTimeout) * time.Second
	}
	if r.config.WriteTimeout > 0 {
		opts.WriteTimeout = time.Duration(r.config.WriteTimeout) * time.Second
	}
	
	// 创建客户端
	r.client = redis.NewClient(opts)
	
	// 测试连接
	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	
	return nil
}

// Disconnect 断开连接
func (r *Redis) Disconnect() error {
	if r.client != nil {
		if err := r.client.Close(); err != nil {
			return fmt.Errorf("failed to disconnect from Redis: %w", err)
		}
	}
	return nil
}

// GetClient 获取客户端
func (r *Redis) GetClient() *redis.Client {
	return r.client
}

// Set 设置键值
func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Del 删除键
func (r *Redis) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (r *Redis) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func (r *Redis) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}