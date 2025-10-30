// Package protocol 统一协议系统
// Author: MMO Server Team
// Created: 2024

package protocol

import (
	"context"
	"fmt"

	// "io"
	"time"
)

// MessageType 消息类型
type MessageType uint16

// 预定义消息类型
const (
	// 系统消息
	MsgTypeHeartbeat MessageType = 1000 + iota
	MsgTypeLogin
	MsgTypeLogout
	MsgTypeError
	MsgTypeNotification

	// 玩家消息
	MsgTypePlayerInfo MessageType = 2000 + iota
	MsgTypePlayerMove
	MsgTypePlayerAction
	MsgTypePlayerChat
	MsgTypePlayerStatus

	// 背包消息
	MsgTypeInventoryList MessageType = 3000 + iota
	MsgTypeInventoryAdd
	MsgTypeInventoryRemove
	MsgTypeInventoryUpdate
	MsgTypeInventoryUse

	// 战斗消息
	MsgTypeBattleStart MessageType = 4000 + iota
	MsgTypeBattleEnd
	MsgTypeBattleAction
	MsgTypeBattleResult
	MsgTypeBattleStatus

	// 社交消息
	MsgTypeFriendList MessageType = 5000 + iota
	MsgTypeFriendAdd
	MsgTypeFriendRemove
	MsgTypeGuildInfo
	MsgTypeGuildJoin

	// 场景消息
	MsgTypeSceneEnter MessageType = 6000 + iota
	MsgTypeSceneLeave
	MsgTypeSceneUpdate
	MsgTypeSceneObject
	MsgTypeSceneNPC

	// 活动消息
	MsgTypeActivityList MessageType = 7000 + iota
	MsgTypeActivityJoin
	MsgTypeActivityReward
	MsgTypeActivityStatus

	// 宠物消息
	MsgTypePetList MessageType = 8000 + iota
	MsgTypePetSummon
	MsgTypePetDismiss
	MsgTypePetUpgrade
	MsgTypePetSkill

	// 建筑消息
	MsgTypeBuildingList MessageType = 9000 + iota
	MsgTypeBuildingBuild
	MsgTypeBuildingUpgrade
	MsgTypeBuildingDestroy
	MsgTypeBuildingCollect
)

// Packet 数据包接口
type Packet interface {
	// GetType 获取消息类型
	GetType() MessageType
	// GetData 获取消息数据
	GetData() []byte
	// SetData 设置消息数据
	SetData(data []byte)
	// GetSize 获取数据包大小
	GetSize() int
	// GetTimestamp 获取时间戳
	GetTimestamp() time.Time
	// SetTimestamp 设置时间戳
	SetTimestamp(timestamp time.Time)
	// GetSequence 获取序列号
	GetSequence() uint32
	// SetSequence 设置序列号
	SetSequence(seq uint32)
	// Validate 验证数据包
	Validate() error
}

// Message 消息接口
type Message interface {
	// GetType 获取消息类型
	GetType() MessageType
	// Marshal 序列化消息
	Marshal() ([]byte, error)
	// Unmarshal 反序列化消息
	Unmarshal(data []byte) error
	// Validate 验证消息
	Validate() error
	// String 字符串表示
	String() string
}

// Codec 编解码器接口
type Codec interface {
	// Encode 编码消息
	Encode(message Message) ([]byte, error)
	// Decode 解码消息
	Decode(data []byte) (Message, error)
	// GetName 获取编解码器名称
	GetName() string
}

// Serializer 序列化器接口
type Serializer interface {
	// Serialize 序列化对象
	Serialize(obj interface{}) ([]byte, error)
	// Deserialize 反序列化对象
	Deserialize(data []byte, obj interface{}) error
	// GetContentType 获取内容类型
	GetContentType() string
}

// Compressor 压缩器接口
type Compressor interface {
	// Compress 压缩数据
	Compress(data []byte) ([]byte, error)
	// Decompress 解压数据
	Decompress(data []byte) ([]byte, error)
	// GetType 获取压缩类型
	GetType() string
}

// Encryptor 加密器接口
type Encryptor interface {
	// Encrypt 加密数据
	Encrypt(data []byte) ([]byte, error)
	// Decrypt 解密数据
	Decrypt(data []byte) ([]byte, error)
	// GetType 获取加密类型
	GetType() string
}

// Handler 消息处理器接口
type Handler interface {
	// Handle 处理消息
	Handle(ctx context.Context, message Message) (Message, error)
	// GetMessageType 获取处理的消息类型
	GetMessageType() MessageType
}

// Middleware 中间件接口
type Middleware interface {
	// Process 处理消息
	Process(ctx context.Context, message Message, next func(context.Context, Message) (Message, error)) (Message, error)
	// GetName 获取中间件名称
	GetName() string
}

// Router 路由器接口
type Router interface {
	// RegisterHandler 注册处理器
	RegisterHandler(msgType MessageType, handler Handler) error
	// UnregisterHandler 注销处理器
	UnregisterHandler(msgType MessageType) error
	// GetHandler 获取处理器
	GetHandler(msgType MessageType) (Handler, bool)
	// Route 路由消息
	Route(ctx context.Context, message Message) (Message, error)
	// AddMiddleware 添加中间件
	AddMiddleware(middleware Middleware)
	// RemoveMiddleware 移除中间件
	RemoveMiddleware(name string)
}

// Connection 连接接口
type Connection interface {
	// ID 获取连接ID
	ID() string
	// RemoteAddr 获取远程地址
	RemoteAddr() string
	// LocalAddr 获取本地地址
	LocalAddr() string
	// Send 发送消息
	Send(message Message) error
	// SendPacket 发送数据包
	SendPacket(packet Packet) error
	// Receive 接收消息
	Receive() (Message, error)
	// ReceivePacket 接收数据包
	ReceivePacket() (Packet, error)
	// Close 关闭连接
	Close() error
	// IsClosed 是否已关闭
	IsClosed() bool
	// GetMetadata 获取元数据
	GetMetadata(key string) interface{}
	// SetMetadata 设置元数据
	SetMetadata(key string, value interface{})
	// GetLastActivity 获取最后活动时间
	GetLastActivity() time.Time
	// UpdateActivity 更新活动时间
	UpdateActivity()
}

// Server 协议服务器接口
type Server interface {
	// Start 启动服务器
	Start(ctx context.Context) error
	// Stop 停止服务器
	Stop(ctx context.Context) error
	// GetAddr 获取监听地址
	GetAddr() string
	// GetConnections 获取所有连接
	GetConnections() []Connection
	// GetConnection 获取指定连接
	GetConnection(id string) (Connection, bool)
	// Broadcast 广播消息
	Broadcast(message Message) error
	// BroadcastToGroup 向组广播消息
	BroadcastToGroup(group string, message Message) error
	// SetRouter 设置路由器
	SetRouter(router Router)
	// GetRouter 获取路由器
	GetRouter() Router
}

// Client 协议客户端接口
type Client interface {
	// Connect 连接服务器
	Connect(ctx context.Context, addr string) error
	// Disconnect 断开连接
	Disconnect() error
	// Send 发送消息
	Send(message Message) error
	// Receive 接收消息
	Receive() (Message, error)
	// IsConnected 是否已连接
	IsConnected() bool
	// GetConnection 获取连接
	GetConnection() Connection
	// SetHandler 设置消息处理器
	SetHandler(handler func(Message) error)
}

// Registry 协议注册表接口
type Registry interface {
	// RegisterMessage 注册消息类型
	RegisterMessage(msgType MessageType, factory func() Message) error
	// UnregisterMessage 注销消息类型
	UnregisterMessage(msgType MessageType) error
	// CreateMessage 创建消息实例
	CreateMessage(msgType MessageType) (Message, error)
	// GetMessageTypes 获取所有消息类型
	GetMessageTypes() []MessageType
	// IsRegistered 检查是否已注册
	IsRegistered(msgType MessageType) bool
}

// Config 协议配置
type Config struct {
	// 协议类型
	Protocol string `yaml:"protocol" json:"protocol"`
	// 监听地址
	Host string `yaml:"host" json:"host"`
	// 监听端口
	Port int `yaml:"port" json:"port"`
	// 缓冲区大小
	BufferSize int `yaml:"buffer_size" json:"buffer_size"`
	// 最大数据包大小
	MaxPacketSize int `yaml:"max_packet_size" json:"max_packet_size"`
	// 读取超时
	ReadTimeout time.Duration `yaml:"read_timeout" json:"read_timeout"`
	// 写入超时
	WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout"`
	// 心跳间隔
	HeartbeatInterval time.Duration `yaml:"heartbeat_interval" json:"heartbeat_interval"`
	// 连接超时
	ConnectionTimeout time.Duration `yaml:"connection_timeout" json:"connection_timeout"`
	// 最大连接数
	MaxConnections int `yaml:"max_connections" json:"max_connections"`
	// 压缩类型
	CompressionType string `yaml:"compression_type" json:"compression_type"`
	// 加密类型
	EncryptionType string `yaml:"encryption_type" json:"encryption_type"`
	// 序列化类型
	SerializationType string `yaml:"serialization_type" json:"serialization_type"`
	// 是否启用压缩
	EnableCompression bool `yaml:"enable_compression" json:"enable_compression"`
	// 是否启用加密
	EnableEncryption bool `yaml:"enable_encryption" json:"enable_encryption"`
	// 是否启用心跳
	EnableHeartbeat bool `yaml:"enable_heartbeat" json:"enable_heartbeat"`
	// TLS配置
	TLS TLSConfig `yaml:"tls" json:"tls"`
}

// TLSConfig TLS配置
type TLSConfig struct {
	// 是否启用TLS
	Enabled bool `yaml:"enabled" json:"enabled"`
	// 证书文件
	CertFile string `yaml:"cert_file" json:"cert_file"`
	// 私钥文件
	KeyFile string `yaml:"key_file" json:"key_file"`
	// CA文件
	CAFile string `yaml:"ca_file" json:"ca_file"`
	// 是否跳过验证
	InsecureSkipVerify bool `yaml:"insecure_skip_verify" json:"insecure_skip_verify"`
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Protocol:          "tcp",
		Host:              "0.0.0.0",
		Port:              8080,
		BufferSize:        4096,
		MaxPacketSize:     65536,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		HeartbeatInterval: 30 * time.Second,
		ConnectionTimeout: 10 * time.Second,
		MaxConnections:    10000,
		CompressionType:   "gzip",
		EncryptionType:    "aes",
		SerializationType: "binary",
		EnableCompression: false,
		EnableEncryption:  false,
		EnableHeartbeat:   true,
		TLS: TLSConfig{
			Enabled:            false,
			InsecureSkipVerify: false,
		},
	}
}

// Factory 协议工厂接口
type Factory interface {
	// CreateServer 创建服务器
	CreateServer(config *Config) (Server, error)
	// CreateClient 创建客户端
	CreateClient(config *Config) (Client, error)
	// CreateCodec 创建编解码器
	CreateCodec(codecType string) (Codec, error)
	// CreateSerializer 创建序列化器
	CreateSerializer(serializerType string) (Serializer, error)
	// CreateCompressor 创建压缩器
	CreateCompressor(compressorType string) (Compressor, error)
	// CreateEncryptor 创建加密器
	CreateEncryptor(encryptorType string) (Encryptor, error)
}

// Manager 协议管理器接口
type Manager interface {
	// GetRegistry 获取注册表
	GetRegistry() Registry
	// GetFactory 获取工厂
	GetFactory() Factory
	// RegisterProtocol 注册协议
	RegisterProtocol(name string, factory Factory) error
	// UnregisterProtocol 注销协议
	UnregisterProtocol(name string) error
	// GetProtocol 获取协议
	GetProtocol(name string) (Factory, bool)
	// CreateServer 创建服务器
	CreateServer(protocol string, config *Config) (Server, error)
	// CreateClient 创建客户端
	CreateClient(protocol string, config *Config) (Client, error)
}

// 预定义错误
var (
	ErrInvalidMessageType    = fmt.Errorf("invalid message type")
	ErrMessageNotRegistered  = fmt.Errorf("message not registered")
	ErrInvalidPacket         = fmt.Errorf("invalid packet")
	ErrPacketTooLarge        = fmt.Errorf("packet too large")
	ErrConnectionClosed      = fmt.Errorf("connection closed")
	ErrConnectionTimeout     = fmt.Errorf("connection timeout")
	ErrHandlerNotFound       = fmt.Errorf("handler not found")
	ErrProtocolNotSupported  = fmt.Errorf("protocol not supported")
	ErrSerializationFailed   = fmt.Errorf("serialization failed")
	ErrDeserializationFailed = fmt.Errorf("deserialization failed")
	ErrCompressionFailed     = fmt.Errorf("compression failed")
	ErrDecompressionFailed   = fmt.Errorf("decompression failed")
	ErrEncryptionFailed      = fmt.Errorf("encryption failed")
	ErrDecryptionFailed      = fmt.Errorf("decryption failed")
)

// 常用常量
const (
	// 协议版本
	ProtocolVersion = 1
	// 魔数
	MagicNumber = 0x12345678
	// 头部大小
	HeaderSize = 16
	// 最小数据包大小
	MinPacketSize = HeaderSize
	// 默认缓冲区大小
	DefaultBufferSize = 4096
	// 默认最大数据包大小
	DefaultMaxPacketSize = 1024 * 1024 // 1MB
)

// PacketHeader 数据包头部
type PacketHeader struct {
	Magic     uint32      // 魔数
	Version   uint16      // 协议版本
	Type      MessageType // 消息类型
	Length    uint32      // 数据长度
	Sequence  uint32      // 序列号
	Timestamp int64       // 时间戳
	Checksum  uint32      // 校验和
}

// BasePacket 基础数据包实现
type BasePacket struct {
	header PacketHeader
	data   []byte
}

// NewBasePacket 创建基础数据包
func NewBasePacket(msgType MessageType, data []byte) *BasePacket {
	return &BasePacket{
		header: PacketHeader{
			Magic:     MagicNumber,
			Version:   ProtocolVersion,
			Type:      msgType,
			Length:    uint32(len(data)),
			Timestamp: time.Now().UnixNano(),
		},
		data: data,
	}
}

func (bp *BasePacket) GetType() MessageType     { return bp.header.Type }
func (bp *BasePacket) GetData() []byte          { return bp.data }
func (bp *BasePacket) SetData(data []byte)      { bp.data = data; bp.header.Length = uint32(len(data)) }
func (bp *BasePacket) GetSize() int             { return HeaderSize + len(bp.data) }
func (bp *BasePacket) GetTimestamp() time.Time  { return time.Unix(0, bp.header.Timestamp) }
func (bp *BasePacket) SetTimestamp(t time.Time) { bp.header.Timestamp = t.UnixNano() }
func (bp *BasePacket) GetSequence() uint32      { return bp.header.Sequence }
func (bp *BasePacket) SetSequence(seq uint32)   { bp.header.Sequence = seq }

func (bp *BasePacket) Validate() error {
	if bp.header.Magic != MagicNumber {
		return fmt.Errorf("invalid magic number: %x", bp.header.Magic)
	}
	if bp.header.Version != ProtocolVersion {
		return fmt.Errorf("unsupported protocol version: %d", bp.header.Version)
	}
	if bp.header.Length != uint32(len(bp.data)) {
		return fmt.Errorf("length mismatch: header=%d, actual=%d", bp.header.Length, len(bp.data))
	}
	return nil
}

// BaseMessage 基础消息实现
type BaseMessage struct {
	msgType MessageType
}

func (bm *BaseMessage) GetType() MessageType { return bm.msgType }
func (bm *BaseMessage) Validate() error      { return nil }
func (bm *BaseMessage) String() string       { return fmt.Sprintf("Message{Type: %d}", bm.msgType) }

// 便捷函数

// IsSystemMessage 检查是否为系统消息
func IsSystemMessage(msgType MessageType) bool {
	return msgType >= 1000 && msgType < 2000
}

// IsPlayerMessage 检查是否为玩家消息
func IsPlayerMessage(msgType MessageType) bool {
	return msgType >= 2000 && msgType < 3000
}

// IsInventoryMessage 检查是否为背包消息
func IsInventoryMessage(msgType MessageType) bool {
	return msgType >= 3000 && msgType < 4000
}

// IsBattleMessage 检查是否为战斗消息
func IsBattleMessage(msgType MessageType) bool {
	return msgType >= 4000 && msgType < 5000
}

// IsSocialMessage 检查是否为社交消息
func IsSocialMessage(msgType MessageType) bool {
	return msgType >= 5000 && msgType < 6000
}

// IsSceneMessage 检查是否为场景消息
func IsSceneMessage(msgType MessageType) bool {
	return msgType >= 6000 && msgType < 7000
}

// IsActivityMessage 检查是否为活动消息
func IsActivityMessage(msgType MessageType) bool {
	return msgType >= 7000 && msgType < 8000
}

// IsPetMessage 检查是否为宠物消息
func IsPetMessage(msgType MessageType) bool {
	return msgType >= 8000 && msgType < 9000
}

// IsBuildingMessage 检查是否为建筑消息
func IsBuildingMessage(msgType MessageType) bool {
	return msgType >= 9000 && msgType < 10000
}

// GetMessageCategory 获取消息分类
func GetMessageCategory(msgType MessageType) string {
	switch {
	case IsSystemMessage(msgType):
		return "system"
	case IsPlayerMessage(msgType):
		return "player"
	case IsInventoryMessage(msgType):
		return "inventory"
	case IsBattleMessage(msgType):
		return "battle"
	case IsSocialMessage(msgType):
		return "social"
	case IsSceneMessage(msgType):
		return "scene"
	case IsActivityMessage(msgType):
		return "activity"
	case IsPetMessage(msgType):
		return "pet"
	case IsBuildingMessage(msgType):
		return "building"
	default:
		return "unknown"
	}
}

// PacketReader 数据包读取器
type PacketReader interface {
	// ReadPacket 读取数据包
	ReadPacket() (Packet, error)
}

// PacketWriter 数据包写入器
type PacketWriter interface {
	// WritePacket 写入数据包
	WritePacket(packet Packet) error
}
