package building

import (
	"fmt"
	"math"
	"time"
)

// 建筑状态相关值对象

// BuildingStatus 建筑状态
type BuildingStatus int32

const (
	BuildingStatusPlanning          BuildingStatus = iota + 1 // 规划中
	BuildingStatusUnderConstruction                           // 建造中
	BuildingStatusActive                                      // 活跃
	BuildingStatusUpgrading                                   // 升级中
	BuildingStatusMaintenance                                 // 维护中
	BuildingStatusDamaged                                     // 受损
	BuildingStatusDestroyed                                   // 被摧毁
	BuildingStatusDemolished                                  // 已拆除
	BuildingStatusCancelled                                   // 已取消
	BuildingStatusInactive                                    // 非活跃
)

// String 返回建筑状态的字符串表示
func (bs BuildingStatus) String() string {
	switch bs {
	case BuildingStatusPlanning:
		return "planning"
	case BuildingStatusUnderConstruction:
		return "under_construction"
	case BuildingStatusActive:
		return "active"
	case BuildingStatusUpgrading:
		return "upgrading"
	case BuildingStatusMaintenance:
		return "maintenance"
	case BuildingStatusDamaged:
		return "damaged"
	case BuildingStatusDestroyed:
		return "destroyed"
	case BuildingStatusDemolished:
		return "demolished"
	case BuildingStatusCancelled:
		return "cancelled"
	case BuildingStatusInactive:
		return "inactive"
	default:
		return "unknown"
	}
}

// IsValid 检查建筑状态是否有效
func (bs BuildingStatus) IsValid() bool {
	return bs >= BuildingStatusPlanning && bs <= BuildingStatusInactive
}

// CanTransitionTo 检查是否可以转换到目标状态
func (bs BuildingStatus) CanTransitionTo(target BuildingStatus) bool {
	switch bs {
	case BuildingStatusPlanning:
		return target == BuildingStatusUnderConstruction || target == BuildingStatusCancelled
	case BuildingStatusUnderConstruction:
		return target == BuildingStatusActive || target == BuildingStatusCancelled
	case BuildingStatusActive:
		return target == BuildingStatusUpgrading || target == BuildingStatusMaintenance ||
			target == BuildingStatusDamaged || target == BuildingStatusDestroyed ||
			target == BuildingStatusDemolished || target == BuildingStatusInactive
	case BuildingStatusUpgrading:
		return target == BuildingStatusActive || target == BuildingStatusCancelled
	case BuildingStatusMaintenance:
		return target == BuildingStatusActive
	case BuildingStatusDamaged:
		return target == BuildingStatusActive || target == BuildingStatusDestroyed ||
			target == BuildingStatusDemolished
	case BuildingStatusDestroyed, BuildingStatusDemolished, BuildingStatusCancelled:
		return false // 终态，不能转换
	case BuildingStatusInactive:
		return target == BuildingStatusActive || target == BuildingStatusDemolished
	default:
		return false
	}
}

// BuildingCategory 建筑分类
type BuildingCategory int32

const (
	BuildingCategoryResidential BuildingCategory = iota + 1 // 住宅
	BuildingCategoryCommercial                              // 商业
	BuildingCategoryIndustrial                              // 工业
	BuildingCategoryMilitary                                // 军事
	BuildingCategoryReligious                               // 宗教
	BuildingCategoryEducational                             // 教育
	BuildingCategoryMedical                                 // 医疗
	BuildingCategoryEntertainment                           // 娱乐
	BuildingCategoryUtility                                 // 公用设施
	BuildingCategoryDecoration                              // 装饰
	BuildingCategorySpecial                                 // 特殊
)

// String 返回建筑分类的字符串表示
func (bc BuildingCategory) String() string {
	switch bc {
	case BuildingCategoryResidential:
		return "residential"
	case BuildingCategoryCommercial:
		return "commercial"
	case BuildingCategoryIndustrial:
		return "industrial"
	case BuildingCategoryMilitary:
		return "military"
	case BuildingCategoryReligious:
		return "religious"
	case BuildingCategoryEducational:
		return "educational"
	case BuildingCategoryMedical:
		return "medical"
	case BuildingCategoryEntertainment:
		return "entertainment"
	case BuildingCategoryUtility:
		return "utility"
	case BuildingCategoryDecoration:
		return "decoration"
	case BuildingCategorySpecial:
		return "special"
	default:
		return "unknown"
	}
}

// IsValid 检查建筑分类是否有效
func (bc BuildingCategory) IsValid() bool {
	return bc >= BuildingCategoryResidential && bc <= BuildingCategorySpecial
}

// 位置和尺寸相关值对象

// Position 位置
type Position struct {
	X int32 `json:"x" bson:"x"`
	Y int32 `json:"y" bson:"y"`
	Z int32 `json:"z" bson:"z"`
}

// NewPosition 创建新位置
func NewPosition(x, y, z int32) *Position {
	return &Position{X: x, Y: y, Z: z}
}

// Distance 计算到另一个位置的距离
func (p *Position) Distance(other *Position) float64 {
	dx := float64(p.X - other.X)
	dy := float64(p.Y - other.Y)
	dz := float64(p.Z - other.Z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// IsAdjacent 检查是否相邻
func (p *Position) IsAdjacent(other *Position) bool {
	dx := abs(p.X - other.X)
	dy := abs(p.Y - other.Y)
	dz := abs(p.Z - other.Z)
	return (dx <= 1 && dy <= 1 && dz <= 1) && !(dx == 0 && dy == 0 && dz == 0)
}

// Validate 验证位置
func (p *Position) Validate() error {
	if p.X < 0 || p.Y < 0 || p.Z < 0 {
		return fmt.Errorf("position coordinates cannot be negative")
	}
	return nil
}

// Clone 克隆位置
func (p *Position) Clone() *Position {
	return &Position{X: p.X, Y: p.Y, Z: p.Z}
}

// Size 尺寸
type Size struct {
	Width  int32 `json:"width" bson:"width"`
	Height int32 `json:"height" bson:"height"`
	Depth  int32 `json:"depth" bson:"depth"`
}

// NewSize 创建新尺寸
func NewSize(width, height, depth int32) *Size {
	return &Size{Width: width, Height: height, Depth: depth}
}

// Volume 计算体积
func (s *Size) Volume() int32 {
	return s.Width * s.Height * s.Depth
}

// Area 计算面积
func (s *Size) Area() int32 {
	return s.Width * s.Height
}

// Validate 验证尺寸
func (s *Size) Validate() error {
	if s.Width <= 0 || s.Height <= 0 || s.Depth <= 0 {
		return fmt.Errorf("size dimensions must be positive")
	}
	return nil
}

// IsValid 检查尺寸是否有效
func (s *Size) IsValid() bool {
	return s != nil && s.Width > 0 && s.Height > 0 && s.Depth > 0
}

// Clone 克隆尺寸
func (s *Size) Clone() *Size {
	return &Size{Width: s.Width, Height: s.Height, Depth: s.Depth}
}

// BoundingBox 边界框
type BoundingBox struct {
	MinX int32 `json:"min_x" bson:"min_x"`
	MinY int32 `json:"min_y" bson:"min_y"`
	MinZ int32 `json:"min_z" bson:"min_z"`
	MaxX int32 `json:"max_x" bson:"max_x"`
	MaxY int32 `json:"max_y" bson:"max_y"`
	MaxZ int32 `json:"max_z" bson:"max_z"`
}

// NewBoundingBox 创建新边界框
func NewBoundingBox(minX, minY, minZ, maxX, maxY, maxZ int32) *BoundingBox {
	return &BoundingBox{
		MinX: minX, MinY: minY, MinZ: minZ,
		MaxX: maxX, MaxY: maxY, MaxZ: maxZ,
	}
}

// Contains 检查是否包含位置
func (bb *BoundingBox) Contains(pos *Position) bool {
	return pos.X >= bb.MinX && pos.X <= bb.MaxX &&
		pos.Y >= bb.MinY && pos.Y <= bb.MaxY &&
		pos.Z >= bb.MinZ && pos.Z <= bb.MaxZ
}

// Intersects 检查是否与另一个边界框相交
func (bb *BoundingBox) Intersects(other *BoundingBox) bool {
	return bb.MinX <= other.MaxX && bb.MaxX >= other.MinX &&
		bb.MinY <= other.MaxY && bb.MaxY >= other.MinY &&
		bb.MinZ <= other.MaxZ && bb.MaxZ >= other.MinZ
}

// Volume 计算体积
func (bb *BoundingBox) Volume() int32 {
	return (bb.MaxX - bb.MinX + 1) * (bb.MaxY - bb.MinY + 1) * (bb.MaxZ - bb.MinZ + 1)
}

// Orientation 朝向
type Orientation int32

const (
	OrientationNorth Orientation = iota + 1 // 北
	OrientationEast                         // 东
	OrientationSouth                        // 南
	OrientationWest                         // 西
	OrientationUp                           // 上
	OrientationDown                         // 下
)

// String 返回朝向的字符串表示
func (o Orientation) String() string {
	switch o {
	case OrientationNorth:
		return "north"
	case OrientationEast:
		return "east"
	case OrientationSouth:
		return "south"
	case OrientationWest:
		return "west"
	case OrientationUp:
		return "up"
	case OrientationDown:
		return "down"
	default:
		return "unknown"
	}
}

// IsValid 检查朝向是否有效
func (o Orientation) IsValid() bool {
	return o >= OrientationNorth && o <= OrientationDown
}

// Opposite 获取相反朝向
func (o Orientation) Opposite() Orientation {
	switch o {
	case OrientationNorth:
		return OrientationSouth
	case OrientationEast:
		return OrientationWest
	case OrientationSouth:
		return OrientationNorth
	case OrientationWest:
		return OrientationEast
	case OrientationUp:
		return OrientationDown
	case OrientationDown:
		return OrientationUp
	default:
		return OrientationNorth
	}
}

// 资源和成本相关值对象

// ResourceCost 资源成本
type ResourceCost struct {
	ResourceType string `json:"resource_type" bson:"resource_type"`
	Amount       int64  `json:"amount" bson:"amount"`
	Optional     bool   `json:"optional" bson:"optional"`
}

// NewResourceCost 创建新资源成本
func NewResourceCost(resourceType string, amount int64) *ResourceCost {
	return &ResourceCost{
		ResourceType: resourceType,
		Amount:       amount,
		Optional:     false,
	}
}

// Validate 验证资源成本
func (rc *ResourceCost) Validate() error {
	if rc.ResourceType == "" {
		return fmt.Errorf("resource type cannot be empty")
	}
	if rc.Amount < 0 {
		return fmt.Errorf("resource amount cannot be negative")
	}
	return nil
}

// Clone 克隆资源成本
func (rc *ResourceCost) Clone() *ResourceCost {
	return &ResourceCost{
		ResourceType: rc.ResourceType,
		Amount:       rc.Amount,
		Optional:     rc.Optional,
	}
}

// 要求相关值对象

// RequirementType 要求类型
type RequirementType int32

const (
	RequirementTypeLevel      RequirementType = iota + 1 // 等级要求
	RequirementTypeResource                               // 资源要求
	RequirementTypeBuilding                               // 建筑要求
	RequirementTypeTechnology                             // 科技要求
	RequirementTypePopulation                             // 人口要求
	RequirementTypeTime                                   // 时间要求
	RequirementTypeCustom                                 // 自定义要求
)

// String 返回要求类型的字符串表示
func (rt RequirementType) String() string {
	switch rt {
	case RequirementTypeLevel:
		return "level"
	case RequirementTypeResource:
		return "resource"
	case RequirementTypeBuilding:
		return "building"
	case RequirementTypeTechnology:
		return "technology"
	case RequirementTypePopulation:
		return "population"
	case RequirementTypeTime:
		return "time"
	case RequirementTypeCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// IsValid 检查要求类型是否有效
func (rt RequirementType) IsValid() bool {
	return rt >= RequirementTypeLevel && rt <= RequirementTypeCustom
}

// Requirement 要求
type Requirement struct {
	Type        RequirementType `json:"type" bson:"type"`
	Target      string          `json:"target" bson:"target"`
	Value       int64           `json:"value" bson:"value"`
	Operator    string          `json:"operator" bson:"operator"`
	Description string          `json:"description" bson:"description"`
	Optional    bool            `json:"optional" bson:"optional"`
	Met         bool            `json:"met" bson:"met"`
}

// NewRequirement 创建新要求
func NewRequirement(reqType RequirementType, target string, value int64, operator, description string) *Requirement {
	return &Requirement{
		Type:        reqType,
		Target:      target,
		Value:       value,
		Operator:    operator,
		Description: description,
		Optional:    false,
		Met:         false,
	}
}

// IsMet 检查要求是否满足
func (r *Requirement) IsMet() bool {
	return r.Met || r.Optional
}

// SetMet 设置要求满足状态
func (r *Requirement) SetMet(met bool) {
	r.Met = met
}

// Validate 验证要求
func (r *Requirement) Validate() error {
	if !r.Type.IsValid() {
		return fmt.Errorf("invalid requirement type: %v", r.Type)
	}
	if r.Target == "" {
		return fmt.Errorf("requirement target cannot be empty")
	}
	if r.Operator == "" {
		return fmt.Errorf("requirement operator cannot be empty")
	}
	return nil
}

// Clone 克隆要求
func (r *Requirement) Clone() *Requirement {
	return &Requirement{
		Type:        r.Type,
		Target:      r.Target,
		Value:       r.Value,
		Operator:    r.Operator,
		Description: r.Description,
		Optional:    r.Optional,
		Met:         r.Met,
	}
}

// 效果相关值对象

// EffectType 效果类型
type EffectType int32

const (
	EffectTypeProduction  EffectType = iota + 1 // 生产效果
	EffectTypeDefense                            // 防御效果
	EffectTypeEfficiency                         // 效率效果
	EffectTypeCapacity                           // 容量效果
	EffectTypeSpeed                              // 速度效果
	EffectTypeCost                               // 成本效果
	EffectTypeHealth                             // 生命值效果
	EffectTypeDurability                         // 耐久度效果
	EffectTypeRange                              // 范围效果
	EffectTypeCustom                             // 自定义效果
)

// String 返回效果类型的字符串表示
func (et EffectType) String() string {
	switch et {
	case EffectTypeProduction:
		return "production"
	case EffectTypeDefense:
		return "defense"
	case EffectTypeEfficiency:
		return "efficiency"
	case EffectTypeCapacity:
		return "capacity"
	case EffectTypeSpeed:
		return "speed"
	case EffectTypeCost:
		return "cost"
	case EffectTypeHealth:
		return "health"
	case EffectTypeDurability:
		return "durability"
	case EffectTypeRange:
		return "range"
	case EffectTypeCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// IsValid 检查效果类型是否有效
func (et EffectType) IsValid() bool {
	return et >= EffectTypeProduction && et <= EffectTypeCustom
}

// BuildingEffect 建筑效果
type BuildingEffect struct {
	Type        EffectType             `json:"type" bson:"type"`
	Target      string                 `json:"target" bson:"target"`
	Value       float64                `json:"value" bson:"value"`
	Duration    *time.Duration         `json:"duration,omitempty" bson:"duration,omitempty"`
	StartTime   *time.Time             `json:"start_time,omitempty" bson:"start_time,omitempty"`
	EndTime     *time.Time             `json:"end_time,omitempty" bson:"end_time,omitempty"`
	Stackable   bool                   `json:"stackable" bson:"stackable"`
	Permanent   bool                   `json:"permanent" bson:"permanent"`
	Conditions  map[string]interface{} `json:"conditions" bson:"conditions"`
	Description string                 `json:"description" bson:"description"`
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updated_at"`
}

// NewBuildingEffect 创建新建筑效果
func NewBuildingEffect(effectType EffectType, target string, value float64) *BuildingEffect {
	now := time.Now()
	return &BuildingEffect{
		Type:        effectType,
		Target:      target,
		Value:       value,
		Stackable:   false,
		Permanent:   true,
		Conditions:  make(map[string]interface{}),
		Description: fmt.Sprintf("%s effect on %s", effectType.String(), target),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// IsActive 检查效果是否活跃
func (be *BuildingEffect) IsActive() bool {
	if be.Permanent {
		return true
	}
	
	if be.EndTime != nil {
		return time.Now().Before(*be.EndTime)
	}
	
	if be.StartTime != nil && be.Duration != nil {
		endTime := be.StartTime.Add(*be.Duration)
		return time.Now().Before(endTime)
	}
	
	return true
}

// SetDuration 设置持续时间
func (be *BuildingEffect) SetDuration(duration time.Duration) {
	be.Duration = &duration
	if be.StartTime != nil {
		endTime := be.StartTime.Add(duration)
		be.EndTime = &endTime
	}
	be.UpdatedAt = time.Now()
}

// Start 开始效果
func (be *BuildingEffect) Start() {
	now := time.Now()
	be.StartTime = &now
	if be.Duration != nil {
		endTime := now.Add(*be.Duration)
		be.EndTime = &endTime
	}
	be.UpdatedAt = now
}

// Validate 验证建筑效果
func (be *BuildingEffect) Validate() error {
	if !be.Type.IsValid() {
		return fmt.Errorf("invalid effect type: %v", be.Type)
	}
	if be.Target == "" {
		return fmt.Errorf("effect target cannot be empty")
	}
	return nil
}

// Clone 克隆建筑效果
func (be *BuildingEffect) Clone() *BuildingEffect {
	clone := &BuildingEffect{
		Type:        be.Type,
		Target:      be.Target,
		Value:       be.Value,
		Stackable:   be.Stackable,
		Permanent:   be.Permanent,
		Conditions:  make(map[string]interface{}),
		Description: be.Description,
		CreatedAt:   be.CreatedAt,
		UpdatedAt:   be.UpdatedAt,
	}
	
	// 深拷贝map
	for k, v := range be.Conditions {
		clone.Conditions[k] = v
	}
	
	// 深拷贝指针
	if be.Duration != nil {
		duration := *be.Duration
		clone.Duration = &duration
	}
	if be.StartTime != nil {
		startTime := *be.StartTime
		clone.StartTime = &startTime
	}
	if be.EndTime != nil {
		endTime := *be.EndTime
		clone.EndTime = &endTime
	}
	
	return clone
}

// 生产相关值对象

// ProductionType 生产类型
type ProductionType int32

const (
	ProductionTypeResource ProductionType = iota + 1 // 资源生产
	ProductionTypeItem                                // 物品生产
	ProductionTypeUnit                                // 单位生产
	ProductionTypeService                             // 服务生产
	ProductionTypeCustom                              // 自定义生产
)

// String 返回生产类型的字符串表示
func (pt ProductionType) String() string {
	switch pt {
	case ProductionTypeResource:
		return "resource"
	case ProductionTypeItem:
		return "item"
	case ProductionTypeUnit:
		return "unit"
	case ProductionTypeService:
		return "service"
	case ProductionTypeCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// IsValid 检查生产类型是否有效
func (pt ProductionType) IsValid() bool {
	return pt >= ProductionTypeResource && pt <= ProductionTypeCustom
}

// ProductionInfo 生产信息
type ProductionInfo struct {
	Type           ProductionType         `json:"type" bson:"type"`
	Outputs        []*ProductionOutput    `json:"outputs" bson:"outputs"`
	Inputs         []*ProductionInput     `json:"inputs" bson:"inputs"`
	Rate           float64                `json:"rate" bson:"rate"`
	Efficiency     float64                `json:"efficiency" bson:"efficiency"`
	Capacity       int32                  `json:"capacity" bson:"capacity"`
	Queue          []*ProductionTask      `json:"queue" bson:"queue"`
	CurrentTask    *ProductionTask        `json:"current_task,omitempty" bson:"current_task,omitempty"`
	AutoProduction bool                   `json:"auto_production" bson:"auto_production"`
	Conditions     map[string]interface{} `json:"conditions" bson:"conditions"`
	CreatedAt      time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at" bson:"updated_at"`
}

// NewProductionInfo 创建新生产信息
func NewProductionInfo(productionType ProductionType) *ProductionInfo {
	now := time.Now()
	return &ProductionInfo{
		Type:           productionType,
		Outputs:        make([]*ProductionOutput, 0),
		Inputs:         make([]*ProductionInput, 0),
		Rate:           1.0,
		Efficiency:     1.0,
		Capacity:       10,
		Queue:          make([]*ProductionTask, 0),
		AutoProduction: false,
		Conditions:     make(map[string]interface{}),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// AddOutput 添加产出
func (pi *ProductionInfo) AddOutput(output *ProductionOutput) {
	pi.Outputs = append(pi.Outputs, output)
	pi.UpdatedAt = time.Now()
}

// AddInput 添加输入
func (pi *ProductionInfo) AddInput(input *ProductionInput) {
	pi.Inputs = append(pi.Inputs, input)
	pi.UpdatedAt = time.Now()
}

// AddTask 添加生产任务
func (pi *ProductionInfo) AddTask(task *ProductionTask) {
	pi.Queue = append(pi.Queue, task)
	pi.UpdatedAt = time.Now()
}

// StartNextTask 开始下一个任务
func (pi *ProductionInfo) StartNextTask() *ProductionTask {
	if len(pi.Queue) == 0 {
		return nil
	}
	
	task := pi.Queue[0]
	pi.Queue = pi.Queue[1:]
	pi.CurrentTask = task
	task.Start()
	pi.UpdatedAt = time.Now()
	return task
}

// CompleteCurrentTask 完成当前任务
func (pi *ProductionInfo) CompleteCurrentTask() *ProductionTask {
	if pi.CurrentTask == nil {
		return nil
	}
	
	task := pi.CurrentTask
	task.Complete()
	pi.CurrentTask = nil
	pi.UpdatedAt = time.Now()
	return task
}

// Validate 验证生产信息
func (pi *ProductionInfo) Validate() error {
	if !pi.Type.IsValid() {
		return fmt.Errorf("invalid production type: %v", pi.Type)
	}
	if pi.Rate <= 0 {
		return fmt.Errorf("production rate must be positive")
	}
	if pi.Efficiency < 0 || pi.Efficiency > 2 {
		return fmt.Errorf("production efficiency must be between 0 and 2")
	}
	if pi.Capacity <= 0 {
		return fmt.Errorf("production capacity must be positive")
	}
	return nil
}

// ProductionOutput 生产产出
type ProductionOutput struct {
	ResourceType string  `json:"resource_type" bson:"resource_type"`
	Amount       int64   `json:"amount" bson:"amount"`
	Rate         float64 `json:"rate" bson:"rate"`
	Quality      float64 `json:"quality" bson:"quality"`
}

// NewProductionOutput 创建新生产产出
func NewProductionOutput(resourceType string, amount int64, rate float64) *ProductionOutput {
	return &ProductionOutput{
		ResourceType: resourceType,
		Amount:       amount,
		Rate:         rate,
		Quality:      1.0,
	}
}

// ProductionInput 生产输入
type ProductionInput struct {
	ResourceType string  `json:"resource_type" bson:"resource_type"`
	Amount       int64   `json:"amount" bson:"amount"`
	Rate         float64 `json:"rate" bson:"rate"`
	Optional     bool    `json:"optional" bson:"optional"`
}

// NewProductionInput 创建新生产输入
func NewProductionInput(resourceType string, amount int64, rate float64) *ProductionInput {
	return &ProductionInput{
		ResourceType: resourceType,
		Amount:       amount,
		Rate:         rate,
		Optional:     false,
	}
}

// ProductionTask 生产任务
type ProductionTask struct {
	ID          string                 `json:"id" bson:"id"`
	Type        ProductionType         `json:"type" bson:"type"`
	Target      string                 `json:"target" bson:"target"`
	Quantity    int32                  `json:"quantity" bson:"quantity"`
	Progress    float64                `json:"progress" bson:"progress"`
	Duration    time.Duration          `json:"duration" bson:"duration"`
	StartedAt   *time.Time             `json:"started_at,omitempty" bson:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
	Status      ProductionTaskStatus   `json:"status" bson:"status"`
	Inputs      []*ProductionInput     `json:"inputs" bson:"inputs"`
	Outputs     []*ProductionOutput    `json:"outputs" bson:"outputs"`
	Metadata    map[string]interface{} `json:"metadata" bson:"metadata"`
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updated_at"`
}

// ProductionTaskStatus 生产任务状态
type ProductionTaskStatus int32

const (
	ProductionTaskStatusPending    ProductionTaskStatus = iota + 1 // 等待中
	ProductionTaskStatusInProgress                                   // 进行中
	ProductionTaskStatusCompleted                                    // 已完成
	ProductionTaskStatusCancelled                                    // 已取消
	ProductionTaskStatusFailed                                       // 失败
)

// String 返回生产任务状态的字符串表示
func (pts ProductionTaskStatus) String() string {
	switch pts {
	case ProductionTaskStatusPending:
		return "pending"
	case ProductionTaskStatusInProgress:
		return "in_progress"
	case ProductionTaskStatusCompleted:
		return "completed"
	case ProductionTaskStatusCancelled:
		return "cancelled"
	case ProductionTaskStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// NewProductionTask 创建新生产任务
func NewProductionTask(taskType ProductionType, target string, quantity int32, duration time.Duration) *ProductionTask {
	now := time.Now()
	return &ProductionTask{
		ID:        fmt.Sprintf("task_%d", now.UnixNano()),
		Type:      taskType,
		Target:    target,
		Quantity:  quantity,
		Progress:  0.0,
		Duration:  duration,
		Status:    ProductionTaskStatusPending,
		Inputs:    make([]*ProductionInput, 0),
		Outputs:   make([]*ProductionOutput, 0),
		Metadata:  make(map[string]interface{}),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Start 开始任务
func (pt *ProductionTask) Start() {
	now := time.Now()
	pt.StartedAt = &now
	pt.Status = ProductionTaskStatusInProgress
	pt.UpdatedAt = now
}

// Complete 完成任务
func (pt *ProductionTask) Complete() {
	now := time.Now()
	pt.CompletedAt = &now
	pt.Progress = 100.0
	pt.Status = ProductionTaskStatusCompleted
	pt.UpdatedAt = now
}

// Cancel 取消任务
func (pt *ProductionTask) Cancel() {
	pt.Status = ProductionTaskStatusCancelled
	pt.UpdatedAt = time.Now()
}

// UpdateProgress 更新进度
func (pt *ProductionTask) UpdateProgress(progress float64) {
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}
	pt.Progress = progress
	pt.UpdatedAt = time.Now()
	
	if progress >= 100 {
		pt.Complete()
	}
}

// 存储相关值对象

// StorageType 存储类型
type StorageType int32

const (
	StorageTypeGeneral   StorageType = iota + 1 // 通用存储
	StorageTypeResource                          // 资源存储
	StorageTypeItem                              // 物品存储
	StorageTypeFood                              // 食物存储
	StorageTypeWeapon                            // 武器存储
	StorageTypeArmor                             // 装备存储
	StorageTypeLiquid                            // 液体存储
	StorageTypeGas                               // 气体存储
	StorageTypeSpecial                           // 特殊存储
)

// String 返回存储类型的字符串表示
func (st StorageType) String() string {
	switch st {
	case StorageTypeGeneral:
		return "general"
	case StorageTypeResource:
		return "resource"
	case StorageTypeItem:
		return "item"
	case StorageTypeFood:
		return "food"
	case StorageTypeWeapon:
		return "weapon"
	case StorageTypeArmor:
		return "armor"
	case StorageTypeLiquid:
		return "liquid"
	case StorageTypeGas:
		return "gas"
	case StorageTypeSpecial:
		return "special"
	default:
		return "unknown"
	}
}

// IsValid 检查存储类型是否有效
func (st StorageType) IsValid() bool {
	return st >= StorageTypeGeneral && st <= StorageTypeSpecial
}

// StorageInfo 存储信息
type StorageInfo struct {
	Type         StorageType            `json:"type" bson:"type"`
	Capacity     int64                  `json:"capacity" bson:"capacity"`
	Used         int64                  `json:"used" bson:"used"`
	Reserved     int64                  `json:"reserved" bson:"reserved"`
	Items        []*StorageItem         `json:"items" bson:"items"`
	Filters      []string               `json:"filters" bson:"filters"`
	AutoSort     bool                   `json:"auto_sort" bson:"auto_sort"`
	AutoCompact  bool                   `json:"auto_compact" bson:"auto_compact"`
	AccessRules  []*AccessRule          `json:"access_rules" bson:"access_rules"`
	Conditions   map[string]interface{} `json:"conditions" bson:"conditions"`
	CreatedAt    time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" bson:"updated_at"`
}

// NewStorageInfo 创建新存储信息
func NewStorageInfo(storageType StorageType, capacity int64) *StorageInfo {
	now := time.Now()
	return &StorageInfo{
		Type:        storageType,
		Capacity:    capacity,
		Used:        0,
		Reserved:    0,
		Items:       make([]*StorageItem, 0),
		Filters:     make([]string, 0),
		AutoSort:    true,
		AutoCompact: true,
		AccessRules: make([]*AccessRule, 0),
		Conditions:  make(map[string]interface{}),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// GetAvailable 获取可用容量
func (si *StorageInfo) GetAvailable() int64 {
	return si.Capacity - si.Used - si.Reserved
}

// GetUsagePercentage 获取使用率
func (si *StorageInfo) GetUsagePercentage() float64 {
	if si.Capacity == 0 {
		return 0.0
	}
	return float64(si.Used) / float64(si.Capacity) * 100.0
}

// IsFull 检查是否已满
func (si *StorageInfo) IsFull() bool {
	return si.GetAvailable() <= 0
}

// CanStore 检查是否可以存储
func (si *StorageInfo) CanStore(itemType string, quantity int64) bool {
	if si.GetAvailable() < quantity {
		return false
	}
	
	// 检查过滤器
	if len(si.Filters) > 0 {
		allowed := false
		for _, filter := range si.Filters {
			if filter == itemType || filter == "*" {
				allowed = true
				break
			}
		}
		if !allowed {
			return false
		}
	}
	
	return true
}

// AddItem 添加物品
func (si *StorageInfo) AddItem(item *StorageItem) error {
	if !si.CanStore(item.ItemType, item.Quantity) {
		return fmt.Errorf("cannot store item: insufficient space or not allowed")
	}
	
	// 查找是否已存在相同物品
	for _, existing := range si.Items {
		if existing.ItemType == item.ItemType && existing.CanStack(item) {
			existing.Quantity += item.Quantity
			existing.UpdatedAt = time.Now()
			si.Used += item.Quantity
			si.UpdatedAt = time.Now()
			return nil
		}
	}
	
	// 添加新物品
	si.Items = append(si.Items, item)
	si.Used += item.Quantity
	si.UpdatedAt = time.Now()
	return nil
}

// RemoveItem 移除物品
func (si *StorageInfo) RemoveItem(itemType string, quantity int64) error {
	for i, item := range si.Items {
		if item.ItemType == itemType {
			if item.Quantity < quantity {
				return fmt.Errorf("insufficient quantity: have %d, need %d", item.Quantity, quantity)
			}
			
			item.Quantity -= quantity
			item.UpdatedAt = time.Now()
			si.Used -= quantity
			
			// 如果数量为0，移除物品
			if item.Quantity <= 0 {
				si.Items = append(si.Items[:i], si.Items[i+1:]...)
			}
			
			si.UpdatedAt = time.Now()
			return nil
		}
	}
	
	return fmt.Errorf("item not found: %s", itemType)
}

// Validate 验证存储信息
func (si *StorageInfo) Validate() error {
	if !si.Type.IsValid() {
		return fmt.Errorf("invalid storage type: %v", si.Type)
	}
	if si.Capacity <= 0 {
		return fmt.Errorf("storage capacity must be positive")
	}
	if si.Used < 0 {
		return fmt.Errorf("storage used cannot be negative")
	}
	if si.Reserved < 0 {
		return fmt.Errorf("storage reserved cannot be negative")
	}
	if si.Used+si.Reserved > si.Capacity {
		return fmt.Errorf("storage used+reserved cannot exceed capacity")
	}
	return nil
}

// StorageItem 存储物品
type StorageItem struct {
	ItemType   string                 `json:"item_type" bson:"item_type"`
	Quantity   int64                  `json:"quantity" bson:"quantity"`
	Quality    float64                `json:"quality" bson:"quality"`
	Durability float64                `json:"durability" bson:"durability"`
	Stackable  bool                   `json:"stackable" bson:"stackable"`
	Metadata   map[string]interface{} `json:"metadata" bson:"metadata"`
	CreatedAt  time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at" bson:"updated_at"`
}

// NewStorageItem 创建新存储物品
func NewStorageItem(itemType string, quantity int64) *StorageItem {
	now := time.Now()
	return &StorageItem{
		ItemType:   itemType,
		Quantity:   quantity,
		Quality:    1.0,
		Durability: 1.0,
		Stackable:  true,
		Metadata:   make(map[string]interface{}),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// CanStack 检查是否可以堆叠
func (si *StorageItem) CanStack(other *StorageItem) bool {
	return si.Stackable && other.Stackable &&
		si.ItemType == other.ItemType &&
		si.Quality == other.Quality &&
		si.Durability == other.Durability
}

// AccessRule 访问规则
type AccessRule struct {
	UserID     *uint64   `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Role       *string   `json:"role,omitempty" bson:"role,omitempty"`
	Permission string    `json:"permission" bson:"permission"`
	ItemTypes  []string  `json:"item_types" bson:"item_types"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
}

// NewAccessRule 创建新访问规则
func NewAccessRule(permission string) *AccessRule {
	return &AccessRule{
		Permission: permission,
		ItemTypes:  make([]string, 0),
		CreatedAt:  time.Now(),
	}
}

// 防御相关值对象

// DamageType 伤害类型
type DamageType int32

const (
	DamageTypePhysical DamageType = iota + 1 // 物理伤害
	DamageTypeFire                           // 火焰伤害
	DamageTypeIce                            // 冰霜伤害
	DamageTypeLightning                      // 闪电伤害
	DamageTypePoison                         // 毒素伤害
	DamageTypeAcid                           // 酸性伤害
	DamageTypeMagic                          // 魔法伤害
	DamageTypeHoly                           // 神圣伤害
	DamageTypeDark                           // 黑暗伤害
	DamageTypeCustom                         // 自定义伤害
)

// String 返回伤害类型的字符串表示
func (dt DamageType) String() string {
	switch dt {
	case DamageTypePhysical:
		return "physical"
	case DamageTypeFire:
		return "fire"
	case DamageTypeIce:
		return "ice"
	case DamageTypeLightning:
		return "lightning"
	case DamageTypePoison:
		return "poison"
	case DamageTypeAcid:
		return "acid"
	case DamageTypeMagic:
		return "magic"
	case DamageTypeHoly:
		return "holy"
	case DamageTypeDark:
		return "dark"
	case DamageTypeCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// IsValid 检查伤害类型是否有效
func (dt DamageType) IsValid() bool {
	return dt >= DamageTypePhysical && dt <= DamageTypeCustom
}

// DefenseInfo 防御信息
type DefenseInfo struct {
	Armor        int32                  `json:"armor" bson:"armor"`
	Resistances  map[DamageType]int32   `json:"resistances" bson:"resistances"`
	Immunities   []DamageType           `json:"immunities" bson:"immunities"`
	Weaknesses   []DamageType           `json:"weaknesses" bson:"weaknesses"`
	Shield       int32                  `json:"shield" bson:"shield"`
	MaxShield    int32                  `json:"max_shield" bson:"max_shield"`
	RegenRate    float64                `json:"regen_rate" bson:"regen_rate"`
	Absorption   float64                `json:"absorption" bson:"absorption"`
	Reflection   float64                `json:"reflection" bson:"reflection"`
	Conditions   map[string]interface{} `json:"conditions" bson:"conditions"`
	CreatedAt    time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" bson:"updated_at"`
}

// NewDefenseInfo 创建新防御信息
func NewDefenseInfo() *DefenseInfo {
	now := time.Now()
	return &DefenseInfo{
		Armor:       0,
		Resistances: make(map[DamageType]int32),
		Immunities:  make([]DamageType, 0),
		Weaknesses:  make([]DamageType, 0),
		Shield:      0,
		MaxShield:   0,
		RegenRate:   0.0,
		Absorption:  0.0,
		Reflection:  0.0,
		Conditions:  make(map[string]interface{}),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// GetDefenseValue 获取对特定伤害类型的防御值
func (di *DefenseInfo) GetDefenseValue(damageType DamageType) int32 {
	// 检查免疫
	for _, immunity := range di.Immunities {
		if immunity == damageType {
			return 999999 // 免疫，返回极高防御值
		}
	}
	
	// 检查弱点
	for _, weakness := range di.Weaknesses {
		if weakness == damageType {
			return -di.Armor // 弱点，负防御
		}
	}
	
	// 基础护甲 + 特定抗性
	defenseValue := di.Armor
	if resistance, exists := di.Resistances[damageType]; exists {
		defenseValue += resistance
	}
	
	return defenseValue
}

// AddResistance 添加抗性
func (di *DefenseInfo) AddResistance(damageType DamageType, value int32) {
	di.Resistances[damageType] = value
	di.UpdatedAt = time.Now()
}

// AddImmunity 添加免疫
func (di *DefenseInfo) AddImmunity(damageType DamageType) {
	// 检查是否已存在
	for _, existing := range di.Immunities {
		if existing == damageType {
			return
		}
	}
	di.Immunities = append(di.Immunities, damageType)
	di.UpdatedAt = time.Now()
}

// AddWeakness 添加弱点
func (di *DefenseInfo) AddWeakness(damageType DamageType) {
	// 检查是否已存在
	for _, existing := range di.Weaknesses {
		if existing == damageType {
			return
		}
	}
	di.Weaknesses = append(di.Weaknesses, damageType)
	di.UpdatedAt = time.Now()
}

// RegenerateShield 恢复护盾
func (di *DefenseInfo) RegenerateShield() {
	if di.Shield < di.MaxShield && di.RegenRate > 0 {
		di.Shield += int32(di.RegenRate)
		if di.Shield > di.MaxShield {
			di.Shield = di.MaxShield
		}
		di.UpdatedAt = time.Now()
	}
}

// Validate 验证防御信息
func (di *DefenseInfo) Validate() error {
	if di.Armor < 0 {
		return fmt.Errorf("armor cannot be negative")
	}
	if di.Shield < 0 || di.Shield > di.MaxShield {
		return fmt.Errorf("shield must be between 0 and max shield")
	}
	if di.MaxShield < 0 {
		return fmt.Errorf("max shield cannot be negative")
	}
	if di.RegenRate < 0 {
		return fmt.Errorf("regen rate cannot be negative")
	}
	if di.Absorption < 0 || di.Absorption > 1 {
		return fmt.Errorf("absorption must be between 0 and 1")
	}
	if di.Reflection < 0 || di.Reflection > 1 {
		return fmt.Errorf("reflection must be between 0 and 1")
	}
	return nil
}

// 工人相关值对象

// WorkerRole 工人角色
type WorkerRole int32

const (
	WorkerRoleGeneral     WorkerRole = iota + 1 // 通用工人
	WorkerRoleBuilder                            // 建造工人
	WorkerRoleMaintenance                        // 维护工人
	WorkerRoleOperator                           // 操作工人
	WorkerRoleGuard                              // 守卫
	WorkerRoleManager                            // 管理员
	WorkerRoleSpecialist                         // 专家
)

// String 返回工人角色的字符串表示
func (wr WorkerRole) String() string {
	switch wr {
	case WorkerRoleGeneral:
		return "general"
	case WorkerRoleBuilder:
		return "builder"
	case WorkerRoleMaintenance:
		return "maintenance"
	case WorkerRoleOperator:
		return "operator"
	case WorkerRoleGuard:
		return "guard"
	case WorkerRoleManager:
		return "manager"
	case WorkerRoleSpecialist:
		return "specialist"
	default:
		return "unknown"
	}
}

// IsValid 检查工人角色是否有效
func (wr WorkerRole) IsValid() bool {
	return wr >= WorkerRoleGeneral && wr <= WorkerRoleSpecialist
}

// WorkerStatus 工人状态
type WorkerStatus int32

const (
	WorkerStatusActive      WorkerStatus = iota + 1 // 活跃
	WorkerStatusIdle                                 // 空闲
	WorkerStatusBusy                                 // 忙碌
	WorkerStatusResting                              // 休息
	WorkerStatusSick                                 // 生病
	WorkerStatusOnLeave                              // 请假
	WorkerStatusDismissed                            // 解雇
)

// String 返回工人状态的字符串表示
func (ws WorkerStatus) String() string {
	switch ws {
	case WorkerStatusActive:
		return "active"
	case WorkerStatusIdle:
		return "idle"
	case WorkerStatusBusy:
		return "busy"
	case WorkerStatusResting:
		return "resting"
	case WorkerStatusSick:
		return "sick"
	case WorkerStatusOnLeave:
		return "on_leave"
	case WorkerStatusDismissed:
		return "dismissed"
	default:
		return "unknown"
	}
}

// IsValid 检查工人状态是否有效
func (ws WorkerStatus) IsValid() bool {
	return ws >= WorkerStatusActive && ws <= WorkerStatusDismissed
}

// WorkerInfo 工人信息
type WorkerInfo struct {
	WorkerID   uint64       `json:"worker_id" bson:"worker_id"`
	Role       WorkerRole   `json:"role" bson:"role"`
	Status     WorkerStatus `json:"status" bson:"status"`
	Efficiency float64      `json:"efficiency" bson:"efficiency"`
	Experience int32        `json:"experience" bson:"experience"`
	Level      int32        `json:"level" bson:"level"`
	Salary     int64        `json:"salary" bson:"salary"`
	AssignedAt time.Time    `json:"assigned_at" bson:"assigned_at"`
	CreatedAt  time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at" bson:"updated_at"`
}

// NewWorkerInfo 创建新工人信息
func NewWorkerInfo(workerID uint64, role WorkerRole) *WorkerInfo {
	now := time.Now()
	return &WorkerInfo{
		WorkerID:   workerID,
		Role:       role,
		Status:     WorkerStatusActive,
		Efficiency: 1.0,
		Experience: 0,
		Level:      1,
		Salary:     100,
		AssignedAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// VisitorInfo 访客信息
type VisitorInfo struct {
	VisitorID uint64    `json:"visitor_id" bson:"visitor_id"`
	Purpose   string    `json:"purpose" bson:"purpose"`
	ArrivedAt time.Time `json:"arrived_at" bson:"arrived_at"`
	LeftAt    *time.Time `json:"left_at,omitempty" bson:"left_at,omitempty"`
	Duration  time.Duration `json:"duration" bson:"duration"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// NewVisitorInfo 创建新访客信息
func NewVisitorInfo(visitorID uint64, purpose string) *VisitorInfo {
	now := time.Now()
	return &VisitorInfo{
		VisitorID: visitorID,
		Purpose:   purpose,
		ArrivedAt: now,
		Duration:  0,
		CreatedAt: now,
	}
}

// Leave 离开
func (vi *VisitorInfo) Leave() {
	now := time.Now()
	vi.LeftAt = &now
	vi.Duration = now.Sub(vi.ArrivedAt)
}

// 维护相关值对象

// MaintenanceType 维护类型
type MaintenanceType int32

const (
	MaintenanceTypeRoutine   MaintenanceType = iota + 1 // 常规维护
	MaintenanceTypePreventive                            // 预防性维护
	MaintenanceTypeEmergency                             // 紧急维护
	MaintenanceTypeRepair                                // 修理维护
	MaintenanceTypeUpgrade                               // 升级维护
	MaintenanceTypeCleaning                              // 清洁维护
	MaintenanceTypeInspection                            // 检查维护
)

// String 返回维护类型的字符串表示
func (mt MaintenanceType) String() string {
	switch mt {
	case MaintenanceTypeRoutine:
		return "routine"
	case MaintenanceTypePreventive:
		return "preventive"
	case MaintenanceTypeEmergency:
		return "emergency"
	case MaintenanceTypeRepair:
		return "repair"
	case MaintenanceTypeUpgrade:
		return "upgrade"
	case MaintenanceTypeCleaning:
		return "cleaning"
	case MaintenanceTypeInspection:
		return "inspection"
	default:
		return "unknown"
	}
}

// IsValid 检查维护类型是否有效
func (mt MaintenanceType) IsValid() bool {
	return mt >= MaintenanceTypeRoutine && mt <= MaintenanceTypeInspection
}

// MaintenanceInfo 维护信息
type MaintenanceInfo struct {
	LastMaintenanceAt *time.Time            `json:"last_maintenance_at,omitempty" bson:"last_maintenance_at,omitempty"`
	NextMaintenanceAt time.Time             `json:"next_maintenance_at" bson:"next_maintenance_at"`
	MaintenanceLevel  int32                 `json:"maintenance_level" bson:"maintenance_level"`
	Costs             []*ResourceCost       `json:"costs" bson:"costs"`
	History           []*MaintenanceRecord  `json:"history" bson:"history"`
	CreatedAt         time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at" bson:"updated_at"`
}

// NewMaintenanceInfo 创建新维护信息
func NewMaintenanceInfo() *MaintenanceInfo {
	now := time.Now()
	return &MaintenanceInfo{
		NextMaintenanceAt: now.Add(24 * time.Hour),
		MaintenanceLevel:  100,
		Costs:             make([]*ResourceCost, 0),
		History:           make([]*MaintenanceRecord, 0),
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// MaintenanceRecord 维护记录
type MaintenanceRecord struct {
	Type        MaintenanceType `json:"type" bson:"type"`
	PerformedAt time.Time       `json:"performed_at" bson:"performed_at"`
	Costs       []*ResourceCost `json:"costs" bson:"costs"`
	Result      string          `json:"result" bson:"result"`
	Notes       string          `json:"notes" bson:"notes"`
}

// NewMaintenanceRecord 创建新维护记录
func NewMaintenanceRecord(maintenanceType MaintenanceType) *MaintenanceRecord {
	return &MaintenanceRecord{
		Type:        maintenanceType,
		PerformedAt: time.Now(),
		Costs:       make([]*ResourceCost, 0),
		Result:      "success",
		Notes:       "",
	}
}

// 辅助函数

// abs 计算绝对值
func abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}