package tcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	
	"greatestworks/application/commands/player"
	"greatestworks/application/commands/battle"
	"greatestworks/application/queries/player"
	"greatestworks/application/handlers"
	"greatestworks/internal/infrastructure/network"
	"greatestworks/internal/interfaces/tcp/protocol"
)

// GameHandler 游戏处理器
type GameHandler struct {
	commandBus *handlers.CommandBus
	queryBus   *handlers.QueryBus
	logger     Logger
}

// Logger 日志接口
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}

// NewGameHandler 创建游戏处理器
func NewGameHandler(commandBus *handlers.CommandBus, queryBus *handlers.QueryBus, logger Logger) *GameHandler {
	return &GameHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
		logger:     logger,
	}
}

// RegisterHandlers 注册处理器
func (h *GameHandler) RegisterHandlers(server network.Server) error {
	// 玩家相关协议
	server.RegisterHandler(protocol.MsgPlayerLogin, h.handlePlayerLogin)
	server.RegisterHandler(protocol.MsgPlayerLogout, h.handlePlayerLogout)
	server.RegisterHandler(protocol.MsgPlayerMove, h.handlePlayerMove)
	server.RegisterHandler(protocol.MsgPlayerInfo, h.handlePlayerInfo)
	server.RegisterHandler(protocol.MsgPlayerCreate, h.handlePlayerCreate)
	
	// 战斗相关协议
	server.RegisterHandler(protocol.MsgCreateBattle, h.handleCreateBattle)
	server.RegisterHandler(protocol.MsgJoinBattle, h.handleJoinBattle)
	server.RegisterHandler(protocol.MsgStartBattle, h.handleStartBattle)
	server.RegisterHandler(protocol.MsgBattleAction, h.handleBattleAction)
	server.RegisterHandler(protocol.MsgLeaveBattle, h.handleLeaveBattle)
	
	// 查询相关协议
	server.RegisterHandler(protocol.MsgGetPlayerInfo, h.handleGetPlayerInfo)
	server.RegisterHandler(protocol.MsgGetOnlinePlayers, h.handleGetOnlinePlayers)
	server.RegisterHandler(protocol.MsgGetBattleInfo, h.handleGetBattleInfo)
	
	h.logger.Info("Game handlers registered successfully")
	return nil
}

// 玩家协议处理器

// handlePlayerLogin 处理玩家登录
func (h *GameHandler) handlePlayerLogin(ctx context.Context, conn network.Connection, packet *network.Packet) error {
	var req protocol.PlayerLoginRequest
	if err := json.Unmarshal(packet.Data, &req); err != nil {
		h.logger.Error("Failed to unmarshal player login request", "error", err)
		return h.sendErrorResponse(conn, protocol.MsgPlayerLogin, "Invalid request format")
	}
	
	// 验证玩家登录（这里简化处理）
	if req.PlayerID == "" {
		return h.sendErrorResponse(conn, protocol.MsgPlayerLogin, "Player ID is required")
	}
	
	// 查询玩家信息
	query := &player.GetPlayerQuery{PlayerID: req.PlayerID}
	result, err := handlers.ExecuteQueryTyped[*player.GetPlayerQuery, *player.GetPlayerResult](ctx, h.queryBus, query)
	if err != nil {
		h.logger.Error("Failed to get player info", "error", err, "player_id", req.PlayerID)
		return h.sendErrorResponse(conn, protocol.MsgPlayerLogin, "Failed to get player info")
	}
	
	if !result.Found {
		return h.sendErrorResponse(conn, protocol.MsgPlayerLogin, "Player not found")
	}
	
	// 设置连接的玩家ID
	conn.SetAttribute("player_id", req.PlayerID)
	
	// 构造响应
	response := &protocol.PlayerLoginResponse{
		BaseResponse: protocol.NewBaseResponse(true, "Login successful"),
		Player: &protocol.PlayerInfo{
			ID:     result.Player.ID,
			Name:   result.Player.Name,
			Level:  result.Player.Level,
			Exp:    result.Player.Exp,
			Status: result.Player.Status,
			Position: protocol.Position{
				X: result.Player.Position.X,
				Y: result.Player.Position.Y,
				Z: result.Player.Position.Z,
			},
			Stats: protocol.Stats{
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
		},
		SessionID:  fmt.Sprintf("session_%s_%d", req.PlayerID, time.Now().Unix()),
		ServerTime: time.Now().Unix(),
	}
	
	return h.sendResponse(conn, protocol.MsgPlayerLogin, response)
}

// handlePlayerMove 处理玩家移动
func (h *GameHandler) handlePlayerMove(ctx context.Context, conn network.Connection, packet *network.Packet) error {
	playerID, ok := conn.GetAttribute("player_id").(string)
	if !ok || playerID == "" {
		return h.sendErrorResponse(conn, protocol.MsgPlayerMove, "Not logged in")
	}
	
	var req protocol.PlayerMoveRequest
	if err := json.Unmarshal(packet.Data, &req); err != nil {
		h.logger.Error("Failed to unmarshal player move request", "error", err)
		return h.sendErrorResponse(conn, protocol.MsgPlayerMove, "Invalid request format")
	}
	
	// 执行移动命令
	cmd := &player.MovePlayerCommand{
		PlayerID: playerID,
		Position: player.Position{
			X: req.Position.X,
			Y: req.Position.Y,
			Z: req.Position.Z,
		},
	}
	
	result, err := handlers.ExecuteTyped[*player.MovePlayerCommand, *player.MovePlayerResult](ctx, h.commandBus, cmd)
	if err != nil {
		h.logger.Error("Failed to move player", "error", err, "player_id", playerID)
		return h.sendErrorResponse(conn, protocol.MsgPlayerMove, "Failed to move player")
	}
	
	// 构造响应
	response := &protocol.PlayerMoveResponse{
		BaseResponse: protocol.NewBaseResponse(result.Success, "Move successful"),
		OldPosition: protocol.Position{
			X: result.OldPosition.X,
			Y: result.OldPosition.Y,
			Z: result.OldPosition.Z,
		},
		NewPosition: protocol.Position{
			X: result.NewPosition.X,
			Y: result.NewPosition.Y,
			Z: result.NewPosition.Z,
		},
		MoveTime: time.Now().Unix(),
	}
	
	return h.sendResponse(conn, protocol.MsgPlayerMove, response)
}

// handlePlayerCreate 处理玩家创建
func (h *GameHandler) handlePlayerCreate(ctx context.Context, conn network.Connection, packet *network.Packet) error {
	var req protocol.PlayerCreateRequest
	if err := json.Unmarshal(packet.Data, &req); err != nil {
		h.logger.Error("Failed to unmarshal create player request", "error", err)
		return h.sendErrorResponse(conn, protocol.MsgPlayerCreate, "Invalid request format")
	}
	
	// 执行创建玩家命令
	cmd := &player.CreatePlayerCommand{
		Name: req.Name,
	}
	
	result, err := handlers.ExecuteTyped[*player.CreatePlayerCommand, *player.CreatePlayerResult](ctx, h.commandBus, cmd)
	if err != nil {
		h.logger.Error("Failed to create player", "error", err, "name", req.Name)
		return h.sendErrorResponse(conn, protocol.MsgPlayerCreate, "Failed to create player")
	}
	
	// 构造响应
	response := &protocol.PlayerCreateResponse{
		BaseResponse: protocol.NewBaseResponse(true, "Player created successfully"),
		PlayerID:     result.PlayerID,
		Name:         result.Name,
		Level:        result.Level,
		CreatedAt:    result.CreatedAt,
	}
	
	return h.sendResponse(conn, protocol.MsgPlayerCreate, response)
}

// handlePlayerInfo 处理获取玩家信息
func (h *GameHandler) handlePlayerInfo(ctx context.Context, conn network.Connection, packet *network.Packet) error {
	playerID, ok := conn.GetAttribute("player_id").(string)
	if !ok || playerID == "" {
		return h.sendErrorResponse(conn, protocol.MsgPlayerInfo, "Not logged in")
	}
	
	// 查询玩家信息
	query := &player.GetPlayerQuery{PlayerID: playerID}
	result, err := handlers.ExecuteQueryTyped[*player.GetPlayerQuery, *player.GetPlayerResult](ctx, h.queryBus, query)
	if err != nil {
		h.logger.Error("Failed to get player info", "error", err, "player_id", playerID)
		return h.sendErrorResponse(conn, protocol.MsgPlayerInfo, "Failed to get player info")
	}
	
	if !result.Found {
		return h.sendErrorResponse(conn, protocol.MsgPlayerInfo, "Player not found")
	}
	
	// 构造响应
	response := &protocol.PlayerInfoResponse{
		BaseResponse: protocol.NewBaseResponse(true, "Player info retrieved successfully"),
		Player: &protocol.PlayerInfo{
			ID:     result.Player.ID,
			Name:   result.Player.Name,
			Level:  result.Player.Level,
			Exp:    result.Player.Exp,
			Status: result.Player.Status,
			Position: protocol.Position{
				X: result.Player.Position.X,
				Y: result.Player.Position.Y,
				Z: result.Player.Position.Z,
			},
			Stats: protocol.Stats{
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
		},
	}
	
	return h.sendResponse(conn, protocol.MsgPlayerInfo, response)
}

// handlePlayerLogout 处理玩家登出
func (h *GameHandler) handlePlayerLogout(ctx context.Context, conn network.Connection, packet *network.Packet) error {
	playerID, ok := conn.GetAttribute("player_id").(string)
	if !ok || playerID == "" {
		return h.sendErrorResponse(conn, protocol.MsgPlayerLogout, "Not logged in")
	}
	
	// 清除连接属性
	conn.SetAttribute("player_id", nil)
	
	// 构造响应
	response := protocol.NewBaseResponse(true, "Logout successful")
	return h.sendResponse(conn, protocol.MsgPlayerLogout, response)
}

// 战斗协议处理器

// handleCreateBattle 处理创建战斗
func (h *GameHandler) handleCreateBattle(ctx context.Context, conn network.Connection, packet *network.Packet) error {
	playerID, ok := conn.GetAttribute("player_id").(string)
	if !ok || playerID == "" {
		return h.sendErrorResponse(conn, protocol.MsgCreateBattle, "Not logged in")
	}
	
	var req protocol.CreateBattleRequest
	if err := json.Unmarshal(packet.Data, &req); err != nil {
		h.logger.Error("Failed to unmarshal create battle request", "error", err)
		return h.sendErrorResponse(conn, protocol.MsgCreateBattle, "Invalid request format")
	}
	
	// 执行创建战斗命令
	cmd := &battle.CreateBattleCommand{
		BattleType: battle.BattleType(req.BattleType),
		CreatorID:  playerID,
	}
	
	result, err := handlers.ExecuteTyped[*battle.CreateBattleCommand, *battle.CreateBattleResult](ctx, h.commandBus, cmd)
	if err != nil {
		h.logger.Error("Failed to create battle", "error", err, "player_id", playerID)
		return h.sendErrorResponse(conn, protocol.MsgCreateBattle, "Failed to create battle")
	}
	
	// 构造响应
	response := &protocol.CreateBattleResponse{
		BaseResponse: protocol.NewBaseResponse(true, "Battle created successfully"),
		BattleID:     result.BattleID,
		BattleType:   int(result.BattleType),
		Status:       result.Status,
		CreatedAt:    result.CreatedAt,
	}
	
	return h.sendResponse(conn, protocol.MsgCreateBattle, response)
}

// 辅助方法

// sendResponse 发送响应
func (h *GameHandler) sendResponse(conn network.Connection, msgType uint32, data interface{}) error {
	responseData, err := json.Marshal(data)
	if err != nil {
		h.logger.Error("Failed to marshal response", "error", err)
		return err
	}
	
	responsePacket := network.NewPacket(msgType, responseData)
	return conn.Send(responsePacket)
}

// sendErrorResponse 发送错误响应
func (h *GameHandler) sendErrorResponse(conn network.Connection, msgType uint32, message string) error {
	errorResponse := protocol.NewErrorResponse(message, 400, "BadRequest")
	return h.sendResponse(conn, msgType, errorResponse)
}

// 其他处理器的占位符实现
func (h *GameHandler) handleJoinBattle(ctx context.Context, conn network.Connection, packet *network.Packet) error {
	return h.sendErrorResponse(conn, protocol.MsgJoinBattle, "Not implemented")
}

func (h *GameHandler) handleStartBattle(ctx context.Context, conn network.Connection, packet *network.Packet) error {
	return h.sendErrorResponse(conn, protocol.MsgStartBattle, "Not implemented")
}

func (h *GameHandler) handleBattleAction(ctx context.Context, conn network.Connection, packet *network.Packet) error {
	return h.sendErrorResponse(conn, protocol.MsgBattleAction, "Not implemented")
}

func (h *GameHandler) handleLeaveBattle(ctx context.Context, conn network.Connection, packet *network.Packet) error {
	return h.sendErrorResponse(conn, protocol.MsgLeaveBattle, "Not implemented")
}

func (h *GameHandler) handleGetPlayerInfo(ctx context.Context, conn network.Connection, packet *network.Packet) error {
	return h.sendErrorResponse(conn, protocol.MsgGetPlayerInfo, "Not implemented")
}

func (h *GameHandler) handleGetOnlinePlayers(ctx context.Context, conn network.Connection, packet *network.Packet) error {
	return h.sendErrorResponse(conn, protocol.MsgGetOnlinePlayers, "Not implemented")
}

func (h *GameHandler) handleGetBattleInfo(ctx context.Context, conn network.Connection, packet *network.Packet) error {
	return h.sendErrorResponse(conn, protocol.MsgGetBattleInfo, "Not implemented")
}