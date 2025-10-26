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
	// Mapping from player entity ID to session
	playerSessions map[int32]*Session
	// Reverse mapping from session ID to player entity ID
	sessionToPlayer map[string]int32
	mutex           sync.RWMutex
	logger          logging.Logger
}

// NewManager 创建连接管理器
func NewManager(logger logging.Logger) *Manager {
	return &Manager{
		connections:     make(map[string]*Session),
		playerSessions:  make(map[int32]*Session),
		sessionToPlayer: make(map[string]int32),
		logger:          logger,
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
		// Also remove any player-session bindings pointing to this session
		for pid, s := range m.playerSessions {
			if s == session {
				delete(m.playerSessions, pid)
			}
		}
		delete(m.sessionToPlayer, sessionID)
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

// BindPlayer binds a player entity ID to a session for targeted sends.
func (m *Manager) BindPlayer(entityID int32, session *Session) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.playerSessions[entityID] = session
	m.sessionToPlayer[session.ID] = entityID
	m.logger.Info("Player bound to session", logging.Fields{
		"entity_id":  entityID,
		"session_id": session.ID,
	})
}

// UnbindPlayer removes the binding between a player entity ID and any session.
func (m *Manager) UnbindPlayer(entityID int32) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if s, ok := m.playerSessions[entityID]; ok {
		delete(m.sessionToPlayer, s.ID)
		delete(m.playerSessions, entityID)
	} else {
		delete(m.playerSessions, entityID)
	}
	m.logger.Info("Player unbound from session", logging.Fields{
		"entity_id": entityID,
	})
}

// GetSessionByPlayer retrieves the session bound to the given player entity ID.
func (m *Manager) GetSessionByPlayer(entityID int32) (*Session, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	s, ok := m.playerSessions[entityID]
	return s, ok
}

// GetPlayerBySession retrieves the bound player entity ID from a session ID.
func (m *Manager) GetPlayerBySession(sessionID string) (int32, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	pid, ok := m.sessionToPlayer[sessionID]
	return pid, ok
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
