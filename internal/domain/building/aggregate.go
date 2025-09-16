package building

import (
	"fmt"
	"time"
)

// BuildingAggregate 建筑聚合根
type BuildingAggregate struct {
	ID               string                 `json:"id" bson:"_id"`
	PlayerID         uint64                 `json:"player_id" bson:"player_id"`
	BuildingTypeID   string                 `json:"building_type_id" bson:"building_type_id"`
	Name             string                 `json:"name" bson:"name"`
	Description      string                 `json:"description" bson:"description"`
	Level            int32                  `json:"level" bson:"level"`
	MaxLevel         int32                  `json:"max_level" bson:"max_level"`
	Status           BuildingStatus         `json:"status" bson:"status"`
	Category         BuildingCategory       `json:"category" bson:"category"`
	Position         *Position              `json:"position,omitempty" bson:"position,omitempty"`
	Size             *Size                  `json:"size" bson:"size"`
	Orientation      Orientation            `json:"orientation" bson:"orientation"`
	Health           int32                  `json:"health" bson:"health"`
	MaxHealth        int32                  `json:"max_health" bson:"max_health"`
	Durability       int32                  `json:"durability" bson:"durability"`
	MaxDurability    int32                  `json:"max_durability" bson:"max_durability"`
	Effects          []*BuildingEffect      `json:"effects" bson:"effects"`
	Requirements     []*Requirement         `json:"requirements" bson:"requirements"`
	UpgradeCosts     []*ResourceCost        `json:"upgrade_costs" bson:"upgrade_costs"`
	MaintenanceCosts []*ResourceCost        `json:"maintenance_costs" bson:"maintenance_costs"`
	Production       *ProductionInfo        `json:"production,omitempty" bson:"production,omitempty"`
	Storage          *StorageInfo           `json:"storage,omitempty" bson:"storage,omitempty"`
	Defense          *DefenseInfo           `json:"defense,omitempty" bson:"defense,omitempty"`
	Construction     *ConstructionInfo      `json:"construction,omitempty" bson:"construction,omitempty"`
	Upgrade          *UpgradeInfo           `json:"upgrade,omitempty" bson:"upgrade,omitempty"`
	Maintenance      *MaintenanceInfo       `json:"maintenance,omitempty" bson:"maintenance,omitempty"`
	Workers          []*WorkerInfo          `json:"workers" bson:"workers"`
	Visitors         []*VisitorInfo         `json:"visitors" bson:"visitors"`
	Tags             []string               `json:"tags" bson:"tags"`
	Metadata         map[string]interface{} `json:"metadata" bson:"metadata"`
	LastActiveAt     time.Time              `json:"last_active_at" bson:"last_active_at"`
	CreatedAt        time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at" bson:"updated_at"`
}

// NewBuildingAggregate 创建新的建筑聚合
func NewBuildingAggregate(playerID uint64, buildingTypeID, name string, category BuildingCategory) *BuildingAggregate {
	now := time.Now()
	return &BuildingAggregate{
		ID:               generateBuildingID(),
		PlayerID:         playerID,
		BuildingTypeID:   buildingTypeID,
		Name:             name,
		Description:      "",
		Level:            1,
		MaxLevel:         10,
		Status:           BuildingStatusPlanning,
		Category:         category,
		Size:             &Size{Width: 1, Height: 1, Depth: 1},
		Orientation:      OrientationNorth,
		Health:           100,
		MaxHealth:        100,
		Durability:       100,
		MaxDurability:    100,
		Effects:          make([]*BuildingEffect, 0),
		Requirements:     make([]*Requirement, 0),
		UpgradeCosts:     make([]*ResourceCost, 0),
		MaintenanceCosts: make([]*ResourceCost, 0),
		Workers:          make([]*WorkerInfo, 0),
		Visitors:         make([]*VisitorInfo, 0),
		Tags:             make([]string, 0),
		Metadata:         make(map[string]interface{}),
		LastActiveAt:     now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// StartConstruction 开始建造
func (b *BuildingAggregate) StartConstruction(duration time.Duration, costs []*ResourceCost) error {
	if b.Status != BuildingStatusPlanning {
		return fmt.Errorf("building must be in planning status to start construction")
	}
	
	// 检查建造要求
	if err := b.checkRequirements(); err != nil {
		return fmt.Errorf("construction requirements not met: %w", err)
	}
	
	// 创建建造信息
	now := time.Now()
	b.Construction = &ConstructionInfo{
		StartedAt:    now,
		Duration:     duration,
		CompletedAt:  nil,
		Progress:     0.0,
		Costs:        costs,
		Workers:      make([]*WorkerAssignment, 0),
		Materials:    make([]*MaterialUsage, 0),
		Status:       ConstructionStatusInProgress,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	
	b.Status = BuildingStatusUnderConstruction
	b.UpdatedAt = now
	return nil
}

// UpdateConstructionProgress 更新建造进度
func (b *BuildingAggregate) UpdateConstructionProgress(progress float64) error {
	if b.Status != BuildingStatusUnderConstruction {
		return fmt.Errorf("building is not under construction")
	}
	
	if b.Construction == nil {
		return fmt.Errorf("construction info not found")
	}
	
	if progress < 0 || progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100")
	}
	
	b.Construction.Progress = progress
	b.Construction.UpdatedAt = time.Now()
	b.UpdatedAt = time.Now()
	
	// 检查是否完成建造
	if progress >= 100 {
		return b.CompleteConstruction()
	}
	
	return nil
}

// CompleteConstruction 完成建造
func (b *BuildingAggregate) CompleteConstruction() error {
	if b.Status != BuildingStatusUnderConstruction {
		return fmt.Errorf("building is not under construction")
	}
	
	if b.Construction == nil {
		return fmt.Errorf("construction info not found")
	}
	
	now := time.Now()
	b.Construction.CompletedAt = &now
	b.Construction.Progress = 100.0
	b.Construction.Status = ConstructionStatusCompleted
	b.Construction.UpdatedAt = now
	
	b.Status = BuildingStatusActive
	b.UpdatedAt = now
	return nil
}

// CancelConstruction 取消建造
func (b *BuildingAggregate) CancelConstruction(reason string) error {
	if b.Status != BuildingStatusUnderConstruction {
		return fmt.Errorf("building is not under construction")
	}
	
	if b.Construction == nil {
		return fmt.Errorf("construction info not found")
	}
	
	now := time.Now()
	b.Construction.Status = ConstructionStatusCancelled
	b.Construction.UpdatedAt = now
	b.Construction.SetMetadata("cancel_reason", reason)
	
	b.Status = BuildingStatusCancelled
	b.UpdatedAt = now
	return nil
}

// StartUpgrade 开始升级
func (b *BuildingAggregate) StartUpgrade(targetLevel int32, duration time.Duration, costs []*ResourceCost) error {
	if b.Status != BuildingStatusActive {
		return fmt.Errorf("building must be active to start upgrade")
	}
	
	if targetLevel <= b.Level {
		return fmt.Errorf("target level must be higher than current level")
	}
	
	if targetLevel > b.MaxLevel {
		return fmt.Errorf("target level exceeds maximum level")
	}
	
	// 检查升级要求
	if err := b.checkUpgradeRequirements(targetLevel); err != nil {
		return fmt.Errorf("upgrade requirements not met: %w", err)
	}
	
	// 创建升级信息
	now := time.Now()
	b.Upgrade = &UpgradeInfo{
		FromLevel:    b.Level,
		ToLevel:      targetLevel,
		StartedAt:    now,
		Duration:     duration,
		CompletedAt:  nil,
		Progress:     0.0,
		Costs:        costs,
		Workers:      make([]*WorkerAssignment, 0),
		Materials:    make([]*MaterialUsage, 0),
		Status:       UpgradeStatusInProgress,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	
	b.Status = BuildingStatusUpgrading
	b.UpdatedAt = now
	return nil
}

// UpdateUpgradeProgress 更新升级进度
func (b *BuildingAggregate) UpdateUpgradeProgress(progress float64) error {
	if b.Status != BuildingStatusUpgrading {
		return fmt.Errorf("building is not upgrading")
	}
	
	if b.Upgrade == nil {
		return fmt.Errorf("upgrade info not found")
	}
	
	if progress < 0 || progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100")
	}
	
	b.Upgrade.Progress = progress
	b.Upgrade.UpdatedAt = time.Now()
	b.UpdatedAt = time.Now()
	
	// 检查是否完成升级
	if progress >= 100 {
		return b.CompleteUpgrade()
	}
	
	return nil
}

// CompleteUpgrade 完成升级
func (b *BuildingAggregate) CompleteUpgrade() error {
	if b.Status != BuildingStatusUpgrading {
		return fmt.Errorf("building is not upgrading")
	}
	
	if b.Upgrade == nil {
		return fmt.Errorf("upgrade info not found")
	}
	
	now := time.Now()
	b.Upgrade.CompletedAt = &now
	b.Upgrade.Progress = 100.0
	b.Upgrade.Status = UpgradeStatusCompleted
	b.Upgrade.UpdatedAt = now
	
	// 更新建筑等级
	b.Level = b.Upgrade.ToLevel
	
	// 应用升级效果
	b.applyUpgradeEffects()
	
	b.Status = BuildingStatusActive
	b.UpdatedAt = now
	return nil
}

// CancelUpgrade 取消升级
func (b *BuildingAggregate) CancelUpgrade(reason string) error {
	if b.Status != BuildingStatusUpgrading {
		return fmt.Errorf("building is not upgrading")
	}
	
	if b.Upgrade == nil {
		return fmt.Errorf("upgrade info not found")
	}
	
	now := time.Now()
	b.Upgrade.Status = UpgradeStatusCancelled
	b.Upgrade.UpdatedAt = now
	b.Upgrade.SetMetadata("cancel_reason", reason)
	
	b.Status = BuildingStatusActive
	b.UpdatedAt = now
	return nil
}

// Repair 修理建筑
func (b *BuildingAggregate) Repair(amount int32, costs []*ResourceCost) error {
	if b.Status != BuildingStatusActive && b.Status != BuildingStatusDamaged {
		return fmt.Errorf("building cannot be repaired in current status")
	}
	
	if b.Health >= b.MaxHealth {
		return fmt.Errorf("building is already at full health")
	}
	
	// 修理建筑
	b.Health += amount
	if b.Health > b.MaxHealth {
		b.Health = b.MaxHealth
	}
	
	// 如果完全修复，更新状态
	if b.Health >= b.MaxHealth && b.Status == BuildingStatusDamaged {
		b.Status = BuildingStatusActive
	}
	
	b.UpdatedAt = time.Now()
	return nil
}

// TakeDamage 受到伤害
func (b *BuildingAggregate) TakeDamage(damage int32, damageType DamageType) error {
	if b.Status == BuildingStatusDestroyed {
		return fmt.Errorf("building is already destroyed")
	}
	
	// 计算实际伤害（考虑防御）
	actualDamage := b.calculateActualDamage(damage, damageType)
	
	// 扣除生命值
	b.Health -= actualDamage
	if b.Health < 0 {
		b.Health = 0
	}
	
	// 扣除耐久度
	b.Durability -= actualDamage / 2
	if b.Durability < 0 {
		b.Durability = 0
	}
	
	// 更新状态
	if b.Health <= 0 {
		b.Status = BuildingStatusDestroyed
	} else if b.Health < b.MaxHealth/2 {
		b.Status = BuildingStatusDamaged
	}
	
	b.UpdatedAt = time.Now()
	return nil
}

// Demolish 拆除建筑
func (b *BuildingAggregate) Demolish(reason string) error {
	if b.Status == BuildingStatusDestroyed {
		return fmt.Errorf("building is already destroyed")
	}
	
	// 停止所有进行中的操作
	if b.Construction != nil && b.Construction.Status == ConstructionStatusInProgress {
		b.Construction.Status = ConstructionStatusCancelled
	}
	
	if b.Upgrade != nil && b.Upgrade.Status == UpgradeStatusInProgress {
		b.Upgrade.Status = UpgradeStatusCancelled
	}
	
	// 清空工人和访客
	b.Workers = make([]*WorkerInfo, 0)
	b.Visitors = make([]*VisitorInfo, 0)
	
	b.Status = BuildingStatusDemolished
	b.SetMetadata("demolish_reason", reason)
	b.UpdatedAt = time.Now()
	return nil
}

// SetPosition 设置位置
func (b *BuildingAggregate) SetPosition(position *Position) error {
	if position == nil {
		return fmt.Errorf("position cannot be nil")
	}
	
	if err := position.Validate(); err != nil {
		return fmt.Errorf("invalid position: %w", err)
	}
	
	b.Position = position
	b.UpdatedAt = time.Now()
	return nil
}

// SetOrientation 设置朝向
func (b *BuildingAggregate) SetOrientation(orientation Orientation) error {
	if !orientation.IsValid() {
		return fmt.Errorf("invalid orientation: %v", orientation)
	}
	
	b.Orientation = orientation
	b.UpdatedAt = time.Now()
	return nil
}

// AddWorker 添加工人
func (b *BuildingAggregate) AddWorker(workerID uint64, role WorkerRole, efficiency float64) error {
	if b.Status != BuildingStatusActive {
		return fmt.Errorf("building must be active to add workers")
	}
	
	// 检查工人是否已存在
	for _, worker := range b.Workers {
		if worker.WorkerID == workerID {
			return fmt.Errorf("worker %d is already assigned to this building", workerID)
		}
	}
	
	// 检查工人容量
	maxWorkers := b.getMaxWorkers()
	if len(b.Workers) >= maxWorkers {
		return fmt.Errorf("building has reached maximum worker capacity: %d", maxWorkers)
	}
	
	// 添加工人
	now := time.Now()
	worker := &WorkerInfo{
		WorkerID:   workerID,
		Role:       role,
		Efficiency: efficiency,
		AssignedAt: now,
		Status:     WorkerStatusActive,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	
	b.Workers = append(b.Workers, worker)
	b.UpdatedAt = now
	return nil
}

// RemoveWorker 移除工人
func (b *BuildingAggregate) RemoveWorker(workerID uint64, reason string) error {
	for i, worker := range b.Workers {
		if worker.WorkerID == workerID {
			// 移除工人
			b.Workers = append(b.Workers[:i], b.Workers[i+1:]...)
			b.UpdatedAt = time.Now()
			return nil
		}
	}
	
	return fmt.Errorf("worker %d not found in building", workerID)
}

// AddEffect 添加效果
func (b *BuildingAggregate) AddEffect(effect *BuildingEffect) error {
	if effect == nil {
		return fmt.Errorf("effect cannot be nil")
	}
	
	if err := effect.Validate(); err != nil {
		return fmt.Errorf("invalid effect: %w", err)
	}
	
	// 检查是否已存在相同效果
	for _, existing := range b.Effects {
		if existing.Type == effect.Type && existing.Target == effect.Target {
			// 更新现有效果
			existing.Value = effect.Value
			existing.Duration = effect.Duration
			existing.UpdatedAt = time.Now()
			b.UpdatedAt = time.Now()
			return nil
		}
	}
	
	// 添加新效果
	effect.CreatedAt = time.Now()
	effect.UpdatedAt = time.Now()
	b.Effects = append(b.Effects, effect)
	b.UpdatedAt = time.Now()
	return nil
}

// RemoveEffect 移除效果
func (b *BuildingAggregate) RemoveEffect(effectType EffectType, target string) error {
	for i, effect := range b.Effects {
		if effect.Type == effectType && effect.Target == target {
			// 移除效果
			b.Effects = append(b.Effects[:i], b.Effects[i+1:]...)
			b.UpdatedAt = time.Now()
			return nil
		}
	}
	
	return fmt.Errorf("effect not found: type=%v, target=%s", effectType, target)
}

// UpdateProduction 更新生产信息
func (b *BuildingAggregate) UpdateProduction(production *ProductionInfo) error {
	if production == nil {
		return fmt.Errorf("production info cannot be nil")
	}
	
	if err := production.Validate(); err != nil {
		return fmt.Errorf("invalid production info: %w", err)
	}
	
	b.Production = production
	b.UpdatedAt = time.Now()
	return nil
}

// UpdateStorage 更新存储信息
func (b *BuildingAggregate) UpdateStorage(storage *StorageInfo) error {
	if storage == nil {
		return fmt.Errorf("storage info cannot be nil")
	}
	
	if err := storage.Validate(); err != nil {
		return fmt.Errorf("invalid storage info: %w", err)
	}
	
	b.Storage = storage
	b.UpdatedAt = time.Now()
	return nil
}

// UpdateDefense 更新防御信息
func (b *BuildingAggregate) UpdateDefense(defense *DefenseInfo) error {
	if defense == nil {
		return fmt.Errorf("defense info cannot be nil")
	}
	
	if err := defense.Validate(); err != nil {
		return fmt.Errorf("invalid defense info: %w", err)
	}
	
	b.Defense = defense
	b.UpdatedAt = time.Now()
	return nil
}

// PerformMaintenance 执行维护
func (b *BuildingAggregate) PerformMaintenance(maintenanceType MaintenanceType, costs []*ResourceCost) error {
	if b.Status != BuildingStatusActive {
		return fmt.Errorf("building must be active to perform maintenance")
	}
	
	// 创建或更新维护信息
	if b.Maintenance == nil {
		now := time.Now()
		b.Maintenance = &MaintenanceInfo{
			LastMaintenanceAt: &now,
			NextMaintenanceAt: now.Add(24 * time.Hour), // 默认24小时后需要维护
			MaintenanceLevel:  100,
			Costs:             make([]*ResourceCost, 0),
			History:           make([]*MaintenanceRecord, 0),
			CreatedAt:         now,
			UpdatedAt:         now,
		}
	}
	
	// 执行维护
	now := time.Now()
	record := &MaintenanceRecord{
		Type:        maintenanceType,
		PerformedAt: now,
		Costs:       costs,
		Result:      "success",
		Notes:       fmt.Sprintf("Performed %s maintenance", maintenanceType.String()),
	}
	
	b.Maintenance.History = append(b.Maintenance.History, record)
	b.Maintenance.LastMaintenanceAt = &now
	b.Maintenance.NextMaintenanceAt = now.Add(24 * time.Hour)
	b.Maintenance.MaintenanceLevel = 100 // 重置维护等级
	b.Maintenance.UpdatedAt = now
	
	// 恢复耐久度
	if b.Durability < b.MaxDurability {
		b.Durability += int32(float64(b.MaxDurability) * 0.1) // 恢复10%耐久度
		if b.Durability > b.MaxDurability {
			b.Durability = b.MaxDurability
		}
	}
	
	b.UpdatedAt = now
	return nil
}

// SetMetadata 设置元数据
func (b *BuildingAggregate) SetMetadata(key string, value interface{}) {
	if b.Metadata == nil {
		b.Metadata = make(map[string]interface{})
	}
	b.Metadata[key] = value
	b.UpdatedAt = time.Now()
}

// GetMetadata 获取元数据
func (b *BuildingAggregate) GetMetadata(key string) (interface{}, bool) {
	if b.Metadata == nil {
		return nil, false
	}
	value, exists := b.Metadata[key]
	return value, exists
}

// AddTag 添加标签
func (b *BuildingAggregate) AddTag(tag string) {
	// 检查标签是否已存在
	for _, existing := range b.Tags {
		if existing == tag {
			return // 已存在，不重复添加
		}
	}
	
	b.Tags = append(b.Tags, tag)
	b.UpdatedAt = time.Now()
}

// RemoveTag 移除标签
func (b *BuildingAggregate) RemoveTag(tag string) {
	for i, existing := range b.Tags {
		if existing == tag {
			b.Tags = append(b.Tags[:i], b.Tags[i+1:]...)
			b.UpdatedAt = time.Now()
			return
		}
	}
}

// HasTag 检查是否有标签
func (b *BuildingAggregate) HasTag(tag string) bool {
	for _, existing := range b.Tags {
		if existing == tag {
			return true
		}
	}
	return false
}

// 查询方法

// IsActive 检查是否活跃
func (b *BuildingAggregate) IsActive() bool {
	return b.Status == BuildingStatusActive
}

// IsUnderConstruction 检查是否在建造中
func (b *BuildingAggregate) IsUnderConstruction() bool {
	return b.Status == BuildingStatusUnderConstruction
}

// IsUpgrading 检查是否在升级中
func (b *BuildingAggregate) IsUpgrading() bool {
	return b.Status == BuildingStatusUpgrading
}

// IsDamaged 检查是否受损
func (b *BuildingAggregate) IsDamaged() bool {
	return b.Status == BuildingStatusDamaged
}

// IsDestroyed 检查是否被摧毁
func (b *BuildingAggregate) IsDestroyed() bool {
	return b.Status == BuildingStatusDestroyed
}

// CanUpgrade 检查是否可以升级
func (b *BuildingAggregate) CanUpgrade() bool {
	return b.Status == BuildingStatusActive && b.Level < b.MaxLevel
}

// CanRepair 检查是否可以修理
func (b *BuildingAggregate) CanRepair() bool {
	return (b.Status == BuildingStatusActive || b.Status == BuildingStatusDamaged) && b.Health < b.MaxHealth
}

// NeedsMaintenance 检查是否需要维护
func (b *BuildingAggregate) NeedsMaintenance() bool {
	if b.Maintenance == nil {
		return true
	}
	return time.Now().After(b.Maintenance.NextMaintenanceAt)
}

// GetEfficiency 获取效率
func (b *BuildingAggregate) GetEfficiency() float64 {
	if b.Status != BuildingStatusActive {
		return 0.0
	}
	
	// 基础效率
	baseEfficiency := 1.0
	
	// 健康度影响
	healthFactor := float64(b.Health) / float64(b.MaxHealth)
	
	// 耐久度影响
	durabilityFactor := float64(b.Durability) / float64(b.MaxDurability)
	
	// 维护影响
	maintenanceFactor := 1.0
	if b.Maintenance != nil {
		maintenanceFactor = float64(b.Maintenance.MaintenanceLevel) / 100.0
	}
	
	// 工人效率影响
	workerFactor := b.getWorkerEfficiencyFactor()
	
	// 效果影响
	effectFactor := b.getEffectFactor(EffectTypeEfficiency)
	
	return baseEfficiency * healthFactor * durabilityFactor * maintenanceFactor * workerFactor * effectFactor
}

// GetTotalUpgradeCost 获取升级总成本
func (b *BuildingAggregate) GetTotalUpgradeCost(targetLevel int32) []*ResourceCost {
	totalCosts := make(map[string]int64)
	
	for level := b.Level + 1; level <= targetLevel; level++ {
		levelCosts := b.getUpgradeCostForLevel(level)
		for _, cost := range levelCosts {
			totalCosts[cost.ResourceType] += cost.Amount
		}
	}
	
	result := make([]*ResourceCost, 0, len(totalCosts))
	for resourceType, amount := range totalCosts {
		result = append(result, &ResourceCost{
			ResourceType: resourceType,
			Amount:       amount,
		})
	}
	
	return result
}

// GetOccupiedArea 获取占用面积
func (b *BuildingAggregate) GetOccupiedArea() int32 {
	if b.Size == nil {
		return 1
	}
	return b.Size.Width * b.Size.Height
}

// GetBoundingBox 获取边界框
func (b *BuildingAggregate) GetBoundingBox() *BoundingBox {
	if b.Position == nil || b.Size == nil {
		return nil
	}
	
	return &BoundingBox{
		MinX: b.Position.X,
		MinY: b.Position.Y,
		MinZ: b.Position.Z,
		MaxX: b.Position.X + b.Size.Width - 1,
		MaxY: b.Position.Y + b.Size.Height - 1,
		MaxZ: b.Position.Z + b.Size.Depth - 1,
	}
}

// 私有方法

// checkRequirements 检查建造要求
func (b *BuildingAggregate) checkRequirements() error {
	for _, req := range b.Requirements {
		if !req.IsMet() {
			return fmt.Errorf("requirement not met: %s", req.Description)
		}
	}
	return nil
}

// checkUpgradeRequirements 检查升级要求
func (b *BuildingAggregate) checkUpgradeRequirements(targetLevel int32) error {
	// 检查基础要求
	if err := b.checkRequirements(); err != nil {
		return err
	}
	
	// 检查等级特定要求
	// TODO: 实现等级特定要求检查
	
	return nil
}

// applyUpgradeEffects 应用升级效果
func (b *BuildingAggregate) applyUpgradeEffects() {
	// 提升最大生命值
	b.MaxHealth = int32(float64(b.MaxHealth) * 1.1)
	b.Health = b.MaxHealth
	
	// 提升最大耐久度
	b.MaxDurability = int32(float64(b.MaxDurability) * 1.1)
	b.Durability = b.MaxDurability
	
	// 应用等级相关的效果
	// TODO: 根据建筑类型和等级应用特定效果
}

// calculateActualDamage 计算实际伤害
func (b *BuildingAggregate) calculateActualDamage(damage int32, damageType DamageType) int32 {
	actualDamage := damage
	
	// 应用防御
	if b.Defense != nil {
		defenseValue := b.Defense.GetDefenseValue(damageType)
		actualDamage -= defenseValue
		if actualDamage < 0 {
			actualDamage = 0
		}
	}
	
	// 应用效果
	effectFactor := b.getEffectFactor(EffectTypeDefense)
	actualDamage = int32(float64(actualDamage) * (1.0 - effectFactor))
	
	return actualDamage
}

// getMaxWorkers 获取最大工人数
func (b *BuildingAggregate) getMaxWorkers() int {
	baseWorkers := 2
	levelBonus := int(b.Level - 1)
	sizeBonus := int(b.GetOccupiedArea() / 4)
	return baseWorkers + levelBonus + sizeBonus
}

// getWorkerEfficiencyFactor 获取工人效率因子
func (b *BuildingAggregate) getWorkerEfficiencyFactor() float64 {
	if len(b.Workers) == 0 {
		return 0.5 // 没有工人时效率降低
	}
	
	totalEfficiency := 0.0
	for _, worker := range b.Workers {
		if worker.Status == WorkerStatusActive {
			totalEfficiency += worker.Efficiency
		}
	}
	
	averageEfficiency := totalEfficiency / float64(len(b.Workers))
	return averageEfficiency
}

// getEffectFactor 获取效果因子
func (b *BuildingAggregate) getEffectFactor(effectType EffectType) float64 {
	factor := 0.0
	for _, effect := range b.Effects {
		if effect.Type == effectType && effect.IsActive() {
			factor += effect.Value
		}
	}
	return factor
}

// getUpgradeCostForLevel 获取指定等级的升级成本
func (b *BuildingAggregate) getUpgradeCostForLevel(level int32) []*ResourceCost {
	// 基础成本随等级增长
	baseCost := int64(100 * level * level)
	
	return []*ResourceCost{
		{ResourceType: "wood", Amount: baseCost},
		{ResourceType: "stone", Amount: baseCost / 2},
		{ResourceType: "metal", Amount: baseCost / 4},
	}
}

// generateBuildingID 生成建筑ID
func generateBuildingID() string {
	return fmt.Sprintf("building_%d", time.Now().UnixNano())
}

// 常量定义

const (
	// 建筑相关常量
	DefaultMaxLevel      = int32(10)  // 默认最大等级
	DefaultMaxHealth     = int32(100) // 默认最大生命值
	DefaultMaxDurability = int32(100) // 默认最大耐久度
	
	// 维护相关常量
	MaintenanceInterval = 24 * time.Hour // 维护间隔
	MaintenanceCost     = int64(50)      // 基础维护成本
	
	// 建造相关常量
	DefaultConstructionTime = 1 * time.Hour // 默认建造时间
	DefaultUpgradeTime      = 30 * time.Minute // 默认升级时间
	
	// 效率相关常量
	MinEfficiency     = 0.1 // 最小效率
	MaxEfficiency     = 2.0 // 最大效率
	BaseEfficiency    = 1.0 // 基础效率
	NoWorkerEfficiency = 0.5 // 无工人时的效率
)

// 验证函数

// ReconstructBuildingAggregate 从持久化数据重建建筑聚合根
func ReconstructBuildingAggregate(
	id string,
	playerID uint64,
	buildingTypeID string,
	name string,
	description string,
	level int32,
	maxLevel int32,
	status BuildingStatus,
	category BuildingCategory,
	position *Position,
	size *Size,
	orientation Orientation,
	health int32,
	maxHealth int32,
	durability int32,
	maxDurability int32,
	effects []*BuildingEffect,
	requirements []*Requirement,
	upgradeCosts []*ResourceCost,
	maintenanceCosts []*ResourceCost,
	tags []string,
	metadata map[string]interface{},
	lastActiveAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) *BuildingAggregate {
	return &BuildingAggregate{
		ID:               id,
		PlayerID:         playerID,
		BuildingTypeID:   buildingTypeID,
		Name:             name,
		Description:      description,
		Level:            level,
		MaxLevel:         maxLevel,
		Status:           status,
		Category:         category,
		Position:         position,
		Size:             size,
		Orientation:      orientation,
		Health:           health,
		MaxHealth:        maxHealth,
		Durability:       durability,
		MaxDurability:    maxDurability,
		Effects:          effects,
		Requirements:     requirements,
		UpgradeCosts:     upgradeCosts,
		MaintenanceCosts: maintenanceCosts,
		Workers:          make([]*WorkerInfo, 0),
		Visitors:         make([]*VisitorInfo, 0),
		Tags:             tags,
		Metadata:         metadata,
		LastActiveAt:     lastActiveAt,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}

// Validate 验证建筑聚合
func (b *BuildingAggregate) Validate() error {
	if b.ID == "" {
		return fmt.Errorf("building ID cannot be empty")
	}
	
	if b.PlayerID == 0 {
		return fmt.Errorf("player ID cannot be zero")
	}
	
	if b.BuildingTypeID == "" {
		return fmt.Errorf("building type ID cannot be empty")
	}
	
	if b.Name == "" {
		return fmt.Errorf("building name cannot be empty")
	}
	
	if !b.Status.IsValid() {
		return fmt.Errorf("invalid building status: %v", b.Status)
	}
	
	if !b.Category.IsValid() {
		return fmt.Errorf("invalid building category: %v", b.Category)
	}
	
	if b.Level < 1 {
		return fmt.Errorf("building level must be at least 1")
	}
	
	if b.Level > b.MaxLevel {
		return fmt.Errorf("building level cannot exceed max level")
	}
	
	if b.Health < 0 || b.Health > b.MaxHealth {
		return fmt.Errorf("building health must be between 0 and max health")
	}
	
	if b.Durability < 0 || b.Durability > b.MaxDurability {
		return fmt.Errorf("building durability must be between 0 and max durability")
	}
	
	if b.Size != nil {
		if err := b.Size.Validate(); err != nil {
			return fmt.Errorf("invalid building size: %w", err)
		}
	}
	
	if b.Position != nil {
		if err := b.Position.Validate(); err != nil {
			return fmt.Errorf("invalid building position: %w", err)
		}
	}
	
	if !b.Orientation.IsValid() {
		return fmt.Errorf("invalid building orientation: %v", b.Orientation)
	}
	
	// 验证效果
	for _, effect := range b.Effects {
		if err := effect.Validate(); err != nil {
			return fmt.Errorf("invalid building effect: %w", err)
		}
	}
	
	// 验证要求
	for _, req := range b.Requirements {
		if err := req.Validate(); err != nil {
			return fmt.Errorf("invalid building requirement: %w", err)
		}
	}
	
	// 验证成本
	for _, cost := range b.UpgradeCosts {
		if err := cost.Validate(); err != nil {
			return fmt.Errorf("invalid upgrade cost: %w", err)
		}
	}
	
	for _, cost := range b.MaintenanceCosts {
		if err := cost.Validate(); err != nil {
			return fmt.Errorf("invalid maintenance cost: %w", err)
		}
	}
	
	// 验证生产信息
	if b.Production != nil {
		if err := b.Production.Validate(); err != nil {
			return fmt.Errorf("invalid production info: %w", err)
		}
	}
	
	// 验证存储信息
	if b.Storage != nil {
		if err := b.Storage.Validate(); err != nil {
			return fmt.Errorf("invalid storage info: %w", err)
		}
	}
	
	// 验证防御信息
	if b.Defense != nil {
		if err := b.Defense.Validate(); err != nil {
			return fmt.Errorf("invalid defense info: %w", err)
		}
	}
	
	return nil
}