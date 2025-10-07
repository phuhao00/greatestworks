package http

import (
	"context"
	"time"
	
	"github.com/gin-gonic/gin"
	
	playerCmd "greatestworks/application/commands/player"
	playerQuery "greatestworks/application/queries/player"
	"greatestworks/application/handlers"
	"greatestworks/internal/domain/player"
	"greatestworks/internal/infrastructure/logger"
)

// PlayerHandler 玩家HTTP处理器
type PlayerHandler struct {
	commandBus *handlers.CommandBus
	queryBus   *handlers.QueryBus
	logger     logger.Logger
}

// NewPlayerHandler 创建玩家处理器
func NewPlayerHandler(commandBus *handlers.CommandBus, queryBus *handlers.QueryBus, logger logger.Logger) *PlayerHandler {
	return &PlayerHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
		logger:     logger,
	}
}

// CreatePlayerRequest 创建玩家请求
type CreatePlayerRequest struct {
	Name   string `json:"name" binding:"required,min=2,max=50"`
	Avatar string `json:"avatar,omitempty"`
	Gender int    `json:"gender,omitempty" binding:"min=0,max=2"`
}

// UpdatePlayerRequest 更新玩家请求
type UpdatePlayerRequest struct {
	Name   string `json:"name,omitempty" binding:"omitempty,min=2,max=50"`
	Avatar string `json:"avatar,omitempty"`
	Gender *int   `json:"gender,omitempty" binding:"omitempty,min=0,max=2"`
}

// MovePlayerRequest 移动玩家请求
type MovePlayerRequest struct {
	X float64 `json:"x" binding:"required"`
	Y float64 `json:"y" binding:"required"`
	Z float64 `json:"z" binding:"required"`
}

// PlayerResponse 玩家响应
type PlayerResponse struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Level     int               `json:"level"`
	Exp       int64             `json:"exp"`
	Status    string            `json:"status"`
	Position  PositionResponse  `json:"position"`
	Stats     StatsResponse     `json:"stats"`
	Avatar    string            `json:"avatar,omitempty"`
	Gender    int               `json:"gender,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
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

// PlayerListRequest 玩家列表请求
type PlayerListRequest struct {
	PaginationRequest
	Name   string `form:"name,omitempty"`
	Status string `form:"status,omitempty"`
	Level  int    `form:"level,omitempty"`
}

// CreatePlayer 创建玩家
func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var req CreatePlayerRequest
	if !BindAndValidate(c, &req) {
		return
	}
	
	ctx := context.Background()
	
	// 执行创建玩家命令
	cmd := &playerCmd.CreatePlayerCommand{
		Name:   req.Name,
		Avatar: req.Avatar,
		Gender: req.Gender,
	}
	
	result, err := handlers.ExecuteTyped[*playerCmd.CreatePlayerCommand, *playerCmd.CreatePlayerResult](ctx, h.commandBus, cmd)
	if err != nil {
		h.logger.Error("Failed to create player", "error", err, "name", req.Name)
		HandleError(c, err)
		return
	}
	
	// 构造响应
	response := &PlayerResponse{
		ID:        result.PlayerID,
		Name:      result.Name,
		Level:     result.Level,
		Exp:       0,
		Status:    "active",
		Position:  PositionResponse{X: 0, Y: 0, Z: 0},
		Stats:     StatsResponse{HP: 100, MaxHP: 100, MP: 50, MaxMP: 50, Attack: 10, Defense: 5, Speed: 10},
		Avatar:    req.Avatar,
		Gender:    req.Gender,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.CreatedAt,
	}
	
	CreatedResponse(c, response, "Player created successfully")
}

// GetPlayer 获取玩家信息
func (h *PlayerHandler) GetPlayer(c *gin.Context) {
	playerID, ok := ValidateID(c, "id")
	if !ok {
		return
	}
	
	ctx := context.Background()
	
	// 查询玩家信息
	query := &playerQuery.GetPlayerQuery{PlayerID: playerID}
	result, err := handlers.ExecuteQueryTyped[*playerQuery.GetPlayerQuery, *playerQuery.GetPlayerResult](ctx, h.queryBus, query)
	if err != nil {
		h.logger.Error("Failed to get player", "error", err, "player_id", playerID)
		HandleError(c, err)
		return
	}
	
	if !result.Found {
		NotFoundResponse(c, "Player not found")
		return
	}
	
	// 构造响应
	response := &PlayerResponse{
		ID:       result.Player.ID,
		Name:     result.Player.Name,
		Level:    result.Player.Level,
		Exp:      result.Player.Exp,
		Status:   result.Player.Status,
		Position: PositionResponse{
			X: result.Player.Position.X,
			Y: result.Player.Position.Y,
			Z: result.Player.Position.Z,
		},
		Stats: StatsResponse{
			HP:      result.Player.Stats.HP,
			MaxHP:   result.Player.Stats.MaxHP,
			MP:      result.Player.Stats.MP,
			MaxMP:   result.Player.Stats.MaxMP,
			Attack:  result.Player.Stats.Attack,
			Defense: result.Player.Stats.Defense,
			Speed:   result.Player.Stats.Speed,
		},
		CreatedAt: result.Player.CreatedAt,
		UpdatedAt: result.Player.UpdatedAt,
	}
	
	SuccessResponse(c, response)
}

// UpdatePlayer 更新玩家信息
func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	playerID, ok := ValidateID(c, "id")
	if !ok {
		return
	}
	
	var req UpdatePlayerRequest
	if !BindAndValidate(c, &req) {
		return
	}
	
	ctx := context.Background()
	
	// 执行更新玩家命令
	cmd := &playerCmd.UpdatePlayerCommand{
		PlayerID: playerID,
		Name:     req.Name,
		Avatar:   req.Avatar,
	}
	
	if req.Gender != nil {
		cmd.Gender = *req.Gender
	}
	
	result, err := handlers.ExecuteTyped[*playerCmd.UpdatePlayerCommand, *playerCmd.UpdatePlayerResult](ctx, h.commandBus, cmd)
	if err != nil {
		h.logger.Error("Failed to update player", "error", err, "player_id", playerID)
		HandleError(c, err)
		return
	}
	
	// 构造响应
	response := &PlayerResponse{
		ID:        result.PlayerID,
		Name:      result.Name,
		Level:     result.Level,
		Exp:       result.Exp,
		Status:    result.Status,
		UpdatedAt: result.UpdatedAt,
	}
	
	SuccessResponse(c, response, "Player updated successfully")
}

// DeletePlayer 删除玩家
func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	playerIDStr, ok := ValidateID(c, "id")
	if !ok {
		return
	}
	
	ctx := context.Background()
	
	// 执行删除玩家命令
	cmd := &playerCmd.DeletePlayerCommand{
		PlayerID: playerIDStr,
	}
	_, err := handlers.ExecuteTyped[*playerCmd.DeletePlayerCommand, *playerCmd.DeletePlayerResult](ctx, h.commandBus, cmd)
	if err != nil {
		h.logger.Error("Failed to delete player", "error", err, "player_id", playerIDStr)
		HandleError(c, err)
		return
	}
	
	NoContentResponse(c, "Player deleted successfully")
}

// ListPlayers 获取玩家列表
func (h *PlayerHandler) ListPlayers(c *gin.Context) {
	var req PlayerListRequest
	if !BindQueryAndValidate(c, &req) {
		return
	}
	
	ctx := context.Background()
	page, pageSize := req.GetPagination()
	
	// 查询玩家列表
	query := &playerQuery.ListPlayersQuery{
		Page:     page,
		PageSize: pageSize,
		Name:     req.Name,
		Status:   req.Status,
		Level:    req.Level,
	}
	
	result, err := handlers.ExecuteQueryTyped[*playerQuery.ListPlayersQuery, *playerQuery.ListPlayersResult](ctx, h.queryBus, query)
	if err != nil {
		h.logger.Error("Failed to list players", "error", err)
		HandleError(c, err)
		return
	}
	
	// 构造响应
	var players []*PlayerResponse
	for _, p := range result.Players {
		players = append(players, &PlayerResponse{
			ID:       p.ID,
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
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}
	
	meta := CreateMeta(page, pageSize, result.Total)
	SuccessResponseWithMeta(c, players, meta)
}

// MovePlayer 移动玩家
func (h *PlayerHandler) MovePlayer(c *gin.Context) {
	playerIDStr, ok := ValidateID(c, "id")
	if !ok {
		return
	}
	
	var req MovePlayerRequest
	if !BindAndValidate(c, &req) {
		return
	}
	
	ctx := context.Background()
	
	// 执行移动玩家命令
	cmd := &playerCmd.MovePlayerCommand{
		PlayerID: playerIDStr,
		Position: playerCmd.Position{
			X: req.X,
			Y: req.Y,
			Z: req.Z,
		},
	}
	
	result, err := handlers.ExecuteTyped[*playerCmd.MovePlayerCommand, *playerCmd.MovePlayerResult](ctx, h.commandBus, cmd)
	if err != nil {
		h.logger.Error("Failed to move player", "error", err, "player_id", playerIDStr)
		HandleError(c, err)
		return
	}
	
	// 构造响应
	response := map[string]interface{}{
		"success": result.Success,
		"old_position": PositionResponse{
			X: result.OldPosition.X,
			Y: result.OldPosition.Y,
			Z: result.OldPosition.Z,
		},
		"new_position": PositionResponse{
			X: result.NewPosition.X,
			Y: result.NewPosition.Y,
			Z: result.NewPosition.Z,
		},
		"moved_at": time.Now(),
	}
	
	SuccessResponse(c, response, "Player moved successfully")
}

// LevelUpPlayer 玩家升级
func (h *PlayerHandler) LevelUpPlayer(c *gin.Context) {
	playerIDStr, ok := ValidateID(c, "id")
	if !ok {
		return
	}
	
	ctx := context.Background()
	
	// 执行玩家升级命令
	cmd := &playerCmd.LevelUpPlayerCommand{PlayerID: playerIDStr}
	result, err := handlers.ExecuteTyped[*playerCmd.LevelUpPlayerCommand, *playerCmd.LevelUpPlayerResult](ctx, h.commandBus, cmd)
	if err != nil {
		h.logger.Error("Failed to level up player", "error", err, "player_id", playerIDStr)
		HandleError(c, err)
		return
	}
	
	// 构造响应
	response := map[string]interface{}{
		"success":    result.LeveledUp,
		"old_level":  result.OldLevel,
		"new_level":  result.NewLevel,
		"old_exp":    result.OldExp,
		"new_exp":    result.NewExp,
		"leveled_up_at": time.Now(),
	}
	
	SuccessResponse(c, response, "Player leveled up successfully")
}

// GetPlayerStats 获取玩家统计信息
func (h *PlayerHandler) GetPlayerStats(c *gin.Context) {
	playerIDStr, ok := ValidateID(c, "id")
	if !ok {
		return
	}
	
	ctx := context.Background()
	
	// 查询玩家统计信息
	query := &playerQuery.GetPlayerStatsQuery{
		PlayerID: player.PlayerIDFromString(playerIDStr),
	}
	result, err := handlers.ExecuteQueryTyped[*playerQuery.GetPlayerStatsQuery, *playerQuery.GetPlayerStatsResult](ctx, h.queryBus, query)
	if err != nil {
		h.logger.Error("Failed to get player stats", "error", err, "player_id", playerIDStr)
		HandleError(c, err)
		return
	}
	
	if !result.Found {
		NotFoundResponse(c, "Player stats not found")
		return
	}
	
	// 构造响应
	response := map[string]interface{}{
		"player_id":     result.PlayerID,
		"total_battles": result.TotalBattles,
		"wins":          result.Wins,
		"losses":        result.Losses,
		"win_rate":      result.WinRate,
		"total_exp":     result.TotalExp,
		"play_time":     result.PlayTime,
		"last_login":    result.LastLogin,
		"achievements":  result.Achievements,
	}
	
	SuccessResponse(c, response)
}