package network

import (
	"context"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logging"
)

// ConnectionManager 连接管理器
type ConnectionManager struct {
	connections map[string]*Connection
	mutex       sync.RWMutex
	logger      logging.Logger
}

// Connection 连接信息
type Connection struct {
	ID        string
	Address   string
	CreatedAt time.Time
	LastSeen  time.Time
	Status    string
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager(logger logging.Logger) *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*Connection),
		logger:      logger,
	}
}

// AddConnection 添加连接
func (cm *ConnectionManager) AddConnection(id, address string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.connections[id] = &Connection{
		ID:        id,
		Address:   address,
		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
		Status:    "active",
	}

	cm.logger.Info("连接已添加", map[string]interface{}{
		"connection_id": id,
		"address":       address,
	})
}

// RemoveConnection 移除连接
func (cm *ConnectionManager) RemoveConnection(id string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if conn, exists := cm.connections[id]; exists {
		delete(cm.connections, id)
		cm.logger.Info("连接已移除", map[string]interface{}{
			"connection_id": id,
			"address":       conn.Address,
		})
	}
}

// GetConnection 获取连接
func (cm *ConnectionManager) GetConnection(id string) (*Connection, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	conn, exists := cm.connections[id]
	return conn, exists
}

// GetAllConnections 获取所有连接
func (cm *ConnectionManager) GetAllConnections() map[string]*Connection {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	// 返回副本
	connections := make(map[string]*Connection)
	for id, conn := range cm.connections {
		connections[id] = conn
	}
	return connections
}

// UpdateLastSeen 更新最后活跃时间
func (cm *ConnectionManager) UpdateLastSeen(id string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if conn, exists := cm.connections[id]; exists {
		conn.LastSeen = time.Now()
	}
}

// CleanupInactiveConnections 清理非活跃连接
func (cm *ConnectionManager) CleanupInactiveConnections(timeout time.Duration) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	now := time.Now()
	for id, conn := range cm.connections {
		if now.Sub(conn.LastSeen) > timeout {
			delete(cm.connections, id)
			cm.logger.Info("清理非活跃连接", map[string]interface{}{
				"connection_id": id,
				"last_seen":     conn.LastSeen,
			})
		}
	}
}

// GetConnectionCount 获取连接数量
func (cm *ConnectionManager) GetConnectionCount() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	return len(cm.connections)
}

// StartCleanupRoutine 启动清理例程
func (cm *ConnectionManager) StartCleanupRoutine(ctx context.Context, interval, timeout time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			cm.logger.Info("连接清理例程已停止")
			return
		case <-ticker.C:
			cm.CleanupInactiveConnections(timeout)
		}
	}
}
