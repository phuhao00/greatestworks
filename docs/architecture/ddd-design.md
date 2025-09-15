# 领域驱动设计 (DDD) 架构文档

## 📖 概述

GreatestWorks 采用领域驱动设计 (Domain-Driven Design) 作为核心架构模式，通过深入理解游戏业务领域，构建了清晰的领域模型和架构边界。

## 🎯 DDD 核心概念

### 战略设计

#### 限界上下文 (Bounded Context)

```
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   玩家上下文      │  │   游戏上下文      │  │   社交上下文      │
│   Player        │  │   Game          │  │   Social        │
│                 │  │                 │  │                 │
│ • 账户管理       │  │ • 场景管理       │  │ • 好友系统       │
│ • 角色信息       │  │ • 战斗系统       │  │ • 聊天系统       │
│ • 等级经验       │  │ • 技能系统       │  │ • 邮件系统       │
└─────────────────┘  └─────────────────┘  └─────────────────┘

┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   物品上下文      │  │   任务上下文      │  │   排行上下文      │
│   Inventory     │  │   Quest         │  │   Ranking       │
│                 │  │                 │  │                 │
│ • 背包管理       │  │ • 任务系统       │  │ • 等级排行       │
│ • 装备系统       │  │ • 成就系统       │  │ • 财富排行       │
│ • 道具合成       │  │ • 奖励发放       │  │ • PVP 排行      │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

#### 上下文映射 (Context Mapping)

```
Player Context ──→ Game Context     (Customer/Supplier)
Player Context ──→ Social Context   (Shared Kernel)
Game Context   ──→ Inventory Context (Published Language)
Quest Context  ──→ Player Context   (Anticorruption Layer)
```

### 战术设计

#### 领域模型层次

```
实体 (Entity)
├── 聚合根 (Aggregate Root)
│   ├── Player (玩家)
│   ├── Scene (场景)
│   ├── Battle (战斗)
│   └── Guild (公会)
├── 值对象 (Value Object)
│   ├── Position (位置)
│   ├── Money (金钱)
│   ├── Experience (经验)
│   └── Attribute (属性)
└── 领域服务 (Domain Service)
    ├── BattleCalculator (战斗计算)
    ├── LevelCalculator (等级计算)
    └── RewardDistributor (奖励分发)
```

## 🏗️ 架构分层

### 四层架构

```go
// 1. 接口层 (Interfaces Layer)
package interfaces

// HTTP 处理器
type PlayerHandler struct {
    playerService *application.PlayerService
}

// TCP 处理器
type GameHandler struct {
    gameService *application.GameService
}

// 2. 应用层 (Application Layer)
package application

// 应用服务
type PlayerService struct {
    playerRepo domain.PlayerRepository
    eventBus   infrastructure.EventBus
}

// 命令对象
type CreatePlayerCommand struct {
    Username string
    Email    string
}

// 3. 领域层 (Domain Layer)
package domain

// 聚合根
type Player struct {
    id       PlayerID
    username Username
    level    Level
    exp      Experience
    events   []DomainEvent
}

// 领域服务
type LevelService interface {
    CalculateLevel(exp Experience) Level
    GetRequiredExp(level Level) Experience
}

// 4. 基础设施层 (Infrastructure Layer)
package infrastructure

// 仓储实现
type MongoPlayerRepository struct {
    collection *mongo.Collection
}

// 事件总线实现
type NATSEventBus struct {
    conn *nats.Conn
}
```

## 🎮 领域模型设计

### 玩家聚合 (Player Aggregate)

```go
// 玩家聚合根
type Player struct {
    // 标识
    id       PlayerID
    username Username
    
    // 基础属性
    level      Level
    experience Experience
    gold       Gold
    
    // 状态信息
    status     PlayerStatus
    location   Location
    
    // 时间信息
    createdAt  time.Time
    lastLogin  time.Time
    
    // 领域事件
    events []DomainEvent
}

// 玩家行为
func (p *Player) GainExperience(exp Experience) error {
    if exp <= 0 {
        return errors.New("experience must be positive")
    }
    
    oldLevel := p.level
    p.experience += exp
    
    // 检查升级
    newLevel := p.calculateLevel()
    if newLevel > oldLevel {
        p.level = newLevel
        p.addEvent(PlayerLevelUpEvent{
            PlayerID: p.id,
            OldLevel: oldLevel,
            NewLevel: newLevel,
        })
    }
    
    return nil
}

func (p *Player) MoveTo(location Location) error {
    if !p.canMoveTo(location) {
        return errors.New("cannot move to location")
    }
    
    oldLocation := p.location
    p.location = location
    
    p.addEvent(PlayerMovedEvent{
        PlayerID:    p.id,
        OldLocation: oldLocation,
        NewLocation: location,
    })
    
    return nil
}
```

### 战斗聚合 (Battle Aggregate)

```go
// 战斗聚合根
type Battle struct {
    id          BattleID
    battleType  BattleType
    participants []Participant
    status      BattleStatus
    startTime   time.Time
    endTime     *time.Time
    result      *BattleResult
    events      []DomainEvent
}

// 战斗行为
func (b *Battle) Start() error {
    if b.status != BattleStatusPending {
        return errors.New("battle already started")
    }
    
    if len(b.participants) < 2 {
        return errors.New("not enough participants")
    }
    
    b.status = BattleStatusInProgress
    b.startTime = time.Now()
    
    b.addEvent(BattleStartedEvent{
        BattleID:     b.id,
        Participants: b.participants,
        StartTime:    b.startTime,
    })
    
    return nil
}

func (b *Battle) Attack(attackerID PlayerID, targetID PlayerID, skillID SkillID) error {
    if b.status != BattleStatusInProgress {
        return errors.New("battle not in progress")
    }
    
    attacker := b.getParticipant(attackerID)
    target := b.getParticipant(targetID)
    
    if attacker == nil || target == nil {
        return errors.New("invalid participants")
    }
    
    // 计算伤害
    damage := b.calculateDamage(attacker, target, skillID)
    target.TakeDamage(damage)
    
    b.addEvent(AttackEvent{
        BattleID:   b.id,
        AttackerID: attackerID,
        TargetID:   targetID,
        SkillID:    skillID,
        Damage:     damage,
    })
    
    // 检查战斗结束
    if target.IsDead() {
        b.end(attacker)
    }
    
    return nil
}
```

### 值对象设计

```go
// 位置值对象
type Position struct {
    X float64 `json:"x"`
    Y float64 `json:"y"`
    Z float64 `json:"z"`
}

func (p Position) DistanceTo(other Position) float64 {
    dx := p.X - other.X
    dy := p.Y - other.Y
    dz := p.Z - other.Z
    return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (p Position) IsValid() bool {
    return p.X >= 0 && p.Y >= 0 && p.Z >= 0
}

// 金钱值对象
type Gold struct {
    amount int64
}

func NewGold(amount int64) (Gold, error) {
    if amount < 0 {
        return Gold{}, errors.New("gold amount cannot be negative")
    }
    return Gold{amount: amount}, nil
}

func (g Gold) Add(other Gold) Gold {
    return Gold{amount: g.amount + other.amount}
}

func (g Gold) Subtract(other Gold) (Gold, error) {
    if g.amount < other.amount {
        return Gold{}, errors.New("insufficient gold")
    }
    return Gold{amount: g.amount - other.amount}, nil
}

// 经验值对象
type Experience struct {
    points int64
}

func (e Experience) Add(points int64) Experience {
    return Experience{points: e.points + points}
}

func (e Experience) ToLevel() Level {
    // 经验转等级的计算逻辑
    level := int(math.Sqrt(float64(e.points)/100)) + 1
    return Level{value: level}
}
```

## 🔄 领域事件

### 事件定义

```go
// 领域事件接口
type DomainEvent interface {
    EventID() string
    EventType() string
    AggregateID() string
    OccurredOn() time.Time
    EventVersion() int
}

// 玩家升级事件
type PlayerLevelUpEvent struct {
    eventID     string
    playerID    PlayerID
    oldLevel    Level
    newLevel    Level
    occurredOn  time.Time
}

func (e PlayerLevelUpEvent) EventID() string     { return e.eventID }
func (e PlayerLevelUpEvent) EventType() string   { return "PlayerLevelUp" }
func (e PlayerLevelUpEvent) AggregateID() string { return e.playerID.String() }
func (e PlayerLevelUpEvent) OccurredOn() time.Time { return e.occurredOn }
func (e PlayerLevelUpEvent) EventVersion() int   { return 1 }

// 战斗结束事件
type BattleEndedEvent struct {
    eventID    string
    battleID   BattleID
    winner     PlayerID
    loser      PlayerID
    duration   time.Duration
    occurredOn time.Time
}
```

### 事件处理器

```go
// 事件处理器接口
type EventHandler interface {
    Handle(event DomainEvent) error
    CanHandle(eventType string) bool
}

// 玩家升级事件处理器
type PlayerLevelUpHandler struct {
    rewardService *RewardService
    notifyService *NotificationService
}

func (h *PlayerLevelUpHandler) Handle(event DomainEvent) error {
    levelUpEvent := event.(PlayerLevelUpEvent)
    
    // 发放升级奖励
    reward := h.rewardService.GetLevelUpReward(levelUpEvent.newLevel)
    err := h.rewardService.GiveReward(levelUpEvent.playerID, reward)
    if err != nil {
        return err
    }
    
    // 发送升级通知
    return h.notifyService.NotifyLevelUp(levelUpEvent.playerID, levelUpEvent.newLevel)
}

func (h *PlayerLevelUpHandler) CanHandle(eventType string) bool {
    return eventType == "PlayerLevelUp"
}
```

## 🏪 仓储模式

### 仓储接口

```go
// 玩家仓储接口
type PlayerRepository interface {
    Save(player *Player) error
    FindByID(id PlayerID) (*Player, error)
    FindByUsername(username string) (*Player, error)
    FindAll(criteria PlayerCriteria) ([]*Player, error)
    Delete(id PlayerID) error
}

// 查询条件
type PlayerCriteria struct {
    MinLevel    *Level
    MaxLevel    *Level
    Status      *PlayerStatus
    LastLoginAfter *time.Time
    Limit       int
    Offset      int
}
```

### 仓储实现

```go
// MongoDB 仓储实现
type MongoPlayerRepository struct {
    collection *mongo.Collection
}

func (r *MongoPlayerRepository) Save(player *Player) error {
    doc := r.toDocument(player)
    
    filter := bson.M{"_id": player.ID()}
    opts := options.Replace().SetUpsert(true)
    
    _, err := r.collection.ReplaceOne(context.Background(), filter, doc, opts)
    return err
}

func (r *MongoPlayerRepository) FindByID(id PlayerID) (*Player, error) {
    var doc playerDocument
    err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&doc)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, ErrPlayerNotFound
        }
        return nil, err
    }
    
    return r.fromDocument(doc), nil
}

// 文档映射
type playerDocument struct {
    ID        string    `bson:"_id"`
    Username  string    `bson:"username"`
    Level     int       `bson:"level"`
    Experience int64    `bson:"experience"`
    Gold      int64     `bson:"gold"`
    Status    string    `bson:"status"`
    CreatedAt time.Time `bson:"created_at"`
    LastLogin time.Time `bson:"last_login"`
}
```

## 🎯 领域服务

### 战斗计算服务

```go
// 战斗计算领域服务
type BattleCalculationService struct {
    skillRepo SkillRepository
}

func (s *BattleCalculationService) CalculateDamage(
    attacker *Player, 
    target *Player, 
    skillID SkillID,
) (Damage, error) {
    skill, err := s.skillRepo.FindByID(skillID)
    if err != nil {
        return 0, err
    }
    
    // 基础伤害计算
    baseDamage := attacker.GetAttack() * skill.GetDamageMultiplier()
    
    // 防御减免
    defense := target.GetDefense()
    actualDamage := baseDamage * (1 - defense/(defense+100))
    
    // 暴击计算
    if s.isCritical(attacker.GetCritRate()) {
        actualDamage *= attacker.GetCritDamage()
    }
    
    // 随机浮动
    variance := 0.1 // 10% 浮动
    randomFactor := 1 + (rand.Float64()-0.5)*variance
    
    return Damage(actualDamage * randomFactor), nil
}

func (s *BattleCalculationService) isCritical(critRate float64) bool {
    return rand.Float64() < critRate
}
```

### 等级计算服务

```go
// 等级计算领域服务
type LevelCalculationService struct{}

func (s *LevelCalculationService) CalculateLevel(exp Experience) Level {
    points := exp.Points()
    
    // 使用平方根公式计算等级
    level := int(math.Sqrt(float64(points)/100)) + 1
    
    // 等级上限
    if level > MaxLevel {
        level = MaxLevel
    }
    
    return NewLevel(level)
}

func (s *LevelCalculationService) GetRequiredExp(level Level) Experience {
    if level.Value() <= 1 {
        return NewExperience(0)
    }
    
    // 计算升到指定等级需要的经验
    required := int64(math.Pow(float64(level.Value()-1), 2) * 100)
    return NewExperience(required)
}

func (s *LevelCalculationService) GetExpToNextLevel(player *Player) Experience {
    currentLevel := player.Level()
    nextLevel := NewLevel(currentLevel.Value() + 1)
    
    requiredExp := s.GetRequiredExp(nextLevel)
    currentExp := player.Experience()
    
    return requiredExp.Subtract(currentExp)
}
```

## 📋 最佳实践

### 1. 聚合设计原则

- **小聚合**: 保持聚合尽可能小
- **一致性边界**: 聚合内强一致性，聚合间最终一致性
- **通过ID引用**: 聚合间通过ID引用，避免对象引用
- **事务边界**: 一个事务只修改一个聚合

### 2. 领域事件使用

- **业务含义**: 事件应该有明确的业务含义
- **不可变**: 事件一旦发生不可修改
- **异步处理**: 事件处理应该异步进行
- **幂等性**: 事件处理器应该是幂等的

### 3. 值对象设计

- **不可变性**: 值对象应该是不可变的
- **相等性**: 基于值的相等性比较
- **验证**: 在构造时进行验证
- **行为丰富**: 包含相关的业务行为

### 4. 仓储实现

- **接口分离**: 领域层定义接口，基础设施层实现
- **聚合完整性**: 保存和加载完整的聚合
- **查询优化**: 针对查询场景优化实现
- **缓存策略**: 合理使用缓存提高性能

## 🔍 代码示例

### 完整的用例实现

```go
// 应用服务：玩家升级用例
type PlayerLevelUpUseCase struct {
    playerRepo    domain.PlayerRepository
    levelService  domain.LevelCalculationService
    eventBus      infrastructure.EventBus
}

func (uc *PlayerLevelUpUseCase) Execute(cmd GainExperienceCommand) error {
    // 1. 加载聚合
    player, err := uc.playerRepo.FindByID(cmd.PlayerID)
    if err != nil {
        return err
    }
    
    // 2. 执行业务逻辑
    err = player.GainExperience(cmd.Experience)
    if err != nil {
        return err
    }
    
    // 3. 保存聚合
    err = uc.playerRepo.Save(player)
    if err != nil {
        return err
    }
    
    // 4. 发布领域事件
    events := player.GetEvents()
    for _, event := range events {
        err = uc.eventBus.Publish(event)
        if err != nil {
            // 记录日志，但不影响主流程
            log.Error("Failed to publish event", "error", err)
        }
    }
    
    player.ClearEvents()
    return nil
}
```

---

*DDD 版本: v1.0.0 | 最后更新: 2024年*