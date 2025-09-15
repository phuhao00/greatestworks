package tcp

import (
	"context"
	"fmt"
	"time"
	
	"greatestworks/aop/logger"
	"greatestworks/application/services"
	"greatestworks/internal/infrastructure/network"
)

// TCPServer TCP服务器
type TCPServer struct {
	server         network.Server
	playerHandler  *PlayerHandler
	sceneHandler   *SceneHandler
	npcHandler     *NPCHandler
	logger         logger.Logger
	config         *TCPServerConfig
}

// TCPServerConfig TCP服务器配置
type TCPServerConfig struct {
	Host                string        `json:"host" yaml:"host"`
	Port                int           `json:"port" yaml:"port"`
	MaxConnections      int           `json:"max_connections" yaml:"max_connections"`
	ReadTimeout         time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout        time.Duration `json:"write_timeout" yaml:"write_timeout"`
	KeepAliveInterval   time.Duration `json:"keep_alive_interval" yaml:"keep_alive_interval"`
	HeartbeatInterval   time.Duration `json:"heartbeat_interval" yaml:"heartbeat_interval"`
	MaxPacketSize       int           `json:"max_packet_size" yaml:"max_packet_size"`
	CompressionEnabled  bool          `json:"compression_enabled" yaml:"compression_enabled"`
	EncryptionEnabled   bool          `json:"encryption_enabled" yaml:"encryption_enabled"`
	EnableMetrics       bool          `json:"enable_metrics" yaml:"enable_metrics"`
}

// ServiceContainer 服务容器
type ServiceContainer struct {
	PlayerService  *services.PlayerService
	HangupService  *services.HangupService
	WeatherService *services.WeatherService
	PlantService   *services.PlantService
	NPCService     *services.NPCService
}

// NewTCPServer 创建TCP服务器
func NewTCPServer(config *TCPServerConfig, serviceContainer *ServiceContainer, logger logger.Logger) (*TCPServer, error) {
	if config == nil {
		config = &TCPServerConfig{
			Host:                "0.0.0.0",
			Port:                8080,
			MaxConnections:      1000,
			ReadTimeout:         30 * time.Second,
			WriteTimeout:        30 * time.Second,
			KeepAliveInterval:   60 * time.Second,
			HeartbeatInterval:   30 * time.Second,
			MaxPacketSize:       1024 * 1024, // 1MB
			CompressionEnabled:  false,
			EncryptionEnabled:   false,
			EnableMetrics:       true,
		}
	}
	
	// 创建网络服务器配置
	serverConfig := &network.ServerConfig{
		Host:                config.Host,
		Port:                config.Port,
		MaxConnections:      config.MaxConnections,
		ReadTimeout:         config.ReadTimeout,
		WriteTimeout:        config.WriteTimeout,
		KeepAliveInterval:   config.KeepAliveInterval,
		HeartbeatInterval:   config.HeartbeatInterval,
		MaxPacketSize:       config.MaxPacketSize,
		CompressionEnabled:  config.CompressionEnabled,
		EncryptionEnabled:   config.EncryptionEnabled,
		EnableMetrics:       config.EnableMetrics,
	}
	
	// 创建网络服务器
	server := network.NewNetcoreServer(serverConfig, logger)
	
	// 创建处理器
	playerHandler := NewPlayerHandler(serviceContainer.PlayerService, serviceContainer.HangupService, logger)
	sceneHandler := NewSceneHandler(serviceContainer.WeatherService, serviceContainer.PlantService, logger)
	npcHandler := NewNPCHandler(serviceContainer.NPCService, logger)
	
	tcpServer := &TCPServer{
		server:        server,
		playerHandler: playerHandler,
		sceneHandler:  sceneHandler,
		npcHandler:    npcHandler,
		logger:        logger,
		config:        config,
	}
	
	// 注册所有处理器
	if err := tcpServer.registerHandlers(); err != nil {
		return nil, fmt.Errorf("failed to register handlers: %w", err)
	}
	
	// 设置连接处理器
	tcpServer.setupConnectionHandlers()
	
	logger.Info("TCP server created successfully", "address", fmt.Sprintf("%s:%d", config.Host, config.Port))
	return tcpServer, nil
}

// Start 启动TCP服务器
func (s *TCPServer) Start(ctx context.Context) error {
	s.logger.Info("Starting TCP server", "address", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
	
	// 启动网络服务器
	if err := s.server.Start(ctx); err != nil {
		s.logger.Error("Failed to start TCP server", "error", err)
		return fmt.Errorf("failed to start TCP server: %w", err)
	}
	
	s.logger.Info("TCP server started successfully")
	return nil
}

// Stop 停止TCP服务器
func (s *TCPServer) Stop() error {
	s.logger.Info("Stopping TCP server")
	
	// 停止网络服务器
	if err := s.server.Stop(); err != nil {
		s.logger.Error("Failed to stop TCP server", "error", err)
		return fmt.Errorf("failed to stop TCP server: %w", err)
	}
	
	s.logger.Info("TCP server stopped successfully")
	return nil
}

// GetStats 获取服务器统计信息
func (s *TCPServer) GetStats() *network.ServerStats {
	return s.server.GetStats()
}

// GetConnections 获取所有连接
func (s *TCPServer) GetConnections() []network.Connection {
	return s.server.GetConnections()
}

// BroadcastMessage 广播消息
func (s *TCPServer) BroadcastMessage(msgType uint32, data interface{}) error {
	// 序列化数据
	packetData, err := s.serializeData(data)
	if err != nil {
		s.logger.Error("Failed to serialize broadcast data", "error", err)
		return fmt.Errorf("failed to serialize broadcast data: %w", err)
	}
	
	// 创建数据包
	packet := network.NewPacket(msgType, packetData)
	
	// 广播消息
	if err := s.server.Broadcast(packet); err != nil {
		s.logger.Error("Failed to broadcast message", "error", err, "msg_type", msgType)
		return fmt.Errorf("failed to broadcast message: %w", err)
	}
	
	s.logger.Debug("Message broadcasted successfully", "msg_type", msgType)
	return nil
}

// SendToConnection 发送消息到指定连接
func (s *TCPServer) SendToConnection(connID string, msgType uint32, data interface{}) error {
	// 序列化数据
	packetData, err := s.serializeData(data)
	if err != nil {
		s.logger.Error("Failed to serialize message data", "error", err)
		return fmt.Errorf("failed to serialize message data: %w", err)
	}
	
	// 创建数据包
	packet := network.NewPacket(msgType, packetData)
	
	// 发送消息
	if err := s.server.SendToConnection(connID, packet); err != nil {
		s.logger.Error("Failed to send message to connection", "error", err, "conn_id", connID, "msg_type", msgType)
		return fmt.Errorf("failed to send message to connection: %w", err)
	}
	
	s.logger.Debug("Message sent to connection successfully", "conn_id", connID, "msg_type", msgType)
	return nil
}

// 私有方法

// registerHandlers 注册所有处理器
func (s *TCPServer) registerHandlers() error {
	// 注册玩家处理器
	if err := s.playerHandler.RegisterHandlers(s.server); err != nil {
		return fmt.Errorf("failed to register player handlers: %w", err)
	}
	
	// 注册场景处理器
	if err := s.sceneHandler.RegisterHandlers(s.server); err != nil {
		return fmt.Errorf("failed to register scene handlers: %w", err)
	}
	
	// 注册NPC处理器
	if err := s.npcHandler.RegisterHandlers(s.server); err != nil {
		return fmt.Errorf("failed to register NPC handlers: %w", err)
	}
	
	// 注册系统处理器
	if err := s.registerSystemHandlers(); err != nil {
		return fmt.Errorf("failed to register system handlers: %w", err)
	}
	
	s.logger.Info("All handlers registered successfully")
	return nil
}

// registerSystemHandlers 注册系统处理器
func (s *TCPServer) registerSystemHandlers() error {
	// 注册心跳处理器
	heartbeatHandler := &HeartbeatHandler{logger: s.logger}
	if err := s.server.RegisterHandler(heartbeatHandler); err != nil {
		return fmt.Errorf("failed to register heartbeat handler: %w", err)
	}
	
	// 注册ping处理器
	pingHandler := &PingHandler{logger: s.logger}
	if err := s.server.RegisterHandler(pingHandler); err != nil {
		return fmt.Errorf("failed to register ping handler: %w", err)
	}
	
	// 注册错误处理器
	errorHandler := &ErrorHandler{logger: s.logger}
	if err := s.server.RegisterHandler(errorHandler); err != nil {
		return fmt.Errorf("failed to register error handler: %w", err)
	}
	
	s.logger.Debug("System handlers registered successfully")
	return nil
}

// setupConnectionHandlers 设置连接处理器
func (s *TCPServer) setupConnectionHandlers() {
	// 这里可以设置连接建立、断开等事件的处理器
	// 由于当前的network.Server接口设计，这些处理器在创建服务器时已经设置
	s.logger.Debug("Connection handlers setup completed")
}

// serializeData 序列化数据
func (s *TCPServer) serializeData(data interface{}) ([]byte, error) {
	// 如果已经是字节数组，直接返回
	if bytes, ok := data.([]byte); ok {
		return bytes, nil
	}
	
	// 如果是字符串，转换为字节数组
	if str, ok := data.(string); ok {
		return []byte(str), nil
	}
	
	// 否则使用JSON序列化
	return json.Marshal(data)
}

// 系统处理器

// HeartbeatHandler 心跳处理器
type HeartbeatHandler struct {
	logger logger.Logger
}

func (h *HeartbeatHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	h.logger.Debug("Heartbeat received", "conn_id", conn.GetID())
	
	// 回复心跳
	response := []byte("heartbeat_ack")
	replyPacket := netcore.NewPacket(0, response)
	
	return conn.Send(replyPacket)
}

func (h *HeartbeatHandler) GetMessageType() uint32 {
	return 0 // 心跳消息类型
}

func (h *HeartbeatHandler) GetHandlerName() string {
	return "HeartbeatHandler"
}

// PingHandler Ping处理器
type PingHandler struct {
	logger logger.Logger
}

func (h *PingHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	h.logger.Debug("Ping received", "conn_id", conn.GetID())
	
	// 回复Pong
	response := []byte("pong")
	replyPacket := netcore.NewPacket(9998, response)
	
	return conn.Send(replyPacket)
}

func (h *PingHandler) GetMessageType() uint32 {
	return 9998 // Ping消息类型
}

func (h *PingHandler) GetHandlerName() string {
	return "PingHandler"
}

// ErrorHandler 错误处理器
type ErrorHandler struct {
	logger logger.Logger
}

func (h *ErrorHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	h.logger.Error("Error message received", "conn_id", conn.GetID(), "data", string(packet.GetData()))
	
	// 记录错误但不回复
	return nil
}

func (h *ErrorHandler) GetMessageType() uint32 {
	return 9999 // 错误消息类型
}

func (h *ErrorHandler) GetHandlerName() string {
	return "ErrorHandler"
}

// 连接管理器集成

// ConnectionManagerIntegration 连接管理器集成
type ConnectionManagerIntegration struct {
	manager network.Manager
	logger  logger.Logger
}

// NewConnectionManagerIntegration 创建连接管理器集成
func NewConnectionManagerIntegration(manager network.Manager, logger logger.Logger) *ConnectionManagerIntegration {
	return &ConnectionManagerIntegration{
		manager: manager,
		logger:  logger,
	}
}

// OnConnect 连接建立时调用
func (c *ConnectionManagerIntegration) OnConnect(conn *netcore.Connection) error {
	// 添加连接到管理器
	// 这里需要从连接中提取用户ID和组ID，简化处理
	userID := "" // 从连接上下文或认证信息中获取
	groupID := "" // 从连接上下文或认证信息中获取
	
	if _, err := c.manager.AddConnection(conn, userID, groupID); err != nil {
		c.logger.Error("Failed to add connection to manager", "error", err, "conn_id", conn.GetID())
		return err
	}
	
	c.logger.Info("Connection added to manager", "conn_id", conn.GetID())
	return nil
}

// OnDisconnect 连接断开时调用
func (c *ConnectionManagerIntegration) OnDisconnect(conn *netcore.Connection, err error) {
	// 从管理器中移除连接
	if removeErr := c.manager.RemoveConnection(conn.GetID()); removeErr != nil {
		c.logger.Error("Failed to remove connection from manager", "error", removeErr, "conn_id", conn.GetID())
	}
	
	c.logger.Info("Connection removed from manager", "conn_id", conn.GetID(), "disconnect_error", err)
}

// OnError 发生错误时调用
func (c *ConnectionManagerIntegration) OnError(conn *netcore.Connection, err error) {
	c.logger.Error("Connection error occurred", "error", err, "conn_id", conn.GetID())
}

// 服务器工厂

// ServerFactory 服务器工厂
type ServerFactory struct {
	logger logger.Logger
}

// NewServerFactory 创建服务器工厂
func NewServerFactory(logger logger.Logger) *ServerFactory {
	return &ServerFactory{
		logger: logger,
	}
}

// CreateTCPServer 创建TCP服务器
func (f *ServerFactory) CreateTCPServer(config *TCPServerConfig, services *ServiceContainer) (*TCPServer, error) {
	return NewTCPServer(config, services, f.logger)
}

// CreateTCPServerWithDefaults 使用默认配置创建TCP服务器
func (f *ServerFactory) CreateTCPServerWithDefaults(services *ServiceContainer) (*TCPServer, error) {
	return NewTCPServer(nil, services, f.logger)
}