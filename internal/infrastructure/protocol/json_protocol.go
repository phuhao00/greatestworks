// Package protocol JSON协议实现
// Author: MMO Server Team
// Created: 2024

package protocol

import (
	"encoding/json"
	"fmt"
	"time"
)

// JSONCodec JSON编解码器
type JSONCodec struct {
	prettyPrint bool
}

// NewJSONCodec 创建JSON编解码器
func NewJSONCodec(prettyPrint bool) *JSONCodec {
	return &JSONCodec{
		prettyPrint: prettyPrint,
	}
}

// GetName 获取编解码器名称
func (jc *JSONCodec) GetName() string {
	return "json"
}

// Encode 编码消息
func (jc *JSONCodec) Encode(message Message) ([]byte, error) {
	// 创建JSON包装器
	wrapper := &JSONMessageWrapper{
		Type:      message.GetType(),
		Timestamp: time.Now().UnixNano(),
		Data:      message,
	}

	// 序列化为JSON
	if jc.prettyPrint {
		return json.MarshalIndent(wrapper, "", "  ")
	}
	return json.Marshal(wrapper)
}

// Decode 解码消息
func (jc *JSONCodec) Decode(data []byte) (Message, error) {
	// 解析JSON包装器
	var wrapper JSONMessageWrapper
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// 创建具体的消息实例
	message, err := jc.createMessage(wrapper.Type)
	if err != nil {
		return nil, err
	}

	// 反序列化消息数据
	if wrapper.Data != nil {
		dataBytes, err := json.Marshal(wrapper.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal message data: %w", err)
		}
		if err := message.Unmarshal(dataBytes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message: %w", err)
		}
	}

	return message, nil
}

// createMessage 创建消息实例
func (jc *JSONCodec) createMessage(msgType MessageType) (Message, error) {
	// 根据消息类型创建对应的消息实例
	switch msgType {
	case MsgTypeHeartbeat:
		return &JSONHeartbeatMessage{BaseMessage: BaseMessage{msgType: msgType}}, nil
	case MsgTypeLogin:
		return &JSONLoginMessage{BaseMessage: BaseMessage{msgType: msgType}}, nil
	case MsgTypeLogout:
		return &JSONLogoutMessage{BaseMessage: BaseMessage{msgType: msgType}}, nil
	case MsgTypeError:
		return &JSONErrorMessage{BaseMessage: BaseMessage{msgType: msgType}}, nil
	case MsgTypePlayerInfo:
		return &JSONPlayerInfoMessage{BaseMessage: BaseMessage{msgType: msgType}}, nil
	case MsgTypePlayerMove:
		return &JSONPlayerMoveMessage{BaseMessage: BaseMessage{msgType: msgType}}, nil
	case MsgTypePlayerChat:
		return &JSONPlayerChatMessage{BaseMessage: BaseMessage{msgType: msgType}}, nil
	default:
		// 默认创建通用JSON消息
		return &JSONGenericMessage{BaseMessage: BaseMessage{msgType: msgType}}, nil
	}
}

// JSONMessageWrapper JSON消息包装器
type JSONMessageWrapper struct {
	Type      MessageType `json:"type"`
	Timestamp int64       `json:"timestamp"`
	Sequence  uint32      `json:"sequence,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

// JSONSerializer JSON序列化器
type JSONSerializer struct {
	prettyPrint bool
}

// NewJSONSerializer 创建JSON序列化器
func NewJSONSerializer(prettyPrint bool) *JSONSerializer {
	return &JSONSerializer{
		prettyPrint: prettyPrint,
	}
}

// GetContentType 获取内容类型
func (js *JSONSerializer) GetContentType() string {
	return "application/json"
}

// Serialize 序列化对象
func (js *JSONSerializer) Serialize(obj interface{}) ([]byte, error) {
	if js.prettyPrint {
		return json.MarshalIndent(obj, "", "  ")
	}
	return json.Marshal(obj)
}

// Deserialize 反序列化对象
func (js *JSONSerializer) Deserialize(data []byte, obj interface{}) error {
	return json.Unmarshal(data, obj)
}

// JSONGenericMessage 通用JSON消息
type JSONGenericMessage struct {
	BaseMessage
	Data map[string]interface{} `json:"data"`
}

// NewJSONGenericMessage 创建通用JSON消息
func NewJSONGenericMessage(msgType MessageType, data map[string]interface{}) *JSONGenericMessage {
	return &JSONGenericMessage{
		BaseMessage: BaseMessage{msgType: msgType},
		Data:        data,
	}
}

// Marshal 序列化消息
func (jgm *JSONGenericMessage) Marshal() ([]byte, error) {
	return json.Marshal(jgm.Data)
}

// Unmarshal 反序列化消息
func (jgm *JSONGenericMessage) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &jgm.Data)
}

// String 字符串表示
func (jgm *JSONGenericMessage) String() string {
	return fmt.Sprintf("JSONGenericMessage{Type: %d, Data: %v}", jgm.msgType, jgm.Data)
}

// JSONHeartbeatMessage JSON心跳消息
type JSONHeartbeatMessage struct {
	BaseMessage
	Timestamp int64  `json:"timestamp"`
	ClientID  string `json:"client_id,omitempty"`
}

// NewJSONHeartbeatMessage 创建JSON心跳消息
func NewJSONHeartbeatMessage(clientID string) *JSONHeartbeatMessage {
	return &JSONHeartbeatMessage{
		BaseMessage: BaseMessage{msgType: MsgTypeHeartbeat},
		Timestamp:   time.Now().UnixNano(),
		ClientID:    clientID,
	}
}

// Marshal 序列化消息
func (jhm *JSONHeartbeatMessage) Marshal() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"timestamp": jhm.Timestamp,
		"client_id": jhm.ClientID,
	})
}

// Unmarshal 反序列化消息
func (jhm *JSONHeartbeatMessage) Unmarshal(data []byte) error {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if timestamp, ok := obj["timestamp"].(float64); ok {
		jhm.Timestamp = int64(timestamp)
	}
	if clientID, ok := obj["client_id"].(string); ok {
		jhm.ClientID = clientID
	}

	return nil
}

// String 字符串表示
func (jhm *JSONHeartbeatMessage) String() string {
	return fmt.Sprintf("JSONHeartbeatMessage{Timestamp: %d, ClientID: %s}", jhm.Timestamp, jhm.ClientID)
}

// JSONLoginMessage JSON登录消息
type JSONLoginMessage struct {
	BaseMessage
	Username string            `json:"username"`
	Password string            `json:"password"`
	Version  uint32            `json:"version"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// NewJSONLoginMessage 创建JSON登录消息
func NewJSONLoginMessage(username, password string, version uint32) *JSONLoginMessage {
	return &JSONLoginMessage{
		BaseMessage: BaseMessage{msgType: MsgTypeLogin},
		Username:    username,
		Password:    password,
		Version:     version,
		Metadata:    make(map[string]string),
	}
}

// Marshal 序列化消息
func (jlm *JSONLoginMessage) Marshal() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"username": jlm.Username,
		"password": jlm.Password,
		"version":  jlm.Version,
		"metadata": jlm.Metadata,
	})
}

// Unmarshal 反序列化消息
func (jlm *JSONLoginMessage) Unmarshal(data []byte) error {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if username, ok := obj["username"].(string); ok {
		jlm.Username = username
	}
	if password, ok := obj["password"].(string); ok {
		jlm.Password = password
	}
	if version, ok := obj["version"].(float64); ok {
		jlm.Version = uint32(version)
	}
	if metadata, ok := obj["metadata"].(map[string]interface{}); ok {
		jlm.Metadata = make(map[string]string)
		for k, v := range metadata {
			if str, ok := v.(string); ok {
				jlm.Metadata[k] = str
			}
		}
	}

	return nil
}

// Validate 验证登录消息
func (jlm *JSONLoginMessage) Validate() error {
	if jlm.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if jlm.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if len(jlm.Username) > 32 {
		return fmt.Errorf("username too long: %d > 32", len(jlm.Username))
	}
	if len(jlm.Password) > 64 {
		return fmt.Errorf("password too long: %d > 64", len(jlm.Password))
	}
	return nil
}

// String 字符串表示
func (jlm *JSONLoginMessage) String() string {
	return fmt.Sprintf("JSONLoginMessage{Username: %s, Version: %d}", jlm.Username, jlm.Version)
}

// JSONLogoutMessage JSON登出消息
type JSONLogoutMessage struct {
	BaseMessage
	Reason string `json:"reason,omitempty"`
}

// NewJSONLogoutMessage 创建JSON登出消息
func NewJSONLogoutMessage(reason string) *JSONLogoutMessage {
	return &JSONLogoutMessage{
		BaseMessage: BaseMessage{msgType: MsgTypeLogout},
		Reason:      reason,
	}
}

// Marshal 序列化消息
func (jlom *JSONLogoutMessage) Marshal() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"reason": jlom.Reason,
	})
}

// Unmarshal 反序列化消息
func (jlom *JSONLogoutMessage) Unmarshal(data []byte) error {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if reason, ok := obj["reason"].(string); ok {
		jlom.Reason = reason
	}

	return nil
}

// String 字符串表示
func (jlom *JSONLogoutMessage) String() string {
	return fmt.Sprintf("JSONLogoutMessage{Reason: %s}", jlom.Reason)
}

// JSONErrorMessage JSON错误消息
type JSONErrorMessage struct {
	BaseMessage
	Code    uint32 `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// NewJSONErrorMessage 创建JSON错误消息
func NewJSONErrorMessage(code uint32, message, details string) *JSONErrorMessage {
	return &JSONErrorMessage{
		BaseMessage: BaseMessage{msgType: MsgTypeError},
		Code:        code,
		Message:     message,
		Details:     details,
	}
}

// Marshal 序列化消息
func (jem *JSONErrorMessage) Marshal() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"code":    jem.Code,
		"message": jem.Message,
		"details": jem.Details,
	})
}

// Unmarshal 反序列化消息
func (jem *JSONErrorMessage) Unmarshal(data []byte) error {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if code, ok := obj["code"].(float64); ok {
		jem.Code = uint32(code)
	}
	if message, ok := obj["message"].(string); ok {
		jem.Message = message
	}
	if details, ok := obj["details"].(string); ok {
		jem.Details = details
	}

	return nil
}

// String 字符串表示
func (jem *JSONErrorMessage) String() string {
	return fmt.Sprintf("JSONErrorMessage{Code: %d, Message: %s}", jem.Code, jem.Message)
}

// JSONPlayerInfoMessage JSON玩家信息消息
type JSONPlayerInfoMessage struct {
	BaseMessage
	PlayerID string  `json:"player_id"`
	Name     string  `json:"name"`
	Level    uint32  `json:"level"`
	Exp      uint64  `json:"exp"`
	Gold     uint64  `json:"gold"`
	HP       uint32  `json:"hp"`
	MP       uint32  `json:"mp"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        float64 `json:"z"`
}

// NewJSONPlayerInfoMessage 创建JSON玩家信息消息
func NewJSONPlayerInfoMessage(playerID, name string) *JSONPlayerInfoMessage {
	return &JSONPlayerInfoMessage{
		BaseMessage: BaseMessage{msgType: MsgTypePlayerInfo},
		PlayerID:    playerID,
		Name:        name,
	}
}

// Marshal 序列化消息
func (jpim *JSONPlayerInfoMessage) Marshal() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"player_id": jpim.PlayerID,
		"name":      jpim.Name,
		"level":     jpim.Level,
		"exp":       jpim.Exp,
		"gold":      jpim.Gold,
		"hp":        jpim.HP,
		"mp":        jpim.MP,
		"x":         jpim.X,
		"y":         jpim.Y,
		"z":         jpim.Z,
	})
}

// Unmarshal 反序列化消息
func (jpim *JSONPlayerInfoMessage) Unmarshal(data []byte) error {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if playerID, ok := obj["player_id"].(string); ok {
		jpim.PlayerID = playerID
	}
	if name, ok := obj["name"].(string); ok {
		jpim.Name = name
	}
	if level, ok := obj["level"].(float64); ok {
		jpim.Level = uint32(level)
	}
	if exp, ok := obj["exp"].(float64); ok {
		jpim.Exp = uint64(exp)
	}
	if gold, ok := obj["gold"].(float64); ok {
		jpim.Gold = uint64(gold)
	}
	if hp, ok := obj["hp"].(float64); ok {
		jpim.HP = uint32(hp)
	}
	if mp, ok := obj["mp"].(float64); ok {
		jpim.MP = uint32(mp)
	}
	if x, ok := obj["x"].(float64); ok {
		jpim.X = x
	}
	if y, ok := obj["y"].(float64); ok {
		jpim.Y = y
	}
	if z, ok := obj["z"].(float64); ok {
		jpim.Z = z
	}

	return nil
}

// String 字符串表示
func (jpim *JSONPlayerInfoMessage) String() string {
	return fmt.Sprintf("JSONPlayerInfoMessage{PlayerID: %s, Name: %s, Level: %d}", jpim.PlayerID, jpim.Name, jpim.Level)
}

// JSONPlayerMoveMessage JSON玩家移动消息
type JSONPlayerMoveMessage struct {
	BaseMessage
	PlayerID  string  `json:"player_id"`
	FromX     float64 `json:"from_x"`
	FromY     float64 `json:"from_y"`
	FromZ     float64 `json:"from_z"`
	ToX       float64 `json:"to_x"`
	ToY       float64 `json:"to_y"`
	ToZ       float64 `json:"to_z"`
	Speed     float64 `json:"speed"`
	Timestamp int64   `json:"timestamp"`
}

// NewJSONPlayerMoveMessage 创建JSON玩家移动消息
func NewJSONPlayerMoveMessage(playerID string, fromX, fromY, fromZ, toX, toY, toZ, speed float64) *JSONPlayerMoveMessage {
	return &JSONPlayerMoveMessage{
		BaseMessage: BaseMessage{msgType: MsgTypePlayerMove},
		PlayerID:    playerID,
		FromX:       fromX,
		FromY:       fromY,
		FromZ:       fromZ,
		ToX:         toX,
		ToY:         toY,
		ToZ:         toZ,
		Speed:       speed,
		Timestamp:   time.Now().UnixNano(),
	}
}

// Marshal 序列化消息
func (jpmm *JSONPlayerMoveMessage) Marshal() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"player_id": jpmm.PlayerID,
		"from_x":    jpmm.FromX,
		"from_y":    jpmm.FromY,
		"from_z":    jpmm.FromZ,
		"to_x":      jpmm.ToX,
		"to_y":      jpmm.ToY,
		"to_z":      jpmm.ToZ,
		"speed":     jpmm.Speed,
		"timestamp": jpmm.Timestamp,
	})
}

// Unmarshal 反序列化消息
func (jpmm *JSONPlayerMoveMessage) Unmarshal(data []byte) error {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if playerID, ok := obj["player_id"].(string); ok {
		jpmm.PlayerID = playerID
	}
	if fromX, ok := obj["from_x"].(float64); ok {
		jpmm.FromX = fromX
	}
	if fromY, ok := obj["from_y"].(float64); ok {
		jpmm.FromY = fromY
	}
	if fromZ, ok := obj["from_z"].(float64); ok {
		jpmm.FromZ = fromZ
	}
	if toX, ok := obj["to_x"].(float64); ok {
		jpmm.ToX = toX
	}
	if toY, ok := obj["to_y"].(float64); ok {
		jpmm.ToY = toY
	}
	if toZ, ok := obj["to_z"].(float64); ok {
		jpmm.ToZ = toZ
	}
	if speed, ok := obj["speed"].(float64); ok {
		jpmm.Speed = speed
	}
	if timestamp, ok := obj["timestamp"].(float64); ok {
		jpmm.Timestamp = int64(timestamp)
	}

	return nil
}

// String 字符串表示
func (jpmm *JSONPlayerMoveMessage) String() string {
	return fmt.Sprintf("JSONPlayerMoveMessage{PlayerID: %s, From: (%.2f,%.2f,%.2f), To: (%.2f,%.2f,%.2f)}",
		jpmm.PlayerID, jpmm.FromX, jpmm.FromY, jpmm.FromZ, jpmm.ToX, jpmm.ToY, jpmm.ToZ)
}

// JSONPlayerChatMessage JSON玩家聊天消息
type JSONPlayerChatMessage struct {
	BaseMessage
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	Channel    string `json:"channel"`
	Message    string `json:"message"`
	Timestamp  int64  `json:"timestamp"`
}

// NewJSONPlayerChatMessage 创建JSON玩家聊天消息
func NewJSONPlayerChatMessage(playerID, playerName, channel, message string) *JSONPlayerChatMessage {
	return &JSONPlayerChatMessage{
		BaseMessage: BaseMessage{msgType: MsgTypePlayerChat},
		PlayerID:    playerID,
		PlayerName:  playerName,
		Channel:     channel,
		Message:     message,
		Timestamp:   time.Now().UnixNano(),
	}
}

// Marshal 序列化消息
func (jpcm *JSONPlayerChatMessage) Marshal() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"player_id":   jpcm.PlayerID,
		"player_name": jpcm.PlayerName,
		"channel":     jpcm.Channel,
		"message":     jpcm.Message,
		"timestamp":   jpcm.Timestamp,
	})
}

// Unmarshal 反序列化消息
func (jpcm *JSONPlayerChatMessage) Unmarshal(data []byte) error {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if playerID, ok := obj["player_id"].(string); ok {
		jpcm.PlayerID = playerID
	}
	if playerName, ok := obj["player_name"].(string); ok {
		jpcm.PlayerName = playerName
	}
	if channel, ok := obj["channel"].(string); ok {
		jpcm.Channel = channel
	}
	if message, ok := obj["message"].(string); ok {
		jpcm.Message = message
	}
	if timestamp, ok := obj["timestamp"].(float64); ok {
		jpcm.Timestamp = int64(timestamp)
	}

	return nil
}

// Validate 验证聊天消息
func (jpcm *JSONPlayerChatMessage) Validate() error {
	if jpcm.PlayerID == "" {
		return fmt.Errorf("player_id cannot be empty")
	}
	if jpcm.Message == "" {
		return fmt.Errorf("message cannot be empty")
	}
	if len(jpcm.Message) > 500 {
		return fmt.Errorf("message too long: %d > 500", len(jpcm.Message))
	}
	return nil
}

// String 字符串表示
func (jpcm *JSONPlayerChatMessage) String() string {
	return fmt.Sprintf("JSONPlayerChatMessage{PlayerID: %s, Channel: %s, Message: %s}",
		jpcm.PlayerID, jpcm.Channel, jpcm.Message)
}
