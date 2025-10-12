package session

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Session 会话接口
type Session interface {
	ID() string
	PlayerID() string
	SetPlayerID(playerID string)
	Conn() net.Conn
	Send(data []byte) error
	Close() error
	IsActive() bool
	LastActivity() time.Time
	UpdateActivity()
	GetAttribute(key string) interface{}
	SetAttribute(key string, value interface{})
}

// TCPSession TCP会话实现
type TCPSession struct {
	id           string
	playerID     string
	conn         net.Conn
	lastActivity time.Time
	active       bool
	attributes   map[string]interface{}
	mutex        sync.RWMutex
}

// NewTCPSession 创建新的TCP会话
func NewTCPSession(id string, conn net.Conn) *TCPSession {
	return &TCPSession{
		id:           id,
		conn:         conn,
		lastActivity: time.Now(),
		active:       true,
		attributes:   make(map[string]interface{}),
	}
}

// ID 获取会话ID
func (s *TCPSession) ID() string {
	return s.id
}

// PlayerID 获取玩家ID
func (s *TCPSession) PlayerID() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.playerID
}

// SetPlayerID 设置玩家ID
func (s *TCPSession) SetPlayerID(playerID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.playerID = playerID
}

// Conn 获取连接
func (s *TCPSession) Conn() net.Conn {
	return s.conn
}

// Send 发送数据
func (s *TCPSession) Send(data []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.active {
		return ErrSessionClosed
	}

	_, err := s.conn.Write(data)
	if err != nil {
		s.active = false
		return err
	}

	s.lastActivity = time.Now()
	return nil
}

// Close 关闭会话
func (s *TCPSession) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.active {
		return nil
	}

	s.active = false
	return s.conn.Close()
}

// IsActive 检查会话是否活跃
func (s *TCPSession) IsActive() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.active
}

// LastActivity 获取最后活动时间
func (s *TCPSession) LastActivity() time.Time {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastActivity
}

// UpdateActivity 更新活动时间
func (s *TCPSession) UpdateActivity() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastActivity = time.Now()
}

// GetAttribute 获取属性
func (s *TCPSession) GetAttribute(key string) interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.attributes[key]
}

// SetAttribute 设置属性
func (s *TCPSession) SetAttribute(key string, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.attributes[key] = value
}

// NewSession 创建新会话（简单实现，用于兼容）
func NewSession(id string) Session {
	return &TCPSession{
		id:           id,
		conn:         nil, // 没有实际连接
		lastActivity: time.Now(),
		active:       true,
		attributes:   make(map[string]interface{}),
	}
}

// 错误定义
var (
	ErrSessionClosed = fmt.Errorf("session is closed")
)
