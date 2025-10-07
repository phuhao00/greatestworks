package network

import (
	"fmt"
	"net"
	"sync"

	"greatestworks/internal/infrastructure/logging"
)

// NetCoreServer 网络核心服务器
type NetCoreServer struct {
	listener     net.Listener
	clients      map[string]*NetCoreClient
	mutex        sync.RWMutex
	logger       logging.Logger
	port         int
	host         string
	onConnect    func(*NetCoreClient)
	onDisconnect func(*NetCoreClient)
	onMessage    func(*NetCoreClient, []byte)
}

// NewNetCoreServer 创建网络核心服务器
func NewNetCoreServer(host string, port int, logger logging.Logger) *NetCoreServer {
	return &NetCoreServer{
		clients: make(map[string]*NetCoreClient),
		logger:  logger,
		port:    port,
		host:    host,
	}
}

// SetOnConnect 设置连接回调
func (s *NetCoreServer) SetOnConnect(callback func(*NetCoreClient)) {
	s.onConnect = callback
}

// SetOnDisconnect 设置断开连接回调
func (s *NetCoreServer) SetOnDisconnect(callback func(*NetCoreClient)) {
	s.onDisconnect = callback
}

// SetOnMessage 设置消息回调
func (s *NetCoreServer) SetOnMessage(callback func(*NetCoreClient, []byte)) {
	s.onMessage = callback
}

// Start 启动服务器
func (s *NetCoreServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("启动服务器失败: %w", err)
	}

	s.listener = listener

	s.logger.Info("Network server started", logging.Fields{
		"address": addr,
	})

	// 启动接受连接循环
	go s.acceptConnections()

	return nil
}

// Stop 停止服务器
func (s *NetCoreServer) Stop() error {
	if s.listener == nil {
		return nil
	}

	s.logger.Info("Network server stopped")

	// 关闭所有客户端连接
	s.mutex.Lock()
	for _, client := range s.clients {
		client.Close()
	}
	s.clients = make(map[string]*NetCoreClient)
	s.mutex.Unlock()

	return s.listener.Close()
}

// acceptConnections 接受连接
func (s *NetCoreServer) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.Error("Failed to accept connection", err)
			continue
		}

		// 创建客户端
		client := NewNetCoreClient(conn, s.logger)
		clientID := conn.RemoteAddr().String()

		// 添加到客户端列表
		s.mutex.Lock()
		s.clients[clientID] = client
		s.mutex.Unlock()

		s.logger.Info("新客户端连接", logging.Fields{
			"client_id":   clientID,
			"remote_addr": conn.RemoteAddr().String(),
		})

		// 调用连接回调
		if s.onConnect != nil {
			s.onConnect(client)
		}

		// 启动客户端处理循环
		go s.handleClient(client)
	}
}

// handleClient 处理客户端
func (s *NetCoreServer) handleClient(client *NetCoreClient) {
	defer func() {
		// 从客户端列表移除
		s.mutex.Lock()
		delete(s.clients, client.GetRemoteAddr())
		s.mutex.Unlock()

		// 调用断开连接回调
		if s.onDisconnect != nil {
			s.onDisconnect(client)
		}

		// 关闭客户端连接
		client.Close()
	}()

	for {
		// 接收消息
		data, err := client.Receive()
		if err != nil {
			s.logger.Error("Failed to receive message", err, logging.Fields{
				"client_id": client.GetRemoteAddr(),
			})
			break
		}

		// 调用消息回调
		if s.onMessage != nil {
			s.onMessage(client, data)
		}
	}
}

// GetClientCount 获取客户端数量
func (s *NetCoreServer) GetClientCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.clients)
}

// GetClients 获取所有客户端
func (s *NetCoreServer) GetClients() map[string]*NetCoreClient {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	clients := make(map[string]*NetCoreClient)
	for id, client := range s.clients {
		clients[id] = client
	}

	return clients
}

// Broadcast 广播消息
func (s *NetCoreServer) Broadcast(message []byte) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, client := range s.clients {
		if err := client.Send(message); err != nil {
			s.logger.Error("Failed to broadcast message", err, logging.Fields{
				"client_id": client.GetRemoteAddr(),
			})
		}
	}
}

// SendToClient 发送消息给指定客户端
func (s *NetCoreServer) SendToClient(clientID string, message []byte) error {
	s.mutex.RLock()
	client, exists := s.clients[clientID]
	s.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("客户端不存在: %s", clientID)
	}

	return client.Send(message)
}

// GetClient 获取客户端
func (s *NetCoreServer) GetClient(clientID string) (*NetCoreClient, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	client, exists := s.clients[clientID]
	return client, exists
}
