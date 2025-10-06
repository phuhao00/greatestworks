package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"greatestworks/application/services"
	"greatestworks/internal/infrastructure/logger"
)

// PlayerRPCService 玩家RPC服务
type PlayerRPCService struct {
	playerService *services.PlayerService
	logger        logger.Logger
}

// NewPlayerRPCService 创建玩家RPC服务
func NewPlayerRPCService(playerService *services.PlayerService, logger logger.Logger) *PlayerRPCService {
	return &PlayerRPCService{
		playerService: playerService,
		logger:        logger,
	}
}

// GetName 获取服务名称
func (s *PlayerRPCService) GetName() string {
	return "PlayerService"
}

// HandleRequest 处理请求
func (s *PlayerRPCService) HandleRequest(ctx context.Context, method string, data []byte) ([]byte, error) {
	switch method {
	case "CreatePlayer":
		return s.handleCreatePlayer(ctx, data)
	case "Login":
		return s.handleLogin(ctx, data)
	case "Logout":
		return s.handleLogout(ctx, data)
	case "GetPlayerInfo":
		return s.handleGetPlayerInfo(ctx, data)
	case "UpdatePlayer":
		return s.handleUpdatePlayer(ctx, data)
	case "MovePlayer":
		return s.handleMovePlayer(ctx, data)
	case "GetOnlinePlayers":
		return s.handleGetOnlinePlayers(ctx, data)
	default:
		return nil, fmt.Errorf("未知方法: %s", method)
	}
}

// handleCreatePlayer 处理创建玩家请求
func (s *PlayerRPCService) handleCreatePlayer(ctx context.Context, data []byte) ([]byte, error) {
	var req services.CreatePlayerCommand
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	result, err := s.playerService.CreatePlayer(ctx, &req)
	if err != nil {
		return nil, err
	}

	return json.Marshal(result)
}

// handleLogin 处理登录请求
func (s *PlayerRPCService) handleLogin(ctx context.Context, data []byte) ([]byte, error) {
	var req struct {
		PlayerID string `json:"player_id"`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	result, err := s.playerService.Login(ctx, req.PlayerID)
	if err != nil {
		return nil, err
	}

	return json.Marshal(result)
}

// handleLogout 处理登出请求
func (s *PlayerRPCService) handleLogout(ctx context.Context, data []byte) ([]byte, error) {
	var req struct {
		PlayerID string `json:"player_id"`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	err := s.playerService.Logout(ctx, req.PlayerID)
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]string{"status": "success"})
}

// handleGetPlayerInfo 处理获取玩家信息请求
func (s *PlayerRPCService) handleGetPlayerInfo(ctx context.Context, data []byte) ([]byte, error) {
	var req struct {
		PlayerID string `json:"player_id"`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	result, err := s.playerService.GetPlayerInfo(ctx, req.PlayerID)
	if err != nil {
		return nil, err
	}

	return json.Marshal(result)
}

// handleUpdatePlayer 处理更新玩家请求
func (s *PlayerRPCService) handleUpdatePlayer(ctx context.Context, data []byte) ([]byte, error) {
	var req struct {
		PlayerID string                 `json:"player_id"`
		Updates  map[string]interface{} `json:"updates"`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	err := s.playerService.UpdatePlayer(ctx, req.PlayerID, req.Updates)
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]string{"status": "success"})
}

// handleMovePlayer 处理移动玩家请求
func (s *PlayerRPCService) handleMovePlayer(ctx context.Context, data []byte) ([]byte, error) {
	var req struct {
		PlayerID string `json:"player_id"`
		Position struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
			Z float64 `json:"z"`
		} `json:"position"`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	// 这里需要根据实际的Position类型进行转换
	// position := player.Position{X: req.Position.X, Y: req.Position.Y, Z: req.Position.Z}
	// err := s.playerService.MovePlayer(ctx, req.PlayerID, position)
	// if err != nil {
	//     return nil, err
	// }

	return json.Marshal(map[string]string{"status": "success"})
}

// handleGetOnlinePlayers 处理获取在线玩家请求
func (s *PlayerRPCService) handleGetOnlinePlayers(ctx context.Context, data []byte) ([]byte, error) {
	var req struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	query := &services.GetOnlinePlayersQuery{
		Limit: req.Limit,
	}
	result, err := s.playerService.GetOnlinePlayers(ctx, query)
	if err != nil {
		return nil, err
	}

	return json.Marshal(result)
}
