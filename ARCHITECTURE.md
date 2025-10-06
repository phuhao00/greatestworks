# Greatest Works - 架构设计文档

## 🎯 架构概述

本项目采用领域驱动设计(Domain-Driven Design, DDD)架构模式，构建了一个高度模块化、可扩展的MMO游戏服务器系统。通过清晰的分层架构和领域划分，实现了业务逻辑与技术实现的有效分离，提供了优秀的可维护性和可测试性。

### 核心设计原则

- **领域驱动**: 以业务领域为核心，将复杂的游戏逻辑按领域进行组织
- **分层架构**: 清晰的职责分离，每层只关注自己的核心职责
- **依赖倒置**: 高层模块不依赖低层模块，都依赖于抽象
- **单一职责**: 每个组件都有明确的单一职责
- **开闭原则**: 对扩展开放，对修改封闭
- **接口隔离**: 客户端不应该依赖它不需要的接口

## 🏗️ DDD分层架构

### 架构分层图

```
┌─────────────────────────────────────────────────────────────┐
│                    接口层 (Interfaces)                      │
│  ┌─────────────┐  ┌─────────────┐                         │
│  │   TCP API   │  │  HTTP API   │                         │
│  └─────────────┘  └─────────────┘                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    应用层 (Application)                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Commands  │  │   Queries   │  │  Services   │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     领域层 (Domain)                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Entities  │  │ Aggregates  │  │   Services  │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │ Value Objs  │  │ Repositories│  │    Events   │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  基础设施层 (Infrastructure)                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │  Database   │  │    Cache    │  │  Messaging  │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Network   │  │   Config    │  │   Logging   │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
```

### 项目目录结构

```
greatestworks/
├── cmd/                              # 应用程序入口
│   └── server/
│       ├── main.go                   # 主程序入口
│       └── bootstrap.go              # 启动引导系统
├── configs/                          # 配置模板
│   ├── config.example.yaml           # 基础配置模板
│   ├── config.dev.yaml.example       # 开发环境配置
│   └── config.prod.yaml.example      # 生产环境配置
├── docs/                             # 项目文档
│   ├── api/                          # API文档
│   ├── architecture/                 # 架构文档
│   └── diagrams/                     # 架构图表
├── scripts/                          # 开发脚本
│   ├── build.sh                      # 构建脚本
│   ├── deploy.sh                     # 部署脚本
│   ├── test.sh                       # 测试脚本
│   └── setup-dev.sh                  # 开发环境设置
├── application/                      # 应用层
│   ├── commands/                     # 命令处理器
│   ├── handlers/                     # 事件处理器
│   ├── queries/                      # 查询处理器
│   └── services/                     # 应用服务
│       ├── player_service.go         # 玩家应用服务
│       ├── battle_service.go         # 战斗应用服务
│       ├── social_service.go         # 社交应用服务
│       └── ...
├── internal/                         # 内部模块
│   ├── domain/                       # 领域层
│   │   ├── player/                   # 玩家领域
│   │   │   ├── aggregate.go          # 聚合根
│   │   │   ├── entity.go             # 实体
│   │   │   ├── value_object.go       # 值对象
│   │   │   ├── repository.go         # 仓储接口
│   │   │   ├── service.go            # 领域服务
│   │   │   └── events.go             # 领域事件
│   │   ├── battle/                   # 战斗领域
│   │   ├── social/                   # 社交领域
│   │   ├── building/                 # 建筑领域
│   │   ├── pet/                      # 宠物领域
│   │   ├── ranking/                  # 排行榜领域
│   │   └── minigame/                 # 小游戏领域
│   ├── infrastructure/               # 基础设施层
│   │   ├── persistence/              # 数据持久化
│   │   │   ├── mongodb/              # MongoDB实现
│   │   │   ├── redis/                # Redis实现
│   │   │   └── repositories/         # 仓储实现
│   │   ├── cache/                    # 缓存服务
│   │   ├── messaging/                # 消息服务
│   │   ├── network/                  # 网络服务
│   │   ├── config/                   # 配置管理
│   │   ├── logging/                  # 日志服务
│   │   ├── monitoring/               # 监控服务
│   │   ├── container/                # 依赖注入容器
│   │   └── protocol/                 # 协议管理
│   └── interfaces/                   # 接口层
│       ├── tcp/                      # TCP接口
│       │   ├── handlers/             # TCP处理器
│       │   └── protocol/             # TCP协议
│       ├── http/                     # HTTP接口
│       │   ├── controllers/          # HTTP控制器
│       │   ├── middleware/           # HTTP中间件
│       │   └── routes/               # 路由定义

├── migrations/                       # 数据库迁移
├── seeds/                           # 种子数据
├── docker-compose.yml               # Docker编排
├── Dockerfile                       # Docker镜像
├── Makefile                        # 构建工具
└── go.mod                          # Go模块定义
```

## 📋 DDD各层职责详解

### 1. 接口层 (Interfaces Layer)

**职责**: 处理外部请求，协议转换，输入验证

#### 核心组件
- **TCP处理器**: 处理游戏客户端的TCP连接和消息
- **HTTP控制器**: 提供RESTful API，主要用于管理后台

- **协议转换器**: 将外部协议转换为内部领域对象

#### 设计原则
- 薄接口层，只负责协议转换和基本验证
- 不包含业务逻辑，所有业务操作委托给应用层
- 统一的错误处理和响应格式
- 支持多种协议和数据格式

```go
// TCP处理器示例
type PlayerHandler struct {
    playerService application.PlayerService
}

func (h *PlayerHandler) HandleLogin(ctx context.Context, req *protocol.LoginRequest) (*protocol.LoginResponse, error) {
    // 协议验证
    if err := req.Validate(); err != nil {
        return nil, err
    }
    
    // 委托给应用服务
    result, err := h.playerService.Login(ctx, &application.LoginCommand{
        Username: req.Username,
        Password: req.Password,
    })
    
    // 转换为协议响应
    return &protocol.LoginResponse{
        Success: result.Success,
        Token:   result.Token,
        Player:  convertToProtocolPlayer(result.Player),
    }, err
}
```

### 2. 应用层 (Application Layer)

**职责**: 协调领域对象，处理业务用例，事务管理

#### 核心组件
- **应用服务**: 实现具体的业务用例
- **命令处理器**: 处理修改操作的命令
- **查询处理器**: 处理只读查询操作
- **事件处理器**: 处理领域事件
- **DTO对象**: 数据传输对象

#### 设计原则
- 薄应用层，主要负责协调和编排
- 事务边界的管理
- 领域事件的发布和处理
- 不包含业务规则，委托给领域层

```go
// 应用服务示例
type PlayerService struct {
    playerRepo   domain.PlayerRepository
    eventBus     events.EventBus
    unitOfWork   persistence.UnitOfWork
}

func (s *PlayerService) CreatePlayer(ctx context.Context, cmd *CreatePlayerCommand) (*CreatePlayerResult, error) {
    return s.unitOfWork.Execute(ctx, func(ctx context.Context) (*CreatePlayerResult, error) {
        // 创建领域对象
        player, err := domain.NewPlayer(cmd.Name, cmd.Class)
        if err != nil {
            return nil, err
        }
        
        // 保存到仓储
        if err := s.playerRepo.Save(ctx, player); err != nil {
            return nil, err
        }
        
        // 发布领域事件
        s.eventBus.Publish(ctx, player.GetEvents()...)
        
        return &CreatePlayerResult{
            PlayerID: player.ID(),
            Name:     player.Name(),
        }, nil
    })
}
```

### 3. 领域层 (Domain Layer)

**职责**: 核心业务逻辑，业务规则，领域模型

#### 核心组件
- **聚合根**: 保证数据一致性的边界
- **实体**: 具有唯一标识的领域对象
- **值对象**: 不可变的描述性对象
- **领域服务**: 跨聚合的业务逻辑
- **仓储接口**: 数据访问抽象
- **领域事件**: 领域内重要事件

#### 领域划分

##### 玩家领域 (Player Domain)
```go
// 玩家聚合根
type Player struct {
    id       PlayerID
    name     string
    level    int
    exp      int64
    stats    PlayerStats
    events   []events.DomainEvent
}

func (p *Player) LevelUp() error {
    if !p.CanLevelUp() {
        return errors.New("insufficient experience")
    }
    
    p.level++
    p.stats = p.stats.RecalculateForLevel(p.level)
    
    // 发布领域事件
    p.AddEvent(&PlayerLevelUpEvent{
        PlayerID: p.id,
        NewLevel: p.level,
        OccurredAt: time.Now(),
    })
    
    return nil
}
```

##### 战斗领域 (Battle Domain)
```go
// 战斗聚合根
type Battle struct {
    id          BattleID
    attacker    PlayerID
    defender    PlayerID
    status      BattleStatus
    rounds      []BattleRound
    result      *BattleResult
}

func (b *Battle) ExecuteRound(attackerAction, defenderAction Action) error {
    if b.status != BattleStatusInProgress {
        return errors.New("battle is not in progress")
    }
    
    round := NewBattleRound(attackerAction, defenderAction)
    round.Execute()
    
    b.rounds = append(b.rounds, round)
    
    if round.IsDecisive() {
        b.EndBattle(round.Winner())
    }
    
    return nil
}
```

### 4. 基础设施层 (Infrastructure Layer)

**职责**: 技术实现，外部系统集成，数据持久化

#### 核心组件
- **仓储实现**: 具体的数据访问实现
- **消息队列**: 异步消息处理
- **缓存服务**: 性能优化
- **配置管理**: 系统配置
- **日志服务**: 系统监控
- **网络服务**: 底层网络通信

```go
// MongoDB仓储实现
type MongoPlayerRepository struct {
    collection *mongo.Collection
}

func (r *MongoPlayerRepository) Save(ctx context.Context, player *domain.Player) error {
    doc := r.toDocument(player)
    
    _, err := r.collection.ReplaceOne(
        ctx,
        bson.M{"_id": player.ID()},
        doc,
        options.Replace().SetUpsert(true),
    )
    
    return err
}

func (r *MongoPlayerRepository) FindByID(ctx context.Context, id domain.PlayerID) (*domain.Player, error) {
    var doc playerDocument
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
    if err != nil {
        return nil, err
    }
    
    return r.toDomain(&doc), nil
}
```

## 🎯 架构决策记录 (ADR)

### ADR-001: 采用DDD架构模式

**状态**: 已接受  
**日期**: 2024-01-15  
**决策者**: 架构团队

#### 背景
项目需要处理复杂的游戏业务逻辑，包括玩家系统、战斗系统、社交系统等多个领域。传统的分层架构难以应对业务复杂性。

#### 决策
采用领域驱动设计(DDD)架构模式，按业务领域组织代码结构。

#### 理由
- **业务复杂性**: 游戏业务逻辑复杂，需要清晰的领域划分
- **团队协作**: 不同团队可以专注于不同的业务领域
- **可维护性**: 业务逻辑集中在领域层，易于理解和维护
- **可测试性**: 领域逻辑与技术实现分离，便于单元测试

#### 后果
- **正面**: 代码组织清晰，业务逻辑集中，易于维护和扩展
- **负面**: 学习成本较高，需要团队对DDD有深入理解

### ADR-002: 使用MongoDB作为主数据库

**状态**: 已接受  
**日期**: 2024-01-15  
**决策者**: 技术团队

#### 背景
需要选择合适的数据库来存储游戏数据，包括玩家信息、游戏状态等。

#### 决策
使用MongoDB作为主数据库，Redis作为缓存。

#### 理由
- **灵活性**: 文档数据库适合存储复杂的游戏对象
- **扩展性**: 支持水平扩展，适合大规模游戏
- **性能**: 读写性能优秀，适合游戏场景
- **开发效率**: 与Go语言集成良好

#### 后果
- **正面**: 开发效率高，性能优秀，扩展性好
- **负面**: 事务支持相对较弱，需要在应用层处理一致性

### ADR-003: 采用事件驱动架构

**状态**: 已接受  
**日期**: 2024-01-15  
**决策者**: 架构团队

#### 背景
游戏系统中存在大量的异步处理需求，如经验获得、成就解锁、排行榜更新等。

#### 决策
在DDD架构基础上，采用事件驱动架构处理异步业务逻辑。

#### 理由
- **解耦**: 事件发布者和订阅者解耦，提高系统灵活性
- **扩展性**: 新功能可以通过订阅事件实现，无需修改现有代码
- **一致性**: 通过事件确保最终一致性
- **审计**: 事件流提供完整的业务操作记录

#### 后果
- **正面**: 系统解耦，易于扩展，支持复杂的业务流程
- **负面**: 调试复杂度增加，需要处理事件的顺序和重复问题

### ADR-004: 使用依赖注入容器

**状态**: 已接受  
**日期**: 2024-01-15  
**决策者**: 开发团队

#### 背景
系统中存在大量的依赖关系，需要一种优雅的方式管理这些依赖。

#### 决策
实现自定义的依赖注入容器，管理服务的生命周期和依赖关系。

#### 理由
- **解耦**: 减少组件间的直接依赖
- **测试**: 便于进行单元测试和集成测试
- **配置**: 集中管理服务配置
- **生命周期**: 统一管理服务的创建和销毁

#### 后果
- **正面**: 代码解耦，易于测试，配置集中
- **负面**: 增加了系统复杂度，需要额外的学习成本

## 🛠️ 开发指南

### 添加新领域的步骤

1. **创建领域目录结构**
```bash
mkdir -p internal/domain/newdomain
touch internal/domain/newdomain/{aggregate.go,entity.go,value_object.go,repository.go,service.go,events.go}
```

2. **定义领域模型**
```go
// internal/domain/newdomain/aggregate.go
type NewDomainAggregate struct {
    id     NewDomainID
    // 其他字段
    events []events.DomainEvent
}

func (a *NewDomainAggregate) DoSomething() error {
    // 业务逻辑
    a.AddEvent(&SomethingHappenedEvent{
        AggregateID: a.id,
        OccurredAt:  time.Now(),
    })
    return nil
}
```

3. **创建应用服务**
```go
// application/services/newdomain_service.go
type NewDomainService struct {
    repo       domain.NewDomainRepository
    eventBus   events.EventBus
    unitOfWork persistence.UnitOfWork
}
```

4. **实现基础设施层**
```go
// internal/infrastructure/persistence/repositories/newdomain_repository.go
type MongoNewDomainRepository struct {
    collection *mongo.Collection
}
```

5. **添加接口层处理器**
```go
// internal/interfaces/tcp/handlers/newdomain_handler.go
type NewDomainHandler struct {
    service application.NewDomainService
}
```

### 代码规范

#### 命名约定
- **包名**: 小写，简短，描述性强
- **接口**: 以"er"结尾，如`Repository`、`Service`
- **结构体**: 帕斯卡命名法，如`PlayerService`
- **方法**: 驼峰命名法，动词开头，如`CreatePlayer`
- **常量**: 全大写，下划线分隔，如`MAX_LEVEL`

#### 错误处理
```go
// 使用自定义错误类型
type DomainError struct {
    Code    string
    Message string
    Cause   error
}

func (e *DomainError) Error() string {
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// 错误包装
func (s *PlayerService) CreatePlayer(ctx context.Context, cmd *CreatePlayerCommand) error {
    if err := s.validateCommand(cmd); err != nil {
        return &DomainError{
            Code:    "INVALID_COMMAND",
            Message: "Invalid create player command",
            Cause:   err,
        }
    }
    // ...
}
```

#### 日志记录
```go
// 结构化日志
logger.WithFields(map[string]interface{}{
    "player_id": playerID,
    "action":    "level_up",
    "old_level": oldLevel,
    "new_level": newLevel,
}).Info("Player leveled up")

// 错误日志
logger.WithError(err).WithField("player_id", playerID).Error("Failed to save player")
```

### 测试策略

#### 单元测试
```go
// 领域层测试
func TestPlayer_LevelUp(t *testing.T) {
    // Given
    player := domain.NewPlayer("TestPlayer", domain.ClassWarrior)
    player.AddExperience(1000)
    
    // When
    err := player.LevelUp()
    
    // Then
    assert.NoError(t, err)
    assert.Equal(t, 2, player.Level())
    assert.Len(t, player.GetEvents(), 1)
}

// 应用服务测试
func TestPlayerService_CreatePlayer(t *testing.T) {
    // Given
    mockRepo := &mocks.PlayerRepository{}
    mockEventBus := &mocks.EventBus{}
    service := application.NewPlayerService(mockRepo, mockEventBus)
    
    // When
    result, err := service.CreatePlayer(context.Background(), &application.CreatePlayerCommand{
        Name:  "TestPlayer",
        Class: "Warrior",
    })
    
    // Then
    assert.NoError(t, err)
    assert.NotEmpty(t, result.PlayerID)
    mockRepo.AssertCalled(t, "Save", mock.Anything, mock.Anything)
}
```

#### 集成测试
```go
func TestPlayerIntegration(t *testing.T) {
    // 设置测试数据库
    testDB := setupTestDatabase(t)
    defer cleanupTestDatabase(t, testDB)
    
    // 创建真实的仓储实现
    repo := persistence.NewMongoPlayerRepository(testDB)
    service := application.NewPlayerService(repo, events.NewInMemoryEventBus())
    
    // 执行集成测试
    result, err := service.CreatePlayer(context.Background(), &application.CreatePlayerCommand{
        Name:  "IntegrationTestPlayer",
        Class: "Mage",
    })
    
    assert.NoError(t, err)
    
    // 验证数据已保存到数据库
    savedPlayer, err := repo.FindByID(context.Background(), result.PlayerID)
    assert.NoError(t, err)
    assert.Equal(t, "IntegrationTestPlayer", savedPlayer.Name())
}
```

## 🎯 性能优化指南

### 数据库优化

#### 索引策略
```javascript
// MongoDB索引创建
db.players.createIndex({ "user_id": 1 }, { unique: true })
db.players.createIndex({ "name": 1 }, { unique: true })
db.players.createIndex({ "level": -1 })
db.players.createIndex({ "guild_id": 1, "level": -1 })

// 复合索引用于复杂查询
db.battles.createIndex({ "player_id": 1, "created_at": -1 })
db.rankings.createIndex({ "type": 1, "score": -1, "updated_at": -1 })
```

#### 查询优化
```go
// 使用投影减少数据传输
func (r *MongoPlayerRepository) FindPlayerSummary(ctx context.Context, id PlayerID) (*PlayerSummary, error) {
    projection := bson.M{
        "name":  1,
        "level": 1,
        "class": 1,
    }
    
    var doc playerSummaryDocument
    err := r.collection.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(projection)).Decode(&doc)
    return r.toSummary(&doc), err
}

// 批量操作
func (r *MongoPlayerRepository) SaveBatch(ctx context.Context, players []*domain.Player) error {
    var operations []mongo.WriteModel
    
    for _, player := range players {
        doc := r.toDocument(player)
        operation := mongo.NewReplaceOneModel()
        operation.SetFilter(bson.M{"_id": player.ID()})
        operation.SetReplacement(doc)
        operation.SetUpsert(true)
        operations = append(operations, operation)
    }
    
    _, err := r.collection.BulkWrite(ctx, operations)
    return err
}
```

### 缓存策略

#### Redis缓存模式
```go
// Cache-Aside模式
func (s *PlayerService) GetPlayer(ctx context.Context, id PlayerID) (*domain.Player, error) {
    // 先查缓存
    if cached, err := s.cache.Get(ctx, fmt.Sprintf("player:%s", id)); err == nil {
        return s.deserializePlayer(cached), nil
    }
    
    // 缓存未命中，查数据库
    player, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 写入缓存
    s.cache.Set(ctx, fmt.Sprintf("player:%s", id), s.serializePlayer(player), time.Hour)
    
    return player, nil
}

// Write-Through模式
func (s *PlayerService) UpdatePlayer(ctx context.Context, player *domain.Player) error {
    // 同时更新数据库和缓存
    if err := s.repo.Save(ctx, player); err != nil {
        return err
    }
    
    return s.cache.Set(ctx, fmt.Sprintf("player:%s", player.ID()), s.serializePlayer(player), time.Hour)
}
```

### 并发控制

#### 乐观锁
```go
type Player struct {
    id      PlayerID
    version int64  // 版本号
    // 其他字段
}

func (r *MongoPlayerRepository) Save(ctx context.Context, player *domain.Player) error {
    filter := bson.M{
        "_id":     player.ID(),
        "version": player.Version(),
    }
    
    update := bson.M{
        "$set": r.toDocument(player),
        "$inc": bson.M{"version": 1},
    }
    
    result, err := r.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }
    
    if result.ModifiedCount == 0 {
        return errors.New("optimistic lock failed")
    }
    
    return nil
}
```

## 🔒 安全最佳实践

### 输入验证
```go
// 使用验证器
type CreatePlayerCommand struct {
    Name  string `validate:"required,min=3,max=20,alphanum"`
    Class string `validate:"required,oneof=warrior mage archer"`
}

func (s *PlayerService) CreatePlayer(ctx context.Context, cmd *CreatePlayerCommand) error {
    if err := s.validator.Struct(cmd); err != nil {
        return &ValidationError{Cause: err}
    }
    // ...
}
```

### 权限控制
```go
// 基于角色的访问控制
type Permission string

const (
    PermissionReadPlayer   Permission = "player:read"
    PermissionWritePlayer  Permission = "player:write"
    PermissionDeletePlayer Permission = "player:delete"
)

func (s *PlayerService) GetPlayer(ctx context.Context, id PlayerID) (*domain.Player, error) {
    if !s.authService.HasPermission(ctx, PermissionReadPlayer) {
        return nil, errors.New("insufficient permissions")
    }
    // ...
}
```

### 敏感数据处理
```go
// 密码哈希
func (s *AuthService) HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// 数据脱敏
func (p *Player) ToPublicView() *PublicPlayerView {
    return &PublicPlayerView{
        ID:    p.id,
        Name:  p.name,
        Level: p.level,
        // 不包含敏感信息如邮箱、IP等
    }
}
```

## 📈 监控与运维

### 应用监控

#### 业务指标
```go
// 定义业务指标
type GameMetrics struct {
    OnlinePlayersGauge    prometheus.Gauge
    LoginCounter          prometheus.Counter
    BattleHistogram       prometheus.Histogram
    LevelUpCounter        prometheus.CounterVec
}

func NewGameMetrics() *GameMetrics {
    return &GameMetrics{
        OnlinePlayersGauge: prometheus.NewGauge(prometheus.GaugeOpts{
            Name: "game_online_players_total",
            Help: "Current number of online players",
        }),
        LoginCounter: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "game_login_total",
            Help: "Total number of player logins",
        }),
        BattleHistogram: prometheus.NewHistogram(prometheus.HistogramOpts{
            Name: "game_battle_duration_seconds",
            Help: "Battle duration in seconds",
            Buckets: prometheus.DefBuckets,
        }),
        LevelUpCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
            Name: "game_level_up_total",
            Help: "Total number of level ups by class",
        }, []string{"class"}),
    }
}
```

#### 健康检查
```go
type HealthChecker struct {
    dbChecker    DatabaseHealthChecker
    cacheChecker CacheHealthChecker
}

func (h *HealthChecker) Check(ctx context.Context) *HealthStatus {
    status := &HealthStatus{
        Status:    "healthy",
        Timestamp: time.Now(),
        Services:  make(map[string]string),
    }
    
    // 检查数据库
    if err := h.dbChecker.Ping(ctx); err != nil {
        status.Status = "unhealthy"
        status.Services["database"] = "unhealthy"
    } else {
        status.Services["database"] = "healthy"
    }
    
    // 检查缓存
    if err := h.cacheChecker.Ping(ctx); err != nil {
        status.Status = "unhealthy"
        status.Services["cache"] = "unhealthy"
    } else {
        status.Services["cache"] = "healthy"
    }
    
    return status
}
```

### 日志管理

#### 结构化日志
```go
// 定义日志字段
type LogFields map[string]interface{}

func (l LogFields) WithPlayerID(id string) LogFields {
    l["player_id"] = id
    return l
}

func (l LogFields) WithAction(action string) LogFields {
    l["action"] = action
    return l
}

// 使用示例
logger.WithFields(LogFields{}.WithPlayerID("12345").WithAction("login")).Info("Player logged in")
```

## 🚀 部署架构

### 微服务部署

```yaml
# docker-compose.yml
version: '3.8'
services:
  game-server:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ENV=production
      - MONGODB_URI=mongodb://mongo:27017/gamedb
      - REDIS_ADDR=redis:6379
    depends_on:
      - mongo
      - redis
      - nats
    
  mongo:
    image: mongo:5.0
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db
    
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    
  nats:
    image: nats:2.9-alpine
    ports:
      - "4222:4222"
    
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    
volumes:
  mongo_data:
  redis_data:
```

### Kubernetes部署

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: game-server
spec:
  replicas: 3
  selector:
    matchLabels:
      app: game-server
  template:
    metadata:
      labels:
        app: game-server
    spec:
      containers:
      - name: game-server
        image: greatestworks:latest
        ports:
        - containerPort: 8080
        env:
        - name: MONGODB_URI
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: mongodb-uri
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
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

## 📚 总结

### 架构优势

1. **清晰的职责分离**: DDD分层架构确保每层都有明确的职责
2. **高度可测试**: 依赖注入和接口抽象使得单元测试和集成测试更容易
3. **易于扩展**: 新功能可以通过添加新的领域或扩展现有领域来实现
4. **技术无关性**: 领域层不依赖具体的技术实现
5. **团队协作**: 不同团队可以并行开发不同的领域

### 技术选型理由

- **Go语言**: 高性能、并发支持好、部署简单
- **MongoDB**: 文档数据库，适合复杂的游戏对象存储
- **Redis**: 高性能缓存，适合游戏场景的实时数据
- **NATS**: 轻量级消息队列，支持高并发
- **Docker**: 容器化部署，环境一致性
- **Kubernetes**: 容器编排，支持自动扩缩容

### 未来规划

1. **微服务拆分**: 将不同领域拆分为独立的微服务
2. **事件溯源**: 实现事件溯源模式，提供完整的审计日志
3. **CQRS**: 实现命令查询职责分离，优化读写性能
4. **分布式缓存**: 实现分布式缓存，支持更大规模
5. **服务网格**: 引入Istio等服务网格，提供更好的服务治理

### 开发团队建议

1. **学习DDD**: 团队成员需要深入理解DDD的概念和实践
2. **代码审查**: 建立严格的代码审查流程，确保架构一致性
3. **文档维护**: 及时更新架构文档和API文档
4. **监控告警**: 建立完善的监控和告警体系
5. **性能测试**: 定期进行性能测试，确保系统性能

---

**本文档将随着项目的发展持续更新，请定期查看最新版本。**