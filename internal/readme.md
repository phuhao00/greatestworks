
# Internal Module Documentation

## 📋 模块说明

### 🏗️ 模块架构设计

每个模块遵循DDD（领域驱动设计）架构，包含以下核心组件：

#### 📦 核心组件
- **Model**: 对应的数据存储和领域模型
- **System**: 该模块的管理系统，负责数据的CRUD操作
- **Owner**: 定义从属模块需要实现的方法接口
- **Handler**: 处理从属模块需要的业务逻辑
- **Abstract**: 模块成员的抽象，接口定义

### 🎯 领域驱动设计 (DDD)

本项目采用完整的DDD架构模式，将业务逻辑分为以下层次：

#### 🏛️ 领域层 (Domain Layer)
```
domain/
├── player/           # 玩家领域
│   ├── beginner/     # 新手引导系统
│   ├── hangup/       # 挂机系统
│   ├── honor/        # 荣誉系统
│   ├── player.go     # 玩家聚合根
│   ├── service.go    # 领域服务
│   └── repository.go # 仓储接口
├── battle/           # 战斗领域
├── social/           # 社交领域 (31个文件)
├── building/         # 建筑领域
├── pet/              # 宠物领域
├── ranking/          # 排行榜领域
├── minigame/         # 小游戏领域
├── npc/              # NPC领域
├── quest/            # 任务领域
├── scene/            # 场景领域 (24个文件)
├── skill/            # 技能领域
├── inventory/        # 背包领域
│   ├── dressup/      # 装扮系统
│   └── synthesis/    # 合成系统
└── events/           # 领域事件
```

#### 🏗️ 基础设施层 (Infrastructure Layer)
```
infrastructure/
├── persistence/      # 数据持久化 (10个文件)
│   ├── base_repository.go    # 基础仓储
│   ├── player_repository.go  # 玩家仓储
│   ├── battle_repository.go  # 战斗仓储
│   ├── hangup_repository.go  # 挂机仓储
│   ├── weather_repository.go # 天气仓储
│   ├── plant_repository.go   # 植物仓储
│   └── npc_repository.go     # NPC仓储
├── cache/            # 缓存服务
├── messaging/        # 消息服务 (5个文件)
│   ├── nats_publisher.go    # NATS发布者
│   ├── nats_subscriber.go   # NATS订阅者
│   ├── event_dispatcher.go  # 事件分发器
│   └── worker_pool.go       # 工作池
├── network/          # 网络服务
├── config/           # 配置管理 (7个文件)
├── logging/          # 日志服务
├── auth/            # 认证服务
├── container/       # 依赖注入容器
└── monitoring/      # 监控服务
```

#### 🌐 接口层 (Interface Layer)
```
interfaces/
├── http/             # HTTP接口 (13个文件)
│   ├── auth/         # 认证接口
│   ├── gm/           # GM管理接口
│   └── server.go     # HTTP服务器
├── tcp/              # TCP接口 (14个文件)
│   ├── handlers/     # TCP处理器
│   ├── connection/   # 连接管理
│   └── protocol/     # 协议定义
└── rpc/              # RPC接口 (4个文件)
```

### 🔧 核心模块说明

#### 📊 模块管理器
- **module_manager.go**: 模块生命周期管理
- **imodule.go**: 模块接口定义
- **base_module.go**: 基础模块实现

#### 🎮 游戏核心
- **game/**: 游戏核心逻辑
- **events/**: 事件系统
- **errors/**: 错误处理

#### 🌐 网络通信
- **network/**: 网络协议处理
- **proto/**: 协议定义

#### 🗄️ 数据存储
- **database/**: 数据库连接
- **config/**: 配置管理
- **auth/**: 认证系统

### 🚀 开发指南

#### 添加新领域模块
1. 在 `domain/` 下创建新领域目录
2. 定义领域实体、值对象和聚合根
3. 实现领域服务和仓储接口
4. 在 `infrastructure/` 下实现具体实现
5. 在 `interfaces/` 下添加接口层

#### 模块开发规范
- 遵循DDD架构原则
- 使用依赖注入容器
- 实现统一的错误处理
- 采用结构化日志
- 编写完整的单元测试

### 📈 性能优化

#### 数据库优化
- 使用连接池管理数据库连接
- 实现读写分离策略
- 合理使用缓存机制

#### 网络优化
- TCP连接复用
- 消息批处理
- 协议压缩

#### 内存优化
- 对象池复用
- 合理的内存分配
- 垃圾回收优化
