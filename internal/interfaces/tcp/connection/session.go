package connection

import (
	"context"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logger"
	"greatestworks/internal/interfaces/tcp/protocol"
)

// SessionState 会话状态
type SessionState int

const (
	SessionStateNew SessionState = iota
	SessionStateConnected
	SessionStateAuthenticated
	SessionStateActive
	SessionStateIdle
	SessionStateDisconnecting
	SessionStateDisconnected
)

// String 返回会话状态的字符串表示
func (s SessionState) String() string {
	switch s {
	case SessionStateNew:
		return "new"
	case SessionStateConnected:
		return "connected"
	case SessionStateAuthenticated:
		return "authenticated"
	case SessionStateActive:
		return "active"
	case SessionStateIdle:
		return "idle"
	case SessionStateDisconnecting:
		return "disconnecting"
	case SessionStateDisconnected:
		return "disconnected"
	default:
		return "unknown"
	}
}

// Session 会话信息
type Session struct {
	ID            string
	ConnectionID  string
	PlayerID      string
	State         SessionState
	CreatedAt     time.Time
	LastActivity  time.Time
	AuthTime      time.Time
	IdleTimeout   time.Duration
	SessionData   map[string]interface{}
	mutex         sync.RWMutex
	logger        logger.Logger
}

// NewSession 创建新会话
func NewSession(id, connectionID string, logger logger.Logger) *Session {
	now := time.Now()
	return &Session{
		ID:           id,
		ConnectionID: connectionID,
		State:        SessionStateNew,
		CreatedAt:    now,
		LastActivity: now,
		IdleTimeout:  30 * time.Minute, // 默认30分钟空闲超时
		SessionData:  make(map[string]interface{}),
		logger:       logger,
	}
}

// UpdateActivity 更新活动时间
func (s *Session) UpdateActivity() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.LastActivity = time.Now()
	if s.State == SessionStateIdle {
		s.State = SessionStateActive
		s.logger.Debug("Session state changed from idle to active", "session_id", s.ID)
	}
}

// SetState 设置会话状态
func (s *Session) SetState(state SessionState) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	oldState := s.State
	s.State = state

	if state == SessionStateAuthenticated {
		s.AuthTime = time.Now()
	}

	s.logger.Debug("Session state changed", 
		"session_id", s.ID,
		"old_state", oldState.String(),
		"new_state", state.String())
}

// GetState 获取会话状态
func (s *Session) GetState() SessionState {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.State
}

// SetPlayerID 设置玩家ID
func (s *Session) SetPlayerID(playerID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.PlayerID = playerID
	s.logger.Info("Player bound to session", "session_id", s.ID, "player_id", playerID)
}

// GetPlayerID 获取玩家ID
func (s *Session) GetPlayerID() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.PlayerID
}

// IsAuthenticated 检查是否已认证
func (s *Session) IsAuthenticated() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.State >= SessionStateAuthenticated
}

// IsActive 检查是否活跃
func (s *Session) IsActive() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.State == SessionStateActive || s.State == SessionStateAuthenticated
}

// IsIdle 检查是否空闲
func (s *Session) IsIdle() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.State == SessionStateDisconnected || s.State == SessionStateDisconnecting {
		return false
	}

	return time.Since(s.LastActivity) > s.IdleTimeout
}

// GetIdleDuration 获取空闲时长
func (s *Session) GetIdleDuration() time.Duration {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return time.Since(s.LastActivity)
}

// SetSessionData 设置会话数据
func (s *Session) SetSessionData(key string, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.SessionData[key] = value
}

// GetSessionData 获取会话数据
func (s *Session) GetSessionData(key string) (interface{}, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	value, exists := s.SessionData[key]
	return value, exists
}

// RemoveSessionData 移除会话数据
func (s *Session) RemoveSessionData(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.SessionData, key)
}

// GetSessionInfo 获取会话信息
func (s *Session) GetSessionInfo() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return map[string]interface{}{
		"id":             s.ID,
		"connection_id":  s.ConnectionID,
		"player_id":      s.PlayerID,
		"state":          s.State.String(),
		"created_at":     s.CreatedAt.Unix(),
		"last_activity":  s.LastActivity.Unix(),
		"auth_time":      s.AuthTime.Unix(),
		"idle_duration":  s.GetIdleDuration().Seconds(),
		"is_authenticated": s.IsAuthenticated(),
		"is_active":      s.IsActive(),
		"is_idle":        s.IsIdle(),
	}
}

// SessionManager 会话管理器
type SessionManager struct {
	sessions      map[string]*Session
	playerSessions map[string]*Session // 玩家ID到会话的映射
	mutex         sync.RWMutex
	logger        logger.Logger
	cleanupTicker *time.Ticker
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewSessionManager 创建会话管理器
func NewSessionManager(logger logger.Logger) *SessionManager {
	ctx, cancel := context.WithCancel(context.Background())

	sm := &SessionManager{
		sessions:       make(map[string]*Session),
		playerSessions: make(map[string]*Session),
		logger:         logger,
		cleanupTicker:  time.NewTicker(5 * time.Minute), // 每5分钟清理一次
		ctx:            ctx,
		cancel:         cancel,
	}

	// 启动清理协程
	go sm.cleanupRoutine()

	return sm
}

// CreateSession 创建会话
func (sm *SessionManager) CreateSession(sessionID, connectionID string) *Session {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	session := NewSession(sessionID, connectionID, sm.logger)
	sm.sessions[sessionID] = session

	sm.logger.Info("Session created", "session_id", sessionID, "connection_id", connectionID)
	return session
}

// GetSession 获取会话
func (sm *SessionManager) GetSession(sessionID string) (*Session, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	session, exists := sm.sessions[sessionID]
	return session, exists
}

// GetSessionByPlayer 根据玩家ID获取会话
func (sm *SessionManager) GetSessionByPlayer(playerID string) (*Session, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	session, exists := sm.playerSessions[playerID]
	return session, exists
}

// BindPlayerToSession 绑定玩家到会话
func (sm *SessionManager) BindPlayerToSession(sessionID, playerID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	// 检查玩家是否已经绑定到其他会话
	if existingSession, exists := sm.playerSessions[playerID]; exists {
		// 如果是同一个会话，直接返回
		if existingSession.ID == sessionID {
			return nil
		}
		// 断开旧会话
		existingSession.SetState(SessionStateDisconnecting)
		delete(sm.playerSessions, playerID)
		sm.logger.Warn("Player switched sessions", 
			"player_id", playerID,
			"old_session", existingSession.ID,
			"new_session", sessionID)
	}

	session.SetPlayerID(playerID)
	sm.playerSessions[playerID] = session

	sm.logger.Info("Player bound to session", "player_id", playerID, "session_id", sessionID)
	return nil
}

// UnbindPlayerFromSession 解绑玩家和会话
func (sm *SessionManager) UnbindPlayerFromSession(playerID string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if session, exists := sm.playerSessions[playerID]; exists {
		delete(sm.playerSessions, playerID)
		sm.logger.Info("Player unbound from session", "player_id", playerID, "session_id", session.ID)
	}
}

// RemoveSession 移除会话
func (sm *SessionManager) RemoveSession(sessionID string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if session, exists := sm.sessions[sessionID]; exists {
		// 解绑玩家
		if session.PlayerID != "" {
			delete(sm.playerSessions, session.PlayerID)
		}
		
		// 设置会话状态为已断开
		session.SetState(SessionStateDisconnected)
		
		// 移除会话
		delete(sm.sessions, sessionID)
		
		sm.logger.Info("Session removed", "session_id", sessionID, "player_id", session.PlayerID)
	}
}

// GetAllSessions 获取所有会话
func (sm *SessionManager) GetAllSessions() []*Session {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	sessions := make([]*Session, 0, len(sm.sessions))
	for _, session := range sm.sessions {
		sessions = append(sessions, session)
	}

	return sessions
}

// GetActiveSessions 获取活跃会话
func (sm *SessionManager) GetActiveSessions() []*Session {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	var activeSessions []*Session
	for _, session := range sm.sessions {
		if session.IsActive() {
			activeSessions = append(activeSessions, session)
		}
	}

	return activeSessions
}

// GetSessionCount 获取会话数量
func (sm *SessionManager) GetSessionCount() int {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	return len(sm.sessions)
}

// GetActiveSessionCount 获取活跃会话数量
func (sm *SessionManager) GetActiveSessionCount() int {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	count := 0
	for _, session := range sm.sessions {
		if session.IsActive() {
			count++
		}
	}

	return count
}

// cleanupRoutine 清理协程
func (sm *SessionManager) cleanupRoutine() {
	for {
		select {
		case <-sm.ctx.Done():
			return
		case <-sm.cleanupTicker.C:
			sm.cleanupIdleSessions()
		}
	}
}

// cleanupIdleSessions 清理空闲会话
func (sm *SessionManager) cleanupIdleSessions() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	var toRemove []string
	for sessionID, session := range sm.sessions {
		if session.IsIdle() || session.GetState() == SessionStateDisconnected {
			toRemove = append(toRemove, sessionID)
		}
	}

	for _, sessionID := range toRemove {
		if session, exists := sm.sessions[sessionID]; exists {
			// 解绑玩家
			if session.PlayerID != "" {
				delete(sm.playerSessions, session.PlayerID)
			}
			
			// 移除会话
			delete(sm.sessions, sessionID)
			
			sm.logger.Info("Idle session cleaned up", 
				"session_id", sessionID,
				"player_id", session.PlayerID,
				"idle_duration", session.GetIdleDuration().String())
		}
	}

	if len(toRemove) > 0 {
		sm.logger.Info("Session cleanup completed", "removed_count", len(toRemove))
	}
}

// Stop 停止会话管理器
func (sm *SessionManager) Stop() {
	sm.cancel()
	sm.cleanupTicker.Stop()
	sm.logger.Info("Session manager stopped")
}

// GetStats 获取统计信息
func (sm *SessionManager) GetStats() map[string]interface{} {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	totalSessions := len(sm.sessions)
	activeSessions := 0
	idleSessions := 0
	authenticatedSessions := 0

	for _, session := range sm.sessions {
		if session.IsActive() {
			activeSessions++
		}
		if session.IsIdle() {
			idleSessions++
		}
		if session.IsAuthenticated() {
			authenticatedSessions++
		}
	}

	return map[string]interface{}{
		"total_sessions":        totalSessions,
		"active_sessions":       activeSessions,
		"idle_sessions":         idleSessions,
		"authenticated_sessions": authenticatedSessions,
		"bound_players":         len(sm.playerSessions),
	}
}