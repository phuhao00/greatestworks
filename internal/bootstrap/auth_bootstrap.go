package bootstrap

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"greatestworks/internal/config"
	"greatestworks/internal/database"
	"greatestworks/internal/events"
	"greatestworks/internal/infrastructure/logging"
	"greatestworks/internal/infrastructure/messaging"
	"greatestworks/internal/infrastructure/monitoring"
	httpiface "greatestworks/internal/interfaces/http"
)

// AuthBootstrap wires infrastructure for the auth service
type AuthBootstrap struct {
	config     atomic.Pointer[config.Config]
	logger     logging.Logger
	httpServer *httpiface.Server
	profiler   *monitoring.Profiler

	// infra
	mongoClient *mongo.Client
	redisClient *redis.Client
	eventBus    *events.EventBus

	ctx    context.Context
	cancel context.CancelFunc
}

func NewAuthBootstrap(cfg *config.Config, logger logging.Logger) *AuthBootstrap {
	ctx, cancel := context.WithCancel(context.Background())
	b := &AuthBootstrap{logger: logger, ctx: ctx, cancel: cancel}
	if cfg != nil {
		b.config.Store(cfg)
	}
	return b
}

func (s *AuthBootstrap) UpdateConfig(cfg *config.Config) {
	if cfg != nil {
		s.config.Store(cfg)
	}
}

func (s *AuthBootstrap) Start() error {
	cfg := s.config.Load()
	if cfg == nil {
		return fmt.Errorf("auth service configuration not loaded")
	}

	s.logger.Info("Starting auth service", logging.Fields{
		"service": cfg.Service.Name,
		"version": cfg.Service.Version,
		"node_id": cfg.Service.NodeID,
	})

	if err := s.initializeInfrastructure(cfg); err != nil {
		return fmt.Errorf("初始化基础设施失败: %w", err)
	}
	if err := s.initializeHTTPServer(cfg); err != nil {
		return fmt.Errorf("初始化HTTP服务器失败: %w", err)
	}

	go func() {
		if err := s.httpServer.Start(); err != nil {
			s.logger.Error("HTTP server start failed", err)
		}
	}()

	s.profiler = monitoring.NewProfiler(s.logger)
	if cfg.Monitoring.Profiling.Enabled {
		host := cfg.Monitoring.Profiling.Host
		if host == "" {
			host = cfg.Server.HTTP.Host
		}
		if cfg.Monitoring.Profiling.Port == 0 {
			s.logger.Warn("pprof未启动: 未配置端口")
		} else if host == cfg.Server.HTTP.Host && cfg.Monitoring.Profiling.Port == cfg.Server.HTTP.Port {
			s.logger.Info("pprof routes enabled on primary HTTP server", logging.Fields{"addr": fmt.Sprintf("%s:%d", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port), "path": "/debug/pprof/"})
		} else if err := s.profiler.Start(host, cfg.Monitoring.Profiling.Port); err != nil {
			s.logger.Error("Failed to start pprof server", err, logging.Fields{"host": host, "port": cfg.Monitoring.Profiling.Port})
		}
	}

	s.logger.Info("Auth service started successfully", logging.Fields{
		"http_addr": fmt.Sprintf("%s:%d", cfg.Server.HTTP.Host, cfg.Server.HTTP.Port),
	})
	return nil
}

func (s *AuthBootstrap) Stop() error {
	s.logger.Info("停止认证服务")
	s.cancel()
	if s.httpServer != nil {
		if err := s.httpServer.Stop(); err != nil {
			s.logger.Error("Failed to stop HTTP server", err)
			return err
		}
	}
	if s.profiler != nil {
		if err := s.profiler.Stop(context.Background()); err != nil {
			s.logger.Error("Failed to stop pprof server", err)
			return err
		}
	}
	if s.mongoClient != nil {
		if err := s.mongoClient.Disconnect(s.ctx); err != nil {
			s.logger.Error("Failed to disconnect MongoDB", err)
		}
	}
	if s.redisClient != nil {
		if err := s.redisClient.Close(); err != nil {
			s.logger.Error("Failed to close Redis", err)
		}
	}
	if s.eventBus != nil {
		s.eventBus.Close()
	}
	s.logger.Info("认证服务已停止")
	return nil
}

func (s *AuthBootstrap) initializeInfrastructure(cfg *config.Config) error {
	s.logger.Info("初始化基础设施层")
	// Mongo
	mongoConfig := &database.MongoConfig{
		URI: cfg.Database.MongoDB.URI, Database: cfg.Database.MongoDB.Database,
		MaxPoolSize: uint64(cfg.Database.MongoDB.MaxPoolSize), MinPoolSize: uint64(cfg.Database.MongoDB.MinPoolSize),
		MaxIdleTime: int(cfg.Database.MongoDB.MaxIdleTime / time.Second), ConnectTimeout: int(cfg.Database.MongoDB.ConnectTimeout / time.Second), SocketTimeout: int(cfg.Database.MongoDB.SocketTimeout / time.Second),
	}
	mongoDB := database.NewMongoDB(mongoConfig)
	if err := mongoDB.Connect(s.ctx); err != nil {
		return fmt.Errorf("连接MongoDB失败: %w", err)
	}
	s.mongoClient = mongoDB.GetClient()
	s.logger.Info("MongoDB连接成功", logging.Fields{"database": mongoConfig.Database})

	// Redis
	redisConfig := &database.RedisConfig{
		Addr: cfg.Database.Redis.Addr, Password: cfg.Database.Redis.Password, DB: cfg.Database.Redis.DB,
		PoolSize: cfg.Database.Redis.PoolSize, MinIdleConns: cfg.Database.Redis.MinIdleConns,
		DialTimeout: int(cfg.Database.Redis.DialTimeout / time.Second), ReadTimeout: int(cfg.Database.Redis.ReadTimeout / time.Second), WriteTimeout: int(cfg.Database.Redis.WriteTimeout / time.Second),
	}
	redisDB := database.NewRedis(redisConfig)
	if err := redisDB.Connect(s.ctx); err != nil {
		return fmt.Errorf("连接Redis失败: %w", err)
	}
	s.redisClient = redisDB.GetClient()
	s.logger.Info("Redis连接成功", logging.Fields{"addr": redisConfig.Addr, "db": redisConfig.DB})

	// Event bus (optional)
	eventLogger := messaging.NewEventLoggerAdapter(s.logger)
	s.eventBus = events.NewEventBus(eventLogger)
	if cfg.Messaging.NATS.URL != "" {
		if err := s.eventBus.ConnectNATS(cfg.Messaging.NATS.URL); err != nil {
			s.logger.Error("连接NATS失败", err, logging.Fields{"url": cfg.Messaging.NATS.URL})
		} else {
			s.logger.Info("NATS连接成功", logging.Fields{"url": cfg.Messaging.NATS.URL})
		}
	}
	s.logger.Info("基础设施层初始化完成")
	return nil
}

func (s *AuthBootstrap) initializeHTTPServer(cfg *config.Config) error {
	s.logger.Info("初始化HTTP服务器")
	httpConfig := &httpiface.ServerConfig{
		Host: cfg.Server.HTTP.Host, Port: cfg.Server.HTTP.Port,
		ReadTimeout: cfg.Server.HTTP.ReadTimeout, WriteTimeout: cfg.Server.HTTP.WriteTimeout, IdleTimeout: cfg.Server.HTTP.IdleTimeout,
	}
	s.httpServer = httpiface.NewServer(httpConfig, s.logger)
	if cfg.Monitoring.Profiling.Enabled && cfg.Monitoring.Profiling.Host == cfg.Server.HTTP.Host && cfg.Monitoring.Profiling.Port == cfg.Server.HTTP.Port {
		s.httpServer.EnableProfiling()
	}
	s.logger.Info("HTTP服务器初始化完成")
	return nil
}

// Done returns a channel that's closed when the service context is canceled.
func (s *AuthBootstrap) Done() <-chan struct{} { return s.ctx.Done() }
