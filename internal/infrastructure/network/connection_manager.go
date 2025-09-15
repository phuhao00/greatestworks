package network

import (
	"context"
	"fmt"
	"sync"
	"time"
	
	"github.com/netcore-go/netcore"
	"greatestworks/aop/logger"
)

// ConnectionManager 连接管理器
type ConnectionManager struct {
	connections map[string]*ManagedConnection
	mu          sync.RWMutex
	logger      logger.Logger
	config      *ConnectionManagerConfig
	stats       *ConnectionManagerStats
	ctx         context.Context
	cancel      context.CancelFunc
	pools       map[string]*ConnectionPool
}

// ConnectionManagerConfig 连接管理器配置
type ConnectionManagerConfig struct {
	MaxConnections      int           `json:"max_connections" yaml:"max_connections"`
	ConnectionTimeout   time.Duration `json:"connection_timeout" yaml:"connection_timeout"`
	IdleTimeout         time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
	CleanupInterval     time.Duration `json:"cleanup_interval" yaml:"cleanup_interval"`
	HeartbeatInterval   time.Duration `json:"heartbeat_interval" yaml:"heartbeat_interval"`
	EnableLoadBalancing bool          `json:"enable_load_balancing" yaml:"enable_load_balancing"`
	EnableMetrics       bool          `json:"enable_metrics" yaml:"enable_metrics"`
	PoolSize            int           `json:"pool_size" yaml:"pool_size"`
	MaxRetries          int           `json:"max_retries" yaml:"max_retries"`
}

// ManagedConnection 托管连接
type ManagedConnection struct {
	conn        *netcore.Connection
	id          string
	groupID     string
	userID      string
	createdAt   time.Time
	lastActive  time.Time
	isActive    bool
	metadata    map[string]interface{}
	stats       *ConnectionStats
	mu          sync.RWMutex
}

// ConnectionPool 连接池
type ConnectionPool struct {
	name        string
	connections chan *netcore.Connection
	factory     ConnectionFactory
	config      *PoolConfig
	logger      logger.Logger
	stats       *PoolStats
	mu          sync.RWMutex
}

// ConnectionFactory 连接工厂接口
type ConnectionFactory interface {
	// CreateConnection 创建连接
	CreateConnection(ctx context.Context) (*netcore.Connection, error)
	
	// ValidateConnection 验证连接
	ValidateConnection(conn *netcore.Connection) bool
	
	// CloseConnection 关闭连接
	CloseConnection(conn *netcore.Connection) error
}

// PoolConfig 连接池配置
type PoolConfig struct {
	MinSize         int           `json:"min_size" yaml:"min_size"`
	MaxSize         int           `json:"max_size" yaml:"max_size"`
	MaxIdleTime     time.Duration `json:"max_idle_time" yaml:"max_idle_time"`
	ValidationQuery string        `json:"validation_query" yaml:"validation_query"`
	TestOnBorrow    bool          `json:"test_on_borrow" yaml:"test_on_borrow"`
	TestOnReturn    bool          `json:"test_on_return" yaml:"test_on_return"`
}

// Manager 连接管理器接口
type Manager interface {
	// AddConnection 添加连接
	AddConnection(conn *netcore.Connection, userID, groupID string) (*ManagedConnection, error)
	
	// RemoveConnection 移除连接
	RemoveConnection(connID string) error
	
	// GetConnection 获取连接
	GetConnection(connID string) (*ManagedConnection, error)
	
	// GetConnectionsByUser 根据用户ID获取连接
	GetConnectionsByUser(userID string) []*ManagedConnection
	
	// GetConnectionsByGroup 根据组ID获取连接
	GetConnectionsByGroup(groupID string) []*ManagedConnection
	
	// BroadcastToGroup 向组广播消息
	BroadcastToGroup(groupID string, packet *netcore.Packet) error
	
	// BroadcastToUser 向用户广播消息
	BroadcastToUser(userID string, packet *netcore.Packet) error
	
	// BroadcastToAll 向所有连接广播消息
	BroadcastToAll(packet *netcore.Packet) error
	
	// GetActiveConnections 获取活跃连接数
	GetActiveConnections() int
	
	// GetStats 获取统计信息
	GetStats() *ConnectionManagerStats
	
	// CreatePool 创建连接池
	CreatePool(name string, factory ConnectionFactory, config *PoolConfig) error
	
	// GetPool 获取连接池
	GetPool(name string) (*ConnectionPool, error)
	
	// Start 启动管理器
	Start(ctx context.Context) error
	
	// Stop 停止管理器
	Stop() error
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager(config *ConnectionManagerConfig, logger logger.Logger) Manager {
	if config == nil {
		config = &ConnectionManagerConfig{
			MaxConnections:      10000,
			ConnectionTimeout:   30 * time.Second,
			IdleTimeout:         300 * time.Second, // 5分钟
			CleanupInterval:     60 * time.Second,  // 1分钟
			HeartbeatInterval:   30 * time.Second,
			EnableLoadBalancing: true,
			EnableMetrics:       true,
			PoolSize:            100,
			MaxRetries:          3,
		}
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	
	m := &ConnectionManager{
		connections: make(map[string]*ManagedConnection),
		logger:      logger,
		config:      config,
		ctx:         ctx,
		cancel:      cancel,
		pools:       make(map[string]*ConnectionPool),
		stats: &ConnectionManagerStats{
			StartTime:   time.Now(),
			ByGroup:     make(map[string]*GroupStats),
			ByUser:      make(map[string]*UserStats),
			ByPool:      make(map[string]*PoolStats),
		},
	}
	
	logger.Info("Connection manager initialized successfully", "max_connections", config.MaxConnections)
	return m
}

// AddConnection 添加连接
func (m *ConnectionManager) AddConnection(conn *netcore.Connection, userID, groupID string) (*ManagedConnection, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// 检查连接数限制
	if len(m.connections) >= m.config.MaxConnections {
		return nil, fmt.Errorf("maximum connections limit reached: %d", m.config.MaxConnections)
	}
	
	connID := conn.GetID()
	
	// 检查连接是否已存在
	if _, exists := m.connections[connID]; exists {
		return nil, fmt.Errorf("connection %s already exists", connID)
	}
	
	// 创建托管连接
	managedConn := &ManagedConnection{
		conn:       conn,
		id:         connID,
		groupID:    groupID,
		userID:     userID,
		createdAt:  time.Now(),
		lastActive: time.Now(),
		isActive:   true,
		metadata:   make(map[string]interface{}),
		stats: &ConnectionStats{
			ConnectTime: time.Now(),
		},
	}
	
	m.connections[connID] = managedConn
	
	// 更新统计信息
	m.stats.TotalConnections++
	m.stats.ActiveConnections++
	
	// 更新组统计
	if groupID != "" {
		groupStats, exists := m.stats.ByGroup[groupID]
		if !exists {
			groupStats = &GroupStats{}
			m.stats.ByGroup[groupID] = groupStats
		}
		groupStats.ConnectionCount++
	}
	
	// 更新用户统计
	if userID != "" {
		userStats, exists := m.stats.ByUser[userID]
		if !exists {
			userStats = &UserStats{}
			m.stats.ByUser[userID] = userStats
		}
		userStats.ConnectionCount++
	}
	
	m.logger.Info("Connection added successfully", "conn_id", connID, "user_id", userID, "group_id", groupID)
	return managedConn, nil
}

// RemoveConnection 移除连接
func (m *ConnectionManager) RemoveConnection(connID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	managedConn, exists := m.connections[connID]
	if !exists {
		return fmt.Errorf("connection %s not found", connID)
	}
	
	// 关闭连接
	if err := managedConn.conn.Close(); err != nil {
		m.logger.Error("Failed to close connection", "error", err, "conn_id", connID)
	}
	
	// 从映射中删除
	delete(m.connections, connID)
	
	// 更新统计信息
	m.stats.ActiveConnections--
	m.stats.TotalDisconnections++
	
	// 更新组统计
	if managedConn.groupID != "" {
		if groupStats, exists := m.stats.ByGroup[managedConn.groupID]; exists {
			groupStats.ConnectionCount--
			if groupStats.ConnectionCount <= 0 {
				delete(m.stats.ByGroup, managedConn.groupID)
			}
		}
	}
	
	// 更新用户统计
	if managedConn.userID != "" {
		if userStats, exists := m.stats.ByUser[managedConn.userID]; exists {
			userStats.ConnectionCount--
			if userStats.ConnectionCount <= 0 {
				delete(m.stats.ByUser, managedConn.userID)
			}
		}
	}
	
	m.logger.Info("Connection removed successfully", "conn_id", connID, "user_id", managedConn.userID, "group_id", managedConn.groupID)
	return nil
}

// GetConnection 获取连接
func (m *ConnectionManager) GetConnection(connID string) (*ManagedConnection, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	managedConn, exists := m.connections[connID]
	if !exists {
		return nil, fmt.Errorf("connection %s not found", connID)
	}
	
	// 更新最后活跃时间
	managedConn.mu.Lock()
	managedConn.lastActive = time.Now()
	managedConn.mu.Unlock()
	
	return managedConn, nil
}

// GetConnectionsByUser 根据用户ID获取连接
func (m *ConnectionManager) GetConnectionsByUser(userID string) []*ManagedConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	var connections []*ManagedConnection
	for _, managedConn := range m.connections {
		if managedConn.userID == userID {
			connections = append(connections, managedConn)
		}
	}
	
	return connections
}

// GetConnectionsByGroup 根据组ID获取连接
func (m *ConnectionManager) GetConnectionsByGroup(groupID string) []*ManagedConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	var connections []*ManagedConnection
	for _, managedConn := range m.connections {
		if managedConn.groupID == groupID {
			connections = append(connections, managedConn)
		}
	}
	
	return connections
}

// BroadcastToGroup 向组广播消息
func (m *ConnectionManager) BroadcastToGroup(groupID string, packet *netcore.Packet) error {
	connections := m.GetConnectionsByGroup(groupID)
	if len(connections) == 0 {
		m.logger.Debug("No connections found for group", "group_id", groupID)
		return nil
	}
	
	var errors []error
	successCount := 0
	
	for _, managedConn := range connections {
		if err := managedConn.conn.Send(packet); err != nil {
			m.logger.Error("Failed to send message to group connection", "error", err, "conn_id", managedConn.id, "group_id", groupID)
			errors = append(errors, err)
		} else {
			successCount++
			// 更新连接统计
			managedConn.mu.Lock()
			managedConn.stats.MessagesSent++
			managedConn.lastActive = time.Now()
			managedConn.mu.Unlock()
		}
	}
	
	m.logger.Debug("Group broadcast completed", "group_id", groupID, "total_connections", len(connections), "success_count", successCount, "error_count", len(errors))
	
	if len(errors) > 0 {
		return fmt.Errorf("group broadcast failed for %d connections: %v", len(errors), errors[0])
	}
	
	return nil
}

// BroadcastToUser 向用户广播消息
func (m *ConnectionManager) BroadcastToUser(userID string, packet *netcore.Packet) error {
	connections := m.GetConnectionsByUser(userID)
	if len(connections) == 0 {
		m.logger.Debug("No connections found for user", "user_id", userID)
		return nil
	}
	
	var errors []error
	successCount := 0
	
	for _, managedConn := range connections {
		if err := managedConn.conn.Send(packet); err != nil {
			m.logger.Error("Failed to send message to user connection", "error", err, "conn_id", managedConn.id, "user_id", userID)
			errors = append(errors, err)
		} else {
			successCount++
			// 更新连接统计
			managedConn.mu.Lock()
			managedConn.stats.MessagesSent++
			managedConn.lastActive = time.Now()
			managedConn.mu.Unlock()
		}
	}
	
	m.logger.Debug("User broadcast completed", "user_id", userID, "total_connections", len(connections), "success_count", successCount, "error_count", len(errors))
	
	if len(errors) > 0 {
		return fmt.Errorf("user broadcast failed for %d connections: %v", len(errors), errors[0])
	}
	
	return nil
}

// BroadcastToAll 向所有连接广播消息
func (m *ConnectionManager) BroadcastToAll(packet *netcore.Packet) error {
	m.mu.RLock()
	connections := make([]*ManagedConnection, 0, len(m.connections))
	for _, managedConn := range m.connections {
		connections = append(connections, managedConn)
	}
	m.mu.RUnlock()
	
	if len(connections) == 0 {
		m.logger.Debug("No connections to broadcast to")
		return nil
	}
	
	var errors []error
	successCount := 0
	
	for _, managedConn := range connections {
		if err := managedConn.conn.Send(packet); err != nil {
			m.logger.Error("Failed to broadcast to connection", "error", err, "conn_id", managedConn.id)
			errors = append(errors, err)
		} else {
			successCount++
			// 更新连接统计
			managedConn.mu.Lock()
			managedConn.stats.MessagesSent++
			managedConn.lastActive = time.Now()
			managedConn.mu.Unlock()
		}
	}
	
	m.logger.Debug("Global broadcast completed", "total_connections", len(connections), "success_count", successCount, "error_count", len(errors))
	
	if len(errors) > 0 {
		return fmt.Errorf("global broadcast failed for %d connections: %v", len(errors), errors[0])
	}
	
	return nil
}

// GetActiveConnections 获取活跃连接数
func (m *ConnectionManager) GetActiveConnections() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.connections)
}

// GetStats 获取统计信息
func (m *ConnectionManager) GetStats() *ConnectionManagerStats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// 创建统计信息副本
	stats := &ConnectionManagerStats{
		ActiveConnections:    int64(len(m.connections)),
		TotalConnections:     m.stats.TotalConnections,
		TotalDisconnections:  m.stats.TotalDisconnections,
		StartTime:            m.stats.StartTime,
		Uptime:               time.Since(m.stats.StartTime),
		ByGroup:              make(map[string]*GroupStats),
		ByUser:               make(map[string]*UserStats),
		ByPool:               make(map[string]*PoolStats),
	}
	
	// 复制组统计
	for groupID, groupStats := range m.stats.ByGroup {
		stats.ByGroup[groupID] = &GroupStats{
			ConnectionCount: groupStats.ConnectionCount,
			MessagesSent:    groupStats.MessagesSent,
			MessagesReceived: groupStats.MessagesReceived,
			LastActivity:    groupStats.LastActivity,
		}
	}
	
	// 复制用户统计
	for userID, userStats := range m.stats.ByUser {
		stats.ByUser[userID] = &UserStats{
			ConnectionCount: userStats.ConnectionCount,
			MessagesSent:    userStats.MessagesSent,
			MessagesReceived: userStats.MessagesReceived,
			LastActivity:    userStats.LastActivity,
		}
	}
	
	// 复制池统计
	for poolName, poolStats := range m.stats.ByPool {
		stats.ByPool[poolName] = &PoolStats{
			ActiveConnections: poolStats.ActiveConnections,
			IdleConnections:   poolStats.IdleConnections,
			TotalCreated:      poolStats.TotalCreated,
			TotalDestroyed:    poolStats.TotalDestroyed,
			BorrowCount:       poolStats.BorrowCount,
			ReturnCount:       poolStats.ReturnCount,
		}
	}
	
	return stats
}

// CreatePool 创建连接池
func (m *ConnectionManager) CreatePool(name string, factory ConnectionFactory, config *PoolConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.pools[name]; exists {
		return fmt.Errorf("pool %s already exists", name)
	}
	
	if config == nil {
		config = &PoolConfig{
			MinSize:         5,
			MaxSize:         50,
			MaxIdleTime:     300 * time.Second,
			TestOnBorrow:    true,
			TestOnReturn:    false,
		}
	}
	
	pool := &ConnectionPool{
		name:        name,
		connections: make(chan *netcore.Connection, config.MaxSize),
		factory:     factory,
		config:      config,
		logger:      m.logger,
		stats:       &PoolStats{},
	}
	
	m.pools[name] = pool
	m.stats.ByPool[name] = pool.stats
	
	// 预创建最小连接数
	go pool.initialize(m.ctx)
	
	m.logger.Info("Connection pool created successfully", "name", name, "min_size", config.MinSize, "max_size", config.MaxSize)
	return nil
}

// GetPool 获取连接池
func (m *ConnectionManager) GetPool(name string) (*ConnectionPool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	pool, exists := m.pools[name]
	if !exists {
		return nil, fmt.Errorf("pool %s not found", name)
	}
	
	return pool, nil
}

// Start 启动管理器
func (m *ConnectionManager) Start(ctx context.Context) error {
	m.logger.Info("Starting connection manager")
	
	// 启动清理任务
	go m.startCleanup()
	
	// 启动心跳检测
	if m.config.HeartbeatInterval > 0 {
		go m.startHeartbeat()
	}
	
	// 启动指标收集
	if m.config.EnableMetrics {
		go m.collectMetrics()
	}
	
	m.logger.Info("Connection manager started successfully")
	
	// 等待上下文取消
	select {
	case <-ctx.Done():
		m.logger.Info("Connection manager context cancelled")
		return ctx.Err()
	case <-m.ctx.Done():
		m.logger.Info("Connection manager stopped")
		return nil
	}
}

// Stop 停止管理器
func (m *ConnectionManager) Stop() error {
	m.logger.Info("Stopping connection manager")
	
	// 取消上下文
	m.cancel()
	
	// 关闭所有连接
	m.mu.Lock()
	for connID := range m.connections {
		m.RemoveConnection(connID)
	}
	
	// 关闭所有连接池
	for _, pool := range m.pools {
		pool.close()
	}
	m.mu.Unlock()
	
	m.logger.Info("Connection manager stopped successfully")
	return nil
}

// 私有方法

// startCleanup 启动清理任务
func (m *ConnectionManager) startCleanup() {
	ticker := time.NewTicker(m.config.CleanupInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			m.cleanup()
		case <-m.ctx.Done():
			return
		}
	}
}

// cleanup 清理空闲连接
func (m *ConnectionManager) cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	now := time.Now()
	var toRemove []string
	
	for connID, managedConn := range m.connections {
		managedConn.mu.RLock()
		idleTime := now.Sub(managedConn.lastActive)
		managedConn.mu.RUnlock()
		
		if idleTime > m.config.IdleTimeout {
			toRemove = append(toRemove, connID)
		}
	}
	
	for _, connID := range toRemove {
		m.logger.Debug("Removing idle connection", "conn_id", connID)
		m.RemoveConnection(connID)
	}
	
	if len(toRemove) > 0 {
		m.logger.Info("Cleanup completed", "removed_connections", len(toRemove))
	}
}

// startHeartbeat 启动心跳检测
func (m *ConnectionManager) startHeartbeat() {
	ticker := time.NewTicker(m.config.HeartbeatInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			m.sendHeartbeat()
		case <-m.ctx.Done():
			return
		}
	}
}

// sendHeartbeat 发送心跳
func (m *ConnectionManager) sendHeartbeat() {
	heartbeatPacket := netcore.NewPacket(0, []byte("heartbeat"))
	
	m.mu.RLock()
	connections := make([]*ManagedConnection, 0, len(m.connections))
	for _, managedConn := range m.connections {
		connections = append(connections, managedConn)
	}
	m.mu.RUnlock()
	
	for _, managedConn := range connections {
		if err := managedConn.conn.Send(heartbeatPacket); err != nil {
			m.logger.Debug("Failed to send heartbeat", "conn_id", managedConn.id, "error", err)
			// 心跳失败，可能需要移除连接
			go m.RemoveConnection(managedConn.id)
		}
	}
	
	m.logger.Debug("Heartbeat sent to all connections", "connection_count", len(connections))
}

// collectMetrics 收集指标
func (m *ConnectionManager) collectMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			stats := m.GetStats()
			m.logger.Debug("Connection manager metrics",
				"active_connections", stats.ActiveConnections,
				"total_connections", stats.TotalConnections,
				"total_disconnections", stats.TotalDisconnections,
				"uptime", stats.Uptime)
		case <-m.ctx.Done():
			return
		}
	}
}

// 连接池方法

// initialize 初始化连接池
func (p *ConnectionPool) initialize(ctx context.Context) {
	for i := 0; i < p.config.MinSize; i++ {
		conn, err := p.factory.CreateConnection(ctx)
		if err != nil {
			p.logger.Error("Failed to create initial connection", "error", err, "pool", p.name)
			continue
		}
		
		select {
		case p.connections <- conn:
			p.mu.Lock()
			p.stats.TotalCreated++
			p.stats.IdleConnections++
			p.mu.Unlock()
		default:
			// 池已满，关闭连接
			p.factory.CloseConnection(conn)
		}
	}
	
	p.logger.Info("Connection pool initialized", "name", p.name, "initial_connections", p.config.MinSize)
}

// close 关闭连接池
func (p *ConnectionPool) close() {
	close(p.connections)
	
	// 关闭所有连接
	for conn := range p.connections {
		p.factory.CloseConnection(conn)
		p.mu.Lock()
		p.stats.TotalDestroyed++
		p.mu.Unlock()
	}
	
	p.logger.Info("Connection pool closed", "name", p.name)
}

// 统计信息结构
type ConnectionManagerStats struct {
	ActiveConnections   int64                  `json:"active_connections"`
	TotalConnections    int64                  `json:"total_connections"`
	TotalDisconnections int64                  `json:"total_disconnections"`
	StartTime           time.Time              `json:"start_time"`
	Uptime              time.Duration          `json:"uptime"`
	ByGroup             map[string]*GroupStats `json:"by_group"`
	ByUser              map[string]*UserStats  `json:"by_user"`
	ByPool              map[string]*PoolStats  `json:"by_pool"`
}

type ConnectionStats struct {
	ConnectTime      time.Time `json:"connect_time"`
	MessagesSent     int64     `json:"messages_sent"`
	MessagesReceived int64     `json:"messages_received"`
	BytesTransferred int64     `json:"bytes_transferred"`
	LastActivity     time.Time `json:"last_activity"`
}

type GroupStats struct {
	ConnectionCount  int64     `json:"connection_count"`
	MessagesSent     int64     `json:"messages_sent"`
	MessagesReceived int64     `json:"messages_received"`
	LastActivity     time.Time `json:"last_activity"`
}

type UserStats struct {
	ConnectionCount  int64     `json:"connection_count"`
	MessagesSent     int64     `json:"messages_sent"`
	MessagesReceived int64     `json:"messages_received"`
	LastActivity     time.Time `json:"last_activity"`
}

type PoolStats struct {
	ActiveConnections int64 `json:"active_connections"`
	IdleConnections   int64 `json:"idle_connections"`
	TotalCreated      int64 `json:"total_created"`
	TotalDestroyed    int64 `json:"total_destroyed"`
	BorrowCount       int64 `json:"borrow_count"`
	ReturnCount       int64 `json:"return_count"`
}