package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	battleCmd "greatestworks/application/commands/battle"
	playerCmd "greatestworks/application/commands/player"
	"greatestworks/application/handlers"
	playerQuery "greatestworks/application/queries/player"
	"greatestworks/internal/infrastructure/logger"
	"greatestworks/internal/interfaces/tcp/connection"
	"greatestworks/internal/interfaces/tcp/protocol"
)

// GameHandler 游戏处理器
type GameHandler struct {
	commandBus *handlers.CommandBus
	queryBus   *handlers.QueryBus
	connMgr    *connection.ConnectionManager
	logger     logger.Logger
}

// NewGameHandler 创建游戏处理器
func NewGameHandler(commandBus *handlers.CommandBus, queryBus *handlers.QueryBus, connMgr *connection.ConnectionManager, logger logger.Logger) *GameHandler {
	return &GameHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
		connMgr:    connMgr,
		logger:     logger,
	}
}

// HandleMessage 处理消息
func (h *GameHandler) HandleMessage(conn *connection.Connection, msg *protocol.Message) error {
	ctx := context.Background()

	// 更新连接活动时间
	conn.UpdateActivity()

	// 根据消息类型分发处理
	switch msg.Header.MessageType {
	// 认证相关
	case protocol.MsgAuth:
		return h.handleAuth(ctx, conn, msg)

	// 玩家相关
	case protocol.MsgPlayerLogin:
		return h.handlePlayerLogin(ctx, conn, msg)
	case protocol.MsgPlayerLogout:
		return h.handlePlayerLogout(ctx, conn, msg)
	case protocol.MsgPlayerMove:
		return h.handlePlayerMove(ctx, conn, msg)
	case protocol.MsgPlayerInfo:
		return h.handlePlayerInfo(ctx, conn, msg)
	case protocol.MsgPlayerCreate:
		return h.handlePlayerCreate(ctx, conn, msg)

	// 战斗相关
	case protocol.MsgCreateBattle:
		return h.handleCreateBattle(ctx, conn, msg)
	case protocol.MsgJoinBattle:
		return h.handleJoinBattle(ctx, conn, msg)
	case protocol.MsgStartBattle:
		return h.handleStartBattle(ctx, conn, msg)
	case protocol.MsgBattleAction:
		return h.handleBattleAction(ctx, conn, msg)
	case protocol.MsgLeaveBattle:
		return h.handleLeaveBattle(ctx, conn, msg)

	// 查询相关
	case protocol.MsgGetPlayerInfo:
		return h.handleGetPlayerInfo(ctx, conn, msg)
	case protocol.MsgGetOnlinePlayers:
		return h.handleGetOnlinePlayers(ctx, conn, msg)
	case protocol.MsgGetBattleInfo:
		return h.handleGetBattleInfo(ctx, conn, msg)

	default:
		h.logger.Warn("Unknown message type", "message_type", msg.Header.MessageType, "conn_id", conn.ID)
		return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeInvalidMessage, "Unknown message type")
	}
}

// 认证处理

func (h *GameHandler) handleAuth(ctx context.Context, conn *connection.Connection, msg *protocol.Message) error {
	var req protocol.AuthRequest
	if err := h.unmarshalPayload(msg.Payload, &req); err != nil {
		h.logger.Error("Failed to unmarshal auth request", "error", err, "conn_id", conn.ID)
		return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeInvalidMessage, "Invalid request format")
	}

	// 验证Token（这里简化处理，实际应该验证JWT）
	if req.Token == "" || req.PlayerID == "" {
		h.logger.Warn("Invalid auth request", "conn_id", conn.ID, "player_id", req.PlayerID)
		return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeAuthFailed, "Invalid token or player ID")
	}

	// 查询玩家信息
	query := &playerQuery.GetPlayerQuery{PlayerID: req.PlayerID}
	result, err := handlers.ExecuteQueryTyped[*playerQuery.GetPlayerQuery, *playerQuery.GetPlayerResult](ctx, h.queryBus, query)
	if err != nil {
		h.logger.Error("Failed to get player info", "error", err, "player_id", req.PlayerID)
		return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodePlayerNotFound, "Player not found")
	}

	if !result.Found {
		return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodePlayerNotFound, "Player not found")
	}

	// 绑定玩家到连接
	if err := h.connMgr.BindPlayerToConnection(conn.ID, req.PlayerID); err != nil {
		h.logger.Error("Failed to bind player to connection", "error", err, "conn_id", conn.ID, "player_id", req.PlayerID)
		return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeServerBusy, "Failed to bind player")
	}

	// 生成会话ID
	sessionID := fmt.Sprintf("session_%s_%d", req.PlayerID, time.Now().Unix())
	conn.SessionID = sessionID

	// 构造响应
	response := &protocol.AuthResponse{
		BaseResponse: protocol.NewBaseResponse(true, "Authentication successful"),
		SessionID:    sessionID,
		PlayerInfo: &protocol.PlayerInfo{
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
		ServerTime: time.Now().Unix(),
	}

	h.logger.Info("Player authenticated", "player_id", req.PlayerID, "conn_id", conn.ID, "session_id", sessionID)
	return h.sendResponse(conn, msg.Header.MessageID, uint16(protocol.MsgAuth), response)
}

// 玩家相关处理

func (h *GameHandler) handlePlayerLogin(ctx context.Context, conn *connection.Connection, msg *protocol.Message) error {
	if !conn.IsAuthenticated {
		return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeAuthFailed, "Not authenticated")
	}

	var req protocol.PlayerLoginRequest
	if err := h.unmarshalPayload(msg.Payload, &req); err != nil {
		h.logger.Error("Failed to unmarshal player login request", "error", err)
		return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeInvalidMessage, "Invalid request format")
	}

	// 验证玩家ID是否匹配
	if req.PlayerID != conn.PlayerID {
		h.logger.Warn("Player ID mismatch", "conn_player_id", conn.PlayerID, "req_player_id", req.PlayerID)
		return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeInvalidPlayer, "Player ID mismatch")
	}

	// 查询玩家信息
	query := &playerQuery.GetPlayerQuery{PlayerID: req.PlayerID}
	result, err := handlers.ExecuteQueryTyped[*playerQuery.GetPlayerQuery, *playerQuery.GetPlayerResult](ctx, h.queryBus, query)
	if err != nil {
		h.logger.Error("Failed to get player info", "error", err, "player_id", req.PlayerID)
		return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodePlayerNotFound, "Failed to get player info")
	}

	if !result.Found {
		return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodePlayerNotFound, "Player not found")
	}

	// 更新玩家登录状态 - 暂时注释掉，因为UpdatePlayerStatusCommand不存在
	// cmd := &playerCmd.UpdatePlayerStatusCommand{
	//	PlayerID: req.PlayerID,
	//	Status:   "online",
	//	LoginTime: time.Now(),
	// }

	// _, err = handlers.ExecuteTyped[*playerCmd.UpdatePlayerStatusCommand, *playerCmd.UpdatePlayerStatusResult](ctx, h.commandBus, cmd)
	if err != nil {
		h.logger.Error("Failed to update player status", "error", err, "player_id", req.PlayerID)
		// 不返回错误，继续处理
	}

	// 构造响应
	response := &protocol.PlayerLoginResponse{
		BaseResponse: protocol.NewBaseResponse(true, "Login successful"),
		Player: &protocol.PlayerInfo{
			ID:     result.Player.ID,
			Name:   result.Player.Name,
			Level:  result.Player.Level,
			Exp:    result.Player.Exp,
			Status: "online",
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
		SessionID:  conn.SessionID,
		ServerTime: time.Now().Unix(),
	}

	h.logger.Info("Player logged in", "player_id", req.PlayerID, "conn_id", conn.ID)

	// 广播玩家上线消息给其他玩家
	h.broadcastPlayerOnline(result.Player)

	return h.sendResponse(conn, msg.Header.MessageID, uint16(protocol.MsgPlayerLogin), response)
}

func (h *GameHandler) handlePlayerMove(ctx context.Context, conn *connection.Connection, msg *protocol.Message) error {
	if !conn.IsAuthenticated {
		return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeAuthFailed, "Not authenticated")
	}

	var req protocol.PlayerMoveRequest
	if err := h.unmarshalPayload(msg.Payload, &req); err != nil {
		h.logger.Error("Failed to unmarshal player move request", "error", err)
		return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeInvalidMessage, "Invalid request format")
	}

	// 执行移动命令
	cmd := &playerCmd.MovePlayerCommand{
		PlayerID: conn.PlayerID,
		Position: playerCmd.Position{
			X: req.Position.X,
			Y: req.Position.Y,
			Z: req.Position.Z,
		},
	}

	result, err := handlers.ExecuteTyped[*playerCmd.MovePlayerCommand, *playerCmd.MovePlayerResult](ctx, h.commandBus, cmd)
	if err != nil {
		h.logger.Error("Failed to move player", "error", err, "player_id", conn.PlayerID)
		return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeUnknown, "Failed to move player")
	}

	// 构造响应
	response := &protocol.PlayerMoveResponse{
		BaseResponse: protocol.NewBaseResponse(result.Success, "Move completed"),
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

	// 广播移动消息给附近玩家
	h.broadcastPlayerMove(conn.PlayerID, result.NewPosition.X, result.NewPosition.Y)

	return h.sendResponse(conn, msg.Header.MessageID, uint16(protocol.MsgPlayerMove), response)
}

// 战斗相关处理

func (h *GameHandler) handleCreateBattle(ctx context.Context, conn *connection.Connection, msg *protocol.Message) error {
	if !conn.IsAuthenticated {
		return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeAuthFailed, "Not authenticated")
	}

	var req protocol.CreateBattleRequest
	if err := h.unmarshalPayload(msg.Payload, &req); err != nil {
		h.logger.Error("Failed to unmarshal create battle request", "error", err)
		return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeInvalidMessage, "Invalid request format")
	}

	// 执行创建战斗命令
	cmd := &battleCmd.CreateBattleCommand{
		CreatorID:  conn.PlayerID,
		BattleType: battleCmd.BattleType(req.BattleType),
	}

	result, err := handlers.ExecuteTyped[*battleCmd.CreateBattleCommand, *battleCmd.CreateBattleResult](ctx, h.commandBus, cmd)
	if err != nil {
		h.logger.Error("Failed to create battle", "error", err, "player_id", conn.PlayerID)
		return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeUnknown, "Failed to create battle")
	}

	// 构造响应
	response := &protocol.CreateBattleResponse{
		BaseResponse: protocol.NewBaseResponse(true, "Battle created successfully"),
		BattleID:     result.BattleID,
		BattleType:   int(result.BattleType),
		Status:       result.Status,
		CreatedAt:    result.CreatedAt,
	}

	h.logger.Info("Battle created", "battle_id", result.BattleID, "creator_id", conn.PlayerID, "battle_type", req.BattleType)
	return h.sendResponse(conn, msg.Header.MessageID, uint16(protocol.MsgCreateBattle), response)
}

// 辅助方法

func (h *GameHandler) unmarshalPayload(payload interface{}, target interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func (h *GameHandler) sendResponse(conn *connection.Connection, messageID uint32, messageType uint16, payload interface{}) error {
	// Convert PlayerID string to uint64
	var playerID uint64
	if conn.PlayerID != "" {
		if id, err := strconv.ParseUint(conn.PlayerID, 10, 64); err == nil {
			playerID = id
		}
	}

	response := &protocol.Message{
		Header: protocol.MessageHeader{
			Magic:       protocol.MessageMagic,
			MessageID:   messageID,
			MessageType: uint32(messageType),
			Flags:       protocol.FlagResponse,
			PlayerID:    playerID,
			Timestamp:   time.Now().Unix(),
			Sequence:    0,
		},
		Payload: payload,
	}

	return conn.SendMessage(response)
}

func (h *GameHandler) sendErrorResponse(conn *connection.Connection, messageID uint32, errorCode int, message string) error {
	// Convert PlayerID string to uint64
	var playerID uint64
	if conn.PlayerID != "" {
		if id, err := strconv.ParseUint(conn.PlayerID, 10, 64); err == nil {
			playerID = id
		}
	}

	errorMsg := &protocol.ErrorResponse{
		BaseResponse: protocol.NewBaseResponse(false, message),
		ErrorCode:    errorCode,
		ErrorType:    "ERROR",
	}

	response := &protocol.Message{
		Header: protocol.MessageHeader{
			Magic:       protocol.MessageMagic,
			MessageID:   messageID,
			MessageType: uint32(protocol.MsgError),
			Flags:       protocol.FlagResponse | protocol.FlagError,
			PlayerID:    playerID,
			Timestamp:   time.Now().Unix(),
		},
		Payload: errorMsg,
	}

	return conn.SendMessage(response)
}

func (h *GameHandler) broadcastPlayerOnline(player *playerQuery.PlayerDTO) {
	// 构造玩家上线广播消息
	broadcastMsg := &protocol.Message{
		Header: protocol.MessageHeader{
			Magic:       protocol.MessageMagic,
			MessageType: uint32(protocol.MsgPlayerStatus),
			Flags:       protocol.FlagBroadcast,
			PlayerID:    0, // 广播消息
			Timestamp:   time.Now().Unix(),
		},
		Payload: map[string]interface{}{
			"event":       "player_online",
			"player_id":   player.ID,
			"player_name": player.Name,
			"timestamp":   time.Now().Unix(),
		},
	}

	h.connMgr.BroadcastMessage(broadcastMsg)
}

func (h *GameHandler) broadcastPlayerMove(playerID string, x, y float64) {
	// 构造玩家移动广播消息
	broadcastMsg := &protocol.Message{
		Header: protocol.MessageHeader{
			Magic:       protocol.MessageMagic,
			MessageType: uint32(protocol.MsgPlayerMove),
			Flags:       protocol.FlagBroadcast,
			PlayerID:    0, // 广播消息
			Timestamp:   time.Now().Unix(),
		},
		Payload: map[string]interface{}{
			"player_id": playerID,
			"x":         x,
			"y":         y,
			"timestamp": time.Now().Unix(),
		},
	}

	h.connMgr.BroadcastMessage(broadcastMsg)
}

// 其他处理方法的占位符实现

func (h *GameHandler) handlePlayerLogout(ctx context.Context, conn *connection.Connection, msg *protocol.Message) error {
	// TODO: 实现玩家登出处理
	return h.sendResponse(conn, msg.Header.MessageID, uint16(protocol.MsgPlayerLogout), protocol.NewBaseResponse(true, "Logout successful"))
}

func (h *GameHandler) handlePlayerInfo(ctx context.Context, conn *connection.Connection, msg *protocol.Message) error {
	// TODO: 实现玩家信息查询
	return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeUnknown, "Not implemented")
}

func (h *GameHandler) handlePlayerCreate(ctx context.Context, conn *connection.Connection, msg *protocol.Message) error {
	// TODO: 实现玩家创建
	return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeUnknown, "Not implemented")
}

func (h *GameHandler) handleJoinBattle(ctx context.Context, conn *connection.Connection, msg *protocol.Message) error {
	// TODO: 实现加入战斗
	return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeUnknown, "Not implemented")
}

func (h *GameHandler) handleStartBattle(ctx context.Context, conn *connection.Connection, msg *protocol.Message) error {
	// TODO: 实现开始战斗
	return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeUnknown, "Not implemented")
}

func (h *GameHandler) handleBattleAction(ctx context.Context, conn *connection.Connection, msg *protocol.Message) error {
	// TODO: 实现战斗行动
	return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeUnknown, "Not implemented")
}

func (h *GameHandler) handleLeaveBattle(ctx context.Context, conn *connection.Connection, msg *protocol.Message) error {
	// TODO: 实现离开战斗
	return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeUnknown, "Not implemented")
}

func (h *GameHandler) handleGetPlayerInfo(ctx context.Context, conn *connection.Connection, msg *protocol.Message) error {
	// TODO: 实现获取玩家信息
	return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeUnknown, "Not implemented")
}

func (h *GameHandler) handleGetOnlinePlayers(ctx context.Context, conn *connection.Connection, msg *protocol.Message) error {
	// TODO: 实现获取在线玩家列表
	return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeUnknown, "Not implemented")
}

func (h *GameHandler) handleGetBattleInfo(ctx context.Context, conn *connection.Connection, msg *protocol.Message) error {
	// TODO: 实现获取战斗信息
	return h.sendErrorResponse(conn, msg.Header.MessageID, protocol.ErrCodeUnknown, "Not implemented")
}
