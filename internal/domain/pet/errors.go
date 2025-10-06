package pet

import (
	"fmt"
)

// 宠物领域错误定义

// PetError 宠物错误基础接口
type PetError interface {
	error
	GetCode() string
	GetMessage() string
	GetDetails() map[string]interface{}
	IsRetryable() bool
	GetSeverity() ErrorSeverity
}

// ErrorSeverity 错误严重程度
type ErrorSeverity int

const (
	ErrorSeverityLow ErrorSeverity = iota
	ErrorSeverityMedium
	ErrorSeverityHigh
	ErrorSeverityCritical
)

// String 返回错误严重程度的字符串表示
func (s ErrorSeverity) String() string {
	switch s {
	case ErrorSeverityLow:
		return "low"
	case ErrorSeverityMedium:
		return "medium"
	case ErrorSeverityHigh:
		return "high"
	case ErrorSeverityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// BasePetError 宠物错误基础结构
type BasePetError struct {
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details"`
	Retryable bool                   `json:"retryable"`
	Severity  ErrorSeverity          `json:"severity"`
}

// Error 实现error接口
func (e *BasePetError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// GetCode 获取错误代码
func (e *BasePetError) GetCode() string {
	return e.Code
}

// GetMessage 获取错误消息
func (e *BasePetError) GetMessage() string {
	return e.Message
}

// GetDetails 获取错误详情
func (e *BasePetError) GetDetails() map[string]interface{} {
	return e.Details
}

// IsRetryable 是否可重试
func (e *BasePetError) IsRetryable() bool {
	return e.Retryable
}

// GetSeverity 获取错误严重程度
func (e *BasePetError) GetSeverity() ErrorSeverity {
	return e.Severity
}

// 宠物相关错误

// PetNotFoundError 宠物未找到错误
type PetNotFoundError struct {
	*BasePetError
	PetID string `json:"pet_id"`
}

// NewPetNotFoundError 创建宠物未找到错误
func NewPetNotFoundError(petID string) *PetNotFoundError {
	return &PetNotFoundError{
		BasePetError: &BasePetError{
			Code:      "PET_NOT_FOUND",
			Message:   fmt.Sprintf("Pet with ID %s not found", petID),
			Details:   map[string]interface{}{"pet_id": petID},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		PetID: petID,
	}
}

// PetAlreadyExistsError 宠物已存在错误
type PetAlreadyExistsError struct {
	*BasePetError
	PetID string `json:"pet_id"`
}

// NewPetAlreadyExistsError 创建宠物已存在错误
func NewPetAlreadyExistsError(petID string) *PetAlreadyExistsError {
	return &PetAlreadyExistsError{
		BasePetError: &BasePetError{
			Code:      "PET_ALREADY_EXISTS",
			Message:   fmt.Sprintf("Pet with ID %s already exists", petID),
			Details:   map[string]interface{}{"pet_id": petID},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		PetID: petID,
	}
}

// PetInvalidStateError 宠物状态无效错误
type PetInvalidStateError struct {
	*BasePetError
	PetID         string   `json:"pet_id"`
	CurrentState  PetState `json:"current_state"`
	RequiredState PetState `json:"required_state"`
	Operation     string   `json:"operation"`
}

// NewPetInvalidStateError 创建宠物状态无效错误
func NewPetInvalidStateError(petID string, currentState, requiredState PetState, operation string) *PetInvalidStateError {
	return &PetInvalidStateError{
		BasePetError: &BasePetError{
			Code:    "PET_INVALID_STATE",
			Message: fmt.Sprintf("Pet %s is in state %s, but %s is required for operation %s", petID, currentState, requiredState, operation),
			Details: map[string]interface{}{
				"pet_id":         petID,
				"current_state":  currentState,
				"required_state": requiredState,
				"operation":      operation,
			},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		PetID:         petID,
		CurrentState:  currentState,
		RequiredState: requiredState,
		Operation:     operation,
	}
}

// PetMaxLevelReachedError 宠物达到最大等级错误
type PetMaxLevelReachedError struct {
	*BasePetError
	PetID    string `json:"pet_id"`
	MaxLevel uint32 `json:"max_level"`
}

// NewPetMaxLevelReachedError 创建宠物达到最大等级错误
func NewPetMaxLevelReachedError(petID string, maxLevel uint32) *PetMaxLevelReachedError {
	return &PetMaxLevelReachedError{
		BasePetError: &BasePetError{
			Code:      "PET_MAX_LEVEL_REACHED",
			Message:   fmt.Sprintf("Pet %s has reached maximum level %d", petID, maxLevel),
			Details:   map[string]interface{}{"pet_id": petID, "max_level": maxLevel},
			Retryable: false,
			Severity:  ErrorSeverityLow,
		},
		PetID:    petID,
		MaxLevel: maxLevel,
	}
}

// PetInsufficientExperienceError 宠物经验不足错误
type PetInsufficientExperienceError struct {
	*BasePetError
	PetID              string `json:"pet_id"`
	CurrentExperience  uint64 `json:"current_experience"`
	RequiredExperience uint64 `json:"required_experience"`
}

// NewPetInsufficientExperienceError 创建宠物经验不足错误
func NewPetInsufficientExperienceError(petID string, current, required uint64) *PetInsufficientExperienceError {
	return &PetInsufficientExperienceError{
		BasePetError: &BasePetError{
			Code:    "PET_INSUFFICIENT_EXPERIENCE",
			Message: fmt.Sprintf("Pet %s has insufficient experience: %d/%d", petID, current, required),
			Details: map[string]interface{}{
				"pet_id":              petID,
				"current_experience":  current,
				"required_experience": required,
			},
			Retryable: false,
			Severity:  ErrorSeverityLow,
		},
		PetID:              petID,
		CurrentExperience:  current,
		RequiredExperience: required,
	}
}

// PetDeadError 宠物死亡错误
type PetDeadError struct {
	*BasePetError
	PetID     string `json:"pet_id"`
	Operation string `json:"operation"`
}

// NewPetDeadError 创建宠物死亡错误
func NewPetDeadError(petID, operation string) *PetDeadError {
	return &PetDeadError{
		BasePetError: &BasePetError{
			Code:      "PET_DEAD",
			Message:   fmt.Sprintf("Pet %s is dead and cannot perform operation: %s", petID, operation),
			Details:   map[string]interface{}{"pet_id": petID, "operation": operation},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		PetID:     petID,
		Operation: operation,
	}
}

// 宠物技能相关错误

// PetSkillNotFoundError 宠物技能未找到错误
type PetSkillNotFoundError struct {
	*BasePetError
	PetID   string `json:"pet_id"`
	SkillID string `json:"skill_id"`
}

// NewPetSkillNotFoundError 创建宠物技能未找到错误
func NewPetSkillNotFoundError(petID, skillID string) *PetSkillNotFoundError {
	return &PetSkillNotFoundError{
		BasePetError: &BasePetError{
			Code:      "PET_SKILL_NOT_FOUND",
			Message:   fmt.Sprintf("Skill %s not found for pet %s", skillID, petID),
			Details:   map[string]interface{}{"pet_id": petID, "skill_id": skillID},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		PetID:   petID,
		SkillID: skillID,
	}
}

// PetSkillOnCooldownError 宠物技能冷却中错误
type PetSkillOnCooldownError struct {
	*BasePetError
	PetID             string `json:"pet_id"`
	SkillID           string `json:"skill_id"`
	RemainingCooldown int64  `json:"remaining_cooldown"`
}

// NewPetSkillOnCooldownError 创建宠物技能冷却中错误
func NewPetSkillOnCooldownError(petID, skillID string, remainingCooldown int64) *PetSkillOnCooldownError {
	return &PetSkillOnCooldownError{
		BasePetError: &BasePetError{
			Code:    "PET_SKILL_ON_COOLDOWN",
			Message: fmt.Sprintf("Skill %s for pet %s is on cooldown for %d seconds", skillID, petID, remainingCooldown),
			Details: map[string]interface{}{
				"pet_id":             petID,
				"skill_id":           skillID,
				"remaining_cooldown": remainingCooldown,
			},
			Retryable: true,
			Severity:  ErrorSeverityLow,
		},
		PetID:             petID,
		SkillID:           skillID,
		RemainingCooldown: remainingCooldown,
	}
}

// PetSkillMaxLevelError 宠物技能达到最大等级错误
type PetSkillMaxLevelError struct {
	*BasePetError
	PetID    string `json:"pet_id"`
	SkillID  string `json:"skill_id"`
	MaxLevel uint32 `json:"max_level"`
}

// NewPetSkillMaxLevelError 创建宠物技能达到最大等级错误
func NewPetSkillMaxLevelError(petID, skillID string, maxLevel uint32) *PetSkillMaxLevelError {
	return &PetSkillMaxLevelError{
		BasePetError: &BasePetError{
			Code:      "PET_SKILL_MAX_LEVEL",
			Message:   fmt.Sprintf("Skill %s for pet %s has reached maximum level %d", skillID, petID, maxLevel),
			Details:   map[string]interface{}{"pet_id": petID, "skill_id": skillID, "max_level": maxLevel},
			Retryable: false,
			Severity:  ErrorSeverityLow,
		},
		PetID:    petID,
		SkillID:  skillID,
		MaxLevel: maxLevel,
	}
}

// 宠物碎片相关错误

// PetFragmentNotFoundError 宠物碎片未找到错误
type PetFragmentNotFoundError struct {
	*BasePetError
	PlayerID   string `json:"player_id"`
	FragmentID uint32 `json:"fragment_id"`
}

// NewPetFragmentNotFoundError 创建宠物碎片未找到错误
func NewPetFragmentNotFoundError(playerID string, fragmentID uint32) *PetFragmentNotFoundError {
	return &PetFragmentNotFoundError{
		BasePetError: &BasePetError{
			Code:      "PET_FRAGMENT_NOT_FOUND",
			Message:   fmt.Sprintf("Fragment %d not found for player %s", fragmentID, playerID),
			Details:   map[string]interface{}{"player_id": playerID, "fragment_id": fragmentID},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		PlayerID:   playerID,
		FragmentID: fragmentID,
	}
}

// PetFragmentInsufficientError 宠物碎片不足错误
type PetFragmentInsufficientError struct {
	*BasePetError
	PlayerID         string `json:"player_id"`
	FragmentID       uint32 `json:"fragment_id"`
	CurrentQuantity  uint64 `json:"current_quantity"`
	RequiredQuantity uint64 `json:"required_quantity"`
}

// NewPetFragmentInsufficientError 创建宠物碎片不足错误
func NewPetFragmentInsufficientError(playerID string, fragmentID uint32, current, required uint64) *PetFragmentInsufficientError {
	return &PetFragmentInsufficientError{
		BasePetError: &BasePetError{
			Code:    "PET_FRAGMENT_INSUFFICIENT",
			Message: fmt.Sprintf("Insufficient fragments %d for player %s: %d/%d", fragmentID, playerID, current, required),
			Details: map[string]interface{}{
				"player_id":         playerID,
				"fragment_id":       fragmentID,
				"current_quantity":  current,
				"required_quantity": required,
			},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		PlayerID:         playerID,
		FragmentID:       fragmentID,
		CurrentQuantity:  current,
		RequiredQuantity: required,
	}
}

// 宠物皮肤相关错误

// PetSkinNotFoundError 宠物皮肤未找到错误
type PetSkinNotFoundError struct {
	*BasePetError
	SkinID string `json:"skin_id"`
}

// NewPetSkinNotFoundError 创建宠物皮肤未找到错误
func NewPetSkinNotFoundError(skinID string) *PetSkinNotFoundError {
	return &PetSkinNotFoundError{
		BasePetError: &BasePetError{
			Code:      "PET_SKIN_NOT_FOUND",
			Message:   fmt.Sprintf("Pet skin %s not found", skinID),
			Details:   map[string]interface{}{"skin_id": skinID},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		SkinID: skinID,
	}
}

// PetSkinNotUnlockedError 宠物皮肤未解锁错误
type PetSkinNotUnlockedError struct {
	*BasePetError
	SkinID   string `json:"skin_id"`
	PlayerID string `json:"player_id"`
}

// NewPetSkinNotUnlockedError 创建宠物皮肤未解锁错误
func NewPetSkinNotUnlockedError(skinID, playerID string) *PetSkinNotUnlockedError {
	return &PetSkinNotUnlockedError{
		BasePetError: &BasePetError{
			Code:      "PET_SKIN_NOT_UNLOCKED",
			Message:   fmt.Sprintf("Pet skin %s is not unlocked for player %s", skinID, playerID),
			Details:   map[string]interface{}{"skin_id": skinID, "player_id": playerID},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		SkinID:   skinID,
		PlayerID: playerID,
	}
}

// PetSkinIncompatibleError 宠物皮肤不兼容错误
type PetSkinIncompatibleError struct {
	*BasePetError
	PetID        string      `json:"pet_id"`
	SkinID       string      `json:"skin_id"`
	PetCategory  PetCategory `json:"pet_category"`
	SkinCategory PetCategory `json:"skin_category"`
}

// NewPetSkinIncompatibleError 创建宠物皮肤不兼容错误
func NewPetSkinIncompatibleError(petID, skinID string, petCategory, skinCategory PetCategory) *PetSkinIncompatibleError {
	return &PetSkinIncompatibleError{
		BasePetError: &BasePetError{
			Code:    "PET_SKIN_INCOMPATIBLE",
			Message: fmt.Sprintf("Skin %s (category: %s) is incompatible with pet %s (category: %s)", skinID, skinCategory, petID, petCategory),
			Details: map[string]interface{}{
				"pet_id":        petID,
				"skin_id":       skinID,
				"pet_category":  petCategory,
				"skin_category": skinCategory,
			},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		PetID:        petID,
		SkinID:       skinID,
		PetCategory:  petCategory,
		SkinCategory: skinCategory,
	}
}

// 宠物羁绊相关错误

// PetBondNotFoundError 宠物羁绊未找到错误
type PetBondNotFoundError struct {
	*BasePetError
	BondID string `json:"bond_id"`
}

// NewPetBondNotFoundError 创建宠物羁绊未找到错误
func NewPetBondNotFoundError(bondID string) *PetBondNotFoundError {
	return &PetBondNotFoundError{
		BasePetError: &BasePetError{
			Code:      "PET_BOND_NOT_FOUND",
			Message:   fmt.Sprintf("Pet bond %s not found", bondID),
			Details:   map[string]interface{}{"bond_id": bondID},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		BondID: bondID,
	}
}

// PetBondRequirementsNotMetError 宠物羁绊条件未满足错误
type PetBondRequirementsNotMetError struct {
	*BasePetError
	BondID       string   `json:"bond_id"`
	RequiredPets []string `json:"required_pets"`
	CurrentPets  []string `json:"current_pets"`
	MissingPets  []string `json:"missing_pets"`
}

// NewPetBondRequirementsNotMetError 创建宠物羁绊条件未满足错误
func NewPetBondRequirementsNotMetError(bondID string, required, current, missing []string) *PetBondRequirementsNotMetError {
	return &PetBondRequirementsNotMetError{
		BasePetError: &BasePetError{
			Code:    "PET_BOND_REQUIREMENTS_NOT_MET",
			Message: fmt.Sprintf("Bond %s requirements not met, missing pets: %v", bondID, missing),
			Details: map[string]interface{}{
				"bond_id":       bondID,
				"required_pets": required,
				"current_pets":  current,
				"missing_pets":  missing,
			},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		BondID:       bondID,
		RequiredPets: required,
		CurrentPets:  current,
		MissingPets:  missing,
	}
}

// 宠物图鉴相关错误

// PetPictorialNotFoundError 宠物图鉴未找到错误
type PetPictorialNotFoundError struct {
	*BasePetError
	PlayerID    string `json:"player_id"`
	PetConfigID uint32 `json:"pet_config_id"`
}

// NewPetPictorialNotFoundError 创建宠物图鉴未找到错误
func NewPetPictorialNotFoundError(playerID string, petConfigID uint32) *PetPictorialNotFoundError {
	return &PetPictorialNotFoundError{
		BasePetError: &BasePetError{
			Code:      "PET_PICTORIAL_NOT_FOUND",
			Message:   fmt.Sprintf("Pet pictorial for config %d not found for player %s", petConfigID, playerID),
			Details:   map[string]interface{}{"player_id": playerID, "pet_config_id": petConfigID},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		PlayerID:    playerID,
		PetConfigID: petConfigID,
	}
}

// PetPictorialAlreadyUnlockedError 宠物图鉴已解锁错误
type PetPictorialAlreadyUnlockedError struct {
	*BasePetError
	PlayerID    string `json:"player_id"`
	PetConfigID uint32 `json:"pet_config_id"`
}

// NewPetPictorialAlreadyUnlockedError 创建宠物图鉴已解锁错误
func NewPetPictorialAlreadyUnlockedError(playerID string, petConfigID uint32) *PetPictorialAlreadyUnlockedError {
	return &PetPictorialAlreadyUnlockedError{
		BasePetError: &BasePetError{
			Code:      "PET_PICTORIAL_ALREADY_UNLOCKED",
			Message:   fmt.Sprintf("Pet pictorial for config %d is already unlocked for player %s", petConfigID, playerID),
			Details:   map[string]interface{}{"player_id": playerID, "pet_config_id": petConfigID},
			Retryable: false,
			Severity:  ErrorSeverityLow,
		},
		PlayerID:    playerID,
		PetConfigID: petConfigID,
	}
}

// 资源相关错误

// PetInsufficientResourcesError 宠物资源不足错误
type PetInsufficientResourcesError struct {
	*BasePetError
	PlayerID        string           `json:"player_id"`
	ResourceType    string           `json:"resource_type"`
	CurrentAmount   int64            `json:"current_amount"`
	RequiredAmount  int64            `json:"required_amount"`
	Operation       string           `json:"operation"`
	AdditionalCosts map[string]int64 `json:"additional_costs,omitempty"`
}

// NewPetInsufficientResourcesError 创建宠物资源不足错误
func NewPetInsufficientResourcesError(playerID, resourceType string, current, required int64, operation string) *PetInsufficientResourcesError {
	return &PetInsufficientResourcesError{
		BasePetError: &BasePetError{
			Code:    "PET_INSUFFICIENT_RESOURCES",
			Message: fmt.Sprintf("Insufficient %s for player %s: %d/%d (operation: %s)", resourceType, playerID, current, required, operation),
			Details: map[string]interface{}{
				"player_id":       playerID,
				"resource_type":   resourceType,
				"current_amount":  current,
				"required_amount": required,
				"operation":       operation,
			},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		PlayerID:       playerID,
		ResourceType:   resourceType,
		CurrentAmount:  current,
		RequiredAmount: required,
		Operation:      operation,
	}
}

// 配置相关错误

// PetConfigNotFoundError 宠物配置未找到错误
type PetConfigNotFoundError struct {
	*BasePetError
	ConfigID uint32 `json:"config_id"`
}

// NewPetConfigNotFoundError 创建宠物配置未找到错误
func NewPetConfigNotFoundError(configID uint32) *PetConfigNotFoundError {
	return &PetConfigNotFoundError{
		BasePetError: &BasePetError{
			Code:      "PET_CONFIG_NOT_FOUND",
			Message:   fmt.Sprintf("Pet configuration %d not found", configID),
			Details:   map[string]interface{}{"config_id": configID},
			Retryable: false,
			Severity:  ErrorSeverityHigh,
		},
		ConfigID: configID,
	}
}

// PetConfigInvalidError 宠物配置无效错误
type PetConfigInvalidError struct {
	*BasePetError
	ConfigID uint32 `json:"config_id"`
	Reason   string `json:"reason"`
}

// NewPetConfigInvalidError 创建宠物配置无效错误
func NewPetConfigInvalidError(configID uint32, reason string) *PetConfigInvalidError {
	return &PetConfigInvalidError{
		BasePetError: &BasePetError{
			Code:      "PET_CONFIG_INVALID",
			Message:   fmt.Sprintf("Pet configuration %d is invalid: %s", configID, reason),
			Details:   map[string]interface{}{"config_id": configID, "reason": reason},
			Retryable: false,
			Severity:  ErrorSeverityHigh,
		},
		ConfigID: configID,
		Reason:   reason,
	}
}

// 业务逻辑错误

// PetOperationNotAllowedError 宠物操作不允许错误
type PetOperationNotAllowedError struct {
	*BasePetError
	PetID     string `json:"pet_id"`
	Operation string `json:"operation"`
	Reason    string `json:"reason"`
}

// NewPetOperationNotAllowedError 创建宠物操作不允许错误
func NewPetOperationNotAllowedError(petID, operation, reason string) *PetOperationNotAllowedError {
	return &PetOperationNotAllowedError{
		BasePetError: &BasePetError{
			Code:      "PET_OPERATION_NOT_ALLOWED",
			Message:   fmt.Sprintf("Operation %s not allowed for pet %s: %s", operation, petID, reason),
			Details:   map[string]interface{}{"pet_id": petID, "operation": operation, "reason": reason},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		PetID:     petID,
		Operation: operation,
		Reason:    reason,
	}
}

// PetLimitExceededError 宠物限制超出错误
type PetLimitExceededError struct {
	*BasePetError
	PlayerID     string `json:"player_id"`
	CurrentCount int32  `json:"current_count"`
	MaxCount     int32  `json:"max_count"`
	LimitType    string `json:"limit_type"`
}

// NewPetLimitExceededError 创建宠物限制超出错误
func NewPetLimitExceededError(playerID string, current, max int32, limitType string) *PetLimitExceededError {
	return &PetLimitExceededError{
		BasePetError: &BasePetError{
			Code:    "PET_LIMIT_EXCEEDED",
			Message: fmt.Sprintf("%s limit exceeded for player %s: %d/%d", limitType, playerID, current, max),
			Details: map[string]interface{}{
				"player_id":     playerID,
				"current_count": current,
				"max_count":     max,
				"limit_type":    limitType,
			},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		PlayerID:     playerID,
		CurrentCount: current,
		MaxCount:     max,
		LimitType:    limitType,
	}
}

// 系统错误

// PetSystemError 宠物系统错误
type PetSystemError struct {
	*BasePetError
	SystemComponent string `json:"system_component"`
	InternalError   error  `json:"internal_error,omitempty"`
}

// NewPetSystemError 创建宠物系统错误
func NewPetSystemError(component, message string, internalErr error) *PetSystemError {
	return &PetSystemError{
		BasePetError: &BasePetError{
			Code:      "PET_SYSTEM_ERROR",
			Message:   fmt.Sprintf("System error in %s: %s", component, message),
			Details:   map[string]interface{}{"system_component": component},
			Retryable: true,
			Severity:  ErrorSeverityCritical,
		},
		SystemComponent: component,
		InternalError:   internalErr,
	}
}

// PetDatabaseError 宠物数据库错误
type PetDatabaseError struct {
	*BasePetError
	Operation     string `json:"operation"`
	Table         string `json:"table"`
	InternalError error  `json:"internal_error,omitempty"`
}

// NewPetDatabaseError 创建宠物数据库错误
func NewPetDatabaseError(operation, table, message string, internalErr error) *PetDatabaseError {
	return &PetDatabaseError{
		BasePetError: &BasePetError{
			Code:      "PET_DATABASE_ERROR",
			Message:   fmt.Sprintf("Database error during %s on table %s: %s", operation, table, message),
			Details:   map[string]interface{}{"operation": operation, "table": table},
			Retryable: true,
			Severity:  ErrorSeverityHigh,
		},
		Operation:     operation,
		Table:         table,
		InternalError: internalErr,
	}
}

// PetCacheError 宠物缓存错误
type PetCacheError struct {
	*BasePetError
	Operation     string `json:"operation"`
	Key           string `json:"key"`
	InternalError error  `json:"internal_error,omitempty"`
}

// NewPetCacheError 创建宠物缓存错误
func NewPetCacheError(operation, key, message string, internalErr error) *PetCacheError {
	return &PetCacheError{
		BasePetError: &BasePetError{
			Code:      "PET_CACHE_ERROR",
			Message:   fmt.Sprintf("Cache error during %s for key %s: %s", operation, key, message),
			Details:   map[string]interface{}{"operation": operation, "key": key},
			Retryable: true,
			Severity:  ErrorSeverityMedium,
		},
		Operation:     operation,
		Key:           key,
		InternalError: internalErr,
	}
}

// 验证错误

// PetValidationError 宠物验证错误
type PetValidationError struct {
	*BasePetError
	Field          string      `json:"field"`
	Value          interface{} `json:"value"`
	Constraint     string      `json:"constraint"`
	ValidationRule string      `json:"validation_rule"`
}

// NewPetValidationError 创建宠物验证错误
func NewPetValidationError(field string, value interface{}, constraint, rule string) *PetValidationError {
	return &PetValidationError{
		BasePetError: &BasePetError{
			Code:    "PET_VALIDATION_ERROR",
			Message: fmt.Sprintf("Validation failed for field %s: %s (rule: %s)", field, constraint, rule),
			Details: map[string]interface{}{
				"field":           field,
				"value":           value,
				"constraint":      constraint,
				"validation_rule": rule,
			},
			Retryable: false,
			Severity:  ErrorSeverityMedium,
		},
		Field:          field,
		Value:          value,
		Constraint:     constraint,
		ValidationRule: rule,
	}
}

// 错误代码常量

const (
	// 宠物相关错误代码
	ErrCodePetNotFound        = "PET_NOT_FOUND"
	ErrCodePetAlreadyExists   = "PET_ALREADY_EXISTS"
	ErrCodePetInvalidState    = "PET_INVALID_STATE"
	ErrCodePetMaxLevelReached = "PET_MAX_LEVEL_REACHED"
	ErrCodePetInsufficientExp = "PET_INSUFFICIENT_EXPERIENCE"
	ErrCodePetDead            = "PET_DEAD"

	// 技能相关错误代码
	ErrCodePetSkillNotFound   = "PET_SKILL_NOT_FOUND"
	ErrCodePetSkillOnCooldown = "PET_SKILL_ON_COOLDOWN"
	ErrCodePetSkillMaxLevel   = "PET_SKILL_MAX_LEVEL"

	// 碎片相关错误代码
	ErrCodePetFragmentNotFound     = "PET_FRAGMENT_NOT_FOUND"
	ErrCodePetFragmentInsufficient = "PET_FRAGMENT_INSUFFICIENT"

	// 皮肤相关错误代码
	ErrCodePetSkinNotFound     = "PET_SKIN_NOT_FOUND"
	ErrCodePetSkinNotUnlocked  = "PET_SKIN_NOT_UNLOCKED"
	ErrCodePetSkinIncompatible = "PET_SKIN_INCOMPATIBLE"

	// 羁绊相关错误代码
	ErrCodePetBondNotFound           = "PET_BOND_NOT_FOUND"
	ErrCodePetBondRequirementsNotMet = "PET_BOND_REQUIREMENTS_NOT_MET"

	// 图鉴相关错误代码
	ErrCodePetPictorialNotFound        = "PET_PICTORIAL_NOT_FOUND"
	ErrCodePetPictorialAlreadyUnlocked = "PET_PICTORIAL_ALREADY_UNLOCKED"

	// 资源相关错误代码
	ErrCodePetInsufficientResources = "PET_INSUFFICIENT_RESOURCES"

	// 配置相关错误代码
	ErrCodePetConfigNotFound = "PET_CONFIG_NOT_FOUND"
	ErrCodePetConfigInvalid  = "PET_CONFIG_INVALID"

	// 业务逻辑错误代码
	ErrCodePetOperationNotAllowed = "PET_OPERATION_NOT_ALLOWED"
	ErrCodePetLimitExceeded       = "PET_LIMIT_EXCEEDED"

	// 系统错误代码
	ErrCodePetSystemError   = "PET_SYSTEM_ERROR"
	ErrCodePetDatabaseError = "PET_DATABASE_ERROR"
	ErrCodePetCacheError    = "PET_CACHE_ERROR"

	// 验证错误代码
	ErrCodePetValidationError = "PET_VALIDATION_ERROR"
)

// 错误工具函数

// IsPetError 检查是否为宠物错误
func IsPetError(err error) bool {
	_, ok := err.(PetError)
	return ok
}

// GetPetErrorCode 获取宠物错误代码
func GetPetErrorCode(err error) string {
	if petErr, ok := err.(PetError); ok {
		return petErr.GetCode()
	}
	return ""
}

// IsRetryablePetError 检查是否为可重试的宠物错误
func IsRetryablePetError(err error) bool {
	if petErr, ok := err.(PetError); ok {
		return petErr.IsRetryable()
	}
	return false
}

// GetPetErrorSeverity 获取宠物错误严重程度
func GetPetErrorSeverity(err error) ErrorSeverity {
	if petErr, ok := err.(PetError); ok {
		return petErr.GetSeverity()
	}
	return ErrorSeverityLow
}

// WrapPetError 包装宠物错误
func WrapPetError(err error, code, message string) PetError {
	return &BasePetError{
		Code:      code,
		Message:   fmt.Sprintf("%s: %v", message, err),
		Details:   map[string]interface{}{"wrapped_error": err.Error()},
		Retryable: IsRetryablePetError(err),
		Severity:  GetPetErrorSeverity(err),
	}
}

// FormatPetError 格式化宠物错误
func FormatPetError(err PetError) string {
	return fmt.Sprintf("[%s][%s] %s", err.GetSeverity(), err.GetCode(), err.GetMessage())
}

// LogPetError 记录宠物错误（占位符函数）
func LogPetError(err PetError) {
	// 实现错误日志记录逻辑
	fmt.Printf("PET_ERROR: %s\n", FormatPetError(err))
}

// 添加缺失的错误定义
var (
	ErrInvalidPetName         = fmt.Errorf("invalid pet name")
	ErrPetIsDead              = fmt.Errorf("pet is dead")
	ErrMaxLevelReached        = fmt.Errorf("max level reached")
	ErrMaxStarReached         = fmt.Errorf("max star reached")
	ErrInvalidStateTransition = fmt.Errorf("invalid state transition")
	ErrPetNotDead             = fmt.Errorf("pet is not dead")
	ErrReviveTimeNotReached   = fmt.Errorf("revive time not reached")
	ErrMaxSkillsReached       = fmt.Errorf("max skills reached")
	ErrSkillAlreadyExists     = fmt.Errorf("skill already exists")
	ErrSkillNotFound          = fmt.Errorf("skill not found")
	ErrSkinAlreadyOwned       = fmt.Errorf("skin already owned")
	ErrSkinNotOwned           = fmt.Errorf("skin not owned")
	ErrInvalidAmount          = fmt.Errorf("invalid amount")
	ErrInvalidFoodType        = fmt.Errorf("invalid food type")
	ErrPetNotIdle             = fmt.Errorf("pet is not idle")
	ErrPetNotTraining         = fmt.Errorf("pet is not training")
	ErrPetNotInBattle         = fmt.Errorf("pet is not in battle")
	ErrInvalidPetID           = fmt.Errorf("invalid pet ID")
	ErrInvalidPlayerID        = fmt.Errorf("invalid player ID")
	ErrInvalidPetLevel        = fmt.Errorf("invalid pet level")
	ErrInvalidPetStar         = fmt.Errorf("invalid pet star")
	ErrInvalidPetAttributes   = fmt.Errorf("invalid pet attributes")
	ErrInsufficientFragments  = fmt.Errorf("insufficient fragments")
	ErrSkinAlreadyUnlocked    = fmt.Errorf("skin already unlocked")
	ErrSkinNotUnlocked        = fmt.Errorf("skin not unlocked")
	ErrSkinAlreadyEquipped    = fmt.Errorf("skin already equipped")
	ErrMaxSkillLevelReached   = fmt.Errorf("max skill level reached")
	ErrInsufficientSkillExperience = fmt.Errorf("insufficient skill experience")
	ErrSkillOnCooldown        = fmt.Errorf("skill on cooldown")
	ErrBondAlreadyActive     = fmt.Errorf("bond already active")
	ErrMaxActiveBondsReached  = fmt.Errorf("max active bonds reached")
	ErrBondNotActive         = fmt.Errorf("bond not active")
	ErrPetTemplateNotFound   = fmt.Errorf("pet template not found")
	ErrNoFragmentsProvided   = fmt.Errorf("no fragments provided")
)
