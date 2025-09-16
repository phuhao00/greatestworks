package tcp

import (
	"fmt"
	"sync"

	"greatestworks/internal/infrastructure/logger"
	"greatestworks/internal/interfaces/tcp/protocol"
	"greatestworks/internal/interfaces/tcp/connection"
	"greatestworks/internal/interfaces/tcp/handlers"
)

// MessageHandler 消息处理器接口
type MessageHandler interface {
	HandleMessage(conn *connection.Connection, msg *protocol.Message) error
}

// Router TCP消息路由器
type Router struct {
	handlers map[uint16]MessageHandler
	mutex    sync.RWMutex
	logger   logger.Logger
}

// NewRouter 创建新的路由器
func NewRouter(logger logger.Logger) *Router {
	return &Router{
		handlers: make(map[uint16]MessageHandler),
		logger:   logger,
	}
}

// RegisterHandler 注册消息处理器
func (r *Router) RegisterHandler(messageType uint16, handler MessageHandler) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.handlers[messageType] = handler
	r.logger.Info("Message handler registered", "message_type", messageType)
}

// RegisterGameHandler 注册游戏处理器的所有消息类型
func (r *Router) RegisterGameHandler(handler *handlers.GameHandler) {
	// 系统消息
	r.RegisterHandler(protocol.MsgHeartbeat, handler)
	r.RegisterHandler(protocol.MsgAuth, handler)
	r.RegisterHandler(protocol.MsgError, handler)

	// 玩家相关消息
	r.RegisterHandler(protocol.MsgPlayerLogin, handler)
	r.RegisterHandler(protocol.MsgPlayerLogout, handler)
	r.RegisterHandler(protocol.MsgPlayerMove, handler)
	r.RegisterHandler(protocol.MsgPlayerInfo, handler)
	r.RegisterHandler(protocol.MsgPlayerCreate, handler)
	r.RegisterHandler(protocol.MsgPlayerStatus, handler)
	r.RegisterHandler(protocol.MsgPlayerStats, handler)
	r.RegisterHandler(protocol.MsgPlayerInventory, handler)
	r.RegisterHandler(protocol.MsgPlayerSkills, handler)
	r.RegisterHandler(protocol.MsgPlayerQuests, handler)

	// 战斗相关消息
	r.RegisterHandler(protocol.MsgCreateBattle, handler)
	r.RegisterHandler(protocol.MsgJoinBattle, handler)
	r.RegisterHandler(protocol.MsgStartBattle, handler)
	r.RegisterHandler(protocol.MsgBattleAction, handler)
	r.RegisterHandler(protocol.MsgLeaveBattle, handler)
	r.RegisterHandler(protocol.MsgBattleStatus, handler)
	r.RegisterHandler(protocol.MsgBattleResult, handler)

	// 宠物相关消息
	r.RegisterHandler(protocol.MsgPetSummon, handler)
	r.RegisterHandler(protocol.MsgPetDismiss, handler)
	r.RegisterHandler(protocol.MsgPetAction, handler)
	r.RegisterHandler(protocol.MsgPetStatus, handler)
	r.RegisterHandler(protocol.MsgPetTrain, handler)
	r.RegisterHandler(protocol.MsgPetEvolve, handler)

	// 建筑相关消息
	r.RegisterHandler(protocol.MsgBuildingCreate, handler)
	r.RegisterHandler(protocol.MsgBuildingUpgrade, handler)
	r.RegisterHandler(protocol.MsgBuildingDestroy, handler)
	r.RegisterHandler(protocol.MsgBuildingStatus, handler)
	r.RegisterHandler(protocol.MsgBuildingList, handler)

	// 社交相关消息
	r.RegisterHandler(protocol.MsgChatSend, handler)
	r.RegisterHandler(protocol.MsgChatReceive, handler)
	r.RegisterHandler(protocol.MsgFriendAdd, handler)
	r.RegisterHandler(protocol.MsgFriendRemove, handler)
	r.RegisterHandler(protocol.MsgFriendList, handler)
	r.RegisterHandler(protocol.MsgGuildJoin, handler)
	r.RegisterHandler(protocol.MsgGuildLeave, handler)
	r.RegisterHandler(protocol.MsgGuildInfo, handler)
	r.RegisterHandler(protocol.MsgTeamCreate, handler)
	r.RegisterHandler(protocol.MsgTeamJoin, handler)
	r.RegisterHandler(protocol.MsgTeamLeave, handler)
	r.RegisterHandler(protocol.MsgTeamInfo, handler)

	// 物品相关消息
	r.RegisterHandler(protocol.MsgItemUse, handler)
	r.RegisterHandler(protocol.MsgItemMove, handler)
	r.RegisterHandler(protocol.MsgItemDrop, handler)
	r.RegisterHandler(protocol.MsgItemPickup, handler)
	r.RegisterHandler(protocol.MsgItemTrade, handler)
	r.RegisterHandler(protocol.MsgItemCraft, handler)

	// 任务相关消息
	r.RegisterHandler(protocol.MsgQuestAccept, handler)
	r.RegisterHandler(protocol.MsgQuestComplete, handler)
	r.RegisterHandler(protocol.MsgQuestCancel, handler)
	r.RegisterHandler(protocol.MsgQuestProgress, handler)
	r.RegisterHandler(protocol.MsgQuestList, handler)

	// 查询相关消息
	r.RegisterHandler(protocol.MsgGetPlayerInfo, handler)
	r.RegisterHandler(protocol.MsgGetOnlinePlayers, handler)
	r.RegisterHandler(protocol.MsgGetBattleInfo, handler)
	r.RegisterHandler(protocol.MsgGetServerInfo, handler)
	r.RegisterHandler(protocol.MsgGetWorldInfo, handler)

	r.logger.Info("Game handler registered with all message types")
}

// RouteMessage 路由消息到对应的处理器
func (r *Router) RouteMessage(conn *connection.Connection, msg *protocol.Message) error {
	r.mutex.RLock()
	handler, exists := r.handlers[msg.Header.MessageType]
	r.mutex.RUnlock()

	if !exists {
		r.logger.Warn("No handler found for message type", 
			"message_type", msg.Header.MessageType, 
			"conn_id", conn.ID,
			"player_id", conn.PlayerID)
		return r.sendUnhandledMessageError(conn, msg)
	}

	// 记录消息处理日志
	r.logger.Debug("Routing message", 
		"message_type", msg.Header.MessageType,
		"message_id", msg.Header.MessageID,
		"conn_id", conn.ID,
		"player_id", conn.PlayerID)

	// 调用处理器
	err := handler.HandleMessage(conn, msg)
	if err != nil {
		r.logger.Error("Message handler error", 
			"error", err,
			"message_type", msg.Header.MessageType,
			"message_id", msg.Header.MessageID,
			"conn_id", conn.ID,
			"player_id", conn.PlayerID)
		return err
	}

	return nil
}

// GetHandler 获取指定消息类型的处理器
func (r *Router) GetHandler(messageType uint16) (MessageHandler, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	handler, exists := r.handlers[messageType]
	return handler, exists
}

// UnregisterHandler 注销消息处理器
func (r *Router) UnregisterHandler(messageType uint16) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.handlers, messageType)
	r.logger.Info("Message handler unregistered", "message_type", messageType)
}

// GetRegisteredMessageTypes 获取所有已注册的消息类型
func (r *Router) GetRegisteredMessageTypes() []uint16 {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	types := make([]uint16, 0, len(r.handlers))
	for msgType := range r.handlers {
		types = append(types, msgType)
	}

	return types
}

// GetHandlerCount 获取已注册处理器数量
func (r *Router) GetHandlerCount() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.handlers)
}

// IsMessageTypeSupported 检查消息类型是否被支持
func (r *Router) IsMessageTypeSupported(messageType uint16) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, exists := r.handlers[messageType]
	return exists
}

// sendUnhandledMessageError 发送未处理消息错误
func (r *Router) sendUnhandledMessageError(conn *connection.Connection, msg *protocol.Message) error {
	errorMsg := &protocol.ErrorMessage{
		ErrorCode: protocol.ErrInvalidMessage,
		Message:   fmt.Sprintf("Unhandled message type: %d", msg.Header.MessageType),
		Timestamp: msg.Header.Timestamp,
	}

	errorResponse := &protocol.Message{
		Header: protocol.MessageHeader{
			Magic:       protocol.MessageMagic,
			MessageID:   msg.Header.MessageID,
			MessageType: protocol.MsgError,
			Flags:       protocol.FlagResponse | protocol.FlagError,
			PlayerID:    msg.Header.PlayerID,
			Timestamp:   msg.Header.Timestamp,
			Sequence:    0,
		},
		Payload: errorMsg,
	}

	return conn.SendMessage(errorResponse)
}

// ValidateMessage 验证消息格式
func (r *Router) ValidateMessage(msg *protocol.Message) error {
	if msg == nil {
		return fmt.Errorf("message is nil")
	}

	// 验证魔数
	if msg.Header.Magic != protocol.MessageMagic {
		return fmt.Errorf("invalid magic number: %d", msg.Header.Magic)
	}

	// 验证消息类型
	if msg.Header.MessageType == 0 {
		return fmt.Errorf("invalid message type: %d", msg.Header.MessageType)
	}

	// 验证消息ID
	if msg.Header.MessageID == 0 {
		return fmt.Errorf("invalid message ID: %d", msg.Header.MessageID)
	}

	// 验证时间戳
	if msg.Header.Timestamp <= 0 {
		return fmt.Errorf("invalid timestamp: %d", msg.Header.Timestamp)
	}

	return nil
}

// LogRouterStats 记录路由器统计信息
func (r *Router) LogRouterStats() {
	r.mutex.RLock()
	handlerCount := len(r.handlers)
	messageTypes := make([]uint16, 0, handlerCount)
	for msgType := range r.handlers {
		messageTypes = append(messageTypes, msgType)
	}
	r.mutex.RUnlock()

	r.logger.Info("Router statistics", 
		"handler_count", handlerCount,
		"message_types", messageTypes)
}