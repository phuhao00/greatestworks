// Package main 游戏服务器主程序
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"greatestworks/application/handlers"
	"greatestworks/internal/infrastructure/logger"
	"greatestworks/internal/interfaces/tcp/connection"

	"greatestworks/internal/interfaces/http"
	"greatestworks/internal/interfaces/tcp"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	HTTP *http.HTTPServerConfig `yaml:"http" json:"http"`
	TCP  *tcp.ServerConfig      `yaml:"tcp" json:"tcp"`
}

// MultiProtocolServer 多协议服务器
type MultiProtocolServer struct {
	config     *ServerConfig
	httpServer *http.HTTPServer
	tcpServer  *tcp.TCPServer

	commandBus   *handlers.CommandBus
	queryBus     *handlers.QueryBus
	logger       logger.Logger
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	shutdownChan chan os.Signal
}

// NewMultiProtocolServer 创建多协议服务器
func NewMultiProtocolServer(config *ServerConfig, logger logger.Logger) *MultiProtocolServer {
	ctx, cancel := context.WithCancel(context.Background())

	// 创建命令和查询总线
	commandBus := handlers.NewCommandBus()
	queryBus := handlers.NewQueryBus()

	// 创建各协议服务器
	services := &http.ServiceContainer{
		CommandBus: commandBus,
		QueryBus:   queryBus,
	}
	httpServer, _ := http.NewHTTPServer(config.HTTP, services, logger)
	tcpServer := tcp.NewTCPServer(config.TCP, commandBus, queryBus, logger)

	return &MultiProtocolServer{
		config:     config,
		httpServer: httpServer,
		tcpServer:  tcpServer,

		commandBus:   commandBus,
		queryBus:     queryBus,
		logger:       logger,
		ctx:          ctx,
		cancel:       cancel,
		shutdownChan: make(chan os.Signal, 1),
	}
}

// Start 启动所有服务器
func (s *MultiProtocolServer) Start() error {
	s.logger.Info("Starting multi-protocol server")

	// 启动HTTP服务器
	if s.config.HTTP != nil {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			if err := s.httpServer.Start(s.ctx); err != nil {
				s.logger.Error("HTTP server failed", "error", err)
			}
		}()
		s.logger.Info("HTTP server started", "address", fmt.Sprintf("%s:%d", s.config.HTTP.Host, s.config.HTTP.Port))
	}

	// 启动TCP服务器
	if s.config.TCP != nil {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			if err := s.tcpServer.Start(); err != nil {
				s.logger.Error("TCP server failed", "error", err)
			}
		}()
		s.logger.Info("TCP server started", "address", s.config.TCP.Addr)
	}

	// 等待一段时间确保所有服务器启动
	time.Sleep(1 * time.Second)

	s.logger.Info("All servers started successfully")
	return nil
}

// Stop 停止所有服务器
func (s *MultiProtocolServer) Stop() error {
	s.logger.Info("Stopping multi-protocol server")

	// 取消上下文
	s.cancel()

	// 停止所有服务器
	var stopErrors []error

	if s.httpServer != nil {
		if err := s.httpServer.Stop(); err != nil {
			s.logger.Error("Failed to stop HTTP server", "error", err)
			stopErrors = append(stopErrors, err)
		}
	}

	if s.tcpServer != nil {
		if err := s.tcpServer.Stop(); err != nil {
			s.logger.Error("Failed to stop TCP server", "error", err)
			stopErrors = append(stopErrors, err)
		}
	}

	// 等待所有协程结束
	s.wg.Wait()

	if len(stopErrors) > 0 {
		s.logger.Error("Some servers failed to stop gracefully", "error_count", len(stopErrors))
		return fmt.Errorf("failed to stop %d servers", len(stopErrors))
	}

	s.logger.Info("All servers stopped successfully")
	return nil
}

// Wait 等待关闭信号
func (s *MultiProtocolServer) Wait() {
	// 监听关闭信号
	signal.Notify(s.shutdownChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	select {
	case sig := <-s.shutdownChan:
		s.logger.Info("Received shutdown signal", "signal", sig.String())
	case <-s.ctx.Done():
		s.logger.Info("Context cancelled")
	}
}

// GetStats 获取所有服务器统计信息
func (s *MultiProtocolServer) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})

	if s.httpServer != nil {
		stats["http"] = map[string]interface{}{
			"status":  "running",
			"address": fmt.Sprintf("%s:%d", s.config.HTTP.Host, s.config.HTTP.Port),
		}
	}

	if s.tcpServer != nil {
		stats["tcp"] = s.tcpServer.GetStats()
	}

	return stats
}

// loadConfig 加载配置
func loadConfig() (*ServerConfig, error) {
	// 默认配置
	defaultConfig := &ServerConfig{
		HTTP: &http.HTTPServerConfig{
			Host:              "0.0.0.0",
			Port:              8080,
			ReadTimeout:       30 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       60 * time.Second,
			MaxHeaderBytes:    1 << 20, // 1MB
			EnableCORS:        true,
			EnableMetrics:     true,
			EnableRequestID:   true,
			EnableLogging:     true,
			EnableRecovery:    true,
			RateLimitEnabled:  true,
			RateLimitRequests: 100,
			RateLimitDuration: time.Minute,
		},
		TCP: &tcp.ServerConfig{
			Addr:           ":9090",
			MaxConnections: 10000,
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			HeartbeatConfig: &connection.HeartbeatConfig{
				Interval: 30 * time.Second,
				Timeout:  60 * time.Second,
			},
			EnableCompression: false,
			BufferSize:        4096,
		},
	}

	// TODO: 从配置文件或环境变量加载配置
	// 这里可以添加配置文件加载逻辑

	return defaultConfig, nil
}

// initializeServices 初始化服务
func initializeServices(logger logger.Logger) error {
	// TODO: 初始化数据库连接
	// TODO: 初始化缓存
	// TODO: 初始化其他依赖服务

	logger.Info("Services initialized successfully")
	return nil
}

// main 主函数
func main() {
	// 创建日志器
	logger := logger.NewLogger()

	logger.Info("Starting Greatest Works Game Server")

	// 加载配置
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化服务
	if err := initializeServices(logger); err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}

	// 创建多协议服务器
	server := NewMultiProtocolServer(config, logger)

	// 启动服务器
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// 打印服务器信息
	logger.Info("Server startup completed",
		"http_addr", fmt.Sprintf("%s:%d", config.HTTP.Host, config.HTTP.Port),
		"tcp_addr", config.TCP.Addr)

	// 等待关闭信号
	server.Wait()

	// 优雅关闭
	logger.Info("Shutting down server...")
	if err := server.Stop(); err != nil {
		logger.Error("Failed to stop server gracefully", "error", err)
		os.Exit(1)
	}

	logger.Info("Server shutdown completed")
}
