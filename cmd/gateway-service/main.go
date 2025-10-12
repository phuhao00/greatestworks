// Package main 网关服务主程序
// 基于DDD架构的分布式网关服务
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"greatestworks/internal/config"
	"greatestworks/internal/infrastructure/logging"
	"greatestworks/internal/infrastructure/monitoring"
	"greatestworks/internal/interfaces/tcp"
)

// GatewayServiceConfig aliases the shared configuration schema for readability.
type GatewayServiceConfig = config.Config

// GatewayService 网关服务
type GatewayService struct {
	config   atomic.Pointer[GatewayServiceConfig]
	logger   logging.Logger
	server   *tcp.TCPServer
	profiler *monitoring.Profiler
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewGatewayService 创建网关服务
func NewGatewayService(cfg *GatewayServiceConfig, logger logging.Logger) *GatewayService {
	ctx, cancel := context.WithCancel(context.Background())

	service := &GatewayService{
		logger:   logger,
		ctx:      ctx,
		cancel:   cancel,
		profiler: monitoring.NewProfiler(logger),
	}

	if cfg != nil {
		service.config.Store(cfg)
	}

	return service
}

// UpdateConfig replaces the in-memory configuration snapshot.
func (s *GatewayService) UpdateConfig(cfg *GatewayServiceConfig) {
	if cfg == nil {
		return
	}
	s.config.Store(cfg)
}

// Start 启动网关服务
func (s *GatewayService) Start() error {
	cfg := s.config.Load()
	if cfg == nil {
		return fmt.Errorf("gateway service configuration not loaded")
	}

	s.logger.Info("Starting gateway service", logging.Fields{
		"service": cfg.Service.Name,
		"version": cfg.Service.Version,
		"node_id": cfg.Service.NodeID,
	})

	if err := s.initializeDatabase(cfg); err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}

	if err := s.initializeGameServices(cfg); err != nil {
		return fmt.Errorf("初始化游戏服务连接失败: %w", err)
	}

	if err := s.initializeTCPServer(cfg); err != nil {
		return fmt.Errorf("初始化TCP服务器失败: %w", err)
	}

	go func() {
		if err := s.server.Start(); err != nil {
			s.logger.Error("TCP server start failed", err)
		}
	}()

	if cfg.Monitoring.Profiling.Enabled {
		host := cfg.Monitoring.Profiling.Host
		if host == "" {
			host = cfg.Server.TCP.Host
		}

		if cfg.Monitoring.Profiling.Port == 0 {
			s.logger.Warn("pprof未启动: 未配置端口")
		} else if err := s.profiler.Start(host, cfg.Monitoring.Profiling.Port); err != nil {
			s.logger.Error("Failed to start pprof server", err, logging.Fields{
				"host": host,
				"port": cfg.Monitoring.Profiling.Port,
			})
		}
	}

	s.logger.Info("Gateway service started successfully", logging.Fields{
		"tcp_addr": fmt.Sprintf("%s:%d", cfg.Server.TCP.Host, cfg.Server.TCP.Port),
	})

	return nil
}

// Stop 停止网关服务
func (s *GatewayService) Stop() error {
	s.logger.Info("停止网关服务")

	s.cancel()

	if s.server != nil {
		if err := s.server.Stop(); err != nil {
			s.logger.Error("Failed to stop TCP server", err)
			return err
		}
	}

	if s.profiler != nil {
		if err := s.profiler.Stop(context.Background()); err != nil {
			s.logger.Error("Failed to stop pprof server", err)
			return err
		}
	}

	s.logger.Info("网关服务已停止")
	return nil
}

// initializeDatabase 初始化数据库连接
func (s *GatewayService) initializeDatabase(cfg *GatewayServiceConfig) error {
	s.logger.Info("初始化数据库连接")

	// TODO: 实现Redis连接，使用 cfg.Database.Redis

	s.logger.Info("数据库连接初始化完成")
	return nil
}

// initializeGameServices 初始化游戏服务连接
func (s *GatewayService) initializeGameServices(cfg *GatewayServiceConfig) error {
	s.logger.Info("初始化游戏服务连接")

	// TODO: 实现服务发现（cfg.Gateway.GameServices.Discovery）
	// TODO: 实现负载均衡（cfg.Gateway.GameServices.LoadBalancer）
	// TODO: 实现健康检查

	s.logger.Info("游戏服务连接初始化完成")
	return nil
}

// initializeTCPServer 初始化TCP服务器
func (s *GatewayService) initializeTCPServer(cfg *GatewayServiceConfig) error {
	s.logger.Info("初始化TCP服务器")

	// TODO: 使用 cfg.Server.TCP 初始化真正的 TCP 服务

	s.logger.Info("TCP服务器初始化完成")
	return nil
}

// loadInitialConfig 加载配置并返回配置与文件来源
func loadInitialConfig() (*GatewayServiceConfig, []string, *config.Loader, error) {
	loader := config.NewLoader(
		config.WithService("gateway-service"),
	)

	cfg, sources, err := loader.Load()
	if err != nil {
		return nil, nil, nil, err
	}

	return cfg, sources, loader, nil
}

// main 主函数
func main() {
	logger := logging.NewBaseLogger(logging.InfoLevel)
	logger.Info("启动网关服务")

	cfg, sources, loader, err := loadInitialConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	logger.Info("配置加载成功", logging.Fields{
		"environment": cfg.App.Environment,
		"sources":     sources,
	})

	manager, err := config.NewManager(loader)
	if err != nil {
		log.Fatalf("创建配置管理器失败: %v", err)
	}
	defer func() {
		_ = manager.Close()
	}()

	runtimeCfg := manager.Config()
	service := NewGatewayService(runtimeCfg, logger)

	manager.OnChange(func(next *config.Config) {
		if next == nil {
			return
		}
		service.UpdateConfig(next)
		logger.Info("网关服务配置已刷新", logging.Fields{
			"service_version": next.Service.Version,
		})
	})

	watchCtx, watchCancel := context.WithCancel(context.Background())
	defer watchCancel()

	if runtimeCfg != nil && runtimeCfg.Environment.HotReload {
		if err := manager.StartWatching(watchCtx); err != nil {
			logger.Error("启动配置热更新监听失败", err, logging.Fields{})
		} else {
			logger.Info("已启用配置热更新", logging.Fields{})
		}
	}

	if err := service.Start(); err != nil {
		log.Fatalf("启动网关服务失败: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	select {
	case sig := <-sigChan:
		logger.Info("收到关闭信号", logging.Fields{
			"signal": sig.String(),
		})
	case <-service.ctx.Done():
		logger.Info("上下文已取消")
	}

	logger.Info("正在关闭网关服务...")
	watchCancel()
	if err := service.Stop(); err != nil {
		logger.Error("关闭网关服务失败", err, logging.Fields{})
		os.Exit(1)
	}

	logger.Info("网关服务已关闭")
}
