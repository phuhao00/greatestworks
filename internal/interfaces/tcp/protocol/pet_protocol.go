package protocol

// 宠物相关消息类型
const (
	// 请求消息类型
	MsgTypeCreatePet     = "create_pet"
	MsgTypeFeedPet       = "feed_pet"
	MsgTypeTrainPet      = "train_pet"
	MsgTypeGetPet        = "get_pet"
	MsgTypeGetPlayerPets = "get_player_pets"
	MsgTypeEvolvePet     = "evolve_pet"
	MsgTypeEquipPetSkin  = "equip_pet_skin"
	MsgTypeSynthesizePet = "synthesize_pet"
	
	// 响应消息类型
	MsgTypeCreatePetResponse     = "create_pet_response"
	MsgTypeFeedPetResponse       = "feed_pet_response"
	MsgTypeTrainPetResponse      = "train_pet_response"
	MsgTypeGetPetResponse        = "get_pet_response"
	MsgTypeGetPlayerPetsResponse = "get_player_pets_response"
	MsgTypeEvolvePetResponse     = "evolve_pet_response"
	MsgTypeEquipPetSkinResponse  = "equip_pet_skin_response"
	MsgTypeSynthesizePetResponse = "synthesize_pet_response"
	
	// 通知消息类型
	MsgTypePetLevelUp      = "pet_level_up"
	MsgTypePetMoodChanged  = "pet_mood_changed"
	MsgTypePetEvolved      = "pet_evolved"
	MsgTypePetSkinEquipped = "pet_skin_equipped"
)

// CreatePetRequest 创建宠物请求
type CreatePetRequest struct {
	PlayerID  uint64 `json:"player_id"`
	SpeciesID string `json:"species_id"`
	Name      string `json:"name"`
	Rarity    string `json:"rarity,omitempty"`
	Quality   string `json:"quality,omitempty"`
}

// CreatePetResponse 创建宠物响应
type CreatePetResponse struct {
	PetID     string `json:"pet_id"`
	PlayerID  uint64 `json:"player_id"`
	SpeciesID string `json:"species_id"`
	Name      string `json:"name"`
	Level     int32  `json:"level"`
	Exp       int64  `json:"exp"`
	MaxExp    int64  `json:"max_exp"`
	Rarity    string `json:"rarity"`
	Quality   string `json:"quality"`
	CreatedAt int64  `json:"created_at"`
}

// FeedPetRequest 喂养宠物请求
type FeedPetRequest struct {
	PetID    string `json:"pet_id"`
	FoodID   string `json:"food_id"`
	Quantity int32  `json:"quantity"`
}

// FeedPetResponse 喂养宠物响应
type FeedPetResponse struct {
	PetID        string `json:"pet_id"`
	OldHunger    int32  `json:"old_hunger"`
	NewHunger    int32  `json:"new_hunger"`
	OldHappiness int32  `json:"old_happiness"`
	NewHappiness int32  `json:"new_happiness"`
	ExpGained    int64  `json:"exp_gained"`
	LevelUp      bool   `json:"level_up"`
	FedAt        int64  `json:"fed_at"`
}

// TrainPetRequest 训练宠物请求
type TrainPetRequest struct {
	PetID        string `json:"pet_id"`
	TrainingType string `json:"training_type"`
	Duration     int32  `json:"duration"`
	Intensity    string `json:"intensity,omitempty"`
}

// TrainPetResponse 训练宠物响应
type TrainPetResponse struct {
	PetID            string            `json:"pet_id"`
	TrainingType     string            `json:"training_type"`
	ExpGained        int64             `json:"exp_gained"`
	AttributeChanges map[string]int64  `json:"attribute_changes"`
	EnergyConsumed   int32             `json:"energy_consumed"`
	LevelUp          bool              `json:"level_up"`
	SkillsLearned    []string          `json:"skills_learned,omitempty"`
	TrainedAt        int64             `json:"trained_at"`
}

// GetPetRequest 获取宠物请求
type GetPetRequest struct {
	PetID string `json:"pet_id"`
}

// GetPetResponse 获取宠物响应
type GetPetResponse struct {
	Found bool     `json:"found"`
	Pet   *PetInfo `json:"pet,omitempty"`
}

// GetPlayerPetsRequest 获取玩家宠物列表请求
type GetPlayerPetsRequest struct {
	PlayerID uint64 `json:"player_id"`
	Rarity   string `json:"rarity,omitempty"`
	Quality  string `json:"quality,omitempty"`
	MinLevel int32  `json:"min_level,omitempty"`
	MaxLevel int32  `json:"max_level,omitempty"`
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
}

// GetPlayerPetsResponse 获取玩家宠物列表响应
type GetPlayerPetsResponse struct {
	PlayerID   uint64     `json:"player_id"`
	Pets       []*PetInfo `json:"pets"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalPages int64      `json:"total_pages"`
}

// EvolvePetRequest 宠物进化请求
type EvolvePetRequest struct {
	PetID         string            `json:"pet_id"`
	TargetSpecies string            `json:"target_species"`
	Materials     map[string]int32  `json:"materials,omitempty"`
}

// EvolvePetResponse 宠物进化响应
type EvolvePetResponse struct {
	PetID          string            `json:"pet_id"`
	OldSpecies     string            `json:"old_species"`
	NewSpecies     string            `json:"new_species"`
	OldRarity      string            `json:"old_rarity"`
	NewRarity      string            `json:"new_rarity"`
	AttributeBonus map[string]int64  `json:"attribute_bonus"`
	NewSkills      []string          `json:"new_skills,omitempty"`
	MaterialsUsed  map[string]int32  `json:"materials_used"`
	EvolvedAt      int64             `json:"evolved_at"`
}

// EquipPetSkinRequest 装备宠物皮肤请求
type EquipPetSkinRequest struct {
	PetID  string `json:"pet_id"`
	SkinID string `json:"skin_id"`
}

// EquipPetSkinResponse 装备宠物皮肤响应
type EquipPetSkinResponse struct {
	PetID         string             `json:"pet_id"`
	OldSkinID     string             `json:"old_skin_id,omitempty"`
	NewSkinID     string             `json:"new_skin_id"`
	EffectChanges map[string]float64 `json:"effect_changes,omitempty"`
	EquippedAt    int64              `json:"equipped_at"`
}

// SynthesizePetRequest 宠物合成请求
type SynthesizePetRequest struct {
	PlayerID   uint64 `json:"player_id"`
	FragmentID string `json:"fragment_id"`
	Quantity   int32  `json:"quantity"`
}

// SynthesizePetResponse 宠物合成响应
type SynthesizePetResponse struct {
	PlayerID      uint64 `json:"player_id"`
	FragmentID    string `json:"fragment_id"`
	QuantityUsed  int32  `json:"quantity_used"`
	PetID         string `json:"pet_id,omitempty"`
	SpeciesID     string `json:"species_id,omitempty"`
	Rarity        string `json:"rarity,omitempty"`
	Success       bool   `json:"success"`
	SynthesizedAt int64  `json:"synthesized_at"`
}

// PetInfo 宠物信息
type PetInfo struct {
	PetID        string            `json:"pet_id"`
	PlayerID     uint64            `json:"player_id"`
	SpeciesID    string            `json:"species_id"`
	Name         string            `json:"name"`
	Level        int32             `json:"level"`
	Exp          int64             `json:"exp"`
	MaxExp       int64             `json:"max_exp"`
	Rarity       string            `json:"rarity"`
	Quality      string            `json:"quality"`
	Attributes   map[string]int64  `json:"attributes"`
	Skills       []string          `json:"skills"`
	EquippedSkin string            `json:"equipped_skin,omitempty"`
	Mood         string            `json:"mood"`
	Hunger       int32             `json:"hunger"`
	Energy       int32             `json:"energy"`
	Health       int32             `json:"health"`
	Happiness    int32             `json:"happiness"`
	IsActive     bool              `json:"is_active"`
	LastFedAt    int64             `json:"last_fed_at"`
	LastPlayedAt int64             `json:"last_played_at"`
	CreatedAt    int64             `json:"created_at"`
	UpdatedAt    int64             `json:"updated_at"`
}

// 通知消息结构

// PetLevelUpNotification 宠物升级通知
type PetLevelUpNotification struct {
	PetID        string            `json:"pet_id"`
	PlayerID     uint64            `json:"player_id"`
	OldLevel     int32             `json:"old_level"`
	NewLevel     int32             `json:"new_level"`
	NewMaxExp    int64             `json:"new_max_exp"`
	AttributeGain map[string]int64 `json:"attribute_gain,omitempty"`
	SkillsLearned []string         `json:"skills_learned,omitempty"`
	LevelUpAt    int64             `json:"level_up_at"`
}

// PetMoodChangedNotification 宠物心情变化通知
type PetMoodChangedNotification struct {
	PetID     string `json:"pet_id"`
	PlayerID  uint64 `json:"player_id"`
	OldMood   string `json:"old_mood"`
	NewMood   string `json:"new_mood"`
	Reason    string `json:"reason"`
	ChangedAt int64  `json:"changed_at"`
}

// PetEvolvedNotification 宠物进化通知
type PetEvolvedNotification struct {
	PetID         string            `json:"pet_id"`
	PlayerID      uint64            `json:"player_id"`
	OldSpecies    string            `json:"old_species"`
	NewSpecies    string            `json:"new_species"`
	OldRarity     string            `json:"old_rarity"`
	NewRarity     string            `json:"new_rarity"`
	AttributeGain map[string]int64  `json:"attribute_gain"`
	NewSkills     []string          `json:"new_skills,omitempty"`
	EvolvedAt     int64             `json:"evolved_at"`
}

// PetSkinEquippedNotification 宠物皮肤装备通知
type PetSkinEquippedNotification struct {
	PetID         string             `json:"pet_id"`
	PlayerID      uint64             `json:"player_id"`
	OldSkinID     string             `json:"old_skin_id,omitempty"`
	NewSkinID     string             `json:"new_skin_id"`
	SkinName      string             `json:"skin_name"`
	EffectChanges map[string]float64 `json:"effect_changes,omitempty"`
	EquippedAt    int64              `json:"equipped_at"`
}

// PetFragmentInfo 宠物碎片信息
type PetFragmentInfo struct {
	FragmentID string `json:"fragment_id"`
	PlayerID   uint64 `json:"player_id"`
	SpeciesID  string `json:"species_id"`
	Quantity   int32  `json:"quantity"`
	Required   int32  `json:"required"`
	Source     string `json:"source"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

// PetSkinInfo 宠物皮肤信息
type PetSkinInfo struct {
	SkinID     string             `json:"skin_id"`
	PlayerID   uint64             `json:"player_id"`
	SpeciesID  string             `json:"species_id"`
	Name       string             `json:"name"`
	Rarity     string             `json:"rarity"`
	Effects    map[string]float64 `json:"effects"`
	IsUnlocked bool               `json:"is_unlocked"`
	UnlockedAt int64              `json:"unlocked_at,omitempty"`
	CreatedAt  int64              `json:"created_at"`
}

// PetBondInfo 宠物羁绊信息
type PetBondInfo struct {
	BondID      string            `json:"bond_id"`
	PlayerID    uint64            `json:"player_id"`
	PetIDs      []string          `json:"pet_ids"`
	BondType    string            `json:"bond_type"`
	Level       int32             `json:"level"`
	Exp         int64             `json:"exp"`
	MaxExp      int64             `json:"max_exp"`
	Effects     map[string]float64 `json:"effects"`
	IsActive    bool              `json:"is_active"`
	ActivatedAt int64             `json:"activated_at,omitempty"`
	CreatedAt   int64             `json:"created_at"`
	UpdatedAt   int64             `json:"updated_at"`
}

// PetPictorialInfo 宠物图鉴信息
type PetPictorialInfo struct {
	PictorialID   string `json:"pictorial_id"`
	PlayerID      uint64 `json:"player_id"`
	SpeciesID     string `json:"species_id"`
	IsUnlocked    bool   `json:"is_unlocked"`
	FirstSeenAt   int64  `json:"first_seen_at,omitempty"`
	FirstOwnedAt  int64  `json:"first_owned_at,omitempty"`
	TotalOwned    int32  `json:"total_owned"`
	HighestLevel  int32  `json:"highest_level"`
	HighestRarity string `json:"highest_rarity"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}