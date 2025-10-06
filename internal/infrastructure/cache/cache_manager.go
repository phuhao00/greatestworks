package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logger"
)

// CacheManager 缓存管理器
type CacheManager struct {
	primary   Cache // 主缓存（通常是Redis）
	secondary Cache // 备用缓存（通常是内存缓存）
	logger    logger.Logger
	config    *CacheManagerConfig
	mu        sync.RWMutex
	stats     *ManagerStats
}

// CacheManagerConfig 缓存管理器配置
type CacheManagerConfig struct {
	UseFallback         bool          `json:"use_fallback" yaml:"use_fallback"`                   // 是否使用备用缓存
	FallbackOnError     bool          `json:"fallback_on_error" yaml:"fallback_on_error"`         // 错误时是否回退到备用缓存
	SyncToSecondary     bool          `json:"sync_to_secondary" yaml:"sync_to_secondary"`         // 是否同步到备用缓存
	HealthCheckInterval time.Duration `json:"health_check_interval" yaml:"health_check_interval"` // 健康检查间隔
	RetryAttempts       int           `json:"retry_attempts" yaml:"retry_attempts"`               // 重试次数
	RetryDelay          time.Duration `json:"retry_delay" yaml:"retry_delay"`                     // 重试延迟
}

// ManagerStats 管理器统计信息
type ManagerStats struct {
	PrimaryHits     int64 `json:"primary_hits"`
	SecondaryHits   int64 `json:"secondary_hits"`
	PrimaryErrors   int64 `json:"primary_errors"`
	SecondaryErrors int64 `json:"secondary_errors"`
	FallbackCount   int64 `json:"fallback_count"`
	SyncCount       int64 `json:"sync_count"`
}

// NewCacheManager 创建缓存管理器
func NewCacheManager(primary, secondary Cache, config *CacheManagerConfig, logger logger.Logger) *CacheManager {
	if config == nil {
		config = &CacheManagerConfig{
			UseFallback:         true,
			FallbackOnError:     true,
			SyncToSecondary:     false,
			HealthCheckInterval: 30 * time.Second,
			RetryAttempts:       3,
			RetryDelay:          100 * time.Millisecond,
		}
	}

	m := &CacheManager{
		primary:   primary,
		secondary: secondary,
		logger:    logger,
		config:    config,
		stats:     &ManagerStats{},
	}

	// 启动健康检查
	if config.HealthCheckInterval > 0 {
		go m.healthCheck()
	}

	logger.Info("Cache manager initialized successfully", "use_fallback", config.UseFallback, "sync_to_secondary", config.SyncToSecondary)
	return m
}

// Set 设置缓存
func (m *CacheManager) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// 尝试设置主缓存
	err := m.executeWithRetry(func() error {
		return m.primary.Set(ctx, key, value, ttl)
	})

	if err != nil {
		m.mu.Lock()
		m.stats.PrimaryErrors++
		m.mu.Unlock()

		m.logger.Error("Failed to set primary cache", "error", err, "key", key)

		// 如果配置了错误时回退，尝试设置备用缓存
		if m.config.FallbackOnError && m.secondary != nil {
			if fallbackErr := m.secondary.Set(ctx, key, value, ttl); fallbackErr != nil {
				m.mu.Lock()
				m.stats.SecondaryErrors++
				m.mu.Unlock()
				m.logger.Error("Failed to set secondary cache", "error", fallbackErr, "key", key)
				return fmt.Errorf("both primary and secondary cache failed: primary=%w, secondary=%v", err, fallbackErr)
			}
			m.mu.Lock()
			m.stats.FallbackCount++
			m.mu.Unlock()
			m.logger.Warn("Used secondary cache as fallback for set", "key", key)
		}
		return err
	}

	m.mu.Lock()
	m.stats.PrimaryHits++
	m.mu.Unlock()

	// 如果配置了同步到备用缓存
	if m.config.SyncToSecondary && m.secondary != nil {
		go func() {
			if syncErr := m.secondary.Set(context.Background(), key, value, ttl); syncErr != nil {
				m.logger.Error("Failed to sync to secondary cache", "error", syncErr, "key", key)
			} else {
				m.mu.Lock()
				m.stats.SyncCount++
				m.mu.Unlock()
			}
		}()
	}

	return nil
}

// Get 获取缓存
func (m *CacheManager) Get(ctx context.Context, key string, dest interface{}) error {
	// 尝试从主缓存获取
	err := m.executeWithRetry(func() error {
		return m.primary.Get(ctx, key, dest)
	})

	if err == nil {
		m.mu.Lock()
		m.stats.PrimaryHits++
		m.mu.Unlock()
		return nil
	}

	m.mu.Lock()
	m.stats.PrimaryErrors++
	m.mu.Unlock()

	// 如果主缓存失败且配置了备用缓存
	if m.config.UseFallback && m.secondary != nil {
		if fallbackErr := m.secondary.Get(ctx, key, dest); fallbackErr == nil {
			m.mu.Lock()
			m.stats.SecondaryHits++
			m.stats.FallbackCount++
			m.mu.Unlock()
			m.logger.Debug("Used secondary cache as fallback for get", "key", key)
			return nil
		} else {
			m.mu.Lock()
			m.stats.SecondaryErrors++
			m.mu.Unlock()
		}
	}

	m.logger.Debug("Cache miss for key", "key", key, "primary_error", err)
	return err
}

// Delete 删除缓存
func (m *CacheManager) Delete(ctx context.Context, key string) error {
	var primaryErr, secondaryErr error

	// 删除主缓存
	primaryErr = m.executeWithRetry(func() error {
		return m.primary.Delete(ctx, key)
	})

	// 删除备用缓存
	if m.secondary != nil {
		secondaryErr = m.secondary.Delete(ctx, key)
	}

	if primaryErr != nil {
		m.mu.Lock()
		m.stats.PrimaryErrors++
		m.mu.Unlock()
		m.logger.Error("Failed to delete from primary cache", "error", primaryErr, "key", key)
	}

	if secondaryErr != nil {
		m.mu.Lock()
		m.stats.SecondaryErrors++
		m.mu.Unlock()
		m.logger.Error("Failed to delete from secondary cache", "error", secondaryErr, "key", key)
	}

	// 如果主缓存删除成功，认为操作成功
	if primaryErr == nil {
		return nil
	}

	// 如果主缓存失败但备用缓存成功，也认为操作成功
	if secondaryErr == nil {
		return nil
	}

	return primaryErr
}

// Exists 检查缓存是否存在
func (m *CacheManager) Exists(ctx context.Context, key string) (bool, error) {
	// 检查主缓存
	exists, err := m.primary.Exists(ctx, key)
	if err == nil {
		return exists, nil
	}

	m.mu.Lock()
	m.stats.PrimaryErrors++
	m.mu.Unlock()

	// 如果主缓存失败，检查备用缓存
	if m.config.UseFallback && m.secondary != nil {
		if fallbackExists, fallbackErr := m.secondary.Exists(ctx, key); fallbackErr == nil {
			m.mu.Lock()
			m.stats.FallbackCount++
			m.mu.Unlock()
			return fallbackExists, nil
		} else {
			m.mu.Lock()
			m.stats.SecondaryErrors++
			m.mu.Unlock()
		}
	}

	return false, err
}

// SetBatch 批量设置缓存
func (m *CacheManager) SetBatch(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
	// 尝试批量设置主缓存
	err := m.executeWithRetry(func() error {
		return m.primary.SetBatch(ctx, items, ttl)
	})

	if err != nil {
		m.mu.Lock()
		m.stats.PrimaryErrors++
		m.mu.Unlock()

		// 如果配置了错误时回退，尝试设置备用缓存
		if m.config.FallbackOnError && m.secondary != nil {
			if fallbackErr := m.secondary.SetBatch(ctx, items, ttl); fallbackErr != nil {
				m.mu.Lock()
				m.stats.SecondaryErrors++
				m.mu.Unlock()
				return fmt.Errorf("both primary and secondary cache failed: primary=%w, secondary=%v", err, fallbackErr)
			}
			m.mu.Lock()
			m.stats.FallbackCount++
			m.mu.Unlock()
		}
		return err
	}

	// 如果配置了同步到备用缓存
	if m.config.SyncToSecondary && m.secondary != nil {
		go func() {
			if syncErr := m.secondary.SetBatch(context.Background(), items, ttl); syncErr != nil {
				m.logger.Error("Failed to sync batch to secondary cache", "error", syncErr, "count", len(items))
			} else {
				m.mu.Lock()
				m.stats.SyncCount++
				m.mu.Unlock()
			}
		}()
	}

	return nil
}

// GetBatch 批量获取缓存
func (m *CacheManager) GetBatch(ctx context.Context, keys []string) (map[string]interface{}, error) {
	// 尝试从主缓存批量获取
	results, err := m.primary.GetBatch(ctx, keys)
	if err == nil {
		m.mu.Lock()
		m.stats.PrimaryHits++
		m.mu.Unlock()
		return results, nil
	}

	m.mu.Lock()
	m.stats.PrimaryErrors++
	m.mu.Unlock()

	// 如果主缓存失败且配置了备用缓存
	if m.config.UseFallback && m.secondary != nil {
		if fallbackResults, fallbackErr := m.secondary.GetBatch(ctx, keys); fallbackErr == nil {
			m.mu.Lock()
			m.stats.SecondaryHits++
			m.stats.FallbackCount++
			m.mu.Unlock()
			return fallbackResults, nil
		} else {
			m.mu.Lock()
			m.stats.SecondaryErrors++
			m.mu.Unlock()
		}
	}

	return nil, err
}

// DeleteBatch 批量删除缓存
func (m *CacheManager) DeleteBatch(ctx context.Context, keys []string) error {
	var primaryErr, secondaryErr error

	// 删除主缓存
	primaryErr = m.executeWithRetry(func() error {
		return m.primary.DeleteBatch(ctx, keys)
	})

	// 删除备用缓存
	if m.secondary != nil {
		secondaryErr = m.secondary.DeleteBatch(ctx, keys)
	}

	if primaryErr != nil {
		m.mu.Lock()
		m.stats.PrimaryErrors++
		m.mu.Unlock()
	}

	if secondaryErr != nil {
		m.mu.Lock()
		m.stats.SecondaryErrors++
		m.mu.Unlock()
	}

	// 如果主缓存删除成功，认为操作成功
	if primaryErr == nil {
		return nil
	}

	return primaryErr
}

// SetTTL 设置缓存过期时间
func (m *CacheManager) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	err := m.executeWithRetry(func() error {
		return m.primary.SetTTL(ctx, key, ttl)
	})

	if err != nil && m.config.FallbackOnError && m.secondary != nil {
		return m.secondary.SetTTL(ctx, key, ttl)
	}

	return err
}

// GetTTL 获取缓存剩余过期时间
func (m *CacheManager) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := m.primary.GetTTL(ctx, key)
	if err == nil {
		return ttl, nil
	}

	if m.config.UseFallback && m.secondary != nil {
		return m.secondary.GetTTL(ctx, key)
	}

	return 0, err
}

// Clear 清理缓存
func (m *CacheManager) Clear(ctx context.Context) error {
	var primaryErr, secondaryErr error

	primaryErr = m.primary.Clear(ctx)
	if m.secondary != nil {
		secondaryErr = m.secondary.Clear(ctx)
	}

	if primaryErr != nil {
		return primaryErr
	}

	return secondaryErr
}

// FlushDB 清空整个缓存
func (m *CacheManager) FlushDB(ctx context.Context) error {
	var primaryErr, secondaryErr error

	primaryErr = m.primary.FlushDB(ctx)
	if m.secondary != nil {
		secondaryErr = m.secondary.FlushDB(ctx)
	}

	if primaryErr != nil {
		return primaryErr
	}

	return secondaryErr
}

// Ping 健康检查
func (m *CacheManager) Ping(ctx context.Context) error {
	err := m.primary.Ping(ctx)
	if err != nil && m.secondary != nil {
		return m.secondary.Ping(ctx)
	}
	return err
}

// Close 关闭缓存管理器
func (m *CacheManager) Close() error {
	var primaryErr, secondaryErr error

	if m.primary != nil {
		primaryErr = m.primary.Close()
	}

	if m.secondary != nil {
		secondaryErr = m.secondary.Close()
	}

	m.logger.Info("Cache manager closed successfully")

	if primaryErr != nil {
		return primaryErr
	}

	return secondaryErr
}

// GetStats 获取管理器统计信息
func (m *CacheManager) GetStats() *ManagerStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 创建副本以避免并发访问问题
	return &ManagerStats{
		PrimaryHits:     m.stats.PrimaryHits,
		SecondaryHits:   m.stats.SecondaryHits,
		PrimaryErrors:   m.stats.PrimaryErrors,
		SecondaryErrors: m.stats.SecondaryErrors,
		FallbackCount:   m.stats.FallbackCount,
		SyncCount:       m.stats.SyncCount,
	}
}

// GetPrimaryCache 获取主缓存
func (m *CacheManager) GetPrimaryCache() Cache {
	return m.primary
}

// GetSecondaryCache 获取备用缓存
func (m *CacheManager) GetSecondaryCache() Cache {
	return m.secondary
}

// 私有方法

// executeWithRetry 带重试的执行
func (m *CacheManager) executeWithRetry(fn func() error) error {
	var lastErr error

	for i := 0; i <= m.config.RetryAttempts; i++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
			if i < m.config.RetryAttempts {
				time.Sleep(m.config.RetryDelay)
			}
		}
	}

	return lastErr
}

// healthCheck 健康检查
func (m *CacheManager) healthCheck() {
	ticker := time.NewTicker(m.config.HealthCheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		// 检查主缓存
		if err := m.primary.Ping(ctx); err != nil {
			m.logger.Error("Primary cache health check failed", "error", err)
		} else {
			m.logger.Debug("Primary cache health check passed")
		}

		// 检查备用缓存
		if m.secondary != nil {
			if err := m.secondary.Ping(ctx); err != nil {
				m.logger.Error("Secondary cache health check failed", "error", err)
			} else {
				m.logger.Debug("Secondary cache health check passed")
			}
		}

		cancel()
	}
}
