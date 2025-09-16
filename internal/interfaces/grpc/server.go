package grpc

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"greatestworks/application/handlers"
	"greatestworks/internal/infrastructure/logger"
	"greatestworks/internal/interfaces/grpc/interceptors"
	"greatestworks/internal/interfaces/grpc/services"
	pb "greatestworks/internal/interfaces/grpc/proto"
)

// ServerConfig gRPC服务器配置
type ServerConfig struct {
	Addr                string
	MaxConnections      int
	ConnectionTimeout   time.Duration
	KeepaliveTime       time.Duration
	KeepaliveTimeout    time.Duration
	MaxConnectionIdle   time.Duration
	MaxConnectionAge    time.Duration
	MaxConnectionAgeGrace time.Duration
	EnableReflection    bool
	EnableHealthCheck   bool
	EnableMetrics       bool
	EnableAuth          bool
}

// DefaultServerConfig 默认gRPC服务器配置
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Addr:                  ":9091",
		MaxConnections:        1000,
		ConnectionTimeout:     10 * time.Second,
		KeepaliveTime:         30 * time.Second,
		KeepaliveTimeout:      5 * time.Second,
		MaxConnectionIdle:     15 * time.Minute,
		MaxConnectionAge:      30 * time.Minute,
		MaxConnectionAgeGrace: 5 * time.Second,
		EnableReflection:      true,
		EnableHealthCheck:     true,
		EnableMetrics:         true,
		EnableAuth:            true,
	}
}

// GRPCServer gRPC服务器
type GRPCServer struct {
	config         *ServerConfig
	server         *grpc.Server
	listener       net.Listener
	healthServer   *health.Server
	playerService  *services.PlayerServiceImpl
	battleService  *services.BattleServiceImpl
	notifyService  *services.NotificationServiceImpl
	logger         logger.Logger
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	running        bool
	mutex          sync.RWMutex
}

// NewGRPCServer 创建gRPC服务器
func NewGRPCServer(config *ServerConfig, commandBus *handlers.CommandBus, queryBus *handlers.QueryBus, logger logger.Logger) *GRPCServer {
	if config == nil {
		config = DefaultServerConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	// 创建gRPC服务器选项
	opts := []grpc.ServerOption{
		// 连接保活设置
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    config.KeepaliveTime,
			Timeout: config.KeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: true,
		}),
		// 连接限制
		grpc.MaxConcurrentStreams(uint32(config.MaxConnections)),
		grpc.ConnectionTimeout(config.ConnectionTimeout),
	}

	// 添加拦截器
	if config.EnableAuth {
		opts = append(opts, grpc.UnaryInterceptor(interceptors.AuthUnaryInterceptor(logger)))
		opts = append(opts, grpc.StreamInterceptor(interceptors.AuthStreamInterceptor(logger)))
	}

	if config.EnableMetrics {
		opts = append(opts, grpc.UnaryInterceptor(interceptors.MetricsUnaryInterceptor(logger)))
		opts = append(opts, grpc.StreamInterceptor(interceptors.MetricsStreamInterceptor(logger)))
	}

	// 添加日志拦截器
	opts = append(opts, grpc.UnaryInterceptor(interceptors.LoggingUnaryInterceptor(logger)))
	opts = append(opts, grpc.StreamInterceptor(interceptors.LoggingStreamInterceptor(logger)))

	// 创建gRPC服务器
	grpcServer := grpc.NewServer(opts...)

	// 创建健康检查服务器
	var healthServer *health.Server
	if config.EnableHealthCheck {
		healthServer = health.NewServer()
	}

	// 创建业务服务
	playerService := services.NewPlayerService(commandBus, queryBus, logger)
	battleService := services.NewBattleService(commandBus, queryBus, logger)
	notifyService := services.NewNotificationService(commandBus, queryBus, logger)

	server := &GRPCServer{
		config:        config,
		server:        grpcServer,
		healthServer:  healthServer,
		playerService: playerService,
		battleService: battleService,
		notifyService: notifyService,
		logger:        logger,
		ctx:           ctx,
		cancel:        cancel,
		running:       false,
	}

	// 注册服务
	server.registerServices()

	return server
}

// registerServices 注册gRPC服务
func (s *GRPCServer) registerServices() {
	// 注册业务服务
	pb.RegisterPlayerServiceServer(s.server, s.playerService)
	pb.RegisterBattleServiceServer(s.server, s.battleService)
	pb.RegisterNotificationServiceServer(s.server, s.notifyService)

	// 注册健康检查服务
	if s.config.EnableHealthCheck && s.healthServer != nil {
		grpc_health_v1.RegisterHealthServer(s.server, s.healthServer)
		
		// 设置服务状态
		s.healthServer.SetServingStatus("greatestworks.grpc.PlayerService", grpc_health_v1.HealthCheckResponse_SERVING)
		s.healthServer.SetServingStatus("greatestworks.grpc.BattleService", grpc_health_v1.HealthCheckResponse_SERVING)
		s.healthServer.SetServingStatus("greatestworks.grpc.NotificationService", grpc_health_v1.HealthCheckResponse_SERVING)
		s.healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING) // 整体服务状态
	}

	// 启用反射（开发环境）
	if s.config.EnableReflection {
		reflection.Register(s.server)
		s.logger.Info("gRPC reflection enabled")
	}

	s.logger.Info("gRPC services registered successfully")
}

// Start 启动gRPC服务器
func (s *GRPCServer) Start() error {
	s.mutex.Lock()
	if s.running {
		s.mutex.Unlock()
		return fmt.Errorf("gRPC server is already running")
	}
	s.mutex.Unlock()

	s.logger.Info("Starting gRPC server", "address", s.config.Addr)

	// 创建监听器
	listener, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		s.logger.Error("Failed to create gRPC listener", "error", err, "address", s.config.Addr)
		return fmt.Errorf("failed to create gRPC listener: %w", err)
	}

	s.listener = listener
	s.mutex.Lock()
	s.running = true
	s.mutex.Unlock()

	// 启动服务器
	s.wg.Add(1)
	go s.serve()

	s.logger.Info("gRPC server started successfully", "address", s.config.Addr)
	return nil
}

// serve 服务协程
func (s *GRPCServer) serve() {
	defer s.wg.Done()

	s.logger.Info("gRPC server serving", "address", s.listener.Addr())

	if err := s.server.Serve(s.listener); err != nil {
		select {
		case <-s.ctx.Done():
			// 正常关闭
			s.logger.Info("gRPC server stopped")
		default:
			// 异常退出
			s.logger.Error("gRPC server serve error", "error", err)
		}
	}
}

// Stop 停止gRPC服务器
func (s *GRPCServer) Stop() error {
	s.mutex.Lock()
	if !s.running {
		s.mutex.Unlock()
		return fmt.Errorf("gRPC server is not running")
	}
	s.running = false
	s.mutex.Unlock()

	s.logger.Info("Stopping gRPC server")

	// 取消上下文
	s.cancel()

	// 设置健康检查状态为不可用
	if s.healthServer != nil {
		s.healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		s.healthServer.SetServingStatus("greatestworks.grpc.PlayerService", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		s.healthServer.SetServingStatus("greatestworks.grpc.BattleService", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		s.healthServer.SetServingStatus("greatestworks.grpc.NotificationService", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	}

	// 优雅关闭服务器
	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	// 等待优雅关闭或强制关闭
	select {
	case <-done:
		s.logger.Info("gRPC server gracefully stopped")
	case <-time.After(30 * time.Second):
		s.logger.Warn("gRPC server graceful stop timeout, forcing stop")
		s.server.Stop()
	}

	// 关闭监听器
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			s.logger.Error("Failed to close gRPC listener", "error", err)
		}
	}

	// 等待所有协程结束
	s.wg.Wait()

	s.logger.Info("gRPC server stopped successfully")
	return nil
}

// IsRunning 检查服务器是否运行中
func (s *GRPCServer) IsRunning() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.running
}

// GetStats 获取服务器统计信息
func (s *GRPCServer) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"running":             s.IsRunning(),
		"address":             s.config.Addr,
		"max_connections":     s.config.MaxConnections,
		"connection_timeout":  s.config.ConnectionTimeout.String(),
		"keepalive_time":      s.config.KeepaliveTime.String(),
		"keepalive_timeout":   s.config.KeepaliveTimeout.String(),
		"enable_reflection":   s.config.EnableReflection,
		"enable_health_check": s.config.EnableHealthCheck,
		"enable_metrics":      s.config.EnableMetrics,
		"enable_auth":         s.config.EnableAuth,
		"services": []string{
			"PlayerService",
			"BattleService",
			"NotificationService",
		},
	}
}

// GetServer 获取原始gRPC服务器实例
func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}

// GetListener 获取监听器
func (s *GRPCServer) GetListener() net.Listener {
	return s.listener
}

// SetHealthStatus 设置服务健康状态
func (s *GRPCServer) SetHealthStatus(service string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	if s.healthServer != nil {
		s.healthServer.SetServingStatus(service, status)
		s.logger.Info("Health status updated", "service", service, "status", status.String())
	}
}

// GetHealthStatus 获取服务健康状态
func (s *GRPCServer) GetHealthStatus(service string) grpc_health_v1.HealthCheckResponse_ServingStatus {
	if s.healthServer == nil {
		return grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN
	}

	// 这里需要实现获取健康状态的逻辑
	// 由于grpc/health包没有直接提供获取状态的方法，我们返回默认状态
	if s.IsRunning() {
		return grpc_health_v1.HealthCheckResponse_SERVING
	}
	return grpc_health_v1.HealthCheckResponse_NOT_SERVING
}

// UpdateConfig 更新服务器配置
func (s *GRPCServer) UpdateConfig(config *ServerConfig) error {
	if s.IsRunning() {
		return fmt.Errorf("cannot update config while server is running")
	}

	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	s.mutex.Lock()
	s.config = config
	s.mutex.Unlock()

	s.logger.Info("gRPC server configuration updated")
	return nil
}

// GetConfig 获取服务器配置
func (s *GRPCServer) GetConfig() *ServerConfig {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// 返回配置副本
	configCopy := *s.config
	return &configCopy
}

// RegisterService 动态注册服务
func (s *GRPCServer) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.server.RegisterService(desc, impl)
	s.logger.Info("Service registered dynamically", "service", desc.ServiceName)
}

// GetServiceInfo 获取已注册的服务信息
func (s *GRPCServer) GetServiceInfo() map[string]grpc.ServiceInfo {
	return s.server.GetServiceInfo()
}

// Shutdown 关闭服务器（Stop的别名）
func (s *GRPCServer) Shutdown() error {
	return s.Stop()
}

// Wait 等待服务器停止
func (s *GRPCServer) Wait() {
	s.wg.Wait()
}