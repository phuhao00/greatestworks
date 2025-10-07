package connection

import (
	"context"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logging"
)

// Manager 连接管理器
type Manager struct {
	connections map[string]*Session
	mutex       sync.RWMutex
	logger      logging.Logger
}

// NewManager 创建连接管理器
func NewManager(logger logging.Logger) *Manager {
	return &Manager{
		connections: make(map[string]*Session),
		logger:      logger,
	}
}

// AddConnection 添加连接
func (m *Manager) AddConnection(session *Session) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.connections[session.ID] = session
	m.logger.Info("Connection added", logging.Fields{
		"session_id": session.ID,
		"address":    session.RemoteAddr,
	})
}

// RemoveConnection 移除连接
func (m *Manager) RemoveConnection(sessionID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if session, exists := m.connections[sessionID]; exists {
		delete(m.connections, sessionID)
		m.logger.Info("Connection removed", logging.Fields{
			"session_id": sessionID,
			"address":    session.RemoteAddr,
		})
	}
}

// GetConnection 获取连接
func (m *Manager) GetConnection(sessionID string) (*Session, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	session, exists := m.connections[sessionID]
	return session, exists
}

// GetAllConnections 获取所有连接
func (m *Manager) GetAllConnections() map[string]*Session {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// 返回副本
	connections := make(map[string]*Session)
	for id, session := range m.connections {
		connections[id] = session
	}
	return connections
}

// GetConnectionCount 获取连接数量
func (m *Manager) GetConnectionCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.connections)
}

// Broadcast 广播消息
func (m *Manager) Broadcast(message []byte) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, session := range m.connections {
		if err := session.Send(message); err != nil {
			m.logger.Error("Failed to broadcast message", err, logging.Fields{
				"session_id": session.ID,
			})
		}
	}
}

// BroadcastToGroup 向指定组广播消息
func (m *Manager) BroadcastToGroup(groupID string, message []byte) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, session := range m.connections {
		if session.GroupID == groupID {
			if err := session.Send(message); err != nil {
				m.logger.Error("Failed to broadcast to group", err, logging.Fields{
					"session_id": session.ID,
					"group_id":   groupID,
				})
			}
		}
	}
}

// CleanupInactiveConnections 清理非活跃连接
func (m *Manager) CleanupInactiveConnections(timeout time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	for id, session := range m.connections {
		if now.Sub(session.LastActivity) > timeout {
			session.Close()
			delete(m.connections, id)
			m.logger.Info("Cleaned up inactive connection", logging.Fields{
				"session_id":    id,
				"last_activity": session.LastActivity,
			})
		}
	}
}

// StartCleanupRoutine 启动清理例程
func (m *Manager) StartCleanupRoutine(ctx context.Context, interval, timeout time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			m.logger.Info("Connection cleanup routine stopped", logging.Fields{})
			return
		case <-ticker.C:
			m.CleanupInactiveConnections(timeout)
		}
	}
}
