package network

import (
	"context"
	"fmt"
	"sync"
	"time"

	"greatestworks/aop/logger"

	"github.com/phuhao00/netcore-go"
)

// NetcoreClient netcore-go TCP客户�?type NetcoreClient struct {
	client   *netcore.Client
	logger   logger.Logger
	config   *ClientConfig
	handlers map[uint32]ClientMessageHandler
	mu       sync.RWMutex
	stats    *ClientStats
	ctx      context.Context
	cancel   context.CancelFunc
	conn     *netcore.Connection
}

// ClientConfig 客户端配�?type ClientConfig struct {
	ServerHost           string        `json:"server_host" yaml:"server_host"`
	ServerPort           int           `json:"server_port" yaml:"server_port"`
	ConnectTimeout       time.Duration `json:"connect_timeout" yaml:"connect_timeout"`
	ReadTimeout          time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout         time.Duration `json:"write_timeout" yaml:"write_timeout"`
	KeepAliveInterval    time.Duration `json:"keep_alive_interval" yaml:"keep_alive_interval"`
	HeartbeatInterval    time.Duration `json:"heartbeat_interval" yaml:"heartbeat_interval"`
	ReconnectInterval    time.Duration `json:"reconnect_interval" yaml:"reconnect_interval"`
	MaxReconnectAttempts int           `json:"max_reconnect_attempts" yaml:"max_reconnect_attempts"`
	MaxPacketSize        int           `json:"max_packet_size" yaml:"max_packet_size"`
	CompressionEnabled   bool          `json:"compression_enabled" yaml:"compression_enabled"`
	EncryptionEnabled    bool          `json:"encryption_enabled" yaml:"encryption_enabled"`
	EnableMetrics        bool          `json:"enable_metrics" yaml:"enable_metrics"`
	AutoReconnect        bool          `json:"auto_reconnect" yaml:"auto_reconnect"`
}

// ClientMessageHandler 客户端消息处理器接口
type ClientMessageHandler interface {
	// Handle 处理消息
	Handle(ctx context.Context, packet *netcore.Packet) error

	// GetMessageType 获取消息类型
	GetMessageType() uint32

	// GetHandlerName 获取处理器名�?	GetHandlerName() string
}

// ClientConnectionHandler 客户端连接处理器接口
type ClientConnectionHandler interface {
	// OnConnect 连接建立时调�?	OnConnect() error

	// OnDisconnect 连接断开时调�?	OnDisconnect(err error)

	// OnReconnect 重连成功时调�?	OnReconnect() error

	// OnError 发生错误时调�?	OnError(err error)
}

// Client TCP客户端接�?type Client interface {
	// RegisterHandler 注册消息处理�?	RegisterHandler(handler ClientMessageHandler) error

	// UnregisterHandler 取消注册消息处理�?	UnregisterHandler(messageType uint32) error

	// SetConnectionHandler 设置连接处理�?	SetConnectionHandler(handler ClientConnectionHandler)

	// Connect 连接到服务器
	Connect(ctx context.Context) error

	// Disconnect 断开连接
	Disconnect() error

	// Send 发送消�?	Send(packet *netcore.Packet) error

	// SendWithTimeout 带超时的发送消�?	SendWithTimeout(packet *netcore.Packet, timeout time.Duration) error

	// SendRequest 发送请求并等待响应
	SendRequest(request *netcore.Packet, timeout time.Duration) (*netcore.Packet, error)

	// IsConnected 检查是否已连接
	IsConnected() bool

	// GetStats 获取客户端统计信�?	GetStats() *ClientStats

	// Start 启动客户�?	Start(ctx context.Context) error

	// Stop 停止客户�?	Stop() error
}

// NewNetcoreClient 创建netcore客户�?func NewNetcoreClient(config *ClientConfig, logger logger.Logger) Client {
	if config == nil {
		config = &ClientConfig{
			ServerHost:           "localhost",
			ServerPort:           8080,
			ConnectTimeout:       10 * time.Second,
			ReadTimeout:          30 * time.Second,
			WriteTimeout:         30 * time.Second,
			KeepAliveInterval:    60 * time.Second,
			HeartbeatInterval:    30 * time.Second,
			ReconnectInterval:    5 * time.Second,
			MaxReconnectAttempts: 10,
			MaxPacketSize:        1024 * 1024, // 1MB
			CompressionEnabled:   false,
			EncryptionEnabled:    false,
			EnableMetrics:        true,
			AutoReconnect:        true,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	// 创建netcore客户端配�?	netcoreConfig := &netcore.ClientConfig{
		ServerAddress:     fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort),
		ConnectTimeout:    config.ConnectTimeout,
		ReadTimeout:       config.ReadTimeout,
		WriteTimeout:      config.WriteTimeout,
		KeepAliveInterval: config.KeepAliveInterval,
		MaxPacketSize:     config.MaxPacketSize,
	}

	// 创建netcore客户�?	netcoreClient := netcore.NewClient(netcoreConfig)

	c := &NetcoreClient{
		client:   netcoreClient,
		logger:   logger,
		config:   config,
		handlers: make(map[uint32]ClientMessageHandler),
		ctx:      ctx,
		cancel:   cancel,
		stats: &ClientStats{
			StartTime:     time.Now(),
			ByMessageType: make(map[uint32]*ClientMessageTypeStats),
		},
	}

	// 设置netcore事件处理�?	c.setupNetcoreHandlers()

	logger.Info("Netcore client initialized successfully", "server_address", netcoreConfig.ServerAddress)
	return c
}

// RegisterHandler 注册消息处理�?func (c *NetcoreClient) RegisterHandler(handler ClientMessageHandler) error {
	msgType := handler.GetMessageType()

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.handlers[msgType]; exists {
		return fmt.Errorf("handler for message type %d already exists", msgType)
	}

	c.handlers[msgType] = handler

	c.logger.Info("Message handler registered successfully", "message_type", msgType, "handler", handler.GetHandlerName())
	return nil
}

// UnregisterHandler 取消注册消息处理�?func (c *NetcoreClient) UnregisterHandler(messageType uint32) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.handlers[messageType]; !exists {
		return fmt.Errorf("handler for message type %d not found", messageType)
	}

	delete(c.handlers, messageType)

	c.logger.Info("Message handler unregistered successfully", "message_type", messageType)
	return nil
}

// SetConnectionHandler 设置连接处理�?func (c *NetcoreClient) SetConnectionHandler(handler ClientConnectionHandler) {
	// 这里可以存储连接处理器，在netcore事件中调�?	// 简化实现，直接在setupNetcoreHandlers中处�?}

// Connect 连接到服务器
func (c *NetcoreClient) Connect(ctx context.Context) error {
	c.logger.Info("Connecting to server", "address", fmt.Sprintf("%s:%d", c.config.ServerHost, c.config.ServerPort))

	conn, err := c.client.Connect()
	if err != nil {
		c.logger.Error("Failed to connect to server", "error", err)
		return fmt.Errorf("failed to connect to server: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.stats.ConnectTime = time.Now()
	c.stats.TotalConnections++
	c.mu.Unlock()

	c.logger.Info("Connected to server successfully", "conn_id", conn.GetID())
	return nil
}

// Disconnect 断开连接
func (c *NetcoreClient) Disconnect() error {
	c.logger.Info("Disconnecting from server")

	c.mu.Lock()
	conn := c.conn
	c.conn = nil
	c.mu.Unlock()

	if conn != nil {
		if err := conn.Close(); err != nil {
			c.logger.Error("Failed to close connection", "error", err)
			return fmt.Errorf("failed to close connection: %w", err)
		}
	}

	c.logger.Info("Disconnected from server successfully")
	return nil
}

// Send 发送消�?func (c *NetcoreClient) Send(packet *netcore.Packet) error {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	if conn == nil {
		return fmt.Errorf("not connected to server")
	}

	err := conn.Send(packet)
	if err != nil {
		c.logger.Error("Failed to send message", "error", err, "message_type", packet.GetType())
		c.updateStats(packet.GetType(), false, 0)
		return fmt.Errorf("failed to send message: %w", err)
	}

	c.updateStats(packet.GetType(), true, 0)
	c.logger.Debug("Message sent successfully", "message_type", packet.GetType(), "size", len(packet.GetData()))
	return nil
}

// SendWithTimeout 带超时的发送消�?func (c *NetcoreClient) SendWithTimeout(packet *netcore.Packet, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(c.ctx, timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- c.Send(packet)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("send timeout after %v", timeout)
	}
}

// SendRequest 发送请求并等待响应
func (c *NetcoreClient) SendRequest(request *netcore.Packet, timeout time.Duration) (*netcore.Packet, error) {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	if conn == nil {
		return nil, fmt.Errorf("not connected to server")
	}

	ctx, cancel := context.WithTimeout(c.ctx, timeout)
	defer cancel()

	// 这里需要实现请�?响应模式
	// 简化实现，直接发送并等待特定类型的响�?	if err := c.Send(request); err != nil {
		return nil, err
	}

	// 等待响应（这里需要根据实际协议实现）
	// 简化处理，返回nil表示需要在实际使用中完�?	c.logger.Debug("Request sent, waiting for response", "request_type", request.GetType())
	return nil, fmt.Errorf("request-response not implemented")
}

// IsConnected 检查是否已连接
func (c *NetcoreClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn != nil && !c.conn.IsClosed()
}

// GetStats 获取客户端统计信�?func (c *NetcoreClient) GetStats() *ClientStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 创建统计信息副本
	stats := &ClientStats{
		IsConnected:      c.IsConnected(),
		TotalConnections: c.stats.TotalConnections,
		TotalMessages:    c.stats.TotalMessages,
		TotalErrors:      c.stats.TotalErrors,
		ReconnectCount:   c.stats.ReconnectCount,
		ConnectTime:      c.stats.ConnectTime,
		StartTime:        c.stats.StartTime,
		Uptime:           time.Since(c.stats.StartTime),
		ByMessageType:    make(map[uint32]*ClientMessageTypeStats),
	}

	// 复制消息类型统计
	for msgType, msgStats := range c.stats.ByMessageType {
		stats.ByMessageType[msgType] = &ClientMessageTypeStats{
			SentCount:      msgStats.SentCount,
			ReceivedCount:  msgStats.ReceivedCount,
			FailedCount:    msgStats.FailedCount,
			LastSent:       msgStats.LastSent,
			LastReceived:   msgStats.LastReceived,
			AvgProcessTime: msgStats.AvgProcessTime,
		}
	}

	return stats
}

// Start 启动客户�?func (c *NetcoreClient) Start(ctx context.Context) error {
	c.logger.Info("Starting netcore client")

	// 连接到服务器
	if err := c.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect during start: %w", err)
	}

	// 启动心跳检�?	if c.config.HeartbeatInterval > 0 {
		go c.startHeartbeat()
	}

	// 启动自动重连
	if c.config.AutoReconnect {
		go c.startAutoReconnect()
	}

	// 启动指标收集
	if c.config.EnableMetrics {
		go c.collectMetrics()
	}

	c.logger.Info("Netcore client started successfully")

	// 等待上下文取�?	select {
	case <-ctx.Done():
		c.logger.Info("Netcore client context cancelled")
		return ctx.Err()
	case <-c.ctx.Done():
		c.logger.Info("Netcore client stopped")
		return nil
	}
}

// Stop 停止客户�?func (c *NetcoreClient) Stop() error {
	c.logger.Info("Stopping netcore client")

	// 取消上下�?	c.cancel()

	// 断开连接
	if err := c.Disconnect(); err != nil {
		c.logger.Error("Failed to disconnect during stop", "error", err)
	}

	c.logger.Info("Netcore client stopped successfully")
	return nil
}

// 私有方法

// setupNetcoreHandlers 设置netcore事件处理�?func (c *NetcoreClient) setupNetcoreHandlers() {
	// 连接建立事件
	c.client.OnConnect(func(conn *netcore.Connection) {
		c.logger.Info("Connected to server", "conn_id", conn.GetID(), "server_addr", conn.GetRemoteAddr())
	})

	// 连接断开事件
	c.client.OnDisconnect(func(conn *netcore.Connection, err error) {
		if err != nil {
			c.logger.Warn("Disconnected from server with error", "conn_id", conn.GetID(), "error", err)
		} else {
			c.logger.Info("Disconnected from server", "conn_id", conn.GetID())
		}

		c.mu.Lock()
		c.conn = nil
		c.mu.Unlock()
	})

	// 消息接收事件
	c.client.OnMessage(func(conn *netcore.Connection, packet *netcore.Packet) {
		c.handleMessage(packet)
	})

	// 错误事件
	c.client.OnError(func(conn *netcore.Connection, err error) {
		c.mu.Lock()
		c.stats.TotalErrors++
		c.mu.Unlock()

		c.logger.Error("Client error", "conn_id", conn.GetID(), "error", err)
	})
}

// handleMessage 处理消息
func (c *NetcoreClient) handleMessage(packet *netcore.Packet) {
	start := time.Now()
	msgType := packet.GetType()

	c.logger.Debug("Received message", "message_type", msgType, "size", len(packet.GetData()))

	// 获取处理�?	c.mu.RLock()
	handler, exists := c.handlers[msgType]
	c.mu.RUnlock()

	if !exists {
		c.logger.Warn("No handler found for message type", "message_type", msgType)
		c.updateStats(msgType, false, 1)
		return
	}

	// 处理消息
	ctx, cancel := context.WithTimeout(c.ctx, c.config.ReadTimeout)
	defer cancel()

	err := handler.Handle(ctx, packet)
	processTime := time.Since(start)

	if err != nil {
		c.logger.Error("Message handling failed", "error", err, "message_type", msgType, "handler", handler.GetHandlerName())
		c.updateStats(msgType, false, 1)
		return
	}

	c.updateStats(msgType, true, 1)
	c.logger.Debug("Message handled successfully", "message_type", msgType, "process_time", processTime)
}

// startHeartbeat 启动心跳检�?func (c *NetcoreClient) startHeartbeat() {
	ticker := time.NewTicker(c.config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.sendHeartbeat()
		case <-c.ctx.Done():
			return
		}
	}
}

// sendHeartbeat 发送心�?func (c *NetcoreClient) sendHeartbeat() {
	if !c.IsConnected() {
		return
	}

	heartbeatPacket := netcore.NewPacket(0, []byte("heartbeat")) // 消息类型0作为心跳

	if err := c.Send(heartbeatPacket); err != nil {
		c.logger.Debug("Failed to send heartbeat", "error", err)
	} else {
		c.logger.Debug("Heartbeat sent successfully")
	}
}

// startAutoReconnect 启动自动重连
func (c *NetcoreClient) startAutoReconnect() {
	ticker := time.NewTicker(c.config.ReconnectInterval)
	defer ticker.Stop()

	reconnectAttempts := 0

	for {
		select {
		case <-ticker.C:
			if !c.IsConnected() && reconnectAttempts < c.config.MaxReconnectAttempts {
				reconnectAttempts++
				c.logger.Info("Attempting to reconnect", "attempt", reconnectAttempts, "max_attempts", c.config.MaxReconnectAttempts)

				ctx, cancel := context.WithTimeout(c.ctx, c.config.ConnectTimeout)
				if err := c.Connect(ctx); err != nil {
					c.logger.Error("Reconnect attempt failed", "error", err, "attempt", reconnectAttempts)
				} else {
					c.logger.Info("Reconnected successfully", "attempt", reconnectAttempts)
					c.mu.Lock()
					c.stats.ReconnectCount++
					c.mu.Unlock()
					reconnectAttempts = 0 // 重置重连计数
				}
				cancel()
			} else if c.IsConnected() {
				reconnectAttempts = 0 // 连接正常，重置计�?			}
		case <-c.ctx.Done():
			return
		}
	}
}

// collectMetrics 收集指标
func (c *NetcoreClient) collectMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats := c.GetStats()
			c.logger.Debug("Client metrics",
				"is_connected", stats.IsConnected,
				"total_connections", stats.TotalConnections,
				"total_messages", stats.TotalMessages,
				"total_errors", stats.TotalErrors,
				"reconnect_count", stats.ReconnectCount,
				"uptime", stats.Uptime)
		case <-c.ctx.Done():
			return
		}
	}
}

// updateStats 更新统计信息
func (c *NetcoreClient) updateStats(msgType uint32, success bool, received int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if success {
		if received > 0 {
			c.stats.TotalMessages++
		}
	} else {
		c.stats.TotalErrors++
	}

	// 更新消息类型统计
	msgStats, exists := c.stats.ByMessageType[msgType]
	if !exists {
		msgStats = &ClientMessageTypeStats{}
		c.stats.ByMessageType[msgType] = msgStats
	}

	if success {
		if received > 0 {
			msgStats.ReceivedCount++
			msgStats.LastReceived = time.Now()
		} else {
			msgStats.SentCount++
			msgStats.LastSent = time.Now()
		}
	} else {
		msgStats.FailedCount++
	}
}

// 统计信息结构
type ClientStats struct {
	IsConnected      bool                               `json:"is_connected"`
	TotalConnections int64                              `json:"total_connections"`
	TotalMessages    int64                              `json:"total_messages"`
	TotalErrors      int64                              `json:"total_errors"`
	ReconnectCount   int64                              `json:"reconnect_count"`
	ConnectTime      time.Time                          `json:"connect_time"`
	StartTime        time.Time                          `json:"start_time"`
	Uptime           time.Duration                      `json:"uptime"`
	ByMessageType    map[uint32]*ClientMessageTypeStats `json:"by_message_type"`
}

type ClientMessageTypeStats struct {
	SentCount      int64         `json:"sent_count"`
	ReceivedCount  int64         `json:"received_count"`
	FailedCount    int64         `json:"failed_count"`
	LastSent       time.Time     `json:"last_sent"`
	LastReceived   time.Time     `json:"last_received"`
	AvgProcessTime time.Duration `json:"avg_process_time"`
}
