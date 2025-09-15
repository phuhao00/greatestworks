// Package main provides bootstrap functionality for the game server
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"greatestworks/application/services"
	"greatestworks/internal/infrastructure/config"
	"greatestworks/internal/infrastructure/logging"
	"greatestworks/internal/infrastructure/monitoring"
	"greatestworks/internal/infrastructure/network"
	"greatestworks/internal/infrastructure/persistence"
	"greatestworks/internal/interfaces/tcp"
)

const (
	// ToWeaveletKey is the environment variable under which the file descriptor
	// for messages sent from envelope to weavelet is stored. For internal use by
	// Service Weaver infrastructure.
	ToWeaveletKey = "ENVELOPE_TO_WEAVELET_FD"

	// ToEnvelopeKey is the environment variable under which the file descriptor
	// for messages sent from weavelet to envelope is stored. For internal use by
	// Service Weaver infrastructure.
	ToEnvelopeKey = "WEAVELET_TO_ENVELOPE_FD"
)

// Bootstrap holds configuration information used to start a process execution.
type Bootstrap struct {
	ToWeaveletFd int    // File descriptor on which to send to weavelet (0 if unset)
	ToEnvelopeFd int    // File descriptor from which to send to envelope (0 if unset)
	TestConfig   string // Configuration passed by user test code to weavertest
	
	// Game server specific fields
	Config       *config.Config
	Logger       logging.Logger
	Metrics      monitoring.MetricsRegistry
	ShutdownChan chan os.Signal
}

// BootstrapKey is the Context key used by weavertest to pass Bootstrap to [weaver.Init].
type BootstrapKey struct{}

// ServerBootstrap manages the complete server bootstrap process
type ServerBootstrap struct {
	bootstrap Bootstrap
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewServerBootstrap creates a new server bootstrap instance
func NewServerBootstrap() *ServerBootstrap {
	ctx, cancel := context.WithCancel(context.Background())
	return &ServerBootstrap{
		ctx:    ctx,
		cancel: cancel,
	}
}

// Initialize performs the complete server initialization
func (sb *ServerBootstrap) Initialize() error {
	log.Println("初始化游戏服务器...")
	
	// Get bootstrap configuration
	bootstrap, err := GetBootstrap(sb.ctx)
	if err != nil {
		return fmt.Errorf("获取启动配置失败: %w", err)
	}
	sb.bootstrap = bootstrap
	
	// Load configuration
	if err := sb.loadConfiguration(); err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}
	
	// Initialize logging
	if err := sb.initializeLogging(); err != nil {
		return fmt.Errorf("初始化日志系统失败: %w", err)
	}
	
	// Initialize metrics
	if err := sb.initializeMetrics(); err != nil {
		return fmt.Errorf("初始化监控系统失败: %w", err)
	}
	
	log.Println("服务器初始化完成")
	return nil
}

// loadConfiguration loads the server configuration
func (sb *ServerBootstrap) loadConfiguration() error {
	loader := config.NewConfigLoader()
	
	// Load configuration based on environment
	envManager := config.GetEnvManager()
	configPath := envManager.GetConfigPath()
	
	cfg, err := loader.LoadFromFile(configPath)
	if err != nil {
		return fmt.Errorf("加载配置文件失败: %w", err)
	}
	
	sb.bootstrap.Config = cfg
	return nil
}

// initializeLogging initializes the logging system
func (sb *ServerBootstrap) initializeLogging() error {
	logConfig := sb.bootstrap.Config.Logging
	
	// Create logger based on configuration
	var logger logging.Logger
	var err error
	
	switch logConfig.Type {
	case "zap":
		logger, err = logging.NewZapLogger(&logConfig)
	case "file":
		logger, err = logging.NewFileLogger(&logConfig)
	case "console":
		logger, err = logging.NewConsoleLogger(&logConfig)
	default:
		logger, err = logging.NewZapLogger(&logConfig)
	}
	
	if err != nil {
		return err
	}
	
	sb.bootstrap.Logger = logger
	return nil
}

// initializeMetrics initializes the monitoring system
func (sb *ServerBootstrap) initializeMetrics() error {
	metricsConfig := sb.bootstrap.Config.Monitoring
	
	// Create metrics registry
	registry, err := monitoring.NewPrometheusRegistry(&metricsConfig)
	if err != nil {
		return err
	}
	
	sb.bootstrap.Metrics = registry
	return nil
}

// StartServer starts the game server with all components
func (sb *ServerBootstrap) StartServer() error {
	log.Println("启动游戏服务器组件...")
	
	// Initialize database
	mongoDB, err := sb.initializeDatabase()
	if err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}
	defer sb.closeDatabase(mongoDB)
	
	// Initialize services
	playerService, err := sb.initializeServices(mongoDB)
	if err != nil {
		return fmt.Errorf("初始化服务失败: %w", err)
	}
	
	// Initialize and start TCP server
	server, err := sb.initializeTCPServer(playerService)
	if err != nil {
		return fmt.Errorf("初始化TCP服务器失败: %w", err)
	}
	
	// Start server
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("启动TCP服务器失败: %v", err)
		}
	}()
	
	log.Printf("游戏服务器启动成功，监听端口: %d", sb.bootstrap.Config.Server.Port)
	
	// Wait for shutdown signal
	sb.waitForShutdown()
	
	// Graceful shutdown
	return sb.gracefulShutdown(server)
}

// initializeDatabase initializes the database connection
func (sb *ServerBootstrap) initializeDatabase() (*persistence.MongoDB, error) {
	mongoConfig := sb.bootstrap.Config.Database.MongoDB
	mongoDB, err := persistence.NewMongoDB(&mongoConfig)
	if err != nil {
		return nil, err
	}
	
	log.Println("MongoDB连接成功")
	return mongoDB, nil
}

// closeDatabase closes the database connection
func (sb *ServerBootstrap) closeDatabase(mongoDB *persistence.MongoDB) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mongoDB.Close(ctx)
}

// initializeServices initializes application services
func (sb *ServerBootstrap) initializeServices(mongoDB *persistence.MongoDB) (*services.PlayerService, error) {
	// Initialize repositories
	playerRepo := persistence.NewPlayerRepository(mongoDB)
	
	// Initialize services
	playerService := services.NewPlayerService(playerRepo)
	
	return playerService, nil
}

// initializeTCPServer initializes the TCP server
func (sb *ServerBootstrap) initializeTCPServer(playerService *services.PlayerService) (*network.Server, error) {
	// Initialize TCP handlers
	playerHandler := tcp.NewPlayerHandler(playerService)
	
	// Create TCP server
	serverConfig := sb.bootstrap.Config.Server
	server := network.NewServer(&serverConfig)
	
	// Register handlers
	playerHandler.RegisterHandlers(server)
	
	// Start heartbeat checker
	server.StartHeartbeatChecker()
	
	return server, nil
}

// waitForShutdown waits for shutdown signal
func (sb *ServerBootstrap) waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	sb.bootstrap.ShutdownChan = sigChan
	
	<-sigChan
	log.Println("收到关闭信号，正在关闭服务器...")
}

// gracefulShutdown performs graceful shutdown
func (sb *ServerBootstrap) gracefulShutdown(server *network.Server) error {
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	
	if err := server.Stop(); err != nil {
		log.Printf("关闭TCP服务器失败: %v", err)
		return err
	}
	
	// Close logger
	if sb.bootstrap.Logger != nil {
		sb.bootstrap.Logger.Close()
	}
	
	select {
	case <-shutdownCtx.Done():
		log.Println("关闭超时")
		return fmt.Errorf("shutdown timeout")
	default:
		log.Println("服务器已优雅关闭")
		return nil
	}
}

// GetBootstrap returns information needed to configure process
// execution. For normal execution, this comes from the environment. For
// weavertest, it comes from a context value.
func GetBootstrap(ctx context.Context) (Bootstrap, error) {
	if val := ctx.Value(BootstrapKey{}); val != nil {
		bootstrap, ok := val.(Bootstrap)
		if !ok {
			return Bootstrap{}, fmt.Errorf("invalid type %T for bootstrap info in context", val)
		}
		return bootstrap, nil
	}

	str1 := os.Getenv(ToWeaveletKey)
	str2 := os.Getenv(ToEnvelopeKey)
	if str1 == "" && str2 == "" {
		return Bootstrap{}, nil
	}
	if str1 == "" || str2 == "" {
		return Bootstrap{}, fmt.Errorf("envelope/weavelet pipe should have 2 file descriptors, got (%s, %s)", str1, str2)
	}
	toWeaveletFd, err := strconv.Atoi(str1)
	if err != nil {
		return Bootstrap{}, fmt.Errorf("unable to parse envelope to weavelet fd: %w", err)
	}
	toEnvelopeFd, err := strconv.Atoi(str2)
	if err != nil {
		return Bootstrap{}, fmt.Errorf("unable to parse weavelet to envelope fd: %w", err)
	}
	return Bootstrap{
		ToWeaveletFd: toWeaveletFd,
		ToEnvelopeFd: toEnvelopeFd,
	}, nil
}

// HasPipes returns true if pipe information has been supplied. This
// is true except in the case of singleprocess.
func (b Bootstrap) HasPipes() bool {
	return b.ToWeaveletFd != 0 && b.ToEnvelopeFd != 0
}

// MakePipes creates pipe reader and writer. It returns an error if pipes are not configured.
func (b Bootstrap) MakePipes() (io.ReadCloser, io.WriteCloser, error) {
	toWeavelet, err := openFileDescriptor(b.ToWeaveletFd)
	if err != nil {
		return nil, nil, fmt.Errorf("open pipe to weavelet: %w", err)
	}
	toEnvelope, err := openFileDescriptor(b.ToEnvelopeFd)
	if err != nil {
		return nil, nil, fmt.Errorf("open pipe to envelope: %w", err)
	}
	return toWeavelet, toEnvelope, nil
}

func openFileDescriptor(fd int) (*os.File, error) {
	if fd == 0 {
		return nil, fmt.Errorf("bad file descriptor %d", fd)
	}
	f := os.NewFile(uintptr(fd), fmt.Sprint("/proc/self/fd/", fd))
	if f == nil {
		return nil, fmt.Errorf("open file descriptor %d: failed", fd)
	}
	return f, nil
}