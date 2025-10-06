package npc

import (
	"fmt"
	"time"
)

// Dialogue 对话实体
type Dialogue struct {
	id          string
	name        string
	description string
	type_       DialogueType
	nodes       map[string]*DialogueNode
	startNodeID string
	conditions  []*DialogueCondition
	rewards     *DialogueReward
	cooldown    time.Duration
	lastUsed    map[string]time.Time
	maxUses     int
	useCount    map[string]int
	createdAt   time.Time
	updatedAt   time.Time
}

// NewDialogue 创建对话
func NewDialogue(id, name, description string, dialogueType DialogueType) *Dialogue {
	now := time.Now()
	return &Dialogue{
		id:          id,
		name:        name,
		description: description,
		type_:       dialogueType,
		nodes:       make(map[string]*DialogueNode),
		conditions:  make([]*DialogueCondition, 0),
		lastUsed:    make(map[string]time.Time),
		useCount:    make(map[string]int),
		maxUses:     -1, // 无限制
		createdAt:   now,
		updatedAt:   now,
	}
}

// GetID 获取ID
func (d *Dialogue) GetID() string {
	return d.id
}

// GetName 获取名称
func (d *Dialogue) GetName() string {
	return d.name
}

// GetDescription 获取描述
func (d *Dialogue) GetDescription() string {
	return d.description
}

// GetType 获取类型
func (d *Dialogue) GetType() DialogueType {
	return d.type_
}

// AddNode 添加对话节点
func (d *Dialogue) AddNode(node *DialogueNode) {
	d.nodes[node.GetID()] = node
	d.updatedAt = time.Now()
}

// SetStartNode 设置开始节点
func (d *Dialogue) SetStartNode(nodeID string) error {
	if _, exists := d.nodes[nodeID]; !exists {
		return fmt.Errorf("node not found: %s", nodeID)
	}
	d.startNodeID = nodeID
	d.updatedAt = time.Now()
	return nil
}

// GetStartNode 获取开始节点
func (d *Dialogue) GetStartNode() *DialogueNode {
	return d.nodes[d.startNodeID]
}

// GetNode 获取节点
func (d *Dialogue) GetNode(nodeID string) *DialogueNode {
	return d.nodes[nodeID]
}

// AddCondition 添加条件
func (d *Dialogue) AddCondition(condition *DialogueCondition) {
	d.conditions = append(d.conditions, condition)
	d.updatedAt = time.Now()
}

// CanStart 检查是否可以开始
func (d *Dialogue) CanStart(playerID string) bool {
	// 检查使用次数
	if d.maxUses > 0 && d.useCount[playerID] >= d.maxUses {
		return false
	}

	// 检查冷却时间
	if d.cooldown > 0 {
		if lastUsed, exists := d.lastUsed[playerID]; exists {
			if time.Since(lastUsed) < d.cooldown {
				return false
			}
		}
	}

	// 检查条件
	for _, condition := range d.conditions {
		if !condition.Check(playerID) {
			return false
		}
	}

	return true
}

// Use 使用对话
func (d *Dialogue) Use(playerID string) {
	d.lastUsed[playerID] = time.Now()
	d.useCount[playerID]++
	d.updatedAt = time.Now()
}

// SetReward 设置奖励
func (d *Dialogue) SetReward(reward *DialogueReward) {
	d.rewards = reward
	d.updatedAt = time.Now()
}

// GetReward 获取奖励
func (d *Dialogue) GetReward() *DialogueReward {
	return d.rewards
}

// DialogueNode 对话节点
type DialogueNode struct {
	id       string
	text     string
	speaker  string
	options  []*DialogueOption
	actions  []*DialogueAction
	nextNode string
}

// NewDialogueNode 创建对话节点
func NewDialogueNode(id, text, speaker string) *DialogueNode {
	return &DialogueNode{
		id:      id,
		text:    text,
		speaker: speaker,
		options: make([]*DialogueOption, 0),
		actions: make([]*DialogueAction, 0),
	}
}

// GetID 获取ID
func (dn *DialogueNode) GetID() string {
	return dn.id
}

// GetText 获取文本
func (dn *DialogueNode) GetText() string {
	return dn.text
}

// GetSpeaker 获取说话者
func (dn *DialogueNode) GetSpeaker() string {
	return dn.speaker
}

// GetOptions 获取选项
func (dn *DialogueNode) GetOptions() []*DialogueOption {
	return dn.options
}

// AddOption 添加选项
func (dn *DialogueNode) AddOption(option *DialogueOption) {
	dn.options = append(dn.options, option)
}

// AddAction 添加动作
func (dn *DialogueNode) AddAction(action *DialogueAction) {
	dn.actions = append(dn.actions, action)
}

// ExecuteActions 执行动作
func (dn *DialogueNode) ExecuteActions(playerID string) error {
	for _, action := range dn.actions {
		if err := action.Execute(playerID); err != nil {
			return err
		}
	}
	return nil
}

// SetNextNode 设置下一个节点
func (dn *DialogueNode) SetNextNode(nodeID string) {
	dn.nextNode = nodeID
}

// GetNextNode 获取下一个节点
func (dn *DialogueNode) GetNextNode() string {
	return dn.nextNode
}

// DialogueOption 对话选项
type DialogueOption struct {
	id         string
	text       string
	targetNode string
	conditions []*DialogueCondition
	actions    []*DialogueAction
}

// NewDialogueOption 创建对话选项
func NewDialogueOption(id, text, targetNode string) *DialogueOption {
	return &DialogueOption{
		id:         id,
		text:       text,
		targetNode: targetNode,
		conditions: make([]*DialogueCondition, 0),
		actions:    make([]*DialogueAction, 0),
	}
}

// GetID 获取ID
func (do *DialogueOption) GetID() string {
	return do.id
}

// GetText 获取文本
func (do *DialogueOption) GetText() string {
	return do.text
}

// GetTargetNode 获取目标节点
func (do *DialogueOption) GetTargetNode() string {
	return do.targetNode
}

// IsAvailable 检查是否可用
func (do *DialogueOption) IsAvailable(playerID string) bool {
	for _, condition := range do.conditions {
		if !condition.Check(playerID) {
			return false
		}
	}
	return true
}

// Execute 执行选项
func (do *DialogueOption) Execute(playerID string) error {
	for _, action := range do.actions {
		if err := action.Execute(playerID); err != nil {
			return err
		}
	}
	return nil
}

// DialogueCondition 对话条件
type DialogueCondition struct {
	type_    ConditionType
	key      string
	operator string
	value    interface{}
	message  string
}

// NewDialogueCondition 创建对话条件
func NewDialogueCondition(conditionType ConditionType, key, operator string, value interface{}, message string) *DialogueCondition {
	return &DialogueCondition{
		type_:    conditionType,
		key:      key,
		operator: operator,
		value:    value,
		message:  message,
	}
}

// Check 检查条件
func (dc *DialogueCondition) Check(playerID string) bool {
	// 这里应该根据条件类型和玩家数据进行检查
	// 简化实现，总是返回true
	return true
}

// DialogueAction 对话动作
type DialogueAction struct {
	type_      ActionType
	parameters map[string]interface{}
}

// NewDialogueAction 创建对话动作
func NewDialogueAction(actionType ActionType, parameters map[string]interface{}) *DialogueAction {
	return &DialogueAction{
		type_:      actionType,
		parameters: parameters,
	}
}

// Execute 执行动作
func (da *DialogueAction) Execute(playerID string) error {
	// 这里应该根据动作类型执行相应的操作
	// 简化实现，直接返回nil
	return nil
}

// DialogueReward 对话奖励
type DialogueReward struct {
	gold       int
	experience int
	items      map[string]int
	special    map[string]interface{}
}

// NewDialogueReward 创建对话奖励
func NewDialogueReward() *DialogueReward {
	return &DialogueReward{
		items:   make(map[string]int),
		special: make(map[string]interface{}),
	}
}

// AddGold 添加金币奖励
func (dr *DialogueReward) AddGold(amount int) {
	dr.gold += amount
}

// AddExperience 添加经验奖励
func (dr *DialogueReward) AddExperience(amount int) {
	dr.experience += amount
}

// AddItem 添加物品奖励
func (dr *DialogueReward) AddItem(itemID string, quantity int) {
	dr.items[itemID] = quantity
}

// DialogueSession 对话会话
type DialogueSession struct {
	npcID       string
	dialogueID  string
	playerID    string
	currentNode string
	startTime   time.Time
	lastUpdate  time.Time
	context     map[string]interface{}
	active      bool
}

// NewDialogueSession 创建对话会话
func NewDialogueSession(npcID, dialogueID, playerID string) *DialogueSession {
	now := time.Now()
	return &DialogueSession{
		npcID:      npcID,
		dialogueID: dialogueID,
		playerID:   playerID,
		startTime:  now,
		lastUpdate: now,
		context:    make(map[string]interface{}),
		active:     true,
	}
}

// GetNPCID 获取NPC ID
func (ds *DialogueSession) GetNPCID() string {
	return ds.npcID
}

// GetDialogueID 获取对话ID
func (ds *DialogueSession) GetDialogueID() string {
	return ds.dialogueID
}

// GetPlayerID 获取玩家ID
func (ds *DialogueSession) GetPlayerID() string {
	return ds.playerID
}

// GetCurrentNode 获取当前节点
func (ds *DialogueSession) GetCurrentNode() string {
	return ds.currentNode
}

// SetCurrentNode 设置当前节点
func (ds *DialogueSession) SetCurrentNode(nodeID string) {
	ds.currentNode = nodeID
	ds.lastUpdate = time.Now()
}

// IsActive 检查是否激活
func (ds *DialogueSession) IsActive() bool {
	return ds.active
}

// End 结束会话
func (ds *DialogueSession) End() {
	ds.active = false
	ds.lastUpdate = time.Now()
}

// GetDuration 获取持续时间
func (ds *DialogueSession) GetDuration() time.Duration {
	return ds.lastUpdate.Sub(ds.startTime)
}

// GetID 获取会话ID
func (ds *DialogueSession) GetID() string {
	return ds.npcID + "_" + ds.dialogueID + "_" + ds.playerID
}

// GetCurrentNodeID 获取当前节点ID
func (ds *DialogueSession) GetCurrentNodeID() string {
	return ds.currentNode
}

// GetStartTime 获取开始时间
func (ds *DialogueSession) GetStartTime() time.Time {
	return ds.startTime
}

// GetEndTime 获取结束时间
func (ds *DialogueSession) GetEndTime() time.Time {
	return ds.lastUpdate
}

// GetContext 获取上下文
func (ds *DialogueSession) GetContext() map[string]interface{} {
	return ds.context
}

// Quest 任务实体
type Quest struct {
	id            string
	name          string
	description   string
	type_         QuestType
	objectives    []*QuestObjective
	rewards       *QuestReward
	prerequisites []*QuestPrerequisite
	timeLimit     time.Duration
	repeatable    bool
	dailyReset    bool
	createdAt     time.Time
	updatedAt     time.Time
}

// NewQuest 创建任务
func NewQuest(id, name, description string, questType QuestType) *Quest {
	now := time.Now()
	return &Quest{
		id:            id,
		name:          name,
		description:   description,
		type_:         questType,
		objectives:    make([]*QuestObjective, 0),
		prerequisites: make([]*QuestPrerequisite, 0),
		createdAt:     now,
		updatedAt:     now,
	}
}

// GetID 获取ID
func (q *Quest) GetID() string {
	return q.id
}

// GetName 获取名称
func (q *Quest) GetName() string {
	return q.name
}

// GetDescription 获取描述
func (q *Quest) GetDescription() string {
	return q.description
}

// GetType 获取类型
func (q *Quest) GetType() QuestType {
	return q.type_
}

// AddObjective 添加目标
func (q *Quest) AddObjective(objective *QuestObjective) {
	q.objectives = append(q.objectives, objective)
	q.updatedAt = time.Now()
}

// GetObjectives 获取目标
func (q *Quest) GetObjectives() []*QuestObjective {
	return q.objectives
}

// SetReward 设置奖励
func (q *Quest) SetReward(reward *QuestReward) {
	q.rewards = reward
	q.updatedAt = time.Now()
}

// GetReward 获取奖励
func (q *Quest) GetReward() *QuestReward {
	return q.rewards
}

// AddPrerequisite 添加前置条件
func (q *Quest) AddPrerequisite(prerequisite *QuestPrerequisite) {
	q.prerequisites = append(q.prerequisites, prerequisite)
	q.updatedAt = time.Now()
}

// CanAccept 检查是否可以接受
func (q *Quest) CanAccept(playerID string) bool {
	// 检查前置条件
	for _, prerequisite := range q.prerequisites {
		if !prerequisite.Check(playerID) {
			return false
		}
	}
	return true
}

// SetTimeLimit 设置时间限制
func (q *Quest) SetTimeLimit(duration time.Duration) {
	q.timeLimit = duration
	q.updatedAt = time.Now()
}

// SetRepeatable 设置是否可重复
func (q *Quest) SetRepeatable(repeatable bool) {
	q.repeatable = repeatable
	q.updatedAt = time.Now()
}

// SetDailyReset 设置每日重置
func (q *Quest) SetDailyReset(dailyReset bool) {
	q.dailyReset = dailyReset
	q.updatedAt = time.Now()
}

// QuestObjective 任务目标
type QuestObjective struct {
	id          string
	description string
	type_       ObjectiveType
	target      string
	required    int
	optional    bool
}

// NewQuestObjective 创建任务目标
func NewQuestObjective(id, description string, objectiveType ObjectiveType, target string, required int) *QuestObjective {
	return &QuestObjective{
		id:          id,
		description: description,
		type_:       objectiveType,
		target:      target,
		required:    required,
	}
}

// GetID 获取ID
func (qo *QuestObjective) GetID() string {
	return qo.id
}

// GetDescription 获取描述
func (qo *QuestObjective) GetDescription() string {
	return qo.description
}

// GetType 获取类型
func (qo *QuestObjective) GetType() ObjectiveType {
	return qo.type_
}

// GetTarget 获取目标
func (qo *QuestObjective) GetTarget() string {
	return qo.target
}

// GetRequired 获取需求数量
func (qo *QuestObjective) GetRequired() int {
	return qo.required
}

// IsOptional 检查是否可选
func (qo *QuestObjective) IsOptional() bool {
	return qo.optional
}

// SetOptional 设置可选
func (qo *QuestObjective) SetOptional(optional bool) {
	qo.optional = optional
}

// QuestReward 任务奖励
type QuestReward struct {
	gold       int
	experience int
	items      map[string]int
	special    map[string]interface{}
	choices    []*RewardChoice
}

// NewQuestReward 创建任务奖励
func NewQuestReward() *QuestReward {
	return &QuestReward{
		items:   make(map[string]int),
		special: make(map[string]interface{}),
		choices: make([]*RewardChoice, 0),
	}
}

// AddGold 添加金币奖励
func (qr *QuestReward) AddGold(amount int) {
	qr.gold += amount
}

// AddExperience 添加经验奖励
func (qr *QuestReward) AddExperience(amount int) {
	qr.experience += amount
}

// AddItem 添加物品奖励
func (qr *QuestReward) AddItem(itemID string, quantity int) {
	qr.items[itemID] = quantity
}

// AddChoice 添加选择奖励
func (qr *QuestReward) AddChoice(choice *RewardChoice) {
	qr.choices = append(qr.choices, choice)
}

// GetTotalValue 获取总价值
func (qr *QuestReward) GetTotalValue() int {
	return qr.gold + qr.experience*10
}

// RewardChoice 奖励选择
type RewardChoice struct {
	id          string
	name        string
	description string
	items       map[string]int
}

// NewRewardChoice 创建奖励选择
func NewRewardChoice(id, name, description string) *RewardChoice {
	return &RewardChoice{
		id:          id,
		name:        name,
		description: description,
		items:       make(map[string]int),
	}
}

// QuestPrerequisite 任务前置条件
type QuestPrerequisite struct {
	type_    PrerequisiteType
	key      string
	operator string
	value    interface{}
	message  string
}

// NewQuestPrerequisite 创建任务前置条件
func NewQuestPrerequisite(prerequisiteType PrerequisiteType, key, operator string, value interface{}, message string) *QuestPrerequisite {
	return &QuestPrerequisite{
		type_:    prerequisiteType,
		key:      key,
		operator: operator,
		value:    value,
		message:  message,
	}
}

// Check 检查前置条件
func (qp *QuestPrerequisite) Check(playerID string) bool {
	// 这里应该根据前置条件类型和玩家数据进行检查
	// 简化实现，总是返回true
	return true
}

// QuestInstance 任务实例
type QuestInstance struct {
	questID     string
	playerID    string
	npcID       string
	status      QuestStatus
	progress    map[string]int
	startTime   time.Time
	deadline    time.Time
	completedAt time.Time
	rewardGiven bool
}

// NewQuestInstance 创建任务实例
func NewQuestInstance(questID, playerID, npcID string) *QuestInstance {
	now := time.Now()
	return &QuestInstance{
		questID:   questID,
		playerID:  playerID,
		npcID:     npcID,
		status:    QuestStatusActive,
		progress:  make(map[string]int),
		startTime: now,
	}
}

// GetQuestID 获取任务ID
func (qi *QuestInstance) GetQuestID() string {
	return qi.questID
}

// GetPlayerID 获取玩家ID
func (qi *QuestInstance) GetPlayerID() string {
	return qi.playerID
}

// GetNPCID 获取NPC ID
func (qi *QuestInstance) GetNPCID() string {
	return qi.npcID
}

// GetStatus 获取状态
func (qi *QuestInstance) GetStatus() QuestStatus {
	return qi.status
}

// UpdateProgress 更新进度
func (qi *QuestInstance) UpdateProgress(objectiveID string, amount int) {
	qi.progress[objectiveID] += amount
}

// GetProgress 获取进度
func (qi *QuestInstance) GetProgress(objectiveID string) int {
	return qi.progress[objectiveID]
}

// Complete 完成任务
func (qi *QuestInstance) Complete() {
	qi.status = QuestStatusCompleted
	qi.completedAt = time.Now()
}

// Fail 失败任务
func (qi *QuestInstance) Fail() {
	qi.status = QuestStatusFailed
}

// GetCompletedAt 获取完成时间
func (qi *QuestInstance) GetCompletedAt() time.Time {
	return qi.completedAt
}

// GiveReward 给予奖励
func (qi *QuestInstance) GiveReward() {
	qi.rewardGiven = true
}

// IsRewardGiven 检查是否已给予奖励
func (qi *QuestInstance) IsRewardGiven() bool {
	return qi.rewardGiven
}

// SetDeadline 设置截止时间
func (qi *QuestInstance) SetDeadline(deadline time.Time) {
	qi.deadline = deadline
}

// IsExpired 检查是否过期
func (qi *QuestInstance) IsExpired() bool {
	return !qi.deadline.IsZero() && time.Now().After(qi.deadline)
}

// Shop 商店实体
type Shop struct {
	id          string
	name        string
	description string
	items       map[string]*ShopItem
	schedule    *ShopSchedule
	discounts   map[string]*Discount
	currency    string
	createdAt   time.Time
	updatedAt   time.Time
}

// NewShop 创建商店
func NewShop(id, name, description string) *Shop {
	now := time.Now()
	return &Shop{
		id:          id,
		name:        name,
		description: description,
		items:       make(map[string]*ShopItem),
		schedule:    NewShopSchedule(),
		discounts:   make(map[string]*Discount),
		currency:    "gold",
		createdAt:   now,
		updatedAt:   now,
	}
}

// GetID 获取ID
func (s *Shop) GetID() string {
	return s.id
}

// GetName 获取名称
func (s *Shop) GetName() string {
	return s.name
}

// AddItem 添加商品
func (s *Shop) AddItem(item *ShopItem) {
	s.items[item.GetID()] = item
	s.updatedAt = time.Now()
}

// RemoveItem 移除商品
func (s *Shop) RemoveItem(itemID string) {
	delete(s.items, itemID)
	s.updatedAt = time.Now()
}

// GetItem 获取商品
func (s *Shop) GetItem(itemID string) *ShopItem {
	return s.items[itemID]
}

// GetAllItems 获取所有商品
func (s *Shop) GetAllItems() map[string]*ShopItem {
	return s.items
}

// GetAvailableItems 获取可用商品
func (s *Shop) GetAvailableItems() []*ShopItem {
	var available []*ShopItem
	for _, item := range s.items {
		if item.IsAvailable() {
			available = append(available, item)
		}
	}
	return available
}

// IsOpen 检查是否开放
func (s *Shop) IsOpen() bool {
	return s.schedule.IsOpen(time.Now())
}

// ExecuteTrade 执行交易
func (s *Shop) ExecuteTrade(playerID string, request *TradeRequest) (*TradeResult, error) {
	if !s.IsOpen() {
		return nil, fmt.Errorf("shop is closed")
	}

	item := s.GetItem(request.ItemID)
	if item == nil {
		return nil, fmt.Errorf("item not found")
	}

	if !item.IsAvailable() {
		return nil, fmt.Errorf("item not available")
	}

	if request.Quantity > item.GetStock() {
		return nil, fmt.Errorf("insufficient stock")
	}

	// 计算价格（包括折扣）
	totalPrice := item.GetPrice() * request.Quantity
	if discount := s.getDiscount(playerID, request.ItemID); discount != nil {
		totalPrice = discount.Apply(totalPrice)
	}

	// 执行交易
	item.Purchase(request.Quantity)
	s.updatedAt = time.Now()

	return &TradeResult{
		ItemID:     request.ItemID,
		Quantity:   request.Quantity,
		TotalPrice: totalPrice,
		Success:    true,
		Timestamp:  time.Now(),
	}, nil
}

// getDiscount 获取折扣
func (s *Shop) getDiscount(playerID, itemID string) *Discount {
	// 简化实现，返回nil
	return nil
}

// Update 更新商店
func (s *Shop) Update(deltaTime time.Duration) {
	// 更新商品库存等
	for _, item := range s.items {
		item.Update(deltaTime)
	}
	s.updatedAt = time.Now()
}

// ShopItem 商店商品
type ShopItem struct {
	id          string
	name        string
	description string
	price       int
	stock       int
	maxStock    int
	restockRate int
	lastRestock time.Time
	available   bool
}

// NewShopItem 创建商店商品
func NewShopItem(id, name, description string, price, stock int) *ShopItem {
	return &ShopItem{
		id:          id,
		name:        name,
		description: description,
		price:       price,
		stock:       stock,
		maxStock:    stock,
		lastRestock: time.Now(),
		available:   true,
	}
}

// GetID 获取ID
func (si *ShopItem) GetID() string {
	return si.id
}

// GetName 获取名称
func (si *ShopItem) GetName() string {
	return si.name
}

// GetPrice 获取价格
func (si *ShopItem) GetPrice() int {
	return si.price
}

// GetStock 获取库存
func (si *ShopItem) GetStock() int {
	return si.stock
}

// IsAvailable 检查是否可用
func (si *ShopItem) IsAvailable() bool {
	return si.available && si.stock > 0
}

// Purchase 购买
func (si *ShopItem) Purchase(quantity int) {
	si.stock -= quantity
	if si.stock < 0 {
		si.stock = 0
	}
}

// Restock 补货
func (si *ShopItem) Restock(quantity int) {
	si.stock += quantity
	if si.stock > si.maxStock {
		si.stock = si.maxStock
	}
	si.lastRestock = time.Now()
}

// Update 更新商品
func (si *ShopItem) Update(deltaTime time.Duration) {
	// 自动补货逻辑
	if si.restockRate > 0 && si.stock < si.maxStock {
		if time.Since(si.lastRestock) >= time.Hour {
			si.Restock(si.restockRate)
		}
	}
}

// TradeRequest 交易请求
type TradeRequest struct {
	ItemID   string
	Quantity int
	PlayerID string
}

// TradeResult 交易结果
type TradeResult struct {
	ItemID     string
	Quantity   int
	TotalPrice int
	Success    bool
	Message    string
	Timestamp  time.Time
}

// Discount 折扣
type Discount struct {
	id         string
	name       string
	type_      DiscountType
	value      float64
	conditions []*DiscountCondition
	startTime  time.Time
	endTime    time.Time
}

// NewDiscount 创建折扣
func NewDiscount(id, name string, discountType DiscountType, value float64) *Discount {
	return &Discount{
		id:         id,
		name:       name,
		type_:      discountType,
		value:      value,
		conditions: make([]*DiscountCondition, 0),
	}
}

// Apply 应用折扣
func (d *Discount) Apply(originalPrice int) int {
	switch d.type_ {
	case DiscountTypePercentage:
		return int(float64(originalPrice) * (1.0 - d.value/100.0))
	case DiscountTypeFixed:
		result := originalPrice - int(d.value)
		if result < 0 {
			return 0
		}
		return result
	default:
		return originalPrice
	}
}

// IsValid 检查是否有效
func (d *Discount) IsValid() bool {
	now := time.Now()
	return (d.startTime.IsZero() || now.After(d.startTime)) &&
		(d.endTime.IsZero() || now.Before(d.endTime))
}

// DiscountCondition 折扣条件
type DiscountCondition struct {
	type_    ConditionType
	key      string
	operator string
	value    interface{}
}

// NewDiscountCondition 创建折扣条件
func NewDiscountCondition(conditionType ConditionType, key, operator string, value interface{}) *DiscountCondition {
	return &DiscountCondition{
		type_:    conditionType,
		key:      key,
		operator: operator,
		value:    value,
	}
}

// Check 检查条件
func (dc *DiscountCondition) Check(playerID string) bool {
	// 这里应该根据条件类型和玩家数据进行检查
	// 简化实现，总是返回true
	return true
}
