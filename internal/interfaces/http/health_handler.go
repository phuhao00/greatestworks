package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"greatestworks/internal/infrastructure/logger"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	logger logger.Logger
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(logger logger.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

// Check 健康检查
func (h *HealthHandler) Check(c *gin.Context) {
	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "greatestworks",
		"version":   "1.0.0",
	}

	c.JSON(http.StatusOK, response)
}

// Ready 就绪检查
func (h *HealthHandler) Ready(c *gin.Context) {
	// TODO: 添加依赖服务检查逻辑
	// 例如：数据库连接、缓存连接等

	response := map[string]interface{}{
		"status":    "ready",
		"timestamp": time.Now().Unix(),
		"checks": map[string]string{
			"database": "ok",
			"cache":    "ok",
		},
	}

	c.JSON(http.StatusOK, response)
}

// Live 存活检查
func (h *HealthHandler) Live(c *gin.Context) {
	response := map[string]interface{}{
		"status":    "alive",
		"timestamp": time.Now().Unix(),
		"uptime":    time.Since(time.Now()).Seconds(), // TODO: 实际启动时间
	}

	c.JSON(http.StatusOK, response)
}