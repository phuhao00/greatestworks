package pet

import (
	"fmt"
	"time"
)

// 宠物相关事件定义

// PetEvent 宠物事件基础接口
type PetEvent interface {
	GetEventID() string
	GetEventType() string
	GetAggregateID() string
	GetPlayerID() string
	GetTimestamp() time.Time
	GetVersion() int
	GetMetadata() map[string]interface{}
}

// BasePetEvent 宠物事件基础结构
type BasePetEvent struct {
	EventID     string                 `json:"event_id"`
	EventType   string                 `json:"event_type"`
	AggregateID string                 `json:"aggregate_id"`
	PlayerID    string                 `json:"player_id"`
	Timestamp   time.Time              `json:"timestamp"`
	Version     int                    `json:"version"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// GetEventID 获取事件ID
func (e *BasePetEvent) GetEventID() string {
	return e.EventID
}

// GetEventType 获取事件类型
func (e *BasePetEvent) GetEventType() string {
	return e.EventType
}

// GetAggregateID 获取聚合ID
func (e *BasePetEvent) GetAggregateID() string {
	return e.AggregateID
}

// GetPlayerID 获取玩家ID
func (e *BasePetEvent) GetPlayerID() string {
	return e.PlayerID
}

// GetTimestamp 获取时间戳
func (e *BasePetEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetVersion 获取版本
func (e *BasePetEvent) GetVersion() int {
	return e.Version
}

// GetMetadata 获取元数据
func (e *BasePetEvent) GetMetadata() map[string]interface{} {
	return e.Metadata
}

// 宠物生命周期事件

// PetCreatedEvent 宠物创建事件
type PetCreatedEvent struct {
	BasePetEvent
	PetID          string      `json:"pet_id"`
	PetName        string      `json:"pet_name"`
	ConfigID       uint32      `json:"config_id"`
	Category       PetCategory `json:"category"`
	Rarity         PetRarity   `json:"rarity"`
	InitialLevel   uint32      `json:"initial_level"`
	InitialStar    uint32      `json:"initial_star"`
	CreationMethod string      `json:"creation_method"` // "summon", "fragment", "purchase", "reward"
}

// PetDeletedEvent 宠物删除事件
type PetDeletedEvent struct {
	BasePetEvent
	PetID        string `json:"pet_id"`
	PetName      string `json:"pet_name"`
	FinalLevel   uint32 `json:"final_level"`
	FinalStar    uint32 `json:"final_star"`
	FinalPower   int64  `json:"final_power"`
	DeleteReason string `json:"delete_reason"`
}

// PetRenamedEvent 宠物重命名事件
type PetRenamedEvent struct {
	BasePetEvent
	PetID   string `json:"pet_id"`
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

// 宠物成长事件

// PetLevelUpEvent 宠物升级事件
type PetLevelUpEvent struct {
	BasePetEvent
	PetID            string           `json:"pet_id"`
	PetName          string           `json:"pet_name"`
	OldLevel         uint32           `json:"old_level"`
	NewLevel         uint32           `json:"new_level"`
	ExperienceGained uint64           `json:"experience_gained"`
	PowerIncrease    int64            `json:"power_increase"`
	AttributeChanges map[string]int64 `json:"attribute_changes"`
}

// PetStarUpEvent 宠物升星事件
type PetStarUpEvent struct {
	BasePetEvent
	PetID             string         `json:"pet_id"`
	PetName           string         `json:"pet_name"`
	OldStar           uint32         `json:"old_star"`
	NewStar           uint32         `json:"new_star"`
	MaterialsUsed     []MaterialUsed `json:"materials_used"`
	PowerIncrease     int64          `json:"power_increase"`
	NewSkillsUnlocked []string       `json:"new_skills_unlocked"`
}

// PetEvolvedEvent 宠物进化事件
type PetEvolvedEvent struct {
	BasePetEvent
	PetID         string         `json:"pet_id"`
	PetName       string         `json:"pet_name"`
	OldConfigID   uint32         `json:"old_config_id"`
	NewConfigID   uint32         `json:"new_config_id"`
	OldCategory   PetCategory    `json:"old_category"`
	NewCategory   PetCategory    `json:"new_category"`
	MaterialsUsed []MaterialUsed `json:"materials_used"`
	PowerIncrease int64          `json:"power_increase"`
	NewAbilities  []string       `json:"new_abilities"`
}

// 宠物状态事件

// PetStateChangedEvent 宠物状态改变事件
type PetStateChangedEvent struct {
	BasePetEvent
	PetID    string         `json:"pet_id"`
	PetName  string         `json:"pet_name"`
	OldState PetState       `json:"old_state"`
	NewState PetState       `json:"new_state"`
	Reason   string         `json:"reason"`
	Duration *time.Duration `json:"duration,omitempty"`
}

// PetMoodChangedEvent 宠物心情改变事件
type PetMoodChangedEvent struct {
	BasePetEvent
	PetID     string  `json:"pet_id"`
	PetName   string  `json:"pet_name"`
	OldMood   PetMood `json:"old_mood"`
	NewMood   PetMood `json:"new_mood"`
	Reason    string  `json:"reason"`
	MoodValue int32   `json:"mood_value"`
}

// PetHealthChangedEvent 宠物健康改变事件
type PetHealthChangedEvent struct {
	BasePetEvent
	PetID      string `json:"pet_id"`
	PetName    string `json:"pet_name"`
	OldHealth  uint32 `json:"old_health"`
	NewHealth  uint32 `json:"new_health"`
	MaxHealth  uint32 `json:"max_health"`
	Reason     string `json:"reason"`
	IsCritical bool   `json:"is_critical"`
}

// 宠物互动事件

// PetFedEvent 宠物喂食事件
type PetFedEvent struct {
	BasePetEvent
	PetID            string   `json:"pet_id"`
	PetName          string   `json:"pet_name"`
	FoodType         FoodType `json:"food_type"`
	FoodQuantity     uint32   `json:"food_quantity"`
	SatietyGained    uint32   `json:"satiety_gained"`
	MoodChange       int32    `json:"mood_change"`
	ExperienceGained uint64   `json:"experience_gained"`
	SpecialEffect    string   `json:"special_effect,omitempty"`
}

// PetTrainedEvent 宠物训练事件
type PetTrainedEvent struct {
	BasePetEvent
	PetID            string             `json:"pet_id"`
	PetName          string             `json:"pet_name"`
	TrainingType     TrainingType       `json:"training_type"`
	TrainingDuration time.Duration      `json:"training_duration"`
	ExperienceGained uint64             `json:"experience_gained"`
	AttributeGains   map[string]int64   `json:"attribute_gains"`
	SkillProgress    map[string]float64 `json:"skill_progress"`
	TrainingCost     int64              `json:"training_cost"`
}

// PetPlayedEvent 宠物游戏事件
type PetPlayedEvent struct {
	BasePetEvent
	PetID        string        `json:"pet_id"`
	PetName      string        `json:"pet_name"`
	GameType     string        `json:"game_type"`
	GameDuration time.Duration `json:"game_duration"`
	MoodIncrease int32         `json:"mood_increase"`
	BondIncrease int32         `json:"bond_increase"`
	Rewards      []GameReward  `json:"rewards"`
}

// 宠物技能事件

// PetSkillLearnedEvent 宠物学习技能事件
type PetSkillLearnedEvent struct {
	BasePetEvent
	PetID       string `json:"pet_id"`
	PetName     string `json:"pet_name"`
	SkillID     string `json:"skill_id"`
	SkillName   string `json:"skill_name"`
	SkillLevel  uint32 `json:"skill_level"`
	LearnMethod string `json:"learn_method"` // "level_up", "training", "item", "evolution"
	Cost        int64  `json:"cost"`
}

// PetSkillUpgradedEvent 宠物技能升级事件
type PetSkillUpgradedEvent struct {
	BasePetEvent
	PetID         string `json:"pet_id"`
	PetName       string `json:"pet_name"`
	SkillID       string `json:"skill_id"`
	SkillName     string `json:"skill_name"`
	OldLevel      uint32 `json:"old_level"`
	NewLevel      uint32 `json:"new_level"`
	PowerIncrease int64  `json:"power_increase"`
	UpgradeCost   int64  `json:"upgrade_cost"`
}

// PetSkillUsedEvent 宠物使用技能事件
type PetSkillUsedEvent struct {
	BasePetEvent
	PetID        string        `json:"pet_id"`
	PetName      string        `json:"pet_name"`
	SkillID      string        `json:"skill_id"`
	SkillName    string        `json:"skill_name"`
	TargetID     string        `json:"target_id,omitempty"`
	DamageDealt  int64         `json:"damage_dealt"`
	Effects      []SkillEffect `json:"effects"`
	CooldownTime time.Duration `json:"cooldown_time"`
	Context      string        `json:"context"` // "battle", "training", "exploration"
}

// 宠物装备事件

// PetSkinEquippedEvent 宠物装备皮肤事件
type PetSkinEquippedEvent struct {
	BasePetEvent
	PetID          string   `json:"pet_id"`
	PetName        string   `json:"pet_name"`
	SkinID         string   `json:"skin_id"`
	SkinName       string   `json:"skin_name"`
	OldSkinID      string   `json:"old_skin_id,omitempty"`
	PowerBonus     int64    `json:"power_bonus"`
	SpecialEffects []string `json:"special_effects"`
}

// PetSkinUnequippedEvent 宠物卸下皮肤事件
type PetSkinUnequippedEvent struct {
	BasePetEvent
	PetID     string `json:"pet_id"`
	PetName   string `json:"pet_name"`
	SkinID    string `json:"skin_id"`
	SkinName  string `json:"skin_name"`
	PowerLoss int64  `json:"power_loss"`
}

// 宠物羁绊事件

// PetBondActivatedEvent 宠物羁绊激活事件
type PetBondActivatedEvent struct {
	BasePetEvent
	BondID       string       `json:"bond_id"`
	BondName     string       `json:"bond_name"`
	PetIDs       []string     `json:"pet_ids"`
	BondLevel    uint32       `json:"bond_level"`
	BonusEffects []BondEffect `json:"bonus_effects"`
	PowerBonus   int64        `json:"power_bonus"`
}

// PetBondUpgradedEvent 宠物羁绊升级事件
type PetBondUpgradedEvent struct {
	BasePetEvent
	BondID        string       `json:"bond_id"`
	BondName      string       `json:"bond_name"`
	OldLevel      uint32       `json:"old_level"`
	NewLevel      uint32       `json:"new_level"`
	PetIDs        []string     `json:"pet_ids"`
	NewEffects    []BondEffect `json:"new_effects"`
	PowerIncrease int64        `json:"power_increase"`
}

// PetBondDeactivatedEvent 宠物羁绊失效事件
type PetBondDeactivatedEvent struct {
	BasePetEvent
	BondID    string   `json:"bond_id"`
	BondName  string   `json:"bond_name"`
	PetIDs    []string `json:"pet_ids"`
	Reason    string   `json:"reason"`
	PowerLoss int64    `json:"power_loss"`
}

// 宠物碎片事件

// PetFragmentObtainedEvent 获得宠物碎片事件
type PetFragmentObtainedEvent struct {
	BasePetEvent
	FragmentID    uint32 `json:"fragment_id"`
	RelatedPetID  uint32 `json:"related_pet_id"`
	Quantity      uint64 `json:"quantity"`
	Source        string `json:"source"` // "battle", "shop", "event", "decompose"
	TotalQuantity uint64 `json:"total_quantity"`
}

// PetFragmentUsedEvent 使用宠物碎片事件
type PetFragmentUsedEvent struct {
	BasePetEvent
	FragmentID        uint32 `json:"fragment_id"`
	RelatedPetID      uint32 `json:"related_pet_id"`
	QuantityUsed      uint64 `json:"quantity_used"`
	Purpose           string `json:"purpose"` // "summon", "upgrade", "evolution"
	RemainingQuantity uint64 `json:"remaining_quantity"`
	Result            string `json:"result"`
}

// 宠物图鉴事件

// PetPictorialUnlockedEvent 宠物图鉴解锁事件
type PetPictorialUnlockedEvent struct {
	BasePetEvent
	PetConfigID  uint32            `json:"pet_config_id"`
	PetName      string            `json:"pet_name"`
	Category     PetCategory       `json:"category"`
	Rarity       PetRarity         `json:"rarity"`
	UnlockMethod string            `json:"unlock_method"` // "obtain", "encounter", "complete"
	Rewards      []PictorialReward `json:"rewards"`
}

// PetPictorialUpdatedEvent 宠物图鉴更新事件
type PetPictorialUpdatedEvent struct {
	BasePetEvent
	PetConfigID     uint32            `json:"pet_config_id"`
	OldHighestLevel uint32            `json:"old_highest_level"`
	NewHighestLevel uint32            `json:"new_highest_level"`
	OldStar         uint32            `json:"old_star"`
	NewStar         uint32            `json:"new_star"`
	ProgressRewards []PictorialReward `json:"progress_rewards"`
}

// 宠物战斗事件

// PetBattleStartedEvent 宠物战斗开始事件
type PetBattleStartedEvent struct {
	BasePetEvent
	BattleID     string              `json:"battle_id"`
	BattleType   string              `json:"battle_type"` // "pve", "pvp", "arena", "boss"
	Participants []BattleParticipant `json:"participants"`
	BattleMode   string              `json:"battle_mode"`
	Location     string              `json:"location"`
}

// PetBattleEndedEvent 宠物战斗结束事件
type PetBattleEndedEvent struct {
	BasePetEvent
	BattleID     string         `json:"battle_id"`
	Result       string         `json:"result"` // "victory", "defeat", "draw"
	Duration     time.Duration  `json:"duration"`
	Participants []BattleResult `json:"participants"`
	Rewards      []BattleReward `json:"rewards"`
	Experience   uint64         `json:"experience"`
}

// PetDefeatedEvent 宠物被击败事件
type PetDefeatedEvent struct {
	BasePetEvent
	PetID       string     `json:"pet_id"`
	PetName     string     `json:"pet_name"`
	BattleID    string     `json:"battle_id"`
	DefeatedBy  string     `json:"defeated_by"`
	FinalDamage int64      `json:"final_damage"`
	ReviveTime  *time.Time `json:"revive_time,omitempty"`
}

// 宠物探索事件

// PetExplorationStartedEvent 宠物探索开始事件
type PetExplorationStartedEvent struct {
	BasePetEvent
	PetID           string        `json:"pet_id"`
	PetName         string        `json:"pet_name"`
	ExplorationID   string        `json:"exploration_id"`
	Location        string        `json:"location"`
	Duration        time.Duration `json:"duration"`
	ExpectedRewards []string      `json:"expected_rewards"`
}

// PetExplorationCompletedEvent 宠物探索完成事件
type PetExplorationCompletedEvent struct {
	BasePetEvent
	PetID          string              `json:"pet_id"`
	PetName        string              `json:"pet_name"`
	ExplorationID  string              `json:"exploration_id"`
	Location       string              `json:"location"`
	ActualDuration time.Duration       `json:"actual_duration"`
	Success        bool                `json:"success"`
	Rewards        []ExplorationReward `json:"rewards"`
	Experience     uint64              `json:"experience"`
	Encounters     []string            `json:"encounters"`
}

// 系统事件

// PetSystemMaintenanceEvent 宠物系统维护事件
type PetSystemMaintenanceEvent struct {
	BasePetEvent
	MaintenanceType  string         `json:"maintenance_type"` // "daily", "weekly", "emergency"
	AffectedFeatures []string       `json:"affected_features"`
	Duration         time.Duration  `json:"duration"`
	Compensation     []SystemReward `json:"compensation"`
}

// PetDataMigrationEvent 宠物数据迁移事件
type PetDataMigrationEvent struct {
	BasePetEvent
	MigrationType   string   `json:"migration_type"`
	FromVersion     string   `json:"from_version"`
	ToVersion       string   `json:"to_version"`
	AffectedRecords int64    `json:"affected_records"`
	MigrationStatus string   `json:"migration_status"`
	Errors          []string `json:"errors,omitempty"`
}

// 事件相关的辅助结构体

// MaterialUsed 使用的材料
type MaterialUsed struct {
	MaterialID   string `json:"material_id"`
	MaterialName string `json:"material_name"`
	Quantity     uint64 `json:"quantity"`
	MaterialType string `json:"material_type"`
}

// GameReward 游戏奖励
type GameReward struct {
	RewardType string `json:"reward_type"`
	RewardID   string `json:"reward_id"`
	Quantity   uint64 `json:"quantity"`
	Rarity     string `json:"rarity"`
}

// SkillEffect and BondEffect are defined in entity.go

// PictorialReward 图鉴奖励
type PictorialReward struct {
	RewardType  string `json:"reward_type"`
	RewardID    string `json:"reward_id"`
	Quantity    uint64 `json:"quantity"`
	Description string `json:"description"`
}

// BattleParticipant 战斗参与者
type BattleParticipant struct {
	PlayerID   string   `json:"player_id"`
	PlayerName string   `json:"player_name"`
	PetIDs     []string `json:"pet_ids"`
	TeamPower  int64    `json:"team_power"`
	Formation  string   `json:"formation"`
}

// BattleResult 战斗结果
type BattleResult struct {
	PlayerID       string `json:"player_id"`
	Result         string `json:"result"`
	DamageDealt    int64  `json:"damage_dealt"`
	DamageReceived int64  `json:"damage_received"`
	PetsLost       int32  `json:"pets_lost"`
	MVPPetID       string `json:"mvp_pet_id"`
}

// BattleReward 战斗奖励
type BattleReward struct {
	RewardType      string  `json:"reward_type"`
	RewardID        string  `json:"reward_id"`
	Quantity        uint64  `json:"quantity"`
	BonusMultiplier float64 `json:"bonus_multiplier"`
}

// ExplorationReward 探索奖励
type ExplorationReward struct {
	RewardType string  `json:"reward_type"`
	RewardID   string  `json:"reward_id"`
	Quantity   uint64  `json:"quantity"`
	Rarity     string  `json:"rarity"`
	FindChance float64 `json:"find_chance"`
}

// SystemReward 系统奖励
type SystemReward struct {
	RewardType string     `json:"reward_type"`
	RewardID   string     `json:"reward_id"`
	Quantity   uint64     `json:"quantity"`
	Reason     string     `json:"reason"`
	ExpireTime *time.Time `json:"expire_time,omitempty"`
}

// 事件常量

const (
	// 生命周期事件类型
	EventTypePetCreated = "pet.created"
	EventTypePetDeleted = "pet.deleted"
	EventTypePetRenamed = "pet.renamed"

	// 成长事件类型
	EventTypePetLevelUp = "pet.level_up"
	EventTypePetStarUp  = "pet.star_up"
	EventTypePetEvolved = "pet.evolved"

	// 状态事件类型
	EventTypePetStateChanged  = "pet.state_changed"
	EventTypePetMoodChanged   = "pet.mood_changed"
	EventTypePetHealthChanged = "pet.health_changed"

	// 互动事件类型
	EventTypePetFed     = "pet.fed"
	EventTypePetTrained = "pet.trained"
	EventTypePetPlayed  = "pet.played"

	// 技能事件类型
	EventTypePetSkillLearned  = "pet.skill_learned"
	EventTypePetSkillUpgraded = "pet.skill_upgraded"
	EventTypePetSkillUsed     = "pet.skill_used"

	// 装备事件类型
	EventTypePetSkinEquipped   = "pet.skin_equipped"
	EventTypePetSkinUnequipped = "pet.skin_unequipped"

	// 羁绊事件类型
	EventTypePetBondActivated   = "pet.bond_activated"
	EventTypePetBondUpgraded    = "pet.bond_upgraded"
	EventTypePetBondDeactivated = "pet.bond_deactivated"

	// 碎片事件类型
	EventTypePetFragmentObtained = "pet.fragment_obtained"
	EventTypePetFragmentUsed     = "pet.fragment_used"

	// 图鉴事件类型
	EventTypePetPictorialUnlocked = "pet.pictorial_unlocked"
	EventTypePetPictorialUpdated  = "pet.pictorial_updated"

	// 战斗事件类型
	EventTypePetBattleStarted = "pet.battle_started"
	EventTypePetBattleEnded   = "pet.battle_ended"
	EventTypePetDefeated      = "pet.defeated"

	// 探索事件类型
	EventTypePetExplorationStarted   = "pet.exploration_started"
	EventTypePetExplorationCompleted = "pet.exploration_completed"

	// 系统事件类型
	EventTypePetSystemMaintenance = "pet.system_maintenance"
	EventTypePetDataMigration     = "pet.data_migration"
)

// 事件工厂函数

// NewPetCreatedEvent 创建宠物创建事件
func NewPetCreatedEvent(petID, playerID, petName string, configID uint32, category PetCategory, rarity PetRarity, level, star uint32, method string) *PetCreatedEvent {
	return &PetCreatedEvent{
		BasePetEvent: BasePetEvent{
			EventID:     generateEventID(),
			EventType:   EventTypePetCreated,
			AggregateID: petID,
			PlayerID:    playerID,
			Timestamp:   time.Now(),
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		PetID:          petID,
		PetName:        petName,
		ConfigID:       configID,
		Category:       category,
		Rarity:         rarity,
		InitialLevel:   level,
		InitialStar:    star,
		CreationMethod: method,
	}
}

// NewPetLevelUpEvent 创建宠物升级事件
func NewPetLevelUpEvent(petID, playerID, petName string, oldLevel, newLevel uint32, expGained uint64, powerIncrease int64, attrChanges map[string]int64) *PetLevelUpEvent {
	return &PetLevelUpEvent{
		BasePetEvent: BasePetEvent{
			EventID:     generateEventID(),
			EventType:   EventTypePetLevelUp,
			AggregateID: petID,
			PlayerID:    playerID,
			Timestamp:   time.Now(),
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		PetID:            petID,
		PetName:          petName,
		OldLevel:         oldLevel,
		NewLevel:         newLevel,
		ExperienceGained: expGained,
		PowerIncrease:    powerIncrease,
		AttributeChanges: attrChanges,
	}
}

// NewPetSkillLearnedEvent 创建宠物学习技能事件
func NewPetSkillLearnedEvent(petID, playerID, petName, skillID, skillName string, skillLevel uint32, method string, cost int64) *PetSkillLearnedEvent {
	return &PetSkillLearnedEvent{
		BasePetEvent: BasePetEvent{
			EventID:     generateEventID(),
			EventType:   EventTypePetSkillLearned,
			AggregateID: petID,
			PlayerID:    playerID,
			Timestamp:   time.Now(),
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		PetID:       petID,
		PetName:     petName,
		SkillID:     skillID,
		SkillName:   skillName,
		SkillLevel:  skillLevel,
		LearnMethod: method,
		Cost:        cost,
	}
}

// 事件处理器接口

// PetEventHandler 宠物事件处理器接口
type PetEventHandler interface {
	Handle(event PetEvent) error
	CanHandle(eventType string) bool
	GetHandlerName() string
}

// PetEventBus 宠物事件总线接口
type PetEventBus interface {
	// 发布事件
	Publish(event PetEvent) error
	PublishBatch(events []PetEvent) error

	// 订阅事件
	Subscribe(eventType string, handler PetEventHandler) error
	Unsubscribe(eventType string, handlerName string) error

	// 事件存储
	Store(event PetEvent) error
	GetEvents(aggregateID string, fromVersion int) ([]PetEvent, error)
	GetEventsByType(eventType string, limit int) ([]PetEvent, error)
	GetEventsByPlayer(playerID string, limit int) ([]PetEvent, error)

	// 事件重放
	Replay(aggregateID string, fromVersion int, handler PetEventHandler) error

	// 快照管理
	CreateSnapshot(aggregateID string, version int, data interface{}) error
	GetSnapshot(aggregateID string) (interface{}, int, error)

	// 事件清理
	CleanupEvents(beforeTime time.Time) error
	ArchiveEvents(beforeTime time.Time) error
}

// 辅助函数

// generateEventID 生成事件ID
func generateEventID() string {
	// 实现事件ID生成逻辑
	return fmt.Sprintf("pet_event_%d", time.Now().UnixNano())
}

// ValidateEvent 验证事件
func ValidateEvent(event PetEvent) error {
	if event.GetEventID() == "" {
		return fmt.Errorf("event ID cannot be empty")
	}
	if event.GetEventType() == "" {
		return fmt.Errorf("event type cannot be empty")
	}
	if event.GetAggregateID() == "" {
		return fmt.Errorf("aggregate ID cannot be empty")
	}
	if event.GetPlayerID() == "" {
		return fmt.Errorf("player ID cannot be empty")
	}
	if event.GetTimestamp().IsZero() {
		return fmt.Errorf("timestamp cannot be zero")
	}
	return nil
}

// SerializeEvent 序列化事件
func SerializeEvent(event PetEvent) ([]byte, error) {
	// 实现事件序列化逻辑
	return nil, nil
}

// DeserializeEvent 反序列化事件
func DeserializeEvent(data []byte, eventType string) (PetEvent, error) {
	// 实现事件反序列化逻辑
	return nil, nil
}
