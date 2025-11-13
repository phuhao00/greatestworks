---
applyTo: '**'
---
# GreatestWorks 项目 AI 协作说明（必读）

本文件为 AI 与协作者在本仓库中进行代码生成、变更评审与文档撰写时的统一约束与上下文指引。请严格遵循以下规范，确保改动与项目架构、风格和运行形态保持一致。

## 工作流程补充要求

每次开始编写新功能或进行较大改动时，需先列出本次开发的 todo list（任务清单），明确拆分为具体、可执行的小步骤，并在开发过程中逐项完成和勾选。todo list 应包含：

- 主要功能点拆解（如领域模型、应用服务、接口适配、配置变更等）
- 相关协议/配置/文档同步
- 单元测试与集成验证
- 变更影响面评估与回滚方案

todo list 可在 PR 描述、评论或协作工具中展示，便于团队成员跟进进度与评审。

## 项目上下文速览

- 语言与版本：Go 1.24（见 `go.mod`）
- 架构风格：DDD（领域驱动）+ 分层架构 + 微服务
- 服务与入口：
	- 认证服务：`cmd/auth-service/main.go`（HTTP: 8080）
	- 网关服务：`cmd/gateway-service/main.go`（TCP: 9090）
	- 游戏服务：`cmd/game-service/main.go`（Go 原生 RPC: 8081）
	- 其他服务：`cmd/scene/main.go`、`cmd/replication/main.go`（保留/扩展）
- 配置位置：`configs/*.yaml`（支持 dev/prod/docker 等环境示例）
- 数据存储：MongoDB（主存储）+ Redis（缓存）
- 协议与通信：
	- 客户端 → 认证：HTTP
	- 客户端 → 网关：TCP（自定义二进制协议，参考 `internal/network` 与 `proto/`）
	- 网关 ↔ 游戏：Go 原生 RPC
- 协议定义：根目录 `proto/*.proto`（以及 `internal/proto/*` 目录结构），存在独立 proto 模块替换：`replace github.com/phuhao00/greatestworks-proto => ../greatestworks-proto`
- 压测/集成测试工具：`tools/simclient`（支持 E2E/压测/功能验证场景）
- 容器与脚本：`Dockerfile`、`docker-compose.yml`、`scripts/*`、`Makefile`

## 代码组织与分层约定（DDD）

- 领域层：`internal/domain/*`
	- 聚合根/实体/值对象仅暴露必要行为，字段使用小写私有化，保持不变式。
	- 领域服务仅聚焦业务规则，不做基础设施细节。
	- 领域事件放在 `internal/domain/events`。
- 应用层：`internal/application/*`
	- 命令/查询处理器与应用服务；通过 `services.ServiceRegistry` 进行装配。
	- 面向用例编排，不直接依赖具体基础设施实现，依赖接口。
- 基础设施层：`internal/infrastructure/*`
	- 持久化、缓存、消息、日志、配置、网络等适配与实现。
	- 仓储实现放于 `infrastructure/persistence/*`，接口在领域层定义。
- 接口层：`internal/interfaces/{http,tcp,rpc}/*`
	- 将传输协议与应用层解耦；在此完成 DTO/协议消息 ↔ 领域模型映射。

新增功能时遵循“自内向外”原则：先领域模型与用例，再接口适配与基础设施实现，最后编排到入口服务。

## 依赖与模块策略

- 优先使用标准库与现有依赖，不轻易新增第三方库；如需新增必须：
	1) 说明动机与替代方案；2) 版本固定并最小化；3) 兼容 Go 1.24。
- 修改依赖后执行 `go mod tidy` 保持清洁；不要引入需要 CGO 的库（跨平台/容器化负担）。
- 保留并尊重 `go.mod` 中的 `replace github.com/phuhao00/greatestworks-proto`，该仓库用于统一管理协议代码；不要在本仓库内随意复制粘贴生成代码覆盖。

## 协议与兼容性（Proto/TCP）

- 修改 `.proto` 前务必确认兼容性策略：
	- 仅追加字段，避免复用/更改现有 tag；删除请改为 `reserved`。
	- 变更需同步网关与游戏服务的处理逻辑，保持向后兼容。
- 生成代码：优先使用 `scripts/generate_proto.bat`（Windows）或 `scripts/generate_proto.sh`（Unix）；不要手工在 `internal/proto` 下散落生成物。
- TCP 二进制协议格式固定（见 README）；非兼容性变更需评审并更新 `internal/network` 与模拟客户端映射。

## 配置管理

- 所有服务读取 `configs/*.yaml`，经 `internal/config` 管理加载；禁止硬编码端口/密钥/连接串。
- 支持环境变量覆盖（在文档/示例配置中体现）；新增配置项需：
	- 更新对应 `configs/*.yaml` 示例；
	- 更新读取/校验逻辑与默认值；
	- 在 README 或服务说明中记录。

## 日志与监控

- 使用 `internal/infrastructure/logging` 提供的 `Logger`，禁止使用 `fmt.Println`。
- 日志规范：结构化（json），带关键字段：`service`, `module`, `player_id`, `trace_id`（如可用）。
- 级别：`debug`（开发调试）、`info`（业务里程碑）、`warn`（可恢复异常）、`error`（失败/降级）、`fatal`（进程退出）。
- 性能剖析：通过配置开启 pprof（参考 README 的 monitoring.profiling）；仅在可信网络或受控环境开放。

## 错误处理

- 领域错误集中于 `internal/errors`（如 `domain_errors.go`）；服务/接口层使用 `errors.Is/As` 判断语义。
- Wrap 原始错误保留栈上下文：`fmt.Errorf("…: %w", err)`。
- HTTP 映射：领域校验错误 → 400，权限问题 → 401/403，找不到 → 404，服务内部错误 → 500。
- TCP/RPC：映射到协议内的错误码或错误消息，保持一致性并记录日志。

## 并发与资源管理

- 所有跨边界函数均传递 `context.Context`，尊重取消与超时；网络/数据库调用设置合理超时。
- 禁止无限制 goroutine 派生；使用工作池（见 `infrastructure/messaging/worker_pool.go`）或事件驱动（`internal/events`）。
- 连接/句柄必须在 `defer` 中关闭；大对象复用需证明收益且不影响可读性。

## 测试与验证

- 单元测试：同包 `_test.go`，采用表驱动，覆盖核心分支；提交前保证 `go test ./...` 通过。
- 集成/E2E：使用 `tools/simclient`；支持 smoke/feature/load 模式：
	- smoke：在服务运行下进行冒烟；可通过环境变量控制跳过。
	- feature：按功能库驱动场景（参考 `tools/simclient/feature_library.go`）。
	- load：并发用户压测，输出 P50/P95/Max 时延统计。
- 新增协议或接口时，应提供最小可运行的集成验证脚本/说明。

## 代码风格与提交

- 风格：`gofmt`/`goimports` 必须；命名符合 Go 惯例（接口以 `er` 结尾、包名简短小写）。
- 变更最小化：倾向小步提交，避免无关重构；保留公开 API 稳定性。
- 提交信息（建议使用 Conventional Commits）：
	- feat: 新功能；fix: 修复；refactor: 重构；docs: 文档；test: 测试；chore: 其它；build/ci: 构建与流水线。
- PR 要求：
	- 说明动机、设计要点、影响面（协议/配置/数据/性能）；
	- 附上验证方式（单测/模拟客户端步骤）；
	- 涉及配置或脚本的，更新对应示例与 README 片段。

## 变更落地清单（AI 执行时务必遵循）

1) 读取相关代码与文档，明确所处层与依赖关系；
2) 若涉及协议/配置/公共接口，先出兼容性策略与影响面；
3) 优先修改领域/应用层契约，再补充基础设施与接口适配；
4) 如需新依赖，给出选择理由与安全性评估，版本固定；
5) 同步更新：示例配置、脚本、文档与模拟客户端用例；
6) 保持改动最小并可回滚；
7) 本地自检：静态检查 + 单测通过（若仓库上下文缺依赖导致无法编译，请在说明中标注并提供可替代验证方式）。

## 不要做（DON’Ts）

- 不要绕过 `internal/infrastructure/logging` 打印日志。
- 不要在接口层直接操作数据库或 Redis；通过应用/领域契约间接访问。
- 不要破坏 `.proto` 的兼容性（更改已有 tag、复用 tag、删除非 reserved）。
- 不要硬编码端口、密钥、连接字符串。
- 不要引入重量级依赖或需要本地原生编译的库（除非有充分理由与替代方案评审）。

## 附：常用位置速查

- 入口：`cmd/*/main.go`
- 配置：`configs/*.yaml`
- 领域：`internal/domain/*`
- 应用：`internal/application/*`（`services/service_registry.go` 装配）
- 基础设施：`internal/infrastructure/*`
- 接口适配：`internal/interfaces/{http,tcp,rpc}`
- 协议：`proto/*.proto`（生成脚本见 `scripts/generate_proto.*`）
- 网络协议实现：`internal/network/*`
- 模拟客户端：`tools/simclient/*`

如对以上规范与上下文有疑问，请在 PR 描述中明确提出并给出权衡与建议方案，以便评审达成一致。