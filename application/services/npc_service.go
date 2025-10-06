package services

import (
	"context"
	"fmt"
	"time"

	"greatestworks/internal/domain/npc"
)

// NPCService NPC应用服务
type NPCService struct {
	npcRepo          npc.NPCRepository
	dialogueRepo     npc.DialogueRepository
	questRepo        npc.QuestRepository
	shopRepo         npc.ShopRepository
	relationshipRepo npc.RelationshipRepository
	statisticsRepo   npc.NPCStatisticsRepository
	cacheRepo        npc.NPCCacheRepository
	npcService       *npc.NPCService
}

// NewNPCService 创建NPC应用服务
func NewNPCService(
	npcRepo npc.NPCRepository,
	dialogueRepo npc.DialogueRepository,
	questRepo npc.QuestRepository,
	shopRepo npc.ShopRepository,
	relationshipRepo npc.RelationshipRepository,
	statisticsRepo npc.NPCStatisticsRepository,
	cacheRepo npc.NPCCacheRepository,
	npcService *npc.NPCService,
) *NPCService {
	return &NPCService{
		npcRepo:          npcRepo,
		dialogueRepo:     dialogueRepo,
		questRepo:        questRepo,
		shopRepo:         shopRepo,
		relationshipRepo: relationshipRepo,
		statisticsRepo:   statisticsRepo,
		cacheRepo:        cacheRepo,
		npcService:       npcService,
	}
}

// GetNPCInfo 获取NPC信息
func (s *NPCService) GetNPCInfo(ctx context.Context, npcID string) (*NPCDTO, error) {
	// 先从缓存获取
	cachedNPC, err := s.cacheRepo.GetNPC(npcID)
	if err == nil && cachedNPC != nil {
		return s.buildNPCDTO(cachedNPC), nil
	}

	// 从数据库获取
	npcAggregate, err := s.npcRepo.FindByID(npcID)
	if err != nil {
		return nil, fmt.Errorf("failed to get NPC info: %w", err)
	}

	// 更新缓存
	if err := s.cacheRepo.SetNPC(npcID, npcAggregate, time.Hour); err != nil {
		// 缓存更新失败不影响主流程
		// TODO: 添加日志记录
	}

	return s.buildNPCDTO(npcAggregate), nil
}

// GetNearbyNPCs 获取附近的NPC
func (s *NPCService) GetNearbyNPCs(ctx context.Context, playerID string, location *npc.Location, radius float64) ([]*NPCDTO, error) {
	// 先从缓存获取
	cachedNPCs, err := s.cacheRepo.GetLocationIndex(location.GetRegion())
	if err == nil && len(cachedNPCs) > 0 {
		// 过滤距离
		nearbyNPCs := s.filterNPCsByDistance(cachedNPCs, location, radius)
		return s.buildNPCDTOs(nearbyNPCs), nil
	}

	// 从数据库获取
	nearbyNPCs, err := s.npcRepo.FindByLocation(location, radius)
	if err != nil {
		return nil, fmt.Errorf("failed to get nearby NPCs: %w", err)
	}

	// 更新缓存
	if err := s.cacheRepo.SetLocationIndex(location.GetRegion(), nearbyNPCs, time.Minute*30); err != nil {
		// 缓存更新失败不影响主流程
		// TODO: 添加日志记录
	}

	return s.buildNPCDTOs(nearbyNPCs), nil
}

// StartDialogue 开始对话
func (s *NPCService) StartDialogue(ctx context.Context, playerID string, npcID string) (*DialogueSessionDTO, error) {
	// 获取NPC信息
	npcAggregate, err := s.npcRepo.FindByID(npcID)
	if err != nil {
		return nil, fmt.Errorf("failed to get NPC info: %w", err)
	}

	// 检查是否已有对话会话
	existingSession, err := s.cacheRepo.GetSession(npcID, playerID)
	if err == nil && existingSession != nil {
		return s.buildDialogueSessionDTO(existingSession), nil
	}

	// 开始新对话
	session, err := s.npcService.StartDialogue(playerID, npcAggregate)
	if err != nil {
		return nil, fmt.Errorf("failed to start dialogue: %w", err)
	}

	// 缓存对话会话
	if err := s.cacheRepo.SetSession(npcID, playerID, session, time.Hour); err != nil {
		// 缓存失败不影响主流程
		// TODO: 添加日志记录
	}

	return s.buildDialogueSessionDTO(session), nil
}

// ContinueDialogue 继续对话
func (s *NPCService) ContinueDialogue(ctx context.Context, playerID string, npcID string, choiceID string) (*DialogueResponseDTO, error) {
	// 获取对话会话
	session, err := s.cacheRepo.GetSession(npcID, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dialogue session: %w", err)
	}

	if session == nil {
		return nil, npc.ErrDialogueSessionNotFound
	}

	// 处理对话选择
	response, err := s.npcService.ProcessDialogueChoice(session, choiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to process dialogue choice: %w", err)
	}

	// 更新会话缓存
	if err := s.cacheRepo.SetSession(npcID, playerID, session, time.Hour); err != nil {
		// 缓存更新失败不影响主流程
		// TODO: 添加日志记录
	}

	return s.buildDialogueResponseDTO(response), nil
}

// EndDialogue 结束对话
func (s *NPCService) EndDialogue(ctx context.Context, playerID string, npcID string) error {
	// 获取对话会话
	session, err := s.cacheRepo.GetSession(npcID, playerID)
	if err != nil {
		return fmt.Errorf("failed to get dialogue session: %w", err)
	}

	if session != nil {
		// 结束对话
		if err := s.npcService.EndDialogue(session); err != nil {
			return fmt.Errorf("failed to end dialogue: %w", err)
		}

		// 更新统计数据
		if err := s.updateDialogueStatistics(ctx, playerID, npcID, session); err != nil {
			// 统计更新失败不影响主流程
			// TODO: 添加日志记录
		}
	}

	// 清除会话缓存
	if err := s.cacheRepo.DeleteSession(npcID, playerID); err != nil {
		// 缓存清除失败不影响主流程
		// TODO: 添加日志记录
	}

	return nil
}

// GetAvailableQuests 获取可用任务
func (s *NPCService) GetAvailableQuests(ctx context.Context, playerID string, npcID string) ([]*QuestDTO, error) {
	// 获取NPC的任务
	quests, err := s.questRepo.FindByNPC(npcID)
	if err != nil {
		return nil, fmt.Errorf("failed to get NPC quests: %w", err)
	}

	// 过滤可用任务
	availableQuests := make([]*npc.Quest, 0)
	for _, quest := range quests {
		if s.npcService.IsQuestAvailable(playerID, quest) {
			availableQuests = append(availableQuests, quest)
		}
	}

	return s.buildQuestDTOs(availableQuests), nil
}

// AcceptQuest 接受任务
func (s *NPCService) AcceptQuest(ctx context.Context, playerID string, questID string) (*QuestInstanceDTO, error) {
	// 获取任务信息
	quest, err := s.questRepo.FindByID(questID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quest info: %w", err)
	}

	// 接受任务
	questInstance, err := s.npcService.AcceptQuest(playerID, quest)
	if err != nil {
		return nil, fmt.Errorf("failed to accept quest: %w", err)
	}

	// 保存任务实例
	if err := s.questRepo.SaveInstance(questInstance); err != nil {
		return nil, fmt.Errorf("failed to save quest instance: %w", err)
	}

	return s.buildQuestInstanceDTO(questInstance), nil
}

// CompleteQuest 完成任务
func (s *NPCService) CompleteQuest(ctx context.Context, playerID string, questID string) (*QuestRewardDTO, error) {
	// 获取任务实例
	questInstance, err := s.questRepo.FindInstance(questID, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quest instance: %w", err)
	}

	// 完成任务
	reward, err := s.npcService.CompleteQuest(questInstance)
	if err != nil {
		return nil, fmt.Errorf("failed to complete quest: %w", err)
	}

	// 更新任务实例
	if err := s.questRepo.UpdateInstance(questInstance); err != nil {
		return nil, fmt.Errorf("failed to update quest instance: %w", err)
	}

	// 更新统计数据
	if err := s.updateQuestStatistics(ctx, playerID, questInstance, reward); err != nil {
		// 统计更新失败不影响主流程
		// TODO: 添加日志记录
	}

	return s.buildQuestRewardDTO(reward), nil
}

// GetShopInfo 获取商店信息
func (s *NPCService) GetShopInfo(ctx context.Context, npcID string) (*ShopDTO, error) {
	// 获取商店信息
	shop, err := s.shopRepo.FindByNPC(npcID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop info: %w", err)
	}

	return s.buildShopDTO(shop), nil
}

// BuyItem 购买物品
func (s *NPCService) BuyItem(ctx context.Context, playerID string, shopID string, itemID string, quantity int) (*TradeResultDTO, error) {
	// 获取商店信息
	shop, err := s.shopRepo.FindByID(shopID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop info: %w", err)
	}

	// 执行购买
	tradeResult, err := s.npcService.BuyItem(playerID, shop, itemID, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to buy item: %w", err)
	}

	// 保存交易记录
	tradeRecord := npc.NewTradeRecord(shopID, playerID, itemID, quantity, tradeResult.Price)
	if err := s.shopRepo.SaveTradeRecord(tradeRecord); err != nil {
		// 交易记录保存失败不影响主流程
		// TODO: 添加日志记录
	}

	// 更新商店
	if err := s.shopRepo.Update(shop); err != nil {
		return nil, fmt.Errorf("failed to update shop: %w", err)
	}

	return s.buildTradeResultDTO(tradeResult), nil
}

// GetRelationship 获取关系信息
func (s *NPCService) GetRelationship(ctx context.Context, playerID string, npcID string) (*RelationshipDTO, error) {
	// 先从缓存获取
	cachedRelationship, err := s.cacheRepo.GetRelationship(playerID, npcID)
	if err == nil && cachedRelationship != nil {
		return s.buildRelationshipDTO(cachedRelationship), nil
	}

	// 从数据库获取
	relationship, err := s.relationshipRepo.FindByID(playerID, npcID)
	if err != nil {
		return nil, fmt.Errorf("failed to get relationship: %w", err)
	}

	// 更新缓存
	if err := s.cacheRepo.SetRelationship(playerID, npcID, relationship, time.Hour*2); err != nil {
		// 缓存更新失败不影响主流程
		// TODO: 添加日志记录
	}

	return s.buildRelationshipDTO(relationship), nil
}

// UpdateRelationship 更新关系
func (s *NPCService) UpdateRelationship(ctx context.Context, playerID string, npcID string, changeType npc.RelationshipChangeType, value int, reason string) error {
	// 获取关系信息
	relationship, err := s.relationshipRepo.FindByID(playerID, npcID)
	if err != nil && !npc.IsNotFoundError(err) {
		return fmt.Errorf("failed to get relationship: %w", err)
	}

	if relationship == nil {
		// 创建新关系
		relationship = npc.NewRelationship(playerID, npcID)
	}

	// 更新关系值
	oldValue := relationship.GetValue()
	oldLevel := relationship.GetLevel()

	if err := relationship.ChangeValue(value, reason); err != nil {
		return fmt.Errorf("failed to change relationship value: %w", err)
	}

	// 保存关系
	if err := s.relationshipRepo.Save(relationship); err != nil {
		return fmt.Errorf("failed to save relationship: %w", err)
	}

	// 记录关系变化事件
	if relationship.GetValue() != oldValue {
		event := npc.NewRelationshipChangedEvent(
			npcID, playerID,
			oldValue, relationship.GetValue(),
			oldLevel, relationship.GetLevel(),
			changeType, reason,
		)
		// TODO: 发布事件
		_ = event
	}

	// 清除缓存
	if err := s.cacheRepo.DeleteRelationship(playerID, npcID); err != nil {
		// 缓存清除失败不影响主流程
		// TODO: 添加日志记录
	}

	return nil
}

// GetNPCStatistics 获取NPC统计
func (s *NPCService) GetNPCStatistics(ctx context.Context, npcID string) (*NPCStatisticsDTO, error) {
	stats, err := s.statisticsRepo.FindStatistics(npcID)
	if err != nil {
		return nil, fmt.Errorf("failed to get NPC statistics: %w", err)
	}

	return s.buildNPCStatisticsDTO(stats), nil
}

// 私有方法

// filterNPCsByDistance 按距离过滤NPC
func (s *NPCService) filterNPCsByDistance(npcs []*npc.NPCAggregate, location *npc.Location, radius float64) []*npc.NPCAggregate {
	filtered := make([]*npc.NPCAggregate, 0)
	for _, npcAgg := range npcs {
		if npcAgg.GetLocation().DistanceTo(location) <= radius {
			filtered = append(filtered, npcAgg)
		}
	}
	return filtered
}

// updateDialogueStatistics 更新对话统计
func (s *NPCService) updateDialogueStatistics(ctx context.Context, playerID string, npcID string, session *npc.DialogueSession) error {
	stats, err := s.statisticsRepo.FindStatistics(npcID)
	if err != nil && !npc.IsNotFoundError(err) {
		return err
	}

	if stats == nil {
		stats = npc.NewNPCStatistics(npcID)
	}

	// 更新统计数据
	stats.AddDialogueSession(playerID, session.GetDuration())
	stats.UpdateLastInteractionTime(session.GetEndTime())

	// 保存统计数据
	return s.statisticsRepo.SaveStatistics(stats)
}

// updateQuestStatistics 更新任务统计
func (s *NPCService) updateQuestStatistics(ctx context.Context, playerID string, questInstance *npc.QuestInstance, reward *npc.QuestReward) error {
	stats, err := s.statisticsRepo.FindStatistics(questInstance.GetNPCID())
	if err != nil && !npc.IsNotFoundError(err) {
		return err
	}

	if stats == nil {
		stats = npc.NewNPCStatistics(questInstance.GetNPCID())
	}

	// 更新统计数据
	stats.AddQuestCompletion(playerID, questInstance.GetQuestID(), reward.GetTotalValue())
	stats.UpdateLastInteractionTime(questInstance.GetCompletedAt())

	// 保存统计数据
	return s.statisticsRepo.SaveStatistics(stats)
}

// 构建DTO方法

// buildNPCDTO 构建NPC DTO
func (s *NPCService) buildNPCDTO(npcAggregate *npc.NPCAggregate) *NPCDTO {
	return &NPCDTO{
		ID:          npcAggregate.GetID(),
		Name:        npcAggregate.GetName(),
		Description: npcAggregate.GetDescription(),
		Type:        string(npcAggregate.GetType()),
		Status:      string(npcAggregate.GetStatus()),
		Location:    s.buildLocationDTO(npcAggregate.GetLocation()),
		Attributes:  s.buildAttributesDTO(npcAggregate.GetAttributes()),
		Behavior:    s.buildBehaviorDTO(npcAggregate.GetBehavior()),
		HasDialogue: npcAggregate.HasDialogue(),
		HasQuests:   npcAggregate.HasQuests(),
		HasShop:     npcAggregate.HasShop(),
		CreatedAt:   npcAggregate.GetCreatedAt(),
		UpdatedAt:   npcAggregate.GetUpdatedAt(),
	}
}

// buildNPCDTOs 构建NPC DTO列表
func (s *NPCService) buildNPCDTOs(npcs []*npc.NPCAggregate) []*NPCDTO {
	dtos := make([]*NPCDTO, len(npcs))
	for i, npcAgg := range npcs {
		dtos[i] = s.buildNPCDTO(npcAgg)
	}
	return dtos
}

// buildLocationDTO 构建位置DTO
func (s *NPCService) buildLocationDTO(location *npc.Location) *LocationDTO {
	return &LocationDTO{
		X:      location.GetX(),
		Y:      location.GetY(),
		Z:      location.GetZ(),
		Region: location.GetRegion(),
		Zone:   location.GetZone(),
	}
}

// buildAttributesDTO 构建属性DTO
func (s *NPCService) buildAttributesDTO(attributes *npc.NPCAttributes) *NPCAttributesDTO {
	return &NPCAttributesDTO{
		Level:        attributes.GetLevel(),
		Health:       attributes.GetHealth(),
		MaxHealth:    attributes.GetMaxHealth(),
		Attack:       attributes.GetAttack(),
		Defense:      attributes.GetDefense(),
		Speed:        attributes.GetSpeed(),
		Intelligence: attributes.GetIntelligence(),
	}
}

// buildBehaviorDTO 构建行为DTO
func (s *NPCService) buildBehaviorDTO(behavior *npc.NPCBehavior) *NPCBehaviorDTO {
	return &NPCBehaviorDTO{
		CurrentAction: string(behavior.GetCurrentAction()),
		NextAction:    string(behavior.GetNextAction()),
		Cooldown:      behavior.GetCooldown(),
		IsActive:      behavior.IsActive(),
	}
}

// buildDialogueSessionDTO 构建对话会话DTO
func (s *NPCService) buildDialogueSessionDTO(session *npc.DialogueSession) *DialogueSessionDTO {
	return &DialogueSessionDTO{
		SessionID:     session.GetID(),
		NPCID:         session.GetNPCID(),
		PlayerID:      session.GetPlayerID(),
		CurrentNodeID: session.GetCurrentNodeID(),
		StartTime:     session.GetStartTime(),
		IsActive:      session.IsActive(),
		Context:       session.GetContext(),
	}
}

// buildDialogueResponseDTO 构建对话响应DTO
func (s *NPCService) buildDialogueResponseDTO(response *npc.DialogueResponse) *DialogueResponseDTO {
	return &DialogueResponseDTO{
		NodeID:     response.GetNodeID(),
		Text:       response.GetText(),
		Choices:    s.buildDialogueChoiceDTOs(response.GetChoices()),
		Actions:    response.GetActions(),
		IsEnd:      response.IsEnd(),
		NextNodeID: response.GetNextNodeID(),
	}
}

// buildDialogueChoiceDTOs 构建对话选择DTO列表
func (s *NPCService) buildDialogueChoiceDTOs(choices []*npc.DialogueChoice) []*DialogueChoiceDTO {
	dtos := make([]*DialogueChoiceDTO, len(choices))
	for i, choice := range choices {
		dtos[i] = &DialogueChoiceDTO{
			ID:          choice.GetID(),
			Text:        choice.GetText(),
			Condition:   choice.GetCondition(),
			NextNodeID:  choice.GetNextNodeID(),
			IsAvailable: choice.IsAvailable(),
		}
	}
	return dtos
}

// buildQuestDTOs 构建任务DTO列表
func (s *NPCService) buildQuestDTOs(quests []*npc.Quest) []*QuestDTO {
	dtos := make([]*QuestDTO, len(quests))
	for i, quest := range quests {
		dtos[i] = &QuestDTO{
			ID:            quest.GetID(),
			Name:          quest.GetName(),
			Description:   quest.GetDescription(),
			Type:          string(quest.GetType()),
			RequiredLevel: quest.GetRequiredLevel(),
			Rewards:       quest.GetRewards(),
			Objectives:    quest.GetObjectives(),
			IsRepeatable:  quest.IsRepeatable(),
			Cooldown:      quest.GetCooldown(),
		}
	}
	return dtos
}

// buildQuestInstanceDTO 构建任务实例DTO
func (s *NPCService) buildQuestInstanceDTO(instance *npc.QuestInstance) *QuestInstanceDTO {
	return &QuestInstanceDTO{
		ID:          instance.GetID(),
		QuestID:     instance.GetQuestID(),
		PlayerID:    instance.GetPlayerID(),
		Status:      string(instance.GetStatus()),
		Progress:    instance.GetProgress(),
		StartTime:   instance.GetStartTime(),
		EndTime:     instance.GetEndTime(),
		IsCompleted: instance.IsCompleted(),
	}
}

// buildQuestRewardDTO 构建任务奖励DTO
func (s *NPCService) buildQuestRewardDTO(reward *npc.QuestReward) *QuestRewardDTO {
	return &QuestRewardDTO{
		Experience: reward.GetExperience(),
		Items:      reward.GetItems(),
		Gold:       reward.GetGold(),
		TotalValue: reward.GetTotalValue(),
	}
}

// buildShopDTO 构建商店DTO
func (s *NPCService) buildShopDTO(shop *npc.Shop) *ShopDTO {
	return &ShopDTO{
		ID:          shop.GetID(),
		NPCID:       shop.GetNPCID(),
		Name:        shop.GetName(),
		Description: shop.GetDescription(),
		Items:       s.buildShopItemDTOs(shop.GetItems()),
		IsOpen:      shop.IsOpen(),
		Schedule:    s.buildShopScheduleDTO(shop.GetSchedule()),
	}
}

// buildShopItemDTOs 构建商店物品DTO列表
func (s *NPCService) buildShopItemDTOs(items []*npc.ShopItem) []*ShopItemDTO {
	dtos := make([]*ShopItemDTO, len(items))
	for i, item := range items {
		dtos[i] = &ShopItemDTO{
			ID:          item.GetID(),
			Name:        item.GetName(),
			Description: item.GetDescription(),
			Price:       item.GetPrice(),
			Stock:       item.GetStock(),
			MaxStock:    item.GetMaxStock(),
			IsAvailable: item.IsAvailable(),
		}
	}
	return dtos
}

// buildShopScheduleDTO 构建商店日程DTO
func (s *NPCService) buildShopScheduleDTO(schedule *npc.ShopSchedule) *ShopScheduleDTO {
	return &ShopScheduleDTO{
		OpenTime:  schedule.GetOpenTime(),
		CloseTime: schedule.GetCloseTime(),
		IsOpen24H: schedule.IsOpen24H(),
		Weekdays:  schedule.GetWeekdays(),
	}
}

// buildTradeResultDTO 构建交易结果DTO
func (s *NPCService) buildTradeResultDTO(result *npc.TradeResult) *TradeResultDTO {
	return &TradeResultDTO{
		ItemID:     result.GetItemID(),
		Quantity:   result.GetQuantity(),
		Price:      result.GetPrice(),
		TotalPrice: result.GetTotalPrice(),
		Success:    result.IsSuccess(),
		Message:    result.GetMessage(),
	}
}

// buildRelationshipDTO 构建关系DTO
func (s *NPCService) buildRelationshipDTO(relationship *npc.Relationship) *RelationshipDTO {
	return &RelationshipDTO{
		PlayerID:    relationship.GetPlayerID(),
		NPCID:       relationship.GetNPCID(),
		Value:       relationship.GetValue(),
		Level:       string(relationship.GetLevel()),
		LastChanged: relationship.GetLastChanged(),
		IsLocked:    relationship.IsLocked(),
	}
}

// buildNPCStatisticsDTO 构建NPC统计DTO
func (s *NPCService) buildNPCStatisticsDTO(stats *npc.NPCStatistics) *NPCStatisticsDTO {
	return &NPCStatisticsDTO{
		NPCID:                  stats.GetNPCID(),
		TotalInteractions:      stats.GetTotalInteractions(),
		DialogueCount:          stats.GetDialogueCount(),
		QuestCount:             stats.GetQuestCount(),
		TradeCount:             stats.GetTradeCount(),
		UniqueVisitors:         stats.GetUniqueVisitors(),
		AverageInteractionTime: stats.GetAverageInteractionTime(),
		LastInteractionTime:    stats.GetLastInteractionTime(),
		PopularityScore:        stats.GetPopularityScore(),
	}
}

// DTO 定义

// NPCDTO NPC DTO
type NPCDTO struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        string            `json:"type"`
	Status      string            `json:"status"`
	Location    *LocationDTO      `json:"location"`
	Attributes  *NPCAttributesDTO `json:"attributes"`
	Behavior    *NPCBehaviorDTO   `json:"behavior"`
	HasDialogue bool              `json:"has_dialogue"`
	HasQuests   bool              `json:"has_quests"`
	HasShop     bool              `json:"has_shop"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// LocationDTO 位置DTO
type LocationDTO struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Z      float64 `json:"z"`
	Region string  `json:"region"`
	Zone   string  `json:"zone"`
}

// NPCAttributesDTO NPC属性DTO
type NPCAttributesDTO struct {
	Level        int     `json:"level"`
	Health       int     `json:"health"`
	MaxHealth    int     `json:"max_health"`
	Attack       int     `json:"attack"`
	Defense      int     `json:"defense"`
	Speed        float64 `json:"speed"`
	Intelligence int     `json:"intelligence"`
}

// NPCBehaviorDTO NPC行为DTO
type NPCBehaviorDTO struct {
	CurrentAction string        `json:"current_action"`
	NextAction    string        `json:"next_action"`
	Cooldown      time.Duration `json:"cooldown"`
	IsActive      bool          `json:"is_active"`
}

// DialogueSessionDTO 对话会话DTO
type DialogueSessionDTO struct {
	SessionID     string                 `json:"session_id"`
	NPCID         string                 `json:"npc_id"`
	PlayerID      string                 `json:"player_id"`
	CurrentNodeID string                 `json:"current_node_id"`
	StartTime     time.Time              `json:"start_time"`
	IsActive      bool                   `json:"is_active"`
	Context       map[string]interface{} `json:"context"`
}

// DialogueResponseDTO 对话响应DTO
type DialogueResponseDTO struct {
	NodeID     string               `json:"node_id"`
	Text       string               `json:"text"`
	Choices    []*DialogueChoiceDTO `json:"choices"`
	Actions    []string             `json:"actions"`
	IsEnd      bool                 `json:"is_end"`
	NextNodeID string               `json:"next_node_id,omitempty"`
}

// DialogueChoiceDTO 对话选择DTO
type DialogueChoiceDTO struct {
	ID          string `json:"id"`
	Text        string `json:"text"`
	Condition   string `json:"condition,omitempty"`
	NextNodeID  string `json:"next_node_id"`
	IsAvailable bool   `json:"is_available"`
}

// QuestDTO 任务DTO
type QuestDTO struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Type          string                 `json:"type"`
	RequiredLevel int                    `json:"required_level"`
	Rewards       map[string]interface{} `json:"rewards"`
	Objectives    []string               `json:"objectives"`
	IsRepeatable  bool                   `json:"is_repeatable"`
	Cooldown      time.Duration          `json:"cooldown"`
}

// QuestInstanceDTO 任务实例DTO
type QuestInstanceDTO struct {
	ID          string         `json:"id"`
	QuestID     string         `json:"quest_id"`
	PlayerID    string         `json:"player_id"`
	Status      string         `json:"status"`
	Progress    map[string]int `json:"progress"`
	StartTime   time.Time      `json:"start_time"`
	EndTime     time.Time      `json:"end_time"`
	IsCompleted bool           `json:"is_completed"`
}

// QuestRewardDTO 任务奖励DTO
type QuestRewardDTO struct {
	Experience int64          `json:"experience"`
	Items      map[string]int `json:"items"`
	Gold       int64          `json:"gold"`
	TotalValue int64          `json:"total_value"`
}

// ShopDTO 商店DTO
type ShopDTO struct {
	ID          string           `json:"id"`
	NPCID       string           `json:"npc_id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Items       []*ShopItemDTO   `json:"items"`
	IsOpen      bool             `json:"is_open"`
	Schedule    *ShopScheduleDTO `json:"schedule"`
}

// ShopItemDTO 商店物品DTO
type ShopItemDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Stock       int    `json:"stock"`
	MaxStock    int    `json:"max_stock"`
	IsAvailable bool   `json:"is_available"`
}

// ShopScheduleDTO 商店日程DTO
type ShopScheduleDTO struct {
	OpenTime  time.Time `json:"open_time"`
	CloseTime time.Time `json:"close_time"`
	IsOpen24H bool      `json:"is_open_24h"`
	Weekdays  []int     `json:"weekdays"`
}

// TradeResultDTO 交易结果DTO
type TradeResultDTO struct {
	ItemID     string `json:"item_id"`
	Quantity   int    `json:"quantity"`
	Price      int64  `json:"price"`
	TotalPrice int64  `json:"total_price"`
	Success    bool   `json:"success"`
	Message    string `json:"message"`
}

// RelationshipDTO 关系DTO
type RelationshipDTO struct {
	PlayerID    string    `json:"player_id"`
	NPCID       string    `json:"npc_id"`
	Value       int       `json:"value"`
	Level       string    `json:"level"`
	LastChanged time.Time `json:"last_changed"`
	IsLocked    bool      `json:"is_locked"`
}

// NPCStatisticsDTO NPC统计DTO
type NPCStatisticsDTO struct {
	NPCID                  string        `json:"npc_id"`
	TotalInteractions      int64         `json:"total_interactions"`
	DialogueCount          int64         `json:"dialogue_count"`
	QuestCount             int64         `json:"quest_count"`
	TradeCount             int64         `json:"trade_count"`
	UniqueVisitors         int64         `json:"unique_visitors"`
	AverageInteractionTime time.Duration `json:"average_interaction_time"`
	LastInteractionTime    time.Time     `json:"last_interaction_time"`
	PopularityScore        float64       `json:"popularity_score"`
}
