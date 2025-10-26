package tcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	appHandlers "greatestworks/internal/application/handlers"
	appServices "greatestworks/internal/application/services"
	"greatestworks/internal/domain/character"
	"greatestworks/internal/infrastructure/logging"
	"greatestworks/internal/interfaces/tcp/connection"
	tcpHandlers "greatestworks/internal/interfaces/tcp/handlers"
	"greatestworks/internal/interfaces/tcp/protocol"
)

// ServerConfig TCP服务器配置
type ServerConfig struct {
	Addr              string
	MaxConnections    int
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	EnableCompression bool
	BufferSize        int
}

// DefaultServerConfig 默认服务器配置
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Addr:              ":9090",
		MaxConnections:    10000,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		EnableCompression: false,
		BufferSize:        4096,
	}
}

// TCPServer TCP服务器
type TCPServer struct {
	config           *ServerConfig
	listener         net.Listener
	gameHandler      *tcpHandlers.GameHandler
	router           *Router
	connManager      *connection.Manager
	heartbeatManager *connection.HeartbeatManager
	logger           logging.Logger
	ctx              context.Context
	cancel           context.CancelFunc
	wg               sync.WaitGroup
	running          bool
	mutex            sync.RWMutex

	// optional references for wiring
	mapService       *appServices.MapService
	fightService     *appServices.FightService
	characterService *appServices.CharacterService
}

// NewTCPServer 创建TCP服务器
func NewTCPServer(config *ServerConfig, commandBus *appHandlers.CommandBus, queryBus *appHandlers.QueryBus, logger logging.Logger) *TCPServer {
	if config == nil {
		config = DefaultServerConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	// 创建连接管理器
	connManager := connection.NewManager(logger)

	// 创建心跳管理器
	heartbeatManager := connection.NewHeartbeatManager(logger, 30*time.Second, 60*time.Second)

	// 创建游戏处理器
	gameHandler := tcpHandlers.NewGameHandler(commandBus, queryBus, connManager, logger)

	// 创建路由器
	router := NewRouter(logger)
	router.RegisterGameHandler(gameHandler)

	server := &TCPServer{
		config:           config,
		gameHandler:      gameHandler,
		router:           router,
		connManager:      connManager,
		heartbeatManager: heartbeatManager,
		logger:           logger,
		ctx:              ctx,
		cancel:           cancel,
		running:          false,
	}

	return server
}

// SetMapService allows injecting MapService for potential handler usage.
func (s *TCPServer) SetMapService(ms *appServices.MapService) {
	s.mapService = ms
	if s.gameHandler != nil {
		s.gameHandler.SetMapService(ms)
	}
}

// SetFightService allows injecting FightService for handler usage.
func (s *TCPServer) SetFightService(fs *appServices.FightService) {
	s.fightService = fs
	if s.gameHandler != nil {
		s.gameHandler.SetFightService(fs)
	}
}

// SetCharacterService allows injecting CharacterService for handler usage.
func (s *TCPServer) SetCharacterService(cs *appServices.CharacterService) {
	s.characterService = cs
	if s.gameHandler != nil {
		s.gameHandler.SetCharacterService(cs)
	}
}

// GetConnectionManager exposes the underlying connection manager for wiring.
func (s *TCPServer) GetConnectionManager() *connection.Manager { return s.connManager }

// Start 启动TCP服务器
func (s *TCPServer) Start() error {
	s.mutex.Lock()
	if s.running {
		s.mutex.Unlock()
		return fmt.Errorf("server is already running")
	}
	s.mutex.Unlock()

	s.logger.Info("Starting TCP server", map[string]interface{}{
		"address": s.config.Addr,
	})

	// 创建监听器
	listener, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		s.logger.Error("Failed to create listener", err, logging.Fields{
			"address": s.config.Addr,
		})
		return fmt.Errorf("failed to create listener: %w", err)
	}

	s.listener = listener
	s.mutex.Lock()
	s.running = true
	s.mutex.Unlock()

	// 启动心跳管理器
	go s.heartbeatManager.Start(s.ctx)

	// 启动接受连接的协程
	s.wg.Add(1)
	go s.acceptConnections()

	s.logger.Info("TCP server started successfully", map[string]interface{}{
		"address": s.config.Addr,
	})
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
			s.logger.Error("Failed to close listener", err)
		}
	}

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
					s.logger.Error("Failed to accept connection", err)
					continue
				}
			}

			// 检查连接数限制
			connectionCount := s.connManager.GetConnectionCount()
			if connectionCount >= s.config.MaxConnections {
				s.logger.Warn("Connection limit reached, rejecting new connection", logging.Fields{
					"current_count":   connectionCount,
					"max_connections": s.config.MaxConnections,
				})
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

	// 创建会话
	session := connection.NewSession(fmt.Sprintf("session_%d", time.Now().UnixNano()), netConn, s.logger)

	// 添加到连接管理器
	s.connManager.AddConnection(session)

	// 添加到心跳管理器
	s.heartbeatManager.AddSession(session)

	s.logger.Info("New connection established", map[string]interface{}{
		"session_id":  session.ID,
		"remote_addr": netConn.RemoteAddr(),
	})

	// 处理连接消息
	defer func() {
		// 清理连接
		s.cleanupConnection(session)
	}()

	// 消息处理循环
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			// 设置读取超时
			netConn.SetReadDeadline(time.Now().Add(s.config.ReadTimeout))

			// 读取消息
			msg, err := s.readMessage(netConn)
			if err != nil {
				if err == io.EOF {
					s.logger.Info("Connection closed by client", map[string]interface{}{
						"session_id": session.ID,
					})
				} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					s.logger.Debug("Read timeout", map[string]interface{}{
						"session_id": session.ID,
					})
					continue
				} else {
					s.logger.Error("Failed to read message", err, logging.Fields{
						"session_id": session.ID,
					})
				}
				return
			}

			// 验证消息
			if err := s.router.ValidateMessage(msg); err != nil {
				s.logger.Error("Invalid message received", err, logging.Fields{
					"session_id": session.ID,
				})
				continue
			}

			// 路由消息
			if err := s.router.RouteMessage(session, msg); err != nil {
				s.logger.Error("Failed to route message", err, logging.Fields{
					"session_id":   session.ID,
					"message_type": msg.Header.MessageType,
				})
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
func (s *TCPServer) cleanupConnection(session *connection.Session) {
	s.logger.Info("Cleaning up connection", map[string]interface{}{
		"session_id": session.ID,
		"user_id":    session.UserID,
	})

	// 从地图中移除与该会话绑定的实体；并保存最后位置
	if s.mapService != nil {
		if entityID, ok := s.connManager.GetPlayerBySession(session.ID); ok {
			var mapID int32 = 1
			if session.GroupID != "" && len(session.GroupID) > 4 && session.GroupID[:4] == "map:" {
				if v, err := strconv.ParseInt(session.GroupID[4:], 10, 32); err == nil {
					mapID = int32(v)
				}
			}
			// 保存位置
			if s.characterService != nil && mapID > 0 {
				if m, err := s.mapService.GetMap(mapID); err == nil && m != nil {
					if e := m.GetEntity(character.EntityID(entityID)); e != nil {
						pos := e.Position()
						_ = s.characterService.UpdateLastLocation(
							s.ctx, int64(entityID), mapID, pos.X, pos.Y, pos.Z,
						)
					}
				}
			}
			_ = s.mapService.LeaveMapByID(s.ctx, mapID, entityID)
		}
	}

	// 从心跳管理器移除
	s.heartbeatManager.RemoveSession(session.ID)

	// 从连接管理器移除
	s.connManager.RemoveConnection(session.ID)

	// 关闭连接
	session.Close()

	s.logger.Debug("Connection cleanup completed", map[string]interface{}{
		"session_id": session.ID,
	})
}

// IsRunning 检查服务器是否运行中
func (s *TCPServer) IsRunning() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.running
}

// GetStats 获取服务器统计信息
func (s *TCPServer) GetStats() map[string]interface{} {
	connectionCount := s.connManager.GetConnectionCount()
	activeSessionCount := s.heartbeatManager.GetActiveSessionCount()

	return map[string]interface{}{
		"running":          s.IsRunning(),
		"address":          s.config.Addr,
		"max_connections":  s.config.MaxConnections,
		"connection_count": connectionCount,
		"active_sessions":  activeSessionCount,
		"router_stats": map[string]interface{}{
			"handler_count": s.router.GetHandlerCount(),
			"message_types": s.router.GetRegisteredMessageTypes(),
		},
	}
}
