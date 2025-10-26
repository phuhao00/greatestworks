package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"greatestworks/internal/application/services"
	"greatestworks/internal/infrastructure/logging"
)

// ReplicationHTTPHandlers 提供副本相关的HTTP处理器
type ReplicationHTTPHandlers struct {
	app    *services.ReplicationService
	logger logging.Logger
}

func NewReplicationHTTPHandlers(app *services.ReplicationService, logger logging.Logger) *ReplicationHTTPHandlers {
	return &ReplicationHTTPHandlers{app: app, logger: logger}
}

// RegisterReplicationRoutes 在给定服务器上注册副本相关路由
func RegisterReplicationRoutes(s *Server, h *ReplicationHTTPHandlers) {
	s.Handle("POST", "/instances/create", h.CreateInstance)
	s.Handle("POST", "/instances/join", h.JoinInstance)
	s.Handle("POST", "/instances/leave", h.LeaveInstance)
	s.Handle("GET", "/instances/info", h.GetInstanceInfo)
	s.Handle("GET", "/instances/active", h.ListActiveInstances)
	s.Handle("POST", "/instances/cleanup", h.CleanupExpiredInstances)
}

// CreateInstance 创建实例
func (h *ReplicationHTTPHandlers) CreateInstance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TemplateID    string `json:"template_id"`
		InstanceType  int    `json:"instance_type"`
		OwnerPlayerID string `json:"owner_player_id"`
		OwnerName     string `json:"owner_name"`
		OwnerLevel    int    `json:"owner_level"`
		MaxPlayers    int    `json:"max_players"`
		Difficulty    int    `json:"difficulty"`
		LifetimeSec   int64  `json:"lifetime_sec"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	dto, err := h.app.CreateInstance(r.Context(), &services.CreateInstanceCommand{
		TemplateID:    req.TemplateID,
		InstanceType:  req.InstanceType,
		OwnerPlayerID: req.OwnerPlayerID,
		OwnerName:     req.OwnerName,
		OwnerLevel:    req.OwnerLevel,
		MaxPlayers:    req.MaxPlayers,
		Difficulty:    req.Difficulty,
		Lifetime:      time.Duration(req.LifetimeSec) * time.Second,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, dto)
}

// JoinInstance 加入实例
func (h *ReplicationHTTPHandlers) JoinInstance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		InstanceID string `json:"instance_id"`
		PlayerID   string `json:"player_id"`
		PlayerName string `json:"player_name"`
		Level      int    `json:"level"`
		Role       string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.app.JoinInstance(r.Context(), &services.JoinInstanceCommand{
		InstanceID: req.InstanceID,
		PlayerID:   req.PlayerID,
		PlayerName: req.PlayerName,
		Level:      req.Level,
		Role:       req.Role,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"})
}

// LeaveInstance 离开实例
func (h *ReplicationHTTPHandlers) LeaveInstance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		InstanceID string `json:"instance_id"`
		PlayerID   string `json:"player_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.app.LeaveInstance(r.Context(), &services.LeaveInstanceCommand{
		InstanceID: req.InstanceID,
		PlayerID:   req.PlayerID,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"})
}

// GetInstanceInfo 获取实例信息（通过查询参数传入instance_id）
func (h *ReplicationHTTPHandlers) GetInstanceInfo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("instance_id")
	if id == "" {
		http.Error(w, "missing instance_id", http.StatusBadRequest)
		return
	}
	dto, err := h.app.GetInstanceInfo(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, dto)
}

// ListActiveInstances 列出活跃实例
func (h *ReplicationHTTPHandlers) ListActiveInstances(w http.ResponseWriter, r *http.Request) {
	dtos, err := h.app.ListActiveInstances(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, dtos)
}

// CleanupExpiredInstances 清理过期实例
func (h *ReplicationHTTPHandlers) CleanupExpiredInstances(w http.ResponseWriter, r *http.Request) {
	// 可选：从查询参数读取limit等
	_ = r.URL.Query().Get("limit")
	_, _ = strconv.Atoi(r.URL.Query().Get("limit"))

	count, err := h.app.CleanupExpiredInstances(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]int{"count": count})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	_ = enc.Encode(v)
}
