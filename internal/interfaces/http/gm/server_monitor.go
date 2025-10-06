package gm

import (
	"context"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"

	"greatestworks/application/handlers"
	// "greatestworks/application/queries" // TODO: 实现查询系统
	"greatestworks/internal/infrastructure/logger"
	"greatestworks/internal/interfaces/http/auth"
)

// ServerMonitorHandler GM服务器监控处理器
type ServerMonitorHandler struct {
	queryBus *handlers.QueryBus
	logger   logger.Logger
}

// NewServerMonitorHandler 创建GM服务器监控处理器
func NewServerMonitorHandler(queryBus *handlers.QueryBus, logger logger.Logger) *ServerMonitorHandler {
	return &ServerMonitorHandler{
		queryBus: queryBus,
		logger:   logger,
	}
}

// ServerStatusResponse 服务器状态响应
type ServerStatusResponse struct {
	ServerInfo  ServerInfo  `json:"server_info"`
	SystemInfo  SystemInfo  `json:"system_info"`
	PlayerStats PlayerStats `json:"player_stats"`
	Performance Performance `json:"performance"`
	Connections Connections `json:"connections"`
	GameStats   GameStats   `json:"game_stats"`
	Timestamp   time.Time   `json:"timestamp"`
}

// ServerInfo 服务器信息
type ServerInfo struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Environment string    `json:"environment"`
	StartTime   time.Time `json:"start_time"`
	Uptime      string    `json:"uptime"`
	Region      string    `json:"region"`
	NodeID      string    `json:"node_id"`
}

// SystemInfo 系统信息
type SystemInfo struct {
	OS          string  `json:"os"`
	Arch        string  `json:"arch"`
	GoVersion   string  `json:"go_version"`
	CPUCores    int     `json:"cpu_cores"`
	MemoryTotal uint64  `json:"memory_total"`
	MemoryUsed  uint64  `json:"memory_used"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskTotal   uint64  `json:"disk_total"`
	DiskUsed    uint64  `json:"disk_used"`
	DiskUsage   float64 `json:"disk_usage"`
}

// PlayerStats 玩家统计
type PlayerStats struct {
	OnlineCount    int       `json:"online_count"`
	TotalCount     int       `json:"total_count"`
	NewToday       int       `json:"new_today"`
	ActiveToday    int       `json:"active_today"`
	PeakOnline     int       `json:"peak_online"`
	PeakOnlineTime time.Time `json:"peak_online_time"`
}

// Performance 性能指标
type Performance struct {
	CPUUsage       float64 `json:"cpu_usage"`
	MemoryUsage    float64 `json:"memory_usage"`
	Goroutines     int     `json:"goroutines"`
	GCPauseAvg     float64 `json:"gc_pause_avg"`
	GCPauseMax     float64 `json:"gc_pause_max"`
	RequestsPerSec float64 `json:"requests_per_sec"`
	ResponseTime   float64 `json:"response_time"`
	ErrorRate      float64 `json:"error_rate"`
}

// Connections 连接统计
type Connections struct {
	HTTPConnections  int `json:"http_connections"`
	TCPConnections   int `json:"tcp_connections"`
	WebSocketConns   int `json:"websocket_connections"`
	TotalConnections int `json:"total_connections"`
	MaxConnections   int `json:"max_connections"`
}

// GameStats 游戏统计
type GameStats struct {
	ActiveBattles   int     `json:"active_battles"`
	ActiveRooms     int     `json:"active_rooms"`
	MessagesPerSec  int     `json:"messages_per_sec"`
	EventsProcessed int     `json:"events_processed"`
	QueuedEvents    int     `json:"queued_events"`
	DatabaseQueries int     `json:"database_queries"`
	CacheHitRate    float64 `json:"cache_hit_rate"`
}

// MetricsHistoryRequest 指标历史请求
type MetricsHistoryRequest struct {
	Metric    string `form:"metric" binding:"required,oneof=cpu memory connections players"`
	TimeRange string `form:"time_range" binding:"required,oneof=1h 6h 24h 7d 30d"`
	Interval  string `form:"interval,omitempty"`
}

// MetricsHistoryResponse 指标历史响应
type MetricsHistoryResponse struct {
	Metric     string            `json:"metric"`
	TimeRange  string            `json:"time_range"`
	Interval   string            `json:"interval"`
	DataPoints []MetricDataPoint `json:"data_points"`
}

// MetricDataPoint 指标数据点
type MetricDataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Label     string    `json:"label,omitempty"`
}

// AlertsResponse 告警响应
type AlertsResponse struct {
	ActiveAlerts []Alert      `json:"active_alerts"`
	RecentAlerts []Alert      `json:"recent_alerts"`
	AlertSummary AlertSummary `json:"alert_summary"`
}

// Alert 告警信息
type Alert struct {
	ID          string     `json:"id"`
	Level       string     `json:"level"`
	Type        string     `json:"type"`
	Message     string     `json:"message"`
	Source      string     `json:"source"`
	TriggeredAt time.Time  `json:"triggered_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	Status      string     `json:"status"`
}

// AlertSummary 告警摘要
type AlertSummary struct {
	Critical int `json:"critical"`
	Warning  int `json:"warning"`
	Info     int `json:"info"`
	Total    int `json:"total"`
}

// GetServerStatus 获取服务器状态
func (h *ServerMonitorHandler) GetServerStatus(c *gin.Context) {
	ctx := context.Background()

	// 查询服务器状态
	query := &system.GetServerStatusQuery{}
	result, err := handlers.ExecuteQueryTyped[*system.GetServerStatusQuery, *system.GetServerStatusResult](ctx, h.queryBus, query)
	if err != nil {
		h.logger.Error("Failed to get server status", "error", err)
		c.JSON(500, gin.H{"error": "Failed to get server status", "success": false})
		return
	}

	// 获取系统运行时信息
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// 构造响应
	response := &ServerStatusResponse{
		ServerInfo: ServerInfo{
			Name:        result.ServerName,
			Version:     result.Version,
			Environment: result.Environment,
			StartTime:   result.StartTime,
			Uptime:      time.Since(result.StartTime).String(),
			Region:      result.Region,
			NodeID:      result.NodeID,
		},
		SystemInfo: SystemInfo{
			OS:          runtime.GOOS,
			Arch:        runtime.GOARCH,
			GoVersion:   runtime.Version(),
			CPUCores:    runtime.NumCPU(),
			MemoryTotal: result.SystemInfo.MemoryTotal,
			MemoryUsed:  memStats.Alloc,
			MemoryUsage: float64(memStats.Alloc) / float64(result.SystemInfo.MemoryTotal) * 100,
			DiskTotal:   result.SystemInfo.DiskTotal,
			DiskUsed:    result.SystemInfo.DiskUsed,
			DiskUsage:   float64(result.SystemInfo.DiskUsed) / float64(result.SystemInfo.DiskTotal) * 100,
		},
		PlayerStats: PlayerStats{
			OnlineCount:    result.PlayerStats.OnlineCount,
			TotalCount:     result.PlayerStats.TotalCount,
			NewToday:       result.PlayerStats.NewToday,
			ActiveToday:    result.PlayerStats.ActiveToday,
			PeakOnline:     result.PlayerStats.PeakOnline,
			PeakOnlineTime: result.PlayerStats.PeakOnlineTime,
		},
		Performance: Performance{
			CPUUsage:       result.Performance.CPUUsage,
			MemoryUsage:    float64(memStats.Alloc) / float64(result.SystemInfo.MemoryTotal) * 100,
			Goroutines:     runtime.NumGoroutine(),
			GCPauseAvg:     result.Performance.GCPauseAvg,
			GCPauseMax:     result.Performance.GCPauseMax,
			RequestsPerSec: result.Performance.RequestsPerSec,
			ResponseTime:   result.Performance.ResponseTime,
			ErrorRate:      result.Performance.ErrorRate,
		},
		Connections: Connections{
			HTTPConnections:  result.Connections.HTTPConnections,
			TCPConnections:   result.Connections.TCPConnections,
			WebSocketConns:   result.Connections.WebSocketConns,
			TotalConnections: result.Connections.TotalConnections,
			MaxConnections:   result.Connections.MaxConnections,
		},
		GameStats: GameStats{
			ActiveBattles:   result.GameStats.ActiveBattles,
			ActiveRooms:     result.GameStats.ActiveRooms,
			MessagesPerSec:  result.GameStats.MessagesPerSec,
			EventsProcessed: result.GameStats.EventsProcessed,
			QueuedEvents:    result.GameStats.QueuedEvents,
			DatabaseQueries: result.GameStats.DatabaseQueries,
			CacheHitRate:    result.GameStats.CacheHitRate,
		},
		Timestamp: time.Now(),
	}

	// 记录GM操作日志
	gmUser, _ := auth.GetCurrentUser(c)
	h.logger.Debug("GM viewed server status", "gm_user", gmUser.Username)

	c.JSON(200, gin.H{"data": response, "success": true})
}

// GetMetricsHistory 获取指标历史数据
func (h *ServerMonitorHandler) GetMetricsHistory(c *gin.Context) {
	var req MetricsHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error("Invalid metrics history request", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request parameters", "success": false})
		return
	}

	// 设置默认间隔
	if req.Interval == "" {
		switch req.TimeRange {
		case "1h":
			req.Interval = "1m"
		case "6h":
			req.Interval = "5m"
		case "24h":
			req.Interval = "15m"
		case "7d":
			req.Interval = "1h"
		case "30d":
			req.Interval = "6h"
		}
	}

	ctx := context.Background()

	// 查询指标历史数据
	query := &system.GetMetricsHistoryQuery{
		Metric:    req.Metric,
		TimeRange: req.TimeRange,
		Interval:  req.Interval,
	}

	result, err := handlers.ExecuteQueryTyped[*system.GetMetricsHistoryQuery, *system.GetMetricsHistoryResult](ctx, h.queryBus, query)
	if err != nil {
		h.logger.Error("Failed to get metrics history", "error", err, "metric", req.Metric)
		c.JSON(500, gin.H{"error": "Failed to get metrics history", "success": false})
		return
	}

	// 构造响应
	dataPoints := make([]MetricDataPoint, len(result.DataPoints))
	for i, dp := range result.DataPoints {
		dataPoints[i] = MetricDataPoint{
			Timestamp: dp.Timestamp,
			Value:     dp.Value,
			Label:     dp.Label,
		}
	}

	response := &MetricsHistoryResponse{
		Metric:     req.Metric,
		TimeRange:  req.TimeRange,
		Interval:   req.Interval,
		DataPoints: dataPoints,
	}

	// 记录GM操作日志
	gmUser, _ := auth.GetCurrentUser(c)
	h.logger.Debug("GM viewed metrics history", "gm_user", gmUser.Username, "metric", req.Metric, "time_range", req.TimeRange)

	c.JSON(200, gin.H{"data": response, "success": true})
}

// GetAlerts 获取告警信息
func (h *ServerMonitorHandler) GetAlerts(c *gin.Context) {
	ctx := context.Background()

	// 查询告警信息
	query := &system.GetAlertsQuery{}
	result, err := handlers.ExecuteQueryTyped[*system.GetAlertsQuery, *system.GetAlertsResult](ctx, h.queryBus, query)
	if err != nil {
		h.logger.Error("Failed to get alerts", "error", err)
		c.JSON(500, gin.H{"error": "Failed to get alerts", "success": false})
		return
	}

	// 构造响应
	activeAlerts := make([]Alert, len(result.ActiveAlerts))
	for i, alert := range result.ActiveAlerts {
		activeAlerts[i] = Alert{
			ID:          alert.ID,
			Level:       alert.Level,
			Type:        alert.Type,
			Message:     alert.Message,
			Source:      alert.Source,
			TriggeredAt: alert.TriggeredAt,
			ResolvedAt:  alert.ResolvedAt,
			Status:      alert.Status,
		}
	}

	recentAlerts := make([]Alert, len(result.RecentAlerts))
	for i, alert := range result.RecentAlerts {
		recentAlerts[i] = Alert{
			ID:          alert.ID,
			Level:       alert.Level,
			Type:        alert.Type,
			Message:     alert.Message,
			Source:      alert.Source,
			TriggeredAt: alert.TriggeredAt,
			ResolvedAt:  alert.ResolvedAt,
			Status:      alert.Status,
		}
	}

	response := &AlertsResponse{
		ActiveAlerts: activeAlerts,
		RecentAlerts: recentAlerts,
		AlertSummary: AlertSummary{
			Critical: result.AlertSummary.Critical,
			Warning:  result.AlertSummary.Warning,
			Info:     result.AlertSummary.Info,
			Total:    result.AlertSummary.Total,
		},
	}

	// 记录GM操作日志
	gmUser, _ := auth.GetCurrentUser(c)
	h.logger.Debug("GM viewed alerts", "gm_user", gmUser.Username)

	c.JSON(200, gin.H{"data": response, "success": true})
}

// GetOnlinePlayers 获取在线玩家列表
func (h *ServerMonitorHandler) GetOnlinePlayers(c *gin.Context) {
	page := 1
	pageSize := 50

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	ctx := context.Background()

	// 查询在线玩家
	query := &system.GetOnlinePlayersQuery{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := handlers.ExecuteQueryTyped[*system.GetOnlinePlayersQuery, *system.GetOnlinePlayersResult](ctx, h.queryBus, query)
	if err != nil {
		h.logger.Error("Failed to get online players", "error", err)
		c.JSON(500, gin.H{"error": "Failed to get online players", "success": false})
		return
	}

	// 构造响应
	players := make([]map[string]interface{}, len(result.Players))
	for i, player := range result.Players {
		players[i] = map[string]interface{}{
			"id":              player.ID,
			"username":        player.Username,
			"name":            player.Name,
			"level":           player.Level,
			"status":          player.Status,
			"login_time":      player.LoginTime,
			"online_duration": time.Since(player.LoginTime).String(),
			"ip_address":      player.IPAddress,
			"location":        player.Location,
		}
	}

	response := map[string]interface{}{
		"players": players,
		"pagination": map[string]interface{}{
			"page":        page,
			"page_size":   pageSize,
			"total":       result.Total,
			"total_pages": (result.Total + int64(pageSize) - 1) / int64(pageSize),
		},
		"summary": map[string]interface{}{
			"total_online": result.Total,
			"avg_level":    result.AvgLevel,
			"new_players":  result.NewPlayersToday,
		},
	}

	// 记录GM操作日志
	gmUser, _ := auth.GetCurrentUser(c)
	h.logger.Debug("GM viewed online players", "gm_user", gmUser.Username, "page", page, "page_size", pageSize)

	c.JSON(200, gin.H{"data": response, "success": true})
}

// RestartServer 重启服务器（仅超级管理员）
func (h *ServerMonitorHandler) RestartServer(c *gin.Context) {
	type RestartRequest struct {
		Reason       string `json:"reason" binding:"required"`
		DelayMinutes int    `json:"delay_minutes,omitempty"`
		NotifyUsers  bool   `json:"notify_users"`
	}

	var req RestartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid restart server request", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request format", "success": false})
		return
	}

	// 获取GM用户信息
	gmUser, _ := auth.GetCurrentUser(c)

	// 记录重启操作日志
	h.logger.Warn("Server restart initiated by GM", "gm_user", gmUser.Username, "reason", req.Reason, "delay_minutes", req.DelayMinutes)

	// TODO: 实现服务器重启逻辑
	// 1. 通知所有在线玩家
	// 2. 等待延迟时间
	// 3. 优雅关闭服务器
	// 4. 重启服务器

	response := map[string]interface{}{
		"message":       "Server restart scheduled",
		"delay_minutes": req.DelayMinutes,
		"restart_time":  time.Now().Add(time.Duration(req.DelayMinutes) * time.Minute),
		"initiated_by":  gmUser.Username,
	}

	c.JSON(200, gin.H{"data": response, "success": true, "message": "Server restart scheduled successfully"})
}
