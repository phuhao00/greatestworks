package connection

import (
	"context"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logger"
	"greatestworks/internal/interfaces/tcp/protocol"
)

// HeartbeatConfig 心跳配置
type HeartbeatConfig struct {
	Interval  time.Duration // 心跳间隔
	Timeout   time.Duration // 心跳超时时间
	MaxMissed int           // 最大丢失心跳次数
	Enabled   bool          // 是否启用心跳
}

// DefaultHeartbeatConfig 默认心跳配置
func DefaultHeartbeatConfig() *HeartbeatConfig {
	return &HeartbeatConfig{
		Interval:  30 * time.Second,
		Timeout:   10 * time.Second,
		MaxMissed: 3,
		Enabled:   true,
	}
}

// HeartbeatStatus 心跳状态
type HeartbeatStatus struct {
	LastSent     time.Time
	LastReceived time.Time
	MissedCount  int
	RTT          time.Duration // 往返时间
	IsAlive      bool
}

// HeartbeatManager 心跳管理器
type HeartbeatManager struct {
	config       *HeartbeatConfig
	connections  map[string]*Connection
	heartbeats   map[string]*HeartbeatStatus
	mutex        sync.RWMutex
	logger       logger.Logger
	ticker       *time.Ticker
	ctx          context.Context
	cancel       context.CancelFunc
	onDisconnect func(connID string) // 断线回调
}

// NewHeartbeatManager 创建心跳管理器
func NewHeartbeatManager(config *HeartbeatConfig, logger logger.Logger) *HeartbeatManager {
	if config == nil {
		config = DefaultHeartbeatConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	hm := &HeartbeatManager{
		config:      config,
		connections: make(map[string]*Connection),
		heartbeats:  make(map[string]*HeartbeatStatus),
		logger:      logger,
		ctx:         ctx,
		cancel:      cancel,
	}

	if config.Enabled {
		hm.ticker = time.NewTicker(config.Interval)
		go hm.heartbeatRoutine()
		hm.logger.Info("Heartbeat manager started",
			"interval", config.Interval,
			"timeout", config.Timeout,
			"max_missed", config.MaxMissed)
	}

	return hm
}

// AddConnection 添加连接到心跳管理
func (hm *HeartbeatManager) AddConnection(conn *Connection) {
	if !hm.config.Enabled {
		return
	}

	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	now := time.Now()
	hm.connections[conn.ID] = conn
	hm.heartbeats[conn.ID] = &HeartbeatStatus{
		LastSent:     now,
		LastReceived: now,
		MissedCount:  0,
		RTT:          0,
		IsAlive:      true,
	}

	hm.logger.Debug("Connection added to heartbeat manager", "conn_id", conn.ID)
}

// RemoveConnection 从心跳管理中移除连接
func (hm *HeartbeatManager) RemoveConnection(connID string) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	delete(hm.connections, connID)
	delete(hm.heartbeats, connID)

	hm.logger.Debug("Connection removed from heartbeat manager", "conn_id", connID)
}

// HandleHeartbeatRequest 处理心跳请求
func (hm *HeartbeatManager) HandleHeartbeatRequest(conn *Connection, msg *protocol.Message) error {
	if !hm.config.Enabled {
		return nil
	}

	hm.mutex.Lock()
	status, exists := hm.heartbeats[conn.ID]
	if exists {
		status.LastReceived = time.Now()
		status.MissedCount = 0
		status.IsAlive = true

		// 计算RTT
		if !status.LastSent.IsZero() {
			status.RTT = time.Since(status.LastSent)
		}
	}
	hm.mutex.Unlock()

	// 发送心跳响应
	return hm.sendHeartbeatResponse(conn, msg)
}

// HandleHeartbeatResponse 处理心跳响应
func (hm *HeartbeatManager) HandleHeartbeatResponse(conn *Connection, msg *protocol.Message) error {
	if !hm.config.Enabled {
		return nil
	}

	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	status, exists := hm.heartbeats[conn.ID]
	if !exists {
		return nil
	}

	status.LastReceived = time.Now()
	status.MissedCount = 0
	status.IsAlive = true

	// 计算RTT
	if !status.LastSent.IsZero() {
		status.RTT = time.Since(status.LastSent)
	}

	hm.logger.Debug("Heartbeat response received",
		"conn_id", conn.ID,
		"rtt", status.RTT)

	return nil
}

// sendHeartbeatResponse 发送心跳响应
func (hm *HeartbeatManager) sendHeartbeatResponse(conn *Connection, requestMsg *protocol.Message) error {
	heartbeatResp := &protocol.HeartbeatResponse{
		BaseResponse: protocol.NewBaseResponse(true, "pong"),
		ServerTime:   time.Now().Unix(),
	}

	response := &protocol.Message{
		Header: protocol.MessageHeader{
			Magic:       protocol.MessageMagic,
			MessageID:   requestMsg.Header.MessageID,
			MessageType: protocol.MsgHeartbeat,
			Flags:       protocol.FlagResponse,
			PlayerID:    requestMsg.Header.PlayerID,
			Timestamp:   time.Now().Unix(),
			Sequence:    0,
		},
		Payload: heartbeatResp,
	}

	return conn.SendMessage(response)
}

// sendHeartbeatRequest 发送心跳请求
func (hm *HeartbeatManager) sendHeartbeatRequest(conn *Connection) error {
	heartbeatReq := &protocol.HeartbeatRequest{
		ClientTime: time.Now().Unix(),
	}
	heartbeatReq.Timestamp = time.Now().Unix()

	request := &protocol.Message{
		Header: protocol.MessageHeader{
			Magic:       protocol.MessageMagic,
			MessageID:   uint32(time.Now().UnixNano() & 0xFFFFFFFF),
			MessageType: protocol.MsgHeartbeat,
			Flags:       protocol.FlagRequest,
			PlayerID:    0, // 心跳消息不需要玩家ID
			Timestamp:   time.Now().Unix(),
			Sequence:    0,
		},
		Payload: heartbeatReq,
	}

	// 更新发送时间
	hm.mutex.Lock()
	if status, exists := hm.heartbeats[conn.ID]; exists {
		status.LastSent = time.Now()
	}
	hm.mutex.Unlock()

	return conn.SendMessage(request)
}

// heartbeatRoutine 心跳检测协程
func (hm *HeartbeatManager) heartbeatRoutine() {
	hm.logger.Info("Heartbeat routine started")

	for {
		select {
		case <-hm.ctx.Done():
			hm.logger.Info("Heartbeat routine stopped")
			return
		case <-hm.ticker.C:
			hm.checkHeartbeats()
		}
	}
}

// checkHeartbeats 检查所有连接的心跳状态
func (hm *HeartbeatManager) checkHeartbeats() {
	hm.mutex.Lock()
	connections := make([]*Connection, 0, len(hm.connections))
	heartbeats := make(map[string]*HeartbeatStatus)

	for connID, conn := range hm.connections {
		connections = append(connections, conn)
		if status, exists := hm.heartbeats[connID]; exists {
			heartbeats[connID] = status
		}
	}
	hm.mutex.Unlock()

	now := time.Now()
	var deadConnections []string

	for _, conn := range connections {
		status, exists := heartbeats[conn.ID]
		if !exists {
			continue
		}

		// 检查是否超时
		timeSinceLastReceived := now.Sub(status.LastReceived)
		if timeSinceLastReceived > hm.config.Timeout {
			status.MissedCount++
			status.IsAlive = false

			hm.logger.Warn("Heartbeat timeout",
				"conn_id", conn.ID,
				"player_id", conn.PlayerID,
				"missed_count", status.MissedCount,
				"timeout_duration", timeSinceLastReceived)

			// 检查是否达到最大丢失次数
			if status.MissedCount >= hm.config.MaxMissed {
				deadConnections = append(deadConnections, conn.ID)
				continue
			}
		}

		// 发送心跳请求
		if err := hm.sendHeartbeatRequest(conn); err != nil {
			hm.logger.Error("Failed to send heartbeat request",
				"error", err,
				"conn_id", conn.ID,
				"player_id", conn.PlayerID)
			status.MissedCount++
		}
	}

	// 处理死连接
	for _, connID := range deadConnections {
		hm.handleDeadConnection(connID)
	}

	if len(deadConnections) > 0 {
		hm.logger.Info("Dead connections detected", "count", len(deadConnections))
	}
}

// handleDeadConnection 处理死连接
func (hm *HeartbeatManager) handleDeadConnection(connID string) {
	hm.mutex.Lock()
	conn, exists := hm.connections[connID]
	if exists {
		delete(hm.connections, connID)
		delete(hm.heartbeats, connID)
	}
	hm.mutex.Unlock()

	if exists {
		hm.logger.Warn("Connection marked as dead due to heartbeat failure",
			"conn_id", connID,
			"player_id", conn.PlayerID)

		// 关闭连接
		conn.Close()
		// 调用断线回调
		if hm.onDisconnect != nil {
			hm.onDisconnect(connID)
		}
	}
}

// SetDisconnectCallback 设置断线回调
func (hm *HeartbeatManager) SetDisconnectCallback(callback func(connID string)) {
	hm.onDisconnect = callback
}

// GetHeartbeatStatus 获取连接的心跳状态
func (hm *HeartbeatManager) GetHeartbeatStatus(connID string) (*HeartbeatStatus, bool) {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	status, exists := hm.heartbeats[connID]
	if !exists {
		return nil, false
	}

	// 返回状态副本
	return &HeartbeatStatus{
		LastSent:     status.LastSent,
		LastReceived: status.LastReceived,
		MissedCount:  status.MissedCount,
		RTT:          status.RTT,
		IsAlive:      status.IsAlive,
	}, true
}

// GetAllHeartbeatStatus 获取所有连接的心跳状态
func (hm *HeartbeatManager) GetAllHeartbeatStatus() map[string]*HeartbeatStatus {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	result := make(map[string]*HeartbeatStatus)
	for connID, status := range hm.heartbeats {
		result[connID] = &HeartbeatStatus{
			LastSent:     status.LastSent,
			LastReceived: status.LastReceived,
			MissedCount:  status.MissedCount,
			RTT:          status.RTT,
			IsAlive:      status.IsAlive,
		}
	}

	return result
}

// GetStats 获取心跳统计信息
func (hm *HeartbeatManager) GetStats() map[string]interface{} {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	totalConnections := len(hm.connections)
	aliveConnections := 0
	deadConnections := 0
	totalRTT := time.Duration(0)
	validRTTCount := 0

	for _, status := range hm.heartbeats {
		if status.IsAlive {
			aliveConnections++
		} else {
			deadConnections++
		}

		if status.RTT > 0 {
			totalRTT += status.RTT
			validRTTCount++
		}
	}

	averageRTT := time.Duration(0)
	if validRTTCount > 0 {
		averageRTT = totalRTT / time.Duration(validRTTCount)
	}

	return map[string]interface{}{
		"enabled":           hm.config.Enabled,
		"interval":          hm.config.Interval.String(),
		"timeout":           hm.config.Timeout.String(),
		"max_missed":        hm.config.MaxMissed,
		"total_connections": totalConnections,
		"alive_connections": aliveConnections,
		"dead_connections":  deadConnections,
		"average_rtt":       averageRTT.String(),
		"valid_rtt_count":   validRTTCount,
	}
}

// IsConnectionAlive 检查连接是否存活
func (hm *HeartbeatManager) IsConnectionAlive(connID string) bool {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	status, exists := hm.heartbeats[connID]
	if !exists {
		return false
	}

	return status.IsAlive
}

// UpdateConfig 更新心跳配置
func (hm *HeartbeatManager) UpdateConfig(config *HeartbeatConfig) {
	if config == nil {
		return
	}

	hm.mutex.Lock()
	oldEnabled := hm.config.Enabled
	hm.config = config
	hm.mutex.Unlock()

	// 如果启用状态发生变化
	if config.Enabled != oldEnabled {
		if config.Enabled {
			// 启用心跳
			if hm.ticker == nil {
				hm.ticker = time.NewTicker(config.Interval)
				go hm.heartbeatRoutine()
			}
		} else {
			// 禁用心跳
			if hm.ticker != nil {
				hm.ticker.Stop()
				hm.ticker = nil
			}
		}
	} else if config.Enabled && hm.ticker != nil {
		// 更新间隔
		hm.ticker.Reset(config.Interval)
	}

	hm.logger.Info("Heartbeat configuration updated",
		"enabled", config.Enabled,
		"interval", config.Interval,
		"timeout", config.Timeout,
		"max_missed", config.MaxMissed)
}

// Stop 停止心跳管理器
func (hm *HeartbeatManager) Stop() {
	hm.cancel()
	if hm.ticker != nil {
		hm.ticker.Stop()
	}
	hm.logger.Info("Heartbeat manager stopped")
}
