// Package network 网络协议定义
// Author: MMO Server Team
// Created: 2024

package network

import (
	"context"
	"encoding/binary"
	"fmt"
	"sync"
	"time"
	// "github.com/phuhao00/netcore-go/pkg/core" // TODO: 暂时注释掉有问题的依赖
)

// MessageType 消息类型定义
type MessageType uint16

const (
	// 系统消息
	MsgTypeHeartbeat MessageType = 1000 + iota
	MsgTypeLogin
	MsgTypeLogout
	MsgTypeAuth
	MsgTypeError

	// 玩家消息
	MsgTypePlayerInfo MessageType = 2000 + iota
	MsgTypePlayerMove
	MsgTypePlayerAction
	MsgTypePlayerChat
	MsgTypePlayerMail

	// 游戏消息
	MsgTypeGameBattle MessageType = 3000 + iota
	MsgTypeGameShop
	MsgTypeGameBag
	MsgTypeGamePet
	MsgTypeGameBuilding

	// RPC消息
	MsgTypeRPCRequest MessageType = 9000 + iota
	MsgTypeRPCResponse
	MsgTypeRPCNotify
)

// MessageHeader 消息头定义
type MessageHeader struct {
	Magic    uint32      // 魔数 0x12345678
	Length   uint32      // 消息总长度（包含头部）
	Type     MessageType // 消息类型
	Sequence uint32      // 序列号
	Flags    uint16      // 标志位
	Checksum uint16      // 校验和
}

// Message 完整消息结构
type Message struct {
	Header MessageHeader
	Body   []byte
}

const (
	MessageMagic      = 0x12345678
	MessageHeaderSize = 20          // 消息头大小
	MaxMessageSize    = 1024 * 1024 // 最大消息大小 1MB
	MinMessageSize    = MessageHeaderSize
)

// MessageFlags 消息标志位
const (
	FlagCompressed = 1 << iota // 压缩标志
	FlagEncrypted              // 加密标志
	FlagFragment               // 分片标志
	FlagAck                    // 需要确认标志
)

// Connection 连接接口
type Connection interface {
	Send(data []byte) error
	Close() error
	RemoteAddr() string
}

// TCPConnection TCP连接封装
type TCPConnection struct {
	conn      Connection // 使用自定义的Connection接口
	readBuf   []byte
	writeBuf  []byte
	mutex     sync.RWMutex
	closed    bool
	lastPing  time.Time
	userID    string
	sessionID string
}

// NewTCPConnection 创建新的TCP连接
func NewTCPConnection(conn Connection) *TCPConnection {
	return &TCPConnection{
		conn:     conn,
		readBuf:  make([]byte, 4096),
		writeBuf: make([]byte, 4096),
		lastPing: time.Now(),
	}
}

// ReadMessage 读取消息 - 注意：netcore-go框架通过回调提供数据，此方法主要用于兼容性
func (c *TCPConnection) ReadMessage() (*Message, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.closed {
		return nil, fmt.Errorf("connection closed")
	}

	// 在netcore-go架构中，消息读取通过onMessage回调处理
	// 这个方法主要用于兼容性，实际的消息处理在parseMessage中进行
	return nil, fmt.Errorf("direct message reading not supported in netcore-go architecture")
}

// WriteMessage 写入消息
func (c *TCPConnection) WriteMessage(msg *Message) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.closed {
		return fmt.Errorf("connection closed")
	}

	// 设置消息长度和校验和
	msg.Header.Length = uint32(MessageHeaderSize + len(msg.Body))
	msg.Header.Magic = MessageMagic
	msg.Header.Checksum = calculateChecksum(&msg.Header, msg.Body)

	// 序列化消息头
	headerBuf := serializeMessageHeader(&msg.Header)

	// 组合完整消息
	fullMsg := make([]byte, len(headerBuf)+len(msg.Body))
	copy(fullMsg, headerBuf)
	copy(fullMsg[len(headerBuf):], msg.Body)

	// 使用netcore-go的Send方法发送消息
	err := c.conn.Send(fullMsg)
	if err != nil {
		return fmt.Errorf("send message failed: %w", err)
	}

	return nil
}

// Close 关闭连接
func (c *TCPConnection) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	return c.conn.Close()
}

// SetUserID 设置用户ID
func (c *TCPConnection) SetUserID(userID string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.userID = userID
}

// GetUserID 获取用户ID
func (c *TCPConnection) GetUserID() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.userID
}

// SetSessionID 设置会话ID
func (c *TCPConnection) SetSessionID(sessionID string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.sessionID = sessionID
}

// GetSessionID 获取会话ID
func (c *TCPConnection) GetSessionID() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.sessionID
}

// UpdatePing 更新心跳时间
func (c *TCPConnection) UpdatePing() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.lastPing = time.Now()
}

// GetLastPing 获取最后心跳时间
func (c *TCPConnection) GetLastPing() time.Time {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.lastPing
}

// parseMessageHeader 解析消息头
func parseMessageHeader(data []byte) (*MessageHeader, error) {
	if len(data) < MessageHeaderSize {
		return nil, fmt.Errorf("header data too short")
	}

	header := &MessageHeader{
		Magic:    binary.BigEndian.Uint32(data[0:4]),
		Length:   binary.BigEndian.Uint32(data[4:8]),
		Type:     MessageType(binary.BigEndian.Uint16(data[8:10])),
		Sequence: binary.BigEndian.Uint32(data[10:14]),
		Flags:    binary.BigEndian.Uint16(data[14:16]),
		Checksum: binary.BigEndian.Uint16(data[16:18]),
	}

	if header.Magic != MessageMagic {
		return nil, fmt.Errorf("invalid magic number: 0x%x", header.Magic)
	}

	return header, nil
}

// serializeMessageHeader 序列化消息头
func serializeMessageHeader(header *MessageHeader) []byte {
	data := make([]byte, MessageHeaderSize)
	binary.BigEndian.PutUint32(data[0:4], header.Magic)
	binary.BigEndian.PutUint32(data[4:8], header.Length)
	binary.BigEndian.PutUint16(data[8:10], uint16(header.Type))
	binary.BigEndian.PutUint32(data[10:14], header.Sequence)
	binary.BigEndian.PutUint16(data[14:16], header.Flags)
	binary.BigEndian.PutUint16(data[16:18], header.Checksum)
	return data
}

// calculateChecksum 计算校验和
func calculateChecksum(header *MessageHeader, body []byte) uint16 {
	sum := uint32(0)

	// 计算头部校验和（除了校验和字段本身）
	sum += uint32(header.Magic)
	sum += uint32(header.Length)
	sum += uint32(header.Type)
	sum += uint32(header.Sequence)
	sum += uint32(header.Flags)

	// 计算消息体校验和
	for i := 0; i < len(body); i += 2 {
		if i+1 < len(body) {
			sum += uint32(binary.BigEndian.Uint16(body[i : i+2]))
		} else {
			sum += uint32(body[i]) << 8
		}
	}

	// 折叠为16位
	for sum>>16 != 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}

	return uint16(^sum)
}

// verifyChecksum 验证校验和
func verifyChecksum(header *MessageHeader, body []byte) bool {
	expected := calculateChecksum(header, body)
	return header.Checksum == expected
}

// NetworkService netcore-go网络服务接口
type NetworkService interface {
	// StartTCPServer 启动TCP服务器
	StartTCPServer(ctx context.Context, addr string) error

	// StopTCPServer 停止TCP服务器
	StopTCPServer(ctx context.Context) error

	// SendMessage 发送消息
	SendMessage(ctx context.Context, userID string, msg *Message) error

	// BroadcastMessage 广播消息
	BroadcastMessage(ctx context.Context, msg *Message) error

	// GetConnectionCount 获取连接数
	GetConnectionCount(ctx context.Context) (int, error)
}

// Logger 日志接口
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// networkServiceImpl 网络服务实现
type networkServiceImpl struct {
	// tcpServer   *tcp.Server  // TODO: 待实现服务器集成
	// connPool    *pool.ConnPool  // TODO: 待实现连接池
	connections map[string]*TCPConnection
	mutex       sync.RWMutex
	running     bool
	logger      Logger
}

// StartTCPServer 启动TCP服务器
func (n *networkServiceImpl) StartTCPServer(ctx context.Context, addr string) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.running {
		return fmt.Errorf("server already running")
	}

	// TODO: 实现TCP服务器创建
	// 创建netcore-go TCP服务器配置
	// config := &tcp.ServerConfig{
	//     Address:      addr,
	//     MaxConn:      10000,
	//     ReadTimeout:  30 * time.Second,
	//     WriteTimeout: 30 * time.Second,
	// }

	// 创建TCP服务器
	// server, err := tcp.NewServer(config)
	// if err != nil {
	//     return fmt.Errorf("create tcp server failed: %w", err)
	// }

	// 设置连接处理器
	// server.SetOnConnect(n.onConnect)
	// server.SetOnMessage(n.onMessage)
	// server.SetOnDisconnect(n.onDisconnect)

	// n.tcpServer = server
	n.connections = make(map[string]*TCPConnection)
	n.running = true

	// 启动服务器
	// go func() {
	//     if err := n.tcpServer.Start(); err != nil {
	//         n.logger.Error("TCP server start failed", "error", err)
	//     }
	// }()

	n.logger.Info(fmt.Sprintf("TCP server started on %s", addr))
	return nil
}

// StopTCPServer 停止TCP服务器
func (n *networkServiceImpl) StopTCPServer(ctx context.Context) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if !n.running {
		return nil
	}

	n.running = false

	// TODO: 停止TCP服务器
	// if n.tcpServer != nil {
	//     n.tcpServer.Stop()
	// }

	// TODO: 关闭连接池
	// if n.connPool != nil {
	//     n.connPool.Close()
	// }

	// 关闭所有连接
	for _, conn := range n.connections {
		conn.Close()
	}

	n.connections = nil
	n.logger.Info("TCP server stopped")
	return nil
}

// SendMessage 发送消息
func (n *networkServiceImpl) SendMessage(ctx context.Context, userID string, msg *Message) error {
	n.mutex.RLock()
	conn, exists := n.connections[userID]
	n.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("user %s not connected", userID)
	}

	return conn.WriteMessage(msg)
}

// BroadcastMessage 广播消息
func (n *networkServiceImpl) BroadcastMessage(ctx context.Context, msg *Message) error {
	n.mutex.RLock()
	connections := make([]*TCPConnection, 0, len(n.connections))
	for _, conn := range n.connections {
		connections = append(connections, conn)
	}
	n.mutex.RUnlock()

	for _, conn := range connections {
		if err := conn.WriteMessage(msg); err != nil {
			n.logger.Error(fmt.Sprintf("broadcast message failed: %v", err))
		}
	}

	return nil
}

// GetConnectionCount 获取连接数
func (n *networkServiceImpl) GetConnectionCount(ctx context.Context) (int, error) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return len(n.connections), nil
}

// onConnect 连接建立回调
func (n *networkServiceImpl) onConnect(conn Connection) {
	n.logger.Info(fmt.Sprintf("new connection established from %s", conn.RemoteAddr()))

	// 创建TCP连接包装器
	tcpConn := &TCPConnection{
		conn:     conn,
		readBuf:  make([]byte, 4096),
		writeBuf: make([]byte, 4096),
		lastPing: time.Now(),
	}

	// 暂时使用连接地址作为临时ID
	tempID := conn.RemoteAddr()
	n.mutex.Lock()
	n.connections[tempID] = tcpConn
	n.mutex.Unlock()
}

// onMessage 消息接收回调
func (n *networkServiceImpl) onMessage(conn Connection, data []byte) {
	// 解析消息
	msg, err := n.parseMessage(data)
	if err != nil {
		n.logger.Error(fmt.Sprintf("parse message failed: %v", err))
		return
	}

	// 处理消息
	// 需要找到对应的TCPConnection和创建context
	n.mutex.RLock()
	var tcpConn *TCPConnection
	for _, c := range n.connections {
		if c.conn == conn {
			tcpConn = c
			break
		}
	}
	n.mutex.RUnlock()

	if tcpConn != nil {
		ctx := context.Background()
		if err := n.handleMessage(ctx, tcpConn, msg); err != nil {
			n.logger.Error(fmt.Sprintf("handle message failed: %v", err))
		}
	}
}

// onDisconnect 连接断开回调
func (n *networkServiceImpl) onDisconnect(conn Connection) {
	remoteAddr := conn.RemoteAddr()
	n.logger.Info(fmt.Sprintf("connection disconnected from %s", remoteAddr))

	n.mutex.Lock()
	defer n.mutex.Unlock()

	// 查找并移除连接
	for userID, tcpConn := range n.connections {
		if tcpConn.conn == conn {
			delete(n.connections, userID)
			break
		}
	}
}

// parseMessage 解析消息
func (n *networkServiceImpl) parseMessage(data []byte) (*Message, error) {
	if len(data) < MessageHeaderSize {
		return nil, fmt.Errorf("message too short")
	}

	// 解析消息头
	header, err := parseMessageHeader(data[:MessageHeaderSize])
	if err != nil {
		return nil, err
	}

	// 验证消息长度
	if int(header.Length) != len(data) {
		return nil, fmt.Errorf("message length mismatch")
	}

	// 提取消息体
	body := data[MessageHeaderSize:]

	// 验证校验和
	if !verifyChecksum(header, body) {
		return nil, fmt.Errorf("checksum verification failed")
	}

	return &Message{
		Header: *header,
		Body:   body,
	}, nil
}

// NewNetworkService 创建网络服务
func NewNetworkService(logger Logger) NetworkService {
	return &networkServiceImpl{
		connections: make(map[string]*TCPConnection),
		logger:      logger,
	}
}

// handleMessage 处理消息
func (n *networkServiceImpl) handleMessage(ctx context.Context, conn *TCPConnection, msg *Message) error {
	switch msg.Header.Type {
	case MsgTypeHeartbeat:
		conn.UpdatePing()
		// 回复心跳
		resp := &Message{
			Header: MessageHeader{
				Type:     MsgTypeHeartbeat,
				Sequence: msg.Header.Sequence,
			},
		}
		return conn.WriteMessage(resp)

	case MsgTypeLogin:
		// 处理登录消息
		return n.handleLogin(ctx, conn, msg)

	case MsgTypeLogout:
		// 处理登出消息
		return n.handleLogout(ctx, conn, msg)

	default:
		// 其他消息类型的处理
		n.logger.Debug(fmt.Sprintf("received message type: %d, size: %d", msg.Header.Type, len(msg.Body)))
	}

	return nil
}

// handleLogin 处理登录
func (n *networkServiceImpl) handleLogin(ctx context.Context, conn *TCPConnection, msg *Message) error {
	// TODO: 实现登录逻辑
	// 这里应该验证用户凭据，设置用户ID等

	// 临时实现：直接设置用户ID
	userID := fmt.Sprintf("user_%d", time.Now().UnixNano())
	conn.SetUserID(userID)

	// 添加到连接映射
	n.mutex.Lock()
	n.connections[userID] = conn
	n.mutex.Unlock()

	// 回复登录成功
	resp := &Message{
		Header: MessageHeader{
			Type:     MsgTypeLogin,
			Sequence: msg.Header.Sequence,
		},
		Body: []byte(fmt.Sprintf(`{"status":"success","userID":"%s"}`, userID)),
	}

	return conn.WriteMessage(resp)
}

// handleLogout 处理登出
func (n *networkServiceImpl) handleLogout(ctx context.Context, conn *TCPConnection, msg *Message) error {
	userID := conn.GetUserID()
	if userID != "" {
		n.mutex.Lock()
		delete(n.connections, userID)
		n.mutex.Unlock()
	}

	// 回复登出成功
	resp := &Message{
		Header: MessageHeader{
			Type:     MsgTypeLogout,
			Sequence: msg.Header.Sequence,
		},
		Body: []byte(`{"status":"success"}`),
	}

	return conn.WriteMessage(resp)
}
