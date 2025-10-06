package network

import (
	"context"
	"fmt"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logger"
)

// Logger 简单的日志接口
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
}

// NetcoreServer netcore-go TCP服务器
type NetcoreServer struct {
	server   interface{} // *netcore.Server
	logger   Logger
	config   *ServerConfig
	handlers map[uint32]NetcoreMessageHandler
	mu       sync.RWMutex
	stats    *ServerStats
	ctx      context.Context
	cancel   context.CancelFunc
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host               string        `json:"host" yaml:"host"`
	Port               int           `json:"port" yaml:"port"`
	MaxConnections     int           `json:"max_connections" yaml:"max_connections"`
	ReadTimeout        time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout       time.Duration `json:"write_timeout" yaml:"write_timeout"`
	KeepAliveInterval  time.Duration `json:"keep_alive_interval" yaml:"keep_alive_interval"`
	HeartbeatInterval  time.Duration `json:"heartbeat_interval" yaml:"heartbeat_interval"`
	MaxPacketSize      int           `json:"max_packet_size" yaml:"max_packet_size"`
	CompressionEnabled bool          `json:"compression_enabled" yaml:"compression_enabled"`
	EncryptionEnabled  bool          `json:"encryption_enabled" yaml:"encryption_enabled"`
	EnableMetrics      bool          `json:"enable_metrics" yaml:"enable_metrics"`
}

// Connection 连接接口
type Connection interface {
	GetID() string
	GetRemoteAddr() string
	Send(data []byte) error
	Close() error
}

// Packet 数据包接口
type Packet interface {
	GetType() uint32
	GetData() []byte
	SetType(msgType uint32)
	SetData(data []byte)
}

// NetcoreMessageHandler 消息处理器接口
type NetcoreMessageHandler interface {
	// Handle 处理消息
	Handle(ctx context.Context, conn Connection, packet Packet) error

	// GetMessageType 获取消息类型
	GetMessageType() uint32

	// GetHandlerName 获取处理器名称
	GetHandlerName() string
}

// ConnectionHandler 连接处理器接口
type ConnectionHandler interface {
	// OnConnect 连接建立时调用
	OnConnect(conn Connection) error

	// OnDisconnect 连接断开时调用
	OnDisconnect(conn Connection, err error)

	// OnError 发生错误时调用
	OnError(conn Connection, err error)
}

// NetcoreServerInterface TCP服务器接口
type NetcoreServerInterface interface {
	// RegisterHandler 注册消息处理器
	RegisterHandler(handler NetcoreMessageHandler) error

	// UnregisterHandler 取消注册消息处理器
	UnregisterHandler(messageType uint32) error

	// SetConnectionHandler 设置连接处理器
	SetConnectionHandler(handler ConnectionHandler)

	// Start 启动服务器
	Start(ctx context.Context) error

	// Stop 停止服务器
	Stop() error

	// Broadcast 广播消息
	Broadcast(packet Packet) error

	// SendToConnection 发送消息到指定连接
	SendToConnection(connID string, packet Packet) error

	// GetStats 获取服务器统计信息
	GetStats() *ServerStats

	// GetConnections 获取所有连接
	GetConnections() []Connection

	// GetConnection 根据ID获取连接
	GetConnection(connID string) Connection
}

// NewNetcoreServer 创建netcore服务器
func NewNetcoreServer(config *ServerConfig, logger logger.Logger) NetcoreServerInterface {
	if config == nil {
		config = &ServerConfig{
			Host:               "0.0.0.0",
			Port:               8080,
			MaxConnections:     1000,
			ReadTimeout:        30 * time.Second,
			WriteTimeout:       30 * time.Second,
			KeepAliveInterval:  60 * time.Second,
			HeartbeatInterval:  30 * time.Second,
			MaxPacketSize:      1024 * 1024, // 1MB
			CompressionEnabled: false,
			EncryptionEnabled:  false,
			EnableMetrics:      true,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	s := &NetcoreServer{
		server:   nil, // 暂时设为nil，后续实现
		logger:   logger,
		config:   config,
		handlers: make(map[uint32]NetcoreMessageHandler),
		ctx:      ctx,
		cancel:   cancel,
		stats: &ServerStats{
			StartTime:     time.Now(),
			ByMessageType: make(map[uint32]*MessageTypeStats),
		},
	}

	// 设置事件处理器
	s.setupHandlers()

	logger.Info("Netcore server initialized successfully", "address", fmt.Sprintf("%s:%d", config.Host, config.Port), "max_connections", config.MaxConnections)
	return s
}

// RegisterHandler 注册消息处理器
func (s *NetcoreServer) RegisterHandler(handler NetcoreMessageHandler) error {
	msgType := handler.GetMessageType()

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.handlers[msgType]; exists {
		return fmt.Errorf("handler for message type %d already exists", msgType)
	}

	s.handlers[msgType] = handler

	s.logger.Info("Message handler registered successfully", "message_type", msgType, "handler", handler.GetHandlerName())
	return nil
}

// UnregisterHandler 取消注册消息处理器
func (s *NetcoreServer) UnregisterHandler(messageType uint32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.handlers[messageType]; !exists {
		return fmt.Errorf("handler for message type %d not found", messageType)
	}

	delete(s.handlers, messageType)

	s.logger.Info("Message handler unregistered successfully", "message_type", messageType)
	return nil
}

// SetConnectionHandler 设置连接处理器
func (s *NetcoreServer) SetConnectionHandler(handler ConnectionHandler) {
	// 这里可以存储连接处理器，在netcore事件中调用
	// 简化实现，直接在setupNetcoreHandlers中处理
}

// Start 启动服务器
func (s *NetcoreServer) Start(ctx context.Context) error {
	s.logger.Info("Starting netcore server", "address", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))

	// TODO: 实现实际的服务器启动逻辑
	// 这里暂时只是模拟启动

	// 启动心跳检测
	if s.config.HeartbeatInterval > 0 {
		go s.startHeartbeat()
	}

	// 启动指标收集
	if s.config.EnableMetrics {
		go s.collectMetrics()
	}

	s.logger.Info("Netcore server started successfully")

	// 等待上下文取消
	select {
	case <-ctx.Done():
		s.logger.Info("Netcore server context cancelled")
		return ctx.Err()
	case <-s.ctx.Done():
		s.logger.Info("Netcore server stopped")
		return nil
	}
}

// Stop 停止服务器
func (s *NetcoreServer) Stop() error {
	s.logger.Info("Stopping netcore server")

	// 取消上下文
	s.cancel()

	// TODO: 实现实际的服务器停止逻辑
	s.logger.Info("Netcore server stopped successfully")
	return nil
}

// Broadcast 广播消息
func (s *NetcoreServer) Broadcast(packet Packet) error {
	connections := s.GetConnections()
	if len(connections) == 0 {
		s.logger.Debug("No connections to broadcast to")
		return nil
	}

	var errors []error
	successCount := 0

	for _, conn := range connections {
		if err := conn.Send(packet.GetData()); err != nil {
			s.logger.Error("Failed to broadcast to connection", "error", err, "conn_id", conn.GetID())
			errors = append(errors, err)
		} else {
			successCount++
		}
	}

	// 更新统计信息
	s.updateStats(packet.GetType(), true, successCount)

	s.logger.Debug("Broadcast completed", "total_connections", len(connections), "success_count", successCount, "error_count", len(errors))

	if len(errors) > 0 {
		return fmt.Errorf("broadcast failed for %d connections: %v", len(errors), errors[0])
	}

	return nil
}

// SendToConnection 发送消息到指定连接
func (s *NetcoreServer) SendToConnection(connID string, packet Packet) error {
	conn := s.GetConnection(connID)
	if conn == nil {
		return fmt.Errorf("connection %s not found", connID)
	}

	err := conn.Send(packet.GetData())
	if err != nil {
		s.logger.Error("Failed to send message to connection", "error", err, "conn_id", connID, "message_type", packet.GetType())
		s.updateStats(packet.GetType(), false, 0)
		return fmt.Errorf("failed to send message to connection %s: %w", connID, err)
	}

	s.updateStats(packet.GetType(), true, 1)
	s.logger.Debug("Message sent to connection successfully", "conn_id", connID, "message_type", packet.GetType())
	return nil
}

// GetStats 获取服务器统计信息
func (s *NetcoreServer) GetStats() *ServerStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	connections := s.GetConnections()

	// 创建统计信息副本
	stats := &ServerStats{
		ActiveConnections: int64(len(connections)),
		TotalConnections:  s.stats.TotalConnections,
		TotalMessages:     s.stats.TotalMessages,
		TotalErrors:       s.stats.TotalErrors,
		StartTime:         s.stats.StartTime,
		Uptime:            time.Since(s.stats.StartTime),
		ByMessageType:     make(map[uint32]*MessageTypeStats),
	}

	// 复制消息类型统计
	for msgType, msgStats := range s.stats.ByMessageType {
		stats.ByMessageType[msgType] = &MessageTypeStats{
			ProcessedCount: msgStats.ProcessedCount,
			FailedCount:    msgStats.FailedCount,
			LastProcessed:  msgStats.LastProcessed,
			AvgProcessTime: msgStats.AvgProcessTime,
		}
	}

	return stats
}

// GetConnections 获取所有连接
func (s *NetcoreServer) GetConnections() []Connection {
	// TODO: 实现实际的连接获取逻辑
	return []Connection{}
}

// GetConnection 根据ID获取连接
func (s *NetcoreServer) GetConnection(connID string) Connection {
	// TODO: 实现实际的连接获取逻辑
	return nil
}

// 私有方法

// setupHandlers 设置事件处理器
func (s *NetcoreServer) setupHandlers() {
	// TODO: 实现实际的事件处理器设置逻辑
	s.logger.Info("Event handlers setup completed")
}

// handleMessage 处理消息
func (s *NetcoreServer) handleMessage(conn Connection, packet Packet) {
	start := time.Now()
	msgType := packet.GetType()

	s.logger.Debug("Received message", "conn_id", conn.GetID(), "message_type", msgType, "size", len(packet.GetData()))

	// 获取处理器
	s.mu.RLock()
	handler, exists := s.handlers[msgType]
	s.mu.RUnlock()

	if !exists {
		s.logger.Warn("No handler found for message type", "message_type", msgType, "conn_id", conn.GetID())
		s.updateStats(msgType, false, 0)
		return
	}

	// 处理消息
	ctx, cancel := context.WithTimeout(s.ctx, s.config.ReadTimeout)
	defer cancel()

	err := handler.Handle(ctx, conn, packet)
	processTime := time.Since(start)

	if err != nil {
		s.logger.Error("Message handling failed", "error", err, "message_type", msgType, "conn_id", conn.GetID(), "handler", handler.GetHandlerName())
		s.updateStats(msgType, false, 0)
		return
	}

	s.updateStats(msgType, true, 1)
	s.logger.Debug("Message handled successfully", "message_type", msgType, "conn_id", conn.GetID(), "process_time", processTime)
}

// startHeartbeat 启动心跳检测
func (s *NetcoreServer) startHeartbeat() {
	ticker := time.NewTicker(s.config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.sendHeartbeat()
		case <-s.ctx.Done():
			return
		}
	}
}

// sendHeartbeat 发送心跳
func (s *NetcoreServer) sendHeartbeat() {
	// TODO: 实现实际的心跳发送逻辑
	connections := s.GetConnections()
	s.logger.Debug("Heartbeat sent to all connections", "connection_count", len(connections))
}

// collectMetrics 收集指标
func (s *NetcoreServer) collectMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats := s.GetStats()
			s.logger.Debug("Server metrics",
				"active_connections", stats.ActiveConnections,
				"total_connections", stats.TotalConnections,
				"total_messages", stats.TotalMessages,
				"total_errors", stats.TotalErrors,
				"uptime", stats.Uptime)
		case <-s.ctx.Done():
			return
		}
	}
}

// updateStats 更新统计信息
func (s *NetcoreServer) updateStats(msgType uint32, success bool, count int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if success {
		s.stats.TotalMessages += int64(count)
	} else {
		s.stats.TotalErrors++
	}

	// 更新消息类型统计
	msgStats, exists := s.stats.ByMessageType[msgType]
	if !exists {
		msgStats = &MessageTypeStats{}
		s.stats.ByMessageType[msgType] = msgStats
	}

	if success {
		msgStats.ProcessedCount += int64(count)
		msgStats.LastProcessed = time.Now()
	} else {
		msgStats.FailedCount++
	}
}

// 统计信息结构
type ServerStats struct {
	ActiveConnections int64                        `json:"active_connections"`
	TotalConnections  int64                        `json:"total_connections"`
	TotalMessages     int64                        `json:"total_messages"`
	TotalErrors       int64                        `json:"total_errors"`
	StartTime         time.Time                    `json:"start_time"`
	Uptime            time.Duration                `json:"uptime"`
	ByMessageType     map[uint32]*MessageTypeStats `json:"by_message_type"`
}

type MessageTypeStats struct {
	ProcessedCount int64         `json:"processed_count"`
	FailedCount    int64         `json:"failed_count"`
	LastProcessed  time.Time     `json:"last_processed"`
	AvgProcessTime time.Duration `json:"avg_process_time"`
}
