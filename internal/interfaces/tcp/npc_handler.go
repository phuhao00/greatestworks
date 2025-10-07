package tcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"greatestworks/application/services"
	"greatestworks/internal/infrastructure/logger"
	"greatestworks/internal/infrastructure/network"
	"greatestworks/internal/network/session"
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
	server.RegisterHandler(0x3001, h.handleNPCDialogue)    // MsgNPCDialogue
	server.RegisterHandler(0x3002, h.handleNPCInteraction) // MsgNPCInteraction
	server.RegisterHandler(0x3003, h.handleNPCQuest)       // MsgNPCQuest
}

// handleNPCDialogue 处理NPC对话
func (h *NPCHandler) handleNPCDialogue(ctx context.Context, session session.Session, packet network.Packet) error {
	var req commonProto.CommonRequest
	if err := proto.Unmarshal(packet.GetData(), &req); err != nil {
		return h.sendError(session, "解析请求失败", err)
	}

	// 开始对话
	npcID := req.Metadata["npc_id"]
	playerID := req.Metadata["player_id"]
	dialogue, err := h.npcService.StartDialogue(ctx, npcID, playerID)
	if err != nil {
		return h.sendError(session, "开始对话失败", err)
	}

	response := NPCResponse{
		Success: true,
		Message: "对话开始",
		Data:    dialogue,
	}

	return h.sendResponse(session, response)
}

// handleNPCInteraction 处理NPC交互
func (h *NPCHandler) handleNPCInteraction(ctx context.Context, session session.Session, packet network.Packet) error {
	var req commonProto.CommonRequest
	if err := proto.Unmarshal(packet.GetData(), &req); err != nil {
		return h.sendError(session, "解析请求失败", err)
	}

	// 处理对话选择
	dialogueID := req.Metadata["dialogue_id"]
	choiceStr := req.Metadata["choice"]
	if dialogueID != "" && choiceStr != "" {
		_, err := strconv.Atoi(choiceStr)
		if err != nil {
			return h.sendError(session, "无效的选择", err)
		}
		// TODO: 实现ProcessDialogueChoice方法
		result := &struct{}{} // 占位符
		err = fmt.Errorf("ProcessDialogueChoice not implemented")
		if err != nil {
			return h.sendError(session, "处理对话选择失败", err)
		}

		response := NPCResponse{
			Success: true,
			Message: "对话选择处理成功",
			Data:    result,
		}

		return h.sendResponse(session, response)
	}

	// 结束对话
	action := req.Metadata["action"]
	if action == "end_dialogue" {
		dialogueID := req.Metadata["dialogue_id"]
		playerID := req.Metadata["player_id"]
		err := h.npcService.EndDialogue(ctx, dialogueID, playerID)
		if err != nil {
			return h.sendError(session, "结束对话失败", err)
		}

		response := NPCResponse{
			Success: true,
			Message: "对话结束",
		}

		return h.sendResponse(session, response)
	}

	return h.sendError(session, "未知的交互类型", fmt.Errorf("unknown action: %s", action))
}

// handleNPCQuest 处理NPC任务
func (h *NPCHandler) handleNPCQuest(ctx context.Context, session session.Session, packet network.Packet) error {
	var req commonProto.CommonRequest
	if err := proto.Unmarshal(packet.GetData(), &req); err != nil {
		return h.sendError(session, "解析请求失败", err)
	}

	// 获取NPC任务
	// TODO: 实现GetNPCQuests方法
	quests := []interface{}{} // 占位符
	err := fmt.Errorf("GetNPCQuests not implemented")
	if err != nil {
		return h.sendError(session, "获取任务失败", err)
	}

	response := NPCResponse{
		Success: true,
		Message: "获取任务成功",
		Data:    quests,
	}

	return h.sendResponse(session, response)
}

// sendResponse 发送响应
func (h *NPCHandler) sendResponse(session session.Session, response NPCResponse) error {
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("序列化响应失败: %w", err)
	}

	return session.Send(data)
}

// sendError 发送错误响应
func (h *NPCHandler) sendError(session session.Session, message string, err error) error {
	response := NPCResponse{
		Success: false,
		Message: fmt.Sprintf("%s: %v", message, err),
	}

	return h.sendResponse(session, response)
}
