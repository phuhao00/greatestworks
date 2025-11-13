package dialogue

import (
	"fmt"
	"sync"
)

// DialogueType 对话类型
type DialogueType int32

const (
	DialogueTypeNormal DialogueType = 0 // 普通对话
	DialogueTypeQuest  DialogueType = 1 // 任务对话
	DialogueTypeShop   DialogueType = 2 // 商店对话
)

// DialogueAction 对话动作
type DialogueAction int32

const (
	DialogueActionNone        DialogueAction = 0 // 无动作
	DialogueActionAcceptQuest DialogueAction = 1 // 接受任务
	DialogueActionSubmitQuest DialogueAction = 2 // 提交任务
	DialogueActionOpenShop    DialogueAction = 3 // 打开商店
)

// DialogueNode 对话节点
type DialogueNode struct {
	ID      int32             // 节点ID
	Text    string            // 对话文本
	Options []*DialogueOption // 选项列表
}

// DialogueOption 对话选项
type DialogueOption struct {
	Text     string         // 选项文本
	NextNode int32          // 下一个节点ID
	Action   DialogueAction // 触发的动作
	ActionID int32          // 动作ID（如任务ID）
}

// DialogueManager 对话管理器
type DialogueManager struct {
	mu sync.RWMutex

	ownerID      int64         // 所属玩家ID
	currentNPCID int32         // 当前交互的NPC ID
	currentNode  *DialogueNode // 当前对话节点
}

// NewDialogueManager 创建对话管理器
func NewDialogueManager(ownerID int64) *DialogueManager {
	return &DialogueManager{
		ownerID: ownerID,
	}
}

// StartDialogue 开始对话
func (dm *DialogueManager) StartDialogue(npcID int32, dialogueID int32) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// TODO: 从DataManager加载对话数据
	dm.currentNPCID = npcID
	dm.currentNode = &DialogueNode{
		ID:   dialogueID,
		Text: "你好，旅行者。",
		Options: []*DialogueOption{
			{Text: "有什么可以帮助你的吗？", NextNode: 2, Action: DialogueActionNone},
			{Text: "再见。", NextNode: -1, Action: DialogueActionNone},
		},
	}

	return nil
}

// SelectOption 选择对话选项
func (dm *DialogueManager) SelectOption(optionIndex int32) (*DialogueNode, DialogueAction, int32, error) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.currentNode == nil {
		return nil, DialogueActionNone, 0, fmt.Errorf("no active dialogue")
	}

	if optionIndex < 0 || optionIndex >= int32(len(dm.currentNode.Options)) {
		return nil, DialogueActionNone, 0, fmt.Errorf("invalid option index")
	}

	option := dm.currentNode.Options[optionIndex]

	// 处理对话结束
	if option.NextNode == -1 {
		dm.currentNode = nil
		dm.currentNPCID = 0
		return nil, option.Action, option.ActionID, nil
	}

	// TODO: 加载下一个节点
	nextNode := &DialogueNode{
		ID:      option.NextNode,
		Text:    "下一步对话...",
		Options: []*DialogueOption{},
	}

	dm.currentNode = nextNode
	return nextNode, option.Action, option.ActionID, nil
}

// EndDialogue 结束对话
func (dm *DialogueManager) EndDialogue() {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dm.currentNode = nil
	dm.currentNPCID = 0
}

// GetCurrentNode 获取当前对话节点
func (dm *DialogueManager) GetCurrentNode() *DialogueNode {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.currentNode
}

// GetCurrentNPCID 获取当前NPC ID
func (dm *DialogueManager) GetCurrentNPCID() int32 {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.currentNPCID
}

// ShopItem 商店物品
type ShopItem struct {
	ItemDefineID int32 // 物品定义ID
	Price        int32 // 价格
	Stock        int32 // 库存（-1表示无限）
}

// Shop 商店
type Shop struct {
	ID    int32       // 商店ID
	Name  string      // 商店名称
	Items []*ShopItem // 商品列表
}

// ShopManager 商店管理器
type ShopManager struct {
	mu sync.RWMutex

	shops map[int32]*Shop // 商店ID -> 商店
}

// NewShopManager 创建商店管理器
func NewShopManager() *ShopManager {
	return &ShopManager{
		shops: make(map[int32]*Shop),
	}
}

// GetShop 获取商店
func (sm *ShopManager) GetShop(shopID int32) *Shop {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.shops[shopID]
}

// Buy 购买物品
func (sm *ShopManager) Buy(shopID int32, itemDefineID int32, count int32) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	shop, exists := sm.shops[shopID]
	if !exists {
		return fmt.Errorf("shop not found: %d", shopID)
	}

	for _, item := range shop.Items {
		if item.ItemDefineID == itemDefineID {
			if item.Stock != -1 && item.Stock < count {
				return fmt.Errorf("insufficient stock")
			}

			// TODO: 扣除玩家金币
			// TODO: 添加物品到背包

			if item.Stock != -1 {
				item.Stock -= count
			}

			return nil
		}
	}

	return fmt.Errorf("item not found in shop")
}

// Sell 出售物品
func (sm *ShopManager) Sell(itemDefineID int32, count int32) error {
	// TODO: 从背包移除物品
	// TODO: 增加玩家金币
	return nil
}
