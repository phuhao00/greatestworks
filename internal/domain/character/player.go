package character

import (
	"context"
	"fmt"
)

// Player 玩家实体 - 聚合根
type Player struct {
	*Actor // 继承Actor

	// 用户关联
	userID int64 // 所属用户ID

	// 角色数据
	characterID int64 // 角色ID（数据库主键）
	exp         int32 // 经验值
	gold        int64 // 金币

	// 背包系统（聚合）
	inventory *Inventory

	// 任务系统（聚合）
	taskManager *TaskManager

	// 对话系统
	dialogueManager *DialogueManager

	// 当前交互的NPC
	interactingNPC *NPC
}

// NewPlayer 创建玩家（工厂方法）
func NewPlayer(
	entityID EntityID,
	characterID int64,
	userID int64,
	unitID int32,
	position Vector3,
	direction Vector3,
	name string,
	level int32,
) *Player {
	actor := NewActor(
		entityID,
		EntityTypePlayer,
		unitID,
		position,
		direction,
		name,
		level,
	)

	player := &Player{
		Actor:       actor,
		userID:      userID,
		characterID: characterID,
	}

	// 初始化子系统
	player.inventory = NewInventory(player)
	player.taskManager = NewTaskManager(player)
	player.dialogueManager = NewDialogueManager(player)

	return player
}

// ========== 身份信息 ==========

// UserID 获取用户ID
func (p *Player) UserID() int64 {
	return p.userID
}

// CharacterID 获取角色ID
func (p *Player) CharacterID() int64 {
	return p.characterID
}

// ========== 角色属性 ==========

// Exp 获取经验值
func (p *Player) Exp() int32 {
	return p.exp
}

// Gold 获取金币
func (p *Player) Gold() int64 {
	return p.gold
}

// GetName 获取角色名称
func (p *Player) GetName() string {
	return p.Name()
}

// GetRace 获取种族
func (p *Player) GetRace() int32 {
	// 从UnitID获取种族信息，这里简化返回0
	// 实际应该从DataManager中的UnitDefine获取
	return 0
}

// GetClass 获取职业
func (p *Player) GetClass() int32 {
	return p.UnitID()
}

// GetExp 获取经验值
func (p *Player) GetExp() int64 {
	return int64(p.exp)
}

// GetGold 获取金币
func (p *Player) GetGold() int64 {
	return p.gold
}

// AddExp 添加经验值
func (p *Player) AddExp(amount int64) {
	p.exp += int32(amount)

	// 检查升级
	for p.CanLevelUp() {
		p.LevelUp()
	}
}

// AddGold 添加金币
func (p *Player) AddGold(amount int64) {
	p.gold += amount
	if p.gold < 0 {
		p.gold = 0
	}
}

// CanLevelUp 是否可以升级
func (p *Player) CanLevelUp() bool {
	currentLevel := p.Level()
	if currentLevel >= 100 { // 最大等级100
		return false
	}

	// 简化的升级经验公式: exp >= level * 1000
	requiredExp := int32(currentLevel * 1000)
	return p.exp >= requiredExp
}

// LevelUp 升级
func (p *Player) LevelUp() {
	currentLevel := p.Level()
	newLevel := currentLevel + 1

	// 设置新等级
	p.SetLevel(newLevel)

	// 扣除升级所需经验
	requiredExp := int32(currentLevel * 1000)
	p.exp -= requiredExp

	// 恢复生命和魔法 - 使用ChangeHP/ChangeMP设置为最大值
	attrs := p.Actor.GetAttributeManager().Final()
	maxHP := attrs.MaxHP
	maxMP := attrs.MaxMP
	currentHP := p.Actor.HP()
	currentMP := p.Actor.MP()
	p.Actor.ChangeHP(maxHP - currentHP)
	p.Actor.ChangeMP(maxMP - currentMP)

	// TODO: 触发升级事件
	// p.PublishEvent(&PlayerLevelUpEvent{...})
}

// ChangeExp 改变经验值
func (p *Player) ChangeExp(amount int32) {
	p.exp += amount

	// 处理升级逻辑
	for p.CanLevelUp() {
		p.LevelUp()
	}

	// TODO: 同步属性变化
	// p.syncAttributeEntry(AttributeTypeExp, p.exp)
}

// ChangeGold 改变金币
func (p *Player) ChangeGold(amount int64) {
	p.gold += amount
	if p.gold < 0 {
		p.gold = 0
	}

	// TODO: 同步属性变化
	// p.syncAttributeEntry(AttributeTypeGold, int32(p.gold))
}

// ========== 子系统访问 ==========

// GetInventory 获取背包
func (p *Player) GetInventory() *Inventory {
	return p.inventory
}

// GetTaskManager 获取任务管理器
func (p *Player) GetTaskManager() *TaskManager {
	return p.taskManager
}

// GetDialogueManager 获取对话管理器
func (p *Player) GetDialogueManager() *DialogueManager {
	return p.dialogueManager
}

// ========== NPC交互 ==========

// SetInteractingNPC 设置当前交互的NPC
func (p *Player) SetInteractingNPC(npc *NPC) {
	p.interactingNPC = npc
}

// GetInteractingNPC 获取当前交互的NPC
func (p *Player) GetInteractingNPC() *NPC {
	return p.interactingNPC
}

// ========== 生命周期 ==========

// Start 初始化玩家
func (p *Player) Start(ctx context.Context) error {
	// 调用Actor的Start
	if err := p.Actor.Start(ctx); err != nil {
		return err
	}

	// 初始化背包
	if err := p.inventory.Start(ctx); err != nil {
		return fmt.Errorf("inventory start failed: %w", err)
	}

	// 初始化任务管理器
	if err := p.taskManager.Start(ctx); err != nil {
		return fmt.Errorf("taskManager start failed: %w", err)
	}

	return nil
}

// Revive 玩家复活
func (p *Player) Revive(ctx context.Context) error {
	// 调用Actor的复活逻辑
	if err := p.Actor.Revive(ctx); err != nil {
		return err
	}

	// 玩家特有的复活逻辑（比如扣除经验等）
	// TODO: 实现玩家复活逻辑

	return nil
}

// String 字符串表示
func (p *Player) String() string {
	return fmt.Sprintf("Player:\"%s(%d)\"[User:%d,Char:%d]",
		p.Name(), p.ID(), p.userID, p.characterID)
}

// ========== 占位符子系统 ==========

// Inventory 背包（占位符）
type Inventory struct {
	owner *Player
}

func NewInventory(owner *Player) *Inventory {
	return &Inventory{owner: owner}
}

func (i *Inventory) Start(ctx context.Context) error {
	return nil
}

// TaskManager 任务管理器（占位符）
type TaskManager struct {
	owner *Player
}

func NewTaskManager(owner *Player) *TaskManager {
	return &TaskManager{owner: owner}
}

func (tm *TaskManager) Start(ctx context.Context) error {
	return nil
}

// DialogueManager 对话管理器（占位符）
type DialogueManager struct {
	owner *Player
}

func NewDialogueManager(owner *Player) *DialogueManager {
	return &DialogueManager{owner: owner}
}
