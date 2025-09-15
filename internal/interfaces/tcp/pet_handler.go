package tcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/greatestworks/application/services"
	"github.com/greatestworks/internal/interfaces/tcp/protocol"
)

// PetHandler 宠物TCP处理器
type PetHandler struct {
	petService *services.PetApplicationService
}

// NewPetHandler 创建宠物处理器
func NewPetHandler(petService *services.PetApplicationService) *PetHandler {
	return &PetHandler{
		petService: petService,
	}
}

// HandleCreatePet 处理创建宠物请求
func (h *PetHandler) HandleCreatePet(ctx context.Context, req *protocol.CreatePetRequest) (*protocol.CreatePetResponse, error) {
	log.Printf("[PetHandler] HandleCreatePet: PlayerID=%d, SpeciesID=%s", req.PlayerID, req.SpeciesID)
	
	// 转换为应用服务请求
	serviceReq := &services.CreatePetRequest{
		PlayerID:  req.PlayerID,
		SpeciesID: req.SpeciesID,
		Name:      req.Name,
		Rarity:    req.Rarity,
		Quality:   req.Quality,
	}
	
	// 调用应用服务
	serviceResp, err := h.petService.CreatePet(ctx, serviceReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create pet: %w", err)
	}
	
	// 转换响应
	return &protocol.CreatePetResponse{
		PetID:     serviceResp.PetID,
		PlayerID:  serviceResp.PlayerID,
		SpeciesID: serviceResp.SpeciesID,
		Name:      serviceResp.Name,
		Level:     serviceResp.Level,
		Exp:       serviceResp.Exp,
		MaxExp:    serviceResp.MaxExp,
		Rarity:    serviceResp.Rarity,
		Quality:   serviceResp.Quality,
		CreatedAt: serviceResp.CreatedAt.Unix(),
	}, nil
}

// HandleFeedPet 处理喂养宠物请求
func (h *PetHandler) HandleFeedPet(ctx context.Context, req *protocol.FeedPetRequest) (*protocol.FeedPetResponse, error) {
	log.Printf("[PetHandler] HandleFeedPet: PetID=%s, FoodID=%s, Quantity=%d", req.PetID, req.FoodID, req.Quantity)
	
	// 转换为应用服务请求
	serviceReq := &services.FeedPetRequest{
		PetID:    req.PetID,
		FoodID:   req.FoodID,
		Quantity: req.Quantity,
	}
	
	// 调用应用服务
	serviceResp, err := h.petService.FeedPet(ctx, serviceReq)
	if err != nil {
		return nil, fmt.Errorf("failed to feed pet: %w", err)
	}
	
	// 转换响应
	return &protocol.FeedPetResponse{
		PetID:        req.PetID,
		OldHunger:    serviceResp.OldHunger,
		NewHunger:    serviceResp.NewHunger,
		OldHappiness: serviceResp.OldHappiness,
		NewHappiness: serviceResp.NewHappiness,
		ExpGained:    serviceResp.ExpGained,
		LevelUp:      serviceResp.LevelUp,
		FedAt:        serviceResp.FedAt.Unix(),
	}, nil
}

// HandleTrainPet 处理训练宠物请求
func (h *PetHandler) HandleTrainPet(ctx context.Context, req *protocol.TrainPetRequest) (*protocol.TrainPetResponse, error) {
	log.Printf("[PetHandler] HandleTrainPet: PetID=%s, TrainingType=%s, Duration=%d", req.PetID, req.TrainingType, req.Duration)
	
	// 转换为应用服务请求
	serviceReq := &services.TrainPetRequest{
		PetID:        req.PetID,
		TrainingType: req.TrainingType,
		Duration:     req.Duration,
		Intensity:    req.Intensity,
	}
	
	// 调用应用服务
	serviceResp, err := h.petService.TrainPet(ctx, serviceReq)
	if err != nil {
		return nil, fmt.Errorf("failed to train pet: %w", err)
	}
	
	// 转换属性变化
	attributeChanges := make(map[string]int64)
	for attr, change := range serviceResp.AttributeChanges {
		attributeChanges[attr] = change
	}
	
	// 转换响应
	return &protocol.TrainPetResponse{
		PetID:            req.PetID,
		TrainingType:     req.TrainingType,
		ExpGained:        serviceResp.ExpGained,
		AttributeChanges: attributeChanges,
		EnergyConsumed:   serviceResp.EnergyConsumed,
		LevelUp:          serviceResp.LevelUp,
		SkillsLearned:    serviceResp.SkillsLearned,
		TrainedAt:        serviceResp.TrainedAt.Unix(),
	}, nil
}

// HandleGetPet 处理获取宠物请求
func (h *PetHandler) HandleGetPet(ctx context.Context, req *protocol.GetPetRequest) (*protocol.GetPetResponse, error) {
	log.Printf("[PetHandler] HandleGetPet: PetID=%s", req.PetID)
	
	// 转换为应用服务请求
	serviceReq := &services.GetPetRequest{
		PetID: req.PetID,
	}
	
	// 调用应用服务
	serviceResp, err := h.petService.GetPet(ctx, serviceReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get pet: %w", err)
	}
	
	if serviceResp.Pet == nil {
		return &protocol.GetPetResponse{
			Found: false,
		}, nil
	}
	
	// 转换宠物信息
	pet := &protocol.PetInfo{
		PetID:        serviceResp.Pet.PetID,
		PlayerID:     serviceResp.Pet.PlayerID,
		SpeciesID:    serviceResp.Pet.SpeciesID,
		Name:         serviceResp.Pet.Name,
		Level:        serviceResp.Pet.Level,
		Exp:          serviceResp.Pet.Exp,
		MaxExp:       serviceResp.Pet.MaxExp,
		Rarity:       serviceResp.Pet.Rarity,
		Quality:      serviceResp.Pet.Quality,
		Attributes:   serviceResp.Pet.Attributes,
		Skills:       serviceResp.Pet.Skills,
		EquippedSkin: serviceResp.Pet.EquippedSkin,
		Mood:         serviceResp.Pet.Mood,
		Hunger:       serviceResp.Pet.Hunger,
		Energy:       serviceResp.Pet.Energy,
		Health:       serviceResp.Pet.Health,
		Happiness:    serviceResp.Pet.Happiness,
		IsActive:     serviceResp.Pet.IsActive,
		LastFedAt:    serviceResp.Pet.LastFedAt.Unix(),
		LastPlayedAt: serviceResp.Pet.LastPlayedAt.Unix(),
		CreatedAt:    serviceResp.Pet.CreatedAt.Unix(),
		UpdatedAt:    serviceResp.Pet.UpdatedAt.Unix(),
	}
	
	// 转换响应
	return &protocol.GetPetResponse{
		Found: true,
		Pet:   pet,
	}, nil
}

// HandleGetPlayerPets 处理获取玩家宠物列表请求
func (h *PetHandler) HandleGetPlayerPets(ctx context.Context, req *protocol.GetPlayerPetsRequest) (*protocol.GetPlayerPetsResponse, error) {
	log.Printf("[PetHandler] HandleGetPlayerPets: PlayerID=%d", req.PlayerID)
	
	// 转换为应用服务请求
	serviceReq := &services.GetPlayerPetsRequest{
		PlayerID: req.PlayerID,
		Rarity:   req.Rarity,
		Quality:  req.Quality,
		MinLevel: req.MinLevel,
		MaxLevel: req.MaxLevel,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	
	// 调用应用服务
	serviceResp, err := h.petService.GetPlayerPets(ctx, serviceReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get player pets: %w", err)
	}
	
	// 转换宠物列表
	pets := make([]*protocol.PetInfo, len(serviceResp.Pets))
	for i, petResp := range serviceResp.Pets {
		pets[i] = &protocol.PetInfo{
			PetID:        petResp.PetID,
			PlayerID:     petResp.PlayerID,
			SpeciesID:    petResp.SpeciesID,
			Name:         petResp.Name,
			Level:        petResp.Level,
			Exp:          petResp.Exp,
			MaxExp:       petResp.MaxExp,
			Rarity:       petResp.Rarity,
			Quality:      petResp.Quality,
			Attributes:   petResp.Attributes,
			Skills:       petResp.Skills,
			EquippedSkin: petResp.EquippedSkin,
			Mood:         petResp.Mood,
			Hunger:       petResp.Hunger,
			Energy:       petResp.Energy,
			Health:       petResp.Health,
			Happiness:    petResp.Happiness,
			IsActive:     petResp.IsActive,
			LastFedAt:    petResp.LastFedAt.Unix(),
			LastPlayedAt: petResp.LastPlayedAt.Unix(),
			CreatedAt:    petResp.CreatedAt.Unix(),
			UpdatedAt:    petResp.UpdatedAt.Unix(),
		}
	}
	
	// 转换响应
	return &protocol.GetPlayerPetsResponse{
		PlayerID:   req.PlayerID,
		Pets:       pets,
		Total:      serviceResp.Total,
		Page:       serviceResp.Page,
		PageSize:   serviceResp.PageSize,
		TotalPages: serviceResp.TotalPages,
	}, nil
}

// HandleEvolvePet 处理宠物进化请求
func (h *PetHandler) HandleEvolvePet(ctx context.Context, req *protocol.EvolvePetRequest) (*protocol.EvolvePetResponse, error) {
	log.Printf("[PetHandler] HandleEvolvePet: PetID=%s, TargetSpecies=%s", req.PetID, req.TargetSpecies)
	
	// 转换为应用服务请求
	serviceReq := &services.EvolvePetRequest{
		PetID:         req.PetID,
		TargetSpecies: req.TargetSpecies,
		Materials:     req.Materials,
	}
	
	// 调用应用服务
	serviceResp, err := h.petService.EvolvePet(ctx, serviceReq)
	if err != nil {
		return nil, fmt.Errorf("failed to evolve pet: %w", err)
	}
	
	// 转换响应
	return &protocol.EvolvePetResponse{
		PetID:           req.PetID,
		OldSpecies:      serviceResp.OldSpecies,
		NewSpecies:      serviceResp.NewSpecies,
		OldRarity:       serviceResp.OldRarity,
		NewRarity:       serviceResp.NewRarity,
		AttributeBonus:  serviceResp.AttributeBonus,
		NewSkills:       serviceResp.NewSkills,
		MaterialsUsed:   serviceResp.MaterialsUsed,
		EvolvedAt:       serviceResp.EvolvedAt.Unix(),
	}, nil
}

// HandleEquipPetSkin 处理装备宠物皮肤请求
func (h *PetHandler) HandleEquipPetSkin(ctx context.Context, req *protocol.EquipPetSkinRequest) (*protocol.EquipPetSkinResponse, error) {
	log.Printf("[PetHandler] HandleEquipPetSkin: PetID=%s, SkinID=%s", req.PetID, req.SkinID)
	
	// 转换为应用服务请求
	serviceReq := &services.EquipPetSkinRequest{
		PetID:  req.PetID,
		SkinID: req.SkinID,
	}
	
	// 调用应用服务
	serviceResp, err := h.petService.EquipPetSkin(ctx, serviceReq)
	if err != nil {
		return nil, fmt.Errorf("failed to equip pet skin: %w", err)
	}
	
	// 转换响应
	return &protocol.EquipPetSkinResponse{
		PetID:         req.PetID,
		OldSkinID:     serviceResp.OldSkinID,
		NewSkinID:     serviceResp.NewSkinID,
		EffectChanges: serviceResp.EffectChanges,
		EquippedAt:    serviceResp.EquippedAt.Unix(),
	}, nil
}

// HandleSynthesizePet 处理宠物合成请求
func (h *PetHandler) HandleSynthesizePet(ctx context.Context, req *protocol.SynthesizePetRequest) (*protocol.SynthesizePetResponse, error) {
	log.Printf("[PetHandler] HandleSynthesizePet: PlayerID=%d, FragmentID=%s, Quantity=%d", req.PlayerID, req.FragmentID, req.Quantity)
	
	// 转换为应用服务请求
	serviceReq := &services.SynthesizePetRequest{
		PlayerID:   req.PlayerID,
		FragmentID: req.FragmentID,
		Quantity:   req.Quantity,
	}
	
	// 调用应用服务
	serviceResp, err := h.petService.SynthesizePet(ctx, serviceReq)
	if err != nil {
		return nil, fmt.Errorf("failed to synthesize pet: %w", err)
	}
	
	// 转换响应
	return &protocol.SynthesizePetResponse{
		PlayerID:        req.PlayerID,
		FragmentID:      req.FragmentID,
		QuantityUsed:    serviceResp.QuantityUsed,
		PetID:           serviceResp.PetID,
		SpeciesID:       serviceResp.SpeciesID,
		Rarity:          serviceResp.Rarity,
		Success:         serviceResp.Success,
		SynthesizedAt:   serviceResp.SynthesizedAt.Unix(),
	}, nil
}

// RegisterHandlers 注册处理器到路由器
func (h *PetHandler) RegisterHandlers(router *TCPRouter) {
	// 注册宠物相关的消息处理器
	router.RegisterHandler(protocol.MsgTypeCreatePet, h.handleCreatePetMessage)
	router.RegisterHandler(protocol.MsgTypeFeedPet, h.handleFeedPetMessage)
	router.RegisterHandler(protocol.MsgTypeTrainPet, h.handleTrainPetMessage)
	router.RegisterHandler(protocol.MsgTypeGetPet, h.handleGetPetMessage)
	router.RegisterHandler(protocol.MsgTypeGetPlayerPets, h.handleGetPlayerPetsMessage)
	router.RegisterHandler(protocol.MsgTypeEvolvePet, h.handleEvolvePetMessage)
	router.RegisterHandler(protocol.MsgTypeEquipPetSkin, h.handleEquipPetSkinMessage)
	router.RegisterHandler(protocol.MsgTypeSynthesizePet, h.handleSynthesizePetMessage)
}

// 消息处理器包装函数

func (h *PetHandler) handleCreatePetMessage(ctx context.Context, conn *TCPConnection, msg *protocol.Message) error {
	var req protocol.CreatePetRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return fmt.Errorf("failed to unmarshal create pet request: %w", err)
	}
	
	resp, err := h.HandleCreatePet(ctx, &req)
	if err != nil {
		return err
	}
	
	return conn.SendResponse(msg.ID, protocol.MsgTypeCreatePetResponse, resp)
}

func (h *PetHandler) handleFeedPetMessage(ctx context.Context, conn *TCPConnection, msg *protocol.Message) error {
	var req protocol.FeedPetRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return fmt.Errorf("failed to unmarshal feed pet request: %w", err)
	}
	
	resp, err := h.HandleFeedPet(ctx, &req)
	if err != nil {
		return err
	}
	
	return conn.SendResponse(msg.ID, protocol.MsgTypeFeedPetResponse, resp)
}

func (h *PetHandler) handleTrainPetMessage(ctx context.Context, conn *TCPConnection, msg *protocol.Message) error {
	var req protocol.TrainPetRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return fmt.Errorf("failed to unmarshal train pet request: %w", err)
	}
	
	resp, err := h.HandleTrainPet(ctx, &req)
	if err != nil {
		return err
	}
	
	return conn.SendResponse(msg.ID, protocol.MsgTypeTrainPetResponse, resp)
}

func (h *PetHandler) handleGetPetMessage(ctx context.Context, conn *TCPConnection, msg *protocol.Message) error {
	var req protocol.GetPetRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return fmt.Errorf("failed to unmarshal get pet request: %w", err)
	}
	
	resp, err := h.HandleGetPet(ctx, &req)
	if err != nil {
		return err
	}
	
	return conn.SendResponse(msg.ID, protocol.MsgTypeGetPetResponse, resp)
}

func (h *PetHandler) handleGetPlayerPetsMessage(ctx context.Context, conn *TCPConnection, msg *protocol.Message) error {
	var req protocol.GetPlayerPetsRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return fmt.Errorf("failed to unmarshal get player pets request: %w", err)
	}
	
	resp, err := h.HandleGetPlayerPets(ctx, &req)
	if err != nil {
		return err
	}
	
	return conn.SendResponse(msg.ID, protocol.MsgTypeGetPlayerPetsResponse, resp)
}

func (h *PetHandler) handleEvolvePetMessage(ctx context.Context, conn *TCPConnection, msg *protocol.Message) error {
	var req protocol.EvolvePetRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return fmt.Errorf("failed to unmarshal evolve pet request: %w", err)
	}
	
	resp, err := h.HandleEvolvePet(ctx, &req)
	if err != nil {
		return err
	}
	
	return conn.SendResponse(msg.ID, protocol.MsgTypeEvolvePetResponse, resp)
}

func (h *PetHandler) handleEquipPetSkinMessage(ctx context.Context, conn *TCPConnection, msg *protocol.Message) error {
	var req protocol.EquipPetSkinRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return fmt.Errorf("failed to unmarshal equip pet skin request: %w", err)
	}
	
	resp, err := h.HandleEquipPetSkin(ctx, &req)
	if err != nil {
		return err
	}
	
	return conn.SendResponse(msg.ID, protocol.MsgTypeEquipPetSkinResponse, resp)
}

func (h *PetHandler) handleSynthesizePetMessage(ctx context.Context, conn *TCPConnection, msg *protocol.Message) error {
	var req protocol.SynthesizePetRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return fmt.Errorf("failed to unmarshal synthesize pet request: %w", err)
	}
	
	resp, err := h.HandleSynthesizePet(ctx, &req)
	if err != nil {
		return err
	}
	
	return conn.SendResponse(msg.ID, protocol.MsgTypeSynthesizePetResponse, resp)
}