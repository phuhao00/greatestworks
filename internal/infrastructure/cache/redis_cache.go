package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"greatestworks/internal/infrastructure/logging"

	"github.com/redis/go-redis/v9"
)

// RedisCache Redis缓存实现
type RedisCache struct {
	client *redis.Client
	logger logging.Logger
}

// NewRedisCache 创建Redis缓存
func NewRedisCache(client *redis.Client, logger logging.Logger) *RedisCache {
	return &RedisCache{
		client: client,
		logger: logger,
	}
}

// Set 设置值
func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}

	err = rc.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("设置缓存失败: %w", err)
	}

	rc.logger.Info("缓存设置成功", map[string]interface{}{
		"key": key,
		"ttl": ttl,
	})

	return nil
}

// Get 获取值
func (rc *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("获取缓存失败: %w", err)
	}

	err = json.Unmarshal([]byte(data), dest)
	if err != nil {
		return fmt.Errorf("反序列化失败: %w", err)
	}

	rc.logger.Info("缓存获取成功", map[string]interface{}{
		"key": key,
	})

	return nil
}

// Delete 删除值
func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	err := rc.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("删除缓存失败: %w", err)
	}

	rc.logger.Info("缓存删除成功", map[string]interface{}{
		"key": key,
	})

	return nil
}

// Exists 检查键是否存在
func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := rc.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("检查缓存存在性失败: %w", err)
	}

	return count > 0, nil
}

// SetBatch 批量设置
func (rc *RedisCache) SetBatch(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
	pipe := rc.client.Pipeline()

	for key, value := range items {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("序列化失败: %w", err)
		}

		pipe.Set(ctx, key, data, ttl)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("批量设置缓存失败: %w", err)
	}

	rc.logger.Info("批量缓存设置成功", map[string]interface{}{
		"count": len(items),
		"ttl":   ttl,
	})

	return nil
}

// GetBatch 批量获取
func (rc *RedisCache) GetBatch(ctx context.Context, keys []string) (map[string]interface{}, error) {
	pipe := rc.client.Pipeline()

	cmds := make([]*redis.StringCmd, len(keys))
	for i, key := range keys {
		cmds[i] = pipe.Get(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("批量获取缓存失败: %w", err)
	}

	result := make(map[string]interface{})
	for i, cmd := range cmds {
		if cmd.Err() == nil {
			var value interface{}
			err := json.Unmarshal([]byte(cmd.Val()), &value)
			if err == nil {
				result[keys[i]] = value
			}
		}
	}

	rc.logger.Info("批量缓存获取成功", map[string]interface{}{
		"requested": len(keys),
		"found":     len(result),
	})

	return result, nil
}

// DeleteBatch 批量删除
func (rc *RedisCache) DeleteBatch(ctx context.Context, keys []string) error {
	err := rc.client.Del(ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("批量删除缓存失败: %w", err)
	}

	rc.logger.Info("批量缓存删除成功", map[string]interface{}{
		"count": len(keys),
	})

	return nil
}

// SetTTL 设置TTL
func (rc *RedisCache) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	err := rc.client.Expire(ctx, key, ttl).Err()
	if err != nil {
		return fmt.Errorf("设置TTL失败: %w", err)
	}

	rc.logger.Info("TTL设置成功", map[string]interface{}{
		"key": key,
		"ttl": ttl,
	})

	return nil
}

// GetTTL 获取TTL
func (rc *RedisCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := rc.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("获取TTL失败: %w", err)
	}

	return ttl, nil
}

// Clear 清空缓存
func (rc *RedisCache) Clear(ctx context.Context) error {
	err := rc.client.FlushDB(ctx).Err()
	if err != nil {
		return fmt.Errorf("清空缓存失败: %w", err)
	}

	rc.logger.Info("缓存清空成功")

	return nil
}

// Size 获取缓存大小
func (rc *RedisCache) Size() int64 {
	// Redis没有直接的方法获取数据库大小
	// 这里返回0，实际使用中可以通过INFO命令获取
	return 0
}

// Keys 获取所有键
func (rc *RedisCache) Keys() []string {
	// 注意：在生产环境中，KEYS命令可能会阻塞Redis
	// 建议使用SCAN命令进行迭代
	keys, err := rc.client.Keys(context.Background(), "*").Result()
	if err != nil {
		rc.logger.Error("Failed to get keys", err)
		return []string{}
	}

	return keys
}
