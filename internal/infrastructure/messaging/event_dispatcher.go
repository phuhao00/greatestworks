package messaging

import (
	"context"
	"fmt"
	"sync"
	"time"

	"greatestworks/internal/events"
	"greatestworks/internal/infrastructure/logging"
)

// EventDispatcher 事件分发器
type EventDispatcher struct {
	publisher  Publisher
	subscriber Subscriber
	logger     logging.Logger
	config     *DispatcherConfig
	handlers   map[string][]events.EventHandler
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

// Publisher 发布者接口
type Publisher interface {
	Publish(ctx context.Context, topic string, event events.Event) error
}

// Subscriber 订阅者接口
type Subscriber interface {
	Subscribe(ctx context.Context, topic string, handler events.EventHandler) error
	Unsubscribe(ctx context.Context, topic string, handler events.EventHandler) error
}

// DispatcherConfig 分发器配置
type DispatcherConfig struct {
	MaxRetries     int           `json:"max_retries" yaml:"max_retries"`
	RetryDelay     time.Duration `json:"retry_delay" yaml:"retry_delay"`
	MaxConcurrency int           `json:"max_concurrency" yaml:"max_concurrency"`
	BufferSize     int           `json:"buffer_size" yaml:"buffer_size"`
	EnableMetrics  bool          `json:"enable_metrics" yaml:"enable_metrics"`
	EnableTracing  bool          `json:"enable_tracing" yaml:"enable_tracing"`
}

// NewEventDispatcher 创建事件分发器
func NewEventDispatcher(publisher Publisher, subscriber Subscriber, config *DispatcherConfig, logger logging.Logger) *EventDispatcher {
	if config == nil {
		config = &DispatcherConfig{
			MaxRetries:     3,
			RetryDelay:     time.Second,
			MaxConcurrency: 10,
			BufferSize:     1000,
			EnableMetrics:  true,
			EnableTracing:  false,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &EventDispatcher{
		publisher:  publisher,
		subscriber: subscriber,
		logger:     logger,
		config:     config,
		handlers:   make(map[string][]events.EventHandler),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start 启动事件分发器
func (d *EventDispatcher) Start() error {
	d.logger.Info("Event dispatcher starting")

	// 启动订阅处理
	go d.subscribeLoop()

	d.logger.Info("Event dispatcher started")
	return nil
}

// Stop 停止事件分发器
func (d *EventDispatcher) Stop() error {
	d.logger.Info("Event dispatcher stopping")

	// 取消上下文
	d.cancel()

	// 等待所有goroutine完成
	d.wg.Wait()

	d.logger.Info("Event dispatcher stopped")
	return nil
}

// Publish 发布事件
func (d *EventDispatcher) Publish(ctx context.Context, event events.Event) error {
	if d.publisher == nil {
		return fmt.Errorf("publisher not configured")
	}

	// 获取事件类型作为主题
	topic := event.GetEventType()

	// 发布事件
	if err := d.publisher.Publish(ctx, topic, event); err != nil {
		d.logger.Error("Failed to publish event", err, logging.Fields{
			"event_type": event.GetEventType(),
		})
		return fmt.Errorf("failed to publish event: %w", err)
	}

	d.logger.Debug("Event published", logging.Fields{
		"event_type":   event.GetEventType(),
		"aggregate_id": event.GetAggregateID(),
	})
	return nil
}

// Subscribe 订阅事件
func (d *EventDispatcher) Subscribe(eventType string, handler events.EventHandler) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 添加到本地处理器列表
	d.handlers[eventType] = append(d.handlers[eventType], handler)

	// 如果配置了订阅者，也添加到远程订阅
	if d.subscriber != nil {
		if err := d.subscriber.Subscribe(d.ctx, eventType, handler); err != nil {
			d.logger.Error("Failed to subscribe to remote event", err, logging.Fields{
				"event_type": eventType,
			})
			return fmt.Errorf("failed to subscribe to remote event: %w", err)
		}
	}

	d.logger.Debug("Event handler subscribed", logging.Fields{
		"event_type": eventType,
	})
	return nil
}

// Unsubscribe 取消订阅事件
func (d *EventDispatcher) Unsubscribe(eventType string, handler events.EventHandler) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 从本地处理器列表移除
	if handlers, exists := d.handlers[eventType]; exists {
		for i, h := range handlers {
			if h == handler {
				d.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}

	// 如果配置了订阅者，也从远程订阅移除
	if d.subscriber != nil {
		if err := d.subscriber.Unsubscribe(d.ctx, eventType, handler); err != nil {
			d.logger.Error("Failed to unsubscribe from remote event", err, logging.Fields{
				"event_type": eventType,
			})
			return fmt.Errorf("failed to unsubscribe from remote event: %w", err)
		}
	}

	d.logger.Debug("Event handler unsubscribed", logging.Fields{
		"event_type": eventType,
	})
	return nil
}

// Dispatch 分发事件到本地处理器
func (d *EventDispatcher) Dispatch(ctx context.Context, event events.Event) error {
	eventType := event.GetEventType()

	d.mu.RLock()
	handlers, exists := d.handlers[eventType]
	d.mu.RUnlock()

	if !exists || len(handlers) == 0 {
		d.logger.Debug("No handlers found for event", logging.Fields{
			"event_type": eventType,
		})
		return nil
	}

	// 并发处理所有处理器
	var wg sync.WaitGroup
	for _, handler := range handlers {
		wg.Add(1)
		go func(h events.EventHandler) {
			defer wg.Done()
			d.handleEvent(ctx, h, event)
		}(handler)
	}

	wg.Wait()
	return nil
}

// 私有方法

// subscribeLoop 订阅循环
func (d *EventDispatcher) subscribeLoop() {
	d.wg.Add(1)
	defer d.wg.Done()

	d.logger.Debug("Event dispatcher subscribe loop started")

	for {
		select {
		case <-d.ctx.Done():
			d.logger.Debug("Event dispatcher subscribe loop stopped")
			return
		default:
			// 这里应该处理来自订阅者的消息
			// 简化实现，实际项目中应该从消息队列接收事件
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// handleEvent 处理事件
func (d *EventDispatcher) handleEvent(ctx context.Context, handler events.EventHandler, event events.Event) {
	// 实现重试逻辑
	for attempt := 0; attempt <= d.config.MaxRetries; attempt++ {
		if err := handler.Handle(ctx, event); err != nil {
			d.logger.Error("Event handler failed", err, logging.Fields{
				"event_type": event.GetEventType(),
				"attempt":    attempt + 1,
			})

			if attempt < d.config.MaxRetries {
				// 等待重试延迟
				time.Sleep(d.config.RetryDelay)
				continue
			}

			// 达到最大重试次数，记录错误
			d.logger.Error("Event handler failed after max retries", err, logging.Fields{
				"event_type": event.GetEventType(),
			})
			return
		}

		// 处理成功
		d.logger.Debug("Event handled successfully", logging.Fields{
			"event_type": event.GetEventType(),
		})
		return
	}
}

// GetStats 获取统计信息
func (d *EventDispatcher) GetStats() map[string]interface{} {
	d.mu.RLock()
	defer d.mu.RUnlock()

	stats := make(map[string]interface{})
	stats["total_event_types"] = len(d.handlers)

	totalHandlers := 0
	for _, handlers := range d.handlers {
		totalHandlers += len(handlers)
	}
	stats["total_handlers"] = totalHandlers

	return stats
}
