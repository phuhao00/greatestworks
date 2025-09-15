// Package network 消息编解码器
// Author: MMO Server Team
// Created: 2024

package network

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sync"
)

// Codec 编解码器接口
type Codec interface {
	// Encode 编码消息
	Encode(msg interface{}) ([]byte, error)
	
	// Decode 解码消息
	Decode(data []byte, msg interface{}) error
	
	// Name 编解码器名称
	Name() string
}

// JSONCodec JSON编解码器
type JSONCodec struct{}

// NewJSONCodec 创建JSON编解码器
func NewJSONCodec() *JSONCodec {
	return &JSONCodec{}
}

// Encode JSON编码
func (c *JSONCodec) Encode(msg interface{}) ([]byte, error) {
	return json.Marshal(msg)
}

// Decode JSON解码
func (c *JSONCodec) Decode(data []byte, msg interface{}) error {
	return json.Unmarshal(data, msg)
}

// Name 编解码器名称
func (c *JSONCodec) Name() string {
	return "json"
}

// BinaryCodec 二进制编解码器
type BinaryCodec struct {
	typeRegistry map[string]reflect.Type
	mutex        sync.RWMutex
}

// NewBinaryCodec 创建二进制编解码器
func NewBinaryCodec() *BinaryCodec {
	return &BinaryCodec{
		typeRegistry: make(map[string]reflect.Type),
	}
}

// RegisterType 注册类型
func (c *BinaryCodec) RegisterType(name string, t reflect.Type) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.typeRegistry[name] = t
}

// Encode 二进制编码
func (c *BinaryCodec) Encode(msg interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	
	// 写入类型名
	typeName := reflect.TypeOf(msg).Name()
	if err := c.writeString(buf, typeName); err != nil {
		return nil, err
	}
	
	// 编码数据
	if err := c.encodeValue(buf, reflect.ValueOf(msg)); err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}

// Decode 二进制解码
func (c *BinaryCodec) Decode(data []byte, msg interface{}) error {
	buf := bytes.NewReader(data)
	
	// 读取类型名
	typeName, err := c.readString(buf)
	if err != nil {
		return err
	}
	
	// 验证类型
	c.mutex.RLock()
	expectedType, exists := c.typeRegistry[typeName]
	c.mutex.RUnlock()
	
	if !exists {
		return fmt.Errorf("unknown type: %s", typeName)
	}
	
	msgValue := reflect.ValueOf(msg)
	if msgValue.Kind() != reflect.Ptr {
		return fmt.Errorf("msg must be a pointer")
	}
	
	msgElem := msgValue.Elem()
	if msgElem.Type() != expectedType {
		return fmt.Errorf("type mismatch: expected %s, got %s", expectedType, msgElem.Type())
	}
	
	// 解码数据
	return c.decodeValue(buf, msgElem)
}

// Name 编解码器名称
func (c *BinaryCodec) Name() string {
	return "binary"
}

// encodeValue 编码值
func (c *BinaryCodec) encodeValue(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return binary.Write(buf, binary.BigEndian, uint8(1))
		}
		return binary.Write(buf, binary.BigEndian, uint8(0))
		
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return binary.Write(buf, binary.BigEndian, v.Int())
		
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return binary.Write(buf, binary.BigEndian, v.Uint())
		
	case reflect.Float32, reflect.Float64:
		return binary.Write(buf, binary.BigEndian, v.Float())
		
	case reflect.String:
		return c.writeString(buf, v.String())
		
	case reflect.Slice:
		// 写入长度
		if err := binary.Write(buf, binary.BigEndian, uint32(v.Len())); err != nil {
			return err
		}
		// 写入元素
		for i := 0; i < v.Len(); i++ {
			if err := c.encodeValue(buf, v.Index(i)); err != nil {
				return err
			}
		}
		return nil
		
	case reflect.Map:
		// 写入长度
		if err := binary.Write(buf, binary.BigEndian, uint32(v.Len())); err != nil {
			return err
		}
		// 写入键值对
		for _, key := range v.MapKeys() {
			if err := c.encodeValue(buf, key); err != nil {
				return err
			}
			if err := c.encodeValue(buf, v.MapIndex(key)); err != nil {
				return err
			}
		}
		return nil
		
	case reflect.Struct:
		// 写入字段数量
		if err := binary.Write(buf, binary.BigEndian, uint32(v.NumField())); err != nil {
			return err
		}
		// 写入字段
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if field.CanInterface() {
				if err := c.encodeValue(buf, field); err != nil {
					return err
				}
			}
		}
		return nil
		
	case reflect.Ptr:
		if v.IsNil() {
			return binary.Write(buf, binary.BigEndian, uint8(0))
		}
		if err := binary.Write(buf, binary.BigEndian, uint8(1)); err != nil {
			return err
		}
		return c.encodeValue(buf, v.Elem())
		
	default:
		return fmt.Errorf("unsupported type: %s", v.Kind())
	}
}

// decodeValue 解码值
func (c *BinaryCodec) decodeValue(buf *bytes.Reader, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Bool:
		var b uint8
		if err := binary.Read(buf, binary.BigEndian, &b); err != nil {
			return err
		}
		v.SetBool(b != 0)
		return nil
		
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var val int64
		if err := binary.Read(buf, binary.BigEndian, &val); err != nil {
			return err
		}
		v.SetInt(val)
		return nil
		
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var val uint64
		if err := binary.Read(buf, binary.BigEndian, &val); err != nil {
			return err
		}
		v.SetUint(val)
		return nil
		
	case reflect.Float32, reflect.Float64:
		var val float64
		if err := binary.Read(buf, binary.BigEndian, &val); err != nil {
			return err
		}
		v.SetFloat(val)
		return nil
		
	case reflect.String:
		str, err := c.readString(buf)
		if err != nil {
			return err
		}
		v.SetString(str)
		return nil
		
	case reflect.Slice:
		// 读取长度
		var length uint32
		if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
			return err
		}
		// 创建切片
		slice := reflect.MakeSlice(v.Type(), int(length), int(length))
		// 读取元素
		for i := 0; i < int(length); i++ {
			if err := c.decodeValue(buf, slice.Index(i)); err != nil {
				return err
			}
		}
		v.Set(slice)
		return nil
		
	case reflect.Map:
		// 读取长度
		var length uint32
		if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
			return err
		}
		// 创建映射
		mapValue := reflect.MakeMap(v.Type())
		// 读取键值对
		for i := 0; i < int(length); i++ {
			key := reflect.New(v.Type().Key()).Elem()
			val := reflect.New(v.Type().Elem()).Elem()
			if err := c.decodeValue(buf, key); err != nil {
				return err
			}
			if err := c.decodeValue(buf, val); err != nil {
				return err
			}
			mapValue.SetMapIndex(key, val)
		}
		v.Set(mapValue)
		return nil
		
	case reflect.Struct:
		// 读取字段数量
		var fieldCount uint32
		if err := binary.Read(buf, binary.BigEndian, &fieldCount); err != nil {
			return err
		}
		// 读取字段
		for i := 0; i < int(fieldCount) && i < v.NumField(); i++ {
			field := v.Field(i)
			if field.CanSet() {
				if err := c.decodeValue(buf, field); err != nil {
					return err
				}
			}
		}
		return nil
		
	case reflect.Ptr:
		var isNil uint8
		if err := binary.Read(buf, binary.BigEndian, &isNil); err != nil {
			return err
		}
		if isNil == 0 {
			v.Set(reflect.Zero(v.Type()))
			return nil
		}
		// 创建新值
		newVal := reflect.New(v.Type().Elem())
		if err := c.decodeValue(buf, newVal.Elem()); err != nil {
			return err
		}
		v.Set(newVal)
		return nil
		
	default:
		return fmt.Errorf("unsupported type: %s", v.Kind())
	}
}

// writeString 写入字符串
func (c *BinaryCodec) writeString(buf *bytes.Buffer, s string) error {
	data := []byte(s)
	if err := binary.Write(buf, binary.BigEndian, uint32(len(data))); err != nil {
		return err
	}
	_, err := buf.Write(data)
	return err
}

// readString 读取字符串
func (c *BinaryCodec) readString(buf *bytes.Reader) (string, error) {
	var length uint32
	if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
		return "", err
	}
	
	data := make([]byte, length)
	if _, err := io.ReadFull(buf, data); err != nil {
		return "", err
	}
	
	return string(data), nil
}

// CompressCodec 压缩编解码器
type CompressCodec struct {
	baseCodec Codec
}

// NewCompressCodec 创建压缩编解码器
func NewCompressCodec(baseCodec Codec) *CompressCodec {
	return &CompressCodec{
		baseCodec: baseCodec,
	}
}

// Encode 压缩编码
func (c *CompressCodec) Encode(msg interface{}) ([]byte, error) {
	// 先用基础编解码器编码
	data, err := c.baseCodec.Encode(msg)
	if err != nil {
		return nil, err
	}
	
	// 压缩数据
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}

// Decode 解压解码
func (c *CompressCodec) Decode(data []byte, msg interface{}) error {
	// 解压数据
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer reader.Close()
	
	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	
	// 用基础编解码器解码
	return c.baseCodec.Decode(decompressed, msg)
}

// Name 编解码器名称
func (c *CompressCodec) Name() string {
	return "compress_" + c.baseCodec.Name()
}

// EncryptCodec 加密编解码器
type EncryptCodec struct {
	baseCodec Codec
	key       []byte
}

// NewEncryptCodec 创建加密编解码器
func NewEncryptCodec(baseCodec Codec, key []byte) *EncryptCodec {
	return &EncryptCodec{
		baseCodec: baseCodec,
		key:       key,
	}
}

// Encode 加密编码
func (c *EncryptCodec) Encode(msg interface{}) ([]byte, error) {
	// 先用基础编解码器编码
	data, err := c.baseCodec.Encode(msg)
	if err != nil {
		return nil, err
	}
	
	// 创建AES加密器
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}
	
	// 生成随机IV
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	
	// 加密数据
	stream := cipher.NewCFBEncrypter(block, iv)
	encrypted := make([]byte, len(data))
	stream.XORKeyStream(encrypted, data)
	
	// 返回IV+加密数据
	result := make([]byte, len(iv)+len(encrypted))
	copy(result, iv)
	copy(result[len(iv):], encrypted)
	
	return result, nil
}

// Decode 解密解码
func (c *EncryptCodec) Decode(data []byte, msg interface{}) error {
	if len(data) < aes.BlockSize {
		return fmt.Errorf("encrypted data too short")
	}
	
	// 提取IV和加密数据
	iv := data[:aes.BlockSize]
	encrypted := data[aes.BlockSize:]
	
	// 创建AES解密器
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return err
	}
	
	// 解密数据
	stream := cipher.NewCFBDecrypter(block, iv)
	decrypted := make([]byte, len(encrypted))
	stream.XORKeyStream(decrypted, encrypted)
	
	// 用基础编解码器解码
	return c.baseCodec.Decode(decrypted, msg)
}

// Name 编解码器名称
func (c *EncryptCodec) Name() string {
	return "encrypt_" + c.baseCodec.Name()
}

// CodecManager 编解码器管理器
type CodecManager struct {
	codecs map[string]Codec
	mutex  sync.RWMutex
}

// NewCodecManager 创建编解码器管理器
func NewCodecManager() *CodecManager {
	m := &CodecManager{
		codecs: make(map[string]Codec),
	}
	
	// 注册默认编解码器
	m.RegisterCodec(NewJSONCodec())
	m.RegisterCodec(NewBinaryCodec())
	
	return m
}

// RegisterCodec 注册编解码器
func (m *CodecManager) RegisterCodec(codec Codec) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.codecs[codec.Name()] = codec
}

// GetCodec 获取编解码器
func (m *CodecManager) GetCodec(name string) (Codec, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	codec, exists := m.codecs[name]
	return codec, exists
}

// ListCodecs 列出所有编解码器
func (m *CodecManager) ListCodecs() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	names := make([]string, 0, len(m.codecs))
	for name := range m.codecs {
		names = append(names, name)
	}
	return names
}

// MessageProcessor 消息处理器
type MessageProcessor struct {
	codecManager *CodecManager
	defaultCodec string
}

// NewMessageProcessor 创建消息处理器
func NewMessageProcessor() *MessageProcessor {
	return &MessageProcessor{
		codecManager: NewCodecManager(),
		defaultCodec: "json",
	}
}

// SetDefaultCodec 设置默认编解码器
func (p *MessageProcessor) SetDefaultCodec(name string) {
	p.defaultCodec = name
}

// RegisterCodec 注册编解码器
func (p *MessageProcessor) RegisterCodec(codec Codec) {
	p.codecManager.RegisterCodec(codec)
}

// EncodeMessage 编码消息
func (p *MessageProcessor) EncodeMessage(msgType MessageType, data interface{}, codecName string) (*Message, error) {
	if codecName == "" {
		codecName = p.defaultCodec
	}
	
	codec, exists := p.codecManager.GetCodec(codecName)
	if !exists {
		return nil, fmt.Errorf("codec %s not found", codecName)
	}
	
	body, err := codec.Encode(data)
	if err != nil {
		return nil, err
	}
	
	msg := &Message{
		Header: MessageHeader{
			Type: msgType,
		},
		Body: body,
	}
	
	return msg, nil
}

// DecodeMessage 解码消息
func (p *MessageProcessor) DecodeMessage(msg *Message, data interface{}, codecName string) error {
	if codecName == "" {
		codecName = p.defaultCodec
	}
	
	codec, exists := p.codecManager.GetCodec(codecName)
	if !exists {
		return fmt.Errorf("codec %s not found", codecName)
	}
	
	return codec.Decode(msg.Body, data)
}

// ProcessMessage 处理消息（自动检测编解码器）
func (p *MessageProcessor) ProcessMessage(msg *Message, data interface{}) error {
	// 尝试不同的编解码器
	codecs := p.codecManager.ListCodecs()
	
	for _, codecName := range codecs {
		if err := p.DecodeMessage(msg, data, codecName); err == nil {
			return nil
		}
	}
	
	return fmt.Errorf("failed to decode message with any codec")
}