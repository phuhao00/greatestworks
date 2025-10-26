package character

import (
	"context"
)

// EntityRepository 实体仓储接口 - 由基础设施层实现
type EntityRepository interface {
	// 通用实体操作
	Register(ctx context.Context, entity *Entity) error
	Unregister(ctx context.Context, entityID EntityID) error
	Get(ctx context.Context, entityID EntityID) (*Entity, error)
	GetAll(ctx context.Context) ([]*Entity, error)

	// 按类型查询
	GetByType(ctx context.Context, entityType EntityType) ([]*Entity, error)
}

// PlayerRepository 玩家仓储接口
type PlayerRepository interface {
	// 保存玩家数据到数据库
	Save(ctx context.Context, player *Player) error

	// 从数据库加载玩家
	Load(ctx context.Context, characterID int64) (*Player, error)

	// 删除玩家
	Delete(ctx context.Context, characterID int64) error

	// 查询用户的所有角色
	FindByUserID(ctx context.Context, userID int64) ([]*Player, error)

	// 检查角色名是否存在
	ExistsByName(ctx context.Context, name string) (bool, error)
}

// MonsterRepository 怪物仓储接口
type MonsterRepository interface {
	// 创建怪物实例
	Create(ctx context.Context, monster *Monster) error

	// 销毁怪物实例
	Destroy(ctx context.Context, entityID EntityID) error

	// 根据刷新点ID获取怪物列表
	GetBySpawnID(ctx context.Context, spawnID int32) ([]*Monster, error)
}

// NPCRepository NPC仓储接口
type NPCRepository interface {
	// 创建NPC实例
	Create(ctx context.Context, npc *NPC) error

	// 销毁NPC实例
	Destroy(ctx context.Context, entityID EntityID) error

	// 根据地图ID获取NPC列表
	GetByMapID(ctx context.Context, mapID int32) ([]*NPC, error)
}

// UnitDefineRepository 单位定义仓储接口（配置数据）
type UnitDefineRepository interface {
	// 获取单位定义
	Get(ctx context.Context, unitID int32) (*UnitDefine, error)

	// 加载所有单位定义
	LoadAll(ctx context.Context) (map[int32]*UnitDefine, error)
}

// UnitDefine 单位定义（从配置文件加载）
type UnitDefine struct {
	ID   int32  // 单位ID
	Name string // 单位名称
	Type string // 单位类型

	// 基础属性
	BaseHP float32 // 基础生命值
	BaseMP float32 // 基础魔法值
	BaseAD float32 // 基础物理攻击
	BaseAP float32 // 基础法术攻击

	// 其他属性
	Speed              float32 // 移动速度
	ViewRange          float32 // 视野范围
	HurtTime           float32 // 受击硬直时间
	DropExpBase        float32 // 基础经验掉落
	DropExpLevelFactor float32 // 经验等级系数

	// 技能列表
	Skills []int32 // 技能ID列表
}
