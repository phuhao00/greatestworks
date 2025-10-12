// Package http 提供HTTP服务器实现
package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"greatestworks/internal/infrastructure/logging"
)

// ServerConfig HTTP服务器配置
type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// Server HTTP服务器
type Server struct {
	config *ServerConfig
	logger logging.Logger
	server *http.Server
	ctx    context.Context
	cancel context.CancelFunc
}

// NewServer 创建HTTP服务器
func NewServer(config *ServerConfig, logger logging.Logger) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		config: config,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start 启动HTTP服务器
func (s *Server) Start() error {
	// 创建路由
	mux := http.NewServeMux()

	// 健康检查端点
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/ready", s.readyHandler)
	mux.HandleFunc("/metrics", s.metricsHandler)

	// 创建HTTP服务器
	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		Handler:      mux,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	// 启动服务器
	go func() {
		s.logger.Info("HTTP server starting", logging.Fields{
			"addr": s.server.Addr,
		})

		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("HTTP server failed", err, logging.Fields{
				"addr": s.server.Addr,
			})
		}
	}()

	return nil
}

// Stop 停止HTTP服务器
func (s *Server) Stop() error {
	s.logger.Info("Stopping HTTP server")

	// 取消上下文
	s.cancel()

	// 创建关闭上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭服务器
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("HTTP server shutdown failed", err)
		return err
	}

	s.logger.Info("HTTP server stopped")
	return nil
}

// healthHandler 健康检查处理器
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
}

// readyHandler 就绪检查处理器
func (s *Server) readyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ready","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
}

// metricsHandler 指标处理器
func (s *Server) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("# HELP mmo_server_info Server information\n# TYPE mmo_server_info gauge\nmmo_server_info{version=\"1.0.0\"} 1\n"))
}
