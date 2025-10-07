package connection

import (
	"fmt"
	"net"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logging"
)

// Session 会话
type Session struct {
	ID           string
	Conn         net.Conn
	RemoteAddr   string
	GroupID      string
	UserID       string
	CreatedAt    time.Time
	LastActivity time.Time
	Status       string
	mutex        sync.RWMutex
	logger       logging.Logger
}

// NewSession 创建会话
func NewSession(id string, conn net.Conn, logger logging.Logger) *Session {
	return &Session{
		ID:           id,
		Conn:         conn,
		RemoteAddr:   conn.RemoteAddr().String(),
		CreatedAt:    time.Now(),
		LastActivity: time.Now(),
		Status:       "active",
		logger:       logger,
	}
}

// Send 发送消息
func (s *Session) Send(data []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Conn == nil {
		return fmt.Errorf("连接已关闭")
	}

	_, err := s.Conn.Write(data)
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}

	s.LastActivity = time.Now()
	s.logger.Info("消息已发送", map[string]interface{}{
		"session_id":  s.ID,
		"data_length": len(data),
		"remote_addr": s.RemoteAddr,
	})

	return nil
}

// Receive 接收消息
func (s *Session) Receive() ([]byte, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Conn == nil {
		return nil, fmt.Errorf("连接已关闭")
	}

	buffer := make([]byte, 4096)
	n, err := s.Conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("接收消息失败: %w", err)
	}

	s.LastActivity = time.Now()
	s.logger.Info("消息已接收", map[string]interface{}{
		"session_id":  s.ID,
		"data_length": n,
		"remote_addr": s.RemoteAddr,
	})

	return buffer[:n], nil
}

// Close 关闭会话
func (s *Session) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Conn == nil {
		return nil
	}

	err := s.Conn.Close()
	s.Conn = nil
	s.Status = "closed"

	s.logger.Info("会话已关闭", map[string]interface{}{
		"session_id":  s.ID,
		"remote_addr": s.RemoteAddr,
	})

	return err
}

// SetUserID 设置用户ID
func (s *Session) SetUserID(userID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.UserID = userID
	s.logger.Info("用户ID已设置", map[string]interface{}{
		"session_id": s.ID,
		"user_id":    userID,
	})
}

// GetUserID 获取用户ID
func (s *Session) GetUserID() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.UserID
}

// SetGroupID 设置组ID
func (s *Session) SetGroupID(groupID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.GroupID = groupID
	s.logger.Info("组ID已设置", map[string]interface{}{
		"session_id": s.ID,
		"group_id":   groupID,
	})
}

// GetGroupID 获取组ID
func (s *Session) GetGroupID() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.GroupID
}

// SetStatus 设置状态
func (s *Session) SetStatus(status string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Status = status
	s.logger.Info("状态已更新", map[string]interface{}{
		"session_id": s.ID,
		"status":     status,
	})
}

// GetStatus 获取状态
func (s *Session) GetStatus() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.Status
}

// IsActive 检查是否活跃
func (s *Session) IsActive() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.Status == "active" && s.Conn != nil
}

// SetReadTimeout 设置读取超时
func (s *Session) SetReadTimeout(timeout time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Conn == nil {
		return fmt.Errorf("连接已关闭")
	}

	return s.Conn.SetReadDeadline(time.Now().Add(timeout))
}

// SetWriteTimeout 设置写入超时
func (s *Session) SetWriteTimeout(timeout time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Conn == nil {
		return fmt.Errorf("连接已关闭")
	}

	return s.Conn.SetWriteDeadline(time.Now().Add(timeout))
}

// GetInfo 获取会话信息
func (s *Session) GetInfo() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return map[string]interface{}{
		"id":            s.ID,
		"remote_addr":   s.RemoteAddr,
		"group_id":      s.GroupID,
		"user_id":       s.UserID,
		"created_at":    s.CreatedAt,
		"last_activity": s.LastActivity,
		"status":        s.Status,
	}
}
