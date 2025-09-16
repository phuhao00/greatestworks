package battle

import (
	"context"
	"greatestworks/internal/domain/player"
)

// Repository 战斗仓储接口
type Repository interface {
	// Save 保存战斗
	Save(ctx context.Context, battle *Battle) error
	
	// FindByID 根据ID查找战斗
	FindByID(ctx context.Context, id BattleID) (*Battle, error)
	
	// Update 更新战斗
	Update(ctx context.Context, battle *Battle) error
	
	// Delete 删除战斗
	Delete(ctx context.Context, id BattleID) error
	
	// FindByPlayerID 根据玩家ID查找战斗
	FindByPlayerID(ctx context.Context, playerID player.PlayerID, limit int) ([]*Battle, error)
	
	// FindActiveBattles 查找进行中的战斗
	FindActiveBattles(ctx context.Context, limit int) ([]*Battle, error)
	
	// FindByStatus 根据状态查找战斗
	FindByStatus(ctx context.Context, status BattleStatus, limit int) ([]*Battle, error)
	
	// FindByType 根据类型查找战斗
	FindByType(ctx context.Context, battleType BattleType, limit int) ([]*Battle, error)
	
	// CountByPlayerID 统计玩家参与的战斗数量
	CountByPlayerID(ctx context.Context, playerID player.PlayerID) (int64, error)
}