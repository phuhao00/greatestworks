package scene

import (
	"context"
	"time"
)

// Repository 场景仓储接口
type Repository interface {
	// 基础CRUD操作
	Save(ctx context.Context, scene *Scene) error
	FindByID(ctx context.Context, sceneID string) (*Scene, error)
	Delete(ctx context.Context, sceneID string) error
	Exists(ctx context.Context, sceneID string) (bool, error)
	
	// 批量操作
	SaveBatch(ctx context.Context, scenes []*Scene) error
	FindByIDs(ctx context.Context, sceneIDs []string) ([]*Scene, error)
	FindAll(ctx context.Context) ([]*Scene, error)
	
	// 场景查询
	FindByType(ctx context.Context, sceneType SceneType) ([]*Scene, error)
	FindByStatus(ctx context.Context, status SceneStatus) ([]*Scene, error)
	FindAvailableScenes(ctx context.Context) ([]*Scene, error)
	FindScenesWithSpace(ctx context.Context, minSpace int) ([]*Scene, error)
	
	// 实体管理
	SaveEntity(ctx context.Context, sceneID string, entity Entity) error
	RemoveEntity(ctx context.Context, sceneID string, entityID string) error
	FindEntitiesByType(ctx context.Context, sceneID string, entityType EntityType) ([]Entity, error)
	FindEntitiesInRadius(ctx context.Context, sceneID string, center *Position, radius float64) ([]Entity, error)
	
	// 玩家管理
	AddPlayerToScene(ctx context.Context, sceneID string, playerID string) error
	RemovePlayerFromScene(ctx context.Context, sceneID string, playerID string) error
	FindPlayerScene(ctx context.Context, playerID string) (*Scene, error)
	GetScenePlayerCount(ctx context.Context, sceneID string) (int, error)
	GetScenePlayers(ctx context.Context, sceneID string) ([]string, error)
	
	// 统计查询
	GetSceneStats(ctx context.Context, sceneID string) (*SceneStats, error)
	GetSceneHistory(ctx context.Context, sceneID string, limit int) ([]*SceneHistoryRecord, error)
	GetPopularScenes(ctx context.Context, limit int) ([]*ScenePopularity, error)
	
	// 配置管理
	GetSceneConfig(ctx context.Context, sceneID string) (*SceneConfig, error)
	SaveSceneConfig(ctx context.Context, config *SceneConfig) error
	GetAllSceneConfigs(ctx context.Context) ([]*SceneConfig, error)
}

// SceneStats 场景统计信息
type SceneStats struct {
	SceneID         string            `json:"scene_id"`
	SceneName       string            `json:"scene_name"`
	SceneType       SceneType         `json:"scene_type"`
	CurrentPlayers  int               `json:"current_players"`
	MaxPlayers      int               `json:"max_players"`
	TotalEntities   int               `json:"total_entities"`
	EntitiesByType  map[EntityType]int `json:"entities_by_type"`
	ActiveMonsters  int               `json:"active_monsters"`
	ActiveNPCs      int               `json:"active_npcs"`
	DroppedItems    int               `json:"dropped_items"`
	AveragePlayerLevel int            `json:"average_player_level"`
	PeakPlayerCount int               `json:"peak_player_count"`
	LastUpdate      time.Time         `json:"last_update"`
	Uptime          time.Duration     `json:"uptime"`
}

// SceneHistoryRecord 场景历史记录
type SceneHistoryRecord struct {
	ID          string      `json:"id"`
	SceneID     string      `json:"scene_id"`
	EventType   string      `json:"event_type"` // player_entered, player_left, monster_spawned, etc.
	EntityID    string      `json:"entity_id"`
	EntityType  EntityType  `json:"entity_type"`
	Position    *Position   `json:"position,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
	OccurredAt  time.Time   `json:"occurred_at"`
}

// ScenePopularity 场景热度
type ScenePopularity struct {
	SceneID       string    `json:"scene_id"`
	SceneName     string    `json:"scene_name"`
	SceneType     SceneType `json:"scene_type"`
	PlayerCount   int       `json:"current_player_count"`
	PeakCount     int       `json:"peak_player_count"`
	TotalVisits   int64     `json:"total_visits"`
	AverageStay   time.Duration `json:"average_stay_duration"`
	PopularityScore float64 `json:"popularity_score"`
	LastActive    time.Time `json:"last_active"`
}

// SceneConfig 场景配置
type SceneConfig struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	SceneType       SceneType              `json:"scene_type"`
	Width           float64                `json:"width"`
	Height          float64                `json:"height"`
	MaxPlayers      int                    `json:"max_players"`
	MinLevel        int                    `json:"min_level"`
	MaxLevel        int                    `json:"max_level"`
	PvPEnabled      bool                   `json:"pvp_enabled"`
	RespawnEnabled  bool                   `json:"respawn_enabled"`
	DropEnabled     bool                   `json:"drop_enabled"`
	SpawnPoints     []*SpawnPointConfig    `json:"spawn_points"`
	Portals         []*PortalConfig        `json:"portals"`
	NPCs            []*NPCConfig           `json:"npcs"`
	Environment     *EnvironmentConfig     `json:"environment"`
	Restrictions    *SceneRestrictions     `json:"restrictions"`
	Enabled         bool                   `json:"enabled"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// SpawnPointConfig 刷新点配置
type SpawnPointConfig struct {
	ID          string        `json:"id"`
	Position    *Position     `json:"position"`
	SpawnType   SpawnType     `json:"spawn_type"`
	TargetID    string        `json:"target_id"`
	Interval    time.Duration `json:"interval"`
	MaxCount    int           `json:"max_count"`
	Active      bool          `json:"active"`
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
}

// PortalConfig 传送门配置
type PortalConfig struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Position        *Position `json:"position"`
	TargetSceneID   string    `json:"target_scene_id"`
	TargetPosition  *Position `json:"target_position"`
	RequiredLevel   int       `json:"required_level"`
	RequiredItems   []string  `json:"required_items"`
	Cost            int64     `json:"cost"`
	Active          bool      `json:"active"`
	Bidirectional   bool      `json:"bidirectional"`
}

// NPCConfig NPC配置
type NPCConfig struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	NPCType     NPCType     `json:"npc_type"`
	Position    *Position   `json:"position"`
	Level       int         `json:"level"`
	Health      int64       `json:"health"`
	AI          *AIConfig   `json:"ai"`
	Dialogues   []string    `json:"dialogues"`
	Shop        *ShopConfig `json:"shop,omitempty"`
	Quests      []string    `json:"quests"`
	Active      bool        `json:"active"`
	Respawn     bool        `json:"respawn"`
}

// AIConfig AI配置
type AIConfig struct {
	BehaviorType BehaviorType `json:"behavior_type"`
	PatrolPath   []*Position  `json:"patrol_path"`
	AggroRange   float64      `json:"aggro_range"`
	ChaseRange   float64      `json:"chase_range"`
	ReturnRange  float64      `json:"return_range"`
	AttackRange  float64      `json:"attack_range"`
	MoveSpeed    float64      `json:"move_speed"`
	AttackSpeed  float64      `json:"attack_speed"`
}

// ShopConfig 商店配置
type ShopConfig struct {
	Items       []string `json:"items"`
	RefreshRate time.Duration `json:"refresh_rate"`
	DiscountRate float64 `json:"discount_rate"`
}

// EnvironmentConfig 环境配置
type EnvironmentConfig struct {
	Weather     string  `json:"weather"`
	TimeOfDay   string  `json:"time_of_day"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
	Visibility  float64 `json:"visibility"`
	AmbientSound string `json:"ambient_sound"`
	BackgroundMusic string `json:"background_music"`
}

// SceneRestrictions 场景限制
type SceneRestrictions struct {
	ClassRestrictions []string `json:"class_restrictions"`
	RaceRestrictions  []string `json:"race_restrictions"`
	GuildOnly         bool     `json:"guild_only"`
	PartyOnly         bool     `json:"party_only"`
	VIPOnly           bool     `json:"vip_only"`
	TimeRestrictions  *TimeRestriction `json:"time_restrictions"`
}

// TimeRestriction 时间限制
type TimeRestriction struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	DaysOfWeek []int    `json:"days_of_week"` // 0=Sunday, 1=Monday, etc.
}

// SceneQueryFilter 场景查询过滤器
type SceneQueryFilter struct {
	SceneTypes      []SceneType   `json:"scene_types,omitempty"`
	Statuses        []SceneStatus `json:"statuses,omitempty"`
	MinPlayers      *int          `json:"min_players,omitempty"`
	MaxPlayers      *int          `json:"max_players,omitempty"`
	MinLevel        *int          `json:"min_level,omitempty"`
	MaxLevel        *int          `json:"max_level,omitempty"`
	PvPEnabled      *bool         `json:"pvp_enabled,omitempty"`
	HasSpace        bool          `json:"has_space"`
	ActiveOnly      bool          `json:"active_only"`
	SortBy          string        `json:"sort_by"` // player_count, popularity, name
	SortOrder       string        `json:"sort_order"` // asc, desc
	Limit           int           `json:"limit"`
	Offset          int           `json:"offset"`
}

// EntityQueryFilter 实体查询过滤器
type EntityQueryFilter struct {
	SceneID     string       `json:"scene_id"`
	EntityTypes []EntityType `json:"entity_types,omitempty"`
	ActiveOnly  bool         `json:"active_only"`
	Center      *Position    `json:"center,omitempty"`
	Radius      *float64     `json:"radius,omitempty"`
	MinLevel    *int         `json:"min_level,omitempty"`
	MaxLevel    *int         `json:"max_level,omitempty"`
	Limit       int          `json:"limit"`
	Offset      int          `json:"offset"`
}