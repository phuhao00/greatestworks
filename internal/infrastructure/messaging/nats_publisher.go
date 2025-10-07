package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"greatestworks/internal/infrastructure/logging"

	"github.com/nats-io/nats.go"
)

// NATSPublisher NATS消息发布器
type NATSPublisher struct {
	conn   *nats.Conn
	logger logging.Logger
	config *PublisherConfig
}

// PublisherConfig 发布器配置
type PublisherConfig struct {
	URL          string        `json:"url" yaml:"url"`
	ClusterID    string        `json:"cluster_id" yaml:"cluster_id"`
	ClientID     string        `json:"client_id" yaml:"client_id"`
	Timeout      time.Duration `json:"timeout" yaml:"timeout"`
	MaxReconn    int           `json:"max_reconn" yaml:"max_reconn"`
	ReconnWait   time.Duration `json:"reconn_wait" yaml:"reconn_wait"`
	PingInterval time.Duration `json:"ping_interval" yaml:"ping_interval"`
	MaxPingsOut  int           `json:"max_pings_out" yaml:"max_pings_out"`
}

// NewNATSPublisher 创建NATS发布器
func NewNATSPublisher(config *PublisherConfig, logger logging.Logger) (*NATSPublisher, error) {
	if config == nil {
		config = &PublisherConfig{
			URL:          "nats://localhost:4222",
			ClusterID:    "test-cluster",
			ClientID:     "publisher",
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

	publisher := &NATSPublisher{
		conn:   conn,
		logger: logger,
		config: config,
	}

	logger.Info("NATS publisher created", logging.Fields{
		"url":       config.URL,
		"client_id": config.ClientID,
	})
	return publisher, nil
}

// Publish 发布消息
func (p *NATSPublisher) Publish(ctx context.Context, subject string, data interface{}) error {
	// 序列化数据
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// 发布消息
	if err := p.conn.Publish(subject, payload); err != nil {
		p.logger.Error("Failed to publish message", err, logging.Fields{
			"subject": subject,
		})
		return fmt.Errorf("failed to publish message: %w", err)
	}

	p.logger.Debug("Message published", logging.Fields{
		"subject": subject,
		"size":    len(payload),
	})
	return nil
}

// PublishAsync 异步发布消息
func (p *NATSPublisher) PublishAsync(ctx context.Context, subject string, data interface{}, callback func(error)) error {
	// 序列化数据
	payload, err := json.Marshal(data)
	if err != nil {
		if callback != nil {
			callback(err)
		}
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// 异步发布消息
	if err := p.conn.Publish(subject, payload); err != nil {
		p.logger.Error("Failed to publish message async", err, logging.Fields{
			"subject": subject,
		})
		return fmt.Errorf("failed to publish message async: %w", err)
	}

	p.logger.Debug("Message published async", logging.Fields{
		"subject": subject,
		"size":    len(payload),
	})
	return nil
}

// PublishWithReply 发布带回复的消息
func (p *NATSPublisher) PublishWithReply(ctx context.Context, subject string, data interface{}, replySubject string) error {
	// 序列化数据
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// 发布消息
	if err := p.conn.PublishRequest(subject, replySubject, payload); err != nil {
		p.logger.Error("Failed to publish message with reply", err, logging.Fields{
			"subject": subject,
			"reply":   replySubject,
		})
		return fmt.Errorf("failed to publish message with reply: %w", err)
	}

	p.logger.Debug("Message published with reply", logging.Fields{
		"subject": subject,
		"reply":   replySubject,
		"size":    len(payload),
	})
	return nil
}

// Request 发送请求并等待回复
func (p *NATSPublisher) Request(ctx context.Context, subject string, data interface{}, timeout time.Duration) ([]byte, error) {
	// 序列化数据
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	// 发送请求
	msg, err := p.conn.Request(subject, payload, timeout)
	if err != nil {
		p.logger.Error("Failed to send request", err, logging.Fields{
			"subject": subject,
		})
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	p.logger.Debug("Request sent and reply received", logging.Fields{
		"subject":    subject,
		"reply_size": len(msg.Data),
	})
	return msg.Data, nil
}

// Close 关闭连接
func (p *NATSPublisher) Close() error {
	if p.conn != nil {
		p.conn.Close()
		p.logger.Info("NATS publisher connection closed")
	}
	return nil
}

// IsConnected 检查连接状态
func (p *NATSPublisher) IsConnected() bool {
	return p.conn != nil && p.conn.IsConnected()
}

// GetStats 获取统计信息
func (p *NATSPublisher) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})

	if p.conn != nil {
		stats["connected"] = p.conn.IsConnected()
		stats["server_url"] = p.conn.ConnectedUrl()
		stats["server_id"] = p.conn.ConnectedServerId()
		stats["server_version"] = p.conn.ConnectedServerVersion()
	}

	return stats
}
