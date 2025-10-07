package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"greatestworks/internal/events"
	"greatestworks/internal/infrastructure/logging"

	"github.com/nats-io/nats.go"
)

// NATSSubscriber NATS消息订阅�?
type NATSSubscriber struct {
	conn          *nats.Conn
	logger        logger.Logger
	config        *SubscriberConfig
	subscriptions map[string]*nats.Subscription
	handlers      map[string][]MessageHandler
	mu            sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
	stats         *SubscriberStats
}

// SubscriberConfig 订阅者配�?
type SubscriberConfig struct {
	SubjectPrefix     string        `json:"subject_prefix" yaml:"subject_prefix"`
	QueueGroup        string        `json:"queue_group" yaml:"queue_group"`
	MaxInFlight       int           `json:"max_in_flight" yaml:"max_in_flight"`
	AckWait           time.Duration `json:"ack_wait" yaml:"ack_wait"`
	MaxDeliver        int           `json:"max_deliver" yaml:"max_deliver"`
	EnableMetrics     bool          `json:"enable_metrics" yaml:"enable_metrics"`
	ErrorRetryDelay   time.Duration `json:"error_retry_delay" yaml:"error_retry_delay"`
	DeadLetterSubject string        `json:"dead_letter_subject" yaml:"dead_letter_subject"`
}

// MessageHandler 消息处理器接�?
type MessageHandler interface {
	// Handle 处理消息
	Handle(ctx context.Context, msg *nats.Msg) error

	// GetHandlerName 获取处理器名�?
	GetHandlerName() string
}

// Subscriber 消息订阅者接�?
type Subscriber interface {
	// Subscribe 订阅主题
	Subscribe(subject string, handler MessageHandler) error

	// SubscribeQueue 订阅队列
	SubscribeQueue(subject, queue string, handler MessageHandler) error

	// SubscribeEvent 订阅领域事件
	SubscribeEvent(eventType string, handler events.EventHandler) error

	// SubscribeEventPattern 订阅事件模式
	SubscribeEventPattern(pattern string, handler events.EventHandler) error

	// Unsubscribe 取消订阅
	Unsubscribe(subject string) error

	// Start 启动订阅�?
	Start(ctx context.Context) error

	// Stop 停止订阅�?
	Stop() error

	// GetStats 获取统计信息
	GetStats() *SubscriberStats
}

// NewNATSSubscriber 创建NATS订阅�?
func NewNATSSubscriber(conn *nats.Conn, config *SubscriberConfig, logger logger.Logger) Subscriber {
	if config == nil {
		config = &SubscriberConfig{
			SubjectPrefix:     "greatestworks",
			QueueGroup:        "greatestworks-workers",
			MaxInFlight:       1000,
			AckWait:           30 * time.Second,
			MaxDeliver:        3,
			EnableMetrics:     true,
			ErrorRetryDelay:   5 * time.Second,
			DeadLetterSubject: "greatestworks.dead-letter",
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	logger.Info("NATS subscriber initialized successfully", "subject_prefix", config.SubjectPrefix, "queue_group", config.QueueGroup)

	return &NATSSubscriber{
		conn:          conn,
		logger:        logger,
		config:        config,
		subscriptions: make(map[string]*nats.Subscription),
		handlers:      make(map[string][]MessageHandler),
		ctx:           ctx,
		cancel:        cancel,
		stats: &SubscriberStats{
			StartTime: time.Now(),
			BySubject: make(map[string]*SubjectStats),
		},
	}
}

// Subscribe 订阅主题
func (s *NATSSubscriber) Subscribe(subject string, handler MessageHandler) error {
	fullSubject := s.buildSubject(subject)

	s.mu.Lock()
	defer s.mu.Unlock()

	// 添加处理�?
	s.handlers[fullSubject] = append(s.handlers[fullSubject], handler)

	// 如果已经有订阅，直接返回
	if _, exists := s.subscriptions[fullSubject]; exists {
		s.logger.Debug("Handler added to existing subscription", "subject", fullSubject, "handler", handler.GetHandlerName())
		return nil
	}

	// 创建新订�?
	sub, err := s.conn.Subscribe(fullSubject, s.createMessageCallback(fullSubject))
	if err != nil {
		s.logger.Error("Failed to subscribe to subject", "error", err, "subject", fullSubject)
		return fmt.Errorf("failed to subscribe to subject %s: %w", fullSubject, err)
	}

	s.subscriptions[fullSubject] = sub
	s.stats.SubscriptionCount++

	s.logger.Info("Subscribed to subject successfully", "subject", fullSubject, "handler", handler.GetHandlerName())
	return nil
}

// SubscribeQueue 订阅队列
func (s *NATSSubscriber) SubscribeQueue(subject, queue string, handler MessageHandler) error {
	fullSubject := s.buildSubject(subject)
	queueName := queue
	if queueName == "" {
		queueName = s.config.QueueGroup
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	subKey := fmt.Sprintf("%s:%s", fullSubject, queueName)

	// 添加处理�?
	s.handlers[subKey] = append(s.handlers[subKey], handler)

	// 如果已经有订阅，直接返回
	if _, exists := s.subscriptions[subKey]; exists {
		s.logger.Debug("Handler added to existing queue subscription", "subject", fullSubject, "queue", queueName, "handler", handler.GetHandlerName())
		return nil
	}

	// 创建新的队列订阅
	sub, err := s.conn.QueueSubscribe(fullSubject, queueName, s.createMessageCallback(subKey))
	if err != nil {
		s.logger.Error("Failed to subscribe to queue", "error", err, "subject", fullSubject, "queue", queueName)
		return fmt.Errorf("failed to subscribe to queue %s on subject %s: %w", queueName, fullSubject, err)
	}

	s.subscriptions[subKey] = sub
	s.stats.SubscriptionCount++

	s.logger.Info("Subscribed to queue successfully", "subject", fullSubject, "queue", queueName, "handler", handler.GetHandlerName())
	return nil
}

// SubscribeEvent 订阅领域事件
func (s *NATSSubscriber) SubscribeEvent(eventType string, handler events.EventHandler) error {
	// 构建事件主题模式
	subject := fmt.Sprintf("%s.events.*.%s", s.config.SubjectPrefix, eventType)

	// 创建事件处理器包装器
	eventWrapper := &eventHandlerWrapper{
		handler: handler,
		logger:  s.logger,
	}

	return s.SubscribeQueue(subject, s.config.QueueGroup, eventWrapper)
}

// SubscribeEventPattern 订阅事件模式
func (s *NATSSubscriber) SubscribeEventPattern(pattern string, handler events.EventHandler) error {
	// 构建完整的事件主题模�?
	fullPattern := fmt.Sprintf("%s.events.%s", s.config.SubjectPrefix, pattern)

	// 创建事件处理器包装器
	eventWrapper := &eventHandlerWrapper{
		handler: handler,
		logger:  s.logger,
	}

	return s.SubscribeQueue(fullPattern, s.config.QueueGroup, eventWrapper)
}

// Unsubscribe 取消订阅
func (s *NATSSubscriber) Unsubscribe(subject string) error {
	fullSubject := s.buildSubject(subject)

	s.mu.Lock()
	defer s.mu.Unlock()

	sub, exists := s.subscriptions[fullSubject]
	if !exists {
		return fmt.Errorf("subscription not found for subject: %s", fullSubject)
	}

	err := sub.Unsubscribe()
	if err != nil {
		s.logger.Error("Failed to unsubscribe from subject", "error", err, "subject", fullSubject)
		return fmt.Errorf("failed to unsubscribe from subject %s: %w", fullSubject, err)
	}

	delete(s.subscriptions, fullSubject)
	delete(s.handlers, fullSubject)
	s.stats.SubscriptionCount--

	s.logger.Info("Unsubscribed from subject successfully", "subject", fullSubject)
	return nil
}

// Start 启动订阅�?
func (s *NATSSubscriber) Start(ctx context.Context) error {
	s.logger.Info("Starting NATS subscriber")

	// 启动指标收集（如果启用）
	if s.config.EnableMetrics {
		go s.collectMetrics()
	}

	// 等待上下文取�?
	select {
	case <-ctx.Done():
		s.logger.Info("NATS subscriber context cancelled")
		return ctx.Err()
	case <-s.ctx.Done():
		s.logger.Info("NATS subscriber stopped")
		return nil
	}
}

// Stop 停止订阅�?
func (s *NATSSubscriber) Stop() error {
	s.logger.Info("Stopping NATS subscriber")

	s.mu.Lock()
	defer s.mu.Unlock()

	// 取消所有订�?
	for subject, sub := range s.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			s.logger.Error("Failed to unsubscribe during stop", "error", err, "subject", subject)
		}
	}

	// 清空订阅和处理器
	s.subscriptions = make(map[string]*nats.Subscription)
	s.handlers = make(map[string][]MessageHandler)

	// 取消上下�?
	s.cancel()

	s.logger.Info("NATS subscriber stopped successfully")
	return nil
}

// GetStats 获取统计信息
func (s *NATSSubscriber) GetStats() *SubscriberStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 创建统计信息副本
	stats := &SubscriberStats{
		SubscriptionCount: s.stats.SubscriptionCount,
		TotalProcessed:    s.stats.TotalProcessed,
		TotalFailed:       s.stats.TotalFailed,
		StartTime:         s.stats.StartTime,
		Uptime:            time.Since(s.stats.StartTime),
		BySubject:         make(map[string]*SubjectStats),
	}

	// 复制主题统计信息
	for subject, subjectStats := range s.stats.BySubject {
		stats.BySubject[subject] = &SubjectStats{
			ProcessedCount: subjectStats.ProcessedCount,
			FailedCount:    subjectStats.FailedCount,
			LastProcessed:  subjectStats.LastProcessed,
			AvgProcessTime: subjectStats.AvgProcessTime,
		}
	}

	return stats
}

// 私有方法

// buildSubject 构建完整主题
func (s *NATSSubscriber) buildSubject(subject string) string {
	if s.config.SubjectPrefix == "" {
		return subject
	}
	return fmt.Sprintf("%s.%s", s.config.SubjectPrefix, subject)
}

// createMessageCallback 创建消息回调
func (s *NATSSubscriber) createMessageCallback(subjectKey string) nats.MsgHandler {
	return func(msg *nats.Msg) {
		start := time.Now()

		// 更新统计信息
		s.updateStats(subjectKey, true, 0)

		// 获取处理�?
		s.mu.RLock()
		handlers := s.handlers[subjectKey]
		s.mu.RUnlock()

		if len(handlers) == 0 {
			s.logger.Warn("No handlers found for subject", "subject", msg.Subject)
			return
		}

		// 处理消息
		for _, handler := range handlers {
			if err := s.handleMessage(handler, msg); err != nil {
				s.logger.Error("Message handling failed", "error", err, "subject", msg.Subject, "handler", handler.GetHandlerName())
				s.updateStats(subjectKey, false, time.Since(start))

				// 发送到死信队列（如果配置了�?
				if s.config.DeadLetterSubject != "" {
					s.sendToDeadLetter(msg, err)
				}
				return
			}
		}

		// 更新成功统计
		s.updateStats(subjectKey, true, time.Since(start))
	}
}

// handleMessage 处理单个消息
func (s *NATSSubscriber) handleMessage(handler MessageHandler, msg *nats.Msg) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.config.AckWait)
	defer cancel()

	return handler.Handle(ctx, msg)
}

// updateStats 更新统计信息
func (s *NATSSubscriber) updateStats(subjectKey string, success bool, processTime time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if success {
		s.stats.TotalProcessed++
	} else {
		s.stats.TotalFailed++
	}

	// 更新主题统计
	subjectStats, exists := s.stats.BySubject[subjectKey]
	if !exists {
		subjectStats = &SubjectStats{}
		s.stats.BySubject[subjectKey] = subjectStats
	}

	if success {
		subjectStats.ProcessedCount++
		subjectStats.LastProcessed = time.Now()

		// 更新平均处理时间
		if subjectStats.AvgProcessTime == 0 {
			subjectStats.AvgProcessTime = processTime
		} else {
			subjectStats.AvgProcessTime = (subjectStats.AvgProcessTime + processTime) / 2
		}
	} else {
		subjectStats.FailedCount++
	}
}

// sendToDeadLetter 发送到死信队列
func (s *NATSSubscriber) sendToDeadLetter(msg *nats.Msg, err error) {
	deadLetterMsg := &DeadLetterMessage{
		OriginalSubject: msg.Subject,
		OriginalData:    msg.Data,
		Error:           err.Error(),
		Timestamp:       time.Now(),
		RetryCount:      s.getRetryCount(msg),
	}

	data, marshalErr := json.Marshal(deadLetterMsg)
	if marshalErr != nil {
		s.logger.Error("Failed to marshal dead letter message", "error", marshalErr)
		return
	}

	if publishErr := s.conn.Publish(s.config.DeadLetterSubject, data); publishErr != nil {
		s.logger.Error("Failed to publish to dead letter queue", "error", publishErr)
	}
}

// getRetryCount 获取重试次数
func (s *NATSSubscriber) getRetryCount(msg *nats.Msg) int {
	// 从消息头部获取重试次�?
	if msg.Header != nil {
		if retryCountStr := msg.Header.Get("Retry-Count"); retryCountStr != "" {
			// 解析重试次数
			// 这里简化处理，实际应该解析字符�?
			return 1
		}
	}
	return 0
}

// collectMetrics 收集指标
func (s *NATSSubscriber) collectMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats := s.GetStats()
			s.logger.Debug("Subscriber metrics",
				"subscriptions", stats.SubscriptionCount,
				"total_processed", stats.TotalProcessed,
				"total_failed", stats.TotalFailed,
				"uptime", stats.Uptime)
		case <-s.ctx.Done():
			return
		}
	}
}

// 事件处理器包装器
type eventHandlerWrapper struct {
	handler events.EventHandler
	logger  logger.Logger
}

func (w *eventHandlerWrapper) Handle(ctx context.Context, msg *nats.Msg) error {
	// 解析事件
	event, err := w.parseEvent(msg.Data)
	if err != nil {
		w.logger.Error("Failed to parse event", "error", err, "subject", msg.Subject)
		return fmt.Errorf("failed to parse event: %w", err)
	}

	// 处理事件
	// 将DomainEvent转换为events.Event
	eventWrapper := &events.BaseEvent{
		ID:        event.GetEventID(),
		Type:      event.GetEventType(),
		Data:      event,
		Timestamp: event.GetTimestamp(),
		UserID:    event.GetAggregateID(),
	}
	return w.handler.Handle(ctx, eventWrapper)
}

func (w *eventHandlerWrapper) GetHandlerName() string {
	return fmt.Sprintf("EventWrapper(%s)", w.handler.GetHandlerName())
}

// parseEvent 解析事件
func (w *eventHandlerWrapper) parseEvent(data []byte) (DomainEvent, error) {
	// 解析事件包装�?
	var eventWrapper map[string]interface{}
	if err := json.Unmarshal(data, &eventWrapper); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event wrapper: %w", err)
	}

	// 这里需要根据事件类型创建具体的事件实例
	// 简化处理，实际应该有事件注册表
	eventType, ok := eventWrapper["event_type"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid event_type")
	}

	// 创建基础事件（实际应该根据类型创建具体事件）
	baseEvent := &events.BaseEvent{}
	if eventData, exists := eventWrapper["data"]; exists {
		if eventBytes, err := json.Marshal(eventData); err == nil {
			json.Unmarshal(eventBytes, baseEvent)
		}
	}

	w.logger.Debug("Event parsed successfully", "event_type", eventType, "event_id", baseEvent.GetID())
	// 将BaseEvent转换为DomainEvent
	domainEvent := &BaseDomainEvent{
		eventID:     baseEvent.GetID(),
		eventType:   baseEvent.GetType(),
		aggregateID: baseEvent.GetPlayerID(),
		occurredAt:  baseEvent.GetTimestamp(),
		version:     1,
	}
	return domainEvent, nil
}

// 统计信息结构
type SubscriberStats struct {
	SubscriptionCount int64                    `json:"subscription_count"`
	TotalProcessed    int64                    `json:"total_processed"`
	TotalFailed       int64                    `json:"total_failed"`
	StartTime         time.Time                `json:"start_time"`
	Uptime            time.Duration            `json:"uptime"`
	BySubject         map[string]*SubjectStats `json:"by_subject"`
}

type SubjectStats struct {
	ProcessedCount int64         `json:"processed_count"`
	FailedCount    int64         `json:"failed_count"`
	LastProcessed  time.Time     `json:"last_processed"`
	AvgProcessTime time.Duration `json:"avg_process_time"`
}

// DeadLetterMessage 死信消息
type DeadLetterMessage struct {
	OriginalSubject string    `json:"original_subject"`
	OriginalData    []byte    `json:"original_data"`
	Error           string    `json:"error"`
	Timestamp       time.Time `json:"timestamp"`
	RetryCount      int       `json:"retry_count"`
}
