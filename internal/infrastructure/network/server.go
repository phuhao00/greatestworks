// Package network 网络通信基础设施
package network

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/phuhao00/netcore-go/network"
	"github.com/phuhao00/netcore-go/protocol"
	"github.com/phuhao00/netcore-go/session"
)

// 使用netcore_server.go中定义的ServerConfig

// DefaultServerConfig 默认服务器配置
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Host:           "0.0.0.0",
		Port:           8080,
		MaxConnections: 1000,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		Heartbeat:      60 * time.Second,
	}
}

// MessageHandler 消息处理器接口
type MessageHandler interface {
	Handle(ctx context.Context, session *session.Session, msg *protocol.Message) error
}

// MessageHandlerFunc 消息处理器函数类型
type MessageHandlerFunc func(ctx context.Context, session *session.Session, msg *protocol.Message) error

// Handle 实现MessageHandler接口
func (f MessageHandlerFunc) Handle(ctx context.Context, session *session.Session, msg *protocol.Message) error {
	return f(ctx, session, msg)
}

// Server TCP服务器
type Server struct {
	config   *ServerConfig
	server   *network.TCPServer
	handlers map[uint32]MessageHandler
	sessions map[string]*session.Session
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewServer 创建新的TCP服务器
func NewServer(config *ServerConfig) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		config:   config,
		handlers: make(map[uint32]MessageHandler),
		sessions: make(map[string]*session.Session),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// RegisterHandler 注册消息处理器
func (s *Server) RegisterHandler(msgType uint32, handler MessageHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[msgType] = handler
}

// RegisterHandlerFunc 注册消息处理器函数
func (s *Server) RegisterHandlerFunc(msgType uint32, handler MessageHandlerFunc) {
	s.RegisterHandler(msgType, handler)
}

// Start 启动服务器
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// 创建netcore-go TCP服务器
	s.server = network.NewTCPServer(addr)

	// 设置连接处理器
	s.server.SetOnConnect(s.onConnect)
	s.server.SetOnDisconnect(s.onDisconnect)
	s.server.SetOnMessage(s.onMessage)

	log.Printf("TCP服务器启动在 %s", addr)
	return s.server.Start()
}

// Stop 停止服务器
func (s *Server) Stop() error {
	s.cancel()
	if s.server != nil {
		return s.server.Stop()
	}
	return nil
}

// onConnect 连接建立处理
func (s *Server) onConnect(sess *session.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionID := sess.ID()
	s.sessions[sessionID] = sess

	log.Printf("客户端连接: %s, 总连接数: %d", sessionID, len(s.sessions))

	// 设置会话属性
	sess.SetAttribute("connect_time", time.Now())
	sess.SetAttribute("last_heartbeat", time.Now())
}

// onDisconnect 连接断开处理
func (s *Server) onDisconnect(sess *session.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionID := sess.ID()
	delete(s.sessions, sessionID)

	log.Printf("客户端断开: %s, 剩余连接数: %d", sessionID, len(s.sessions))
}

// onMessage 消息处理
func (s *Server) onMessage(sess *session.Session, data []byte) {
	// 解析消息
	msg, err := protocol.DecodeMessage(data)
	if err != nil {
		log.Printf("消息解析失败: %v", err)
		return
	}

	// 更新心跳时间
	sess.SetAttribute("last_heartbeat", time.Now())

	// 查找处理器
	s.mu.RLock()
	handler, exists := s.handlers[msg.Type]
	s.mu.RUnlock()

	if !exists {
		log.Printf("未找到消息类型 %d 的处理器", msg.Type)
		return
	}

	// 处理消息
	go func() {
		ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
		defer cancel()

		if err := handler.Handle(ctx, sess, msg); err != nil {
			log.Printf("消息处理失败: %v", err)
		}
	}()
}

// Broadcast 广播消息给所有连接
func (s *Server) Broadcast(msg *protocol.Message) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := protocol.EncodeMessage(msg)
	if err != nil {
		return fmt.Errorf("消息编码失败: %w", err)
	}

	for _, sess := range s.sessions {
		if err := sess.Send(data); err != nil {
			log.Printf("发送消息到会话 %s 失败: %v", sess.ID(), err)
		}
	}

	return nil
}

// SendToSession 发送消息给指定会话
func (s *Server) SendToSession(sessionID string, msg *protocol.Message) error {
	s.mu.RLock()
	sess, exists := s.sessions[sessionID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("会话 %s 不存在", sessionID)
	}

	data, err := protocol.EncodeMessage(msg)
	if err != nil {
		return fmt.Errorf("消息编码失败: %w", err)
	}

	return sess.Send(data)
}

// GetSession 获取指定会话
func (s *Server) GetSession(sessionID string) (*session.Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, exists := s.sessions[sessionID]
	return sess, exists
}

// GetSessionCount 获取当前连接数
func (s *Server) GetSessionCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.sessions)
}

// GetAllSessions 获取所有会话
func (s *Server) GetAllSessions() []*session.Session {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sessions := make([]*session.Session, 0, len(s.sessions))
	for _, sess := range s.sessions {
		sessions = append(sessions, sess)
	}
	return sessions
}

// StartHeartbeatChecker 启动心跳检查器
func (s *Server) StartHeartbeatChecker() {
	ticker := time.NewTicker(s.config.Heartbeat)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ticker.C:
				s.checkHeartbeat()
			}
		}
	}()
}

// checkHeartbeat 检查心跳超时
func (s *Server) checkHeartbeat() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	timeout := s.config.Heartbeat * 2 // 心跳超时时间为心跳间隔的2倍

	for sessionID, sess := range s.sessions {
		lastHeartbeat, ok := sess.GetAttribute("last_heartbeat").(time.Time)
		if !ok {
			continue
		}

		if now.Sub(lastHeartbeat) > timeout {
			log.Printf("会话 %s 心跳超时，断开连接", sessionID)
			sess.Close()
		}
	}
}
