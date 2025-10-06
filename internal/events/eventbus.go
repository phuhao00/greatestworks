// Package events 事件系统和消息队列
// Author: MMO Server Team
// Created: 2024

package events

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	// "github.com/phuhao00/netcore-go/core" // 暂时注释掉缺失的包
)

// Logger 简单的日志接口
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// Event 事件接口
type Event interface {
	GetType() string
	GetData() interface{}
	GetTimestamp() time.Time
}

// Handler 事件处理器
type Handler func(ctx context.Context, event Event) error

// EventBus 事件总线
type EventBus struct {
	handlers map[string][]Handler
	mutex    sync.RWMutex
	logger   Logger
	nats     *nats.Conn
}

// BaseEvent 基础事件结构
type BaseEvent struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	UserID    string      `json:"user_id,omitempty"`
	SessionID string      `json:"session_id,omitempty"`
}

// GetType 获取事件类型
func (e *BaseEvent) GetType() string {
	return e.Type
}

// GetData 获取事件数据
func (e *BaseEvent) GetData() interface{} {
	return e.Data
}

// GetTimestamp 获取时间戳
func (e *BaseEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// NewEventBus 创建事件总线
func NewEventBus(logger Logger) *EventBus {
	return &EventBus{
		handlers: make(map[string][]Handler),
		logger:   logger,
	}
}

// ConnectNATS 连接到NATS
func (eb *EventBus) ConnectNATS(url string) error {
	nc, err := nats.Connect(url)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}

	eb.nats = nc
	eb.logger.Info("Connected to NATS", "url", url)
	return nil
}

// Subscribe 订阅事件
func (eb *EventBus) Subscribe(eventType string, handler Handler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
	eb.logger.Debug("Event handler subscribed", "type", eventType)
}

// Unsubscribe 取消订阅
func (eb *EventBus) Unsubscribe(eventType string, handler Handler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	handlers := eb.handlers[eventType]
	for i, h := range handlers {
		if reflect.ValueOf(h).Pointer() == reflect.ValueOf(handler).Pointer() {
			eb.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}

	eb.logger.Debug("Event handler unsubscribed", "type", eventType)
}

// Publish 发布事件
func (eb *EventBus) Publish(ctx context.Context, event Event) error {
	// 本地处理
	if err := eb.publishLocal(ctx, event); err != nil {
		eb.logger.Error("Local event publish failed", "error", err)
	}

	// 远程发布（通过NATS）
	if eb.nats != nil {
		if err := eb.publishRemote(event); err != nil {
			eb.logger.Error("Remote event publish failed", "error", err)
			return err
		}
	}

	return nil
}

// publishLocal 本地发布事件
func (eb *EventBus) publishLocal(ctx context.Context, event Event) error {
	eb.mutex.RLock()
	handlers := eb.handlers[event.GetType()]
	eb.mutex.RUnlock()

	for _, handler := range handlers {
		go func(h Handler) {
			if err := h(ctx, event); err != nil {
				eb.logger.Error("Event handler error", "type", event.GetType(), "error", err)
			}
		}(handler)
	}

	eb.logger.Debug("Event published locally", "type", event.GetType())
	return nil
}

// publishRemote 远程发布事件
func (eb *EventBus) publishRemote(event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	subject := fmt.Sprintf("events.%s", event.GetType())
	if err := eb.nats.Publish(subject, data); err != nil {
		return fmt.Errorf("failed to publish event to NATS: %w", err)
	}

	eb.logger.Debug("Event published remotely", "type", event.GetType(), "subject", subject)
	return nil
}

// SubscribeRemote 订阅远程事件
func (eb *EventBus) SubscribeRemote(eventType string, handler Handler) error {
	if eb.nats == nil {
		return fmt.Errorf("NATS connection not established")
	}

	subject := fmt.Sprintf("events.%s", eventType)
	_, err := eb.nats.Subscribe(subject, func(msg *nats.Msg) {
		var event BaseEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			eb.logger.Error("Failed to unmarshal remote event", "error", err)
			return
		}

		ctx := context.Background()
		if err := handler(ctx, &event); err != nil {
			eb.logger.Error("Remote event handler error", "type", eventType, "error", err)
		}
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to remote events: %w", err)
	}

	eb.logger.Info("Subscribed to remote events", "type", eventType, "subject", subject)
	return nil
}

// Close 关闭事件总线
func (eb *EventBus) Close() {
	if eb.nats != nil {
		eb.nats.Close()
		eb.logger.Info("NATS connection closed")
	}
}

// GameEvents 游戏事件类型常量
const (
	EventPlayerLogin   = "player.login"
	EventPlayerLogout  = "player.logout"
	EventPlayerMove    = "player.move"
	EventPlayerChat    = "player.chat"
	EventBattleStart   = "battle.start"
	EventBattleEnd     = "battle.end"
	EventSceneEnter    = "scene.enter"
	EventSceneLeave    = "scene.leave"
	EventActivityJoin  = "activity.join"
	EventActivityLeave = "activity.leave"
)

// PlayerLoginEvent 玩家登录事件
type PlayerLoginEvent struct {
	*BaseEvent
	PlayerID string `json:"player_id"`
	IP       string `json:"ip"`
}

// NewPlayerLoginEvent 创建玩家登录事件
func NewPlayerLoginEvent(playerID, ip string) *PlayerLoginEvent {
	return &PlayerLoginEvent{
		BaseEvent: &BaseEvent{
			Type:      EventPlayerLogin,
			Timestamp: time.Now(),
			UserID:    playerID,
		},
		PlayerID: playerID,
		IP:       ip,
	}
}

// PlayerMoveEvent 玩家移动事件
type PlayerMoveEvent struct {
	*BaseEvent
	PlayerID string  `json:"player_id"`
	SceneID  string  `json:"scene_id"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        float64 `json:"z"`
}

// NewPlayerMoveEvent 创建玩家移动事件
func NewPlayerMoveEvent(playerID, sceneID string, x, y, z float64) *PlayerMoveEvent {
	return &PlayerMoveEvent{
		BaseEvent: &BaseEvent{
			Type:      EventPlayerMove,
			Timestamp: time.Now(),
			UserID:    playerID,
		},
		PlayerID: playerID,
		SceneID:  sceneID,
		X:        x,
		Y:        y,
		Z:        z,
	}
}

// ChatEvent 聊天事件
type ChatEvent struct {
	*BaseEvent
	PlayerID string `json:"player_id"`
	Channel  string `json:"channel"`
	Message  string `json:"message"`
}

// NewChatEvent 创建聊天事件
func NewChatEvent(playerID, channel, message string) *ChatEvent {
	return &ChatEvent{
		BaseEvent: &BaseEvent{
			Type:      EventPlayerChat,
			Timestamp: time.Now(),
			UserID:    playerID,
		},
		PlayerID: playerID,
		Channel:  channel,
		Message:  message,
	}
}

// BattleStartEvent 战斗开始事件
type BattleStartEvent struct {
	*BaseEvent
	BattleID string   `json:"battle_id"`
	Players  []string `json:"players"`
}

// NewBattleStartEvent 创建战斗开始事件
func NewBattleStartEvent(battleID string, players []string) *BattleStartEvent {
	return &BattleStartEvent{
		BaseEvent: &BaseEvent{
			Type:      EventBattleStart,
			Timestamp: time.Now(),
		},
		BattleID: battleID,
		Players:  players,
	}
}

// EventManager 事件管理器
type EventManager struct {
	eventBus *EventBus
	logger   Logger
}

// NewEventManager 创建事件管理器
func NewEventManager(logger Logger) *EventManager {
	return &EventManager{
		eventBus: NewEventBus(logger),
		logger:   logger,
	}
}

// Initialize 初始化事件管理器
func (em *EventManager) Initialize(natsURL string) error {
	if err := em.eventBus.ConnectNATS(natsURL); err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// 注册默认事件处理器
	em.registerDefaultHandlers()

	em.logger.Info("Event manager initialized")
	return nil
}

// registerDefaultHandlers 注册默认事件处理器
func (em *EventManager) registerDefaultHandlers() {
	// 玩家登录事件处理
	em.eventBus.Subscribe(EventPlayerLogin, func(ctx context.Context, event Event) error {
		em.logger.Info("Player login event", "event", event.GetData())
		return nil
	})

	// 玩家移动事件处理
	em.eventBus.Subscribe(EventPlayerMove, func(ctx context.Context, event Event) error {
		em.logger.Debug("Player move event", "event", event.GetData())
		return nil
	})

	// 聊天事件处理
	em.eventBus.Subscribe(EventPlayerChat, func(ctx context.Context, event Event) error {
		em.logger.Info("Chat event", "event", event.GetData())
		return nil
	})
}

// GetEventBus 获取事件总线
func (em *EventManager) GetEventBus() *EventBus {
	return em.eventBus
}

// Close 关闭事件管理器
func (em *EventManager) Close() {
	em.eventBus.Close()
	em.logger.Info("Event manager closed")
}
