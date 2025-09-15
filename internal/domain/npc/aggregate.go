package npc

import (
	"fmt"
	"time"
)

// NPCAggregate NPC聚合根
type NPCAggregate struct {
	id          string
	name        string
	description string
	npcType     NPCType
	status      NPCStatus
	location    *Location
	attributes  *NPCAttributes
	behavior    *NPCBehavior
	dialogues   map[string]*Dialogue
	quests      map[string]*Quest
	shop        *Shop
	relationships map[string]*Relationship
	schedule    *NPCSchedule
	createdAt   time.Time
	updatedAt   time.Time
	version     int
	events      []DomainEvent
}

// NewNPCAggregate 创建NPC聚合根
func NewNPCAggregate(id, name, description string, npcType NPCType) *NPCAggregate {
	now := time.Now()
	return &NPCAggregate{
		id:          id,
		name:        name,
		description: description,
		npcType:     npcType,
		status:      NPCStatusActive,
		location:    NewLocation(0, 0, 0, "default", "default"),
		attributes:  NewNPCAttributes(),
		behavior:    NewNPCBehavior(),
		dialogues:   make(map[string]*Dialogue),
		quests:      make(map[string]*Quest),
		relationships: make(map[string]*Relationship),
		schedule:    NewNPCSchedule(),
		createdAt:   now,
		updatedAt:   now,
		version:     1,
		events:      make([]DomainEvent, 0),
	}
}

// GetID 获取ID
func (n *NPCAggregate) GetID() string {
	return n.id
}

// GetName 获取名称
func (n *NPCAggregate) GetName() string {
	return n.name
}

// GetDescription 获取描述
func (n *NPCAggregate) GetDescription() string {
	return n.description
}

// GetType 获取类型
func (n *NPCAggregate) GetType() NPCType {
	return n.npcType
}

// GetStatus 获取状态
func (n *NPCAggregate) GetStatus() NPCStatus {
	return n.status
}

// GetLocation 获取位置
func (n *NPCAggregate) GetLocation() *Location {
	return n.location
}

// GetAttributes 获取属性
func (n *NPCAggregate) GetAttributes() *NPCAttributes {
	return n.attributes
}

// GetBehavior 获取行为
func (n *NPCAggregate) GetBehavior() *NPCBehavior {
	return n.behavior
}

// GetShop 获取商店
func (n *NPCAggregate) GetShop() *Shop {
	return n.shop
}

// GetSchedule 获取日程
func (n *NPCAggregate) GetSchedule() *NPCSchedule {
	return n.schedule
}

// GetVersion 获取版本
func (n *NPCAggregate) GetVersion() int {
	return n.version
}

// GetEvents 获取领域事件
func (n *NPCAggregate) GetEvents() []DomainEvent {
	return n.events
}

// ClearEvents 清除领域事件
func (n *NPCAggregate) ClearEvents() {
	n.events = make([]DomainEvent, 0)
}

// SetName 设置名称
func (n *NPCAggregate) SetName(name string) error {
	if name == "" {
		return ErrInvalidNPCName
	}
	
	oldName := n.name
	n.name = name
	n.updatedAt = time.Now()
	n.version++
	
	// 发布名称变更事件
	event := NewNPCNameChangedEvent(n.id, oldName, name)
	n.addEvent(event)
	
	return nil
}

// SetDescription 设置描述
func (n *NPCAggregate) SetDescription(description string) {
	n.description = description
	n.updatedAt = time.Now()
	n.version++
}

// SetStatus 设置状态
func (n *NPCAggregate) SetStatus(status NPCStatus) error {
	if !status.IsValid() {
		return ErrInvalidNPCStatus
	}
	
	oldStatus := n.status
	n.status = status
	n.updatedAt = time.Now()
	n.version++
	
	// 发布状态变更事件
	event := NewNPCStatusChangedEvent(n.id, oldStatus, status)
	n.addEvent(event)
	
	return nil
}

// MoveTo 移动到指定位置
func (n *NPCAggregate) MoveTo(location *Location) error {
	if location == nil {
		return ErrInvalidLocation
	}
	
	// 检查是否可以移动
	if !n.CanMove() {
		return ErrNPCCannotMove
	}
	
	oldLocation := n.location
	n.location = location
	n.updatedAt = time.Now()
	n.version++
	
	// 发布移动事件
	event := NewNPCMovedEvent(n.id, oldLocation, location)
	n.addEvent(event)
	
	return nil
}

// CanMove 检查是否可以移动
func (n *NPCAggregate) CanMove() bool {
	return n.status == NPCStatusActive && n.behavior.CanMove()
}

// AddDialogue 添加对话
func (n *NPCAggregate) AddDialogue(dialogue *Dialogue) error {
	if dialogue == nil {
		return ErrInvalidDialogue
	}
	
	if _, exists := n.dialogues[dialogue.GetID()]; exists {
		return ErrDialogueAlreadyExists
	}
	
	n.dialogues[dialogue.GetID()] = dialogue
	n.updatedAt = time.Now()
	n.version++
	
	// 发布对话添加事件
	event := NewDialogueAddedEvent(n.id, dialogue.GetID(), dialogue.GetType())
	n.addEvent(event)
	
	return nil
}

// RemoveDialogue 移除对话
func (n *NPCAggregate) RemoveDialogue(dialogueID string) error {
	dialogue, exists := n.dialogues[dialogueID]
	if !exists {
		return ErrDialogueNotFound
	}
	
	delete(n.dialogues, dialogueID)
	n.updatedAt = time.Now()
	n.version++
	
	// 发布对话移除事件
	event := NewDialogueRemovedEvent(n.id, dialogueID, dialogue.GetType())
	n.addEvent(event)
	
	return nil
}

// StartDialogue 开始对话
func (n *NPCAggregate) StartDialogue(dialogueID, playerID string) (*DialogueSession, error) {
	dialogue, exists := n.dialogues[dialogueID]
	if !exists {
		return nil, ErrDialogueNotFound
	}
	
	// 检查对话条件
	if !dialogue.CanStart(playerID) {
		return nil, ErrDialogueConditionsNotMet
	}
	
	// 创建对话会话
	session := NewDialogueSession(n.id, dialogueID, playerID)
	
	n.updatedAt = time.Now()
	n.version++
	
	// 发布对话开始事件
	event := NewDialogueStartedEvent(n.id, dialogueID, playerID)
	n.addEvent(event)
	
	return session, nil
}

// AddQuest 添加任务
func (n *NPCAggregate) AddQuest(quest *Quest) error {
	if quest == nil {
		return ErrInvalidQuest
	}
	
	if _, exists := n.quests[quest.GetID()]; exists {
		return ErrQuestAlreadyExists
	}
	
	n.quests[quest.GetID()] = quest
	n.updatedAt = time.Now()
	n.version++
	
	// 发布任务添加事件
	event := NewQuestAddedEvent(n.id, quest.GetID(), quest.GetType())
	n.addEvent(event)
	
	return nil
}

// RemoveQuest 移除任务
func (n *NPCAggregate) RemoveQuest(questID string) error {
	quest, exists := n.quests[questID]
	if !exists {
		return ErrQuestNotFound
	}
	
	delete(n.quests, questID)
	n.updatedAt = time.Now()
	n.version++
	
	// 发布任务移除事件
	event := NewQuestRemovedEvent(n.id, questID, quest.GetType())
	n.addEvent(event)
	
	return nil
}

// GiveQuest 给予任务
func (n *NPCAggregate) GiveQuest(questID, playerID string) (*QuestInstance, error) {
	quest, exists := n.quests[questID]
	if !exists {
		return nil, ErrQuestNotFound
	}
	
	// 检查任务条件
	if !quest.CanAccept(playerID) {
		return nil, ErrQuestConditionsNotMet
	}
	
	// 创建任务实例
	instance := NewQuestInstance(questID, playerID, n.id)
	
	n.updatedAt = time.Now()
	n.version++
	
	// 发布任务给予事件
	event := NewQuestGivenEvent(n.id, questID, playerID)
	n.addEvent(event)
	
	return instance, nil
}

// SetShop 设置商店
func (n *NPCAggregate) SetShop(shop *Shop) error {
	if shop == nil {
		return ErrInvalidShop
	}
	
	// 检查NPC类型是否支持商店
	if !n.npcType.CanHaveShop() {
		return ErrNPCCannotHaveShop
	}
	
	n.shop = shop
	n.updatedAt = time.Now()
	n.version++
	
	// 发布商店设置事件
	event := NewShopSetEvent(n.id, shop.GetID())
	n.addEvent(event)
	
	return nil
}

// Trade 交易
func (n *NPCAggregate) Trade(playerID string, tradeRequest *TradeRequest) (*TradeResult, error) {
	if n.shop == nil {
		return nil, ErrNPCHasNoShop
	}
	
	// 执行交易
	result, err := n.shop.ExecuteTrade(playerID, tradeRequest)
	if err != nil {
		return nil, err
	}
	
	n.updatedAt = time.Now()
	n.version++
	
	// 发布交易事件
	event := NewTradeExecutedEvent(n.id, playerID, tradeRequest, result)
	n.addEvent(event)
	
	return result, nil
}

// UpdateRelationship 更新关系
func (n *NPCAggregate) UpdateRelationship(playerID string, change int, reason string) error {
	relationship, exists := n.relationships[playerID]
	if !exists {
		relationship = NewRelationship(playerID, n.id)
		n.relationships[playerID] = relationship
	}
	
	oldLevel := relationship.GetLevel()
	err := relationship.ChangeValue(change, reason)
	if err != nil {
		return err
	}
	
	n.updatedAt = time.Now()
	n.version++
	
	// 如果关系等级发生变化，发布事件
	if relationship.GetLevel() != oldLevel {
		event := NewRelationshipChangedEvent(n.id, playerID, oldLevel, relationship.GetLevel(), change)
		n.addEvent(event)
	}
	
	return nil
}

// GetRelationship 获取关系
func (n *NPCAggregate) GetRelationship(playerID string) *Relationship {
	return n.relationships[playerID]
}

// GetDialogue 获取对话
func (n *NPCAggregate) GetDialogue(dialogueID string) (*Dialogue, error) {
	dialogue, exists := n.dialogues[dialogueID]
	if !exists {
		return nil, ErrDialogueNotFound
	}
	return dialogue, nil
}

// GetAllDialogues 获取所有对话
func (n *NPCAggregate) GetAllDialogues() map[string]*Dialogue {
	return n.dialogues
}

// GetAvailableDialogues 获取可用对话
func (n *NPCAggregate) GetAvailableDialogues(playerID string) []*Dialogue {
	var available []*Dialogue
	for _, dialogue := range n.dialogues {
		if dialogue.CanStart(playerID) {
			available = append(available, dialogue)
		}
	}
	return available
}

// GetQuest 获取任务
func (n *NPCAggregate) GetQuest(questID string) (*Quest, error) {
	quest, exists := n.quests[questID]
	if !exists {
		return nil, ErrQuestNotFound
	}
	return quest, nil
}

// GetAllQuests 获取所有任务
func (n *NPCAggregate) GetAllQuests() map[string]*Quest {
	return n.quests
}

// GetAvailableQuests 获取可用任务
func (n *NPCAggregate) GetAvailableQuests(playerID string) []*Quest {
	var available []*Quest
	for _, quest := range n.quests {
		if quest.CanAccept(playerID) {
			available = append(available, quest)
		}
	}
	return available
}

// Update 更新NPC状态
func (n *NPCAggregate) Update(deltaTime time.Duration) {
	// 更新行为
	n.behavior.Update(deltaTime)
	
	// 更新日程
	n.schedule.Update(time.Now())
	
	// 更新商店（如果有）
	if n.shop != nil {
		n.shop.Update(deltaTime)
	}
	
	n.updatedAt = time.Now()
	n.version++
}

// IsActive 检查是否激活
func (n *NPCAggregate) IsActive() bool {
	return n.status == NPCStatusActive
}

// CanInteract 检查是否可以交互
func (n *NPCAggregate) CanInteract(playerID string) bool {
	if !n.IsActive() {
		return false
	}
	
	// 检查关系是否允许交互
	relationship := n.GetRelationship(playerID)
	if relationship != nil && relationship.GetLevel() == RelationshipLevelHostile {
		return false
	}
	
	return true
}

// GetInteractionOptions 获取交互选项
func (n *NPCAggregate) GetInteractionOptions(playerID string) []InteractionOption {
	var options []InteractionOption
	
	if !n.CanInteract(playerID) {
		return options
	}
	
	// 添加对话选项
	availableDialogues := n.GetAvailableDialogues(playerID)
	for _, dialogue := range availableDialogues {
		options = append(options, InteractionOption{
			Type:        InteractionTypeDialogue,
			ID:          dialogue.GetID(),
			Name:        dialogue.GetName(),
			Description: dialogue.GetDescription(),
		})
	}
	
	// 添加任务选项
	availableQuests := n.GetAvailableQuests(playerID)
	for _, quest := range availableQuests {
		options = append(options, InteractionOption{
			Type:        InteractionTypeQuest,
			ID:          quest.GetID(),
			Name:        quest.GetName(),
			Description: quest.GetDescription(),
		})
	}
	
	// 添加商店选项
	if n.shop != nil && n.shop.IsOpen() {
		options = append(options, InteractionOption{
			Type:        InteractionTypeShop,
			ID:          n.shop.GetID(),
			Name:        "商店",
			Description: "浏览商品",
		})
	}
	
	return options
}

// Activate 激活NPC
func (n *NPCAggregate) Activate() error {
	return n.SetStatus(NPCStatusActive)
}

// Deactivate 停用NPC
func (n *NPCAggregate) Deactivate() error {
	return n.SetStatus(NPCStatusInactive)
}

// Hide 隐藏NPC
func (n *NPCAggregate) Hide() error {
	return n.SetStatus(NPCStatusHidden)
}

// Busy 设置忙碌状态
func (n *NPCAggregate) Busy() error {
	return n.SetStatus(NPCStatusBusy)
}

// GetStatistics 获取统计信息
func (n *NPCAggregate) GetStatistics() *NPCStatistics {
	totalDialogues := len(n.dialogues)
	totalQuests := len(n.quests)
	totalRelationships := len(n.relationships)
	
	var averageRelationship float64
	if totalRelationships > 0 {
		sum := 0
		for _, rel := range n.relationships {
			sum += rel.GetValue()
		}
		averageRelationship = float64(sum) / float64(totalRelationships)
	}
	
	return &NPCStatistics{
		NPCID:               n.id,
		Name:                n.name,
		Type:                n.npcType,
		Status:              n.status,
		TotalDialogues:      totalDialogues,
		TotalQuests:         totalQuests,
		TotalRelationships:  totalRelationships,
		AverageRelationship: averageRelationship,
		CreatedAt:           n.createdAt,
		LastActiveAt:        n.updatedAt,
	}
}

// addEvent 添加领域事件
func (n *NPCAggregate) addEvent(event DomainEvent) {
	n.events = append(n.events, event)
}

// ToMap 转换为映射
func (n *NPCAggregate) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":          n.id,
		"name":        n.name,
		"description": n.description,
		"type":        n.npcType.String(),
		"status":      n.status.String(),
		"location":    n.location.ToMap(),
		"attributes":  n.attributes.ToMap(),
		"dialogues":   len(n.dialogues),
		"quests":      len(n.quests),
		"relationships": len(n.relationships),
		"has_shop":    n.shop != nil,
		"created_at":  n.createdAt,
		"updated_at":  n.updatedAt,
		"version":     n.version,
	}
}

// InteractionOption 交互选项
type InteractionOption struct {
	Type        InteractionType
	ID          string
	Name        string
	Description string
	Icon        string
	Enabled     bool
}

// InteractionType 交互类型
type InteractionType int

const (
	InteractionTypeDialogue InteractionType = iota + 1 // 对话
	InteractionTypeQuest                                // 任务
	InteractionTypeShop                                 // 商店
	InteractionTypeTrade                                // 交易
	InteractionTypeService                              // 服务
)

// String 返回交互类型字符串
func (it InteractionType) String() string {
	switch it {
	case InteractionTypeDialogue:
		return "dialogue"
	case InteractionTypeQuest:
		return "quest"
	case InteractionTypeShop:
		return "shop"
	case InteractionTypeTrade:
		return "trade"
	case InteractionTypeService:
		return "service"
	default:
		return "unknown"
	}
}

// NPCStatistics NPC统计信息
type NPCStatistics struct {
	NPCID               string
	Name                string
	Type                NPCType
	Status              NPCStatus
	TotalDialogues      int
	TotalQuests         int
	TotalRelationships  int
	AverageRelationship float64
	CreatedAt           time.Time
	LastActiveAt        time.Time
}

// 相关错误定义
var (
	ErrInvalidNPCName           = fmt.Errorf("invalid NPC name")
	ErrInvalidNPCStatus         = fmt.Errorf("invalid NPC status")
	ErrInvalidLocation          = fmt.Errorf("invalid location")
	ErrNPCCannotMove            = fmt.Errorf("NPC cannot move")
	ErrInvalidDialogue          = fmt.Errorf("invalid dialogue")
	ErrDialogueAlreadyExists    = fmt.Errorf("dialogue already exists")
	ErrDialogueNotFound         = fmt.Errorf("dialogue not found")
	ErrDialogueConditionsNotMet = fmt.Errorf("dialogue conditions not met")
	ErrInvalidQuest             = fmt.Errorf("invalid quest")
	ErrQuestAlreadyExists       = fmt.Errorf("quest already exists")
	ErrQuestNotFound            = fmt.Errorf("quest not found")
	ErrQuestConditionsNotMet    = fmt.Errorf("quest conditions not met")
	ErrInvalidShop              = fmt.Errorf("invalid shop")
	ErrNPCCannotHaveShop        = fmt.Errorf("NPC cannot have shop")
	ErrNPCHasNoShop             = fmt.Errorf("NPC has no shop")
)