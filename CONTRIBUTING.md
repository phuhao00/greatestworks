# 🤝 贡献指南

感谢您对 Greatest Works 项目的关注！我们欢迎所有形式的贡献，包括但不限于代码、文档、测试、问题报告和功能建议。

## 📋 目录

- [开发环境搭建](#开发环境搭建)
- [代码规范](#代码规范)
- [提交规范](#提交规范)
- [Pull Request 流程](#pull-request-流程)
- [测试要求](#测试要求)
- [文档更新](#文档更新)
- [问题报告](#问题报告)
- [功能建议](#功能建议)

## 🛠️ 开发环境搭建

### 前置要求

- Go 1.21 或更高版本
- Docker 和 Docker Compose
- Git
- MongoDB 5.0+
- Redis 6.0+

### 快速开始

1. **克隆项目**
   ```bash
   git clone https://github.com/your-org/greatestworks.git
   cd greatestworks
   ```

2. **安装依赖**
   ```bash
   go mod download
   ```

3. **设置环境变量**
   ```bash
   cp .env.example .env
   # 编辑 .env 文件，配置数据库连接等信息
   ```

4. **启动开发环境**
   ```bash
   # 使用 Docker Compose 启动依赖服务
   docker-compose up -d mongo redis nats
   
   # 运行数据库迁移
   make migrate
   
   # 启动开发服务器
   make dev
   ```

5. **验证安装**
   ```bash
   curl http://localhost:8080/health
   ```

## 📝 代码规范

### Go 代码规范

我们遵循标准的 Go 代码规范，并使用以下工具确保代码质量：

- `gofmt` - 代码格式化
- `golint` - 代码风格检查
- `go vet` - 静态分析
- `golangci-lint` - 综合代码检查

#### 命名规范

```go
// ✅ 正确的命名
type PlayerService struct {}
func (s *PlayerService) GetPlayerByID(id string) (*Player, error) {}
const MaxPlayersPerRoom = 100
var ErrPlayerNotFound = errors.New("player not found")

// ❌ 错误的命名
type playerservice struct {}
func (s *playerservice) getPlayerById(id string) (*Player, error) {}
const max_players_per_room = 100
var errPlayerNotFound = errors.New("player not found")
```

#### 包结构规范

遵循 DDD（领域驱动设计）架构：

```
internal/
├── application/     # 应用层
│   ├── command/     # 命令处理器
│   ├── query/       # 查询处理器
│   └── service/     # 应用服务
├── domain/          # 领域层
│   ├── player/      # 玩家领域
│   ├── game/        # 游戏领域
│   └── social/      # 社交领域
└── infrastructure/  # 基础设施层
    ├── persistence/ # 数据持久化
    ├── messaging/   # 消息传递
    └── config/      # 配置管理
```

#### 错误处理

```go
// ✅ 正确的错误处理
func (s *PlayerService) CreatePlayer(req *CreatePlayerRequest) (*Player, error) {
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }
    
    player, err := s.repo.Create(req.ToPlayer())
    if err != nil {
        return nil, fmt.Errorf("failed to create player: %w", err)
    }
    
    return player, nil
}

// ❌ 错误的错误处理
func (s *PlayerService) CreatePlayer(req *CreatePlayerRequest) (*Player, error) {
    req.Validate() // 忽略错误
    player, _ := s.repo.Create(req.ToPlayer()) // 忽略错误
    return player, nil
}
```

#### 接口设计

```go
// ✅ 正确的接口设计
type PlayerRepository interface {
    Create(ctx context.Context, player *Player) error
    GetByID(ctx context.Context, id string) (*Player, error)
    Update(ctx context.Context, player *Player) error
    Delete(ctx context.Context, id string) error
}

// ❌ 过于宽泛的接口
type Repository interface {
    Save(interface{}) error
    Load(string) (interface{}, error)
    Delete(string) error
}
```

### 注释规范

```go
// Package player 提供玩家相关的领域模型和业务逻辑
package player

// Player 表示游戏中的玩家实体
type Player struct {
    ID       string    `json:"id" bson:"_id"`
    Username string    `json:"username" bson:"username"`
    Level    int       `json:"level" bson:"level"`
    CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// NewPlayer 创建一个新的玩家实例
// 参数 username 必须是唯一的且长度在 3-20 个字符之间
func NewPlayer(username string) (*Player, error) {
    if len(username) < 3 || len(username) > 20 {
        return nil, ErrInvalidUsername
    }
    
    return &Player{
        ID:        generateID(),
        Username:  username,
        Level:     1,
        CreatedAt: time.Now(),
    }, nil
}
```

## 📤 提交规范

我们使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

### 提交消息格式

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### 提交类型

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式化（不影响代码运行的变动）
- `refactor`: 重构（既不是新增功能，也不是修复 bug 的代码变动）
- `perf`: 性能优化
- `test`: 增加测试
- `chore`: 构建过程或辅助工具的变动
- `ci`: CI/CD 相关变动

### 提交示例

```bash
# 新功能
git commit -m "feat(player): add player level up system"

# 修复 bug
git commit -m "fix(auth): resolve JWT token expiration issue"

# 文档更新
git commit -m "docs: update API documentation for player endpoints"

# 重构
git commit -m "refactor(game): extract battle logic to separate service"

# 破坏性变更
git commit -m "feat(api)!: change player creation endpoint structure

BREAKING CHANGE: player creation now requires email field"
```

## 🔄 Pull Request 流程

### 1. 创建分支

```bash
# 从 main 分支创建新分支
git checkout main
git pull origin main
git checkout -b feature/player-inventory-system
```

### 2. 开发和测试

```bash
# 开发过程中频繁提交
git add .
git commit -m "feat(inventory): add basic inventory structure"

# 运行测试
make test
make lint
```

### 3. 推送分支

```bash
git push origin feature/player-inventory-system
```

### 4. 创建 Pull Request

在 GitHub 上创建 Pull Request，请确保：

- **标题清晰**: 简洁描述变更内容
- **描述详细**: 包含变更原因、实现方式、测试情况
- **关联 Issue**: 如果相关，请关联对应的 Issue
- **截图/演示**: 如果是 UI 变更，请提供截图或演示

#### PR 模板

```markdown
## 📝 变更描述

简要描述这个 PR 的变更内容。

## 🔗 相关 Issue

Closes #123

## 🧪 测试

- [ ] 单元测试通过
- [ ] 集成测试通过
- [ ] 手动测试完成

## 📋 检查清单

- [ ] 代码遵循项目规范
- [ ] 添加了必要的测试
- [ ] 更新了相关文档
- [ ] 没有破坏现有功能

## 📸 截图（如适用）

<!-- 添加截图或 GIF 演示 -->
```

### 5. 代码审查

- 至少需要一个维护者的审查批准
- 解决所有审查意见
- 确保 CI/CD 检查通过

### 6. 合并

- 使用 "Squash and merge" 合并方式
- 删除已合并的分支

## 🧪 测试要求

### 单元测试

每个公共函数都应该有对应的单元测试：

```go
func TestPlayer_LevelUp(t *testing.T) {
    tests := []struct {
        name     string
        player   *Player
        expected int
        wantErr  bool
    }{
        {
            name:     "normal level up",
            player:   &Player{Level: 1},
            expected: 2,
            wantErr:  false,
        },
        {
            name:     "max level reached",
            player:   &Player{Level: 100},
            expected: 100,
            wantErr:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.player.LevelUp()
            if (err != nil) != tt.wantErr {
                t.Errorf("LevelUp() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if tt.player.Level != tt.expected {
                t.Errorf("LevelUp() level = %v, expected %v", tt.player.Level, tt.expected)
            }
        })
    }
}
```

### 集成测试

重要的业务流程需要集成测试：

```go
func TestPlayerService_CreateAndRetrieve(t *testing.T) {
    // 设置测试数据库
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    service := NewPlayerService(db)
    
    // 创建玩家
    req := &CreatePlayerRequest{
        Username: "testplayer",
        Email:    "test@example.com",
    }
    
    player, err := service.CreatePlayer(context.Background(), req)
    require.NoError(t, err)
    require.NotEmpty(t, player.ID)
    
    // 检索玩家
    retrieved, err := service.GetPlayerByID(context.Background(), player.ID)
    require.NoError(t, err)
    assert.Equal(t, player.Username, retrieved.Username)
}
```

### 测试覆盖率

- 新代码的测试覆盖率应该达到 80% 以上
- 核心业务逻辑的覆盖率应该达到 90% 以上

```bash
# 运行测试并生成覆盖率报告
make test-coverage

# 查看覆盖率报告
go tool cover -html=coverage.out
```

## 📚 文档更新

### API 文档

如果变更涉及 API，请更新 OpenAPI 规范：

```yaml
# api/openapi.yaml
paths:
  /api/v1/players:
    post:
      summary: Create a new player
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreatePlayerRequest'
      responses:
        '201':
          description: Player created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Player'
```

### 架构文档

重大架构变更需要更新 `ARCHITECTURE.md`：

- 添加新的架构决策记录 (ADR)
- 更新架构图
- 说明变更原因和影响

### README 更新

如果变更影响项目的使用方式，请更新 `README.md`。

## 🐛 问题报告

### 报告 Bug

使用 GitHub Issues 报告 Bug，请包含：

1. **环境信息**
   - Go 版本
   - 操作系统
   - 数据库版本

2. **重现步骤**
   - 详细的操作步骤
   - 预期结果
   - 实际结果

3. **相关日志**
   - 错误日志
   - 堆栈跟踪

4. **最小重现示例**
   - 如果可能，提供最小的代码示例

### Bug 报告模板

```markdown
## 🐛 Bug 描述

简要描述遇到的问题。

## 🔄 重现步骤

1. 执行 '...'
2. 点击 '....'
3. 滚动到 '....'
4. 看到错误

## 🎯 预期行为

描述你期望发生的情况。

## 📸 截图

如果适用，添加截图来帮助解释你的问题。

## 🖥️ 环境信息

- OS: [e.g. macOS 12.0]
- Go Version: [e.g. 1.21.0]
- Database: [e.g. MongoDB 5.0]

## 📋 附加信息

添加任何其他相关的信息。
```

## 💡 功能建议

### 提出新功能

使用 GitHub Issues 提出功能建议，请包含：

1. **功能描述**: 清晰描述建议的功能
2. **使用场景**: 说明什么情况下会用到这个功能
3. **预期收益**: 这个功能能带来什么价值
4. **实现建议**: 如果有想法，可以提供实现建议

### 功能建议模板

```markdown
## 🚀 功能建议

简要描述你建议的功能。

## 🎯 问题描述

描述当前存在的问题或不便。

## 💡 解决方案

描述你希望看到的解决方案。

## 🔄 替代方案

描述你考虑过的其他替代解决方案。

## 📋 附加信息

添加任何其他相关的信息或截图。
```

## 🏆 贡献者认可

我们感谢每一位贡献者的努力！贡献者将会：

- 在项目 README 中被列出
- 获得项目贡献者徽章
- 参与项目重要决策的讨论

## 📞 联系我们

如果你有任何问题或需要帮助，可以通过以下方式联系我们：

- 创建 GitHub Issue
- 发送邮件到 [maintainers@greatestworks.com](mailto:maintainers@greatestworks.com)
- 加入我们的 Discord 社区

---

再次感谢您对 Greatest Works 项目的贡献！🎉