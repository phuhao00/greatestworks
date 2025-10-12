package inventory

import "errors"

var (
	// 背包相关错误
	ErrInventoryFull   = errors.New("inventory is full")
	ErrInvalidCapacity = errors.New("invalid inventory capacity")
	ErrSlotNotFound    = errors.New("inventory slot not found")

	// 物品相关错误
	ErrItemNotFound         = errors.New("item not found")
	ErrInvalidQuantity      = errors.New("invalid item quantity")
	ErrInsufficientQuantity = errors.New("insufficient item quantity")
	ErrExceedsMaxStack      = errors.New("exceeds maximum stack size")
	ErrItemNotUsable        = errors.New("item is not usable")
	ErrItemNotEquippable    = errors.New("item is not equippable")
	ErrItemOnCooldown       = errors.New("item is on cooldown")
	ErrItemExpired          = errors.New("item has expired")
	ErrInvalidItemType      = errors.New("invalid item type")

	// 装备相关错误
	ErrInvalidEquipment  = errors.New("invalid equipment")
	ErrEquipmentDamaged  = errors.New("equipment is damaged")
	ErrInsufficientLevel = errors.New("insufficient level to equip")
	ErrClassRestriction  = errors.New("class restriction for equipment")

	// 宝石相关错误
	ErrGemSlotFull    = errors.New("gem slot is full")
	ErrInvalidGemType = errors.New("invalid gem type")
	ErrGemNotFound    = errors.New("gem not found")

	// 交易相关错误
	ErrItemNotTradeable = errors.New("item is not tradeable")
	ErrTradeRestricted  = errors.New("trade is restricted")
)
