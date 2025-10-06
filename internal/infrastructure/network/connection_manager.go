package network

import (
	"context"
	"fmt"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logger"
)

// ConnectionManager 连接管理器
type ConnectionManager struct {
	connections map[string]*ManagedConnection
	mu          sync.RWMutex
	logger      logger.Logger
	config      *ConnectionManagerConfig
	stats       *ConnectionManagerStats
	ctx         context.Context
	cancel      context.CancelFunc
	pools       map[string]*ConnectionPool
}

// ConnectionManagerConfig 连接管理器配置
type ConnectionManagerConfig struct {
	MaxConnections    int           `json:"max_connections" yaml:"max_connections"`
	ConnectionTimeout time.Duration `json:"connection_timeout" yaml:"connection_timeout"`
	IdleTimeout       time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
	CleanupInterval   time.Duration `json:"cleanup_interval" yaml:"cleanup_interval"`
	HeartbeatInterval time.Duration `json:"heartbeat_interval" yaml:"heartbeat_interval"`
	EnableMetrics     bool          `json:"enable_metrics" yaml:"enable_metrics"`
	EnablePooling     bool          `json:"enable_pooling" yaml:"enable_pooling"`
	PoolSize          int           `json:"pool_size" yaml:"pool_size"`
	PoolMaxIdle       int           `json:"pool_max_idle" yaml:"pool_max_idle"`
	PoolMaxLifetime   time.Duration `json:"pool_max_lifetime" yaml:"pool_max_lifetime"`
}

// ManagedConnection 管理的连接
type ManagedConnection struct {
	ID         string
	Connection Connection
	CreatedAt  time.Time
	LastUsedAt time.Time
	IsActive   bool
	Stats      *ConnectionStats
}

// ConnectionStats 连接统计信息
type ConnectionStats struct {
	MessagesSent     int64         `json:"messages_sent"`
	MessagesReceived int64         `json:"messages_received"`
	BytesSent        int64         `json:"bytes_sent"`
	BytesReceived    int64         `json:"bytes_received"`
	Errors           int64         `json:"errors"`
	LastActivity     time.Time     `json:"last_activity"`
	Uptime           time.Duration `json:"uptime"`
}

// ConnectionPool 连接池
type ConnectionPool struct {
	ID          string
	Connections []*ManagedConnection
	MaxSize     int
	MinSize     int
	mu          sync.RWMutex
}

// ConnectionManagerStats 连接管理器统计信息
type ConnectionManagerStats struct {
	TotalConnections  int64            `json:"total_connections"`
	ActiveConnections int64            `json:"active_connections"`
	IdleConnections   int64            `json:"idle_connections"`
	TotalMessages     int64            `json:"total_messages"`
	TotalErrors       int64            `json:"total_errors"`
	StartTime         time.Time        `json:"start_time"`
	Uptime            time.Duration    `json:"uptime"`
	ByConnectionType  map[string]int64 `json:"by_connection_type"`
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager(config *ConnectionManagerConfig, logger logger.Logger) *ConnectionManager {
	if config == nil {
		config = &ConnectionManagerConfig{
			MaxConnections:    1000,
			ConnectionTimeout: 30 * time.Second,
			IdleTimeout:       5 * time.Minute,
			CleanupInterval:   1 * time.Minute,
			HeartbeatInterval: 30 * time.Second,
			EnableMetrics:     true,
			EnablePooling:     false,
			PoolSize:          10,
			PoolMaxIdle:       5,
			PoolMaxLifetime:   1 * time.Hour,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	manager := &ConnectionManager{
		connections: make(map[string]*ManagedConnection),
		logger:      logger,
		config:      config,
		ctx:         ctx,
		cancel:      cancel,
		pools:       make(map[string]*ConnectionPool),
		stats: &ConnectionManagerStats{
			StartTime:        time.Now(),
			ByConnectionType: make(map[string]int64),
		},
	}

	// 启动清理协程
	if config.CleanupInterval > 0 {
		go manager.startCleanup()
	}

	// 启动心跳检测
	if config.HeartbeatInterval > 0 {
		go manager.startHeartbeat()
	}

	logger.Info("Connection manager initialized successfully")
	return manager
}

// AddConnection 添加连接
func (cm *ConnectionManager) AddConnection(id string, conn Connection) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.connections[id]; exists {
		return fmt.Errorf("connection %s already exists", id)
	}

	managedConn := &ManagedConnection{
		ID:         id,
		Connection: conn,
		CreatedAt:  time.Now(),
		LastUsedAt: time.Now(),
		IsActive:   true,
		Stats: &ConnectionStats{
			LastActivity: time.Now(),
		},
	}

	cm.connections[id] = managedConn
	cm.stats.TotalConnections++
	cm.stats.ActiveConnections++

	cm.logger.Info("Connection added successfully", "connection_id", id)
	return nil
}

// RemoveConnection 移除连接
func (cm *ConnectionManager) RemoveConnection(id string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	conn, exists := cm.connections[id]
	if !exists {
		return fmt.Errorf("connection %s not found", id)
	}

	conn.IsActive = false
	delete(cm.connections, id)
	cm.stats.ActiveConnections--

	cm.logger.Info("Connection removed successfully", "connection_id", id)
	return nil
}

// GetConnection 获取连接
func (cm *ConnectionManager) GetConnection(id string) *ManagedConnection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return cm.connections[id]
}

// GetAllConnections 获取所有连接
func (cm *ConnectionManager) GetAllConnections() map[string]*ManagedConnection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// 创建副本
	connections := make(map[string]*ManagedConnection)
	for id, conn := range cm.connections {
		connections[id] = conn
	}

	return connections
}

// GetActiveConnections 获取活跃连接
func (cm *ConnectionManager) GetActiveConnections() map[string]*ManagedConnection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	activeConnections := make(map[string]*ManagedConnection)
	for id, conn := range cm.connections {
		if conn.IsActive {
			activeConnections[id] = conn
		}
	}

	return activeConnections
}

// Broadcast 广播消息
func (cm *ConnectionManager) Broadcast(data []byte) error {
	cm.mu.RLock()
	connections := make([]*ManagedConnection, 0, len(cm.connections))
	for _, conn := range cm.connections {
		if conn.IsActive {
			connections = append(connections, conn)
		}
	}
	cm.mu.RUnlock()

	if len(connections) == 0 {
		cm.logger.Debug("No active connections to broadcast to")
		return nil
	}

	var errors []error
	successCount := 0

	for _, conn := range connections {
		if err := conn.Connection.Send(data); err != nil {
			cm.logger.Error("Failed to broadcast to connection", "error", err, "connection_id", conn.ID)
			errors = append(errors, err)
			conn.Stats.Errors++
		} else {
			successCount++
			conn.Stats.MessagesSent++
			conn.Stats.BytesSent += int64(len(data))
			conn.Stats.LastActivity = time.Now()
			conn.LastUsedAt = time.Now()
		}
	}

	cm.stats.TotalMessages += int64(successCount)
	if len(errors) > 0 {
		cm.stats.TotalErrors += int64(len(errors))
	}

	cm.logger.Debug("Broadcast completed", "total_connections", len(connections), "success_count", successCount, "error_count", len(errors))

	if len(errors) > 0 {
		return fmt.Errorf("broadcast failed for %d connections: %v", len(errors), errors[0])
	}

	return nil
}

// SendToConnection 发送消息到指定连接
func (cm *ConnectionManager) SendToConnection(id string, data []byte) error {
	conn := cm.GetConnection(id)
	if conn == nil {
		return fmt.Errorf("connection %s not found", id)
	}

	if !conn.IsActive {
		return fmt.Errorf("connection %s is not active", id)
	}

	err := conn.Connection.Send(data)
	if err != nil {
		conn.Stats.Errors++
		cm.stats.TotalErrors++
		cm.logger.Error("Failed to send message to connection", "error", err, "connection_id", id)
		return fmt.Errorf("failed to send message to connection %s: %w", id, err)
	}

	conn.Stats.MessagesSent++
	conn.Stats.BytesSent += int64(len(data))
	conn.Stats.LastActivity = time.Now()
	conn.LastUsedAt = time.Now()
	cm.stats.TotalMessages++

	cm.logger.Debug("Message sent to connection successfully", "connection_id", id)
	return nil
}

// GetStats 获取统计信息
func (cm *ConnectionManager) GetStats() *ConnectionManagerStats {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// 计算活跃连接数
	activeCount := int64(0)
	for _, conn := range cm.connections {
		if conn.IsActive {
			activeCount++
		}
	}

	// 创建统计信息副本
	stats := &ConnectionManagerStats{
		TotalConnections:  cm.stats.TotalConnections,
		ActiveConnections: activeCount,
		IdleConnections:   cm.stats.TotalConnections - activeCount,
		TotalMessages:     cm.stats.TotalMessages,
		TotalErrors:       cm.stats.TotalErrors,
		StartTime:         cm.stats.StartTime,
		Uptime:            time.Since(cm.stats.StartTime),
		ByConnectionType:  make(map[string]int64),
	}

	// 复制连接类型统计
	for connType, count := range cm.stats.ByConnectionType {
		stats.ByConnectionType[connType] = count
	}

	return stats
}

// Close 关闭连接管理器
func (cm *ConnectionManager) Close() error {
	cm.logger.Info("Closing connection manager")

	// 取消上下文
	cm.cancel()

	// 关闭所有连接
	cm.mu.Lock()
	for _, conn := range cm.connections {
		if conn.IsActive {
			conn.Connection.Close()
			conn.IsActive = false
		}
	}
	cm.mu.Unlock()

	cm.logger.Info("Connection manager closed successfully")
	return nil
}

// 私有方法

// startCleanup 启动清理协程
func (cm *ConnectionManager) startCleanup() {
	ticker := time.NewTicker(cm.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cm.cleanupIdleConnections()
		case <-cm.ctx.Done():
			return
		}
	}
}

// cleanupIdleConnections 清理空闲连接
func (cm *ConnectionManager) cleanupIdleConnections() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	now := time.Now()
	cleanedCount := 0

	for connID, conn := range cm.connections {
		if conn.IsActive && now.Sub(conn.LastUsedAt) > cm.config.IdleTimeout {
			conn.Connection.Close()
			conn.IsActive = false
			cleanedCount++
			cm.logger.Debug("Cleaned up idle connection", "connection_id", connID, "idle_time", now.Sub(conn.LastUsedAt))
		}
	}

	if cleanedCount > 0 {
		cm.stats.ActiveConnections -= int64(cleanedCount)
		cm.logger.Info("Cleaned up idle connections", "count", cleanedCount)
	}
}

// startHeartbeat 启动心跳检测
func (cm *ConnectionManager) startHeartbeat() {
	ticker := time.NewTicker(cm.config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cm.sendHeartbeat()
		case <-cm.ctx.Done():
			return
		}
	}
}

// sendHeartbeat 发送心跳
func (cm *ConnectionManager) sendHeartbeat() {
	heartbeatData := []byte("heartbeat")

	cm.mu.RLock()
	connections := make([]*ManagedConnection, 0, len(cm.connections))
	for _, conn := range cm.connections {
		if conn.IsActive {
			connections = append(connections, conn)
		}
	}
	cm.mu.RUnlock()

	for _, conn := range connections {
		if err := conn.Connection.Send(heartbeatData); err != nil {
			cm.logger.Debug("Failed to send heartbeat", "connection_id", conn.ID, "error", err)
		}
	}

	cm.logger.Debug("Heartbeat sent to all connections", "connection_count", len(connections))
}
