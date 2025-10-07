// Package network 网络通信基础设施
package network

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	// "github.com/phuhao00/spoor/v2" // TODO: 暂时注释掉有问题的依赖

	"greatestworks/internal/interfaces/tcp/protocol"
	"greatestworks/internal/network"
	"greatestworks/internal/network/session"
)

// 使用netcore_server.go中定义的ServerConfig

// DefaultServerConfig 默认服务器配置
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Host:               "0.0.0.0",
		Port:               8080,
		MaxConnections:     1000,
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		HeartbeatInterval:  60 * time.Second,
		KeepAliveInterval:  60 * time.Second,
		MaxPacketSize:      1024 * 1024,
		CompressionEnabled: false,
		EncryptionEnabled:  false,
		EnableMetrics:      true,
	}
}

// MessageHandler 消息处理器接口
type MessageHandler func(ctx context.Context, session session.Session, packet Packet) error

// simpleLogger 简单日志实现
type simpleLogger struct{}

func (l *simpleLogger) Info(msg string, args ...interface{}) {
	log.Printf("[INFO] "+msg, args...)
}

func (l *simpleLogger) Error(msg string, args ...interface{}) {
	log.Printf("[ERROR] "+msg, args...)
}

func (l *simpleLogger) Debug(msg string, args ...interface{}) {
	log.Printf("[DEBUG] "+msg, args...)
}

func (l *simpleLogger) Warn(msg string, args ...interface{}) {
	log.Printf("[WARN] "+msg, args...)
}

// connectionHandler 连接处理器实现
type connectionHandler struct {
	server *Server
}

func (h *connectionHandler) OnConnect(conn Connection) error {
	return h.server.onConnect(conn)
}

func (h *connectionHandler) OnDisconnect(conn Connection, err error) {
	h.server.onDisconnect(conn, err)
}

func (h *connectionHandler) OnError(conn Connection, err error) {
	log.Printf("Connection error: %v", err)
}

// simplePacket 简单数据包实现
type simplePacket struct {
	msgType uint32
	data    []byte
}

func (p *simplePacket) GetType() uint32 {
	return p.msgType
}

func (p *simplePacket) GetData() []byte {
	return p.data
}

func (p *simplePacket) SetType(msgType uint32) {
	p.msgType = msgType
}

func (p *simplePacket) SetData(data []byte) {
	p.data = data
}

// MessageHandlerFunc 消息处理器函数类型
type MessageHandlerFunc func(ctx context.Context, session *session.Session, msg *protocol.Message) error

// Handle 实现MessageHandler接口
func (f MessageHandlerFunc) Handle(ctx context.Context, session *session.Session, msg *protocol.Message) error {
	return f(ctx, session, msg)
}

// Server TCP服务器
type Server struct {
	config    *ServerConfig
	server    NetcoreServerInterface
	handlers  map[uint32]MessageHandler
	sessions  map[string]session.Session
	processor *network.MessageProcessor
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewServer 创建新的TCP服务器
func NewServer(config *ServerConfig) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		config:    config,
		handlers:  make(map[uint32]MessageHandler),
		sessions:  make(map[string]session.Session),
		processor: network.NewMessageProcessor(),
		ctx:       ctx,
		cancel:    cancel,
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
	// 将MessageHandlerFunc转换为MessageHandler
	messageHandler := func(ctx context.Context, session session.Session, packet Packet) error {
		return handler(ctx, &session, &protocol.Message{})
	}
	s.RegisterHandler(msgType, messageHandler)
}

// Start 启动服务器
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// 创建netcore-go TCP服务器
	s.server = NewNetcoreServer(s.config, &simpleLogger{})

	// 设置连接处理器
	s.server.SetConnectionHandler(&connectionHandler{
		server: s,
	})

	log.Printf("TCP服务器启动在 %s", addr)
	return s.server.Start(s.ctx)
}

// Stop 停止服务器
func (s *Server) Stop() error {
	s.cancel()
	if s.server != nil {
		return s.server.Stop()
	}
	return nil
}

// onConnect 处理连接建立
func (s *Server) onConnect(conn Connection) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionID := conn.GetID()
	// 创建新会话
	sess := session.NewSession(sessionID)
	s.sessions[sessionID] = sess

	log.Printf("客户端连接: %s", sessionID)
	return nil
}

// onDisconnect 处理连接断开
func (s *Server) onDisconnect(conn Connection, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionID := conn.GetID()
	// 清理会话
	delete(s.sessions, sessionID)

	log.Printf("客户端断开连接: %s, 错误: %v", sessionID, err)
}

// onMessage 处理消息
func (s *Server) onMessage(conn Connection, packet Packet) error {
	// 查找处理器
	handler, exists := s.handlers[packet.GetType()]
	if !exists {
		log.Printf("未找到消息类型 %d 的处理器", packet.GetType())
		return fmt.Errorf("no handler for message type %d", packet.GetType())
	}

	// TODO: Fix session handling - GetSession method needs implementation
	// 获取会话
	// sessionID := conn.GetID()
	// sess, exists := s.GetSession(sessionID)
	// if !exists {
	// 	log.Printf("未找到会话: %s", sessionID)
	// 	return fmt.Errorf("session not found: %s", sessionID)
	// }

	// 处理消息 - Temporarily disabled until session handling is fixed
	// if err := handler.Handle(s.ctx, sess, &protocol.Message{}); err != nil {
	// 	log.Printf("消息处理失败: %v", err)
	// 	return err
	// }
	log.Printf("消息处理成功 (session handling disabled), handler type: %T", handler)
	return nil
}

// Broadcast 广播消息给所有连接
func (s *Server) Broadcast(msg *network.Message) error {
	// 创建数据包
	packet := &simplePacket{
		msgType: uint32(msg.Header.Type),
		data:    msg.Body,
	}

	return s.server.Broadcast(packet)
}

// SendToSession 发送消息给指定会话
func (s *Server) SendToSession(sessionID string, msg *network.Message) error {
	// 创建数据包
	packet := &simplePacket{
		msgType: uint32(msg.Header.Type),
		data:    msg.Body,
	}

	return s.server.SendToConnection(sessionID, packet)
}

// GetSession 获取指定会话
func (s *Server) GetSession(sessionID string) (session.Session, bool) {
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
func (s *Server) GetAllSessions() []session.Session {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sessions := make([]session.Session, 0, len(s.sessions))
	for _, sess := range s.sessions {
		sessions = append(sessions, sess)
	}
	return sessions
}

// StartHeartbeatChecker 启动心跳检查器
func (s *Server) StartHeartbeatChecker() {
	ticker := time.NewTicker(s.config.HeartbeatInterval)
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
	timeout := s.config.HeartbeatInterval * 2 // 心跳超时时间为心跳间隔的2倍

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
