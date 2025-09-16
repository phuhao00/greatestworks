package tcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"greatestworks/application/handlers"
	"greatestworks/internal/infrastructure/logger"
	"greatestworks/internal/interfaces/tcp/connection"
	"greatestworks/internal/interfaces/tcp/handlers"
	"greatestworks/internal/interfaces/tcp/protocol"
)

// ServerConfig TCP服务器配置
type ServerConfig struct {
	Addr               string
	MaxConnections     int
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	HeartbeatConfig    *connection.HeartbeatConfig
	EnableCompression  bool
	BufferSize         int
}

// DefaultServerConfig 默认服务器配置
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Addr:               ":9090",
		MaxConnections:     10000,
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		HeartbeatConfig:    connection.DefaultHeartbeatConfig(),
		EnableCompression:  false,
		BufferSize:         4096,
	}
}

// TCPServer TCP服务器
type TCPServer struct {
	config           *ServerConfig
	listener         net.Listener
	gameHandler      *handlers.GameHandler
	router           *Router
	connManager      *connection.ConnectionManager
	sessionManager   *connection.SessionManager
	heartbeatManager *connection.HeartbeatManager
	logger           logger.Logger
	ctx              context.Context
	cancel           context.CancelFunc
	wg               sync.WaitGroup
	running          bool
	mutex            sync.RWMutex
}

// NewTCPServer 创建TCP服务器
func NewTCPServer(config *ServerConfig, commandBus *handlers.CommandBus, queryBus *handlers.QueryBus, logger logger.Logger) *TCPServer {
	if config == nil {
		config = DefaultServerConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	// 创建连接管理器
	connManager := connection.NewConnectionManager(config.MaxConnections, logger)

	// 创建会话管理器
	sessionManager := connection.NewSessionManager(logger)

	// 创建心跳管理器
	heartbeatManager := connection.NewHeartbeatManager(config.HeartbeatConfig, logger)

	// 创建游戏处理器
	gameHandler := handlers.NewGameHandler(commandBus, queryBus, connManager, logger)

	// 创建路由器
	router := NewRouter(logger)
	router.RegisterGameHandler(gameHandler)

	server := &TCPServer{
		config:           config,
		gameHandler:      gameHandler,
		router:           router,
		connManager:      connManager,
		sessionManager:   sessionManager,
		heartbeatManager: heartbeatManager,
		logger:           logger,
		ctx:              ctx,
		cancel:           cancel,
		running:          false,
	}

	// 设置心跳管理器的断线回调
	heartbeatManager.SetDisconnectCallback(server.handleConnectionDisconnect)

	return server
}

// Start 启动TCP服务器
func (s *TCPServer) Start() error {
	s.mutex.Lock()
	if s.running {
		s.mutex.Unlock()
		return fmt.Errorf("server is already running")
	}
	s.mutex.Unlock()

	s.logger.Info("Starting TCP server", "address", s.config.Addr)

	// 创建监听器
	listener, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		s.logger.Error("Failed to create listener", "error", err, "address", s.config.Addr)
		return fmt.Errorf("failed to create listener: %w", err)
	}

	s.listener = listener
	s.mutex.Lock()
	s.running = true
	s.mutex.Unlock()

	// 启动接受连接的协程
	s.wg.Add(1)
	go s.acceptConnections()

	s.logger.Info("TCP server started successfully", "address", s.config.Addr)
	return nil
}

// Stop 停止TCP服务器
func (s *TCPServer) Stop() error {
	s.mutex.Lock()
	if !s.running {
		s.mutex.Unlock()
		return fmt.Errorf("server is not running")
	}
	s.running = false
	s.mutex.Unlock()

	s.logger.Info("Stopping TCP server")

	// 取消上下文
	s.cancel()

	// 关闭监听器
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			s.logger.Error("Failed to close listener", "error", err)
		}
	}

	// 停止心跳管理器
	s.heartbeatManager.Stop()

	// 停止会话管理器
	s.sessionManager.Stop()

	// 关闭所有连接
	s.connManager.CloseAllConnections()

	// 等待所有协程结束
	s.wg.Wait()

	s.logger.Info("TCP server stopped successfully")
	return nil
}

// acceptConnections 接受连接
func (s *TCPServer) acceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("Accept connections routine stopped")
			return
		default:
			// 设置接受超时
			if tcpListener, ok := s.listener.(*net.TCPListener); ok {
				tcpListener.SetDeadline(time.Now().Add(1 * time.Second))
			}

			conn, err := s.listener.Accept()
			if err != nil {
				// 检查是否是超时错误
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}

				// 检查服务器是否正在关闭
				select {
				case <-s.ctx.Done():
					return
				default:
					s.logger.Error("Failed to accept connection", "error", err)
					continue
				}
			}

			// 检查连接数限制
			if s.connManager.GetConnectionCount() >= s.config.MaxConnections {
				s.logger.Warn("Connection limit reached, rejecting new connection", 
					"current_count", s.connManager.GetConnectionCount(),
					"max_connections", s.config.MaxConnections)
				conn.Close()
				continue
			}

			// 处理新连接
			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}
}

// handleConnection 处理连接
func (s *TCPServer) handleConnection(netConn net.Conn) {
	defer s.wg.Done()

	// 设置连接超时
	if tcpConn, ok := netConn.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)
	}

	// 创建连接对象
	conn := connection.NewConnection(netConn, s.config.BufferSize, s.logger)

	// 添加到连接管理器
	if err := s.connManager.AddConnection(conn); err != nil {
		s.logger.Error("Failed to add connection to manager", "error", err, "conn_id", conn.ID)
		conn.Close()
		return
	}

	// 添加到心跳管理器
	s.heartbeatManager.AddConnection(conn)

	s.logger.Info("New connection established", "conn_id", conn.ID, "remote_addr", netConn.RemoteAddr())

	// 处理连接消息
	defer func() {
		// 清理连接
		s.cleanupConnection(conn)
	}()

	// 消息处理循环
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-conn.Done():
			return
		default:
			// 设置读取超时
			netConn.SetReadDeadline(time.Now().Add(s.config.ReadTimeout))

			// 读取消息
			msg, err := s.readMessage(netConn)
			if err != nil {
				if err == io.EOF {
					s.logger.Info("Connection closed by client", "conn_id", conn.ID)
				} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					s.logger.Debug("Read timeout", "conn_id", conn.ID)
					continue
				} else {
					s.logger.Error("Failed to read message", "error", err, "conn_id", conn.ID)
				}
				return
			}

			// 验证消息
			if err := s.router.ValidateMessage(msg); err != nil {
				s.logger.Error("Invalid message received", "error", err, "conn_id", conn.ID)
				continue
			}

			// 路由消息
			if err := s.router.RouteMessage(conn, msg); err != nil {
				s.logger.Error("Failed to route message", "error", err, "conn_id", conn.ID, "message_type", msg.Header.MessageType)
			}
		}
	}
}

// readMessage 读取消息
func (s *TCPServer) readMessage(conn net.Conn) (*protocol.Message, error) {
	// 读取消息头
	headerBytes := make([]byte, protocol.MessageHeaderSize)
	if _, err := io.ReadFull(conn, headerBytes); err != nil {
		return nil, err
	}

	// 解析消息头
	header, err := protocol.ParseMessageHeader(headerBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse message header: %w", err)
	}

	// 读取消息体
	var payload interface{}
	if header.Length > 0 {
		payloadBytes := make([]byte, header.Length)
		if _, err := io.ReadFull(conn, payloadBytes); err != nil {
			return nil, fmt.Errorf("failed to read message payload: %w", err)
		}

		// 解析消息体
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message payload: %w", err)
		}
	}

	return &protocol.Message{
		Header:  *header,
		Payload: payload,
	}, nil
}

// cleanupConnection 清理连接
func (s *TCPServer) cleanupConnection(conn *connection.Connection) {
	s.logger.Info("Cleaning up connection", "conn_id", conn.ID, "player_id", conn.PlayerID)

	// 从心跳管理器移除
	s.heartbeatManager.RemoveConnection(conn.ID)

	// 从会话管理器移除
	if conn.PlayerID != "" {
		s.sessionManager.UnbindPlayerFromSession(conn.PlayerID)
	}

	// 从连接管理器移除
	s.connManager.RemoveConnection(conn.ID)

	// 关闭连接
	conn.Close()

	s.logger.Debug("Connection cleanup completed", "conn_id", conn.ID)
}

// handleConnectionDisconnect 处理连接断开
func (s *TCPServer) handleConnectionDisconnect(connID string) {
	s.logger.Info("Handling connection disconnect", "conn_id", connID)

	// 获取连接
	conn, exists := s.connManager.GetConnection(connID)
	if !exists {
		s.logger.Warn("Connection not found for disconnect handling", "conn_id", connID)
		return
	}

	// 清理连接
	s.cleanupConnection(conn)
}

// IsRunning 检查服务器是否运行中
func (s *TCPServer) IsRunning() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.running
}

// GetStats 获取服务器统计信息
func (s *TCPServer) GetStats() map[string]interface{} {
	connStats := s.connManager.GetStats()
	sessionStats := s.sessionManager.GetStats()
	heartbeatStats := s.heartbeatManager.GetStats()

	return map[string]interface{}{
		"running":           s.IsRunning(),
		"address":           s.config.Addr,
		"max_connections":   s.config.MaxConnections,
		"connection_stats":  connStats,
		"session_stats":     sessionStats,
		"heartbeat_stats":   heartbeatStats,
		"router_stats": map[string]interface{}{
			"handler_count": s.router.GetHandlerCount(),
			"message_types": s.router.GetRegisteredMessageTypes(),
		},
	}
}

// BroadcastMessage 广播消息
func (s *TCPServer) BroadcastMessage(msg *protocol.Message) error {
	return s.connManager.BroadcastMessage(msg)
}

// SendToPlayer 发送消息给指定玩家
func (s *TCPServer) SendToPlayer(playerID string, msg *protocol.Message) error {
	conn, exists := s.connManager.GetConnectionByPlayer(playerID)
	if !exists {
		return fmt.Errorf("player connection not found: %s", playerID)
	}

	return conn.SendMessage(msg)
}

// GetConnectionCount 获取连接数
func (s *TCPServer) GetConnectionCount() int {
	return s.connManager.GetConnectionCount()
}

// GetActiveSessionCount 获取活跃会话数
func (s *TCPServer) GetActiveSessionCount() int {
	return s.sessionManager.GetActiveSessionCount()
}