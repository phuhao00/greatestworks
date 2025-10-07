// Package http 提供HTTP服务器实现
// 基于DDD架构的分布式认证服务HTTP接口
package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"greatestworks/internal/infrastructure/logging"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// HTTPServerConfig HTTP服务器配置
type HTTPServerConfig struct {
	Host           string          `yaml:"host"`
	Port           int             `yaml:"port"`
	ReadTimeout    time.Duration   `yaml:"read_timeout"`
	WriteTimeout   time.Duration   `yaml:"write_timeout"`
	IdleTimeout    time.Duration   `yaml:"idle_timeout"`
	MaxHeaderBytes int             `yaml:"max_header_bytes"`
	EnableCORS     bool            `yaml:"enable_cors"`
	EnableMetrics  bool            `yaml:"enable_metrics"`
	EnableSwagger  bool            `yaml:"enable_swagger"`
	RateLimit      RateLimitConfig `yaml:"rate_limit"`
	CORS           CORSConfig      `yaml:"cors"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	RequestsPerSecond int `yaml:"requests_per_second"`
	Burst             int `yaml:"burst"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
}

// HTTPServer HTTP服务器
type HTTPServer struct {
	config *HTTPServerConfig
	logger logging.Logger
	server *http.Server
	router *gin.Engine
	ctx    context.Context
	cancel context.CancelFunc
}

// NewHTTPServer 创建HTTP服务器
func NewHTTPServer(config *HTTPServerConfig, logger logging.Logger) *HTTPServer {
	ctx, cancel := context.WithCancel(context.Background())

	return &HTTPServer{
		config: config,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start 启动HTTP服务器
func (s *HTTPServer) Start() error {
	// 创建路由器
	s.router = s.createRouter()

	// 创建HTTP服务器
	s.server = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		Handler:        s.router,
		ReadTimeout:    s.config.ReadTimeout,
		WriteTimeout:   s.config.WriteTimeout,
		IdleTimeout:    s.config.IdleTimeout,
		MaxHeaderBytes: s.config.MaxHeaderBytes,
	}

	// 启动服务器
	go func() {
		s.logger.Info("HTTP服务器启动", "address", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("HTTP服务器运行失败", "error", err)
		}
	}()

	return nil
}

// Stop 停止HTTP服务器
func (s *HTTPServer) Stop() error {
	s.logger.Info("停止HTTP服务器")

	// 取消上下文
	s.cancel()

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("HTTP服务器关闭失败", "error", err)
		return err
	}

	s.logger.Info("HTTP服务器已停止")
	return nil
}

// createRouter 创建路由器
func (s *HTTPServer) createRouter() *gin.Engine {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由器
	router := gin.New()

	// 中间件
	s.setupMiddleware(router)

	// 路由
	s.setupRoutes(router)

	return router
}

// setupMiddleware 设置中间件
func (s *HTTPServer) setupMiddleware(router *gin.Engine) {
	// 日志中间件
	router.Use(gin.Logger())

	// 恢复中间件
	router.Use(gin.Recovery())

	// 请求ID中间件
	router.Use(requestid.New())

	// 超时中间件
	router.Use(timeout.New(
		timeout.WithTimeout(30*time.Second),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(func(c *gin.Context) {
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error": "Request timeout",
			})
		}),
	))

	// CORS中间件
	if s.config.EnableCORS {
		corsConfig := cors.Config{
			AllowOrigins:     s.config.CORS.AllowedOrigins,
			AllowMethods:     s.config.CORS.AllowedMethods,
			AllowHeaders:     s.config.CORS.AllowedHeaders,
			AllowCredentials: s.config.CORS.AllowCredentials,
		}
		router.Use(cors.New(corsConfig))
	}

	// 限流中间件
	if s.config.RateLimit.RequestsPerSecond > 0 {
		limiter := rate.NewLimiter(
			rate.Limit(s.config.RateLimit.RequestsPerSecond),
			s.config.RateLimit.Burst,
		)
		router.Use(func(c *gin.Context) {
			if !limiter.Allow() {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "Rate limit exceeded",
				})
				c.Abort()
				return
			}
			c.Next()
		})
	}
}

// setupRoutes 设置路由
func (s *HTTPServer) setupRoutes(router *gin.Engine) {
	// 健康检查
	router.GET("/health", s.healthCheck)

	// 指标端点
	if s.config.EnableMetrics {
		router.GET("/metrics", s.metrics)
	}

	// API路由组
	api := router.Group("/api/v1")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/login", s.login)
			auth.POST("/register", s.register)
			auth.POST("/logout", s.logout)
			auth.POST("/refresh", s.refresh)
			auth.GET("/profile", s.getProfile)
		}

		// 用户路由
		users := api.Group("/users")
		{
			users.GET("/:id", s.getUser)
			users.PUT("/:id", s.updateUser)
			users.DELETE("/:id", s.deleteUser)
		}
	}

	// Swagger文档
	if s.config.EnableSwagger {
		router.GET("/swagger/*any", s.swagger)
	}
}

// healthCheck 健康检查
func (s *HTTPServer) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "auth-service",
		"version":   "1.0.0",
	})
}

// metrics 指标端点
func (s *HTTPServer) metrics(c *gin.Context) {
	// TODO: 实现指标收集
	c.JSON(http.StatusOK, gin.H{
		"metrics": "not implemented",
	})
}

// login 登录
func (s *HTTPServer) login(c *gin.Context) {
	// TODO: 实现登录逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "login endpoint",
	})
}

// register 注册
func (s *HTTPServer) register(c *gin.Context) {
	// TODO: 实现注册逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "register endpoint",
	})
}

// logout 登出
func (s *HTTPServer) logout(c *gin.Context) {
	// TODO: 实现登出逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "logout endpoint",
	})
}

// refresh 刷新令牌
func (s *HTTPServer) refresh(c *gin.Context) {
	// TODO: 实现刷新令牌逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "refresh endpoint",
	})
}

// getProfile 获取用户资料
func (s *HTTPServer) getProfile(c *gin.Context) {
	// TODO: 实现获取用户资料逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "get profile endpoint",
	})
}

// getUser 获取用户
func (s *HTTPServer) getUser(c *gin.Context) {
	// TODO: 实现获取用户逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "get user endpoint",
	})
}

// updateUser 更新用户
func (s *HTTPServer) updateUser(c *gin.Context) {
	// TODO: 实现更新用户逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "update user endpoint",
	})
}

// deleteUser 删除用户
func (s *HTTPServer) deleteUser(c *gin.Context) {
	// TODO: 实现删除用户逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "delete user endpoint",
	})
}

// swagger Swagger文档
func (s *HTTPServer) swagger(c *gin.Context) {
	// TODO: 实现Swagger文档
	c.JSON(http.StatusOK, gin.H{
		"message": "swagger endpoint",
	})
}

// GetStats 获取服务器统计信息
func (s *HTTPServer) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})

	stats["status"] = "running"
	stats["address"] = s.server.Addr
	stats["read_timeout"] = s.config.ReadTimeout.String()
	stats["write_timeout"] = s.config.WriteTimeout.String()
	stats["idle_timeout"] = s.config.IdleTimeout.String()

	return stats
}

// DefaultHTTPServerConfig 默认HTTP服务器配置
func DefaultHTTPServerConfig() *HTTPServerConfig {
	return &HTTPServerConfig{
		Host:           "0.0.0.0",
		Port:           8080,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
		EnableCORS:     true,
		EnableMetrics:  true,
		EnableSwagger:  true,
		RateLimit: RateLimitConfig{
			RequestsPerSecond: 100,
			Burst:             200,
		},
		CORS: CORSConfig{
			AllowedOrigins:   []string{"http://localhost:3000", "https://greatestworks.com"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With"},
			AllowCredentials: true,
		},
	}
}
