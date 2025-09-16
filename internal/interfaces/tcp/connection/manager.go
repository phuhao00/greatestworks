package connection

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logger"
	"greatestworks/internal/interfaces/tcp/protocol"
)

// ConnectionManager 连接管理器
type ConnectionManager struct {
	connections    map[string]*Connection
	playerConns    map[string]string // playerID -> connectionID
	mutex          sync.RWMutex
	logger         logger.Logger
	maxConnections int
	heartbeatInterval time.Duration
	readTimeout    time.Duration
	writeTimeout   time.Duration
	cleanupTicker  *time.Ticker
	ctx            context.Context
	cancel         context.CancelFunc
	stats          *ConnectionStats
}

// Connection 连接信息
type Connection struct {
	ID            string
	Conn          net.Conn
	PlayerID      string
	SessionID     string
	RemoteAddr    string
	ConnectedAt   time.Time
	LastHeartbeat time.Time
	LastActivity  time.Time
	IsAuthenticated bool
	Attributes    map[string]interface{}
	SendChan      chan *protocol.Message
	ReceiveChan   chan *protocol.Message
	CloseChan     chan struct{}
	mutex         sync.RWMutex
	logger        logger.Logger
}

// ConnectionStats 连接统计
type ConnectionStats struct {
	TotalConnections    int64
	ActiveConnections   int64
	PeakConnections     int64
	TotalMessages       int64
	MessagesPerSecond   float64
	BytesReceived       int64
	BytesSent           int64
	ConnectionsAccepted int64
	ConnectionsRejected int64
	ConnectionsClosed   int64
	mutex               sync.RWMutex
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager(maxConnections int, heartbeatInterval, readTimeout, writeTimeout time.Duration, logger logger.Logger) *ConnectionManager {
	ctx, cancel := context.WithCancel(context.Background())
	
	cm := &ConnectionManager{
		connections:       make(map[string]*Connection),
		playerConns:       make(map[string]string),
		logger:            logger,
		maxConnections:    maxConnections,
		heartbeatInterval: heartbeatInterval,
		readTimeout:       readTimeout,
		writeTimeout:      writeTimeout,
		ctx:               ctx,
		cancel:            cancel,
		stats:             &ConnectionStats{},
	}

	// 启动清理协程
	cm.cleanupTicker = time.NewTicker(30 * time.Second)
	go cm.cleanupRoutine()

	return cm
}

// AddConnection 添加连接
func (cm *ConnectionManager) AddConnection(conn net.Conn) (*Connection, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// 检查连接数限制
	if len(cm.connections) >= cm.maxConnections {
		cm.stats.mutex.Lock()
		cm.stats.ConnectionsRejected++
		cm.stats.mutex.Unlock()
		return nil, fmt.Errorf("maximum connections reached: %d", cm.maxConnections)
	}

	// 生成连接ID
	connID := cm.generateConnectionID()

	// 创建连接对象
	connection := &Connection{
		ID:              connID,
		Conn:            conn,
		RemoteAddr:      conn.RemoteAddr().String(),
		ConnectedAt:     time.Now(),
		LastHeartbeat:   time.Now(),
		LastActivity:    time.Now(),
		IsAuthenticated: false,
		Attributes:      make(map[string]interface{}),
		SendChan:        make(chan *protocol.Message, 100),
		ReceiveChan:     make(chan *protocol.Message, 100),
		CloseChan:       make(chan struct{}),
		logger:          cm.logger,
	}

	// 添加到连接映射
	cm.connections[connID] = connection

	// 更新统计
	cm.stats.mutex.Lock()
	cm.stats.TotalConnections++
	cm.stats.ActiveConnections++
	cm.stats.ConnectionsAccepted++
	if cm.stats.ActiveConnections > cm.stats.PeakConnections {
		cm.stats.PeakConnections = cm.stats.ActiveConnections
	}
	cm.stats.mutex.Unlock()

	// 启动连接处理协程
	go connection.handleConnection(cm)

	cm.logger.Info("Connection added", "conn_id", connID, "remote_addr", connection.RemoteAddr, "active_connections", len(cm.connections))
	return connection, nil
}

// RemoveConnection 移除连接
func (cm *ConnectionManager) RemoveConnection(connID string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	connection, exists := cm.connections[connID]
	if !exists {
		return
	}

	// 从玩家连接映射中移除
	if connection.PlayerID != "" {
		delete(cm.playerConns, connection.PlayerID)
	}

	// 关闭连接
	connection.Close()

	// 从连接映射中移除
	delete(cm.connections, connID)

	// 更新统计
	cm.stats.mutex.Lock()
	cm.stats.ActiveConnections--
	cm.stats.ConnectionsClosed++
	cm.stats.mutex.Unlock()

	cm.logger.Info("Connection removed", "conn_id", connID, "player_id", connection.PlayerID, "active_connections", len(cm.connections))
}

// GetConnection 获取连接
func (cm *ConnectionManager) GetConnection(connID string) (*Connection, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	connection, exists := cm.connections[connID]
	return connection, exists
}

// GetConnectionByPlayer 根据玩家ID获取连接
func (cm *ConnectionManager) GetConnectionByPlayer(playerID string) (*Connection, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	connID, exists := cm.playerConns[playerID]
	if !exists {
		return nil, false
	}

	connection, exists := cm.connections[connID]
	return connection, exists
}

// BindPlayerToConnection 绑定玩家到连接
func (cm *ConnectionManager) BindPlayerToConnection(connID, playerID string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	connection, exists := cm.connections[connID]
	if !exists {
		return fmt.Errorf("connection not found: %s", connID)
	}

	// 检查玩家是否已经有其他连接
	if existingConnID, exists := cm.playerConns[playerID]; exists {
		// 关闭旧连接
		if existingConn, exists := cm.connections[existingConnID]; exists {
			cm.logger.Warn("Player has existing connection, closing old connection", "player_id", playerID, "old_conn_id", existingConnID, "new_conn_id", connID)
			existingConn.Close()
			delete(cm.connections, existingConnID)
		}
	}

	// 绑定玩家到新连接
	connection.PlayerID = playerID
	connection.IsAuthenticated = true
	cm.playerConns[playerID] = connID

	cm.logger.Info("Player bound to connection", "player_id", playerID, "conn_id", connID)
	return nil
}

// UnbindPlayerFromConnection 解绑玩家和连接
func (cm *ConnectionManager) UnbindPlayerFromConnection(playerID string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	connID, exists := cm.playerConns[playerID]
	if !exists {
		return
	}

	if connection, exists := cm.connections[connID]; exists {
		connection.PlayerID = ""
		connection.IsAuthenticated = false
	}

	delete(cm.playerConns, playerID)
	cm.logger.Info("Player unbound from connection", "player_id", playerID, "conn_id", connID)
}

// BroadcastMessage 广播消息
func (cm *ConnectionManager) BroadcastMessage(msg *protocol.Message) {
	cm.mutex.RLock()
	connections := make([]*Connection, 0, len(cm.connections))
	for _, conn := range cm.connections {
		if conn.IsAuthenticated {
			connections = append(connections, conn)
		}
	}
	cm.mutex.RUnlock()

	for _, conn := range connections {
		select {
		case conn.SendChan <- msg:
			// 消息发送成功
		default:
			// 发送缓冲区满，跳过此连接
			cm.logger.Warn("Failed to broadcast message, send buffer full", "conn_id", conn.ID, "player_id", conn.PlayerID)
		}
	}

	cm.logger.Debug("Message broadcasted", "message_type", msg.Header.MessageType, "recipients", len(connections))
}

// SendToPlayer 发送消息给指定玩家
func (cm *ConnectionManager) SendToPlayer(playerID string, msg *protocol.Message) error {
	connection, exists := cm.GetConnectionByPlayer(playerID)
	if !exists {
		return fmt.Errorf("player not connected: %s", playerID)
	}

	return connection.SendMessage(msg)
}

// GetAllConnections 获取所有连接
func (cm *ConnectionManager) GetAllConnections() []*Connection {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	connections := make([]*Connection, 0, len(cm.connections))
	for _, conn := range cm.connections {
		connections = append(connections, conn)
	}

	return connections
}

// GetOnlinePlayerIDs 获取所有在线玩家ID
func (cm *ConnectionManager) GetOnlinePlayerIDs() []string {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	playerIDs := make([]string, 0, len(cm.playerConns))
	for playerID := range cm.playerConns {
		playerIDs = append(playerIDs, playerID)
	}

	return playerIDs
}

// GetStats 获取连接统计
func (cm *ConnectionManager) GetStats() *ConnectionStats {
	cm.stats.mutex.RLock()
	defer cm.stats.mutex.RUnlock()

	// 复制统计数据
	stats := &ConnectionStats{
		TotalConnections:    cm.stats.TotalConnections,
		ActiveConnections:   cm.stats.ActiveConnections,
		PeakConnections:     cm.stats.PeakConnections,
		TotalMessages:       cm.stats.TotalMessages,
		MessagesPerSecond:   cm.stats.MessagesPerSecond,
		BytesReceived:       cm.stats.BytesReceived,
		BytesSent:           cm.stats.BytesSent,
		ConnectionsAccepted: cm.stats.ConnectionsAccepted,
		ConnectionsRejected: cm.stats.ConnectionsRejected,
		ConnectionsClosed:   cm.stats.ConnectionsClosed,
	}

	return stats
}

// Close 关闭连接管理器
func (cm *ConnectionManager) Close() {
	cm.cancel()
	cm.cleanupTicker.Stop()

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// 关闭所有连接
	for _, conn := range cm.connections {
		conn.Close()
	}

	cm.connections = make(map[string]*Connection)
	cm.playerConns = make(map[string]string)

	cm.logger.Info("Connection manager closed")
}

// 私有方法

// generateConnectionID 生成连接ID
func (cm *ConnectionManager) generateConnectionID() string {
	return fmt.Sprintf("conn_%d_%d", time.Now().UnixNano(), len(cm.connections))
}

// cleanupRoutine 清理协程
func (cm *ConnectionManager) cleanupRoutine() {
	for {
		select {
		case <-cm.ctx.Done():
			return
		case <-cm.cleanupTicker.C:
			cm.cleanupInactiveConnections()
		}
	}
}

// cleanupInactiveConnections 清理非活跃连接
func (cm *ConnectionManager) cleanupInactiveConnections() {
	now := time.Now()
	inactiveTimeout := 5 * time.Minute
	heartbeatTimeout := 2 * cm.heartbeatInterval

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	var toRemove []string

	for connID, conn := range cm.connections {
		// 检查心跳超时
		if now.Sub(conn.LastHeartbeat) > heartbeatTimeout {
			cm.logger.Warn("Connection heartbeat timeout", "conn_id", connID, "player_id", conn.PlayerID, "last_heartbeat", conn.LastHeartbeat)
			toRemove = append(toRemove, connID)
			continue
		}

		// 检查活动超时
		if now.Sub(conn.LastActivity) > inactiveTimeout {
			cm.logger.Warn("Connection inactive timeout", "conn_id", connID, "player_id", conn.PlayerID, "last_activity", conn.LastActivity)
			toRemove = append(toRemove, connID)
			continue
		}
	}

	// 移除超时连接
	for _, connID := range toRemove {
		if conn, exists := cm.connections[connID]; exists {
			if conn.PlayerID != "" {
				delete(cm.playerConns, conn.PlayerID)
			}
			conn.Close()
			delete(cm.connections, connID)

			// 更新统计
			cm.stats.mutex.Lock()
			cm.stats.ActiveConnections--
			cm.stats.ConnectionsClosed++
			cm.stats.mutex.Unlock()
		}
	}

	if len(toRemove) > 0 {
		cm.logger.Info("Cleaned up inactive connections", "removed_count", len(toRemove), "active_connections", len(cm.connections))
	}
}

// Connection 方法

// SendMessage 发送消息
func (c *Connection) SendMessage(msg *protocol.Message) error {
	select {
	case c.SendChan <- msg:
		return nil
	default:
		return fmt.Errorf("send buffer full")
	}
}

// SetAttribute 设置连接属性
func (c *Connection) SetAttribute(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Attributes[key] = value
}

// GetAttribute 获取连接属性
func (c *Connection) GetAttribute(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, exists := c.Attributes[key]
	return value, exists
}

// UpdateHeartbeat 更新心跳时间
func (c *Connection) UpdateHeartbeat() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.LastHeartbeat = time.Now()
	c.LastActivity = time.Now()
}

// UpdateActivity 更新活动时间
func (c *Connection) UpdateActivity() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.LastActivity = time.Now()
}

// Close 关闭连接
func (c *Connection) Close() {
	select {
	case <-c.CloseChan:
		// 已经关闭
		return
	default:
		close(c.CloseChan)
	}

	if c.Conn != nil {
		c.Conn.Close()
	}

	c.logger.Debug("Connection closed", "conn_id", c.ID, "player_id", c.PlayerID)
}

// IsClosed 检查连接是否已关闭
func (c *Connection) IsClosed() bool {
	select {
	case <-c.CloseChan:
		return true
	default:
		return false
	}
}

// handleConnection 处理连接
func (c *Connection) handleConnection(cm *ConnectionManager) {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error("Connection handler panic", "conn_id", c.ID, "error", r)
		}
		cm.RemoveConnection(c.ID)
	}()

	// 启动发送协程
	go c.sendLoop(cm)

	// 启动接收协程
	go c.receiveLoop(cm)

	// 等待连接关闭
	<-c.CloseChan
}

// sendLoop 发送循环
func (c *Connection) sendLoop(cm *ConnectionManager) {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error("Send loop panic", "conn_id", c.ID, "error", r)
		}
	}()

	for {
		select {
		case <-c.CloseChan:
			return
		case msg := <-c.SendChan:
			if err := c.writeMessage(msg); err != nil {
				c.logger.Error("Failed to write message", "conn_id", c.ID, "error", err)
				c.Close()
				return
			}
			c.UpdateActivity()
		}
	}
}

// receiveLoop 接收循环
func (c *Connection) receiveLoop(cm *ConnectionManager) {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error("Receive loop panic", "conn_id", c.ID, "error", r)
		}
	}()

	for {
		select {
		case <-c.CloseChan:
			return
		default:
			msg, err := c.readMessage()
			if err != nil {
				c.logger.Error("Failed to read message", "conn_id", c.ID, "error", err)
				c.Close()
				return
			}

			c.UpdateActivity()

			// 处理心跳消息
			if msg.Header.MessageType == protocol.MsgHeartbeat {
				c.UpdateHeartbeat()
				// 发送心跳响应
				response := &protocol.Message{
					Header: protocol.MessageHeader{
						Magic:       protocol.MessageMagic,
						MessageID:   msg.Header.MessageID,
						MessageType: protocol.MsgPong,
						PlayerID:    msg.Header.PlayerID,
						Timestamp:   time.Now().Unix(),
						Sequence:    msg.Header.Sequence,
					},
					Payload: &protocol.HeartbeatResponse{
						BaseResponse: protocol.NewBaseResponse(true, "pong"),
						ServerTime:   time.Now().Unix(),
					},
				}
				c.SendMessage(response)
				continue
			}

			// 将消息发送到接收通道
			select {
			case c.ReceiveChan <- msg:
				// 消息发送成功
			default:
				// 接收缓冲区满
				c.logger.Warn("Receive buffer full, dropping message", "conn_id", c.ID, "message_type", msg.Header.MessageType)
			}
		}
	}
}

// writeMessage 写入消息
func (c *Connection) writeMessage(msg *protocol.Message) error {
	// TODO: 实现消息序列化和写入
	// 这里应该将消息序列化为字节流并写入连接
	return nil
}

// readMessage 读取消息
func (c *Connection) readMessage() (*protocol.Message, error) {
	// TODO: 实现消息读取和反序列化
	// 这里应该从连接读取字节流并反序列化为消息
	return nil, fmt.Errorf("not implemented")
}