package tcp

import (
	"fmt"

	"greatestworks/internal/infrastructure/logging"
	"greatestworks/internal/interfaces/tcp/connection"
	"greatestworks/internal/interfaces/tcp/protocol"
)

// NPCHandler NPC处理器
type NPCHandler struct {
	logger logging.Logger
}

// NewNPCHandler 创建NPC处理器
func NewNPCHandler(logger logging.Logger) *NPCHandler {
	return &NPCHandler{
		logger: logger,
	}
}

// HandleMessage 处理消息
func (h *NPCHandler) HandleMessage(session *connection.Session, message *protocol.Message) error {
	h.logger.Info("处理NPC消息", map[string]interface{}{
		"message_type": message.Header.MessageType,
		"player_id":    message.Header.PlayerID,
	})

	switch message.Header.MessageType {
	case uint32(protocol.MsgPlayerInfo):
		return h.handleNPCInteraction(session, message)
	case uint32(protocol.MsgQuestAccept):
		return h.handleNPCQuest(session, message)
	case uint32(protocol.MsgItemTrade):
		return h.handleNPCTrade(session, message)
	default:
		return fmt.Errorf("未知的NPC消息类型: %d", message.Header.MessageType)
	}
}

// handleNPCInteraction 处理NPC交互
func (h *NPCHandler) handleNPCInteraction(session *connection.Session, message *protocol.Message) error {
	h.logger.Info("处理NPC交互", map[string]interface{}{
		"player_id": message.Header.PlayerID,
	})

	// TODO: 实现NPC交互逻辑
	// 1. 验证NPC存在
	// 2. 检查交互条件
	// 3. 执行交互逻辑
	// 4. 发送响应

	return nil
}

// handleNPCQuest 处理NPC任务
func (h *NPCHandler) handleNPCQuest(session *connection.Session, message *protocol.Message) error {
	h.logger.Info("处理NPC任务", map[string]interface{}{
		"player_id": message.Header.PlayerID,
	})

	// TODO: 实现NPC任务逻辑
	// 1. 验证任务存在
	// 2. 检查任务条件
	// 3. 执行任务逻辑
	// 4. 发送响应

	return nil
}

// handleNPCTrade 处理NPC交易
func (h *NPCHandler) handleNPCTrade(session *connection.Session, message *protocol.Message) error {
	h.logger.Info("处理NPC交易", map[string]interface{}{
		"player_id": message.Header.PlayerID,
	})

	// TODO: 实现NPC交易逻辑
	// 1. 验证交易物品
	// 2. 检查交易条件
	// 3. 执行交易逻辑
	// 4. 发送响应

	return nil
}
