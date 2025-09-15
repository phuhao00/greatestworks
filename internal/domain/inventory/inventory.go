// Package inventory 背包领域
package inventory

import (
	"errors"
	"fmt"
	"time"
	"github.com/google/uuid"
	"greatestworks/internal/domain/player"
)

// InventoryID 背包ID值对象
type InventoryID struct {
	value string
}

// NewInventoryID 创建新的背包ID
func NewInventoryID() InventoryID {
	return InventoryID{value: uuid.New().String()}
}

// String 返回字符串表示
func (id InventoryID) String() string {
	return id.value
}

// ItemID 物品ID值对象
type ItemID struct {
	value string
}

// NewItemID 创建新的物品ID
func NewItemID() ItemID {
	return ItemID{value: uuid.New().String()}
}

// String 返回字符串表示
func (id ItemID) String() string {
	return id.value
}

// ItemType 物品类型枚举
type ItemType int

const (
	ItemTypeWeapon ItemType = iota
	ItemTypeArmor
	ItemTypeConsumable
	ItemTypeMaterial
	ItemTypeQuest
	ItemTypeCurrency
)

// ItemRarity 物品稀有度枚举
type ItemRarity int

const (
	ItemRarityCommon ItemRarity = iota
	ItemRarityUncommon
	ItemRarityRare
	ItemRarityEpic
	ItemRarityLegendary
)

// Item 物品实体
type Item struct {
	id          ItemID
	name        string
	description string
	itemType    ItemType
	rarity      ItemRarity
	maxStack    int
	value       int
	attributes  map[string]int
	createdAt   time.Time
}

// NewItem 创建新物品
func NewItem(name, description string, itemType ItemType, rarity ItemRarity, maxStack, value int) *Item {
	return &Item{
		id:          NewItemID(),
		name:        name,
		description: description,
		itemType:    itemType,
		rarity:      rarity,
		maxStack:    maxStack,
		value:       value,
		attributes:  make(map[string]int),
		createdAt:   time.Now(),
	}
}

// ID 获取物品ID
func (i *Item) ID() ItemID {
	return i.id
}

// Name 获取物品名称
func (i *Item) Name() string {
	return i.name
}

// Type 获取物品类型
func (i *Item) Type() ItemType {
	return i.itemType
}

// Rarity 获取物品稀有度
func (i *Item) Rarity() ItemRarity {
	return i.rarity
}

// MaxStack 获取最大堆叠数量
func (i *Item) MaxStack() int {
	return i.maxStack
}

// Value 获取物品价值
func (i *Item) Value() int {
	return i.value
}

// InventorySlot 背包槽位
type InventorySlot struct {
	SlotIndex int     `json:"slot_index"`
	ItemID    *ItemID `json:"item_id,omitempty"`
	Quantity  int     `json:"quantity"`
	Locked    bool    `json:"locked"`
}

// IsEmpty 是否为空槽位
func (s *InventorySlot) IsEmpty() bool {
	return s.ItemID == nil || s.Quantity <= 0
}

// CanStack 是否可以堆叠指定物品
func (s *InventorySlot) CanStack(itemID ItemID, item *Item) bool {
	if s.IsEmpty() {
		return true
	}
	if s.ItemID == nil || *s.ItemID != itemID {
		return false
	}
	return s.Quantity < item.MaxStack()
}

// Inventory 背包聚合根
type Inventory struct {
	id       InventoryID
	playerID player.PlayerID
	slots    []*InventorySlot
	capacity int
	createdAt time.Time
	updatedAt time.Time
	version  int64
}

// NewInventory 创建新背包
func NewInventory(playerID player.PlayerID, capacity int) *Inventory {
	now := time.Now()
	slots := make([]*InventorySlot, capacity)
	for i := 0; i < capacity; i++ {
		slots[i] = &InventorySlot{
			SlotIndex: i,
			Quantity:  0,
			Locked:    false,
		}
	}
	
	return &Inventory{
		id:       NewInventoryID(),
		playerID: playerID,
		slots:    slots,
		capacity: capacity,
		createdAt: now,
		updatedAt: now,
		version:  1,
	}
}

// ID 获取背包ID
func (inv *Inventory) ID() InventoryID {
	return inv.id
}

// PlayerID 获取玩家ID
func (inv *Inventory) PlayerID() string {
	return inv.playerID
}

// Capacity 获取背包容量
func (inv *Inventory) Capacity() int {
	return inv.capacity
}

// UsedSlots 获取已使用槽位
func (inv *Inventory) UsedSlots() int {
	return inv.usedSlots
}

// Items 获取所有物品
func (inv *Inventory) Items() map[string]*Item {
	return inv.items
}

// GetItem 获取指定物品
func (inv *Inventory) GetItem(itemID string) (*Item, bool) {
	item, exists := inv.items[itemID]
	return item, exists
}

// Slots 获取所有槽位
func (inv *Inventory) Slots() []*InventorySlot {
	return inv.slots
}

// AddItem 添加物品
func (inv *Inventory) AddItem(item *Item, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	
	remainingQuantity := quantity
	
	// 首先尝试堆叠到现有槽位
	for _, slot := range inv.slots {
		if slot.CanStack(item.ID(), item) && !slot.Locked {
			if slot.IsEmpty() {
				// 空槽位
				addQuantity := remainingQuantity
				if addQuantity > item.MaxStack() {
					addQuantity = item.MaxStack()
				}
				slot.ItemID = &item.id
				slot.Quantity = addQuantity
				remainingQuantity -= addQuantity
			} else {
				// 已有相同物品的槽位
				canAdd := item.MaxStack() - slot.Quantity
				if canAdd > 0 {
					addQuantity := remainingQuantity
					if addQuantity > canAdd {
						addQuantity = canAdd
					}
					slot.Quantity += addQuantity
					remainingQuantity -= addQuantity
				}
			}
			
			if remainingQuantity <= 0 {
				break
			}
		}
	}
	
	if remainingQuantity > 0 {
		return ErrInventoryFull
	}
	
	inv.updatedAt = time.Now()
	inv.version++
	return nil
}

// RemoveItem 移除物品
func (inv *Inventory) RemoveItem(itemID ItemID, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	
	// 检查是否有足够的物品
	totalQuantity := inv.GetItemQuantity(itemID)
	if totalQuantity < quantity {
		return ErrInsufficientItems
	}
	
	remainingToRemove := quantity
	
	// 从槽位中移除物品
	for _, slot := range inv.slots {
		if !slot.IsEmpty() && slot.ItemID != nil && *slot.ItemID == itemID && !slot.Locked {
			if slot.Quantity <= remainingToRemove {
				// 移除整个槽位的物品
				remainingToRemove -= slot.Quantity
				slot.ItemID = nil
				slot.Quantity = 0
			} else {
				// 部分移除
				slot.Quantity -= remainingToRemove
				remainingToRemove = 0
			}
			
			if remainingToRemove <= 0 {
				break
			}
		}
	}
	
	inv.updatedAt = time.Now()
	inv.version++
	return nil
}

// GetItemQuantity 获取物品总数量
func (inv *Inventory) GetItemQuantity(itemID ItemID) int {
	total := 0
	for _, slot := range inv.slots {
		if !slot.IsEmpty() && slot.ItemID != nil && *slot.ItemID == itemID {
			total += slot.Quantity
		}
	}
	return total
}

// HasItem 检查是否拥有指定数量的物品
func (inv *Inventory) HasItem(itemID ItemID, quantity int) bool {
	return inv.GetItemQuantity(itemID) >= quantity
}

// GetEmptySlotCount 获取空槽位数量
func (inv *Inventory) GetEmptySlotCount() int {
	count := 0
	for _, slot := range inv.slots {
		if slot.IsEmpty() && !slot.Locked {
			count++
		}
	}
	return count
}

// IsFull 检查背包是否已满
func (inv *Inventory) IsFull() bool {
	return inv.GetEmptySlotCount() == 0
}

// MoveItem 移动物品到指定槽位
func (inv *Inventory) MoveItem(fromSlot, toSlot int) error {
	if fromSlot < 0 || fromSlot >= inv.capacity || toSlot < 0 || toSlot >= inv.capacity {
		return ErrInvalidSlot
	}
	
	if fromSlot == toSlot {
		return nil
	}
	
	from := inv.slots[fromSlot]
	to := inv.slots[toSlot]
	
	if from.IsEmpty() {
		return ErrSlotEmpty
	}
	
	if from.Locked || to.Locked {
		return ErrSlotLocked
	}
	
	// 交换槽位内容
	from.ItemID, to.ItemID = to.ItemID, from.ItemID
	from.Quantity, to.Quantity = to.Quantity, from.Quantity
	
	inv.updatedAt = time.Now()
	inv.version++
	return nil
}

// LockSlot 锁定槽位
func (inv *Inventory) LockSlot(slotIndex int) error {
	if slotIndex < 0 || slotIndex >= inv.capacity {
		return ErrInvalidSlot
	}
	
	inv.slots[slotIndex].Locked = true
	inv.updatedAt = time.Now()
	inv.version++
	return nil
}

// UnlockSlot 解锁槽位
func (inv *Inventory) UnlockSlot(slotIndex int) error {
	if slotIndex < 0 || slotIndex >= inv.capacity {
		return ErrInvalidSlot
	}
	
	inv.slots[slotIndex].Locked = false
	inv.updatedAt = time.Now()
	inv.version++
	return nil
}

// Version 获取版本号
func (inv *Inventory) Version() int64 {
	return inv.version
}