# WebSocket API 文档

## 📖 概述

GreatestWorks WebSocket API 提供实时双向通信功能，主要用于游戏内的实时消息推送、聊天系统、实时状态更新等场景。

## 🔌 连接建立

### 连接 URL

```
ws://localhost:8081/ws
wss://game.example.com/ws  # 生产环境
```

### 连接参数

```javascript
const ws = new WebSocket('ws://localhost:8081/ws', {
  headers: {
    'Authorization': 'Bearer ' + jwt_token,
    'X-Client-Version': '1.0.0'
  }
});
```

### 连接状态

| 状态 | 说明 |
|------|------|
| CONNECTING | 正在连接 |
| OPEN | 连接已建立 |
| CLOSING | 正在关闭 |
| CLOSED | 连接已关闭 |

## 📦 消息格式

### 基础消息结构

```json
{
  "type": "message_type",
  "id": "unique_message_id",
  "timestamp": 1642234567890,
  "data": {}
}
```

### 消息类型

- **type**: 消息类型标识符
- **id**: 唯一消息 ID (用于消息确认)
- **timestamp**: 消息时间戳 (毫秒)
- **data**: 消息数据载荷

## 🎯 消息类型

### 系统消息

#### 连接确认

**服务器 -> 客户端**

```json
{
  "type": "connection_ack",
  "id": "msg_001",
  "timestamp": 1642234567890,
  "data": {
    "session_id": "sess_123456",
    "player_id": "player_789",
    "server_time": 1642234567890
  }
}
```

#### 心跳

**客户端 <-> 服务器**

```json
{
  "type": "ping",
  "id": "ping_001",
  "timestamp": 1642234567890,
  "data": {}
}
```

```json
{
  "type": "pong",
  "id": "pong_001",
  "timestamp": 1642234567890,
  "data": {
    "ping_id": "ping_001"
  }
}
```

### 聊天消息

#### 发送聊天消息

**客户端 -> 服务器**

```json
{
  "type": "chat_send",
  "id": "chat_001",
  "timestamp": 1642234567890,
  "data": {
    "channel": "world",
    "content": "Hello, world!",
    "target_id": null
  }
}
```

#### 接收聊天消息

**服务器 -> 客户端**

```json
{
  "type": "chat_message",
  "id": "chat_002",
  "timestamp": 1642234567890,
  "data": {
    "channel": "world",
    "sender_id": "player_123",
    "sender_name": "PlayerName",
    "content": "Hello, world!",
    "message_id": "msg_456"
  }
}
```

### 实时状态更新

#### 玩家状态变化

**服务器 -> 客户端**

```json
{
  "type": "player_status_update",
  "id": "status_001",
  "timestamp": 1642234567890,
  "data": {
    "player_id": "player_123",
    "status": "online",
    "level": 50,
    "location": {
      "scene_id": "scene_001",
      "x": 100.5,
      "y": 200.3
    }
  }
}
```

#### 好友上线通知

**服务器 -> 客户端**

```json
{
  "type": "friend_online",
  "id": "friend_001",
  "timestamp": 1642234567890,
  "data": {
    "friend_id": "player_456",
    "friend_name": "FriendName",
    "login_time": 1642234567890
  }
}
```

### 游戏事件

#### 战斗结果通知

**服务器 -> 客户端**

```json
{
  "type": "battle_result",
  "id": "battle_001",
  "timestamp": 1642234567890,
  "data": {
    "battle_id": "battle_123",
    "result": "victory",
    "rewards": {
      "experience": 1000,
      "gold": 500,
      "items": [
        {"id": "item_001", "quantity": 1}
      ]
    }
  }
}
```

#### 系统公告

**服务器 -> 客户端**

```json
{
  "type": "system_announcement",
  "id": "announce_001",
  "timestamp": 1642234567890,
  "data": {
    "title": "系统维护通知",
    "content": "服务器将于今晚 22:00 进行维护",
    "priority": "high",
    "duration": 300000
  }
}
```

## 🔄 消息确认机制

### 消息确认

**客户端 -> 服务器**

```json
{
  "type": "message_ack",
  "id": "ack_001",
  "timestamp": 1642234567890,
  "data": {
    "message_id": "chat_002"
  }
}
```

### 重要消息重发

服务器会对重要消息进行重发，直到收到客户端确认或达到最大重试次数。

## 🚫 错误处理

### 错误消息格式

```json
{
  "type": "error",
  "id": "error_001",
  "timestamp": 1642234567890,
  "data": {
    "code": "INVALID_MESSAGE",
    "message": "消息格式无效",
    "original_message_id": "chat_001"
  }
}
```

### 错误代码

| 错误代码 | 说明 |
|----------|------|
| INVALID_MESSAGE | 消息格式无效 |
| UNAUTHORIZED | 未授权访问 |
| RATE_LIMITED | 频率限制 |
| SERVER_ERROR | 服务器内部错误 |
| CONNECTION_LOST | 连接丢失 |

## 💓 心跳机制

- **心跳间隔**: 30 秒
- **超时检测**: 90 秒无响应则断开连接
- **自动重连**: 客户端应实现自动重连机制

```javascript
// 心跳实现示例
setInterval(() => {
  if (ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({
      type: 'ping',
      id: generateId(),
      timestamp: Date.now(),
      data: {}
    }));
  }
}, 30000);
```

## 🔒 安全机制

### 认证

- 连接时需要提供有效的 JWT Token
- Token 过期后需要重新认证

### 频率限制

- **聊天消息**: 10 条/分钟
- **一般消息**: 100 条/分钟
- **心跳消息**: 不限制

### 消息验证

- 服务器验证所有接收到的消息格式
- 过滤恶意内容和非法字符

## 📱 客户端实现示例

### JavaScript

```javascript
class GameWebSocket {
  constructor(url, token) {
    this.url = url;
    this.token = token;
    this.ws = null;
    this.messageHandlers = new Map();
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
  }

  connect() {
    this.ws = new WebSocket(this.url);
    
    this.ws.onopen = () => {
      console.log('WebSocket connected');
      this.reconnectAttempts = 0;
      this.startHeartbeat();
    };

    this.ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      this.handleMessage(message);
    };

    this.ws.onclose = () => {
      console.log('WebSocket disconnected');
      this.stopHeartbeat();
      this.attemptReconnect();
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
  }

  send(type, data) {
    if (this.ws.readyState === WebSocket.OPEN) {
      const message = {
        type,
        id: this.generateId(),
        timestamp: Date.now(),
        data
      };
      this.ws.send(JSON.stringify(message));
    }
  }

  onMessage(type, handler) {
    this.messageHandlers.set(type, handler);
  }

  handleMessage(message) {
    const handler = this.messageHandlers.get(message.type);
    if (handler) {
      handler(message.data);
    }
  }

  generateId() {
    return 'msg_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
  }
}

// 使用示例
const gameWS = new GameWebSocket('ws://localhost:8081/ws', token);
gameWS.connect();

gameWS.onMessage('chat_message', (data) => {
  console.log('收到聊天消息:', data);
});

gameWS.send('chat_send', {
  channel: 'world',
  content: 'Hello, WebSocket!'
});
```

### Go 客户端

```go
package main

import (
    "encoding/json"
    "log"
    "github.com/gorilla/websocket"
)

type GameClient struct {
    conn *websocket.Conn
    token string
}

type Message struct {
    Type      string      `json:"type"`
    ID        string      `json:"id"`
    Timestamp int64       `json:"timestamp"`
    Data      interface{} `json:"data"`
}

func (c *GameClient) Connect(url string) error {
    conn, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        return err
    }
    c.conn = conn
    
    go c.readMessages()
    return nil
}

func (c *GameClient) Send(msgType string, data interface{}) error {
    msg := Message{
        Type:      msgType,
        ID:        generateID(),
        Timestamp: time.Now().UnixMilli(),
        Data:      data,
    }
    
    return c.conn.WriteJSON(msg)
}

func (c *GameClient) readMessages() {
    for {
        var msg Message
        err := c.conn.ReadJSON(&msg)
        if err != nil {
            log.Printf("读取消息错误: %v", err)
            break
        }
        
        c.handleMessage(msg)
    }
}
```

## 📊 性能指标

- **连接延迟**: < 100ms
- **消息延迟**: < 50ms
- **并发连接**: 50,000+
- **消息吞吐**: 500,000 msg/s
- **内存使用**: < 1MB per 1000 connections

---

*API 版本: v1.0.0 | 最后更新: 2024年*