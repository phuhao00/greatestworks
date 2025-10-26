// Package rpc 提供Go原生RPC服务器实现
// 基于DDD架构的分布式游戏服务RPC接口
package rpc

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"time"

	"greatestworks/internal/application/handlers"
	"greatestworks/internal/infrastructure/logging"
)

// RPCServerConfig RPC服务器配置
type RPCServerConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	MaxConnections  int           `yaml:"max_connections"`
	Timeout         time.Duration `yaml:"timeout"`
	KeepAlive       bool          `yaml:"keep_alive"`
	KeepAlivePeriod time.Duration `yaml:"keep_alive_period"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
}

// RPCServer RPC服务器
type RPCServer struct {
	config     *RPCServerConfig
	logger     logging.Logger
	server     *rpc.Server
	commandBus *handlers.CommandBus
	queryBus   *handlers.QueryBus
	listener   net.Listener
	ctx        context.Context
	cancel     context.CancelFunc
	extras     []interface{}
}

// NewRPCServer 创建RPC服务器
func NewRPCServer(
	config *RPCServerConfig,
	commandBus *handlers.CommandBus,
	queryBus *handlers.QueryBus,
	logger logging.Logger,
) *RPCServer {
	ctx, cancel := context.WithCancel(context.Background())

	return &RPCServer{
		config:     config,
		logger:     logger,
		commandBus: commandBus,
		queryBus:   queryBus,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start 启动RPC服务器
func (s *RPCServer) Start() error {
	// 创建RPC服务器
	s.server = rpc.NewServer()

	// 注册服务
	s.registerServices()

	// 注册额外服务（由引导器注入）
	for _, svc := range s.extras {
		if err := s.server.Register(svc); err != nil {
			s.logger.Error("failed to register extra RPC service", err, logging.Fields{})
		}
	}

	// 创建监听器
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("创建监听器失败: %w", err)
	}
	s.listener = listener

	// 启动服务器
	go func() {
		s.logger.Info("RPC server started", logging.Fields{
			"address": addr,
		})
		s.server.Accept(listener)
	}()

	return nil
}

// RegisterService 允许引导器注册额外的RPC服务
func (s *RPCServer) RegisterService(service interface{}) {
	s.extras = append(s.extras, service)
}

// Stop 停止RPC服务器
func (s *RPCServer) Stop() error {
	s.logger.Info("停止RPC服务器")

	// 取消上下文
	s.cancel()

	// 关闭监听器
	if s.listener != nil {
		s.listener.Close()
	}

	s.logger.Info("RPC服务器已停止")
	return nil
}

// registerServices 注册服务
func (s *RPCServer) registerServices() {
	s.logger.Info("注册RPC服务")

	// 注册玩家服务
	playerService := NewPlayerRPCService(s.commandBus, s.queryBus, s.logger)
	s.server.Register(playerService)

	// 注册战斗服务
	battleService := NewBattleRPCService(s.commandBus, s.queryBus, s.logger)
	s.server.Register(battleService)

	// 注册排行榜服务
	rankingService := NewRankingRPCService(s.commandBus, s.queryBus, s.logger)
	s.server.Register(rankingService)

	// 注册其他领域服务
	// TODO: 注册更多服务

	s.logger.Info("RPC服务注册完成")
}

// GetStats 获取服务器统计信息
func (s *RPCServer) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})

	stats["status"] = "running"
	stats["address"] = fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	stats["max_connections"] = s.config.MaxConnections
	stats["timeout"] = s.config.Timeout.String()

	return stats
}

// DefaultRPCServerConfig 默认RPC服务器配置
func DefaultRPCServerConfig() *RPCServerConfig {
	return &RPCServerConfig{
		Host:            "0.0.0.0",
		Port:            8081,
		MaxConnections:  1000,
		Timeout:         30 * time.Second,
		KeepAlive:       true,
		KeepAlivePeriod: 30 * time.Second,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
	}
}
