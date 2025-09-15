package chat

import (
	"strings"
	"time"
)

// Message 聊天消息实体
type Message struct {
	ID        string
	ChannelID string
	SenderID  string
	Content   string
	Type      MessageType
	Timestamp time.Time
	Metadata  map[string]interface{}
}

// MessageType 消息类型
type MessageType int

const (
	MessageTypeText   MessageType = iota // 文本消息
	MessageTypeImage                     // 图片消息
	MessageTypeSystem                    // 系统消息
	MessageTypeEmoji                     // 表情消息
	MessageTypeItem                      // 物品链接消息
)

// NewMessage 创建新消息
func NewMessage(channelID, senderID, content string, msgType MessageType) *Message {
	return &Message{
		ID:        generateMessageID(),
		ChannelID: channelID,
		SenderID:  senderID,
		Content:   content,
		Type:      msgType,
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
	}
}

// NewSystemMessage 创建系统消息
func NewSystemMessage(channelID, content string) *Message {
	return &Message{
		ID:        generateMessageID(),
		ChannelID: channelID,
		SenderID:  "system",
		Content:   content,
		Type:      MessageTypeSystem,
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
	}
}

// Validate 验证消息
func (m *Message) Validate() error {
	if m.ChannelID == "" {
		return ErrInvalidChannelID
	}
	
	if m.SenderID == "" {
		return ErrInvalidSenderID
	}
	
	if strings.TrimSpace(m.Content) == "" {
		return ErrEmptyContent
	}
	
	if len(m.Content) > MaxMessageLength {
		return ErrMessageTooLong
	}
	
	return nil
}

// SetMetadata 设置元数据
func (m *Message) SetMetadata(key string, value interface{}) {
	m.Metadata[key] = value
}

// GetMetadata 获取元数据
func (m *Message) GetMetadata(key string) (interface{}, bool) {
	value, exists := m.Metadata[key]
	return value, exists
}

// IsSystemMessage 是否为系统消息
func (m *Message) IsSystemMessage() bool {
	return m.Type == MessageTypeSystem
}

// GetAge 获取消息年龄
func (m *Message) GetAge() time.Duration {
	return time.Since(m.Timestamp)
}

const (
	MaxMessageLength = 500 // 最大消息长度
)

// generateMessageID 生成消息ID
func generateMessageID() string {
	// 简单的ID生成，实际项目中应使用更robust的方案
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}