# Greatest Works 网络架构重构文档

## 概述

本文档描述了 Greatest Works 游戏服务器的网络架构重构，实现了HTTP、TCP和gRPC多协议支持，明确了各协议的职责分工。

## 架构设计原则

### 协议职责分工

- **HTTP服务器** (端口: 8080)
  - 用户注册和登录
  - JWT Token管理
  - GM后台管理接口
  - RESTful API
  - Web管理界面

- **TCP服务器** (端口: 9090)
  - 游戏核心逻辑处理
  - 实时游戏数据同步
  - 玩家移动和战斗
  - 长连接维护
  - 心跳检测

- **gRPC服务器** (端口: 9091)
  - 服务间通信
  - 微服务架构支持
  - 高性能RPC调用
  - 类型安全的接口

## 目录结构

```
internal/interfaces/
├── http/                    # HTTP接口层
│   ├── auth/               # 认证相关
│   │   ├── login_handler.go
│   │   ├── register_handler.go
│   │   ├── token_handler.go
│   │   └── middleware.go
│   ├── gm/                 # GM管理接口
│   │   ├── player_management.go
│   │   └── server_monitor.go
│   └── server.go           # HTTP服务器
├── tcp/                    # TCP接口层
│   ├── handlers/           # 消息处理器
│   │   └── game_handler.go
│   ├── protocol/           # 协议定义
│   │   └── message_types.go
│   ├── connection/         # 连接管理
│   │   ├── manager.go
│   │   ├── session.go
│   │   └── heartbeat.go
│   ├── router.go           # 消息路由
│   └── server.go           # TCP服务器
└── grpc/                   # gRPC接口层
    ├── proto/              # Protocol Buffers定义
    │   ├── common.proto
    │   ├── player.proto
    │   ├── battle.proto
    │   └── notification.proto
    ├── services/           # gRPC服务实现
    │   └── player_service.go
    ├── interceptors/       # 拦截器
    │   ├── auth.go
    │   ├── logging.go
    │   └── metrics.go
    └── server.go           # gRPC服务器
```

## 核心功能

### 1. HTTP接口层

#### 认证系统
- JWT Token生成和验证
- 用户注册和登录
- Token刷新机制
- 角色权限管理

#### GM管理系统
- 玩家数据管理
- 服务器状态监控
- 系统配置管理
- 日志查询接口

### 2. TCP接口层

#### 消息协议
```go
type MessageHeader struct {
    Magic       uint32 // 魔数标识
    MessageID   uint32 // 消息ID
    MessageType uint16 // 消息类型
    PlayerID    uint64 // 玩家ID
    Timestamp   int64  // 时间戳
    Sequence    uint32 // 序列号
    Length      uint32 // 消息长度
}
```

#### 连接管理
- 连接池管理
- 会话状态维护
- 心跳检测机制
- 断线重连处理

### 3. gRPC接口层

#### 服务定义
- PlayerService: 玩家数据服务
- BattleService: 战斗逻辑服务
- NotificationService: 通知推送服务

#### 拦截器
- 认证拦截器: JWT验证
- 日志拦截器: 请求日志记录
- 监控拦截器: 性能指标收集

## 认证机制

### JWT认证流程

1. **用户登录** (HTTP)
   ```
   POST /api/auth/login
   {
     "username": "user",
     "password": "password"
   }
   ```

2. **获取Token**
   ```json
   {
     "access_token": "eyJ...",
     "refresh_token": "eyJ...",
     "token_type": "Bearer",
     "expires_in": 3600
   }
   ```

3. **TCP认证**
   ```go
   type AuthRequest struct {
     Token    string `json:"token"`
     PlayerID string `json:"player_id"`
   }
   ```

4. **gRPC认证**
   ```
   Authorization: Bearer eyJ...
   ```

## 监控和日志

### 指标收集
- HTTP请求计数和延迟
- TCP连接数和消息统计
- gRPC调用性能指标
- 系统资源使用情况
- 错误率统计

### 日志记录
- 结构化日志格式
- 请求追踪ID
- 错误堆栈信息
- 性能监控数据

### 错误处理
- 统一错误码定义
- 分层错误处理
- 错误上报机制
- 优雅降级策略

## 配置管理

### 服务器配置 (config/server.yaml)
```yaml
http:
  enabled: true
  addr: ":8080"
  read_timeout: 30s
  write_timeout: 30s

tcp:
  enabled: true
  addr: ":9090"
  max_connections: 10000
  heartbeat:
    enabled: true
    interval: 30s

grpc:
  enabled: true
  addr: ":9091"
  max_connections: 1000
  enable_reflection: true
```

## 部署和运维

### 启动服务器
```bash
go run cmd/server/main.go
```

### 健康检查
- HTTP: `GET /health`
- gRPC: `grpc.health.v1.Health/Check`
- TCP: 心跳检测

### 监控端点
- 指标: `GET /metrics`
- 状态: `GET /status`
- 调试: `GET /debug/pprof`

## 性能优化

### 连接管理
- 连接池复用
- 长连接维护
- 优雅关闭机制

### 消息处理
- 异步消息处理
- 消息队列缓冲
- 批量操作优化

### 缓存策略
- 内存缓存
- Redis分布式缓存
- 数据预加载

## 安全措施

### 认证安全
- JWT签名验证
- Token过期检查
- 刷新Token机制
- 会话管理

### 网络安全
- TLS加密支持
- 请求限流
- DDoS防护
- 输入验证

### 数据安全
- 敏感数据加密
- SQL注入防护
- XSS攻击防护
- CSRF保护

## 扩展性设计

### 水平扩展
- 无状态服务设计
- 负载均衡支持
- 服务发现机制
- 分布式会话

### 微服务架构
- gRPC服务间通信
- 服务注册发现
- 配置中心集成
- 分布式追踪

## 测试策略

### 单元测试
- 业务逻辑测试
- 接口功能测试
- 错误处理测试

### 集成测试
- 多协议交互测试
- 端到端流程测试
- 性能压力测试

### 监控测试
- 健康检查测试
- 指标收集测试
- 告警机制测试

## 故障排查

### 常见问题
1. **连接超时**: 检查网络配置和防火墙
2. **认证失败**: 验证JWT配置和密钥
3. **消息丢失**: 检查TCP缓冲区和心跳
4. **性能问题**: 查看监控指标和日志

### 调试工具
- 日志分析
- 性能分析
- 网络抓包
- 压力测试

## 未来规划

### 功能扩展
- WebSocket支持
- 消息队列集成
- 分布式缓存
- 服务网格

### 性能优化
- 协议优化
- 内存管理
- 并发优化
- 数据库优化

### 运维改进
- 自动化部署
- 容器化支持
- 监控告警
- 日志聚合

---

## 总结

本次网络架构重构成功实现了：

✅ **协议分工明确**: HTTP负责认证管理，TCP处理游戏逻辑，gRPC支持服务通信

✅ **统一认证机制**: JWT Token在所有协议中统一使用

✅ **完善监控体系**: 指标收集、日志记录、错误处理一体化

✅ **高可用设计**: 心跳检测、优雅关闭、错误恢复机制

✅ **扩展性良好**: 支持水平扩展和微服务架构

✅ **安全可靠**: 多层安全防护和数据保护

该架构为 Greatest Works 游戏提供了稳定、高性能、可扩展的网络基础设施。