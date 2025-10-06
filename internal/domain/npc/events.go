package npc

import (
	"fmt"
	"time"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	GetEventID() string
	GetEventType() string
	GetAggregateID() string
	GetOccurredAt() time.Time
	GetVersion() int
	GetData() map[string]interface{}
	Validate() error
}

// BaseDomainEvent 基础领域事件
type BaseDomainEvent struct {
	EventID     string
	EventType   string
	AggregateID string
	OccurredAt  time.Time
	Version     int
	Data        map[string]interface{}
}

// GetEventID 获取事件ID
func (e *BaseDomainEvent) GetEventID() string {
	return e.EventID
}

// GetEventType 获取事件类型
func (e *BaseDomainEvent) GetEventType() string {
	return e.EventType
}

// GetAggregateID 获取聚合根ID
func (e *BaseDomainEvent) GetAggregateID() string {
	return e.AggregateID
}

// GetOccurredAt 获取发生时间
func (e *BaseDomainEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetVersion 获取版本
func (e *BaseDomainEvent) GetVersion() int {
	return e.Version
}

// GetData 获取数据
func (e *BaseDomainEvent) GetData() map[string]interface{} {
	return e.Data
}

// Validate 验证事件
func (e *BaseDomainEvent) Validate() error {
	if e.EventID == "" {
		return fmt.Errorf("event ID cannot be empty")
	}
	if e.EventType == "" {
		return fmt.Errorf("event type cannot be empty")
	}
	if e.AggregateID == "" {
		return fmt.Errorf("aggregate ID cannot be empty")
	}
	if e.OccurredAt.IsZero() {
		return fmt.Errorf("occurred at cannot be zero")
	}
	return nil
}

// NPC相关事件

// NPCCreatedEvent NPC创建事件
type NPCCreatedEvent struct {
	*BaseDomainEvent
	NPCID     string
	Name      string
	Type      NPCType
	Location  *Location
	CreatedBy string
}

// NewNPCCreatedEvent 创建NPC创建事件
func NewNPCCreatedEvent(npcID, name string, npcType NPCType, location *Location, createdBy string) *NPCCreatedEvent {
	return &NPCCreatedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("npc_created_%d", time.Now().UnixNano()),
			EventType:   "NPCCreated",
			AggregateID: npcID,
			OccurredAt:  time.Now(),
			Version:     1,
			Data: map[string]interface{}{
				"npc_id":     npcID,
				"name":       name,
				"type":       npcType,
				"location":   location,
				"created_by": createdBy,
			},
		},
		NPCID:     npcID,
		Name:      name,
		Type:      npcType,
		Location:  location,
		CreatedBy: createdBy,
	}
}

// NPCStatusChangedEvent NPC状态变更事件
type NPCStatusChangedEvent struct {
	*BaseDomainEvent
	NPCID     string
	OldStatus NPCStatus
	NewStatus NPCStatus
	Reason    string
	ChangedBy string
}

// NewNPCStatusChangedEvent 创建NPC状态变更事件
func NewNPCStatusChangedEvent(npcID string, oldStatus, newStatus NPCStatus, reason, changedBy string) *NPCStatusChangedEvent {
	return &NPCStatusChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("npc_status_changed_%d", time.Now().UnixNano()),
			EventType:   "NPCStatusChanged",
			AggregateID: npcID,
			OccurredAt:  time.Now(),
			Version:     1,
			Data: map[string]interface{}{
				"npc_id":     npcID,
				"old_status": oldStatus,
				"new_status": newStatus,
				"reason":     reason,
				"changed_by": changedBy,
			},
		},
		NPCID:     npcID,
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Reason:    reason,
		ChangedBy: changedBy,
	}
}

// NPCLocationChangedEvent NPC位置变更事件
type NPCLocationChangedEvent struct {
	*BaseDomainEvent
	NPCID       string
	OldLocation *Location
	NewLocation *Location
	MoveReason  string
}

// NewNPCLocationChangedEvent 创建NPC位置变更事件
func NewNPCLocationChangedEvent(npcID string, oldLocation, newLocation *Location, moveReason string) *NPCLocationChangedEvent {
	return &NPCLocationChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("npc_location_changed_%d", time.Now().UnixNano()),
			EventType:   "NPCLocationChanged",
			AggregateID: npcID,
			OccurredAt:  time.Now(),
			Version:     1,
			Data: map[string]interface{}{
				"npc_id":       npcID,
				"old_location": oldLocation,
				"new_location": newLocation,
				"move_reason":  moveReason,
			},
		},
		NPCID:       npcID,
		OldLocation: oldLocation,
		NewLocation: newLocation,
		MoveReason:  moveReason,
	}
}

// 对话相关事件

// DialogueStartedEvent 对话开始事件
type DialogueStartedEvent struct {
	*BaseDomainEvent
	NPCID      string
	PlayerID   string
	DialogueID string
	SessionID  string
}

// NewDialogueStartedEvent 创建对话开始事件
func NewDialogueStartedEvent(npcID, playerID, dialogueID, sessionID string) *DialogueStartedEvent {
	return &DialogueStartedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("dialogue_started_%d", time.Now().UnixNano()),
			EventType:   "DialogueStarted",
			AggregateID: npcID,
			OccurredAt:  time.Now(),
			Version:     1,
			Data: map[string]interface{}{
				"npc_id":      npcID,
				"player_id":   playerID,
				"dialogue_id": dialogueID,
				"session_id":  sessionID,
			},
		},
		NPCID:      npcID,
		PlayerID:   playerID,
		DialogueID: dialogueID,
		SessionID:  sessionID,
	}
}

// DialogueEndedEvent 对话结束事件
type DialogueEndedEvent struct {
	*BaseDomainEvent
	NPCID      string
	PlayerID   string
	DialogueID string
	SessionID  string
	Duration   time.Duration
	EndReason  string
}

// NewDialogueEndedEvent 创建对话结束事件
func NewDialogueEndedEvent(npcID, playerID, dialogueID, sessionID string, duration time.Duration, endReason string) *DialogueEndedEvent {
	return &DialogueEndedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("dialogue_ended_%d", time.Now().UnixNano()),
			EventType:   "DialogueEnded",
			AggregateID: npcID,
			OccurredAt:  time.Now(),
			Version:     1,
			Data: map[string]interface{}{
				"npc_id":      npcID,
				"player_id":   playerID,
				"dialogue_id": dialogueID,
				"session_id":  sessionID,
				"duration":    duration,
				"end_reason":  endReason,
			},
		},
		NPCID:      npcID,
		PlayerID:   playerID,
		DialogueID: dialogueID,
		SessionID:  sessionID,
		Duration:   duration,
		EndReason:  endReason,
	}
}

// 任务相关事件

// QuestAcceptedEvent 任务接受事件
type QuestAcceptedEvent struct {
	*BaseDomainEvent
	NPCID    string
	PlayerID string
	QuestID  string
}

// NewQuestAcceptedEvent 创建任务接受事件
func NewQuestAcceptedEvent(npcID, playerID, questID string) *QuestAcceptedEvent {
	return &QuestAcceptedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("quest_accepted_%d", time.Now().UnixNano()),
			EventType:   "QuestAccepted",
			AggregateID: npcID,
			OccurredAt:  time.Now(),
			Version:     1,
			Data: map[string]interface{}{
				"npc_id":    npcID,
				"player_id": playerID,
				"quest_id":  questID,
			},
		},
		NPCID:    npcID,
		PlayerID: playerID,
		QuestID:  questID,
	}
}

// QuestCompletedEvent 任务完成事件
type QuestCompletedEvent struct {
	*BaseDomainEvent
	NPCID      string
	PlayerID   string
	QuestID    string
	Rewards    []QuestReward
	Completion time.Duration
}

// NewQuestCompletedEvent 创建任务完成事件
func NewQuestCompletedEvent(npcID, playerID, questID string, rewards []QuestReward, completion time.Duration) *QuestCompletedEvent {
	return &QuestCompletedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("quest_completed_%d", time.Now().UnixNano()),
			EventType:   "QuestCompleted",
			AggregateID: npcID,
			OccurredAt:  time.Now(),
			Version:     1,
			Data: map[string]interface{}{
				"npc_id":     npcID,
				"player_id":  playerID,
				"quest_id":   questID,
				"rewards":    rewards,
				"completion": completion,
			},
		},
		NPCID:      npcID,
		PlayerID:   playerID,
		QuestID:    questID,
		Rewards:    rewards,
		Completion: completion,
	}
}

// 商店相关事件

// TradeCompletedEvent 交易完成事件
type TradeCompletedEvent struct {
	*BaseDomainEvent
	NPCID      string
	PlayerID   string
	ShopID     string
	ItemID     string
	Quantity   int
	Price      int
	TotalPrice int
}

// NewTradeCompletedEvent 创建交易完成事件
func NewTradeCompletedEvent(npcID, playerID, shopID, itemID string, quantity, price, totalPrice int) *TradeCompletedEvent {
	return &TradeCompletedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("trade_completed_%d", time.Now().UnixNano()),
			EventType:   "TradeCompleted",
			AggregateID: npcID,
			OccurredAt:  time.Now(),
			Version:     1,
			Data: map[string]interface{}{
				"npc_id":      npcID,
				"player_id":   playerID,
				"shop_id":     shopID,
				"item_id":     itemID,
				"quantity":    quantity,
				"price":       price,
				"total_price": totalPrice,
			},
		},
		NPCID:      npcID,
		PlayerID:   playerID,
		ShopID:     shopID,
		ItemID:     itemID,
		Quantity:   quantity,
		Price:      price,
		TotalPrice: totalPrice,
	}
}

// 关系相关事件

// RelationshipChangedEvent 关系变更事件
type RelationshipChangedEvent struct {
	*BaseDomainEvent
	NPCID      string
	PlayerID   string
	OldValue   int
	NewValue   int
	OldLevel   RelationshipLevel
	NewLevel   RelationshipLevel
	ChangeType RelationshipChangeType
	Reason     string
}

// NewRelationshipChangedEvent 创建关系变更事件
func NewRelationshipChangedEvent(npcID, playerID string, oldValue, newValue int, oldLevel, newLevel RelationshipLevel, changeType RelationshipChangeType, reason string) *RelationshipChangedEvent {
	return &RelationshipChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("relationship_changed_%d", time.Now().UnixNano()),
			EventType:   "RelationshipChanged",
			AggregateID: npcID,
			OccurredAt:  time.Now(),
			Version:     1,
			Data: map[string]interface{}{
				"npc_id":      npcID,
				"player_id":   playerID,
				"old_value":   oldValue,
				"new_value":   newValue,
				"old_level":   oldLevel,
				"new_level":   newLevel,
				"change_type": changeType,
				"reason":      reason,
			},
		},
		NPCID:      npcID,
		PlayerID:   playerID,
		OldValue:   oldValue,
		NewValue:   newValue,
		OldLevel:   oldLevel,
		NewLevel:   newLevel,
		ChangeType: changeType,
		Reason:     reason,
	}
}

// 事件处理器接口

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(event DomainEvent) error
	CanHandle(eventType string) bool
	GetHandlerName() string
}

// EventBus 事件总线接口
type EventBus interface {
	// 发布事件
	Publish(event DomainEvent) error
	PublishBatch(events []DomainEvent) error

	// 订阅事件
	Subscribe(eventType string, handler EventHandler) error
	Unsubscribe(eventType string, handlerName string) error

	// 获取订阅者
	GetSubscribers(eventType string) []EventHandler

	// 启动和停止
	Start() error
	Stop() error

	// 健康检查
	HealthCheck() error
}

// EventStore 事件存储接口
type EventStore interface {
	// 保存事件
	Save(event DomainEvent) error
	SaveBatch(events []DomainEvent) error

	// 查询事件
	FindByAggregateID(aggregateID string) ([]DomainEvent, error)
	FindByEventType(eventType string) ([]DomainEvent, error)
	FindByTimeRange(start, end time.Time) ([]DomainEvent, error)

	// 分页查询
	FindWithPagination(query *EventQuery) (*EventPageResult, error)

	// 统计
	Count() (int64, error)
	CountByType(eventType string) (int64, error)
	CountByAggregateID(aggregateID string) (int64, error)

	// 清理
	CleanupOldEvents(before time.Time) (int64, error)
}

// EventQuery 事件查询条件
type EventQuery struct {
	AggregateID string
	EventType   string
	StartTime   *time.Time
	EndTime     *time.Time
	OrderBy     string
	OrderDesc   bool
	Offset      int
	Limit       int
}

// EventPageResult 事件分页结果
type EventPageResult struct {
	Events  []DomainEvent
	Total   int64
	Offset  int
	Limit   int
	HasMore bool
}

// 事件验证器

// EventValidator 事件验证器接口
type EventValidator interface {
	Validate(event DomainEvent) error
	ValidateType(eventType string) error
	ValidateData(eventType string, data map[string]interface{}) error
}

// DefaultEventValidator 默认事件验证器
type DefaultEventValidator struct {
	validationRules map[string]func(DomainEvent) error
}

// NewDefaultEventValidator 创建默认事件验证器
func NewDefaultEventValidator() *DefaultEventValidator {
	return &DefaultEventValidator{
		validationRules: make(map[string]func(DomainEvent) error),
	}
}

// RegisterRule 注册验证规则
func (v *DefaultEventValidator) RegisterRule(eventType string, rule func(DomainEvent) error) {
	v.validationRules[eventType] = rule
}

// Validate 验证事件
func (v *DefaultEventValidator) Validate(event DomainEvent) error {
	// 基础验证
	if err := event.Validate(); err != nil {
		return err
	}

	// 类型特定验证
	if rule, exists := v.validationRules[event.GetEventType()]; exists {
		return rule(event)
	}

	return nil
}

// ValidateType 验证事件类型
func (v *DefaultEventValidator) ValidateType(eventType string) error {
	if eventType == "" {
		return fmt.Errorf("event type cannot be empty")
	}
	return nil
}

// ValidateData 验证事件数据
func (v *DefaultEventValidator) ValidateData(eventType string, data map[string]interface{}) error {
	if data == nil {
		return fmt.Errorf("event data cannot be nil")
	}
	return nil
}

// 事件监控器

// EventMonitor 事件监控器接口
type EventMonitor interface {
	// 记录事件指标
	RecordEvent(event DomainEvent) error
	RecordEventProcessed(eventType string, duration time.Duration) error
	RecordEventFailed(eventType string, err error) error

	// 获取指标
	GetEventCount(eventType string) (int64, error)
	GetProcessingTime(eventType string) (time.Duration, error)
	GetFailureRate(eventType string) (float64, error)

	// 健康检查
	GetHealthStatus() (*EventHealthStatus, error)
}

// EventHealthStatus 事件健康状态
type EventHealthStatus struct {
	TotalEvents     int64
	ProcessedEvents int64
	FailedEvents    int64
	AverageLatency  time.Duration
	ErrorRate       float64
	LastEventTime   time.Time
	Status          string
}

// 事件重放器

// EventReplayer 事件重放器接口
type EventReplayer interface {
	// 重放事件
	Replay(aggregateID string, fromVersion int) error
	ReplayAll(fromTime time.Time) error
	ReplayByType(eventType string, fromTime time.Time) error

	// 重建聚合
	RebuildAggregate(aggregateID string) error
	RebuildAllAggregates() error

	// 快照管理
	CreateSnapshot(aggregateID string) error
	LoadFromSnapshot(aggregateID string) error
}

// 事件投影器

// EventProjector 事件投影器接口
type EventProjector interface {
	// 处理事件
	Project(event DomainEvent) error
	ProjectBatch(events []DomainEvent) error

	// 重建投影
	Rebuild() error
	RebuildFrom(fromTime time.Time) error

	// 获取投影名称
	GetProjectionName() string

	// 健康检查
	HealthCheck() error
}

// 事件序列化器

// EventSerializer 事件序列化器接口
type EventSerializer interface {
	// 序列化
	Serialize(event DomainEvent) ([]byte, error)
	SerializeBatch(events []DomainEvent) ([]byte, error)

	// 反序列化
	Deserialize(data []byte) (DomainEvent, error)
	DeserializeBatch(data []byte) ([]DomainEvent, error)

	// 获取内容类型
	GetContentType() string
}

// 事件过滤器

// EventFilter 事件过滤器接口
type EventFilter interface {
	// 过滤事件
	Filter(event DomainEvent) bool

	// 获取过滤器名称
	GetFilterName() string
}

// EventTypeFilter 事件类型过滤器
type EventTypeFilter struct {
	allowedTypes map[string]bool
}

// NewEventTypeFilter 创建事件类型过滤器
func NewEventTypeFilter(allowedTypes []string) *EventTypeFilter {
	typeMap := make(map[string]bool)
	for _, eventType := range allowedTypes {
		typeMap[eventType] = true
	}
	return &EventTypeFilter{
		allowedTypes: typeMap,
	}
}

// Filter 过滤事件
func (f *EventTypeFilter) Filter(event DomainEvent) bool {
	return f.allowedTypes[event.GetEventType()]
}

// GetFilterName 获取过滤器名称
func (f *EventTypeFilter) GetFilterName() string {
	return "EventTypeFilter"
}

// 事件聚合器

// EventAggregator 事件聚合器接口
type EventAggregator interface {
	// 聚合事件
	Aggregate(events []DomainEvent) (map[string]interface{}, error)

	// 获取聚合器名称
	GetAggregatorName() string
}

// 事件调度器

// EventScheduler 事件调度器接口
type EventScheduler interface {
	// 调度事件
	Schedule(event DomainEvent, delay time.Duration) error
	ScheduleAt(event DomainEvent, at time.Time) error

	// 取消调度
	Cancel(eventID string) error

	// 获取调度状态
	GetScheduledEvents() ([]ScheduledEvent, error)

	// 启动和停止
	Start() error
	Stop() error
}

// ScheduledEvent 调度事件
type ScheduledEvent struct {
	ID          string
	Event       DomainEvent
	ScheduledAt time.Time
	ExecuteAt   time.Time
	Status      string
	RetryCount  int
	MaxRetries  int
	LastError   string
}

// 事件工厂

// EventFactory 事件工厂接口
type EventFactory interface {
	// 创建事件
	CreateEvent(eventType string, aggregateID string, data map[string]interface{}) (DomainEvent, error)

	// 注册事件类型
	RegisterEventType(eventType string, factory func(string, map[string]interface{}) (DomainEvent, error)) error

	// 获取支持的事件类型
	GetSupportedEventTypes() []string
}

// NPCNameChangedEvent NPC名称变更事件
type NPCNameChangedEvent struct {
	*BaseDomainEvent
	OldName string
	NewName string
}

// NewNPCNameChangedEvent 创建NPC名称变更事件
func NewNPCNameChangedEvent(npcID, oldName, newName string) *NPCNameChangedEvent {
	return &NPCNameChangedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("npc_name_changed_%d", time.Now().UnixNano()),
			EventType:   "NPCNameChanged",
			AggregateID: npcID,
			OccurredAt:  time.Now(),
			Version:     1,
			Data:        make(map[string]interface{}),
		},
		OldName: oldName,
		NewName: newName,
	}
}

// NPCMovedEvent NPC移动事件
type NPCMovedEvent struct {
	*BaseDomainEvent
	OldLocation *Location
	NewLocation *Location
}

// NewNPCMovedEvent 创建NPC移动事件
func NewNPCMovedEvent(npcID string, oldLocation, newLocation *Location) *NPCMovedEvent {
	return &NPCMovedEvent{
		BaseDomainEvent: &BaseDomainEvent{
			EventID:     fmt.Sprintf("npc_moved_%d", time.Now().UnixNano()),
			EventType:   "NPCMoved",
			AggregateID: npcID,
			OccurredAt:  time.Now(),
			Version:     1,
			Data:        make(map[string]interface{}),
		},
		OldLocation: oldLocation,
		NewLocation: newLocation,
	}
}
