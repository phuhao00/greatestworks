package tcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	
	"github.com/netcore-go/netcore"
	"greatestworks/aop/logger"
	"greatestworks/application/services"
	"greatestworks/internal/infrastructure/network"
)

// NPCHandler NPC TCP处理器
type NPCHandler struct {
	npcService *services.NPCService
	logger     logger.Logger
}

// NPCRequest NPC请求
type NPCRequest struct {
	Action string                  `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

// NPCResponse NPC响应
type NPCResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// 消息类型常量
const (
	// NPC交互
	MsgTypeNPCList      uint32 = 3001
	MsgTypeNPCInfo      uint32 = 3002
	MsgTypeNPCInteract  uint32 = 3003
	MsgTypeNPCDialogue  uint32 = 3004
	MsgTypeNPCTrade     uint32 = 3005
	MsgTypeNPCQuest     uint32 = 3006
	MsgTypeNPCShop      uint32 = 3007
	MsgTypeNPCBattle    uint32 = 3008
	MsgTypeNPCStatus    uint32 = 3009
	MsgTypeNPCMove      uint32 = 3010
)

// NewNPCHandler 创建NPC处理器
func NewNPCHandler(npcService *services.NPCService, logger logger.Logger) *NPCHandler {
	return &NPCHandler{
		npcService: npcService,
		logger:     logger,
	}
}

// RegisterHandlers 注册处理器
func (h *NPCHandler) RegisterHandlers(server network.Server) error {
	// 注册NPC相关处理器
	if err := server.RegisterHandler(&NPCListHandler{h}); err != nil {
		return fmt.Errorf("failed to register NPC list handler: %w", err)
	}
	
	if err := server.RegisterHandler(&NPCInfoHandler{h}); err != nil {
		return fmt.Errorf("failed to register NPC info handler: %w", err)
	}
	
	if err := server.RegisterHandler(&NPCInteractHandler{h}); err != nil {
		return fmt.Errorf("failed to register NPC interact handler: %w", err)
	}
	
	if err := server.RegisterHandler(&NPCDialogueHandler{h}); err != nil {
		return fmt.Errorf("failed to register NPC dialogue handler: %w", err)
	}
	
	if err := server.RegisterHandler(&NPCTradeHandler{h}); err != nil {
		return fmt.Errorf("failed to register NPC trade handler: %w", err)
	}
	
	if err := server.RegisterHandler(&NPCQuestHandler{h}); err != nil {
		return fmt.Errorf("failed to register NPC quest handler: %w", err)
	}
	
	if err := server.RegisterHandler(&NPCShopHandler{h}); err != nil {
		return fmt.Errorf("failed to register NPC shop handler: %w", err)
	}
	
	if err := server.RegisterHandler(&NPCBattleHandler{h}); err != nil {
		return fmt.Errorf("failed to register NPC battle handler: %w", err)
	}
	
	if err := server.RegisterHandler(&NPCStatusHandler{h}); err != nil {
		return fmt.Errorf("failed to register NPC status handler: %w", err)
	}
	
	if err := server.RegisterHandler(&NPCMoveHandler{h}); err != nil {
		return fmt.Errorf("failed to register NPC move handler: %w", err)
	}
	
	h.logger.Info("NPC handlers registered successfully")
	return nil
}

// NPC列表处理器
type NPCListHandler struct {
	*NPCHandler
}

func (h *NPCListHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req NPCRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal NPC list request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}
	
	// 提取区域ID
	regionID, ok := req.Data["region_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing region_id")
	}
	
	// 调用服务层获取NPC列表
	npcs, err := h.npcService.GetNPCsByRegion(ctx, regionID)
	if err != nil {
		h.logger.Error("Failed to get NPC list", "error", err, "region_id", regionID)
		return h.sendErrorResponse(conn, "Failed to get NPC list: "+err.Error())
	}
	
	// 发送成功响应
	response := NPCResponse{
		Success: true,
		Message: "NPC list retrieved successfully",
		Data:    npcs,
	}
	
	return h.sendResponse(conn, MsgTypeNPCList, response)
}

func (h *NPCListHandler) GetMessageType() uint32 {
	return MsgTypeNPCList
}

func (h *NPCListHandler) GetHandlerName() string {
	return "NPCListHandler"
}

// NPC信息处理器
type NPCInfoHandler struct {
	*NPCHandler
}

func (h *NPCInfoHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req NPCRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal NPC info request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}
	
	// 提取NPC ID
	npcID, ok := req.Data["npc_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing npc_id")
	}
	
	// 调用服务层获取NPC信息
	npc, err := h.npcService.GetNPCInfo(ctx, npcID)
	if err != nil {
		h.logger.Error("Failed to get NPC info", "error", err, "npc_id", npcID)
		return h.sendErrorResponse(conn, "Failed to get NPC info: "+err.Error())
	}
	
	// 发送成功响应
	response := NPCResponse{
		Success: true,
		Message: "NPC info retrieved successfully",
		Data:    npc,
	}
	
	return h.sendResponse(conn, MsgTypeNPCInfo, response)
}

func (h *NPCInfoHandler) GetMessageType() uint32 {
	return MsgTypeNPCInfo
}

func (h *NPCInfoHandler) GetHandlerName() string {
	return "NPCInfoHandler"
}

// NPC交互处理器
type NPCInteractHandler struct {
	*NPCHandler
}

func (h *NPCInteractHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req NPCRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal NPC interact request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}
	
	// 提取参数
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}
	
	npcID, ok := req.Data["npc_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing npc_id")
	}
	
	interactionType, ok := req.Data["interaction_type"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing interaction_type")
	}
	
	// 调用服务层进行NPC交互
	result, err := h.npcService.InteractWithNPC(ctx, playerID, npcID, interactionType)
	if err != nil {
		h.logger.Error("Failed to interact with NPC", "error", err, "player_id", playerID, "npc_id", npcID, "interaction_type", interactionType)
		return h.sendErrorResponse(conn, "Failed to interact with NPC: "+err.Error())
	}
	
	// 发送成功响应
	response := NPCResponse{
		Success: true,
		Message: "NPC interaction completed successfully",
		Data:    result,
	}
	
	h.logger.Info("NPC interaction completed", "player_id", playerID, "npc_id", npcID, "interaction_type", interactionType)
	return h.sendResponse(conn, MsgTypeNPCInteract, response)
}

func (h *NPCInteractHandler) GetMessageType() uint32 {
	return MsgTypeNPCInteract
}

func (h *NPCInteractHandler) GetHandlerName() string {
	return "NPCInteractHandler"
}

// NPC对话处理器
type NPCDialogueHandler struct {
	*NPCHandler
}

func (h *NPCDialogueHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req NPCRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal NPC dialogue request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}
	
	// 提取参数
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}
	
	npcID, ok := req.Data["npc_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing npc_id")
	}
	
	dialogueID, ok := req.Data["dialogue_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing dialogue_id")
	}
	
	// 调用服务层进行对话
	dialogue, err := h.npcService.StartDialogue(ctx, playerID, npcID, dialogueID)
	if err != nil {
		h.logger.Error("Failed to start dialogue", "error", err, "player_id", playerID, "npc_id", npcID, "dialogue_id", dialogueID)
		return h.sendErrorResponse(conn, "Failed to start dialogue: "+err.Error())
	}
	
	// 发送成功响应
	response := NPCResponse{
		Success: true,
		Message: "Dialogue started successfully",
		Data:    dialogue,
	}
	
	h.logger.Info("Dialogue started", "player_id", playerID, "npc_id", npcID, "dialogue_id", dialogueID)
	return h.sendResponse(conn, MsgTypeNPCDialogue, response)
}

func (h *NPCDialogueHandler) GetMessageType() uint32 {
	return MsgTypeNPCDialogue
}

func (h *NPCDialogueHandler) GetHandlerName() string {
	return "NPCDialogueHandler"
}

// NPC交易处理器
type NPCTradeHandler struct {
	*NPCHandler
}

func (h *NPCTradeHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req NPCRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal NPC trade request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}
	
	// 提取参数
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}
	
	npcID, ok := req.Data["npc_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing npc_id")
	}
	
	tradeData, ok := req.Data["trade_data"].(map[string]interface{})
	if !ok {
		return h.sendErrorResponse(conn, "Missing trade_data")
	}
	
	// 调用服务层进行交易
	result, err := h.npcService.TradeWithNPC(ctx, playerID, npcID, tradeData)
	if err != nil {
		h.logger.Error("Failed to trade with NPC", "error", err, "player_id", playerID, "npc_id", npcID)
		return h.sendErrorResponse(conn, "Failed to trade with NPC: "+err.Error())
	}
	
	// 发送成功响应
	response := NPCResponse{
		Success: true,
		Message: "Trade completed successfully",
		Data:    result,
	}
	
	h.logger.Info("Trade completed", "player_id", playerID, "npc_id", npcID)
	return h.sendResponse(conn, MsgTypeNPCTrade, response)
}

func (h *NPCTradeHandler) GetMessageType() uint32 {
	return MsgTypeNPCTrade
}

func (h *NPCTradeHandler) GetHandlerName() string {
	return "NPCTradeHandler"
}

// NPC任务处理器
type NPCQuestHandler struct {
	*NPCHandler
}

func (h *NPCQuestHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req NPCRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal NPC quest request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}
	
	// 提取参数
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}
	
	npcID, ok := req.Data["npc_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing npc_id")
	}
	
	questAction, ok := req.Data["quest_action"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing quest_action")
	}
	
	// 调用服务层处理任务
	result, err := h.npcService.HandleQuest(ctx, playerID, npcID, questAction)
	if err != nil {
		h.logger.Error("Failed to handle quest", "error", err, "player_id", playerID, "npc_id", npcID, "quest_action", questAction)
		return h.sendErrorResponse(conn, "Failed to handle quest: "+err.Error())
	}
	
	// 发送成功响应
	response := NPCResponse{
		Success: true,
		Message: "Quest handled successfully",
		Data:    result,
	}
	
	h.logger.Info("Quest handled", "player_id", playerID, "npc_id", npcID, "quest_action", questAction)
	return h.sendResponse(conn, MsgTypeNPCQuest, response)
}

func (h *NPCQuestHandler) GetMessageType() uint32 {
	return MsgTypeNPCQuest
}

func (h *NPCQuestHandler) GetHandlerName() string {
	return "NPCQuestHandler"
}

// NPC商店处理器
type NPCShopHandler struct {
	*NPCHandler
}

func (h *NPCShopHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req NPCRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal NPC shop request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}
	
	// 提取参数
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}
	
	npcID, ok := req.Data["npc_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing npc_id")
	}
	
	shopAction, ok := req.Data["shop_action"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing shop_action")
	}
	
	// 调用服务层处理商店操作
	result, err := h.npcService.HandleShop(ctx, playerID, npcID, shopAction, req.Data)
	if err != nil {
		h.logger.Error("Failed to handle shop", "error", err, "player_id", playerID, "npc_id", npcID, "shop_action", shopAction)
		return h.sendErrorResponse(conn, "Failed to handle shop: "+err.Error())
	}
	
	// 发送成功响应
	response := NPCResponse{
		Success: true,
		Message: "Shop operation completed successfully",
		Data:    result,
	}
	
	h.logger.Info("Shop operation completed", "player_id", playerID, "npc_id", npcID, "shop_action", shopAction)
	return h.sendResponse(conn, MsgTypeNPCShop, response)
}

func (h *NPCShopHandler) GetMessageType() uint32 {
	return MsgTypeNPCShop
}

func (h *NPCShopHandler) GetHandlerName() string {
	return "NPCShopHandler"
}

// NPC战斗处理器
type NPCBattleHandler struct {
	*NPCHandler
}

func (h *NPCBattleHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req NPCRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal NPC battle request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}
	
	// 提取参数
	playerID, ok := req.Data["player_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing player_id")
	}
	
	npcID, ok := req.Data["npc_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing npc_id")
	}
	
	battleAction, ok := req.Data["battle_action"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing battle_action")
	}
	
	// 调用服务层处理战斗
	result, err := h.npcService.BattleWithNPC(ctx, playerID, npcID, battleAction)
	if err != nil {
		h.logger.Error("Failed to battle with NPC", "error", err, "player_id", playerID, "npc_id", npcID, "battle_action", battleAction)
		return h.sendErrorResponse(conn, "Failed to battle with NPC: "+err.Error())
	}
	
	// 发送成功响应
	response := NPCResponse{
		Success: true,
		Message: "Battle action completed successfully",
		Data:    result,
	}
	
	h.logger.Info("Battle action completed", "player_id", playerID, "npc_id", npcID, "battle_action", battleAction)
	return h.sendResponse(conn, MsgTypeNPCBattle, response)
}

func (h *NPCBattleHandler) GetMessageType() uint32 {
	return MsgTypeNPCBattle
}

func (h *NPCBattleHandler) GetHandlerName() string {
	return "NPCBattleHandler"
}

// NPC状态处理器
type NPCStatusHandler struct {
	*NPCHandler
}

func (h *NPCStatusHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req NPCRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal NPC status request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}
	
	// 提取NPC ID
	npcID, ok := req.Data["npc_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing npc_id")
	}
	
	// 调用服务层获取NPC状态
	status, err := h.npcService.GetNPCStatus(ctx, npcID)
	if err != nil {
		h.logger.Error("Failed to get NPC status", "error", err, "npc_id", npcID)
		return h.sendErrorResponse(conn, "Failed to get NPC status: "+err.Error())
	}
	
	// 发送成功响应
	response := NPCResponse{
		Success: true,
		Message: "NPC status retrieved successfully",
		Data:    status,
	}
	
	return h.sendResponse(conn, MsgTypeNPCStatus, response)
}

func (h *NPCStatusHandler) GetMessageType() uint32 {
	return MsgTypeNPCStatus
}

func (h *NPCStatusHandler) GetHandlerName() string {
	return "NPCStatusHandler"
}

// NPC移动处理器
type NPCMoveHandler struct {
	*NPCHandler
}

func (h *NPCMoveHandler) Handle(ctx context.Context, conn *netcore.Connection, packet *netcore.Packet) error {
	var req NPCRequest
	if err := json.Unmarshal(packet.GetData(), &req); err != nil {
		h.logger.Error("Failed to unmarshal NPC move request", "error", err)
		return h.sendErrorResponse(conn, "Invalid request format")
	}
	
	// 提取参数
	npcID, ok := req.Data["npc_id"].(string)
	if !ok {
		return h.sendErrorResponse(conn, "Missing npc_id")
	}
	
	targetPosition, ok := req.Data["target_position"].(map[string]interface{})
	if !ok {
		return h.sendErrorResponse(conn, "Missing target_position")
	}
	
	// 调用服务层移动NPC
	result, err := h.npcService.MoveNPC(ctx, npcID, targetPosition)
	if err != nil {
		h.logger.Error("Failed to move NPC", "error", err, "npc_id", npcID)
		return h.sendErrorResponse(conn, "Failed to move NPC: "+err.Error())
	}
	
	// 发送成功响应
	response := NPCResponse{
		Success: true,
		Message: "NPC moved successfully",
		Data:    result,
	}
	
	h.logger.Info("NPC moved successfully", "npc_id", npcID)
	return h.sendResponse(conn, MsgTypeNPCMove, response)
}

func (h *NPCMoveHandler) GetMessageType() uint32 {
	return MsgTypeNPCMove
}

func (h *NPCMoveHandler) GetHandlerName() string {
	return "NPCMoveHandler"
}

// 辅助方法

// sendResponse 发送响应
func (h *NPCHandler) sendResponse(conn *netcore.Connection, msgType uint32, response NPCResponse) error {
	data, err := json.Marshal(response)
	if err != nil {
		h.logger.Error("Failed to marshal response", "error", err)
		return err
	}
	
	packet := netcore.NewPacket(msgType, data)
	return conn.Send(packet)
}

// sendErrorResponse 发送错误响应
func (h *NPCHandler) sendErrorResponse(conn *netcore.Connection, errorMsg string) error {
	response := NPCResponse{
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