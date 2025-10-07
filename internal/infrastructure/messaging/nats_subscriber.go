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

// NATSSubscriber NATS消息订阅器
type NATSSubscriber struct {
	conn          *nats.Conn
	logger        logging.Logger
	config        *SubscriberConfig
	subscriptions map[string]*nats.Subscription
	mu            sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
}

// SubscriberConfig 订阅器配置
type SubscriberConfig struct {
	URL          string        `json:"url" yaml:"url"`
	ClusterID    string        `json:"cluster_id" yaml:"cluster_id"`
	ClientID     string        `json:"client_id" yaml:"client_id"`
	Timeout      time.Duration `json:"timeout" yaml:"timeout"`
	MaxReconn    int           `json:"max_reconn" yaml:"max_reconn"`
	ReconnWait   time.Duration `json:"reconn_wait" yaml:"reconn_wait"`
	PingInterval time.Duration `json:"ping_interval" yaml:"ping_interval"`
	MaxPingsOut  int           `json:"max_pings_out" yaml:"max_pings_out"`
}

// NewNATSSubscriber 创建NATS订阅器
func NewNATSSubscriber(config *SubscriberConfig, logger logging.Logger) (*NATSSubscriber, error) {
	if config == nil {
		config = &SubscriberConfig{
			URL:          "nats://localhost:4222",
			ClusterID:    "test-cluster",
			ClientID:     "subscriber",
			Timeout:      30 * time.Second,
			MaxReconn:    10,
			ReconnWait:   2 * time.Second,
			PingInterval: 2 * time.Minute,
			MaxPingsOut:  2,
		}
	}

	// 连接NATS服务器
	conn, err := nats.Connect(config.URL, nats.Name(config.ClientID))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	subscriber := &NATSSubscriber{
		conn:          conn,
		logger:        logger,
		config:        config,
		subscriptions: make(map[string]*nats.Subscription),
		ctx:           ctx,
		cancel:        cancel,
	}

	logger.Info("NATS subscriber created", logging.Fields{
		"url":       config.URL,
		"client_id": config.ClientID,
	})
	return subscriber, nil
}

// Subscribe 订阅消息
func (s *NATSSubscriber) Subscribe(subject string, handler events.EventHandler) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查是否已经订阅
	if _, exists := s.subscriptions[subject]; exists {
		return fmt.Errorf("already subscribed to subject: %s", subject)
	}

	// 创建订阅
	subscription, err := s.conn.Subscribe(subject, func(msg *nats.Msg) {
		s.handleMessage(subject, msg, handler)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to subject %s: %w", subject, err)
	}

	// 存储订阅
	s.subscriptions[subject] = subscription

	s.logger.Info("Subscribed to subject", logging.Fields{
		"subject": subject,
	})
	return nil
}

// SubscribeWithQueue 使用队列组订阅消息
func (s *NATSSubscriber) SubscribeWithQueue(subject, queue string, handler events.EventHandler) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查是否已经订阅
	key := fmt.Sprintf("%s:%s", subject, queue)
	if _, exists := s.subscriptions[key]; exists {
		return fmt.Errorf("already subscribed to subject with queue: %s", key)
	}

	// 创建队列订阅
	subscription, err := s.conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		s.handleMessage(subject, msg, handler)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to subject %s with queue %s: %w", subject, queue, err)
	}

	// 存储订阅
	s.subscriptions[key] = subscription

	s.logger.Info("Subscribed to subject with queue", logging.Fields{
		"subject": subject,
		"queue":   queue,
	})
	return nil
}

// Unsubscribe 取消订阅
func (s *NATSSubscriber) Unsubscribe(subject string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	subscription, exists := s.subscriptions[subject]
	if !exists {
		return fmt.Errorf("not subscribed to subject: %s", subject)
	}

	// 取消订阅
	if err := subscription.Unsubscribe(); err != nil {
		return fmt.Errorf("failed to unsubscribe from subject %s: %w", subject, err)
	}

	// 从订阅列表中移除
	delete(s.subscriptions, subject)

	s.logger.Info("Unsubscribed from subject", logging.Fields{
		"subject": subject,
	})
	return nil
}

// Close 关闭订阅器
func (s *NATSSubscriber) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 取消所有订阅
	for subject, subscription := range s.subscriptions {
		if err := subscription.Unsubscribe(); err != nil {
			s.logger.Error("Failed to unsubscribe", err, logging.Fields{
				"subject": subject,
			})
		}
	}

	// 清空订阅列表
	s.subscriptions = make(map[string]*nats.Subscription)

	// 关闭连接
	if s.conn != nil {
		s.conn.Close()
		s.logger.Info("NATS subscriber connection closed")
	}

	return nil
}

// IsConnected 检查连接状态
func (s *NATSSubscriber) IsConnected() bool {
	return s.conn != nil && s.conn.IsConnected()
}

// GetSubscriptions 获取订阅列表
func (s *NATSSubscriber) GetSubscriptions() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	subjects := make([]string, 0, len(s.subscriptions))
	for subject := range s.subscriptions {
		subjects = append(subjects, subject)
	}
	return subjects
}

// GetStats 获取统计信息
func (s *NATSSubscriber) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})

	if s.conn != nil {
		stats["connected"] = s.conn.IsConnected()
		stats["server_url"] = s.conn.ConnectedUrl()
		stats["server_id"] = s.conn.ConnectedServerId()
		stats["server_version"] = s.conn.ConnectedServerVersion()
	}

	s.mu.RLock()
	stats["subscriptions"] = len(s.subscriptions)
	s.mu.RUnlock()

	return stats
}

// 私有方法

// handleMessage 处理接收到的消息
func (s *NATSSubscriber) handleMessage(subject string, msg *nats.Msg, handler events.EventHandler) {
	s.logger.Debug("Message received", logging.Fields{
		"subject": subject,
		"size":    len(msg.Data),
	})

	// 解析消息
	var event events.Event
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		s.logger.Error("Failed to unmarshal message", err, logging.Fields{
			"subject": subject,
		})
		return
	}

	// 处理事件
	if err := handler.Handle(s.ctx, event); err != nil {
		s.logger.Error("Failed to handle event", err, logging.Fields{
			"subject":    subject,
			"event_type": event.GetEventType(),
		})
		return
	}

	s.logger.Debug("Event handled successfully", logging.Fields{
		"subject":    subject,
		"event_type": event.GetEventType(),
	})
}
