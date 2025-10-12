// Package main 认证服务主程序
// 基于DDD架构的分布式认证服务
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
	"greatestworks/internal/interfaces/http"
)

// AuthServiceConfig aliases the shared configuration schema for readability.
type AuthServiceConfig = config.Config

// AuthService 认证服务
type AuthService struct {
	config atomic.Pointer[AuthServiceConfig]
	logger logging.Logger
	server *http.Server
	ctx    context.Context
	cancel context.CancelFunc
}

// NewAuthService 创建认证服务
func NewAuthService(cfg *AuthServiceConfig, logger logging.Logger) *AuthService {
	ctx, cancel := context.WithCancel(context.Background())

	service := &AuthService{
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}

	if cfg != nil {
		service.config.Store(cfg)
	}

	return service
}

// UpdateConfig replaces the in-memory configuration snapshot.
func (s *AuthService) UpdateConfig(cfg *AuthServiceConfig) {
	if cfg == nil {
		return
	}
	s.config.Store(cfg)
}

// Start 启动认证服务
func (s *AuthService) Start() error {
	cfg := s.config.Load()
	if cfg == nil {
		return fmt.Errorf("auth service configuration not loaded")
	}

	s.logger.Info("Starting auth service", logging.Fields{
		"service": cfg.Service.Name,
		"version": cfg.Service.Version,
		"node_id": cfg.Service.NodeID,
	})

	if err := s.initializeDatabase(cfg); err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}

	if err := s.initializeHTTPServer(cfg); err != nil {
		return fmt.Errorf("初始化HTTP服务器失败: %w", err)
	}

	go func() {
		if err := s.server.Start(); err != nil {
			s.logger.Error("HTTP server start failed", err)
		}
	}()

	s.logger.Info("Auth service started successfully", logging.Fields{
		"http_addr": fmt.Sprintf("%s:%d", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port),
	})

	return nil
}

// Stop 停止认证服务
func (s *AuthService) Stop() error {
	s.logger.Info("停止认证服务")

	// 取消上下文
	s.cancel()

	// 停止HTTP服务器
	if s.server != nil {
		if err := s.server.Stop(); err != nil {
			s.logger.Error("Failed to stop HTTP server", err)
			return err
		}
	}

	s.logger.Info("认证服务已停止")
	return nil
}

// initializeDatabase 初始化数据库连接
func (s *AuthService) initializeDatabase(cfg *AuthServiceConfig) error {
	s.logger.Info("初始化数据库连接")

	// TODO: 实现MongoDB连接，使用 cfg.Database.MongoDB
	// TODO: 实现Redis连接，使用 cfg.Database.Redis

	s.logger.Info("数据库连接初始化完成")
	return nil
}

// initializeHTTPServer 初始化HTTP服务器
func (s *AuthService) initializeHTTPServer(cfg *AuthServiceConfig) error {
	s.logger.Info("初始化HTTP服务器")

	httpConfig := &http.ServerConfig{
		Host:         cfg.Server.HTTP.Host,
		Port:         cfg.Server.HTTP.Port,
		ReadTimeout:  cfg.Server.HTTP.ReadTimeout,
		WriteTimeout: cfg.Server.HTTP.WriteTimeout,
		IdleTimeout:  cfg.Server.HTTP.IdleTimeout,
	}

	// TODO: 根据配置启用CORS、Swagger、限流等特性

	s.server = http.NewServer(httpConfig, s.logger)

	s.logger.Info("HTTP服务器初始化完成")
	return nil
}

// loadInitialConfig 加载配置并返回配置与文件来源
func loadInitialConfig() (*AuthServiceConfig, []string, *config.Loader, error) {
	loader := config.NewLoader(
		config.WithService("auth-service"),
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
	logger.Info("启动认证服务")

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
	service := NewAuthService(runtimeCfg, logger)

	manager.OnChange(func(next *config.Config) {
		if next == nil {
			return
		}
		service.UpdateConfig(next)
		logger.Info("认证服务配置已刷新", logging.Fields{
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
		log.Fatalf("启动认证服务失败: %v", err)
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

	logger.Info("正在关闭认证服务...")
	watchCancel()
	if err := service.Stop(); err != nil {
		logger.Error("关闭认证服务失败", err, logging.Fields{})
		os.Exit(1)
	}

	logger.Info("认证服务已关闭")
}
