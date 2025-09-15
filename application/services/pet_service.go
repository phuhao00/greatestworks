package services

import (
	"context"
	"fmt"
	"time"

	"github.com/greatestworks/internal/domain/pet"
)

// PetApplicationService 宠物应用服务
type PetApplicationService struct {
	petRepo         pet.PetRepository
	fragmentRepo    pet.PetFragmentRepository
	skinRepo        pet.PetSkinRepository
	bondRepo        pet.PetBondRepository
	pictorialRepo   pet.PetPictorialRepository
	petService      *pet.PetService
	eventBus        pet.PetEventBus
}

// NewPetApplicationService 创建宠物应用服务
func NewPetApplicationService(
	petRepo pet.PetRepository,
	fragmentRepo pet.PetFragmentRepository,
	skinRepo pet.PetSkinRepository,
	bondRepo pet.PetBondRepository,
	pictorialRepo pet.PetPictorialRepository,
	petService *pet.PetService,
	eventBus pet.PetEventBus,
) *PetApplicationService {
	return &PetApplicationService{
		petRepo:       petRepo,
		fragmentRepo:  fragmentRepo,
		skinRepo:      skinRepo,
		bondRepo:      bondRepo,
		pictorialRepo: pictorialRepo,
		petService:    petService,
		eventBus:      eventBus,
	}
}

// CreatePetRequest 创建宠物请求
type CreatePetRequest struct {
	OwnerID   uint64 `json:"owner_id"`
	PetType   string `json:"pet_type"`
	Name      string `json:"name"`
	Rarity    string `json:"rarity"`
	Source    string `json:"source"`
}

// CreatePetResponse 创建宠物响应
type CreatePetResponse struct {
	PetID     string    `json:"pet_id"`
	Name      string    `json:"name"`
	PetType   string    `json:"pet_type"`
	Rarity    string    `json:"rarity"`
	Level     int32     `json:"level"`
	Exp       int64     `json:"exp"`
	CreatedAt time.Time `json:"created_at"`
}

// CreatePet 创建宠物
func (s *PetApplicationService) CreatePet(ctx context.Context, req *CreatePetRequest) (*CreatePetResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	
	if err := s.validateCreatePetRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	
	// 转换稀有度
	rarity, err := s.parseRarity(req.Rarity)
	if err != nil {
		return nil, fmt.Errorf("invalid rarity: %w", err)
	}
	
	// 转换来源
	source, err := s.parseSource(req.Source)
	if err != nil {
		return nil, fmt.Errorf("invalid source: %w", err)
	}
	
	// 创建宠物聚合根
	petAggregate := pet.NewPetAggregate(req.OwnerID, req.PetType, req.Name)
	petAggregate.SetRarity(rarity)
	petAggregate.SetSource(source)
	
	// 保存宠物
	if err := s.petRepo.Save(ctx, petAggregate); err != nil {
		return nil, fmt.Errorf("failed to save pet: %w", err)
	}
	
	// 发布事件
	event := pet.NewPetCreatedEvent(petAggregate.GetID(), req.OwnerID, req.PetType, req.Name)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		// 记录错误但不影响主流程
		fmt.Printf("failed to publish pet created event: %v\n", err)
	}
	
	return &CreatePetResponse{
		PetID:     petAggregate.GetID(),
		Name:      petAggregate.GetName(),
		PetType:   petAggregate.GetPetType(),
		Rarity:    petAggregate.GetRarity().String(),
		Level:     petAggregate.GetLevel(),
		Exp:       petAggregate.GetExp(),
		CreatedAt: petAggregate.GetCreatedAt(),
	}, nil
}

// FeedPetRequest 喂养宠物请求
type FeedPetRequest struct {
	PetID    string `json:"pet_id"`
	FoodType string `json:"food_type"`
	Amount   int32  `json:"amount"`
}

// FeedPetResponse 喂养宠物响应
type FeedPetResponse struct {
	PetID       string `json:"pet_id"`
	ExpGained   int64  `json:"exp_gained"`
	LeveledUp   bool   `json:"leveled_up"`
	NewLevel    int32  `json:"new_level"`
	NewExp      int64  `json:"new_exp"`
	Happiness   int32  `json:"happiness"`
}

// FeedPet 喂养宠物
func (s *PetApplicationService) FeedPet(ctx context.Context, req *FeedPetRequest) (*FeedPetResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	
	if err := s.validateFeedPetRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	
	// 获取宠物
	petAggregate, err := s.petRepo.FindByID(ctx, req.PetID)
	if err != nil {
		return nil, fmt.Errorf("failed to find pet: %w", err)
	}
	if petAggregate == nil {
		return nil, fmt.Errorf("pet not found")
	}
	
	// 计算经验值
	expGained := s.calculateFoodExp(req.FoodType, req.Amount)
	
	// 喂养宠物
	leveledUp := petAggregate.AddExp(expGained)
	petAggregate.Feed(req.FoodType, req.Amount)
	
	// 保存宠物
	if err := s.petRepo.Save(ctx, petAggregate); err != nil {
		return nil, fmt.Errorf("failed to save pet: %w", err)
	}
	
	// 发布事件
	event := pet.NewPetFedEvent(petAggregate.GetID(), req.FoodType, req.Amount, expGained)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish pet fed event: %v\n", err)
	}
	
	if leveledUp {
		levelUpEvent := pet.NewPetLevelUpEvent(petAggregate.GetID(), petAggregate.GetLevel()-1, petAggregate.GetLevel())
		if err := s.eventBus.Publish(ctx, levelUpEvent); err != nil {
			fmt.Printf("failed to publish pet level up event: %v\n", err)
		}
	}
	
	return &FeedPetResponse{
		PetID:     petAggregate.GetID(),
		ExpGained: expGained,
		LeveledUp: leveledUp,
		NewLevel:  petAggregate.GetLevel(),
		NewExp:    petAggregate.GetExp(),
		Happiness: petAggregate.GetHappiness(),
	}, nil
}

// TrainPetRequest 训练宠物请求
type TrainPetRequest struct {
	PetID        string `json:"pet_id"`
	TrainingType string `json:"training_type"`
	Duration     int32  `json:"duration"` // 训练时长（分钟）
}

// TrainPetResponse 训练宠物响应
type TrainPetResponse struct {
	PetID           string            `json:"pet_id"`
	TrainingType    string            `json:"training_type"`
	AttributeGains  map[string]int32  `json:"attribute_gains"`
	ExpGained       int64             `json:"exp_gained"`
	LeveledUp       bool              `json:"leveled_up"`
	NewLevel        int32             `json:"new_level"`
	SkillsLearned   []string          `json:"skills_learned"`
}

// TrainPet 训练宠物
func (s *PetApplicationService) TrainPet(ctx context.Context, req *TrainPetRequest) (*TrainPetResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	
	if err := s.validateTrainPetRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	
	// 获取宠物
	petAggregate, err := s.petRepo.FindByID(ctx, req.PetID)
	if err != nil {
		return nil, fmt.Errorf("failed to find pet: %w", err)
	}
	if petAggregate == nil {
		return nil, fmt.Errorf("pet not found")
	}
	
	// 计算训练收益
	attributeGains := s.calculateTrainingGains(req.TrainingType, req.Duration, petAggregate.GetLevel())
	expGained := s.calculateTrainingExp(req.TrainingType, req.Duration)
	
	// 训练宠物
	leveledUp := petAggregate.AddExp(expGained)
	for attr, gain := range attributeGains {
		petAggregate.AddAttribute(attr, gain)
	}
	
	// 检查是否学会新技能
	skillsLearned := s.checkSkillLearning(petAggregate, req.TrainingType)
	for _, skill := range skillsLearned {
		petAggregate.LearnSkill(skill)
	}
	
	// 保存宠物
	if err := s.petRepo.Save(ctx, petAggregate); err != nil {
		return nil, fmt.Errorf("failed to save pet: %w", err)
	}
	
	// 发布事件
	event := pet.NewPetTrainedEvent(petAggregate.GetID(), req.TrainingType, attributeGains, expGained)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish pet trained event: %v\n", err)
	}
	
	return &TrainPetResponse{
		PetID:          petAggregate.GetID(),
		TrainingType:   req.TrainingType,
		AttributeGains: attributeGains,
		ExpGained:      expGained,
		LeveledUp:      leveledUp,
		NewLevel:       petAggregate.GetLevel(),
		SkillsLearned:  skillsLearned,
	}, nil
}

// GetPetRequest 获取宠物请求
type GetPetRequest struct {
	PetID string `json:"pet_id"`
}

// GetPetResponse 获取宠物响应
type GetPetResponse struct {
	PetID      string            `json:"pet_id"`
	OwnerID    uint64            `json:"owner_id"`
	Name       string            `json:"name"`
	PetType    string            `json:"pet_type"`
	Rarity     string            `json:"rarity"`
	Level      int32             `json:"level"`
	Exp        int64             `json:"exp"`
	MaxExp     int64             `json:"max_exp"`
	Happiness  int32             `json:"happiness"`
	Health     int32             `json:"health"`
	Attributes map[string]int32  `json:"attributes"`
	Skills     []string          `json:"skills"`
	Skins      []string          `json:"skins"`
	CurrentSkin string           `json:"current_skin"`
	Bonds      []string          `json:"bonds"`
	Status     string            `json:"status"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

// GetPet 获取宠物信息
func (s *PetApplicationService) GetPet(ctx context.Context, req *GetPetRequest) (*GetPetResponse, error) {
	if req == nil || req.PetID == "" {
		return nil, fmt.Errorf("pet ID is required")
	}
	
	// 获取宠物
	petAggregate, err := s.petRepo.FindByID(ctx, req.PetID)
	if err != nil {
		return nil, fmt.Errorf("failed to find pet: %w", err)
	}
	if petAggregate == nil {
		return nil, fmt.Errorf("pet not found")
	}
	
	return &GetPetResponse{
		PetID:       petAggregate.GetID(),
		OwnerID:     petAggregate.GetOwnerID(),
		Name:        petAggregate.GetName(),
		PetType:     petAggregate.GetPetType(),
		Rarity:      petAggregate.GetRarity().String(),
		Level:       petAggregate.GetLevel(),
		Exp:         petAggregate.GetExp(),
		MaxExp:      petAggregate.GetMaxExp(),
		Happiness:   petAggregate.GetHappiness(),
		Health:      petAggregate.GetHealth(),
		Attributes:  petAggregate.GetAttributes(),
		Skills:      petAggregate.GetSkills(),
		Skins:       petAggregate.GetUnlockedSkins(),
		CurrentSkin: petAggregate.GetCurrentSkin(),
		Bonds:       petAggregate.GetBonds(),
		Status:      petAggregate.GetStatus().String(),
		CreatedAt:   petAggregate.GetCreatedAt(),
		UpdatedAt:   petAggregate.GetUpdatedAt(),
	}, nil
}

// GetPlayerPetsRequest 获取玩家宠物列表请求
type GetPlayerPetsRequest struct {
	OwnerID  uint64 `json:"owner_id"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	SortBy   string `json:"sort_by"`   // level, exp, happiness, created_at
	SortOrder string `json:"sort_order"` // asc, desc
}

// GetPlayerPetsResponse 获取玩家宠物列表响应
type GetPlayerPetsResponse struct {
	Pets       []*GetPetResponse `json:"pets"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int64             `json:"total_pages"`
}

// GetPlayerPets 获取玩家宠物列表
func (s *PetApplicationService) GetPlayerPets(ctx context.Context, req *GetPlayerPetsRequest) (*GetPlayerPetsResponse, error) {
	if req == nil || req.OwnerID == 0 {
		return nil, fmt.Errorf("owner ID is required")
	}
	
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}
	
	// 构建查询
	query := pet.NewPetQuery().
		WithOwner(req.OwnerID).
		WithSort(req.SortBy, req.SortOrder).
		WithPagination(req.Page, req.PageSize)
	
	// 查询宠物
	pets, total, err := s.petRepo.FindByQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find pets: %w", err)
	}
	
	// 转换响应
	petResponses := make([]*GetPetResponse, len(pets))
	for i, petAggregate := range pets {
		petResponses[i] = &GetPetResponse{
			PetID:       petAggregate.GetID(),
			OwnerID:     petAggregate.GetOwnerID(),
			Name:        petAggregate.GetName(),
			PetType:     petAggregate.GetPetType(),
			Rarity:      petAggregate.GetRarity().String(),
			Level:       petAggregate.GetLevel(),
			Exp:         petAggregate.GetExp(),
			MaxExp:      petAggregate.GetMaxExp(),
			Happiness:   petAggregate.GetHappiness(),
			Health:      petAggregate.GetHealth(),
			Attributes:  petAggregate.GetAttributes(),
			Skills:      petAggregate.GetSkills(),
			Skins:       petAggregate.GetUnlockedSkins(),
			CurrentSkin: petAggregate.GetCurrentSkin(),
			Bonds:       petAggregate.GetBonds(),
			Status:      petAggregate.GetStatus().String(),
			CreatedAt:   petAggregate.GetCreatedAt(),
			UpdatedAt:   petAggregate.GetUpdatedAt(),
		}
	}
	
	totalPages := (total + int64(req.PageSize) - 1) / int64(req.PageSize)
	
	return &GetPlayerPetsResponse{
		Pets:       petResponses,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// 私有方法

// validateCreatePetRequest 验证创建宠物请求
func (s *PetApplicationService) validateCreatePetRequest(req *CreatePetRequest) error {
	if req.OwnerID == 0 {
		return fmt.Errorf("owner ID is required")
	}
	if req.PetType == "" {
		return fmt.Errorf("pet type is required")
	}
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(req.Name) > 50 {
		return fmt.Errorf("name too long (max 50 characters)")
	}
	return nil
}

// validateFeedPetRequest 验证喂养宠物请求
func (s *PetApplicationService) validateFeedPetRequest(req *FeedPetRequest) error {
	if req.PetID == "" {
		return fmt.Errorf("pet ID is required")
	}
	if req.FoodType == "" {
		return fmt.Errorf("food type is required")
	}
	if req.Amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	return nil
}

// validateTrainPetRequest 验证训练宠物请求
func (s *PetApplicationService) validateTrainPetRequest(req *TrainPetRequest) error {
	if req.PetID == "" {
		return fmt.Errorf("pet ID is required")
	}
	if req.TrainingType == "" {
		return fmt.Errorf("training type is required")
	}
	if req.Duration <= 0 {
		return fmt.Errorf("duration must be positive")
	}
	return nil
}

// parseRarity 解析稀有度
func (s *PetApplicationService) parseRarity(rarityStr string) (pet.PetRarity, error) {
	switch rarityStr {
	case "common":
		return pet.RarityCommon, nil
	case "uncommon":
		return pet.RarityUncommon, nil
	case "rare":
		return pet.RarityRare, nil
	case "epic":
		return pet.RarityEpic, nil
	case "legendary":
		return pet.RarityLegendary, nil
	default:
		return pet.RarityCommon, fmt.Errorf("unknown rarity: %s", rarityStr)
	}
}

// parseSource 解析来源
func (s *PetApplicationService) parseSource(sourceStr string) (pet.PetSource, error) {
	switch sourceStr {
	case "shop":
		return pet.SourceShop, nil
	case "wild":
		return pet.SourceWild, nil
	case "breed":
		return pet.SourceBreed, nil
	case "event":
		return pet.SourceEvent, nil
	case "gift":
		return pet.SourceGift, nil
	default:
		return pet.SourceShop, fmt.Errorf("unknown source: %s", sourceStr)
	}
}

// calculateFoodExp 计算食物经验值
func (s *PetApplicationService) calculateFoodExp(foodType string, amount int32) int64 {
	baseExp := map[string]int64{
		"basic_food":   10,
		"premium_food": 25,
		"luxury_food":  50,
		"special_food": 100,
	}
	
	exp, exists := baseExp[foodType]
	if !exists {
		exp = 10 // 默认经验值
	}
	
	return exp * int64(amount)
}

// calculateTrainingGains 计算训练收益
func (s *PetApplicationService) calculateTrainingGains(trainingType string, duration int32, level int32) map[string]int32 {
	gains := make(map[string]int32)
	
	baseGains := map[string]map[string]int32{
		"strength": {
			"attack":  2,
			"defense": 1,
		},
		"agility": {
			"speed":    2,
			"accuracy": 1,
		},
		"intelligence": {
			"magic_attack": 2,
			"mana":         1,
		},
	}
	
	if baseGain, exists := baseGains[trainingType]; exists {
		for attr, gain := range baseGain {
			// 基础收益 * 时长倍数 * 等级倍数
			multiplier := float64(duration) / 60.0 * (1.0 + float64(level)*0.1)
			gains[attr] = int32(float64(gain) * multiplier)
		}
	}
	
	return gains
}

// calculateTrainingExp 计算训练经验值
func (s *PetApplicationService) calculateTrainingExp(trainingType string, duration int32) int64 {
	baseExp := int64(5) // 每分钟5经验
	return baseExp * int64(duration)
}

// checkSkillLearning 检查技能学习
func (s *PetApplicationService) checkSkillLearning(petAggregate *pet.PetAggregate, trainingType string) []string {
	skills := make([]string, 0)
	
	// 简单的技能学习逻辑
	level := petAggregate.GetLevel()
	learnedSkills := petAggregate.GetSkills()
	
	// 根据等级和训练类型判断可学习的技能
	potentialSkills := map[string][]string{
		"strength": {"power_strike", "berserker_rage", "iron_defense"},
		"agility":  {"quick_attack", "dodge", "critical_strike"},
		"intelligence": {"magic_missile", "heal", "mana_shield"},
	}
	
	if skillList, exists := potentialSkills[trainingType]; exists {
		for i, skill := range skillList {
			requiredLevel := int32((i + 1) * 10) // 10, 20, 30级学习
			if level >= requiredLevel {
				// 检查是否已学会
				alreadyLearned := false
				for _, learned := range learnedSkills {
					if learned == skill {
						alreadyLearned = true
						break
					}
				}
				if !alreadyLearned {
					skills = append(skills, skill)
				}
			}
		}
	}
	
	return skills
}