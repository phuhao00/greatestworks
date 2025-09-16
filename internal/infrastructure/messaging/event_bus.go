package messaging

import (
	"context"
	"fmt"
	"sync"
	"time"
	
	"greatestworks/internal/domain/player"
	"greatestworks/internal/domain/battle"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	EventID() string
	EventType() string
	AggregateID() string
	OccurredAt() time.Time
	Version() int
}

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
	GetHandlerName() string
	GetEventTypes() []string
}

// EventBus 事件总线
type EventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
	stats    *EventBusStats
}

// EventBusStats 事件总线统计
type EventBusStats struct {
	TotalPublished int64            `json:"total_published"`
	TotalHandled   int64            `json:"total_handled"`
	TotalFailed    int64            `json:"total_failed"`
	ByEventType    map[string]int64 `json:"by_event_type"`
	StartTime      time.Time        `json:"start_time"`
}

// NewEventBus 创建事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
		stats: &EventBusStats{
			ByEventType: make(map[string]int64),
			StartTime:   time.Now(),
		},
	}
}

// Subscribe 订阅事件
func (bus *EventBus) Subscribe(handler EventHandler) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	
	for _, eventType := range handler.GetEventTypes() {
		// 检查处理器是否已存在
		for _, existingHandler := range bus.handlers[eventType] {
			if existingHandler.GetHandlerName() == handler.GetHandlerName() {
				return fmt.Errorf("handler %s already subscribed to event %s", handler.GetHandlerName(), eventType)
			}
		}
		
		bus.handlers[eventType] = append(bus.handlers[eventType], handler)
	}
	
	return nil
}

// Unsubscribe 取消订阅事件
func (bus *EventBus) Unsubscribe(handlerName string, eventType string) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	
	handlers := bus.handlers[eventType]
	for i, handler := range handlers {
		if handler.GetHandlerName() == handlerName {
			bus.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			return nil
		}
	}
	
	return fmt.Errorf("handler %s not found for event %s", handlerName, eventType)
}

// Publish 发布事件
func (bus *EventBus) Publish(ctx context.Context, event DomainEvent) error {
	bus.mu.RLock()
	handlers := bus.handlers[event.EventType()]
	bus.mu.RUnlock()
	
	bus.stats.TotalPublished++
	bus.stats.ByEventType[event.EventType()]++
	
	if len(handlers) == 0 {
		// 没有处理器，直接返回
		return nil
	}
	
	// 并发处理所有处理器
	var wg sync.WaitGroup
	errorChan := make(chan error, len(handlers))
	
	for _, handler := range handlers {
		wg.Add(1)
		go func(h EventHandler) {
			defer wg.Done()
			
			if err := h.Handle(ctx, event); err != nil {
				errorChan <- fmt.Errorf("handler %s failed: %w", h.GetHandlerName(), err)
				bus.stats.TotalFailed++
			} else {
				bus.stats.TotalHandled++
			}
		}(handler)
	}
	
	wg.Wait()
	close(errorChan)
	
	// 收集错误
	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("event handling failed: %v", errors)
	}
	
	return nil
}

// PublishAsync 异步发布事件
func (bus *EventBus) PublishAsync(ctx context.Context, event DomainEvent) {
	go func() {
		if err := bus.Publish(ctx, event); err != nil {
			// 这里应该记录日志
			fmt.Printf("Async event handling failed: %v\n", err)
		}
	}()
}

// GetStats 获取统计信息
func (bus *EventBus) GetStats() *EventBusStats {
	bus.mu.RLock()
	defer bus.mu.RUnlock()
	
	stats := &EventBusStats{
		TotalPublished: bus.stats.TotalPublished,
		TotalHandled:   bus.stats.TotalHandled,
		TotalFailed:    bus.stats.TotalFailed,
		ByEventType:    make(map[string]int64),
		StartTime:      bus.stats.StartTime,
	}
	
	for eventType, count := range bus.stats.ByEventType {
		stats.ByEventType[eventType] = count
	}
	
	return stats
}

// PlayerEventHandler 玩家事件处理器基类
type PlayerEventHandler struct {
	name string
}

// NewPlayerEventHandler 创建玩家事件处理器
func NewPlayerEventHandler(name string) *PlayerEventHandler {
	return &PlayerEventHandler{name: name}
}

// GetHandlerName 获取处理器名称
func (h *PlayerEventHandler) GetHandlerName() string {
	return h.name
}

// GetEventTypes 获取处理的事件类型
func (h *PlayerEventHandler) GetEventTypes() []string {
	return []string{
		"PlayerCreated",
		"PlayerLevelUp",
		"PlayerOnline",
		"PlayerOffline",
		"PlayerMoved",
		"PlayerDied",
	}
}

// Handle 处理事件
func (h *PlayerEventHandler) Handle(ctx context.Context, event DomainEvent) error {
	switch event.EventType() {
	case "PlayerCreated":
		return h.handlePlayerCreated(ctx, event)
	case "PlayerLevelUp":
		return h.handlePlayerLevelUp(ctx, event)
	case "PlayerOnline":
		return h.handlePlayerOnline(ctx, event)
	case "PlayerOffline":
		return h.handlePlayerOffline(ctx, event)
	case "PlayerMoved":
		return h.handlePlayerMoved(ctx, event)
	case "PlayerDied":
		return h.handlePlayerDied(ctx, event)
	default:
		return fmt.Errorf("unknown event type: %s", event.EventType())
	}
}

// handlePlayerCreated 处理玩家创建事件
func (h *PlayerEventHandler) handlePlayerCreated(ctx context.Context, event DomainEvent) error {
	// 实现玩家创建后的逻辑，比如发送欢迎消息、初始化数据等
	fmt.Printf("Player created: %s\n", event.AggregateID())
	return nil
}

// handlePlayerLevelUp 处理玩家升级事件
func (h *PlayerEventHandler) handlePlayerLevelUp(ctx context.Context, event DomainEvent) error {
	// 实现玩家升级后的逻辑，比如发送奖励、更新排行榜等
	fmt.Printf("Player leveled up: %s\n", event.AggregateID())
	return nil
}

// handlePlayerOnline 处理玩家上线事件
func (h *PlayerEventHandler) handlePlayerOnline(ctx context.Context, event DomainEvent) error {
	// 实现玩家上线后的逻辑，比如通知好友、更新在线状态等
	fmt.Printf("Player online: %s\n", event.AggregateID())
	return nil
}

// handlePlayerOffline 处理玩家下线事件
func (h *PlayerEventHandler) handlePlayerOffline(ctx context.Context, event DomainEvent) error {
	// 实现玩家下线后的逻辑，比如保存数据、通知好友等
	fmt.Printf("Player offline: %s\n", event.AggregateID())
	return nil
}

// handlePlayerMoved 处理玩家移动事件
func (h *PlayerEventHandler) handlePlayerMoved(ctx context.Context, event DomainEvent) error {
	// 实现玩家移动后的逻辑，比如更新位置、检查触发器等
	fmt.Printf("Player moved: %s\n", event.AggregateID())
	return nil
}

// handlePlayerDied 处理玩家死亡事件
func (h *PlayerEventHandler) handlePlayerDied(ctx context.Context, event DomainEvent) error {
	// 实现玩家死亡后的逻辑，比如掉落物品、复活处理等
	fmt.Printf("Player died: %s\n", event.AggregateID())
	return nil
}

// BattleEventHandler 战斗事件处理器
type BattleEventHandler struct {
	name string
}

// NewBattleEventHandler 创建战斗事件处理器
func NewBattleEventHandler(name string) *BattleEventHandler {
	return &BattleEventHandler{name: name}
}

// GetHandlerName 获取处理器名称
func (h *BattleEventHandler) GetHandlerName() string {
	return h.name
}

// GetEventTypes 获取处理的事件类型
func (h *BattleEventHandler) GetEventTypes() []string {
	return []string{
		"BattleStarted",
		"BattleEnded",
		"PlayerJoinedBattle",
		"PlayerLeftBattle",
		"BattleActionExecuted",
	}
}

// Handle 处理事件
func (h *BattleEventHandler) Handle(ctx context.Context, event DomainEvent) error {
	switch event.EventType() {
	case "BattleStarted":
		return h.handleBattleStarted(ctx, event)
	case "BattleEnded":
		return h.handleBattleEnded(ctx, event)
	case "PlayerJoinedBattle":
		return h.handlePlayerJoinedBattle(ctx, event)
	case "PlayerLeftBattle":
		return h.handlePlayerLeftBattle(ctx, event)
	case "BattleActionExecuted":
		return h.handleBattleActionExecuted(ctx, event)
	default:
		return fmt.Errorf("unknown battle event type: %s", event.EventType())
	}
}

// handleBattleStarted 处理战斗开始事件
func (h *BattleEventHandler) handleBattleStarted(ctx context.Context, event DomainEvent) error {
	fmt.Printf("Battle started: %s\n", event.AggregateID())
	return nil
}

// handleBattleEnded 处理战斗结束事件
func (h *BattleEventHandler) handleBattleEnded(ctx context.Context, event DomainEvent) error {
	fmt.Printf("Battle ended: %s\n", event.AggregateID())
	return nil
}

// handlePlayerJoinedBattle 处理玩家加入战斗事件
func (h *BattleEventHandler) handlePlayerJoinedBattle(ctx context.Context, event DomainEvent) error {
	fmt.Printf("Player joined battle: %s\n", event.AggregateID())
	return nil
}

// handlePlayerLeftBattle 处理玩家离开战斗事件
func (h *BattleEventHandler) handlePlayerLeftBattle(ctx context.Context, event DomainEvent) error {
	fmt.Printf("Player left battle: %s\n", event.AggregateID())
	return nil
}

// handleBattleActionExecuted 处理战斗行动执行事件
func (h *BattleEventHandler) handleBattleActionExecuted(ctx context.Context, event DomainEvent) error {
	fmt.Printf("Battle action executed: %s\n", event.AggregateID())
	return nil
}