// Package battle 战斗领域
package battle

import (
	"time"
	"github.com/google/uuid"
	"greatestworks/internal/domain/player"
)

// BattleID 战斗ID值对象
type BattleID struct {
	value string
}

// NewBattleID 创建新的战斗ID
func NewBattleID() BattleID {
	return BattleID{value: uuid.New().String()}
}

// String 返回字符串表示
func (id BattleID) String() string {
	return id.value
}

// BattleStatus 战斗状态枚举
type BattleStatus int

const (
	BattleStatusWaiting BattleStatus = iota
	BattleStatusInProgress
	BattleStatusFinished
	BattleStatusCancelled
)

// BattleType 战斗类型枚举
type BattleType int

const (
	BattleTypePvP BattleType = iota
	BattleTypePvE
	BattleTypeTeamPvP
	BattleTypeRaid
)

// Battle 战斗聚合根
type Battle struct {
	id          BattleID
	battleType  BattleType
	status      BattleStatus
	participants []*BattleParticipant
	rounds      []*BattleRound
	winner      *player.PlayerID
	startTime   time.Time
	endTime     *time.Time
	createdAt   time.Time
	updatedAt   time.Time
	version     int64
}

// BattleParticipant 战斗参与者
type BattleParticipant struct {
	PlayerID     player.PlayerID `json:"player_id"`
	Team         int             `json:"team"`
	CurrentHP    int             `json:"current_hp"`
	CurrentMP    int             `json:"current_mp"`
	IsAlive      bool            `json:"is_alive"`
	DamageDealt  int             `json:"damage_dealt"`
	DamageTaken  int             `json:"damage_taken"`
	JoinedAt     time.Time       `json:"joined_at"`
}

// BattleRound 战斗回合
type BattleRound struct {
	RoundNumber int                `json:"round_number"`
	Actions     []*BattleAction    `json:"actions"`
	StartTime   time.Time          `json:"start_time"`
	EndTime     *time.Time         `json:"end_time"`
}

// BattleAction 战斗行动
type BattleAction struct {
	ActionID    string           `json:"action_id"`
	ActorID     player.PlayerID  `json:"actor_id"`
	TargetID    *player.PlayerID `json:"target_id,omitempty"`
	ActionType  ActionType       `json:"action_type"`
	SkillID     *string          `json:"skill_id,omitempty"`
	Damage      int              `json:"damage"`
	Healing     int              `json:"healing"`
	Critical    bool             `json:"critical"`
	Timestamp   time.Time        `json:"timestamp"`
}

// ActionType 行动类型枚举
type ActionType int

const (
	ActionTypeAttack ActionType = iota
	ActionTypeSkill
	ActionTypeDefend
	ActionTypeHeal
	ActionTypeEscape
)

// NewBattle 创建新战斗
func NewBattle(battleType BattleType) *Battle {
	now := time.Now()
	return &Battle{
		id:          NewBattleID(),
		battleType:  battleType,
		status:      BattleStatusWaiting,
		participants: make([]*BattleParticipant, 0),
		rounds:      make([]*BattleRound, 0),
		createdAt:   now,
		updatedAt:   now,
		version:     1,
	}
}

// ID 获取战斗ID
func (b *Battle) ID() BattleID {
	return b.id
}

// Status 获取战斗状态
func (b *Battle) Status() BattleStatus {
	return b.status
}

// BattleType 获取战斗类型
func (b *Battle) GetBattleType() BattleType {
	return b.battleType
}

// Participants 获取参与者
func (b *Battle) Participants() []*BattleParticipant {
	return b.participants
}

// AddParticipant 添加参与者
func (b *Battle) AddParticipant(playerID player.PlayerID, team int, hp, mp int) error {
	if b.status != BattleStatusWaiting {
		return ErrBattleAlreadyStarted
	}
	
	// 检查玩家是否已经参与
	for _, p := range b.participants {
		if p.PlayerID == playerID {
			return ErrPlayerAlreadyInBattle
		}
	}
	
	participant := &BattleParticipant{
		PlayerID:    playerID,
		Team:        team,
		CurrentHP:   hp,
		CurrentMP:   mp,
		IsAlive:     true,
		DamageDealt: 0,
		DamageTaken: 0,
		JoinedAt:    time.Now(),
	}
	
	b.participants = append(b.participants, participant)
	b.updatedAt = time.Now()
	b.version++
	
	return nil
}

// Start 开始战斗
func (b *Battle) Start() error {
	if b.status != BattleStatusWaiting {
		return ErrBattleAlreadyStarted
	}
	
	if len(b.participants) < 2 {
		return ErrInsufficientParticipants
	}
	
	b.status = BattleStatusInProgress
	b.startTime = time.Now()
	b.updatedAt = time.Now()
	b.version++
	
	return nil
}

// ExecuteAction 执行战斗行动
func (b *Battle) ExecuteAction(actorID player.PlayerID, targetID *player.PlayerID, actionType ActionType, skillID *string) (*BattleAction, error) {
	if b.status != BattleStatusInProgress {
		return nil, ErrBattleNotInProgress
	}
	
	// 查找行动者
	actor := b.findParticipant(actorID)
	if actor == nil {
		return nil, ErrPlayerNotInBattle
	}
	
	if !actor.IsAlive {
		return nil, ErrPlayerDead
	}
	
	// 创建行动
	action := &BattleAction{
		ActionID:   uuid.New().String(),
		ActorID:    actorID,
		TargetID:   targetID,
		ActionType: actionType,
		SkillID:    skillID,
		Timestamp:  time.Now(),
	}
	
	// 执行行动逻辑
	switch actionType {
	case ActionTypeAttack:
		b.executeAttack(action, actor)
	case ActionTypeDefend:
		b.executeDefend(action, actor)
	case ActionTypeHeal:
		b.executeHeal(action, actor)
	}
	
	// 添加到当前回合
	b.addActionToCurrentRound(action)
	
	// 检查战斗是否结束
	b.checkBattleEnd()
	
	b.updatedAt = time.Now()
	b.version++
	
	return action, nil
}

// executeAttack 执行攻击
func (b *Battle) executeAttack(action *BattleAction, actor *BattleParticipant) {
	if action.TargetID == nil {
		return
	}
	
	target := b.findParticipant(*action.TargetID)
	if target == nil || !target.IsAlive {
		return
	}
	
	// 计算伤害（简化版本）
	baseDamage := 20 // 基础攻击力
	damage := baseDamage
	
	// 暴击判断
	if b.rollCritical() {
		damage *= 2
		action.Critical = true
	}
	
	action.Damage = damage
	target.CurrentHP -= damage
	target.DamageTaken += damage
	actor.DamageDealt += damage
	
	if target.CurrentHP <= 0 {
		target.CurrentHP = 0
		target.IsAlive = false
	}
}

// executeDefend 执行防御
func (b *Battle) executeDefend(action *BattleAction, actor *BattleParticipant) {
	// 防御状态，下次受到伤害减半（简化实现）
}

// executeHeal 执行治疗
func (b *Battle) executeHeal(action *BattleAction, actor *BattleParticipant) {
	healAmount := 30 // 基础治疗量
	actor.CurrentHP += healAmount
	action.Healing = healAmount
}

// rollCritical 暴击判断
func (b *Battle) rollCritical() bool {
	// 简化版本：20%暴击率
	return time.Now().UnixNano()%5 == 0
}

// findParticipant 查找参与者
func (b *Battle) findParticipant(playerID player.PlayerID) *BattleParticipant {
	for _, p := range b.participants {
		if p.PlayerID == playerID {
			return p
		}
	}
	return nil
}

// addActionToCurrentRound 添加行动到当前回合
func (b *Battle) addActionToCurrentRound(action *BattleAction) {
	if len(b.rounds) == 0 {
		// 创建第一个回合
		round := &BattleRound{
			RoundNumber: 1,
			Actions:     make([]*BattleAction, 0),
			StartTime:   time.Now(),
		}
		b.rounds = append(b.rounds, round)
	}
	
	currentRound := b.rounds[len(b.rounds)-1]
	currentRound.Actions = append(currentRound.Actions, action)
}

// checkBattleEnd 检查战斗是否结束
func (b *Battle) checkBattleEnd() {
	// 统计各队伍存活人数
	teamAlive := make(map[int]int)
	for _, p := range b.participants {
		if p.IsAlive {
			teamAlive[p.Team]++
		}
	}
	
	// 如果只有一个队伍有存活者，战斗结束
	aliveTeams := 0
	winnerTeam := -1
	for team, count := range teamAlive {
		if count > 0 {
			aliveTeams++
			winnerTeam = team
		}
	}
	
	if aliveTeams <= 1 {
		b.endBattle(winnerTeam)
	}
}

// endBattle 结束战斗
func (b *Battle) endBattle(winnerTeam int) {
	b.status = BattleStatusFinished
	now := time.Now()
	b.endTime = &now
	
	// 设置获胜者（简化版本：取获胜队伍第一个存活玩家）
	for _, p := range b.participants {
		if p.Team == winnerTeam && p.IsAlive {
			b.winner = &p.PlayerID
			break
		}
	}
	
	// 结束当前回合
	if len(b.rounds) > 0 {
		currentRound := b.rounds[len(b.rounds)-1]
		if currentRound.EndTime == nil {
			currentRound.EndTime = &now
		}
	}
}

// Winner 获取获胜者
func (b *Battle) Winner() *player.PlayerID {
	return b.winner
}

// IsFinished 是否已结束
func (b *Battle) IsFinished() bool {
	return b.status == BattleStatusFinished || b.status == BattleStatusCancelled
}

// Version 获取版本号
func (b *Battle) Version() int64 {
	return b.version
}