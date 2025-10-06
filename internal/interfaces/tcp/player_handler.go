package tcp

import (
	"context"
	"encoding/json"
	"fmt"

	"greatestworks/aop/logger"
	"greatestworks/application/services"
	"greatestworks/internal/infrastructure/network"

	"github.com/netcore-go/netcore"
)

// PlayerHandler 玩家TCP处理器
type PlayerHandler struct {
	playerService *services.PlayerService
	hangupService *services.HangupService
	logger        logger.Logger
}

// PlayerRequest 玩家请求
type PlayerRequest struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

// PlayerResponse 玩家响应
type PlayerResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// 消息类型常量
const (
	MsgTypePlayerLogin  uint32 = 1001
	MsgTypePlayerLogout uint32 = 1002
	MsgTypePlayerInfo   uint32 = 1003
	MsgTypePlayerUpdate uint32 = 1004
	MsgTypePlayerStats  uint32 = 1005
	MsgTypeHangupStart  uint32 = 1101
	MsgTypeHangupStop   uint32 = 1102
	MsgTypeHangupStatus uint32 = 1103
	MsgTypeHangupReward uint32 = 1104
)

// NewPlayerHandler 创建玩家处理器
func NewPlayerHandler(playerService *services.PlayerService, hangupService *services.HangupService, logger logger.Logger) *PlayerHandler {
	return &PlayerHandler{
		playerService: playerService,
		hangupService: hangupService,
		logger:        logger,
	}
}

// RegisterHandlers 注册处理器
func (h *PlayerHandler) RegisterHandlers(server network.Server) error {
	// 注册玩家相关处理器
	if err := server.RegisterHandler(&PlayerLoginHandler{h}); err != nil {
		return fmt.Errorf("failed to register player login handler: %w", err)
	}

	if err := server.RegisterHandler(&PlayerLogoutHandler{h}); err != nil {
		return fmt.Errorf("failed to register player logout handler: %w", err)
	}

	if err := server.RegisterHandler(&PlayerInfoHandler{h}); err != nil {
		return fmt.Errorf("failed to register player info handler: %w", err)
	}

	if err := server.RegisterHandler(&PlayerUpdateHandler{h}); err != nil {
		return fmt.Errorf("failed to register player update handler: %w", err)
	}

	if err := server.RegisterHandler(&PlayerStatsHandler{h}); err != nil {
		return fmt.Errorf("failed to register player stats handler: %w", err)
	}

	// 注册挂机相关处理器
	if err := server.RegisterHandler(&HangupStartHandler{h}); err != nil {
		return fmt.Errorf("failed to register hangup start handler: %w", err)
	}

	if err := server.RegisterHandler(&HangupStopHandler{h}); err != nil {
		return fmt.Errorf("failed to register hangup stop handler: %w", err)
	}

	if err := server.RegisterHandler(&HangupStatusHandler{h}); err != nil {
		return fmt.Errorf("failed to register hangup status handler: %w", err)
	}

	if err := server.RegisterHandler(&HangupRewardHandler{h}); err != nil {
		return fmt.Errorf("failed to register hangup reward handler: %w", err)
	}

	h.logger.Info("Player handlers registered successfully")
	return nil
}

// 玩家登录处理器
type PlayerLoginHandler struct {
	*PlayerHandler
}

func (h *PlayerLoginHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req PlayerRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal player login request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取登录信息
	userID, ok := req.Data["user_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing user_id")
	}

	password, ok := req.Data["password"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing password")
	}

	// 调用服务层进行登录验证
	player, err := h.playerService.Login(ctx, userID, password)
	if err != nil {
		h.logger.Error("Player login failed", "error", err, "user_id", userID)
		return h.sendErrorResponse(conn, "Login failed: "+err.Error())
	}

	// 发送成功响应
	response := PlayerResponse{
		Success: true,
		Message: "Login successful",
		Data:    player,
	}

	h.logger.Info("Player logged in successfully", "user_id", userID, "player_id", player.ID)
	return h.sendResponse(conn, MsgTypePlayerLogin, response)
}

func (h *PlayerLoginHandler) GetMessageType() uint32 {
	return MsgTypePlayerLogin
}

func (h *PlayerLoginHandler) GetHandlerName() string {
	return "PlayerLoginHandler"
}

// 玩家登出处理器
type PlayerLogoutHandler struct {
	*PlayerHandler
}

func (h *PlayerLogoutHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req PlayerRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal player logout request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取玩家ID
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}

	// 调用服务层进行登出
	if err := h.playerService.Logout(ctx, playerID); err != nil {
		h.logger.Error("Player logout failed", "error", err, "player_id", playerID)
		return h.sendErrorResponse(conn, "Logout failed: "+err.Error())
	}

	// 发送成功响应
	response := PlayerResponse{
		Success: true,
		Message: "Logout successful",
	}

	h.logger.Info("Player logged out successfully", "player_id", playerID)
	return h.sendResponse(conn, MsgTypePlayerLogout, response)
}

func (h *PlayerLogoutHandler) GetMessageType() uint32 {
	return MsgTypePlayerLogout
}

func (h *PlayerLogoutHandler) GetHandlerName() string {
	return "PlayerLogoutHandler"
}

// 玩家信息处理器
type PlayerInfoHandler struct {
	*PlayerHandler
}

func (h *PlayerInfoHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req PlayerRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal player info request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取玩家ID
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}

	// 调用服务层获取玩家信息
	player, err := h.playerService.GetPlayerInfo(ctx, playerID)
	if err != nil {
		h.logger.Error("Failed to get player info", "error", err, "player_id", playerID)
		return h.sendErrorResponse(conn, "Failed to get player info: "+err.Error())
	}

	// 发送成功响应
	response := PlayerResponse{
		Success: true,
		Message: "Player info retrieved successfully",
		Data:    player,
	}

	return h.sendResponse(conn, MsgTypePlayerInfo, response)
}

func (h *PlayerInfoHandler) GetMessageType() uint32 {
	return MsgTypePlayerInfo
}

func (h *PlayerInfoHandler) GetHandlerName() string {
	return "PlayerInfoHandler"
}

// 玩家更新处理器
type PlayerUpdateHandler struct {
	*PlayerHandler
}

func (h *PlayerUpdateHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req PlayerRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal player update request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取玩家ID
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}

	// 提取更新数据
	updateData, ok := req.Data["update_data"].(map[string]interface{})
	if !ok {
		return h.sendErrorResponse(conn, "Missing update_data")
	}

	// 调用服务层更新玩家信息
	if err := h.playerService.UpdatePlayer(ctx, playerID, updateData); err != nil {
		h.logger.Error("Failed to update player", "error", err, "player_id", playerID)
		return h.sendErrorResponse(conn, "Failed to update player: "+err.Error())
	}

	// 发送成功响应
	response := PlayerResponse{
		Success: true,
		Message: "Player updated successfully",
	}

	h.logger.Info("Player updated successfully", "player_id", playerID)
	return h.sendResponse(conn, MsgTypePlayerUpdate, response)
}

func (h *PlayerUpdateHandler) GetMessageType() uint32 {
	return MsgTypePlayerUpdate
}

func (h *PlayerUpdateHandler) GetHandlerName() string {
	return "PlayerUpdateHandler"
}

// 玩家统计处理器
type PlayerStatsHandler struct {
	*PlayerHandler
}

func (h *PlayerStatsHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req PlayerRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal player stats request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取玩家ID
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}

	// 调用服务层获取玩家统计
	stats, err := h.playerService.GetPlayerStats(ctx, playerID)
	if err != nil {
		h.logger.Error("Failed to get player stats", "error", err, "player_id", playerID)
		return h.sendErrorResponse(conn, "Failed to get player stats: "+err.Error())
	}

	// 发送成功响应
	response := PlayerResponse{
		Success: true,
		Message: "Player stats retrieved successfully",
		Data:    stats,
	}

	return h.sendResponse(conn, MsgTypePlayerStats, response)
}

func (h *PlayerStatsHandler) GetMessageType() uint32 {
	return MsgTypePlayerStats
}

func (h *PlayerStatsHandler) GetHandlerName() string {
	return "PlayerStatsHandler"
}

// 挂机开始处理器
type HangupStartHandler struct {
	*PlayerHandler
}

func (h *HangupStartHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req PlayerRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal hangup start request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取玩家ID和挂机地点
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}

	locationID, ok := req.Data["location_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing location_id")
	}

	// 调用服务层开始挂机
	hangupInfo, err := h.hangupService.StartHangup(ctx, playerID, locationID)
	if err != nil {
		h.logger.Error("Failed to start hangup", "error", err, "player_id", playerID, "location_id", locationID)
		return h.sendErrorResponse(conn, "Failed to start hangup: "+err.Error())
	}

	// 发送成功响应
	response := PlayerResponse{
		Success: true,
		Message: "Hangup started successfully",
		Data:    hangupInfo,
	}

	h.logger.Info("Hangup started successfully", "player_id", playerID, "location_id", locationID)
	return h.sendResponse(conn, MsgTypeHangupStart, response)
}

func (h *HangupStartHandler) GetMessageType() uint32 {
	return MsgTypeHangupStart
}

func (h *HangupStartHandler) GetHandlerName() string {
	return "HangupStartHandler"
}

// 挂机停止处理器
type HangupStopHandler struct {
	*PlayerHandler
}

func (h *HangupStopHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req PlayerRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal hangup stop request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取玩家ID
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}

	// 调用服务层停止挂机
	rewards, err := h.hangupService.StopHangup(ctx, playerID)
	if err != nil {
		h.logger.Error("Failed to stop hangup", "error", err, "player_id", playerID)
		return h.sendErrorResponse(conn, "Failed to stop hangup: "+err.Error())
	}

	// 发送成功响应
	response := PlayerResponse{
		Success: true,
		Message: "Hangup stopped successfully",
		Data:    rewards,
	}

	h.logger.Info("Hangup stopped successfully", "player_id", playerID)
	return h.sendResponse(conn, MsgTypeHangupStop, response)
}

func (h *HangupStopHandler) GetMessageType() uint32 {
	return MsgTypeHangupStop
}

func (h *HangupStopHandler) GetHandlerName() string {
	return "HangupStopHandler"
}

// 挂机状态处理器
type HangupStatusHandler struct {
	*PlayerHandler
}

func (h *HangupStatusHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req PlayerRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal hangup status request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取玩家ID
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}

	// 调用服务层获取挂机状态
	status, err := h.hangupService.GetHangupStatus(ctx, playerID)
	if err != nil {
		h.logger.Error("Failed to get hangup status", "error", err, "player_id", playerID)
		return h.sendErrorResponse(conn, "Failed to get hangup status: "+err.Error())
	}

	// 发送成功响应
	response := PlayerResponse{
		Success: true,
		Message: "Hangup status retrieved successfully",
		Data:    status,
	}

	return h.sendResponse(conn, MsgTypeHangupStatus, response)
}

func (h *HangupStatusHandler) GetMessageType() uint32 {
	return MsgTypeHangupStatus
}

func (h *HangupStatusHandler) GetHandlerName() string {
	return "HangupStatusHandler"
}

// 挂机奖励处理器
type HangupRewardHandler struct {
	*PlayerHandler
}

func (h *HangupRewardHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req PlayerRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal hangup reward request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}

	// 提取玩家ID
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}

	// 调用服务层收集挂机奖励
	rewards, err := h.hangupService.CollectRewards(ctx, playerID)
	if err != nil {
		h.logger.Error("Failed to collect hangup rewards", "error", err, "player_id", playerID)
		return h.sendErrorResponse(conn, "Failed to collect hangup rewards: "+err.Error())
	}

	// 发送成功响应
	response := PlayerResponse{
		Success: true,
		Message: "Hangup rewards collected successfully",
		Data:    rewards,
	}

	h.logger.Info("Hangup rewards collected successfully", "player_id", playerID)
	return h.sendResponse(conn, MsgTypeHangupReward, response)
}

func (h *HangupRewardHandler) GetMessageType() uint32 {
	return MsgTypeHangupReward
}

func (h *HangupRewardHandler) GetHandlerName() string {
	return "HangupRewardHandler"
}

// 辅助方法

// sendResponse 发送响应
func (h *PlayerHandler) sendResponse(conn *netcore.Connection, msgType uint32, response PlayerResponse) error {
	data, err := json.Marshal(response)
	if err != nil {
		h.logger.Error("Failed to marshal response", "error", err)
		return err
	}

	packet := netcore.NewPacket(msgType, data)
	return conn.Send(packet)
}

// sendErrorResponse 发送错误响应
func (h *PlayerHandler) sendErrorResponse(conn *netcore.Connection, errorMsg string) error {
	response := PlayerResponse{
		Success: false,
		Message: "Request failed",
		Error:   errorMsg,
	}

	data, err := json.Marshal(response)
	if err != nil {
		h.logger.Error("Failed to marshal error response", "error", err)
		return err
	}

	// 使用通用错误消息类型
	packet := netcore.NewPacket(9999, data)
	return conn.Send(packet)
}
