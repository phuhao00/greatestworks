// Package main 游戏服务主程序
// 基于DDD架构的分布式游戏服务
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"greatestworks/application/handlers"
	"greatestworks/internal/infrastructure/logging"
	"greatestworks/internal/interfaces/rpc"
)

// GameServiceConfig 游戏服务配置
type GameServiceConfig struct {
	Service struct {
		Name        string `yaml:"name"`
		Version     string `yaml:"version"`
		Environment string `yaml:"environment"`
		NodeID      string `yaml:"node_id"`
	} `yaml:"service"`

	Server struct {
		RPC struct {
			Host            string        `yaml:"host"`
			Port            int           `yaml:"port"`
			MaxConnections  int           `yaml:"max_connections"`
			Timeout         time.Duration `yaml:"timeout"`
			KeepAlive       bool          `yaml:"keep_alive"`
			KeepAlivePeriod time.Duration `yaml:"keep_alive_period"`
			ReadTimeout     time.Duration `yaml:"read_timeout"`
			WriteTimeout    time.Duration `yaml:"write_timeout"`
		} `yaml:"rpc"`
	} `yaml:"server"`

	Database struct {
		MongoDB struct {
			URI            string        `yaml:"uri"`
			Database       string        `yaml:"database"`
			Username       string        `yaml:"username"`
			Password       string        `yaml:"password"`
			AuthSource     string        `yaml:"auth_source"`
			MaxPoolSize    int           `yaml:"max_pool_size"`
			MinPoolSize    int           `yaml:"min_pool_size"`
			MaxIdleTime    time.Duration `yaml:"max_idle_time"`
			ConnectTimeout time.Duration `yaml:"connect_timeout"`
			SocketTimeout  time.Duration `yaml:"socket_timeout"`
		} `yaml:"mongodb"`

		Redis struct {
			Addr         string        `yaml:"addr"`
			Password     string        `yaml:"password"`
			DB           int           `yaml:"db"`
			PoolSize     int           `yaml:"pool_size"`
			MinIdleConns int           `yaml:"min_idle_conns"`
			MaxRetries   int           `yaml:"max_retries"`
			DialTimeout  time.Duration `yaml:"dial_timeout"`
			ReadTimeout  time.Duration `yaml:"read_timeout"`
			WriteTimeout time.Duration `yaml:"write_timeout"`
			PoolTimeout  time.Duration `yaml:"pool_timeout"`
			IdleTimeout  time.Duration `yaml:"idle_timeout"`
		} `yaml:"redis"`
	} `yaml:"database"`

	Messaging struct {
		NATS struct {
			URL           string        `yaml:"url"`
			ClusterID     string        `yaml:"cluster_id"`
			ClientID      string        `yaml:"client_id"`
			MaxReconnect  int           `yaml:"max_reconnect"`
			ReconnectWait time.Duration `yaml:"reconnect_wait"`
			Timeout       time.Duration `yaml:"timeout"`
			JetStream     struct {
				Enabled bool   `yaml:"enabled"`
				Domain  string `yaml:"domain"`
			} `yaml:"jetstream"`
			Subjects struct {
				PlayerEvents string `yaml:"player_events"`
				BattleEvents string `yaml:"battle_events"`
				SystemEvents string `yaml:"system_events"`
				DomainEvents string `yaml:"domain_events"`
			} `yaml:"subjects"`
		} `yaml:"nats"`
	} `yaml:"messaging"`

	Domain struct {
		Player struct {
			MaxLevel          int           `yaml:"max_level"`
			InitialGold       int           `yaml:"initial_gold"`
			InitialExperience int           `yaml:"initial_experience"`
			MaxInventorySlots int           `yaml:"max_inventory_slots"`
			SessionTimeout    time.Duration `yaml:"session_timeout"`
		} `yaml:"player"`

		Battle struct {
			MaxBattleTime      time.Duration `yaml:"max_battle_time"`
			DamageVariance     float64       `yaml:"damage_variance"`
			CriticalRateBase   float64       `yaml:"critical_rate_base"`
			CriticalDamageBase float64       `yaml:"critical_damage_base"`
			MaxParticipants    int           `yaml:"max_participants"`
			TurnTimeout        time.Duration `yaml:"turn_timeout"`
		} `yaml:"battle"`

		Experience struct {
			BaseExpPerLevel int     `yaml:"base_exp_per_level"`
			ExpMultiplier   float64 `yaml:"exp_multiplier"`
			MaxExpBonus     float64 `yaml:"max_exp_bonus"`
		} `yaml:"experience"`

		Chat struct {
			MaxMessageLength int      `yaml:"max_message_length"`
			RateLimit        int      `yaml:"rate_limit"`
			BannedWords      []string `yaml:"banned_words"`
		} `yaml:"chat"`

		Ranking struct {
			MaxEntries     int           `yaml:"max_entries"`
			UpdateInterval time.Duration `yaml:"update_interval"`
			CacheTTL       time.Duration `yaml:"cache_ttl"`
		} `yaml:"ranking"`

		Weather struct {
			UpdateInterval  time.Duration `yaml:"update_interval"`
			ForecastDays    int           `yaml:"forecast_days"`
			SeasonalEffects bool          `yaml:"seasonal_effects"`
		} `yaml:"weather"`

		Plant struct {
			GrowthSpeed  float64 `yaml:"growth_speed"`
			HarvestBonus float64 `yaml:"harvest_bonus"`
			MaxFarmSize  int     `yaml:"max_farm_size"`
		} `yaml:"plant"`
	} `yaml:"domain"`

	Application struct {
		CommandBus struct {
			Timeout       time.Duration `yaml:"timeout"`
			RetryAttempts int           `yaml:"retry_attempts"`
			RetryDelay    time.Duration `yaml:"retry_delay"`
		} `yaml:"command_bus"`

		QueryBus struct {
			Timeout  time.Duration `yaml:"timeout"`
			CacheTTL time.Duration `yaml:"cache_ttl"`
		} `yaml:"query_bus"`

		EventBus struct {
			Timeout         time.Duration `yaml:"timeout"`
			RetryAttempts   int           `yaml:"retry_attempts"`
			RetryDelay      time.Duration `yaml:"retry_delay"`
			DeadLetterQueue bool          `yaml:"dead_letter_queue"`
		} `yaml:"event_bus"`
	} `yaml:"application"`

	Logging struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
		Output string `yaml:"output"`
		Fields struct {
			Service string `yaml:"service"`
			Version string `yaml:"version"`
		} `yaml:"fields"`
	} `yaml:"logging"`
}

// GameService 游戏服务
type GameService struct {
	config     *GameServiceConfig
	logger     logging.Logger
	server     *rpc.RPCServer
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

	return &GameService{
		config:     config,
		logger:     logger,
		commandBus: commandBus,
		queryBus:   queryBus,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start 启动游戏服务
func (s *GameService) Start() error {
	s.logger.Info("Starting game service", logging.Fields{
		"service": s.config.Service.Name,
		"version": s.config.Service.Version,
		"node_id": s.config.Service.NodeID,
	})

	// 初始化数据库连接
	if err := s.initializeDatabase(); err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}

	// 初始化消息队列
	if err := s.initializeMessaging(); err != nil {
		return fmt.Errorf("初始化消息队列失败: %w", err)
	}

	// 初始化应用服务
	if err := s.initializeApplicationServices(); err != nil {
		return fmt.Errorf("初始化应用服务失败: %w", err)
	}

	// 初始化RPC服务器
	if err := s.initializeRPCServer(); err != nil {
		return fmt.Errorf("初始化RPC服务器失败: %w", err)
	}

	// 启动RPC服务器
	go func() {
		if err := s.server.Start(); err != nil {
			s.logger.Error("RPC server start failed", err)
		}
	}()

	s.logger.Info("Game service started successfully", logging.Fields{
		"rpc_addr": fmt.Sprintf("%s:%d", s.config.Server.RPC.Host, s.config.Server.RPC.Port),
	})

	return nil
}

// Stop 停止游戏服务
func (s *GameService) Stop() error {
	s.logger.Info("停止游戏服务")

	// 取消上下文
	s.cancel()

	// 停止RPC服务器
	if s.server != nil {
		if err := s.server.Stop(); err != nil {
			s.logger.Error("Failed to stop RPC server", err)
			return err
		}
	}

	s.logger.Info("游戏服务已停止")
	return nil
}

// initializeDatabase 初始化数据库连接
func (s *GameService) initializeDatabase() error {
	s.logger.Info("初始化数据库连接")

	// TODO: 实现MongoDB连接
	// TODO: 实现Redis连接

	s.logger.Info("数据库连接初始化完成")
	return nil
}

// initializeMessaging 初始化消息队列
func (s *GameService) initializeMessaging() error {
	s.logger.Info("初始化消息队列")

	// TODO: 实现NATS连接
	// TODO: 实现JetStream配置

	s.logger.Info("消息队列初始化完成")
	return nil
}

// initializeApplicationServices 初始化应用服务
func (s *GameService) initializeApplicationServices() error {
	s.logger.Info("初始化应用服务")

	// TODO: 初始化领域服务
	// TODO: 初始化仓储
	// TODO: 初始化事件处理器

	s.logger.Info("应用服务初始化完成")
	return nil
}

// initializeRPCServer 初始化RPC服务器
func (s *GameService) initializeRPCServer() error {
	s.logger.Info("初始化RPC服务器")

	// 创建RPC服务器配置
	rpcConfig := &rpc.RPCServerConfig{
		Host:            s.config.Server.RPC.Host,
		Port:            s.config.Server.RPC.Port,
		MaxConnections:  s.config.Server.RPC.MaxConnections,
		Timeout:         s.config.Server.RPC.Timeout,
		KeepAlive:       s.config.Server.RPC.KeepAlive,
		KeepAlivePeriod: s.config.Server.RPC.KeepAlivePeriod,
		ReadTimeout:     s.config.Server.RPC.ReadTimeout,
		WriteTimeout:    s.config.Server.RPC.WriteTimeout,
	}

	// 创建RPC服务器
	s.server = rpc.NewRPCServer(rpcConfig, s.commandBus, s.queryBus, s.logger)

	s.logger.Info("RPC服务器初始化完成")
	return nil
}

// loadConfig 加载配置
func loadConfig() (*GameServiceConfig, error) {
	// TODO: 从配置文件加载配置
	// 这里先返回默认配置
	config := &GameServiceConfig{
		Service: struct {
			Name        string `yaml:"name"`
			Version     string `yaml:"version"`
			Environment string `yaml:"environment"`
			NodeID      string `yaml:"node_id"`
		}{
			Name:        "game-service",
			Version:     "1.0.0",
			Environment: "development",
			NodeID:      "game-node-1",
		},
	}

	return config, nil
}

// main 主函数
func main() {
	// 创建日志器
	logger := logging.NewBaseLogger(logging.InfoLevel)

	logger.Info("启动游戏服务", logging.Fields{})

	// 加载配置
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 创建游戏服务
	service := NewGameService(config, logger)

	// 启动服务
	if err := service.Start(); err != nil {
		log.Fatalf("启动游戏服务失败: %v", err)
	}

	// 等待关闭信号
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

	// 优雅关闭
	logger.Info("正在关闭游戏服务...", logging.Fields{})
	if err := service.Stop(); err != nil {
		logger.Error("关闭游戏服务失败", err, logging.Fields{})
		os.Exit(1)
	}

	logger.Info("游戏服务已关闭", logging.Fields{})
}
