# REST API 接口文档

## 📖 概述

GreatestWorks REST API 提供了游戏管理、玩家数据查询、系统配置等功能的 HTTP 接口。

## 🔐 认证

所有 API 请求都需要在请求头中包含有效的 JWT Token：

```http
Authorization: Bearer <jwt_token>
```

## 🎯 API 端点

### 认证相关

#### 用户登录

```http
POST /api/auth/login
```

**请求体:**
```json
{
  "username": "string",
  "password": "string"
}
```

**响应:**
```json
{
  "success": true,
  "data": {
    "token": "jwt_token_string",
    "expires_in": 3600,
    "user_id": "uuid"
  }
}
```

#### 刷新令牌

```http
POST /api/auth/refresh
```

### 玩家管理

#### 获取玩家信息

```http
GET /api/player/{player_id}
```

**响应:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "username": "string",
    "level": 50,
    "experience": 125000,
    "gold": 10000,
    "created_at": "2024-01-01T00:00:00Z",
    "last_login": "2024-01-15T12:30:00Z"
  }
}
```

#### 更新玩家信息

```http
PUT /api/player/{player_id}
```

**请求体:**
```json
{
  "nickname": "string",
  "avatar": "string"
}
```

#### 获取玩家列表

```http
GET /api/players?page=1&limit=20&sort=level&order=desc
```

**查询参数:**
- `page`: 页码 (默认: 1)
- `limit`: 每页数量 (默认: 20, 最大: 100)
- `sort`: 排序字段 (level, created_at, last_login)
- `order`: 排序方向 (asc, desc)

### 游戏数据

#### 获取排行榜

```http
GET /api/ranking/{type}?limit=100
```

**路径参数:**
- `type`: 排行榜类型 (level, gold, pvp_score)

**响应:**
```json
{
  "success": true,
  "data": {
    "type": "level",
    "updated_at": "2024-01-15T12:00:00Z",
    "rankings": [
      {
        "rank": 1,
        "player_id": "uuid",
        "username": "string",
        "value": 100,
        "change": "+2"
      }
    ]
  }
}
```

#### 获取服务器状态

```http
GET /api/server/status
```

**响应:**
```json
{
  "success": true,
  "data": {
    "server_id": "server-001",
    "status": "online",
    "online_players": 1250,
    "max_players": 2000,
    "uptime": 86400,
    "version": "1.0.0",
    "last_restart": "2024-01-14T00:00:00Z"
  }
}
```

### 管理接口

#### 系统配置

```http
GET /api/admin/config
PUT /api/admin/config
```

#### 服务器管理

```http
POST /api/admin/server/restart
POST /api/admin/server/maintenance
```

#### 玩家管理

```http
POST /api/admin/player/{player_id}/ban
POST /api/admin/player/{player_id}/unban
POST /api/admin/player/{player_id}/kick
```

## 📊 响应格式

### 成功响应

```json
{
  "success": true,
  "data": {},
  "message": "操作成功",
  "timestamp": "2024-01-15T12:30:00Z"
}
```

### 错误响应

```json
{
  "success": false,
  "error": {
    "code": "INVALID_PARAMETER",
    "message": "参数验证失败",
    "details": {
      "field": "username",
      "reason": "用户名不能为空"
    }
  },
  "timestamp": "2024-01-15T12:30:00Z"
}
```

## 🔄 分页

支持分页的接口使用统一的分页格式：

```json
{
  "success": true,
  "data": {
    "items": [],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 100,
      "pages": 5,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

## 🚦 限流

- **普通用户**: 100 请求/分钟
- **管理员**: 1000 请求/分钟
- **系统接口**: 10 请求/分钟

## 📝 示例代码

### JavaScript (Fetch)

```javascript
const response = await fetch('/api/player/123', {
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});

const data = await response.json();
```

### Go

```go
req, _ := http.NewRequest("GET", "/api/player/123", nil)
req.Header.Set("Authorization", "Bearer "+token)
req.Header.Set("Content-Type", "application/json")

resp, err := client.Do(req)
```

---

*API 版本: v1.0.0 | 最后更新: 2024年*