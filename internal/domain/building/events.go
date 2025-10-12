package building

import (
	"context"
	"fmt"
	"time"
)

// BuildingEvent 建筑事件接口
type BuildingEvent interface {
	GetEventID() string
	GetEventType() string
	GetBuildingID() string
	GetTimestamp() time.Time
	GetPayload() interface{}
	Validate() error
}

// BaseBuildingEvent 建筑事件基础结构
type BaseBuildingEvent struct {
	EventID    string                 `json:"event_id"`
	EventType  string                 `json:"event_type"`
	BuildingID string                 `json:"building_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Payload    map[string]interface{} `json:"payload"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// GetEventID 获取事件ID
func (e *BaseBuildingEvent) GetEventID() string {
	return e.EventID
}

// GetEventType 获取事件类型
func (e *BaseBuildingEvent) GetEventType() string {
	return e.EventType
}

// GetBuildingID 获取建筑ID
func (e *BaseBuildingEvent) GetBuildingID() string {
	return e.BuildingID
}

// GetTimestamp 获取时间戳
func (e *BaseBuildingEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetPayload 获取载荷
func (e *BaseBuildingEvent) GetPayload() interface{} {
	return e.Payload
}

// Validate 验证事件
func (e *BaseBuildingEvent) Validate() error {
	if e.EventID == "" {
		return fmt.Errorf("event ID cannot be empty")
	}
	if e.EventType == "" {
		return fmt.Errorf("event type cannot be empty")
	}
	if e.BuildingID == "" {
		return fmt.Errorf("building ID cannot be empty")
	}
	if e.Timestamp.IsZero() {
		return fmt.Errorf("timestamp cannot be zero")
	}
	return nil
}

// SetMetadata 设置元数据
func (e *BaseBuildingEvent) SetMetadata(key string, value interface{}) {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
}

// GetMetadata 获取元数据
func (e *BaseBuildingEvent) GetMetadata(key string) (interface{}, bool) {
	if e.Metadata == nil {
		return nil, false
	}
	value, exists := e.Metadata[key]
	return value, exists
}

// 建筑生命周期事件

// BuildingCreatedEvent 建筑创建事件
type BuildingCreatedEvent struct {
	*BaseBuildingEvent
	Name     string           `json:"name"`
	Type     BuildingType     `json:"type"`
	Category BuildingCategory `json:"category"`
	OwnerID  uint64           `json:"owner_id"`
	Position *Position        `json:"position,omitempty"`
	Size     *Size            `json:"size,omitempty"`
}

// NewBuildingCreatedEvent 创建建筑创建事件
func NewBuildingCreatedEvent(buildingID, name string, buildingType BuildingType, ownerID uint64) *BuildingCreatedEvent {
	return &BuildingCreatedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeBuildingCreated,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		Name:    name,
		Type:    buildingType,
		OwnerID: ownerID,
	}
}

// BuildingUpdatedEvent 建筑更新事件
type BuildingUpdatedEvent struct {
	*BaseBuildingEvent
	Changes map[string]interface{} `json:"changes"`
}

// NewBuildingUpdatedEvent 创建建筑更新事件
func NewBuildingUpdatedEvent(buildingID string, changes map[string]interface{}) *BuildingUpdatedEvent {
	return &BuildingUpdatedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeBuildingUpdated,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		Changes: changes,
	}
}

// BuildingDeletedEvent 建筑删除事件
type BuildingDeletedEvent struct {
	*BaseBuildingEvent
	Reason string `json:"reason"`
}

// NewBuildingDeletedEvent 创建建筑删除事件
func NewBuildingDeletedEvent(buildingID, reason string) *BuildingDeletedEvent {
	return &BuildingDeletedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeBuildingDeleted,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		Reason: reason,
	}
}

// 建筑状态事件

// BuildingStatusChangedEvent 建筑状态变更事件
type BuildingStatusChangedEvent struct {
	*BaseBuildingEvent
	OldStatus BuildingStatus `json:"old_status"`
	NewStatus BuildingStatus `json:"new_status"`
	Reason    string         `json:"reason"`
}

// NewBuildingStatusChangedEvent 创建建筑状态变更事件
func NewBuildingStatusChangedEvent(buildingID string, oldStatus, newStatus BuildingStatus, reason string) *BuildingStatusChangedEvent {
	return &BuildingStatusChangedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeBuildingStatusChanged,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Reason:    reason,
	}
}

// BuildingHealthChangedEvent 建筑健康度变更事件
type BuildingHealthChangedEvent struct {
	*BaseBuildingEvent
	OldHealth float64 `json:"old_health"`
	NewHealth float64 `json:"new_health"`
	Reason    string  `json:"reason"`
}

// NewBuildingHealthChangedEvent 创建建筑健康度变更事件
func NewBuildingHealthChangedEvent(buildingID string, oldHealth, newHealth float64, reason string) *BuildingHealthChangedEvent {
	return &BuildingHealthChangedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeBuildingHealthChanged,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		OldHealth: oldHealth,
		NewHealth: newHealth,
		Reason:    reason,
	}
}

// BuildingLevelChangedEvent 建筑等级变更事件
type BuildingLevelChangedEvent struct {
	*BaseBuildingEvent
	OldLevel int32  `json:"old_level"`
	NewLevel int32  `json:"new_level"`
	Reason   string `json:"reason"`
}

// NewBuildingLevelChangedEvent 创建建筑等级变更事件
func NewBuildingLevelChangedEvent(buildingID string, oldLevel, newLevel int32, reason string) *BuildingLevelChangedEvent {
	return &BuildingLevelChangedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeBuildingLevelChanged,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		OldLevel: oldLevel,
		NewLevel: newLevel,
		Reason:   reason,
	}
}

// 建造相关事件

// ConstructionStartedEvent 建造开始事件
type ConstructionStartedEvent struct {
	*BaseBuildingEvent
	ConstructionID string          `json:"construction_id"`
	Duration       time.Duration   `json:"duration"`
	Costs          []*ResourceCost `json:"costs,omitempty"`
}

// NewConstructionStartedEvent 创建建造开始事件
func NewConstructionStartedEvent(buildingID, constructionID string, duration time.Duration) *ConstructionStartedEvent {
	return &ConstructionStartedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeConstructionStarted,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		ConstructionID: constructionID,
		Duration:       duration,
	}
}

// ConstructionProgressUpdatedEvent 建造进度更新事件
type ConstructionProgressUpdatedEvent struct {
	*BaseBuildingEvent
	ConstructionID string  `json:"construction_id"`
	Progress       float64 `json:"progress"`
	OldProgress    float64 `json:"old_progress"`
}

// NewConstructionProgressUpdatedEvent 创建建造进度更新事件
func NewConstructionProgressUpdatedEvent(buildingID, constructionID string, progress float64) *ConstructionProgressUpdatedEvent {
	return &ConstructionProgressUpdatedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeConstructionProgressUpdated,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		ConstructionID: constructionID,
		Progress:       progress,
	}
}

// ConstructionCompletedEvent 建造完成事件
type ConstructionCompletedEvent struct {
	*BaseBuildingEvent
	ConstructionID string          `json:"construction_id"`
	Duration       time.Duration   `json:"duration"`
	ActualCosts    []*ResourceCost `json:"actual_costs,omitempty"`
}

// NewConstructionCompletedEvent 创建建造完成事件
func NewConstructionCompletedEvent(buildingID, constructionID string) *ConstructionCompletedEvent {
	return &ConstructionCompletedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeConstructionCompleted,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		ConstructionID: constructionID,
	}
}

// ConstructionCancelledEvent 建造取消事件
type ConstructionCancelledEvent struct {
	*BaseBuildingEvent
	ConstructionID string `json:"construction_id"`
	Reason         string `json:"reason"`
}

// NewConstructionCancelledEvent 创建建造取消事件
func NewConstructionCancelledEvent(buildingID, constructionID, reason string) *ConstructionCancelledEvent {
	return &ConstructionCancelledEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeConstructionCancelled,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		ConstructionID: constructionID,
		Reason:         reason,
	}
}

// 升级相关事件

// UpgradeStartedEvent 升级开始事件
type UpgradeStartedEvent struct {
	*BaseBuildingEvent
	UpgradeID string          `json:"upgrade_id"`
	FromLevel int32           `json:"from_level"`
	ToLevel   int32           `json:"to_level"`
	Duration  time.Duration   `json:"duration"`
	Costs     []*ResourceCost `json:"costs,omitempty"`
}

// NewUpgradeStartedEvent 创建升级开始事件
func NewUpgradeStartedEvent(buildingID, upgradeID string, fromLevel, toLevel int32) *UpgradeStartedEvent {
	return &UpgradeStartedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeUpgradeStarted,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		UpgradeID: upgradeID,
		FromLevel: fromLevel,
		ToLevel:   toLevel,
	}
}

// UpgradeProgressUpdatedEvent 升级进度更新事件
type UpgradeProgressUpdatedEvent struct {
	*BaseBuildingEvent
	UpgradeID   string  `json:"upgrade_id"`
	Progress    float64 `json:"progress"`
	OldProgress float64 `json:"old_progress"`
}

// NewUpgradeProgressUpdatedEvent 创建升级进度更新事件
func NewUpgradeProgressUpdatedEvent(buildingID, upgradeID string, progress float64) *UpgradeProgressUpdatedEvent {
	return &UpgradeProgressUpdatedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeUpgradeProgressUpdated,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		UpgradeID: upgradeID,
		Progress:  progress,
	}
}

// UpgradeCompletedEvent 升级完成事件
type UpgradeCompletedEvent struct {
	*BaseBuildingEvent
	UpgradeID   string            `json:"upgrade_id"`
	FromLevel   int32             `json:"from_level"`
	ToLevel     int32             `json:"to_level"`
	Duration    time.Duration     `json:"duration"`
	ActualCosts []*ResourceCost   `json:"actual_costs,omitempty"`
	Benefits    []*UpgradeBenefit `json:"benefits,omitempty"`
}

// NewUpgradeCompletedEvent 创建升级完成事件
func NewUpgradeCompletedEvent(buildingID, upgradeID string, fromLevel, toLevel int32) *UpgradeCompletedEvent {
	return &UpgradeCompletedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeUpgradeCompleted,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		UpgradeID: upgradeID,
		FromLevel: fromLevel,
		ToLevel:   toLevel,
	}
}

// UpgradeCancelledEvent 升级取消事件
type UpgradeCancelledEvent struct {
	*BaseBuildingEvent
	UpgradeID string `json:"upgrade_id"`
	Reason    string `json:"reason"`
}

// NewUpgradeCancelledEvent 创建升级取消事件
func NewUpgradeCancelledEvent(buildingID, upgradeID, reason string) *UpgradeCancelledEvent {
	return &UpgradeCancelledEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeUpgradeCancelled,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		UpgradeID: upgradeID,
		Reason:    reason,
	}
}

// 维护相关事件

// BuildingRepairedEvent 建筑修复事件
type BuildingRepairedEvent struct {
	*BaseBuildingEvent
	OldHealth    float64         `json:"old_health"`
	NewHealth    float64         `json:"new_health"`
	RepairAmount float64         `json:"repair_amount"`
	Costs        []*ResourceCost `json:"costs,omitempty"`
}

// NewBuildingRepairedEvent 创建建筑修复事件
func NewBuildingRepairedEvent(buildingID string, oldHealth, newHealth float64) *BuildingRepairedEvent {
	return &BuildingRepairedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeBuildingRepaired,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		OldHealth:    oldHealth,
		NewHealth:    newHealth,
		RepairAmount: newHealth - oldHealth,
	}
}

// BuildingDamagedEvent 建筑损坏事件
type BuildingDamagedEvent struct {
	*BaseBuildingEvent
	OldHealth    float64 `json:"old_health"`
	NewHealth    float64 `json:"new_health"`
	DamageAmount float64 `json:"damage_amount"`
	DamageType   string  `json:"damage_type"`
	Reason       string  `json:"reason"`
}

// NewBuildingDamagedEvent 创建建筑损坏事件
func NewBuildingDamagedEvent(buildingID string, oldHealth, newHealth float64, damageType, reason string) *BuildingDamagedEvent {
	return &BuildingDamagedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeBuildingDamaged,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		OldHealth:    oldHealth,
		NewHealth:    newHealth,
		DamageAmount: oldHealth - newHealth,
		DamageType:   damageType,
		Reason:       reason,
	}
}

// BuildingDestroyedEvent 建筑摧毁事件
type BuildingDestroyedEvent struct {
	*BaseBuildingEvent
	Reason string `json:"reason"`
}

// NewBuildingDestroyedEvent 创建建筑摧毁事件
func NewBuildingDestroyedEvent(buildingID, reason string) *BuildingDestroyedEvent {
	return &BuildingDestroyedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeBuildingDestroyed,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		Reason: reason,
	}
}

// 工人相关事件

// WorkerAssignedEvent 工人分配事件
type WorkerAssignedEvent struct {
	*BaseBuildingEvent
	WorkerID       uint64     `json:"worker_id"`
	Role           WorkerRole `json:"role"`
	Task           string     `json:"task"`
	ConstructionID string     `json:"construction_id,omitempty"`
	UpgradeID      string     `json:"upgrade_id,omitempty"`
}

// NewWorkerAssignedEvent 创建工人分配事件
func NewWorkerAssignedEvent(buildingID string, workerID uint64, role WorkerRole, task string) *WorkerAssignedEvent {
	return &WorkerAssignedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeWorkerAssigned,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		WorkerID: workerID,
		Role:     role,
		Task:     task,
	}
}

// WorkerUnassignedEvent 工人取消分配事件
type WorkerUnassignedEvent struct {
	*BaseBuildingEvent
	WorkerID       uint64 `json:"worker_id"`
	Reason         string `json:"reason"`
	ConstructionID string `json:"construction_id,omitempty"`
	UpgradeID      string `json:"upgrade_id,omitempty"`
}

// NewWorkerUnassignedEvent 创建工人取消分配事件
func NewWorkerUnassignedEvent(buildingID string, workerID uint64, reason string) *WorkerUnassignedEvent {
	return &WorkerUnassignedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeWorkerUnassigned,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		WorkerID: workerID,
		Reason:   reason,
	}
}

// 蓝图相关事件

// BlueprintCreatedEvent 蓝图创建事件
type BlueprintCreatedEvent struct {
	*BaseBuildingEvent
	BlueprintID string           `json:"blueprint_id"`
	Name        string           `json:"name"`
	Category    BuildingCategory `json:"category"`
	Author      string           `json:"author"`
	Difficulty  int32            `json:"difficulty"`
}

// NewBlueprintCreatedEvent 创建蓝图创建事件
func NewBlueprintCreatedEvent(blueprintID, name string, category BuildingCategory) *BlueprintCreatedEvent {
	return &BlueprintCreatedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeBlueprintCreated,
			BuildingID: "", // 蓝图事件可能没有关联的建筑ID
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		BlueprintID: blueprintID,
		Name:        name,
		Category:    category,
	}
}

// BlueprintUsedEvent 蓝图使用事件
type BlueprintUsedEvent struct {
	*BaseBuildingEvent
	BlueprintID string `json:"blueprint_id"`
	UserID      uint64 `json:"user_id"`
}

// NewBlueprintUsedEvent 创建蓝图使用事件
func NewBlueprintUsedEvent(buildingID, blueprintID string, userID uint64) *BlueprintUsedEvent {
	return &BlueprintUsedEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeBlueprintUsed,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		BlueprintID: blueprintID,
		UserID:      userID,
	}
}

// 系统事件

// BuildingSystemErrorEvent 建筑系统错误事件
type BuildingSystemErrorEvent struct {
	*BaseBuildingEvent
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Severity     string `json:"severity"`
	StackTrace   string `json:"stack_trace,omitempty"`
}

// NewBuildingSystemErrorEvent 创建建筑系统错误事件
func NewBuildingSystemErrorEvent(buildingID, errorCode, errorMessage, severity string) *BuildingSystemErrorEvent {
	return &BuildingSystemErrorEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeBuildingSystemError,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
		Severity:     severity,
	}
}

// BuildingMaintenanceScheduledEvent 建筑维护计划事件
type BuildingMaintenanceScheduledEvent struct {
	*BaseBuildingEvent
	MaintenanceType string        `json:"maintenance_type"`
	ScheduledAt     time.Time     `json:"scheduled_at"`
	Duration        time.Duration `json:"duration"`
	Description     string        `json:"description"`
}

// NewBuildingMaintenanceScheduledEvent 创建建筑维护计划事件
func NewBuildingMaintenanceScheduledEvent(buildingID, maintenanceType string, scheduledAt time.Time, duration time.Duration) *BuildingMaintenanceScheduledEvent {
	return &BuildingMaintenanceScheduledEvent{
		BaseBuildingEvent: &BaseBuildingEvent{
			EventID:    generateEventID(),
			EventType:  EventTypeBuildingMaintenanceScheduled,
			BuildingID: buildingID,
			Timestamp:  time.Now(),
			Payload:    make(map[string]interface{}),
			Metadata:   make(map[string]interface{}),
		},
		MaintenanceType: maintenanceType,
		ScheduledAt:     scheduledAt,
		Duration:        duration,
	}
}

// 事件常量

const (
	// 建筑生命周期事件
	EventTypeBuildingCreated = "building.created"
	EventTypeBuildingUpdated = "building.updated"
	EventTypeBuildingDeleted = "building.deleted"

	// 建筑状态事件
	EventTypeBuildingStatusChanged = "building.status_changed"
	EventTypeBuildingHealthChanged = "building.health_changed"
	EventTypeBuildingLevelChanged  = "building.level_changed"

	// 建造相关事件
	EventTypeConstructionStarted         = "construction.started"
	EventTypeConstructionProgressUpdated = "construction.progress_updated"
	EventTypeConstructionCompleted       = "construction.completed"
	EventTypeConstructionCancelled       = "construction.cancelled"

	// 升级相关事件
	EventTypeUpgradeStarted         = "upgrade.started"
	EventTypeUpgradeProgressUpdated = "upgrade.progress_updated"
	EventTypeUpgradeCompleted       = "upgrade.completed"
	EventTypeUpgradeCancelled       = "upgrade.cancelled"

	// 维护相关事件
	EventTypeBuildingRepaired  = "building.repaired"
	EventTypeBuildingDamaged   = "building.damaged"
	EventTypeBuildingDestroyed = "building.destroyed"

	// 工人相关事件
	EventTypeWorkerAssigned   = "worker.assigned"
	EventTypeWorkerUnassigned = "worker.unassigned"

	// 蓝图相关事件
	EventTypeBlueprintCreated = "blueprint.created"
	EventTypeBlueprintUsed    = "blueprint.used"

	// 系统事件
	EventTypeBuildingSystemError          = "building.system_error"
	EventTypeBuildingMaintenanceScheduled = "building.maintenance_scheduled"
)

// 事件处理器接口

// BuildingEventHandler 建筑事件处理器接口
type BuildingEventHandler interface {
	Handle(ctx context.Context, event BuildingEvent) error
	CanHandle(eventType string) bool
	GetHandlerName() string
}

// BuildingEventBus 建筑事件总线接口
type BuildingEventBus interface {
	Publish(ctx context.Context, event BuildingEvent) error
	Subscribe(eventType string, handler BuildingEventHandler) error
	Unsubscribe(eventType string, handlerName string) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// 事件聚合器

// BuildingEventAggregator 建筑事件聚合器
type BuildingEventAggregator struct {
	buildingID string
	events     []BuildingEvent
	version    int64
	createdAt  time.Time
	updatedAt  time.Time
}

// NewBuildingEventAggregator 创建新建筑事件聚合器
func NewBuildingEventAggregator(buildingID string) *BuildingEventAggregator {
	now := time.Now()
	return &BuildingEventAggregator{
		buildingID: buildingID,
		events:     make([]BuildingEvent, 0),
		version:    0,
		createdAt:  now,
		updatedAt:  now,
	}
}

// AddEvent 添加事件
func (ea *BuildingEventAggregator) AddEvent(event BuildingEvent) {
	ea.events = append(ea.events, event)
	ea.version++
	ea.updatedAt = time.Now()
}

// GetEvents 获取所有事件
func (ea *BuildingEventAggregator) GetEvents() []BuildingEvent {
	return ea.events
}

// GetEventsByType 根据类型获取事件
func (ea *BuildingEventAggregator) GetEventsByType(eventType string) []BuildingEvent {
	var result []BuildingEvent
	for _, event := range ea.events {
		if event.GetEventType() == eventType {
			result = append(result, event)
		}
	}
	return result
}

// GetEventsAfter 获取指定时间后的事件
func (ea *BuildingEventAggregator) GetEventsAfter(timestamp time.Time) []BuildingEvent {
	var result []BuildingEvent
	for _, event := range ea.events {
		if event.GetTimestamp().After(timestamp) {
			result = append(result, event)
		}
	}
	return result
}

// GetEventCount 获取事件数量
func (ea *BuildingEventAggregator) GetEventCount() int {
	return len(ea.events)
}

// GetEventCountByType 根据类型获取事件数量
func (ea *BuildingEventAggregator) GetEventCountByType(eventType string) int {
	count := 0
	for _, event := range ea.events {
		if event.GetEventType() == eventType {
			count++
		}
	}
	return count
}

// Clear 清空事件
func (ea *BuildingEventAggregator) Clear() {
	ea.events = make([]BuildingEvent, 0)
	ea.version = 0
	ea.updatedAt = time.Now()
}

// GetBuildingID 获取建筑ID
func (ea *BuildingEventAggregator) GetBuildingID() string {
	return ea.buildingID
}

// GetVersion 获取版本
func (ea *BuildingEventAggregator) GetVersion() int64 {
	return ea.version
}

// GetCreatedAt 获取创建时间
func (ea *BuildingEventAggregator) GetCreatedAt() time.Time {
	return ea.createdAt
}

// GetUpdatedAt 获取更新时间
func (ea *BuildingEventAggregator) GetUpdatedAt() time.Time {
	return ea.updatedAt
}

// 事件统计

// BuildingEventStatistics 建筑事件统计
type BuildingEventStatistics struct {
	BuildingID    string           `json:"building_id"`
	TotalEvents   int64            `json:"total_events"`
	EventsByType  map[string]int64 `json:"events_by_type"`
	EventsByHour  map[string]int64 `json:"events_by_hour"`
	EventsByDay   map[string]int64 `json:"events_by_day"`
	LastEventTime time.Time        `json:"last_event_time"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

// NewBuildingEventStatistics 创建新建筑事件统计
func NewBuildingEventStatistics(buildingID string) *BuildingEventStatistics {
	return &BuildingEventStatistics{
		BuildingID:   buildingID,
		TotalEvents:  0,
		EventsByType: make(map[string]int64),
		EventsByHour: make(map[string]int64),
		EventsByDay:  make(map[string]int64),
		UpdatedAt:    time.Now(),
	}
}

// AddEvent 添加事件到统计
func (es *BuildingEventStatistics) AddEvent(event BuildingEvent) {
	es.TotalEvents++
	es.EventsByType[event.GetEventType()]++

	timestamp := event.GetTimestamp()
	hourKey := timestamp.Format("2006-01-02-15")
	dayKey := timestamp.Format("2006-01-02")

	es.EventsByHour[hourKey]++
	es.EventsByDay[dayKey]++

	if timestamp.After(es.LastEventTime) {
		es.LastEventTime = timestamp
	}

	es.UpdatedAt = time.Now()
}

// GetMostFrequentEventType 获取最频繁的事件类型
func (es *BuildingEventStatistics) GetMostFrequentEventType() string {
	maxCount := int64(0)
	mostFrequent := ""

	for eventType, count := range es.EventsByType {
		if count > maxCount {
			maxCount = count
			mostFrequent = eventType
		}
	}

	return mostFrequent
}

// GetEventsInLastHour 获取最近一小时的事件数量
func (es *BuildingEventStatistics) GetEventsInLastHour() int64 {
	hourKey := time.Now().Format("2006-01-02-15")
	return es.EventsByHour[hourKey]
}

// GetEventsInLastDay 获取最近一天的事件数量
func (es *BuildingEventStatistics) GetEventsInLastDay() int64 {
	dayKey := time.Now().Format("2006-01-02")
	return es.EventsByDay[dayKey]
}

// 辅助函数

// generateEventID 生成事件ID
func generateEventID() string {
	return fmt.Sprintf("event_%d", time.Now().UnixNano())
}

// ValidateEventType 验证事件类型
func ValidateEventType(eventType string) bool {
	validTypes := []string{
		EventTypeBuildingCreated,
		EventTypeBuildingUpdated,
		EventTypeBuildingDeleted,
		EventTypeBuildingStatusChanged,
		EventTypeBuildingHealthChanged,
		EventTypeBuildingLevelChanged,
		EventTypeConstructionStarted,
		EventTypeConstructionProgressUpdated,
		EventTypeConstructionCompleted,
		EventTypeConstructionCancelled,
		EventTypeUpgradeStarted,
		EventTypeUpgradeProgressUpdated,
		EventTypeUpgradeCompleted,
		EventTypeUpgradeCancelled,
		EventTypeBuildingRepaired,
		EventTypeBuildingDamaged,
		EventTypeBuildingDestroyed,
		EventTypeWorkerAssigned,
		EventTypeWorkerUnassigned,
		EventTypeBlueprintCreated,
		EventTypeBlueprintUsed,
		EventTypeBuildingSystemError,
		EventTypeBuildingMaintenanceScheduled,
	}

	for _, validType := range validTypes {
		if eventType == validType {
			return true
		}
	}
	return false
}

// GetEventCategory 获取事件分类
func GetEventCategory(eventType string) string {
	switch eventType {
	case EventTypeBuildingCreated, EventTypeBuildingUpdated, EventTypeBuildingDeleted:
		return "lifecycle"
	case EventTypeBuildingStatusChanged, EventTypeBuildingHealthChanged, EventTypeBuildingLevelChanged:
		return "status"
	case EventTypeConstructionStarted, EventTypeConstructionProgressUpdated, EventTypeConstructionCompleted, EventTypeConstructionCancelled:
		return "construction"
	case EventTypeUpgradeStarted, EventTypeUpgradeProgressUpdated, EventTypeUpgradeCompleted, EventTypeUpgradeCancelled:
		return "upgrade"
	case EventTypeBuildingRepaired, EventTypeBuildingDamaged, EventTypeBuildingDestroyed:
		return "maintenance"
	case EventTypeWorkerAssigned, EventTypeWorkerUnassigned:
		return "worker"
	case EventTypeBlueprintCreated, EventTypeBlueprintUsed:
		return "blueprint"
	case EventTypeBuildingSystemError, EventTypeBuildingMaintenanceScheduled:
		return "system"
	default:
		return "unknown"
	}
}

// IsSystemEvent 检查是否为系统事件
func IsSystemEvent(eventType string) bool {
	return GetEventCategory(eventType) == "system"
}

// IsUserEvent 检查是否为用户事件
func IsUserEvent(eventType string) bool {
	category := GetEventCategory(eventType)
	return category != "system" && category != "unknown"
}

// GetEventPriority 获取事件优先级
func GetEventPriority(eventType string) int {
	switch eventType {
	case EventTypeBuildingSystemError:
		return 1 // 最高优先级
	case EventTypeBuildingDestroyed, EventTypeBuildingDamaged:
		return 2 // 高优先级
	case EventTypeConstructionCompleted, EventTypeUpgradeCompleted:
		return 3 // 中等优先级
	case EventTypeBuildingCreated, EventTypeConstructionStarted, EventTypeUpgradeStarted:
		return 4 // 普通优先级
	default:
		return 5 // 低优先级
	}
}
