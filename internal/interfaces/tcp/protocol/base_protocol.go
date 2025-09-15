package protocol

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

// Message 基础消息结构
type Message struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data,omitempty"`
	Timestamp int64           `json:"timestamp"`
	PlayerID  uint64          `json:"player_id,omitempty"`
}

// Response 响应消息结构
type Response struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// ErrorInfo 错误信息
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Notification 通知消息结构
type Notification struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
	PlayerID  uint64      `json:"player_id,omitempty"`
}

// 通用消息类型
const (
	// 连接相关
	MsgTypeConnect    = "connect"
	MsgTypeDisconnect = "disconnect"
	MsgTypePing       = "ping"
	MsgTypePong       = "pong"
	
	// 认证相关
	MsgTypeAuth        = "auth"
	MsgTypeAuthSuccess = "auth_success"
	MsgTypeAuthFailed  = "auth_failed"
	
	// 错误相关
	MsgTypeError = "error"
	
	// 通知相关
	MsgTypeNotification = "notification"
)

// 错误代码
const (
	ErrorCodeInvalidMessage   = "INVALID_MESSAGE"
	ErrorCodeUnauthorized     = "UNAUTHORIZED"
	ErrorCodeNotFound         = "NOT_FOUND"
	ErrorCodeInternalError    = "INTERNAL_ERROR"
	ErrorCodeInvalidParameter = "INVALID_PARAMETER"
	ErrorCodePermissionDenied = "PERMISSION_DENIED"
	ErrorCodeRateLimit        = "RATE_LIMIT"
	ErrorCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
)

// ConnectRequest 连接请求
type ConnectRequest struct {
	ClientID    string `json:"client_id"`
	Version     string `json:"version"`
	Platform    string `json:"platform,omitempty"`
	DeviceID    string `json:"device_id,omitempty"`
	Compression bool   `json:"compression,omitempty"`
}

// ConnectResponse 连接响应
type ConnectResponse struct {
	SessionID   string `json:"session_id"`
	ServerTime  int64  `json:"server_time"`
	Heartbeat   int32  `json:"heartbeat"`
	Compression bool   `json:"compression"`
}

// AuthRequest 认证请求
type AuthRequest struct {
	Token    string `json:"token"`
	PlayerID uint64 `json:"player_id,omitempty"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	PlayerID    uint64 `json:"player_id"`
	DisplayName string `json:"display_name"`
	Permissions []string `json:"permissions,omitempty"`
	AuthTime    int64  `json:"auth_time"`
}

// PingRequest Ping请求
type PingRequest struct {
	Timestamp int64 `json:"timestamp"`
}

// PongResponse Pong响应
type PongResponse struct {
	Timestamp   int64 `json:"timestamp"`
	ServerTime  int64 `json:"server_time"`
	RoundTrip   int64 `json:"round_trip,omitempty"`
}

// TCPConnection TCP连接封装
type TCPConnection struct {
	conn        net.Conn
	sessionID   string
	playerID    uint64
	isAuth      bool
	lastPing    time.Time
	mutex       sync.RWMutex
	closed      bool
	ctx         context.Context
	cancel      context.CancelFunc
	messageChan chan []byte
}

// NewTCPConnection 创建TCP连接
func NewTCPConnection(conn net.Conn, sessionID string) *TCPConnection {
	ctx, cancel := context.WithCancel(context.Background())
	return &TCPConnection{
		conn:        conn,
		sessionID:   sessionID,
		lastPing:    time.Now(),
		ctx:         ctx,
		cancel:      cancel,
		messageChan: make(chan []byte, 100),
	}
}

// GetSessionID 获取会话ID
func (c *TCPConnection) GetSessionID() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.sessionID
}

// GetPlayerID 获取玩家ID
func (c *TCPConnection) GetPlayerID() uint64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.playerID
}

// SetPlayerID 设置玩家ID
func (c *TCPConnection) SetPlayerID(playerID uint64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.playerID = playerID
}

// IsAuthenticated 检查是否已认证
func (c *TCPConnection) IsAuthenticated() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.isAuth
}

// SetAuthenticated 设置认证状态
func (c *TCPConnection) SetAuthenticated(auth bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.isAuth = auth
}

// IsClosed 检查连接是否已关闭
func (c *TCPConnection) IsClosed() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.closed
}

// UpdatePing 更新Ping时间
func (c *TCPConnection) UpdatePing() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.lastPing = time.Now()
}

// GetLastPing 获取最后Ping时间
func (c *TCPConnection) GetLastPing() time.Time {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.lastPing
}

// SendMessage 发送消息
func (c *TCPConnection) SendMessage(msg *Message) error {
	c.mutex.RLock()
	if c.closed {
		c.mutex.RUnlock()
		return fmt.Errorf("connection is closed")
	}
	c.mutex.RUnlock()
	
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	
	select {
	case c.messageChan <- data:
		return nil
	case <-c.ctx.Done():
		return fmt.Errorf("connection context cancelled")
	default:
		return fmt.Errorf("message channel is full")
	}
}

// SendResponse 发送响应
func (c *TCPConnection) SendResponse(msgID, msgType string, data interface{}) error {
	resp := &Response{
		ID:        msgID,
		Type:      msgType,
		Success:   true,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
	
	return c.SendMessage(&Message{
		ID:        msgID,
		Type:      msgType,
		Data:      mustMarshal(resp),
		Timestamp: time.Now().Unix(),
		PlayerID:  c.GetPlayerID(),
	})
}

// SendError 发送错误响应
func (c *TCPConnection) SendError(msgID, msgType, errorCode, errorMessage string) error {
	resp := &Response{
		ID:      msgID,
		Type:    msgType,
		Success: false,
		Error: &ErrorInfo{
			Code:    errorCode,
			Message: errorMessage,
		},
		Timestamp: time.Now().Unix(),
	}
	
	return c.SendMessage(&Message{
		ID:        msgID,
		Type:      MsgTypeError,
		Data:      mustMarshal(resp),
		Timestamp: time.Now().Unix(),
		PlayerID:  c.GetPlayerID(),
	})
}

// SendNotification 发送通知
func (c *TCPConnection) SendNotification(notificationType string, data interface{}) error {
	notification := &Notification{
		Type:      notificationType,
		Data:      data,
		Timestamp: time.Now().Unix(),
		PlayerID:  c.GetPlayerID(),
	}
	
	return c.SendMessage(&Message{
		Type:      MsgTypeNotification,
		Data:      mustMarshal(notification),
		Timestamp: time.Now().Unix(),
		PlayerID:  c.GetPlayerID(),
	})
}

// ReadMessage 读取消息
func (c *TCPConnection) ReadMessage() (*Message, error) {
	c.mutex.RLock()
	if c.closed {
		c.mutex.RUnlock()
		return nil, fmt.Errorf("connection is closed")
	}
	c.mutex.RUnlock()
	
	// 设置读取超时
	c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	
	// 读取消息长度（4字节）
	lengthBytes := make([]byte, 4)
	if _, err := c.conn.Read(lengthBytes); err != nil {
		return nil, fmt.Errorf("failed to read message length: %w", err)
	}
	
	// 解析消息长度
	length := int(lengthBytes[0])<<24 | int(lengthBytes[1])<<16 | int(lengthBytes[2])<<8 | int(lengthBytes[3])
	if length <= 0 || length > 1024*1024 { // 最大1MB
		return nil, fmt.Errorf("invalid message length: %d", length)
	}
	
	// 读取消息内容
	msgBytes := make([]byte, length)
	if _, err := c.conn.Read(msgBytes); err != nil {
		return nil, fmt.Errorf("failed to read message data: %w", err)
	}
	
	// 解析消息
	var msg Message
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}
	
	return &msg, nil
}

// WriteMessage 写入消息
func (c *TCPConnection) WriteMessage(data []byte) error {
	c.mutex.RLock()
	if c.closed {
		c.mutex.RUnlock()
		return fmt.Errorf("connection is closed")
	}
	c.mutex.RUnlock()
	
	// 设置写入超时
	c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	
	// 写入消息长度（4字节）
	length := len(data)
	lengthBytes := []byte{
		byte(length >> 24),
		byte(length >> 16),
		byte(length >> 8),
		byte(length),
	}
	
	if _, err := c.conn.Write(lengthBytes); err != nil {
		return fmt.Errorf("failed to write message length: %w", err)
	}
	
	// 写入消息内容
	if _, err := c.conn.Write(data); err != nil {
		return fmt.Errorf("failed to write message data: %w", err)
	}
	
	return nil
}

// StartMessageWriter 启动消息写入协程
func (c *TCPConnection) StartMessageWriter() {
	go func() {
		for {
			select {
			case data := <-c.messageChan:
				if err := c.WriteMessage(data); err != nil {
					// 写入失败，关闭连接
					c.Close()
					return
				}
			case <-c.ctx.Done():
				return
			}
		}
	}()
}

// Close 关闭连接
func (c *TCPConnection) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	if c.closed {
		return nil
	}
	
	c.closed = true
	c.cancel()
	close(c.messageChan)
	
	return c.conn.Close()
}

// GetRemoteAddr 获取远程地址
func (c *TCPConnection) GetRemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// 辅助函数

func mustMarshal(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal: %v", err))
	}
	return data
}

// MessageHandler 消息处理器接口
type MessageHandler func(ctx context.Context, conn *TCPConnection, msg *Message) error

// TCPRouter TCP路由器
type TCPRouter struct {
	handlers map[string]MessageHandler
	mutex    sync.RWMutex
}

// NewTCPRouter 创建TCP路由器
func NewTCPRouter() *TCPRouter {
	return &TCPRouter{
		handlers: make(map[string]MessageHandler),
	}
}

// RegisterHandler 注册消息处理器
func (r *TCPRouter) RegisterHandler(msgType string, handler MessageHandler) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.handlers[msgType] = handler
}

// HandleMessage 处理消息
func (r *TCPRouter) HandleMessage(ctx context.Context, conn *TCPConnection, msg *Message) error {
	r.mutex.RLock()
	handler, exists := r.handlers[msg.Type]
	r.mutex.RUnlock()
	
	if !exists {
		return conn.SendError(msg.ID, msg.Type, ErrorCodeNotFound, fmt.Sprintf("unknown message type: %s", msg.Type))
	}
	
	return handler(ctx, conn, msg)
}

// GetHandlers 获取所有处理器
func (r *TCPRouter) GetHandlers() map[string]MessageHandler {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	handlers := make(map[string]MessageHandler)
	for k, v := range r.handlers {
		handlers[k] = v
	}
	return handlers
}