# Greatest Works - MMO Game Server

基于Go语言和领域驱动设计(DDD)架构开发的大型多人在线游戏服务器，采用现代化微服务设计，支持高并发和分布式部署。

## 🎯 项目概述

这是一个企业级的MMO游戏服务器项目，采用领域驱动设计(Domain-Driven Design)架构模式，提供高性能、可扩展、易维护的游戏服务器解决方案。项目包含完整的游戏系统，如玩家管理、社交系统、战斗系统、建筑系统、宠物系统等。

## ✨ 核心特性

- 🏗️ **DDD架构**: 采用领域驱动设计，清晰的架构分层和职责分离
- 🚀 **高性能网络**: 基于netcore-go的TCP网络框架，支持高并发连接
- 🔧 **微服务设计**: 模块化设计，支持独立部署和扩展
- 💾 **多数据库支持**: MongoDB + Redis 混合存储策略
- 🔐 **安全认证**: JWT认证系统，保障用户数据安全
- 🎮 **完整游戏功能**: 涵盖现代MMO游戏的核心系统
- 📊 **实时同步**: 高频率的游戏状态同步和事件处理
- 🛡️ **容错设计**: 完善的错误处理、监控和恢复机制
- 🐳 **容器化部署**: Docker和Kubernetes支持
- 📚 **完整文档**: 详细的API文档和架构说明

## 🏗️ DDD架构设计

本项目采用领域驱动设计(Domain-Driven Design)架构，将复杂的游戏业务逻辑按照领域进行划分，实现高内聚、低耦合的系统设计。

### 架构分层

- **接口层 (Interfaces)**: 处理外部请求，包括TCP、HTTP、gRPC接口
- **应用层 (Application)**: 协调领域对象，处理业务用例
- **领域层 (Domain)**: 核心业务逻辑和领域模型
- **基础设施层 (Infrastructure)**: 技术实现，如数据库、缓存、消息队列

## 📁 项目结构

```
greatestworks/
├── cmd/                        # 应用程序入口
│   └── server/
│       ├── bootstrap.go        # 启动引导
│       └── main.go            # 主程序
├── configs/                    # 配置模板
│   ├── config.example.yaml    # 基础配置模板
│   ├── config.dev.yaml.example # 开发环境配置
│   ├── config.prod.yaml.example # 生产环境配置
│   └── docker.yaml            # Docker环境配置
├── docs/                       # 项目文档
│   ├── api/                   # API文档
│   ├── architecture/          # 架构文档
│   ├── deployment/            # 部署文档
│   └── diagrams/              # 架构图表
├── application/                # 应用层
│   ├── commands/              # 命令处理器
│   ├── handlers/              # 事件处理器
│   ├── queries/               # 查询处理器
│   └── services/              # 应用服务
├── internal/                   # 内部模块
│   ├── domain/                # 领域层
│   │   ├── player/           # 玩家领域
│   │   ├── battle/           # 战斗领域
│   │   ├── social/           # 社交领域
│   │   ├── building/         # 建筑领域
│   │   ├── pet/              # 宠物领域
│   │   ├── ranking/          # 排行榜领域
│   │   └── minigame/         # 小游戏领域
│   ├── infrastructure/        # 基础设施层
│   │   ├── persistence/      # 数据持久化
│   │   ├── cache/            # 缓存服务
│   │   ├── messaging/        # 消息服务
│   │   ├── network/          # 网络服务
│   │   ├── config/           # 配置管理
│   │   └── logging/          # 日志服务
│   └── interfaces/            # 接口层
│       ├── tcp/              # TCP接口
│       ├── http/             # HTTP接口
│       └── grpc/             # gRPC接口
├── scripts/                    # 开发脚本
│   ├── build.sh              # 构建脚本
│   ├── deploy.sh             # 部署脚本
│   └── test.sh               # 测试脚本
├── docker-compose.yml          # Docker编排
├── Dockerfile                  # Docker镜像
├── Makefile                   # 构建工具
├── go.mod                     # Go模块定义
└── README.md                  # 项目说明
```

## 🛠️ 技术栈

### 核心技术
- **语言**: Go 1.21+
- **架构模式**: 领域驱动设计 (DDD)
- **网络框架**: netcore-go (TCP) + HTTP + gRPC
- **数据库**: MongoDB (主数据库) + Redis (缓存)
- **消息队列**: NATS
- **认证**: JWT + 自定义认证
- **协议**: 自定义二进制协议 + JSON + Protobuf

### 开发工具
- **构建工具**: Make + Go Modules
- **容器化**: Docker + Docker Compose
- **编排**: Kubernetes
- **代码质量**: golangci-lint + 自定义规范
- **文档**: Markdown + 架构图

### 监控与运维
- **日志**: 结构化日志 + 分级输出
- **监控**: 自定义指标收集
- **健康检查**: HTTP健康检查接口
- **配置管理**: YAML配置 + 环境变量

## 🚀 快速开始

### 📋 环境要求

- **Go**: 1.21 或更高版本
- **MongoDB**: 4.4+ (推荐 5.0+)
- **Redis**: 6.0+ (推荐 7.0+)
- **NATS**: 2.9+ (可选，用于消息队列)
- **Docker**: 20.10+ (可选，用于容器化部署)

### 📦 安装依赖

```bash
# 克隆项目
git clone https://github.com/your-org/greatestworks.git
cd greatestworks

# 安装Go依赖
go mod tidy

# 使用Make命令安装开发工具
make setup
```

### ⚙️ 配置文件

复制配置模板并根据环境进行配置：

```bash
# 开发环境
cp configs/config.dev.yaml.example config.yaml

# 生产环境
cp configs/config.prod.yaml.example config.yaml
```

基础配置示例：

```yaml
# 服务器配置
server:
  port: 8080
  host: "0.0.0.0"
  max_connections: 10000
  read_timeout: 30s
  write_timeout: 30s
  shutdown_timeout: 10s

# 数据库配置
database:
  mongodb:
    uri: "mongodb://localhost:27017"
    database: "mmo_game"
    max_pool_size: 100
    connect_timeout: 10s
  redis:
    addr: "localhost:6379"
    password: ""
    db: 0
    pool_size: 100
    dial_timeout: 5s

# 消息队列配置
messaging:
  nats:
    url: "nats://localhost:4222"
    max_reconnects: 10
    reconnect_wait: 2s

# 认证配置
auth:
  jwt:
    secret: "your-super-secret-key-change-this-in-production"
    expire: 24h
    refresh_expire: 168h

# 日志配置
logging:
  level: "info"
  format: "json"
  output: "stdout"

# 游戏配置
game:
  max_level: 100
  max_players: 1000
  tick_rate: 20
  save_interval: 300s
```

### 🎮 启动服务器

#### 开发环境启动
```bash
# 使用Make命令启动开发服务器
make dev

# 或者直接运行
go run cmd/server/main.go
```

#### 生产环境启动
```bash
# 构建二进制文件
make build

# 启动服务器
./bin/server -config=config.yaml
```

#### Docker启动
```bash
# 使用Docker Compose启动完整环境
docker-compose up -d

# 仅启动游戏服务器
docker run -d -p 8080:8080 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  greatestworks:latest
```

## 🏛️ DDD领域架构

### 核心领域 (Core Domains)

#### 🎮 玩家领域 (Player Domain)
- **职责**: 玩家基础信息、等级经验、属性管理
- **核心实体**: Player, PlayerStats, PlayerProfile
- **主要功能**: 玩家创建、升级、属性计算、状态管理

#### ⚔️ 战斗领域 (Battle Domain)
- **职责**: 战斗逻辑、技能系统、伤害计算
- **核心实体**: Battle, Skill, Damage, BattleResult
- **主要功能**: PvP/PvE战斗、技能释放、战斗结算

#### 🏠 社交领域 (Social Domain)
- **职责**: 聊天、好友、家族、队伍系统
- **核心实体**: Chat, Friend, Guild, Team, Mail
- **主要功能**: 社交互动、组队协作、消息通信

#### 🏗️ 建筑领域 (Building Domain)
- **职责**: 建筑系统、家园管理、建筑升级
- **核心实体**: Building, BuildingTemplate, BuildingUpgrade
- **主要功能**: 建筑建造、升级、功能管理

#### 🐾 宠物领域 (Pet Domain)
- **职责**: 宠物系统、宠物培养、宠物战斗
- **核心实体**: Pet, PetTemplate, PetSkill
- **主要功能**: 宠物获取、培养、进化、战斗辅助

#### 🏆 排行榜领域 (Ranking Domain)
- **职责**: 各类排行榜、积分统计、奖励发放
- **核心实体**: Ranking, RankingEntry, RankingReward
- **主要功能**: 排名计算、榜单更新、奖励分发

#### 🎯 小游戏领域 (Minigame Domain)
- **职责**: 各种小游戏、活动玩法、特殊奖励
- **核心实体**: Minigame, MinigameSession, MinigameReward
- **主要功能**: 小游戏逻辑、积分计算、奖励发放

### 支撑领域 (Supporting Domains)

#### 🔐 认证与授权
- JWT令牌管理
- 用户权限控制
- 安全策略实施

#### 📊 监控与日志
- 性能指标收集
- 业务日志记录
- 系统健康检查

#### ⚙️ 配置管理
- 多环境配置
- 动态配置更新
- 配置验证

## 🌐 网络协议设计

### 多协议支持
- **TCP**: 主要游戏协议，低延迟、高可靠性
- **HTTP**: RESTful API，用于管理后台和第三方集成
- **gRPC**: 微服务间通信，高性能RPC调用
- **WebSocket**: Web客户端支持，实时双向通信

### TCP协议格式
```
+--------+--------+--------+----------+
| Magic  | Length | Type   | Data     |
| 2bytes | 4bytes | 2bytes | Variable |
+--------+--------+--------+----------+
```

### 消息分类
- **0x1xxx**: 系统消息 (登录、心跳、错误)
- **0x2xxx**: 玩家消息 (属性、状态、升级)
- **0x3xxx**: 社交消息 (聊天、好友、邮件)
- **0x4xxx**: 战斗消息 (技能、伤害、结果)
- **0x5xxx**: 建筑消息 (建造、升级、管理)
- **0x6xxx**: 宠物消息 (获取、培养、战斗)
- **0x7xxx**: 排行榜消息 (查询、更新、奖励)
- **0x8xxx**: 小游戏消息 (开始、操作、结算)
- **0x9xxx**: 管理消息 (GM命令、系统公告)

## 🗄️ 数据存储设计

### MongoDB 集合设计

#### 核心业务集合
- **players**: 玩家基础信息和状态
- **player_stats**: 玩家统计数据和属性
- **battles**: 战斗记录和结果
- **guilds**: 公会信息和成员关系
- **buildings**: 建筑数据和状态
- **pets**: 宠物信息和属性
- **rankings**: 排行榜数据和历史
- **minigames**: 小游戏记录和积分

#### 配置和模板集合
- **game_configs**: 游戏配置参数
- **item_templates**: 物品模板数据
- **skill_templates**: 技能模板配置
- **building_templates**: 建筑模板信息
- **pet_templates**: 宠物模板数据

#### 日志和审计集合
- **player_logs**: 玩家操作日志
- **battle_logs**: 战斗详细日志
- **admin_logs**: 管理操作日志
- **system_events**: 系统事件记录

### Redis 缓存策略

#### 热点数据缓存
- **在线玩家**: `online:players:{server_id}`
- **玩家会话**: `session:{player_id}`
- **排行榜**: `ranking:{type}:{period}`
- **公会信息**: `guild:{guild_id}`

#### 临时数据缓存
- **战斗状态**: `battle:{battle_id}`
- **队伍信息**: `team:{team_id}`
- **聊天频道**: `chat:{channel_id}`
- **活动状态**: `event:{event_id}`

#### 性能优化缓存
- **查询结果**: `query:{hash}` (TTL: 5分钟)
- **计算结果**: `calc:{type}:{id}` (TTL: 1小时)
- **配置数据**: `config:{key}` (TTL: 24小时)

## 👨‍💻 开发指南

### 🏗️ DDD开发模式

#### 添加新领域
1. 在 `internal/domain/` 下创建领域目录
2. 定义领域实体、值对象和聚合根
3. 实现领域服务和仓储接口
4. 创建对应的应用服务
5. 实现基础设施层的具体实现
6. 添加接口层的处理器

#### 领域开发规范
```go
// 领域实体示例
type Player struct {
    id       PlayerID
    name     string
    level    int
    exp      int64
    stats    PlayerStats
    // 领域行为
}

func (p *Player) LevelUp() error {
    // 领域逻辑实现
}
```

### 🔧 开发工具使用

#### Make命令
```bash
make setup      # 初始化开发环境
make dev        # 启动开发服务器
make build      # 构建生产版本
make test       # 运行测试
make lint       # 代码质量检查
make clean      # 清理构建产物
make docs       # 生成文档
```

#### 代码生成
```bash
# 生成领域模板
scripts/generate-domain.sh <domain_name>

# 生成API接口
scripts/generate-api.sh <api_name>

# 生成数据库迁移
scripts/generate-migration.sh <migration_name>
```

### 📊 性能优化策略

#### 数据库优化
- **连接池管理**: 合理配置MongoDB和Redis连接池
- **索引优化**: 为查询频繁的字段创建合适索引
- **分片策略**: 大数据量集合采用分片存储
- **读写分离**: 读操作使用从库，写操作使用主库

#### 缓存策略
- **多级缓存**: 内存缓存 + Redis缓存 + 数据库
- **缓存预热**: 服务启动时预加载热点数据
- **缓存更新**: 采用Cache-Aside模式更新缓存
- **缓存穿透**: 使用布隆过滤器防止缓存穿透

#### 网络优化
- **连接复用**: TCP连接池和HTTP Keep-Alive
- **消息批处理**: 批量处理非实时消息
- **压缩传输**: 大数据包启用压缩
- **协议优化**: 使用二进制协议减少传输开销

### 📈 监控与运维

#### 日志管理
- **结构化日志**: 使用JSON格式便于解析
- **日志分级**: ERROR/WARN/INFO/DEBUG四个级别
- **日志轮转**: 按大小和时间自动轮转
- **敏感信息**: 避免记录密码等敏感数据

#### 指标监控
- **业务指标**: 在线人数、注册量、收入等
- **性能指标**: 响应时间、吞吐量、错误率
- **系统指标**: CPU、内存、磁盘、网络使用率
- **自定义指标**: 游戏特定的业务指标

#### 健康检查
```go
// HTTP健康检查接口
GET /health
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "services": {
    "database": "healthy",
    "redis": "healthy",
    "nats": "healthy"
  }
}
```

## 🚀 部署指南

### 🐳 Docker部署

#### 单容器部署
```bash
# 构建镜像
make docker-build

# 运行容器
docker run -d \
  --name greatestworks \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  -e ENV=production \
  greatestworks:latest
```

#### Docker Compose部署
```bash
# 启动完整环境（包含MongoDB、Redis、NATS）
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f greatestworks
```

### ☸️ Kubernetes部署

#### 基础部署
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: greatestworks
  namespace: gaming
spec:
  replicas: 3
  selector:
    matchLabels:
      app: greatestworks
  template:
    metadata:
      labels:
        app: greatestworks
    spec:
      containers:
      - name: server
        image: greatestworks:latest
        ports:
        - containerPort: 8080
        env:
        - name: ENV
          value: "production"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
```

#### 服务暴露
```yaml
apiVersion: v1
kind: Service
metadata:
  name: greatestworks-service
spec:
  selector:
    app: greatestworks
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

### 🔧 生产环境配置

#### 环境变量
```bash
# 服务配置
export SERVER_PORT=8080
export SERVER_HOST=0.0.0.0

# 数据库配置
export MONGODB_URI="mongodb://mongo-cluster:27017/gamedb"
export REDIS_ADDR="redis-cluster:6379"

# 消息队列
export NATS_URL="nats://nats-cluster:4222"

# 认证配置
export JWT_SECRET="your-production-secret-key"

# 日志配置
export LOG_LEVEL=info
export LOG_FORMAT=json
```

## 📚 API文档

详细的API文档请参考：
- [REST API文档](docs/api/rest-api.md)
- [TCP协议文档](docs/api/tcp-protocol.md)
- [WebSocket API文档](docs/api/websocket-api.md)

## 🏗️ 架构文档

深入了解系统架构：
- [DDD设计文档](docs/architecture/ddd-design.md)
- [数据库设计](docs/architecture/database-design.md)
- [微服务架构](docs/architecture/microservices.md)

## 🤝 贡献指南

我们欢迎所有形式的贡献！请阅读 [CONTRIBUTING.md](CONTRIBUTING.md) 了解详细信息。

### 贡献流程
1. **Fork** 项目到你的GitHub账户
2. **创建** 功能分支 (`git checkout -b feature/amazing-feature`)
3. **提交** 你的更改 (`git commit -m 'Add some amazing feature'`)
4. **推送** 到分支 (`git push origin feature/amazing-feature`)
5. **创建** Pull Request

### 开发规范
- 遵循 [Go代码规范](https://golang.org/doc/effective_go.html)
- 编写单元测试，保持测试覆盖率 > 80%
- 更新相关文档
- 通过所有CI检查

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 联系我们

- **项目主页**: [https://github.com/your-org/greatestworks](https://github.com/your-org/greatestworks)
- **问题反馈**: [GitHub Issues](https://github.com/your-org/greatestworks/issues)
- **讨论交流**: [GitHub Discussions](https://github.com/your-org/greatestworks/discussions)
- **邮箱**: dev@greatestworks.com
- **文档站点**: [https://docs.greatestworks.com](https://docs.greatestworks.com)

## 📈 项目状态

![Build Status](https://github.com/your-org/greatestworks/workflows/CI/badge.svg)
![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Docker Pulls](https://img.shields.io/docker/pulls/greatestworks/server.svg)

## 🎯 路线图

### v2.0.0 (计划中)
- [ ] 微服务拆分和服务网格
- [ ] GraphQL API支持
- [ ] 实时数据分析和BI
- [ ] 多语言客户端SDK
- [ ] 云原生部署优化

### v1.5.0 (开发中)
- [ ] WebSocket API完善
- [ ] 管理后台界面
- [ ] 性能监控面板
- [ ] 自动化测试覆盖

### v1.0.0 ✅ (已发布)
- [x] DDD架构重构完成
- [x] 核心游戏系统实现
- [x] Docker容器化支持
- [x] 基础监控和日志
- [x] 完整文档体系

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给我们一个Star！⭐**

*Built with ❤️ by the Greatest Works Team*

</div>