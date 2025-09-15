# Greatest Works - MMO Game Server

基于Go语言开发的大型多人在线游戏服务器架构，采用微服务设计，支持高并发和分布式部署。

## 项目特性

- 🚀 **高性能网络架构**: 基于netcore-go的TCP网络框架
- 🔧 **微服务架构**: 网关、场景、战斗、活动服务器分离
- 💾 **多数据库支持**: MongoDB + Redis 混合存储
- 🔐 **JWT认证**: 安全的用户认证系统
- 🎮 **完整游戏功能**: 玩家系统、场景管理、战斗系统、活动系统
- 📊 **实时同步**: 高频率的游戏状态同步
- 🛡️ **容错设计**: 完善的错误处理和恢复机制

## 项目结构

```
greatestworks/
├── cmd/
│   └── server/          # 服务器启动入口
│       └── main.go
├── config/              # 配置管理
│   └── config.go
├── internal/            # 内部模块
│   ├── auth/           # 认证模块
│   │   └── jwt.go
│   ├── database/       # 数据库模块
│   │   ├── mongodb.go
│   │   └── redis.go
│   └── game/           # 游戏逻辑模块
│       └── player.go
├── server/             # 服务器实现
│   ├── gateway/        # 网关服务器
│   │   └── server.go
│   ├── scene/          # 场景服务器
│   │   └── server.go
│   ├── battle/         # 战斗服务器
│   │   └── server.go
│   └── activity/       # 活动服务器
│       └── server.go
├── protocol/           # 网络协议定义
│   └── protocol.go
├── go.mod              # Go模块定义
├── go.work             # Go工作空间配置
└── README.md           # 项目说明
```

## 技术栈

- **语言**: Go 1.21+
- **网络框架**: netcore-go (TCP)
- **数据库**: MongoDB + Redis
- **认证**: JWT
- **协议**: 自定义二进制协议 + JSON
- **架构**: 微服务 + 分布式

## 快速开始

### 环境要求

- Go 1.21 或更高版本
- MongoDB 4.4+
- Redis 6.0+

### 安装依赖

```bash
go mod tidy
```

### 配置文件

创建 `config/config.yaml` 配置文件：

```yaml
server:
  gateway:
    port: 8080
    host: "0.0.0.0"
    max_connections: 10000
    read_timeout: 30
    write_timeout: 30
    heartbeat_time: 60
  scene:
    port: 8081
    host: "0.0.0.0"
    max_players: 1000
    tick_rate: 20
    sync_interval: 100
  battle:
    port: 8082
    host: "0.0.0.0"
    max_battles: 100
    tick_rate: 30
    battle_time: 300
  activity:
    port: 8083
    host: "0.0.0.0"
    max_activities: 50
    update_interval: 1000

mongodb:
  uri: "mongodb://localhost:27017"
  database: "mmo_game"
  max_pool_size: 100
  min_pool_size: 10
  max_idle_time: 300
  connect_timeout: 10
  socket_timeout: 30

redis:
  addr: "localhost:6379"
  password: ""
  db: 0
  pool_size: 100
  min_idle_conns: 10
  max_idle_conns: 50
  conn_max_age: 3600
  dial_timeout: 5
  read_timeout: 3
  write_timeout: 3

jwt:
  secret_key: "your-super-secret-key-change-this-in-production"
  token_duration: 24
  refresh_time: 1

log:
  level: "info"
  format: "json"
  output: "stdout"
  max_size: 100
  max_backups: 3
  max_age: 28
  compress: true

game:
  max_level: 100
  exp_multiplier: 1.0
  gold_multiplier: 1.0
  drop_rate: 0.1
  pk_enabled: true
  guild_enabled: true
  trade_enabled: true

network:
  protocol: "tcp"
  buffer_size: 4096
  max_packet_size: 65536
  compression_type: "none"
  encryption_type: "none"
```

### 启动服务器

#### 启动网关服务器
```bash
go run cmd/server/main.go -type=gateway -port=8080
```

#### 启动场景服务器
```bash
go run cmd/server/main.go -type=scene -port=8081
```

#### 启动战斗服务器
```bash
go run cmd/server/main.go -type=battle -port=8082
```

#### 启动活动服务器
```bash
go run cmd/server/main.go -type=activity -port=8083
```

## 服务器架构

### 网关服务器 (Gateway)
- 处理客户端连接和认证
- 消息路由和转发
- 负载均衡
- 心跳检测

### 场景服务器 (Scene)
- 管理游戏场景和地图
- 处理玩家移动和交互
- NPC和怪物AI
- 场景同步

### 战斗服务器 (Battle)
- 处理PvP和PvE战斗
- 技能系统
- 伤害计算
- 战斗奖励

### 活动服务器 (Activity)
- 管理游戏活动
- 任务系统
- 排行榜
- 奖励发放

## 网络协议

### 数据包格式
```
+--------+--------+----------+
| Length | Type   | Data     |
| 4bytes | 2bytes | Variable |
+--------+--------+----------+
```

### 消息类型
- **1xxx**: 基础消息 (登录、心跳等)
- **2xxx**: 玩家相关消息
- **3xxx**: 聊天相关消息
- **4xxx**: 场景相关消息
- **5xxx**: 战斗相关消息
- **6xxx**: 活动相关消息
- **7xxx**: 物品相关消息
- **8xxx**: 交易相关消息
- **9xxx**: 公会相关消息

## 数据库设计

### MongoDB 集合
- `players`: 玩家基础数据
- `guilds`: 公会信息
- `activities`: 活动数据
- `battles`: 战斗记录
- `items`: 物品模板
- `quests`: 任务数据

### Redis 缓存
- 在线玩家列表
- 会话信息
- 排行榜数据
- 临时战斗数据

## 开发指南

### 添加新的消息类型

1. 在 `protocol/protocol.go` 中定义消息常量和结构体
2. 在相应的服务器中注册消息处理器
3. 实现消息处理逻辑

### 添加新的服务器类型

1. 在 `server/` 目录下创建新的服务器包
2. 实现 `Server` 接口
3. 在 `cmd/server/main.go` 中添加启动逻辑

### 数据库操作

使用 `internal/database/` 包中的封装方法进行数据库操作，避免直接使用原生驱动。

## 性能优化

- 使用连接池管理数据库连接
- Redis缓存热点数据
- 消息批量处理
- 异步日志记录
- 内存池复用对象

## 监控和日志

- 结构化日志输出
- 性能指标收集
- 错误追踪
- 健康检查接口

## 部署

### Docker 部署
```bash
# 构建镜像
docker build -t mmo-server .

# 运行容器
docker run -d -p 8080:8080 mmo-server
```

### Kubernetes 部署
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mmo-gateway
spec:
  replicas: 3
  selector:
    matchLabels:
      app: mmo-gateway
  template:
    metadata:
      labels:
        app: mmo-gateway
    spec:
      containers:
      - name: gateway
        image: mmo-server:latest
        args: ["-type=gateway"]
        ports:
        - containerPort: 8080
```

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License

## 联系方式

- 项目地址: https://github.com/your-org/greatestworks
- 问题反馈: https://github.com/your-org/greatestworks/issues
- 邮箱: dev@greatestworks.com

## 更新日志

### v1.0.0 (2024-01-15)
- 初始版本发布
- 基础服务器架构
- 玩家系统
- 场景管理
- 战斗系统
- 活动系统

---

**注意**: 这是一个开发中的项目，部分功能可能还不完善。欢迎贡献代码和提出建议！