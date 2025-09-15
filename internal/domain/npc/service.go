package npc

import (
	"fmt"
	"math/rand"
	"time"
)

// NPCService NPC领域服务
type NPCService struct {
	npcTemplates      map[NPCType]*NPCTemplate
	dialogueTemplates map[DialogueType]*DialogueTemplate
	questTemplates    map[QuestType]*QuestTemplate
	behaviorRules     map[BehaviorType]*BehaviorRule
	relationshipRules *RelationshipRules
	aiEngine          *AIEngine
}

// NewNPCService 创建NPC服务
func NewNPCService() *NPCService {
	service := &NPCService{
		npcTemplates:      make(map[NPCType]*NPCTemplate),
		dialogueTemplates: make(map[DialogueType]*DialogueTemplate),
		questTemplates:    make(map[QuestType]*QuestTemplate),
		behaviorRules:     make(map[BehaviorType]*BehaviorRule),
		relationshipRules: NewRelationshipRules(),
		aiEngine:          NewAIEngine(),
	}
	
	// 初始化默认模板和规则
	service.initializeDefaultTemplates()
	service.initializeBehaviorRules()
	
	return service
}

// CreateNPC 创建NPC
func (s *NPCService) CreateNPC(id, name, description string, npcType NPCType, location *Location) (*NPCAggregate, error) {
	if id == "" || name == "" {
		return nil, fmt.Errorf("invalid parameters for NPC creation")
	}
	
	npc := NewNPCAggregate(id, name, description, npcType)
	
	// 设置位置
	if location != nil {
		npc.MoveTo(location)
	}
	
	// 应用模板
	if template, exists := s.npcTemplates[npcType]; exists {
		s.applyTemplate(npc, template)
	}
	
	// 生成默认对话
	defaultDialogues := s.generateDefaultDialogues(npcType)
	for _, dialogue := range defaultDialogues {
		npc.AddDialogue(dialogue)
	}
	
	// 生成默认任务（如果适用）
	if npcType.CanGiveQuests() {
		defaultQuests := s.generateDefaultQuests(npcType)
		for _, quest := range defaultQuests {
			npc.AddQuest(quest)
		}
	}
	
	// 创建商店（如果适用）
	if npcType.CanHaveShop() {
		shop := s.createDefaultShop(npcType, id)
		npc.SetShop(shop)
	}
	
	return npc, nil
}

// GenerateDialogue 生成对话
func (s *NPCService) GenerateDialogue(dialogueType DialogueType, npcType NPCType, context map[string]interface{}) (*Dialogue, error) {
	template, exists := s.dialogueTemplates[dialogueType]
	if !exists {
		return nil, fmt.Errorf("dialogue template not found for type: %s", dialogueType.String())
	}
	
	// 生成唯一ID
	id := fmt.Sprintf("dialogue_%s_%s_%d", dialogueType.String(), npcType.String(), time.Now().UnixNano())
	
	// 根据模板创建对话
	dialogue := NewDialogue(
		id,
		template.GenerateName(npcType, context),
		template.GenerateDescription(npcType, context),
		dialogueType,
	)
	
	// 生成对话节点
	nodes := template.GenerateNodes(npcType, context)
	for _, node := range nodes {
		dialogue.AddNode(node)
	}
	
	// 设置开始节点
	if len(nodes) > 0 {
		dialogue.SetStartNode(nodes[0].GetID())
	}
	
	// 添加条件
	conditions := template.GenerateConditions(npcType, context)
	for _, condition := range conditions {
		dialogue.AddCondition(condition)
	}
	
	// 设置奖励
	if reward := template.GenerateReward(npcType, context); reward != nil {
		dialogue.SetReward(reward)
	}
	
	return dialogue, nil
}

// GenerateQuest 生成任务
func (s *NPCService) GenerateQuest(questType QuestType, npcType NPCType, playerLevel int) (*Quest, error) {
	template, exists := s.questTemplates[questType]
	if !exists {
		return nil, fmt.Errorf("quest template not found for type: %s", questType.String())
	}
	
	// 生成唯一ID
	id := fmt.Sprintf("quest_%s_%s_%d", questType.String(), npcType.String(), time.Now().UnixNano())
	
	// 根据模板创建任务
	quest := NewQuest(
		id,
		template.GenerateName(npcType, playerLevel),
		template.GenerateDescription(npcType, playerLevel),
		questType,
	)
	
	// 生成目标
	objectives := template.GenerateObjectives(npcType, playerLevel)
	for _, objective := range objectives {
		quest.AddObjective(objective)
	}
	
	// 设置奖励
	reward := template.GenerateReward(npcType, playerLevel)
	quest.SetReward(reward)
	
	// 添加前置条件
	prerequisites := template.GeneratePrerequisites(npcType, playerLevel)
	for _, prerequisite := range prerequisites {
		quest.AddPrerequisite(prerequisite)
	}
	
	// 设置时间限制
	if timeLimit := template.GetTimeLimit(playerLevel); timeLimit > 0 {
		quest.SetTimeLimit(timeLimit)
	}
	
	// 设置重复性
	quest.SetRepeatable(template.IsRepeatable())
	quest.SetDailyReset(template.IsDailyReset())
	
	return quest, nil
}

// ProcessDialogue 处理对话
func (s *NPCService) ProcessDialogue(npc *NPCAggregate, playerID string, dialogueID string, optionID string) (*DialogueResponse, error) {
	dialogue, err := npc.GetDialogue(dialogueID)
	if err != nil {
		return nil, err
	}
	
	// 检查是否可以开始对话
	if !dialogue.CanStart(playerID) {
		return nil, fmt.Errorf("cannot start dialogue")
	}
	
	// 开始对话会话
	session, err := npc.StartDialogue(dialogueID, playerID)
	if err != nil {
		return nil, err
	}
	
	// 获取当前节点
	currentNode := dialogue.GetStartNode()
	if session.GetCurrentNode() != "" {
		currentNode = dialogue.GetNode(session.GetCurrentNode())
	}
	
	if currentNode == nil {
		return nil, fmt.Errorf("dialogue node not found")
	}
	
	// 处理选项
	if optionID != "" {
		for _, option := range currentNode.GetOptions() {
			if option.GetID() == optionID && option.IsAvailable(playerID) {
				// 执行选项动作
				if err := option.Execute(playerID); err != nil {
					return nil, err
				}
				
				// 移动到目标节点
				if option.GetTargetNode() != "" {
					currentNode = dialogue.GetNode(option.GetTargetNode())
					session.SetCurrentNode(option.GetTargetNode())
				}
				break
			}
		}
	}
	
	// 执行节点动作
	if err := currentNode.ExecuteActions(playerID); err != nil {
		return nil, err
	}
	
	// 使用对话
	dialogue.Use(playerID)
	
	// 创建响应
	response := &DialogueResponse{
		NPCID:       npc.GetID(),
		DialogueID:  dialogueID,
		NodeID:      currentNode.GetID(),
		Text:        currentNode.GetText(),
		Speaker:     currentNode.GetSpeaker(),
		Options:     s.filterAvailableOptions(currentNode.GetOptions(), playerID),
		CanContinue: currentNode.GetNextNode() != "",
		Completed:   currentNode.GetNextNode() == "",
	}
	
	return response, nil
}

// UpdateNPCBehavior 更新NPC行为
func (s *NPCService) UpdateNPCBehavior(npc *NPCAggregate, deltaTime time.Duration) error {
	behavior := npc.GetBehavior()
	behaviorRule, exists := s.behaviorRules[behavior.Type]
	if !exists {
		return fmt.Errorf("behavior rule not found for type: %s", behavior.Type.String())
	}
	
	// 应用行为规则
	behaviorRule.Apply(npc, deltaTime)
	
	// 更新NPC
	npc.Update(deltaTime)
	
	return nil
}

// CalculateRelationshipChange 计算关系变化
func (s *NPCService) CalculateRelationshipChange(npc *NPCAggregate, playerID string, action string, context map[string]interface{}) int {
	return s.relationshipRules.CalculateChange(npc.GetType(), action, context)
}

// GenerateAIResponse 生成AI响应
func (s *NPCService) GenerateAIResponse(npc *NPCAggregate, playerID string, input string) (string, error) {
	return s.aiEngine.GenerateResponse(npc, playerID, input)
}

// ValidateQuestCompletion 验证任务完成
func (s *NPCService) ValidateQuestCompletion(quest *Quest, questInstance *QuestInstance, playerData map[string]interface{}) bool {
	objectives := quest.GetObjectives()
	for _, objective := range objectives {
		if !objective.IsOptional() {
			required := objective.GetRequired()
			current := questInstance.GetProgress(objective.GetID())
			if current < required {
				return false
			}
		}
	}
	return true
}

// GetRecommendedQuests 获取推荐任务
func (s *NPCService) GetRecommendedQuests(npc *NPCAggregate, playerID string, playerLevel int) []*Quest {
	availableQuests := npc.GetAvailableQuests(playerID)
	var recommended []*Quest
	
	for _, quest := range availableQuests {
		// 根据玩家等级和任务类型推荐
		if s.isQuestRecommended(quest, playerLevel) {
			recommended = append(recommended, quest)
		}
	}
	
	return recommended
}

// GetOptimalDialogue 获取最佳对话
func (s *NPCService) GetOptimalDialogue(npc *NPCAggregate, playerID string, context map[string]interface{}) *Dialogue {
	availableDialogues := npc.GetAvailableDialogues(playerID)
	
	// 根据上下文选择最合适的对话
	for _, dialogue := range availableDialogues {
		if s.isDialogueOptimal(dialogue, context) {
			return dialogue
		}
	}
	
	// 返回默认对话
	if len(availableDialogues) > 0 {
		return availableDialogues[0]
	}
	
	return nil
}

// 私有方法

// applyTemplate 应用模板
func (s *NPCService) applyTemplate(npc *NPCAggregate, template *NPCTemplate) {
	// 设置属性
	attributes := npc.GetAttributes()
	attributes.SetLevel(template.Level)
	attributes.Strength = template.Attributes.Strength
	attributes.Agility = template.Attributes.Agility
	attributes.Intelligence = template.Attributes.Intelligence
	attributes.Charisma = template.Attributes.Charisma
	attributes.Luck = template.Attributes.Luck
	attributes.MoveSpeed = template.Attributes.MoveSpeed
	attributes.ViewRange = template.Attributes.ViewRange
	attributes.HearRange = template.Attributes.HearRange
	
	// 设置行为
	behavior := npc.GetBehavior()
	behavior.SetBehaviorType(template.DefaultBehavior)
	behavior.MoveSpeed = template.Attributes.MoveSpeed
	behavior.CanMove = template.CanMove
	behavior.CanTalk = template.CanTalk
	behavior.CanFight = template.CanFight
	
	// 设置巡逻点
	for _, point := range template.PatrolPoints {
		behavior.AddPatrolPoint(point)
	}
}

// generateDefaultDialogues 生成默认对话
func (s *NPCService) generateDefaultDialogues(npcType NPCType) []*Dialogue {
	var dialogues []*Dialogue
	
	// 生成问候对话
	if greeting, err := s.GenerateDialogue(DialogueTypeGreeting, npcType, nil); err == nil {
		dialogues = append(dialogues, greeting)
	}
	
	// 根据NPC类型生成特定对话
	switch npcType {
	case NPCTypeMerchant:
		if trade, err := s.GenerateDialogue(DialogueTypeTrade, npcType, nil); err == nil {
			dialogues = append(dialogues, trade)
		}
	case NPCTypeQuestGiver:
		if quest, err := s.GenerateDialogue(DialogueTypeQuest, npcType, nil); err == nil {
			dialogues = append(dialogues, quest)
		}
	case NPCTypeVillager:
		if info, err := s.GenerateDialogue(DialogueTypeInformation, npcType, nil); err == nil {
			dialogues = append(dialogues, info)
		}
		if rumor, err := s.GenerateDialogue(DialogueTypeRumor, npcType, nil); err == nil {
			dialogues = append(dialogues, rumor)
		}
	}
	
	return dialogues
}

// generateDefaultQuests 生成默认任务
func (s *NPCService) generateDefaultQuests(npcType NPCType) []*Quest {
	var quests []*Quest
	
	// 根据NPC类型生成不同的任务
	switch npcType {
	case NPCTypeQuestGiver:
		// 生成各种类型的任务
		questTypes := []QuestType{QuestTypeKill, QuestTypeCollect, QuestTypeDeliver}
		for _, questType := range questTypes {
			if quest, err := s.GenerateQuest(questType, npcType, 1); err == nil {
				quests = append(quests, quest)
			}
		}
	case NPCTypeVillager:
		// 生成简单任务
		if quest, err := s.GenerateQuest(QuestTypeTalk, npcType, 1); err == nil {
			quests = append(quests, quest)
		}
	case NPCTypeGuard:
		// 生成巡逻或保护任务
		if quest, err := s.GenerateQuest(QuestTypeEscort, npcType, 1); err == nil {
			quests = append(quests, quest)
		}
	}
	
	return quests
}

// createDefaultShop 创建默认商店
func (s *NPCService) createDefaultShop(npcType NPCType, npcID string) *Shop {
	shopID := fmt.Sprintf("shop_%s_%s", npcType.String(), npcID)
	shop := NewShop(shopID, fmt.Sprintf("%s商店", npcType.String()), "默认商店")
	
	// 根据NPC类型添加商品
	switch npcType {
	case NPCTypeMerchant:
		// 添加一般商品
		shop.AddItem(NewShopItem("item_potion_health", "生命药水", "恢复生命值", 50, 10))
		shop.AddItem(NewShopItem("item_potion_mana", "法力药水", "恢复法力值", 30, 10))
	case NPCTypeBlacksmith:
		// 添加武器装备
		shop.AddItem(NewShopItem("weapon_sword", "铁剑", "基础武器", 100, 5))
		shop.AddItem(NewShopItem("armor_leather", "皮甲", "基础护甲", 80, 5))
	case NPCTypeInnkeeper:
		// 添加食物和住宿
		shop.AddItem(NewShopItem("food_bread", "面包", "基础食物", 10, 20))
		shop.AddItem(NewShopItem("service_room", "房间", "住宿服务", 25, 10))
	}
	
	return shop
}

// filterAvailableOptions 过滤可用选项
func (s *NPCService) filterAvailableOptions(options []*DialogueOption, playerID string) []*DialogueOption {
	var available []*DialogueOption
	for _, option := range options {
		if option.IsAvailable(playerID) {
			available = append(available, option)
		}
	}
	return available
}

// isQuestRecommended 检查任务是否推荐
func (s *NPCService) isQuestRecommended(quest *Quest, playerLevel int) bool {
	// 简单的推荐逻辑：任务类型和玩家等级匹配
	switch quest.GetType() {
	case QuestTypeDaily:
		return true // 日常任务总是推荐
	case QuestTypeKill, QuestTypeCollect:
		return playerLevel >= 5 // 需要一定等级
	case QuestTypeDeliver, QuestTypeTalk:
		return playerLevel >= 1 // 低等级任务
	default:
		return playerLevel >= 10 // 高等级任务
	}
}

// isDialogueOptimal 检查对话是否最佳
func (s *NPCService) isDialogueOptimal(dialogue *Dialogue, context map[string]interface{}) bool {
	// 根据上下文判断对话是否合适
	if mood, exists := context["mood"]; exists {
		if mood == "friendly" && dialogue.GetType() == DialogueTypeGreeting {
			return true
		}
		if mood == "business" && dialogue.GetType() == DialogueTypeTrade {
			return true
		}
	}
	
	// 默认返回false，使用第一个可用对话
	return false
}

// initializeDefaultTemplates 初始化默认模板
func (s *NPCService) initializeDefaultTemplates() {
	// 初始化NPC模板
	s.npcTemplates[NPCTypeVillager] = &NPCTemplate{
		Level:           1,
		Attributes:      NewNPCAttributes(),
		DefaultBehavior: BehaviorTypeWander,
		CanMove:         true,
		CanTalk:         true,
		CanFight:        false,
		PatrolPoints:    make([]*Location, 0),
	}
	
	s.npcTemplates[NPCTypeMerchant] = &NPCTemplate{
		Level:           5,
		Attributes:      NewNPCAttributes(),
		DefaultBehavior: BehaviorTypeStationary,
		CanMove:         false,
		CanTalk:         true,
		CanFight:        false,
		PatrolPoints:    make([]*Location, 0),
	}
	
	s.npcTemplates[NPCTypeGuard] = &NPCTemplate{
		Level:           10,
		Attributes:      NewNPCAttributes(),
		DefaultBehavior: BehaviorTypePatrol,
		CanMove:         true,
		CanTalk:         true,
		CanFight:        true,
		PatrolPoints:    make([]*Location, 0),
	}
	
	// 初始化对话模板
	s.dialogueTemplates[DialogueTypeGreeting] = &DialogueTemplate{
		Name:        "问候",
		Description: "基础问候对话",
		Nodes:       make([]*DialogueNodeTemplate, 0),
	}
	
	s.dialogueTemplates[DialogueTypeTrade] = &DialogueTemplate{
		Name:        "交易",
		Description: "商店交易对话",
		Nodes:       make([]*DialogueNodeTemplate, 0),
	}
	
	// 初始化任务模板
	s.questTemplates[QuestTypeKill] = &QuestTemplate{
		Name:         "击杀任务",
		Description:  "击杀指定目标",
		Objectives:   make([]*QuestObjectiveTemplate, 0),
		BaseReward:   NewQuestReward(),
		TimeLimit:    time.Hour * 24,
		Repeatable:   false,
		DailyReset:   false,
	}
	
	s.questTemplates[QuestTypeCollect] = &QuestTemplate{
		Name:         "收集任务",
		Description:  "收集指定物品",
		Objectives:   make([]*QuestObjectiveTemplate, 0),
		BaseReward:   NewQuestReward(),
		TimeLimit:    time.Hour * 12,
		Repeatable:   true,
		DailyReset:   true,
	}
}

// initializeBehaviorRules 初始化行为规则
func (s *NPCService) initializeBehaviorRules() {
	s.behaviorRules[BehaviorTypeIdle] = &BehaviorRule{
		Type:        BehaviorTypeIdle,
		Description: "空闲行为",
		ApplyFunc: func(npc *NPCAggregate, deltaTime time.Duration) {
			// 空闲状态不需要特殊处理
		},
	}
	
	s.behaviorRules[BehaviorTypePatrol] = &BehaviorRule{
		Type:        BehaviorTypePatrol,
		Description: "巡逻行为",
		ApplyFunc: func(npc *NPCAggregate, deltaTime time.Duration) {
			behavior := npc.GetBehavior()
			behavior.Update(deltaTime)
		},
	}
	
	s.behaviorRules[BehaviorTypeWander] = &BehaviorRule{
		Type:        BehaviorTypeWander,
		Description: "漫游行为",
		ApplyFunc: func(npc *NPCAggregate, deltaTime time.Duration) {
			behavior := npc.GetBehavior()
			behavior.Update(deltaTime)
			
			// 随机移动逻辑
			if rand.Float64() < 0.1 { // 10%概率改变方向
				currentLocation := npc.GetLocation()
				newX := currentLocation.X + (rand.Float64()-0.5)*10
				newY := currentLocation.Y + (rand.Float64()-0.5)*10
				newLocation := NewLocation(newX, newY, currentLocation.Z, currentLocation.Region, currentLocation.Zone)
				npc.MoveTo(newLocation)
			}
		},
	}
	
	s.behaviorRules[BehaviorTypeStationary] = &BehaviorRule{
		Type:        BehaviorTypeStationary,
		Description: "固定行为",
		ApplyFunc: func(npc *NPCAggregate, deltaTime time.Duration) {
			// 固定位置，不移动
		},
	}
}

// 辅助结构体

// NPCTemplate NPC模板
type NPCTemplate struct {
	Level           int
	Attributes      *NPCAttributes
	DefaultBehavior BehaviorType
	CanMove         bool
	CanTalk         bool
	CanFight        bool
	PatrolPoints    []*Location
}

// DialogueTemplate 对话模板
type DialogueTemplate struct {
	Name        string
	Description string
	Nodes       []*DialogueNodeTemplate
}

// GenerateName 生成名称
func (dt *DialogueTemplate) GenerateName(npcType NPCType, context map[string]interface{}) string {
	return fmt.Sprintf("%s - %s", dt.Name, npcType.String())
}

// GenerateDescription 生成描述
func (dt *DialogueTemplate) GenerateDescription(npcType NPCType, context map[string]interface{}) string {
	return fmt.Sprintf("%s (%s)", dt.Description, npcType.String())
}

// GenerateNodes 生成节点
func (dt *DialogueTemplate) GenerateNodes(npcType NPCType, context map[string]interface{}) []*DialogueNode {
	var nodes []*DialogueNode
	
	// 根据模板生成节点
	for i, nodeTemplate := range dt.Nodes {
		nodeID := fmt.Sprintf("node_%d", i)
		node := NewDialogueNode(nodeID, nodeTemplate.Text, nodeTemplate.Speaker)
		nodes = append(nodes, node)
	}
	
	// 如果没有模板节点，生成默认节点
	if len(nodes) == 0 {
		defaultNode := NewDialogueNode("node_0", s.getDefaultDialogueText(npcType), "NPC")
		nodes = append(nodes, defaultNode)
	}
	
	return nodes
}

// GenerateConditions 生成条件
func (dt *DialogueTemplate) GenerateConditions(npcType NPCType, context map[string]interface{}) []*DialogueCondition {
	// 简化实现，返回空条件
	return make([]*DialogueCondition, 0)
}

// GenerateReward 生成奖励
func (dt *DialogueTemplate) GenerateReward(npcType NPCType, context map[string]interface{}) *DialogueReward {
	// 简化实现，返回nil
	return nil
}

// getDefaultDialogueText 获取默认对话文本
func (s *NPCService) getDefaultDialogueText(npcType NPCType) string {
	switch npcType {
	case NPCTypeVillager:
		return "你好，旅行者！欢迎来到我们的村庄。"
	case NPCTypeMerchant:
		return "欢迎光临！我这里有各种商品，看看有什么需要的吗？"
	case NPCTypeGuard:
		return "站住！请出示你的通行证。"
	case NPCTypeQuestGiver:
		return "勇敢的冒险者，我这里有个任务需要你的帮助。"
	default:
		return "你好。"
	}
}

// DialogueNodeTemplate 对话节点模板
type DialogueNodeTemplate struct {
	Text    string
	Speaker string
	Options []*DialogueOptionTemplate
}

// DialogueOptionTemplate 对话选项模板
type DialogueOptionTemplate struct {
	Text       string
	TargetNode string
	Conditions []*DialogueCondition
	Actions    []*DialogueAction
}

// QuestTemplate 任务模板
type QuestTemplate struct {
	Name         string
	Description  string
	Objectives   []*QuestObjectiveTemplate
	BaseReward   *QuestReward
	TimeLimit    time.Duration
	Repeatable   bool
	DailyReset   bool
}

// GenerateName 生成名称
func (qt *QuestTemplate) GenerateName(npcType NPCType, playerLevel int) string {
	return fmt.Sprintf("%s (Lv.%d)", qt.Name, playerLevel)
}

// GenerateDescription 生成描述
func (qt *QuestTemplate) GenerateDescription(npcType NPCType, playerLevel int) string {
	return fmt.Sprintf("%s - 适合等级 %d 的玩家", qt.Description, playerLevel)
}

// GenerateObjectives 生成目标
func (qt *QuestTemplate) GenerateObjectives(npcType NPCType, playerLevel int) []*QuestObjective {
	var objectives []*QuestObjective
	
	for i, objTemplate := range qt.Objectives {
		objID := fmt.Sprintf("obj_%d", i)
		objective := NewQuestObjective(
			objID,
			objTemplate.Description,
			objTemplate.Type,
			objTemplate.Target,
			objTemplate.Required,
		)
		objectives = append(objectives, objective)
	}
	
	return objectives
}

// GenerateReward 生成奖励
func (qt *QuestTemplate) GenerateReward(npcType NPCType, playerLevel int) *QuestReward {
	reward := NewQuestReward()
	
	// 根据玩家等级调整奖励
	multiplier := float64(playerLevel)
	reward.AddGold(int(float64(qt.BaseReward.gold) * multiplier))
	reward.AddExperience(int(float64(qt.BaseReward.experience) * multiplier))
	
	// 复制物品奖励
	for itemID, quantity := range qt.BaseReward.items {
		reward.AddItem(itemID, quantity)
	}
	
	return reward
}

// GeneratePrerequisites 生成前置条件
func (qt *QuestTemplate) GeneratePrerequisites(npcType NPCType, playerLevel int) []*QuestPrerequisite {
	// 简化实现，返回空前置条件
	return make([]*QuestPrerequisite, 0)
}

// GetTimeLimit 获取时间限制
func (qt *QuestTemplate) GetTimeLimit(playerLevel int) time.Duration {
	return qt.TimeLimit
}

// IsRepeatable 检查是否可重复
func (qt *QuestTemplate) IsRepeatable() bool {
	return qt.Repeatable
}

// IsDailyReset 检查是否每日重置
func (qt *QuestTemplate) IsDailyReset() bool {
	return qt.DailyReset
}

// QuestObjectiveTemplate 任务目标模板
type QuestObjectiveTemplate struct {
	Description string
	Type        ObjectiveType
	Target      string
	Required    int
}

// BehaviorRule 行为规则
type BehaviorRule struct {
	Type        BehaviorType
	Description string
	ApplyFunc   func(*NPCAggregate, time.Duration)
}

// Apply 应用规则
func (br *BehaviorRule) Apply(npc *NPCAggregate, deltaTime time.Duration) {
	if br.ApplyFunc != nil {
		br.ApplyFunc(npc, deltaTime)
	}
}

// RelationshipRules 关系规则
type RelationshipRules struct {
	rules map[string]*RelationshipRule
}

// NewRelationshipRules 创建关系规则
func NewRelationshipRules() *RelationshipRules {
	rules := &RelationshipRules{
		rules: make(map[string]*RelationshipRule),
	}
	
	// 初始化默认规则
	rules.rules["quest_complete"] = &RelationshipRule{Action: "quest_complete", BaseChange: 10}
	rules.rules["quest_fail"] = &RelationshipRule{Action: "quest_fail", BaseChange: -5}
	rules.rules["trade"] = &RelationshipRule{Action: "trade", BaseChange: 1}
	rules.rules["dialogue"] = &RelationshipRule{Action: "dialogue", BaseChange: 1}
	
	return rules
}

// CalculateChange 计算关系变化
func (rr *RelationshipRules) CalculateChange(npcType NPCType, action string, context map[string]interface{}) int {
	rule, exists := rr.rules[action]
	if !exists {
		return 0
	}
	
	change := rule.BaseChange
	
	// 根据NPC类型调整
	switch npcType {
	case NPCTypeMerchant:
		if action == "trade" {
			change *= 2 // 商人更重视交易
		}
	case NPCTypeQuestGiver:
		if action == "quest_complete" {
			change *= 2 // 任务发布者更重视任务完成
		}
	}
	
	return change
}

// RelationshipRule 关系规则
type RelationshipRule struct {
	Action     string
	BaseChange int
}

// AIEngine AI引擎
type AIEngine struct {
	responseTemplates map[string][]string
}

// NewAIEngine 创建AI引擎
func NewAIEngine() *AIEngine {
	engine := &AIEngine{
		responseTemplates: make(map[string][]string),
	}
	
	// 初始化响应模板
	engine.responseTemplates["greeting"] = []string{
		"你好！",
		"欢迎！",
		"很高兴见到你！",
	}
	
	engine.responseTemplates["farewell"] = []string{
		"再见！",
		"祝你好运！",
		"期待下次见面！",
	}
	
	return engine
}

// GenerateResponse 生成响应
func (ai *AIEngine) GenerateResponse(npc *NPCAggregate, playerID string, input string) (string, error) {
	// 简化的AI响应逻辑
	if templates, exists := ai.responseTemplates["greeting"]; exists {
		index := rand.Intn(len(templates))
		return templates[index], nil
	}
	
	return "我不明白你在说什么。", nil
}

// DialogueResponse 对话响应
type DialogueResponse struct {
	NPCID       string
	DialogueID  string
	NodeID      string
	Text        string
	Speaker     string
	Options     []*DialogueOption
	CanContinue bool
	Completed   bool
}