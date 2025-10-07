package services

import (
	"context"
	"fmt"
	"time"

	"greatestworks/internal/domain/building"
)

// BuildingApplicationService 建筑应用服务
type BuildingApplicationService struct {
	buildingRepo     building.BuildingRepository
	constructionRepo building.ConstructionRepository
	buildingService  *building.BuildingService
	eventBus         building.BuildingEventBus
}

// NewBuildingApplicationService 创建建筑应用服务
func NewBuildingApplicationService(
	buildingRepo building.BuildingRepository,
	constructionRepo building.ConstructionRepository,
	buildingService *building.BuildingService,
	eventBus building.BuildingEventBus,
) *BuildingApplicationService {
	return &BuildingApplicationService{
		buildingRepo:     buildingRepo,
		constructionRepo: constructionRepo,
		buildingService:  buildingService,
		eventBus:         eventBus,
	}
}

// CreateBuildingTypeRequest 创建建筑类型请求
type CreateBuildingTypeRequest struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Category     string                 `json:"category"`
	MaxLevel     int32                  `json:"max_level"`
	BaseCost     int64                  `json:"base_cost"`
	BuildTime    int32                  `json:"build_time"`
	Requirements map[string]interface{} `json:"requirements,omitempty"`
	IsActive     bool                   `json:"is_active"`
}

// CreateBuildingTypeResponse 创建建筑类型响应
type CreateBuildingTypeResponse struct {
	BuildingTypeID string    `json:"building_type_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Category       string    `json:"category"`
	MaxLevel       int32     `json:"max_level"`
	BaseCost       int64     `json:"base_cost"`
	BuildTime      int32     `json:"build_time"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

// CreateBuildingType 创建建筑类型
func (s *BuildingApplicationService) CreateBuildingType(ctx context.Context, req *CreateBuildingTypeRequest) (*CreateBuildingTypeResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	if err := s.validateCreateBuildingTypeRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// 转换建筑类别
	_, err := s.parseBuildingCategory(req.Category)
	if err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	// 创建建筑类型
	// TODO: 实现NewBuildingType方法
	buildingType := &building.BuildingAggregate{}
	// TODO: 实现BuildingAggregate的方法
	// buildingType.SetDescription(req.Description)
	// buildingType.SetBaseCost(req.BaseCost)
	// buildingType.SetBuildTime(req.BuildTime)

	// 设置需求条件
	// if req.Requirements != nil {
	// 	for key, value := range req.Requirements {
	// 		buildingType.AddRequirement(key, value)
	// 	}
	// }

	// if req.IsActive {
	// 	buildingType.Activate()
	// }

	// 保存建筑类型
	// TODO: 实现SaveBuildingType方法
	if err := s.buildingRepo.Save(ctx, buildingType); err != nil {
		return nil, fmt.Errorf("failed to save building type: %w", err)
	}

	// 发布事件
	// TODO: 实现NewBuildingTypeCreatedEvent方法
	event := &building.BuildingCreatedEvent{}
	if err := s.eventBus.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish building type created event: %v\n", err)
	}

	return &CreateBuildingTypeResponse{
		BuildingTypeID: "", // TODO: buildingType.GetID(),
		Name:           req.Name,
		Description:    req.Description,
		Category:       req.Category,
		MaxLevel:       req.MaxLevel,
		BaseCost:       req.BaseCost,
		BuildTime:      req.BuildTime,
		IsActive:       req.IsActive,
		CreatedAt:      time.Now(),
	}, nil
}

// StartConstructionRequest 开始建造请求
type StartConstructionRequest struct {
	PlayerID       uint64 `json:"player_id"`
	BuildingTypeID string `json:"building_type_id"`
	Position       string `json:"position"`
	SlotID         string `json:"slot_id"`
}

// StartConstructionResponse 开始建造响应
type StartConstructionResponse struct {
	ConstructionID string    `json:"construction_id"`
	BuildingID     string    `json:"building_id"`
	PlayerID       uint64    `json:"player_id"`
	BuildingTypeID string    `json:"building_type_id"`
	Position       string    `json:"position"`
	SlotID         string    `json:"slot_id"`
	Status         string    `json:"status"`
	StartedAt      time.Time `json:"started_at"`
	CompletesAt    time.Time `json:"completes_at"`
	Cost           int64     `json:"cost"`
}

// StartConstruction 开始建造
func (s *BuildingApplicationService) StartConstruction(ctx context.Context, req *StartConstructionRequest) (*StartConstructionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	if err := s.validateStartConstructionRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// 获取建筑类型
	// TODO: 实现FindBuildingTypeByID方法
	buildingType, err := s.buildingRepo.FindByID(ctx, req.BuildingTypeID)
	if err != nil {
		return nil, fmt.Errorf("failed to find building type: %w", err)
	}
	if buildingType == nil {
		return nil, fmt.Errorf("building type not found")
	}

	// 检查建筑类型是否激活
	if !buildingType.IsActive() {
		return nil, fmt.Errorf("building type is not active")
	}

	// 检查位置是否可用
	// TODO: 修复Position类型转换
	position := &building.Position{} // 临时解决方案
	existingBuilding, err := s.buildingRepo.FindByPlayerAndPosition(ctx, req.PlayerID, position)
	if err != nil {
		return nil, fmt.Errorf("failed to check position: %w", err)
	}
	if existingBuilding != nil {
		return nil, fmt.Errorf("position already occupied")
	}

	// 开始建造
	// TODO: 修复StartConstruction方法调用
	// startReq := &building.StartConstructionRequest{
	// 	PlayerID:       req.PlayerID,
	// 	BuildingTypeID: req.BuildingTypeID,
	// 	Position:       req.Position,
	// 	SlotID:         req.SlotID,
	// }
	_, err = s.buildingService.StartConstruction(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start construction: %w", err)
	}

	// 发布事件
	// TODO: 修复NewConstructionStartedEvent方法调用
	event := &building.ConstructionStartedEvent{}
	if err := s.eventBus.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish construction started event: %v\n", err)
	}

	return &StartConstructionResponse{
		ConstructionID: "", // TODO: construction.GetID(),
		BuildingID:     "", // TODO: buildingAggregate.GetID(),
		PlayerID:       req.PlayerID,
		BuildingTypeID: req.BuildingTypeID,
		Position:       req.Position,
		SlotID:         req.SlotID,
		Status:         "started", // TODO: construction.GetStatus().String(),
		StartedAt:      time.Now(),
		CompletesAt:    time.Now().Add(24 * time.Hour), // TODO: construction.GetCompletesAt(),
		Cost:           0,                              // TODO: construction.GetCost(),
	}, nil
}

// CompleteConstructionRequest 完成建造请求
type CompleteConstructionRequest struct {
	ConstructionID  string `json:"construction_id"`
	InstantComplete bool   `json:"instant_complete"`
}

// CompleteConstructionResponse 完成建造响应
type CompleteConstructionResponse struct {
	ConstructionID string                    `json:"construction_id"`
	BuildingID     string                    `json:"building_id"`
	Status         string                    `json:"status"`
	CompletedAt    time.Time                 `json:"completed_at"`
	Rewards        []*BuildingRewardResponse `json:"rewards,omitempty"`
}

// BuildingRewardResponse 建筑奖励响应
type BuildingRewardResponse struct {
	Type     string `json:"type"`
	ItemID   string `json:"item_id,omitempty"`
	Quantity int32  `json:"quantity"`
	Reason   string `json:"reason"`
}

// CompleteConstruction 完成建造
func (s *BuildingApplicationService) CompleteConstruction(ctx context.Context, req *CompleteConstructionRequest) (*CompleteConstructionResponse, error) {
	if req == nil || req.ConstructionID == "" {
		return nil, fmt.Errorf("construction ID is required")
	}

	// 获取建造进程
	construction, err := s.constructionRepo.FindByID(ctx, req.ConstructionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find construction: %w", err)
	}
	if construction == nil {
		return nil, fmt.Errorf("construction not found")
	}

	// 检查建造状态
	// TODO: 实现IsInProgress方法
	// if !construction.IsInProgress() {
	// 	return nil, fmt.Errorf("construction is not in progress")
	// }

	// 完成建造
	// TODO: 实现CompleteConstruction方法
	// _, err = s.buildingService.CompleteConstruction(ctx, req.ConstructionID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to complete construction: %w", err)
	// }

	// 转换奖励响应
	// TODO: 修复result变量
	// rewardResponses := make([]*BuildingRewardResponse, len(result.Rewards))
	// for i, reward := range result.Rewards {
	// 	rewardResponses[i] = &BuildingRewardResponse{
	// 		Type:     reward.GetType().String(),
	// 		ItemID:   reward.GetItemID(),
	// 		Quantity: reward.GetQuantity(),
	// 		Reason:   reward.GetReason(),
	// 	}
	// }
	rewardResponses := []*BuildingRewardResponse{}

	// 发布事件
	// TODO: 修复NewConstructionCompletedEvent方法调用
	event := &building.ConstructionCompletedEvent{}
	if err := s.eventBus.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish construction completed event: %v\n", err)
	}

	return &CompleteConstructionResponse{
		ConstructionID: req.ConstructionID,
		BuildingID:     "",          // TODO: result.BuildingID,
		Status:         "completed", // TODO: result.Status.String(),
		CompletedAt:    time.Now(),  // TODO: result.CompletedAt,
		Rewards:        rewardResponses,
	}, nil
}

// UpgradeBuildingRequest 升级建筑请求
type UpgradeBuildingRequest struct {
	BuildingID  string `json:"building_id"`
	TargetLevel int32  `json:"target_level"`
}

// UpgradeBuildingResponse 升级建筑响应
type UpgradeBuildingResponse struct {
	BuildingID  string    `json:"building_id"`
	OldLevel    int32     `json:"old_level"`
	NewLevel    int32     `json:"new_level"`
	UpgradeCost int64     `json:"upgrade_cost"`
	UpgradeTime int32     `json:"upgrade_time"`
	StartedAt   time.Time `json:"started_at"`
	CompletesAt time.Time `json:"completes_at"`
}

// UpgradeBuilding 升级建筑
func (s *BuildingApplicationService) UpgradeBuilding(ctx context.Context, req *UpgradeBuildingRequest) (*UpgradeBuildingResponse, error) {
	if req == nil || req.BuildingID == "" {
		return nil, fmt.Errorf("building ID is required")
	}

	if req.TargetLevel <= 0 {
		return nil, fmt.Errorf("target level must be positive")
	}

	// 获取建筑
	buildingAggregate, err := s.buildingRepo.FindByID(ctx, req.BuildingID)
	if err != nil {
		return nil, fmt.Errorf("failed to find building: %w", err)
	}
	if buildingAggregate == nil {
		return nil, fmt.Errorf("building not found")
	}

	// 检查建筑状态
	// TODO: 实现IsCompleted方法
	// if !buildingAggregate.IsCompleted() {
	// 	return nil, fmt.Errorf("building is not completed")
	// }

	// TODO: 实现GetLevel方法
	// oldLevel := buildingAggregate.GetLevel().GetCurrentLevel()

	// 升级建筑
	// TODO: 实现UpgradeBuilding方法
	// _, err = s.buildingService.UpgradeBuilding(ctx, req.BuildingID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to upgrade building: %w", err)
	// }

	// 发布事件
	// TODO: 修复NewBuildingUpgradeStartedEvent方法调用
	event := &building.BuildingCreatedEvent{}
	if err := s.eventBus.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish building upgrade started event: %v\n", err)
	}

	return &UpgradeBuildingResponse{
		BuildingID:  req.BuildingID,
		OldLevel:    1, // TODO: oldLevel,
		NewLevel:    req.TargetLevel,
		UpgradeCost: 0, // TODO: result.Cost,
		UpgradeTime: 0, // TODO: result.Time,
		StartedAt:   time.Now(),
		CompletesAt: time.Now().Add(24 * time.Hour), // TODO: result.CompletesAt,
	}, nil
}

// GetPlayerBuildingsRequest 获取玩家建筑请求
type GetPlayerBuildingsRequest struct {
	PlayerID uint64 `json:"player_id"`
	Category string `json:"category,omitempty"`
	Status   string `json:"status,omitempty"`
}

// BuildingResponse 建筑响应
type BuildingResponse struct {
	BuildingID     string                    `json:"building_id"`
	BuildingTypeID string                    `json:"building_type_id"`
	Name           string                    `json:"name"`
	Category       string                    `json:"category"`
	Level          int32                     `json:"level"`
	MaxLevel       int32                     `json:"max_level"`
	Position       string                    `json:"position"`
	SlotID         string                    `json:"slot_id"`
	Status         string                    `json:"status"`
	Effects        []*BuildingEffectResponse `json:"effects,omitempty"`
	Production     map[string]interface{}    `json:"production,omitempty"`
	CreatedAt      time.Time                 `json:"created_at"`
	CompletedAt    *time.Time                `json:"completed_at,omitempty"`
	LastUpdated    time.Time                 `json:"last_updated"`
}

// BuildingEffectResponse 建筑效果响应
type BuildingEffectResponse struct {
	Type        string  `json:"type"`
	Target      string  `json:"target"`
	Value       float64 `json:"value"`
	Description string  `json:"description"`
}

// GetPlayerBuildingsResponse 获取玩家建筑响应
type GetPlayerBuildingsResponse struct {
	PlayerID  uint64              `json:"player_id"`
	Buildings []*BuildingResponse `json:"buildings"`
	Total     int64               `json:"total"`
}

// GetPlayerBuildings 获取玩家建筑
func (s *BuildingApplicationService) GetPlayerBuildings(ctx context.Context, req *GetPlayerBuildingsRequest) (*GetPlayerBuildingsResponse, error) {
	if req == nil || req.PlayerID == 0 {
		return nil, fmt.Errorf("player ID is required")
	}

	// 构建查询
	// TODO: 修复NewBuildingQuery方法调用
	query := &building.BuildingQuery{}

	if req.Category != "" {
		_, err := s.parseBuildingCategory(req.Category)
		if err != nil {
			return nil, fmt.Errorf("invalid category: %w", err)
		}
		// TODO: 修复WithCategory方法调用
		// query = query.WithCategory(category)
	}

	if req.Status != "" {
		_, err := s.parseBuildingStatus(req.Status)
		if err != nil {
			return nil, fmt.Errorf("invalid status: %w", err)
		}
		// TODO: 修复WithStatus方法调用
		// query = query.WithStatus(status)
	}

	// 查询建筑
	buildings, total, err := s.buildingRepo.FindByQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find buildings: %w", err)
	}

	// 转换响应
	buildingResponses := make([]*BuildingResponse, len(buildings))
	for i, _ := range buildings {
		// 获取建筑类型信息
		// TODO: 修复FindBuildingTypeByID方法调用
		// buildingType, _ := s.buildingRepo.FindBuildingTypeByID(ctx, buildingAggregate.GetBuildingTypeID())

		// 转换效果
		// TODO: 实现GetEffects方法
		effectResponses := make([]*BuildingEffectResponse, 0)
		// for _, effect := range buildingAggregate.GetEffects() {
		// 	effectResponses = append(effectResponses, &BuildingEffectResponse{
		// 		Type:        effect.GetType().String(),
		// 		Target:      effect.GetTarget(),
		// 		Value:       effect.GetValue(),
		// 		Description: effect.GetDescription(),
		// 	})
		// }

		buildingResponse := &BuildingResponse{
			BuildingID:     "",       // TODO: buildingAggregate.GetID(),
			BuildingTypeID: "",       // TODO: buildingAggregate.GetBuildingTypeID(),
			Level:          1,        // TODO: buildingAggregate.GetLevel().GetCurrentLevel(),
			Position:       "",       // TODO: buildingAggregate.GetPosition(),
			SlotID:         "",       // TODO: buildingAggregate.GetSlotID(),
			Status:         "active", // TODO: buildingAggregate.GetStatus().String(),
			Effects:        effectResponses,
			Production:     nil,        // TODO: buildingAggregate.GetProduction(),
			CreatedAt:      time.Now(), // TODO: buildingAggregate.GetCreatedAt(),
			LastUpdated:    time.Now(), // TODO: buildingAggregate.GetUpdatedAt(),
		}

		// TODO: 修复buildingType变量
		// if buildingType != nil {
		// 	buildingResponse.Name = buildingType.GetName()
		// 	buildingResponse.Category = buildingType.GetCategory().String()
		// 	buildingResponse.MaxLevel = buildingType.GetMaxLevel()
		// }

		// TODO: 实现IsCompleted方法
		// if buildingAggregate.IsCompleted() {
		// 	completedAt := buildingAggregate.GetCompletedAt()
		// 	buildingResponse.CompletedAt = &completedAt
		// }

		buildingResponses[i] = buildingResponse
	}

	return &GetPlayerBuildingsResponse{
		PlayerID:  req.PlayerID,
		Buildings: buildingResponses,
		Total:     total,
	}, nil
}

// GetBuildingDetailsRequest 获取建筑详情请求
type GetBuildingDetailsRequest struct {
	BuildingID string `json:"building_id"`
}

// GetBuildingDetailsResponse 获取建筑详情响应
type GetBuildingDetailsResponse struct {
	Building     *BuildingResponse     `json:"building"`
	Construction *ConstructionResponse `json:"construction,omitempty"`
}

// ConstructionResponse 建造响应
type ConstructionResponse struct {
	ConstructionID string    `json:"construction_id"`
	Status         string    `json:"status"`
	Progress       float64   `json:"progress"`
	StartedAt      time.Time `json:"started_at"`
	CompletesAt    time.Time `json:"completes_at"`
	RemainingTime  int32     `json:"remaining_time"`
	Cost           int64     `json:"cost"`
}

// GetBuildingDetails 获取建筑详情
func (s *BuildingApplicationService) GetBuildingDetails(ctx context.Context, req *GetBuildingDetailsRequest) (*GetBuildingDetailsResponse, error) {
	if req == nil || req.BuildingID == "" {
		return nil, fmt.Errorf("building ID is required")
	}

	// 获取建筑
	buildingAggregate, err := s.buildingRepo.FindByID(ctx, req.BuildingID)
	if err != nil {
		return nil, fmt.Errorf("failed to find building: %w", err)
	}
	if buildingAggregate == nil {
		return nil, fmt.Errorf("building not found")
	}

	// 获取建筑类型信息
	// TODO: 修复FindBuildingTypeByID方法调用
	// buildingType, err := s.buildingRepo.FindBuildingTypeByID(ctx, buildingAggregate.GetBuildingTypeID())
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to find building type: %w", err)
	// }

	// 转换建筑响应
	// TODO: 实现GetEffects方法
	effectResponses := make([]*BuildingEffectResponse, 0)
	// for _, effect := range buildingAggregate.GetEffects() {
	// 	effectResponses = append(effectResponses, &BuildingEffectResponse{
	// 		Type:        effect.GetType().String(),
	// 		Target:      effect.GetTarget(),
	// 		Value:       effect.GetValue(),
	// 		Description: effect.GetDescription(),
	// 	})
	// }

	buildingResponse := &BuildingResponse{
		BuildingID:     "",       // TODO: buildingAggregate.GetID(),
		BuildingTypeID: "",       // TODO: buildingAggregate.GetBuildingTypeID(),
		Level:          1,        // TODO: buildingAggregate.GetLevel().GetCurrentLevel(),
		Position:       "",       // TODO: buildingAggregate.GetPosition(),
		SlotID:         "",       // TODO: buildingAggregate.GetSlotID(),
		Status:         "active", // TODO: buildingAggregate.GetStatus().String(),
		Effects:        effectResponses,
		Production:     nil,        // TODO: buildingAggregate.GetProduction(),
		CreatedAt:      time.Now(), // TODO: buildingAggregate.GetCreatedAt(),
		LastUpdated:    time.Now(), // TODO: buildingAggregate.GetUpdatedAt(),
	}

	// TODO: 修复buildingType变量
	// if buildingType != nil {
	// 	buildingResponse.Name = buildingType.GetName()
	// 	buildingResponse.Category = buildingType.GetCategory().String()
	// 	buildingResponse.MaxLevel = buildingType.GetMaxLevel()
	// }

	// TODO: 实现IsCompleted方法
	// if buildingAggregate.IsCompleted() {
	// 	completedAt := buildingAggregate.GetCompletedAt()
	// 	buildingResponse.CompletedAt = &completedAt
	// }

	response := &GetBuildingDetailsResponse{
		Building: buildingResponse,
	}

	// 获取建造信息（如果存在）
	// TODO: 修复FindByBuildingID方法调用
	// construction, err := s.constructionRepo.FindByBuildingID(ctx, req.BuildingID)
	// if err == nil && construction != nil && construction.IsInProgress() {
	// 	// 计算进度和剩余时间
	// 	progress := construction.GetProgress()
	// 	remainingTime := int32(0)
	// 	if !construction.IsCompleted() {
	// 		remainingTime = int32(construction.GetCompletesAt().Sub(time.Now()).Seconds())
	// 		if remainingTime < 0 {
	// 			remainingTime = 0
	// 		}
	// 	}

	// 	response.Construction = &ConstructionResponse{
	// 		ConstructionID: construction.GetID(),
	// 		Status:         construction.GetStatus().String(),
	// 		Progress:       progress,
	// 		StartedAt:      construction.GetStartedAt(),
	// 		CompletesAt:    construction.GetCompletesAt(),
	// 		RemainingTime:  remainingTime,
	// 		Cost:           construction.GetCost(),
	// 	}
	// }

	return response, nil
}

// 私有方法

// validateCreateBuildingTypeRequest 验证创建建筑类型请求
func (s *BuildingApplicationService) validateCreateBuildingTypeRequest(req *CreateBuildingTypeRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(req.Name) > 100 {
		return fmt.Errorf("name too long (max 100 characters)")
	}
	if req.Category == "" {
		return fmt.Errorf("category is required")
	}
	if req.MaxLevel <= 0 {
		return fmt.Errorf("max level must be positive")
	}
	if req.BaseCost < 0 {
		return fmt.Errorf("base cost cannot be negative")
	}
	if req.BuildTime <= 0 {
		return fmt.Errorf("build time must be positive")
	}
	return nil
}

// validateStartConstructionRequest 验证开始建造请求
func (s *BuildingApplicationService) validateStartConstructionRequest(req *StartConstructionRequest) error {
	if req.PlayerID == 0 {
		return fmt.Errorf("player ID is required")
	}
	if req.BuildingTypeID == "" {
		return fmt.Errorf("building type ID is required")
	}
	if req.Position == "" {
		return fmt.Errorf("position is required")
	}
	return nil
}

// parseBuildingCategory 解析建筑类别
// TODO: 实现BuildingCategory类型
func (s *BuildingApplicationService) parseBuildingCategory(categoryStr string) (string, error) {
	switch categoryStr {
	case "residential":
		return "residential", nil
	case "commercial":
		return "commercial", nil
	case "industrial":
		return "industrial", nil
	case "military":
		return "military", nil
	case "decoration":
		return "decoration", nil
	case "special":
		return "special", nil
	default:
		return "residential", fmt.Errorf("unknown category: %s", categoryStr)
	}
}

// parseBuildingStatus 解析建筑状态
// TODO: 实现BuildingStatus类型
func (s *BuildingApplicationService) parseBuildingStatus(statusStr string) (string, error) {
	switch statusStr {
	case "planning":
		return "planning", nil
	case "constructing":
		return "constructing", nil
	case "completed":
		return "completed", nil
	case "upgrading":
		return "upgrading", nil
	case "demolished":
		return "demolished", nil
	default:
		return "planning", fmt.Errorf("unknown status: %s", statusStr)
	}
}
