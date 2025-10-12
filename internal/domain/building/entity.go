package building

import (
	"fmt"
	"time"
)

// ConstructionInfo 建造信息实体
type ConstructionInfo struct {
	ID           string                 `json:"id" bson:"_id"`
	BuildingID   string                 `json:"building_id" bson:"building_id"`
	StartedAt    time.Time              `json:"started_at" bson:"started_at"`
	Duration     time.Duration          `json:"duration" bson:"duration"`
	CompletedAt  *time.Time             `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
	Progress     float64                `json:"progress" bson:"progress"`
	Costs        []*ResourceCost        `json:"costs" bson:"costs"`
	Workers      []*WorkerAssignment    `json:"workers" bson:"workers"`
	Materials    []*MaterialUsage       `json:"materials" bson:"materials"`
	Status       ConstructionStatus     `json:"status" bson:"status"`
	Blueprint    *Blueprint             `json:"blueprint,omitempty" bson:"blueprint,omitempty"`
	Phases       []*ConstructionPhase   `json:"phases" bson:"phases"`
	CurrentPhase *ConstructionPhase     `json:"current_phase,omitempty" bson:"current_phase,omitempty"`
	QualityScore float64                `json:"quality_score" bson:"quality_score"`
	SafetyScore  float64                `json:"safety_score" bson:"safety_score"`
	Metadata     map[string]interface{} `json:"metadata" bson:"metadata"`
	CreatedAt    time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" bson:"updated_at"`
}

// ConstructionStatus 建造状态
type ConstructionStatus int32

const (
	ConstructionStatusPlanning   ConstructionStatus = iota + 1 // 规划中
	ConstructionStatusInProgress                               // 进行中
	ConstructionStatusPaused                                   // 暂停
	ConstructionStatusCompleted                                // 已完成
	ConstructionStatusCancelled                                // 已取消
	ConstructionStatusFailed                                   // 失败
)

// String 返回建造状态的字符串表示
func (cs ConstructionStatus) String() string {
	switch cs {
	case ConstructionStatusPlanning:
		return "planning"
	case ConstructionStatusInProgress:
		return "in_progress"
	case ConstructionStatusPaused:
		return "paused"
	case ConstructionStatusCompleted:
		return "completed"
	case ConstructionStatusCancelled:
		return "cancelled"
	case ConstructionStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// IsValid 检查建造状态是否有效
func (cs ConstructionStatus) IsValid() bool {
	return cs >= ConstructionStatusPlanning && cs <= ConstructionStatusFailed
}

// NewConstructionInfo 创建新建造信息
func NewConstructionInfo(buildingID string, duration time.Duration) *ConstructionInfo {
	now := time.Now()
	return &ConstructionInfo{
		ID:           generateConstructionID(),
		BuildingID:   buildingID,
		StartedAt:    now,
		Duration:     duration,
		Progress:     0.0,
		Costs:        make([]*ResourceCost, 0),
		Workers:      make([]*WorkerAssignment, 0),
		Materials:    make([]*MaterialUsage, 0),
		Status:       ConstructionStatusPlanning,
		Phases:       make([]*ConstructionPhase, 0),
		QualityScore: 100.0,
		SafetyScore:  100.0,
		Metadata:     make(map[string]interface{}),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// AddWorker 添加工人
func (ci *ConstructionInfo) AddWorker(assignment *WorkerAssignment) error {
	if assignment == nil {
		return fmt.Errorf("worker assignment cannot be nil")
	}

	// 检查工人是否已分配
	for _, existing := range ci.Workers {
		if existing.WorkerID == assignment.WorkerID {
			return fmt.Errorf("worker %d is already assigned", assignment.WorkerID)
		}
	}

	ci.Workers = append(ci.Workers, assignment)
	ci.UpdatedAt = time.Now()
	return nil
}

// RemoveWorker 移除工人
func (ci *ConstructionInfo) RemoveWorker(workerID uint64) error {
	for i, worker := range ci.Workers {
		if worker.WorkerID == workerID {
			ci.Workers = append(ci.Workers[:i], ci.Workers[i+1:]...)
			ci.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("worker %d not found", workerID)
}

// AddMaterial 添加材料
func (ci *ConstructionInfo) AddMaterial(material *MaterialUsage) error {
	if material == nil {
		return fmt.Errorf("material usage cannot be nil")
	}

	ci.Materials = append(ci.Materials, material)
	ci.UpdatedAt = time.Now()
	return nil
}

// AddPhase 添加阶段
func (ci *ConstructionInfo) AddPhase(phase *ConstructionPhase) error {
	if phase == nil {
		return fmt.Errorf("construction phase cannot be nil")
	}

	ci.Phases = append(ci.Phases, phase)
	ci.UpdatedAt = time.Now()
	return nil
}

// StartNextPhase 开始下一阶段
func (ci *ConstructionInfo) StartNextPhase() *ConstructionPhase {
	for _, phase := range ci.Phases {
		if phase.Status == PhaseStatusPending {
			phase.Start()
			ci.CurrentPhase = phase
			ci.UpdatedAt = time.Now()
			return phase
		}
	}
	return nil
}

// CompleteCurrentPhase 完成当前阶段
func (ci *ConstructionInfo) CompleteCurrentPhase() error {
	if ci.CurrentPhase == nil {
		return fmt.Errorf("no current phase to complete")
	}

	ci.CurrentPhase.Complete()
	ci.CurrentPhase = nil
	ci.UpdatedAt = time.Now()

	// 检查是否所有阶段都完成
	allCompleted := true
	for _, phase := range ci.Phases {
		if phase.Status != PhaseStatusCompleted {
			allCompleted = false
			break
		}
	}

	if allCompleted {
		ci.Status = ConstructionStatusCompleted
		now := time.Now()
		ci.CompletedAt = &now
		ci.Progress = 100.0
	}

	return nil
}

// UpdateProgress 更新进度
func (ci *ConstructionInfo) UpdateProgress(progress float64) error {
	if progress < 0 || progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100")
	}

	ci.Progress = progress
	ci.UpdatedAt = time.Now()

	if progress >= 100 {
		ci.Status = ConstructionStatusCompleted
		now := time.Now()
		ci.CompletedAt = &now
	}

	return nil
}

// SetMetadata 设置元数据
func (ci *ConstructionInfo) SetMetadata(key string, value interface{}) {
	if ci.Metadata == nil {
		ci.Metadata = make(map[string]interface{})
	}
	ci.Metadata[key] = value
	ci.UpdatedAt = time.Now()
}

// GetMetadata 获取元数据
func (ci *ConstructionInfo) GetMetadata(key string) (interface{}, bool) {
	if ci.Metadata == nil {
		return nil, false
	}
	value, exists := ci.Metadata[key]
	return value, exists
}

// GetEstimatedCompletionTime 获取预计完成时间
func (ci *ConstructionInfo) GetEstimatedCompletionTime() time.Time {
	if ci.Progress <= 0 {
		return ci.StartedAt.Add(ci.Duration)
	}

	elapsed := time.Since(ci.StartedAt)
	estimatedTotal := time.Duration(float64(elapsed) / (ci.Progress / 100.0))
	return ci.StartedAt.Add(estimatedTotal)
}

// GetEfficiency 获取建造效率
func (ci *ConstructionInfo) GetEfficiency() float64 {
	if ci.Progress <= 0 {
		return 0.0
	}

	elapsed := time.Since(ci.StartedAt)
	expectedProgress := float64(elapsed) / float64(ci.Duration) * 100.0

	if expectedProgress <= 0 {
		return 0.0
	}

	return ci.Progress / expectedProgress
}

// UpgradeInfo 升级信息实体
type UpgradeInfo struct {
	ID           string                 `json:"id" bson:"_id"`
	BuildingID   string                 `json:"building_id" bson:"building_id"`
	FromLevel    int32                  `json:"from_level" bson:"from_level"`
	ToLevel      int32                  `json:"to_level" bson:"to_level"`
	StartedAt    time.Time              `json:"started_at" bson:"started_at"`
	Duration     time.Duration          `json:"duration" bson:"duration"`
	CompletedAt  *time.Time             `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
	Progress     float64                `json:"progress" bson:"progress"`
	Costs        []*ResourceCost        `json:"costs" bson:"costs"`
	Workers      []*WorkerAssignment    `json:"workers" bson:"workers"`
	Materials    []*MaterialUsage       `json:"materials" bson:"materials"`
	Status       UpgradeStatus          `json:"status" bson:"status"`
	Benefits     []*UpgradeBenefit      `json:"benefits" bson:"benefits"`
	Requirements []*Requirement         `json:"requirements" bson:"requirements"`
	Metadata     map[string]interface{} `json:"metadata" bson:"metadata"`
	CreatedAt    time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" bson:"updated_at"`
}

// UpgradeStatus 升级状态
type UpgradeStatus int32

const (
	UpgradeStatusPlanning   UpgradeStatus = iota + 1 // 规划中
	UpgradeStatusInProgress                          // 进行中
	UpgradeStatusPaused                              // 暂停
	UpgradeStatusCompleted                           // 已完成
	UpgradeStatusCancelled                           // 已取消
	UpgradeStatusFailed                              // 失败
)

// String 返回升级状态的字符串表示
func (us UpgradeStatus) String() string {
	switch us {
	case UpgradeStatusPlanning:
		return "planning"
	case UpgradeStatusInProgress:
		return "in_progress"
	case UpgradeStatusPaused:
		return "paused"
	case UpgradeStatusCompleted:
		return "completed"
	case UpgradeStatusCancelled:
		return "cancelled"
	case UpgradeStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// IsValid 检查升级状态是否有效
func (us UpgradeStatus) IsValid() bool {
	return us >= UpgradeStatusPlanning && us <= UpgradeStatusFailed
}

// NewUpgradeInfo 创建新升级信息
func NewUpgradeInfo(buildingID string, fromLevel, toLevel int32, duration time.Duration) *UpgradeInfo {
	now := time.Now()
	return &UpgradeInfo{
		ID:           generateUpgradeID(),
		BuildingID:   buildingID,
		FromLevel:    fromLevel,
		ToLevel:      toLevel,
		StartedAt:    now,
		Duration:     duration,
		Progress:     0.0,
		Costs:        make([]*ResourceCost, 0),
		Workers:      make([]*WorkerAssignment, 0),
		Materials:    make([]*MaterialUsage, 0),
		Status:       UpgradeStatusPlanning,
		Benefits:     make([]*UpgradeBenefit, 0),
		Requirements: make([]*Requirement, 0),
		Metadata:     make(map[string]interface{}),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// AddBenefit 添加升级收益
func (ui *UpgradeInfo) AddBenefit(benefit *UpgradeBenefit) error {
	if benefit == nil {
		return fmt.Errorf("upgrade benefit cannot be nil")
	}

	ui.Benefits = append(ui.Benefits, benefit)
	ui.UpdatedAt = time.Now()
	return nil
}

// AddRequirement 添加升级要求
func (ui *UpgradeInfo) AddRequirement(requirement *Requirement) error {
	if requirement == nil {
		return fmt.Errorf("requirement cannot be nil")
	}

	ui.Requirements = append(ui.Requirements, requirement)
	ui.UpdatedAt = time.Now()
	return nil
}

// CheckRequirements 检查升级要求
func (ui *UpgradeInfo) CheckRequirements() bool {
	for _, req := range ui.Requirements {
		if !req.IsMet() {
			return false
		}
	}
	return true
}

// UpdateProgress 更新升级进度
func (ui *UpgradeInfo) UpdateProgress(progress float64) error {
	if progress < 0 || progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100")
	}

	ui.Progress = progress
	ui.UpdatedAt = time.Now()

	if progress >= 100 {
		ui.Status = UpgradeStatusCompleted
		now := time.Now()
		ui.CompletedAt = &now
	}

	return nil
}

// SetMetadata 设置元数据
func (ui *UpgradeInfo) SetMetadata(key string, value interface{}) {
	if ui.Metadata == nil {
		ui.Metadata = make(map[string]interface{})
	}
	ui.Metadata[key] = value
	ui.UpdatedAt = time.Now()
}

// GetMetadata 获取元数据
func (ui *UpgradeInfo) GetMetadata(key string) (interface{}, bool) {
	if ui.Metadata == nil {
		return nil, false
	}
	value, exists := ui.Metadata[key]
	return value, exists
}

// WorkerAssignment 工人分配实体
type WorkerAssignment struct {
	ID         string                 `json:"id" bson:"_id"`
	WorkerID   uint64                 `json:"worker_id" bson:"worker_id"`
	Role       WorkerRole             `json:"role" bson:"role"`
	Task       string                 `json:"task" bson:"task"`
	Efficiency float64                `json:"efficiency" bson:"efficiency"`
	StartTime  time.Time              `json:"start_time" bson:"start_time"`
	EndTime    *time.Time             `json:"end_time,omitempty" bson:"end_time,omitempty"`
	Status     WorkerAssignmentStatus `json:"status" bson:"status"`
	Progress   float64                `json:"progress" bson:"progress"`
	Metadata   map[string]interface{} `json:"metadata" bson:"metadata"`
	CreatedAt  time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at" bson:"updated_at"`
}

// WorkerAssignmentStatus 工人分配状态
type WorkerAssignmentStatus int32

const (
	WorkerAssignmentStatusAssigned  WorkerAssignmentStatus = iota + 1 // 已分配
	WorkerAssignmentStatusWorking                                     // 工作中
	WorkerAssignmentStatusPaused                                      // 暂停
	WorkerAssignmentStatusCompleted                                   // 已完成
	WorkerAssignmentStatusCancelled                                   // 已取消
)

// String 返回工人分配状态的字符串表示
func (was WorkerAssignmentStatus) String() string {
	switch was {
	case WorkerAssignmentStatusAssigned:
		return "assigned"
	case WorkerAssignmentStatusWorking:
		return "working"
	case WorkerAssignmentStatusPaused:
		return "paused"
	case WorkerAssignmentStatusCompleted:
		return "completed"
	case WorkerAssignmentStatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

// IsValid 检查工人分配状态是否有效
func (was WorkerAssignmentStatus) IsValid() bool {
	return was >= WorkerAssignmentStatusAssigned && was <= WorkerAssignmentStatusCancelled
}

// NewWorkerAssignment 创建新工人分配
func NewWorkerAssignment(workerID uint64, role WorkerRole, task string) *WorkerAssignment {
	now := time.Now()
	return &WorkerAssignment{
		ID:         generateWorkerAssignmentID(),
		WorkerID:   workerID,
		Role:       role,
		Task:       task,
		Efficiency: 1.0,
		StartTime:  now,
		Status:     WorkerAssignmentStatusAssigned,
		Progress:   0.0,
		Metadata:   make(map[string]interface{}),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// Start 开始工作
func (wa *WorkerAssignment) Start() {
	wa.Status = WorkerAssignmentStatusWorking
	wa.UpdatedAt = time.Now()
}

// Pause 暂停工作
func (wa *WorkerAssignment) Pause() {
	wa.Status = WorkerAssignmentStatusPaused
	wa.UpdatedAt = time.Now()
}

// Resume 恢复工作
func (wa *WorkerAssignment) Resume() {
	wa.Status = WorkerAssignmentStatusWorking
	wa.UpdatedAt = time.Now()
}

// Complete 完成工作
func (wa *WorkerAssignment) Complete() {
	now := time.Now()
	wa.Status = WorkerAssignmentStatusCompleted
	wa.EndTime = &now
	wa.Progress = 100.0
	wa.UpdatedAt = now
}

// Cancel 取消工作
func (wa *WorkerAssignment) Cancel() {
	now := time.Now()
	wa.Status = WorkerAssignmentStatusCancelled
	wa.EndTime = &now
	wa.UpdatedAt = now
}

// UpdateProgress 更新进度
func (wa *WorkerAssignment) UpdateProgress(progress float64) error {
	if progress < 0 || progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100")
	}

	wa.Progress = progress
	wa.UpdatedAt = time.Now()

	if progress >= 100 {
		wa.Complete()
	}

	return nil
}

// GetDuration 获取工作持续时间
func (wa *WorkerAssignment) GetDuration() time.Duration {
	if wa.EndTime != nil {
		return wa.EndTime.Sub(wa.StartTime)
	}
	return time.Since(wa.StartTime)
}

// MaterialUsage 材料使用实体
type MaterialUsage struct {
	ID           string                 `json:"id" bson:"_id"`
	MaterialType string                 `json:"material_type" bson:"material_type"`
	Quantity     int64                  `json:"quantity" bson:"quantity"`
	Used         int64                  `json:"used" bson:"used"`
	Wasted       int64                  `json:"wasted" bson:"wasted"`
	Quality      float64                `json:"quality" bson:"quality"`
	Cost         int64                  `json:"cost" bson:"cost"`
	Supplier     string                 `json:"supplier" bson:"supplier"`
	DeliveredAt  time.Time              `json:"delivered_at" bson:"delivered_at"`
	UsedAt       *time.Time             `json:"used_at,omitempty" bson:"used_at,omitempty"`
	Status       MaterialUsageStatus    `json:"status" bson:"status"`
	Metadata     map[string]interface{} `json:"metadata" bson:"metadata"`
	CreatedAt    time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" bson:"updated_at"`
}

// MaterialUsageStatus 材料使用状态
type MaterialUsageStatus int32

const (
	MaterialUsageStatusOrdered   MaterialUsageStatus = iota + 1 // 已订购
	MaterialUsageStatusDelivered                                // 已交付
	MaterialUsageStatusInUse                                    // 使用中
	MaterialUsageStatusUsed                                     // 已使用
	MaterialUsageStatusWasted                                   // 已浪费
	MaterialUsageStatusReturned                                 // 已退回
)

// String 返回材料使用状态的字符串表示
func (mus MaterialUsageStatus) String() string {
	switch mus {
	case MaterialUsageStatusOrdered:
		return "ordered"
	case MaterialUsageStatusDelivered:
		return "delivered"
	case MaterialUsageStatusInUse:
		return "in_use"
	case MaterialUsageStatusUsed:
		return "used"
	case MaterialUsageStatusWasted:
		return "wasted"
	case MaterialUsageStatusReturned:
		return "returned"
	default:
		return "unknown"
	}
}

// IsValid 检查材料使用状态是否有效
func (mus MaterialUsageStatus) IsValid() bool {
	return mus >= MaterialUsageStatusOrdered && mus <= MaterialUsageStatusReturned
}

// NewMaterialUsage 创建新材料使用
func NewMaterialUsage(materialType string, quantity int64, cost int64) *MaterialUsage {
	now := time.Now()
	return &MaterialUsage{
		ID:           generateMaterialUsageID(),
		MaterialType: materialType,
		Quantity:     quantity,
		Used:         0,
		Wasted:       0,
		Quality:      1.0,
		Cost:         cost,
		Supplier:     "",
		DeliveredAt:  now,
		Status:       MaterialUsageStatusOrdered,
		Metadata:     make(map[string]interface{}),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// Use 使用材料
func (mu *MaterialUsage) Use(amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("use amount must be positive")
	}

	if mu.Used+amount > mu.Quantity {
		return fmt.Errorf("insufficient material: have %d, need %d", mu.Quantity-mu.Used, amount)
	}

	mu.Used += amount
	mu.Status = MaterialUsageStatusInUse

	if mu.Used >= mu.Quantity {
		mu.Status = MaterialUsageStatusUsed
		now := time.Now()
		mu.UsedAt = &now
	}

	mu.UpdatedAt = time.Now()
	return nil
}

// Waste 浪费材料
func (mu *MaterialUsage) Waste(amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("waste amount must be positive")
	}

	if mu.Used+mu.Wasted+amount > mu.Quantity {
		return fmt.Errorf("waste amount exceeds available material")
	}

	mu.Wasted += amount
	mu.UpdatedAt = time.Now()
	return nil
}

// GetRemaining 获取剩余材料
func (mu *MaterialUsage) GetRemaining() int64 {
	return mu.Quantity - mu.Used - mu.Wasted
}

// GetUsageRate 获取使用率
func (mu *MaterialUsage) GetUsageRate() float64 {
	if mu.Quantity == 0 {
		return 0.0
	}
	return float64(mu.Used) / float64(mu.Quantity) * 100.0
}

// GetWasteRate 获取浪费率
func (mu *MaterialUsage) GetWasteRate() float64 {
	if mu.Quantity == 0 {
		return 0.0
	}
	return float64(mu.Wasted) / float64(mu.Quantity) * 100.0
}

// ConstructionPhase 建造阶段实体
type ConstructionPhase struct {
	ID           string                 `json:"id" bson:"_id"`
	Name         string                 `json:"name" bson:"name"`
	Description  string                 `json:"description" bson:"description"`
	Order        int32                  `json:"order" bson:"order"`
	Duration     time.Duration          `json:"duration" bson:"duration"`
	StartedAt    *time.Time             `json:"started_at,omitempty" bson:"started_at,omitempty"`
	CompletedAt  *time.Time             `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
	Progress     float64                `json:"progress" bson:"progress"`
	Status       PhaseStatus            `json:"status" bson:"status"`
	Requirements []*Requirement         `json:"requirements" bson:"requirements"`
	Tasks        []*PhaseTask           `json:"tasks" bson:"tasks"`
	Workers      []*WorkerAssignment    `json:"workers" bson:"workers"`
	Materials    []*MaterialUsage       `json:"materials" bson:"materials"`
	Metadata     map[string]interface{} `json:"metadata" bson:"metadata"`
	CreatedAt    time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" bson:"updated_at"`
}

// PhaseStatus 阶段状态
type PhaseStatus int32

const (
	PhaseStatusPending    PhaseStatus = iota + 1 // 等待中
	PhaseStatusInProgress                        // 进行中
	PhaseStatusPaused                            // 暂停
	PhaseStatusCompleted                         // 已完成
	PhaseStatusCancelled                         // 已取消
	PhaseStatusFailed                            // 失败
)

// String 返回阶段状态的字符串表示
func (ps PhaseStatus) String() string {
	switch ps {
	case PhaseStatusPending:
		return "pending"
	case PhaseStatusInProgress:
		return "in_progress"
	case PhaseStatusPaused:
		return "paused"
	case PhaseStatusCompleted:
		return "completed"
	case PhaseStatusCancelled:
		return "cancelled"
	case PhaseStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// IsValid 检查阶段状态是否有效
func (ps PhaseStatus) IsValid() bool {
	return ps >= PhaseStatusPending && ps <= PhaseStatusFailed
}

// NewConstructionPhase 创建新建造阶段
func NewConstructionPhase(name, description string, order int32, duration time.Duration) *ConstructionPhase {
	now := time.Now()
	return &ConstructionPhase{
		ID:           generatePhaseID(),
		Name:         name,
		Description:  description,
		Order:        order,
		Duration:     duration,
		Progress:     0.0,
		Status:       PhaseStatusPending,
		Requirements: make([]*Requirement, 0),
		Tasks:        make([]*PhaseTask, 0),
		Workers:      make([]*WorkerAssignment, 0),
		Materials:    make([]*MaterialUsage, 0),
		Metadata:     make(map[string]interface{}),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// Start 开始阶段
func (cp *ConstructionPhase) Start() {
	now := time.Now()
	cp.StartedAt = &now
	cp.Status = PhaseStatusInProgress
	cp.UpdatedAt = now
}

// Complete 完成阶段
func (cp *ConstructionPhase) Complete() {
	now := time.Now()
	cp.CompletedAt = &now
	cp.Progress = 100.0
	cp.Status = PhaseStatusCompleted
	cp.UpdatedAt = now
}

// Pause 暂停阶段
func (cp *ConstructionPhase) Pause() {
	cp.Status = PhaseStatusPaused
	cp.UpdatedAt = time.Now()
}

// Resume 恢复阶段
func (cp *ConstructionPhase) Resume() {
	cp.Status = PhaseStatusInProgress
	cp.UpdatedAt = time.Now()
}

// Cancel 取消阶段
func (cp *ConstructionPhase) Cancel() {
	cp.Status = PhaseStatusCancelled
	cp.UpdatedAt = time.Now()
}

// UpdateProgress 更新进度
func (cp *ConstructionPhase) UpdateProgress(progress float64) error {
	if progress < 0 || progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100")
	}

	cp.Progress = progress
	cp.UpdatedAt = time.Now()

	if progress >= 100 {
		cp.Complete()
	}

	return nil
}

// AddTask 添加任务
func (cp *ConstructionPhase) AddTask(task *PhaseTask) error {
	if task == nil {
		return fmt.Errorf("phase task cannot be nil")
	}

	cp.Tasks = append(cp.Tasks, task)
	cp.UpdatedAt = time.Now()
	return nil
}

// PhaseTask 阶段任务实体
type PhaseTask struct {
	ID          string                 `json:"id" bson:"_id"`
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	Type        TaskType               `json:"type" bson:"type"`
	Priority    TaskPriority           `json:"priority" bson:"priority"`
	Duration    time.Duration          `json:"duration" bson:"duration"`
	StartedAt   *time.Time             `json:"started_at,omitempty" bson:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
	Progress    float64                `json:"progress" bson:"progress"`
	Status      TaskStatus             `json:"status" bson:"status"`
	AssignedTo  *uint64                `json:"assigned_to,omitempty" bson:"assigned_to,omitempty"`
	DependsOn   []string               `json:"depends_on" bson:"depends_on"`
	Metadata    map[string]interface{} `json:"metadata" bson:"metadata"`
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updated_at"`
}

// TaskType 任务类型
type TaskType int32

const (
	TaskTypeFoundation  TaskType = iota + 1 // 地基
	TaskTypeFramework                       // 框架
	TaskTypeWalls                           // 墙体
	TaskTypeRoof                            // 屋顶
	TaskTypeElectrical                      // 电气
	TaskTypePlumbing                        // 管道
	TaskTypeInterior                        // 内装
	TaskTypeExterior                        // 外装
	TaskTypeLandscaping                     // 景观
	TaskTypeCustom                          // 自定义
)

// String 返回任务类型的字符串表示
func (tt TaskType) String() string {
	switch tt {
	case TaskTypeFoundation:
		return "foundation"
	case TaskTypeFramework:
		return "framework"
	case TaskTypeWalls:
		return "walls"
	case TaskTypeRoof:
		return "roof"
	case TaskTypeElectrical:
		return "electrical"
	case TaskTypePlumbing:
		return "plumbing"
	case TaskTypeInterior:
		return "interior"
	case TaskTypeExterior:
		return "exterior"
	case TaskTypeLandscaping:
		return "landscaping"
	case TaskTypeCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// IsValid 检查任务类型是否有效
func (tt TaskType) IsValid() bool {
	return tt >= TaskTypeFoundation && tt <= TaskTypeCustom
}

// TaskPriority 任务优先级
type TaskPriority int32

const (
	TaskPriorityLow      TaskPriority = iota + 1 // 低优先级
	TaskPriorityNormal                           // 普通优先级
	TaskPriorityHigh                             // 高优先级
	TaskPriorityCritical                         // 关键优先级
)

// String 返回任务优先级的字符串表示
func (tp TaskPriority) String() string {
	switch tp {
	case TaskPriorityLow:
		return "low"
	case TaskPriorityNormal:
		return "normal"
	case TaskPriorityHigh:
		return "high"
	case TaskPriorityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// IsValid 检查任务优先级是否有效
func (tp TaskPriority) IsValid() bool {
	return tp >= TaskPriorityLow && tp <= TaskPriorityCritical
}

// TaskStatus 任务状态
type TaskStatus int32

const (
	TaskStatusPending    TaskStatus = iota + 1 // 等待中
	TaskStatusInProgress                       // 进行中
	TaskStatusPaused                           // 暂停
	TaskStatusCompleted                        // 已完成
	TaskStatusCancelled                        // 已取消
	TaskStatusFailed                           // 失败
)

// String 返回任务状态的字符串表示
func (ts TaskStatus) String() string {
	switch ts {
	case TaskStatusPending:
		return "pending"
	case TaskStatusInProgress:
		return "in_progress"
	case TaskStatusPaused:
		return "paused"
	case TaskStatusCompleted:
		return "completed"
	case TaskStatusCancelled:
		return "cancelled"
	case TaskStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// IsValid 检查任务状态是否有效
func (ts TaskStatus) IsValid() bool {
	return ts >= TaskStatusPending && ts <= TaskStatusFailed
}

// NewPhaseTask 创建新阶段任务
func NewPhaseTask(name, description string, taskType TaskType, priority TaskPriority, duration time.Duration) *PhaseTask {
	now := time.Now()
	return &PhaseTask{
		ID:          generateTaskID(),
		Name:        name,
		Description: description,
		Type:        taskType,
		Priority:    priority,
		Duration:    duration,
		Progress:    0.0,
		Status:      TaskStatusPending,
		DependsOn:   make([]string, 0),
		Metadata:    make(map[string]interface{}),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Start 开始任务
func (pt *PhaseTask) Start() {
	now := time.Now()
	pt.StartedAt = &now
	pt.Status = TaskStatusInProgress
	pt.UpdatedAt = now
}

// Complete 完成任务
func (pt *PhaseTask) Complete() {
	now := time.Now()
	pt.CompletedAt = &now
	pt.Progress = 100.0
	pt.Status = TaskStatusCompleted
	pt.UpdatedAt = now
}

// Assign 分配任务
func (pt *PhaseTask) Assign(workerID uint64) {
	pt.AssignedTo = &workerID
	pt.UpdatedAt = time.Now()
}

// Unassign 取消分配
func (pt *PhaseTask) Unassign() {
	pt.AssignedTo = nil
	pt.UpdatedAt = time.Now()
}

// AddDependency 添加依赖
func (pt *PhaseTask) AddDependency(taskID string) {
	// 检查是否已存在
	for _, existing := range pt.DependsOn {
		if existing == taskID {
			return
		}
	}
	pt.DependsOn = append(pt.DependsOn, taskID)
	pt.UpdatedAt = time.Now()
}

// RemoveDependency 移除依赖
func (pt *PhaseTask) RemoveDependency(taskID string) {
	for i, dep := range pt.DependsOn {
		if dep == taskID {
			pt.DependsOn = append(pt.DependsOn[:i], pt.DependsOn[i+1:]...)
			pt.UpdatedAt = time.Now()
			return
		}
	}
}

// UpdateProgress 更新进度
func (pt *PhaseTask) UpdateProgress(progress float64) error {
	if progress < 0 || progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100")
	}

	pt.Progress = progress
	pt.UpdatedAt = time.Now()

	if progress >= 100 {
		pt.Complete()
	}

	return nil
}

// Blueprint 蓝图实体
type Blueprint struct {
	ID          string                 `json:"id" bson:"_id"`
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	Version     string                 `json:"version" bson:"version"`
	Author      string                 `json:"author" bson:"author"`
	Category    BuildingCategory       `json:"category" bson:"category"`
	Size        *Size                  `json:"size" bson:"size"`
	Layers      []*BlueprintLayer      `json:"layers" bson:"layers"`
	Materials   []*MaterialRequirement `json:"materials" bson:"materials"`
	Costs       []*ResourceCost        `json:"costs" bson:"costs"`
	Duration    time.Duration          `json:"duration" bson:"duration"`
	Difficulty  int32                  `json:"difficulty" bson:"difficulty"`
	Tags        []string               `json:"tags" bson:"tags"`
	Metadata    map[string]interface{} `json:"metadata" bson:"metadata"`
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updated_at"`
}

// NewBlueprint 创建新蓝图
func NewBlueprint(name, description string, category BuildingCategory) *Blueprint {
	now := time.Now()
	return &Blueprint{
		ID:          generateBlueprintID(),
		Name:        name,
		Description: description,
		Version:     "1.0.0",
		Author:      "",
		Category:    category,
		Size:        NewSize(1, 1, 1),
		Layers:      make([]*BlueprintLayer, 0),
		Materials:   make([]*MaterialRequirement, 0),
		Costs:       make([]*ResourceCost, 0),
		Duration:    1 * time.Hour,
		Difficulty:  1,
		Tags:        make([]string, 0),
		Metadata:    make(map[string]interface{}),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AddLayer 添加层
func (bp *Blueprint) AddLayer(layer *BlueprintLayer) error {
	if layer == nil {
		return fmt.Errorf("blueprint layer cannot be nil")
	}

	bp.Layers = append(bp.Layers, layer)
	bp.UpdatedAt = time.Now()
	return nil
}

// AddMaterial 添加材料需求
func (bp *Blueprint) AddMaterial(material *MaterialRequirement) error {
	if material == nil {
		return fmt.Errorf("material requirement cannot be nil")
	}

	bp.Materials = append(bp.Materials, material)
	bp.UpdatedAt = time.Now()
	return nil
}

// AddCost 添加成本
func (bp *Blueprint) AddCost(cost *ResourceCost) error {
	if cost == nil {
		return fmt.Errorf("resource cost cannot be nil")
	}

	bp.Costs = append(bp.Costs, cost)
	bp.UpdatedAt = time.Now()
	return nil
}

// AddTag 添加标签
func (bp *Blueprint) AddTag(tag string) {
	// 检查是否已存在
	for _, existing := range bp.Tags {
		if existing == tag {
			return
		}
	}
	bp.Tags = append(bp.Tags, tag)
	bp.UpdatedAt = time.Now()
}

// BlueprintLayer 蓝图层实体
type BlueprintLayer struct {
	ID          string                 `json:"id" bson:"_id"`
	Name        string                 `json:"name" bson:"name"`
	Level       int32                  `json:"level" bson:"level"`
	Blocks      []*BlueprintBlock      `json:"blocks" bson:"blocks"`
	Connections []*BlueprintConnection `json:"connections" bson:"connections"`
	Metadata    map[string]interface{} `json:"metadata" bson:"metadata"`
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updated_at"`
}

// NewBlueprintLayer 创建新蓝图层
func NewBlueprintLayer(name string, level int32) *BlueprintLayer {
	now := time.Now()
	return &BlueprintLayer{
		ID:          generateLayerID(),
		Name:        name,
		Level:       level,
		Blocks:      make([]*BlueprintBlock, 0),
		Connections: make([]*BlueprintConnection, 0),
		Metadata:    make(map[string]interface{}),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// BlueprintBlock 蓝图块实体
type BlueprintBlock struct {
	ID          string                 `json:"id" bson:"_id"`
	Type        string                 `json:"type" bson:"type"`
	Position    *Position              `json:"position" bson:"position"`
	Size        *Size                  `json:"size" bson:"size"`
	Orientation Orientation            `json:"orientation" bson:"orientation"`
	Material    string                 `json:"material" bson:"material"`
	Properties  map[string]interface{} `json:"properties" bson:"properties"`
	Metadata    map[string]interface{} `json:"metadata" bson:"metadata"`
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updated_at"`
}

// NewBlueprintBlock 创建新蓝图块
func NewBlueprintBlock(blockType string, position *Position, size *Size) *BlueprintBlock {
	now := time.Now()
	return &BlueprintBlock{
		ID:          generateBlockID(),
		Type:        blockType,
		Position:    position,
		Size:        size,
		Orientation: OrientationNorth,
		Material:    "stone",
		Properties:  make(map[string]interface{}),
		Metadata:    make(map[string]interface{}),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// BlueprintConnection 蓝图连接实体
type BlueprintConnection struct {
	ID         string                 `json:"id" bson:"_id"`
	FromBlock  string                 `json:"from_block" bson:"from_block"`
	ToBlock    string                 `json:"to_block" bson:"to_block"`
	Type       string                 `json:"type" bson:"type"`
	Properties map[string]interface{} `json:"properties" bson:"properties"`
	Metadata   map[string]interface{} `json:"metadata" bson:"metadata"`
	CreatedAt  time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at" bson:"updated_at"`
}

// NewBlueprintConnection 创建新蓝图连接
func NewBlueprintConnection(fromBlock, toBlock, connectionType string) *BlueprintConnection {
	now := time.Now()
	return &BlueprintConnection{
		ID:         generateConnectionID(),
		FromBlock:  fromBlock,
		ToBlock:    toBlock,
		Type:       connectionType,
		Properties: make(map[string]interface{}),
		Metadata:   make(map[string]interface{}),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// MaterialRequirement 材料需求实体
type MaterialRequirement struct {
	MaterialType string   `json:"material_type" bson:"material_type"`
	Quantity     int64    `json:"quantity" bson:"quantity"`
	Quality      float64  `json:"quality" bson:"quality"`
	Optional     bool     `json:"optional" bson:"optional"`
	Alternatives []string `json:"alternatives" bson:"alternatives"`
}

// NewMaterialRequirement 创建新材料需求
func NewMaterialRequirement(materialType string, quantity int64) *MaterialRequirement {
	return &MaterialRequirement{
		MaterialType: materialType,
		Quantity:     quantity,
		Quality:      1.0,
		Optional:     false,
		Alternatives: make([]string, 0),
	}
}

// UpgradeBenefit 升级收益实体
type UpgradeBenefit struct {
	Type        string  `json:"type" bson:"type"`
	Target      string  `json:"target" bson:"target"`
	Value       float64 `json:"value" bson:"value"`
	Description string  `json:"description" bson:"description"`
}

// NewUpgradeBenefit 创建新升级收益
func NewUpgradeBenefit(benefitType, target string, value float64, description string) *UpgradeBenefit {
	return &UpgradeBenefit{
		Type:        benefitType,
		Target:      target,
		Value:       value,
		Description: description,
	}
}

// 辅助函数

// generateConstructionID 生成建造ID
func generateConstructionID() string {
	return fmt.Sprintf("construction_%d", time.Now().UnixNano())
}

// generateUpgradeID 生成升级ID
func generateUpgradeID() string {
	return fmt.Sprintf("upgrade_%d", time.Now().UnixNano())
}

// generateWorkerAssignmentID 生成工人分配ID
func generateWorkerAssignmentID() string {
	return fmt.Sprintf("assignment_%d", time.Now().UnixNano())
}

// generateMaterialUsageID 生成材料使用ID
func generateMaterialUsageID() string {
	return fmt.Sprintf("material_%d", time.Now().UnixNano())
}

// generatePhaseID 生成阶段ID
func generatePhaseID() string {
	return fmt.Sprintf("phase_%d", time.Now().UnixNano())
}

// generateTaskID 生成任务ID
func generateTaskID() string {
	return fmt.Sprintf("task_%d", time.Now().UnixNano())
}

// generateBlueprintID 生成蓝图ID
func generateBlueprintID() string {
	return fmt.Sprintf("blueprint_%d", time.Now().UnixNano())
}

// generateLayerID 生成层ID
func generateLayerID() string {
	return fmt.Sprintf("layer_%d", time.Now().UnixNano())
}

// generateBlockID 生成块ID
func generateBlockID() string {
	return fmt.Sprintf("block_%d", time.Now().UnixNano())
}

// generateConnectionID 生成连接ID
func generateConnectionID() string {
	return fmt.Sprintf("connection_%d", time.Now().UnixNano())
}
