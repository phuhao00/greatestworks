// Package http 提供HTTP服务器实现
package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"greatestworks/internal/infrastructure/logging"
	"greatestworks/internal/infrastructure/monitoring"
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
	config           *ServerConfig
	logger           logging.Logger
	server           *http.Server
	mux              *http.ServeMux
	ctx              context.Context
	cancel           context.CancelFunc
	profilingEnabled bool
	routes           []route
}

type route struct {
	method  string
	path    string
	handler http.HandlerFunc
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

// EnableProfiling 注册标准pprof处理器。
func (s *Server) EnableProfiling() {
	s.profilingEnabled = true
}

// Handle 注册业务路由
func (s *Server) Handle(method, path string, handler http.HandlerFunc) {
	s.routes = append(s.routes, route{method: method, path: path, handler: handler})
}

// Start 启动HTTP服务器
func (s *Server) Start() error {
	// 创建路由
	s.mux = http.NewServeMux()

	// 健康检查端点
	s.mux.HandleFunc("/health", s.healthHandler)
	s.mux.HandleFunc("/ready", s.readyHandler)

	if s.profilingEnabled {
		monitoring.RegisterHandlers(s.mux)
	}

	// 注册业务路由
	for _, r := range s.routes {
		handler := r.handler
		method := r.method
		s.mux.HandleFunc(r.path, func(w http.ResponseWriter, req *http.Request) {
			if method != "" && req.Method != method {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			handler(w, req)
		})
	}

	// 创建HTTP服务器
	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		Handler:      s.mux,
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
