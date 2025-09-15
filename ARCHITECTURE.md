# 项目架构文档

## 概述

本项目已按照领域驱动设计(DDD)架构进行重构，将原有的 `config`、`aop`、`protocol` 目录整合到统一的基础设施层中。

## 目录结构

```
greatestworks/
├── cmd/                           # 应用程序入口
│   └── server/                    # 游戏服务器
│       ├── main.go               # 主入口文件
│       └── bootstrap.go          # 启动引导系统
├── application/                   # 应用层
│   └── services/                 # 应用服务
├── internal/                     # 内部模块
│   └── infrastructure/           # 基础设施层
│       ├── config/              # 配置管理
│       ├── logging/             # 日志系统
│       ├── monitoring/          # 监控系统
│       ├── protocol/            # 协议管理
│       ├── weave/               # Service Weaver集成
│       ├── container/           # 依赖注入容器
│       ├── network/             # 网络层
│       └── persistence/         # 持久化层
└── interfaces/                   # 接口层
    └── tcp/                     # TCP接口
```

## 重构内容

### 1. 配置管理 (Configuration Management)

**原位置**: `config/`, `aop/config/`  
**新位置**: `internal/infrastructure/config/`

#### 功能特性:
- 统一配置加载器 (`loader.go`)
- 环境管理器 (`env_manager.go`)
- 配置验证器 (`validation.go`)
- 热重载支持 (`hot_reload.go`)
- 多环境配置文件 (`environments/`)

#### 配置文件结构:
```
internal/infrastructure/config/
├── config.go              # 配置结构定义
├── loader.go              # 配置加载器
├── env_manager.go         # 环境管理
├── validation.go          # 配置验证
├── hot_reload.go          # 热重载
└── environments/          # 环境配置
    ├── config.dev.yaml   # 开发环境
    ├── config.prod.yaml  # 生产环境
    └── config.test.yaml  # 测试环境
```

### 2. 日志系统 (Logging System)

**原位置**: `aop/logger/`  
**新位置**: `internal/infrastructure/logging/`

#### 功能特性:
- 多种日志实现 (Zap, File, Console)
- 结构化日志支持
- 日志格式化器
- HTTP/Game中间件
- 日志级别管理

#### 文件结构:
```
internal/infrastructure/logging/
├── logger.go             # 日志接口定义
├── zap_logger.go         # Zap日志实现
├── file_logger.go        # 文件日志
├── console_logger.go     # 控制台日志
├── formatter.go          # 格式化器
└── middleware.go         # 日志中间件
```

### 3. 监控系统 (Monitoring System)

**原位置**: `aop/metrics/`  
**新位置**: `internal/infrastructure/monitoring/`

#### 功能特性:
- Prometheus集成
- 多种指标类型 (Counter, Gauge, Histogram, Summary)
- 系统指标收集器
- 游戏指标收集器
- 指标导出和服务器

#### 文件结构:
```
internal/infrastructure/monitoring/
├── metrics.go            # 指标接口定义
├── prometheus.go         # Prometheus实现
└── collectors.go         # 指标收集器
```

### 4. 协议管理 (Protocol Management)

**原位置**: `protocol/`  
**新位置**: `internal/infrastructure/protocol/`

#### 功能特性:
- 二进制协议支持
- JSON协议支持
- 协议编解码器
- 消息验证
- 协议管理器

#### 文件结构:
```
internal/infrastructure/protocol/
├── protocol.go           # 协议接口定义
├── binary_protocol.go    # 二进制协议
└── json_protocol.go      # JSON协议
```

### 5. Service Weaver集成 (Weave Integration)

**原位置**: `aop/weavelet.go`  
**新位置**: `internal/infrastructure/weave/`

#### 功能特性:
- Weavelet管理器
- 配置验证
- 标签管理
- 生命周期管理

#### 文件结构:
```
internal/infrastructure/weave/
└── weavelet.go           # Weave集成
```

### 6. 依赖注入容器 (DI Container)

**新增功能**  
**位置**: `internal/infrastructure/container/`

#### 功能特性:
- 服务注册和解析
- 生命周期管理 (Singleton, Transient, Scoped)
- 依赖注入
- 服务提供者模式
- 作用域管理

#### 文件结构:
```
internal/infrastructure/container/
├── container.go          # DI容器实现
└── providers.go          # 服务提供者
```

### 7. 启动引导系统 (Bootstrap System)

**原位置**: `aop/bootstrap.go`  
**新位置**: `cmd/server/bootstrap.go`

#### 功能特性:
- 统一启动流程
- 服务初始化
- 优雅关闭
- 错误处理
- Service Weaver支持

## 使用方式

### 1. 启动服务器

```go
// cmd/server/main.go
func main() {
    // 创建服务器启动器
    bootstrap := NewServerBootstrap()
    
    // 初始化服务器
    if err := bootstrap.Initialize(); err != nil {
        log.Fatalf("服务器初始化失败: %v", err)
    }
    
    // 启动服务器
    if err := bootstrap.StartServer(); err != nil {
        log.Fatalf("服务器启动失败: %v", err)
    }
}
```

### 2. 配置管理

```go
// 加载配置
loader := config.NewConfigLoader()
cfg, err := loader.LoadFromFile("config.yaml")

// 环境管理
envManager := config.GetEnvManager()
configPath := envManager.GetConfigPath()

// 热重载
hotReloader := config.NewHotReloader(&config.HotReloadConfig{
    Enabled: true,
    WatchPaths: []string{"config.yaml"},
})
```

### 3. 依赖注入

```go
// 创建容器
container := container.NewContainer()

// 注册服务提供者
provider := container.NewAllProvidersProvider("config.yaml")
container.RegisterProvider(provider)

// 解析服务
logger, err := container.Resolve("logger")
config, err := container.Resolve("config")
```

### 4. 日志记录

```go
// 创建日志器
logger, err := logging.NewZapLogger(&logConfig)

// 记录日志
logger.Info("服务器启动")
logger.WithField("port", 8080).Info("监听端口")
logger.WithError(err).Error("启动失败")
```

### 5. 监控指标

```go
// 创建指标注册表
registry, err := monitoring.NewPrometheusRegistry(&metricsConfig)

// 创建指标
counter := registry.NewCounter("requests_total", "Total requests")
gauge := registry.NewGauge("active_connections", "Active connections")

// 记录指标
counter.Inc()
gauge.Set(100)
```

## 配置文件示例

### 开发环境配置 (config.dev.yaml)

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

database:
  mongodb:
    uri: "mongodb://localhost:27017"
    database: "greatestworks_dev"
    timeout: 10s

logging:
  type: "zap"
  level: "debug"
  format: "json"
  output: ["stdout", "file"]
  file:
    path: "logs/app.log"
    max_size: 100
    max_backups: 3
    max_age: 7

monitoring:
  enabled: true
  port: 9090
  path: "/metrics"
```

## 迁移指南

### 更新导入路径

原有的导入路径需要更新为新的路径:

```go
// 旧的导入
import "greatestworks/config"
import "greatestworks/aop/logger"
import "greatestworks/aop/metrics"
import "greatestworks/protocol"

// 新的导入
import "greatestworks/internal/infrastructure/config"
import "greatestworks/internal/infrastructure/logging"
import "greatestworks/internal/infrastructure/monitoring"
import "greatestworks/internal/infrastructure/protocol"
```

### 配置文件迁移

1. 将原有的配置文件移动到 `internal/infrastructure/config/environments/`
2. 统一配置文件格式为 YAML
3. 更新配置结构以匹配新的配置定义

### 代码重构

1. 使用新的启动引导系统
2. 采用依赖注入容器管理服务
3. 更新日志和监控调用
4. 使用新的协议管理器

## 优势

1. **统一管理**: 所有基础设施组件集中管理
2. **依赖注入**: 松耦合的服务架构
3. **配置标准化**: 统一的配置管理和验证
4. **可扩展性**: 易于添加新的基础设施组件
5. **可测试性**: 更好的单元测试支持
6. **维护性**: 清晰的代码组织和职责分离

## 注意事项

1. 确保所有导入路径已更新
2. 配置文件格式需要统一为 YAML
3. 使用依赖注入容器管理服务生命周期
4. 遵循 DDD 架构原则进行后续开发