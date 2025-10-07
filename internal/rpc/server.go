package rpc

import (
	"context"
	"fmt"
	"sync"

	"greatestworks/internal/infrastructure/logging"
)

// RPCServer netcore-go RPC服务�?
type RPCServer struct {
	services map[string]Service
	logger   logger.Logger
	mu       sync.RWMutex
}

// Service RPC服务接口
type Service interface {
	GetName() string
	HandleRequest(ctx context.Context, method string, data []byte) ([]byte, error)
}

// NewRPCServer 创建RPC服务�?
func NewRPCServer(logger logger.Logger) *RPCServer {
	return &RPCServer{
		services: make(map[string]Service),
		logger:   logger,
	}
}

// RegisterService 注册服务
func (s *RPCServer) RegisterService(service Service) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.services[service.GetName()] = service
	s.logger.Info("RPC服务已注�?, "service", service.GetName())
}

// UnregisterService 注销服务
func (s *RPCServer) UnregisterService(serviceName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.services, serviceName)
	s.logger.Info("RPC服务已注销", "service", serviceName)
}

// HandleRequest 处理RPC请求
func (s *RPCServer) HandleRequest(ctx context.Context, serviceName, method string, data []byte) ([]byte, error) {
	s.mu.RLock()
	service, exists := s.services[serviceName]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("服务不存�? %s", serviceName)
	}

	return service.HandleRequest(ctx, method, data)
}

// GetServices 获取所有注册的服务
func (s *RPCServer) GetServices() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	services := make([]string, 0, len(s.services))
	for name := range s.services {
		services = append(services, name)
	}

	return services
}

// Start 启动RPC服务�?
func (s *RPCServer) Start() error {
	s.logger.Info("RPC服务器启�?, "services", s.GetServices())
	return nil
}

// Stop 停止RPC服务�?
func (s *RPCServer) Stop() error {
	s.logger.Info("RPC服务器停�?)
	return nil
}
