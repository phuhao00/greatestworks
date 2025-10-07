# Proto 常量管理

本项目将所有交互相关的枚举、常量和协议都定义在proto文件中，通过代码生成来统一管理。

## 文件结构

```
internal/proto/
├── battle/          # 战斗相关proto
├── common/          # 通用proto
├── errors/          # 错误码proto
├── messages/        # 消息号proto
├── pet/            # 宠物相关proto
├── player/         # 玩家相关proto
└── protocol/       # 协议常量proto
```

## Proto文件说明

### 1. errors.proto - 错误码管理
包含所有错误码和错误消息：
- `CommonErrorCode`: 通用错误 (1000-1999)
- `BattleErrorCode`: 战斗相关错误 (2000-2999)
- `PetErrorCode`: 宠物相关错误 (3000-3999)
- `ItemErrorCode`: 物品相关错误 (4000-4999)
- `BuildingErrorCode`: 建筑相关错误 (5000-5999)
- `SocialErrorCode`: 社交相关错误 (6000-6999)
- `QuestErrorCode`: 任务相关错误 (7000-7999)
- `SystemErrorCode`: 系统相关错误 (8000-8999)

### 2. messages.proto - 消息号管理
包含所有消息类型常量：
- `SystemMessageID`: 系统消息 (0x0000-0x00FF)
- `PlayerMessageID`: 玩家相关消息 (0x0100-0x01FF)
- `BattleMessageID`: 战斗相关消息 (0x0200-0x02FF)
- `PetMessageID`: 宠物相关消息 (0x0300-0x03FF)
- `BuildingMessageID`: 建筑相关消息 (0x0400-0x04FF)
- `SocialMessageID`: 社交相关消息 (0x0500-0x05FF)
- `ItemMessageID`: 物品相关消息 (0x0600-0x06FF)
- `QuestMessageID`: 任务相关消息 (0x0700-0x07FF)
- `QueryMessageID`: 查询相关消息 (0x0800-0x08FF)
- `AdminMessageID`: 系统管理消息 (0x0900-0x09FF)

### 3. protocol.proto - 协议常量
包含协议相关的枚举和常量：
- 各种消息类型枚举
- 错误码枚举
- 状态枚举
- 消息标志位枚举
- 协议常量定义

### 4. 其他模块proto文件
- `battle.proto`: 战斗相关枚举和消息
- `player.proto`: 玩家相关枚举和消息
- `pet.proto`: 宠物相关枚举和消息
- `common.proto`: 通用枚举和消息

## 使用方法

### 导入生成的包
```go
import (
    "greatestworks/internal/proto/errors"
    "greatestworks/internal/proto/messages"
    "greatestworks/internal/proto/protocol"
    "greatestworks/internal/proto/common"
    "greatestworks/internal/proto/battle"
    "greatestworks/internal/proto/player"
    "greatestworks/internal/proto/pet"
)
```

### 使用错误码
```go
// 使用通用错误码
if err != nil {
    return errors.CommonErrorCode_ERR_PLAYER_NOT_FOUND
}

// 使用战斗错误码
if battleFull {
    return errors.BattleErrorCode_ERR_BATTLE_FULL
}
```

### 使用消息号
```go
// 创建消息头
header := &messages.MessageHeader{
    MessageId:   uint32(messages.PlayerMessageID_MSG_PLAYER_LOGIN),
    MessageType: uint32(messages.SystemMessageID_MSG_AUTH),
    Flags:       uint32(messages.MessageFlag_MESSAGE_FLAG_REQUEST),
}
```

### 使用状态枚举
```go
// 检查玩家状态
if player.Status == protocol.PlayerStatus_PLAYER_STATUS_ONLINE {
    // 处理在线玩家
}

// 检查战斗状态
if battle.Status == protocol.BattleStatus_BATTLE_STATUS_ACTIVE {
    // 处理进行中的战斗
}
```

## 代码生成

使用以下命令重新生成proto代码：

```bash
# 生成所有proto文件
.\protoc\bin\protoc.exe --go_out=. --go_opt=paths=source_relative proto/*.proto

# 或者使用脚本
scripts/generate_proto.sh
```

## 优势

1. **类型安全**: 所有常量都有明确的类型定义
2. **统一管理**: 所有枚举和常量都在proto文件中定义
3. **代码生成**: 自动生成Go代码，避免手动维护
4. **版本控制**: proto文件可以版本化，便于管理
5. **多语言支持**: 可以生成多种语言的代码
6. **文档化**: proto文件本身就是很好的文档

## 注意事项

1. 修改proto文件后需要重新生成代码
2. 枚举值不能重复，需要合理分配范围
3. 添加新的枚举值时需要考虑向后兼容性
4. 错误码和消息号需要按模块分类管理
