package protocol

import (
	"encoding/binary"
	"fmt"
	"time"
)

// Message 协议消息结构
type Message struct {
	Header  MessageHeader `json:"header"`
	Payload []byte        `json:"payload"`
}

// MessageHeader 消息头
type MessageHeader struct {
	MessageType uint16 `json:"message_type"`
	MessageID   uint32 `json:"message_id"`
	Timestamp   int64  `json:"timestamp"`
	Length      uint32 `json:"length"`
}

// NewMessage 创建新消息
func NewMessage(messageType uint16, payload []byte) *Message {
	return &Message{
		Header: MessageHeader{
			MessageType: messageType,
			MessageID:   uint32(time.Now().UnixNano()),
			Timestamp:   time.Now().Unix(),
			Length:      uint32(len(payload)),
		},
		Payload: payload,
	}
}

// Serialize 序列化消息
func (m *Message) Serialize() ([]byte, error) {
	// 计算总长度
	totalLength := 2 + 4 + 8 + 4 + len(m.Payload) // header + payload

	// 创建缓冲区
	buf := make([]byte, totalLength)
	offset := 0

	// 写入消息类型 (2 bytes)
	binary.LittleEndian.PutUint16(buf[offset:], m.Header.MessageType)
	offset += 2

	// 写入消息ID (4 bytes)
	binary.LittleEndian.PutUint32(buf[offset:], m.Header.MessageID)
	offset += 4

	// 写入时间戳 (8 bytes)
	binary.LittleEndian.PutUint64(buf[offset:], uint64(m.Header.Timestamp))
	offset += 8

	// 写入长度 (4 bytes)
	binary.LittleEndian.PutUint32(buf[offset:], m.Header.Length)
	offset += 4

	// 写入载荷
	copy(buf[offset:], m.Payload)

	return buf, nil
}

// Deserialize 反序列化消息
func Deserialize(data []byte) (*Message, error) {
	if len(data) < 18 { // 最小消息头长度
		return nil, fmt.Errorf("消息数据太短")
	}

	offset := 0

	// 读取消息类型
	messageType := binary.LittleEndian.Uint16(data[offset:])
	offset += 2

	// 读取消息ID
	messageID := binary.LittleEndian.Uint32(data[offset:])
	offset += 4

	// 读取时间戳
	timestamp := int64(binary.LittleEndian.Uint64(data[offset:]))
	offset += 8

	// 读取长度
	length := binary.LittleEndian.Uint32(data[offset:])
	offset += 4

	// 读取载荷
	payload := make([]byte, length)
	copy(payload, data[offset:offset+int(length)])

	return &Message{
		Header: MessageHeader{
			MessageType: messageType,
			MessageID:   messageID,
			Timestamp:   timestamp,
			Length:      length,
		},
		Payload: payload,
	}, nil
}
