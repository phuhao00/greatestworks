package dressup

import "errors"

// 换装系统相关错误
var (
	ErrInvalidOutfit         = errors.New("invalid outfit")
	ErrOutfitNotFound        = errors.New("outfit not found")
	ErrInvalidSlot           = errors.New("invalid slot")
	ErrInvalidSetName        = errors.New("invalid set name")
	ErrSetNotFound           = errors.New("set not found")
	ErrInvalidFashionSet     = errors.New("invalid fashion set")
	ErrInvalidDyeColor       = errors.New("invalid dye color")
	ErrDyeColorNotUnlocked   = errors.New("dye color not unlocked")
	ErrInvalidStyle          = errors.New("invalid style")
	ErrStyleNotFound         = errors.New("style not found")
	ErrAutoEquipDisabled     = errors.New("auto equip disabled")
	ErrOutfitAlreadyEquipped = errors.New("outfit already equipped")
	ErrOutfitLocked          = errors.New("outfit is locked")
	ErrInsufficientLevel     = errors.New("insufficient level")
	ErrInsufficientExp       = errors.New("insufficient experience")
	ErrMaxEnhanceLevel       = errors.New("max enhance level reached")
	ErrInvalidFilter         = errors.New("invalid filter")
	ErrNoOutfitsFound        = errors.New("no outfits found")
)