package player

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// PlayerName 玩家名称值对象
type PlayerName struct {
	value string
}

// NewPlayerName 创建玩家名称
func NewPlayerName(name string) (PlayerName, error) {
	if err := validatePlayerName(name); err != nil {
		return PlayerName{}, err
	}
	return PlayerName{value: strings.TrimSpace(name)}, nil
}

// String 返回字符串表示
func (n PlayerName) String() string {
	return n.value
}

// Equals 比较是否相等
func (n PlayerName) Equals(other PlayerName) bool {
	return n.value == other.value
}

// validatePlayerName 验证玩家名称
func validatePlayerName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("player name cannot be empty")
	}
	if len(name) < 2 {
		return errors.New("player name must be at least 2 characters")
	}
	if len(name) > 20 {
		return errors.New("player name cannot exceed 20 characters")
	}

	// 只允许字母、数字和下划线
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", name)
	if !matched {
		return errors.New("player name can only contain letters, numbers and underscores")
	}

	return nil
}

// Level 等级值对象
type Level struct {
	value int
}

// NewLevel 创建等级
func NewLevel(level int) (Level, error) {
	if level < 1 {
		return Level{}, errors.New("level must be at least 1")
	}
	if level > 100 {
		return Level{}, errors.New("level cannot exceed 100")
	}
	return Level{value: level}, nil
}

// Value 获取等级值
func (l Level) Value() int {
	return l.value
}

// String 返回字符串表示
func (l Level) String() string {
	return fmt.Sprintf("Level %d", l.value)
}

// Equals 比较是否相等
func (l Level) Equals(other Level) bool {
	return l.value == other.value
}

// CanLevelUp 是否可以升级
func (l Level) CanLevelUp() bool {
	return l.value < 100
}

// NextLevel 获取下一级
func (l Level) NextLevel() (Level, error) {
	if !l.CanLevelUp() {
		return Level{}, errors.New("already at max level")
	}
	return NewLevel(l.value + 1)
}

// Experience 经验值对象
type Experience struct {
	value int64
}

// NewExperience 创建经验
func NewExperience(exp int64) (Experience, error) {
	if exp < 0 {
		return Experience{}, errors.New("experience cannot be negative")
	}
	return Experience{value: exp}, nil
}

// Value 获取经验值
func (e Experience) Value() int64 {
	return e.value
}

// String 返回字符串表示
func (e Experience) String() string {
	return fmt.Sprintf("%d EXP", e.value)
}

// Add 增加经验
func (e Experience) Add(amount int64) (Experience, error) {
	if amount < 0 {
		return Experience{}, errors.New("cannot add negative experience")
	}
	return NewExperience(e.value + amount)
}

// Equals 比较是否相等
func (e Experience) Equals(other Experience) bool {
	return e.value == other.value
}

// HealthPoints 生命值对象
type HealthPoints struct {
	current int
	max     int
}

// NewHealthPoints 创建生命值
func NewHealthPoints(current, max int) (HealthPoints, error) {
	if max <= 0 {
		return HealthPoints{}, errors.New("max health must be positive")
	}
	if current < 0 {
		return HealthPoints{}, errors.New("current health cannot be negative")
	}
	if current > max {
		current = max
	}
	return HealthPoints{current: current, max: max}, nil
}

// Current 获取当前生命值
func (hp HealthPoints) Current() int {
	return hp.current
}

// Max 获取最大生命值
func (hp HealthPoints) Max() int {
	return hp.max
}

// Percentage 获取生命值百分比
func (hp HealthPoints) Percentage() float64 {
	if hp.max == 0 {
		return 0
	}
	return float64(hp.current) / float64(hp.max) * 100
}

// IsAlive 是否存活
func (hp HealthPoints) IsAlive() bool {
	return hp.current > 0
}

// IsFull 是否满血
func (hp HealthPoints) IsFull() bool {
	return hp.current == hp.max
}

// TakeDamage 受到伤害
func (hp HealthPoints) TakeDamage(damage int) HealthPoints {
	if damage < 0 {
		damage = 0
	}
	newCurrent := hp.current - damage
	if newCurrent < 0 {
		newCurrent = 0
	}
	return HealthPoints{current: newCurrent, max: hp.max}
}

// Heal 治疗
func (hp HealthPoints) Heal(amount int) HealthPoints {
	if amount < 0 {
		amount = 0
	}
	newCurrent := hp.current + amount
	if newCurrent > hp.max {
		newCurrent = hp.max
	}
	return HealthPoints{current: newCurrent, max: hp.max}
}

// String 返回字符串表示
func (hp HealthPoints) String() string {
	return fmt.Sprintf("%d/%d HP (%.1f%%)", hp.current, hp.max, hp.Percentage())
}

// ManaPoints 魔法值对象
type ManaPoints struct {
	current int
	max     int
}

// NewManaPoints 创建魔法值
func NewManaPoints(current, max int) (ManaPoints, error) {
	if max <= 0 {
		return ManaPoints{}, errors.New("max mana must be positive")
	}
	if current < 0 {
		return ManaPoints{}, errors.New("current mana cannot be negative")
	}
	if current > max {
		current = max
	}
	return ManaPoints{current: current, max: max}, nil
}

// Current 获取当前魔法值
func (mp ManaPoints) Current() int {
	return mp.current
}

// Max 获取最大魔法值
func (mp ManaPoints) Max() int {
	return mp.max
}

// Percentage 获取魔法值百分比
func (mp ManaPoints) Percentage() float64 {
	if mp.max == 0 {
		return 0
	}
	return float64(mp.current) / float64(mp.max) * 100
}

// HasEnough 是否有足够魔法值
func (mp ManaPoints) HasEnough(required int) bool {
	return mp.current >= required
}

// Consume 消耗魔法值
func (mp ManaPoints) Consume(amount int) (ManaPoints, error) {
	if amount < 0 {
		return mp, errors.New("cannot consume negative mana")
	}
	if mp.current < amount {
		return mp, errors.New("insufficient mana")
	}
	return ManaPoints{current: mp.current - amount, max: mp.max}, nil
}

// Restore 恢复魔法值
func (mp ManaPoints) Restore(amount int) ManaPoints {
	if amount < 0 {
		amount = 0
	}
	newCurrent := mp.current + amount
	if newCurrent > mp.max {
		newCurrent = mp.max
	}
	return ManaPoints{current: newCurrent, max: mp.max}
}

// String 返回字符串表示
func (mp ManaPoints) String() string {
	return fmt.Sprintf("%d/%d MP (%.1f%%)", mp.current, mp.max, mp.Percentage())
}
