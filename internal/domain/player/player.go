// Package player 玩家领域
package player

import (
	"time"

	"github.com/google/uuid"
)

// PlayerID 玩家ID值对象
type PlayerID struct {
	value string
}

// NewPlayerID 创建新的玩家ID
func NewPlayerID() PlayerID {
	return PlayerID{value: uuid.New().String()}
}

// String 返回字符串表示
func (id PlayerID) String() string {
	return id.value
}

// PlayerIDFromString 从字符串创建PlayerID
func PlayerIDFromString(value string) PlayerID {
	return PlayerID{value: value}
}

// PlayerStatus 玩家状态枚举
type PlayerStatus int

const (
	PlayerStatusOffline PlayerStatus = iota
	PlayerStatusOnline
	PlayerStatusInBattle
	PlayerStatusInScene
)

// Player 玩家聚合根
type Player struct {
	id        PlayerID
	name      string
	level     int
	exp       int64
	status    PlayerStatus
	position  Position
	lastMapID int32 // 上次所在地图ID
	stats     PlayerStats
	createdAt time.Time
	updatedAt time.Time
	version   int64 // 乐观锁版本号
}

// Position 位置值对象
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// PlayerStats 玩家属性值对象
type PlayerStats struct {
	HP      int `json:"hp"`
	MaxHP   int `json:"max_hp"`
	MP      int `json:"mp"`
	MaxMP   int `json:"max_mp"`
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	Speed   int `json:"speed"`
}

// NewPlayer 创建新玩家
func NewPlayer(name string) *Player {
	now := time.Now()
	return &Player{
		id:        NewPlayerID(),
		name:      name,
		level:     1,
		exp:       0,
		status:    PlayerStatusOffline,
		position:  Position{X: 0, Y: 0, Z: 0},
		lastMapID: 1001, // 默认新手地图
		stats:     PlayerStats{HP: 100, MaxHP: 100, MP: 50, MaxMP: 50, Attack: 10, Defense: 5, Speed: 10},
		createdAt: now,
		updatedAt: now,
		version:   1,
	}
}

// ID 获取玩家ID
func (p *Player) ID() PlayerID {
	return p.id
}

// Name 获取玩家名称
func (p *Player) Name() string {
	return p.name
}

// Level 获取玩家等级
func (p *Player) Level() int {
	return p.level
}

// Status 获取玩家状态
func (p *Player) Status() PlayerStatus {
	return p.status
}

// Position 获取玩家位置
func (p *Player) GetPosition() Position {
	return p.position
}

// LastMapID 获取上次所在地图ID
func (p *Player) LastMapID() int32 {
	return p.lastMapID
}

// Stats 获取玩家属性
func (p *Player) Stats() PlayerStats {
	return p.stats
}

// SetOnline 设置玩家上线
func (p *Player) SetOnline() {
	p.status = PlayerStatusOnline
	p.updatedAt = time.Now()
	p.version++
}

// SetOffline 设置玩家下线
func (p *Player) SetOffline() {
	p.status = PlayerStatusOffline
	p.updatedAt = time.Now()
	p.version++
}

// MoveTo 移动到指定位置
func (p *Player) MoveTo(pos Position) error {
	if p.status == PlayerStatusOffline {
		return ErrPlayerOffline
	}
	p.position = pos
	p.updatedAt = time.Now()
	p.version++
	return nil
}

// SetLastLocation 设置上次位置（用于登出保存）
func (p *Player) SetLastLocation(mapID int32, pos Position) {
	p.lastMapID = mapID
	p.position = pos
	p.updatedAt = time.Now()
	p.version++
}

// GainExp 获得经验值
func (p *Player) GainExp(exp int64) {
	p.exp += exp
	// 检查是否升级
	for p.exp >= p.getExpForNextLevel() {
		p.levelUp()
	}
	p.updatedAt = time.Now()
	p.version++
}

// levelUp 升级
func (p *Player) levelUp() {
	p.level++
	// 升级时增加属性
	p.stats.MaxHP += 20
	p.stats.HP = p.stats.MaxHP
	p.stats.MaxMP += 10
	p.stats.MP = p.stats.MaxMP
	p.stats.Attack += 2
	p.stats.Defense += 1
}

// getExpForNextLevel 获取下一级所需经验
func (p *Player) getExpForNextLevel() int64 {
	return int64(p.level * 100)
}

// TakeDamage 受到伤害
func (p *Player) TakeDamage(damage int) bool {
	if damage <= 0 {
		return false
	}

	actualDamage := damage - p.stats.Defense
	if actualDamage < 1 {
		actualDamage = 1
	}

	p.stats.HP -= actualDamage
	if p.stats.HP < 0 {
		p.stats.HP = 0
	}

	p.updatedAt = time.Now()
	p.version++

	return p.stats.HP == 0 // 返回是否死亡
}

// Heal 治疗
func (p *Player) Heal(amount int) {
	p.stats.HP += amount
	if p.stats.HP > p.stats.MaxHP {
		p.stats.HP = p.stats.MaxHP
	}
	p.updatedAt = time.Now()
	p.version++
}

// IsAlive 是否存活
func (p *Player) IsAlive() bool {
	return p.stats.HP > 0
}

// CreatedAt 获取创建时间
func (p *Player) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt 获取更新时间
func (p *Player) UpdatedAt() time.Time {
	return p.updatedAt
}

// Version 获取版本号
func (p *Player) Version() int64 {
	return p.version
}

// Exp 获取经验值
func (p *Player) Exp() int64 {
	return p.exp
}

// ReconstructPlayer 从持久化数据重建玩家聚合根
func ReconstructPlayer(id PlayerID, name string, level int, exp int64, status PlayerStatus, position Position, lastMapID int32, stats PlayerStats, createdAt, updatedAt time.Time, version int64) *Player {
	return &Player{
		id:        id,
		name:      name,
		level:     level,
		exp:       exp,
		status:    status,
		position:  position,
		lastMapID: lastMapID,
		stats:     stats,
		createdAt: createdAt,
		updatedAt: updatedAt,
		version:   version,
	}
}
