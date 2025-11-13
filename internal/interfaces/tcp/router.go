package tcp

import (
	"encoding/json"
	"fmt"
	"sync"

	"greatestworks/internal/infrastructure/logging"
	"greatestworks/internal/interfaces/tcp/connection"
	"greatestworks/internal/interfaces/tcp/handlers"
	"greatestworks/internal/interfaces/tcp/protocol"
)

// 使用其他地方定义的Logger接口

// MessageHandler 消息处理器接口
type MessageHandler interface {
	HandleMessage(session *connection.Session, msg *protocol.Message) error
}

// Router TCP消息路由器
type Router struct {
	handlers map[uint16]MessageHandler
	mutex    sync.RWMutex
	logger   logging.Logger
}

// NewRouter 创建新的路由器
func NewRouter(logger logging.Logger) *Router {
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
	r.logger.Info("Message handler registered", logging.Fields{
		"message_type": messageType,
	})
}

// RegisterGameHandler 注册游戏处理器的所有消息类型
func (r *Router) RegisterGameHandler(handler *handlers.GameHandler) {
	// 系统消息
	r.RegisterHandler(uint16(protocol.MsgHeartbeat), handler)
	r.RegisterHandler(uint16(protocol.MsgAuth), handler)
	r.RegisterHandler(uint16(protocol.MsgError), handler)

	// 玩家相关消息
	r.RegisterHandler(uint16(protocol.MsgPlayerLogin), handler)
	r.RegisterHandler(uint16(protocol.MsgPlayerLogout), handler)
	r.RegisterHandler(uint16(protocol.MsgPlayerMove), handler)
	r.RegisterHandler(uint16(protocol.MsgPlayerStatus), handler)
	r.RegisterHandler(uint16(protocol.MsgBattleSkill), handler)
	r.RegisterHandler(uint16(protocol.MsgPlayerInfo), handler)
	r.RegisterHandler(uint16(protocol.MsgPlayerCreate), handler)
	r.RegisterHandler(uint16(protocol.MsgPlayerStatus), handler)
	r.RegisterHandler(uint16(protocol.MsgPlayerStats), handler)
	//r.RegisterHandler(uint16(protocol.MsgPlayerInventory), handler)
	//r.RegisterHandler(uint16(protocol.MsgPlayerSkills), handler)
	//r.RegisterHandler(uint16(protocol.MsgPlayerQuests), handler)

	// 战斗相关消息
	r.RegisterHandler(uint16(protocol.MsgCreateBattle), handler)
	r.RegisterHandler(uint16(protocol.MsgJoinBattle), handler)
	r.RegisterHandler(uint16(protocol.MsgStartBattle), handler)
	r.RegisterHandler(uint16(protocol.MsgBattleAction), handler)
	r.RegisterHandler(uint16(protocol.MsgBattleSkill), handler)
	r.RegisterHandler(uint16(protocol.MsgLeaveBattle), handler)
	r.RegisterHandler(uint16(protocol.MsgBattleStatus), handler)
	r.RegisterHandler(uint16(protocol.MsgBattleResult), handler)

	// 宠物相关消息
	r.RegisterHandler(uint16(protocol.MsgPetSummon), handler)
	r.RegisterHandler(uint16(protocol.MsgPetDismiss), handler)
	r.RegisterHandler(uint16(protocol.MsgPetAction), handler)
	r.RegisterHandler(uint16(protocol.MsgPetStatus), handler)
	r.RegisterHandler(uint16(protocol.MsgPetTrain), handler)
	//r.RegisterHandler(uint16(protocol.MsgPetEvolve), handler)

	// 建筑相关消息
	r.RegisterHandler(uint16(protocol.MsgBuildingCreate), handler)
	r.RegisterHandler(uint16(protocol.MsgBuildingUpgrade), handler)
	r.RegisterHandler(uint16(protocol.MsgBuildingDestroy), handler)
	r.RegisterHandler(uint16(protocol.MsgBuildingStatus), handler)
	//r.RegisterHandler(uint16(protocol.MsgBuildingList), handler)

	// 社交相关消息
	// 聊天与社交
	r.RegisterHandler(uint16(protocol.MsgChatMessage), handler)
	//r.RegisterHandler(uint16(protocol.MsgFriendAdd), handler)
	r.RegisterHandler(uint16(protocol.MsgFriendRemove), handler)
	r.RegisterHandler(uint16(protocol.MsgFriendList), handler)
	r.RegisterHandler(uint16(protocol.MsgGuildJoin), handler)
	r.RegisterHandler(uint16(protocol.MsgGuildLeave), handler)
	r.RegisterHandler(uint16(protocol.MsgGuildInfo), handler)
	r.RegisterHandler(uint16(protocol.MsgTeamCreate), handler)
	r.RegisterHandler(uint16(protocol.MsgTeamJoin), handler)
	r.RegisterHandler(uint16(protocol.MsgTeamLeave), handler)
	r.RegisterHandler(uint16(protocol.MsgTeamInfo), handler)

	// 物品相关消息
	r.RegisterHandler(uint16(protocol.MsgItemUse), handler)
	//r.RegisterHandler(uint16(protocol.MsgItemMove), handler)
	r.RegisterHandler(uint16(protocol.MsgItemDrop), handler)
	r.RegisterHandler(uint16(protocol.MsgItemPickup), handler)
	r.RegisterHandler(uint16(protocol.MsgItemTrade), handler)
	r.RegisterHandler(uint16(protocol.MsgItemCraft), handler)

	// 任务相关消息
	r.RegisterHandler(uint16(protocol.MsgQuestAccept), handler)
	r.RegisterHandler(uint16(protocol.MsgQuestComplete), handler)
	r.RegisterHandler(uint16(protocol.MsgQuestCancel), handler)
	r.RegisterHandler(uint16(protocol.MsgQuestProgress), handler)
	r.RegisterHandler(uint16(protocol.MsgQuestList), handler)

	// 查询相关消息
	r.RegisterHandler(uint16(protocol.MsgGetPlayerInfo), handler)
	r.RegisterHandler(uint16(protocol.MsgGetOnlinePlayers), handler)
	r.RegisterHandler(uint16(protocol.MsgGetBattleInfo), handler)
	r.RegisterHandler(uint16(protocol.MsgGetServerInfo), handler)
	//r.RegisterHandler(uint16(protocol.MsgGetWorldInfo), handler)

	r.logger.Info("Game handler registered with all message types")
}

// RouteMessage 路由消息到对应的处理器
func (r *Router) RouteMessage(session *connection.Session, msg *protocol.Message) error {
	r.mutex.RLock()
	handler, exists := r.handlers[uint16(msg.Header.MessageType)]
	r.mutex.RUnlock()

	if !exists {
		r.logger.Info("No handler found for message type", logging.Fields{
			"message_type": msg.Header.MessageType,
			"session_id":   session.ID,
			"user_id":      session.UserID,
		})
		return r.sendUnhandledMessageError(session, msg)
	}

	// 记录消息处理日志
	r.logger.Debug("Routing message", logging.Fields{
		"message_type": msg.Header.MessageType,
		"message_id":   msg.Header.MessageID,
		"session_id":   session.ID,
		"user_id":      session.UserID,
	})

	// 调用处理器
	err := handler.HandleMessage(session, msg)
	if err != nil {
		r.logger.Error("Message handler error", err, logging.Fields{
			"message_type": msg.Header.MessageType,
			"message_id":   msg.Header.MessageID,
			"session_id":   session.ID,
			"user_id":      session.UserID,
		})
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
	r.logger.Info("Message handler unregistered", logging.Fields{
		"message_type": messageType,
	})
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
func (r *Router) sendUnhandledMessageError(session *connection.Session, msg *protocol.Message) error {
	errorMsg := &protocol.ErrorResponse{
		BaseResponse: protocol.NewBaseResponse(false, fmt.Sprintf("Unhandled message type: %d", msg.Header.MessageType)),
		ErrorCode:    int(protocol.ErrCodeInvalidMessage),
		ErrorType:    "UNHANDLED_MESSAGE",
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

	// 序列化消息并发送
	data, err := json.Marshal(errorResponse)
	if err != nil {
		return fmt.Errorf("序列化错误消息失败: %w", err)
	}

	return session.Send(data)
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

	r.logger.Info("Router statistics", logging.Fields{
		"handler_count": handlerCount,
		"message_types": messageTypes,
	})
}
