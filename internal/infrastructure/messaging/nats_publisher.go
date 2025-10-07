package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"greatestworks/internal/infrastructure/logging"

	"github.com/nats-io/nats.go"
	// "greatestworks/internal/domain/events" // TODO: 实现事件系统
)

// NATSPublisher NATS消息发布�?
type NATSPublisher struct {
	conn   *nats.Conn
	logger logger.Logger
	config *PublisherConfig
}

// PublisherConfig 发布者配�?
type PublisherConfig struct {
	SubjectPrefix   string        `json:"subject_prefix" yaml:"subject_prefix"`
	Timeout         time.Duration `json:"timeout" yaml:"timeout"`
	RetryAttempts   int           `json:"retry_attempts" yaml:"retry_attempts"`
	RetryDelay      time.Duration `json:"retry_delay" yaml:"retry_delay"`
	EnableMetrics   bool          `json:"enable_metrics" yaml:"enable_metrics"`
	CompressionType string        `json:"compression_type" yaml:"compression_type"`
}

// Publisher 消息发布者接�?
type Publisher interface {
	// PublishEvent 发布领域事件
	PublishEvent(ctx context.Context, event DomainEvent) error

	// PublishEventAsync 异步发布领域事件
	PublishEventAsync(ctx context.Context, event DomainEvent) error

	// PublishMessage 发布普通消�?
	PublishMessage(ctx context.Context, subject string, data interface{}) error

	// PublishMessageWithReply 发布带回复的消息
	PublishMessageWithReply(ctx context.Context, subject string, data interface{}, timeout time.Duration) (*nats.Msg, error)

	// PublishBatch 批量发布消息
	PublishBatch(ctx context.Context, messages []BatchMessage) error

	// Close 关闭发布�?
	Close() error
}

// BatchMessage 批量消息
type BatchMessage struct {
	Subject string            `json:"subject"`
	Data    interface{}       `json:"data"`
	Headers map[string]string `json:"headers,omitempty"`
}

// NewNATSPublisher 创建NATS发布�?
func NewNATSPublisher(conn *nats.Conn, config *PublisherConfig, logger logger.Logger) Publisher {
	if config == nil {
		config = &PublisherConfig{
			SubjectPrefix:   "greatestworks",
			Timeout:         5 * time.Second,
			RetryAttempts:   3,
			RetryDelay:      100 * time.Millisecond,
			EnableMetrics:   true,
			CompressionType: "none",
		}
	}

	logger.Info("NATS publisher initialized successfully", "subject_prefix", config.SubjectPrefix)

	return &NATSPublisher{
		conn:   conn,
		logger: logger,
		config: config,
	}
}

// PublishEvent 发布领域事件
func (p *NATSPublisher) PublishEvent(ctx context.Context, event DomainEvent) error {
	subject := p.buildEventSubject(event)

	// 序列化事�?
	data, err := p.serializeEvent(event)
	if err != nil {
		p.logger.Error("Failed to serialize event", "error", err, "event_type", event.GetEventType())
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	// 发布消息
	err = p.publishWithRetry(ctx, subject, data)
	if err != nil {
		p.logger.Error("Failed to publish event", "error", err, "subject", subject, "event_type", event.GetEventType())
		return fmt.Errorf("failed to publish event: %w", err)
	}

	p.logger.Debug("Event published successfully", "subject", subject, "event_type", event.GetEventType(), "event_id", event.GetEventID())
	return nil
}

// PublishEventAsync 异步发布领域事件
func (p *NATSPublisher) PublishEventAsync(ctx context.Context, event DomainEvent) error {
	go func() {
		if err := p.PublishEvent(context.Background(), event); err != nil {
			p.logger.Error("Failed to publish event asynchronously", "error", err, "event_type", event.GetEventType())
		}
	}()

	return nil
}

// PublishMessage 发布普通消�?
func (p *NATSPublisher) PublishMessage(ctx context.Context, subject string, data interface{}) error {
	fullSubject := p.buildSubject(subject)

	// 序列化数�?
	payload, err := p.serializeData(data)
	if err != nil {
		p.logger.Error("Failed to serialize message data", "error", err, "subject", subject)
		return fmt.Errorf("failed to serialize message data: %w", err)
	}

	// 发布消息
	err = p.publishWithRetry(ctx, fullSubject, payload)
	if err != nil {
		p.logger.Error("Failed to publish message", "error", err, "subject", fullSubject)
		return fmt.Errorf("failed to publish message: %w", err)
	}

	p.logger.Debug("Message published successfully", "subject", fullSubject)
	return nil
}

// PublishMessageWithReply 发布带回复的消息
func (p *NATSPublisher) PublishMessageWithReply(ctx context.Context, subject string, data interface{}, timeout time.Duration) (*nats.Msg, error) {
	fullSubject := p.buildSubject(subject)

	// 序列化数�?
	payload, err := p.serializeData(data)
	if err != nil {
		p.logger.Error("Failed to serialize request data", "error", err, "subject", subject)
		return nil, fmt.Errorf("failed to serialize request data: %w", err)
	}

	// 发送请求并等待回复
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	reply, err := p.conn.RequestWithContext(ctx, fullSubject, payload)
	if err != nil {
		p.logger.Error("Failed to send request with reply", "error", err, "subject", fullSubject)
		return nil, fmt.Errorf("failed to send request with reply: %w", err)
	}

	p.logger.Debug("Request with reply sent successfully", "subject", fullSubject, "reply_size", len(reply.Data))
	return reply, nil
}

// PublishBatch 批量发布消息
func (p *NATSPublisher) PublishBatch(ctx context.Context, messages []BatchMessage) error {
	if len(messages) == 0 {
		return nil
	}

	// 使用NATS的批量发布功�?
	var errors []error

	for _, msg := range messages {
		fullSubject := p.buildSubject(msg.Subject)

		// 序列化数�?
		payload, err := p.serializeData(msg.Data)
		if err != nil {
			p.logger.Error("Failed to serialize batch message data", "error", err, "subject", msg.Subject)
			errors = append(errors, fmt.Errorf("failed to serialize message for subject %s: %w", msg.Subject, err))
			continue
		}

		// 创建NATS消息
		natsMsg := &nats.Msg{
			Subject: fullSubject,
			Data:    payload,
		}

		// 添加头部信息
		if len(msg.Headers) > 0 {
			natsMsg.Header = make(nats.Header)
			for k, v := range msg.Headers {
				natsMsg.Header.Set(k, v)
			}
		}

		// 发布消息
		if err := p.conn.PublishMsg(natsMsg); err != nil {
			p.logger.Error("Failed to publish batch message", "error", err, "subject", fullSubject)
			errors = append(errors, fmt.Errorf("failed to publish message for subject %s: %w", msg.Subject, err))
		}
	}

	// 刷新连接以确保消息发�?
	if err := p.conn.Flush(); err != nil {
		p.logger.Error("Failed to flush batch messages", "error", err)
		errors = append(errors, fmt.Errorf("failed to flush batch messages: %w", err))
	}

	if len(errors) > 0 {
		p.logger.Error("Batch publish completed with errors", "total_messages", len(messages), "error_count", len(errors))
		return fmt.Errorf("batch publish failed with %d errors: %v", len(errors), errors[0])
	}

	p.logger.Debug("Batch messages published successfully", "count", len(messages))
	return nil
}

// Close 关闭发布�?
func (p *NATSPublisher) Close() error {
	if p.conn != nil && !p.conn.IsClosed() {
		p.conn.Close()
	}

	p.logger.Info("NATS publisher closed successfully")
	return nil
}

// 私有方法

// buildEventSubject 构建事件主题
func (p *NATSPublisher) buildEventSubject(event DomainEvent) string {
	return fmt.Sprintf("%s.events.%s.%s", p.config.SubjectPrefix, event.GetAggregateType(), event.GetEventType())
}

// buildSubject 构建普通主�?
func (p *NATSPublisher) buildSubject(subject string) string {
	if p.config.SubjectPrefix == "" {
		return subject
	}
	return fmt.Sprintf("%s.%s", p.config.SubjectPrefix, subject)
}

// serializeEvent 序列化事�?
func (p *NATSPublisher) serializeEvent(event DomainEvent) ([]byte, error) {
	// 创建事件包装�?
	eventWrapper := map[string]interface{}{
		"event_id":       event.GetEventID(),
		"event_type":     event.GetEventType(),
		"aggregate_id":   event.GetAggregateID(),
		"aggregate_type": event.GetAggregateType(),
		"version":        event.GetVersion(),
		"timestamp":      event.GetTimestamp(),
		"data":           event,
		"metadata":       event.GetMetadata(),
	}

	return json.Marshal(eventWrapper)
}

// serializeData 序列化数�?
func (p *NATSPublisher) serializeData(data interface{}) ([]byte, error) {
	// 如果已经是字节数组，直接返回
	if bytes, ok := data.([]byte); ok {
		return bytes, nil
	}

	// 如果是字符串，转换为字节数组
	if str, ok := data.(string); ok {
		return []byte(str), nil
	}

	// 否则使用JSON序列�?
	return json.Marshal(data)
}

// publishWithRetry 带重试的发布
func (p *NATSPublisher) publishWithRetry(ctx context.Context, subject string, data []byte) error {
	var lastErr error

	for i := 0; i <= p.config.RetryAttempts; i++ {
		// 检查上下文是否已取�?
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 尝试发布
		err := p.conn.Publish(subject, data)
		if err == nil {
			// 刷新连接以确保消息发�?
			if flushErr := p.conn.Flush(); flushErr == nil {
				return nil
			} else {
				lastErr = flushErr
			}
		} else {
			lastErr = err
		}

		// 如果不是最后一次尝试，等待后重�?
		if i < p.config.RetryAttempts {
			p.logger.Warn("Publish attempt failed, retrying", "attempt", i+1, "error", lastErr, "subject", subject)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(p.config.RetryDelay):
				// 继续重试
			}
		}
	}

	return lastErr
}

// EventMetrics 事件指标
type EventMetrics struct {
	PublishedCount int64            `json:"published_count"`
	FailedCount    int64            `json:"failed_count"`
	ByEventType    map[string]int64 `json:"by_event_type"`
	BySubject      map[string]int64 `json:"by_subject"`
	LastPublished  time.Time        `json:"last_published"`
}

// PublisherStats 发布者统计信�?
type PublisherStats struct {
	TotalPublished int64                    `json:"total_published"`
	TotalFailed    int64                    `json:"total_failed"`
	EventMetrics   map[string]*EventMetrics `json:"event_metrics"`
	Uptime         time.Duration            `json:"uptime"`
	StartTime      time.Time                `json:"start_time"`
}
