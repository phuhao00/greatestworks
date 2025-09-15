package synthesis

import "errors"

// 合成系统相关错误
var (
	ErrInvalidRecipe        = errors.New("invalid recipe")
	ErrRecipeNotFound       = errors.New("recipe not found")
	ErrInvalidMaterial      = errors.New("invalid material")
	ErrMaterialNotFound     = errors.New("material not found")
	ErrInsufficientMaterial = errors.New("insufficient material")
	ErrInvalidQuantity      = errors.New("invalid quantity")
	ErrRecipeAlreadyExists  = errors.New("recipe already exists")
	ErrInsufficientLevel    = errors.New("insufficient level")
	ErrConditionNotMet      = errors.New("crafting condition not met")
	ErrCraftingInProgress   = errors.New("crafting already in progress")
	ErrInvalidCategory      = errors.New("invalid recipe category")
	ErrInvalidQuality       = errors.New("invalid material quality")
	ErrMaxStackExceeded     = errors.New("max stack size exceeded")
	ErrSynthesisFailed      = errors.New("synthesis failed")
	ErrInvalidBonus         = errors.New("invalid synthesis bonus")
)