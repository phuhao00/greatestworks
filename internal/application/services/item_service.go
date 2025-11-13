package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"greatestworks/internal/infrastructure/datamanager"
	"greatestworks/internal/infrastructure/persistence"
)

// ItemService 物品服务
type ItemService struct {
	itemRepo *persistence.ItemRepository
}

// NewItemService 创建物品服务
func NewItemService(itemRepo *persistence.ItemRepository) *ItemService {
	return &ItemService{
		itemRepo: itemRepo,
	}
}

// CreateItem 创建物品
func (s *ItemService) CreateItem(ctx context.Context, characterID int64, itemID, count, slot, location int32) (int64, error) {
	// 获取物品配置
	itemDefine := datamanager.GetInstance().GetItem(itemID)
	if itemDefine == nil {
		return 0, errors.New("item not found")
	}

	// 生成物品唯一ID
	itemUID := time.Now().UnixNano()

	// 创建物品
	item := &persistence.DbItem{
		ItemUID:     itemUID,
		CharacterID: characterID,
		ItemID:      itemID,
		Count:       count,
		Slot:        slot,
		Location:    location,
		Bound:       false,
		Expire:      0,
	}

	if err := s.itemRepo.Create(ctx, item); err != nil {
		return 0, fmt.Errorf("failed to create item: %w", err)
	}

	return itemUID, nil
}

// GetItem 获取物品
func (s *ItemService) GetItem(ctx context.Context, itemUID int64) (*persistence.DbItem, error) {
	return s.itemRepo.FindByUID(ctx, itemUID)
}

// GetCharacterItems 获取角色的所有物品
func (s *ItemService) GetCharacterItems(ctx context.Context, characterID int64) ([]*persistence.DbItem, error) {
	return s.itemRepo.FindByCharacterID(ctx, characterID)
}

// UseItem 使用物品
func (s *ItemService) UseItem(ctx context.Context, itemUID int64) error {
	item, err := s.itemRepo.FindByUID(ctx, itemUID)
	if err != nil {
		return err
	}

	// 获取物品配置
	itemDefine := datamanager.GetInstance().GetItem(item.ItemID)
	if itemDefine == nil {
		return errors.New("item not found")
	}

	// TODO: 根据物品类型执行不同的使用逻辑

	// 减少数量
	item.Count--
	if item.Count <= 0 {
		// 删除物品
		return s.itemRepo.Delete(ctx, itemUID)
	}

	// 更新物品
	return s.itemRepo.Update(ctx, item)
}

// MoveItem 移动物品
func (s *ItemService) MoveItem(ctx context.Context, itemUID int64, newSlot, newLocation int32) error {
	item, err := s.itemRepo.FindByUID(ctx, itemUID)
	if err != nil {
		return err
	}

	item.Slot = newSlot
	item.Location = newLocation

	return s.itemRepo.Update(ctx, item)
}

// DeleteItem 删除物品
func (s *ItemService) DeleteItem(ctx context.Context, itemUID int64) error {
	return s.itemRepo.Delete(ctx, itemUID)
}

// SplitItem 拆分物品
func (s *ItemService) SplitItem(ctx context.Context, itemUID int64, splitCount int32) (int64, error) {
	item, err := s.itemRepo.FindByUID(ctx, itemUID)
	if err != nil {
		return 0, err
	}

	if item.Count <= splitCount {
		return 0, errors.New("invalid split count")
	}

	// 减少原物品数量
	item.Count -= splitCount
	if err := s.itemRepo.Update(ctx, item); err != nil {
		return 0, err
	}

	// 创建新物品
	newItemUID := time.Now().UnixNano()
	newItem := &persistence.DbItem{
		ItemUID:     newItemUID,
		CharacterID: item.CharacterID,
		ItemID:      item.ItemID,
		Count:       splitCount,
		Slot:        -1, // 未放置槽位
		Location:    item.Location,
		Bound:       item.Bound,
		Expire:      item.Expire,
	}

	if err := s.itemRepo.Create(ctx, newItem); err != nil {
		return 0, err
	}

	return newItemUID, nil
}

// MergeItem 合并物品
func (s *ItemService) MergeItem(ctx context.Context, fromUID, toUID int64) error {
	fromItem, err := s.itemRepo.FindByUID(ctx, fromUID)
	if err != nil {
		return err
	}

	toItem, err := s.itemRepo.FindByUID(ctx, toUID)
	if err != nil {
		return err
	}

	// 检查是否可以合并
	if fromItem.ItemID != toItem.ItemID {
		return errors.New("cannot merge different items")
	}

	// 获取物品配置
	itemDefine := datamanager.GetInstance().GetItem(fromItem.ItemID)
	if itemDefine == nil {
		return errors.New("item not found")
	}

	// 合并数量
	totalCount := fromItem.Count + toItem.Count
	if totalCount > itemDefine.MaxStack {
		// 超过堆叠上限
		toItem.Count = itemDefine.MaxStack
		fromItem.Count = totalCount - itemDefine.MaxStack
		if err := s.itemRepo.Update(ctx, fromItem); err != nil {
			return err
		}
	} else {
		// 全部合并
		toItem.Count = totalCount
		if err := s.itemRepo.Delete(ctx, fromUID); err != nil {
			return err
		}
	}

	return s.itemRepo.Update(ctx, toItem)
}
