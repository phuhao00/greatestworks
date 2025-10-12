// Package main 游戏服务主程序
// 基于DDD架构的分布式游戏服务
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"greatestworks/internal/application/handlers"
	"greatestworks/internal/config"
	"greatestworks/internal/infrastructure/logging"
	"greatestworks/internal/infrastructure/monitoring"
	"greatestworks/internal/interfaces/http"
	"greatestworks/internal/interfaces/rpc"
)

// GameServiceConfig 游戏服务配置
type GameServiceConfig = config.Config

// GameService 游戏服务
type GameService struct {
	config     atomic.Pointer[GameServiceConfig]
	logger     logging.Logger
	httpServer *http.Server
	rpcServer  *rpc.RPCServer
	profiler   *monitoring.Profiler
	commandBus *handlers.CommandBus
	queryBus   *handlers.QueryBus
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewGameService 创建游戏服务
func NewGameService(config *GameServiceConfig, logger logging.Logger) *GameService {
	ctx, cancel := context.WithCancel(context.Background())

	// 创建命令和查询总线
	commandBus := handlers.NewCommandBus()
	queryBus := handlers.NewQueryBus()

	service := &GameService{
		logger:     logger,
		commandBus: commandBus,
		queryBus:   queryBus,
		ctx:        ctx,
		cancel:     cancel,
		profiler:   monitoring.NewProfiler(logger),
	}

	if config != nil {
		service.config.Store(config)
	}

	return service
}

// UpdateConfig replaces the active configuration snapshot.
func (s *GameService) UpdateConfig(cfg *GameServiceConfig) {
	if cfg == nil {
		return
	}
	s.config.Store(cfg)
}

// Start 启动游戏服务
func (s *GameService) Start() error {
	cfg := s.config.Load()
	if cfg == nil {
		return fmt.Errorf("game service configuration not loaded")
	}

	s.logger.Info("Starting game service", logging.Fields{
		"service": cfg.Service.Name,
		"version": cfg.Service.Version,
		"node_id": cfg.Service.NodeID,
	})

	// 初始化数据库连接
	if err := s.initializeDatabase(cfg); err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}

	// 初始化消息队列
	if err := s.initializeMessaging(cfg); err != nil {
		return fmt.Errorf("初始化消息队列失败: %w", err)
	}

	// 初始化应用服务
	if err := s.initializeApplicationServices(cfg); err != nil {
		return fmt.Errorf("初始化应用服务失败: %w", err)
	}

	// 初始化HTTP服务器
	if err := s.initializeHTTPServer(cfg); err != nil {
		return fmt.Errorf("初始化HTTP服务器失败: %w", err)
	}

	// 初始化RPC服务器
	if err := s.initializeRPCServer(cfg); err != nil {
		return fmt.Errorf("初始化RPC服务器失败: %w", err)
	}

	// 启动HTTP服务器
	go func() {
		if err := s.httpServer.Start(); err != nil {
			s.logger.Error("HTTP server start failed", err)
		}
	}()

	// 启动RPC服务器
	go func() {
		if err := s.rpcServer.Start(); err != nil {
			s.logger.Error("RPC server start failed", err)
		}
	}()

	if cfg.Monitoring.Profiling.Enabled {
		host := cfg.Monitoring.Profiling.Host
		if host == "" {
			host = cfg.Server.HTTP.Host
		}

		if cfg.Monitoring.Profiling.Port == 0 {
			s.logger.Warn("pprof未启动: 未配置端口")
		} else if host == cfg.Server.HTTP.Host && cfg.Monitoring.Profiling.Port == cfg.Server.HTTP.Port {
			s.logger.Info("pprof routes enabled on primary HTTP server", logging.Fields{
				"addr": fmt.Sprintf("%s:%d", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port),
				"path": "/debug/pprof/",
			})
		} else if err := s.profiler.Start(host, cfg.Monitoring.Profiling.Port); err != nil {
			s.logger.Error("Failed to start pprof server", err, logging.Fields{
				"host": host,
				"port": cfg.Monitoring.Profiling.Port,
			})
		}
	}

	s.logger.Info("Game service started successfully", logging.Fields{
		"http_addr": fmt.Sprintf("%s:%d", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port),
		"rpc_addr":  fmt.Sprintf("%s:%d", cfg.Server.RPC.Host, cfg.Server.RPC.Port),
	})

	return nil
}

// Stop 停止游戏服务
func (s *GameService) Stop() error {
	s.logger.Info("停止游戏服务")

	// 取消上下文
	s.cancel()

	// 停止HTTP服务器
	if s.httpServer != nil {
		if err := s.httpServer.Stop(); err != nil {
			s.logger.Error("Failed to stop HTTP server", err)
			return err
		}
	}

	// 停止RPC服务器
	if s.rpcServer != nil {
		if err := s.rpcServer.Stop(); err != nil {
			s.logger.Error("Failed to stop RPC server", err)
			return err
		}
	}

	if s.profiler != nil {
		if err := s.profiler.Stop(context.Background()); err != nil {
			s.logger.Error("Failed to stop pprof server", err)
			return err
		}
	}

	s.logger.Info("游戏服务已停止")
	return nil
}

// initializeDatabase 初始化数据库连接
func (s *GameService) initializeDatabase(cfg *GameServiceConfig) error {
	_ = cfg
	s.logger.Info("初始化数据库连接")

	// TODO: 实现MongoDB连接
	// TODO: 实现Redis连接

	s.logger.Info("数据库连接初始化完成")
	return nil
}

// initializeMessaging 初始化消息队列
func (s *GameService) initializeMessaging(cfg *GameServiceConfig) error {
	_ = cfg
	s.logger.Info("初始化消息队列")

	// TODO: 实现NATS连接
	// TODO: 实现JetStream配置

	s.logger.Info("消息队列初始化完成")
	return nil
}

// initializeApplicationServices 初始化应用服务
func (s *GameService) initializeApplicationServices(cfg *GameServiceConfig) error {
	_ = cfg
	s.logger.Info("初始化应用服务")

	// TODO: 初始化领域服务
	// TODO: 初始化仓储
	// TODO: 初始化事件处理器

	s.logger.Info("应用服务初始化完成")
	return nil
}

// initializeHTTPServer 初始化HTTP服务器
func (s *GameService) initializeHTTPServer(cfg *GameServiceConfig) error {
	s.logger.Info("初始化HTTP服务器")

	// 创建HTTP服务器配置
	httpConfig := &http.ServerConfig{
		Host:         cfg.Server.HTTP.Host,
		Port:         cfg.Server.HTTP.Port,
		ReadTimeout:  cfg.Server.HTTP.ReadTimeout,
		WriteTimeout: cfg.Server.HTTP.WriteTimeout,
		IdleTimeout:  cfg.Server.HTTP.IdleTimeout,
	}

	// 创建HTTP服务器
	s.httpServer = http.NewServer(httpConfig, s.logger)

	if cfg.Monitoring.Profiling.Enabled &&
		cfg.Monitoring.Profiling.Host == cfg.Server.HTTP.Host &&
		cfg.Monitoring.Profiling.Port == cfg.Server.HTTP.Port {
		s.httpServer.EnableProfiling()
	}

	s.logger.Info("HTTP服务器初始化完成")
	return nil
}

// initializeRPCServer 初始化RPC服务器
func (s *GameService) initializeRPCServer(cfg *GameServiceConfig) error {
	s.logger.Info("初始化RPC服务器")

	// 创建RPC服务器配置
	rpcConfig := &rpc.RPCServerConfig{
		Host:            cfg.Server.RPC.Host,
		Port:            cfg.Server.RPC.Port,
		MaxConnections:  cfg.Server.RPC.MaxConnections,
		Timeout:         cfg.Server.RPC.Timeout,
		KeepAlive:       cfg.Server.RPC.KeepAlive,
		KeepAlivePeriod: cfg.Server.RPC.KeepAlivePeriod,
		ReadTimeout:     cfg.Server.RPC.ReadTimeout,
		WriteTimeout:    cfg.Server.RPC.WriteTimeout,
	}

	// 创建RPC服务器
	s.rpcServer = rpc.NewRPCServer(rpcConfig, s.commandBus, s.queryBus, s.logger)

	s.logger.Info("RPC服务器初始化完成")
	return nil
}

// loadConfig 加载配置
func loadInitialConfig() (*GameServiceConfig, []string, *config.Loader, error) {
	loader := config.NewLoader(
		config.WithService("game-service"),
	)

	cfg, files, err := loader.Load()
	if err != nil {
		return nil, nil, nil, err
	}

	return cfg, files, loader, nil
}

// main 主函数
func main() {
	logger := logging.NewBaseLogger(logging.InfoLevel)
	logger.Info("启动游戏服务", logging.Fields{})

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
	service := NewGameService(runtimeCfg, logger)

	manager.OnChange(func(next *config.Config) {
		if next == nil {
			return
		}
		service.UpdateConfig(next)
		logger.Info("游戏服务配置已刷新", logging.Fields{
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
		log.Fatalf("启动游戏服务失败: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	select {
	case sig := <-sigChan:
		logger.Info("收到关闭信号", logging.Fields{
			"signal": sig.String(),
		})
	case <-service.ctx.Done():
		logger.Info("上下文已取消", logging.Fields{})
	}

	logger.Info("正在关闭游戏服务...", logging.Fields{})
	watchCancel()
	if err := service.Stop(); err != nil {
		logger.Error("关闭游戏服务失败", err, logging.Fields{})
		os.Exit(1)
	}

	logger.Info("游戏服务已关闭", logging.Fields{})
}
