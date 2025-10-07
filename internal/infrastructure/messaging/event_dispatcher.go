package messaging

import (
	"context"
	"fmt"
	"sync"
	"time"

	"greatestworks/internal/events"
	"greatestworks/internal/infrastructure/logger"
)

// EventDispatcher 事件分发器
type EventDispatcher struct {
	publisher  Publisher
	subscriber Subscriber
	logger     logger.Logger
	config     *DispatcherConfig
	handlers   map[string][]events.EventHandler
	mu         sync.RWMutex
	stats      *DispatcherStats
	ctx        context.Context
	cancel     context.CancelFunc
	eventQueue chan *EventMessage
	workerPool *WorkerPool
}

// DispatcherConfig 分发器配置
type DispatcherConfig struct {
	WorkerCount       int           `json:"worker_count" yaml:"worker_count"`
	QueueSize         int           `json:"queue_size" yaml:"queue_size"`
	BatchSize         int           `json:"batch_size" yaml:"batch_size"`
	BatchTimeout      time.Duration `json:"batch_timeout" yaml:"batch_timeout"`
	RetryAttempts     int           `json:"retry_attempts" yaml:"retry_attempts"`
	RetryDelay        time.Duration `json:"retry_delay" yaml:"retry_delay"`
	EnableMetrics     bool          `json:"enable_metrics" yaml:"enable_metrics"`
	EnableDeadLetter  bool          `json:"enable_dead_letter" yaml:"enable_dead_letter"`
	MaxProcessingTime time.Duration `json:"max_processing_time" yaml:"max_processing_time"`
}

// EventMessage 事件消息
type EventMessage struct {
	Event      DomainEvent       `json:"event"`
	Timestamp  time.Time         `json:"timestamp"`
	RetryCount int               `json:"retry_count"`
	Metadata   map[string]string `json:"metadata"`
}

// Dispatcher 事件分发器接口
type Dispatcher interface {
	// RegisterHandler 注册事件处理器
	RegisterHandler(eventType string, handler events.EventHandler) error

	// UnregisterHandler 取消注册事件处理器
	UnregisterHandler(eventType string, handlerName string) error

	// Dispatch 分发事件
	Dispatch(ctx context.Context, event DomainEvent) error

	// DispatchAsync 异步分发事件
	DispatchAsync(ctx context.Context, event DomainEvent) error

	// DispatchBatch 批量分发事件
	DispatchBatch(ctx context.Context, events []DomainEvent) error

	// Start 启动分发器
	Start(ctx context.Context) error

	// Stop 停止分发器
	Stop() error

	// GetStats 获取统计信息
	GetStats() *DispatcherStats
}

// NewEventDispatcher 创建事件分发器
func NewEventDispatcher(publisher Publisher, subscriber Subscriber, config *DispatcherConfig, logger logger.Logger) Dispatcher {
	if config == nil {
		config = &DispatcherConfig{
			WorkerCount:       10,
			QueueSize:         1000,
			BatchSize:         100,
			BatchTimeout:      5 * time.Second,
			RetryAttempts:     3,
			RetryDelay:        1 * time.Second,
			EnableMetrics:     true,
			EnableDeadLetter:  true,
			MaxProcessingTime: 30 * time.Second,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	d := &EventDispatcher{
		publisher:  publisher,
		subscriber: subscriber,
		logger:     logger,
		config:     config,
		handlers:   make(map[string][]events.EventHandler),
		ctx:        ctx,
		cancel:     cancel,
		eventQueue: make(chan *EventMessage, config.QueueSize),
		stats: &DispatcherStats{
			StartTime:   time.Now(),
			ByEventType: make(map[string]*EventTypeStats),
			ByHandler:   make(map[string]*HandlerStats),
		},
	}

	// 创建工作池
	d.workerPool = NewWorkerPool(config.WorkerCount, d.processEvent, logger)

	logger.Info("Event dispatcher initialized successfully", "worker_count", config.WorkerCount, "queue_size", config.QueueSize)
	return d
}

// RegisterHandler 注册事件处理器
func (d *EventDispatcher) RegisterHandler(eventType string, handler events.EventHandler) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 检查处理器是否已存在
	for _, existingHandler := range d.handlers[eventType] {
		if existingHandler.GetHandlerName() == handler.GetHandlerName() {
			return fmt.Errorf("handler %s already registered for event type %s", handler.GetHandlerName(), eventType)
		}
	}

	// 添加处理器
	d.handlers[eventType] = append(d.handlers[eventType], handler)

	// 订阅事件
	if err := d.subscriber.SubscribeEvent(eventType, handler); err != nil {
		d.logger.Error("Failed to subscribe to event", "error", err, "event_type", eventType, "handler", handler.GetHandlerName())
		return fmt.Errorf("failed to subscribe to event %s: %w", eventType, err)
	}

	d.logger.Info("Event handler registered successfully", "event_type", eventType, "handler", handler.GetHandlerName())
	return nil
}

// UnregisterHandler 取消注册事件处理器
func (d *EventDispatcher) UnregisterHandler(eventType string, handlerName string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	handlers := d.handlers[eventType]
	for i, handler := range handlers {
		if handler.GetHandlerName() == handlerName {
			// 移除处理器
			d.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)

			// 如果没有更多处理器，取消订阅
			if len(d.handlers[eventType]) == 0 {
				delete(d.handlers, eventType)
				// 这里应该取消订阅，但NATS订阅是按主题的，可能有其他处理器
			}

			d.logger.Info("Event handler unregistered successfully", "event_type", eventType, "handler", handlerName)
			return nil
		}
	}

	return fmt.Errorf("handler %s not found for event type %s", handlerName, eventType)
}

// Dispatch 分发事件
func (d *EventDispatcher) Dispatch(ctx context.Context, event DomainEvent) error {
	// 直接发布事件
	err := d.publisher.PublishEvent(ctx, event)
	if err != nil {
		d.updateStats(event.GetEventType(), false, 0, "publish_error")
		d.logger.Error("Failed to dispatch event", "error", err, "event_type", event.GetEventType(), "event_id", event.GetEventID())
		return fmt.Errorf("failed to dispatch event: %w", err)
	}

	d.updateStats(event.GetEventType(), true, 0, "dispatched")
	d.logger.Debug("Event dispatched successfully", "event_type", event.GetEventType(), "event_id", event.GetEventID())
	return nil
}

// DispatchAsync 异步分发事件
func (d *EventDispatcher) DispatchAsync(ctx context.Context, event DomainEvent) error {
	eventMsg := &EventMessage{
		Event:      event,
		Timestamp:  time.Now(),
		RetryCount: 0,
		Metadata:   make(map[string]string),
	}

	// 添加到队列
	select {
	case d.eventQueue <- eventMsg:
		d.updateStats(event.GetEventType(), true, 0, "queued")
		d.logger.Debug("Event queued for async dispatch", "event_type", event.GetEventType(), "event_id", event.GetEventID())
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		d.updateStats(event.GetEventType(), false, 0, "queue_full")
		return fmt.Errorf("event queue is full")
	}
}

// DispatchBatch 批量分发事件
func (d *EventDispatcher) DispatchBatch(ctx context.Context, events []DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	// 构建批量消息
	batchMessages := make([]BatchMessage, len(events))
	for i, event := range events {
		batchMessages[i] = BatchMessage{
			Subject: fmt.Sprintf("events.%s.%s", event.GetAggregateType(), event.GetEventType()),
			Data:    event,
		}
	}

	// 批量发布
	err := d.publisher.PublishBatch(ctx, batchMessages)
	if err != nil {
		for _, event := range events {
			d.updateStats(event.GetEventType(), false, 0, "batch_error")
		}
		d.logger.Error("Failed to dispatch batch events", "error", err, "count", len(events))
		return fmt.Errorf("failed to dispatch batch events: %w", err)
	}

	// 更新统计
	for _, event := range events {
		d.updateStats(event.GetEventType(), true, 0, "batch_dispatched")
	}

	d.logger.Debug("Batch events dispatched successfully", "count", len(events))
	return nil
}

// Start 启动分发器
func (d *EventDispatcher) Start(ctx context.Context) error {
	d.logger.Info("Starting event dispatcher")

	// 启动工作池
	if err := d.workerPool.Start(d.ctx); err != nil {
		return fmt.Errorf("failed to start worker pool: %w", err)
	}

	// 启动队列处理器
	go d.processQueue()

	// 启动批处理器
	go d.processBatch()

	// 启动指标收集
	if d.config.EnableMetrics {
		go d.collectMetrics()
	}

	// 启动订阅者
	go func() {
		if err := d.subscriber.Start(d.ctx); err != nil {
			d.logger.Error("Subscriber stopped with error", "error", err)
		}
	}()

	// 等待上下文取消
	select {
	case <-ctx.Done():
		d.logger.Info("Event dispatcher context cancelled")
		return ctx.Err()
	case <-d.ctx.Done():
		d.logger.Info("Event dispatcher stopped")
		return nil
	}
}

// Stop 停止分发器
func (d *EventDispatcher) Stop() error {
	d.logger.Info("Stopping event dispatcher")

	// 取消上下文
	d.cancel()

	// 停止工作池
	if err := d.workerPool.Stop(); err != nil {
		d.logger.Error("Failed to stop worker pool", "error", err)
	}

	// 停止订阅者
	if err := d.subscriber.Stop(); err != nil {
		d.logger.Error("Failed to stop subscriber", "error", err)
	}

	// 关闭发布者
	if err := d.publisher.Close(); err != nil {
		d.logger.Error("Failed to close publisher", "error", err)
	}

	// 关闭事件队列
	close(d.eventQueue)

	d.logger.Info("Event dispatcher stopped successfully")
	return nil
}

// GetStats 获取统计信息
func (d *EventDispatcher) GetStats() *DispatcherStats {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// 创建统计信息副本
	stats := &DispatcherStats{
		TotalDispatched: d.stats.TotalDispatched,
		TotalFailed:     d.stats.TotalFailed,
		QueueSize:       int64(len(d.eventQueue)),
		WorkerCount:     int64(d.config.WorkerCount),
		StartTime:       d.stats.StartTime,
		Uptime:          time.Since(d.stats.StartTime),
		ByEventType:     make(map[string]*EventTypeStats),
		ByHandler:       make(map[string]*HandlerStats),
	}

	// 复制事件类型统计
	for eventType, eventStats := range d.stats.ByEventType {
		stats.ByEventType[eventType] = &EventTypeStats{
			DispatchedCount: eventStats.DispatchedCount,
			FailedCount:     eventStats.FailedCount,
			LastDispatched:  eventStats.LastDispatched,
			AvgProcessTime:  eventStats.AvgProcessTime,
		}
	}

	// 复制处理器统计
	for handlerName, handlerStats := range d.stats.ByHandler {
		stats.ByHandler[handlerName] = &HandlerStats{
			ProcessedCount: handlerStats.ProcessedCount,
			FailedCount:    handlerStats.FailedCount,
			LastProcessed:  handlerStats.LastProcessed,
			AvgProcessTime: handlerStats.AvgProcessTime,
		}
	}

	return stats
}

// 私有方法

// processQueue 处理事件队列
func (d *EventDispatcher) processQueue() {
	for {
		select {
		case eventMsg := <-d.eventQueue:
			if eventMsg != nil {
				// 提交到工作池
				d.workerPool.Submit(eventMsg)
			}
		case <-d.ctx.Done():
			return
		}
	}
}

// processBatch 处理批量事件
func (d *EventDispatcher) processBatch() {
	ticker := time.NewTicker(d.config.BatchTimeout)
	defer ticker.Stop()

	batch := make([]*EventMessage, 0, d.config.BatchSize)

	for {
		select {
		case eventMsg := <-d.eventQueue:
			if eventMsg != nil {
				batch = append(batch, eventMsg)

				// 如果批次满了，立即处理
				if len(batch) >= d.config.BatchSize {
					d.processBatchEvents(batch)
					batch = batch[:0] // 重置批次
				}
			}
		case <-ticker.C:
			// 定时处理批次
			if len(batch) > 0 {
				d.processBatchEvents(batch)
				batch = batch[:0] // 重置批次
			}
		case <-d.ctx.Done():
			// 处理剩余批次
			if len(batch) > 0 {
				d.processBatchEvents(batch)
			}
			return
		}
	}
}

// processBatchEvents 处理批量事件
func (d *EventDispatcher) processBatchEvents(batch []*EventMessage) {
	events := make([]DomainEvent, len(batch))
	for i, eventMsg := range batch {
		events[i] = eventMsg.Event
	}

	ctx, cancel := context.WithTimeout(d.ctx, d.config.MaxProcessingTime)
	defer cancel()

	if err := d.DispatchBatch(ctx, events); err != nil {
		d.logger.Error("Failed to process batch events", "error", err, "count", len(batch))

		// 重试单个事件
		for _, eventMsg := range batch {
			d.retryEvent(eventMsg)
		}
	}
}

// processEvent 处理单个事件
func (d *EventDispatcher) processEvent(data interface{}) error {
	eventMsg, ok := data.(*EventMessage)
	if !ok {
		return fmt.Errorf("invalid event message type")
	}

	start := time.Now()

	ctx, cancel := context.WithTimeout(d.ctx, d.config.MaxProcessingTime)
	defer cancel()

	err := d.Dispatch(ctx, eventMsg.Event)
	processTime := time.Since(start)

	if err != nil {
		d.updateStats(eventMsg.Event.GetEventType(), false, processTime, "process_error")

		// 重试逻辑
		if eventMsg.RetryCount < d.config.RetryAttempts {
			d.retryEvent(eventMsg)
		} else {
			d.logger.Error("Event processing failed after max retries", "error", err, "event_type", eventMsg.Event.GetEventType(), "retry_count", eventMsg.RetryCount)
		}
		return err
	}

	d.updateStats(eventMsg.Event.GetEventType(), true, processTime, "processed")
	return nil
}

// retryEvent 重试事件
func (d *EventDispatcher) retryEvent(eventMsg *EventMessage) {
	eventMsg.RetryCount++

	// 延迟重试
	go func() {
		time.Sleep(d.config.RetryDelay * time.Duration(eventMsg.RetryCount))

		select {
		case d.eventQueue <- eventMsg:
			d.logger.Debug("Event queued for retry", "event_type", eventMsg.Event.GetEventType(), "retry_count", eventMsg.RetryCount)
		case <-d.ctx.Done():
			return
		default:
			d.logger.Error("Failed to queue event for retry - queue full", "event_type", eventMsg.Event.GetEventType())
		}
	}()
}

// updateStats 更新统计信息
func (d *EventDispatcher) updateStats(eventType string, success bool, processTime time.Duration, operation string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if success {
		d.stats.TotalDispatched++
	} else {
		d.stats.TotalFailed++
	}

	// 更新事件类型统计
	eventStats, exists := d.stats.ByEventType[eventType]
	if !exists {
		eventStats = &EventTypeStats{}
		d.stats.ByEventType[eventType] = eventStats
	}

	if success {
		eventStats.DispatchedCount++
		eventStats.LastDispatched = time.Now()

		// 更新平均处理时间
		if eventStats.AvgProcessTime == 0 {
			eventStats.AvgProcessTime = processTime
		} else {
			eventStats.AvgProcessTime = (eventStats.AvgProcessTime + processTime) / 2
		}
	} else {
		eventStats.FailedCount++
	}
}

// collectMetrics 收集指标
func (d *EventDispatcher) collectMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats := d.GetStats()
			d.logger.Debug("Dispatcher metrics",
				"total_dispatched", stats.TotalDispatched,
				"total_failed", stats.TotalFailed,
				"queue_size", stats.QueueSize,
				"worker_count", stats.WorkerCount,
				"uptime", stats.Uptime)
		case <-d.ctx.Done():
			return
		}
	}
}

// 统计信息结构
type DispatcherStats struct {
	TotalDispatched int64                      `json:"total_dispatched"`
	TotalFailed     int64                      `json:"total_failed"`
	QueueSize       int64                      `json:"queue_size"`
	WorkerCount     int64                      `json:"worker_count"`
	StartTime       time.Time                  `json:"start_time"`
	Uptime          time.Duration              `json:"uptime"`
	ByEventType     map[string]*EventTypeStats `json:"by_event_type"`
	ByHandler       map[string]*HandlerStats   `json:"by_handler"`
}

type EventTypeStats struct {
	DispatchedCount int64         `json:"dispatched_count"`
	FailedCount     int64         `json:"failed_count"`
	LastDispatched  time.Time     `json:"last_dispatched"`
	AvgProcessTime  time.Duration `json:"avg_process_time"`
}

type HandlerStats struct {
	ProcessedCount int64         `json:"processed_count"`
	FailedCount    int64         `json:"failed_count"`
	LastProcessed  time.Time     `json:"last_processed"`
	AvgProcessTime time.Duration `json:"avg_process_time"`
}
