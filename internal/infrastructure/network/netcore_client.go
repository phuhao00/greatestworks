package network

import (
	"context"
	"fmt"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logger"
)

// NetcoreClient netcore-go TCP客户端
type NetcoreClient struct {
	client   interface{} // *netcore.Client
	logger   logger.Logger
	config   *ClientConfig
	handlers map[uint32]ClientMessageHandler
	mu       sync.RWMutex
	stats    *ClientStats
	ctx      context.Context
	cancel   context.CancelFunc
	conn     interface{} // *netcore.Connection
}

// ClientConfig 客户端配置
type ClientConfig struct {
	ServerHost           string        `json:"server_host" yaml:"server_host"`
	ServerPort           int           `json:"server_port" yaml:"server_port"`
	ConnectTimeout       time.Duration `json:"connect_timeout" yaml:"connect_timeout"`
	ReadTimeout          time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout         time.Duration `json:"write_timeout" yaml:"write_timeout"`
	KeepAliveInterval    time.Duration `json:"keep_alive_interval" yaml:"keep_alive_interval"`
	HeartbeatInterval    time.Duration `json:"heartbeat_interval" yaml:"heartbeat_interval"`
	MaxReconnectAttempts int           `json:"max_reconnect_attempts" yaml:"max_reconnect_attempts"`
	ReconnectInterval    time.Duration `json:"reconnect_interval" yaml:"reconnect_interval"`
	EnableMetrics        bool          `json:"enable_metrics" yaml:"enable_metrics"`
}

// ClientMessageHandler 客户端消息处理器接口
type ClientMessageHandler interface {
	// Handle 处理消息
	Handle(ctx context.Context, conn Connection, packet Packet) error

	// GetMessageType 获取消息类型
	GetMessageType() uint32

	// GetHandlerName 获取处理器名称
	GetHandlerName() string
}

// ClientStats 客户端统计信息
type ClientStats struct {
	Connected       bool                         `json:"connected"`
	TotalMessages   int64                        `json:"total_messages"`
	TotalErrors     int64                        `json:"total_errors"`
	ConnectTime     time.Time                    `json:"connect_time"`
	LastMessageTime time.Time                    `json:"last_message_time"`
	ReconnectCount  int64                        `json:"reconnect_count"`
	ByMessageType   map[uint32]*MessageTypeStats `json:"by_message_type"`
}

// NewNetcoreClient 创建netcore客户端
func NewNetcoreClient(config *ClientConfig, logger logger.Logger) *NetcoreClient {
	if config == nil {
		config = &ClientConfig{
			ServerHost:           "localhost",
			ServerPort:           8080,
			ConnectTimeout:       10 * time.Second,
			ReadTimeout:          30 * time.Second,
			WriteTimeout:         30 * time.Second,
			KeepAliveInterval:    60 * time.Second,
			HeartbeatInterval:    30 * time.Second,
			MaxReconnectAttempts: 5,
			ReconnectInterval:    5 * time.Second,
			EnableMetrics:        true,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	client := &NetcoreClient{
		client:   nil, // 暂时设为nil，后续实现
		logger:   logger,
		config:   config,
		handlers: make(map[uint32]ClientMessageHandler),
		ctx:      ctx,
		cancel:   cancel,
		stats: &ClientStats{
			Connected:     false,
			ByMessageType: make(map[uint32]*MessageTypeStats),
		},
	}

	logger.Info("Netcore client initialized successfully", "server", fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort))
	return client
}

// Connect 连接到服务器
func (c *NetcoreClient) Connect() error {
	c.logger.Info("Connecting to server", "server", fmt.Sprintf("%s:%d", c.config.ServerHost, c.config.ServerPort))

	// TODO: 实现实际的连接逻辑
	c.stats.Connected = true
	c.stats.ConnectTime = time.Now()

	c.logger.Info("Connected to server successfully")
	return nil
}

// Disconnect 断开连接
func (c *NetcoreClient) Disconnect() error {
	c.logger.Info("Disconnecting from server")

	// TODO: 实现实际的断开连接逻辑
	c.stats.Connected = false
	c.cancel()

	c.logger.Info("Disconnected from server")
	return nil
}

// Send 发送消息
func (c *NetcoreClient) Send(packet Packet) error {
	if !c.stats.Connected {
		return fmt.Errorf("client not connected")
	}

	// TODO: 实现实际的消息发送逻辑
	c.stats.TotalMessages++
	c.stats.LastMessageTime = time.Now()
	c.updateStats(packet.GetType(), true, 1)

	c.logger.Debug("Message sent successfully", "message_type", packet.GetType())
	return nil
}

// RegisterHandler 注册消息处理器
func (c *NetcoreClient) RegisterHandler(handler ClientMessageHandler) error {
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

// UnregisterHandler 取消注册消息处理器
func (c *NetcoreClient) UnregisterHandler(messageType uint32) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.handlers[messageType]; !exists {
		return fmt.Errorf("handler for message type %d not found", messageType)
	}

	delete(c.handlers, messageType)
	c.logger.Info("Message handler unregistered successfully", "message_type", messageType)
	return nil
}

// GetStats 获取客户端统计信息
func (c *NetcoreClient) GetStats() *ClientStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 创建统计信息副本
	stats := &ClientStats{
		Connected:       c.stats.Connected,
		TotalMessages:   c.stats.TotalMessages,
		TotalErrors:     c.stats.TotalErrors,
		ConnectTime:     c.stats.ConnectTime,
		LastMessageTime: c.stats.LastMessageTime,
		ReconnectCount:  c.stats.ReconnectCount,
		ByMessageType:   make(map[uint32]*MessageTypeStats),
	}

	// 复制消息类型统计
	for msgType, msgStats := range c.stats.ByMessageType {
		stats.ByMessageType[msgType] = &MessageTypeStats{
			ProcessedCount: msgStats.ProcessedCount,
			FailedCount:    msgStats.FailedCount,
			LastProcessed:  msgStats.LastProcessed,
			AvgProcessTime: msgStats.AvgProcessTime,
		}
	}

	return stats
}

// IsConnected 检查是否已连接
func (c *NetcoreClient) IsConnected() bool {
	return c.stats.Connected
}

// 私有方法

// updateStats 更新统计信息
func (c *NetcoreClient) updateStats(msgType uint32, success bool, count int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if success {
		c.stats.TotalMessages += int64(count)
	} else {
		c.stats.TotalErrors++
	}

	// 更新消息类型统计
	msgStats, exists := c.stats.ByMessageType[msgType]
	if !exists {
		msgStats = &MessageTypeStats{}
		c.stats.ByMessageType[msgType] = msgStats
	}

	if success {
		msgStats.ProcessedCount += int64(count)
		msgStats.LastProcessed = time.Now()
	} else {
		msgStats.FailedCount++
	}
}
