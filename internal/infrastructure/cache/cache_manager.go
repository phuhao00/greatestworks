package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logging"
)

// Cache 缓存接口
type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Clear(ctx context.Context) error
}

// CacheManager 缓存管理器
type CacheManager struct {
	primary   Cache
	secondary Cache
	logger    logging.Logger
	config    *CacheManagerConfig
	mu        sync.RWMutex
	stats     *ManagerStats
}

// CacheManagerConfig 缓存管理器配置
type CacheManagerConfig struct {
	UseFallback             bool          `json:"use_fallback"`
	FallbackOnError         bool          `json:"fallback_on_error"`
	SyncToSecondary         bool          `json:"sync_to_secondary"`
	SyncInterval            time.Duration `json:"sync_interval"`
	MaxRetries              int           `json:"max_retries"`
	RetryDelay              time.Duration `json:"retry_delay"`
	CircuitBreakerThreshold int           `json:"circuit_breaker_threshold"`
}

// ManagerStats 管理器统计
type ManagerStats struct {
	TotalRequests int64 `json:"total_requests"`
	SuccessCount  int64 `json:"success_count"`
	ErrorCount    int64 `json:"error_count"`
	FallbackCount int64 `json:"fallback_count"`
	SyncCount     int64 `json:"sync_count"`
}

// NewCacheManager 创建缓存管理器
func NewCacheManager(primary, secondary Cache, config *CacheManagerConfig, logger logging.Logger) *CacheManager {
	if config == nil {
		config = &CacheManagerConfig{
			UseFallback:             true,
			FallbackOnError:         true,
			SyncToSecondary:         false,
			SyncInterval:            5 * time.Minute,
			MaxRetries:              3,
			RetryDelay:              time.Second,
			CircuitBreakerThreshold: 10,
		}
	}

	return &CacheManager{
		primary:   primary,
		secondary: secondary,
		logger:    logger,
		config:    config,
		stats:     &ManagerStats{},
	}
}

// Set 设置值
func (m *CacheManager) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stats.TotalRequests++

	// 尝试设置主缓存
	err := m.primary.Set(ctx, key, value, ttl)
	if err != nil {
		m.stats.ErrorCount++
		m.logger.Error("Primary cache set failed", err, logging.Fields{
			"key": key,
		})

		// 如果启用备用缓存，尝试设置备用缓存
		if m.config.UseFallback && m.secondary != nil {
			if err := m.secondary.Set(ctx, key, value, ttl); err != nil {
				m.logger.Error("Secondary cache set failed", err, logging.Fields{
					"key": key,
				})
				return fmt.Errorf("主缓存和备用缓存都设置失败: %w", err)
			}
			m.stats.FallbackCount++
			m.logger.Info("Using secondary cache", logging.Fields{
				"key": key,
			})
		} else {
			return err
		}
	} else {
		m.stats.SuccessCount++

		// 如果启用同步到备用缓存
		if m.config.SyncToSecondary && m.secondary != nil {
			go func() {
				if err := m.secondary.Set(context.Background(), key, value, ttl); err != nil {
					m.logger.Error("Failed to sync to secondary cache", err, logging.Fields{
						"key": key,
					})
				} else {
					m.stats.SyncCount++
				}
			}()
		}
	}

	return nil
}

// Get 获取值
func (m *CacheManager) Get(ctx context.Context, key string, dest interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stats.TotalRequests++

	// 尝试从主缓存获取
	err := m.primary.Get(ctx, key, dest)
	if err != nil {
		m.stats.ErrorCount++
		m.logger.Error("Primary cache get failed", err, logging.Fields{
			"key": key,
		})

		// 如果启用备用缓存，尝试从备用缓存获取
		if m.config.UseFallback && m.secondary != nil {
			if err := m.secondary.Get(ctx, key, dest); err != nil {
				m.logger.Error("Secondary cache get failed", err, logging.Fields{
					"key": key,
				})
				return fmt.Errorf("主缓存和备用缓存都获取失败: %w", err)
			}
			m.stats.FallbackCount++
			m.logger.Info("Got from secondary cache", logging.Fields{
				"key": key,
			})
		} else {
			return err
		}
	} else {
		m.stats.SuccessCount++
	}

	return nil
}

// Delete 删除值
func (m *CacheManager) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stats.TotalRequests++

	// 删除主缓存
	err := m.primary.Delete(ctx, key)
	if err != nil {
		m.stats.ErrorCount++
		m.logger.Error("Primary cache delete failed", err, logging.Fields{
			"key": key,
		})
	}

	// 删除备用缓存
	if m.secondary != nil {
		if err := m.secondary.Delete(ctx, key); err != nil {
			m.logger.Error("Secondary cache delete failed", err, logging.Fields{
				"key": key,
			})
		}
	}

	if err != nil {
		return err
	}

	m.stats.SuccessCount++
	return nil
}

// Exists 检查键是否存在
func (m *CacheManager) Exists(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 检查主缓存
	exists, err := m.primary.Exists(ctx, key)
	if err != nil {
		m.logger.Error("Primary cache exists check failed", err, logging.Fields{
			"key": key,
		})

		// 如果启用备用缓存，检查备用缓存
		if m.config.UseFallback && m.secondary != nil {
			exists, err = m.secondary.Exists(ctx, key)
			if err != nil {
				m.logger.Error("Secondary cache exists check failed", err, logging.Fields{
					"key": key,
				})
				return false, err
			}
		} else {
			return false, err
		}
	}

	return exists, nil
}

// Clear 清空缓存
func (m *CacheManager) Clear(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 清空主缓存
	err := m.primary.Clear(ctx)
	if err != nil {
		m.logger.Error("Primary cache clear failed", err)
	}

	// 清空备用缓存
	if m.secondary != nil {
		if err := m.secondary.Clear(ctx); err != nil {
			m.logger.Error("Secondary cache clear failed", err)
		}
	}

	return err
}

// GetStats 获取统计信息
func (m *CacheManager) GetStats() *ManagerStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 返回统计信息的副本
	return &ManagerStats{
		TotalRequests: m.stats.TotalRequests,
		SuccessCount:  m.stats.SuccessCount,
		ErrorCount:    m.stats.ErrorCount,
		FallbackCount: m.stats.FallbackCount,
		SyncCount:     m.stats.SyncCount,
	}
}

// ResetStats 重置统计信息
func (m *CacheManager) ResetStats() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stats = &ManagerStats{}
}

// GetConfig 获取配置
func (m *CacheManager) GetConfig() *CacheManagerConfig {
	return m.config
}

// UpdateConfig 更新配置
func (m *CacheManager) UpdateConfig(config *CacheManagerConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config = config
}
