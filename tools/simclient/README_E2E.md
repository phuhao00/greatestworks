# SimClient E2E 测试使用说明

## 概述

`tools/simclient` 提供了端到端（E2E）测试能力，用于验证网关服务的完整功能链路：
- **登录认证**（可选）：通过 HTTP 认证服务获取 token
- **TCP 连接**：连接到网关 TCP 服务器
- **游戏协议**：发送登录、移动、技能释放、登出等消息
- **AOI 验证**：观察服务器响应与广播

## 快速开始

### 1. 准备服务器

确保以下服务正在运行：
```powershell
# 启动网关服务（默认端口 9090）
go run ./cmd/gateway-service

# （可选）启动认证服务（默认端口 8080）
go run ./cmd/auth-service
```

### 2. 运行单次 E2E 测试

```powershell
# 使用提供的 E2E 配置
go run ./tools/simclient/cmd/simclient -config=tools/simclient/e2e.yaml -mode=integration

# 或使用命令行参数覆盖
go run ./tools/simclient/cmd/simclient -config=tools/simclient/e2e.yaml -mode=integration -debug
```

**预期输出示例**：
```
Scenario: e2e-login-move-skill
Duration: 1.234s
Success: true
Actions:
  gateway.connect            123ms  OK
  gateway.msg.login          45ms   OK
  login.response             12ms   OK
  gateway.msg.move           8ms    OK
  move.response              10ms   OK
  gateway.msg.skill          9ms    OK
  skill.response             11ms   OK
  gateway.msg.move           7ms    OK
  move2.response             10ms   OK
  gateway.msg.logout         6ms    OK
  logout.response            5ms    OK
```

### 3. 运行压力测试

```powershell
# 50 个虚拟玩家，并发 10，每个玩家执行 3 次完整流程
go run ./tools/simclient/cmd/simclient -config=tools/simclient/e2e_load.yaml -mode=load

# 命令行自定义参数
go run ./tools/simclient/cmd/simclient ^
  -config=tools/simclient/e2e_load.yaml ^
  -mode=load ^
  -users=100 ^
  -concurrency=20 ^
  -iterations=5
```

**预期输出示例**：
```
Load Scenario: e2e-load-test
Users: 50  Concurrency: 10  Iterations/User: 3
Total Duration: 12.456s
Scenarios: 150 (success: 148, failures: 2)
Action Metrics:
  gateway.connect          count= 150 success= 150 fail=  0  min=   45ms avg=  123ms p95=  234ms max=  456ms
  gateway.msg.login        count= 150 success= 150 fail=  0  min=    5ms avg=   12ms p95=   23ms max=   45ms
  gateway.msg.move         count= 300 success= 300 fail=  0  min=    3ms avg=    8ms p95=   15ms max=   34ms
  gateway.msg.skill        count= 150 success= 150 fail=  0  min=    4ms avg=    9ms p95=   18ms max=   28ms
  ...
```

## E2E 场景流程

新增的 `E2EScenario` (`tools/simclient/e2e_scenario.go`) 执行以下步骤：

1. **认证**（如果 `auth.enabled=true`）
   - 向 HTTP 认证服务发送登录请求
   - 获取并记录 token（当前未附加到后续消息，可扩展）

2. **连接网关**
   - TCP 连接到 `gateway.host:gateway.port`
   - 设置读写超时

3. **发送登录包**
   ```json
   {
     "player_id": "123456",
     "map_id": 1
   }
   ```
   - 消息类型：`MsgPlayerLogin`

4. **发送移动包**
   ```json
   {
     "position": {"x": 100.0, "y": 50.0, "z": 10.0}
   }
   ```
   - 消息类型：`MsgPlayerMove`

5. **发送技能释放包**
   ```json
   {
     "skill_id": 1001,
     "target_id": 2001
   }
   ```
   - 消息类型：`MsgBattleSkill`

6. **再次移动**
   - 验证多次操作

7. **发送登出包**
   - 消息类型：`MsgPlayerLogout`

每个步骤后会尝试读取服务器响应（200ms 超时），记录接收到的字节数或超时情况。

## 配置说明

### scenario 配置
- `type: e2e` - 使用 E2E 场景（必须）
- `name` - 场景名称（用于报告）
- `player_prefix` - 玩家名前缀（用于生成唯一 player_id）
- `action_interval` - 动作间隔（E2E 场景自动控制，可忽略）
- `stop_on_error` - 遇到错误是否立即停止

### auth 配置
- `enabled` - 是否启用 HTTP 认证
- `base_url` - 认证服务地址（如 `http://localhost:8080`）
- `login_path` - 登录路径（如 `/api/v1/auth/login`）
- `username/password` - 认证凭据

### gateway 配置
- `host` - 网关 TCP 服务器地址
- `port` - 网关 TCP 端口（默认 9090）
- `connect_timeout` - 连接超时
- `read_timeout/write_timeout` - 读写超时

### load 配置（压测模式）
- `enabled: true` - 启用压测
- `virtual_users` - 虚拟玩家总数
- `concurrency` - 并发执行数
- `iterations` - 每个玩家执行完整流程的次数
- `ramp_up` - 启动所有玩家的时间（平滑启动）
- `stop_on_error` - 遇到错误是否停止整个压测

## 高级用法

### 自定义场景

如需修改流程，编辑 `tools/simclient/e2e_scenario.go`：

```go
// 修改技能 ID 或目标 ID
func (s *E2EScenario) Execute(ctx context.Context, client *SimulatorClient) (*ScenarioResult, error) {
    // ... 现有代码 ...
    
    // 自定义：释放多个技能
    s.sendSkillCast(result, client, conn, 1001, 2001)
    s.sendSkillCast(result, client, conn, 1002, 2001)
    
    // ... 继续 ...
}
```

### 验证 AOI 广播

当多个 simclient 同时运行时，可观察是否收到来自其他玩家的广播消息：

```powershell
# 终端 1
go run ./tools/simclient/cmd/simclient -config=tools/simclient/e2e.yaml -mode=integration -debug

# 终端 2（稍后启动）
go run ./tools/simclient/cmd/simclient -config=tools/simclient/e2e.yaml -mode=integration -debug
```

在 debug 日志中查看是否收到 `entity_move`、`entity_appear` 等 AOI 消息。

### 集成到 CI/CD

```yaml
# .github/workflows/test.yml
- name: Run E2E Test
  run: |
    go run ./cmd/gateway-service &
    sleep 5
    go run ./tools/simclient/cmd/simclient -config=tools/simclient/e2e.yaml -mode=integration
    killall gateway-service
```

## 故障排查

### 连接失败
- **错误**：`dial gateway localhost:9090: connection refused`
- **解决**：确保网关服务已启动并监听正确端口

### 认证失败
- **错误**：`auth request returned 401`
- **解决**：检查 `auth.username` 和 `password` 是否正确；或设置 `auth.enabled: false` 跳过认证

### 超时
- **错误**：大量 `timeout: true` 记录
- **原因**：服务器未及时响应
- **解决**：
  - 增加 `gateway.read_timeout`
  - 检查服务器日志是否有错误
  - 确认服务器正常处理消息

### 编译错误
```powershell
# 重新编译 simclient
cd c:\Users\HHaou\greatestworks
go build -o simclient.exe ./tools/simclient/cmd/simclient
./simclient.exe -config=tools/simclient/e2e.yaml
```

## 扩展与贡献

欢迎扩展 E2E 场景，如：
- 添加队伍组队测试
- 添加聊天消息测试
- 添加背包/交易测试
- 验证 AOI 广播正确性（多玩家视野同步）

修改 `tools/simclient/e2e_scenario.go` 并提交 PR。

---

**相关文件**：
- `tools/simclient/e2e_scenario.go` - E2E 场景实现
- `tools/simclient/e2e.yaml` - 单次测试配置
- `tools/simclient/e2e_load.yaml` - 压测配置
- `tools/simclient/cmd/simclient/main.go` - CLI 入口
