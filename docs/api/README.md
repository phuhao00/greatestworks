# API 文档

GreatestWorks MMO 游戏服务器提供多种 API 接口，支持不同的通信协议和使用场景。

## 📋 API 概览

### 通信协议

| 协议类型 | 用途 | 文档链接 |
|---------|------|----------|
| REST API | 管理接口、配置查询 | [REST API](./rest-api.md) |
| TCP 协议 | 游戏核心逻辑通信 | [TCP 协议](./tcp-protocol.md) |
| WebSocket | 实时消息推送 | [WebSocket API](./websocket-api.md) |

### 认证方式

- **JWT Token**: 用于 REST API 认证
- **Session ID**: 用于 TCP 连接认证
- **WebSocket Token**: 用于 WebSocket 连接认证

## 🔧 API 使用指南

### 1. 获取访问令牌

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "player123",
  "password": "password"
}
```

### 2. 使用令牌访问 API

```http
GET /api/player/profile
Authorization: Bearer <jwt_token>
```

### 3. 建立 TCP 连接

```go
// 连接游戏服务器
conn, err := net.Dial("tcp", "localhost:8080")
// 发送认证消息
// ...
```

## 📊 API 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未授权访问 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 🔍 调试工具

- **Postman**: REST API 测试
- **WebSocket King**: WebSocket 连接测试
- **Telnet**: TCP 协议调试

## 📝 更新日志

- **v1.0.0**: 初始版本发布
- 支持基础的玩家认证和游戏操作
- 提供完整的 REST API 接口

---

*更多详细信息请查看具体的协议文档*