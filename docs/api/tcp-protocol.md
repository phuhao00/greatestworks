# TCP 协议文档

## 📖 概述

GreatestWorks 使用自定义的 TCP 协议进行游戏客户端与服务器之间的实时通信。该协议基于二进制格式，提供高效的数据传输和低延迟的游戏体验。

## 🔌 连接流程

### 1. 建立连接

```
客户端 -> 服务器: TCP 连接请求 (端口 8080)
服务器 -> 客户端: 连接确认
```

### 2. 协议握手

```
客户端 -> 服务器: HANDSHAKE_REQUEST
服务器 -> 客户端: HANDSHAKE_RESPONSE
```

### 3. 用户认证

```
客户端 -> 服务器: AUTH_REQUEST
服务器 -> 客户端: AUTH_RESPONSE
```

## 📦 消息格式

### 消息头 (Header)

```
+--------+--------+--------+--------+
| Length |  Type  |   ID   | Flags  |
+--------+--------+--------+--------+
|   4B   |   2B   |   2B   |   1B   |
+--------+--------+--------+--------+
```

- **Length**: 消息总长度 (包含头部)
- **Type**: 消息类型
- **ID**: 消息序列号
- **Flags**: 消息标志位

### 消息体 (Body)

消息体采用 Protocol Buffers 格式，具体结构根据消息类型而定。

## 🎯 消息类型

### 系统消息 (0x0000 - 0x00FF)

| 类型码 | 名称 | 说明 |
|--------|------|------|
| 0x0001 | HANDSHAKE_REQUEST | 握手请求 |
| 0x0002 | HANDSHAKE_RESPONSE | 握手响应 |
| 0x0003 | HEARTBEAT | 心跳包 |
| 0x0004 | DISCONNECT | 断开连接 |

### 认证消息 (0x0100 - 0x01FF)

| 类型码 | 名称 | 说明 |
|--------|------|------|
| 0x0101 | AUTH_REQUEST | 认证请求 |
| 0x0102 | AUTH_RESPONSE | 认证响应 |
| 0x0103 | LOGOUT_REQUEST | 登出请求 |
| 0x0104 | LOGOUT_RESPONSE | 登出响应 |

### 玩家消息 (0x1000 - 0x1FFF)

| 类型码 | 名称 | 说明 |
|--------|------|------|
| 0x1001 | PLAYER_LOGIN | 玩家登录 |
| 0x1002 | PLAYER_LOGOUT | 玩家登出 |
| 0x1003 | PLAYER_MOVE | 玩家移动 |
| 0x1004 | PLAYER_ATTACK | 玩家攻击 |
| 0x1005 | PLAYER_CHAT | 玩家聊天 |
| 0x1006 | PLAYER_UPDATE | 玩家信息更新 |

### 游戏消息 (0x2000 - 0x2FFF)

| 类型码 | 名称 | 说明 |
|--------|------|------|
| 0x2001 | SCENE_ENTER | 进入场景 |
| 0x2002 | SCENE_LEAVE | 离开场景 |
| 0x2003 | ITEM_USE | 使用物品 |
| 0x2004 | SKILL_CAST | 释放技能 |
| 0x2005 | BATTLE_START | 战斗开始 |
| 0x2006 | BATTLE_END | 战斗结束 |

## 🔐 认证流程

### 认证请求

```protobuf
message AuthRequest {
  string username = 1;
  string password = 2;
  string client_version = 3;
  string device_id = 4;
}
```

### 认证响应

```protobuf
message AuthResponse {
  enum Result {
    SUCCESS = 0;
    INVALID_CREDENTIALS = 1;
    ACCOUNT_BANNED = 2;
    SERVER_FULL = 3;
    VERSION_MISMATCH = 4;
  }
  
  Result result = 1;
  string session_id = 2;
  PlayerInfo player_info = 3;
  string message = 4;
}
```

## 🎮 游戏协议示例

### 玩家移动

```protobuf
message PlayerMove {
  string player_id = 1;
  Position from = 2;
  Position to = 3;
  float speed = 4;
  uint64 timestamp = 5;
}

message Position {
  float x = 1;
  float y = 2;
  float z = 3;
}
```

### 玩家攻击

```protobuf
message PlayerAttack {
  string attacker_id = 1;
  string target_id = 2;
  uint32 skill_id = 3;
  Position target_position = 4;
  uint64 timestamp = 5;
}
```

### 聊天消息

```protobuf
message ChatMessage {
  enum Channel {
    WORLD = 0;
    GUILD = 1;
    TEAM = 2;
    PRIVATE = 3;
    SYSTEM = 4;
  }
  
  Channel channel = 1;
  string sender_id = 2;
  string sender_name = 3;
  string content = 4;
  string target_id = 5;  // 私聊目标
  uint64 timestamp = 6;
}
```

## 💓 心跳机制

- **心跳间隔**: 30 秒
- **超时时间**: 90 秒
- **重连机制**: 自动重连，最多尝试 3 次

```protobuf
message Heartbeat {
  uint64 timestamp = 1;
  uint32 sequence = 2;
}
```

## 🔄 消息确认

重要消息需要客户端确认收到：

```protobuf
message MessageAck {
  uint32 message_id = 1;
  uint64 timestamp = 2;
}
```

## 🚫 错误处理

### 错误响应

```protobuf
message ErrorResponse {
  enum ErrorCode {
    UNKNOWN = 0;
    INVALID_MESSAGE = 1;
    PERMISSION_DENIED = 2;
    RATE_LIMITED = 3;
    SERVER_ERROR = 4;
  }
  
  ErrorCode code = 1;
  string message = 2;
  uint32 original_message_id = 3;
}
```

## 📊 性能指标

- **消息处理延迟**: < 10ms
- **并发连接数**: 10,000+
- **消息吞吐量**: 100,000 msg/s
- **网络带宽**: 优化后平均 < 1KB/s per player

## 🛠️ 开发工具

### 协议调试

```bash
# 使用 telnet 连接服务器
telnet localhost 8080

# 使用 netcat 发送二进制数据
echo -ne '\x00\x00\x00\x09\x00\x03\x00\x01\x00' | nc localhost 8080
```

### 消息解析工具

```go
// Go 示例代码
func parseMessage(data []byte) (*Message, error) {
    if len(data) < 9 {
        return nil, errors.New("message too short")
    }
    
    length := binary.BigEndian.Uint32(data[0:4])
    msgType := binary.BigEndian.Uint16(data[4:6])
    msgID := binary.BigEndian.Uint16(data[6:8])
    flags := data[8]
    
    return &Message{
        Length: length,
        Type:   msgType,
        ID:     msgID,
        Flags:  flags,
        Body:   data[9:],
    }, nil
}
```

## 🔒 安全考虑

- **消息加密**: 敏感数据使用 AES 加密
- **防重放攻击**: 消息包含时间戳和序列号
- **频率限制**: 限制客户端消息发送频率
- **输入验证**: 服务器端严格验证所有输入

---

*协议版本: v1.0.0 | 最后更