package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"greatestworks/internal/application/handlers"
	appServices "greatestworks/internal/application/services"
	"greatestworks/internal/config"
	"greatestworks/internal/database"
	"greatestworks/internal/domain/character"
	"greatestworks/internal/infrastructure/logging"
	"greatestworks/internal/infrastructure/monitoring"
	"greatestworks/internal/infrastructure/persistence"
	"greatestworks/internal/interfaces/tcp"
	tcpProtocol "greatestworks/internal/interfaces/tcp/protocol"
)

// GatewayBootstrap wires infrastructure for the gateway service
type GatewayBootstrap struct {
	config    atomic.Pointer[config.Config]
	logger    logging.Logger
	tcpServer *tcp.TCPServer
	profiler  *monitoring.Profiler

	// infra
	mongoClient *mongo.Client
	redisClient *redis.Client

	// buses
	commandBus *handlers.CommandBus
	queryBus   *handlers.QueryBus

	// app services
	mapService       *appServices.MapService
	fightService     *appServices.FightService
	characterService *appServices.CharacterService
	updateMgr        *appServices.UpdateManager
	spawnMgr         *appServices.SpawnManager

	ctx    context.Context
	cancel context.CancelFunc
}

func NewGatewayBootstrap(cfg *config.Config, logger logging.Logger) *GatewayBootstrap {
	ctx, cancel := context.WithCancel(context.Background())
	b := &GatewayBootstrap{logger: logger, ctx: ctx, cancel: cancel}
	if cfg != nil {
		b.config.Store(cfg)
	}
	return b
}

func (s *GatewayBootstrap) UpdateConfig(cfg *config.Config) {
	if cfg != nil {
		s.config.Store(cfg)
	}
}

func (s *GatewayBootstrap) Start() error {
	cfg := s.config.Load()
	if cfg == nil {
		return fmt.Errorf("gateway service configuration not loaded")
	}

	s.logger.Info("Starting gateway service", logging.Fields{"service": cfg.Service.Name, "version": cfg.Service.Version, "node_id": cfg.Service.NodeID})

	if err := s.initializeInfrastructure(cfg); err != nil {
		return fmt.Errorf("初始化基础设施失败: %w", err)
	}
	if err := s.initializeApplicationLayer(cfg); err != nil {
		return fmt.Errorf("初始化应用服务层失败: %w", err)
	}
	if err := s.initializeTCPServer(cfg); err != nil {
		return fmt.Errorf("初始化TCP服务器失败: %w", err)
	}

	// Start runtime managers
	if s.updateMgr != nil {
		// Register map tick into update loop
		if s.mapService != nil {
			s.updateMgr.Register("map.tick", appServices.UpdateFunc(func(ctx context.Context, d time.Duration) error {
				s.mapService.Tick(ctx, d)
				return nil
			}))
		}
		s.updateMgr.Start(s.ctx)
	}
	if s.spawnMgr != nil {
		// default 2 workers; can be made configurable later
		s.spawnMgr.Start(s.ctx, 2)
	}

	go func() {
		if err := s.tcpServer.Start(); err != nil {
			s.logger.Error("TCP server start failed", err)
		}
	}()

	s.profiler = monitoring.NewProfiler(s.logger)
	if cfg.Monitoring.Profiling.Enabled {
		host := cfg.Monitoring.Profiling.Host
		if host == "" {
			host = cfg.Server.TCP.Host
		}
		if cfg.Monitoring.Profiling.Port == 0 {
			s.logger.Warn("pprof未启动: 未配置端口")
		} else if err := s.profiler.Start(host, cfg.Monitoring.Profiling.Port); err != nil {
			s.logger.Error("Failed to start pprof server", err, logging.Fields{"host": host, "port": cfg.Monitoring.Profiling.Port})
		}
	}

	s.logger.Info("Gateway service started successfully", logging.Fields{"tcp_addr": fmt.Sprintf("%s:%d", cfg.Server.TCP.Host, cfg.Server.TCP.Port)})
	return nil
}

func (s *GatewayBootstrap) Stop() error {
	s.logger.Info("停止网关服务")
	s.cancel()
	if s.updateMgr != nil {
		s.updateMgr.Stop()
	}
	if s.spawnMgr != nil {
		s.spawnMgr.Stop()
	}
	if s.tcpServer != nil {
		if err := s.tcpServer.Stop(); err != nil {
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
	if s.redisClient != nil {
		if err := s.redisClient.Close(); err != nil {
			s.logger.Error("Failed to close Redis", err)
		}
	}
	s.logger.Info("网关服务已停止")
	return nil
}

func (s *GatewayBootstrap) initializeInfrastructure(cfg *config.Config) error {
	s.logger.Info("初始化基础设施层")
	// Mongo
	mongoConfig := &database.MongoConfig{
		URI:            cfg.Database.MongoDB.URI,
		Database:       cfg.Database.MongoDB.Database,
		MaxPoolSize:    uint64(cfg.Database.MongoDB.MaxPoolSize),
		MinPoolSize:    uint64(cfg.Database.MongoDB.MinPoolSize),
		MaxIdleTime:    int(cfg.Database.MongoDB.MaxIdleTime / time.Second),
		ConnectTimeout: int(cfg.Database.MongoDB.ConnectTimeout / time.Second),
		SocketTimeout:  int(cfg.Database.MongoDB.SocketTimeout / time.Second),
	}
	mongoDB := database.NewMongoDB(mongoConfig)
	if err := mongoDB.Connect(s.ctx); err != nil {
		return fmt.Errorf("连接MongoDB失败: %w", err)
	}
	s.mongoClient = mongoDB.GetClient()
	s.logger.Info("MongoDB连接成功", logging.Fields{"database": mongoConfig.Database})

	// Redis
	redisConfig := &database.RedisConfig{
		Addr:         cfg.Database.Redis.Addr,
		Password:     cfg.Database.Redis.Password,
		DB:           cfg.Database.Redis.DB,
		PoolSize:     cfg.Database.Redis.PoolSize,
		MinIdleConns: cfg.Database.Redis.MinIdleConns,
		DialTimeout:  int(cfg.Database.Redis.DialTimeout / time.Second),
		ReadTimeout:  int(cfg.Database.Redis.ReadTimeout / time.Second),
		WriteTimeout: int(cfg.Database.Redis.WriteTimeout / time.Second),
	}
	redisDB := database.NewRedis(redisConfig)
	if err := redisDB.Connect(s.ctx); err != nil {
		return fmt.Errorf("连接Redis失败: %w", err)
	}
	s.redisClient = redisDB.GetClient()
	s.logger.Info("Redis连接成功", logging.Fields{"addr": redisConfig.Addr, "db": redisConfig.DB})
	s.logger.Info("基础设施层初始化完成")
	return nil
}

func (s *GatewayBootstrap) initializeApplicationLayer(cfg *config.Config) error {
	_ = cfg
	s.logger.Info("初始化应用服务层")
	s.commandBus = handlers.NewCommandBus()
	s.queryBus = handlers.NewQueryBus()

	// 创建仓储
	db := s.mongoClient.Database(cfg.Database.MongoDB.Database)
	characterRepo := persistence.NewCharacterRepository(db)
	itemRepo := persistence.NewItemRepository(db)
	questRepo := persistence.NewQuestRepository(db)

	// Instantiate application services
	s.mapService = appServices.NewMapService()
	s.fightService = appServices.NewFightService(nil)
	s.characterService = appServices.NewCharacterService(characterRepo, itemRepo, questRepo)
	s.updateMgr = appServices.NewUpdateManager(s.logger, 50*time.Millisecond)
	s.spawnMgr = appServices.NewSpawnManager(s.logger, 1024)
	// Wiring: map service uses spawn manager for async tasks
	s.mapService.SetSpawnManager(s.spawnMgr)
	s.logger.Info("应用服务层初始化完成")
	return nil
}

func (s *GatewayBootstrap) initializeTCPServer(cfg *config.Config) error {
	s.logger.Info("初始化TCP服务器")
	tcpCfg := &tcp.ServerConfig{Addr: fmt.Sprintf("%s:%d", cfg.Server.TCP.Host, cfg.Server.TCP.Port), MaxConnections: cfg.Server.TCP.MaxConnections, ReadTimeout: cfg.Server.TCP.ReadTimeout, WriteTimeout: cfg.Server.TCP.WriteTimeout, EnableCompression: cfg.Server.TCP.CompressionEnabled, BufferSize: cfg.Server.TCP.BufferSize}
	s.tcpServer = tcp.NewTCPServer(tcpCfg, s.commandBus, s.queryBus, s.logger)
	// Provide services to TCP server for handlers
	s.tcpServer.SetMapService(s.mapService)
	s.tcpServer.SetFightService(s.fightService)
	s.tcpServer.SetCharacterService(s.characterService)

	// Inject broadcaster from TCP server into MapService
	if s.mapService != nil {
		connMgr := s.tcpServer.GetConnectionManager()
		s.mapService.SetBroadcaster(func(recipients []character.EntityID, topic string, payload interface{}) {
			// Choose a message type based on topic
			var msgType uint32
			switch topic {
			case "entity_move":
				msgType = uint32(tcpProtocol.MsgPlayerMove)
			case "entity_appear", "entity_disappear":
				msgType = uint32(tcpProtocol.MsgPlayerStatusSync)
			case "skill_cast":
				msgType = uint32(tcpProtocol.MsgBattleSkill)
			default:
				msgType = uint32(tcpProtocol.MsgPlayerStatus)
			}

			msg := &tcpProtocol.Message{
				Header: tcpProtocol.MessageHeader{
					Magic:       tcpProtocol.MessageMagic,
					MessageID:   0,
					MessageType: msgType,
					Flags:       tcpProtocol.FlagBroadcast | tcpProtocol.FlagAsync,
					PlayerID:    0,
					Timestamp:   time.Now().Unix(),
					Sequence:    0,
				},
				Payload: map[string]interface{}{
					"topic":   topic,
					"payload": payload,
				},
			}
			if data, err := json.Marshal(msg); err == nil {
				for _, id := range recipients {
					if session, ok := connMgr.GetSessionByPlayer(int32(id)); ok {
						_ = session.Send(data)
					}
				}
			} else {
				s.logger.Error("广播消息序列化失败", err, logging.Fields{"topic": topic})
			}
		})
	}
	s.logger.Info("TCP服务器初始化完成")
	return nil
}

// Done returns a channel that's closed when the service context is canceled.
func (s *GatewayBootstrap) Done() <-chan struct{} { return s.ctx.Done() }
