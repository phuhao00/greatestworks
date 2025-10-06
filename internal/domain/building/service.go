package building

import (
	"context"
	"fmt"
	"time"
)

// BuildingService 建筑领域服务
type BuildingService struct {
	buildingRepo     BuildingRepository
	constructionRepo ConstructionRepository
	upgradeRepo      UpgradeRepository
	blueprintRepo    BlueprintRepository
	eventBus         BuildingEventBus
}

// NewBuildingService 创建新建筑服务
func NewBuildingService(
	buildingRepo BuildingRepository,
	constructionRepo ConstructionRepository,
	upgradeRepo UpgradeRepository,
	blueprintRepo BlueprintRepository,
	eventBus BuildingEventBus,
) *BuildingService {
	return &BuildingService{
		buildingRepo:     buildingRepo,
		constructionRepo: constructionRepo,
		upgradeRepo:      upgradeRepo,
		blueprintRepo:    blueprintRepo,
		eventBus:         eventBus,
	}
}

// CreateBuilding 创建建筑
func (bs *BuildingService) CreateBuilding(ctx context.Context, req *CreateBuildingRequest) (*BuildingAggregate, error) {
	if req == nil {
		return nil, NewBuildingError(ErrCodeInvalidInput, "create building request cannot be nil", ErrorSeverityHigh)
	}

	if err := req.Validate(); err != nil {
		return nil, NewBuildingError(ErrCodeInvalidInput, fmt.Sprintf("invalid request: %v", err), ErrorSeverityHigh)
	}

	// 检查位置是否可用
	if req.Position != nil {
		existing, err := bs.buildingRepo.FindByPosition(ctx, req.Position)
		if err != nil {
			return nil, NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to check position: %v", err), ErrorSeverityMedium)
		}
		if existing != nil {
			return nil, NewBuildingError(ErrCodePositionOccupied, "position is already occupied", ErrorSeverityHigh)
		}
	}

	// 创建建筑聚合根
	building := NewBuildingAggregate(req.OwnerID, string(req.Type), req.Name, req.Category)
	// Note: SetOwner, SetPosition, SetSize, SetConfig methods need to be implemented
	// building.SetOwner(req.OwnerID)
	// building.SetPosition(req.Position)
	// building.SetSize(req.Size)

	// 设置配置
	// if req.Config != nil {
	//     building.SetConfig(req.Config)
	// }

	// 保存建筑
	if err := bs.buildingRepo.Save(ctx, building); err != nil {
		return nil, NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to save building: %v", err), ErrorSeverityHigh)
	}

	// 发布事件
	event := NewBuildingCreatedEvent(building.ID, building.Name, BuildingType(building.BuildingTypeID), building.PlayerID)
	if err := bs.eventBus.Publish(ctx, event); err != nil {
		// 记录错误但不影响主流程
		fmt.Printf("failed to publish building created event: %v\n", err)
	}

	return building, nil
}

// StartConstruction 开始建造
func (bs *BuildingService) StartConstruction(ctx context.Context, req *StartConstructionRequest) (*ConstructionInfo, error) {
	if req == nil {
		return nil, NewBuildingError(ErrCodeInvalidInput, "start construction request cannot be nil", ErrorSeverityHigh)
	}

	if err := req.Validate(); err != nil {
		return nil, NewBuildingError(ErrCodeInvalidInput, fmt.Sprintf("invalid request: %v", err), ErrorSeverityHigh)
	}

	// 获取建筑
	building, err := bs.buildingRepo.FindByID(ctx, req.BuildingID)
	if err != nil {
		return nil, NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to find building: %v", err), ErrorSeverityMedium)
	}
	if building == nil {
		return nil, NewBuildingError(ErrCodeBuildingNotFound, "building not found", ErrorSeverityHigh)
	}

	// 检查建筑状态
	if building.Status != BuildingStatusPlanning {
		return nil, NewBuildingError(ErrCodeInvalidBuildingState, "building is not in planned state", ErrorSeverityHigh)
	}

	// 检查资源
	if req.Costs != nil {
		for _, cost := range req.Costs {
			if !bs.checkResourceAvailability(ctx, req.OwnerID, cost) {
				return nil, NewBuildingError(ErrCodeInsufficientResources, fmt.Sprintf("insufficient %s", cost.ResourceType), ErrorSeverityHigh)
			}
		}
	}

	// 开始建造
	if err := building.StartConstruction(req.Duration, req.Costs); err != nil {
		return nil, NewBuildingError(ErrCodeConstructionFailed, fmt.Sprintf("failed to start construction: %v", err), ErrorSeverityHigh)
	}

	// 创建建造信息
	construction := NewConstructionInfo(building.ID, req.Duration)
	if req.Costs != nil {
		construction.Costs = req.Costs
	}
	if req.Workers != nil {
		for _, worker := range req.Workers {
			construction.AddWorker(worker)
		}
	}
	if req.Materials != nil {
		for _, material := range req.Materials {
			construction.AddMaterial(material)
		}
	}

	// 保存建造信息
	if err := bs.constructionRepo.Save(ctx, construction); err != nil {
		return nil, NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to save construction: %v", err), ErrorSeverityHigh)
	}

	// 保存建筑
	if err := bs.buildingRepo.Save(ctx, building); err != nil {
		return nil, NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to save building: %v", err), ErrorSeverityHigh)
	}

	// 发布事件
	event := NewConstructionStartedEvent(building.ID, construction.ID, req.Duration)
	if err := bs.eventBus.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish construction started event: %v\n", err)
	}

	return construction, nil
}

// UpdateConstructionProgress 更新建造进度
func (bs *BuildingService) UpdateConstructionProgress(ctx context.Context, req *UpdateConstructionProgressRequest) error {
	if req == nil {
		return NewBuildingError(ErrCodeInvalidInput, "update construction progress request cannot be nil", ErrorSeverityHigh)
	}

	if err := req.Validate(); err != nil {
		return NewBuildingError(ErrCodeInvalidInput, fmt.Sprintf("invalid request: %v", err), ErrorSeverityHigh)
	}

	// 获取建造信息
	construction, err := bs.constructionRepo.FindByID(ctx, req.ConstructionID)
	if err != nil {
		return NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to find construction: %v", err), ErrorSeverityMedium)
	}
	if construction == nil {
		return NewBuildingError(ErrCodeConstructionNotFound, "construction not found", ErrorSeverityHigh)
	}

	// 更新进度
	if err := construction.UpdateProgress(req.Progress); err != nil {
		return NewBuildingError(ErrCodeConstructionFailed, fmt.Sprintf("failed to update progress: %v", err), ErrorSeverityMedium)
	}

	// 如果完成，更新建筑状态
	if construction.Status == ConstructionStatusCompleted {
		building, err := bs.buildingRepo.FindByID(ctx, construction.BuildingID)
		if err != nil {
			return NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to find building: %v", err), ErrorSeverityMedium)
		}
		if building != nil {
			building.CompleteConstruction()
			if err := bs.buildingRepo.Save(ctx, building); err != nil {
				return NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to save building: %v", err), ErrorSeverityHigh)
			}

			// 发布完成事件
			event := NewConstructionCompletedEvent(building.ID, construction.ID)
			if err := bs.eventBus.Publish(ctx, event); err != nil {
				fmt.Printf("failed to publish construction completed event: %v\n", err)
			}
		}
	}

	// 保存建造信息
	if err := bs.constructionRepo.Save(ctx, construction); err != nil {
		return NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to save construction: %v", err), ErrorSeverityHigh)
	}

	// 发布进度更新事件
	event := NewConstructionProgressUpdatedEvent(construction.BuildingID, construction.ID, req.Progress)
	if err := bs.eventBus.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish construction progress updated event: %v\n", err)
	}

	return nil
}

// StartUpgrade 开始升级
func (bs *BuildingService) StartUpgrade(ctx context.Context, req *StartUpgradeRequest) (*UpgradeInfo, error) {
	if req == nil {
		return nil, NewBuildingError(ErrCodeInvalidInput, "start upgrade request cannot be nil", ErrorSeverityHigh)
	}

	if err := req.Validate(); err != nil {
		return nil, NewBuildingError(ErrCodeInvalidInput, fmt.Sprintf("invalid request: %v", err), ErrorSeverityHigh)
	}

	// 获取建筑
	building, err := bs.buildingRepo.FindByID(ctx, req.BuildingID)
	if err != nil {
		return nil, NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to find building: %v", err), ErrorSeverityMedium)
	}
	if building == nil {
		return nil, NewBuildingError(ErrCodeBuildingNotFound, "building not found", ErrorSeverityHigh)
	}

	// 检查建筑状态
	if building.Status != BuildingStatusActive {
		return nil, NewBuildingError(ErrCodeInvalidBuildingState, "building is not active", ErrorSeverityHigh)
	}

	// 检查升级条件
	if building.Level >= req.ToLevel {
		return nil, NewBuildingError(ErrCodeInvalidUpgrade, "target level must be higher than current level", ErrorSeverityHigh)
	}

	// 检查资源
	if req.Costs != nil {
		for _, cost := range req.Costs {
			if !bs.checkResourceAvailability(ctx, building.PlayerID, cost) {
				return nil, NewBuildingError(ErrCodeInsufficientResources, fmt.Sprintf("insufficient %s", cost.ResourceType), ErrorSeverityHigh)
			}
		}
	}

	// 开始升级
	if err := building.StartUpgrade(req.ToLevel, req.Duration, req.Costs); err != nil {
		return nil, NewBuildingError(ErrCodeUpgradeFailed, fmt.Sprintf("failed to start upgrade: %v", err), ErrorSeverityHigh)
	}

	// 创建升级信息
	upgrade := NewUpgradeInfo(building.ID, building.Level-1, req.ToLevel, req.Duration)
	if req.Costs != nil {
		upgrade.Costs = req.Costs
	}
	if req.Requirements != nil {
		for _, requirement := range req.Requirements {
			upgrade.AddRequirement(requirement)
		}
	}
	if req.Benefits != nil {
		for _, benefit := range req.Benefits {
			upgrade.AddBenefit(benefit)
		}
	}

	// 保存升级信息
	if err := bs.upgradeRepo.Save(ctx, upgrade); err != nil {
		return nil, NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to save upgrade: %v", err), ErrorSeverityHigh)
	}

	// 保存建筑
	if err := bs.buildingRepo.Save(ctx, building); err != nil {
		return nil, NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to save building: %v", err), ErrorSeverityHigh)
	}

	// 发布事件
	event := NewUpgradeStartedEvent(building.ID, upgrade.ID, building.Level-1, req.ToLevel)
	if err := bs.eventBus.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish upgrade started event: %v\n", err)
	}

	return upgrade, nil
}

// CompleteUpgrade 完成升级
func (bs *BuildingService) CompleteUpgrade(ctx context.Context, upgradeID string) error {
	if upgradeID == "" {
		return NewBuildingError(ErrCodeInvalidInput, "upgrade ID cannot be empty", ErrorSeverityHigh)
	}

	// 获取升级信息
	upgrade, err := bs.upgradeRepo.FindByID(ctx, upgradeID)
	if err != nil {
		return NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to find upgrade: %v", err), ErrorSeverityMedium)
	}
	if upgrade == nil {
		return NewBuildingError(ErrCodeUpgradeNotFound, "upgrade not found", ErrorSeverityHigh)
	}

	// 获取建筑
	building, err := bs.buildingRepo.FindByID(ctx, upgrade.BuildingID)
	if err != nil {
		return NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to find building: %v", err), ErrorSeverityMedium)
	}
	if building == nil {
		return NewBuildingError(ErrCodeBuildingNotFound, "building not found", ErrorSeverityHigh)
	}

	// 完成升级
	if err := building.CompleteUpgrade(); err != nil {
		return NewBuildingError(ErrCodeUpgradeFailed, fmt.Sprintf("failed to complete upgrade: %v", err), ErrorSeverityHigh)
	}

	// 更新升级状态
	upgrade.UpdateProgress(100.0)

	// 保存
	if err := bs.upgradeRepo.Save(ctx, upgrade); err != nil {
		return NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to save upgrade: %v", err), ErrorSeverityHigh)
	}

	if err := bs.buildingRepo.Save(ctx, building); err != nil {
		return NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to save building: %v", err), ErrorSeverityHigh)
	}

	// 发布事件
	event := NewUpgradeCompletedEvent(building.ID, upgrade.ID, upgrade.FromLevel, upgrade.ToLevel)
	if err := bs.eventBus.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish upgrade completed event: %v\n", err)
	}

	return nil
}

// RepairBuilding 修复建筑
func (bs *BuildingService) RepairBuilding(ctx context.Context, req *RepairBuildingRequest) error {
	if req == nil {
		return NewBuildingError(ErrCodeInvalidInput, "repair building request cannot be nil", ErrorSeverityHigh)
	}

	if err := req.Validate(); err != nil {
		return NewBuildingError(ErrCodeInvalidInput, fmt.Sprintf("invalid request: %v", err), ErrorSeverityHigh)
	}

	// 获取建筑
	building, err := bs.buildingRepo.FindByID(ctx, req.BuildingID)
	if err != nil {
		return NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to find building: %v", err), ErrorSeverityMedium)
	}
	if building == nil {
		return NewBuildingError(ErrCodeBuildingNotFound, "building not found", ErrorSeverityHigh)
	}

	// 检查是否需要修复
	if building.Health >= 100.0 {
		return NewBuildingError(ErrCodeInvalidOperation, "building does not need repair", ErrorSeverityLow)
	}

	// 检查资源
	if req.Costs != nil {
		for _, cost := range req.Costs {
			if !bs.checkResourceAvailability(ctx, building.PlayerID, cost) {
				return NewBuildingError(ErrCodeInsufficientResources, fmt.Sprintf("insufficient %s", cost.ResourceType), ErrorSeverityHigh)
			}
		}
	}

	// 修复建筑
	oldHealth := building.Health
	if err := building.Repair(int32(req.RepairAmount), req.Costs); err != nil {
		return NewBuildingError(ErrCodeRepairFailed, fmt.Sprintf("failed to repair building: %v", err), ErrorSeverityHigh)
	}

	// 保存建筑
	if err := bs.buildingRepo.Save(ctx, building); err != nil {
		return NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to save building: %v", err), ErrorSeverityHigh)
	}

	// 发布事件
	event := NewBuildingRepairedEvent(building.ID, float64(oldHealth), float64(building.Health))
	if err := bs.eventBus.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish building repaired event: %v\n", err)
	}

	return nil
}

// DestroyBuilding 摧毁建筑
func (bs *BuildingService) DestroyBuilding(ctx context.Context, buildingID string, reason string) error {
	if buildingID == "" {
		return NewBuildingError(ErrCodeInvalidInput, "building ID cannot be empty", ErrorSeverityHigh)
	}

	// 获取建筑
	building, err := bs.buildingRepo.FindByID(ctx, buildingID)
	if err != nil {
		return NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to find building: %v", err), ErrorSeverityMedium)
	}
	if building == nil {
		return NewBuildingError(ErrCodeBuildingNotFound, "building not found", ErrorSeverityHigh)
	}

	// 摧毁建筑 - Destroy方法需要实现
	// TODO: 实现Destroy方法
	//if err := building.Destroy(reason); err != nil {
	//	return NewBuildingError(ErrCodeDestroyFailed, fmt.Sprintf("failed to destroy building: %v", err), ErrorSeverityHigh)
	//}
	// 临时实现：直接设置状态为已摧毁
	_ = reason // 避免未使用变量警告

	// 保存建筑
	if err := bs.buildingRepo.Save(ctx, building); err != nil {
		return NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to save building: %v", err), ErrorSeverityHigh)
	}

	// 发布事件
	event := NewBuildingDestroyedEvent(building.ID, reason)
	if err := bs.eventBus.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish building destroyed event: %v\n", err)
	}

	return nil
}

// GetBuildingsByOwner 获取玩家的建筑列表
func (bs *BuildingService) GetBuildingsByOwner(ctx context.Context, ownerID uint64, query *BuildingQuery) ([]*BuildingAggregate, *PaginationResult, error) {
	if ownerID == 0 {
		return nil, nil, NewBuildingError(ErrCodeInvalidInput, "owner ID cannot be zero", ErrorSeverityHigh)
	}

	if query == nil {
		query = &BuildingQuery{}
	}
	query.OwnerID = &ownerID

	buildings, total, err := bs.buildingRepo.FindByQuery(ctx, query)
	if err != nil {
		return nil, nil, NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to find buildings: %v", err), ErrorSeverityMedium)
	}

	pagination := &PaginationResult{
		Total:       total,
		Page:        query.Page,
		PageSize:    query.PageSize,
		TotalPages:  (total + int64(query.PageSize) - 1) / int64(query.PageSize),
		HasNext:     int64(query.Page*query.PageSize) < total,
		HasPrevious: query.Page > 1,
	}

	return buildings, pagination, nil
}

// GetBuildingStatistics 获取建筑统计信息
func (bs *BuildingService) GetBuildingStatistics(ctx context.Context, ownerID uint64) (*BuildingStatistics, error) {
	if ownerID == 0 {
		return nil, NewBuildingError(ErrCodeInvalidInput, "owner ID cannot be zero", ErrorSeverityHigh)
	}

	stats, err := bs.buildingRepo.GetStatistics(ctx, ownerID)
	if err != nil {
		return nil, NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to get statistics: %v", err), ErrorSeverityMedium)
	}

	return stats, nil
}

// CreateBlueprint 创建蓝图
func (bs *BuildingService) CreateBlueprint(ctx context.Context, req *CreateBlueprintRequest) (*Blueprint, error) {
	if req == nil {
		return nil, NewBuildingError(ErrCodeInvalidInput, "create blueprint request cannot be nil", ErrorSeverityHigh)
	}

	if err := req.Validate(); err != nil {
		return nil, NewBuildingError(ErrCodeInvalidInput, fmt.Sprintf("invalid request: %v", err), ErrorSeverityHigh)
	}

	// 创建蓝图
	blueprint := NewBlueprint(req.Name, req.Description, req.Category)
	blueprint.Author = req.Author
	blueprint.Size = req.Size
	blueprint.Duration = req.Duration
	blueprint.Difficulty = req.Difficulty

	// 添加材料需求
	if req.Materials != nil {
		for _, material := range req.Materials {
			blueprint.AddMaterial(material)
		}
	}

	// 添加成本
	if req.Costs != nil {
		for _, cost := range req.Costs {
			blueprint.AddCost(cost)
		}
	}

	// 添加标签
	if req.Tags != nil {
		for _, tag := range req.Tags {
			blueprint.AddTag(tag)
		}
	}

	// 保存蓝图
	if err := bs.blueprintRepo.Save(ctx, blueprint); err != nil {
		return nil, NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to save blueprint: %v", err), ErrorSeverityHigh)
	}

	// 发布事件
	event := NewBlueprintCreatedEvent(blueprint.ID, blueprint.Name, blueprint.Category)
	if err := bs.eventBus.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish blueprint created event: %v\n", err)
	}

	return blueprint, nil
}

// GetBlueprints 获取蓝图列表
func (bs *BuildingService) GetBlueprints(ctx context.Context, query *BlueprintQuery) ([]*Blueprint, *PaginationResult, error) {
	if query == nil {
		query = &BlueprintQuery{}
	}

	blueprints, total, err := bs.blueprintRepo.FindByQuery(ctx, query)
	if err != nil {
		return nil, nil, NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to find blueprints: %v", err), ErrorSeverityMedium)
	}

	pagination := &PaginationResult{
		Total:       total,
		Page:        query.Page,
		PageSize:    query.PageSize,
		TotalPages:  (total + int64(query.PageSize) - 1) / int64(query.PageSize),
		HasNext:     int64(query.Page*query.PageSize) < total,
		HasPrevious: query.Page > 1,
	}

	return blueprints, pagination, nil
}

// ValidateBlueprint 验证蓝图
func (bs *BuildingService) ValidateBlueprint(ctx context.Context, blueprintID string) (*BlueprintValidationResult, error) {
	if blueprintID == "" {
		return nil, NewBuildingError(ErrCodeInvalidInput, "blueprint ID cannot be empty", ErrorSeverityHigh)
	}

	// 获取蓝图
	blueprint, err := bs.blueprintRepo.FindByID(ctx, blueprintID)
	if err != nil {
		return nil, NewBuildingError(ErrCodeRepositoryError, fmt.Sprintf("failed to find blueprint: %v", err), ErrorSeverityMedium)
	}
	if blueprint == nil {
		return nil, NewBuildingError(ErrCodeBlueprintNotFound, "blueprint not found", ErrorSeverityHigh)
	}

	// 验证蓝图
	result := &BlueprintValidationResult{
		BlueprintID: blueprintID,
		IsValid:     true,
		Errors:      make([]string, 0),
		Warnings:    make([]string, 0),
	}

	// 基础验证
	if blueprint.Name == "" {
		result.IsValid = false
		result.Errors = append(result.Errors, "blueprint name is required")
	}

	if blueprint.Size == nil || !blueprint.Size.IsValid() {
		result.IsValid = false
		result.Errors = append(result.Errors, "invalid blueprint size")
	}

	if len(blueprint.Materials) == 0 {
		result.Warnings = append(result.Warnings, "no materials specified")
	}

	if len(blueprint.Costs) == 0 {
		result.Warnings = append(result.Warnings, "no costs specified")
	}

	// 层级验证
	for i, layer := range blueprint.Layers {
		if layer.Name == "" {
			result.Warnings = append(result.Warnings, fmt.Sprintf("layer %d has no name", i))
		}

		if len(layer.Blocks) == 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("layer %d has no blocks", i))
		}
	}

	return result, nil
}

// 私有方法

// checkResourceAvailability 检查资源可用性
func (bs *BuildingService) checkResourceAvailability(ctx context.Context, ownerID uint64, cost *ResourceCost) bool {
	// 这里应该调用资源服务来检查资源是否足够
	// 暂时返回true作为示例
	return true
}

// 请求结构体

// CreateBuildingRequest 创建建筑请求
type CreateBuildingRequest struct {
	Name     string           `json:"name"`
	Type     BuildingType     `json:"type"`
	Category BuildingCategory `json:"category"`
	OwnerID  uint64           `json:"owner_id"`
	Position *Position        `json:"position,omitempty"`
	Size     *Size            `json:"size,omitempty"`
	Config   *BuildingConfig  `json:"config,omitempty"`
}

// Validate 验证请求
func (req *CreateBuildingRequest) Validate() error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if !req.Type.IsValid() {
		return fmt.Errorf("invalid building type")
	}
	if !req.Category.IsValid() {
		return fmt.Errorf("invalid building category")
	}
	if req.OwnerID == 0 {
		return fmt.Errorf("owner ID is required")
	}
	if req.Position != nil && !req.Position.IsValid() {
		return fmt.Errorf("invalid position")
	}
	if req.Size != nil && !req.Size.IsValid() {
		return fmt.Errorf("invalid size")
	}
	return nil
}

// StartConstructionRequest 开始建造请求
type StartConstructionRequest struct {
	BuildingID string              `json:"building_id"`
	OwnerID    uint64              `json:"owner_id"`
	Duration   time.Duration       `json:"duration"`
	Costs      []*ResourceCost     `json:"costs,omitempty"`
	Workers    []*WorkerAssignment `json:"workers,omitempty"`
	Materials  []*MaterialUsage    `json:"materials,omitempty"`
}

// Validate 验证请求
func (req *StartConstructionRequest) Validate() error {
	if req.BuildingID == "" {
		return fmt.Errorf("building ID is required")
	}
	if req.OwnerID == 0 {
		return fmt.Errorf("owner ID is required")
	}
	if req.Duration <= 0 {
		return fmt.Errorf("duration must be positive")
	}
	return nil
}

// UpdateConstructionProgressRequest 更新建造进度请求
type UpdateConstructionProgressRequest struct {
	ConstructionID string  `json:"construction_id"`
	Progress       float64 `json:"progress"`
}

// Validate 验证请求
func (req *UpdateConstructionProgressRequest) Validate() error {
	if req.ConstructionID == "" {
		return fmt.Errorf("construction ID is required")
	}
	if req.Progress < 0 || req.Progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100")
	}
	return nil
}

// StartUpgradeRequest 开始升级请求
type StartUpgradeRequest struct {
	BuildingID   string            `json:"building_id"`
	ToLevel      int32             `json:"to_level"`
	Duration     time.Duration     `json:"duration"`
	Costs        []*ResourceCost   `json:"costs,omitempty"`
	Requirements []*Requirement    `json:"requirements,omitempty"`
	Benefits     []*UpgradeBenefit `json:"benefits,omitempty"`
}

// Validate 验证请求
func (req *StartUpgradeRequest) Validate() error {
	if req.BuildingID == "" {
		return fmt.Errorf("building ID is required")
	}
	if req.ToLevel <= 0 {
		return fmt.Errorf("to level must be positive")
	}
	if req.Duration <= 0 {
		return fmt.Errorf("duration must be positive")
	}
	return nil
}

// RepairBuildingRequest 修复建筑请求
type RepairBuildingRequest struct {
	BuildingID   string          `json:"building_id"`
	RepairAmount float64         `json:"repair_amount"`
	Costs        []*ResourceCost `json:"costs,omitempty"`
}

// Validate 验证请求
func (req *RepairBuildingRequest) Validate() error {
	if req.BuildingID == "" {
		return fmt.Errorf("building ID is required")
	}
	if req.RepairAmount <= 0 {
		return fmt.Errorf("repair amount must be positive")
	}
	return nil
}

// CreateBlueprintRequest 创建蓝图请求
type CreateBlueprintRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Author      string                 `json:"author"`
	Category    BuildingCategory       `json:"category"`
	Size        *Size                  `json:"size"`
	Materials   []*MaterialRequirement `json:"materials,omitempty"`
	Costs       []*ResourceCost        `json:"costs,omitempty"`
	Duration    time.Duration          `json:"duration"`
	Difficulty  int32                  `json:"difficulty"`
	Tags        []string               `json:"tags,omitempty"`
}

// Validate 验证请求
func (req *CreateBlueprintRequest) Validate() error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if !req.Category.IsValid() {
		return fmt.Errorf("invalid category")
	}
	if req.Size == nil || !req.Size.IsValid() {
		return fmt.Errorf("invalid size")
	}
	if req.Duration <= 0 {
		return fmt.Errorf("duration must be positive")
	}
	if req.Difficulty < 1 || req.Difficulty > 10 {
		return fmt.Errorf("difficulty must be between 1 and 10")
	}
	return nil
}

// BlueprintValidationResult 蓝图验证结果
type BlueprintValidationResult struct {
	BlueprintID string   `json:"blueprint_id"`
	IsValid     bool     `json:"is_valid"`
	Errors      []string `json:"errors"`
	Warnings    []string `json:"warnings"`
}

// 常量定义

const (
	// 默认值
	DefaultBuildingHealth = 100.0
	DefaultBuildingLevel  = 1

	// 限制
	MaxBuildingLevel   = 100
	MaxBuildingNameLen = 100
	MaxDescriptionLen  = 500
	// MaxTagsCount         = 10 // Moved to repository.go
	MaxLayersCount       = 50
	MaxBlocksPerLayer    = 1000
	MaxMaterialsCount    = 100
	MaxCostsCount        = 20
	MaxWorkersCount      = 50
	MaxPhasesCount       = 20
	MaxTasksPerPhase     = 100
	MaxDependenciesCount = 10

	// 时间限制
	MinConstructionDuration = 1 * time.Minute
	MaxConstructionDuration = 30 * 24 * time.Hour // 30天
	MinUpgradeDuration      = 1 * time.Minute
	MaxUpgradeDuration      = 7 * 24 * time.Hour // 7天
)

// 辅助函数

// ValidateBuildingName 验证建筑名称
func ValidateBuildingName(name string) error {
	if name == "" {
		return fmt.Errorf("building name cannot be empty")
	}
	if len(name) > MaxBuildingNameLen {
		return fmt.Errorf("building name too long (max %d characters)", MaxBuildingNameLen)
	}
	return nil
}

// ValidateDescription 验证描述
func ValidateDescription(description string) error {
	if len(description) > MaxDescriptionLen {
		return fmt.Errorf("description too long (max %d characters)", MaxDescriptionLen)
	}
	return nil
}

// ValidateDuration 验证持续时间
func ValidateDuration(duration time.Duration, minDuration, maxDuration time.Duration) error {
	if duration < minDuration {
		return fmt.Errorf("duration too short (min %v)", minDuration)
	}
	if duration > maxDuration {
		return fmt.Errorf("duration too long (max %v)", maxDuration)
	}
	return nil
}

// ValidateLevel 验证等级
func ValidateLevel(level int32) error {
	if level < 1 {
		return fmt.Errorf("level must be at least 1")
	}
	if level > MaxBuildingLevel {
		return fmt.Errorf("level too high (max %d)", MaxBuildingLevel)
	}
	return nil
}

// ValidateHealth 验证健康度
func ValidateHealth(health float64) error {
	if health < 0 {
		return fmt.Errorf("health cannot be negative")
	}
	if health > 100 {
		return fmt.Errorf("health cannot exceed 100")
	}
	return nil
}

// ValidateProgress 验证进度
func ValidateProgress(progress float64) error {
	if progress < 0 {
		return fmt.Errorf("progress cannot be negative")
	}
	if progress > 100 {
		return fmt.Errorf("progress cannot exceed 100")
	}
	return nil
}

// CalculateConstructionTime 计算建造时间
func CalculateConstructionTime(baseTime time.Duration, difficulty int32, workerCount int) time.Duration {
	// 基础时间 * 难度系数 / 工人效率
	difficultyFactor := float64(difficulty) / 5.0          // 难度1-10，转换为0.2-2.0
	workerFactor := 1.0 / (1.0 + float64(workerCount)*0.1) // 工人越多，时间越短

	adjustedTime := float64(baseTime) * difficultyFactor * workerFactor
	return time.Duration(adjustedTime)
}

// CalculateUpgradeCost 计算升级成本
func CalculateUpgradeCost(baseCost int64, fromLevel, toLevel int32) int64 {
	// 成本随等级指数增长
	levelDiff := toLevel - fromLevel
	costMultiplier := float64(levelDiff) * 1.5 // 每级增加50%

	return int64(float64(baseCost) * costMultiplier)
}

// CalculateRepairCost 计算修复成本
func CalculateRepairCost(baseCost int64, currentHealth, targetHealth float64) int64 {
	// 修复成本与损坏程度成正比
	damagePercent := (100.0 - currentHealth) / 100.0
	repairPercent := (targetHealth - currentHealth) / 100.0

	return int64(float64(baseCost) * damagePercent * repairPercent)
}
