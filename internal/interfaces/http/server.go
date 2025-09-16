package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	
	"greatestworks/application/handlers"
	"greatestworks/internal/infrastructure/logger"
)

// HTTPServer HTTP服务器
type HTTPServer struct {
	server     *http.Server
	router     *gin.Engine
	config     *HTTPServerConfig
	logger     logger.Logger
	services   *ServiceContainer
	handlers   *HandlerContainer
}

// HTTPServerConfig HTTP服务器配置
type HTTPServerConfig struct {
	Host               string        `json:"host" yaml:"host"`
	Port               int           `json:"port" yaml:"port"`
	ReadTimeout        time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout       time.Duration `json:"write_timeout" yaml:"write_timeout"`
	IdleTimeout        time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
	MaxHeaderBytes     int           `json:"max_header_bytes" yaml:"max_header_bytes"`
	EnableCORS         bool          `json:"enable_cors" yaml:"enable_cors"`
	EnableRequestID    bool          `json:"enable_request_id" yaml:"enable_request_id"`
	EnableLogging      bool          `json:"enable_logging" yaml:"enable_logging"`
	EnableRecovery     bool          `json:"enable_recovery" yaml:"enable_recovery"`
	EnableMetrics      bool          `json:"enable_metrics" yaml:"enable_metrics"`
	TrustedProxies     []string      `json:"trusted_proxies" yaml:"trusted_proxies"`
	AllowedOrigins     []string      `json:"allowed_origins" yaml:"allowed_origins"`
	AllowedMethods     []string      `json:"allowed_methods" yaml:"allowed_methods"`
	AllowedHeaders     []string      `json:"allowed_headers" yaml:"allowed_headers"`
	RateLimitEnabled   bool          `json:"rate_limit_enabled" yaml:"rate_limit_enabled"`
	RateLimitRequests  int           `json:"rate_limit_requests" yaml:"rate_limit_requests"`
	RateLimitDuration  time.Duration `json:"rate_limit_duration" yaml:"rate_limit_duration"`
}

// ServiceContainer 服务容器
type ServiceContainer struct {
	CommandBus *handlers.CommandBus
	QueryBus   *handlers.QueryBus
}

// HandlerContainer 处理器容器
type HandlerContainer struct {
	PlayerHandler   *PlayerHandler
	BattleHandler   *BattleHandler
	PetHandler      *PetHandler
	BuildingHandler *BuildingHandler
	HealthHandler   *HealthHandler
}

// NewHTTPServer 创建HTTP服务器
func NewHTTPServer(config *HTTPServerConfig, services *ServiceContainer, logger logger.Logger) (*HTTPServer, error) {
	if config == nil {
		config = &HTTPServerConfig{
			Host:               "0.0.0.0",
			Port:               8080,
			ReadTimeout:        30 * time.Second,
			WriteTimeout:       30 * time.Second,
			IdleTimeout:        60 * time.Second,
			MaxHeaderBytes:     1 << 20, // 1MB
			EnableCORS:         true,
			EnableRequestID:    true,
			EnableLogging:      true,
			EnableRecovery:     true,
			EnableMetrics:      true,
			AllowedOrigins:     []string{"*"},
			AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:     []string{"*"},
			RateLimitEnabled:   true,
			RateLimitRequests:  100,
			RateLimitDuration:  time.Minute,
		}
	}
	
	if services == nil {
		return nil, fmt.Errorf("services container is required")
	}
	
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)
	
	// 创建路由器
	router := gin.New()
	
	// 创建处理器
	handlers := &HandlerContainer{
		PlayerHandler:   NewPlayerHandler(services.CommandBus, services.QueryBus, logger),
		BattleHandler:   NewBattleHandler(services.CommandBus, services.QueryBus, logger),
		PetHandler:      NewPetHandler(services.CommandBus, services.QueryBus, logger),
		BuildingHandler: NewBuildingHandler(services.CommandBus, services.QueryBus, logger),
		HealthHandler:   NewHealthHandler(logger),
	}
	
	httpServer := &HTTPServer{
		config:   config,
		router:   router,
		logger:   logger,
		services: services,
		handlers: handlers,
	}
	
	// 设置中间件
	httpServer.setupMiddleware()
	
	// 设置路由
	httpServer.setupRoutes()
	
	// 创建HTTP服务器
	httpServer.server = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler:        router,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		IdleTimeout:    config.IdleTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}
	
	logger.Info("HTTP server created successfully", "address", httpServer.server.Addr)
	return httpServer, nil
}

// setupMiddleware 设置中间件
func (s *HTTPServer) setupMiddleware() {
	// 恢复中间件
	if s.config.EnableRecovery {
		s.router.Use(gin.Recovery())
	}
	
	// 请求ID中间件
	if s.config.EnableRequestID {
		s.router.Use(requestid.New())
	}
	
	// 日志中间件
	if s.config.EnableLogging {
		s.router.Use(s.loggingMiddleware())
	}
	
	// CORS中间件
	if s.config.EnableCORS {
		corsConfig := cors.Config{
			AllowOrigins:     s.config.AllowedOrigins,
			AllowMethods:     s.config.AllowedMethods,
			AllowHeaders:     s.config.AllowedHeaders,
			ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}
		s.router.Use(cors.New(corsConfig))
	}
	
	// 速率限制中间件
	if s.config.RateLimitEnabled {
		s.router.Use(s.rateLimitMiddleware())
	}
	
	// 指标中间件
	if s.config.EnableMetrics {
		s.router.Use(s.metricsMiddleware())
	}
	
	// 设置信任的代理
	if len(s.config.TrustedProxies) > 0 {
		s.router.SetTrustedProxies(s.config.TrustedProxies)
	}
}

// setupRoutes 设置路由
func (s *HTTPServer) setupRoutes() {
	// API版本组
	v1 := s.router.Group("/api/v1")
	{
		// 健康检查
		health := v1.Group("/health")
		{
			health.GET("/", s.handlers.HealthHandler.Check)
			health.GET("/ready", s.handlers.HealthHandler.Ready)
			health.GET("/live", s.handlers.HealthHandler.Live)
		}
		
		// 玩家相关路由
		players := v1.Group("/players")
		{
			players.POST("/", s.handlers.PlayerHandler.CreatePlayer)
			players.GET("/:id", s.handlers.PlayerHandler.GetPlayer)
			players.PUT("/:id", s.handlers.PlayerHandler.UpdatePlayer)
			players.DELETE("/:id", s.handlers.PlayerHandler.DeletePlayer)
			players.GET("/", s.handlers.PlayerHandler.ListPlayers)
			players.POST("/:id/move", s.handlers.PlayerHandler.MovePlayer)
			players.POST("/:id/level-up", s.handlers.PlayerHandler.LevelUpPlayer)
			players.GET("/:id/stats", s.handlers.PlayerHandler.GetPlayerStats)
		}
		
		// 战斗相关路由
		battles := v1.Group("/battles")
		{
			battles.POST("/", s.handlers.BattleHandler.CreateBattle)
			battles.GET("/:id", s.handlers.BattleHandler.GetBattle)
			battles.PUT("/:id", s.handlers.BattleHandler.UpdateBattle)
			battles.DELETE("/:id", s.handlers.BattleHandler.DeleteBattle)
			battles.GET("/", s.handlers.BattleHandler.ListBattles)
			battles.POST("/:id/join", s.handlers.BattleHandler.JoinBattle)
			battles.POST("/:id/leave", s.handlers.BattleHandler.LeaveBattle)
			battles.POST("/:id/start", s.handlers.BattleHandler.StartBattle)
			battles.POST("/:id/actions", s.handlers.BattleHandler.ExecuteAction)
			battles.GET("/:id/status", s.handlers.BattleHandler.GetBattleStatus)
		}
		
		// 宠物相关路由
		pets := v1.Group("/pets")
		{
			pets.POST("/", s.handlers.PetHandler.CreatePet)
			pets.GET("/:id", s.handlers.PetHandler.GetPet)
			pets.PUT("/:id", s.handlers.PetHandler.UpdatePet)
			pets.DELETE("/:id", s.handlers.PetHandler.DeletePet)
			pets.GET("/", s.handlers.PetHandler.ListPets)
			pets.POST("/:id/feed", s.handlers.PetHandler.FeedPet)
			pets.POST("/:id/train", s.handlers.PetHandler.TrainPet)
			pets.POST("/:id/upgrade", s.handlers.PetHandler.UpgradePet)
			pets.POST("/:id/revive", s.handlers.PetHandler.RevivePet)
			pets.GET("/player/:player_id", s.handlers.PetHandler.GetPlayerPets)
		}
		
		// 建筑相关路由
		buildings := v1.Group("/buildings")
		{
			buildings.POST("/", s.handlers.BuildingHandler.CreateBuilding)
			buildings.GET("/:id", s.handlers.BuildingHandler.GetBuilding)
			buildings.PUT("/:id", s.handlers.BuildingHandler.UpdateBuilding)
			buildings.DELETE("/:id", s.handlers.BuildingHandler.DeleteBuilding)
			buildings.GET("/", s.handlers.BuildingHandler.ListBuildings)
			buildings.POST("/:id/construct", s.handlers.BuildingHandler.StartConstruction)
			buildings.POST("/:id/upgrade", s.handlers.BuildingHandler.StartUpgrade)
			buildings.POST("/:id/repair", s.handlers.BuildingHandler.RepairBuilding)
			buildings.POST("/:id/demolish", s.handlers.BuildingHandler.DemolishBuilding)
			buildings.GET("/player/:player_id", s.handlers.BuildingHandler.GetPlayerBuildings)
		}
	}
	
	// 根路径重定向到API文档
	s.router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Greatest Works API Server",
			"version": "v1.0.0",
			"docs":    "/api/v1/docs",
		})
	})
	
	// 404处理
	s.router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "The requested resource was not found",
			"path":    c.Request.URL.Path,
		})
	})
}

// Start 启动HTTP服务器
func (s *HTTPServer) Start(ctx context.Context) error {
	s.logger.Info("Starting HTTP server", "address", s.server.Addr)
	
	// 启动服务器
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Failed to start HTTP server", "error", err)
		}
	}()
	
	// 等待上下文取消
	<-ctx.Done()
	
	// 优雅关闭
	return s.Stop()
}

// Stop 停止HTTP服务器
func (s *HTTPServer) Stop() error {
	s.logger.Info("Stopping HTTP server")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("Failed to stop HTTP server gracefully", "error", err)
		return fmt.Errorf("failed to stop HTTP server: %w", err)
	}
	
	s.logger.Info("HTTP server stopped successfully")
	return nil
}

// GetRouter 获取路由器
func (s *HTTPServer) GetRouter() *gin.Engine {
	return s.router
}

// GetServer 获取HTTP服务器
func (s *HTTPServer) GetServer() *http.Server {
	return s.server
}

// 中间件实现

// loggingMiddleware 日志中间件
func (s *HTTPServer) loggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		s.logger.Info("HTTP Request",
			"method", param.Method,
			"path", param.Path,
			"status", param.StatusCode,
			"latency", param.Latency,
			"client_ip", param.ClientIP,
			"user_agent", param.Request.UserAgent(),
			"request_id", param.Request.Header.Get("X-Request-ID"),
		)
		return ""
	})
}

// rateLimitMiddleware 速率限制中间件
func (s *HTTPServer) rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现速率限制逻辑
		c.Next()
	}
}

// metricsMiddleware 指标中间件
func (s *HTTPServer) metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		
		// TODO: 记录指标
		s.logger.Debug("Request metrics",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", duration,
		)
	}
}