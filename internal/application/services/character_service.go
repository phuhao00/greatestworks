package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"greatestworks/internal/domain/character"
	"greatestworks/internal/infrastructure/datamanager"
	"greatestworks/internal/infrastructure/persistence"
)

// CharacterService 角色服务
type CharacterService struct {
	characterRepo *persistence.CharacterRepository
	itemRepo      *persistence.ItemRepository
	questRepo     *persistence.QuestRepository
}

// NewCharacterService 创建角色服务
func NewCharacterService(
	characterRepo *persistence.CharacterRepository,
	itemRepo *persistence.ItemRepository,
	questRepo *persistence.QuestRepository,
) *CharacterService {
	return &CharacterService{
		characterRepo: characterRepo,
		itemRepo:      itemRepo,
		questRepo:     questRepo,
	}
}

// CreateCharacter 创建角色
func (s *CharacterService) CreateCharacter(ctx context.Context, userID int64, name string, race, class int32) (int64, error) {
	// 生成角色ID
	characterID := time.Now().UnixNano()

	// 获取角色初始配置
	unitDefine := datamanager.GetInstance().GetUnit(class)
	if unitDefine == nil {
		return 0, errors.New("invalid character class")
	}

	// 创建角色
	dbChar := &persistence.DbCharacter{
		CharacterID: characterID,
		UserID:      userID,
		Name:        name,
		Race:        race,
		Class:       class,
		Level:       1,
		Exp:         0,
		Gold:        1000, // 初始金币

		MapID:     1, // 默认地图
		PositionX: 0.0,
		PositionY: 0.0,
		PositionZ: 0.0,
		Direction: 0.0,

		HP:    unitDefine.MaxHP,
		MP:    unitDefine.MaxMP,
		MaxHP: unitDefine.MaxHP,
		MaxMP: unitDefine.MaxMP,

		STR: unitDefine.STR,
		INT: unitDefine.INT,
		AGI: unitDefine.AGI,
		VIT: unitDefine.VIT,
		SPR: unitDefine.SPR,

		AD:  unitDefine.AD,
		AP:  unitDefine.AP,
		DEF: unitDefine.DEF,
		RES: unitDefine.RES,
		SPD: unitDefine.SPD,

		CRI:     50,  // 默认暴击率 5%
		CRID:    150, // 默认暴击伤害 150%
		HitRate: 950, // 默认命中率 95%
		Dodge:   50,  // 默认闪避率 5%
	}

	if err := s.characterRepo.Create(ctx, dbChar); err != nil {
		return 0, fmt.Errorf("failed to create character: %w", err)
	}

	// TODO: 创建初始物品

	return characterID, nil
}

// GetCharacter 获取角色信息
func (s *CharacterService) GetCharacter(ctx context.Context, characterID int64) (*persistence.DbCharacter, error) {
	return s.characterRepo.FindByID(ctx, characterID)
}

// GetCharactersByUser 获取用户的所有角色
func (s *CharacterService) GetCharactersByUser(ctx context.Context, userID int64) ([]*persistence.DbCharacter, error) {
	return s.characterRepo.FindByUserID(ctx, userID)
}

// DeleteCharacter 删除角色
func (s *CharacterService) DeleteCharacter(ctx context.Context, characterID int64) error {
	return s.characterRepo.Delete(ctx, characterID)
}

// LoadCharacter 加载角色到内存（转换为领域对象）
func (s *CharacterService) LoadCharacter(ctx context.Context, characterID int64) (*character.Player, error) {
	dbChar, err := s.characterRepo.FindByID(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to load character: %w", err)
	}

	// 创建领域对象
	// 构造实体所需的位置信息与朝向（方向先用默认前向）
	pos := character.NewVector3(dbChar.PositionX, dbChar.PositionY, dbChar.PositionZ)
	dir := character.NewVector3(0, 0, 1)
	player := character.NewPlayer(
		character.EntityID(int32(dbChar.CharacterID)), // 实体ID（简化转换）
		dbChar.CharacterID,                            // 角色ID
		dbChar.UserID,                                 // 用户ID
		dbChar.Class,                                  // 使用职业作为unitID
		pos,
		dir,
		dbChar.Name,
		dbChar.Level,
	)

	// 设置属性
	player.SetLevel(dbChar.Level)

	// 设置基础属性到属性管理器
	am := player.GetAttributeManager()
	base := character.Attributes{
		MaxHP:       float32(dbChar.MaxHP),
		MaxMP:       float32(dbChar.MaxMP),
		HPRegen:     0,
		MPRegen:     0,
		AD:          float32(dbChar.AD),
		AP:          float32(dbChar.AP),
		Def:         float32(dbChar.DEF),
		MDef:        float32(dbChar.RES),
		Cri:         float32(dbChar.CRI) / 1000.0,     // 依据存储比例进行简单换算（占位）
		Crd:         float32(dbChar.CRID) / 100.0,     // 依据存储比例进行简单换算（占位）
		HitRate:     float32(dbChar.HitRate) / 1000.0, // 占位
		DodgeRate:   float32(dbChar.Dodge) / 1000.0,   // 占位
		Speed:       float32(dbChar.SPD),
		AttackSpeed: 1,
	}
	am.SetBase(base)

	// 设置当前HP/MP（以增量方式设置到目标数值）
	if dbChar.HP > 0 {
		player.ChangeHP(float32(dbChar.HP))
	}
	if dbChar.MP > 0 {
		player.ChangeMP(float32(dbChar.MP))
	}

	// 设置基础属性
	// 由于当前领域模型未包含STR/INT/AGI/VIT/SPR等细分属性，暂不映射这些字段

	// 加载物品
	items, err := s.itemRepo.FindByCharacterID(ctx, characterID)
	if err == nil {
		// TODO: 加载物品到背包
		_ = items
	}

	// 加载任务
	quests, err := s.questRepo.FindByCharacterID(ctx, characterID)
	if err == nil {
		// TODO: 加载任务到任务管理器
		_ = quests
	}

	return player, nil
}

// SaveCharacter 保存角色到数据库
func (s *CharacterService) SaveCharacter(ctx context.Context, player *character.Player) error {
	// 读取最终属性用于持久化
	attrs := player.GetAttributeManager().Final()

	dbChar := &persistence.DbCharacter{
		CharacterID: player.CharacterID(),
		Name:        player.GetName(),
		Race:        player.GetRace(),
		Class:       player.GetClass(),
		Level:       player.Level(),
		Exp:         player.GetExp(),
		Gold:        player.GetGold(),

		HP:    int32(player.HP()),
		MP:    int32(player.MP()),
		MaxHP: int32(attrs.MaxHP),
		MaxMP: int32(attrs.MaxMP),

		// 领域模型暂不包含以下基础属性的独立管理，使用0占位
		STR: 0,
		INT: 0,
		AGI: 0,
		VIT: 0,
		SPR: 0,

		AD:  int32(attrs.AD),
		AP:  int32(attrs.AP),
		DEF: int32(attrs.Def),
		RES: int32(attrs.MDef),
		SPD: int32(attrs.Speed),

		CRI:     int32(attrs.Cri * 1000),     // 与Load时的占位换算保持一致
		CRID:    int32(attrs.Crd * 100),      // 占位
		HitRate: int32(attrs.HitRate * 1000), // 占位
		Dodge:   int32(attrs.DodgeRate * 1000),
	}

	return s.characterRepo.Update(ctx, dbChar)
}

// UpdatePosition 更新角色位置
func (s *CharacterService) UpdatePosition(ctx context.Context, characterID int64, mapID int32, x, y, z, dir float32) error {
	return s.characterRepo.UpdatePosition(ctx, characterID, mapID, x, y, z, dir)
}

// UpdateLastLocation 更新角色上次位置（用于登出保存）
func (s *CharacterService) UpdateLastLocation(ctx context.Context, characterID int64, mapID int32, x, y, z float32) error {
	// 直接调用底层仓储更新位置，direction 传0即可
	return s.characterRepo.UpdatePosition(ctx, characterID, mapID, x, y, z, 0)
}

// AddExp 添加经验
func (s *CharacterService) AddExp(ctx context.Context, player *character.Player, exp int64) error {
	player.AddExp(exp)

	// 检查升级
	for player.CanLevelUp() {
		player.LevelUp()
	}

	return s.SaveCharacter(ctx, player)
}

// AddGold 添加金币
func (s *CharacterService) AddGold(ctx context.Context, player *character.Player, gold int64) error {
	player.AddGold(gold)
	return s.SaveCharacter(ctx, player)
}
