// Package protocol 二进制协议实现
// Author: MMO Server Team
// Created: 2024

package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"time"
)

// BinaryCodec 二进制编解码器
type BinaryCodec struct {
	byteOrder binary.ByteOrder
}

// NewBinaryCodec 创建二进制编解码器
func NewBinaryCodec() *BinaryCodec {
	return &BinaryCodec{
		byteOrder: binary.LittleEndian,
	}
}

// GetName 获取编解码器名称
func (bc *BinaryCodec) GetName() string {
	return "binary"
}

// Encode 编码消息
func (bc *BinaryCodec) Encode(message Message) ([]byte, error) {
	// 序列化消息数据
	msgData, err := message.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	// 创建数据包
	packet := NewBasePacket(message.GetType(), msgData)

	// 编码数据包
	return bc.EncodePacket(packet)
}

// Decode 解码消息
func (bc *BinaryCodec) Decode(data []byte) (Message, error) {
	// 解码数据包
	packet, err := bc.DecodePacket(data)
	if err != nil {
		return nil, err
	}

	// 创建消息实例（这里需要消息注册表支持）
	message := &BinaryMessage{
		BaseMessage: BaseMessage{msgType: packet.GetType()},
		data:        packet.GetData(),
	}

	return message, nil
}

// EncodePacket 编码数据包
func (bc *BinaryCodec) EncodePacket(packet Packet) ([]byte, error) {
	basePacket, ok := packet.(*BasePacket)
	if !ok {
		return nil, fmt.Errorf("unsupported packet type: %T", packet)
	}

	// 计算校验和
	basePacket.header.Checksum = bc.calculateChecksum(basePacket.data)

	// 创建缓冲区
	buf := new(bytes.Buffer)

	// 写入头部
	if err := bc.writeHeader(buf, &basePacket.header); err != nil {
		return nil, fmt.Errorf("failed to write header: %w", err)
	}

	// 写入数据
	if len(basePacket.data) > 0 {
		if _, err := buf.Write(basePacket.data); err != nil {
			return nil, fmt.Errorf("failed to write data: %w", err)
		}
	}

	return buf.Bytes(), nil
}

// DecodePacket 解码数据包
func (bc *BinaryCodec) DecodePacket(data []byte) (Packet, error) {
	if len(data) < HeaderSize {
		return nil, fmt.Errorf("packet too small: %d < %d", len(data), HeaderSize)
	}

	// 读取头部
	buf := bytes.NewReader(data)
	header, err := bc.readHeader(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// 验证头部
	if err := bc.validateHeader(header); err != nil {
		return nil, fmt.Errorf("invalid header: %w", err)
	}

	// 检查数据长度
	expectedSize := HeaderSize + int(header.Length)
	if len(data) != expectedSize {
		return nil, fmt.Errorf("packet size mismatch: expected %d, got %d", expectedSize, len(data))
	}

	// 读取数据
	var msgData []byte
	if header.Length > 0 {
		msgData = make([]byte, header.Length)
		if _, err := buf.Read(msgData); err != nil {
			return nil, fmt.Errorf("failed to read data: %w", err)
		}

		// 验证校验和
		if header.Checksum != bc.calculateChecksum(msgData) {
			return nil, fmt.Errorf("checksum mismatch")
		}
	}

	// 创建数据包
	packet := &BasePacket{
		header: *header,
		data:   msgData,
	}

	return packet, nil
}

// writeHeader 写入头部
func (bc *BinaryCodec) writeHeader(w io.Writer, header *PacketHeader) error {
	if err := binary.Write(w, bc.byteOrder, header.Magic); err != nil {
		return err
	}
	if err := binary.Write(w, bc.byteOrder, header.Version); err != nil {
		return err
	}
	if err := binary.Write(w, bc.byteOrder, header.Type); err != nil {
		return err
	}
	if err := binary.Write(w, bc.byteOrder, header.Length); err != nil {
		return err
	}
	if err := binary.Write(w, bc.byteOrder, header.Sequence); err != nil {
		return err
	}
	if err := binary.Write(w, bc.byteOrder, header.Timestamp); err != nil {
		return err
	}
	if err := binary.Write(w, bc.byteOrder, header.Checksum); err != nil {
		return err
	}
	return nil
}

// readHeader 读取头部
func (bc *BinaryCodec) readHeader(r io.Reader) (*PacketHeader, error) {
	header := &PacketHeader{}

	if err := binary.Read(r, bc.byteOrder, &header.Magic); err != nil {
		return nil, err
	}
	if err := binary.Read(r, bc.byteOrder, &header.Version); err != nil {
		return nil, err
	}
	if err := binary.Read(r, bc.byteOrder, &header.Type); err != nil {
		return nil, err
	}
	if err := binary.Read(r, bc.byteOrder, &header.Length); err != nil {
		return nil, err
	}
	if err := binary.Read(r, bc.byteOrder, &header.Sequence); err != nil {
		return nil, err
	}
	if err := binary.Read(r, bc.byteOrder, &header.Timestamp); err != nil {
		return nil, err
	}
	if err := binary.Read(r, bc.byteOrder, &header.Checksum); err != nil {
		return nil, err
	}

	return header, nil
}

// validateHeader 验证头部
func (bc *BinaryCodec) validateHeader(header *PacketHeader) error {
	if header.Magic != MagicNumber {
		return fmt.Errorf("invalid magic number: %x", header.Magic)
	}
	if header.Version != ProtocolVersion {
		return fmt.Errorf("unsupported protocol version: %d", header.Version)
	}
	if header.Length > DefaultMaxPacketSize {
		return fmt.Errorf("packet too large: %d > %d", header.Length, DefaultMaxPacketSize)
	}
	return nil
}

// calculateChecksum 计算校验和
func (bc *BinaryCodec) calculateChecksum(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

// BinaryMessage 二进制消息实现
type BinaryMessage struct {
	BaseMessage
	data []byte
}

// NewBinaryMessage 创建二进制消息
func NewBinaryMessage(msgType MessageType, data []byte) *BinaryMessage {
	return &BinaryMessage{
		BaseMessage: BaseMessage{msgType: msgType},
		data:        data,
	}
}

// Marshal 序列化消息
func (bm *BinaryMessage) Marshal() ([]byte, error) {
	return bm.data, nil
}

// Unmarshal 反序列化消息
func (bm *BinaryMessage) Unmarshal(data []byte) error {
	bm.data = make([]byte, len(data))
	copy(bm.data, data)
	return nil
}

// GetData 获取数据
func (bm *BinaryMessage) GetData() []byte {
	return bm.data
}

// SetData 设置数据
func (bm *BinaryMessage) SetData(data []byte) {
	bm.data = make([]byte, len(data))
	copy(bm.data, data)
}

// String 字符串表示
func (bm *BinaryMessage) String() string {
	return fmt.Sprintf("BinaryMessage{Type: %d, Size: %d}", bm.msgType, len(bm.data))
}

// BinarySerializer 二进制序列化器
type BinarySerializer struct {
	byteOrder binary.ByteOrder
}

// NewBinarySerializer 创建二进制序列化器
func NewBinarySerializer() *BinarySerializer {
	return &BinarySerializer{
		byteOrder: binary.LittleEndian,
	}
}

// GetContentType 获取内容类型
func (bs *BinarySerializer) GetContentType() string {
	return "application/octet-stream"
}

// Serialize 序列化对象
func (bs *BinarySerializer) Serialize(obj interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, bs.byteOrder, obj); err != nil {
		return nil, fmt.Errorf("failed to serialize object: %w", err)
	}
	return buf.Bytes(), nil
}

// Deserialize 反序列化对象
func (bs *BinarySerializer) Deserialize(data []byte, obj interface{}) error {
	buf := bytes.NewReader(data)
	if err := binary.Read(buf, bs.byteOrder, obj); err != nil {
		return fmt.Errorf("failed to deserialize object: %w", err)
	}
	return nil
}

// BinaryPacketReader 二进制数据包读取器
type BinaryPacketReader struct {
	reader io.Reader
	codec  *BinaryCodec
	buffer []byte
	offset int
}

// NewBinaryPacketReader 创建二进制数据包读取器
func NewBinaryPacketReader(reader io.Reader) *BinaryPacketReader {
	return &BinaryPacketReader{
		reader: reader,
		codec:  NewBinaryCodec(),
		buffer: make([]byte, DefaultBufferSize),
	}
}

// ReadPacket 读取数据包
func (bpr *BinaryPacketReader) ReadPacket() (Packet, error) {
	// 确保有足够的数据读取头部
	if err := bpr.ensureBytes(HeaderSize); err != nil {
		return nil, err
	}

	// 读取头部
	headerData := bpr.buffer[bpr.offset : bpr.offset+HeaderSize]
	header, err := bpr.codec.readHeader(bytes.NewReader(headerData))
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// 验证头部
	if err := bpr.codec.validateHeader(header); err != nil {
		return nil, fmt.Errorf("invalid header: %w", err)
	}

	// 计算总包大小
	totalSize := HeaderSize + int(header.Length)

	// 确保有足够的数据读取整个数据包
	if err := bpr.ensureBytes(totalSize); err != nil {
		return nil, err
	}

	// 读取完整数据包
	packetData := bpr.buffer[bpr.offset : bpr.offset+totalSize]
	packet, err := bpr.codec.DecodePacket(packetData)
	if err != nil {
		return nil, err
	}

	// 更新偏移量
	bpr.offset += totalSize

	return packet, nil
}

// ensureBytes 确保缓冲区有足够的字节
func (bpr *BinaryPacketReader) ensureBytes(needed int) error {
	available := len(bpr.buffer) - bpr.offset
	if available >= needed {
		return nil
	}

	// 移动现有数据到缓冲区开头
	if bpr.offset > 0 {
		copy(bpr.buffer, bpr.buffer[bpr.offset:])
		bpr.buffer = bpr.buffer[:available]
		bpr.offset = 0
	}

	// 扩展缓冲区如果需要
	if cap(bpr.buffer) < needed {
		newBuffer := make([]byte, needed*2)
		copy(newBuffer, bpr.buffer)
		bpr.buffer = newBuffer[:len(bpr.buffer)]
	}

	// 读取更多数据
	for len(bpr.buffer) < needed {
		n, err := bpr.reader.Read(bpr.buffer[len(bpr.buffer):cap(bpr.buffer)])
		if err != nil {
			return err
		}
		bpr.buffer = bpr.buffer[:len(bpr.buffer)+n]
	}

	return nil
}

// BinaryPacketWriter 二进制数据包写入器
type BinaryPacketWriter struct {
	writer io.Writer
	codec  *BinaryCodec
}

// NewBinaryPacketWriter 创建二进制数据包写入器
func NewBinaryPacketWriter(writer io.Writer) *BinaryPacketWriter {
	return &BinaryPacketWriter{
		writer: writer,
		codec:  NewBinaryCodec(),
	}
}

// WritePacket 写入数据包
func (bpw *BinaryPacketWriter) WritePacket(packet Packet) error {
	// 编码数据包
	data, err := bpw.codec.EncodePacket(packet)
	if err != nil {
		return fmt.Errorf("failed to encode packet: %w", err)
	}

	// 写入数据
	if _, err := bpw.writer.Write(data); err != nil {
		return fmt.Errorf("failed to write packet: %w", err)
	}

	return nil
}

// BinaryPacketReadWriter 二进制数据包读写器
type BinaryPacketReadWriter struct {
	*BinaryPacketReader
	*BinaryPacketWriter
	closer io.Closer
}

// NewBinaryPacketReadWriter 创建二进制数据包读写器
func NewBinaryPacketReadWriter(rw io.ReadWriter) *BinaryPacketReadWriter {
	var closer io.Closer
	if c, ok := rw.(io.Closer); ok {
		closer = c
	}

	return &BinaryPacketReadWriter{
		BinaryPacketReader: NewBinaryPacketReader(rw),
		BinaryPacketWriter: NewBinaryPacketWriter(rw),
		closer:             closer,
	}
}

// Close 关闭读写器
func (bprw *BinaryPacketReadWriter) Close() error {
	if bprw.closer != nil {
		return bprw.closer.Close()
	}
	return nil
}

// 预定义的二进制消息类型

// HeartbeatMessage 心跳消息
type HeartbeatMessage struct {
	BaseMessage
	Timestamp int64
}

// NewHeartbeatMessage 创建心跳消息
func NewHeartbeatMessage() *HeartbeatMessage {
	return &HeartbeatMessage{
		BaseMessage: BaseMessage{msgType: MsgTypeHeartbeat},
		Timestamp:   time.Now().UnixNano(),
	}
}

// Marshal 序列化心跳消息
func (hm *HeartbeatMessage) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, hm.Timestamp); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal 反序列化心跳消息
func (hm *HeartbeatMessage) Unmarshal(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, &hm.Timestamp)
}

// String 字符串表示
func (hm *HeartbeatMessage) String() string {
	return fmt.Sprintf("HeartbeatMessage{Timestamp: %d}", hm.Timestamp)
}

// LoginMessage 登录消息
type LoginMessage struct {
	BaseMessage
	Username string
	Password string
	Version  uint32
}

// NewLoginMessage 创建登录消息
func NewLoginMessage(username, password string, version uint32) *LoginMessage {
	return &LoginMessage{
		BaseMessage: BaseMessage{msgType: MsgTypeLogin},
		Username:    username,
		Password:    password,
		Version:     version,
	}
}

// Marshal 序列化登录消息
func (lm *LoginMessage) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	// 写入版本
	if err := binary.Write(buf, binary.LittleEndian, lm.Version); err != nil {
		return nil, err
	}

	// 写入用户名长度和内容
	usernameBytes := []byte(lm.Username)
	if err := binary.Write(buf, binary.LittleEndian, uint16(len(usernameBytes))); err != nil {
		return nil, err
	}
	if _, err := buf.Write(usernameBytes); err != nil {
		return nil, err
	}

	// 写入密码长度和内容
	passwordBytes := []byte(lm.Password)
	if err := binary.Write(buf, binary.LittleEndian, uint16(len(passwordBytes))); err != nil {
		return nil, err
	}
	if _, err := buf.Write(passwordBytes); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Unmarshal 反序列化登录消息
func (lm *LoginMessage) Unmarshal(data []byte) error {
	buf := bytes.NewReader(data)

	// 读取版本
	if err := binary.Read(buf, binary.LittleEndian, &lm.Version); err != nil {
		return err
	}

	// 读取用户名
	var usernameLen uint16
	if err := binary.Read(buf, binary.LittleEndian, &usernameLen); err != nil {
		return err
	}
	usernameBytes := make([]byte, usernameLen)
	if _, err := buf.Read(usernameBytes); err != nil {
		return err
	}
	lm.Username = string(usernameBytes)

	// 读取密码
	var passwordLen uint16
	if err := binary.Read(buf, binary.LittleEndian, &passwordLen); err != nil {
		return err
	}
	passwordBytes := make([]byte, passwordLen)
	if _, err := buf.Read(passwordBytes); err != nil {
		return err
	}
	lm.Password = string(passwordBytes)

	return nil
}

// Validate 验证登录消息
func (lm *LoginMessage) Validate() error {
	if lm.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if lm.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if len(lm.Username) > 32 {
		return fmt.Errorf("username too long: %d > 32", len(lm.Username))
	}
	if len(lm.Password) > 64 {
		return fmt.Errorf("password too long: %d > 64", len(lm.Password))
	}
	return nil
}

// String 字符串表示
func (lm *LoginMessage) String() string {
	return fmt.Sprintf("LoginMessage{Username: %s, Version: %d}", lm.Username, lm.Version)
}

// ErrorMessage 错误消息
type ErrorMessage struct {
	BaseMessage
	Code    uint32
	Message string
}

// NewErrorMessage 创建错误消息
func NewErrorMessage(code uint32, message string) *ErrorMessage {
	return &ErrorMessage{
		BaseMessage: BaseMessage{msgType: MsgTypeError},
		Code:        code,
		Message:     message,
	}
}

// Marshal 序列化错误消息
func (em *ErrorMessage) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	// 写入错误码
	if err := binary.Write(buf, binary.LittleEndian, em.Code); err != nil {
		return nil, err
	}

	// 写入消息长度和内容
	messageBytes := []byte(em.Message)
	if err := binary.Write(buf, binary.LittleEndian, uint16(len(messageBytes))); err != nil {
		return nil, err
	}
	if _, err := buf.Write(messageBytes); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Unmarshal 反序列化错误消息
func (em *ErrorMessage) Unmarshal(data []byte) error {
	buf := bytes.NewReader(data)

	// 读取错误码
	if err := binary.Read(buf, binary.LittleEndian, &em.Code); err != nil {
		return err
	}

	// 读取消息
	var messageLen uint16
	if err := binary.Read(buf, binary.LittleEndian, &messageLen); err != nil {
		return err
	}
	messageBytes := make([]byte, messageLen)
	if _, err := buf.Read(messageBytes); err != nil {
		return err
	}
	em.Message = string(messageBytes)

	return nil
}

// String 字符串表示
func (em *ErrorMessage) String() string {
	return fmt.Sprintf("ErrorMessage{Code: %d, Message: %s}", em.Code, em.Message)
}