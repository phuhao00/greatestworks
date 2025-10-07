package http

import (
	"context"
	"net/http"
	"time"

	"greatestworks/application/handlers"
	"greatestworks/internal/infrastructure/logging"

	"github.com/gin-gonic/gin"
)

// Server HTTP服务器
type Server struct {
	router          *gin.Engine
	server          *http.Server
	commandBus      *handlers.CommandBus
	queryBus        *handlers.QueryBus
	logger          logging.Logger
	playerHandler   *PlayerHandler
	battleHandler   *BattleHandler
	petHandler      *PetHandler
	buildingHandler *BuildingHandler
}

// NewServer 创建HTTP服务器
func NewServer(commandBus *handlers.CommandBus, queryBus *handlers.QueryBus, logger logging.Logger) *Server {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// 添加中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	server := &Server{
		router:     router,
		commandBus: commandBus,
		queryBus:   queryBus,
		logger:     logger,
	}

	// 初始化处理器
	server.playerHandler = NewPlayerHandler(commandBus, queryBus, logger)
	server.battleHandler = NewBattleHandler(commandBus, queryBus, logger)
	server.petHandler = NewPetHandler(commandBus, queryBus, logger)
	server.buildingHandler = NewBuildingHandler(commandBus, queryBus, logger)

	// 注册路由
	server.registerRoutes()

	return server
}

// registerRoutes 注册路由
func (s *Server) registerRoutes() {
	// 健康检查
	s.router.GET("/health", s.healthCheck)

	// API路由
	api := s.router.Group("/api/v1")
	{
		s.playerHandler.RegisterRoutes(api)
		s.battleHandler.RegisterRoutes(api)
		s.petHandler.RegisterRoutes(api)
		s.buildingHandler.RegisterRoutes(api)
	}
}

// healthCheck 健康检查
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"message": "服务器运行正常",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// Start 启动服务器
func (s *Server) Start(addr string) error {
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	s.logger.Info("HTTP服务器启动", map[string]interface{}{
		"address": addr,
	})

	return s.server.ListenAndServe()
}

// Stop 停止服务器
func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	s.logger.Info("HTTP服务器停止")

	return s.server.Shutdown(ctx)
}

// GetRouter 获取路由器
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
