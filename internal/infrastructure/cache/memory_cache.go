package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	
	"greatestworks/aop/logger"
)

// MemoryCache 内存缓存实现
type MemoryCache struct {
	mu      sync.RWMutex
	items   map[string]*cacheItem
	logger  logger.Logger
	prefix  string
	cleaner *cacheCleaner
}

// cacheItem 缓存项
type cacheItem struct {
	value      interface{}
	expiration int64
	createdAt  time.Time
	accessedAt time.Time
	accessCount int64
}

// cacheCleaner 缓存清理器
type cacheCleaner struct {
	interval time.Duration
	stop     chan bool
	cache    *MemoryCache
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache(logger logger.Logger, prefix string, cleanupInterval time.Duration) Cache {
	c := &MemoryCache{
		items:  make(map[string]*cacheItem),
		logger: logger,
		prefix: prefix,
	}
	
	// 启动清理器
	if cleanupInterval > 0 {
		c.cleaner = &cacheCleaner{
			interval: cleanupInterval,
			stop:     make(chan bool),
			cache:    c,
		}
		go c.cleaner.run()
	}
	
	return c
}

// Set 设置缓存
func (m *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	fullKey := m.buildKey(key)
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}
	
	m.items[fullKey] = &cacheItem{
		value:       value,
		expiration:  expiration,
		createdAt:   time.Now(),
		accessedAt:  time.Now(),
		accessCount: 0,
	}
	
	m.logger.Debug("Memory cache set successfully", "key", key, "ttl", ttl)
	return nil
}

// Get 获取缓存
func (m *MemoryCache) Get(ctx context.Context, key string, dest interface{}) error {
	fullKey := m.buildKey(key)
	
	m.mu.RLock()
	item, found := m.items[fullKey]
	m.mu.RUnlock()
	
	if !found {
		return ErrCacheNotFound
	}
	
	// 检查是否过期
	if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
		m.mu.Lock()
		delete(m.items, fullKey)
		m.mu.Unlock()
		return ErrCacheExpired
	}
	
	// 更新访问信息
	m.mu.Lock()
	item.accessedAt = time.Now()
	item.accessCount++
	m.mu.Unlock()
	
	// 复制值到目标
	err := m.copyValue(item.value, dest)
	if err != nil {
		m.logger.Error("Failed to copy cache value", "error", err, "key", key)
		return fmt.Errorf("failed to copy cache value: %w", err)
	}
	
	m.logger.Debug("Memory cache get successfully", "key", key)
	return nil
}

// Delete 删除缓存
func (m *MemoryCache) Delete(ctx context.Context, key string) error {
	fullKey := m.buildKey(key)
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	delete(m.items, fullKey)
	
	m.logger.Debug("Memory cache deleted successfully", "key", key)
	return nil
}

// Exists 检查缓存是否存在
func (m *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	fullKey := m.buildKey(key)
	
	m.mu.RLock()
	item, found := m.items[fullKey]
	m.mu.RUnlock()
	
	if !found {
		return false, nil
	}
	
	// 检查是否过期
	if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
		m.mu.Lock()
		delete(m.items, fullKey)
		m.mu.Unlock()
		return false, nil
	}
	
	return true, nil
}

// SetBatch 批量设置缓存
func (m *MemoryCache) SetBatch(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}
	
	now := time.Now()
	for key, value := range items {
		fullKey := m.buildKey(key)
		m.items[fullKey] = &cacheItem{
			value:       value,
			expiration:  expiration,
			createdAt:   now,
			accessedAt:  now,
			accessCount: 0,
		}
	}
	
	m.logger.Debug("Memory batch cache set successfully", "count", len(items), "ttl", ttl)
	return nil
}

// GetBatch 批量获取缓存
func (m *MemoryCache) GetBatch(ctx context.Context, keys []string) (map[string]interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	results := make(map[string]interface{})
	now := time.Now().UnixNano()
	
	for _, key := range keys {
		fullKey := m.buildKey(key)
		item, found := m.items[fullKey]
		
		if !found {
			continue
		}
		
		// 检查是否过期
		if item.expiration > 0 && now > item.expiration {
			continue
		}
		
		// 更新访问信息（在读锁中不能修改，这里只是获取值）
		results[key] = item.value
	}
	
	m.logger.Debug("Memory batch cache get successfully", "requested", len(keys), "found", len(results))
	return results, nil
}

// DeleteBatch 批量删除缓存
func (m *MemoryCache) DeleteBatch(ctx context.Context, keys []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	for _, key := range keys {
		fullKey := m.buildKey(key)
		delete(m.items, fullKey)
	}
	
	m.logger.Debug("Memory batch cache deleted successfully", "count", len(keys))
	return nil
}

// SetTTL 设置缓存过期时间
func (m *MemoryCache) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	fullKey := m.buildKey(key)
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	item, found := m.items[fullKey]
	if !found {
		return ErrCacheNotFound
	}
	
	if ttl > 0 {
		item.expiration = time.Now().Add(ttl).UnixNano()
	} else {
		item.expiration = 0 // 永不过期
	}
	
	m.logger.Debug("Memory cache TTL set successfully", "key", key, "ttl", ttl)
	return nil
}

// GetTTL 获取缓存剩余过期时间
func (m *MemoryCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	fullKey := m.buildKey(key)
	
	m.mu.RLock()
	item, found := m.items[fullKey]
	m.mu.RUnlock()
	
	if !found {
		return 0, ErrCacheNotFound
	}
	
	if item.expiration == 0 {
		return -1, nil // 永不过期
	}
	
	ttl := time.Duration(item.expiration - time.Now().UnixNano())
	if ttl <= 0 {
		return 0, ErrCacheExpired
	}
	
	return ttl, nil
}

// Clear 清理所有缓存（根据前缀）
func (m *MemoryCache) Clear(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.prefix == "" {
		// 清空所有
		m.items = make(map[string]*cacheItem)
	} else {
		// 只清空匹配前缀的
		prefix := m.prefix + ":"
		for key := range m.items {
			if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
				delete(m.items, key)
			}
		}
	}
	
	m.logger.Info("Memory cache cleared successfully")
	return nil
}

// FlushDB 清空整个缓存
func (m *MemoryCache) FlushDB(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.items = make(map[string]*cacheItem)
	
	m.logger.Info("Memory cache flushed successfully")
	return nil
}

// Ping 健康检查
func (m *MemoryCache) Ping(ctx context.Context) error {
	// 内存缓存总是可用的
	return nil
}

// Close 关闭缓存
func (m *MemoryCache) Close() error {
	if m.cleaner != nil {
		m.cleaner.stop <- true
	}
	
	m.mu.Lock()
	m.items = nil
	m.mu.Unlock()
	
	m.logger.Info("Memory cache closed successfully")
	return nil
}

// GetStats 获取缓存统计信息
func (m *MemoryCache) GetStats() *CacheStats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	stats := &CacheStats{
		TotalItems: int64(len(m.items)),
		HitCount:   0,
		MissCount:  0,
	}
	
	now := time.Now().UnixNano()
	for _, item := range m.items {
		stats.HitCount += item.accessCount
		
		if item.expiration > 0 && now > item.expiration {
			stats.ExpiredItems++
		}
	}
	
	return stats
}

// 私有方法

// buildKey 构建完整的缓存键
func (m *MemoryCache) buildKey(key string) string {
	if m.prefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", m.prefix, key)
}

// copyValue 复制值到目标
func (m *MemoryCache) copyValue(src, dest interface{}) error {
	// 如果是相同类型的指针，直接赋值
	switch d := dest.(type) {
	case *string:
		if s, ok := src.(string); ok {
			*d = s
			return nil
		}
	case *int:
		if i, ok := src.(int); ok {
			*d = i
			return nil
		}
	case *int64:
		if i, ok := src.(int64); ok {
			*d = i
			return nil
		}
	case *float64:
		if f, ok := src.(float64); ok {
			*d = f
			return nil
		}
	case *bool:
		if b, ok := src.(bool); ok {
			*d = b
			return nil
		}
	}
	
	// 使用JSON进行深拷贝
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, dest)
}

// 清理器运行
func (c *cacheCleaner) run() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stop:
			return
		}
	}
}

// cleanup 清理过期项
func (c *cacheCleaner) cleanup() {
	c.cache.mu.Lock()
	defer c.cache.mu.Unlock()
	
	now := time.Now().UnixNano()
	expiredCount := 0
	
	for key, item := range c.cache.items {
		if item.expiration > 0 && now > item.expiration {
			delete(c.cache.items, key)
			expiredCount++
		}
	}
	
	if expiredCount > 0 {
		c.cache.logger.Debug("Cleaned up expired cache items", "count", expiredCount)
	}
}

// CacheStats 缓存统计信息
type CacheStats struct {
	TotalItems   int64 `json:"total_items"`
	ExpiredItems int64 `json:"expired_items"`
	HitCount     int64 `json:"hit_count"`
	MissCount    int64 `json:"miss_count"`
}

// MemoryCacheConfig 内存缓存配置
type MemoryCacheConfig struct {
	Prefix          string        `json:"prefix" yaml:"prefix"`
	CleanupInterval time.Duration `json:"cleanup_interval" yaml:"cleanup_interval"`
	MaxItems        int           `json:"max_items" yaml:"max_items"`
	DefaultTTL      time.Duration `json:"default_ttl" yaml:"default_ttl"`
}

// NewMemoryCacheFromConfig 从配置创建内存缓存
func NewMemoryCacheFromConfig(config *MemoryCacheConfig, logger logger.Logger) Cache {
	// 设置默认值
	if config.CleanupInterval == 0 {
		config.CleanupInterval = 10 * time.Minute
	}
	
	logger.Info("Memory cache initialized successfully", "prefix", config.Prefix, "cleanup_interval", config.CleanupInterval)
	
	return NewMemoryCache(logger, config.Prefix, config.CleanupInterval)
}