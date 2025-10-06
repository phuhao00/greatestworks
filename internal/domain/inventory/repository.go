package inventory

import (
	"context"
	"time"
)

// Quality 物品品质
type Quality int

const (
	QualityCommon    Quality = iota + 1 // 普通
	QualityUncommon                     // 不常见
	QualityRare                         // 稀有
	QualityEpic                         // 史诗
	QualityLegendary                    // 传说
	QualityMythic                       // 神话
)

// Repository 背包仓储接口
type Repository interface {
	// 基础CRUD操作
	Save(ctx context.Context, inventory *Inventory) error
	FindByPlayerID(ctx context.Context, playerID string) (*Inventory, error)
	Delete(ctx context.Context, playerID string) error
	Exists(ctx context.Context, playerID string) (bool, error)

	// 批量操作
	SaveBatch(ctx context.Context, inventories []*Inventory) error
	FindByPlayerIDs(ctx context.Context, playerIDs []string) ([]*Inventory, error)

	// 查询操作
	FindItemsByType(ctx context.Context, playerID string, itemType ItemType) ([]*Item, error)
	FindExpiredItems(ctx context.Context, playerID string, before time.Time) ([]*Item, error)
	CountItemsByType(ctx context.Context, playerID string, itemType ItemType) (int64, error)

	// 统计操作
	GetInventoryStats(ctx context.Context, playerID string) (*InventoryStats, error)
	GetPlayerItemHistory(ctx context.Context, playerID string, limit int) ([]*ItemHistory, error)
}

// InventoryStats 背包统计信息
type InventoryStats struct {
	PlayerID       string             `json:"player_id"`
	TotalItems     int64              `json:"total_items"`
	UsedSlots      int                `json:"used_slots"`
	Capacity       int                `json:"capacity"`
	ItemsByType    map[ItemType]int64 `json:"items_by_type"`
	ItemsByQuality map[Quality]int64  `json:"items_by_quality"`
	LastUpdate     time.Time          `json:"last_update"`
}

// ItemHistory 物品历史记录
type ItemHistory struct {
	ID         string                 `json:"id"`
	PlayerID   string                 `json:"player_id"`
	ItemID     string                 `json:"item_id"`
	Action     string                 `json:"action"` // add, remove, use, trade
	Quantity   int64                  `json:"quantity"`
	Reason     string                 `json:"reason"`
	OccurredAt time.Time              `json:"occurred_at"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// ItemQueryFilter 物品查询过滤器
type ItemQueryFilter struct {
	PlayerID    string     `json:"player_id"`
	ItemTypes   []ItemType `json:"item_types,omitempty"`
	Qualities   []Quality  `json:"qualities,omitempty"`
	MinQuantity *int64     `json:"min_quantity,omitempty"`
	MaxQuantity *int64     `json:"max_quantity,omitempty"`
	ExpiredOnly bool       `json:"expired_only"`
	UsableOnly  bool       `json:"usable_only"`
	Limit       int        `json:"limit"`
	Offset      int        `json:"offset"`
}

// ItemRepository 物品仓储接口
type ItemRepository interface {
	// Save 保存物品
	Save(ctx context.Context, item *Item) error

	// FindByID 根据ID查找物品
	FindByID(ctx context.Context, id ItemID) (*Item, error)

	// FindByType 根据类型查找物品
	FindByType(ctx context.Context, itemType ItemType, limit int) ([]*Item, error)

	// FindByRarity 根据稀有度查找物品
	FindByRarity(ctx context.Context, rarity ItemRarity, limit int) ([]*Item, error)

	// Update 更新物品
	Update(ctx context.Context, item *Item) error

	// Delete 删除物品
	Delete(ctx context.Context, itemID string) error
}
