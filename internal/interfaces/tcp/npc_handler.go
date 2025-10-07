package tcp

import (
	"context"
	"fmt"

	"greatestworks/application/services"
	"greatestworks/internal/infrastructure/logger"
	"greatestworks/internal/infrastructure/network"
	"greatestworks/internal/interfaces/tcp/protocol"
	commonProto "greatestworks/internal/proto/common"
	"google.golang.org/protobuf/proto"
)

// NPCHandler NPC TCP处理器
type NPCHandler struct {
	npcService *services.NPCService
	logger     logger.Logger
}

// NPCRequest NPC请求
type NPCRequest struct {
	NPCID      string `json:"npc_id"`
	Action     string `json:"action"`
	PlayerID   string `json:"player_id"`
	DialogueID string `json:"dialogue_id,omitempty"`
	Choice     int    `json:"choice,omitempty"`
}

// NPCResponse NPC响应
type NPCResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewNPCHandler 创建NPC处理器
func NewNPCHandler(npcService *services.NPCService, logger logger.Logger) *NPCHandler {
	return &NPCHandler{
		npcService: npcService,
		logger:     logger,
	}
}

// RegisterHandlers 注册处理器
func (h *NPCHandler) RegisterHandlers(server *network.Server) {
	// 注册NPC相关消息处理器
	server.RegisterHandler(network.MsgNPCDialogue, h.handleNPCDialogue)
	server.RegisterHandler(network.MsgNPCInteraction, h.handleNPCInteraction)
	server.RegisterHandler(network.MsgNPCQuest, h.handleNPCQuest)
}

// handleNPCDialogue 处理NPC对话
func (h *NPCHandler) handleNPCDialogue(ctx context.Context, conn network.Connection, msg protocol.Message) error {
	var req commonProto.CommonRequest
	if err := proto.Unmarshal(msg.Payload, &req); err != nil {
		return h.sendError(conn, "解析请求失败", err)
	}

	// 开始对话
	dialogue, err := h.npcService.StartDialogue(ctx, req.NPCID, req.PlayerID)
	if err != nil {
		return h.sendError(conn, "开始对话失败", err)
	}

	response := NPCResponse{
		Success: true,
		Message: "对话开始",
		Data:    dialogue,
	}

	return h.sendResponse(conn, response)
}

// handleNPCInteraction 处理NPC交互
func (h *NPCHandler) handleNPCInteraction(ctx context.Context, conn network.Connection, msg protocol.Message) error {
	var req commonProto.CommonRequest
	if err := proto.Unmarshal(msg.Payload, &req); err != nil {
		return h.sendError(conn, "解析请求失败", err)
	}

	// 处理对话选择
	if req.DialogueID != "" && req.Choice >= 0 {
		result, err := h.npcService.ProcessDialogueChoice(ctx, req.DialogueID, req.Choice, req.PlayerID)
		if err != nil {
			return h.sendError(conn, "处理对话选择失败", err)
		}

		response := NPCResponse{
			Success: true,
			Message: "对话选择处理成功",
			Data:    result,
		}

		return h.sendResponse(conn, response)
	}

	// 结束对话
	if req.Action == "end_dialogue" {
		err := h.npcService.EndDialogue(ctx, req.DialogueID, req.PlayerID)
		if err != nil {
			return h.sendError(conn, "结束对话失败", err)
		}

		response := NPCResponse{
			Success: true,
			Message: "对话结束",
		}

		return h.sendResponse(conn, response)
	}

	return h.sendError(conn, "未知的交互类型", fmt.Errorf("unknown action: %s", req.Action))
}

// handleNPCQuest 处理NPC任务
func (h *NPCHandler) handleNPCQuest(ctx context.Context, conn network.Connection, msg protocol.Message) error {
	var req commonProto.CommonRequest
	if err := proto.Unmarshal(msg.Payload, &req); err != nil {
		return h.sendError(conn, "解析请求失败", err)
	}

	// 获取NPC任务
	quests, err := h.npcService.GetNPCQuests(ctx, req.NPCID, req.PlayerID)
	if err != nil {
		return h.sendError(conn, "获取任务失败", err)
	}

	response := NPCResponse{
		Success: true,
		Message: "获取任务成功",
		Data:    quests,
	}

	return h.sendResponse(conn, response)
}

// sendResponse 发送响应
func (h *NPCHandler) sendResponse(conn network.Connection, response NPCResponse) error {
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("序列化响应失败: %w", err)
	}

	packet := &network.Packet{
		Type: network.MsgNPCResponse,
		Data: data,
	}

	return conn.Send(packet)
}

// sendError 发送错误响应
func (h *NPCHandler) sendError(conn network.Connection, message string, err error) error {
	response := NPCResponse{
		Success: false,
		Message: fmt.Sprintf("%s: %v", message, err),
	}

	return h.sendResponse(conn, response)
}
