package gm

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"greatestworks/application/commands/player"
	"greatestworks/application/queries/player"
	"greatestworks/application/handlers"
	"greatestworks/internal/infrastructure/logger"
	"greatestworks/internal/interfaces/http/auth"
)

// PlayerManagementHandler GM玩家管理处理器
type PlayerManagementHandler struct {
	commandBus *handlers.CommandBus
	queryBus   *handlers.QueryBus
	logger     logger.Logger
}

// NewPlayerManagementHandler 创建GM玩家管理处理器
func NewPlayerManagementHandler(commandBus *handlers.CommandBus, queryBus *handlers.QueryBus, logger logger.Logger) *PlayerManagementHandler {
	return &PlayerManagementHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
		logger:     logger,
	}
}

// PlayerSearchRequest 玩家搜索请求
type PlayerSearchRequest struct {
	Keyword    string `form:"keyword,omitempty"`
	PlayerID   string `form:"player_id,omitempty"`
	Username   string `form:"username,omitempty"`
	Email      string `form:"email,omitempty"`
	Status     string `form:"status,omitempty"`
	MinLevel   int    `form:"min_level,omitempty"`
	MaxLevel   int    `form:"max_level,omitempty"`
	Page       int    `form:"page,omitempty" binding:"min=1"`
	PageSize   int    `form:"page_size,omitempty" binding:"min=1,max=100"`
	SortBy     string `form:"sort_by,omitempty"`
	SortOrder  string `form:"sort_order,omitempty"`
}

// PlayerUpdateRequest GM玩家更新请求
type PlayerUpdateRequest struct {
	Name     *string `json:"name,omitempty"`
	Level    *int    `json:"level,omitempty"`
	Exp      *int64  `json:"exp,omitempty"`
	Status   *string `json:"status,omitempty"`
	HP       *int    `json:"hp,omitempty"`
	MaxHP    *int    `json:"max_hp,omitempty"`
	MP       *int    `json:"mp,omitempty"`
	MaxMP    *int    `json:"max_mp,omitempty"`
	Attack   *int    `json:"attack,omitempty"`
	Defense  *int    `json:"defense,omitempty"`
	Speed    *int    `json:"speed,omitempty"`
	Reason   string  `json:"reason" binding:"required"`
}

// PlayerBanRequest 玩家封禁请求
type PlayerBanRequest struct {
	PlayerID  string    `json:"player_id" binding:"required"`
	Reason    string    `json:"reason" binding:"required"`
	Duration  int       `json:"duration"` // 封禁时长（小时），0表示永久
	BanType   string    `json:"ban_type" binding:"required,oneof=login chat trade all"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// PlayerUnbanRequest 玩家解封请求
type PlayerUnbanRequest struct {
	PlayerID string `json:"player_id" binding:"required"`
	Reason   string `json:"reason" binding:"required"`
}

// GMPlayerResponse GM玩家响应
type GMPlayerResponse struct {
	ID            string            `json:"id"`
	Username      string            `json:"username"`
	Email         string            `json:"email"`
	Name          string            `json:"name"`
	Level         int               `json:"level"`
	Exp           int64             `json:"exp"`
	Status        string            `json:"status"`
	Position      PositionResponse  `json:"position"`
	Stats         StatsResponse     `json:"stats"`
	Avatar        string            `json:"avatar,omitempty"`
	Gender        int               `json:"gender,omitempty"`
	LastLoginAt   *time.Time        `json:"last_login_at,omitempty"`
	LastLogoutAt  *time.Time        `json:"last_logout_at,omitempty"`
	OnlineTime    int64             `json:"online_time"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	BanInfo       *BanInfo          `json:"ban_info,omitempty"`
}

// PositionResponse 位置响应
type PositionResponse struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// StatsResponse 属性响应
type StatsResponse struct {
	HP      int `json:"hp"`
	MaxHP   int `json:"max_hp"`
	MP      int `json:"mp"`
	MaxMP   int `json:"max_mp"`
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	Speed   int `json:"speed"`
}

// BanInfo 封禁信息
type BanInfo struct {
	IsBanned  bool       `json:"is_banned"`
	BanType   string     `json:"ban_type,omitempty"`
	Reason    string     `json:"reason,omitempty"`
	BannedBy  string     `json:"banned_by,omitempty"`
	BannedAt  *time.Time `json:"banned_at,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// SearchPlayers 搜索玩家
func (h *PlayerManagementHandler) SearchPlayers(c *gin.Context) {
	var req PlayerSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error("Invalid search players request", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request parameters", "success": false})
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	ctx := context.Background()

	// 执行搜索查询
	query := &player.SearchPlayersQuery{
		Keyword:   req.Keyword,
		PlayerID:  req.PlayerID,
		Username:  req.Username,
		Email:     req.Email,
		Status:    req.Status,
		MinLevel:  req.MinLevel,
		MaxLevel:  req.MaxLevel,
		Page:      req.Page,
		PageSize:  req.PageSize,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	result, err := handlers.ExecuteQueryTyped[*player.SearchPlayersQuery, *player.SearchPlayersResult](ctx, h.queryBus, query)
	if err != nil {
		h.logger.Error("Failed to search players", "error", err)
		c.JSON(500, gin.H{"error": "Failed to search players", "success": false})
		return
	}

	// 构造响应
	players := make([]*GMPlayerResponse, len(result.Players))
	for i, p := range result.Players {
		players[i] = &GMPlayerResponse{
			ID:       p.ID,
			Username: p.Username,
			Email:    p.Email,
			Name:     p.Name,
			Level:    p.Level,
			Exp:      p.Exp,
			Status:   p.Status,
			Position: PositionResponse{
				X: p.Position.X,
				Y: p.Position.Y,
				Z: p.Position.Z,
			},
			Stats: StatsResponse{
				HP:      p.Stats.HP,
				MaxHP:   p.Stats.MaxHP,
				MP:      p.Stats.MP,
				MaxMP:   p.Stats.MaxMP,
				Attack:  p.Stats.Attack,
				Defense: p.Stats.Defense,
				Speed:   p.Stats.Speed,
			},
			Avatar:       p.Avatar,
			Gender:       p.Gender,
			LastLoginAt:  p.LastLoginAt,
			LastLogoutAt: p.LastLogoutAt,
			OnlineTime:   p.OnlineTime,
			CreatedAt:    p.CreatedAt,
			UpdatedAt:    p.UpdatedAt,
		}

		// 添加封禁信息
		if p.BanInfo != nil {
			players[i].BanInfo = &BanInfo{
				IsBanned:  p.BanInfo.IsBanned,
				BanType:   p.BanInfo.BanType,
				Reason:    p.BanInfo.Reason,
				BannedBy:  p.BanInfo.BannedBy,
				BannedAt:  p.BanInfo.BannedAt,
				ExpiresAt: p.BanInfo.ExpiresAt,
			}
		}
	}

	response := map[string]interface{}{
		"players": players,
		"pagination": map[string]interface{}{
			"page":       result.Page,
			"page_size":  result.PageSize,
			"total":      result.Total,
			"total_pages": (result.Total + int64(result.PageSize) - 1) / int64(result.PageSize),
		},
	}

	// 记录GM操作日志
	gmUser, _ := auth.GetCurrentUser(c)
	h.logger.Info("GM searched players", "gm_user", gmUser.Username, "search_params", req)

	c.JSON(200, gin.H{"data": response, "success": true})
}

// GetPlayerDetail 获取玩家详细信息
func (h *PlayerManagementHandler) GetPlayerDetail(c *gin.Context) {
	playerID := c.Param("id")
	if playerID == "" {
		c.JSON(400, gin.H{"error": "Player ID is required", "success": false})
		return
	}

	ctx := context.Background()

	// 查询玩家详细信息
	query := &player.GetPlayerDetailQuery{PlayerID: playerID}
	result, err := handlers.ExecuteQueryTyped[*player.GetPlayerDetailQuery, *player.GetPlayerDetailResult](ctx, h.queryBus, query)
	if err != nil {
		h.logger.Error("Failed to get player detail", "error", err, "player_id", playerID)
		c.JSON(500, gin.H{"error": "Failed to get player detail", "success": false})
		return
	}

	if !result.Found {
		c.JSON(404, gin.H{"error": "Player not found", "success": false})
		return
	}

	// 构造详细响应
	p := result.Player
	response := &GMPlayerResponse{
		ID:       p.ID,
		Username: p.Username,
		Email:    p.Email,
		Name:     p.Name,
		Level:    p.Level,
		Exp:      p.Exp,
		Status:   p.Status,
		Position: PositionResponse{
			X: p.Position.X,
			Y: p.Position.Y,
			Z: p.Position.Z,
		},
		Stats: StatsResponse{
			HP:      p.Stats.HP,
			MaxHP:   p.Stats.MaxHP,
			MP:      p.Stats.MP,
			MaxMP:   p.Stats.MaxMP,
			Attack:  p.Stats.Attack,
			Defense: p.Stats.Defense,
			Speed:   p.Stats.Speed,
		},
		Avatar:       p.Avatar,
		Gender:       p.Gender,
		LastLoginAt:  p.LastLoginAt,
		LastLogoutAt: p.LastLogoutAt,
		OnlineTime:   p.OnlineTime,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}

	// 添加封禁信息
	if p.BanInfo != nil {
		response.BanInfo = &BanInfo{
			IsBanned:  p.BanInfo.IsBanned,
			BanType:   p.BanInfo.BanType,
			Reason:    p.BanInfo.Reason,
			BannedBy:  p.BanInfo.BannedBy,
			BannedAt:  p.BanInfo.BannedAt,
			ExpiresAt: p.BanInfo.ExpiresAt,
		}
	}

	// 记录GM操作日志
	gmUser, _ := auth.GetCurrentUser(c)
	h.logger.Info("GM viewed player detail", "gm_user", gmUser.Username, "player_id", playerID)

	c.JSON(200, gin.H{"data": response, "success": true})
}

// UpdatePlayer GM更新玩家信息
func (h *PlayerManagementHandler) UpdatePlayer(c *gin.Context) {
	playerID := c.Param("id")
	if playerID == "" {
		c.JSON(400, gin.H{"error": "Player ID is required", "success": false})
		return
	}

	var req PlayerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid update player request", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request format", "success": false})
		return
	}

	ctx := context.Background()

	// 获取GM用户信息
	gmUser, _ := auth.GetCurrentUser(c)

	// 执行更新命令
	cmd := &player.GMUpdatePlayerCommand{
		PlayerID: playerID,
		GMUserID: gmUser.PlayerID,
		GMUser:   gmUser.Username,
		Reason:   req.Reason,
		Updates: map[string]interface{}{},
	}

	// 添加需要更新的字段
	if req.Name != nil {
		cmd.Updates["name"] = *req.Name
	}
	if req.Level != nil {
		cmd.Updates["level"] = *req.Level
	}
	if req.Exp != nil {
		cmd.Updates["exp"] = *req.Exp
	}
	if req.Status != nil {
		cmd.Updates["status"] = *req.Status
	}
	if req.HP != nil {
		cmd.Updates["hp"] = *req.HP
	}
	if req.MaxHP != nil {
		cmd.Updates["max_hp"] = *req.MaxHP
	}
	if req.MP != nil {
		cmd.Updates["mp"] = *req.MP
	}
	if req.MaxMP != nil {
		cmd.Updates["max_mp"] = *req.MaxMP
	}
	if req.Attack != nil {
		cmd.Updates["attack"] = *req.Attack
	}
	if req.Defense != nil {
		cmd.Updates["defense"] = *req.Defense
	}
	if req.Speed != nil {
		cmd.Updates["speed"] = *req.Speed
	}

	result, err := handlers.ExecuteTyped[*player.GMUpdatePlayerCommand, *player.GMUpdatePlayerResult](ctx, h.commandBus, cmd)
	if err != nil {
		h.logger.Error("Failed to update player", "error", err, "player_id", playerID, "gm_user", gmUser.Username)
		c.JSON(500, gin.H{"error": "Failed to update player", "success": false})
		return
	}

	// 记录GM操作日志
	h.logger.Info("GM updated player", "gm_user", gmUser.Username, "player_id", playerID, "updates", cmd.Updates, "reason", req.Reason)

	c.JSON(200, gin.H{"data": result, "success": true, "message": "Player updated successfully"})
}

// BanPlayer 封禁玩家
func (h *PlayerManagementHandler) BanPlayer(c *gin.Context) {
	var req PlayerBanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid ban player request", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request format", "success": false})
		return
	}

	ctx := context.Background()

	// 获取GM用户信息
	gmUser, _ := auth.GetCurrentUser(c)

	// 计算封禁过期时间
	var expiresAt *time.Time
	if req.Duration > 0 {
		expiry := time.Now().Add(time.Duration(req.Duration) * time.Hour)
		expiresAt = &expiry
	} else if !req.ExpiresAt.IsZero() {
		expiresAt = &req.ExpiresAt
	}

	// 执行封禁命令
	cmd := &player.BanPlayerCommand{
		PlayerID:  req.PlayerID,
		BannedBy:  gmUser.PlayerID,
		BanType:   req.BanType,
		Reason:    req.Reason,
		ExpiresAt: expiresAt,
	}

	result, err := handlers.ExecuteTyped[*player.BanPlayerCommand, *player.BanPlayerResult](ctx, h.commandBus, cmd)
	if err != nil {
		h.logger.Error("Failed to ban player", "error", err, "player_id", req.PlayerID, "gm_user", gmUser.Username)
		c.JSON(500, gin.H{"error": "Failed to ban player", "success": false})
		return
	}

	// 记录GM操作日志
	h.logger.Info("GM banned player", "gm_user", gmUser.Username, "player_id", req.PlayerID, "ban_type", req.BanType, "reason", req.Reason, "expires_at", expiresAt)

	c.JSON(200, gin.H{"data": result, "success": true, "message": "Player banned successfully"})
}

// UnbanPlayer 解封玩家
func (h *PlayerManagementHandler) UnbanPlayer(c *gin.Context) {
	var req PlayerUnbanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid unban player request", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request format", "success": false})
		return
	}

	ctx := context.Background()

	// 获取GM用户信息
	gmUser, _ := auth.GetCurrentUser(c)

	// 执行解封命令
	cmd := &player.UnbanPlayerCommand{
		PlayerID:   req.PlayerID,
		UnbannedBy: gmUser.PlayerID,
		Reason:     req.Reason,
	}

	result, err := handlers.ExecuteTyped[*player.UnbanPlayerCommand, *player.UnbanPlayerResult](ctx, h.commandBus, cmd)
	if err != nil {
		h.logger.Error("Failed to unban player", "error", err, "player_id", req.PlayerID, "gm_user", gmUser.Username)
		c.JSON(500, gin.H{"error": "Failed to unban player", "success": false})
		return
	}

	// 记录GM操作日志
	h.logger.Info("GM unbanned player", "gm_user", gmUser.Username, "player_id", req.PlayerID, "reason", req.Reason)

	c.JSON(200, gin.H{"data": result, "success": true, "message": "Player unbanned successfully"})
}