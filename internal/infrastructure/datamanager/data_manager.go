package datamanager

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// UnitDefine 单位定义
type UnitDefine struct {
	ID        int32   `json:"id"`
	Name      string  `json:"name"`
	Type      int32   `json:"type"` // 1=玩家 2=怪物 3=NPC
	Level     int32   `json:"level"`
	MaxHP     int32   `json:"max_hp"`
	MaxMP     int32   `json:"max_mp"`
	STR       int32   `json:"str"`
	INT       int32   `json:"int"`
	AGI       int32   `json:"agi"`
	VIT       int32   `json:"vit"`
	SPR       int32   `json:"spr"`
	AD        int32   `json:"ad"`
	AP        int32   `json:"ap"`
	DEF       int32   `json:"def"`
	RES       int32   `json:"res"`
	SPD       int32   `json:"spd"`
	MoveSpeed float32 `json:"move_speed"`
	Skills    []int32 `json:"skills"`
	AIType    int32   `json:"ai_type"`
	NPCType   int32   `json:"npc_type"`
}

// SkillDefine 技能定义
type SkillDefine struct {
	ID         int32   `json:"id"`
	Name       string  `json:"name"`
	Type       int32   `json:"type"` // 技能类型
	BaseDamage int32   `json:"base_damage"`
	ScaleAD    float32 `json:"scale_ad"`
	ScaleAP    float32 `json:"scale_ap"`
	DamageType int32   `json:"damage_type"` // 1=物理 2=魔法 3=真实
	Cooldown   float32 `json:"cooldown"`
	CastTime   float32 `json:"cast_time"`
	Range      float32 `json:"range"`
	MPCost     int32   `json:"mp_cost"`
	TargetType int32   `json:"target_type"`
	BuffID     int32   `json:"buff_id"`
}

// ItemDefine 物品定义
type ItemDefine struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Type        int32  `json:"type"` // 1=消耗品 2=装备 3=材料
	Quality     int32  `json:"quality"`
	MaxStack    int32  `json:"max_stack"`
	Price       int32  `json:"price"`
	SellPrice   int32  `json:"sell_price"`
	EquipSlot   int32  `json:"equip_slot"`
	Description string `json:"description"`
}

// MapDefine 地图定义
type MapDefine struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	Width  int32  `json:"width"`
	Height int32  `json:"height"`
}

// QuestDefine 任务定义
type QuestDefine struct {
	ID          int32            `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Level       int32            `json:"level"`
	Objectives  []QuestObjective `json:"objectives"`
}

// QuestObjective 任务目标
type QuestObjective struct {
	Type     int32 `json:"type"`
	TargetID int32 `json:"target_id"`
	Required int32 `json:"required"`
}

// DataManager 数据管理器
type DataManager struct {
	mu sync.RWMutex

	unitDefines  map[int32]*UnitDefine
	skillDefines map[int32]*SkillDefine
	itemDefines  map[int32]*ItemDefine
	mapDefines   map[int32]*MapDefine
	questDefines map[int32]*QuestDefine
}

var instance *DataManager
var once sync.Once

// GetInstance 获取单例实例
func GetInstance() *DataManager {
	once.Do(func() {
		instance = &DataManager{
			unitDefines:  make(map[int32]*UnitDefine),
			skillDefines: make(map[int32]*SkillDefine),
			itemDefines:  make(map[int32]*ItemDefine),
			mapDefines:   make(map[int32]*MapDefine),
			questDefines: make(map[int32]*QuestDefine),
		}
	})
	return instance
}

// LoadAll 加载所有配置
func (dm *DataManager) LoadAll(configPath string) error {
	if err := dm.LoadUnits(configPath + "/units.json"); err != nil {
		return fmt.Errorf("load units failed: %w", err)
	}
	if err := dm.LoadSkills(configPath + "/skills.json"); err != nil {
		return fmt.Errorf("load skills failed: %w", err)
	}
	if err := dm.LoadItems(configPath + "/items.json"); err != nil {
		return fmt.Errorf("load items failed: %w", err)
	}
	if err := dm.LoadMaps(configPath + "/maps.json"); err != nil {
		return fmt.Errorf("load maps failed: %w", err)
	}
	if err := dm.LoadQuests(configPath + "/quests.json"); err != nil {
		return fmt.Errorf("load quests failed: %w", err)
	}
	return nil
}

// LoadUnits 加载单位配置
func (dm *DataManager) LoadUnits(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var units []*UnitDefine
	if err := json.Unmarshal(data, &units); err != nil {
		return err
	}

	dm.mu.Lock()
	defer dm.mu.Unlock()

	for _, unit := range units {
		dm.unitDefines[unit.ID] = unit
	}

	return nil
}

// LoadSkills 加载技能配置
func (dm *DataManager) LoadSkills(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var skills []*SkillDefine
	if err := json.Unmarshal(data, &skills); err != nil {
		return err
	}

	dm.mu.Lock()
	defer dm.mu.Unlock()

	for _, skill := range skills {
		dm.skillDefines[skill.ID] = skill
	}

	return nil
}

// LoadItems 加载物品配置
func (dm *DataManager) LoadItems(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var items []*ItemDefine
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}

	dm.mu.Lock()
	defer dm.mu.Unlock()

	for _, item := range items {
		dm.itemDefines[item.ID] = item
	}

	return nil
}

// LoadMaps 加载地图配置
func (dm *DataManager) LoadMaps(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var maps []*MapDefine
	if err := json.Unmarshal(data, &maps); err != nil {
		return err
	}

	dm.mu.Lock()
	defer dm.mu.Unlock()

	for _, m := range maps {
		dm.mapDefines[m.ID] = m
	}

	return nil
}

// LoadQuests 加载任务配置
func (dm *DataManager) LoadQuests(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var quests []*QuestDefine
	if err := json.Unmarshal(data, &quests); err != nil {
		return err
	}

	dm.mu.Lock()
	defer dm.mu.Unlock()

	for _, quest := range quests {
		dm.questDefines[quest.ID] = quest
	}

	return nil
}

// GetUnitDefine 获取单位定义
func (dm *DataManager) GetUnitDefine(id int32) *UnitDefine {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.unitDefines[id]
}

// GetSkillDefine 获取技能定义
func (dm *DataManager) GetSkillDefine(id int32) *SkillDefine {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.skillDefines[id]
}

// GetItemDefine 获取物品定义
func (dm *DataManager) GetItemDefine(id int32) *ItemDefine {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.itemDefines[id]
}

// GetMapDefine 获取地图定义
func (dm *DataManager) GetMapDefine(id int32) *MapDefine {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.mapDefines[id]
}

// GetQuestDefine 获取任务定义
func (dm *DataManager) GetQuestDefine(id int32) *QuestDefine {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.questDefines[id]
}

// GetUnit 获取单位定义（简短别名）
func (dm *DataManager) GetUnit(id int32) *UnitDefine {
	return dm.GetUnitDefine(id)
}

// GetSkill 获取技能定义（简短别名）
func (dm *DataManager) GetSkill(id int32) *SkillDefine {
	return dm.GetSkillDefine(id)
}

// GetItem 获取物品定义（简短别名）
func (dm *DataManager) GetItem(id int32) *ItemDefine {
	return dm.GetItemDefine(id)
}

// GetMap 获取地图定义（简短别名）
func (dm *DataManager) GetMap(id int32) *MapDefine {
	return dm.GetMapDefine(id)
}

// GetQuest 获取任务定义（简短别名）
func (dm *DataManager) GetQuest(id int32) *QuestDefine {
	return dm.GetQuestDefine(id)
}
