package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logging"
)

// MemoryCache 内存缓存实现
type MemoryCache struct {
	data            map[string]*cacheItem
	mutex           sync.RWMutex
	logger          logging.Logger
	maxSize         int64
	cleanupInterval time.Duration
}

// cacheItem 缓存项
type cacheItem struct {
	value     interface{}
	expiresAt time.Time
	createdAt time.Time
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache(logger logging.Logger, maxSize int64, cleanupInterval time.Duration) *MemoryCache {
	cache := &MemoryCache{
		data:            make(map[string]*cacheItem),
		logger:          logger,
		maxSize:         maxSize,
		cleanupInterval: cleanupInterval,
	}

	// 启动清理例程
	go cache.startCleanupRoutine()

	return cache
}

// Get 获取值
func (c *MemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	// 检查是否过期
	if time.Now().After(item.expiresAt) {
		delete(c.data, key)
		return nil, fmt.Errorf("key expired: %s", key)
	}

	c.logger.Info("缓存命中", map[string]interface{}{
		"key": key,
	})

	return item.value, nil
}

// Set 设置值
func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 检查缓存大小
	if int64(len(c.data)) >= c.maxSize {
		c.evictOldest()
	}

	c.data[key] = &cacheItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
		createdAt: time.Now(),
	}

	c.logger.Info("缓存设置", map[string]interface{}{
		"key": key,
		"ttl": ttl,
	})

	return nil
}

// Delete 删除值
func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)

	c.logger.Info("缓存删除", map[string]interface{}{
		"key": key,
	})

	return nil
}

// Exists 检查键是否存在
func (c *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return false, nil
	}

	// 检查是否过期
	if time.Now().After(item.expiresAt) {
		delete(c.data, key)
		return false, nil
	}

	return true, nil
}

// Clear 清空缓存
func (c *MemoryCache) Clear(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]*cacheItem)

	c.logger.Info("缓存清空")

	return nil
}

// Size 获取缓存大小
func (c *MemoryCache) Size() int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return int64(len(c.data))
}

// Keys 获取所有键
func (c *MemoryCache) Keys() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	keys := make([]string, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}

	return keys
}

// evictOldest 驱逐最旧的项
func (c *MemoryCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, item := range c.data {
		if oldestKey == "" || item.createdAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.createdAt
		}
	}

	if oldestKey != "" {
		delete(c.data, oldestKey)
		c.logger.Info("驱逐最旧缓存项", map[string]interface{}{
			"key": oldestKey,
		})
	}
}

// startCleanupRoutine 启动清理例程
func (c *MemoryCache) startCleanupRoutine() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanupExpired()
	}
}

// cleanupExpired 清理过期项
func (c *MemoryCache) cleanupExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	expiredCount := 0

	for key, item := range c.data {
		if now.After(item.expiresAt) {
			delete(c.data, key)
			expiredCount++
		}
	}

	if expiredCount > 0 {
		c.logger.Info("清理过期缓存项", map[string]interface{}{
			"expired_count": expiredCount,
		})
	}
}
