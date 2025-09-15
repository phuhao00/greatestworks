package npc

import (
	"time"
)

// NPCRepository NPC仓储接口
type NPCRepository interface {
	// 基础CRUD操作
	Save(npc *NPCAggregate) error
	FindByID(id string) (*NPCAggregate, error)
	FindByType(npcType NPCType) ([]*NPCAggregate, error)
	FindByStatus(status NPCStatus) ([]*NPCAggregate, error)
	Update(npc *NPCAggregate) error
	Delete(id string) error
	
	// 位置相关查询
	FindByLocation(location *Location, radius float64) ([]*NPCAggregate, error)
	FindByRegion(region string) ([]*NPCAggregate, error)
	FindByZone(zone string) ([]*NPCAggregate, error)
	
	// 分页查询
	FindWithPagination(query *NPCQuery) (*NPCPageResult, error)
	
	// 统计操作
	Count() (int64, error)
	CountByType(npcType NPCType) (int64, error)
	CountByStatus(status NPCStatus) (int64, error)
	CountByRegion(region string) (int64, error)
	
	// 批量操作
	SaveBatch(npcs []*NPCAggregate) error
	DeleteBatch(ids []string) error
	
	// 高级查询
	FindActiveNPCs() ([]*NPCAggregate, error)
	FindNPCsWithShops() ([]*NPCAggregate, error)
	FindNPCsWithQuests() ([]*NPCAggregate, error)
	FindNearbyNPCs(location *Location, radius float64, npcType NPCType) ([]*NPCAggregate, error)
}

// DialogueRepository 对话仓储接口
type DialogueRepository interface {
	// 基础CRUD操作
	Save(dialogue *Dialogue) error
	FindByID(id string) (*Dialogue, error)
	FindByNPC(npcID string) ([]*Dialogue, error)
	FindByType(dialogueType DialogueType) ([]*Dialogue, error)
	Update(dialogue *Dialogue) error
	Delete(id string) error
	
	// 分页查询
	FindWithPagination(query *DialogueQuery) (*DialoguePageResult, error)
	
	// 统计操作
	Count() (int64, error)
	CountByType(dialogueType DialogueType) (int64, error)
	CountByNPC(npcID string) (int64, error)
	
	// 会话相关
	SaveSession(session *DialogueSession) error
	FindSession(npcID, playerID string) (*DialogueSession, error)
	FindActiveSessions(playerID string) ([]*DialogueSession, error)
	DeleteSession(npcID, playerID string) error
	CleanupExpiredSessions() error
}

// QuestRepository 任务仓储接口
type QuestRepository interface {
	// 基础CRUD操作
	Save(quest *Quest) error
	FindByID(id string) (*Quest, error)
	FindByNPC(npcID string) ([]*Quest, error)
	FindByType(questType QuestType) ([]*Quest, error)
	Update(quest *Quest) error
	Delete(id string) error
	
	// 分页查询
	FindWithPagination(query *QuestQuery) (*QuestPageResult, error)
	
	// 统计操作
	Count() (int64, error)
	CountByType(questType QuestType) (int64, error)
	CountByNPC(npcID string) (int64, error)
	
	// 任务实例相关
	SaveInstance(instance *QuestInstance) error
	FindInstance(questID, playerID string) (*QuestInstance, error)
	FindInstancesByPlayer(playerID string) ([]*QuestInstance, error)
	FindInstancesByQuest(questID string) ([]*QuestInstance, error)
	FindInstancesByStatus(status QuestStatus) ([]*QuestInstance, error)
	UpdateInstance(instance *QuestInstance) error
	DeleteInstance(questID, playerID string) error
	
	// 任务进度
	UpdateProgress(questID, playerID, objectiveID string, progress int) error
	GetProgress(questID, playerID string) (map[string]int, error)
	
	// 任务完成统计
	GetCompletionStats(questID string) (*QuestCompletionStats, error)
	GetPlayerQuestStats(playerID string) (*PlayerQuestStats, error)
}

// ShopRepository 商店仓储接口
type ShopRepository interface {
	// 基础CRUD操作
	Save(shop *Shop) error
	FindByID(id string) (*Shop, error)
	FindByNPC(npcID string) (*Shop, error)
	Update(shop *Shop) error
	Delete(id string) error
	
	// 商品相关
	SaveItem(shopID string, item *ShopItem) error
	FindItem(shopID, itemID string) (*ShopItem, error)
	FindItemsByShop(shopID string) ([]*ShopItem, error)
	UpdateItem(shopID string, item *ShopItem) error
	DeleteItem(shopID, itemID string) error
	
	// 交易记录
	SaveTradeRecord(record *TradeRecord) error
	FindTradeRecords(shopID string, limit int) ([]*TradeRecord, error)
	FindPlayerTradeRecords(playerID string, limit int) ([]*TradeRecord, error)
	
	// 统计操作
	Count() (int64, error)
	GetShopStats(shopID string) (*ShopStatistics, error)
	GetTradeStats(shopID string, startTime, endTime time.Time) (*TradeStatistics, error)
}

// RelationshipRepository 关系仓储接口
type RelationshipRepository interface {
	// 基础CRUD操作
	Save(relationship *Relationship) error
	FindByID(playerID, npcID string) (*Relationship, error)
	FindByPlayer(playerID string) ([]*Relationship, error)
	FindByNPC(npcID string) ([]*Relationship, error)
	Update(relationship *Relationship) error
	Delete(playerID, npcID string) error
	
	// 关系等级查询
	FindByLevel(level RelationshipLevel) ([]*Relationship, error)
	FindByValueRange(minValue, maxValue int) ([]*Relationship, error)
	
	// 分页查询
	FindWithPagination(query *RelationshipQuery) (*RelationshipPageResult, error)
	
	// 统计操作
	Count() (int64, error)
	CountByLevel(level RelationshipLevel) (int64, error)
	GetAverageRelationship(npcID string) (float64, error)
	
	// 关系历史
	SaveRelationshipEvent(event *RelationshipEvent) error
	FindRelationshipHistory(playerID, npcID string, limit int) ([]*RelationshipEvent, error)
	
	// 排行榜
	GetTopRelationships(npcID string, limit int) ([]*Relationship, error)
	GetPlayerRanking(playerID, npcID string) (int, error)
}

// NPCStatisticsRepository NPC统计仓储接口
type NPCStatisticsRepository interface {
	// 保存统计数据
	SaveStatistics(stats *NPCStatistics) error
	UpdateStatistics(stats *NPCStatistics) error
	
	// 查询统计数据
	FindStatistics(npcID string) (*NPCStatistics, error)
	FindStatisticsByType(npcType NPCType) ([]*NPCStatistics, error)
	
	// 全局统计
	GetGlobalStatistics() (*GlobalNPCStatistics, error)
	GetTypeStatistics(npcType NPCType) (*TypeNPCStatistics, error)
	
	// 趋势分析
	GetInteractionTrend(npcID string, days int) ([]*InteractionTrendData, error)
	GetPopularityTrend(npcType NPCType, days int) ([]*PopularityTrendData, error)
	
	// 活跃度统计
	GetActiveNPCCount(timeRange time.Duration) (int64, error)
	GetMostActiveNPCs(limit int) ([]*NPCStatistics, error)
}

// 查询条件结构体

// NPCQuery NPC查询条件
type NPCQuery struct {
	Name       string
	Type       *NPCType
	Status     *NPCStatus
	Region     string
	Zone       string
	Location   *Location
	Radius     *float64
	HasShop    *bool
	HasQuests  *bool
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	UpdatedAfter  *time.Time
	UpdatedBefore *time.Time
	OrderBy    string
	OrderDesc  bool
	Offset     int
	Limit      int
}

// DialogueQuery 对话查询条件
type DialogueQuery struct {
	NPCID      string
	Type       *DialogueType
	PlayerID   string
	Available  *bool
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	OrderBy    string
	OrderDesc  bool
	Offset     int
	Limit      int
}

// QuestQuery 任务查询条件
type QuestQuery struct {
	NPCID      string
	Type       *QuestType
	PlayerID   string
	Status     *QuestStatus
	Repeatable *bool
	DailyReset *bool
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	OrderBy    string
	OrderDesc  bool
	Offset     int
	Limit      int
}

// RelationshipQuery 关系查询条件
type RelationshipQuery struct {
	PlayerID   string
	NPCID      string
	Level      *RelationshipLevel
	MinValue   *int
	MaxValue   *int
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	UpdatedAfter  *time.Time
	UpdatedBefore *time.Time
	OrderBy    string
	OrderDesc  bool
	Offset     int
	Limit      int
}

// 分页结果结构体

// NPCPageResult NPC分页结果
type NPCPageResult struct {
	Items      []*NPCAggregate
	Total      int64
	Offset     int
	Limit      int
	HasMore    bool
}

// DialoguePageResult 对话分页结果
type DialoguePageResult struct {
	Items      []*Dialogue
	Total      int64
	Offset     int
	Limit      int
	HasMore    bool
}

// QuestPageResult 任务分页结果
type QuestPageResult struct {
	Items      []*Quest
	Total      int64
	Offset     int
	Limit      int
	HasMore    bool
}

// RelationshipPageResult 关系分页结果
type RelationshipPageResult struct {
	Items      []*Relationship
	Total      int64
	Offset     int
	Limit      int
	HasMore    bool
}

// 统计数据结构体

// QuestCompletionStats 任务完成统计
type QuestCompletionStats struct {
	QuestID        string
	TotalAttempts  int64
	TotalCompleted int64
	CompletionRate float64
	AverageTime    time.Duration
	LastCompleted  time.Time
}

// PlayerQuestStats 玩家任务统计
type PlayerQuestStats struct {
	PlayerID         string
	TotalQuests      int64
	CompletedQuests  int64
	FailedQuests     int64
	ActiveQuests     int64
	CompletionRate   float64
	AverageTime      time.Duration
	FavoriteType     QuestType
	LastQuestTime    time.Time
}

// ShopStatistics 商店统计
type ShopStatistics struct {
	ShopID         string
	TotalTrades    int64
	TotalRevenue   int64
	TotalItems     int64
	PopularItem    string
	AveragePrice   float64
	LastTradeTime  time.Time
	CreatedAt      time.Time
}

// TradeStatistics 交易统计
type TradeStatistics struct {
	ShopID       string
	PeriodStart  time.Time
	PeriodEnd    time.Time
	TotalTrades  int64
	TotalRevenue int64
	TotalItems   int64
	TopItems     []string
	TopCustomers []string
}

// GlobalNPCStatistics 全局NPC统计
type GlobalNPCStatistics struct {
	TotalNPCs        int64
	ActiveNPCs       int64
	NPCsByType       map[NPCType]int64
	NPCsByStatus     map[NPCStatus]int64
	TotalDialogues   int64
	TotalQuests      int64
	TotalShops       int64
	TotalRelationships int64
	AverageRelationship float64
	MostPopularNPC   string
	MostActiveRegion string
	UpdatedAt        time.Time
}

// TypeNPCStatistics 类型NPC统计
type TypeNPCStatistics struct {
	NPCType          NPCType
	TotalCount       int64
	ActiveCount      int64
	AverageLevel     float64
	TotalDialogues   int64
	TotalQuests      int64
	TotalShops       int64
	AverageRelationship float64
	MostPopularNPC   string
	UpdatedAt        time.Time
}

// InteractionTrendData 交互趋势数据
type InteractionTrendData struct {
	Date            time.Time
	DialogueCount   int64
	QuestCount      int64
	TradeCount      int64
	UniqueVisitors  int64
}

// PopularityTrendData 受欢迎程度趋势数据
type PopularityTrendData struct {
	Date           time.Time
	NPCType        NPCType
	InteractionCount int64
	UniquePlayers  int64
	AverageRating  float64
}

// TradeRecord 交易记录
type TradeRecord struct {
	ID         string
	ShopID     string
	PlayerID   string
	ItemID     string
	Quantity   int
	Price      int
	TotalPrice int
	Timestamp  time.Time
}

// NewTradeRecord 创建交易记录
func NewTradeRecord(shopID, playerID, itemID string, quantity, price int) *TradeRecord {
	return &TradeRecord{
		ID:         fmt.Sprintf("trade_%d", time.Now().UnixNano()),
		ShopID:     shopID,
		PlayerID:   playerID,
		ItemID:     itemID,
		Quantity:   quantity,
		Price:      price,
		TotalPrice: quantity * price,
		Timestamp:  time.Now(),
	}
}

// 缓存接口

// NPCCacheRepository NPC缓存仓储接口
type NPCCacheRepository interface {
	// NPC缓存
	SetNPC(id string, npc *NPCAggregate, ttl time.Duration) error
	GetNPC(id string) (*NPCAggregate, error)
	DeleteNPC(id string) error
	
	// 对话缓存
	SetDialogue(id string, dialogue *Dialogue, ttl time.Duration) error
	GetDialogue(id string) (*Dialogue, error)
	DeleteDialogue(id string) error
	
	// 任务缓存
	SetQuest(id string, quest *Quest, ttl time.Duration) error
	GetQuest(id string) (*Quest, error)
	DeleteQuest(id string) error
	
	// 关系缓存
	SetRelationship(playerID, npcID string, relationship *Relationship, ttl time.Duration) error
	GetRelationship(playerID, npcID string) (*Relationship, error)
	DeleteRelationship(playerID, npcID string) error
	
	// 会话缓存
	SetSession(npcID, playerID string, session *DialogueSession, ttl time.Duration) error
	GetSession(npcID, playerID string) (*DialogueSession, error)
	DeleteSession(npcID, playerID string) error
	
	// 统计缓存
	SetStatistics(key string, stats interface{}, ttl time.Duration) error
	GetStatistics(key string, result interface{}) error
	DeleteStatistics(key string) error
	
	// 位置索引缓存
	SetLocationIndex(region string, npcs []*NPCAggregate, ttl time.Duration) error
	GetLocationIndex(region string) ([]*NPCAggregate, error)
	DeleteLocationIndex(region string) error
	
	// 批量操作
	SetBatch(items map[string]interface{}, ttl time.Duration) error
	GetBatch(keys []string) (map[string]interface{}, error)
	DeleteBatch(keys []string) error
	
	// 缓存管理
	Clear() error
	Exists(key string) (bool, error)
	SetTTL(key string, ttl time.Duration) error
	GetTTL(key string) (time.Duration, error)
}

// 事务接口

// NPCTransactionRepository NPC事务仓储接口
type NPCTransactionRepository interface {
	// 事务管理
	BeginTransaction() (NPCTransaction, error)
	CommitTransaction(tx NPCTransaction) error
	RollbackTransaction(tx NPCTransaction) error
	
	// 在事务中执行操作
	ExecuteInTransaction(fn func(tx NPCTransaction) error) error
}

// NPCTransaction NPC事务接口
type NPCTransaction interface {
	// NPC操作
	SaveNPC(npc *NPCAggregate) error
	UpdateNPC(npc *NPCAggregate) error
	DeleteNPC(id string) error
	
	// 对话操作
	SaveDialogue(dialogue *Dialogue) error
	UpdateDialogue(dialogue *Dialogue) error
	DeleteDialogue(id string) error
	
	// 任务操作
	SaveQuest(quest *Quest) error
	UpdateQuest(quest *Quest) error
	DeleteQuest(id string) error
	
	// 任务实例操作
	SaveQuestInstance(instance *QuestInstance) error
	UpdateQuestInstance(instance *QuestInstance) error
	DeleteQuestInstance(questID, playerID string) error
	
	// 关系操作
	SaveRelationship(relationship *Relationship) error
	UpdateRelationship(relationship *Relationship) error
	DeleteRelationship(playerID, npcID string) error
	
	// 商店操作
	SaveShop(shop *Shop) error
	UpdateShop(shop *Shop) error
	DeleteShop(id string) error
	
	// 交易记录
	SaveTradeRecord(record *TradeRecord) error
	
	// 统计操作
	UpdateStatistics(stats *NPCStatistics) error
	
	// 事务状态
	IsActive() bool
	GetID() string
}

// 仓储工厂接口

// NPCRepositoryFactory NPC仓储工厂接口
type NPCRepositoryFactory interface {
	// 创建仓储实例
	CreateNPCRepository() NPCRepository
	CreateDialogueRepository() DialogueRepository
	CreateQuestRepository() QuestRepository
	CreateShopRepository() ShopRepository
	CreateRelationshipRepository() RelationshipRepository
	CreateStatisticsRepository() NPCStatisticsRepository
	CreateCacheRepository() NPCCacheRepository
	CreateTransactionRepository() NPCTransactionRepository
	
	// 健康检查
	HealthCheck() error
	
	// 关闭连接
	Close() error
}

// 搜索接口

// NPCSearchRepository NPC搜索仓储接口
type NPCSearchRepository interface {
	// 全文搜索
	SearchNPCs(query string, filters map[string]interface{}) ([]*NPCAggregate, error)
	SearchDialogues(query string, filters map[string]interface{}) ([]*Dialogue, error)
	SearchQuests(query string, filters map[string]interface{}) ([]*Quest, error)
	
	// 地理搜索
	SearchNearbyNPCs(location *Location, radius float64, filters map[string]interface{}) ([]*NPCAggregate, error)
	
	// 智能推荐
	RecommendNPCs(playerID string, limit int) ([]*NPCAggregate, error)
	RecommendQuests(playerID string, limit int) ([]*Quest, error)
	RecommendDialogues(playerID string, npcID string, limit int) ([]*Dialogue, error)
	
	// 索引管理
	RebuildIndex() error
	UpdateIndex(entityType string, entityID string, data interface{}) error
	DeleteFromIndex(entityType string, entityID string) error
}

// 导入fmt包
import "fmt"