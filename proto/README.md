# Proto文件说明

本目录包含GreatestWorks项目的Protocol Buffers定义文件，支持Go和C#客户端。

## 文件结构

```
proto/
├── common.proto      # 通用消息定义
├── player.proto      # 玩家相关服务
├── battle.proto      # 战斗相关服务
├── pet.proto         # 宠物相关服务
└── README.md         # 本文件
```

## 生成代码

### 方法1: 使用Makefile (推荐)

```bash
# 生成所有Proto文件
make proto-all

# 只生成Go代码
make proto-go

# 只生成C#代码
make proto-csharp
```

### 方法2: 使用脚本

```bash
# Linux/macOS
./scripts/generate_proto.sh

# Windows
scripts\generate_proto.bat
```

### 方法3: 手动生成

```bash
# 生成Go代码
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/*.proto

# 生成C#代码
protoc --csharp_out=csharp --csharp_opt=file_extension=.g.cs \
       proto/*.proto
```

## 输出目录

### Go代码
- `internal/proto/player/` - 玩家服务
- `internal/proto/battle/` - 战斗服务
- `internal/proto/pet/` - 宠物服务
- `internal/proto/common/` - 通用消息

### C#代码
- `csharp/GreatestWorks.Player/` - 玩家服务
- `csharp/GreatestWorks.Battle/` - 战斗服务
- `csharp/GreatestWorks.Pet/` - 宠物服务
- `csharp/GreatestWorks.Common/` - 通用消息

## 服务定义

### PlayerService (玩家服务)
- `CreatePlayer` - 创建玩家
- `Login` - 玩家登录
- `Logout` - 玩家登出
- `GetPlayerInfo` - 获取玩家信息
- `UpdatePlayer` - 更新玩家信息
- `MovePlayer` - 移动玩家
- `GetOnlinePlayers` - 获取在线玩家列表

### BattleService (战斗服务)
- `CreateBattle` - 创建战斗
- `JoinBattle` - 加入战斗
- `LeaveBattle` - 离开战斗
- `ExecuteAction` - 执行战斗动作
- `GetBattleInfo` - 获取战斗信息
- `GetBattleList` - 获取战斗列表

### PetService (宠物服务)
- `CreatePet` - 创建宠物
- `GetPetInfo` - 获取宠物信息
- `UpdatePet` - 更新宠物信息
- `LevelUpPet` - 宠物升级
- `EvolvePet` - 宠物进化
- `GetPlayerPets` - 获取玩家宠物列表

## 使用示例

### Go客户端示例

```go
package main

import (
    "context"
    "log"
    
    pb "greatestworks/internal/proto/player"
    "google.golang.org/grpc"
)

func main() {
    conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    client := pb.NewPlayerServiceClient(conn)
    
    // 创建玩家
    resp, err := client.CreatePlayer(context.Background(), &pb.CreatePlayerRequest{
        Name: "测试玩家",
        AccountId: "account123",
        InitialLevel: 1,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("创建玩家成功: %s", resp.Player.Name)
}
```

### C#客户端示例

```csharp
using GreatestWorks.Player;
using Grpc.Core;

class Program
{
    static void Main(string[] args)
    {
        var channel = new Channel("localhost:8080", ChannelCredentials.Insecure);
        var client = new PlayerService.PlayerServiceClient(channel);
        
        // 创建玩家
        var request = new CreatePlayerRequest
        {
            Name = "测试玩家",
            AccountId = "account123",
            InitialLevel = 1
        };
        
        var response = client.CreatePlayer(request);
        Console.WriteLine($"创建玩家成功: {response.Player.Name}");
        
        channel.ShutdownAsync().Wait();
    }
}
```

## 注意事项

1. **不使用gRPC**: 本项目使用netcore-go RPC，不是标准的gRPC
2. **编码格式**: 所有字符串使用UTF-8编码
3. **版本兼容**: 确保protoc版本 >= 3.0
4. **依赖管理**: 生成后需要运行 `go mod tidy` 更新依赖

## 故障排除

### 常见问题

1. **protoc未找到**
   ```bash
   # 安装protoc
   # Windows: 下载 https://github.com/protocolbuffers/protobuf/releases
   # macOS: brew install protobuf
   # Ubuntu: sudo apt-get install protobuf-compiler
   ```

2. **Go插件未找到**
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

3. **C#插件未找到**
   - 确保安装了完整的Protocol Buffers工具链
   - 检查PATH环境变量

4. **权限问题**
   ```bash
   # Linux/macOS
   chmod +x scripts/generate_proto.sh
   ```

## 更新日志

- v1.0.0: 初始版本，支持玩家、战斗、宠物服务
- 支持Go和C#客户端
- 使用netcore-go RPC架构
