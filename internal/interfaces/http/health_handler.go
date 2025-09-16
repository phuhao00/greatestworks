package http

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
	
	"github.com/gin-gonic/gin"
	
	"greatestworks/internal/infrastructure/logger"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	logger    logger.Logger
	startTime time.Time
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(logger logger.Logger) *HealthHandler {
	return &HealthHandler{
		logger:    logger,
		startTime: time.Now(),
	}
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Uptime    string            `json:"uptime"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services,omitempty"`
	System    *SystemInfo       `json:"system,omitempty"`
}

// SystemInfo 系统信息
type SystemInfo struct {
	GoVersion    string `json:"go_version"`
	Goroutines   int    `json:"goroutines"`
	MemoryUsage  string `json:"memory_usage"`
	CPUCount     int    `json:"cpu_count"`
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
}

// Check 基础健康检查
func (h *HealthHandler) Check(c *gin.Context) {
	uptime := time.Since(h.startTime)
	
	response := &HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    uptime.String(),
		Version:   "1.0.0",
	}
	
	c.JSON(http.StatusOK, response)
}

// Ready 就绪检查
func (h *HealthHandler) Ready(c *gin.Context) {
	uptime := time.Since(h.startTime)
	
	// 检查各个服务的就绪状态
	services := map[string]string{
		"database": h.checkDatabase(),
		"cache":    h.checkCache(),
		"queue":    h.checkQueue(),
	}
	
	// 判断整体状态
	status := "ready"
	statusCode := http.StatusOK
	for _, serviceStatus := range services {
		if serviceStatus != "healthy" {
			status = "not_ready"
			statusCode = http.StatusServiceUnavailable
			break
		}
	}
	
	response := &HealthResponse{
		Status:    status,
		Timestamp: time.Now(),
		Uptime:    uptime.String(),
		Version:   "1.0.0",
		Services:  services,
	}
	
	c.JSON(statusCode, response)
}

// Live 存活检查
func (h *HealthHandler) Live(c *gin.Context) {
	uptime := time.Since(h.startTime)
	
	// 获取系统信息
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	systemInfo := &SystemInfo{
		GoVersion:    runtime.Version(),
		Goroutines:   runtime.NumGoroutine(),
		MemoryUsage:  formatBytes(memStats.Alloc),
		CPUCount:     runtime.NumCPU(),
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
	}
	
	response := &HealthResponse{
		Status:    "alive",
		Timestamp: time.Now(),
		Uptime:    uptime.String(),
		Version:   "1.0.0",
		System:    systemInfo,
	}
	
	c.JSON(http.StatusOK, response)
}

// 辅助方法

// checkDatabase 检查数据库连接
func (h *HealthHandler) checkDatabase() string {
	// TODO: 实现数据库连接检查
	// 这里应该尝试连接数据库并执行简单查询
	return "healthy"
}

// checkCache 检查缓存连接
func (h *HealthHandler) checkCache() string {
	// TODO: 实现缓存连接检查
	// 这里应该尝试连接Redis并执行简单操作
	return "healthy"
}

// checkQueue 检查消息队列连接
func (h *HealthHandler) checkQueue() string {
	// TODO: 实现消息队列连接检查
	// 这里应该检查消息队列的连接状态
	return "healthy"
}

// formatBytes 格式化字节数
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}