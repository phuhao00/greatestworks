package player

import "context"

// Repository 玩家仓储接口
type Repository interface {
	// Save 保存玩家
	Save(ctx context.Context, player *Player) error
	
	// FindByID 根据ID查找玩家
	FindByID(ctx context.Context, id PlayerID) (*Player, error)
	
	// FindByName 根据名称查找玩家
	FindByName(ctx context.Context, name string) (*Player, error)
	
	// Update 更新玩家
	Update(ctx context.Context, player *Player) error
	
	// Delete 删除玩家
	Delete(ctx context.Context, id PlayerID) error
	
	// FindOnlinePlayers 查找在线玩家
	FindOnlinePlayers(ctx context.Context, limit int) ([]*Player, error)
	
	// FindPlayersByLevel 根据等级范围查找玩家
	FindPlayersByLevel(ctx context.Context, minLevel, maxLevel int) ([]*Player, error)
	
	// ExistsByName 检查名称是否存在
	ExistsByName(ctx context.Context, name string) (