package gm

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	playerCmd "greatestworks/application/commands/player"
	"greatestworks/application/handlers"
	playerQuery "greatestworks/application/queries/player"
	"greatestworks/internal/infrastructure/logging"
)

// PlayerManagementHandler GM玩家管理处理器
type PlayerManagementHandler struct {
	commandBus *handlers.CommandBus
	queryBus   *handlers.QueryBus
	logger     logging.Logger
}

// NewPlayerManagementHandler 创建GM玩家管理处理器
func NewPlayerManagementHandler(commandBus *handlers.CommandBus, queryBus *handlers.QueryBus, logger logging.Logger) *PlayerManagementHandler {
	return &PlayerManagementHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
		logger:     logger,
	}
}

// CreatePlayer 创建玩家
func (h *PlayerManagementHandler) CreatePlayer(c *gin.Context) {
	var req CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid create player request", logging.Fields{
			"error": err,
		})
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	// 创建命令
	cmd := &playerCmd.CreatePlayerCommand{
		Name:   req.Name,
		Avatar: req.Avatar,
		Gender: req.Gender,
	}

	// 执行命令
	result, err := h.commandBus.Execute(c.Request.Context(), cmd)
	if err != nil {
		h.logger.Error("Failed to create player", err, logging.Fields{
			"name": req.Name,
		})
		c.JSON(500, gin.H{"error": "Failed to create player"})
		return
	}

	h.logger.Info("Player created successfully", logging.Fields{
		"player_id": result.(*playerCmd.CreatePlayerResult).PlayerID,
		"name":      req.Name,
	})
	c.JSON(200, result)
}

// GetPlayer 获取玩家信息
func (h *PlayerManagementHandler) GetPlayer(c *gin.Context) {
	playerID := c.Param("id")
	if playerID == "" {
		c.JSON(400, gin.H{"error": "Player ID is required"})
		return
	}

	// 创建查询
	query := &playerQuery.GetPlayerQuery{
		PlayerID: playerID,
	}

	// 执行查询
	result, err := h.queryBus.Execute(c.Request.Context(), query)
	if err != nil {
		h.logger.Error("Failed to get player", err, logging.Fields{
			"player_id": playerID,
		})
		c.JSON(500, gin.H{"error": "Failed to get player"})
		return
	}

	c.JSON(200, result)
}

// ListPlayers 获取玩家列表
func (h *PlayerManagementHandler) ListPlayers(c *gin.Context) {
	// 解析查询参数
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")
	search := c.Query("search")

	// 创建查询
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)
	query := &playerQuery.ListPlayersQuery{
		Page:     pageInt,
		PageSize: limitInt,
		Name:     search,
	}

	// 执行查询
	result, err := h.queryBus.Execute(c.Request.Context(), query)
	if err != nil {
		h.logger.Error("Failed to list players", err, logging.Fields{})
		c.JSON(500, gin.H{"error": "Failed to list players"})
		return
	}

	c.JSON(200, result)
}

// UpdatePlayer 更新玩家信息
func (h *PlayerManagementHandler) UpdatePlayer(c *gin.Context) {
	playerID := c.Param("id")
	if playerID == "" {
		c.JSON(400, gin.H{"error": "Player ID is required"})
		return
	}

	var req UpdatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid update player request", logging.Fields{
			"error": err,
		})
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	// 创建命令
	cmd := &playerCmd.UpdatePlayerCommand{
		PlayerID: playerID,
		Name:     req.Name,
	}

	// 执行命令
	result, err := h.commandBus.Execute(c.Request.Context(), cmd)
	if err != nil {
		h.logger.Error("Failed to update player", err, logging.Fields{
			"player_id": playerID,
		})
		c.JSON(500, gin.H{"error": "Failed to update player"})
		return
	}

	h.logger.Info("Player updated successfully", logging.Fields{
		"player_id": playerID,
	})
	c.JSON(200, result)
}

// DeletePlayer 删除玩家
func (h *PlayerManagementHandler) DeletePlayer(c *gin.Context) {
	playerID := c.Param("id")
	if playerID == "" {
		c.JSON(400, gin.H{"error": "Player ID is required"})
		return
	}

	// 创建命令
	cmd := &playerCmd.DeletePlayerCommand{
		PlayerID: playerID,
	}

	// 执行命令
	result, err := h.commandBus.Execute(c.Request.Context(), cmd)
	if err != nil {
		h.logger.Error("Failed to delete player", err, logging.Fields{
			"player_id": playerID,
		})
		c.JSON(500, gin.H{"error": "Failed to delete player"})
		return
	}

	h.logger.Info("Player deleted successfully", logging.Fields{
		"player_id": playerID,
	})
	c.JSON(200, result)
}

// BanPlayer 封禁玩家
func (h *PlayerManagementHandler) BanPlayer(c *gin.Context) {
	playerID := c.Param("id")
	if playerID == "" {
		c.JSON(400, gin.H{"error": "Player ID is required"})
		return
	}

	var req BanPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid ban player request", logging.Fields{"error": err})
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	// 创建命令
	banType := "temporary"
	if req.Permanent {
		banType = "permanent"
	}
	cmd := &playerCmd.BanPlayerCommand{
		PlayerID:     playerID,
		BannedBy:     "GM",
		BannedByName: "GameMaster",
		Reason:       req.Reason,
		BanType:      banType,
		BanUntil:     time.Now().Add(req.Duration),
	}

	// 执行命令
	result, err := h.commandBus.Execute(c.Request.Context(), cmd)
	if err != nil {
		h.logger.Error("Failed to ban player", err, logging.Fields{"player_id": playerID})
		c.JSON(500, gin.H{"error": "Failed to ban player"})
		return
	}

	h.logger.Info("Player banned successfully", logging.Fields{"player_id": playerID, "reason": req.Reason})
	c.JSON(200, result)
}

// UnbanPlayer 解封玩家
func (h *PlayerManagementHandler) UnbanPlayer(c *gin.Context) {
	playerID := c.Param("id")
	if playerID == "" {
		c.JSON(400, gin.H{"error": "Player ID is required"})
		return
	}

	// 创建命令
	cmd := &playerCmd.UnbanPlayerCommand{
		PlayerID: playerID,
	}

	// 执行命令
	result, err := h.commandBus.Execute(c.Request.Context(), cmd)
	if err != nil {
		h.logger.Error("Failed to unban player", err, logging.Fields{"player_id": playerID})
		c.JSON(500, gin.H{"error": "Failed to unban player"})
		return
	}

	h.logger.Info("Player unbanned successfully", logging.Fields{"player_id": playerID})
	c.JSON(200, result)
}

// MovePlayer 移动玩家
func (h *PlayerManagementHandler) MovePlayer(c *gin.Context) {
	playerID := c.Param("id")
	if playerID == "" {
		c.JSON(400, gin.H{"error": "Player ID is required"})
		return
	}

	var req MovePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid move player request", logging.Fields{"error": err})
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	// 创建命令
	cmd := &playerCmd.MovePlayerCommand{
		PlayerID: playerID,
		Position: playerCmd.Position{X: req.X, Y: req.Y, Z: req.Z},
	}

	// 执行命令
	result, err := h.commandBus.Execute(c.Request.Context(), cmd)
	if err != nil {
		h.logger.Error("Failed to move player", err, logging.Fields{"player_id": playerID})
		c.JSON(500, gin.H{"error": "Failed to move player"})
		return
	}

	h.logger.Info("Player moved successfully", logging.Fields{
		"player_id": playerID,
		"position":  req,
	})
	c.JSON(200, result)
}

// LevelUpPlayer 升级玩家
func (h *PlayerManagementHandler) LevelUpPlayer(c *gin.Context) {
	playerID := c.Param("id")
	if playerID == "" {
		c.JSON(400, gin.H{"error": "Player ID is required"})
		return
	}

	var req LevelUpPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid level up player request", logging.Fields{
			"error": err,
		})
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	// 创建命令
	cmd := &playerCmd.LevelUpPlayerCommand{
		PlayerID: playerID,
		ExpGain:  int64(req.Levels * 1000), // 假设每级需要1000经验
	}

	// 执行命令
	result, err := h.commandBus.Execute(c.Request.Context(), cmd)
	if err != nil {
		h.logger.Error("Failed to level up player", err, logging.Fields{
			"player_id": playerID,
		})
		c.JSON(500, gin.H{"error": "Failed to level up player"})
		return
	}

	h.logger.Info("Player leveled up successfully", logging.Fields{
		"player_id": playerID,
		"levels":    req.Levels,
	})
	c.JSON(200, result)
}

// 请求和响应结构体

// CreatePlayerRequest 创建玩家请求
type CreatePlayerRequest struct {
	Name   string `json:"name" binding:"required"`
	Avatar string `json:"avatar,omitempty"`
	Gender int    `json:"gender,omitempty"`
}

// UpdatePlayerRequest 更新玩家请求
type UpdatePlayerRequest struct {
	Name  string `json:"name,omitempty"`
	Level int    `json:"level,omitempty"`
	Exp   int64  `json:"exp,omitempty"`
}

// BanPlayerRequest 封禁玩家请求
type BanPlayerRequest struct {
	Reason    string        `json:"reason" binding:"required"`
	Duration  time.Duration `json:"duration,omitempty"`
	Permanent bool          `json:"permanent,omitempty"`
}

// MovePlayerRequest 移动玩家请求
type MovePlayerRequest struct {
	X float64 `json:"x" binding:"required"`
	Y float64 `json:"y" binding:"required"`
	Z float64 `json:"z" binding:"required"`
}

// LevelUpPlayerRequest 升级玩家请求
type LevelUpPlayerRequest struct {
	Levels int `json:"levels" binding:"required,min=1"`
}
