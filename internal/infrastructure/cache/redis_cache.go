package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	
	"github.com/redis/go-redis/v9"
	"greatestworks/aop/logger"
)

// Cache 缓存接口
type Cache interface {
	// 基础操作
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	
	// 批量操作
	SetBatch(ctx context.Context, items map[string]interface{}, ttl time.Duration) error
	GetBatch(ctx context.Context, keys []string) (map[string]interface{}, error)
	DeleteBatch(ctx context.Context, keys []string) error
	
	// TTL操作
	SetTTL(ctx context.Context, key string, ttl time.Duration) error
	GetTTL(ctx context.Context, key string) (time.Duration, error)
	
	// 清理操作
	Clear(ctx context.Context) error
	FlushDB(ctx context.Context) error
	
	// 健康检查
	Ping(ctx context.Context) error
	Close() error
}

// RedisCache Redis缓存实现
type RedisCache struct {
	client *redis.Client
	logger logger.Logger
	prefix string
}

// NewRedisCache 创建Redis缓存
func NewRedisCache(client *redis.Client, logger logger.Logger, prefix string) Cache {
	return &RedisCache{
		client: client,
		logger: logger,
		prefix: prefix,
	}
}

// Set 设置缓存
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	fullKey := r.buildKey(key)
	
	// 序列化值
	data, err := r.serialize(value)
	if err != nil {
		r.logger.Error("Failed to serialize cache value", "error", err, "key", key)
		return fmt.Errorf("failed to serialize cache value: %w", err)
	}
	
	// 设置到Redis
	err = r.client.Set(ctx, fullKey, data, ttl).Err()
	if err != nil {
		r.logger.Error("Failed to set cache", "error", err, "key", key, "ttl", ttl)
		return fmt.Errorf("failed to set cache: %w", err)
	}
	
	r.logger.Debug("Cache set successfully", "key", key, "ttl", ttl)
	return nil
}

// Get 获取缓存
func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	fullKey := r.buildKey(key)
	
	// 从Redis获取
	data, err := r.client.Get(ctx, fullKey).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheNotFound
		}
		r.logger.Error("Failed to get cache", "error", err, "key", key)
		return fmt.Errorf("failed to get cache: %w", err)
	}
	
	// 反序列化值
	err = r.deserialize([]byte(data), dest)
	if err != nil {
		r.logger.Error("Failed to deserialize cache value", "error", err, "key", key)
		return fmt.Errorf("failed to deserialize cache value: %w", err)
	}
	
	r.logger.Debug("Cache get successfully", "key", key)
	return nil
}

// Delete 删除缓存
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	fullKey := r.buildKey(key)
	
	err := r.client.Del(ctx, fullKey).Err()
	if err != nil {
		r.logger.Error("Failed to delete cache", "error", err, "key", key)
		return fmt.Errorf("failed to delete cache: %w", err)
	}
	
	r.logger.Debug("Cache deleted successfully", "key", key)
	return nil
}

// Exists 检查缓存是否存在
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	fullKey := r.buildKey(key)
	
	count, err := r.client.Exists(ctx, fullKey).Result()
	if err != nil {
		r.logger.Error("Failed to check cache existence", "error", err, "key", key)
		return false, fmt.Errorf("failed to check cache existence: %w", err)
	}
	
	return count > 0, nil
}

// SetBatch 批量设置缓存
func (r *RedisCache) SetBatch(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
	pipe := r.client.Pipeline()
	
	for key, value := range items {
		fullKey := r.buildKey(key)
		
		// 序列化值
		data, err := r.serialize(value)
		if err != nil {
			r.logger.Error("Failed to serialize batch cache value", "error", err, "key", key)
			continue
		}
		
		pipe.Set(ctx, fullKey, data, ttl)
	}
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		r.logger.Error("Failed to set batch cache", "error", err, "count", len(items))
		return fmt.Errorf("failed to set batch cache: %w", err)
	}
	
	r.logger.Debug("Batch cache set successfully", "count", len(items), "ttl", ttl)
	return nil
}

// GetBatch 批量获取缓存
func (r *RedisCache) GetBatch(ctx context.Context, keys []string) (map[string]interface{}, error) {
	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = r.buildKey(key)
	}
	
	results, err := r.client.MGet(ctx, fullKeys...).Result()
	if err != nil {
		r.logger.Error("Failed to get batch cache", "error", err, "count", len(keys))
		return nil, fmt.Errorf("failed to get batch cache: %w", err)
	}
	
	batchResults := make(map[string]interface{})
	for i, result := range results {
		if result != nil {
			var value interface{}
			if err := r.deserialize([]byte(result.(string)), &value); err == nil {
				batchResults[keys[i]] = value
			}
		}
	}
	
	r.logger.Debug("Batch cache get successfully", "requested", len(keys), "found", len(batchResults))
	return batchResults, nil
}

// DeleteBatch 批量删除缓存
func (r *RedisCache) DeleteBatch(ctx context.Context, keys []string) error {
	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = r.buildKey(key)
	}
	
	err := r.client.Del(ctx, fullKeys...).Err()
	if err != nil {
		r.logger.Error("Failed to delete batch cache", "error", err, "count", len(keys))
		return fmt.Errorf("failed to delete batch cache: %w", err)
	}
	
	r.logger.Debug("Batch cache deleted successfully", "count", len(keys))
	return nil
}

// SetTTL 设置缓存过期时间
func (r *RedisCache) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	fullKey := r.buildKey(key)
	
	err := r.client.Expire(ctx, fullKey, ttl).Err()
	if err != nil {
		r.logger.Error("Failed to set cache TTL", "error", err, "key", key, "ttl", ttl)
		return fmt.Errorf("failed to set cache TTL: %w", err)
	}
	
	r.logger.Debug("Cache TTL set successfully", "key", key, "ttl", ttl)
	return nil
}

// GetTTL 获取缓存剩余过期时间
func (r *RedisCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	fullKey := r.buildKey(key)
	
	ttl, err := r.client.TTL(ctx, fullKey).Result()
	if err != nil {
		r.logger.Error("Failed to get cache TTL", "error", err, "key", key)
		return 0, fmt.Errorf("failed to get cache TTL: %w", err)
	}
	
	return ttl, nil
}

// Clear 清理所有缓存（根据前缀）
func (r *RedisCache) Clear(ctx context.Context) error {
	pattern := r.buildKey("*")
	
	// 使用SCAN命令获取所有匹配的键
	var cursor uint64
	var keys []string
	
	for {
		var scanKeys []string
		var err error
		
		scanKeys, cursor, err = r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			r.logger.Error("Failed to scan cache keys", "error", err, "pattern", pattern)
			return fmt.Errorf("failed to scan cache keys: %w", err)
		}
		
		keys = append(keys, scanKeys...)
		
		if cursor == 0 {
			break
		}
	}
	
	// 批量删除键
	if len(keys) > 0 {
		err := r.client.Del(ctx, keys...).Err()
		if err != nil {
			r.logger.Error("Failed to clear cache", "error", err, "count", len(keys))
			return fmt.Errorf("failed to clear cache: %w", err)
		}
	}
	
	r.logger.Info("Cache cleared successfully", "count", len(keys))
	return nil
}

// FlushDB 清空整个数据库
func (r *RedisCache) FlushDB(ctx context.Context) error {
	err := r.client.FlushDB(ctx).Err()
	if err != nil {
		r.logger.Error("Failed to flush database", "error", err)
		return fmt.Errorf("failed to flush database: %w", err)
	}
	
	r.logger.Info("Database flushed successfully")
	return nil
}

// Ping 健康检查
func (r *RedisCache) Ping(ctx context.Context) error {
	err := r.client.Ping(ctx).Err()
	if err != nil {
		r.logger.Error("Redis ping failed", "error", err)
		return fmt.Errorf("redis ping failed: %w", err)
	}
	
	return nil
}

// Close 关闭连接
func (r *RedisCache) Close() error {
	err := r.client.Close()
	if err != nil {
		r.logger.Error("Failed to close Redis client", "error", err)
		return fmt.Errorf("failed to close Redis client: %w", err)
	}
	
	r.logger.Info("Redis client closed successfully")
	return nil
}

// 私有方法

// buildKey 构建完整的缓存键
func (r *RedisCache) buildKey(key string) string {
	if r.prefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", r.prefix, key)
}

// serialize 序列化值
func (r *RedisCache) serialize(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		return json.Marshal(v)
	default:
		return json.Marshal(v)
	}
}

// deserialize 反序列化值
func (r *RedisCache) deserialize(data []byte, dest interface{}) error {
	switch dest.(type) {
	case *string:
		*(dest.(*string)) = string(data)
		return nil
	case *[]byte:
		*(dest.(*[]byte)) = data
		return nil
	default:
		return json.Unmarshal(data, dest)
	}
}

// 错误定义
var (
	ErrCacheNotFound = fmt.Errorf("cache not found")
	ErrCacheExpired  = fmt.Errorf("cache expired")
	ErrCacheInvalid  = fmt.Errorf("cache data invalid")
)

// CacheConfig Redis缓存配置
type CacheConfig struct {
	Addr         string        `json:"addr" yaml:"addr"`
	Password     string        `json:"password" yaml:"password"`
	DB           int           `json:"db" yaml:"db"`
	PoolSize     int           `json:"pool_size" yaml:"pool_size"`
	MinIdleConns int           `json:"min_idle_conns" yaml:"min_idle_conns"`
	MaxRetries   int           `json:"max_retries" yaml:"max_retries"`
	DialTimeout  time.Duration `json:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
	Prefix       string        `json:"prefix" yaml:"prefix"`
}

// NewRedisCacheFromConfig 从配置创建Redis缓存
func NewRedisCacheFromConfig(config *CacheConfig, logger logger.Logger) (Cache, error) {
	// 设置默认值
	if config.PoolSize == 0 {
		config.PoolSize = 10
	}
	if config.MinIdleConns == 0 {
		config.MinIdleConns = 5
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.DialTimeout == 0 {
		config.DialTimeout = 5 * time.Second
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = 3 * time.Second
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 3 * time.Second
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = 5 * time.Minute
	}
	
	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	})
	
	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	
	logger.Info("Redis cache initialized successfully", "addr", config.Addr, "db", config.DB, "prefix", config.Prefix)
	
	return NewRedisCache(client, logger, config.Prefix), nil
}